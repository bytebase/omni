// Package parser - grant.go implements T-SQL GRANT/REVOKE/DENY statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseGrantStmt parses a GRANT statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/grant-transact-sql
//
//	GRANT { ALL [ PRIVILEGES ] }
//	      | permission [ ( column [ , ...n ] ) ] [ , ...n ]
//	      [ ON [ class :: ] securable ] TO principal [ , ...n ]
//	      [ WITH GRANT OPTION ] [ AS principal ]
//
//	<class> ::=
//	{
//	    LOGIN | DATABASE | OBJECT | ROLE | SCHEMA | USER
//	  | APPLICATION ROLE | ASSEMBLY | ASYMMETRIC KEY | CERTIFICATE
//	  | CONTRACT | ENDPOINT | FULLTEXT CATALOG | FULLTEXT STOPLIST
//	  | MESSAGE TYPE | REMOTE SERVICE BINDING | ROUTE | SEARCH PROPERTY LIST
//	  | SERVER | SERVER ROLE | SERVICE | SYMMETRIC KEY | TYPE
//	  | XML SCHEMA COLLECTION | AVAILABILITY GROUP
//	}
func (p *Parser) parseGrantStmt() (*nodes.GrantStmt, error) {
	loc := p.pos()
	p.advance() // consume GRANT

	// Completion: after GRANT → permission keywords
	if p.collectMode() {
		p.addTokenCandidate(kwSELECT)
		p.addTokenCandidate(kwINSERT)
		p.addTokenCandidate(kwUPDATE)
		p.addTokenCandidate(kwDELETE)
		p.addTokenCandidate(kwEXECUTE)
		p.addTokenCandidate(kwALTER)
		p.addTokenCandidate(kwREFERENCES)
		p.addTokenCandidate(kwALL)
		p.addRuleCandidate("permission")
		return nil, errCollecting
	}

	stmt := &nodes.GrantStmt{
		StmtType: nodes.GrantTypeGrant,
		Loc:      nodes.Loc{Start: loc, End: -1},
	}

	// Privileges
	stmt.Privileges = p.parsePrivilegeList()

	// ON [class ::] securable
	if _, ok := p.match(kwON); ok {
		// Completion: after GRANT ... ON → table_ref
		if p.collectMode() {
			p.addRuleCandidate("table_ref")
			return nil, errCollecting
		}
		p.parseGrantOnClause(stmt)
	}

	// TO principals
	if _, ok := p.match(kwTO); ok {
		// Completion: after GRANT ... ON t TO → user_ref / role_ref
		if p.collectMode() {
			p.addRuleCandidate("user_ref")
			p.addRuleCandidate("role_ref")
			return nil, errCollecting
		}
		stmt.Principals = p.parsePrincipalList()
	}

	// WITH GRANT OPTION
	if p.cur.Type == kwWITH {
		next := p.peekNext()
		if next.Type == kwGRANT {
			p.advance() // WITH
			p.advance() // GRANT
			if p.cur.Type == kwOPTION {
				p.advance()
			}
			stmt.WithGrant = true
		}
	}

	// AS principal
	if p.cur.Type == kwAS {
		p.advance() // consume AS
		if p.isIdentLike() {
			stmt.AsPrincipal = p.cur.Str
			p.advance()
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseRevokeStmt parses a REVOKE statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/revoke-transact-sql
//
//	REVOKE [ GRANT OPTION FOR ]
//	      { ALL [ PRIVILEGES ] }
//	      | permission [ ( column [ , ...n ] ) ] [ , ...n ]
//	      [ ON [ class :: ] securable ]
//	      { TO | FROM } principal [ , ...n ]
//	      [ CASCADE ] [ AS principal ]
func (p *Parser) parseRevokeStmt() (*nodes.GrantStmt, error) {
	loc := p.pos()
	p.advance() // consume REVOKE

	// Completion: after REVOKE → permission keywords
	if p.collectMode() {
		p.addTokenCandidate(kwSELECT)
		p.addTokenCandidate(kwINSERT)
		p.addTokenCandidate(kwUPDATE)
		p.addTokenCandidate(kwDELETE)
		p.addTokenCandidate(kwEXECUTE)
		p.addTokenCandidate(kwALTER)
		p.addTokenCandidate(kwREFERENCES)
		p.addTokenCandidate(kwALL)
		p.addRuleCandidate("permission")
		return nil, errCollecting
	}

	stmt := &nodes.GrantStmt{
		StmtType: nodes.GrantTypeRevoke,
		Loc:      nodes.Loc{Start: loc, End: -1},
	}

	// GRANT OPTION FOR
	if p.cur.Type == kwGRANT {
		next := p.peekNext()
		if next.Type == kwOPTION {
			p.advance() // consume GRANT
			p.advance() // consume OPTION
			if p.cur.Type == kwFOR {
				p.advance() // consume FOR
			}
			stmt.GrantOptionFor = true
		}
	}

	// Privileges
	stmt.Privileges = p.parsePrivilegeList()

	// ON [class ::] securable
	if _, ok := p.match(kwON); ok {
		p.parseGrantOnClause(stmt)
	}

	// FROM or TO principals
	if _, ok := p.match(kwFROM); ok {
		stmt.Principals = p.parsePrincipalList()
	} else if _, ok := p.match(kwTO); ok {
		stmt.Principals = p.parsePrincipalList()
	}

	// CASCADE
	if _, ok := p.match(kwCASCADE); ok {
		stmt.CascadeOpt = true
	}

	// AS principal
	if p.cur.Type == kwAS {
		p.advance() // consume AS
		if p.isIdentLike() {
			stmt.AsPrincipal = p.cur.Str
			p.advance()
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseDenyStmt parses a DENY statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/deny-transact-sql
//
//	DENY { ALL [ PRIVILEGES ] }
//	     | permission [ ( column [ , ...n ] ) ] [ , ...n ]
//	     [ ON [ class :: ] securable ] TO principal [ , ...n ]
//	     [ CASCADE ] [ AS principal ]
func (p *Parser) parseDenyStmt() (*nodes.GrantStmt, error) {
	loc := p.pos()
	p.advance() // consume DENY

	// Completion: after DENY → permission keywords
	if p.collectMode() {
		p.addTokenCandidate(kwSELECT)
		p.addTokenCandidate(kwINSERT)
		p.addTokenCandidate(kwUPDATE)
		p.addTokenCandidate(kwDELETE)
		p.addTokenCandidate(kwEXECUTE)
		p.addTokenCandidate(kwALTER)
		p.addTokenCandidate(kwREFERENCES)
		p.addTokenCandidate(kwALL)
		p.addRuleCandidate("permission")
		return nil, errCollecting
	}

	stmt := &nodes.GrantStmt{
		StmtType: nodes.GrantTypeDeny,
		Loc:      nodes.Loc{Start: loc, End: -1},
	}

	// Privileges
	stmt.Privileges = p.parsePrivilegeList()

	// ON [class ::] securable
	if _, ok := p.match(kwON); ok {
		p.parseGrantOnClause(stmt)
	}

	// TO principals
	if _, ok := p.match(kwTO); ok {
		stmt.Principals = p.parsePrincipalList()
	}

	// CASCADE
	if _, ok := p.match(kwCASCADE); ok {
		stmt.CascadeOpt = true
	}

	// AS principal
	if p.cur.Type == kwAS {
		p.advance() // consume AS
		if p.isIdentLike() {
			stmt.AsPrincipal = p.cur.Str
			p.advance()
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseGrantOnClause parses ON [ class :: ] securable.
// It detects the class :: pattern by checking for tokCOLONCOLON after the first identifier.
// Securable classes can be multi-word like "XML SCHEMA COLLECTION", "APPLICATION ROLE", etc.
func (p *Parser) parseGrantOnClause(stmt *nodes.GrantStmt) {
	// Try to detect class :: securable pattern.
	// Look ahead: if cur is an identifier and the next or next-next token is ::, it's a class.
	// Securable classes include Core keywords (SCHEMA, DATABASE, etc.)
	// so we must accept any keyword token, not just identifiers.
	if p.isKeywordOrIdent() {
		// Check for single-word class followed by ::
		next := p.peekNext()
		if next.Type == tokCOLONCOLON {
			// Single-word class: SCHEMA, OBJECT, DATABASE, LOGIN, USER, ROLE, etc.
			stmt.OnType = strings.ToUpper(p.cur.Str)
			p.advance() // consume class
			p.advance() // consume ::
			stmt.OnName , _ = p.parseTableRef()
			return
		}
		// Check for multi-word class: peek further
		// Multi-word classes: APPLICATION ROLE, ASYMMETRIC KEY, FULLTEXT CATALOG,
		// FULLTEXT STOPLIST, MESSAGE TYPE, REMOTE SERVICE BINDING, SEARCH PROPERTY LIST,
		// SERVER ROLE, SYMMETRIC KEY, XML SCHEMA COLLECTION, AVAILABILITY GROUP
		if p.isSecurableClassStart() {
			className := p.tryParseSecurableClass()
			if className != "" {
				stmt.OnType = className
				stmt.OnName , _ = p.parseTableRef()
				return
			}
		}
	}

	// No class :: prefix, just parse the securable as a table ref
	stmt.OnName , _ = p.parseTableRef()
}

// isSecurableClassStart checks if the current token could be the start of a multi-word securable class.
func (p *Parser) isSecurableClassStart() bool {
	if !p.isIdentLike() {
		return false
	}
	kw := strings.ToUpper(p.cur.Str)
	switch kw {
	case "APPLICATION", "ASYMMETRIC", "FULLTEXT", "MESSAGE", "REMOTE",
		"SEARCH", "SERVER", "SYMMETRIC", "XML", "AVAILABILITY":
		return true
	}
	return false
}

// tryParseSecurableClass tries to parse a multi-word securable class followed by ::.
// Returns the class name if successful, or empty string if not (without consuming tokens).
// Since we only have 1-token lookahead, we consume tokens greedily and use :: as confirmation.
func (p *Parser) tryParseSecurableClass() string {
	// Save position info for potential rollback - but we don't have rollback,
	// so we'll consume greedily and check for :: after building the class name.
	var parts []string
	parts = append(parts, strings.ToUpper(p.cur.Str))
	p.advance()

	// Continue consuming identifier-like tokens until we see :: or non-identifier
	for p.isKeywordOrIdent() {
		next := p.peekNext()
		if next.Type == tokCOLONCOLON {
			// This identifier is the last word of the class
			parts = append(parts, strings.ToUpper(p.cur.Str))
			p.advance() // consume last word
			p.advance() // consume ::
			return strings.Join(parts, " ")
		}
		parts = append(parts, strings.ToUpper(p.cur.Str))
		p.advance()
	}

	// If we hit :: directly
	if p.cur.Type == tokCOLONCOLON {
		p.advance() // consume ::
		return strings.Join(parts, " ")
	}

	// Not a class :: pattern - the tokens are already consumed.
	// This shouldn't happen if isSecurableClassStart was checked properly.
	return ""
}

// parsePrivilegeList parses a comma-separated list of privileges.
// Privileges can include column lists: permission (col1, col2, ...)
func (p *Parser) parsePrivilegeList() *nodes.List {
	var items []nodes.Node
	for {
		if !p.isIdentLike() && p.cur.Type != kwSELECT && p.cur.Type != kwINSERT &&
			p.cur.Type != kwUPDATE && p.cur.Type != kwDELETE && p.cur.Type != kwEXEC &&
			p.cur.Type != kwEXECUTE && p.cur.Type != kwREFERENCES && p.cur.Type != kwALL &&
			p.cur.Type != kwCREATE && p.cur.Type != kwALTER {
			break
		}
		name := strings.ToUpper(p.cur.Str)
		p.advance()

		// Handle multi-word privileges: ALTER ANY xxx, CREATE xxx, VIEW xxx, etc.
		// Also ALL PRIVILEGES
		if name == "ALL" && p.cur.Type == kwPRIVILEGES {
			p.advance()
			name = "ALL PRIVILEGES"
		} else if (name == "ALTER" || name == "CREATE" || name == "VIEW" || name == "CONTROL" ||
			name == "TAKE" || name == "IMPERSONATE" || name == "BACKUP" || name == "CONNECT" ||
			name == "ADMINISTER" || name == "AUTHENTICATE" || name == "SHOWPLAN" ||
			name == "SUBSCRIBE" || name == "RECEIVE" || name == "SEND" || name == "UNMASK") &&
			p.isKeywordOrIdent() && p.cur.Type != kwON && p.cur.Type != kwTO && p.cur.Type != kwFROM {
			// Multi-word: ALTER ANY DATABASE, CREATE TABLE, VIEW DEFINITION, etc.
			for p.isKeywordOrIdent() && p.cur.Type != kwON && p.cur.Type != kwTO && p.cur.Type != kwFROM &&
				p.cur.Type != ',' && p.cur.Type != ';' && p.cur.Type != tokEOF {
				name += " " + strings.ToUpper(p.cur.Str)
				p.advance()
			}
		}

		// Column list: permission (col1, col2, ...)
		if p.cur.Type == '(' {
			p.advance() // consume (
			var cols []string
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				if p.isIdentLike() {
					cols = append(cols, p.cur.Str)
					p.advance()
				}
				if p.cur.Type == ',' {
					p.advance()
				} else {
					break
				}
			}
			p.match(')') // consume )
			if len(cols) > 0 {
				name += "(" + strings.Join(cols, ", ") + ")"
			}
		}

		items = append(items, &nodes.String{Str: name})
		if _, ok := p.match(','); !ok {
			break
		}
	}
	return &nodes.List{Items: items}
}

// parsePrincipalList parses a comma-separated list of principals.
func (p *Parser) parsePrincipalList() *nodes.List {
	var items []nodes.Node
	for {
		if !p.isIdentLike() && p.cur.Type != kwPUBLIC {
			break
		}
		name := p.cur.Str
		p.advance()
		items = append(items, &nodes.String{Str: name})
		if _, ok := p.match(','); !ok {
			break
		}
	}
	return &nodes.List{Items: items}
}
