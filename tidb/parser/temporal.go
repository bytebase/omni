package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// The reader for the temporal functions' arguments. Most of these are not
// ordinary calls: each has its own production in the grammar (pingcap/tidb
// pkg/parser/parser.y):
//
//	FunctionNameDatetimePrecision FuncDatetimePrec         (parser.y:8176)
//	FuncDatetimePrec: %empty | '(' ')' | '(' intLit ')'    (parser.y:8707)
//	builtinCurTime '(' FuncDatetimePrecListOpt ')'         (parser.y:8231)
//	builtinSysDate '(' FuncDatetimePrecListOpt ')'         (parser.y:8235)
//	FunctionNameOptionalBraces OptionalBraces              (parser.y:8168)
//	builtinCurDate '(' ')'                                 (parser.y:8172)
//	FunctionNameDateArith '(' Expression ',' "INTERVAL" Expression TimeUnit ')' (parser.y:8261)
//	FunctionNameDateArithMultiForms '(' Expression ',' Expression ')'           (parser.y:8239)
//	FunctionNameDateArithMultiForms '(' Expression ',' "INTERVAL" Expression TimeUnit ')' (parser.y:8250)
//
// So a precision is an unsigned integer literal and nothing else: TiDB v8.5.5
// answers 1064 for curtime(-1), current_time(1.0), utc_time('1'),
// current_timestamp(NULL) and curdate(1). NOW is the odd one out, because the
// grammar lists it in FunctionNameConflict (parser.y:8123) and so reads an
// ExpressionListOpt: now(1.0) and now('1') run and now(1, 2) answers 1582.
//
// A few of these answer 1064 from past the parser, with no position in the message:
// now(-1), curtime(2147483648), and `curtime`(-1), whose quoted name is not even the
// keyword. The grammar took the text, so Parse does too. Every row is held to the
// engine by TestTiDBTemporalFormsMatchEngine, which tells those two 1064s apart.

// checkSpacedFuncWord applies the lexer's rule for a function word with a space
// or a comment before its '(' (misc.go, isTokenIdentifier), where only a word
// followed directly by '(' takes its btFuncTokenMap token. p.prev is the word and
// p.cur its '('. It reports whether the call is an ordinary generic one.
//
// The aggregate and datetime builtins reach the grammar as a per-name token the
// lexer folds out of `name(` (SumExpr's builtinCount/builtinSum/..., parser.y:8461;
// builtinCount '(' '*' ')', parser.y:8546; the temporal productions listed above).
// With a space the word is no longer that token, so none of those special forms
// (the '*' and DISTINCT/ALL quantifiers, the FROM/SEPARATOR/INTERVAL clauses, a
// precision) is reachable, and how the spaced word then reads splits in two:
//
//   - A word that is also a column-nameable keyword reverts to that keyword, and
//     FunctionCallGeneric (identifier '(' ExpressionListOpt ')', parser.y:8684)
//     does not take a keyword as the name, so spaced it is a syntax error whatever
//     follows. TiDB v8.5.5 answers 1064 for now (), date_add (a, 10), sum ( a ),
//     max ( distinct a ), group_concat ( a ), bit_and ( a ), substring ( a from 1 )
//     and trim ( leading 'x' from a ) — the ordinary arguments included.
//   - A word that is only the builtin token (not a keyword) reverts to a plain
//     identifier, so spaced it is an ordinary generic call. TiDB runs count ( a ),
//     count ( a, a ), mid ( a ) and sysdate (b), and answers 1064 for count ( * )
//     and count ( distinct a ) because '*' and the quantifiers are no Expression
//     in a generic argument list.
//
// AVG and JSON_ARRAYAGG/JSON_OBJECTAGG are keyword literals in SumExpr
// ("AVG" '(' ..., parser.y:8462), not btFuncTokenMap tokens, so a space before the
// '(' changes nothing and they stay aggregate calls: avg ( distinct a ) runs. MID
// has no builtin production in omni, so it is already an ordinary call either way;
// ADDDATE and SUBDATE are only in tokenMap, so spacing changes nothing for them.
func (p *Parser) checkSpacedFuncWord(upperName string) (ordinary bool, err error) {
	if p.prev.End == p.cur.Loc {
		return false, nil
	}
	switch upperName {
	// Only the builtin token, so spaced the word is a plain identifier and the
	// call an ordinary generic one.
	case "COUNT", "SYSDATE":
		return true, nil
	// Also a column-nameable keyword, so spaced the word cannot name a generic
	// call: a syntax error whatever follows.
	case "NOW", "CURDATE", "CURTIME", "DATE_ADD", "DATE_SUB",
		"SUM", "MIN", "MAX", "GROUP_CONCAT",
		"STD", "STDDEV", "STDDEV_POP", "STDDEV_SAMP",
		"VAR_POP", "VAR_SAMP", "VARIANCE",
		"BIT_AND", "BIT_OR", "BIT_XOR",
		"APPROX_COUNT_DISTINCT", "APPROX_PERCENTILE",
		"SUBSTR", "SUBSTRING", "TRIM", "POSITION":
		return false, p.syntaxErrorAtCur()
	}
	return false, nil
}

