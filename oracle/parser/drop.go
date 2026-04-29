package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseDropStmt parses a DROP statement.
//
// BNF: oracle/parser/bnf/DROP-TABLE.bnf, DROP-VIEW.bnf, DROP-INDEX.bnf,
//
//	     DROP-SEQUENCE.bnf, DROP-SYNONYM.bnf, DROP-DATABASE-LINK.bnf,
//	     DROP-FUNCTION.bnf, DROP-PROCEDURE.bnf, DROP-PACKAGE.bnf,
//	     DROP-TRIGGER.bnf, DROP-TYPE.bnf
//
//		DROP TABLE [ IF EXISTS ] [ schema. ] table [ CASCADE CONSTRAINTS ] [ PURGE ] ;
//		DROP VIEW [ IF EXISTS ] [ schema. ] view [ CASCADE CONSTRAINTS ] ;
//		DROP INDEX [ IF EXISTS ] [ schema. ] index_name [ ONLINE ] [ FORCE ]
//		    [ { DEFERRED | IMMEDIATE } INVALIDATION ] ;
//		DROP SEQUENCE [ IF EXISTS ] [ schema. ] sequence_name ;
//		DROP [ PUBLIC ] SYNONYM [ IF EXISTS ] [ schema. ] synonym [ FORCE ] ;
//		DROP [ PUBLIC ] DATABASE LINK [ IF EXISTS ] dblink ;
//		DROP FUNCTION [ IF EXISTS ] [ schema. ] function_name ;
//		DROP PROCEDURE [ IF EXISTS ] [ schema. ] procedure ;
//		DROP PACKAGE [ BODY ] [ IF EXISTS ] [ schema. ] package ;
//		DROP TRIGGER [ IF EXISTS ] [ schema. ] trigger ;
//		DROP TYPE [ IF EXISTS ] [ schema. ] type_name [ FORCE | VALIDATE ] ;
func (p *Parser) parseDropStmt() (nodes.StmtNode, error) {
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
		p.advance() // consume MATERIALIZED
		// Check for MATERIALIZED ZONEMAP
		if p.isIdentLike() && p.cur.Str == "ZONEMAP" {
			p.advance() // consume ZONEMAP
			return p.parseDropSimpleStmt(nodes.OBJECT_MATERIALIZED_ZONEMAP, start)
		}
		if p.cur.Type == kwVIEW {
			p.advance() // consume VIEW
			// Check for MATERIALIZED VIEW LOG
			if p.cur.Type == kwLOG {
				p.advance() // consume LOG
				return p.parseDropMaterializedViewLogStmt(start)
			}
		}
		stmt.ObjectType = nodes.OBJECT_MATERIALIZED_VIEW
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
			stmt.ObjectType = nodes.OBJECT_DATABASE_LINK
		} else {
			// DROP DATABASE (no LINK)
			stmt.ObjectType = nodes.OBJECT_DATABASE
			stmt.Loc.End = p.prev.End
			return stmt, nil
		}
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
	case kwAUDIT:
		// DROP AUDIT POLICY
		p.advance() // consume AUDIT
		if p.isIdentLikeStr("POLICY") {
			p.advance() // consume POLICY
		}
		return p.parseDropAuditPolicyStmt(start)
	case kwJSON:
		// DROP JSON RELATIONAL DUALITY VIEW
		p.advance() // consume JSON
		if p.isIdentLike() && p.cur.Str == "RELATIONAL" {
			p.advance() // consume RELATIONAL
		}
		if p.isIdentLike() && p.cur.Str == "DUALITY" {
			p.advance() // consume DUALITY
		}
		if p.cur.Type == kwVIEW {
			p.advance() // consume VIEW
		}
		return p.parseDropSimpleStmt(nodes.OBJECT_JSON_DUALITY_VIEW, start)
	case kwUSER, kwROLE, kwPROFILE,
		kwTABLESPACE, kwDIRECTORY, kwCONTEXT,
		kwCLUSTER, kwJAVA, kwLIBRARY:
		adminStmt, err := p.parseDropAdminObject(start)
		if err != nil {
			return nil, err
		}
		if adminStmt != nil {
			return adminStmt, nil
		}
		stmt.ObjectType = nodes.OBJECT_TABLE
	default:
		// Check for DIMENSION, FLASHBACK ARCHIVE
		if p.isIdentLike() {
			adminStmt, err := p.parseDropAdminObject(start)
			if err != nil {
				return nil, err
			}
			if adminStmt != nil {
				return adminStmt, nil
			}
		}
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
	name, parseErr692 := p.parseObjectName()
	if parseErr692 != nil {
		return nil, parseErr692
	}
	stmt.Names.Items = append(stmt.Names.Items, name)

	// Optional trailing options: CASCADE CONSTRAINTS, PURGE, FORCE, ONLINE, DEFERRED/IMMEDIATE INVALIDATION
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
			stmt.Force = true
		} else if p.isIdentLikeStr("VALIDATE") {
			p.advance()
			stmt.Validate = true
		} else if p.cur.Type == kwONLINE {
			p.advance()
			stmt.Online = true
		} else if p.cur.Type == kwDEFERRED {
			stmt.Invalidation = "DEFERRED"
			p.advance() // consume DEFERRED
			if p.isIdentLikeStr("INVALIDATION") {
				p.advance() // consume INVALIDATION
			}
		} else if p.cur.Type == kwIMMEDIATE {
			stmt.Invalidation = "IMMEDIATE"
			p.advance() // consume IMMEDIATE
			if p.isIdentLikeStr("INVALIDATION") {
				p.advance() // consume INVALIDATION
			}
		} else {
			break
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseDropSimpleStmt parses a simple DROP statement for objects like INDEXTYPE and OPERATOR.
// Called after DROP and the object type keyword have been consumed.
//
//	DROP { INDEXTYPE | OPERATOR } [ IF EXISTS ] [ schema. ] name [ FORCE ]
func (p *Parser) parseDropSimpleStmt(objType nodes.ObjectType, start int) (*nodes.DropStmt, error) {
	stmt := &nodes.DropStmt{
		ObjectType: objType,
		Names:      &nodes.List{},
		Loc:        nodes.Loc{Start: start},
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
	name, parseErr693 := p.parseObjectName()
	if parseErr693 != nil {
		return nil, parseErr693
	}
	stmt.Names.Items = append(stmt.Names.Items, name)

	// Optional FORCE
	if p.cur.Type == kwFORCE {
		p.advance()
		stmt.Force = true
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseDropMaterializedViewLogStmt parses a DROP MATERIALIZED VIEW LOG statement.
// Called after DROP MATERIALIZED VIEW LOG has been consumed.
//
// BNF: oracle/parser/bnf/DROP-MATERIALIZED-VIEW-LOG.bnf
//
//	DROP MATERIALIZED VIEW LOG [ IF EXISTS ] ON [ schema. ] table_name ;
func (p *Parser) parseDropMaterializedViewLogStmt(start int) (*nodes.DropStmt, error) {
	stmt := &nodes.DropStmt{
		ObjectType: nodes.OBJECT_MATERIALIZED_VIEW_LOG,
		Names:      &nodes.List{},
		Loc:        nodes.Loc{Start: start},
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

	// ON keyword
	if p.cur.Type == kwON {
		p.advance() // consume ON
	}

	// Parse table name
	name, parseErr694 := p.parseObjectName()
	if parseErr694 != nil {
		return nil, parseErr694
	}
	stmt.Names.Items = append(stmt.Names.Items, name)

	stmt.Loc.End = p.prev.End
	return stmt, nil
}
