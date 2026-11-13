package v1review

import (
	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementAddColumnWithoutPosition reports each ADD COLUMN command of a
// top-level ALTER TABLE that positions its column with FIRST or AFTER;
// MODIFY or CHANGE positions and statements in a routine body are not
// judged. Findings anchor nothing (the V1 advisor took the statement's first
// content line). ColumnPositioned is the only code; Name is "FIRST" or
// "AFTER". Ported from Bytebase's V1 rule
// STATEMENT_ADD_COLUMN_WITHOUT_POSITION.
func StatementAddColumnWithoutPosition(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		alter, ok := stmt.(*ast.AlterTableStmt)
		if !ok || alter.Table == nil {
			continue
		}
		for _, cmd := range alter.Commands {
			if cmd.Type != ast.ATAddColumn {
				continue
			}
			switch {
			case cmd.First:
				findings = append(findings, v1review.Finding{Code: v1review.ColumnPositioned, Name: "FIRST"})
			case cmd.After != "":
				findings = append(findings, v1review.Finding{Code: v1review.ColumnPositioned, Name: "AFTER"})
			}
		}
	}
	return findings
}
