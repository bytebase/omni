package v1review

import (
	"slices"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnRequired judges the statements together for tables that lack a
// required column, each table and column named without a qualifier and
// compared bytewise.
//
// A top-level CREATE TABLE starts its table afresh with the columns it
// defines, so one created by LIKE or AS SELECT has none. An ALTER TABLE of a
// table no earlier statement created is taken to have every required
// column. ADD COLUMN, in both its forms, adds a column; DROP COLUMN drops
// one; RENAME COLUMN and CHANGE COLUMN drop the old name and add the new. A
// DROP TABLE forgets the table. MODIFY COLUMN is not judged.
//
// Each finding is one missing column, in table name order and then column
// name order. Its anchor is the whole statement that last did one of these
// to the table: CREATE TABLE; RENAME COLUMN, whatever columns it names; DROP
// COLUMN of a required column; CHANGE COLUMN from a required column.
// Statement says which statement holds it. A DROP TABLE does not clear the
// anchor. A table never so anchored is anchored at 0,0 in statement 0, but
// such a table never lacks a column. Ported from Bytebase's V1 rule
// COLUMN_REQUIRED.
func ColumnRequired(stmts []ast.Node, required []string) []v1review.Finding {
	r := &columnRequired{
		required: make(map[string]bool),
		tables:   make(map[string]map[string]bool),
		anchors:  make(map[string]int),
	}
	for _, column := range required {
		r.required[column] = true
	}
	for i, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.Table == nil {
				continue
			}
			table := n.Table.Name
			r.anchors[table] = i
			r.tables[table] = make(map[string]bool)
			for _, col := range n.Columns {
				if col != nil {
					r.add(table, col.Name)
				}
			}
		case *ast.DropTableStmt:
			for _, ref := range n.Tables {
				if ref != nil {
					delete(r.tables, ref.Name)
				}
			}
		case *ast.AlterTableStmt:
			if n.Table == nil {
				continue
			}
			table := n.Table.Name
			for _, cmd := range n.Commands {
				if cmd == nil {
					continue
				}
				switch cmd.Type {
				case ast.ATRenameColumn:
					r.rename(table, cmd.Name, cmd.NewName)
					r.anchors[table] = i
				case ast.ATAddColumn:
					for _, col := range addedColumns(cmd) {
						if col != nil {
							r.add(table, col.Name)
						}
					}
				case ast.ATDropColumn:
					if r.drop(table, cmd.Name) {
						r.anchors[table] = i
					}
				case ast.ATChangeColumn:
					if cmd.Column != nil && r.rename(table, cmd.Name, cmd.Column.Name) {
						r.anchors[table] = i
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
		statement, ok := r.anchors[table]
		var start, end int
		if ok {
			start, end = statementSpan(stmts[statement])
		}
		for _, column := range missing {
			findings = append(findings, v1review.Finding{
				Code:      v1review.ColumnMissing,
				Name:      column,
				Table:     table,
				Start:     start,
				End:       end,
				Statement: statement,
			})
		}
	}
	return findings
}

// statementSpan returns the byte range of a CREATE TABLE or ALTER TABLE
// statement.
func statementSpan(stmt ast.Node) (int, int) {
	switch n := stmt.(type) {
	case *ast.CreateTableStmt:
		return span(n.Loc)
	case *ast.AlterTableStmt:
		return span(n.Loc)
	default:
		return 0, 0
	}
}

type columnRequired struct {
	required map[string]bool
	// tables holds, per table, whether each required column is present.
	tables map[string]map[string]bool
	// anchors holds, per table, the index of the statement anchoring it.
	anchors map[string]int
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
