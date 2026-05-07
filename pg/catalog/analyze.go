package catalog

import (
	"fmt"
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
)

// analyzeSelectStmt transforms a raw pgparser SelectStmt into an analyzed Query.
//
// pg: src/backend/parser/analyze.c — transformSelectStmt
func (c *Catalog) analyzeSelectStmt(stmt *nodes.SelectStmt) (*Query, error) {
	if stmt == nil {
		return nil, fmt.Errorf("NULL select statement")
	}

	// Set operation.
	if stmt.Op != nodes.SETOP_NONE {
		q, err := c.analyzeSetOp(stmt)
		if err != nil {
			return nil, err
		}
		// WITH clause on set-op queries attaches at the top level.
		// pg: src/backend/parser/analyze.c — transformSetOperationStmt
		if stmt.WithClause != nil {
			ac := &analyzeCtx{catalog: c, query: q}
			if err := ac.analyzeCTEs(stmt.WithClause); err != nil {
				return nil, err
			}
		}
		return q, nil
	}

	return c.analyzeSimpleSelect(stmt)
}

// analyzeSetOp handles UNION/INTERSECT/EXCEPT.
//
// pg: src/backend/parser/analyze.c — transformSetOperationStmt
func (c *Catalog) analyzeSetOp(stmt *nodes.SelectStmt) (*Query, error) {
	lQuery, err := c.analyzeSelectStmt(stmt.Larg)
	if err != nil {
		return nil, err
	}
	rQuery, err := c.analyzeSelectStmt(stmt.Rarg)
	if err != nil {
		return nil, err
	}

	if len(lQuery.TargetList) != len(rQuery.TargetList) {
		return nil, &Error{
			Code:    CodeDatatypeMismatch,
			Message: fmt.Sprintf("each %s query must have the same number of columns", setOpKeyword(stmt.Op)),
		}
	}

	op := convertSetOpType(stmt.Op)

	// Determine common types for each column position.
	// pg: src/backend/parser/analyze.c — transformSetOperationTree
	// For non-UNKNOWN types: verify coercibility but don't modify the branch expr
	// (the deparser shows branches as-is). For UNKNOWN Const: coerce by changing
	// the Const's type, matching PG's coerce_type behavior.
	tlist := make([]*TargetEntry, len(lQuery.TargetList))
	for i := range lQuery.TargetList {
		lte := lQuery.TargetList[i]
		rte := rQuery.TargetList[i]

		lTypeOID := lte.Expr.exprType()
		rTypeOID := rte.Expr.exprType()

		common, err := c.selectCommonType([]AnalyzedExpr{lte.Expr, rte.Expr}, setOpKeyword(stmt.Op))
		if err != nil {
			return nil, err
		}

		// Left branch coercion.
		// pg: src/backend/parser/analyze.c — transformSetOperationTree (coercion loop)
		if lTypeOID != UNKNOWNOID {
			if _, err := c.coerceToTargetType(lte.Expr, lTypeOID, common, 'i'); err != nil {
				return nil, err
			}
		} else if _, ok := lte.Expr.(*ConstExpr); ok {
			coerced, err := c.coerceToTargetType(lte.Expr, lTypeOID, common, 'i')
			if err != nil {
				return nil, err
			}
			lte.Expr = coerced
		}

		// Right branch coercion.
		if rTypeOID != UNKNOWNOID {
			if _, err := c.coerceToTargetType(rte.Expr, rTypeOID, common, 'i'); err != nil {
				return nil, err
			}
		} else if _, ok := rte.Expr.(*ConstExpr); ok {
			coerced, err := c.coerceToTargetType(rte.Expr, rTypeOID, common, 'i')
			if err != nil {
				return nil, err
			}
			rte.Expr = coerced
		}

		// Align right branch column names with left branch.
		// PG does this via the resultDesc mechanism in the deparser.
		rte.ResName = lte.ResName

		// Top-level target entry with common type for column metadata.
		// Uses a placeholder expression — never shown by the deparser.
		tlist[i] = &TargetEntry{
			Expr:    &VarExpr{TypeOID: common, TypeMod: -1},
			ResNo:   int16(i + 1),
			ResName: lte.ResName,
		}
	}

	return &Query{
		TargetList: tlist,
		SetOp:      op,
		AllSetOp:   stmt.All,
		LArg:       lQuery,
		RArg:       rQuery,
	}, nil
}

// analyzeSimpleSelect handles a simple (non-set-op) SELECT.
//
// pg: src/backend/parser/analyze.c — transformSelectStmt
func (c *Catalog) analyzeSimpleSelect(stmt *nodes.SelectStmt) (*Query, error) {
	q := &Query{}

	// 0. Process WITH clause (CTEs) before anything else.
	// pg: src/backend/parser/analyze.c — transformSelectStmt calls analyzeCTEList
	ac := &analyzeCtx{
		catalog: c,
		query:   q,
	}

	if stmt.WithClause != nil {
		if err := ac.analyzeCTEs(stmt.WithClause); err != nil {
			return nil, err
		}
	}

	if stmt.FromClause != nil {
		for _, item := range stmt.FromClause.Items {
			jn, err := ac.transformFromClauseItem(item)
			if err != nil {
				return nil, err
			}
			if q.JoinTree == nil {
				q.JoinTree = &JoinTree{}
			}
			q.JoinTree.FromList = append(q.JoinTree.FromList, jn)
		}
	}

	// 2. Process SELECT target list.
	if stmt.TargetList != nil {
		for _, item := range stmt.TargetList.Items {
			rt, ok := item.(*nodes.ResTarget)
			if !ok {
				continue
			}
			entries, err := ac.transformTargetEntry(rt, len(q.TargetList))
			if err != nil {
				return nil, err
			}
			q.TargetList = append(q.TargetList, entries...)
		}
	}

	// 2b. Coerce UNKNOWN target entries to TEXT.
	// pg: src/backend/rewrite/rewriteHandler.c — fireRIRrules
	// PG's rewriter coerces UNKNOWN output columns to TEXT for views.
	for _, te := range q.TargetList {
		if !te.ResJunk && te.Expr.exprType() == UNKNOWNOID {
			if c, ok := te.Expr.(*ConstExpr); ok {
				c.TypeOID = TEXTOID
			}
		}
	}

	// 3. Process WHERE clause.
	if stmt.WhereClause != nil {
		qual, err := ac.transformExpr(stmt.WhereClause)
		if err != nil {
			return nil, err
		}
		if q.JoinTree == nil {
			q.JoinTree = &JoinTree{}
		}
		q.JoinTree.Quals = qual
	}

	// 4. Process GROUP BY clause.
	//
	// pg: src/backend/parser/parse_clause.c — transformGroupClause
	if stmt.GroupClause != nil {
		hasGroupingSets := false
		for _, item := range stmt.GroupClause.Items {
			if _, ok := item.(*nodes.GroupingSet); ok {
				hasGroupingSets = true
				break
			}
		}
		if hasGroupingSets {
			for _, item := range stmt.GroupClause.Items {
				gs, err := ac.transformGroupingSet(item, q.TargetList)
				if err != nil {
					return nil, err
				}
				q.GroupingSets = append(q.GroupingSets, gs)
				// Also add flat refs to GroupClause for deparse of simple columns.
				ac.flattenGroupingSetRefs(gs, q)
			}
		} else {
			for _, item := range stmt.GroupClause.Items {
				tle, err := ac.findTargetlistEntry(item, q.TargetList)
				if err != nil {
					return nil, err
				}
				ref := ac.assignSortGroupRef(tle, q.TargetList)
				// Check for duplicate GROUP BY entries.
				isDup := false
				for _, existing := range q.GroupClause {
					if existing.TLESortGroupRef == ref {
						isDup = true
						break
					}
				}
				if !isDup {
					q.GroupClause = append(q.GroupClause, &SortGroupClause{
						TLESortGroupRef: ref,
					})
				}
			}
		}
	}

	// 5. Process HAVING clause.
	if stmt.HavingClause != nil {
		having, err := ac.transformExpr(stmt.HavingClause)
		if err != nil {
			return nil, err
		}
		q.HavingQual = having
	}

	// 5b. Process WINDOW clause (named windows).
	// pg: src/backend/parser/parse_clause.c — transformWindowDefinitions
	if stmt.WindowClause != nil {
		for _, item := range stmt.WindowClause.Items {
			wd, ok := item.(*nodes.WindowDef)
			if !ok {
				continue
			}
			_, err := ac.resolveWindowDef(wd)
			if err != nil {
				return nil, err
			}
		}
	}

	// 6. Process ORDER BY clause (as junk target entries).
	//
	// pg: src/backend/parser/parse_clause.c — transformSortClause
	if stmt.SortClause != nil {
		for _, item := range stmt.SortClause.Items {
			sb, ok := item.(*nodes.SortBy)
			if !ok {
				continue
			}
			tle, err := ac.findTargetlistEntry(sb.Node, q.TargetList)
			if err != nil {
				return nil, err
			}
			ref := ac.assignSortGroupRef(tle, q.TargetList)
			// Check for duplicate ORDER BY entries.
			isDup := false
			for _, existing := range q.SortClause {
				if existing.TLESortGroupRef == ref {
					isDup = true
					break
				}
			}
			if !isDup {
				// Determine sort direction.
				// pg: addTargetToSortList — SORTBY_ASC/DESC/DEFAULT + NULLS
				desc := sb.SortbyDir == nodes.SORTBY_DESC
				nullsFirst := false
				switch sb.SortbyNulls {
				case nodes.SORTBY_NULLS_FIRST:
					nullsFirst = true
				case nodes.SORTBY_NULLS_LAST:
					nullsFirst = false
				default:
					// SORTBY_NULLS_DEFAULT: DESC gets nulls first, ASC gets nulls last
					nullsFirst = desc
				}
				q.SortClause = append(q.SortClause, &SortGroupClause{
					TLESortGroupRef: ref,
					Descending:      desc,
					NullsFirst:      nullsFirst,
				})
			}
		}
	}

	// 7. Process LIMIT/OFFSET.
	if stmt.LimitCount != nil {
		lc, err := ac.transformExpr(stmt.LimitCount)
		if err != nil {
			return nil, err
		}
		q.LimitCount = lc
	}
	if stmt.LimitOffset != nil {
		lo, err := ac.transformExpr(stmt.LimitOffset)
		if err != nil {
			return nil, err
		}
		q.LimitOffset = lo
	}

	// 8. DISTINCT / DISTINCT ON.
	//
	// pg: src/backend/parser/analyze.c — transformDistinctClause / transformDistinctOnClause
	// pgparser: plain DISTINCT → list with one nil item; DISTINCT ON → list of expressions
	if stmt.DistinctClause != nil {
		hasRealExprs := false
		for _, item := range stmt.DistinctClause.Items {
			if item != nil {
				hasRealExprs = true
				break
			}
		}
		if hasRealExprs {
			// DISTINCT ON (col1, col2, ...)
			q.Distinct = true
			for _, item := range stmt.DistinctClause.Items {
				if item == nil {
					continue
				}
				tle, err := ac.findTargetlistEntry(item, q.TargetList)
				if err != nil {
					return nil, err
				}
				ref := ac.assignSortGroupRef(tle, q.TargetList)
				q.DistinctOn = append(q.DistinctOn, &SortGroupClause{
					TLESortGroupRef: ref,
				})
			}
		} else {
			// Plain DISTINCT
			q.Distinct = true
		}
	}

	// 9. FOR UPDATE/SHARE (LockingClause).
	if stmt.LockingClause != nil {
		for _, item := range stmt.LockingClause.Items {
			lc, ok := item.(*nodes.LockingClause)
			if !ok || lc == nil {
				continue
			}
			lcq := &LockingClauseQ{
				Strength:   LockClauseStrength(lc.Strength),
				WaitPolicy: LockWaitPolicy(lc.WaitPolicy),
			}
			if lc.LockedRels != nil {
				for _, rel := range lc.LockedRels.Items {
					if rv, ok := rel.(*nodes.RangeVar); ok && rv != nil {
						lcq.Tables = append(lcq.Tables, rv.Relname)
					}
				}
			}
			q.LockingClauses = append(q.LockingClauses, lcq)
		}
	}

	return q, nil
}

// analyzeCtx holds state during query analysis.
type analyzeCtx struct {
	catalog             *Catalog
	query               *Query
	parent              *analyzeCtx // for correlated subqueries
	domainConstraint    bool        // true when analyzing a domain CHECK constraint
	domainBaseTypeOID   uint32      // base type OID for domain VALUE keyword
	domainBaseTypMod    int32       // base type modifier for domain VALUE keyword
	domainBaseCollation uint32      // base type collation for domain VALUE keyword
}

// transformFromClauseItem processes a FROM clause item.
//
// pg: src/backend/parser/parse_clause.c — transformFromClauseItem
func (ac *analyzeCtx) transformFromClauseItem(n nodes.Node) (JoinNode, error) {
	switch v := n.(type) {
	case *nodes.RangeVar:
		return ac.transformRangeVar(v)
	case *nodes.JoinExpr:
		return ac.transformJoinExpr(v)
	case *nodes.RangeSubselect:
		return ac.transformRangeSubselect(v)
	case *nodes.RangeFunction:
		return ac.transformRangeFunction(v)
	case *nodes.RangeTableSample:
		return ac.transformRangeTableSampleItem(v)
	default:
		return nil, fmt.Errorf("unsupported FROM clause item: %T", n)
	}
}

// transformRangeVar processes a table reference in FROM.
//
// pg: src/backend/parser/parse_clause.c — transformTableEntry
func (ac *analyzeCtx) transformRangeVar(rv *nodes.RangeVar) (JoinNode, error) {
	// Check if this is a CTE reference before looking up tables.
	// pg: src/backend/parser/parse_clause.c — getRTEForSpecialRelationTypes
	if rv.Schemaname == "" {
		for i, cte := range ac.query.CTEList {
			if cte.Name == rv.Relname {
				alias := ""
				if rv.Alias != nil {
					alias = rv.Alias.Aliasname
				}
				return ac.transformCTERef(rv.Relname, alias, i)
			}
		}
		// Also check visibleCTEs for recursive CTE references.
		// During recursive CTE analysis, the partially-defined CTE is stored
		// on the Catalog so that the recursive term (analyzed via a fresh
		// analyzeCtx) can find it.
		// pg: src/backend/parser/analyze.c — determineRecursiveColTypes
		for i, cte := range ac.catalog.visibleCTEs {
			if cte.Name == rv.Relname {
				alias := ""
				if rv.Alias != nil {
					alias = rv.Alias.Aliasname
				}
				return ac.transformVisibleCTERef(cte, rv.Relname, alias, i)
			}
		}
	}

	_, rel, err := ac.catalog.findRelation(rv.Schemaname, rv.Relname)
	if err != nil {
		return nil, err
	}

	alias := ""
	if rv.Alias != nil {
		alias = rv.Alias.Aliasname
	}

	eref := rv.Relname
	if alias != "" {
		eref = alias
	}

	colNames := make([]string, len(rel.Columns))
	colTypes := make([]uint32, len(rel.Columns))
	colTypMods := make([]int32, len(rel.Columns))
	colCollations := make([]uint32, len(rel.Columns))
	for i, col := range rel.Columns {
		colNames[i] = col.Name
		colTypes[i] = col.TypeOID
		colTypMods[i] = col.TypeMod
		colCollations[i] = col.Collation
	}

	rtIdx := len(ac.query.RangeTable)
	rte := &RangeTableEntry{
		Kind:          RTERelation,
		RelOID:        rel.OID,
		RelName:       rv.Relname,
		SchemaName:    rv.Schemaname,
		Alias:         alias,
		ERef:          eref,
		ColNames:      colNames,
		ColTypes:      colTypes,
		ColTypMods:    colTypMods,
		ColCollations: colCollations,
	}
	ac.query.RangeTable = append(ac.query.RangeTable, rte)

	return &RangeTableRef{RTIndex: rtIdx}, nil
}

// transformJoinExpr processes a JOIN expression.
//
// pg: src/backend/parser/parse_clause.c — transformJoinOnClause
func (ac *analyzeCtx) transformJoinExpr(je *nodes.JoinExpr) (JoinNode, error) {
	left, err := ac.transformFromClauseItem(je.Larg)
	if err != nil {
		return nil, err
	}
	right, err := ac.transformFromClauseItem(je.Rarg)
	if err != nil {
		return nil, err
	}

	jt := convertJoinType(je.Jointype, je.IsNatural)

	// Transform ON clause or USING clause.
	// pg: src/backend/parser/parse_clause.c — transformFromClauseItem (USING handling)
	var quals AnalyzedExpr
	var usingClause []string
	if je.UsingClause != nil && len(je.UsingClause.Items) > 0 {
		// USING clause: extract column names and build equality quals.
		usingClause = make([]string, 0, len(je.UsingClause.Items))
		for _, item := range je.UsingClause.Items {
			if s, ok := item.(*nodes.String); ok {
				usingClause = append(usingClause, s.Str)
			}
		}
		quals = ac.buildUsingQuals(left, right, usingClause)
	} else if je.Quals != nil {
		quals, err = ac.transformExpr(je.Quals)
		if err != nil {
			return nil, err
		}
	}

	// Create RTE for the join itself.
	// Per SQL standard and PostgreSQL behavior, USING columns appear once
	// (coalesced from both sides). We collect all left columns, then collect
	// right columns excluding any that are named in the USING clause.
	// pg: src/backend/parser/parse_clause.c — extractRemainingColumns
	var colNames []string
	var colTypes []uint32
	var colTypMods []int32
	var colCollations []uint32
	collectJoinColumns(ac.query.RangeTable, left, &colNames, &colTypes, &colTypMods, &colCollations)
	if len(usingClause) > 0 {
		collectJoinColumnsExcluding(ac.query.RangeTable, right, &colNames, &colTypes, &colTypMods, &colCollations, usingClause)
	} else {
		collectJoinColumns(ac.query.RangeTable, right, &colNames, &colTypes, &colTypMods, &colCollations)
	}

	rtIdx := len(ac.query.RangeTable)
	rte := &RangeTableEntry{
		Kind:          RTEJoin,
		JoinType:      jt,
		ERef:          "",
		ColNames:      colNames,
		ColTypes:      colTypes,
		ColTypMods:    colTypMods,
		ColCollations: colCollations,
	}
	ac.query.RangeTable = append(ac.query.RangeTable, rte)

	return &JoinExprNode{
		JoinType:    jt,
		Left:        left,
		Right:       right,
		Quals:       quals,
		UsingClause: usingClause,
		RTIndex:     rtIdx,
	}, nil
}

// buildUsingQuals builds equality expressions for JOIN USING columns.
//
// pg: src/backend/parser/parse_clause.c — transformFromClauseItem (USING)
func (ac *analyzeCtx) buildUsingQuals(left, right JoinNode, cols []string) AnalyzedExpr {
	var result AnalyzedExpr
	for _, colName := range cols {
		leftVar := ac.findVarInJoinNode(left, colName)
		rightVar := ac.findVarInJoinNode(right, colName)
		if leftVar == nil || rightVar == nil {
			continue
		}
		eq := &OpExpr{
			OpOID:      ac.catalog.findEqualityOp(leftVar.exprType(), rightVar.exprType()),
			ResultType: BOOLOID,
			Left:       leftVar,
			Right:      rightVar,
			OpName:     "=",
		}
		if result == nil {
			result = eq
		} else {
			result = &BoolExprQ{
				Op:   BoolAnd,
				Args: []AnalyzedExpr{result, eq},
			}
		}
	}
	return result
}

