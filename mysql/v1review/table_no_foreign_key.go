package v1review

import (
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// TableNoForeignKey judges every top-level CREATE TABLE and ALTER TABLE
// that names its table; statements nested in a routine body are not
// judged. Each table-level FOREIGN KEY of a CREATE TABLE, and each ALTER
// TABLE ... ADD [CONSTRAINT] FOREIGN KEY action, is one TableForeignKey
// anchoring that constraint, with Name the table name without its
// database. A column-level REFERENCES clause is not judged. Ported from
// Bytebase's V1 rule TABLE_NO_FOREIGN_KEY.
//
// A nil CREATE TABLE constraint panics, as the V1 advisor did; the parser
// builds none.
func TableNoForeignKey(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.Table == nil {
				continue
			}
			for _, constraint := range n.Constraints {
				if constraint.Type == ast.ConstrForeignKey {
					start, end := span(constraint.Loc)
					findings = append(findings, v1review.Finding{Code: v1review.TableForeignKey, Name: n.Table.Name, Start: start, End: end})
				}
			}
		case *ast.AlterTableStmt:
			if n.Table == nil {
				continue
			}
			for _, cmd := range n.Commands {
				if cmd != nil && cmd.Type == ast.ATAddConstraint && cmd.Constraint != nil && cmd.Constraint.Type == ast.ConstrForeignKey {
					start, end := span(cmd.Constraint.Loc)
					findings = append(findings, v1review.Finding{Code: v1review.TableForeignKey, Name: n.Table.Name, Start: start, End: end})
				}
			}
		default:
		}
	}
	return findings
}
