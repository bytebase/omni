package v1review

import (
	"slices"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnRequired judges the statements together for tables that lack a
// required column, each table named without its database qualifier.
//
// A CREATE TABLE starts its table afresh with the columns it defines, so one
// created by LIKE or AS SELECT has none. An ALTER TABLE of a table no earlier
// statement created is taken to have every required column. ADD COLUMN, in
// both its forms, adds a column; DROP COLUMN drops one; RENAME COLUMN and
// CHANGE COLUMN drop the old name and add the new. A DROP TABLE forgets the
// table. MODIFY COLUMN is not judged, and neither are statements nested in a
// routine body.
//
// Each finding is one missing column, in table name order and then column
// name order, both compared bytewise. Its anchor is the last of these that
// touched the table: the CREATE TABLE statement, or the DROP, RENAME or
// CHANGE COLUMN action that removed a required column; Statement says which
// statement holds it. A table whose only anchor was never set, one altered
// without ever losing a required column, is anchored at 0,0 in statement 0,
// but such a table never lacks a column. Ported from Bytebase's V1 rule
// COLUMN_REQUIRED.
func ColumnRequired(stmts []ast.Node, required []string) []v1review.Finding {
	r := &columnRequired{
		required: make(map[string]bool),
		tables:   make(map[string]map[string]bool),
		anchors:  make(map[string]columnAnchor),
	}
	for _, column := range required {
		r.required[column] = true
	}
	for i, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			table := refName(n.Table)
			if table == "" {
				continue
			}
			r.anchors[table] = columnAnchor{statement: i, loc: n.Loc}
			r.tables[table] = make(map[string]bool)
			for _, col := range n.Columns {
				r.add(table, col.Name)
			}
		case *ast.DropTableStmt:
			for _, ref := range n.Tables {
				delete(r.tables, refName(ref))
			}
		case *ast.AlterTableStmt:
			table := refName(n.Table)
			if table == "" {
				continue
			}
			for _, cmd := range n.Commands {
				anchor := columnAnchor{statement: i, loc: cmd.Loc}
				switch cmd.Type {
				case ast.ATAddColumn:
					if cmd.Column != nil {
						r.add(table, cmd.Column.Name)
					}
					for _, col := range cmd.Columns {
						r.add(table, col.Name)
					}
				case ast.ATDropColumn:
					if r.drop(table, cmd.Name) {
						r.anchors[table] = anchor
					}
				case ast.ATRenameColumn, ast.ATChangeColumn:
					if r.rename(table, cmd.Name, cmd.NewName) {
						r.anchors[table] = anchor
					}
				default:
				}
			}
		}
	}

	var findings []v1review.Finding
	tables := make([]string, 0, len(r.tables))
	for table := range r.tables {
		tables = append(tables, table)
	}
	slices.Sort(tables)
	for _, table := range tables {
		var missing []string
		for column := range r.required {
			if !r.tables[table][column] {
				missing = append(missing, column)
			}
		}
		slices.Sort(missing)
		anchor := r.anchors[table]
		start, end := span(anchor.loc)
		for _, column := range missing {
			findings = append(findings, v1review.Finding{
				Code:      v1review.ColumnMissing,
				Name:      column,
				Table:     table,
				Start:     start,
				End:       end,
				Statement: anchor.statement,
			})
		}
	}
	return findings
}

// columnAnchor is the node a cross-statement finding takes its line from, and
// the index of the statement holding it.
type columnAnchor struct {
	statement int
	loc       ast.Loc
}

type columnRequired struct {
	required map[string]bool
	// tables holds, per table, whether each required column is present.
	tables  map[string]map[string]bool
	anchors map[string]columnAnchor
}

// add records a required column as present; a table not yet seen is taken
// to have every required column.
func (r *columnRequired) add(table, column string) {
	if !r.required[column] {
		return
	}
	if columns, ok := r.tables[table]; ok {
		columns[column] = true
		return
	}
	r.full(table)
}

// drop records a required column as absent and reports whether it was one.
func (r *columnRequired) drop(table, column string) bool {
	if !r.required[column] {
		return false
	}
	columns, ok := r.tables[table]
	if !ok {
		columns = r.full(table)
	}
	columns[column] = false
	return true
}

// rename records the old name absent and the new present, and reports
// whether the old name was a required column.
func (r *columnRequired) rename(table, oldColumn, newColumn string) bool {
	oldNeed, newNeed := r.required[oldColumn], r.required[newColumn]
	if !oldNeed && !newNeed {
		return false
	}
	columns, ok := r.tables[table]
	if !ok {
		columns = r.full(table)
	}
	if oldNeed {
		columns[oldColumn] = false
	}
	if newNeed {
		columns[newColumn] = true
	}
	return oldNeed
}

// full starts a table with every required column present.
func (r *columnRequired) full(table string) map[string]bool {
	columns := make(map[string]bool, len(r.required))
	for column := range r.required {
		columns[column] = true
	}
	r.tables[table] = columns
	return columns
}

func refName(ref *ast.TableRef) string {
	if ref == nil {
		return ""
	}
	return ref.Name
}