// findVarInJoinNode finds a column variable from a join node by name.
// (pgddl helper — PG resolves USING columns via transformFromClauseItem)
func (ac *analyzeCtx) findVarInJoinNode(jn JoinNode, colName string) *VarExpr {
	switch v := jn.(type) {
	case *RangeTableRef:
		rte := ac.query.RangeTable[v.RTIndex]
		for i, name := range rte.ColNames {
			if name == colName {
				var collation uint32
				if i < len(rte.ColCollations) {
					collation = rte.ColCollations[i]
				}
				return &VarExpr{
					RangeIdx:  v.RTIndex,
					AttNum:    int16(i + 1),
					TypeOID:   rte.ColTypes[i],
					TypeMod:   rte.ColTypMods[i],
					Collation: collation,
				}
			}
		}
	case *JoinExprNode:
		rte := ac.query.RangeTable[v.RTIndex]
		for i, name := range rte.ColNames {
			if name == colName {
				var collation uint32
				if i < len(rte.ColCollations) {
					collation = rte.ColCollations[i]
				}
				return &VarExpr{
					RangeIdx:  v.RTIndex,
					AttNum:    int16(i + 1),
					TypeOID:   rte.ColTypes[i],
					TypeMod:   rte.ColTypMods[i],
					Collation: collation,
				}
			}
		}
	}
	return nil
}

// transformRangeSubselect processes a subquery in FROM.
//
// pg: src/backend/parser/parse_clause.c — transformRangeSubselect
func (ac *analyzeCtx) transformRangeSubselect(rs *nodes.RangeSubselect) (JoinNode, error) {
	sub, ok := rs.Subquery.(*nodes.SelectStmt)
	if !ok {
		return nil, fmt.Errorf("subquery in FROM is not a SELECT")
	}

	// LATERAL subqueries can reference columns from preceding FROM items.
	// pg: src/backend/parser/parse_clause.c — transformRangeSubselect
	var subQuery *Query
	var err error
	if rs.Lateral {
		// LATERAL: analyze with parent context for correlation.
		subQuery, err = ac.analyzeSubSelect(sub)
	} else {
		subQuery, err = ac.catalog.analyzeSelectStmt(sub)
	}
	if err != nil {
		return nil, err
	}

	alias := ""
	if rs.Alias != nil {
		alias = rs.Alias.Aliasname
	}

	colNames := make([]string, len(subQuery.TargetList))
	colTypes := make([]uint32, len(subQuery.TargetList))
	colTypMods := make([]int32, len(subQuery.TargetList))
	colCollations := make([]uint32, len(subQuery.TargetList))
	for i, te := range subQuery.TargetList {
		colNames[i] = te.ResName
		colTypes[i] = te.Expr.exprType()
		colTypMods[i] = -1 // typmod lost through subquery
		colCollations[i] = te.Expr.exprCollation()
	}

	rtIdx := len(ac.query.RangeTable)
	rte := &RangeTableEntry{
		Kind:          RTESubquery,
		Alias:         alias,
		ERef:          alias,
		ColNames:      colNames,
		ColTypes:      colTypes,
		ColTypMods:    colTypMods,
		ColCollations: colCollations,
		Subquery:      subQuery,
		Lateral:       rs.Lateral,
	}
	ac.query.RangeTable = append(ac.query.RangeTable, rte)

	return &RangeTableRef{RTIndex: rtIdx}, nil
}

// transformRangeFunction processes a function call in FROM clause.
//
// pg: src/backend/parser/parse_clause.c — transformRangeFunction
func (ac *analyzeCtx) transformRangeFunction(rf *nodes.RangeFunction) (JoinNode, error) {
	if rf.Functions == nil || len(rf.Functions.Items) == 0 {
		return nil, fmt.Errorf("empty function list in RangeFunction")
	}

	// Each item in Functions is a List of {funcexpr, coldeflist}.
	var funcExprs []AnalyzedExpr
	var funcNames []string
	var funcColdefLists []*nodes.List
	for _, item := range rf.Functions.Items {
		funcItem, ok := item.(*nodes.List)
		if !ok || len(funcItem.Items) < 1 {
			return nil, fmt.Errorf("invalid function list item in RangeFunction")
		}

		fexpr := funcItem.Items[0]
		analyzed, err := ac.transformExpr(fexpr)
		if err != nil {
			return nil, err
		}
		funcExprs = append(funcExprs, analyzed)
		funcNames = append(funcNames, figureColname(fexpr))
		var coldeflist *nodes.List
		if len(funcItem.Items) > 1 {
			coldeflist, _ = funcItem.Items[1].(*nodes.List)
		}
		funcColdefLists = append(funcColdefLists, coldeflist)
	}

	return ac.addRangeTableEntryForFunction(rf, funcExprs, funcNames, funcColdefLists)
}

// addRangeTableEntryForFunction builds a RTEFunction RangeTableEntry.
//
// pg: src/backend/parser/parse_relation.c — addRangeTableEntryForFunction
func (ac *analyzeCtx) addRangeTableEntryForFunction(
	rf *nodes.RangeFunction, funcExprs []AnalyzedExpr, funcNames []string, funcColdefLists []*nodes.List,
) (JoinNode, error) {
	// Choose alias name.
	alias := ""
	if rf.Alias != nil {
		alias = rf.Alias.Aliasname
	}
	eref := alias
	if eref == "" && len(funcNames) > 0 {
		eref = funcNames[0]
	}

	// Build column info from each function's return type.
	var colNames []string
	var colTypes []uint32
	var colTypMods []int32
	var colCollations []uint32

	for i, fexpr := range funcExprs {
		retType := fexpr.exprType()
		cls := ac.catalog.getTypeFuncClass(retType)

		switch cls {
		case typeFuncScalar:
			// Single column with name from function.
			cname := chooseScalarFunctionAlias(funcNames[i], alias, len(funcExprs))
			colNames = append(colNames, cname)
			colTypes = append(colTypes, retType)
			colTypMods = append(colTypMods, fexpr.exprTypMod())
			var coll uint32
			if ac.catalog.typeCollation(retType) != 0 {
				coll = fexpr.exprCollation()
				if coll == 0 {
					coll = ac.catalog.typeCollation(retType)
				}
			}
			colCollations = append(colCollations, coll)

		case typeFuncComposite:
			// Extract columns from composite type's relation.
			bt := ac.catalog.typeByOID[retType]
			if bt == nil || bt.RelID == 0 {
				return nil, fmt.Errorf("composite type %d has no relation", retType)
			}
			rel := ac.catalog.findRelByOID(bt.RelID)
			if rel == nil {
				return nil, fmt.Errorf("relation %d for composite type not found", bt.RelID)
			}
			for _, col := range rel.Columns {
				colNames = append(colNames, col.Name)
				colTypes = append(colTypes, col.TypeOID)
				colTypMods = append(colTypMods, col.TypeMod)
				colCollations = append(colCollations, col.Collation)
			}

		case typeFuncRecord:
			coldeflist := recordFunctionColdefList(rf, funcColdefLists, i, len(funcExprs))
			if coldeflist != nil {
				if err := ac.appendColumnsFromColdefList(coldeflist, &colNames, &colTypes, &colTypMods, &colCollations); err != nil {
					return nil, err
				}
				continue
			}
			if ac.appendOutputParamColumns(fexpr, &colNames, &colTypes, &colTypMods, &colCollations) {
				continue
			}
			return nil, &Error{
				Code:    CodeFeatureNotSupported,
				Message: "function returning record requires column definition list",
			}
		default:
			return nil, &Error{
				Code:    CodeFeatureNotSupported,
				Message: fmt.Sprintf("unsupported return type for function in FROM: %s", ac.catalog.typeName(retType)),
			}
		}
	}

	// WITH ORDINALITY adds an extra column.
	if rf.Ordinality {
		colNames = append(colNames, "ordinality")
		colTypes = append(colTypes, INT8OID)
		colTypMods = append(colTypMods, -1)
		colCollations = append(colCollations, 0)
	}

	// Apply explicit column aliases from rf.Alias.Colnames.
	if rf.Alias != nil && rf.Alias.Colnames != nil {
		for i, item := range rf.Alias.Colnames.Items {
			if i >= len(colNames) {
				break
			}
			if s, ok := item.(*nodes.String); ok {
				colNames[i] = s.Str
			}
		}
	}

	rtIdx := len(ac.query.RangeTable)
	rte := &RangeTableEntry{
		Kind:          RTEFunction,
		Alias:         alias,
		ERef:          eref,
		ColNames:      colNames,
		ColTypes:      colTypes,
		ColTypMods:    colTypMods,
		ColCollations: colCollations,
		FuncExprs:     funcExprs,
		Ordinality:    rf.Ordinality,
		Lateral:       rf.Lateral,
	}
	ac.query.RangeTable = append(ac.query.RangeTable, rte)

	return &RangeTableRef{RTIndex: rtIdx}, nil
}

func recordFunctionColdefList(rf *nodes.RangeFunction, funcColdefLists []*nodes.List, idx, nfuncs int) *nodes.List {
	if idx < len(funcColdefLists) && funcColdefLists[idx] != nil {
		return funcColdefLists[idx]
	}
	if nfuncs == 1 && rf.Coldeflist != nil {
		return rf.Coldeflist
	}
	return nil
}

func (ac *analyzeCtx) appendColumnsFromColdefList(
	coldeflist *nodes.List,
	colNames *[]string,
	colTypes *[]uint32,
	colTypMods *[]int32,
	colCollations *[]uint32,
) error {
	for _, item := range coldeflist.Items {
		coldef, ok := item.(*nodes.ColumnDef)
		if !ok {
			return fmt.Errorf("invalid column definition in function column definition list")
		}
		typ := convertTypeNameToInternal(coldef.TypeName)
		typeOID, typmod, err := ac.catalog.ResolveType(typ)
		if err != nil {
			return err
		}
		*colNames = append(*colNames, coldef.Colname)
		*colTypes = append(*colTypes, typeOID)
		*colTypMods = append(*colTypMods, typmod)
		*colCollations = append(*colCollations, ac.catalog.typeCollation(typeOID))
	}
	return nil
}

func (ac *analyzeCtx) appendOutputParamColumns(
	fexpr AnalyzedExpr,
	colNames *[]string,
	colTypes *[]uint32,
	colTypMods *[]int32,
	colCollations *[]uint32,
) bool {
	call, ok := fexpr.(*FuncCallExpr)
	if !ok || call.FuncOID == 0 {
		return false
	}

	var modes []byte
	var names []string
	var allArgTypes []uint32
	if up := ac.catalog.userProcs[call.FuncOID]; up != nil {
		modes = up.ArgModes
		names = up.ArgNames
		allArgTypes = up.AllArgTypes
	} else if proc := ac.catalog.procByOID[call.FuncOID]; proc != nil {
		modes = parseProcArgModes(proc.ArgModes)
		names = parseProcArgNames(proc.ArgNames)
		allArgTypes = proc.AllArgTypes
	}
	if len(modes) == 0 || len(allArgTypes) == 0 {
		return false
	}

	outCount := 0
	for i, mode := range modes {
		if mode != 'o' && mode != 'b' && mode != 't' {
			continue
		}
		if i >= len(allArgTypes) {
			return false
		}
		outCount++
		name := ""
		if i < len(names) {
			name = names[i]
		}
		if name == "" {
			name = fmt.Sprintf("column%d", outCount)
		}
		typeOID := allArgTypes[i]
		*colNames = append(*colNames, name)
		*colTypes = append(*colTypes, typeOID)
		*colTypMods = append(*colTypMods, int32(-1))
		*colCollations = append(*colCollations, ac.catalog.typeCollation(typeOID))
	}
	return outCount > 0
}

func parseProcArgModes(raw string) []byte {
	elems := parseProcArrayElements(raw)
	if len(elems) == 0 {
		return nil
	}
	modes := make([]byte, 0, len(elems))
	for _, elem := range elems {
		if elem == "" {
			continue
		}
		modes = append(modes, elem[0])
	}
	return modes
}

func parseProcArgNames(raw string) []string {
	return parseProcArrayElements(raw)
}

func parseProcArrayElements(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	raw = strings.TrimPrefix(raw, "{")
	raw = strings.TrimSuffix(raw, "}")
	if raw == "" {
		return nil
	}

	var elems []string
	var b strings.Builder
	inQuotes := false
	escaped := false
	for _, r := range raw {
		switch {
		case escaped:
			b.WriteRune(r)
			escaped = false
		case inQuotes && r == '\\':
			escaped = true
		case r == '"':
			inQuotes = !inQuotes
		case r == ',' && !inQuotes:
			elems = append(elems, strings.TrimSpace(b.String()))
			b.Reset()
		default:
			b.WriteRune(r)
		}
	}
	elems = append(elems, strings.TrimSpace(b.String()))
	return elems
}

// chooseScalarFunctionAlias picks the column name for a scalar function in FROM.
//
// pg: src/backend/parser/parse_relation.c — chooseScalarFunctionAlias
func chooseScalarFunctionAlias(funcName, rteAlias string, nfuncs int) string {
	if nfuncs == 1 && rteAlias != "" {
		return rteAlias
	}
	return funcName
}

// figureColname extracts the function name from a raw expression for FROM alias.
//
// pg: src/backend/parser/parse_target.c — FigureColname
func figureColname(n nodes.Node) string {
	if n == nil {
		return "?column?"
	}
	switch v := n.(type) {
	case *nodes.FuncCall:
		if v.Funcname != nil && len(v.Funcname.Items) > 0 {
			return stringVal(v.Funcname.Items[len(v.Funcname.Items)-1])
		}
	}
	return "?column?"
}

// transformTargetEntry processes a single SELECT target list entry.
// A star expression expands into multiple entries.
//
// pg: src/backend/parser/parse_target.c — transformTargetEntry
func (ac *analyzeCtx) transformTargetEntry(rt *nodes.ResTarget, startIdx int) ([]*TargetEntry, error) {
	// Check for star expression.
	if cr, ok := rt.Val.(*nodes.ColumnRef); ok {
		if isStar(cr) {
			return ac.expandStar(cr, startIdx)
		}
	}

	expr, err := ac.transformExpr(rt.Val)
	if err != nil {
		return nil, err
	}

	name := rt.Name
	if name == "" {
		name = figureColName(rt.Val)
	}

	te := &TargetEntry{
		Expr:    expr,
		ResNo:   int16(startIdx + 1),
		ResName: name,
	}

	// Track provenance for simple column refs.
	if v, ok := expr.(*VarExpr); ok {
		rte := ac.query.RangeTable[v.RangeIdx]
		te.ResOrigTbl = rte.RelOID
		te.ResOrigCol = v.AttNum
	}

	return []*TargetEntry{te}, nil
}

// expandStar expands a * or table.* reference.
//
// pg: src/backend/parser/parse_target.c — ExpandAllTables / expandNSItemAttrs
func (ac *analyzeCtx) expandStar(cr *nodes.ColumnRef, startIdx int) ([]*TargetEntry, error) {
	tableName := starTableName(cr)

	var entries []*TargetEntry

	if tableName == "" && ac.query.JoinTree != nil {
		// Unqualified *: walk the join tree so that JOIN USING deduplication
		// is respected. Each top-level FROM item contributes its columns;
		// for JoinExprNode the RTE already has the correct (deduplicated) list.
		for _, jn := range ac.query.JoinTree.FromList {
			ac.expandStarFromJoinNode(jn, &entries, startIdx)
		}
	} else {
		// Qualified table.* or no join tree: iterate RTEs directly.
		for rtIdx, rte := range ac.query.RangeTable {
			if rte.Kind == RTEJoin {
				continue
			}
			if tableName != "" && rte.ERef != tableName {
				continue
			}
			for colIdx, colName := range rte.ColNames {
				var coll uint32
				if colIdx < len(rte.ColCollations) {
					coll = rte.ColCollations[colIdx]
				}
				te := &TargetEntry{
					Expr: &VarExpr{
						RangeIdx:  rtIdx,
						AttNum:    int16(colIdx + 1),
						TypeOID:   rte.ColTypes[colIdx],
						TypeMod:   rte.ColTypMods[colIdx],
						Collation: coll,
					},
					ResNo:      int16(startIdx + len(entries) + 1),
					ResName:    colName,
					ResOrigTbl: rte.RelOID,
					ResOrigCol: int16(colIdx + 1),
				}
				entries = append(entries, te)
			}
			if tableName != "" {
				break
			}
		}
	}

	if len(entries) == 0 {
		if tableName != "" {
			return nil, errUndefinedTable(tableName)
		}
		return nil, &Error{Code: CodeFeatureNotSupported, Message: "SELECT * with no tables specified is not valid"}
	}

	return entries, nil
}

// expandStarFromJoinNode expands columns from a JoinNode for unqualified *.
// For JoinExprNode with USING, the left side is expanded fully and the right
// side skips USING columns, matching PostgreSQL's behavior. Vars still point
// to the base table RTEs so that lineage tracking is preserved.
func (ac *analyzeCtx) expandStarFromJoinNode(jn JoinNode, entries *[]*TargetEntry, startIdx int) {
	switch v := jn.(type) {
	case *RangeTableRef:
		rte := ac.query.RangeTable[v.RTIndex]
		for colIdx, colName := range rte.ColNames {
			var coll uint32
			if colIdx < len(rte.ColCollations) {
				coll = rte.ColCollations[colIdx]
			}
			te := &TargetEntry{
				Expr: &VarExpr{
					RangeIdx:  v.RTIndex,
					AttNum:    int16(colIdx + 1),
					TypeOID:   rte.ColTypes[colIdx],
					TypeMod:   rte.ColTypMods[colIdx],
					Collation: coll,
				},
				ResNo:      int16(startIdx + len(*entries) + 1),
				ResName:    colName,
				ResOrigTbl: rte.RelOID,
				ResOrigCol: int16(colIdx + 1),
			}
			*entries = append(*entries, te)
		}
	case *JoinExprNode:
		// Expand left side fully.
		ac.expandStarFromJoinNode(v.Left, entries, startIdx)
		if len(v.UsingClause) > 0 {
			// JOIN USING: expand right side but skip USING columns.
			ac.expandStarFromJoinNodeExcluding(v.Right, entries, startIdx, v.UsingClause)
		} else {
			// JOIN ON or CROSS JOIN: expand right side normally.
			ac.expandStarFromJoinNode(v.Right, entries, startIdx)
		}
	}
}

// expandStarFromJoinNodeExcluding expands columns from a JoinNode, skipping
// the first occurrence of each column name in the exclude list. This implements
// the PostgreSQL behavior where USING columns from the right side are omitted.
func (ac *analyzeCtx) expandStarFromJoinNodeExcluding(jn JoinNode, entries *[]*TargetEntry, startIdx int, exclude []string) {
	// Collect all entries from this node first, then filter.
	var rightEntries []*TargetEntry
	ac.expandStarFromJoinNode(jn, &rightEntries, startIdx)

	excludeSet := make(map[string]bool, len(exclude))
	for _, name := range exclude {
		excludeSet[name] = true
	}
	for _, te := range rightEntries {
		if excludeSet[te.ResName] {
			delete(excludeSet, te.ResName)
			continue
		}
		te.ResNo = int16(startIdx + len(*entries) + 1)
		*entries = append(*entries, te)
	}
}

