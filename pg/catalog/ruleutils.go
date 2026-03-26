package catalog

import (
	"fmt"
	"strings"
)

// GetViewDefinition returns the SQL definition of a view,
// matching pg_get_viewdef() output format.
//
// When pretty is false, matches pg_get_viewdef(oid) (PRETTYFLAG_INDENT only).
// When pretty is true, matches pg_get_viewdef(oid, true) (PRETTYFLAG_INDENT | PRETTYFLAG_PAREN).
//
// pg: src/backend/utils/adt/ruleutils.c — pg_get_viewdef_worker
func (c *Catalog) GetViewDefinition(schema, name string, pretty ...bool) (string, error) {
	rel := c.GetRelation(schema, name)
	if rel == nil {
		return "", errUndefinedTable(name)
	}
	if rel.RelKind != 'v' {
		return "", errWrongObjectType(name, "a view")
	}
	p := false
	if len(pretty) > 0 {
		p = pretty[0]
	}
	return c.deparseRelationQuery(rel, p)
}

// GetMatViewDefinition returns the SQL definition of a materialized view,
// matching pg_get_viewdef() output format (PG's pg_get_viewdef works on matviews too).
//
// When pretty is false, matches pg_get_viewdef(oid) (PRETTYFLAG_INDENT only).
// When pretty is true, matches pg_get_viewdef(oid, true) (PRETTYFLAG_INDENT | PRETTYFLAG_PAREN).
//
// pg: src/backend/utils/adt/ruleutils.c — pg_get_viewdef_worker
func (c *Catalog) GetMatViewDefinition(schema, name string, pretty ...bool) (string, error) {
	rel := c.GetRelation(schema, name)
	if rel == nil {
		return "", errUndefinedTable(name)
	}
	if rel.RelKind != 'm' {
		return "", errWrongObjectType(name, "a materialized view")
	}
	p := false
	if len(pretty) > 0 {
		p = pretty[0]
	}
	return c.deparseRelationQuery(rel, p)
}

// deparseRelationQuery deparses the AnalyzedQuery of a view or matview.
//
// pretty=false: PRETTYFLAG_INDENT only (matches pg_get_viewdef(oid))
// pretty=true: PRETTYFLAG_INDENT | PRETTYFLAG_PAREN (matches pg_get_viewdef(oid, true))
//
// (pgddl helper — shared implementation for GetViewDefinition/GetMatViewDefinition)
func (c *Catalog) deparseRelationQuery(rel *Relation, pretty bool) (string, error) {
	if rel.AnalyzedQuery == nil {
		return "", nil
	}

	ctx := &deparseCtx{
		catalog:     c,
		buf:         &strings.Builder{},
		indentLevel: 0,
		// pg: src/backend/utils/adt/ruleutils.c — pg_get_viewdef_worker
		prettyIndent:       true,
		prettyParen:        pretty,
		hasParentNamespace: false,
		wrapColumn:         0, // 0 = always wrap (PG's WRAP_COLUMN_DEFAULT for pg_get_viewdef)
		colNamesVisible:    true,
	}

	ctx.getQueryDef(rel.AnalyzedQuery)
	// pg: src/backend/utils/adt/ruleutils.c — make_viewdef (line 5413)
	// PG always appends ';' after the view definition.
	ctx.buf.WriteByte(';')
	return ctx.buf.String(), nil
}

// DeparseExpr deparses an analyzed expression to SQL text,
// matching pg_get_expr() / pg_get_constraintdef() output format.
//
// pg: src/backend/utils/adt/ruleutils.c — deparse_expression_pretty
func (c *Catalog) DeparseExpr(expr AnalyzedExpr, rtable []*RangeTableEntry, prettyParen bool) string {
	if expr == nil {
		return ""
	}
	ctx := &deparseCtx{
		catalog:      c,
		buf:          &strings.Builder{},
		prettyIndent: true,
		prettyParen:  prettyParen,
		query:        &Query{RangeTable: rtable},
	}
	ctx.getRuleExpr(expr, false)
	return ctx.buf.String()
}

// deparseCtx holds state during query deparsing.
//
// pg: src/backend/utils/adt/ruleutils.c — deparse_context
type deparseCtx struct {
	catalog            *Catalog
	buf                *strings.Builder
	indentLevel        int
	prettyIndent       bool
	prettyParen        bool
	hasParentNamespace bool // true when inside a subquery/setop (affects varprefix)
	wrapColumn         int  // line-wrap threshold: 0=always wrap, -1=never wrap
	colNamesVisible    bool // true for views (always show AS for non-Var)
	query              *Query
	parentCtx          *deparseCtx // for resolving outer Var references (LevelsUp > 0)
}

const (
	prettyIndentStd  = 8
	prettyIndentJoin = 4
	prettyIndentVar  = 4
)

// getQueryDef is the top-level query deparser.
//
// pg: src/backend/utils/adt/ruleutils.c — get_query_def
func (d *deparseCtx) getQueryDef(query *Query) {
	d.query = query
	d.getSelectQueryDef(query)
}

// getSelectQueryDef deparses a SELECT query.
//
// pg: src/backend/utils/adt/ruleutils.c — get_select_query_def
func (d *deparseCtx) getSelectQueryDef(query *Query) {
	if query.SetOp != SetOpNone {
		d.getSetOpQuery(query)
		return
	}
	d.getBasicSelectQuery(query)
}

// getBasicSelectQuery deparses a simple SELECT ... FROM ... WHERE ... query.
//
// pg: src/backend/utils/adt/ruleutils.c — get_basic_select_query
func (d *deparseCtx) getBasicSelectQuery(query *Query) {
	savedQuery := d.query
	d.query = query

	// WITH clause (CTEs).
	if len(query.CTEList) > 0 {
		d.getWithClause(query)
	}

	if d.prettyIndent {
		d.indentLevel += prettyIndentStd
		d.buf.WriteByte(' ')
	}

	d.buf.WriteString("SELECT")

	if query.Distinct {
		if len(query.DistinctOn) > 0 {
			// DISTINCT ON (col1, col2, ...)
			d.buf.WriteString(" DISTINCT ON (")
			for i, dc := range query.DistinctOn {
				if i > 0 {
					d.buf.WriteString(", ")
				}
				d.getSortGroupClause(dc, query.TargetList, false)
			}
			d.buf.WriteByte(')')
		} else {
			d.buf.WriteString(" DISTINCT")
		}
	}

	// Target list.
	d.getTargetList(query.TargetList)

	// FROM clause.
	if query.JoinTree != nil && len(query.JoinTree.FromList) > 0 {
		d.getFromClause(query)
	}

	// WHERE clause.
	if query.JoinTree != nil && query.JoinTree.Quals != nil {
		d.appendContextKeyword(" WHERE ", -prettyIndentStd, prettyIndentStd, 1)
		d.getRuleExpr(query.JoinTree.Quals, false)
	}

	// GROUP BY.
	if len(query.GroupingSets) > 0 {
		d.appendContextKeyword(" GROUP BY ", -prettyIndentStd, prettyIndentStd, 1)
		for i, gs := range query.GroupingSets {
			if i > 0 {
				d.buf.WriteString(", ")
			}
			d.getGroupingSet(gs, query.TargetList)
		}
	} else if len(query.GroupClause) > 0 {
		d.appendContextKeyword(" GROUP BY ", -prettyIndentStd, prettyIndentStd, 1)
		for i, grp := range query.GroupClause {
			if i > 0 {
				d.buf.WriteString(", ")
			}
			d.getSortGroupClause(grp, query.TargetList, false)
		}
	}

	// HAVING.
	if query.HavingQual != nil {
		d.appendContextKeyword(" HAVING ", -prettyIndentStd, prettyIndentStd, 0)
		d.getRuleExpr(query.HavingQual, false)
	}

	// ORDER BY.
	//
	// pg: src/backend/utils/adt/ruleutils.c — get_rule_orderby
	if len(query.SortClause) > 0 {
		d.appendContextKeyword(" ORDER BY ", -prettyIndentStd, prettyIndentStd, 1)
		for i, srt := range query.SortClause {
			if i > 0 {
				d.buf.WriteString(", ")
			}
			d.getSortGroupClause(srt, query.TargetList, false)
			// Emit direction and nulls.
			if srt.Descending {
				d.buf.WriteString(" DESC")
				// DESC defaults to NULLS FIRST, so only show NULLS LAST.
				if !srt.NullsFirst {
					d.buf.WriteString(" NULLS LAST")
				}
			} else {
				// ASC is default, so only show NULLS FIRST.
				if srt.NullsFirst {
					d.buf.WriteString(" NULLS FIRST")
				}
			}
		}
	}

	// OFFSET (PG outputs OFFSET before LIMIT).
	if query.LimitOffset != nil {
		d.appendContextKeyword(" OFFSET ", -prettyIndentStd, prettyIndentStd, 0)
		d.getRuleExpr(query.LimitOffset, false)
	}

	// LIMIT.
	if query.LimitCount != nil {
		d.appendContextKeyword(" LIMIT ", -prettyIndentStd, prettyIndentStd, 0)
		d.getRuleExpr(query.LimitCount, false)
	}

	if d.prettyIndent {
		d.indentLevel -= prettyIndentStd
	}

	d.query = savedQuery
}

