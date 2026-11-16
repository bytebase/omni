package driver

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/shlex"
	"github.com/redis/go-redis/v9"

	"github.com/bytebase/omni/driver"
)

// Session returns one session. Redis has no session state this driver sets —
// the database index is chosen when the client is opened — so a session is the
// driver itself.
func (d *Driver) Session(context.Context) (driver.Session, error) {
	return &Session{driver: d}, nil
}

// Session runs a Redis script's commands. It implements driver.Session, and adds
// QueryAllNodes for the all-nodes choice Bytebase's query API offers per request,
// which the contract's QueryOptions has no field for.
type Session struct {
	driver *Driver
}

// Close closes the session, which holds nothing of its own.
func (*Session) Close() error { return nil }

// Query runs each command of script in order, one command per line, and returns
// one Result per command. A command that fails ends the script on its Result; the
// error returned is nil, since the script ran as far as it could. Only text whose
// quoting cannot be read at all is an error for the call, since then there is no
// command to send.
func (s *Session) Query(ctx context.Context, script string, opts driver.QueryOptions) ([]*driver.Result, error) {
	if opts.Explain {
		return nil, errors.New("Redis does not support EXPLAIN")
	}
	inputs, commands, err := splitCommands(script)
	if err != nil {
		return nil, err
	}

	startTime := time.Now()
	cmds, err := s.driver.client.Do(ctx, inputs)
	if err != nil {
		return nil, err
	}

	results := make([]*driver.Result, 0, len(cmds))
	for i, cmd := range cmds {
		if i >= len(commands) {
			break
		}
		result := newResult(commands[i])
		result.Latency = time.Since(startTime)
		results = append(results, result)
		if cmd.Err() != nil && cmd.Err() != redis.Nil {
			result.Err = cmd.Err()
			break
		}
		appendReplyRows(result, cmd, opts.MaxResultBytes)
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
	inputs, commands, err := splitCommands(script)
	if err != nil {
		return nil, err
	}

	startTime := time.Now()
	shards, err := s.driver.client.DoEachShard(ctx, inputs)
	if err != nil {
		// A client of one server says so itself, without the wrap a failure of
		// the cluster's own state takes.
		var notCluster *notClusterError
		if errors.As(err, &notCluster) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to query for each shard: %w", err)
	}

	results := make([]*driver.Result, len(inputs))
	for i := range inputs {
		results[i] = newResult(commands[i])
	}
	for _, shard := range shards {
		for i, cmd := range shard {
			if i >= len(results) {
				break
			}
			if cmd.Err() != nil && cmd.Err() != redis.Nil {
				results[i].Err = cmd.Err()
				return results[:i+1], nil
			}
			appendReplyRows(results[i], cmd, opts.MaxResultBytes)
		}
	}
	for _, result := range results {
		result.Latency = time.Since(startTime)
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
// result, and marks the result truncated.
func appendReplyRows(result *driver.Result, cmd *redis.Cmd, maxBytes int64) {
	if result.Truncated {
		return
	}
	totalBytes := driver.HeaderByteSize(result.Columns)
	add := func(row []driver.Value) {
		result.Rows = append(result.Rows, row)
		totalBytes += driver.RowByteSize(row)
		if maxBytes > 0 && totalBytes > maxBytes {
			result.Truncated = true
		}
	}
	if cmd.Err() == redis.Nil {
		add([]driver.Value{{Kind: driver.Int32, Int64: 1}, {Kind: driver.Null}})
		return
	}
	if list, ok := cmd.Val().([]any); ok {
		for i, value := range list {
			if result.Truncated {
				return
			}
			add(replyRow(i+1, value))
		}
		return
	}
	add(replyRow(1, cmd.Val()))
}

// replyRow is one row of an answer: the element's position and its value as a
// string.
func replyRow(i int, value any) []driver.Value {
	return []driver.Value{
		{Kind: driver.Int32, Int64: int64(i)},
		{Kind: driver.String, String: fmt.Sprintf("%v", value)},
	}
}

// splitCommands reads one command per line, skipping lines that hold none, and
// returns each command's arguments beside the line it came from. A result is
// reported per command rather than per line, so the two are built together to
// keep them aligned. A line is parsed as a shell parses one, so a quoted
// argument may hold spaces, and a line's `#` starts a comment.
func splitCommands(statement string) ([][]any, []string, error) {
	var inputs [][]any
	var commands []string
	for _, line := range strings.Split(statement, "\n") {
		fields, err := shlex.Split(line)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to split command %s: %w", line, err)
		}
		if len(fields) == 0 {
			continue
		}
		var input []any
		for _, v := range fields {
			input = append(input, v)
		}
		inputs = append(inputs, input)
		commands = append(commands, line)
	}
	return inputs, commands, nil
}
