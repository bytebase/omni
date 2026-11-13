package v1review

import (
	"github.com/bytebase/omni/snowflake/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementWhereRequireUpdateDelete reports every UPDATE and DELETE without
// a WHERE clause anywhere in the statements, in walk order, statements
// nested in scripting blocks included. MERGE's WHEN ... THEN UPDATE and
// DELETE actions are not UPDATE or DELETE statements and are not judged. A
// finding anchors the statement. WhereMissing is the only code; Name is
// "UPDATE" or "DELETE". Ported from Bytebase's V1 rule
// STATEMENT_WHERE_REQUIRE_UPDATE_DELETE.
func StatementWhereRequireUpdateDelete(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		ast.Inspect(stmt, func(n ast.Node) bool {
			switch s := n.(type) {
			case *ast.UpdateStmt:
				if s.Where == nil {
					start, end := stmtAnchor(s.Loc)
					findings = append(findings, v1review.Finding{Code: v1review.WhereMissing, Name: "UPDATE", Start: start, End: end})
				}
			case *ast.DeleteStmt:
				if s.Where == nil {
					start, end := stmtAnchor(s.Loc)
					findings = append(findings, v1review.Finding{Code: v1review.WhereMissing, Name: "DELETE", Start: start, End: end})
				}
			}
			return true
		})
	}
	return findings
}
