package catalog

import (
	"math"
	"strconv"
	"strings"
	"unicode/utf8"

	nodes "github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/mysql/deparse"
)

// analyzeExpr is the main expression analysis dispatcher.
func analyzeExpr(c *Catalog, expr nodes.ExprNode, scope *analyzerScope) (AnalyzedExpr, error) {
	switch n := expr.(type) {
	case *nodes.ColumnRef:
		return analyzeColumnRef(n, scope)
	case *nodes.IntLit:
		bt := BaseTypeBigInt
		if n.Value >= math.MinInt32 && n.Value <= math.MaxInt32 {
			bt = BaseTypeInt
		}
		return &ConstExprQ{Value: strconv.FormatInt(n.Value, 10), Type: &ResolvedType{BaseType: bt}}, nil
	case *nodes.StringLit:
		return &ConstExprQ{
			Value:      n.Value,
			IsString:   true,
			Introducer: n.Charset,
			Type:       &ResolvedType{BaseType: BaseTypeVarchar, Length: utf8.RuneCountInString(n.Value)},
		}, nil
	case *nodes.FloatLit:
		return &ConstExprQ{Value: n.Value, Type: floatLitResolvedType(n.Value)}, nil
	case *nodes.NullLit:
		return &ConstExprQ{IsNull: true, Value: "NULL", Type: &ResolvedType{BaseType: BaseTypeBinary}}, nil
	case *nodes.BoolLit:
		val := "FALSE"
		if n.Value {
			val = "TRUE"
		}
		return &ConstExprQ{Value: val, Type: &ResolvedType{BaseType: BaseTypeInt}}, nil
	case *nodes.KeywordArg:
		return &ConstExprQ{Value: n.Keyword, Keyword: true}, nil
	case *nodes.ParamMarkerExpr:
		// A `?` parameter marker cannot appear in a view definition — the only
		// thing this analyzer processes — and the engine rejects such a view.
		// Reject it here rather than carry a marker through the normalizing
		// deparse, which cannot round-trip its positional binding (a comma-form
		// LIMIT ?,? reverses the count/offset bindings).
		return nil, &Error{
			Code:     1351,
			SQLState: "HY000",
			Message:  "View's SELECT contains a '?' (parameter marker)",
		}
	case *nodes.FuncCallExpr:
		return analyzeFuncCall(c, n, scope)
	case *nodes.ParenExpr:
		return analyzeExpr(c, n.Expr, scope)
	case *nodes.BinaryExpr:
		return analyzeBinaryExpr(c, n, scope)
	case *nodes.UnaryExpr:
		return analyzeUnaryExpr(c, n, scope)
	case *nodes.InExpr:
		return analyzeInExpr(c, n, scope)
	case *nodes.BetweenExpr:
		return analyzeBetweenExpr(c, n, scope)
	case *nodes.CollateExpr:
		arg, err := analyzeExpr(c, n.Expr, scope)
		if err != nil {
			return nil, err
		}
		return &CollateExprQ{Arg: arg, Collation: toLower(n.Collation)}, nil
	case *nodes.LikeExpr:
		return analyzeLikeExpr(c, n, scope)
	case *nodes.IsExpr:
		return analyzeIsExpr(c, n, scope)
	case *nodes.CaseExpr:
		return analyzeCaseExpr(c, n, scope)
	case *nodes.CastExpr:
		return analyzeCastExpr(c, n, scope)
	case *nodes.PositionExpr:
		return analyzePositionExpr(c, n, scope)
	case *nodes.IntervalExpr:
		val, err := analyzeExpr(c, n.Value, scope)
		if err != nil {
			return nil, err
		}
		return &IntervalExprQ{Value: val, Unit: n.Unit}, nil
	case *nodes.SubqueryExpr:
		return analyzeScalarSubquery(c, n, scope)
	case *nodes.ExistsExpr:
		return analyzeExistsSubquery(c, n, scope)
	default:
		return nil, &Error{
			Code:     0,
			SQLState: "HY000",
			Message:  "unsupported expression type in analyzer",
		}
	}
}

