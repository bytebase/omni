package parser

import "github.com/bytebase/omni/mongo/ast"

// parseNativeFunctionCall parses a top-level native function call like sleep(1000).
func (p *Parser) parseNativeFunctionCall() (ast.Node, error) {
	return nil, p.syntaxErrorAtCur()
}
