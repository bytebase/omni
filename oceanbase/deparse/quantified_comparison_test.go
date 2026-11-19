package deparse

import (
	"testing"

	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/oceanbase/parser"
)

// TestQuantifiedComparison pins that a comparison keeps its quantifier through the
// rewrite and the writer, in OceanBase 4.3.5's stored form. Dropping ALL reads `a > ALL
// (q)` as a scalar comparison against q, and folding NOT into the operator reads `NOT
// a = ANY (q)` as `a <> ANY (q)`; both parse and print without complaint, and a
// view stored through either means something else.
func TestQuantifiedComparison(t *testing.T) {
	for _, c := range []struct{ input, want string }{
		{"a = ANY (SELECT 1)", "(`a` = any (select 1))"},
		{"a = SOME (SELECT 1)", "(`a` = any (select 1))"},
		{"a = ALL (SELECT 1)", "(`a` = all (select 1))"},
		{"a <> ANY (SELECT 1)", "(`a` <> any (select 1))"},
		{"a <> SOME (SELECT 1)", "(`a` <> any (select 1))"},
		{"a <> ALL (SELECT 1)", "(`a` <> all (select 1))"},
		{"a > ANY (SELECT 1)", "(`a` > any (select 1))"},
		{"a > SOME (SELECT 1)", "(`a` > any (select 1))"},
		{"a > ALL (SELECT 1)", "(`a` > all (select 1))"},
		{"a <= ALL (SELECT 1)", "(`a` <= all (select 1))"},
		{"NOT a = ANY (SELECT 1)", "(not((`a` = any (select 1))))"},
		{"NOT a > ALL (SELECT 1)", "(not((`a` > all (select 1))))"},
	} {
		t.Run(c.input, func(t *testing.T) {
			node := parseExpr(t, c.input)
			rewritten := RewriteExpr(node)
			if got := Deparse(rewritten); got != c.want {
				t.Errorf("Deparse(RewriteExpr(%q)) = %q, want %q", c.input, got, c.want)
			}
		})
	}
}

// TestNestedQuantifiedComparison pins that a nested quantified comparison still means
// what was written: its text parses, and printing that parse again gives the same text.
// These writers leave `a > ALL (q)` and `a IN (q)` bare where the grammar allows it, so
// a missing pair of parentheses is either a syntax error or, as in `1 = (a > ALL (q))`
// read back as `(1 = a) > ALL (q)`, a view that means something else.
func TestNestedQuantifiedComparison(t *testing.T) {
	for _, input := range nestedQuantifiedShapes {
		t.Run(input, func(t *testing.T) {
			once := Deparse(RewriteExpr(parseExpr(t, input)))
			if twice := Deparse(RewriteExpr(parseExpr(t, once))); twice != once {
				t.Errorf("Deparse(%q) = %q, which reads back as %q", input, once, twice)
			}
		})
	}
}

var nestedQuantifiedShapes = []string{
	"(a > ALL (SELECT 1)) + 1",
	"1 + (a > ALL (SELECT 1))",
	"1 = (a > ALL (SELECT 1))",
	"(a > ALL (SELECT 1)) = 1",
	"(a = ANY (SELECT 1)) + 1",
	"1 = (a = ANY (SELECT 1))",
	"(a <> ALL (SELECT 1)) + 1",
	"1 = (a <> ALL (SELECT 1))",
	"(a > ANY (SELECT 1)) AND b",
	"b OR (a > ALL (SELECT 1))",
	"(a = ANY (SELECT 1)) AND b",
	"(a > ALL (SELECT 1)) IS NULL",
	"(a = ANY (SELECT 1)) IS NULL",
	"(a > ALL (SELECT 1)) IN (1, 2)",
	"(a = ANY (SELECT 1)) IN (1, 2)",
	"(a > ALL (SELECT 1)) IN (SELECT 1)",
	"(a = ANY (SELECT 1)) IN (SELECT 1)",
	"(a > ALL (SELECT 1)) BETWEEN 1 AND 2",
	"1 BETWEEN (a > ALL (SELECT 1)) AND 2",
	"1 BETWEEN 0 AND (a > ALL (SELECT 1))",
	"(a > ALL (SELECT 1)) LIKE 'x'",
	"'x' LIKE (a > ALL (SELECT 1))",
	"(a > ALL (SELECT 1)) COLLATE utf8mb4_bin",
	"(a > ALL (SELECT 1)) * (b < ANY (SELECT 1))",
	"-(a > ALL (SELECT 1))",
	"+(a > ALL (SELECT 1)) + 1",
	"(a > ALL (SELECT 1)) > ALL (SELECT 1)",
	"(a = ANY (SELECT 1)) = ANY (SELECT 1)",
	"(a <> ALL (SELECT 1)) > ALL (SELECT 1)",
	"(a IN (SELECT 1)) + 1",
	"(a NOT IN (SELECT 1)) + 1",
	"1 = (a IN (SELECT 1))",
	"(a IN (SELECT 1)) IN (SELECT 1)",
	"(a IN (SELECT 1)) BETWEEN 1 AND 2",
	"(a IN (SELECT 1)) LIKE 'x'",
	"(a IN (SELECT 1)) COLLATE utf8mb4_bin",
	"NOT a IN (SELECT 1)",
	"abs(a > ALL (SELECT 1))",
	"CASE WHEN a > ALL (SELECT 1) THEN 1 ELSE (a > ANY (SELECT 1)) + 1 END",
}

