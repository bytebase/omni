package v1review

import (
	"github.com/bytebase/omni/googlesql/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementWhereRequireUpdateDelete reports each UPDATE or DELETE with no
// WHERE clause anywhere in the statements, nested ones included (script
// blocks such as BEGIN ... END, IF and loops). The UPDATE and DELETE
// actions of a MERGE are other nodes and are not judged. Name is the
// statement's kind, "UPDATE" or "DELETE". Each finding anchors the UPDATE or
// DELETE. WhereMissing is the only code. Ported from Bytebase's V1 rule
// STATEMENT_WHERE_REQUIRE_UPDATE_DELETE.
func StatementWhereRequireUpdateDelete(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		ast.Inspect(stmt, func(n ast.Node) bool {
			switch s := n.(type) {
			case *ast.UpdateStmt:
				if s.Where == nil {
					start, end := anchor(s.Loc)
					findings = append(findings, v1review.Finding{Code: v1review.WhereMissing, Name: "UPDATE", Start: start, End: end})
				}
			case *ast.DeleteStmt:
				if s.Where == nil {
					start, end := anchor(s.Loc)
					findings = append(findings, v1review.Finding{Code: v1review.WhereMissing, Name: "DELETE", Start: start, End: end})
				}
			}
			return true
		})
	}
	return findings
}
