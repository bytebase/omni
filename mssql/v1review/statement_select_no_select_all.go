package v1review

import (
	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementSelectNoSelectAll reports, in walk order, each * or qualifier.*
// select-list item of every SELECT anywhere in the statements, nested
// SELECTs and procedure bodies included. A star inside an expression, such
// as COUNT(*), is not a select-list item and is not judged. A finding
// anchors the star. SelectAll is the only code. Ported from Bytebase's V1
// rule STATEMENT_SELECT_NO_SELECT_ALL.
func StatementSelectNoSelectAll(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		ast.Inspect(stmt, func(n ast.Node) bool {
			sel, ok := n.(*ast.SelectStmt)
			if !ok || sel.TargetList == nil {
				return true
			}
			for _, item := range sel.TargetList.Items {
				var star *ast.StarExpr
				switch v := item.(type) {
				case *ast.StarExpr:
					star = v
				case *ast.ResTarget:
					star, _ = v.Val.(*ast.StarExpr)
				}
				if star == nil {
					continue
				}
				start, end := stmtAnchor(star.Loc)
				findings = append(findings, v1review.Finding{Code: v1review.SelectAll, Start: start, End: end})
			}
			return true
		})
	}
	return findings
}
