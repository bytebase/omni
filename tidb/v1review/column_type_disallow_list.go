package v1review

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnTypeDisallowList reports the columns whose type, rendered compactly
// and upper-cased, equals an entry of disallowed upper-cased. The rendering
// keeps a length and scale, an ENUM or SET value list, and the display width
// MySQL gives a bare integer type, and renders BOOLEAN as TINYINT(1) and
// NUMERIC or FIXED as DECIMAL; it drops UNSIGNED, ZEROFILL, a character set,
// and a FLOAT, DOUBLE or REAL precision written without a scale. VARCHAR(255),
// INT(11), ENUM('A','B') and JSON are such renderings. A column with no type
// is not judged.
//
// The columns judged are those of top-level statements only: every column a
// CREATE TABLE defines, anchored at the column definition, and the columns
// ALTER TABLE's ADD, CHANGE and MODIFY COLUMN write, anchored at the whole
// statement. A statement naming no table is skipped. Name is the column,
// Table the table, both without a qualifier, and Value the rendered type.
// Ported from Bytebase's V1 rule COLUMN_TYPE_DISALLOW_LIST.
func ColumnTypeDisallowList(stmts []ast.Node, disallowed []string) []v1review.Finding {
	deny := make(map[string]bool, len(disallowed))
	for _, tp := range disallowed {
		deny[strings.ToUpper(tp)] = true
	}
	return columnValueRule(stmts, v1review.ColumnTypeDisallowed, func(col *ast.ColumnDef) (string, bool) {
		if col.TypeName == nil {
			return "", false
		}
		tp := strings.ToUpper(compactType(col.TypeName))
		return tp, deny[tp]
	})
}

// compactType renders a type in lower case the way TiDB's parser rendered a
// column type compactly before omni.
func compactType(dt *ast.DataType) string {
	lower := strings.ToLower(dt.Name)
	switch {
	case lower == "enum" || lower == "set":
		values := make([]string, len(dt.EnumValues))
		for i, v := range dt.EnumValues {
			values[i] = "'" + strings.ReplaceAll(v, "'", "''") + "'"
		}
		return fmt.Sprintf("%s(%s)", lower, strings.Join(values, ","))
	case dt.Length > 0 && (lower == "decimal" || lower == "numeric" || lower == "fixed"):
		return fmt.Sprintf("decimal(%d,%d)", dt.Length, dt.Scale)
	case dt.Length > 0 && dt.Scale > 0:
		return fmt.Sprintf("%s(%d,%d)", lower, dt.Length, dt.Scale)
	case dt.Length > 0 && lower != "float" && lower != "double" && lower != "real":
		return fmt.Sprintf("%s(%d)", lower, dt.Length)
	}
	switch lower {
	case "decimal", "numeric", "fixed":
		return "decimal(10,0)"
	case "tinyint":
		return "tinyint(4)"
	case "smallint":
		return "smallint(6)"
	case "mediumint":
		return "mediumint(9)"
	case "int", "integer":
		return "int(11)"
	case "bigint":
		return "bigint(20)"
	case "bit":
		return "bit(1)"
	case "binary":
		return "binary(1)"
	case "year":
		return "year(4)"
	case "boolean", "bool":
		return "tinyint(1)"
	default:
		return lower
	}
}
