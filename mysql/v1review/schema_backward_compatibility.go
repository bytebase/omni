package v1review

import (
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// SchemaBackwardCompatibility judges which top-level statements may break
// the existing data or code, and reports at most one finding a statement,
// with Statement set: it takes every statement of the change in one call, since
// a table the change itself just created may be altered freely. Statements
// nested in a routine body are not judged.
//
// Reported are every DROP DATABASE, DROP TABLE and RENAME TABLE; a CREATE
// UNIQUE INDEX; and an ALTER TABLE, by the first of its actions that is a
// RENAME or DROP COLUMN, a RENAME of the table, an added primary key,
// unique key, foreign key or enforced check, an ALTER CHECK or ALTER
// CONSTRAINT enforcement change, or a MODIFY or CHANGE COLUMN.
//
// An ALTER TABLE or CREATE UNIQUE INDEX is not reported when its table has
// the name of the last table created by an earlier CREATE TABLE that
// defines columns. Only that one table is remembered, by its name as
// written and without its database, so a table created two statements back
// is judged as existing, and a table of the same name in another database
// is taken for the new one. A DROP VIEW, DROP INDEX and DROP of any other
// object is not reported.
//
// Every finding anchors the whole statement. Ported from Bytebase's V1 rule
// SCHEMA_BACKWARD_COMPATIBILITY.
func SchemaBackwardCompatibility(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	lastCreated := ""
	for i, stmt := range stmts {
		var code v1review.Code
		var loc ast.Loc
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.Table != nil && len(n.Columns) > 0 {
				lastCreated = n.Table.Name
			}
		case *ast.DropDatabaseStmt:
			code, loc = v1review.IncompatibleDropDatabase, n.Loc
		case *ast.RenameTableStmt:
			code, loc = v1review.IncompatibleRenameTable, n.Loc
		case *ast.DropTableStmt:
			code, loc = v1review.IncompatibleDropTable, n.Loc
		case *ast.AlterTableStmt:
			if n.Table != nil && n.Table.Name != lastCreated {
				code, loc = incompatibleAction(n), n.Loc
			}
		case *ast.CreateIndexStmt:
			if n.Unique && n.Table != nil && n.Table.Name != lastCreated {
				code, loc = v1review.IncompatibleAddUniqueKey, n.Loc
			}
		}
		if code == "" {
			continue
		}
		start, end := span(loc)
		findings = append(findings, v1review.Finding{Code: code, Start: start, End: end, Statement: i})
	}
	return findings
}

// incompatibleAction returns the code of the first action of an ALTER TABLE
// that SchemaBackwardCompatibility reports, "" when it has none.
func incompatibleAction(n *ast.AlterTableStmt) v1review.Code {
	for _, cmd := range n.Commands {
		if cmd == nil {
			continue
		}
		switch cmd.Type {
		case ast.ATRenameColumn:
			return v1review.IncompatibleRenameColumn
		case ast.ATDropColumn:
			return v1review.IncompatibleDropColumn
		case ast.ATRenameTable:
			return v1review.IncompatibleRenameTable
		case ast.ATAddConstraint:
			if cmd.Constraint == nil {
				continue
			}
			switch cmd.Constraint.Type {
			case ast.ConstrPrimaryKey:
				return v1review.IncompatibleAddPrimaryKey
			case ast.ConstrUnique:
				return v1review.IncompatibleAddUniqueKey
			case ast.ConstrForeignKey:
				return v1review.IncompatibleAddForeignKey
			case ast.ConstrCheck:
				if !cmd.Constraint.NotEnforced {
					return v1review.IncompatibleAddCheck
				}
			}
		case ast.ATAlterCheckEnforced:
			return v1review.IncompatibleAlterCheck
		case ast.ATModifyColumn, ast.ATChangeColumn:
			return v1review.IncompatibleAlterColumn
		}
	}
	return ""
}
