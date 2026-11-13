package v1review

import (
	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/v1review"
)

// TableRequirePK judges each table a top-level CREATE TABLE creates as the
// top-level statements leave it, reporting after reading them all;
// statements nested in a routine body or an IF or BEGIN block are not read.
// Tables are told apart as TableNoForeignKey tells them apart.
//
// A CREATE TABLE gives its table no primary key, then one if a column's
// constraint or a table constraint is a PRIMARY KEY. An ALTER TABLE gives
// its table one for each ADD of a column with a PRIMARY KEY constraint and
// each ADD CONSTRAINT ... PRIMARY KEY, whether or not the script created
// the table, so a later CREATE TABLE of it starts again with none. Nothing
// removes one: neither ALTER TABLE ... DROP CONSTRAINT nor DROP TABLE.
//
// Each created table left without a primary key is one
// TablePrimaryKeyMissing, anchoring the table name of its last CREATE
// TABLE; Name is that name as written without brackets, with the
// qualifiers it was written with, and Statement that statement's index.
//
// The V1 advisor reported tables in Go map order, which varies from run to
// run; the findings come in the order the tables were first created or
// given a primary key. Ported from Bytebase's V1 rule TABLE_REQUIRE_PK.
func TableRequirePK(stmts []ast.Node) []v1review.Finding {
	tables := &tableStates{}
	for i, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.Name == nil {
				continue
			}
			t := tables.get(joinTableRef(n.Name, "dbo"))
			t.flag, t.name = createTableHasPrimaryKey(n), tableRefText(n.Name)
			t.anchor(i, n.Name.Loc)
		case *ast.AlterTableStmt:
			if n.Name == nil || n.Actions == nil {
				continue
			}
			key := joinTableRef(n.Name, "dbo")
			for _, item := range n.Actions.Items {
				action, ok := item.(*ast.AlterTableAction)
				if ok && alterTableAddsPrimaryKey(action) {
					tables.get(key).flag = true
				}
			}
		default:
		}
	}
	var findings []v1review.Finding
	for _, t := range tables.list {
		if !t.flag {
			findings = append(findings, t.finding(v1review.TablePrimaryKeyMissing))
		}
	}
	return findings
}

func createTableHasPrimaryKey(n *ast.CreateTableStmt) bool {
	if n.Columns != nil {
		for _, item := range n.Columns.Items {
			col, ok := item.(*ast.ColumnDef)
			if ok && listHasPrimaryKey(col.Constraints) {
				return true
			}
		}
	}
	return listHasPrimaryKey(n.Constraints)
}

func alterTableAddsPrimaryKey(action *ast.AlterTableAction) bool {
	switch action.Type {
	case ast.ATAddColumn:
		return action.Column != nil && listHasPrimaryKey(action.Column.Constraints)
	case ast.ATAddConstraint:
		return action.Constraint != nil && action.Constraint.Type == ast.ConstraintPrimaryKey
	default:
		return false
	}
}

// listHasPrimaryKey reports whether a constraint list holds a PRIMARY KEY.
func listHasPrimaryKey(constraints *ast.List) bool {
	if constraints == nil {
		return false
	}
	for _, item := range constraints.Items {
		if cd, ok := item.(*ast.ConstraintDef); ok && cd.Type == ast.ConstraintPrimaryKey {
			return true
		}
	}
	return false
}
