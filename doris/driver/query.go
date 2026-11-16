package driver

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/bytebase/omni/doris/analysis"
	"github.com/bytebase/omni/doris/parser"
	"github.com/bytebase/omni/driver"
)

// splitCommands splits script into the statements a query request runs, in
// order, each with the text and byte range the task log reports.
func splitCommands(script string) []driver.Command {
	segments := parser.Split(script)
	commands := make([]driver.Command, 0, len(segments))
	for _, seg := range segments {
		if seg.Empty() {
			continue
		}
		end := seg.ByteEnd
		if end < len(script) && script[end] == ';' {
			end++
		}
		commands = append(commands, driver.Command{
			Text:      script[seg.ByteStart:end],
			ByteStart: seg.ByteStart,
			ByteEnd:   end,
		})
	}
	return commands
}

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
// connection. Errors returned by Doris unwrap to *mysql.MySQLError.
func (s *session) Query(ctx context.Context, script string, opts driver.QueryOptions) ([]*driver.Result, error) {
	commands := splitCommands(script)
	if len(commands) == 0 {
		return nil, nil
	}

	// Rewrite every statement before running any, so a rewrite that fails
	// leaves no statement results behind.
	type preparedStatement struct {
		text        string
		returnsData bool
	}
	prepared := make([]preparedStatement, 0, len(commands))
	for _, cmd := range commands {
		statement := cmd.Text
		if opts.Explain {
			explained, err := analysis.ExplainStatement(statement)
			if err != nil {
				return nil, err
			}
			statement = explained
		} else if opts.Limit > 0 {
			statement = analysis.RewriteResultLimit(statement, opts.Limit)
		}

		_, returnsData, err := analysis.ValidateQuery(statement)
		if err != nil {
			// A statement omni cannot parse runs as a query, as Bytebase runs
			// it, and the server's answer decides.
			returnsData = true
		}
		prepared = append(prepared, preparedStatement{text: statement, returnsData: returnsData})
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
		if res.Err == nil {
			if format, executed, ok := analysis.DescribeExplain(stmt.text); ok {
				res.Plan = &driver.Plan{Format: format, Executed: executed}
			}
		}
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
			if typeName == "TIMESTAMP" || typeName == "DATETIME" {
				_, scale, _ := columnType.DecimalSize()
				return parseTimestampOrStringValue(raw.String, int32(scale))
			}
			return driver.Value{Kind: driver.String, String: raw.String}
		}
	case *sql.NullInt64:
		if raw.Valid {
			return driver.Value{Kind: driver.Int64, Int64: raw.Int64}
		}
	case *[]byte:
		// nil is SQL NULL; an empty non-nil slice is a zero-byte value.
		if *raw != nil {
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
	var (
		t   time.Time
		err error
	)
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
