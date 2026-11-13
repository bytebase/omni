package v1review

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
	"github.com/bytebase/omni/v1review"
)

// varcharDefaultLength is the length Snowflake gives a VARCHAR declared
// without one.
// https://docs.snowflake.com/en/sql-reference/data-types-text
const varcharDefaultLength = 16_777_216

// ColumnMaximumVarcharLength judges every type name anywhere in the
// statements, nested ones included (a CAST, an element type, a Snowflake
// Scripting block), spelled VARCHAR in any letter case, for a length above
// maximum. The synonyms STRING, TEXT, CHAR VARYING and NVARCHAR are not
// judged. The length is the type's first parameter, or 16,777,216 when it
// has none. Each finding has code ColumnVarcharTooLong, leaves Name and Value
// empty, since V1 printed only the maximum, and anchors the type name; only
// the range's start is meaningful. Ported from Bytebase's V1 rule
// COLUMN_MAXIMUM_VARCHAR_LENGTH.
func ColumnMaximumVarcharLength(stmts []ast.Node, maximum int) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		ast.Inspect(stmt, func(n ast.Node) bool {
			typeName, ok := n.(*ast.TypeName)
			if !ok || !strings.EqualFold(typeName.Name, "VARCHAR") {
				return true
			}
			length := varcharDefaultLength
			if len(typeName.Params) > 0 {
				length = typeName.Params[0]
			}
			if length > maximum {
				start, end := columnAnchor(typeName.Loc)
				findings = append(findings, v1review.Finding{Code: v1review.ColumnVarcharTooLong, Start: start, End: end})
			}
			return true
		})
	}
	return findings
}

// columnAnchor returns a Loc as a finding's byte range: 0,0 when the AST kept no
// start, and a missing end collapsed to the start. The column rules take
// their line from the start alone.
func columnAnchor(loc ast.Loc) (int, int) {
	if loc.Start < 0 {
		return 0, 0
	}
	return loc.Start, max(loc.End, loc.Start)
}
