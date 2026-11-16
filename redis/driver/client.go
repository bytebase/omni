package driver

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/redis/go-redis/v9"
)

// client is what the driver asks of a Redis client: the commands sync reads and
// the pipelines query and execute run. It is the seam a test replaces with a
// tape, so the driver's own code is what a recording and a replay exercise; the
// live implementation below is the one call per method it takes to reach
// go-redis.
type client interface {
	Ping(ctx context.Context) error
	// Info reads one INFO section.
	Info(ctx context.Context, section string) (string, error)
	// ConfigGet reads one CONFIG GET parameter.
	ConfigGet(ctx context.Context, parameter string) (map[string]string, error)
	// Do sends commands as one pipeline on the driver's connection and returns
	// each command's answer, in order. An error means the pipeline could not be
	// sent at all, so no command has an answer; a server's refusal is the
	// command's own answer, which its Cmd carries.
	Do(ctx context.Context, commands [][]any) ([]*redis.Cmd, error)
	// DoEachShard sends commands as one pipeline to every node of a cluster,
	// ordered by node address, and returns each node's answers.
	DoEachShard(ctx context.Context, commands [][]any) ([][]*redis.Cmd, error)
	Close() error
}

// liveClient is the client the driver opens against a server.
type liveClient struct {
	client redis.UniversalClient
}

func (c *liveClient) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

func (c *liveClient) Close() error {
	return c.client.Close()
}

func (c *liveClient) Info(ctx context.Context, section string) (string, error) {
	return c.client.Info(ctx, section).Result()
}

func (c *liveClient) ConfigGet(ctx context.Context, parameter string) (map[string]string, error) {
	return c.client.ConfigGet(ctx, parameter).Result()
}

func (c *liveClient) Do(ctx context.Context, commands [][]any) ([]*redis.Cmd, error) {
	cmds := make([]*redis.Cmd, 0, len(commands))
	if _, err := c.client.Pipelined(ctx, func(p redis.Pipeliner) error {
		for _, command := range commands {
			cmds = append(cmds, p.Do(ctx, command...))
		}
		return nil
	}); err != nil && !isServerError(err) {
		return nil, err
	}
	return cmds, nil
}

// isServerError reports whether err is the server's answer to a command rather
// than a failure to send one. Both a refusal and a missing key are answers;
// redis.Nil is the latter.
func isServerError(err error) bool {
	var serverError redis.Error
	return errors.As(err, &serverError)
}

// notClusterError reports that an all-nodes query asked a client of one server.
// go-redis has a cluster client and a client of one server, and only the first
// can ask every node; the message is the one Bytebase answered with, naming the
// client that was asked.
type notClusterError struct {
	client any
}

func (e *notClusterError) Error() string {
	return fmt.Sprintf("expect *redis.ClusterClient, but get %T", e.client)
}

// shardCmds is one node's answers, kept with its address so that the shards come
// back in an order a recording repeats.
type shardCmds struct {
	addr string
	cmds []*redis.Cmd
}

func (c *liveClient) DoEachShard(ctx context.Context, commands [][]any) ([][]*redis.Cmd, error) {
	cluster, ok := c.client.(*redis.ClusterClient)
	if !ok {
		return nil, &notClusterError{client: c.client}
	}
	// ForEachShard runs the callbacks concurrently and returns once they have
	// all finished, so the mutex is what orders the appends, and the sort below
	// is what makes their order the node addresses' rather than the run's.
	var mu sync.Mutex
	var shards []shardCmds
	err := cluster.ForEachShard(ctx, func(ctx context.Context, client *redis.Client) error {
		cmds := make([]*redis.Cmd, 0, len(commands))
		_, _ = client.Pipelined(ctx, func(p redis.Pipeliner) error {
			for _, command := range commands {
				cmds = append(cmds, p.Do(ctx, command...))
			}
			return nil
		})
		mu.Lock()
		defer mu.Unlock()
		shards = append(shards, shardCmds{addr: client.Options().Addr, cmds: cmds})
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.SliceStable(shards, func(i, j int) bool { return shards[i].addr < shards[j].addr })
	answers := make([][]*redis.Cmd, len(shards))
	for i, shard := range shards {
		answers[i] = shard.cmds
	}
	return answers, nil
}
