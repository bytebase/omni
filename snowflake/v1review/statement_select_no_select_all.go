package v1review

import (
	"slices"

	"github.com/bytebase/omni/snowflake/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementSelectNoSelectAll reports each * or qualifier.* select-list item
// of every SELECT anywhere in the statements: subqueries, CTE bodies,
// INSERT ... SELECT, and statements nested in scripting blocks included. A
// star inside an expression, such as COUNT(*), is not a select-list item and
// is not judged. A finding anchors the * itself, the last byte of the star
// expression, or the select item when it has no star expression; findings
// are in anchor order, not walk order. SelectAll is the only code. Ported
// from Bytebase's V1 rule STATEMENT_SELECT_NO_SELECT_ALL.
func StatementSelectNoSelectAll(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		ast.Inspect(stmt, func(n ast.Node) bool {
			sel, ok := n.(*ast.SelectStmt)
			if !ok {
				return true
			}
			for _, target := range sel.Targets {
				if target == nil || !target.Star {
					continue
				}
				start, end := stmtAnchor(target.Loc)
				if star, ok := target.Expr.(*ast.StarExpr); ok && star.Loc.IsValid() {
					start, end = stmtAnchor(ast.Loc{Start: star.Loc.End - 1, End: star.Loc.End})
				}
				findings = append(findings, v1review.Finding{Code: v1review.SelectAll, Start: start, End: end})
			}
			return true
		})
	}
	slices.SortStableFunc(findings, func(a, b v1review.Finding) int { return a.Start - b.Start })
	return findings
}
