//go:build oracle

package parser

import "testing"

// TestParenOracleFromJoined transcribes SCENARIOS-pg-paren-dispatch.md §2.4
// (FROM-clause joined_table-shape corpus) into oracle probes. Each case
// asserts the OmniStatus omni must emit to be PG-17 aligned for a paren-
// wrapped joined_table FROM item and for its nested/aliased/LATERAL
// variants.
//
// Coverage outline:
//   - wrapping depth: double-, triple-wrapped paren joined_table
//   - nesting topology: BYT-9315 left-nested, right-side paren-join,
//     both-sides paren-join, deeply nested left/right
//   - mixed / chained joins: CROSS+JOIN mix, chained LEFT/RIGHT outer
//   - outer-alias variants: AS jt, AS jt(col1), alias-then-outer-join,
//     natural-then-using
//   - USING variants exercising gram.y `join_type` branch: FULL OUTER
//     USING, multi-column USING
//   - pg_get_viewdef() round-trip shape: double-paren ON-expr
//   - LATERAL-inside-joined_table: LATERAL subquery as right operand at
//     first and second join levels
//
// Fixture tables T/U/V/W and function probe_f() are created by
// StartParenOracle (see paren_oracle_test.go).
func TestParenOracleFromJoined(t *testing.T) {
	o := StartParenOracle(t)

	cases := []struct {
		name     string
		sql      string
		expected OmniStatus
	}{
		// --- wrapping depth ---
		{
			name:     "double-wrapped joined_table",
			sql:      `SELECT * FROM ((T JOIN U ON TRUE))`,
			expected: OmniJoinedTable,
		},
		{
			name:     "triple-wrapped joined_table",
			sql:      `SELECT * FROM (((T JOIN U ON TRUE)))`,
			expected: OmniJoinedTable,
		},

		// --- nesting topology ---
		{
			name:     "BYT-9315 left-nested paren-join",
			sql:      `SELECT * FROM ((T JOIN U ON TRUE) JOIN V ON TRUE)`,
			expected: OmniJoinedTable,
		},
		{
			name:     "right-side paren-join",
			sql:      `SELECT * FROM (T JOIN (U JOIN V ON TRUE) ON TRUE)`,
			expected: OmniJoinedTable,
		},
		{
			name:     "both-sides paren-join",
			sql:      `SELECT * FROM ((T JOIN U ON TRUE) JOIN (V JOIN W ON TRUE) ON TRUE)`,
			expected: OmniJoinedTable,
		},
		{
			name:     "deeply nested left",
			sql:      `SELECT * FROM (((T JOIN U ON TRUE) JOIN V ON TRUE) JOIN W ON TRUE)`,
			expected: OmniJoinedTable,
		},
		{
			name:     "deeply nested right",
			sql:      `SELECT * FROM (T JOIN (U JOIN (V JOIN W ON TRUE) ON TRUE) ON TRUE)`,
			expected: OmniJoinedTable,
		},

		// --- mixed / chained joins ---
		{
			name:     "mixed CROSS + JOIN",
			sql:      `SELECT * FROM (T CROSS JOIN U JOIN V ON TRUE)`,
			expected: OmniJoinedTable,
		},
		{
			name:     "chained LEFT and RIGHT outer joins",
			sql:      `SELECT * FROM (T LEFT JOIN U ON T.a = U.a RIGHT JOIN V ON U.b = V.b)`,
			expected: OmniJoinedTable,
		},

		// --- outer alias on paren-joined ---
		// Note: `(joined_table) AS jt [JOIN ...]` — the top-level FROM
		// item is still a JoinExpr because PG grammar attaches alias_clause
		// directly to joined_table (not wrapping it in a RangeSubselect).
		{
			name:     "outer alias on paren-joined",
			sql:      `SELECT * FROM ((T JOIN U ON TRUE) JOIN V ON TRUE) AS jt`,
			expected: OmniJoinedTable,
		},
		{
			name:     "outer column-list alias on paren-joined",
			sql:      `SELECT * FROM ((T JOIN U ON TRUE) JOIN V ON TRUE) AS jt(col1)`,
			expected: OmniJoinedTable,
		},
		{
			name:     "alias then outer join",
			sql:      `SELECT * FROM (T JOIN U ON TRUE) AS jt JOIN V ON TRUE`,
			expected: OmniJoinedTable,
		},
		{
			name:     "natural then using",
			sql:      `SELECT * FROM (T NATURAL JOIN U) JOIN V USING (a)`,
			expected: OmniJoinedTable,
		},

		// --- USING variants exercising join_type ---
		{
			name:     "FULL OUTER JOIN USING single column",
			sql:      `SELECT * FROM (T FULL OUTER JOIN U USING (a))`,
			expected: OmniJoinedTable,
		},
		{
			name:     "LEFT OUTER JOIN USING multi-column",
			sql:      `SELECT * FROM (T LEFT OUTER JOIN U USING (a, b))`,
			expected: OmniJoinedTable,
		},

		// --- pg_get_viewdef() round-trip shape ---
		// PG's view deparser emits ON-expressions wrapped in an extra pair
		// of parens: `ON ((T.a = U.a))`. The whole joined_table itself is
		// also paren-wrapped. This is the exact shape our parser sees
		// after a CREATE VIEW round-trip.
		{
			name:     "pg_get_viewdef shape",
			sql:      `SELECT * FROM ((T JOIN U ON ((T.a = U.a))) JOIN V ON ((U.b = V.b)))`,
			expected: OmniJoinedTable,
		},

		// --- LATERAL inside joined_table ---
		// LATERAL is legal as the right operand of a JOIN when the right
		// operand is a subquery/function. The outer paren wraps the
		// joined_table, not a subquery.
		{
			name:     "LATERAL subquery as first JOIN right operand",
			sql:      `SELECT * FROM (T JOIN LATERAL (SELECT U.a FROM U WHERE U.x = T.x) v ON TRUE)`,
			expected: OmniJoinedTable,
		},
		{
			name:     "LATERAL subquery as second JOIN right operand",
			sql:      `SELECT * FROM ((T JOIN U ON T.a = U.a) JOIN LATERAL (SELECT 1) l ON TRUE)`,
			expected: OmniJoinedTable,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			assertParenParity(t, o, tc.sql, tc.expected)
		})
	}
}
