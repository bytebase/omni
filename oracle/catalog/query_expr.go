package catalog

import (
	"strings"

	"github.com/bytebase/omni/oracle/ast"
)

// exprClause is the clause an expression stands in, which decides whether
// aggregate and analytic functions may appear.
type exprClause int

const (
	clauseSelect exprClause = iota
	clauseWhere
	clauseOn
	clauseGroup
	clauseHaving
	clauseOrder
	clauseFrom
	// clauseAggArg is an aggregate's or analytic function's argument, which
	// may hold an aggregate (nested aggregation) but not an analytic call.
	clauseAggArg
)

func (c exprClause) allowsAggregate() bool {
	return c == clauseSelect || c == clauseHaving || c == clauseOrder || c == clauseAggArg
}

func (c exprClause) allowsWindow() bool {
	return c == clauseSelect || c == clauseOrder
}

// condition analyzes a condition.
func (a *queryAnalyzer) condition(x ast.ExprNode, sc *qscope, clause exprClause) error {
	_, _, err := a.expr(x, sc, clause)
	return err
}

// expr types an expression: its data type and whether it is never NULL.
func (a *queryAnalyzer) expr(x ast.ExprNode, sc *qscope, clause exprClause) (qtype, bool, error) {
	switch n := x.(type) {
	case nil:
		return nullType(), false, nil
	case *ast.ParenExpr:
		return a.expr(n.Expr, sc, clause)
	case *ast.NumberLiteral:
		v := strings.ToUpper(n.Val)
		switch {
		case strings.HasSuffix(v, "D"):
			return simpleType("BINARY_DOUBLE"), false, nil
		case strings.HasSuffix(v, "F"):
			return simpleType("BINARY_FLOAT"), false, nil
		}
		return numberType(), false, nil
	case *ast.StringLiteral:
		if n.Val == "" {
			return nullType(), false, nil
		}
		if n.IsNChar {
			return characterType("NCHAR", len([]rune(n.Val)), false), false, nil
		}
		return charLiteralType(n.Val), false, nil
	case *ast.NullLiteral:
		return nullType(), false, nil
	case *ast.DateTimeLiteral:
		switch strings.ToUpper(n.TypeName) {
		case "DATE":
			return simpleType("DATE"), false, nil
		case "TIMESTAMP":
			if strings.ContainsAny(n.Val, "+") || strings.Count(n.Val, "-") > 2 || strings.ContainsAny(n.Val, "/") {
				return timestampType("TIMESTAMP WITH TIME ZONE", 9), false, nil
			}
			return timestampType("TIMESTAMP", 9), false, nil
		}
		return simpleType(strings.ToUpper(n.TypeName)), false, nil
	case *ast.IntervalExpr:
		if strings.EqualFold(n.From, "YEAR") || strings.EqualFold(n.From, "MONTH") {
			t := simpleType("INTERVAL YEAR TO MONTH")
			t.lead = 2
			return t, false, nil
		}
		t := simpleType("INTERVAL DAY TO SECOND")
		t.lead, t.frac = 2, 6
		return t, false, nil
	case *ast.BindVariable:
		return qtype{}, false, oraError("01027", "bind variables not allowed for data definition operations")
	case *ast.ColumnRef:
		return a.columnRef(n, sc)
	case *ast.PseudoColumn:
		switch n.Type {
		case ast.PSEUDO_ROWID:
			return simpleType("ROWID"), false, nil
		case ast.PSEUDO_ROWNUM, ast.PSEUDO_LEVEL:
			return numberType(), false, nil
		case ast.PSEUDO_SYSDATE:
			return simpleType("DATE"), false, nil
		case ast.PSEUDO_SYSTIMESTAMP:
			return timestampType("TIMESTAMP WITH TIME ZONE", 6), false, nil
		case ast.PSEUDO_USER:
			return varchar2Bytes(128), false, nil
		}
		return numberType(), false, nil
	case *ast.UnaryExpr:
		t, nn, err := a.expr(n.Operand, sc, clause)
		if err != nil {
			return qtype{}, false, err
		}
		switch strings.ToUpper(n.Op) {
		case "PRIOR", "CONNECT_BY_ROOT":
			return t, false, nil
		case "NOT":
			return simpleType("BOOLEAN"), false, nil
		case "-", "+":
			if t.kind == "NUMBER" || t.isCharacter() || t.isNull() {
				return numberType(), nn, nil
			}
			return t, nn, nil
		}
		return t, false, nil
	case *ast.BinaryExpr:
		return a.binary(n, sc, clause)
	case *ast.BoolExpr:
		for _, arg := range listItems(n.Args) {
			if ex, ok := arg.(ast.ExprNode); ok {
				if _, _, err := a.expr(ex, sc, clause); err != nil {
					return qtype{}, false, err
				}
			}
		}
		return simpleType("BOOLEAN"), false, nil
	case *ast.BetweenExpr:
		for _, ex := range []ast.ExprNode{n.Expr, n.Low, n.High} {
			if _, _, err := a.expr(ex, sc, clause); err != nil {
				return qtype{}, false, err
			}
		}
		return simpleType("BOOLEAN"), false, nil
	case *ast.InExpr:
		if _, _, err := a.expr(n.Expr, sc, clause); err != nil {
			return qtype{}, false, err
		}
		width := operandWidth(n.Expr)
		for _, item := range listItems(n.List) {
			if ex, ok := item.(ast.ExprNode); ok {
				if sq, ok := ex.(*ast.SubqueryExpr); ok {
					if err := a.subquery(sq, sc, width); err != nil {
						return qtype{}, false, err
					}
					continue
				}
				if _, _, err := a.expr(ex, sc, clause); err != nil {
					return qtype{}, false, err
				}
			}
		}
		return simpleType("BOOLEAN"), false, nil
	case *ast.LikeExpr:
		for _, ex := range []ast.ExprNode{n.Expr, n.Pattern, n.Escape} {
			if ex == nil {
				continue
			}
			if _, _, err := a.expr(ex, sc, clause); err != nil {
				return qtype{}, false, err
			}
		}
		return simpleType("BOOLEAN"), false, nil
	case *ast.IsExpr:
		if _, _, err := a.expr(n.Expr, sc, clause); err != nil {
			return qtype{}, false, err
		}
		return simpleType("BOOLEAN"), false, nil
	case *ast.ExistsExpr:
		if sel, ok := n.Subquery.(*ast.SelectStmt); ok {
			if _, err := a.query(sel, sc, sc.ctes, false); err != nil {
				return qtype{}, false, err
			}
		}
		return simpleType("BOOLEAN"), false, nil
	case *ast.SubqueryExpr:
		sel, ok := n.Subquery.(*ast.SelectStmt)
		if !ok {
			return simpleType("UNKNOWN"), false, nil
		}
		res, err := a.query(sel, sc, sc.ctes, false)
		if err != nil {
			return qtype{}, false, err
		}
		if len(res.cols) > 1 {
			return qtype{}, false, oraError("00913", "too many values")
		}
		if len(res.cols) == 0 {
			return nullType(), false, nil
		}
		return res.cols[0].typ, false, nil
	case *ast.CursorExpr:
		return qtype{}, false, oraError("22902", "CURSOR expression is not allowed in this context.")
	case *ast.CastExpr:
		if _, _, err := a.expr(n.Arg, sc, clause); err != nil {
			return qtype{}, false, err
		}
		return typeNameType(n.TypeName), false, nil
	case *ast.ExtractExpr:
		if _, _, err := a.expr(n.Expr, sc, clause); err != nil {
			return qtype{}, false, err
		}
		switch strings.ToUpper(n.Field) {
		case "TIMEZONE_REGION", "TIMEZONE_ABBR":
			return varchar2Bytes(64), false, nil
		}
		return numberType(), false, nil
	case *ast.CaseExpr:
		return a.caseExpr(n, sc, clause)
	case *ast.DecodeExpr:
		return a.decodeExpr(n, sc, clause)
	case *ast.FuncCallExpr:
		return a.funcCall(n, sc, clause)
	case *ast.TreatExpr:
		if _, _, err := a.expr(n.Expr, sc, clause); err != nil {
			return qtype{}, false, err
		}
		return typeNameType(n.TypeName), false, nil
	case *ast.Star:
		return simpleType("UNKNOWN"), false, nil
	}
	// Expressions this analyzer does not type are left untyped; their
	// operands are still resolved.
	var err error
	walkExpr(x, func(child ast.Node) bool {
		if err != nil || child == ast.Node(x) {
			return err == nil
		}
		if ex, ok := child.(ast.ExprNode); ok {
			_, _, err = a.expr(ex, sc, clause)
			return false
		}
		return true
	})
	if err != nil {
		return qtype{}, false, err
	}
	return simpleType("UNKNOWN"), false, nil
}

