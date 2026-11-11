package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// The reader for the temporal functions' arguments. None of these is an
// ordinary call: each has its own production in the grammar (MariaDB 11.8,
// sql_yacc.yy, function_call_nonkeyword and function_call_keyword):
//
//	CURTIME, LOCALTIMESTAMP, NOW_SYM, UTC_TIME_SYM, UTC_TIMESTAMP_SYM opt_time_precision
//	SYSDATE | SYSDATE '(' ')' | SYSDATE '(' real_ulong_num ')'
//	opt_time_precision: %empty | '(' ')' | '(' real_ulong_num ')'
//	real_ulong_num: NUM | HEX_NUM | LONG_NUM | ULONGLONG_NUM | dec_num_error
//	CURDATE, UTC_DATE_SYM, CURRENT_USER optional_braces
//	DATE_ADD_INTERVAL '(' expr ',' INTERVAL_SYM expr interval ')'
//	ADDDATE_SYM '(' expr ',' expr ')'
//	ADDDATE_SYM '(' expr ',' INTERVAL_SYM expr interval ')'
//
// So a precision is an unsigned integer literal or a 0x hex one and nothing
// else: mariadb:11.8.8 answers 1064 for now(-1), curtime(1.0),
// current_time('1'), utc_time(NULL) and curdate(1), and parses now(0x1). A
// precision above the maximum is not a syntax error: now(7) answers 1426 after
// parsing. Every row is held to the engine by TestMariaDBTemporalFormsMatchEngine.

// spacedTemporalWords are the words the lexer gives as their function token only
// when '(' follows with nothing between. With a space or a comment first the word
// is an identifier, so the call is an ordinary one to a stored function of that
// name, which takes any expressions and no window: mariadb:11.8.8 answers 1630
// for curtime (-1) and date_add (a, 10), and 1064 for date_add (a, INTERVAL 1 DAY)
// and now (1) OVER (). SYSDATE is not among them: it is a keyword wherever it
// stands, so sysdate (1) runs and sysdate (b) is still a 1064.
var spacedTemporalWords = map[string]bool{
	"ADDDATE": true, "CURDATE": true, "CURTIME": true, "DATE_ADD": true,
	"DATE_SUB": true, "NOW": true, "SUBDATE": true,
}

// spacedFuncWords are the aggregate and scalar builtins whose special argument
// syntax (COUNT(*), the DISTINCT/ALL quantifier, GROUP_CONCAT's SEPARATOR,
// SUBSTRING/TRIM's FROM, POSITION's IN) the lexer reaches only through a token it
// folds out of `name(` when the '(' sits directly after the name. With a space or
// a comment first the word is a plain identifier, so the call is an ordinary
// generic one: it takes expression arguments and none of those special forms.
// mariadb:11.8.8 parses sum (a), count (a, a) and max ( ), and answers 1064 for
// sum (*), max (distinct a), count (*), group_concat (a separator ','),
// substring (a from 1), trim (leading 'x' from a) and position ('a' in b). Every
// one is held to the engine by TestMariaDBSpacedFuncWordsMatchEngine.
//
// AVG, WEIGHT_STRING, GET_FORMAT, TIMESTAMPDIFF and TIMESTAMPADD are
// reserved-keyword productions, not folded tokens, so a space before the '('
// changes nothing and they keep their special forms: avg (distinct a) still runs.
var spacedFuncWords = map[string]bool{
	"COUNT": true,
	"SUM":   true, "MIN": true, "MAX": true, "GROUP_CONCAT": true,
	"STD": true, "STDDEV": true, "STDDEV_POP": true, "STDDEV_SAMP": true,
	"VAR_POP": true, "VAR_SAMP": true, "VARIANCE": true,
	"BIT_AND": true, "BIT_OR": true, "BIT_XOR": true,
	"JSON_ARRAYAGG": true, "JSON_OBJECTAGG": true,
	"SUBSTR": true, "SUBSTRING": true, "TRIM": true, "POSITION": true,
}

// spacedFuncWordIsName reports whether the function word just consumed is a plain
// identifier here, so the call is an ordinary generic one rather than the
// builtin's special form. p.prev is the word and p.cur its '('. Both maps gate
// their word the same way: only a word directly before '(' takes its folded
// token.
func (p *Parser) spacedFuncWordIsName(upperName string) bool {
	return (spacedTemporalWords[upperName] || spacedFuncWords[upperName]) && p.prev.End != p.cur.Loc
}

// atDatetimePrecision reports whether p.cur can be a function's fractional
// seconds precision, which is real_ulong_num. The lexer gives a digit run past 64
// bits as DECIMAL_NUM, which that refuses: mariadb:11.8.8 answers 1426 for
// now(18446744073709551615) and 1064 for now(18446744073709551616). HEX_NUM is the 0x spelling only, at any length: X'01' and
// 0X1 lex as something else and answer 1064. The lexer gives all three as
// tokXCONST, and a string's contents can begin 0x too, so this reads the source.
func (p *Parser) atDatetimePrecision() bool {
	switch p.cur.Type {
	case tokICONST:
		return true
	case tokXCONST:
		return strings.HasPrefix(p.lexer.input[p.cur.Loc-p.lexer.baseOffset:], "0x")
	}
	return false
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
