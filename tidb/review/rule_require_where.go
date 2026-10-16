package review

import (
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/tidb/ast"
)

// checkRequireWhere reports an UPDATE or DELETE without a WHERE clause,
// wherever it sits: top level, the body of a routine, trigger, or event,
// a BATCH, or EXPLAIN ANALYZE, which runs the statement. A statement under
// a plain EXPLAIN, or a BATCH dry run, which runs nothing, is not judged.
func checkRequireWhere(s *statement, r *reporter) {
	ast.Inspect(s.node, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.ExplainStmt:
			return v.Analyze
		case *ast.BatchStmt:
			return v.DryRun == ast.BatchDryRunNone
		case *ast.UpdateStmt:
			if v.Where == nil {
				r.report(review.RequireWhere, s.index, rangeOf(v.Loc), "UPDATE "+tableExprs(v.Tables)+" has no WHERE clause")
			}
		case *ast.DeleteStmt:
			if v.Where == nil {
				r.report(review.RequireWhere, s.index, rangeOf(v.Loc), "DELETE FROM "+tableExprs(deleted(v))+" has no WHERE clause")
			}
		}
		return true
	})
}

// deleted returns the tables a DELETE names as the ones it deletes from:
// the list after DELETE or FROM, which in the multiple-table form lists
// only the targets.
func deleted(v *ast.DeleteStmt) []ast.TableExpr {
	return v.Tables
}
