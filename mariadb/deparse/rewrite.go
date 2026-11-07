// Package deparse — rewrite.go implements AST rewrites applied before deparsing.
// These rewrites match MySQL 8.0's resolver behavior (SHOW CREATE VIEW output).
package deparse

import (
	ast "github.com/bytebase/omni/mariadb/ast"
)

// RewriteExpr applies MySQL 8.0 resolver rewrites to an expression AST.
// Currently implements NOT folding: NOT(comparison) → inverted comparison.
// The rewrite is recursive — children are processed first (bottom-up).
func RewriteExpr(node ast.ExprNode) ast.ExprNode {
	if node == nil {
		return nil
	}
	return rewriteExpr(node)
}

func rewriteExpr(node ast.ExprNode) ast.ExprNode {
	switch n := node.(type) {
	case *ast.UnaryExpr:
		// First rewrite the operand recursively
		n.Operand = rewriteExpr(n.Operand)
		if n.Op == ast.UnaryNot {
			return rewriteNot(n.Operand)
		}
		return n

	case *ast.BinaryExpr:
		n.Left = rewriteExpr(n.Left)
		n.Right = rewriteExpr(n.Right)
		// MySQL 8.0 rewrites SOUNDS LIKE to soundex(left) = soundex(right)
		if n.Op == ast.BinOpSoundsLike {
			return &ast.ParenExpr{
				Expr: &ast.BinaryExpr{
					Op: ast.BinOpEq,
					Left: &ast.FuncCallExpr{
						Name:      "soundex",
						HasParens: true,
						Args:      []ast.ExprNode{n.Left},
					},
					Right: &ast.FuncCallExpr{
						Name:      "soundex",
						HasParens: true,
						Args:      []ast.ExprNode{n.Right},
					},
				},
			}
		}
		// Boolean context wrapping: AND/OR/XOR operands that are not boolean
		// expressions get wrapped in (0 <> expr).
		if n.Op == ast.BinOpAnd || n.Op == ast.BinOpOr {
			if _, ok := unwrapParen(n.Left).(*ast.ColumnRef); ok {
				n.Left = wrapBooleanContext(n.Left)
			}
			if _, ok := unwrapParen(n.Right).(*ast.ColumnRef); ok {
				n.Right = wrapBooleanContext(n.Right)
			}
		}
		return n

	case *ast.ParenExpr:
		n.Expr = rewriteExpr(n.Expr)
		return n

	case *ast.InExpr:
		n.Expr = rewriteExpr(n.Expr)
		for i, item := range n.List {
			n.List[i] = rewriteExpr(item)
		}
		return n

	case *ast.BetweenExpr:
		n.Expr = rewriteExpr(n.Expr)
		n.Low = rewriteExpr(n.Low)
		n.High = rewriteExpr(n.High)
		return n

	case *ast.LikeExpr:
		n.Expr = rewriteExpr(n.Expr)
		n.Pattern = rewriteExpr(n.Pattern)
		if n.Escape != nil {
			n.Escape = rewriteExpr(n.Escape)
		}
		return n

	case *ast.IsExpr:
		n.Expr = rewriteExpr(n.Expr)
		return n

	case *ast.CaseExpr:
		if n.Operand != nil {
			n.Operand = rewriteExpr(n.Operand)
		}
		for _, w := range n.Whens {
			w.Cond = rewriteExpr(w.Cond)
			w.Result = rewriteExpr(w.Result)
		}
		if n.Default != nil {
			n.Default = rewriteExpr(n.Default)
		}
		return n

	case *ast.FuncCallExpr:
		for i, arg := range n.Args {
			n.Args[i] = rewriteExpr(arg)
		}
		return n

	case *ast.PositionExpr:
		n.Substr = rewriteExpr(n.Substr)
		n.Str = rewriteExpr(n.Str)
		return n

	case *ast.CastExpr:
		n.Expr = rewriteExpr(n.Expr)
		return n

	case *ast.ConvertExpr:
		n.Expr = rewriteExpr(n.Expr)
		return n

	case *ast.CollateExpr:
		n.Expr = rewriteExpr(n.Expr)
		return n

	default:
		// Leaf nodes (literals, column refs, etc.) — no rewriting needed
		return node
	}
}

