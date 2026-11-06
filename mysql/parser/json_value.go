package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// parseJSONValueFunc reads what follows `JSON_VALUE(`, which has its own arm in the grammar and
// not the generic call's expression list:
//
//	JSON_VALUE ( simple_expr , string_literal [RETURNING cast_type]
//	    [response ON EMPTY] [response ON ERROR] )
//	response: NULL | ERROR | DEFAULT signed_literal
//
// So MySQL 8.0 answers 1064 for `JSON_VALUE(a)`, `JSON_VALUE(1,2)`, `JSON_VALUE(a+1,'$')`,
// `JSON_VALUE(a,@p)` and `JSON_VALUE(a,'$' NULL ON ERROR NULL ON EMPTY)`, and reads
// `JSON_VALUE(-a,'$')`, `JSON_VALUE((1+1),'$')` and `JSON_VALUE(a,N'$')`. Quoted or qualified,
// the name is an ordinary call.
//
// Ref: JSON_VALUE_SYM in MySQL 8.0 sql/sql_yacc.yy.
func (p *Parser) parseJSONValueFunc(fc *nodes.FuncCallExpr) (nodes.ExprNode, error) {
	if p.cur.Type == '*' {
		return nil, p.syntaxErrorAtCur()
	}
	doc, err := p.parseExprPrec(precCollate)
	if err != nil {
		return nil, err
	}
	if doc == nil {
		return nil, p.syntaxErrorAtCur()
	}
	if _, err := p.expect(','); err != nil {
		return nil, err
	}
	pathStart := p.cur
	// The path is a text_literal: an introducer before a hex or bit literal makes a string too,
	// but not one the path takes (1064 on mysql:8.0.46).
	if p.cur.Type == tokIDENT && isCharsetIntroducer(p.cur.Str) && !p.backtickQuotedAt(p.cur.Loc) {
		if next := p.peekNext().Type; next == tokXCONST || next == tokBCONST {
			return nil, &ParseError{Message: "JSON_VALUE takes a string literal as its path", Position: pathStart.Loc}
		}
	}
	national := p.skipNationalPrefix()
	path, err := p.parsePrimaryExpr()
	if err != nil {
		return nil, err
	}
	pathLit, ok := path.(*nodes.StringLit)
	if !ok {
		return nil, &ParseError{Message: "JSON_VALUE takes a string literal as its path", Position: pathStart.Loc}
	}
	pathLit.National = national
	if national {
		// The literal starts at its N, which the reader stepped over.
		pathLit.Loc.Start = pathStart.Loc
	}
	fc.Args = []nodes.ExprNode{doc, path}

	opts := &nodes.JSONValueOptions{Loc: nodes.Loc{Start: p.pos()}}
	if _, ok := p.match(kwRETURNING); ok {
		dt, err := p.parseCastDataType()
		if err != nil {
			return nil, err
		}
		// A geometry type reads as a cast type here and answers 4032 while parsing.
		switch dt.Name {
		case "GEOMETRY", "POINT", "LINESTRING", "POLYGON", "MULTIPOINT", "MULTILINESTRING",
			"MULTIPOLYGON", "GEOMETRYCOLLECTION", "GEOMCOLLECTION":
			return nil, &ParseError{Message: "Invalid JSON type in argument 1 to function json_value", Position: dt.Loc.Start}
		}
		opts.Returning = dt
	}
	if p.startsJSONValueResponse() {
		resp, err := p.parseJSONValueResponse()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(kwON); err != nil {
			return nil, err
		}
		if _, ok := p.match(kwEMPTY); ok {
			opts.OnEmpty = resp
			if p.startsJSONValueResponse() {
				if opts.OnError, err = p.parseJSONValueResponse(); err != nil {
					return nil, err
				}
				if _, err := p.expect(kwON); err != nil {
					return nil, err
				}
				if _, err := p.expect(kwERROR_KW); err != nil {
					return nil, err
				}
			}
		} else {
			if _, err := p.expect(kwERROR_KW); err != nil {
				return nil, err
			}
			opts.OnError = resp
		}
	}
	fc.JSONValueArm = true
	if p.prev.End > opts.Loc.Start {
		opts.Loc.End = p.prev.End
		fc.JSONValue = opts
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	fc.Loc.End = p.prev.End
	return fc, nil
}

// skipNationalPrefix steps over the N of N'...', which is a string literal too and which the
// expression reader has no arm for: it would read the N as a column. The quote has to follow at
// once, since `N 'x'` is a name and a string. It reports whether there was one.
func (p *Parser) skipNationalPrefix() bool {
	if strings.EqualFold(p.cur.Str, "N") && p.cur.Type == tokIDENT && !p.backtickQuotedAt(p.cur.Loc) {
		if next := p.peekNext(); next.Type == tokSCONST && next.Loc == p.cur.End {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) startsJSONValueResponse() bool {
	return p.cur.Type == kwNULL || p.cur.Type == kwERROR_KW || p.cur.Type == kwDEFAULT
}

// parseJSONValueResponse reads NULL, ERROR, or DEFAULT and a literal, which may carry a sign and
// is no NULL, no parenthesized value, no variable and no expression.
func (p *Parser) parseJSONValueResponse() (*nodes.JSONValueResponse, error) {
	begin := p.pos()
	if _, ok := p.match(kwNULL); ok {
		return &nodes.JSONValueResponse{Loc: nodes.Loc{Start: begin, End: p.prev.End}}, nil
	}
	if _, ok := p.match(kwERROR_KW); ok {
		return &nodes.JSONValueResponse{Loc: nodes.Loc{Start: begin, End: p.prev.End}, Error: true}, nil
	}
	p.advance() // consume DEFAULT
	start := p.cur.Loc
	// signed_literal is one sign and a number straight after it, token by token: - 1, +1.5 and
	// -1e2 read, and +~1, --1, -(1) and -0x61 are 1064 (mysql:8.0.46). The expression reader would
	// fold the first sign into what follows, so the sign is read here.
	if p.cur.Type == '-' || p.cur.Type == '+' {
		sign := p.advance()
		if p.cur.Type != tokICONST && p.cur.Type != tokFCONST {
			return nil, &ParseError{Message: "JSON_VALUE takes a literal after DEFAULT", Position: start}
		}
		number, err := p.parsePrimaryExpr()
		if err != nil {
			return nil, err
		}
		value := number
		if sign.Type == '-' {
			value = &nodes.UnaryExpr{Loc: nodes.Loc{Start: sign.Loc, End: p.prev.End}, Op: nodes.UnaryMinus, Operand: number}
		}
		return &nodes.JSONValueResponse{Loc: nodes.Loc{Start: begin, End: p.prev.End}, Default: value}, nil
	}
	national := p.skipNationalPrefix()
	// OFF reads as a boolean in an expression, and signed_literal has no arm for it: DEFAULT OFF
	// is a 1064 where DEFAULT FALSE reads.
	if p.cur.Type == kwOFF {
		return nil, &ParseError{Message: "JSON_VALUE takes a literal after DEFAULT", Position: start}
	}
	// A hex or bit literal is held to its own digits, which the expression reader takes on
	// trust: b'2', X'6' and X'GG' are 1064.
	if p.cur.Type == tokXCONST || p.cur.Type == tokBCONST {
		if _, err := p.decodeByteLiteral(p.cur); err != nil {
			return nil, err
		}
	}
	value, err := p.parseExprPrec(precUnary)
	if err != nil {
		return nil, err
	}
	if !isJSONValueDefault(value) {
		return nil, &ParseError{Message: "JSON_VALUE takes a literal after DEFAULT", Position: start}
	}
	if lit, ok := value.(*nodes.StringLit); ok {
		lit.National = national
		if national {
			lit.Loc.Start = start
		}
	}
	return &nodes.JSONValueResponse{Loc: nodes.Loc{Start: begin, End: p.prev.End}, Default: value}, nil
}

// isJSONValueDefault reports whether value is a literal that takes no sign.
func isJSONValueDefault(value nodes.ExprNode) bool {
	switch value.(type) {
	case *nodes.IntLit, *nodes.FloatLit, *nodes.StringLit, *nodes.BoolLit, *nodes.HexLit, *nodes.BitLit, *nodes.TemporalLit:
		return true
	}
	return false
}
