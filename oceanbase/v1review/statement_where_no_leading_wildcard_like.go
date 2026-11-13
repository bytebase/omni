package v1review

import (
	"strings"

	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementWhereNoLeadingWildcardLike reports each LIKE anywhere in a
// statement, routine bodies and select lists included, whose pattern is a
// string literal starting with %; NOT LIKE is judged too, and a pattern that
// is not a plain string literal is not. Each finding anchors the LIKE
// expression. LeadingWildcardLike is the only code. Ported from Bytebase's V1
// rule STATEMENT_WHERE_NO_LEADING_WILDCARD_LIKE.
func StatementWhereNoLeadingWildcardLike(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		ast.Inspect(stmt, func(n ast.Node) bool {
			like, ok := n.(*ast.LikeExpr)
			if !ok {
				return true
			}
			if str, ok := like.Pattern.(*ast.StringLit); ok && strings.HasPrefix(str.Value, "%") {
				start, end := anchor(like.Loc)
				findings = append(findings, v1review.Finding{Code: v1review.LeadingWildcardLike, Start: start, End: end})
			}
			return true
		})
	}
	return findings
}
