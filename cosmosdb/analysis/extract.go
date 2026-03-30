package analysis

import (
	"github.com/bytebase/omni/cosmosdb/ast"
)

// extractFieldPaths recursively collects all field paths referenced in an expression.
func extractFieldPaths(expr ast.ExprNode, aliases aliasMap) []FieldPath {
	if expr == nil {
		return nil
	}

	switch e := expr.(type) {
	case *ast.ColumnRef:
		if resolved, ok := aliases[e.Name]; ok {
			return []FieldPath{copyPath(resolved)}
		}
		return []FieldPath{{ItemSelector(e.Name)}}

	case *ast.DotAccessExpr:
		bases := extractFieldPaths(e.Expr, aliases)
		for i := range bases {
			bases[i] = append(bases[i], ItemSelector(e.Property))
		}
		return bases

	case *ast.BracketAccessExpr:
		bases := extractFieldPaths(e.Expr, aliases)
		sel := bracketSelector(e.Index)
		for i := range bases {
			bases[i] = append(bases[i], sel)
		}
		return bases

	case *ast.BinaryExpr:
		return collectFromExprs(aliases, e.Left, e.Right)

	case *ast.UnaryExpr:
		return extractFieldPaths(e.Operand, aliases)

	case *ast.TernaryExpr:
		return collectFromExprs(aliases, e.Cond, e.Then, e.Else)

	case *ast.InExpr:
		paths := extractFieldPaths(e.Expr, aliases)
		for _, item := range e.List {
			paths = append(paths, extractFieldPaths(item, aliases)...)
		}
		return paths

	case *ast.BetweenExpr:
		return collectFromExprs(aliases, e.Expr, e.Low, e.High)

	case *ast.LikeExpr:
		paths := collectFromExprs(aliases, e.Expr, e.Pattern)
		if e.Escape != nil {
			paths = append(paths, extractFieldPaths(e.Escape, aliases)...)
		}
		return paths

	case *ast.FuncCall:
		return collectFromExprSlice(e.Args, aliases)

	case *ast.UDFCall:
		return collectFromExprSlice(e.Args, aliases)

	case *ast.CreateArrayExpr:
		return collectFromExprSlice(e.Elements, aliases)

	case *ast.CreateObjectExpr:
		var paths []FieldPath
		for _, f := range e.Fields {
			paths = append(paths, extractFieldPaths(f.Value, aliases)...)
		}
		return paths

	case *ast.StringLit, *ast.NumberLit, *ast.BoolLit,
		*ast.NullLit, *ast.UndefinedLit, *ast.InfinityLit,
		*ast.NanLit, *ast.ParamRef:
		return nil

	case *ast.SubLink, *ast.ExistsExpr, *ast.ArrayExpr:
		return nil

	default:
		return nil
	}
}

// collectFromExprs collects field paths from multiple expressions.
func collectFromExprs(aliases aliasMap, exprs ...ast.ExprNode) []FieldPath {
	var all []FieldPath
	for _, e := range exprs {
		all = append(all, extractFieldPaths(e, aliases)...)
	}
	return all
}

// collectFromExprSlice collects field paths from a slice of expressions.
func collectFromExprSlice(exprs []ast.ExprNode, aliases aliasMap) []FieldPath {
	var all []FieldPath
	for _, e := range exprs {
		all = append(all, extractFieldPaths(e, aliases)...)
	}
	return all
}
