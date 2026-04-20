package catalog

import (
	"strings"
)

// DeparseQuery converts an analyzed Query IR back to canonical SQL text.
func DeparseQuery(q *Query) string {
	if q == nil {
		return ""
	}
	if q.SetOp != SetOpNone {
		return deparseSetOpQuery(q)
	}
	return deparseSimpleQuery(q)
}

func deparseSimpleQuery(q *Query) string {
	var b strings.Builder

	if len(q.CTEList) > 0 {
		deparseCTEsQ(&b, q)
	}

	b.WriteString("select ")

	if q.Distinct {
		b.WriteString("distinct ")
	}

	deparseTargetListQ(&b, q)

	if len(q.RangeTable) > 0 && q.JoinTree != nil && len(q.JoinTree.FromList) > 0 {
		b.WriteString(" from ")
		deparseFromClauseQ(&b, q)
	}

	if q.JoinTree != nil && q.JoinTree.Quals != nil {
		b.WriteString(" where ")
		b.WriteString(DeparseAnalyzedExpr(q.JoinTree.Quals, q))
	}

	if len(q.GroupClause) > 0 {
		b.WriteString(" group by ")
		deparseGroupByQ(&b, q)
	}

	if q.HavingQual != nil {
		b.WriteString(" having ")
		b.WriteString(DeparseAnalyzedExpr(q.HavingQual, q))
	}

	if len(q.SortClause) > 0 {
		b.WriteString(" order by ")
		deparseOrderByQ(&b, q)
	}

	if q.LimitCount != nil {
		b.WriteString(" limit ")
		b.WriteString(DeparseAnalyzedExpr(q.LimitCount, q))
		if q.LimitOffset != nil {
			b.WriteString(" offset ")
			b.WriteString(DeparseAnalyzedExpr(q.LimitOffset, q))
		}
	}

	return b.String()
}

func deparseCTEsQ(b *strings.Builder, q *Query) {
	b.WriteString("with ")
	if q.IsRecursive {
		b.WriteString("recursive ")
	}
	for i, cte := range q.CTEList {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(backtickIDQ(cte.Name))
		if len(cte.ColumnNames) > 0 {
			b.WriteByte('(')
			for j, col := range cte.ColumnNames {
				if j > 0 {
					b.WriteString(", ")
				}
				b.WriteString(backtickIDQ(col))
			}
			b.WriteByte(')')
		}
		b.WriteString(" as (")
		b.WriteString(DeparseQuery(cte.Query))
		b.WriteString(") ")
	}
}

func deparseTargetListQ(b *strings.Builder, q *Query) {
	first := true
	for _, te := range q.TargetList {
		if te.ResJunk {
			continue
		}
		if !first {
			b.WriteString(", ")
		}
		first = false

		exprText := DeparseAnalyzedExpr(te.Expr, q)
		b.WriteString(exprText)

		// Omit alias when expression is a simple column ref whose name matches ResName.
		needAlias := true
		if v, ok := te.Expr.(*VarExprQ); ok {
			if v.RangeIdx >= 0 && v.RangeIdx < len(q.RangeTable) {
				rte := q.RangeTable[v.RangeIdx]
				if v.AttNum >= 1 && v.AttNum <= len(rte.ColNames) {
					if strings.EqualFold(rte.ColNames[v.AttNum-1], te.ResName) {
						needAlias = false
					}
				}
			}
		}
		if needAlias {
			b.WriteString(" as ")
			b.WriteString(backtickIDQ(te.ResName))
		}
	}
}

func deparseFromClauseQ(b *strings.Builder, q *Query) {
	for i, node := range q.JoinTree.FromList {
		if i > 0 {
			b.WriteString(", ")
		}
		deparseJoinNodeQ(b, node, q)
	}
}

func deparseJoinNodeQ(b *strings.Builder, node JoinNode, q *Query) {
	switch n := node.(type) {
	case *RangeTableRefQ:
		deparseRTERefQ(b, n.RTIndex, q)
	case *JoinExprNodeQ:
		deparseJoinExprNodeQ(b, n, q)
	}
}

