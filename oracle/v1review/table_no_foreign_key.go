package v1review

import (
	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/v1review"
)

// TableNoForeignKey judges, after reading every statement, which tables
// top-level statements gave a foreign key: a CREATE TABLE with a FOREIGN
// KEY table constraint or column REFERENCES clause, and an ALTER TABLE
// command carrying a FOREIGN KEY constraint, whatever its action. A column
// an ALTER TABLE adds with a REFERENCES clause is not judged, nor is a
// statement nested in CREATE SCHEMA or a PL/SQL block. Each such table is
// one TableForeignKey finding, in the order the tables were first given a
// foreign key; V1 reported them in Go map order. Name is the table's
// qualified name (see qualifiedTable). A finding anchors the last foreign
// key the table was given: the CREATE TABLE statement, or the constraint
// clause of the ALTER TABLE; Statement is that statement's index in stmts.
// A statement may be passed bare or wrapped in an *ast.RawStmt. Ported from
// Bytebase's V1 rule TABLE_NO_FOREIGN_KEY.
func TableNoForeignKey(stmts []ast.Node, currentSchema string) []v1review.Finding {
	var names []string
	anchors := make(map[string]tableAnchor)
	mark := func(name string, a tableAnchor) {
		if _, ok := anchors[name]; !ok {
			names = append(names, name)
		}
		anchors[name] = a
	}
	for i, stmt := range stmts {
		if raw, ok := stmt.(*ast.RawStmt); ok {
			stmt = raw.Stmt
		}
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if createTableHasConstraint(n, ast.CONSTRAINT_FOREIGN) {
				mark(qualifiedTable(n.Name, currentSchema), tableAnchor{i, n.Loc})
			}
		case *ast.AlterTableStmt:
			if n.Actions == nil {
				continue
			}
			for _, item := range n.Actions.Items {
				if cmd, ok := item.(*ast.AlterTableCmd); ok && cmd.Constraint != nil && cmd.Constraint.Type == ast.CONSTRAINT_FOREIGN {
					mark(qualifiedTable(n.Name, currentSchema), tableAnchor{i, cmd.Constraint.Loc})
				}
			}
		}
	}
	var findings []v1review.Finding
	for _, name := range names {
		findings = append(findings, anchors[name].finding(v1review.TableForeignKey, name))
	}
	return findings
}

// createTableHasConstraint reports whether a CREATE TABLE has a constraint
// of that type, on a column or on the table.
func createTableHasConstraint(n *ast.CreateTableStmt, typ ast.ConstraintType) bool {
	if n.Columns != nil {
		for _, item := range n.Columns.Items {
			col, ok := item.(*ast.ColumnDef)
			if !ok || col.Constraints == nil {
				continue
			}
			for _, c := range col.Constraints.Items {
				if c, ok := c.(*ast.ColumnConstraint); ok && c.Type == typ {
					return true
				}
			}
		}
	}
	for _, c := range tableConstraints(n.Constraints) {
		if c.Type == typ {
			return true
		}
	}
	return false
}
