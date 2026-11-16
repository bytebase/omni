package driver

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/driver"
)

// TestExecute runs a script and reads its result back. Cassandra reports no
// affected rows, so the count is 0 whatever the script did, and neither the
// transaction mode nor the isolation level changes anything: each statement
// commits itself.
func TestExecute(t *testing.T) {
	d := openDriver(t)
	affected, err := d.Execute(context.Background(), `
		CREATE TABLE exec_items (id int PRIMARY KEY, note text);
		INSERT INTO exec_items (id, note) VALUES (1, 'a');
		INSERT INTO exec_items (id, note) VALUES (2, 'b');
		UPDATE exec_items SET note = 'z' WHERE id = 1;
		DELETE FROM exec_items WHERE id = 2;`, driver.ExecuteOptions{
		Transaction: driver.Transactional,
		Isolation:   driver.Serializable,
	})
	require.NoError(t, err)
	require.Zero(t, affected)

	sess := openSession(t, d)
	results, err := sess.Query(context.Background(), "SELECT id, note FROM exec_items", driver.QueryOptions{})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.Len(t, results[0].Rows, 1)
	require.Equal(t, int64(1), results[0].Rows[0][0].Int64)
	require.Equal(t, "z", results[0].Rows[0][1].String)
}

// TestExecuteStopsAtAFailingStatement holds a script's failure to its own
// error, and the statements before it stay: Cassandra has no transaction to
// roll them back.
func TestExecuteStopsAtAFailingStatement(t *testing.T) {
	d := openDriver(t)
	_, err := d.Execute(context.Background(), `
		CREATE TABLE kept (id int PRIMARY KEY);
		INSERT INTO kept (id) VALUES (1);
		INSERT INTO missing (id) VALUES (2);`, driver.ExecuteOptions{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to execute")

	sess := openSession(t, d)
	results, err := sess.Query(context.Background(), "SELECT id FROM kept", driver.QueryOptions{})
	require.NoError(t, err)
	require.Len(t, results[0].Rows, 1)
}

type stepLogger struct {
	steps []string
}

func (l *stepLogger) CommandExecute(cmd driver.Command) {
	l.steps = append(l.steps, fmt.Sprintf("exec[%d:%d]%s", cmd.ByteStart, cmd.ByteEnd, cmd.Text))
}

func (l *stepLogger) CommandResponse(affected int64, perResult []int64, err error) {
	if err != nil {
		l.steps = append(l.steps, "resp:err")
		return
	}
	l.steps = append(l.steps, fmt.Sprintf("resp:%d", affected))
}

func (*stepLogger) TransactionControl(driver.TransactionOp, error) {}

func (*stepLogger) Retry(error, int, int) {}

// TestExecuteLogsItsSteps holds Execute to the contract's logger: one command
// with its range and one response for each statement it runs.
func TestExecuteLogsItsSteps(t *testing.T) {
	d := openDriver(t)
	log := &stepLogger{}
	_, err := d.Execute(context.Background(), "CREATE TABLE logged (id int PRIMARY KEY); INSERT INTO logged (id) VALUES (1)", driver.ExecuteOptions{Log: log})
	require.NoError(t, err)
	require.Len(t, log.steps, 4)
	require.Contains(t, log.steps[0], "CREATE TABLE logged")
	require.Equal(t, "resp:0", log.steps[1])
	require.Contains(t, log.steps[2], "INSERT INTO logged")
	require.Equal(t, "resp:0", log.steps[3])

	// A failing statement reports its response with the error and ends the
	// script.
	log = &stepLogger{}
	_, err = d.Execute(context.Background(), "INSERT INTO logged (id) VALUES (2); INSERT INTO missing (id) VALUES (3); INSERT INTO logged (id) VALUES (4)", driver.ExecuteOptions{Log: log})
	require.Error(t, err)
	require.Len(t, log.steps, 4)
	require.Equal(t, "resp:0", log.steps[1])
	require.Equal(t, "resp:err", log.steps[3])
}

// TestExecuteBatch holds a batch to one request: its inner statements are
// separated by semicolons that do not end it, and the server applies them
// together.
func TestExecuteBatch(t *testing.T) {
	d := openDriver(t)
	_, err := d.Execute(context.Background(), `
		CREATE TABLE batched (id int PRIMARY KEY, n int);
		BEGIN BATCH
			INSERT INTO batched (id, n) VALUES (1, 10);
			INSERT INTO batched (id, n) VALUES (2, 20);
		APPLY BATCH;`, driver.ExecuteOptions{})
	require.NoError(t, err)

	sess := openSession(t, d)
	results, err := sess.Query(context.Background(), "SELECT id, n FROM batched", driver.QueryOptions{})
	require.NoError(t, err)
	require.Len(t, results[0].Rows, 2)
}

// TestExecuteSkipsACommentOnlyStatement holds a script that ends in a comment to
// its statements: the comment is not sent as a query, so the script succeeds.
func TestExecuteSkipsACommentOnlyStatement(t *testing.T) {
	d := openDriver(t)
	affected, err := d.Execute(context.Background(), `
		CREATE TABLE commented (id int PRIMARY KEY);
		INSERT INTO commented (id) VALUES (1);
		-- done
		/* and a block */`, driver.ExecuteOptions{})
	require.NoError(t, err)
	require.Zero(t, affected)

	sess := openSession(t, d)
	results, err := sess.Query(context.Background(), "SELECT id FROM commented", driver.QueryOptions{})
	require.NoError(t, err)
	require.Len(t, results[0].Rows, 1)
}

// TestExecuteCorpusRecordsOnce is the corpus every other test records, run
// through Execute: a script of one statement per line is one request per
// statement on the tape.
func TestExecuteCorpus(t *testing.T) {
	d := openDriver(t)
	executeCorpus(t, d, mustRead(t, filepath.Join("testdata", "corpus.cql")))
	sess := openSession(t, d)
	results, err := sess.Query(context.Background(), "SELECT count(*) FROM omni_rows", driver.QueryOptions{})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.Equal(t, int64(10), results[0].Rows[0][0].Int64)
}
