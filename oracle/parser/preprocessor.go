package parser

import (
	"math/big"
	"strconv"
	"strings"
	"unicode/utf8"

	nodes "github.com/bytebase/omni/oracle/ast"
)

// staticCondition evaluates the conditional compilation condition at
// source[loc.Start:loc.End] as far as the text settles it. Oracle 23ai
// folds the Boolean operators over TRUE, FALSE, and NULL, short-circuiting
// past what it cannot settle (TRUE OR $$x is TRUE), compares numbers as
// its NUMBER holds them (0.1 + 0.2 = 0.3) and Booleans (FALSE < TRUE), and
// finds a comparison with NULL NULL; inquiry directives and package
// constants have values the text does not hold.
//
// It reports false for a condition that is no static Boolean expression,
// which the parser rejects and Split lexes past.
func staticCondition(source string, loc nodes.Loc) (ccValue, bool) {
	sub := &Parser{lexer: NewLexerRange(source, loc.Start, loc.End), source: source}
	sub.advance()
	expr, err := sub.parseExpr()
	if err != nil || expr == nil || sub.cur.Type != tokEOF || sub.staticConditionFault(expr) != nil {
		return ccUnknown, false
	}
	switch staticTruthOf(source, expr) {
	case truthTrue:
		return ccTrue, true
	case truthFalse, truthNull:
		return ccFalse, true
	}
	return ccUnknown, true
}

// staticTruth is the three-valued truth of a static condition, or unknown.
type staticTruth int

const (
	truthUnknown staticTruth = iota
	truthTrue
	truthFalse
	truthNull
)

// staticTruthOf evaluates expr as a static Boolean expression.
func staticTruthOf(source string, expr nodes.ExprNode) staticTruth {
	switch e := expr.(type) {
	case *nodes.ParenExpr:
		return staticTruthOf(source, e.Expr)
	case *nodes.NullLiteral:
		return truthNull
	case *nodes.ColumnRef:
		// TRUE and FALSE, not a quoted "TRUE".
		if e.Schema != "" || e.Table != "" || e.Loc.Start < 0 || e.Loc.Start >= len(source) || source[e.Loc.Start] == '"' {
			return truthUnknown
		}
		switch e.Column {
		case "TRUE":
			return truthTrue
		case "FALSE":
			return truthFalse
		}
	case *nodes.BoolExpr:
		return staticBoolExpr(source, e)
	case *nodes.FuncCallExpr:
		// STANDARD's NVL and XOR, the Boolean functions a static
		// expression may call.
		if !isStaticBooleanFunction(e) {
			return truthUnknown
		}
		a := staticTruthOf(source, e.Args.Items[0].(nodes.ExprNode))
		b := staticTruthOf(source, e.Args.Items[1].(nodes.ExprNode))
		switch e.FuncName.Name {
		case "NVL":
			if a == truthNull {
				return b
			}
			return a
		case "XOR":
			switch {
			case a == truthUnknown || b == truthUnknown:
				return truthUnknown
			case a == truthNull || b == truthNull:
				return truthNull
			case a != b:
				return truthTrue
			}
			return truthFalse
		}
	case *nodes.IsExpr:
		if e.Test != "NULL" {
			return truthUnknown
		}
		_, null, ok := staticComparable(source, e.Expr)
		switch {
		case !ok:
			return truthUnknown
		case null != e.Not:
			return truthTrue
		}
		return truthFalse
	case *nodes.BinaryExpr:
		left, lnull, lok := staticComparable(source, e.Left)
		right, rnull, rok := staticComparable(source, e.Right)
		switch {
		case !isComparison(e.Op) || !lok || !rok:
			return truthUnknown
		case lnull || rnull:
			return truthNull
		}
		var holds bool
		switch cmp := left.Cmp(right); e.Op {
		case "=":
			holds = cmp == 0
		case "!=", "<>", "~=", "^=":
			holds = cmp != 0
		case "<":
			holds = cmp < 0
		case ">":
			holds = cmp > 0
		case "<=":
			holds = cmp <= 0
		case ">=":
			holds = cmp >= 0
		default:
			return truthUnknown
		}
		if holds {
			return truthTrue
		}
		return truthFalse
	}
	return truthUnknown
}

