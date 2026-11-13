package v1review

import (
	"strings"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnRequireDefault judges each column a top-level CREATE TABLE defines and
// each column an ALTER ... ADD COLUMN adds, whatever relation the ALTER names,
// for a DEFAULT. A column has one when it carries a raw default or a DEFAULT
// constraint, or when the last part of its type's name, lower-cased, is
// serial, bigserial, smallserial, serial2, serial4 or serial8. Identity and
// generated columns have none, and neither does a column of a typed or
// partition table that writes no type. A statement whose relation has no name
// is skipped, and a statement in a routine body is not judged.
//
// Each finding reports ColumnNoDefault with Name the column and Table the
// relation's name without its schema, in the order the columns are written.
// It anchors the whole statement and sets Search to the column name:
// Bytebase's V1 advisor reported the line of the name's first occurrence in
// the statement text. Ported from Bytebase's V1 rule COLUMN_REQUIRE_DEFAULT.
func ColumnRequireDefault(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		var table string
		var loc ast.Loc
		var columns []*ast.ColumnDef
		switch n := statementNode(stmt).(type) {
		case *ast.CreateStmt:
			table, loc, columns = relname(n.Relation), n.Loc, tableColumns(n)
		case *ast.AlterTableStmt:
			table, loc = relname(n.Relation), n.Loc
			for _, cmd := range alterCommands(n) {
				if ast.AlterTableType(cmd.Subtype) != ast.AT_AddColumn {
					continue
				}
				if col, ok := cmd.Def.(*ast.ColumnDef); ok && col != nil {
					columns = append(columns, col)
				}
			}
		default:
			continue
		}
		if table == "" {
			continue
		}
		for _, col := range columns {
			if hasDefault(col) {
				continue
			}
			start, end := span(loc)
			findings = append(findings, v1review.Finding{
				Code:   v1review.ColumnNoDefault,
				Name:   col.Colname,
				Table:  table,
				Search: col.Colname,
				Start:  start,
				End:    end,
			})
		}
	}
	return findings
}

func hasDefault(col *ast.ColumnDef) bool {
	if col.RawDefault != nil {
		return true
	}
	switch strings.ToLower(lastTypeName(col.TypeName)) {
	case "serial", "bigserial", "smallserial", "serial4", "serial8", "serial2":
		return true
	default:
	}
	for _, c := range columnConstraints(col) {
		if c.Contype == ast.CONSTR_DEFAULT {
			return true
		}
	}
	return false
}
