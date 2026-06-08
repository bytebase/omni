package analysis

import (
	"fmt"
	"sort"
	"strings"

	"github.com/bytebase/omni/redshift"
	"github.com/bytebase/omni/redshift/ast"
	"github.com/bytebase/omni/redshift/catalog"
)

// ColumnResource identifies a source or predicate column.
type ColumnResource struct {
	Schema string
	Table  string
	Column string
}

// QuerySpanResult describes one output column.
type QuerySpanResult struct {
	Name          string
	SourceColumns []ColumnResource
	IsPlainField  bool
}

// QuerySpan contains Redshift SELECT lineage.
type QuerySpan struct {
	Results          []QuerySpanResult
	SourceColumns    []ColumnResource
	PredicateColumns []ColumnResource
}

// GetQuerySpan parses and analyzes a single Redshift SELECT statement.
func GetQuerySpan(c *catalog.Catalog, sql string) (*QuerySpan, error) {
	if strings.Contains(strings.ToLower(sql), "connect by") {
		return nil, fmt.Errorf("Redshift CONNECT BY query span is not supported")
	}
	if c == nil {
		return nil, fmt.Errorf("catalog is nil")
	}
	stmts, err := redshift.Parse(sql)
	if err != nil {
		return nil, err
	}
	if len(stmts) != 1 || stmts[0].Empty() {
		return nil, fmt.Errorf("expected exactly one SELECT statement")
	}
	stmt, ok := stmts[0].AST.(*ast.SelectStmt)
	if !ok {
		return nil, fmt.Errorf("expected SELECT statement, got %T", stmts[0].AST)
	}
	if stmt.IntoClause != nil {
		return nil, fmt.Errorf("SELECT INTO creates a table and is not a read-only query span")
	}
	q, err := c.AnalyzeSelectStmt(stmt)
	if err != nil {
		return nil, err
	}
	return buildQuerySpan(c, q), nil
}

func buildQuerySpan(c *catalog.Catalog, q *catalog.Query) *QuerySpan {
	span := &QuerySpan{
		Results:          collectResults(c, q),
		SourceColumns:    collectAllSourceColumns(c, q),
		PredicateColumns: collectPredicateColumns(c, q),
	}
	return span
}

func collectResults(c *catalog.Catalog, q *catalog.Query) []QuerySpanResult {
	if q == nil {
		return nil
	}
	if q.SetOp != catalog.SetOpNone {
		return collectSetOpResults(c, q)
	}
	var results []QuerySpanResult
	for _, te := range q.TargetList {
		if te == nil || te.ResJunk {
			continue
		}
		_, isVar := te.Expr.(*catalog.VarExpr)
		results = append(results, QuerySpanResult{
			Name:          te.ResName,
			SourceColumns: collectExprColumns(c, q, te.Expr),
			IsPlainField:  isVar,
		})
	}
	return results
}

func collectSetOpResults(c *catalog.Catalog, q *catalog.Query) []QuerySpanResult {
	left := collectResults(c, q.LArg)
	right := collectResults(c, q.RArg)
	resultCount := len(left)
	if len(q.TargetList) > resultCount {
		resultCount = len(q.TargetList)
	}
	results := make([]QuerySpanResult, 0, resultCount)
	for i := 0; i < resultCount; i++ {
		result := QuerySpanResult{}
		if i < len(q.TargetList) && q.TargetList[i] != nil {
			result.Name = q.TargetList[i].ResName
		} else if i < len(left) {
			result.Name = left[i].Name
		}
		if i < len(left) {
			result.SourceColumns = appendUniqueColumns(result.SourceColumns, left[i].SourceColumns...)
			result.IsPlainField = left[i].IsPlainField
		}
		if i < len(right) {
			result.SourceColumns = appendUniqueColumns(result.SourceColumns, right[i].SourceColumns...)
			result.IsPlainField = result.IsPlainField && right[i].IsPlainField
		}
		results = append(results, result)
	}
	return results
}

