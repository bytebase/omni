package catalog

import (
	"strings"
)

// DeparseQuery converts an analyzed Query IR back to canonical SQL text.
func DeparseQuery(q *Query) string {
	var b strings.Builder
	writeQueryQ(&b, q)
	return b.String()
}

// writeQueryQ writes a query into the builder of the whole call (see
// writeAnalyzedExpr); a UNION ALL chain nests one Query per arm.
func writeQueryQ(b *strings.Builder, q *Query) {
	if q == nil {
		return
	}
	if q.SetOp != SetOpNone {
		deparseSetOpQuery(b, q)
		return
	}
	deparseSimpleQuery(b, q)
}

func deparseSimpleQuery(b *strings.Builder, q *Query) {
	if len(q.CTEList) > 0 {
		deparseCTEsQ(b, q)
	}

	b.WriteString("select ")

	if q.Distinct {
		b.WriteString("distinct ")
	}

	deparseTargetListQ(b, q)

	if len(q.RangeTable) > 0 && q.JoinTree != nil && len(q.JoinTree.FromList) > 0 {
		b.WriteString(" from ")
		deparseFromClauseQ(b, q)
	}

	if q.JoinTree != nil && q.JoinTree.Quals != nil {
		b.WriteString(" where ")
		writeAnalyzedExpr(b, q.JoinTree.Quals, q)
	}

	if len(q.GroupClause) > 0 {
		b.WriteString(" group by ")
		deparseGroupByQ(b, q)
	}

	if q.HavingQual != nil {
		b.WriteString(" having ")
		writeAnalyzedExpr(b, q.HavingQual, q)
	}

	if len(q.SortClause) > 0 {
		b.WriteString(" order by ")
		deparseOrderByQ(b, q)
	}

	if q.LimitCount != nil {
		b.WriteString(" limit ")
		writeAnalyzedExpr(b, q.LimitCount, q)
		if q.LimitOffset != nil {
			b.WriteString(" offset ")
			writeAnalyzedExpr(b, q.LimitOffset, q)
		}
	}
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
		writeQueryQ(b, cte.Query)
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

		writeAnalyzedExpr(b, te.Expr, q)

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
		writeQueryQ(b, rte.Subquery)
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
		writeAnalyzedExpr(b, j.Quals, q)
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
			writeAnalyzedExpr(b, q.TargetList[sc.TargetIdx-1].Expr, q)
		}
	}
}

func deparseOrderByQ(b *strings.Builder, q *Query) {
	for i, sc := range q.SortClause {
		if i > 0 {
			b.WriteString(", ")
		}
		if sc.TargetIdx >= 1 && sc.TargetIdx <= len(q.TargetList) {
			writeAnalyzedExpr(b, q.TargetList[sc.TargetIdx-1].Expr, q)
		}
		if sc.Descending {
			b.WriteString(" desc")
		}
	}
}

func deparseSetOpQuery(b *strings.Builder, q *Query) {
	if len(q.CTEList) > 0 {
		deparseCTEsQ(b, q)
	}

	writeQueryQ(b, q.LArg)
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
	writeQueryQ(b, q.RArg)

	if len(q.SortClause) > 0 {
		b.WriteString(" order by ")
		deparseOrderByQ(b, q)
	}

	if q.LimitCount != nil {
		b.WriteString(" limit ")
		writeAnalyzedExpr(b, q.LimitCount, q)
		if q.LimitOffset != nil {
			b.WriteString(" offset ")
			writeAnalyzedExpr(b, q.LimitOffset, q)
		}
	}
}

// backtickIDQ wraps an identifier in backticks.
func backtickIDQ(name string) string {
	return "`" + strings.ReplaceAll(name, "`", "``") + "`"
}
