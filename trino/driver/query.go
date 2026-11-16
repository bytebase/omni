package driver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	trinodb "github.com/trinodb/trino-go-client/trino"

	"github.com/bytebase/omni/driver"
	"github.com/bytebase/omni/trino/ast"
	"github.com/bytebase/omni/trino/parser"
)

// Session opens one server connection, held until Close, so what a statement
// sets on it (the schema, a session property) holds for the next.
func (d *Driver) Session(ctx context.Context) (driver.Session, error) {
	conn, err := d.db.Conn(ctx)
	if err != nil {
		return nil, err
	}
	return &session{conn: conn}, nil
}

type session struct {
	conn *sql.Conn
}

func (s *session) Close() error {
	return s.conn.Close()
}

// Query runs each statement of script in order on the session's connection.
// Errors returned by Trino unwrap to *trinodb.ErrQueryFailed, whose Reason is
// the server's *trinodb.ErrTrino.
func (s *session) Query(ctx context.Context, script string, opts driver.QueryOptions) ([]*driver.Result, error) {
	segments := parser.Split(script)

	if opts.Schema != "" {
		escaped := strings.ReplaceAll(opts.Schema, `"`, `""`)
		// The schema names the connection's own setting, not anything in the
		// statement, so it is quoted rather than parameterized.
		if _, err := s.conn.ExecContext(ctx, fmt.Sprintf("USE \"%s\"", escaped)); err != nil {
			return nil, fmt.Errorf("failed to set schema: %w", err)
		}
	}

	var results []*driver.Result
	for _, segment := range segments {
		if segment.Empty() {
			continue
		}
		statement := trimStatement(segment.Text)
		startTime := time.Now()
		result, err := s.runStatement(ctx, statement, opts)
		stop := false
		if err != nil {
			if result == nil {
				result = &driver.Result{Statement: statement}
			}
			result.Err = err
			stop = true
		}
		result.Latency = time.Since(startTime)
		results = append(results, result)
		if stop {
			break
		}
	}
	return results, nil
}

func (s *session) runStatement(ctx context.Context, statement string, opts driver.QueryOptions) (*driver.Result, error) {
	// Row-returning statements are told from the rest by their leading keyword,
	// not by parsing, so a statement this rule misreads runs on the exec path,
	// which discards the rows it returns.
	returnsData := returnsRows(statement)
	if opts.Explain {
		explained, err := explainStatement(statement, opts.ExplainFormat)
		if err != nil {
			return &driver.Result{Statement: statement}, err
		}
		statement = explained
		returnsData = true
	} else if returnsData && opts.Limit > 0 && !strings.Contains(strings.ToUpper(stripComments(statement)), " LIMIT ") {
		// A plan of a statement the caller did not write answers the wrong
		// question, so the limit applies only when the statement runs.
		statement = fmt.Sprintf("%s LIMIT %d", statement, opts.Limit)
	}

	if returnsData {
		rows, err := s.conn.QueryContext(ctx, statement)
		if err != nil {
			return &driver.Result{Statement: statement}, err
		}
		defer rows.Close()

		result, err := readRows(rows, opts.MaxResultBytes)
		if err != nil {
			return &driver.Result{Statement: statement}, err
		}
		if err := rows.Err(); err != nil {
			return &driver.Result{Statement: statement}, err
		}
		result.Statement = statement
		return result, nil
	}

	execResult, err := s.conn.ExecContext(ctx, statement)
	if err != nil {
		return &driver.Result{Statement: statement}, err
	}
	affectedRows, err := execResult.RowsAffected()
	if err != nil {
		// Trino does not always report a count, especially for DDL.
		affectedRows = 0
	}
	return &driver.Result{Statement: statement, AffectedRows: &affectedRows}, nil
}

