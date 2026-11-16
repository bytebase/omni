package driver

import (
	"context"
	"testing"

	"github.com/bytebase/omni/driver"
)

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