// subquery analyzes a subquery in a condition. want is the number of columns
// it must return, or -1 for any.
func (a *queryAnalyzer) subquery(sq *ast.SubqueryExpr, sc *qscope, want int) error {
	sel, ok := sq.Subquery.(*ast.SelectStmt)
	if !ok {
		return nil
	}
	res, err := a.query(sel, sc, sc.ctes, false)
	if err != nil {
		return err
	}
	switch {
	case want >= 0 && len(res.cols) > want:
		return oraError("00913", "too many values")
	case want >= 0 && len(res.cols) < want:
		return oraError("00947", "not enough values")
	}
	return nil
}

// keywordValues are the words Oracle reads as a value when no column of
// that name is in scope.
var keywordValues = map[string]qtype{
	"TRUE":               simpleType("BOOLEAN"),
	"FALSE":              simpleType("BOOLEAN"),
	"CURRENT_DATE":       simpleType("DATE"),
	"CURRENT_TIMESTAMP":  timestampType("TIMESTAMP WITH TIME ZONE", 6),
	"LOCALTIMESTAMP":     timestampType("TIMESTAMP", 6),
	"SYSDATE":            simpleType("DATE"),
	"SYSTIMESTAMP":       timestampType("TIMESTAMP WITH TIME ZONE", 6),
	"UID":                numberType(),
	"USER":               varchar2Bytes(128),
	"CONNECT_BY_ISLEAF":  numberType(),
	"CONNECT_BY_ISCYCLE": numberType(),
	"ORA_ROWSCN":         numberType(),
	"SESSIONTIMEZONE":    varchar2Bytes(75),
	"DBTIMEZONE":         varchar2Bytes(6),
	"ROWNUM":             numberType(),
	"LEVEL":              numberType(),
}

