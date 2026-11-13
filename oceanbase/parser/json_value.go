package parser

import (
	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// parseJSONValueFunc reads what follows `JSON_VALUE(`, which has its own arm in the grammar and
// not the generic call's expression list:
//
//	JSON_VALUE ( simple_expr , string_literal [RETURNING cast_type]
//	    [response ON EMPTY] [response ON ERROR] )
//	response: NULL | ERROR | DEFAULT signed_literal
//
// OceanBase also takes [TRUNCATE] [ASCII] between the type and the responses, in that order.
// So OceanBase 4.3.5.6 answers 1064 for `JSON_VALUE(a)`, `JSON_VALUE(1,2)`, `JSON_VALUE(a+1,'$')`,
// `JSON_VALUE(a,@p)` and `JSON_VALUE(a,'$' NULL ON ERROR NULL ON EMPTY)`, and reads
// `JSON_VALUE(-a,'$')`, `JSON_VALUE((1+1),'$')` and `JSON_VALUE(a,N'$')`. Quoted or qualified,
// the name is an ordinary call.
//
// Ref: json_value_expr in sql_parser_mysql_mode.y (v4.3.5_CE_BP6).
func (p *Parser) parseJSONValueFunc(fc *nodes.FuncCallExpr) (nodes.ExprNode, error) {
	if err := p.parseJSONDocAndPath(fc); err != nil {
		return nil, err
	}

	opts := &nodes.JSONValueOptions{Loc: nodes.Loc{Start: p.pos()}}
	if _, ok := p.match(kwRETURNING); ok {
		dt, err := p.parseCastDataType()
		if err != nil {
			return nil, err
		}
		opts.Returning = dt
		// The type reader takes the ASCII of CHAR(3) ASCII, and here it is this clause's. In
		// CHAR(3) CHARSET ascii the word is a charset's name.
		opts.ASCII = dt.Charset == "" && p.prev.Type == kwASCII
	}
	if _, ok := p.match(kwTRUNCATE); ok {
		opts.Truncate = true
	}
	if _, ok := p.match(kwASCII); ok {
		opts.ASCII = true
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

// parseJSONDocAndPath reads what JSON_VALUE and JSON_QUERY both open with, into fc.Args: a
// document that is a simple_expr, a comma, and a path that is a string literal.
//
// Ref: `simple_expr ',' complex_string_literal` in json_value_expr and json_query_expr,
// sql_parser_mysql_mode.y (v4.3.5_CE_BP6).
func (p *Parser) parseJSONDocAndPath(fc *nodes.FuncCallExpr) error {
	if p.cur.Type == '*' {
		return p.syntaxErrorAtCur()
	}
	doc, err := p.parseExprPrec(precCollate)
	if err != nil {
		return err
	}
	if doc == nil {
		return p.syntaxErrorAtCur()
	}
	if _, err := p.expect(','); err != nil {
		return err
	}
	pathStart := p.cur
	// Unlike MySQL, OceanBase 4.3.5.6 takes an introduced hex or bit literal as the path:
	// JSON_VALUE('{}', _latin1 X'24') reads.
	national := p.skipNationalPrefix()
	path, err := p.parsePrimaryExpr()
	if err != nil {
		return err
	}
	pathLit, ok := path.(*nodes.StringLit)
	if !ok {
		return &ParseError{Message: fc.Name + " takes a string literal as its path", Position: pathStart.Loc}
	}
	if national {
		// The literal starts at its N, which the reader stepped over.
		pathLit.Loc.Start = pathStart.Loc
	}
	fc.Args = []nodes.ExprNode{doc, path}
	return nil
}

// skipNationalPrefix steps over the N of N'...', which is a string literal too and which the
// expression reader has no arm for: it would read the N as a column. The quote has to follow at
// once, since `N 'x'` is a name and a string.
func (p *Parser) skipNationalPrefix() bool {
	if p.isBareWord("N") && p.peekNext().Type == tokSCONST && p.peekNext().Loc == p.cur.End {
		p.advance()
		return true
	}
	return false
}

func (p *Parser) startsJSONValueResponse() bool {
	return p.cur.Type == kwNULL || p.cur.Type == kwERROR_KW || p.cur.Type == kwDEFAULT
}

// parseJSONValueResponse reads NULL, ERROR, or DEFAULT and a literal, which may carry a sign and
// is no parenthesized value, no variable and no expression. A DEFAULT NULL reads and is refused
// later, with 1149.
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
	if lit, ok := value.(*nodes.StringLit); ok && national {
		// The literal starts at its N, which the reader stepped over.
		lit.Loc.Start = start
	}
	return &nodes.JSONValueResponse{Loc: nodes.Loc{Start: begin, End: p.prev.End}, Default: value}, nil
}

// isJSONValueDefault reports whether value is a literal that takes no sign.
func isJSONValueDefault(value nodes.ExprNode) bool {
	switch value.(type) {
	case *nodes.IntLit, *nodes.FloatLit, *nodes.StringLit, *nodes.BoolLit, *nodes.HexLit, *nodes.BitLit, *nodes.TemporalLit,
		*nodes.NullLit:
		return true
	}
	return false
}
