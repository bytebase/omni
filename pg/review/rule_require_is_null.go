package review

import (
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/review"
)

// checkRequireIsNull reports a comparison of the form x = NULL or
// x <> NULL (the lexer spells != as <>), which is null for every row and so
// never true. The literal may sit on either side and may carry a cast,
// as in NULL::int. The finding anchors the comparison and quotes it.
func checkRequireIsNull(sql string, s *statement, r *reporter) {
	ast.Inspect(s.node, func(n ast.Node) bool {
		e, ok := n.(*ast.A_Expr)
		if !ok || e.Kind != ast.AEXPR_OP {
			return true
		}
		op := operatorName(e.Name)
		if op != "=" && op != "<>" {
			return true
		}
		if !isNullLiteral(e.Lexpr) && !isNullLiteral(e.Rexpr) {
			return true
		}
		fix := "IS NULL"
		if op == "<>" {
			fix = "IS NOT NULL"
		}
		rng := rangeOf(e.Loc)
		text := op + " NULL"
		if rng.Start < rng.End && rng.End <= len(sql) {
			text = sql[rng.Start:rng.End]
		}
		r.report(review.RequireIsNull, s.index, rng, `"`+text+`" is never true; use `+fix)
		return true
	})
}

// operatorName returns the unqualified operator of an A_Expr name list.
func operatorName(name *ast.List) string {
	parts := nameParts(name)
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}

// isNullLiteral reports whether the expression is the NULL constant,
// possibly under casts.
func isNullLiteral(n ast.Node) bool {
	for {
		switch v := n.(type) {
		case *ast.A_Const:
			return v.Isnull
		case *ast.TypeCast:
			n = v.Arg
		default:
			return false
		}
	}
}
