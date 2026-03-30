package parser

import (
	"fmt"

	nodes "github.com/bytebase/omni/cosmosdb/ast"
)

// parseFuncCall parses a function call where the name has already been consumed.
// Expects to see ( [*|args] ).
func (p *Parser) parseFuncCall(name string, startLoc int) (nodes.ExprNode, error) {
	if _, err := p.expect(tokLPAREN); err != nil {
		return nil, err
	}

	fc := &nodes.FuncCall{Name: name}

	// Check for COUNT(*).
	if p.cur.Type == tokSTAR {
		fc.Star = true
		p.advance()
	} else if p.cur.Type != tokRPAREN {
		// Parse argument list.
		for {
			arg, err := p.parseExpr(0)
			if err != nil {
				return nil, err
			}
			fc.Args = append(fc.Args, arg)
			if !p.match(tokCOMMA) {
				break
			}
		}
	}

	rparen, err := p.expect(tokRPAREN)
	if err != nil {
		return nil, err
	}

	fc.Loc = nodes.Loc{Start: startLoc, End: rparen.Loc + 1}
	return fc, nil
}

// parseUDFCall parses UDF.name(args).
func (p *Parser) parseUDFCall() (nodes.ExprNode, error) {
	startLoc := p.cur.Loc
	p.advance() // consume UDF

	if _, err := p.expect(tokDOT); err != nil {
		return nil, err
	}

	// The function name after UDF.
	if p.cur.Type != tokIDENT && !isIdentLike(p.cur.Type) {
		return nil, &ParseError{
			Message: fmt.Sprintf("expected function name after UDF., got %q", p.cur.Str),
			Pos:     p.cur.Loc,
		}
	}
	name := p.cur.Str
	p.advance()

	if _, err := p.expect(tokLPAREN); err != nil {
		return nil, err
	}

	var args []nodes.ExprNode
	if p.cur.Type != tokRPAREN {
		for {
			arg, err := p.parseExpr(0)
			if err != nil {
				return nil, err
			}
			args = append(args, arg)
			if !p.match(tokCOMMA) {
				break
			}
		}
	}

	rparen, err := p.expect(tokRPAREN)
	if err != nil {
		return nil, err
	}

	return &nodes.UDFCall{
		Name: name,
		Args: args,
		Loc:  nodes.Loc{Start: startLoc, End: rparen.Loc + 1},
	}, nil
}
