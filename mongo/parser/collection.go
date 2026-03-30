package parser

import "github.com/bytebase/omni/mongo/ast"

// parseCollectionStatement parses db.collection.method().cursor()... chains.
// The collection name, its location, access method, and statement start have already been determined.
func (p *Parser) parseCollectionStatement(collName string, collLoc ast.Loc, accessMethod string, stmtStart int) (ast.Node, error) {
	// Check for bulk operations
	if p.cur.Type == kwInitializeOrderedBulkOp || p.cur.Type == kwInitializeUnorderedBulkOp {
		return p.parseBulkStatement(collName, collLoc, accessMethod, stmtStart)
	}

	// Check for plan cache
	if p.cur.Type == kwGetPlanCache {
		return p.parsePlanCacheStatement(collName, collLoc, accessMethod, stmtStart)
	}

	// Check for explain prefix
	var explain bool
	var explainArgs []ast.Node
	if p.cur.Type == kwExplain {
		explain = true
		p.advance() // consume 'explain'
		eArgs, err := p.parseArguments()
		if err != nil {
			return nil, err
		}
		explainArgs = eArgs

		// After explain(), if no '.', this is a standalone explain call
		if p.cur.Type != '.' {
			return &ast.CollectionStatement{
				Collection:    collName,
				CollectionLoc: collLoc,
				AccessMethod:  accessMethod,
				Method:        "explain",
				Args:          explainArgs,
				Explain:       false,
				Loc:           ast.Loc{Start: stmtStart, End: p.prev.End},
			}, nil
		}
		p.advance() // consume '.'
	}

	// Parse primary method: consume method name
	methodName, methodTok, err := p.expectIdent()
	if err != nil {
		return nil, err
	}
	_ = methodTok

	// Expect '(' and parse arguments
	args, err := p.parseArguments()
	if err != nil {
		return nil, err
	}

	// Parse cursor method chain: while next is '.', parse chained methods
	var cursorMethods []ast.CursorMethod
	for p.cur.Type == '.' {
		dotTok := p.advance() // consume '.'
		_ = dotTok

		chainName, chainTok, err := p.expectIdent()
		if err != nil {
			return nil, err
		}
		chainStart := chainTok.Loc

		chainArgs, err := p.parseArguments()
		if err != nil {
			return nil, err
		}

		cursorMethods = append(cursorMethods, ast.CursorMethod{
			Method: chainName,
			Args:   chainArgs,
			Loc:    ast.Loc{Start: chainStart, End: p.prev.End},
		})
	}

	return &ast.CollectionStatement{
		Collection:    collName,
		CollectionLoc: collLoc,
		AccessMethod:  accessMethod,
		Method:        methodName,
		Args:          args,
		CursorMethods: cursorMethods,
		Explain:       explain,
		ExplainArgs:   explainArgs,
		Loc:           ast.Loc{Start: stmtStart, End: p.prev.End},
	}, nil
}
