package driver

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/driver"
	"github.com/bytebase/omni/redis/internal/testcontainer"
)

// TestExecute runs a change script as one pipeline and holds its result: a
// Redis command reports no row count, so the rows affected are 0.
func TestExecute(t *testing.T) {
	ctx := context.Background()
	d := openDriver(t, testcontainer.Login{}, "0")
	affected, err := d.Execute(ctx, "SET greeting hello\nINCR counter", driver.ExecuteOptions{})
	require.NoError(t, err)
	require.Zero(t, affected)

	results := query(t, d, "GET greeting\nGET counter", driver.QueryOptions{})
	require.Len(t, results, 2)
	_, value := row(t, results[0].Rows[0])
	require.Equal(t, "hello", value)
	_, value = row(t, results[1].Rows[0])
	require.Equal(t, "1", value)
}

// TestExecuteMissingKeyIsNotAnError covers a read whose answer is a missing key:
// redis.Nil is an answer, not a failed script.
func TestExecuteMissingKeyIsNotAnError(t *testing.T) {
	d := openDriver(t, testcontainer.Login{}, "0")
	_, err := d.Execute(context.Background(), "GET missing", driver.ExecuteOptions{})
	require.NoError(t, err)
}

// TestExecuteRefusedCommand covers a command the server refuses: it ends the
// script with its error.
func TestExecuteRefusedCommand(t *testing.T) {
	ctx := context.Background()
	d := openDriver(t, testcontainer.Login{}, "0")
	_, err := d.Execute(ctx, "SET key text\nLPUSH key x", driver.ExecuteOptions{})
	require.ErrorContains(t, err, "WRONGTYPE")
}

// TestExecuteCreateDatabase holds the error a database-creation step answers:
// Redis has no statement that creates one.
func TestExecuteCreateDatabase(t *testing.T) {
	_, err := (&Driver{}).Execute(context.Background(), "SET key value", driver.ExecuteOptions{CreateDatabase: true})
	require.ErrorContains(t, err, "redis: cannot create database")
}

// TestExecuteUnreadableQuoting covers a line whose quoting cannot be read: there
// is no command to send, so the call fails before reaching the server.
func TestExecuteUnreadableQuoting(t *testing.T) {
	_, err := (&Driver{}).Execute(context.Background(), `SET key "unclosed`, driver.ExecuteOptions{})
	require.ErrorContains(t, err, "failed to split command")
}

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
			name:      "no commands",
			statement: "\n  \n",
			wantArgs:  nil,
			wantLines: nil,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			args, lines, err := splitCommands(tc.statement)
			require.NoError(t, err)
			require.Equal(t, tc.wantArgs, args)
			require.Equal(t, tc.wantLines, lines)
			require.Len(t, lines, len(args))
		})
	}
}

func TestSplitCommandsRejectsAnUnclosedQuote(t *testing.T) {
	_, _, err := splitCommands(`SET a "unclosed`)
	require.Error(t, err)
}