// transformExpr transforms a raw expression node into an analyzed expression.
//
// pg: src/backend/parser/parse_expr.c — transformExprRecurse
func (ac *analyzeCtx) transformExpr(n nodes.Node) (AnalyzedExpr, error) {
	if n == nil {
		return &ConstExpr{TypeOID: UNKNOWNOID, TypeMod: -1, IsNull: true}, nil
	}

	switch v := n.(type) {
	case *nodes.ColumnRef:
		return ac.transformColumnRef(v)
	case *nodes.A_Const:
		return ac.transformAConst(v)
	case *nodes.FuncCall:
		return ac.transformFuncCall(v)
	case *nodes.A_Expr:
		return ac.transformAExpr(v)
	case *nodes.TypeCast:
		return ac.transformTypeCast(v)
	case *nodes.BoolExpr:
		return ac.transformBoolExpr(v)
	case *nodes.CaseExpr:
		return ac.transformCaseExpr(v)
	case *nodes.CoalesceExpr:
		return ac.transformCoalesceExpr(v)
	case *nodes.SubLink:
		return ac.transformSubLink(v)
	case *nodes.NullTest:
		return ac.transformNullTest(v)
	case *nodes.MinMaxExpr:
		return ac.transformMinMaxExpr(v)
	case *nodes.BooleanTest:
		return ac.transformBooleanTest(v)
	case *nodes.SQLValueFunction:
		return ac.transformSQLValueFunction(v)
	case *nodes.ArrayExpr:
		return ac.transformArrayExpr(v)
	case *nodes.A_ArrayExpr:
		return ac.transformAArrayExpr(v)
	case *nodes.RowExpr:
		return ac.transformRowExpr(v)
	case *nodes.CollateClause:
		return ac.transformCollateClause(v)
	case *nodes.Integer:
		return ac.transformAConst(&nodes.A_Const{Val: v})
	case *nodes.Float:
		return ac.transformAConst(&nodes.A_Const{Val: v})
	case *nodes.String:
		return ac.transformAConst(&nodes.A_Const{Val: v})
	case *nodes.Boolean:
		return ac.transformAConst(&nodes.A_Const{Val: v})
	case *nodes.A_Indirection:
		return ac.transformIndirection(v)
	case *nodes.NamedArgExpr:
		return ac.transformNamedArgExpr(v)
	case *nodes.ParamRef:
		return ac.transformParamRef(v)
	case *nodes.GroupingFunc:
		return ac.transformGroupingFunc(v)
	case *nodes.XmlExpr:
		return ac.transformXmlExpr(v)
	case *nodes.XmlSerialize:
		return ac.transformXmlSerialize(v)
	case *nodes.JsonIsPredicate:
		return ac.transformJsonIsPredicate(v)
	case *nodes.JsonFuncExpr:
		return ac.transformJsonFuncExpr(v)
	case *nodes.JsonObjectConstructor:
		return ac.transformJsonObjectConstructor(v)
	case *nodes.JsonArrayConstructor:
		return ac.transformJsonArrayConstructor(v)
	case *nodes.JsonValueExpr:
		return ac.transformJsonValueExpr(v)
	case *nodes.JsonScalarExpr:
		return ac.transformJsonScalarExpr(v)
	case *nodes.JsonParseExpr:
		return ac.transformJsonParseExpr(v)
	case *nodes.JsonSerializeExpr:
		return ac.transformJsonSerializeExpr(v)
	default:
		return &ConstExpr{TypeOID: UNKNOWNOID, TypeMod: -1, IsNull: true}, nil
	}
}

// transformColumnRef resolves a column reference to a Var.
//
// pg: src/backend/parser/parse_expr.c — transformColumnRef
func (ac *analyzeCtx) transformColumnRef(cr *nodes.ColumnRef) (AnalyzedExpr, error) {
	if cr.Fields == nil {
		return nil, errUndefinedColumn("")
	}

	items := cr.Fields.Items
	var tableName, colName string

	switch len(items) {
	case 1:
		colName = stringVal(items[0])
	case 2:
		tableName = stringVal(items[0])
		colName = stringVal(items[1])
	default:
		// schema.table.column — take last two parts.
		colName = stringVal(items[len(items)-1])
		tableName = stringVal(items[len(items)-2])
	}

	// pg: src/backend/commands/typecmds.c — replace_domain_constraint_value
	// In domain CHECK constraints, "value" is replaced with CoerceToDomainValue.
	if ac.domainConstraint && tableName == "" && strings.EqualFold(colName, "value") {
		return &CoerceToDomainValueExpr{
			TypeOID:   ac.domainBaseTypeOID,
			TypeMod:   ac.domainBaseTypMod,
			Collation: ac.domainBaseCollation,
		}, nil
	}

	return ac.resolveColumnRef(tableName, colName)
}

// resolveColumnRef resolves a column name to a VarExpr.
// If the column is not found in the current query's range table, and a parent
// context exists, it searches the parent for correlated subquery references.
func (ac *analyzeCtx) resolveColumnRef(tableName, colName string) (*VarExpr, error) {
	var found *VarExpr
	for rtIdx, rte := range ac.query.RangeTable {
		if tableName != "" && rte.ERef != tableName {
			continue
		}
		for colIdx, cn := range rte.ColNames {
			if cn == colName {
				var coll uint32
				if colIdx < len(rte.ColCollations) {
					coll = rte.ColCollations[colIdx]
				}
				v := &VarExpr{
					RangeIdx:  rtIdx,
					AttNum:    int16(colIdx + 1),
					TypeOID:   rte.ColTypes[colIdx],
					TypeMod:   rte.ColTypMods[colIdx],
					Collation: coll,
				}
				if found != nil {
					return nil, errAmbiguousColumn(colName)
				}
				found = v
				if tableName != "" {
					return found, nil // exact table match, no ambiguity possible
				}
			}
		}
	}

	if found == nil && ac.parent != nil {
		// Try resolving in parent context (correlated subquery).
		// pg: src/backend/parser/parse_expr.c — transformColumnRef
		// PG walks up the namespace chain to find outer references.
		outerVar, err := ac.parent.resolveColumnRef(tableName, colName)
		if err != nil {
			return nil, err
		}
		// Increment LevelsUp to indicate this is an outer reference.
		return &VarExpr{
			RangeIdx:  outerVar.RangeIdx,
			AttNum:    outerVar.AttNum,
			TypeOID:   outerVar.TypeOID,
			TypeMod:   outerVar.TypeMod,
			Collation: outerVar.Collation,
			LevelsUp:  outerVar.LevelsUp + 1,
		}, nil
	}

	if found == nil {
		if tableName != "" {
			// Check if table exists at all in the range table.
			tableFound := false
			for _, rte := range ac.query.RangeTable {
				if rte.ERef == tableName {
					tableFound = true
					break
				}
			}
			if !tableFound {
				return nil, errUndefinedTable(tableName)
			}
		}
		return nil, errUndefinedColumn(colName)
	}
	return found, nil
}

// transformAConst transforms a constant value.
//
// pg: src/backend/parser/parse_expr.c — make_const
func (ac *analyzeCtx) transformAConst(ac2 *nodes.A_Const) (AnalyzedExpr, error) {
	if ac2.Isnull {
		return &ConstExpr{TypeOID: UNKNOWNOID, TypeMod: -1, IsNull: true}, nil
	}

	switch v := ac2.Val.(type) {
	case *nodes.Integer:
		// pg: make_const — PG uses int4 for values that fit, int8 for larger.
		// pgparser uses Go int64 for Integer, so check if it overflows int32.
		if v.Ival < -2147483648 || v.Ival > 2147483647 {
			return &ConstExpr{TypeOID: INT8OID, TypeMod: -1, Value: fmt.Sprintf("%d", v.Ival)}, nil
		}
		return &ConstExpr{TypeOID: INT4OID, TypeMod: -1, Value: fmt.Sprintf("%d", v.Ival)}, nil
	case *nodes.Float:
		// PG parser puts large integers and real floats here.
		// pg: src/backend/parser/parse_node.c — make_const
		// Check if it looks like an integer (no decimal point or exponent).
		if !strings.Contains(v.Fval, ".") && !strings.ContainsAny(v.Fval, "eE") {
			// Large integer that didn't fit int4. Try int8 first.
			// PG: scanint8 check — if it fits int8, use INT8OID.
			_, err := strconv.ParseInt(v.Fval, 10, 64)
			if err == nil {
				return &ConstExpr{TypeOID: INT8OID, TypeMod: -1, Value: v.Fval}, nil
			}
			// Doesn't fit int8 either → NUMERIC
			return &ConstExpr{TypeOID: NUMERICOID, TypeMod: -1, Value: v.Fval}, nil
		}
		return &ConstExpr{TypeOID: NUMERICOID, TypeMod: -1, Value: v.Fval}, nil
	case *nodes.String:
		return &ConstExpr{TypeOID: UNKNOWNOID, TypeMod: -1, Collation: DEFAULT_COLLATION_OID, Value: v.Str}, nil
	case *nodes.Boolean:
		val := "false"
		if v.Boolval {
			val = "true"
		}
		return &ConstExpr{TypeOID: BOOLOID, TypeMod: -1, Value: val}, nil
	default:
		return &ConstExpr{TypeOID: UNKNOWNOID, TypeMod: -1}, nil
	}
}

// transformFuncCall transforms a function call.
//
// pg: src/backend/parser/parse_func.c — ParseFuncOrColumn
func (ac *analyzeCtx) transformFuncCall(fc *nodes.FuncCall) (AnalyzedExpr, error) {
	_, funcName := qualifiedName(fc.Funcname)
	funcName = strings.ToLower(funcName)

	// Transform arguments.
	var args []AnalyzedExpr
	var argTypes []uint32
	if fc.Args != nil {
		for _, a := range fc.Args.Items {
			expr, err := ac.transformExpr(a)
			if err != nil {
				return nil, err
			}
			args = append(args, expr)
			argTypes = append(argTypes, expr.exprType())
		}
	}

	// count(*) special case.
	if fc.AggStar && funcName == "count" {
		// If OVER is present, it's a window function.
		if fc.Over != nil {
			winRef, err := ac.resolveWindowDef(fc.Over)
			if err != nil {
				return nil, err
			}
			var aggFilter AnalyzedExpr
			if fc.AggFilter != nil {
				aggFilter, err = ac.transformExpr(fc.AggFilter)
				if err != nil {
					return nil, err
				}
			}
			return &WindowFuncExpr{
				FuncOID:    0,
				FuncName:   "count",
				ResultType: INT8OID,
				Args:       nil,
				AggStar:    true,
				AggFilter:  aggFilter,
				WinRef:     winRef,
			}, nil
		}
		return &AggExpr{
			AggFuncOID: 0,
			AggName:    "count",
			ResultType: INT8OID,
			AggStar:    true,
		}, nil
	}

	// Resolve function.
	procs := ac.catalog.procByName[funcName]
	if len(procs) == 0 {
		return nil, errUndefinedFunction(funcName, argTypes)
	}

	proc, matchedArgTypes, err := ac.catalog.resolveFuncOverload(procs, argTypes, fc.AggStar)
	if err != nil {
		return nil, err
	}

	// Insert coercions for mismatched argument types.
	// Skip coercion when the target type is polymorphic — polymorphic types
	// accept any actual type without coercion.
	for i := range args {
		if i < int(proc.NArgs) && argTypes[i] != matchedArgTypes[i] && !isPolymorphic(matchedArgTypes[i]) {
			args[i], err = ac.catalog.coerceToTargetType(args[i], argTypes[i], matchedArgTypes[i], 'i')
			if err != nil {
				return nil, err
			}
		}
	}

	retType := ac.catalog.resolveReturnType(proc, argTypes)

	// Window function: has OVER clause.
	// pg: src/backend/parser/parse_func.c — ParseFuncOrColumn (window function path)
	if fc.Over != nil {
		winRef, err := ac.resolveWindowDef(fc.Over)
		if err != nil {
			return nil, err
		}
		var aggFilter AnalyzedExpr
		if fc.AggFilter != nil {
			aggFilter, err = ac.transformExpr(fc.AggFilter)
			if err != nil {
				return nil, err
			}
		}
		// Derive collation for window function result.
		var winColl uint32
		if ac.catalog.typeCollation(retType) != 0 {
			winColl = resolveCollation(args...)
			if winColl == 0 {
				winColl = ac.catalog.typeCollation(retType)
			}
		}
		return &WindowFuncExpr{
			FuncOID:    proc.OID,
			FuncName:   funcName,
			ResultType: retType,
			Collation:  winColl,
			Args:       args,
			AggStar:    fc.AggStar,
			AggFilter:  aggFilter,
			WinRef:     winRef,
		}, nil
	}

	// Determine if this is an aggregate.
	if proc.Kind == 'a' {
		var aggColl uint32
		if ac.catalog.typeCollation(retType) != 0 {
			aggColl = resolveCollation(args...)
			if aggColl == 0 {
				aggColl = ac.catalog.typeCollation(retType)
			}
		}
		return &AggExpr{
			AggFuncOID:  proc.OID,
			AggName:     funcName,
			ResultType:  retType,
			Collation:   aggColl,
			Args:        args,
			AggStar:     fc.AggStar,
			AggDistinct: fc.AggDistinct,
		}, nil
	}

	// Determine result collation: if result type is collatable, derive from args.
	var funcColl uint32
	if ac.catalog.typeCollation(retType) != 0 {
		funcColl = resolveCollation(args...)
		if funcColl == 0 {
			funcColl = ac.catalog.typeCollation(retType)
		}
	}

	return &FuncCallExpr{
		FuncOID:      proc.OID,
		FuncName:     funcName,
		ResultType:   retType,
		ResultTypMod: -1,
		Collation:    funcColl,
		Args:         args,
	}, nil
}

// transformAExpr transforms an operator expression.
//
// pg: src/backend/parser/parse_expr.c — transformAExprOp
func (ac *analyzeCtx) transformAExpr(ae *nodes.A_Expr) (AnalyzedExpr, error) {
	// Handle special A_Expr kinds before normal operator resolution.
	// pg: src/backend/parser/parse_expr.c — transformExprRecurse (A_Expr cases)
	switch ae.Kind {
	case nodes.AEXPR_NULLIF:
		return ac.transformNullIf(ae)
	case nodes.AEXPR_IN:
		return ac.transformAExprIn(ae)
	case nodes.AEXPR_BETWEEN, nodes.AEXPR_NOT_BETWEEN:
		return ac.transformBetween(ae)
	case nodes.AEXPR_BETWEEN_SYM, nodes.AEXPR_NOT_BETWEEN_SYM:
		return ac.transformBetween(ae)
	case nodes.AEXPR_DISTINCT, nodes.AEXPR_NOT_DISTINCT:
		return ac.transformDistinct(ae)
	case nodes.AEXPR_OP, nodes.AEXPR_LIKE, nodes.AEXPR_ILIKE, nodes.AEXPR_SIMILAR:
		// These all use normal operator resolution — fall through.
	case nodes.AEXPR_OP_ANY:
		return ac.transformAExprOpAny(ae)
	case nodes.AEXPR_OP_ALL:
		return ac.transformAExprOpAll(ae)
	}

	opName := ""
	if ae.Name != nil && len(ae.Name.Items) > 0 {
		opName = stringVal(ae.Name.Items[len(ae.Name.Items)-1])
	}

	var left AnalyzedExpr
	var leftOID uint32
	if ae.Lexpr != nil {
		var err error
		left, err = ac.transformExpr(ae.Lexpr)
		if err != nil {
			return nil, err
		}
		leftOID = left.exprType()
	}

	right, err := ac.transformExpr(ae.Rexpr)
	if err != nil {
		return nil, err
	}
	rightOID := right.exprType()

	// Row comparison: (a,b) op (c,d) → RowCompareExpr.
	// pg: src/backend/parser/parse_expr.c — make_row_comparison_op
	if leftRow, lok := left.(*RowExprQ); lok {
		if rightRow, rok := right.(*RowExprQ); rok {
			rcType := opNameToRowCompareType(opName)
			if rcType != 0 {
				// Resolve operators per-column.
				var opNos []uint32
				for i := range leftRow.Args {
					if i >= len(rightRow.Args) {
						break
					}
					lType := leftRow.Args[i].exprType()
					rType := rightRow.Args[i].exprType()
					colOp, _, _, err := ac.catalog.resolveOperator(opName, lType, rType, false)
					if err != nil {
						return nil, err
					}
					opNos = append(opNos, colOp.OID)
				}
				return &RowCompareExprQ{
					RCType: rcType,
					LArgs:  leftRow.Args,
					RArgs:  rightRow.Args,
					OpNos:  opNos,
				}, nil
			}
		}
	}

	isPrefix := ae.Lexpr == nil

	// Resolve operator.
	op, resolvedLeft, resolvedRight, err := ac.catalog.resolveOperator(opName, leftOID, rightOID, isPrefix)
	if err != nil {
		return nil, err
	}

	// Insert coercions if needed.
	if left != nil && leftOID != resolvedLeft {
		left, err = ac.catalog.coerceToTargetType(left, leftOID, resolvedLeft, 'i')
		if err != nil {
			return nil, err
		}
	}
	if rightOID != resolvedRight {
		right, err = ac.catalog.coerceToTargetType(right, rightOID, resolvedRight, 'i')
		if err != nil {
			return nil, err
		}
	}

	// Determine result collation for operator.
	var opColl uint32
	if ac.catalog.typeCollation(op.Result) != 0 {
		opColl = resolveCollation(left, right)
		if opColl == 0 {
			opColl = ac.catalog.typeCollation(op.Result)
		}
	}

	return &OpExpr{
		OpOID:      op.OID,
		OpName:     op.Name,
		ResultType: op.Result,
		Collation:  opColl,
		Left:       left,
		Right:      right,
	}, nil
}

// transformTypeCast transforms an explicit type cast.
//
// pg: src/backend/parser/parse_expr.c — transformTypeCast
func (ac *analyzeCtx) transformTypeCast(tc *nodes.TypeCast) (AnalyzedExpr, error) {
	arg, err := ac.transformExpr(tc.Arg)
	if err != nil {
		return nil, err
	}

	targetOID, targetTypMod, err := ac.catalog.resolveTypeName(tc.TypeName)
	if err != nil {
		return nil, err
	}

	// If the argument is a NULL constant, directly change its type rather than
	// wrapping in a coercion node. PG does this in coerce_type() for Const nodes.
	if c, ok := arg.(*ConstExpr); ok && c.IsNull {
		return &ConstExpr{
			TypeOID: targetOID,
			TypeMod: targetTypMod,
			IsNull:  true,
		}, nil
	}

	srcType := arg.exprType()
	if srcType == targetOID {
		// Same type — may still need typmod coercion.
		if targetTypMod != -1 {
			return &RelabelExpr{
				Arg:        arg,
				ResultType: targetOID,
				TypeMod:    targetTypMod,
				Format:     'e',
			}, nil
		}
		return arg, nil
	}

	result, err := ac.catalog.coerceToTargetType(arg, srcType, targetOID, 'e')
	if err != nil {
		return nil, err
	}
	// Apply typmod if present (e.g., varchar(50) → typmod=54).
	// pg: coerce_type_typmod applies length coercion after type coercion.
	if targetTypMod != -1 {
		switch r := result.(type) {
		case *RelabelExpr:
			r.TypeMod = targetTypMod
		case *FuncCallExpr:
			r.ResultTypMod = targetTypMod
		case *CoerceViaIOExpr:
			// CoerceViaIO doesn't carry typmod directly; wrap in RelabelExpr.
			result = &RelabelExpr{
				Arg:        result,
				ResultType: targetOID,
				TypeMod:    targetTypMod,
				Format:     'i', // implicit
			}
		}
	}
	return result, nil
}

// transformBoolExpr transforms AND/OR/NOT expressions.
//
// pg: src/backend/parser/parse_expr.c — transformBoolExpr
func (ac *analyzeCtx) transformBoolExpr(be *nodes.BoolExpr) (AnalyzedExpr, error) {
	op := BoolAnd
	switch be.Boolop {
	case nodes.OR_EXPR:
		op = BoolOr
	case nodes.NOT_EXPR:
		op = BoolNot
	}

	var args []AnalyzedExpr
	if be.Args != nil {
		for _, a := range be.Args.Items {
			expr, err := ac.transformExpr(a)
			if err != nil {
				return nil, err
			}
			args = append(args, expr)
		}
	}

	return &BoolExprQ{Op: op, Args: args}, nil
}

