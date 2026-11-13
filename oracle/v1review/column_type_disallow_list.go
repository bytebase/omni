package v1review

import (
	"strings"

	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnTypeDisallowList judges every column definition anywhere in the
// statements, nested ones included, for a type in disallowed. The type is
// its name parts joined by dots and upper-cased, without its modifiers, so
// "VARCHAR2" matches VARCHAR2(10) and "VARCHAR2(10)" matches nothing; it is
// compared with each entry case-insensitively (Unicode simple folding). A
// column whose type has no name parts matches an empty entry. A column
// matching several entries is reported once. Each finding has code
// ColumnTypeDisallowed, Name the column as the parser stores it, Value the
// type as described, and anchors the column's type. Ported from Bytebase's
// V1 rule COLUMN_TYPE_DISALLOW_LIST.
func ColumnTypeDisallowList(stmts []ast.Node, disallowed []string) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		ast.Inspect(stmt, func(n ast.Node) bool {
			col, ok := n.(*ast.ColumnDef)
			if !ok || col.TypeName == nil {
				return true
			}
			name := typeName(col.TypeName)
			for _, entry := range disallowed {
				if strings.EqualFold(name, entry) {
					start, end := span(col.TypeName.Loc)
					findings = append(findings, v1review.Finding{Code: v1review.ColumnTypeDisallowed, Name: col.Name, Value: name, Start: start, End: end})
					break
				}
			}
			return true
		})
	}
	return findings
}