func (a *queryAnalyzer) columnRef(n *ast.ColumnRef, sc *qscope) (qtype, bool, error) {
	c, found, err := a.resolveColumn(n, sc)
	if err != nil {
		return qtype{}, false, err
	}
	if found {
		for s := sc; s != nil; s = s.parent {
			if a.resolvesLocally(n, s) {
				a.owner[n] = s.from
				break
			}
		}
		return c.typ, c.notNull, nil
	}
	if n.Table == "" {
		if t, ok := keywordValues[n.Column]; ok {
			return t, false, nil
		}
		return qtype{}, false, oraError("00904", "%q: invalid identifier", n.Column)
	}
	if n.Schema != "" {
		return qtype{}, false, oraError("00904", "%q.%q.%q: invalid identifier", n.Schema, n.Table, n.Column)
	}
	return qtype{}, false, oraError("00904", "%q.%q: invalid identifier", n.Table, n.Column)
}

// binary types an operator expression.
func (a *queryAnalyzer) binary(n *ast.BinaryExpr, sc *qscope, clause exprClause) (qtype, bool, error) {
	lt, lnn, err := a.expr(n.Left, sc, clause)
	if err != nil {
		return qtype{}, false, err
	}
	if isComparison(n.Op) {
		// A subquery compared with a value, or quantified by ANY, SOME or
		// ALL, returns as many columns as the left side holds.
		if sq := comparedSubquery(n.Right); sq != nil {
			if err := a.subquery(sq, sc, operandWidth(n.Left)); err != nil {
				return qtype{}, false, err
			}
			return simpleType("BOOLEAN"), false, nil
		}
	}
	rt, rnn, err := a.expr(n.Right, sc, clause)
	if err != nil {
		return qtype{}, false, err
	}
	nn := lnn && rnn
	switch strings.ToUpper(n.Op) {
	case "||":
		return concatType(lt, rt), false, nil
	case "+", "-":
		return arithmeticType(n.Op, lt, rt), false, nil
	case "*", "/", "**":
		return multiplyType(lt, rt), false, nil
	case "=", "!=", "<>", "^=", "~=", "<", ">", "<=", ">=":
		if err := comparable(lt, rt, n.Right); err != nil {
			return qtype{}, false, err
		}
		return simpleType("BOOLEAN"), false, nil
	case "AND", "OR":
		return simpleType("BOOLEAN"), false, nil
	}
	_ = nn
	return simpleType("UNKNOWN"), false, nil
}

