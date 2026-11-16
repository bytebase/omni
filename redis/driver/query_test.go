package driver

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/driver"
	"github.com/bytebase/omni/redis/internal/testcontainer"
)

// query runs script on d's session and returns the results.
func query(t *testing.T, d *Driver, script string, opts driver.QueryOptions) []*driver.Result {
	t.Helper()
	session, err := d.Session(context.Background())
	require.NoError(t, err)
	defer func() { require.NoError(t, session.Close()) }()
	results, err := session.Query(context.Background(), script, opts)
	require.NoError(t, err)
	return results
}

// queryAllNodes is query through Session's all-nodes path.
func queryAllNodes(t *testing.T, d *Driver, script string, opts driver.QueryOptions) []*driver.Result {
	t.Helper()
	session, err := d.Session(context.Background())
	require.NoError(t, err)
	defer func() { require.NoError(t, session.Close()) }()
	results, err := session.(*Session).QueryAllNodes(context.Background(), script, opts)
	require.NoError(t, err)
	return results
}

// columns are the two columns every Redis command's answer is reported in.
var columns = []driver.Column{{Name: "#", TypeName: "INT"}, {Name: "Value", TypeName: "TEXT"}}

// row reads one row back as its position and its value.
func row(t *testing.T, r []driver.Value) (int64, string) {
	t.Helper()
	require.Len(t, r, 2)
	require.Equal(t, driver.Int32, r[0].Kind)
	if r[1].Kind == driver.Null {
		return r[0].Int64, "<nil>"
	}
	require.Equal(t, driver.String, r[1].Kind)
	return r[0].Int64, r[1].String
}

// TestQuery runs a script of commands, one per line, with a quoted argument, a
// missing key and a blank line, and holds each command's answer to a Result.
func TestQuery(t *testing.T) {
	ctx := context.Background()
	d := openDriver(t, testcontainer.Login{}, "0")
	_, err := d.Execute(ctx, "SET greeting \"hello world\"\nRPUSH queue first second", driver.ExecuteOptions{})
	require.NoError(t, err)

	results := query(t, d, "GET greeting\n\nGET missing\nLRANGE queue 0 -1\nDBSIZE", driver.QueryOptions{})
	require.Len(t, results, 4)
	for i, statement := range []string{"GET greeting", "GET missing", "LRANGE queue 0 -1", "DBSIZE"} {
		require.Equal(t, statement, results[i].Statement)
		require.Equal(t, columns, results[i].Columns)
		require.NoError(t, results[i].Err)
		require.Positive(t, results[i].Latency, "a result reports how long its commands took")
		require.LessOrEqual(t, results[i].Latency, time.Minute)
	}

	// A quoted argument is one value, so the answer is the string it holds.
	require.Len(t, results[0].Rows, 1)
	position, value := row(t, results[0].Rows[0])
	require.Equal(t, int64(1), position)
	require.Equal(t, "hello world", value)

	// A missing key is a row holding NULL, not an error.
	require.Len(t, results[1].Rows, 1)
	_, value = row(t, results[1].Rows[0])
	require.Equal(t, "<nil>", value)
	require.Equal(t, driver.Null, results[1].Rows[0][1].Kind)

	// An array answer is one row per element.
	require.Len(t, results[2].Rows, 2)
	for i, want := range []string{"first", "second"} {
		position, value := row(t, results[2].Rows[i])
		require.Equal(t, int64(i+1), position)
		require.Equal(t, want, value)
	}

	// Any other answer is one row.
	require.Len(t, results[3].Rows, 1)
	_, value = row(t, results[3].Rows[0])
	require.Equal(t, "2", value)
}

// TestQueryCommandError covers a command the server refuses: its answer is the
// Result's error, and the commands after it are not reported.
func TestQueryCommandError(t *testing.T) {
	ctx := context.Background()
	d := openDriver(t, testcontainer.Login{}, "0")
	_, err := d.Execute(ctx, "SET counter text", driver.ExecuteOptions{})
	require.NoError(t, err)

	results := query(t, d, "INCR counter\nDBSIZE", driver.QueryOptions{})
	require.Len(t, results, 1)
	require.Equal(t, "INCR counter", results[0].Statement)
	require.Equal(t, columns, results[0].Columns)
	require.ErrorContains(t, results[0].Err, "ERR value is not an integer or out of range")
	require.Empty(t, results[0].Rows)
}

// TestQueryExplainRefused holds an explain to the error Redis's driver answers
// it with: Redis states no plan.
func TestQueryExplainRefused(t *testing.T) {
	session, err := (&Driver{}).Session(context.Background())
	require.NoError(t, err)
	_, err = session.Query(context.Background(), "GET key", driver.QueryOptions{Explain: true})
	require.ErrorContains(t, err, "Redis does not support EXPLAIN")
}

// TestQueryMaxResultBytes holds a command's rows to the cap and the truncation
// flag it sets.
func TestQueryMaxResultBytes(t *testing.T) {
	ctx := context.Background()
	d := openDriver(t, testcontainer.Login{}, "0")
	_, err := d.Execute(ctx,
		"RPUSH list aaaaaaaaaa bbbbbbbbbb cccccccccc dddddddddd eeeeeeeeee ffffffffff gggggggggg hhhhhhhhhh", driver.ExecuteOptions{})
	require.NoError(t, err)

	results := query(t, d, "LRANGE list 0 -1", driver.QueryOptions{MaxResultBytes: 120})
	require.Len(t, results, 1)
	require.NoError(t, results[0].Err)
	require.True(t, results[0].Truncated)
	require.NotEmpty(t, results[0].Rows)
	require.Less(t, len(results[0].Rows), 8)
}

// TestQueryAllNodes runs a script on every node of the cluster. The pinned
// server is a cluster of one node, so each answer is that node's.
func TestQueryAllNodes(t *testing.T) {
	ctx := context.Background()
	d := openClusterDriver(t, "0")
	_, err := d.Execute(ctx, "SET allnodes:key allnodes-value", driver.ExecuteOptions{})
	require.NoError(t, err)

	results := queryAllNodes(t, d, "GET allnodes:key\nDBSIZE", driver.QueryOptions{})
	require.Len(t, results, 2)
	for i, want := range []string{"allnodes-value", "1"} {
		require.Equal(t, columns, results[i].Columns)
		require.NoError(t, results[i].Err)
		require.Len(t, results[i].Rows, 1, "statement %q returned no row", results[i].Statement)
		_, value := row(t, results[i].Rows[0])
		require.Equal(t, want, value)
	}
}
