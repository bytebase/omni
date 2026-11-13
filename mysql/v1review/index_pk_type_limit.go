package v1review

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// IndexPKTypeLimit reports the primary key columns of a script whose type is
// neither INT nor BIGINT: a column defined with the PRIMARY KEY attribute,
// and each column a table-level PRIMARY KEY lists, in a CREATE TABLE or an
// ALTER TABLE. How the script is walked, how a listed column's type is
// found, and when the function returns v1review.ErrSchemaMissing are
// judgeKeyColumns'.
//
// A defined column's type is INT for INT and INTEGER, BIGINT for BIGINT, and
// for any other the type's name in lower case with its length, or length and
// scale, when they are above zero; UNSIGNED and every other attribute are
// dropped, so INT UNSIGNED passes. A type read from schema is judged as the
// snapshot spells it, INT and BIGINT in any letter case, so `int(11)` and
// `bigint unsigned` do not pass.
//
// Name is the column, Table the table's name without its database
// qualifier, and Value the type as judged. Ported from Bytebase's V1 rule
// INDEX_PK_TYPE_LIMIT.
func IndexPKTypeLimit(stmts []ast.Node, schema v1review.Schema) ([]v1review.Finding, error) {
	return judgeKeyColumns(stmts, schema, keyColumnRule{
		code: v1review.PrimaryKeyTypeNotAllowed,
		typeOf: func(t *ast.DataType) string {
			switch strings.ToUpper(t.Name) {
			case "INT", "INTEGER":
				return "INT"
			case "BIGINT":
				return "BIGINT"
			}
			name := strings.ToLower(t.Name)
			switch {
			case t.Length > 0 && t.Scale > 0:
				return fmt.Sprintf("%s(%d,%d)", name, t.Length, t.Scale)
			case t.Length > 0:
				return fmt.Sprintf("%s(%d)", name, t.Length)
			}
			return name
		},
		attribute:  func(t ast.ColumnConstraintType) bool { return t == ast.ColConstrPrimaryKey },
		constraint: func(t ast.ConstraintType) bool { return t == ast.ConstrPrimaryKey },
		allowed: func(columnType string) bool {
			return strings.EqualFold(columnType, "INT") || strings.EqualFold(columnType, "BIGINT")
		},
	})
}
