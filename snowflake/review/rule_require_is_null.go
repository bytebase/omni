package review

import (
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/snowflake/ast"
)

// checkRequireIsNull reports a comparison of the form x = NULL, x <> NULL,
// or x != NULL that a predicate tests for truth. Snowflake's = treats a NULL
// as an unknown value, so the comparison is NULL for every row and never
// true. The NULL may sit on either side, parenthesized or not. Only a
// comparison the predicate reads as its condition counts: the predicate
// itself, or an operand of AND, OR, or NOT in it, parenthesized or not.
//
// sql: https://docs.snowflake.com/en/sql-reference/functions/equal_null
func checkRequireIsNull(sql string, s *statement, r *reporter) {
	seen := make(map[*ast.BinaryExpr]bool)
	var visit func(n ast.Node)
	visit = func(n ast.Node) {
		switch v := n.(type) {
		case *ast.ParenExpr:
			visit(v.Expr)
		case *ast.UnaryExpr:
			if v.Op == ast.UnaryNot {
				visit(v.Expr)
			}
		case *ast.BinaryExpr:
			switch v.Op {
			case ast.BinAnd, ast.BinOr:
				visit(v.Left)
				visit(v.Right)
			case ast.BinEq, ast.BinNe:
				if seen[v] || !isNull(v.Left) && !isNull(v.Right) {
					return
				}
				seen[v] = true
				fix := "IS NULL"
				if v.Op == ast.BinNe {
					fix = "IS NOT NULL"
				}
				rng := rangeOf(span(v))
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

// isNull reports whether an expression is the NULL literal, parenthesized
// or not.
func isNull(n ast.Node) bool {
	switch v := n.(type) {
	case *ast.Literal:
		return v.Kind == ast.LitNull
	case *ast.ParenExpr:
		return isNull(v.Expr)
	}
	return false
}

// span is the range of an expression's text: from the first byte any node
// in it covers to the last. The parser starts a binary expression's
// location at its operator, so its own location leaves the left operand
// out.
func span(n ast.Node) ast.Loc {
	out := ast.Loc{Start: -1, End: -1}
	ast.Inspect(n, func(m ast.Node) bool {
		loc := ast.NodeLoc(m)
		if loc.Start >= 0 && loc.End > loc.Start {
			if out.Start < 0 || loc.Start < out.Start {
				out.Start = loc.Start
			}
			if loc.End > out.End {
				out.End = loc.End
			}
		}
		return true
	})
	return out
}

// predicates lists the expressions of the statement that are evaluated for
// truth: WHERE, HAVING, QUALIFY, START WITH, and CONNECT BY; join
// conditions; MERGE's ON and its WHEN ... AND conditions; the WHEN
// conditions of INSERT ALL and INSERT FIRST; the conditions of a searched
// CASE and of IFF; and the conditions of Snowflake Scripting's IF, ELSEIF,
// WHILE, REPEAT ... UNTIL, and searched CASE. Subqueries and nested
// statements are visited too, so their own predicates are listed.
func predicates(root ast.Node) []ast.Node {
	var roots []ast.Node
	add := func(n ast.Node) {
		if n != nil {
			roots = append(roots, n)
		}
	}
	ast.Inspect(root, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.CreatePolicyStmt:
			// A row access policy's body decides, for truth, which rows show.
			if v.Kind == ast.PolicyRowAccess {
				add(v.Body)
			}
		case *ast.AlterPolicyStmt:
			if v.Kind == ast.PolicyRowAccess {
				add(v.Body)
			}
		case *ast.CreateDynamicTableStmt:
			add(v.ImmutableWhere)
		case *ast.CreateTaskStmt:
			add(v.When)
		case *ast.AlterTaskStmt:
			add(v.When)
		case *ast.MatchDefine:
			add(v.Cond)
		case *ast.SelectStmt:
			add(v.Where)
			add(v.Having)
			add(v.Qualify)
			add(v.StartWith)
			for _, c := range v.ConnectBy {
				add(c)
			}
		case *ast.UpdateStmt:
			add(v.Where)
		case *ast.DeleteStmt:
			add(v.Where)
		case *ast.JoinExpr:
			add(v.On)
		case *ast.MergeStmt:
			add(v.On)
			for _, w := range v.Whens {
				if w != nil {
					add(w.AndCond)
				}
			}
		case *ast.InsertMultiStmt:
			for _, b := range v.Branches {
				if b != nil {
					add(b.When)
				}
			}
		case *ast.CaseExpr:
			// A simple CASE compares its WHEN values with the operand;
			// only a searched CASE tests its WHEN expressions.
			if v.Kind == ast.CaseSearched {
				for _, w := range v.Whens {
					if w != nil {
						add(w.Cond)
					}
				}
			}
		case *ast.IffExpr:
			add(v.Cond)
		case *ast.ScriptIfBranch:
			add(v.Cond)
		case *ast.ScriptWhileStmt:
			add(v.Cond)
		case *ast.ScriptRepeatStmt:
			add(v.Cond)
		case *ast.ScriptCaseStmt:
			if v.Operand == nil {
				for _, item := range v.Whens {
					if w, ok := item.(*ast.ScriptCaseWhen); ok {
						add(w.Match)
					}
				}
			}
		}
		return true
	})
	return roots
}
