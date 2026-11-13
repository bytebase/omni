package v1review

import (
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnRequireCollation judges the columns top-level statements define for
// a CHAR, VARCHAR, TINYTEXT, TEXT, MEDIUMTEXT or LONGTEXT type (name in any
// case) that names no COLLATE: every column a CREATE TABLE defines, and the
// columns of ADD COLUMN, in both its forms, in an ALTER TABLE. MODIFY and
// CHANGE COLUMN are not judged. A column with no type is not judged; a
// CREATE TABLE or ALTER TABLE with no table is not judged, though one whose
// table has an empty name is. Statements nested in a routine body are not
// judged.
//
// Each finding names the column. A CREATE TABLE finding anchors the column
// definition; an ALTER TABLE finding anchors the whole ADD COLUMN action,
// whichever line the column is on. Findings are in statement order, then
// column order. Ported from Bytebase's V1 rule COLUMN_REQUIRE_COLLATION.
func ColumnRequireCollation(stmts []ast.Node) []v1review.Finding {
	return unnamedTypeOption(stmts, v1review.ColumnNoCollation, func(dt *ast.DataType) string { return dt.Collate })
}
