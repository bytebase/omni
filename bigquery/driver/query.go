package driver

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/big"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"

	"github.com/bytebase/omni/driver"
	googlesqlanalysis "github.com/bytebase/omni/googlesql/analysis"
	"github.com/bytebase/omni/googlesql/parser"
)

// Session returns one session. BigQuery has no session state a statement can
// set, because every statement runs as a job of its own, so a session is the
// driver itself.
func (d *Driver) Session(_ context.Context) (driver.Session, error) {
	return &session{driver: d}, nil
}

type session struct {
	driver *Driver
}

func (*session) Close() error { return nil }

// Query runs each statement of script in order and returns one Result per
// statement run. A statement that fails ends the script with its error on its own
// Result; the error returned is nil, since the script ran as far as it could.
func (s *session) Query(ctx context.Context, script string, opts driver.QueryOptions) ([]*driver.Result, error) {
	if opts.Explain {
		return s.dryRun(ctx, script)
	}

	var results []*driver.Result
	for _, statement := range statements(script) {
		startTime := time.Now()
		result, err := s.run(ctx, statement, opts)
		if err != nil {
			result = &driver.Result{Err: err}
		}
		// A query reports the text the rewrite left, which is the text BigQuery
		// ran.
		if result.Statement == "" {
			result.Statement = statement
		}
		result.Latency = time.Since(startTime)
		results = append(results, result)
		if result.Err != nil {
			break
		}
	}
	return results, nil
}

// statements splits script at its top-level semicolons with the BigQuery
// splitter, which takes every ';' outside a literal, a quoted identifier and a
// comment as a statement boundary. A segment that holds no statement asks
// BigQuery nothing and is skipped.
func statements(script string) []string {
	segments := parser.SplitFlat(script)
	texts := make([]string, 0, len(segments))
	for _, segment := range segments {
		if segment.Empty() {
			continue
		}
		texts = append(texts, strings.TrimSpace(segment.Text))
	}
	return texts
}

// returnsRows reports whether statement returns a result set, which
// googlesql/analysis decides for the BigQuery dialect. A statement it does not
// classify as a query is run as a job that reports rows affected, as Bytebase's
// driver did; one it cannot classify at all takes that path too, and
// session.command reads the result set when the service states one.
func returnsRows(statement string) bool {
	switch googlesqlanalysis.ClassifySQL(statement, googlesqlanalysis.DialectBigQuery) {
	case googlesqlanalysis.Select, googlesqlanalysis.SelectInfoSchema:
		return true
	default:
		return false
	}
}

// limitStatement wraps statement so BigQuery returns at most limit rows, through
// googlesql/analysis, which is BigQuery's CTE wrap: BigQuery has no LIMIT to
// splice into an arbitrary query.
func limitStatement(statement string, limit int) string {
	return googlesqlanalysis.StatementWithResultLimit(statement, limit, googlesqlanalysis.DialectBigQuery)
}

// run runs one statement: a query reads it, and anything else runs as a job whose
// statistics report the rows it affected.
func (s *session) run(ctx context.Context, statement string, opts driver.QueryOptions) (*driver.Result, error) {
	if !returnsRows(statement) {
		return s.command(ctx, statement, opts, unclassified(statement))
	}
	text := statement
	if opts.Limit > 0 {
		text = limitStatement(text, opts.Limit)
	}
	result, err := s.queryStatement(ctx, text, opts)
	if err != nil {
		return nil, err
	}
	result.Statement = text
	return result, nil
}

// queryStatement reads one statement's rows.
func (s *session) queryStatement(ctx context.Context, statement string, opts driver.QueryOptions) (*driver.Result, error) {
	it, err := s.driver.client.Read(ctx, s.driver.databaseName, statement)
	if err != nil {
		return nil, err
	}
	return s.readRows(ctx, it, opts)
}

