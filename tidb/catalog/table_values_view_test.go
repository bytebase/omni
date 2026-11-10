package catalog

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// mustParseSelect parses a query expression and returns its SelectStmt.
func mustParseSelect(t *testing.T, sql string) *nodes.SelectStmt {
	t.Helper()
	stmt, ok := parseFirst(t, sql).(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("not a SelectStmt: %T", parseFirst(t, sql))
	}
	return stmt
}

// A view body may be a TABLE or VALUES query primary (MySQL 8.0.19+), and so
// may its set-operation operands. The parser accepting one is not enough: the
// catalog has to give the view the columns TiDB gives it and keep the body
// text, and both go through code paths that read a target list the primary
// does not have. Container-verified on pingcap/tidb:v8.5.5:
//
//	CREATE VIEW vt AS TABLE t                  → (`a`, `b`), body "TABLE `test`.`t`"
//	CREATE VIEW vu AS TABLE t UNION TABLE t    → (`a`, `b`)
//	CREATE VIEW vv AS VALUES ROW(1,2)          → error 1051, Unknown table ''
//
// The catalog stores a VALUES body in TiDB's restore format although TiDB
// refuses one. These run without a container.

// TestDefineView_TablePrimaryBodyColumns: a TABLE primary is SELECT * FROM t,
// so the view's columns are the table's. Reading the primary's (absent) target
// list derived none, which left the view with an empty column list that no
// consumer could expand.
func TestDefineView_TablePrimaryBodyColumns(t *testing.T) {
	c := newCatalogWithDB(t, "mydb")
	if err := c.DefineTable(mustParseTable(t, "CREATE TABLE t (a INT, b INT)")); err != nil {
		t.Fatal(err)
	}
	cases := []struct{ view, ddl string }{
		{"v1", "CREATE VIEW v1 AS TABLE t"},
		{"v2", "CREATE VIEW v2 AS (TABLE t)"},
		{"v3", "CREATE VIEW v3 AS TABLE t ORDER BY a LIMIT 1"},
		{"v4", "CREATE VIEW v4 AS TABLE t UNION TABLE t"},
		{"v5", "CREATE VIEW v5 AS SELECT a, b FROM t UNION TABLE t"},
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
		if len(view.Columns) != 2 || view.Columns[0] != "a" || view.Columns[1] != "b" {
			t.Errorf("view %s (%s) Columns = %v, want [a b]", tc.view, tc.ddl, view.Columns)
		}
	}
}

// TestDefineView_TablePrimaryBodyDefinition: the canonical definition keeps the
// TABLE primary, as TiDB's SHOW CREATE VIEW does. A deparse that did not know
// the primary wrote "select " and dropped the table, so the stored definition
// named no relation at all.
func TestDefineView_TablePrimaryBodyDefinition(t *testing.T) {
	c := newCatalogWithDB(t, "mydb")
	if err := c.DefineTable(mustParseTable(t, "CREATE TABLE t (a INT, b INT)")); err != nil {
		t.Fatal(err)
	}
	cases := []struct{ view, ddl, want string }{
		{"v1", "CREATE VIEW v1 AS TABLE t", "TABLE `mydb`.`t`"},
		{"v2", "CREATE VIEW v2 AS TABLE t ORDER BY a LIMIT 1", "TABLE `mydb`.`t` ORDER BY `a` LIMIT 1"},
		{"v3", "CREATE VIEW v3 AS TABLE t UNION TABLE t", "TABLE `mydb`.`t` UNION TABLE `mydb`.`t`"},
		{"v4", "CREATE VIEW v4 AS VALUES ROW(1,2)", "VALUES ROW(1,2)"},
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
		if view.Definition != tc.want {
			t.Errorf("view %s (%s) Definition = %q, want %q", tc.view, tc.ddl, view.Definition, tc.want)
		}
		if strings.Contains(view.Definition, "unsupported") {
			t.Errorf("view %s Definition = %q, which names a node deparse does not write", tc.view, view.Definition)
		}
	}
}

