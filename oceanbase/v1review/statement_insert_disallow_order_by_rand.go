package v1review

import (
	"strings"

	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementInsertDisallowOrderByRand reports each ORDER BY item that is a
// call named RAND (any case, any schema, any arguments; a call nested in an
// expression is not judged) in the SELECT a top-level INSERT ... SELECT
// inserts from. It descends into a parenthesized SELECT before judging the
// outer ORDER BY, then into both sides of a set operation; subqueries
// elsewhere are not judged. Findings anchor nothing (the V1 advisor took the
// statement's first content line). InsertOrderByRand is the only code.
// Ported from Bytebase's V1 rule STATEMENT_INSERT_DISALLOW_ORDER_BY_RAND.
func StatementInsertDisallowOrderByRand(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		ins, ok := stmt.(*ast.InsertStmt)
		if !ok || ins.Select == nil {
			continue
		}
		findings = appendOrderByRand(findings, ins.Select)
	}
	return findings
}

func appendOrderByRand(findings []v1review.Finding, sel *ast.SelectStmt) []v1review.Finding {
	if sel == nil {
		return findings
	}
	if sel.ParenSource != nil {
		findings = appendOrderByRand(findings, sel.ParenSource)
	}
	for _, item := range sel.OrderBy {
		if fc, ok := item.Expr.(*ast.FuncCallExpr); ok && strings.EqualFold(fc.Name, "rand") {
			findings = append(findings, v1review.Finding{Code: v1review.InsertOrderByRand})
		}
	}
	if sel.SetOp != ast.SetOpNone {
		findings = appendOrderByRand(findings, sel.Left)
		findings = appendOrderByRand(findings, sel.Right)
	}
	return findings
}