func analyzeLikeExpr(c *Catalog, expr *nodes.LikeExpr, scope *analyzerScope) (AnalyzedExpr, error) {
	arg, err := analyzeExpr(c, expr.Expr, scope)
	if err != nil {
		return nil, err
	}
	pattern, err := analyzeExpr(c, expr.Pattern, scope)
	if err != nil {
		return nil, err
	}
	var escape AnalyzedExpr
	if expr.Escape != nil {
		escape, err = analyzeExpr(c, expr.Escape, scope)
		if err != nil {
			return nil, err
		}
	}
	return &LikeExprQ{Arg: arg, Pattern: pattern, Escape: escape, Negated: expr.Not}, nil
}

func floatLitResolvedType(raw string) *ResolvedType {
	if strings.ContainsAny(raw, "eE") {
		return &ResolvedType{BaseType: BaseTypeDouble}
	}
	digits := strings.TrimLeft(raw, "+-")
	intPart, fracPart, hasDot := strings.Cut(digits, ".")
	if !hasDot {
		return &ResolvedType{BaseType: BaseTypeBigInt}
	}
	return &ResolvedType{
		BaseType:  BaseTypeDecimal,
		Precision: len(intPart) + len(fracPart),
		Scale:     len(fracPart),
	}
}

// analyzeColumnRef resolves a column reference against the scope.
// Uses the Full resolution methods to support correlated subquery references
// (setting LevelsUp > 0 when the column comes from a parent scope).
func analyzeColumnRef(ref *nodes.ColumnRef, scope *analyzerScope) (AnalyzedExpr, error) {
	if ref.Table != "" {
		rteIdx, attNum, levelsUp, err := scope.resolveSchemaQualifiedColumnFull(ref.Schema, ref.Table, ref.Column)
		if err != nil {
			return nil, err
		}
		return &VarExprQ{RangeIdx: rteIdx, AttNum: attNum, LevelsUp: levelsUp}, nil
	}
	rteIdx, attNum, levelsUp, aliasExpr, err := scope.resolveColumnOrOuterAlias(ref.Column)
	if err != nil {
		return nil, err
	}
	if aliasExpr != nil {
		return shiftExprLevelsUp(aliasExpr, levelsUp, 0), nil
	}
	return &VarExprQ{RangeIdx: rteIdx, AttNum: attNum, LevelsUp: levelsUp}, nil
}

