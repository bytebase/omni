package parser

import (
	"math/big"

	nodes "github.com/bytebase/omni/oracle/ast"
)

// staticCondition evaluates the conditional compilation condition at
// source[loc.Start:loc.End] as far as the text settles it. Oracle 23ai
// folds the Boolean operators over TRUE, FALSE, and NULL, short-circuiting
// past what it cannot settle (TRUE OR $$x is TRUE), and compares numeric
// literals and arithmetic over them in decimal (0.1 + 0.2 = 0.3); inquiry
// directives and package constants have values the text does not hold.
//
// It reports false for a condition that is no expression, which the
// parser rejects and Split lexes past.
func staticCondition(source string, loc nodes.Loc) (ccValue, bool) {
	sub := &Parser{lexer: NewLexerRange(source, loc.Start, loc.End), source: source}
	sub.advance()
	expr, err := sub.parseExpr()
	if err != nil || expr == nil || sub.cur.Type != tokEOF {
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
	case *nodes.BinaryExpr:
		left, lok := staticNumber(e.Left)
		right, rok := staticNumber(e.Right)
		if !lok || !rok {
			return truthUnknown
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

// staticNumber returns the value of a numeric literal, possibly signed, or
// of arithmetic over them (Oracle 23ai selects the $ELSE of $IF 1 + 1 = 3).
// Values are exact decimals, as Oracle's NUMBER holds them; a quotient
// that no decimal holds exactly (1 / 3) is left unknown, since NUMBER
// rounds it (Oracle 23ai finds 1 / 3 * 3 = 1 FALSE).
func staticNumber(expr nodes.ExprNode) (*big.Rat, bool) {
	switch e := expr.(type) {
	case *nodes.ParenExpr:
		return staticNumber(e.Expr)
	case *nodes.NumberLiteral:
		return new(big.Rat).SetString(e.Val)
	case *nodes.UnaryExpr:
		v, ok := staticNumber(e.Operand)
		switch {
		case !ok:
			return nil, false
		case e.Op == "-":
			return v.Neg(v), true
		case e.Op == "+":
			return v, true
		}
	case *nodes.FuncCallExpr:
		return staticNumericCall(e)
	case *nodes.BinaryExpr:
		left, lok := staticNumber(e.Left)
		right, rok := staticNumber(e.Right)
		if !lok || !rok {
			return nil, false
		}
		switch e.Op {
		case "+":
			return left.Add(left, right), true
		case "-":
			return left.Sub(left, right), true
		case "*":
			return left.Mul(left, right), true
		case "/":
			if right.Sign() == 0 {
				return nil, false
			}
			q := left.Quo(left, right)
			return q, terminatingDecimal(q)
		case "**":
			// A small whole power, exact.
			if !right.IsInt() || right.Sign() < 0 || right.Num().Cmp(big.NewInt(64)) > 0 {
				return nil, false
			}
			result := big.NewRat(1, 1)
			for i := right.Num().Int64(); i > 0; i-- {
				result.Mul(result, left)
			}
			return result, true
		}
	}
	return nil, false
}

// staticNumericCall folds the static numeric functions whose exact value
// the text settles: ABS, SIGN, CEIL, FLOOR, TRUNC and ROUND of one
// argument, MOD (Oracle's, with the dividend's sign), and POWER to a whole
// exponent (Oracle 23ai selects the $ELSE of $IF ABS(-1) = 2).
func staticNumericCall(call *nodes.FuncCallExpr) (*big.Rat, bool) {
	if call.FuncName == nil || call.Args == nil {
		return nil, false
	}
	switch call.FuncName.Schema {
	case "", "STANDARD", "SYS.STANDARD":
	default:
		return nil, false
	}
	var args []*big.Rat
	for _, item := range call.Args.Items {
		v, ok := staticNumber(item.(nodes.ExprNode))
		if !ok {
			return nil, false
		}
		args = append(args, v)
	}
	one := len(args) == 1
	switch name := call.FuncName.Name; {
	case name == "ABS" && one:
		return args[0].Abs(args[0]), true
	case name == "SIGN" && one:
		return big.NewRat(int64(args[0].Sign()), 1), true
	case name == "CEIL" && one, name == "FLOOR" && one, name == "TRUNC" && one, name == "ROUND" && one:
		return roundRat(args[0], name), true
	case name == "MOD" && len(args) == 2:
		if args[1].Sign() == 0 {
			return args[0], true
		}
		q := roundRat(new(big.Rat).Quo(args[0], args[1]), "TRUNC")
		return args[0].Sub(args[0], q.Mul(q, args[1])), true
	case name == "POWER" && len(args) == 2:
		exp := args[1]
		if !exp.IsInt() || exp.Sign() < 0 || exp.Num().Cmp(big.NewInt(64)) > 0 {
			return nil, false
		}
		result := big.NewRat(1, 1)
		for i := exp.Num().Int64(); i > 0; i-- {
			result.Mul(result, args[0])
		}
		return result, true
	}
	return nil, false
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

// terminatingDecimal reports whether a decimal holds r exactly: its reduced
// denominator has no prime factor but 2 and 5.
func terminatingDecimal(r *big.Rat) bool {
	d := new(big.Int).Set(r.Denom())
	for _, f := range []int64{2, 5} {
		factor, rem := big.NewInt(f), new(big.Int)
		for {
			q, m := new(big.Int).QuoRem(d, factor, rem)
			if m.Sign() != 0 {
				break
			}
			d = q
		}
	}
	return d.Cmp(big.NewInt(1)) == 0
}

// nonBooleanOperand returns the first operand in a Boolean position of a
// condition, the condition itself or an operand of NOT, AND, or OR, that
// cannot be a static Boolean whatever its operands hold, or nil. Oracle 23ai
// rejects a number or a string (PLS-00382 for $IF 'x', $IF 1 + 2, $IF $$x +
// 1), and a call to anything but the STANDARD functions NVL and XOR: a
// numeric or character one (PLS-00382 for $IF ABS(1)) or one that is not
// static (PLS-00174 for a Boolean function of the user's), and CASE and
// DECODE (PLS-00174, PLS-00306).
func nonBooleanOperand(expr nodes.ExprNode) nodes.ExprNode {
	switch e := expr.(type) {
	case *nodes.ParenExpr:
		return nonBooleanOperand(e.Expr)
	case *nodes.NumberLiteral, *nodes.StringLiteral, *nodes.CaseExpr, *nodes.DecodeExpr:
		return expr
	case *nodes.UnaryExpr:
		if e.Op == "-" || e.Op == "+" {
			return expr
		}
	case *nodes.BinaryExpr:
		switch e.Op {
		case "+", "-", "*", "/", "**", "||":
			return expr
		}
		// A static comparison compares no character values (PLS-00174 for
		// 'a' = 'a', $$x = 'a', and $$PLSQL_UNIT = 'P').
		if characterOperand(e.Left) || characterOperand(e.Right) {
			return expr
		}
	case *nodes.FuncCallExpr:
		if !isStaticBooleanFunction(e) {
			return expr
		}
		// NVL and XOR are Boolean over Boolean operands (PLS-00382 for
		// NVL(1, 2), PLS-00306 for NVL(TRUE, 1)).
		for _, item := range e.Args.Items {
			if bad := nonBooleanOperand(item.(nodes.ExprNode)); bad != nil {
				return bad
			}
		}
	case *nodes.BoolExpr:
		if e.Args != nil {
			for _, item := range e.Args.Items {
				if bad := nonBooleanOperand(item.(nodes.ExprNode)); bad != nil {
					return bad
				}
			}
		}
	}
	return nil
}

// isStaticBooleanFunction reports whether call is STANDARD's NVL or XOR,
// with two arguments: the Boolean functions a static expression may call.
func isStaticBooleanFunction(call *nodes.FuncCallExpr) bool {
	if call.FuncName == nil || call.Args == nil || call.Args.Len() != 2 {
		return false
	}
	switch call.FuncName.Schema {
	case "", "STANDARD", "SYS.STANDARD":
	default:
		return false
	}
	return call.FuncName.Name == "NVL" || call.FuncName.Name == "XOR"
}

// characterOperand reports whether expr is character data in a static
// expression: a string literal, a concatenation, or the inquiry directive
// $$PLSQL_UNIT, $$PLSQL_UNIT_OWNER, or $$PLSQL_UNIT_TYPE.
func characterOperand(expr nodes.ExprNode) bool {
	switch e := expr.(type) {
	case *nodes.ParenExpr:
		return characterOperand(e.Expr)
	case *nodes.StringLiteral:
		return true
	case *nodes.BinaryExpr:
		return e.Op == "||"
	case *nodes.ColumnRef:
		switch e.Column {
		case "$$PLSQL_UNIT", "$$PLSQL_UNIT_OWNER", "$$PLSQL_UNIT_TYPE":
			return true
		}
	}
	return false
}

// staticFunctions are the functions Oracle 23ai lets a static expression
// call; any other call, a builtin such as INSTR, GREATEST, or TO_NUMBER or
// a function of the user's, is PLS-00174 (or PLS-00201, PLS-00204).
var staticFunctions = map[string]bool{
	"ABS": true, "CEIL": true, "FLOOR": true, "ROUND": true, "TRUNC": true, "POWER": true,
	"SQRT": true, "MOD": true, "REMAINDER": true, "SIGN": true, "EXP": true, "LN": true,
	"LOG": true, "BITAND": true, "ATAN": true, "COS": true, "SIN": true, "TAN": true,
	"COSH": true, "SINH": true, "TANH": true, "LENGTH": true, "LENGTHB": true, "LENGTHC": true,
	"SUBSTR": true, "SUBSTRB": true, "SUBSTRC": true, "LOWER": true, "UPPER": true,
	"INITCAP": true, "ASCII": true, "CHR": true, "CONCAT": true, "RPAD": true, "LPAD": true,
	"TO_CHAR": true, "TO_NCHAR": true, "NVL": true, "XOR": true,
}

// nonStaticCall returns the first call in expr that no static expression
// may make, or nil: a function outside staticFunctions, a pseudo-column
// (SYSDATE, USER), DECODE, CASE, or CAST (PLS-00174 on Oracle 23ai).
func nonStaticCall(expr nodes.ExprNode) nodes.ExprNode {
	var bad nodes.ExprNode
	nodes.Inspect(expr, func(n nodes.Node) bool {
		if bad != nil {
			return false
		}
		switch e := n.(type) {
		case *nodes.PseudoColumn, *nodes.DecodeExpr, *nodes.CaseExpr, *nodes.CastExpr:
			bad = e.(nodes.ExprNode)
			return false
		}
		if call, ok := n.(*nodes.FuncCallExpr); ok {
			name := ""
			if call.FuncName != nil {
				switch call.FuncName.Schema {
				case "", "STANDARD", "SYS.STANDARD":
					name = call.FuncName.Name
				}
			}
			if !staticFunctions[name] {
				bad = call
				return false
			}
		}
		return true
	})
	return bad
}

// isInquiryRef reports whether expr is an inquiry directive, $$name, as
// opposed to a quoted column whose name starts with $$ ("$$X").
func (p *Parser) isInquiryRef(expr nodes.ExprNode) bool {
	ref, ok := expr.(*nodes.ColumnRef)
	return ok && ref.Schema == "" && ref.Table == "" && ref.Loc.Start >= 0 &&
		ref.Loc.Start < len(p.source) && p.source[ref.Loc.Start] == '$'
}

// checkDirectiveCondition parses the condition at loc of an $IF or $ELSIF
// the lexer evaluated: an expression (PLS-00103) whose Boolean operands can
// be static Booleans (nonBooleanOperand).
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
	if bad := nonBooleanOperand(expr); bad != nil {
		return p.syntaxErrorAtNode(bad)
	}
	if bad := nonStaticCall(expr); bad != nil {
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
