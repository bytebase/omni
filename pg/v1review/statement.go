package v1review

import (
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// anchor returns a Loc as the byte range a statement rule's finding anchors.
// The V1 statement advisors took a node's line from where it starts alone,
// so a node with a start and no end anchors at its start, and only a node
// with no start anchors nowhere (0,0).
func anchor(loc ast.Loc) (int, int) {
	if loc.Start < 0 {
		return 0, 0
	}
	if loc.End < loc.Start {
		return loc.Start, loc.Start
	}
	return loc.Start, loc.End
}

// unwrapRaw returns the statement a RawStmt wraps. Bytebase's V1 advisors
// were handed the statement itself; the parser wraps each in a RawStmt.
func unwrapRaw(stmt ast.Node) ast.Node {
	if raw, ok := stmt.(*ast.RawStmt); ok {
		return raw.Stmt
	}
	return stmt
}

// alterTableCmds returns the AlterTableCmd items of an ALTER TABLE.
func alterTableCmds(alter *ast.AlterTableStmt) []*ast.AlterTableCmd {
	if alter.Cmds == nil {
		return nil
	}
	var cmds []*ast.AlterTableCmd
	for _, item := range alter.Cmds.Items {
		if cmd, ok := item.(*ast.AlterTableCmd); ok {
			cmds = append(cmds, cmd)
		}
	}
	return cmds
}

// selectWalk is the walk Bytebase's V1 SELECT advisors shared: from a
// top-level SELECT, or the query of a top-level INSERT ... SELECT, CREATE
// VIEW, or CREATE TABLE AS / CREATE MATERIALIZED VIEW, it visits each
// SELECT that is not a set operation, the subqueries of a SELECT before the
// SELECT itself. A set operation is not visited, only its branches. The
// subqueries walked are WITH queries, and the SubLinks and FROM subqueries
// reached through select-list, FROM, WHERE, HAVING, GROUP BY, ORDER BY,
// LIMIT and OFFSET expressions by way of function arguments, operators,
// AND/OR/NOT, COALESCE, CASE, IS NULL, casts and joins; a subquery under any
// other expression is not.
type selectWalk struct {
	// visit judges one SELECT; top is set only for the top-level SELECT
	// statement itself.
	visit func(sel *ast.SelectStmt, top bool)
}

func (w *selectWalk) statement(stmt ast.Node) {
	switch n := unwrapRaw(stmt).(type) {
	case *ast.SelectStmt:
		w.selectStmt(n, true)
	case *ast.InsertStmt:
		if sel, ok := n.SelectStmt.(*ast.SelectStmt); ok {
			w.selectStmt(sel, false)
		}
	case *ast.ViewStmt:
		if sel, ok := n.Query.(*ast.SelectStmt); ok {
			w.selectStmt(sel, false)
		}
	case *ast.CreateTableAsStmt:
		if sel, ok := n.Query.(*ast.SelectStmt); ok {
			w.selectStmt(sel, false)
		}
	}
}

func (w *selectWalk) selectStmt(sel *ast.SelectStmt, top bool) {
	if sel == nil {
		return
	}
	if sel.Op != ast.SETOP_NONE {
		w.selectStmt(sel.Larg, false)
		w.selectStmt(sel.Rarg, false)
		return
	}
	if sel.WithClause != nil && sel.WithClause.Ctes != nil {
		for _, item := range sel.WithClause.Ctes.Items {
			if cte, ok := item.(*ast.CommonTableExpr); ok {
				if sub, ok := cte.Ctequery.(*ast.SelectStmt); ok {
					w.selectStmt(sub, false)
				}
			}
		}
	}
	w.list(sel.TargetList)
	w.list(sel.FromClause)
	w.node(sel.WhereClause)
	w.node(sel.HavingClause)
	w.list(sel.GroupClause)
	w.list(sel.SortClause)
	w.node(sel.LimitCount)
	w.node(sel.LimitOffset)
	w.visit(sel, top)
}

func (w *selectWalk) list(list *ast.List) {
	if list == nil {
		return
	}
	for _, item := range list.Items {
		w.node(item)
	}
}

func (w *selectWalk) node(node ast.Node) {
	if node == nil {
		return
	}
	switch n := node.(type) {
	case *ast.SubLink:
		if sub, ok := n.Subselect.(*ast.SelectStmt); ok {
			w.selectStmt(sub, false)
		}
	case *ast.RangeSubselect:
		if sub, ok := n.Subquery.(*ast.SelectStmt); ok {
			w.selectStmt(sub, false)
		}
	case *ast.List:
		w.list(n)
	case *ast.ResTarget:
		w.node(n.Val)
	case *ast.FuncCall:
		w.list(n.Args)
	case *ast.A_Expr:
		w.node(n.Lexpr)
		w.node(n.Rexpr)
	case *ast.BoolExpr:
		w.list(n.Args)
	case *ast.CoalesceExpr:
		w.list(n.Args)
	case *ast.CaseExpr:
		w.node(n.Arg)
		w.list(n.Args)
		w.node(n.Defresult)
	case *ast.CaseWhen:
		w.node(n.Expr)
		w.node(n.Result)
	case *ast.NullTest:
		w.node(n.Arg)
	case *ast.TypeCast:
		w.node(n.Arg)
	case *ast.JoinExpr:
		w.node(n.Larg)
		w.node(n.Rarg)
		w.node(n.Quals)
	case *ast.SortBy:
		w.node(n.Node)
	}
}

// selectFinding returns a finding on a SELECT the walk visited: one on the
// top-level SELECT statement anchors nothing, since Bytebase's V1 advisors
// printed the statement's text for it, and one on a nested SELECT anchors
// that SELECT, whose text they printed.
func selectFinding(code v1review.Code, sel *ast.SelectStmt, top bool) v1review.Finding {
	if top {
		return v1review.Finding{Code: code}
	}
	start, end := anchor(sel.Loc)
	return v1review.Finding{Code: code, Start: start, End: end}
}