// transformCaseExpr transforms a CASE expression.
//
// pg: src/backend/parser/parse_expr.c — transformCaseExpr
func (ac *analyzeCtx) transformCaseExpr(ce *nodes.CaseExpr) (AnalyzedExpr, error) {
	// Transform the CASE argument for simple CASE (CASE arg WHEN val THEN ...).
	// pg: src/backend/parser/parse_expr.c — transformCaseExpr
	var caseArg AnalyzedExpr
	if ce.Arg != nil {
		var err error
		caseArg, err = ac.transformExpr(ce.Arg)
		if err != nil {
			return nil, err
		}
	}

	var whens []*CaseWhenQ
	var resultExprs []AnalyzedExpr

	if ce.Args != nil {
		for _, a := range ce.Args.Items {
			cw, ok := a.(*nodes.CaseWhen)
			if !ok {
				continue
			}
			cond, err := ac.transformExpr(cw.Expr)
			if err != nil {
				return nil, err
			}
			result, err := ac.transformExpr(cw.Result)
			if err != nil {
				return nil, err
			}
			whens = append(whens, &CaseWhenQ{Condition: cond, Result: result})
			resultExprs = append(resultExprs, result)
		}
	}

	var defResult AnalyzedExpr
	if ce.Defresult != nil {
		var err error
		defResult, err = ac.transformExpr(ce.Defresult)
		if err != nil {
			return nil, err
		}
		resultExprs = append(resultExprs, defResult)
	} else {
		resultExprs = append(resultExprs, &ConstExpr{TypeOID: UNKNOWNOID})
	}

	common, err := ac.catalog.selectCommonType(resultExprs, "CASE")
	if err != nil {
		return nil, err
	}

	// Coerce results to common type.
	// pg: src/backend/parser/parse_expr.c — transformCaseExpr
	// PG coerces ALL branch results including UNKNOWN (string literals).
	for i, w := range whens {
		rt := w.Result.exprType()
		if rt != common {
			whens[i].Result, err = ac.catalog.coerceToTargetType(w.Result, rt, common, 'i')
			if err != nil {
				return nil, err
			}
		}
	}
	if defResult != nil {
		dt := defResult.exprType()
		if dt != common {
			defResult, err = ac.catalog.coerceToTargetType(defResult, dt, common, 'i')
			if err != nil {
				return nil, err
			}
		}
	} else {
		// pg: transformCaseExpr — PG always adds an implicit ELSE NULL::resultType
		// when there is no ELSE clause.
		defResult = &ConstExpr{TypeOID: common, IsNull: true}
	}

	// Derive collation from all result expressions.
	var caseColl uint32
	if ac.catalog.typeCollation(common) != 0 {
		var allResults []AnalyzedExpr
		for _, w := range whens {
			allResults = append(allResults, w.Result)
		}
		if defResult != nil {
			allResults = append(allResults, defResult)
		}
		caseColl = resolveCollation(allResults...)
		if caseColl == 0 {
			caseColl = ac.catalog.typeCollation(common)
		}
	}

	return &CaseExprQ{
		Arg:        caseArg,
		When:       whens,
		Default:    defResult,
		ResultType: common,
		Collation:  caseColl,
	}, nil
}

// transformCoalesceExpr transforms a COALESCE expression.
//
// pg: src/backend/parser/parse_expr.c — transformCoalesceExpr
func (ac *analyzeCtx) transformCoalesceExpr(ce *nodes.CoalesceExpr) (AnalyzedExpr, error) {
	var args []AnalyzedExpr

	if ce.Args != nil {
		for _, a := range ce.Args.Items {
			expr, err := ac.transformExpr(a)
			if err != nil {
				return nil, err
			}
			args = append(args, expr)
		}
	}

	common, err := ac.catalog.selectCommonType(args, "COALESCE")
	if err != nil {
		return nil, err
	}

	// Coerce args to common type.
	// pg: src/backend/parser/parse_expr.c — transformCoalesceExpr
	// PG coerces ALL args including UNKNOWN (string literals).
	for i, arg := range args {
		at := arg.exprType()
		if at != common {
			args[i], err = ac.catalog.coerceToTargetType(arg, at, common, 'i')
			if err != nil {
				return nil, err
			}
		}
	}

	// Derive collation from args.
	var coalColl uint32
	if ac.catalog.typeCollation(common) != 0 {
		coalColl = resolveCollation(args...)
		if coalColl == 0 {
			coalColl = ac.catalog.typeCollation(common)
		}
	}

	return &CoalesceExprQ{Args: args, ResultType: common, Collation: coalColl}, nil
}

// transformSubLink transforms a subquery expression (scalar or EXISTS).
//
// pg: src/backend/parser/parse_expr.c — transformSubLink
func (ac *analyzeCtx) transformSubLink(sl *nodes.SubLink) (AnalyzedExpr, error) {
	sub, ok := sl.Subselect.(*nodes.SelectStmt)
	if !ok {
		return nil, fmt.Errorf("subquery is not a SELECT")
	}

	// Analyze subquery with parent context for correlated subqueries.
	subQuery, err := ac.analyzeSubSelect(sub)
	if err != nil {
		return nil, err
	}

	switch nodes.SubLinkType(sl.SubLinkType) {
	case nodes.EXISTS_SUBLINK:
		return &SubLinkExpr{
			SubLinkType: SubLinkExistsType,
			SubQuery:    subQuery,
			ResultType:  BOOLOID,
		}, nil
	case nodes.ANY_SUBLINK:
		if len(subQuery.TargetList) != 1 {
			return nil, &Error{Code: CodeTooManyColumns, Message: "subquery must return only one column"}
		}
		var testExpr AnalyzedExpr
		if sl.Testexpr != nil {
			testExpr, err = ac.transformExpr(sl.Testexpr)
			if err != nil {
				return nil, err
			}
		}
		return &SubLinkExpr{
			SubLinkType: SubLinkAnyType,
			TestExpr:    testExpr,
			SubQuery:    subQuery,
			ResultType:  BOOLOID,
		}, nil
	case nodes.ALL_SUBLINK:
		if len(subQuery.TargetList) != 1 {
			return nil, &Error{Code: CodeTooManyColumns, Message: "subquery must return only one column"}
		}
		var testExpr AnalyzedExpr
		if sl.Testexpr != nil {
			testExpr, err = ac.transformExpr(sl.Testexpr)
			if err != nil {
				return nil, err
			}
		}
		return &SubLinkExpr{
			SubLinkType: SubLinkAllType,
			TestExpr:    testExpr,
			SubQuery:    subQuery,
			ResultType:  BOOLOID,
		}, nil
	default:
		// EXPR_SUBLINK (scalar subquery)
		if len(subQuery.TargetList) != 1 {
			return nil, &Error{Code: CodeTooManyColumns, Message: "subquery must return only one column"}
		}
		return &SubLinkExpr{
			SubLinkType: SubLinkExprType,
			SubQuery:    subQuery,
			ResultType:  subQuery.TargetList[0].Expr.exprType(),
		}, nil
	}
}

// transformNullTest transforms IS [NOT] NULL.
//
// pg: src/backend/parser/parse_expr.c — transformNullTest (simplified)
func (ac *analyzeCtx) transformNullTest(nt *nodes.NullTest) (AnalyzedExpr, error) {
	arg, err := ac.transformExpr(nt.Arg)
	if err != nil {
		return nil, err
	}

	return &NullTestExpr{
		Arg:    arg,
		IsNull: nt.Nulltesttype == nodes.IS_NULL,
	}, nil
}

// analyzeSubSelect analyzes a subquery SELECT with parent context for correlation.
//
// pg: src/backend/parser/analyze.c — transformSelectStmt (with parent pstate)
func (ac *analyzeCtx) analyzeSubSelect(stmt *nodes.SelectStmt) (*Query, error) {
	if stmt == nil {
		return nil, fmt.Errorf("NULL select statement")
	}

	// Set operation — analyze recursively without parent context on branches.
	if stmt.Op != nodes.SETOP_NONE {
		return ac.catalog.analyzeSetOp(stmt)
	}

	// Simple SELECT with parent context.
	q := &Query{}
	subAc := &analyzeCtx{
		catalog: ac.catalog,
		query:   q,
		parent:  ac,
	}

	// Process WITH clause if present.
	if stmt.WithClause != nil {
		if err := subAc.analyzeCTEs(stmt.WithClause); err != nil {
			return nil, err
		}
	}

	if stmt.FromClause != nil {
		for _, item := range stmt.FromClause.Items {
			jn, err := subAc.transformFromClauseItem(item)
			if err != nil {
				return nil, err
			}
			if q.JoinTree == nil {
				q.JoinTree = &JoinTree{}
			}
			q.JoinTree.FromList = append(q.JoinTree.FromList, jn)
		}
	}

	if stmt.TargetList != nil {
		for _, item := range stmt.TargetList.Items {
			rt, ok := item.(*nodes.ResTarget)
			if !ok {
				continue
			}
			entries, err := subAc.transformTargetEntry(rt, len(q.TargetList))
			if err != nil {
				return nil, err
			}
			q.TargetList = append(q.TargetList, entries...)
		}
	}

	if stmt.WhereClause != nil {
		qual, err := subAc.transformExpr(stmt.WhereClause)
		if err != nil {
			return nil, err
		}
		if q.JoinTree == nil {
			q.JoinTree = &JoinTree{}
		}
		q.JoinTree.Quals = qual
	}

	// GROUP BY.
	if stmt.GroupClause != nil {
		hasGroupingSets := false
		for _, item := range stmt.GroupClause.Items {
			if _, ok := item.(*nodes.GroupingSet); ok {
				hasGroupingSets = true
				break
			}
		}
		if hasGroupingSets {
			for _, item := range stmt.GroupClause.Items {
				gs, err := subAc.transformGroupingSet(item, q.TargetList)
				if err != nil {
					return nil, err
				}
				q.GroupingSets = append(q.GroupingSets, gs)
				subAc.flattenGroupingSetRefs(gs, q)
			}
		} else {
			for _, item := range stmt.GroupClause.Items {
				tle, err := subAc.findTargetlistEntry(item, q.TargetList)
				if err != nil {
					return nil, err
				}
				ref := subAc.assignSortGroupRef(tle, q.TargetList)
				isDup := false
				for _, existing := range q.GroupClause {
					if existing.TLESortGroupRef == ref {
						isDup = true
						break
					}
				}
				if !isDup {
					q.GroupClause = append(q.GroupClause, &SortGroupClause{TLESortGroupRef: ref})
				}
			}
		}
	}

	if stmt.HavingClause != nil {
		having, err := subAc.transformExpr(stmt.HavingClause)
		if err != nil {
			return nil, err
		}
		q.HavingQual = having
	}

	if stmt.SortClause != nil {
		for _, item := range stmt.SortClause.Items {
			sb, ok := item.(*nodes.SortBy)
			if !ok {
				continue
			}
			tle, err := subAc.findTargetlistEntry(sb.Node, q.TargetList)
			if err != nil {
				return nil, err
			}
			ref := subAc.assignSortGroupRef(tle, q.TargetList)
			isDup := false
			for _, existing := range q.SortClause {
				if existing.TLESortGroupRef == ref {
					isDup = true
					break
				}
			}
			if !isDup {
				desc := sb.SortbyDir == nodes.SORTBY_DESC
				nullsFirst := false
				switch sb.SortbyNulls {
				case nodes.SORTBY_NULLS_FIRST:
					nullsFirst = true
				case nodes.SORTBY_NULLS_LAST:
					nullsFirst = false
				default:
					nullsFirst = desc
				}
				q.SortClause = append(q.SortClause, &SortGroupClause{
					TLESortGroupRef: ref,
					Descending:      desc,
					NullsFirst:      nullsFirst,
				})
			}
		}
	}

	if stmt.LimitCount != nil {
		lc, err := subAc.transformExpr(stmt.LimitCount)
		if err != nil {
			return nil, err
		}
		q.LimitCount = lc
	}
	if stmt.LimitOffset != nil {
		lo, err := subAc.transformExpr(stmt.LimitOffset)
		if err != nil {
			return nil, err
		}
		q.LimitOffset = lo
	}

	if stmt.DistinctClause != nil {
		q.Distinct = true
	}

	return q, nil
}

// transformMinMaxExpr transforms a GREATEST/LEAST expression.
//
// pg: src/backend/parser/parse_expr.c — transformMinMaxExpr
func (ac *analyzeCtx) transformMinMaxExpr(mme *nodes.MinMaxExpr) (AnalyzedExpr, error) {
	var args []AnalyzedExpr

	if mme.Args != nil {
		for _, a := range mme.Args.Items {
			expr, err := ac.transformExpr(a)
			if err != nil {
				return nil, err
			}
			args = append(args, expr)
		}
	}

	mmContext := "GREATEST"
	if mme.Op == nodes.IS_LEAST {
		mmContext = "LEAST"
	}
	common, err := ac.catalog.selectCommonType(args, mmContext)
	if err != nil {
		return nil, err
	}

	// Coerce args to common type.
	for i, arg := range args {
		at := arg.exprType()
		if at != common {
			args[i], err = ac.catalog.coerceToTargetType(arg, at, common, 'i')
			if err != nil {
				return nil, err
			}
		}
	}

	op := MinMaxGreatest
	if mme.Op == nodes.IS_LEAST {
		op = MinMaxLeast
	}

	// Derive collation from args.
	var mmColl uint32
	if ac.catalog.typeCollation(common) != 0 {
		mmColl = resolveCollation(args...)
		if mmColl == 0 {
			mmColl = ac.catalog.typeCollation(common)
		}
	}

	return &MinMaxExprQ{
		Op:         op,
		Args:       args,
		ResultType: common,
		Collation:  mmColl,
	}, nil
}

// transformBooleanTest transforms IS [NOT] TRUE/FALSE/UNKNOWN.
//
// pg: src/backend/parser/parse_expr.c — transformBooleanTest
func (ac *analyzeCtx) transformBooleanTest(bt *nodes.BooleanTest) (AnalyzedExpr, error) {
	arg, err := ac.transformExpr(bt.Arg)
	if err != nil {
		return nil, err
	}

	var testType BoolTestType
	switch bt.Booltesttype {
	case nodes.IS_TRUE:
		testType = BoolIsTrue
	case nodes.IS_NOT_TRUE:
		testType = BoolIsNotTrue
	case nodes.IS_FALSE:
		testType = BoolIsFalse
	case nodes.IS_NOT_FALSE:
		testType = BoolIsNotFalse
	case nodes.IS_UNKNOWN:
		testType = BoolIsUnknown
	case nodes.IS_NOT_UNKNOWN:
		testType = BoolIsNotUnknown
	}

	return &BooleanTestExpr{
		Arg:      arg,
		TestType: testType,
	}, nil
}

// transformSQLValueFunction transforms CURRENT_DATE, CURRENT_TIMESTAMP, etc.
//
// pg: src/backend/parser/parse_expr.c — transformSQLValueFunction
func (ac *analyzeCtx) transformSQLValueFunction(svf *nodes.SQLValueFunction) (AnalyzedExpr, error) {
	var op SVFOp
	var typeOID uint32
	typMod := int32(-1)

	switch svf.Op {
	case nodes.SVFOP_CURRENT_DATE:
		op = SVFCurrentDate
		typeOID = DATEOID
	case nodes.SVFOP_CURRENT_TIME:
		op = SVFCurrentTime
		typeOID = TIMETZOID
	case nodes.SVFOP_CURRENT_TIME_N:
		op = SVFCurrentTimeN
		typeOID = TIMETZOID
		typMod = int32(svf.Typmod)
	case nodes.SVFOP_CURRENT_TIMESTAMP:
		op = SVFCurrentTimestamp
		typeOID = TIMESTAMPTZOID
	case nodes.SVFOP_CURRENT_TIMESTAMP_N:
		op = SVFCurrentTimestampN
		typeOID = TIMESTAMPTZOID
		typMod = int32(svf.Typmod)
	case nodes.SVFOP_LOCALTIME:
		op = SVFLocaltime
		typeOID = TIMEOID
	case nodes.SVFOP_LOCALTIME_N:
		op = SVFLocaltimeN
		typeOID = TIMEOID
		typMod = int32(svf.Typmod)
	case nodes.SVFOP_LOCALTIMESTAMP:
		op = SVFLocaltimestamp
		typeOID = TIMESTAMPOID
	case nodes.SVFOP_LOCALTIMESTAMP_N:
		op = SVFLocaltimestampN
		typeOID = TIMESTAMPOID
		typMod = int32(svf.Typmod)
	case nodes.SVFOP_CURRENT_ROLE:
		op = SVFCurrentRole
		typeOID = NAMEOID
	case nodes.SVFOP_CURRENT_USER:
		op = SVFCurrentUser
		typeOID = NAMEOID
	case nodes.SVFOP_USER:
		op = SVFUser
		typeOID = NAMEOID
	case nodes.SVFOP_SESSION_USER:
		op = SVFSessionUser
		typeOID = NAMEOID
	case nodes.SVFOP_CURRENT_CATALOG:
		op = SVFCurrentCatalog
		typeOID = NAMEOID
	case nodes.SVFOP_CURRENT_SCHEMA:
		op = SVFCurrentSchema
		typeOID = NAMEOID
	default:
		return nil, fmt.Errorf("unrecognized SQLValueFunction op: %d", svf.Op)
	}

	return &SQLValueFuncExpr{
		Op:      op,
		TypeOID: typeOID,
		TypeMod: typMod,
	}, nil
}

// transformNullIf transforms a NULLIF(a, b) expression.
//
// pg: src/backend/parser/parse_expr.c — transformAExprNullIf
func (ac *analyzeCtx) transformNullIf(ae *nodes.A_Expr) (AnalyzedExpr, error) {
	left, err := ac.transformExpr(ae.Lexpr)
	if err != nil {
		return nil, err
	}
	right, err := ac.transformExpr(ae.Rexpr)
	if err != nil {
		return nil, err
	}

	leftOID := left.exprType()
	rightOID := right.exprType()

	// Resolve the equality operator for the two args.
	op, resolvedLeft, resolvedRight, err := ac.catalog.resolveOperator("=", leftOID, rightOID, false)
	if err != nil {
		return nil, err
	}

	if leftOID != resolvedLeft {
		left, err = ac.catalog.coerceToTargetType(left, leftOID, resolvedLeft, 'i')
		if err != nil {
			return nil, err
		}
	}
	if rightOID != resolvedRight {
		right, err = ac.catalog.coerceToTargetType(right, rightOID, resolvedRight, 'i')
		if err != nil {
			return nil, err
		}
	}

	// NULLIF result type is the type of the first argument.
	// pg: src/backend/parser/parse_expr.c — transformAExprNullIf
	return &NullIfExprQ{
		OpOID:      op.OID,
		Args:       []AnalyzedExpr{left, right},
		ResultType: left.exprType(),
	}, nil
}