// invertOp maps comparison operators to their NOT-inverted counterparts.
var invertOp = map[ast.BinaryOp]ast.BinaryOp{
	ast.BinOpGt: ast.BinOpLe,
	ast.BinOpLt: ast.BinOpGe,
	ast.BinOpGe: ast.BinOpLt,
	ast.BinOpLe: ast.BinOpGt,
	ast.BinOpEq: ast.BinOpNe,
	ast.BinOpNe: ast.BinOpEq,
}

// isComparisonOp returns true if op is a comparison operator that can be inverted.
func isComparisonOp(op ast.BinaryOp) bool {
	_, ok := invertOp[op]
	return ok
}

// unwrapParen strips ParenExpr wrappers to get the inner expression.
func unwrapParen(node ast.ExprNode) ast.ExprNode {
	for {
		p, ok := node.(*ast.ParenExpr)
		if !ok {
			return node
		}
		node = p.Expr
	}
}

// rewriteNot applies NOT folding to the operand of a NOT expression.
// MySQL 8.0's resolver:
//   - NOT (comparison) → inverted comparison (e.g., NOT(a > 0) → (a <= 0))
//   - NOT (LIKE) → not((expr like pattern)) — wraps in not(), doesn't fold
//   - NOT (non-boolean) → (0 = expr) — e.g., NOT(a+1) → (0 = (a+1)), NOT(col) → (0 = col)
func rewriteNot(operand ast.ExprNode) ast.ExprNode {
	// NOT of a quantified comparison inverts the operator and swaps the quantifier:
	// `NOT a > ALL (q)` is `a <= ANY (q)`. Flipping the operator alone would turn `NOT
	// a = ANY (q)` into `a <> ANY (q)`, which holds whenever any row of q differs from
	// a, where the original holds only when every row does. `= ANY` is the exception
	// the server makes, since it stores that one as IN: `(not(a in (q)))`.
	if binExpr, ok := unwrapParen(operand).(*ast.BinaryExpr); ok {
		if sub, ok := binExpr.Right.(*ast.SubqueryExpr); ok && sub.Quantifier != "" {
			inverted, canInvert := invertOp[binExpr.Op]
			if !canInvert || (binExpr.Op == ast.BinOpEq && sub.Quantifier != "ALL") {
				return &ast.UnaryExpr{Op: ast.UnaryNot, Operand: operand}
			}
			dual := *sub
			dual.Quantifier = "ALL"
			if sub.Quantifier == "ALL" {
				dual.Quantifier = "ANY"
			}
			return &ast.BinaryExpr{Loc: binExpr.Loc, Op: inverted, Left: binExpr.Left, Right: &dual}
		}
	}

	inner := unwrapParen(operand)

	// Case 1: NOT(comparison) → invert the comparison operator
	if binExpr, ok := inner.(*ast.BinaryExpr); ok {
		if inverted, canInvert := invertOp[binExpr.Op]; canInvert {
			return &ast.BinaryExpr{
				Loc:   binExpr.Loc,
				Op:    inverted,
				Left:  binExpr.Left,
				Right: binExpr.Right,
			}
		}
	}

	// Case 2a: NOT(LIKE) — keep as not() wrapping (don't fold into the LIKE)
	// The deparsing of UnaryNot already produces (not(...)), which is the correct
	// MySQL 8.0 output. So we return the UnaryNot as-is.
	if _, ok := inner.(*ast.LikeExpr); ok {
		return &ast.UnaryExpr{
			Op:      ast.UnaryNot,
			Operand: operand,
		}
	}

	// Case 2b: NOT(REGEXP) — keep as not() wrapping.
	// MySQL 8.0 rewrites NOT REGEXP to (not(regexp_like(...))).
	if binExpr, ok := inner.(*ast.BinaryExpr); ok && binExpr.Op == ast.BinOpRegexp {
		return &ast.UnaryExpr{
			Op:      ast.UnaryNot,
			Operand: operand,
		}
	}

	// Case 2c: NOT(a IN (subquery)) — keep as not() wrapping, which is what the server
	// stores, and what `a <> ALL (q)` prints as. `0 = a in (q)` would mean the same but
	// would not read back as what was written.
	if in, ok := inner.(*ast.InExpr); ok && in.Select != nil && !in.Not {
		return &ast.UnaryExpr{
			Op:      ast.UnaryNot,
			Operand: operand,
		}
	}

	if _, ok := inner.(*ast.ColumnRef); ok {
		return &ast.BinaryExpr{Op: ast.BinOpEq, Left: operand, Right: &ast.IntLit{Value: 0}}
	}
	return &ast.UnaryExpr{Op: ast.UnaryNot, Operand: operand}

}