func collectAllSourceColumns(c *catalog.Catalog, q *catalog.Query) []ColumnResource {
	seen := make(map[ColumnResource]bool)
	var cols []ColumnResource
	for _, result := range collectResults(c, q) {
		appendSeenColumns(&cols, seen, result.SourceColumns...)
	}
	appendSeenColumns(&cols, seen, collectPredicateColumns(c, q)...)
	return sortedColumns(cols)
}

func collectPredicateColumns(c *catalog.Catalog, q *catalog.Query) []ColumnResource {
	seen := make(map[ColumnResource]bool)
	var cols []ColumnResource
	collectPredicateColumnsInto(c, q, seen, &cols)
	return sortedColumns(cols)
}

func collectPredicateColumnsInto(c *catalog.Catalog, q *catalog.Query, seen map[ColumnResource]bool, cols *[]ColumnResource) {
	if q == nil {
		return
	}
	if q.SetOp != catalog.SetOpNone {
		collectPredicateColumnsInto(c, q.LArg, seen, cols)
		collectPredicateColumnsInto(c, q.RArg, seen, cols)
		return
	}
	if q.JoinTree != nil {
		appendSeenColumns(cols, seen, collectExprColumns(c, q, q.JoinTree.Quals)...)
		for _, node := range q.JoinTree.FromList {
			collectJoinPredicateColumns(c, q, node, seen, cols)
		}
	}
	appendSeenColumns(cols, seen, collectExprColumns(c, q, q.HavingQual)...)
	appendSeenColumns(cols, seen, collectExprColumns(c, q, q.QualifyQual)...)
}

func collectJoinPredicateColumns(c *catalog.Catalog, q *catalog.Query, node catalog.JoinNode, seen map[ColumnResource]bool, cols *[]ColumnResource) {
	switch n := node.(type) {
	case *catalog.JoinExprNode:
		appendSeenColumns(cols, seen, collectExprColumns(c, q, n.Quals)...)
		collectJoinPredicateColumns(c, q, n.Left, seen, cols)
		collectJoinPredicateColumns(c, q, n.Right, seen, cols)
	}
}

func collectExprColumns(c *catalog.Catalog, q *catalog.Query, expr catalog.AnalyzedExpr) []ColumnResource {
	seen := make(map[ColumnResource]bool)
	var cols []ColumnResource
	walkExpr(c, q, expr, seen, &cols)
	return sortedColumns(cols)
}

func walkExpr(c *catalog.Catalog, q *catalog.Query, expr catalog.AnalyzedExpr, seen map[ColumnResource]bool, cols *[]ColumnResource) {
	switch v := expr.(type) {
	case nil:
	case *catalog.VarExpr:
		resolveVar(c, q, v, seen, cols)
	case *catalog.FuncCallExpr:
		for _, arg := range v.Args {
			walkExpr(c, q, arg, seen, cols)
		}
	case *catalog.AggExpr:
		for _, arg := range v.Args {
			walkExpr(c, q, arg, seen, cols)
		}
	case *catalog.OpExpr:
		walkExpr(c, q, v.Left, seen, cols)
		walkExpr(c, q, v.Right, seen, cols)
	case *catalog.BoolExprQ:
		for _, arg := range v.Args {
			walkExpr(c, q, arg, seen, cols)
		}
	case *catalog.CaseExprQ:
		walkExpr(c, q, v.Arg, seen, cols)
		for _, w := range v.When {
			walkExpr(c, q, w.Condition, seen, cols)
			walkExpr(c, q, w.Result, seen, cols)
		}
		walkExpr(c, q, v.Default, seen, cols)
	case *catalog.CoalesceExprQ:
		for _, arg := range v.Args {
			walkExpr(c, q, arg, seen, cols)
		}
	case *catalog.SubLinkExpr:
		walkExpr(c, q, v.TestExpr, seen, cols)
		if v.SubQuery != nil {
			for _, te := range v.SubQuery.TargetList {
				if te != nil && !te.ResJunk {
					walkExpr(c, v.SubQuery, te.Expr, seen, cols)
				}
			}
		}
	case *catalog.RelabelExpr:
		walkExpr(c, q, v.Arg, seen, cols)
	case *catalog.CoerceViaIOExpr:
		walkExpr(c, q, v.Arg, seen, cols)
	case *catalog.RowExprQ:
		for _, arg := range v.Args {
			walkExpr(c, q, arg, seen, cols)
		}
	case *catalog.WindowFuncExpr:
		for _, arg := range v.Args {
			walkExpr(c, q, arg, seen, cols)
		}
		walkExpr(c, q, v.AggFilter, seen, cols)
	case *catalog.NullIfExprQ:
		for _, arg := range v.Args {
			walkExpr(c, q, arg, seen, cols)
		}
	case *catalog.MinMaxExprQ:
		for _, arg := range v.Args {
			walkExpr(c, q, arg, seen, cols)
		}
	case *catalog.ArrayExprQ:
		for _, elem := range v.Elements {
			walkExpr(c, q, elem, seen, cols)
		}
	case *catalog.ScalarArrayOpExpr:
		walkExpr(c, q, v.Left, seen, cols)
		walkExpr(c, q, v.Right, seen, cols)
	case *catalog.CollateExprQ:
		walkExpr(c, q, v.Arg, seen, cols)
	}
}

