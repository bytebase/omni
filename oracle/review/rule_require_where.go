package review

import (
	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/review"
)

// checkRequireWhere reports an UPDATE or DELETE without a WHERE clause,
// wherever it sits: top level, the body of a procedure, function,
// package, type, or trigger, an anonymous block, or a FORALL. The
// statement EXPLAIN PLAN explains does not run, so it is not judged. A
// MERGE has no WHERE of its own to require, and the text EXECUTE
// IMMEDIATE runs is not read.
func checkRequireWhere(s *statement, r *reporter) {
	ast.Inspect(s.node, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.ExplainPlanStmt:
			return false
		case *ast.UpdateStmt:
			if v.WhereClause == nil {
				r.report(review.RequireWhere, s.index, rangeOf(v.Loc), "UPDATE "+dmlTarget(v.Target)+" has no WHERE clause")
			}
		case *ast.DeleteStmt:
			if v.WhereClause == nil {
				r.report(review.RequireWhere, s.index, rangeOf(v.Loc), "DELETE FROM "+dmlTarget(v.Target)+" has no WHERE clause")
			}
		}
		return true
	})
}
