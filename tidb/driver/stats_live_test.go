package driver

import (
	"context"
	"strings"
	"testing"
)

// TestCountAffectedRowsEstimate pins the planner estimate the review rules read
// as a statement's row count: it follows the operator that feeds the DML, so a
// filter over a hundred matching rows estimates whole-row counts far below the
// thousand-row scan beneath it, and a statement that changes no rows is refused
// before it reaches the server.
func TestCountAffectedRowsEstimate(t *testing.T) {
	ctx := context.Background()
	d := openTestDriver(t, "omni_estimate_test")
	for _, statement := range []string{
		"CREATE TABLE t (id INT PRIMARY KEY, c INT)",
		"CREATE TABLE archive (id INT PRIMARY KEY, c INT)",
		`INSERT INTO t
			WITH RECURSIVE seq (n) AS (SELECT 1 UNION ALL SELECT n + 1 FROM seq WHERE n < 1000)
			SELECT n, n % 10 FROM seq`,
		// By default ANALYZE collects only the columns earlier queries filtered
		// on, which would leave c without statistics.
		"ANALYZE TABLE t ALL COLUMNS",
	} {
		if _, err := d.db.ExecContext(ctx, statement); err != nil {
			t.Fatalf("%s: %v", statement, err)
		}
	}

	for _, statement := range []string{
		"UPDATE t SET c = 0 WHERE c = 1",
		"DELETE FROM t WHERE c = 1",
		"INSERT INTO archive SELECT * FROM t WHERE c = 1",
	} {
		got, err := d.CountAffectedRows(ctx, statement)
		if err != nil {
			t.Fatalf("CountAffectedRows(%q) = %v", statement, err)
		}
		// 100 of the 1000 rows match c = 1, while the scan beneath the filter
		// estimates all 1000.
		if got <= 0 || got >= 1000 {
			t.Errorf("CountAffectedRows(%q) = %d, want an estimate between 1 and 999", statement, got)
		}
	}

	for _, statement := range []string{
		"SELECT * FROM t",
		"DELETE FROM t WHERE c = 1; SELECT 1",
	} {
		if _, err := d.CountAffectedRows(ctx, statement); err == nil || !strings.Contains(err.Error(), "expected") {
			t.Errorf("CountAffectedRows(%q) err = %v, want a refusal", statement, err)
		}
	}
}
