package v1review

import (
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementDisallowOrderBy reports each top-level UPDATE or DELETE with an
// ORDER BY clause. Name is the statement's kind, "UPDATE" or "DELETE". Each
// finding anchors the statement. DMLOrderBy is the only code. Ported from
// Bytebase's V1 rule STATEMENT_DISALLOW_ORDER_BY.
func StatementDisallowOrderBy(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.DeleteStmt:
			if len(n.OrderBy) > 0 {
				start, end := anchor(n.Loc)
				findings = append(findings, v1review.Finding{Code: v1review.DMLOrderBy, Name: "DELETE", Start: start, End: end})
			}
		case *ast.UpdateStmt:
			if len(n.OrderBy) > 0 {
				start, end := anchor(n.Loc)
				findings = append(findings, v1review.Finding{Code: v1review.DMLOrderBy, Name: "UPDATE", Start: start, End: end})
			}
		}
	}
	return findings
}
