package parser

import "github.com/bytebase/omni/mongo/ast"

// parseSpStatement parses sp.method() and sp.x.method() calls.
// Grammar: SP DOT identifier LPAREN arguments? RPAREN
//        | SP DOT identifier DOT identifier LPAREN arguments? RPAREN
func (p *Parser) parseSpStatement() (ast.Node, error) {
	stmtStart := p.cur.Loc
	p.advance() // consume 'sp'

	if _, err := p.expect('.'); err != nil {
		return nil, err
	}

	firstName, _, err := p.expectIdent()
	if err != nil {
		return nil, err
	}

	// Check if there's a sub-method: sp.processorName.method()
	if p.cur.Type == '.' {
		p.advance() // consume '.'
		subMethod, _, err := p.expectIdent()
		if err != nil {
			return nil, err
		}

		args, err := p.parseArguments()
		if err != nil {
			return nil, err
		}

		return &ast.SpStatement{
			MethodName: firstName,
			SubMethod:  subMethod,
			Args:       args,
			Loc:        ast.Loc{Start: stmtStart, End: p.prev.End},
		}, nil
	}

	// Simple form: sp.method()
	args, err := p.parseArguments()
	if err != nil {
		return nil, err
	}

	return &ast.SpStatement{
		MethodName: firstName,
		Args:       args,
		Loc:        ast.Loc{Start: stmtStart, End: p.prev.End},
	}, nil
}
