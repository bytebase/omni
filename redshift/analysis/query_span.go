package analysis

import (
	"fmt"
	"sort"
	"strings"

	"github.com/bytebase/omni/redshift"
	"github.com/bytebase/omni/redshift/ast"
	"github.com/bytebase/omni/redshift/catalog"
)

// ColumnResource identifies a source or predicate column. Database is empty
// when the analysis has no database context.
type ColumnResource struct {
	Database string
	Schema   string
	Table    string
	Column   string
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

// parseRedshiftSpan is redshift.Parse behind a variable so a test can prove a
// statement is parsed once per analysis.
var parseRedshiftSpan = redshift.Parse

// GetQuerySpan parses and analyzes a single Redshift SELECT statement.
func GetQuerySpan(c *catalog.Catalog, sql string) (*QuerySpan, error) {
	// Keep the pre-parse rejections in place: a statement the span cannot
	// support fails the same way whether or not it also parses.
	if strings.Contains(strings.ToLower(sql), "connect by") {
		return nil, fmt.Errorf("Redshift CONNECT BY query span is not supported")
	}
	if c == nil {
		return nil, fmt.Errorf("catalog is nil")
	}
	stmts, err := parseRedshiftSpan(sql)
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
	return getQuerySpanFromSelect(c, sql, stmt)
}

// getQuerySpanFromSelect analyzes an already-parsed SELECT. Callers that parsed
// a statement to classify it pass the same AST here, so it is parsed once.
func getQuerySpanFromSelect(c *catalog.Catalog, sql string, stmt *ast.SelectStmt) (*QuerySpan, error) {
	if strings.Contains(strings.ToLower(sql), "connect by") {
		return nil, fmt.Errorf("Redshift CONNECT BY query span is not supported")
	}
	if c == nil {
		return nil, fmt.Errorf("catalog is nil")
	}
	if stmt == nil {
		return nil, fmt.Errorf("expected SELECT statement")
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
		Results:          collectResults(c, nil, q),
		SourceColumns:    collectAllSourceColumns(c, q),
		PredicateColumns: collectPredicateColumns(c, q),
	}
	return span
}

// collectResults returns the output columns of q, nested in the queries of
// parents: the set operations whose branch it is, whose CTEs it may read.
func collectResults(c *catalog.Catalog, parents []*catalog.Query, q *catalog.Query) []QuerySpanResult {
	if q == nil {
		return nil
	}
	if q.SetOp != catalog.SetOpNone {
		return collectSetOpResults(c, parents, q)
	}
	stack := pushQuery(parents, q)
	var results []QuerySpanResult
	for _, te := range q.TargetList {
		if te == nil || te.ResJunk {
			continue
		}
		_, isVar := te.Expr.(*catalog.VarExpr)
		results = append(results, QuerySpanResult{
			Name:          te.ResName,
			SourceColumns: collectExprColumns(c, stack, te.Expr),
			IsPlainField:  isVar,
		})
	}
	return results
}

func collectSetOpResults(c *catalog.Catalog, parents []*catalog.Query, q *catalog.Query) []QuerySpanResult {
	branchParents := pushQuery(parents, q)
	left := collectResults(c, branchParents, q.LArg)
	right := collectResults(c, branchParents, q.RArg)
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
	for _, result := range collectResults(c, nil, q) {
		appendSeenColumns(&cols, seen, result.SourceColumns...)
	}
	appendSeenColumns(&cols, seen, collectPredicateColumns(c, q)...)
	return sortedColumns(cols)
}

func collectPredicateColumns(c *catalog.Catalog, q *catalog.Query) []ColumnResource {
	seen := make(map[ColumnResource]bool)
	var cols []ColumnResource
	collectPredicateColumnsInto(c, []*catalog.Query{q}, seen, &cols)
	return sortedColumns(cols)
}

// collectPredicateColumnsInto collects the predicate columns of the last
// query of stack.
func collectPredicateColumnsInto(c *catalog.Catalog, stack []*catalog.Query, seen map[ColumnResource]bool, cols *[]ColumnResource) {
	q := currentQuery(stack)
	if q == nil {
		return
	}
	if q.SetOp != catalog.SetOpNone {
		collectPredicateColumnsInto(c, pushQuery(stack, q.LArg), seen, cols)
		collectPredicateColumnsInto(c, pushQuery(stack, q.RArg), seen, cols)
		return
	}
	if q.JoinTree != nil {
		appendSeenColumns(cols, seen, collectExprColumns(c, stack, q.JoinTree.Quals)...)
		for _, node := range q.JoinTree.FromList {
			collectJoinPredicateColumns(c, stack, node, seen, cols)
		}
	}
	appendSeenColumns(cols, seen, collectExprColumns(c, stack, q.HavingQual)...)
	appendSeenColumns(cols, seen, collectExprColumns(c, stack, q.QualifyQual)...)
}

func collectJoinPredicateColumns(c *catalog.Catalog, stack []*catalog.Query, node catalog.JoinNode, seen map[ColumnResource]bool, cols *[]ColumnResource) {
	switch n := node.(type) {
	case *catalog.JoinExprNode:
		appendSeenColumns(cols, seen, collectExprColumns(c, stack, n.Quals)...)
		collectJoinPredicateColumns(c, stack, n.Left, seen, cols)
		collectJoinPredicateColumns(c, stack, n.Right, seen, cols)
	}
}

// collectExprColumns returns the columns expr, an expression of the last
// query of stack, reads. The queries before it enclose it, so an outer
// reference or a CTE an enclosing query defines resolves against them.
func collectExprColumns(c *catalog.Catalog, stack []*catalog.Query, expr catalog.AnalyzedExpr) []ColumnResource {
	w := &exprWalker{c: c, seen: make(map[ColumnResource]bool), visited: make(map[queryExpr]bool)}
	w.walkExpr(stack, expr)
	return sortedColumns(w.cols)
}

// exprWalker collects the columns expressions read.
type exprWalker struct {
	c    *catalog.Catalog
	seen map[ColumnResource]bool
	cols []ColumnResource
	// visited stops a walk that comes back to an expression of the same
	// query, as a recursive CTE's reference to itself does.
	visited map[queryExpr]bool
}

type queryExpr struct {
	query *catalog.Query
	expr  catalog.AnalyzedExpr
}

func (w *exprWalker) walkExpr(stack []*catalog.Query, expr catalog.AnalyzedExpr) {
	q := currentQuery(stack)
	key := queryExpr{query: q, expr: expr}
	if expr == nil || w.visited[key] {
		return
	}
	w.visited[key] = true
	catalog.WalkAnalyzedExpr(expr, func(e catalog.AnalyzedExpr) bool {
		switch v := e.(type) {
		case *catalog.VarExpr:
			w.resolveVar(stack, v)
		case *catalog.SetOpColumnExpr:
			// An output column of a set operation reads that column of
			// both branches.
			if q == nil {
				return true
			}
			for _, branch := range []*catalog.Query{q.LArg, q.RArg} {
				if te := outputEntry(branch, int(v.ResNo)); te != nil {
					w.walkExpr(pushQuery(stack, branch), te.Expr)
				}
			}
		case *catalog.SubLinkExpr:
			w.walkExpr(stack, v.TestExpr)
			if v.SubQuery != nil {
				for _, te := range v.SubQuery.TargetList {
					if te != nil && !te.ResJunk {
						w.walkExpr(pushQuery(stack, v.SubQuery), te.Expr)
					}
				}
			}
			return false
		}
		return true
	})
}

// outputEntry returns the target entry of output column resNo (1-based) of
// q, skipping junk entries.
func outputEntry(q *catalog.Query, resNo int) *catalog.TargetEntry {
	if q == nil {
		return nil
	}
	for _, te := range q.TargetList {
		if te == nil || te.ResJunk {
			continue
		}
		if resNo--; resNo == 0 {
			return te
		}
	}
	return nil
}

func (w *exprWalker) resolveVar(stack []*catalog.Query, v *catalog.VarExpr) {
	c := w.c
	// A Var LevelsUp levels out reads the query that many levels up. That
	// stack is a prefix of a stack still in use, so it is copied to extend.
	level := len(stack) - 1 - v.LevelsUp
	if level < 0 || level >= len(stack) {
		return
	}
	extend := pushQuery
	if level < len(stack)-1 {
		extend = appendQuery
	}
	stack = stack[:level+1]
	q := stack[level]
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
			w.walkExpr([]*catalog.Query{rel.AnalyzedQuery}, rel.AnalyzedQuery.TargetList[colIdx].Expr)
			return
		}
		// The table's own column name, which an alias column list may hide.
		colName := ""
		if colIdx >= 0 && colIdx < len(rel.Columns) {
			colName = rel.Columns[colIdx].Name
		}
		appendSeenColumns(&w.cols, w.seen, ColumnResource{Database: rte.CatalogName, Schema: rel.Schema.Name, Table: rel.Name, Column: colName})
	case catalog.RTESubquery:
		if rte.Subquery != nil && colIdx >= 0 && colIdx < len(rte.Subquery.TargetList) {
			w.walkExpr(extend(stack, rte.Subquery), rte.Subquery.TargetList[colIdx].Expr)
		}
	case catalog.RTECTE:
		// rte.CTEIndex indexes the list of the query that defines the CTE,
		// which may enclose q, so the name is matched from q outward.
		for i := level; i >= 0; i-- {
			if stack[i] == nil {
				continue
			}
			for _, cte := range stack[i].CTEList {
				if cte == nil || cte.Name != rte.CTEName {
					continue
				}
				if cte.Query != nil && colIdx >= 0 && colIdx < len(cte.Query.TargetList) {
					w.walkExpr(appendQuery(stack[:i+1], cte.Query), cte.Query.TargetList[colIdx].Expr)
				}
				return
			}
		}
	case catalog.RTEJoin:
		// Only a FULL JOIN USING column references the join itself; it
		// reads both inputs.
		if colIdx >= 0 && colIdx < len(rte.JoinAliasVars) {
			w.walkExpr(stack, rte.JoinAliasVars[colIdx])
		}
	}
}

func currentQuery(stack []*catalog.Query) *catalog.Query {
	if len(stack) == 0 {
		return nil
	}
	return stack[len(stack)-1]
}

// pushQuery returns stack with q appended, reusing stack's array. The caller
// passes its whole stack and is done with the result before it pushes
// another query, so no stack in use shares the slot q takes.
func pushQuery(stack []*catalog.Query, q *catalog.Query) []*catalog.Query {
	return append(stack, q)
}

// appendQuery returns stack with q appended in a new array, for a stack that
// is a prefix of one still in use.
func appendQuery(stack []*catalog.Query, q *catalog.Query) []*catalog.Query {
	return append(stack[:len(stack):len(stack)], q)
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
