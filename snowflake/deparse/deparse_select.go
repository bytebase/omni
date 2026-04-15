package deparse

import (
	"fmt"

	"github.com/bytebase/omni/snowflake/ast"
)

// writeQueryExpr writes a query expression (SelectStmt or SetOperationStmt).
// This is the recursive entry point used by subqueries as well.
func (w *writer) writeQueryExpr(node ast.Node) error {
	switch n := node.(type) {
	case *ast.SelectStmt:
		return w.writeSelectStmt(n)
	case *ast.SetOperationStmt:
		return w.writeSetOperationStmt(n)
	default:
		return fmt.Errorf("deparse: expected query expression, got %T", node)
	}
}

// writeCTEs writes the WITH clause (all CTEs).
func (w *writer) writeCTEs(ctes []*ast.CTE) error {
	if len(ctes) == 0 {
		return nil
	}
	// Check if any CTE is recursive.
	hasRecursive := false
	for _, cte := range ctes {
		if cte.Recursive {
			hasRecursive = true
			break
		}
	}
	if hasRecursive {
		w.buf.WriteString("WITH RECURSIVE ")
	} else {
		w.buf.WriteString("WITH ")
	}
	for i, cte := range ctes {
		if i > 0 {
			w.buf.WriteString(", ")
		}
		if err := w.writeCTE(cte); err != nil {
			return err
		}
	}
	w.buf.WriteByte(' ')
	return nil
}

func (w *writer) writeCTE(cte *ast.CTE) error {
	w.buf.WriteString(cte.Name.String())
	if len(cte.Columns) > 0 {
		w.buf.WriteByte('(')
		for i, col := range cte.Columns {
			if i > 0 {
				w.buf.WriteString(", ")
			}
			w.buf.WriteString(col.String())
		}
		w.buf.WriteByte(')')
	}
	w.buf.WriteString(" AS (")
	if err := w.writeQueryExpr(cte.Query); err != nil {
		return err
	}
	w.buf.WriteByte(')')
	return nil
}

// writeSelectStmt writes a complete SELECT statement.
func (w *writer) writeSelectStmt(n *ast.SelectStmt) error {
	// WITH clause
	if err := w.writeCTEs(n.With); err != nil {
		return err
	}

	w.buf.WriteString("SELECT")

	// TOP n (before DISTINCT/ALL)
	if n.Top != nil {
		w.buf.WriteString(" TOP ")
		if err := writeExprNoLeadSpace(w, n.Top); err != nil {
			return err
		}
	}

	if n.Distinct {
		w.buf.WriteString(" DISTINCT")
	} else if n.All {
		w.buf.WriteString(" ALL")
	}

	// SELECT list
	if len(n.Targets) == 0 {
		w.buf.WriteString(" *")
	} else {
		for i, target := range n.Targets {
			if i == 0 {
				w.buf.WriteByte(' ')
			} else {
				w.buf.WriteString(", ")
			}
			if err := w.writeSelectTarget(target); err != nil {
				return err
			}
		}
	}

	// FROM clause
	if len(n.From) > 0 {
		w.buf.WriteString(" FROM ")
		for i, src := range n.From {
			if i > 0 {
				w.buf.WriteString(", ")
			}
			if err := w.writeFromSource(src); err != nil {
				return err
			}
		}
	}

	// WHERE
	if n.Where != nil {
		w.buf.WriteString(" WHERE ")
		if err := writeExprNoLeadSpace(w, n.Where); err != nil {
			return err
		}
	}

	// GROUP BY
	if n.GroupBy != nil {
		if err := w.writeGroupByClause(n.GroupBy); err != nil {
			return err
		}
	}

	// HAVING
	if n.Having != nil {
		w.buf.WriteString(" HAVING ")
		if err := writeExprNoLeadSpace(w, n.Having); err != nil {
			return err
		}
	}

	// QUALIFY (Snowflake-specific)
	if n.Qualify != nil {
		w.buf.WriteString(" QUALIFY ")
		if err := writeExprNoLeadSpace(w, n.Qualify); err != nil {
			return err
		}
	}

	// ORDER BY
	if len(n.OrderBy) > 0 {
		w.buf.WriteString(" ORDER BY ")
		for i, item := range n.OrderBy {
			if i > 0 {
				w.buf.WriteString(", ")
			}
			if err := w.writeOrderItem(item); err != nil {
				return err
			}
		}
	}

	// LIMIT
	if n.Limit != nil {
		w.buf.WriteString(" LIMIT ")
		if err := writeExprNoLeadSpace(w, n.Limit); err != nil {
			return err
		}
	}

	// OFFSET
	if n.Offset != nil {
		w.buf.WriteString(" OFFSET ")
		if err := writeExprNoLeadSpace(w, n.Offset); err != nil {
			return err
		}
	}

	// FETCH FIRST/NEXT n ROWS ONLY
	if n.Fetch != nil {
		w.buf.WriteString(" FETCH FIRST ")
		if err := writeExprNoLeadSpace(w, n.Fetch.Count); err != nil {
			return err
		}
		w.buf.WriteString(" ROWS ONLY")
	}

	return nil
}

