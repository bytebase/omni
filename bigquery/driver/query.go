package driver

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"

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
		text := strings.TrimSpace(segment.Text)
		if text == "" {
			continue
		}
		texts = append(texts, text)
	}
	return texts
}

// returnsRows reports whether statement returns a result set, which
// googlesql/analysis decides for the BigQuery dialect. A statement it does not
// classify as a query is run as a job that reports rows affected, as Bytebase's
// driver did; one it cannot classify at all takes the same path, and the service
// is then what refuses it.
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
		return s.command(ctx, statement)
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
	it, err := s.driver.query(statement).Read(ctx)
	if err != nil {
		return nil, err
	}
	result := &driver.Result{}
	var totalBytes int64
	for {
		var values []bigquery.Value
		err := it.Next(&values)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		// The schema arrives with the first row the client decodes, not with the
		// iterator.
		if result.Columns == nil && len(it.Schema) > 0 {
			result.Columns = columnsOf(it.Schema)
			totalBytes = driver.HeaderByteSize(result.Columns)
		}
		row := make([]driver.Value, 0, len(values))
		for i, value := range values {
			var fieldType bigquery.FieldType
			if i < len(it.Schema) {
				fieldType = it.Schema[i].Type
			}
			row = append(row, convertValue(value, fieldType))
		}
		result.Rows = append(result.Rows, row)
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
	if result.Columns == nil && len(it.Schema) > 0 {
		result.Columns = columnsOf(it.Schema)
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
func (s *session) command(ctx context.Context, statement string) (*driver.Result, error) {
	job, err := s.driver.query(statement).Run(ctx)
	if err != nil {
		return nil, err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return nil, err
	}
	if err := status.Err(); err != nil {
		return nil, err
	}
	stats, ok := statisticsOf(status)
	if !ok {
		return nil, errors.New("the job's status carries no query statistics")
	}
	affected := stats.NumDMLAffectedRows
	return &driver.Result{AffectedRows: &affected}, nil
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
	q := s.driver.query(statement)
	q.DryRun = true
	job, err := q.Run(ctx)
	if err != nil {
		return nil, err
	}
	// A dry run completes before Run returns, so its status is the finished one.
	status := job.LastStatus()
	if status == nil {
		return nil, errors.New("failed to read the dry run's job status")
	}
	if err := status.Err(); err != nil {
		return nil, err
	}
	result := &driver.Result{}
	stats, ok := statisticsOf(status)
	if !ok {
		return result, nil
	}
	result.Columns = []driver.Column{{Name: "Validation Result", TypeName: "STRING"}}
	result.Rows = [][]driver.Value{{{
		Kind: driver.String,
		String: fmt.Sprintf(
			"Query successfully validated. Assuming the tables are not modified, "+
				"running this query will process %d bytes (%.2f MB).",
			stats.TotalBytesProcessed,
			float64(stats.TotalBytesProcessed)/(1024*1024),
		),
	}}}
	return result, nil
}

// convertValue converts one cell of a result row. A cell whose declared type has
// no kind here falls back to Go's own rendering of the value, which is what
// Bytebase's driver returned.
func convertValue(value bigquery.Value, fieldType bigquery.FieldType) driver.Value {
	if value == nil {
		return driver.Value{Kind: driver.Null}
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
	default:
		// A TIMESTAMP, DATE, TIME, DATETIME, NUMERIC, BIGNUMERIC, GEOGRAPHY,
		// INTERVAL, JSON, RANGE, a RECORD and an ARRAY fall back to Go's own
		// rendering of the value.
	}
	return driver.Value{Kind: driver.String, String: fmt.Sprintf("%v", value)}
}
