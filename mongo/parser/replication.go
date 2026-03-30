package parser

import "github.com/bytebase/omni/mongo/ast"

// parseRsStatement parses rs.method() calls.
func (p *Parser) parseRsStatement() (ast.Node, error) {
	return nil, p.syntaxErrorAtCur()
}
