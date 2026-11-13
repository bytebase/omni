package v1review

import (
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementDisallowAddNotNull reports each ALTER COLUMN ... SET NOT NULL of
// a top-level ALTER TABLE, one NotNullSet finding per command with Name the
// column as the parser stores it. A NOT NULL column added or created is not
// judged. Findings anchor nothing: Bytebase's V1 advisor placed them on the
// statement's first line. Ported from Bytebase's V1 rule
// STATEMENT_DISALLOW_ADD_NOT_NULL.
func StatementDisallowAddNotNull(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		alter, ok := unwrapRaw(stmt).(*ast.AlterTableStmt)
		if !ok {
			continue
		}
		for _, cmd := range alterTableCmds(alter) {
			if ast.AlterTableType(cmd.Subtype) == ast.AT_SetNotNull {
				findings = append(findings, v1review.Finding{Code: v1review.NotNullSet, Name: cmd.Name})
			}
		}
	}
	return findings
}
