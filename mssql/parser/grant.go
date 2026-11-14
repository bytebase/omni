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

	if err := p.parseGrantPermissionsOnTo(stmt, kwTO); err != nil {
		return nil, err
	}

	// WITH GRANT OPTION
	if p.cur.Type == kwWITH {
		p.advance() // WITH
		if _, err := p.expect(kwGRANT); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwOPTION); err != nil {
			return nil, err
		}
		stmt.WithGrant = true
	}

	if err := p.parseGrantAsPrincipal(stmt); err != nil {
		return nil, err
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
			if _, err := p.expect(kwFOR); err != nil {
				return nil, err
			}
			stmt.GrantOptionFor = true
		}
	}

	if err := p.parseGrantPermissionsOnTo(stmt, kwFROM); err != nil {
		return nil, err
	}

	// CASCADE
	if _, ok := p.match(kwCASCADE); ok {
		stmt.CascadeOpt = true
	}

	if err := p.parseGrantAsPrincipal(stmt); err != nil {
		return nil, err
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

	if err := p.parseGrantPermissionsOnTo(stmt, kwTO); err != nil {
		return nil, err
	}

	// CASCADE
	if _, ok := p.match(kwCASCADE); ok {
		stmt.CascadeOpt = true
	}

	if err := p.parseGrantAsPrincipal(stmt); err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseGrantOnClause parses ON [ class :: ] securable.
// It detects the class :: pattern by checking for tokCOLONCOLON after the first identifier.
// Securable classes can be multi-word like "XML SCHEMA COLLECTION", "APPLICATION ROLE", etc.
func (p *Parser) parseGrantOnClause(stmt *nodes.GrantStmt) error {
	// Try to detect class :: securable pattern.
	// Securable classes include Core keywords (SCHEMA, DATABASE, etc.)
	// so we must accept any keyword token, not just identifiers.
	if p.isKeywordOrIdent() {
		if p.peekNext().Type == tokCOLONCOLON {
			// Single-word class: SCHEMA, OBJECT, DATABASE, LOGIN, USER, ROLE, etc.
			stmt.OnType = strings.ToUpper(p.cur.Str)
			p.advance() // consume class
			p.advance() // consume ::
		} else if p.isSecurableClassStart() {
			// Multi-word classes: APPLICATION ROLE, ASYMMETRIC KEY, ...
			className := p.tryParseSecurableClass()
			if className == "" {
				return p.syntaxErrorAtCur()
			}
			stmt.OnType = className
		}
	}
	start := p.pos()
	ref, err := p.parseTableRef()
	if err != nil {
		return err
	}
	// A securable has at most three parts (Msg 117).
	if ref == nil || ref.Server != "" {
		return &ParseError{Message: "syntax error in the securable name", Position: start}
	}
	stmt.OnName = ref
	// An object's columns: ON object ( column [ , ...n ] ).
	if p.cur.Type == '(' {
		cols, err := p.parseGrantColumnList()
		if err != nil {
			return err
		}
		for i, item := range stmt.Privileges.Items {
			if s, ok := item.(*nodes.String); ok && !strings.Contains(s.Str, "(") {
				stmt.Privileges.Items[i] = &nodes.String{Str: s.Str + "(" + strings.Join(cols, ", ") + ")"}
			}
		}
	}
	return nil
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

// parseGrantPermissionsOnTo parses the permissions, the optional
// ON [ class :: ] securable, and TO (or, for REVOKE, also FROM) and the
// principals, all of which SQL Server requires but ON.
func (p *Parser) parseGrantPermissionsOnTo(stmt *nodes.GrantStmt, alt int) error {
	privs, needsOn, err := p.parsePrivilegeList()
	if err != nil {
		return err
	}
	stmt.Privileges = privs

	// ON [class ::] securable
	if _, ok := p.match(kwON); ok {
		// Completion: after GRANT ... ON → table_ref
		if p.collectMode() {
			p.addRuleCandidate("table_ref")
			return errCollecting
		}
		if err := p.parseGrantOnClause(stmt); err != nil {
			return err
		}
	} else if needsOn {
		return p.syntaxErrorAtCur()
	}

	// TO principals (REVOKE: TO or FROM)
	if _, ok := p.match(kwTO); !ok {
		if _, ok := p.match(alt); !ok {
			return p.syntaxErrorAtCur()
		}
	}
	// Completion: after GRANT ... ON t TO → user_ref / role_ref
	if p.collectMode() {
		p.addRuleCandidate("user_ref")
		p.addRuleCandidate("role_ref")
		return errCollecting
	}
	stmt.Principals, err = p.parsePrincipalList()
	return err
}

// parseGrantAsPrincipal parses an optional AS principal, a name and not
// PUBLIC (Msg 156).
func (p *Parser) parseGrantAsPrincipal(stmt *nodes.GrantStmt) error {
	if p.cur.Type != kwAS {
		return nil
	}
	p.advance() // consume AS
	name, ok := p.parseIdentifier()
	if !ok {
		return p.syntaxErrorAtCur()
	}
	stmt.AsPrincipal = name
	return nil
}

// grantPermissions holds the permission names SQL Server 2022 takes: every
// permission_name of sys.fn_builtin_permissions(DEFAULT) on the pinned server
// (../testdata/pins.tsv), each checked as GRANT <name> TO PUBLIC, plus ALL
// [ PRIVILEGES ] and EXEC. SQL Server's parser refuses any other name (Msg
// 102), and EXECUTE EXTERNAL SCRIPT even though the function lists it.
var grantPermissions = func() map[string]bool {
	m := map[string]bool{"ALL": true, "ALL PRIVILEGES": true, "EXEC": true}
	for _, name := range []string{
		"ADMINISTER BULK OPERATIONS",
		"ADMINISTER DATABASE BULK OPERATIONS",
		"ALTER",
		"ALTER ANY APPLICATION ROLE",
		"ALTER ANY ASSEMBLY",
		"ALTER ANY ASYMMETRIC KEY",
		"ALTER ANY AVAILABILITY GROUP",
		"ALTER ANY CERTIFICATE",
		"ALTER ANY COLUMN ENCRYPTION KEY",
		"ALTER ANY COLUMN MASTER KEY",
		"ALTER ANY CONNECTION",
		"ALTER ANY CONTRACT",
		"ALTER ANY CREDENTIAL",
		"ALTER ANY DATABASE",
		"ALTER ANY DATABASE AUDIT",
		"ALTER ANY DATABASE DDL TRIGGER",
		"ALTER ANY DATABASE EVENT NOTIFICATION",
		"ALTER ANY DATABASE EVENT SESSION",
		"ALTER ANY DATABASE EVENT SESSION ADD EVENT",
		"ALTER ANY DATABASE EVENT SESSION ADD TARGET",
		"ALTER ANY DATABASE EVENT SESSION DISABLE",
		"ALTER ANY DATABASE EVENT SESSION DROP EVENT",
		"ALTER ANY DATABASE EVENT SESSION DROP TARGET",
		"ALTER ANY DATABASE EVENT SESSION ENABLE",
		"ALTER ANY DATABASE EVENT SESSION OPTION",
		"ALTER ANY DATABASE SCOPED CONFIGURATION",
		"ALTER ANY DATASPACE",
		"ALTER ANY ENDPOINT",
		"ALTER ANY EVENT NOTIFICATION",
		"ALTER ANY EVENT SESSION",
		"ALTER ANY EVENT SESSION ADD EVENT",
		"ALTER ANY EVENT SESSION ADD TARGET",
		"ALTER ANY EVENT SESSION DISABLE",
		"ALTER ANY EVENT SESSION DROP EVENT",
		"ALTER ANY EVENT SESSION DROP TARGET",
		"ALTER ANY EVENT SESSION ENABLE",
		"ALTER ANY EVENT SESSION OPTION",
		"ALTER ANY EXTERNAL DATA SOURCE",
		"ALTER ANY EXTERNAL FILE FORMAT",
		"ALTER ANY EXTERNAL JOB",
		"ALTER ANY EXTERNAL LANGUAGE",
		"ALTER ANY EXTERNAL LIBRARY",
		"ALTER ANY EXTERNAL STREAM",
		"ALTER ANY FULLTEXT CATALOG",
		"ALTER ANY LINKED SERVER",
		"ALTER ANY LOGIN",
		"ALTER ANY MASK",
		"ALTER ANY MESSAGE TYPE",
		"ALTER ANY REMOTE SERVICE BINDING",
		"ALTER ANY ROLE",
		"ALTER ANY ROUTE",
		"ALTER ANY SCHEMA",
		"ALTER ANY SECURITY POLICY",
		"ALTER ANY SENSITIVITY CLASSIFICATION",
		"ALTER ANY SERVER AUDIT",
		"ALTER ANY SERVER ROLE",
		"ALTER ANY SERVICE",
		"ALTER ANY SYMMETRIC KEY",
		"ALTER ANY USER",
		"ALTER LEDGER",
		"ALTER LEDGER CONFIGURATION",
		"ALTER RESOURCES",
		"ALTER SERVER STATE",
		"ALTER SETTINGS",
		"ALTER TRACE",
		"AUTHENTICATE",
		"AUTHENTICATE SERVER",
		"BACKUP DATABASE",
		"BACKUP LOG",
		"CHECKPOINT",
		"CONNECT",
		"CONNECT ANY DATABASE",
		"CONNECT REPLICATION",
		"CONNECT SQL",
		"CONTROL",
		"CONTROL SERVER",
		"CREATE AGGREGATE",
		"CREATE ANY DATABASE",
		"CREATE ANY DATABASE EVENT SESSION",
		"CREATE ANY EVENT SESSION",
		"CREATE ASSEMBLY",
		"CREATE ASYMMETRIC KEY",
		"CREATE AVAILABILITY GROUP",
		"CREATE CERTIFICATE",
		"CREATE CONTRACT",
		"CREATE DATABASE",
		"CREATE DATABASE DDL EVENT NOTIFICATION",
		"CREATE DDL EVENT NOTIFICATION",
		"CREATE DEFAULT",
		"CREATE ENDPOINT",
		"CREATE EXTERNAL LANGUAGE",
		"CREATE EXTERNAL LIBRARY",
		"CREATE FULLTEXT CATALOG",
		"CREATE FUNCTION",
		"CREATE LOGIN",
		"CREATE MESSAGE TYPE",
		"CREATE PROCEDURE",
		"CREATE QUEUE",
		"CREATE REMOTE SERVICE BINDING",
		"CREATE ROLE",
		"CREATE ROUTE",
		"CREATE RULE",
		"CREATE SCHEMA",
		"CREATE SEQUENCE",
		"CREATE SERVER ROLE",
		"CREATE SERVICE",
		"CREATE SYMMETRIC KEY",
		"CREATE SYNONYM",
		"CREATE TABLE",
		"CREATE TRACE EVENT NOTIFICATION",
		"CREATE TYPE",
		"CREATE USER",
		"CREATE VIEW",
		"CREATE XML SCHEMA COLLECTION",
		"DELETE",
		"DROP ANY DATABASE EVENT SESSION",
		"DROP ANY EVENT SESSION",
		"ENABLE LEDGER",
		"EXECUTE",
		"EXECUTE ANY EXTERNAL ENDPOINT",
		"EXECUTE ANY EXTERNAL SCRIPT",
		"EXECUTE EXTERNAL SCRIPT",
		"EXTERNAL ACCESS ASSEMBLY",
		"IMPERSONATE",
		"IMPERSONATE ANY LOGIN",
		"INSERT",
		"KILL DATABASE CONNECTION",
		"RECEIVE",
		"REFERENCES",
		"SELECT",
		"SELECT ALL USER SECURABLES",
		"SEND",
		"SHOWPLAN",
		"SHUTDOWN",
		"SUBSCRIBE QUERY NOTIFICATIONS",
		"TAKE OWNERSHIP",
		"UNMASK",
		"UNSAFE ASSEMBLY",
		"UPDATE",
		"VIEW ANY COLUMN ENCRYPTION KEY DEFINITION",
		"VIEW ANY COLUMN MASTER KEY DEFINITION",
		"VIEW ANY CRYPTOGRAPHICALLY SECURED DEFINITION",
		"VIEW ANY DATABASE",
		"VIEW ANY DEFINITION",
		"VIEW ANY ERROR LOG",
		"VIEW ANY PERFORMANCE DEFINITION",
		"VIEW ANY SECURITY DEFINITION",
		"VIEW ANY SENSITIVITY CLASSIFICATION",
		"VIEW CHANGE TRACKING",
		"VIEW CRYPTOGRAPHICALLY SECURED DEFINITION",
		"VIEW DATABASE PERFORMANCE STATE",
		"VIEW DATABASE SECURITY AUDIT",
		"VIEW DATABASE SECURITY STATE",
		"VIEW DATABASE STATE",
		"VIEW DEFINITION",
		"VIEW LEDGER CONTENT",
		"VIEW PERFORMANCE DEFINITION",
		"VIEW SECURITY DEFINITION",
		"VIEW SERVER PERFORMANCE STATE",
		"VIEW SERVER SECURITY AUDIT",
		"VIEW SERVER SECURITY STATE",
		"VIEW SERVER STATE",
	} {
		if name != "EXECUTE EXTERNAL SCRIPT" {
			m[name] = true
		}
	}
	return m
}()

// grantPermissionsNeedingOn are the permissions SQL Server takes only with
// an ON clause (Msg 102 without one, on the pinned server).
var grantPermissionsNeedingOn = map[string]bool{
	"CREATE SEQUENCE": true, "IMPERSONATE": true, "RECEIVE": true, "SEND": true, "VIEW CHANGE TRACKING": true,
}

// grantPermissionPrefixes holds every proper word prefix of a permission
// name, so that a name is read word by word without backtracking.
var grantPermissionPrefixes = func() map[string]bool {
	m := map[string]bool{}
	for name := range grantPermissions {
		words := strings.Fields(name)
		for i := 1; i < len(words); i++ {
			m[strings.Join(words[:i], " ")] = true
		}
	}
	return m
}()

// grantPermissionWord returns the current token's text, upper-cased, when it
// is a keyword or an undelimited identifier, and "" otherwise.
func (p *Parser) grantPermissionWord() string {
	if p.cur.Str == "" || p.isDelimited(p.cur) || !(p.cur.Type == tokIDENT || p.cur.Type >= kwABSENT) {
		return ""
	}
	return strings.ToUpper(p.cur.Str)
}

// parsePrivilegeList parses a comma-separated list of one or more
// permissions, each a name of grantPermissions and an optional column list.
// It reports whether one of them needs an ON clause.
func (p *Parser) parsePrivilegeList() (*nodes.List, bool, error) {
	var items []nodes.Node
	needsOn := false
	for {
		name := p.grantPermissionWord()
		if name == "" {
			return nil, false, p.syntaxErrorAtCur()
		}
		p.advance()
		for grantPermissionPrefixes[name] {
			next := p.grantPermissionWord()
			if next == "" || !grantPermissions[name+" "+next] && !grantPermissionPrefixes[name+" "+next] {
				break
			}
			name += " " + next
			p.advance()
		}
		if !grantPermissions[name] {
			return nil, false, &ParseError{Message: "unknown permission " + name, Position: p.prev.Loc}
		}
		needsOn = needsOn || grantPermissionsNeedingOn[name]

		// Column list: permission ( column [ , ...n ] )
		if p.cur.Type == '(' {
			cols, err := p.parseGrantColumnList()
			if err != nil {
				return nil, false, err
			}
			name += "(" + strings.Join(cols, ", ") + ")"
		}

		items = append(items, &nodes.String{Str: name})
		if _, ok := p.match(','); !ok {
			break
		}
	}
	return &nodes.List{Items: items}, needsOn, nil
}

// parseGrantColumnList parses ( column [ , ...n ] ), one column at least.
func (p *Parser) parseGrantColumnList() ([]string, error) {
	p.advance() // consume (
	var cols []string
	for {
		col, ok := p.parseIdentifier()
		if !ok {
			return nil, p.syntaxErrorAtCur()
		}
		cols = append(cols, col)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return cols, nil
}

// parsePrincipalList parses a comma-separated list of one or more
// principals, each a name or PUBLIC.
func (p *Parser) parsePrincipalList() (*nodes.List, error) {
	var items []nodes.Node
	for {
		if !p.isIdentLike() && p.cur.Type != kwPUBLIC {
			return nil, p.syntaxErrorAtCur()
		}
		name := p.cur.Str
		p.advance()
		items = append(items, &nodes.String{Str: name})
		if _, ok := p.match(','); !ok {
			break
		}
	}
	return &nodes.List{Items: items}, nil
}
