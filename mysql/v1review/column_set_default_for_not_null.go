package v1review

import (
	"strings"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnSetDefaultForNotNull judges the columns top-level statements define
// for being NOT NULL, by a NOT NULL or PRIMARY KEY constraint of their own,
// yet having no DEFAULT: every column a CREATE TABLE defines except its
// primary key columns (inline or by table constraint, names compared as
// written), and in an ALTER TABLE the columns of ADD COLUMN in both its
// forms, MODIFY COLUMN and CHANGE COLUMN. Inline PRIMARY KEY and
// AUTO_INCREMENT columns are exempt, and so are columns of type BLOB, TEXT,
// JSON or GEOMETRY (the tiny, medium and long variants included, name in
// any case) or SERIAL; the spatial subtypes are not. A statement whose table
// has no name is not judged, and neither are statements nested in a routine
// body.
//
// Each finding names the column and its table, without the database
// qualifier, and anchors the column definition. Findings are in statement
// order, then column order. Ported from Bytebase's V1 rule
// COLUMN_SET_DEFAULT_FOR_NOT_NULL.
func ColumnSetDefaultForNotNull(stmts []ast.Node) []v1review.Finding {
	return judgeDefaults(stmts, func(col *ast.ColumnDef) bool {
		return !isNullable(col) && !hasDefault(col) && needsDefaultNotNull(col)
	}, v1review.ColumnNotNullNoDefault)
}

// needsDefaultNotNull reports whether a NOT NULL column without a DEFAULT
// must have one, as V1 decided it.
func needsDefaultNotNull(col *ast.ColumnDef) bool {
	if col == nil {
		return false
	}
	if isNullable(col) || hasDefault(col) || isPrimaryKey(col) || col.AutoIncrement {
		return false
	}
	if col.TypeName != nil {
		switch strings.ToUpper(col.TypeName.Name) {
		case "BLOB", "TINYBLOB", "MEDIUMBLOB", "LONGBLOB",
			"TEXT", "TINYTEXT", "MEDIUMTEXT", "LONGTEXT",
			"JSON", "GEOMETRY", "SERIAL":
			return false
		}
	}
	return true
}
