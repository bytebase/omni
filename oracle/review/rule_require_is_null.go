package review

import (
	"strings"

	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/review"
)

// nullFixes maps each comparison operator that is never true against NULL
// to the test that says what it means: = to IS NULL, and the inequality
// operators, != <> ^= and PL/SQL's ~=, to IS NOT NULL.
var nullFixes = map[string]string{
	"=":  "IS NULL",
	"!=": "IS NOT NULL",
	"<>": "IS NOT NULL",
	"^=": "IS NOT NULL",
	"~=": "IS NOT NULL",
}

// checkRequireIsNull reports a comparison of the form x = NULL, x <> NULL,
// x != NULL, or x ^= NULL that a condition tests for truth, where it is
// unknown for every row and so never true. Oracle treats a zero-length
// string literal as NULL, so x = ” (or N”, or q'[]') is never true
// either, in SQL and in PL/SQL alike. The NULL may sit on either side, and
// in parentheses. Only a comparison the condition reads as its truth
// counts: the condition itself, or an operand of AND, OR, or NOT in it,
// parenthesized or not.
//
// The conditions read are those of WHERE, HAVING, QUALIFY, START WITH and
// CONNECT BY, join ON, MERGE ON and its WHERE clauses, a multitable
// INSERT's WHEN, a searched CASE's WHEN, a CHECK constraint of a table or
// a domain, a trigger's
// WHEN, a bitmap join index's WHERE, and PL/SQL's IF, ELSIF, WHILE, EXIT
// WHEN, CONTINUE WHEN, and searched CASE statement. A simple CASE's WHEN
// NULL, an IN list holding NULL, and a comparison used as a value are not
// judged.
//
// oracle: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/Nulls.html
func checkRequireIsNull(sql string, s *statement, r *reporter) {
	seen := make(map[*ast.BinaryExpr]bool)
	var visit func(n ast.Node)
	visit = func(n ast.Node) {
		switch v := n.(type) {
		case *ast.ParenExpr:
			visit(v.Expr)
		case *ast.BoolExpr:
			if v.Args != nil {
				for _, arg := range v.Args.Items {
					visit(arg)
				}
			}
		case *ast.UnaryExpr:
			if strings.EqualFold(v.Op, "NOT") {
				visit(v.Operand)
			}
		case *ast.BinaryExpr:
			fix, ok := nullFixes[v.Op]
			if !ok || seen[v] || !isNull(v.Left) && !isNull(v.Right) {
				return
			}
			seen[v] = true
			rng := comparisonRange(v)
			text := strings.Join(strings.Fields(sql[rng.Start:rng.End]), " ")
			r.report(review.RequireIsNull, s.index, rng, `"`+escapeLines(text)+`" is never true; use `+fix)
		}
	}
	for _, p := range conditions(s.node) {
		visit(p)
	}
}

// comparisonRange returns the range a comparison covers. The parser starts
// a binary expression's location at its operator, so the left operand's
// location widens it to the whole comparison.
func comparisonRange(v *ast.BinaryExpr) review.Range {
	loc := v.Loc
	if l := ast.NodeLoc(v.Left); l.Start >= 0 && l.End > l.Start && l.Start < loc.Start {
		loc.Start = l.Start
	}
	if r := ast.NodeLoc(v.Right); r.Start >= 0 && r.End > loc.End {
		loc.End = r.End
	}
	return rangeOf(loc)
}

// isNull reports whether an expression is the NULL literal or a
// zero-length string literal, which Oracle treats as NULL, parenthesized
// or not.
//
// oracle: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/Nulls.html
func isNull(n ast.Node) bool {
	switch v := n.(type) {
	case *ast.NullLiteral:
		return true
	case *ast.StringLiteral:
		return v.Val == ""
	case *ast.ParenExpr:
		return isNull(v.Expr)
	}
	return false
}

// conditions lists the expressions of the statement that are evaluated
// for truth. Subqueries and PL/SQL bodies are visited too, so their own
// conditions are listed.
func conditions(root ast.Node) []ast.Node {
	var out []ast.Node
	add := func(n ast.Node) {
		if n != nil {
			out = append(out, n)
		}
	}
	ast.Inspect(root, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.SelectStmt:
			add(v.WhereClause)
			add(v.HavingClause)
			add(v.QualifyClause)
		case *ast.HierarchicalClause:
			add(v.StartWith)
			add(v.ConnectBy)
		case *ast.ModelRulesClause:
			add(v.Until)
		case *ast.MatchRecognizeClause:
			if v.Definitions != nil {
				for _, item := range v.Definitions.Items {
					if t, ok := item.(*ast.ResTarget); ok {
						add(t.Expr)
					}
				}
			}
		case *ast.UpdateStmt:
			add(v.WhereClause)
		case *ast.DeleteStmt:
			add(v.WhereClause)
		case *ast.JoinClause:
			add(v.On)
		case *ast.MergeStmt:
			add(v.On)
		case *ast.MergeClause:
			add(v.Condition)
			add(v.UpdateWhere)
			add(v.DeleteWhere)
			add(v.InsertWhere)
		case *ast.InsertIntoClause:
			add(v.When)
		case *ast.CaseExpr:
			// A simple CASE compares its WHEN values with the argument;
			// only a searched CASE tests its WHEN expressions.
			if v.Arg == nil && v.Whens != nil {
				for _, item := range v.Whens.Items {
					if w, ok := item.(*ast.CaseWhen); ok {
						add(w.Condition)
					}
				}
			}
		case *ast.ColumnConstraint:
			if v.Type == ast.CONSTRAINT_CHECK {
				add(v.Expr)
			}
		case *ast.TableConstraint:
			if v.Type == ast.CONSTRAINT_CHECK {
				add(v.Expr)
			}
		case *ast.DomainConstraint:
			if strings.EqualFold(v.Type, "CHECK") {
				add(v.CheckExpr)
			}
		case *ast.CreateTriggerStmt:
			add(v.When)
		case *ast.CreateIndexStmt:
			add(v.Where)
		case *ast.PLSQLIf:
			add(v.Condition)
		case *ast.PLSQLElsIf:
			add(v.Condition)
		case *ast.PLSQLLoop:
			if v.Type == ast.LOOP_WHILE {
				add(v.Condition)
			}
		case *ast.PLSQLExit:
			add(v.Condition)
		case *ast.PLSQLContinue:
			add(v.Condition)
		case *ast.PLSQLCase:
			if v.Expr == nil {
				for _, w := range v.Whens {
					if w != nil {
						add(w.Expr)
					}
				}
			}
		}
		return true
	})
	return out
}