// makeScalarArrayOp builds a ScalarArrayOpExpr from a scalar and an array expression.
// This is the shared core for IN-list, = ANY(array), and op ALL(array).
//
// pg: src/backend/parser/parse_oper.c — make_scalar_array_op
func (ac *analyzeCtx) makeScalarArrayOp(opName string, useOr bool,
	left, right AnalyzedExpr) (AnalyzedExpr, error) {
	leftOID := left.exprType()
	rightOID := right.exprType()

	// Extract element type from array.
	// pg: make_scalar_array_op — get_base_element_type(rtypeId)
	elemOID := ac.catalog.getBaseElementType(rightOID)
	if elemOID == 0 {
		if rightOID == UNKNOWNOID {
			elemOID = UNKNOWNOID
		} else {
			return nil, &Error{
				Code:    CodeDatatypeMismatch,
				Message: fmt.Sprintf("op %s requires array on right side, got %s", opName, ac.catalog.typeName(rightOID)),
			}
		}
	}

	// Resolve operator on element types (not array types).
	op, resolvedLeft, resolvedRight, err := ac.catalog.resolveOperator(opName, leftOID, elemOID, false)
	if err != nil {
		return nil, err
	}

	// Validate return type is boolean.
	if op.Result != BOOLOID {
		return nil, &Error{
			Code:    CodeDatatypeMismatch,
			Message: fmt.Sprintf("op %s must return boolean for ANY/ALL, got %s", opName, ac.catalog.typeName(op.Result)),
		}
	}

	// Coerce left argument if needed.
	if leftOID != resolvedLeft {
		left, err = ac.catalog.coerceToTargetType(left, leftOID, resolvedLeft, 'i')
		if err != nil {
			return nil, err
		}
	}

	// Map resolved right type back to its array type for coercion.
	if resolvedRight != elemOID && resolvedRight != UNKNOWNOID {
		arrayType := ac.catalog.findArrayType(resolvedRight)
		if arrayType != 0 && arrayType != rightOID {
			right, err = ac.catalog.coerceToTargetType(right, rightOID, arrayType, 'i')
			if err != nil {
				return nil, err
			}
		}
	}

	return &ScalarArrayOpExpr{
		OpOID:  op.OID,
		OpName: op.Name,
		UseOr:  useOr,
		Left:   left,
		Right:  right,
	}, nil
}

// transformAExprIn transforms x IN (a, b, c) to ScalarArrayOpExpr.
//
// pg: src/backend/parser/parse_expr.c — transformAExprIn
func (ac *analyzeCtx) transformAExprIn(ae *nodes.A_Expr) (AnalyzedExpr, error) {
	left, err := ac.transformExpr(ae.Lexpr)
	if err != nil {
		return nil, err
	}

	// The right side is a List of expressions.
	opName := "="
	if ae.Name != nil && len(ae.Name.Items) > 0 {
		opName = stringVal(ae.Name.Items[len(ae.Name.Items)-1])
	}

	useOr := opName == "=" // IN → = ANY, NOT IN → <> ALL

	var rightExprs []AnalyzedExpr
	switch r := ae.Rexpr.(type) {
	case *nodes.List:
		for _, item := range r.Items {
			expr, err := ac.transformExpr(item)
			if err != nil {
				return nil, err
			}
			rightExprs = append(rightExprs, expr)
		}
	default:
		expr, err := ac.transformExpr(ae.Rexpr)
		if err != nil {
			return nil, err
		}
		rightExprs = append(rightExprs, expr)
	}

	// Determine common element type and build ArrayExprQ.
	// pg: src/backend/parser/parse_expr.c — transformAExprIn
	// PG includes the left expression: allexprs = list_concat(list_make1(lexpr), rnonvars)
	allExprs := make([]AnalyzedExpr, 0, 1+len(rightExprs))
	allExprs = append(allExprs, left)
	allExprs = append(allExprs, rightExprs...)
	elemType, err := ac.catalog.selectCommonType(allExprs, "")
	if err != nil {
		return nil, err
	}

	// Coerce elements to common type.
	for i, elem := range rightExprs {
		et := elem.exprType()
		if et != elemType {
			rightExprs[i], err = ac.catalog.coerceToTargetType(elem, et, elemType, 'i')
			if err != nil {
				return nil, err
			}
		}
	}

	arrayType := ac.catalog.findArrayType(elemType)
	arrayExpr := &ArrayExprQ{
		ElementType: elemType,
		ArrayType:   arrayType,
		Elements:    rightExprs,
	}

	return ac.makeScalarArrayOp(opName, useOr, left, arrayExpr)
}

// transformBetween transforms x BETWEEN a AND b (or NOT BETWEEN / SYMMETRIC).
// PG transforms BETWEEN to (x >= a AND x <= b).
//
// pg: src/backend/parser/parse_expr.c — transformAExprBetween
func (ac *analyzeCtx) transformBetween(ae *nodes.A_Expr) (AnalyzedExpr, error) {
	left, err := ac.transformExpr(ae.Lexpr)
	if err != nil {
		return nil, err
	}

	// Rexpr is a list of exactly 2 items: [low, high].
	rList, ok := ae.Rexpr.(*nodes.List)
	if !ok || len(rList.Items) != 2 {
		return nil, fmt.Errorf("BETWEEN requires exactly 2 arguments on the right side")
	}

	low, err := ac.transformExpr(rList.Items[0])
	if err != nil {
		return nil, err
	}
	high, err := ac.transformExpr(rList.Items[1])
	if err != nil {
		return nil, err
	}

	leftOID := left.exprType()
	lowOID := low.exprType()
	highOID := high.exprType()

	// Resolve >= operator.
	geOp, resolvedLeftGe, resolvedRightGe, err := ac.catalog.resolveOperator(">=", leftOID, lowOID, false)
	if err != nil {
		return nil, err
	}
	// Resolve <= operator.
	leOp, resolvedLeftLe, resolvedRightLe, err := ac.catalog.resolveOperator("<=", leftOID, highOID, false)
	if err != nil {
		return nil, err
	}

	// Coerce if needed.
	leftForGe := left
	if leftOID != resolvedLeftGe {
		leftForGe, err = ac.catalog.coerceToTargetType(left, leftOID, resolvedLeftGe, 'i')
		if err != nil {
			return nil, err
		}
	}
	if lowOID != resolvedRightGe {
		low, err = ac.catalog.coerceToTargetType(low, lowOID, resolvedRightGe, 'i')
		if err != nil {
			return nil, err
		}
	}

	leftForLe := left
	if leftOID != resolvedLeftLe {
		leftForLe, err = ac.catalog.coerceToTargetType(left, leftOID, resolvedLeftLe, 'i')
		if err != nil {
			return nil, err
		}
	}
	if highOID != resolvedRightLe {
		high, err = ac.catalog.coerceToTargetType(high, highOID, resolvedRightLe, 'i')
		if err != nil {
			return nil, err
		}
	}

	geExpr := &OpExpr{OpOID: geOp.OID, OpName: geOp.Name, ResultType: geOp.Result, Left: leftForGe, Right: low}
	leExpr := &OpExpr{OpOID: leOp.OID, OpName: leOp.Name, ResultType: leOp.Result, Left: leftForLe, Right: high}

	result := AnalyzedExpr(&BoolExprQ{Op: BoolAnd, Args: []AnalyzedExpr{geExpr, leExpr}})

	// NOT BETWEEN / NOT BETWEEN SYMMETRIC wraps in NOT.
	if ae.Kind == nodes.AEXPR_NOT_BETWEEN || ae.Kind == nodes.AEXPR_NOT_BETWEEN_SYM {
		result = &BoolExprQ{Op: BoolNot, Args: []AnalyzedExpr{result}}
	}

	return result, nil
}

// transformDistinct transforms IS [NOT] DISTINCT FROM.
//
// pg: src/backend/parser/parse_expr.c — transformAExprDistinct
func (ac *analyzeCtx) transformDistinct(ae *nodes.A_Expr) (AnalyzedExpr, error) {
	left, err := ac.transformExpr(ae.Lexpr)
	if err != nil {
		return nil, err
	}
	right, err := ac.transformExpr(ae.Rexpr)
	if err != nil {
		return nil, err
	}

	leftOID := left.exprType()
	rightOID := right.exprType()

	// Resolve the equality operator.
	op, resolvedLeft, resolvedRight, err := ac.catalog.resolveOperator("=", leftOID, rightOID, false)
	if err != nil {
		return nil, err
	}

	if leftOID != resolvedLeft {
		left, err = ac.catalog.coerceToTargetType(left, leftOID, resolvedLeft, 'i')
		if err != nil {
			return nil, err
		}
	}
	if rightOID != resolvedRight {
		right, err = ac.catalog.coerceToTargetType(right, rightOID, resolvedRight, 'i')
		if err != nil {
			return nil, err
		}
	}

	return &DistinctExprQ{
		OpOID:      op.OID,
		OpName:     op.Name,
		ResultType: BOOLOID,
		Left:       left,
		Right:      right,
		IsNot:      ae.Kind == nodes.AEXPR_NOT_DISTINCT,
	}, nil
}

// transformAExprOpAny transforms op ANY (array) expressions.
//
// pg: src/backend/parser/parse_expr.c — transformAExprOpAny
func (ac *analyzeCtx) transformAExprOpAny(ae *nodes.A_Expr) (AnalyzedExpr, error) {
	left, err := ac.transformExpr(ae.Lexpr)
	if err != nil {
		return nil, err
	}
	right, err := ac.transformExpr(ae.Rexpr)
	if err != nil {
		return nil, err
	}

	opName := ""
	if ae.Name != nil && len(ae.Name.Items) > 0 {
		opName = stringVal(ae.Name.Items[len(ae.Name.Items)-1])
	}

	return ac.makeScalarArrayOp(opName, true, left, right)
}

// transformAExprOpAll transforms op ALL (array) expressions.
//
// pg: src/backend/parser/parse_expr.c — transformAExprOpAll
func (ac *analyzeCtx) transformAExprOpAll(ae *nodes.A_Expr) (AnalyzedExpr, error) {
	left, err := ac.transformExpr(ae.Lexpr)
	if err != nil {
		return nil, err
	}
	right, err := ac.transformExpr(ae.Rexpr)
	if err != nil {
		return nil, err
	}

	opName := ""
	if ae.Name != nil && len(ae.Name.Items) > 0 {
		opName = stringVal(ae.Name.Items[len(ae.Name.Items)-1])
	}

	return ac.makeScalarArrayOp(opName, false, left, right)
}

// transformAArrayExpr transforms a raw ARRAY[...] constructor (from parser).
//
// pg: src/backend/parser/parse_expr.c — transformArrayExpr
func (ac *analyzeCtx) transformAArrayExpr(ae *nodes.A_ArrayExpr) (AnalyzedExpr, error) {
	var elems []AnalyzedExpr

	if ae.Elements != nil {
		for _, item := range ae.Elements.Items {
			expr, err := ac.transformExpr(item)
			if err != nil {
				return nil, err
			}
			elems = append(elems, expr)
		}
	}

	// Determine common element type.
	var elemType uint32
	if len(elems) > 0 {
		var err error
		elemType, err = ac.catalog.selectCommonType(elems, "ARRAY")
		if err != nil {
			return nil, err
		}
	} else {
		elemType = TEXTOID
	}

	// Coerce elements to common type.
	for i, elem := range elems {
		et := elem.exprType()
		if et != elemType {
			var err error
			elems[i], err = ac.catalog.coerceToTargetType(elem, et, elemType, 'i')
			if err != nil {
				return nil, err
			}
		}
	}

	arrayType := ac.catalog.findArrayType(elemType)

	return &ArrayExprQ{
		ElementType: elemType,
		ArrayType:   arrayType,
		Elements:    elems,
	}, nil
}

// transformArrayExpr transforms an ARRAY[...] constructor.
//
// pg: src/backend/parser/parse_expr.c — transformArrayExpr
func (ac *analyzeCtx) transformArrayExpr(ae *nodes.ArrayExpr) (AnalyzedExpr, error) {
	var elems []AnalyzedExpr

	if ae.Elements != nil {
		for _, item := range ae.Elements.Items {
			expr, err := ac.transformExpr(item)
			if err != nil {
				return nil, err
			}
			elems = append(elems, expr)
		}
	}

	// Determine common element type.
	var elemType uint32
	if len(elems) > 0 {
		var err error
		elemType, err = ac.catalog.selectCommonType(elems, "ARRAY")
		if err != nil {
			return nil, err
		}
	} else {
		elemType = TEXTOID
	}

	// Coerce elements to common type.
	for i, elem := range elems {
		et := elem.exprType()
		if et != elemType {
			var err error
			elems[i], err = ac.catalog.coerceToTargetType(elem, et, elemType, 'i')
			if err != nil {
				return nil, err
			}
		}
	}

	// Find array type for element type.
	arrayType := ac.catalog.findArrayType(elemType)

	return &ArrayExprQ{
		ElementType: elemType,
		ArrayType:   arrayType,
		Elements:    elems,
	}, nil
}

// transformRowExpr transforms a ROW(...) expression.
//
// pg: src/backend/parser/parse_expr.c — transformRowExpr
func (ac *analyzeCtx) transformRowExpr(re *nodes.RowExpr) (AnalyzedExpr, error) {
	var args []AnalyzedExpr
	if re.Args != nil {
		for _, item := range re.Args.Items {
			expr, err := ac.transformExpr(item)
			if err != nil {
				return nil, err
			}
			args = append(args, expr)
		}
	}

	var format byte
	if re.RowFormat == nodes.COERCE_EXPLICIT_CALL {
		format = 'e' // explicit ROW keyword
	}

	return &RowExprQ{
		Args:       args,
		ResultType: RECORDOID,
		RowFormat:  format,
	}, nil
}

// transformIndirection transforms an A_Indirection node into either
// a FieldSelectExprQ (for field access like rec.field) or a SubscriptingRefExpr
// (for array subscript like arr[1]).
//
// pg: src/backend/parser/parse_expr.c — transformIndirection
func (ac *analyzeCtx) transformIndirection(ind *nodes.A_Indirection) (AnalyzedExpr, error) {
	// Transform the base expression.
	result, err := ac.transformExpr(ind.Arg)
	if err != nil {
		return nil, err
	}

	if ind.Indirection == nil {
		return result, nil
	}

	for _, item := range ind.Indirection.Items {
		switch idx := item.(type) {
		case *nodes.A_Indices:
			// Array subscript: arr[expr] or arr[lower:upper]
			containerType := result.exprType()
			// Determine element type.
			elemType := ac.catalog.getArrayElementType(containerType)
			if elemType == 0 {
				elemType = containerType // fallback
			}

			if idx.IsSlice {
				// Slice: arr[lower:upper]
				var lower, upper AnalyzedExpr
				if idx.Lidx != nil {
					lower, err = ac.transformExpr(idx.Lidx)
					if err != nil {
						return nil, err
					}
				}
				if idx.Uidx != nil {
					upper, err = ac.transformExpr(idx.Uidx)
					if err != nil {
						return nil, err
					}
				}
				result = &SubscriptingRefExpr{
					ContainerExpr:  result,
					SubscriptExprs: []AnalyzedExpr{upper},
					LowerExprs:     []AnalyzedExpr{lower},
					ResultType:     containerType, // slice returns array type
					IsSlice:        true,
				}
			} else {
				// Simple subscript: arr[expr]
				sub, err := ac.transformExpr(idx.Uidx)
				if err != nil {
					return nil, err
				}
				result = &SubscriptingRefExpr{
					ContainerExpr:  result,
					SubscriptExprs: []AnalyzedExpr{sub},
					ResultType:     elemType,
					IsSlice:        false,
				}
			}
		case *nodes.String:
			// Field access: rec.field — already handled as FieldSelectExprQ
			fieldName := idx.Str
			resultType := result.exprType()
			// Try to resolve the field from composite type.
			fieldType, fieldTypMod, fieldColl, fieldNum := ac.catalog.resolveCompositeField(resultType, fieldName)
			result = &FieldSelectExprQ{
				Arg:        result,
				FieldNum:   int16(fieldNum),
				FieldName:  fieldName,
				ResultType: fieldType,
				TypeMod:    fieldTypMod,
				Collation:  fieldColl,
			}
		default:
			// A_Star or other — skip
		}
	}

	return result, nil
}

// transformNamedArgExpr transforms a named argument expression (arg_name => value).
//
// pg: src/backend/parser/parse_expr.c — transformExprRecurse (NamedArgExpr case)
func (ac *analyzeCtx) transformNamedArgExpr(nae *nodes.NamedArgExpr) (AnalyzedExpr, error) {
	arg, err := ac.transformExpr(nae.Arg)
	if err != nil {
		return nil, err
	}
	return &NamedArgExprQ{
		Name:   nae.Name,
		Arg:    arg,
		ArgNum: nae.Argnumber,
	}, nil
}

// transformParamRef transforms a $N parameter reference.
//
// pg: src/backend/parser/parse_param.c — variable_paramref_hook
func (ac *analyzeCtx) transformParamRef(pr *nodes.ParamRef) (AnalyzedExpr, error) {
	return &ParamExpr{
		ParamNum:  pr.Number,
		ParamType: UNKNOWNOID,
	}, nil
}

// transformCollateClause transforms a COLLATE clause.
//
// pg: src/backend/parser/parse_expr.c — transformCollateClause
func (ac *analyzeCtx) transformCollateClause(cc *nodes.CollateClause) (AnalyzedExpr, error) {
	arg, err := ac.transformExpr(cc.Arg)
	if err != nil {
		return nil, err
	}

	collName := ""
	if cc.Collname != nil && len(cc.Collname.Items) > 0 {
		collName = stringVal(cc.Collname.Items[len(cc.Collname.Items)-1])
	}

	return &CollateExprQ{
		Arg:      arg,
		CollName: collName,
	}, nil
}

// resolveWindowDef resolves a WINDOW clause reference, creating or finding the WindowClauseQ.
//
// pg: src/backend/parser/parse_clause.c — transformWindowDefinitions
func (ac *analyzeCtx) resolveWindowDef(overNode nodes.Node) (uint32, error) {
	wd, ok := overNode.(*nodes.WindowDef)
	if !ok {
		return 0, fmt.Errorf("unsupported OVER clause type: %T", overNode)
	}

	// If referencing a named window, find it.
	if wd.Refname != "" {
		for i, wc := range ac.query.WindowClause {
			if wc.Name == wd.Refname {
				return uint32(i), nil
			}
		}
		return 0, &Error{Code: CodeUndefinedObject, Message: fmt.Sprintf("window %q does not exist", wd.Refname)}
	}

	// Build WindowClauseQ.
	wc := &WindowClauseQ{
		Name:         wd.Name,
		FrameOptions: wd.FrameOptions,
	}

	// PARTITION BY.
	if wd.PartitionClause != nil {
		for _, item := range wd.PartitionClause.Items {
			tle, err := ac.findTargetlistEntry(item, ac.query.TargetList)
			if err != nil {
				return 0, err
			}
			ref := ac.assignSortGroupRef(tle, ac.query.TargetList)
			wc.PartitionBy = append(wc.PartitionBy, &SortGroupClause{
				TLESortGroupRef: ref,
			})
		}
	}

	// ORDER BY.
	if wd.OrderClause != nil {
		for _, item := range wd.OrderClause.Items {
			sb, ok := item.(*nodes.SortBy)
			if !ok {
				continue
			}
			tle, err := ac.findTargetlistEntry(sb.Node, ac.query.TargetList)
			if err != nil {
				return 0, err
			}
			ref := ac.assignSortGroupRef(tle, ac.query.TargetList)
			desc := sb.SortbyDir == nodes.SORTBY_DESC
			nullsFirst := false
			switch sb.SortbyNulls {
			case nodes.SORTBY_NULLS_FIRST:
				nullsFirst = true
			case nodes.SORTBY_NULLS_LAST:
				nullsFirst = false
			default:
				nullsFirst = desc
			}
			wc.OrderBy = append(wc.OrderBy, &SortGroupClause{
				TLESortGroupRef: ref,
				Descending:      desc,
				NullsFirst:      nullsFirst,
			})
		}
	}

	// Frame offsets.
	if wd.StartOffset != nil {
		var err error
		wc.StartOffset, err = ac.transformExpr(wd.StartOffset)
		if err != nil {
			return 0, err
		}
	}
	if wd.EndOffset != nil {
		var err error
		wc.EndOffset, err = ac.transformExpr(wd.EndOffset)
		if err != nil {
			return 0, err
		}
	}

	// Check if this matches an existing unnamed window clause.
	idx := uint32(len(ac.query.WindowClause))
	ac.query.WindowClause = append(ac.query.WindowClause, wc)
	return idx, nil
}

