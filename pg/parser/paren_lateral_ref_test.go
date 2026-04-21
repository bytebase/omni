package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// TestParenLateralRefDispatch covers pg-paren-dispatch §1.2:
// parseLateralTableRef's choice among LATERAL's four valid grammar
// productions.
//
// PG grammar (gram.y:13611-13620, table_ref LATERAL alternatives):
//
//	table_ref:
//	    ...
//	  | LATERAL_P func_table func_alias_clause
//	  | LATERAL_P select_with_parens opt_alias_clause
//	  | LATERAL_P xmltable opt_alias_clause
//	  | LATERAL_P json_table opt_alias_clause
//
// All four alternatives have disjoint lead tokens after LATERAL:
// `(` → select_with_parens; XMLTABLE → xmltable; JSON_TABLE → json_table;
// everything else routes through func_expr_windowless (which admits
// ROWS FROM and func_expr_common_subexpr starters).  Notably absent:
// bare `LATERAL <relation>` — relation_expr is not a LATERAL variant.
//
// Before this fix omni's dispatch routed `(` to select_with_parens but
// silently accepted bare `LATERAL u` (a ColumnRef wrapped in a
// RangeFunction — not a valid func_application).  The T6 dispatch arms
// now explicitly route XMLTABLE / JSON_TABLE / ROWS / `(` and the
// func_table fallthrough rejects bare ColumnRef results.
func TestParenLateralRefDispatch(t *testing.T) {
	cases := []struct {
		sql       string
		wantParse bool
	}{
		// --- LATERAL select_with_parens --------------------------------------
		{`SELECT * FROM t, LATERAL (SELECT 1) x`, true},
		{`SELECT * FROM t, LATERAL (SELECT * FROM u WHERE u.x = t.x) y`, true},

		// --- LATERAL xmltable ------------------------------------------------
		{`SELECT * FROM t, LATERAL XMLTABLE('/root' PASSING x COLUMNS a int PATH 'a') as xt`, true},

		// --- LATERAL json_table ----------------------------------------------
		{`SELECT * FROM t, LATERAL JSON_TABLE(x, '$' COLUMNS(a int PATH '$.a')) as jt`, true},

		// --- LATERAL func_table (func_application) ---------------------------
		{`SELECT * FROM t, LATERAL f(t.x)`, true},
		{`SELECT * FROM t, LATERAL f(t.x) WITH ORDINALITY`, true},
		{`SELECT * FROM t, LATERAL f(t.x) AS fa(a int, b text)`, true},

		// --- LATERAL func_table (ROWS FROM) ----------------------------------
		{`SELECT * FROM t, LATERAL ROWS FROM (f(t.x), g(t.y))`, true},
		{`SELECT * FROM t, LATERAL ROWS FROM (f(t.x)) WITH ORDINALITY`, true},

		// --- LATERAL func_expr_common_subexpr (sql_value_function) -----------
		// USER / CURRENT_USER are reachable via func_expr_common_subexpr
		// and therefore a valid func_table. Tests the T6 dispatch
		// doesn't accidentally exclude this arm.
		{`SELECT * FROM t, LATERAL CURRENT_USER`, true},

		// --- Rejected: LATERAL joined_table ----------------------------------
		// PG grammar has no LATERAL joined_table production.  The `(`
		// arm routes to select_with_parens, which parses `a JOIN b ON
		// TRUE` as a SELECT statement — that fails because the inner
		// content is not a valid simple_select lead.
		{`SELECT * FROM t, LATERAL (a JOIN b ON TRUE)`, false},

		// --- Rejected: bare LATERAL relation ---------------------------------
		// Neither func_application (needs parens) nor
		// func_expr_common_subexpr (needs a keyword lead) admits a
		// plain identifier.  Before this fix omni silently accepted it.
		{`SELECT * FROM t, LATERAL u`, false},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			_, err := Parse(tc.sql)
			if tc.wantParse && err != nil {
				t.Fatalf("expected parse success, got error: %v", err)
			}
			if !tc.wantParse && err == nil {
				t.Fatalf("expected parse error, got nil")
			}
		})
	}
}

