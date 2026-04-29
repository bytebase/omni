package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseGrantStmt parses a GRANT statement.
//
// BNF: oracle/parser/bnf/GRANT.bnf
//
//	grant::=
//	    grant_system_privileges
//	  | grant_schema_privileges
//	  | grant_object_privileges
//	  | grant_roles_to_programs
//
//	grant_system_privileges::=
//	    GRANT { system_privilege [, system_privilege ]...
//	          | ALL PRIVILEGES
//	          | role [, role ]...
//	          }
//	        TO { grantee_clause [, grantee_clause ]... }
//	        [ IDENTIFIED BY password [, password ]... ]
//	        [ WITH { ADMIN | DELEGATE } OPTION ]
//	        [ CONTAINER = { CURRENT | ALL } ]
//
//	grant_schema_privileges::=
//	    GRANT { schema_privilege [, schema_privilege ]...
//	          | ALL [ PRIVILEGES ]
//	          }
//	        ON SCHEMA schema
//	        TO { grantee_clause [, grantee_clause ]... }
//	        [ WITH ADMIN OPTION ]
//	        [ CONTAINER = { CURRENT | ALL } ]
//
//	grant_object_privileges::=
//	    GRANT { object_privilege [ ( column [, column ]... ) ]
//	              [, object_privilege [ ( column [, column ]... ) ] ]...
//	          | ALL [ PRIVILEGES ]
//	          }
//	        on_object_clause
//	        TO { grantee_clause [, grantee_clause ]... }
//	        [ WITH { GRANT | HIERARCHY } OPTION ]
//	        [ CONTAINER = { CURRENT | ALL } ]
//
//	on_object_clause::=
//	    ON [ schema. ] object
//	  | ON USER user [, user ]...
//	  | ON DIRECTORY directory_name
//	  | ON EDITION edition_name
//	  | ON MINING MODEL [ schema. ] mining_model_name
//	  | ON JAVA { SOURCE | RESOURCE } [ schema. ] object
//	  | ON SQL TRANSLATION PROFILE [ schema. ] profile
//
//	grant_roles_to_programs::=
//	    GRANT role [, role ]...
//	        TO program_unit [, program_unit ]...
//	        [ WITH { ADMIN | DELEGATE } OPTION ]
//	        [ CONTAINER = { CURRENT | ALL } ]
//
//	grantee_clause::=
//	    { user | role | PUBLIC }
func (p *Parser) parseGrantStmt() (nodes.StmtNode, error) {
	start := p.pos()
	p.advance() // consume GRANT

	stmt := &nodes.GrantStmt{
		Privileges: &nodes.List{},
		Grantees:   &nodes.List{},
		Loc:        nodes.Loc{Start: start},
	}
	var parseErr763 error

	// Parse privileges or ALL [PRIVILEGES].
	stmt.AllPriv, stmt.Privileges, parseErr763 = p.parsePrivilegeList()
	if parseErr763 !=

		// ON clause (optional — absent means role grant or system privilege).
		nil {
		return nil, parseErr763
	}

	if p.cur.Type == kwON {
		p.advance()
		var // consume ON
		parseErr764 error

		// Optional object type keyword.
		stmt.OnType, parseErr764 = p.parseOptionalObjectType()
		if parseErr764 !=

			// Object name.
			nil {
			return nil, parseErr764
		}
		var parseErr765 error

		stmt.OnObject, parseErr765 = p.parseObjectName()
		if parseErr765 !=

			// TO grantee [, grantee ...]
			nil {
			return nil, parseErr765
		}
	}

	if p.cur.Type == kwTO {
		p.advance()
		var // consume TO
		parseErr766 error
		stmt.Grantees, parseErr766 = p.parseIdentList()
		if parseErr766 !=

			// WITH GRANT OPTION / WITH ADMIN OPTION
			nil {
			return nil, parseErr766
		}
	}

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

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseRevokeStmt parses a REVOKE statement.
//
// BNF: oracle/parser/bnf/REVOKE.bnf
//
//	REVOKE
//	  { revoke_system_privileges
//	  | revoke_schema_privileges
//	  | revoke_object_privileges
//	  | revoke_roles_from_programs
//	  }
//	  [ CONTAINER = { CURRENT | ALL } ] ;
//
//	revoke_system_privileges:
//	    REVOKE { system_privilege [, system_privilege ]...
//	           | role [, role ]...
//	           | ALL PRIVILEGES }
//	    FROM revokee_clause
//
//	revoke_schema_privileges:
//	    REVOKE { schema_privilege [, schema_privilege ]...
//	           | ALL [ PRIVILEGES ] }
//	    ON SCHEMA schema_name
//	    FROM revokee_clause
//
//	revoke_object_privileges:
//	    REVOKE { object_privilege [, object_privilege ]...
//	           | ALL [ PRIVILEGES ] }
//	    on_object_clause
//	    FROM revokee_clause
//	    [ CASCADE CONSTRAINTS ]
//	    [ FORCE ]
//
//	on_object_clause:
//	    ON { [ schema. ] object
//	       | USER user [, user ]...
//	       | DIRECTORY directory_name
//	       | EDITION edition_name
//	       | MINING MODEL [ schema. ] mining_model_name
//	       | JAVA { SOURCE | RESOURCE } [ schema. ] object
//	       | SQL TRANSLATION PROFILE [ schema. ] profile_name
//	       }
//
//	revokee_clause:
//	    { user | role | PUBLIC } [, { user | role | PUBLIC } ]...
func (p *Parser) parseRevokeStmt() (nodes.StmtNode, error) {
	start := p.pos()
	p.advance() // consume REVOKE

	stmt := &nodes.RevokeStmt{
		Privileges: &nodes.List{},
		Grantees:   &nodes.List{},
		Loc:        nodes.Loc{Start: start},
	}
	var parseErr767 error

	// Parse privileges or ALL [PRIVILEGES].
	stmt.AllPriv, stmt.Privileges, parseErr767 = p.parsePrivilegeList()
	if parseErr767 !=

		// ON clause (optional).
		nil {
		return nil, parseErr767
	}

	if p.cur.Type == kwON {
		p.advance()
		var // consume ON
		parseErr768 error

		// Optional object type keyword.
		stmt.OnType, parseErr768 = p.parseOptionalObjectType()
		if parseErr768 !=

			// Object name.
			nil {
			return nil, parseErr768
		}
		var parseErr769 error

		stmt.OnObject, parseErr769 = p.parseObjectName()
		if parseErr769 !=

			// FROM grantee [, grantee ...]
			nil {
			return nil, parseErr769
		}
		if stmt.OnObject == nil || stmt.OnObject.Name == "" {
			return nil, p.syntaxErrorAtCur()
		}
	}

	if p.cur.Type == kwFROM {
		p.advance()
		var // consume FROM
		parseErr770 error
		stmt.Grantees, parseErr770 = p.parseIdentList()
		if parseErr770 != nil {
			return nil, parseErr770
		}
		if stmt.Grantees.Len() == 0 {
			return nil, p.syntaxErrorAtCur()
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parsePrivilegeList parses a comma-separated list of privilege names or ALL [PRIVILEGES].
// Returns (allPriv bool, privileges *List).
func (p *Parser) parsePrivilegeList() (bool, *nodes.List, error) {
	privs := &nodes.List{}

	// ALL [PRIVILEGES]
	if p.cur.Type == kwALL {
		p.advance() // consume ALL
		if p.cur.Type == kwPRIVILEGES {
			p.advance() // consume PRIVILEGES
		}
		return true, privs, nil
	}

	// Comma-separated list of privilege names.
	// Privileges can be multi-word (e.g., "CREATE TABLE", "ALTER ANY TABLE").
	// We parse each privilege as a sequence of identifiers up to a comma, ON, TO, or FROM.
	for {
		priv, parseErr771 := p.parsePrivilegeName()
		if parseErr771 != nil {
			return false, nil, parseErr771
		}
		if priv == "" {
			break
		}
		privs.Items = append(privs.Items, &nodes.String{Str: priv})

		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}

	return false, privs, nil
}

// parsePrivilegeName parses a single privilege name, which may be multi-word
// (e.g., "SELECT", "CREATE TABLE", "ALTER ANY TABLE").
// Stops at comma, ON, TO, FROM, semicolon, or EOF.
func (p *Parser) parsePrivilegeName() (string, error) {
	if !p.isIdentLike() {
		return "", nil
	}

	// Don't consume ON, TO, FROM as privilege names.
	if p.cur.Type == kwON || p.cur.Type == kwTO || p.cur.Type == kwFROM {
		return "", nil
	}

	name, parseErr772 := p.parseIdentifier()
	if parseErr772 != nil {
		return "", parseErr772

		// Accumulate multi-word privilege names.
	}
	if name == "" {
		return "", nil
	}

	for p.isIdentLike() &&
		p.cur.Type != kwON &&
		p.cur.Type != kwTO &&
		p.cur.Type != kwFROM &&
		p.cur.Type != ',' &&
		p.cur.Type != ';' {
		word, parseErr773 := p.parseIdentifier()
		if parseErr773 != nil {
			return "", parseErr773
		}
		if word == "" {
			break
		}
		name += " " + word
	}

	return name, nil
}

// parseOptionalObjectType checks if the current token is an object type keyword
// and returns the corresponding ObjectType. Returns OBJECT_TABLE as default.
func (p *Parser) parseOptionalObjectType() (nodes.ObjectType, error) {
	switch p.cur.Type {
	case kwTABLE:
		p.advance()
		return nodes.OBJECT_TABLE, nil
	case kwVIEW:
		p.advance()
		return nodes.OBJECT_VIEW, nil
	case kwSEQUENCE:
		p.advance()
		return nodes.OBJECT_SEQUENCE, nil
	case kwPROCEDURE:
		p.advance()
		return nodes.OBJECT_PROCEDURE, nil
	case kwFUNCTION:
		p.advance()
		return nodes.OBJECT_FUNCTION, nil
	case kwPACKAGE:
		p.advance()
		return nodes.OBJECT_PACKAGE, nil
	case kwTYPE:
		p.advance()
		return nodes.OBJECT_TYPE, nil
	case kwINDEX:
		p.advance()
		return nodes.OBJECT_INDEX, nil
	default:
		// No explicit object type — default to TABLE for object privileges.
		return nodes.OBJECT_TABLE, nil
	}
}

// parseIdentList parses a comma-separated list of identifiers and returns
// a *List of *String nodes.
func (p *Parser) parseIdentList() (*nodes.List, error) {
	list := &nodes.List{}
	for {
		if !p.isIdentLike() {
			break
		}
		name, parseErr774 := p.parseIdentifier()
		if parseErr774 != nil {
			return nil, parseErr774
		}
		if name == "" {
			break
		}
		list.Items = append(list.Items, &nodes.String{Str: name})

		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	return list, nil
}
