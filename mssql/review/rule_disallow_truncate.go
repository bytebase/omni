package review

import (
	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/review"
)

// checkDisallowTruncate reports a TRUNCATE TABLE, of the whole table or
// of partitions, at top level or inside a routine body.
func checkDisallowTruncate(s *statement, r *reporter) {
	ast.Inspect(s.node, func(n ast.Node) bool {
		if v, ok := n.(*ast.TruncateStmt); ok {
			what := object(v.Table)
			if v.Partitions != nil && len(v.Partitions.Items) > 0 {
				what = "partitions of " + what
			}
			r.report(review.DisallowTruncate, s.index, rangeOf(v.Loc), "truncates "+what)
		}
		return true
	})
}
