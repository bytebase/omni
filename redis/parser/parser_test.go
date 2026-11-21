package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		want    []string
		wantErr bool
	}{
		{name: "plain", line: "GET a", want: []string{"GET", "a"}},
		{name: "blanks", line: " \t GET  a \t ", want: []string{"GET", "a"}},
		{name: "quoted spaces", line: `SET a "two words"`, want: []string{"SET", "a", "two words"}},
		{name: "single quotes keep backslash", line: `SET a 'c:\path'`, want: []string{"SET", "a", `c:\path`}},
		{name: "escaped quote in double quotes", line: `SET a "say \"hi\""`, want: []string{"SET", "a", `say "hi"`}},
		{name: "double quote escapes", line: `SET a "a\nb\rc\td\be\af\\g\x41"`, want: []string{"SET", "a", "a\nb\rc\td\be\af\\gA"}},
		{name: "non-hex backslash x", line: `SET a "\xZZ"`, want: []string{"SET", "a", "xZZ"}},
		{name: "quotes join one argument", line: `SET a b"c d"e`, want: []string{"SET", "a", "bc de"}},
		{name: "escaped blank", line: `SET a two\ words`, want: []string{"SET", "a", "two words"}},
		{name: "empty double-quoted argument", line: `SET a ""`, want: []string{"SET", "a", ""}},
		{name: "empty single-quoted argument", line: `SET a ''`, want: []string{"SET", "a", ""}},
		{name: "hash is argument text", line: `SET a #channel`, want: []string{"SET", "a", "#channel"}},
		{name: "hash inside quotes", line: `SET a "#channel"`, want: []string{"SET", "a", "#channel"}},
		{name: "no arguments", line: "   \t\r ", want: nil},
		{name: "unclosed double quote", line: `SET a "open`, wantErr: true},
		{name: "unclosed single quote", line: `SET a 'open`, wantErr: true},
		{name: "trailing backslash", line: `SET a value\`, wantErr: true},
		{name: "backslash before closing quote", line: `SET a "value\"`, wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Tokenize(tc.line)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestParseCommand(t *testing.T) {
	valid := []struct {
		line    string
		wantCmd string
		wantSub string
		wantArg []string
	}{
		{line: "GET mykey", wantCmd: "get", wantArg: []string{"GET", "mykey"}},
		{line: "gEt mykey", wantCmd: "get", wantArg: []string{"gEt", "mykey"}},
		{line: "SET k v EX 10", wantCmd: "set", wantArg: []string{"SET", "k", "v", "EX", "10"}},
		{line: "PING", wantCmd: "ping", wantArg: []string{"PING"}},
		{line: "PING hello", wantCmd: "ping", wantArg: []string{"PING", "hello"}},
		{line: "DBSIZE", wantCmd: "dbsize", wantArg: []string{"DBSIZE"}},
		{line: "COMMAND", wantCmd: "command", wantArg: []string{"COMMAND"}},
		{line: "COMMAND COUNT", wantCmd: "command", wantSub: "count", wantArg: []string{"COMMAND", "COUNT"}},
		{line: "COMMAND INFO GET SET", wantCmd: "command", wantSub: "info", wantArg: []string{"COMMAND", "INFO", "GET", "SET"}},
		{line: "ACL LIST", wantCmd: "acl", wantSub: "list", wantArg: []string{"ACL", "LIST"}},
		{line: "acl getuser default", wantCmd: "acl", wantSub: "getuser", wantArg: []string{"acl", "getuser", "default"}},
		{line: "CLIENT NO-EVICT ON", wantCmd: "client", wantSub: "no-evict", wantArg: []string{"CLIENT", "NO-EVICT", "ON"}},
		{line: "MEMORY MALLOC-STATS", wantCmd: "memory", wantSub: "malloc-stats", wantArg: []string{"MEMORY", "MALLOC-STATS"}},
		{line: "RESTORE-ASKING k 0 val", wantCmd: "restore-asking", wantArg: []string{"RESTORE-ASKING", "k", "0", "val"}},
		{line: `SET k "hello\nworld"`, wantCmd: "set", wantArg: []string{"SET", "k", "hello\nworld"}},
		{line: "GET #channel", wantCmd: "get", wantArg: []string{"GET", "#channel"}},
	}

	for _, tc := range valid {
		t.Run("valid/"+tc.line, func(t *testing.T) {
			cmd, err := ParseCommand(tc.line)
			require.NoError(t, err)
			require.Equal(t, tc.wantCmd, cmd.Name)
			require.Equal(t, tc.wantSub, cmd.Subcommand)
			require.Equal(t, tc.wantArg, cmd.Args)
			require.Equal(t, tc.line, cmd.Text)
			require.Equal(t, 0, cmd.Start)
			require.Equal(t, len(tc.line), cmd.End)
		})
	}

	invalid := []struct {
		name string
		line string
	}{
		{name: "empty", line: ""},
		{name: "whitespace", line: "   \t "},
		{name: "multiline", line: "GET a\nGET b"},
		{name: "unknown command", line: "NOTACOMMAND a"},
		{name: "pipe subcommand syntax", line: "acl|list"},
		{name: "hash not a comment", line: "# comment"},
		{name: "missing required arg", line: "GET"},
		{name: "too many args for fixed arity", line: "GET a b"},
		{name: "extra arg on zero-arg command", line: "DBSIZE extra"},
		{name: "too few args for variadic", line: "SET onlykey"},
		{name: "container without subcommand", line: "ACL"},
		{name: "container unknown subcommand", line: "ACL NOPE"},
		{name: "command unknown subcommand", line: "COMMAND NOPE"},
		{name: "subcommand too many args", line: "ACL LIST extra"},
		{name: "command subcommand too many args", line: "COMMAND COUNT extra"},
		{name: "subcommand too few args", line: "ACL GETUSER"},
		{name: "unclosed quote", line: `GET "unclosed`},
		{name: "trailing backslash", line: `GET a\`},
	}

	for _, tc := range invalid {
		t.Run("invalid/"+tc.name, func(t *testing.T) {
			_, err := ParseCommand(tc.line)
			require.Error(t, err)
			var pe *ParseError
			require.ErrorAs(t, err, &pe)
		})
	}
}

func TestParseScript(t *testing.T) {
	script := "GET a\n\nSET b \"two words\"\r\nACL LIST\n   \nPING"
	cmds, err := Parse(script)
	require.NoError(t, err)
	require.Len(t, cmds, 4)

	require.Equal(t, "get", cmds[0].Name)
	require.Equal(t, []string{"GET", "a"}, cmds[0].Args)
	require.Equal(t, cmds[0].Text, script[cmds[0].Start:cmds[0].End])

	require.Equal(t, "set", cmds[1].Name)
	require.Equal(t, []string{"SET", "b", "two words"}, cmds[1].Args)
	require.Equal(t, cmds[1].Text, script[cmds[1].Start:cmds[1].End])

	require.Equal(t, "acl", cmds[2].Name)
	require.Equal(t, "list", cmds[2].Subcommand)
	require.Equal(t, cmds[2].Text, script[cmds[2].Start:cmds[2].End])

	require.Equal(t, "ping", cmds[3].Name)
	require.Equal(t, cmds[3].Text, script[cmds[3].Start:cmds[3].End])

	// Empty / whitespace-only script returns zero commands and no error.
	empty, err := Parse("   \n\t\r\n")
	require.NoError(t, err)
	require.Empty(t, empty)

	// Error reports accurate 1-based line number and byte position.
	_, err = Parse("GET a\n\nGET a b")
	require.Error(t, err)
	var pe *ParseError
	require.ErrorAs(t, err, &pe)
	require.Equal(t, 3, pe.Line)
	require.Equal(t, len("GET a\n\n"), pe.Position)
}

func TestLookupReturnsCopy(t *testing.T) {
	spec, ok := Lookup("ACL")
	require.True(t, ok)
	require.Equal(t, "acl", spec.Name)
	require.Equal(t, 2, spec.Subcommands["list"])

	// Mutating the returned copy must not affect the package registry.
	spec.Subcommands["list"] = 99
	fresh, ok := Lookup("acl")
	require.True(t, ok)
	require.Equal(t, 2, fresh.Subcommands["list"])

	require.Len(t, Commands(), 250)
}

// TestRegistryCoversAnalysisReadOnlyCommands ensures every core read-only
// command in redis/analysis/testdata/core_readonly_commands.txt is known to
// redis/parser's Redis 7.4.12 registry.
func TestRegistryCoversAnalysisReadOnlyCommands(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("..", "analysis", "testdata", "core_readonly_commands.txt"))
	require.NoError(t, err)
	for line := range strings.SplitSeq(string(raw), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Fields(line)
		spec, ok := Lookup(parts[0])
		require.True(t, ok, "missing top-level command %q from core_readonly_commands.txt", line)
		if len(parts) == 2 {
			_, subOK := spec.Subcommands[parts[1]]
			require.True(t, subOK, "missing subcommand %q from core_readonly_commands.txt", line)
		}
	}
}
