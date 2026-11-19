package parser

import (
	"testing"

	"github.com/bytebase/omni/mssql/ast"
)

// TestTableSourceShapes pins the FROM-clause shapes the pinned SQL Server
// answered under SET PARSEONLY ON: testdata/compat.tsv's
// selectTableReferenceElementWithoutJoinParenthesis, fromClause,
// simpleTableHints, fulltextTableOptions, forPathClause and derivedTable rows.
func TestTableSourceShapes(t *testing.T) {
	for _, sql := range []string{
		// A table variable's method call takes an alias and its column list.
		"SELECT * FROM @a.a() a(a)",
		// A derived table takes an alias.
		"SELECT * FROM(SELECT 1 AS a) d",
		"SELECT * FROM(SELECT 1 AS a) FOR PATH d",
		// FOR PATH follows a table name before its alias.
		"SELECT * FROM a FOR PATH",
		"SELECT * FROM a FOR PATH x",
		// WITH (n) is the index hint that holds one index.
		"SELECT * FROM a WITH(1)",
		"SELECT * FROM a WITH(INDEX(1))",
		// An ODBC rowset function's arguments.
		"SELECT * FROM OPENXML(@a, 'a')",
		"SELECT * FROM OPENXML(@a, 'a', 1)",
		"SELECT * FROM OPENXML(@a, @b, @c)",
		"SELECT * FROM OPENQUERY(srv, 'SELECT 1')",
		"SELECT * FROM OPENROWSET('SQLNCLI', 'srv', 'SELECT 1')",
		"SELECT * FROM CONTAINSTABLE(a, a, 'a', 1)",
	} {
		ParseAndCheck(t, sql)
	}

	for _, sql := range []string{
		// The alias and the column list are both required, and the list is
		// not empty: FROM @a.a()(a), FROM @a.a() a, FROM @a.a() a().
		"SELECT * FROM @a.a()(a)",
		"SELECT * FROM @a.a() a",
		"SELECT * FROM @a.a() a()",
		// A comma in the FROM list promises another table source.
		"SELECT * FROM @a,",
		// A derived table without an alias.
		"SELECT * FROM(SELECT 1 AS a)",
		// tableHintsBody holds a hint: WITH () is empty.
		"SELECT * FROM a WITH()",
		// OPENXML takes a variable handle, and OPENQUERY an identifier and a
		// string literal.
		"SELECT * FROM OPENXML(a)",
		"SELECT * FROM OPENXML(a, 'a')",
		"SELECT * FROM OPENXML(@a, 1)",
		"SELECT * FROM OPENXML(@a, 'a', 0x01)",
		"SELECT * FROM OPENQUERY(1)",
		"SELECT * FROM OPENQUERY('srv', 'SELECT 1')",
		"SELECT * FROM OPENQUERY(srv)",
		// OPENROWSET's first argument is the provider, a string literal:
		// an identifier there is Msg 156, near OPENROWSET.
		"SELECT * FROM OPENROWSET(a)",
		"SELECT * FROM OPENROWSET(a) zzz",
		// fulltextTableOptions' top-N argument is an integer or a variable,
		// and it follows LANGUAGE, which it may not precede.
		"SELECT * FROM CONTAINSTABLE(a, a, 'a', 1, LANGUAGE 0x00)",
		"SELECT * FROM CONTAINSTABLE(a, a, 'a', 0x00)",
		// SPATIAL_WINDOW_MAX_CELLS takes a value.
		"DELETE a WITH(SPATIAL_WINDOW_MAX_CELLS =)",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) = nil error, want SQL Server's rejection", sql)
		}
	}
}

// TestTableSourceAST pins the fields these shapes set.
func TestTableSourceAST(t *testing.T) {
	from := func(t *testing.T, sql string) ast.Node {
		t.Helper()
		list, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", sql, err)
		}
		return list.Items[0].(*ast.SelectStmt).FromClause.Items[0]
	}

	mc, ok := from(t, "SELECT * FROM @a.a() a(c)").(*ast.TableVarMethodCallRef)
	if !ok {
		t.Fatalf("FROM = %T, want *TableVarMethodCallRef", from(t, "SELECT * FROM @a.a() a(c)"))
	}
	if mc.Alias != "a" || len(mc.Columns) != 1 || mc.Columns[0] != "c" {
		t.Errorf("table variable method call = alias %q columns %v, want a and [c]", mc.Alias, mc.Columns)
	}

	ref, ok := from(t, "SELECT * FROM a FOR PATH").(*ast.TableRef)
	if !ok {
		t.Fatalf("FROM = %T, want *TableRef", from(t, "SELECT * FROM a FOR PATH"))
	}
	if !ref.ForPath {
		t.Error("ForPath = false, want true")
	}

	derived, ok := from(t, "SELECT * FROM(SELECT 1 AS a) FOR PATH d").(*ast.AliasedTableRef)
	if !ok {
		t.Fatalf("FROM = %T, want *AliasedTableRef", from(t, "SELECT * FROM(SELECT 1 AS a) FOR PATH d"))
	}
	if !derived.ForPath || derived.Alias != "d" {
		t.Errorf("derived table = forPath %t alias %q, want true and d", derived.ForPath, derived.Alias)
	}

	hinted, ok := from(t, "SELECT * FROM a WITH(1)").(*ast.TableRef)
	if !ok {
		t.Fatalf("FROM = %T, want *TableRef", from(t, "SELECT * FROM a WITH(1)"))
	}
	hint := hinted.Hints.Items[0].(*ast.TableHint)
	if !hint.OldForceIndex || hint.IndexValues == nil || hint.IndexValues.Len() != 1 {
		t.Errorf("hint = %+v, want WITH(1) as the old force index", hint)
	}
}
