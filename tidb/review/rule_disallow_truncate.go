package review

import (
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/tidb/ast"
)

// checkDisallowTruncate reports a TRUNCATE, naming every table it lists,
// and ALTER TABLE ... TRUNCATE PARTITION, at top level or inside a routine
// body.
func checkDisallowTruncate(s *statement, r *reporter) {
	ast.Inspect(s.node, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.ExplainStmt, *ast.TraceStmt, *ast.PlanReplayerStmt:
			// A plain EXPLAIN reports a plan, TRACE PLAN traces the
			// optimizer and a plain PLAN REPLAYER DUMP writes a dump, so
			// what any of the three carries is no change of this kind. The
			// form that runs the statement is followed; wrapped says which
			// form that is.
			return runsWrapped(v)
		case *ast.TruncateStmt:
			r.report(review.TruncatesData, s.index, rangeOf(v.Loc), "Truncates "+review.Names(tables(v.Tables)...)+".")
		case *ast.AlterTableStmt:
			for _, cmd := range v.Commands {
				if cmd == nil || cmd.Type != ast.ATTruncatePartition {
					continue
				}
				partitions := "all partitions"
				if !cmd.AllPartitions {
					names := make([]string, len(cmd.PartitionNames))
					for i, name := range cmd.PartitionNames {
						names[i] = ident(name)
					}
					partitions = review.Plural("partition", len(names)) + " " + review.Names(names...)
				}
				r.report(review.TruncatesData, s.index, rangeOf(cmd.Loc), "Truncates "+partitions+" of "+review.Code(table(v.Table))+".")
			}
		}
		return true
	})
}