// isComparison reports whether op compares its operands.
func isComparison(op string) bool {
	switch op {
	case "=", "!=", "<>", "~=", "^=", "<", ">", "<=", ">=":
		return true
	}
	return false
}

// staticComparable returns the value of a static number or Boolean that
// a comparison or IS NULL tests, a Boolean as 0 for FALSE and 1 for TRUE;
// null when it is NULL; and false when the text does not settle it.
func staticComparable(source string, expr nodes.ExprNode) (value *big.Rat, null bool, ok bool) {
	if v, ok := staticNumber(expr); ok {
		return v, v == nil, true
	}
	switch staticTruthOf(source, expr) {
	case truthTrue:
		return big.NewRat(1, 1), false, true
	case truthFalse:
		return new(big.Rat), false, true
	case truthNull:
		return nil, true, true
	}
	return nil, false, false
}

// staticBoolExpr evaluates NOT, AND, and OR in three-valued logic.
func staticBoolExpr(source string, e *nodes.BoolExpr) staticTruth {
	if e.Args == nil || e.Args.Len() == 0 {
		return truthUnknown
	}
	if e.Boolop == nodes.BOOL_NOT {
		switch staticTruthOf(source, e.Args.Items[0].(nodes.ExprNode)) {
		case truthTrue:
			return truthFalse
		case truthFalse:
			return truthTrue
		case truthNull:
			return truthNull
		}
		return truthUnknown
	}
	// AND is FALSE when an operand is, OR is TRUE when one is; otherwise
	// an unknown operand leaves it unknown, and a NULL one makes it NULL.
	dominant, other := truthFalse, truthTrue
	if e.Boolop == nodes.BOOL_OR {
		dominant, other = truthTrue, truthFalse
	}
	result := other
	for _, item := range e.Args.Items {
		switch v := staticTruthOf(source, item.(nodes.ExprNode)); v {
		case dominant:
			return dominant
		case truthUnknown:
			result = truthUnknown
		case truthNull:
			if result != truthUnknown {
				result = truthNull
			}
		}
	}
	return result
}

// staticNumber returns the value of a static numeric expression, nil for
// NULL, or false when the text does not settle it: a numeric literal,
// arithmetic over values it settles (Oracle 23ai selects the $ELSE of
// $IF 1 + 1 = 3), or a static numeric function (staticNumericCall). Each
// value is rounded as Oracle's NUMBER holds it (oracleNumber).
func staticNumber(expr nodes.ExprNode) (*big.Rat, bool) {
	switch e := expr.(type) {
	case *nodes.ParenExpr:
		return staticNumber(e.Expr)
	case *nodes.NullLiteral:
		return nil, true
	case *nodes.NumberLiteral:
		return staticLiteral(e.Val)
	case *nodes.UnaryExpr:
		v, ok := staticNumber(e.Operand)
		switch {
		case !ok || e.Op != "-" && e.Op != "+":
			return nil, false
		case v == nil || e.Op == "+":
			return v, true
		}
		return v.Neg(v), true
	case *nodes.FuncCallExpr:
		return staticNumericCall(e)
	case *nodes.BinaryExpr:
		switch e.Op {
		case "+", "-", "*", "/", "**":
		default:
			return nil, false
		}
		left, lok := staticNumber(e.Left)
		right, rok := staticNumber(e.Right)
		switch {
		case !lok || !rok:
			return nil, false
		case left == nil || right == nil:
			return nil, true
		}
		switch e.Op {
		case "+":
			return oracleNumber(left.Add(left, right))
		case "-":
			return oracleNumber(left.Sub(left, right))
		case "*":
			return oracleNumber(left.Mul(left, right))
		case "/":
			if right.Sign() == 0 {
				return nil, false
			}
			return oracleNumber(left.Quo(left, right))
		}
		return staticPower(left, right)
	}
	return nil, false
}

