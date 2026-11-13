package v1review

import (
	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnAutoIncrementMustUnsigned reports the AUTO_INCREMENT columns that
// are neither UNSIGNED nor ZEROFILL; a column with no type is reported.
//
// The columns judged are those of top-level statements only: every column a
// CREATE TABLE defines, anchored at the column definition, and the columns
// ALTER TABLE's ADD, CHANGE and MODIFY COLUMN write, anchored at the whole
// statement. Name is the column and Table the table, both without a
// qualifier. Ported from Bytebase's V1 rule
// COLUMN_AUTO_INCREMENT_MUST_UNSIGNED.
func ColumnAutoIncrementMustUnsigned(stmts []ast.Node) []v1review.Finding {
	return columnRule(stmts, v1review.ColumnAutoIncrementNotUnsigned, func(col *ast.ColumnDef) bool {
		return col.AutoIncrement && (col.TypeName == nil || !col.TypeName.Unsigned && !col.TypeName.Zerofill)
	})
}
