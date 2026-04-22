package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// TestParenArrayExprDispatch covers pg-paren-dispatch §1.3:
// parseArrayCExpr's choice between the ARRAY array-constructor form
// and the ARRAY sublink form.
//
// PG grammar (gram.y:15440-15459 + gram.y:16583-16595):
//
//	c_expr:
//	    ARRAY select_with_parens   → SubLink{subLinkType=ARRAY_SUBLINK}
//	  | ARRAY array_expr           → A_ArrayExpr
//
//	array_expr:
//	    '[' expr_list ']'
//	  | '[' array_expr_list ']'
//	  | '[' ']'
//
// The two productions diverge on the first token after ARRAY: `[`
// commits to array_expr, `(` to select_with_parens. omni's
// parseArrayCExpr performs this T1 peek at expr.go:1741 ('[') and
// expr.go:1754 ('(').
//
// Notes on the "SELECT ARRAY[]" case: PG's grammar accepts '[' ']' as
// a valid array_expr at parse time; the "cannot determine type of
// empty array" error is produced later in transformAArrayExpr (post-
// analyze), not by the parser. Since omni's parser is raw-parse only
// (no analyze phase), we mirror PG's parser acceptance and document
// the semantic rejection as out-of-scope for parser-level testing.
// With an explicit cast (`ARRAY[]::int[]`), both PG and the analyze
// phase accept.
func TestParenArrayExprDispatch(t *testing.T) {
	cases := []struct {
		sql       string
		wantParse bool
	}{
		// --- array_expr branch (ARRAY[...]) ---------------------------------
		{`SELECT ARRAY[1, 2, 3]`, true},
		{`SELECT ARRAY['a', 'b']`, true},
		// Empty with cast: array_expr produces A_ArrayExpr{Elements: nil},
		// the ::int[] cast then resolves the element type.
		{`SELECT ARRAY[]::int[]`, true},
		// Empty without cast: parser accepts (array_expr: '[' ']'). PG
		// rejects this semantically in transformAArrayExpr ("cannot
		// determine type of empty array") — a post-parse check omni's
		// raw parser does not perform. Keeping this as parse-accept
		// tracks PG's gram.y exactly.
		{`SELECT ARRAY[]`, true},
		// Nested array constructor: array_expr_list → array_expr.
		// At the raw-parse level this produces an outer A_ArrayExpr
		// whose Elements are two inner A_ArrayExpr nodes. Analyze
		// would flag Multidims=true on the resulting ArrayExpr.
		{`SELECT ARRAY[[1,2],[3,4]]`, true},
		// Array concat in expr context: each ARRAY[...] parses as
		// A_ArrayExpr, the `||` is an A_Expr operator between them.
		{`SELECT ARRAY[1, 2] || ARRAY[3]`, true},

		// --- ARRAY sublink branch (ARRAY(...)) -------------------------------
		{`SELECT ARRAY(SELECT 1)`, true},
		{`SELECT ARRAY(SELECT x FROM t)`, true},
		{`SELECT ARRAY(SELECT DISTINCT x FROM t ORDER BY x)`, true},
		{`SELECT ARRAY(VALUES (1), (2))`, true},
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

// TestParenArrayExprConstructorAST asserts ARRAY[1,2] produces
// A_ArrayExpr with a 2-element expr_list. At the raw-parse level
// there is no Multidims field on A_ArrayExpr — PG's Multidims is
// set during transformArrayExpr on the analyzed ArrayExpr. We assert
// the parser-level shape.
func TestParenArrayExprConstructorAST(t *testing.T) {
	stmts, err := Parse(`SELECT ARRAY[1, 2]`)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	sel := stmts.Items[0].(*nodes.RawStmt).Stmt.(*nodes.SelectStmt)
	tl := sel.TargetList.Items[0].(*nodes.ResTarget)
	arr, ok := tl.Val.(*nodes.A_ArrayExpr)
	if !ok {
		t.Fatalf("expected *A_ArrayExpr, got %T", tl.Val)
	}
	if arr.Elements == nil || len(arr.Elements.Items) != 2 {
		t.Fatalf("expected 2 elements, got %+v", arr.Elements)
	}
}

// TestParenArrayExprNestedConstructorAST asserts ARRAY[[1,2],[3,4]]
// produces an outer A_ArrayExpr whose Elements are two inner
// A_ArrayExpr nodes — i.e. the raw-parse shape that analyze will
// later flag as Multidims=true.
func TestParenArrayExprNestedConstructorAST(t *testing.T) {
	stmts, err := Parse(`SELECT ARRAY[[1,2],[3,4]]`)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	sel := stmts.Items[0].(*nodes.RawStmt).Stmt.(*nodes.SelectStmt)
	tl := sel.TargetList.Items[0].(*nodes.ResTarget)
	arr, ok := tl.Val.(*nodes.A_ArrayExpr)
	if !ok {
		t.Fatalf("expected outer *A_ArrayExpr, got %T", tl.Val)
	}
	if arr.Elements == nil || len(arr.Elements.Items) != 2 {
		t.Fatalf("expected 2 outer elements, got %+v", arr.Elements)
	}
	for i, item := range arr.Elements.Items {
		inner, ok := item.(*nodes.A_ArrayExpr)
		if !ok {
			t.Fatalf("expected inner element %d to be *A_ArrayExpr, got %T", i, item)
		}
		if inner.Elements == nil || len(inner.Elements.Items) != 2 {
			t.Fatalf("expected inner element %d to have 2 items, got %+v", i, inner.Elements)
		}
	}
}

// TestParenArrayExprSubqueryContract covers pg-paren-dispatch §1.4: the
// content contract of `ARRAY '(' ... ')'`. Per PG 17 gram.y:15440-15451,
// the production is `ARRAY select_with_parens`, which restricts the
// paren content to `select_with_parens` — i.e. another wrapped
// select_with_parens or a select_no_parens. select_no_parens admits:
//
//   - simple_select (SELECT, VALUES, TABLE form)
//   - select_clause sort_clause / select_clause for_locking / etc.
//   - with_clause select_clause ...
//
// It does NOT admit `ROWS FROM (...)` (func_table-only), bare
// expressions like `ARRAY(1)`, or empty parens `ARRAY()`. omni's
// parseArrayCExpr defers to parseSelectStmtForExpr (= parseSelectNoParens)
// which returns a typed-nil *SelectStmt for the latter cases; the
// fix at expr.go:1610 (T7 post-parse content contract check)
// rejects when the subquery comes back nil.
//
// These six scenarios provide the empirical proof that locks in
// PAREN_AUDIT row expr.go:1610 from `aligned: unclear` → `aligned: yes`.
func TestParenArrayExprSubqueryContract(t *testing.T) {
	cases := []struct {
		sql       string
		wantParse bool
		why       string
	}{
		// --- ACCEPT: select_no_parens shapes ---------------------------------
		// TABLE foo is a simple_select form (gram.y:13183-13186):
		//   simple_select: TABLE relation_expr
		// Reachable via parseSimpleSelectLeaf's TABLE arm.
		{
			`SELECT ARRAY(TABLE foo)`,
			true,
			"TABLE is simple_select → select_no_parens → select_with_parens",
		},
		// WITH ... SELECT is select_no_parens's with_clause+select_clause
		// alternative (gram.y:13165-13168). parseSelectNoParens picks up
		// the WITH lead, then parses select_clause.
		{
			`SELECT ARRAY(WITH cte AS (SELECT 1) SELECT * FROM cte)`,
			true,
			"WITH + select_clause is a select_no_parens form",
		},
		// VALUES is simple_select (gram.y:13180): VALUES_P ... .
		{
			`SELECT ARRAY(VALUES (1), (2))`,
			true,
			"VALUES is simple_select → select_no_parens",
		},

		// --- REJECT: not a select_no_parens lead -----------------------------
		// ROWS FROM is reachable only through func_table (gram.y:13468);
		// it is NOT in select_no_parens. parseSelectStmtForExpr returns
		// nil on the unexpected ROWS lead, then expect(')') sees ROWS
		// and emits a syntax error.
		{
			`SELECT ARRAY(ROWS FROM (f(1)))`,
			false,
			"ROWS FROM is func_table-only, not in select_no_parens",
		},
		// Bare expression: parseSimpleSelectLeaf rejects integer-literal
		// lead (returns nil). The T7 post-parse check then errors before
		// expect(')').
		{
			`SELECT ARRAY(1)`,
			false,
			"bare expr is not a select_no_parens",
		},
		// Empty parens: parseSimpleSelectLeaf returns nil for `)` lead.
		// Without the T7 nil check, omni would silently produce a
		// SubLink with Subselect=nil (which is what it did pre-1.4).
		{
			`SELECT ARRAY()`,
			false,
			"empty parens are not a select_no_parens",
		},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			_, err := Parse(tc.sql)
			if tc.wantParse && err != nil {
				t.Fatalf("expected parse success (%s), got error: %v", tc.why, err)
			}
			if !tc.wantParse && err == nil {
				t.Fatalf("expected parse error (%s), got nil", tc.why)
			}
		})
	}
}

// TestParenArrayExprSublinkAST asserts ARRAY(SELECT 1) produces a
// SubLink with SubLinkType=ARRAY_SUBLINK, matching PG's gram.y
// yacc action at gram.y:15440-15451.
func TestParenArrayExprSublinkAST(t *testing.T) {
	stmts, err := Parse(`SELECT ARRAY(SELECT 1)`)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	sel := stmts.Items[0].(*nodes.RawStmt).Stmt.(*nodes.SelectStmt)
	tl := sel.TargetList.Items[0].(*nodes.ResTarget)
	sub, ok := tl.Val.(*nodes.SubLink)
	if !ok {
		t.Fatalf("expected *SubLink, got %T", tl.Val)
	}
	if sub.SubLinkType != int(nodes.ARRAY_SUBLINK) {
		t.Fatalf("expected SubLinkType=ARRAY_SUBLINK (%d), got %d",
			int(nodes.ARRAY_SUBLINK), sub.SubLinkType)
	}
	if _, ok := sub.Subselect.(*nodes.SelectStmt); !ok {
		t.Fatalf("expected Subselect to be *SelectStmt, got %T", sub.Subselect)
	}
}
