//go:build oracle

package parser

import (
	"fmt"
	"strings"
	"testing"
)

// TestParenOracleIn transcribes SCENARIOS-pg-paren-dispatch.md §3.1
// (parseInExpr oracle corpus) into PG 17 testcontainer-backed parity
// probes. The §1.1 handwritten tests in paren_in_expr_test.go cover
// omni's accept/reject + AST-shape obligations. §3.1 extends those to
// full oracle parity: every `IN (...)` shape below is actually
// submitted to PG 17 so any future parser change that drifts from PG's
// accept/reject contract lights up as a failing probe here.
//
// Scope caveat (from the dispatch prompt):
//   - classifyOmni inspects FROM-clause items to categorize subquery vs
//     joined_table vs other. For §3.1 probes the FROM clause is always
//     bare `T`, so Items[0] is a *nodes.RangeVar → OmniOther on accept
//     cases. The exception is scenario 11 (IN inside a JOIN ON), where
//     FROM carries a *nodes.JoinExpr → OmniJoinedTable.
//   - The oracle value here is reject-parity (does omni reject what PG
//     rejects?) and accept-parity (does omni accept what PG accepts?),
//     not routing classification — the §1.1 handwritten suite carries
//     the routing-shape proofs.
//
// Fixture (from paren_oracle_test.go):
//
//	T(a int, b int, x int, y int, doc xml)
//	U(a int, b int, x int, y int)
//	V(a int, b int)
//	foo(a int)
func TestParenOracleIn(t *testing.T) {
	o := StartParenOracle(t)

	type probe struct {
		name     string
		sql      string
		expected OmniStatus
	}

	var cases []probe

	// --- Scenarios 1-2: oracle parity for the two basic shapes ---
	// §1.1 handwritten tests prove the raw AST shape; §3.1 confirms PG
	// agrees on accept/reject through a live probe.
	cases = append(cases,
		probe{
			name:     "expr_list basic IN (1,2,3)",
			sql:      `SELECT * FROM T WHERE a IN (1, 2, 3)`,
			expected: OmniOther,
		},
		probe{
			name:     "subquery basic IN (SELECT 1)",
			sql:      `SELECT * FROM T WHERE a IN (SELECT 1)`,
			expected: OmniOther,
		},
	)

	// --- Scenario 5: list-size variants (single, 2-element, 10-element,
	// 100-element) — table-driven. PG's grammar imposes no upper bound on
	// expr_list so all four sizes must accept.
	for _, n := range []int{1, 2, 10, 100} {
		parts := make([]string, n)
		for i := 0; i < n; i++ {
			parts[i] = fmt.Sprintf("%d", i)
		}
		cases = append(cases, probe{
			name:     fmt.Sprintf("list size %d", n),
			sql:      fmt.Sprintf(`SELECT * FROM T WHERE a IN (%s)`, strings.Join(parts, ", ")),
			expected: OmniOther,
		})
	}

	// --- Scenario 6: literal-kind variants (int, float, string, bool,
	// null). All must accept — PG's a_expr on the RHS side of IN accepts
	// any constant.
	literalKinds := []struct {
		name string
		lit  string
	}{
		{"int literal", `1, 2, 3`},
		{"float literal", `1.5, 2.5, 3.25`},
		{"string literal", `'a', 'b', 'c'`},
		{"bool literal", `true, false`},
		{"null literal", `NULL, NULL`},
	}
	for _, lk := range literalKinds {
		cases = append(cases, probe{
			name:     "literal " + lk.name,
			sql:      fmt.Sprintf(`SELECT * FROM T WHERE a IN (%s)`, lk.lit),
			expected: OmniOther,
		})
	}

	// --- Scenario 7: subquery-kind variants (SELECT, VALUES,
	// WITH...SELECT, TABLE, SELECT UNION SELECT).
	subqueryKinds := []struct {
		name  string
		inRHS string
	}{
		{"SELECT subquery", `SELECT 1`},
		{"VALUES subquery", `VALUES (1), (2)`},
		{"WITH...SELECT subquery", `WITH cte AS (SELECT 1) SELECT * FROM cte`},
		{"TABLE subquery", `TABLE foo`},
		{"SELECT UNION SELECT subquery", `SELECT 1 UNION SELECT 2`},
	}
	for _, sk := range subqueryKinds {
		cases = append(cases, probe{
			name:     "subquery " + sk.name,
			sql:      fmt.Sprintf(`SELECT * FROM T WHERE a IN (%s)`, sk.inRHS),
			expected: OmniOther,
		})
	}

	// --- Scenario 8: row-constructor LHS — list RHS + subquery RHS.
	cases = append(cases,
		probe{
			name:     "row constructor LHS, list RHS",
			sql:      `SELECT * FROM T WHERE (a, b) IN ((1, 2), (3, 4))`,
			expected: OmniOther,
		},
		probe{
			name:     "row constructor LHS, subquery RHS",
			sql:      `SELECT * FROM T WHERE (a, b) IN (SELECT a, b FROM U)`,
			expected: OmniOther,
		},
	)

	// --- Scenario 9: NOT IN expr_list.
	cases = append(cases, probe{
		name:     "NOT IN expr_list",
		sql:      `SELECT * FROM T WHERE a NOT IN (1, 2, 3)`,
		expected: OmniOther,
	})

	// --- Scenario 10: NOT IN subquery.
	cases = append(cases, probe{
		name:     "NOT IN subquery",
		sql:      `SELECT * FROM T WHERE a NOT IN (SELECT a FROM U)`,
		expected: OmniOther,
	})

	// --- Scenario 12: IN in HAVING.
	cases = append(cases, probe{
		name:     "IN in HAVING",
		sql:      `SELECT count(*) FROM T GROUP BY a HAVING count(*) IN (1, 2, 3)`,
		expected: OmniOther,
	})

	// --- Scenario 13: IN in CASE WHEN.
	cases = append(cases, probe{
		name:     "IN in CASE WHEN",
		sql:      `SELECT CASE WHEN a IN (1, 2) THEN 'y' ELSE 'n' END FROM T`,
		expected: OmniOther,
	})

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			assertParenParity(t, o, tc.sql, tc.expected)
		})
	}

	// --- Scenario 11: IN inside a JOIN ON — FROM becomes a JoinExpr so
	// Items[0] classifies as OmniJoinedTable. Kept as a separate
	// subtest so the expected-status override is explicit.
	t.Run("IN in JOIN ON", func(t *testing.T) {
		sql := `SELECT * FROM T JOIN U ON T.a IN (SELECT a FROM V)`
		assertParenParity(t, o, sql, OmniJoinedTable)
	})
}
