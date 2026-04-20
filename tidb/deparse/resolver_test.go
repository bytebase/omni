package deparse_test

import (
	"strings"
	"testing"

	ast "github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/tidb/catalog"
	"github.com/bytebase/omni/tidb/deparse"
	"github.com/bytebase/omni/tidb/parser"
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
func catalogLookup(c *catalog.Catalog) deparse.TableLookup {
	return func(tableName string) *deparse.ResolverTable {
		db := c.GetDatabase(c.CurrentDatabase())
		if db == nil {
			return nil
		}
		tbl := db.GetTable(tableName)
		if tbl == nil {
			return nil
		}
		cols := make([]deparse.ResolverColumn, len(tbl.Columns))
		for i, col := range tbl.Columns {
			cols[i] = deparse.ResolverColumn{Name: col.Name, Position: col.Position}
		}
		return &deparse.ResolverTable{Name: tbl.Name, Columns: cols}
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
	deparse.RewriteSelectStmt(sel)

	// Get the database default charset for the resolver
	defaultCharset := ""
	db := cat.GetDatabase(cat.CurrentDatabase())
	if db != nil {
		defaultCharset = db.Charset
	}

	resolver := &deparse.Resolver{Lookup: catalogLookup(cat), DefaultCharset: defaultCharset}
	resolved := resolver.Resolve(sel)
	return deparse.DeparseSelect(resolved)
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
			"select `t`.`a` AS `a`,count(0) AS `COUNT(*)` from `t` group by `t`.`a` having (`t`.`a` > 0)",
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
	rt := &deparse.ResolverTable{
		Name: "t",
		Columns: []deparse.ResolverColumn{
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

func TestResolver_Section_6_2_SelectStarExpansion(t *testing.T) {
	cat := setupCatalog(t)

	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"simple_star",
			"SELECT * FROM t",
			"select `t`.`a` AS `a`,`t`.`b` AS `b`,`t`.`c` AS `c` from `t`",
		},
		{
			"star_alias_per_column",
			// Each expanded column gets `table`.`col` AS `col` format
			"SELECT * FROM t1",
			"select `t1`.`a` AS `a`,`t1`.`b` AS `b` from `t1`",
		},
		{
			"star_ordering",
			// Columns in table definition order (Column.Position)
			"SELECT * FROM t",
			"select `t`.`a` AS `a`,`t`.`b` AS `b`,`t`.`c` AS `c` from `t`",
		},
		{
			"star_with_where",
			"SELECT * FROM t WHERE a > 0",
			"select `t`.`a` AS `a`,`t`.`b` AS `b`,`t`.`c` AS `c` from `t` where (`t`.`a` > 0)",
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

func TestResolver_Section_6_3_AutoAliasGeneration(t *testing.T) {
	cat := setupCatalog(t)

	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"column_ref",
			"SELECT a FROM t",
			"select `t`.`a` AS `a` from `t`",
		},
		{
			"qualified_column",
			"SELECT t.a FROM t",
			"select `t`.`a` AS `a` from `t`",
		},
		{
			"expression_auto_alias",
			"SELECT a + b FROM t",
			"select (`t`.`a` + `t`.`b`) AS `a + b` from `t`",
		},
		{
			"literal_int",
			"SELECT 1",
			"select 1 AS `1`",
		},
		{
			"literal_string",
			"SELECT 'hello'",
			"select 'hello' AS `hello`",
		},
		{
			"literal_null",
			"SELECT NULL",
			"select NULL AS `NULL`",
		},
		{
			"complex_expression_name_exp",
			// Expression alias > 64 chars triggers Name_exp_N pattern
			"SELECT CONCAT(CONCAT(CONCAT(CONCAT(CONCAT(CONCAT(CONCAT(a, b), c), a), b), c), a), b) FROM t",
			"select concat(concat(concat(concat(concat(concat(concat(`t`.`a`,`t`.`b`),`t`.`c`),`t`.`a`),`t`.`b`),`t`.`c`),`t`.`a`),`t`.`b`) AS `Name_exp_1` from `t`",
		},
		{
			"explicit_alias_preserved",
			"SELECT a AS x FROM t",
			"select `t`.`a` AS `x` from `t`",
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

// setupCatalogForJoins creates a catalog with tables that have overlapping columns
// for testing NATURAL JOIN expansion.
// Schema: testdb.t1(a INT, b INT), testdb.t2(a INT, d INT), testdb.t3(a INT, b INT)
func setupCatalogForJoins(t *testing.T) *catalog.Catalog {
	t.Helper()
	c := catalog.New()
	sqls := []string{
		"CREATE DATABASE testdb",
		"USE testdb",
		"CREATE TABLE t1 (a INT, b INT)",
		"CREATE TABLE t2 (a INT, d INT)",
		"CREATE TABLE t3 (a INT, b INT)",
	}
	for _, sql := range sqls {
		_, err := c.Exec(sql, nil)
		if err != nil {
			t.Fatalf("catalog setup failed on %q: %v", sql, err)
		}
	}
	return c
}

func TestResolver_Section_6_4_JoinNormalization(t *testing.T) {
	t.Run("natural_join", func(t *testing.T) {
		// t1(a, b) NATURAL JOIN t2(a, d) → common column: a
		// Test with explicit column selection to verify ON expansion
		cat := setupCatalogForJoins(t)
		got := resolveAndDeparse(t, cat, "SELECT t1.a, t1.b, t2.d FROM t1 NATURAL JOIN t2")
		expected := "select `t1`.`a` AS `a`,`t1`.`b` AS `b`,`t2`.`d` AS `d` from (`t1` join `t2` on((`t1`.`a` = `t2`.`a`)))"
		if got != expected {
			t.Errorf("NATURAL JOIN:\n  got:  %q\n  want: %q", got, expected)
		}
	})

	t.Run("natural_join_multi_common", func(t *testing.T) {
		// t1(a, b) NATURAL JOIN t3(a, b) → common columns: a, b
		cat := setupCatalogForJoins(t)
		got := resolveAndDeparse(t, cat, "SELECT t1.a, t1.b FROM t1 NATURAL JOIN t3")
		expected := "select `t1`.`a` AS `a`,`t1`.`b` AS `b` from (`t1` join `t3` on(((`t1`.`a` = `t3`.`a`) and (`t1`.`b` = `t3`.`b`))))"
		if got != expected {
			t.Errorf("NATURAL JOIN multi common:\n  got:  %q\n  want: %q", got, expected)
		}
	})

	t.Run("using_single_column", func(t *testing.T) {
		// USING (a) → on((`t1`.`a` = `t2`.`a`))
		cat := setupCatalogForJoins(t)
		got := resolveAndDeparse(t, cat, "SELECT t1.a, t1.b, t2.d FROM t1 JOIN t2 USING (a)")
		expected := "select `t1`.`a` AS `a`,`t1`.`b` AS `b`,`t2`.`d` AS `d` from (`t1` join `t2` on((`t1`.`a` = `t2`.`a`)))"
		if got != expected {
			t.Errorf("USING single:\n  got:  %q\n  want: %q", got, expected)
		}
	})

	t.Run("using_multiple_columns", func(t *testing.T) {
		// USING (a, b) via AST since parser may not support multi-column USING
		cat := setupCatalogForJoins(t)
		// Build manually: t1 JOIN t3 USING (a, b)
		join := &ast.JoinClause{
			Type:  ast.JoinInner,
			Left:  &ast.TableRef{Name: "t1"},
			Right: &ast.TableRef{Name: "t3"},
			Condition: &ast.UsingCondition{
				Columns: []string{"a", "b"},
			},
		}
		sel := &ast.SelectStmt{
			TargetList: []ast.ExprNode{&ast.StarExpr{}},
			From:       []ast.TableExpr{join},
		}
		deparse.RewriteSelectStmt(sel)
		resolver := &deparse.Resolver{Lookup: catalogLookup(cat)}
		resolved := resolver.Resolve(sel)
		got := deparse.DeparseSelect(resolved)
		expected := "select `t1`.`a` AS `a`,`t1`.`b` AS `b`,`t3`.`a` AS `a`,`t3`.`b` AS `b` from (`t1` join `t3` on(((`t1`.`a` = `t3`.`a`) and (`t1`.`b` = `t3`.`b`))))"
		if got != expected {
			t.Errorf("USING multi:\n  got:  %q\n  want: %q", got, expected)
		}
	})

	t.Run("right_join_to_left_join", func(t *testing.T) {
		// RIGHT JOIN → LEFT JOIN with table swap
		cat := setupCatalogForJoins(t)
		got := resolveAndDeparse(t, cat, "SELECT t1.a, t2.d FROM t1 RIGHT JOIN t2 ON t1.a = t2.a")
		expected := "select `t1`.`a` AS `a`,`t2`.`d` AS `d` from (`t2` left join `t1` on((`t1`.`a` = `t2`.`a`)))"
		if got != expected {
			t.Errorf("RIGHT JOIN:\n  got:  %q\n  want: %q", got, expected)
		}
	})

	t.Run("cross_join_to_join", func(t *testing.T) {
		// CROSS JOIN → plain join (no ON)
		cat := setupCatalogForJoins(t)
		got := resolveAndDeparse(t, cat, "SELECT t1.a, t2.d FROM t1 CROSS JOIN t2")
		expected := "select `t1`.`a` AS `a`,`t2`.`d` AS `d` from (`t1` join `t2`)"
		if got != expected {
			t.Errorf("CROSS JOIN:\n  got:  %q\n  want: %q", got, expected)
		}
	})

	t.Run("implicit_cross_join", func(t *testing.T) {
		// FROM t1, t2 → FROM (`t1` join `t2`)
		cat := setupCatalogForJoins(t)
		got := resolveAndDeparse(t, cat, "SELECT t1.a, t2.d FROM t1, t2")
		expected := "select `t1`.`a` AS `a`,`t2`.`d` AS `d` from (`t1` join `t2`)"
		if got != expected {
			t.Errorf("Implicit cross join:\n  got:  %q\n  want: %q", got, expected)
		}
	})
}

// setupCatalogWithCharset creates a catalog with a database using the given charset.
func setupCatalogWithCharset(t *testing.T, charset string) *catalog.Catalog {
	t.Helper()
	c := catalog.New()
	createDB := "CREATE DATABASE testdb"
	if charset != "" {
		createDB += " CHARACTER SET " + charset
	}
	sqls := []string{
		createDB,
		"USE testdb",
		"CREATE TABLE t (a INT, b VARCHAR(100))",
	}
	for _, sql := range sqls {
		_, err := c.Exec(sql, nil)
		if err != nil {
			t.Fatalf("catalog setup failed on %q: %v", sql, err)
		}
	}
	return c
}

func TestResolver_Section_6_5_CastCharsetFromCatalog(t *testing.T) {
	t.Run("cast_char_uses_database_default_charset_utf8mb4", func(t *testing.T) {
		// Default database charset is utf8mb4 — CAST to CHAR should use charset utf8mb4
		cat := setupCatalog(t) // uses default charset (utf8mb4)
		got := resolveAndDeparse(t, cat, "SELECT CAST(a AS CHAR) FROM t")
		expected := "select cast(`t`.`a` as char charset utf8mb4) AS `CAST(a AS CHAR)` from `t`"
		if got != expected {
			t.Errorf("CAST CHAR utf8mb4:\n  got:  %q\n  want: %q", got, expected)
		}
	})

	t.Run("cast_char_latin1_database", func(t *testing.T) {
		// Database with latin1 charset — CAST to CHAR should use charset latin1
		cat := setupCatalogWithCharset(t, "latin1")
		got := resolveAndDeparse(t, cat, "SELECT CAST(a AS CHAR) FROM t")
		expected := "select cast(`t`.`a` as char charset latin1) AS `CAST(a AS CHAR)` from `t`"
		if got != expected {
			t.Errorf("CAST CHAR latin1:\n  got:  %q\n  want: %q", got, expected)
		}
	})

	t.Run("cast_char_n_latin1_database", func(t *testing.T) {
		// Database with latin1 charset — CAST to CHAR(10) should use charset latin1
		cat := setupCatalogWithCharset(t, "latin1")
		got := resolveAndDeparse(t, cat, "SELECT CAST(a AS CHAR(10)) FROM t")
		expected := "select cast(`t`.`a` as char(10) charset latin1) AS `CAST(a AS CHAR(10))` from `t`"
		if got != expected {
			t.Errorf("CAST CHAR(10) latin1:\n  got:  %q\n  want: %q", got, expected)
		}
	})

	t.Run("cast_binary_unaffected", func(t *testing.T) {
		// CAST to BINARY should always use charset binary, regardless of database charset
		cat := setupCatalogWithCharset(t, "latin1")
		got := resolveAndDeparse(t, cat, "SELECT CAST(a AS BINARY) FROM t")
		expected := "select cast(`t`.`a` as char charset binary) AS `CAST(a AS BINARY)` from `t`"
		if got != expected {
			t.Errorf("CAST BINARY (latin1 db):\n  got:  %q\n  want: %q", got, expected)
		}
	})

	t.Run("convert_type_latin1_database", func(t *testing.T) {
		// CONVERT(expr, CHAR) — rewritten to CAST, should use database charset
		// MySQL 8.0 auto-alias preserves the original CONVERT(a, CHAR) form.
		cat := setupCatalogWithCharset(t, "latin1")
		got := resolveAndDeparse(t, cat, "SELECT CONVERT(a, CHAR) FROM t")
		expected := "select cast(`t`.`a` as char charset latin1) AS `CONVERT(a, CHAR)` from `t`"
		if got != expected {
			t.Errorf("CONVERT CHAR latin1:\n  got:  %q\n  want: %q", got, expected)
		}
	})
}

// Ensure the Resolver handles an unused import gracefully.
var _ = strings.ToLower
