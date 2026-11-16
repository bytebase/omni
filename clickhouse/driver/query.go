package driver

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkt"
	"github.com/shopspring/decimal"

	"github.com/bytebase/omni/driver"
)

const aggregateFunctionErrMarker = `unsupported column type "AggregateFunction(`

var (
	_ driver.Querier = (*Driver)(nil)
	_ driver.Session = (*session)(nil)
)

// session is one server session, held on one connection so that what a
// statement sets holds for the next.
type session struct {
	conn *sql.Conn
}

// Session returns one server session, held until Close.
func (d *Driver) Session(ctx context.Context) (driver.Session, error) {
	conn, err := d.db.Conn(ctx)
	if err != nil {
		return nil, err
	}
	return &session{conn: conn}, nil
}

// Close releases the session's connection.
func (s *session) Close() error {
	return s.conn.Close()
}

// Query runs each statement of script in order and returns one Result per
// statement run. A statement that fails ends the script.
func (s *session) Query(ctx context.Context, script string, opts driver.QueryOptions) ([]*driver.Result, error) {
	statements, err := splitSQL(script)
	if err != nil {
		return nil, err
	}
	statements = filterEmptyStatements(statements)
	if len(statements) == 0 {
		return nil, nil
	}

	var results []*driver.Result
	for _, single := range statements {
		statement := single.Text
		if opts.Explain {
			statement = explainStatement(statement, opts.ExplainFormat)
		} else if opts.Limit > 0 {
			statement = statementWithResultLimit(statement, opts.Limit)
		}
		startTime := time.Now()
		allQuery, err := validateSQLForEditor(statement)
		if err != nil {
			results = append(results, &driver.Result{Statement: statement, Err: err, Latency: time.Since(startTime)})
			break
		}

		var result *driver.Result
		var queryErr error
		if allQuery {
			result, queryErr = s.runQuery(ctx, statement, opts.MaxResultBytes)
		} else {
			result, queryErr = s.runExec(ctx, statement)
		}
		stop := false
		if queryErr != nil {
			queryErr = translateAggregateFunctionError(queryErr)
			result = &driver.Result{Err: queryErr}
			stop = true
		}
		result.Statement = statement
		result.Latency = time.Since(startTime)
		results = append(results, result)
		if stop {
			break
		}
	}
	return results, nil
}

