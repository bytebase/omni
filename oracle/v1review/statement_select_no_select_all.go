package v1review

import (
	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementSelectNoSelectAll reports, in walk order, each bare * select-list
// item anywhere in the statements, PL/SQL bodies included. A qualified
// table.* is not judged, nor is a star inside an expression such as
// COUNT(*). A finding anchors the select item. SelectAll is the only code.
// Ported from Bytebase's V1 rule STATEMENT_SELECT_NO_SELECT_ALL.
func StatementSelectNoSelectAll(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		ast.Inspect(stmt, func(n ast.Node) bool {
			target, ok := n.(*ast.ResTarget)
			if !ok {
				return true
			}
			if _, ok := target.Expr.(*ast.Star); ok {
				start, end := span(target.Loc)
				findings = append(findings, v1review.Finding{Code: v1review.SelectAll, Start: start, End: end})
			}
			return true
		})
	}
	return findings
}
