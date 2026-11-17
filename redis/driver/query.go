package driver

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/bytebase/omni/driver"
)

// Session borrows one connection, so a command that sets connection state (SELECT,
// CLIENT SETNAME) changes it for the session's next command. A cluster has no
// connection to lend, because it routes each command by its key.
func (d *Driver) Session(ctx context.Context) (driver.Session, error) {
	number, err := d.client.Session(ctx)
	if err != nil {
		return nil, err
	}
	return &Session{driver: d, session: number}, nil
}

// Session runs a Redis script's commands on one connection. It implements
// driver.Session, and driver.AllNodesQuerier for the all-nodes choice Bytebase's
// query API offers per request, which the contract's QueryOptions has no field
// for.
type Session struct {
	driver  *Driver
	session int
	// mu serializes the session's commands: the connection it holds is not safe
	// for concurrent use.
	mu sync.Mutex
}

var _ driver.AllNodesQuerier = (*Session)(nil)

// Close returns the session's connection to the pool.
func (s *Session) Close() error {
	return s.driver.client.CloseSession(context.Background(), s.session)
}

// Query runs each command of script in order, one command per line, and returns
// one Result per command. Each command is its own request, so a command the server
// refuses ends the script and the commands after it do not run; the error
// returned is nil, since the script ran as far as it could. Only text whose
// quoting cannot be read at all is an error for the call, since then there is no
// command to send.
func (s *Session) Query(ctx context.Context, script string, opts driver.QueryOptions) ([]*driver.Result, error) {
	if opts.Explain {
		return nil, errors.New("Redis does not support EXPLAIN")
	}
	commands, err := splitCommands(script)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	results := make([]*driver.Result, 0, len(commands))
	for _, command := range commands {
		startTime := time.Now()
		cmds, err := s.driver.client.Do(ctx, s.session, [][]any{command.args})
		if err != nil {
			// The commands before this one ran, so their results stand beside the
			// error that kept the rest from being sent.
			return results, err
		}
		result := newResult(command.text)
		result.Latency = time.Since(startTime)
		if len(cmds) == 0 {
			results = append(results, result)
			continue
		}
		if cmds[0].Err() != nil && cmds[0].Err() != redis.Nil {
			result.Err = cmds[0].Err()
			results = append(results, result)
			break
		}
		appendReplyRows(result, cmds[0], opts.MaxResultBytes, driver.HeaderByteSize(result.Columns))
		results = append(results, result)
	}
	return results, nil
}

// QueryAllNodes runs each command of script on every node of the cluster and
// returns one Result per command, whose rows are every node's answers in node
// address order. A command that fails on any node ends the script on its Result,
// as it does in Query.
func (s *Session) QueryAllNodes(ctx context.Context, script string, opts driver.QueryOptions) ([]*driver.Result, error) {
	if opts.Explain {
		return nil, errors.New("Redis does not support EXPLAIN")
	}
	commands, err := splitCommands(script)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	results := make([]*driver.Result, 0, len(commands))
	for _, command := range commands {
		startTime := time.Now()
		shards, err := s.driver.client.DoEachShard(ctx, [][]any{command.args})
		if err != nil {
			// The commands before this one ran, so their results stand beside the
			// error that kept the rest from being sent. A client of one server
			// says so itself, without the wrap a failure of the cluster's own
			// state takes.
			var notCluster *notClusterError
			if errors.As(err, &notCluster) {
				return results, err
			}
			return results, fmt.Errorf("failed to query for each shard: %w", err)
		}
		// One command's rows are every node's answers, in node order, before the
		// next command is sent, so a command that fails on one node leaves the
		// earlier commands whole and the commands after it unsent.
		result := newResult(command.text)
		size := driver.HeaderByteSize(result.Columns)
		for _, shard := range shards {
			if len(shard) == 0 {
				continue
			}
			if shard[0].Err() != nil && shard[0].Err() != redis.Nil {
				result.Err = shard[0].Err()
				break
			}
			size = appendReplyRows(result, shard[0], opts.MaxResultBytes, size)
		}
		result.Latency = time.Since(startTime)
		results = append(results, result)
		if result.Err != nil {
			break
		}
	}
	return results, nil
}

// newResult is the empty result of one command: the two columns Bytebase showed
// a Redis command's answer in.
func newResult(statement string) *driver.Result {
	return &driver.Result{
		Statement: statement,
		Columns: []driver.Column{
			{Name: "#", TypeName: "INT"},
			{Name: "Value", TypeName: "TEXT"},
		},
	}
}

// appendReplyRows appends the rows one command's answer holds: an array becomes
// one row per element, a missing key one row holding NULL, and any other answer
// one row. It stops at maxBytes, which is positive when the caller capped the
// result, and marks the result truncated. It takes the result's byte count so far
// and returns the updated one, because an all-nodes result collects rows from
// every node.
func appendReplyRows(result *driver.Result, cmd *redis.Cmd, maxBytes int64, totalBytes int64) int64 {
	if result.Truncated {
		return totalBytes
	}
	add := func(row []driver.Value) {
		result.Rows = append(result.Rows, row)
		totalBytes += driver.RowByteSize(row)
		if maxBytes > 0 && totalBytes > maxBytes {
			result.Truncated = true
		}
	}
	if cmd.Err() == redis.Nil {
		add([]driver.Value{{Kind: driver.Int32, Int64: 1}, {Kind: driver.Null}})
		return totalBytes
	}
	if list, ok := cmd.Val().([]any); ok {
		for i, value := range list {
			if result.Truncated {
				return totalBytes
			}
			add(replyRow(i+1, value))
		}
		return totalBytes
	}
	add(replyRow(1, cmd.Val()))
	return totalBytes
}

// replyRow is one row of an answer: the element's position and its value as a
// string, or NULL for an element that holds none — a missing key in an MGET, a
// missing field in an HMGET.
func replyRow(i int, value any) []driver.Value {
	if value == nil {
		return []driver.Value{
			{Kind: driver.Int32, Int64: int64(i)},
			{Kind: driver.Null},
		}
	}
	return []driver.Value{
		{Kind: driver.Int32, Int64: int64(i)},
		{Kind: driver.String, String: fmt.Sprintf("%v", value)},
	}
}
