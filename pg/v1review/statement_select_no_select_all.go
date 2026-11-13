package v1review

import (
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementSelectNoSelectAll reports each SELECT whose select list holds a
// * or t.*, one SelectAll
// finding per SELECT, visiting SELECTs as selectWalk describes, so a
// subquery's finding comes before its enclosing SELECT's. A * as a function
// argument, such as count(*), is not judged. A finding on the top-level
// SELECT statement anchors nothing, and one on any other SELECT anchors that
// SELECT: Bytebase's V1 advisor printed the statement's text or that
// SELECT's, and placed every finding on the statement's first line. Ported
// from Bytebase's V1 rule STATEMENT_SELECT_NO_SELECT_ALL.
func StatementSelectNoSelectAll(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	w := &selectWalk{visit: func(sel *ast.SelectStmt, top bool) {
		if selectsAll(sel) {
			findings = append(findings, selectFinding(v1review.SelectAll, sel, top))
		}
	}}
	for _, stmt := range stmts {
		w.statement(stmt)
	}
	return findings
}

func selectsAll(sel *ast.SelectStmt) bool {
	if sel.TargetList == nil {
		return false
	}
	for _, item := range sel.TargetList.Items {
		target, ok := item.(*ast.ResTarget)
		if !ok {
			continue
		}
		ref, ok := target.Val.(*ast.ColumnRef)
		if !ok || ref.Fields == nil {
			continue
		}
		for _, field := range ref.Fields.Items {
			if _, ok := field.(*ast.A_Star); ok {
				return true
			}
		}
	}
	return false
}
