package v1review

import (
	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementDisallowLimit reports each LIMIT in a top-level UPDATE or DELETE,
// and each LIMIT in the SELECT a top-level INSERT ... SELECT inserts from:
// the outer LIMIT of a parenthesized SELECT or set operation, then the LIMITs
// inside it (left side before right). Subqueries elsewhere and statements in a
// routine body are not judged. Findings anchor nothing (the V1 advisor took
// the statement's first content line). DMLLimit is the only code; Name is
// "UPDATE", "DELETE" or "INSERT". Ported from Bytebase's V1 rule
// STATEMENT_DISALLOW_LIMIT.
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
			findings = appendInsertLimits(findings, n.Select)
		}
	}
	return findings
}

func appendInsertLimits(findings []v1review.Finding, sel *ast.SelectStmt) []v1review.Finding {
	if sel == nil {
		return findings
	}
	if sel.Limit != nil {
		findings = append(findings, v1review.Finding{Code: v1review.DMLLimit, Name: "INSERT"})
	}
	if sel.ParenSource != nil {
		return appendInsertLimits(findings, sel.ParenSource)
	}
	if sel.SetOp != ast.SetOpNone {
		findings = appendInsertLimits(findings, sel.Left)
		findings = appendInsertLimits(findings, sel.Right)
	}
	return findings
}
