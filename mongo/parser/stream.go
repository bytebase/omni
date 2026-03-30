package parser

import "github.com/bytebase/omni/mongo/ast"

// parseSpStatement parses sp.method() and sp.x.method() calls.
func (p *Parser) parseSpStatement() (ast.Node, error) {
	return nil, p.syntaxErrorAtCur()
}
