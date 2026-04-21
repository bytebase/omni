//go:build oracle

package parser

import "testing"

// TestParenOracleFromMixed transcribes SCENARIOS-pg-paren-dispatch.md §2.5
// (FROM-clause mixed shapes corpus) into oracle probes. Each case asserts
// the OmniStatus omni must emit to be PG-17 aligned for paren-wrapped
// joined_tables whose operands are subqueries, VALUES, TABLE, or bare
// relations — the `((subquery) JOIN relation)` family, which is the exact
// shape that motivated the Phase 0 BYT-9315 fix.
//
// Coverage outline:
//   - subquery as left-hand JOIN operand (plain JOIN, CROSS, NATURAL)
//   - both operands subquery
//   - VALUES with column-alias as left-hand operand
//   - TABLE as left-hand operand — verified ACCEPT against PG 17 oracle;
//     the SCENARIOS "check — may reject" is resolved to accept.
//   - paren-wrapped single relation as left-hand operand — verified REJECT
//     against PG 17 oracle; the SCENARIOS "verify" resolves to reject with
//     42601 at the `)` following the lone relation name.
//   - subquery as right-hand JOIN operand
//
// Fixture tables T/U/V/W are created by StartParenOracle
// (see paren_oracle_test.go).
func TestParenOracleFromMixed(t *testing.T) {
	o := StartParenOracle(t)

	cases := []struct {
		name     string
		sql      string
		expected OmniStatus
	}{
		// --- subquery as left-hand JOIN operand ---
		{
			name:     "subquery JOIN relation",
			sql:      `SELECT * FROM ((SELECT 1) x JOIN U ON TRUE)`,
			expected: OmniJoinedTable,
		},
		{
			name:     "subquery CROSS JOIN relation",
			sql:      `SELECT * FROM ((SELECT 1) x CROSS JOIN U)`,
			expected: OmniJoinedTable,
		},
		{
			name:     "subquery NATURAL JOIN relation",
			sql:      `SELECT * FROM ((SELECT 1) x NATURAL JOIN U)`,
			expected: OmniJoinedTable,
		},

		// --- both operands subquery ---
		// On-clause references the aliases of both subqueries. PG accepts
		// both operands being select_with_parens inside a joined_table.
		{
			name:     "both operands subquery with ON qual",
			sql:      `SELECT * FROM ((SELECT 1 AS a) x JOIN (SELECT 2 AS a) y ON x.a = y.a)`,
			expected: OmniJoinedTable,
		},

		// --- VALUES as left-hand operand ---
		// VALUES in parens with column-list alias is a valid
		// select_with_parens table reference; PG routes `((VALUES ...) v(a)
		// JOIN ...)` through joined_table.
		{
			name:     "VALUES with column alias JOIN relation",
			sql:      `SELECT * FROM ((VALUES (1)) v(a) JOIN U ON TRUE)`,
			expected: OmniJoinedTable,
		},

		// --- TABLE as left-hand operand ---
		// Verified against PG 17 oracle: `((TABLE T) JOIN U ON TRUE)` is
		// accepted and produces a joined_table whose left operand is a
		// select_with_parens wrapping `TABLE T` (itself a simple_select).
		// The SCENARIOS "check — may reject" note is resolved to ACCEPT.
		{
			name:     "TABLE form JOIN relation",
			sql:      `SELECT * FROM ((TABLE T) JOIN U ON TRUE)`,
			expected: OmniJoinedTable,
		},

		// --- paren-wrapped single relation as left-hand operand ---
		// Verified against PG 17 oracle: `((T) JOIN U ON TRUE)` is rejected
		// with 42601 (syntax error at or near `)` after the lone `T`). PG's
		// grammar has no production for `( table_ref )` — a paren-wrapped
		// single relation is not a valid joined_table operand. The
		// SCENARIOS "verify" resolves to REJECT.
		{
			name:     "paren-single-relation JOIN relation rejected",
			sql:      `SELECT * FROM ((T) JOIN U ON TRUE)`,
			expected: OmniRejected,
		},

		// --- subquery as right-hand operand ---
		// The outer paren wraps the joined_table; the right operand is a
		// select_with_parens table reference with column-list alias.
		{
			name:     "relation JOIN subquery with column alias",
			sql:      `SELECT * FROM (T JOIN (SELECT 1) s(x) ON TRUE)`,
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
