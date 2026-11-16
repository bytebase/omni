package driver

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/driver"
)

// TestPing asks the endpoint for its root document.
func TestPing(t *testing.T) {
	d := openDriver(t)
	require.NoError(t, d.Ping(context.Background()))
}

// TestQueryReadsBothAnswerShapes reads a search, whose answer is an object, and
// a cat request, whose answer is an array, in one script.
func TestQueryReadsBothAnswerShapes(t *testing.T) {
	d := openDriver(t)
	executeCorpus(t, d, mustRead(t, filepath.Join("testdata", "corpus.http")))
	sess := openSession(t, d)

	results, err := sess.Query(context.Background(), `
GET /alpha/_search
{
  "query": {"match_all": {}}
}

GET /_cat/indices?format=json`, driver.QueryOptions{})
	require.NoError(t, err)
	require.Len(t, results, 2)
	require.NoError(t, results[0].Err)
	require.NoError(t, results[1].Err)

	// The search's answer is one row of one column per key, the keys in name
	// order so one query shows its columns the same way every run.
	search := results[0]
	var names []string
	for _, column := range search.Columns {
		names = append(names, column.Name)
	}
	require.Equal(t, []string{"_shards", "hits", "timed_out", "took"}, names)
	values := map[string]string{}
	for i, column := range search.Columns {
		require.Empty(t, column.TypeName, "Elasticsearch states no column type")
		values[column.Name] = search.Rows[0][i].String
	}
	require.Equal(t, "0", values["took"])
	require.Equal(t, "false", values["timed_out"])
	require.Contains(t, values["hits"], `"total":{"relation":"eq","value":5}`)
	require.Len(t, search.Rows, 1)

	// The cat answer is an array: one column holding it.
	cat := results[1]
	require.Equal(t, []driver.Column{{Name: "result"}}, cat.Columns)
	require.Len(t, cat.Rows, 1)
	require.Contains(t, cat.Rows[0][0].String, `"index":"alpha"`)
	require.Contains(t, cat.Rows[0][0].String, `"docs.count":"5"`)
}

// TestQueryTextPlain reads an answer the endpoint sends as plain text.
func TestQueryTextPlain(t *testing.T) {
	d := openDriver(t)
	executeCorpus(t, d, mustRead(t, filepath.Join("testdata", "corpus.http")))
	sess := openSession(t, d)

	// The cat API's plain-text rendering names the shards the indices run on,
	// which a container makes up afresh; asking for one column keeps the answer
	// one a tape can hold.
	results, err := sess.Query(context.Background(), "GET /_cat/indices?h=index", driver.QueryOptions{})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.NoError(t, results[0].Err)
	require.Equal(t, []driver.Column{{Name: "text/plain"}}, results[0].Columns)
	require.Len(t, results[0].Rows, 1)
	require.Contains(t, results[0].Rows[0][0].String, "alpha")
	require.Contains(t, results[0].Rows[0][0].String, "beta")
}

// TestQueryExplain holds an explain to a search that ran with profiling on,
// whose hits do not come back and whose result says it is a plan.
func TestQueryExplain(t *testing.T) {
	d := openDriver(t)
	executeCorpus(t, d, mustRead(t, filepath.Join("testdata", "corpus.http")))
	sess := openSession(t, d)

	results, err := sess.Query(context.Background(), `
GET /alpha/_search
{"query": {"match_all": {}}}`, driver.QueryOptions{Explain: true})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.NoError(t, results[0].Err)
	require.Equal(t, []driver.Column{{Name: "profile"}}, results[0].Columns)
	require.Len(t, results[0].Rows, 1)
	require.Contains(t, results[0].Rows[0][0].String, `"profile"`)
	require.Contains(t, results[0].Rows[0][0].String, `"took"`)
	require.NotContains(t, results[0].Rows[0][0].String, `"hits"`)
	require.Equal(t, &driver.Plan{Format: "JSON", Executed: true}, results[0].Plan)
	require.Contains(t, results[0].Statement, `"profile":true`)
}

// TestQueryExplainRefusesAWrite holds an explain of anything but a search to
// the request's own refusal, so an explain can never carry out a write.
func TestQueryExplainRefusesAWrite(t *testing.T) {
	d := openDriver(t)
	sess := openSession(t, d)

	for _, script := range []string{
		`POST /alpha/_doc/1
{"name": "row 1"}`,
		`DELETE /alpha`,
		`POST /alpha/_bulk
{"index": {}}
{"name": "row 1"}`,
	} {
		results, err := sess.Query(context.Background(), script, driver.QueryOptions{Explain: true})
		require.NoError(t, err)
		require.Len(t, results, 1)
		require.Error(t, results[0].Err)
		require.Contains(t, results[0].Err.Error(), "is not a search to explain")
	}
}

// TestQueryError holds a failing request to its own Result, and the script
// stops there.
func TestQueryError(t *testing.T) {
	d := openDriver(t)
	sess := openSession(t, d)

	results, err := sess.Query(context.Background(), `
GET /missing-index/_search
{"query": {"match_all": {}}}

GET /`, driver.QueryOptions{})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.Error(t, results[0].Err)
	require.Contains(t, results[0].Err.Error(), "request failed with status code 404")
}

