package review

import (
	"strings"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/review"
)

// checkDisallowTruncate reports a TRUNCATE, naming every table it lists.
func checkDisallowTruncate(s *statement, r *reporter) {
	v, ok := s.node.(*ast.TruncateStmt)
	if !ok {
		return
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
}
