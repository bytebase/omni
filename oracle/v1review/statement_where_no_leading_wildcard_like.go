package v1review

import (
	"strings"

	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementWhereNoLeadingWildcardLike reports, in walk order, each LIKE
// anywhere in the statements, in any clause and in PL/SQL bodies, whose
// pattern is a string literal starting with %. NOT LIKE and the LIKEC,
// LIKE2 and LIKE4 variants count; a pattern that is not a literal is not
// judged. A finding anchors the LIKE expression. LeadingWildcardLike is the
// only code. Ported from Bytebase's V1 rule
// STATEMENT_WHERE_NO_LEADING_WILDCARD_LIKE.
func StatementWhereNoLeadingWildcardLike(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		ast.Inspect(stmt, func(n ast.Node) bool {
			like, ok := n.(*ast.LikeExpr)
			if !ok {
				return true
			}
			if pattern, ok := like.Pattern.(*ast.StringLiteral); ok && strings.HasPrefix(pattern.Val, "%") {
				start, end := span(like.Loc)
				findings = append(findings, v1review.Finding{Code: v1review.LeadingWildcardLike, Start: start, End: end})
			}
			return true
		})
	}
	return findings
}
