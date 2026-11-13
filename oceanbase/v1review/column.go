package v1review

import (
	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/v1review"
)

// alterColumns returns the column definitions an ALTER TABLE action carries,
// as the V1 column rules read them: ADD COLUMN's single column, else every
// entry of its parenthesized form, nil entries included; MODIFY and CHANGE
// COLUMN's column; nothing for a nil action or any other action. Unlike
// commandColumns it keeps nil entries, so a rule that dereferences each
// column fails on one exactly as V1 did.
func alterColumns(cmd *ast.AlterTableCmd) []*ast.ColumnDef {
	if cmd == nil {
		return nil
	}
	switch cmd.Type {
	case ast.ATAddColumn:
		if cmd.Column != nil {
			return []*ast.ColumnDef{cmd.Column}
		}
		return cmd.Columns
	case ast.ATModifyColumn, ast.ATChangeColumn:
		if cmd.Column != nil {
			return []*ast.ColumnDef{cmd.Column}
		}
	default:
	}
	return nil
}

// addColumns returns the columns ADD COLUMN carries, its single column and
// then every entry of its parenthesized form, nil entries included. Rules
// that read both fields of the action, as V1's did, use it.
func addColumns(cmd *ast.AlterTableCmd) []*ast.ColumnDef {
	var cols []*ast.ColumnDef
	if cmd.Column != nil {
		cols = append(cols, cmd.Column)
	}
	return append(cols, cmd.Columns...)
}

// isNullable reports whether a column can hold NULL: it has no NOT NULL and
// no PRIMARY KEY constraint of its own. A nil column is nullable.
func isNullable(col *ast.ColumnDef) bool {
	if col == nil {
		return true
	}
	for _, c := range col.Constraints {
		if c.Type == ast.ColConstrNotNull || c.Type == ast.ColConstrPrimaryKey {
			return false
		}
	}
	return true
}

// hasDefault reports whether a column has a DEFAULT clause.
func hasDefault(col *ast.ColumnDef) bool {
	if col == nil {
		return false
	}
	if col.DefaultValue != nil {
		return true
	}
	for _, c := range col.Constraints {
		if c.Type == ast.ColConstrDefault {
			return true
		}
	}
	return false
}

// isPrimaryKey reports whether a column has a PRIMARY KEY constraint of its
// own.
func isPrimaryKey(col *ast.ColumnDef) bool {
	if col == nil {
		return false
	}
	for _, c := range col.Constraints {
		if c.Type == ast.ColConstrPrimaryKey {
			return true
		}
	}
	return false
}

// primaryKeyColumns returns the names of a CREATE TABLE's primary key
// columns: those with an inline PRIMARY KEY and those a PRIMARY KEY table
// constraint lists, names compared as written.
func primaryKeyColumns(n *ast.CreateTableStmt) map[string]bool {
	result := make(map[string]bool)
	for _, col := range n.Columns {
		if isPrimaryKey(col) {
			result[col.Name] = true
		}
	}
	for _, c := range n.Constraints {
		if c.Type == ast.ConstrPrimaryKey {
			for _, name := range c.Columns {
				result[name] = true
			}
		}
	}
	return result
}

// columnFinding returns a finding for a column, anchored at its definition.
func columnFinding(code v1review.Code, table string, col *ast.ColumnDef) v1review.Finding {
	start, end := span(col.Loc)
	return v1review.Finding{Code: code, Name: col.Name, Table: table, Start: start, End: end}
}

// definedColumns returns the table a top-level statement defines columns of,
// without its database qualifier, and those columns, nil entries included:
// every column of a CREATE TABLE, and the alterColumns of each action of an
// ALTER TABLE in order. A statement of another kind, or one whose table has
// no name, defines none.
func definedColumns(stmt ast.Node) (string, []*ast.ColumnDef) {
	switch n := stmt.(type) {
	case *ast.CreateTableStmt:
		table := refName(n.Table)
		if table == "" {
			return "", nil
		}
		return table, n.Columns
	case *ast.AlterTableStmt:
		table := refName(n.Table)
		if table == "" {
			return "", nil
		}
		var cols []*ast.ColumnDef
		for _, cmd := range n.Commands {
			cols = append(cols, alterColumns(cmd)...)
		}
		return table, cols
	default:
		return "", nil
	}
}

// modifiedColumns returns the columns an ALTER TABLE action defines, read
// the way V1's comment and default rules read them: ADD COLUMN's single
// column and then every entry of its parenthesized form, nil entries
// included; MODIFY and CHANGE COLUMN's column; nothing for any other action.
// A nil action is dereferenced.
func modifiedColumns(cmd *ast.AlterTableCmd) []*ast.ColumnDef {
	switch cmd.Type {
	case ast.ATAddColumn:
		return addColumns(cmd)
	case ast.ATModifyColumn, ast.ATChangeColumn:
		if cmd.Column != nil {
			return []*ast.ColumnDef{cmd.Column}
		}
	default:
	}
	return nil
}

// wholeStatementFinding returns a finding anchored at a whole statement, for a
// rule whose V1 advisor printed the statement and placed its advice at the
// statement's first line; the caller keeps that computation.
func wholeStatementFinding(code v1review.Code, name string, loc ast.Loc) v1review.Finding {
	start, end := span(loc)
	return v1review.Finding{Code: code, Name: name, Start: start, End: end}
}
