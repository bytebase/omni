package v1review

import (
	"strings"

	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnTypeDisallowList judges the columns top-level statements define for
// a type whose name, upper-cased, is one of types, also upper-cased: every
// column a CREATE TABLE defines, and in an ALTER TABLE the columns of ADD
// COLUMN in both its forms, MODIFY COLUMN and CHANGE COLUMN. Only the bare
// type name is compared, so "VARCHAR(20)" in types matches nothing and
// "VARCHAR" matches every VARCHAR. A column with no type is not judged; a
// statement whose table has no name is not judged, and neither are
// statements nested in a routine body.
//
// Each finding names the column and its table, without the database
// qualifier, holds the upper-cased type name in Value, and anchors the
// column definition. Findings are in statement order, then column order.
// Ported from Bytebase's V1 rule COLUMN_TYPE_DISALLOW_LIST.
func ColumnTypeDisallowList(stmts []ast.Node, types []string) []v1review.Finding {
	disallowed := make(map[string]bool, len(types))
	for _, t := range types {
		disallowed[strings.ToUpper(t)] = true
	}
	var findings []v1review.Finding
	for _, stmt := range stmts {
		table, cols := definedColumns(stmt)
		for _, col := range cols {
			dt := col.TypeName
			if dt == nil {
				continue
			}
			columnType := strings.ToUpper(strings.ToLower(dt.Name))
			if disallowed[columnType] {
				finding := columnFinding(v1review.ColumnTypeDisallowed, table, col)
				finding.Value = columnType
				findings = append(findings, finding)
			}
		}
	}
	return findings
}