func deparseRTERefQ(b *strings.Builder, idx int, q *Query) {
	if idx < 0 || idx >= len(q.RangeTable) {
		b.WriteString("?")
		return
	}
	rte := q.RangeTable[idx]
	switch rte.Kind {
	case RTERelation:
		if rte.DBName != "" {
			b.WriteString(backtickIDQ(rte.DBName))
			b.WriteByte('.')
		}
		b.WriteString(backtickIDQ(rte.TableName))
		if rte.Alias != "" && rte.Alias != rte.TableName {
			b.WriteString(" as ")
			b.WriteString(backtickIDQ(rte.Alias))
		}
	case RTESubquery:
		b.WriteByte('(')
		b.WriteString(DeparseQuery(rte.Subquery))
		b.WriteString(") as ")
		b.WriteString(backtickIDQ(rte.ERef))
	case RTECTE:
		b.WriteString(backtickIDQ(rte.CTEName))
		if rte.Alias != "" && rte.Alias != rte.CTEName {
			b.WriteString(" as ")
			b.WriteString(backtickIDQ(rte.Alias))
		}
	case RTEJoin:
		// Synthetic; actual join structure is in JoinExprNodeQ.
	case RTEFunction:
		b.WriteString("/* function-in-FROM */")
	}
}

func deparseJoinExprNodeQ(b *strings.Builder, j *JoinExprNodeQ, q *Query) {
	deparseJoinNodeQ(b, j.Left, q)
	b.WriteByte(' ')

	if j.Natural {
		b.WriteString("natural ")
	}

	switch j.JoinType {
	case JoinLeft:
		b.WriteString("left join ")
	case JoinRight:
		b.WriteString("right join ")
	case JoinCross:
		b.WriteString("cross join ")
	case JoinStraight:
		b.WriteString("straight_join ")
	default: // JoinInner
		b.WriteString("join ")
	}

	deparseJoinNodeQ(b, j.Right, q)

	if j.Quals != nil {
		b.WriteString(" on ")
		b.WriteString(DeparseAnalyzedExpr(j.Quals, q))
	}

	if len(j.UsingClause) > 0 && j.Quals == nil {
		b.WriteString(" using (")
		for i, col := range j.UsingClause {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(backtickIDQ(col))
		}
		b.WriteByte(')')
	}
}

func deparseGroupByQ(b *strings.Builder, q *Query) {
	for i, sc := range q.GroupClause {
		if i > 0 {
			b.WriteString(", ")
		}
		if sc.TargetIdx >= 1 && sc.TargetIdx <= len(q.TargetList) {
			b.WriteString(DeparseAnalyzedExpr(q.TargetList[sc.TargetIdx-1].Expr, q))
		}
	}
}

func deparseOrderByQ(b *strings.Builder, q *Query) {
	for i, sc := range q.SortClause {
		if i > 0 {
			b.WriteString(", ")
		}
		if sc.TargetIdx >= 1 && sc.TargetIdx <= len(q.TargetList) {
			b.WriteString(DeparseAnalyzedExpr(q.TargetList[sc.TargetIdx-1].Expr, q))
		}
		if sc.Descending {
			b.WriteString(" desc")
		}
	}
}

func deparseSetOpQuery(q *Query) string {
	var b strings.Builder

	if len(q.CTEList) > 0 {
		deparseCTEsQ(&b, q)
	}

	b.WriteString(DeparseQuery(q.LArg))
	b.WriteByte(' ')

	switch q.SetOp {
	case SetOpUnion:
		b.WriteString("union")
	case SetOpIntersect:
		b.WriteString("intersect")
	case SetOpExcept:
		b.WriteString("except")
	}
	if q.AllSetOp {
		b.WriteString(" all")
	}

	b.WriteByte(' ')
	b.WriteString(DeparseQuery(q.RArg))

	if len(q.SortClause) > 0 {
		b.WriteString(" order by ")
		deparseOrderByQ(&b, q)
	}

	if q.LimitCount != nil {
		b.WriteString(" limit ")
		b.WriteString(DeparseAnalyzedExpr(q.LimitCount, q))
		if q.LimitOffset != nil {
			b.WriteString(" offset ")
			b.WriteString(DeparseAnalyzedExpr(q.LimitOffset, q))
		}
	}

	return b.String()
}

// backtickIDQ wraps an identifier in backticks.
func backtickIDQ(name string) string {
	return "`" + strings.ReplaceAll(name, "`", "``") + "`"
}
