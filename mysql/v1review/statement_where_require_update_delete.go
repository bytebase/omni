package v1review

import (
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementWhereRequireUpdateDelete reports each top-level UPDATE or DELETE
// with no WHERE clause; statements nested in a routine body are not judged.
// Name is the statement's kind, "UPDATE" or "DELETE". Findings anchor
// nothing (the V1 advisor placed them on the statement's first line).
// WhereMissing is the only code. Ported from Bytebase's V1 rule
// STATEMENT_WHERE_REQUIRE_UPDATE_DELETE.
func StatementWhereRequireUpdateDelete(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.DeleteStmt:
			if n.Where == nil {
				findings = append(findings, v1review.Finding{Code: v1review.WhereMissing, Name: "DELETE"})
			}
		case *ast.UpdateStmt:
			if n.Where == nil {
				findings = append(findings, v1review.Finding{Code: v1review.WhereMissing, Name: "UPDATE"})
			}
		}
	}
	return findings
}
