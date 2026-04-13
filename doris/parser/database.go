package parser

import (
	"github.com/bytebase/omni/doris/ast"
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
		} else if p.cur.Kind == kwQUOTA {
			// Consume the rest of the SET QUOTA clause as a best-effort skip.
			// Quota value can be a number + unit (e.g., 10GB) — consume one token
			// as the value and leave anything else for the next statement boundary.
			p.advance() // consume QUOTA keyword
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
		keyEnd := p.cur.Loc.End
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

		_ = keyEnd
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
