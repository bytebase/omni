package parser

import "github.com/bytebase/omni/mongo/ast"

// parsePlanCacheStatement parses db.collection.getPlanCache() chains.
// Already called with collName, collLoc, accessMethod, stmtStart determined.
// Current token is kwGetPlanCache.
func (p *Parser) parsePlanCacheStatement(collName string, collLoc ast.Loc, accessMethod string, stmtStart int) (ast.Node, error) {
	p.advance() // consume 'getPlanCache'

	// Expect ()
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	// Parse optional method chain
	var chain []ast.MethodCall
	for p.cur.Type == '.' {
		p.advance() // consume '.'
		mc, err := p.parseMethodCall()
		if err != nil {
			return nil, err
		}
		chain = append(chain, *mc)
	}

	return &ast.PlanCacheStatement{
		Collection:     collName,
		AccessMethod:   accessMethod,
		ChainedMethods: chain,
		Loc:            ast.Loc{Start: stmtStart, End: p.prev.End},
	}, nil
}
