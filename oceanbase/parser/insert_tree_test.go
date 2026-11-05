package parser

import (
	"errors"
	"regexp"
	"testing"

	"github.com/bytebase/omni/oceanbase/ast"
)

// The tree of the forms insert_forms_test.go holds to the server. A verdict cannot
// see any of this: a parser could accept every row there and still lose the
// difference between () and no list, or hang the rows on the wrong clause.

func insertTree(t *testing.T, sql string) *ast.InsertStmt {
	t.Helper()
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("Parse(%q) read %d statements, want 1", sql, len(list.Items))
	}
	stmt, ok := list.Items[0].(*ast.InsertStmt)
	if !ok {
		t.Fatalf("Parse(%q) read a %T", sql, list.Items[0])
	}
	return stmt
}

func rowLengths(rows [][]ast.ExprNode) []int {
	out := make([]int, len(rows))
	for i, r := range rows {
		out[i] = len(r)
	}
	return out
}

func sameInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// TestInsertColumnListAndRows pins how an empty list reads. No list is nil and ()
// is non-nil and empty; a row of () is a row with no elements.
func TestInsertColumnListAndRows(t *testing.T) {
	const noList = -1
	for _, c := range []struct {
		sql     string
		columns int
		rows    []int
	}{
		{"INSERT INTO t VALUES ()", noList, []int{0}},
		{"INSERT INTO t VALUE ()", noList, []int{0}},
		{"REPLACE INTO t VALUES ()", noList, []int{0}},
		{"INSERT INTO t () VALUES ()", 0, []int{0}},
		{"REPLACE INTO t () VALUE ()", 0, []int{0}},
		{"INSERT INTO t () VALUES (1), ()", 0, []int{1, 0}},
		{"INSERT INTO t VALUES (), (1, 2), ()", noList, []int{0, 2, 0}},
		{"INSERT INTO t (a) VALUES ()", 1, []int{0}},
		{"INSERT INTO t (a, b) VALUES (1, 2), ()", 2, []int{2, 0}},
	} {
		t.Run(c.sql, func(t *testing.T) {
			stmt := insertTree(t, c.sql)
			switch {
			case c.columns == noList && stmt.Columns != nil:
				t.Errorf("Columns = %v, want nil", stmt.Columns)
			case c.columns != noList && stmt.Columns == nil:
				t.Errorf("Columns = nil, want %d columns", c.columns)
			case c.columns != noList && len(stmt.Columns) != c.columns:
				t.Errorf("len(Columns) = %d, want %d", len(stmt.Columns), c.columns)
			}
			if got := rowLengths(stmt.Values); !sameInts(got, c.rows) {
				t.Errorf("rows = %v, want %v", got, c.rows)
			}
			if stmt.Select != nil {
				t.Errorf("Select = %+v, want nil", stmt.Select)
			}
			if stmt.Loc.Start != 0 || stmt.Loc.End != len(c.sql) {
				t.Errorf("Loc = %+v, want the whole statement of %d bytes", stmt.Loc, len(c.sql))
			}
		})
	}

	stmt := insertTree(t, "INSERT INTO t () SELECT 1")
	if stmt.Columns == nil || len(stmt.Columns) != 0 {
		t.Errorf("() before a query: Columns = %#v, want non-nil and empty", stmt.Columns)
	}
	if stmt.Select == nil || stmt.Values != nil {
		t.Errorf("() before a query: Select = %v, Values = %v", stmt.Select, stmt.Values)
	}
	if stmt := insertTree(t, "INSERT INTO t SELECT 1"); stmt.Columns != nil {
		t.Errorf("no list before a query: Columns = %#v, want nil", stmt.Columns)
	}
}

// TestInsertColumnListPrints keeps the two apart in the printed tree too.
func TestInsertColumnListPrints(t *testing.T) {
	columns := regexp.MustCompile(` :columns `)
	empty := regexp.MustCompile(` :columns \(\)`)
	if got := ast.NodeToString(insertTree(t, "INSERT INTO t () VALUES ()")); !empty.MatchString(got) {
		t.Errorf("() does not print as an empty list: %s", got)
	}
	if got := ast.NodeToString(insertTree(t, "INSERT INTO t VALUES ()")); columns.MatchString(got) {
		t.Errorf("no list prints as a list: %s", got)
	}
	got := ast.NodeToString(insertTree(t, "INSERT INTO t (a) VALUES (1)"))
	if !columns.MatchString(got) || empty.MatchString(got) {
		t.Errorf("(a) does not print its column: %s", got)
	}
}