// comparable checks that two values can be compared: a datetime does not
// compare with a number.
func comparable(l, r qtype, right ast.ExprNode) error {
	if (l.isDatetime() && r.isNumeric()) || (l.isNumeric() && r.isDatetime()) {
		number, other := r, l
		if l.isNumeric() {
			number, other = l, r
		}
		_ = number
		return oraError("00932", "expression is of data type NUMBER, which is incompatible with expected data type %s", other.kind)
	}
	return nil
}

func concatType(l, r qtype) qtype {
	if l.kind == "CLOB" || r.kind == "CLOB" || l.kind == "NCLOB" || r.kind == "NCLOB" {
		return simpleType("CLOB")
	}
	if l.isNational() || r.isNational() {
		return characterType("NVARCHAR2", l.byteLength()/dbCharBytes+r.byteLength()/dbCharBytes, false)
	}
	if (l.kind == "CHAR" || l.isNull()) && (r.kind == "CHAR" || r.isNull()) && !(l.isNull() && r.isNull()) {
		// Two CHAR operands concatenate to a CHAR.
		return characterType("CHAR", min(l.byteLength()+r.byteLength(), maxCharBytes), false)
	}
	return varchar2Bytes(concatLength(l) + concatLength(r))
}

// concatLength is the bytes a value takes as a concatenation operand.
func concatLength(t qtype) int {
	switch {
	case t.isCharacter():
		return t.byteLength()
	case t.isNumeric():
		return 40
	case t.kind == "DATE":
		return 75
	case t.isDatetime():
		return 75
	case t.isNull():
		return 0
	}
	return t.byteLength()
}

func arithmeticType(op string, l, r qtype) qtype {
	switch {
	case l.kind == "DATE" && r.kind == "DATE" && op == "-":
		t := numberType()
		t.scale = 0
		return t
	case l.isDatetime() && r.isDatetime() && op == "-":
		t := simpleType("INTERVAL DAY TO SECOND")
		t.lead, t.frac = 9, 6
		return t
	case l.isDatetime() && (r.isNumeric() || r.isCharacter() || r.isNull()):
		return simpleType("DATE")
	case r.isDatetime() && l.isNumeric():
		return simpleType("DATE")
	case l.isDatetime() && r.isInterval():
		return l
	case l.isInterval() && r.isInterval():
		return l
	}
	return numericResult(l, r)
}

func multiplyType(l, r qtype) qtype {
	switch {
	case l.kind == "INTERVAL DAY TO SECOND" || r.kind == "INTERVAL DAY TO SECOND":
		t := simpleType("INTERVAL DAY TO SECOND")
		t.lead, t.frac = 9, 9
		return t
	case l.kind == "INTERVAL YEAR TO MONTH" || r.kind == "INTERVAL YEAR TO MONTH":
		t := simpleType("INTERVAL YEAR TO MONTH")
		t.lead = 9
		return t
	}
	return numericResult(l, r)
}

func numericResult(l, r qtype) qtype {
	switch {
	case l.kind == "BINARY_DOUBLE" || r.kind == "BINARY_DOUBLE":
		return simpleType("BINARY_DOUBLE")
	case l.kind == "BINARY_FLOAT" || r.kind == "BINARY_FLOAT":
		return simpleType("BINARY_FLOAT")
	}
	return numberType()
}

