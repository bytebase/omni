package analysis

import (
	"github.com/bytebase/omni/cosmosdb/ast"
)

// exprPaths appends the field paths expr's value comes from to dst, in the
// order they occur, and returns the extended slice. Every expression
// appends to the one slice: returning a slice per subexpression and
// copying it into the parent's re-copied every path at each level, which
// made a chain of n ORs take time quadratic in n (16,000 took seconds).
func (a *analyzer) exprPaths(dst []FieldPath, expr ast.ExprNode, sc *scope) []FieldPath {
	if expr == nil {
		return dst
	}

	switch e := expr.(type) {
	case *ast.ColumnRef, *ast.DotAccessExpr, *ast.BracketAccessExpr,
		*ast.SubLink, *ast.ArrayExpr, *ast.SubqueryFuncExpr, *ast.ExistsExpr:
		// Names, property paths and subqueries resolve through the shapes
		// of the values they name.
		return a.shapeOf(e, sc).paths(dst)

	case *ast.BinaryExpr:
		dst = a.exprPaths(dst, e.Left, sc)
		return a.exprPaths(dst, e.Right, sc)

	case *ast.UnaryExpr:
		return a.exprPaths(dst, e.Operand, sc)

	case *ast.TernaryExpr:
		dst = a.exprPaths(dst, e.Cond, sc)
		dst = a.exprPaths(dst, e.Then, sc)
		return a.exprPaths(dst, e.Else, sc)

	case *ast.InExpr:
		dst = a.exprPaths(dst, e.Expr, sc)
		for _, item := range e.List {
			dst = a.exprPaths(dst, item, sc)
		}
		return dst

	case *ast.BetweenExpr:
		dst = a.exprPaths(dst, e.Expr, sc)
		dst = a.exprPaths(dst, e.Low, sc)
		return a.exprPaths(dst, e.High, sc)

	case *ast.LikeExpr:
		dst = a.exprPaths(dst, e.Expr, sc)
		dst = a.exprPaths(dst, e.Pattern, sc)
		return a.exprPaths(dst, e.Escape, sc)

	case *ast.FuncCall:
		return a.exprPathsOf(dst, e.Args, sc)

	case *ast.UDFCall:
		return a.exprPathsOf(dst, e.Args, sc)

	case *ast.CreateArrayExpr:
		return a.exprPathsOf(dst, e.Elements, sc)

	case *ast.CreateObjectExpr:
		for _, f := range e.Fields {
			dst = a.exprPaths(dst, f.Value, sc)
		}
		return dst

	default:
		// Literals and parameters read no field.
		return dst
	}
}

// exprPathsOf appends the field paths of each expression in exprs.
func (a *analyzer) exprPathsOf(dst []FieldPath, exprs []ast.ExprNode, sc *scope) []FieldPath {
	for _, e := range exprs {
		dst = a.exprPaths(dst, e, sc)
	}
	return dst
}
