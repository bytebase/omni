package v1review

import (
	"strings"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementInsertDisallowOrderByRand reports each ORDER BY item that calls
// RAND (any case, any arguments) in the query of a top-level INSERT ...
// SELECT. It looks at the query's own ORDER BY, the inner query of a
// parenthesized query (whose outer ORDER BY is then judged too), and both
// sides of a set operation, recursively; subqueries and CTEs are not judged.
// Findings anchor nothing (the V1 advisor placed them on the statement's
// first line). InsertOrderByRand is the only code. Ported from Bytebase's
// V1 rule STATEMENT_INSERT_DISALLOW_ORDER_BY_RAND.
func StatementInsertDisallowOrderByRand(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		ins, ok := stmt.(*ast.InsertStmt)
		if !ok || ins.Select == nil {
			continue
		}
		findings = orderByRand(ins.Select, findings)
	}
	return findings
}

func orderByRand(sel *ast.SelectStmt, findings []v1review.Finding) []v1review.Finding {
	if sel == nil {
		return findings
	}
	if sel.ParenSource != nil {
		findings = orderByRand(sel.ParenSource, findings)
	}
	for _, item := range sel.OrderBy {
		if fc, ok := item.Expr.(*ast.FuncCallExpr); ok && strings.EqualFold(fc.Name, "rand") {
			findings = append(findings, v1review.Finding{Code: v1review.InsertOrderByRand})
		}
	}
	if sel.SetOp != ast.SetOpNone {
		findings = orderByRand(sel.Left, findings)
		findings = orderByRand(sel.Right, findings)
	}
	return findings
}
