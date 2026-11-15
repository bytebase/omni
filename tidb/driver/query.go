package driver

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/bytebase/omni/driver"
	tidbanalysis "github.com/bytebase/omni/tidb/analysis"
	tidbparser "github.com/bytebase/omni/tidb/parser"
)

// Session opens one server connection held until Close so session variables
// and temporary tables persist across Query calls.
func (d *Driver) Session(ctx context.Context) (driver.Session, error) {
	conn, err := d.db.Conn(ctx)
	if err != nil {
		return nil, err
	}
	connID, err := d.readConnectionID(ctx, conn)
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	return &session{driver: d, conn: conn, connID: connID}, nil
}

type session struct {
	driver *Driver
	conn   *sql.Conn
	connID uint64
}

func (s *session) Close() error {
	return s.conn.Close()
}

// Query runs each non-empty statement of script in order on the session's
// connection. Errors returned by TiDB unwrap to *mysql.MySQLError.
func (s *session) Query(ctx context.Context, script string, opts driver.QueryOptions) ([]*driver.Result, error) {
	rawStmts := splitCommands(script)
	if len(rawStmts) == 0 {
		return nil, nil
	}

	type preparedStmt struct {
		text        string
		returnsData bool
	}
	prepared := make([]preparedStmt, 0, len(rawStmts))
	for _, cmd := range rawStmts {
		statement := cmd.Text
		if opts.Explain {
			explained, err := tidbanalysis.ExplainStatement(statement, strings.EqualFold(opts.ExplainFormat, "JSON"))
			if err != nil {
				return nil, err
			}
			statement = explained
		} else if opts.Limit > 0 {
			if rewritten, err := tidbanalysis.RewriteResultLimit(statement, opts.Limit); err == nil {
				statement = rewritten
			} else {
				statement = fmt.Sprintf("WITH result AS (%s) SELECT * FROM result LIMIT %d;", trimStatement(statement), opts.Limit)
			}
		}

		parsed, err := tidbparser.Parse(statement)
		if err != nil {
			return nil, err
		}
		_, returnsData := tidbanalysis.ValidateQuery(parsed.Items)
		prepared = append(prepared, preparedStmt{text: statement, returnsData: returnsData})
	}

	var results []*driver.Result
	for _, stmt := range prepared {
		startTime := time.Now()
		res, err := s.runStatement(ctx, stmt.text, stmt.returnsData, opts.MaxResultBytes)
		stop := false
		if err != nil {
			s.driver.killQueryOnCancel(err, s.connID)
			res = &driver.Result{Err: err}
			stop = true
		}
		res.Statement = stmt.text
		res.Latency = time.Since(startTime)
		results = append(results, res)
		if stop {
			break
		}
	}
	return results, nil
}

func (s *session) runStatement(ctx context.Context, statement string, returnsData bool, maxResultBytes int64) (*driver.Result, error) {
	if returnsData {
		rows, err := s.conn.QueryContext(ctx, s.driver.withAppComment(statement))
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		res, err := readRows(rows, maxResultBytes)
		if err != nil {
			return nil, err
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
		if format, executed, ok := tidbanalysis.DescribeExplain(statement); ok {
			res.Plan = &driver.Plan{
				Format:   strings.ToUpper(format),
				Executed: executed,
			}
		}
		return res, nil
	}

	sqlResult, err := s.conn.ExecContext(ctx, statement)
	if err != nil {
		return nil, err
	}
	affectedRows, err := sqlResult.RowsAffected()
	if err != nil {
		affectedRows = 0
	}
	return &driver.Result{AffectedRows: &affectedRows}, nil
}

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
	for i, ct := range columnTypes {
		columns[i] = driver.Column{
			Name:     columnNames[i],
			TypeName: strings.ToUpper(ct.DatabaseTypeName()),
		}
	}
	res := &driver.Result{Columns: columns}
	totalBytes := driver.HeaderByteSize(columns)

	if len(columns) > 0 {
		for rows.Next() {
			dest := make([]any, len(columns))
			for i, col := range columns {
				dest[i] = makeValueByTypeName(col.TypeName)
			}
			if err := rows.Scan(dest...); err != nil {
				return nil, err
			}
			row := make([]driver.Value, len(columns))
			for i, col := range columns {
				row[i] = convertValue(col.TypeName, columnTypes[i], dest[i])
			}
			res.Rows = append(res.Rows, row)
			if maxResultBytes > 0 {
				totalBytes += driver.RowByteSize(row)
				if totalBytes > maxResultBytes {
					res.Truncated = true
					break
				}
			}
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func makeValueByTypeName(typeName string) any {
	switch typeName {
	case "VARCHAR", "TEXT", "UUID", "DATETIME", "TIMESTAMP":
		return new(sql.NullString)
	case "BOOL":
		return new(sql.NullBool)
	case "INT", "INTEGER", "TINYINT", "SMALLINT", "MEDIUMINT", "BIGINT", "INT2", "INT4", "INT8":
		return new(sql.NullInt64)
	case "FLOAT", "DOUBLE", "FLOAT4", "FLOAT8":
		return new(sql.NullFloat64)
	case "BIT", "VARBIT":
		return new([]byte)
	default:
		return new(sql.NullString)
	}
}

func convertValue(typeName string, columnType *sql.ColumnType, value any) driver.Value {
	switch raw := value.(type) {
	case *sql.NullString:
		if raw.Valid {
			_, scale, _ := columnType.DecimalSize()
			if typeName == "TIMESTAMP" || typeName == "DATETIME" {
				return parseTimestampOrStringValue(raw.String, int32(scale))
			}
			return driver.Value{Kind: driver.String, String: raw.String}
		}
	case *sql.NullInt64:
		if raw.Valid {
			return driver.Value{Kind: driver.Int64, Int64: raw.Int64}
		}
	case *[]byte:
		if len(*raw) > 0 {
			return driver.Value{Kind: driver.Bytes, Bytes: *raw}
		}
	case *sql.NullBool:
		if raw.Valid {
			return driver.Value{Kind: driver.Bool, Bool: raw.Bool}
		}
	case *sql.NullFloat64:
		if raw.Valid {
			return driver.Value{Kind: driver.Float64, Float64: raw.Float64}
		}
	}
	return driver.Value{Kind: driver.Null}
}

func parseTimestampOrStringValue(raw string, scale int32) driver.Value {
	var t time.Time
	var err error
	if strings.Contains(raw, ".") {
		t, err = time.Parse("2006-01-02 15:04:05.999999999", raw)
	} else {
		t, err = time.Parse(time.DateTime, raw)
	}
	if err != nil || t.Year() < 1 || t.Year() > 9999 {
		return driver.Value{Kind: driver.String, String: raw}
	}
	return driver.Value{Kind: driver.Timestamp, Time: t, Precision: scale}
}

func trimStatement(statement string) string {
	return strings.TrimLeftFunc(strings.TrimRightFunc(statement, func(r rune) bool {
		return unicode.IsSpace(r) || r == ';'
	}), unicode.IsSpace)
}
