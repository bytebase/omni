package v1review

import (
	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// TableNoForeignKey judges the foreign keys top-level statements add: each
// table-level FOREIGN KEY constraint of a CREATE TABLE, and each ALTER TABLE
// ADD [CONSTRAINT] FOREIGN KEY, is TableForeignKey. A column's own
// REFERENCES clause, a statement whose table the AST does not name, and
// statements nested in a routine body or carried by EXPLAIN are not judged.
// Table is the table name as parsed, without its database. A CREATE TABLE
// finding anchors its constraint; an ALTER TABLE finding anchors the whole
// statement. Findings are in statement order. Ported from Bytebase's V1 rule
// TABLE_NO_FOREIGN_KEY.
func TableNoForeignKey(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.Table == nil {
				continue
			}
			for _, constraint := range n.Constraints {
				if constraint != nil && constraint.Type == ast.ConstrForeignKey {
					start, end := span(constraint.Loc)
					findings = append(findings, v1review.Finding{Code: v1review.TableForeignKey, Table: n.Table.Name, Start: start, End: end})
				}
			}
		case *ast.AlterTableStmt:
			if n.Table == nil {
				continue
			}
			start, end := span(n.Loc)
			for _, cmd := range n.Commands {
				if cmd != nil && cmd.Type == ast.ATAddConstraint && cmd.Constraint != nil && cmd.Constraint.Type == ast.ConstrForeignKey {
					findings = append(findings, v1review.Finding{Code: v1review.TableForeignKey, Table: n.Table.Name, Start: start, End: end})
				}
			}
		}
	}
	return findings
}
