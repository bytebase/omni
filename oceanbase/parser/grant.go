package parser

import (
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/oceanbase/ast"
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
		if err := p.requireRoles(privs, start); err != nil {
			return nil, err
		}
		p.advance() // consume TO
		users, err := p.parseUserList()
		if err != nil {
			return nil, err
		}
		stmt := &nodes.GrantRoleStmt{
			Loc:   nodes.Loc{Start: start},
			Roles: privNames(privs),
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
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	// Privilege grant
	stmt := &nodes.GrantStmt{Loc: nodes.Loc{Start: start}}
	stmt.Privileges = privNames(privs)
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
		stmt.Loc.End = p.prev.End
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

	users, err := p.parseGranteeList()
	if err != nil {
		return nil, err
	}
	stmt.To = users

	// [WITH GRANT OPTION]
	//
	// A GRANT carries neither a REQUIRE clause nor a resource option on
	// this engine, however CREATE USER and ALTER USER do: oceanbase-ce:4.3.5
	// answers 1064 for `GRANT ALL ON *.* TO u REQUIRE SSL` and for `GRANT
	// ALL ON *.* TO u WITH MAX_QUERIES_PER_HOUR 10`, and runs both as
	// clauses of CREATE USER. MariaDB and TiDB take them on a GRANT, which
	// is why their forks of this file still read them here.
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

	// [AS user [WITH ROLE {DEFAULT | NONE | ALL | ALL EXCEPT role [, role] ... | role [, role] ...}]]
	if p.cur.Type == kwAS {
		p.advance() // consume AS
		asUser, err := p.parseUserName()
		if err != nil {
			return nil, err
		}
		stmt.AsUser = asUser

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

	stmt.Loc.End = p.prev.End
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
			stmt.Loc.End = p.prev.End
			return stmt, nil
		}
	}

	// If we see FROM, it's a role revoke; if ON, it's a privilege revoke.
	if p.cur.Type == kwFROM {
		// Role revoke: REVOKE role [, role] ... FROM user [, user] ...
		if err := p.requireRoles(privs, start); err != nil {
			return nil, err
		}
		p.advance() // consume FROM
		users, err := p.parseUserList()
		if err != nil {
			return nil, err
		}
		stmt := &nodes.RevokeRoleStmt{
			Loc:      nodes.Loc{Start: start},
			IfExists: ifExists,
			Roles:    privNames(privs),
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
		stmt.Loc.End = p.prev.End
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
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	// Privilege revoke
	stmt := &nodes.RevokeStmt{Loc: nodes.Loc{Start: start}}
	stmt.IfExists = ifExists
	stmt.Privileges = privNames(privs)
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

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// privOrRole is one element of the list GRANT and REVOKE read before the ON,
// TO or FROM that settles what it was. A bare word can be either a privilege
// or a role, so the element carries which readings are open to it.
type privOrRole struct {
	name string
	// privOnly marks an element the role reading is closed to: a role is
	// one name, so a privilege written as two words is no role name.
	// oceanbase-ce:4.3.5 answers 1064 for `GRANT GRANT OPTION TO u1`, for `GRANT CREATE
	// VIEW TO u1` and for `GRANT SHOW DATABASES TO u1`.
	privOnly bool
}

// privNames drops the priv-or-role distinction, for the statement fields that
// hold plain names.
func privNames(elems []privOrRole) []string {
	if len(elems) == 0 {
		return nil
	}
	names := make([]string, len(elems))
	for i, e := range elems {
		names[i] = e.name
	}
	return names
}

// requireRoles rejects a list that holds an element only a privilege reading
// is open to.
func (p *Parser) requireRoles(elems []privOrRole, at int) error {
	for _, e := range elems {
		if e.privOnly {
			return &ParseError{
				Message:  "expected role name",
				Position: at,
			}
		}
	}
	return nil
}

// parsePrivilegeList parses the comma-separated list GRANT and REVOKE read
// before the ON, TO or FROM that settles whether it held privileges or roles.
// Returns the list and whether ALL [PRIVILEGES] was specified.
func (p *Parser) parsePrivilegeList() ([]privOrRole, bool, error) {
	if p.cur.Type == kwALL {
		p.advance()
		// Optional PRIVILEGES keyword
		if p.cur.Type == kwPRIVILEGES {
			p.advance()
		}
		return nil, true, nil
	}

	var privs []privOrRole
	for {
		elem, err := p.parsePrivOrRoleElem()
		if err != nil {
			return nil, false, err
		}
		privs = append(privs, elem)

		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}

	return privs, false, nil
}

// parsePrivOrRoleElem parses one element of that list.
//
// A privilege name is a keyword or an identifier; a role name is one of those
// or a string literal, and a role may carry an @host of its own.
// oceanbase-ce:4.3.5 runs `GRANT 'role1' ON *.* TO u1` and `GRANT
// role1@localhost ON *.* TO u1` as well as the TO form of each, so neither
// spelling is role-only the way it is on TiDB.
func (p *Parser) parsePrivOrRoleElem() (privOrRole, error) {
	if p.cur.Type == tokSCONST {
		name := p.cur.Str
		p.advance()
		host, hasHost, err := p.parseUserHostOpt()
		if err != nil {
			return privOrRole{}, err
		}
		if hasHost {
			name += "@" + host
		}
		return privOrRole{name: name}, nil
	}
	if p.isIdentToken() {
		if next := p.peekNext(); next.Type == tokIDENT && strings.HasPrefix(next.Str, "@") &&
			!strings.HasPrefix(p.inputText(next.Loc, next.End), "`") {
			name, _, err := p.parseIdent()
			if err != nil {
				return privOrRole{}, err
			}
			host, _, err := p.parseUserHostOpt()
			if err != nil {
				return privOrRole{}, err
			}
			return privOrRole{name: name + "@" + host}, nil
		}
	}
	name, word, err := p.parsePrivilegeName()
	if err != nil {
		return privOrRole{}, err
	}
	// GRANT on its own is no privilege name, and the word GRANT is no role
	// name either, so a bare GRANT is neither: the server answers 1064 for
	// `GRANT GRANT ON *.* TO u` and for `GRANT GRANT TO u`.
	if word.keyword && name == "GRANT" {
		return privOrRole{}, p.syntaxErrorAtCur()
	}
	return privOrRole{name: name, privOnly: word.multiToken}, nil
}

// privWord records what parsePrivilegeName read. The caller needs it because
// a privilege word and a name that happens to spell one are not the same
// thing: oceanbase-ce:4.3.5, mysql:8.0 and v8.5.5 all run `GRANT `grant` TO u` and `GRANT `create view` TO u`,
// where the backquoted word is a role name, and answer 1064 for `GRANT
// `create` VIEW ON *.* TO u`.
type privWord struct {
	// keyword is true when the first token was a keyword rather than an
	// identifier. Only a keyword starts a privilege word, which is why the
	// multi-word completions below are closed to identifiers.
	keyword bool
	// multiToken is true when the name spans more than one token. No role
	// name does; a backquoted name holding a space is still one token.
	multiToken bool
}

// parsePrivilegeName parses a single privilege name.
// Privilege names can be keywords like SELECT, INSERT, UPDATE, DELETE, etc.
// Multi-word: CREATE VIEW, ALTER ROUTINE, SHOW DATABASES, CREATE TEMPORARY TABLES,
// LOCK TABLES, REPLICATION CLIENT, REPLICATION SLAVE, CREATE TABLESPACE,
// CREATE USER, CREATE ROLE, DROP ROLE.
// Column-level: SELECT (col1, col2), INSERT (col1).
func (p *Parser) parsePrivilegeName() (string, privWord, error) {
	var name string
	var word privWord
	switch p.cur.Type {
	case kwSELECT, kwINSERT, kwUPDATE, kwDELETE, kwCREATE, kwDROP, kwALTER,
		kwINDEX, kwEXECUTE, kwGRANT, kwREFERENCES, kwTRIGGER, kwEVENT,
		kwSHOW, kwLOCK:
		tok := p.advance()
		name = strings.ToUpper(tok.Str)
		word.keyword = true
	case tokIDENT:
		tok := p.advance()
		name = strings.ToUpper(tok.Str)
	default:
		// Try accepting any keyword as a privilege name
		if p.cur.Type >= 700 {
			tok := p.advance()
			name = strings.ToUpper(tok.Str)
			word.keyword = true
		} else {
			return "", word, &ParseError{
				Message:  "expected privilege name",
				Position: p.cur.Loc,
			}
		}
	}

	// The multi-word privilege names. Every second word here is a keyword
	// in this lexer, so the completions are reached only from a keyword
	// first word: a backquoted `create` is a name, and the VIEW after it is
	// not part of it.
	if word.keyword {
		before := name
		name = p.completePrivilegeName(name)
		word.multiToken = name != before
	}
	return p.parsePrivilegeColumns(name, word)
}

// completePrivilegeName reads the second word of a two-word privilege name,
// where there is one.
func (p *Parser) completePrivilegeName(name string) string {
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
	case "GRANT":
		// GRANT OPTION is the privilege; GRANT on its own is not one, so
		// an unmatched OPTION leaves the name as GRANT and the caller
		// rejects it. The server agrees: `GRANT GRANT ON *.* TO u` is
		// 1064. REVOKE's `ALL PRIVILEGES, GRANT OPTION` form does not
		// come through here -- its caller reads the two words itself.
		if _, ok := p.match(kwOPTION); ok {
			name = "GRANT OPTION"
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
	return name
}

// parsePrivilegeColumns reads the optional column list of a column-level
// privilege: priv_type (col_list).
func (p *Parser) parsePrivilegeColumns(name string, word privWord) (string, privWord, error) {
	if p.cur.Type != '(' {
		return name, word, nil
	}
	p.advance()
	var cols []string
	for {
		col, _, err := p.parseIdent()
		if err != nil {
			return "", word, err
		}
		cols = append(cols, col)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if _, err := p.expect(')'); err != nil {
		return "", word, err
	}
	return name + "(" + strings.Join(cols, ", ") + ")", word, nil
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
				Loc:    nodes.Loc{Start: start, End: p.prev.End},
				Schema: "*",
				Name:   "*",
			}
			if p.cur.Type == '*' {
				p.advance()
				ref.Loc.End = p.prev.End
			} else {
				// *.tbl_name (unusual but handle)
				name, _, err := p.parseIdent()
				if err != nil {
					return nil, err
				}
				ref.Name = name
				ref.Loc.End = p.prev.End
			}
			target.Name = ref
		} else {
			// Just *
			target.Name = &nodes.TableRef{
				Loc:  nodes.Loc{Start: start, End: p.prev.End},
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

	target.Loc.End = p.prev.End
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

	ref.Loc.End = p.prev.End
	return ref, nil
}

// parseUserList parses a comma-separated list of user identifiers.
// User can be 'user'@'host' or just a name.
// parseGranteeList parses the grantee list of a GRANT, where each grantee may
// carry an auth option of its own: `GRANT ALL ON db.* TO 'u'@'localhost'
// IDENTIFIED BY '...'` is the spelling an old dump writes. MySQL 8.0 dropped
// the tail and oceanbase-ce:4.3.5 runs it, with the BY PASSWORD spelling too,
// on the privilege form only:
// oceanbase-ce:4.3.5 answers 1064 for `GRANT 'r' TO u IDENTIFIED BY '...'`,
// for the same tail on GRANT PROXY, and for `REVOKE SELECT ON *.* FROM u
// IDENTIFIED BY '...'`, so every other arm reads parseUserList instead.
//
// The AST models a grantee as a name. Nothing downstream reads a credential
// off a GRANT -- the catalog models no users and the review walk claims
// nothing from one -- so the auth option is parsed and dropped rather than
// stored.
func (p *Parser) parseGranteeList() ([]string, error) {
	var users []string
	for {
		user, err := p.parseUserName()
		if err != nil {
			return nil, err
		}
		users = append(users, user)

		if p.cur.Type == kwIDENTIFIED {
			if err := p.parseUserAuthOption(&nodes.UserSpec{}, true); err != nil {
				return nil, err
			}
		}

		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	return users, nil
}

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

// parseRoleRefList parses a comma-separated list of role references. It is
// parseUserList with CURRENT_USER refused: CURRENT_USER is a user and not a
// role, so oceanbase-ce:4.3.5 answers 1064 for it wherever a role is named --
// CREATE USER u DEFAULT ROLE current_user, CREATE/SET ROLE current_user, SET
// DEFAULT ROLE current_user TO u. DROP ROLE current_user is the one the engine
// reads, so it keeps parseUserList. A quoted 'CURRENT_USER' is an ordinary role
// name.
func (p *Parser) parseRoleRefList() ([]string, error) {
	var roles []string
	for {
		if p.cur.Type == kwCURRENT_USER {
			return nil, &ParseError{
				Message:  "expected role name",
				Position: p.cur.Loc,
			}
		}
		role, err := p.parseUserName()
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)

		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	return roles, nil
}

// parseUserName parses a user identifier, which may be 'user'@'host'.
func (p *Parser) parseUserName() (string, error) {
	name, host, hasHost, err := p.parseUserAndHost()
	if err != nil {
		return "", err
	}
	if hasHost {
		return name + "@" + host, nil
	}
	return name, nil
}

// parseUserAndHost parses a user name, then optionally the @host half,
// reporting whether one was written.
//
// Both halves are a string literal or an identifier, backquoted or not, and a
// reserved word is not one: oceanbase-ce:4.3.5 answers 1064 for
// CREATE USER select@localhost and runs CREATE USER comment@localhost.
//
// CURRENT_USER is the one reserved word a user accepts, as its own alternative
// `CURRENT_USER optional_braces`: it names the connected user, takes no @host,
// and allows an optional empty "()" like the function spelling. It reaches
// every statement that names a user through here, so all accept it where
// oceanbase-ce:4.3.5 does.
func (p *Parser) parseUserAndHost() (name, host string, hasHost bool, err error) {
	switch {
	case p.cur.Type == kwCURRENT_USER:
		p.advance() // consume CURRENT_USER
		if p.cur.Type == '(' {
			p.advance() // consume '('
			if _, err = p.expect(')'); err != nil {
				return "", "", false, err
			}
		}
		return "CURRENT_USER", "", false, nil
	case p.cur.Type == tokSCONST:
		name = p.cur.Str
		p.advance()
	case p.isIdentToken():
		name, _, err = p.parseIdent()
		if err != nil {
			return "", "", false, err
		}
	default:
		return "", "", false, &ParseError{
			Message:  "expected user name",
			Position: p.cur.Loc,
		}
	}
	host, hasHost, err = p.parseUserHostOpt()
	return name, host, hasHost, err
}

// atAtIdent reports whether the current token is an '@' the lexer scanned as
// the start of a name, with or without a host attached.
//
// The spelling has to be checked against the source and not the token's text,
// because scanBacktickIdent strips the backquotes: "`@localhost`" arrives as
// an identifier whose text is "@localhost", exactly like the "@localhost" of
// `u@localhost`, and it is not a host. oceanbase-ce:4.3.5 answers 1064 for
// `CREATE USER 'u' `@localhost“.
func (p *Parser) atAtIdent() bool {
	if p.cur.Type != tokIDENT || !strings.HasPrefix(p.cur.Str, "@") {
		return false
	}
	return !strings.HasPrefix(p.inputText(p.cur.Loc, p.cur.End), "`")
}

// parseUserHostOpt consumes the optional @host half of a user name.
//
// The lexer folds `@word` and "@`word`" into a single identifier token, so a
// host written without quotes arrives attached to its '@'. A quoted host
// arrives as a separate string literal, and there the '@' has to touch it:
// oceanbase-ce:4.3.5 answers 1064 for `CREATE USER u@ 'h'` and
// runs `CREATE USER u @'h'`, so the gap that matters is the one after the
// '@'.
// A '@' with nothing attached is an empty host, which the server takes.
func (p *Parser) parseUserHostOpt() (string, bool, error) {
	if !p.atAtIdent() {
		return "", false, nil
	}
	if len(p.cur.Str) > 1 {
		host := strings.TrimPrefix(p.cur.Str, "@")
		at := p.cur.Loc
		p.advance()
		return host, true, p.checkHostChars(host, at)
	}
	p.advance() // consume the bare '@'
	if p.cur.Loc != p.prev.End {
		return "", true, nil
	}
	if p.cur.Type == tokSCONST {
		host, at := p.cur.Str, p.cur.Loc
		p.advance()
		return host, true, p.checkHostChars(host, at)
	}
	if p.isIdentToken() {
		at := p.cur.Loc
		host, _, err := p.parseIdent()
		if err != nil {
			return host, true, err
		}
		return host, true, p.checkHostChars(host, at)
	}
	return "", true, nil
}

// checkHostChars rejects a host holding a character this engine's host
// production does not take. OceanBase's host is narrower than a quoted
// identifier: oceanbase-ce:4.3.5 runs `u@'a_b'`, `u@'a.b'`, `u@'a%b'`,
// `u@'a$b'`, `u@'a/b'`, `u@'127.0.0.1'`, a host of non-ASCII letters and an
// empty one, and answers 1064 for `u@'a b'`, `u@'a-b'`, `u@'a:b'`, `u@'a+b'`,
// `u@'a*b'` and `u@'@'` -- so `my-host.example.com` is no host at all there.
//
// The set below is what that server was measured to take rather than a
// published grammar, so a character outside it that the server does take is a
// gap in this parser and not a statement it wrongly runs.
func (p *Parser) checkHostChars(host string, at int) error {
	for _, r := range host {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
		case r == '_', r == '.', r == '%', r == '/', r == '$':
		case r > 0x7f:
		default:
			return &ParseError{
				Message:  "host name holds a character this engine does not take",
				Position: at,
			}
		}
	}
	return nil
}

// parseAuthPlugin parses the plugin name of an IDENTIFIED WITH clause, which
// is a string literal as well as an identifier: `IDENTIFIED WITH
// 'caching_sha2_password'` is the spelling a mysqldump writes.
func (p *Parser) parseAuthPlugin() (string, error) {
	if p.cur.Type == tokSCONST {
		name := p.cur.Str
		p.advance()
		return name, nil
	}
	if p.isIdentToken() {
		name, _, err := p.parseIdent()
		return name, err
	}
	return "", p.syntaxErrorAtCur()
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
func (p *Parser) parseCreateUserStmt(start int) (*nodes.CreateUserStmt, error) {
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
		spec, err := p.parseUserSpec(true)
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
			roles, err := p.parseRoleRefList()
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
	if err := p.parseUserAccountOptions(
		&stmt.PasswordExpire, &stmt.PasswordHistory, &stmt.PasswordReuseInterval,
		&stmt.PasswordRequireCurrent, &stmt.FailedLoginAttempts, &stmt.HasFailedLogin,
		&stmt.PasswordLockTime, &stmt.AccountLock, &stmt.AccountUnlock,
		&stmt.Comment, &stmt.Attribute,
	); err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseDropUserStmt parses a DROP USER statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/drop-user.html
//
//	DROP USER [IF EXISTS] user [, user] ...
func (p *Parser) parseDropUserStmt(start int) (*nodes.DropUserStmt, error) {
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

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterUserStmt parses an ALTER USER statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/alter-user.html
//
//	ALTER USER user_spec [, user_spec] ...
//	    [REQUIRE {NONE | tls_option [[AND] tls_option] ...}]
//	    [WITH resource_option [resource_option] ...]
func (p *Parser) parseAlterUserStmt(start int) (*nodes.AlterUserStmt, error) {
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
		userStart := p.pos()
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
			stmt.Loc.End = p.prev.End
			return stmt, nil
		}

		// user_func_auth_option: IDENTIFIED BY 'auth_string' [REPLACE ...] [RETAIN ...] | DISCARD OLD PASSWORD
		userSpec := &nodes.UserSpec{Loc: nodes.Loc{Start: userStart}, Name: "USER()"}
		if p.cur.Type == kwIDENTIFIED {
			if err := p.parseUserAuthOption(userSpec, false); err != nil {
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
		userSpec.Loc.End = p.prev.End
		stmt.Users = append(stmt.Users, userSpec)
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	// Check for ALTER USER user DEFAULT ROLE ...
	// We parse the first user, then check if DEFAULT follows.
	spec, err := p.parseUserSpec(false)
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
			stmt.Loc.End = p.prev.End
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
			stmt.Loc.End = p.prev.End
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
				roles, err := p.parseRoleRefList()
				if err != nil {
					return nil, err
				}
				stmt.DefaultRoles = roles
			}
			stmt.Loc.End = p.prev.End
			return stmt, nil
		}
	}

	// More user_specs
	for p.cur.Type == ',' {
		p.advance() // consume ','
		spec, err := p.parseUserSpec(false)
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
	if err := p.parseUserAccountOptions(
		&stmt.PasswordExpire, &stmt.PasswordHistory, &stmt.PasswordReuseInterval,
		&stmt.PasswordRequireCurrent, &stmt.FailedLoginAttempts, &stmt.HasFailedLogin,
		&stmt.PasswordLockTime, &stmt.AccountLock, &stmt.AccountUnlock,
		&stmt.Comment, &stmt.Attribute,
	); err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prev.End
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

	op.Loc.End = p.prev.End
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

		op.Loc.End = p.prev.End
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

// parseUserSpec parses a user specification. byPassword says whether this
// position takes the IDENTIFIED BY PASSWORD spelling, which
// oceanbase-ce:4.3.5 runs only where a user is created.
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
func (p *Parser) parseUserSpec(byPassword bool) (*nodes.UserSpec, error) {
	start := p.pos()
	spec := &nodes.UserSpec{Loc: nodes.Loc{Start: start}}

	name, host, _, err := p.parseUserAndHost()
	if err != nil {
		return nil, err
	}
	spec.Name, spec.Host = name, host

	// Optional IDENTIFIED BY/WITH
	if p.cur.Type == kwIDENTIFIED {
		if err := p.parseUserAuthOption(spec, byPassword); err != nil {
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
		if err := p.parseUserAuthOption(factor, byPassword); err != nil {
			return nil, err
		}
		factor.Loc.End = p.prev.End
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

	spec.Loc.End = p.prev.End
	return spec, nil
}

// parseCreateRoleStmt parses a CREATE ROLE statement.
//
//	CREATE ROLE [IF NOT EXISTS] role [, role] ...
func (p *Parser) parseCreateRoleStmt(start int) (*nodes.CreateRoleStmt, error) {
	p.advance() // consume ROLE

	stmt := &nodes.CreateRoleStmt{Loc: nodes.Loc{Start: start}}

	// [IF NOT EXISTS]
	if p.cur.Type == kwIF {
		p.advance()
		p.match(kwNOT)
		p.match(kwEXISTS_KW)
		stmt.IfNotExists = true
	}

	// role [, role] ... An unquoted role_ident refuses a word that is a name elsewhere, so
	// `CREATE ROLE event` is a 1064. CURRENT_USER is a user and not a role, so it is one too.
	for {
		if p.cur.Type == kwCURRENT_USER {
			return nil, p.syntaxErrorAtCur()
		}
		if isAmbiguousRoleWord(p.cur.Type) {
			return nil, p.syntaxErrorAtCur()
		}
		role, err := p.parseUserName()
		if err != nil {
			return nil, err
		}
		stmt.Roles = append(stmt.Roles, role)
		if _, ok := p.match(','); !ok {
			break
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseDropRoleStmt parses a DROP ROLE statement.
//
//	DROP ROLE [IF EXISTS] role [, role] ...
func (p *Parser) parseDropRoleStmt(start int) (*nodes.DropRoleStmt, error) {
	p.advance() // consume ROLE

	stmt := &nodes.DropRoleStmt{Loc: nodes.Loc{Start: start}}

	// [IF EXISTS]
	if p.cur.Type == kwIF {
		p.advance()
		p.match(kwEXISTS_KW)
		stmt.IfExists = true
	}

	// role [, role] ... oceanbase-ce:4.3.5 reads DROP ROLE current_user (unlike
	// CREATE/SET ROLE), so this list keeps parseUserList and accepts CURRENT_USER.
	roles, err := p.parseUserList()
	if err != nil {
		return nil, err
	}
	stmt.Roles = roles

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseUserAuthOption parses the IDENTIFIED portion of a user spec. Caller
// has verified p.cur.Type == kwIDENTIFIED.
//
//	IDENTIFIED BY 'password'
//	IDENTIFIED BY PASSWORD 'hash'                 (CREATE USER only)
//	IDENTIFIED WITH auth_plugin BY 'password'
//	IDENTIFIED WITH auth_plugin BY PASSWORD 'hash'
//	IDENTIFIED WITH auth_plugin AS 'hash'
//
// BY PASSWORD names the stored hash rather than the password, after a plugin
// as well as on its own. MySQL 8.0 dropped both spellings;
// oceanbase-ce:4.3.5, which reports itself as 5.7.25, runs them -- but only
// where a user is created: it answers 1064 for `ALTER USER u@'h' IDENTIFIED
// BY PASSWORD 'hash'` while running the CREATE USER of the same clause, which
// is what byPassword carries.
//
// A plugin is never the whole clause: oceanbase-ce:4.3.5 answers 1064 for
// `IDENTIFIED WITH p` and for `IDENTIFIED WITH 'p'`, so a BY or an AS has to
// follow. The hash is a string literal and not a hex one.
//
// RANDOM PASSWORD is MySQL 8.0's, but it is left where it was:
// oceanbase-ce:4.3.5 answers 1064 for `IDENTIFIED BY RANDOM PASSWORD`, so
// taking it is a false accept, and one to retire on its own rather than
// alongside the clauses this change adds.
// A bare IDENTIFIED is not a clause: oceanbase-ce:4.3.5 answers
// 1064 for `CREATE USER u@'h' IDENTIFIED` and for the same tail on a GRANT
// grantee, so an arm has to follow.
func (p *Parser) parseUserAuthOption(spec *nodes.UserSpec, byPassword bool) error {
	p.advance() // consume IDENTIFIED
	switch p.cur.Type {
	case kwBY:
		p.advance()
		if err := p.parseAuthByTail(spec, byPassword); err != nil {
			return err
		}
	case kwWITH:
		p.advance()
		plugin, err := p.parseAuthPlugin()
		if err != nil {
			return err
		}
		spec.AuthPlugin = plugin
		switch p.cur.Type {
		case kwBY:
			p.advance()
			if err := p.parseAuthByTail(spec, true); err != nil {
				return err
			}
		case kwAS:
			p.advance()
			if p.cur.Type != tokSCONST {
				return p.syntaxErrorAtCur()
			}
			spec.AuthHash = p.cur.Str
			p.advance()
		default:
			return p.syntaxErrorAtCur()
		}
	default:
		return p.syntaxErrorAtCur()
	}
	return nil
}

// parseAuthByTail parses what follows a BY in an auth option. Caller has
// consumed the BY, and byPassword says whether the PASSWORD spelling is
// allowed in this position.
func (p *Parser) parseAuthByTail(spec *nodes.UserSpec, byPassword bool) error {
	if byPassword && p.cur.Type == kwPASSWORD {
		p.advance()
		if p.cur.Type != tokSCONST {
			return p.syntaxErrorAtCur()
		}
		spec.AuthHash = p.cur.Str
		p.advance()
		return nil
	}
	if p.cur.Type == kwRANDOM {
		p.advance() // consume RANDOM
		p.advance() // consume PASSWORD
		spec.PasswordRandom = true
		return nil
	}
	if p.cur.Type != tokSCONST {
		return p.syntaxErrorAtCur()
	}
	spec.Password = p.cur.Str
	p.advance()
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
		roles, err := p.parseRoleRefList()
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

	stmt.Loc.End = p.prev.End
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
			roles, err := p.parseRoleRefList()
			if err != nil {
				return nil, err
			}
			stmt.AllExcept = roles
		}
	default:
		roles, err := p.parseRoleRefList()
		if err != nil {
			return nil, err
		}
		stmt.Roles = roles
	}

	stmt.Loc.End = p.prev.End
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

		pair.Loc.End = p.prev.End
		stmt.Pairs = append(stmt.Pairs, pair)

		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseRenameUserPart parses a user part (name[@host]) in a RENAME USER statement.
// Supports both quoted and unquoted user names.
func (p *Parser) parseRenameUserPart() (string, string, error) {
	// oceanbase-ce:4.3.5 answers 1064 for CURRENT_USER as a RENAME USER target
	// (measured), while accepting it as a user in CREATE/ALTER/DROP USER,
	// GRANT/REVOKE and SET DEFAULT ROLE. The source half is refused with it, so
	// neither half of a rename takes the keyword.
	if p.cur.Type == kwCURRENT_USER {
		return "", "", &ParseError{
			Message:  "CURRENT_USER is not allowed in RENAME USER",
			Position: p.cur.Loc,
		}
	}
	name, host, _, err := p.parseUserAndHost()
	if err != nil {
		return "", "", err
	}
	return name, host, nil
}

// parseRequireClause parses a REQUIRE clause for TLS options.
//
//	REQUIRE NONE
//	REQUIRE SSL
//	REQUIRE X509
//	REQUIRE tls_option [[AND] tls_option] ...
//
//	tls_option: {
//	    CIPHER 'cipher'
//	  | ISSUER 'issuer'
//	  | SUBJECT 'subject'
//	}
//
// NONE, SSL and X509 each stand alone and join no list, so `REQUIRE SSL AND
// X509`, `REQUIRE X509 ISSUER 'i'` and `REQUIRE NONE SSL` are all 1064 on
// oceanbase-ce:4.3.5. The list is never empty, every element
// carries a string literal, and an AND is never trailing.
func (p *Parser) parseRequireClause() (*nodes.RequireClause, error) {
	start := p.pos()
	p.advance() // consume REQUIRE

	req := &nodes.RequireClause{Loc: nodes.Loc{Start: start}}

	switch p.cur.Type {
	case kwNONE:
		req.None = true
	case kwSSL:
		req.SSL = true
	case kwX509:
		req.X509 = true
	}
	if req.None || req.SSL || req.X509 {
		p.advance()
		req.Loc.End = p.prev.End
		return req, nil
	}

	// An element is mandatory to start with, and again after an AND; between
	// two elements the AND is optional.
	for required := true; ; {
		var out *string
		switch p.cur.Type {
		case kwCIPHER:
			out = &req.Cipher
		case kwISSUER:
			out = &req.Issuer
		case kwSUBJECT:
			out = &req.Subject
		default:
			if required {
				return nil, p.syntaxErrorAtCur()
			}
			req.Loc.End = p.prev.End
			return req, nil
		}
		p.advance() // consume the option keyword
		if p.cur.Type != tokSCONST {
			return nil, p.syntaxErrorAtCur()
		}
		*out = p.cur.Str
		p.advance()

		required = p.cur.Type == kwAND
		if required {
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
		res.Loc.End = p.prev.End
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
) error {
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
				return nil
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

		// COMMENT and ATTRIBUTE take a quoted string; the server answers 1064
		// without it.
		case p.cur.Type == kwCOMMENT:
			p.advance()
			if p.cur.Type != tokSCONST {
				return p.syntaxErrorAtCur()
			}
			*comment = p.cur.Str
			p.advance()

		case p.cur.Type == kwATTRIBUTE:
			p.advance()
			if p.cur.Type != tokSCONST {
				return p.syntaxErrorAtCur()
			}
			*attribute = p.cur.Str
			p.advance()

		default:
			return nil
		}
	}
}
