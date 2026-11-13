package v1review

import (
	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementDisallowLimit reports each top-level UPDATE or DELETE with a
// LIMIT, and each top-level INSERT ... SELECT whose outermost query has one;
// a LIMIT in a subquery or a parenthesized set-operation arm is not judged,
// nor is a statement nested in a routine body. Findings anchor nothing (the
// V1 advisor placed them on the statement's first line). DMLLimit is the
// only code; Name is "INSERT", "UPDATE" or "DELETE". Ported from Bytebase's
// V1 rule STATEMENT_DISALLOW_LIMIT.
func StatementDisallowLimit(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		var kind string
		switch n := stmt.(type) {
		case *ast.UpdateStmt:
			if n.Limit != nil {
				kind = "UPDATE"
			}
		case *ast.DeleteStmt:
			if n.Limit != nil {
				kind = "DELETE"
			}
		case *ast.InsertStmt:
			if n.Select != nil && n.Select.Limit != nil {
				kind = "INSERT"
			}
		}
		if kind != "" {
			findings = append(findings, v1review.Finding{Code: v1review.DMLLimit, Name: kind})
		}
	}
	return findings
}
