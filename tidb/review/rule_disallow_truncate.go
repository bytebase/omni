package review

import (
	"strings"

	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/tidb/ast"
)

// checkDisallowTruncate reports a TRUNCATE, naming every table it lists,
// and ALTER TABLE ... TRUNCATE PARTITION, at top level or inside a routine
// body.
func checkDisallowTruncate(s *statement, r *reporter) {
	ast.Inspect(s.node, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.TruncateStmt:
			r.report(review.DisallowTruncate, s.index, rangeOf(v.Loc), "truncates "+tables(v.Tables))
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
					partitions = "partition " + strings.Join(names, ", ")
				}
				r.report(review.DisallowTruncate, s.index, rangeOf(cmd.Loc), "truncates "+partitions+" of "+table(v.Table))
			}
		}
		return true
	})
}
