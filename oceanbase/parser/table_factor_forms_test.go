package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// OceanBase's table_factor has alternatives MySQL's lacks: a derived table with
// a flashback, RB_ITERATE(...), UNNEST(...), and a JSON_TABLE without an alias.
// This file holds each to the candidates and rejections the pinned server gave,
// recorded under table_factor, table_subquery, json_table_expr, json_table_*,
// rb_iterate_expr, unnest_expr and relation_name_list in testdata/compat.tsv.

// TestTableFactorForms pins which table factors parse.
func TestTableFactorForms(t *testing.T) {
	for _, sql := range []string{
		// select_with_parens use_flashback, plain and LATERAL.
		"SELECT a FROM(SELECT a) AS OF SNAPSHOT a",
		"SELECT a FROM(SELECT a) AS OF SNAPSHOT a a",
		"SELECT a FROM(SELECT a) AS OF SNAPSHOT a AS a",
		"SELECT a FROM LATERAL(SELECT a) AS OF SNAPSHOT a",
		// JSON_TABLE without an alias, and with each column kind.
		"SELECT a FROM JSON_TABLE(a, 's' COLUMNS(a FOR ORDINALITY))",
		"SELECT a FROM JSON_TABLE(a, 's' COLUMNS(a FOR ORDINALITY)) AS jt",
		"SELECT a FROM JSON_TABLE(a, 's' COLUMNS(a TINYINT EXISTS PATH 's'))",
		"SELECT a FROM JSON_TABLE(a, 's' COLUMNS(a TINYINT PATH 's'))",
		"SELECT a FROM JSON_TABLE(a, 's' COLUMNS(NESTED PATH 's' COLUMNS(a FOR ORDINALITY)))",
		// RB_ITERATE and UNNEST, with and without an alias and column list.
		"SELECT a FROM RB_ITERATE(a)",
		"SELECT a FROM RB_ITERATE(a) a",
		"SELECT a FROM RB_ITERATE(a) AS a",
		"SELECT a FROM RB_ITERATE(a) a(a)",
		"SELECT a FROM RB_ITERATE(a) AS a(a)",
		"SELECT a FROM UNNEST(a)",
		"SELECT a FROM UNNEST(a, a)",
		"SELECT a FROM UNNEST(a) AS a",
		"SELECT a FROM UNNEST(a) a(a)",
		"SELECT a FROM UNNEST(a) AS a(a, a)",
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v", sql, err)
		}
	}
	for _, sql := range []string{
		"SELECT a FROM(SELECT a) AS OF SNAPSHOT",
		"SELECT a FROM(SELECT a) AS OF",
		// ORDINALITY is required, and NESTED needs its PATH: the pinned server
		// answers 1064 for `a FOR` and for `NESTED 's' COLUMNS (...)`.
		"SELECT a FROM JSON_TABLE(a, 's' COLUMNS(a FOR)) a",
		"SELECT a FROM JSON_TABLE(a, 's' COLUMNS(a FOR)) AS a",
		"SELECT a FROM JSON_TABLE(a, 's' COLUMNS(NESTED 's' COLUMNS(a FOR ORDINALITY)))",
		"SELECT a FROM JSON_TABLE(a, 's' COLUMNS())",
		"SELECT a FROM RB_ITERATE()",
		"SELECT a FROM UNNEST()",
		"SELECT a FROM UNNEST(a) a()",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) was accepted", sql)
		}
	}
}

// TestTableFactorTree pins what the new table factors parse into: the derived
// table's flashback, and the two table functions' arguments and columns.
func TestTableFactorTree(t *testing.T) {
	list, err := Parse("SELECT a FROM(SELECT a) AS OF SNAPSHOT 1 t")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	sel := list.Items[0].(*nodes.SelectStmt)
	sub, ok := sel.From[0].(*nodes.SubqueryExpr)
	if !ok {
		t.Fatalf("got %T, want *nodes.SubqueryExpr", sel.From[0])
	}
	if sub.Flashback == nil || sub.Flashback.Snapshot == nil || sub.Flashback.Proctime {
		t.Errorf("flashback = %#v, want a snapshot", sub.Flashback)
	}
	if sub.Alias != "t" {
		t.Errorf("alias = %q, want t", sub.Alias)
	}

	list, err = Parse("SELECT a FROM UNNEST(a, b) AS u(c, d)")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	sel = list.Items[0].(*nodes.SelectStmt)
	unnest, ok := sel.From[0].(*nodes.UnnestExpr)
	if !ok {
		t.Fatalf("got %T, want *nodes.UnnestExpr", sel.From[0])
	}
	if len(unnest.Args) != 2 || unnest.Alias != "u" || len(unnest.Columns) != 2 {
		t.Errorf("unnest = %#v, want two arguments, alias u and two columns", unnest)
	}
}

// TestTableFactorJoinLoc pins that a join over the new table factors covers
// them: tableExprLoc reads each factor's own Loc, so the join's range does not
// fall back to NoLoc.
func TestTableFactorJoinLoc(t *testing.T) {
	const sql = "SELECT a FROM UNNEST(a) u JOIN RB_ITERATE(a) r JOIN SOURCE(TYPE = 's')"
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	sel := list.Items[0].(*nodes.SelectStmt)
	join, ok := sel.From[0].(*nodes.JoinClause)
	if !ok {
		t.Fatalf("got %T, want *nodes.JoinClause", sel.From[0])
	}
	if join.Loc.Start != len("SELECT a FROM ") || join.Loc.End != len(sql) {
		t.Errorf("join Loc = %+v, want the whole FROM list", join.Loc)
	}
}
