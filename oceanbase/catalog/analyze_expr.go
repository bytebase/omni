package catalog

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode/utf8"

	nodes "github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/oceanbase/deparse"
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
	case *nodes.TemporalLit:
		var rt *ResolvedType
		switch strings.ToUpper(n.Type) {
		case "DATE":
			rt = &ResolvedType{BaseType: BaseTypeDate}
		case "TIME":
			rt = &ResolvedType{BaseType: BaseTypeTime}
		case "TIMESTAMP":
			rt = &ResolvedType{BaseType: BaseTypeTimestamp}
		}
		return &RawExprQ{Text: deparse.Deparse(n), ResultType: rt}, nil
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
	case *nodes.QuantifiedExpr:
		// any_expr's `ALL (expr)` / `ANY (expr)` keeps its quantifier around the
		// analyzed expression so `a = ANY(b)` does not collapse to `a = b`.
		arg, err := analyzeExpr(c, n.Expr, scope)
		if err != nil {
			return nil, err
		}
		return &QuantifiedExprQ{Quantifier: strings.ToUpper(n.Quantifier), Arg: arg}, nil
	case *nodes.OdbcEscapeExpr:
		// The escape's word says how its expression is read (`{d '2026-01-01'}`
		// is a date), and its expression is analyzed like any other, so a column
		// inside it keeps its lineage.
		arg, err := analyzeExpr(c, n.Expr, scope)
		if err != nil {
			return nil, err
		}
		return &OdbcEscapeExprQ{Name: n.Name, Arg: arg}, nil
	case *nodes.RemoteSequenceExpr:
		// A sequence read through a database link is a value of its own, which
		// the IR keeps as written: it is no column of a local table, and it is
		// no literal either.
		return &RawExprQ{Text: deparse.Deparse(n)}, nil
	case *nodes.LambdaExpr:
		body, err := analyzeExpr(c, n.Body, scope.withMaskedColumns(n.Params))
		if err != nil {
			return nil, err
		}
		return &LambdaExprQ{Params: append([]string(nil), n.Params...), Body: body}, nil
	case *nodes.ResTarget:
		val, err := analyzeExpr(c, n.Val, scope)
		if err != nil {
			return nil, err
		}
		if n.Name == "" {
			return val, nil
		}
		return &ResTargetExprQ{Name: n.Name, Val: val}, nil
	case *nodes.GroupingSet:
		// A GROUP BY item's grouping set keeps its kind and its items, so a view
		// stores ROLLUP (...) and not the expressions it groups by.
		items := make([]AnalyzedExpr, 0, len(n.Items))
		for _, item := range n.Items {
			analyzed, err := analyzeExpr(c, item, scope)
			if err != nil {
				return nil, err
			}
			items = append(items, analyzed)
		}
		return &GroupingSetExprQ{Kind: n.Kind, Args: items}, nil
	case *nodes.RowExpr:
		items := make([]AnalyzedExpr, 0, len(n.Items))
		for _, item := range n.Items {
			analyzed, err := analyzeExpr(c, item, scope)
			if err != nil {
				return nil, err
			}
			items = append(items, analyzed)
		}
		return &RowExprQ{Args: items}, nil
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
	case *nodes.MemberOfExpr:
		val, err := analyzeExpr(c, n.Value, scope)
		if err != nil {
			return nil, err
		}
		arr, err := analyzeExpr(c, n.Array, scope)
		if err != nil {
			return nil, err
		}
		return &RawExprQ{Text: deparse.Deparse(n), Args: []AnalyzedExpr{val, arr}}, nil
	case *nodes.MatchExpr:
		var args []AnalyzedExpr
		for _, col := range n.Columns {
			a, err := analyzeColumnRef(col, scope)
			if err != nil {
				return nil, err
			}
			args = append(args, a)
		}
		if n.Against != nil {
			a, err := analyzeExpr(c, n.Against, scope)
			if err != nil {
				return nil, err
			}
			args = append(args, a)
		}
		return &RawExprQ{Text: deparse.Deparse(n), Args: args}, nil
	case *nodes.ConvertExpr:
		var args []AnalyzedExpr
		if n.Expr != nil {
			a, err := analyzeExpr(c, n.Expr, scope)
			if err != nil {
				return nil, err
			}
			args = append(args, a)
		}
		if err := checkQuotedCharset(n.Charset, n.CharsetQuoted); err != nil {
			return nil, err
		}
		return &RawExprQ{Text: deparse.Deparse(n), Args: args}, nil
	case *nodes.ExtractExpr:
		var args []AnalyzedExpr
		if n.Expr != nil {
			a, err := analyzeExpr(c, n.Expr, scope)
			if err != nil {
				return nil, err
			}
			args = append(args, a)
		}
		return &RawExprQ{Text: deparse.Deparse(n), Args: args}, nil
	case *nodes.DefaultExpr:
		var args []AnalyzedExpr
		if n.Column != "" {
			a, err := analyzeColumnRef(&nodes.ColumnRef{Column: n.Column}, scope)
			if err != nil {
				return nil, err
			}
			args = append(args, a)
		}
		return &RawExprQ{Text: deparse.Deparse(n), Args: args}, nil
	case *nodes.WeightStringExpr:
		var args []AnalyzedExpr
		if n.Expr != nil {
			a, err := analyzeExpr(c, n.Expr, scope)
			if err != nil {
				return nil, err
			}
			args = append(args, a)
		}
		return &RawExprQ{Text: deparse.Deparse(n), Args: args}, nil
	case *nodes.VariableRef:
		return &RawExprQ{Text: deparse.Deparse(n)}, nil
	case *nodes.HexLit:
		text := n.Value
		if !strings.HasPrefix(text, "0x") && !strings.HasPrefix(text, "0X") {
			text = "X'" + text + "'"
		}
		return &RawExprQ{Text: text, ResultType: &ResolvedType{BaseType: BaseTypeBinary}}, nil
	case *nodes.BitLit:
		text := n.Value
		if !strings.HasPrefix(text, "0b") && !strings.HasPrefix(text, "0B") {
			text = "b'" + text + "'"
		}
		return &RawExprQ{Text: text, ResultType: &ResolvedType{BaseType: BaseTypeBinary}}, nil
	default:
		return nil, &Error{
			Code:     0,
			SQLState: "HY000",
			Message:  "unsupported expression type in analyzer",
		}
	}
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
		return scope.resolveQualifiedRefFull(ref.Catalog, ref.Schema, ref.Table, ref.Column)
	}
	if maskLevels, masked := scope.lambdaMaskLevels(ref.Column); masked {
		if maskLevels == 0 {
			return &LambdaParamExprQ{Name: ref.Column}, nil
		}
		if rteIdx, attNum, levelsUp, found, err := scope.resolveColumnBeforeLevel(ref.Column, maskLevels); err != nil {
			return nil, err
		} else if found {
			if levelsUp == 0 {
				if merged, ok := scope.fullJoinColumn(ref.Column); ok {
					return merged, nil
				}
			} else if merged, ok := scope.rebaseFullJoinColumn(ref.Column, levelsUp, rteIdx); ok {
				return merged, nil
			}
			return &VarExprQ{RangeIdx: rteIdx, AttNum: attNum, LevelsUp: levelsUp}, nil
		}
		aliasTE, aliasLevels := scope.selectAlias(ref.Column)
		if aliasTE != nil && aliasLevels == 0 {
			return aliasTE.Expr, nil
		}
		if idx, unkLevels, ok := scope.soleUnknownColumnInChain(); ok && unkLevels < maskLevels {
			return &VarExprQ{RangeIdx: idx, Name: ref.Column, LevelsUp: unkLevels}, nil
		}
		if aliasTE != nil && aliasLevels < maskLevels {
			return nil, errOuterAlias
		}
		if scope.uncertainBeforeLevel(maskLevels) {
			return nil, scope.unsure(errNoSuchColumn(ref.Column, "field list"))
		}
		return &LambdaParamExprQ{Name: ref.Column}, nil
	}
	rteIdx, attNum, levelsUp, err := scope.resolveColumnFull(ref.Column)
	if err != nil {
		// Check same-scope SELECT aliases first in GROUP BY / HAVING / ORDER BY
		// before attributing an unresolved name to an unknown remote relation,
		// and compare nesting distance when an enclosing scope holds either.
		aliasTE, aliasLevels := scope.selectAlias(ref.Column)
		if aliasTE != nil && aliasLevels == 0 {
			return aliasTE.Expr, nil
		}
		if idx, unkLevels, ok := scope.soleUnknownColumnInChain(); ok {
			if aliasTE == nil || unkLevels <= aliasLevels {
				return &VarExprQ{RangeIdx: idx, Name: ref.Column, LevelsUp: unkLevels}, nil
			}
		}
		if aliasTE != nil {
			return nil, errOuterAlias
		}
		return nil, err
	}
	// A FULL join's key read unqualified is COALESCE(left, right): an r-only row
	// would project NULL from the surviving left column alone. A correlated
	// reference to an outer query's key lifts that expression's variables by the
	// levels it crossed.
	if levelsUp == 0 {
		if merged, ok := scope.fullJoinColumn(ref.Column); ok {
			return merged, nil
		}
	} else if merged, ok := scope.rebaseFullJoinColumn(ref.Column, levelsUp, rteIdx); ok {
		return merged, nil
	}
	return &VarExprQ{RangeIdx: rteIdx, AttNum: attNum, LevelsUp: levelsUp}, nil
}