// analyzeCTEs processes a WITH clause.
//
// pg: src/backend/parser/parse_cte.c — analyzeCTEList
func (ac *analyzeCtx) analyzeCTEs(withClause *nodes.WithClause) error {
	ac.query.IsRecursive = withClause.Recursive

	if withClause.Ctes == nil {
		return nil
	}

	for _, item := range withClause.Ctes.Items {
		cte, ok := item.(*nodes.CommonTableExpr)
		if !ok {
			continue
		}

		sub, ok := cte.Ctequery.(*nodes.SelectStmt)
		if !ok {
			return fmt.Errorf("CTE query is not a SELECT")
		}

		// Collect column aliases.
		var aliases []string
		if cte.Aliascolnames != nil {
			for _, item := range cte.Aliascolnames.Items {
				aliases = append(aliases, stringVal(item))
			}
		}

		// pgparser sets Cterecursive=false on individual CTEs; only
		// WithClause.Recursive is set by the parser. PG's analyzer determines
		// whether each CTE actually self-references. We use a simple heuristic:
		// if the WITH clause is RECURSIVE and the CTE body is a UNION, treat
		// it as potentially recursive.
		isRecursive := withClause.Recursive && sub.Op != nodes.SETOP_NONE && sub.Larg != nil
		if isRecursive {
			// Recursive CTE: two-phase analysis.
			// pg: src/backend/parser/analyze.c — determineRecursiveColTypes
			//
			// 1. The CTE query must be a UNION (PG requires this).
			if sub.Op == nodes.SETOP_NONE || sub.Larg == nil {
				return fmt.Errorf("recursive query \"%s\" does not have the form non-recursive-term UNION [ALL] recursive-term", cte.Ctename)
			}

			// 2. Analyze only the non-recursive term (larg).
			nrQuery, err := ac.catalog.analyzeSelectStmt(sub.Larg)
			if err != nil {
				return err
			}

			// 3. Build a temporary CTE entry with non-recursive term types
			//    so the recursive term can reference it.
			tmpCTE := &CommonTableExprQ{
				Name:      cte.Ctename,
				Aliases:   aliases,
				Query:     nrQuery,
				Recursive: true,
			}
			ac.query.CTEList = append(ac.query.CTEList, tmpCTE)

			// 4. Set visibleCTEs so nested analyzeSelectStmt calls can see
			//    the partially-defined CTE.
			ac.catalog.visibleCTEs = ac.query.CTEList

			// 5. Now analyze the full UNION — the recursive term can see the CTE.
			fullQuery, err := ac.catalog.analyzeSetOp(sub)

			// Clear visibleCTEs regardless of error.
			ac.catalog.visibleCTEs = nil

			if err != nil {
				return err
			}

			// 6. Update the CTE entry with the full query.
			tmpCTE.Query = fullQuery
			tmpCTE.Materialized = cte.Ctematerialized
		} else {
			// Non-recursive CTE: analyze normally.
			subQuery, err := ac.catalog.analyzeSelectStmt(sub)
			if err != nil {
				return err
			}

			cteQ := &CommonTableExprQ{
				Name:         cte.Ctename,
				Aliases:      aliases,
				Query:        subQuery,
				Recursive:    cte.Cterecursive,
				Materialized: cte.Ctematerialized,
			}
			ac.query.CTEList = append(ac.query.CTEList, cteQ)
		}
	}

	return nil
}

// transformCTERef creates a range table entry for a CTE reference.
//
// pg: src/backend/parser/parse_clause.c — getRTEForSpecialRelationTypes (CTE case)
func (ac *analyzeCtx) transformCTERef(cteName, alias string, cteIdx int) (JoinNode, error) {
	cte := ac.query.CTEList[cteIdx]
	subQuery := cte.Query

	eref := cteName
	if alias != "" {
		eref = alias
	}

	colNames := make([]string, 0, len(subQuery.TargetList))
	colTypes := make([]uint32, 0, len(subQuery.TargetList))
	colTypMods := make([]int32, 0, len(subQuery.TargetList))
	colCollations := make([]uint32, 0, len(subQuery.TargetList))
	for i, te := range subQuery.TargetList {
		if te.ResJunk {
			continue
		}
		name := te.ResName
		if i < len(cte.Aliases) {
			name = cte.Aliases[i]
		}
		colNames = append(colNames, name)
		colTypes = append(colTypes, te.Expr.exprType())
		colTypMods = append(colTypMods, int32(-1))
		colCollations = append(colCollations, te.Expr.exprCollation())
	}

	rtIdx := len(ac.query.RangeTable)
	rte := &RangeTableEntry{
		Kind:          RTECTE,
		Alias:         alias,
		ERef:          eref,
		ColNames:      colNames,
		ColTypes:      colTypes,
		ColTypMods:    colTypMods,
		ColCollations: colCollations,
		CTEName:       cteName,
		CTEIndex:      cteIdx,
	}
	ac.query.RangeTable = append(ac.query.RangeTable, rte)

	return &RangeTableRef{RTIndex: rtIdx}, nil
}

// transformVisibleCTERef creates a range table entry for a CTE reference found
// in the catalog's visibleCTEs list (used during recursive CTE analysis).
//
// (pgddl helper — bridges the gap when the recursive term is analyzed in a
// fresh analyzeCtx that doesn't have the parent's CTEList.)
func (ac *analyzeCtx) transformVisibleCTERef(cte *CommonTableExprQ, cteName, alias string, cteIdx int) (JoinNode, error) {
	subQuery := cte.Query

	eref := cteName
	if alias != "" {
		eref = alias
	}

	colNames := make([]string, 0, len(subQuery.TargetList))
	colTypes := make([]uint32, 0, len(subQuery.TargetList))
	colTypMods := make([]int32, 0, len(subQuery.TargetList))
	colCollations := make([]uint32, 0, len(subQuery.TargetList))
	for i, te := range subQuery.TargetList {
		if te.ResJunk {
			continue
		}
		name := te.ResName
		if i < len(cte.Aliases) {
			name = cte.Aliases[i]
		}
		colNames = append(colNames, name)
		colTypes = append(colTypes, te.Expr.exprType())
		colTypMods = append(colTypMods, int32(-1))
		colCollations = append(colCollations, te.Expr.exprCollation())
	}

	rtIdx := len(ac.query.RangeTable)
	rte := &RangeTableEntry{
		Kind:          RTECTE,
		Alias:         alias,
		ERef:          eref,
		ColNames:      colNames,
		ColTypes:      colTypes,
		ColTypMods:    colTypMods,
		ColCollations: colCollations,
		CTEName:       cteName,
		CTEIndex:      cteIdx,
	}
	ac.query.RangeTable = append(ac.query.RangeTable, rte)

	return &RangeTableRef{RTIndex: rtIdx}, nil
}

// findArrayType returns the array type OID for a given element type.
// Returns 0 if no array type is found.
//
// (pgddl helper — PG uses get_array_type)
func (c *Catalog) findArrayType(elemType uint32) uint32 {
	if bt := c.typeByOID[elemType]; bt != nil && bt.Array != 0 {
		return bt.Array
	}
	return 0
}

// opNameToRowCompareType maps operator names to RowCompare type constants.
func opNameToRowCompareType(opName string) int {
	switch opName {
	case "<":
		return RowCompareLT
	case "<=":
		return RowCompareLE
	case "=":
		return RowCompareEQ
	case ">=":
		return RowCompareGE
	case ">":
		return RowCompareGT
	case "<>", "!=":
		return RowCompareNE
	default:
		return 0
	}
}

// rowCompareOpName returns the operator name for a RowCompare type.
func rowCompareOpName(rcType int) string {
	switch rcType {
	case RowCompareLT:
		return "<"
	case RowCompareLE:
		return "<="
	case RowCompareEQ:
		return "="
	case RowCompareGE:
		return ">="
	case RowCompareGT:
		return ">"
	case RowCompareNE:
		return "<>"
	default:
		return "="
	}
}

// getArrayElementType returns the element type OID for an array type, or 0 if not an array.
func (c *Catalog) getArrayElementType(arrayType uint32) uint32 {
	bt := c.typeByOID[arrayType]
	if bt == nil {
		return 0
	}
	if bt.Category == 'A' && bt.Elem != 0 {
		return bt.Elem
	}
	return 0
}

// resolveCompositeField resolves a field name within a composite type.
// Returns the field's type OID, type modifier, collation, and field number (1-based).
// If the composite type or field cannot be resolved, returns UNKNOWNOID with reasonable defaults.
func (c *Catalog) resolveCompositeField(compositeType uint32, fieldName string) (uint32, int32, uint32, int) {
	// Look up composite type's relation.
	bt := c.typeByOID[compositeType]
	if bt == nil {
		return UNKNOWNOID, -1, 0, 0
	}
	// Find the relation that defines this composite type.
	rel := c.typeRelation(compositeType)
	if rel == nil {
		return UNKNOWNOID, -1, 0, 0
	}
	for i, col := range rel.Columns {
		if col.Name == fieldName {
			return col.TypeOID, col.TypeMod, col.Collation, i + 1
		}
	}
	return UNKNOWNOID, -1, 0, 0
}

// typeRelation looks up the relation that defines a composite type.
func (c *Catalog) typeRelation(typeOID uint32) *Relation {
	for _, s := range c.schemas {
		for _, r := range s.Relations {
			if r.RowTypeOID == typeOID {
				return r
			}
		}
	}
	return nil
}

// --- Helper functions ---

// coerceToTargetType creates a coercion node for the given pathway.
func (c *Catalog) coerceToTargetType(expr AnalyzedExpr, srcType, targetType uint32, context byte) (AnalyzedExpr, error) {
	if srcType == targetType {
		return expr, nil
	}
	if srcType == UNKNOWNOID {
		// pg: src/backend/parser/parse_coerce.c — coerce_type
		// For UNKNOWN Const, PG directly changes the Const's type rather than
		// wrapping in a coercion node.
		if c, ok := expr.(*ConstExpr); ok {
			return &ConstExpr{
				TypeOID: targetType,
				TypeMod: -1,
				IsNull:  c.IsNull,
				Value:   c.Value,
			}, nil
		}
		return &RelabelExpr{
			Arg:        expr,
			ResultType: targetType,
			TypeMod:    -1,
			Format:     context,
		}, nil
	}

	pathway, funcOID := c.FindCoercionPathway(srcType, targetType, context)
	switch pathway {
	case CoercionRelabel:
		return &RelabelExpr{
			Arg:        expr,
			ResultType: targetType,
			TypeMod:    -1,
			Format:     context,
		}, nil
	case CoercionFunc:
		proc := c.procByOID[funcOID]
		funcName := ""
		if proc != nil {
			funcName = proc.Name
		}
		return &FuncCallExpr{
			FuncOID:      funcOID,
			FuncName:     funcName,
			ResultType:   targetType,
			ResultTypMod: -1,
			Args:         []AnalyzedExpr{expr},
			CoerceFormat: context,
		}, nil
	case CoercionIO:
		return &CoerceViaIOExpr{
			Arg:        expr,
			ResultType: targetType,
			Format:     context,
		}, nil
	default:
		return nil, errDatatypeMismatch(fmt.Sprintf(
			"cannot cast type %s to %s",
			c.typeName(srcType), c.typeName(targetType),
		))
	}
}

// resolveOperator resolves an operator expression and returns the operator
// along with the resolved types for left and right operands.
func (c *Catalog) resolveOperator(name string, leftType, rightType uint32, isPrefix bool) (*BuiltinOperator, uint32, uint32, error) {
	if isPrefix {
		// Prefix operator.
		if ops := c.operByKey[operKey{name: name, left: 0, right: rightType}]; len(ops) > 0 {
			return ops[0], 0, rightType, nil
		}
		if rightType == UNKNOWNOID {
			for _, commonType := range []uint32{TEXTOID, INT4OID, FLOAT8OID} {
				if ops := c.operByKey[operKey{name: name, left: 0, right: commonType}]; len(ops) > 0 {
					return ops[0], 0, commonType, nil
				}
			}
		}
		return nil, 0, 0, errUndefinedFunction(name, []uint32{rightType})
	}

	// Binary operator: exact match.
	if ops := c.operByKey[operKey{name: name, left: leftType, right: rightType}]; len(ops) > 0 {
		return ops[0], leftType, rightType, nil
	}

	// If one side is UNKNOWN, try substituting.
	if leftType == UNKNOWNOID && rightType != UNKNOWNOID {
		if ops := c.operByKey[operKey{name: name, left: rightType, right: rightType}]; len(ops) > 0 {
			return ops[0], rightType, rightType, nil
		}
		// pg: src/backend/parser/parse_func.c — func_select_candidate (step 4d)
		// Try the preferred type of the same category.
		if pref := c.preferredTypeForCategory(rightType); pref != 0 && pref != rightType {
			if ops := c.operByKey[operKey{name: name, left: pref, right: pref}]; len(ops) > 0 {
				return ops[0], pref, pref, nil
			}
		}
	}
	if rightType == UNKNOWNOID && leftType != UNKNOWNOID {
		if ops := c.operByKey[operKey{name: name, left: leftType, right: leftType}]; len(ops) > 0 {
			return ops[0], leftType, leftType, nil
		}
		// pg: src/backend/parser/parse_func.c — func_select_candidate (step 4d)
		if pref := c.preferredTypeForCategory(leftType); pref != 0 && pref != leftType {
			if ops := c.operByKey[operKey{name: name, left: pref, right: pref}]; len(ops) > 0 {
				return ops[0], pref, pref, nil
			}
		}
	}
	if leftType == UNKNOWNOID && rightType == UNKNOWNOID {
		if ops := c.operByKey[operKey{name: name, left: TEXTOID, right: TEXTOID}]; len(ops) > 0 {
			return ops[0], TEXTOID, TEXTOID, nil
		}
	}

	// Try implicit coercion.
	return c.resolveOpWithCoercionFull(name, leftType, rightType)
}

// funcCandidate represents a candidate during overload resolution.
//
// pg: src/backend/parser/parse_func.c — FuncCandidateList
type funcCandidate struct {
	argTypes []uint32
	op       *BuiltinOperator // set for operator candidates
	proc     *BuiltinProc     // set for function candidates
}

// canCoerceArgs checks whether inputTypes can be implicitly coerced to candidateTypes.
// UNKNOWN inputs are coercible to anything; polymorphic params use IsBinaryCoercible
// to check type-specific compatibility (e.g., ANYENUM only matches enum types).
//
// pg: src/backend/parser/parse_func.c — func_match_argtypes (per-candidate check)
func (c *Catalog) canCoerceArgs(inputTypes, candidateTypes []uint32) bool {
	for i := range inputTypes {
		if inputTypes[i] == UNKNOWNOID {
			continue
		}
		if inputTypes[i] == candidateTypes[i] {
			continue
		}
		// pg: func_match_argtypes uses can_coerce_type which internally
		// calls IsBinaryCoercible for polymorphic types.
		if isPolymorphic(candidateTypes[i]) {
			if c.IsBinaryCoercible(inputTypes[i], candidateTypes[i]) {
				continue
			}
			return false
		}
		if !c.CanCoerce(inputTypes[i], candidateTypes[i], 'i') {
			return false
		}
	}
	return true
}

// funcSelectCandidate resolves ambiguous function/operator overloads using
// PG's multi-phase heuristic algorithm.
//
// pg: src/backend/parser/parse_func.c — func_select_candidate
func (c *Catalog) funcSelectCandidate(inputTypes []uint32, candidates []funcCandidate) *funcCandidate {
	nargs := len(inputTypes)
	if len(candidates) <= 1 {
		if len(candidates) == 1 {
			return &candidates[0]
		}
		return nil
	}

	// Phase 1: Domain unwrap + count unknowns.
	// pg: func_select_candidate lines 1051-1062
	inputBase := make([]uint32, nargs)
	nunknowns := 0
	for i := 0; i < nargs; i++ {
		if inputTypes[i] != UNKNOWNOID {
			inputBase[i] = c.getBaseType(inputTypes[i])
		} else {
			inputBase[i] = UNKNOWNOID
			nunknowns++
		}
	}

	// Phase 2: Keep candidates with the most exact type matches at non-UNKNOWN positions.
	// pg: func_select_candidate lines 1068-1106
	nbestMatch := -1
	var filtered []funcCandidate
	for ci := range candidates {
		nmatch := 0
		for i := 0; i < nargs; i++ {
			if inputBase[i] != UNKNOWNOID && candidates[ci].argTypes[i] == inputBase[i] {
				nmatch++
			}
		}
		if nmatch > nbestMatch {
			nbestMatch = nmatch
			filtered = []funcCandidate{candidates[ci]}
		} else if nmatch == nbestMatch {
			filtered = append(filtered, candidates[ci])
		}
	}
	candidates = filtered
	if len(candidates) == 1 {
		return &candidates[0]
	}

	// Phase 3: Prefer candidates with preferred types in the same category.
	// pg: func_select_candidate lines 1115-1155
	slotCategory := make([]byte, nargs)
	for i := 0; i < nargs; i++ {
		if bt := c.typeByOID[inputBase[i]]; bt != nil {
			slotCategory[i] = bt.Category
		}
	}

	nbestMatch = -1
	filtered = nil
	for ci := range candidates {
		nmatch := 0
		for i := 0; i < nargs; i++ {
			if inputBase[i] != UNKNOWNOID {
				if candidates[ci].argTypes[i] == inputBase[i] {
					nmatch++
				} else if ct := c.typeByOID[candidates[ci].argTypes[i]]; ct != nil &&
					ct.IsPreferred && ct.Category == slotCategory[i] {
					nmatch++
				}
			}
		}
		if nmatch > nbestMatch {
			nbestMatch = nmatch
			filtered = []funcCandidate{candidates[ci]}
		} else if nmatch == nbestMatch {
			filtered = append(filtered, candidates[ci])
		}
	}
	candidates = filtered
	if len(candidates) == 1 {
		return &candidates[0]
	}

	// Phase 4: Resolve unknown argument positions using type category heuristics.
	// STRING category wins on conflict (unknown literals look like strings).
	// pg: func_select_candidate lines 1163-1300
	if nunknowns == 0 {
		return nil // no more heuristics without unknowns
	}

	slotCategoryU := make([]byte, nargs)
	slotHasPreferred := make([]bool, nargs)
	resolvedUnknowns := false
	for i := 0; i < nargs; i++ {
		if inputBase[i] != UNKNOWNOID {
			continue
		}
		resolvedUnknowns = true
		slotCategoryU[i] = 0 // TYPCATEGORY_INVALID
		slotHasPreferred[i] = false
		haveConflict := false
		for ci := range candidates {
			ct := c.typeByOID[candidates[ci].argTypes[i]]
			if ct == nil {
				continue
			}
			curCategory := ct.Category
			curIsPreferred := ct.IsPreferred
			if slotCategoryU[i] == 0 {
				// First candidate for this position.
				slotCategoryU[i] = curCategory
				slotHasPreferred[i] = curIsPreferred
			} else if curCategory == slotCategoryU[i] {
				// Same category — merge preferred.
				slotHasPreferred[i] = slotHasPreferred[i] || curIsPreferred
			} else {
				// Category conflict.
				if curCategory == 'S' {
					// STRING always wins.
					slotCategoryU[i] = curCategory
					slotHasPreferred[i] = curIsPreferred
				} else {
					haveConflict = true
				}
			}
		}
		if haveConflict && slotCategoryU[i] != 'S' {
			resolvedUnknowns = false
			break
		}
	}

	if resolvedUnknowns {
		// Filter candidates that don't match resolved categories.
		filtered = nil
		for ci := range candidates {
			keepit := true
			for i := 0; i < nargs; i++ {
				if inputBase[i] != UNKNOWNOID {
					continue
				}
				ct := c.typeByOID[candidates[ci].argTypes[i]]
				if ct == nil {
					keepit = false
					break
				}
				if ct.Category != slotCategoryU[i] {
					keepit = false
					break
				}
				if slotHasPreferred[i] && !ct.IsPreferred {
					keepit = false
					break
				}
			}
			if keepit {
				filtered = append(filtered, candidates[ci])
			}
		}
		if len(filtered) > 0 {
			candidates = filtered
		}
		if len(candidates) == 1 {
			return &candidates[0]
		}
	}

	// Phase 5: Last gasp — if all known-type inputs have the same base type,
	// assume unknowns are that type and see if exactly one candidate matches.
	// pg: func_select_candidate lines 1312-1357
	if nunknowns < nargs {
		knownType := uint32(UNKNOWNOID)
		for i := 0; i < nargs; i++ {
			if inputBase[i] == UNKNOWNOID {
				continue
			}
			if knownType == UNKNOWNOID {
				knownType = inputBase[i]
			} else if knownType != inputBase[i] {
				knownType = UNKNOWNOID
				break
			}
		}

		if knownType != UNKNOWNOID {
			// Pretend all inputs are knownType.
			allKnown := make([]uint32, nargs)
			for i := range allKnown {
				allKnown[i] = knownType
			}
			var lastMatch *funcCandidate
			nmatches := 0
			for ci := range candidates {
				if c.canCoerceArgs(allKnown, candidates[ci].argTypes) {
					nmatches++
					if nmatches > 1 {
						break // not unique
					}
					lastMatch = &candidates[ci]
				}
			}
			if nmatches == 1 {
				return lastMatch
			}
		}
	}

	return nil // failed to select a best candidate
}