// staticNumericCall folds the static numeric functions whose value the
// text settles: ABS, SIGN, CEIL, FLOOR, TRUNC and ROUND (to a scale, which
// rounds to a whole number: Oracle 23ai finds ROUND(1.25, 1.5) = 1.25),
// MOD (Oracle's, with the dividend's sign), POWER to a whole exponent,
// LENGTH and LENGTHC of a string literal, and NVL of numbers. A NULL
// argument makes each NULL (Oracle 23ai selects the $ELSE of
// $IF ABS(-1) = 2).
func staticNumericCall(call *nodes.FuncCallExpr) (*big.Rat, bool) {
	if call.FuncName == nil || call.Args == nil {
		return nil, false
	}
	switch call.FuncName.Schema {
	case "", "STANDARD", "SYS.STANDARD":
	default:
		return nil, false
	}
	name := call.FuncName.Name
	// LENGTH and LENGTHC of a string literal count its characters; '' is
	// NULL.
	if (name == "LENGTH" || name == "LENGTHC") && call.Args.Len() == 1 {
		if lit, ok := call.Args.Items[0].(*nodes.StringLiteral); ok {
			if lit.Val == "" {
				return nil, true
			}
			return big.NewRat(int64(utf8.RuneCountInString(lit.Val)), 1), true
		}
	}
	if name == "NVL" && call.Args.Len() == 2 {
		v, ok := staticNumber(call.Args.Items[0].(nodes.ExprNode))
		if !ok || v != nil {
			return v, ok
		}
		return staticNumber(call.Args.Items[1].(nodes.ExprNode))
	}
	var args []*big.Rat
	null := false
	for _, item := range call.Args.Items {
		v, ok := staticNumber(item.(nodes.ExprNode))
		if !ok {
			return nil, false
		}
		null = null || v == nil
		args = append(args, v)
	}
	one := len(args) == 1
	switch {
	case name == "ABS" && one, name == "SIGN" && one, name == "CEIL" && one, name == "FLOOR" && one,
		name == "TRUNC" && len(args) <= 2, name == "ROUND" && len(args) <= 2,
		name == "MOD" && len(args) == 2, name == "POWER" && len(args) == 2:
		if null {
			return nil, true
		}
	default:
		return nil, false
	}
	switch name {
	case "ABS":
		return args[0].Abs(args[0]), true
	case "SIGN":
		return big.NewRat(int64(args[0].Sign()), 1), true
	case "CEIL", "FLOOR":
		return oracleNumber(roundRat(args[0], name))
	case "TRUNC", "ROUND":
		if one {
			return oracleNumber(roundRat(args[0], name))
		}
		scale := roundRat(args[1], "ROUND")
		switch {
		case scale.Cmp(big.NewRat(200, 1)) > 0:
			// Finer than any NUMBER.
			return args[0], true
		case scale.Cmp(big.NewRat(-200, 1)) < 0:
			// Coarser than any NUMBER.
			return new(big.Rat), true
		}
		unit := pow10(int(scale.Num().Int64()))
		v := roundRat(new(big.Rat).Mul(args[0], unit), name)
		return oracleNumber(v.Quo(v, unit))
	case "MOD":
		if args[1].Sign() == 0 {
			return args[0], true
		}
		q := roundRat(new(big.Rat).Quo(args[0], args[1]), "TRUNC")
		return oracleNumber(args[0].Sub(args[0], q.Mul(q, args[1])))
	}
	return staticPower(args[0], args[1])
}

// staticPower raises base to a whole exponent of at most 64, when NUMBER
// holds the exact power; past that, Oracle's rounding of the steps is not
// the text's to settle.
func staticPower(base, exp *big.Rat) (*big.Rat, bool) {
	if !exp.IsInt() || exp.Sign() < 0 || exp.Num().Cmp(big.NewInt(64)) > 0 {
		return nil, false
	}
	result := big.NewRat(1, 1)
	for i := exp.Num().Int64(); i > 0; i-- {
		result.Mul(result, base)
	}
	if rounded, ok := oracleNumber(new(big.Rat).Set(result)); !ok || rounded.Cmp(result) != 0 {
		return nil, false
	}
	return result, true
}

