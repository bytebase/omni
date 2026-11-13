package v1review

import (
	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnMaximumVarcharLength judges every column definition anywhere in the
// statements, nested ones included, whose type is VARCHAR or VARCHAR2 (its
// name parts joined by dots and upper-cased, so NVARCHAR2 and a qualified
// name are not judged) for a length above maximum. The length is the type's
// first modifier when it is an integer; a column with no such modifier is
// not judged, and no unit (BYTE, CHAR) is considered. Each finding has code
// ColumnVarcharTooLong, Name the column as the parser stores it, and anchors
// the column's type. Ported from Bytebase's V1 rule
// COLUMN_MAXIMUM_VARCHAR_LENGTH.
func ColumnMaximumVarcharLength(stmts []ast.Node, maximum int) []v1review.Finding {
	return columnLengthFindings(stmts, maximum, v1review.ColumnVarcharTooLong, "VARCHAR", "VARCHAR2")
}
