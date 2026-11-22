package driver

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplitCommands(t *testing.T) {
	tests := []struct {
		name      string
		statement string
		wantArgs  [][]any
		wantLines []string
	}{
		{
			name:      "one command per line",
			statement: "GET a\nGET b",
			wantArgs:  [][]any{{"GET", "a"}, {"GET", "b"}},
			wantLines: []string{"GET a", "GET b"},
		},
		{
			name:      "a blank line does not shift the commands after it",
			statement: "GET a\n\nGET b\n   \nGET c",
			wantArgs:  [][]any{{"GET", "a"}, {"GET", "b"}, {"GET", "c"}},
			wantLines: []string{"GET a", "GET b", "GET c"},
		},
		{
			name:      "quoted argument",
			statement: `SET a "two words"`,
			wantArgs:  [][]any{{"SET", "a", "two words"}},
			wantLines: []string{`SET a "two words"`},
		},
		{
			name:      "a hash is an argument's text, as Redis has no comment",
			statement: "SET a #channel\nGET a#b",
			wantArgs:  [][]any{{"SET", "a", "#channel"}, {"GET", "a#b"}},
			wantLines: []string{"SET a #channel", "GET a#b"},
		},
		{
			name:      "a line ends with a carriage return",
			statement: "GET a\r\nGET b\r",
			wantArgs:  [][]any{{"GET", "a"}, {"GET", "b"}},
			wantLines: []string{"GET a\r", "GET b\r"},
		},
		{
			name:      "no commands",
			statement: "\n  \n",
			wantArgs:  nil,
			wantLines: nil,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			commands, err := splitCommands(tc.statement)
			require.NoError(t, err)
			var args [][]any
			var lines []string
			for _, command := range commands {
				args = append(args, command.args)
				lines = append(lines, command.text)
			}
			require.Equal(t, tc.wantArgs, args)
			require.Equal(t, tc.wantLines, lines)
		})
	}
}

// TestSplitCommandsRanges holds each command's byte range to the line it was
// written on: a task run's log locates the command by it.
func TestSplitCommandsRanges(t *testing.T) {
	script := "GET a\n\nSET b \"two words\"\nGET c"
	commands, err := splitCommands(script)
	require.NoError(t, err)
	require.Len(t, commands, 3)
	for _, command := range commands {
		require.Equal(t, command.text, script[command.start:command.end])
	}
	require.Equal(t, 0, commands[0].start)
	require.Equal(t, len("GET a"), commands[0].end)
	require.Equal(t, len("GET a\n\n"), commands[1].start)
	require.Equal(t, "SET b \"two words\"", commands[1].text)
	require.Equal(t, len(script)-len("GET c"), commands[2].start)
}

func TestSplitCommandsRejectsAnUnclosedQuote(t *testing.T) {
	_, err := splitCommands(`SET a "unclosed`)
	require.ErrorContains(t, err, "EOF found when expecting closing quote")

	_, err = splitCommands(`SET a 'unclosed`)
	require.ErrorContains(t, err, "EOF found when expecting closing quote")
}

func TestSplitLine(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		want    []string
		wantErr bool
	}{
		{name: "plain", line: "GET a", want: []string{"GET", "a"}},
		{name: "blanks", line: " \t GET  a \t ", want: []string{"GET", "a"}},
		{name: "quoted spaces", line: `SET a "two words"`, want: []string{"SET", "a", "two words"}},
		{name: "single quotes keep a backslash", line: `SET a 'c:\path'`, want: []string{"SET", "a", `c:\path`}},
		{name: "an escaped quote inside double quotes", line: `SET a "say \"hi\""`, want: []string{"SET", "a", `say "hi"`}},
		{name: "quotes join one argument", line: `SET a b"c d"e`, want: []string{"SET", "a", "bc de"}},
		{name: "an escaped blank", line: `SET a two\ words`, want: []string{"SET", "a", "two words"}},
		{name: "an empty quoted argument", line: `SET a ""`, want: []string{"SET", "a", ""}},
		{name: "a hash is text", line: `SET a #channel`, want: []string{"SET", "a", "#channel"}},
		{name: "a hash inside quotes", line: `SET a "#channel"`, want: []string{"SET", "a", "#channel"}},
		{name: "no arguments", line: "   ", want: nil},
		{name: "an unclosed double quote", line: `SET a "open`, wantErr: true},
		{name: "an unclosed single quote", line: `SET a 'open`, wantErr: true},
		{name: "a trailing backslash", line: `SET a value\`, wantErr: true},
		{name: "a backslash before the closing quote", line: `SET a "value\"`, wantErr: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmds, err := splitCommands(tc.line)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tc.want == nil {
				require.Empty(t, cmds)
				return
			}
			require.Len(t, cmds, 1)
			got := make([]string, len(cmds[0].args))
			for i, a := range cmds[0].args {
				got[i] = a.(string)
			}
			require.Equal(t, tc.want, got)
		})
	}
}
