package driver

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/bytebase/omni/driver"
)

// Session opens one server connection held until Close, so a database USE
// selects stays selected for the next Query.
func (d *Driver) Session(ctx context.Context) (driver.Session, error) {
	conn, err := d.db.Conn(ctx)
	if err != nil {
		return nil, err
	}
	if err := d.selectDatabase(ctx, conn); err != nil {
		_ = conn.Close()
		return nil, err
	}
	return &session{conn: conn}, nil
}

// selectDatabase starts conn from the database this driver was opened on, or
// from Hive's default when none was named. The pool has no reset hook for Hive,
// so a connection another session left on its own database arrives still there.
func (d *Driver) selectDatabase(ctx context.Context, conn *sql.Conn) error {
	name := d.databaseName
	if name == "" {
		name = "default"
	}
	if _, err := conn.ExecContext(ctx, "USE "+quoteBacktick(name)); err != nil {
		return fmt.Errorf("failed to select database %s: %w", name, err)
	}
	return nil
}

type session struct {
	conn *sql.Conn
}

func (s *session) Close() error { return s.conn.Close() }

// Query runs each statement of script in order on the session's connection.
//
// Every statement runs as a query, as Bytebase ran it: Hive answers a
// statement that returns rows with a result set and any other with an
// operation that has none, and reading that operation's rows is where Hive's
// error arrives, as the statement's Result.Err. Change scripts go through
// Execute instead.
func (s *session) Query(ctx context.Context, script string, opts driver.QueryOptions) ([]*driver.Result, error) {
	commands, err := splitCommands(script)
	if err != nil {
		return nil, fmt.Errorf("failed to split statements: %w", err)
	}
	if len(commands) == 0 {
		return nil, nil
	}
	// Select the database the caller names every time: a statement of the
	// script may select another one, and its USE is a statement like any other.
	if opts.Schema != "" {
		if _, err := s.conn.ExecContext(ctx, "USE "+quoteBacktick(opts.Schema)); err != nil {
			return nil, fmt.Errorf("failed to select database %s: %w", opts.Schema, err)
		}
	}

	results := make([]*driver.Result, 0, len(commands))
	for _, command := range commands {
		statement := trimStatement(command.Text)
		if opts.Explain {
			statement = explainStatement(statement, opts.ExplainFormat)
		}
		start := time.Now()
		result, err := queryRows(ctx, s.conn, statement, opts.MaxResultBytes)
		if err != nil {
			// A statement that fails ends the script.
			results = append(results, &driver.Result{Statement: statement, Err: err, Latency: time.Since(start)})
			break
		}
		if format, ok := planFormat(statement); ok {
			result.Plan = &driver.Plan{Format: format}
		}
		results = append(results, result)
	}
	return results, nil
}

// queryStatement runs one statement on the driver's pool, for the sync.
func (d *Driver) queryStatement(ctx context.Context, statement string, maxResultBytes int64) (*driver.Result, error) {
	return queryRows(ctx, d.db, statement, maxResultBytes)
}

// queryer runs a statement: the driver's pool, or one session's connection.
type queryer interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

// queryRows runs one statement and reads its rows, stopping once the rows pass
// maxResultBytes.
func queryRows(ctx context.Context, q queryer, statement string, maxResultBytes int64) (*driver.Result, error) {
	startTime := time.Now()

	rows, err := q.QueryContext(ctx, statement)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, fmt.Errorf("failed to get column types: %w", err)
	}
	result := &driver.Result{Statement: statement}
	result.Columns = make([]driver.Column, len(columnTypes))
	for i, col := range columnTypes {
		result.Columns[i] = driver.Column{Name: col.Name(), TypeName: strings.ToUpper(col.DatabaseTypeName())}
	}

	values := make([]any, len(result.Columns))
	valuePtrs := make([]any, len(result.Columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	totalBytes := driver.HeaderByteSize(result.Columns)
	for rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		row := make([]driver.Value, len(result.Columns))
		for i, col := range result.Columns {
			row[i] = parseValueType(values[i], col.TypeName)
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
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	result.Latency = time.Since(startTime)
	return result, nil
}

// parseValueType maps one cell gohive decoded onto the contract's kinds.
// A float column becomes text with twenty fractional digits, as Bytebase
// mapped it; the driver's own DOUBLE and FLOAT values are the follow-up fix.
func parseValueType(value any, typeName string) driver.Value {
	if value == nil {
		return driver.Value{Kind: driver.Null}
	}

	switch typeName {
	case "BOOLEAN_TYPE", "BOOLEAN":
		if b, ok := value.(bool); ok {
			return driver.Value{Kind: driver.Bool, Bool: b}
		}
	case "TINYINT_TYPE", "TINYINT", "SMALLINT_TYPE", "SMALLINT", "INT_TYPE", "INT":
		if i, ok := value.(int64); ok {
			return driver.Value{Kind: driver.Int32, Int64: int64(int32(i))}
		}
	case "BIGINT_TYPE", "BIGINT":
		if i, ok := value.(int64); ok {
			return driver.Value{Kind: driver.Int64, Int64: i}
		}
	case "DOUBLE_TYPE", "DOUBLE":
		if f, ok := value.(float64); ok {
			return driver.Value{Kind: driver.Float64, Float64: f}
		}
	case "FLOAT_TYPE", "FLOAT":
		if f, ok := value.(float64); ok {
			// FLOAT is four bytes on the server; gohive decodes every number as a
			// float64, so narrow it to the value Hive stored.
			return driver.Value{Kind: driver.Float32, Float64: float64(float32(f))}
		}
	case "BINARY_TYPE", "BINARY":
		if b, ok := value.([]byte); ok {
			return driver.Value{Kind: driver.Bytes, Bytes: b}
		}
	}

	// Every other type arrives as the Go value its cell holds, and a string
	// is what a caller reads: gohive decodes a string column as a string.
	switch v := value.(type) {
	case []byte:
		return driver.Value{Kind: driver.String, String: string(v)}
	case string:
		return driver.Value{Kind: driver.String, String: v}
	case int64:
		return driver.Value{Kind: driver.String, String: strconv.FormatInt(v, 10)}
	case float64:
		return driver.Value{Kind: driver.String, String: strconv.FormatFloat(v, 'f', -1, 64)}
	case bool:
		return driver.Value{Kind: driver.String, String: strconv.FormatBool(v)}
	default:
		return driver.Value{Kind: driver.String, String: fmt.Sprintf("%v", v)}
	}
}

// trimStatement drops the whitespace and semicolons around a statement, so a
// split segment's terminator does not reach the server.
func trimStatement(statement string) string {
	return strings.TrimLeftFunc(strings.TrimRightFunc(statement, isSpaceOrSemicolon), unicode.IsSpace)
}

func isSpaceOrSemicolon(r rune) bool {
	return unicode.IsSpace(r) || r == ';'
}

func quoteBacktick(name string) string { return "`" + strings.ReplaceAll(name, "`", "``") + "`" }
