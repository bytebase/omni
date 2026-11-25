package analysis

import (
	"github.com/bytebase/omni/cosmosdb/ast"
)

// exprPaths appends the field paths expr's value comes from to dst, in the
// order they occur, and returns the extended slice.
func (a *analyzer) exprPaths(dst []FieldPath, expr ast.ExprNode, sc *scope) []FieldPath {
	return union(a.exprShapes(nil, expr, sc)).paths(dst)
}

// exprShapes appends the shapes of the values expr is computed from to dst.
// Every expression appends to the one slice: returning a slice per
// subexpression and copying it into the parent's re-copied it at each
// level, which made a chain of n ORs take time quadratic in n (16,000 took
// seconds). The shapes are the values' own, shared rather than flattened, so
// s + s over a binding that is itself s + s stays one shape per level.
func (a *analyzer) exprShapes(dst []*shape, expr ast.ExprNode, sc *scope) []*shape {
	if expr == nil {
		return dst
	}

	switch e := expr.(type) {
	case *ast.ColumnRef, *ast.DotAccessExpr, *ast.BracketAccessExpr,
		*ast.SubLink, *ast.ArrayExpr, *ast.SubqueryFuncExpr, *ast.ExistsExpr,
		*ast.CreateObjectExpr, *ast.CreateArrayExpr:
		// Names, property paths, literals and subqueries have shapes of
		// their own.
		return append(dst, a.shapeOf(e, sc))

	case *ast.BinaryExpr:
		dst = a.exprShapes(dst, e.Left, sc)
		return a.exprShapes(dst, e.Right, sc)

	case *ast.UnaryExpr:
		return a.exprShapes(dst, e.Operand, sc)

	case *ast.TernaryExpr:
		dst = a.exprShapes(dst, e.Cond, sc)
		dst = a.exprShapes(dst, e.Then, sc)
		return a.exprShapes(dst, e.Else, sc)

	case *ast.InExpr:
		dst = a.exprShapes(dst, e.Expr, sc)
		for _, item := range e.List {
			dst = a.exprShapes(dst, item, sc)
		}
		return dst

	case *ast.BetweenExpr:
		dst = a.exprShapes(dst, e.Expr, sc)
		dst = a.exprShapes(dst, e.Low, sc)
		return a.exprShapes(dst, e.High, sc)

	case *ast.LikeExpr:
		dst = a.exprShapes(dst, e.Expr, sc)
		dst = a.exprShapes(dst, e.Pattern, sc)
		return a.exprShapes(dst, e.Escape, sc)

	case *ast.FuncCall:
		return a.exprShapesOf(dst, e.Args, sc)

	case *ast.UDFCall:
		return a.exprShapesOf(dst, e.Args, sc)

	default:
		// Literals and parameters read no field.
		return dst
	}
}

// exprShapesOf appends the shapes of each expression in exprs.
func (a *analyzer) exprShapesOf(dst []*shape, exprs []ast.ExprNode, sc *scope) []*shape {
	for _, e := range exprs {
		dst = a.exprShapes(dst, e, sc)
	}
	return dst
}
