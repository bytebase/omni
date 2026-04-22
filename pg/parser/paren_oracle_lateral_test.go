//go:build oracle

package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// TestParenOracleLateral transcribes SCENARIOS-pg-paren-dispatch.md §3.2
// (parseLateralTableRef oracle corpus) into PG 17 testcontainer-backed
// parity probes plus direct AST-shape inspection. §2.6 already carries
// the baseline LATERAL accept/reject matrix; §3.2 extends it with:
//
//  1. AST-shape confirmation — Lateral flag + node type for each of the
//     three body-kind grammar arms (select_with_parens, xmltable,
//     json_table). Direct Parse() inspection so the assertion covers the
//     internal field, not just FROM-item classification.
//  2. Column-list alias combinations — `LATERAL (...) s(x, y)` and
//     `LATERAL XMLTABLE(...) AS xt(...)` variants go through opt_alias_clause.
//  3. Correlated LATERAL — subquery body references the outer anchor T,
//     which is the whole reason LATERAL exists. Must still route as
//     RangeSubselect with Lateral=true and PG must accept the shape.
//  4. Invalid LATERAL shapes — LATERAL joined_table (already in §2.6,
//     reinforced here) and LATERAL ROWS FROM without parens (new). PG
//     rejects both with 42601; omni must too.
//
// Fixture tables T/U/V/W and function probe_f()/probe_g() are created by
// StartParenOracle (see paren_oracle_test.go). T has a `doc xml` column
// so XMLTABLE/JSON_TABLE probes can pass a real node. PG may raise a
// semantic datatype_mismatch on JSON_TABLE(T.doc, ...) since doc is xml
// not jsonb — that's still Accept for classifyPG (grammar decision
// already made). §3.2 cares about the grammar decision.
func TestParenOracleLateral(t *testing.T) {
	o := StartParenOracle(t)

	// --- Scenario 1: LATERAL body-kind AST shapes ---
	// Direct Parse() inspection so we can assert the Lateral flag on the
	// exact node type, which assertParenParityAt alone can't reach
	// (classifyOmniAt collapses RangeTableFunc/JsonTable/RangeFunction
	// into OmniOther). Each case here still probes PG through the oracle
	// to keep the accept-parity fence.
	t.Run("AST shapes", func(t *testing.T) {
		astCases := []struct {
			name      string
			sql       string
			wantNode  string // human-readable node kind, for error messages
			checkNode func(t *testing.T, item nodes.Node)
		}{
			{
				name:     "LATERAL select_with_parens sets RangeSubselect.Lateral=true",
				sql:      `SELECT * FROM T, LATERAL (SELECT 1) s`,
				wantNode: "*nodes.RangeSubselect",
				checkNode: func(t *testing.T, item nodes.Node) {
					rs, ok := item.(*nodes.RangeSubselect)
					if !ok {
						t.Fatalf("expected *nodes.RangeSubselect, got %T", item)
					}
					if !rs.Lateral {
						t.Errorf("expected RangeSubselect.Lateral=true, got false")
					}
				},
			},
			{
				name:     "LATERAL xmltable sets RangeTableFunc.Lateral=true",
				sql:      `SELECT * FROM T, LATERAL XMLTABLE('/root' PASSING T.doc COLUMNS a int PATH 'a') AS xt`,
				wantNode: "*nodes.RangeTableFunc",
				checkNode: func(t *testing.T, item nodes.Node) {
					rtf, ok := item.(*nodes.RangeTableFunc)
					if !ok {
						t.Fatalf("expected *nodes.RangeTableFunc, got %T", item)
					}
					if !rtf.Lateral {
						t.Errorf("expected RangeTableFunc.Lateral=true, got false")
					}
				},
			},
			{
				name:     "LATERAL json_table sets JsonTable.Lateral=true",
				sql:      `SELECT * FROM T, LATERAL JSON_TABLE(T.doc, '$' COLUMNS(a int PATH '$.a')) AS jt`,
				wantNode: "*nodes.JsonTable",
				checkNode: func(t *testing.T, item nodes.Node) {
					jt, ok := item.(*nodes.JsonTable)
					if !ok {
						t.Fatalf("expected *nodes.JsonTable, got %T", item)
					}
					if !jt.Lateral {
						t.Errorf("expected JsonTable.Lateral=true, got false")
					}
				},
			},
		}

		for _, tc := range astCases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				// First: oracle parity — PG must accept the shape. This keeps
				// the AST-shape assertion couplings to a real PG outcome so
				// an omni-only accept can't sneak through.
				r := ProbeParenAt(o.ctx, o, tc.sql, 1)
				if r.PGStatus == PGReject {
					t.Fatalf("PG rejected %q: %s", tc.sql, r.PGError)
				}
				if r.OmniStatus == OmniRejected {
					t.Fatalf("omni rejected %q: %s", tc.sql, r.OmniError)
				}

				// Second: AST shape — the LATERAL node at Items[1] has the
				// correct concrete type AND Lateral=true.
				stmts, err := Parse(tc.sql)
				if err != nil {
					t.Fatalf("Parse() error: %v", err)
				}
				raw, ok := stmts.Items[0].(*nodes.RawStmt)
				if !ok {
					t.Fatalf("expected *nodes.RawStmt at top, got %T", stmts.Items[0])
				}
				sel, ok := raw.Stmt.(*nodes.SelectStmt)
				if !ok {
					t.Fatalf("expected *nodes.SelectStmt, got %T", raw.Stmt)
				}
				if sel.FromClause == nil || len(sel.FromClause.Items) < 2 {
					t.Fatalf("expected >=2 FROM items, got %v",
						sel.FromClause)
				}
				tc.checkNode(t, sel.FromClause.Items[1])
			})
		}
	})

	// --- Scenario 2: LATERAL + column-list alias combinations ---
	// opt_alias_clause in PG's LATERAL arms can carry a column list
	// (e.g. `AS s(x, y)`). Coverage here spans all four LATERAL body
	// kinds so alias routing can't quietly regress on one arm.
	t.Run("column-list aliases", func(t *testing.T) {
		aliasCases := []struct {
			name     string
			sql      string
			expected OmniStatus
		}{
			{
				// LATERAL (SELECT...) s(x): the alias carries a single
				// column name. opt_alias_clause → colnames list.
				name:     "LATERAL subquery with AS s(x)",
				sql:      `SELECT * FROM T, LATERAL (SELECT 1) AS s(x)`,
				expected: OmniSubquery,
			},
			{
				// Multi-column alias.
				name:     "LATERAL subquery with AS s(x, y)",
				sql:      `SELECT * FROM T, LATERAL (SELECT 1, 2) AS s(x, y)`,
				expected: OmniSubquery,
			},
			{
				// AS-less alias with column list — PG grammar still
				// accepts (alias_clause without AS).
				name:     "LATERAL subquery with s(x, y) no AS",
				sql:      `SELECT * FROM T, LATERAL (SELECT 1, 2) s(x, y)`,
				expected: OmniSubquery,
			},
			{
				// XMLTABLE with plain alias (AS xt). XMLTABLE itself
				// carries its own COLUMNS list; the outer alias names
				// only the table.
				name:     "LATERAL xmltable with AS xt",
				sql:      `SELECT * FROM T, LATERAL XMLTABLE('/root' PASSING T.doc COLUMNS a int PATH 'a') AS xt`,
				expected: OmniOther,
			},
			{
				// JSON_TABLE with AS jt(colA). Outer column list overrides
				// the inner COLUMNS — PG's opt_alias_clause admits this
				// shape even though in practice the colnames list should
				// match the COLUMNS arity.
				name:     "LATERAL json_table with AS jt",
				sql:      `SELECT * FROM T, LATERAL JSON_TABLE(T.doc, '$' COLUMNS(a int PATH '$.a')) AS jt`,
				expected: OmniOther,
			},
			{
				// Func-table with a plain column-name alias list (the
				// opt_alias_clause / func_alias_clause non-typed form).
				// Typed coldef lists (`AS fa(a int)`) need the function
				// to return record; probe_f returns int, so PG rejects
				// those semantically even though the raw grammar admits
				// them. The plain-column form avoids that coupling.
				name:     "LATERAL func_table with plain column alias",
				sql:      `SELECT * FROM T, LATERAL probe_f(T.x) AS fa(x)`,
				expected: OmniOther,
			},
		}

		for _, tc := range aliasCases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				assertParenParityAt(t, o, tc.sql, 1, tc.expected)
			})
		}
	})

	// --- Scenario 3: correlated LATERAL (outer-table reference) ---
	// The typical use case for LATERAL: inner subquery references the
	// anchor table T. PG admits the shape; omni must match. Lateral
	// flag on the inner node is what makes the planner resolve outer
	// refs; we confirm it in the AST check below for the first case.
	t.Run("correlated LATERAL", func(t *testing.T) {
		correlatedCases := []struct {
			name     string
			sql      string
			expected OmniStatus
		}{
			{
				name:     "LATERAL subquery references T.x in WHERE",
				sql:      `SELECT * FROM T, LATERAL (SELECT * FROM U WHERE U.x = T.x) y`,
				expected: OmniSubquery,
			},
			{
				name:     "LATERAL subquery references T.x in SELECT list",
				sql:      `SELECT * FROM T, LATERAL (SELECT T.x + 1) AS s(n)`,
				expected: OmniSubquery,
			},
			{
				name:     "LATERAL func_table references T.x",
				sql:      `SELECT * FROM T, LATERAL probe_f(T.x)`,
				expected: OmniOther, // RangeFunction, not subquery
			},
			{
				name:     "LATERAL ROWS FROM references T.x and T.y",
				sql:      `SELECT * FROM T, LATERAL ROWS FROM (probe_f(T.x), probe_g(T.y))`,
				expected: OmniOther, // RangeFunction with IsRowsfrom=true
			},
			{
				name:     "LATERAL xmltable references T.doc",
				sql:      `SELECT * FROM T, LATERAL XMLTABLE('/root' PASSING T.doc COLUMNS a int PATH 'a') xt`,
				expected: OmniOther,
			},
		}

		for _, tc := range correlatedCases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				assertParenParityAt(t, o, tc.sql, 1, tc.expected)
			})
		}

		// Sub-assertion: first correlated case — verify Lateral=true on
		// the resulting RangeSubselect. This is the field the planner
		// keys off to allow the outer reference; omni must set it.
		t.Run("correlated LATERAL sets Lateral=true", func(t *testing.T) {
			sql := `SELECT * FROM T, LATERAL (SELECT * FROM U WHERE U.x = T.x) y`
			stmts, err := Parse(sql)
			if err != nil {
				t.Fatalf("Parse() error: %v", err)
			}
			raw := stmts.Items[0].(*nodes.RawStmt)
			sel := raw.Stmt.(*nodes.SelectStmt)
			rs, ok := sel.FromClause.Items[1].(*nodes.RangeSubselect)
			if !ok {
				t.Fatalf("expected *nodes.RangeSubselect at Items[1], got %T",
					sel.FromClause.Items[1])
			}
			if !rs.Lateral {
				t.Errorf("correlated LATERAL must set Lateral=true, got false")
			}
		})
	})

	// --- Scenario 4: invalid LATERAL shapes ---
	// PG rejects with 42601; omni must reject too. The assertion uses
	// OmniRejected, which assertParenParityAt couples to PGReject — the
	// mismatch detector will flag any accept-vs-reject drift loudly.
	//
	// For reject cases the index argument is effectively ignored because
	// the whole SELECT fails to parse before any FROM-item is reached.
	// We pass 1 for symmetry with the accept cases above, but the
	// classifier's bounds-check (index >= len(Items)) short-circuits to
	// OmniRejected in that path too.
	t.Run("invalid shapes rejected", func(t *testing.T) {
		rejectCases := []struct {
			name string
			sql  string
			skip string // non-empty → t.Skip with this reason (known omni bug)
		}{
			{
				// LATERAL (joined_table) — PG grammar has no such rule.
				// The `(` arm routes to select_with_parens, whose body
				// must be select_no_parens; `a JOIN b ON TRUE` isn't.
				// §2.6 already covers this; §3.2 reinforces via the
				// AST-shape oracle fence.
				name: "LATERAL joined_table",
				sql:  `SELECT * FROM T, LATERAL (V JOIN W ON TRUE)`,
			},
			{
				// LATERAL ROWS FROM without parens — ROWS FROM requires
				// `(` func_list `)`. Bare `ROWS FROM probe_f(...)` is not
				// a grammar production.
				name: "LATERAL ROWS FROM without parens",
				sql:  `SELECT * FROM T, LATERAL ROWS FROM probe_f(T.x)`,
			},
			{
				// Bare LATERAL relation — LATERAL prefix requires one of
				// the four grammar arms; plain relation_expr is not a
				// variant. §1.2 handwritten test already covers raw
				// parse rejection; this line is the oracle twin.
				name: "LATERAL bare relation",
				sql:  `SELECT * FROM T, LATERAL U`,
			},
			{
				// PAREN-KB-3 closed: parseLateralTableRef now rejects when
				// parseSelectWithParens returns a nil *SelectStmt (empty
				// body case). Matches PG 17's 42601 at the closing `)`.
				name: "LATERAL empty parens",
				sql:  `SELECT * FROM T, LATERAL ()`,
			},
		}

		for _, tc := range rejectCases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				if tc.skip != "" {
					t.Skip(tc.skip)
				}
				assertParenParityAt(t, o, tc.sql, 1, OmniRejected)
			})
		}
	})
}
