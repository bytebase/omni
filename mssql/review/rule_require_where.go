package review

import (
	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/review"
)

// checkRequireWhere reports an UPDATE or DELETE without a WHERE clause,
// wherever it sits: top level, or the body of a procedure, function, or
// trigger, or a control-of-flow block.
func checkRequireWhere(s *statement, r *reporter) {
	ast.Inspect(s.node, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.UpdateStmt:
			if v.WhereClause == nil {
				r.report(review.RequireWhere, s.index, rangeOf(v.Loc), "`UPDATE` of "+whereTarget(v.Relation)+" has no `WHERE` clause.")
			}
		case *ast.DeleteStmt:
			if v.WhereClause == nil {
				r.report(review.RequireWhere, s.index, rangeOf(v.Loc), "`DELETE` from "+whereTarget(v.Relation)+" has no `WHERE` clause.")
			}
		}
		return true
	})
}

// whereTarget writes the table an UPDATE or DELETE names as its target,
// as a code span.
func whereTarget(e ast.TableExpr) string {
	if name := tableExpr(e); name != "" {
		return review.Code(name)
	}
	return "the target"
}