// rebaseAnalyzedLevels returns expr with every variable of the query it belongs
// to moved levels queries out. A merged FULL-join key is a COALESCE over that
// query's own variables, which is what a correlated reference to it needs.
func rebaseAnalyzedLevels(expr AnalyzedExpr, levels int) AnalyzedExpr {
	switch e := expr.(type) {
	case *VarExprQ:
		if e.LevelsUp != 0 {
			return e
		}
		rebased := *e
		rebased.LevelsUp = levels
		return &rebased
	case *CoalesceExprQ:
		args := make([]AnalyzedExpr, len(e.Args))
		for i, arg := range e.Args {
			args[i] = rebaseAnalyzedLevels(arg, levels)
		}
		return &CoalesceExprQ{Args: args}
	}
	return expr
}

// checkDistinctAllowed refuses DISTINCT or UNIQUE in JSON_ARRAYAGG and PERCENTILE_CONT
// with OceanBase's 5798. The grammar takes opt_distinct_or_all in both
// (sql_parser_mysql_mode.y, v4.3.5_CE_BP6), and the resolver refuses it before it
// resolves the arguments (1054) or the window's name (5929).
func checkDistinctAllowed(fc *nodes.FuncCallExpr) error {
	if fc.Schema != "" || !fc.Quantifier.RemovesDuplicates() {
		return nil
	}
	switch strings.ToLower(fc.Name) {
	case "json_arrayagg", "percentile_cont":
		return &Error{Code: ErrOBDistinctNotAllowed, SQLState: "HY000", Message: "DISTINCT not allowed here"}
	}
	return nil
}

