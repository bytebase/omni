package parser

import "github.com/bytebase/omni/mongo/ast"

// parseBulkStatement parses db.collection.initializeOrderedBulkOp()/initializeUnorderedBulkOp() chains.
// Already called with collName, collLoc, accessMethod, stmtStart determined.
// Current token is kwInitializeOrderedBulkOp or kwInitializeUnorderedBulkOp.
func (p *Parser) parseBulkStatement(collName string, collLoc ast.Loc, accessMethod string, stmtStart int) (ast.Node, error) {
	ordered := p.cur.Type == kwInitializeOrderedBulkOp
	p.advance() // consume initializeOrderedBulkOp / initializeUnorderedBulkOp

	// Expect ()
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	// Parse operation chain: .method(args)...
	var ops []ast.BulkOperation
	for p.cur.Type == '.' {
		p.advance() // consume '.'

		opName, opTok, err := p.expectIdent()
		if err != nil {
			return nil, err
		}
		opStart := opTok.Loc

		args, err := p.parseArguments()
		if err != nil {
			return nil, err
		}

		ops = append(ops, ast.BulkOperation{
			Method: opName,
			Args:   args,
			Loc:    ast.Loc{Start: opStart, End: p.prev.End},
		})
	}

	return &ast.BulkStatement{
		Collection:   collName,
		AccessMethod: accessMethod,
		Ordered:      ordered,
		Operations:   ops,
		Loc:          ast.Loc{Start: stmtStart, End: p.prev.End},
	}, nil
}
