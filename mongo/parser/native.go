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

// parseQualifiedNativeCall parses top-level calls on built-in shell namespace
// objects such as config.get("...") and EJSON.parse("...").
func (p *Parser) parseQualifiedNativeCall() (ast.Node, error) {
	stmtStart := p.cur.Loc
	owner := p.cur.Str
	p.advance() // consume owner (config, EJSON)

	if _, err := p.expect('.'); err != nil {
		return nil, err
	}
	mc, err := p.parseMethodCall()
	if err != nil {
		return nil, err
	}

	return &ast.NativeFunctionCall{
		Name: owner + "." + mc.Name,
		Args: mc.Args,
		Loc:  ast.Loc{Start: stmtStart, End: mc.Loc.End},
	}, nil
}
