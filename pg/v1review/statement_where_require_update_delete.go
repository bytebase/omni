package v1review

import (
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementWhereRequireUpdateDelete reports each top-level UPDATE or DELETE
// without a WHERE clause, one WhereMissing finding per statement with Name
// "UPDATE" or "DELETE". WHERE CURRENT OF counts as a WHERE clause; an
// UPDATE or DELETE in a WITH query or a routine body is not judged.
// Findings anchor nothing: Bytebase's V1 advisor placed them on the
// statement's first line. Ported from Bytebase's V1 rule
// STATEMENT_WHERE_REQUIRE_UPDATE_DELETE.
func StatementWhereRequireUpdateDelete(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		switch n := unwrapRaw(stmt).(type) {
		case *ast.UpdateStmt:
			if n.WhereClause == nil {
				findings = append(findings, v1review.Finding{Code: v1review.WhereMissing, Name: "UPDATE"})
			}
		case *ast.DeleteStmt:
			if n.WhereClause == nil {
				findings = append(findings, v1review.Finding{Code: v1review.WhereMissing, Name: "DELETE"})
			}
		}
	}
	return findings
}
