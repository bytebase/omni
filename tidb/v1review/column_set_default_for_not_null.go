package v1review

import (
	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnSetDefaultForNotNull reports the columns declared NOT NULL or
// PRIMARY KEY that have no DEFAULT but need one, as ColumnRequireDefault
// decides need.
//
// The columns judged are those of top-level statements only: every column a
// CREATE TABLE defines, anchored at the column definition, except those its
// table-level PRIMARY KEY lists, matched exactly; and the columns ALTER
// TABLE's ADD, CHANGE and MODIFY COLUMN write, anchored at the whole
// statement. A statement naming no table is skipped. Name is the column and
// Table the table, both without a qualifier. Ported from Bytebase's V1 rule
// COLUMN_SET_DEFAULT_FOR_NOT_NULL.
func ColumnSetDefaultForNotNull(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		create, ok := stmt.(*ast.CreateTableStmt)
		if !ok {
			findings = append(findings, columnRule([]ast.Node{stmt}, v1review.ColumnNotNullNoDefault, notNullLacksDefault)...)
			continue
		}
		if create.Table == nil {
			continue
		}
		primary := make(map[string]bool)
		for _, name := range primaryKeyColumns(create) {
			primary[name] = true
		}
		for _, col := range create.Columns {
			if col != nil && !primary[col.Name] && notNullLacksDefault(col) {
				findings = append(findings, columnFinding(v1review.ColumnNotNullNoDefault, create.Table.Name, col.Name, col.Loc))
			}
		}
	}
	return findings
}

func notNullLacksDefault(col *ast.ColumnDef) bool {
	notNull := hasColumnConstraint(col, ast.ColConstrNotNull) || hasColumnConstraint(col, ast.ColConstrPrimaryKey)
	return notNull && lacksDefault(col)
}
