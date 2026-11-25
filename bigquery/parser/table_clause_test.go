package parser

import (
	"testing"

	"github.com/bytebase/omni/bigquery/ast"
)

// Tests for the query-side table_clause: `TABLE <path|tvf> …` as a query primary
// (query_primary: table_clause_reserved, googlesql.tm:5674).
//
// CORRECTNESS BASIS: unlike the other statement families of the same compat
// rows, the Spanner emulator PARSES a TABLE statement in full and only then
// answers "TABLE statements are not supported" — it rejects `TABLE a zzz`,
// `TABLE a AS x`, `TABLE a PIVOT(…)` and `TABLE a(` on syntax, and accepts the
// tvf argument shapes. So the `TABLE …` rows of testdata/compat.tsv are a real
// parse verdict, and the cases below mirror the emulator's answers.
//
// The rule is shared with the pipe set operation's TABLE operand (pipe.go's
// parseTableClauseOperand); only the tvf form takes an alias and PIVOT/UNPIVOT.

// tableClauseOf parses sql and returns the TableClause body of its query.
func tableClauseOf(t *testing.T, sql string) *ast.TableClause {
	t.Helper()
	n := parseOneStmt(t, sql)
	q, ok := n.(*ast.QueryStmt)
	if !ok {
		t.Fatalf("parseOneStmt(%q): got %T, want *ast.QueryStmt", sql, n)
	}
	tc, ok := q.Body.(*ast.TableClause)
	if !ok {
		t.Fatalf("parseOneStmt(%q): body is %T, want *ast.TableClause", sql, q.Body)
	}
	return tc
}

func TestParseTableClause_Path(t *testing.T) {
	tc := tableClauseOf(t, "TABLE t")
	if tc.Table == nil || tc.Table.Path == nil {
		t.Fatalf("Table = %+v, want a table path", tc.Table)
	}
	if got := len(tc.Table.Path.Parts); got != 1 || tc.Table.Path.Parts[0] != "t" {
		t.Errorf("path = %+v, want [t]", tc.Table.Path.Parts)
	}
	if tc.Where != nil {
		t.Errorf("Where = %+v, want nil", tc.Where)
	}
	if tc.Loc.Start != 0 || tc.Loc.End != len("TABLE t") {
		t.Errorf("Loc = %+v, want (0, %d)", tc.Loc, len("TABLE t"))
	}

	// A dotted path, and a backtick-quoted one.
	if got := tableClauseOf(t, "TABLE a.b.c").Table.Path.Parts; len(got) != 3 {
		t.Errorf("path parts = %+v, want 3", got)
	}
	if got := tableClauseOf(t, "TABLE `a-b`").Table.Path.Parts; len(got) != 1 || got[0] != "a-b" {
		t.Errorf("path parts = %+v, want [a-b]", got)
	}

	// The query's trailing modifiers bind to the whole query, and the TABLE
	// primary takes a set operation like any other query_primary.
	n := parseOneStmt(t, "TABLE t ORDER BY a LIMIT 1").(*ast.QueryStmt)
	if len(n.OrderBy) != 1 || n.Limit == nil {
		t.Errorf("ORDER BY / LIMIT = (%+v, %+v), want one of each", n.OrderBy, n.Limit)
	}
	so := parseOneStmt(t, "TABLE a UNION ALL TABLE b").(*ast.QueryStmt)
	if _, ok := so.Body.(*ast.SetOperation); !ok {
		t.Errorf("body = %T, want *ast.SetOperation", so.Body)
	}
	// A TABLE primary is a query_primary, so a parenthesized query, a FROM source
	// and an expression subquery can all hold one (the emulator parses each of
	// these and stops at "TABLE statements are not supported").
	expr := parseOneStmt(t, "SELECT (TABLE t)").(*ast.QueryStmt).Body.(*ast.SelectStmt).Items[0].Expr
	sq, ok := expr.(*ast.SubqueryExpr)
	if !ok {
		t.Fatalf("select item = %T, want *ast.SubqueryExpr", expr)
	}
	if q, ok := sq.Query.(*ast.QueryStmt); !ok {
		t.Errorf("subquery = %T, want *ast.QueryStmt", sq.Query)
	} else if _, ok := q.Body.(*ast.TableClause); !ok {
		t.Errorf("subquery body = %T, want *ast.TableClause", q.Body)
	}
	from := parseOneStmt(t, "SELECT * FROM (TABLE t) AS x").(*ast.QueryStmt)
	sel := from.Body.(*ast.SelectStmt)
	te, ok := sel.From[0].(*ast.TableExpr)
	if !ok || te.Subquery == nil {
		t.Fatalf("FROM = %+v, want a table_subquery", sel.From[0])
	}
	sub, ok := te.Subquery.(*ast.QueryStmt)
	if !ok {
		t.Fatalf("subquery = %T, want *ast.QueryStmt", te.Subquery)
	}
	if _, ok := sub.Body.(*ast.TableClause); !ok {
		t.Errorf("subquery body = %T, want *ast.TableClause", sub.Body)
	}

	// TABLE is non-reserved, so a bare TABLE in an expression position is the
	// keyword as an identifier, not a query head: `(TABLE)` and `(TABLE.x)` stay
	// expressions, and `(TABLE t)` is the query.
	exprTable := parseOneStmt(t, "SELECT (TABLE)").(*ast.QueryStmt).Body.(*ast.SelectStmt).Items[0].Expr
	if _, ok := exprTable.(*ast.ParenExpr); !ok {
		t.Errorf("SELECT (TABLE): item = %T, want a parenthesized expression", exprTable)
	}
	if _, ok := parseOneStmt(t, "SELECT (TABLE.x)").(*ast.QueryStmt).Body.(*ast.SelectStmt).Items[0].Expr.(*ast.ParenExpr); !ok {
		t.Error("SELECT (TABLE.x): want a parenthesized expression, not a query")
	}

	// A leading WITH clause takes a TABLE primary too.
	with := parseOneStmt(t, "WITH c AS (SELECT 1 AS n) TABLE c").(*ast.QueryStmt)
	if with.With == nil || len(with.With.CTEs) != 1 {
		t.Errorf("With = %+v, want one CTE", with.With)
	}
}

