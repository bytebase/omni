package v1review

import (
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// TableNoForeignKey judges the foreign keys top-level statements add, at
// most one TableForeignKey finding per statement: a CREATE TABLE with a
// FOREIGN KEY table constraint or a column REFERENCES clause in its element
// list, and an ALTER statement with an ADD CONSTRAINT ... FOREIGN KEY
// command, whatever object kind it alters. ALTER TABLE ... ADD COLUMN ...
// REFERENCES, CREATE TABLE AS, and statements under EXPLAIN or in a routine
// body are not judged. Table is the table name as parsed, empty when the AST
// has none; Value is its schema as written, "public" when unqualified.
// Each finding anchors the whole statement. A statement may be passed bare
// or wrapped in an *ast.RawStmt. Ported from Bytebase's V1 rule
// TABLE_NO_FOREIGN_KEY.
func TableNoForeignKey(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		if raw, ok := stmt.(*ast.RawStmt); ok {
			stmt = raw.Stmt
		}
		var relation *ast.RangeVar
		var loc ast.Loc
		switch n := stmt.(type) {
		case *ast.CreateStmt:
			if !createHasForeignKey(n) {
				continue
			}
			relation, loc = n.Relation, n.Loc
		case *ast.AlterTableStmt:
			if !alterTableHas(n, func(cmd *ast.AlterTableCmd) bool {
				c, ok := cmd.Def.(*ast.Constraint)
				return ast.AlterTableType(cmd.Subtype) == ast.AT_AddConstraint && ok && c != nil && c.Contype == ast.CONSTR_FOREIGN
			}) {
				continue
			}
			relation, loc = n.Relation, n.Loc
		default:
			continue
		}
		start, end := span(loc)
		findings = append(findings, v1review.Finding{Code: v1review.TableForeignKey, Table: relname(relation), Value: schemaOrDefault(relation), Start: start, End: end})
	}
	return findings
}

// createHasForeignKey reports whether a CREATE TABLE's element list holds a
// foreign key, as a table constraint or a column's constraint.
func createHasForeignKey(n *ast.CreateStmt) bool {
	if n.TableElts == nil {
		return false
	}
	for _, item := range n.TableElts.Items {
		switch e := item.(type) {
		case *ast.Constraint:
			if e.Contype == ast.CONSTR_FOREIGN {
				return true
			}
		case *ast.ColumnDef:
			for _, c := range columnConstraints(e) {
				if c.Contype == ast.CONSTR_FOREIGN {
					return true
				}
			}
		}
	}
	return false
}
