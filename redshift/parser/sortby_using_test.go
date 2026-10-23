package parser

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
)

// TestSortByUsingOperator checks that ORDER BY ... USING keeps its operator
// and that a missing operator is a syntax error, as in the pg parser this
// package forks.
//
// gram.y: sortby — a_expr USING qual_all_Op opt_nulls_order
func TestSortByUsingOperator(t *testing.T) {
	for sql, want := range map[string]string{
		"SELECT 1 ORDER BY 1 USING >":                      ">",
		"SELECT 1 ORDER BY 1 USING OPERATOR(pg_catalog.<)": "pg_catalog.<",
	} {
		list, err := Parse(sql)
		if err != nil {
			t.Fatalf("%s: %v", sql, err)
		}
		sb := list.Items[0].(*nodes.RawStmt).Stmt.(*nodes.SelectStmt).SortClause.Items[0].(*nodes.SortBy)
		var names []string
		if sb.UseOp != nil {
			for _, item := range sb.UseOp.Items {
				names = append(names, item.(*nodes.String).Str)
			}
		}
		if got := strings.Join(names, "."); got != want {
			t.Errorf("%s: UseOp %q, want %q", sql, got, want)
		}
	}
	for _, sql := range []string{"SELECT 1 ORDER BY 1 USING", "SELECT 1 ORDER BY 1 USING NULLS FIRST"} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("%s: parsed, want a syntax error", sql)
		}
	}
}

// TestWithinGroupClause checks that WITHIN GROUP requires its ORDER BY and
// rejects another ORDER BY and VARIADIC, as in the pg parser, but takes
// DISTINCT, as Redshift's LISTAGG does.
//
// gram.y: within_group_clause; func_expr (the checks on WITHIN GROUP)
// https://docs.aws.amazon.com/redshift/latest/dg/r_LISTAGG.html
func TestWithinGroupClause(t *testing.T) {
	for _, sql := range []string{
		"SELECT percentile_cont(0.5) WITHIN GROUP (ORDER BY 1)",
		"SELECT LISTAGG(DISTINCT product_category, ',') WITHIN GROUP (ORDER BY product_category) FROM t",
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("%s: %v", sql, err)
		}
	}
	for _, sql := range []string{
		"SELECT percentile_cont(0.5) WITHIN GROUP ()",
		"SELECT percentile_cont(VARIADIC ARRAY[0.5]) WITHIN GROUP (ORDER BY 1)",
		"SELECT percentile_cont(0.5 ORDER BY 1) WITHIN GROUP (ORDER BY 1)",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("%s: parsed, want a syntax error", sql)
		}
	}
}