func (w *writer) writeSelectTarget(t *ast.SelectTarget) error {
	if t.Star {
		// Expr is always a *StarExpr when Star is true.
		if se, ok := t.Expr.(*ast.StarExpr); ok && se != nil && se.Qualifier != nil {
			// qualifier.*
			w.writeObjectNameNoSpace(se.Qualifier)
			w.buf.WriteString(".*")
		} else {
			w.buf.WriteByte('*')
		}
		// EXCLUDE (col, ...)
		if len(t.Exclude) > 0 {
			w.buf.WriteString(" EXCLUDE (")
			for i, col := range t.Exclude {
				if i > 0 {
					w.buf.WriteString(", ")
				}
				w.buf.WriteString(col.String())
			}
			w.buf.WriteByte(')')
		}
		return nil
	}
	if err := writeExprNoLeadSpace(w, t.Expr); err != nil {
		return err
	}
	if !t.Alias.IsEmpty() {
		w.buf.WriteString(" AS ")
		w.buf.WriteString(t.Alias.String())
	}
	// EXCLUDE only makes sense on *, but guard anyway
	if len(t.Exclude) > 0 {
		w.buf.WriteString(" EXCLUDE (")
		for i, col := range t.Exclude {
			if i > 0 {
				w.buf.WriteString(", ")
			}
			w.buf.WriteString(col.String())
		}
		w.buf.WriteByte(')')
	}
	return nil
}

func (w *writer) writeFromSource(node ast.Node) error {
	switch n := node.(type) {
	case *ast.TableRef:
		return w.writeTableRef(n)
	case *ast.JoinExpr:
		return w.writeJoinExpr(n)
	default:
		return fmt.Errorf("deparse: unsupported FROM source node type %T", node)
	}
}

func (w *writer) writeTableRef(n *ast.TableRef) error {
	if n.Lateral {
		w.buf.WriteString("LATERAL ")
	}
	if n.Name != nil {
		w.buf.WriteString(n.Name.String())
	} else if n.Subquery != nil {
		w.buf.WriteByte('(')
		if err := w.writeQueryExpr(n.Subquery); err != nil {
			return err
		}
		w.buf.WriteByte(')')
	} else if n.FuncCall != nil {
		if err := writeExprNoLeadSpace(w, n.FuncCall); err != nil {
			return err
		}
	}
	if !n.Alias.IsEmpty() {
		w.buf.WriteString(" AS ")
		w.buf.WriteString(n.Alias.String())
	}
	return nil
}

