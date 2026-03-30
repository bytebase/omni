package parser

import "github.com/bytebase/omni/mongo/ast"

// parseCollectionMethod parses db.collection.method().cursor()... chains.
func (p *Parser) parseCollectionMethod(database, collection string, startLoc int) (ast.Node, error) {
	return nil, p.syntaxErrorAtCur()
}
