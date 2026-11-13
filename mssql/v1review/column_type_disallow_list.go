package v1review

import (
	"slices"
	"strings"

	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnTypeDisallowList judges every data type anywhere in the statements,
// nested ones included (a CAST, a DECLARE, a routine body), for a name that,
// upper-cased, is one of disallowed, each entry also upper-cased. Only the
// bare name is compared, without its schema or length, so "VARCHAR(10)"
// matches nothing. A data type node the walk reaches twice is reported
// twice. Each finding has code ColumnTypeDisallowed, Value the upper-cased
// type name, which is all V1 printed, and anchors the data type. Ported from
// Bytebase's V1 rule COLUMN_TYPE_DISALLOW_LIST.
func ColumnTypeDisallowList(stmts []ast.Node, disallowed []string) []v1review.Finding {
	upper := make([]string, 0, len(disallowed))
	for _, t := range disallowed {
		upper = append(upper, strings.ToUpper(t))
	}
	var findings []v1review.Finding
	for _, stmt := range stmts {
		ast.Inspect(stmt, func(n ast.Node) bool {
			dt, ok := n.(*ast.DataType)
			if !ok || dt == nil {
				return true
			}
			name := strings.ToUpper(dt.Name)
			if slices.Contains(upper, name) {
				start, end := span(dt.Loc)
				findings = append(findings, v1review.Finding{Code: v1review.ColumnTypeDisallowed, Value: name, Start: start, End: end})
			}
			return true
		})
	}
	return findings
}