func TestQuantifiedExprDeparse(t *testing.T) {
	for _, c := range []struct{ input, want string }{
		{"a = ANY(b)", "(`a` = any (`b`))"},
		{"a = SOME(b)", "(`a` = any (`b`))"},
		{"a = ALL(b)", "(`a` = all (`b`))"},
		{"NOT (a = ANY(b))", "(not((`a` = any (`b`))))"},
	} {
		t.Run(c.input, func(t *testing.T) {
			got := Deparse(RewriteExpr(parseExpr(t, c.input)))
			if got != c.want {
				t.Errorf("Deparse(RewriteExpr(%q)) = %q, want %q", c.input, got, c.want)
			}
		})
	}
}

// TestWindowModifierAlias pins that an unaliased window call names its column
// after the modifiers the statement wrote, not the default FROM FIRST/RESPECT
// NULLS form.
func TestWindowModifierAlias(t *testing.T) {
	for _, c := range []struct{ input, want string }{
		{
			"SELECT NTH_VALUE(a, 1) FROM LAST IGNORE NULLS OVER () FROM t",
			"select nth_value(`t`.`a`,1) from last ignore nulls OVER ( )  AS `NTH_VALUE(a, 1) FROM LAST IGNORE NULLS OVER ()` from `t`",
		},
		{
			"SELECT LEAD(a IGNORE NULLS, 2) OVER () FROM t",
			"select lead(`t`.`a` ignore nulls,2) OVER ( )  AS `LEAD(a IGNORE NULLS, 2) OVER ()` from `t`",
		},
	} {
		t.Run(c.input, func(t *testing.T) {
			list, err := parser.Parse(c.input)
			if err != nil {
				t.Fatalf("Parse(%q): %v", c.input, err)
			}
			sel := list.Items[0].(*ast.SelectStmt)
			resolved := (&Resolver{Lookup: func(string) *ResolverTable {
				return &ResolverTable{Name: "t", Columns: []ResolverColumn{{Name: "a", Position: 1}}}
			}}).Resolve(sel)
			if got := DeparseSelect(resolved); got != c.want {
				t.Errorf("DeparseSelect(%q)\n got %s\nwant %s", c.input, got, c.want)
			}
		})
	}
}

// TestCastIgnoreAndGroupingSetRewrite pin that a cast's IGNORE reaches the alias
// the server would generate, and that a grouping set's items go through the
// ordinary rewrites (IF to CASE here).
func TestCastIgnoreAndGroupingSetRewrite(t *testing.T) {
	list, err := parser.Parse("SELECT CAST(a AS SIGNED IGNORE) FROM t")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	resolved := (&Resolver{Lookup: func(string) *ResolverTable {
		return &ResolverTable{Name: "t", Columns: []ResolverColumn{{Name: "a", Position: 1}}}
	}}).Resolve(list.Items[0].(*ast.SelectStmt))
	if got, want := DeparseSelect(resolved), "select cast(`t`.`a` as signed ignore) AS `CAST(a AS SIGNED IGNORE)` from `t`"; got != want {
		t.Errorf("DeparseSelect(cast)\n got %s\nwant %s", got, want)
	}

	// ROLLUP(IF(c, a, b)) rewrites the IF to a CASE inside the grouping set.
	list, err = parser.Parse("SELECT a FROM t GROUP BY ROLLUP(IF(c, a, b))")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	sel := list.Items[0].(*ast.SelectStmt)
	RewriteSelectStmt(sel)
	set, ok := sel.GroupBy[0].Expr.(*ast.GroupingSet)
	if !ok {
		t.Fatalf("GROUP BY item = %T, want *ast.GroupingSet", sel.GroupBy[0].Expr)
	}
	if _, ok := set.Items[0].(*ast.CaseExpr); !ok {
		t.Errorf("grouping set item = %T, want the IF rewritten to *ast.CaseExpr", set.Items[0])
	}
}

