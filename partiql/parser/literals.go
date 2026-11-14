package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/bytebase/omni/partiql/ast"
)

// parseLiteral dispatches on the current token to produce one of the
// literal AST nodes. Handles the 6 scalar forms (NULL/MISSING/TRUE/
// FALSE/string/number/Ion) plus the DATE and TIME literal forms
// directly. TIMESTAMP is routed here by the dispatcher but is NOT a
// PartiQL literal (see parseDateTimeLiteral) and is rejected.
//
// Grammar: literal (PartiQLParser.g4 lines 661-672).
func (p *Parser) parseLiteral() (ast.ExprNode, error) {
	switch p.cur.Type {
	case tokNULL:
		loc := p.cur.Loc
		p.advance()
		return &ast.NullLit{Loc: loc}, nil

	case tokMISSING:
		loc := p.cur.Loc
		p.advance()
		return &ast.MissingLit{Loc: loc}, nil

	case tokTRUE:
		loc := p.cur.Loc
		p.advance()
		return &ast.BoolLit{Val: true, Loc: loc}, nil

	case tokFALSE:
		loc := p.cur.Loc
		p.advance()
		return &ast.BoolLit{Val: false, Loc: loc}, nil

	case tokSCONST:
		val := p.cur.Str
		loc := p.cur.Loc
		p.advance()
		return &ast.StringLit{Val: val, Loc: loc}, nil

	case tokICONST, tokFCONST:
		// NumberLit.Val preserves the raw source text so callers can
		// distinguish integer/decimal/scientific forms. Token.Str is
		// already the raw slice (scanNumber in the lexer preserves it).
		val := p.cur.Str
		loc := p.cur.Loc
		p.advance()
		return &ast.NumberLit{Val: val, Loc: loc}, nil

	case tokION_LITERAL:
		// Lexer's scanIonLiteral delivers the verbatim inner content
		// (between backticks) in Token.Str. No further decoding at
		// this layer — that is DAG node 17's job.
		text := p.cur.Str
		loc := p.cur.Loc
		p.advance()
		return &ast.IonLit{Text: text, Loc: loc}, nil

	case tokDATE, tokTIME, tokTIMESTAMP:
		return p.parseDateTimeLiteral()
	}

	return nil, &ParseError{
		Message:  "expected literal",
		Position: p.cur.Loc.Start, End: p.cur.Loc.End,
	}
}

// parseDateTimeLiteral parses the two PartiQL date/time literal forms
// and rejects the (non-existent) TIMESTAMP literal form.
//
// Grammar (PartiQLParser.g4 lines 670-671):
//
//	DATE LITERAL_STRING                                                  # LiteralDate
//	TIME ( PAREN_LEFT LITERAL_INTEGER PAREN_RIGHT )? (WITH TIME ZONE)? LITERAL_STRING   # LiteralTime
//
// TIMESTAMP is intentionally NOT a literal: the canonical PartiQL
// grammar (partiql.org) has no TIMESTAMP alternative in the `literal`
// rule, and AWS DynamoDB PartiQL has no temporal types at all — the
// only way to write a timestamp value is an Ion literal or
// CAST(... AS TIMESTAMP). TIMESTAMP appears solely as a type-name
// keyword in the `type` rule (line 677). The dispatcher routes
// TIMESTAMP here only so a future grammar revision has a single seam to
// extend; today it is a precise syntax error.
//
// The grammar takes any string after DATE or TIME, and the reference
// parser's visitor refuses one that is no date or no time
// (validDateLiteral, validTimeLiteral), so omni does too.
func (p *Parser) parseDateTimeLiteral() (ast.ExprNode, error) {
	switch p.cur.Type {
	case tokDATE:
		return p.parseDateLiteral()
	case tokTIME:
		return p.parseTimeLiteral()
	case tokTIMESTAMP:
		// Fail-fast: the error Loc points at the TIMESTAMP keyword
		// (cur). No advance needed — parsing aborts on this error.
		return nil, &ParseError{
			Message: "TIMESTAMP literal is not supported in PartiQL; " +
				"use an Ion literal (`...`) or CAST(... AS TIMESTAMP)",
			Position: p.cur.Loc.Start, End: p.cur.Loc.End,
		}
	}
	// Unreachable: callers gate on the three tokens above.
	return nil, &ParseError{
		Message:  "expected DATE, TIME, or TIMESTAMP literal",
		Position: p.cur.Loc.Start, End: p.cur.Loc.End,
	}
}