// TestParenLateralRefSubqueryAST asserts that LATERAL (SELECT ...)
// produces a RangeSubselect with Lateral=true.  The flag is what lets
// the analyzer permit outer-table references inside the subquery; a
// plain RangeSubselect (no LATERAL) would fail analysis at the first
// outer reference.
func TestParenLateralRefSubqueryAST(t *testing.T) {
	stmts, err := Parse(`SELECT * FROM t, LATERAL (SELECT 1) x`)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	sel := stmts.Items[0].(*nodes.RawStmt).Stmt.(*nodes.SelectStmt)
	if len(sel.FromClause.Items) < 2 {
		t.Fatalf("expected 2 FROM items, got %d", len(sel.FromClause.Items))
	}
	rs, ok := sel.FromClause.Items[1].(*nodes.RangeSubselect)
	if !ok {
		t.Fatalf("expected RangeSubselect, got %T", sel.FromClause.Items[1])
	}
	if !rs.Lateral {
		t.Fatalf("expected Lateral=true on RangeSubselect, got false")
	}
}

// TestParenLateralRefFuncTableAST asserts that LATERAL f(t.x) produces
// a RangeFunction with Lateral=true. Covers the func_application leg
// of func_table.
func TestParenLateralRefFuncTableAST(t *testing.T) {
	stmts, err := Parse(`SELECT * FROM t, LATERAL f(t.x)`)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	sel := stmts.Items[0].(*nodes.RawStmt).Stmt.(*nodes.SelectStmt)
	rf, ok := sel.FromClause.Items[1].(*nodes.RangeFunction)
	if !ok {
		t.Fatalf("expected RangeFunction, got %T", sel.FromClause.Items[1])
	}
	if !rf.Lateral {
		t.Fatalf("expected Lateral=true on RangeFunction, got false")
	}
	if rf.IsRowsfrom {
		t.Fatalf("expected IsRowsfrom=false (func_application, not ROWS FROM), got true")
	}
}

// TestParenLateralRefRowsFromAST asserts that LATERAL ROWS FROM (...)
// produces a RangeFunction with Lateral=true and IsRowsfrom=true.
func TestParenLateralRefRowsFromAST(t *testing.T) {
	stmts, err := Parse(`SELECT * FROM t, LATERAL ROWS FROM (f(t.x), g(t.y)) WITH ORDINALITY`)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	sel := stmts.Items[0].(*nodes.RawStmt).Stmt.(*nodes.SelectStmt)
	rf, ok := sel.FromClause.Items[1].(*nodes.RangeFunction)
	if !ok {
		t.Fatalf("expected RangeFunction, got %T", sel.FromClause.Items[1])
	}
	if !rf.Lateral {
		t.Fatalf("expected Lateral=true on RangeFunction (ROWS FROM), got false")
	}
	if !rf.IsRowsfrom {
		t.Fatalf("expected IsRowsfrom=true, got false")
	}
	if !rf.Ordinality {
		t.Fatalf("expected Ordinality=true, got false")
	}
}

// TestParenLateralRefXmlTableAST asserts that LATERAL XMLTABLE(...)
// produces a RangeTableFunc with Lateral=true.
func TestParenLateralRefXmlTableAST(t *testing.T) {
	stmts, err := Parse(`SELECT * FROM t, LATERAL XMLTABLE('/root' PASSING x COLUMNS a int PATH 'a') as xt`)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	sel := stmts.Items[0].(*nodes.RawStmt).Stmt.(*nodes.SelectStmt)
	rtf, ok := sel.FromClause.Items[1].(*nodes.RangeTableFunc)
	if !ok {
		t.Fatalf("expected RangeTableFunc, got %T", sel.FromClause.Items[1])
	}
	if !rtf.Lateral {
		t.Fatalf("expected Lateral=true on RangeTableFunc, got false")
	}
}

