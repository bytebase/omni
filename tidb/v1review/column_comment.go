package v1review

import (
	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnComment judges the comments of the columns top-level CREATE TABLE
// and ALTER TABLE statements define: every column a CREATE TABLE defines,
// anchored at the column definition, and the columns ALTER TABLE's ADD,
// CHANGE and MODIFY COLUMN write, anchored at the whole statement. A
// statement naming no table is skipped.
//
// When required is set, a column with no COMMENT clause is reported
// ColumnCommentMissing; a COMMENT clause with an empty string counts as one.
// A column whose comment is longer than maxLength bytes is reported
// ColumnCommentTooLong, a column with no COMMENT clause having the empty
// comment, so a negative maxLength reports every column; a caller that means
// no limit drops those findings. A column with both is reported missing
// first. Name is the column and Table the table, both without a qualifier.
// Ported from Bytebase's V1 rule COLUMN_COMMENT.
func ColumnComment(stmts []ast.Node, required bool, maxLength int) []v1review.Finding {
	var findings []v1review.Finding
	walkColumns(stmts, func(table string, col *ast.ColumnDef, loc ast.Loc) {
		if required && !hasColumnConstraint(col, ast.ColConstrComment) {
			findings = append(findings, columnFinding(v1review.ColumnCommentMissing, table, col.Name, loc))
		}
		if len(col.Comment) > maxLength {
			findings = append(findings, columnFinding(v1review.ColumnCommentTooLong, table, col.Name, loc))
		}
	})
	return findings
}
