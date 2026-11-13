package v1review

import (
	"slices"

	"github.com/bytebase/omni/snowflake/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnNoNull judges each top-level statement on its own for nullable
// columns. A plain CREATE TABLE (not AS SELECT, LIKE or CLONE) records each
// column it defines, in order, then forgets one that says NOT NULL or
// carries an inline PRIMARY KEY, so a later definition of the same name
// wins; then it forgets every column a PRIMARY KEY table constraint names.
// An ALTER TABLE walks its actions in order: ADD COLUMN records and forgets
// columns the same way, ALTER COLUMN ... DROP NOT NULL records the column,
// and ADD of a PRIMARY KEY constraint forgets its columns. Columns are keyed
// by name, an unquoted name upper-cased and a quoted one as written.
// Statements nested in a Snowflake Scripting block or under EXPLAIN are not
// judged.
//
// Each finding is one recorded column, with code ColumnNullable, Name its
// key, and the anchor of its last recording: the column definition, or for
// DROP NOT NULL the column name, of which only the start is meaningful. V1
// reported a statement's columns in Go map order, which is random; they come
// here in key order, compared bytewise. Ported from Bytebase's V1 rule
// COLUMN_NO_NULL.
func ColumnNoNull(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		nullable := make(map[string]ast.Loc)
		addColumn := func(column *ast.ColumnDef) {
			name := column.Name.Normalize()
			nullable[name] = column.Loc
			if column.NotNull || (column.InlineConstraint != nil && column.InlineConstraint.Type == ast.ConstrPrimaryKey) {
				delete(nullable, name)
			}
		}
		forgetPrimaryKey := func(constraint *ast.TableConstraint) {
			if constraint.Type != ast.ConstrPrimaryKey {
				return
			}
			for _, column := range constraint.Columns {
				delete(nullable, column.Normalize())
			}
		}
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.AsSelect != nil || n.Like != nil || n.Clone != nil {
				continue
			}
			for _, column := range n.Columns {
				addColumn(column)
			}
			for _, constraint := range n.Constraints {
				forgetPrimaryKey(constraint)
			}
		case *ast.AlterTableStmt:
			for _, action := range n.Actions {
				switch action.Kind {
				case ast.AlterTableAddColumn:
					for _, column := range action.Columns {
						addColumn(column)
					}
				case ast.AlterTableAlterColumn:
					for _, alter := range action.ColumnAlters {
						if alter.Kind == ast.ColumnAlterDropNotNull {
							nullable[alter.Column.Normalize()] = alter.Column.Loc
						}
					}
				case ast.AlterTableAddConstraint:
					if action.Constraint != nil {
						forgetPrimaryKey(action.Constraint)
					}
				default:
				}
			}
		default:
			continue
		}
		names := make([]string, 0, len(nullable))
		for name := range nullable {
			names = append(names, name)
		}
		slices.Sort(names)
		for _, name := range names {
			start, end := columnAnchor(nullable[name])
			findings = append(findings, v1review.Finding{Code: v1review.ColumnNullable, Name: name, Start: start, End: end})
		}
	}
	return findings
}
