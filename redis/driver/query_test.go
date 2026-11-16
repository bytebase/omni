package driver

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
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
// Result's error, the commands after it are not reported, and they do not run.
func TestQueryCommandError(t *testing.T) {
	ctx := context.Background()
	d := openDriver(t, testcontainer.Login{}, "0")
	_, err := d.Execute(ctx, "SET counter text", driver.ExecuteOptions{})
	require.NoError(t, err)

	results := query(t, d, "INCR counter\nSET after value", driver.QueryOptions{})
	require.Len(t, results, 1)
	require.Equal(t, "INCR counter", results[0].Statement)
	require.Equal(t, columns, results[0].Columns)
	require.ErrorContains(t, results[0].Err, "ERR value is not an integer or out of range")
	require.Empty(t, results[0].Rows)

	after := query(t, d, "GET after", driver.QueryOptions{})
	require.Len(t, after, 1)
	require.Len(t, after[0].Rows, 1)
	require.Equal(t, driver.Null, after[0].Rows[0][1].Kind, "the command after the refused one ran")
}

// TestQueryExplainRefused holds an explain to the error Redis's driver answers
// it with: Redis states no plan.
func TestQueryExplainRefused(t *testing.T) {
	d := openDriver(t, testcontainer.Login{}, "0")
	session, err := d.Session(context.Background())
	require.NoError(t, err)
	defer func() { require.NoError(t, session.Close()) }()
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

// shardClient answers the all-nodes path, which the pinned server cannot: it is
// one node. Each shard's answers are the Val each request's command answered with,
// in request order.
type shardClient struct {
	answers [][]any
	errs    [][]error
	calls   int
}

func (*shardClient) Ping(context.Context) error { return errors.New("not asked") }

func (*shardClient) Info(context.Context, string) (string, error) { return "", errors.New("not asked") }

func (*shardClient) ConfigGet(context.Context, string) (map[string]string, error) {
	return nil, errors.New("not asked")
}

func (*shardClient) Session(context.Context) (int, error) { return 0, nil }

func (*shardClient) Do(context.Context, int, [][]any) ([]*redis.Cmd, error) {
	return nil, errors.New("not asked")
}

func (*shardClient) CloseSession(context.Context, int) error { return nil }

func (*shardClient) Close() error { return nil }

func (c *shardClient) DoEachShard(ctx context.Context, commands [][]any) ([][]*redis.Cmd, error) {
	call := c.calls
	c.calls++
	shards := make([][]*redis.Cmd, len(c.answers))
	for s, answers := range c.answers {
		cmd := redis.NewCmd(ctx, commands[0]...)
		if call < len(answers) {
			cmd.SetVal(answers[call])
		}
		if s < len(c.errs) && call < len(c.errs[s]) {
			cmd.SetErr(c.errs[s][call])
		}
		shards[s] = []*redis.Cmd{cmd}
	}
	return shards, nil
}

// queryShards runs script through the all-nodes path on a fake cluster.
func queryShards(t *testing.T, client client, script string, opts driver.QueryOptions) []*driver.Result {
	t.Helper()
	d := &Driver{client: client}
	session, err := d.Session(context.Background())
	require.NoError(t, err)
	results, err := session.(*Session).QueryAllNodes(context.Background(), script, opts)
	require.NoError(t, err)
	return results
}

// TestQueryAllNodesAcrossShards holds a command's rows to every node's answers, in
// node order, which one node's client never shows.
func TestQueryAllNodesAcrossShards(t *testing.T) {
	results := queryShards(t, &shardClient{
		answers: [][]any{
			{"a", "c"},
			{"b", "d"},
		},
		errs: [][]error{{nil, nil}, {nil, nil}},
	}, "DBSIZE\nPING", driver.QueryOptions{})

	require.Len(t, results, 2)
	for _, result := range results {
		require.NoError(t, result.Err)
	}
	require.Len(t, results[0].Rows, 2)
	require.Len(t, results[1].Rows, 2)
	for i, want := range []string{"a", "b"} {
		_, value := row(t, results[0].Rows[i])
		require.Equal(t, want, value)
	}
}

// TestQueryAllNodesStopsAfterAFailingCommand holds a command one node refused: the
// commands before it keep every node's answers, and the commands after it are not
// reported.
func TestQueryAllNodesStopsAfterAFailingCommand(t *testing.T) {
	refused := errors.New("ERR unknown command")
	results := queryShards(t, &shardClient{
		answers: [][]any{
			{"a"},
			{"b"},
		},
		errs: [][]error{{nil, refused}, {nil, refused}},
	}, "DBSIZE\nNOTACOMMAND\nDBSIZE", driver.QueryOptions{})

	require.Len(t, results, 2)
	require.NoError(t, results[0].Err)
	require.Len(t, results[0].Rows, 2, "the command before the failing one lost a node's answer")
	require.ErrorContains(t, results[1].Err, "unknown command")
}

// TestQueryAllNodesCountsEveryShardsRows holds the result-size cap across nodes:
// every node's rows count toward it, so two nodes cannot each pass the cap while
// their sum does not.
func TestQueryAllNodesCountsEveryShardsRows(t *testing.T) {
	items := []any{"aaaaaaaaaa", "bbbbbbbbbb", "cccccccccc", "dddddddddd"}
	results := queryShards(t, &shardClient{
		answers: [][]any{
			{items},
			{items},
		},
		errs: [][]error{{nil}, {nil}},
	}, "LRANGE list 0 -1", driver.QueryOptions{MaxResultBytes: 120})

	require.Len(t, results, 1)
	require.NoError(t, results[0].Err)
	require.True(t, results[0].Truncated)
	require.Len(t, results[0].Rows, 5, "the cap was counted again for the second node")
}

// scriptClient answers a script's commands from a list and fails to send the
// command numbered failAt, which is a call error rather than an answer.
type scriptClient struct {
	answers []any
	errs    []error
	sent    int
	failAt  int
}

func (*scriptClient) Ping(context.Context) error { return errors.New("not asked") }

func (*scriptClient) Info(context.Context, string) (string, error) {
	return "", errors.New("not asked")
}

func (*scriptClient) ConfigGet(context.Context, string) (map[string]string, error) {
	return nil, errors.New("not asked")
}

func (*scriptClient) Session(context.Context) (int, error) { return 1, nil }

func (c *scriptClient) Do(ctx context.Context, _ int, commands [][]any) ([]*redis.Cmd, error) {
	if c.sent == c.failAt {
		return nil, errors.New("connection reset by peer")
	}
	cmd := redis.NewCmd(ctx, commands[0]...)
	cmd.SetVal(c.answers[c.sent])
	cmd.SetErr(c.errs[c.sent])
	c.sent++
	return []*redis.Cmd{cmd}, nil
}

func (*scriptClient) CloseSession(context.Context, int) error { return nil }

func (*scriptClient) DoEachShard(context.Context, [][]any) ([][]*redis.Cmd, error) {
	return nil, errors.New("not asked")
}

func (*scriptClient) Close() error { return nil }

// TestQueryKeepsCompletedResultsWhenALaterCommandCannotBeSent holds a script whose
// connection broke midway: the results of the commands that ran come back beside
// the error that kept the rest from being sent.
func TestQueryKeepsCompletedResultsWhenALaterCommandCannotBeSent(t *testing.T) {
	d := &Driver{client: &scriptClient{answers: []any{"a"}, errs: []error{nil}, failAt: 1}}
	session, err := d.Session(context.Background())
	require.NoError(t, err)
	results, err := session.Query(context.Background(), "PING\nPING", driver.QueryOptions{})
	require.ErrorContains(t, err, "connection reset by peer")
	require.Len(t, results, 1)
	require.NoError(t, results[0].Err)
	require.Len(t, results[0].Rows, 1)
}