// parseDateLiteral parses `DATE LITERAL_STRING`. cur is tokDATE on entry.
func (p *Parser) parseDateLiteral() (ast.ExprNode, error) {
	start := p.cur.Loc.Start
	p.advance() // consume DATE
	body, err := p.expectStringLiteral("DATE")
	if err != nil {
		return nil, err
	}
	if !validDateLiteral(body.Str) {
		return nil, &ParseError{
			Message:  fmt.Sprintf("DATE string %q is not a date of the format yyyy-MM-dd", body.Str),
			Position: body.Loc.Start, End: body.Loc.End,
		}
	}
	return &ast.DateLit{
		Val: body.Str,
		Loc: ast.Loc{Start: start, End: body.Loc.End},
	}, nil
}

// parseTimeLiteral parses
//
//	TIME ( PAREN_LEFT LITERAL_INTEGER PAREN_RIGHT )? (WITH TIME ZONE)? LITERAL_STRING
//
// cur is tokTIME on entry.
func (p *Parser) parseTimeLiteral() (ast.ExprNode, error) {
	start := p.cur.Loc.Start
	p.advance() // consume TIME

	// Optional ( LITERAL_INTEGER ) precision.
	var precision *int
	if p.cur.Type == tokPAREN_LEFT {
		p.advance() // consume (
		intTok, err := p.expect(tokICONST)
		if err != nil {
			return nil, &ParseError{
				Message:  fmt.Sprintf("expected integer precision in TIME(...), got %q", p.cur.Str),
				Position: p.cur.Loc.Start, End: p.cur.Loc.End,
			}
		}
		n, perr := strconv.Atoi(intTok.Str)
		if perr != nil {
			return nil, &ParseError{
				Message:  fmt.Sprintf("invalid TIME precision %q: %v", intTok.Str, perr),
				Position: intTok.Loc.Start, End: intTok.Loc.End,
			}
		}
		if _, err := p.expect(tokPAREN_RIGHT); err != nil {
			return nil, err
		}
		precision = &n
	}

	// Optional WITH TIME ZONE.
	withTZ := false
	if p.cur.Type == tokWITH {
		p.advance() // consume WITH
		if _, err := p.expect(tokTIME); err != nil {
			return nil, &ParseError{
				Message:  fmt.Sprintf("expected TIME after WITH, got %q", p.cur.Str),
				Position: p.cur.Loc.Start, End: p.cur.Loc.End,
			}
		}
		if _, err := p.expect(tokZONE); err != nil {
			return nil, &ParseError{
				Message:  fmt.Sprintf("expected ZONE after WITH TIME, got %q", p.cur.Str),
				Position: p.cur.Loc.Start, End: p.cur.Loc.End,
			}
		}
		withTZ = true
	}

	body, err := p.expectStringLiteral("TIME")
	if err != nil {
		return nil, err
	}
	if msg := checkTimeLiteral(body.Str, precision); msg != "" {
		return nil, &ParseError{
			Message:  fmt.Sprintf("TIME string %q: %s", body.Str, msg),
			Position: body.Loc.Start, End: body.Loc.End,
		}
	}
	return &ast.TimeLit{
		Val:          body.Str,
		Precision:    precision,
		WithTimeZone: withTZ,
		Loc:          ast.Loc{Start: start, End: body.Loc.End},
	}, nil
}

