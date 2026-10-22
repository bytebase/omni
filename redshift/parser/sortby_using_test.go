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