// shiftExprLevelsUp returns a copy of expr with every VarExprQ crossing above
// minDepth shifted by +delta levels, so an outer select-list alias expression
// keeps valid LevelsUp coordinates inside a nested subquery.
func shiftExprLevelsUp(expr AnalyzedExpr, delta, minDepth int) AnalyzedExpr {
	if expr == nil || delta == 0 {
		return expr
	}
	switch e := expr.(type) {
	case *VarExprQ:
		if e.LevelsUp < minDepth {
			return e
		}
		cp := *e
		cp.LevelsUp += delta
		return &cp
	case *ConstExprQ:
		return e
	case *FuncCallExprQ:
		cp := *e
		if len(e.Args) > 0 {
			cp.Args = make([]AnalyzedExpr, len(e.Args))
			for i, a := range e.Args {
				cp.Args[i] = shiftExprLevelsUp(a, delta, minDepth)
			}
		}
		if e.Over != nil && len(e.Over.PartitionBy) > 0 {
			ov := *e.Over
			ov.PartitionBy = make([]AnalyzedExpr, len(e.Over.PartitionBy))
			for i, p := range e.Over.PartitionBy {
				ov.PartitionBy[i] = shiftExprLevelsUp(p, delta, minDepth)
			}
			cp.Over = &ov
		}
		return &cp
	case *OpExprQ:
		cp := *e
		cp.Left = shiftExprLevelsUp(e.Left, delta, minDepth)
		cp.Right = shiftExprLevelsUp(e.Right, delta, minDepth)
		return &cp
	case *BoolExprQ:
		cp := *e
		if len(e.Args) > 0 {
			cp.Args = make([]AnalyzedExpr, len(e.Args))
			for i, a := range e.Args {
				cp.Args[i] = shiftExprLevelsUp(a, delta, minDepth)
			}
		}
		return &cp
	case *CaseExprQ:
		cp := *e
		cp.TestExpr = shiftExprLevelsUp(e.TestExpr, delta, minDepth)
		cp.Default = shiftExprLevelsUp(e.Default, delta, minDepth)
		if len(e.Args) > 0 {
			cp.Args = make([]*CaseWhenQ, len(e.Args))
			for i, w := range e.Args {
				cp.Args[i] = &CaseWhenQ{
					Cond: shiftExprLevelsUp(w.Cond, delta, minDepth),
					Then: shiftExprLevelsUp(w.Then, delta, minDepth),
				}
			}
		}
		return &cp
	case *CoalesceExprQ:
		cp := *e
		if len(e.Args) > 0 {
			cp.Args = make([]AnalyzedExpr, len(e.Args))
			for i, a := range e.Args {
				cp.Args[i] = shiftExprLevelsUp(a, delta, minDepth)
			}
		}
		return &cp
	case *NullTestExprQ:
		cp := *e
		cp.Arg = shiftExprLevelsUp(e.Arg, delta, minDepth)
		return &cp
	case *SubLinkExprQ:
		cp := *e
		cp.TestExpr = shiftExprLevelsUp(e.TestExpr, delta, minDepth)
		cp.Subquery = shiftQueryLevelsUp(e.Subquery, delta, minDepth+1)
		return &cp
	case *InListExprQ:
		cp := *e
		cp.Arg = shiftExprLevelsUp(e.Arg, delta, minDepth)
		if len(e.List) > 0 {
			cp.List = make([]AnalyzedExpr, len(e.List))
			for i, a := range e.List {
				cp.List[i] = shiftExprLevelsUp(a, delta, minDepth)
			}
		}
		return &cp
	case *BetweenExprQ:
		cp := *e
		cp.Arg = shiftExprLevelsUp(e.Arg, delta, minDepth)
		cp.Lower = shiftExprLevelsUp(e.Lower, delta, minDepth)
		cp.Upper = shiftExprLevelsUp(e.Upper, delta, minDepth)
		return &cp
	case *LikeExprQ:
		cp := *e
		cp.Arg = shiftExprLevelsUp(e.Arg, delta, minDepth)
		cp.Pattern = shiftExprLevelsUp(e.Pattern, delta, minDepth)
		cp.Escape = shiftExprLevelsUp(e.Escape, delta, minDepth)
		return &cp
	case *CollateExprQ:
		cp := *e
		cp.Arg = shiftExprLevelsUp(e.Arg, delta, minDepth)
		return &cp
	case *RowExprQ:
		cp := *e
		if len(e.Args) > 0 {
			cp.Args = make([]AnalyzedExpr, len(e.Args))
			for i, a := range e.Args {
				cp.Args[i] = shiftExprLevelsUp(a, delta, minDepth)
			}
		}
		return &cp
	case *CastExprQ:
		cp := *e
		cp.Arg = shiftExprLevelsUp(e.Arg, delta, minDepth)
		return &cp
	case *IntervalExprQ:
		cp := *e
		cp.Value = shiftExprLevelsUp(e.Value, delta, minDepth)
		return &cp
	default:
		return expr
	}
}

