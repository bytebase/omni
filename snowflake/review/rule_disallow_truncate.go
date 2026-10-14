package review

import (
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/snowflake/ast"
)

// checkDisallowTruncate reports a TRUNCATE in each of its forms: TRUNCATE
// [TABLE] of a table, of the error table of a table (ERROR_TABLE(t)), and
// TRUNCATE MATERIALIZED VIEW, which removes every row of the view; at top
// level or inside a Snowflake Scripting block, a task's body, or an alert's
// action. IF EXISTS does not exempt a statement.
func checkDisallowTruncate(s *statement, r *reporter) {
	ast.Inspect(s.node, func(n ast.Node) bool {
		if _, ok := n.(*ast.ExplainStmt); ok {
			// EXPLAIN compiles its statement without running it.
			return false
		}
		if v, ok := n.(*ast.TruncateStmt); ok {
			what := object(v.Name)
			switch {
			case v.MaterializedView:
				what = "materialized view " + what
			case v.ErrorTable:
				what = "the error table of " + what
			}
			r.report(review.DisallowTruncate, s.index, anchor(s, v, v.Loc), "truncates "+what)
		}
		return true
	})
}