// resolveOpWithCoercionFull resolves a binary operator by collecting all
// candidates with the matching name, filtering to coercible ones, and
// delegating to funcSelectCandidate.
//
// pg: src/backend/parser/parse_oper.c — oper_select_candidate (via func_select_candidate)
func (c *Catalog) resolveOpWithCoercionFull(name string, leftType, rightType uint32) (*BuiltinOperator, uint32, uint32, error) {
	inputTypes := []uint32{leftType, rightType}

	// 1. Collect all binary operators with matching name.
	var allCandidates []funcCandidate
	for key, ops := range c.operByKey {
		if key.name != name {
			continue
		}
		for _, op := range ops {
			if op.Kind == 'l' {
				continue // skip prefix operators
			}
			allCandidates = append(allCandidates, funcCandidate{
				argTypes: []uint32{op.Left, op.Right},
				op:       op,
			})
		}
	}

	// 2. Filter to coercible candidates (func_match_argtypes equivalent).
	var coercible []funcCandidate
	for _, cand := range allCandidates {
		if c.canCoerceArgs(inputTypes, cand.argTypes) {
			coercible = append(coercible, cand)
		}
	}
	if len(coercible) == 0 {
		return nil, 0, 0, errUndefinedFunction(name, []uint32{leftType, rightType})
	}
	if len(coercible) == 1 {
		best := &coercible[0]
		resolvedLeft := best.op.Left
		resolvedRight := best.op.Right
		if isPolymorphic(resolvedLeft) {
			resolvedLeft = leftType
		}
		if isPolymorphic(resolvedRight) {
			resolvedRight = rightType
		}
		return best.op, resolvedLeft, resolvedRight, nil
	}

	// 3. Delegate to funcSelectCandidate for multi-phase resolution.
	best := c.funcSelectCandidate(inputTypes, coercible)
	if best == nil {
		return nil, 0, 0, errUndefinedFunction(name, []uint32{leftType, rightType})
	}

	// Resolve polymorphic types: replace with actual input types.
	resolvedLeft := best.op.Left
	resolvedRight := best.op.Right
	if isPolymorphic(resolvedLeft) {
		resolvedLeft = leftType
	}
	if isPolymorphic(resolvedRight) {
		resolvedRight = rightType
	}
	return best.op, resolvedLeft, resolvedRight, nil
}

// resolveFuncOverload resolves function overloading and returns the best match proc
// and the resolved argument types (after resolving UNKNOWN).
//
// pg: src/backend/parser/parse_func.c — FuncnameGetCandidates + func_select_candidate
func (c *Catalog) resolveFuncOverload(procs []*BuiltinProc, argTypes []uint32, hasStar bool) (*BuiltinProc, []uint32, error) {
	// 1. Build candidates filtered by arg count.
	var candidates []funcCandidate
	for _, p := range procs {
		if int(p.NArgs) != len(argTypes) {
			continue
		}
		candidates = append(candidates, funcCandidate{
			argTypes: p.ArgTypes[:p.NArgs],
			proc:     p,
		})
	}

	if len(candidates) == 0 {
		return nil, nil, errUndefinedFunction(procs[0].Name, argTypes)
	}

	// 2. Check for exact match first.
	for ci := range candidates {
		exact := true
		for i, argOID := range argTypes {
			if argOID != candidates[ci].argTypes[i] {
				exact = false
				break
			}
		}
		if exact {
			return candidates[ci].proc, buildResolvedArgs(argTypes, candidates[ci].proc), nil
		}
	}

	// 3. Filter to coercible candidates (func_match_argtypes).
	var coercible []funcCandidate
	for _, cand := range candidates {
		if c.canCoerceArgs(argTypes, cand.argTypes) {
			coercible = append(coercible, cand)
		}
	}
	if len(coercible) == 0 {
		return nil, nil, errUndefinedFunction(procs[0].Name, argTypes)
	}
	if len(coercible) == 1 {
		return coercible[0].proc, buildResolvedArgs(argTypes, coercible[0].proc), nil
	}

	// 4. Delegate to funcSelectCandidate for multi-phase resolution.
	best := c.funcSelectCandidate(argTypes, coercible)
	if best == nil {
		return nil, nil, errUndefinedFunction(procs[0].Name, argTypes)
	}
	return best.proc, buildResolvedArgs(argTypes, best.proc), nil
}

// buildResolvedArgs constructs the resolved argument type slice from the matched proc.
func buildResolvedArgs(argTypes []uint32, proc *BuiltinProc) []uint32 {
	resolved := make([]uint32, len(argTypes))
	for i := range argTypes {
		if i < int(proc.NArgs) {
			resolved[i] = proc.ArgTypes[i]
		} else {
			resolved[i] = argTypes[i]
		}
	}
	return resolved
}

// figureColName determines the output column name for an expression.
//
// pg: src/backend/parser/parse_target.c — FigureColname
func figureColName(n nodes.Node) string {
	if n == nil {
		return "?column?"
	}
	switch v := n.(type) {
	case *nodes.ColumnRef:
		if v.Fields != nil && len(v.Fields.Items) > 0 {
			last := v.Fields.Items[len(v.Fields.Items)-1]
			if s, ok := last.(*nodes.String); ok {
				return s.Str
			}
		}
		return "?column?"
	case *nodes.FuncCall:
		if v.Funcname != nil && len(v.Funcname.Items) > 0 {
			return stringVal(v.Funcname.Items[len(v.Funcname.Items)-1])
		}
		return "?column?"
	case *nodes.TypeCast:
		// Try inner expression first, then fall back to type name.
		inner := figureColName(v.Arg)
		if inner != "?column?" {
			return inner
		}
		if v.TypeName != nil {
			_, name := typeNameParts(v.TypeName)
			if name != "" {
				return name
			}
		}
		return "?column?"
	case *nodes.A_Expr:
		return "?column?"
	case *nodes.SubLink:
		// scalar subquery — use the subquery's first column name
		if sub, ok := v.Subselect.(*nodes.SelectStmt); ok {
			if sub.TargetList != nil && len(sub.TargetList.Items) > 0 {
				if rt, ok := sub.TargetList.Items[0].(*nodes.ResTarget); ok {
					if rt.Name != "" {
						return rt.Name
					}
					return figureColName(rt.Val)
				}
			}
		}
		return "?column?"
	case *nodes.CaseExpr:
		return "case"
	case *nodes.CoalesceExpr:
		return "coalesce"
	case *nodes.MinMaxExpr:
		if v.Op == nodes.IS_GREATEST {
			return "greatest"
		}
		return "least"
	case *nodes.NullIfExpr:
		return "nullif"
	case *nodes.SQLValueFunction:
		switch v.Op {
		case nodes.SVFOP_CURRENT_DATE:
			return "current_date"
		case nodes.SVFOP_CURRENT_TIME, nodes.SVFOP_CURRENT_TIME_N:
			return "current_time"
		case nodes.SVFOP_CURRENT_TIMESTAMP, nodes.SVFOP_CURRENT_TIMESTAMP_N:
			return "current_timestamp"
		case nodes.SVFOP_LOCALTIME, nodes.SVFOP_LOCALTIME_N:
			return "localtime"
		case nodes.SVFOP_LOCALTIMESTAMP, nodes.SVFOP_LOCALTIMESTAMP_N:
			return "localtimestamp"
		case nodes.SVFOP_CURRENT_ROLE:
			return "current_role"
		case nodes.SVFOP_CURRENT_USER:
			return "current_user"
		case nodes.SVFOP_USER:
			return "current_user"
		case nodes.SVFOP_SESSION_USER:
			return "session_user"
		case nodes.SVFOP_CURRENT_CATALOG:
			return "current_catalog"
		case nodes.SVFOP_CURRENT_SCHEMA:
			return "current_schema"
		}
		return "?column?"
	case *nodes.BooleanTest:
		return figureColName(v.Arg)
	case *nodes.A_Const:
		return "?column?"
	case *nodes.ArrayExpr:
		return "array"
	case *nodes.RowExpr:
		return "row"
	case *nodes.CollateClause:
		return figureColName(v.Arg)
	case *nodes.XmlExpr:
		switch v.Op {
		case nodes.IS_XMLCONCAT:
			return "xmlconcat"
		case nodes.IS_XMLELEMENT:
			return "xmlelement"
		case nodes.IS_XMLFOREST:
			return "xmlforest"
		case nodes.IS_XMLPARSE:
			return "xmlparse"
		case nodes.IS_XMLPI:
			return "xmlpi"
		case nodes.IS_XMLROOT:
			return "xmlroot"
		case nodes.IS_XMLSERIALIZE:
			return "xmlserialize"
		}
		return "?column?"
	case *nodes.XmlSerialize:
		return "xmlserialize"
	case *nodes.GroupingFunc:
		return "grouping"
	case *nodes.A_Indirection:
		if v.Indirection != nil && len(v.Indirection.Items) > 0 {
			last := v.Indirection.Items[len(v.Indirection.Items)-1]
			if s, ok := last.(*nodes.String); ok {
				return s.Str
			}
		}
		return "?column?"
	default:
		return "?column?"
	}
}

// isStar checks if a ColumnRef is a star expression.
func isStar(cr *nodes.ColumnRef) bool {
	if cr.Fields == nil {
		return false
	}
	for _, item := range cr.Fields.Items {
		if _, ok := item.(*nodes.A_Star); ok {
			return true
		}
	}
	return false
}

// starTableName extracts the table name from a qualified star expression (table.*).
func starTableName(cr *nodes.ColumnRef) string {
	if cr.Fields == nil || len(cr.Fields.Items) < 2 {
		return ""
	}
	return stringVal(cr.Fields.Items[0])
}

// collectJoinColumns collects column names/types/collations from a join node's constituent RTEs.
func collectJoinColumns(rangeTable []*RangeTableEntry, jn JoinNode, names *[]string, types *[]uint32, typMods *[]int32, colls *[]uint32) {
	switch v := jn.(type) {
	case *RangeTableRef:
		rte := rangeTable[v.RTIndex]
		*names = append(*names, rte.ColNames...)
		*types = append(*types, rte.ColTypes...)
		*typMods = append(*typMods, rte.ColTypMods...)
		if colls != nil {
			*colls = append(*colls, rte.ColCollations...)
		}
	case *JoinExprNode:
		// Use the RTE's already-computed column list, which correctly
		// excludes deduplicated USING columns for JOIN USING.
		rte := rangeTable[v.RTIndex]
		*names = append(*names, rte.ColNames...)
		*types = append(*types, rte.ColTypes...)
		*typMods = append(*typMods, rte.ColTypMods...)
		if colls != nil {
			*colls = append(*colls, rte.ColCollations...)
		}
	}
}

// collectJoinColumnsExcluding collects columns from a join node, skipping
// columns whose names appear in the exclude list. This implements the
// PostgreSQL behavior where USING columns from the right side are omitted.
// pg: src/backend/parser/parse_clause.c — extractRemainingColumns
func collectJoinColumnsExcluding(rangeTable []*RangeTableEntry, jn JoinNode, names *[]string, types *[]uint32, typMods *[]int32, colls *[]uint32, exclude []string) {
	excludeSet := make(map[string]bool, len(exclude))
	for _, name := range exclude {
		excludeSet[name] = true
	}

	var allNames []string
	var allTypes []uint32
	var allTypMods []int32
	var allColls []uint32
	collectJoinColumns(rangeTable, jn, &allNames, &allTypes, &allTypMods, &allColls)

	for i, name := range allNames {
		if excludeSet[name] {
			// Skip this USING column (only skip the first occurrence).
			delete(excludeSet, name)
			continue
		}
		*names = append(*names, name)
		*types = append(*types, allTypes[i])
		*typMods = append(*typMods, allTypMods[i])
		if colls != nil && i < len(allColls) {
			*colls = append(*colls, allColls[i])
		}
	}
}

// convertJoinType converts a pgparser JoinType to the catalog JoinType.
func convertJoinType(jt nodes.JoinType, isNatural bool) JoinType {
	switch jt {
	case nodes.JOIN_LEFT:
		return JoinLeft
	case nodes.JOIN_RIGHT:
		return JoinRight
	case nodes.JOIN_FULL:
		return JoinFull
	default:
		if isNatural && jt == nodes.JOIN_INNER {
			return JoinCross
		}
		return JoinInner
	}
}

// convertSetOpType converts pgparser SetOperation to the catalog SetOpType.
func convertSetOpType(op nodes.SetOperation) SetOpType {
	switch op {
	case nodes.SETOP_UNION:
		return SetOpUnion
	case nodes.SETOP_INTERSECT:
		return SetOpIntersect
	case nodes.SETOP_EXCEPT:
		return SetOpExcept
	default:
		return SetOpNone
	}
}

// findTargetlistEntry finds (or creates) a target list entry for a GROUP BY/ORDER BY expression.
// Supports SQL92 (integer position, bare column name) and SQL99 (expression matching).
//
// pg: src/backend/parser/parse_clause.c — findTargetlistEntrySQL92 + findTargetlistEntrySQL99
func (ac *analyzeCtx) findTargetlistEntry(n nodes.Node, tlist []*TargetEntry) (*TargetEntry, error) {
	// SQL92: integer constant → positional reference.
	if ac2, ok := n.(*nodes.A_Const); ok && !ac2.Isnull {
		if iv, ok := ac2.Val.(*nodes.Integer); ok {
			pos := int(iv.Ival)
			// Count non-junk entries.
			idx := 0
			for _, te := range tlist {
				if te.ResJunk {
					continue
				}
				idx++
				if idx == pos {
					return te, nil
				}
			}
			return nil, &Error{
				Code:    CodeFeatureNotSupported,
				Message: fmt.Sprintf("ORDER/GROUP BY position %d is not in select list", pos),
			}
		}
	}

	// SQL92: bare column name → match by ResName.
	if cr, ok := n.(*nodes.ColumnRef); ok {
		if cr.Fields != nil && len(cr.Fields.Items) == 1 {
			colName := stringVal(cr.Fields.Items[0])
			for _, te := range tlist {
				if !te.ResJunk && te.ResName == colName {
					return te, nil
				}
			}
		}
	}

	// SQL99: transform and match by expression.
	expr, err := ac.transformExpr(n)
	if err != nil {
		return nil, err
	}

	// Search for matching expression in existing target list.
	for _, te := range tlist {
		if exprEqual(te.Expr, expr) {
			return te, nil
		}
	}

	// No match — create a resjunk target entry.
	// pg: makeTargetEntry with resjunk=true, appended to tlist
	resNo := int16(len(ac.query.TargetList) + 1)
	name := figureColName(n)
	te := &TargetEntry{
		Expr:    expr,
		ResNo:   resNo,
		ResName: name,
		ResJunk: true,
	}
	ac.query.TargetList = append(ac.query.TargetList, te)
	return te, nil
}

// assignSortGroupRef assigns a sort/group reference number to a target entry.
//
// pg: src/backend/parser/parse_clause.c — assignSortGroupRef
func (ac *analyzeCtx) assignSortGroupRef(tle *TargetEntry, tlist []*TargetEntry) uint32 {
	if tle.ResSortGroupRef != 0 {
		return tle.ResSortGroupRef
	}
	var maxRef uint32
	for _, te := range tlist {
		if te.ResSortGroupRef > maxRef {
			maxRef = te.ResSortGroupRef
		}
	}
	tle.ResSortGroupRef = maxRef + 1
	return tle.ResSortGroupRef
}