// staticLiteral returns the value of a numeric literal as NUMBER holds it,
// or false for a binary floating-point literal (1.5f) or one too long or
// far out of range to build exactly.
func staticLiteral(val string) (*big.Rat, bool) {
	if len(val) > 1000 {
		return nil, false
	}
	if i := strings.IndexAny(val, "eE"); i >= 0 {
		exp, err := strconv.Atoi(val[i+1:])
		if err != nil || exp > 1000 || exp < -1000 {
			return nil, false
		}
	}
	r, ok := new(big.Rat).SetString(val)
	if !ok {
		return nil, false
	}
	return oracleNumber(r)
}

// oracleNumber rounds r as Oracle's NUMBER holds a value: to 20 base-100
// digits, half away from zero, and to zero below 1e-130. Oracle 23ai finds
// 1 / 3 = 0.3333333333333333333333333333333333333333 (40 digits), 1 / 3 *
// 3 = 1 FALSE, and 2 / 3 * 3 = 2 TRUE; a literal rounds the same way. It
// reports false from 1e126 on, past NUMBER's range.
func oracleNumber(r *big.Rat) (*big.Rat, bool) {
	if r.Sign() == 0 {
		return r, true
	}
	abs := new(big.Rat).Abs(r)
	if abs.Cmp(pow10(-130)) < 0 {
		return new(big.Rat), true
	}
	if abs.Cmp(pow10(126)) >= 0 {
		return nil, false
	}
	// 10^e <= abs < 10^(e+1); the leading base-100 digit is 100^floor(e/2).
	e := len(abs.Num().String()) - len(abs.Denom().String())
	for abs.Cmp(pow10(e)) < 0 {
		e--
	}
	for abs.Cmp(pow10(e+1)) >= 0 {
		e++
	}
	lead := e / 2
	if e < 0 && e%2 != 0 {
		lead--
	}
	unit := pow10(2*lead - 38)
	n := roundRat(abs.Quo(abs, unit), "ROUND")
	n.Mul(n, unit)
	if n.Cmp(pow10(126)) >= 0 {
		return nil, false
	}
	if r.Sign() < 0 {
		n.Neg(n)
	}
	return n, true
}

// pow10 returns 10^e.
func pow10(e int) *big.Rat {
	if e < 0 {
		return new(big.Rat).Inv(pow10(-e))
	}
	return new(big.Rat).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(e)), nil))
}

// roundRat rounds r to a whole number the way CEIL, FLOOR, TRUNC, or ROUND
// (half away from zero) does.
func roundRat(r *big.Rat, mode string) *big.Rat {
	num, den := r.Num(), r.Denom()
	q, m := new(big.Int).QuoRem(num, den, new(big.Int)) // truncates toward zero
	if m.Sign() != 0 {
		switch mode {
		case "CEIL":
			if r.Sign() > 0 {
				q.Add(q, big.NewInt(1))
			}
		case "FLOOR":
			if r.Sign() < 0 {
				q.Sub(q, big.NewInt(1))
			}
		case "ROUND":
			twice := new(big.Int).Mul(new(big.Int).Abs(m), big.NewInt(2))
			if twice.Cmp(den) >= 0 {
				q.Add(q, big.NewInt(int64(r.Sign())))
			}
		}
	}
	return new(big.Rat).SetInt(q)
}

// isStaticBooleanFunction reports whether call is STANDARD's NVL or XOR,
// with two arguments: the Boolean functions a static expression may call.
func isStaticBooleanFunction(call *nodes.FuncCallExpr) bool {
	if call.FuncName == nil || call.Args == nil || call.Args.Len() != 2 {
		return false
	}
	return staticFunctionName(call) == "NVL" || staticFunctionName(call) == "XOR"
}

// staticFunctionName returns the name of call when it may name a STANDARD
// function: unqualified, or qualified by STANDARD or SYS.STANDARD.
func staticFunctionName(call *nodes.FuncCallExpr) string {
	if call.FuncName == nil {
		return ""
	}
	switch call.FuncName.Schema {
	case "", "STANDARD", "SYS.STANDARD":
		return call.FuncName.Name
	}
	return ""
}

// staticKind is the datatype family of a static expression.
type staticKind int

const (
	kindAny       staticKind = iota // NULL, a constant, or an inquiry directive: the text does not say
	kindBoolean                     // TRUE, FALSE, a comparison, NOT, AND, OR, IS NULL
	kindNumber                      // a numeric literal, arithmetic, a numeric function
	kindCharacter                   // a string literal, ||, a character function
)