// readRows reads a result set into a Result, stopping at maxResultBytes when
// that is positive.
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
	totalBytes := driver.HeaderByteSize(columns)

	if len(columns) > 0 {
		for rows.Next() {
			dest := make([]any, len(columns))
			for i, column := range columns {
				dest[i] = makeValueByTypeName(column.TypeName)
			}
			if err := rows.Scan(dest...); err != nil {
				return nil, err
			}
			row := make([]driver.Value, len(columns))
			for i, column := range columns {
				row[i] = convertValue(column.TypeName, dest[i])
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
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

// explainStatement returns the statement that plans statement. Trino spells
// its structured plan (FORMAT JSON); any other format gets the engine's own
// text plan.
func explainStatement(statement string, format string) (string, error) {
	node, err := singleStatement(statement)
	if err != nil {
		return "", err
	}
	planned := statement
	if explain, ok := node.(*parser.ExplainStmt); ok {
		if explain.Analyze {
			return "", errExplainAnalyzeExecutes
		}
		if inner, ok := explainedText(explain, statement); ok {
			planned = inner
		} else {
			// Nothing to rebuild around: run what the user wrote rather than
			// plan a plan.
			return statement, nil
		}
	}
	if strings.EqualFold(format, "JSON") {
		return "EXPLAIN (FORMAT JSON) " + planned, nil
	}
	return "EXPLAIN " + planned, nil
}

// errExplainAnalyzeExecutes is what an EXPLAIN ANALYZE request gets: planning it
// without the ANALYZE would answer a question the caller did not ask.
var errExplainAnalyzeExecutes = errors.New("EXPLAIN ANALYZE executes the statement it plans, which an explain request never does: run it as a query, or drop ANALYZE to plan the statement")

// singleStatement returns the one statement in text.
func singleStatement(text string) (ast.Node, error) {
	file, err := parser.Parse(text)
	if err != nil {
		return nil, err
	}
	if file == nil || len(file.Stmts) == 0 || file.Stmts[0] == nil {
		return nil, fmt.Errorf("EXPLAIN needs a statement to plan")
	}
	if len(file.Stmts) > 1 {
		return nil, fmt.Errorf("EXPLAIN plans a single statement")
	}
	return file.Stmts[0], nil
}

// explainedText returns the text within text of the statement explain plans,
// and whether it could be located.
func explainedText(explain *parser.ExplainStmt, text string) (string, bool) {
	inner, ok := nodeSpan(explain.Statement)
	if !ok {
		return "", false
	}
	// The plan ends where the EXPLAIN does, which is past clauses that the
	// planned statement's own location can leave out.
	if inner.Start < 0 || inner.Start >= explain.Loc.End || explain.Loc.End > len(text) {
		return "", false
	}
	return text[inner.Start:explain.Loc.End], true
}

func nodeSpan(node ast.Node) (ast.Loc, bool) {
	if query, ok := node.(*parser.QueryStmt); ok {
		return query.Loc, true
	}
	if spanned, ok := node.(interface{ Span() ast.Loc }); ok {
		return spanned.Span(), true
	}
	return ast.Loc{}, false
}

// returnsRows reports whether statement returns rows: Bytebase's rule, the
// leading keyword of the statement after any comment.
func returnsRows(statement string) bool {
	upper := strings.ToUpper(strings.TrimSpace(stripComments(statement)))
	return strings.HasPrefix(upper, "SELECT") ||
		strings.HasPrefix(upper, "SHOW") ||
		strings.HasPrefix(upper, "DESCRIBE") ||
		strings.HasPrefix(upper, "EXPLAIN")
}

// stripComments removes -- and /* */ comments from statement, respecting
// quoted text, so the leading-keyword and LIMIT checks see what the server
// would parse. It is the rule Bytebase's splitter applied to every statement
// before the driver saw it.
func stripComments(statement string) string {
	var out strings.Builder
	out.Grow(len(statement))
	const (
		singleQuote = '\''
		doubleQuote = '"'
		backtick    = '`'
	)
	var quote rune
	escaped := false
	for i := 0; i < len(statement); i++ {
		c := rune(statement[i])
		if quote != 0 {
			out.WriteByte(statement[i])
			switch {
			case escaped:
				escaped = false
			case c == '\\':
				escaped = true
			case c == quote:
				quote = 0
			}
			continue
		}
		switch {
		case c == '-' && i+1 < len(statement) && statement[i+1] == '-':
			for i < len(statement) && statement[i] != '\n' {
				i++
			}
			if i < len(statement) {
				out.WriteByte('\n')
			}
		case c == '/' && i+1 < len(statement) && statement[i+1] == '*':
			i += 2
			for i+1 < len(statement) && !(statement[i] == '*' && statement[i+1] == '/') {
				i++
			}
			i++
		case c == singleQuote || c == doubleQuote || c == backtick:
			quote = c
			out.WriteByte(statement[i])
		default:
			out.WriteByte(statement[i])
		}
	}
	return out.String()
}

// trimStatement trims a segment the way Bytebase's splitter trimmed the
// statements it returned, so a caller reads the text without its padding.
func trimStatement(statement string) string {
	return strings.Trim(statement, " \n\t")
}

// trinoTimeLayout renders a TIME the way Trino and PostgreSQL print one, which
// is what a SQL client shows: the fraction when the column's type carries one,
// and nothing when it does not. The nines trim trailing zeros and take the
// decimal point with them, so TIME keeps 12:34:56 while TIME(3) keeps
// 12:34:56.789 instead of losing it.
//
// A zoned time keeps its offset as well, since dropping it moves the instant.
const (
	trinoTimeLayout   = "15:04:05.999999999"
	trinoTimeTZLayout = trinoTimeLayout + "-07:00"
)

func makeValueByTypeName(typeName string) any {
	switch typeName {
	case "VARCHAR", "TEXT", "CHAR", "TIMESTAMP", "TIMESTAMP WITH TIME ZONE", "TIMESTAMP WITHOUT TIME ZONE":
		return new(sql.NullString)
	case "BOOLEAN":
		return new(sql.NullBool)
	case "TINYINT", "SMALLINT", "INTEGER", "BIGINT":
		return new(sql.NullInt64)
	case "DECIMAL", "REAL", "DOUBLE":
		return new(sql.NullFloat64)
	case "VARBINARY", "BINARY":
		return new([]byte)
	// The driver hands a VARIANT column a trino.Variant, which the default
	// *sql.NullString below cannot hold, failing the whole query.
	case "VARIANT":
		return new(trinodb.Variant)
	// The driver decodes these to time.Time. Left to the default below, Go
	// renders them through *sql.NullString as RFC3339Nano, so a DATE came back
	// as 2026-01-02T00:00:00Z and a TIME as 0000-01-01T12:34:56Z.
	case "DATE", "TIME", "TIME WITH TIME ZONE":
		return new(sql.NullTime)
	default:
		return new(sql.NullString)
	}
}

func convertValue(typeName string, value any) driver.Value {
	switch raw := value.(type) {
	case *sql.NullString:
		if raw.Valid {
			if strings.Contains(typeName, "TIMESTAMP") {
				if t, err := time.Parse(time.RFC3339, raw.String); err == nil {
					return driver.Value{Kind: driver.Timestamp, Time: t, Precision: 6}
				}
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
	case *sql.NullTime:
		if raw.Valid {
			layout := trinoTimeLayout
			switch typeName {
			case "DATE":
				layout = time.DateOnly
			case "TIME WITH TIME ZONE":
				layout = trinoTimeTZLayout
			}
			return driver.Value{Kind: driver.String, String: raw.Time.Format(layout)}
		}
	case *trinodb.Variant:
		// String is the JSON text the server returns for CAST(v AS JSON). A
		// VARIANT null is valid and renders as "null", unlike a SQL NULL.
		if raw.Valid {
			return driver.Value{Kind: driver.String, String: raw.String()}
		}
	}
	return driver.Value{Kind: driver.Null}
}
