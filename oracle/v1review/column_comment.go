package v1review

import (
	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnComment judges the top-level statements together for column
// comments. Each column is keyed table.column, the table as schema.name
// with the schema defaulting to currentDatabase. The columns a CREATE TABLE
// defines, and those an ADD action of an ALTER TABLE defines (its column
// list, then its single column), are judged, each time one is defined, so
// a column defined twice is judged twice. A COMMENT ON COLUMN anywhere among
// the statements, before or after the definition, gives a column its
// comment, the last one winning. A column with no comment, or whose last
// comment is empty, has code ColumnCommentMissing when required; one whose
// comment is longer than maxLength bytes has code ColumnCommentTooLong, so
// a negative maxLength reports every commented column. Nothing removes a
// column: DROP COLUMN and DROP TABLE are not judged.
//
// Findings come in definition order, Table the table's key and Name the
// column as the parser stores it. Every finding of one key takes its
// Statement and range from the last definition of that key. Ported from
// Bytebase's V1 rule COLUMN_COMMENT.
func ColumnComment(stmts []ast.Node, currentDatabase string, required bool, maxLength int) []v1review.Finding {
	type column struct{ table, name, key string }
	var columns []column
	anchors := make(map[string]columnAnchor)
	comments := make(map[string]string)
	define := func(i int, table string, col *ast.ColumnDef) {
		key := table + "." + col.Name
		columns = append(columns, column{table: table, name: col.Name, key: key})
		anchors[key] = columnAnchor{statement: i, loc: col.Loc}
	}
	for i, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			table := objectName(n.Name, currentDatabase)
			for _, col := range columnDefs(n.Columns) {
				define(i, table, col)
			}
		case *ast.AlterTableStmt:
			table := objectName(n.Name, currentDatabase)
			for _, cmd := range alterTableCmds(n) {
				if cmd.Action != ast.AT_ADD_COLUMN {
					continue
				}
				for _, col := range actionColumns(cmd) {
					if col != nil {
						define(i, table, col)
					}
				}
			}
		case *ast.CommentStmt:
			if n.ObjectType != ast.OBJECT_TABLE || n.Column == "" {
				continue
			}
			comments[objectName(n.Object, currentDatabase)+"."+n.Column] = n.Comment
		}
	}

	var findings []v1review.Finding
	for _, c := range columns {
		var code v1review.Code
		if comment := comments[c.key]; comment == "" {
			if !required {
				continue
			}
			code = v1review.ColumnCommentMissing
		} else if len(comment) > maxLength {
			code = v1review.ColumnCommentTooLong
		} else {
			continue
		}
		anchor := anchors[c.key]
		start, end := span(anchor.loc)
		findings = append(findings, v1review.Finding{Code: code, Name: c.name, Table: c.table, Start: start, End: end, Statement: anchor.statement})
	}
	return findings
}
