package parser

import (
	"testing"
)

// TestParenSublinkTypedNilCallers covers the typed-nil propagation
// fix for the four parseSelectStmtForExpr callers:
//
//   - expr.go parseInExpr       (IN (subquery))
//   - expr.go parseSubqueryOp   (ANY/ALL/SOME (subquery))
//   - expr.go parseParenExprOrRow (SELECT ... (subquery))
//   - expr.go parseExistsExpr   (EXISTS (subquery))
//
// Each caller previously built a SubLink with Subselect boxed from a
// typed-nil *SelectStmt, which is NOT == nil when stored in a
// nodes.Node interface. The fix, mirrored from §1.4 parseArrayCExpr,
// type-asserts the returned Node and rejects when the underlying
// pointer is nil.
//
// One scenario per callsite. Each is a minimal malformed input that
// would reach the parseSelectStmtForExpr call but produce a typed-nil
// *SelectStmt (the paren content is not a valid select_no_parens
// lead — e.g. empty parens, or a bare non-SELECT token).
func TestParenSublinkTypedNilCallers(t *testing.T) {
	cases := []struct {
		sql  string
		site string
	}{
		// parseExistsExpr — no isSelectStart gate; empty parens reach
		// parseSelectStmtForExpr directly and return typed-nil.
		{`SELECT EXISTS()`, "parseExistsExpr"},
		// parseInExpr subquery branch — `((SELECT ...))` triggers the
		// lookaheadParenContentIsSubquery scan; with a malformed inner (WITH cte
		// AS (SELECT 1)) the lookahead still classifies as subquery and
		// parseSelectStmtForExpr returns typed-nil on no trailing
		// select_clause.
		{`SELECT * FROM t WHERE x IN (WITH cte AS (SELECT 1))`, "parseInExpr"},
		// parseSubqueryOp — ANY with a malformed WITH-only subquery.
		{`SELECT 1 = ANY (WITH cte AS (SELECT 1))`, "parseSubqueryOp"},
		// parseParenExprOrRow — same malformed WITH-only content at
		// the start of a parenthesized expression position.
		{`SELECT (WITH cte AS (SELECT 1))`, "parseParenExprOrRow"},
	}
	for _, tc := range cases {
		t.Run(tc.site+"/"+tc.sql, func(t *testing.T) {
			_, err := Parse(tc.sql)
			if err == nil {
				t.Fatalf("expected parse error (%s malformed subquery), got nil", tc.site)
			}
		})
	}
}
