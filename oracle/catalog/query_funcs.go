package catalog

import (
	"github.com/bytebase/omni/oracle/ast"
)

// funcCall types a function call: a built-in function, an aggregate or
// analytic function, or a function a WITH clause declares.
func (a *queryAnalyzer) funcCall(n *ast.FuncCallExpr, sc *qscope, clause exprClause) (qtype, bool, error) {
	if n.FuncName == nil {
		return simpleType("UNKNOWN"), false, nil
	}
	name := n.FuncName.Name
	if n.FuncName.Schema == "" {
		switch name {
		case "":
			// A row constructor: (a, b) IN (...).
			for _, item := range listItems(n.Args) {
				if ex, ok := item.(ast.ExprNode); ok {
					if _, _, err := a.expr(ex, sc, clause); err != nil {
						return qtype{}, false, err
					}
				}
			}
			return simpleType("UNKNOWN"), false, nil
		case "ALL", "ANY", "SOME":
			// A quantified comparison's right side: a subquery or a list.
			var t qtype
			for i, item := range listItems(n.Args) {
				ex, ok := item.(ast.ExprNode)
				if !ok {
					continue
				}
				if sq, ok := ex.(*ast.SubqueryExpr); ok {
					if err := a.subquery(sq, sc, -1); err != nil {
						return qtype{}, false, err
					}
					t = simpleType("UNKNOWN")
					continue
				}
				xt, _, err := a.expr(ex, sc, clause)
				if err != nil {
					return qtype{}, false, err
				}
				if i == 0 {
					t = xt
				}
			}
			return t, false, nil
		case "XMLELEMENT":
			// XMLELEMENT's first argument names the element.
			for i, item := range listItems(n.Args) {
				ex, ok := item.(ast.ExprNode)
				if !ok || i == 0 {
					continue
				}
				if _, _, err := a.expr(ex, sc, clause); err != nil {
					return qtype{}, false, err
				}
			}
			return simpleType("XMLTYPE"), false, nil
		}
	}
	aggregate := isAggregateCall(n)
	if aggregate && !clause.allowsAggregate() {
		return qtype{}, false, oraError("00934", "group function is not allowed here")
	}
	if n.Over != nil && !clause.allowsWindow() {
		return qtype{}, false, oraError("30483", "window  functions are not allowed here")
	}
	argClause := clause
	if aggregate || n.Over != nil {
		// An aggregate's argument may hold another aggregate (nested
		// aggregation), and an analytic function's an aggregate.
		argClause = clauseSelect
	}
	var args []qtype
	var argNN []bool
	for _, item := range listItems(n.Args) {
		ex, ok := item.(ast.ExprNode)
		if !ok {
			continue
		}
		if na, ok := ex.(*ast.NamedArgExpr); ok {
			ex = na.Expr
		}
		t, nn, err := a.expr(ex, sc, argClause)
		if err != nil {
			return qtype{}, false, err
		}
		args = append(args, t)
		argNN = append(argNN, nn)
	}
	for _, item := range listItems(n.OrderBy) {
		if sb, ok := item.(*ast.SortBy); ok {
			if _, _, err := a.expr(sb.Expr, sc, clauseSelect); err != nil {
				return qtype{}, false, err
			}
		}
	}
	if n.KeepClause != nil {
		for _, item := range listItems(n.KeepClause.OrderBy) {
			if sb, ok := item.(*ast.SortBy); ok {
				if _, _, err := a.expr(sb.Expr, sc, clauseSelect); err != nil {
					return qtype{}, false, err
				}
			}
		}
	}
	if n.Over != nil {
		if err := a.windowSpec(n.Over, sc, map[string]bool{}); err != nil {
			return qtype{}, false, err
		}
	}
	arg := func(i int) qtype {
		if i < len(args) {
			return args[i]
		}
		return nullType()
	}
	if n.FuncName.Schema == "" {
		if t, ok := a.plsqlFuncs[name]; ok {
			return t, false, nil
		}
	}
	if aggregate || n.Over != nil {
		return a.aggregateType(n, name, args)
	}
	if n.FuncName.Schema != "" {
		if t, ok := collectionElementType(n); ok {
			_ = t
			return simpleType(n.FuncName.Schema + "." + name), false, nil
		}
		return simpleType("UNKNOWN"), false, nil
	}
	switch name {
	case "UPPER", "LOWER", "INITCAP", "NLS_UPPER", "NLS_LOWER", "NLS_INITCAP":
		return arg(0), false, nil
	case "TRIM", "LTRIM", "RTRIM":
		t := arg(0)
		if t.kind == "CHAR" {
			t.kind = "VARCHAR2"
		}
		if t.isNumeric() {
			t = varchar2Bytes(40)
		}
		return t, false, nil
	case "SUBSTR", "SUBSTRB", "SUBSTRC", "SUBSTR2", "SUBSTR4":
		t := arg(0)
		if t.isLOB() {
			return t, false, nil
		}
		if t.isNumeric() {
			t = varchar2Bytes(40)
		}
		if len(args) > 2 {
			if lit, ok := literalInt(argNode(n, 2)); ok {
				if t.isNational() {
					return characterType("NVARCHAR2", lit, false), false, nil
				}
				return varchar2Bytes(lit * dbCharBytes), false, nil
			}
		}
		if t.kind == "CHAR" {
			t.kind = "VARCHAR2"
		}
		return t, false, nil
	case "LPAD", "RPAD":
		t := arg(0)
		if len(args) > 1 {
			if lit, ok := literalInt(argNode(n, 1)); ok {
				if t.isNational() {
					return characterType("NVARCHAR2", lit, false), false, nil
				}
				return varchar2Bytes(lit * dbCharBytes), false, nil
			}
		}
		return varchar2Bytes(4000), false, nil
	case "REPLACE":
		t := arg(0)
		if t.isLOB() {
			return t, false, nil
		}
		if t.kind == "CHAR" {
			t.kind = "VARCHAR2"
		}
		return t, false, nil
	case "TRANSLATE":
		t := arg(0)
		if t.isLOB() {
			return t, false, nil
		}
		return varchar2Bytes(t.byteLength() * dbCharBytes), false, nil
	case "REGEXP_SUBSTR":
		t := arg(0)
		if t.isLOB() {
			return t, false, nil
		}
		return varchar2Bytes(t.byteLength() * dbCharBytes), false, nil
	case "REGEXP_REPLACE":
		t := arg(0)
		if t.isLOB() {
			return t, false, nil
		}
		return varchar2Bytes(4000), false, nil
	case "LENGTH", "LENGTHB", "LENGTHC", "INSTR", "INSTRB", "INSTRC", "ASCII", "VSIZE",
		"REGEXP_INSTR", "REGEXP_COUNT", "MOD", "REMAINDER", "ROUND", "TRUNC", "ABS", "CEIL", "FLOOR",
		"POWER", "SQRT", "SIGN", "LN", "EXP", "LOG", "SIN", "COS", "TAN", "ASIN", "ACOS", "ATAN", "ATAN2",
		"SINH", "COSH", "TANH", "BITAND", "WIDTH_BUCKET", "MONTHS_BETWEEN", "ORA_HASH", "NANVL":
		switch name {
		case "ROUND", "TRUNC":
			if arg(0).isDatetime() {
				return simpleType("DATE"), false, nil
			}
		}
		if t := arg(0); t.kind == "BINARY_DOUBLE" || t.kind == "BINARY_FLOAT" {
			switch name {
			case "ABS", "CEIL", "FLOOR", "ROUND", "TRUNC", "SQRT", "EXP", "LN", "POWER", "MOD", "SIGN":
				return t, false, nil
			}
		}
		return numberType(), false, nil
	case "CHR":
		return varchar2Bytes(1), false, nil
	case "NCHR":
		return characterType("NVARCHAR2", 1, false), false, nil
	case "DUMP":
		return varchar2Bytes(4000), false, nil
	case "CONCAT":
		return concatType(arg(0), arg(1)), false, nil
	case "TO_CHAR":
		t := arg(0)
		if t.isCharacter() || t.isLOB() {
			if t.kind == "NCLOB" || t.isNational() {
				return varchar2Bytes(t.byteLength()), false, nil
			}
			return t, false, nil
		}
		if len(args) > 1 {
			if s, ok := literalString(argNode(n, 1)); ok {
				if t.isNumeric() {
					return varchar2Bytes(len(s) + 1), false, nil
				}
				return varchar2Bytes(len(s)), false, nil
			}
		}
		if t.kind == "DATE" {
			return varchar2Bytes(75), false, nil
		}
		return varchar2Bytes(40), false, nil
	case "TO_NUMBER":
		return numberType(), false, nil
	case "TO_BINARY_DOUBLE":
		return simpleType("BINARY_DOUBLE"), false, nil
	case "TO_BINARY_FLOAT":
		return simpleType("BINARY_FLOAT"), false, nil
	case "TO_DATE", "ADD_MONTHS", "LAST_DAY", "NEXT_DAY", "NEW_TIME":
		return simpleType("DATE"), false, nil
	case "TO_TIMESTAMP":
		return timestampType("TIMESTAMP", 9), false, nil
	case "TO_TIMESTAMP_TZ":
		return timestampType("TIMESTAMP WITH TIME ZONE", 9), false, nil
	case "TO_CLOB":
		return simpleType("CLOB"), false, nil
	case "TO_NCLOB":
		return simpleType("NCLOB"), false, nil
	case "TO_BLOB":
		return simpleType("BLOB"), false, nil
	case "TO_NCHAR":
		t := arg(0)
		return characterType("NVARCHAR2", t.byteLength()*5/2, false), false, nil
	case "RAWTOHEX":
		return varchar2Bytes(arg(0).bytes * 2), false, nil
	case "HEXTORAW":
		t := simpleType("RAW")
		if s, ok := literalString(argNode(n, 0)); ok {
			t.bytes = (len(s) + 1) / 2
		} else {
			t.bytes = arg(0).byteLength() / 2
		}
		return t, false, nil
	case "SYS_GUID":
		t := simpleType("RAW")
		t.bytes = 16
		return t, false, nil
	case "USERENV":
		return varchar2Bytes(52), false, nil
	case "SYS_CONTEXT":
		return varchar2Bytes(256), false, nil
	case "SYS_CONNECT_BY_PATH":
		return varchar2Bytes(4000), false, nil
	case "NVL", "IFNULL":
		return widenFirst(args), false, nil
	case "NVL2":
		return widenFirst(args[min(1, len(args)):]), false, nil
	case "COALESCE":
		t, err := commonResultType(args, true)
		if err != nil {
			return qtype{}, false, err
		}
		if t.kind == "NUMBER" {
			same := true
			for _, x := range args {
				if !x.isNull() && !sameType(x, t) {
					same = false
				}
			}
			if !same {
				t = numberType()
			}
		}
		return t, false, nil
	case "NULLIF":
		t := arg(0)
		if t.kind == "NUMBER" {
			return numberType(), false, nil
		}
		if t.kind == "CHAR" {
			t.kind = "VARCHAR2"
		}
		return t, false, nil
	case "GREATEST", "LEAST":
		return widenFirst(args), false, nil
	case "JSON_OBJECT", "JSON_ARRAY", "JSON_VALUE", "JSON_QUERY", "JSON_SERIALIZE", "JSON_MERGEPATCH", "JSON_TRANSFORM":
		return varchar2Bytes(4000), false, nil
	case "XMLELEMENT", "XMLFOREST", "XMLCONCAT", "XMLTYPE", "XMLPARSE", "XMLROOT", "XMLQUERY", "XMLCOMMENT", "XMLPI", "XMLCDATA", "XMLSEQUENCE":
		return simpleType("XMLTYPE"), false, nil
	case "STANDARD_HASH":
		t := simpleType("RAW")
		t.bytes = 20
		return t, false, nil
	case "EMPTY_CLOB":
		return simpleType("CLOB"), false, nil
	case "EMPTY_BLOB":
		return simpleType("BLOB"), false, nil
	case "SYSDATE", "CURRENT_DATE":
		return simpleType("DATE"), false, nil
	case "SYSTIMESTAMP", "CURRENT_TIMESTAMP":
		frac := 6
		if v, ok := firstLiteralInt(n); ok {
			frac = v
		}
		return timestampType("TIMESTAMP WITH TIME ZONE", frac), false, nil
	case "LOCALTIMESTAMP":
		frac := 6
		if v, ok := firstLiteralInt(n); ok {
			frac = v
		}
		return timestampType("TIMESTAMP", frac), false, nil
	case "NUMTODSINTERVAL", "TO_DSINTERVAL":
		t := simpleType("INTERVAL DAY TO SECOND")
		t.lead, t.frac = 9, 9
		return t, false, nil
	case "NUMTOYMINTERVAL", "TO_YMINTERVAL":
		t := simpleType("INTERVAL YEAR TO MONTH")
		t.lead = 9
		return t, false, nil
	case "UID":
		return numberType(), false, nil
	case "USER":
		return varchar2Bytes(128), false, nil
	}
	if analyticOnly[name] || rankingNames[name] {
		return qtype{}, false, oraError("30484", "missing window specification for this function")
	}
	// A function this table does not list may be a built-in it lacks, a
	// stand-alone function or one a synonym names: leave it untyped rather
	// than refuse a view Oracle takes.
	return simpleType("UNKNOWN"), false, nil
}

