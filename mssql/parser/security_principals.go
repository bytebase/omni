// Package parser - security_principals.go implements T-SQL security principal
// statement parsing: CREATE/ALTER/DROP USER, LOGIN, ROLE, APPLICATION ROLE,
// and ADD/DROP ROLE MEMBER.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseSecurityUserStmt parses CREATE/ALTER/DROP USER.
//
// BNF: mssql/parser/bnf/create-user-transact-sql.bnf
//
//	CREATE USER user_name
//	    [ { FOR | FROM } LOGIN login_name ]
//	    [ WITH <option_list> [ ,... ] ]
//	CREATE USER user_name
//	    { FOR | FROM } CERTIFICATE cert_name
//	CREATE USER user_name
//	    { FOR | FROM } ASYMMETRIC KEY asym_key_name
//	CREATE USER user_name WITHOUT LOGIN [ WITH <limited_options_list> ]
//	CREATE USER user_name FROM EXTERNAL PROVIDER [ WITH OBJECT_ID = 'objectid' ]
//	ALTER USER user_name WITH <option_list> [ ,... ]
//	DROP USER [ IF EXISTS ] user_name
//
//	<option_list> ::=
//	    PASSWORD = 'password' [ OLD_PASSWORD = 'oldpassword' ]
//	  | DEFAULT_SCHEMA = { schema_name | NULL }
//	  | DEFAULT_LANGUAGE = { NONE | lcid | language_name | language_alias }
//	  | ALLOW_ENCRYPTED_VALUE_MODIFICATIONS = { ON | OFF }
//	  | NAME = new_user_name
//	  | LOGIN = new_login_name
//	  | SID = sid
func (p *Parser) parseSecurityUserStmt(action string) (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     action,
		ObjectType: "USER",
		Loc:        nodes.Loc{Start: loc, End: -1},
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

	// FOR/FROM or WITHOUT LOGIN (CREATE)
	var opts []nodes.Node
	if action == "CREATE" {
		if p.cur.Type == kwFOR || p.cur.Type == kwFROM {
			optLoc := p.pos()
			p.advance() // consume FOR / FROM
			if p.cur.Type == kwCERTIFICATE {
				p.advance() // consume CERTIFICATE
				if name, ok := p.parseIdentifier(); ok {
					opts = append(opts, &nodes.SecurityPrincipalOption{
						Name:  "CERTIFICATE",
						Value: name,
						Loc:   nodes.Loc{Start: optLoc, End: p.prevEnd()},
					})
				}
			} else if p.cur.Type == kwASYMMETRIC {
				p.advance() // consume ASYMMETRIC
				if p.cur.Type == kwKEY {
					p.advance() // consume KEY
				}
				if name, ok := p.parseIdentifier(); ok {
					opts = append(opts, &nodes.SecurityPrincipalOption{
						Name:  "ASYMMETRIC KEY",
						Value: name,
						Loc:   nodes.Loc{Start: optLoc, End: p.prevEnd()},
					})
				}
			} else if p.cur.Type == kwEXTERNAL {
				p.advance() // consume EXTERNAL
				if p.cur.Type == kwPROVIDER {
					p.advance() // consume PROVIDER
				}
				opts = append(opts, &nodes.SecurityPrincipalOption{
					Name: "EXTERNAL PROVIDER",
					Loc:  nodes.Loc{Start: optLoc, End: p.prevEnd()},
				})
			} else {
				// LOGIN (explicit or implicit)
				if p.cur.Type == kwLOGIN {
					p.advance() // consume LOGIN
				}
				if name, ok := p.parseIdentifier(); ok {
					opts = append(opts, &nodes.SecurityPrincipalOption{
						Name:  "LOGIN",
						Value: name,
						Loc:   nodes.Loc{Start: optLoc, End: p.prevEnd()},
					})
				}
			}
		} else if p.cur.Type == kwWITHOUT {
			optLoc := p.pos()
			p.advance() // consume WITHOUT
			if p.cur.Type == kwLOGIN {
				p.advance() // consume LOGIN
			}
			opts = append(opts, &nodes.SecurityPrincipalOption{
				Name: "WITHOUT LOGIN",
				Loc:  nodes.Loc{Start: optLoc, End: p.prevEnd()},
			})
		}
	}

	// WITH options
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		opts = append(opts, p.parseSecurityPrincipalWithOptions()...)
	}

	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseSecurityLoginStmt parses CREATE/ALTER/DROP LOGIN.