func shiftQueryLevelsUp(q *Query, delta, minDepth int) *Query {
	if q == nil || delta == 0 {
		return q
	}
	cp := *q
	if len(q.RangeTable) > 0 {
		cp.RangeTable = make([]*RangeTableEntryQ, len(q.RangeTable))
		for i, rte := range q.RangeTable {
			if rte != nil && rte.Subquery != nil {
				rteCp := *rte
				rteCp.Subquery = shiftQueryLevelsUp(rte.Subquery, delta, minDepth+1)
				cp.RangeTable[i] = &rteCp
			} else {
				cp.RangeTable[i] = rte
			}
		}
	}
	if len(q.CTEList) > 0 {
		cp.CTEList = make([]*CommonTableExprQ, len(q.CTEList))
		for i, cte := range q.CTEList {
			if cte != nil && cte.Query != nil {
				cteCp := *cte
				cteCp.Query = shiftQueryLevelsUp(cte.Query, delta, minDepth)
				cp.CTEList[i] = &cteCp
			} else {
				cp.CTEList[i] = cte
			}
		}
	}
	if len(q.TargetList) > 0 {
		cp.TargetList = make([]*TargetEntryQ, len(q.TargetList))
		for i, te := range q.TargetList {
			teCp := *te
			teCp.Expr = shiftExprLevelsUp(te.Expr, delta, minDepth)
			cp.TargetList[i] = &teCp
		}
	}
	if q.JoinTree != nil {
		jt := *q.JoinTree
		jt.Quals = shiftExprLevelsUp(jt.Quals, delta, minDepth)
		cp.JoinTree = &jt
	}
	cp.HavingQual = shiftExprLevelsUp(q.HavingQual, delta, minDepth)
	cp.LimitCount = shiftExprLevelsUp(q.LimitCount, delta, minDepth)
	cp.LimitOffset = shiftExprLevelsUp(q.LimitOffset, delta, minDepth)
	cp.LArg = shiftQueryLevelsUp(q.LArg, delta, minDepth)
	cp.RArg = shiftQueryLevelsUp(q.RArg, delta, minDepth)
	return &cp
}

// analyzeFuncCall resolves a function call expression.
func analyzeFuncCall(c *Catalog, fc *nodes.FuncCallExpr, scope *analyzerScope) (AnalyzedExpr, error) {
	args := make([]AnalyzedExpr, 0, len(fc.Args))
	for _, arg := range fc.Args {
		a, err := analyzeExpr(c, arg, scope)
		if err != nil {
			return nil, err
		}
		args = append(args, a)
	}
	lower := strings.ToLower(fc.Name)
	if lower == "coalesce" || lower == "ifnull" {
		return &CoalesceExprQ{Args: args, IsIfNull: lower == "ifnull"}, nil
	}
	result := &FuncCallExprQ{
		Name: lower,
		Args: args,
		// A spaced gated keyword is a generic call, not the builtin, so it is
		// never the aggregate: `count (a)` is a stored-function call, not COUNT.
		IsAggregate: isAggregateFunc(fc.Name) && !fc.SpacedParen,
		Distinct:    fc.Quantifier.RemovesDuplicates(),
		SpacedParen: fc.SpacedParen,
	}
	// Phase 3: populate return type from function type table.
	result.ResultType = functionReturnType(result.Name, result.Args)
	if fc.JSONValueArm {
		result.JSONValueClauses = deparse.JSONValueClauses(fc.JSONValue, "")
		if fc.JSONValue != nil && fc.JSONValue.Returning != nil {
			result.ResultType = dataTypeToResolvedType(fc.JSONValue.Returning)
		}
	}
	return result, nil
}

// analyzePositionExpr resolves POSITION(substr IN str). The node has to be
// analyzed and not skipped: an unanalyzed expression fails the whole view
// body, createView then stores a nil AnalyzedQuery, and ExpandMergeViews can
// neither expand the view nor carry its column lineage.
//
// mysql:8.0 resolves POSITION(a IN b) to locate(a, b) and stores it that way,
// so that is the call this records; functionReturnType resolves its BIGINT type.
func analyzePositionExpr(c *Catalog, pe *nodes.PositionExpr, scope *analyzerScope) (AnalyzedExpr, error) {
	substr, err := analyzeExpr(c, pe.Substr, scope)
	if err != nil {
		return nil, err
	}
	str, err := analyzeExpr(c, pe.Str, scope)
	if err != nil {
		return nil, err
	}
	result := &FuncCallExprQ{Name: "locate", Args: []AnalyzedExpr{substr, str}}
	result.ResultType = functionReturnType(result.Name, result.Args)
	return result, nil
}

