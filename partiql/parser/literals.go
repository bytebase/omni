package parser

import (
	"fmt"
	"strconv"

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
		Message: "expected literal",
		Loc:     p.cur.Loc,
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
// The string body is NOT validated at the syntax layer — any string is
// grammatically accepted after DATE/TIME, matching the legacy ANTLR
// parser. Calendar/clock validity is a downstream semantic concern.
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
			Loc: p.cur.Loc,
		}
	}
	// Unreachable: callers gate on the three tokens above.
	return nil, &ParseError{
		Message: "expected DATE, TIME, or TIMESTAMP literal",
		Loc:     p.cur.Loc,
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
				Message: fmt.Sprintf("expected integer precision in TIME(...), got %q", p.cur.Str),
				Loc:     p.cur.Loc,
			}
		}
		n, perr := strconv.Atoi(intTok.Str)
		if perr != nil {
			return nil, &ParseError{
				Message: fmt.Sprintf("invalid TIME precision %q: %v", intTok.Str, perr),
				Loc:     intTok.Loc,
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
				Message: fmt.Sprintf("expected TIME after WITH, got %q", p.cur.Str),
				Loc:     p.cur.Loc,
			}
		}
		if _, err := p.expect(tokZONE); err != nil {
			return nil, &ParseError{
				Message: fmt.Sprintf("expected ZONE after WITH TIME, got %q", p.cur.Str),
				Loc:     p.cur.Loc,
			}
		}
		withTZ = true
	}

	body, err := p.expectStringLiteral("TIME")
	if err != nil {
		return nil, err
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
func (p *Parser) expectStringLiteral(kind string) (Token, error) {
	if p.cur.Type != tokSCONST {
		return Token{}, &ParseError{
			Message: fmt.Sprintf("expected string literal after %s, got %q", kind, p.cur.Str),
			Loc:     p.cur.Loc,
		}
	}
	tok := p.cur
	p.advance()
	return tok, nil
}