// readRows decodes the rows the iterator returns, up to opts.Limit rows and
// opts.MaxResultBytes bytes, and carries the plan of the job that returned them.
// A query the analysis read was rewritten to cap its own rows; a statement it
// could not classify runs as it was written, so the cap is applied here too.
func (s *session) readRows(ctx context.Context, it rowIterator, opts driver.QueryOptions) (*driver.Result, error) {
	result := &driver.Result{}
	var totalBytes int64
	for {
		values, err := it.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		schema := it.Schema()
		// The schema arrives with the first row the client decodes, not with the
		// iterator.
		if result.Columns == nil && len(schema) > 0 {
			result.Columns = columnsOf(schema)
			totalBytes = driver.HeaderByteSize(result.Columns)
		}
		row := make([]driver.Value, 0, len(values))
		for i, value := range values {
			var fieldType bigquery.FieldType
			if i < len(schema) {
				fieldType = schema[i].Type
			}
			row = append(row, convertValue(value, fieldType))
		}
		result.Rows = append(result.Rows, row)
		if opts.Limit > 0 && len(result.Rows) >= opts.Limit {
			break
		}
		if opts.MaxResultBytes > 0 {
			totalBytes += driver.RowByteSize(row)
			if totalBytes > opts.MaxResultBytes {
				result.Truncated = true
				break
			}
		}
	}
	// A query that returned no rows has no row to learn the schema from, and its
	// columns are still the ones the statement declares.
	if result.Columns == nil && len(it.Schema()) > 0 {
		result.Columns = columnsOf(it.Schema())
	}

	// The job that returned these rows reports the stages it ran as it goes;
	// BigQuery has no EXPLAIN to ask for them instead.
	result.Plan = queryPlan(ctx, it.SourceJob())
	return result, nil
}

// columnsOf maps a result schema onto the contract's columns.
func columnsOf(schema bigquery.Schema) []driver.Column {
	columns := make([]driver.Column, 0, len(schema))
	for _, field := range schema {
		if field == nil {
			continue
		}
		columns = append(columns, driver.Column{
			Name:     field.Name,
			TypeName: strings.ToUpper(string(field.Type)),
		})
	}
	return columns
}

// command runs a statement that returns no rows and reports the rows it affected.
//
// A statement googlesql/analysis cannot read runs this way too, and omni's
// parser has gaps a valid statement may hit (a pipe query is one). When the job
// answered a result set, its rows are read rather than dropped; one that answered
// none reports the rows it affected, and reading the job's result is what tells the
// two apart.
func (s *session) command(ctx context.Context, statement string, opts driver.QueryOptions, unclassified bool) (*driver.Result, error) {
	job, err := s.driver.client.Run(ctx, s.driver.databaseName, statement)
	if err != nil {
		return nil, err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return nil, err
	}
	if status == nil {
		return nil, errors.New("the job's status carries no query statistics")
	}
	if status.Err != nil {
		return nil, status.Err
	}
	if status.Stats == nil {
		return nil, errors.New("the job's status carries no query statistics")
	}
	it, err := job.Read(ctx)
	if err != nil {
		return nil, err
	}
	if jobReturnsRows(status.Stats.StatementType, len(it.Schema()) > 0, unclassified) {
		return s.readRows(ctx, it, opts)
	}
	affected := status.Stats.NumDMLAffectedRows
	return &driver.Result{AffectedRows: &affected}, nil
}

// selectStatementType is the statement type BigQuery states for a query that
// answered rows.
const selectStatementType = "SELECT"

// unclassified reports whether googlesql/analysis could read nothing about
// statement, which is when the service's own answer is what says whether it
// returned rows: a pipe query, a CALL and a BEGIN block are all valid BigQuery
// that omni's parser does not read yet.
func unclassified(statement string) bool {
	return googlesqlanalysis.ClassifySQL(statement, googlesqlanalysis.DialectBigQuery) == googlesqlanalysis.Unknown
}

// jobReturnsRows reports whether a finished job answered a result set.
//
// A plain query's job statistics state no result schema at all, so its statement
// type is what says it is a query, which is the answer the client's own storage
// read uses; it says so whether the query returned rows or none.
//
// A statement the analysis read as DDL or DML is a command: a DDL job states a
// result schema of its own that no rows came back for (a CREATE TABLE stated one
// column), and a DML states the count it affected, so neither is read as rows.
//
// A statement the analysis could read nothing about is the service's to answer. Its
// result schema is what says it answered a result set — one that came back with no
// rows still reports its columns — and a statement that answered none, as an ASSERT
// does, reports zero rows affected.
func jobReturnsRows(statementType string, hasResultSchema bool, unclassified bool) bool {
	if statementType == selectStatementType {
		return true
	}
	return unclassified && hasResultSchema
}

