package catalog

import (
	"strings"
	"testing"
)

// A view body may be a TABLE or VALUES query primary (MySQL 8.0.19+), and so
// may a CTE or derived table inside one. The parser accepting it is not
// enough: the view's column list and its canonical definition are both read
// off a target list the primary does not have, so deriving nothing left the
// view with no columns and a definition of "select  from `c`" — invalid SQL.
// Container-verified on mysql:8.0:
//
//	CREATE VIEW v AS WITH c AS (TABLE t) SELECT * FROM c  → (`a`, `b`)
//	CREATE VIEW v AS SELECT * FROM (VALUES ROW(1,2)) d    → (column_0, column_1)
//	CREATE VIEW v AS VALUES ROW(1,2) UNION VALUES ROW(3,4) → (column_0, column_1)
//
// These run without a container.
func TestDefineView_QueryPrimaryBodyColumns(t *testing.T) {
	c := newCatalogWithDB(t, "mydb")
	if err := c.DefineTable(mustParseTable(t, "CREATE TABLE t (a INT, b INT)")); err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		view, ddl string
		cols      []string
	}{
		{"v1", "CREATE VIEW v1 AS TABLE t", []string{"a", "b"}},
		{"v2", "CREATE VIEW v2 AS TABLE t UNION TABLE t", []string{"a", "b"}},
		{"v3", "CREATE VIEW v3 AS WITH c AS (TABLE t) SELECT * FROM c", []string{"a", "b"}},
		{"v4", "CREATE VIEW v4 AS SELECT * FROM (TABLE t) d", []string{"a", "b"}},
		{"v5", "CREATE VIEW v5 AS VALUES ROW(1,2)", []string{"column_0", "column_1"}},
		{"v6", "CREATE VIEW v6 AS SELECT * FROM (VALUES ROW(1,2)) d", []string{"column_0", "column_1"}},
		{"v7", "CREATE VIEW v7 AS WITH c AS (VALUES ROW(1,2)) SELECT * FROM c", []string{"column_0", "column_1"}},
		{"v8", "CREATE VIEW v8 AS VALUES ROW(1,2) UNION VALUES ROW(3,4)", []string{"column_0", "column_1"}},
		{"v9", "CREATE VIEW v9 AS VALUES ROW(1,2) UNION SELECT a, b FROM t", []string{"column_0", "column_1"}},
		{"v10", "CREATE VIEW v10 AS SELECT a, b FROM t UNION VALUES ROW(1,2)", []string{"a", "b"}},
	}
	for _, tc := range cases {
		if err := c.DefineView(mustParseView(t, tc.ddl)); err != nil {
			t.Fatalf("%s: DefineView: %v", tc.view, err)
		}
	}
	db := c.GetDatabase("mydb")
	for _, tc := range cases {
		view := db.Views[tc.view]
		if view == nil {
			t.Fatalf("view %s not registered", tc.view)
		}
		if strings.Join(view.Columns, ",") != strings.Join(tc.cols, ",") {
			t.Errorf("view %s (%s) Columns = %v, want %v", tc.view, tc.ddl, view.Columns, tc.cols)
		}
		// A definition that lost the primary's projection is invalid SQL; the
		// empty select list it leaves behind is the symptom to catch.
		if strings.Contains(view.Definition, "select  ") {
			t.Errorf("view %s (%s) Definition = %q, which selects nothing", tc.view, tc.ddl, view.Definition)
		}
	}
}

