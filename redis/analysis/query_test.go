package analysis

import (
	"errors"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	redisparser "github.com/bytebase/omni/redis/parser"
)

func TestValidateQuery(t *testing.T) {
	tests := []struct {
		statement   string
		readOnly    bool
		returnsData bool
		wantErr     bool
	}{
		{
			statement:   "",
			readOnly:    true,
			returnsData: true,
		},
		{
			statement:   "   \n\t\r\n",
			readOnly:    true,
			returnsData: true,
		},
		{
			statement:   "get hello",
			readOnly:    true,
			returnsData: true,
		},
		{
			statement:   "GET hello\n\nHGET h f\r\nOBJECT ENCODING hello",
			readOnly:    true,
			returnsData: true,
		},
		{
			statement:   "set hello 1",
			readOnly:    false,
			returnsData: false,
		},
		{
			statement:   "GET hello\nSET hello 1",
			readOnly:    false,
			returnsData: false,
		},
		{
			statement:   "memory usage hello",
			readOnly:    true,
			returnsData: true,
		},
		{
			statement:   "MEMORY DOCTOR",
			readOnly:    false,
			returnsData: false,
		},
		{
			statement:   "OBJECT ENCODING key",
			readOnly:    true,
			returnsData: true,
		},
		{
			statement:   "OBJECT HELP",
			readOnly:    false,
			returnsData: false,
		},
		{
			statement:   "XINFO STREAM mystream",
			readOnly:    true,
			returnsData: true,
		},
		{
			statement:   "XINFO GROUPS mystream",
			readOnly:    true,
			returnsData: true,
		},
		{
			statement:   "XINFO CONSUMERS mystream mygroup",
			readOnly:    true,
			returnsData: true,
		},
		{
			statement:   "XINFO HELP",
			readOnly:    false,
			returnsData: false,
		},
		{
			statement:   "CONFIG GET *",
			readOnly:    false,
			returnsData: false,
		},
		{
			statement:   "CLIENT LIST",
			readOnly:    false,
			returnsData: false,
		},
		{
			statement:   "ACL LIST",
			readOnly:    false,
			returnsData: false,
		},
		{
			statement:   "SCRIPT EXISTS abc",
			readOnly:    false,
			returnsData: false,
		},
		// RedisJSON reads. The command namespaces with a dot inside one token,
		// so it can only match readCommands/moduleReadCommands, never the
		// container-subcommand fallback.
		{
			statement:   "JSON.GET doc",
			readOnly:    true,
			returnsData: true,
		},
		{
			statement:   "json.get doc $.name",
			readOnly:    true,
			returnsData: true,
		},
		{
			statement:   "JSON.DEBUG MEMORY doc",
			readOnly:    true,
			returnsData: true,
		},
		{
			statement:   `JSON.SET doc $ {"a":1}`,
			readOnly:    false,
			returnsData: false,
		},
		{
			statement:   "JSON.DEL doc $.a",
			readOnly:    false,
			returnsData: false,
		},
		// Other Redis Stack modules.
		{
			statement:   "FT.SEARCH idx *",
			readOnly:    true,
			returnsData: true,
		},
		{
			statement:   "TS.GET series",
			readOnly:    true,
			returnsData: true,
		},
		{
			statement:   "BF.EXISTS filter item",
			readOnly:    true,
			returnsData: true,
		},
		// FT.DICTADD mutates a dictionary even though Redis flags it readonly.
		{
			statement:   "FT.DICTADD dict term",
			readOnly:    false,
			returnsData: false,
		},
		// Core read commands added in Redis 7.4.
		{
			statement:   "HTTL h FIELDS 1 f",
			readOnly:    true,
			returnsData: true,
		},
		// Read-only script mode only blocks commands flagged "write", so Lua
		// still reaches SELECT and the readonly-flagged module commands that
		// mutate state. See knownExclusions in
		// TestReadCommandsMatchServerFlags.
		{
			statement:   "EVAL_RO \"return redis.call('GET','k')\" 0",
			readOnly:    false,
			returnsData: false,
		},
		{
			statement:   "EVALSHA_RO abc 0",
			readOnly:    false,
			returnsData: false,
		},
		{
			statement:   "FCALL_RO fn 0",
			readOnly:    false,
			returnsData: false,
		},
		// SELECT switches the connection to another logical database, so a
		// later read escapes the database the request was authorized against.
		{
			statement:   "select 3",
			readOnly:    false,
			returnsData: false,
		},
		{
			statement:   "select 3\nkeys *",
			readOnly:    false,
			returnsData: false,
		},
		// Redis has no comment syntax: '#' is argument text, matching redis/driver.
		{
			statement:   "GET #channel",
			readOnly:    true,
			returnsData: true,
		},
		{
			statement:   "# not a comment\nGET k",
			readOnly:    false,
			returnsData: false,
		},
		// Quoted and escaped arguments.
		{
			statement:   `GET "my \"key\""`,
			readOnly:    true,
			returnsData: true,
		},
		{
			statement:   `GET 'c:\path'`,
			readOnly:    true,
			returnsData: true,
		},
		{
			statement:   `GET two\ words`,
			readOnly:    true,
			returnsData: true,
		},
		{
			statement:   `GET "line1\nline2\t\x41"`,
			readOnly:    true,
			returnsData: true,
		},
		{
			statement:   `"GET" key`,
			readOnly:    true,
			returnsData: true,
		},
		{
			statement:   `'OBJECT' "ENCODING" key`,
			readOnly:    true,
			returnsData: true,
		},
		{
			statement:   "GET k1\r\nHGET h f\r\nJSON.GET doc\r\n",
			readOnly:    true,
			returnsData: true,
		},
		// Malformed quoting and escapes.
		{
			statement: `GET "unclosed`,
			wantErr:   true,
		},
		{
			statement: `GET 'unclosed`,
			wantErr:   true,
		},
		{
			statement: `GET trailing\`,
			wantErr:   true,
		},
		{
			statement: `GET "escaped\"`,
			wantErr:   true,
		},
		{
			statement: "SET a 1\nGET \"unclosed",
			wantErr:   true,
		},
	}

	for _, tc := range tests {
		gotReadOnly, gotReturnsData, err := ValidateQuery(tc.statement)
		if tc.wantErr {
			require.Error(t, err, tc.statement)
			var pe *redisparser.ParseError
			require.True(t, errors.As(err, &pe), "expected *redis/parser.ParseError for %q, got %T: %v", tc.statement, err, err)
			require.False(t, gotReadOnly, tc.statement)
			require.False(t, gotReturnsData, tc.statement)
		} else {
			require.NoError(t, err, tc.statement)
			require.Equal(t, tc.readOnly, gotReadOnly, tc.statement)
			require.Equal(t, tc.returnsData, gotReturnsData, tc.statement)
		}
	}
}