// getSetOpQuery deparses a UNION/INTERSECT/EXCEPT query.
//
// pg: src/backend/utils/adt/ruleutils.c — get_setop_query
func (d *deparseCtx) getSetOpQuery(query *Query) {
	savedQuery := d.query
	d.query = query

	// Left branch.
	// pg: get_setop_query — only wrap LHS in parens when the left child is a
	// different set operation (e.g., UNION inside EXCEPT). Same op+all left-associates naturally.
	needParens := false
	if query.LArg != nil && query.LArg.SetOp != SetOpNone {
		if query.SetOp != query.LArg.SetOp || query.AllSetOp != query.LArg.AllSetOp {
			needParens = true
		}
	}

	subindent := 0
	if needParens {
		d.buf.WriteByte('(')
		subindent = prettyIndentStd
		d.appendContextKeyword("", subindent, 0, 0)
	}
	if query.LArg != nil {
		subCtx := d.subContext()
		subCtx.hasParentNamespace = true
		subCtx.colNamesVisible = d.colNamesVisible
		subCtx.getSelectQueryDef(query.LArg)
		d.indentLevel = subCtx.indentLevel
	}

	if needParens {
		d.appendContextKeyword(") ", -subindent, 0, 0)
	} else if d.prettyIndent {
		d.appendContextKeyword("", -subindent, 0, 0)
	} else {
		d.buf.WriteByte(' ')
	}

	// Set operation keyword.
	switch query.SetOp {
	case SetOpUnion:
		d.buf.WriteString("UNION ")
	case SetOpIntersect:
		d.buf.WriteString("INTERSECT ")
	case SetOpExcept:
		d.buf.WriteString("EXCEPT ")
	}
	if query.AllSetOp {
		d.buf.WriteString("ALL ")
	}

	// Right branch.
	// pg: Always parenthesize if RHS is another setop.
	needParens = query.RArg != nil && query.RArg.SetOp != SetOpNone

	subindent = 0
	if needParens {
		d.buf.WriteByte('(')
		subindent = prettyIndentStd
	}

	// Newline before right branch.
	d.appendContextKeyword("", subindent, 0, 0)

	if query.RArg != nil {
		subCtx := d.subContext()
		subCtx.hasParentNamespace = true
		subCtx.colNamesVisible = false // PG sets colNamesVisible=false for right branch
		subCtx.getSelectQueryDef(query.RArg)
		d.indentLevel = subCtx.indentLevel
	}

	if d.prettyIndent {
		d.indentLevel -= subindent
	}
	if needParens {
		d.appendContextKeyword(")", 0, 0, 0)
	}

	d.query = savedQuery
}

// getTargetList deparses the SELECT target list.
//
// pg: src/backend/utils/adt/ruleutils.c — get_target_list
func (d *deparseCtx) getTargetList(tlist []*TargetEntry) {
	sep := " "
	colno := 0
	for _, te := range tlist {
		if te.ResJunk {
			continue
		}
		colno++

		d.buf.WriteString(sep)
		sep = ", "

		// Render expression to a temporary buffer so we can check for wrapping.
		var targetBuf strings.Builder
		tmpCtx := d.subContextBuf(&targetBuf)

		// pg: get_target_list — Var uses get_variable directly
		var attname string
		if v, ok := te.Expr.(*VarExpr); ok {
			attname = tmpCtx.getVariableReturningName(v)
		} else {
			tmpCtx.getRuleExpr(te.Expr, true)
			if d.colNamesVisible {
				attname = "" // force AS for non-Var when colNamesVisible
			} else {
				attname = "?column?"
			}
		}

		// Add AS label if needed.
		// pg: show AS when attname != colname
		colname := te.ResName
		if colname != "" && attname != colname {
			targetBuf.WriteString(" AS ")
			targetBuf.WriteString(quoteIdentifier(colname))
		}

		// Propagate indent level changes.
		d.indentLevel = tmpCtx.indentLevel

		target := targetBuf.String()

		// Check for line-wrapping.
		if d.prettyIndent && d.wrapColumn >= 0 {
			leadingNL := len(target) > 0 && target[0] == '\n'

			if leadingNL {
				// Expression starts with newline (e.g., CASE) — remove trailing spaces.
				trimBufTrailingSpaces(d.buf)
			} else if colno > 1 {
				// Wrap to new line.
				d.appendContextKeyword("", -prettyIndentStd, prettyIndentStd, prettyIndentVar)
			}
		}

		d.buf.WriteString(target)
	}
}

// getFromClause deparses the FROM clause.
//
// pg: src/backend/utils/adt/ruleutils.c — get_from_clause
func (d *deparseCtx) getFromClause(query *Query) {
	if query.JoinTree == nil || len(query.JoinTree.FromList) == 0 {
		return
	}

	first := true
	for _, jn := range query.JoinTree.FromList {
		if first {
			d.appendContextKeyword(" FROM ", -prettyIndentStd, prettyIndentStd, 2)
			first = false
			d.getFromClauseItem(jn, query)
		} else {
			d.buf.WriteString(", ")

			// Render to temp buffer for wrapping check.
			var itemBuf strings.Builder
			tmpCtx := d.subContextBuf(&itemBuf)
			tmpCtx.getFromClauseItem(jn, query)
			d.indentLevel = tmpCtx.indentLevel

			item := itemBuf.String()

			if d.prettyIndent && d.wrapColumn >= 0 {
				if len(item) > 0 && item[0] == '\n' {
					trimBufTrailingSpaces(d.buf)
				} else {
					d.appendContextKeyword("", -prettyIndentStd, prettyIndentStd, prettyIndentVar)
				}
			}

			d.buf.WriteString(item)
		}
	}
}

// getFromClauseItem deparses a single FROM clause item.
//
// pg: src/backend/utils/adt/ruleutils.c — get_from_clause_item
func (d *deparseCtx) getFromClauseItem(jn JoinNode, query *Query) {
	switch v := jn.(type) {
	case *RangeTableRef:
		rte := query.RangeTable[v.RTIndex]

		// LATERAL prefix.
		if rte.Lateral {
			d.buf.WriteString("LATERAL ")
		}

		switch rte.Kind {
		case RTERelation:
			d.buf.WriteString(d.generateRelationName(rte))
			if rte.Alias != "" {
				d.buf.WriteByte(' ')
				d.buf.WriteString(quoteIdentifier(rte.Alias))
			}
			if rte.Tablesample != nil {
				d.getTablesample(rte.Tablesample)
			}
		case RTESubquery:
			d.buf.WriteByte('(')
			subCtx := d.subContext()
			subCtx.hasParentNamespace = true
			subCtx.getQueryDef(rte.Subquery)
			d.indentLevel = subCtx.indentLevel
			d.buf.WriteByte(')')
			if rte.Alias != "" {
				d.buf.WriteByte(' ')
				d.buf.WriteString(quoteIdentifier(rte.Alias))
			}
		case RTECTE:
			d.buf.WriteString(quoteIdentifier(rte.CTEName))
			if rte.Alias != "" && rte.Alias != rte.CTEName {
				d.buf.WriteByte(' ')
				d.buf.WriteString(quoteIdentifier(rte.Alias))
			}
		case RTEFunction:
			// pg: src/backend/utils/adt/ruleutils.c — get_from_clause_item RTE_FUNCTION case
			if len(rte.FuncExprs) == 1 {
				d.getRuleExpr(rte.FuncExprs[0], false)
			} else {
				d.buf.WriteString("ROWS FROM(")
				for i, fe := range rte.FuncExprs {
					if i > 0 {
						d.buf.WriteString(", ")
					}
					d.getRuleExpr(fe, false)
				}
				d.buf.WriteByte(')')
			}
			if rte.Ordinality {
				d.buf.WriteString(" WITH ORDINALITY")
			}
			// pg: src/backend/utils/adt/ruleutils.c — get_rte_alias + get_column_alias_list
			// For function RTEs, PG always prints alias and column aliases
			// (printaliases=true) to protect against column name instability.
			refname := rte.Alias
			if refname == "" {
				refname = rte.ERef
			}
			if refname != "" {
				d.buf.WriteByte(' ')
				d.buf.WriteString(quoteIdentifier(refname))
				// pg: get_column_alias_list — always print for function RTEs
				if len(rte.ColNames) > 0 {
					d.buf.WriteByte('(')
					for i, cn := range rte.ColNames {
						if i > 0 {
							d.buf.WriteString(", ")
						}
						d.buf.WriteString(quoteIdentifier(cn))
					}
					d.buf.WriteByte(')')
				}
			}
		}

	case *JoinExprNode:
		// pg: get_from_clause_item wraps joins in parens when !PRETTY_PAREN
		if !d.prettyParen {
			d.buf.WriteByte('(')
		}

		d.getFromClauseItem(v.Left, query)
		switch v.JoinType {
		case JoinInner:
			if v.Quals != nil {
				d.appendContextKeyword(" JOIN ", -prettyIndentStd, prettyIndentStd, prettyIndentJoin)
			} else {
				d.appendContextKeyword(" CROSS JOIN ", -prettyIndentStd, prettyIndentStd, prettyIndentJoin)
			}
		case JoinLeft:
			d.appendContextKeyword(" LEFT JOIN ", -prettyIndentStd, prettyIndentStd, prettyIndentJoin)
		case JoinFull:
			d.appendContextKeyword(" FULL JOIN ", -prettyIndentStd, prettyIndentStd, prettyIndentJoin)
		case JoinRight:
			d.appendContextKeyword(" RIGHT JOIN ", -prettyIndentStd, prettyIndentStd, prettyIndentJoin)
		case JoinCross:
			d.appendContextKeyword(" CROSS JOIN ", -prettyIndentStd, prettyIndentStd, prettyIndentJoin)
		}
		d.getFromClauseItem(v.Right, query)

		// pg: src/backend/utils/adt/ruleutils.c — get_from_clause_item (USING vs ON)
		if len(v.UsingClause) > 0 {
			d.buf.WriteString(" USING (")
			for i, col := range v.UsingClause {
				if i > 0 {
					d.buf.WriteString(", ")
				}
				d.buf.WriteString(quoteIdentifier(col))
			}
			d.buf.WriteByte(')')
		} else if v.Quals != nil {
			d.buf.WriteString(" ON ")
			if !d.prettyParen {
				d.buf.WriteByte('(')
			}
			d.getRuleExpr(v.Quals, false)
			if !d.prettyParen {
				d.buf.WriteByte(')')
			}
		}

		if !d.prettyParen {
			d.buf.WriteByte(')')
		}
	}
}

