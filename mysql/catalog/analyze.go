package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// AnalyzeSelectStmt performs semantic analysis on a parsed SELECT statement,
// returning a resolved Query IR.
func (c *Catalog) AnalyzeSelectStmt(stmt *nodes.SelectStmt) (*Query, error) {
	q := &Query{
		CommandType: CmdSelect,
		JoinTree:    &JoinTreeQ{},
	}

	scope := newScope()

	// Step 1: Analyze FROM clause → populate RangeTable and scope.
	if err := analyzeFromClause(c, stmt.From, q, scope); err != nil {
		return nil, err
	}

	// Step 2: Analyze target list (SELECT expressions).
	if err := analyzeTargetList(stmt.TargetList, q, scope); err != nil {
		return nil, err
	}

	// Step 3: Analyze WHERE clause.
	if stmt.Where != nil {
		analyzed, err := analyzeExpr(stmt.Where, scope)
		if err != nil {
			return nil, err
		}
		q.JoinTree.Quals = analyzed
	}

	// Step 4: GROUP BY
	if len(stmt.GroupBy) > 0 {
		if err := c.analyzeGroupBy(stmt.GroupBy, q, scope); err != nil {
			return nil, err
		}
	}

	// Step 5: HAVING
	if stmt.Having != nil {
		analyzed, err := analyzeExpr(stmt.Having, scope)
		if err != nil {
			return nil, err
		}
		q.HavingQual = analyzed
	}

	// Step 6: Detect aggregates in target list and having
	q.HasAggs = detectAggregates(q)

	// Step 7: ORDER BY
	if len(stmt.OrderBy) > 0 {
		if err := c.analyzeOrderBy(stmt.OrderBy, q, scope); err != nil {
			return nil, err
		}
	}

	// Step 8: LIMIT / OFFSET
	if stmt.Limit != nil {
		if err := c.analyzeLimitOffset(stmt.Limit, q, scope); err != nil {
			return nil, err
		}
	}

	// Step 9: DISTINCT
	q.Distinct = stmt.DistinctKind != nodes.DistinctNone && stmt.DistinctKind != nodes.DistinctAll

	return q, nil
}

// analyzeFromClause processes the FROM clause, populating the query's
// RangeTable, JoinTree.FromList, and the scope for column resolution.
func analyzeFromClause(c *Catalog, from []nodes.TableExpr, q *Query, scope *analyzerScope) error {
	for _, te := range from {
		switch ref := te.(type) {
		case *nodes.TableRef:
			rte, cols, err := analyzeTableRef(c, ref)
			if err != nil {
				return err
			}
			idx := len(q.RangeTable)
			q.RangeTable = append(q.RangeTable, rte)
			q.JoinTree.FromList = append(q.JoinTree.FromList, &RangeTableRefQ{RTIndex: idx})

			// Determine the scope name: alias if present, else table name.
			scopeName := rte.ERef
			scope.add(scopeName, idx, cols)
		default:
			return fmt.Errorf("unsupported FROM clause element: %T", te)
		}
	}
	return nil
}

// analyzeTableRef resolves a table reference from the FROM clause against
// the catalog, returning the RTE and the column list.
func analyzeTableRef(c *Catalog, ref *nodes.TableRef) (*RangeTableEntryQ, []*Column, error) {
	dbName := ref.Schema
	if dbName == "" {
		dbName = c.CurrentDatabase()
	}
	if dbName == "" {
		return nil, nil, errNoDatabaseSelected()
	}

	db := c.GetDatabase(dbName)
	if db == nil {
		return nil, nil, errUnknownDatabase(dbName)
	}

	// Check for a table first, then a view.
	tbl := db.GetTable(ref.Name)
	if tbl != nil {
		eref := ref.Name
		if ref.Alias != "" {
			eref = ref.Alias
		}
		colNames := make([]string, len(tbl.Columns))
		for i, col := range tbl.Columns {
			colNames[i] = col.Name
		}
		rte := &RangeTableEntryQ{
			Kind:      RTERelation,
			DBName:    db.Name,
			TableName: tbl.Name,
			Alias:     ref.Alias,
			ERef:      eref,
			ColNames:  colNames,
		}
		return rte, tbl.Columns, nil
	}

	// Check views.
	view := db.Views[toLower(ref.Name)]
	if view != nil {
		eref := ref.Name
		if ref.Alias != "" {
			eref = ref.Alias
		}
		// Build stub columns from view column names.
		cols := make([]*Column, len(view.Columns))
		colNames := make([]string, len(view.Columns))
		for i, name := range view.Columns {
			cols[i] = &Column{Position: i + 1, Name: name}
			colNames[i] = name
		}
		rte := &RangeTableEntryQ{
			Kind:          RTERelation,
			DBName:        db.Name,
			TableName:     view.Name,
			Alias:         ref.Alias,
			ERef:          eref,
			ColNames:      colNames,
			IsView:        true,
			ViewAlgorithm: viewAlgorithmFromString(view.Algorithm),
		}
		return rte, cols, nil
	}

	return nil, nil, errNoSuchTable(dbName, ref.Name)
}