//
// BNF: mssql/parser/bnf/create-login-transact-sql.bnf
// BNF: mssql/parser/bnf/alter-login-transact-sql.bnf
// BNF: mssql/parser/bnf/drop-login-transact-sql.bnf
//
//	CREATE LOGIN login_name
//	    { WITH PASSWORD = { 'password' | hashed_password HASHED } [ MUST_CHANGE ]
//	        [ , DEFAULT_DATABASE = database ]
//	        [ , DEFAULT_LANGUAGE = language ]
//	        [ , CHECK_EXPIRATION = { ON | OFF } ]
//	        [ , CHECK_POLICY = { ON | OFF } ]
//	        [ , CREDENTIAL = credential_name ]
//	        [ , SID = sid ]
//	    | FROM WINDOWS [ WITH DEFAULT_DATABASE = database [, DEFAULT_LANGUAGE = language ] ]
//	    | FROM CERTIFICATE certname
//	    | FROM ASYMMETRIC KEY asym_key_name
//	    | FROM EXTERNAL PROVIDER [ WITH OBJECT_ID = 'objectid' ]
//	    }
//	ALTER LOGIN login_name
//	    { ENABLE | DISABLE
//	    | WITH PASSWORD = 'password' [ OLD_PASSWORD = 'old' | { MUST_CHANGE | UNLOCK } [...] ]
//	        [ , DEFAULT_DATABASE = database ] [ , DEFAULT_LANGUAGE = language ]
//	        [ , NAME = login_name ] [ , CHECK_POLICY = { ON | OFF } ]
//	        [ , CHECK_EXPIRATION = { ON | OFF } ] [ , CREDENTIAL = credential_name ]
//	        [ , NO CREDENTIAL ]
//	    | ADD CREDENTIAL credential_name
//	    | DROP CREDENTIAL credential_name
//	    }
//	DROP LOGIN login_name
func (p *Parser) parseSecurityLoginStmt(action string) (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     action,
		ObjectType: "LOGIN",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// login name
	if name, ok := p.parseIdentifier(); ok {
		stmt.Name = name
	}

	var opts []nodes.Node

	// ENABLE / DISABLE / ADD CREDENTIAL / DROP CREDENTIAL (ALTER LOGIN)
	if action == "ALTER" {
		if p.cur.Type == kwENABLE {
			optLoc := p.pos()
			p.advance()
			opts = append(opts, &nodes.SecurityPrincipalOption{
				Name: "ENABLE",
				Loc:  nodes.Loc{Start: optLoc, End: p.prevEnd()},
			})
		} else if p.cur.Type == kwDISABLE {
			optLoc := p.pos()
			p.advance()
			opts = append(opts, &nodes.SecurityPrincipalOption{
				Name: "DISABLE",
				Loc:  nodes.Loc{Start: optLoc, End: p.prevEnd()},
			})
		} else if p.cur.Type == kwADD {
			optLoc := p.pos()
			p.advance() // consume ADD
			if p.cur.Type == kwCREDENTIAL {
				p.advance() // consume CREDENTIAL
				if name, ok := p.parseIdentifier(); ok {
					opts = append(opts, &nodes.SecurityPrincipalOption{
						Name:  "ADD CREDENTIAL",
						Value: name,
						Loc:   nodes.Loc{Start: optLoc, End: p.prevEnd()},
					})
				}
			}
		} else if p.cur.Type == kwDROP {
			optLoc := p.pos()
			p.advance() // consume DROP
			if p.cur.Type == kwCREDENTIAL {
				p.advance() // consume CREDENTIAL
				if name, ok := p.parseIdentifier(); ok {
					opts = append(opts, &nodes.SecurityPrincipalOption{
						Name:  "DROP CREDENTIAL",
						Value: name,
						Loc:   nodes.Loc{Start: optLoc, End: p.prevEnd()},
					})
				}
			}
		}
	}

	// WITH options
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		opts = append(opts, p.parseSecurityPrincipalWithOptions()...)
	} else if p.cur.Type == kwFROM {
		optLoc := p.pos()
		p.advance() // consume FROM
		if p.isAnyKeywordIdent() {
			src := strings.ToUpper(p.cur.Str)
			switch src {
			case "WINDOWS":
				p.advance()
				opts = append(opts, &nodes.SecurityPrincipalOption{
					Name:  "FROM",
					Value: "WINDOWS",
					Loc:   nodes.Loc{Start: optLoc, End: p.prevEnd()},
				})
			case "CERTIFICATE":
				p.advance()
				certName := ""
				if id, ok := p.parseIdentifier(); ok {
					certName = id
				}
				opts = append(opts, &nodes.SecurityPrincipalOption{
					Name:  "FROM",
					Value: "CERTIFICATE " + certName,
					Loc:   nodes.Loc{Start: optLoc, End: p.prevEnd()},
				})
			case "ASYMMETRIC":
				p.advance() // consume ASYMMETRIC
				if p.cur.Type == kwKEY {
					p.advance() // consume KEY
				}
				keyName := ""
				if id, ok := p.parseIdentifier(); ok {
					keyName = id
				}
				opts = append(opts, &nodes.SecurityPrincipalOption{
					Name:  "FROM",
					Value: "ASYMMETRIC KEY " + keyName,
					Loc:   nodes.Loc{Start: optLoc, End: p.prevEnd()},
				})
			case "EXTERNAL":
				p.advance() // consume EXTERNAL
				if p.cur.Type == kwPROVIDER {
					p.advance()
				}
				opts = append(opts, &nodes.SecurityPrincipalOption{
					Name:  "FROM",
					Value: "EXTERNAL PROVIDER",
					Loc:   nodes.Loc{Start: optLoc, End: p.prevEnd()},
				})
			default:
				// unknown FROM source
				p.advance()
				opts = append(opts, &nodes.SecurityPrincipalOption{
					Name:  "FROM",
					Value: src,
					Loc:   nodes.Loc{Start: optLoc, End: p.prevEnd()},
				})
			}
		}
		// After FROM WINDOWS, optional WITH clause
		if p.cur.Type == kwWITH {
			p.advance()
			opts = append(opts, p.parseSecurityPrincipalWithOptions()...)
		}
	}

	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseSecurityRoleStmt parses CREATE/ALTER/DROP ROLE and ADD/DROP ROLE MEMBER.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-role-transact-sql
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-role-transact-sql
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-role-transact-sql
//
//	CREATE ROLE role_name [ AUTHORIZATION owner_name ]
//	ALTER  ROLE role_name { ADD MEMBER member_name | DROP MEMBER member_name | WITH NAME = new_name }
//	DROP   ROLE [ IF EXISTS ] role_name
func (p *Parser) parseSecurityRoleStmt(action string) (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     action,
		ObjectType: "ROLE",
		Loc:        nodes.Loc{Start: loc, End: -1},
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
		if p.cur.Type == kwAUTHORIZATION {
			optLoc := p.pos()
			p.advance()
			if owner, ok := p.parseIdentifier(); ok {
				opts = append(opts, &nodes.SecurityPrincipalOption{
					Name:  "AUTHORIZATION",
					Value: owner,
					Loc:   nodes.Loc{Start: optLoc, End: p.prevEnd()},
				})
			}
		}
	} else if action == "ALTER" {
		// ADD MEMBER | DROP MEMBER | WITH NAME = new_name
		if p.cur.Type == kwADD {
			optLoc := p.pos()
			p.advance() // consume ADD
			if p.cur.Type == kwMEMBER {
				p.advance() // consume MEMBER
			}
			if member, ok := p.parseIdentifier(); ok {
				opts = append(opts, &nodes.SecurityPrincipalOption{
					Name:  "ADD MEMBER",
					Value: member,
					Loc:   nodes.Loc{Start: optLoc, End: p.prevEnd()},
				})
			}
		} else if p.cur.Type == kwDROP {
			optLoc := p.pos()
			p.advance() // consume DROP
			if p.cur.Type == kwMEMBER {
				p.advance() // consume MEMBER
			}
			if member, ok := p.parseIdentifier(); ok {
				opts = append(opts, &nodes.SecurityPrincipalOption{
					Name:  "DROP MEMBER",
					Value: member,
					Loc:   nodes.Loc{Start: optLoc, End: p.prevEnd()},
				})
			}
		} else if p.cur.Type == kwWITH {
			p.advance() // consume WITH
			opts = append(opts, p.parseSecurityPrincipalWithOptions()...)
		}
	}

	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseSecurityApplicationRoleStmt parses CREATE/ALTER/DROP APPLICATION ROLE.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-application-role-transact-sql
