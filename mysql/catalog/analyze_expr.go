package catalog

import (
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// analyzeExpr is the main expression analysis dispatcher.
func analyzeExpr(expr nodes.ExprNode, scope *analyzerScope) (AnalyzedExpr, error) {
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
		return analyzeFuncCall(n, scope)
	case *nodes.ParenExpr:
		return analyzeExpr(n.Expr, scope)
	default:
		return nil, &Error{
			Code:     0,
			SQLState: "HY000",
			Message:  "unsupported expression type in analyzer",
		}
	}
}

// analyzeColumnRef resolves a column reference against the scope.
func analyzeColumnRef(ref *nodes.ColumnRef, scope *analyzerScope) (AnalyzedExpr, error) {
	if ref.Table != "" {
		rteIdx, attNum, err := scope.resolveQualifiedColumn(ref.Table, ref.Column)
		if err != nil {
			return nil, err
		}
		return &VarExprQ{RangeIdx: rteIdx, AttNum: attNum}, nil
	}
	rteIdx, attNum, err := scope.resolveColumn(ref.Column)
	if err != nil {
		return nil, err
	}
	return &VarExprQ{RangeIdx: rteIdx, AttNum: attNum}, nil
}

// analyzeFuncCall resolves a function call expression.
func analyzeFuncCall(fc *nodes.FuncCallExpr, scope *analyzerScope) (AnalyzedExpr, error) {
	args := make([]AnalyzedExpr, 0, len(fc.Args))
	for _, arg := range fc.Args {
		a, err := analyzeExpr(arg, scope)
		if err != nil {
			return nil, err
		}
		args = append(args, a)
	}
	return &FuncCallExprQ{
		Name:        strings.ToLower(fc.Name),
		Args:        args,
		IsAggregate: isAggregateFunc(fc.Name),
		Distinct:    fc.Distinct,
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
