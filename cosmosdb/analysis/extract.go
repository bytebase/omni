package analysis

import (
	"github.com/bytebase/omni/cosmosdb/ast"
)

// extractFieldPaths collects all field paths referenced in an expression.
func extractFieldPaths(expr ast.ExprNode, aliases aliasMap) []FieldPath {
	return appendFieldPaths(nil, expr, aliases)
}

// appendFieldPaths appends the field paths referenced in expr to dst, in
// the order they occur, and returns the extended slice. Every expression
// appends to the one slice: returning a slice per subexpression and
// copying it into the parent's re-copied every path at each level, which
// made a chain of n ORs take time quadratic in n (16,000 took seconds).
func appendFieldPaths(dst []FieldPath, expr ast.ExprNode, aliases aliasMap) []FieldPath {
	if expr == nil {
		return dst
	}

	switch e := expr.(type) {
	case *ast.ColumnRef:
		if resolved, ok := aliases[e.Name]; ok {
			return append(dst, copyPath(resolved))
		}
		return append(dst, FieldPath{ItemSelector(e.Name)})

	case *ast.DotAccessExpr:
		start := len(dst)
		dst = appendFieldPaths(dst, e.Expr, aliases)
		for i := start; i < len(dst); i++ {
			dst[i] = append(dst[i], ItemSelector(e.Property))
		}
		return dst

	case *ast.BracketAccessExpr:
		start := len(dst)
		dst = appendFieldPaths(dst, e.Expr, aliases)
		sel := bracketSelector(e.Index)
		for i := start; i < len(dst); i++ {
			dst[i] = append(dst[i], sel)
		}
		return dst

	case *ast.BinaryExpr:
		dst = appendFieldPaths(dst, e.Left, aliases)
		return appendFieldPaths(dst, e.Right, aliases)

	case *ast.UnaryExpr:
		return appendFieldPaths(dst, e.Operand, aliases)

	case *ast.TernaryExpr:
		dst = appendFieldPaths(dst, e.Cond, aliases)
		dst = appendFieldPaths(dst, e.Then, aliases)
		return appendFieldPaths(dst, e.Else, aliases)

	case *ast.InExpr:
		dst = appendFieldPaths(dst, e.Expr, aliases)
		for _, item := range e.List {
			dst = appendFieldPaths(dst, item, aliases)
		}
		return dst

	case *ast.BetweenExpr:
		dst = appendFieldPaths(dst, e.Expr, aliases)
		dst = appendFieldPaths(dst, e.Low, aliases)
		return appendFieldPaths(dst, e.High, aliases)

	case *ast.LikeExpr:
		dst = appendFieldPaths(dst, e.Expr, aliases)
		dst = appendFieldPaths(dst, e.Pattern, aliases)
		return appendFieldPaths(dst, e.Escape, aliases)

	case *ast.FuncCall:
		return appendFieldPathsOf(dst, e.Args, aliases)

	case *ast.UDFCall:
		return appendFieldPathsOf(dst, e.Args, aliases)

	case *ast.CreateArrayExpr:
		return appendFieldPathsOf(dst, e.Elements, aliases)

	case *ast.CreateObjectExpr:
		for _, f := range e.Fields {
			dst = appendFieldPaths(dst, f.Value, aliases)
		}
		return dst

	case *ast.StringLit, *ast.NumberLit, *ast.BoolLit,
		*ast.NullLit, *ast.UndefinedLit, *ast.InfinityLit,
		*ast.NanLit, *ast.ParamRef:
		return dst

	case *ast.SubLink, *ast.ExistsExpr, *ast.ArrayExpr:
		return dst

	default:
		return dst
	}
}

// appendFieldPathsOf appends the field paths of each expression in exprs.
func appendFieldPathsOf(dst []FieldPath, exprs []ast.ExprNode, aliases aliasMap) []FieldPath {
	for _, e := range exprs {
		dst = appendFieldPaths(dst, e, aliases)
	}
	return dst
}
