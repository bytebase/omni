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
	case *ast.ValuesClause:
		return w.writeValuesClause(n)
	default:
		return fmt.Errorf("deparse: expected query expression, got %T", node)
	}
}

// writeResultScanStmt emits the result-pipe: <source> ->> <query>. The source
// node may be any deparsable statement (a non-deparsable source — e.g. SHOW,
// which the deparser does not render — surfaces as the underlying writeNode
// error rather than a panic).
func (w *writer) writeResultScanStmt(n *ast.ResultScanStmt) error {
	if err := w.writeNode(n.Source); err != nil {
		return err
	}
	w.buf.WriteString(" ->> ")
	return w.writeQueryExpr(n.Query)
}

// writeValuesClause emits a VALUES row source: VALUES (e, …), (e, …).
func (w *writer) writeValuesClause(n *ast.ValuesClause) error {
	w.buf.WriteString("VALUES ")
	for i, row := range n.Rows {
		if i > 0 {
			w.buf.WriteString(", ")
		}
		w.buf.WriteByte('(')
		for j, val := range row {
			if j > 0 {
				w.buf.WriteString(", ")
			}
			if err := writeExprNoLeadSpace(w, val); err != nil {
				return err
			}
		}
		w.buf.WriteByte(')')
	}
	return nil
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
		// REPLACE (expr AS col, ...)
		if len(t.Replace) > 0 {
			w.buf.WriteString(" REPLACE (")
			for i, r := range t.Replace {
				if i > 0 {
					w.buf.WriteString(", ")
				}
				if err := writeExprNoLeadSpace(w, r.Expr); err != nil {
					return err
				}
				w.buf.WriteString(" AS ")
				w.buf.WriteString(r.Col.String())
			}
			w.buf.WriteByte(')')
		}
		// RENAME (col AS alias, ...)
		if len(t.Rename) > 0 {
			w.buf.WriteString(" RENAME (")
			for i, r := range t.Rename {
				if i > 0 {
					w.buf.WriteString(", ")
				}
				w.buf.WriteString(r.Col.String())
				w.buf.WriteString(" AS ")
				w.buf.WriteString(r.Alias.String())
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
	switch {
	case n.Nested != nil:
		// Chained PIVOT/UNPIVOT: the source is itself a pivoted ref.
		if err := w.writeTableRef(n.Nested); err != nil {
			return err
		}
	case n.Name != nil:
		w.buf.WriteString(n.Name.String())
	case n.Subquery != nil:
		w.buf.WriteByte('(')
		if err := w.writeQueryExpr(n.Subquery); err != nil {
			return err
		}
		w.buf.WriteByte(')')
	case n.FuncCall != nil:
		if err := writeExprNoLeadSpace(w, n.FuncCall); err != nil {
			return err
		}
	case n.DollarN != nil:
		// $N result-set table reference (e.g. FROM $1).
		if err := writeExprNoLeadSpace(w, n.DollarN); err != nil {
			return err
		}
	}
	// PIVOT / UNPIVOT come before the alias position in the documented
	// clause order (each clause carries its own trailing alias).
	if n.Pivot != nil {
		w.buf.WriteByte(' ')
		if err := w.writePivotClause(n.Pivot); err != nil {
			return err
		}
	}
	if n.Unpivot != nil {
		w.buf.WriteByte(' ')
		if err := w.writeUnpivotClause(n.Unpivot); err != nil {
			return err
		}
	}
	if !n.Alias.IsEmpty() {
		w.buf.WriteString(" AS ")
		w.buf.WriteString(n.Alias.String())
	}
	// Derived column list: AS v (c1, c2).
	if len(n.Columns) > 0 {
		w.buf.WriteString(" (")
		for i, col := range n.Columns {
			if i > 0 {
				w.buf.WriteString(", ")
			}
			w.buf.WriteString(col.String())
		}
		w.buf.WriteByte(')')
	}
	return nil
}

// writePivotClause emits PIVOT (agg [AS a] FOR col IN (...) [DEFAULT ON NULL
// (expr)]) [AS alias].
func (w *writer) writePivotClause(n *ast.PivotClause) error {
	w.buf.WriteString("PIVOT (")
	if err := writeExprNoLeadSpace(w, n.Agg); err != nil {
		return err
	}
	if !n.AggAlias.IsEmpty() {
		w.buf.WriteString(" AS ")
		w.buf.WriteString(n.AggAlias.String())
	}
	w.buf.WriteString(" FOR ")
	if err := writeExprNoLeadSpace(w, n.ForColumn); err != nil {
		return err
	}
	w.buf.WriteString(" IN (")
	if n.In == nil {
		return fmt.Errorf("deparse: PIVOT clause without IN list")
	}
	switch n.In.Kind {
	case ast.PivotInAny:
		w.buf.WriteString("ANY")
		if len(n.In.OrderBy) > 0 {
			w.buf.WriteString(" ORDER BY ")
			for i, item := range n.In.OrderBy {
				if i > 0 {
					w.buf.WriteString(", ")
				}
				if err := w.writeOrderItem(item); err != nil {
					return err
				}
			}
		}
	case ast.PivotInSubquery:
		if err := w.writeQueryExpr(n.In.Subquery); err != nil {
			return err
		}
	default: // PivotInValues
		for i, val := range n.In.Values {
			if i > 0 {
				w.buf.WriteString(", ")
			}
			if err := writeExprNoLeadSpace(w, val.Value); err != nil {
				return err
			}
			if !val.Alias.IsEmpty() {
				w.buf.WriteString(" AS ")
				w.buf.WriteString(val.Alias.String())
			}
		}
	}
	w.buf.WriteByte(')')
	if n.DefaultVal != nil {
		w.buf.WriteString(" DEFAULT ON NULL (")
		if err := writeExprNoLeadSpace(w, n.DefaultVal); err != nil {
			return err
		}
		w.buf.WriteByte(')')
	}
	w.buf.WriteByte(')')
	if !n.Alias.IsEmpty() {
		w.buf.WriteString(" AS ")
		w.buf.WriteString(n.Alias.String())
	}
	return nil
}

// writeUnpivotClause emits UNPIVOT [INCLUDE|EXCLUDE NULLS] (val FOR name IN
// (col [AS a], ...)) [AS alias].
func (w *writer) writeUnpivotClause(n *ast.UnpivotClause) error {
	w.buf.WriteString("UNPIVOT ")
	switch n.NullsMode {
	case ast.UnpivotIncludeNulls:
		w.buf.WriteString("INCLUDE NULLS ")
	case ast.UnpivotExcludeNulls:
		w.buf.WriteString("EXCLUDE NULLS ")
	}
	w.buf.WriteByte('(')
	w.buf.WriteString(n.ValueColumn.String())
	w.buf.WriteString(" FOR ")
	w.buf.WriteString(n.NameColumn.String())
	w.buf.WriteString(" IN (")
	for i, col := range n.Columns {
		if i > 0 {
			w.buf.WriteString(", ")
		}
		w.buf.WriteString(col.Column.String())
		if !col.Alias.IsEmpty() {
			w.buf.WriteString(" AS ")
			w.buf.WriteString(col.Alias.String())
		}
	}
	w.buf.WriteString("))")
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
