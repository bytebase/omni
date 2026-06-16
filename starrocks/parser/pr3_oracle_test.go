package parser

import "testing"

// TestPR3Parity is the container-grounded conformance matrix for the PR3 DML
// features: INSERT divergences (OVERWRITE no-TABLE, BY NAME, FILES target),
// CTE-prefixed DELETE/UPDATE, and the async-MV REFRESH ASYNC EVERY schedule.
// Each construct is asserted accepted by BOTH the omni parser and StarRocks 3.4,
// and each sibling-arm negative rejected by both. Short-gated via startStarRocks.
//
// Reject probes fail INSIDE the construct (the parser does not enforce
// end-of-input, so trailing-garbage rejects would be false-accepted).
func TestPR3Parity(t *testing.T) {
	c := startStarRocks(t)

	cases := []struct {
		name       string
		sql        string
		wantAccept bool
	}{
		// INSERT divergences.
		{"insert_overwrite_no_table", "INSERT OVERWRITE t SELECT a FROM src", true},
		{"insert_by_name", "INSERT INTO t BY NAME SELECT a, b FROM s", true},
		{"insert_into_files", "INSERT INTO FILES('path'='s3://x', 'format'='parquet') SELECT a FROM s", true},
		{"insert_overwrite_files", "INSERT OVERWRITE FILES('path'='s3://x', 'format'='parquet') SELECT a FROM s", true},
		{"insert_by_without_name", "INSERT INTO t BY SELECT a FROM s", false},
		{"insert_files_no_value", "INSERT INTO FILES('path') SELECT a FROM s", false},
		{"insert_by_name_with_label", "INSERT INTO t BY NAME WITH LABEL lbl SELECT a, b FROM s", true}, // any modifier order
		{"insert_files_with_label", "INSERT INTO FILES('path'='s3://x') WITH LABEL lbl SELECT a FROM s", true},
		{"insert_by_name_collist", "INSERT INTO t BY NAME (a, b) SELECT a, b FROM s", false}, // mutually exclusive

		// CTE-prefixed DELETE / UPDATE (StarRocks allows WITH before both; not INSERT).
		{"cte_delete", "WITH c AS (SELECT id FROM s) DELETE FROM t WHERE id IN (SELECT id FROM c)", true},
		{"cte_update", "WITH c AS (SELECT id FROM s) UPDATE t SET v = 1 WHERE id IN (SELECT id FROM c)", true},
		{"cte_insert_rejected", "WITH c AS (SELECT 1) INSERT INTO t SELECT x FROM c", false},

		// Async-MV REFRESH ASYNC scheduling.
		{"mv_refresh_async_every", "CREATE MATERIALIZED VIEW mv REFRESH ASYNC EVERY (INTERVAL 1 DAY) AS SELECT a FROM t", true},
		{"mv_refresh_async_start_every", "CREATE MATERIALIZED VIEW mv REFRESH ASYNC START('2024-12-01 20:30:00') EVERY (INTERVAL 1 DAY) AS SELECT a FROM t", true},
		{"mv_refresh_immediate_async", "CREATE MATERIALIZED VIEW mv REFRESH IMMEDIATE ASYNC EVERY (INTERVAL 5 MINUTE) AS SELECT a FROM t", true},
		{"mv_refresh_async_plain", "CREATE MATERIALIZED VIEW mv REFRESH ASYNC AS SELECT a FROM t", true},
		{"mv_refresh_every_no_parens", "CREATE MATERIALIZED VIEW mv REFRESH ASYNC EVERY INTERVAL 1 DAY AS SELECT a FROM t", false},
		{"mv_refresh_every_no_interval_kw", "CREATE MATERIALIZED VIEW mv REFRESH ASYNC EVERY (1 DAY) AS SELECT a FROM t", false},
		{"mv_refresh_manual_with_every", "CREATE MATERIALIZED VIEW mv REFRESH MANUAL EVERY (INTERVAL 1 DAY) AS SELECT a FROM t", false},         // tail is ASYNC-only
		{"mv_refresh_async_start_no_every", "CREATE MATERIALIZED VIEW mv REFRESH ASYNC START('2024-12-01 20:30:00') AS SELECT a FROM t", false}, // EVERY mandatory

		// Known accepted divergences (omni superset, left as-is per skip-prune):
		//   - INSERT OVERWRITE TABLE t ...: omni keeps the doris form (26 inherited
		//     tests), but StarRocks 3.4 has no TABLE arm and rejects it — parse-only,
		//     not a parity case.
		//   - UPDATE ... FROM c ...: omni over-accepts the doris UPDATE…FROM form
		//     that StarRocks has no grammar for. Pre-existing, not from PR3.
		//   - WITH RECURSIVE ... {DELETE|UPDATE|SELECT}: omni tolerates RECURSIVE
		//     (StarRocks 3.4 has no RECURSIVE) — pre-existing for SELECT (an inherited
		//     test relies on it); PR3 extends the same tolerance to DELETE/UPDATE.
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c.assertParity(t, tc.sql, tc.wantAccept)
		})
	}
}
