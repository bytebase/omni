package parser

import (
	"github.com/bytebase/omni/starrocks/ast"
)

// parseCreateCatalog parses:
//
//	CREATE [EXTERNAL] CATALOG [IF NOT EXISTS] catalog_name
//	    [COMMENT 'comment']
//	    [PROPERTIES ("key"="value", ...)]
//	    [WITH RESOURCE resource_name]
//
// The CREATE keyword has already been consumed by the caller.
// cur may be kwEXTERNAL or kwCATALOG on entry.
func (p *Parser) parseCreateCatalog() (ast.Node, error) {
	startLoc := p.prev.Loc // loc of CREATE token

	stmt := &ast.CreateCatalogStmt{}

	// Optional EXTERNAL
	// createExternalCatalogStatement (StarRocks.g4:760) is the only CREATE
	// CATALOG rule: the pinned server answers `CREATE CATALOG a` with "No
	// viable statement for input 'CREATE CATALOG'", and its PROPERTIES are not
	// optional.
	if p.cur.Kind != kwEXTERNAL {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.External = true
	p.advance() // consume EXTERNAL

	// CATALOG keyword
	if _, err := p.expect(kwCATALOG); err != nil {
		return nil, err
	}

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

	// Catalog name — single identifier
	name, nameLoc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	endLoc := nameLoc

	// Optional clauses: COMMENT, PROPERTIES, WITH RESOURCE (in any order per Doris docs)
	for {
		switch p.cur.Kind {
		case kwCOMMENT:
			p.advance() // consume COMMENT
			if p.cur.Kind != tokString {
				return nil, p.syntaxErrorAtCur()
			}
			stmt.Comment = p.cur.Str
			endLoc = p.cur.Loc
			p.advance()
		case kwPROPERTIES:
			props, err := p.parseProperties()
			if err != nil {
				return nil, err
			}
			stmt.Properties = props
			if len(props) > 0 {
				endLoc = ast.NodeLoc(props[len(props)-1])
			}
		default:
			goto done
		}
	}
done:
	if stmt.Properties == nil {
		// properties is not optional in createExternalCatalogStatement: the
		// pinned server refuses `CREATE EXTERNAL CATALOG a`.
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseAlterCatalog parses alterCatalogStatement (StarRocks.g4:2152):
//
//	ALTER CATALOG catalog_name SET '(' property [',' property]* ')'
//
// The ALTER keyword has already been consumed; cur is CATALOG.
func (p *Parser) parseAlterCatalog() (ast.Node, error) {
	startLoc := p.prev.Loc // loc of ALTER token

	// Consume CATALOG
	p.advance()

	stmt := &ast.AlterCatalogStmt{}

	name, _, err := p.parseIdentifierOrString()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	if _, err := p.expect(kwSET); err != nil {
		return nil, err
	}
	props, err := p.parsePropertyList()
	if err != nil {
		return nil, err
	}
	stmt.Properties = props

	stmt.Loc = startLoc.Merge(p.prev.Loc)
	return stmt, nil
}

// parseDropCatalog parses:
//
//	DROP CATALOG [IF EXISTS] catalog_name
//
// The DROP keyword has already been consumed; cur is CATALOG.
func (p *Parser) parseDropCatalog() (ast.Node, error) {
	startLoc := p.prev.Loc // loc of DROP token

	// Consume CATALOG
	p.advance()

	stmt := &ast.DropCatalogStmt{}

	// Optional IF EXISTS
	if p.cur.Kind == kwIF {
		p.advance() // consume IF
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}

	// Catalog name — single identifier
	name, nameLoc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	stmt.Loc = startLoc.Merge(nameLoc)
	return stmt, nil
}

// parseRefreshCatalog parses:
//
//	REFRESH CATALOG catalog_name [PROPERTIES(...)]
//
// The REFRESH keyword has already been consumed; cur is CATALOG.
func (p *Parser) parseRefreshCatalog() (ast.Node, error) {
	startLoc := p.prev.Loc // loc of REFRESH token

	// Consume CATALOG
	p.advance()

	stmt := &ast.RefreshCatalogStmt{}

	// Catalog name — single identifier
	name, nameLoc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	endLoc := nameLoc

	// Optional PROPERTIES
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

// parseRefreshExternalTable parses refreshTableStatement (StarRocks.g4:550):
//
//	REFRESH EXTERNAL TABLE qualifiedName [(PARTITION '(' string [',' string]* ')')]
//
// On entry REFRESH has been consumed; cur is EXTERNAL.
func (p *Parser) parseRefreshExternalTable(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume EXTERNAL
	if _, err := p.expect(kwTABLE); err != nil {
		return nil, err
	}
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	// The AST builder reads the name as a table name, which is a table, a
	// database and a table, or a catalog, a database and a table: it answers
	// `REFRESH EXTERNAL TABLE a.b.c.d` with "Invalid table name format".
	if len(name.Parts) > 3 {
		return nil, &ParseError{
			Position: name.Loc.Start, End: name.Loc.End,
			Message: "REFRESH EXTERNAL TABLE names at most a catalog, a database and a table",
		}
	}
	stmt := &ast.RefreshExternalTableStmt{Name: name}
	endLoc := name.Loc
	if p.cur.Kind == kwPARTITION {
		p.advance() // consume PARTITION
		if _, err := p.expect(int('(')); err != nil {
			return nil, err
		}
		for {
			if p.cur.Kind != tokString {
				return nil, p.syntaxErrorAtCur()
			}
			stmt.Partitions = append(stmt.Partitions, p.cur.Str)
			endLoc = p.cur.Loc
			p.advance()
			if p.cur.Kind != int(',') {
				break
			}
			p.advance() // consume ','
		}
		closeTok, err := p.expect(int(')'))
		if err != nil {
			return nil, err
		}
		endLoc = closeTok.Loc
	}
	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}
