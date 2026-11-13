package v1review

import (
	"slices"
	"strings"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnNoNull judges the statements together for columns they define that
// final, the schema after the change, holds as nullable. The columns judged
// are every column a top-level CREATE TABLE defines, anchored at the column
// definition, and the columns ALTER TABLE's ADD and CHANGE COLUMN write,
// anchored at the whole statement; MODIFY COLUMN is not judged. A statement
// naming no table is skipped. A column is judged once, at its first
// definition, keyed by its table and its name joined with a dot, so table
// "a.b" with column "c" and table "a" with column "b.c" are one column.
//
// The column is looked up in final, in the schema named "", by the table's
// name without a qualifier, column names compared without regard to case;
// one final lacks is not reported. Name is the column and Table the table as
// written; Statement says which statement holds the anchor. Findings are in
// the order of their anchors, ties broken by column name.
//
// V1 ordered by line, then column name; a caller placing findings on lines
// re-sorts stably by line and Name to match it, two statements sharing a
// line. Two tables' columns of the same name on one line came out in either
// order from V1; here the earlier statement comes first. Ported from
// Bytebase's V1 rule COLUMN_NO_NULL.
func ColumnNoNull(stmts []ast.Node, final v1review.Schema) []v1review.Finding {
	type candidate struct {
		table, column string
		statement     int
		loc           ast.Loc
	}
	seen := make(map[string]bool)
	var candidates []candidate
	add := func(i int, loc ast.Loc, table, column string) {
		key := table + "." + column
		if seen[key] {
			return
		}
		seen[key] = true
		candidates = append(candidates, candidate{table, column, i, loc})
	}
	for i, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.Table == nil {
				continue
			}
			for _, col := range n.Columns {
				if col != nil {
					add(i, col.Loc, n.Table.Name, col.Name)
				}
			}
		case *ast.AlterTableStmt:
			if n.Table == nil {
				continue
			}
			for _, cmd := range n.Commands {
				if cmd == nil {
					continue
				}
				var cols []*ast.ColumnDef
				switch cmd.Type {
				case ast.ATAddColumn:
					cols = addedColumns(cmd)
				case ast.ATChangeColumn:
					if cmd.Column != nil {
						cols = []*ast.ColumnDef{cmd.Column}
					}
				default:
				}
				for _, col := range cols {
					if col != nil {
						add(i, n.Loc, n.Table.Name, col.Name)
					}
				}
			}
		}
	}
	slices.SortStableFunc(candidates, func(a, b candidate) int {
		if a.statement != b.statement {
			return a.statement - b.statement
		}
		if a.loc.Start != b.loc.Start {
			return a.loc.Start - b.loc.Start
		}
		return strings.Compare(a.column, b.column)
	})
	var findings []v1review.Finding
	for _, c := range candidates {
		column := final.TableColumn("", c.table, c.column, false)
		if column == nil || !column.GetNullable() {
			continue
		}
		start, end := span(c.loc)
		findings = append(findings, v1review.Finding{
			Code:      v1review.ColumnNullable,
			Name:      c.column,
			Table:     c.table,
			Start:     start,
			End:       end,
			Statement: c.statement,
		})
	}
	return findings
}