// caseExpr types a CASE expression: its results' common type, which keeps
// a type all results share and otherwise widens characters to VARCHAR2.
func (a *queryAnalyzer) caseExpr(n *ast.CaseExpr, sc *qscope, clause exprClause) (qtype, bool, error) {
	if n.Arg != nil {
		if _, _, err := a.expr(n.Arg, sc, clause); err != nil {
			return qtype{}, false, err
		}
	}
	var results []qtype
	for _, item := range listItems(n.Whens) {
		w, ok := item.(*ast.CaseWhen)
		if !ok {
			continue
		}
		if _, _, err := a.expr(w.Condition, sc, clause); err != nil {
			return qtype{}, false, err
		}
		t, _, err := a.expr(w.Result, sc, clause)
		if err != nil {
			return qtype{}, false, err
		}
		results = append(results, t)
	}
	if n.Default != nil {
		t, _, err := a.expr(n.Default, sc, clause)
		if err != nil {
			return qtype{}, false, err
		}
		results = append(results, t)
	}
	t, err := commonResultType(results, true)
	return t, false, err
}

// commonResultType is the type of CASE's results. strict refuses results
// of different families (ORA-00932).
func commonResultType(results []qtype, strict bool) (qtype, error) {
	var out qtype
	have := false
	for _, t := range results {
		if t.isNull() {
			continue
		}
		if !have {
			out, have = t, true
			continue
		}
		if sameType(out, t) {
			continue
		}
		switch {
		case out.isCharacter() && t.isCharacter():
			out = varchar2Bytes(max(out.byteLength(), t.byteLength()))
		case out.isNumeric() && t.isNumeric():
			out = numericResult(out, t)
		case out.isDatetime() && t.isDatetime():
			if out.kind == "DATE" {
				out = t
			}
		case strict && typeFamily(out) != "" && typeFamily(t) != "":
			return qtype{}, oraError("00932", "expression is of data type %s, which is incompatible with expected data type %s", t.kind, out.kind)
		default:
			// A type this analyzer does not model proves no mismatch.
			return simpleType("UNKNOWN"), nil
		}
	}
	if !have {
		return nullType(), nil
	}
	return out, nil
}

// decodeExpr types DECODE: its first result's family, character results
// widened to VARCHAR2 and numbers to NUMBER.
func (a *queryAnalyzer) decodeExpr(n *ast.DecodeExpr, sc *qscope, clause exprClause) (qtype, bool, error) {
	if _, _, err := a.expr(n.Arg, sc, clause); err != nil {
		return qtype{}, false, err
	}
	var results []qtype
	for _, item := range listItems(n.Pairs) {
		p, ok := item.(*ast.DecodePair)
		if !ok {
			continue
		}
		if _, _, err := a.expr(p.Search, sc, clause); err != nil {
			return qtype{}, false, err
		}
		t, _, err := a.expr(p.Result, sc, clause)
		if err != nil {
			return qtype{}, false, err
		}
		results = append(results, t)
	}
	if n.Default != nil {
		t, _, err := a.expr(n.Default, sc, clause)
		if err != nil {
			return qtype{}, false, err
		}
		results = append(results, t)
	}
	return widenFirst(results), false, nil
}

// widenFirst is the result type of DECODE, NVL and the like: the first
// non-NULL argument's family, with character arguments widened to a
// VARCHAR2 of the longest and numbers to NUMBER.
func widenFirst(args []qtype) qtype {
	var first qtype
	have := false
	for _, t := range args {
		if !t.isNull() {
			first, have = t, true
			break
		}
	}
	if !have {
		return nullType()
	}
	switch {
	case first.isCharacter():
		if first.isNational() {
			n := 0
			for _, t := range args {
				n = max(n, t.chars)
			}
			return characterType("NVARCHAR2", n, false)
		}
		n := 0
		for _, t := range args {
			n = max(n, t.byteLength())
		}
		out := varchar2Bytes(n)
		if first.charSem {
			out.charSem = true
			out.chars = n
		}
		return out
	case first.isNumeric():
		// The argument of the highest numeric precedence decides.
		out := numberType()
		for _, t := range args {
			if t.isNumeric() {
				out = numericResult(out, t)
			}
		}
		return out
	}
	return first
}

