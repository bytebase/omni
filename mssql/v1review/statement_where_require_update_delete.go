package v1review

import (
	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementWhereRequireUpdateDelete reports each top-level UPDATE and DELETE
// without a WHERE clause; one nested in a procedure body, a control-flow
// block or a MERGE is not judged. A FROM clause does not count as a WHERE
// clause. A finding anchors the statement. WhereMissing is the only code;
// Name is "UPDATE" or "DELETE". Ported from Bytebase's V1 rule
// STATEMENT_WHERE_REQUIRE_UPDATE_DELETE.
func StatementWhereRequireUpdateDelete(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.DeleteStmt:
			if n.WhereClause == nil {
				start, end := stmtAnchor(n.Loc)
				findings = append(findings, v1review.Finding{Code: v1review.WhereMissing, Name: "DELETE", Start: start, End: end})
			}
		case *ast.UpdateStmt:
			if n.WhereClause == nil {
				start, end := stmtAnchor(n.Loc)
				findings = append(findings, v1review.Finding{Code: v1review.WhereMissing, Name: "UPDATE", Start: start, End: end})
			}
		}
	}
	return findings
}
