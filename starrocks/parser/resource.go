package parser

import (
	"github.com/bytebase/omni/starrocks/ast"
)

// ---------------------------------------------------------------------------
// RESOURCE
// ---------------------------------------------------------------------------

// parseCreateResource parses:
//
//	CREATE [EXTERNAL] RESOURCE [IF NOT EXISTS] name PROPERTIES(...)
//
// CREATE has already been consumed; cur is RESOURCE (for non-EXTERNAL) or
// EXTERNAL (the caller already peeked). The RESOURCE (and optional EXTERNAL)
// token(s) are consumed here. external indicates whether EXTERNAL was already
// seen/consumed by the caller.
func (p *Parser) parseCreateResource(startLoc ast.Loc, external bool) (ast.Node, error) {
	// cur is RESOURCE
	if _, err := p.expect(kwRESOURCE); err != nil {
		return nil, err
	}

	stmt := &ast.CreateResourceStmt{External: external}

	// Optional IF NOT EXISTS
	if p.cur.Kind == kwIF {
		p.advance() // consume IF
		if _, err := p.expect(kwNOT); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfNotExists = true
	}

	// Resource name — can be a string literal or identifier
	name, nameLoc, err := p.parseIdentifierOrString()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	endLoc := nameLoc

	// PROPERTIES clause
	if p.cur.Kind == kwPROPERTIES {
		props, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		stmt.Properties = props
		if len(props) > 0 {
			endLoc = ast.NodeLoc(props[len(props)-1])
		}
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseAlterResource parses alterResourceStatement (StarRocks.g4:1388):
//
//	ALTER RESOURCE resourceName=identifierOrString SET properties
//
// ALTER has already been consumed; cur is RESOURCE. RESOURCE is consumed here.
func (p *Parser) parseAlterResource(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume RESOURCE

	stmt := &ast.AlterResourceStmt{}

	// Resource name
	name, nameLoc, err := p.parseIdentifierOrString()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	endLoc := nameLoc

	if _, err := p.expect(kwSET); err != nil {
		return nil, err
	}
	if p.cur.Kind != kwPROPERTIES {
		return nil, p.syntaxErrorAtCur()
	}
	props, err := p.parseProperties()
	if err != nil {
		return nil, err
	}
	stmt.Properties = props
	endLoc = p.prev.Loc

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseDropResource parses:
//
//	DROP RESOURCE [IF EXISTS] name
//
// DROP has already been consumed; cur is RESOURCE. RESOURCE is consumed here.
func (p *Parser) parseDropResource(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume RESOURCE

	stmt := &ast.DropResourceStmt{}

	// Optional IF EXISTS
	if p.cur.Kind == kwIF {
		p.advance() // consume IF
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}

	// Resource name
	name, nameLoc, err := p.parseIdentifierOrString()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	stmt.Loc = startLoc.Merge(nameLoc)
	return stmt, nil
}
