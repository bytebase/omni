package v1review

import (
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
	"strconv"
)

// StatementWhereMaximumLogicalOperatorCount reports WHERE clauses of
// top-level statements that hold more than maximum operands of one logical
// operator. It judges the WHERE of a SELECT (each side of a set operation,
// and the inner query of a parenthesized one, whose outer clauses are then
// skipped), of an INSERT ... SELECT's query, and of an UPDATE or DELETE.
// For each WHERE it counts the operands of the OR chain at its root, only
// down the left side (a OR (b OR c) is two), and anchors the deepest OR of
// that chain; a WHERE whose root is not an OR counts zero and anchors
// nothing, so a negative maximum reports every WHERE. It then reports every
// IN list longer than maximum anywhere in the WHERE, subqueries included,
// but not an IN nested inside an IN, each anchored at its IN; an IN over a
// subquery counts zero. Name is "OR" or "IN" and Value the operands.
// LogicalOperatorsExceeded is the only code. Ported from Bytebase's V1 rule
// STATEMENT_WHERE_MAXIMUM_LOGICAL_OPERATOR_COUNT.
func StatementWhereMaximumLogicalOperatorCount(stmts []ast.Node, maximum int) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.SelectStmt:
			findings = logicalOperatorsInSelect(n, maximum, findings)
		case *ast.UpdateStmt:
			findings = logicalOperatorsInWhere(n.Where, maximum, findings)
		case *ast.DeleteStmt:
			findings = logicalOperatorsInWhere(n.Where, maximum, findings)
		case *ast.InsertStmt:
			findings = logicalOperatorsInSelect(n.Select, maximum, findings)
		}
	}
	return findings
}

func logicalOperatorsInSelect(sel *ast.SelectStmt, maximum int, findings []v1review.Finding) []v1review.Finding {
	if sel == nil {
		return findings
	}
	if sel.ParenSource != nil {
		return logicalOperatorsInSelect(sel.ParenSource, maximum, findings)
	}
	if sel.SetOp != ast.SetOpNone {
		findings = logicalOperatorsInSelect(sel.Left, maximum, findings)
		return logicalOperatorsInSelect(sel.Right, maximum, findings)
	}
	return logicalOperatorsInWhere(sel.Where, maximum, findings)
}

func logicalOperatorsInWhere(where ast.ExprNode, maximum int, findings []v1review.Finding) []v1review.Finding {
	if where == nil {
		return findings
	}
	count, start, end := 0, 0, 0
	if or, ok := where.(*ast.BinaryExpr); ok && or.Op == ast.BinOpOr {
		count = 1
		for {
			count++
			left, ok := or.Left.(*ast.BinaryExpr)
			if !ok || left.Op != ast.BinOpOr {
				break
			}
			or = left
		}
		start, end = anchor(or.Loc)
	}
	if count > maximum {
		findings = append(findings, v1review.Finding{Code: v1review.LogicalOperatorsExceeded, Name: "OR", Value: strconv.Itoa(count), Start: start, End: end})
	}
	ast.Inspect(where, func(n ast.Node) bool {
		in, ok := n.(*ast.InExpr)
		if !ok {
			return true
		}
		if len(in.List) > maximum {
			start, end := anchor(in.Loc)
			findings = append(findings, v1review.Finding{Code: v1review.LogicalOperatorsExceeded, Name: "IN", Value: strconv.Itoa(len(in.List)), Start: start, End: end})
		}
		return false
	})
	return findings
}
