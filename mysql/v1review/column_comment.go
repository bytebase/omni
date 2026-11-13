package v1review

import (
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnComment judges the comments of the columns top-level statements
// define: every column a CREATE TABLE defines, and in an ALTER TABLE the
// columns of ADD COLUMN in both its forms, MODIFY COLUMN and CHANGE COLUMN.
// A comment longer than maxLength bytes is ColumnCommentTooLong; a negative
// maxLength sets no limit, and 0 makes every comment too long. A column
// without a comment, or with an empty one, is ColumnCommentMissing when
// required is set. A statement whose table has no name is not judged, and
// neither are statements nested in a routine body.
//
// Each finding names the column and its table, without the database
// qualifier, and anchors the column definition. Findings are in statement
// order, then column order. Every column judged yields a finding when
// required is set and maxLength is 0, which is how a caller can tell that a
// column was judged at all. Ported from Bytebase's V1 rule COLUMN_COMMENT.
func ColumnComment(stmts []ast.Node, required bool, maxLength int) []v1review.Finding {
	var findings []v1review.Finding
	check := func(table string, col *ast.ColumnDef) {
		comment := ""
		if col != nil {
			comment = col.Comment
		}
		if comment != "" && maxLength >= 0 && len(comment) > maxLength {
			findings = append(findings, columnFinding(v1review.ColumnCommentTooLong, table, col))
		}
		if comment == "" && required {
			findings = append(findings, columnFinding(v1review.ColumnCommentMissing, table, col))
		}
	}
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			table := refName(n.Table)
			if table == "" {
				continue
			}
			for _, col := range n.Columns {
				check(table, col)
			}
		case *ast.AlterTableStmt:
			table := refName(n.Table)
			if table == "" {
				continue
			}
			for _, cmd := range n.Commands {
				for _, col := range modifiedColumns(cmd) {
					check(table, col)
				}
			}
		}
	}
	return findings
}
