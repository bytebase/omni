package parser

import (
	"github.com/bytebase/omni/starrocks/ast"
)

// ---------------------------------------------------------------------------
// WAREHOUSE
// ---------------------------------------------------------------------------

// parseCreateWarehouse parses createWarehouseStatement (StarRocks.g4:2038):
//
//	CREATE WAREHOUSE [IF NOT EXISTS] name [COMMENT 'text'] [PROPERTIES (...)]
//
// CREATE has already been consumed; cur is WAREHOUSE.
func (p *Parser) parseCreateWarehouse(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume WAREHOUSE

	stmt := &ast.CreateWarehouseStmt{}

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

	name, nameLoc, err := p.parseIdentifierOrString()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	endLoc := nameLoc

	// comment (StarRocks.g4:2734): COMMENT string.
	if p.cur.Kind == kwCOMMENT {
		p.advance() // consume COMMENT
		if p.cur.Kind != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Comment = p.cur.Str
		endLoc = p.cur.Loc
		p.advance()
	}

	if p.cur.Kind == kwPROPERTIES {
		props, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		stmt.Properties = props
		endLoc = p.prev.Loc
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseDropWarehouse parses dropWarehouseStatement (StarRocks.g4:2043):
//
//	DROP WAREHOUSE [IF EXISTS] name
//
// DROP has already been consumed; cur is WAREHOUSE.
func (p *Parser) parseDropWarehouse(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume WAREHOUSE

	ifExists, err := p.parseIfExists()
	if err != nil {
		return nil, err
	}

	name, nameLoc, err := p.parseIdentifierOrString()
	if err != nil {
		return nil, err
	}

	stmt := &ast.DropWarehouseStmt{
		Name:     name,
		IfExists: ifExists,
		Loc:      startLoc.Merge(nameLoc),
	}
	return stmt, nil
}

// parseSuspendWarehouse parses suspendWarehouseStatement (StarRocks.g4:2047):
//
//	SUSPEND WAREHOUSE [IF EXISTS] name
//
// SUSPEND has already been consumed; cur is WAREHOUSE.
func (p *Parser) parseSuspendWarehouse(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume WAREHOUSE

	ifExists, err := p.parseIfExists()
	if err != nil {
		return nil, err
	}

	// suspendWarehouseStatement takes a plain identifier, not identifierOrString.
	name, nameLoc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	stmt := &ast.SuspendWarehouseStmt{
		Name:     name,
		IfExists: ifExists,
		Loc:      startLoc.Merge(nameLoc),
	}
	return stmt, nil
}

// parseResumeWarehouse parses resumeWarehouseStatement (StarRocks.g4:2051):
//
//	RESUME WAREHOUSE [IF EXISTS] name
//
// RESUME has already been consumed; cur is WAREHOUSE.
func (p *Parser) parseResumeWarehouse(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume WAREHOUSE

	ifExists, err := p.parseIfExists()
	if err != nil {
		return nil, err
	}

	// resumeWarehouseStatement takes a plain identifier, not identifierOrString.
	name, nameLoc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	stmt := &ast.ResumeWarehouseStmt{
		Name:     name,
		IfExists: ifExists,
		Loc:      startLoc.Merge(nameLoc),
	}
	return stmt, nil
}

// parseIfExists reads the optional IF EXISTS of a statement and reports
// whether it was there.
func (p *Parser) parseIfExists() (bool, error) {
	if p.cur.Kind != kwIF {
		return false, nil
	}
	p.advance() // consume IF
	if _, err := p.expect(kwEXISTS); err != nil {
		return false, err
	}
	return true, nil
}