func TestIsReadOnlyCommand(t *testing.T) {
	require.False(t, IsReadOnlyCommand(nil))
	require.False(t, IsReadOnlyCommand([]string{}))
	require.True(t, IsReadOnlyCommand([]string{"GET", "k"}))
	require.True(t, IsReadOnlyCommand([]string{"JSON.GET", "doc"}))
	require.True(t, IsReadOnlyCommand([]string{"OBJECT", "ENCODING", "k"}))
	require.False(t, IsReadOnlyCommand([]string{"OBJECT"}))
	require.False(t, IsReadOnlyCommand([]string{"MEMORY"}))
	require.False(t, IsReadOnlyCommand([]string{"XINFO"}))
	require.False(t, IsReadOnlyCommand([]string{"SET", "k", "v"}))
}

// TestReadCommandsMatchServerFlags diffs the core allowlist against a snapshot
// of the flags Redis reports for itself, so a Redis upgrade that ships new read
// commands fails here instead of reaching a caller. See the snapshot header
// for how to regenerate it.
func TestReadCommandsMatchServerFlags(t *testing.T) {
	// knownExclusions holds commands Redis flags "readonly" that are
	// deliberately kept out of the gate. Every entry needs a reason: the flag
	// means "does not write the keyspace", not "has no side effects".
	//
	// SELECT is not listed. Redis does not flag it readonly, so the reverse
	// check below is what keeps it out of the allowlist.
	knownExclusions := map[string]string{
		// Read-only script mode only rejects commands flagged "write", so a
		// script still reaches SELECT and every readonly-flagged module command
		// that mutates state. Verified on redis-stack-server 7.4.7: running
		// EVAL_RO "redis.call('SELECT',1); return redis.call('GET','marker')"
		// against database 0 returns database 1's value, and a no-writes
		// function called through FCALL_RO both crosses databases and adds a
		// RediSearch dictionary term via FT.DICTADD.
		"eval_ro":    "Lua reaches SELECT and readonly-flagged mutating module commands",
		"evalsha_ro": "same as eval_ro",
		"fcall_ro":   "same as eval_ro, through a no-writes function",
	}

	content, err := os.ReadFile("testdata/core_readonly_commands.txt")
	require.NoError(t, err)

	serverReadOnly := map[string]bool{}
	for line := range strings.SplitSeq(string(content), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		serverReadOnly[line] = true
	}
	require.NotEmpty(t, serverReadOnly, "snapshot has no commands")

	allowed := map[string]bool{}
	for command := range readCommands {
		allowed[command] = true
	}
	for container, subcommands := range doubleNameReadCommands {
		for subcommand := range subcommands {
			allowed[container+" "+subcommand] = true
		}
	}

	var rejected []string
	for command := range serverReadOnly {
		if allowed[command] {
			continue
		}
		if _, ok := knownExclusions[command]; ok {
			continue
		}
		rejected = append(rejected, command)
	}
	slices.Sort(rejected)
	// assert, not require: a run that breaks both directions must report both.
	assert.Emptyf(t, rejected,
		"Redis flags these commands readonly but the gate rejects them: %v. "+
			"Add each to readCommands or doubleNameReadCommands, or to knownExclusions with a reason.",
		rejected)

	var unflagged []string
	for command := range allowed {
		if !serverReadOnly[command] {
			unflagged = append(unflagged, command)
		}
	}
	slices.Sort(unflagged)
	assert.Emptyf(t, unflagged,
		"the gate accepts these but Redis does not flag them readonly: %v. "+
			"Either the name is wrong, or the command has side effects and must be removed.",
		unflagged)
}

