package parser

import "github.com/bytebase/omni/mongo/ast"

// parseDocument parses a { key: value, ... } document literal.
func (p *Parser) parseDocument() (ast.Node, error) {
	return nil, p.syntaxErrorAtCur()
}

// parseArray parses a [ elem, ... ] array literal.
func (p *Parser) parseArray() (ast.Node, error) {
	return nil, p.syntaxErrorAtCur()
}
