package v1review

import (
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementWhereRequireSelect reports each SELECT that has a FROM clause and
// no WHERE clause, one WhereMissing finding per SELECT, visiting SELECTs as
// selectWalk describes, so a subquery's finding comes before its enclosing
// SELECT's. A SELECT without FROM, such as SELECT 1, is not judged. Name is
// empty: Bytebase's V1 advisor printed no kind. A finding on the top-level
// SELECT statement anchors nothing, and one on any other SELECT anchors that
// SELECT: the V1 advisor printed the statement's text or that SELECT's
// (widened by the parentheses around it), and placed every finding on the
// statement's first line. Ported from Bytebase's V1 rule
// STATEMENT_WHERE_REQUIRE_SELECT.
func StatementWhereRequireSelect(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	w := &selectWalk{visit: func(sel *ast.SelectStmt, top bool) {
		if sel.FromClause == nil || len(sel.FromClause.Items) == 0 || sel.WhereClause != nil {
			return
		}
		findings = append(findings, selectFinding(v1review.WhereMissing, sel, top))
	}}
	for _, stmt := range stmts {
		w.statement(stmt)
	}
	return findings
}
