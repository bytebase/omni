package v1review

import (
	"strings"

	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnAutoIncrementMustInteger judges the AUTO_INCREMENT columns top-level
// statements define for a type other than INT, TINYINT, SMALLINT, MEDIUMINT,
// BIGINT or INTEGER, the type name matched in any case: every column a
// CREATE TABLE defines, and in an ALTER TABLE the columns of ADD COLUMN in
// both its forms, MODIFY COLUMN and CHANGE COLUMN. A column with no type is
// reported. A statement whose table has no name is not judged, and neither
// are statements nested in a routine body.
//
// Each finding names the column and its table, without the database
// qualifier, and anchors the column definition. Findings are in statement
// order, then column order. Ported from Bytebase's V1 rule
// COLUMN_AUTO_INCREMENT_MUST_INTEGER.
func ColumnAutoIncrementMustInteger(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		table, cols := definedColumns(stmt)
		for _, col := range cols {
			if col == nil || !col.AutoIncrement || isIntegerType(col.TypeName) {
				continue
			}
			findings = append(findings, columnFinding(v1review.ColumnAutoIncrementNotInteger, table, col))
		}
	}
	return findings
}

func isIntegerType(dt *ast.DataType) bool {
	if dt == nil {
		return false
	}
	switch strings.ToUpper(dt.Name) {
	case "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "BIGINT", "INTEGER":
		return true
	default:
		return false
	}
}