// TestQueryUnreadableScript holds a script the REST console parser cannot read
// to the call's error, since then there is no request to send.
func TestQueryUnreadableScript(t *testing.T) {
	d := openDriver(t)
	sess := openSession(t, d)
	_, err := sess.Query(context.Background(), "GET /alpha\n{not json", driver.QueryOptions{})
	require.Error(t, err)
}

// TestQueryHead holds an existence check to its status: a HEAD answer carries
// no body, and a successful one is a result with no columns and no rows.
func TestQueryHead(t *testing.T) {
	d := openDriver(t)
	executeCorpus(t, d, mustRead(t, filepath.Join("testdata", "corpus.http")))
	sess := openSession(t, d)

	results, err := sess.Query(context.Background(), "HEAD /alpha", driver.QueryOptions{})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.NoError(t, results[0].Err)
	require.Empty(t, results[0].Columns)
	require.Empty(t, results[0].Rows)

	// A HEAD that fails still states its error in the status.
	results, err = sess.Query(context.Background(), "HEAD /missing", driver.QueryOptions{})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.Error(t, results[0].Err)
	require.Contains(t, results[0].Err.Error(), "request failed with status code 404")
}

// TestQueryReadsADocumentThatLooksLikeABulkAnswer holds a document whose own
// fields are named like a bulk answer's to the document: only a bulk request's
// answer states a bulk verdict, and a document read under a path holding the
// word is still a document.
func TestQueryReadsADocumentThatLooksLikeABulkAnswer(t *testing.T) {
	d := openDriver(t)
	executeCorpus(t, d, mustRead(t, filepath.Join("testdata", "corpus.http")))
	sess := openSession(t, d)

	results, err := sess.Query(context.Background(), "GET /alpha/_source/app_bulk", driver.QueryOptions{})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.NoError(t, results[0].Err)
	require.Len(t, results[0].Rows, 1)
	values := map[string]string{}
	for i, column := range results[0].Columns {
		values[column.Name] = results[0].Rows[0][i].String
	}
	require.Equal(t, "true", values["errors"])
	require.Contains(t, values["items"], "logged application error")
}

// TestQueryBulkFailure holds a bulk whose items failed to the answer's own
// verdict: the request answers 200 either way, and the failures are in its body.
func TestQueryBulkFailure(t *testing.T) {
	d := openDriver(t)
	executeCorpus(t, d, mustRead(t, filepath.Join("testdata", "corpus.http")))
	sess := openSession(t, d)

	results, err := sess.Query(context.Background(), `
POST /_bulk
{"update":{"_index":"alpha","_id":"missing"}}
{"doc":{"n":1}}
`, driver.QueryOptions{})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.Error(t, results[0].Err)
	require.Contains(t, results[0].Err.Error(), "1 of 1 bulk items failed")
	require.Contains(t, results[0].Err.Error(), "document missing")
}

// TestQueryKeepsLargeNumbers holds a long past 2^53 to its digits: read as a
// float64 it would come back rounded.
func TestQueryKeepsLargeNumbers(t *testing.T) {
	d := openDriver(t)
	executeCorpus(t, d, mustRead(t, filepath.Join("testdata", "corpus.http")))
	sess := openSession(t, d)

	results, err := sess.Query(context.Background(), "GET /alpha/_doc/big", driver.QueryOptions{})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.NoError(t, results[0].Err)
	values := map[string]string{}
	for i, column := range results[0].Columns {
		values[column.Name] = results[0].Rows[0][i].String
	}
	require.Contains(t, values["_source"], `"big":9007199254740993`)
}

type stepLogger struct {
	steps []string
}

func (l *stepLogger) CommandExecute(cmd driver.Command) {
	l.steps = append(l.steps, fmt.Sprintf("exec[%d:%d]%s", cmd.ByteStart, cmd.ByteEnd, cmd.Text))
}

func (l *stepLogger) CommandResponse(affected int64, perResult []int64, err error) {
	if err != nil {
		l.steps = append(l.steps, "resp:err")
		return
	}
	l.steps = append(l.steps, fmt.Sprintf("resp:%d", affected))
}

func (*stepLogger) TransactionControl(driver.TransactionOp, error) {}

func (*stepLogger) Retry(error, int, int) {}

// TestExecuteLogsItsSteps holds Execute to the contract's logger: one command
// with its range and one response for each request it sends.
func TestExecuteLogsItsSteps(t *testing.T) {
	d := openDriver(t)
	log := &stepLogger{}
	_, err := d.Execute(context.Background(), "PUT /logged\n{\"settings\": {\"number_of_shards\": 1}}\n\nGET /missing-index/_search", driver.ExecuteOptions{Log: log})
	require.Error(t, err)
	require.Len(t, log.steps, 4)
	require.Contains(t, log.steps[0], "PUT /logged")
	require.Contains(t, log.steps[1], "resp:0")
	require.Contains(t, log.steps[2], "GET /missing-index/_search")
	require.Equal(t, "resp:err", log.steps[3])
}

func openSession(t *testing.T, d *Driver) driver.Session {
	t.Helper()
	session, err := d.Session(context.Background())
	require.NoError(t, err)
	t.Cleanup(func() { _ = session.Close() })
	return session
}