// viewAlgorithmFromString converts a string algorithm value to ViewAlgorithm.
func viewAlgorithmFromString(s string) ViewAlgorithm {
	switch strings.ToUpper(s) {
	case "MERGE":
		return ViewAlgMerge
	case "TEMPTABLE":
		return ViewAlgTemptable
	case "UNDEFINED", "":
		return ViewAlgUndefined
	default:
		return ViewAlgUndefined
	}
}

// analyzeGroupBy processes the GROUP BY clause, populating q.GroupClause.
func (c *Catalog) analyzeGroupBy(groupBy []nodes.ExprNode, q *Query, scope *analyzerScope) error {
	for _, expr := range groupBy {
		switch n := expr.(type) {
		case *nodes.IntLit:
			// Ordinal reference: GROUP BY 1 means first SELECT column.
			idx := int(n.Value)
			if idx < 1 || idx > len(q.TargetList) {
				return fmt.Errorf("GROUP BY position %d is not in select list", idx)
			}
			q.GroupClause = append(q.GroupClause, &SortGroupClauseQ{
				TargetIdx: idx,
			})
		case *nodes.ColumnRef:
			// Resolve the column ref, then find matching target entry.
			analyzed, err := analyzeExpr(n, scope)
			if err != nil {
				return err
			}
			varExpr, ok := analyzed.(*VarExprQ)
			if !ok {
				return fmt.Errorf("GROUP BY column reference resolved to unexpected type %T", analyzed)
			}
			targetIdx := findMatchingTarget(q.TargetList, varExpr)
			if targetIdx == 0 {
				// Not found in target list — add as junk entry.
				te := &TargetEntryQ{
					Expr:    varExpr,
					ResNo:   len(q.TargetList) + 1,
					ResName: n.Column,
					ResJunk: true,
				}
				q.TargetList = append(q.TargetList, te)
				targetIdx = te.ResNo
			}
			q.GroupClause = append(q.GroupClause, &SortGroupClauseQ{
				TargetIdx: targetIdx,
			})
		default:
			// General expression — analyze and try to match to target list.
			analyzed, err := analyzeExpr(expr, scope)
			if err != nil {
				return err
			}
			targetIdx := 0
			for _, te := range q.TargetList {
				if exprEqual(te.Expr, analyzed) {
					targetIdx = te.ResNo
					break
				}
			}
			if targetIdx == 0 {
				te := &TargetEntryQ{
					Expr:    analyzed,
					ResNo:   len(q.TargetList) + 1,
					ResJunk: true,
				}
				q.TargetList = append(q.TargetList, te)
				targetIdx = te.ResNo
			}
			q.GroupClause = append(q.GroupClause, &SortGroupClauseQ{
				TargetIdx: targetIdx,
			})
		}
	}
	return nil
}

// findMatchingTarget finds a TargetEntryQ whose Expr is a VarExprQ matching
// the given VarExprQ (same RangeIdx and AttNum). Returns ResNo (1-based) or 0.
func findMatchingTarget(tl []*TargetEntryQ, v *VarExprQ) int {
	for _, te := range tl {
		if tv, ok := te.Expr.(*VarExprQ); ok {
			if tv.RangeIdx == v.RangeIdx && tv.AttNum == v.AttNum {
				return te.ResNo
			}
		}
	}
	return 0
}

// exprEqual compares two AnalyzedExpr values for structural equality.
// Phase 1a: only VarExprQ is compared; other types return false.
func exprEqual(a, b AnalyzedExpr) bool {
	va, okA := a.(*VarExprQ)
	vb, okB := b.(*VarExprQ)
	if okA && okB {
		return va.RangeIdx == vb.RangeIdx && va.AttNum == vb.AttNum
	}
	return false
}

// detectAggregates returns true if any aggregate function call exists in the
// query's TargetList or HavingQual.
func detectAggregates(q *Query) bool {
	for _, te := range q.TargetList {
		if exprContainsAggregate(te.Expr) {
			return true
		}
	}
	if q.HavingQual != nil {
		return exprContainsAggregate(q.HavingQual)
	}
	return false
}

