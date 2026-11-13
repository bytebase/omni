package v1review

import (
	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementWhereRequireUpdateDelete reports, in walk order, every UPDATE and
// DELETE without a WHERE clause anywhere in the statements, EXPLAIN and
// routine bodies included. Findings anchor nothing (the V1 advisor placed
// them on the statement's first line). WhereMissing is the only code; Name
// is not set, as the V1 advisor did not tell the kinds apart. Ported from
// Bytebase's V1 rule STATEMENT_WHERE_REQUIRE_UPDATE_DELETE.
func StatementWhereRequireUpdateDelete(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		ast.Inspect(stmt, func(n ast.Node) bool {
			switch s := n.(type) {
			case *ast.UpdateStmt:
				if s.Where == nil {
					findings = append(findings, v1review.Finding{Code: v1review.WhereMissing})
				}
			case *ast.DeleteStmt:
				if s.Where == nil {
					findings = append(findings, v1review.Finding{Code: v1review.WhereMissing})
				}
			}
			return true
		})
	}
	return findings
}
