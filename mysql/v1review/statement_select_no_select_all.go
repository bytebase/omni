package v1review

import (
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementSelectNoSelectAll reports each * in the select list of any SELECT
// in a top-level statement, subqueries, CTEs, INSERT ... SELECT and CREATE
// ... AS SELECT included. A qualified t.* is not reported, nor is a * in
// a function call such as COUNT(*). Each finding anchors the star.
// SelectAll is the only code. Ported from Bytebase's V1 rule
// STATEMENT_SELECT_NO_SELECT_ALL.
func StatementSelectNoSelectAll(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		ast.Inspect(stmt, func(n ast.Node) bool {
			sel, ok := n.(*ast.SelectStmt)
			if !ok || sel.SetOp != ast.SetOpNone {
				return true
			}
			for _, target := range sel.TargetList {
				var star *ast.StarExpr
				switch e := target.(type) {
				case *ast.StarExpr:
					star = e
				case *ast.ResTarget:
					star, _ = e.Val.(*ast.StarExpr)
				}
				if star != nil {
					start, end := anchor(star.Loc)
					findings = append(findings, v1review.Finding{Code: v1review.SelectAll, Start: start, End: end})
				}
			}
			return true
		})
	}
	return findings
}
