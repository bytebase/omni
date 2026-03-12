// Package parser - security_audit.go implements T-SQL audit statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseCreateServerAuditStmt parses CREATE SERVER AUDIT.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-server-audit-transact-sql
//
//	CREATE SERVER AUDIT audit_name
//	    {
//	        TO { FILE ( <file_options> [,...n] ) | APPLICATION_LOG | SECURITY_LOG | URL | EXTERNAL_MONITOR }
//	        [ WITH ( <audit_options> [,...n] ) ]
//	        [ WHERE <predicate_expression> ]
//	    }
//	<file_options> ::=
//	    { FILEPATH = 'os_file_path'
//	      [, MAXSIZE = { max_size { MB | GB | TB } | UNLIMITED } ]
//	      [, { MAX_ROLLOVER_FILES = { integer | UNLIMITED } } | { MAX_FILES = integer } ]
//	      [, RESERVE_DISK_SPACE = { ON | OFF } ]
//	    }
//	<audit_options> ::=
//	    { QUEUE_DELAY = integer
//	      [, ON_FAILURE = { CONTINUE | SHUTDOWN | FAIL_OPERATION } ]
//	      [, AUDIT_GUID = uniqueidentifier
//	    }
func (p *Parser) parseCreateServerAuditStmt() *nodes.SecurityStmt {
	loc := p.pos()
	// SERVER AUDIT keywords already consumed by caller

	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "SERVER AUDIT",
		Loc:        nodes.Loc{Start: loc},
	}

	// audit_name
	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// Consume rest of statement (TO, WITH, WHERE clauses)
	stmt.Options = p.parseAuditOptions()

	stmt.Loc.End = p.pos()
	return stmt
}

// parseAlterServerAuditStmt parses ALTER SERVER AUDIT.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-server-audit-transact-sql
//
//	ALTER SERVER AUDIT audit_name
//	    {
//	        [ TO { FILE ( <file_options> [,...n] ) | APPLICATION_LOG | SECURITY_LOG | URL | EXTERNAL_MONITOR } ]
//	        [ WITH ( <audit_options> [,...n] ) ]
//	        [ WHERE <predicate_expression> ]
//	    }
//	  | { REMOVE WHERE }
//	  | { MODIFY NAME = new_audit_name }
func (p *Parser) parseAlterServerAuditStmt() *nodes.SecurityStmt {
	loc := p.pos()
	// SERVER AUDIT keywords already consumed by caller

	stmt := &nodes.SecurityStmt{
		Action:     "ALTER",
		ObjectType: "SERVER AUDIT",
		Loc:        nodes.Loc{Start: loc},
	}

	// audit_name
	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Options = p.parseAuditOptions()

	stmt.Loc.End = p.pos()
	return stmt
}

// parseDropServerAuditStmt parses DROP SERVER AUDIT.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-server-audit-transact-sql
//
//	DROP SERVER AUDIT audit_name
func (p *Parser) parseDropServerAuditStmt() *nodes.SecurityStmt {
	loc := p.pos()
	// SERVER AUDIT keywords already consumed by caller

	stmt := &nodes.SecurityStmt{
		Action:     "DROP",
		ObjectType: "SERVER AUDIT",
		Loc:        nodes.Loc{Start: loc},
	}

	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseCreateServerAuditSpecStmt parses CREATE SERVER AUDIT SPECIFICATION.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-server-audit-specification-transact-sql
//
//	CREATE SERVER AUDIT SPECIFICATION audit_specification_name
//	    FOR SERVER AUDIT audit_name
//	    { { ADD ( audit_action_group_name ) } [,...n] }
//	    [ WITH ( STATE = { ON | OFF } ) ]
func (p *Parser) parseCreateServerAuditSpecStmt() *nodes.SecurityStmt {
	loc := p.pos()
	// SPECIFICATION keyword already consumed by caller

	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "SERVER AUDIT SPECIFICATION",
		Loc:        nodes.Loc{Start: loc},
	}

	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Options = p.parseAuditSpecOptions()

	stmt.Loc.End = p.pos()
	return stmt
}

