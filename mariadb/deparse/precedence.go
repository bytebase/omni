package deparse

import (
	ast "github.com/bytebase/omni/mariadb/ast"
	"strings"
)

// MariaDB's stored expressions omit redundant parentheses while preserving
// operator binding, including a right-hand expression with equal precedence.
// https://mariadb.com/docs/server/reference/sql-structure/operators/operator-precedence
func expressionPrecedence(expr ast.ExprNode) int {
	switch n := unwrapParen(expr).(type) {
	case *ast.BinaryExpr:
		if sub, ok := n.Right.(*ast.SubqueryExpr); ok && sub.Quantifier != "" {
			if n.Op == ast.BinOpEq && sub.Quantifier != "ALL" {
				return 7
			} // stored as IN
			if n.Op == ast.BinOpNe && sub.Quantifier == "ALL" {
				return 14
			} // stored as !(IN)
		}
		switch n.Op {
		case ast.BinOpAssign:
			return 0
		case ast.BinOpOr:
			return 1
		case ast.BinOpXor:
			return 2
		case ast.BinOpAnd:
			return 3
		case ast.BinOpEq, ast.BinOpNe, ast.BinOpLt, ast.BinOpGt, ast.BinOpLe, ast.BinOpGe, ast.BinOpNullSafeEq, ast.BinOpSoundsLike:
			return 5
		case ast.BinOpRegexp:
			return 7
		case ast.BinOpBitOr:
			return 8
		case ast.BinOpBitAnd:
			return 9
		case ast.BinOpShiftLeft, ast.BinOpShiftRight:
			return 10
		case ast.BinOpAdd, ast.BinOpSub:
			return 11
		case ast.BinOpMul, ast.BinOpDiv, ast.BinOpDivInt, ast.BinOpMod:
			return 12
		case ast.BinOpBitXor:
			return 13
		}
	case *ast.UnaryExpr:
		return 14
	case *ast.IsExpr:
		return 5
	case *ast.BetweenExpr:
		return 6
	case *ast.InExpr, *ast.LikeExpr:
		return 7
	}
	return 100
}

func writeOperand(b *strings.Builder, expr ast.ExprNode, minimum int) {
	parens := expressionPrecedence(expr) < minimum
	if parens {
		b.WriteByte('(')
	}
	writeExpr(b, expr)
	if parens {
		b.WriteByte(')')
	}
}