// letter is the letter staticSignature spells kind with.
func (k staticKind) letter() string {
	switch k {
	case kindBoolean:
		return "b"
	case kindNumber:
		return "n"
	case kindCharacter:
		return "c"
	}
	return ""
}

// staticSignature is a static function's arity and what it takes and
// returns; args lists each argument's accepted kinds as a string of
// "n" (number), "c" (character), or "b" (Boolean), the last repeated up
// to max arguments. NVL returns its arguments' kind.
type staticSignature struct {
	min, max int
	args     []string
	result   staticKind
}

// staticFunctions are the functions Oracle 23ai lets a static expression
// call, with the signatures it checks (PLS-00306 for ABS(TRUE) or
// ABS(1, 2), PLS-00174 for ABS('1') or TO_CHAR(1, '9')). Any other call,
// a builtin such as INSTR, GREATEST, or TO_NUMBER or a function of the
// user's, is PLS-00174 (or PLS-00201, PLS-00204).
var staticFunctions = map[string]staticSignature{
	"ABS": {1, 1, []string{"n"}, kindNumber}, "CEIL": {1, 1, []string{"n"}, kindNumber},
	"FLOOR": {1, 1, []string{"n"}, kindNumber}, "SQRT": {1, 1, []string{"n"}, kindNumber},
	"SIGN": {1, 1, []string{"n"}, kindNumber}, "EXP": {1, 1, []string{"n"}, kindNumber},
	"LN": {1, 1, []string{"n"}, kindNumber}, "ATAN": {1, 1, []string{"n"}, kindNumber},
	"COS": {1, 1, []string{"n"}, kindNumber}, "SIN": {1, 1, []string{"n"}, kindNumber},
	"TAN": {1, 1, []string{"n"}, kindNumber}, "COSH": {1, 1, []string{"n"}, kindNumber},
	"SINH": {1, 1, []string{"n"}, kindNumber}, "TANH": {1, 1, []string{"n"}, kindNumber},
	"ROUND": {1, 2, []string{"n", "n"}, kindNumber}, "TRUNC": {1, 2, []string{"n", "n"}, kindNumber},
	"POWER": {2, 2, []string{"n", "n"}, kindNumber}, "MOD": {2, 2, []string{"n", "n"}, kindNumber},
	"REMAINDER": {2, 2, []string{"n", "n"}, kindNumber}, "LOG": {2, 2, []string{"n", "n"}, kindNumber},
	"BITAND": {2, 2, []string{"n", "n"}, kindNumber},
	"LENGTH": {1, 1, []string{"nc"}, kindNumber}, "LENGTHB": {1, 1, []string{"nc"}, kindNumber},
	"LENGTHC": {1, 1, []string{"nc"}, kindNumber}, "ASCII": {1, 1, []string{"nc"}, kindNumber},
	"SUBSTR": {2, 3, []string{"nc", "n", "n"}, kindCharacter}, "SUBSTRB": {2, 3, []string{"nc", "n", "n"}, kindCharacter},
	"SUBSTRC": {2, 3, []string{"nc", "n", "n"}, kindCharacter},
	"LOWER":   {1, 1, []string{"nc"}, kindCharacter}, "UPPER": {1, 1, []string{"nc"}, kindCharacter},
	"INITCAP": {1, 1, []string{"nc"}, kindCharacter}, "CHR": {1, 1, []string{"n"}, kindCharacter},
	"CONCAT": {2, 2, []string{"nc", "nc"}, kindCharacter},
	"RPAD":   {2, 3, []string{"nc", "n", "nc"}, kindCharacter}, "LPAD": {2, 3, []string{"nc", "n", "nc"}, kindCharacter},
	"TO_CHAR": {1, 1, []string{"nc"}, kindCharacter}, "TO_NCHAR": {1, 1, []string{"nc"}, kindCharacter},
	"NVL": {2, 2, []string{"nbc", "nbc"}, kindAny}, "XOR": {2, 2, []string{"b", "b"}, kindBoolean},
}