// isBooleanExpr returns true if the expression is inherently boolean-valued.
// MySQL's is_bool_func() returns true for: comparisons (=,<>,<,>,<=,>=,<=>),
// IN, BETWEEN, LIKE, IS NULL/IS NOT NULL, AND, OR, NOT, XOR, EXISTS,
// TRUE/FALSE literals. Everything else (column refs, arithmetic, functions,
// CASE, IF, subqueries, literals) is NOT boolean and gets wrapped in (0 <> expr)
// when used as an operand of AND/OR/XOR.
func isBooleanExpr(node ast.ExprNode) bool {
	inner := unwrapParen(node)
	switch n := inner.(type) {
	case *ast.BinaryExpr:
		switch n.Op {
		case ast.BinOpEq, ast.BinOpNe, ast.BinOpLt, ast.BinOpGt,
			ast.BinOpLe, ast.BinOpGe, ast.BinOpNullSafeEq,
			ast.BinOpAnd, ast.BinOpOr, ast.BinOpXor, ast.BinOpSoundsLike:
			return true
		}
		return false
	case *ast.InExpr:
		return true
	case *ast.BetweenExpr:
		return true
	case *ast.LikeExpr:
		return true
	case *ast.IsExpr:
		return true
	case *ast.UnaryExpr:
		if n.Op == ast.UnaryNot {
			return true
		}
		return false
	case *ast.ExistsExpr:
		return true
	case *ast.BoolLit:
		return true
	default:
		return false
	}
}

// wrapBooleanContext wraps a non-boolean expression in (0 <> expr) for
// boolean context (AND/OR/XOR operands, IS TRUE/IS FALSE operands).
func wrapBooleanContext(node ast.ExprNode) ast.ExprNode {
	return &ast.BinaryExpr{Op: ast.BinOpNe, Left: node, Right: &ast.IntLit{Value: 0}}

}

// RewriteSelectStmt applies RewriteExpr to all expression positions in a SelectStmt.
// This should be called after resolver but before deparsing.
func RewriteSelectStmt(stmt *ast.SelectStmt) {
	if stmt == nil {
		return
	}
	if stmt.SetOp != ast.SetOpNone {
		RewriteSelectStmt(stmt.Left)
		RewriteSelectStmt(stmt.Right)
		return
	}
	if stmt.ParenSource != nil {
		RewriteSelectStmt(stmt.ParenSource)
		for _, item := range stmt.OrderBy {
			item.Expr = RewriteExpr(item.Expr)
		}
		return
	}
	for i, target := range stmt.TargetList {
		if rt, ok := target.(*ast.ResTarget); ok {
			rt.Val = RewriteExpr(rt.Val)
		} else {
			stmt.TargetList[i] = RewriteExpr(target)
		}
	}
	if stmt.Where != nil {
		stmt.Where = RewriteExpr(stmt.Where)
	}
	for i, expr := range stmt.GroupBy {
		stmt.GroupBy[i] = RewriteExpr(expr)
	}
	if stmt.Having != nil {
		stmt.Having = RewriteExpr(stmt.Having)
	}
	for _, item := range stmt.OrderBy {
		item.Expr = RewriteExpr(item.Expr)
	}
}
