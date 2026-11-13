package v1review

import (
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementDisallowLimit reports LIMIT clauses of top-level DML: one for an
// UPDATE or DELETE with a LIMIT, and for an INSERT ... SELECT one for each
// LIMIT on its query, on the outer clauses of a parenthesized query and its
// inner query, and on a set operation and each of its sides, recursively;
// subqueries and CTEs are not judged. Name is the statement's kind,
// "INSERT", "UPDATE" or "DELETE". Findings anchor nothing (the V1 advisor
// placed them on the statement's first line). DMLLimit is the only code.
// Ported from Bytebase's V1 rule STATEMENT_DISALLOW_LIMIT.
func StatementDisallowLimit(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.DeleteStmt:
			if n.Limit != nil {
				findings = append(findings, v1review.Finding{Code: v1review.DMLLimit, Name: "DELETE"})
			}
		case *ast.UpdateStmt:
			if n.Limit != nil {
				findings = append(findings, v1review.Finding{Code: v1review.DMLLimit, Name: "UPDATE"})
			}
		case *ast.InsertStmt:
			findings = insertSelectLimits(n.Select, findings)
		}
	}
	return findings
}

func insertSelectLimits(sel *ast.SelectStmt, findings []v1review.Finding) []v1review.Finding {
	if sel == nil {
		return findings
	}
	if sel.Limit != nil {
		findings = append(findings, v1review.Finding{Code: v1review.DMLLimit, Name: "INSERT"})
	}
	if sel.ParenSource != nil {
		return insertSelectLimits(sel.ParenSource, findings)
	}
	if sel.SetOp != ast.SetOpNone {
		findings = insertSelectLimits(sel.Left, findings)
		return insertSelectLimits(sel.Right, findings)
	}
	return findings
}