// TestInsertEmptyListsWalk reads through an empty list and row without visiting a
// node that is not there.
func TestInsertEmptyListsWalk(t *testing.T) {
	stmt := insertTree(t, "INSERT INTO t () VALUES ()")
	count := 0
	ast.Inspect(stmt, func(n ast.Node) bool {
		if n == nil {
			t.Error("Inspect visited a nil node")
		}
		count++
		return true
	})
	if count != 2 { // the statement and its table
		t.Errorf("visited %d nodes, want the statement and its table", count)
	}
}

// TestParenthesizedQuerySourceTakesClauses: the source is read like any other, so
// what follows it is not left for a second statement.
func TestParenthesizedQuerySourceTakesClauses(t *testing.T) {
	const sql = "INSERT INTO t (SELECT 1) ON DUPLICATE KEY UPDATE a = 1"
	stmt := insertTree(t, sql)
	if stmt.Columns != nil {
		t.Errorf("Columns = %#v, want nil", stmt.Columns)
	}
	if stmt.Select == nil || len(stmt.OnDuplicateKey) != 1 {
		t.Errorf("Select = %v, OnDuplicateKey = %v", stmt.Select, stmt.OnDuplicateKey)
	}
	if stmt.Loc.End != len(sql) {
		t.Errorf("Loc.End = %d, want %d", stmt.Loc.End, len(sql))
	}
}

// TestInsertFormErrorsPoint: a refusal names the token at fault.
func TestInsertFormErrorsPoint(t *testing.T) {
	for _, c := range []struct {
		sql string
		at  int
	}{
		{"INSERT INTO t (a) SET a = 1", 18},
		{"INSERT INTO t VALUES (*)", 22},
		{"UPDATE t SET @a = 1", 13},
		{"INSERT INTO t VALUES (1), ROW(2)", 26},
		{"INSERT INTO t VALUE ROW(1)", 20},
		{"INSERT INTO t (@a) VALUES (1)", 15},
		{"UPDATE t SET a := 1", 15},
		{"INSERT INTO t VALUES (,,1)", 23},
	} {
		t.Run(c.sql, func(t *testing.T) {
			_, err := Parse(c.sql)
			var perr *ParseError
			if !errors.As(err, &perr) {
				t.Fatalf("Parse(%q) = %v, want a ParseError", c.sql, err)
			}
			if perr.Position != c.at {
				t.Errorf("Parse(%q) points at %d, want %d (%v)", c.sql, perr.Position, c.at, err)
			}
		})
	}
}

// TestRowCompletion: the caret inside an empty row still offers what a value can be.
func TestRowCompletion(t *testing.T) {
	for _, c := range []struct {
		sql string
		at  int
	}{
		{"insert into t values ()", 22},
		{"insert into t () values ()", 25},
		{"insert into t ()", 15},
	} {
		if got := Collect(c.sql, c.at); len(got.Rules) == 0 {
			t.Errorf("Collect(%q, %d) offered no rules", c.sql, c.at)
		}
	}
}

// TestQuotedAtNameIsAColumn: only the bare @a is a variable.
func TestQuotedAtNameIsAColumn(t *testing.T) {
	list, err := Parse("UPDATE t SET `@a` = 1")
	if err != nil {
		t.Fatal(err)
	}
	if got := list.Items[0].(*ast.UpdateStmt).SetList[0].Column.Column; got != "@a" {
		t.Errorf("column = %q, want @a", got)
	}
	if got := insertTree(t, "INSERT INTO t (`@a`) VALUES (1)").Columns[0].Column; got != "@a" {
		t.Errorf("column = %q, want @a", got)
	}
}