// getRuleExpr deparses an analyzed expression node.
//
// pg: src/backend/utils/adt/ruleutils.c — get_rule_expr
func (d *deparseCtx) getRuleExpr(expr AnalyzedExpr, showImplicit bool) {
	if expr == nil {
		return
	}

	switch v := expr.(type) {
	case *VarExpr:
		d.getVariable(v)
	case *ConstExpr:
		d.getConstExpr(v, 0)
	case *FuncCallExpr:
		d.getFuncExpr(v, showImplicit)
	case *AggExpr:
		d.getAggExpr(v)
	case *OpExpr:
		d.getOpExpr(v)
	case *RelabelExpr:
		d.getRelabelExpr(v, showImplicit)
	case *CoerceViaIOExpr:
		d.getCoerceViaIOExpr(v, showImplicit)
	case *CaseExprQ:
		d.getCaseExpr(v)
	case *CoalesceExprQ:
		d.getCoalesceExpr(v)
	case *BoolExprQ:
		d.getBoolExpr(v)
	case *NullTestExpr:
		d.getNullTest(v)
	case *SubLinkExpr:
		d.getSubLink(v)
	case *NullIfExprQ:
		d.getNullIf(v)
	case *MinMaxExprQ:
		d.getMinMax(v)
	case *BooleanTestExpr:
		d.getBooleanTest(v)
	case *SQLValueFuncExpr:
		d.getSQLValueFunc(v)
	case *DistinctExprQ:
		d.getDistinctExpr(v)
	case *ScalarArrayOpExpr:
		d.getScalarArrayOp(v)
	case *ArrayExprQ:
		d.getArrayExpr(v)
	case *RowExprQ:
		d.getRowExpr(v)
	case *CollateExprQ:
		d.getCollateExpr(v)
	case *FieldSelectExprQ:
		d.getFieldSelect(v)
	case *WindowFuncExpr:
		d.getWindowFunc(v)
	case *CoerceToDomainValueExpr:
		// pg: src/backend/utils/adt/ruleutils.c — get_domaincheck_expr
		d.buf.WriteString("VALUE")
	case *SubscriptingRefExpr:
		d.getSubscriptingRef(v)
	case *NamedArgExprQ:
		d.getNamedArgExpr(v)
	case *ArrayCoerceExprQ:
		d.getArrayCoerceExpr(v, showImplicit)
	case *CaseTestExprQ:
		// pg: CaseTestExpr is not emitted directly by ruleutils
		// It's an internal placeholder; nothing to output.
	case *CoerceToDomainExpr:
		// pg: CoerceToDomain — deparse the inner arg (domain coercion is implicit)
		d.getRuleExpr(v.Arg, showImplicit)
	case *RowCompareExprQ:
		d.getRowCompareExpr(v)
	case *ParamExpr:
		fmt.Fprintf(d.buf, "$%d", v.ParamNum)
	case *SetToDefaultExpr:
		d.buf.WriteString("DEFAULT")
	case *NextValueExprQ:
		// pg: deparse as nextval('seqname'::regclass)
		fmt.Fprintf(d.buf, "nextval('%s'::regclass)", v.SeqName)
	case *ConvertRowtypeExprQ:
		// pg: ConvertRowtypeExpr — deparse the inner arg (conversion is implicit)
		d.getRuleExpr(v.Arg, showImplicit)
	case *FieldStoreExprQ:
		d.getFieldStore(v)
	case *GroupingFuncExpr:
		d.getGroupingFunc(v)
	case *XmlExprQ:
		d.getXmlExpr(v)
	default:
		d.buf.WriteString("???")
	}
}

// getVariable deparses a column reference.
//
// pg: src/backend/utils/adt/ruleutils.c — get_variable
func (d *deparseCtx) getVariable(v *VarExpr) {
	d.getVariableReturningName(v)
}

// getVariableReturningName deparses a column reference and returns the column name.
// pg: src/backend/utils/adt/ruleutils.c — get_variable
func (d *deparseCtx) getVariableReturningName(v *VarExpr) string {
	// For outer references (LevelsUp > 0), walk up the context chain.
	// pg: src/backend/utils/adt/ruleutils.c — get_variable uses
	// deparse_namespace which tracks multiple levels.
	ctx := d
	for i := 0; i < v.LevelsUp; i++ {
		if ctx.parentCtx == nil {
			break
		}
		ctx = ctx.parentCtx
	}

	if ctx.query == nil || v.RangeIdx >= len(ctx.query.RangeTable) {
		d.buf.WriteString("???")
		return "???"
	}

	rte := ctx.query.RangeTable[v.RangeIdx]
	colIdx := int(v.AttNum - 1)
	if colIdx < 0 || colIdx >= len(rte.ColNames) {
		d.buf.WriteString("???")
		return "???"
	}

	colName := rte.ColNames[colIdx]

	// For outer references, always use prefix.
	// For local references: pg: varprefix = (parentnamespace != NIL || list_length(query->rtable) != 1)
	needPrefix := v.LevelsUp > 0 || d.hasParentNamespace || len(d.query.RangeTable) != 1

	if needPrefix {
		d.buf.WriteString(quoteIdentifier(rte.ERef))
		d.buf.WriteByte('.')
	}
	d.buf.WriteString(quoteIdentifier(colName))

	return colName
}

// getConstExpr deparses a constant.
//
// pg: src/backend/utils/adt/ruleutils.c — get_const_expr
func (d *deparseCtx) getConstExpr(c *ConstExpr, showType int) {
	if c.IsNull {
		d.buf.WriteString("NULL")
		if showType >= 0 {
			d.buf.WriteString("::")
			d.buf.WriteString(d.catalog.formatType(c.TypeOID, c.TypeMod))
		}
		return
	}

	// pg: regclassout — strip schema prefix when relation is visible
	constValue := c.Value
	if c.TypeOID == REGCLASSOID && constValue != "" {
		constValue = d.catalog.regclassout(constValue)
	}

	needLabel := false

	switch c.TypeOID {
	case INT4OID:
		if len(c.Value) > 0 && c.Value[0] != '-' {
			d.buf.WriteString(c.Value)
		} else {
			fmt.Fprintf(d.buf, "'%s'", c.Value)
			needLabel = true
		}
	case INT8OID:
		// pg: get_const_expr — INT8 constants always shown as 'value'::bigint
		// because bare integers would be parsed as int4 and could overflow.
		simpleQuoteLiteral(d.buf, c.Value)
		needLabel = true
	case NUMERICOID:
		if len(c.Value) > 0 && c.Value[0] != '-' {
			d.buf.WriteString(c.Value)
		} else {
			fmt.Fprintf(d.buf, "'%s'", c.Value)
			needLabel = true
		}
	case BOOLOID:
		d.buf.WriteString(c.Value)
	default:
		simpleQuoteLiteral(d.buf, constValue)
		needLabel = true
	}

	if showType < 0 {
		return
	}

	// pg: get_const_expr — decide whether to show type label
	switch c.TypeOID {
	case BOOLOID, UNKNOWNOID:
		needLabel = false
	case INT4OID:
		// Already determined above.
	case INT8OID:
		needLabel = true
	case NUMERICOID:
		needLabel = needLabel || c.TypeMod >= 0
	default:
		needLabel = true
	}

	if needLabel || showType > 0 {
		d.buf.WriteString("::")
		d.buf.WriteString(d.catalog.formatType(c.TypeOID, c.TypeMod))
	}
}

// getFuncExpr deparses a function call expression.
//
// pg: src/backend/utils/adt/ruleutils.c — get_func_expr
func (d *deparseCtx) getFuncExpr(f *FuncCallExpr, showImplicit bool) {
	// Implicit cast: just show the argument.
	if f.CoerceFormat == 'i' && !showImplicit {
		if len(f.Args) > 0 {
			d.getRuleExprParen(f.Args[0], false, f)
		}
		return
	}

	// Explicit cast or shown implicit cast: show as expression::type.
	if f.CoerceFormat == 'e' || (f.CoerceFormat == 'i' && showImplicit) {
		if len(f.Args) > 0 {
			d.getCoercionExpr(f.Args[0], f.ResultType, f.ResultTypMod, f)
		}
		return
	}

	// AT TIME ZONE — PG deparses timezone() calls with special syntax.
	// Note reversed argument order: timezone(zone, ts) → (ts AT TIME ZONE zone).
	//
	// pg: src/backend/utils/adt/ruleutils.c — get_func_expr (timezone cases)
	switch f.FuncOID {
	case 1026, 1159, 2037, 2038, 2069, 2070:
		// 2-arg timezone: (arg2 AT TIME ZONE arg1)
		if len(f.Args) == 2 {
			d.buf.WriteByte('(')
			d.getRuleExprParen(f.Args[1], false, f)
			d.buf.WriteString(" AT TIME ZONE ")
			d.getRuleExprParen(f.Args[0], false, f)
			d.buf.WriteByte(')')
			return
		}
	case 6334, 6335, 6336:
		// 1-arg timezone: (arg1 AT LOCAL)
		if len(f.Args) == 1 {
			d.buf.WriteByte('(')
			d.getRuleExprParen(f.Args[0], false, f)
			d.buf.WriteString(" AT LOCAL)")
			return
		}
	}

	// Normal function call.
	d.buf.WriteString(d.generateFuncName(f))
	d.buf.WriteByte('(')
	for i, arg := range f.Args {
		if i > 0 {
			d.buf.WriteString(", ")
		}
		d.getRuleExpr(arg, true)
	}
	d.buf.WriteByte(')')
}

// getAggExpr deparses an aggregate function call.
//
// pg: src/backend/utils/adt/ruleutils.c — get_agg_expr
func (d *deparseCtx) getAggExpr(a *AggExpr) {
	d.buf.WriteString(a.AggName)
	d.buf.WriteByte('(')
	if a.AggStar {
		d.buf.WriteByte('*')
	} else {
		if a.AggDistinct {
			d.buf.WriteString("DISTINCT ")
		}
		for i, arg := range a.Args {
			if i > 0 {
				d.buf.WriteString(", ")
			}
			d.getRuleExpr(arg, true)
		}
	}
	d.buf.WriteByte(')')
}

// getOpExpr deparses an operator expression.
//
// pg: src/backend/utils/adt/ruleutils.c — get_oper_expr
func (d *deparseCtx) getOpExpr(o *OpExpr) {
	if !d.prettyParen {
		d.buf.WriteByte('(')
	}

	if o.Left != nil {
		// Binary operator.
		d.getRuleExprParen(o.Left, true, o)
		d.buf.WriteByte(' ')
		d.buf.WriteString(d.generateOpName(o))
		d.buf.WriteByte(' ')
		d.getRuleExprParen(o.Right, true, o)
	} else {
		// Prefix operator.
		d.buf.WriteString(d.generateOpName(o))
		d.buf.WriteByte(' ')
		d.getRuleExprParen(o.Right, true, o)
	}

	if !d.prettyParen {
		d.buf.WriteByte(')')
	}
}

// getRelabelExpr deparses a binary-compatible cast.
//
// pg: src/backend/utils/adt/ruleutils.c — RelabelType case
func (d *deparseCtx) getRelabelExpr(r *RelabelExpr, showImplicit bool) {
	if r.Format == 'i' && !showImplicit {
		d.getRuleExprParen(r.Arg, false, r)
		return
	}
	d.getCoercionExpr(r.Arg, r.ResultType, r.TypeMod, r)
}

