package driver

import (
	"context"
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
