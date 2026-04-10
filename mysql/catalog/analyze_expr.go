package catalog

import (
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// analyzeExpr is the main expression analysis dispatcher.
func analyzeExpr(c *Catalog, expr nodes.ExprNode, scope *analyzerScope) (AnalyzedExpr, error) {
	switch n := expr.(type) {
	case *nodes.ColumnRef:
		return analyzeColumnRef(n, scope)
	case *nodes.IntLit:
		return &ConstExprQ{Value: strconv.FormatInt(n.Value, 10)}, nil
	case *nodes.StringLit:
		return &ConstExprQ{Value: n.Value}, nil
	case *nodes.FloatLit:
		return &ConstExprQ{Value: n.Value}, nil
	case *nodes.NullLit:
		return &ConstExprQ{IsNull: true, Value: "NULL"}, nil
	case *nodes.BoolLit:
		if n.Value {
			return &ConstExprQ{Value: "TRUE"}, nil
		}
		return &ConstExprQ{Value: "FALSE"}, nil
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
	case *nodes.IsExpr:
		return analyzeIsExpr(c, n, scope)
	case *nodes.CaseExpr:
		return analyzeCaseExpr(c, n, scope)
	case *nodes.CastExpr:
		return analyzeCastExpr(c, n, scope)
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

// analyzeColumnRef resolves a column reference against the scope.
// Uses the Full resolution methods to support correlated subquery references
// (setting LevelsUp > 0 when the column comes from a parent scope).
func analyzeColumnRef(ref *nodes.ColumnRef, scope *analyzerScope) (AnalyzedExpr, error) {
	if ref.Table != "" {
		rteIdx, attNum, levelsUp, err := scope.resolveQualifiedColumnFull(ref.Table, ref.Column)
		if err != nil {
			return nil, err
		}
		return &VarExprQ{RangeIdx: rteIdx, AttNum: attNum, LevelsUp: levelsUp}, nil
	}
	rteIdx, attNum, levelsUp, err := scope.resolveColumnFull(ref.Column)
	if err != nil {
		return nil, err
	}
	return &VarExprQ{RangeIdx: rteIdx, AttNum: attNum, LevelsUp: levelsUp}, nil
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
		return &CoalesceExprQ{Args: args}, nil
	}
	return &FuncCallExprQ{
		Name:        lower,
		Args:        args,
		IsAggregate: isAggregateFunc(fc.Name),
		Distinct:    fc.Distinct,
	}, nil
}

// analyzeBinaryExpr resolves a binary expression.
func analyzeBinaryExpr(c *Catalog, expr *nodes.BinaryExpr, scope *analyzerScope) (AnalyzedExpr, error) {
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
		innerQ, err := c.analyzeSelectStmtInternal(expr.Select, scope)
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

	return &CastExprQ{Arg: arg, TargetType: targetType}, nil
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
	innerQ, err := c.analyzeSelectStmtInternal(subq.Select, scope)
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
	innerQ, err := c.analyzeSelectStmtInternal(expr.Select, scope)
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
