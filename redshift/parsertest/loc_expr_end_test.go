package parsertest

import (
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
)

// Expression nodes must end at their own last token, not at the start of the
// token that follows them. Every case here puts whitespace between the
// expression and the next token so a parser that records the next token's
// start would include the trailing space.
func TestLocExprEndsAtLastToken(t *testing.T) {
	whereExpr := func(s *nodes.SelectStmt) nodes.Node { return s.WhereClause }
	target := func(i int) func(s *nodes.SelectStmt) nodes.Node {
		return func(s *nodes.SelectStmt) nodes.Node {
			return s.TargetList.Items[i].(*nodes.ResTarget).Val
		}
	}
	cases := []struct {
		name string
		sql  string
		pick func(s *nodes.SelectStmt) nodes.Node
		want string
	}{
		{"A_Expr before OR", "SELECT 1 FROM t WHERE a = NULL OR b = 1",
			func(s *nodes.SelectStmt) nodes.Node { return s.WhereClause.(*nodes.BoolExpr).Args.Items[0] }, "a = NULL"},
		{"A_Expr last in WHERE", "SELECT 1 FROM t WHERE a = NULL OR b = 1",
			func(s *nodes.SelectStmt) nodes.Node { return s.WhereClause.(*nodes.BoolExpr).Args.Items[1] }, "b = 1"},
		{"BoolExpr OR", "SELECT 1 FROM t WHERE a = NULL OR b = 1 ORDER BY a", whereExpr, "a = NULL OR b = 1"},
		{"A_Expr before semicolon", "SELECT 1 FROM t WHERE a = NULL;", whereExpr, "a = NULL"},
		{"A_Expr before newline and comment", "SELECT 1 FROM t WHERE a = NULL \n-- c\nORDER BY a", whereExpr, "a = NULL"},
		{"A_Expr in target list", "SELECT a + 1 , b FROM t", target(0), "a + 1"},
		{"NullTest", "SELECT a IS NULL , b FROM t", target(0), "a IS NULL"},
		{"NullTest NOT NULL in WHERE", "SELECT 1 FROM t WHERE a IS NOT NULL   ORDER BY a", whereExpr, "a IS NOT NULL"},
		{"TypeCast", "SELECT a::int , b FROM t", target(0), "a::int"},
		{"TypeCast CAST()", "SELECT CAST(a AS int) , b FROM t", target(0), "CAST(a AS int)"},
		{"FuncCall", "SELECT f(a) , b FROM t", target(0), "f(a)"},
		{"FuncCall in WHERE", "SELECT 1 FROM t WHERE f(a) \n", whereExpr, "f(a)"},
		{"ColumnRef", "SELECT a , b FROM t", target(0), "a"},
		{"A_Const", "SELECT 1 , b FROM t", target(0), "1"},
		{"BoolExpr NOT", "SELECT NOT a , b FROM t", target(0), "NOT a"},
		{"unary minus", "SELECT - a , b FROM t", target(0), "- a"},
		{"prefix operator", "SELECT @ a , b FROM t", target(0), "@ a"},
		{"prefix operator on b_expr", "SELECT 1 FROM t WHERE x BETWEEN @ a AND b   ORDER BY x",
			func(s *nodes.SelectStmt) nodes.Node {
				return s.WhereClause.(*nodes.A_Expr).Rexpr.(*nodes.List).Items[0]
			}, "@ a"},
		{"parenthesized expression", "SELECT (a + 1) , b FROM t", target(0), "(a + 1)"},
		{"SubLink", "SELECT 1 FROM t WHERE a IN (SELECT 1) \n ORDER BY a", whereExpr, "a IN (SELECT 1)"},
		{"CaseExpr", "SELECT CASE WHEN a THEN 1 END , b FROM t", target(0), "CASE WHEN a THEN 1 END"},
		{"SortBy", "SELECT a FROM t ORDER BY a DESC , b",
			func(s *nodes.SelectStmt) nodes.Node { return s.SortClause.Items[0] }, "a DESC"},
		{"WindowDef by name", "SELECT sum(a) OVER w , b FROM t WINDOW w AS ()",
			func(s *nodes.SelectStmt) nodes.Node { return target(0)(s).(*nodes.FuncCall).Over }, "w"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			list, err := parse(tc.sql)
			if err != nil {
				t.Fatal(err)
			}
			n := tc.pick(list.Items[0].(*nodes.SelectStmt))
			loc := nodes.NodeLoc(n)
			if loc == nodes.NoLoc() {
				t.Fatalf("%T has no Loc", n)
			}
			if got := tc.sql[loc.Start:loc.End]; got != tc.want {
				t.Errorf("%T text = %q (Loc %+v), want %q", n, got, loc, tc.want)
			}
		})
	}
}