func TestParseTableClause_TVF(t *testing.T) {
	// tvf: (path_expression | IF) '(' (tvf_argument , …)* ')' hint?
	tc := tableClauseOf(t, "TABLE ds.f(a, NULL)")
	if tc.Table == nil || tc.Table.Func == nil {
		t.Fatalf("Table = %+v, want a tvf", tc.Table)
	}
	if got := len(tc.Table.Func.Args); got != 2 {
		t.Errorf("args = %d, want 2", got)
	}

	// Every tvf_argument alternative the grammar has besides a bare expression.
	tc = tableClauseOf(t, "TABLE f(DESCRIPTOR(a, b), MODEL m, CONNECTION c, INPUT TABLE, x => TABLE u)")
	var kinds []ast.CallArgKind
	for _, arg := range tc.Table.Func.Args {
		switch a := arg.(type) {
		case *ast.CallArg:
			kinds = append(kinds, a.Kind)
			if a.Kind == ast.CallArgDescriptor && len(a.Columns) != 2 {
				t.Errorf("descriptor columns = %+v, want 2", a.Columns)
			}
			if a.Kind == ast.CallArgConnection && (a.Path == nil || a.Path.Parts[0] != "c") {
				t.Errorf("connection = %+v, want c", a)
			}
		case *ast.NamedArg:
			if a.Name != "x" {
				t.Errorf("named argument = %q, want x", a.Name)
			}
			if _, ok := a.Value.(*ast.TableClause); !ok {
				t.Errorf("named argument value = %T, want *ast.TableClause", a.Value)
			}
		}
	}
	want := []ast.CallArgKind{ast.CallArgDescriptor, ast.CallArgModel, ast.CallArgConnection, ast.CallArgInputTable}
	if len(kinds) != len(want) {
		t.Fatalf("call args = %+v, want %+v", kinds, want)
	}
	for i := range want {
		if kinds[i] != want[i] {
			t.Fatalf("call args = %+v, want %+v", kinds, want)
		}
	}

	// `TABLE IF(...)`: IF stands in for the tvf name.
	if tc := tableClauseOf(t, "TABLE IF()"); tc.Table.Func == nil || tc.Table.Func.Name == nil || tc.Table.Func.Name.Parts[0] != "IF" {
		t.Errorf("TABLE IF() = %+v, want a tvf named IF", tc.Table)
	}

	// The tvf form takes an alias (bare or AS) and a PIVOT/UNPIVOT clause.
	tc = tableClauseOf(t, "TABLE f() AS x")
	if tc.Table.Alias != "x" {
		t.Errorf("Alias = %q, want x", tc.Table.Alias)
	}
	tc = tableClauseOf(t, "TABLE f() y")
	if tc.Table.Alias != "y" {
		t.Errorf("Alias = %q, want y", tc.Table.Alias)
	}
	if tc := tableClauseOf(t, "TABLE f() PIVOT(SUM(x) FOR y IN ('a'))"); tc.Table.Pivot == nil {
		t.Error("Pivot = nil, want a PIVOT clause")
	}
	if tc := tableClauseOf(t, "TABLE f() UNPIVOT(x FOR y IN (a))"); tc.Table.Unpivot == nil {
		t.Error("Unpivot = nil, want an UNPIVOT clause")
	}

	// A hint after the argument list.
	if tc := tableClauseOf(t, "TABLE f()@1"); tc.Table.Func == nil {
		t.Error("a hinted tvf call did not parse")
	}
	// named_argument's fourth alternative: `identifier => input_table_argument`.
	tc = tableClauseOf(t, "TABLE f(a => INPUT TABLE)")
	if got := tc.Table.Func.Args; len(got) != 1 {
		t.Fatalf("args = %+v, want 1", got)
	}
	if na, ok := tc.Table.Func.Args[0].(*ast.NamedArg); !ok || na.Name != "a" {
		t.Fatalf("arg = %T %+v, want the named argument a", tc.Table.Func.Args[0], tc.Table.Func.Args[0])
	} else if ca, ok := na.Value.(*ast.CallArg); !ok || ca.Kind != ast.CallArgInputTable {
		t.Errorf("named argument value = %T, want an INPUT TABLE CallArg", na.Value)
	}

	// A bare keyword is also a plain identifier expression argument.
	if tc := tableClauseOf(t, "TABLE f(CONNECTION)"); tc.Table.Func == nil {
		t.Error("a bare CONNECTION argument did not parse")
	}
}

func TestParseTableClause_Rejects(t *testing.T) {
	for _, sql := range []string{
		"TABLE",
		"TABLE a zzz",                   // a bare path takes no alias
		"TABLE a AS x",                  // … not even AS
		"TABLE a PIVOT(x FOR y IN (1))", // only the tvf form takes PIVOT
		"TABLE a WHERE",
		"TABLE a(",
		"TABLE a() a zzz",
		"TABLE a(MODEL a zzz)",
		"TABLE a(a => MODEL a)", // a named argument's value is an expression
		"TABLE a JOIN b",
		"TABLE FUNCTION a", // `TABLE FUNCTION METADATA …` only
		"TABLE a SET AS 'a'",
	} {
		assertReject(t, sql)
	}
}
