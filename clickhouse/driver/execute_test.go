package driver

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/bytebase/omni/driver"
)

type stepLogger struct {
	steps []string
}

func (l *stepLogger) CommandExecute(command driver.Command) {
	l.steps = append(l.steps, fmt.Sprintf("exec[%d:%d] %s", command.ByteStart, command.ByteEnd, command.Text))
}

func (l *stepLogger) CommandResponse(affected int64, perResult []int64, err error) {
	if err != nil {
		l.steps = append(l.steps, "resp:err")
		return
	}
	l.steps = append(l.steps, fmt.Sprintf("resp:%d:%v", affected, perResult))
}

func (l *stepLogger) TransactionControl(op driver.TransactionOp, err error) {
	if err != nil {
		l.steps = append(l.steps, fmt.Sprintf("tx:%s:err", op))
		return
	}
	l.steps = append(l.steps, fmt.Sprintf("tx:%s", op))
}

func (l *stepLogger) Retry(_ error, attempt, max int) {
	l.steps = append(l.steps, fmt.Sprintf("retry:%d/%d", attempt, max))
}

const executeSetup = `
CREATE TABLE e (id Int32, note String) ENGINE = Memory;
`

func setupExecute(t *testing.T) *Driver {
	t.Helper()
	d := openTestDriver(t)
	execCorpus(t, d.db, executeSetup)
	return d
}

func TestExecute(t *testing.T) {
	d := setupExecute(t)
	logger := &stepLogger{}
	affected, err := d.Execute(context.Background(), "INSERT INTO e VALUES (1, 'a');\nINSERT INTO e VALUES (2, 'b')", driver.ExecuteOptions{
		Transaction: driver.Transactional,
		Log:         logger,
	})
	if err != nil {
		t.Fatal(err)
	}
	if affected != 0 {
		t.Errorf("affected = %d, want the client's 0", affected)
	}
	if got, want := logger.steps, []string{
		"tx:BEGIN",
		"exec[0:30] INSERT INTO e VALUES (1, 'a');",
		"resp:0:[]",
		"exec[30:60] \nINSERT INTO e VALUES (2, 'b')",
		"resp:0:[]",
		"tx:COMMIT",
	}; !equalStrings(got, want) {
		t.Errorf("steps = %q, want %q", got, want)
	}
	if got := countRows(t, d, "e"); got != 2 {
		t.Errorf("rows = %d, want 2", got)
	}
}

func TestExecuteAutoCommit(t *testing.T) {
	d := setupExecute(t)
	logger := &stepLogger{}
	affected, err := d.Execute(context.Background(), "INSERT INTO e VALUES (1, 'a');\nINSERT INTO e VALUES (2, 'b')", driver.ExecuteOptions{
		Transaction: driver.AutoCommit,
		Log:         logger,
	})
	if err != nil {
		t.Fatal(err)
	}
	if affected != 0 {
		t.Errorf("affected = %d, want the client's 0", affected)
	}
	if got, want := logger.steps, []string{
		"exec[0:30] INSERT INTO e VALUES (1, 'a');",
		"resp:0:[]",
		"exec[30:60] \nINSERT INTO e VALUES (2, 'b')",
		"resp:0:[]",
	}; !equalStrings(got, want) {
		t.Errorf("steps = %q, want %q", got, want)
	}
}

func TestExecuteStopsAtTheFailingStatement(t *testing.T) {
	d := setupExecute(t)
	affected, err := d.Execute(context.Background(), "INSERT INTO e VALUES (1, 'a');\nINSERT INTO no_such_table VALUES (2, 'b');\nINSERT INTO e VALUES (3, 'c')", driver.ExecuteOptions{
		Transaction: driver.AutoCommit,
	})
	if err == nil || !strings.Contains(err.Error(), "no_such_table") {
		t.Fatalf("err = %v", err)
	}
	if affected != 0 {
		t.Errorf("affected = %d, want 0", affected)
	}
	if got := countRows(t, d, "e"); got != 1 {
		t.Errorf("rows = %d, want the first insert to remain", got)
	}
}

func TestExecuteTransactionFailureReturnsZero(t *testing.T) {
	d := setupExecute(t)
	logger := &stepLogger{}
	affected, err := d.Execute(context.Background(), "INSERT INTO e VALUES (1, 'a');\nINSERT INTO no_such_table VALUES (2, 'b')", driver.ExecuteOptions{
		Transaction: driver.Transactional,
		Log:         logger,
	})
	if err == nil || !strings.Contains(err.Error(), "no_such_table") {
		t.Fatalf("err = %v", err)
	}
	if affected != 0 {
		t.Errorf("affected = %d, want 0", affected)
	}
	if got, want := logger.steps, []string{
		"tx:BEGIN",
		"exec[0:30] INSERT INTO e VALUES (1, 'a');",
		"resp:0:[]",
		"exec[30:72] \nINSERT INTO no_such_table VALUES (2, 'b')",
		"resp:err",
		"tx:ROLLBACK",
	}; !equalStrings(got, want) {
		t.Errorf("steps = %q, want %q", got, want)
	}
	// clickhouse-go's begin is a no-op, so the first insert is still there.
	if got := countRows(t, d, "e"); got != 1 {
		t.Errorf("rows = %d, want the first insert to remain", got)
	}
}

func countRows(t *testing.T, d *Driver, table string) int64 {
	t.Helper()
	results := query(t, d, "SELECT count() FROM "+table, driver.QueryOptions{})
	if len(results) != 1 || results[0].Err != nil {
		t.Fatalf("count: %+v", results)
	}
	return int64(results[0].Rows[0][0].Uint64)
}

func equalStrings(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}