// TestDefineView_PrimaryThroughCTEOrDerivedTable: a CTE or derived table over
// a query primary must expose the primary's columns, or SELECT * over it
// expands to nothing and the stored definition is the invalid
// "select  from `c`". Container-verified on pingcap/tidb:v8.5.5:
//
//	CREATE VIEW v AS WITH c AS (TABLE t) SELECT * FROM c  → (`a`, `b`)
//	CREATE VIEW v AS SELECT * FROM (TABLE t) d            → (`a`, `b`)
//
// A VALUES primary gets ValuesColumnName's columns. TiDB cannot run a VALUES
// query at all (error 1051, Unknown table ”), so the names are MySQL's, whose
// VALUES statement TiDB's grammar implements and which answers column_0 and
// column_1 for each of these.
func TestDefineView_PrimaryThroughCTEOrDerivedTable(t *testing.T) {
	c := newCatalogWithDB(t, "mydb")
	if err := c.DefineTable(mustParseTable(t, "CREATE TABLE t (a INT, b INT)")); err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		view, ddl string
		cols      []string
		def       string
	}{
		{"v1", "CREATE VIEW v1 AS WITH c AS (TABLE t) SELECT * FROM c", []string{"a", "b"},
			"WITH `c` AS (TABLE `mydb`.`t`) SELECT `c`.`a` AS `a`,`c`.`b` AS `b` FROM `c`"},
		{"v2", "CREATE VIEW v2 AS SELECT * FROM (TABLE t) d", []string{"a", "b"},
			"SELECT `mydb`.`d`.`a` AS `a`,`mydb`.`d`.`b` AS `b` FROM (TABLE `mydb`.`t`) AS `d`"},
		{"v3", "CREATE VIEW v3 AS SELECT * FROM (VALUES ROW(1,2)) d", []string{"column_0", "column_1"},
			"SELECT `d`.`column_0` AS `column_0`,`d`.`column_1` AS `column_1` FROM (VALUES ROW(1,2)) AS `d`"},
		{"v4", "CREATE VIEW v4 AS WITH c AS (VALUES ROW(1,2)) SELECT * FROM c", []string{"column_0", "column_1"},
			"WITH `c` AS (VALUES ROW(1,2)) SELECT `c`.`column_0` AS `column_0`,`c`.`column_1` AS `column_1` FROM `c`"},
		{"v5", "CREATE VIEW v5 AS VALUES ROW(1,2) UNION VALUES ROW(3,4)", []string{"column_0", "column_1"},
			"VALUES ROW(1,2) UNION VALUES ROW(3,4)"},
		{"v6", "CREATE VIEW v6 AS VALUES ROW(1,2) UNION SELECT a, b FROM t", []string{"column_0", "column_1"},
			"VALUES ROW(1,2) UNION SELECT `a` AS `a`,`b` AS `b` FROM `mydb`.`t`"},
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
		if view.Definition != tc.def {
			t.Errorf("view %s (%s) Definition =\n  %q\nwant:\n  %q", tc.view, tc.ddl, view.Definition, tc.def)
		}
	}
}

// TestDefineView_PrimaryBodyNamingALocalCTE: a view body that is a TABLE primary may
// name a CTE of the view's own query, and that name may shadow a real table.
// The column list read a database-only lookup, so it answered with the stored
// table's columns -- a wrong list, not a missing one. Container-verified on
// pingcap/tidb:v8.5.5:
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
		// A nested WITH shadowing an outer name: v8.5.5 answers with the OUTER
		// declaration, z, and drops the inner WITH from its own SHOW CREATE VIEW
		// body. mysql:8.0 answers q, the innermost; that is the delta between the
		// two trees (see findCTE).
		{"c12", "CREATE VIEW c12 AS SELECT * FROM (WITH t AS (SELECT 1 AS z) (WITH t AS (SELECT 2 AS q) TABLE t)) d",
			[]string{"z"}},
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

// TestDefineView_PrimaryNamingALocalCTE: a TABLE primary may name a CTE of its
// own query expression, which the enclosing lookup does not hold. Resolving
// the primary against the caller's tables alone left the derived table with no
// columns, so the stored definition was again the invalid "select  from (...)".
func TestDefineView_PrimaryNamingALocalCTE(t *testing.T) {
	c := newCatalogWithDB(t, "mydb")
	cases := []struct{ view, ddl, def string }{
		{"v1", "CREATE VIEW v1 AS SELECT * FROM (WITH x AS (SELECT 1 AS a) TABLE x) d",
			"SELECT `d`.`a` AS `a` FROM (WITH `x` AS (SELECT 1 AS `a`) TABLE `x`) AS `d`"},
		{"v2", "CREATE VIEW v2 AS SELECT * FROM (WITH x AS (SELECT 1 AS a, 2 AS b) TABLE x) d",
			"SELECT `d`.`a` AS `a`,`d`.`b` AS `b` FROM (WITH `x` AS (SELECT 1 AS `a`,2 AS `b`) TABLE `x`) AS `d`"},
		{"v3", "CREATE VIEW v3 AS SELECT * FROM (WITH x (p, q) AS (SELECT 1, 2) TABLE x) d",
			"SELECT `d`.`p` AS `p`,`d`.`q` AS `q` FROM (WITH `x` (`p`, `q`) AS (SELECT 1 AS `1`,2 AS `2`) TABLE `x`) AS `d`"},
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
		if view.Definition != tc.def {
			t.Errorf("view %s (%s) Definition =\n  %q\nwant:\n  %q", tc.view, tc.ddl, view.Definition, tc.def)
		}
	}

	// A recursive CTE whose body names itself must not send the search round
	// in a circle. It derives no columns, as before; what matters is that it
	// returns.
	if err := c.DefineView(mustParseView(t,
		"CREATE VIEW vr AS SELECT * FROM (WITH RECURSIVE x AS (TABLE x) TABLE x) d")); err != nil {
		t.Fatalf("recursive CTE body: DefineView: %v", err)
	}
}

