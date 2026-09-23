package review

import (
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/review"
)

// checkRequireIsNull reports a comparison of the form x = NULL or
// x <> NULL (the lexer spells != as <>) in a predicate, where it is null
// for every row and so never true. The literal may sit on either side and
// may carry a cast or a collation, as in NULL::text COLLATE "C". Only the
// built-in operator counts: OPERATOR(s.=) names a user's operator, which
// may not be strict. A comparison that produces a value, as in
// SELECT a = NULL AS c, is not a predicate and is left alone. The finding
// anchors the comparison and quotes it.
func checkRequireIsNull(sql string, s *statement, r *reporter) {
	seen := make(map[*ast.A_Expr]bool)
	for _, root := range predicates(s.node) {
		ast.Inspect(root, func(n ast.Node) bool {
			if _, ok := n.(*ast.SelectStmt); ok {
				// A subquery's own predicates are roots of their own;
				// its other expressions produce values.
				return false
			}
			e, ok := n.(*ast.A_Expr)
			if !ok || e.Kind != ast.AEXPR_OP || seen[e] {
				return true
			}
			op, ok := builtinOperator(e.Name)
			if !ok || (op != "=" && op != "<>") {
				return true
			}
			if !isNullLiteral(e.Lexpr) && !isNullLiteral(e.Rexpr) {
				return true
			}
			seen[e] = true
			fix := "IS NULL"
			if op == "<>" {
				fix = "IS NOT NULL"
			}
			rng := rangeOf(e.Loc)
			text := op + " NULL"
			if rng.Start < rng.End && rng.End <= len(sql) {
				text = collapseSpace(sql[rng.Start:rng.End])
			}
			r.report(review.RequireIsNull, s.index, rng, `"`+text+`" is never true; use `+fix)
			return true
		})
	}
}

// predicates lists the expressions of the statement that are evaluated
// for truth: WHERE and HAVING, join conditions, MERGE conditions, the
// conditions of a searched CASE, aggregate FILTER, partial index and
// constraint predicates, CHECK expressions, policy qualifications, rule
// and COPY conditions. Subqueries are visited too, so their own
// predicates are listed.
func predicates(root ast.Node) []ast.Node {
	var roots []ast.Node
	add := func(nodes ...ast.Node) {
		for _, n := range nodes {
			if n != nil {
				roots = append(roots, n)
			}
		}
	}
	ast.Inspect(root, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.SelectStmt:
			add(v.WhereClause, v.HavingClause)
		case *ast.UpdateStmt:
			add(v.WhereClause)
		case *ast.DeleteStmt:
			add(v.WhereClause)
		case *ast.MergeStmt:
			add(v.JoinCondition)
		case *ast.MergeWhenClause:
			add(v.Condition)
		case *ast.JoinExpr:
			add(v.Quals)
		case *ast.CaseExpr:
			// A simple CASE compares its WHEN values with the argument;
			// only a searched CASE tests its WHEN expressions.
			if v.Arg == nil && v.Args != nil {
				for _, item := range v.Args.Items {
					if w, ok := item.(*ast.CaseWhen); ok {
						add(w.Expr)
					}
				}
			}
		case *ast.FuncCall:
			add(v.AggFilter)
		case *ast.IndexStmt:
			add(v.WhereClause)
		case *ast.Constraint:
			add(v.RawExpr, v.WhereClause)
		case *ast.OnConflictClause:
			add(v.WhereClause)
		case *ast.InferClause:
			add(v.WhereClause)
		case *ast.CreatePolicyStmt:
			add(v.Qual, v.WithCheck)
		case *ast.AlterPolicyStmt:
			add(v.Qual, v.WithCheck)
		case *ast.RuleStmt:
			add(v.WhereClause)
		case *ast.CopyStmt:
			add(v.WhereClause)
		}
		return true
	})
	return roots
}

// builtinOperator returns the operator of an A_Expr name list when it
// names a built-in one: unqualified, or qualified with pg_catalog.
func builtinOperator(name *ast.List) (string, bool) {
	parts := nameParts(name)
	switch len(parts) {
	case 1:
		return parts[0], true
	case 2:
		return parts[1], parts[0] == "pg_catalog"
	default:
		return "", false
	}
}

// isNullLiteral reports whether the expression is the NULL constant,
// possibly under casts and collations, as in NULL::text COLLATE "C".
func isNullLiteral(n ast.Node) bool {
	for {
		switch v := n.(type) {
		case *ast.A_Const:
			return v.Isnull
		case *ast.TypeCast:
			n = v.Arg
		case *ast.CollateClause:
			n = v.Arg
		default:
			return false
		}
	}
}
