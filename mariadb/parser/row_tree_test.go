package parser

import (
	"errors"
	"strings"
	"testing"

	"github.com/bytebase/omni/mariadb/ast"
)

// selectFirstField parses SELECT <expr> and returns the expression.
func selectFirstField(t *testing.T, expr string) ast.ExprNode {
	t.Helper()
	sql := "SELECT " + expr
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("Parse(%q) read %d statements", sql, len(list.Items))
	}
	target := list.Items[0].(*ast.SelectStmt).TargetList[0]
	if rt, ok := target.(*ast.ResTarget); ok {
		return rt.Val
	}
	return target.(ast.ExprNode)
}

// TestRowHoldsItsElements: ROW(...) and a parenthesized row are the same node, whose
// Loc covers the whole of it.
func TestRowHoldsItsElements(t *testing.T) {
	for _, c := range []struct {
		expr  string
		items int
	}{
		{"ROW(1, 2)", 2},
		{"ROW(1, 2, 3)", 3},
		{"(1, 2)", 2},
		{"(1, 2, 3)", 3},
		{"ROW(ROW(1, 2), 3)", 2},
		{"((1, 2), 3)", 2},
		{"ROW(1, DEFAULT(a))", 2},
	} {
		row, ok := selectFirstField(t, c.expr).(*ast.RowExpr)
		if !ok {
			t.Errorf("%s: not a RowExpr", c.expr)
			continue
		}
		if len(row.Items) != c.items {
			t.Errorf("%s: %d items, want %d", c.expr, len(row.Items), c.items)
		}
		const prefix = len("SELECT ")
		if row.Loc.Start != prefix || row.Loc.End != prefix+len(c.expr) {
			t.Errorf("%s: Loc = %+v", c.expr, row.Loc)
		}
	}
	if _, ok := selectFirstField(t, "(1)").(*ast.ParenExpr); !ok {
		t.Error("(1) is a parenthesized expression, not a row")
	}
}

// TestRowRefusalsPoint: a refusal names the token at fault.
func TestRowRefusalsPoint(t *testing.T) {
	for _, c := range []struct {
		sql string
		at  int
	}{
		{"SELECT ROW(1)", 12},
		{"SELECT ROW(DEFAULT, 1)", 11},
		{"SELECT ROW(1, DEFAULT)", 14},
		{"SELECT ROW(*, 1)", 11},
		{"SELECT (DEFAULT, 1)", 8},
		{"SELECT (*, 1)", 8},
	} {
		_, err := Parse(c.sql)
		var perr *ParseError
		if !errors.As(err, &perr) {
			t.Errorf("Parse(%q) = %v, want a ParseError", c.sql, err)
			continue
		}
		if perr.Position != c.at {
			t.Errorf("Parse(%q) points at %d, want %d (%v)", c.sql, perr.Position, c.at, err)
		}
	}
}

// TestRowIsAName: with no '(' after it, ROW is a name wherever a name goes.
func TestRowIsAName(t *testing.T) {
	ref, ok := selectFirstField(t, "row").(*ast.ColumnRef)
	if !ok || ref.Column != "row" {
		t.Errorf("SELECT row: %+v", ref)
	}
	ref, ok = selectFirstField(t, "row.a").(*ast.ColumnRef)
	if !ok || ref.Table != "row" || ref.Column != "a" {
		t.Errorf("SELECT row.a: %+v", ref)
	}
	if _, ok := selectFirstField(t, "row(1, 2)").(*ast.RowExpr); !ok {
		t.Error("SELECT row(1, 2) is the constructor")
	}
	list, err := Parse("INSERT INTO t (row) VALUES (row)")
	if err != nil {
		t.Fatal(err)
	}
	insert := list.Items[0].(*ast.InsertStmt)
	if insert.Columns[0].Column != "row" {
		t.Errorf("Columns = %+v", insert.Columns)
	}
	if v, ok := insert.Values[0][0].(*ast.ColumnRef); !ok || v.Column != "row" {
		t.Errorf("Values = %+v", insert.Values)
	}
}

// TestIgnoreIsAValue: IGNORE reads as its own node at every position that takes one,
// and is not DEFAULT.
func TestIgnoreIsAValue(t *testing.T) {
	isIgnore := func(sql string, e ast.ExprNode, at int) {
		t.Helper()
		ig, ok := e.(*ast.IgnoreExpr)
		if !ok {
			t.Errorf("%s: %T, want *ast.IgnoreExpr", sql, e)
			return
		}
		if ig.Loc.Start != at || ig.Loc.End != at+len("ignore") {
			t.Errorf("%s: Loc = %+v, want IGNORE at %d", sql, ig.Loc, at)
		}
	}
	parse := func(sql string) ast.Node {
		t.Helper()
		list, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", sql, err)
		}
		return list.Items[0]
	}
	sql := "INSERT INTO t VALUES (ignore, default)"
	row := parse(sql).(*ast.InsertStmt).Values[0]
	isIgnore(sql, row[0], 22)
	if _, ok := row[1].(*ast.DefaultExpr); !ok {
		t.Errorf("%s: %T, want *ast.DefaultExpr", sql, row[1])
	}
	sql = "INSERT INTO t SET a = ignore"
	isIgnore(sql, parse(sql).(*ast.InsertStmt).SetList[0].Value, 22)
	sql = "INSERT INTO t VALUES (1) ON DUPLICATE KEY UPDATE a := ignore"
	isIgnore(sql, parse(sql).(*ast.InsertStmt).OnDuplicateKey[0].Value, 54)
	sql = "UPDATE t SET a = ignore, b = 1"
	isIgnore(sql, parse(sql).(*ast.UpdateStmt).SetList[0].Value, 17)
	sql = "LOAD DATA INFILE 'f' INTO TABLE t SET a = ignore"
	isIgnore(sql, parse(sql).(*ast.LoadDataStmt).SetList[0].Value, 42)
	sql = "EXECUTE s USING ignore, 1"
	isIgnore(sql, parse(sql).(*ast.ExecuteStmt).Params[0], 16)
	sql = "SELECT * FROM (VALUES (ignore)) v"
	if got := ast.NodeToString(parse(sql)); !strings.Contains(got, "{IGNORE :loc 23}") {
		t.Errorf("%s prints %s", sql, got)
	}
}
