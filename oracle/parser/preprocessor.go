package parser

import (
	"math"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"unicode"
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
// which the parser rejects and Split lexes past, with the nesting error for
// one nested too deeply: the evaluators walk its tree recursively.
func staticCondition(source string, loc nodes.Loc) (ccValue, bool, *ParseError) {
	sub := &Parser{lexer: NewLexerRange(source, loc.Start, loc.End), source: source}
	sub.advance()
	expr, err := sub.parseExpr()
	if err != nil || expr == nil || sub.cur.Type != tokEOF {
		return ccUnknown, false, nestingError(err)
	}
	if treeTooHigh(expr, loc.End-loc.Start) {
		return ccUnknown, false, nestedTooDeeply(loc.Start)
	}
	if sub.staticConditionFault(expr) != nil {
		return ccUnknown, false, nil
	}
	switch newStaticEval(source).staticTruthOf(expr) {
	case truthTrue:
		return ccTrue, true, nil
	case truthFalse, truthNull:
		return ccFalse, true, nil
	}
	return ccUnknown, true, nil
}

// staticTruth is the three-valued truth of a static condition, or unknown.
type staticTruth int

const (
	truthUnknown staticTruth = iota
	truthTrue
	truthFalse
	truthNull
)

// staticEval evaluates the static expressions of one condition, each node
// once. A node's value, type, and nullness follow from its operands', and
// staticRaise asks about every node: computing them anew for each made a
// condition of 401 summed literals take 2.8 s, and nested NVL calls, which
// ask several questions of each argument, exponential time. It keeps no
// character value, only its length and nullness (textInfo), so nested
// character functions over a long value hold one copy at a time.
type staticEval struct {
	source   string
	truths   map[nodes.ExprNode]staticTruth
	nulls    map[nodes.ExprNode]staticNullness
	numbers  map[nodes.ExprNode]staticValue
	numTypes map[nodes.ExprNode]numType
	texts    map[nodes.ExprNode]textInfo
	empties  map[nodes.ExprNode]bool
}

// staticNullness is what staticNull finds of a node.
type staticNullness struct{ null, ok bool }

// staticValue is what staticNumber finds of a node.
type staticValue struct {
	value *big.Rat
	ok    bool
}

// textInfo is what staticText finds of a node, but for the value itself:
// its length in bytes, whether it is NULL, and whether the text settles it.
type textInfo struct {
	length   int
	null, ok bool
}

// newStaticEval returns an evaluator for a condition in source.
func newStaticEval(source string) *staticEval {
	return &staticEval{
		source:   source,
		truths:   map[nodes.ExprNode]staticTruth{},
		nulls:    map[nodes.ExprNode]staticNullness{},
		numbers:  map[nodes.ExprNode]staticValue{},
		numTypes: map[nodes.ExprNode]numType{},
		texts:    map[nodes.ExprNode]textInfo{},
		empties:  map[nodes.ExprNode]bool{},
	}
}

// staticTruthOf evaluates expr as a static Boolean expression.
func (ev *staticEval) staticTruthOf(expr nodes.ExprNode) staticTruth {
	t, seen := ev.truths[expr]
	if !seen {
		t = ev.truthOf(expr)
		ev.truths[expr] = t
	}
	return t
}

// truthOf computes staticTruthOf.
func (ev *staticEval) truthOf(expr nodes.ExprNode) staticTruth {
	switch e := expr.(type) {
	case *nodes.ParenExpr:
		return ev.staticTruthOf(e.Expr)
	case *nodes.NullLiteral:
		return truthNull
	case *nodes.ColumnRef:
		// TRUE and FALSE, not a quoted "TRUE".
		if e.Schema != "" || e.Table != "" || e.Loc.Start < 0 || e.Loc.Start >= len(ev.source) || ev.source[e.Loc.Start] == '"' {
			return truthUnknown
		}
		switch e.Column {
		case "TRUE":
			return truthTrue
		case "FALSE":
			return truthFalse
		}
	case *nodes.BoolExpr:
		return ev.staticBoolExpr(e)
	case *nodes.FuncCallExpr:
		// STANDARD's NVL and XOR, the Boolean functions a static
		// expression may call.
		if !isStaticBooleanFunction(e) {
			return truthUnknown
		}
		a := ev.staticTruthOf(e.Args.Items[0].(nodes.ExprNode))
		b := ev.staticTruthOf(e.Args.Items[1].(nodes.ExprNode))
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
		null, ok := ev.staticNull(e.Expr)
		switch {
		case !ok:
			return truthUnknown
		case null != e.Not:
			return truthTrue
		}
		return truthFalse
	case *nodes.BinaryExpr:
		left, lnull, lok := ev.staticComparable(e.Left)
		right, rnull, rok := ev.staticComparable(e.Right)
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
func (ev *staticEval) staticComparable(expr nodes.ExprNode) (value *big.Rat, null bool, ok bool) {
	if v, ok := ev.staticNumber(expr); ok {
		return v, v == nil, true
	}
	switch ev.staticTruthOf(expr) {
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
func (ev *staticEval) staticNull(expr nodes.ExprNode) (null bool, ok bool) {
	r, seen := ev.nulls[expr]
	if !seen {
		r.null, r.ok = ev.nullOf(expr)
		ev.nulls[expr] = r
	}
	return r.null, r.ok
}

// nullOf computes staticNull.
func (ev *staticEval) nullOf(expr nodes.ExprNode) (null bool, ok bool) {
	if _, null, ok := ev.staticComparable(expr); ok {
		return null, true
	}
	if info := ev.textInfo(expr); info.ok {
		return info.null, true
	}
	switch e := expr.(type) {
	case *nodes.ParenExpr:
		return ev.staticNull(e.Expr)
	case *nodes.StringLiteral:
		return e.Val == "", true
	case *nodes.BinaryExpr:
		if e.Op == "||" {
			return ev.staticConcatNull(e.Left, e.Right)
		}
	case *nodes.FuncCallExpr:
		if e.Args == nil {
			return false, false
		}
		args := e.Args.Items
		switch name := staticFunctionName(e); {
		case name == "CONCAT" && len(args) == 2:
			return ev.staticConcatNull(args[0].(nodes.ExprNode), args[1].(nodes.ExprNode))
		case name == "NVL" && len(args) == 2:
			first, firstOK := ev.staticNull(args[0].(nodes.ExprNode))
			second, secondOK := ev.staticNull(args[1].(nodes.ExprNode))
			switch {
			case firstOK && !first:
				return false, true
			case firstOK || secondOK && !second:
				return second, secondOK
			}
		case len(args) == 1 && (name == "UPPER" || name == "LOWER" || name == "INITCAP" ||
			name == "TO_CHAR" || name == "TO_NCHAR" || name == "CHR"):
			return ev.staticNull(args[0].(nodes.ExprNode))
		case (name == "RPAD" || name == "LPAD") && (len(args) == 2 || len(args) == 3):
			// NULL when an argument is or the length is under 1.
			nonNull := true
			for i, arg := range args {
				null, ok := ev.staticNull(arg.(nodes.ExprNode))
				if ok && null {
					return true, true
				}
				if i == 1 {
					if n, ok := ev.staticNumber(arg.(nodes.ExprNode)); ok && n != nil && roundRat(n, "ROUND").Sign() <= 0 {
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

// staticText returns the value of a static character expression, null
// when it is NULL, and false for ok when the text does not settle it. It
// settles ASCII values alone, which every character set and SUBSTR or
// LENGTH variant counts alike: a string literal, a PLS_INTEGER literal
// converted, ||, CONCAT, NVL, UPPER, LOWER, INITCAP (a word starts after
// any character but a letter or digit), TO_CHAR, TO_NCHAR, CHR, SUBSTR
// (position 0 counts as 1 and a negative one from the end), and RPAD and
// LPAD (which cut the text to their length); numeric arguments round to
// whole numbers, as Oracle 23ai rounds them (SUBSTR('abc', 3.9) is NULL).
func (ev *staticEval) staticText(expr nodes.ExprNode) (text string, null bool, ok bool) {
	text, null, ok = ev.textOf(expr)
	ev.texts[expr] = textInfo{length: len(text), null: null, ok: ok}
	return text, null, ok
}

// textInfo returns what staticText finds of expr but the value: what it
// found before, or what the operands make of || and CONCAT, which join
// them, and of UPPER, LOWER, INITCAP, TO_CHAR, and TO_NCHAR, which keep
// their argument's length and nullness, without building any value.
func (ev *staticEval) textInfo(expr nodes.ExprNode) textInfo {
	if info, seen := ev.texts[expr]; seen {
		return info
	}
	var info textInfo
	if left, right, ok := concatOperands(expr); ok {
		l, r := ev.textInfo(left), ev.textInfo(right)
		if l.ok && r.ok && l.length+r.length <= 32767 {
			info = textInfo{length: l.length + r.length, null: l.null && r.null, ok: true}
		}
	} else if arg := sameLengthArg(expr); arg != nil {
		info = ev.textInfo(arg)
	} else {
		ev.staticText(expr)
		return ev.texts[expr]
	}
	ev.texts[expr] = info
	return info
}

// sameLengthArg returns the argument of UPPER, LOWER, INITCAP, TO_CHAR, or
// TO_NCHAR of one, whose value has the length and nullness of the call's.
func sameLengthArg(expr nodes.ExprNode) nodes.ExprNode {
	call, ok := unparen(expr).(*nodes.FuncCallExpr)
	if !ok || call.Args == nil || call.Args.Len() != 1 {
		return nil
	}
	switch staticFunctionName(call) {
	case "UPPER", "LOWER", "INITCAP", "TO_CHAR", "TO_NCHAR":
		return call.Args.Items[0].(nodes.ExprNode)
	}
	return nil
}

// textOf computes staticText.
func (ev *staticEval) textOf(expr nodes.ExprNode) (text string, null bool, ok bool) {
	switch e := expr.(type) {
	case *nodes.ParenExpr:
		return ev.staticText(e.Expr)
	case *nodes.NullLiteral:
		return "", true, true
	case *nodes.StringLiteral:
		for i := 0; i < len(e.Val); i++ {
			if e.Val[i] >= utf8.RuneSelf {
				return "", false, false
			}
		}
		return e.Val, e.Val == "", true
	case *nodes.NumberLiteral, *nodes.UnaryExpr, *nodes.BinaryExpr:
		if b, ok := e.(*nodes.BinaryExpr); ok && b.Op == "||" {
			return ev.concatText(b)
		}
		if ev.staticNumType(e) != numPLSInteger {
			break
		}
		if v, ok := ev.staticNumber(e); ok && v != nil && plsInteger(v) {
			return strconv.FormatInt(v.Num().Int64(), 10), false, true
		}
	case *nodes.FuncCallExpr:
		if ev.staticNumType(e) == numPLSInteger {
			if v, ok := ev.staticNumber(e); ok && v != nil && plsInteger(v) {
				return strconv.FormatInt(v.Num().Int64(), 10), false, true
			}
		}
		if e.Args == nil {
			break
		}
		var args []nodes.ExprNode
		for _, item := range e.Args.Items {
			args = append(args, item.(nodes.ExprNode))
		}
		switch name := staticFunctionName(e); {
		case name == "CONCAT" && len(args) == 2:
			return ev.concatText(e)
		case name == "NVL" && len(args) == 2:
			first, null, ok := ev.staticText(args[0])
			if !ok || !null {
				return first, null, ok
			}
			return ev.staticText(args[1])
		case (name == "UPPER" || name == "LOWER" || name == "INITCAP" || name == "TO_CHAR" || name == "TO_NCHAR") && len(args) == 1:
			text, null, ok := ev.staticText(args[0])
			if !ok || null {
				return "", null, ok
			}
			switch name {
			case "UPPER":
				text = strings.ToUpper(text)
			case "LOWER":
				text = strings.ToLower(text)
			case "INITCAP":
				out := []byte(text)
				for i, c := range out {
					word := i > 0 && isASCIIAlnum(out[i-1])
					if word {
						out[i] = byte(unicode.ToLower(rune(c)))
					} else {
						out[i] = byte(unicode.ToUpper(rune(c)))
					}
				}
				text = string(out)
			}
			return text, false, true
		case name == "CHR" && len(args) == 1:
			n, null, ok := ev.staticWhole(args[0])
			if ok && !null && (n < 1 || n >= utf8.RuneSelf) {
				return "", false, false
			}
			return string(rune(n)), null, ok
		case (name == "SUBSTR" || name == "SUBSTRB" || name == "SUBSTRC") && (len(args) == 2 || len(args) == 3):
			text, null, ok := ev.staticText(args[0])
			if !ok || null {
				return "", null, ok
			}
			var whole []int64
			for _, arg := range args[1:] {
				n, null, ok := ev.staticWhole(arg)
				if !ok || null {
					return "", null, ok
				}
				whole = append(whole, n)
			}
			length, start := int64(len(text)), whole[0]-1
			switch {
			case whole[0] == 0:
				start = 0
			case whole[0] < 0:
				start = length + whole[0]
			}
			if start < 0 || start >= length || len(whole) == 2 && whole[1] < 1 {
				return "", true, true
			}
			end := length
			if len(whole) == 2 && start+whole[1] < end {
				end = start + whole[1]
			}
			return text[start:end], false, true
		case (name == "RPAD" || name == "LPAD") && (len(args) == 2 || len(args) == 3):
			text, null, ok := ev.staticText(args[0])
			if !ok || null {
				return "", null, ok
			}
			n, null, ok := ev.staticWhole(args[1])
			switch {
			case !ok || null:
				return "", null, ok
			case n < 1:
				return "", true, true
			case n > 32767:
				return "", false, false
			}
			pad := " "
			if len(args) == 3 {
				if pad, null, ok = ev.staticText(args[2]); !ok || null {
					return "", null, ok
				}
			}
			if int64(len(text)) >= n {
				return text[:n], false, true
			}
			fill := strings.Repeat(pad, int(n)/len(pad)+1)[:n-int64(len(text))]
			if name == "RPAD" {
				return text + fill, false, true
			}
			return fill + text, false, true
		}
	}
	return "", false, false
}

// concatText returns the value of a || or CONCAT expression: NULL when
// every operand is. Past 32767 bytes it raises (staticRaises) and is no
// value. The operands of a chain of them join in one pass: building each
// prefix of 'a' || 'b' || ... in turn copied the text once per operand.
func (ev *staticEval) concatText(expr nodes.ExprNode) (text string, null bool, ok bool) {
	var b strings.Builder
	null = true
	stack := []nodes.ExprNode{expr}
	for len(stack) > 0 {
		e := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if left, right, ok := concatOperands(e); ok {
			stack = append(stack, right, left)
			continue
		}
		part, partNull, ok := ev.staticText(e)
		if !ok || b.Len()+len(part) > 32767 {
			return "", false, false
		}
		b.WriteString(part)
		null = null && partNull
	}
	return b.String(), null, true
}

// concatOperands returns the operands of expr when it is || or CONCAT of
// two, parenthesized or not.
func concatOperands(expr nodes.ExprNode) (left, right nodes.ExprNode, ok bool) {
	switch e := unparen(expr).(type) {
	case *nodes.BinaryExpr:
		if e.Op == "||" {
			return e.Left, e.Right, true
		}
	case *nodes.FuncCallExpr:
		if staticFunctionName(e) == "CONCAT" && e.Args != nil && e.Args.Len() == 2 {
			return e.Args.Items[0].(nodes.ExprNode), e.Args.Items[1].(nodes.ExprNode), true
		}
	}
	return nil, nil, false
}

// staticWhole returns a static number rounded to a whole PLS_INTEGER, as
// a character function's numeric argument rounds.
func (ev *staticEval) staticWhole(expr nodes.ExprNode) (n int64, null bool, ok bool) {
	v, ok := ev.staticNumber(expr)
	switch {
	case !ok:
		return 0, false, false
	case v == nil:
		return 0, true, true
	}
	w := roundRat(v, "ROUND")
	if !plsInteger(w) {
		return 0, false, false
	}
	return w.Num().Int64(), false, true
}

// numType is the static datatype of a numeric expression.
type numType int

const (
	numUnknown    numType = iota // NULL, a constant, or an inquiry directive
	numPLSInteger                // the PLS_INTEGER a character context converts
	numNumber                    // NUMBER, which it does not
)

// staticNumType types a static numeric expression as Oracle 23ai does: an
// integer literal up to 2147483647 is a PLS_INTEGER, and so are +, -, and
// * of PLS_INTEGERs, ABS, SIGN, BITAND, and NVL of them, and LENGTH and
// ASCII; any other literal, /, **, MOD, ROUND, and the other functions are
// NUMBER (TO_CHAR(1 + 1) and TO_CHAR(LENGTH('ab')) compile, TO_CHAR(4 / 2),
// TO_CHAR(MOD(5, 2)), and TO_CHAR(-2147483648) raise).
func (ev *staticEval) staticNumType(expr nodes.ExprNode) numType {
	t, seen := ev.numTypes[expr]
	if !seen {
		t = ev.numTypeOf(expr)
		ev.numTypes[expr] = t
	}
	return t
}

// numTypeOf computes staticNumType.
func (ev *staticEval) numTypeOf(expr nodes.ExprNode) numType {
	switch e := expr.(type) {
	case *nodes.ParenExpr:
		return ev.staticNumType(e.Expr)
	case *nodes.NumberLiteral:
		if _, ok := plsIntegerLiteral(e); ok {
			return numPLSInteger
		}
		return numNumber
	case *nodes.UnaryExpr:
		if e.Op == "-" || e.Op == "+" {
			return ev.staticNumType(e.Operand)
		}
	case *nodes.BinaryExpr:
		switch e.Op {
		case "+", "-", "*":
			return joinNumTypes(ev.staticNumType(e.Left), ev.staticNumType(e.Right))
		case "/", "**":
			return numNumber
		}
	case *nodes.FuncCallExpr:
		if e.Args == nil {
			return numUnknown
		}
		args := e.Args.Items
		switch name := staticFunctionName(e); name {
		case "LENGTH", "LENGTHB", "LENGTHC", "ASCII":
			return numPLSInteger
		case "ABS", "SIGN", "BITAND", "NVL":
			t := ev.staticNumType(args[0].(nodes.ExprNode))
			for _, arg := range args[1:] {
				t = joinNumTypes(t, ev.staticNumType(arg.(nodes.ExprNode)))
			}
			return t
		default:
			if _, ok := staticFunctions[name]; ok && staticFunctions[name].result == kindNumber {
				return numNumber
			}
		}
	}
	return numUnknown
}

// joinNumTypes types an operation over operands of types a and b.
func joinNumTypes(a, b numType) numType {
	switch {
	case a == numNumber || b == numNumber:
		return numNumber
	case a == numPLSInteger && b == numPLSInteger:
		return numPLSInteger
	}
	return numUnknown
}

// plsIntegerLiteral returns the value of an integer literal, possibly
// signed, that is a PLS_INTEGER (Oracle 23ai raises for TO_CHAR(0.5),
// TO_CHAR(1.0), TO_CHAR(2147483648), and TO_CHAR(-2147483648)).
func plsIntegerLiteral(expr nodes.ExprNode) (int64, bool) {
	switch e := expr.(type) {
	case *nodes.ParenExpr:
		return plsIntegerLiteral(e.Expr)
	case *nodes.NumberLiteral:
		if e.IsFloat || e.Ival > 1<<31-1 || strings.ContainsAny(e.Val, ".eE") {
			return 0, false
		}
		return e.Ival, true
	case *nodes.UnaryExpr:
		v, ok := plsIntegerLiteral(e.Operand)
		switch {
		case !ok:
			return 0, false
		case e.Op == "-":
			return -v, true
		case e.Op == "+":
			return v, true
		}
	}
	return 0, false
}

// isASCIIAlnum reports whether c is an ASCII letter or digit.
func isASCIIAlnum(c byte) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9'
}

// staticConcatNull reports whether left || right is NULL: when both are.
func (ev *staticEval) staticConcatNull(left, right nodes.ExprNode) (null bool, ok bool) {
	lnull, lok := ev.staticNull(left)
	rnull, rok := ev.staticNull(right)
	switch {
	case lok && !lnull, rok && !rnull:
		return false, true
	case lok && rok:
		return true, true
	}
	return false, false
}

// staticBoolExpr evaluates NOT, AND, and OR in three-valued logic.
func (ev *staticEval) staticBoolExpr(e *nodes.BoolExpr) staticTruth {
	if e.Args == nil || e.Args.Len() == 0 {
		return truthUnknown
	}
	if e.Boolop == nodes.BOOL_NOT {
		switch ev.staticTruthOf(e.Args.Items[0].(nodes.ExprNode)) {
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
		switch v := ev.staticTruthOf(item.(nodes.ExprNode)); v {
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
func (ev *staticEval) staticNumber(expr nodes.ExprNode) (*big.Rat, bool) {
	r, seen := ev.numbers[expr]
	if !seen {
		r.value, r.ok = ev.numberOf(expr)
		ev.numbers[expr] = r
	}
	if r.value == nil {
		return nil, r.ok
	}
	// A copy, as the callers compute in place.
	return new(big.Rat).Set(r.value), r.ok
}

// numberOf computes staticNumber.
func (ev *staticEval) numberOf(expr nodes.ExprNode) (*big.Rat, bool) {
	switch e := expr.(type) {
	case *nodes.ParenExpr:
		return ev.staticNumber(e.Expr)
	case *nodes.NullLiteral:
		return nil, true
	case *nodes.NumberLiteral:
		return staticLiteral(e.Val)
	case *nodes.UnaryExpr:
		v, ok := ev.staticNumber(e.Operand)
		switch {
		case !ok || e.Op != "-" && e.Op != "+":
			return nil, false
		case v == nil || e.Op == "+":
			return v, true
		}
		return v.Neg(v), true
	case *nodes.FuncCallExpr:
		return ev.staticNumericCall(e)
	case *nodes.BinaryExpr:
		switch e.Op {
		case "+", "-", "*", "/", "**":
		default:
			return nil, false
		}
		left, lok := ev.staticNumber(e.Left)
		right, rok := ev.staticNumber(e.Right)
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
func (ev *staticEval) staticNumericCall(call *nodes.FuncCallExpr) (*big.Rat, bool) {
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
		if null, ok := ev.staticNull(arg); ok && null {
			return nil, true
		}
		if lit, ok := arg.(*nodes.StringLiteral); ok && (name == "LENGTH" || name == "LENGTHC") {
			return big.NewRat(int64(utf8.RuneCountInString(lit.Val)), 1), true
		}
		info := ev.textInfo(arg)
		switch {
		case !info.ok:
			return nil, false
		case info.null:
			return nil, true
		case name == "ASCII":
			text, _, _ := ev.staticText(arg)
			return big.NewRat(int64(text[0]), 1), true
		}
		return big.NewRat(int64(info.length), 1), true
	}
	if name == "NVL" && call.Args.Len() == 2 {
		v, ok := ev.staticNumber(call.Args.Items[0].(nodes.ExprNode))
		if !ok || v != nil {
			return v, ok
		}
		return ev.staticNumber(call.Args.Items[1].(nodes.ExprNode))
	}
	var args []*big.Rat
	null := false
	for _, item := range call.Args.Items {
		v, ok := ev.staticNumber(item.(nodes.ExprNode))
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

// staticPower raises base to a whole exponent (Oracle 23ai finds
// POWER(4, -2) = 0.0625, 1 ** 65 = 1, and 2 ** 100 exact) when NUMBER
// holds the exact power; past that, Oracle's rounding of the steps is not
// the text's to settle.
func staticPower(base, exp *big.Rat) (*big.Rat, bool) {
	if !exp.IsInt() || exp.Sign() < 0 && base.Sign() == 0 {
		return nil, false
	}
	abs := new(big.Rat).Abs(base)
	odd := exp.Num().Bit(0) == 1
	switch {
	case base.Sign() == 0:
		if exp.Sign() == 0 {
			return big.NewRat(1, 1), true
		}
		return new(big.Rat), true
	case abs.Cmp(big.NewRat(1, 1)) == 0:
		if base.Sign() < 0 && odd {
			return big.NewRat(-1, 1), true
		}
		return big.NewRat(1, 1), true
	}
	// Any other base outgrows NUMBER's 40 digits or its range long before
	// the bound: 2 ** 4096 is far past 1e126.
	n := new(big.Int).Abs(exp.Num())
	if n.Cmp(big.NewInt(4096)) > 0 {
		return nil, false
	}
	num := new(big.Int).Exp(base.Num(), n, nil)
	den := new(big.Int).Exp(base.Denom(), n, nil)
	result := new(big.Rat).SetFrac(num, den)
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

// pow10 returns 10^e, which its callers only read: from a table for the
// exponents NUMBER's rounding and ROUND's scales use, as computing one per
// rounding was most of a long condition's cost.
func pow10(e int) *big.Rat {
	pow10Once.Do(func() {
		ten, p := big.NewInt(10), big.NewInt(1)
		for i := 0; i <= pow10Bound; i++ {
			pow10Table[pow10Bound+i] = new(big.Rat).SetInt(p)
			pow10Table[pow10Bound-i] = new(big.Rat).SetFrac(big.NewInt(1), p)
			p = new(big.Int).Mul(p, ten)
		}
	})
	if e >= -pow10Bound && e <= pow10Bound {
		return pow10Table[pow10Bound+e]
	}
	if e < 0 {
		return new(big.Rat).Inv(pow10(-e))
	}
	return new(big.Rat).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(e)), nil))
}

// pow10Bound bounds the exponents pow10Table holds, past NUMBER's range
// (1e-130 to 1e126) and ROUND's scales (-200 to 200) both.
const pow10Bound = 400

var (
	pow10Once  sync.Once
	pow10Table [2*pow10Bound + 1]*big.Rat
)

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
	return newStaticEval(p.source).staticRaise(expr)
}

// staticRaise returns the operation in expr that certainly raises an
// exception, which makes the condition no static expression (PLS-00174 on
// Oracle 23ai for $IF FALSE AND 1 / 0 = 0, short-circuit or not).
func (ev *staticEval) staticRaise(expr nodes.ExprNode) nodes.ExprNode {
	var bad nodes.ExprNode
	nodes.Inspect(expr, func(n nodes.Node) bool {
		if e, ok := n.(nodes.ExprNode); ok && bad == nil && ev.staticRaises(e) {
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
// RPAD or LPAD length, past PLS_INTEGER, and a pad past 32767; an empty
// string as an operand (staticEmpty); and a numeric literal past NUMBER's
// range (numberOverflows; PLS-00569 for $IF 1e126 > 0). EXP, SINH, and
// COSH past NUMBER's range raise nothing.
func (ev *staticEval) staticRaises(e nodes.ExprNode) bool {
	var name string
	var args []nodes.ExprNode
	switch e := e.(type) {
	case *nodes.NumberLiteral:
		return numberOverflows(e.Val)
	case *nodes.BinaryExpr:
		switch e.Op {
		case "/":
			name = "REMAINDER"
		case "**":
			name = "POWER"
		case "||", "+", "-", "*":
			name = e.Op
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
			if ev.staticEmpty(arg) {
				return true
			}
		}
	}
	// A character context converts a PLS_INTEGER alone (Oracle 23ai raises
	// for TO_CHAR(0.5), TO_CHAR(4 / 2), 1e3 || 'a', and UPPER(2147483648)).
	var converted []nodes.ExprNode
	switch name {
	case "TO_CHAR", "TO_NCHAR", "||", "CONCAT", "UPPER", "LOWER", "INITCAP":
		converted = args
	case "SUBSTR", "SUBSTRB", "SUBSTRC":
		converted = args[:1]
	case "RPAD", "LPAD":
		converted = args[:1]
		if len(args) == 3 {
			converted = append(converted, args[2])
		}
	}
	for _, arg := range converted {
		if ev.staticNumType(arg) == numNumber {
			return true
		}
	}
	// PLS_INTEGER arithmetic stays in 32 bits (Oracle 23ai raises for
	// 2147483647 + 1).
	// A static character value holds 32767 bytes at most (Oracle 23ai raises
	// for RPAD('a', 32767) || 'b').
	if (name == "||" || name == "CONCAT") && len(args) == 2 {
		l, r := ev.textInfo(args[0]), ev.textInfo(args[1])
		if l.ok && r.ok && l.length+r.length > 32767 {
			return true
		}
	}
	if (name == "+" || name == "-" || name == "*") && ev.staticNumType(args[0]) == numPLSInteger && ev.staticNumType(args[1]) == numPLSInteger {
		a, aok := ev.staticNumber(args[0])
		b, bok := ev.staticNumber(args[1])
		if aok && bok && a != nil && b != nil {
			exact := new(big.Rat).Mul(a, b)
			switch name {
			case "+":
				exact = new(big.Rat).Add(a, b)
			case "-":
				exact = new(big.Rat).Sub(a, b)
			}
			return !plsInteger(exact)
		}
	}
	// The arguments' values, when the text settles each and none is NULL.
	values := make([]*big.Rat, len(args))
	for i, arg := range args {
		if v, ok := ev.staticNumber(arg); ok && v != nil {
			values[i] = v
		} else if null, ok := ev.staticNull(arg); !ok || null {
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

// unparen returns expr without its parentheses.
func unparen(expr nodes.ExprNode) nodes.ExprNode {
	for paren, ok := expr.(*nodes.ParenExpr); ok; paren, ok = expr.(*nodes.ParenExpr) {
		expr = paren.Expr
	}
	return expr
}

// staticEmpty reports whether expr is an empty string: an empty literal,
// plain, q-quoted, or national, or TO_CHAR, TO_NCHAR, or NVL that yield
// one.
func (ev *staticEval) staticEmpty(expr nodes.ExprNode) bool {
	empty, seen := ev.empties[expr]
	if !seen {
		empty = ev.emptyOf(expr)
		ev.empties[expr] = empty
	}
	return empty
}

// emptyOf computes staticEmpty.
func (ev *staticEval) emptyOf(expr nodes.ExprNode) bool {
	switch e := expr.(type) {
	case *nodes.ParenExpr:
		return ev.staticEmpty(e.Expr)
	case *nodes.StringLiteral:
		return e.Val == ""
	case *nodes.FuncCallExpr:
		if e.Args == nil {
			return false
		}
		args := e.Args.Items
		switch staticFunctionName(e) {
		case "TO_CHAR", "TO_NCHAR":
			return len(args) == 1 && ev.staticEmpty(args[0].(nodes.ExprNode))
		case "NVL":
			if len(args) != 2 {
				return false
			}
			null, ok := ev.staticNull(args[0].(nodes.ExprNode))
			return ok && null && ev.staticEmpty(args[1].(nodes.ExprNode))
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
	if treeTooHigh(expr, loc.End-loc.Start) {
		return nestedTooDeeply(loc.Start)
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
	// Before or in the name, or the first token after it is the directive
	// itself (ORA-00911 for PROCEDURE $IF TRUE $THEN p IS $END and for
	// PROCEDURE s.$IF TRUE $THEN p IS $END).
	return loc < name.Loc.End || NewLexerRange(p.source, name.Loc.End, len(p.source)).nextRaw().Loc == loc
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