func (w *writer) writeJoinExpr(n *ast.JoinExpr) error {
	if err := w.writeFromSource(n.Left); err != nil {
		return err
	}
	// JOIN keyword
	switch n.Type {
	case ast.JoinInner:
		if n.Natural {
			w.buf.WriteString(" NATURAL JOIN ")
		} else {
			w.buf.WriteString(" JOIN ")
		}
	case ast.JoinLeft:
		if n.Natural {
			w.buf.WriteString(" NATURAL LEFT JOIN ")
		} else {
			w.buf.WriteString(" LEFT JOIN ")
		}
	case ast.JoinRight:
		if n.Natural {
			w.buf.WriteString(" NATURAL RIGHT JOIN ")
		} else {
			w.buf.WriteString(" RIGHT JOIN ")
		}
	case ast.JoinFull:
		if n.Natural {
			w.buf.WriteString(" NATURAL FULL JOIN ")
		} else {
			w.buf.WriteString(" FULL JOIN ")
		}
	case ast.JoinCross:
		w.buf.WriteString(" CROSS JOIN ")
	case ast.JoinAsof:
		w.buf.WriteString(" ASOF JOIN ")
	}
	if err := w.writeFromSource(n.Right); err != nil {
		return err
	}
	if n.On != nil {
		w.buf.WriteString(" ON ")
		if err := writeExprNoLeadSpace(w, n.On); err != nil {
			return err
		}
	} else if len(n.Using) > 0 {
		w.buf.WriteString(" USING (")
		for i, col := range n.Using {
			if i > 0 {
				w.buf.WriteString(", ")
			}
			w.buf.WriteString(col.String())
		}
		w.buf.WriteByte(')')
	}
	if n.MatchCondition != nil {
		w.buf.WriteString(" MATCH_CONDITION (")
		if err := writeExprNoLeadSpace(w, n.MatchCondition); err != nil {
			return err
		}
		w.buf.WriteByte(')')
	}
	return nil
}

func (w *writer) writeGroupByClause(gb *ast.GroupByClause) error {
	switch gb.Kind {
	case ast.GroupByAll:
		w.buf.WriteString(" GROUP BY ALL")
		return nil
	case ast.GroupByCube:
		w.buf.WriteString(" GROUP BY CUBE (")
	case ast.GroupByRollup:
		w.buf.WriteString(" GROUP BY ROLLUP (")
	case ast.GroupByGroupingSets:
		w.buf.WriteString(" GROUP BY GROUPING SETS (")
	default:
		w.buf.WriteString(" GROUP BY ")
		for i, item := range gb.Items {
			if i > 0 {
				w.buf.WriteString(", ")
			}
			if err := writeExprNoLeadSpace(w, item); err != nil {
				return err
			}
		}
		return nil
	}
	// For CUBE/ROLLUP/GROUPING SETS, wrap items in parens
	for i, item := range gb.Items {
		if i > 0 {
			w.buf.WriteString(", ")
		}
		if err := writeExprNoLeadSpace(w, item); err != nil {
			return err
		}
	}
	w.buf.WriteByte(')')
	return nil
}

// writeSetOperationStmt writes a UNION/EXCEPT/INTERSECT query.
func (w *writer) writeSetOperationStmt(n *ast.SetOperationStmt) error {
	if err := w.writeQueryExpr(n.Left); err != nil {
		return err
	}
	switch n.Op {
	case ast.SetOpUnion:
		if n.All {
			w.buf.WriteString(" UNION ALL")
		} else {
			w.buf.WriteString(" UNION")
		}
	case ast.SetOpExcept:
		if n.All {
			w.buf.WriteString(" EXCEPT ALL")
		} else {
			w.buf.WriteString(" EXCEPT")
		}
	case ast.SetOpIntersect:
		if n.All {
			w.buf.WriteString(" INTERSECT ALL")
		} else {
			w.buf.WriteString(" INTERSECT")
		}
	}
	if n.ByName {
		w.buf.WriteString(" BY NAME")
	}
	w.buf.WriteByte(' ')
	return w.writeQueryExpr(n.Right)
}
