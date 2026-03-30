package parser

import "github.com/bytebase/omni/mongo/ast"

// parseBulkStatement parses db.collection.initializeOrderedBulkOp()/initializeUnorderedBulkOp() chains.
func (p *Parser) parseBulkStatement(database, collection string, ordered bool, startLoc int) (ast.Node, error) {
	return nil, p.syntaxErrorAtCur()
}