// analyzeBinaryExpr resolves a binary expression.
func analyzeBinaryExpr(c *Catalog, expr *nodes.BinaryExpr, scope *analyzerScope) (AnalyzedExpr, error) {
	if sub, ok := expr.Right.(*nodes.SubqueryExpr); ok && sub.Quantifier != "" {
		return analyzeQuantifiedComparison(c, expr, sub, scope)
	}
	left, err := analyzeExpr(c, expr.Left, scope)
	if err != nil {
		return nil, err
	}
	right, err := analyzeExpr(c, expr.Right, scope)
	if err != nil {
		return nil, err
	}

	switch expr.Op {
	case nodes.BinOpAnd:
		return &BoolExprQ{Op: BoolAnd, Args: []AnalyzedExpr{left, right}}, nil
	case nodes.BinOpOr:
		return &BoolExprQ{Op: BoolOr, Args: []AnalyzedExpr{left, right}}, nil
	default:
		return &OpExprQ{Op: binaryOpToString(expr.Op), Left: left, Right: right}, nil
	}
}

// analyzeQuantifiedComparison resolves `x op {ANY|SOME|ALL} (SELECT ...)`. The
// parser keeps the quantifier on the subquery; read without it, the comparison
// is one against a scalar subquery, and the SQL deparsed from that fails as soon
// as the subquery returns a second row. SOME is a synonym for ANY.
func analyzeQuantifiedComparison(c *Catalog, expr *nodes.BinaryExpr, sub *nodes.SubqueryExpr, scope *analyzerScope) (AnalyzedExpr, error) {
	var kind SubLinkKind
	switch strings.ToUpper(sub.Quantifier) {
	case "ANY", "SOME":
		kind = SubLinkAny
	case "ALL":
		kind = SubLinkAll
	default:
		return nil, &Error{
			Code:     0,
			SQLState: "HY000",
			Message:  "unsupported subquery quantifier in analyzer: " + sub.Quantifier,
		}
	}
	left, err := analyzeExpr(c, expr.Left, scope)
	if err != nil {
		return nil, err
	}
	innerQ, err := c.analyzeSelectStmtIn(sub.Select, scope, scope.ctes, scope.reachableCTEs, scope.source)
	if err != nil {
		return nil, err
	}
	return &SubLinkExprQ{
		Kind:     kind,
		TestExpr: left,
		Op:       binaryOpToString(expr.Op),
		Subquery: innerQ,
	}, nil
}

// analyzeUnaryExpr resolves a unary expression.
func analyzeUnaryExpr(c *Catalog, expr *nodes.UnaryExpr, scope *analyzerScope) (AnalyzedExpr, error) {
	operand, err := analyzeExpr(c, expr.Operand, scope)
	if err != nil {
		return nil, err
	}

	switch expr.Op {
	case nodes.UnaryPlus:
		return operand, nil
	case nodes.UnaryMinus:
		return &OpExprQ{Op: "-", Right: operand}, nil
	case nodes.UnaryNot:
		return &BoolExprQ{Op: BoolNot, Args: []AnalyzedExpr{operand}}, nil
	case nodes.UnaryBitNot:
		return &OpExprQ{Op: "~", Right: operand}, nil
	case nodes.UnaryBinary:
		return &OpExprQ{Op: "BINARY", Right: operand}, nil
	default:
		return nil, &Error{
			Code:     0,
			SQLState: "HY000",
			Message:  "unsupported unary operator in analyzer",
		}
	}
}

