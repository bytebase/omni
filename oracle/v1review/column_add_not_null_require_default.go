package v1review

import (
	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnAddNotNullRequireDefault judges the columns top-level ALTER TABLE
// statements add for NOT NULL without a DEFAULT. Every column of an ADD
// action is judged, its column list and then its single column. A column is
// NOT NULL when it says NOT NULL or carries a NOT NULL constraint; a
// PRIMARY KEY alone does not make it so. Columns a CREATE TABLE or a MODIFY
// action defines are not judged. Each finding has code
// ColumnNotNullNoDefault, Name the column as the parser stores it, and
// anchors the column definition. Ported from Bytebase's V1 rule
// COLUMN_ADD_NOT_NULL_REQUIRE_DEFAULT.
func ColumnAddNotNullRequireDefault(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		n, ok := stmt.(*ast.AlterTableStmt)
		if !ok {
			continue
		}
		for _, cmd := range alterTableCmds(n) {
			if cmd.Action != ast.AT_ADD_COLUMN {
				continue
			}
			for _, col := range actionColumns(cmd) {
				if col == nil || col.Default != nil {
					continue
				}
				if col.NotNull || columnHasConstraint(col, ast.CONSTRAINT_NOT_NULL) {
					start, end := span(col.Loc)
					findings = append(findings, v1review.Finding{Code: v1review.ColumnNotNullNoDefault, Name: col.Name, Start: start, End: end})
				}
			}
		}
	}
	return findings
}
