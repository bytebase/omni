// Package deparse — rewrite.go implements AST rewrites applied before deparsing.
// These rewrites match MariaDB's resolver behavior, which is what
// information_schema.VIEWS holds for a view the server accepted.
package deparse

import (
	"strings"

	ast "github.com/bytebase/omni/mariadb/ast"
)

// RewriteExpr applies MariaDB resolver rewrites to an expression AST.
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
		// SOUNDS LIKE is stored as soundex(left) = soundex(right).
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
		// Boolean context wrapping: an AND or OR operand that is a bare column
		// is stored as `col <> 0`. Measured on mariadb:11.8.8: `a AND b` over
		// integer columns is `a <> 0 and b <> 0`, while `a AND (b + 1)` keeps
		// `b + 1`, `ifnull(a,0) AND b` keeps the call, `1 AND a` keeps the
		// literal, and XOR wraps neither side. MySQL 8.0 instead wraps every
		// operand that is not itself a boolean expression.
		if n.Op == ast.BinOpAnd || n.Op == ast.BinOpOr {
			if needsBooleanContext(n.Left) {
				n.Left = wrapBooleanContext(n.Left)
			}
			if needsBooleanContext(n.Right) {
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
		// MariaDB leaves an IS TRUE / IS FALSE operand as written: `a IS TRUE`
		// over an integer column is stored as `a is true`, where MySQL 8.0
		// wraps the operand. Measured on mariadb:11.8.8.
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

	// Case 0: De Morgan. MariaDB pushes a NOT through AND and OR rather than
	// writing one over the whole conjunction: `NOT (a > 0 AND b > 0)` is stored
	// as `a <= 0 or b <= 0`, and the push recurses. MySQL 8.0 writes
	// `(not((...)))`. Measured on mariadb:11.8.8.
	if binExpr, ok := inner.(*ast.BinaryExpr); ok {
		switch binExpr.Op {
		case ast.BinOpAnd, ast.BinOpOr:
			dual := ast.BinOpOr
			if binExpr.Op == ast.BinOpOr {
				dual = ast.BinOpAnd
			}
			left, right := binExpr.Left, binExpr.Right
			if needsBooleanContext(left) {
				left = wrapBooleanContext(left)
			}
			if needsBooleanContext(right) {
				right = wrapBooleanContext(right)
			}
			return &ast.BinaryExpr{
				Loc:   binExpr.Loc,
				Op:    dual,
				Left:  rewriteNot(left),
				Right: rewriteNot(right),
			}
		}
	}

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

	// Case 1b: a predicate that carries a NOT of its own takes it, rather than
	// a `!` over the whole: `NOT a IS TRUE` → `a is not true`,
	// `NOT a BETWEEN 1 AND 2` → `a not between 1 and 2`, and
	// `NOT a IN (1,2)` → `a not in (1,2)`. Measured on mariadb:11.8.8. An IN
	// over a subquery is the exception below, which the server stores as
	// `!(a in (q))`.
	switch pred := inner.(type) {
	case *ast.IsExpr:
		flipped := *pred
		flipped.Not = !pred.Not
		return &flipped
	case *ast.BetweenExpr:
		flipped := *pred
		flipped.Not = !pred.Not
		return &flipped
	case *ast.InExpr:
		if pred.Select == nil {
			flipped := *pred
			flipped.Not = !pred.Not
			return &flipped
		}
	}

	// Case 2a: NOT(LIKE) folds into the LIKE, which carries a NOT of its own:
	// `a NOT LIKE 'x'`. Measured on mariadb:11.8.8; MySQL 8.0 wraps it in
	// `not(...)` instead.
	if like, ok := inner.(*ast.LikeExpr); ok {
		flipped := *like
		flipped.Not = !like.Not
		return &flipped
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

	// Case 3: NOT over a bare column is stored as `col = 0`, the operand first,
	// where MySQL 8.0 writes the literal first. Anything else keeps the NOT,
	// which prints as `!`: `NOT (a + 1)` → `!(a + 1)` and
	// `NOT ifnull(a,0)` → `!ifnull(a,0)`. Measured on mariadb:11.8.8.
	if needsBooleanContext(operand) {
		return &ast.BinaryExpr{
			Op:    ast.BinOpEq,
			Left:  operand,
			Right: &ast.IntLit{Value: 0},
		}
	}
	return &ast.UnaryExpr{Op: ast.UnaryNot, Operand: operand}
}

// needsBooleanContext reports whether node is an operand MariaDB stores with an
// explicit `<> 0`: a bare column reference, and nothing else. Measured on
// mariadb:11.8.8 — see the AND/OR rewrite above.
func needsBooleanContext(node ast.ExprNode) bool {
	_, ok := unwrapParen(node).(*ast.ColumnRef)
	return ok
}

// wrapBooleanContext wraps a non-boolean expression in `expr <> 0` for boolean
// context (AND/OR/XOR operands, IS TRUE/IS FALSE operands). Measured on
// mariadb:11.8.8: `a AND b` over integer columns is stored as
// `a <> 0 and b <> 0`, the operand first, where MySQL 8.0 writes the literal
// first.
func wrapBooleanContext(node ast.ExprNode) ast.ExprNode {
	return &ast.BinaryExpr{
		Op:    ast.BinOpNe,
		Left:  node,
		Right: &ast.IntLit{Value: 0},
	}
}

// inlineNamedWindows replaces each OVER reference to a window the statement's
// WINDOW clause defines with that definition and drops the clause. MariaDB
// stores no WINDOW clause: `rank() OVER w ... WINDOW w AS (ORDER BY a)` comes
// back as `rank() over ( order by a)`. Measured on mariadb:11.8.8.
func inlineNamedWindows(stmt *ast.SelectStmt) {
	if len(stmt.WindowClause) == 0 {
		return
	}
	byName := make(map[string]*ast.WindowDef, len(stmt.WindowClause))
	for _, wd := range stmt.WindowClause {
		byName[strings.ToLower(wd.Name)] = wd
	}
	for _, target := range stmt.TargetList {
		expr, ok := target.(*ast.ResTarget)
		if !ok {
			continue
		}
		call, ok := expr.Val.(*ast.FuncCallExpr)
		if !ok || call.Over == nil || call.Over.RefName == "" {
			continue
		}
		if def, ok := byName[strings.ToLower(call.Over.RefName)]; ok {
			inlined := *def
			inlined.Name = ""
			call.Over = &inlined
		}
	}
	stmt.WindowClause = nil
}

// applyCTEColumnNames renames the output aliases of each CTE's query to the
// names the CTE declares, which is where MariaDB keeps them:
// `WITH cte(x, y) AS (SELECT a, b FROM t)` comes back as
// `with cte(`x`,`y`) as (select `t`.`a` AS `x`,`t`.`b` AS `y` from `t`)`.
// Measured on mariadb:11.8.8.
func applyCTEColumnNames(ctes []*ast.CommonTableExpr) {
	for _, cte := range ctes {
		if cte == nil || len(cte.Columns) == 0 || cte.Select == nil {
			continue
		}
		leaf := ast.LeftmostQueryLeaf(cte.Select)
		if leaf == nil {
			continue
		}
		for i, target := range leaf.TargetList {
			if i >= len(cte.Columns) {
				break
			}
			if rt, ok := target.(*ast.ResTarget); ok {
				rt.Name = cte.Columns[i]
			}
		}
	}
}

// RewriteSelectStmt applies RewriteExpr to all expression positions in a SelectStmt.
// This should be called after resolver but before deparsing.
func RewriteSelectStmt(stmt *ast.SelectStmt) {
	if stmt == nil {
		return
	}
	inlineNamedWindows(stmt)
	applyCTEColumnNames(stmt.CTEs)
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
	for _, item := range stmt.GroupBy {
		item.Expr = RewriteExpr(item.Expr)
	}
	if stmt.Having != nil {
		stmt.Having = RewriteExpr(stmt.Having)
	}
	for _, item := range stmt.OrderBy {
		item.Expr = RewriteExpr(item.Expr)
	}
}
