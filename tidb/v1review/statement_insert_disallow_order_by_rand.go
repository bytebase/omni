package v1review

import (
	"strings"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementInsertDisallowOrderByRand reports each top-level INSERT ... SELECT
// whose outermost ORDER BY holds a RAND() call, in any letter case, as an
// item of its own, once per statement however many such items it has. An
// ORDER BY in a subquery or a parenthesized set-operation arm is not judged,
// nor is RAND() inside a larger expression. Findings anchor nothing (the V1
// advisor placed them on the statement's first line). InsertOrderByRand is
// the only code. Ported from Bytebase's V1 rule
// STATEMENT_INSERT_DISALLOW_ORDER_BY_RAND.
func StatementInsertDisallowOrderByRand(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		n, ok := stmt.(*ast.InsertStmt)
		if !ok || n.Select == nil {
			continue
		}
		for _, item := range n.Select.OrderBy {
			if item == nil {
				continue
			}
			if fc, ok := item.Expr.(*ast.FuncCallExpr); ok && strings.EqualFold(fc.Name, "RAND") {
				findings = append(findings, v1review.Finding{Code: v1review.InsertOrderByRand})
				break
			}
		}
	}
	return findings
}
