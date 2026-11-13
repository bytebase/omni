package v1review

import (
	"math"
	"strings"

	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnMaximumVarcharLength judges every data type anywhere in the
// statements, nested ones included (a CAST, a DECLARE, a routine body), whose
// name lower-cased is varchar, nvarchar, char or nchar, for a length above
// maximum. A data type node reached twice in one statement is judged once.
// MAX counts as math.MaxInt32; otherwise the length is the integer value of a
// literal length, so a type with no length, or a length that is not a
// literal, counts as 0 and is never reported, and a string literal counts as
// its integer value, 0. Each finding has code ColumnVarcharTooLong, leaves
// Name and Value empty, since V1 printed only the maximum, and anchors the
// data type. Ported from Bytebase's V1 rule COLUMN_MAXIMUM_VARCHAR_LENGTH.
func ColumnMaximumVarcharLength(stmts []ast.Node, maximum int) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		seen := make(map[*ast.DataType]bool)
		ast.Inspect(stmt, func(n ast.Node) bool {
			dt, ok := n.(*ast.DataType)
			if !ok || dt == nil || seen[dt] {
				return true
			}
			seen[dt] = true
			switch strings.ToLower(dt.Name) {
			case "varchar", "nvarchar", "char", "nchar":
			default:
				return true
			}
			length := 0
			if dt.MaxLength {
				length = math.MaxInt32
			} else if lit, ok := dt.Length.(*ast.Literal); ok {
				length = int(lit.Ival)
			}
			if length > maximum {
				start, end := span(dt.Loc)
				findings = append(findings, v1review.Finding{Code: v1review.ColumnVarcharTooLong, Start: start, End: end})
			}
			return true
		})
	}
	return findings
}
