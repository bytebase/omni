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
		if !stop {
			result.Plan = describeExplain(result.Statement)
		}
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
	} else if limitable(statement) && opts.Limit > 0 {
		// A plan of a statement the caller did not write answers the wrong
		// question, so the limit applies only when the statement runs.
		statement = applyLimit(statement, opts.Limit)
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
	if !reportsCount(statement) {
		return &driver.Result{Statement: statement}, nil
	}
	affectedRows, err := execResult.RowsAffected()
	if err != nil {
		// Trino does not always report a count, especially for DDL.
		affectedRows = 0
	}
	return &driver.Result{Statement: statement, AffectedRows: &affectedRows}, nil
}

// reportsCount reports whether Trino states an update count for statement.
// Trino answers INSERT, UPDATE, DELETE, MERGE and CREATE TABLE AS with an
// updateCount, which may be 0, and every other statement without one. The
// client reads an absent count as 0 and returns it with no error, so the
// statement's kind, which the parser already holds, is what tells them apart.
// A statement the parser cannot read keeps the client's count: it may be one
// the server accepts and reports.
func reportsCount(statement string) bool {
	file, err := parser.Parse(statement)
	if err != nil || file == nil || len(file.Stmts) != 1 {
		return true
	}
	switch file.Stmts[0].(type) {
	case *parser.InsertStmt, *parser.UpdateStmt, *parser.DeleteStmt, *parser.MergeStmt,
		*parser.CreateTableAsStmt:
		return true
	case *parser.ExecuteStmt, *parser.ExecuteImmediateStmt, *parser.CallStmt:
		// The count depends on what these run, which the text does not state.
		return true
	}
	return false
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

// describeExplain reports the plan a statement answered with, when it is one:
// an explain request rewrites the statement, and a caller's own EXPLAIN runs as
// written. Trino states the plan's format in the statement, and only EXPLAIN
// ANALYZE runs what it plans.
func describeExplain(statement string) *driver.Plan {
	file, err := parser.Parse(statement)
	if err != nil || file == nil || len(file.Stmts) != 1 {
		return nil
	}
	explain, ok := file.Stmts[0].(*parser.ExplainStmt)
	if !ok {
		return nil
	}
	// TYPE VALIDATE answers whether the statement validates, which is one boolean
	// column and not a plan of anything.
	if explain.Type == parser.ExplainTypeValidate {
		return nil
	}
	format := "TEXT"
	switch explain.Format {
	case parser.ExplainFormatJSON:
		format = "JSON"
	case parser.ExplainFormatGraphviz:
		format = "GRAPHVIZ"
	case parser.ExplainFormatUnset:
		// The IO plan is a JSON document unless a FORMAT option says otherwise.
		if explain.Type == parser.ExplainTypeIO {
			format = "JSON"
		}
	}
	return &driver.Plan{Format: format, Executed: explain.Analyze}
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

// codeEnd returns the offset just past statement's last token, before any
// trailing comment or whitespace.
func codeEnd(statement string) int {
	tokens, _ := parser.Tokenize(statement)
	end := 0
	for i, token := range tokens {
		if i == len(tokens)-1 {
			// The list ends with the end-of-input token.
			break
		}
		end = token.Loc.End
	}
	return end
}

// keywordKinds resolves the keywords the row-returning and limit rules read.
var keywordKinds = func() map[string]parser.TokenKind {
	kinds := map[string]parser.TokenKind{}
	for _, name := range []string{"SELECT", "WITH", "TABLE", "VALUES", "SHOW", "DESCRIBE", "DESC", "EXPLAIN"} {
		kind, ok := parser.KeywordToken(name)
		if !ok {
			panic("trino/driver: unknown keyword: " + name)
		}
		kinds[name] = kind
	}
	return kinds
}()

// firstToken returns statement's first token, skipping comments and whitespace.
func firstToken(statement string) parser.Token {
	return parser.NewLexer(statement).NextToken()
}

// returnsRows reports whether statement returns a result set. The SELECT family
// leads with a keyword or a parenthesized query; SHOW, DESCRIBE/DESC and EXPLAIN
// return rows too. SET and RESET SESSION are the SELECT family's exceptions:
// they report no rows, so they run on the exec path.
func returnsRows(statement string) bool {
	token := firstToken(statement)
	if token.Kind == int('(') {
		return true
	}
	for _, name := range []string{"SELECT", "WITH", "TABLE", "VALUES", "SHOW", "DESCRIBE", "DESC", "EXPLAIN"} {
		if token.Kind == keywordKinds[name] {
			return true
		}
	}
	return false
}

// limitable reports whether statement accepts a LIMIT clause: the SELECT family,
// which Trino's SHOW, DESCRIBE and EXPLAIN do not.
func limitable(statement string) bool {
	token := firstToken(statement)
	if token.Kind == int('(') {
		return true
	}
	for _, name := range []string{"SELECT", "WITH", "TABLE", "VALUES"} {
		if token.Kind == keywordKinds[name] {
			return true
		}
	}
	return false
}

// applyLimit caps statement's rows by appending a LIMIT, unless the statement
// holds one of its own. Trino's OFFSET precedes its LIMIT, so the clause goes
// last; it lands just past the last token so a trailing comment does not
// swallow it.
func applyLimit(statement string, limit int) string {
	hasCap, known := rowCap(statement)
	// A statement the parser cannot read is left as written: the rewrite would
	// risk a second clause on a tree we cannot see (Trino accepts syntax the
	// local parser does not, and a syntax error needs no cap).
	if !known || hasCap {
		return statement
	}
	at := codeEnd(statement)
	return fmt.Sprintf("%s LIMIT %d%s", statement[:at], limit, statement[at:])
}

// rowCap reports whether statement already caps its rows with a LIMIT or FETCH
// clause, and whether the parser could tell. It asks the parser because LIMIT
// and FETCH are non-reserved in Trino: `SELECT t.limit FROM t(limit)` names a
// column, and only the tree tells the clause from the name.
func rowCap(statement string) (hasCap, known bool) {
	file, err := parser.Parse(statement)
	if err != nil || file == nil || len(file.Stmts) != 1 {
		return false, false
	}
	query, ok := file.Stmts[0].(*parser.QueryStmt)
	if !ok || query.Query == nil {
		return false, false
	}
	return query.Query.Limit != nil || query.Query.Fetch != nil, true
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
	case "REAL", "DOUBLE":
		return new(sql.NullFloat64)
	// DECIMAL keeps its exact value as text: a float64 cannot hold every
	// decimal Trino returns.
	case "DECIMAL":
		return new(sql.NullString)
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

func convertValue(typeName string, columnType *sql.ColumnType, value any) driver.Value {
	switch raw := value.(type) {
	case *sql.NullString:
		if raw.Valid {
			if strings.Contains(typeName, "TIMESTAMP") {
				if t, err := time.Parse(time.RFC3339, raw.String); err == nil {
					kind := driver.Timestamp
					if strings.Contains(typeName, "TIME ZONE") {
						kind = driver.TimestampTZ
					}
					// A tape recorded before it carried precision has none, so keep the
					// 6 the driver used to report.
					precision := int64(6)
					if p, _, ok := columnType.DecimalSize(); ok {
						precision = p
					}
					return driver.Value{Kind: kind, Time: t, Precision: int32(precision)}
				}
			}
			return driver.Value{Kind: driver.String, String: raw.String}
		}
	case *sql.NullInt64:
		if raw.Valid {
			return driver.Value{Kind: driver.Int64, Int64: raw.Int64}
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
