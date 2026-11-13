// Package deparse — rewrite.go implements AST rewrites applied before deparsing.
// These rewrites match MySQL 8.0's resolver behavior (SHOW CREATE VIEW output).
package deparse

import (
	"strings"

	ast "github.com/bytebase/omni/oceanbase/ast"
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
		if n.Op == ast.BinOpAnd || n.Op == ast.BinOpOr || n.Op == ast.BinOpXor {
			if !isBooleanExpr(n.Left) {
				n.Left = wrapBooleanContext(n.Left)
			}
			if !isBooleanExpr(n.Right) {
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
		// x IN (subquery) is stored as x = any (subquery) and x NOT IN (subquery) as
		// x <> all (subquery), whatever parentheses the subquery is written in
		// (catalog's TestOracle_ViewIdempotence, in-subquery).
		if n.Select != nil {
			op, quantifier := ast.BinOpEq, "ANY"
			if n.Not {
				op, quantifier = ast.BinOpNe, "ALL"
			}
			return &ast.BinaryExpr{
				Loc:   n.Loc,
				Op:    op,
				Left:  n.Expr,
				Right: &ast.SubqueryExpr{Loc: n.Loc, Select: n.Select, Quantifier: quantifier},
			}
		}
		return n

	case *ast.BetweenExpr:
		n.Expr = rewriteExpr(n.Expr)
		n.Low = rewriteExpr(n.Low)
		n.High = rewriteExpr(n.High)
		// BETWEEN over two constant bounds is stored as two comparisons,
		// ((x >= lo) and (x <= hi)), and NOT BETWEEN as their negation. A bound that
		// reads a column keeps BETWEEN (catalog's TestOracle_CheckDiffIdempotence,
		// between and between-forms).
		if both, folded := foldBetween(n); folded {
			if n.Not {
				return rewriteNot(&ast.ParenExpr{Loc: n.Loc, Expr: both})
			}
			return both
		}
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
		// IS TRUE / IS FALSE on non-boolean: wrap operand in (0 <> expr)
		if (n.Test == ast.IsTrue || n.Test == ast.IsFalse) && !n.Not && !isBooleanExpr(n.Expr) {
			n.Expr = wrapBooleanContext(n.Expr)
		}
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
		// IF(c, a, b) is stored as (case when c then a else b end) (catalog's
		// TestOracle_ViewIdempotence, paren-subquery-compare-if).
		if n.Schema == "" && strings.EqualFold(n.Name, "if") && len(n.Args) == 3 {
			return &ast.CaseExpr{
				Loc:     n.Loc,
				Whens:   []*ast.CaseWhen{{Loc: n.Loc, Cond: n.Args[0], Result: n.Args[1]}},
				Default: n.Args[2],
			}
		}
		return rewriteDateArith(n)

	case *ast.ExtractExpr:
		n.Expr = rewriteExpr(n.Expr)
		return n

	case *ast.IntervalExpr:
		n.Value = rewriteExpr(n.Value)
		return n

	case *ast.WeightStringExpr:
		n.Expr = rewriteExpr(n.Expr)
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

// foldBetween returns the two comparisons OceanBase stores for a BETWEEN over constant
// bounds, without the negation of a NOT BETWEEN.
//
// The fold writes the tested operand twice, so BETWEENs nested in the operand double the
// text at every level: twenty of them turn 368 bytes into 25 MB. An operand past
// maxFoldedOperand nodes therefore keeps its BETWEEN, which costs such a view a
// difference from its readback and bounds what one statement can make the catalog hold.
func foldBetween(n *ast.BetweenExpr) (ast.ExprNode, bool) {
	if !isConstantBound(n.Low) || !isConstantBound(n.High) || largerThan(n.Expr, maxFoldedOperand) {
		return nil, false
	}
	return &ast.BinaryExpr{
		Loc:   n.Loc,
		Op:    ast.BinOpAnd,
		Left:  &ast.BinaryExpr{Loc: n.Loc, Op: ast.BinOpGe, Left: n.Expr, Right: n.Low},
		Right: &ast.BinaryExpr{Loc: n.Loc, Op: ast.BinOpLe, Left: n.Expr, Right: n.High},
	}, true
}

// FoldBetween rewrites every BETWEEN over constant bounds in an expression as OceanBase
// stores it, and nothing else: the one rewrite it makes of a CHECK and of a generated
// column's expression too, which RewriteExpr is not applied to.
func FoldBetween(node ast.ExprNode) ast.ExprNode {
	switch n := node.(type) {
	case *ast.BetweenExpr:
		n.Expr, n.Low, n.High = FoldBetween(n.Expr), FoldBetween(n.Low), FoldBetween(n.High)
		both, folded := foldBetween(n)
		switch {
		case !folded:
			return n
		case n.Not:
			return &ast.UnaryExpr{Loc: n.Loc, Op: ast.UnaryNot, Operand: &ast.ParenExpr{Loc: n.Loc, Expr: both}}
		}
		return both
	case *ast.BinaryExpr:
		n.Left, n.Right = FoldBetween(n.Left), FoldBetween(n.Right)
	case *ast.UnaryExpr:
		n.Operand = FoldBetween(n.Operand)
	case *ast.ParenExpr:
		n.Expr = FoldBetween(n.Expr)
	case *ast.IsExpr:
		n.Expr = FoldBetween(n.Expr)
	case *ast.CaseExpr:
		if n.Operand != nil {
			n.Operand = FoldBetween(n.Operand)
		}
		for _, w := range n.Whens {
			w.Cond, w.Result = FoldBetween(w.Cond), FoldBetween(w.Result)
		}
		if n.Default != nil {
			n.Default = FoldBetween(n.Default)
		}
	case *ast.FuncCallExpr:
		for i, arg := range n.Args {
			n.Args[i] = FoldBetween(arg)
		}
	}
	return node
}

// maxFoldedOperand is the size of the largest operand a BETWEEN is folded over, in nodes.
// It holds the operands people write, a column or an expression over a few, and stops
// the third level of a nest.
const maxFoldedOperand = 32

// largerThan reports whether an expression has more than limit nodes, counting a subtree
// two nodes share once for each, as it is written. It stops at the limit, so it costs no
// more than the limit however far a fold has already spread the tree.
func largerThan(node ast.ExprNode, limit int) bool {
	count := 0
	ast.Inspect(node, func(ast.Node) bool {
		count++
		return count <= limit
	})
	return count > limit
}

// isConstantBound reports whether a bound of BETWEEN is a literal, or arithmetic over
// literals: what OceanBase folds the BETWEEN for.
func isConstantBound(node ast.ExprNode) bool {
	switch n := node.(type) {
	case *ast.IntLit, *ast.FloatLit, *ast.StringLit, *ast.BoolLit:
		return true
	case *ast.ParenExpr:
		return isConstantBound(n.Expr)
	case *ast.UnaryExpr:
		return n.Op != ast.UnaryNot && isConstantBound(n.Operand)
	case *ast.BinaryExpr:
		switch n.Op {
		case ast.BinOpAdd, ast.BinOpSub, ast.BinOpMul, ast.BinOpDiv:
			return isConstantBound(n.Left) && isConstantBound(n.Right)
		}
	}
	return false
}

// rewriteDateArith folds the date-arithmetic function family into INTERVAL
// arithmetic, matching the engine's resolver: TIMESTAMPADD, DATE_ADD/ADDDATE
// and DATE_SUB/SUBDATE never appear in a stored body — they are stored as
// (expr + interval n unit) / (expr - interval n unit) (oracle 8.0.32 +
// 5.7.25, all units). ADDDATE/SUBDATE with a bare count means days:
// adddate(d, 31) stores as (d + interval 31 day). A schema-qualified name
// is a stored function, not the builtin, and is left alone.
func rewriteDateArith(n *ast.FuncCallExpr) ast.ExprNode {
	if n.Schema != "" {
		return n
	}
	switch strings.ToUpper(n.Name) {
	case "TIMESTAMPADD":
		if len(n.Args) == 3 {
			if unit, ok := n.Args[0].(*ast.KeywordArg); ok {
				return &ast.BinaryExpr{
					Op:   ast.BinOpAdd,
					Left: n.Args[2],
					Right: &ast.IntervalExpr{
						Value: n.Args[1],
						Unit:  unit.Keyword,
					},
				}
			}
		}
	case "DATE_ADD", "ADDDATE":
		return rewriteDateAddSub(n, ast.BinOpAdd)
	case "DATE_SUB", "SUBDATE":
		return rewriteDateAddSub(n, ast.BinOpSub)
	}
	return n
}

// rewriteDateAddSub folds DATE_ADD/DATE_SUB/ADDDATE/SUBDATE into interval
// arithmetic. DATE_ADD/DATE_SUB require INTERVAL syntax for the second
// argument; only ADDDATE/SUBDATE also have the bare-days form.
func rewriteDateAddSub(n *ast.FuncCallExpr, op ast.BinaryOp) ast.ExprNode {
	if len(n.Args) != 2 {
		return n
	}
	iv, ok := n.Args[1].(*ast.IntervalExpr)
	if !ok {
		name := strings.ToUpper(n.Name)
		if name != "ADDDATE" && name != "SUBDATE" {
			return n
		}
		iv = &ast.IntervalExpr{Value: n.Args[1], Unit: "DAY"}
	}
	return &ast.BinaryExpr{Op: op, Left: n.Args[0], Right: iv}
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
	// NOT of a quantified comparison stays a NOT: flipping the operator alone would
	// turn `NOT a = ANY (q)` into `a <> ANY (q)`, which holds whenever any row of q
	// differs from a, where the original holds only when every row does.
	if binExpr, ok := unwrapParen(operand).(*ast.BinaryExpr); ok {
		if sub, ok := binExpr.Right.(*ast.SubqueryExpr); ok && sub.Quantifier != "" {
			return &ast.UnaryExpr{Op: ast.UnaryNot, Operand: operand}
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

	// Case 3: NOT(non-boolean) → (0 = expr)
	// This handles: NOT(column), NOT(a+1), NOT(func()), etc.
	// MySQL rewrites these as (0 = expr).
	return &ast.BinaryExpr{
		Op:    ast.BinOpEq,
		Left:  &ast.IntLit{Value: 0},
		Right: operand,
	}
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
	return &ast.BinaryExpr{
		Op:    ast.BinOpNe,
		Left:  &ast.IntLit{Value: 0},
		Right: node,
	}
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
		orderByPosition(stmt)
		return
	}
	if stmt.ParenSource != nil {
		RewriteSelectStmt(stmt.ParenSource)
		for _, item := range stmt.OrderBy {
			item.Expr = RewriteExpr(item.Expr)
		}
		orderByPosition(stmt)
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
	orderByPosition(stmt)
}

// orderByPosition writes an ORDER BY item that is one of the query's select items as that
// item's position, which is how OceanBase stores it: ORDER BY x, for a select item aliased
// x or for the column a select item is, reads back as order by 1 (catalog's
// TestOracle_ViewIdempotence, order-by-*). The select items of a set operation, or of a
// query in parentheses, are those of its leftmost SELECT.
func orderByPosition(stmt *ast.SelectStmt) {
	if len(stmt.OrderBy) == 0 {
		return
	}
	leaf := ast.LeftmostQueryLeaf(stmt)
	if leaf == nil {
		return
	}
	for _, item := range stmt.OrderBy {
		if pos := selectItemPosition(item.Expr, leaf.TargetList); pos > 0 {
			item.Expr = &ast.IntLit{Value: int64(pos)}
		}
	}
}

// selectItemPosition returns the 1-based position of the select item an ORDER BY
// expression names, or 0: a column written with no table that is an item's alias, or the
// very column an item is.
func selectItemPosition(expr ast.ExprNode, items []ast.ExprNode) int {
	ref, ok := expr.(*ast.ColumnRef)
	if !ok || ref.Star {
		return 0
	}
	if ref.Table == "" && ref.Schema == "" {
		for i, item := range items {
			if rt, ok := item.(*ast.ResTarget); ok && rt.Name != "" && strings.EqualFold(rt.Name, ref.Column) {
				return i + 1
			}
		}
	}
	for i, item := range items {
		if rt, ok := item.(*ast.ResTarget); ok {
			item = rt.Val
		}
		if col, ok := item.(*ast.ColumnRef); ok && !col.Star &&
			strings.EqualFold(col.Column, ref.Column) && strings.EqualFold(col.Table, ref.Table) && strings.EqualFold(col.Schema, ref.Schema) {
			return i + 1
		}
	}
	return 0
}
