package catalog

import (
	"strings"

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
		case "TRIM":
			// TRIM([mode] [char] FROM source): the parser keeps the mode
			// word as a column reference and the source last.
			items := listItems(n.Args)
			if len(items) > 0 {
				if ref, ok := items[0].(*ast.ColumnRef); ok && ref.Table == "" && isTrimMode(ref.Column) && len(items) > 1 {
					items = items[1:]
				}
			}
			var source qtype
			for i, item := range items {
				ex, ok := item.(ast.ExprNode)
				if !ok {
					continue
				}
				t, _, err := a.expr(ex, sc, clause)
				if err != nil {
					return qtype{}, false, err
				}
				if i == len(items)-1 {
					source = t
				}
			}
			if source.kind == "CHAR" {
				source.kind = "VARCHAR2"
			}
			if source.isNumeric() || source.isDatetime() || source.isInterval() {
				// Trimmed after conversion to characters.
				source = varchar2Bytes(concatLength(source))
			}
			return source, false, nil
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
		// aggregation), and an analytic function's an aggregate; neither
		// may hold an analytic call (ORA-30483).
		argClause = clauseAggArg
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
	// withinGroup is the type of the first WITHIN GROUP ordering expression,
	// which PERCENTILE_CONT and PERCENTILE_DISC return.
	withinGroup := simpleType("UNKNOWN")
	for i, item := range listItems(n.OrderBy) {
		if sb, ok := item.(*ast.SortBy); ok {
			t, _, err := a.expr(sb.Expr, sc, clauseSelect)
			if err != nil {
				return qtype{}, false, err
			}
			if i == 0 {
				withinGroup = t
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
	if n.Returning != nil {
		// A SQL/JSON function returns the type its RETURNING clause names;
		// VARCHAR2 without a length is VARCHAR2(4000).
		t := typeNameType(n.Returning)
		if t.kind == "VARCHAR2" && listItems(n.Returning.TypeMods) == nil {
			t = varchar2Bytes(4000)
		}
		return t, false, nil
	}
	if aggregate || n.Over != nil {
		return a.aggregateType(n, name, args, withinGroup)
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
		t := arg(0)
		if !t.isCharacter() && !t.isLOB() && !t.isNull() && t.kind != "UNKNOWN" {
			// A number or date is converted to characters first.
			t = varchar2Bytes(concatLength(t))
		}
		return t, false, nil
	case "LTRIM", "RTRIM":
		t := arg(0)
		if t.kind == "CHAR" {
			t.kind = "VARCHAR2"
		}
		if t.isNumeric() || t.isDatetime() || t.isInterval() {
			t = varchar2Bytes(concatLength(t))
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
		switch t.kind {
		case "CHAR":
			t.kind = "VARCHAR2"
		case "NCHAR":
			t.kind = "NVARCHAR2"
		}
		return t, false, nil
	case "LPAD", "RPAD":
		t := arg(0)
		if t.isLOB() {
			// A LOB pads to a LOB.
			return t, false, nil
		}
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
		if !t.isCharacter() && !t.isNull() && t.kind != "UNKNOWN" {
			// A number or date is converted to characters first.
			t = varchar2Bytes(concatLength(t))
		}
		// A replacement no longer than what it replaces keeps the source's
		// length; a longer or unknown one leaves the length unknown.
		search, sok := literalString(argNode(n, 1))
		if len(args) > 2 {
			repl, rok := literalString(argNode(n, 2))
			if !sok || !rok || len(repl) > len(search) {
				return simpleType("UNKNOWN"), false, nil
			}
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
		switch name {
		case "SIN", "COS", "TAN", "SINH", "COSH", "TANH", "ASIN", "ACOS", "ATAN":
			// BINARY_FLOAT promotes to BINARY_DOUBLE; others keep theirs.
			switch arg(0).kind {
			case "BINARY_FLOAT", "BINARY_DOUBLE":
				return simpleType("BINARY_DOUBLE"), false, nil
			}
			return numberType(), false, nil
		case "ROUND", "TRUNC":
			// Without a precision a float keeps its type; with one the
			// result is NUMBER.
			if t := arg(0); len(args) < 2 && (t.kind == "BINARY_DOUBLE" || t.kind == "BINARY_FLOAT") {
				return t, false, nil
			}
		case "ABS", "CEIL", "FLOOR", "SQRT", "EXP", "LN", "POWER", "MOD", "REMAINDER", "ATAN2", "NANVL":
			// The arguments' highest numeric precedence.
			out := numberType()
			for _, t := range args {
				if t.isNumeric() {
					out = numericResult(out, t)
				}
			}
			return out, false, nil
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
			// TO_CHAR of characters returns VARCHAR2, a LOB's of the
			// longest VARCHAR2.
			if t.isLOB() {
				return varchar2Bytes(maxVarchar2Bytes), false, nil
			}
			return varchar2Bytes(t.byteLength()), false, nil
		}
		if len(args) > 1 {
			// A format of fixed-width elements gives the length; one that
			// expands (MONTH, DAY, FM, ...) or is not a literal is left
			// untyped rather than given a short length.
			s, ok := literalString(argNode(n, 1))
			if !ok {
				return simpleType("UNKNOWN"), false, nil
			}
			if t.isNumeric() {
				if n, ok := numberFormatWidth(s); ok {
					return varchar2Bytes(n + 1), false, nil
				}
				return simpleType("UNKNOWN"), false, nil
			}
			if n, ok := datetimeFormatWidth(s); ok {
				return varchar2Bytes(n), false, nil
			}
			return simpleType("UNKNOWN"), false, nil
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
		switch {
		case t.isCharacter():
			return characterType("NVARCHAR2", t.byteLength()*5/2, false), false, nil
		case t.isDatetime() && len(args) > 1:
			// A format of fixed-width elements gives the length.
			if f, ok := literalString(argNode(n, 1)); ok {
				if w, ok := datetimeFormatWidth(f); ok {
					return characterType("NVARCHAR2", w, false), false, nil
				}
			}
		}
		return simpleType("UNKNOWN"), false, nil
	case "RAWTOHEX":
		// Two hex digits per byte of a RAW or character value; another
		// type's encoded size is left unknown.
		switch t := arg(0); {
		case t.kind == "RAW", t.kind == "VARCHAR2", t.kind == "CHAR":
			return varchar2Bytes(t.bytes * 2), false, nil
		}
		return simpleType("UNKNOWN"), false, nil
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
		// One type all arguments share is kept; otherwise the first
		// argument's family decides, numbers by numeric precedence (a
		// character argument converts implicitly).
		var first qtype
		same := true
		for _, x := range args {
			if x.isNull() {
				continue
			}
			if first.kind == "" {
				first = x
				continue
			}
			if !sameType(first, x) {
				same = false
			}
		}
		if first.kind == "" {
			return nullType(), false, nil
		}
		if same {
			return first, false, nil
		}
		return widenFirst(args), false, nil
	case "NULLIF":
		t := arg(0)
		if t.isNumeric() {
			// The higher numeric precedence of the two.
			return numericResult(numberType(), numericResult(t, arg(1))), false, nil
		}
		if t.kind == "CHAR" {
			t.kind = "VARCHAR2"
		}
		return t, false, nil
	case "GREATEST", "LEAST":
		return widenFirst(args), false, nil
	case "JSON_QUERY", "JSON_MERGEPATCH", "JSON_TRANSFORM":
		// A JSON target returns JSON by default.
		if arg(0).kind == "JSON" {
			return simpleType("JSON"), false, nil
		}
		return varchar2Bytes(4000), false, nil
	case "JSON_OBJECT", "JSON_ARRAY", "JSON_VALUE", "JSON_SERIALIZE":
		return varchar2Bytes(4000), false, nil
	case "XMLELEMENT", "XMLFOREST", "XMLCONCAT", "XMLTYPE", "XMLPARSE", "XMLROOT", "XMLQUERY", "XMLCOMMENT", "XMLPI", "XMLCDATA", "XMLSEQUENCE":
		return simpleType("XMLTYPE"), false, nil
	case "STANDARD_HASH":
		// The digest's length is its method's: SHA1 when none is named.
		method := "SHA1"
		if len(args) > 1 {
			m, ok := literalString(argNode(n, 1))
			if !ok {
				return simpleType("UNKNOWN"), false, nil
			}
			method = strings.ToUpper(m)
		}
		bytes := map[string]int{"SHA1": 20, "SHA256": 32, "SHA384": 48, "SHA512": 64, "MD5": 16}[method]
		if bytes == 0 {
			return simpleType("UNKNOWN"), false, nil
		}
		t := simpleType("RAW")
		t.bytes = bytes
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
func (a *queryAnalyzer) aggregateType(n *ast.FuncCallExpr, name string, args []qtype, withinGroup qtype) (qtype, bool, error) {
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
			t = withinGroup
		}
		switch t.kind {
		case "BINARY_DOUBLE", "BINARY_FLOAT":
			return t, false, nil
		}
		switch {
		case t.isDatetime():
			// MEDIAN and the percentiles interpolate datetimes; the sums
			// and statistics take numbers (and SUM and AVG intervals).
			switch name {
			case "MEDIAN", "PERCENTILE_CONT", "APPROX_MEDIAN", "APPROX_PERCENTILE":
				return t, false, nil
			}
			return qtype{}, false, oraError("00932", "expression is of data type %s, which is incompatible with expected data type NUMBER", t.kind)
		case t.isInterval():
			return t, false, nil
		}
		return numberType(), false, nil
	case "PERCENTILE_DISC":
		if withinGroup.kind == "NUMBER" {
			return numberType(), false, nil
		}
		return withinGroup, false, nil
	case "LISTAGG":
		if arg(0).kind == "RAW" {
			// A RAW measure aggregates to RAW (2000 bytes under
			// MAX_STRING_SIZE=STANDARD).
			t := simpleType("RAW")
			t.bytes = 2000
			return t, false, nil
		}
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
			if _, _, err := a.expr(ex, sc, clauseAggArg); err != nil {
				return err
			}
		}
	}
	var orderTypes []qtype
	for _, item := range listItems(w.OrderBy) {
		sb, ok := item.(*ast.SortBy)
		if !ok {
			continue
		}
		t, _, err := a.expr(sb.Expr, sc, clauseAggArg)
		if err != nil {
			return err
		}
		if t.isLOB() {
			return oraError("22848", "cannot use %s type as comparison key", t.kind)
		}
		orderTypes = append(orderTypes, t)
	}
	if f := w.Frame; f != nil {
		for _, b := range []*ast.WindowBound{f.Start, f.End} {
			if b == nil || b.Value == nil {
				continue
			}
			if _, _, err := a.expr(b.Value, sc, clauseAggArg); err != nil {
				return err
			}
			if f.Type != ast.WINDOW_RANGE {
				continue
			}
			// A RANGE offset needs one ordering key, numeric or datetime,
			// to add it to.
			if len(orderTypes) > 1 {
				return oraError("30486", "invalid window aggregation group in the window specification")
			}
			if len(orderTypes) == 1 {
				key := orderTypes[0]
				if typeFamily(key) != "" && !key.isNumeric() && !key.isDatetime() {
					return oraError("00902", "invalid datatype")
				}
			}
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
	return oraError("32784", "window %s is not defined", w.WindowName)
}

func isTrimMode(word string) bool {
	return word == "BOTH" || word == "LEADING" || word == "TRAILING"
}

// datetimeFormatWidth is the length of a datetime format made only of
// fixed-width elements and punctuation, or false.
func datetimeFormatWidth(format string) (int, bool) {
	f := strings.ToUpper(format)
	widths := []struct {
		elem  string
		width int
	}{
		{"SYYYY", 5}, {"YYYY", 4}, {"RRRR", 4}, {"HH24", 2}, {"HH12", 2}, {"SSSSS", 5},
		{"YYY", 3}, {"DDD", 3}, {"YY", 2}, {"RR", 2}, {"MM", 2}, {"DD", 2},
		{"HH", 2}, {"MI", 2}, {"SS", 2}, {"WW", 2}, {"IW", 2}, {"CC", 2},
		{"AM", 2}, {"PM", 2}, {"Y", 1}, {"Q", 1}, {"D", 1}, {"W", 1},
	}
	n := 0
	for i := 0; i < len(f); {
		ch := f[i]
		if strings.ContainsRune(" -/,.;:", rune(ch)) {
			n++
			i++
			continue
		}
		matched := false
		for _, w := range widths {
			if strings.HasPrefix(f[i:], w.elem) {
				n += w.width
				i += len(w.elem)
				matched = true
				break
			}
		}
		if !matched {
			return 0, false
		}
	}
	return n, n > 0
}

// numberFormatWidth is the length of a number format made only of digits,
// group and decimal separators and a dollar sign, or false; the sign takes
// one more.
func numberFormatWidth(format string) (int, bool) {
	if format == "" {
		return 0, false
	}
	for _, ch := range strings.ToUpper(format) {
		if !strings.ContainsRune("90.,$GD", ch) {
			return 0, false
		}
	}
	return len(format), true
}
