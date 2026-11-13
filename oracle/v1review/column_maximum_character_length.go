package v1review

import (
	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnMaximumCharacterLength judges every column definition anywhere in
// the statements, nested ones included, whose type is CHAR or CHARACTER
// (its name parts joined by dots and upper-cased, so NCHAR and a qualified
// name are not judged) for a length above maximum. The length is the type's
// first modifier when it is an integer; a column with no such modifier is
// not judged, and no unit (BYTE, CHAR) is considered. Each finding has code
// ColumnCharTooLong, Name the column as the parser stores it, and anchors
// the column's type. Ported from Bytebase's V1 rule
// COLUMN_MAXIMUM_CHARACTER_LENGTH.
func ColumnMaximumCharacterLength(stmts []ast.Node, maximum int) []v1review.Finding {
	return columnLengthFindings(stmts, maximum, v1review.ColumnCharTooLong, "CHAR", "CHARACTER")
}

// columnLengthFindings reports each column definition in the statements
// whose type name is typeA or typeB and whose first type modifier is an
// integer above maximum.
func columnLengthFindings(stmts []ast.Node, maximum int, code v1review.Code, typeA, typeB string) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		ast.Inspect(stmt, func(n ast.Node) bool {
			col, ok := n.(*ast.ColumnDef)
			if !ok || col.TypeName == nil {
				return true
			}
			name := typeName(col.TypeName)
			if name != typeA && name != typeB {
				return true
			}
			length, ok := firstTypeModInt(col.TypeName)
			if !ok || length <= maximum {
				return true
			}
			start, end := span(col.TypeName.Loc)
			findings = append(findings, v1review.Finding{Code: code, Name: col.Name, Start: start, End: end})
			return true
		})
	}
	return findings
}
