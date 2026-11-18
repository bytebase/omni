package parser

import (
	"github.com/bytebase/omni/starrocks/ast"
)

// parseCreateDatabase parses createDbStatement (StarRocks.g4:354):
//
//	CREATE (DATABASE | SCHEMA) [IF NOT EXISTS] db_name
//	    [charsetDesc] [collateDesc] [PROPERTIES ("key"="value", ...)]
//
// The CREATE keyword has already been consumed by the caller; cur is
// DATABASE or SCHEMA.
func (p *Parser) parseCreateDatabase() (ast.Node, error) {
	startLoc := p.prev.Loc // loc of CREATE token

	// Consume DATABASE or SCHEMA
	p.advance()

	stmt := &ast.CreateDatabaseStmt{}

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

	// Database name — accepts a plain identifier (or non-reserved keyword as name)
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	endLoc := ast.NodeLoc(name)

	// Optional charsetDesc, then optional collateDesc. DEFAULT belongs to the
	// clause it precedes, so each is looked for on its own.
	if p.atCharsetDesc() {
		charset, loc, err := p.parseCharsetDesc()
		if err != nil {
			return nil, err
		}
		stmt.Charset = charset
		endLoc = loc
	}
	if p.atCollateDesc() {
		collate, loc, err := p.parseCollateDesc()
		if err != nil {
			return nil, err
		}
		stmt.Collate = collate
		endLoc = loc
	}

	// Optional PROPERTIES clause
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

// atCharsetDesc reports whether charsetDesc (StarRocks.g4:424)
//
//	[DEFAULT] (CHAR SET | CHARSET | CHARACTER SET) [=] identifierOrString
//
// starts at the current token. CHAR is a type name too, so CHAR SET is read as
// one word only when SET follows it.
func (p *Parser) atCharsetDesc() bool {
	switch p.cur.Kind {
	case kwCHARSET:
		return true
	case kwCHAR, kwCHARACTER:
		return p.peekNext().Kind == kwSET
	case kwDEFAULT:
		// DEFAULT belongs to the clause after it, so the clause is looked for
		// past it: a speculative advance and rollback, since CHAR SET needs
		// two tokens.
		mark := p.save()
		p.advance() // consume DEFAULT
		ok := p.atCharsetDesc()
		p.restore(mark)
		return ok
	}
	return false
}

// atCollateDesc reports whether collateDesc (StarRocks.g4:428)
//
//	[DEFAULT] COLLATE [=] identifierOrString
//
// starts at the current token.
func (p *Parser) atCollateDesc() bool {
	if p.cur.Kind == kwCOLLATE {
		return true
	}
	return p.cur.Kind == kwDEFAULT && p.peekNext().Kind == kwCOLLATE
}

// parseCharsetDesc parses charsetDesc (StarRocks.g4:424). cur is one of the
// tokens atCharsetDesc accepted.
func (p *Parser) parseCharsetDesc() (string, ast.Loc, error) {
	p.match(kwDEFAULT)
	switch p.cur.Kind {
	case kwCHARSET:
		p.advance() // consume CHARSET
	case kwCHAR, kwCHARACTER:
		p.advance() // consume CHAR / CHARACTER
		if _, err := p.expect(kwSET); err != nil {
			return "", ast.NoLoc(), err
		}
	}
	p.match(int('='))
	return p.parseIdentifierOrString()
}

// parseCollateDesc parses collateDesc (StarRocks.g4:428). cur is one of the
// tokens atCollateDesc accepted.
func (p *Parser) parseCollateDesc() (string, ast.Loc, error) {
	p.match(kwDEFAULT)
	if _, err := p.expect(kwCOLLATE); err != nil {
		return "", ast.NoLoc(), err
	}
	p.match(int('='))
	return p.parseIdentifierOrString()
}

// parseAlterDatabase parses:
//
//	ALTER DATABASE db_name RENAME [TO] new_name
//	ALTER DATABASE db_name SET PROPERTIES ("key"="value", ...)
//	ALTER DATABASE db_name SET QUOTA ...  (quota value consumed as best-effort)
//
// The ALTER keyword has already been consumed; cur is DATABASE.
func (p *Parser) parseAlterDatabase() (ast.Node, error) {
	startLoc := p.prev.Loc // loc of ALTER token

	// Consume DATABASE
	p.advance()

	stmt := &ast.AlterDatabaseStmt{}

	// Database name
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	endLoc := ast.NodeLoc(name)

	switch p.cur.Kind {
	case kwRENAME:
		p.advance() // consume RENAME
		// Optional TO keyword
		p.match(kwTO)
		newName, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.NewName = newName
		endLoc = ast.NodeLoc(newName)

	case kwSET:
		p.advance() // consume SET
		if p.cur.Kind == kwPROPERTIES {
			props, err := p.parseProperties()
			if err != nil {
				return nil, err
			}
			stmt.Properties = props
			if len(props) > 0 {
				endLoc = ast.NodeLoc(props[len(props)-1])
			}
		} else if p.cur.Kind == kwDATA || p.cur.Kind == kwREPLICA {
			// StarRocks.g4:348 alterDbQuotaStatement:
			//	ALTER DATABASE identifier SET DATA QUOTA identifier
			//	ALTER DATABASE identifier SET REPLICA QUOTA INTEGER_VALUE
			data := p.cur.Kind == kwDATA
			p.advance() // consume DATA / REPLICA
			if _, err := p.expect(kwQUOTA); err != nil {
				return nil, err
			}
			if data {
				// The quota is an identifier (10T lexes as a DIGIT_IDENTIFIER);
				// a bare byte count, as in SET DATA QUOTA 10995116277760, is
				// kept: no compat row settles it.
				if p.cur.Kind == tokInt {
					p.advance()
				} else if _, _, err := p.parseIdentifier(); err != nil {
					return nil, err
				}
			} else if _, err := p.alterExpectInt(); err != nil {
				return nil, err
			}
			endLoc = p.prev.Loc
		} else {
			return nil, p.syntaxErrorAtCur()
		}

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseDropDatabase parses:
//
//	DROP (DATABASE | SCHEMA) [IF EXISTS] db_name [FORCE]
//
// The DROP keyword has already been consumed; cur is DATABASE or SCHEMA.
func (p *Parser) parseDropDatabase() (ast.Node, error) {
	startLoc := p.prev.Loc // loc of DROP token

	// Consume DATABASE or SCHEMA
	p.advance()

	stmt := &ast.DropDatabaseStmt{}

	// Optional IF EXISTS
	if p.cur.Kind == kwIF {
		p.advance() // consume IF
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}

	// Database name
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	endLoc := ast.NodeLoc(name)

	// Optional FORCE
	if p.cur.Kind == kwFORCE {
		endLoc = p.cur.Loc
		p.advance()
		stmt.Force = true
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseProperties parses:
//
//	PROPERTIES ("key"="value" [, "key"="value" ...])
//
// cur must be kwPROPERTIES on entry; it is consumed here.
func (p *Parser) parseProperties() ([]*ast.Property, error) {
	p.advance() // consume PROPERTIES
	return p.parsePropertyList()
}

// parsePropertyList parses a bare ('key'='value', ...) list; cur must be '('
// on entry. Shared by parseProperties (after the PROPERTIES keyword) and the
// bare index-property list (StarRocks `USING GIN ('k'='v')`).
func (p *Parser) parsePropertyList() ([]*ast.Property, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}

	return p.parsePropertyItems(p.parseStringProperty)
}

// parsePropertyItems parses the body of a property list -- the grammar's
//
//	propertyList : property (COMMA property)*
//
// -- from the first item through the closing ')', which it consumes. cur must
// be on the first item's first token, the caller having consumed the '('.
//
// So the list carries at least one item, items are comma-separated, and it
// does not end on a comma. The engine refuses all three departures --
// PROPERTIES(), PROPERTIES("k"="v" "k2"="v2") and PROPERTIES("k"="v",) --
// in every clause that takes a property list, container-verified across
// CREATE TABLE / DATABASE / EXTERNAL CATALOG / MATERIALIZED VIEW,
// ALTER TABLE ... SET, ROUTINE LOAD's job and data-source lists,
// ADMIN SET FRONTEND CONFIG, the bare USING GIN index list, and
// SELECT ... INTO OUTFILE. Each of those clauses had its own copy of this
// loop treating the comma as optional, so each accepted statements the engine
// rejects; the shape lives here once instead.
//
// parseItem reads one item. The clauses differ in what they take for a key
// and a value, not in the shape of the list.
func (p *Parser) parsePropertyItems(parseItem func() (*ast.Property, error)) ([]*ast.Property, error) {
	var props []*ast.Property
	for {
		if p.cur.Kind == int(')') || p.cur.Kind == tokEOF {
			// An empty list, or one that ended on its separator. The grammar
			// wants an item here, and this is where the engine reports it.
			return nil, p.syntaxErrorAtCur()
		}
		item, err := parseItem()
		if err != nil {
			return nil, err
		}
		props = append(props, item)
		if p.cur.Kind != int(',') {
			// No separator, so the list ends here: a second item with nothing
			// between them is left to expect(')') to refuse.
			break
		}
		p.advance() // consume ','
	}

	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}

	return props, nil
}

// parseStringProperty reads one "key"="value" item, both sides string
// literals. This is the common property item.
func (p *Parser) parseStringProperty() (*ast.Property, error) {
	startLoc := p.cur.Loc

	if p.cur.Kind != tokString {
		return nil, p.syntaxErrorAtCur()
	}
	key := p.cur.Str
	p.advance()

	if _, err := p.expect(int('=')); err != nil {
		return nil, err
	}

	if p.cur.Kind != tokString {
		return nil, p.syntaxErrorAtCur()
	}
	val := p.cur.Str
	endLoc := p.cur.Loc
	p.advance()

	return &ast.Property{
		Key:   key,
		Value: val,
		Loc:   ast.Loc{Start: startLoc.Start, End: endLoc.End},
	}, nil
}
