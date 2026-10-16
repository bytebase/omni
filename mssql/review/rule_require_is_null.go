package review

import (
	"slices"
	"strings"

	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/review"
)

// checkRequireIsNull reports a comparison of the form x = NULL, x <> NULL,
// or x != NULL that a predicate tests for truth, where with ANSI_NULLS ON,
// the server's default, it is unknown for every row and so never true.
// The NULL may sit on either side. Only a comparison the predicate reads
// as its condition counts: the predicate itself, or an operand of AND,
// OR, or NOT in it, parenthesized or not. While the change has set
// ANSI_NULLS OFF, under which = NULL matches a null, or may have, nothing
// is judged, and neither is a statement that sets it within.
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
					text = oneLine(sql[rng.Start:rng.End])
				}
				r.report(review.RequireIsNull, s.index, rng, `"`+text+`" is never true; use `+fix)
			}
		}
	}
	for _, p := range predicates(s.node) {
		visit(p)
	}
}

// ansiNulls is the ANSI_NULLS setting as the change runs: known is false
// once a SET may or may not have run.
type ansiNulls struct {
	known bool
	off   bool
}

// after returns the setting after a top-level statement runs, and whether
// the statement may set it: a SET itself, the statements of a BEGIN ...
// END block in the order they run, and a SET under IF, WHILE, or TRY,
// after which the setting is not known. A routine body runs under the
// setting of its creation, so a SET in it does not count.
func (a ansiNulls) after(n ast.Node) (ansiNulls, bool) {
	if set, ok := ansiNullsSet(n); ok {
		return ansiNulls{known: true, off: !set}, true
	}
	switch v := n.(type) {
	case *ast.BeginEndStmt:
		sets := false
		if v.Stmts != nil {
			for _, item := range v.Stmts.Items {
				var s bool
				a, s = a.after(item)
				sets = sets || s
			}
		}
		return a, sets
	case *ast.CreateProcedureStmt, *ast.CreateFunctionStmt, *ast.CreateTriggerStmt:
		return a, false
	}
	sets := false
	ast.Inspect(n, func(m ast.Node) bool {
		switch m.(type) {
		case *ast.CreateProcedureStmt, *ast.CreateFunctionStmt, *ast.CreateTriggerStmt:
			return false
		}
		if _, ok := ansiNullsSet(m); ok {
			sets = true
		}
		return !sets
	})
	if sets {
		return ansiNulls{}, true
	}
	return a, false
}

// ansiNullsSet returns the ANSI_NULLS setting a statement makes: SET
// ANSI_NULLS, or SET ANSI_DEFAULTS, which sets it too, alone or in a list
// of options sharing one ON or OFF. ok is false for another statement.
func ansiNullsSet(n ast.Node) (on, ok bool) {
	v, isSet := n.(*ast.SetOptionStmt)
	if !isSet || !slices.ContainsFunc(v.Options, func(name string) bool {
		return strings.EqualFold(name, "ANSI_NULLS") || strings.EqualFold(name, "ANSI_DEFAULTS")
	}) {
		return false, false
	}
	var value string
	switch val := v.Value.(type) {
	case *ast.ColumnRef:
		value = val.Column
	case *ast.Literal:
		value = val.Str
	}
	switch {
	case strings.EqualFold(value, "ON"):
		return true, true
	case strings.EqualFold(value, "OFF"):
		return false, true
	}
	return false, false
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
// for truth: WHERE and HAVING, join conditions, MERGE conditions, the
// conditions of a searched CASE and of IIF, CHECK expressions, and the
// conditions of IF and WHILE.
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
		case *ast.MergeStmt:
			if v.OnCondition != nil {
				add(v.OnCondition)
			}
		case *ast.MergeWhenClause:
			if v.Condition != nil {
				add(v.Condition)
			}
		case *ast.IifExpr:
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
