package v1review

import (
	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementWhereRequireSelect reports, in walk order, every SELECT with a
// FROM clause and no WHERE clause anywhere in the statements: subqueries,
// CTE bodies, INSERT ... SELECT, and SELECTs nested in procedure bodies and
// control-flow blocks included. A SELECT without FROM, such as SELECT 1, is
// not judged. A finding anchors the SELECT. WhereMissing is the only code;
// Name is not set. Ported from Bytebase's V1 rule
// STATEMENT_WHERE_REQUIRE_SELECT.
func StatementWhereRequireSelect(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		ast.Inspect(stmt, func(n ast.Node) bool {
			sel, ok := n.(*ast.SelectStmt)
			if !ok {
				return true
			}
			if sel.FromClause != nil && sel.WhereClause == nil {
				start, end := stmtAnchor(sel.Loc)
				findings = append(findings, v1review.Finding{Code: v1review.WhereMissing, Start: start, End: end})
			}
			return true
		})
	}
	return findings
}
