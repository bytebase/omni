package parser

import "github.com/bytebase/omni/mongo/ast"

// parseNativeFunctionCall parses a top-level native function call like sleep(1000).
// Grammar: identifier LPAREN arguments? RPAREN
func (p *Parser) parseNativeFunctionCall() (ast.Node, error) {
	stmtStart := p.cur.Loc
	name := p.cur.Str
	p.advance() // consume function name

	args, err := p.parseArguments()
	if err != nil {
		return nil, err
	}

	return &ast.NativeFunctionCall{
		Name: name,
		Args: args,
		Loc:  ast.Loc{Start: stmtStart, End: p.prev.End},
	}, nil
}
