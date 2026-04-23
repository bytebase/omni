package parser

import (
	"github.com/bytebase/omni/doris/ast"
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
	if p.cur.Kind == kwEXTERNAL {
		stmt.External = true
		p.advance()
	}

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
		case kwWITH:
			p.advance() // consume WITH
			if _, err := p.expect(kwRESOURCE); err != nil {
				return nil, err
			}
			resourceName, resourceLoc, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.WithResource = resourceName
			endLoc = resourceLoc
		default:
			goto done
		}
	}
done:
	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseAlterCatalog parses:
//
//	ALTER CATALOG catalog_name
//	    { RENAME new_name
//	    | SET PROPERTIES ("key"="value", ...)
//	    | MODIFY COMMENT 'text'
//	    | PROPERTY ("key"="value") }
//
// The ALTER keyword has already been consumed; cur is CATALOG.
func (p *Parser) parseAlterCatalog() (ast.Node, error) {
	startLoc := p.prev.Loc // loc of ALTER token

	// Consume CATALOG
	p.advance()

	stmt := &ast.AlterCatalogStmt{}

	// Catalog name — single identifier
	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	endLoc := startLoc

	switch p.cur.Kind {
	case kwRENAME:
		p.advance() // consume RENAME
		newName, newNameLoc, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterCatalogRename
		stmt.NewName = newName
		endLoc = newNameLoc

	case kwSET:
		p.advance() // consume SET
		if p.cur.Kind != kwPROPERTIES {
			return nil, p.syntaxErrorAtCur()
		}
		props, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterCatalogSetProperties
		stmt.Properties = props
		if len(props) > 0 {
			endLoc = ast.NodeLoc(props[len(props)-1])
		}

	case kwMODIFY:
		p.advance() // consume MODIFY
		if _, err := p.expect(kwCOMMENT); err != nil {
			return nil, err
		}
		if p.cur.Kind != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Action = ast.AlterCatalogModifyComment
		stmt.Comment = p.cur.Str
		endLoc = p.cur.Loc
		p.advance()

	case kwPROPERTY:
		// PROPERTY ("key"="value") — singular, no SET prefix
		props, err := p.parsePropertyClause()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterCatalogSetProperty
		stmt.Properties = props
		if len(props) > 0 {
			endLoc = ast.NodeLoc(props[len(props)-1])
		}

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc = startLoc.Merge(endLoc)
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

// parsePropertyClause parses:
//
//	PROPERTY ("key"="value" [, "key"="value" ...])
//
// cur must be kwPROPERTY on entry; it is consumed here.
// This is the singular form used in ALTER CATALOG ... PROPERTY (...).
func (p *Parser) parsePropertyClause() ([]*ast.Property, error) {
	p.advance() // consume PROPERTY

	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}

	var props []*ast.Property

	for p.cur.Kind != int(')') && p.cur.Kind != tokEOF {
		startLoc := p.cur.Loc

		// Key — must be a string literal
		if p.cur.Kind != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		key := p.cur.Str
		p.advance()

		// '='
		if _, err := p.expect(int('=')); err != nil {
			return nil, err
		}

		// Value — must be a string literal
		if p.cur.Kind != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		val := p.cur.Str
		endLoc := p.cur.Loc
		p.advance()

		props = append(props, &ast.Property{
			Key:   key,
			Value: val,
			Loc:   ast.Loc{Start: startLoc.Start, End: endLoc.End},
		})

		// Optional comma separator
		if p.cur.Kind == int(',') {
			p.advance()
		}
	}

	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}

	return props, nil
}
