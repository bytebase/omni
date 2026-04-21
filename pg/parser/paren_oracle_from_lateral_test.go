//go:build oracle

package parser

import "testing"

// TestParenOracleFromLateral transcribes SCENARIOS-pg-paren-dispatch.md §2.6
// (FROM-clause LATERAL interactions corpus) into oracle probes. Each case
// asserts the OmniStatus omni must emit to be PG-17 aligned for a FROM list
// whose second item is LATERAL-prefixed.
//
// Classifier note: the harness's classifyOmni inspects `FromClause.Items[0]`
// only. For every scenario in §2.6 the FROM list starts with `T` (a plain
// RangeVar), so the accept cases all classify as OmniOther — "parsed, but not
// a RangeSubselect or JoinExpr at Items[0]". The PG-reject case (LATERAL
// with a joined_table operand) must surface as OmniRejected because the
// whole SELECT fails to parse. This matches the invariant assertParenParity
// enforces: PGReject iff OmniRejected.
//
// Coverage outline:
//   - LATERAL select_with_parens (bare, aliased, double-wrapped, set-op)
//   - LATERAL select_with_parens with a joined_table body → PG reject
//   - LATERAL xmltable
//   - LATERAL json_table
//
// Fixture table T with columns (a, b, x, y, doc xml) is created by
// StartParenOracle (see paren_oracle_test.go). The JSON_TABLE scenario
// references T.doc, which is xml — PG will raise a datatype_mismatch
// (semantic, not 42601), so classifyPG still reports Accept. The grammar
// decision is what §2.6 cares about.
func TestParenOracleFromLateral(t *testing.T) {
	o := StartParenOracle(t)

	cases := []struct {
		name     string
		sql      string
		expected OmniStatus
	}{
		// --- LATERAL select_with_parens ---
		// Items[0] is T (RangeVar) → OmniOther. The LATERAL arm is at
		// Items[1], not inspected by classifyOmni. These probes still pull
		// their weight: they prove PG accepts the LATERAL shape (no 42601)
		// and omni round-trips without error.
		{
			name:     "LATERAL bare subquery",
			sql:      `SELECT * FROM T, LATERAL (SELECT 1)`,
			expected: OmniOther,
		},
		{
			name:     "LATERAL subquery with alias",
			sql:      `SELECT * FROM T, LATERAL (SELECT 1) x`,
			expected: OmniOther,
		},
		{
			name:     "LATERAL double-wrapped subquery",
			sql:      `SELECT * FROM T, LATERAL ((SELECT 1))`,
			expected: OmniOther,
		},
		// BYT-9315-adjacent: parenthesized set-op operands under LATERAL.
		// Inner parens must route to select_with_parens on the way to the
		// UNION's simple_select arms; outer paren routes the whole UNION
		// as LATERAL's select_with_parens body.
		{
			name:     "LATERAL parenthesized UNION",
			sql:      `SELECT * FROM T, LATERAL ((SELECT 1) UNION (SELECT 2))`,
			expected: OmniOther,
		},

		// --- LATERAL joined_table is NOT a grammar production ---
		// PG grammar only admits LATERAL + {func_table, select_with_parens,
		// xmltable, json_table}. A bare joined_table body like
		// `a JOIN b ON ...` inside the parens is rejected because
		// select_with_parens requires a select_no_parens / values_clause,
		// not a joined_table. Expected outcome: 42601 at `JOIN`.
		{
			name:     "LATERAL joined_table rejected",
			sql:      `SELECT * FROM T, LATERAL (a JOIN b ON t.x = a.x)`,
			expected: OmniRejected,
		},

		// --- LATERAL xmltable ---
		// T.doc is xml, so the PASSING operand type-checks; PG parses and
		// executes successfully against an empty T. omni produces a
		// RangeTableFunc at Items[1], but classifyOmni only looks at
		// Items[0] (T, RangeVar) → OmniOther.
		{
			name:     "LATERAL xmltable",
			sql:      `SELECT * FROM T, LATERAL XMLTABLE('/root' PASSING T.doc COLUMNS a int PATH 'a')`,
			expected: OmniOther,
		},

		// --- LATERAL json_table ---
		// T.doc is xml, not jsonb — PG will raise a semantic type error
		// (42804 / 42883) when the statement runs. That's still Accept
		// for classifyPG: the grammar decision was made before name and
		// type resolution. Omni emits a JsonTable node at Items[1];
		// Items[0] stays T → OmniOther.
		{
			name:     "LATERAL json_table",
			sql:      `SELECT * FROM T, LATERAL JSON_TABLE(T.doc, '$' COLUMNS(a int PATH '$.a'))`,
			expected: OmniOther,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			assertParenParity(t, o, tc.sql, tc.expected)
		})
	}
}