// getCoerceViaIOExpr deparses an I/O conversion cast.
//
// pg: src/backend/utils/adt/ruleutils.c — CoerceViaIO case
func (d *deparseCtx) getCoerceViaIOExpr(c *CoerceViaIOExpr, showImplicit bool) {
	if c.Format == 'i' && !showImplicit {
		d.getRuleExprParen(c.Arg, false, c)
		return
	}
	d.getCoercionExpr(c.Arg, c.ResultType, -1, c)
}

// getCoercionExpr deparses a coercion as (expression)::type.
//
// pg: src/backend/utils/adt/ruleutils.c — get_coercion_expr
func (d *deparseCtx) getCoercionExpr(arg AnalyzedExpr, resultType uint32, typmod int32, parent AnalyzedExpr) {
	// pg: get_coercion_expr wraps arg in parens when !PRETTY_PAREN
	if !d.prettyParen {
		d.buf.WriteByte('(')
	}
	d.getRuleExprParen(arg, false, parent)
	if !d.prettyParen {
		d.buf.WriteByte(')')
	}
	d.buf.WriteString("::")
	d.buf.WriteString(d.catalog.formatType(resultType, typmod))
}

// getCaseExpr deparses a CASE expression.
//
// pg: src/backend/utils/adt/ruleutils.c — CaseExpr case
func (d *deparseCtx) getCaseExpr(c *CaseExprQ) {
	// pg: appendContextKeyword(context, "CASE", 0, PRETTYINDENT_VAR, 0)
	d.appendContextKeyword("CASE", 0, prettyIndentVar, 0)
	if c.Arg != nil {
		d.buf.WriteByte(' ')
		d.getRuleExpr(c.Arg, false)
	}
	for _, w := range c.When {
		// pg: appendContextKeyword(context, "WHEN ", 0, 0, 0)
		d.appendContextKeyword("WHEN ", 0, 0, 0)
		d.getRuleExpr(w.Condition, false)
		d.buf.WriteString(" THEN ")
		d.getRuleExpr(w.Result, true)
	}
	if c.Default != nil {
		d.appendContextKeyword("ELSE ", 0, 0, 0)
		d.getRuleExpr(c.Default, true)
	}
	// pg: appendContextKeyword(context, "END", -PRETTYINDENT_VAR, 0, 0)
	d.appendContextKeyword("END", -prettyIndentVar, 0, 0)
}

// getCoalesceExpr deparses a COALESCE expression.
//
// pg: src/backend/utils/adt/ruleutils.c — CoalesceExpr case
func (d *deparseCtx) getCoalesceExpr(c *CoalesceExprQ) {
	d.buf.WriteString("COALESCE(")
	for i, arg := range c.Args {
		if i > 0 {
			d.buf.WriteString(", ")
		}
		d.getRuleExpr(arg, true)
	}
	d.buf.WriteByte(')')
}

// getBoolExpr deparses AND/OR/NOT expressions.
//
// pg: src/backend/utils/adt/ruleutils.c — BoolExpr case
func (d *deparseCtx) getBoolExpr(b *BoolExprQ) {
	switch b.Op {
	case BoolAnd:
		if !d.prettyParen {
			d.buf.WriteByte('(')
		}
		for i, arg := range b.Args {
			if i > 0 {
				d.buf.WriteString(" AND ")
			}
			d.getRuleExprParen(arg, false, b)
		}
		if !d.prettyParen {
			d.buf.WriteByte(')')
		}
	case BoolOr:
		if !d.prettyParen {
			d.buf.WriteByte('(')
		}
		for i, arg := range b.Args {
			if i > 0 {
				d.buf.WriteString(" OR ")
			}
			d.getRuleExprParen(arg, false, b)
		}
		if !d.prettyParen {
			d.buf.WriteByte(')')
		}
	case BoolNot:
		if !d.prettyParen {
			d.buf.WriteByte('(')
		}
		d.buf.WriteString("NOT ")
		if len(b.Args) > 0 {
			d.getRuleExprParen(b.Args[0], false, b)
		}
		if !d.prettyParen {
			d.buf.WriteByte(')')
		}
	}
}

// getNullTest deparses IS [NOT] NULL.
//
// pg: src/backend/utils/adt/ruleutils.c — NullTest case
func (d *deparseCtx) getNullTest(n *NullTestExpr) {
	if !d.prettyParen {
		d.buf.WriteByte('(')
	}
	d.getRuleExprParen(n.Arg, true, n)
	if n.IsNull {
		d.buf.WriteString(" IS NULL")
	} else {
		d.buf.WriteString(" IS NOT NULL")
	}
	if !d.prettyParen {
		d.buf.WriteByte(')')
	}
}

// getSubLink deparses a subquery expression (scalar or EXISTS).
//
// pg: src/backend/utils/adt/ruleutils.c — get_sublink_expr
func (d *deparseCtx) getSubLink(s *SubLinkExpr) {
	switch s.SubLinkType {
	case SubLinkExistsType:
		// pg: get_sublink_expr — EXISTS wraps in parens when !prettyParen
		if !d.prettyParen {
			d.buf.WriteByte('(')
		}
		d.buf.WriteString("EXISTS (")
		subCtx := d.subContext()
		subCtx.hasParentNamespace = true
		subCtx.colNamesVisible = false // EXISTS doesn't care about column names
		subCtx.getQueryDef(s.SubQuery)
		d.indentLevel = subCtx.indentLevel
		d.buf.WriteByte(')')
		if !d.prettyParen {
			d.buf.WriteByte(')')
		}
	case SubLinkAnyType:
		// pg: get_sublink_expr — (testexpr IN ( subquery ))
		d.buf.WriteByte('(')
		if s.TestExpr != nil {
			d.getRuleExpr(s.TestExpr, true)
			d.buf.WriteString(" IN (")
		}
		subCtx := d.subContext()
		subCtx.hasParentNamespace = true
		subCtx.getQueryDef(s.SubQuery)
		d.indentLevel = subCtx.indentLevel
		if s.TestExpr != nil {
			d.buf.WriteString("))")
		} else {
			d.buf.WriteByte(')')
		}
	case SubLinkAllType:
		// pg: get_sublink_expr — (testexpr op ALL ( subquery ))
		d.buf.WriteByte('(')
		if s.TestExpr != nil {
			d.getRuleExpr(s.TestExpr, true)
			d.buf.WriteString(" <> ALL (")
		}
		subCtx := d.subContext()
		subCtx.hasParentNamespace = true
		subCtx.getQueryDef(s.SubQuery)
		d.indentLevel = subCtx.indentLevel
		if s.TestExpr != nil {
			d.buf.WriteString("))")
		} else {
			d.buf.WriteByte(')')
		}
	default:
		// EXPR_SUBLINK (scalar subquery)
		d.buf.WriteByte('(')
		subCtx := d.subContext()
		subCtx.hasParentNamespace = true
		subCtx.getQueryDef(s.SubQuery)
		d.indentLevel = subCtx.indentLevel
		d.buf.WriteByte(')')
	}
}

// getNullIf deparses a NULLIF expression.
//
// pg: src/backend/utils/adt/ruleutils.c — NullIfExpr case
func (d *deparseCtx) getNullIf(n *NullIfExprQ) {
	d.buf.WriteString("NULLIF(")
	if len(n.Args) >= 1 {
		d.getRuleExpr(n.Args[0], true)
	}
	if len(n.Args) >= 2 {
		d.buf.WriteString(", ")
		d.getRuleExpr(n.Args[1], true)
	}
	d.buf.WriteByte(')')
}

// getMinMax deparses a GREATEST/LEAST expression.
//
// pg: src/backend/utils/adt/ruleutils.c — MinMaxExpr case
func (d *deparseCtx) getMinMax(m *MinMaxExprQ) {
	switch m.Op {
	case MinMaxGreatest:
		d.buf.WriteString("GREATEST(")
	case MinMaxLeast:
		d.buf.WriteString("LEAST(")
	}
	for i, arg := range m.Args {
		if i > 0 {
			d.buf.WriteString(", ")
		}
		d.getRuleExpr(arg, true)
	}
	d.buf.WriteByte(')')
}

// getBooleanTest deparses IS [NOT] TRUE/FALSE/UNKNOWN.
//
// pg: src/backend/utils/adt/ruleutils.c — BooleanTest case
func (d *deparseCtx) getBooleanTest(b *BooleanTestExpr) {
	if !d.prettyParen {
		d.buf.WriteByte('(')
	}
	d.getRuleExprParen(b.Arg, true, b)
	switch b.TestType {
	case BoolIsTrue:
		d.buf.WriteString(" IS TRUE")
	case BoolIsNotTrue:
		d.buf.WriteString(" IS NOT TRUE")
	case BoolIsFalse:
		d.buf.WriteString(" IS FALSE")
	case BoolIsNotFalse:
		d.buf.WriteString(" IS NOT FALSE")
	case BoolIsUnknown:
		d.buf.WriteString(" IS UNKNOWN")
	case BoolIsNotUnknown:
		d.buf.WriteString(" IS NOT UNKNOWN")
	}
	if !d.prettyParen {
		d.buf.WriteByte(')')
	}
}

// getSQLValueFunc deparses CURRENT_DATE, CURRENT_TIMESTAMP, etc.
//
// pg: src/backend/utils/adt/ruleutils.c — SQLValueFunction case
func (d *deparseCtx) getSQLValueFunc(s *SQLValueFuncExpr) {
	switch s.Op {
	case SVFCurrentDate:
		d.buf.WriteString("CURRENT_DATE")
	case SVFCurrentTime:
		d.buf.WriteString("CURRENT_TIME")
	case SVFCurrentTimeN:
		fmt.Fprintf(d.buf, "CURRENT_TIME(%d)", s.TypeMod)
	case SVFCurrentTimestamp:
		d.buf.WriteString("CURRENT_TIMESTAMP")
	case SVFCurrentTimestampN:
		fmt.Fprintf(d.buf, "CURRENT_TIMESTAMP(%d)", s.TypeMod)
	case SVFLocaltime:
		d.buf.WriteString("LOCALTIME")
	case SVFLocaltimeN:
		fmt.Fprintf(d.buf, "LOCALTIME(%d)", s.TypeMod)
	case SVFLocaltimestamp:
		d.buf.WriteString("LOCALTIMESTAMP")
	case SVFLocaltimestampN:
		fmt.Fprintf(d.buf, "LOCALTIMESTAMP(%d)", s.TypeMod)
	case SVFCurrentRole:
		d.buf.WriteString("CURRENT_ROLE")
	case SVFCurrentUser:
		d.buf.WriteString("CURRENT_USER")
	case SVFUser:
		d.buf.WriteString("CURRENT_USER")
	case SVFSessionUser:
		d.buf.WriteString("SESSION_USER")
	case SVFCurrentCatalog:
		d.buf.WriteString("CURRENT_CATALOG")
	case SVFCurrentSchema:
		d.buf.WriteString("CURRENT_SCHEMA")
	}
}

