package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseCreateMaterializedOrView distinguishes between:
// - CREATE MATERIALIZED VIEW LOG ON ... (mview log)
// - CREATE MATERIALIZED VIEW ... (regular mview)
func (p *Parser) parseCreateMaterializedOrView(start int, orReplace bool) nodes.StmtNode {
	p.advance() // consume MATERIALIZED
	if p.cur.Type == kwVIEW {
		next := p.peekNext()
		if next.Type == kwLOG {
			p.advance() // consume VIEW
			p.advance() // consume LOG
			return p.parseAdminDDLStmt("CREATE", nodes.OBJECT_MATERIALIZED_VIEW_LOG, start)
		}
	}
	// It's a regular MATERIALIZED VIEW — but we already consumed MATERIALIZED.
	// parseCreateViewStmt expects to see kwMATERIALIZED. Since we consumed it,
	// we need to handle VIEW directly.
	stmt := &nodes.CreateViewStmt{
		OrReplace:    orReplace,
		Materialized: true,
		Loc:          nodes.Loc{Start: start},
	}
	if p.cur.Type == kwVIEW {
		p.advance()
	}
	// Delegate the rest to the existing view parsing logic.
	return p.finishCreateViewStmt(stmt)
}

// parseAdminDDLStmt parses generic administrative DDL statements by consuming
// the action keyword (CREATE/ALTER/DROP) and object type keyword, then parsing
// the object name and skipping remaining options until semicolon/EOF.
//
// This handles: TABLESPACE, DIRECTORY, CONTEXT, CLUSTER, DIMENSION,
// FLASHBACK ARCHIVE, JAVA, LIBRARY
func (p *Parser) parseAdminDDLStmt(action string, objType nodes.ObjectType, start int) nodes.StmtNode {
	stmt := &nodes.AdminDDLStmt{
		Action:     action,
		ObjectType: objType,
		Loc:        nodes.Loc{Start: start},
	}

	stmt.Name = p.parseObjectName()

	// Skip remaining tokens until ;/EOF
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseCreateAdminObject handles CREATE dispatches for admin DDL objects
// (called from parseCreateStmt after consuming CREATE and modifiers).
func (p *Parser) parseCreateAdminObject(start int) nodes.StmtNode {
	switch p.cur.Type {
	case kwUSER:
		p.advance()
		return p.parseCreateUserStmt(start)
	case kwROLE:
		p.advance()
		return p.parseCreateRoleStmt(start)
	case kwPROFILE:
		p.advance()
		return p.parseCreateProfileStmt(start)
	case kwTABLESPACE:
		p.advance()
		return p.parseAdminDDLStmt("CREATE", nodes.OBJECT_TABLESPACE, start)
	case kwDIRECTORY:
		p.advance()
		return p.parseAdminDDLStmt("CREATE", nodes.OBJECT_DIRECTORY, start)
	case kwCONTEXT:
		p.advance()
		return p.parseAdminDDLStmt("CREATE", nodes.OBJECT_CONTEXT, start)
	case kwCLUSTER:
		p.advance()
		return p.parseAdminDDLStmt("CREATE", nodes.OBJECT_CLUSTER, start)
	case kwJAVA:
		p.advance()
		return p.parseAdminDDLStmt("CREATE", nodes.OBJECT_JAVA, start)
	case kwLIBRARY:
		p.advance()
		return p.parseAdminDDLStmt("CREATE", nodes.OBJECT_LIBRARY, start)
	case kwSCHEMA:
		p.advance()
		return p.parseAdminDDLStmt("CREATE", nodes.OBJECT_TABLE, start) // reuse TABLE type for SCHEMA
	default:
		// DIMENSION, FLASHBACK ARCHIVE handled via identifiers
		if p.isIdentLike() {
			switch p.cur.Str {
			case "DIMENSION":
				p.advance()
				return p.parseAdminDDLStmt("CREATE", nodes.OBJECT_DIMENSION, start)
			case "FLASHBACK":
				p.advance()
				if p.isIdentLike() && p.cur.Str == "ARCHIVE" {
					p.advance()
				}
				return p.parseAdminDDLStmt("CREATE", nodes.OBJECT_FLASHBACK_ARCHIVE, start)
			}
		}
		return nil
	}
}

// parseDropAdminObject handles DROP dispatches for admin DDL objects
// (called from parseDropStmt for object types not handled there).
func (p *Parser) parseDropAdminObject(start int) nodes.StmtNode {
	switch p.cur.Type {
	case kwUSER:
		p.advance()
		stmt := &nodes.AdminDDLStmt{
			Action:     "DROP",
			ObjectType: nodes.OBJECT_USER,
			Loc:        nodes.Loc{Start: start},
		}
		stmt.Name = p.parseObjectName()
		// CASCADE
		if p.cur.Type == kwCASCADE {
			p.advance()
		}
		stmt.Loc.End = p.pos()
		return stmt
	case kwROLE:
		p.advance()
		return p.parseAdminDDLStmt("DROP", nodes.OBJECT_ROLE, start)
	case kwPROFILE:
		p.advance()
		stmt := &nodes.AdminDDLStmt{
			Action:     "DROP",
			ObjectType: nodes.OBJECT_PROFILE,
			Loc:        nodes.Loc{Start: start},
		}
		stmt.Name = p.parseObjectName()
		if p.cur.Type == kwCASCADE {
			p.advance()
		}
		stmt.Loc.End = p.pos()
		return stmt
	case kwTABLESPACE:
		p.advance()
		return p.parseAdminDDLStmt("DROP", nodes.OBJECT_TABLESPACE, start)
	case kwDIRECTORY:
		p.advance()
		return p.parseAdminDDLStmt("DROP", nodes.OBJECT_DIRECTORY, start)
	case kwCONTEXT:
		p.advance()
		return p.parseAdminDDLStmt("DROP", nodes.OBJECT_CONTEXT, start)
	case kwCLUSTER:
		p.advance()
		return p.parseAdminDDLStmt("DROP", nodes.OBJECT_CLUSTER, start)
	case kwJAVA:
		p.advance()
		return p.parseAdminDDLStmt("DROP", nodes.OBJECT_JAVA, start)
	case kwLIBRARY:
		p.advance()
		return p.parseAdminDDLStmt("DROP", nodes.OBJECT_LIBRARY, start)
	default:
		if p.isIdentLike() {
			switch p.cur.Str {
			case "DIMENSION":
				p.advance()
				return p.parseAdminDDLStmt("DROP", nodes.OBJECT_DIMENSION, start)
			case "FLASHBACK":
				p.advance()
				if p.isIdentLike() && p.cur.Str == "ARCHIVE" {
					p.advance()
				}
				return p.parseAdminDDLStmt("DROP", nodes.OBJECT_FLASHBACK_ARCHIVE, start)
			}
		}
		return nil
	}
}

// parseAlterAdminObject handles ALTER dispatches for admin DDL objects.
func (p *Parser) parseAlterAdminObject(start int) nodes.StmtNode {
	switch p.cur.Type {
	case kwUSER:
		p.advance()
		return p.parseAlterUserStmt(start)
	case kwROLE:
		p.advance()
		return p.parseAdminDDLStmt("ALTER", nodes.OBJECT_ROLE, start)
	case kwPROFILE:
		p.advance()
		return p.parseAdminDDLStmt("ALTER", nodes.OBJECT_PROFILE, start)
	case kwTABLESPACE:
		p.advance()
		return p.parseAdminDDLStmt("ALTER", nodes.OBJECT_TABLESPACE, start)
	case kwCLUSTER:
		p.advance()
		return p.parseAdminDDLStmt("ALTER", nodes.OBJECT_CLUSTER, start)
	case kwJAVA:
		p.advance()
		return p.parseAdminDDLStmt("ALTER", nodes.OBJECT_JAVA, start)
	case kwLIBRARY:
		p.advance()
		return p.parseAdminDDLStmt("ALTER", nodes.OBJECT_LIBRARY, start)
	default:
		if p.isIdentLike() {
			switch p.cur.Str {
			case "DIMENSION":
				p.advance()
				return p.parseAdminDDLStmt("ALTER", nodes.OBJECT_DIMENSION, start)
			}
		}
		return nil
	}
}
