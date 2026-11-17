package driver

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"

	gomssqldb "github.com/microsoft/go-mssqldb"

	"github.com/bytebase/omni/driver"
)

// errorLogger keeps the error each command answered with, which is what a
// task-run log shows beside the command.
type errorLogger struct {
	errs []string
}

func (*errorLogger) CommandExecute(driver.Command) {}

func (l *errorLogger) CommandResponse(_ int64, _ []int64, err error) {
	if err != nil {
		l.errs = append(l.errs, err.Error())
	}
}

func (*errorLogger) TransactionControl(driver.TransactionOp, error) {}

func (*errorLogger) Retry(error, int, int) {}

// TestExecuteNamesTheFailingStatementOfABatch holds a batch's error to the
// statement the server reported, in the error and in the step the log reports.
func TestExecuteNamesTheFailingStatementOfABatch(t *testing.T) {
	ctx := context.Background()
	d := openTestDriver(t)
	execCorpus(t, d.db, "CREATE TABLE b_items (id INT PRIMARY KEY, note NVARCHAR(64));")
	if _, err := d.Execute(ctx, "INSERT INTO b_items VALUES (1, N'a');", driver.ExecuteOptions{Transaction: driver.Transactional}); err != nil {
		t.Fatal(err)
	}

	const batch = "INSERT INTO b_items VALUES (2, N'b');\nINSERT INTO b_items VALUES (1, N'duplicate');"
	log := &errorLogger{}
	_, err := d.Execute(ctx, batch, driver.ExecuteOptions{Transaction: driver.Transactional, Log: log})
	if err == nil {
		t.Fatal("expected a duplicate key error")
	}
	// The statement the server reported is named, and the one before it is not.
	failing := `"INSERT INTO b_items VALUES (1, N'duplicate');"`
	if !strings.Contains(err.Error(), "failed to execute statement "+failing+": ") {
		t.Fatalf("error = %v, want failed to execute statement %s", err, failing)
	}
	if strings.Contains(err.Error(), "VALUES (2, N'b')") {
		t.Fatalf("error = %v, want no part of the batch", err)
	}
	if len(log.errs) != 1 || !strings.Contains(log.errs[0], failing) {
		t.Fatalf("logged errors = %v, want %s named", log.errs, failing)
	}
	// The client's own error stays reachable.
	var server gomssqldb.Error
	if !errors.As(err, &server) || server.Number != 2627 {
		t.Fatalf("error = %v, want SQL Server's 2627", err)
	}
}

// TestExecuteKeepsTheBatchWordingWithoutALine holds a batch whose error names
// no line to the batch's wording: there is no statement to point at.
func TestExecuteKeepsTheBatchWordingWithoutALine(t *testing.T) {
	noLine := gomssqldb.Error{Number: 50000, Class: 16, State: 1, Message: "statement failed"}
	batch := driver.Command{Text: "INSERT INTO b_items VALUES (3, N'c');\nINSERT INTO b_items VALUES (4, N'd');"}
	failing := func(context.Context, string, ...any) (sql.Result, error) { return nil, noLine }

	_, err := executeStatement(context.Background(), failing, batch)
	if err == nil {
		t.Fatal("expected an error")
	}
	if err.Error() != "Msg 50000, Level 16, State 1, Server , Line 0\nstatement failed" {
		t.Fatalf("error = %v, want the server's messages alone", err)
	}
	if strings.Contains(err.Error(), "INSERT INTO") {
		t.Fatalf("error = %v, want no statement named", err)
	}
	// An error that is not the server's names no line either.
	_, err = executeStatement(context.Background(), func(context.Context, string, ...any) (sql.Result, error) {
		return nil, errors.New("connection lost")
	}, batch)
	if err == nil || err.Error() != "failed to execute statement: connection lost" {
		t.Fatalf("error = %v, want the batch's wording", err)
	}
}

// TestExecuteAutoCommitKeepsSessionState pins that an auto-commit script keeps
// one connection: what a batch sets up for the next one is still there.
func TestExecuteAutoCommitKeepsSessionState(t *testing.T) {
	ctx := context.Background()
	d := openTestDriver(t)
	execCorpus(t, d.db, "CREATE TABLE s_items (id INT PRIMARY KEY);")

	if _, err := d.Execute(ctx, "CREATE TABLE #scratch (id INT);\nINSERT INTO #scratch VALUES (1), (2);\nGO\nINSERT INTO s_items SELECT id FROM #scratch;", driver.ExecuteOptions{
		Transaction: driver.AutoCommit,
	}); err != nil {
		t.Fatalf("a temporary table did not outlive its batch: %v", err)
	}
	counted, err := d.Session(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer counted.Close()
	results, err := counted.Query(ctx, "SELECT COUNT(*) FROM s_items;", driver.QueryOptions{})
	if err != nil || len(results) != 1 || results[0].Rows[0][0].Int64 != 2 {
		t.Fatalf("rows = %+v, %v", results, err)
	}
}