// dryRun answers what a statement would do without running it, which is what
// BigQuery offers in place of EXPLAIN: whether the statement is valid, and how
// many bytes it would process. Each statement's answer is one row of a STRING
// column named "Validation Result", the shape Bytebase's driver returned.
//
// A statement the service refuses ends the call with its error on its own
// Result, which is what a query's failure does.
func (s *session) dryRun(ctx context.Context, script string) ([]*driver.Result, error) {
	var results []*driver.Result
	for _, statement := range statements(script) {
		startTime := time.Now()
		result, err := s.dryRunStatement(ctx, statement)
		if err != nil {
			result = &driver.Result{Err: err}
		}
		result.Statement = statement
		result.Latency = time.Since(startTime)
		results = append(results, result)
		if result.Err != nil {
			break
		}
	}
	return results, nil
}

// dryRunStatement asks the service to validate one statement.
func (s *session) dryRunStatement(ctx context.Context, statement string) (*driver.Result, error) {
	status, err := s.driver.client.DryRun(ctx, s.driver.databaseName, statement)
	if err != nil {
		return nil, err
	}
	if status == nil {
		return nil, errors.New("failed to read the dry run's job status")
	}
	if status.Err != nil {
		return nil, status.Err
	}
	result := &driver.Result{}
	if status.Stats == nil {
		return result, nil
	}
	result.Columns = []driver.Column{{Name: "Validation Result", TypeName: "STRING"}}
	result.Rows = [][]driver.Value{{{
		Kind: driver.String,
		String: fmt.Sprintf(
			"Query successfully validated. Assuming the tables are not modified, "+
				"running this query will process %d bytes (%.2f MB).",
			status.Stats.TotalBytesProcessed,
			float64(status.Stats.TotalBytesProcessed)/(1024*1024),
		),
	}}}
	return result, nil
}

// convertValue converts one cell of a result row. A cell whose declared type has
// no kind here falls back to Go's own rendering of the value, which is what
// Bytebase's driver returned.
//
// A NUMERIC and a BIGNUMERIC arrive as a *big.Rat, whose own rendering is a ratio
// ("25/2" for 12.5), so they are written as the exact decimal text BigQuery's
// literal states. A TIMESTAMP is an instant and a DATETIME a wall clock, which are
// the contract's TimestampTZ and Timestamp; both hold microseconds, and BigQuery
// states no other scale for either. A RANGE is the half-open literal the service
// states, because Go's rendering of a *bigquery.RangeValue is "&{2026-01-01
// 2026-02-01}", which names neither the type nor the interval.
//
// DATE, TIME, JSON, GEOGRAPHY and INTERVAL keep Go's rendering. An ARRAY and a
// RECORD are rendered element by element, because Go's rendering of the client's
// []bigquery.Value prints a pointer element as an address: an ARRAY<RANGE> cell
// arrived as "[0x27d0b631da40]", which is neither readable nor reproducible. A
// RECORD stays lossy in a way this cannot fix: the client hands over its values in
// the field's order with no names, and the names live in the field's nested schema,
// which this function is not given, so a record reads as a list.
func convertValue(value bigquery.Value, fieldType bigquery.FieldType) driver.Value {
	if value == nil {
		return driver.Value{Kind: driver.Null}
	}
	if elements, ok := value.([]bigquery.Value); ok {
		return driver.Value{Kind: driver.String, String: elementsText(elements)}
	}
	switch fieldType {
	case bigquery.StringFieldType:
		if s, ok := value.(string); ok {
			return driver.Value{Kind: driver.String, String: s}
		}
	case bigquery.BytesFieldType:
		if bytes, ok := value.([]byte); ok {
			return driver.Value{Kind: driver.Bytes, Bytes: bytes}
		}
	case bigquery.IntegerFieldType:
		if i, ok := value.(int64); ok {
			return driver.Value{Kind: driver.Int64, Int64: i}
		}
	case bigquery.FloatFieldType:
		if f, ok := value.(float64); ok {
			return driver.Value{Kind: driver.Float64, Float64: f}
		}
	case bigquery.BooleanFieldType:
		if b, ok := value.(bool); ok {
			return driver.Value{Kind: driver.Bool, Bool: b}
		}
	case bigquery.NumericFieldType, bigquery.BigNumericFieldType:
		if rat, ok := value.(*big.Rat); ok {
			return driver.Value{Kind: driver.String, String: decimalText(rat)}
		}
	case bigquery.TimestampFieldType:
		if instant, ok := value.(time.Time); ok {
			return driver.Value{Kind: driver.TimestampTZ, Time: instant, Precision: bigqueryFractionalDigits}
		}
	case bigquery.DateTimeFieldType:
		if wallClock, ok := value.(civil.DateTime); ok {
			return driver.Value{Kind: driver.Timestamp, Time: wallClock.In(time.UTC), Precision: bigqueryFractionalDigits}
		}
	case bigquery.RangeFieldType:
		if r, ok := value.(*bigquery.RangeValue); ok {
			return driver.Value{Kind: driver.String, String: rangeText(r)}
		}
	default:
		// A DATE, TIME, JSON, GEOGRAPHY and INTERVAL fall back to Go's own
		// rendering of the value.
	}
	return driver.Value{Kind: driver.String, String: fmt.Sprintf("%v", value)}
}

