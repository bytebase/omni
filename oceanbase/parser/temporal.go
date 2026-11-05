package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// The reader for the temporal functions' arguments. None of these is an
// ordinary call: each has its own production in the grammar (oceanbase
// v4.3.5_CE, src/sql/parser/sql_parser_mysql_mode.y):
//
//	cur_timestamp_func: NOW '(' ')' | NOW '(' INTNUM ')' | now_synonyms_func opt_time_func_fsp_i
//	now_synonyms_func: CURRENT_TIMESTAMP | LOCALTIME | LOCALTIMESTAMP
//	opt_time_func_fsp_i: '(' INTNUM ')' | '(' ')' | %empty
//	cur_time_func: CURTIME '(' ')' | CURTIME '(' INTNUM ')' | CURRENT_TIME opt_time_func_fsp_i
//	SYSDATE '(' ')' | SYSDATE '(' INTNUM ')'
//	UTC_TIME | UTC_TIME '(' ')' | UTC_TIME '(' INTNUM ')', and UTC_TIMESTAMP alike
//	cur_date_func: CURDATE '(' ')' | CURRENT_DATE | CURRENT_DATE '(' ')'
//	DATE_ADD '(' date_params ')'
//	date_params: expr ',' INTERVAL expr date_unit
//	ADDDATE '(' date_params ')' | ADDDATE '(' expr ',' expr ')'
//
// These are keywords wherever they stand, so a space before the '(' changes
// nothing: now (6) and date_add (a, INTERVAL 1 DAY) parse, and now (-1) and
// date_add (a, 10) are syntax errors, where MySQL reads the spaced word as the name
// of an ordinary call. CURRENT_USER and UTC_DATE take at most an empty pair too.
//
// So a precision is an unsigned integer literal and nothing else: OceanBase
// 4.3.5 answers 1064 for now(-1), curtime(1.0), current_time('1'),
// utc_time(NULL) and curdate(1). A precision above the maximum is not a syntax
// error: now(7) answers 1426 after parsing. Every row is held to the engine by
// TestOceanBaseTemporalFormsMatchEngine.

// maxPrecisionLiteral is the largest literal the lexer gives as INTNUM. Past it
// the literal is a decimal, which the grammar refuses: OceanBase 4.3.5 runs
// now(18446744073709551615) and answers 1064 for now(18446744073709551616).
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
	case "NOW", "CURTIME", "SYSDATE":
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
