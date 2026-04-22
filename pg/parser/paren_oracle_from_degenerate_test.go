//go:build oracle

package parser

import "testing"

// TestParenOracleFromDegenerate transcribes SCENARIOS-pg-paren-dispatch.md §2.7
// (FROM-clause degenerate / malformed corpus) into oracle probes. Every case
// here is a syntactically invalid paren arrangement in FROM: empty parens,
// unbalanced parens, missing join operands, qualifiers on join types that
// disallow them (CROSS / NATURAL), and other grammar violations.
//
// All probes assert OmniRejected — PG must raise 42601 syntax_error and omni
// must produce a parse error. This section is the reject-side fence that
// complements §2.2–§2.6's accept/joined-table/subquery discrimination: any
// input that slips through omni but PG rejects (or vice versa) is a bug.
//
// Fixture tables T/U are created by StartParenOracle (see paren_oracle_test.go).
func TestParenOracleFromDegenerate(t *testing.T) {
	o := StartParenOracle(t)

	cases := []struct {
		name     string
		sql      string
		expected OmniStatus
		skip     string // non-empty means t.Skip with this reason
	}{
		// --- empty / unbalanced parens ---
		{
			name:     "empty parens",
			sql:      `SELECT * FROM ()`,
			expected: OmniRejected,
		},
		{
			name:     "unclosed open paren",
			sql:      `SELECT * FROM (`,
			expected: OmniRejected,
		},
		{
			name:     "stray close paren",
			sql:      `SELECT * FROM )`,
			expected: OmniRejected,
		},
		{
			name:     "unclosed after subquery",
			sql:      `SELECT * FROM (SELECT 1`,
			expected: OmniRejected,
		},

		// --- malformed joined_table content ---
		{
			name:     "JOIN without right operand",
			sql:      `SELECT * FROM (T JOIN`,
			expected: OmniRejected,
		},
		{
			// PAREN-KB-1 closed: parseJoinQual now requires ON/USING for
			// every non-CROSS/non-NATURAL JOIN. Matches PG 17's 42601 at
			// the closing `)` for `FROM (T JOIN U)`.
			name:     "inner JOIN missing qual",
			sql:      `SELECT * FROM (T JOIN U)`,
			expected: OmniRejected,
		},
		{
			name:     "CROSS JOIN with ON qual",
			sql:      `SELECT * FROM (T CROSS JOIN U ON TRUE)`,
			expected: OmniRejected,
		},
		{
			name:     "NATURAL JOIN with ON qual",
			sql:      `SELECT * FROM (T NATURAL JOIN U ON TRUE)`,
			expected: OmniRejected,
		},

		// --- malformed subquery content ---
		{
			// DIVERGENCE FROM SCENARIOS: the §2.7 line says "reject" but PG
			// 17 actually ACCEPTS `SELECT * FROM (SELECT)` — an empty
			// target-list SELECT is valid syntax in PG (the "null-column"
			// form). omni also accepts as a RangeSubselect. Expectation is
			// resolved to OmniSubquery so the oracle fence reflects reality;
			// the SCENARIOS line documentation is imprecise and will be
			// updated by the driver at merge.
			name:     "incomplete SELECT (PG accepts empty target list)",
			sql:      `SELECT * FROM (SELECT)`,
			expected: OmniSubquery,
		},
		{
			name:     "extra close paren",
			sql:      `SELECT * FROM (SELECT 1))`,
			expected: OmniRejected,
		},
		{
			name:     "missing close paren",
			sql:      `SELECT * FROM ((SELECT 1)`,
			expected: OmniRejected,
		},
		{
			name:     "unclosed FROM inside subquery",
			sql:      `SELECT * FROM (SELECT 1 FROM)`,
			expected: OmniRejected,
		},

		// --- leading empty parens + stray SELECT ---
		{
			name:     "leading empty parens with stray SELECT",
			sql:      `SELECT * FROM ( ) SELECT 1`,
			expected: OmniRejected,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if tc.skip != "" {
				t.Skip(tc.skip)
			}
			assertParenParity(t, o, tc.sql, tc.expected)
		})
	}
}
