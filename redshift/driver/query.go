package driver

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/stdlib"

	"github.com/bytebase/omni/driver"
	"github.com/bytebase/omni/redshift"
	redshiftanalysis "github.com/bytebase/omni/redshift/analysis"
	redshiftparser "github.com/bytebase/omni/redshift/parser"
)

// Session opens one server connection held until Close, so session state such
// as the search path and temporary tables persists across Query calls.
func (d *Driver) Session(ctx context.Context) (driver.Session, error) {
	conn, err := d.db.Conn(ctx)
	if err != nil {
		return nil, err
	}
	return &session{driver: d, conn: conn}, nil
}

type session struct {
	driver *Driver
	conn   *sql.Conn
}

func (s *session) Close() error {
	return s.conn.Close()
}

// Query runs each non-empty statement of script in order on the session's
// connection, setting the search path before each one so that a schema set for
// this call does not follow the session into the next. Errors are the client's
// and the server's: pgx's *pgconn.PgError for a statement Redshift refused.
func (s *session) Query(ctx context.Context, script string, opts driver.QueryOptions) ([]*driver.Result, error) {
	var statements []string
	for _, segment := range redshift.Split(script) {
		if !segment.Empty() {
			statements = append(statements, segment.Text)
		}
	}
	if len(statements) == 0 {
		return nil, nil
	}

	type preparedStmt struct {
		text        string
		returnsData bool
	}
	prepared := make([]preparedStmt, 0, len(statements))
	for _, statement := range statements {
		if s.driver.datashare {
			statement = stripDatashareQualifier(statement, s.driver.databaseName)
		}
		if opts.Explain {
			explained, err := redshiftanalysis.ExplainStatement(statement)
			if err != nil {
				return nil, err
			}
			statement = explained
		} else if opts.Limit > 0 {
			statement = redshiftanalysis.StatementWithResultLimit(statement, opts.Limit)
		}

		returnsData, err := classifyQuery(statement)
		if err != nil {
			return nil, err
		}
		prepared = append(prepared, preparedStmt{text: statement, returnsData: returnsData})
	}

	var results []*driver.Result
	for _, stmt := range prepared {
		if opts.Schema != "" {
			if _, err := s.conn.ExecContext(ctx, fmt.Sprintf("set search_path to %s;", quoteIdent(opts.Schema))); err != nil {
				return nil, err
			}
		}

		startTime := time.Now()
		res, err := s.runStatement(ctx, stmt.text, stmt.returnsData, opts)
		stop := false
		if err != nil {
			res = &driver.Result{Err: err}
			stop = true
		}
		res.Statement = stmt.text
		res.Latency = time.Since(startTime)
		if stop {
			results = append(results, res)
			break
		}
		if opts.Explain {
			res.Plan = &driver.Plan{Format: "TEXT"}
		}
		results = append(results, res)
	}
	return results, nil
}

// classifyQuery reports whether statement returns rows, which is how the
// session tells a SELECT from the rest. It reads redshift/analysis's decision
// rather than its own: the analysis also says whether a statement is read-only,
// which a session does not need.
func classifyQuery(statement string) (bool, error) {
	_, returnsData, err := redshiftanalysis.ValidateQuery(statement)
	if err != nil {
		return false, err
	}
	return returnsData, nil
}

// runStatement runs one statement as a query when it returns rows and as a
// command otherwise.
func (s *session) runStatement(ctx context.Context, statement string, returnsData bool, opts driver.QueryOptions) (*driver.Result, error) {
	if returnsData {
		location := pgSessionLocation(s.conn)
		rows, err := s.conn.QueryContext(ctx, statement)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		return readRows(rows, location, opts.MaxResultBytes)
	}

	sqlResult, err := s.conn.ExecContext(ctx, statement)
	if err != nil {
		return nil, err
	}
	affectedRows, err := sqlResult.RowsAffected()
	if err != nil {
		// RowsAffected is only meaningful for DML, and the editor cannot tell
		// it from DDL here.
		affectedRows = 0
	}
	return &driver.Result{AffectedRows: &affectedRows}, nil
}

