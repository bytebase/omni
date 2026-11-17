package driver

import (
	"context"
	"database/sql"
	sqldriver "database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/bytebase/omni/driver"
	mysqlanalysis "github.com/bytebase/omni/mysql/analysis"
	"github.com/bytebase/omni/mysql/ast"
	mysqlparser "github.com/bytebase/omni/mysql/parser"
)

// Session opens one server connection held until Close so session variables
// and temporary tables persist across Query calls. The connection's default
// database is the driver's, which is why QueryOptions.Schema is ignored: the
// caller selects a schema for its own analysis, and a statement that means
// another database names it.
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

// Close discards the physical connection: see discardConn.
func (s *session) Close() error {
	discardConn(s.conn)
	return nil
}

// discardConn closes conn's physical connection instead of returning it to the
// pool. MySQL's session state outlives a statement — the default database,
// session variables, temporary tables, an open transaction — and
// go-sql-driver's ResetSession only checks that the connection is alive, so a
// pooled connection would hand that state to the next caller. Returning
// ErrBadConn from Raw makes database/sql close it, and marks the Conn used up,
// so no further call is needed.
func discardConn(conn *sql.Conn) {
	_ = conn.Raw(func(any) error { return sqldriver.ErrBadConn })
}

// Query runs each non-empty statement of script in order on the session's
// connection. Errors returned by MySQL unwrap to *mysql.MySQLError. A script
// the parser cannot read is the call's error and no statement of it runs:
// every statement is parsed before the first is sent.
//
// Masking is the caller's work: every row here comes from a statement whose
// columns its analysis traces, so QueryOptions.MaskingEnabled changes nothing.
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
		parsed, err := parseStatement(script, cmd)
		if err != nil {
			return nil, err
		}
		// An executable comment holding no statement, as in /*! */, parses to
		// nothing: there is no statement to run. A command that parses to
		// several statements, as one statement terminator inside a DELIMITER
		// block does, is refused: the connection would send them together and
		// only the first one's result would come back.
		if len(parsed.Items) == 0 {
			continue
		}
		if len(parsed.Items) > 1 {
			return nil, fmt.Errorf("a command must hold one statement, but %q holds %d", strings.TrimSpace(cmd.Text), len(parsed.Items))
		}
		node := parsed.Items[0]
		_, returnsData := mysqlanalysis.ValidateQuery(parsed.Items)

		statement := cmd.Text
		if opts.Explain {
			explained, err := mysqlanalysis.ExplainStatement(statement, strings.EqualFold(opts.ExplainFormat, "JSON"))
			if err != nil {
				return nil, scriptParseError(script, cmd.ByteStart, err)
			}
			statement = explained
			// An explain request answers with a plan, or with the description
			// of a table that EXPLAIN names.
			returnsData = true
		} else if opts.Limit > 0 {
			statement = mysqlanalysis.StatementWithResultLimit(statement, opts.Limit)
		}

		// A statement the read-only whitelist does not name may still answer
		// with rows, and masking withholds the rows of a statement whose
		// columns the caller's analysis cannot trace.
		if serverReturnsRows(node) {
			returnsData = !opts.MaskingEnabled
		}
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

// serverReturnsRows reports whether the server answers node with a result set
// that mysql/analysis's read-only whitelist does not name: a stored procedure's
// result sets, a table check's diagnosis, and a handler's read. It is what the
// pinned servers answer with columns; DO, SET, FLUSH, a cache load without a
// cache, and every write answer none.
func serverReturnsRows(node ast.Node) bool {
	switch node.(type) {
	case *ast.CallStmt, *ast.CheckTableStmt, *ast.ChecksumTableStmt,
		*ast.AnalyzeTableStmt, *ast.OptimizeTableStmt, *ast.RepairTableStmt,
		*ast.LoadIndexIntoCacheStmt, *ast.HandlerReadStmt:
		return true
	}
	return false
}

// parseStatement parses the statement at cmd, reporting a parse error at the
// script's coordinates so a caller locates it in the text it supplied.
func parseStatement(script string, cmd driver.Command) (*ast.List, error) {
	parsed, err := mysqlparser.Parse(cmd.Text)
	if err != nil {
		return nil, scriptParseError(script, cmd.ByteStart, err)
	}
	return parsed, nil
}

// scriptParseError moves a statement's parse error onto the script the
// statement came from. A segment parsed on its own reports positions relative
// to itself, and a caller holds the script.
func scriptParseError(script string, offset int, err error) error {
	var parseErr *mysqlparser.ParseError
	if offset == 0 || !errors.As(err, &parseErr) {
		return err
	}
	moved := *parseErr
	moved.Position += offset
	if moved.Position < 0 || moved.Position > len(script) {
		return &moved
	}
	moved.Line = 1 + strings.Count(script[:moved.Position], "\n")
	lineStart := strings.LastIndexByte(script[:moved.Position], '\n') + 1
	lineEnd := strings.IndexByte(script[moved.Position:], '\n')
	if lineEnd < 0 {
		lineEnd = len(script)
	} else {
		lineEnd += moved.Position
	}
	moved.Column = utf8.RuneCountInString(script[lineStart:moved.Position]) + 1
	moved.RelatedText = script[lineStart:lineEnd]
	return &moved
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
		// Only the first result set is the statement's result. The rest are
		// drained rather than closed, because a later set can carry the
		// statement's error, as a procedure that selects and then signals does.
		for rows.NextResultSet() {
			for rows.Next() {
			}
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
		if format, executed, ok := mysqlanalysis.DescribeExplain(statement); ok {
			res.Plan = &driver.Plan{
				Format:   strings.ToUpper(format),
				Executed: executed,
			}
		}
		return res, nil
	}

	// A statement that returns no rows is sent without the application
	// comment, which is what Bytebase's MySQL driver does.
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
	case "FLOAT", "DOUBLE", "FLOAT4", "FLOAT8":
		return new(sql.NullFloat64)
	case "INT", "INTEGER", "TINYINT", "SMALLINT", "MEDIUMINT", "BIGINT", "INT2", "INT4", "INT8":
		return new(sql.NullInt64)
	case "UNSIGNED TINYINT", "UNSIGNED SMALLINT", "UNSIGNED MEDIUMINT", "UNSIGNED INT", "UNSIGNED BIGINT":
		return new(sql.Null[uint64])
	case "BIT", "VARBIT", "BINARY", "VARBINARY", "BLOB", "TINYBLOB", "MEDIUMBLOB", "LONGBLOB",
		"GEOMETRY", "POINT", "LINESTRING", "POLYGON", "MULTIPOINT", "MULTILINESTRING", "MULTIPOLYGON", "GEOMETRYCOLLECTION":
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
	case *sql.Null[uint64]:
		if raw.Valid {
			return driver.Value{Kind: driver.Uint64, Uint64: raw.V}
		}
	case *[]byte:
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
	var t time.Time
	var err error
	switch {
	case strings.Contains(raw, "T"):
		// A caller that sets the client's parseTime parameter is handed a
		// time.Time, which database/sql renders as RFC 3339 into the string
		// the column scans into.
		t, err = time.Parse(time.RFC3339Nano, raw)
	case strings.Contains(raw, "."):
		t, err = time.Parse("2006-01-02 15:04:05.999999999", raw)
	default:
		t, err = time.Parse(time.DateTime, raw)
	}
	if err != nil || t.Year() < 1 || t.Year() > 9999 {
		return driver.Value{Kind: driver.String, String: raw}
	}
	return driver.Value{Kind: driver.Timestamp, Time: t, Precision: scale}
}
