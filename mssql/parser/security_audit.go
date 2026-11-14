// Package parser - security_audit.go implements T-SQL audit statement parsing.
//
// The grammar is ScriptDOM's (mssql/testdata/TSql160.g, cited by line), held
// to what SQL Server 2022 accepts where the two differ: SQL Server refuses the
// Azure-only TO URL and TO EXTERNAL_MONITOR targets and the OPERATOR_AUDIT
// option (Msg 40517, "not supported in this version of SQL Server").
package parser

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseCreateServerAuditStmt parses CREATE SERVER AUDIT; the caller has
// consumed CREATE SERVER AUDIT.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-server-audit-transact-sql
// Grammar: TSql160.g:1690 createServerAuditStatement
//
//	CREATE SERVER AUDIT audit_name
//	    TO { FILE ( <file_option> [ ,...n ] ) | APPLICATION_LOG | SECURITY_LOG }
//	    [ WITH ( <audit_option> [ ,...n ] ) ]
//	    [ WHERE <predicate_expression> ]
//
//	<audit_option> ::= QUEUE_DELAY = integer
//	    | ON_FAILURE = { CONTINUE | SHUTDOWN | FAIL_OPERATION }
//	    | AUDIT_GUID = 'uniqueidentifier'
func (p *Parser) parseCreateServerAuditStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "SERVER AUDIT",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}
	name, err := p.auditExpectIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	var opts []nodes.Node
	if !p.auditIsWord(p.cur, "TO") {
		return nil, p.syntaxErrorAtCur()
	}
	target, err := p.auditParseTarget(true)
	if err != nil {
		return nil, err
	}
	opts = append(opts, target...)

	if p.cur.Type == kwWITH {
		with, err := p.auditParseWithOptions(auditCreateOptions)
		if err != nil {
			return nil, err
		}
		opts = append(opts, with...)
	}
	if p.cur.Type == kwWHERE {
		p.advance()
		where, err := p.auditParsePredicateOr()
		if err != nil {
			return nil, err
		}
		stmt.WhereClause = where
	}
	stmt.Options = &nodes.List{Items: opts}
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterServerAuditStmt parses ALTER SERVER AUDIT; the caller has
// consumed ALTER SERVER AUDIT.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-server-audit-transact-sql
// Grammar: TSql160.g:2194 alterServerAuditStatement. Its action (TSql160.g:2231)
// refuses a statement with no target, option or predicate.
//
//	ALTER SERVER AUDIT audit_name
//	{   [ TO { FILE ( <file_option> [ ,...n ] ) | APPLICATION_LOG | SECURITY_LOG } ]
//	    [ WITH ( <audit_option> [ ,...n ] ) ]
//	    [ WHERE <predicate_expression> ]
//	  | REMOVE WHERE
//	  | MODIFY NAME = new_audit_name }
//
//	<audit_option> ::= QUEUE_DELAY = integer
//	    | ON_FAILURE = { CONTINUE | SHUTDOWN | FAIL_OPERATION }
//	    | STATE = { ON | OFF }
func (p *Parser) parseAlterServerAuditStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "ALTER",
		ObjectType: "SERVER AUDIT",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}
	name, err := p.auditExpectIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	var opts []nodes.Node
	switch {
	case p.cur.Type == kwREMOVE:
		p.advance()
		if _, err := p.expect(kwWHERE); err != nil {
			return nil, err
		}
		opts = append(opts, &nodes.String{Str: "REMOVE WHERE"})
	case p.cur.Type == kwMODIFY:
		p.advance()
		if !p.auditIsWord(p.cur, "NAME") {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		newName, err := p.auditExpectIdentifier()
		if err != nil {
			return nil, err
		}
		opts = append(opts, &nodes.String{Str: "MODIFY NAME=" + newName})
	default:
		if p.auditIsWord(p.cur, "TO") {
			target, err := p.auditParseTarget(false)
			if err != nil {
				return nil, err
			}
			opts = append(opts, target...)
		}
		if p.cur.Type == kwWITH {
			with, err := p.auditParseWithOptions(auditAlterOptions)
			if err != nil {
				return nil, err
			}
			opts = append(opts, with...)
		}
		if p.cur.Type == kwWHERE {
			p.advance()
			where, err := p.auditParsePredicateOr()
			if err != nil {
				return nil, err
			}
			stmt.WhereClause = where
		}
		if len(opts) == 0 && stmt.WhereClause == nil {
			return nil, p.syntaxErrorAtCur()
		}
	}
	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseDropServerAuditStmt parses DROP SERVER AUDIT; the caller has consumed
// DROP SERVER AUDIT.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-server-audit-transact-sql
// Grammar: TSql160.g:14859 dropServerAuditStatement
//
//	DROP SERVER AUDIT audit_name
func (p *Parser) parseDropServerAuditStmt() (*nodes.SecurityStmt, error) {
	return p.auditParseDrop("SERVER AUDIT")
}

// parseDropAuditStmt parses DROP AUDIT; the caller has consumed DROP AUDIT.
// No document or ScriptDOM rule names it, but SQL Server 2022 accepts
// DROP AUDIT followed by one name and refuses a string or a list.
//
//	DROP AUDIT audit_name
func (p *Parser) parseDropAuditStmt() (*nodes.SecurityStmt, error) {
	return p.auditParseDrop("AUDIT")
}

// parseCreateServerAuditSpecStmt parses CREATE SERVER AUDIT SPECIFICATION; the
// caller has consumed CREATE SERVER AUDIT SPECIFICATION.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-server-audit-specification-transact-sql
// Grammar: TSql160.g:2109 createServerAuditSpecificationStatement
//
//	CREATE SERVER AUDIT SPECIFICATION audit_specification_name
//	    FOR SERVER AUDIT audit_name
//	    [ ADD ( audit_action_group_name ) [ ,...n ] ]
//	    [ WITH ( STATE = { ON | OFF } [ ,...n ] ) ]
func (p *Parser) parseCreateServerAuditSpecStmt() (*nodes.SecurityStmt, error) {
	return p.auditParseSpec("CREATE", "SERVER AUDIT SPECIFICATION", false)
}

// parseAlterServerAuditSpecStmt parses ALTER SERVER AUDIT SPECIFICATION; the
// caller has consumed ALTER SERVER AUDIT SPECIFICATION.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-server-audit-specification-transact-sql
// Grammar: TSql160.g:2167 alterServerAuditSpecificationStatement
//
//	ALTER SERVER AUDIT SPECIFICATION audit_specification_name
//	    [ FOR SERVER AUDIT audit_name ]
//	    [ { ADD | DROP } ( audit_action_group_name ) [ ,...n ] ]
//	    [ WITH ( STATE = { ON | OFF } [ ,...n ] ) ]
func (p *Parser) parseAlterServerAuditSpecStmt() (*nodes.SecurityStmt, error) {
	return p.auditParseSpec("ALTER", "SERVER AUDIT SPECIFICATION", false)
}

// parseDropServerAuditSpecStmt parses DROP SERVER AUDIT SPECIFICATION; the
// caller has consumed DROP SERVER AUDIT SPECIFICATION.
//
// Grammar: TSql160.g:14848 dropServerAuditSpecificationStatement
//
//	DROP SERVER AUDIT SPECIFICATION audit_specification_name
func (p *Parser) parseDropServerAuditSpecStmt() (*nodes.SecurityStmt, error) {
	return p.auditParseDrop("SERVER AUDIT SPECIFICATION")
}

// parseCreateDatabaseAuditSpecStmt parses CREATE DATABASE AUDIT SPECIFICATION;
// the caller has consumed CREATE DATABASE AUDIT SPECIFICATION.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-database-audit-specification-transact-sql
// Grammar: TSql160.g:5267 createDatabaseAuditSpecification
//
//	CREATE DATABASE AUDIT SPECIFICATION audit_specification_name
//	    FOR SERVER AUDIT audit_name
//	    [ ADD ( { <audit_action_specification> | audit_action_group_name } ) [ ,...n ] ]
//	    [ WITH ( STATE = { ON | OFF } [ ,...n ] ) ]
func (p *Parser) parseCreateDatabaseAuditSpecStmt() (*nodes.SecurityStmt, error) {
	return p.auditParseSpec("CREATE", "DATABASE AUDIT SPECIFICATION", true)
}

// parseAlterDatabaseAuditSpecStmt parses ALTER DATABASE AUDIT SPECIFICATION;
// the caller has consumed ALTER DATABASE AUDIT SPECIFICATION.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-database-audit-specification-transact-sql
// Grammar: TSql160.g:4206 alterDatabaseAuditSpecification
//
//	ALTER DATABASE AUDIT SPECIFICATION audit_specification_name
//	    [ FOR SERVER AUDIT audit_name ]
//	    [ { ADD | DROP } ( { <audit_action_specification> | audit_action_group_name } ) [ ,...n ] ]
//	    [ WITH ( STATE = { ON | OFF } [ ,...n ] ) ]
func (p *Parser) parseAlterDatabaseAuditSpecStmt() (*nodes.SecurityStmt, error) {
	return p.auditParseSpec("ALTER", "DATABASE AUDIT SPECIFICATION", true)
}

// parseDropDatabaseAuditSpecStmt parses DROP DATABASE AUDIT SPECIFICATION; the
// caller has consumed DROP DATABASE AUDIT SPECIFICATION.
//
// Grammar: TSql160.g:14906 dropDatabaseAuditSpecificationStatement
//
//	DROP DATABASE AUDIT SPECIFICATION audit_specification_name
func (p *Parser) parseDropDatabaseAuditSpecStmt() (*nodes.SecurityStmt, error) {
	return p.auditParseDrop("DATABASE AUDIT SPECIFICATION")
}

// auditParseDrop reads the one name a DROP of an audit object takes.
func (p *Parser) auditParseDrop(objectType string) (*nodes.SecurityStmt, error) {
	stmt := &nodes.SecurityStmt{
		Action:     "DROP",
		ObjectType: objectType,
		Loc:        nodes.Loc{Start: p.pos(), End: -1},
	}
	name, err := p.auditExpectIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// auditIsWord reports whether tok is the word w written bare, as a keyword or
// an identifier: SQL Server refuses an option word in brackets or quotes
// ([MAXSIZE] = 1 MB), as ScriptDOM's Match on the identifier's text does.
// w is upper-case; a word the lexer knows is matched by its token type.
func (p *Parser) auditIsWord(tok Token, w string) bool {
	bare := tok.Type == tokIDENT && !p.isDelimited(tok)
	if kw := lookupKeyword(w); kw != tokIDENT {
		return tok.Type == kw || (bare && lookupKeyword(tok.Str) == kw)
	}
	return bare && strings.ToUpper(tok.Str) == w
}

// auditBareWord returns the upper-case text of tok when it is a bare word (a
// keyword or an undelimited identifier), and "" otherwise.
func (p *Parser) auditBareWord(tok Token) string {
	if (tok.Type == tokIDENT && !p.isDelimited(tok)) || (tok.Type >= kwABSENT && tok.Str != "") {
		return strings.ToUpper(tok.Str)
	}
	return ""
}

// auditExpectIdentifier consumes an identifier (TSql160.g identifier: a
// regular, bracketed or quoted name, never a string).
func (p *Parser) auditExpectIdentifier() (string, error) {
	name, ok := p.parseIdentifier()
	if !ok {
		return "", p.syntaxErrorAtCur()
	}
	return name, nil
}

// auditExpectInt32 consumes an integer literal that fits an int: SQL Server
// refuses a decimal, a real, a hex literal, a sign, and a value past
// 2,147,483,647 (Msg 102 or 1080) wherever the audit grammar says integer.
func (p *Parser) auditExpectInt32() (string, error) {
	tok := p.cur
	if tok.Type != tokICONST || strings.HasPrefix(tok.Str, "0x") || strings.HasPrefix(tok.Str, "0X") || tok.Ival > 2147483647 {
		return "", p.syntaxErrorAtCur()
	}
	p.advance()
	return tok.Str, nil
}

// auditExpectString consumes a string literal, N'...' included.
func (p *Parser) auditExpectString() (string, error) {
	if p.cur.Type != tokSCONST && p.cur.Type != tokNSCONST {
		return "", p.syntaxErrorAtCur()
	}
	return p.advance().Str, nil
}

// auditExpectOnOff consumes ON or OFF (TSql160.g optionOnOff).
func (p *Parser) auditExpectOnOff() (string, error) {
	if p.cur.Type != kwON && p.cur.Type != kwOFF {
		return "", p.syntaxErrorAtCur()
	}
	return strings.ToUpper(p.advance().Str), nil
}

// auditExpectOptionName consumes the option word and the = after it, and
// returns the word upper-cased.
func (p *Parser) auditExpectOptionName(names ...string) (string, error) {
	word := p.auditBareWord(p.cur)
	for _, n := range names {
		if word == n {
			p.advance()
			if _, err := p.expect('='); err != nil {
				return "", err
			}
			return word, nil
		}
	}
	return "", p.syntaxErrorAtCur()
}

// auditNotSupported is SQL Server 2022's refusal of an Azure-only audit word
// (Msg 40517).
func (p *Parser) auditNotSupported() error {
	return auditError(p.cur.Loc, fmt.Sprintf("%s is not supported in this version of SQL Server", strings.ToUpper(p.cur.Str)))
}

// auditParseTarget parses the TO clause; the current token is TO. CREATE
// requires FILEPATH among the file options (pathRequired), ALTER does not.
//
// Grammar: TSql160.g:1717 auditTargetClause, TSql160.g:1798 auditFileOption
//
//	TO { FILE ( <file_option> [ ,...n ] ) | APPLICATION_LOG | SECURITY_LOG }
//
//	<file_option> ::= FILEPATH = 'os_file_path'
//	    | MAXSIZE = { integer { MB | GB | TB } | UNLIMITED }
//	    | MAX_ROLLOVER_FILES = { integer | UNLIMITED }
//	    | MAX_FILES = integer
//	    | RESERVE_DISK_SPACE = { ON | OFF }
//
// SQL Server 2022 refuses URL and EXTERNAL_MONITOR, which ScriptDOM takes for
// Azure SQL Managed Instance, and FILE's PATH, which only URL takes.
func (p *Parser) auditParseTarget(pathRequired bool) ([]nodes.Node, error) {
	p.advance() // TO
	switch word := p.auditBareWord(p.cur); word {
	case "APPLICATION_LOG", "SECURITY_LOG":
		p.advance()
		return []nodes.Node{&nodes.String{Str: "TO=" + word}}, nil
	case "URL", "EXTERNAL_MONITOR":
		return nil, p.auditNotSupported()
	case "FILE":
	default:
		return nil, p.syntaxErrorAtCur()
	}
	fileLoc := p.pos()
	p.advance() // FILE
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	opts := []nodes.Node{&nodes.String{Str: "TO=FILE"}}
	sawPath := false
	for {
		name, err := p.auditExpectOptionName("FILEPATH", "MAXSIZE", "MAX_ROLLOVER_FILES", "MAX_FILES", "RESERVE_DISK_SPACE")
		if err != nil {
			return nil, err
		}
		var val string
		switch name {
		case "FILEPATH":
			sawPath = true
			val, err = p.auditExpectString()
		case "MAXSIZE":
			if p.auditIsWord(p.cur, "UNLIMITED") {
				p.advance()
				val = "UNLIMITED"
				break
			}
			if val, err = p.auditExpectInt32(); err != nil {
				return nil, err
			}
			unit := p.auditBareWord(p.cur)
			if unit != "MB" && unit != "GB" && unit != "TB" {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
			val += unit
		case "MAX_ROLLOVER_FILES":
			if p.auditIsWord(p.cur, "UNLIMITED") {
				p.advance()
				val = "UNLIMITED"
				break
			}
			val, err = p.auditExpectInt32()
		case "MAX_FILES":
			val, err = p.auditExpectInt32()
		case "RESERVE_DISK_SPACE":
			val, err = p.auditExpectOnOff()
		}
		if err != nil {
			return nil, err
		}
		opts = append(opts, &nodes.String{Str: name + "=" + val})
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	// TSql160.g:1762: CREATE needs FILEPATH (SQL46056).
	if pathRequired && !sawPath {
		return nil, auditError(fileLoc, "TO FILE requires a FILEPATH option")
	}
	return opts, nil
}

// auditCreateOptions and auditAlterOptions are the options of CREATE and
// ALTER SERVER AUDIT's WITH: TSql160.g:1996 auditCreateOption and
// TSql160.g:2011 auditOption, less OPERATOR_AUDIT, which SQL Server 2022
// refuses (Msg 40517).
var (
	auditCreateOptions = []string{"QUEUE_DELAY", "ON_FAILURE", "AUDIT_GUID"}
	auditAlterOptions  = []string{"QUEUE_DELAY", "ON_FAILURE", "STATE"}
)

// auditParseWithOptions parses WITH ( <audit_option> [ ,...n ] ); the current
// token is WITH.
//
// Grammar: TSql160.g:1957 auditCreateWithClause, TSql160.g:1976 auditWithClause
//
//	QUEUE_DELAY = integer
//	ON_FAILURE = { CONTINUE | SHUTDOWN | FAIL_OPERATION }
//	AUDIT_GUID = 'uniqueidentifier'  (CREATE)
//	STATE = { ON | OFF }             (ALTER)
func (p *Parser) auditParseWithOptions(allowed []string) ([]nodes.Node, error) {
	p.advance() // WITH
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	var opts []nodes.Node
	for {
		if p.auditIsWord(p.cur, "OPERATOR_AUDIT") {
			return nil, p.auditNotSupported()
		}
		name, err := p.auditExpectOptionName(allowed...)
		if err != nil {
			return nil, err
		}
		var val string
		switch name {
		case "QUEUE_DELAY":
			val, err = p.auditExpectInt32()
		case "ON_FAILURE":
			val = p.auditBareWord(p.cur)
			if val != "CONTINUE" && val != "SHUTDOWN" && val != "FAIL_OPERATION" {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
		case "AUDIT_GUID":
			guidLoc := p.pos()
			if val, err = p.auditExpectString(); err != nil {
				return nil, err
			}
			if !auditIsGUIDText(val) {
				// Msg 8169: the string does not convert to uniqueidentifier.
				return nil, auditError(guidLoc, "AUDIT_GUID is not a uniqueidentifier")
			}
		case "STATE":
			val, err = p.auditExpectOnOff()
		}
		if err != nil {
			return nil, err
		}
		opts = append(opts, &nodes.String{Str: name + "=" + val})
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return opts, nil
}

// auditIsGUIDText reports whether SQL Server converts s to uniqueidentifier:
// 8-4-4-4-12 hex digits, optionally after '{'. SQL Server ignores what
// follows the 36th character, so '{...}' and a trailing blank convert too.
func auditIsGUIDText(s string) bool {
	s = strings.TrimPrefix(s, "{")
	if len(s) < 36 {
		return false
	}
	for i := 0; i < 36; i++ {
		c := s[i]
		switch i {
		case 8, 13, 18, 23:
			if c != '-' {
				return false
			}
		default:
			if !isHexDigit(c) {
				return false
			}
		}
	}
	return true
}

// auditServerActionGroups and auditDatabaseActionGroups are the action groups
// SQL Server 2022 accepts in a server and a database audit specification:
// the groups sys.dm_audit_actions lists with class SERVER and DATABASE, each
// checked against the server's parser. ScriptDOM reads the group as an
// identifier (TSql160.g:5393, 5402) and checks it against a list that lacks
// the newer groups (LEDGER_OPERATION_GROUP, EXTGOV_OPERATION_GROUP, ...).
//
// Ref: https://learn.microsoft.com/en-us/sql/relational-databases/security/auditing/sql-server-audit-action-groups-and-actions
var (
	auditDatabaseActionGroups = auditGroupSet(
		"APPLICATION_ROLE_CHANGE_PASSWORD_GROUP", "AUDIT_CHANGE_GROUP", "BACKUP_RESTORE_GROUP",
		"BATCH_COMPLETED_GROUP", "BATCH_STARTED_GROUP", "DATABASE_CHANGE_GROUP", "DATABASE_LOGOUT_GROUP",
		"DATABASE_OBJECT_ACCESS_GROUP", "DATABASE_OBJECT_CHANGE_GROUP", "DATABASE_OBJECT_OWNERSHIP_CHANGE_GROUP",
		"DATABASE_OBJECT_PERMISSION_CHANGE_GROUP", "DATABASE_OPERATION_GROUP", "DATABASE_OWNERSHIP_CHANGE_GROUP",
		"DATABASE_PERMISSION_CHANGE_GROUP", "DATABASE_PRINCIPAL_CHANGE_GROUP", "DATABASE_PRINCIPAL_IMPERSONATION_GROUP",
		"DATABASE_ROLE_MEMBER_CHANGE_GROUP", "DBCC_GROUP", "EXTGOV_OPERATION_GROUP",
		"FAILED_DATABASE_AUTHENTICATION_GROUP", "LEDGER_OPERATION_GROUP", "SCHEMA_OBJECT_ACCESS_GROUP",
		"SCHEMA_OBJECT_CHANGE_GROUP", "SCHEMA_OBJECT_OWNERSHIP_CHANGE_GROUP", "SCHEMA_OBJECT_PERMISSION_CHANGE_GROUP",
		"SENSITIVE_BATCH_COMPLETED_GROUP", "SENSITIVITY_CLASSIFICATION_CHANGE_GROUP", "STATEMENT_ROLLBACK_GROUP",
		"SUCCESSFUL_DATABASE_AUTHENTICATION_GROUP", "TRANSACTION_BEGIN_GROUP", "TRANSACTION_COMMIT_GROUP",
		"TRANSACTION_GROUP", "TRANSACTION_ROLLBACK_GROUP", "USER_CHANGE_PASSWORD_GROUP", "USER_DEFINED_AUDIT_GROUP",
	)
	auditServerActionGroups = auditGroupSet(
		"BROKER_LOGIN_GROUP", "DATABASE_MIRRORING_LOGIN_GROUP", "FAILED_LOGIN_GROUP", "FULLTEXT_GROUP",
		"GLOBAL_TRANSACTIONS_LOGIN_GROUP", "LOGIN_CHANGE_PASSWORD_GROUP", "LOGOUT_GROUP",
		"SERVER_OBJECT_CHANGE_GROUP", "SERVER_OBJECT_OWNERSHIP_CHANGE_GROUP", "SERVER_OBJECT_PERMISSION_CHANGE_GROUP",
		"SERVER_OPERATION_GROUP", "SERVER_PERMISSION_CHANGE_GROUP", "SERVER_PRINCIPAL_CHANGE_GROUP",
		"SERVER_PRINCIPAL_IMPERSONATION_GROUP", "SERVER_ROLE_MEMBER_CHANGE_GROUP", "SERVER_STATE_CHANGE_GROUP",
		"STORAGE_LOGIN_GROUP", "SUCCESSFUL_LOGIN_GROUP", "TRACE_CHANGE_GROUP",
	).union(auditDatabaseActionGroups)
)

type auditGroups map[string]bool

func auditGroupSet(names ...string) auditGroups {
	s := make(auditGroups, len(names))
	for _, n := range names {
		s[n] = true
	}
	return s
}

func (s auditGroups) union(o auditGroups) auditGroups {
	out := make(auditGroups, len(s)+len(o))
	for k := range s {
		out[k] = true
	}
	for k := range o {
		out[k] = true
	}
	return out
}

// auditParseSpec parses the body of a server or database audit specification
// after its keywords. CREATE requires FOR SERVER AUDIT and takes only ADD.
//
// Grammar: TSql160.g:5295 auditSpecificationForClause, TSql160.g:5491
// auditSpecificationStateOpt. SQL Server takes a list of STATE options where
// the grammar takes one.
func (p *Parser) auditParseSpec(action, objectType string, database bool) (*nodes.SecurityStmt, error) {
	stmt := &nodes.SecurityStmt{
		Action:     action,
		ObjectType: objectType,
		Loc:        nodes.Loc{Start: p.pos(), End: -1},
	}
	name, err := p.auditExpectIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	create := action == "CREATE"

	var opts []nodes.Node
	// FOR SERVER AUDIT audit_name
	if p.cur.Type == kwFOR {
		p.advance()
		if !p.auditIsWord(p.cur, "SERVER") {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		if !p.auditIsWord(p.cur, "AUDIT") {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		auditName, err := p.auditExpectIdentifier()
		if err != nil {
			return nil, err
		}
		opts = append(opts, &nodes.String{Str: "FOR_AUDIT=" + auditName})
	} else if create {
		return nil, p.syntaxErrorAtCur()
	}

	// { ADD | DROP } ( ... ) [ ,...n ]
	if p.cur.Type == kwADD || (!create && p.cur.Type == kwDROP) {
		for {
			if p.cur.Type != kwADD && (create || p.cur.Type != kwDROP) {
				return nil, p.syntaxErrorAtCur()
			}
			detail, err := p.auditParseSpecDetail(database)
			if err != nil {
				return nil, err
			}
			opts = append(opts, detail)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}

	// WITH ( STATE = { ON | OFF } [ ,...n ] )
	if p.cur.Type == kwWITH {
		p.advance()
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		for {
			if _, err := p.auditExpectOptionName("STATE"); err != nil {
				return nil, err
			}
			val, err := p.auditExpectOnOff()
			if err != nil {
				return nil, err
			}
			opts = append(opts, &nodes.String{Str: "STATE=" + val})
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// auditActionWords are the actions of an audit action specification.
//
// Grammar: TSql160.g:5445 actionWithQual
var auditActionWords = map[string]bool{
	"SELECT": true, "UPDATE": true, "INSERT": true, "DELETE": true,
	"EXECUTE": true, "RECEIVE": true, "REFERENCES": true,
}

// auditParseSpecDetail parses one ADD or DROP of a specification; the current
// token is ADD or DROP.
//
// Grammar: TSql160.g:5307 createAuditSpecificationDetailDb, TSql160.g:5326
// createAuditSpecificationDetail, TSql160.g:5411 auditActionSpecification
//
//	{ ADD | DROP } ( audit_action_group_name )
//	{ ADD | DROP } ( <audit_action_specification> )   (database only)
//
//	<audit_action_specification> ::=
//	    action [ ,...n ] ON [ { OBJECT | SCHEMA | DATABASE } :: ] securable
//	    BY principal [ ,...n ]
func (p *Parser) auditParseSpecDetail(database bool) (*nodes.AuditSpecAction, error) {
	loc := p.pos()
	node := &nodes.AuditSpecAction{
		Action: strings.ToUpper(p.advance().Str),
		Loc:    nodes.Loc{Start: loc, End: -1},
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	word := p.auditBareWord(p.cur)
	if database && auditActionWords[word] {
		for {
			word := p.auditBareWord(p.cur)
			if !auditActionWords[word] {
				return nil, p.syntaxErrorAtCur()
			}
			node.Actions = append(node.Actions, word)
			p.advance()
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if err := p.auditParseSecurable(node); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		for {
			// TSql160.g:13986 principal: PUBLIC | NULL | identifier.
			switch {
			case p.cur.Type == kwPUBLIC || p.cur.Type == kwNULL:
				node.Principals = append(node.Principals, p.advance().Str)
			default:
				principal, err := p.auditExpectIdentifier()
				if err != nil {
					return nil, err
				}
				node.Principals = append(node.Principals, principal)
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	} else {
		groups := auditServerActionGroups
		if database {
			groups = auditDatabaseActionGroups
		}
		if !groups[word] {
			return nil, p.syntaxErrorAtCur()
		}
		node.GroupName = word
		p.advance()
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	node.Loc.End = p.prevEnd()
	return node, nil
}

// auditParseSecurable parses ON [ class :: ] securable of an audit action
// specification.
//
// Grammar: TSql160.g:13924 securityTargetObjectCommon. SQL Server takes only
// the classes OBJECT, SCHEMA and DATABASE here (Msg 33231), a name of up to
// four parts, empty middle parts included (a..b), and no column list.
func (p *Parser) auditParseSecurable(node *nodes.AuditSpecAction) error {
	if _, err := p.expect(kwON); err != nil {
		return err
	}
	if p.peekNext().Type == tokCOLONCOLON {
		class := p.auditBareWord(p.cur)
		if class != "OBJECT" && class != "SCHEMA" && class != "DATABASE" {
			return p.syntaxErrorAtCur()
		}
		node.ClassName = class
		p.advance() // class
		p.advance() // ::
	}
	first, err := p.auditExpectIdentifier()
	if err != nil {
		return err
	}
	parts := []string{first}
	for p.cur.Type == '.' {
		if len(parts) == 4 {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		if p.cur.Type == '.' {
			parts = append(parts, "")
			continue
		}
		part, err := p.auditExpectIdentifier()
		if err != nil {
			return err
		}
		parts = append(parts, part)
	}
	node.Securable = strings.Join(parts, ".")
	return nil
}

// auditParsePredicateOr parses a server audit's WHERE predicate, ScriptDOM's
// event predicate.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-server-audit-transact-sql
// Grammar: TSql160.g:8142 eventBooleanExpression, TSql160.g:8147 eventBooleanExpressionOr
//
//	<predicate> ::= <and> [ OR <and> ]...
func (p *Parser) auditParsePredicateOr() (nodes.ExprNode, error) {
	left, err := p.auditParsePredicateAnd()
	if err != nil {
		return nil, err
	}
	for p.cur.Type == kwOR {
		p.advance()
		right, err := p.auditParsePredicateAnd()
		if err != nil {
			return nil, err
		}
		left = &nodes.BinaryExpr{
			Op: nodes.BinOpOr, Left: left, Right: right,
			Loc: nodes.Loc{Start: auditExprStart(left), End: p.prevEnd()},
		}
	}
	return left, nil
}

// auditParsePredicateAnd parses <unary> [ AND <unary> ]...
//
// Grammar: TSql160.g:8161 eventBooleanExpressionAnd
func (p *Parser) auditParsePredicateAnd() (nodes.ExprNode, error) {
	left, err := p.auditParsePredicateUnary()
	if err != nil {
		return nil, err
	}
	for p.cur.Type == kwAND {
		p.advance()
		right, err := p.auditParsePredicateUnary()
		if err != nil {
			return nil, err
		}
		left = &nodes.BinaryExpr{
			Op: nodes.BinOpAnd, Left: left, Right: right,
			Loc: nodes.Loc{Start: auditExprStart(left), End: p.prevEnd()},
		}
	}
	return left, nil
}

// auditParsePredicateUnary parses NOT <unary>, ( <predicate> ), or a leaf.
//
// Grammar: TSql160.g:8175 eventBooleanExpressionUnary, TSql160.g:8193
// eventBooleanExpressionParenthesis
func (p *Parser) auditParsePredicateUnary() (nodes.ExprNode, error) {
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
	loc := p.pos()
	switch p.cur.Type {
	case kwNOT:
		p.advance()
		operand, err := p.auditParsePredicateUnary()
		if err != nil {
			return nil, err
		}
		return &nodes.UnaryExpr{Op: nodes.UnaryNot, Operand: operand, Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil
	case '(':
		p.advance()
		inner, err := p.auditParsePredicateOr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return &nodes.ParenExpr{Expr: inner, Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil
	}
	return p.auditParsePredicateLeaf()
}

// auditParsePredicateLeaf parses a comparison of an event field, or a call of
// a predicate compare function.
//
// Grammar: TSql160.g:8206 eventDeclarationPredicateParameter, TSql160.g:8246
// eventDeclarationCompareFunction, TSql160.g:8260 eventDeclarationComparisonPredicate.
// The function's name has two or three parts (TSql160.g:8220).
//
//	field { = | <> | != | > | >= | < | <= | !< | !> | LIKE | NOT LIKE } <value>
//	package.compare_function ( field , <value> )
func (p *Parser) auditParsePredicateLeaf() (nodes.ExprNode, error) {
	loc := p.pos()
	parts, err := p.auditParseEventName()
	if err != nil {
		return nil, err
	}
	nameEnd := p.prevEnd()
	if p.cur.Type == '(' {
		if len(parts) < 2 {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		argLoc := p.pos()
		argParts, err := p.auditParseEventName()
		if err != nil {
			return nil, err
		}
		arg := auditColumnRef(argParts, argLoc, p.prevEnd())
		if _, err := p.expect(','); err != nil {
			return nil, err
		}
		val, err := p.auditParsePredicateValue()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		name := &nodes.TableRef{Object: parts[len(parts)-1], Schema: parts[len(parts)-2], Loc: nodes.Loc{Start: loc, End: nameEnd}}
		if len(parts) == 3 {
			name.Database = parts[0]
		}
		return &nodes.FuncCallExpr{
			Name: name,
			Args: &nodes.List{Items: []nodes.Node{arg, val}},
			Loc:  nodes.Loc{Start: loc, End: p.prevEnd()},
		}, nil
	}
	field := auditColumnRef(parts, loc, nameEnd)

	if p.cur.Type == kwLIKE || p.cur.Type == kwNOT {
		not := p.cur.Type == kwNOT
		p.advance()
		if not {
			if _, err := p.expect(kwLIKE); err != nil {
				return nil, err
			}
		}
		val, err := p.auditParsePredicateValue()
		if err != nil {
			return nil, err
		}
		return &nodes.LikeExpr{Expr: field, Pattern: val, Not: not, Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil
	}
	op, err := p.auditParseComparisonOp()
	if err != nil {
		return nil, err
	}
	val, err := p.auditParsePredicateValue()
	if err != nil {
		return nil, err
	}
	return &nodes.BinaryExpr{Op: op, Left: field, Right: val, Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil
}

// auditParseEventName parses a name of one to three parts.
//
// Grammar: TSql160.g:30475 oneOrTwoOrThreePartObjectName
func (p *Parser) auditParseEventName() ([]string, error) {
	first, err := p.auditExpectIdentifier()
	if err != nil {
		return nil, err
	}
	parts := []string{first}
	for p.cur.Type == '.' && len(parts) < 3 {
		p.advance()
		part, err := p.auditExpectIdentifier()
		if err != nil {
			return nil, err
		}
		parts = append(parts, part)
	}
	return parts, nil
}

// auditColumnRef makes the column reference of an event field's name.
func auditColumnRef(parts []string, start, end int) *nodes.ColumnRef {
	ref := &nodes.ColumnRef{Column: parts[len(parts)-1], Loc: nodes.Loc{Start: start, End: end}}
	if len(parts) >= 2 {
		ref.Table = parts[len(parts)-2]
	}
	if len(parts) == 3 {
		ref.Schema = parts[0]
	}
	return ref
}

// auditParseComparisonOp parses a comparison operator, which may be written
// as two tokens (< >, ! =).
//
// Grammar: TSql160.g:31264 comparisonOperator
//
//	= | > [ = ] | < [ = | > ] | ! { = | < | > }
func (p *Parser) auditParseComparisonOp() (nodes.BinaryOp, error) {
	if op, ok := comparisonOp(p.cur.Type); ok {
		first := p.advance().Type
		switch {
		case first == '<' && p.cur.Type == '=':
			p.advance()
			return nodes.BinOpLte, nil
		case first == '<' && p.cur.Type == '>':
			p.advance()
			return nodes.BinOpNeq, nil
		case first == '>' && p.cur.Type == '=':
			p.advance()
			return nodes.BinOpGte, nil
		}
		return op, nil
	}
	if p.cur.Type == '!' {
		p.advance()
		var op nodes.BinaryOp
		switch p.cur.Type {
		case '=':
			op = nodes.BinOpNeq
		case '<':
			op = nodes.BinOpNotLt
		case '>':
			op = nodes.BinOpNotGt
		default:
			return 0, p.syntaxErrorAtCur()
		}
		p.advance()
		return op, nil
	}
	return 0, p.syntaxErrorAtCur()
}

// auditParsePredicateValue parses the value of an event predicate.
//
// Grammar: TSql160.g:8102 eventDeclarationValue, TSql160.g:8129
// eventSessionExpressionParenthesis
//
//   - { integer | real | numeric } | ( <value> ) | integer | real | numeric | 'string'
//
// SQL Server also takes a binary literal (0x10), but not after a minus.
func (p *Parser) auditParsePredicateValue() (nodes.ExprNode, error) {
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
	loc := p.pos()
	switch p.cur.Type {
	case '(':
		p.advance()
		inner, err := p.auditParsePredicateValue()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return &nodes.ParenExpr{Expr: inner, Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil
	case '-':
		p.advance()
		if p.cur.Type != tokFCONST && (p.cur.Type != tokICONST || strings.HasPrefix(p.cur.Str, "0x") || strings.HasPrefix(p.cur.Str, "0X")) {
			return nil, p.syntaxErrorAtCur()
		}
		num := p.auditLiteral()
		return &nodes.UnaryExpr{Op: nodes.UnaryMinus, Operand: num, Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil
	case tokICONST, tokFCONST, tokSCONST, tokNSCONST:
		return p.auditLiteral(), nil
	}
	return nil, p.syntaxErrorAtCur()
}

// auditLiteral consumes the current number or string token as a literal.
func (p *Parser) auditLiteral() *nodes.Literal {
	tok := p.advance()
	lit := &nodes.Literal{Str: tok.Str, Loc: nodes.Loc{Start: tok.Loc, End: p.prevEnd()}}
	switch tok.Type {
	case tokICONST:
		lit.Type = nodes.LitInteger
		lit.Ival = tok.Ival
	case tokFCONST:
		lit.Type = nodes.LitFloat
	case tokNSCONST:
		lit.Type = nodes.LitString
		lit.IsNChar = true
	default:
		lit.Type = nodes.LitString
	}
	return lit
}

// auditExprStart returns the start offset of an audit predicate node.
func auditExprStart(e nodes.ExprNode) int {
	switch n := e.(type) {
	case *nodes.BinaryExpr:
		return n.Loc.Start
	case *nodes.UnaryExpr:
		return n.Loc.Start
	case *nodes.ParenExpr:
		return n.Loc.Start
	case *nodes.LikeExpr:
		return n.Loc.Start
	case *nodes.FuncCallExpr:
		return n.Loc.Start
	}
	return -1
}

// auditError is a ParseError that reports its own message rather than the
// generic "syntax error at or near" text, which needs Source unset.
func auditError(pos int, msg string) *ParseError {
	return &ParseError{Message: msg, Position: pos}
}