// getDistinctExpr deparses IS [NOT] DISTINCT FROM.
// PG stores NOT DISTINCT FROM as NOT(DistinctExpr), so the deparser for
// IS NOT DISTINCT FROM shows: NOT (a IS DISTINCT FROM b).
//
// pg: src/backend/utils/adt/ruleutils.c — DistinctExpr case
func (d *deparseCtx) getDistinctExpr(de *DistinctExprQ) {
	if de.IsNot {
		// pg: IS NOT DISTINCT FROM is stored as BoolExpr(NOT, [DistinctExpr]).
		// PG deparses as: NOT <left> IS DISTINCT FROM <right> (no inner parens).
		if !d.prettyParen {
			d.buf.WriteByte('(')
		}
		d.buf.WriteString("NOT ")
		d.getRuleExprParen(de.Left, true, de)
		d.buf.WriteString(" IS DISTINCT FROM ")
		d.getRuleExprParen(de.Right, true, de)
		if !d.prettyParen {
			d.buf.WriteByte(')')
		}
		return
	}

	if !d.prettyParen {
		d.buf.WriteByte('(')
	}
	d.getRuleExprParen(de.Left, true, de)
	d.buf.WriteString(" IS DISTINCT FROM ")
	d.getRuleExprParen(de.Right, true, de)
	if !d.prettyParen {
		d.buf.WriteByte(')')
	}
}

// getScalarArrayOp deparses op ANY/ALL (array) expressions.
//
// pg: src/backend/utils/adt/ruleutils.c — ScalarArrayOpExpr case
func (d *deparseCtx) getScalarArrayOp(s *ScalarArrayOpExpr) {
	if !d.prettyParen {
		d.buf.WriteByte('(')
	}
	d.getRuleExprParen(s.Left, true, s)
	d.buf.WriteByte(' ')
	d.buf.WriteString(s.OpName)
	d.buf.WriteByte(' ')
	if s.UseOr {
		d.buf.WriteString("ANY (")
	} else {
		d.buf.WriteString("ALL (")
	}
	d.getRuleExpr(s.Right, true)
	d.buf.WriteByte(')')
	if !d.prettyParen {
		d.buf.WriteByte(')')
	}
}

// getWithClause deparses a WITH clause (CTEs).
//
// pg: src/backend/utils/adt/ruleutils.c — get_with_clause
func (d *deparseCtx) getWithClause(query *Query) {
	// pg: indent and prepend space before WITH keyword
	if d.prettyIndent {
		d.indentLevel += prettyIndentStd
		d.buf.WriteByte(' ')
	}

	sep := "WITH "
	if query.IsRecursive {
		sep = "WITH RECURSIVE "
	}

	for _, cte := range query.CTEList {
		d.buf.WriteString(sep)
		d.buf.WriteString(quoteIdentifier(cte.Name))

		// Column aliases.
		if len(cte.Aliases) > 0 {
			d.buf.WriteByte('(')
			for j, a := range cte.Aliases {
				if j > 0 {
					d.buf.WriteString(", ")
				}
				d.buf.WriteString(quoteIdentifier(a))
			}
			d.buf.WriteByte(')')
		}

		d.buf.WriteString(" AS ")

		// MATERIALIZED / NOT MATERIALIZED.
		switch cte.Materialized {
		case 1:
			d.buf.WriteString("MATERIALIZED ")
		case 2:
			d.buf.WriteString("NOT MATERIALIZED ")
		}

		d.buf.WriteByte('(')
		// pg: newline after '(' for pretty indentation
		if d.prettyIndent {
			d.appendContextKeyword("", 0, 0, 0)
		}
		subCtx := d.subContext()
		subCtx.hasParentNamespace = true
		subCtx.getQueryDef(cte.Query)
		d.indentLevel = subCtx.indentLevel
		// pg: newline before ')' for pretty indentation
		if d.prettyIndent {
			d.appendContextKeyword("", 0, 0, 0)
		}
		d.buf.WriteByte(')')

		sep = ", "
	}

	// pg: restore indent and add newline before the main SELECT
	if d.prettyIndent {
		d.indentLevel -= prettyIndentStd
		d.appendContextKeyword("", 0, 0, 0)
	} else {
		d.buf.WriteByte(' ')
	}
}

// getArrayExpr deparses an ARRAY[...] constructor.
//
// pg: src/backend/utils/adt/ruleutils.c — ArrayExpr case
func (d *deparseCtx) getArrayExpr(a *ArrayExprQ) {
	d.buf.WriteString("ARRAY[")
	for i, elem := range a.Elements {
		if i > 0 {
			d.buf.WriteString(", ")
		}
		d.getRuleExpr(elem, true)
	}
	d.buf.WriteByte(']')
}

// getRowExpr deparses a ROW(...) expression.
//
// pg: src/backend/utils/adt/ruleutils.c — RowExpr case
func (d *deparseCtx) getRowExpr(r *RowExprQ) {
	if r.RowFormat == 'e' {
		d.buf.WriteString("ROW(")
	} else {
		d.buf.WriteByte('(')
	}
	for i, arg := range r.Args {
		if i > 0 {
			d.buf.WriteString(", ")
		}
		d.getRuleExpr(arg, true)
	}
	d.buf.WriteByte(')')
}

// getCollateExpr deparses a COLLATE clause.
//
// pg: src/backend/utils/adt/ruleutils.c — CollateExpr case
func (d *deparseCtx) getCollateExpr(c *CollateExprQ) {
	if !d.prettyParen {
		d.buf.WriteByte('(')
	}
	d.getRuleExprParen(c.Arg, true, c)
	d.buf.WriteString(" COLLATE ")
	d.buf.WriteString(quoteIdentifier(c.CollName))
	if !d.prettyParen {
		d.buf.WriteByte(')')
	}
}

// getFieldSelect deparses a composite.field selection.
//
// pg: src/backend/utils/adt/ruleutils.c — FieldSelect case
func (d *deparseCtx) getFieldSelect(f *FieldSelectExprQ) {
	d.getRuleExprParen(f.Arg, true, f)
	d.buf.WriteByte('.')
	d.buf.WriteString(quoteIdentifier(f.FieldName))
}

// getWindowFunc deparses a window function call.
//
// pg: src/backend/utils/adt/ruleutils.c — WindowFunc case
func (d *deparseCtx) getWindowFunc(w *WindowFuncExpr) {
	d.buf.WriteString(quoteIdentifier(w.FuncName))
	d.buf.WriteByte('(')
	if w.AggStar {
		d.buf.WriteByte('*')
	} else {
		for i, arg := range w.Args {
			if i > 0 {
				d.buf.WriteString(", ")
			}
			d.getRuleExpr(arg, true)
		}
	}
	d.buf.WriteByte(')')

	// FILTER clause.
	if w.AggFilter != nil {
		d.buf.WriteString(" FILTER (WHERE ")
		d.getRuleExpr(w.AggFilter, false)
		d.buf.WriteByte(')')
	}

	// OVER clause.
	d.buf.WriteString(" OVER ")
	if d.query != nil && int(w.WinRef) < len(d.query.WindowClause) {
		wc := d.query.WindowClause[w.WinRef]
		if wc.Name != "" {
			d.buf.WriteString(quoteIdentifier(wc.Name))
			return
		}
	}

	// Inline window specification.
	d.buf.WriteByte('(')
	if d.query != nil && int(w.WinRef) < len(d.query.WindowClause) {
		d.getWindowSpec(d.query.WindowClause[w.WinRef])
	}
	d.buf.WriteByte(')')
}

// getWindowSpec deparses a window specification.
//
// pg: src/backend/utils/adt/ruleutils.c — get_rule_windowspec
func (d *deparseCtx) getWindowSpec(wc *WindowClauseQ) {
	needSpace := false

	if len(wc.PartitionBy) > 0 {
		d.buf.WriteString("PARTITION BY ")
		for i, pb := range wc.PartitionBy {
			if i > 0 {
				d.buf.WriteString(", ")
			}
			d.getSortGroupClause(pb, d.query.TargetList, false)
		}
		needSpace = true
	}

	if len(wc.OrderBy) > 0 {
		if needSpace {
			d.buf.WriteByte(' ')
		}
		d.buf.WriteString("ORDER BY ")
		for i, ob := range wc.OrderBy {
			if i > 0 {
				d.buf.WriteString(", ")
			}
			d.getSortGroupClause(ob, d.query.TargetList, false)
			if ob.Descending {
				d.buf.WriteString(" DESC")
				if !ob.NullsFirst {
					d.buf.WriteString(" NULLS LAST")
				}
			} else {
				if ob.NullsFirst {
					d.buf.WriteString(" NULLS FIRST")
				}
			}
		}
		needSpace = true
	}

	// Frame clause.
	if wc.FrameOptions != 0 {
		if needSpace {
			d.buf.WriteByte(' ')
		}
		d.deparseFrameOptions(wc)
	}
}

