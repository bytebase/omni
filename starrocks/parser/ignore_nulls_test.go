package parser

import (
	"testing"

	"github.com/bytebase/omni/starrocks/ast"
)

// IGNORE NULLS null-treatment in window functions. Two grammar positions:
// inside the parens after the first argument, and after the closing paren.

func ignoreNullsFunc(t *testing.T, input string) *ast.FuncCallExpr {
	t.Helper()
	stmt := mustParseSelect(t, input)
	fc, ok := stmt.Items[0].Expr.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("item expr = %T, want *ast.FuncCallExpr", stmt.Items[0].Expr)
	}
	return fc
}

// TestIgnoreNullsInside covers IGNORE NULLS after the (only) argument.
func TestIgnoreNullsInside(t *testing.T) {
	fc := ignoreNullsFunc(t, "SELECT FIRST_VALUE(v IGNORE NULLS) OVER (PARTITION BY g ORDER BY ts) FROM t")
	if !fc.IgnoreNulls {
		t.Error("IgnoreNulls = false, want true")
	}
	if fc.Over == nil {
		t.Error("Over = nil, want the OVER clause to be parsed")
	}
	if len(fc.Args) != 1 {
		t.Errorf("args = %d, want 1", len(fc.Args))
	}
}

// TestIgnoreNullsOutside covers IGNORE NULLS after the closing paren. (Without
// the feature, IGNORE is silently mis-consumed as a column alias.)
func TestIgnoreNullsOutside(t *testing.T) {
	fc := ignoreNullsFunc(t, "SELECT LEAD(v, 1) IGNORE NULLS OVER (ORDER BY ts) FROM t")
	if !fc.IgnoreNulls {
		t.Error("IgnoreNulls = false, want true")
	}
	if fc.Over == nil {
		t.Error("Over = nil, want the OVER clause to be parsed")
	}
	if stmt := mustParseSelect(t, "SELECT LEAD(v, 1) IGNORE NULLS OVER (ORDER BY ts) FROM t"); stmt.Items[0].Alias != "" {
		t.Errorf("alias = %q, want empty (IGNORE must not be consumed as an alias)", stmt.Items[0].Alias)
	}
}

// TestIgnoreNullsMidArgs covers IGNORE NULLS after the first arg with more args
// following (LAG(v IGNORE NULLS, 1)).
func TestIgnoreNullsMidArgs(t *testing.T) {
	fc := ignoreNullsFunc(t, "SELECT LAG(v IGNORE NULLS, 1) OVER (ORDER BY ts) FROM t")
	if !fc.IgnoreNulls {
		t.Error("IgnoreNulls = false, want true")
	}
	if len(fc.Args) != 2 {
		t.Errorf("args = %d, want 2 (v, 1)", len(fc.Args))
	}
}

// TestFuncCallNoIgnoreNulls is the regression guard: a plain window function
// still parses with IgnoreNulls=false.
func TestFuncCallNoIgnoreNulls(t *testing.T) {
	fc := ignoreNullsFunc(t, "SELECT FIRST_VALUE(v) OVER (PARTITION BY g ORDER BY ts) FROM t")
	if fc.IgnoreNulls {
		t.Error("IgnoreNulls = true, want false")
	}
}
