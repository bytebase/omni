package v1review

import (
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementDisallowAddColumnWithDefault reports a top-level ALTER TABLE
// that adds a column with a DEFAULT, at most one ColumnWithDefault finding
// per statement however many such columns it adds. Any DEFAULT counts,
// constant or volatile. Findings anchor nothing: Bytebase's V1 advisor
// placed them on the statement's first line. Ported from Bytebase's V1 rule
// STATEMENT_DISALLOW_ADD_COLUMN_WITH_DEFAULT.
func StatementDisallowAddColumnWithDefault(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		alter, ok := unwrapRaw(stmt).(*ast.AlterTableStmt)
		if !ok {
			continue
		}
		if addsColumnWithDefault(alter) {
			findings = append(findings, v1review.Finding{Code: v1review.ColumnWithDefault})
		}
	}
	return findings
}

func addsColumnWithDefault(alter *ast.AlterTableStmt) bool {
	for _, cmd := range alterTableCmds(alter) {
		if ast.AlterTableType(cmd.Subtype) != ast.AT_AddColumn {
			continue
		}
		col, ok := cmd.Def.(*ast.ColumnDef)
		if !ok {
			continue
		}
		if col.RawDefault != nil {
			return true
		}
		for _, c := range columnConstraints(col) {
			if c.Contype == ast.CONSTR_DEFAULT {
				return true
			}
		}
	}
	return false
}