// exprContainsAggregate recursively walks an AnalyzedExpr looking for
// FuncCallExprQ with IsAggregate=true.
func exprContainsAggregate(expr AnalyzedExpr) bool {
	if expr == nil {
		return false
	}
	switch e := expr.(type) {
	case *FuncCallExprQ:
		if e.IsAggregate {
			return true
		}
		for _, arg := range e.Args {
			if exprContainsAggregate(arg) {
				return true
			}
		}
	case *OpExprQ:
		return exprContainsAggregate(e.Left) || exprContainsAggregate(e.Right)
	case *BoolExprQ:
		for _, arg := range e.Args {
			if exprContainsAggregate(arg) {
				return true
			}
		}
	case *InListExprQ:
		if exprContainsAggregate(e.Arg) {
			return true
		}
		for _, item := range e.List {
			if exprContainsAggregate(item) {
				return true
			}
		}
	case *BetweenExprQ:
		return exprContainsAggregate(e.Arg) || exprContainsAggregate(e.Lower) || exprContainsAggregate(e.Upper)
	case *NullTestExprQ:
		return exprContainsAggregate(e.Arg)
	case *VarExprQ, *ConstExprQ:
		return false
	}
	return false
}

// analyzeOrderBy processes the ORDER BY clause, populating q.SortClause.
// When an ORDER BY expression is not in the SELECT list, a junk TargetEntryQ
// is added (ResJunk=true).
func (c *Catalog) analyzeOrderBy(orderBy []*nodes.OrderByItem, q *Query, scope *analyzerScope) error {
	for _, item := range orderBy {
		desc := item.Desc
		// MySQL default: ASC → NullsFirst=true, DESC → NullsFirst=false.
		nullsFirst := !desc
		if item.NullsFirst != nil {
			nullsFirst = *item.NullsFirst
		}

		switch n := item.Expr.(type) {
		case *nodes.IntLit:
			// Ordinal reference: ORDER BY 1 means first SELECT column.
			idx := int(n.Value)
			if idx < 1 || idx > len(q.TargetList) {
				return fmt.Errorf("ORDER BY position %d is not in select list", idx)
			}
			q.SortClause = append(q.SortClause, &SortGroupClauseQ{
				TargetIdx:  idx,
				Descending: desc,
				NullsFirst: nullsFirst,
			})
		case *nodes.ColumnRef:
			// Resolve the column ref, then find matching target entry.
			analyzed, err := analyzeExpr(n, scope)
			if err != nil {
				return err
			}
			varExpr, ok := analyzed.(*VarExprQ)
			if !ok {
				return fmt.Errorf("ORDER BY column reference resolved to unexpected type %T", analyzed)
			}
			targetIdx := findMatchingTarget(q.TargetList, varExpr)
			if targetIdx == 0 {
				// Not found in target list — add as junk entry.
				te := &TargetEntryQ{
					Expr:    varExpr,
					ResNo:   len(q.TargetList) + 1,
					ResName: n.Column,
					ResJunk: true,
				}
				q.TargetList = append(q.TargetList, te)
				targetIdx = te.ResNo
			}
			q.SortClause = append(q.SortClause, &SortGroupClauseQ{
				TargetIdx:  targetIdx,
				Descending: desc,
				NullsFirst: nullsFirst,
			})
		default:
			// General expression — analyze and try to match to target list.
			analyzed, err := analyzeExpr(item.Expr, scope)
			if err != nil {
				return err
			}
			targetIdx := 0
			for _, te := range q.TargetList {
				if exprEqual(te.Expr, analyzed) {
					targetIdx = te.ResNo
					break
				}
			}
			if targetIdx == 0 {
				te := &TargetEntryQ{
					Expr:    analyzed,
					ResNo:   len(q.TargetList) + 1,
					ResJunk: true,
				}
				q.TargetList = append(q.TargetList, te)
				targetIdx = te.ResNo
			}
			q.SortClause = append(q.SortClause, &SortGroupClauseQ{
				TargetIdx:  targetIdx,
				Descending: desc,
				NullsFirst: nullsFirst,
			})
		}
	}
	return nil
}

// analyzeLimitOffset processes the LIMIT/OFFSET clause, populating
// q.LimitCount and q.LimitOffset.
func (c *Catalog) analyzeLimitOffset(limit *nodes.Limit, q *Query, scope *analyzerScope) error {
	if limit.Count != nil {
		analyzed, err := analyzeExpr(limit.Count, scope)
		if err != nil {
			return err
		}
		q.LimitCount = analyzed
	}
	if limit.Offset != nil {
		analyzed, err := analyzeExpr(limit.Offset, scope)
		if err != nil {
			return err
		}
		q.LimitOffset = analyzed
	}
	return nil
}