// atDatetimePrecision reports whether p.cur can be a function's fractional
// seconds precision: an intLit, which the lexer gives up to 18446744073709551615.
// Past it the literal is a decimal, which the grammar refuses: TiDB v8.5.5 answers
// 1064 with a position for curtime(18446744073709551616). Up to it the grammar takes
// the literal, including the ones refused from past the parser (see the file comment).
func (p *Parser) atDatetimePrecision() bool {
	return p.cur.Type == tokICONST
}

// parseTemporalKeyword parses a reserved word that names a function and needs no
// parentheses. CURRENT_DATE, UTC_DATE and CURRENT_USER take at most an empty pair
// of them, and the others a precision as well. The word is reserved, so a space
// before the '(' changes nothing, unlike the function words parseFuncCall reads.
func (p *Parser) parseTemporalKeyword() (nodes.ExprNode, error) {
	tok := p.advance()
	fc := &nodes.FuncCallExpr{Loc: nodes.Loc{Start: tok.Loc}, Name: strings.ToUpper(tok.Str)}
	if p.cur.Type == '(' {
		fc.HasParens = true
		p.advance()
		var err error
		switch tok.Type {
		case kwCURRENT_DATE, kwUTC_DATE, kwCURRENT_USER:
			_, err = p.expect(')')
		default:
			err = p.parseDatetimePrecision(fc)
		}
		if err != nil {
			return nil, err
		}
	}
	fc.Loc.End = p.prev.End
	return fc, nil
}

// parseDatetimePrecision reads an optional precision and the closing ')', the '('
// already consumed.
func (p *Parser) parseDatetimePrecision(fc *nodes.FuncCallExpr) error {
	if p.cur.Type != ')' {
		if !p.atDatetimePrecision() {
			return p.syntaxErrorAtCur()
		}
		arg, err := p.parsePrimaryExpr()
		if err != nil {
			return err
		}
		fc.Args = append(fc.Args, arg)
	}
	_, err := p.expect(')')
	return err
}

// parseTemporalCall parses the rest of a call to one of the temporal builtins with
// a production of its own, the '(' already consumed. None of those productions
// has a window. It reports false, reading nothing, for any other name.
func (p *Parser) parseTemporalCall(fc *nodes.FuncCallExpr) (bool, error) {
	var err error
	switch fc.Name {
	case "CURTIME", "SYSDATE":
		err = p.parseDatetimePrecision(fc)
	case "CURDATE":
		_, err = p.expect(')')
	case "DATE_ADD", "DATE_SUB":
		err = p.parseDateArgs(fc, true)
	case "ADDDATE", "SUBDATE":
		err = p.parseDateArgs(fc, false)
	default:
		return false, nil
	}
	if err != nil {
		return true, err
	}
	fc.Loc.End = p.prev.End
	return true, nil
}

// isBareInterval reports whether expr is an INTERVAL literal standing on its own,
// parenthesized or not. The literal is not an expression: every engine in this
// family reads it only as an operand of + and -, and as the second argument of the
// date arithmetic builtins. So abs(INTERVAL 1 DAY) and abs((INTERVAL 1 DAY)) are
// syntax errors, where abs(INTERVAL 1 DAY + 1) parses.
func isBareInterval(expr nodes.ExprNode) bool {
	for {
		switch e := expr.(type) {
		case *nodes.IntervalExpr:
			return true
		case *nodes.ParenExpr:
			expr = e.Expr
		default:
			return false
		}
	}
}

