package v1review

import (
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementWhereRequireSelect reports each SELECT with a FROM clause and no
// WHERE in a top-level SELECT statement: the statement's query, its CTEs,
// the inner query of a parenthesized query (whose outer clauses are then
// skipped), each side of a set operation, and subqueries in a judged
// SELECT's WHERE, FROM and select list, recursively. A SELECT without FROM
// is not judged, nor are SELECTs inside other statements (INSERT ...
// SELECT, UPDATE, CREATE ... AS SELECT). Findings come in the V1 advisor's
// order: CTEs, then the query, then its subqueries. Name is "SELECT". Each
// finding anchors the SELECT. WhereMissing is the only code. Ported from
// Bytebase's V1 rule STATEMENT_WHERE_REQUIRE_SELECT.
func StatementWhereRequireSelect(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		if sel, ok := stmt.(*ast.SelectStmt); ok {
			findings = selectWhereMissing(sel, findings)
		}
	}
	return findings
}

func selectWhereMissing(sel *ast.SelectStmt, findings []v1review.Finding) []v1review.Finding {
	if sel == nil {
		return findings
	}
	for _, cte := range sel.CTEs {
		if cte.Select != nil {
			findings = selectWhereMissing(cte.Select, findings)
		}
	}
	if sel.ParenSource != nil {
		return selectWhereMissing(sel.ParenSource, findings)
	}
	if sel.SetOp != ast.SetOpNone {
		findings = selectWhereMissing(sel.Left, findings)
		return selectWhereMissing(sel.Right, findings)
	}
	if len(sel.From) == 0 {
		return findings
	}
	if sel.Where == nil {
		start, end := anchor(sel.Loc)
		findings = append(findings, v1review.Finding{Code: v1review.WhereMissing, Name: "SELECT", Start: start, End: end})
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
				findings = selectWhereMissing(sub, findings)
				return false
			}
			return true
		})
	}
	return findings
}