// aggregateNames are Oracle's aggregate functions, which are analytic when
// they carry OVER.
var aggregateNames = map[string]bool{
	"AVG": true, "COUNT": true, "MAX": true, "MIN": true, "SUM": true,
	"STDDEV": true, "STDDEV_POP": true, "STDDEV_SAMP": true,
	"VARIANCE": true, "VAR_POP": true, "VAR_SAMP": true, "MEDIAN": true,
	"LISTAGG": true, "COLLECT": true, "CORR": true, "COVAR_POP": true, "COVAR_SAMP": true,
	"PERCENTILE_CONT": true, "PERCENTILE_DISC": true, "GROUPING": true, "GROUPING_ID": true, "GROUP_ID": true,
	"JSON_ARRAYAGG": true, "JSON_OBJECTAGG": true, "XMLAGG": true, "APPROX_COUNT_DISTINCT": true,
	"ANY_VALUE": true, "STATS_MODE": true, "REGR_SLOPE": true, "REGR_INTERCEPT": true,
	"REGR_COUNT": true, "REGR_R2": true, "REGR_AVGX": true, "REGR_AVGY": true,
	"REGR_SXX": true, "REGR_SYY": true, "REGR_SXY": true, "KURTOSIS_POP": true, "KURTOSIS_SAMP": true,
	"SKEWNESS_POP": true, "SKEWNESS_SAMP": true, "BIT_AND_AGG": true, "BIT_OR_AGG": true, "BIT_XOR_AGG": true,
	"CHECKSUM": true, "APPROX_MEDIAN": true, "APPROX_PERCENTILE": true, "APPROX_SUM": true, "APPROX_COUNT": true,
}

// analyticOnly are the functions that need OVER.
var analyticOnly = map[string]bool{
	"ROW_NUMBER": true, "NTILE": true, "LAG": true, "LEAD": true,
	"FIRST_VALUE": true, "LAST_VALUE": true, "NTH_VALUE": true, "RATIO_TO_REPORT": true,
}

// rankingNames are RANK and its kin: aggregate with WITHIN GROUP, analytic
// with OVER.
var rankingNames = map[string]bool{"RANK": true, "DENSE_RANK": true, "PERCENT_RANK": true, "CUME_DIST": true}

func isAggregateName(name string) bool { return aggregateNames[name] || rankingNames[name] }

// isAggregateCall reports an aggregate function call (not analytic).
func isAggregateCall(fc *ast.FuncCallExpr) bool {
	if fc == nil || fc.FuncName == nil || fc.Over != nil || fc.FuncName.Schema != "" {
		return false
	}
	name := fc.FuncName.Name
	if rankingNames[name] {
		return fc.OrderBy != nil
	}
	return aggregateNames[name]
}

// containsAggregate reports an aggregate call in x outside subqueries.
func containsAggregate(x ast.Node, _ bool) bool {
	found := false
	walkExpr(x, func(n ast.Node) bool {
		if found {
			return false
		}
		switch v := n.(type) {
		case *ast.SubqueryExpr, *ast.ExistsExpr, *ast.SelectStmt:
			return false
		case *ast.FuncCallExpr:
			if isAggregateCall(v) {
				found = true
				return false
			}
		}
		return true
	})
	return found
}

// containsNestedAggregate reports an aggregate whose argument holds one.
func containsNestedAggregate(x ast.Node) bool {
	found := false
	walkExpr(x, func(n ast.Node) bool {
		if found {
			return false
		}
		switch v := n.(type) {
		case *ast.SubqueryExpr, *ast.ExistsExpr, *ast.SelectStmt:
			return false
		case *ast.FuncCallExpr:
			if isAggregateCall(v) {
				for _, arg := range listItems(v.Args) {
					if containsAggregate(arg, false) {
						found = true
						return false
					}
				}
				return false
			}
		}
		return true
	})
	return found
}

// containsWindow reports an analytic call in x outside subqueries.
func containsWindow(x ast.Node) bool {
	found := false
	walkExpr(x, func(n ast.Node) bool {
		if found {
			return false
		}
		switch v := n.(type) {
		case *ast.SubqueryExpr, *ast.ExistsExpr, *ast.SelectStmt:
			return false
		case *ast.FuncCallExpr:
			if v.Over != nil {
				found = true
				return false
			}
		}
		return true
	})
	return found
}

