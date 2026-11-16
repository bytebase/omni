package driver

import (
	"context"
	"database/sql"
	"encoding"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"net"
	"reflect"
	"strconv"
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
		if opts.Explain && queryErr == nil {
			result.Plan = &driver.Plan{Format: explainPlanFormat(opts.ExplainFormat)}
		}
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
func convertValue(typeName string, _ *sql.ColumnType, value any) driver.Value {
	// TUPLE, ARRAY, MAP, JSON
	if v, ok := value.(*any); ok && v != nil {
		if *v == nil {
			return driver.Value{Kind: driver.Null}
		}
		b, err := json.Marshal(normalizeJSONValue(*v))
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
		return timestampValue(typeName, *v)
	}
	if v, ok := value.(*time.Duration); ok && v != nil {
		return driver.Value{Kind: driver.String, String: formatDuration(*v)}
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
		return timestampValue(typeName, **v)
	}
	if v, ok := value.(**time.Duration); ok && *v != nil {
		return driver.Value{Kind: driver.String, String: formatDuration(**v)}
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

// timestampValue is a Date, DateTime or DateTime64 in the contract's temporal
// kinds: a wall clock for a type without a zone, an instant in its zone for one
// with a zone, and the fractional digits the type names.
func timestampValue(typeName string, t time.Time) driver.Value {
	precision := temporalPrecision(typeName)
	if isZonedTemporal(typeName) {
		return driver.Value{Kind: driver.TimestampTZ, Time: t, Precision: precision}
	}
	return driver.Value{
		Kind:      driver.Timestamp,
		Time:      time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.UTC),
		Precision: precision,
	}
}

func baseTemporalType(typeName string) string {
	if strings.HasPrefix(strings.ToUpper(typeName), "NULLABLE(") && strings.HasSuffix(typeName, ")") {
		return typeName[len("Nullable(") : len(typeName)-1]
	}
	return typeName
}

func isZonedTemporal(typeName string) bool {
	base := strings.ToUpper(baseTemporalType(typeName))
	return strings.HasPrefix(base, "DATETIME") && strings.Contains(base, "'")
}

func temporalPrecision(typeName string) int32 {
	base := baseTemporalType(typeName)
	if !strings.HasPrefix(strings.ToUpper(base), "DATETIME64(") {
		return 0
	}
	rest := base[len("DateTime64("):]
	end := strings.IndexAny(rest, ",)")
	if end < 0 {
		return 0
	}
	digits, err := strconv.Atoi(strings.TrimSpace(rest[:end]))
	if err != nil {
		return 0
	}
	return int32(digits)
}

// formatDuration renders a ClickHouse Time or Time64 as a clock, which is what
// the value is: a duration since midnight.
func formatDuration(d time.Duration) string {
	sign := ""
	if d < 0 {
		sign = "-"
		d = -d
	}
	hours := d / time.Hour
	minutes := (d % time.Hour) / time.Minute
	seconds := (d % time.Minute) / time.Second
	out := fmt.Sprintf("%s%02d:%02d:%02d", sign, hours, minutes, seconds)
	if nanos := d % time.Second; nanos != 0 {
		out += "." + strings.TrimRight(fmt.Sprintf("%09d", nanos), "0")
	}
	return out
}

// normalizeJSONValue turns typed maps into string-keyed maps, which is what a
// document can hold: a ClickHouse Map(Float64, ...) arrives as
// map[float64]..., and encoding/json refuses such a key.
func normalizeJSONValue(value any) any {
	switch value.(type) {
	case json.Marshaler, encoding.TextMarshaler:
		return value
	}
	rv := reflect.ValueOf(value)
	if !rv.IsValid() {
		return nil
	}
	switch rv.Kind() {
	case reflect.Map:
		out := make(map[string]any, rv.Len())
		iter := rv.MapRange()
		for iter.Next() {
			out[fmt.Sprint(iter.Key().Interface())] = normalizeJSONValue(iter.Value().Interface())
		}
		return out
	case reflect.Float32, reflect.Float64:
		f := rv.Float()
		switch {
		case math.IsNaN(f):
			return "nan"
		case math.IsInf(f, 1):
			return "inf"
		case math.IsInf(f, -1):
			return "-inf"
		default:
			return value
		}
	case reflect.Slice, reflect.Array:
		out := make([]any, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			out[i] = normalizeJSONValue(rv.Index(i).Interface())
		}
		return out
	default:
		return value
	}
}

func explainPlanFormat(format string) string {
	if strings.EqualFold(format, "JSON") {
		return "JSON"
	}
	return "TEXT"
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
