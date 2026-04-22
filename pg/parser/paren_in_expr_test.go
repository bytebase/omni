package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// TestParenInExprDispatch covers pg-paren-dispatch §1.1: parseInExpr's
// choice between IN (expr_list) and IN (select_with_parens).
//
// PG grammar (gram.y:14973-14998, in_expr):
//
//	in_expr:
//	    select_with_parens
//	    | '(' expr_list ')'
//
// Both branches begin with '(', and select_with_parens admits
// '(' select_with_parens ')' (nested parens around a subquery), so a
// 1-token peek after '(' cannot reliably distinguish the branches.
// Before this fix, omni's peek routed IN ((SELECT 1) UNION (SELECT 2))
// to the expr_list branch and failed on UNION; now the parser walks
// past consecutive '(' tokens to check for a SELECT / VALUES / WITH /
// TABLE lead before committing to a branch.
func TestParenInExprDispatch(t *testing.T) {
	cases := []struct {
		sql       string
		wantParse bool
	}{
		// --- expr_list branch -------------------------------------------------
		{`SELECT * FROM t WHERE x IN (1, 2, 3)`, true},
		{`SELECT * FROM t WHERE x IN (a, b, c)`, true},
		{`SELECT * FROM t WHERE x IN ('a', 'b')`, true},
		{`SELECT * FROM t WHERE x IN (1)`, true},
		// --- subquery branch --------------------------------------------------
		{`SELECT * FROM t WHERE x IN (SELECT y FROM u)`, true},
		{`SELECT * FROM t WHERE x IN (SELECT 1)`, true},
		{`SELECT * FROM t WHERE x IN (SELECT 1 UNION SELECT 2)`, true},
		// select_with_parens operand case — the scenario Phase 0's
		// parseSelectWithParens left-factoring makes "just work" once
		// we scan past the outer '(' to find SELECT.
		{`SELECT * FROM t WHERE x IN ((SELECT 1) UNION (SELECT 2))`, true},
		{`SELECT * FROM t WHERE x IN (VALUES (1), (2))`, true},
		{`SELECT * FROM t WHERE x IN (WITH cte AS (SELECT 1) SELECT * FROM cte)`, true},
		{`SELECT * FROM t WHERE x IN (TABLE foo)`, true},
		// --- rejected empty list ---------------------------------------------
		// PG rejects IN () with a syntax error — in_expr has no empty
		// production. Before this fix omni accepted it silently.
		{`SELECT * FROM t WHERE x IN ()`, false},
		// --- row-constructor LHS ---------------------------------------------
		{`SELECT * FROM t WHERE (x, y) IN (SELECT a, b FROM u)`, true},
		{`SELECT * FROM t WHERE (x, y) IN ((1, 2), (3, 4))`, true},
		// PG's parser accepts this; arity mismatch is caught in
		// analyze, not parse.
		{`SELECT * FROM t WHERE (x, y) IN (SELECT a FROM u)`, true},
		// --- NOT IN variants --------------------------------------------------
		{`SELECT * FROM t WHERE x NOT IN (1, 2, 3)`, true},
		{`SELECT * FROM t WHERE x NOT IN (SELECT y FROM u)`, true},
		{`SELECT * FROM t WHERE x NOT IN (SELECT 1 UNION SELECT 2)`, true},
		{`SELECT * FROM t WHERE (x, y) NOT IN (SELECT a, b FROM u)`, true},
		// --- regression: list of parenthesized sub-select expressions ---
		// PG parses this as IN (expr_list) with two scalar-subquery
		// elements (`(select 1)` and `(select 2)` each parse as an
		// a_expr via the select_with_parens production on the expr
		// side). Before the lookahead reworked to check depth-0 ',',
		// the parser mis-routed this to the subquery branch. Upstream
		// shape from pg regress/partition_prune.sql.
		{`SELECT * FROM rangep WHERE b IN ((SELECT 1), (SELECT 2))`, true},
		// Extra parens around a single subquery — select_with_parens
		// admits '(' select_with_parens ')' recursion. Walk past the
		// outer '(' to find SELECT; end of IN paren is a single element.
		{`SELECT * FROM t WHERE x IN ((SELECT 1))`, true},
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

// TestParenInExprSubqueryAST asserts that IN (SELECT ...) produces a
// SubLink with SubLinkType=ANY_SUBLINK.  This is the raw-parse shape
// PG's gram.y emits; the analyze phase later transforms it to
// ScalarArrayOpExpr (or keeps it as SubLink for the subquery form),
// but omni's parser is raw-parse only, so we assert the raw-parse
// shape.
func TestParenInExprSubqueryAST(t *testing.T) {
	stmts, err := Parse(`SELECT * FROM t WHERE x IN (SELECT y FROM u)`)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	sel := stmts.Items[0].(*nodes.RawStmt).Stmt.(*nodes.SelectStmt)
	sub, ok := sel.WhereClause.(*nodes.SubLink)
	if !ok {
		t.Fatalf("expected SubLink in WhereClause, got %T", sel.WhereClause)
	}
	if sub.SubLinkType != int(nodes.ANY_SUBLINK) {
		t.Fatalf("expected SubLinkType=ANY_SUBLINK (%d), got %d", int(nodes.ANY_SUBLINK), sub.SubLinkType)
	}
	if sub.Testexpr == nil {
		t.Fatalf("expected Testexpr to be the LHS ColumnRef, got nil")
	}
	if _, ok := sub.Subselect.(*nodes.SelectStmt); !ok {
		t.Fatalf("expected Subselect to be SelectStmt, got %T", sub.Subselect)
	}
}

// TestParenInExprExprListAST asserts that IN (expr_list) produces an
// A_Expr with Kind=AEXPR_IN and operator name "=".  PG's post-analyze
// tree turns this into ScalarArrayOpExpr (useOr=true for IN), but at
// the parser level we carry A_Expr AEXPR_IN with name "=".
func TestParenInExprExprListAST(t *testing.T) {
	stmts, err := Parse(`SELECT * FROM t WHERE x IN (1, 2, 3)`)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	sel := stmts.Items[0].(*nodes.RawStmt).Stmt.(*nodes.SelectStmt)
	ae, ok := sel.WhereClause.(*nodes.A_Expr)
	if !ok {
		t.Fatalf("expected A_Expr in WhereClause, got %T", sel.WhereClause)
	}
	if ae.Kind != nodes.AEXPR_IN {
		t.Fatalf("expected Kind=AEXPR_IN, got %v", ae.Kind)
	}
	if ae.Name == nil || len(ae.Name.Items) != 1 {
		t.Fatalf("expected single-item Name, got %+v", ae.Name)
	}
	if s, ok := ae.Name.Items[0].(*nodes.String); !ok || s.Str != "=" {
		t.Fatalf("expected Name[0]='=', got %+v", ae.Name.Items[0])
	}
	rlist, ok := ae.Rexpr.(*nodes.List)
	if !ok {
		t.Fatalf("expected Rexpr to be *List, got %T", ae.Rexpr)
	}
	if len(rlist.Items) != 3 {
		t.Fatalf("expected 3 Rexpr items, got %d", len(rlist.Items))
	}
}

// TestParenInExprNotInSubqueryAST asserts that NOT IN (SELECT ...)
// wraps the ANY_SUBLINK in a BoolExpr NOT_EXPR, matching PG's gram.y
// yacc action (see gram.y:14991-14997 a_expr_in_expr NOT case).
func TestParenInExprNotInSubqueryAST(t *testing.T) {
	stmts, err := Parse(`SELECT * FROM t WHERE x NOT IN (SELECT y FROM u)`)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	sel := stmts.Items[0].(*nodes.RawStmt).Stmt.(*nodes.SelectStmt)
	be, ok := sel.WhereClause.(*nodes.BoolExpr)
	if !ok {
		t.Fatalf("expected BoolExpr in WhereClause, got %T", sel.WhereClause)
	}
	if be.Boolop != nodes.NOT_EXPR {
		t.Fatalf("expected Boolop=NOT_EXPR, got %v", be.Boolop)
	}
	if be.Args == nil || len(be.Args.Items) != 1 {
		t.Fatalf("expected 1 arg, got %+v", be.Args)
	}
	sub, ok := be.Args.Items[0].(*nodes.SubLink)
	if !ok {
		t.Fatalf("expected SubLink inside NOT BoolExpr, got %T", be.Args.Items[0])
	}
	if sub.SubLinkType != int(nodes.ANY_SUBLINK) {
		t.Fatalf("expected ANY_SUBLINK (%d), got %d", int(nodes.ANY_SUBLINK), sub.SubLinkType)
	}
}

// TestParenInExprNotInExprListAST asserts that NOT IN (expr_list)
// produces an A_Expr with Kind=AEXPR_IN and operator name "<>".
// PG's post-analyze tree turns this into ScalarArrayOpExpr with
// useOr=false and operator "<>"; at the parser level we carry the
// A_Expr with name "<>".
func TestParenInExprNotInExprListAST(t *testing.T) {
	stmts, err := Parse(`SELECT * FROM t WHERE x NOT IN (1, 2, 3)`)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	sel := stmts.Items[0].(*nodes.RawStmt).Stmt.(*nodes.SelectStmt)
	ae, ok := sel.WhereClause.(*nodes.A_Expr)
	if !ok {
		t.Fatalf("expected A_Expr in WhereClause, got %T", sel.WhereClause)
	}
	if ae.Kind != nodes.AEXPR_IN {
		t.Fatalf("expected Kind=AEXPR_IN, got %v", ae.Kind)
	}
	if s, ok := ae.Name.Items[0].(*nodes.String); !ok || s.Str != "<>" {
		t.Fatalf("expected Name[0]='<>', got %+v", ae.Name.Items[0])
	}
}

// Out-of-scope marker (§1.1): x = ANY/SOME/ALL (...) is not a '('
// dispatch ambiguity — it goes through sub_type (gram.y:14976-14998)
// which is keyword-led and routed by the operator. Tracked separately
// by the audit; not part of parseInExpr.
