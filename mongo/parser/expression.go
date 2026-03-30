package parser

import "github.com/bytebase/omni/mongo/ast"

// parseExpression parses a value expression (document, array, literal, helper, identifier).
func (p *Parser) parseExpression() (ast.Node, error) {
	return nil, p.syntaxErrorAtCur()
}