// staticKindOf classifies expr as Oracle 23ai's static expressions are
// typed, returning its kind and, when expr is no static expression, the
// node at fault. The grammar takes literals, NULL, TRUE and FALSE,
// package constants and inquiry directives, arithmetic and || , comparisons of
// numbers or of Booleans (no character values: PLS-00174 for 'a' = 'a'
// and $$PLSQL_UNIT = 'P'), NOT, AND, OR, IS [ NOT ] NULL, and the static
// functions; IN, BETWEEN, LIKE, EXISTS, CASE, DECODE, CAST, a
// pseudo-column, and a subquery are PLS-00174.
func (p *Parser) staticKindOf(expr nodes.ExprNode) (staticKind, nodes.ExprNode) {
	switch e := expr.(type) {
	case *nodes.ParenExpr:
		return p.staticKindOf(e.Expr)
	case *nodes.NullLiteral:
		return kindAny, nil
	case *nodes.NumberLiteral:
		return kindNumber, nil
	case *nodes.StringLiteral:
		return kindCharacter, nil
	case *nodes.ColumnRef:
		if e.Schema == "" && e.Table == "" && (e.Column == "TRUE" || e.Column == "FALSE") &&
			e.Loc.Start >= 0 && e.Loc.Start < len(p.source) && p.source[e.Loc.Start] != '"' {
			return kindBoolean, nil
		}
		if p.isInquiryRef(e) {
			switch e.Column {
			case "$$PLSQL_UNIT", "$$PLSQL_UNIT_OWNER", "$$PLSQL_UNIT_TYPE":
				return kindCharacter, nil
			}
			return kindAny, nil
		}
		// A constant is a package's, named through it (PLS-00201 on Oracle
		// 23ai for $IF debug, a local constant or one of the unit's own
		// package specification among them).
		if e.Table == "" {
			return kindAny, expr
		}
		return kindAny, nil
	case *nodes.UnaryExpr:
		kind, bad := p.staticKindOf(e.Operand)
		if bad != nil {
			return kind, bad
		}
		switch e.Op {
		case "-", "+":
			if kind != kindNumber && kind != kindAny {
				return kind, expr
			}
			return kindNumber, nil
		case "NOT":
			if kind != kindBoolean && kind != kindAny {
				return kind, expr
			}
			return kindBoolean, nil
		}
		return kind, expr
	case *nodes.BoolExpr:
		if e.Args != nil {
			for _, item := range e.Args.Items {
				kind, bad := p.staticKindOf(item.(nodes.ExprNode))
				if bad != nil {
					return kind, bad
				}
				if kind != kindBoolean && kind != kindAny {
					return kind, item.(nodes.ExprNode)
				}
			}
		}
		return kindBoolean, nil
	case *nodes.IsExpr:
		if _, bad := p.staticKindOf(e.Expr); bad != nil {
			return kindBoolean, bad
		}
		if e.Test != "NULL" {
			return kindBoolean, expr
		}
		return kindBoolean, nil
	case *nodes.BinaryExpr:
		left, bad := p.staticKindOf(e.Left)
		if bad != nil {
			return left, bad
		}
		right, bad := p.staticKindOf(e.Right)
		if bad != nil {
			return right, bad
		}
		switch e.Op {
		case "+", "-", "*", "/", "**":
			if left != kindNumber && left != kindAny || right != kindNumber && right != kindAny {
				return kindNumber, expr
			}
			return kindNumber, nil
		case "||":
			if left == kindBoolean || right == kindBoolean {
				return kindCharacter, expr
			}
			return kindCharacter, nil
		case "=", "!=", "<>", "~=", "^=", "<", ">", "<=", ">=":
			if left == kindCharacter || right == kindCharacter ||
				left != kindAny && right != kindAny && left != right {
				return kindBoolean, expr
			}
			return kindBoolean, nil
		}
		return kindAny, expr
	case *nodes.FuncCallExpr:
		sig, ok := staticFunctions[staticFunctionName(e)]
		n := 0
		if e.Args != nil {
			n = e.Args.Len()
		}
		if !ok || n < sig.min || n > sig.max {
			return kindAny, expr
		}
		result := sig.result
		for i := 0; i < n; i++ {
			kind, bad := p.staticKindOf(e.Args.Items[i].(nodes.ExprNode))
			if bad != nil {
				return kind, bad
			}
			accepted := sig.args[len(sig.args)-1]
			if i < len(sig.args) {
				accepted = sig.args[i]
			}
			if kind != kindAny && !strings.Contains(accepted, kind.letter()) {
				return kindAny, expr
			}
			// NVL takes and returns one kind.
			if sig.result == kindAny && kind != kindAny {
				if result != kindAny && result != kind {
					return kindAny, expr
				}
				result = kind
			}
		}
		return result, nil
	}
	return kindAny, expr
}