// parseDateArgs reads the two arguments of DATE_ADD, DATE_SUB, ADDDATE or SUBDATE
// and the closing ')'. The first is an expression. The second of DATE_ADD and
// DATE_SUB is INTERVAL value unit, written directly and with nothing after it, so
// date_add(a, 10), date_add(a, (INTERVAL 1 DAY)) and date_add(a, INTERVAL 1 DAY + 1)
// are all syntax errors. ADDDATE and SUBDATE take that or any expression, which is
// why adddate(a, 10) and adddate(a, INTERVAL 1 DAY + 1) parse, and a parenthesized
// INTERVAL literal is no expression there either.
func (p *Parser) parseDateArgs(fc *nodes.FuncCallExpr, intervalOnly bool) error {
	first, err := p.parseExpr()
	if err != nil {
		return err
	}
	if first == nil || isStarArg(first) || isBareInterval(first) {
		return p.syntaxErrorAtCur()
	}
	if _, err := p.expect(','); err != nil {
		return err
	}

	var second nodes.ExprNode
	if intervalOnly {
		p.checkCursor()
		if p.collectMode() {
			p.addTokenCandidate(kwINTERVAL)
			return &ParseError{Message: "collecting"}
		}
		if p.cur.Type != kwINTERVAL {
			return p.syntaxErrorAtCur()
		}
		second, err = p.parseIntervalExpr()
		if err != nil {
			return err
		}
		// The comma form, INTERVAL(a, b), is the function and not the literal.
		if _, literal := second.(*nodes.IntervalExpr); !literal {
			return p.syntaxErrorAtCur()
		}
	} else {
		written := p.cur.Type == kwINTERVAL
		second, err = p.parseExpr()
		if err != nil {
			return err
		}
		if second == nil || isStarArg(second) || (!written && isBareInterval(second)) {
			return p.syntaxErrorAtCur()
		}
	}

	fc.Args = append(fc.Args, first, second)
	_, err = p.expect(')')
	return err
}

// parseODBCEscape parses the ODBC escape { Identifier Expression } (SimpleExpr,
// parser.y:7941, v8.5.5). The action compares the name as written, so only d, t
// and ts, in lower case and backquoted or not, make a DATE, TIME or TIMESTAMP
// literal; {D '2020-01-01'} and {abc 1} are the expression alone. Verified on
// TiDB v8.5.5: a view stores {d '2020-01-01'} as DATE '2020-01-01' typed date,
// and {D '2020-01-01'} as _UTF8MB4'2020-01-01' typed var_string.
//
// The expression alone keeps its own Loc, so an unaliased {abc 1} is named 1
// where TiDB names it after the whole field, {abc 1}.
func (p *Parser) parseODBCEscape() (nodes.ExprNode, error) {
	start := p.cur.Loc
	p.advance() // consume {
	if !p.isIdentToken() {
		return nil, &ParseError{Message: "syntax error, expected an identifier after {", Position: p.cur.Loc}
	}
	nameTok := p.advance()
	exprStart := p.pos()
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	// The argument is an Expression, which takes neither the bare DEFAULT nor
	// the * that parseExpr also reads: {d *} and {abc DEFAULT} are 1064s on
	// v8.5.5.
	if err := checkRowItem(expr, exprStart); err != nil {
		return nil, err
	}
	if _, err := p.expect('}'); err != nil {
		return nil, err
	}
	var typ string
	if nameTok.Type == tokIDENT {
		switch nameTok.Str {
		case "d":
			typ = "DATE"
		case "t":
			typ = "TIME"
		case "ts":
			typ = "TIMESTAMP"
		}
	}
	if typ == "" {
		return expr, nil
	}
	lit := &nodes.TemporalLit{Loc: nodes.Loc{Start: start, End: p.prev.End}, Type: typ}
	if s, ok := expr.(*nodes.StringLit); ok {
		lit.Value = s.Value
	} else {
		lit.Expr = expr
	}
	return lit, nil
}