// deparseFrameOptions deparses a window frame clause.
//
// pg: src/backend/utils/adt/ruleutils.c — get_rule_windowclause (frame part)
func (d *deparseCtx) deparseFrameOptions(wc *WindowClauseQ) {
	// Frame options are a bitmask. These constants match PG's windowapi.h.
	const (
		frameOptRange       = 0x0002
		frameOptRows        = 0x0004
		frameOptGroups      = 0x0008
		frameOptStartUnbPre = 0x0010
		frameOptEndUnbFol   = 0x0020
		frameOptStartCurRow = 0x0040
		frameOptEndCurRow   = 0x0080
		frameOptStartOffset = 0x0100
		frameOptEndOffset   = 0x0200
		frameOptStartOffPre = 0x0800
		frameOptStartOffFol = 0x1000
		frameOptEndOffPre   = 0x2000
		frameOptEndOffFol   = 0x4000
		frameOptExcCurRow   = 0x8000
		frameOptExcGroup    = 0x10000
		frameOptExcTies     = 0x20000
		frameOptBetween     = 0x0400
	)

	fo := wc.FrameOptions

	// Type.
	switch {
	case fo&frameOptRange != 0:
		d.buf.WriteString("RANGE ")
	case fo&frameOptRows != 0:
		d.buf.WriteString("ROWS ")
	case fo&frameOptGroups != 0:
		d.buf.WriteString("GROUPS ")
	default:
		return // no frame
	}

	if fo&frameOptBetween != 0 {
		d.buf.WriteString("BETWEEN ")
	}

	// Start bound.
	switch {
	case fo&frameOptStartUnbPre != 0:
		d.buf.WriteString("UNBOUNDED PRECEDING ")
	case fo&frameOptStartCurRow != 0:
		d.buf.WriteString("CURRENT ROW ")
	case fo&frameOptStartOffset != 0:
		if wc.StartOffset != nil {
			d.getRuleExpr(wc.StartOffset, false)
		}
		if fo&frameOptStartOffPre != 0 {
			d.buf.WriteString(" PRECEDING ")
		} else {
			d.buf.WriteString(" FOLLOWING ")
		}
	}

	if fo&frameOptBetween != 0 {
		d.buf.WriteString("AND ")

		// End bound.
		switch {
		case fo&frameOptEndUnbFol != 0:
			d.buf.WriteString("UNBOUNDED FOLLOWING")
		case fo&frameOptEndCurRow != 0:
			d.buf.WriteString("CURRENT ROW")
		case fo&frameOptEndOffset != 0:
			if wc.EndOffset != nil {
				d.getRuleExpr(wc.EndOffset, false)
			}
			if fo&frameOptEndOffPre != 0 {
				d.buf.WriteString(" PRECEDING")
			} else {
				d.buf.WriteString(" FOLLOWING")
			}
		}
	}

	// Exclusion.
	switch {
	case fo&frameOptExcCurRow != 0:
		d.buf.WriteString(" EXCLUDE CURRENT ROW")
	case fo&frameOptExcGroup != 0:
		d.buf.WriteString(" EXCLUDE GROUP")
	case fo&frameOptExcTies != 0:
		d.buf.WriteString(" EXCLUDE TIES")
	}
}

// getSortGroupClause deparses a GROUP BY / ORDER BY reference.
//
// pg: src/backend/utils/adt/ruleutils.c — get_rule_sortgroupclause
func (d *deparseCtx) getSortGroupClause(sgc *SortGroupClause, tlist []*TargetEntry, forceColNo bool) {
	// Find the referenced target entry.
	var tle *TargetEntry
	for _, te := range tlist {
		if te.ResSortGroupRef == sgc.TLESortGroupRef {
			tle = te
			break
		}
	}
	if tle == nil {
		d.buf.WriteString("???")
		return
	}

	if forceColNo {
		// For set operations, use column number.
		fmt.Fprintf(d.buf, "%d", tle.ResNo)
		return
	}

	// pg: get_rule_sortgroupclause — for Const/Var, emit directly;
	// for anything else (function calls, aggregates, etc.), wrap in parens
	// to prevent misinterpretation as cube()/rollup() constructs.
	switch tle.Expr.(type) {
	case *ConstExpr:
		d.getRuleExpr(tle.Expr, true)
	case *VarExpr:
		d.getRuleExpr(tle.Expr, true)
	default:
		// Force parens for function-like expressions even when !prettyParen.
		needParen := d.prettyParen
		switch tle.Expr.(type) {
		case *FuncCallExpr, *AggExpr:
			needParen = true
		}
		if needParen {
			d.buf.WriteByte('(')
		}
		d.getRuleExpr(tle.Expr, true)
		if needParen {
			d.buf.WriteByte(')')
		}
	}
}

// getSubscriptingRef deparses an array subscript expression.
//
// pg: src/backend/utils/adt/ruleutils.c — SubscriptingRef case
func (d *deparseCtx) getSubscriptingRef(s *SubscriptingRefExpr) {
	d.getRuleExpr(s.ContainerExpr, true)
	if s.IsSlice {
		d.buf.WriteByte('[')
		if len(s.LowerExprs) > 0 && s.LowerExprs[0] != nil {
			d.getRuleExpr(s.LowerExprs[0], false)
		}
		d.buf.WriteByte(':')
		if len(s.SubscriptExprs) > 0 && s.SubscriptExprs[0] != nil {
			d.getRuleExpr(s.SubscriptExprs[0], false)
		}
		d.buf.WriteByte(']')
	} else {
		for _, sub := range s.SubscriptExprs {
			d.buf.WriteByte('[')
			d.getRuleExpr(sub, false)
			d.buf.WriteByte(']')
		}
	}
}

// getNamedArgExpr deparses a named argument expression.
//
// pg: src/backend/utils/adt/ruleutils.c — NamedArgExpr case
func (d *deparseCtx) getNamedArgExpr(n *NamedArgExprQ) {
	d.buf.WriteString(quoteIdentifier(n.Name))
	d.buf.WriteString(" => ")
	d.getRuleExpr(n.Arg, true)
}

// getArrayCoerceExpr deparses an array coercion expression.
//
// pg: src/backend/utils/adt/ruleutils.c — ArrayCoerceExpr case
func (d *deparseCtx) getArrayCoerceExpr(a *ArrayCoerceExprQ, showImplicit bool) {
	if a.Format == 'i' && !showImplicit {
		d.getRuleExprParen(a.Arg, false, a)
		return
	}
	d.getCoercionExpr(a.Arg, a.ResultType, -1, a)
}

// getRowCompareExpr deparses a row comparison expression.
//
// pg: src/backend/utils/adt/ruleutils.c — RowCompareExpr case
func (d *deparseCtx) getRowCompareExpr(r *RowCompareExprQ) {
	if !d.prettyParen {
		d.buf.WriteByte('(')
	}
	// Left row.
	d.buf.WriteByte('(')
	for i, arg := range r.LArgs {
		if i > 0 {
			d.buf.WriteString(", ")
		}
		d.getRuleExpr(arg, true)
	}
	d.buf.WriteString(") ")

	// Operator.
	d.buf.WriteString(rowCompareOpName(r.RCType))
	d.buf.WriteByte(' ')

	// Right row.
	d.buf.WriteByte('(')
	for i, arg := range r.RArgs {
		if i > 0 {
			d.buf.WriteString(", ")
		}
		d.getRuleExpr(arg, true)
	}
	d.buf.WriteByte(')')
	if !d.prettyParen {
		d.buf.WriteByte(')')
	}
}

// getFieldStore deparses a composite field store expression.
//
// pg: src/backend/utils/adt/ruleutils.c — FieldStore case
func (d *deparseCtx) getFieldStore(f *FieldStoreExprQ) {
	d.buf.WriteString("ROW(")
	for i, val := range f.NewVals {
		if i > 0 {
			d.buf.WriteString(", ")
		}
		d.getRuleExpr(val, true)
	}
	d.buf.WriteByte(')')
}

// getGroupingFunc deparses a GROUPING(col1, col2, ...) expression.
//
// pg: src/backend/utils/adt/ruleutils.c — GroupingFunc case
func (d *deparseCtx) getGroupingFunc(g *GroupingFuncExpr) {
	d.buf.WriteString("GROUPING(")
	for i, arg := range g.Args {
		if i > 0 {
			d.buf.WriteString(", ")
		}
		d.getRuleExpr(arg, true)
	}
	d.buf.WriteByte(')')
}

// getGroupingSet deparses a GROUPING SETS / ROLLUP / CUBE clause.
//
// pg: src/backend/utils/adt/ruleutils.c — get_rule_groupingset
func (d *deparseCtx) getGroupingSet(gs *GroupingSetQ, tlist []*TargetEntry) {
	switch gs.Kind {
	case GroupingSetSimple:
		if len(gs.Content) == 0 {
			// Empty grouping set: ()
			d.buf.WriteString("()")
		} else if len(gs.Content) == 1 {
			d.getSortGroupClause(gs.Content[0], tlist, false)
		} else {
			d.buf.WriteByte('(')
			for i, sgc := range gs.Content {
				if i > 0 {
					d.buf.WriteString(", ")
				}
				d.getSortGroupClause(sgc, tlist, false)
			}
			d.buf.WriteByte(')')
		}
	case GroupingSetRollup:
		d.buf.WriteString("ROLLUP(")
		for i, sgc := range gs.Content {
			if i > 0 {
				d.buf.WriteString(", ")
			}
			d.getSortGroupClause(sgc, tlist, false)
		}
		d.buf.WriteByte(')')
	case GroupingSetCube:
		d.buf.WriteString("CUBE(")
		for i, sgc := range gs.Content {
			if i > 0 {
				d.buf.WriteString(", ")
			}
			d.getSortGroupClause(sgc, tlist, false)
		}
		d.buf.WriteByte(')')
	case GroupingSetSets:
		d.buf.WriteString("GROUPING SETS(")
		for i, sub := range gs.Sets {
			if i > 0 {
				d.buf.WriteString(", ")
			}
			d.getGroupingSet(sub, tlist)
		}
		d.buf.WriteByte(')')
	}
}

