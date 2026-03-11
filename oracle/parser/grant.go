package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseGrantStmt parses a GRANT statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/GRANT.html
//
//	GRANT { priv [, priv ...] | ALL [PRIVILEGES] }
//	  ON [object_type] object
//	  TO grantee [, grantee ...]
//	  [WITH GRANT OPTION]
//	  [WITH ADMIN OPTION]
//
//	GRANT role [, role ...] TO user [, user ...]  (role grants — no ON clause)
func (p *Parser) parseGrantStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume GRANT

	stmt := &nodes.GrantStmt{
		Privileges: &nodes.List{},
		Grantees:   &nodes.List{},
		Loc:        nodes.Loc{Start: start},
	}

	// Parse privileges or ALL [PRIVILEGES].
	stmt.AllPriv, stmt.Privileges = p.parsePrivilegeList()

	// ON clause (optional — absent means role grant or system privilege).
	if p.cur.Type == kwON {
		p.advance() // consume ON

		// Optional object type keyword.
		stmt.OnType = p.parseOptionalObjectType()

		// Object name.
		stmt.OnObject = p.parseObjectName()
	}

	// TO grantee [, grantee ...]
	if p.cur.Type == kwTO {
		p.advance() // consume TO
		stmt.Grantees = p.parseIdentList()
	}

	// WITH GRANT OPTION / WITH ADMIN OPTION
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		if p.cur.Type == kwGRANT {
			p.advance() // consume GRANT
			if p.cur.Type == kwOPTION {
				p.advance() // consume OPTION
			}
			stmt.WithGrant = true
		} else if p.isIdentLike() && p.cur.Str == "ADMIN" {
			p.advance() // consume ADMIN
			if p.cur.Type == kwOPTION {
				p.advance() // consume OPTION
			}
			stmt.WithAdmin = true
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseRevokeStmt parses a REVOKE statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/REVOKE.html
//
//	REVOKE { priv [, priv ...] | ALL [PRIVILEGES] }
//	  ON [object_type] object
//	  FROM grantee [, grantee ...]
//
//	REVOKE role [, role ...] FROM user [, user ...]
func (p *Parser) parseRevokeStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume REVOKE

	stmt := &nodes.RevokeStmt{
		Privileges: &nodes.List{},
		Grantees:   &nodes.List{},
		Loc:        nodes.Loc{Start: start},
	}

	// Parse privileges or ALL [PRIVILEGES].
	stmt.AllPriv, stmt.Privileges = p.parsePrivilegeList()

	// ON clause (optional).
	if p.cur.Type == kwON {
		p.advance() // consume ON

		// Optional object type keyword.
		stmt.OnType = p.parseOptionalObjectType()

		// Object name.
		stmt.OnObject = p.parseObjectName()
	}

	// FROM grantee [, grantee ...]
	if p.cur.Type == kwFROM {
		p.advance() // consume FROM
		stmt.Grantees = p.parseIdentList()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parsePrivilegeList parses a comma-separated list of privilege names or ALL [PRIVILEGES].
// Returns (allPriv bool, privileges *List).
func (p *Parser) parsePrivilegeList() (bool, *nodes.List) {
	privs := &nodes.List{}

	// ALL [PRIVILEGES]
	if p.cur.Type == kwALL {
		p.advance() // consume ALL
		if p.cur.Type == kwPRIVILEGES {
			p.advance() // consume PRIVILEGES
		}
		return true, privs
	}

	// Comma-separated list of privilege names.
	// Privileges can be multi-word (e.g., "CREATE TABLE", "ALTER ANY TABLE").
	// We parse each privilege as a sequence of identifiers up to a comma, ON, TO, or FROM.
	for {
		priv := p.parsePrivilegeName()
		if priv == "" {
			break
		}
		privs.Items = append(privs.Items, &nodes.String{Str: priv})

		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}

	return false, privs
}

// parsePrivilegeName parses a single privilege name, which may be multi-word
// (e.g., "SELECT", "CREATE TABLE", "ALTER ANY TABLE").
// Stops at comma, ON, TO, FROM, semicolon, or EOF.
func (p *Parser) parsePrivilegeName() string {
	if !p.isIdentLike() {
		return ""
	}

	// Don't consume ON, TO, FROM as privilege names.
	if p.cur.Type == kwON || p.cur.Type == kwTO || p.cur.Type == kwFROM {
		return ""
	}

	name := p.parseIdentifier()
	if name == "" {
		return ""
	}

	// Accumulate multi-word privilege names.
	for p.isIdentLike() &&
		p.cur.Type != kwON &&
		p.cur.Type != kwTO &&
		p.cur.Type != kwFROM &&
		p.cur.Type != ',' &&
		p.cur.Type != ';' {
		word := p.parseIdentifier()
		if word == "" {
			break
		}
		name += " " + word
	}

	return name
}

// parseOptionalObjectType checks if the current token is an object type keyword
// and returns the corresponding ObjectType. Returns OBJECT_TABLE as default.
func (p *Parser) parseOptionalObjectType() nodes.ObjectType {
	switch p.cur.Type {
	case kwTABLE:
		p.advance()
		return nodes.OBJECT_TABLE
	case kwVIEW:
		p.advance()
		return nodes.OBJECT_VIEW
	case kwSEQUENCE:
		p.advance()
		return nodes.OBJECT_SEQUENCE
	case kwPROCEDURE:
		p.advance()
		return nodes.OBJECT_PROCEDURE
	case kwFUNCTION:
		p.advance()
		return nodes.OBJECT_FUNCTION
	case kwPACKAGE:
		p.advance()
		return nodes.OBJECT_PACKAGE
	case kwTYPE:
		p.advance()
		return nodes.OBJECT_TYPE
	case kwINDEX:
		p.advance()
		return nodes.OBJECT_INDEX
	default:
		// No explicit object type — default to TABLE for object privileges.
		return nodes.OBJECT_TABLE
	}
}

// parseIdentList parses a comma-separated list of identifiers and returns
// a *List of *String nodes.
func (p *Parser) parseIdentList() *nodes.List {
	list := &nodes.List{}
	for {
		if !p.isIdentLike() {
			break
		}
		name := p.parseIdentifier()
		if name == "" {
			break
		}
		list.Items = append(list.Items, &nodes.String{Str: name})

		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	return list
}