// parseAlterServerAuditSpecStmt parses ALTER SERVER AUDIT SPECIFICATION.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-server-audit-specification-transact-sql
//
//	ALTER SERVER AUDIT SPECIFICATION audit_specification_name
//	    FOR SERVER AUDIT audit_name
//	    { { ADD ( audit_action_group_name ) }
//	      | { DROP ( audit_action_group_name ) } } [,...n]
//	    [ WITH ( STATE = { ON | OFF } ) ]
func (p *Parser) parseAlterServerAuditSpecStmt() *nodes.SecurityStmt {
	loc := p.pos()

	stmt := &nodes.SecurityStmt{
		Action:     "ALTER",
		ObjectType: "SERVER AUDIT SPECIFICATION",
		Loc:        nodes.Loc{Start: loc},
	}

	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Options = p.parseAuditSpecOptions()

	stmt.Loc.End = p.pos()
	return stmt
}

// parseDropServerAuditSpecStmt parses DROP SERVER AUDIT SPECIFICATION.
//
//	DROP SERVER AUDIT SPECIFICATION audit_specification_name
func (p *Parser) parseDropServerAuditSpecStmt() *nodes.SecurityStmt {
	loc := p.pos()

	stmt := &nodes.SecurityStmt{
		Action:     "DROP",
		ObjectType: "SERVER AUDIT SPECIFICATION",
		Loc:        nodes.Loc{Start: loc},
	}

	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseCreateDatabaseAuditSpecStmt parses CREATE DATABASE AUDIT SPECIFICATION.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-database-audit-specification-transact-sql
//
//	CREATE DATABASE AUDIT SPECIFICATION audit_specification_name
//	    FOR SERVER AUDIT audit_name
//	    { { ADD ( { <audit_action_specification> | audit_action_group_name }
//	        ) } [,...n] }
//	    [ WITH ( STATE = { ON | OFF } ) ]
//
//	<audit_action_specification> ::=
//	    action [ ,...n ] ON [ class :: ] securable BY principal [ ,...n ]
func (p *Parser) parseCreateDatabaseAuditSpecStmt() *nodes.SecurityStmt {
	loc := p.pos()
	// SPECIFICATION keyword already consumed by caller

	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "DATABASE AUDIT SPECIFICATION",
		Loc:        nodes.Loc{Start: loc},
	}

	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Options = p.parseAuditSpecOptions()

	stmt.Loc.End = p.pos()
	return stmt
}

// parseAlterDatabaseAuditSpecStmt parses ALTER DATABASE AUDIT SPECIFICATION.
//
//	ALTER DATABASE AUDIT SPECIFICATION audit_specification_name
//	    FOR SERVER AUDIT audit_name
//	    { { ADD | DROP } ( { <audit_action_specification> | audit_action_group_name } ) } [,...n]
//	    [ WITH ( STATE = { ON | OFF } ) ]
func (p *Parser) parseAlterDatabaseAuditSpecStmt() *nodes.SecurityStmt {
	loc := p.pos()

	stmt := &nodes.SecurityStmt{
		Action:     "ALTER",
		ObjectType: "DATABASE AUDIT SPECIFICATION",
		Loc:        nodes.Loc{Start: loc},
	}

	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Options = p.parseAuditSpecOptions()

	stmt.Loc.End = p.pos()
	return stmt
}

