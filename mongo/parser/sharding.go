package parser

import "github.com/bytebase/omni/mongo/ast"

// parseShStatement parses sh.method() calls.
// Grammar: sh DOT identifier LPAREN arguments? RPAREN
func (p *Parser) parseShStatement() (ast.Node, error) {
	stmtStart := p.cur.Loc
	p.advance() // consume 'sh'

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

	return &ast.ShStatement{
		MethodName: methodName,
		Args:       args,
		Loc:        ast.Loc{Start: stmtStart, End: p.prev.End},
	}, nil
}
