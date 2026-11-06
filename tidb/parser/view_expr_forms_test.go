package parser

import (
	"testing"

	ast "github.com/bytebase/omni/tidb/ast"
)

// Expression forms a TiDB view body holds, found while writing its restored
// text (deparse.Restore): SOUNDS LIKE, N'...', TRIM with a direction and no
// remstr, CHAR ... USING, and the (a, b) row.
//
// The verdicts were measured against pingcap/tidb:v8.5.5 through PREPARE on
// 2026-11-06. TestViewExprFormsMatchEngine holds them to the server, and
// TestViewExprFormTrees holds the tree each accepted form builds.

var viewExprFormCases = []rowFormCase{
	// TiDB's grammar has no SOUNDS LIKE: sounds is an identifier.
	{`select a sounds like b from t`, false},
	{`select a not sounds like b from t`, false},
	{`select a sounds from t`, true},

	// N or n directly before a quote is the introducer _utf8
	// (startWithNn, pkg/parser/lexer.go), which takes one string.
	{`select N'nat'`, true},
	{`select n'nat'`, true},
	{`select N''`, true},
	{`select N'x' 'y'`, true},
	{`select N'x' 'y' 'z'`, false},
	{`select N'x' as y`, true},
	{`select N'x' collate utf8_bin`, true},
	{`select N 'x' from t`, true},
	{`select N 'x' as y from t`, false},
	{"select `N`'x' from t", true},
	{`select N"x" from t`, true},
	{`select N"x" as y from t`, false},
	{`select _utf8'x' 'y'`, true},
	{`select _utf8'x' 'y' as z`, false},

	// TRIM takes a direction with or without remstr, and FROM after either
	// (parser.y:8336-8366).
	{`select trim(both from c) from t`, true},
	{`select trim(leading from c) from t`, true},
	{`select trim(trailing from c) from t`, true},
	{`select trim(both 'x' from c) from t`, true},
	{`select trim('x' from c) from t`, true},
	{`select trim(c) from t`, true},
	{`select trim(leading 'x')`, false},
	{`select trim(both)`, false},
	{`select trim(both from)`, false},
	{`select trim(both from 'x' from 'y')`, false},

	// CHAR takes an expression list and an optional USING CharsetName
	// (parser.y:8184-8201).
	{`select char(65 using utf8mb4)`, true},
	{`select char(65, 66 using 'latin1')`, true},
	{`select char(65 using binary)`, true},
	{"select char(65 using `utf8`)", true},
	{`select char(65 using utf8mb3)`, true},
	{`select char(65)`, true},
	{`select char()`, false},
	{`select char(using utf8)`, false},
	{`select char(65 using)`, false},
	{`select char(65 using utf8mb4 collate utf8mb4_bin)`, false},
	{`select char(65 using utf8mb4, 66)`, false},

	// '(' ExpressionList ',' Expression ')' is a row (parser.y:7925).
	{`select (a, b) in ((1, 2)) from t`, true},
	{`select (a, b) = (1, 2) from t`, true},
	{`select row(a, b) < (1, 2) from t`, true},
}

// TestViewExprForms holds the parser to the verdicts above.
func TestViewExprForms(t *testing.T) {
	for _, c := range viewExprFormCases {
		t.Run(c.sql, func(t *testing.T) {
			list, err := Parse(c.sql)
			if got := err == nil; got != c.accept {
				t.Fatalf("Parse(%q) accepts = %v, the engine answered %v (err: %v)", c.sql, got, c.accept, err)
			}
			if err == nil && len(list.Items) != 1 {
				t.Errorf("Parse(%q) read %d statements, want 1", c.sql, len(list.Items))
			}
		})
	}
}

// TestViewExprFormTrees pins what each form carries: the restored view text
// is written from these fields.
func TestViewExprFormTrees(t *testing.T) {
	for _, tc := range []struct {
		sql, want string
	}{
		{`select N'nat'`, `{SELECT :loc 0 :targets ({STRING_LIT :val "nat" :loc 7 :charset _utf8})}`},
		{`select N'x' 'y'`, `{SELECT :loc 0 :targets ({RES_TARGET :loc 7 :name y :val {STRING_LIT :val "x" :loc 7 :charset _utf8}})}`},
		{`select N 'x'`, `{SELECT :loc 0 :targets ({RES_TARGET :loc 7 :name x :val {COLREF :loc 7 :col N}})}`},
		{`select trim(both from c)`, `{SELECT :loc 0 :targets ({FUNCCALL :name TRIM_BOTH :loc 7 :args ({COLREF :loc 22 :col c})})}`},
		{`select trim(leading 'x' from c)`, `{SELECT :loc 0 :targets ({FUNCCALL :name TRIM_LEADING :loc 7 :args ({STRING_LIT :val "x" :loc 20} {COLREF :loc 29 :col c})})}`},
		{`select char(65, 66 using 'latin1')`, `{SELECT :loc 0 :targets ({FUNCCALL :name CHAR :loc 7 :using latin1 :args ({INT_LIT :val 65 :loc 12} {INT_LIT :val 66 :loc 16})})}`},
		{`select char(65)`, `{SELECT :loc 0 :targets ({FUNCCALL :name CHAR :loc 7 :args ({INT_LIT :val 65 :loc 12})})}`},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			list, err := Parse(tc.sql)
			if err != nil {
				t.Fatal(err)
			}
			if got := ast.NodeToString(list.Items[0]); got != tc.want {
				t.Errorf("tree\n got %s\nwant %s", got, tc.want)
			}
		})
	}
}

// TestViewExprFormLocs: a form's Loc ends at its last token.
func TestViewExprFormLocs(t *testing.T) {
	for _, tc := range []struct{ sql, text string }{
		{`select N'nat' from t`, `N'nat'`},
		{`select trim(both from c) from t`, `trim(both from c)`},
		{`select char(65 using utf8mb4) from t`, `char(65 using utf8mb4)`},
	} {
		list, err := Parse(tc.sql)
		if err != nil {
			t.Fatalf("%s: %v", tc.sql, err)
		}
		e := list.Items[0].(*ast.SelectStmt).TargetList[0]
		if rt, ok := e.(*ast.ResTarget); ok {
			e = rt.Val
		}
		var loc ast.Loc
		switch n := e.(type) {
		case *ast.StringLit:
			loc = n.Loc
		case *ast.FuncCallExpr:
			loc = n.Loc
		default:
			t.Fatalf("%s: unexpected %T", tc.sql, e)
		}
		if got := tc.sql[loc.Start:loc.End]; got != tc.text {
			t.Errorf("%s: Loc covers %q, want %q", tc.sql, got, tc.text)
		}
	}
}
