// Package parser - grant.go implements T-SQL GRANT/REVOKE/DENY statement parsing.
package parser

import (
	nodes "github.com/bytebase/omni/tsql/ast"
)

// parseGrantStmt parses a GRANT statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/grant-transact-sql
//
//	GRANT permission [, ...] [ON [object_type ::] object] TO principal [, ...] [WITH GRANT OPTION]
func (p *Parser) parseGrantStmt() *nodes.GrantStmt {
	loc := p.pos()
	p.advance() // consume GRANT

	stmt := &nodes.GrantStmt{
		StmtType: nodes.GrantTypeGrant,
		Loc:      nodes.Loc{Start: loc},
	}

	// Privileges
	stmt.Privileges = p.parsePrivilegeList()

	// ON [type ::] object
	if _, ok := p.match(kwON); ok {
		stmt.OnName = p.parseTableRef()
	}

	// TO principals
	if _, ok := p.match(kwTO); ok {
		stmt.Principals = p.parsePrincipalList()
	}

	// WITH GRANT OPTION
	if p.cur.Type == kwWITH {
		next := p.peekNext()
		if next.Type == kwGRANT {
			p.advance() // WITH
			p.advance() // GRANT
			// OPTION
			if p.cur.Type == kwOPTION {
				p.advance()
			}
			stmt.WithGrant = true
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseRevokeStmt parses a REVOKE statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/revoke-transact-sql
//
//	REVOKE permission [, ...] [ON object] FROM principal [, ...] [CASCADE]
func (p *Parser) parseRevokeStmt() *nodes.GrantStmt {
	loc := p.pos()
	p.advance() // consume REVOKE

	stmt := &nodes.GrantStmt{
		StmtType: nodes.GrantTypeRevoke,
		Loc:      nodes.Loc{Start: loc},
	}

	// Privileges
	stmt.Privileges = p.parsePrivilegeList()

	// ON object
	if _, ok := p.match(kwON); ok {
		stmt.OnName = p.parseTableRef()
	}

	// FROM principals
	if _, ok := p.match(kwFROM); ok {
		stmt.Principals = p.parsePrincipalList()
	}

	// CASCADE
	if _, ok := p.match(kwCASCADE); ok {
		stmt.CascadeOpt = true
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseDenyStmt parses a DENY statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/deny-transact-sql
//
//	DENY permission [, ...] [ON object] TO principal [, ...] [CASCADE]
func (p *Parser) parseDenyStmt() *nodes.GrantStmt {
	loc := p.pos()
	p.advance() // consume DENY

	stmt := &nodes.GrantStmt{
		StmtType: nodes.GrantTypeDeny,
		Loc:      nodes.Loc{Start: loc},
	}

	// Privileges
	stmt.Privileges = p.parsePrivilegeList()

	// ON object
	if _, ok := p.match(kwON); ok {
		stmt.OnName = p.parseTableRef()
	}

	// TO principals
	if _, ok := p.match(kwTO); ok {
		stmt.Principals = p.parsePrincipalList()
	}

	// CASCADE
	if _, ok := p.match(kwCASCADE); ok {
		stmt.CascadeOpt = true
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parsePrivilegeList parses a comma-separated list of privileges.
func (p *Parser) parsePrivilegeList() *nodes.List {
	var items []nodes.Node
	for {
		if !p.isIdentLike() && p.cur.Type != kwSELECT && p.cur.Type != kwINSERT &&
			p.cur.Type != kwUPDATE && p.cur.Type != kwDELETE && p.cur.Type != kwEXEC &&
			p.cur.Type != kwEXECUTE && p.cur.Type != kwREFERENCES && p.cur.Type != kwALL {
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