// TestValuesRowSource: bare ROW rows stay in Values, where a reader finds the rows of
// any plain INSERT. Anything that makes them more than rows makes them the query
// they are, in Select.
func TestValuesRowSource(t *testing.T) {
	for _, c := range []struct {
		sql  string
		rows []int
	}{
		{"INSERT INTO t VALUES ROW(1), ROW(2, 3)", []int{1, 2}},
		{"INSERT INTO t VALUES ROW()", []int{0}},
		{"REPLACE INTO t (a) VALUES ROW(1)", []int{1}},
		{"INSERT INTO t () VALUES ROW(1) ON DUPLICATE KEY UPDATE a = 1", []int{1}},
	} {
		stmt := insertTree(t, c.sql)
		if got := rowLengths(stmt.Values); !sameInts(got, c.rows) || stmt.Select != nil {
			t.Errorf("%s: rows = %v, Select = %v, want rows %v", c.sql, got, stmt.Select, c.rows)
		}
	}

	check := func(sql string, ok func(*ast.SelectStmt) bool) {
		t.Helper()
		stmt := insertTree(t, sql)
		if stmt.Values != nil || stmt.Select == nil || !ok(stmt.Select) {
			t.Errorf("%s: Values = %v, Select = %+v", sql, stmt.Values, stmt.Select)
		}
	}
	check("INSERT INTO t VALUES ROW(1) ORDER BY 1", func(s *ast.SelectStmt) bool {
		return s.ValuesSource != nil && len(s.OrderBy) == 1
	})
	check("INSERT INTO t VALUES ROW(1) LIMIT 1", func(s *ast.SelectStmt) bool {
		return s.ValuesSource != nil && s.Limit != nil
	})
	check("INSERT INTO t VALUES ROW(1) UNION VALUES ROW(2)", func(s *ast.SelectStmt) bool {
		return s.SetOp != ast.SetOpNone && s.Left.ValuesSource != nil && s.Right.ValuesSource != nil
	})
	check("INSERT INTO t (VALUES ROW(1))", func(s *ast.SelectStmt) bool {
		return s.ParenSource != nil && s.ParenSource.ValuesSource != nil
	})
}

// TestRowAliasHoldsItsColumns: the alias stays with the rows or SET it follows.
func TestRowAliasHoldsItsColumns(t *testing.T) {
	stmt := insertTree(t, "INSERT INTO t () VALUES () AS new (a, b) ON DUPLICATE KEY UPDATE a = new.a")
	if stmt.RowAlias != "new" || len(stmt.ColAliases) != 2 || len(stmt.OnDuplicateKey) != 1 || len(stmt.Values) != 1 {
		t.Errorf("%+v", stmt)
	}
	if stmt := insertTree(t, "INSERT INTO t SET a = 1 AS new"); stmt.RowAlias != "new" || len(stmt.SetList) != 1 {
		t.Errorf("%+v", stmt)
	}
}

// TestRowLeadingComma: insert_vals may open with its empty alternative, which holds
// no value, so the row is what follows the comma.
func TestRowLeadingComma(t *testing.T) {
	for _, c := range []struct {
		sql  string
		rows []int
	}{
		{"INSERT INTO t VALUES (,1)", []int{1}},
		{"INSERT INTO t VALUES (,1,2)", []int{2}},
		{"INSERT INTO t VALUES (,1), (,2)", []int{1, 1}},
		{"INSERT INTO t VALUES ROW(,1)", []int{1}},
	} {
		if got := rowLengths(insertTree(t, c.sql).Values); !sameInts(got, c.rows) {
			t.Errorf("%s: rows = %v, want %v", c.sql, got, c.rows)
		}
	}
	lit, ok := insertTree(t, "INSERT INTO t VALUES (,1)").Values[0][0].(*ast.IntLit)
	if !ok || lit.Value != 1 {
		t.Errorf("the row holds %v, want 1", insertTree(t, "INSERT INTO t VALUES (,1)").Values[0][0])
	}
}

// TestReplaceTakesWhatInsertTakes: OceanBase's REPLACE rule is replace_with_opt_hint
// opt_low_priority opt_ignore opt_into single_table_insert, so it takes IGNORE and a
// row alias, which MySQL's does not.
func TestReplaceTakesWhatInsertTakes(t *testing.T) {
	if stmt := insertTree(t, "REPLACE IGNORE INTO t VALUES (1)"); !stmt.IsReplace || !stmt.Ignore {
		t.Errorf("%+v", stmt)
	}
	if stmt := insertTree(t, "REPLACE INTO t VALUES (1) AS new"); !stmt.IsReplace || stmt.RowAlias != "new" {
		t.Errorf("%+v", stmt)
	}
}

// TestRowStarOfTable: OceanBase's parser takes t.* in a row, and its resolver refuses
// it, with 1149 rather than 1064.
func TestRowStarOfTable(t *testing.T) {
	ref, ok := insertTree(t, "INSERT INTO t VALUES (t.*)").Values[0][0].(*ast.ColumnRef)
	if !ok || !ref.Star || ref.Table != "t" {
		t.Errorf("the row holds %+v", ref)
	}
}