// checkQuotedCharset checks a CONVERT or CHAR ... USING character set written as a
// string. charset_name's STRING_VALUE arm makes it a parameter (sql_parser_mysql_mode.y,
// v4.3.5_CE_BP6), so OceanBase reads the argument first: CONVERT(a USING 's') with no
// column a is 1054, where CONVERT(a USING s) is 1115. A name written bare is checked
// before the query is resolved (validateViewQueryPreAnalysis).
func checkQuotedCharset(charset string, quoted bool) error {
	if charset == "" || !quoted {
		return nil
	}
	return checkCharsetName(charset)
}

// analyzeFuncCall resolves a function call expression.
func analyzeFuncCall(c *Catalog, fc *nodes.FuncCallExpr, scope *analyzerScope) (AnalyzedExpr, error) {
	lower := strings.ToLower(fc.Name)
	if fc.Schema == "" {
		// OceanBase 4.3.5.6 rejects aliased arguments in native function calls
		// (e.g. STATISTICS(a zzz)) with 1583 before resolving columns.
		for _, arg := range fc.Args {
			if rt, ok := arg.(*nodes.ResTarget); ok && rt.Name != "" {
				return nil, &Error{
					Code:     1583,
					SQLState: "42000",
					Message:  fmt.Sprintf("Incorrect parameters in the call to native function '%s'", lower),
				}
			}
		}
		// OceanBase 4.3.5.6 validates that NTH_VALUE's second argument is a
		// constant (not a column reference) before resolving window names (5929).
		if lower == "nth_value" && len(fc.Args) == 2 && exprContainsColumnRef(fc.Args[1]) {
			return nil, &Error{
				Code:     1210,
				SQLState: "HY000",
				Message:  "Incorrect arguments to nth_value",
			}
		}
		// OceanBase 4.3.5.6 validates CALC_PARTITION_ID's first argument as a
		// constant before resolving column references (1054).
		if lower == "calc_partition_id" && (len(fc.Args) < 2 || len(fc.Args) > 3 || exprContainsColumnRef(fc.Args[0])) {
			return nil, &Error{
				Code:     ErrWrongArguments,
				SQLState: sqlState(ErrWrongArguments),
				Message:  "Invalid argument",
			}
		}
		// Zero-argument calls to built-in functions that always require at least
		// one argument fail with ER_WRONG_PARAMCOUNT_TO_NATIVE_FCT (1582).
		if len(fc.Args) == 0 && builtinRequiresArgs(lower) {
			return nil, errWrongParamCount(lower)
		}
		if err := checkDistinctAllowed(fc); err != nil {
			return nil, err
		}
	}
	args := make([]AnalyzedExpr, 0, len(fc.Args))
	for _, arg := range fc.Args {
		a, err := analyzeExpr(c, arg, scope)
		if err != nil {
			return nil, err
		}
		args = append(args, a)
	}
	if err := checkQuotedCharset(fc.Using, fc.UsingQuoted); err != nil {
		return nil, err
	}
	if lower == "coalesce" || lower == "ifnull" {
		return &CoalesceExprQ{Args: args, IsIfNull: lower == "ifnull"}, nil
	}
	result := &FuncCallExprQ{
		Name:                lower,
		Args:                args,
		IsAggregate:         isAggregateFunc(fc.Name),
		Distinct:            fc.Quantifier.RemovesDuplicates(),
		FromFirstLast:       fc.FromFirstLast,
		NullTreatment:       fc.NullTreatment,
		NullTreatmentInArgs: fc.NullTreatmentInArgs,
	}
	// Phase 3: populate return type from function type table.
	result.ResultType = functionReturnType(result.Name, result.Args)
	if fc.JSONValueArm {
		result.JSONValueClauses = deparse.JSONValueClauses(fc.JSONValue)
		if fc.JSONValue != nil && fc.JSONValue.Returning != nil {
			result.ResultType = dataTypeToResolvedType(fc.JSONValue.Returning)
		}
	}
	if fc.JSONQueryArm {
		result.JSONQueryClauses = deparse.JSONQueryClauses(fc.JSONQuery)
		if fc.JSONQuery != nil && fc.JSONQuery.Returning != nil {
			result.ResultType = dataTypeToResolvedType(fc.JSONQuery.Returning)
		}
	}
	if fc.Over != nil {
		over, err := analyzeWindowDef(c, fc.Over, scope, false)
		if err != nil {
			return nil, err
		}
		result.Over = over
	}
	return result, nil
}