// TestFullJoinKeepsMergedKey pins that the deparser's own resolver merges a FULL
// join's key both in `SELECT *` and in an explicit unqualified projection, so a
// right-only row keeps its value.
func TestFullJoinKeepsMergedKey(t *testing.T) {
	lookup := func(name string) *ResolverTable {
		columns := map[string][]ResolverColumn{
			"l": {{Name: "id", Position: 1}, {Name: "lx", Position: 2}},
			"r": {{Name: "id", Position: 1}, {Name: "rx", Position: 2}},
		}[name]
		if columns == nil {
			return nil
		}
		return &ResolverTable{Name: name, Columns: columns}
	}
	for _, c := range []struct{ input, want string }{
		{
			"SELECT * FROM l NATURAL FULL JOIN r",
			"select coalesce(`l`.`id`,`r`.`id`) AS `id`,`l`.`lx` AS `lx`,`r`.`rx` AS `rx` from (`l` full join `r` on((`l`.`id` = `r`.`id`)))",
		},
		{
			"SELECT id FROM l FULL JOIN r USING(id)",
			"select coalesce(`l`.`id`,`r`.`id`) AS `id` from (`l` full join `r` on((`l`.`id` = `r`.`id`)))",
		},
		{
			// A qualified star stays side-specific.
			"SELECT l.* FROM l FULL JOIN r USING(id)",
			"select `l`.`id` AS `id`,`l`.`lx` AS `lx` from (`l` full join `r` on((`l`.`id` = `r`.`id`)))",
		},
	} {
		t.Run(c.input, func(t *testing.T) {
			list, err := parser.Parse(c.input)
			if err != nil {
				t.Fatalf("Parse(%q): %v", c.input, err)
			}
			resolved := (&Resolver{Lookup: lookup}).Resolve(list.Items[0].(*ast.SelectStmt))
			if got := DeparseSelect(resolved); got != c.want {
				t.Errorf("DeparseSelect(%q)\n got %s\nwant %s", c.input, got, c.want)
			}
		})
	}
}

// TestQuotedCollateRoundTrips pins that a string-valued collation keeps its
// quotes, so the deparsed expression parses back to the same node.
func TestQuotedCollateRoundTrips(t *testing.T) {
	const input = "SELECT a COLLATE 'x y' FROM t"
	list, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	sel := list.Items[0].(*ast.SelectStmt)
	if got, want := DeparseSelect(sel), "select (`a` collate 'x y') AS `a COLLATE 'x y'` from `t`"; got != want {
		t.Fatalf("DeparseSelect(%q) = %s, want %s", input, got, want)
	}
	if _, err := parser.Parse("SELECT a COLLATE 'x y' FROM t"); err != nil {
		t.Fatalf("re-parse: %v", err)
	}
}

// TestOrderByPositionIncludesCatalog pins that an ORDER BY naming a four-part
// column matches the target with the same catalog, not a same-named column of
// another catalog.
func TestOrderByPositionIncludesCatalog(t *testing.T) {
	list, err := parser.Parse("SELECT c1.d.t.a, c2.d.t.a FROM c1.d.t, c2.d.t ORDER BY c2.d.t.a")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	sel := list.Items[0].(*ast.SelectStmt)
	RewriteSelectStmt(sel)
	lit, ok := sel.OrderBy[0].Expr.(*ast.IntLit)
	if !ok || lit.Value != 2 {
		t.Fatalf("ORDER BY rewrote to %#v, want position 2", sel.OrderBy[0].Expr)
	}
}

// TestRemoteSequenceAndEscapeAliases pins that the rebuilt aliases of the new
// wrappers carry their source-style, unquoted text.
func TestRemoteSequenceAndEscapeAliases(t *testing.T) {
	lookup := func(string) *ResolverTable {
		return &ResolverTable{Name: "t", Columns: []ResolverColumn{{Name: "seq", Position: 1}}}
	}
	for _, c := range []struct{ input, want string }{
		{
			"SELECT t.seq@link FROM t",
			"select `t`.`seq`@link AS `t.seq@link` from `t`",
		},
		{
			"SELECT {fn CONCAT(t.seq, 'x')} FROM t",
			"select {`fn` concat(`t`.`seq`,'x')} AS `{fn CONCAT(t.seq, 'x')}` from `t`",
		},
	} {
		t.Run(c.input, func(t *testing.T) {
			list, err := parser.Parse(c.input)
			if err != nil {
				t.Fatalf("Parse(%q): %v", c.input, err)
			}
			resolved := (&Resolver{Lookup: lookup}).Resolve(list.Items[0].(*ast.SelectStmt))
			if got := DeparseSelect(resolved); got != c.want {
				t.Errorf("DeparseSelect(%q)\n got %s\nwant %s", c.input, got, c.want)
			}
		})
	}
}