// TestDefineView_PrimaryBodyNamingALocalCTE: a view body that is a TABLE primary may
// name a CTE of the view's own query, and that name may shadow a real table.
// The column list read a database-only lookup, so it answered with the stored
// table's columns -- a wrong list, not a missing one. Container-verified on
// mysql:8.0:
//
//	WITH x AS (SELECT 1 AS a) TABLE x        → (`a`)
//	WITH x AS (SELECT 1 AS a, 2 AS b) TABLE x → (`a`, `b`)
//	WITH t AS (SELECT 9 AS z) TABLE t        → (`z`), though a table t exists
//	WITH x (p) AS (SELECT 1) TABLE x         → (`p`)
//	WITH x AS (SELECT 1 AS a) TABLE x UNION TABLE x → (`a`)
//	WITH t AS (SELECT 9 AS z) TABLE mydb.t   → (`a`, `b`), the stored table:
//	                                           a qualified name is never a CTE
func TestDefineView_PrimaryBodyNamingALocalCTE(t *testing.T) {
	c := newCatalogWithDB(t, "mydb")
	if err := c.DefineTable(mustParseTable(t, "CREATE TABLE t (a INT, b INT)")); err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		view, ddl string
		cols      []string
	}{
		{"c1", "CREATE VIEW c1 AS WITH x AS (SELECT 1 AS a) TABLE x", []string{"a"}},
		{"c2", "CREATE VIEW c2 AS WITH x AS (SELECT 1 AS a, 2 AS b) TABLE x", []string{"a", "b"}},
		// The CTE shadows the stored table t, on both servers and here.
		{"c3", "CREATE VIEW c3 AS WITH t AS (SELECT 9 AS z) TABLE t", []string{"z"}},
		{"c4", "CREATE VIEW c4 AS WITH x (p) AS (SELECT 1) TABLE x", []string{"p"}},
		{"c5", "CREATE VIEW c5 AS WITH x AS (SELECT 1 AS a) TABLE x UNION TABLE x", []string{"a"}},
		// No CTE in sight: the stored table answers, as it did before.
		{"c6", "CREATE VIEW c6 AS TABLE t", []string{"a", "b"}},
		// A schema-qualified name is never a CTE, so the stored table answers
		// even though a CTE shadows the unqualified name. Consulting the CTE
		// scope here rewrote `SELECT *` over the derived table to a column the
		// table does not have.
		{"c7", "CREATE VIEW c7 AS WITH t AS (SELECT 9 AS z) TABLE mydb.t", []string{"a", "b"}},
		{"c8", "CREATE VIEW c8 AS SELECT * FROM (WITH t AS (SELECT 9 AS z) TABLE mydb.t) x", []string{"a", "b"}},
		// A CTE of the RIGHT operand of a set operation is local to it: the
		// left operand's t is the stored table on both servers. Searching the
		// sibling's WITH gave the union the right operand's columns.
		{"c9", "CREATE VIEW c9 AS SELECT * FROM (TABLE t UNION (WITH t AS (SELECT 9 AS z, 10 AS y) SELECT * FROM t)) d",
			[]string{"a", "b"}},
		// A CTE of the LEFT operand does enclose the primary, and still does.
		{"c10", "CREATE VIEW c10 AS SELECT * FROM ((WITH x AS (SELECT 9 AS z) TABLE x) UNION SELECT 1) d",
			[]string{"z"}},
		// A CTE body that is itself a primary is followed into the sibling it
		// names: both servers give this view the column a.
		{"c11", "CREATE VIEW c11 AS SELECT * FROM (WITH x AS (SELECT 1 AS a), y AS (TABLE x) TABLE y) d",
			[]string{"a"}},
		// A nested WITH shadows an outer declaration of the same name: mysql:8.0
		// answers q, the innermost. v8.5.5 answers z, the outer one, which is the
		// delta between the two trees (see findCTE).
		{"c12", "CREATE VIEW c12 AS SELECT * FROM (WITH t AS (SELECT 1 AS z) (WITH t AS (SELECT 2 AS q) TABLE t)) d",
			[]string{"q"}},
		// A body that names itself, and a pair that name each other: the search
		// stops instead of recurring. Neither query runs on either server
		// (error 3573 and 1146), so no column list is the honest answer.
		{"c13", "CREATE VIEW c13 AS SELECT * FROM (WITH RECURSIVE x AS (TABLE x) TABLE x) d", nil},
		{"c14", "CREATE VIEW c14 AS SELECT * FROM (WITH x AS (TABLE y), y AS (TABLE x) TABLE y) d", nil},
		// A body that declares a CTE of its own name is not a cycle: the inner
		// declaration is a different one, and both servers give this view a. The
		// guard keys on the declaration, not the name, which is what tells the
		// two apart (and what keeps c13 and c14 above from recurring).
		{"c15", "CREATE VIEW c15 AS SELECT * FROM (WITH x AS (WITH x AS (SELECT 1 AS a) TABLE x) TABLE x) d",
			[]string{"a"}},
		{"c16", "CREATE VIEW c16 AS SELECT * FROM (WITH x AS (WITH x AS (SELECT 1 AS a) SELECT * FROM x) TABLE x) d",
			[]string{"a"}},
	}
	for _, tc := range cases {
		if err := c.DefineView(mustParseView(t, tc.ddl)); err != nil {
			t.Fatalf("%s: DefineView: %v", tc.view, err)
		}
	}
	db := c.GetDatabase("mydb")
	for _, tc := range cases {
		view := db.Views[tc.view]
		if view == nil {
			t.Fatalf("view %s not registered", tc.view)
		}
		if strings.Join(view.Columns, ",") != strings.Join(tc.cols, ",") {
			t.Errorf("view %s (%s) Columns = %v, want %v", tc.view, tc.ddl, view.Columns, tc.cols)
		}
	}
}
