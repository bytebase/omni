package v1review

import (
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnDisallowSetCharset judges the columns top-level statements define
// for a CHARACTER SET other than binary: every column a CREATE TABLE
// defines, and in an ALTER TABLE the columns of ADD COLUMN in both its
// forms, CHANGE COLUMN and MODIFY COLUMN. The character set is compared as
// the parser kept it, so only the exact text "binary" is allowed. A column
// with no type is not judged; a CREATE TABLE or ALTER TABLE with no table is
// not judged, though one whose table has an empty name is. Statements nested
// in a routine body are not judged.
//
// Each finding names the column, holds the character set in Value, and
// anchors the whole statement, since V1 printed the statement and placed the
// advice at its first line. Findings are in statement order, then column
// order. Ported from Bytebase's V1 rule COLUMN_DISALLOW_SET_CHARSET.
func ColumnDisallowSetCharset(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	check := func(loc ast.Loc, col *ast.ColumnDef) {
		if col == nil || col.TypeName == nil {
			return
		}
		switch col.TypeName.Charset {
		case "", "binary":
			return
		}
		finding := statementFinding(v1review.ColumnCharsetSet, col.Name, loc)
		finding.Value = col.TypeName.Charset
		findings = append(findings, finding)
	}
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.Table == nil {
				continue
			}
			for _, col := range n.Columns {
				check(n.Loc, col)
			}
		case *ast.AlterTableStmt:
			if n.Table == nil {
				continue
			}
			for _, cmd := range n.Commands {
				if cmd == nil {
					continue
				}
				switch cmd.Type {
				case ast.ATAddColumn:
					for _, col := range addColumns(cmd) {
						check(n.Loc, col)
					}
				case ast.ATChangeColumn, ast.ATModifyColumn:
					check(n.Loc, cmd.Column)
				default:
				}
			}
		}
	}
	return findings
}
