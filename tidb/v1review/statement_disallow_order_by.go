package v1review

import (
	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementDisallowOrderBy reports each top-level UPDATE or DELETE with an
// ORDER BY; one nested in a routine body is not judged. Findings anchor
// nothing (the V1 advisor placed them on the statement's first line).
// DMLOrderBy is the only code; Name is "UPDATE" or "DELETE". Ported from
// Bytebase's V1 rule STATEMENT_DISALLOW_ORDER_BY.
func StatementDisallowOrderBy(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.UpdateStmt:
			if len(n.OrderBy) > 0 {
				findings = append(findings, v1review.Finding{Code: v1review.DMLOrderBy, Name: "UPDATE"})
			}
		case *ast.DeleteStmt:
			if len(n.OrderBy) > 0 {
				findings = append(findings, v1review.Finding{Code: v1review.DMLOrderBy, Name: "DELETE"})
			}
		}
	}
	return findings
}
