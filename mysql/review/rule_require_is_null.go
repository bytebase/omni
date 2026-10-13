package review

import (
	"strings"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/review"
)

// checkRequireIsNull reports a comparison of the form x = NULL, x <> NULL,
// or x != NULL that a predicate tests for truth, where it is null for
// every row and so never true. The NULL may sit on either side. Only a
// comparison the predicate reads as its condition counts: the predicate
// itself, or an operand of AND, OR, XOR, or NOT in it, parenthesized or
// not. A comparison that produces a value, as an argument or in a select
// list, is left alone, and so is the null-safe <=>. The finding anchors
// the comparison and quotes it.
func checkRequireIsNull(sql string, s *statement, r *reporter) {
	seen := make(map[*ast.BinaryExpr]bool)
	var visit func(n ast.Node)
	visit = func(n ast.Node) {
		switch v := n.(type) {
		case *ast.ParenExpr:
			visit(v.Expr)
		case *ast.UnaryExpr:
			if v.Op == ast.UnaryNot {
				visit(v.Operand)
			}
		case *ast.BinaryExpr:
			switch v.Op {
			case ast.BinOpAnd, ast.BinOpOr, ast.BinOpXor:
				visit(v.Left)
				visit(v.Right)
			case ast.BinOpEq, ast.BinOpNe:
				if seen[v] || !isNull(v.Left) && !isNull(v.Right) {
					return
				}
				seen[v] = true
				fix := "IS NULL"
				if v.Op == ast.BinOpNe {
					fix = "IS NOT NULL"
				}
				rng := spanOf(sql, v)
				text := "= NULL"
				if v.Op == ast.BinOpNe {
					text = "<> NULL"
				}
				if rng.Start < rng.End && rng.End <= len(sql) {
					text = strings.Join(strings.Fields(sql[rng.Start:rng.End]), " ")
				}
				r.report(review.RequireIsNull, s.index, rng, `"`+text+`" is never true; use `+fix)
			}
		}
	}
	for _, p := range predicates(s.node) {
		visit(p)
	}
}

// isNull reports whether an expression is the NULL literal, parenthesized
// or not.
func isNull(n ast.Node) bool {
	switch v := n.(type) {
	case *ast.NullLit:
		return true
	case *ast.ParenExpr:
		return isNull(v.Expr)
	}
	return false
}

// predicates lists the expressions of the statement that are evaluated
// for truth: WHERE and HAVING, join conditions, the conditions of a
// searched CASE, CHECK expressions, and the conditions of IF, ELSEIF,
// WHILE, and REPEAT in a routine body. Subqueries are visited too, so
// their own predicates are listed.
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
		case *ast.FuncCallExpr:
			// IF(condition, a, b) tests its first argument.
			if v.Schema == "" && strings.EqualFold(v.Name, "IF") && len(v.Args) == 3 {
				add(v.Args[0])
			}
		case *ast.SelectStmt:
			if v.Where != nil {
				add(v.Where)
			}
			if v.Having != nil {
				add(v.Having)
			}
		case *ast.ShowStmt:
			if v.Where != nil {
				add(v.Where)
			}
		case *ast.UpdateStmt:
			if v.Where != nil {
				add(v.Where)
			}
		case *ast.DeleteStmt:
			if v.Where != nil {
				add(v.Where)
			}
		case *ast.OnCondition:
			if v.Expr != nil {
				add(v.Expr)
			}
		case *ast.CaseExpr:
			// A simple CASE compares its WHEN values with the operand;
			// only a searched CASE tests its WHEN expressions.
			if v.Operand == nil {
				for _, w := range v.Whens {
					if w != nil && w.Cond != nil {
						add(w.Cond)
					}
				}
			}
		case *ast.Constraint:
			if v.Type == ast.ConstrCheck && v.Expr != nil {
				add(v.Expr)
			}
		case *ast.ColumnConstraint:
			if v.Type == ast.ColConstrCheck && v.Expr != nil {
				add(v.Expr)
			}
		case *ast.IfStmt:
			if v.Cond != nil {
				add(v.Cond)
			}
		case *ast.CaseStmtNode:
			// Only a searched CASE statement tests its WHEN conditions.
			if v.Operand == nil {
				for _, w := range v.Whens {
					if w != nil && w.Cond != nil {
						add(w.Cond)
					}
				}
			}
		case *ast.ElseIf:
			if v.Cond != nil {
				add(v.Cond)
			}
		case *ast.WhileStmt:
			if v.Cond != nil {
				add(v.Cond)
			}
		case *ast.RepeatStmt:
			if v.Cond != nil {
				add(v.Cond)
			}
		}
		return true
	})
	return roots
}
