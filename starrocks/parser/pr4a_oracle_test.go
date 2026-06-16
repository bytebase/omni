package parser

import "testing"

// TestPR4aParity is the container-grounded conformance matrix for the PR4a
// FROM-side query features: the VALUES inline-table constructor and (added in
// the LATERAL task) LATERAL table functions. Each StarRocks construct is
// asserted accepted by BOTH the omni parser and StarRocks 3.4, and each
// sibling-arm negative rejected by both. Short-gated via startStarRocks.
//
// Reject probes are chosen so the failure occurs INSIDE the construct: the omni
// parser does not enforce trailing-EOF, so trailing-garbage rejects (e.g. a
// column-alias list with no alias name) would be false-accepted and are
// deliberately omitted.
func TestPR4aParity(t *testing.T) {
	c := startStarRocks(t)

	cases := []struct {
		name       string
		sql        string
		wantAccept bool
	}{
		// VALUES inline-table — accept probes.
		{
			"values_alias_collist",
			"SELECT * FROM (VALUES (1,'a'),(2,'b')) AS v(id,name)",
			true,
		},
		{
			"values_alias_no_collist",
			"SELECT * FROM (VALUES (1,'a'),(2,'b')) v",
			true,
		},
		{
			"values_no_alias",
			"SELECT * FROM (VALUES (1,'a'),(2,'b'))",
			true,
		},
		{
			"values_single_column",
			"SELECT id FROM (VALUES (1),(2),(3)) AS v(id)",
			true,
		},
		{
			"values_row_expressions",
			"SELECT * FROM (VALUES (1+1, concat('a','b'))) AS v(x,y)",
			true,
		},
		{
			"values_joined",
			"SELECT * FROM t JOIN (VALUES (1),(2)) AS v(id) ON t.id = v.id",
			true,
		},
		// VALUES inline-table — sibling-arm negatives (fail inside the construct).
		{
			"values_default_in_row", // DEFAULT is INSERT-only, not a FROM-VALUES expr
			"SELECT * FROM (VALUES (1, DEFAULT)) AS v(a,b)",
			false,
		},
		{
			"values_empty", // at least one row constructor is required
			"SELECT * FROM (VALUES) AS v",
			false,
		},
		{
			"values_missing_parens", // the wrapping parens are mandatory
			"SELECT * FROM VALUES (1),(2)",
			false,
		},
		{
			"values_bare_toplevel", // VALUES is not a top-level statement in StarRocks
			"VALUES (1,'a'),(2,'b')",
			false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c.assertParity(t, tc.sql, tc.wantAccept)
		})
	}
}
