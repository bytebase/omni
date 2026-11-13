package v1review

import (
	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnDisallowChangingOrder judges top-level ALTER TABLE statements for
// MODIFY COLUMN and CHANGE COLUMN actions that give a position, FIRST or
// AFTER a column, one finding per action, whatever the table. ADD COLUMN
// with a position is not judged, and neither are statements nested in a
// routine body.
//
// Each finding names the column the action modifies and anchors the whole
// statement, since V1 printed the statement and placed the advice at its
// first line. Findings are in statement order, then action order. Ported
// from Bytebase's V1 rule COLUMN_DISALLOW_CHANGING_ORDER.
func ColumnDisallowChangingOrder(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		n, ok := stmt.(*ast.AlterTableStmt)
		if !ok {
			continue
		}
		for _, cmd := range n.Commands {
			if cmd == nil || (cmd.Type != ast.ATModifyColumn && cmd.Type != ast.ATChangeColumn) {
				continue
			}
			if cmd.First || cmd.After != "" {
				findings = append(findings, statementFinding(v1review.ColumnOrderChanged, cmd.Name, n.Loc))
			}
		}
	}
	return findings
}
