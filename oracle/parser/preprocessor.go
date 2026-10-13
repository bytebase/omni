package parser

import (
	"math"
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
		null, ok := staticNull(source, e.Expr)
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

// staticNull reports whether a static expression is NULL, and false for
// ok when the text does not settle it. Beyond numbers and Booleans, Oracle
// 23ai settles it for a string literal (NULL when empty), || and CONCAT
// (NULL when both operands are), NVL, and the functions that keep their
// argument's nullness: UPPER, LOWER, INITCAP, TO_CHAR, TO_NCHAR, and CHR
// ('x' IS NULL and 'a' || NULL IS NULL are FALSE).
func staticNull(source string, expr nodes.ExprNode) (null bool, ok bool) {
	if _, null, ok := staticComparable(source, expr); ok {
		return null, true
	}
	switch e := expr.(type) {
	case *nodes.ParenExpr:
		return staticNull(source, e.Expr)
	case *nodes.StringLiteral:
		return e.Val == "", true
	case *nodes.BinaryExpr:
		if e.Op == "||" {
			return staticConcatNull(source, e.Left, e.Right)
		}
	case *nodes.FuncCallExpr:
		if e.Args == nil {
			return false, false
		}
		args := e.Args.Items
		switch name := staticFunctionName(e); {
		case name == "CONCAT" && len(args) == 2:
			return staticConcatNull(source, args[0].(nodes.ExprNode), args[1].(nodes.ExprNode))
		case name == "NVL" && len(args) == 2:
			first, firstOK := staticNull(source, args[0].(nodes.ExprNode))
			second, secondOK := staticNull(source, args[1].(nodes.ExprNode))
			switch {
			case firstOK && !first:
				return false, true
			case firstOK || secondOK && !second:
				return second, secondOK
			}
		case len(args) == 1 && (name == "UPPER" || name == "LOWER" || name == "INITCAP" ||
			name == "TO_CHAR" || name == "TO_NCHAR" || name == "CHR"):
			return staticNull(source, args[0].(nodes.ExprNode))
		case (name == "SUBSTR" || name == "SUBSTRB" || name == "SUBSTRC") && (len(args) == 2 || len(args) == 3):
			return staticSubstrNull(args)
		case (name == "RPAD" || name == "LPAD") && (len(args) == 2 || len(args) == 3):
			// NULL when an argument is or the length is under 1.
			nonNull := true
			for i, arg := range args {
				null, ok := staticNull(source, arg.(nodes.ExprNode))
				if ok && null {
					return true, true
				}
				if i == 1 {
					if n, ok := staticNumber(arg.(nodes.ExprNode)); ok && n != nil && roundRat(n, "ROUND").Sign() <= 0 {
						return true, true
					}
				}
				nonNull = nonNull && ok
			}
			if nonNull {
				return false, true
			}
		}
	}
	return false, false
}

// staticSubstrNull reports whether SUBSTR(text, position [, length]) is
// NULL, as Oracle 23ai finds it: position 0 counts as 1 and a negative one
// from the end, each argument rounds to a whole number (SUBSTR('abc', 3.9)
// is NULL), and the result is NULL when it starts past either end or the
// length is under 1. The text is an ASCII literal or an integer, which
// every character set and SUBSTR variant counts alike.
func staticSubstrNull(args []nodes.Node) (null bool, ok bool) {
	var text string
	switch v := args[0].(type) {
	case *nodes.StringLiteral:
		text = v.Val
	default:
		n, ok := staticNumber(v.(nodes.ExprNode))
		if !ok || n == nil || !n.IsInt() {
			if ok && n == nil {
				return true, true
			}
			return false, false
		}
		text = n.Num().String()
	}
	for i := 0; i < len(text); i++ {
		if text[i] >= utf8.RuneSelf {
			return false, false
		}
	}
	var whole []int64
	for _, arg := range args[1:] {
		v, ok := staticNumber(arg.(nodes.ExprNode))
		switch {
		case !ok:
			return false, false
		case v == nil:
			return true, true
		}
		n := roundRat(v, "ROUND")
		if !plsInteger(n) {
			return false, false
		}
		whole = append(whole, n.Num().Int64())
	}
	length := int64(len(text))
	start := whole[0] - 1
	switch {
	case whole[0] == 0:
		start = 0
	case whole[0] < 0:
		start = length + whole[0]
	}
	if start < 0 || start >= length || len(whole) == 2 && whole[1] < 1 {
		return true, true
	}
	return false, true
}

// staticConcatNull reports whether left || right is NULL: when both are.
func staticConcatNull(source string, left, right nodes.ExprNode) (null bool, ok bool) {
	lnull, lok := staticNull(source, left)
	rnull, rok := staticNull(source, right)
	switch {
	case lok && !lnull, rok && !rnull:
		return false, true
	case lok && rok:
		return true, true
	}
	return false, false
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
	// LENGTH, LENGTHB, LENGTHC, and ASCII are NULL of NULL. LENGTH and
	// LENGTHC of a string literal count its characters, and LENGTHB and
	// ASCII read an ASCII one alike in every character set.
	if (name == "LENGTH" || name == "LENGTHB" || name == "LENGTHC" || name == "ASCII") && call.Args.Len() == 1 {
		arg := call.Args.Items[0].(nodes.ExprNode)
		if null, ok := staticNull("", arg); ok && null {
			return nil, true
		}
		lit, ok := arg.(*nodes.StringLiteral)
		if !ok {
			return nil, false
		}
		ascii := true
		for i := 0; i < len(lit.Val); i++ {
			ascii = ascii && lit.Val[i] < utf8.RuneSelf
		}
		switch {
		case name == "LENGTH" || name == "LENGTHC":
			return big.NewRat(int64(utf8.RuneCountInString(lit.Val)), 1), true
		case !ascii:
			return nil, false
		case name == "LENGTHB":
			return big.NewRat(int64(len(lit.Val)), 1), true
		}
		return big.NewRat(int64(lit.Val[0]), 1), true
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
		name == "MOD" && len(args) == 2, name == "POWER" && len(args) == 2,
		name == "BITAND" && len(args) == 2, name == "REMAINDER" && len(args) == 2,
		name == "SQRT" && one, name == "EXP" && one, name == "LN" && one, name == "COS" && one,
		name == "COSH" && one, name == "SIN" && one, name == "TAN" && one, name == "ATAN" && one,
		name == "SINH" && one, name == "TANH" && one, name == "LOG" && len(args) == 2:
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
	case "REMAINDER":
		// n2 - n1 * N, N the integer nearest n2 / n1, the even one at a
		// tie (Oracle 23ai finds REMAINDER(14, 4) = -2).
		if args[1].Sign() == 0 {
			return nil, false
		}
		q := new(big.Rat).Quo(args[0], args[1])
		n := roundRat(q, "ROUND")
		if diff := new(big.Rat).Sub(q, roundRat(q, "TRUNC")); diff.Abs(diff).Cmp(big.NewRat(1, 2)) == 0 && n.Num().Bit(0) == 1 {
			n.Sub(n, big.NewRat(int64(q.Sign()), 1))
		}
		return oracleNumber(args[0].Sub(args[0], n.Mul(n, args[1])))
	case "SQRT":
		// The exact root of an exact square (Oracle 23ai finds SQRT(2.25) =
		// 1.5); any other root is NUMBER's rounding.
		if args[0].Sign() < 0 {
			return nil, false
		}
		num, den := new(big.Int).Sqrt(args[0].Num()), new(big.Int).Sqrt(args[0].Denom())
		root := new(big.Rat).SetFrac(num, den)
		if new(big.Rat).Mul(root, root).Cmp(args[0]) != 0 {
			return nil, false
		}
		return root, true
	case "EXP", "LN", "COS", "COSH", "SIN", "TAN", "ATAN", "SINH", "TANH", "LOG":
		// The exact values: EXP(0), COS(0), and COSH(0) are 1; LN(1),
		// LOG(b, 1), and the others at 0 are 0.
		x, at := args[0], new(big.Rat)
		switch name {
		case "LN":
			at = big.NewRat(1, 1)
		case "LOG":
			x, at = args[1], big.NewRat(1, 1)
			if args[0].Sign() <= 0 || args[0].Cmp(at) == 0 {
				return nil, false
			}
		}
		if x.Cmp(at) != 0 {
			return nil, false
		}
		if name == "EXP" || name == "COS" || name == "COSH" {
			return big.NewRat(1, 1), true
		}
		return new(big.Rat), true
	case "BITAND":
		// Two's complement over the whole parts, toward zero (Oracle 23ai
		// finds BITAND(2.9, 3) = 2 and BITAND(-6, -3) = -8).
		if !bitandOperand(args[0]) || !bitandOperand(args[1]) {
			return nil, false
		}
		a, b := roundRat(args[0], "TRUNC").Num(), roundRat(args[1], "TRUNC").Num()
		return new(big.Rat).SetInt(new(big.Int).And(a, b)), true
	}
	return staticPower(args[0], args[1])
}

// bitandOperand reports whether BITAND takes v: -2^127 through 2^127 - 1
// (PLS-00174 on Oracle 23ai past them).
func bitandOperand(v *big.Rat) bool {
	limit := new(big.Rat).SetInt(new(big.Int).Lsh(big.NewInt(1), 127))
	return v.Cmp(new(big.Rat).Neg(limit)) >= 0 && v.Cmp(limit) < 0
}

// staticPower raises base to a whole exponent of at most 64 either way
// (Oracle 23ai finds POWER(4, -2) = 0.0625), when NUMBER holds the exact
// power; past that, Oracle's rounding of the steps is not the text's to
// settle.
func staticPower(base, exp *big.Rat) (*big.Rat, bool) {
	if !exp.IsInt() || exp.Num().CmpAbs(big.NewInt(64)) > 0 || exp.Sign() < 0 && base.Sign() == 0 {
		return nil, false
	}
	result := big.NewRat(1, 1)
	for i := new(big.Int).Abs(exp.Num()).Int64(); i > 0; i-- {
		result.Mul(result, base)
	}
	if exp.Sign() < 0 {
		result.Inv(result)
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
		// An inquiry directive's kind is not the text's: PLSQL_CCFLAGS may
		// set even $$PLSQL_LINE or $$PLSQL_UNIT (Oracle 23ai takes
		// $IF $$PLSQL_LINE under PLSQL_LINE:TRUE).
		if p.isInquiryRef(e) {
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

// isInquiryRef reports whether expr is an inquiry directive, $$name,
// possibly parenthesized, as opposed to a quoted column whose name starts
// with $$ ("$$X").
func (p *Parser) isInquiryRef(expr nodes.ExprNode) bool {
	for paren, ok := expr.(*nodes.ParenExpr); ok; paren, ok = expr.(*nodes.ParenExpr) {
		expr = paren.Expr
	}
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
	return staticRaise(expr)
}

// staticRaise returns the operation in expr that certainly raises an
// exception, which makes the condition no static expression (PLS-00174 on
// Oracle 23ai for $IF FALSE AND 1 / 0 = 0, short-circuit or not).
func staticRaise(expr nodes.ExprNode) nodes.ExprNode {
	var bad nodes.ExprNode
	nodes.Inspect(expr, func(n nodes.Node) bool {
		if e, ok := n.(nodes.ExprNode); ok && bad == nil && staticRaises(e) {
			bad = e
		}
		return bad == nil
	})
	return bad
}

// staticRaises reports whether evaluating e raises once the text settles
// its operands, none of them NULL (NULL / 0 is NULL), as Oracle 23ai finds:
// a division or REMAINDER by zero; zero to a negative power, a negative
// number to a fractional one, and a power from 1e126 on; SQRT of a negative
// number, LN of one not positive, and LOG of either; BITAND past -2^127
// through 2^127 - 1; CHR of less than 1; a SUBSTR position or length, or an
// RPAD or LPAD length, past PLS_INTEGER, and a pad past 32767; and an
// empty string as an operand (staticEmpty). EXP, SINH, and COSH past
// NUMBER's range raise nothing.
func staticRaises(e nodes.ExprNode) bool {
	var name string
	var args []nodes.ExprNode
	switch e := e.(type) {
	case *nodes.BinaryExpr:
		switch e.Op {
		case "/":
			name = "REMAINDER"
		case "**":
			name = "POWER"
		case "||":
			name = "||"
		default:
			return false
		}
		args = []nodes.ExprNode{e.Left, e.Right}
	case *nodes.FuncCallExpr:
		if e.Args == nil {
			return false
		}
		name = staticFunctionName(e)
		for _, item := range e.Args.Items {
			args = append(args, item.(nodes.ExprNode))
		}
	default:
		return false
	}
	// An empty string, a NULL of its own, raises as an operand but of
	// TO_CHAR, TO_NCHAR, and NVL (Oracle 23ai: UPPER(''), '' || 'a', and
	// SUBSTR('a', '') raise, TO_CHAR('') IS NULL and UPPER(NULL) do not).
	if name != "TO_CHAR" && name != "TO_NCHAR" && name != "NVL" {
		for _, arg := range args {
			if staticEmpty(arg) {
				return true
			}
		}
	}
	// The arguments' values, when the text settles each and none is NULL.
	values := make([]*big.Rat, len(args))
	for i, arg := range args {
		if v, ok := staticNumber(arg); ok && v != nil {
			values[i] = v
		} else if null, ok := staticNull("", arg); !ok || null {
			return false
		}
	}
	one := big.NewRat(1, 1)
	switch {
	case (name == "REMAINDER") && len(args) == 2 && values[0] != nil && values[1] != nil:
		return values[1].Sign() == 0
	case name == "POWER" && len(args) == 2 && values[0] != nil && values[1] != nil:
		a, b := values[0], values[1]
		return a.Sign() == 0 && b.Sign() < 0 || a.Sign() < 0 && !b.IsInt() || powerOverflows(a, b)
	case name == "SQRT" && len(args) == 1 && values[0] != nil:
		return values[0].Sign() < 0
	case name == "LN" && len(args) == 1 && values[0] != nil:
		return values[0].Sign() <= 0
	case name == "LOG" && len(args) == 2 && values[0] != nil && values[1] != nil:
		return values[0].Sign() <= 0 || values[0].Cmp(one) == 0 || values[1].Sign() <= 0
	case name == "BITAND" && len(args) == 2 && values[0] != nil && values[1] != nil:
		return !bitandOperand(values[0]) || !bitandOperand(values[1])
	case name == "CHR" && len(args) == 1 && values[0] != nil:
		n := roundRat(values[0], "ROUND")
		return n.Sign() <= 0 || !plsInteger(n)
	case (name == "SUBSTR" || name == "SUBSTRB" || name == "SUBSTRC") && len(args) >= 2:
		for _, v := range values[1:] {
			if v != nil && !plsInteger(roundRat(v, "ROUND")) {
				return true
			}
		}
	case (name == "RPAD" || name == "LPAD") && len(args) >= 2 && values[1] != nil:
		n := roundRat(values[1], "ROUND")
		return !plsInteger(n) || n.Cmp(big.NewRat(32767, 1)) > 0
	}
	return false
}

// staticEmpty reports whether expr is an empty string: an empty literal,
// plain, q-quoted, or national, or TO_CHAR, TO_NCHAR, or NVL that yield
// one.
func staticEmpty(expr nodes.ExprNode) bool {
	switch e := expr.(type) {
	case *nodes.ParenExpr:
		return staticEmpty(e.Expr)
	case *nodes.StringLiteral:
		return e.Val == ""
	case *nodes.FuncCallExpr:
		if e.Args == nil {
			return false
		}
		args := e.Args.Items
		switch staticFunctionName(e) {
		case "TO_CHAR", "TO_NCHAR":
			return len(args) == 1 && staticEmpty(args[0].(nodes.ExprNode))
		case "NVL":
			if len(args) != 2 {
				return false
			}
			_, null := args[0].(*nodes.NullLiteral)
			return (null || staticEmpty(args[0].(nodes.ExprNode))) && staticEmpty(args[1].(nodes.ExprNode))
		}
	}
	return false
}

// plsInteger reports whether the whole number n is a PLS_INTEGER.
func plsInteger(n *big.Rat) bool {
	return n.Cmp(big.NewRat(-1<<31, 1)) >= 0 && n.Cmp(big.NewRat(1<<31-1, 1)) <= 0
}

// powerOverflows reports whether a to the whole power b reaches 1e126,
// past NUMBER's range (Oracle 23ai raises for POWER(10, 126) and not for
// POWER(10, 125)); a power the estimate cannot settle is taken not to.
func powerOverflows(a, b *big.Rat) bool {
	if !b.IsInt() || a.Sign() == 0 {
		return false
	}
	exp := new(big.Int).Abs(b.Num())
	if exp.Cmp(big.NewInt(512)) <= 0 {
		// Exactly, as the boundary needs.
		num := new(big.Int).Exp(new(big.Int).Abs(a.Num()), exp, nil)
		den := new(big.Int).Exp(a.Denom(), exp, nil)
		if b.Sign() < 0 {
			num, den = den, num
		}
		return new(big.Rat).SetFrac(num, den).Cmp(pow10(126)) >= 0
	}
	af, _ := new(big.Rat).Abs(a).Float64()
	bf, _ := b.Float64()
	return bf*math.Log10(af) >= 126+1e-9
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
// at loc, stands where Oracle 23ai rejects it with ORA-00911: before or
// right after the name of a stored unit, before any token of the unit's
// own, or before the DECLARE, BEGIN, or COMPOUND TRIGGER of a trigger's
// body.
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
	// Before the name, or the first token after it is the directive itself
	// (ORA-00911 for PROCEDURE $IF TRUE $THEN p IS $END).
	return loc < name.Loc.Start || NewLexerRange(p.source, name.Loc.End, len(p.source)).nextRaw().Loc == loc
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
