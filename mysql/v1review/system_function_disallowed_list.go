package v1review

import (
	"slices"
	"strings"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// SystemFunctionDisallowedList reports each call, in top-level statements,
// of a function in the disallowed list, compared without regard to letter
// case and without its schema qualifier. Statements nested in a routine body
// are not judged.
//
// Only some expressions are searched, each one whole, subqueries included:
// a SELECT's select list, WHERE, GROUP BY, HAVING and ORDER BY, and those of
// its CTEs, set operands and parenthesized body, but not its FROM clause or
// LIMIT; an INSERT's VALUES rows and its SELECT, but not ON DUPLICATE KEY
// UPDATE; an UPDATE's SET values and WHERE, but not its ORDER BY; a
// DELETE's WHERE; and the DEFAULT and ON UPDATE expressions of the columns a
// CREATE TABLE defines and an ALTER TABLE's ADD, MODIFY or CHANGE COLUMN
// gives, plus an ALTER COLUMN ... SET DEFAULT expression. A statement under
// EXPLAIN, a CREATE TABLE's AS SELECT and a generated column are not
// searched.
//
// Name is the function's name as the parser holds it, upper-cased; each
// finding anchors the call.
// Ported from Bytebase's V1 rule SYSTEM_FUNCTION_DISALLOWED_LIST.
func SystemFunctionDisallowedList(stmts []ast.Node, disallowed []string) []v1review.Finding {
	c := &functionCalls{}
	for _, fn := range disallowed {
		c.disallowed = append(c.disallowed, strings.ToUpper(fn))
	}
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.SelectStmt:
			c.selectStmt(n)
		case *ast.InsertStmt:
			for _, row := range n.Values {
				for _, expr := range row {
					c.expr(expr)
				}
			}
			if n.Select != nil {
				c.selectStmt(n.Select)
			}
		case *ast.UpdateStmt:
			for _, assign := range n.SetList {
				if assign != nil {
					c.expr(assign.Value)
				}
			}
			c.expr(n.Where)
		case *ast.DeleteStmt:
			c.expr(n.Where)
		case *ast.CreateTableStmt:
			for _, col := range n.Columns {
				if col != nil {
					c.expr(col.DefaultValue)
					c.expr(col.OnUpdate)
				}
			}
		case *ast.AlterTableStmt:
			for _, cmd := range n.Commands {
				if cmd == nil {
					continue
				}
				for _, col := range commandColumns(cmd) {
					c.expr(col.DefaultValue)
					c.expr(col.OnUpdate)
				}
				c.expr(cmd.DefaultExpr)
			}
		}
	}
	return c.findings
}

// functionCalls collects the calls of disallowed functions, the list
// upper-cased.
type functionCalls struct {
	disallowed []string
	findings   []v1review.Finding
}

func (c *functionCalls) selectStmt(n *ast.SelectStmt) {
	if n == nil {
		return
	}
	for _, cte := range n.CTEs {
		if cte != nil {
			c.selectStmt(cte.Select)
		}
	}
	if n.ParenSource != nil {
		c.selectStmt(n.ParenSource)
	}
	if n.SetOp != ast.SetOpNone {
		c.selectStmt(n.Left)
		c.selectStmt(n.Right)
	}
	for _, expr := range n.TargetList {
		c.expr(expr)
	}
	c.expr(n.Where)
	for _, expr := range n.GroupBy {
		c.expr(expr)
	}
	c.expr(n.Having)
	for _, item := range n.OrderBy {
		if item != nil {
			c.expr(item.Expr)
		}
	}
}

func (c *functionCalls) expr(expr ast.ExprNode) {
	if expr == nil {
		return
	}
	ast.Inspect(expr, func(node ast.Node) bool {
		if fn, ok := node.(*ast.FuncCallExpr); ok && slices.Contains(c.disallowed, strings.ToUpper(fn.Name)) {
			start, end := span(fn.Loc)
			c.findings = append(c.findings, v1review.Finding{Code: v1review.FunctionDisallowed, Name: fn.Name, Start: start, End: end})
		}
		return true
	})
}
