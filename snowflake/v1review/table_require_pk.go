package v1review

import (
	"github.com/bytebase/omni/snowflake/ast"
	"github.com/bytebase/omni/v1review"
)

// TableRequirePK judges each table as the top-level statements leave it,
// reporting after reading them all. Tables are told apart as
// TableNoForeignKey tells them apart.
//
// A plain CREATE TABLE gives its table no primary key, then one if a
// column's inline constraint or a table constraint is a PRIMARY KEY; CREATE
// TABLE ... AS SELECT, LIKE and CLONE are not read. A DROP TABLE forgets
// the table. In an ALTER TABLE, an ADD of a PRIMARY KEY constraint gives
// its table one, whether or not the script created it, and DROP PRIMARY KEY
// takes away the primary key of a table read before; DROP CONSTRAINT by
// name does not, since it does not say the constraint's kind, and a column
// added with an inline PRIMARY KEY gives none.
//
// Each table read and left without a primary key is one
// TablePrimaryKeyMissing. Name is the table as written, quotes kept, in the
// last CREATE TABLE or ALTER TABLE read that adds, drops or renames a
// constraint. The finding anchors the last DROP PRIMARY KEY read for the
// table, else its CREATE TABLE; Statement is that statement's index.
//
// The V1 advisor reported tables in Go map order, which varies from run to
// run; the findings come in the order the tables were first read, a table
// read again after a DROP TABLE counting as new. Ported from Bytebase's V1
// rule TABLE_REQUIRE_PK.
func TableRequirePK(stmts []ast.Node) []v1review.Finding {
	tables := newTableStates()
	for i, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.AsSelect != nil || n.Like != nil || n.Clone != nil {
				continue
			}
			key := tableStateKey(n.Name)
			t := tables.get(key)
			t.hasPK = false
			tables.names[key] = n.Name.String()
			t.anchor(i, n.Loc)
			for _, column := range n.Columns {
				if column.InlineConstraint != nil && column.InlineConstraint.Type == ast.ConstrPrimaryKey {
					t.hasPK = true
				}
			}
			for _, constraint := range n.Constraints {
				if constraint.Type == ast.ConstrPrimaryKey {
					t.hasPK = true
				}
			}
		case *ast.DropStmt:
			if n.Kind == ast.DropTable {
				tables.remove(tableStateKey(n.Name))
			}
		case *ast.AlterTableStmt:
			if !alterTableHasConstraintAction(n) {
				continue
			}
			key := tableStateKey(n.Name)
			tables.names[key] = n.Name.String()
			for _, action := range n.Actions {
				switch action.Kind {
				case ast.AlterTableAddConstraint:
					if action.Constraint != nil && action.Constraint.Type == ast.ConstrPrimaryKey {
						tables.get(key).hasPK = true
					}
				case ast.AlterTableDropConstraint:
					if t, ok := tables.byKey[key]; ok && action.IsPrimaryKey {
						t.hasPK = false
						t.anchor(i, action.Loc)
					}
				default:
				}
			}
		default:
		}
	}
	var findings []v1review.Finding
	for _, t := range tables.list {
		if !t.hasPK {
			findings = append(findings, tables.finding(t, v1review.TablePrimaryKeyMissing))
		}
	}
	return findings
}
