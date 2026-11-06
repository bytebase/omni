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

// checkSpacedTemporalWord applies the lexer's rule for a function word with a
// space or a comment before its '(' (misc.go, isTokenIdentifier), where only a
// word followed directly by '(' takes its btFuncTokenMap token. p.prev is the word
// and p.cur its '('. It reports whether the call is an ordinary one.
//
// NOW, CURDATE, CURTIME, DATE_ADD and DATE_SUB are also in tokenMap, under tokens
// that can name a column, and the name of an ordinary call is the plain identifier
// token (FunctionCallGeneric, parser.y:8684). So spaced they are a syntax error
// whatever follows: TiDB v8.5.5 answers 1064 for now (), curdate (1) and
// date_add (a, 10), with a position. SYSDATE is in no other map, so it is an
// identifier and the call an ordinary one: sysdate (b) answers 1054. ADDDATE and
// SUBDATE are only in tokenMap, under the tokens the grammar reads, so spacing changes
// nothing for them.
func (p *Parser) checkSpacedTemporalWord(upperName string) (ordinary bool, err error) {
	if p.prev.End == p.cur.Loc {
		return false, nil
	}
	switch upperName {
	case "NOW", "CURDATE", "CURTIME", "DATE_ADD", "DATE_SUB":
		return false, p.syntaxErrorAtCur()
	case "SYSDATE":
		return true, nil
	}
	return false, nil
}

// maxPrecisionLiteral is the largest literal the lexer gives as intLit. Past it the
// literal is no longer an intLit, which the grammar refuses: TiDB v8.5.5 answers 1064
// with a position for curtime(18446744073709551616). Up to it the grammar takes the
// literal, including the ones refused from past the parser (see the file comment).
const maxPrecisionLiteral = "18446744073709551615"

// atDatetimePrecision reports whether p.cur can be a function's fractional
// seconds precision.
func (p *Parser) atDatetimePrecision() bool {
	return p.cur.Type == tokICONST && intLiteralFits(p.cur.Str, maxPrecisionLiteral)
}

// intLiteralFits reports whether a run of decimal digits, leading zeros aside, is
// no larger than max.
func intLiteralFits(digits, max string) bool {
	digits = strings.TrimLeft(digits, "0")
	if len(digits) != len(max) {
		return len(digits) < len(max)
	}
	return digits <= max
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
