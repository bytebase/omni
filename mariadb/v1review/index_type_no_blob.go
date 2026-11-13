package v1review

import (
	"strings"

	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/v1review"
)

// IndexTypeNoBlob reports the key columns of a script whose type is BLOB,
// TINYBLOB, MEDIUMBLOB or LONGBLOB: a column defined with the PRIMARY KEY or
// UNIQUE attribute; each column a table-level PRIMARY KEY, UNIQUE, INDEX or
// FOREIGN KEY definition lists, in a CREATE TABLE or an ALTER TABLE; and each
// column of a CREATE INDEX that is neither FULLTEXT nor SPATIAL, CREATE VECTOR
// INDEX included. The TEXT types are not reported, nor are the columns of a
// FULLTEXT, SPATIAL or VECTOR definition. How the script is walked, how a listed column's type is found,
// and when the function returns v1review.ErrSchemaMissing are
// judgeKeyColumns'.
//
// A defined column's type is its type's name in lower case, as the parser
// names it, so LONG VARBINARY is `mediumblob`. A type read from
// schema is the snapshot's whole spelling in lower case, so `blob(100)` is
// not reported.
//
// Name is the column, Table the table's name without its database
// qualifier, and Value the type as judged. Ported from Bytebase's V1 rule
// INDEX_TYPE_NO_BLOB.
func IndexTypeNoBlob(stmts []ast.Node, schema v1review.Schema) ([]v1review.Finding, error) {
	return judgeKeyColumns(stmts, schema, keyColumnRule{
		code:   v1review.IndexColumnTypeNotAllowed,
		typeOf: lowerTypeName,
		attribute: func(t ast.ColumnConstraintType) bool {
			return t == ast.ColConstrPrimaryKey || t == ast.ColConstrUnique
		},
		constraint: func(t ast.ConstraintType) bool {
			switch t {
			case ast.ConstrPrimaryKey, ast.ConstrUnique, ast.ConstrIndex, ast.ConstrForeignKey:
				return true
			default:
				return false
			}
		},
		createIndex: true,
		lower:       true,
		allowed: func(columnType string) bool {
			switch strings.ToLower(columnType) {
			case "blob", "tinyblob", "mediumblob", "longblob":
				return false
			default:
				return true
			}
		},
	})
}
