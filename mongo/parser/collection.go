package parser

import "github.com/bytebase/omni/mongo/ast"

// parseCollectionStatement parses db.collection.method().cursor()... chains.
func (p *Parser) parseCollectionStatement(collName string, collLoc ast.Loc, accessMethod string, stmtStart int) (ast.Node, error) {
	return nil, p.syntaxErrorAtCur()
}
