package review

import (
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/tidb/ast"
)

// checkRequireWhere reports an UPDATE or DELETE without a WHERE clause,
// wherever it sits: top level, the body of a routine, trigger, or event,
// a BATCH, or one of the three wrappers in the form that runs the statement
// -- EXPLAIN ANALYZE, TRACE, and PLAN REPLAYER DUMP EXPLAIN ANALYZE. A
// statement under a plain EXPLAIN, a TRACE PLAN, a plain PLAN REPLAYER DUMP,
// a BATCH dry run, or a SQL plan binding, none of which runs it, is not
// judged.
func checkRequireWhere(s *statement, r *reporter) {
	ast.Inspect(s.node, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.ExplainStmt, *ast.TraceStmt, *ast.PlanReplayerStmt:
			return runsWrapped(v)
		case *ast.BatchStmt:
			return v.DryRun == ast.BatchDryRunNone
		case *ast.CreateBindingStmt, *ast.DropBindingStmt, *ast.SetBindingStmt:
			// A SQL plan binding records a plan for a statement shape and
			// runs neither statement it carries, so the UPDATE or DELETE
			// inside one updates and deletes nothing.
			return false
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