// aggregateType types an aggregate or analytic call.
func (a *queryAnalyzer) aggregateType(n *ast.FuncCallExpr, name string, args []qtype) (qtype, bool, error) {
	arg := func(i int) qtype {
		if i < len(args) {
			return args[i]
		}
		return nullType()
	}
	switch name {
	case "COUNT", "APPROX_COUNT_DISTINCT", "APPROX_COUNT", "REGR_COUNT":
		if arg(0).isLOB() {
			return qtype{}, false, oraError("22849", "Type %s is not supported for this function or operator.", arg(0).kind)
		}
		return numberType(), false, nil
	case "MIN", "MAX", "ANY_VALUE", "FIRST_VALUE", "LAST_VALUE", "NTH_VALUE", "STATS_MODE":
		t := arg(0)
		if t.isLOB() {
			return qtype{}, false, oraError("22849", "Type %s is not supported for this function or operator.", t.kind)
		}
		if (name == "MIN" || name == "MAX") && t.kind == "NUMBER" {
			return numberType(), false, nil
		}
		return t, false, nil
	case "LAG", "LEAD":
		t := arg(0)
		if t.kind == "NUMBER" {
			return numberType(), false, nil
		}
		return t, false, nil
	case "SUM", "AVG", "STDDEV", "STDDEV_POP", "STDDEV_SAMP", "VARIANCE", "VAR_POP", "VAR_SAMP", "MEDIAN",
		"PERCENTILE_CONT", "APPROX_MEDIAN", "APPROX_PERCENTILE", "APPROX_SUM":
		t := arg(0)
		if name == "PERCENTILE_CONT" || name == "APPROX_PERCENTILE" {
			t = numberType()
		}
		switch t.kind {
		case "BINARY_DOUBLE", "BINARY_FLOAT":
			return t, false, nil
		}
		if t.isDatetime() || t.isInterval() {
			return t, false, nil
		}
		return numberType(), false, nil
	case "LISTAGG":
		return varchar2Bytes(4000), false, nil
	case "JSON_ARRAYAGG", "JSON_OBJECTAGG":
		return varchar2Bytes(4000), false, nil
	case "XMLAGG":
		return simpleType("XMLTYPE"), false, nil
	}
	return numberType(), false, nil
}

