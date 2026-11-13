package v1review

import (
	"strings"

	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnRequireDefault judges the columns top-level statements define for
// having no DEFAULT: every column a CREATE TABLE defines except its primary
// key columns (inline or by table constraint, names compared as written),
// and in an ALTER TABLE the columns of ADD COLUMN in both its forms, MODIFY
// COLUMN and CHANGE COLUMN, primary key or not. Columns of type BLOB, TEXT,
// JSON, GEOMETRY or a spatial type (the tiny, medium and long variants and
// the multi- and collection types included, name in any case) are exempt;
// AUTO_INCREMENT and generated columns are not. A statement whose table has
// no name is not judged, and neither are statements nested in a routine
// body.
//
// Each finding names the column and its table, without the database
// qualifier, and anchors the column definition. Findings are in statement
// order, then column order. Ported from Bytebase's V1 rule
// COLUMN_REQUIRE_DEFAULT.
func ColumnRequireDefault(stmts []ast.Node) []v1review.Finding {
	return judgeDefaults(stmts, func(col *ast.ColumnDef) bool {
		return !hasDefault(col) && !isDefaultExemptType(col)
	}, v1review.ColumnNoDefault)
}

// judgeDefaults reports the columns violates holds for, read as V1's default
// rules read them: every column a CREATE TABLE defines except its primary
// key columns, and in an ALTER TABLE the modifiedColumns of each action.
func judgeDefaults(stmts []ast.Node, violates func(*ast.ColumnDef) bool, code v1review.Code) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			table := refName(n.Table)
			if table == "" {
				continue
			}
			pk := primaryKeyColumns(n)
			for _, col := range n.Columns {
				if pk[col.Name] {
					continue
				}
				if violates(col) {
					findings = append(findings, columnFinding(code, table, col))
				}
			}
		case *ast.AlterTableStmt:
			table := refName(n.Table)
			if table == "" {
				continue
			}
			for _, cmd := range n.Commands {
				for _, col := range modifiedColumns(cmd) {
					if violates(col) {
						findings = append(findings, columnFinding(code, table, col))
					}
				}
			}
		}
	}
	return findings
}

func isDefaultExemptType(col *ast.ColumnDef) bool {
	if col == nil || col.TypeName == nil {
		return false
	}
	switch strings.ToUpper(col.TypeName.Name) {
	case "BLOB", "TINYBLOB", "MEDIUMBLOB", "LONGBLOB",
		"TEXT", "TINYTEXT", "MEDIUMTEXT", "LONGTEXT",
		"JSON", "GEOMETRY", "POINT", "LINESTRING", "POLYGON",
		"MULTIPOINT", "MULTILINESTRING", "MULTIPOLYGON", "GEOMETRYCOLLECTION",
		"GEOMCOLLECTION":
		return true
	default:
		return false
	}
}
