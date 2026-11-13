package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestLiveProbes(t *testing.T) {
	stmts := parseStatements(t, `INSERT INTO t VALUES (1), (2), (3);
SELECT 1;
UPDATE t SET a = 1;
EXPLAIN ANALYZE DELETE FROM t;
EXPLAIN (ANALYZE false) DELETE FROM t;
WITH d AS (DELETE FROM t RETURNING *) SELECT * FROM d;
WITH d AS (UPDATE t SET a = 1 RETURNING *) INSERT INTO u SELECT * FROM d;
INSERT INTO u SELECT * FROM t;
CREATE TABLE c AS WITH d AS (DELETE FROM t RETURNING *) SELECT * FROM d;
MERGE INTO t USING s ON t.id = s.id WHEN MATCHED THEN DELETE;
CREATE TABLE x (a int)`)
	if diff := cmp.Diff([]int{0, 2, 3, 5, 6, 7, 9}, StatementDMLDryRunProbes(stmts)); diff != "" {
		t.Errorf("dry-run probes (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff([]int{2, 3, 5, 6, 8, 9}, StatementAffectedRowLimitProbes(stmts)); diff != "" {
		t.Errorf("affected-row probes (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff([]int{6, 7}, StatementInsertRowLimitProbes(stmts)); diff != "" {
		t.Errorf("insert-row probes (-want +got):\n%s", diff)
	}
	got := StatementInsertRowLimit(stmts, map[int]v1review.Probe{6: {CountError: "no ModifyTable"}, 7: {Rows: 2}}, 2)
	want := []v1review.Finding{
		{Code: v1review.RowLimitExceeded, Value: "3", Statement: 0},
		{Code: v1review.RowCountFailed, Value: "no ModifyTable", Statement: 6},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("insert-row findings (-want +got):\n%s", diff)
	}
	got = StatementDMLDryRun(stmts, map[int]v1review.Probe{3: {Error: "permission denied"}, 4: {Error: "not probed"}})
	if diff := cmp.Diff([]v1review.Finding{{Code: v1review.DryRunFailed, Value: "permission denied", Statement: 3}}, got); diff != "" {
		t.Errorf("dry-run findings (-want +got):\n%s", diff)
	}
	got = StatementAffectedRowLimit(stmts, map[int]v1review.Probe{2: {Rows: 3}, 9: {Rows: 2}}, 2)
	if diff := cmp.Diff([]v1review.Finding{{Code: v1review.RowLimitExceeded, Value: "3", Statement: 2}}, got); diff != "" {
		t.Errorf("affected-row findings (-want +got):\n%s", diff)
	}
}
