package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseDropStmt parses a DROP statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/DROP-TABLE.html
//
//	DROP { TABLE | INDEX | VIEW | SEQUENCE | SYNONYM | PACKAGE [BODY] |
//	       PROCEDURE | FUNCTION | TRIGGER | TYPE [BODY] |
//	       MATERIALIZED VIEW | DATABASE LINK }
//	     [IF EXISTS] name
//	     [CASCADE CONSTRAINTS | PURGE]
func (p *Parser) parseDropStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume DROP

	stmt := &nodes.DropStmt{
		Names: &nodes.List{},
		Loc:   nodes.Loc{Start: start},
	}

	// Parse object type
	switch p.cur.Type {
	case kwTABLE:
		stmt.ObjectType = nodes.OBJECT_TABLE
		p.advance()
	case kwINDEX:
		stmt.ObjectType = nodes.OBJECT_INDEX
		p.advance()
	case kwVIEW:
		stmt.ObjectType = nodes.OBJECT_VIEW
		p.advance()
	case kwMATERIALIZED:
		stmt.ObjectType = nodes.OBJECT_MATERIALIZED_VIEW
		p.advance() // consume MATERIALIZED
		if p.cur.Type == kwVIEW {
			p.advance() // consume VIEW
		}
	case kwSEQUENCE:
		stmt.ObjectType = nodes.OBJECT_SEQUENCE
		p.advance()
	case kwSYNONYM:
		stmt.ObjectType = nodes.OBJECT_SYNONYM
		p.advance()
	case kwPACKAGE:
		p.advance() // consume PACKAGE
		if p.cur.Type == kwBODY {
			stmt.ObjectType = nodes.OBJECT_PACKAGE_BODY
			p.advance()
		} else {
			stmt.ObjectType = nodes.OBJECT_PACKAGE
		}
	case kwPROCEDURE:
		stmt.ObjectType = nodes.OBJECT_PROCEDURE
		p.advance()
	case kwFUNCTION:
		stmt.ObjectType = nodes.OBJECT_FUNCTION
		p.advance()
	case kwTRIGGER:
		stmt.ObjectType = nodes.OBJECT_TRIGGER
		p.advance()
	case kwTYPE:
		p.advance() // consume TYPE
		if p.cur.Type == kwBODY {
			stmt.ObjectType = nodes.OBJECT_TYPE_BODY
			p.advance()
		} else {
			stmt.ObjectType = nodes.OBJECT_TYPE
		}
	case kwDATABASE:
		p.advance() // consume DATABASE
		if p.cur.Type == kwLINK {
			p.advance() // consume LINK
		}
		stmt.ObjectType = nodes.OBJECT_DATABASE_LINK
	case kwPUBLIC:
		p.advance() // consume PUBLIC
		// PUBLIC SYNONYM or PUBLIC DATABASE LINK
		if p.cur.Type == kwSYNONYM {
			stmt.ObjectType = nodes.OBJECT_SYNONYM
			p.advance()
		} else if p.cur.Type == kwDATABASE {
			p.advance() // consume DATABASE
			if p.cur.Type == kwLINK {
				p.advance() // consume LINK
			}
			stmt.ObjectType = nodes.OBJECT_DATABASE_LINK
		}
	default:
		// Unknown object type; consume as identifier
		stmt.ObjectType = nodes.OBJECT_TABLE
		p.advance()
	}

	// Optional IF EXISTS
	if p.cur.Type == kwIF {
		next := p.peekNext()
		if next.Type == kwEXISTS {
			p.advance() // consume IF
			p.advance() // consume EXISTS
			stmt.IfExists = true
		}
	}

	// Parse object name
	name := p.parseObjectName()
	stmt.Names.Items = append(stmt.Names.Items, name)

	// Optional trailing options: CASCADE CONSTRAINTS, PURGE, FORCE
	for {
		if p.cur.Type == kwCASCADE {
			p.advance() // consume CASCADE
			stmt.Cascade = true
			if p.cur.Type == kwCONSTRAINTS {
				p.advance() // consume CONSTRAINTS
			}
		} else if p.cur.Type == kwPURGE {
			p.advance()
			stmt.Purge = true
		} else if p.cur.Type == kwFORCE {
			p.advance()
		} else {
			break
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}
