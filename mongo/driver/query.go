package driver

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bytebase/gomongo"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/bytebase/omni/driver"
	"github.com/bytebase/omni/mongo"
	mongoanalysis "github.com/bytebase/omni/mongo/analysis"
)

// Session opens one session. MongoDB has no per-session state this driver
// sets, so a session is the driver itself.
func (d *Driver) Session(_ context.Context) (driver.Session, error) {
	return &session{driver: d}, nil
}

type session struct {
	driver *Driver
}

func (*session) Close() error { return nil }

// Query runs each statement of script in order and returns one Result per
// statement run. Errors MongoDB's client returned stay the client's types (a
// mongo.CommandError, a gomongo parse error).
//
// A statement that fails ends the script with its error on its own Result.
// Only text the mongo parser cannot split at all is an error for the call,
// since then there is no statement to run.
func (s *session) Query(ctx context.Context, script string, opts driver.QueryOptions) ([]*driver.Result, error) {
	statements, err := mongo.Parse(script)
	if err != nil {
		return []*driver.Result{{Statement: trimStatement(script), Err: err}}, nil
	}
	var results []*driver.Result
	for _, statement := range statements {
		text := trimStatement(statement.Text)
		startTime := time.Now()
		var result *driver.Result
		if opts.Explain {
			result, err = s.explain(ctx, text)
		} else {
			result, err = s.run(ctx, text, rowLimit(statement, opts), opts.MaxResultBytes)
			if err == nil {
				result.Plan = describePlan(text)
			}
		}
		if err != nil {
			result.Err = err
		}
		result.Latency = time.Since(startTime)
		results = append(results, result)
		if err != nil {
			break
		}
	}
	return results, nil
}

// rowLimit is the rows QueryOptions.Limit caps a statement to. A count gets
// none: the client applies the cap to the count itself, so a limit meant for
// the rows an editor shows would change what the count counted.
func rowLimit(statement mongo.Statement, opts driver.QueryOptions) int {
	if opts.Limit <= 0 {
		return 0
	}
	if analyzed := mongoanalysis.Analyze(statement.AST); analyzed != nil && analyzed.Operation == mongoanalysis.OpCount {
		return 0
	}
	return opts.Limit
}

// run runs one statement, capping its rows at limit when that is positive and
// at maxResultBytes when that is.
func (s *session) run(ctx context.Context, statement string, limit int, maxResultBytes int64) (*driver.Result, error) {
	result, err := s.driver.client.Execute(ctx, s.driver.databaseName, statement, int64(limit))
	if err != nil {
		return &driver.Result{Statement: statement}, err
	}
	return convertResult(result, statement, maxResultBytes), nil
}

// explain plans one statement and reports the result as a plan. A plan of a
// statement the caller did not write answers the wrong question, so no row
// limit applies: limiting the query a plan covers would change the plan rather
// than the reading of it.
func (s *session) explain(ctx context.Context, statement string) (*driver.Result, error) {
	explained, err := explainStatement(statement, "")
	if err != nil {
		return &driver.Result{Statement: statement}, err
	}
	result, err := s.run(ctx, explained, 0, 0)
	if err != nil {
		return result, err
	}
	result.Plan = describePlan(explained)
	return result, nil
}

// describePlan reports the plan a statement hands back when it ran as a query,
// which is only when the statement asks for its own plan.
func describePlan(statement string) *driver.Plan {
	format, executed, ok := describeExplain(statement)
	if !ok {
		return nil
	}
	return &driver.Plan{Format: strings.ToUpper(format), Executed: executed}
}

// convertResult renders one mongosh execution as a result: one TEXT column
// whose rows are the values as relaxed Extended JSON, which is what the editor
// shows. maxResultBytes caps the rows read and sets Result.Truncated when that
// is positive.
func convertResult(res *gomongo.Result, statement string, maxResultBytes int64) *driver.Result {
	result := &driver.Result{
		Statement: statement,
		Columns:   []driver.Column{{Name: "result", TypeName: "TEXT"}},
	}
	totalBytes := driver.HeaderByteSize(result.Columns)
	for _, v := range res.Value {
		str, err := marshalValueToExtJSON(v)
		if err != nil {
			// Bytebase logged this and showed the value's Go spelling; the
			// driver contract has no logger, so the reading stays.
			str = fmt.Sprintf("%v", v)
		}
		row := []driver.Value{{Kind: driver.String, String: str}}
		result.Rows = append(result.Rows, row)
		if maxResultBytes > 0 {
			totalBytes += driver.RowByteSize(row)
			if totalBytes > maxResultBytes {
				result.Truncated = true
				break
			}
		}
	}
	return result
}

// trimStatement trims a segment the way Bytebase's query path trimmed the
// statement it sent, so a caller reads the text without its padding.
func trimStatement(statement string) string {
	return strings.Trim(statement, " \t\n\r\f;")
}

// marshalValueToExtJSON marshals a value to Extended JSON (relaxed) format.
// The frontend expects each row to be a JSON object with key-value pairs, so a
// value that is not a document is wrapped in an object with a "value" key:
// Extended JSON has no top-level scalar, and a distinct value would otherwise
// be rendered as its Go spelling.
func marshalValueToExtJSON(v any) (string, error) {
	toMarshal := v
	if !isDocument(v) {
		toMarshal = bson.M{"value": v}
	}
	jsonBytes, err := bson.MarshalExtJSON(toMarshal, false, false)
	if err != nil {
		return "", err
	}
	return string(normalizeExtJSONNumbers(jsonBytes)), nil
}

// isDocument reports whether v renders as a document or an array, the values a
// row shows as themselves.
func isDocument(v any) bool {
	switch v.(type) {
	case bson.D, bson.M, bson.A, []any, map[string]any:
		return true
	}
	return false
}
