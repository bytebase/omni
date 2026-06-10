package parser

import "testing"

// TestPR2bParity is the container-grounded conformance matrix for the PR2b
// CREATE/ALTER TABLE divergences (#14–17): each StarRocks construct is asserted
// accepted by BOTH the omni parser and StarRocks 3.4, plus regressions proving
// the new branches don't break the pre-existing forms. Short-gated via
// startStarRocks.
func TestPR2bParity(t *testing.T) {
	c := startStarRocks(t)

	cases := []struct {
		name       string
		sql        string
		wantAccept bool
	}{
		{
			"inline_index_gin_props",
			"CREATE TABLE ix (id BIGINT, name VARCHAR(64), tags VARCHAR(255), INDEX idx_name (name) USING BITMAP COMMENT 'x', INDEX idx_tags (tags) USING GIN ('parser'='english')) DUPLICATE KEY(id) DISTRIBUTED BY HASH(id)",
			true,
		},
		{
			"rollup_from",
			"CREATE TABLE rl (k1 INT, k2 INT, v BIGINT SUM) AGGREGATE KEY(k1,k2) DISTRIBUTED BY HASH(k1) ROLLUP (r1 (k1,v), r2 (k2,v) FROM r1)",
			true,
		},
		{
			"ctas_name_list",
			"CREATE TABLE ck (new_id, new_name) DISTRIBUTED BY HASH(new_id) BUCKETS 8 AS SELECT k AS new_id, v AS new_name FROM smoke_t",
			true,
		},
		{
			"ctas_name_list_with_index",
			"CREATE TABLE ck2 (a, b, INDEX idx_b (b) USING BITMAP) DISTRIBUTED BY HASH(a) AS SELECT k AS a, v AS b FROM smoke_t",
			true,
		},
		{
			"full_defs_regression", // full column-defs still parse, not name-list
			"CREATE TABLE reg (a INT, b VARCHAR(20)) DUPLICATE KEY(a) DISTRIBUTED BY HASH(a) BUCKETS 1",
			true,
		},
		{
			"alter_add_field",
			"ALTER TABLE t MODIFY COLUMN s ADD FIELD f4 DOUBLE AFTER f2",
			true,
		},
		{
			"alter_drop_field",
			"ALTER TABLE t MODIFY COLUMN s DROP FIELD nested",
			true,
		},
		{
			"alter_modify_column_regression", // plain MODIFY COLUMN col TYPE still parses
			"ALTER TABLE t MODIFY COLUMN c BIGINT",
			true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c.assertParity(t, tc.sql, tc.wantAccept)
		})
	}
}
