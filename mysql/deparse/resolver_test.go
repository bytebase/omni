package deparse

import (
	"strings"
	"testing"

	ast "github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/mysql/catalog"
	"github.com/bytebase/omni/mysql/parser"
)

// setupCatalog creates a catalog with a test database and tables for resolver tests.
// Schema: testdb.t(a INT, b INT, c INT), testdb.t1(a INT, b INT), testdb.t2(a INT, d INT)
func setupCatalog(t *testing.T) *catalog.Catalog {
	t.Helper()
	c := catalog.New()
	sqls := []string{
		"CREATE DATABASE testdb",
		"USE testdb",
		"CREATE TABLE t (a INT, b INT, c INT)",
		"CREATE TABLE t1 (a INT, b INT)",
		"CREATE TABLE t2 (a INT, d INT)",
	}
	for _, sql := range sqls {
		_, err := c.Exec(sql, nil)
		if err != nil {
			t.Fatalf("catalog setup failed on %q: %v", sql, err)
		}
	}
	return c
}

// catalogLookup creates a TableLookup function from a catalog.Catalog.
func catalogLookup(c *catalog.Catalog) TableLookup {
	return func(tableName string) *ResolverTable {
		db := c.GetDatabase(c.CurrentDatabase())
		if db == nil {
			return nil
		}
		tbl := db.GetTable(tableName)
		if tbl == nil {
			return nil
		}
		cols := make([]ResolverColumn, len(tbl.Columns))
		for i, col := range tbl.Columns {
			cols[i] = ResolverColumn{Name: col.Name, Position: col.Position}
		}
		return &ResolverTable{Name: tbl.Name, Columns: cols}
	}
}

// resolveAndDeparse parses SQL, resolves column refs, rewrites, and deparses.
func resolveAndDeparse(t *testing.T, cat *catalog.Catalog, sql string) string {
	t.Helper()
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("failed to parse %q: %v", sql, err)
	}
	if list.Len() == 0 {
		t.Fatalf("no statements parsed from %q", sql)
	}
	sel, ok := list.Items[0].(*ast.SelectStmt)
	if !ok {
		t.Fatalf("expected SelectStmt, got %T", list.Items[0])
	}

	// Apply rewrites first (NOT folding, boolean context), then resolve
	rewriteSelectStmt(sel)

	resolver := &Resolver{Lookup: catalogLookup(cat)}
	resolved := resolver.Resolve(sel)
	return DeparseSelect(resolved)
}

// rewriteSelectStmt applies RewriteExpr to all expression positions in a SelectStmt.
func rewriteSelectStmt(stmt *ast.SelectStmt) {
	if stmt == nil {
		return
	}
	if stmt.SetOp != ast.SetOpNone {
		rewriteSelectStmt(stmt.Left)
		rewriteSelectStmt(stmt.Right)
		return
	}
	for i, target := range stmt.TargetList {
		if rt, ok := target.(*ast.ResTarget); ok {
			rt.Val = RewriteExpr(rt.Val)
		} else {
			stmt.TargetList[i] = RewriteExpr(target)
		}
	}
	if stmt.Where != nil {
		stmt.Where = RewriteExpr(stmt.Where)
	}
	for i, expr := range stmt.GroupBy {
		stmt.GroupBy[i] = RewriteExpr(expr)
	}
	if stmt.Having != nil {
		stmt.Having = RewriteExpr(stmt.Having)
	}
	for _, item := range stmt.OrderBy {
		item.Expr = RewriteExpr(item.Expr)
	}
}

func TestResolver_Section_6_1_ColumnQualification(t *testing.T) {
	cat := setupCatalog(t)

	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"single_table_unqualified",
			"SELECT a FROM t",
			"select `t`.`a` AS `a` from `t`",
		},
		{
			"multiple_columns",
			"SELECT a, b, c FROM t",
			"select `t`.`a` AS `a`,`t`.`b` AS `b`,`t`.`c` AS `c` from `t`",
		},
		{
			"qualified_column_preserved",
			"SELECT t.a FROM t",
			"select `t`.`a` AS `a` from `t`",
		},
		{
			"table_alias",
			"SELECT a FROM t AS x",
			"select `x`.`a` AS `a` from `t` `x`",
		},
		{
			"table_alias_no_as",
			"SELECT a FROM t x",
			"select `x`.`a` AS `a` from `t` `x`",
		},
		{
			"column_in_where",
			"SELECT a FROM t WHERE a > 0",
			"select `t`.`a` AS `a` from `t` where (`t`.`a` > 0)",
		},
		{
			"column_in_order_by",
			"SELECT a FROM t ORDER BY a",
			"select `t`.`a` AS `a` from `t` order by `t`.`a`",
		},
		{
			"column_in_group_by",
			"SELECT a FROM t GROUP BY a",
			"select `t`.`a` AS `a` from `t` group by `t`.`a`",
		},
		{
			"column_in_having",
			"SELECT a, COUNT(*) FROM t GROUP BY a HAVING a > 0",
			"select `t`.`a` AS `a`,count(0) AS `count(0)` from `t` group by `t`.`a` having (`t`.`a` > 0)",
		},
		{
			"column_in_on_condition",
			"SELECT t1.a, t2.d FROM t1 JOIN t2 ON t1.a = t2.a",
			"select `t1`.`a` AS `a`,`t2`.`d` AS `d` from (`t1` join `t2` on((`t1`.`a` = `t2`.`a`)))",
		},
		{
			"qualified_star",
			"SELECT t1.*, t2.d FROM t1 JOIN t2 ON t1.a = t2.a",
			"select `t1`.`a` AS `a`,`t1`.`b` AS `b`,`t2`.`d` AS `d` from (`t1` join `t2` on((`t1`.`a` = `t2`.`a`)))",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := resolveAndDeparse(t, cat, tc.input)
			if got != tc.expected {
				t.Errorf("resolveAndDeparse(%q) =\n  %q\nwant:\n  %q", tc.input, got, tc.expected)
			}
		})
	}
}

func TestResolver_Section_6_1_AmbiguousColumn(t *testing.T) {
	cat := setupCatalog(t)

	// Column 'a' exists in both t1 and t2 — should still work (resolve to first match)
	// For MySQL, an ambiguous unqualified column in a multi-table query is an error.
	// Our resolver currently resolves to the first match for simplicity.
	t.Run("ambiguous_column_two_tables", func(t *testing.T) {
		got := resolveAndDeparse(t, cat, "SELECT a FROM t1 JOIN t2 ON t1.a = t2.a")
		// The unqualified 'a' in SELECT will match t1 first (insertion order)
		expected := "select `t1`.`a` AS `a` from (`t1` join `t2` on((`t1`.`a` = `t2`.`a`)))"
		if got != expected {
			t.Errorf("resolveAndDeparse(ambiguous) =\n  %q\nwant:\n  %q", got, expected)
		}
	})
}

func TestResolverTable_GetColumn(t *testing.T) {
	rt := &ResolverTable{
		Name: "t",
		Columns: []ResolverColumn{
			{Name: "a", Position: 1},
			{Name: "b", Position: 2},
		},
	}

	if rt.GetColumn("a") == nil {
		t.Error("expected to find column 'a'")
	}
	if rt.GetColumn("A") == nil {
		t.Error("expected case-insensitive match for 'A'")
	}
	if rt.GetColumn("c") != nil {
		t.Error("expected nil for nonexistent column 'c'")
	}
}

// Ensure the Resolver handles an unused import gracefully.
var _ = strings.ToLower
