package parser

import "github.com/bytebase/omni/mongo/ast"

// parseRsStatement parses rs.method() calls.
// Grammar: rs DOT identifier LPAREN arguments? RPAREN
func (p *Parser) parseRsStatement() (ast.Node, error) {
	stmtStart := p.cur.Loc
	p.advance() // consume 'rs'

	if _, err := p.expect('.'); err != nil {
		return nil, err
	}

	methodName, _, err := p.expectIdent()
	if err != nil {
		return nil, err
	}

	args, err := p.parseArguments()
	if err != nil {
		return nil, err
	}

	return &ast.RsStatement{
		MethodName: methodName,
		Args:       args,
		Loc:        ast.Loc{Start: stmtStart, End: p.prev.End},
	}, nil
}
