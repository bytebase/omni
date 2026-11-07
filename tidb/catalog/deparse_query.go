package catalog

import (
	"strings"
)

// DeparseQuery converts an analyzed Query IR back to canonical SQL text.
func DeparseQuery(q *Query) string {
	var b strings.Builder
	writeQueryQ(&b, q, nil)
	return b.String()
}

// deparseScope is the query an expression is written in, together with the
// queries that enclose it, so that a column of an enclosing query resolves.
// The chain follows the analyzer's scopes (analyzeSelectStmtWithCTEs): a
// subquery in an expression is one level inside the query that holds it, and
// a LATERAL derived table is one level inside the query whose FROM holds it.
// A set operation adds no level: its arms see what it sees. A CTE body, and a
// derived table that is not LATERAL, sees the queries that enclose the query
// holding it, not that query.
type deparseScope struct {
	*Query
	outer *deparseScope
}

// rteOf returns the range table entry a column reference points to, LevelsUp
// queries out from s, or nil when the reference does not resolve.
func (s *deparseScope) rteOf(v *VarExprQ) *RangeTableEntryQ {
	for i := 0; s != nil && i < v.LevelsUp; i++ {
		s = s.outer
	}
	if s == nil || s.Query == nil || v.LevelsUp < 0 || v.RangeIdx < 0 || v.RangeIdx >= len(s.RangeTable) {
		return nil
	}
	return s.RangeTable[v.RangeIdx]
}

// writeQueryQ writes a query into the builder of the whole call (see
// writeAnalyzedExpr); a UNION ALL chain nests one Query per arm. outer is the
// scope q is nested in, nil for the top query.
func writeQueryQ(b *strings.Builder, q *Query, outer *deparseScope) {
	if q == nil {
		return
	}
	if q.SetOp != SetOpNone {
		deparseSetOpQuery(b, q, outer)
		return
	}
	deparseSimpleQuery(b, &deparseScope{Query: q, outer: outer})
}

func deparseSimpleQuery(b *strings.Builder, q *deparseScope) {
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

	writeLimitQ(b, q)
}

func deparseCTEsQ(b *strings.Builder, q *deparseScope) {
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
		// A CTE is analyzed in the scope that encloses its WITH clause.
		writeQueryQ(b, cte.Query, q.outer)
		b.WriteString(") ")
	}
}

func deparseTargetListQ(b *strings.Builder, q *deparseScope) {
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
			if rte := q.rteOf(v); rte != nil {
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

func deparseFromClauseQ(b *strings.Builder, q *deparseScope) {
	for i, node := range q.JoinTree.FromList {
		if i > 0 {
			b.WriteString(", ")
		}
		deparseJoinNodeQ(b, node, q)
	}
}

func deparseJoinNodeQ(b *strings.Builder, node JoinNode, q *deparseScope) {
	switch n := node.(type) {
	case *RangeTableRefQ:
		deparseRTERefQ(b, n.RTIndex, q)
	case *JoinExprNodeQ:
		deparseJoinExprNodeQ(b, n, q)
	}
}

func deparseRTERefQ(b *strings.Builder, idx int, q *deparseScope) {
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
		writeRTEClausesQ(b, rte)
	case RTESubquery:
		// Only a LATERAL derived table sees the query whose FROM holds it.
		outer := q.outer
		if rte.Lateral {
			b.WriteString("lateral ")
			outer = q
		}
		b.WriteByte('(')
		writeQueryQ(b, rte.Subquery, outer)
		b.WriteString(") as ")
		b.WriteString(backtickIDQ(rte.ERef))
	case RTECTE:
		b.WriteString(backtickIDQ(rte.CTEName))
		if rte.Alias != "" && rte.Alias != rte.CTEName {
			b.WriteString(" as ")
			b.WriteString(backtickIDQ(rte.Alias))
		}
		writeRTEClausesQ(b, rte)
	case RTEJoin:
		// Synthetic; actual join structure is in JoinExprNodeQ.
	case RTEFunction:
		b.WriteString("/* function-in-FROM */")
	}
}

// writeRTEClausesQ writes the AS OF and TABLESAMPLE clauses of a range table entry,
// in that order, if it has them. Without them the query still reads, and reads the
// current data or the whole table.
func writeRTEClausesQ(b *strings.Builder, rte *RangeTableEntryQ) {
	if rte.AsOf != "" {
		b.WriteByte(' ')
		b.WriteString(rte.AsOf)
	}
	if rte.Sample != "" {
		b.WriteByte(' ')
		b.WriteString(rte.Sample)
	}
}

func deparseJoinExprNodeQ(b *strings.Builder, j *JoinExprNodeQ, q *deparseScope) {
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

func deparseGroupByQ(b *strings.Builder, q *deparseScope) {
	for i, sc := range q.GroupClause {
		if i > 0 {
			b.WriteString(", ")
		}
		if sc.TargetIdx >= 1 && sc.TargetIdx <= len(q.TargetList) {
			writeAnalyzedExpr(b, q.TargetList[sc.TargetIdx-1].Expr, q)
		}
	}
}

func deparseOrderByQ(b *strings.Builder, q *deparseScope) {
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

func deparseSetOpQuery(b *strings.Builder, setOp *Query, outer *deparseScope) {
	q := &deparseScope{Query: setOp, outer: outer}
	if len(q.CTEList) > 0 {
		deparseCTEsQ(b, q)
	}

	// Each arm is analyzed in the scope that encloses the set operation.
	writeQueryQ(b, q.LArg, outer)
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
	writeQueryQ(b, q.RArg, outer)

	if len(q.SortClause) > 0 {
		b.WriteString(" order by ")
		deparseOrderByQ(b, q)
	}

	writeLimitQ(b, q)
}

// writeLimitQ emits the LIMIT/OFFSET clause. Count-first `limit count offset
// offset` is the normal form, but a comma-form `?` marker binds by position, so
// when a marker is present and the source used the comma form the operands are
// emitted as `limit offset,count` to keep the original binding order.
func writeLimitQ(b *strings.Builder, q *deparseScope) {
	if q.LimitCount == nil {
		return
	}
	b.WriteString(" limit ")
	if q.LimitOffset != nil && !q.LimitOffsetKeyword && queryLimitHasParamMarker(q.Query) {
		writeAnalyzedExpr(b, q.LimitOffset, q)
		b.WriteString(",")
		writeAnalyzedExpr(b, q.LimitCount, q)
		return
	}
	writeAnalyzedExpr(b, q.LimitCount, q)
	if q.LimitOffset != nil {
		b.WriteString(" offset ")
		writeAnalyzedExpr(b, q.LimitOffset, q)
	}
}

// queryLimitHasParamMarker reports whether the LIMIT count or offset is a `?`
// parameter marker.
func queryLimitHasParamMarker(q *Query) bool {
	if _, ok := q.LimitCount.(*ParamMarkerExprQ); ok {
		return true
	}
	_, ok := q.LimitOffset.(*ParamMarkerExprQ)
	return ok
}

// backtickIDQ wraps an identifier in backticks.
func backtickIDQ(name string) string {
	return "`" + strings.ReplaceAll(name, "`", "``") + "`"
}
