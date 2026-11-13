package v1review

import (
	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementSelectNoSelectAll reports each bare * in the select list of every
// SELECT anywhere in a statement, routine bodies, subqueries and set
// operation sides included; a qualified t.* is not judged. Each finding
// anchors the star. SelectAll is the only code. Ported from Bytebase's V1
// rule STATEMENT_SELECT_NO_SELECT_ALL.
func StatementSelectNoSelectAll(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		ast.Inspect(stmt, func(n ast.Node) bool {
			sel, ok := n.(*ast.SelectStmt)
			if !ok || sel.SetOp != ast.SetOpNone {
				return true
			}
			for _, target := range sel.TargetList {
				if star, ok := target.(*ast.StarExpr); ok {
					start, end := statementAnchor(star.Loc)
					findings = append(findings, v1review.Finding{Code: v1review.SelectAll, Start: start, End: end})
				}
			}
			return true
		})
	}
	return findings
}
