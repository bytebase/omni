// Package parser - security_audit.go implements T-SQL audit statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseCreateServerAuditStmt parses CREATE SERVER AUDIT.
//
// BNF: mssql/parser/bnf/create-server-audit-transact-sql.bnf
//
//	CREATE SERVER AUDIT audit_name
//	{
//	    TO { FILE ( <file_options> [,...n] ) | APPLICATION_LOG | SECURITY_LOG
//	       | URL ( <file_options> [,...n] ) | EXTERNAL_MONITOR }
//	    [ WITH ( <audit_options> [,...n] ) ]
//	    [ WHERE <predicate_expression> ]
//	}
//
//	<file_options> ::=
//	{
//	    FILEPATH = 'os_file_path'
//	    [, MAXSIZE = { max_size { MB | GB | TB } | UNLIMITED } ]
//	    [, { MAX_ROLLOVER_FILES = { integer | UNLIMITED } } | { MAX_FILES = integer } ]
//	    [, RESERVE_DISK_SPACE = { ON | OFF } ]
//	}
//
//	<audit_options> ::=
//	{
//	    [ QUEUE_DELAY = integer ]
//	    [, ON_FAILURE = { CONTINUE | SHUTDOWN | FAIL_OPERATION } ]
//	    [, AUDIT_GUID = uniqueidentifier ]
//	    [, OPERATOR_AUDIT = { ON | OFF } ]
//	}
//
//	<predicate_expression> ::=
//	{
//	    [ NOT ] <predicate_factor>
//	    [ { AND | OR } [ NOT ] { <predicate_factor> } ]
//	    [, ...n]
//	}
//
//	<predicate_factor> ::=
//	    event_field_name { = | <> | != | > | >= | < | <= | LIKE } { number | 'string' }
func (p *Parser) parseCreateServerAuditStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	// SERVER AUDIT keywords already consumed by caller

	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "SERVER AUDIT",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// audit_name
	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// Consume rest of statement (TO, WITH, WHERE clauses)
	stmt.Options, stmt.WhereClause = p.parseAuditOptions()

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterServerAuditStmt parses ALTER SERVER AUDIT.
//
// BNF: mssql/parser/bnf/alter-server-audit-transact-sql.bnf
//
//	ALTER SERVER AUDIT audit_name
//	{
//	    [ TO { FILE ( <file_options> [,...n] ) | APPLICATION_LOG | SECURITY_LOG
//	         | URL ( <file_options> [,...n] ) | EXTERNAL_MONITOR } ]
//	    [ WITH ( <audit_options> [,...n] ) ]
//	    [ WHERE <predicate_expression> ]
//	}
//	| REMOVE WHERE
//	| MODIFY NAME = new_audit_name
//
//	<audit_options> ::=
//	{
//	    [ QUEUE_DELAY = integer ]
//	    [, ON_FAILURE = { CONTINUE | SHUTDOWN | FAIL_OPERATION } ]
//	    [, STATE = { ON | OFF } ]
//	    [, OPERATOR_AUDIT = { ON | OFF } ]
//	}
func (p *Parser) parseAlterServerAuditStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	// SERVER AUDIT keywords already consumed by caller

	stmt := &nodes.SecurityStmt{
		Action:     "ALTER",
		ObjectType: "SERVER AUDIT",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// audit_name
	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Options, stmt.WhereClause = p.parseAuditOptions()

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseDropServerAuditStmt parses DROP SERVER AUDIT.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-server-audit-transact-sql
//
//	DROP SERVER AUDIT audit_name
func (p *Parser) parseDropServerAuditStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	// SERVER AUDIT keywords already consumed by caller

	stmt := &nodes.SecurityStmt{
		Action:     "DROP",
		ObjectType: "SERVER AUDIT",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseCreateServerAuditSpecStmt parses CREATE SERVER AUDIT SPECIFICATION.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-server-audit-specification-transact-sql
//
//	CREATE SERVER AUDIT SPECIFICATION audit_specification_name
//	    FOR SERVER AUDIT audit_name
//	    { { ADD ( audit_action_group_name ) } [,...n] }
//	    [ WITH ( STATE = { ON | OFF } ) ]
func (p *Parser) parseCreateServerAuditSpecStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	// SPECIFICATION keyword already consumed by caller

	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "SERVER AUDIT SPECIFICATION",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Options = p.parseAuditSpecOptions()

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
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
func (p *Parser) parseAlterServerAuditSpecStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()

	stmt := &nodes.SecurityStmt{
		Action:     "ALTER",
		ObjectType: "SERVER AUDIT SPECIFICATION",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Options = p.parseAuditSpecOptions()

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseDropServerAuditSpecStmt parses DROP SERVER AUDIT SPECIFICATION.
//
//	DROP SERVER AUDIT SPECIFICATION audit_specification_name
func (p *Parser) parseDropServerAuditSpecStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()

	stmt := &nodes.SecurityStmt{
		Action:     "DROP",
		ObjectType: "SERVER AUDIT SPECIFICATION",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseCreateDatabaseAuditSpecStmt parses CREATE DATABASE AUDIT SPECIFICATION.
//
// BNF: mssql/parser/bnf/create-database-audit-specification-transact-sql.bnf
//
//	CREATE DATABASE AUDIT SPECIFICATION audit_specification_name
//	    FOR SERVER AUDIT audit_name
//	    { { ADD ( { <audit_action_specification> | audit_action_group_name }
//	        ) } [,...n] }
//	    [ WITH ( STATE = { ON | OFF } ) ]
//
//	<audit_action_specification> ::=
//	    action [ ,...n ] ON [ class :: ] securable BY principal [ ,...n ]
func (p *Parser) parseCreateDatabaseAuditSpecStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	// SPECIFICATION keyword already consumed by caller

	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "DATABASE AUDIT SPECIFICATION",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Options = p.parseAuditSpecOptions()

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterDatabaseAuditSpecStmt parses ALTER DATABASE AUDIT SPECIFICATION.
//
//	ALTER DATABASE AUDIT SPECIFICATION audit_specification_name
//	    FOR SERVER AUDIT audit_name
//	    { { ADD | DROP } ( { <audit_action_specification> | audit_action_group_name } ) } [,...n]
//	    [ WITH ( STATE = { ON | OFF } ) ]
func (p *Parser) parseAlterDatabaseAuditSpecStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()

	stmt := &nodes.SecurityStmt{
		Action:     "ALTER",
		ObjectType: "DATABASE AUDIT SPECIFICATION",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Options = p.parseAuditSpecOptions()

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseDropDatabaseAuditSpecStmt parses DROP DATABASE AUDIT SPECIFICATION.
//
//	DROP DATABASE AUDIT SPECIFICATION audit_specification_name
func (p *Parser) parseDropDatabaseAuditSpecStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()

	stmt := &nodes.SecurityStmt{
		Action:     "DROP",
		ObjectType: "DATABASE AUDIT SPECIFICATION",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAuditOptions parses the TO/WITH/WHERE portions of a SERVER AUDIT statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-server-audit-transact-sql
//
//	TO { FILE ( <file_options> [, ...n] ) | APPLICATION_LOG | SECURITY_LOG | URL | EXTERNAL_MONITOR }
//	[ WITH ( <audit_options> [, ...n] ) ]
//	[ WHERE <predicate_expression> ]
//	| REMOVE WHERE
//	| MODIFY NAME = new_audit_name
//
//	<file_options> ::=
//	    FILEPATH = 'os_file_path'
//	    [, MAXSIZE = { max_size { MB | GB | TB } | UNLIMITED } ]
//	    [, { MAX_ROLLOVER_FILES = { integer | UNLIMITED } } | { MAX_FILES = integer } ]
//	    [, RESERVE_DISK_SPACE = { ON | OFF } ]
//
//	<audit_options> ::=
//	    [ QUEUE_DELAY = integer ]
//	    [, ON_FAILURE = { CONTINUE | SHUTDOWN | FAIL_OPERATION } ]
//	    [, AUDIT_GUID = uniqueidentifier ]
//	    [, OPERATOR_AUDIT = { ON | OFF } ]
//	    [, STATE = ON | OFF ]
//
//	<predicate_expression> ::=
//	    [ NOT ] <predicate_factor> [ { AND | OR } [ NOT ] <predicate_factor> ] [, ...n]
//	<predicate_factor> ::=
//	    event_field_name { = | <> | != | > | >= | < | <= | LIKE } { number | 'string' }
func (p *Parser) parseAuditOptions() (*nodes.List, nodes.ExprNode) {
	var opts []nodes.Node

	// REMOVE WHERE (ALTER only)
	if p.cur.Type == kwREMOVE {
		p.advance()
		if p.cur.Type == kwWHERE {
			p.advance()
		}
		opts = append(opts, &nodes.String{Str: "REMOVE WHERE"})
		if len(opts) == 0 {
			return nil, nil
		}
		return &nodes.List{Items: opts}, nil
	}

	// MODIFY NAME = new_name (ALTER only)
	if p.cur.Type == kwMODIFY {
		p.advance()
		if p.cur.Type == kwNAME {
			p.advance()
		}
		if p.cur.Type == '=' {
			p.advance()
		}
		if p.isIdentLike() || p.cur.Type == tokSCONST {
			opts = append(opts, &nodes.String{Str: "MODIFY NAME=" + p.cur.Str})
			p.advance()
		}
		if len(opts) == 0 {
			return nil, nil
		}
		return &nodes.List{Items: opts}, nil
	}

	// TO clause
	if p.cur.Type == kwTO {
		p.advance()
		if p.cur.Type == kwFILE || p.cur.Type == kwURL {
			target := strings.ToUpper(p.cur.Str)
			p.advance()
			opts = append(opts, &nodes.String{Str: "TO=" + target})
			// ( <file_options> ) -- both FILE and URL support parenthesized options
			if p.cur.Type == '(' {
				p.advance()
				for p.cur.Type != ')' && p.cur.Type != tokEOF {
					if p.cur.Type == ',' {
						p.advance()
						continue
					}
					if p.isIdentLike() || p.cur.Type == kwON || p.cur.Type == kwOFF {
						name := strings.ToUpper(p.cur.Str)
						p.advance()
						if p.cur.Type == '=' {
							p.advance()
							val := ""
							if p.cur.Type == tokSCONST {
								val = p.cur.Str
								p.advance()
							} else if p.cur.Type == tokICONST {
								val = p.cur.Str
								p.advance()
								// MAXSIZE may have MB|GB|TB suffix
								if p.cur.Type == kwMB ||
									p.cur.Type == kwGB ||
									p.cur.Type == kwTB {
									val += strings.ToUpper(p.cur.Str)
									p.advance()
								}
							} else if p.isIdentLike() {
								val = strings.ToUpper(p.cur.Str)
								p.advance()
							} else if p.cur.Type == kwON {
								val = "ON"
								p.advance()
							} else if p.cur.Type == kwOFF {
								val = "OFF"
								p.advance()
							}
							opts = append(opts, &nodes.String{Str: name + "=" + val})
						}
					} else {
						p.advance()
					}
				}
				p.match(')')
			}
		} else if p.isIdentLike() {
			target := strings.ToUpper(p.cur.Str)
			p.advance()
			opts = append(opts, &nodes.String{Str: "TO=" + target})
		}
	}

	// WITH clause
	if p.cur.Type == kwWITH {
		p.advance()
		if p.cur.Type == '(' {
			p.advance()
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				if p.cur.Type == ',' {
					p.advance()
					continue
				}
				if p.isIdentLike() || p.cur.Type == kwON || p.cur.Type == kwOFF {
					name := strings.ToUpper(p.cur.Str)
					p.advance()
					if p.cur.Type == '=' {
						p.advance()
						val := ""
						if p.cur.Type == tokSCONST {
							val = p.cur.Str
							p.advance()
						} else if p.cur.Type == tokICONST {
							val = p.cur.Str
							p.advance()
						} else if p.isIdentLike() {
							val = strings.ToUpper(p.cur.Str)
							p.advance()
						} else if p.cur.Type == kwON {
							val = "ON"
							p.advance()
						} else if p.cur.Type == kwOFF {
							val = "OFF"
							p.advance()
						}
						opts = append(opts, &nodes.String{Str: name + "=" + val})
					} else {
						opts = append(opts, &nodes.String{Str: name})
					}
				} else {
					p.advance()
				}
			}
			p.match(')')
		}
	}

	// WHERE clause - parse as a proper expression
	var whereClause nodes.ExprNode
	if p.cur.Type == kwWHERE {
		p.advance()
		whereClause, _ = p.parseExpr()
	}

	if len(opts) == 0 {
		return nil, whereClause
	}
	return &nodes.List{Items: opts}, whereClause
}

// parseAuditSpecOptions parses FOR SERVER AUDIT / ADD / DROP / WITH STATE portions.
func (p *Parser) parseAuditSpecOptions() *nodes.List {
	var opts []nodes.Node

	// FOR SERVER AUDIT audit_name
	if p.cur.Type == kwFOR {
		p.advance()
		// SERVER
		if p.cur.Type == kwSERVER {
			p.advance()
		}
		// AUDIT
		if p.cur.Type == kwAUDIT {
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
			opts = append(opts, p.parseAuditSpecAction())
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

// parseAuditSpecAction parses a single ADD or DROP action in an audit specification.
//
//	ADD ( { <audit_action_specification> | audit_action_group_name } )
//	DROP ( { <audit_action_specification> | audit_action_group_name } )
//
//	<audit_action_specification> ::=
//	    action [ ,...n ] ON [ class :: ] securable BY principal [ ,...n ]
func (p *Parser) parseAuditSpecAction() *nodes.AuditSpecAction {
	loc := p.pos()
	action := strings.ToUpper(p.cur.Str)
	p.advance()

	node := &nodes.AuditSpecAction{
		Action: action,
		Loc:    nodes.Loc{Start: loc, End: -1},
	}

	if p.cur.Type != '(' {
		node.Loc.End = p.prevEnd()
		return node
	}
	p.advance()

	// Parse the contents inside parentheses.
	// It's either:
	//   1. audit_action_group_name (e.g., FAILED_LOGIN_GROUP)
	//   2. action [,...n] ON [class ::] securable BY principal [,...n]
	//
	// We distinguish by looking for ON keyword after the first identifier(s).

	// Collect action/group names until ON, BY, or )
	var names []string
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		if p.cur.Type == kwON {
			break
		}
		if p.isIdentLike() || p.cur.Type == kwSELECT || p.cur.Type == kwINSERT ||
			p.cur.Type == kwUPDATE || p.cur.Type == kwDELETE ||
			p.cur.Type == kwEXECUTE || p.cur.Type == kwEXEC {
			names = append(names, strings.ToUpper(p.cur.Str))
			p.advance()
			if p.cur.Type == ',' {
				p.advance()
				continue
			}
		} else {
			break
		}
	}

	if p.cur.Type == kwON {
		// This is an audit_action_specification: actions ON [class::] securable BY principals
		node.Actions = names
		p.advance() // consume ON

		// Parse [class ::] securable
		// Could be: OBJECT::dbo.MyTable, SCHEMA::dbo, DATABASE::mydb, or just dbo.MyTable
		var securableParts []string
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			if p.isAnyKeywordIdent() || p.cur.Type == kwSELECT || p.cur.Type == kwDELETE {
				word := p.cur.Str
				p.advance()
				if p.cur.Type == tokCOLONCOLON {
					// This is a class name
					node.ClassName = strings.ToUpper(word)
					p.advance() // consume ::
					continue
				}
				securableParts = append(securableParts, word)
				// Check for dotted name
				for p.cur.Type == '.' {
					p.advance()
					if p.isIdentLike() || p.cur.Type == kwDEFAULT {
						securableParts = append(securableParts, p.cur.Str)
						p.advance()
					}
				}
			}
			break
		}
		if len(securableParts) > 0 {
			node.Securable = strings.Join(securableParts, ".")
		}

		// BY principal [,...n]
		if p.cur.Type == kwBY {
			p.advance()
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				if p.isIdentLike() || p.cur.Type == kwPUBLIC {
					node.Principals = append(node.Principals, p.cur.Str)
					p.advance()
				}
				if p.cur.Type == ',' {
					p.advance()
					continue
				}
				break
			}
		}
	} else {
		// Simple audit action group name
		if len(names) > 0 {
			node.GroupName = strings.Join(names, "_")
		}
	}

	p.match(')')
	node.Loc.End = p.prevEnd()
	return node
}