// rangeText renders a RANGE as the literal the service states: a half-open
// interval, "[START, END)", with UNBOUNDED for a bound the range does not have.
// BigQuery accepts no other form: a closed END is refused with "range must be
// formatted exactly as [START, END)".
//
// A bound keeps the fractional seconds the value holds. The service's own
// rendering of a RANGE<TIMESTAMP> drops them, and that loss is not copied here.
func rangeText(r *bigquery.RangeValue) string {
	return "[" + rangeBoundText(r.Start) + ", " + rangeBoundText(r.End) + ")"
}

// rangeBoundText renders one bound, and UNBOUNDED for the one a range does not
// have.
func rangeBoundText(bound bigquery.Value) string {
	if bound == nil {
		return "UNBOUNDED"
	}
	return elementText(bound)
}

// elementsText renders an ARRAY or a RECORD: its elements in order, each one by
// the rule the cell of its own type states, inside one pair of brackets.
func elementsText(elements []bigquery.Value) string {
	texts := make([]string, 0, len(elements))
	for _, element := range elements {
		texts = append(texts, elementText(element))
	}
	return "[" + strings.Join(texts, ", ") + "]"
}

// elementText renders one element of an ARRAY or a RECORD, or one bound of a
// RANGE. The Go type the client decoded says which it is, since an element has no
// field of its own here: a *big.Rat is a decimal, a *bigquery.RangeValue its
// literal, a civil or time value the text the service states for it, and a nested
// []bigquery.Value another array or record.
func elementText(element bigquery.Value) string {
	switch e := element.(type) {
	case nil:
		return "NULL"
	case *big.Rat:
		return decimalText(e)
	case *bigquery.RangeValue:
		return rangeText(e)
	case civil.Date:
		return e.String()
	case civil.Time:
		return e.String()
	case civil.DateTime:
		return e.In(time.UTC).Format("2006-01-02T15:04:05.999999")
	case time.Time:
		return e.UTC().Format("2006-01-02 15:04:05.999999")
	case []bigquery.Value:
		return elementsText(e)
	}
	return fmt.Sprintf("%v", element)
}

// bigqueryFractionalDigits is the fractional-second digits a BigQuery TIMESTAMP
// and DATETIME hold. The service states no per-column scale for either.
const bigqueryFractionalDigits = 6

// decimalText returns rat's exact decimal text: the shortest string of digits that
// names it, which is how BigQuery states a NUMERIC and a BIGNUMERIC. A denominator
// the client only ever builds from a decimal literal is a product of twos and
// fives, so the text is finite; one that is not is rounded to the digits the
// service allows rather than rendered as a ratio.
func decimalText(rat *big.Rat) string {
	denominator := new(big.Int).Set(rat.Denom())
	two, five, remainder := big.NewInt(2), big.NewInt(5), new(big.Int)
	scale := 0
	for denominator.Cmp(big.NewInt(1)) != 0 {
		switch {
		case remainder.Mod(denominator, two).Sign() == 0:
			denominator.Div(denominator, two)
		case remainder.Mod(denominator, five).Sign() == 0:
			denominator.Div(denominator, five)
		default:
			return rat.FloatString(scale + 38)
		}
		scale++
	}
	text := rat.FloatString(scale)
	if !strings.Contains(text, ".") {
		return text
	}
	text = strings.TrimRight(text, "0")
	return strings.TrimSuffix(text, ".")
}
