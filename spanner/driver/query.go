package driver

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	sppb "cloud.google.com/go/spanner/apiv1/spannerpb"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/bytebase/omni/driver"
	spanneranalysis "github.com/bytebase/omni/spanner/analysis"
	"github.com/bytebase/omni/spanner/parser"
)

// Session returns one session. Spanner has no session state a statement can
// set, so a session is the driver itself.
func (d *Driver) Session(_ context.Context) (driver.Session, error) {
	return &session{driver: d}, nil
}

type session struct {
	driver *Driver
}

func (*session) Close() error { return nil }

// statementKind is how the driver runs one statement.
type statementKind int

const (
	// query is read in a single-use read-only transaction.
	query statementKind = iota
	// ddl is applied through the database admin API.
	ddl
	// dml is run in a read-write transaction.
	dml
)

// statement is one statement of a script and how the driver runs it.
type statement struct {
	Text string
	Kind statementKind
}

// statements splits script into the statements the driver runs, in the order
// they are written, by GoogleSQL's grammar. A segment that holds only a
// comment, and asks Spanner nothing, is skipped.
func statements(script string) []statement {
	var statements []statement
	for _, segment := range parser.Split(script) {
		// spanner/parser's Split marks a segment that holds no statement.
		if segment.Empty() {
			continue
		}
		text := strings.TrimSpace(segment.Text)
		if text == "" {
			continue
		}
		kind, ok := classifyStatement(text)
		if !ok {
			continue
		}
		statements = append(statements, statement{Text: text, Kind: kind})
	}
	return statements
}

// classifyStatement routes one statement by its first keyword and reports false
// for a segment that holds no statement. Comments before the keyword are the
// lexer's to skip.
//
// The keywords are Bytebase's, with one addition: Spanner's GRANT and REVOKE are
// DDL statements the admin API takes, and Bytebase's prefix test sent them down
// the DML path, where Spanner refuses them as "not valid DML". A statement that
// changes privileges changes the schema's access, so it is a schema change.
func classifyStatement(text string) (statementKind, bool) {
	tokens := firstTokens(text, 1)
	if len(tokens) == 0 {
		return 0, false
	}
	switch strings.ToUpper(tokens[0]) {
	case "SELECT", "WITH":
		return query, true
	case "CREATE", "DROP", "ALTER", "GRANT", "REVOKE":
		return ddl, true
	default:
		return dml, true
	}
}

// firstTwoKeywords returns text's first two keywords as they are written, upper
// case, with the comments before and between them skipped. A text with fewer
// words answers what it has.
func firstTwoKeywords(text string) [2]string {
	var keywords [2]string
	tokens := firstTokens(text, len(keywords))
	for i := range keywords {
		if i < len(tokens) {
			keywords[i] = strings.ToUpper(tokens[i])
		}
	}
	return keywords
}

// firstTokens returns text's first count tokens as they are written, with
// GoogleSQL's lexer skipping its comments, and stops at the end of the text.
func firstTokens(text string, count int) []string {
	tokens := make([]string, 0, count)
	lexer := parser.NewLexer(text)
	for len(tokens) < count {
		tok := lexer.NextToken()
		if tok.Loc.End <= tok.Loc.Start {
			break
		}
		tokens = append(tokens, text[tok.Loc.Start:tok.Loc.End])
	}
	return tokens
}

// Query runs each statement of script in order and returns one Result per
// statement run. A statement that fails ends the script with its error on its
// own Result; the error returned is nil, since the script ran as far as it
// could. A statement's error is the client's, as the package doc describes.
func (s *session) Query(ctx context.Context, script string, opts driver.QueryOptions) ([]*driver.Result, error) {
	if opts.Explain {
		return s.explainStatement(ctx, script, opts)
	}

	var results []*driver.Result
	for _, statement := range statements(script) {
		startTime := time.Now()
		result, err := s.run(ctx, statement, opts)
		stop := false
		if err != nil {
			result = &driver.Result{Err: err}
			stop = true
		}
		// A query reports the text the rewrite left, which is the text Spanner
		// ran.
		if result.Statement == "" {
			result.Statement = statement.Text
		}
		result.Latency = time.Since(startTime)
		results = append(results, result)
		if stop {
			break
		}
	}
	return results, nil
}

