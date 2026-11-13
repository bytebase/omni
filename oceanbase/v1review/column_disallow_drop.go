package v1review

import (
	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnDisallowDrop judges top-level ALTER TABLE statements for DROP COLUMN
// actions, one finding per action. An ALTER TABLE with no table is not
// judged, though one whose table has an empty name is. Statements nested in a
// routine body are not judged.
//
// Each finding names the dropped column as written and its table, without
// the database qualifier, and anchors the DROP COLUMN action. Findings are in
// statement order, then action order. Ported from Bytebase's V1 rule
// COLUMN_DISALLOW_DROP, which only OceanBase registers.
func ColumnDisallowDrop(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		n, ok := stmt.(*ast.AlterTableStmt)
		if !ok || n.Table == nil {
			continue
		}
		for _, cmd := range n.Commands {
			if cmd == nil || cmd.Type != ast.ATDropColumn {
				continue
			}
			start, end := span(cmd.Loc)
			findings = append(findings, v1review.Finding{
				Code:  v1review.ColumnDropped,
				Name:  cmd.Name,
				Table: n.Table.Name,
				Start: start,
				End:   end,
			})
		}
	}
	return findings
}
