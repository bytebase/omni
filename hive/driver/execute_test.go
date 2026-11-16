package driver

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/driver"
	"github.com/bytebase/omni/hive/internal/testcontainer"
)

// TestExecute runs a change script on one connection, then reads its effect
// back through a query. Hive reports no affected-row count for either
// statement, so the total is zero.
func TestExecute(t *testing.T) {
	d := openTestCorpus(t, "omni_execute")

	affected, err := d.Execute(context.Background(), `
		CREATE TABLE exec_t (id INT, note STRING) STORED AS PARQUET;
		INSERT INTO exec_t VALUES (1, 'a'), (2, 'b');`, driver.ExecuteOptions{})
	require.NoError(t, err)
	require.Equal(t, int64(0), affected)

	sess := openSession(t, d)
	results, err := sess.Query(context.Background(), "SELECT id, note FROM exec_t", driver.QueryOptions{})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.NoError(t, results[0].Err)
	require.Len(t, results[0].Rows, 2)
	require.Equal(t, int64(1), results[0].Rows[0][0].Int64)
	require.Equal(t, "a", results[0].Rows[0][1].String)
	require.Equal(t, int64(2), results[0].Rows[1][0].Int64)
	require.Equal(t, "b", results[0].Rows[1][1].String)
}

// TestExecuteKeepsTheScriptsDatabase: a USE in the script holds for the
// statements after it, because the script runs on one connection.
func TestExecuteKeepsTheScriptsDatabase(t *testing.T) {
	d := openTestCorpus(t, "omni_execute_use")

	_, err := d.Execute(context.Background(), `
		USE omni_execute_use;
		CREATE TABLE use_t (id INT) STORED AS PARQUET;
		INSERT INTO use_t VALUES (7);`, driver.ExecuteOptions{})
	require.NoError(t, err)

	sess := openSession(t, d)
	results, err := sess.Query(context.Background(), "SELECT id FROM use_t", driver.QueryOptions{})
	require.NoError(t, err)
	require.NoError(t, results[0].Err)
	require.Len(t, results[0].Rows, 1)
	require.Equal(t, int64(7), results[0].Rows[0][0].Int64)
}

