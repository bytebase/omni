package testcontainer

import (
	"context"
	"encoding/json"
	"math"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

// everyKind is one value of every reply kind a Redis command can answer with.
func everyKind() []any {
	return []any{
		nil,
		"",
		"a string with \"quotes\", a \\ backslash and é",
		int64(0),
		int64(-7),
		int64(math.MaxInt64),
		float64(0),
		float64(-1.5),
		float64(math.Inf(1)),
		true,
		false,
		[]any{},
		[]any{"a", int64(1), []any{nil, "b"}},
		map[any]any{},
		map[any]any{"b": int64(2), "a": "one", "nested": map[any]any{"x": true}},
	}
}

// TestTheTapeHoldsEveryReplyKind holds the tape's encoding to what a command can
// answer: a replay hands the driver the value a recording's client decoded.
func TestTheTapeHoldsEveryReplyKind(t *testing.T) {
	for i, kind := range everyKind() {
		encoded, err := encodeValue(kind)
		require.NoError(t, err, "value %d", i)
		text, err := json.Marshal(encoded)
		require.NoError(t, err)
		var read value
		require.NoError(t, unmarshalJSON(text, &read), "value %d", i)
		decoded, err := read.decode()
		require.NoError(t, err, "value %d", i)
		require.Equal(t, kind, decoded, "value %d", i)
		// The same value writes the same bytes, whatever run it is recorded in.
		again, err := json.Marshal(read)
		require.NoError(t, err)
		require.Equal(t, string(text), string(again), "value %d", i)
	}
}

// TestTheTapeHoldsNoOtherKind holds the encoding to failing a recording on a
// reply it has no form for, rather than writing a tape a replay cannot read.
func TestTheTapeHoldsNoOtherKind(t *testing.T) {
	for _, kind := range []any{int(1), uint64(1), []byte("bytes"), map[string]any{"a": 1}} {
		_, err := encodeValue(kind)
		require.Error(t, err)
	}
}

// TestTheTapeHoldsCommandErrors holds a command's error to the error go-redis
// hands the driver: redis.Nil comes back as the sentinel the driver compares by
// identity, and a refusal keeps the text HasErrorPrefix reads.
func TestTheTapeHoldsCommandErrors(t *testing.T) {
	value := redis.NewCmd(context.Background(), "GET", "key")
	value.SetVal("answer")
	missing := redis.NewCmd(context.Background(), "GET", "missing")
	missing.SetErr(redis.Nil)
	refused := redis.NewCmd(context.Background(), "INCR", "key")
	refused.SetErr(serverError("NOPERM User noconfig has no permissions to run the 'config|get' command"))

	encoded, err := encodeCommands([]*redis.Cmd{value, missing, refused})
	require.NoError(t, err)
	text, err := json.Marshal(answer{Commands: encoded})
	require.NoError(t, err)
	var read answer
	require.NoError(t, unmarshalJSON(text, &read))
	cmds, err := read.commands(context.Background(), [][]any{{"GET", "key"}, {"GET", "missing"}, {"INCR", "key"}})
	require.NoError(t, err)

	require.Len(t, cmds, 3)
	require.NoError(t, cmds[0].Err())
	require.Equal(t, "answer", cmds[0].Val())
	require.Equal(t, redis.Nil, cmds[1].Err())
	require.True(t, cmds[1].Err() == redis.Nil, "a replay's missing key is the sentinel itself")
	require.ErrorContains(t, cmds[2].Err(), "NOPERM")
	require.True(t, isServerError(cmds[2].Err()))
}

// TestTheTapeNarrowsInfo holds an INFO answer to the lines the driver reads: the
// counters of a whole section move between runs.
func TestTheTapeNarrowsInfo(t *testing.T) {
	server := "# Server\r\nredis_version:7.4.12\r\nuptime_in_seconds:43\r\n"
	require.Equal(t, "redis_version:7.4.12\r", infoLines("server", server))
	keyspace := "# Keyspace\r\ndb0:keys=3,expires=0,avg_ttl=0\r\ndb2:keys=1,expires=0,avg_ttl=0\r\n"
	require.Equal(t, "db0:keys=3,expires=0,avg_ttl=0\r\ndb2:keys=1,expires=0,avg_ttl=0\r", infoLines("keyspace", keyspace))
	require.Empty(t, infoLines("memory", "# Memory\r\nused_memory:1000\r\n"))
}
