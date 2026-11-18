package parser

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/oceanbase/deparse"
)

// TestFunctionExpressionGaps holds the func_expr and window_function forms the
// compat rows name: COUNT's UNIQUE spelling, CHARACTER, CHECK, WEIGHT_STRING's
// internal five-argument call, the array literal, ARRAY_AGG's ORDER BY, and
// the window forms of APPROX_COUNT_DISTINCT, TOP_K_FRE_HIST and HYBRID_HIST.
func TestFunctionExpressionGaps(t *testing.T) {
	accept := []string{
		"SELECT COUNT(UNIQUE a)",
		"SELECT CHARACTER(a)",
		"SELECT CHARACTER(a USING utf8mb4)",
		"SELECT CHECK(a)",
		"SELECT WEIGHT_STRING(a, 1, 1, 1, 1)",
		"SELECT [ a ]",
		"SELECT ARRAY_AGG(a ORDER BY a)",
		"SELECT ARRAY_AGG(DISTINCT a ORDER BY a DESC NULLS LAST)",
		"SELECT APPROX_COUNT_DISTINCT(a) OVER a",
		"SELECT APPROX_COUNT_DISTINCT_SYNOPSIS(a) OVER a",
		"SELECT APPROX_COUNT_DISTINCT_SYNOPSIS_MERGE(a) OVER a",
		"SELECT TOP_K_FRE_HIST(1.5, a, 1, 's') OVER a",
		"SELECT HYBRID_HIST(a, 1) OVER a",
		"SELECT LISTAGG(a) OVER a",
	}
	reject := []string{
		"SELECT DATE()",
		"SELECT ISNULL()",
		"SELECT CHARACTER()",
		"SELECT CHECK()",
		"SELECT ARRAY_AGG()",
		"SELECT [ ]",
		"SELECT [INTERVAL 1 DAY]",
		"SELECT DENSE_RANK() a",
		"SELECT LISTAGG(a) a",
		"SELECT LAST_VALUE(1)",
		"SELECT APPROX_COUNT_DISTINCT(a) OVER",
		"SELECT ARBITRATION(a)",
		"SELECT COUNT(UNIQUE ALL a)",
		"SELECT CHECK(a, b)",
		"SELECT TOP_K_FRE_HIST() OVER a",
		"SELECT HYBRID_HIST(a) OVER a",
		"SELECT APPROX_COUNT_DISTINCT_SYNOPSIS_MERGE(a, b) OVER a",
		"SELECT TOP_K_FRE_HIST(1.5, a, b, 's') OVER w",
		"SELECT HYBRID_HIST(a, b) OVER w",
	}
	for _, sql := range accept {
		if _, err := Parse(sql); err != nil {
			t.Errorf("rejects %q: %v", sql, err)
		}
	}
	for _, sql := range reject {
		if _, err := Parse(sql); err == nil {
			t.Errorf("accepts %q", sql)
		}
	}
}

// TestFunctionExpressionGapTrees pins the shapes those forms build.
func TestFunctionExpressionGapTrees(t *testing.T) {
	fc := firstFuncCall(t, "SELECT COUNT(UNIQUE a)")
	if fc.Quantifier != nodes.AggQuantUnique || len(fc.Args) != 1 {
		t.Errorf("COUNT(UNIQUE a): quantifier %v, %d args", fc.Quantifier, len(fc.Args))
	}

	fc = firstFuncCall(t, "SELECT ARRAY_AGG(a ORDER BY a DESC NULLS FIRST)")
	if len(fc.Args) != 1 || len(fc.OrderBy) != 1 {
		t.Fatalf("ARRAY_AGG: %d args, %d order items", len(fc.Args), len(fc.OrderBy))
	}
	if !fc.OrderBy[0].Desc || fc.OrderBy[0].NullsFirst == nil || !*fc.OrderBy[0].NullsFirst {
		t.Errorf("ARRAY_AGG order item: %+v", fc.OrderBy[0])
	}

	arr, ok := firstExpr(t, "SELECT [ a, 1 ]").(*nodes.FuncCallExpr)
	if !ok || arr.Name != "ARRAY" || len(arr.Args) != 2 {
		t.Fatalf("array literal: %#v", arr)
	}
	if arr.Loc.Start != strings.Index("SELECT [ a, 1 ]", "[") {
		t.Errorf("array literal Loc.Start = %d", arr.Loc.Start)
	}

	fc = firstFuncCall(t, "SELECT WEIGHT_STRING(a, 1, 1, 1, 1)")
	if len(fc.Args) != 5 {
		t.Errorf("WEIGHT_STRING: %d args, want 5", len(fc.Args))
	}

	fc = firstFuncCall(t, "SELECT APPROX_COUNT_DISTINCT(a) OVER a")
	if fc.Over == nil || fc.Over.RefName != "a" {
		t.Errorf("window: Over = %+v", fc.Over)
	}
}

// TestFunctionExpressionGapDeparse asserts the deparser keeps the forms it can
// read: the array literal is ARRAY's call, CHARACTER is char, and ARRAY_AGG
// keeps its ORDER BY.
func TestFunctionExpressionGapDeparse(t *testing.T) {
	for _, c := range []struct{ sql, want string }{
		{"SELECT CHARACTER(a)", "char(`a`)"},
		{"SELECT [ a, 1 ]", "array(`a`,1)"},
		{"SELECT WEIGHT_STRING(a, 1, 1, 1, 1)", "weight_string(`a`,1,1,1,1)"},
		{"SELECT ARRAY_AGG(a ORDER BY a DESC NULLS FIRST)", "array_agg(`a` order by `a` DESC NULLS FIRST)"},
		{"SELECT s.array_agg(a, b)", "s.array_agg(`a`,`b`)"},
		{"SELECT s.character(a)", "s.character(`a`)"},
	} {
		got := deparse.Deparse(firstExpr(t, c.sql))
		if got != c.want {
			t.Errorf("%q deparsed to %q, want %q", c.sql, got, c.want)
		}
		if _, err := Parse("SELECT " + got); err != nil {
			t.Errorf("%q deparsed to %q, which does not parse: %v", c.sql, got, err)
		}
	}
}
