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
	// N'...' is a string literal too, which the expression reader has no arm for.
	if p.isBareWord("N") && p.peekNext().Type == tokSCONST && p.peekNext().Loc == p.cur.End {
		p.advance()
	}
	path, err := p.parsePrimaryExpr()
	if err != nil {
		return nil, err
	}
	if _, ok := path.(*nodes.StringLit); !ok {
		return nil, &ParseError{Message: "JSON_VALUE takes a string literal as its path", Position: pathStart.Loc}
	}
	fc.Args = []nodes.ExprNode{doc, path}

	opts := &nodes.JSONValueOptions{Loc: nodes.Loc{Start: p.pos()}}
	any := false
	if _, ok := p.match(kwRETURNING); ok {
		dt, err := p.parseCastDataType()
		if err != nil {
			return nil, err
		}
		opts.Returning, any = dt, true
	}
	if _, ok := p.match(kwTRUNCATE); ok {
		opts.Truncate, any = true, true
	}
	if _, ok := p.match(kwASCII); ok {
		opts.ASCII, any = true, true
	}
	if p.startsJSONValueResponse() {
		resp, err := p.parseJSONValueResponse()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(kwON); err != nil {
			return nil, err
		}
		any = true
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
	if any {
		opts.Loc.End = p.prev.End
		fc.JSONValue = opts
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	fc.Loc.End = p.prev.End
	return fc, nil
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
	signed := p.cur.Type == '-' || p.cur.Type == '+'
	value, err := p.parseExprPrec(precUnary)
	if err != nil {
		return nil, err
	}
	if !isJSONValueDefault(value, signed) {
		return nil, &ParseError{Message: "JSON_VALUE takes a literal after DEFAULT", Position: start}
	}
	return &nodes.JSONValueResponse{Loc: nodes.Loc{Start: begin, End: p.prev.End}, Default: value}, nil
}

// isJSONValueDefault reports whether value is what signed_literal reads. A sign takes a number,
// and a plus sign leaves no node behind.
func isJSONValueDefault(value nodes.ExprNode, signed bool) bool {
	if unary, ok := value.(*nodes.UnaryExpr); ok {
		if !signed {
			return false
		}
		value = unary.Operand
	}
	switch value.(type) {
	case *nodes.IntLit, *nodes.FloatLit:
		return true
	case *nodes.StringLit, *nodes.BoolLit, *nodes.HexLit, *nodes.BitLit, *nodes.TemporalLit, *nodes.NullLit:
		return !signed
	}
	return false
}
