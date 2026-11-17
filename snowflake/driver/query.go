package driver

import (
	"context"
	"database/sql"
	sqldriver "database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/snowflakedb/gosnowflake/v2"

	"github.com/bytebase/omni/driver"
	snowflakeanalysis "github.com/bytebase/omni/snowflake/analysis"
)

// Session opens one client connection, held until Close, so session state such
// as the current schema and a temporary table persists across Query calls.
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
	discardConn(s.conn)
	return nil
}

// discardConn closes conn and drops its physical connection rather than
// returning it to the pool. A Snowflake session carries the database, schema,
// role, warehouse and transaction state its statements set, and the next
// borrower of the same connection would inherit them: a script that ran USE
// DATABASE would decide where another caller's unqualified write lands.
func discardConn(conn *sql.Conn) {
	// database/sql closes the physical connection when Raw's function reports
	// that the connection is unusable.
	_ = conn.Raw(func(any) error { return sqldriver.ErrBadConn })
	_ = conn.Close()
}

// Query runs each statement of script in order on the session's connection,
// setting the current schema before each one so that a schema set for this call
// does not follow the session into the next. Errors are the client's and the
// server's: gosnowflake's *gosnowflake.SnowflakeError for a statement Snowflake
// refused.
func (s *session) Query(ctx context.Context, script string, opts driver.QueryOptions) ([]*driver.Result, error) {
	var statements []string
	for _, statement := range snowflakeanalysis.SplitStatements(script) {
		if statement.Empty {
			continue
		}
		statements = append(statements, statement.Text)
	}
	if len(statements) == 0 {
		return nil, nil
	}

	type preparedStatement struct {
		text        string
		returnsData bool
	}
	prepared := make([]preparedStatement, 0, len(statements))
	for _, statement := range statements {
		if opts.Explain {
			explained, err := snowflakeanalysis.ExplainStatement(statement, strings.EqualFold(opts.ExplainFormat, "JSON"))
			if err != nil {
				return nil, err
			}
			statement = explained
		} else if opts.Limit > 0 {
			if rewritten, err := snowflakeanalysis.RewriteResultLimit(statement, opts.Limit); err == nil {
				statement = rewritten
			} else {
				// The newline ends a trailing line comment inside the wrapper,
				// which would otherwise swallow the closing parenthesis.
				statement = fmt.Sprintf("SELECT * FROM (%s\n) LIMIT %d", trimStatement(statement), opts.Limit)
			}
		}
		returns, err := classify(statement, opts.MaskingEnabled)
		if err != nil {
			return nil, err
		}
		prepared = append(prepared, preparedStatement{text: statement, returnsData: returns})
	}

	var results []*driver.Result
	for _, statement := range prepared {
		if opts.Schema != "" {
			if _, err := s.conn.ExecContext(ctx, "USE SCHEMA "+quoteIdent(opts.Schema)); err != nil {
				return nil, err
			}
		} else {
			// A session whose earlier statement ran USE DATABASE would otherwise
			// decide where this call's unqualified names resolve. PUBLIC may not
			// exist in the database, and there is nothing to do about it here.
			_, _ = s.conn.ExecContext(ctx, "USE SCHEMA PUBLIC")
		}

		startTime := time.Now()
		result, err := s.runStatement(ctx, statement.text, statement.returnsData, opts)
		if err != nil {
			result = &driver.Result{Err: err}
		}
		result.Statement = statement.text
		result.Latency = time.Since(startTime)
		if format, executed, ok := snowflakeanalysis.DescribeExplain(statement.text); ok {
			result.Plan = &driver.Plan{Format: format, Executed: executed}
		}
		results = append(results, result)
		if err != nil {
			break
		}
	}
	return results, nil
}

