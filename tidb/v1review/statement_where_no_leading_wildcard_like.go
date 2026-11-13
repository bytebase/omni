package v1review

import (
	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementWhereNoLeadingWildcardLike reports each statement holding,
// anywhere in it, a LIKE or NOT LIKE whose pattern is a string literal
// starting with %, once per statement however many it holds. A pattern
// wrapped in COLLATE is judged by its literal; any other pattern that is not
// a literal is not judged. Findings anchor nothing (the V1 advisor placed
// them on the statement's first line). LeadingWildcardLike is the only code.
// Ported from Bytebase's V1 rule STATEMENT_WHERE_NO_LEADING_WILDCARD_LIKE.
func StatementWhereNoLeadingWildcardLike(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		found := false
		ast.Inspect(stmt, func(n ast.Node) bool {
			if found {
				return false
			}
			like, ok := n.(*ast.LikeExpr)
			if !ok {
				return true
			}
			pattern := like.Pattern
			for {
				c, ok := pattern.(*ast.CollateExpr)
				if !ok {
					break
				}
				pattern = c.Expr
			}
			if lit, ok := pattern.(*ast.StringLit); ok && len(lit.Value) > 0 && lit.Value[0] == '%' {
				found = true
				return false
			}
			return true
		})
		if found {
			findings = append(findings, v1review.Finding{Code: v1review.LeadingWildcardLike})
		}
	}
	return findings
}
