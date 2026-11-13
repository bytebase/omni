package v1review

import (
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// SchemaBackwardCompatibility judges which top-level statements may break
// the existing data or code, and reports at most one finding a statement,
// with Statement set: it takes every statement of the change in one call, since
// a table the change itself just created may be altered freely. A
// statement under EXPLAIN or in a routine body is not judged.
//
// Reported are every DROP DATABASE; a DROP TABLE or DROP MATERIALIZED VIEW;
// a RENAME of a column, unless on the last created table; a RENAME of any
// other object, a schema or database included; a CREATE UNIQUE INDEX; and
// an ALTER TABLE, by the first of its actions that drops a column, alters a
// column's type, or adds a primary key, unique, foreign key or check
// constraint, a check added NOT VALID excepted. An ALTER of a foreign
// table, view, index or sequence is judged as an ALTER TABLE. A DROP VIEW,
// DROP INDEX and DROP of any other object is not reported.
//
// An ALTER TABLE, CREATE UNIQUE INDEX or column RENAME is not reported when
// its table has the name of the last table an earlier CREATE TABLE created.
// Only that one table is remembered, by the name the parser stores without
// its schema, so a table created two statements back is judged as
// existing, and a table of the same name in another schema is taken for
// the new one.
//
// Every finding anchors the whole statement. Ported from Bytebase's V1 rule
// SCHEMA_BACKWARD_COMPATIBILITY.
func SchemaBackwardCompatibility(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	lastCreated := ""
	for i, stmt := range stmts {
		if raw, ok := stmt.(*ast.RawStmt); ok {
			stmt = raw.Stmt
		}
		var code v1review.Code
		var loc ast.Loc
		switch n := stmt.(type) {
		case *ast.CreateStmt:
			if n.Relation != nil {
				lastCreated = n.Relation.Relname
			}
		case *ast.DropdbStmt:
			code, loc = v1review.IncompatibleDropDatabase, n.Loc
		case *ast.DropStmt:
			if removeType := ast.ObjectType(n.RemoveType); removeType == ast.OBJECT_TABLE || removeType == ast.OBJECT_MATVIEW {
				code, loc = v1review.IncompatibleDropTable, n.Loc
			}
		case *ast.RenameStmt:
			if n.RenameType != ast.OBJECT_COLUMN {
				code, loc = v1review.IncompatibleRenameTable, n.Loc
			} else if n.Relation != nil && n.Relation.Relname != lastCreated {
				code, loc = v1review.IncompatibleRenameColumn, n.Loc
			}
		case *ast.AlterTableStmt:
			if n.Relation != nil && n.Relation.Relname != lastCreated {
				code, loc = alterTableIncompatibility(n), n.Loc
			}
		case *ast.IndexStmt:
			if n.Unique && n.Relation != nil && n.Relation.Relname != lastCreated {
				code, loc = v1review.IncompatibleAddUniqueKey, n.Loc
			}
		}
		if code != "" {
			start, end := span(loc)
			findings = append(findings, v1review.Finding{Code: code, Statement: i, Start: start, End: end})
		}
	}
	return findings
}

// alterTableIncompatibility returns the code of the first action of n that
// may break the existing data or code, "" when none does.
func alterTableIncompatibility(n *ast.AlterTableStmt) v1review.Code {
	if n.Cmds == nil {
		return ""
	}
	for _, item := range n.Cmds.Items {
		cmd, ok := item.(*ast.AlterTableCmd)
		if !ok {
			continue
		}
		switch ast.AlterTableType(cmd.Subtype) {
		case ast.AT_DropColumn:
			return v1review.IncompatibleDropColumn
		case ast.AT_AlterColumnType:
			return v1review.IncompatibleAlterColumn
		case ast.AT_AddConstraint:
			constraint, ok := cmd.Def.(*ast.Constraint)
			if !ok {
				continue
			}
			switch constraint.Contype {
			case ast.CONSTR_PRIMARY:
				return v1review.IncompatibleAddPrimaryKey
			case ast.CONSTR_UNIQUE:
				return v1review.IncompatibleAddUniqueKey
			case ast.CONSTR_FOREIGN:
				return v1review.IncompatibleAddForeignKey
			case ast.CONSTR_CHECK:
				if !constraint.SkipValidation {
					return v1review.IncompatibleAddCheck
				}
			}
		}
	}
	return ""
}
