package v1review

import (
	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementWhereRequireSelect reports each SELECT with a FROM clause and no
// WHERE clause in a top-level SELECT statement: the statement itself, its
// CTEs, a parenthesized SELECT, both sides of a set operation, and SELECTs
// nested in the WHERE clause, FROM clause or select list of a judged SELECT.
// A SELECT with no FROM (SELECT 1) and SELECTs in other statement kinds are
// not judged. CTEs are judged before the SELECT that owns them; each finding
// anchors its SELECT. WhereMissing is the only code. Ported from Bytebase's
// V1 rule STATEMENT_WHERE_REQUIRE_SELECT.
func StatementWhereRequireSelect(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		if sel, ok := stmt.(*ast.SelectStmt); ok {
			findings = appendSelectsWithoutWhere(findings, sel)
		}
	}
	return findings
}

func appendSelectsWithoutWhere(findings []v1review.Finding, sel *ast.SelectStmt) []v1review.Finding {
	if sel == nil {
		return findings
	}
	for _, cte := range sel.CTEs {
		if cte.Select != nil {
			findings = appendSelectsWithoutWhere(findings, cte.Select)
		}
	}
	if sel.ParenSource != nil {
		return appendSelectsWithoutWhere(findings, sel.ParenSource)
	}
	if sel.SetOp != ast.SetOpNone {
		findings = appendSelectsWithoutWhere(findings, sel.Left)
		return appendSelectsWithoutWhere(findings, sel.Right)
	}
	if len(sel.From) == 0 {
		return findings
	}
	if sel.Where == nil {
		start, end := anchor(sel.Loc)
		findings = append(findings, v1review.Finding{Code: v1review.WhereMissing, Start: start, End: end})
	}
	subtrees := []ast.Node{sel.Where}
	for _, from := range sel.From {
		subtrees = append(subtrees, from)
	}
	for _, target := range sel.TargetList {
		subtrees = append(subtrees, target)
	}
	for _, subtree := range subtrees {
		ast.Inspect(subtree, func(n ast.Node) bool {
			if sub, ok := n.(*ast.SelectStmt); ok {
				findings = appendSelectsWithoutWhere(findings, sub)
				return false
			}
			return true
		})
	}
	return findings
}
