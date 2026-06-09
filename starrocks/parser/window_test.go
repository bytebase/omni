package parser

import (
	"testing"

	"github.com/bytebase/omni/starrocks/ast"
)

func parseOneSelect(t *testing.T, input string) *ast.SelectStmt {
	t.Helper()
	file, errs := Parse(input)
	if len(errs) > 0 {
		t.Fatalf("Parse(%q) errors: %v", input, errs)
	}
	if len(file.Stmts) == 0 {
		t.Fatalf("Parse(%q): no statements", input)
	}
	sel, ok := file.Stmts[0].(*ast.SelectStmt)
	if !ok {
		t.Fatalf("Parse(%q): got %T, want *ast.SelectStmt", input, file.Stmts[0])
	}
	return sel
}

// A window function in a top-level SELECT must be fully consumed, so the
// trailing FROM clause is reachable and the OVER spec is captured.
func TestWindowFunctionTopLevel(t *testing.T) {
	sel := parseOneSelect(t, "SELECT region, ROW_NUMBER() OVER (PARTITION BY region ORDER BY amount DESC) AS rn FROM sales")

	if len(sel.From) != 1 {
		t.Fatalf("From = %d items, want 1 (FROM must parse after the window function)", len(sel.From))
	}
	if len(sel.Items) != 2 {
		t.Fatalf("Items = %d, want 2", len(sel.Items))
	}
	fc, ok := sel.Items[1].Expr.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("Items[1].Expr = %T, want *ast.FuncCallExpr", sel.Items[1].Expr)
	}
	if fc.Over == nil {
		t.Fatal("FuncCallExpr.Over = nil, want non-nil window spec")
	}
	if len(fc.Over.PartitionBy) != 1 {
		t.Errorf("Over.PartitionBy = %d, want 1", len(fc.Over.PartitionBy))
	}
	if len(fc.Over.OrderBy) != 1 {
		t.Errorf("Over.OrderBy = %d, want 1", len(fc.Over.OrderBy))
	}
	if sel.Items[1].Alias != "rn" {
		t.Errorf("Items[1].Alias = %q, want %q", sel.Items[1].Alias, "rn")
	}
}

// Regression for BYT-9636: a window function inside a CTE body used to fail
// with "syntax error at or near OVER" because OVER was never consumed, so the
// CTE's closing ')' collided with it.
func TestWindowFunctionInCTE(t *testing.T) {
	input := "WITH ranked AS (SELECT region, ROW_NUMBER() OVER (PARTITION BY region ORDER BY amount DESC) AS rn FROM sales) SELECT region FROM ranked WHERE rn = 1"
	file, errs := Parse(input)
	if len(errs) > 0 {
		t.Fatalf("Parse(%q) errors: %v", input, errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("Stmts = %d, want 1", len(file.Stmts))
	}
}

func TestWindowFunctionFrame(t *testing.T) {
	sel := parseOneSelect(t, "SELECT SUM(amount) OVER (PARTITION BY region ORDER BY sale_date ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) FROM sales")
	fc, ok := sel.Items[0].Expr.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("Items[0].Expr = %T, want *ast.FuncCallExpr", sel.Items[0].Expr)
	}
	if fc.Over == nil || fc.Over.Frame == nil {
		t.Fatal("expected non-nil Over.Frame")
	}
	if fc.Over.Frame.Unit != "ROWS" {
		t.Errorf("Frame.Unit = %q, want ROWS", fc.Over.Frame.Unit)
	}
}

func TestWindowFunctionEmptyOver(t *testing.T) {
	sel := parseOneSelect(t, "SELECT RANK() OVER () FROM sales")
	fc, ok := sel.Items[0].Expr.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("Items[0].Expr = %T, want *ast.FuncCallExpr", sel.Items[0].Expr)
	}
	if fc.Over == nil {
		t.Fatal("Over = nil, want non-nil for empty OVER ()")
	}
}
