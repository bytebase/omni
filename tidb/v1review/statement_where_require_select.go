package v1review

import (
	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementWhereRequireSelect reports, in walk order, every SELECT with a
// FROM clause and no WHERE clause anywhere in the statements: subqueries,
// set-operation arms, CTE bodies, INSERT ... SELECT and EXPLAIN included. A
// SELECT without FROM, such as SELECT 1, is not judged. Findings anchor
// nothing (the V1 advisor placed them on the statement's first line).
// WhereMissing is the only code; Name is not set. Ported from Bytebase's V1
// rule STATEMENT_WHERE_REQUIRE_SELECT.
func StatementWhereRequireSelect(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		ast.Inspect(stmt, func(n ast.Node) bool {
			if sel, ok := n.(*ast.SelectStmt); ok && sel.Where == nil && len(sel.From) > 0 {
				findings = append(findings, v1review.Finding{Code: v1review.WhereMissing})
			}
			return true
		})
	}
	return findings
}
