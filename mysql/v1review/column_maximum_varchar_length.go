package v1review

import (
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnMaximumVarcharLength judges the VARCHAR columns (type name in any
// case) top-level statements define for a length above maximum: every
// column a CREATE TABLE defines, and in an ALTER TABLE the columns of ADD
// COLUMN in both its forms, MODIFY COLUMN and CHANGE COLUMN. A VARCHAR
// without a length, or VARCHAR(0), counts as length 1. A maximum of 0 or
// less judges nothing. A statement whose table has no name is not judged,
// and neither are statements nested in a routine body.
//
// Each finding names the column and its table, without the database
// qualifier, and anchors the column definition. Findings are in statement
// order, then column order. Ported from Bytebase's V1 rule
// COLUMN_MAXIMUM_VARCHAR_LENGTH.
func ColumnMaximumVarcharLength(stmts []ast.Node, maximum int) []v1review.Finding {
	return typeLengthAbove(stmts, "VARCHAR", maximum, v1review.ColumnVarcharTooLong)
}
