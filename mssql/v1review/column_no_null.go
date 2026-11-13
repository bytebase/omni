package v1review

import (
	"strings"

	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnNoNull judges each top-level statement on its own for nullable
// columns: every column a CREATE TABLE defines, unless the table has no
// name, and the column of each ADD COLUMN and ALTER COLUMN action of an
// ALTER TABLE. A column is not nullable when it says NOT NULL or carries a
// PRIMARY KEY or NOT NULL constraint, or, in a CREATE TABLE, when a PRIMARY
// KEY table constraint names it, compared lower-cased, so an ALTER COLUMN
// that does not say NOT NULL is reported. Statements nested in a routine
// body are not judged.
//
// Each finding has code ColumnNullable, Name the column lower-cased, and
// anchors the column definition. Findings come in statement order, then
// column or action order. Ported from Bytebase's V1 rule COLUMN_NO_NULL.
func ColumnNoNull(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	check := func(col *ast.ColumnDef, pkColumns map[string]bool) {
		if pkColumns[strings.ToLower(col.Name)] || columnNotNull(col) {
			return
		}
		start, end := span(col.Loc)
		findings = append(findings, v1review.Finding{Code: v1review.ColumnNullable, Name: strings.ToLower(col.Name), Start: start, End: end})
	}
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.Name == nil {
				continue
			}
			pkColumns := make(map[string]bool)
			for _, cd := range listOf[*ast.ConstraintDef](n.Constraints) {
				if cd.Type != ast.ConstraintPrimaryKey {
					continue
				}
				for _, ic := range listOf[*ast.IndexColumn](cd.Columns) {
					pkColumns[strings.ToLower(ic.Name)] = true
				}
			}
			for _, col := range listOf[*ast.ColumnDef](n.Columns) {
				check(col, pkColumns)
			}
		case *ast.AlterTableStmt:
			for _, action := range listOf[*ast.AlterTableAction](n.Actions) {
				if action.Type != ast.ATAddColumn && action.Type != ast.ATAlterColumn {
					continue
				}
				if action.Column != nil {
					check(action.Column, nil)
				}
			}
		}
	}
	return findings
}

// columnNotNull reports whether a column says NOT NULL or carries a PRIMARY
// KEY or NOT NULL constraint.
func columnNotNull(col *ast.ColumnDef) bool {
	for _, cd := range listOf[*ast.ConstraintDef](col.Constraints) {
		if cd.Type == ast.ConstraintPrimaryKey || cd.Type == ast.ConstraintNotNull {
			return true
		}
	}
	return col.Nullable != nil && col.Nullable.NotNull
}
