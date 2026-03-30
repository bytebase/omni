package parser

import "github.com/bytebase/omni/mongo/ast"

// parseConnectionStatement parses Mongo() and connect() chains.
// Grammar: MONGO LPAREN arguments? RPAREN connectionMethodChain?
//        | CONNECT LPAREN arguments? RPAREN connectionMethodChain?
func (p *Parser) parseConnectionStatement() (ast.Node, error) {
	stmtStart := p.cur.Loc
	constructor := p.cur.Str
	p.advance() // consume 'Mongo' or 'connect'

	args, err := p.parseArguments()
	if err != nil {
		return nil, err
	}

	// Parse optional method chain: .method().method()...
	var chain []ast.MethodCall
	for p.cur.Type == '.' {
		p.advance() // consume '.'
		mc, err := p.parseMethodCall()
		if err != nil {
			return nil, err
		}
		chain = append(chain, *mc)
	}

	return &ast.ConnectionStatement{
		Constructor:    constructor,
		Args:           args,
		ChainedMethods: chain,
		Loc:            ast.Loc{Start: stmtStart, End: p.prev.End},
	}, nil
}
