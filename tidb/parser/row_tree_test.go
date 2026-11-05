package parser

import (
	"errors"
	"testing"

	"github.com/bytebase/omni/tidb/ast"
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

// TestParenthesizedRowOperands: the corpus rows this reading closes hold a row, not a
// statement that happened to parse.
func TestParenthesizedRowOperands(t *testing.T) {
	cmp, ok := selectFirstField(t, "(1, 1) > (1, 1)").(*ast.BinaryExpr)
	if !ok {
		t.Fatal("not a comparison")
	}
	if _, ok := cmp.Left.(*ast.RowExpr); !ok {
		t.Errorf("Left = %T", cmp.Left)
	}
	if _, ok := cmp.Right.(*ast.RowExpr); !ok {
		t.Errorf("Right = %T", cmp.Right)
	}
	list, err := Parse("SELECT * FROM t GROUP BY (a, b) WITH ROLLUP")
	if err != nil {
		t.Fatal(err)
	}
	sel := list.Items[0].(*ast.SelectStmt)
	if len(sel.GroupBy) != 1 {
		t.Fatalf("GroupBy = %+v", sel.GroupBy)
	}
	if _, ok := sel.GroupBy[0].(*ast.RowExpr); !ok {
		t.Errorf("GroupBy[0] = %T", sel.GroupBy[0])
	}
}
