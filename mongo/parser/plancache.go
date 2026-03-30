package parser

import "github.com/bytebase/omni/mongo/ast"

// parsePlanCacheStatement parses db.collection.getPlanCache() chains.
func (p *Parser) parsePlanCacheStatement(collName string, collLoc ast.Loc, accessMethod string, stmtStart int) (ast.Node, error) {
	return nil, p.syntaxErrorAtCur()
}
