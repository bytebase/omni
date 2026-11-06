package review

import (
	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/review"
)

// checkRequireWhere reports an UPDATE or DELETE without a WHERE clause,
// wherever it sits: top level, or the body of a routine, trigger, or
// event. A statement under EXPLAIN is not judged, but one under ANALYZE is: MariaDB runs it
// (mariadb:11.8.8 updates the rows of ANALYZE UPDATE t SET a = 2).
func checkRequireWhere(s *statement, r *reporter) {
	ast.Inspect(s.node, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.ExplainStmt:
			return v.Analyze
		case *ast.UpdateStmt:
			if v.Where == nil {
				r.report(review.MissingWhere, s.index, rangeOf(v.Loc), noWhere("UPDATE", "of", tableExprs(v.Tables)))
			}
		case *ast.DeleteStmt:
			if v.Where == nil {
				r.report(review.MissingWhere, s.index, rangeOf(v.Loc), noWhere("DELETE", "from", tableExprs(deleted(v))))
			}
		}
		return true
	})
}

// noWhere writes the finding for an UPDATE or DELETE without a WHERE
// clause: "`UPDATE` of `t` has no `WHERE` clause.", naming the tables it
// updates or deletes from, and none when it names no table.
func noWhere(verb, prep string, names []string) string {
	subject := review.Code(verb)
	if len(names) > 0 {
		subject += " " + prep + " " + review.Names(names...)
	}
	return subject + " has no " + review.Code("WHERE") + " clause."
}

// deleted returns the tables a DELETE names as the ones it deletes from:
// the list after DELETE or FROM, which in the multiple-table form lists
// only the targets.
func deleted(v *ast.DeleteStmt) []ast.TableExpr {
	return v.Tables
}