func (s *session) runQuery(ctx context.Context, statement string, maxResultBytes int64) (*driver.Result, error) {
	rows, err := s.conn.QueryContext(ctx, statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result, err := readRows(rows, maxResultBytes)
	if err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *session) runExec(ctx context.Context, statement string) (*driver.Result, error) {
	sqlResult, err := s.conn.ExecContext(ctx, statement)
	if err != nil {
		return nil, err
	}
	affectedRows, err := sqlResult.RowsAffected()
	if err != nil {
		// A DDL statement reports no count, which the client reports as an
		// error; Bytebase ignored it and recorded zero.
		affectedRows = 0
	}
	return &driver.Result{AffectedRows: &affectedRows}, nil
}

// readRows reads one result set, stopping at maxResultBytes.
func readRows(rows *sql.Rows, maxResultBytes int64) (*driver.Result, error) {
	columnNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	columns := make([]driver.Column, len(columnNames))
	for i, columnType := range columnTypes {
		columns[i] = driver.Column{
			Name:     columnNames[i],
			TypeName: strings.ToUpper(columnType.DatabaseTypeName()),
		}
	}
	result := &driver.Result{Columns: columns}
	if len(columns) == 0 {
		return result, nil
	}
	totalBytes := driver.HeaderByteSize(columns)
	for rows.Next() {
		dest := make([]any, len(columns))
		for i, column := range columns {
			dest[i] = makeValueByTypeName(column.TypeName, columnTypes[i])
		}
		if err := rows.Scan(dest...); err != nil {
			return nil, err
		}
		row := make([]driver.Value, len(columns))
		for i, column := range columns {
			row[i] = convertValue(column.TypeName, columnTypes[i], dest[i])
		}
		result.Rows = append(result.Rows, row)
		if maxResultBytes > 0 {
			totalBytes += driver.RowByteSize(row)
			if totalBytes > maxResultBytes {
				result.Truncated = true
				break
			}
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

// makeValueByTypeName prepares a scan destination. ClickHouse's tuple, array,
// map and JSON values are scanned into an any, which the client fills with its
// own Go value; every other type scans into the client's own type.
func makeValueByTypeName(typeName string, columnType *sql.ColumnType) any {
	if strings.HasPrefix(typeName, "TUPLE") || strings.HasPrefix(typeName, "ARRAY") || strings.HasPrefix(typeName, "MAP") || strings.HasPrefix(typeName, "JSON") {
		var it any
		return &it
	}
	// https://github.com/ClickHouse/clickhouse-go/blob/main/TYPES.md
	return reflect.New(columnType.ScanType()).Interface()
}

// convertValue turns one scanned cell into a contract Value. The branches are
// Bytebase's, split between the non-nullable and nullable forms the client
// hands back.
func convertValue(_ string, _ *sql.ColumnType, value any) driver.Value {
	// TUPLE, ARRAY, MAP, JSON
	if v, ok := value.(*any); ok && v != nil {
		b, err := json.Marshal(v)
		if err != nil {
			return driver.Value{Kind: driver.Null}
		}
		return driver.Value{Kind: driver.String, String: string(b)}
	}

	// not nullable
	if v, ok := value.(*int); ok && v != nil {
		return driver.Value{Kind: driver.Int64, Int64: int64(*v)}
	}
	if v, ok := value.(*int8); ok && v != nil {
		return driver.Value{Kind: driver.Int32, Int64: int64(*v)}
	}
	if v, ok := value.(*int16); ok && v != nil {
		return driver.Value{Kind: driver.Int32, Int64: int64(*v)}
	}
	if v, ok := value.(*int32); ok && v != nil {
		return driver.Value{Kind: driver.Int32, Int64: int64(*v)}
	}
	if v, ok := value.(*int64); ok && v != nil {
		return driver.Value{Kind: driver.Int64, Int64: *v}
	}
	if v, ok := value.(*uint); ok && v != nil {
		return driver.Value{Kind: driver.Uint64, Uint64: uint64(*v)}
	}
	if v, ok := value.(*uint8); ok && v != nil {
		return driver.Value{Kind: driver.Uint32, Uint64: uint64(*v)}
	}
	if v, ok := value.(*uint16); ok && v != nil {
		return driver.Value{Kind: driver.Uint32, Uint64: uint64(*v)}
	}
	if v, ok := value.(*uint32); ok && v != nil {
		return driver.Value{Kind: driver.Uint32, Uint64: uint64(*v)}
	}
	if v, ok := value.(*uint64); ok && v != nil {
		return driver.Value{Kind: driver.Uint64, Uint64: *v}
	}
	if v, ok := value.(*float32); ok && v != nil {
		return driver.Value{Kind: driver.Float32, Float64: float64(*v)}
	}
	if v, ok := value.(*float64); ok && v != nil {
		return driver.Value{Kind: driver.Float64, Float64: *v}
	}
	if v, ok := value.(*string); ok && v != nil {
		return driver.Value{Kind: driver.String, String: *v}
	}
	if v, ok := value.(*bool); ok && v != nil {
		return driver.Value{Kind: driver.Bool, Bool: *v}
	}
	if v, ok := value.(*time.Time); ok && v != nil {
		return driver.Value{Kind: driver.String, String: v.Format(time.RFC3339Nano)}
	}
	if v, ok := value.(*net.IP); ok && v != nil {
		return driver.Value{Kind: driver.String, String: v.String()}
	}
	if v, ok := value.(*big.Int); ok && v != nil {
		return driver.Value{Kind: driver.String, String: v.String()}
	}
	if v, ok := value.(*decimal.Decimal); ok && v != nil {
		return driver.Value{Kind: driver.String, String: v.String()}
	}
	if v, ok := value.(*uuid.UUID); ok && v != nil {
		return driver.Value{Kind: driver.String, String: v.String()}
	}
	if v, ok := value.(*orb.Point); ok && v != nil {
		return driver.Value{Kind: driver.String, String: wkt.MarshalString(*v)}
	}
	if v, ok := value.(*orb.Polygon); ok && v != nil {
		return driver.Value{Kind: driver.String, String: wkt.MarshalString(*v)}
	}
	if v, ok := value.(*orb.Ring); ok && v != nil {
		return driver.Value{Kind: driver.String, String: wkt.MarshalString(*v)}
	}
	if v, ok := value.(*orb.MultiPolygon); ok && v != nil {
		return driver.Value{Kind: driver.String, String: wkt.MarshalString(*v)}
	}

	// nullable
	if v, ok := value.(**int); ok && *v != nil {
		return driver.Value{Kind: driver.Int64, Int64: int64(**v)}
	}
	if v, ok := value.(**int8); ok && *v != nil {
		return driver.Value{Kind: driver.Int32, Int64: int64(**v)}
	}
	if v, ok := value.(**int16); ok && *v != nil {
		return driver.Value{Kind: driver.Int32, Int64: int64(**v)}
	}
	if v, ok := value.(**int32); ok && *v != nil {
		return driver.Value{Kind: driver.Int32, Int64: int64(**v)}
	}
	if v, ok := value.(**int64); ok && *v != nil {
		return driver.Value{Kind: driver.Int64, Int64: **v}
	}
	if v, ok := value.(**uint); ok && *v != nil {
		return driver.Value{Kind: driver.Uint64, Uint64: uint64(**v)}
	}
	if v, ok := value.(**uint8); ok && *v != nil {
		return driver.Value{Kind: driver.Uint32, Uint64: uint64(**v)}
	}
	if v, ok := value.(**uint16); ok && *v != nil {
		return driver.Value{Kind: driver.Uint32, Uint64: uint64(**v)}
	}
	if v, ok := value.(**uint32); ok && *v != nil {
		return driver.Value{Kind: driver.Uint32, Uint64: uint64(**v)}
	}
	if v, ok := value.(**uint64); ok && *v != nil {
		return driver.Value{Kind: driver.Uint64, Uint64: **v}
	}
	if v, ok := value.(**float32); ok && *v != nil {
		return driver.Value{Kind: driver.Float32, Float64: float64(**v)}
	}
	if v, ok := value.(**float64); ok && *v != nil {
		return driver.Value{Kind: driver.Float64, Float64: **v}
	}
	if v, ok := value.(**string); ok && *v != nil {
		return driver.Value{Kind: driver.String, String: **v}
	}
	if v, ok := value.(**bool); ok && *v != nil {
		return driver.Value{Kind: driver.Bool, Bool: **v}
	}
	if v, ok := value.(**time.Time); ok && *v != nil {
		return driver.Value{Kind: driver.String, String: (*v).Format(time.RFC3339Nano)}
	}
	if v, ok := value.(**net.IP); ok && *v != nil {
		return driver.Value{Kind: driver.String, String: (*v).String()}
	}
	if v, ok := value.(**big.Int); ok && *v != nil {
		return driver.Value{Kind: driver.String, String: (*v).String()}
	}
	if v, ok := value.(**decimal.Decimal); ok && *v != nil {
		return driver.Value{Kind: driver.String, String: (*v).String()}
	}
	if v, ok := value.(**uuid.UUID); ok && *v != nil {
		return driver.Value{Kind: driver.String, String: (*v).String()}
	}
	return driver.Value{Kind: driver.Null}
}

// translateAggregateFunctionError rewrites the client's opaque "unsupported
// column type \"AggregateFunction(...)\"" error into actionable guidance
// pointing at -Merge combinators or finalizeAggregation(), since the client
// cannot decode AggregateFunction intermediate state directly.
func translateAggregateFunctionError(err error) error {
	if err == nil || !strings.Contains(err.Error(), aggregateFunctionErrMarker) {
		return err
	}
	return fmt.Errorf("this query selects an AggregateFunction column, which stores aggregator intermediate state and cannot be decoded directly. "+
		"Wrap each AggregateFunction column with its -Merge combinator under GROUP BY (for example, maxMerge(col)), or with "+
		"finalizeAggregation(col) for per-row finalization: %w", err)
}
