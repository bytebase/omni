package parser

import (
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// parseGrantStmt parses a GRANT statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/grant.html
//
//	GRANT priv_type [, priv_type] ... ON [object_type] priv_level TO user [, user] ... [WITH GRANT OPTION]
//	GRANT ALL [PRIVILEGES] ON [object_type] priv_level TO user [, user] ... [WITH GRANT OPTION]
//	GRANT role [, role] ... TO user [, user] ... [WITH ADMIN OPTION]
func (p *Parser) parseGrantStmt() (nodes.Node, error) {
	start := p.pos()
	p.advance() // consume GRANT

	// Completion: after GRANT, offer privilege keyword candidates.
	p.checkCursor()
	if p.collectMode() {
		for _, t := range []int{kwALL, kwSELECT, kwINSERT, kwUPDATE, kwDELETE, kwCREATE, kwALTER, kwDROP, kwINDEX, kwEXECUTE} {
			p.addTokenCandidate(t)
		}
		return nil, &ParseError{Message: "collecting"}
	}

	// Detect PROXY keyword token before parsing the privilege list so we can
	// avoid an eqFold string comparison after the fact.
	isProxy := p.cur.Type == kwPROXY

	// Parse the names (could be privileges or roles).
	// We distinguish by looking for ON (privilege grant) vs TO (role grant).
	privs, allPriv, err := p.parsePrivilegeList()
	if err != nil {
		return nil, err
	}

	// If we see ON, it's a privilege grant; if TO, it's a role grant.
	if p.cur.Type == kwTO {
		// Role grant: GRANT role [, role] ... TO user [, user] ... [WITH ADMIN OPTION]
		p.advance() // consume TO
		users, err := p.parseUserList()
		if err != nil {
			return nil, err
		}
		stmt := &nodes.GrantRoleStmt{
			Loc:   nodes.Loc{Start: start},
			Roles: privs,
			To:    users,
		}
		// [WITH ADMIN OPTION]
		if p.cur.Type == kwWITH {
			p.advance()
			if p.cur.Type == kwADMIN {
				p.advance()
				if p.cur.Type == kwOPTION {
					p.advance()
				}
				stmt.WithAdmin = true
			}
		}
		stmt.Loc.End = p.pos()
		return stmt, nil
	}

	// Privilege grant
	stmt := &nodes.GrantStmt{Loc: nodes.Loc{Start: start}}
	stmt.Privileges = privs
	stmt.AllPriv = allPriv

	// GRANT PROXY ON user TO user [WITH GRANT OPTION]
	if len(privs) == 1 && isProxy && p.cur.Type == kwON {
		p.advance() // consume ON
		proxyUser, err := p.parseUserName()
		if err != nil {
			return nil, err
		}
		stmt.ProxyUser = proxyUser
		// TO user [, user] ...
		if _, err := p.expect(kwTO); err != nil {
			return nil, err
		}
		users, err := p.parseUserList()
		if err != nil {
			return nil, err
		}
		stmt.To = users
		// [WITH GRANT OPTION]
		if p.cur.Type == kwWITH {
			p.advance()
			if p.cur.Type == kwGRANT {
				p.advance()
				if p.cur.Type == kwOPTION {
					p.advance()
				}
				stmt.WithGrant = true
			}
		}
		stmt.Loc.End = p.pos()
		return stmt, nil
	}

	// ON [object_type] priv_level
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}

	target, err := p.parseGrantTarget()
	if err != nil {
		return nil, err
	}
	stmt.On = target

	// TO user [, user] ...
	if _, err := p.expect(kwTO); err != nil {
		return nil, err
	}

	users, err := p.parseUserList()
	if err != nil {
		return nil, err
	}
	stmt.To = users

	// [REQUIRE {NONE | tls_option [[AND] tls_option] ...}]
	if p.cur.Type == kwREQUIRE {
		req, err := p.parseRequireClause()
		if err != nil {
			return nil, err
		}
		stmt.Require = req
	}

	// [WITH [GRANT OPTION] [resource_option ...]]
	if p.cur.Type == kwWITH {
		p.advance()
		if p.cur.Type == kwGRANT {
			p.advance()
			if p.cur.Type == kwOPTION {
				p.advance()
			}
			stmt.WithGrant = true
		}
		// resource options after GRANT OPTION or standalone
		res := p.parseResourceOptions()
		if res != nil {
			stmt.Resource = res
		}
	}

	// [AS user [WITH ROLE {DEFAULT | NONE | ALL | ALL EXCEPT role [, role] ... | role [, role] ...}]]
	if p.cur.Type == kwAS {
		p.advance() // consume AS
		// Parse user name
		if p.cur.Type == tokSCONST {
			stmt.AsUser = p.cur.Str
			p.advance()
		} else if p.isIdentToken() {
			name, _, err := p.parseIdent()
			if err != nil {
				return nil, err
			}
			stmt.AsUser = name
		}
		// Optional @host
		if p.cur.Type == tokIDENT && p.cur.Str == "@" {
			p.advance()
			if p.cur.Type == tokSCONST {
				stmt.AsUser += "@" + p.cur.Str
				p.advance()
			} else if p.isIdentToken() {
				host, _, err := p.parseIdent()
				if err != nil {
					return nil, err
				}
				stmt.AsUser += "@" + host
			}
		}

		// [WITH ROLE ...]
		if p.cur.Type == kwWITH {
			p.advance() // consume WITH
			if p.cur.Type == kwROLE {
				p.advance() // consume ROLE
				switch {
				case p.cur.Type == kwDEFAULT:
					stmt.WithRoleType = "DEFAULT"
					p.advance()
				case p.cur.Type == kwNONE:
					stmt.WithRoleType = "NONE"
					p.advance()
				case p.cur.Type == kwALL:
					p.advance() // consume ALL
					if p.cur.Type == kwEXCEPT {
						p.advance() // consume EXCEPT
						stmt.WithRoleType = "ALL EXCEPT"
						// Parse role list
						for {
							role, _, err := p.parseRoleIdent()
							if err != nil {
								return nil, err
							}
							stmt.WithRoles = append(stmt.WithRoles, role)
							if p.cur.Type != ',' {
								break
							}
							p.advance()
						}
					} else {
						stmt.WithRoleType = "ALL"
					}
				default:
					// Role list
					for {
						role, _, err := p.parseRoleIdent()
						if err != nil {
							return nil, err
						}
						stmt.WithRoles = append(stmt.WithRoles, role)
						if p.cur.Type != ',' {
							break
						}
						p.advance()
					}
				}
			}
		}
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseRevokeStmt parses a REVOKE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/revoke.html
//
//	REVOKE [IF EXISTS] priv_type [(column_list)] [, priv_type [(column_list)]] ...
//	    ON [object_type] priv_level FROM user_or_role [, user_or_role] ... [IGNORE UNKNOWN USER]
//	REVOKE [IF EXISTS] ALL [PRIVILEGES], GRANT OPTION FROM user_or_role [, user_or_role] ... [IGNORE UNKNOWN USER]
//	REVOKE [IF EXISTS] PROXY ON user_or_role FROM user_or_role [, user_or_role] ... [IGNORE UNKNOWN USER]
//	REVOKE [IF EXISTS] role [, role] ... FROM user_or_role [, user_or_role] ... [IGNORE UNKNOWN USER]
func (p *Parser) parseRevokeStmt() (nodes.Node, error) {
	start := p.pos()
	p.advance() // consume REVOKE

	// [IF EXISTS]
	ifExists := false
	if p.cur.Type == kwIF {
		p.advance()
		p.match(kwEXISTS_KW)
		ifExists = true
	}

	// Detect PROXY keyword token before parsing the privilege list so we can
	// avoid an eqFold string comparison after the fact.
	isProxy := p.cur.Type == kwPROXY

	// Parse the names (could be privileges or roles).
	privs, allPriv, err := p.parsePrivilegeList()
	if err != nil {
		return nil, err
	}

	// REVOKE ALL PRIVILEGES, GRANT OPTION FROM user
	if allPriv && p.cur.Type == ',' {
		p.advance() // consume ','
		if p.cur.Type == kwGRANT {
			p.advance() // consume GRANT
			if p.cur.Type == kwOPTION {
				p.advance() // consume OPTION
			}
			if _, err := p.expect(kwFROM); err != nil {
				return nil, err
			}
			users, err := p.parseUserList()
			if err != nil {
				return nil, err
			}
			stmt := &nodes.RevokeStmt{
				Loc:         nodes.Loc{Start: start},
				IfExists:    ifExists,
				AllPriv:     true,
				GrantOption: true,
				From:        users,
			}
			// [IGNORE UNKNOWN USER]
			if p.cur.Type == kwIGNORE {
				p.advance()
				if p.cur.Type == kwUNKNOWN {
					p.advance()
				}
				if p.cur.Type == kwUSER {
					p.advance()
				}
				stmt.IgnoreUnknownUser = true
			}
			stmt.Loc.End = p.pos()
			return stmt, nil
		}
	}

	// If we see FROM, it's a role revoke; if ON, it's a privilege revoke.
	if p.cur.Type == kwFROM {
		// Role revoke: REVOKE role [, role] ... FROM user [, user] ...
		p.advance() // consume FROM
		users, err := p.parseUserList()
		if err != nil {
			return nil, err
		}
		stmt := &nodes.RevokeRoleStmt{
			Loc:      nodes.Loc{Start: start},
			IfExists: ifExists,
			Roles:    privs,
			From:     users,
		}
		// [IGNORE UNKNOWN USER]
		if p.cur.Type == kwIGNORE {
			p.advance()
			if p.cur.Type == kwUNKNOWN {
				p.advance()
			}
			if p.cur.Type == kwUSER {
				p.advance()
			}
			stmt.IgnoreUnknownUser = true
		}
		stmt.Loc.End = p.pos()
		return stmt, nil
	}

	// REVOKE PROXY ON user FROM user
	if len(privs) == 1 && isProxy && p.cur.Type == kwON {
		p.advance() // consume ON
		proxyUser, err := p.parseUserName()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(kwFROM); err != nil {
			return nil, err
		}
		users, err := p.parseUserList()
		if err != nil {
			return nil, err
		}
		stmt := &nodes.RevokeStmt{
			Loc:        nodes.Loc{Start: start},
			IfExists:   ifExists,
			Privileges: []string{"PROXY(" + proxyUser + ")"},
			From:       users,
		}
		// [IGNORE UNKNOWN USER]
		if p.cur.Type == kwIGNORE {
			p.advance()
			if p.cur.Type == kwUNKNOWN {
				p.advance()
			}
			if p.cur.Type == kwUSER {
				p.advance()
			}
			stmt.IgnoreUnknownUser = true
		}
		stmt.Loc.End = p.pos()
		return stmt, nil
	}

	// Privilege revoke
	stmt := &nodes.RevokeStmt{Loc: nodes.Loc{Start: start}}
	stmt.IfExists = ifExists
	stmt.Privileges = privs
	stmt.AllPriv = allPriv

	// ON [object_type] priv_level
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}

	target, err := p.parseGrantTarget()
	if err != nil {
		return nil, err
	}
	stmt.On = target

	// FROM user [, user] ...
	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}

	users, err := p.parseUserList()
	if err != nil {
		return nil, err
	}
	stmt.From = users

	// [IGNORE UNKNOWN USER]
	if p.cur.Type == kwIGNORE {
		p.advance()
		if p.cur.Type == kwUNKNOWN {
			p.advance()
		}
		if p.cur.Type == kwUSER {
			p.advance()
		}
		stmt.IgnoreUnknownUser = true
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parsePrivilegeList parses a comma-separated list of privilege names.
// Returns the list and whether ALL [PRIVILEGES] was specified.
func (p *Parser) parsePrivilegeList() ([]string, bool, error) {
	if p.cur.Type == kwALL {
		p.advance()
		// Optional PRIVILEGES keyword
		if p.cur.Type == kwPRIVILEGES {
			p.advance()
		}
		return nil, true, nil
	}

	var privs []string
	for {
		name, err := p.parsePrivilegeName()
		if err != nil {
			return nil, false, err
		}
		privs = append(privs, name)

		if p.cur.Type != ',' {
			break
		}
		// Peek: if comma is followed by ON, it's not a separator for privileges
		// but actually the ON keyword is after the list. But we need to check
		// if the next token after comma could be a privilege or ON.
		p.advance() // consume ','
	}

	return privs, false, nil
}

// parsePrivilegeName parses a single privilege name.
// Privilege names can be keywords like SELECT, INSERT, UPDATE, DELETE, etc.
// Multi-word: CREATE VIEW, ALTER ROUTINE, SHOW DATABASES, CREATE TEMPORARY TABLES,
// LOCK TABLES, REPLICATION CLIENT, REPLICATION SLAVE, CREATE TABLESPACE,
// CREATE USER, CREATE ROLE, DROP ROLE.
// Column-level: SELECT (col1, col2), INSERT (col1).
func (p *Parser) parsePrivilegeName() (string, error) {
	var name string
	switch p.cur.Type {
	case kwSELECT, kwINSERT, kwUPDATE, kwDELETE, kwCREATE, kwDROP, kwALTER,
		kwINDEX, kwEXECUTE, kwGRANT, kwREFERENCES, kwTRIGGER, kwEVENT,
		kwSHOW, kwLOCK:
		tok := p.advance()
		name = strings.ToUpper(tok.Str)
	case tokIDENT:
		tok := p.advance()
		name = strings.ToUpper(tok.Str)
	default:
		// Try accepting any keyword as a privilege name
		if p.cur.Type >= 700 {
			tok := p.advance()
			name = strings.ToUpper(tok.Str)
		} else {
			return "", &ParseError{
				Message:  "expected privilege name",
				Position: p.cur.Loc,
			}
		}
	}

	// Handle multi-word privilege names
	switch name {
	case "CREATE":
		// CREATE VIEW, CREATE TEMPORARY TABLES, CREATE TABLESPACE, CREATE USER, CREATE ROLE, CREATE ROUTINE
		switch {
		case p.cur.Type == kwVIEW:
			p.advance()
			name = "CREATE VIEW"
		case p.cur.Type == kwTEMPORARY:
			p.advance()
			if p.cur.Type == kwTABLES {
				p.advance()
			}
			name = "CREATE TEMPORARY TABLES"
		case p.cur.Type == kwTABLESPACE:
			p.advance()
			name = "CREATE TABLESPACE"
		case p.cur.Type == kwUSER:
			p.advance()
			name = "CREATE USER"
		case p.cur.Type == kwROLE:
			p.advance()
			name = "CREATE ROLE"
		case p.cur.Type == kwROUTINE:
			p.advance()
			name = "CREATE ROUTINE"
		}
	case "ALTER":
		// ALTER ROUTINE
		if p.cur.Type == kwROUTINE {
			p.advance()
			name = "ALTER ROUTINE"
		}
	case "SHOW":
		// SHOW DATABASES, SHOW VIEW
		if p.cur.Type == kwDATABASES || p.cur.Type == kwDATABASE {
			p.advance()
			name = "SHOW DATABASES"
		} else if p.cur.Type == kwVIEW {
			p.advance()
			name = "SHOW VIEW"
		}
	case "LOCK":
		// LOCK TABLES
		if p.cur.Type == kwTABLES {
			p.advance()
			name = "LOCK TABLES"
		}
	case "DROP":
		// DROP ROLE
		if p.cur.Type == kwROLE {
			p.advance()
			name = "DROP ROLE"
		}
	case "REPLICATION":
		// REPLICATION CLIENT, REPLICATION SLAVE
		if p.cur.Type == kwCLIENT {
			p.advance()
			name = "REPLICATION CLIENT"
		} else if p.cur.Type == kwSLAVE {
			p.advance()
			name = "REPLICATION SLAVE"
		}
	}

	// Column-level privilege: priv_type (col_list)
	if p.cur.Type == '(' {
		p.advance()
		var cols []string
		for {
			col, _, err := p.parseIdent()
			if err != nil {
				return "", err
			}
			cols = append(cols, col)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if _, err := p.expect(')'); err != nil {
			return "", err
		}
		name = name + "(" + strings.Join(cols, ", ") + ")"
	}

	return name, nil
}

// parseGrantTarget parses the ON clause target of a GRANT/REVOKE.
//
//	[object_type] priv_level
//	object_type: TABLE | FUNCTION | PROCEDURE
//	priv_level: *.* | db_name.* | db_name.tbl_name | tbl_name | *
func (p *Parser) parseGrantTarget() (*nodes.GrantTarget, error) {
	start := p.pos()
	target := &nodes.GrantTarget{Loc: nodes.Loc{Start: start}}

	// Completion: after ON, offer table_ref candidates.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("table_ref")
		return nil, &ParseError{Message: "collecting"}
	}

	// Optional object type
	switch p.cur.Type {
	case kwTABLE:
		target.Type = "TABLE"
		p.advance()
	case kwFUNCTION:
		target.Type = "FUNCTION"
		p.advance()
	case kwPROCEDURE:
		target.Type = "PROCEDURE"
		p.advance()
	}

	// Parse priv_level: *.* | db.* | db.tbl | tbl | *
	if p.cur.Type == '*' {
		p.advance()
		if p.cur.Type == '.' {
			p.advance() // consume '.'
			// *.*
			ref := &nodes.TableRef{
				Loc:    nodes.Loc{Start: start, End: p.pos()},
				Schema: "*",
				Name:   "*",
			}
			if p.cur.Type == '*' {
				p.advance()
				ref.Loc.End = p.pos()
			} else {
				// *.tbl_name (unusual but handle)
				name, _, err := p.parseIdent()
				if err != nil {
					return nil, err
				}
				ref.Name = name
				ref.Loc.End = p.pos()
			}
			target.Name = ref
		} else {
			// Just *
			target.Name = &nodes.TableRef{
				Loc:  nodes.Loc{Start: start, End: p.pos()},
				Name: "*",
			}
		}
	} else {
		// db.* | db.tbl | tbl
		ref, err := p.parseGrantPrivLevel()
		if err != nil {
			return nil, err
		}
		target.Name = ref
	}

	target.Loc.End = p.pos()
	return target, nil
}

// parseGrantPrivLevel parses a priv_level starting with an identifier.
// Handles: db_name.* | db_name.tbl_name | tbl_name
func (p *Parser) parseGrantPrivLevel() (*nodes.TableRef, error) {
	start := p.pos()
	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	ref := &nodes.TableRef{
		Loc:  nodes.Loc{Start: start},
		Name: name,
	}

	if p.cur.Type == '.' {
		p.advance() // consume '.'
		if p.cur.Type == '*' {
			p.advance()
			ref.Schema = name
			ref.Name = "*"
		} else {
			name2, _, err := p.parseIdent()
			if err != nil {
				return nil, err
			}
			ref.Schema = name
			ref.Name = name2
		}
	}

	ref.Loc.End = p.pos()
	return ref, nil
}

// parseUserList parses a comma-separated list of user identifiers.
// User can be 'user'@'host' or just a name.
func (p *Parser) parseUserList() ([]string, error) {
	var users []string
	for {
		user, err := p.parseUserName()
		if err != nil {
			return nil, err
		}
		users = append(users, user)

		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	return users, nil
}

// parseUserName parses a user identifier, which may be 'user'@'host'.
func (p *Parser) parseUserName() (string, error) {
	var name string

	if p.cur.Type == tokSCONST {
		name = p.cur.Str
		p.advance()
	} else if p.isIdentToken() {
		n, _, err := p.parseIdent()
		if err != nil {
			return "", err
		}
		name = n
	} else {
		return "", &ParseError{
			Message:  "expected user name",
			Position: p.cur.Loc,
		}
	}

	// Check for @'host' - the lexer scans @ as tokIDENT with Str="@"
	if p.cur.Type == tokIDENT && p.cur.Str == "@" {
		p.advance() // consume '@'
		var host string
		if p.cur.Type == tokSCONST {
			host = p.cur.Str
			p.advance()
		} else if p.isIdentToken() {
			h, _, err := p.parseIdent()
			if err != nil {
				return "", err
			}
			host = h
		} else {
			return "", &ParseError{
				Message:  "expected host name after @",
				Position: p.cur.Loc,
			}
		}
		return name + "@" + host, nil
	}

	return name, nil
}

// parseCreateUserStmt parses a CREATE USER statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/create-user.html
//
//	CREATE USER [IF NOT EXISTS]
//	    user [auth_option] [, user [auth_option]] ...
//	    DEFAULT ROLE role [, role ] ...
//	    [REQUIRE {NONE | tls_option [[AND] tls_option] ...}]
//	    [WITH resource_option [resource_option] ...]
//	    [password_option | lock_option] ...
//	    [COMMENT 'comment_string' | ATTRIBUTE 'json_object']
func (p *Parser) parseCreateUserStmt() (*nodes.CreateUserStmt, error) {
	start := p.pos()
	p.advance() // consume USER

	stmt := &nodes.CreateUserStmt{Loc: nodes.Loc{Start: start}}

	// [IF NOT EXISTS]
	if p.cur.Type == kwIF {
		p.advance()
		p.match(kwNOT)
		p.match(kwEXISTS_KW)
		stmt.IfNotExists = true
	}

	// user_spec [, user_spec] ...
	for {
		spec, err := p.parseUserSpec()
		if err != nil {
			return nil, err
		}
		stmt.Users = append(stmt.Users, spec)

		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}

	// [DEFAULT ROLE role [, role] ...]
	if p.cur.Type == kwDEFAULT {
		next := p.peekNext()
		if next.Type == kwROLE {
			p.advance() // consume DEFAULT
			p.advance() // consume ROLE
			roles, err := p.parseUserList()
			if err != nil {
				return nil, err
			}
			stmt.DefaultRoles = roles
		}
	}

	// [REQUIRE {NONE | tls_option [[AND] tls_option] ...}]
	if p.cur.Type == kwREQUIRE {
		req, err := p.parseRequireClause()
		if err != nil {
			return nil, err
		}
		stmt.Require = req
	}

	// [WITH resource_option ...]
	if p.cur.Type == kwWITH {
		p.advance()
		res := p.parseResourceOptions()
		if res != nil {
			stmt.Resource = res
		}
	}

	// [password_option | lock_option] ...
	p.parseUserAccountOptions(
		&stmt.PasswordExpire, &stmt.PasswordHistory, &stmt.PasswordReuseInterval,
		&stmt.PasswordRequireCurrent, &stmt.FailedLoginAttempts, &stmt.HasFailedLogin,
		&stmt.PasswordLockTime, &stmt.AccountLock, &stmt.AccountUnlock,
		&stmt.Comment, &stmt.Attribute,
	)

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseDropUserStmt parses a DROP USER statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/drop-user.html
//
//	DROP USER [IF EXISTS] user [, user] ...
func (p *Parser) parseDropUserStmt() (*nodes.DropUserStmt, error) {
	start := p.pos()
	p.advance() // consume USER

	stmt := &nodes.DropUserStmt{Loc: nodes.Loc{Start: start}}

	// [IF EXISTS]
	if p.cur.Type == kwIF {
		p.advance()
		p.match(kwEXISTS_KW)
		stmt.IfExists = true
	}

	// user [, user] ...
	users, err := p.parseUserList()
	if err != nil {
		return nil, err
	}
	stmt.Users = users

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseAlterUserStmt parses an ALTER USER statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/alter-user.html
//
//	ALTER USER user_spec [, user_spec] ...
//	    [REQUIRE {NONE | tls_option [[AND] tls_option] ...}]
//	    [WITH resource_option [resource_option] ...]
func (p *Parser) parseAlterUserStmt() (*nodes.AlterUserStmt, error) {
	start := p.pos()
	p.advance() // consume USER

	stmt := &nodes.AlterUserStmt{Loc: nodes.Loc{Start: start}}

	// [IF EXISTS]
	if p.cur.Type == kwIF {
		p.advance()
		p.match(kwEXISTS_KW)
		stmt.IfExists = true
	}

	// Check for ALTER USER USER() ... (user_func_auth_option or registration_option)
	if p.cur.Type == kwUSER && p.peekNext().Type == '(' {
		stmt.IsUserFunc = true
		p.advance() // consume USER
		p.advance() // consume (
		p.match(')')

		// Check for registration_option: {2|3} FACTOR ...
		if p.cur.Type == tokICONST && (p.cur.Str == "2" || p.cur.Str == "3") {
			op, err := p.parseRegistrationOption()
			if err != nil {
				return nil, err
			}
			stmt.RegistrationOp = op
			stmt.Loc.End = p.pos()
			return stmt, nil
		}

		// user_func_auth_option: IDENTIFIED BY 'auth_string' [REPLACE ...] [RETAIN ...] | DISCARD OLD PASSWORD
		userSpec := &nodes.UserSpec{Loc: nodes.Loc{Start: start}, Name: "USER()"}
		if p.cur.Type == kwIDENTIFIED {
			if err := p.parseUserAuthOption(userSpec); err != nil {
				return nil, err
			}
		}
		// [REPLACE 'current_auth_string']
		if p.cur.Type == kwREPLACE {
			p.advance()
			if p.cur.Type == tokSCONST {
				userSpec.Replace = p.cur.Str
				p.advance()
			}
		}
		// RETAIN CURRENT PASSWORD
		if p.cur.Type == kwRETAIN {
			p.advance()
			p.advance()
			p.advance()
			userSpec.RetainCurrentPassword = true
		}
		// DISCARD OLD PASSWORD
		if p.cur.Type == kwDISCARD {
			p.advance()
			p.advance()
			p.advance()
			userSpec.DiscardOldPassword = true
		}
		userSpec.Loc.End = p.pos()
		stmt.Users = append(stmt.Users, userSpec)
		stmt.Loc.End = p.pos()
		return stmt, nil
	}

	// Check for ALTER USER user DEFAULT ROLE ...
	// We parse the first user, then check if DEFAULT follows.
	spec, err := p.parseUserSpec()
	if err != nil {
		return nil, err
	}
	stmt.Users = append(stmt.Users, spec)

	// Check for registration_option after user: {2|3} FACTOR ...
	if p.cur.Type == tokICONST && (p.cur.Str == "2" || p.cur.Str == "3") {
		next := p.peekNext()
		if next.Type == kwFACTOR {
			op, err := p.parseRegistrationOption()
			if err != nil {
				return nil, err
			}
			stmt.RegistrationOp = op
			stmt.Loc.End = p.pos()
			return stmt, nil
		}
	}

	// Check for ADD/MODIFY/DROP factor operations in auth_option
	if p.cur.Type == kwADD || p.cur.Type == kwMODIFY || p.cur.Type == kwDROP {
		next := p.peekNext()
		if next.Type == tokICONST && (next.Str == "2" || next.Str == "3") {
			ops, err := p.parseFactorOps()
			if err != nil {
				return nil, err
			}
			stmt.FactorOps = ops
			stmt.Loc.End = p.pos()
			return stmt, nil
		}
	}

	if p.cur.Type == kwDEFAULT {
		next := p.peekNext()
		if next.Type == kwROLE {
			p.advance() // consume DEFAULT
			p.advance() // consume ROLE
			stmt.DefaultRoleUser = spec.Name
			if spec.Host != "" {
				stmt.DefaultRoleUser += "@" + spec.Host
			}
			switch {
			case p.cur.Type == kwNONE:
				stmt.DefaultRoleType = "NONE"
				p.advance()
			case p.cur.Type == kwALL:
				stmt.DefaultRoleType = "ALL"
				p.advance()
			default:
				roles, err := p.parseUserList()
				if err != nil {
					return nil, err
				}
				stmt.DefaultRoles = roles
			}
			stmt.Loc.End = p.pos()
			return stmt, nil
		}
	}

	// More user_specs
	for p.cur.Type == ',' {
		p.advance() // consume ','
		spec, err := p.parseUserSpec()
		if err != nil {
			return nil, err
		}
		stmt.Users = append(stmt.Users, spec)
	}

	// [REQUIRE {NONE | tls_option [[AND] tls_option] ...}]
	if p.cur.Type == kwREQUIRE {
		req, err := p.parseRequireClause()
		if err != nil {
			return nil, err
		}
		stmt.Require = req
	}

	// [WITH resource_option ...]
	if p.cur.Type == kwWITH {
		p.advance()
		res := p.parseResourceOptions()
		if res != nil {
			stmt.Resource = res
		}
	}

	// [password_option | lock_option] ...
	p.parseUserAccountOptions(
		&stmt.PasswordExpire, &stmt.PasswordHistory, &stmt.PasswordReuseInterval,
		&stmt.PasswordRequireCurrent, &stmt.FailedLoginAttempts, &stmt.HasFailedLogin,
		&stmt.PasswordLockTime, &stmt.AccountLock, &stmt.AccountUnlock,
		&stmt.Comment, &stmt.Attribute,
	)

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseRegistrationOption parses a registration_option:
//
//	factor INITIATE REGISTRATION
//	factor FINISH REGISTRATION SET CHALLENGE_RESPONSE AS 'auth_string'
//	factor UNREGISTER
//
//	factor: {2 | 3} FACTOR
func (p *Parser) parseRegistrationOption() (*nodes.RegistrationOp, error) {
	start := p.pos()
	op := &nodes.RegistrationOp{Loc: nodes.Loc{Start: start}}

	// {2|3}
	if p.cur.Type == tokICONST {
		if p.cur.Str == "2" {
			op.Factor = 2
		} else {
			op.Factor = 3
		}
		p.advance()
	}
	// FACTOR
	if p.cur.Type == kwFACTOR {
		p.advance()
	}

	// INITIATE REGISTRATION | FINISH REGISTRATION ... | UNREGISTER
	if p.cur.Type == kwINITIATE {
		op.Action = "INITIATE"
		p.advance()
		// REGISTRATION
		if p.cur.Type == kwREGISTRATION {
			p.advance()
		}
	} else if p.cur.Type == kwFINISH {
		op.Action = "FINISH"
		p.advance()
		// REGISTRATION
		if p.cur.Type == kwREGISTRATION {
			p.advance()
		}
		// SET CHALLENGE_RESPONSE AS 'auth_string'
		if p.cur.Type == kwSET {
			p.advance()
			if p.cur.Type == kwCHALLENGE_RESPONSE {
				p.advance()
			}
			p.match(kwAS)
			if p.cur.Type == tokSCONST {
				op.ChallengeResponse = p.cur.Str
				p.advance()
			}
		}
	} else if p.cur.Type == kwUNREGISTER {
		op.Action = "UNREGISTER"
		p.advance()
	}

	op.Loc.End = p.pos()
	return op, nil
}

// parseFactorOps parses ADD/MODIFY/DROP factor operations:
//
//	ADD factor factor_auth_option [ADD factor factor_auth_option]
//	MODIFY factor factor_auth_option [MODIFY factor factor_auth_option]
//	DROP factor [DROP factor]
func (p *Parser) parseFactorOps() ([]*nodes.FactorOp, error) {
	var ops []*nodes.FactorOp
	for p.cur.Type == kwADD || p.cur.Type == kwMODIFY || p.cur.Type == kwDROP {
		start := p.pos()
		op := &nodes.FactorOp{Loc: nodes.Loc{Start: start}}

		switch p.cur.Type {
		case kwADD:
			op.Action = "ADD"
		case kwMODIFY:
			op.Action = "MODIFY"
		case kwDROP:
			op.Action = "DROP"
		}
		p.advance() // consume ADD/MODIFY/DROP

		// {2|3}
		if p.cur.Type == tokICONST && (p.cur.Str == "2" || p.cur.Str == "3") {
			if p.cur.Str == "2" {
				op.Factor = 2
			} else {
				op.Factor = 3
			}
			p.advance()
		}
		// FACTOR
		if p.cur.Type == kwFACTOR {
			p.advance()
		}

		// For ADD and MODIFY: parse factor_auth_option
		if op.Action != "DROP" {
			if err := p.parseFactorAuthOption(op); err != nil {
				return nil, err
			}
		}

		op.Loc.End = p.pos()
		ops = append(ops, op)
	}
	return ops, nil
}

// parseFactorAuthOption parses factor_auth_option:
//
//	IDENTIFIED BY 'auth_string'
//	IDENTIFIED BY RANDOM PASSWORD
//	IDENTIFIED WITH auth_plugin BY 'auth_string'
//	IDENTIFIED WITH auth_plugin BY RANDOM PASSWORD
//	IDENTIFIED WITH auth_plugin AS 'auth_string'
func (p *Parser) parseFactorAuthOption(op *nodes.FactorOp) error {
	if p.cur.Type != kwIDENTIFIED {
		return nil
	}
	p.advance() // consume IDENTIFIED

	if p.cur.Type == kwBY {
		p.advance()
		if p.cur.Type == kwRANDOM {
			p.advance() // consume RANDOM
			p.advance() // consume PASSWORD
			op.PasswordRandom = true
		} else if p.cur.Type == tokSCONST {
			op.Password = p.cur.Str
			p.advance()
		}
	} else if p.cur.Type == kwWITH {
		p.advance()
		if p.isIdentToken() {
			plugin, _, err := p.parseIdent()
			if err != nil {
				return err
			}
			op.AuthPlugin = plugin
		} else {
			return p.syntaxErrorAtCur()
		}
		if p.cur.Type == kwBY {
			p.advance()
			if p.cur.Type == kwRANDOM {
				p.advance() // consume RANDOM
				p.advance() // consume PASSWORD
				op.PasswordRandom = true
			} else if p.cur.Type == tokSCONST {
				op.Password = p.cur.Str
				p.advance()
			}
		} else if p.cur.Type == kwAS {
			p.advance()
			if p.cur.Type == tokSCONST {
				op.AuthHash = p.cur.Str
				p.advance()
			}
		}
	}
	return nil
}

// parseUserSpec parses a user specification.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/create-user.html
//
//	user [auth_option] [, user [auth_option]] ...
//
//	auth_option: {
//	    IDENTIFIED BY 'auth_string' [AND 2fa_auth_option]
//	  | IDENTIFIED BY RANDOM PASSWORD [AND 2fa_auth_option]
//	  | IDENTIFIED WITH auth_plugin [AND 2fa_auth_option]
//	  | IDENTIFIED WITH auth_plugin BY 'auth_string' [AND 2fa_auth_option]
//	  | IDENTIFIED WITH auth_plugin BY RANDOM PASSWORD [AND 2fa_auth_option]
//	  | IDENTIFIED WITH auth_plugin AS 'auth_string' [AND 2fa_auth_option]
//	  | IDENTIFIED WITH auth_plugin [initial_auth_option]
//	}
//
// ALTER USER also supports [REPLACE 'current_auth_string'] [RETAIN CURRENT PASSWORD]
// after IDENTIFIED BY forms.
func (p *Parser) parseUserSpec() (*nodes.UserSpec, error) {
	start := p.pos()
	spec := &nodes.UserSpec{Loc: nodes.Loc{Start: start}}

	// Parse user name
	if p.cur.Type == tokSCONST {
		spec.Name = p.cur.Str
		p.advance()
	} else if p.isIdentToken() {
		name, _, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		spec.Name = name
	} else {
		return nil, &ParseError{
			Message:  "expected user name",
			Position: p.cur.Loc,
		}
	}

	// Optional @host — the lexer scans @ as tokIDENT with Str="@"
	if p.cur.Type == tokIDENT && p.cur.Str == "@" {
		p.advance()
		if p.cur.Type == tokSCONST {
			spec.Host = p.cur.Str
			p.advance()
		} else if p.isIdentToken() {
			host, _, err := p.parseIdent()
			if err != nil {
				return nil, err
			}
			spec.Host = host
		}
	}

	// Optional IDENTIFIED BY/WITH
	if p.cur.Type == kwIDENTIFIED {
		if err := p.parseUserAuthOption(spec); err != nil {
			return nil, err
		}
	}

	// [REPLACE 'current_auth_string'] (ALTER USER)
	if p.cur.Type == kwREPLACE {
		p.advance()
		if p.cur.Type == tokSCONST {
			spec.Replace = p.cur.Str
			p.advance()
		}
	}

	// RETAIN CURRENT PASSWORD
	if p.cur.Type == kwRETAIN {
		p.advance() // consume RETAIN
		p.advance() // consume CURRENT
		p.advance() // consume PASSWORD
		spec.RetainCurrentPassword = true
	}

	// DISCARD OLD PASSWORD
	if p.cur.Type == kwDISCARD {
		p.advance() // consume DISCARD
		p.advance() // consume OLD
		p.advance() // consume PASSWORD
		spec.DiscardOldPassword = true
	}

	// [AND 2fa_auth_option [AND 3fa_auth_option]]
	for p.cur.Type == kwAND {
		next := p.peekNext()
		if next.Type != kwIDENTIFIED {
			break
		}
		p.advance() // consume AND
		factor := &nodes.UserSpec{Loc: nodes.Loc{Start: p.pos()}}
		if err := p.parseUserAuthOption(factor); err != nil {
			return nil, err
		}
		factor.Loc.End = p.pos()
		spec.AuthFactors = append(spec.AuthFactors, factor)
	}

	// [INITIAL AUTHENTICATION ...]
	if p.cur.Type == kwINITIAL {
		next := p.peekNext()
		if next.Type == kwAUTHENTICATION {
			p.advance() // consume INITIAL
			p.advance() // consume AUTHENTICATION
			if p.cur.Type == kwIDENTIFIED {
				p.advance() // consume IDENTIFIED
				if p.cur.Type == kwBY {
					p.advance() // consume BY
					if p.cur.Type == kwRANDOM {
						p.advance() // consume RANDOM
						p.advance() // consume PASSWORD
						spec.InitialAuthRandom = true
					} else if p.cur.Type == tokSCONST {
						spec.InitialAuthString = p.cur.Str
						p.advance()
					}
				} else if p.cur.Type == kwWITH {
					p.advance() // consume WITH
					if p.isIdentToken() {
						plugin, _, err := p.parseIdent()
						if err != nil {
							return nil, err
						}
						spec.InitialAuthPlugin = plugin
					} else {
						return nil, p.syntaxErrorAtCur()
					}
					if p.cur.Type == kwAS {
						p.advance()
						if p.cur.Type == tokSCONST {
							spec.InitialAuthString = p.cur.Str
							p.advance()
						}
					}
				}
			}
		}
	}

	spec.Loc.End = p.pos()
	return spec, nil
}

// parseCreateRoleStmt parses a CREATE ROLE statement.
//
//	CREATE ROLE [IF NOT EXISTS] role [, role] ...
func (p *Parser) parseCreateRoleStmt() (*nodes.CreateRoleStmt, error) {
	start := p.pos()
	p.advance() // consume ROLE

	stmt := &nodes.CreateRoleStmt{Loc: nodes.Loc{Start: start}}

	// [IF NOT EXISTS]
	if p.cur.Type == kwIF {
		p.advance()
		p.match(kwNOT)
		p.match(kwEXISTS_KW)
		stmt.IfNotExists = true
	}

	// role [, role] ...
	roles, err := p.parseUserList()
	if err != nil {
		return nil, err
	}
	stmt.Roles = roles

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseDropRoleStmt parses a DROP ROLE statement.
//
//	DROP ROLE [IF EXISTS] role [, role] ...
func (p *Parser) parseDropRoleStmt() (*nodes.DropRoleStmt, error) {
	start := p.pos()
	p.advance() // consume ROLE

	stmt := &nodes.DropRoleStmt{Loc: nodes.Loc{Start: start}}

	// [IF EXISTS]
	if p.cur.Type == kwIF {
		p.advance()
		p.match(kwEXISTS_KW)
		stmt.IfExists = true
	}

	// role [, role] ...
	roles, err := p.parseUserList()
	if err != nil {
		return nil, err
	}
	stmt.Roles = roles

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseUserAuthOption parses the IDENTIFIED BY/WITH portion of a user spec.
// Caller has verified p.cur.Type == kwIDENTIFIED.
func (p *Parser) parseUserAuthOption(spec *nodes.UserSpec) error {
	p.advance() // consume IDENTIFIED
	if p.cur.Type == kwBY {
		// IDENTIFIED BY 'password' | IDENTIFIED BY RANDOM PASSWORD
		p.advance()
		if p.cur.Type == kwRANDOM {
			p.advance() // consume RANDOM
			p.advance() // consume PASSWORD
			spec.PasswordRandom = true
		} else if p.cur.Type == tokSCONST {
			spec.Password = p.cur.Str
			p.advance()
		}
	} else if p.cur.Type == kwWITH {
		// IDENTIFIED WITH auth_plugin [BY 'password' | BY RANDOM PASSWORD | AS 'hash']
		p.advance()
		if p.isIdentToken() {
			plugin, _, err := p.parseIdent()
			if err != nil {
				return err
			}
			spec.AuthPlugin = plugin
		} else {
			return p.syntaxErrorAtCur()
		}
		if p.cur.Type == kwBY {
			p.advance()
			if p.cur.Type == kwRANDOM {
				p.advance() // consume RANDOM
				p.advance() // consume PASSWORD
				spec.PasswordRandom = true
			} else if p.cur.Type == tokSCONST {
				spec.Password = p.cur.Str
				p.advance()
			}
		} else if p.cur.Type == kwAS {
			p.advance()
			if p.cur.Type == tokSCONST {
				spec.AuthHash = p.cur.Str
				p.advance()
			}
		}
	}
	return nil
}

// parseSetDefaultRoleStmt parses a SET DEFAULT ROLE statement.
//
//	SET DEFAULT ROLE {NONE | ALL | role [, role] ...} TO user [, user] ...
func (p *Parser) parseSetDefaultRoleStmt(start int) (*nodes.SetDefaultRoleStmt, error) {
	// DEFAULT and ROLE already consumed
	stmt := &nodes.SetDefaultRoleStmt{Loc: nodes.Loc{Start: start}}

	switch {
	case p.cur.Type == kwNONE:
		stmt.None = true
		p.advance()
	case p.cur.Type == kwALL:
		stmt.All = true
		p.advance()
	default:
		roles, err := p.parseUserList()
		if err != nil {
			return nil, err
		}
		stmt.Roles = roles
	}

	// TO user [, user] ...
	if _, err := p.expect(kwTO); err != nil {
		return nil, err
	}

	users, err := p.parseUserList()
	if err != nil {
		return nil, err
	}
	stmt.To = users

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseSetRoleStmt parses a SET ROLE statement.
//
//	SET ROLE {DEFAULT | NONE | ALL | ALL EXCEPT role [, role] ... | role [, role] ...}
func (p *Parser) parseSetRoleStmt(start int) (*nodes.SetRoleStmt, error) {
	// ROLE already consumed
	stmt := &nodes.SetRoleStmt{Loc: nodes.Loc{Start: start}}

	switch {
	case p.cur.Type == kwDEFAULT:
		stmt.Default = true
		p.advance()
	case p.cur.Type == kwNONE:
		stmt.None = true
		p.advance()
	case p.cur.Type == kwALL:
		stmt.All = true
		p.advance()
		// ALL EXCEPT role [, role] ...
		if p.cur.Type == kwEXCEPT {
			p.advance()
			roles, err := p.parseUserList()
			if err != nil {
				return nil, err
			}
			stmt.AllExcept = roles
		}
	default:
		roles, err := p.parseUserList()
		if err != nil {
			return nil, err
		}
		stmt.Roles = roles
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseRenameUserStmt parses RENAME USER old_user TO new_user [, old_user TO new_user] ...
// RENAME already consumed. p.cur is USER.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/rename-user.html
//
//	RENAME USER old_user TO new_user [, old_user TO new_user] ...
func (p *Parser) parseRenameUserStmt(start int) (*nodes.RenameUserStmt, error) {
	p.advance() // consume USER

	stmt := &nodes.RenameUserStmt{
		Loc: nodes.Loc{Start: start},
	}

	for {
		pairStart := p.pos()
		pair := &nodes.RenameUserPair{Loc: nodes.Loc{Start: pairStart}}

		// old_user — can be 'user'@'host' or identifier@identifier
		var err error
		pair.OldUser, pair.OldHost, err = p.parseRenameUserPart()
		if err != nil {
			return nil, err
		}

		// TO
		p.match(kwTO)

		// new_user
		pair.NewUser, pair.NewHost, err = p.parseRenameUserPart()
		if err != nil {
			return nil, err
		}

		pair.Loc.End = p.pos()
		stmt.Pairs = append(stmt.Pairs, pair)

		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseRenameUserPart parses a user part (name[@host]) in a RENAME USER statement.
// Supports both quoted and unquoted user names.
func (p *Parser) parseRenameUserPart() (string, string, error) {
	var name, host string
	if p.cur.Type == tokSCONST {
		name = p.cur.Str
		p.advance()
	} else if p.isIdentToken() {
		var err error
		name, _, err = p.parseIdent()
		if err != nil {
			return "", "", err
		}
	} else {
		return "", "", p.syntaxErrorAtCur()
	}
	if p.cur.Type == tokIDENT && p.cur.Str == "@" {
		p.advance()
		if p.cur.Type == tokSCONST {
			host = p.cur.Str
			p.advance()
		} else if p.isIdentToken() {
			var err error
			host, _, err = p.parseIdent()
			if err != nil {
				return "", "", err
			}
		}
	}
	return name, host, nil
}

// parseRequireClause parses a REQUIRE clause for TLS options.
//
//	REQUIRE {NONE | tls_option [[AND] tls_option] ...}
//
//	tls_option: {
//	    SSL
//	  | X509
//	  | CIPHER 'cipher'
//	  | ISSUER 'issuer'
//	  | SUBJECT 'subject'
//	}
func (p *Parser) parseRequireClause() (*nodes.RequireClause, error) {
	start := p.pos()
	p.advance() // consume REQUIRE

	req := &nodes.RequireClause{Loc: nodes.Loc{Start: start}}

	for {
		switch {
		case p.cur.Type == kwNONE:
			req.None = true
			p.advance()
		case p.cur.Type == kwSSL:
			req.SSL = true
			p.advance()
		case p.cur.Type == kwX509:
			req.X509 = true
			p.advance()
		case p.cur.Type == kwCIPHER:
			p.advance()
			if p.cur.Type == tokSCONST {
				req.Cipher = p.cur.Str
				p.advance()
			}
		case p.cur.Type == kwISSUER:
			p.advance()
			if p.cur.Type == tokSCONST {
				req.Issuer = p.cur.Str
				p.advance()
			}
		case p.cur.Type == kwSUBJECT:
			p.advance()
			if p.cur.Type == tokSCONST {
				req.Subject = p.cur.Str
				p.advance()
			}
		default:
			req.Loc.End = p.pos()
			return req, nil
		}

		// Optional AND between tls_option items
		if p.cur.Type == kwAND {
			p.advance()
		}
	}
}

// parseResourceOptions parses resource limit options.
// Caller has already consumed WITH.
//
//	resource_option: {
//	    MAX_QUERIES_PER_HOUR count
//	  | MAX_UPDATES_PER_HOUR count
//	  | MAX_CONNECTIONS_PER_HOUR count
//	  | MAX_USER_CONNECTIONS count
//	}
func (p *Parser) parseResourceOptions() *nodes.ResourceOption {
	var res *nodes.ResourceOption
	for p.cur.Type == kwMAX_QUERIES_PER_HOUR || p.cur.Type == kwMAX_UPDATES_PER_HOUR ||
		p.cur.Type == kwMAX_CONNECTIONS_PER_HOUR || p.cur.Type == kwMAX_USER_CONNECTIONS {
		if res == nil {
			res = &nodes.ResourceOption{Loc: nodes.Loc{Start: p.pos()}}
		}
		switch p.cur.Type {
		case kwMAX_QUERIES_PER_HOUR:
			p.advance()
			if p.cur.Type == tokICONST {
				res.MaxQueriesPerHour = int(p.cur.Ival)
				res.HasMaxQueries = true
				p.advance()
			}
		case kwMAX_UPDATES_PER_HOUR:
			p.advance()
			if p.cur.Type == tokICONST {
				res.MaxUpdatesPerHour = int(p.cur.Ival)
				res.HasMaxUpdates = true
				p.advance()
			}
		case kwMAX_CONNECTIONS_PER_HOUR:
			p.advance()
			if p.cur.Type == tokICONST {
				res.MaxConnectionsPerHour = int(p.cur.Ival)
				res.HasMaxConnections = true
				p.advance()
			}
		case kwMAX_USER_CONNECTIONS:
			p.advance()
			if p.cur.Type == tokICONST {
				res.MaxUserConnections = int(p.cur.Ival)
				res.HasMaxUserConn = true
				p.advance()
			}
		}
	}
	if res != nil {
		res.Loc.End = p.pos()
	}
	return res
}

// parseUserAccountOptions parses password options, lock options, comment, and attribute
// that can follow CREATE USER or ALTER USER statements.
//
//	password_option: {
//	    PASSWORD EXPIRE [DEFAULT | NEVER | INTERVAL N DAY]
//	  | PASSWORD HISTORY {DEFAULT | N}
//	  | PASSWORD REUSE INTERVAL {DEFAULT | N DAY}
//	  | PASSWORD REQUIRE CURRENT [DEFAULT | OPTIONAL]
//	  | FAILED_LOGIN_ATTEMPTS N
//	  | PASSWORD_LOCK_TIME {N | UNBOUNDED}
//	}
//
//	lock_option: {
//	    ACCOUNT LOCK
//	  | ACCOUNT UNLOCK
//	}
func (p *Parser) parseUserAccountOptions(
	passwordExpire, passwordHistory, passwordReuseInterval, passwordRequireCurrent *string,
	failedLoginAttempts *int, hasFailedLogin *bool,
	passwordLockTime *string,
	accountLock, accountUnlock *bool,
	comment, attribute *string,
) {
	for {
		switch {
		case p.cur.Type == kwPASSWORD:
			next := p.peekNext()
			switch {
			case next.Type == kwEXPIRE:
				p.advance() // consume PASSWORD
				p.advance() // consume EXPIRE
				if p.cur.Type == kwDEFAULT {
					*passwordExpire = "DEFAULT"
					p.advance()
				} else if p.cur.Type == kwNEVER {
					*passwordExpire = "NEVER"
					p.advance()
				} else if p.cur.Type == kwINTERVAL {
					p.advance() // consume INTERVAL
					n := "0"
					if p.cur.Type == tokICONST {
						n = p.cur.Str
						if n == "" {
							n = strconv.FormatInt(p.cur.Ival, 10)
						}
						p.advance()
					}
					// DAY
					if p.cur.Type == kwDAY {
						p.advance()
					}
					*passwordExpire = "INTERVAL " + n + " DAY"
				} else {
					*passwordExpire = "EXPIRE"
				}
			case next.Type == kwHISTORY:
				p.advance() // consume PASSWORD
				p.advance() // consume HISTORY
				if p.cur.Type == kwDEFAULT {
					*passwordHistory = "DEFAULT"
					p.advance()
				} else if p.cur.Type == tokICONST {
					*passwordHistory = strconv.FormatInt(p.cur.Ival, 10)
					p.advance()
				}
			case next.Type == kwREUSE:
				p.advance() // consume PASSWORD
				p.advance() // consume REUSE
				// INTERVAL
				if p.cur.Type == kwINTERVAL {
					p.advance()
				}
				if p.cur.Type == kwDEFAULT {
					*passwordReuseInterval = "DEFAULT"
					p.advance()
				} else if p.cur.Type == tokICONST {
					n := strconv.FormatInt(p.cur.Ival, 10)
					p.advance()
					if p.cur.Type == kwDAY {
						p.advance()
					}
					*passwordReuseInterval = n + " DAY"
				}
			case next.Type == kwREQUIRE:
				p.advance() // consume PASSWORD
				p.advance() // consume REQUIRE
				if p.cur.Type == kwDEFAULT {
					*passwordRequireCurrent = "DEFAULT"
					p.advance()
				} else if p.cur.Type == kwCURRENT {
					p.advance() // consume CURRENT
					// Check for DEFAULT or OPTIONAL after CURRENT
					if p.cur.Type == kwDEFAULT {
						*passwordRequireCurrent = "DEFAULT"
						p.advance()
					} else if p.cur.Type == kwOPTIONAL {
						*passwordRequireCurrent = "OPTIONAL"
						p.advance()
					} else {
						*passwordRequireCurrent = "CURRENT"
					}
				}
			default:
				return
			}

		case p.cur.Type == kwFAILED_LOGIN_ATTEMPTS:
			p.advance()
			if p.cur.Type == tokICONST {
				*failedLoginAttempts = int(p.cur.Ival)
				*hasFailedLogin = true
				p.advance()
			}

		case p.cur.Type == kwPASSWORD_LOCK_TIME:
			p.advance()
			if p.cur.Type == tokICONST {
				*passwordLockTime = strconv.FormatInt(p.cur.Ival, 10)
				p.advance()
			} else if p.cur.Type == kwUNBOUNDED {
				*passwordLockTime = "UNBOUNDED"
				p.advance()
			}

		case p.cur.Type == kwACCOUNT:
			p.advance() // consume ACCOUNT
			if p.cur.Type == kwLOCK {
				*accountLock = true
				p.advance()
			} else if p.cur.Type == kwUNLOCK {
				*accountUnlock = true
				p.advance()
			}

		case p.cur.Type == kwCOMMENT:
			p.advance()
			if p.cur.Type == tokSCONST {
				*comment = p.cur.Str
				p.advance()
			}

		case p.cur.Type == kwATTRIBUTE:
			p.advance()
			if p.cur.Type == tokSCONST {
				*attribute = p.cur.Str
				p.advance()
			}

		default:
			return
		}
	}
}

// parseSetResourceGroupStmt parses SET RESOURCE GROUP group_name [FOR thread_id [, thread_id] ...].
// SET already consumed. p.cur is RESOURCE.
func (p *Parser) parseSetResourceGroupStmt(start int) (*nodes.SetResourceGroupStmt, error) {
	p.advance() // consume RESOURCE
	p.advance() // consume GROUP (identifier)

	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	stmt := &nodes.SetResourceGroupStmt{
		Loc:  nodes.Loc{Start: start},
		Name: name,
	}

	// Optional: FOR thread_id [, thread_id] ...
	if p.cur.Type == kwFOR {
		p.advance() // consume FOR
		for {
			if p.cur.Type == tokICONST {
				id, _ := strconv.ParseInt(p.cur.Str, 10, 64)
				stmt.ThreadIDs = append(stmt.ThreadIDs, id)
				p.advance()
			} else {
				break
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}
