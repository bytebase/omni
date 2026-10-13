package review

import (
	"strings"

	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/review"
)

// checkRequireIsNull reports a comparison of the form x = NULL, x <> NULL,
// or x != NULL that a predicate tests for truth, where with ANSI_NULLS ON,
// the server's default, it is unknown for every row and so never true.
// The NULL may sit on either side. Only a comparison the predicate reads
// as its condition counts: the predicate itself, or an operand of AND,
// OR, or NOT in it, parenthesized or not. Once the change sets ANSI_NULLS
// OFF, under which = NULL matches a null, nothing after is judged.
func checkRequireIsNull(sql string, s *statement, ansiNullsOff bool, r *reporter) {
	if ansiNullsOff {
		return
	}
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
			case ast.BinOpAnd, ast.BinOpOr:
				visit(v.Left)
				visit(v.Right)
			case ast.BinOpEq, ast.BinOpNeq:
				if seen[v] || !isNull(v.Left) && !isNull(v.Right) {
					return
				}
				seen[v] = true
				fix := "IS NULL"
				if v.Op == ast.BinOpNeq {
					fix = "IS NOT NULL"
				}
				rng := rangeOf(v.Loc)
				text := "= NULL"
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

// setsAnsiNullsOff reports whether a statement is SET ANSI_NULLS OFF.
func setsAnsiNullsOff(n ast.Node) bool {
	v, ok := n.(*ast.SetOptionStmt)
	if !ok || !strings.EqualFold(v.Option, "ANSI_NULLS") {
		return false
	}
	switch val := v.Value.(type) {
	case *ast.ColumnRef:
		return strings.EqualFold(val.Column, "OFF")
	case *ast.Literal:
		return strings.EqualFold(val.Str, "OFF")
	}
	return false
}

// isNull reports whether an expression is the NULL literal, parenthesized
// or not.
func isNull(n ast.Node) bool {
	switch v := n.(type) {
	case *ast.Literal:
		return v.Type == ast.LitNull
	case *ast.ParenExpr:
		return isNull(v.Expr)
	}
	return false
}

// predicates lists the expressions of the statement that are evaluated
// for truth: WHERE and HAVING, join conditions, the conditions of a
// searched CASE, CHECK expressions, and the conditions of IF and WHILE.
// Subqueries are visited too, so their own predicates are listed.
func predicates(root ast.Node) []ast.Node {
	var roots []ast.Node
	add := func(n ast.Node) {
		if n != nil {
			roots = append(roots, n)
		}
	}
	ast.Inspect(root, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.SelectStmt:
			if v.WhereClause != nil {
				add(v.WhereClause)
			}
			if v.HavingClause != nil {
				add(v.HavingClause)
			}
		case *ast.UpdateStmt:
			if v.WhereClause != nil {
				add(v.WhereClause)
			}
		case *ast.DeleteStmt:
			if v.WhereClause != nil {
				add(v.WhereClause)
			}
		case *ast.JoinClause:
			if v.Condition != nil {
				add(v.Condition)
			}
		case *ast.CaseExpr:
			// A simple CASE compares its WHEN values with the argument;
			// only a searched CASE tests its WHEN expressions.
			if v.Arg == nil && v.WhenList != nil {
				for _, item := range v.WhenList.Items {
					if w, ok := item.(*ast.CaseWhen); ok && w.Condition != nil {
						add(w.Condition)
					}
				}
			}
		case *ast.ConstraintDef:
			if v.Type == ast.ConstraintCheck && v.Expr != nil {
				add(v.Expr)
			}
		case *ast.IfStmt:
			if v.Condition != nil {
				add(v.Condition)
			}
		case *ast.WhileStmt:
			if v.Condition != nil {
				add(v.Condition)
			}
		}
		return true
	})
	return roots
}
