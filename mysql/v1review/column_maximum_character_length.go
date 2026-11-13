package v1review

import (
	"strings"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnMaximumCharacterLength judges the CHAR columns (type name in any
// case) top-level statements define for a length above maximum: every column
// a CREATE TABLE defines, and in an ALTER TABLE the columns of ADD COLUMN in
// both its forms, MODIFY COLUMN and CHANGE COLUMN. A CHAR written without a
// length, or as CHAR(0), counts as length 1. A maximum of 0 or less judges
// nothing. A statement whose table has no name is not judged, and neither
// are statements nested in a routine body.
//
// Each finding names the column and its table, without the database
// qualifier, and anchors the column definition. Findings are in statement
// order, then column order. Ported from Bytebase's V1 rule
// COLUMN_MAXIMUM_CHARACTER_LENGTH.
func ColumnMaximumCharacterLength(stmts []ast.Node, maximum int) []v1review.Finding {
	return typeLengthAbove(stmts, "CHAR", maximum, v1review.ColumnCharTooLong)
}

// typeLengthAbove reports the columns of the named type, matched in any
// case, whose length is above maximum; a length of 0 counts as 1.
func typeLengthAbove(stmts []ast.Node, typeName string, maximum int, code v1review.Code) []v1review.Finding {
	if maximum <= 0 {
		return nil
	}
	var findings []v1review.Finding
	for _, stmt := range stmts {
		table, cols := definedColumns(stmt)
		for _, col := range cols {
			dt := col.TypeName
			if dt == nil || !strings.EqualFold(dt.Name, typeName) {
				continue
			}
			length := dt.Length
			if length == 0 {
				length = 1
			}
			if length > maximum {
				findings = append(findings, columnFinding(code, table, col))
			}
		}
	}
	return findings
}