// classify reports whether the session runs statement as a query rather than as
// a command reporting affected rows. The decision is snowflake/analysis's; this
// function holds the wiring to it.
//
// A statement omni cannot parse runs as a query, as Bytebase ran it: the answer
// is a hint about how to send a statement, and the server decides whether the
// statement is valid. A caller that masks columns by tracing them through the
// statement gets the parse error instead, because a statement the analysis
// could not read is one whose answer it could not withhold rows from.
//
// A procedure call is a query when masking is off: CALL returns a value or a
// table, which the analysis counts as a statement that returns nothing.
func classify(statement string, maskingEnabled bool) (bool, error) {
	_, returns, err := snowflakeanalysis.ValidateQuery(statement)
	if err != nil {
		if maskingEnabled {
			return false, err
		}
		return true, nil
	}
	if isProcedureCall(statement) {
		return !maskingEnabled, nil
	}
	return returns, nil
}

// isProcedureCall reports whether statement calls a stored procedure.
func isProcedureCall(statement string) bool {
	fields := strings.Fields(statement)
	return len(fields) > 0 && strings.EqualFold(fields[0], "CALL")
}

// statementContext returns ctx with the client's one-statement limit lifted. A
// scripting block (BEGIN ... END) is one statement to the splitter and several to
// the wire, and the client refuses more than one per request unless told to
// allow them, which is what Bytebase's driver asked for too.
func statementContext(ctx context.Context) context.Context {
	return gosnowflake.WithMultiStatement(ctx, 0)
}

// trimStatement trims a statement's surrounding whitespace and its trailing
// semicolons, so it can be wrapped in a subquery.
func trimStatement(statement string) string {
	return strings.TrimLeftFunc(strings.TrimRightFunc(statement, isSpaceOrSemicolon), unicode.IsSpace)
}

// isSpaceOrSemicolon reports whether r ends a statement without being part of
// one.
func isSpaceOrSemicolon(r rune) bool {
	return unicode.IsSpace(r) || r == ';'
}