// expectStringLiteral consumes a required LITERAL_STRING token, returning
// a descriptive error keyed on the literal kind when the next token is
// not a string. The returned Token's Str is the already-decoded body.
//
// The string body is the FINAL token of every DATE/TIME literal form, so
// the advance that consumes it is where the lexer first scans whatever
// follows the literal. If that trailing content is itself un-lexable
// (an unrecognized character like `#`, an unterminated block comment
// `/*`, ...), the lexer's first-error-and-stop contract sets Lexer.Err
// and makes advance() yield tokEOF rather than the offending token.
// A bare EOF check by the caller would then mistake the malformed input
// for a clean end and accept the literal, silently dropping the lexer
// error. We surface it here so e.g. `DATE '2026-01-01' #` and
// `TIME '12:00:00' /*` become parse errors (matching the antlr oracle,
// which rejects the same inputs as unexpected trailing input) instead of
// accepted literals. Clean trailing content — EOF, or a real token the
// EOF check will catch — leaves Lexer.Err nil and is unaffected.
func (p *Parser) expectStringLiteral(kind string) (Token, error) {
	if p.cur.Type != tokSCONST {
		return Token{}, &ParseError{
			Message:  fmt.Sprintf("expected string literal after %s, got %q", kind, p.cur.Str),
			Position: p.cur.Loc.Start, End: p.cur.Loc.End,
		}
	}
	tok := p.cur
	p.advance()
	if err := p.checkLexerErr(); err != nil {
		return Token{}, err
	}
	return tok, nil
}

// The patterns of the reference parser at v0.14.9:
// PartiQLParserDefault.kt:483-485 and util/DateTimeUtils.kt:14-15.
var (
	datePatternRegex  = regexp.MustCompile(`^\d\d\d\d-\d\d-\d\d$`)
	genericTimeRegex  = regexp.MustCompile(`\d\d:\d\d:\d\d(\.\d*)?([+|-]\d\d:\d\d)?`)
	timeLiteralRegexp = regexp.MustCompile(`^(\d{2,}):(\d{2,}):(\d{2,})(?:\.(\d+))?\s*(?:([+-])(\d\d):(\d\d)|[Zz])?$`)
)

// validDateLiteral reports whether s is a DATE literal's body the reference
// parser takes: yyyy-MM-dd and a day of the calendar
// (PartiQLParserDefault.kt:2005-2018, ISO_LOCAL_DATE).
func validDateLiteral(s string) bool {
	if !datePatternRegex.MatchString(s) {
		return false
	}
	_, err := time.Parse("2006-01-02", s)
	return err == nil
}

// checkTimeLiteral returns why the reference parser refuses s as a TIME
// literal's body with the given precision, or "" when it takes it
// (PartiQLParserDefault.kt:2020-2028, 2177-2206; DateTimeUtils.kt:34-56).
func checkTimeLiteral(s string, precision *int) string {
	if precision == nil {
		if !genericTimeRegex.MatchString(s) {
			return "unable to parse precision"
		}
	} else if *precision < 0 || *precision > 9 {
		return "precision out of bounds"
	}
	m := timeLiteralRegexp.FindStringSubmatch(s)
	if m == nil {
		return "expected the format HH:mm:ss[.ddd][+|-hh:mm|z]"
	}
	hour, herr := strconv.Atoi(m[1])
	minute, merr := strconv.Atoi(m[2])
	second, serr := strconv.Atoi(m[3])
	if herr != nil || merr != nil || serr != nil || hour > 23 || minute > 59 || second > 59 {
		return "invalid time"
	}
	if m[5] != "" {
		tzHour, _ := strconv.Atoi(m[6])
		tzMinute, _ := strconv.Atoi(m[7])
		if tzHour > 23 || tzMinute > 59 {
			return "invalid time zone offset"
		}
	}
	return ""
}
