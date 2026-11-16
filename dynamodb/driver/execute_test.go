package driver

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/driver"
	"github.com/bytebase/omni/dynamodb/internal/testcontainer"
)

// recorder is an ExecuteLogger that keeps what Execute reported, in order.
type recorder struct {
	commands     []driver.Command
	responses    []response
	transactions []driver.TransactionOp
}

type response struct {
	affected  int64
	perResult []int64
	err       error
}

func (r *recorder) CommandExecute(command driver.Command) { r.commands = append(r.commands, command) }
func (r *recorder) CommandResponse(affected int64, perResult []int64, err error) {
	r.responses = append(r.responses, response{affected: affected, perResult: perResult, err: err})
}
func (r *recorder) TransactionControl(op driver.TransactionOp, _ error) {
	r.transactions = append(r.transactions, op)
}
func (r *recorder) Retry(error, int, int) {}

// TestExecute runs a script's statements in order, logs each, and reports the
// state they left.
func TestExecute(t *testing.T) {
	d := newDriver(t)
	log := &recorder{}
	script := `DELETE FROM "writes" WHERE id = 'w1';
INSERT INTO "writes" VALUE {'id':'w1','name':'first'};`
	affected, err := d.Execute(context.Background(), script, driver.ExecuteOptions{Log: log})
	require.NoError(t, err)
	require.Zero(t, affected, "DynamoDB answers no affected row count")

	require.Len(t, log.commands, 2)
	require.Equal(t, `DELETE FROM "writes" WHERE id = 'w1';`, log.commands[0].Text)
	require.Equal(t, 0, log.commands[0].ByteStart)
	require.Equal(t, len(`DELETE FROM "writes" WHERE id = 'w1';`), log.commands[0].ByteEnd)
	require.Equal(t, "\nINSERT INTO \"writes\" VALUE {'id':'w1','name':'first'};", log.commands[1].Text)
	require.Equal(t, len(`DELETE FROM "writes" WHERE id = 'w1';`), log.commands[1].ByteStart)
	require.Len(t, log.responses, 2)
	for _, response := range log.responses {
		require.NoError(t, response.err)
		require.Zero(t, response.affected)
		require.Equal(t, []int64{0}, response.perResult)
	}
	require.Empty(t, log.transactions, "the driver opens no transaction of its own")

	// The statements left one item, which a read of the table returns.
	session, err := d.Session(context.Background())
	require.NoError(t, err)
	results, err := session.Query(context.Background(), `SELECT * FROM "writes" WHERE id = 'w1'`, driver.QueryOptions{})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.NoError(t, results[0].Err)
	require.Equal(t, [][]string{{"w1", "first"}}, rowTexts(results[0]))

	cleanupWrites(t, d, "w1")
}

// TestExecuteEndsAtTheFirstFailure holds that a statement whose transaction
// failed is the call's error and that the statements after it do not run.
func TestExecuteEndsAtTheFirstFailure(t *testing.T) {
	d := newDriver(t)
	log := &recorder{}
	script := `DELETE FROM "writes" WHERE id = 'w2';
INSERT INTO "writes" VALUE {'id':'w2','name':'first'};
INSERT INTO "writes" VALUE {'id':'w2','name':'second'};
DELETE FROM "writes" WHERE id = 'w2';`
	affected, err := d.Execute(context.Background(), script, driver.ExecuteOptions{Log: log})
	require.Zero(t, affected)
	require.Error(t, err)
	// DynamoDB answers a write of an item that exists as a cancelled
	// transaction, whose per-statement reasons it reports as a list.
	require.Contains(t, err.Error(), "TransactionCanceledException")
	require.Contains(t, err.Error(), "[ValidationError]")

	require.Len(t, log.commands, 3, "the statement after the failure did not run")
	require.Len(t, log.responses, 3)
	require.NoError(t, log.responses[0].err)
	require.NoError(t, log.responses[1].err)
	require.Error(t, log.responses[2].err)
	require.Contains(t, log.responses[2].err.Error(), "TransactionCanceledException")

	cleanupWrites(t, d, "w2")
}

// TestExecuteRunsOnlyStatements holds that a script's comment-only segments
// hold no statement: partiql.Split decides, and nothing is sent for them.
func TestExecuteRunsOnlyStatements(t *testing.T) {
	d := newDriver(t)
	log := &recorder{}
	affected, err := d.Execute(context.Background(),
		"-- nothing here\n;\nDELETE FROM \"writes\" WHERE id = 'w3';\n/* nor here */",
		driver.ExecuteOptions{Log: log})
	require.NoError(t, err)
	require.Zero(t, affected)
	require.Len(t, log.commands, 1)
	require.Equal(t, "\nDELETE FROM \"writes\" WHERE id = 'w3';", log.commands[0].Text)
	require.Len(t, log.responses, 1)
	require.NoError(t, log.responses[0].err)
}

// TestExecuteWithoutALogger holds that a script runs the same with no logger.
func TestExecuteWithoutALogger(t *testing.T) {
	d := newDriver(t)
	affected, err := d.Execute(context.Background(), `DELETE FROM "writes" WHERE id = 'w4'`, driver.ExecuteOptions{})
	require.NoError(t, err)
	require.Zero(t, affected)
}

// TestExecuteWithoutStatements holds that a script of no statements wakes
// nothing.
func TestExecuteWithoutStatements(t *testing.T) {
	// No client: the driver answers without asking the endpoint.
	d := &Driver{}
	for _, script := range []string{"", "  \n\t", "; ;", "-- only a comment"} {
		log := &recorder{}
		affected, err := d.Execute(context.Background(), script, driver.ExecuteOptions{Log: log})
		require.NoError(t, err, "script %q", script)
		require.Zero(t, affected)
		require.Empty(t, log.commands, "script %q", script)
	}
}

// cleanupWrites deletes the rows a test left in the writes table, through the
// driver, so the tape holds the cleanup and a second recording finds the table
// as the first did.
func cleanupWrites(t *testing.T, d *Driver, ids ...string) {
	t.Helper()
	for _, id := range ids {
		_, err := d.Execute(context.Background(), `DELETE FROM "`+testcontainer.Writes+`" WHERE id = '`+id+`'`, driver.ExecuteOptions{})
		require.NoError(t, err)
	}
}