// analyzeInExpr resolves an IN expression.
func analyzeInExpr(c *Catalog, expr *nodes.InExpr, scope *analyzerScope) (AnalyzedExpr, error) {
	arg, err := analyzeExpr(c, expr.Expr, scope)
	if err != nil {
		return nil, err
	}

	if expr.Select != nil {
		// IN (SELECT ...) / NOT IN (SELECT ...)
		innerQ, err := c.analyzeSelectStmtIn(expr.Select, scope, scope.ctes, scope.reachableCTEs, scope.source)
		if err != nil {
			return nil, err
		}
		subLink := &SubLinkExprQ{
			Kind:     SubLinkIn,
			TestExpr: arg,
			Op:       "=",
			Subquery: innerQ,
		}
		if expr.Not {
			return &BoolExprQ{Op: BoolNot, Args: []AnalyzedExpr{subLink}}, nil
		}
		return subLink, nil
	}

	list := make([]AnalyzedExpr, 0, len(expr.List))
	for _, item := range expr.List {
		a, err := analyzeExpr(c, item, scope)
		if err != nil {
			return nil, err
		}
		list = append(list, a)
	}

	return &InListExprQ{Arg: arg, List: list, Negated: expr.Not}, nil
}

// analyzeBetweenExpr resolves a BETWEEN expression.
func analyzeBetweenExpr(c *Catalog, expr *nodes.BetweenExpr, scope *analyzerScope) (AnalyzedExpr, error) {
	arg, err := analyzeExpr(c, expr.Expr, scope)
	if err != nil {
		return nil, err
	}
	lower, err := analyzeExpr(c, expr.Low, scope)
	if err != nil {
		return nil, err
	}
	upper, err := analyzeExpr(c, expr.High, scope)
	if err != nil {
		return nil, err
	}

	return &BetweenExprQ{Arg: arg, Lower: lower, Upper: upper, Negated: expr.Not}, nil
}

// analyzeIsExpr resolves an IS expression.
func analyzeIsExpr(c *Catalog, expr *nodes.IsExpr, scope *analyzerScope) (AnalyzedExpr, error) {
	arg, err := analyzeExpr(c, expr.Expr, scope)
	if err != nil {
		return nil, err
	}

	switch expr.Test {
	case nodes.IsNull:
		return &NullTestExprQ{Arg: arg, IsNull: !expr.Not}, nil
	case nodes.IsTrue:
		op := "IS TRUE"
		if expr.Not {
			op = "IS NOT TRUE"
		}
		return &OpExprQ{Op: op, Left: arg}, nil
	case nodes.IsFalse:
		op := "IS FALSE"
		if expr.Not {
			op = "IS NOT FALSE"
		}
		return &OpExprQ{Op: op, Left: arg}, nil
	case nodes.IsUnknown:
		op := "IS UNKNOWN"
		if expr.Not {
			op = "IS NOT UNKNOWN"
		}
		return &OpExprQ{Op: op, Left: arg}, nil
	default:
		return nil, &Error{
			Code:     0,
			SQLState: "HY000",
			Message:  "unsupported IS test type in analyzer",
		}
	}
}

// analyzeCaseExpr resolves a CASE expression.
func analyzeCaseExpr(c *Catalog, expr *nodes.CaseExpr, scope *analyzerScope) (AnalyzedExpr, error) {
	var testExpr AnalyzedExpr
	if expr.Operand != nil {
		var err error
		testExpr, err = analyzeExpr(c, expr.Operand, scope)
		if err != nil {
			return nil, err
		}
	}

	whens := make([]*CaseWhenQ, 0, len(expr.Whens))
	for _, w := range expr.Whens {
		cond, err := analyzeExpr(c, w.Cond, scope)
		if err != nil {
			return nil, err
		}
		then, err := analyzeExpr(c, w.Result, scope)
		if err != nil {
			return nil, err
		}
		whens = append(whens, &CaseWhenQ{Cond: cond, Then: then})
	}

	var def AnalyzedExpr
	if expr.Default != nil {
		var err error
		def, err = analyzeExpr(c, expr.Default, scope)
		if err != nil {
			return nil, err
		}
	}

	return &CaseExprQ{TestExpr: testExpr, Args: whens, Default: def}, nil
}