// TestParenLateralRefJsonTableAST asserts that LATERAL JSON_TABLE(...)
// produces a JsonTable with Lateral=true.
func TestParenLateralRefJsonTableAST(t *testing.T) {
	stmts, err := Parse(`SELECT * FROM t, LATERAL JSON_TABLE(x, '$' COLUMNS(a int PATH '$.a')) as jt`)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	sel := stmts.Items[0].(*nodes.RawStmt).Stmt.(*nodes.SelectStmt)
	jt, ok := sel.FromClause.Items[1].(*nodes.JsonTable)
	if !ok {
		t.Fatalf("expected JsonTable, got %T", sel.FromClause.Items[1])
	}
	if !jt.Lateral {
		t.Fatalf("expected Lateral=true on JsonTable, got false")
	}
}

// TestParenLateralRefFuncTableAliasColdefAST asserts that
// `LATERAL f(t.x) AS fa(a int, b text)` populates the RangeFunction
// with Alias.Aliasname="fa" AND a 2-entry Coldeflist (the typed-column
// form of func_alias_clause: `AS ColId '(' TableFuncElementList ')'`).
// Parse-accept alone (TestParenLateralRefDispatch) is not enough — a
// regression that dropped the alias or failed to populate Coldeflist
// would still parse-accept but produce a structurally wrong AST.
func TestParenLateralRefFuncTableAliasColdefAST(t *testing.T) {
	stmts, err := Parse(`SELECT * FROM t, LATERAL f(t.x) AS fa(a int, b text)`)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	sel := stmts.Items[0].(*nodes.RawStmt).Stmt.(*nodes.SelectStmt)
	rf, ok := sel.FromClause.Items[1].(*nodes.RangeFunction)
	if !ok {
		t.Fatalf("expected RangeFunction, got %T", sel.FromClause.Items[1])
	}
	if !rf.Lateral {
		t.Fatalf("expected Lateral=true, got false")
	}
	if rf.Alias == nil {
		t.Fatalf("expected non-nil Alias for `AS fa(...)`, got nil")
	}
	if rf.Alias.Aliasname != "fa" {
		t.Fatalf("expected Alias.Aliasname=\"fa\", got %q", rf.Alias.Aliasname)
	}
	if rf.Coldeflist == nil {
		t.Fatalf("expected non-nil Coldeflist for typed-column alias form, got nil")
	}
	if len(rf.Coldeflist.Items) != 2 {
		t.Fatalf("expected 2 Coldeflist items, got %d", len(rf.Coldeflist.Items))
	}
}

// TestParenLateralRefFuncTableOrdinalityAST asserts that
// `LATERAL f(t.x) WITH ORDINALITY` sets Ordinality=true on the
// RangeFunction (func_table's opt_ordinality production,
// gram.y:13730). Mirrors the ROWS FROM ordinality assertion above but
// for the func_application arm.
func TestParenLateralRefFuncTableOrdinalityAST(t *testing.T) {
	stmts, err := Parse(`SELECT * FROM t, LATERAL f(t.x) WITH ORDINALITY`)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	sel := stmts.Items[0].(*nodes.RawStmt).Stmt.(*nodes.SelectStmt)
	rf, ok := sel.FromClause.Items[1].(*nodes.RangeFunction)
	if !ok {
		t.Fatalf("expected RangeFunction, got %T", sel.FromClause.Items[1])
	}
	if !rf.Lateral {
		t.Fatalf("expected Lateral=true, got false")
	}
	if !rf.Ordinality {
		t.Fatalf("expected Ordinality=true on LATERAL f(...) WITH ORDINALITY, got false")
	}
	if rf.IsRowsfrom {
		t.Fatalf("expected IsRowsfrom=false (func_application, not ROWS FROM), got true")
	}
}