// run runs one statement: a query reads it, a DDL statement applies it through
// the admin API, and DML runs it in a read-write transaction.
func (s *session) run(ctx context.Context, statement statement, opts driver.QueryOptions) (*driver.Result, error) {
	switch statement.Kind {
	case query:
		text := statement.Text
		if opts.Limit > 0 {
			text = spanneranalysis.RewriteResultLimit(text, opts.Limit)
		}
		result, err := s.queryStatement(ctx, text, opts)
		if result != nil {
			result.Statement = text
		}
		return result, err
	case ddl:
		if err := s.driver.client.UpdateDatabaseDdl(ctx, s.driver.databasePath(""), []string{statement.Text}); err != nil {
			return nil, err
		}
		return &driver.Result{}, nil
	default:
		rowCount, err := s.driver.client.Update(ctx, statement.Text)
		if err != nil {
			return nil, err
		}
		return &driver.Result{AffectedRows: &rowCount}, nil
	}
}

func (s *session) queryStatement(ctx context.Context, statement string, opts driver.QueryOptions) (*driver.Result, error) {
	iter, err := s.driver.client.Read(ctx, statement)
	if err != nil {
		return nil, err
	}
	defer iter.Stop()

	cells, err := iter.Next()
	if errors.Is(err, io.EOF) {
		// An empty result still states its columns.
		columns, columnErr := resultColumns(iter)
		if columnErr != nil {
			return nil, columnErr
		}
		return &driver.Result{Columns: columns}, nil
	}
	if err != nil {
		return nil, err
	}

	columns, err := resultColumns(iter)
	if err != nil {
		return nil, err
	}
	result := &driver.Result{Columns: columns}
	totalBytes := driver.HeaderByteSize(columns)
	for {
		row := readRow(cells)
		result.Rows = append(result.Rows, row)
		totalBytes += driver.RowByteSize(row)
		n := len(result.Rows)
		if opts.MaxResultBytes > 0 && totalBytes > opts.MaxResultBytes {
			result.Truncated = true
			break
		}
		if opts.Limit > 0 && n >= opts.Limit {
			break
		}

		cells, err = iter.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

// explainStatement returns the plan Spanner would run for each statement of
// script, as a JSON document in the result's one row, which is what Bytebase
// showed.
func (s *session) explainStatement(ctx context.Context, script string, _ driver.QueryOptions) ([]*driver.Result, error) {
	var results []*driver.Result
	for _, statement := range statements(script) {
		startTime := time.Now()
		result, err := func() (*driver.Result, error) {
			// Spanner refuses DML in a single-use read-only transaction, so
			// anything but a query is analysed in a read-write transaction that is
			// rolled back.
			plan, err := s.driver.client.AnalyzeQuery(ctx, statement.Text, statement.Kind != query)
			if err != nil {
				return nil, err
			}
			planJSON, err := convertQueryPlanToJSON(plan)
			if err != nil {
				return nil, fmt.Errorf("failed to convert query plan to JSON: %w", err)
			}
			return &driver.Result{
				Columns: []driver.Column{{Name: "QUERY PLAN", TypeName: "JSON"}},
				Rows: [][]driver.Value{
					{{Kind: driver.String, String: planJSON}},
				},
				Plan: &driver.Plan{Format: "JSON"},
			}, nil
		}()
		stop := false
		if err != nil {
			result = &driver.Result{Err: err}
			stop = true
		}
		result.Statement = statement.Text
		result.Latency = time.Since(startTime)
		results = append(results, result)
		if stop {
			// A statement that failed ends the script, as it does for a query:
			// asking for the plan of a statement after the failure is asking the
			// server what the script never ran.
			break
		}
	}
	return results, nil
}

// resultColumns names the columns a result declares. A Spanner STRUCT has no
// column type name to report, and a result that declares one is refused, which
// is what Bytebase did.
func resultColumns(iter RowIterator) ([]driver.Column, error) {
	declared := iter.Columns()
	columns := make([]driver.Column, 0, len(declared))
	for _, column := range declared {
		typeName, err := getColumnTypeName(column.Type)
		if err != nil {
			return nil, err
		}
		columns = append(columns, driver.Column{Name: column.Name, TypeName: typeName})
	}
	return columns, nil
}

func getColumnTypeName(columnType *sppb.Type) (string, error) {
	if columnType.GetCode() == sppb.TypeCode_STRUCT {
		return "", errors.New("spanner STRUCT type is not supported")
	}
	if columnType.GetCode() == sppb.TypeCode_ARRAY {
		if columnType.GetArrayElementType().GetCode() == sppb.TypeCode_STRUCT {
			return "", errors.New("spanner STRUCT type is not supported")
		}
		return fmt.Sprintf("[]%s", columnType.GetArrayElementType().GetCode().String()), nil
	}
	return columnType.GetCode().String(), nil
}

func readRow(cells []Cell) []driver.Value {
	row := make([]driver.Value, 0, len(cells))
	for _, cell := range cells {
		row = append(row, convertSpannerValue(cell.Type, cell.Value))
	}
	return row
}

// convertSpannerValue converts a value the client decoded into the contract's
// cell, with the type metadata preserving INT64 precision and every Spanner
// type read as Bytebase read it.
func convertSpannerValue(columnType *sppb.Type, v *structpb.Value) driver.Value {
	if v == nil || v.Kind == nil {
		return driver.Value{Kind: driver.Null}
	}

	switch v.Kind.(type) {
	case *structpb.Value_NullValue:
		return driver.Value{Kind: driver.Null}
	case *structpb.Value_StringValue:
		stringValue := v.GetStringValue()
		if columnType == nil {
			return driver.Value{Kind: driver.String, String: stringValue}
		}

		switch columnType.GetCode() {
		case sppb.TypeCode_INT64:
			// Spanner encodes INT64 as a string to preserve precision. A value
			// that does not parse stays the text it came as; Bytebase logged the
			// failure, and the driver contract has no logger.
			val, err := strconv.ParseInt(stringValue, 10, 64)
			if err != nil {
				return driver.Value{Kind: driver.String, String: stringValue}
			}
			return driver.Value{Kind: driver.Int64, Int64: val}
		case sppb.TypeCode_TIMESTAMP:
			// Spanner encodes TIMESTAMP as an RFC 3339 string with nanosecond
			// precision.
			t, err := time.Parse(time.RFC3339Nano, stringValue)
			if err != nil {
				return driver.Value{Kind: driver.String, String: stringValue}
			}
			return driver.Value{Kind: driver.Timestamp, Time: t, Precision: timestampAccuracy(stringValue)}
		case sppb.TypeCode_BYTES:
			// Spanner encodes BYTES as a base64 string.
			bytes, err := base64.StdEncoding.DecodeString(stringValue)
			if err != nil {
				return driver.Value{Kind: driver.String, String: stringValue}
			}
			return driver.Value{Kind: driver.Bytes, Bytes: bytes}
		default:
			// DATE, JSON, NUMERIC and STRING all stay as their text.
			return driver.Value{Kind: driver.String, String: stringValue}
		}
	case *structpb.Value_NumberValue:
		// An INT64 the client read as a number keeps its precision only where a
		// float64 can hold it; Spanner writes INT64 as a string, so this is the
		// fallback for a value that came as a number.
		if columnType.GetCode() == sppb.TypeCode_INT64 {
			return driver.Value{Kind: driver.Int64, Int64: int64(v.GetNumberValue())}
		}
		return driver.Value{Kind: driver.Float64, Float64: v.GetNumberValue()}
	case *structpb.Value_BoolValue:
		return driver.Value{Kind: driver.Bool, Bool: v.GetBoolValue()}
	case *structpb.Value_ListValue, *structpb.Value_StructValue:
		// Flatten complex types to JSON strings: the contract has no array or
		// struct cell kind, and Bytebase showed them as their JSON document.
		goValue := v.AsInterface()
		jsonBytes, err := json.Marshal(goValue)
		if err != nil {
			return driver.Value{Kind: driver.String, String: fmt.Sprintf("%v", goValue)}
		}
		return driver.Value{Kind: driver.String, String: string(jsonBytes)}
	default:
		return driver.Value{Kind: driver.Null}
	}
}

// timestampAccuracy is the fractional-second digits the server wrote: a
// TIMESTAMP with none keeps Bytebase's default of microseconds, and one with
// more than nanoseconds of them is capped at nine.
func timestampAccuracy(timestamp string) int32 {
	accuracy := int32(6)
	dotIndex := strings.Index(timestamp, ".")
	if dotIndex < 0 {
		return accuracy
	}
	// The fraction ends at the zone, which is a Z or an offset.
	endIndex := strings.IndexAny(timestamp[dotIndex:], "Z+-")
	if endIndex > 0 {
		accuracy = int32(endIndex - 1)
		if accuracy > 9 {
			accuracy = 9
		}
	}
	return accuracy
}

// PlanNode represents a node in the Spanner query plan for JSON serialization.
type PlanNode struct {
	Index               int32                `json:"index"`
	Kind                string               `json:"kind"`
	DisplayName         string               `json:"displayName"`
	ChildLinks          []ChildLink          `json:"childLinks,omitempty"`
	ShortRepresentation *ShortRepresentation `json:"shortRepresentation,omitempty"`
	Metadata            map[string]any       `json:"metadata,omitempty"`
	ExecutionStats      map[string]any       `json:"executionStats,omitempty"`
}

// ChildLink represents a child link in the query plan.
type ChildLink struct {
	ChildIndex int32  `json:"childIndex"`
	Type       string `json:"type,omitempty"`
	Variable   string `json:"variable,omitempty"`
}

// ShortRepresentation represents the short representation of a scalar node.
type ShortRepresentation struct {
	Description string           `json:"description"`
	Subqueries  map[string]int32 `json:"subqueries,omitempty"`
}

// QueryPlan represents the full query plan for JSON serialization.
type QueryPlan struct {
	PlanNodes []PlanNode `json:"planNodes"`
}

// convertQueryPlanToJSON converts a Spanner QueryPlan to JSON string.
func convertQueryPlanToJSON(plan *sppb.QueryPlan) (string, error) {
	if plan == nil {
		return "{}", nil
	}

	queryPlan := QueryPlan{
		PlanNodes: make([]PlanNode, 0, len(plan.PlanNodes)),
	}

	for _, node := range plan.PlanNodes {
		planNode := PlanNode{
			Index:       node.Index,
			Kind:        node.Kind.String(),
			DisplayName: node.DisplayName,
		}

		if len(node.ChildLinks) > 0 {
			planNode.ChildLinks = make([]ChildLink, 0, len(node.ChildLinks))
			for _, link := range node.ChildLinks {
				planNode.ChildLinks = append(planNode.ChildLinks, ChildLink{
					ChildIndex: link.ChildIndex,
					Type:       link.Type,
					Variable:   link.Variable,
				})
			}
		}

		if node.ShortRepresentation != nil {
			planNode.ShortRepresentation = &ShortRepresentation{
				Description: node.ShortRepresentation.Description,
				Subqueries:  node.ShortRepresentation.Subqueries,
			}
		}

		if node.Metadata != nil {
			planNode.Metadata = node.Metadata.AsMap()
		}

		if node.ExecutionStats != nil {
			planNode.ExecutionStats = node.ExecutionStats.AsMap()
		}

		queryPlan.PlanNodes = append(queryPlan.PlanNodes, planNode)
	}

	jsonBytes, err := json.Marshal(queryPlan)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}