// TestDefineView_SubqueryInsideAValuesRow: a VALUES row may hold a subquery,
// and the row is the only place its SELECT can be reached from — a VALUES
// primary has no target list and no FROM, so nothing else in the resolver
// walks it. An unresolved row left a star unexpanded, and the definition
// stored for the view was not SQL.
func TestDefineView_SubqueryInsideAValuesRow(t *testing.T) {
	c := newCatalogWithDB(t, "mydb")
	if err := c.DefineTable(mustParseTable(t, "CREATE TABLE t (a INT)")); err != nil {
		t.Fatal(err)
	}
	cases := []struct{ view, ddl, def string }{
		{"v1", "CREATE VIEW v1 AS VALUES ROW((SELECT * FROM t))",
			"VALUES ROW((SELECT `mydb`.`t`.`a` AS `a` FROM `mydb`.`t`))"},
		{"v2", "CREATE VIEW v2 AS VALUES ROW(1), ROW((SELECT * FROM t))",
			"VALUES ROW(1), ROW((SELECT `mydb`.`t`.`a` AS `a` FROM `mydb`.`t`))"},
		{"v3", "CREATE VIEW v3 AS SELECT * FROM (VALUES ROW((SELECT * FROM t))) d",
			"SELECT `d`.`column_0` AS `column_0` FROM (VALUES ROW((SELECT `mydb`.`t`.`a` AS `a` FROM `mydb`.`t`))) AS `d`"},
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
		if view.Definition != tc.def {
			t.Errorf("view %s (%s) Definition =\n  %q\nwant:\n  %q", tc.view, tc.ddl, view.Definition, tc.def)
		}
		if strings.Contains(view.Definition, "unsupported") {
			t.Errorf("view %s Definition = %q, which names a node deparse does not write", tc.view, view.Definition)
		}
	}
}

// TestAnalyzeValuesLaterRows: a multi-row VALUES is the UNION ALL of one
// SELECT per row, which is what a table value constructor is. The output names
// and arity are the first row's, as a union's first arm gives them, and every
// row's expressions are real target entries of their own arm -- so a row that
// names a table keeps that dependency, and the IR deparses back to a query of
// the same cardinality. Analyzing only the first row dropped the rest; folding
// them into the first as junk entries kept the dependency but still deparsed
// to one row.
func TestAnalyzeValuesLaterRows(t *testing.T) {
	c := newCatalogWithDB(t, "mydb")
	if err := c.DefineTable(mustParseTable(t, "CREATE TABLE t (a INT, b INT)")); err != nil {
		t.Fatal(err)
	}

	q, err := c.AnalyzeSelectStmt(mustParseSelect(t, "VALUES ROW(1,2), ROW(3,(SELECT a FROM t))"))
	if err != nil {
		t.Fatalf("AnalyzeSelectStmt: %v", err)
	}
	if q.SetOp != SetOpUnion || !q.AllSetOp {
		t.Fatalf("SetOp = %v all=%v, want a UNION ALL", q.SetOp, q.AllSetOp)
	}
	if q.LArg == nil || q.RArg == nil {
		t.Fatalf("arms = %v / %v, want one per row", q.LArg, q.RArg)
	}
	for _, arm := range []*Query{q.LArg, q.RArg} {
		if len(arm.TargetList) != 2 {
			t.Errorf("arm has %d target entries, want 2 (the row's two expressions)", len(arm.TargetList))
		}
		for _, te := range arm.TargetList {
			if te.ResJunk {
				t.Errorf("arm entry %q is junk; a row's expressions are its output", te.ResName)
			}
		}
	}
	// The second row's subquery is the one that names a table; it has to be in
	// the IR, or lineage loses t.a.
	var sublinks int
	for _, te := range q.RArg.TargetList {
		if _, ok := te.Expr.(*SubLinkExprQ); ok {
			sublinks++
		}
	}
	if sublinks != 1 {
		t.Errorf("the second arm holds %d subqueries, want the row's one", sublinks)
	}

	// The IR deparses back to a query of the same cardinality, with no empty
	// alias in the later arms.
	cases := []struct{ sql, want string }{
		{"VALUES ROW(1)", "select 1 as `column_0`"},
		{"VALUES ROW(1), ROW(2)", "select 1 as `column_0` union all select 2 as `column_0`"},
		{"VALUES ROW(1,2), ROW(3,4)",
			"select 1 as `column_0`, 2 as `column_1` union all select 3 as `column_0`, 4 as `column_1`"},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			q, err := c.AnalyzeSelectStmt(mustParseSelect(t, tc.sql))
			if err != nil {
				t.Fatalf("AnalyzeSelectStmt: %v", err)
			}
			if got := DeparseQuery(q); got != tc.want {
				t.Errorf("DeparseQuery =\n  %q\nwant:\n  %q", got, tc.want)
			}
		})
	}
}

