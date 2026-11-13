package v1review

import (
	"strings"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnRequireDefault reports the columns that have no DEFAULT but need
// one. A column needs none when it is AUTO_INCREMENT, generated, declared
// PRIMARY KEY itself, has no type, or is of a BLOB, TEXT, JSON or spatial
// type, the type name compared without regard to case.
//
// The columns judged are those of top-level statements only: every column a
// CREATE TABLE defines, anchored at the column definition, except those its
// table-level PRIMARY KEY lists, matched without regard to case; and the
// columns ALTER TABLE's ADD, CHANGE and MODIFY COLUMN write, anchored at the
// whole statement. A statement naming no table is skipped. Name is the
// column and Table the table, both without a qualifier. Ported from
// Bytebase's V1 rule COLUMN_REQUIRE_DEFAULT.
func ColumnRequireDefault(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		create, ok := stmt.(*ast.CreateTableStmt)
		if !ok {
			findings = append(findings, columnRule([]ast.Node{stmt}, v1review.ColumnNoDefault, lacksDefault)...)
			continue
		}
		if create.Table == nil {
			continue
		}
		primary := make(map[string]bool)
		for _, name := range primaryKeyColumns(create) {
			primary[strings.ToLower(name)] = true
		}
		for _, col := range create.Columns {
			if col != nil && !primary[strings.ToLower(col.Name)] && lacksDefault(col) {
				findings = append(findings, columnFinding(v1review.ColumnNoDefault, create.Table.Name, col.Name, col.Loc))
			}
		}
	}
	return findings
}

// lacksDefault reports whether a column has no DEFAULT but needs one.
func lacksDefault(col *ast.ColumnDef) bool {
	return col.DefaultValue == nil && needDefault(col)
}

// needDefault reports whether a column needs a DEFAULT, as
// ColumnRequireDefault describes.
func needDefault(col *ast.ColumnDef) bool {
	if col.AutoIncrement || col.Generated != nil || hasColumnConstraint(col, ast.ColConstrPrimaryKey) || col.TypeName == nil {
		return false
	}
	switch strings.ToUpper(col.TypeName.Name) {
	case "BLOB", "TINYBLOB", "MEDIUMBLOB", "LONGBLOB",
		"TEXT", "TINYTEXT", "MEDIUMTEXT", "LONGTEXT",
		"JSON",
		"GEOMETRY", "POINT", "LINESTRING", "POLYGON",
		"MULTIPOINT", "MULTILINESTRING", "MULTIPOLYGON", "GEOMETRYCOLLECTION":
		return false
	default:
		return true
	}
}

// primaryKeyColumns returns the columns a CREATE TABLE's table-level PRIMARY
// KEY constraints list, as written.
func primaryKeyColumns(create *ast.CreateTableStmt) []string {
	var columns []string
	for _, c := range create.Constraints {
		if c != nil && c.Type == ast.ConstrPrimaryKey {
			columns = append(columns, c.Columns...)
		}
	}
	return columns
}