// runStatement runs one statement as a query when it returns rows and as a
// command otherwise.
func (s *session) runStatement(ctx context.Context, statement string, returnsData bool, opts driver.QueryOptions) (*driver.Result, error) {
	if returnsData {
		rows, err := s.conn.QueryContext(statementContext(ctx), statement)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		return readRows(rows, opts.MaxResultBytes)
	}

	sqlResult, err := s.conn.ExecContext(statementContext(ctx), statement)
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
	metas := make([]columnMeta, len(columnTypes))
	for i, columnType := range columnTypes {
		metas[i] = newColumnMeta(columnType)
		columns[i] = driver.Column{
			Name:     columnNames[i],
			TypeName: strings.ToUpper(metas[i].typeName),
		}
	}
	result := &driver.Result{Columns: columns}
	totalBytes := driver.HeaderByteSize(columns)
	if len(columns) == 0 {
		return result, nil
	}

	for rows.Next() {
		// One destination for every column: the client answers a column's value
		// as the Go type its own type map picks, and the conversion below reads
		// that type together with the column's type name.
		values := make([]any, len(columns))
		dest := make([]any, len(columns))
		for i := range values {
			dest[i] = &values[i]
		}
		if err := rows.Scan(dest...); err != nil {
			return nil, err
		}
		row := make([]driver.Value, len(columns))
		for i := range columns {
			row[i] = convertValue(metas[i], values[i])
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

// columnMeta is what the conversions read off a result column besides its name:
// the client's type name, and the exactness the client states for a numeric or
// temporal column.
//
// It is a value rather than the *sql.ColumnType it came from so that a
// conversion can be tested without a result set, which a package with no
// account cannot make.
type columnMeta struct {
	// typeName is the client's own name for the column's type: FIXED, REAL,
	// TEXT, DATE, TIME, TIMESTAMP_NTZ, TIMESTAMP_LTZ, TIMESTAMP_TZ, VARIANT,
	// OBJECT, ARRAY, BINARY, BOOLEAN or DECFLOAT.
	typeName string
	// precision and scale are what the client reports for the type. For a FIXED
	// column they are the declared precision and scale; for a TIME column the
	// fractional-second digits arrive as the precision; for a TIMESTAMP column
	// the client reports neither.
	precision, scale int32
	hasPrecision     bool
}

// newColumnMeta reads meta from a result column.
func newColumnMeta(columnType *sql.ColumnType) columnMeta {
	precision, scale, ok := columnType.DecimalSize()
	return columnMeta{
		typeName:     columnType.DatabaseTypeName(),
		precision:    int32(precision),
		scale:        int32(scale),
		hasPrecision: ok,
	}
}

// fractionalDigits is the number of fractional-second digits a TIME or
// TIMESTAMP column declares. The client states a TIME column's as its
// precision, and none for a TIMESTAMP column, which Snowflake defaults to 9.
func (m columnMeta) fractionalDigits() int32 {
	if m.typeName == "TIME" && m.hasPrecision {
		return m.precision
	}
	return 9
}

// convertValue maps one scanned cell to the contract's Value.
//
// The client answers a column as the Go type its type map picks: a string for a
// FIXED, REAL, TEXT, VARIANT, OBJECT, ARRAY, DECFLOAT or BOOLEAN value (unless
// the caller asked for higher-precision or structured types), a time.Time for a
// DATE, TIME or TIMESTAMP, and []byte for a BINARY. A FIXED value keeps its
// exact digits as text when it declares a scale, as a decimal column does in
// trino's and SQL Server's drivers.
func convertValue(meta columnMeta, value any) driver.Value {
	switch raw := value.(type) {
	case nil:
		return driver.Value{Kind: driver.Null}
	case int64:
		return driver.Value{Kind: driver.Int64, Int64: raw}
	case float64:
		return driver.Value{Kind: driver.Float64, Float64: raw}
	case bool:
		return driver.Value{Kind: driver.Bool, Bool: raw}
	case []byte:
		return driver.Value{Kind: driver.Bytes, Bytes: raw}
	case time.Time:
		return convertTime(meta, raw)
	case string:
		return convertText(meta, raw)
	default:
		// A value kind this driver does not know — a structured OBJECT or ARRAY
		// a caller asked the client to build — is reported as its text rather
		// than dropped.
		return driver.Value{Kind: driver.String, String: fmt.Sprint(raw)}
	}
}

// convertText maps a value the client stated as text.
func convertText(meta columnMeta, text string) driver.Value {
	switch meta.typeName {
	case "FIXED":
		// A FIXED with a scale is a decimal, and its digits are exact only as
		// text.
		if meta.scale == 0 {
			if number, err := strconv.ParseInt(text, 10, 64); err == nil {
				return driver.Value{Kind: driver.Int64, Int64: number}
			}
		}
	case "REAL":
		if number, err := strconv.ParseFloat(text, 64); err == nil {
			return driver.Value{Kind: driver.Float64, Float64: number}
		}
	case "BOOLEAN":
		if boolean, err := strconv.ParseBool(text); err == nil {
			return driver.Value{Kind: driver.Bool, Bool: boolean}
		}
	}
	return driver.Value{Kind: driver.String, String: text}
}

// convertTime maps a value the client stated as an instant or a wall clock.
func convertTime(meta columnMeta, value time.Time) driver.Value {
	switch meta.typeName {
	case "DATE":
		return driver.Value{Kind: driver.String, String: value.Format(time.DateOnly)}
	case "TIME":
		return driver.Value{Kind: driver.String, String: formatClock(value, meta.fractionalDigits())}
	case "TIMESTAMP_NTZ":
		// A wall clock, which the client states in UTC fields.
		return driver.Value{Kind: driver.Timestamp, Time: value, Precision: meta.fractionalDigits()}
	case "TIMESTAMP_LTZ", "TIMESTAMP_TZ":
		// An instant: the client renders LTZ in the session's zone and keeps
		// TZ's own offset.
		return driver.Value{Kind: driver.TimestampTZ, Time: value, Precision: meta.fractionalDigits()}
	default:
		return driver.Value{Kind: driver.String, String: value.Format(time.RFC3339Nano)}
	}
}

// formatClock writes a TIME value, whose own date is a stand-in, with digits
// fractional-second digits.
func formatClock(value time.Time, digits int32) string {
	clock := value.Format("15:04:05")
	if digits <= 0 {
		return clock
	}
	nanoseconds := fmt.Sprintf("%09d", value.Nanosecond())
	if digits > int32(len(nanoseconds)) {
		digits = int32(len(nanoseconds))
	}
	return clock + "." + nanoseconds[:digits]
}