// TestExecuteStopsAtAFailingStatement returns the rows the statements before
// the failure affected, with Hive's answer to the failure.
func TestExecuteStopsAtAFailingStatement(t *testing.T) {
	d := openTestCorpus(t, "omni_execute_error")

	_, err := d.Execute(context.Background(), `
		CREATE TABLE exec_e (id INT) STORED AS PARQUET;
		INSERT INTO no_such_table VALUES (1);`, driver.ExecuteOptions{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to execute statement:")
	require.Contains(t, err.Error(), "INSERT INTO no_such_table VALUES (1)")

	var hiveErr *testcontainer.Error
	require.True(t, errors.As(err, &hiveErr), "the error is not Hive's answer: %v", err)
	require.Equal(t, 40000, hiveErr.Code)
	require.Contains(t, hiveErr.Message, "InvalidTableException: Table not found no_such_table")

	// The statements before the failure committed themselves.
	sess := openSession(t, d)
	results, err := sess.Query(context.Background(), "SHOW TABLES", driver.QueryOptions{})
	require.NoError(t, err)
	require.NoError(t, results[0].Err)
	tables := make([]string, 0, len(results[0].Rows))
	for _, row := range results[0].Rows {
		tables = append(tables, row[0].String)
	}
	require.Equal(t, []string{"events", "exec_e", "ext_users", "orders", "v_users"}, tables)
}

// TestExecuteStartsFromTheConfiguredDatabase: a session that selected another
// database leaves that choice on the pooled connection, and the script after it
// still runs in the database the driver was opened on.
func TestExecuteStartsFromTheConfiguredDatabase(t *testing.T) {
	testcontainer.Reset(t, func(ctx context.Context, db *sql.DB) error {
		for _, stmt := range []string{
			"DROP DATABASE IF EXISTS omni_execute_home CASCADE",
			"DROP DATABASE IF EXISTS omni_execute_other CASCADE",
			"CREATE DATABASE omni_execute_home",
			"CREATE DATABASE omni_execute_other",
		} {
			if _, err := db.ExecContext(ctx, stmt); err != nil {
				return err
			}
		}
		return nil
	})
	db := testcontainer.Open(t, "omni_execute_home")
	d := &Driver{db: db, databaseName: "omni_execute_home"}

	// A session moves the one pooled connection to another database.
	sess := openSession(t, d)
	results, err := sess.Query(context.Background(), "USE omni_execute_other", driver.QueryOptions{})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.Error(t, results[0].Err)
	require.NoError(t, sess.Close())

	_, err = d.Execute(context.Background(), "CREATE TABLE exec_home (id INT) STORED AS PARQUET", driver.ExecuteOptions{})
	require.NoError(t, err)

	// The table is in the configured database, not the session's.
	other := openSession(t, d)
	results, err = other.Query(context.Background(), "SHOW TABLES", driver.QueryOptions{Schema: "omni_execute_other"})
	require.NoError(t, err)
	require.NoError(t, results[0].Err)
	require.Empty(t, results[0].Rows)
	results, err = other.Query(context.Background(), "SHOW TABLES", driver.QueryOptions{Schema: "omni_execute_home"})
	require.NoError(t, err)
	require.NoError(t, results[0].Err)
	require.Len(t, results[0].Rows, 1)
	require.Equal(t, "exec_home", results[0].Rows[0][0].String)
}

// TestExecuteLogsItsSteps reports every statement and its answer to the
// logger, in order, with the statement's place in the script.
func TestExecuteLogsItsSteps(t *testing.T) {
	d := openTestCorpus(t, "omni_execute_log")
	log := &recordingLogger{}

	script := `CREATE TABLE exec_l (id INT) STORED AS PARQUET;
INSERT INTO exec_l VALUES (1);`
	_, err := d.Execute(context.Background(), script, driver.ExecuteOptions{Log: log})
	require.NoError(t, err)

	// A command carries the segment the splitter cut, terminator included.
	require.Equal(t, []string{"CREATE TABLE exec_l (id INT) STORED AS PARQUET;", "\nINSERT INTO exec_l VALUES (1);"}, log.commands)
	require.Equal(t, []string{"0 <nil>", "0 <nil>"}, log.responses)
	require.Equal(t, "CREATE TABLE exec_l (id INT) STORED AS PARQUET;", script[log.starts[0]:log.ends[0]])
	require.Equal(t, "\nINSERT INTO exec_l VALUES (1);", script[log.starts[1]:log.ends[1]])
}

// TestExecuteKeepsTheCommentAboveAStatement: a comment line is part of the
// statement it introduces, which is what the splitter sends, and a comment
// after the last statement belongs to no statement.
func TestExecuteKeepsTheCommentAboveAStatement(t *testing.T) {
	d := openTestCorpus(t, "omni_execute_comment")
	log := &recordingLogger{}

	_, err := d.Execute(context.Background(), `-- a comment
CREATE TABLE exec_c (id INT) STORED AS PARQUET;
-- another comment
`, driver.ExecuteOptions{Log: log})
	require.NoError(t, err)
	require.Equal(t, []string{"-- a comment\nCREATE TABLE exec_c (id INT) STORED AS PARQUET;"}, log.commands)
}

type recordingLogger struct {
	commands  []string
	responses []string
	starts    []int
	ends      []int
}

func (l *recordingLogger) CommandExecute(cmd driver.Command) {
	l.commands = append(l.commands, cmd.Text)
	l.starts = append(l.starts, cmd.ByteStart)
	l.ends = append(l.ends, cmd.ByteEnd)
}

func (l *recordingLogger) CommandResponse(affectedRows int64, perResult []int64, err error) {
	l.responses = append(l.responses, strconv.FormatInt(affectedRows, 10)+" "+errText(err))
}

func (l *recordingLogger) TransactionControl(driver.TransactionOp, error) {}

func (l *recordingLogger) Retry(error, int, int) {}

func errText(err error) string {
	if err == nil {
		return "<nil>"
	}
	return err.Error()
}
