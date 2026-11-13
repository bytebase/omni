package v1review

import (
	"strings"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementWhereNoLeadingWildcardLike reports each LIKE (or NOT LIKE) whose
// pattern is a string literal starting with %, anywhere in a top-level
// statement, not only in a WHERE clause. Each finding anchors the LIKE
// expression. LeadingWildcardLike is the only code. Ported from Bytebase's
// V1 rule STATEMENT_WHERE_NO_LEADING_WILDCARD_LIKE.
func StatementWhereNoLeadingWildcardLike(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		ast.Inspect(stmt, func(n ast.Node) bool {
			like, ok := n.(*ast.LikeExpr)
			if !ok {
				return true
			}
			if pattern, ok := like.Pattern.(*ast.StringLit); ok && strings.HasPrefix(pattern.Value, "%") {
				start, end := anchor(like.Loc)
				findings = append(findings, v1review.Finding{Code: v1review.LeadingWildcardLike, Start: start, End: end})
			}
			return true
		})
	}
	return findings
}