// getXmlExpr deparses an XML expression.
//
// pg: src/backend/utils/adt/ruleutils.c — XmlExpr case
func (d *deparseCtx) getXmlExpr(x *XmlExprQ) {
	switch x.Op {
	case XmlOpConcat:
		d.buf.WriteString("XMLCONCAT(")
		for i, arg := range x.Args {
			if i > 0 {
				d.buf.WriteString(", ")
			}
			d.getRuleExpr(arg, true)
		}
		d.buf.WriteByte(')')

	case XmlOpElement:
		d.buf.WriteString("XMLELEMENT(NAME ")
		d.buf.WriteString(quoteIdentifier(x.Name))
		// XMLATTRIBUTES
		if len(x.NamedArgs) > 0 {
			d.buf.WriteString(", XMLATTRIBUTES(")
			for i, na := range x.NamedArgs {
				if i > 0 {
					d.buf.WriteString(", ")
				}
				d.getRuleExpr(na, true)
				if i < len(x.ArgNames) && x.ArgNames[i] != "" {
					d.buf.WriteString(" AS ")
					d.buf.WriteString(quoteIdentifier(x.ArgNames[i]))
				}
			}
			d.buf.WriteByte(')')
		}
		// Content args
		for _, arg := range x.Args {
			d.buf.WriteString(", ")
			d.getRuleExpr(arg, true)
		}
		d.buf.WriteByte(')')

	case XmlOpForest:
		d.buf.WriteString("XMLFOREST(")
		for i, na := range x.NamedArgs {
			if i > 0 {
				d.buf.WriteString(", ")
			}
			d.getRuleExpr(na, true)
			if i < len(x.ArgNames) && x.ArgNames[i] != "" {
				d.buf.WriteString(" AS ")
				d.buf.WriteString(quoteIdentifier(x.ArgNames[i]))
			}
		}
		d.buf.WriteByte(')')

	case XmlOpParse:
		d.buf.WriteString("XMLPARSE(")
		if x.Xmloption == XmlOptionDocument {
			d.buf.WriteString("DOCUMENT ")
		} else {
			d.buf.WriteString("CONTENT ")
		}
		if len(x.Args) > 0 {
			d.getRuleExpr(x.Args[0], false)
		}
		d.buf.WriteByte(')')

	case XmlOpPI:
		d.buf.WriteString("XMLPI(NAME ")
		d.buf.WriteString(quoteIdentifier(x.Name))
		if len(x.Args) > 0 {
			d.buf.WriteString(", ")
			d.getRuleExpr(x.Args[0], true)
		}
		d.buf.WriteByte(')')

	case XmlOpRoot:
		d.buf.WriteString("XMLROOT(")
		if len(x.Args) > 0 {
			d.getRuleExpr(x.Args[0], true)
		}
		if len(x.Args) > 1 {
			d.buf.WriteString(", VERSION ")
			// PG uses VERSION NO VALUE for null
			if c, ok := x.Args[1].(*ConstExpr); ok && c.IsNull {
				d.buf.WriteString("NO VALUE")
			} else {
				d.getRuleExpr(x.Args[1], false)
			}
		}
		if len(x.Args) > 2 {
			if c, ok := x.Args[2].(*ConstExpr); ok && !c.IsNull {
				switch c.Value {
				case "0":
					d.buf.WriteString(", STANDALONE YES")
				case "1":
					d.buf.WriteString(", STANDALONE NO")
				case "2":
					d.buf.WriteString(", STANDALONE NO VALUE")
				}
			}
		}
		d.buf.WriteByte(')')

	case XmlOpSerialize:
		d.buf.WriteString("XMLSERIALIZE(")
		if x.Xmloption == XmlOptionDocument {
			d.buf.WriteString("DOCUMENT ")
		} else {
			d.buf.WriteString("CONTENT ")
		}
		if len(x.Args) > 0 {
			d.getRuleExpr(x.Args[0], false)
		}
		d.buf.WriteString(" AS ")
		d.buf.WriteString(d.catalog.formatType(x.TypeOID, x.TypeMod))
		d.buf.WriteByte(')')

	case XmlOpDocument:
		// expr IS DOCUMENT
		if !d.prettyParen {
			d.buf.WriteByte('(')
		}
		if len(x.Args) > 0 {
			d.getRuleExprParen(x.Args[0], true, x)
		}
		d.buf.WriteString(" IS DOCUMENT")
		if !d.prettyParen {
			d.buf.WriteByte(')')
		}
	}
}

// getTablesample deparses a TABLESAMPLE clause.
//
// pg: src/backend/utils/adt/ruleutils.c — get_tablesample_def
func (d *deparseCtx) getTablesample(ts *TablesampleClauseQ) {
	d.buf.WriteString(" TABLESAMPLE ")
	d.buf.WriteString(ts.Method)
	d.buf.WriteByte('(')
	for i, arg := range ts.Args {
		if i > 0 {
			d.buf.WriteString(", ")
		}
		d.getRuleExpr(arg, false)
	}
	d.buf.WriteByte(')')
	if ts.Repeatable != nil {
		d.buf.WriteString(" REPEATABLE (")
		d.getRuleExpr(ts.Repeatable, false)
		d.buf.WriteByte(')')
	}
}

// getRuleExprParen wraps expression in parentheses if needed.
//
// pg: src/backend/utils/adt/ruleutils.c — get_rule_expr_paren
func (d *deparseCtx) getRuleExprParen(expr AnalyzedExpr, showImplicit bool, parent AnalyzedExpr) {
	needParen := d.prettyParen && !isSimpleNode(expr, parent, d.prettyParen)
	if needParen {
		d.buf.WriteByte('(')
	}
	d.getRuleExpr(expr, showImplicit)
	if needParen {
		d.buf.WriteByte(')')
	}
}

// isSimpleNode determines if an expression is "simple" and doesn't need parens.
//
// pg: src/backend/utils/adt/ruleutils.c — isSimpleNode
func isSimpleNode(expr AnalyzedExpr, parent AnalyzedExpr, prettyParen bool) bool {
	if expr == nil {
		return false
	}

	switch e := expr.(type) {
	case *VarExpr, *ConstExpr, *CoerceToDomainValueExpr, *SQLValueFuncExpr,
		*ParamExpr, *SetToDefaultExpr, *CaseTestExprQ:
		// single words: always simple
		return true

	case *ArrayExprQ, *RowExprQ, *CoalesceExprQ, *MinMaxExprQ,
		*NullIfExprQ, *AggExpr, *WindowFuncExpr, *FuncCallExpr,
		*SubscriptingRefExpr, *NextValueExprQ, *GroupingFuncExpr,
		*XmlExprQ:
		// function-like: name(..) or name[..] — always simple
		return true

	case *CaseExprQ:
		// CASE keywords act as parentheses
		return true

	case *FieldSelectExprQ:
		// pg: appears simple since . has top precedence, unless parent is
		// FieldSelect itself!
		_, parentIsFieldSelect := parent.(*FieldSelectExprQ)
		return !parentIsFieldSelect

	case *RelabelExpr:
		// pg: maybe simple, check args
		return isSimpleNode(e.Arg, expr, prettyParen)

	case *CoerceViaIOExpr:
		// pg: maybe simple, check args
		return isSimpleNode(e.Arg, expr, prettyParen)

	case *ArrayCoerceExprQ:
		return isSimpleNode(e.Arg, expr, prettyParen)

	case *CoerceToDomainExpr:
		return isSimpleNode(e.Arg, expr, prettyParen)

	case *ConvertRowtypeExprQ:
		return isSimpleNode(e.Arg, expr, prettyParen)

	case *OpExpr:
		// pg: depends on parent node type; needs further checking
		if prettyParen {
			if parentOp, ok := parent.(*OpExpr); ok {
				return isSimpleOpExpr(e, parentOp)
			}
		}
		// FALLTHROUGH to SubLink handling
		return isSimpleInParent(parent)

	case *SubLinkExpr, *NullTestExpr, *BooleanTestExpr, *DistinctExprQ:
		return isSimpleInParent(parent)

	case *BoolExprQ:
		switch p := parent.(type) {
		case *BoolExprQ:
			if prettyParen {
				switch e.Op {
				case BoolNot, BoolAnd:
					if p.Op == BoolAnd || p.Op == BoolOr {
						return true
					}
				case BoolOr:
					if p.Op == BoolOr {
						return true
					}
				}
			}
			return false
		case *FuncCallExpr:
			if p.CoerceFormat == 'e' || p.CoerceFormat == 'i' || p.CoerceFormat == 's' {
				return false
			}
			return true
		case *CoalesceExprQ, *MinMaxExprQ, *NullIfExprQ,
			*ArrayExprQ, *RowExprQ, *AggExpr, *WindowFuncExpr,
			*CaseExprQ:
			return true
		default:
			return false
		}
	}

	// those we don't know: in dubio complexo
	return false
}

// isSimpleInParent checks if OpExpr/SubLink/NullTest/BooleanTest/DistinctExpr
// is simple in the context of a given parent node type.
//
// pg: src/backend/utils/adt/ruleutils.c — isSimpleNode (T_SubLink case)
func isSimpleInParent(parent AnalyzedExpr) bool {
	switch p := parent.(type) {
	case *FuncCallExpr:
		// special handling for casts and COERCE_SQL_SYNTAX
		if p.CoerceFormat == 'e' || p.CoerceFormat == 'i' || p.CoerceFormat == 's' {
			return false
		}
		return true
	case *BoolExprQ: // lower precedence
		return true
	case *ArrayExprQ, *RowExprQ: // other separators
		return true
	case *CoalesceExprQ, *MinMaxExprQ: // own parentheses
		return true
	case *NullIfExprQ: // other separators
		return true
	case *AggExpr, *WindowFuncExpr: // own parentheses
		return true
	case *CaseExprQ: // other separators
		return true
	default:
		return false
	}
}

// isSimpleOpExpr implements operator precedence logic for OpExpr inside OpExpr.
//
// pg: src/backend/utils/adt/ruleutils.c — isSimpleNode (T_OpExpr with PRETTYFLAG_PAREN)
func isSimpleOpExpr(node *OpExpr, parent *OpExpr) bool {
	op := getSimpleBinaryOpName(node)
	if op == 0 {
		return false
	}

	isLoPriOp := op == '+' || op == '-'
	isHiPriOp := op == '*' || op == '/' || op == '%'
	if !isLoPriOp && !isHiPriOp {
		return false
	}

	parentOp := getSimpleBinaryOpName(parent)
	if parentOp == 0 {
		return false
	}

	isLoPriParent := parentOp == '+' || parentOp == '-'
	isHiPriParent := parentOp == '*' || parentOp == '/' || parentOp == '%'
	if !isLoPriParent && !isHiPriParent {
		return false
	}

	if isHiPriOp && isLoPriParent {
		return true // op binds tighter than parent
	}

	if isLoPriOp && isHiPriParent {
		return false
	}

	// Same priority — can skip parens only if we are the left arg: (a - b) - c, not a - (b - c).
	// pg: node == (Node *) linitial(((OpExpr *) parentNode)->args)
	if parent.Left == node {
		return true
	}

	return false
}

// getSimpleBinaryOpName returns the single-char operator name for a binary OpExpr,
// or 0 if the operator name is not a single character or the op is not binary.
//
// pg: src/backend/utils/adt/ruleutils.c — get_simple_binary_op_name
func getSimpleBinaryOpName(op *OpExpr) byte {
	// Must be binary (has both Left and Right).
	if op.Left == nil || op.Right == nil {
		return 0
	}
	if len(op.OpName) == 1 {
		return op.OpName[0]
	}
	return 0
}

// --- Sub-context helpers ---

