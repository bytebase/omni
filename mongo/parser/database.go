package parser

import "github.com/bytebase/omni/mongo/ast"

// parseDbStatement parses db.method() and db.collection.method() statements.
// Dispatches to collection, bulk, encryption, or plan cache parsers as needed.
func (p *Parser) parseDbStatement() (ast.Node, error) {
	return nil, p.syntaxErrorAtCur()
}
