package v1review

import (
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnDisallowChange judges top-level ALTER TABLE statements for CHANGE
// COLUMN actions, one finding per action, whatever the table. Statements
// nested in a routine body are not judged.
//
// Each finding names the column the action changes and anchors the whole
// statement, since V1 printed the statement and placed the advice at its
// first line. Findings are in statement order, then action order. Ported
// from Bytebase's V1 rule COLUMN_DISALLOW_CHANGE.
func ColumnDisallowChange(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		n, ok := stmt.(*ast.AlterTableStmt)
		if !ok {
			continue
		}
		for _, cmd := range n.Commands {
			if cmd != nil && cmd.Type == ast.ATChangeColumn {
				findings = append(findings, wholeStatementFinding(v1review.ColumnChanged, cmd.Name, n.Loc))
			}
		}
	}
	return findings
}
