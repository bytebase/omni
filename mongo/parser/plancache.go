package parser

import "github.com/bytebase/omni/mongo/ast"

// parsePlanCacheStatement parses db.collection.getPlanCache() chains.
func (p *Parser) parsePlanCacheStatement(database, collection string, startLoc int) (ast.Node, error) {
	return nil, p.syntaxErrorAtCur()
}
