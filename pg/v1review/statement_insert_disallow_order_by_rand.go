package v1review

import (
	"strings"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementInsertDisallowOrderByRand reports each top-level INSERT ...
// SELECT whose SELECT orders by a call to random or random_between, one
// InsertOrderByRand finding per statement. A function counts when any part
// of its qualified name is one of those, in any letter case. The ORDER BY
// of each branch of a set operation is judged, not one on the set
// operation as a whole, nor one in a subquery or WITH query. Findings
// anchor nothing: Bytebase's V1 advisor placed them on the statement's
// first line. Ported from Bytebase's V1 rule
// STATEMENT_INSERT_DISALLOW_ORDER_BY_RAND.
func StatementInsertDisallowOrderByRand(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		insert, ok := unwrapRaw(stmt).(*ast.InsertStmt)
		if !ok {
			continue
		}
		sel, ok := insert.SelectStmt.(*ast.SelectStmt)
		if !ok || sel == nil {
			continue
		}
		if ordersByRandom(sel) {
			findings = append(findings, v1review.Finding{Code: v1review.InsertOrderByRand})
		}
	}
	return findings
}

func ordersByRandom(sel *ast.SelectStmt) bool {
	if sel == nil {
		return false
	}
	if sel.Op != ast.SETOP_NONE {
		return ordersByRandom(sel.Larg) || ordersByRandom(sel.Rarg)
	}
	if sel.SortClause == nil {
		return false
	}
	for _, item := range sel.SortClause.Items {
		sortBy, ok := item.(*ast.SortBy)
		if !ok {
			continue
		}
		call, ok := sortBy.Node.(*ast.FuncCall)
		if !ok || call.Funcname == nil {
			continue
		}
		for _, part := range call.Funcname.Items {
			s, ok := part.(*ast.String)
			if !ok {
				continue
			}
			if name := strings.ToLower(s.Str); name == "random" || name == "random_between" {
				return true
			}
		}
	}
	return false
}
