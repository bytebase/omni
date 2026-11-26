package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/bytebase/omni/dynamodb/ast"
)

// parseLiteral dispatches on the current token to produce one of the
// literal AST nodes. Handles the 6 scalar forms (NULL/MISSING/TRUE/
// FALSE/string/number/Ion) plus the DATE, TIME and TIMESTAMP literal
// forms.
//
// Grammar: literal (PartiQLParser.g4:775-787).
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

// parseDateTimeLiteral parses the three PartiQL date/time literal forms.
//
// Grammar (PartiQLParser.g4:784-786):
//
//	DATE LITERAL_STRING                                                  # LiteralDate
//	TIME ( PAREN_LEFT LITERAL_INTEGER PAREN_RIGHT )? (WITH TIME ZONE)? LITERAL_STRING   # LiteralTime
//	TIMESTAMP ( PAREN_LEFT LITERAL_INTEGER PAREN_RIGHT )? (WITH TIME ZONE)? LITERAL_STRING   # LiteralTimestamp
//
// The grammar takes any string after DATE, TIME or TIMESTAMP, and the
// reference parser's visitor refuses one that is no date, no time or no
// timestamp (validDateLiteral, checkTimeLiteral, checkTimestampLiteral), so
// omni does too.
func (p *Parser) parseDateTimeLiteral() (ast.ExprNode, error) {
	switch p.cur.Type {
	case tokDATE:
		return p.parseDateLiteral()
	case tokTIME, tokTIMESTAMP:
		return p.parseTimeLiteral()
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
//	TIMESTAMP ( PAREN_LEFT LITERAL_INTEGER PAREN_RIGHT )? (WITH TIME ZONE)? LITERAL_STRING
//
// cur is tokTIME or tokTIMESTAMP on entry.
func (p *Parser) parseTimeLiteral() (ast.ExprNode, error) {
	start := p.cur.Loc.Start
	kind := "TIME"
	check := checkTimeLiteral
	if p.cur.Type == tokTIMESTAMP {
		kind = "TIMESTAMP"
		check = checkTimestampLiteral
	}
	p.advance() // consume TIME or TIMESTAMP

	// Optional ( LITERAL_INTEGER ) precision.
	var precision *int
	if p.cur.Type == tokPAREN_LEFT {
		p.advance() // consume (
		intTok, err := p.expect(tokICONST)
		if err != nil {
			return nil, &ParseError{
				Message:  fmt.Sprintf("expected integer precision in %s(...), got %q", kind, p.cur.Str),
				Position: p.cur.Loc.Start, End: p.cur.Loc.End,
			}
		}
		n, perr := strconv.Atoi(intTok.Str)
		if perr != nil {
			return nil, &ParseError{
				Message:  fmt.Sprintf("invalid %s precision %q: %v", kind, intTok.Str, perr),
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

	body, err := p.expectStringLiteral(kind)
	if err != nil {
		return nil, err
	}
	if msg := check(body.Str, precision); msg != "" {
		return nil, &ParseError{
			Message:  fmt.Sprintf("%s string %q: %s", kind, body.Str, msg),
			Position: body.Loc.Start, End: body.Loc.End,
		}
	}
	if kind == "TIMESTAMP" {
		return &ast.TimestampLit{
			Val:          body.Str,
			Precision:    precision,
			WithTimeZone: withTZ,
			Loc:          ast.Loc{Start: start, End: body.Loc.End},
		}, nil
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
// PartiQLParserDefault.kt:483-485 and util/DateTimeUtils.kt:14-18. Java's
// \s also matches a vertical tab, which Go's does not, hence javaSpace.
const (
	javaSpace       = `[ \t\n\x0B\f\r]`
	timePatternBody = `(\d{2,}):(\d{2,}):(\d{2,})(?:\.(\d+))?` + javaSpace + `*(?:([+-])(\d\d):(\d\d)|[Zz])?`
)

var (
	datePatternRegex       = regexp.MustCompile(`^\d\d\d\d-\d\d-\d\d$`)
	genericTimeRegex       = regexp.MustCompile(`\d\d:\d\d:\d\d(\.\d*)?([+|-]\d\d:\d\d)?`)
	timeLiteralRegexp      = regexp.MustCompile(`^` + timePatternBody + `$`)
	timestampLiteralRegexp = regexp.MustCompile(`^(\d{4,})-(\d{2,})-(\d{2,})(?:` + javaSpace + `+|[Tt])` + timePatternBody + `$`)
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
	if msg := checkTimePrecision(s, precision); msg != "" {
		return msg
	}
	m := timeLiteralRegexp.FindStringSubmatch(s)
	if m == nil {
		return "expected the format HH:mm:ss[.ddd][+|-hh:mm|z]"
	}
	return checkTimeParts(m[1:])
}

// checkTimestampLiteral returns why the reference parser refuses s as a
// TIMESTAMP literal's body with the given precision, or "" when it takes
// it: the precision as for TIME, then a date, a space or T, and a time
// (PartiQLParserDefault.kt:2033-2042; DateTimeUtils.kt:21-31, 58-66).
func checkTimestampLiteral(s string, precision *int) string {
	if msg := checkTimePrecision(s, precision); msg != "" {
		return msg
	}
	m := timestampLiteralRegexp.FindStringSubmatch(s)
	if m == nil {
		return "expected the format YYYY-MM-DD[ |T]HH:mm:ss[.ddd][+|-hh:mm|z]"
	}
	year, yerr := strconv.Atoi(m[1])
	month, merr := strconv.Atoi(m[2])
	day, derr := strconv.Atoi(m[3])
	// java.time's years run to 999,999,999.
	if yerr != nil || merr != nil || derr != nil || year > 999999999 ||
		month < 1 || month > 12 || day < 1 || day > daysIn(year, month) {
		return "invalid date"
	}
	return checkTimeParts(m[4:])
}

// daysIn returns the number of days of the month of the year.
func daysIn(year, month int) int {
	return time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// checkTimePrecision returns why the reference parser refuses a TIME or
// TIMESTAMP literal's precision: an explicit one out of 0..9, or, with none,
// a string in which it finds no time to read one from
// (PartiQLParserDefault.kt:2177-2206).
func checkTimePrecision(s string, precision *int) string {
	if precision == nil {
		if !genericTimeRegex.MatchString(s) {
			return "unable to parse precision"
		}
	} else if *precision < 0 || *precision > 9 {
		return "precision out of bounds"
	}
	return ""
}

// checkTimeParts checks the groups of timePatternBody, hour, minute,
// second, fraction, zone sign, zone hour and zone minute, as the reference
// parser's time values do.
func checkTimeParts(m []string) string {
	hour, herr := strconv.Atoi(m[0])
	minute, merr := strconv.Atoi(m[1])
	second, serr := strconv.Atoi(m[2])
	if herr != nil || merr != nil || serr != nil || hour > 23 || minute > 59 || second > 59 {
		return "invalid time"
	}
	if m[4] != "" {
		tzHour, _ := strconv.Atoi(m[5])
		tzMinute, _ := strconv.Atoi(m[6])
		if tzHour > 23 || tzMinute > 59 {
			return "invalid time zone offset"
		}
	}
	return ""
}