//
//	CREATE APPLICATION ROLE role_name WITH PASSWORD = 'password' [ , DEFAULT_SCHEMA = schema_name ]
//	ALTER  APPLICATION ROLE role_name WITH NAME = new_name | PASSWORD = '...' | DEFAULT_SCHEMA = ...
//	DROP   APPLICATION ROLE role_name
func (p *Parser) parseSecurityApplicationRoleStmt(action string) (*nodes.SecurityStmt, error) {
	loc := p.pos()
	// APPLICATION keyword already consumed by caller
	// Consume ROLE
	if p.cur.Type == kwROLE {
		p.advance()
	}

	stmt := &nodes.SecurityStmt{
		Action:     action,
		ObjectType: "APPLICATION ROLE",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// role name
	if name, ok := p.parseIdentifier(); ok {
		stmt.Name = name
	}

	var opts []nodes.Node

	// WITH options
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		opts = append(opts, p.parseSecurityPrincipalWithOptions()...)
	}

	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseSecurityPrincipalWithOptions parses comma-separated key=value options
// after WITH has been consumed. Used by CREATE/ALTER USER, LOGIN, ROLE, APPLICATION ROLE.
//
//	option ::=
//	    PASSWORD = 'password' [ HASHED ] [ MUST_CHANGE ] [ OLD_PASSWORD = 'old' ]
//	  | DEFAULT_SCHEMA = schema_name
//	  | DEFAULT_LANGUAGE = { NONE | lcid | language_name }
//	  | DEFAULT_DATABASE = database_name
//	  | CHECK_EXPIRATION = { ON | OFF }
//	  | CHECK_POLICY = { ON | OFF }
//	  | CREDENTIAL = credential_name
//	  | SID = sid
//	  | NAME = new_name
//	  | LOGIN = login_name
//	  | ALLOW_ENCRYPTED_VALUE_MODIFICATIONS = { ON | OFF }
func (p *Parser) parseSecurityPrincipalWithOptions() []nodes.Node {
	var opts []nodes.Node

	for {
		if !p.isAnyKeywordIdent() {
			break
		}

		optLoc := p.pos()
		key := strings.ToUpper(p.cur.Str)
		p.advance() // consume key

		opt := &nodes.SecurityPrincipalOption{
			Name: key,
			Loc:  nodes.Loc{Start: optLoc, End: -1},
		}

		if p.cur.Type == '=' {
			p.advance() // consume =
			// Value: string constant, number, identifier, ON/OFF
			if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST {
				opt.Value = p.cur.Str
				p.advance()
			} else if p.cur.Type == tokICONST || p.cur.Type == tokFCONST {
				opt.Value = p.cur.Str
				p.advance()
			} else if p.cur.Type == kwON {
				opt.Value = "ON"
				p.advance()
			} else if p.cur.Type == kwOFF {
				opt.Value = "OFF"
				p.advance()
			} else if p.isAnyKeywordIdent() {
				upper := strings.ToUpper(p.cur.Str)
				if upper == "ON" {
					opt.Value = "ON"
				} else if upper == "OFF" {
					opt.Value = "OFF"
				} else {
					opt.Value = p.cur.Str
				}
				p.advance()
			}
		}

		// PASSWORD sub-options: HASHED, MUST_CHANGE, UNLOCK, OLD_PASSWORD
		if key == "PASSWORD" {
			for {
				if p.cur.Type == kwHASHED {
					opt.Hashed = true
					p.advance()
				} else if p.cur.Type == kwMUST_CHANGE {
					opt.MustChange = true
					p.advance()
				} else if p.cur.Type == kwUNLOCK {
					opt.Unlock = true
					p.advance()
				} else if p.cur.Type == kwOLD_PASSWORD {
					p.advance() // consume OLD_PASSWORD
					if p.cur.Type == '=' {
						p.advance()
						if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST {
							opt.OldPassword = p.cur.Str
							p.advance()
						}
					}
				} else {
					break
				}
			}
		}

		// Handle NO CREDENTIAL (two-word option in ALTER LOGIN WITH)
		if key == "NO" && p.cur.Type == kwCREDENTIAL {
			opt.Name = "NO CREDENTIAL"
			p.advance() // consume CREDENTIAL
		}

		opt.Loc.End = p.prevEnd()
		opts = append(opts, opt)

		if p.cur.Type == ',' {
			p.advance()
		} else {
			break
		}
	}

	return opts
}

// parseExecuteAsStmt parses a standalone EXECUTE AS statement (context switching).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/execute-as-transact-sql
//
//	EXECUTE AS { LOGIN | USER | CALLER | SELF | OWNER } [ = 'name' ]
//	    [ WITH { NO REVERT | COOKIE INTO @cookie_variable } ]
func (p *Parser) parseExecuteAsStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	p.advance() // consume EXECUTE
	p.advance() // consume AS

	stmt := &nodes.SecurityStmt{
		Action: "EXECUTE AS",
		Loc:    nodes.Loc{Start: loc, End: -1},
	}

	// { LOGIN | USER | CALLER | SELF | OWNER }
	if p.cur.Type == kwLOGIN {
		stmt.ObjectType = "LOGIN"
		p.advance()
	} else if p.cur.Type == kwUSER {
		stmt.ObjectType = "USER"
		p.advance()
	} else if p.cur.Type == kwCALLER {
		stmt.ObjectType = "CALLER"
		p.advance()
		stmt.Loc.End = p.prevEnd()
		return stmt, nil
	} else if p.cur.Type == kwSELF {
		stmt.ObjectType = "SELF"
		p.advance()
		stmt.Loc.End = p.prevEnd()
		return stmt, nil
	} else if p.cur.Type == kwOWNER {
		stmt.ObjectType = "OWNER"
		p.advance()
		stmt.Loc.End = p.prevEnd()
		return stmt, nil
	}

	// = 'name'
	if p.cur.Type == '=' {
		p.advance()
		if p.cur.Type == tokSCONST || p.isIdentLike() {
			stmt.Name = p.cur.Str
			p.advance()
		}
	}

	// WITH { NO REVERT | COOKIE INTO @cookie_variable }
	if p.cur.Type == kwWITH {
		p.advance()
		var opts []nodes.Node
		if p.cur.Type == kwNO {
			p.advance()
			if p.cur.Type == kwREVERT {
				p.advance()
				opts = append(opts, &nodes.String{Str: "NO REVERT"})
			}
		} else if p.cur.Type == kwCOOKIE {
			p.advance()
			if p.cur.Type == kwINTO {
				p.advance()
			}
			if p.cur.Type == tokVARIABLE {
				opts = append(opts, &nodes.String{Str: "COOKIE=" + p.cur.Str})
				p.advance()
			}
		}
		if len(opts) > 0 {
			stmt.Options = &nodes.List{Items: opts}
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseRevertStmt parses the REVERT statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/revert-transact-sql
//
//	REVERT [ WITH COOKIE = @cookie_variable ]
func (p *Parser) parseRevertStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	p.advance() // consume REVERT

	stmt := &nodes.SecurityStmt{
		Action:     "REVERT",
		ObjectType: "CONTEXT",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// WITH COOKIE = @cookie_variable
	if p.cur.Type == kwWITH {
		p.advance()
		if p.cur.Type == kwCOOKIE {
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
				if p.cur.Type == tokVARIABLE {
					var opts []nodes.Node
					opts = append(opts, &nodes.String{Str: "COOKIE=" + p.cur.Str})
					stmt.Options = &nodes.List{Items: opts}
					p.advance()
				}
			}
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterAuthorizationStmt parses ALTER AUTHORIZATION.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-authorization-transact-sql
//
//	ALTER AUTHORIZATION
//	    ON [ <class_type> :: ] entity_name
//	    TO { principal_name | SCHEMA OWNER }
//
//	<class_type> ::=
//	    { OBJECT | ASSEMBLY | ASYMMETRIC KEY | CERTIFICATE | CONTRACT
//	    | DATABASE | ENDPOINT | FULLTEXT CATALOG | FULLTEXT STOPLIST
//	    | MESSAGE TYPE | REMOTE SERVICE BINDING | ROLE | ROUTE
//	    | SCHEMA | SEARCH PROPERTY LIST | SERVER ROLE | SERVICE
//	    | SYMMETRIC KEY | TYPE | XML SCHEMA COLLECTION }
func (p *Parser) parseAlterAuthorizationStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	// AUTHORIZATION keyword already consumed by caller

	stmt := &nodes.SecurityStmt{
		Action:     "ALTER AUTHORIZATION",
		ObjectType: "ALTER AUTHORIZATION",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// ON [ entity_type :: ] entity_name
	if p.cur.Type == kwON {
		p.advance()
		// Try to detect entity_type :: pattern
		// Entity types can be multi-word: XML SCHEMA COLLECTION, SEARCH PROPERTY LIST, etc.
		entityType := p.tryParseAlterAuthEntityType()
		if entityType != "" {
			// We have entity_type, now expect ::
			if p.cur.Type == tokCOLONCOLON {
				p.advance() // consume ::
			}
			stmt.Name = entityType + "::" + p.parseAlterAuthEntityName()
		} else {
			// No entity_type, just parse entity_name
			stmt.Name = p.parseAlterAuthEntityName()
		}
	}

	// TO { principal_name | SCHEMA OWNER }
	if p.cur.Type == kwTO {
		p.advance()
		var opts []nodes.Node
		if p.cur.Type == kwSCHEMA {
			p.advance()
			if p.cur.Type == kwOWNER {
				p.advance()
				opts = append(opts, &nodes.String{Str: "SCHEMA OWNER"})
			}
		} else if p.isIdentLike() || p.cur.Type == tokSCONST {
			opts = append(opts, &nodes.String{Str: p.cur.Str})
			p.advance()
		}
		if len(opts) > 0 {
			stmt.Options = &nodes.List{Items: opts}
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// tryParseAlterAuthEntityType checks if the current position has an entity_type :: pattern.
// If so, consumes the entity type tokens and returns the entity type string.
// Otherwise returns empty string without consuming anything.
//
// Known single-word entity types: OBJECT, ASSEMBLY, CONTRACT, CERTIFICATE, ENDPOINT,
// ROUTE, SERVICE, TYPE, SCHEMA, DATABASE, ROLE
// Multi-word: XML SCHEMA COLLECTION, MESSAGE TYPE, REMOTE SERVICE BINDING,
// FULLTEXT CATALOG, FULLTEXT STOPLIST, SEARCH PROPERTY LIST,
// SYMMETRIC KEY, ASYMMETRIC KEY, AVAILABILITY GROUP, SERVER ROLE
func (p *Parser) tryParseAlterAuthEntityType() string {
	// For single-word entity types, check if next token is ':'
	if (p.isAnyKeywordIdent() || p.cur.Type == kwDATABASE || p.cur.Type == kwXML ||
		p.cur.Type == kwROLE || p.cur.Type == kwSCHEMA || p.cur.Type == kwTYPE) &&
		p.peekNext().Type == tokCOLONCOLON {
		word := strings.ToUpper(p.cur.Str)
		p.advance() // consume entity type word
		return word
	}
	// For multi-word entity types: check first word, then peek second
	if p.isAnyKeywordIdent() || p.cur.Type == kwXML || p.cur.Type == kwSCHEMA {
		first := strings.ToUpper(p.cur.Str)
		switch first {
		case "XML":
			return p.tryParseMultiWordEntityType(first, []string{"SCHEMA", "COLLECTION"})
		case "MESSAGE":
			return p.tryParseMultiWordEntityType(first, []string{"TYPE"})
		case "REMOTE":
			return p.tryParseMultiWordEntityType(first, []string{"SERVICE", "BINDING"})
		case "FULLTEXT":
			next := p.peekNext()
			if next.Type == kwCATALOG || next.Type == kwSTOPLIST {
				return p.tryParseMultiWordEntityType(first, []string{strings.ToUpper(next.Str)})
			}
		case "SEARCH":
			return p.tryParseMultiWordEntityType(first, []string{"PROPERTY", "LIST"})
		case "SYMMETRIC", "ASYMMETRIC":
			return p.tryParseMultiWordEntityType(first, []string{"KEY"})
		case "AVAILABILITY":
			return p.tryParseMultiWordEntityType(first, []string{"GROUP"})
		case "SERVER":
			return p.tryParseMultiWordEntityType(first, []string{"ROLE"})
		}
	}
	return ""
}

// tryParseMultiWordEntityType consumes a multi-word entity type if the next word matches.
// Validates by peeking before consuming. Returns the entity type or "".
func (p *Parser) tryParseMultiWordEntityType(first string, remaining []string) string {
	if len(remaining) == 0 {
		return ""
	}
	// Peek to confirm next token matches first remaining word
	next := p.peekNext()
	if next.Type != lookupKeyword(remaining[0]) {
		return ""
	}

	// Pattern looks good, consume first word
	p.advance()
	parts := []string{first}
	for _, word := range remaining {
		if (p.isAnyKeywordIdent() || p.cur.Type == kwSCHEMA || p.cur.Type == kwTYPE) &&
			p.cur.Type == lookupKeyword(word) {
			parts = append(parts, strings.ToUpper(p.cur.Str))
			p.advance()
		} else {
			break
		}
	}
	return strings.Join(parts, " ")
}

// parseAlterAuthEntityName parses the entity name (possibly qualified with dots) until TO keyword.
func (p *Parser) parseAlterAuthEntityName() string {
	var parts []string
	for p.cur.Type != tokEOF && p.cur.Type != ';' {
		if p.cur.Type == kwTO {
			break
		}
		if p.isAnyKeywordIdent() || p.cur.Type == kwDATABASE || p.cur.Type == kwXML ||
			p.cur.Type == kwROLE || p.cur.Type == kwSCHEMA || p.cur.Type == kwTYPE {
			parts = append(parts, p.cur.Str)
		} else if p.cur.Type == '.' {
			parts = append(parts, ".")
		} else {
			break
		}
		p.advance()
	}
	return strings.Join(parts, " ")
}