func resolveVar(c *catalog.Catalog, q *catalog.Query, v *catalog.VarExpr, seen map[ColumnResource]bool, cols *[]ColumnResource) {
	if q == nil || v.RangeIdx < 0 || v.RangeIdx >= len(q.RangeTable) {
		return
	}
	rte := q.RangeTable[v.RangeIdx]
	colIdx := int(v.AttNum - 1)
	switch rte.Kind {
	case catalog.RTERelation:
		rel := c.GetRelationByOID(rte.RelOID)
		if rel == nil || rel.Schema == nil {
			return
		}
		if rel.AnalyzedQuery != nil && colIdx >= 0 && colIdx < len(rel.AnalyzedQuery.TargetList) {
			walkExpr(c, rel.AnalyzedQuery, rel.AnalyzedQuery.TargetList[colIdx].Expr, seen, cols)
			return
		}
		colName := ""
		if colIdx >= 0 && colIdx < len(rte.ColNames) {
			colName = rte.ColNames[colIdx]
		}
		appendSeenColumns(cols, seen, ColumnResource{Schema: rel.Schema.Name, Table: rel.Name, Column: colName})
	case catalog.RTESubquery:
		if rte.Subquery != nil && colIdx >= 0 && colIdx < len(rte.Subquery.TargetList) {
			walkExpr(c, rte.Subquery, rte.Subquery.TargetList[colIdx].Expr, seen, cols)
		}
	case catalog.RTECTE:
		if rte.CTEIndex >= 0 && rte.CTEIndex < len(q.CTEList) {
			cte := q.CTEList[rte.CTEIndex]
			if cte.Query != nil && colIdx >= 0 && colIdx < len(cte.Query.TargetList) {
				walkExpr(c, cte.Query, cte.Query.TargetList[colIdx].Expr, seen, cols)
			}
		}
	}
}

func appendUniqueColumns(cols []ColumnResource, items ...ColumnResource) []ColumnResource {
	seen := make(map[ColumnResource]bool)
	for _, col := range cols {
		seen[col] = true
	}
	appendSeenColumns(&cols, seen, items...)
	return cols
}

func appendSeenColumns(cols *[]ColumnResource, seen map[ColumnResource]bool, items ...ColumnResource) {
	for _, col := range items {
		if col.Column == "" || seen[col] {
			continue
		}
		seen[col] = true
		*cols = append(*cols, col)
	}
}

func sortedColumns(cols []ColumnResource) []ColumnResource {
	sort.Slice(cols, func(i, j int) bool {
		if cols[i].Schema != cols[j].Schema {
			return cols[i].Schema < cols[j].Schema
		}
		if cols[i].Table != cols[j].Table {
			return cols[i].Table < cols[j].Table
		}
		return cols[i].Column < cols[j].Column
	})
	return cols
}