// isInquiryRef reports whether expr is an inquiry directive, $$name, as
// opposed to a quoted column whose name starts with $$ ("$$X").
func (p *Parser) isInquiryRef(expr nodes.ExprNode) bool {
	ref, ok := expr.(*nodes.ColumnRef)
	return ok && ref.Schema == "" && ref.Table == "" && ref.Loc.Start >= 0 &&
		ref.Loc.Start < len(p.source) && p.source[ref.Loc.Start] == '$'
}

// staticConditionFault returns the node at fault when expr is no static
// Boolean expression: one staticKindOf rejects, or a number or a string
// (PLS-00382 for $IF 1 or $IF 'x'), or nil.
func (p *Parser) staticConditionFault(expr nodes.ExprNode) nodes.ExprNode {
	kind, bad := p.staticKindOf(expr)
	if bad != nil {
		return bad
	}
	if kind != kindBoolean && kind != kindAny {
		return expr
	}
	return nil
}

// checkDirectiveCondition parses the condition at loc of an $IF or $ELSIF
// the lexer evaluated: an expression (PLS-00103) that is a static Boolean
// expression (staticConditionFault).
func (p *Parser) checkDirectiveCondition(loc nodes.Loc) error {
	sub := &Parser{lexer: NewLexerRange(p.source, loc.Start, loc.End), source: p.source}
	sub.advance()
	expr, err := sub.parseExpr()
	if err != nil {
		return err
	}
	if expr == nil || sub.cur.Type != tokEOF {
		return sub.syntaxErrorAtCur()
	}
	if bad := p.staticConditionFault(expr); bad != nil {
		return p.syntaxErrorAtNode(bad)
	}
	return nil
}

// firstDirectiveTooEarly reports whether the unit stmt's first directive,
// at loc, stands where Oracle 23ai rejects it with ORA-00911: right after
// the name of a stored unit, before any token of the unit's own, or before
// the DECLARE, BEGIN, or COMPOUND TRIGGER of a trigger's body.
func (p *Parser) firstDirectiveTooEarly(stmt nodes.StmtNode, loc int) bool {
	var name *nodes.ObjectName
	switch s := stmt.(type) {
	case *nodes.CreateProcedureStmt:
		name = s.Name
	case *nodes.CreateFunctionStmt:
		name = s.Name
	case *nodes.CreatePackageStmt:
		name = s.Name
	case *nodes.CreateTypeStmt:
		name = s.Name
	case *nodes.CreateTriggerStmt:
		return loc < p.triggerBodyStart(s)
	}
	if name == nil || name.Loc.End <= 0 {
		return false
	}
	// The first token after the name is the directive itself.
	return NewLexerRange(p.source, name.Loc.End, len(p.source)).nextRaw().Loc == loc
}

// triggerBodyStart returns where the body of trigger stmt starts: its
// DECLARE or BEGIN, or the COMPOUND of COMPOUND TRIGGER.
func (p *Parser) triggerBodyStart(stmt *nodes.CreateTriggerStmt) int {
	if stmt.Compound {
		// COMPOUND TRIGGER, not a trigger, table, or column named COMPOUND.
		lex := NewLexerRange(p.source, stmt.Loc.Start, stmt.Loc.End)
		prev := Token{}
		for tok := lex.nextRaw(); tok.Type != tokEOF; tok = lex.nextRaw() {
			if prev.Type != tokQIDENT && prev.Str == "COMPOUND" && tok.Type == kwTRIGGER {
				return prev.Loc
			}
			prev = tok
		}
	}
	if stmt.Body == nil {
		return stmt.Loc.End
	}
	return nodes.NodeLoc(stmt.Body).Start
}