// exprEqual checks if two analyzed expressions are structurally equal.
// Used for matching ORDER BY/GROUP BY expressions to target list entries.
//
// (pgddl helper — PG uses equal() on stripped expression nodes)
func exprEqual(a, b AnalyzedExpr) bool {
	if a == nil || b == nil {
		return a == b
	}
	switch va := a.(type) {
	case *VarExpr:
		vb, ok := b.(*VarExpr)
		return ok && va.RangeIdx == vb.RangeIdx && va.AttNum == vb.AttNum
	case *ConstExpr:
		cb, ok := b.(*ConstExpr)
		return ok && va.TypeOID == cb.TypeOID && va.Value == cb.Value && va.IsNull == cb.IsNull
	case *FuncCallExpr:
		fb, ok := b.(*FuncCallExpr)
		if !ok || va.FuncOID != fb.FuncOID || len(va.Args) != len(fb.Args) {
			return false
		}
		for i := range va.Args {
			if !exprEqual(va.Args[i], fb.Args[i]) {
				return false
			}
		}
		return true
	case *OpExpr:
		ob, ok := b.(*OpExpr)
		if !ok || va.OpOID != ob.OpOID {
			return false
		}
		return exprEqual(va.Left, ob.Left) && exprEqual(va.Right, ob.Right)
	case *AggExpr:
		ab, ok := b.(*AggExpr)
		if !ok || va.AggFuncOID != ab.AggFuncOID || va.AggStar != ab.AggStar || len(va.Args) != len(ab.Args) {
			return false
		}
		for i := range va.Args {
			if !exprEqual(va.Args[i], ab.Args[i]) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

// resolveCollation picks a collation from a list of analyzed expressions.
// Simplified from PG's assign_collations_walker: first non-zero collation wins.
//
// (pgddl helper — PG uses select_common_collation / assign_collations_walker)
func resolveCollation(args ...AnalyzedExpr) uint32 {
	for _, arg := range args {
		if arg == nil {
			continue
		}
		if c := arg.exprCollation(); c != 0 {
			return c
		}
	}
	return 0
}

// AnalyzeSelectStmt transforms a raw SelectStmt AST node into a semantically
// analyzed Query. The catalog must already contain the schema (tables, types,
// functions) referenced by the query — load it via Exec() beforehand.
//
// The returned Query contains resolved column references (VarExpr), type
// information, and provenance tracking (TargetEntry.ResOrigTbl/ResOrigCol)
// that can be used to extract column lineage.
//
// pg: src/backend/parser/analyze.c — transformSelectStmt
func (c *Catalog) AnalyzeSelectStmt(stmt *nodes.SelectStmt) (*Query, error) {
	return c.analyzeSelectStmt(stmt)
}

// AnalyzeExprNoContext analyzes a raw expression node without any range table
// context. Used for function parameter default expressions.
//
// pg: src/backend/parser/parse_expr.c — transformExpr (no pstate relations)
func (c *Catalog) AnalyzeExprNoContext(expr nodes.Node) (AnalyzedExpr, error) {
	if expr == nil {
		return nil, nil
	}
	ac := &analyzeCtx{
		catalog: c,
		query:   &Query{},
	}
	return ac.transformExpr(expr)
}

// AnalyzeTriggerWhenExpr analyzes a trigger WHEN clause in the context of
// OLD (varno=1) and NEW (varno=2) pseudo-relations.
//
// pg: src/backend/commands/trigger.c — CreateTrigger (WHEN clause setup)
func (c *Catalog) AnalyzeTriggerWhenExpr(expr nodes.Node, rel *Relation) (AnalyzedExpr, error) {
	if expr == nil {
		return nil, nil
	}
	oldRTE := c.buildRelationRTE(rel)
	oldRTE.ERef = "old"
	newRTE := c.buildRelationRTE(rel)
	newRTE.ERef = "new"
	ac := &analyzeCtx{
		catalog: c,
		query:   &Query{RangeTable: []*RangeTableEntry{oldRTE, newRTE}},
	}
	return ac.transformExpr(expr)
}

// AnalyzeStandaloneExpr analyzes a raw expression node in the context of a
// relation's columns. Used for CHECK constraints, column defaults, and policy
// expressions.
//
// pg: src/backend/parser/parse_expr.c — transformExpr (standalone context)
func (c *Catalog) AnalyzeStandaloneExpr(expr nodes.Node, rel *Relation) (AnalyzedExpr, error) {
	if expr == nil {
		return nil, nil
	}
	rte := c.buildRelationRTE(rel)
	ac := &analyzeCtx{
		catalog: c,
		query:   &Query{RangeTable: []*RangeTableEntry{rte}},
	}
	return ac.transformExpr(expr)
}

// AnalyzeDomainExpr analyzes a raw expression node in the context of a domain
// type. The "VALUE" keyword is intercepted in transformColumnRef and replaced
// with CoerceToDomainValueExpr.
//
// pg: src/backend/commands/typecmds.c — domainAddCheckConstraint
func (c *Catalog) AnalyzeDomainExpr(expr nodes.Node, domainBaseTypeOID uint32, domainBaseTypMod int32) (AnalyzedExpr, error) {
	if expr == nil {
		return nil, nil
	}
	bt := c.typeByOID[domainBaseTypeOID]
	var coll uint32
	if bt != nil {
		coll = bt.Collation
	}
	ac := &analyzeCtx{
		catalog:             c,
		query:               &Query{},
		domainConstraint:    true,
		domainBaseTypeOID:   domainBaseTypeOID,
		domainBaseTypMod:    domainBaseTypMod,
		domainBaseCollation: coll,
	}
	return ac.transformExpr(expr)
}

// buildRelationRTE builds a RangeTableEntry from a Relation's columns.
//
// (pgddl helper — PG uses addRangeTableEntryForRelation)
func (c *Catalog) buildRelationRTE(rel *Relation) *RangeTableEntry {
	rte := &RangeTableEntry{
		Kind:    RTERelation,
		RelOID:  rel.OID,
		RelName: rel.Name,
		ERef:    rel.Name,
	}
	if rel.Schema != nil {
		rte.SchemaName = rel.Schema.Name
	}
	for _, col := range rel.Columns {
		rte.ColNames = append(rte.ColNames, col.Name)
		rte.ColTypes = append(rte.ColTypes, col.TypeOID)
		rte.ColTypMods = append(rte.ColTypMods, col.TypeMod)
		rte.ColCollations = append(rte.ColCollations, col.Collation)
	}
	return rte
}

// transformGroupingSet transforms a GroupingSet or plain column reference into a GroupingSetQ.
//
// pg: src/backend/parser/parse_clause.c — transformGroupingSet
func (ac *analyzeCtx) transformGroupingSet(n nodes.Node, tlist []*TargetEntry) (*GroupingSetQ, error) {
	gs, ok := n.(*nodes.GroupingSet)
	if !ok {
		// Plain column reference — wrap as GroupingSetSimple.
		tle, err := ac.findTargetlistEntry(n, tlist)
		if err != nil {
			return nil, err
		}
		ref := ac.assignSortGroupRef(tle, tlist)
		return &GroupingSetQ{
			Kind:    GroupingSetSimple,
			Content: []*SortGroupClause{{TLESortGroupRef: ref}},
		}, nil
	}

	switch gs.Kind {
	case nodes.GROUPING_SET_ROLLUP:
		gsq := &GroupingSetQ{Kind: GroupingSetRollup}
		if gs.Content != nil {
			for _, item := range gs.Content.Items {
				tle, err := ac.findTargetlistEntry(item, tlist)
				if err != nil {
					return nil, err
				}
				ref := ac.assignSortGroupRef(tle, tlist)
				gsq.Content = append(gsq.Content, &SortGroupClause{TLESortGroupRef: ref})
			}
		}
		return gsq, nil
	case nodes.GROUPING_SET_CUBE:
		gsq := &GroupingSetQ{Kind: GroupingSetCube}
		if gs.Content != nil {
			for _, item := range gs.Content.Items {
				tle, err := ac.findTargetlistEntry(item, tlist)
				if err != nil {
					return nil, err
				}
				ref := ac.assignSortGroupRef(tle, tlist)
				gsq.Content = append(gsq.Content, &SortGroupClause{TLESortGroupRef: ref})
			}
		}
		return gsq, nil
	case nodes.GROUPING_SET_SETS:
		gsq := &GroupingSetQ{Kind: GroupingSetSets}
		if gs.Content != nil {
			for _, item := range gs.Content.Items {
				sub, err := ac.transformGroupingSet(item, tlist)
				if err != nil {
					return nil, err
				}
				gsq.Sets = append(gsq.Sets, sub)
			}
		}
		return gsq, nil
	case nodes.GROUPING_SET_EMPTY:
		return &GroupingSetQ{Kind: GroupingSetSimple}, nil
	default:
		// GROUPING_SET_SIMPLE — plain column
		gsq := &GroupingSetQ{Kind: GroupingSetSimple}
		if gs.Content != nil {
			for _, item := range gs.Content.Items {
				tle, err := ac.findTargetlistEntry(item, tlist)
				if err != nil {
					return nil, err
				}
				ref := ac.assignSortGroupRef(tle, tlist)
				gsq.Content = append(gsq.Content, &SortGroupClause{TLESortGroupRef: ref})
			}
		}
		return gsq, nil
	}
}

// flattenGroupingSetRefs adds all column refs from a GroupingSetQ to the
// query's GroupClause for target list resolution.
func (ac *analyzeCtx) flattenGroupingSetRefs(gs *GroupingSetQ, q *Query) {
	for _, sgc := range gs.Content {
		isDup := false
		for _, existing := range q.GroupClause {
			if existing.TLESortGroupRef == sgc.TLESortGroupRef {
				isDup = true
				break
			}
		}
		if !isDup {
			q.GroupClause = append(q.GroupClause, sgc)
		}
	}
	for _, sub := range gs.Sets {
		ac.flattenGroupingSetRefs(sub, q)
	}
}

// transformGroupingFunc transforms a GROUPING(...) expression.
//
// pg: src/backend/parser/parse_expr.c — transformGroupingFunc
func (ac *analyzeCtx) transformGroupingFunc(gf *nodes.GroupingFunc) (AnalyzedExpr, error) {
	var args []AnalyzedExpr
	var refs []int
	if gf.Args != nil {
		for _, a := range gf.Args.Items {
			expr, err := ac.transformExpr(a)
			if err != nil {
				return nil, err
			}
			args = append(args, expr)
		}
	}
	if gf.Refs != nil {
		for _, r := range gf.Refs.Items {
			if iv, ok := r.(*nodes.Integer); ok {
				refs = append(refs, int(iv.Ival))
			}
		}
	}
	return &GroupingFuncExpr{
		Args: args,
		Refs: refs,
	}, nil
}

// transformXmlExpr transforms a XmlExpr node into an analyzed XmlExprQ.
//
// pg: src/backend/parser/parse_expr.c — transformXmlExpr
func (ac *analyzeCtx) transformXmlExpr(xe *nodes.XmlExpr) (AnalyzedExpr, error) {
	// Translate AST XmlOptionType (0=DOCUMENT,1=CONTENT) to catalog constants.
	xmlopt := XmlOptionContent
	if xe.Xmloption == nodes.XMLOPTION_DOCUMENT {
		xmlopt = XmlOptionDocument
	}

	result := &XmlExprQ{
		Op:        int(xe.Op),
		Name:      xe.Name,
		Xmloption: xmlopt,
		TypeOID:   XMLOID,
		TypeMod:   -1,
	}

	// IS DOCUMENT returns boolean.
	if xe.Op == nodes.IS_XMLSERIALIZE {
		if xe.Type != 0 {
			result.TypeOID = uint32(xe.Type)
		} else {
			result.TypeOID = TEXTOID
		}
		result.TypeMod = xe.Typmod
	}

	// Named args (xml_attributes / XMLFOREST elements).
	// Parser stores these as ResTarget nodes; extract the expression and
	// the AS alias (ResTarget.Name) so ArgNames stays in sync.
	if xe.NamedArgs != nil {
		for _, a := range xe.NamedArgs.Items {
			if rt, ok := a.(*nodes.ResTarget); ok {
				expr, err := ac.transformExpr(rt.Val)
				if err != nil {
					return nil, err
				}
				result.NamedArgs = append(result.NamedArgs, expr)
				result.ArgNames = append(result.ArgNames, rt.Name)
			} else {
				expr, err := ac.transformExpr(a)
				if err != nil {
					return nil, err
				}
				result.NamedArgs = append(result.NamedArgs, expr)
				result.ArgNames = append(result.ArgNames, "")
			}
		}
	}

	// Arg names (populated from AST when parser fills ArgNames directly).
	if xe.ArgNames != nil {
		for _, a := range xe.ArgNames.Items {
			if s, ok := a.(*nodes.String); ok {
				result.ArgNames = append(result.ArgNames, s.Str)
			}
		}
	}

	// Args.
	if xe.Args != nil {
		for _, a := range xe.Args.Items {
			expr, err := ac.transformExpr(a)
			if err != nil {
				return nil, err
			}
			result.Args = append(result.Args, expr)
		}
	}

	return result, nil
}

// transformXmlSerialize transforms an XMLSERIALIZE expression into an XmlExprQ.
//
// pg: src/backend/parser/parse_expr.c — transformXmlSerialize
func (ac *analyzeCtx) transformXmlSerialize(xs *nodes.XmlSerialize) (AnalyzedExpr, error) {
	// Transform the inner expression.
	arg, err := ac.transformExpr(xs.Expr)
	if err != nil {
		return nil, err
	}

	// Resolve target type.
	var typeOID uint32 = TEXTOID
	var typMod int32 = -1
	if xs.TypeName != nil {
		typeOID, typMod, err = ac.catalog.resolveTypeName(xs.TypeName)
		if err != nil {
			return nil, err
		}
	}

	// Translate AST XmlOptionType to catalog constants.
	serializeOpt := XmlOptionContent
	if xs.Xmloption == nodes.XMLOPTION_DOCUMENT {
		serializeOpt = XmlOptionDocument
	}

	return &XmlExprQ{
		Op:        XmlOpSerialize,
		Xmloption: serializeOpt,
		Args:      []AnalyzedExpr{arg},
		TypeOID:   typeOID,
		TypeMod:   typMod,
	}, nil
}

// transformRangeTableSampleItem processes a RangeTableSample FROM clause item.
// It first transforms the relation, then attaches the tablesample clause.
//
// pg: src/backend/parser/parse_clause.c — transformRangeTableSample
func (ac *analyzeCtx) transformRangeTableSampleItem(rts *nodes.RangeTableSample) (JoinNode, error) {
	// Transform the inner relation.
	jn, err := ac.transformFromClauseItem(rts.Relation)
	if err != nil {
		return nil, err
	}
	// The inner item should be a RangeTableRef.
	rtRef, ok := jn.(*RangeTableRef)
	if !ok {
		return nil, fmt.Errorf("TABLESAMPLE on non-table item")
	}
	rte := ac.query.RangeTable[rtRef.RTIndex]
	if err := ac.transformRangeTableSample(rts, rte); err != nil {
		return nil, err
	}
	return jn, nil
}

// transformRangeTableSample processes a TABLESAMPLE clause and attaches it to the RTE.
//
// pg: src/backend/parser/parse_clause.c — transformRangeTableSample
func (ac *analyzeCtx) transformRangeTableSample(rts *nodes.RangeTableSample, rte *RangeTableEntry) error {
	ts := &TablesampleClauseQ{}

	// Method name.
	if rts.Method != nil {
		for _, item := range rts.Method.Items {
			if s, ok := item.(*nodes.String); ok {
				ts.Method = s.Str
			}
		}
	}

	// Args.
	if rts.Args != nil {
		for _, a := range rts.Args.Items {
			expr, err := ac.transformExpr(a)
			if err != nil {
				return err
			}
			ts.Args = append(ts.Args, expr)
		}
	}

	// REPEATABLE.
	if rts.Repeatable != nil {
		expr, err := ac.transformExpr(rts.Repeatable)
		if err != nil {
			return err
		}
		ts.Repeatable = expr
	}

	rte.Tablesample = ts
	return nil
}

// setOpKeyword returns the SQL keyword for a set operation.
func setOpKeyword(op nodes.SetOperation) string {
	switch op {
	case nodes.SETOP_UNION:
		return "UNION"
	case nodes.SETOP_INTERSECT:
		return "INTERSECT"
	case nodes.SETOP_EXCEPT:
		return "EXCEPT"
	default:
		return "SET"
	}
}

// --- JSON Expression Transformers (PG16+) ---

// transformJsonIsPredicate transforms expr IS [NOT] JSON [OBJECT|ARRAY|SCALAR].
//
// pg: src/backend/parser/parse_expr.c — transformJsonIsPredicate
func (ac *analyzeCtx) transformJsonIsPredicate(n *nodes.JsonIsPredicate) (AnalyzedExpr, error) {
	expr, err := ac.transformExpr(n.Expr)
	if err != nil {
		return nil, err
	}
	return &JsonIsPredicateExpr{
		Expr:       expr,
		ItemType:   int(n.ItemType),
		UniqueKeys: n.UniqueKeys,
		IsNot:      false, // PG wraps NOT externally via BoolExpr
	}, nil
}

// transformJsonFuncExpr transforms JSON_VALUE, JSON_QUERY, JSON_EXISTS.
//
// pg: src/backend/parser/parse_expr.c — transformJsonFuncExpr
func (ac *analyzeCtx) transformJsonFuncExpr(n *nodes.JsonFuncExpr) (AnalyzedExpr, error) {
	var expr AnalyzedExpr
	if n.ContextItem != nil {
		var err error
		expr, err = ac.transformExpr(n.ContextItem.RawExpr)
		if err != nil {
			return nil, err
		}
	}

	var pathStr string
	if n.Pathspec != nil {
		if s, ok := n.Pathspec.(*nodes.A_Const); ok && s.Val != nil {
			if sv, ok := s.Val.(*nodes.String); ok {
				pathStr = sv.Str
			}
		}
		// Also try JsonTablePathSpec
		if ps, ok := n.Pathspec.(*nodes.JsonTablePathSpec); ok && ps.String != nil {
			if s, ok := ps.String.(*nodes.A_Const); ok && s.Val != nil {
				if sv, ok := s.Val.(*nodes.String); ok {
					pathStr = sv.Str
				}
			}
		}
	}

	var op JsonExprOpQ
	var resultType uint32
	switch n.Op {
	case nodes.JSON_EXISTS_OP:
		op = JsonExprExists
		resultType = BOOLOID
	case nodes.JSON_QUERY_OP:
		op = JsonExprQuery
		resultType = JSONBOID
	case nodes.JSON_VALUE_OP:
		op = JsonExprValue
		resultType = TEXTOID
	default:
		resultType = TEXTOID
	}

	return &JsonExprQ{
		Op:         op,
		Expr:       expr,
		Path:       pathStr,
		ResultType: resultType,
	}, nil
}

// transformJsonObjectConstructor transforms JSON_OBJECT(...).
func (ac *analyzeCtx) transformJsonObjectConstructor(n *nodes.JsonObjectConstructor) (AnalyzedExpr, error) {
	var args []AnalyzedExpr
	if n.Exprs != nil {
		for _, item := range n.Exprs.Items {
			if kv, ok := item.(*nodes.JsonKeyValue); ok {
				// Transform key
				keyExpr, err := ac.transformExpr(kv.Key)
				if err != nil {
					return nil, err
				}
				args = append(args, keyExpr)
				// Transform value
				if kv.Value != nil {
					valExpr, err := ac.transformExpr(kv.Value.RawExpr)
					if err != nil {
						return nil, err
					}
					args = append(args, valExpr)
				}
			}
		}
	}
	return &JsonConstructorExprQ{
		ConstructorType: JsonConstructorObject,
		Args:            args,
		ResultType:      JSONOID,
	}, nil
}

// transformJsonArrayConstructor transforms JSON_ARRAY(...).
func (ac *analyzeCtx) transformJsonArrayConstructor(n *nodes.JsonArrayConstructor) (AnalyzedExpr, error) {
	var args []AnalyzedExpr
	if n.Exprs != nil {
		for _, item := range n.Exprs.Items {
			if jve, ok := item.(*nodes.JsonValueExpr); ok {
				expr, err := ac.transformExpr(jve.RawExpr)
				if err != nil {
					return nil, err
				}
				args = append(args, expr)
			} else {
				expr, err := ac.transformExpr(item)
				if err != nil {
					return nil, err
				}
				args = append(args, expr)
			}
		}
	}
	return &JsonConstructorExprQ{
		ConstructorType: JsonConstructorArray,
		Args:            args,
		ResultType:      JSONOID,
	}, nil
}

// transformJsonValueExpr transforms a JsonValueExpr (inner expression).
func (ac *analyzeCtx) transformJsonValueExpr(n *nodes.JsonValueExpr) (AnalyzedExpr, error) {
	expr, err := ac.transformExpr(n.RawExpr)
	if err != nil {
		return nil, err
	}
	return &JsonValueExprQ{
		Expr:       expr,
		ResultType: expr.exprType(),
	}, nil
}

// transformJsonScalarExpr transforms JSON_SCALAR().
func (ac *analyzeCtx) transformJsonScalarExpr(n *nodes.JsonScalarExpr) (AnalyzedExpr, error) {
	expr, err := ac.transformExpr(n.Expr)
	if err != nil {
		return nil, err
	}
	return &JsonConstructorExprQ{
		ConstructorType: JsonConstructorScalar,
		Args:            []AnalyzedExpr{expr},
		ResultType:      JSONOID,
	}, nil
}

// transformJsonParseExpr transforms JSON() parse expression.
func (ac *analyzeCtx) transformJsonParseExpr(n *nodes.JsonParseExpr) (AnalyzedExpr, error) {
	var inner AnalyzedExpr
	if n.Expr != nil {
		var err error
		inner, err = ac.transformExpr(n.Expr.RawExpr)
		if err != nil {
			return nil, err
		}
	}
	return &JsonConstructorExprQ{
		ConstructorType: JsonConstructorParse,
		Args:            []AnalyzedExpr{inner},
		ResultType:      JSONOID,
	}, nil
}

// transformJsonSerializeExpr transforms JSON_SERIALIZE().
func (ac *analyzeCtx) transformJsonSerializeExpr(n *nodes.JsonSerializeExpr) (AnalyzedExpr, error) {
	var inner AnalyzedExpr
	if n.Expr != nil {
		var err error
		inner, err = ac.transformExpr(n.Expr.RawExpr)
		if err != nil {
			return nil, err
		}
	}
	return &JsonConstructorExprQ{
		ConstructorType: JsonConstructorSerialize,
		Args:            []AnalyzedExpr{inner},
		ResultType:      TEXTOID,
	}, nil
}
