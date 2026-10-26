package catalog

// AnalyzedExprChildren returns the expressions nested directly in expr, in
// field order, leaving out nil ones. It does not enter a query: the
// expressions of SubLinkExpr.SubQuery belong to that *Query, so a caller that
// follows subqueries reads them from it with its own range table.
//
// Every AnalyzedExpr node is handled here.
// TestAnalyzedExprChildrenCoversEveryField fails when a node, or an
// expression field of one, is added without it.
//
// pg: src/backend/nodes/nodeFuncs.c — expression_tree_walker
func AnalyzedExprChildren(expr AnalyzedExpr) []AnalyzedExpr {
	var out []AnalyzedExpr
	add := func(exprs ...AnalyzedExpr) {
		for _, e := range exprs {
			if e != nil {
				out = append(out, e)
			}
		}
	}
	switch v := expr.(type) {
	case *FuncCallExpr:
		add(v.Args...)
	case *AggExpr:
		add(v.Args...)
		for _, sb := range v.AggOrder {
			if sb != nil {
				add(sb.Expr)
			}
		}
		add(v.AggFilter)
	case *WindowFuncExpr:
		add(v.Args...)
		add(v.AggFilter)
	case *OpExpr:
		add(v.Left, v.Right)
	case *DistinctExprQ:
		add(v.Left, v.Right)
	case *ScalarArrayOpExpr:
		add(v.Left, v.Right)
	case *NullIfExprQ:
		add(v.Args...)
	case *RelabelExpr:
		add(v.Arg)
	case *CoerceViaIOExpr:
		add(v.Arg)
	case *CaseExprQ:
		add(v.Arg)
		for _, w := range v.When {
			if w != nil {
				add(w.Condition, w.Result)
			}
		}
		add(v.Default)
	case *CoalesceExprQ:
		add(v.Args...)
	case *BoolExprQ:
		add(v.Args...)
	case *NullTestExpr:
		add(v.Arg)
	case *BooleanTestExpr:
		add(v.Arg)
	case *SubLinkExpr:
		add(v.TestExpr)
	case *MinMaxExprQ:
		add(v.Args...)
	case *ArrayExprQ:
		add(v.Elements...)
	case *RowExprQ:
		add(v.Args...)
	case *CollateExprQ:
		add(v.Arg)
	case *FieldSelectExprQ:
		add(v.Arg)
	case *SubscriptingRefExpr:
		add(v.ContainerExpr)
		add(v.SubscriptExprs...)
		add(v.LowerExprs...)
	case *NamedArgExprQ:
		add(v.Arg)
	case *ArrayCoerceExprQ:
		add(v.Arg, v.ElemExpr)
	case *CoerceToDomainExpr:
		add(v.Arg)
	case *RowCompareExprQ:
		add(v.LArgs...)
		add(v.RArgs...)
	case *ConvertRowtypeExprQ:
		add(v.Arg)
	case *FieldStoreExprQ:
		add(v.Arg)
		add(v.NewVals...)
	case *GroupingFuncExpr:
		add(v.Args...)
	case *XmlExprQ:
		add(v.NamedArgs...)
		add(v.Args...)
	case *JsonConstructorExprQ:
		add(v.Args...)
	case *JsonExprQ:
		add(v.Expr)
	case *JsonIsPredicateExpr:
		add(v.Expr)
	case *JsonValueExprQ:
		add(v.Expr)
	case *InferenceElemExpr:
		add(v.Expr)
	case *PartitionBoundSpecExpr:
		add(v.LowerBound...)
		add(v.UpperBound...)
		add(v.ListValues...)
	case *VarExpr, *ConstExpr, *SQLValueFuncExpr, *CoerceToDomainValueExpr,
		*CaseTestExprQ, *ParamExpr, *SetToDefaultExpr, *NextValueExprQ,
		*SetOpColumnExpr, *SubPlanExpr, *AlternativeSubPlanExpr, *MergeSupportFuncExpr:
		// Leaves.
	}
	return out
}

// WalkAnalyzedExpr calls fn for expr and, when fn returns true, walks each
// expression nested in it the same way, depth first in field order. It does
// not enter subqueries; see AnalyzedExprChildren.
func WalkAnalyzedExpr(expr AnalyzedExpr, fn func(AnalyzedExpr) bool) {
	if expr == nil || !fn(expr) {
		return
	}
	for _, child := range AnalyzedExprChildren(expr) {
		WalkAnalyzedExpr(child, fn)
	}
}
