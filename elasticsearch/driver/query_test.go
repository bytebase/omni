package driver

import (
	"context"
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

	// The search's answer is one row of one column per key. Its keys come out
	// in no fixed order, since the answer is a Go map.
	search := results[0]
	values := map[string]string{}
	for i, column := range search.Columns {
		require.Empty(t, column.TypeName, "Elasticsearch states no column type")
		values[column.Name] = search.Rows[0][i].String
	}
	require.Equal(t, "0", values["took"])
	require.Equal(t, "false", values["timed_out"])
	require.Contains(t, values["hits"], `"total":{"relation":"eq","value":3}`)
	require.Len(t, search.Rows, 1)

	// The cat answer is an array: one column holding it.
	cat := results[1]
	require.Equal(t, []driver.Column{{Name: "result"}}, cat.Columns)
	require.Len(t, cat.Rows, 1)
	require.Contains(t, cat.Rows[0][0].String, `"index":"alpha"`)
	require.Contains(t, cat.Rows[0][0].String, `"docs.count":"3"`)
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

func openSession(t *testing.T, d *Driver) driver.Session {
	t.Helper()
	session, err := d.Session(context.Background())
	require.NoError(t, err)
	t.Cleanup(func() { _ = session.Close() })
	return session
}
