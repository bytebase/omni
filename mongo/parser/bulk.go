package parser

import "github.com/bytebase/omni/mongo/ast"

// parseBulkStatement parses db.collection.initializeOrderedBulkOp()/initializeUnorderedBulkOp() chains.
func (p *Parser) parseBulkStatement(collName string, collLoc ast.Loc, accessMethod string, stmtStart int) (ast.Node, error) {
	return nil, p.syntaxErrorAtCur()
}
