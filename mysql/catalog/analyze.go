package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// AnalyzeSelectStmt performs semantic analysis on a parsed SELECT statement,
// returning a resolved Query IR.
func (c *Catalog) AnalyzeSelectStmt(stmt *nodes.SelectStmt) (*Query, error) {
	return c.analyzeSelectStmtInternal(stmt, nil)
}

// analyzeSelectStmtInternal is the core analysis routine. parentScope is non-nil
// when analyzing a subquery (for correlated column resolution).
func (c *Catalog) analyzeSelectStmtInternal(stmt *nodes.SelectStmt, parentScope *analyzerScope) (*Query, error) {
	return c.analyzeSelectStmtWithCTEs(stmt, parentScope, nil)
}

// analyzeSelectStmtWithCTEs is the full internal analysis routine.
// inheritedCTEMap provides CTE definitions inherited from an enclosing context
// (e.g., a recursive CTE body referencing itself).
func (c *Catalog) analyzeSelectStmtWithCTEs(stmt *nodes.SelectStmt, parentScope *analyzerScope, inheritedCTEMap map[string]*CommonTableExprQ) (*Query, error) {
	// Handle set operations (UNION/INTERSECT/EXCEPT).
	if stmt.SetOp != nodes.SetOpNone {
		return c.analyzeSetOpWithCTEs(stmt, parentScope, inheritedCTEMap)
	}

	q := &Query{
		CommandType: CmdSelect,
		JoinTree:    &JoinTreeQ{},
	}

	scope := newScopeWithParent(parentScope)

	// Step 0: Process CTEs (WITH clause).
	cteMap, err := c.analyzeCTEs(stmt.CTEs, q, parentScope)
	if err != nil {
		return nil, err
	}

	// Merge inherited CTE map (for recursive CTE self-references).
	if inheritedCTEMap != nil {
		if cteMap == nil {
			cteMap = make(map[string]*CommonTableExprQ)
		}
		for k, v := range inheritedCTEMap {
			if _, exists := cteMap[k]; !exists {
				cteMap[k] = v
			}
		}
	}

	// Step 1: Analyze FROM clause → populate RangeTable and scope.
	if err := analyzeFromClause(c, stmt.From, q, scope, cteMap); err != nil {
		return nil, err
	}

	// Step 2: Analyze target list (SELECT expressions).
	if err := analyzeTargetList(c, stmt.TargetList, q, scope); err != nil {
		return nil, err
	}

	// Step 3: Analyze WHERE clause.
	if stmt.Where != nil {
		analyzed, err := analyzeExpr(c, stmt.Where, scope)
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
		analyzed, err := analyzeExpr(c, stmt.Having, scope)
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
func analyzeFromClause(c *Catalog, from []nodes.TableExpr, q *Query, scope *analyzerScope, cteMap map[string]*CommonTableExprQ) error {
	for _, te := range from {
		joinNode, err := analyzeTableExpr(c, te, q, scope, cteMap)
		if err != nil {
			return err
		}
		q.JoinTree.FromList = append(q.JoinTree.FromList, joinNode)
	}
	return nil
}

// analyzeTableExpr recursively processes a table expression (TableRef,
// JoinClause, or SubqueryExpr used as a derived table), creating RTEs and
// scope entries as appropriate, and returning a JoinNode for the join tree.
func analyzeTableExpr(c *Catalog, te nodes.TableExpr, q *Query, scope *analyzerScope, cteMap map[string]*CommonTableExprQ) (JoinNode, error) {
	switch ref := te.(type) {
	case *nodes.TableRef:
		// Check if this references a CTE before looking up catalog tables.
		if cteMap != nil && ref.Schema == "" {
			lower := strings.ToLower(ref.Name)
			if cteQ, ok := cteMap[lower]; ok {
				return analyzeCTERef(ref, cteQ, q, scope)
			}
		}
		rte, cols, err := analyzeTableRef(c, ref)
		if err != nil {
			return nil, err
		}
		idx := len(q.RangeTable)
		q.RangeTable = append(q.RangeTable, rte)
		scope.add(rte.ERef, idx, cols)
		return &RangeTableRefQ{RTIndex: idx}, nil

	case *nodes.JoinClause:
		return analyzeJoinClause(c, ref, q, scope, cteMap)

	case *nodes.SubqueryExpr:
		return analyzeFromSubquery(c, ref, q, scope)

	default:
		return nil, fmt.Errorf("unsupported FROM clause element: %T", te)
	}
}

// analyzeJoinClause processes a JOIN clause, creating RTEs for the join and
// its children, and returning a JoinExprNodeQ.
func analyzeJoinClause(c *Catalog, jc *nodes.JoinClause, q *Query, scope *analyzerScope, cteMap map[string]*CommonTableExprQ) (JoinNode, error) {
	// Recursively process left and right sides.
	left, err := analyzeTableExpr(c, jc.Left, q, scope, cteMap)
	if err != nil {
		return nil, err
	}
	right, err := analyzeTableExpr(c, jc.Right, q, scope, cteMap)
	if err != nil {
		return nil, err
	}

	// Map AST JoinType to IR JoinTypeQ and detect NATURAL.
	var joinType JoinTypeQ
	natural := false
	switch jc.Type {
	case nodes.JoinInner:
		joinType = JoinInner
	case nodes.JoinLeft:
		joinType = JoinLeft
	case nodes.JoinRight:
		joinType = JoinRight
	case nodes.JoinCross:
		joinType = JoinCross
	case nodes.JoinStraight:
		joinType = JoinStraight
	case nodes.JoinNatural:
		joinType = JoinInner
		natural = true
	case nodes.JoinNaturalLeft:
		joinType = JoinLeft
		natural = true
	case nodes.JoinNaturalRight:
		joinType = JoinRight
		natural = true
	default:
		return nil, fmt.Errorf("unsupported join type: %d", jc.Type)
	}

	// Collect left and right column names for USING/NATURAL coalescing.
	leftCols := collectJoinNodeColNames(left, q)
	rightCols := collectJoinNodeColNames(right, q)

	var usingCols []string
	var quals AnalyzedExpr

	switch cond := jc.Condition.(type) {
	case *nodes.OnCondition:
		quals, err = analyzeExpr(c, cond.Expr, scope)
		if err != nil {
			return nil, err
		}
	case *nodes.UsingCondition:
		usingCols = cond.Columns
	case nil:
		// CROSS JOIN or condition resolved via NATURAL below.
	default:
		return nil, fmt.Errorf("unsupported join condition type: %T", cond)
	}

	// For NATURAL JOIN, compute shared columns.
	if natural {
		usingCols = computeNaturalJoinColumns(leftCols, rightCols)
	}

	// Build coalesced column names for the RTEJoin.
	coalescedCols := buildCoalescedColNames(leftCols, rightCols, usingCols)

	// Mark right-side USING columns as coalesced for star expansion.
	if len(usingCols) > 0 {
		markCoalescedColumns(right, q, scope, usingCols)
	}

	// Create the RTEJoin entry.
	rteJoin := &RangeTableEntryQ{
		Kind:      RTEJoin,
		JoinType:  joinType,
		JoinUsing: usingCols,
		ColNames:  coalescedCols,
	}
	rtIdx := len(q.RangeTable)
	q.RangeTable = append(q.RangeTable, rteJoin)

	joinExpr := &JoinExprNodeQ{
		JoinType:    joinType,
		Left:        left,
		Right:       right,
		Quals:       quals,
		UsingClause: usingCols,
		Natural:     natural,
		RTIndex:     rtIdx,
	}

	return joinExpr, nil
}

// collectJoinNodeColNames returns the column names contributed by a JoinNode.
func collectJoinNodeColNames(node JoinNode, q *Query) []string {
	switch n := node.(type) {
	case *RangeTableRefQ:
		return q.RangeTable[n.RTIndex].ColNames
	case *JoinExprNodeQ:
		return q.RangeTable[n.RTIndex].ColNames
	}
	return nil
}

// computeNaturalJoinColumns finds column names shared between left and right
// (preserving left-side order).
func computeNaturalJoinColumns(leftCols, rightCols []string) []string {
	rightSet := make(map[string]bool, len(rightCols))
	for _, c := range rightCols {
		rightSet[strings.ToLower(c)] = true
	}
	var shared []string
	for _, c := range leftCols {
		if rightSet[strings.ToLower(c)] {
			shared = append(shared, c)
		}
	}
	return shared
}

// buildCoalescedColNames builds the coalesced column list for a JOIN:
// USING columns first (from left), then remaining left columns, then remaining right columns.
func buildCoalescedColNames(leftCols, rightCols, usingCols []string) []string {
	if len(usingCols) == 0 {
		// No coalescing — just concatenate all columns.
		result := make([]string, 0, len(leftCols)+len(rightCols))
		result = append(result, leftCols...)
		result = append(result, rightCols...)
		return result
	}

	usingSet := make(map[string]bool, len(usingCols))
	for _, c := range usingCols {
		usingSet[strings.ToLower(c)] = true
	}

	result := make([]string, 0, len(leftCols)+len(rightCols)-len(usingCols))
	// USING columns first (in USING order, from left side).
	result = append(result, usingCols...)
	// Remaining left columns.
	for _, c := range leftCols {
		if !usingSet[strings.ToLower(c)] {
			result = append(result, c)
		}
	}
	// Remaining right columns (skip USING columns).
	for _, c := range rightCols {
		if !usingSet[strings.ToLower(c)] {
			result = append(result, c)
		}
	}
	return result
}

// markCoalescedColumns marks right-side USING columns as coalesced in scope
// so that star expansion skips them (avoiding duplicate columns).
func markCoalescedColumns(rightNode JoinNode, q *Query, scope *analyzerScope, usingCols []string) {
	usingSet := make(map[string]bool, len(usingCols))
	for _, c := range usingCols {
		usingSet[strings.ToLower(c)] = true
	}

	// Walk the right node's base tables and mark their USING columns.
	markCoalescedInNode(rightNode, q, scope, usingSet)
}

// markCoalescedInNode recursively marks coalesced columns on base tables
// within a join node.
func markCoalescedInNode(node JoinNode, q *Query, scope *analyzerScope, usingSet map[string]bool) {
	switch n := node.(type) {
	case *RangeTableRefQ:
		rte := q.RangeTable[n.RTIndex]
		scopeName := strings.ToLower(rte.ERef)
		for _, colName := range rte.ColNames {
			if usingSet[strings.ToLower(colName)] {
				scope.markCoalesced(scopeName, colName)
			}
		}
	case *JoinExprNodeQ:
		markCoalescedInNode(n.Left, q, scope, usingSet)
		markCoalescedInNode(n.Right, q, scope, usingSet)
	}
}

// analyzeFromSubquery processes a subquery used as a derived table in FROM.
func analyzeFromSubquery(c *Catalog, subq *nodes.SubqueryExpr, q *Query, scope *analyzerScope) (JoinNode, error) {
	// Recursively analyze the inner SELECT. FROM subqueries are not correlated,
	// so we pass nil as parent scope (unless LATERAL).
	var parentScope *analyzerScope
	if subq.Lateral {
		parentScope = scope
	}
	innerQ, err := c.analyzeSelectStmtInternal(subq.Select, parentScope)
	if err != nil {
		return nil, err
	}

	// Derive column names from the inner query's non-junk target list.
	var colNames []string
	for _, te := range innerQ.TargetList {
		if !te.ResJunk {
			colNames = append(colNames, te.ResName)
		}
	}

	// If explicit column aliases are provided, use those instead.
	if len(subq.Columns) > 0 {
		colNames = subq.Columns
	}

	alias := subq.Alias
	if alias == "" {
		alias = fmt.Sprintf("__subquery_%d", len(q.RangeTable))
	}

	rte := &RangeTableEntryQ{
		Kind:     RTESubquery,
		Alias:    alias,
		ERef:     alias,
		ColNames: colNames,
		Subquery: innerQ,
		Lateral:  subq.Lateral,
	}

	idx := len(q.RangeTable)
	q.RangeTable = append(q.RangeTable, rte)

	// Build stub columns for scope resolution.
	cols := make([]*Column, len(colNames))
	for i, name := range colNames {
		cols[i] = &Column{Position: i + 1, Name: name}
	}
	scope.add(alias, idx, cols)

	return &RangeTableRefQ{RTIndex: idx}, nil
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
			analyzed, err := analyzeExpr(c, n, scope)
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
			analyzed, err := analyzeExpr(c, expr, scope)
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
			analyzed, err := analyzeExpr(c, n, scope)
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
			analyzed, err := analyzeExpr(c, item.Expr, scope)
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
		analyzed, err := analyzeExpr(c, limit.Count, scope)
		if err != nil {
			return err
		}
		q.LimitCount = analyzed
	}
	if limit.Offset != nil {
		analyzed, err := analyzeExpr(c, limit.Offset, scope)
		if err != nil {
			return err
		}
		q.LimitOffset = analyzed
	}
	return nil
}

// analyzeCTEs processes WITH clause CTEs, returning a map for CTE name lookup.
func (c *Catalog) analyzeCTEs(ctes []*nodes.CommonTableExpr, q *Query, parentScope *analyzerScope) (map[string]*CommonTableExprQ, error) {
	if len(ctes) == 0 {
		return nil, nil
	}

	cteMap := make(map[string]*CommonTableExprQ)
	for i, cte := range ctes {
		if cte.Recursive {
			q.IsRecursive = true
		}

		var innerQ *Query
		var err error

		if cte.Recursive && cte.Select.SetOp != nodes.SetOpNone {
			// Recursive CTE: analyze the left arm first to establish columns,
			// then register the CTE, then analyze the right arm.
			innerQ, err = c.analyzeRecursiveCTE(cte, parentScope, cteMap)
		} else {
			innerQ, err = c.analyzeSelectStmtInternal(cte.Select, parentScope)
		}
		if err != nil {
			return nil, err
		}

		// Derive column names.
		colNames := cte.Columns
		if len(colNames) == 0 {
			for _, te := range innerQ.TargetList {
				if !te.ResJunk {
					colNames = append(colNames, te.ResName)
				}
			}
		}

		cteQ := &CommonTableExprQ{
			Name:        cte.Name,
			ColumnNames: colNames,
			Query:       innerQ,
			Recursive:   cte.Recursive,
		}
		q.CTEList = append(q.CTEList, cteQ)
		cteMap[strings.ToLower(cte.Name)] = cteQ
		_ = i
	}
	return cteMap, nil
}

// analyzeRecursiveCTE handles WITH RECURSIVE where the CTE body is a set
// operation. The left arm establishes column signatures; the right arm may
// reference the CTE itself.
func (c *Catalog) analyzeRecursiveCTE(cte *nodes.CommonTableExpr, parentScope *analyzerScope, cteMap map[string]*CommonTableExprQ) (*Query, error) {
	stmt := cte.Select

	// Analyze left arm to establish base columns.
	larg, err := c.analyzeSelectStmtInternal(stmt.Left, parentScope)
	if err != nil {
		return nil, err
	}

	// Derive column names from left arm (or explicit column list).
	colNames := cte.Columns
	if len(colNames) == 0 {
		for _, te := range larg.TargetList {
			if !te.ResJunk {
				colNames = append(colNames, te.ResName)
			}
		}
	}

	// Register a temporary CTE entry so the right arm can reference it.
	tempCTE := &CommonTableExprQ{
		Name:        cte.Name,
		ColumnNames: colNames,
		Query:       larg, // temporary — will be replaced
		Recursive:   true,
	}
	cteMap[strings.ToLower(cte.Name)] = tempCTE

	// Analyze right arm with the CTE visible via inherited CTE map.
	rarg, err := c.analyzeSelectStmtWithCTEs(stmt.Right, parentScope, cteMap)
	if err != nil {
		return nil, err
	}

	// Map SetOperation to SetOpType.
	var setOp SetOpType
	switch stmt.SetOp {
	case nodes.SetOpUnion:
		setOp = SetOpUnion
	case nodes.SetOpIntersect:
		setOp = SetOpIntersect
	case nodes.SetOpExcept:
		setOp = SetOpExcept
	}

	// Build result columns from left arm.
	targetList := make([]*TargetEntryQ, 0, len(larg.TargetList))
	for _, te := range larg.TargetList {
		if !te.ResJunk {
			targetList = append(targetList, &TargetEntryQ{
				Expr:    te.Expr,
				ResNo:   te.ResNo,
				ResName: te.ResName,
			})
		}
	}

	q := &Query{
		CommandType: CmdSelect,
		SetOp:       setOp,
		AllSetOp:    stmt.SetAll,
		LArg:        larg,
		RArg:        rarg,
		TargetList:  targetList,
		JoinTree:    &JoinTreeQ{},
	}
	return q, nil
}

// analyzeCTERef creates an RTE and scope entry for a CTE reference in FROM.
func analyzeCTERef(ref *nodes.TableRef, cteQ *CommonTableExprQ, q *Query, scope *analyzerScope) (JoinNode, error) {
	eref := ref.Name
	if ref.Alias != "" {
		eref = ref.Alias
	}

	// Find the CTE's index in the current query's CTEList.
	cteIndex := -1
	for i, c := range q.CTEList {
		if strings.EqualFold(c.Name, cteQ.Name) {
			cteIndex = i
			break
		}
	}
	if cteIndex < 0 {
		cteIndex = 0 // fallback for recursive self-ref during analysis
	}

	colNames := cteQ.ColumnNames

	rte := &RangeTableEntryQ{
		Kind:     RTECTE,
		Alias:    ref.Alias,
		ERef:     eref,
		ColNames: colNames,
		CTEIndex: cteIndex,
		CTEName:  cteQ.Name,
		Subquery: cteQ.Query,
	}

	idx := len(q.RangeTable)
	q.RangeTable = append(q.RangeTable, rte)

	// Build stub columns for scope resolution.
	cols := make([]*Column, len(colNames))
	for i, name := range colNames {
		cols[i] = &Column{Position: i + 1, Name: name}
	}
	scope.add(eref, idx, cols)

	return &RangeTableRefQ{RTIndex: idx}, nil
}

// analyzeSetOp processes a set operation (UNION/INTERSECT/EXCEPT).
func (c *Catalog) analyzeSetOp(stmt *nodes.SelectStmt, parentScope *analyzerScope) (*Query, error) {
	return c.analyzeSetOpWithCTEs(stmt, parentScope, nil)
}

// analyzeSetOpWithCTEs processes a set operation with inherited CTE definitions.
func (c *Catalog) analyzeSetOpWithCTEs(stmt *nodes.SelectStmt, parentScope *analyzerScope, inheritedCTEMap map[string]*CommonTableExprQ) (*Query, error) {
	// Process CTEs if present on the outer set-op node.
	q := &Query{
		CommandType: CmdSelect,
		JoinTree:    &JoinTreeQ{},
	}

	cteMap, err := c.analyzeCTEs(stmt.CTEs, q, parentScope)
	if err != nil {
		return nil, err
	}

	// Merge inherited CTEs.
	if inheritedCTEMap != nil {
		if cteMap == nil {
			cteMap = make(map[string]*CommonTableExprQ)
		}
		for k, v := range inheritedCTEMap {
			if _, exists := cteMap[k]; !exists {
				cteMap[k] = v
			}
		}
	}

	larg, err := c.analyzeSelectStmtWithCTEs(stmt.Left, parentScope, cteMap)
	if err != nil {
		return nil, err
	}
	rarg, err := c.analyzeSelectStmtWithCTEs(stmt.Right, parentScope, cteMap)
	if err != nil {
		return nil, err
	}

	// Map SetOperation to SetOpType.
	var setOp SetOpType
	switch stmt.SetOp {
	case nodes.SetOpUnion:
		setOp = SetOpUnion
	case nodes.SetOpIntersect:
		setOp = SetOpIntersect
	case nodes.SetOpExcept:
		setOp = SetOpExcept
	}

	// Result columns come from the left arm (MySQL convention).
	targetList := make([]*TargetEntryQ, 0, len(larg.TargetList))
	for _, te := range larg.TargetList {
		if !te.ResJunk {
			targetList = append(targetList, &TargetEntryQ{
				Expr:    te.Expr,
				ResNo:   te.ResNo,
				ResName: te.ResName,
			})
		}
	}

	q.SetOp = setOp
	q.AllSetOp = stmt.SetAll
	q.LArg = larg
	q.RArg = rarg
	q.TargetList = targetList

	// Handle ORDER BY / LIMIT on the outer set-op query using a
	// scope built from result columns.
	if len(stmt.OrderBy) > 0 || stmt.Limit != nil {
		setScope := newScope()
		// Build stub columns from target list for ORDER BY resolution.
		cols := make([]*Column, len(targetList))
		colNames := make([]string, len(targetList))
		for i, te := range targetList {
			cols[i] = &Column{Position: i + 1, Name: te.ResName}
			colNames[i] = te.ResName
		}
		// Add a virtual table entry for unqualified column resolution.
		setScope.add("__setop__", 0, cols)

		if len(stmt.OrderBy) > 0 {
			if err := c.analyzeOrderBy(stmt.OrderBy, q, setScope); err != nil {
				return nil, err
			}
		}
		if stmt.Limit != nil {
			if err := c.analyzeLimitOffset(stmt.Limit, q, setScope); err != nil {
				return nil, err
			}
		}
	}

	return q, nil
}

// AnalyzeStandaloneExpr analyzes an expression in the context of a single table.
// Used for CHECK constraints, DEFAULT expressions, and GENERATED column expressions.
func (c *Catalog) AnalyzeStandaloneExpr(expr nodes.ExprNode, table *Table) (AnalyzedExpr, error) {
	scope := newScope()
	// Register the table's columns into the scope at RTE index 0.
	scope.add(table.Name, 0, table.Columns)
	return analyzeExpr(c, expr, scope)
}
