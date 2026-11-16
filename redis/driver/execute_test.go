package driver

import (
	"context"
	"fmt"
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
// script with its error, so the command after it does not run.
func TestExecuteRefusedCommand(t *testing.T) {
	ctx := context.Background()
	d := openDriver(t, testcontainer.Login{}, "0")
	_, err := d.Execute(ctx, "SET key text\nLPUSH key x\nSET after value", driver.ExecuteOptions{})
	require.ErrorContains(t, err, "WRONGTYPE")

	results := query(t, d, "GET after", driver.QueryOptions{})
	require.Len(t, results, 1)
	require.Len(t, results[0].Rows, 1)
	require.Equal(t, driver.Null, results[0].Rows[0][1].Kind, "the command after the refused one ran")
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

// logRecorder records what Execute reports.
type logRecorder struct {
	commands  []driver.Command
	responses []string
}

func (l *logRecorder) CommandExecute(command driver.Command) {
	l.commands = append(l.commands, command)
}

func (l *logRecorder) CommandResponse(affectedRows int64, perResult []int64, err error) {
	response := fmt.Sprintf("%d", affectedRows)
	if err != nil {
		response += ":" + err.Error()
	}
	l.responses = append(l.responses, response)
}

func (*logRecorder) TransactionControl(driver.TransactionOp, error) {}

func (*logRecorder) Retry(error, int, int) {}

// TestExecuteLogsItsSteps holds a change script's task run log to every command it
// ran, with the command's range in the script, and to the error of the command
// that ended it.
func TestExecuteLogsItsSteps(t *testing.T) {
	ctx := context.Background()
	d := openDriver(t, testcontainer.Login{}, "0")
	script := "SET greeting \"hello world\"\n\nINCR counter"
	log := &logRecorder{}
	_, err := d.Execute(ctx, script, driver.ExecuteOptions{Log: log})
	require.NoError(t, err)
	require.Len(t, log.commands, 2)
	for _, command := range log.commands {
		require.Equal(t, command.Text, script[command.ByteStart:command.ByteEnd])
	}
	require.Equal(t, "SET greeting \"hello world\"", log.commands[0].Text)
	require.Equal(t, "INCR counter", log.commands[1].Text)
	require.Equal(t, []string{"0", "0"}, log.responses)
}

// TestExecuteLogsTheCommandThatFailed holds the step that ended a script: its
// command and its error both reach the log.
func TestExecuteLogsTheCommandThatFailed(t *testing.T) {
	ctx := context.Background()
	d := openDriver(t, testcontainer.Login{}, "0")
	log := &logRecorder{}
	_, err := d.Execute(ctx, "SET key text\nLPUSH key x", driver.ExecuteOptions{Log: log})
	require.ErrorContains(t, err, "WRONGTYPE")
	require.Len(t, log.commands, 2)
	require.Len(t, log.responses, 2)
	require.Equal(t, "0", log.responses[0])
	require.Contains(t, log.responses[1], "WRONGTYPE")
}
