package parser

import (
	"github.com/bytebase/omni/starrocks/ast"
)

// parseCreateDatabase parses:
//
//	CREATE (DATABASE | SCHEMA) [IF NOT EXISTS] db_name
//	    [PROPERTIES ("key"="value", ...)]
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

	// Optional PROPERTIES clause
	if p.cur.Kind == kwPROPERTIES {
		props, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		stmt.Properties = props
	}

	stmt.Loc = startLoc.Merge(ast.NodeLoc(stmt.Name))
	return stmt, nil
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
		} else if p.cur.Kind == kwQUOTA || p.cur.Kind == kwDATA || p.cur.Kind == kwREPLICA {
			// Consume the rest of the SET [DATA|REPLICA] QUOTA clause as a
			// best-effort skip. Quota value can be a number + unit (e.g., 10GB).
			p.advance() // consume DATA/REPLICA/QUOTA keyword
			if p.cur.Kind == kwQUOTA {
				p.advance() // consume QUOTA after DATA/REPLICA
			}
			if p.cur.Kind != tokEOF && p.cur.Kind != int(';') {
				p.advance() // consume the quota value token
			}
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