// TestModuleReadCommandsWellFormed guards the module allowlist, which has no
// snapshot because module presence varies per instance. A key that is not a
// lowercase dotted name can never match, since IsReadOnlyCommand lowercases the
// first token and module commands carry the dot inside it.
func TestModuleReadCommandsWellFormed(t *testing.T) {
	require.NotEmpty(t, moduleReadCommands)
	for command := range moduleReadCommands {
		require.Equal(t, strings.ToLower(command), command, "module command must be lowercase")
		require.Contains(t, command, ".", "module command must be dot-namespaced")
		require.False(t, readCommands[command], "module command duplicated in readCommands: %s", command)
	}
}

// TestReadCommandsRegisteredInParser verifies every core read-only command and
// subcommand in redis/analysis is present in redis/parser's command table.
func TestReadCommandsRegisteredInParser(t *testing.T) {
	for command := range readCommands {
		spec, ok := redisparser.Lookup(command)
		require.True(t, ok, "readCommands entry %q missing from redis/parser", command)
		require.NotZero(t, spec.StandaloneArity, "readCommands entry %q has zero standalone arity in redis/parser", command)
	}
	for container, subcommands := range doubleNameReadCommands {
		spec, ok := redisparser.Lookup(container)
		require.True(t, ok, "doubleNameReadCommands container %q missing from redis/parser", container)
		for sub := range subcommands {
			_, subOK := spec.Subcommands[sub]
			require.True(t, subOK, "doubleNameReadCommands entry %q %q missing from redis/parser", container, sub)
		}
	}
}
