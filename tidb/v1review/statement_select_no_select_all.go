package v1review

import (
	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementSelectNoSelectAll reports, in walk order, each SELECT anywhere in
// the statements whose select list holds * or table.*, once per SELECT
// however many stars it holds: subqueries, set-operation arms, INSERT ...
// SELECT and EXPLAIN included. A star inside an expression, such as
// COUNT(*), is not judged. Findings anchor nothing (the V1 advisor placed
// them on the statement's first line). SelectAll is the only code. Ported
// from Bytebase's V1 rule STATEMENT_SELECT_NO_SELECT_ALL.
func StatementSelectNoSelectAll(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		ast.Inspect(stmt, func(n ast.Node) bool {
			if sel, ok := n.(*ast.SelectStmt); ok && selectListHasStar(sel.TargetList) {
				findings = append(findings, v1review.Finding{Code: v1review.SelectAll})
			}
			return true
		})
	}
	return findings
}

// selectListHasStar reports whether a select list holds * (a StarExpr) or
// table.* (a ColumnRef with Star).
func selectListHasStar(targets []ast.ExprNode) bool {
	for _, expr := range targets {
		switch e := expr.(type) {
		case *ast.StarExpr:
			return true
		case *ast.ColumnRef:
			if e.Star {
				return true
			}
		}
	}
	return false
}
