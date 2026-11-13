package v1review

import (
	"strings"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// IndexPrimaryKeyTypeAllowlist reports the primary key columns of a script
// whose type is not in allow: a column defined with the PRIMARY KEY
// attribute, and each column a table-level PRIMARY KEY lists, in a CREATE
// TABLE or an ALTER TABLE. How the script is walked, how a listed column's
// type is found, and when the function returns v1review.ErrSchemaMissing are
// judgeKeyColumns'.
//
// A defined column's type is its type's name in lower case, without length,
// scale or attributes, so VARCHAR(20) is `varchar`; a column with no type is
// the empty type. A type read from schema is the snapshot's whole spelling
// in lower case, so `int(11)` is not `int`. allow's entries are compared in
// lower case.
//
// Name is the column, Table the table's name without its database
// qualifier, and Value the type as judged. Ported from Bytebase's V1 rule
// INDEX_PRIMARY_KEY_TYPE_ALLOWLIST.
func IndexPrimaryKeyTypeAllowlist(stmts []ast.Node, allow []string, schema v1review.Schema) ([]v1review.Finding, error) {
	allowed := make(map[string]bool, len(allow))
	for _, columnType := range allow {
		allowed[strings.ToLower(columnType)] = true
	}
	primaryKey := func(t ast.ColumnConstraintType) bool { return t == ast.ColConstrPrimaryKey }
	return judgeKeyColumns(stmts, schema, keyColumnRule{
		code:            v1review.PrimaryKeyTypeNotAllowed,
		typeOf:          lowerTypeName,
		attribute:       primaryKey,
		changeAttribute: primaryKey,
		constraint:      func(t ast.ConstraintType) bool { return t == ast.ConstrPrimaryKey },
		lower:           true,
		allowed:         func(columnType string) bool { return allowed[columnType] },
	})
}
