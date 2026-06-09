package parser

import "testing"

// TestPR2aParity is the container-grounded conformance matrix for the PR2a
// CREATE TABLE divergences (#10–13): each StarRocks construct is asserted to be
// accepted by BOTH the omni parser and StarRocks 3.4 (and the sibling-arm
// negative rejected by both). Short-gated via startStarRocks.
func TestPR2aParity(t *testing.T) {
	c := startStarRocks(t)

	cases := []struct {
		name       string
		sql        string
		wantAccept bool
	}{
		{
			"percentile_agg",
			"CREATE TABLE agg_pct (k INT, p PERCENTILE PERCENTILE_UNION) AGGREGATE KEY(k) DISTRIBUTED BY HASH(k) BUCKETS 1",
			true,
		},
		{
			"batch_partition_date",
			"CREATE TABLE bp (dt DATE, k INT) DUPLICATE KEY(dt) PARTITION BY RANGE(dt) (START ('2024-01-01') END ('2024-04-01') EVERY (INTERVAL 1 MONTH)) DISTRIBUTED BY HASH(k)",
			true,
		},
		{
			"batch_partition_numeric",
			"CREATE TABLE bpn (id INT, v BIGINT) DUPLICATE KEY(id) PARTITION BY RANGE(id) (START ('1') END ('100') EVERY (10)) DISTRIBUTED BY HASH(id)",
			true,
		},
		{
			"batch_partition_requires_every", // sibling-arm negative
			"CREATE TABLE bpno (dt DATE) DUPLICATE KEY(dt) PARTITION BY RANGE(dt) (START ('2024-01-01') END ('2024-04-01')) DISTRIBUTED BY HASH(dt)",
			false,
		},
		{
			"generated_column_as_parens",
			"CREATE TABLE gc (id BIGINT, price DECIMAL(18,2), qty INT, total DECIMAL(18,2) AS (price*qty)) DUPLICATE KEY(id) DISTRIBUTED BY HASH(id)",
			true,
		},
		{
			"generated_column_as_unparenthesized", // #288 P2: AS <expr> without outer parens
			"CREATE TABLE g1 (a INT, b INT, c BIGINT AS abs(a)) DUPLICATE KEY(a) DISTRIBUTED BY HASH(a) BUCKETS 1",
			true,
		},
		{
			"default_expr",
			"CREATE TABLE dft (id BIGINT, created DATETIME DEFAULT CURRENT_TIMESTAMP, token VARCHAR(64) DEFAULT (uuid())) PRIMARY KEY(id) DISTRIBUTED BY HASH(id)",
			true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c.assertParity(t, tc.sql, tc.wantAccept)
		})
	}
}
