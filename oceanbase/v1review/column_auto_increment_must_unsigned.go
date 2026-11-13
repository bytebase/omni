package v1review

import (
	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnAutoIncrementMustUnsigned judges the AUTO_INCREMENT columns
// top-level statements define for a type that is neither UNSIGNED nor
// ZEROFILL (which implies UNSIGNED): every column a CREATE TABLE defines,
// and in an ALTER TABLE the columns of ADD COLUMN in both its forms, MODIFY
// COLUMN and CHANGE COLUMN. A column with no type is reported, whatever its
// type would be, so a non-numeric AUTO_INCREMENT column is reported too. A
// statement whose table has no name is not judged, and neither are
// statements nested in a routine body.
//
// Each finding names the column and its table, without the database
// qualifier, and anchors the column definition. Findings are in statement
// order, then column order. Ported from Bytebase's V1 rule
// COLUMN_AUTO_INCREMENT_MUST_UNSIGNED.
func ColumnAutoIncrementMustUnsigned(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		table, cols := definedColumns(stmt)
		for _, col := range cols {
			if col == nil || !col.AutoIncrement {
				continue
			}
			if col.TypeName == nil || (!col.TypeName.Unsigned && !col.TypeName.Zerofill) {
				findings = append(findings, columnFinding(v1review.ColumnAutoIncrementNotUnsigned, table, col))
			}
		}
	}
	return findings
}