// readRows decodes rows into a Result, stopping at maxResultBytes when that is
// positive and setting Result.Truncated.
func readRows(rows *sql.Rows, location *time.Location, maxResultBytes int64) (*driver.Result, error) {
	columnNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	columns := make([]driver.Column, len(columnNames))
	metas := make([]columnMeta, len(columnTypes))
	for i, ct := range columnTypes {
		metas[i] = newColumnMeta(ct)
		columns[i] = driver.Column{
			Name:     columnNames[i],
			TypeName: strings.ToUpper(metas[i].databaseTypeName),
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
				row[i] = convertValue(col.TypeName, metas[i], dest[i], location)
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

// makeValueByTypeName returns the scan destination for a column of typeName, as
// Redshift names it.
func makeValueByTypeName(typeName string) any {
	switch typeName {
	case "VARCHAR", "TEXT", "UUID":
		return new(sql.NullString)
	case "BOOL":
		return new(sql.NullBool)
	case "INT", "INTEGER", "TINYINT", "SMALLINT", "MEDIUMINT", "BIGINT", "INT2", "INT4", "INT8":
		return new(sql.NullInt64)
	case "FLOAT", "DOUBLE", "FLOAT4", "FLOAT8":
		return new(sql.NullFloat64)
	case "DATE":
		return new(pgtype.Date)
	case "TIMESTAMP", "TIMESTAMPTZ":
		return new(pgtype.Timestamptz)
	case "BIT", "VARBIT":
		return new([]byte)
	default:
		return new(sql.NullString)
	}
}

var timeTzOID = fmt.Sprintf("%d", pgtype.TimetzOID)

// columnMeta is what the conversions read off a result column besides its
// name: the engine's own type name, and a temporal column's fractional-second
// precision, which pgx reports as the column's scale.
//
// It is a value rather than the *sql.ColumnType it came from so that a
// conversion can be tested without a result set, which a package with no server
// cannot make.
type columnMeta struct {
	databaseTypeName string
	scale            int32
	hasScale         bool
}

// newColumnMeta reads meta from a result column.
func newColumnMeta(columnType *sql.ColumnType) columnMeta {
	_, scale, ok := columnType.DecimalSize()
	return columnMeta{
		databaseTypeName: columnType.DatabaseTypeName(),
		scale:            int32(scale),
		hasScale:         ok,
	}
}

// convertValue maps one scanned cell to the contract's Value. A TimestampTZ is
// rendered in location, the session's zone, and a TIME, TIMETZ or INTERVAL keeps
// the six fractional digits the contract states.
func convertValue(typeName string, meta columnMeta, value any, location *time.Location) driver.Value {
	switch raw := value.(type) {
	case *sql.NullString:
		if raw.Valid {
			// Redshift reports TIME and INTERVAL columns by name, and a TIME
			// WITH TIME ZONE column by its pgx OID.
			dbType := meta.databaseTypeName
			if dbType == "TIME" || dbType == timeTzOID || dbType == "INTERVAL" {
				return driver.Value{Kind: driver.String, String: padZeroes(raw.String, 6)}
			}
			return driver.Value{Kind: driver.String, String: raw.String}
		}
	case *sql.NullInt64:
		if raw.Valid {
			return driver.Value{Kind: driver.Int64, Int64: raw.Int64}
		}
	case *[]byte:
		// An empty BIT or VARBIT value is NULL to the editor, which is what
		// Bytebase's driver reported.
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
	case *pgtype.Date:
		if raw.Valid {
			if raw.InfinityModifier != pgtype.Finite {
				return driver.Value{Kind: driver.String, String: raw.InfinityModifier.String()}
			}
			if typeName == "DATE" {
				return driver.Value{Kind: driver.String, String: raw.Time.Format(time.DateOnly)}
			}
		}
	case *pgtype.Timestamptz:
		if raw.Valid {
			if raw.InfinityModifier != pgtype.Finite {
				return driver.Value{Kind: driver.String, String: raw.InfinityModifier.String()}
			}
			// An unspecified precision is already normalized to PostgreSQL's
			// default of 6 (jackc/pgx#2704).
			scale := meta.scale
			if !meta.hasScale {
				scale = 6
			}
			if typeName == "TIMESTAMP" {
				return driver.Value{Kind: driver.Timestamp, Time: raw.Time, Precision: scale}
			}
			return driver.Value{Kind: driver.TimestampTZ, Time: raw.Time.In(location), Precision: scale}
		}
	}
	return driver.Value{Kind: driver.Null}
}

// padZeroes pads the fractional seconds of a TIME, TIMETZ or INTERVAL text value
// to acc digits, since PostgreSQL omits trailing fractional zeros. The value
// arrives as text because it cannot be formatted into a time.Time.
func padZeroes(rawStr string, acc int) string {
	dotIndex := strings.Index(rawStr, ".")
	if dotIndex < 0 {
		return rawStr
	}
	// End index is used to cut off the time zone information. A value's own
	// sign is not one: an INTERVAL of "-00:02:45.25" starts with it.
	endIndex := len(rawStr)
	if plusIndex := strings.Index(rawStr, "+"); plusIndex >= 0 {
		endIndex = plusIndex
	} else if minusIndex := strings.Index(rawStr, "-"); minusIndex > dotIndex {
		endIndex = minusIndex
	}
	if endIndex <= dotIndex {
		return rawStr
	}
	decimalPart := rawStr[dotIndex+1 : endIndex]
	if len(decimalPart) < acc {
		rawStr = rawStr[:endIndex] + strings.Repeat("0", acc-len(decimalPart)) + rawStr[endIndex:]
	}
	return rawStr
}

var pgLocations sync.Map

// pgSessionLocation reads the session's TimeZone parameter immediately before a
// statement runs, so every cell of its result set renders in one location.
func pgSessionLocation(conn *sql.Conn) *time.Location {
	location := time.UTC
	_ = conn.Raw(func(driverConn any) error {
		pgxConn, ok := driverConn.(*stdlib.Conn)
		if !ok {
			return nil
		}
		location = pgLocation(pgxConn.Conn().PgConn().ParameterStatus("TimeZone"))
		return nil
	})
	return location
}

// pgLocation resolves a TimeZone setting, caching by name because
// time.LoadLocation reparses the zoneinfo file on every call. A POSIX
// specification Go's zoneinfo rejects falls back to UTC, keeping the instant
// correct rather than failing every timestamptz read.
func pgLocation(name string) *time.Location {
	if name == "" || name == "UTC" {
		return time.UTC
	}
	if cached, ok := pgLocations.Load(name); ok {
		if location, ok := cached.(*time.Location); ok {
			return location
		}
	}
	location, err := time.LoadLocation(name)
	if err != nil {
		location = time.UTC
	}
	pgLocations.Store(name, location)
	return location
}

func quoteIdent(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

// stripDatashareQualifier removes the database's name from a shared object's
// qualified name, because a shared object is addressed without the database it
// came from. Only a relation a FROM or JOIN names is rewritten: a literal, a
// comment, a column qualifier and a longer identifier that ends in the same
// text stay as they are.
func stripDatashareQualifier(statement, databaseName string) string {
	if databaseName == "" {
		return statement
	}
	tokens := redshiftparser.Tokenize(statement)
	var b strings.Builder
	last := 0
	// inRelation is true while the tokens are a relation list, which is where
	// the database qualifier belongs. parens records, for each open
	// parenthesis, whether it holds an expression — where a FROM starts an
	// argument (EXTRACT(x FROM y)) rather than a relation list — and the
	// relation list the parenthesis was opened in, which continues after it.
	inRelation := false
	type paren struct {
		expression bool
		inRelation bool
	}
	var parens []paren
	for i, token := range tokens {
		switch token.Type {
		case '(':
			parens = append(parens, paren{expression: i > 0 && opensExpression(tokens[i-1]), inRelation: inRelation})
			continue
		case ')':
			if len(parens) > 0 {
				inRelation = parens[len(parens)-1].inRelation
				parens = parens[:len(parens)-1]
			}
			continue
		}
		if relationListStart[token.Type] {
			inRelation = len(parens) == 0 || !parens[len(parens)-1].expression
			continue
		}
		if relationListEnd[token.Type] {
			inRelation = false
			continue
		}
		if !inRelation || token.Type != redshiftparser.IDENT || i+2 >= len(tokens) {
			continue
		}
		if i > 0 && tokens[i-1].Type == '.' {
			// A schema or column qualifier: the database is a relation's first
			// component, so shared.shared.facts loses only its database.
			continue
		}
		referenced := tokens[i+2]
		if tokens[i+1].Type != '.' || !redshiftparser.IsIdentifierTokenType(referenced.Type) {
			continue
		}
		if token.Loc < len(statement) && statement[token.Loc] == '"' {
			// A quoted identifier keeps its case, so only the same spelling
			// names the database.
			if token.Str != databaseName {
				continue
			}
		} else if !strings.EqualFold(token.Str, databaseName) {
			continue
		}
		b.WriteString(statement[last:token.Loc])
		last = tokens[i+1].End
	}
	if last == 0 {
		return statement
	}
	b.WriteString(statement[last:])
	return b.String()
}

// relationListStart and relationListEnd are the keywords that open and close a
// relation list: a FROM, JOIN, INTO or UPDATE names the relations after it, and
// the next clause or SELECT ends them.
var (
	relationListStart = keywordTokens("FROM", "JOIN", "INTO", "UPDATE")
	relationListEnd   = keywordTokens("SELECT", "WHERE", "ON", "USING", "SET", "HAVING", "QUALIFY", "WINDOW", "RETURNING", "UNION", "EXCEPT", "INTERSECT", "GROUP", "ORDER", "LIMIT", "OFFSET", "FETCH", "FOR", "CONNECT", "START", "VALUES")
)

// keywordTokens maps each keyword's token type, as the lexer reports it.
func keywordTokens(names ...string) map[int]bool {
	tokens := make(map[int]bool, len(names))
	for _, name := range names {
		if keyword := redshiftparser.LookupKeyword(name); keyword != nil {
			tokens[keyword.Token] = true
		}
	}
	return tokens
}

// opensExpression reports whether the token before a parenthesis makes it an
// expression's rather than a subquery's. A name, a literal or a closing
// parenthesis does, so EXTRACT(x FROM y) and TRIM(BOTH ... FROM y) are
// expressions; only the keywords and commas that open a query or a relation
// list do not.
func opensExpression(previous redshiftparser.Token) bool {
	if previous.Type == ',' {
		return false
	}
	return !subqueryParenBefore[previous.Type]
}

// subqueryParenBefore are the tokens a subquery or a grouping follows.
var subqueryParenBefore = keywordTokens("SELECT", "WITH", "VALUES", "FROM", "JOIN", "IN", "EXISTS", "ANY", "ALL", "SOME", "UNION", "INTERSECT", "EXCEPT", "WHERE", "AND", "OR", "NOT", "ON", "USING", "TABLE")