// analyzeCastExpr resolves a CAST expression.
func analyzeCastExpr(c *Catalog, expr *nodes.CastExpr, scope *analyzerScope) (AnalyzedExpr, error) {
	arg, err := analyzeExpr(c, expr.Expr, scope)
	if err != nil {
		return nil, err
	}

	var targetType *ResolvedType
	if expr.TypeName != nil {
		targetType = dataTypeToResolvedType(expr.TypeName)
	}

	return &CastExprQ{
		Arg:                arg,
		TargetType:         targetType,
		AtTimeZone:         expr.AtTimeZone,
		AtTimeZoneInterval: expr.AtTimeZoneInterval,
	}, nil
}

// dataTypeToResolvedType maps a parser DataType to a ResolvedType for CAST targets.
func dataTypeToResolvedType(dt *nodes.DataType) *ResolvedType {
	name := strings.ToLower(dt.Name)
	switch name {
	case "signed", "signed integer":
		return &ResolvedType{BaseType: BaseTypeBigInt}
	case "unsigned", "unsigned integer":
		return &ResolvedType{BaseType: BaseTypeBigInt, Unsigned: true}
	case "char":
		rt := &ResolvedType{BaseType: BaseTypeChar}
		if dt.Length > 0 {
			rt.Length = dt.Length
		}
		return rt
	case "binary":
		rt := &ResolvedType{BaseType: BaseTypeBinary}
		if dt.Length > 0 {
			rt.Length = dt.Length
		}
		return rt
	case "decimal":
		return &ResolvedType{BaseType: BaseTypeDecimal, Precision: dt.Length, Scale: dt.Scale}
	case "date":
		return &ResolvedType{BaseType: BaseTypeDate}
	case "datetime":
		return &ResolvedType{BaseType: BaseTypeDateTime}
	case "time":
		return &ResolvedType{BaseType: BaseTypeTime}
	case "json":
		return &ResolvedType{BaseType: BaseTypeJSON}
	case "float":
		return &ResolvedType{BaseType: BaseTypeFloat}
	case "double":
		return &ResolvedType{BaseType: BaseTypeDouble}
	default:
		return &ResolvedType{BaseType: BaseTypeUnknown}
	}
}

// analyzeScalarSubquery resolves a scalar subquery expression: (SELECT ...).
func analyzeScalarSubquery(c *Catalog, subq *nodes.SubqueryExpr, scope *analyzerScope) (AnalyzedExpr, error) {
	innerQ, err := c.analyzeSelectStmtIn(subq.Select, scope, scope.ctes, scope.reachableCTEs, scope.source)
	if err != nil {
		return nil, err
	}
	return &SubLinkExprQ{
		Kind:     SubLinkScalar,
		Subquery: innerQ,
	}, nil
}

// analyzeExistsSubquery resolves an EXISTS (SELECT ...) expression.
func analyzeExistsSubquery(c *Catalog, expr *nodes.ExistsExpr, scope *analyzerScope) (AnalyzedExpr, error) {
	innerQ, err := c.analyzeSelectStmtIn(expr.Select, scope, scope.ctes, scope.reachableCTEs, scope.source)
	if err != nil {
		return nil, err
	}
	return &SubLinkExprQ{
		Kind:     SubLinkExists,
		Subquery: innerQ,
	}, nil
}

// isAggregateFunc returns true if the function name is a known aggregate.
func isAggregateFunc(name string) bool {
	switch strings.ToLower(name) {
	case "count", "sum", "avg", "min", "max",
		"group_concat", "json_arrayagg", "json_objectagg",
		"bit_and", "bit_or", "bit_xor",
		"std", "stddev", "stddev_pop", "stddev_samp",
		"var_pop", "var_samp", "variance",
		"any_value":
		return true
	}
	return false
}