// subContext creates a child deparse context sharing the same output buffer.
func (d *deparseCtx) subContext() *deparseCtx {
	return &deparseCtx{
		catalog:            d.catalog,
		buf:                d.buf,
		indentLevel:        d.indentLevel,
		prettyIndent:       d.prettyIndent,
		prettyParen:        d.prettyParen,
		hasParentNamespace: d.hasParentNamespace,
		wrapColumn:         d.wrapColumn,
		colNamesVisible:    d.colNamesVisible,
		parentCtx:          d, // link to parent for outer Var resolution
	}
}

// subContextBuf creates a child deparse context writing to a different buffer.
func (d *deparseCtx) subContextBuf(buf *strings.Builder) *deparseCtx {
	return &deparseCtx{
		catalog:            d.catalog,
		buf:                buf,
		indentLevel:        d.indentLevel,
		prettyIndent:       d.prettyIndent,
		prettyParen:        d.prettyParen,
		hasParentNamespace: d.hasParentNamespace,
		wrapColumn:         d.wrapColumn,
		colNamesVisible:    d.colNamesVisible,
		query:              d.query,
		parentCtx:          d.parentCtx,
	}
}

// --- Name generation helpers ---

// generateRelationName generates the table name for deparse.
//
// pg: src/backend/utils/adt/ruleutils.c — generate_relation_name
func (d *deparseCtx) generateRelationName(rte *RangeTableEntry) string {
	if rte.SchemaName != "" && rte.SchemaName != "public" {
		return quoteIdentifier(rte.SchemaName) + "." + quoteIdentifier(rte.RelName)
	}
	return quoteIdentifier(rte.RelName)
}

// generateFuncName generates a function name for deparse.
//
// pg: src/backend/utils/adt/ruleutils.c — generate_function_name
func (d *deparseCtx) generateFuncName(f *FuncCallExpr) string {
	return quoteIdentifier(f.FuncName)
}

// generateOpName generates an operator name for deparse.
func (d *deparseCtx) generateOpName(o *OpExpr) string {
	return o.OpName
}

// formatType returns the SQL-standard type name for a type OID.
//
// pg: src/backend/utils/adt/format_type.c — format_type_with_typemod
func (c *Catalog) formatType(typeOID uint32, typmod int32) string {
	switch typeOID {
	case BOOLOID:
		return "boolean"
	case INT2OID:
		return "smallint"
	case INT4OID:
		return "integer"
	case INT8OID:
		return "bigint"
	case FLOAT4OID:
		return "real"
	case FLOAT8OID:
		return "double precision"
	case NUMERICOID:
		if typmod >= 0 {
			precision := ((typmod - 4) >> 16) & 0xFFFF
			scale := (typmod - 4) & 0xFFFF
			if scale > 0 {
				return fmt.Sprintf("numeric(%d,%d)", precision, scale)
			}
			return fmt.Sprintf("numeric(%d)", precision)
		}
		return "numeric"
	case TEXTOID:
		return "text"
	case VARCHAROID:
		if typmod >= 0 {
			return fmt.Sprintf("character varying(%d)", typmod-4)
		}
		return "character varying"
	case BPCHAROID:
		if typmod >= 0 {
			return fmt.Sprintf("character(%d)", typmod-4)
		}
		return "character"
	case DATEOID:
		return "date"
	case TIMEOID:
		return "time without time zone"
	case TIMETZOID:
		return "time with time zone"
	case TIMESTAMPOID:
		return "timestamp without time zone"
	case TIMESTAMPTZOID:
		return "timestamp with time zone"
	case INTERVALOID:
		return "interval"
	case BYTEAOID:
		return "bytea"
	case UUIDOID:
		return "uuid"
	case JSONOID:
		return "json"
	case JSONBOID:
		return "jsonb"
	case XMLOID:
		return "xml"
	case INETOID:
		return "inet"
	case CIDROID:
		return "cidr"
	case MACADDROID:
		return "macaddr"
	case MONEYOID:
		return "money"
	case BITOID:
		if typmod >= 0 {
			return fmt.Sprintf("bit(%d)", typmod)
		}
		return "bit"
	case VARBITOID:
		if typmod >= 0 {
			return fmt.Sprintf("bit varying(%d)", typmod)
		}
		return "bit varying"
	case UNKNOWNOID:
		return "unknown"
	case NAMEOID:
		return "name"
	case OIDOID:
		return "oid"
	case REGCLASSOID:
		return "regclass"
	case REGTYPEOID:
		return "regtype"
	case TSVECTOROID:
		return "tsvector"
	case TSQUERYOID:
		return "tsquery"
	}

	// For array types, check elem type.
	if t := c.typeByOID[typeOID]; t != nil {
		if t.Category == 'A' && t.Elem != 0 {
			return c.formatType(t.Elem, -1) + "[]"
		}
		// pg: format_type_with_typemod — schema-qualify if not visible
		if !c.TypeIsVisible(typeOID) {
			for _, s := range c.schemaByName {
				if s.OID == t.Namespace {
					return quoteIdentifier(s.Name) + "." + quoteIdentifier(t.TypeName)
				}
			}
		}
		return t.TypeName
	}

	return fmt.Sprintf("oid(%d)", typeOID)
}

// formatTypeOID returns the display name for a type OID without typmod.
//
// (pgddl helper — convenience wrapper around formatType)
func (c *Catalog) formatTypeOID(typeOID uint32) string {
	return c.formatType(typeOID, -1)
}

// appendContextKeyword appends a keyword with proper indentation.
//
// pg: src/backend/utils/adt/ruleutils.c — appendContextKeyword
func (d *deparseCtx) appendContextKeyword(str string, indentBefore, indentAfter, indentPlus int) {
	if d.prettyIndent {
		d.indentLevel += indentBefore
		// pg: removeStringInfoSpaces — trim trailing spaces before newline
		trimBufTrailingSpaces(d.buf)
		d.buf.WriteByte('\n')
		indentAmount := d.indentLevel + indentPlus
		if indentAmount < 0 {
			indentAmount = 0
		}
		for i := 0; i < indentAmount; i++ {
			d.buf.WriteByte(' ')
		}
		d.indentLevel += indentAfter
	}
	d.buf.WriteString(str)
}

// trimBufTrailingSpaces removes trailing space characters from a strings.Builder.
// pg: src/backend/utils/adt/ruleutils.c — removeStringInfoSpaces
func trimBufTrailingSpaces(buf *strings.Builder) {
	s := buf.String()
	trimmed := strings.TrimRight(s, " ")
	if len(trimmed) != len(s) {
		buf.Reset()
		buf.WriteString(trimmed)
	}
}

// quoteIdentifier quotes an identifier if needed.
//
// pg: src/backend/utils/adt/ruleutils.c — quote_identifier
func quoteIdentifier(name string) string {
	if name == "" {
		return name
	}
	if needsQuoting(name) {
		return "\"" + strings.ReplaceAll(name, "\"", "\"\"") + "\""
	}
	return name
}

// needsQuoting checks if an identifier needs quoting.
func needsQuoting(name string) bool {
	if len(name) == 0 {
		return true
	}
	if !isIdentStart(name[0]) {
		return true
	}
	for i := 1; i < len(name); i++ {
		if !isIdentCont(name[i]) {
			return true
		}
	}
	if isReservedWord(name) {
		return true
	}
	for _, c := range name {
		if c >= 'A' && c <= 'Z' {
			return true
		}
	}
	return false
}

func isIdentStart(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isIdentCont(c byte) bool {
	return isIdentStart(c) || (c >= '0' && c <= '9') || c == '$'
}

// isReservedWord checks if an identifier needs quoting per PG's quote_identifier.
// PG quotes everything that is not UNRESERVED_KEYWORD. This includes
// RESERVED_KEYWORD, COL_NAME_KEYWORD, and TYPE_FUNC_NAME_KEYWORD.
//
// pg: src/include/parser/kwlist.h
func isReservedWord(name string) bool {
	switch strings.ToLower(name) {
	// RESERVED_KEYWORD
	case "all", "analyse", "analyze", "and", "any", "array", "as", "asc",
		"asymmetric", "both", "case", "cast", "check", "collate",
		"column", "constraint", "create", "current_catalog",
		"current_date", "current_role", "current_time",
		"current_timestamp", "current_user", "default", "deferrable",
		"desc", "distinct", "do", "else", "end", "except", "false",
		"fetch", "for", "foreign", "from", "grant", "group",
		"having", "in", "initially", "intersect", "into",
		"lateral", "leading", "limit",
		"localtime", "localtimestamp", "not", "null",
		"offset", "on", "only", "or", "order",
		"placing", "primary", "references", "returning",
		"select", "session_user", "some", "symmetric",
		"table", "then", "to", "trailing", "true",
		"union", "unique", "user", "using", "variadic",
		"when", "where", "window", "with":
		return true
	// COL_NAME_KEYWORD
	case "between", "bigint", "bit", "boolean", "char", "character",
		"coalesce", "dec", "decimal", "exists", "extract",
		"float", "greatest", "grouping", "inout", "int", "integer",
		"interval", "json",
		"json_array", "json_arrayagg", "json_exists",
		"json_object", "json_objectagg", "json_query",
		"json_scalar", "json_serialize", "json_table", "json_value",
		"least", "national", "nchar", "none", "normalize", "nullif",
		"numeric", "out", "overlay", "position", "precision",
		"real", "row", "setof", "smallint", "substring", "time",
		"timestamp", "treat", "trim", "values", "varchar",
		"xmlattributes", "xmlconcat", "xmlelement", "xmlexists",
		"xmlforest", "xmlnamespaces", "xmlparse", "xmlpi",
		"xmlroot", "xmlserialize", "xmltable":
		return true
	// TYPE_FUNC_NAME_KEYWORD
	case "authorization", "binary", "collation", "concurrently",
		"cross", "current_schema", "freeze", "full", "ilike",
		"inner", "is", "isnull", "join", "left", "like",
		"natural", "notnull", "outer", "overlaps", "right",
		"similar", "tablesample", "verbose":
		return true
	}
	return false
}

// simpleQuoteLiteral writes a single-quoted literal string.
func simpleQuoteLiteral(buf *strings.Builder, val string) {
	buf.WriteByte('\'')
	for _, c := range val {
		if c == '\'' {
			buf.WriteString("''")
		} else {
			buf.WriteRune(c)
		}
	}
	buf.WriteByte('\'')
}
