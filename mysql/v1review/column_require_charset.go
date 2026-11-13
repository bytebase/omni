package v1review

import (
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnRequireCharset judges the columns top-level statements define for a
// CHAR, VARCHAR, TINYTEXT, TEXT, MEDIUMTEXT or LONGTEXT type (name in any
// case) that names no CHARACTER SET: every column a CREATE TABLE defines,
// and the columns of ADD COLUMN, in both its forms, in an ALTER TABLE.
// MODIFY and CHANGE COLUMN are not judged. A column with no type is not
// judged; a CREATE TABLE or ALTER TABLE with no table is not judged, though
// one whose table has an empty name is. Statements nested in a routine body
// are not judged.
//
// Each finding names the column. A CREATE TABLE finding anchors the column
// definition; an ALTER TABLE finding anchors the whole ADD COLUMN action,
// whichever line the column is on. Findings are in statement order, then
// column order. Ported from Bytebase's V1 rule COLUMN_REQUIRE_CHARSET.
func ColumnRequireCharset(stmts []ast.Node) []v1review.Finding {
	return unnamedTypeOption(stmts, v1review.ColumnNoCharset, func(dt *ast.DataType) string { return dt.Charset })
}

// unnamedTypeOption reports the character columns whose type leaves the
// option option reads empty, as ColumnRequireCharset describes.
func unnamedTypeOption(stmts []ast.Node, code v1review.Code, option func(*ast.DataType) string) []v1review.Finding {
	var findings []v1review.Finding
	check := func(loc ast.Loc, col *ast.ColumnDef) {
		if col == nil || col.TypeName == nil || !isCharsetType(col.TypeName) || option(col.TypeName) != "" {
			return
		}
		start, end := span(loc)
		findings = append(findings, v1review.Finding{Code: code, Name: col.Name, Start: start, End: end})
	}
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.Table == nil {
				continue
			}
			for _, col := range n.Columns {
				check(colLoc(col), col)
			}
		case *ast.AlterTableStmt:
			if n.Table == nil {
				continue
			}
			for _, cmd := range n.Commands {
				if cmd == nil || cmd.Type != ast.ATAddColumn {
					continue
				}
				for _, col := range addColumns(cmd) {
					check(cmd.Loc, col)
				}
			}
		}
	}
	return findings
}

// colLoc returns a column's range, the zero Loc for a nil column.
func colLoc(col *ast.ColumnDef) ast.Loc {
	if col == nil {
		return ast.Loc{}
	}
	return col.Loc
}
