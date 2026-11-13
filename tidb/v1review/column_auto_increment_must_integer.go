package v1review

import (
	"strings"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnAutoIncrementMustInteger reports the AUTO_INCREMENT columns whose
// type is not an integer type: TINYINT, SMALLINT, MEDIUMINT, INT, INTEGER,
// BIGINT, or BOOLEAN (which BOOL also parses to), compared without regard
// to case. A column with no type is reported.
//
// The columns judged are those of top-level statements only: every column a
// CREATE TABLE defines, anchored at the column definition, and the columns
// ALTER TABLE's ADD, CHANGE and MODIFY COLUMN write, anchored at the whole
// statement. Name is the column and Table the table, both without a
// qualifier. Ported from Bytebase's V1 rule
// COLUMN_AUTO_INCREMENT_MUST_INTEGER.
func ColumnAutoIncrementMustInteger(stmts []ast.Node) []v1review.Finding {
	return columnRule(stmts, v1review.ColumnAutoIncrementNotInteger, func(col *ast.ColumnDef) bool {
		return col.AutoIncrement && !isIntegerType(col.TypeName)
	})
}

func isIntegerType(dt *ast.DataType) bool {
	if dt == nil {
		return false
	}
	switch strings.ToUpper(dt.Name) {
	case "TINYINT", "SMALLINT", "MEDIUMINT", "INT", "INTEGER", "BIGINT", "BOOLEAN":
		return true
	default:
		return false
	}
}
