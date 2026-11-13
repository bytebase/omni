package v1review

import (
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementWhereNoEqualNull reports each = or != (<>) comparison whose right
// operand is a NULL literal inside the WHERE of any SELECT in a top-level
// statement, subqueries and CTEs included; NULL on the left and <=> are not
// judged. A comparison in a subquery's WHERE that itself sits in an outer
// SELECT's WHERE is reported twice, once per enclosing WHERE walked, as the
// V1 advisor did. Each finding anchors the comparison. EqualNull is the
// only code. Ported from Bytebase's V1 rule STATEMENT_WHERE_NO_EQUAL_NULL.
func StatementWhereNoEqualNull(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		ast.Inspect(stmt, func(n ast.Node) bool {
			sel, ok := n.(*ast.SelectStmt)
			if !ok || sel.Where == nil {
				return true
			}
			ast.Inspect(sel.Where, func(wn ast.Node) bool {
				bin, ok := wn.(*ast.BinaryExpr)
				if !ok || (bin.Op != ast.BinOpEq && bin.Op != ast.BinOpNe) {
					return true
				}
				if _, ok := bin.Right.(*ast.NullLit); ok {
					start, end := anchor(bin.Loc)
					findings = append(findings, v1review.Finding{Code: v1review.EqualNull, Start: start, End: end})
				}
				return true
			})
			return true
		})
	}
	return findings
}
