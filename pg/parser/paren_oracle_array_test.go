//go:build oracle

package parser

import (
	"fmt"
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// TestParenOracleArray transcribes SCENARIOS-pg-paren-dispatch.md §3.3
// (parseArrayCExpr oracle corpus) into PG 17 testcontainer-backed
// parity probes. The §1.3 and §1.4 handwritten tests in
// paren_array_expr_test.go carry omni's AST-shape obligations
// (A_ArrayExpr vs SubLink{ARRAY_SUBLINK}); §3.3 extends those to full
// oracle accept/reject parity: every ARRAY[...] / ARRAY(...) shape
// below is submitted to PG 17 so any future parser drift from PG's
// accept/reject contract lights up as a failing probe here.
//
// Scope caveat (from the dispatch prompt):
//   - The ARRAY expression always lives in the SELECT list (target
//     list), never in FROM. With `SELECT ARRAY[...] FROM T` the FROM
//     item is a bare *nodes.RangeVar → OmniOther. With bare
//     `SELECT ARRAY[...]` (no FROM clause at all) classifyOmni also
//     returns OmniOther (FromClause nil branch). Either way Items[0]
//     is NOT the probed node — the §3.3 value is reject-parity and
//     accept-parity on the grammar decision (parseArrayCExpr's
//     `[` vs `(` split), not FROM routing classification.
//   - The AST-shape assertion for "ARRAY[...] produces A_ArrayExpr"
//     vs "ARRAY(...) produces SubLink{ARRAY_SUBLINK}" lives in the
//     handwritten §1.3 tests (paren_array_expr_test.go). §3.3 can
//     only verify "PG accepts → omni accepts"; the deep AST-shape
//     check is carried by §1.3.
//   - `SELECT ARRAY[]` is accepted by PG's raw parser; the post-
//     analyze "cannot determine type of empty array" error is
//     SQLSTATE 42P18 (indeterminate_datatype), NOT 42601. classifyPG
//     only reads 42601 as Reject — any other SQLSTATE is Accept.
//     So `SELECT ARRAY[]` lands as PG=Accept, omni=Accept, expected
//     OmniOther. The explicit cast variant `ARRAY[]::int[]` is fully
//     accepted by PG (cast resolves element type at analyze time).
//
// Fixture (from paren_oracle_test.go): table `foo(a int)` is used
// for `ARRAY(TABLE foo)`; tables T/U/V/W cover SELECT-list anchors.
func TestParenOracleArray(t *testing.T) {
	o := StartParenOracle(t)

	// --- Scenario 1: ARRAY[...] A_ArrayExpr vs ARRAY(...) SubLink ---
	// Parity probe for the two branches of parseArrayCExpr. PG must
	// accept both; omni must too. Deep AST-shape assertion lives in
	// the §1.3 handwritten tests — here we also do one direct Parse()
	// check to confirm the two shapes remain distinguishable at the
	// omni side (guards against a future refactor that collapses them).
	t.Run("dispatch shapes", func(t *testing.T) {
		shapeCases := []struct {
			name     string
			sql      string
			wantNode string
			check    func(t *testing.T, val nodes.Node)
		}{
			{
				name:     "ARRAY[...] produces A_ArrayExpr",
				sql:      `SELECT ARRAY[1, 2, 3] FROM T`,
				wantNode: "*nodes.A_ArrayExpr",
				check: func(t *testing.T, val nodes.Node) {
					if _, ok := val.(*nodes.A_ArrayExpr); !ok {
						t.Errorf("expected *nodes.A_ArrayExpr, got %T", val)
					}
				},
			},
			{
				name:     "ARRAY(SELECT ...) produces SubLink",
				sql:      `SELECT ARRAY(SELECT 1) FROM T`,
				wantNode: "*nodes.SubLink",
				check: func(t *testing.T, val nodes.Node) {
					sl, ok := val.(*nodes.SubLink)
					if !ok {
						t.Fatalf("expected *nodes.SubLink, got %T", val)
					}
					if sl.SubLinkType != int(nodes.ARRAY_SUBLINK) {
						t.Errorf("expected SubLinkType=ARRAY_SUBLINK (%d), got %d",
							int(nodes.ARRAY_SUBLINK), sl.SubLinkType)
					}
				},
			},
		}

		for _, tc := range shapeCases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				// Oracle parity first: PG must accept; omni must accept.
				// Items[0] is T (RangeVar) so we expect OmniOther.
				assertParenParity(t, o, tc.sql, OmniOther)

				// AST shape confirmation: parse locally and inspect the
				// target-list expression. This guards against omni silently
				// collapsing the two productions into one node type.
				stmts, err := Parse(tc.sql)
				if err != nil {
					t.Fatalf("omni Parse(%q) failed: %v", tc.sql, err)
				}
				sel := stmts.Items[0].(*nodes.RawStmt).Stmt.(*nodes.SelectStmt)
				tl := sel.TargetList.Items[0].(*nodes.ResTarget)
				tc.check(t, tl.Val)
			})
		}
	})

	// --- Scenario 2: nested ARRAY[ARRAY[...]] constructions ---
	// PG admits nested array_expr_list arbitrarily deep. Each level is
	// a separate A_ArrayExpr in the raw-parse AST; the semantic
	// Multidims flag is set by transformArrayExpr post-analyze.
	t.Run("nested ARRAY constructors", func(t *testing.T) {
		nestedCases := []struct {
			name string
			sql  string
		}{
			{
				name: "2-level ARRAY[ARRAY[...],ARRAY[...]]",
				sql:  `SELECT ARRAY[ARRAY[1, 2], ARRAY[3, 4]] FROM T`,
			},
			{
				name: "3-level ARRAY[ARRAY[ARRAY[...]]]",
				sql:  `SELECT ARRAY[ARRAY[ARRAY[1, 2]]] FROM T`,
			},
			{
				name: "implicit nested array_expr_list [[1,2],[3,4]]",
				sql:  `SELECT ARRAY[[1, 2], [3, 4]] FROM T`,
			},
			{
				name: "asymmetric nesting mix (raw parse accepts)",
				// PG's raw parser accepts asymmetric nested arrays;
				// transformArrayExpr would later reject with
				// "multidimensional arrays must have array expressions
				// with matching dimensions". classifyPG sees grammar
				// acceptance → PGAccept.
				sql: `SELECT ARRAY[ARRAY[1, 2], ARRAY[3]] FROM T`,
			},
			{
				name: "nested sublinks ARRAY(SELECT ARRAY(SELECT 1))",
				sql:  `SELECT ARRAY(SELECT ARRAY(SELECT 1)) FROM T`,
			},
		}

		for _, tc := range nestedCases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				assertParenParity(t, o, tc.sql, OmniOther)
			})
		}
	})

	// --- Scenario 3: ARRAY with type cast combinations ---
	// Type casts wrap the ARRAY expression; the parser's `(` / `[`
	// dispatch at parseArrayCExpr is unaffected by the outer cast.
	// Each variant must accept on both sides.
	t.Run("type cast combinations", func(t *testing.T) {
		castCases := []struct {
			name string
			sql  string
		}{
			{
				name: "ARRAY[...]::int[] postfix cast",
				sql:  `SELECT ARRAY[1, 2]::int[] FROM T`,
			},
			{
				name: "ARRAY[...]::text[] postfix cast",
				sql:  `SELECT ARRAY['a', 'b']::text[] FROM T`,
			},
			{
				name: "CAST(ARRAY[...] AS int[]) prefix form",
				sql:  `SELECT CAST(ARRAY[1, 2] AS int[]) FROM T`,
			},
			{
				name: "ARRAY(SELECT ...)::int[] sublink + cast",
				sql:  `SELECT ARRAY(SELECT 1)::int[] FROM T`,
			},
			{
				name: "parenthesized ARRAY then cast",
				sql:  `SELECT (ARRAY[1, 2])::int[] FROM T`,
			},
			{
				name: "inner element casts ARRAY[1::text, 2::text]",
				sql:  `SELECT ARRAY[1::text, 2::text] FROM T`,
			},
			{
				name: "empty ARRAY with cast ARRAY[]::int[]",
				// Cast resolves the indeterminate element type —
				// PG accepts end-to-end. Without the cast PG raises
				// 42P18 (see empty-no-cast case below, still Accept
				// from classifyPG's POV).
				sql: `SELECT ARRAY[]::int[] FROM T`,
			},
			{
				name: "multi-dim cast ARRAY[[1,2],[3,4]]::int[][]",
				sql:  `SELECT ARRAY[[1, 2], [3, 4]]::int[][] FROM T`,
			},
		}

		for _, tc := range castCases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				assertParenParity(t, o, tc.sql, OmniOther)
			})
		}
	})

	// --- Scenario 4: ARRAY sublink with VALUES/TABLE/WITH/UNION/ORDER BY ---
	// PG's `ARRAY select_with_parens` production accepts any
	// select_no_parens body, which includes VALUES, TABLE, WITH, and
	// set-op combinations. §1.4 handwritten tests cover the omni side;
	// §3.3 verifies PG agrees via live oracle probe.
	//
	// Note on semantic outcomes: PG may raise non-42601 errors (e.g.
	// "subquery must return only one column" 42601 or
	// "subquery used as an expression returned more than one row"
	// 21000 at exec time). Those SQLSTATEs that classifyPG treats as
	// Accept still count as "parser accepted the shape" — which is
	// precisely what §3.3 is probing. The one exception is VALUES with
	// multi-column rows (`VALUES (1,2)`): that actually passes grammar
	// AND is valid input to a scalar sublink only if the outer target
	// can accept a row — for an int[] target, PG raises 42804
	// "subquery has too many columns" at analyze, still not 42601.
	t.Run("sublink body variants", func(t *testing.T) {
		sublinkCases := []struct {
			name string
			sql  string
		}{
			{
				name: "SELECT body",
				sql:  `SELECT ARRAY(SELECT 1) FROM T`,
			},
			{
				name: "VALUES body",
				sql:  `SELECT ARRAY(VALUES (1), (2)) FROM T`,
			},
			{
				name: "TABLE body",
				// `foo` is a pre-created fixture (see paren_oracle_test.go).
				sql: `SELECT ARRAY(TABLE foo) FROM T`,
			},
			{
				name: "WITH ... SELECT body",
				sql:  `SELECT ARRAY(WITH cte AS (SELECT 1) SELECT * FROM cte) FROM T`,
			},
			{
				name: "SELECT UNION SELECT body",
				sql:  `SELECT ARRAY(SELECT 1 UNION SELECT 2) FROM T`,
			},
			{
				name: "SELECT ... ORDER BY body",
				sql:  `SELECT ARRAY(SELECT a FROM U ORDER BY a) FROM T`,
			},
			{
				name: "SELECT ... LIMIT body",
				sql:  `SELECT ARRAY(SELECT a FROM U LIMIT 10) FROM T`,
			},
			{
				name: "SELECT DISTINCT body",
				sql:  `SELECT ARRAY(SELECT DISTINCT a FROM U) FROM T`,
			},
		}

		for _, tc := range sublinkCases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				assertParenParity(t, o, tc.sql, OmniOther)
			})
		}
	})

	// --- Scenario 5: negative cases — PG rejects, omni must too ---
	// Each SQL here violates the parseArrayCExpr content contract. PG
	// emits 42601 at the offending token; classifyPG returns
	// PGReject; assertParenParity requires omni to also reject. The
	// ARRAY[] case is the documented exception: PG's raw parser
	// ACCEPTS it (see paren_array_expr_test.go §1.3 notes); the
	// semantic error at analyze is 42P18, not 42601. So ARRAY[] is
	// NOT in this reject group — it sits in the accept-shapes list
	// above (via the `ARRAY[]::int[]` cast variant) and in a
	// dedicated accept probe below.
	t.Run("negative cases", func(t *testing.T) {
		rejectCases := []struct {
			name string
			sql  string
		}{
			{
				// ARRAY() — empty parens violate the `ARRAY
				// select_with_parens` production (select_with_parens
				// needs a select_no_parens body; empty is not one).
				// §1.4 covers omni's reject; this is the oracle twin.
				name: "ARRAY() empty parens",
				sql:  `SELECT ARRAY() FROM T`,
			},
			{
				// ARRAY[SELECT 1] — array_expr's content is expr_list
				// (a_expr list); bare SELECT is not a_expr.
				name: "ARRAY[SELECT 1]",
				sql:  `SELECT ARRAY[SELECT 1] FROM T`,
			},
			{
				// Trailing ARRAY with no body at all — parseArrayCExpr
				// requires `[` or `(` next; EOF is neither.
				name: "ARRAY with no body",
				sql:  `SELECT ARRAY FROM T`,
			},
			{
				// ARRAY(1) — single bare expression in parens is not
				// a select_no_parens. §1.4 content-contract coverage.
				name: "ARRAY(1) bare expr",
				sql:  `SELECT ARRAY(1) FROM T`,
			},
			{
				// ARRAY(1,2) — comma-list of bare exprs is not a
				// select_no_parens either.
				name: "ARRAY(1, 2) bare expr list",
				sql:  `SELECT ARRAY(1, 2) FROM T`,
			},
			{
				// ARRAY[,] — empty slot in expr_list.
				name: "ARRAY[,] empty slot",
				sql:  `SELECT ARRAY[,] FROM T`,
			},
			{
				// ARRAY[1 2] — missing comma between expr_list items.
				name: "ARRAY[1 2] missing comma",
				sql:  `SELECT ARRAY[1 2] FROM T`,
			},
			{
				// ROWS FROM is func_table-only (gram.y:13468); not a
				// select_no_parens lead. §1.4 handwritten test covers
				// omni's reject path via T7 post-parse nil check.
				name: "ARRAY(ROWS FROM (probe_f(1)))",
				sql:  `SELECT ARRAY(ROWS FROM (probe_f(1))) FROM T`,
			},
		}

		for _, tc := range rejectCases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				assertParenParity(t, o, tc.sql, OmniRejected)
			})
		}
	})

	// --- Scenario 5 continuation: ARRAY[] raw-parse accept ---
	// Pinned as its own subtest because the expected parity is
	// Accept/Accept with a SEMANTIC PG error that classifyPG maps to
	// PGAccept (42P18 indeterminate_datatype, NOT 42601). This is the
	// one "grammar-accepts-but-analyze-rejects" entry in §3.3.
	t.Run("ARRAY[] raw parse accepts (analyze-only reject)", func(t *testing.T) {
		// Run the probe manually so we can assert classifyPG's
		// PG-Accept classification (the mismatch detector's parity
		// check would also catch a drift, but this explicit assertion
		// documents the intent).
		sql := `SELECT ARRAY[] FROM T`
		r := ProbeParen(o.ctx, o, sql)
		if r.PGStatus != PGAccept {
			t.Errorf("expected PGAccept (42P18 is not 42601), got PG=%s err=%q",
				r.PGStatus, r.PGError)
		}
		if r.OmniStatus != OmniOther {
			t.Errorf("expected OmniOther (FROM item is RangeVar T), got omni=%s err=%q",
				r.OmniStatus, r.OmniError)
		}
	})

	// --- Extra accept-parity probe: ARRAY with no FROM clause ---
	// Covers the `FromClause == nil` branch of classifyOmniAt. The
	// parity test still applies: PG accepts and omni accepts, both
	// report OmniOther via different classifier paths (PG:
	// grammar accept; omni: nil-FROM → OmniOther default).
	t.Run("ARRAY without FROM clause", func(t *testing.T) {
		noFromCases := []string{
			`SELECT ARRAY[1, 2, 3]`,
			`SELECT ARRAY(SELECT 1)`,
			`SELECT ARRAY[]::int[]`,
			`SELECT ARRAY[ARRAY[1, 2]]`,
		}
		for i, sql := range noFromCases {
			sql := sql
			t.Run(fmt.Sprintf("case %d", i+1), func(t *testing.T) {
				assertParenParity(t, o, sql, OmniOther)
			})
		}
	})
}
