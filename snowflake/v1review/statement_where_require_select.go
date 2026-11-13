package v1review

import (
	"slices"

	"github.com/bytebase/omni/snowflake/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementWhereRequireSelect reports every SELECT with a FROM clause and no
// WHERE clause anywhere in the statements: CTE bodies, derived tables,
// subqueries, MERGE sources, and statements nested in scripting blocks
// included. A SELECT without FROM, such as SELECT 1, is not judged. A finding
// anchors the SELECT; for a WITH ... SELECT, whose range starts at its first
// CTE, it anchors the empty range at the end of the last CTE, because V1
// took the line of the body's SELECT keyword, the first token after it,
// which the AST does not record. Findings are in anchor order, not walk
// order. WhereMissing is the only code; Name is not set. Ported from
// Bytebase's V1 rule STATEMENT_WHERE_REQUIRE_SELECT.
func StatementWhereRequireSelect(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		ast.Inspect(stmt, func(n ast.Node) bool {
			sel, ok := n.(*ast.SelectStmt)
			if !ok || sel.Where != nil || sel.From == nil {
				return true
			}
			start, end := stmtAnchor(sel.Loc)
			if k := len(sel.With); k > 0 && sel.With[k-1] != nil && sel.With[k-1].Loc.End >= 0 {
				start, end = sel.With[k-1].Loc.End, sel.With[k-1].Loc.End
			}
			findings = append(findings, v1review.Finding{Code: v1review.WhereMissing, Start: start, End: end})
			return true
		})
	}
	slices.SortStableFunc(findings, func(a, b v1review.Finding) int { return a.Start - b.Start })
	return findings
}
