package v1review

import (
	"strings"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// IndexTypeNoBlob reports the key columns of a script whose type is BLOB,
// TINYBLOB, MEDIUMBLOB, LONGBLOB, TEXT, TINYTEXT, MEDIUMTEXT or LONGTEXT: a
// column a CREATE TABLE or an ADD COLUMN defines with the UNIQUE attribute,
// but not with PRIMARY KEY; a column a MODIFY COLUMN or a CHANGE COLUMN
// defines with either; each column a table-level PRIMARY KEY, UNIQUE or
// INDEX definition lists, in a CREATE TABLE or an ALTER TABLE; and each
// column of a CREATE INDEX of any kind, FULLTEXT and SPATIAL included. The
// columns of a FOREIGN KEY, FULLTEXT, SPATIAL or VECTOR definition are not
// reported. How the script is walked, how a listed column's type is found,
// and when the function returns v1review.ErrSchemaMissing are
// judgeKeyColumns'.
//
// A defined column's type is its type's name in lower case; a column with
// no type is the empty type. A type read from schema is the snapshot's
// whole spelling in lower case, so `blob(100)` is not reported.
//
// Name is the column, Table the table's name without its database
// qualifier, and Value the type as judged. Ported from Bytebase's V1 rule
// INDEX_TYPE_NO_BLOB.
func IndexTypeNoBlob(stmts []ast.Node, schema v1review.Schema) ([]v1review.Finding, error) {
	return judgeKeyColumns(stmts, schema, keyColumnRule{
		code:      v1review.IndexColumnTypeNotAllowed,
		typeOf:    lowerTypeName,
		attribute: func(t ast.ColumnConstraintType) bool { return t == ast.ColConstrUnique },
		changeAttribute: func(t ast.ColumnConstraintType) bool {
			return t == ast.ColConstrPrimaryKey || t == ast.ColConstrUnique
		},
		constraint: func(t ast.ConstraintType) bool {
			switch t {
			case ast.ConstrPrimaryKey, ast.ConstrUnique, ast.ConstrIndex:
				return true
			default:
				return false
			}
		},
		createIndex: true,
		lower:       true,
		allowed: func(columnType string) bool {
			switch strings.ToLower(columnType) {
			case "blob", "tinyblob", "mediumblob", "longblob", "text", "tinytext", "mediumtext", "longtext":
				return false
			default:
				return true
			}
		},
	})
}
