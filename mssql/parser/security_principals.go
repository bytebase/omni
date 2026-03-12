// Package parser - security_principals.go implements T-SQL security principal
// statement parsing: CREATE/ALTER/DROP USER, LOGIN, ROLE, APPLICATION ROLE,
// and ADD/DROP ROLE MEMBER.
package parser

import (
	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseSecurityUserStmt parses CREATE/ALTER/DROP USER.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-user-transact-sql
//
//	CREATE USER user_name [FOR|FROM LOGIN login_name] [WITH ...]
//	ALTER  USER user_name WITH NAME = new_name [, ...]
//	DROP   USER [IF EXISTS] user_name
func (p *Parser) parseSecurityUserStmt(action string) *nodes.SecurityStmt {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     action,
		ObjectType: "USER",
		Loc:        nodes.Loc{Start: loc},
	}

	// IF EXISTS (DROP only)
	if action == "DROP" {
		if p.cur.Type == kwIF {
			next := p.peekNext()
			if next.Type == kwEXISTS {
				p.advance() // IF
				p.advance() // EXISTS
			}
		}
	}

	// user name
	if name, ok := p.parseIdentifier(); ok {
		stmt.Name = name
	}

	// FOR/FROM LOGIN login_name  (CREATE)
	var opts []nodes.Node
	if action == "CREATE" {
		if p.cur.Type == kwFOR || p.cur.Type == kwFROM {
			p.advance() // consume FOR / FROM
			if p.isIdentLike() && (matchesKeywordCI(p.cur.Str, "LOGIN") || p.cur.Type == kwLOGIN) {
				p.advance() // consume LOGIN
			}
			if name, ok := p.parseIdentifier(); ok {
				opts = append(opts, &nodes.String{Str: "LOGIN=" + name})
			}
		}
	}

	// WITH options
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		for {
			if !p.isIdentLike() {
				break
			}
			key := p.cur.Str
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
				if val, ok := p.parseIdentifier(); ok {
					opts = append(opts, &nodes.String{Str: key + "=" + val})
				}
			} else {
				opts = append(opts, &nodes.String{Str: key})
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
	}

	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseSecurityLoginStmt parses CREATE/ALTER/DROP LOGIN.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-login-transact-sql
//
//	CREATE LOGIN login_name { WITH PASSWORD = 'password' [, ...] | FROM WINDOWS [WITH ...] }
//	ALTER  LOGIN login_name { ENABLE | DISABLE | WITH ... }
//	DROP   LOGIN login_name
func (p *Parser) parseSecurityLoginStmt(action string) *nodes.SecurityStmt {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     action,
		ObjectType: "LOGIN",
		Loc:        nodes.Loc{Start: loc},
	}

	// login name
	if name, ok := p.parseIdentifier(); ok {
		stmt.Name = name
	}

	var opts []nodes.Node

	// ENABLE / DISABLE (ALTER LOGIN)
	if action == "ALTER" {
		if matchesKeywordCI(p.cur.Str, "ENABLE") && p.isIdentLike() {
			opts = append(opts, &nodes.String{Str: "ENABLE"})
			p.advance()
		} else if matchesKeywordCI(p.cur.Str, "DISABLE") && p.isIdentLike() {
			opts = append(opts, &nodes.String{Str: "DISABLE"})
			p.advance()
		}
	}

	// WITH / FROM options
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		for {
			if !p.isIdentLike() {
				break
			}
			key := p.cur.Str
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
				// Password can be a string literal
				if p.cur.Type == tokSCONST {
					val := p.cur.Str
					p.advance()
					opts = append(opts, &nodes.String{Str: key + "='" + val + "'"})
				} else if val, ok := p.parseIdentifier(); ok {
					opts = append(opts, &nodes.String{Str: key + "=" + val})
				}
			} else {
				opts = append(opts, &nodes.String{Str: key})
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
	} else if p.cur.Type == kwFROM {
		p.advance() // consume FROM
		if p.isIdentLike() {
			src := p.cur.Str
			p.advance()
			opts = append(opts, &nodes.String{Str: "FROM=" + src})
		}
		if p.cur.Type == kwWITH {
			p.advance()
			for {
				if !p.isIdentLike() {
					break
				}
				key := p.cur.Str
				p.advance()
				if p.cur.Type == '=' {
					p.advance()
					if val, ok := p.parseIdentifier(); ok {
						opts = append(opts, &nodes.String{Str: key + "=" + val})
					}
				} else {
					opts = append(opts, &nodes.String{Str: key})
				}
				if _, ok := p.match(','); !ok {
					break
				}
			}
		}
	}

	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseSecurityRoleStmt parses CREATE/ALTER/DROP ROLE and ADD/DROP ROLE MEMBER.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-role-transact-sql
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-role-transact-sql
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-role-transact-sql
//
//	CREATE ROLE role_name [AUTHORIZATION owner_name]
//	ALTER  ROLE role_name { ADD MEMBER member_name | DROP MEMBER member_name | WITH NAME = new_name }
//	DROP   ROLE [IF EXISTS] role_name
func (p *Parser) parseSecurityRoleStmt(action string) *nodes.SecurityStmt {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     action,
		ObjectType: "ROLE",
		Loc:        nodes.Loc{Start: loc},
	}

	// IF EXISTS (DROP only)
	if action == "DROP" {
		if p.cur.Type == kwIF {
			next := p.peekNext()
			if next.Type == kwEXISTS {
				p.advance() // IF
				p.advance() // EXISTS
			}
		}
	}

	// role name
	if name, ok := p.parseIdentifier(); ok {
		stmt.Name = name
	}

	var opts []nodes.Node

	if action == "CREATE" {
		// AUTHORIZATION owner_name
		if matchesKeywordCI(p.cur.Str, "AUTHORIZATION") && p.isIdentLike() {
			p.advance()
			if owner, ok := p.parseIdentifier(); ok {
				opts = append(opts, &nodes.String{Str: "AUTHORIZATION=" + owner})
			}
		}
	} else if action == "ALTER" {
		// ADD MEMBER | DROP MEMBER | WITH NAME = new_name
		if p.cur.Type == kwADD {
			p.advance() // consume ADD
			if matchesKeywordCI(p.cur.Str, "MEMBER") && p.isIdentLike() {
				p.advance() // consume MEMBER
			}
			if member, ok := p.parseIdentifier(); ok {
				opts = append(opts, &nodes.String{Str: "ADD MEMBER=" + member})
			}
		} else if p.cur.Type == kwDROP {
			p.advance() // consume DROP
			if matchesKeywordCI(p.cur.Str, "MEMBER") && p.isIdentLike() {
				p.advance() // consume MEMBER
			}
			if member, ok := p.parseIdentifier(); ok {
				opts = append(opts, &nodes.String{Str: "DROP MEMBER=" + member})
			}
		} else if p.cur.Type == kwWITH {
			p.advance() // consume WITH
			for {
				if !p.isIdentLike() {
					break
				}
				key := p.cur.Str
				p.advance()
				if p.cur.Type == '=' {
					p.advance()
					if val, ok := p.parseIdentifier(); ok {
						opts = append(opts, &nodes.String{Str: key + "=" + val})
					}
				} else {
					opts = append(opts, &nodes.String{Str: key})
				}
				if _, ok := p.match(','); !ok {
					break
				}
			}
		}
	}

	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseSecurityApplicationRoleStmt parses CREATE/ALTER/DROP APPLICATION ROLE.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-application-role-transact-sql
//
//	CREATE APPLICATION ROLE role_name WITH PASSWORD = 'password' [, DEFAULT_SCHEMA = schema_name]
//	ALTER  APPLICATION ROLE role_name WITH NAME = new_name | PASSWORD = '...' | DEFAULT_SCHEMA = ...
//	DROP   APPLICATION ROLE role_name
func (p *Parser) parseSecurityApplicationRoleStmt(action string) *nodes.SecurityStmt {
	loc := p.pos()
	// APPLICATION keyword already consumed by caller
	// Consume ROLE
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "ROLE") {
		p.advance()
	} else if p.cur.Type == kwROLE {
		p.advance()
	}

	stmt := &nodes.SecurityStmt{
		Action:     action,
		ObjectType: "APPLICATION ROLE",
		Loc:        nodes.Loc{Start: loc},
	}

	// role name
	if name, ok := p.parseIdentifier(); ok {
		stmt.Name = name
	}

	var opts []nodes.Node

	// WITH options
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		for {
			if !p.isIdentLike() {
				break
			}
			key := p.cur.Str
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
				if p.cur.Type == tokSCONST {
					val := p.cur.Str
					p.advance()
					opts = append(opts, &nodes.String{Str: key + "='" + val + "'"})
				} else if val, ok := p.parseIdentifier(); ok {
					opts = append(opts, &nodes.String{Str: key + "=" + val})
				}
			} else {
				opts = append(opts, &nodes.String{Str: key})
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
	}

	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// matchesKeywordCI returns true if s case-insensitively equals keyword.
// Helper used to check string tokens against context-sensitive keywords.
func matchesKeywordCI(s, keyword string) bool {
	if len(s) != len(keyword) {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		k := keyword[i]
		if c >= 'a' && c <= 'z' {
			c -= 32
		}
		if k >= 'a' && k <= 'z' {
			k -= 32
		}
		if c != k {
			return false
		}
	}
	return true
}