func exprContainsColumnRef(expr nodes.ExprNode) bool {
	found := false
	nodes.Inspect(expr, func(n nodes.Node) bool {
		if found {
			return false
		}
		if _, ok := n.(*nodes.ColumnRef); ok {
			found = true
			return false
		}
		return true
	})
	return found
}

// BuiltinRequiresArgs reports whether a builtin function of the lower-cased
// name fails a call with no argument (1582).
func BuiltinRequiresArgs(lower string) bool { return builtinRequiresArgs(lower) }

func builtinRequiresArgs(lower string) bool {
	switch lower {
	case "coalesce", "charset", "collation", "dump",
		"repeat", "reverse", "right", "replace",
		"truncate", "format", "normal", "password", "encrypt":
		return true
	}
	return false
}

func analyzeWindowDef(c *Catalog, wd *nodes.WindowDef, scope *analyzerScope, isDecl bool) (*WindowDefQ, error) {
	if wd == nil {
		return nil, nil
	}
	out := &WindowDefQ{
		Name:    wd.Name,
		RefName: wd.RefName,
	}
	if !isDecl && out.Name == "" && out.RefName != "" && len(wd.PartitionBy) == 0 && len(wd.OrderBy) == 0 && wd.Frame == nil {
		out.Name = out.RefName
		out.RefName = ""
	}
	for _, expr := range wd.PartitionBy {
		a, err := analyzeExpr(c, expr, scope)
		if err != nil {
			return nil, err
		}
		out.PartitionBy = append(out.PartitionBy, a)
	}
	for _, item := range wd.OrderBy {
		a, err := analyzeExpr(c, item.Expr, scope)
		if err != nil {
			return nil, err
		}
		out.OrderExprs = append(out.OrderExprs, &WindowOrderItemQ{
			Expr:       a,
			Descending: item.Desc,
		})
	}
	if wd.Frame != nil {
		if wd.Frame.Start != nil && wd.Frame.Start.Offset != nil {
			if _, err := analyzeExpr(c, wd.Frame.Start.Offset, scope); err != nil {
				return nil, err
			}
		}
		if wd.Frame.End != nil && wd.Frame.End.Offset != nil {
			if _, err := analyzeExpr(c, wd.Frame.End.Offset, scope); err != nil {
				return nil, err
			}
		}
		out.FrameClause = deparse.DeparseWindowFrame(wd.Frame)
	}
	return out, nil
}

// analyzePositionExpr resolves POSITION(substr IN str). The node has to be
// analyzed and not skipped: an unanalyzed expression fails the whole view
// body, createView then stores a nil AnalyzedQuery, and ExpandMergeViews can
// neither expand the view nor carry its column lineage.
//
// MySQL 8.0, whose stored form this catalog follows, resolves
// POSITION(a IN b) to locate(a, b), so that is the call this records;
// functionReturnType resolves its BIGINT type.
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
	innerQ, err := c.analyzeSelectStmtInternal(sub.Select, scope)
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

	return &CastExprQ{Arg: arg, TargetType: targetType, Ignore: expr.Ignore}, nil
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
