package v1review

import (
	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnDisallowChange reports each top-level ALTER TABLE statement holding
// a CHANGE COLUMN action, once however many it holds, anchored at the whole
// statement. Ported from Bytebase's V1 rule COLUMN_DISALLOW_CHANGE.
func ColumnDisallowChange(stmts []ast.Node) []v1review.Finding {
	return alterRule(stmts, v1review.ColumnChanged, func(cmd *ast.AlterTableCmd) bool {
		return cmd.Type == ast.ATChangeColumn
	})
}

// alterRule reports, once per top-level ALTER TABLE statement and anchored at
// the whole statement, the statements holding an action bad holds for. The
// statement need not name a table.
func alterRule(stmts []ast.Node, code v1review.Code, bad func(*ast.AlterTableCmd) bool) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		alter, ok := stmt.(*ast.AlterTableStmt)
		if !ok {
			continue
		}
		for _, cmd := range alter.Commands {
			if cmd != nil && bad(cmd) {
				start, end := span(alter.Loc)
				findings = append(findings, v1review.Finding{Code: code, Start: start, End: end})
				break
			}
		}
	}
	return findings
}