// checkGrouped checks an expression of a grouped query block: each column
// it reads outside an aggregate is a GROUP BY expression.
func checkGrouped(x ast.Node, keys []string, keyOf func(ast.Node) string, owned func(*ast.ColumnRef) bool, hasGroupBy, positional, nested bool) error {
	if ex, ok := x.(ast.ExprNode); ok {
		key := keyOf(ex)
		for _, k := range keys {
			if k == key {
				if nested {
					return oraError("00937", "not a single-group group function")
				}
				return nil
			}
		}
	}
	switch n := x.(type) {
	case *ast.FuncCallExpr:
		if isAggregateCall(n) {
			return nil
		}
	case *ast.ColumnRef:
		if n.Column == "*" {
			return nil
		}
		if _, ok := keywordValues[n.Column]; ok && n.Table == "" {
			return nil
		}
		if keyOf(n) == outerColumn {
			// A column of an enclosing query block is a constant here.
			return nil
		}
		switch {
		case !hasGroupBy || nested:
			return oraError("00937", "not a single-group group function")
		case positional:
			return oraError("03162", "%q: must appear in the GROUP BY clause or be used in an aggregate function as 'group_by_position_enabled' is FALSE", n.Column)
		}
		return oraError("00979", "%q: must appear in the GROUP BY clause or be used in an aggregate function", n.Column)
	case *ast.SubqueryExpr, *ast.ExistsExpr:
		// A subquery is grouped by its own rules, but a column of this
		// block it reads (a correlated reference) is held to this block's.
		var err error
		walkExpr(x, func(child ast.Node) bool {
			if err != nil {
				return false
			}
			if ref, ok := child.(*ast.ColumnRef); ok && owned(ref) {
				err = checkGrouped(ref, keys, keyOf, owned, hasGroupBy, positional, nested)
			}
			return true
		})
		return err
	case *ast.NumberLiteral, *ast.StringLiteral, *ast.NullLiteral, *ast.DateTimeLiteral, *ast.PseudoColumn, *ast.IntervalExpr:
		return nil
	}
	var err error
	walkChildren(x, func(child ast.Node) {
		if err == nil {
			err = checkGrouped(child, keys, keyOf, owned, hasGroupBy, positional, nested)
		}
	})
	return err
}

// walkChildren calls fn on each direct AST child of n.
func walkChildren(n ast.Node, fn func(ast.Node)) {
	first := true
	walkExpr(n, func(child ast.Node) bool {
		if first {
			first = false
			return true
		}
		if _, ok := child.(*ast.List); ok {
			return true
		}
		fn(child)
		return false
	})
}

// typeFamily is the family of types that convert into one another in CASE,
// or "" for a type this analyzer does not know well enough to compare.
func typeFamily(t qtype) string {
	switch {
	case t.isCharacter():
		return "character"
	case t.isNumeric():
		return "numeric"
	case t.isDatetime():
		return "datetime"
	case t.isInterval():
		return t.kind
	case t.kind == "BOOLEAN", t.kind == "RAW":
		return t.kind
	}
	return ""
}

func isComparison(op string) bool {
	switch strings.ToUpper(op) {
	case "=", "!=", "<>", "^=", "~=", "<", ">", "<=", ">=":
		return true
	}
	return false
}

// comparedSubquery is the subquery a comparison's right side is, directly or
// under ANY, SOME or ALL, or nil.
func comparedSubquery(x ast.ExprNode) *ast.SubqueryExpr {
	switch v := unparen(x).(type) {
	case *ast.SubqueryExpr:
		return v
	case *ast.FuncCallExpr:
		if v.FuncName == nil || v.FuncName.Schema != "" {
			return nil
		}
		switch v.FuncName.Name {
		case "ANY", "SOME", "ALL":
			items := listItems(v.Args)
			if len(items) == 1 {
				if sq, ok := items[0].(*ast.SubqueryExpr); ok {
					return sq
				}
			}
		}
	}
	return nil
}

// operandWidth is how many values an operand holds: a row constructor's
// width, else one.
func operandWidth(x ast.ExprNode) int {
	if row, ok := unparen(x).(*ast.FuncCallExpr); ok && (row.FuncName == nil || row.FuncName.Name == "") {
		return len(listItems(row.Args))
	}
	return 1
}
