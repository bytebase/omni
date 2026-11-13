package parser

import (
	"math"
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// parseGrantStmt parses a GRANT statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/grant.html
//
//	GRANT priv_type [, priv_type] ... ON [object_type] priv_level TO user [, user] ... [WITH GRANT OPTION]
//	GRANT ALL [PRIVILEGES] ON [object_type] priv_level TO user [, user] ... [WITH GRANT OPTION]
//	GRANT role [, role] ... TO user [, user] ... [WITH ADMIN OPTION]
//
// MySQL 8.0 ends a privilege grant with AS user [WITH ROLE ...]. MariaDB has no such
// clause: grant_command ends at opt_grant_options in MariaDB 11.8 sql/sql_yacc.yy,
// and the server answers 1064 for it.
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
		// Role grant: GRANT role TO user [, user] ... [WITH ADMIN OPTION]
		if allPriv {
			return nil, p.syntaxErrorAtCur()
		}
		if err := p.requireSingleRole(privs, start); err != nil {
			return nil, err
		}
		if err := p.requireRoles(privs, start); err != nil {
			return nil, err
		}
		p.advance() // consume TO
		users, err := p.parseGranteeList()
		if err != nil {
			return nil, err
		}
		stmt := &nodes.GrantRoleStmt{
			Loc:   nodes.Loc{Start: start},
			Roles: privNames(privs),
			To:    users,
		}
		// [WITH ADMIN OPTION]
		stmt.WithAdmin, err = p.parseWithOption(kwADMIN)
		if err != nil {
			return nil, err
		}
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	// Privilege grant
	if err := p.requirePrivileges(privs, start); err != nil {
		return nil, err
	}
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
		users, err := p.parseGranteeList()
		if err != nil {
			return nil, err
		}
		stmt.To = users
		// [WITH GRANT OPTION]
		stmt.WithGrant, err = p.parseWithOption(kwGRANT)
		if err != nil {
			return nil, err
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

	// [REQUIRE {NONE | tls_option [[AND] tls_option] ...}]
	if p.cur.Type == kwREQUIRE {
		req, err := p.parseRequireClause()
		if err != nil {
			return nil, err
		}
		stmt.Require = req
	}

	// [WITH grant_option ...]: GRANT OPTION and the resource options, in any
	// order, repeated or not, and at least one of them. MariaDB 11.8 parses `WITH
	// GRANT OPTION GRANT OPTION` and `WITH MAX_QUERIES_PER_HOUR 1 GRANT OPTION`,
	// and answers 1064 for a bare `WITH` and for `WITH GRANT`.
	//
	// Ref: opt_grant_options, grant_option_list and grant_option in MariaDB 11.8
	// sql/sql_yacc.yy.
	if _, ok := p.match(kwWITH); ok {
		for first := true; ; first = false {
			if _, ok := p.match(kwGRANT); ok {
				if _, err := p.expect(kwOPTION); err != nil {
					return nil, err
				}
				stmt.WithGrant = true
				continue
			}
			if !p.startsResourceOption() {
				if first {
					return nil, p.syntaxErrorAtCur()
				}
				break
			}
			if stmt.Resource == nil {
				stmt.Resource = &nodes.ResourceOption{Loc: nodes.Loc{Start: p.pos()}}
			}
			if err := p.parseResourceOption(stmt.Resource); err != nil {
				return nil, err
			}
			stmt.Resource.Loc.End = p.prev.End
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseWithOption reads `WITH word OPTION`, the clause that closes GRANT PROXY (GRANT
// OPTION) and a role grant (ADMIN OPTION). It reports whether the clause is there. A WITH
// takes both words: MariaDB 11.8 answers 1064 for `GRANT PROXY ON a TO b WITH`, `GRANT
// PROXY ON a TO b WITH GRANT` and `GRANT r TO u WITH ADMIN`. A privilege grant's WITH is
// a list instead (opt_grant_options), which parseGrantStmt reads.
//
// Ref: opt_grant_option and opt_with_admin_option in MariaDB 11.8 sql/sql_yacc.yy.
func (p *Parser) parseWithOption(word int) (bool, error) {
	if _, ok := p.match(kwWITH); !ok {
		return false, nil
	}
	if _, err := p.expect(word); err != nil {
		return false, err
	}
	if _, err := p.expect(kwOPTION); err != nil {
		return false, err
	}
	return true, nil
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
			if _, err := p.expect(kwOPTION); err != nil {
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
		// Role revoke: REVOKE role FROM user [, user] ...
		if err := p.requireSingleRole(privs, start); err != nil {
			return nil, err
		}
		if allPriv {
			return nil, p.syntaxErrorAtCur()
		}
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
	if err := p.requirePrivileges(privs, start); err != nil {
		return nil, err
	}
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
	// roleOnly marks an element the privilege reading is closed to.
	roleOnly bool
	// privOnly marks an element the role reading is closed to: a role is
	// one name, so a privilege written as two words is no role name.
	// mariadb:11.8 answers 1064 for `GRANT GRANT OPTION TO u1`, for `GRANT CREATE
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

// requirePrivileges rejects a list that holds an element only a role reading
// is open to.
func (p *Parser) requirePrivileges(elems []privOrRole, at int) error {
	for _, e := range elems {
		if e.roleOnly {
			return &ParseError{
				Message:  "expected privilege name",
				Position: at,
			}
		}
	}
	return nil
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
// or a string literal, written on its own. A MariaDB role carries no host and
// a role grant names one role: mariadb:11.8 answers 1064 for `GRANT
// 'role1'@'%' TO u1`, for `GRANT role1, role2 TO u1` and for `REVOKE 'r1',
// 'r2' FROM u1`, the last two of which requireSingleRole reports. A quoted
// name is a role and nothing else, so `GRANT 'role1' ON *.* TO u1` is 1064
// too, where mysql:8.0 runs it.
func (p *Parser) parsePrivOrRoleElem() (privOrRole, error) {
	if p.cur.Type == tokSCONST {
		name := p.cur.Str
		p.advance()
		return privOrRole{name: name, roleOnly: true}, nil
	}
	name, word, err := p.parsePrivilegeName()
	if err != nil {
		return privOrRole{}, err
	}
	// A static privilege word is a privilege and no role, and neither is a privilege
	// with columns: role_or_privilege gives each its own arm, and the action for TO
	// answers 1064 for any of them, where a backquoted or dynamic name is a role.
	if (word.keyword && staticPrivilegeWord(name)) || word.columns {
		word.multiToken = true
	}
	// Nor is a reserved word a role name: 11.8.8 answers 1064 for `GRANT READ TO u`, where
	// `GRANT BINLOG TO u` grants a role.
	if word.keyword && !word.multiToken && IsReservedKeyword(name) {
		word.multiToken = true
	}
	// GRANT on its own is no privilege name, and the word GRANT is no role
	// name either, so a bare GRANT is neither: the server answers 1064 for
	// `GRANT GRANT ON *.* TO u` and for `GRANT GRANT TO u`.
	if word.keyword && name == "GRANT" {
		return privOrRole{}, p.syntaxErrorAtCur()
	}
	return privOrRole{name: name, privOnly: word.multiToken, roleOnly: !isPrivilegeName(name, word)}, nil
}

// isPrivilegeName reports a privilege word object_privilege lists. Any other name is no
// privilege, so `GRANT r1 ON *.* TO u` is a 1064 where `GRANT r1 TO u` runs, and a name that
// spells one is no privilege either: `GRANT `select` ON *.* TO u`.
//
// Ref: object_privilege in MariaDB 11.8 sql/sql_yacc.yy.
func isPrivilegeName(name string, word privWord) bool {
	if !word.keyword && !word.multiToken {
		return false
	}
	if base, _, ok := strings.Cut(name, " ("); ok {
		name = base
	}
	if base, _, ok := strings.Cut(name, "("); ok {
		name = base
	}
	switch name {
	case "SELECT", "INSERT", "UPDATE", "DELETE", "CREATE", "DROP", "RELOAD", "SHUTDOWN", "PROCESS", "FILE",
		"GRANT OPTION", "REFERENCES", "INDEX", "ALTER", "SHOW DATABASES", "SUPER", "CREATE TEMPORARY TABLES",
		"LOCK TABLES", "EXECUTE", "REPLICATION SLAVE", "REPLICATION CLIENT", "BINLOG MONITOR", "CREATE VIEW",
		"SHOW VIEW", "CREATE ROUTINE", "ALTER ROUTINE", "CREATE USER", "EVENT", "TRIGGER", "CREATE TABLESPACE",
		"DELETE HISTORY", "SET USER", "FEDERATED ADMIN", "CONNECTION ADMIN", "READ_ONLY ADMIN", "READ ONLY ADMIN",
		"REPLICATION SLAVE ADMIN", "REPLICATION MASTER ADMIN", "BINLOG ADMIN", "BINLOG REPLAY", "SLAVE MONITOR",
		"REPLICA MONITOR", "SHOW CREATE ROUTINE", "USAGE", "PROXY":
		return true
	}
	return false
}

// requireSingleRole rejects a role list with more than one element. A
// MariaDB role grant names one role where a MySQL one names a list:
// mariadb:11.8 answers 1064 for `GRANT role1, role2 TO u1` and for `REVOKE
// 'r1', 'r2' FROM u1`, and runs each with one role and several grantees.
func (p *Parser) requireSingleRole(roles []privOrRole, at int) error {
	if len(roles) <= 1 {
		return nil
	}
	return &ParseError{
		Message:  "a role grant names one role",
		Position: at,
	}
}

// staticPrivilegeWord reports a reserved word that is a privilege, which no role name can be.
//
// Ref: object_privilege and role_name in MariaDB 11.8 sql/sql_yacc.yy.
func staticPrivilegeWord(name string) bool {
	switch name {
	case "SELECT", "INSERT", "UPDATE", "DELETE", "CREATE", "DROP", "REFERENCES", "INDEX", "ALTER", "SHOW", "LOCK", "USAGE", "TRIGGER":
		return true
	}
	return false
}

// parseUserLimit reads a number of a user's password or lock options, which the action
// holds to a range with 1525, ER_WRONG_VALUE. A number past 64 bits is the 1064 it is on
// any other integer here; both are refused while the statement parses.
//
// Ref: connect_option and the PASSWORD and account lock options in MariaDB 11.8
// sql/sql_yacc.yy.
func (p *Parser) parseUserLimit(minimum, maximum int64, what string) (string, error) {
	if p.cur.Type != tokICONST {
		return "", p.syntaxErrorAtCur()
	}
	if _, err := strconv.ParseInt(p.cur.Str, 10, 64); err != nil {
		return "", p.syntaxErrorAtCur()
	}
	if p.cur.Ival < minimum || p.cur.Ival > maximum {
		return "", &ParseError{Message: "Incorrect " + what + " value: '" + p.cur.Str + "'", Position: p.cur.Loc}
	}
	return p.advance().Str, nil
}

// privWord records what parsePrivilegeName read. The caller needs it because
// a privilege word and a name that happens to spell one are not the same
// thing: mariadb:11.8, mysql:8.0 and v8.5.5 all run `GRANT `grant` TO u` and `GRANT `create view` TO u`,
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
	// columns is true when a column list follows the name. It is read off the tokens, since a
	// backquoted role may hold a parenthesis in its name.
	columns bool
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
	if (!word.multiToken || name == "REPLICATION SLAVE") && !p.backtickQuotedAt(p.prev.Loc) {
		before := name
		name = p.completeAdminPrivilege(name)
		word.multiToken = word.multiToken || name != before
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

// completeAdminPrivilege reads the rest of a privilege MariaDB spells in words of its own.
// The words are matched as written, plain, since they are names in every other place.
//
// Ref: object_privilege in MariaDB 11.8 sql/sql_yacc.yy.
func (p *Parser) completeAdminPrivilege(name string) string {
	// next reads the words only all together: on a partial match it reads none, so `READ ONLY TO`
	// stops at ONLY and is the 1064 11.8.8 answers, rather than a role READ with ONLY dropped.
	next := func(words ...string) bool {
		saved := p.saveScan()
		for _, word := range words {
			if !strings.EqualFold(p.cur.Str, word) || p.backtickQuotedAt(p.cur.Loc) {
				p.restoreScan(saved)
				return false
			}
			p.advance()
		}
		return true
	}
	switch name {
	case "BINLOG":
		for _, word := range []string{"ADMIN", "MONITOR", "REPLAY"} {
			if next(word) {
				return "BINLOG " + word
			}
		}
	case "CONNECTION", "FEDERATED", "READ_ONLY":
		if next("ADMIN") {
			return name + " ADMIN"
		}
	case "DELETE":
		if next("HISTORY") {
			return "DELETE HISTORY"
		}
	case "READ":
		if next("ONLY", "ADMIN") {
			return "READ ONLY ADMIN"
		}
	case "REPLICA", "SLAVE":
		if next("MONITOR") {
			return name + " MONITOR"
		}
	case "REPLICATION":
		if next("MASTER", "ADMIN") {
			return "REPLICATION MASTER ADMIN"
		}
	case "REPLICATION SLAVE":
		if next("ADMIN") {
			return "REPLICATION SLAVE ADMIN"
		}
	case "SET":
		if next("USER") {
			return "SET USER"
		}
	case "SHOW":
		if next("CREATE", "ROUTINE") {
			return "SHOW CREATE ROUTINE"
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
	word.columns = true
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

	// grant_ident holds a table to check_table_name and its database to check_db_name. A
	// `*` after the database leaves it alone, which is why `a `.* runs.
	if p.cur.Type != '.' {
		if identifierNameRefused(name) {
			return nil, &ParseError{Message: "Incorrect table name '" + name + "'", Position: start}
		}
	} else if p.peekNext().Type != '*' {
		if err := refuseDatabaseName(name, start); err != nil {
			return nil, err
		}
	}

	if p.cur.Type == '.' {
		p.advance() // consume '.'
		if p.cur.Type == '*' {
			p.advance()
			ref.Schema = name
			ref.Name = "*"
		} else {
			nameStart := p.pos()
			name2, _, err := p.parseIdent()
			if err != nil {
				return nil, err
			}
			if identifierNameRefused(name2) {
				return nil, &ParseError{Message: "Incorrect table name '" + name2 + "'", Position: nameStart}
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
// the tail and mariadb:11.8 runs it, on the role and PROXY forms as well as
// the privilege one. REVOKE does not take it: mariadb:11.8 answers 1064 for
// `REVOKE SELECT ON *.* FROM u IDENTIFIED BY '...'`, so the REVOKE arms read
// parseUserList instead.
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
			if err := p.parseUserAuthOption(&nodes.UserSpec{}); err != nil {
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
// role, so mariadb:11.8 answers 1064 for it wherever a role is named --
// CREATE USER u DEFAULT ROLE current_user, CREATE/DROP/SET ROLE current_user.
// A quoted 'CURRENT_USER' is an ordinary role name.
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
// reserved word is not one: mariadb:11.8 answers 1064 for
// CREATE USER select@localhost and runs CREATE USER comment@localhost.
//
// CURRENT_USER is the one reserved word a user accepts, as its own alternative
// `CURRENT_USER optional_braces`: it names the connected user, takes no @host,
// and allows an optional empty "()" like the function spelling. It reaches
// every statement that names a user through here, so all accept it where
// mariadb:11.8 does.
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
// `u@localhost`, and it is not a host. mariadb:11.8 answers 1064 for
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
// mariadb:11.8 answers 1064 for `CREATE USER u@ 'h'` and runs
// `CREATE USER u @'h'`, so the gap that matters is the one after the '@'.
// A '@' with nothing attached is an empty host, which the server takes.
func (p *Parser) parseUserHostOpt() (string, bool, error) {
	if !p.atAtIdent() {
		return "", false, nil
	}
	if len(p.cur.Str) > 1 {
		host := strings.TrimPrefix(p.cur.Str, "@")
		p.advance()
		return host, true, nil
	}
	p.advance() // consume the bare '@'
	if p.cur.Loc != p.prev.End {
		return "", true, nil
	}
	if p.cur.Type == tokSCONST {
		host := p.cur.Str
		p.advance()
		return host, true, nil
	}
	if p.isIdentToken() {
		host, _, err := p.parseIdent()
		return host, true, err
	}
	return "", true, nil
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
// Ref: https://mariadb.com/kb/en/create-user/
//
//	CREATE [OR REPLACE] USER [IF NOT EXISTS]
//	    user [auth_option] [, user [auth_option]] ...
//	    DEFAULT ROLE role [, role ] ...
//	    [REQUIRE {NONE | tls_option [[AND] tls_option] ...}]
//	    [WITH resource_option [resource_option] ...]
//	    [password_option | lock_option] ...
//	    [COMMENT 'comment_string' | ATTRIBUTE 'json_object']
func (p *Parser) parseCreateUserStmt(start int, orReplace bool) (*nodes.CreateUserStmt, error) {
	p.advance() // consume USER

	stmt := &nodes.CreateUserStmt{Loc: nodes.Loc{Start: start}, OrReplace: orReplace}

	ifNotExists, err := p.parseCreateIfNotExists(orReplace)
	if err != nil {
		return nil, err
	}
	stmt.IfNotExists = ifNotExists

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
		res, err := p.parseResourceOptions()
		if err != nil {
			return nil, err
		}
		stmt.Resource = res
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
		userSpec.Loc.End = p.prev.End
		stmt.Users = append(stmt.Users, userSpec)
		stmt.Loc.End = p.prev.End
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
		res, err := p.parseResourceOptions()
		if err != nil {
			return nil, err
		}
		stmt.Resource = res
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

	name, host, _, err := p.parseUserAndHost()
	if err != nil {
		return nil, err
	}
	spec.Name, spec.Host = name, host

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
//	CREATE [OR REPLACE] ROLE [IF NOT EXISTS] role [, role] ...
func (p *Parser) parseCreateRoleStmt(start int, orReplace bool) (*nodes.CreateRoleStmt, error) {
	p.advance() // consume ROLE

	stmt := &nodes.CreateRoleStmt{Loc: nodes.Loc{Start: start}, OrReplace: orReplace}

	ifNotExists, err := p.parseCreateIfNotExists(orReplace)
	if err != nil {
		return nil, err
	}
	stmt.IfNotExists = ifNotExists

	// role [, role] ...
	roles, err := p.parseRoleRefList()
	if err != nil {
		return nil, err
	}
	stmt.Roles = roles

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

	// role [, role] ...
	roles, err := p.parseRoleRefList()
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
//	IDENTIFIED BY PASSWORD 'hash'
//	IDENTIFIED WITH auth_plugin [AS 'hash']
//
// A plugin takes no password: mariadb:11.8 answers 1064 for `IDENTIFIED WITH
// p BY 'pw'`, which MySQL 8.0 runs, and the hash it takes is a string literal
// and not a hex one. BY PASSWORD, which MySQL 8.0 dropped, it runs.
//
// RANDOM PASSWORD is MySQL 8.0's and not MariaDB's, but it is left where it
// was: mariadb:11.8 answers 1064 for `IDENTIFIED BY RANDOM PASSWORD`, so
// taking it is a false accept, and one to retire on its own rather than
// alongside the clauses this change adds.
// A bare IDENTIFIED is not a clause: mariadb:11.8 answers
// 1064 for `CREATE USER u@'h' IDENTIFIED` and for the same tail on a GRANT
// grantee, so an arm has to follow.
func (p *Parser) parseUserAuthOption(spec *nodes.UserSpec) error {
	p.advance() // consume IDENTIFIED
	switch p.cur.Type {
	case kwBY:
		p.advance()
		// BY PASSWORD names the stored hash rather than the password.
		if p.cur.Type == kwPASSWORD {
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
	case kwWITH:
		p.advance()
		plugin, err := p.parseAuthPlugin()
		if err != nil {
			return err
		}
		spec.AuthPlugin = plugin
		if p.cur.Type == kwAS {
			p.advance()
			if p.cur.Type != tokSCONST {
				return p.syntaxErrorAtCur()
			}
			spec.AuthHash = p.cur.Str
			p.advance()
		}
	default:
		return p.syntaxErrorAtCur()
	}
	return nil
}

// parseSetDefaultRoleStmt parses a SET DEFAULT ROLE statement.
//
//	SET DEFAULT ROLE {role | NONE} [FOR user]
//
// MariaDB's form is not MySQL's: it names one role rather than a list, has no
// ALL, and says FOR where MySQL says TO, with the user optional because the
// statement defaults to the current one. mariadb:11.8 runs `SET DEFAULT ROLE
// role1`, `SET DEFAULT ROLE NONE` and `SET DEFAULT ROLE 'role1' FOR
// u@localhost`, and answers 1064 for `SET DEFAULT ROLE role1 TO u`, for `SET
// DEFAULT ROLE role1, role2 FOR u` and for `SET DEFAULT ROLE ALL FOR u`.
func (p *Parser) parseSetDefaultRoleStmt(start int) (*nodes.SetDefaultRoleStmt, error) {
	// DEFAULT and ROLE already consumed
	stmt := &nodes.SetDefaultRoleStmt{Loc: nodes.Loc{Start: start}}

	if p.cur.Type == kwNONE {
		stmt.None = true
		p.advance()
	} else {
		// The role operand is a Rolename, not a Username: CURRENT_USER is a user
		// and not a role, so mariadb:11.8 answers 1064 for SET DEFAULT ROLE
		// current_user (the quoted 'CURRENT_USER' stays an ordinary role name).
		// The FOR user below is a Username and does read CURRENT_USER.
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
		stmt.Roles = []string{role}
	}

	// [FOR user]
	if p.cur.Type == kwFOR {
		p.advance()
		user, err := p.parseUserName()
		if err != nil {
			return nil, err
		}
		stmt.To = []string{user}
	}

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
		if _, err := p.expect(kwTO); err != nil {
			return nil, err
		}

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
// mariadb:11.8. The list is never empty, every element
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
	// two elements the AND is optional. Each comes once (ER_DUP_ARGUMENT).
	seen := map[int]bool{}
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
		if seen[p.cur.Type] {
			return nil, &ParseError{Message: "Option '" + strings.ToUpper(p.cur.Str) + "' used twice in statement", Position: p.cur.Loc}
		}
		seen[p.cur.Type] = true
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
func (p *Parser) startsResourceOption() bool {
	return p.cur.Type == kwMAX_QUERIES_PER_HOUR || p.cur.Type == kwMAX_UPDATES_PER_HOUR ||
		p.cur.Type == kwMAX_CONNECTIONS_PER_HOUR || p.cur.Type == kwMAX_USER_CONNECTIONS
}

func (p *Parser) parseResourceOptions() (*nodes.ResourceOption, error) {
	// WITH takes at least one: `CREATE USER u WITH` is a 1064.
	if !p.startsResourceOption() {
		return nil, p.syntaxErrorAtCur()
	}
	res := &nodes.ResourceOption{Loc: nodes.Loc{Start: p.pos()}}
	for p.startsResourceOption() {
		if err := p.parseResourceOption(res); err != nil {
			return nil, err
		}
	}
	res.Loc.End = p.prev.End
	return res, nil
}

// parseResourceOption reads the resource option at the current token and its
// count into res.
func (p *Parser) parseResourceOption(res *nodes.ResourceOption) error {
	option := p.advance().Type
	if option == kwMAX_USER_CONNECTIONS {
		// NUM, where the other three are ulong_num and take any number:
		// mariadb:11.8.8 answers 1064 for MAX_USER_CONNECTIONS 2147483648.
		if !p.cur.isNum() {
			return p.syntaxErrorAtCur()
		}
	} else if p.cur.Type != tokICONST && p.cur.Type != tokFCONST {
		return p.syntaxErrorAtCur()
	}
	n := int(p.advance().Ival)
	switch option {
	case kwMAX_QUERIES_PER_HOUR:
		res.MaxQueriesPerHour = n
		res.HasMaxQueries = true
	case kwMAX_UPDATES_PER_HOUR:
		res.MaxUpdatesPerHour = n
		res.HasMaxUpdates = true
	case kwMAX_CONNECTIONS_PER_HOUR:
		res.MaxConnectionsPerHour = n
		res.HasMaxConnections = true
	case kwMAX_USER_CONNECTIONS:
		res.MaxUserConnections = n
		res.HasMaxUserConn = true
	}
	return nil
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
					// The day count is NUM: mariadb:11.8.8 answers 1064 for INTERVAL
					// 2147483648 DAY.
					n, err := p.parseUserLimit(1, math.MaxInt32, "PASSWORD EXPIRE INTERVAL")
					if err != nil {
						return err
					}
					if _, err := p.expect(kwDAY); err != nil {
						return err
					}
					*passwordExpire = "INTERVAL " + n + " DAY"
				} else {
					*passwordExpire = "EXPIRE"
				}
			case next.Type == kwHISTORY, next.Type == kwREUSE, next.Type == kwREQUIRE:
				// MariaDB's PASSWORD takes EXPIRE and nothing else.
				return p.syntaxErrorAtCur()
			case next.Type == kwHISTORY:
				p.advance() // consume PASSWORD
				p.advance() // consume HISTORY
				if p.cur.Type == kwDEFAULT {
					*passwordHistory = "DEFAULT"
					p.advance()
				} else if p.cur.Type == tokICONST {
					*passwordHistory = strconv.FormatInt(p.cur.Ival, 10)
					p.advance()
				} else {
					return p.syntaxErrorAtCur()
				}
			case next.Type == kwREUSE:
				p.advance() // consume PASSWORD
				p.advance() // consume REUSE
				// INTERVAL
				if _, err := p.expect(kwINTERVAL); err != nil {
					return err
				}
				if p.cur.Type == kwDEFAULT {
					*passwordReuseInterval = "DEFAULT"
					p.advance()
				} else if p.cur.Type == tokICONST {
					n := strconv.FormatInt(p.cur.Ival, 10)
					p.advance()
					if _, err := p.expect(kwDAY); err != nil {
						return err
					}
					*passwordReuseInterval = n + " DAY"
				} else {
					return p.syntaxErrorAtCur()
				}
			case next.Type == kwREQUIRE:
				p.advance() // consume PASSWORD
				p.advance() // consume REQUIRE
				if p.cur.Type != kwCURRENT {
					return p.syntaxErrorAtCur()
				}
				{
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

		case p.cur.Type == kwFAILED_LOGIN_ATTEMPTS, p.cur.Type == kwPASSWORD_LOCK_TIME:
			// No such option on MariaDB.
			return p.syntaxErrorAtCur()

		case p.cur.Type == kwFAILED_LOGIN_ATTEMPTS:
			p.advance()
			n, err := p.parseUserLimit(0, 32767, "FAILED_LOGIN_ATTEMPTS")
			if err != nil {
				return err
			}
			*failedLoginAttempts, _ = strconv.Atoi(n)
			*hasFailedLogin = true

		case p.cur.Type == kwPASSWORD_LOCK_TIME:
			p.advance()
			if p.cur.Type == kwUNBOUNDED {
				*passwordLockTime = "UNBOUNDED"
				p.advance()
			} else {
				n, err := p.parseUserLimit(0, 32767, "PASSWORD_LOCK_TIME")
				if err != nil {
					return err
				}
				*passwordLockTime = n
			}

		case p.cur.Type == kwACCOUNT:
			p.advance() // consume ACCOUNT
			if p.cur.Type == kwLOCK {
				*accountLock = true
				p.advance()
			} else if p.cur.Type == kwUNLOCK {
				*accountUnlock = true
				p.advance()
			} else {
				return p.syntaxErrorAtCur()
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