// argNode is a call's ith argument, or nil.
func argNode(n *ast.FuncCallExpr, i int) ast.Node {
	items := listItems(n.Args)
	if i < 0 || i >= len(items) {
		return nil
	}
	return items[i]
}

func literalInt(n ast.Node) (int, bool) {
	if ex, ok := n.(ast.ExprNode); ok {
		if lit, ok := unparen(ex).(*ast.NumberLiteral); ok && !lit.IsFloat {
			return int(lit.Ival), true
		}
	}
	return 0, false
}

func literalString(n ast.Node) (string, bool) {
	if ex, ok := n.(ast.ExprNode); ok {
		if lit, ok := unparen(ex).(*ast.StringLiteral); ok {
			return lit.Val, true
		}
	}
	return "", false
}

func firstLiteralInt(n *ast.FuncCallExpr) (int, bool) {
	items := listItems(n.Args)
	if len(items) == 0 {
		return 0, false
	}
	return literalInt(items[0])
}

// windowSpec checks an analytic function's window: its partitioning and
// ordering expressions, and those of the WINDOW clause definition it names.
func (a *queryAnalyzer) windowSpec(w *ast.WindowSpec, sc *qscope, seen map[string]bool) error {
	for _, item := range listItems(w.PartitionBy) {
		if ex, ok := item.(ast.ExprNode); ok {
			if _, _, err := a.expr(ex, sc, clauseSelect); err != nil {
				return err
			}
		}
	}
	for _, item := range listItems(w.OrderBy) {
		sb, ok := item.(*ast.SortBy)
		if !ok {
			continue
		}
		t, _, err := a.expr(sb.Expr, sc, clauseSelect)
		if err != nil {
			return err
		}
		if t.isLOB() {
			return oraError("22848", "cannot use %s type as comparison key", t.kind)
		}
	}
	if w.WindowName == "" || seen[w.WindowName] {
		return nil
	}
	for s := sc; s != nil; s = s.parent {
		if def := s.windows[w.WindowName]; def != nil {
			seen[w.WindowName] = true
			return a.windowSpec(def, s, seen)
		}
	}
	return nil
}
