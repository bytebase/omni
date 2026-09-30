package review

import (
	"strings"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/review"
)

// checkDisallowTruncate reports a TRUNCATE, naming every table it lists,
// at top level or inside a SQL-standard routine body.
func checkDisallowTruncate(s *statement, r *reporter) {
	ast.Inspect(s.node, func(n ast.Node) bool {
		v, ok := n.(*ast.TruncateStmt)
		if !ok {
			return true
		}
		var names []string
		if v.Relations != nil {
			for _, item := range v.Relations.Items {
				if rv, ok := item.(*ast.RangeVar); ok {
					names = append(names, relation(rv))
				}
			}
		}
		r.report(review.DisallowTruncate, s.index, rangeOf(v.Loc), "truncates "+strings.Join(names, ", "))
		return true
	})
}
