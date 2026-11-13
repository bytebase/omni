package v1review

import (
	"strings"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementWhereNoLeadingWildcardLike reports a top-level statement that
// holds, anywhere in its tree, a LIKE or ILIKE (NOT LIKE and NOT ILIKE
// included) whose pattern is a string constant starting with % or _: at
// most one LeadingWildcardLike finding per statement, wherever the
// expression is, not only in WHERE. SIMILAR TO, a pattern that is not a
// plain constant, and a routine body are not judged. Findings anchor
// nothing: Bytebase's V1 advisor placed them on the statement's first
// line. Ported from Bytebase's V1 rule
// STATEMENT_WHERE_NO_LEADING_WILDCARD_LIKE.
func StatementWhereNoLeadingWildcardLike(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		if hasLeadingWildcardLike(unwrapRaw(stmt)) {
			findings = append(findings, v1review.Finding{Code: v1review.LeadingWildcardLike})
		}
	}
	return findings
}

func hasLeadingWildcardLike(node ast.Node) bool {
	found := false
	ast.Inspect(node, func(n ast.Node) bool {
		if found {
			return false
		}
		expr, ok := n.(*ast.A_Expr)
		if !ok || (expr.Kind != ast.AEXPR_LIKE && expr.Kind != ast.AEXPR_ILIKE) {
			return true
		}
		if c, ok := expr.Rexpr.(*ast.A_Const); ok {
			if s, ok := c.Val.(*ast.String); ok && (strings.HasPrefix(s.Str, "%") || strings.HasPrefix(s.Str, "_")) {
				found = true
				return false
			}
		}
		return true
	})
	return found
}