// TestAnalyzeQueryPrimaryResultColumns: the analyzed query of a primary has the
// result columns the primary outputs, wherever it stands. A primary left
// unanalyzed produced a successful Query with no columns, which reads as a
// query that selects nothing.
func TestAnalyzeQueryPrimaryResultColumns(t *testing.T) {
	c := newCatalogWithDB(t, "mydb")
	if err := c.DefineTable(mustParseTable(t, "CREATE TABLE t (a INT, b INT)")); err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		sql  string
		want []string
	}{
		{"TABLE t", []string{"a", "b"}},
		{"TABLE t UNION TABLE t", []string{"a", "b"}},
		{"SELECT * FROM (TABLE t) d", []string{"a", "b"}},
		{"WITH c AS (TABLE t) SELECT * FROM c", []string{"a", "b"}},
		{"VALUES ROW(1,2)", []string{"column_0", "column_1"}},
		{"VALUES ROW(1,2), ROW(3,4) LIMIT 1", []string{"column_0", "column_1"}},
		{"VALUES ROW(1,2) UNION VALUES ROW(3,4)", []string{"column_0", "column_1"}},
		{"SELECT * FROM (VALUES ROW(1,2)) d", []string{"column_0", "column_1"}},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			q, err := c.AnalyzeSelectStmt(mustParseSelect(t, tc.sql))
			if err != nil {
				t.Fatalf("AnalyzeSelectStmt: %v", err)
			}
			var got []string
			for _, te := range q.TargetList {
				if te.ResJunk {
					continue // a later VALUES row's expression, not an output column
				}
				got = append(got, te.ResName)
			}
			if strings.Join(got, ",") != strings.Join(tc.want, ",") {
				t.Errorf("result columns = %v, want %v", got, tc.want)
			}
		})
	}
}

// TestDefineView_TablePrimaryBodyAnalyzed: the analyzed view body resolves the
// TABLE primary's relation, so lineage and the SDL diff see the dependency. A
// primary left unanalyzed produced a Query with an empty range table, which
// reads as a view that depends on nothing.
func TestDefineView_TablePrimaryBodyAnalyzed(t *testing.T) {
	c := newCatalogWithDB(t, "mydb")
	if err := c.DefineTable(mustParseTable(t, "CREATE TABLE t (a INT, b INT)")); err != nil {
		t.Fatal(err)
	}
	if err := c.DefineView(mustParseView(t, "CREATE VIEW v AS TABLE t")); err != nil {
		t.Fatalf("DefineView: %v", err)
	}
	view := c.GetDatabase("mydb").Views["v"]
	if view == nil {
		t.Fatal("view not registered")
	}
	q := view.AnalyzedQuery
	if q == nil {
		t.Fatal("AnalyzedQuery = nil, want the analyzed TABLE primary")
	}
	if len(q.RangeTable) != 1 {
		t.Fatalf("RangeTable holds %d entries, want 1 (the table of the TABLE primary)", len(q.RangeTable))
	}
	if len(q.TargetList) != 2 {
		t.Errorf("TargetList holds %d entries, want 2 (the columns of t)", len(q.TargetList))
	}
}
