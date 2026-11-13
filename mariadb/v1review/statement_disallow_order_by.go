package v1review

import (
	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementDisallowOrderBy reports each top-level UPDATE or DELETE that has
// an ORDER BY clause; statements in a routine body are not judged. Each
// finding anchors the statement. DMLOrderBy is the only code; Name is
// "UPDATE" or "DELETE". Ported from Bytebase's V1 rule
// STATEMENT_DISALLOW_ORDER_BY.
func StatementDisallowOrderBy(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.DeleteStmt:
			if len(n.OrderBy) > 0 {
				start, end := statementAnchor(n.Loc)
				findings = append(findings, v1review.Finding{Code: v1review.DMLOrderBy, Name: "DELETE", Start: start, End: end})
			}
		case *ast.UpdateStmt:
			if len(n.OrderBy) > 0 {
				start, end := statementAnchor(n.Loc)
				findings = append(findings, v1review.Finding{Code: v1review.DMLOrderBy, Name: "UPDATE", Start: start, End: end})
			}
		}
	}
	return findings
}
