package review

import (
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/snowflake/ast"
)

// checkRequireWhere reports an UPDATE or DELETE without a WHERE clause,
// wherever the tree holds it: top level, or inside a Snowflake Scripting
// block, a task's body, or an alert's action. A routine's body and an
// EXECUTE IMMEDIATE string are text the parser does not read.
func checkRequireWhere(s *statement, r *reporter) {
	ast.Inspect(s.node, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.UpdateStmt:
			if v.Where == nil {
				r.report(review.RequireWhere, s.index, anchor(s, v, v.Loc), "UPDATE "+object(v.Target)+" has no WHERE clause")
			}
		case *ast.DeleteStmt:
			if v.Where == nil {
				r.report(review.RequireWhere, s.index, anchor(s, v, v.Loc), "DELETE FROM "+object(v.Target)+" has no WHERE clause")
			}
		}
		return true
	})
}