// parseDropDatabaseAuditSpecStmt parses DROP DATABASE AUDIT SPECIFICATION.
//
//	DROP DATABASE AUDIT SPECIFICATION audit_specification_name
func (p *Parser) parseDropDatabaseAuditSpecStmt() *nodes.SecurityStmt {
	loc := p.pos()

	stmt := &nodes.SecurityStmt{
		Action:     "DROP",
		ObjectType: "DATABASE AUDIT SPECIFICATION",
		Loc:        nodes.Loc{Start: loc},
	}

	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseAuditOptions parses the TO/WITH/WHERE portions of a SERVER AUDIT statement.
func (p *Parser) parseAuditOptions() *nodes.List {
	var opts []nodes.Node

	// Consume rest of statement until ; or EOF or GO
	for p.cur.Type != ';' && p.cur.Type != tokEOF && p.cur.Type != kwGO {
		if p.cur.Type == '(' {
			// Skip parenthesized sub-clause
			p.advance()
			depth := 1
			for depth > 0 && p.cur.Type != tokEOF {
				if p.cur.Type == '(' {
					depth++
				} else if p.cur.Type == ')' {
					depth--
				}
				if depth > 0 {
					p.advance()
				}
			}
			p.match(')')
		} else if p.isIdentLike() || p.cur.Type == kwWITH || p.cur.Type == kwWHERE ||
			p.cur.Type == kwON || p.cur.Type == kwOFF || p.cur.Type == kwTO {
			optStr := strings.ToUpper(p.cur.Str)
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
				if p.isIdentLike() || p.cur.Type == tokSCONST || p.cur.Type == tokICONST ||
					p.cur.Type == kwON || p.cur.Type == kwOFF {
					optStr += "=" + strings.ToUpper(p.cur.Str)
					p.advance()
				}
			}
			opts = append(opts, &nodes.String{Str: optStr})
		} else {
			p.advance()
		}
	}

	if len(opts) == 0 {
		return nil
	}
	return &nodes.List{Items: opts}
}

// parseAuditSpecOptions parses FOR SERVER AUDIT / ADD / DROP / WITH STATE portions.
func (p *Parser) parseAuditSpecOptions() *nodes.List {
	var opts []nodes.Node

	// FOR SERVER AUDIT audit_name
	if p.cur.Type == kwFOR {
		p.advance()
		// SERVER
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SERVER") {
			p.advance()
		}
		// AUDIT
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "AUDIT") {
			p.advance()
		}
		// audit_name
		if p.isIdentLike() || p.cur.Type == tokSCONST {
			opts = append(opts, &nodes.String{Str: "FOR_AUDIT=" + p.cur.Str})
			p.advance()
		}
	}

	// ADD/DROP ( ... ) clauses and WITH
	for p.cur.Type != ';' && p.cur.Type != tokEOF && p.cur.Type != kwGO {
		if p.cur.Type == kwADD || p.cur.Type == kwDROP {
			action := strings.ToUpper(p.cur.Str)
			p.advance()
			if p.cur.Type == '(' {
				p.advance()
				// Consume until matching )
				depth := 1
				var innerParts []string
				for depth > 0 && p.cur.Type != tokEOF {
					if p.cur.Type == '(' {
						depth++
					} else if p.cur.Type == ')' {
						depth--
						if depth == 0 {
							break
						}
					}
					if p.isIdentLike() || p.cur.Type == kwON || p.cur.Type == kwSELECT ||
						p.cur.Type == kwINSERT || p.cur.Type == kwUPDATE || p.cur.Type == kwDELETE ||
						p.cur.Type == kwEXECUTE || p.cur.Type == kwEXEC {
						innerParts = append(innerParts, strings.ToUpper(p.cur.Str))
					}
					p.advance()
				}
				p.match(')')
				if len(innerParts) > 0 {
					opts = append(opts, &nodes.String{Str: action + "(" + strings.Join(innerParts, " ") + ")"})
				}
			}
		} else if p.cur.Type == kwWITH {
			p.advance()
			if p.cur.Type == '(' {
				p.advance()
				for p.cur.Type != ')' && p.cur.Type != tokEOF {
					if p.isIdentLike() || p.cur.Type == kwON || p.cur.Type == kwOFF {
						optName := strings.ToUpper(p.cur.Str)
						p.advance()
						if p.cur.Type == '=' {
							p.advance()
							if p.isIdentLike() || p.cur.Type == kwON || p.cur.Type == kwOFF {
								optName += "=" + strings.ToUpper(p.cur.Str)
								p.advance()
							}
						}
						opts = append(opts, &nodes.String{Str: optName})
					} else {
						p.advance()
					}
					p.match(',')
				}
				p.match(')')
			}
			break
		} else {
			break
		}
		p.match(',')
	}

	if len(opts) == 0 {
		return nil
	}
	return &nodes.List{Items: opts}
}
