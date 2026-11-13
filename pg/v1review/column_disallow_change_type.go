package v1review

import (
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnDisallowChangeType judges top-level ALTER statements, whatever
// relation they name, for an ALTER COLUMN ... TYPE action, reporting one
// ColumnTypeChanged per statement however many such actions it holds. Name
// is the first such action's column and Table the relation's name without
// its schema, "" when the statement names none. A statement in a routine body
// is not judged. The finding anchors the whole statement: Bytebase's V1
// advisor printed the statement's text and reported the line it starts on.
// Ported from Bytebase's V1 rule COLUMN_DISALLOW_CHANGE_TYPE.
func ColumnDisallowChangeType(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		n, ok := statementNode(stmt).(*ast.AlterTableStmt)
		if !ok {
			continue
		}
		for _, cmd := range alterCommands(n) {
			if ast.AlterTableType(cmd.Subtype) != ast.AT_AlterColumnType {
				continue
			}
			start, end := span(n.Loc)
			findings = append(findings, v1review.Finding{
				Code:  v1review.ColumnTypeChanged,
				Name:  cmd.Name,
				Table: relname(n.Relation),
				Start: start,
				End:   end,
			})
			break
		}
	}
	return findings
}
