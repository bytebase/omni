package parser

import "github.com/bytebase/omni/mongo/ast"

// parseShStatement parses sh.method() calls.
func (p *Parser) parseShStatement() (ast.Node, error) {
	return nil, p.syntaxErrorAtCur()
}
