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
	var privs []privOrRole
	var allPriv bool
	if isProxy && p.peekNext().Type == kwON {
		// PROXY is no PrivType: GrantProxyStmt spells it out, and only for GRANT. Only ON
		// makes it that form, because `GRANT proxy@h TO u` names a role.
		p.advance()
		privs = []privOrRole{{name: "PROXY", privOnly: true}}
	} else {
		var err error
		privs, allPriv, err = p.parsePrivilegeList()
		if err != nil {
			return nil, err
		}
	}

	// If we see ON, it's a privilege grant; if TO, it's a role grant.
	if p.cur.Type == kwTO {
		// ALL [PRIVILEGES] is no role, and it takes an ON: `GRANT ALL TO u` is a 1064.
		if allPriv {
			return nil, p.syntaxErrorAtCur()
		}
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

	// TO user_spec [, user_spec] ...
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
		res, err := p.parseResourceOptions()
		if err != nil {
			return nil, err
		}
		if res != nil {
			stmt.Resource = res
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
		// ALL [PRIVILEGES] is no role, and it takes an ON: `REVOKE ALL FROM u` is a 1064.
		if allPriv {
			return nil, p.syntaxErrorAtCur()
		}
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

	// FROM user_spec [, user_spec] ...
	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}

	users, err := p.parseGranteeList()
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
// TO or FROM that settles what it was: RoleOrPrivElem (TiDB v8.5.5
// parser.y:14759). A bare word can be either a privilege or a role, so the
// element carries which readings are open to it.
type privOrRole struct {
	name string
	// roleOnly marks an element the privilege reading is closed to: a
	// string literal and a name@host are RolenameWithoutIdent
	// (parser.y:11385) and nothing else, which is why `GRANT 'role1' ON *.*
	// TO u1` and `GRANT role1@localhost ON *.* TO u1` are both 1064 on
	// v8.5.5 while the TO forms of each run.
	roleOnly bool
	// privOnly marks an element the role reading is closed to: a keyword
	// reduces to PrivElem, and Rolename takes a string literal or an
	// identifier (parser.y:11395), never a keyword. So `GRANT SELECT TO u1`
	// and `GRANT GRANT OPTION TO u1` are 1064 where `GRANT role1 TO u1` and
	// `GRANT INVOKE LAMBDA TO u1` run. A keyword with an @host after it is
	// a role after all, and comes through as roleOnly instead.
	privOnly bool
}

// parsePrivilegeList parses a comma-separated RoleOrPrivElemList
// (parser.y:14793). It returns the elements and whether ALL [PRIVILEGES] was
// written.
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
// is open to, which is where v8.5.5's convertToPriv answers 1064.
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
// is open to, which is where v8.5.5's convertToRole answers 1064.
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

// parsePrivOrRoleElem parses one RoleOrPrivElem (parser.y:14759).
func (p *Parser) parsePrivOrRoleElem() (privOrRole, error) {
	// A leading '@variable' is a user/system variable token, not a
	// RoleOrPrivElem: a role is a stringLit, an identifier, or name@host
	// (Rolename, parser.y:11159) and a privilege is a keyword or a run of
	// identifiers (PrivElem/ExtendedPriv, parser.y:14501), so v8.5.5 answers
	// 1064 for `GRANT @v TO u` and `REVOKE @v FROM u`. The lexer folds `@v`
	// into one tokIDENT whose Str starts with '@'; a backquoted `@v` is an
	// ordinary name and atAtIdent lets it through.
	if p.atAtIdent() {
		return privOrRole{}, &ParseError{
			Message:  "expected privilege or role name",
			Position: p.cur.Loc,
		}
	}

	// A string literal is a role name, and may carry its own @host.
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
		return privOrRole{name: name, roleOnly: true}, nil
	}

	// LOAD FROM S3 and SELECT INTO S3 are spelled out in the grammar
	// (parser.y:14780) because LOAD, FROM, SELECT and INTO are reserved and
	// so cannot reach ExtendedPriv's run of identifiers. S3 is the word the
	// rule names and not a name, there or here: v8.5.5 answers 1064 for a
	// backquoted S3, so isBareWord is what matches it.
	for _, form := range [...]struct {
		first, second int
		third, name   string
	}{
		{kwLOAD, kwFROM, "S3", "LOAD FROM S3"},
		{kwSELECT, kwINTO, "S3", "SELECT INTO S3"},
	} {
		if p.cur.Type != form.first {
			continue
		}
		snap := p.snapshotTokenStream()
		p.advance()
		if p.cur.Type == form.second {
			p.advance()
			if p.isBareWord(form.third) {
				p.advance()
				return privOrRole{name: form.name}, nil
			}
		}
		p.restoreTokenStream(snap)
		break
	}

	// A single word with an @host after it is a role and nothing else
	// (RolenameComposed, parser.y:11377), whether the word is a keyword or
	// an identifier.
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
			return privOrRole{name: name + "@" + host, roleOnly: true}, nil
		}
	}

	name, word, err := p.parsePrivilegeName()
	if err != nil {
		return privOrRole{}, err
	}
	return privOrRole{name: name, privOnly: word.keyword}, nil
}

// privWord records what parsePrivilegeName read. The caller needs it because
// a privilege word and a name that happens to spell one are not the same
// thing: v8.5.5, mysql:8.0 and mariadb:11.8 all run `GRANT `grant` TO u` and `GRANT `create view` TO u`,
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

// tidbPrivilegeNames is PrivType, parser.y:14568: the privileges TiDB v8.5.5 spells with
// keywords. It is a closed list, so `GRANT ADMIN ON *.* TO u` is a 1064. A longer privilege
// name is a run of bare identifiers instead (ExtendedPriv), which does not come through here.
var tidbPrivilegeNames = map[string]bool{
	"ALL": true, "ALL PRIVILEGES": true, "ALTER": true, "CREATE": true, "CREATE USER": true,
	"CREATE TABLESPACE": true, "TRIGGER": true, "DELETE": true, "DROP": true, "PROCESS": true,
	"EXECUTE": true, "INDEX": true, "INSERT": true, "SELECT": true, "SUPER": true,
	"SHOW DATABASES": true, "UPDATE": true, "GRANT OPTION": true, "REFERENCES": true,
	"REPLICATION SLAVE": true, "REPLICATION CLIENT": true, "USAGE": true, "RELOAD": true,
	"FILE": true, "CONFIG": true, "CREATE TEMPORARY TABLES": true, "LOCK TABLES": true,
	"CREATE VIEW": true, "SHOW VIEW": true, "CREATE ROLE": true, "DROP ROLE": true,
	"CREATE ROUTINE": true, "ALTER ROUTINE": true, "EVENT": true, "SHUTDOWN": true,
}

// parsePrivilegeName parses a single privilege name.
// Privilege names can be keywords like SELECT, INSERT, UPDATE, DELETE, etc.
// Multi-word: CREATE VIEW, ALTER ROUTINE, SHOW DATABASES, CREATE TEMPORARY TABLES,
// LOCK TABLES, REPLICATION CLIENT, REPLICATION SLAVE, CREATE TABLESPACE,
// CREATE USER, CREATE ROLE, DROP ROLE.
// Column-level: SELECT (col1, col2), INSERT (col1).
func (p *Parser) parsePrivilegeName() (string, privWord, error) {
	start := p.pos()
	var name string
	var word privWord
	first := p.cur.Type
	if p.isPlainIdentifier() {
		first = tokIDENT
	}
	switch first {
	case kwSELECT, kwINSERT, kwUPDATE, kwDELETE, kwCREATE, kwDROP, kwALTER,
		kwINDEX, kwEXECUTE, kwGRANT, kwREFERENCES, kwTRIGGER, kwEVENT,
		kwSHOW, kwLOCK:
		tok := p.advance()
		name = strings.ToUpper(tok.Str)
		word.keyword = true
	case tokIDENT:
		tok := p.advance()
		name = strings.ToUpper(tok.Str)
		// ExtendedPriv (parser.y:14749) is a run of identifiers, which is
		// how the long privilege names are spelled: INVOKE LAMBDA, INVOKE
		// SAGEMAKER, INVOKE COMPREHEND. The run takes identifier tokens,
		// backquoted or not, so it stops at ON, at a comma and at any
		// keyword, while `GRANT INVOKE `LAMBDA` ON *.* TO u` runs on
		// v8.5.5 as the bare spelling does.
		//
		// The run is unbounded, so the words are joined once rather than
		// concatenated per word: a 20,000-word run is a 40 KB statement,
		// and growing the name in place would copy it 20,000 times.
		if p.isPlainIdentifier() && !p.atAtIdent() {
			var b strings.Builder
			b.WriteString(name)
			for p.isPlainIdentifier() && !p.atAtIdent() {
				b.WriteByte(' ')
				b.WriteString(strings.ToUpper(p.advance().Str))
				word.multiToken = true
			}
			name = b.String()
		}
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
		if !tidbPrivilegeNames[name] {
			return "", word, &ParseError{
				Message:  "expected privilege name",
				Position: start,
			}
		}
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
			// TABLES is part of the name. Without it the name is no privilege, which the
			// caller refuses.
			if p.cur.Type != kwTABLES {
				return "CREATE TEMPORARY"
			}
			p.advance()
			return "CREATE TEMPORARY TABLES"
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
	// Columns follow a PrivElem. ExtendedPriv, a run of identifiers, takes none, so `GRANT
	// a(b) ON t TO u` is a 1064.
	if p.cur.Type != '(' || !word.keyword {
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

// parseRoleRefList parses a comma-separated list of role references where the
// grammar wants a RolenameList but the surrounding code reads each as a user
// name (the DEFAULT ROLE clause of CREATE/ALTER USER). It is parseUserList with
// CURRENT_USER refused: CURRENT_USER is a Username (parser.y:11346) and not a
// Rolename, so v8.5.5 answers 1064 for CREATE USER u DEFAULT ROLE current_user.
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

// parseGranteeList parses the UserSpecList (TiDB v8.5.5 parser.y:14037) the
// privilege forms of GRANT and REVOKE take, where each grantee may carry an
// AuthOption of its own: `GRANT ALL ON db.* TO 'u'@'localhost' IDENTIFIED BY
// '...'` is the spelling an old dump writes.
//
// The role forms take a UsernameList instead (parser.y:14712 and 15035), so
// `GRANT 'r' TO u IDENTIFIED BY '...'` is 1064, as is the same tail on GRANT
// PROXY; only the privilege forms come through here.
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

// parseRoleList parses a comma-separated list of role names, the grammar's
// RolenameList (TiDB v8.5.5 parser.y:11405).
func (p *Parser) parseRoleList() ([]string, error) {
	var roles []string
	for {
		role, err := p.parseRoleName()
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

// parseRoleName parses Rolename (TiDB v8.5.5 parser.y:11395), which is not
// Username (parser.y:11333). A role written on its own is RoleNameString
// (parser.y:11371) -- a string literal or a plain identifier -- where a user
// written on its own is StringName (parser.y:13843) and so may be any
// unreserved keyword. That is the whole difference, and it is why v8.5.5
// takes CREATE USER resource and answers 1064 for CREATE ROLE RESOURCE while
// running the same statement with the role name quoted. Only the role@host
// form reaches StringName on both sides, so a keyword is a role name there.
//
// GRANT and REVOKE read their role list as a privilege list first and settle
// which it was on the TO or FROM that follows, so their names do not come
// through here.
func (p *Parser) parseRoleName() (string, error) {
	if err := p.checkRoleName(); err != nil {
		return "", err
	}
	return p.parseUserName()
}

// checkRoleName refuses a keyword upstream's scanner would not hand over as a bare
// identifier, unless the role@host form lets a StringName through.
func (p *Parser) checkRoleName() error {
	// CURRENT_USER is a Username (parser.y:11346) but not a Rolename
	// (parser.y:11395 has no CURRENT_USER alternative, and the keyword is not a
	// StringName), so it is refused wherever a role is named: v8.5.5 answers
	// 1064 for CREATE/DROP/SET ROLE current_user, SET DEFAULT ROLE current_user
	// TO u, and SHOW GRANTS ... USING current_user.
	if p.cur.Type == kwCURRENT_USER {
		return &ParseError{
			Message:  "expected role name",
			Position: p.cur.Loc,
		}
	}
	if p.cur.Type >= 700 && !isReserved(p.cur.Type) && !p.isPlainIdentifier() {
		next := p.peekNext()
		if next.Type != tokIDENT || !strings.HasPrefix(next.Str, "@") {
			return &ParseError{
				Message:  "expected role name",
				Position: p.cur.Loc,
			}
		}
	}
	return nil
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

// parseUserAndHost parses Username (TiDB v8.5.5 parser.y:11333): a name, then
// optionally the @host half, reporting whether one was written.
//
// StringName (parser.y:13843) covers the name, so it is a string literal or an
// identifier, backquoted or not, and a reserved word is not one: v8.5.5,
// mysql:8.0 and mariadb:11.8 all answer 1064 for CREATE USER select@localhost.
//
// CURRENT_USER is the one reserved word a Username accepts, as its own
// alternative `"CURRENT_USER" OptionalBraces` (parser.y:11346): it names the
// connected user, takes no @host, and allows an optional empty "()" like the
// function spelling. It reaches every account-management statement through
// here -- CREATE/ALTER/DROP/RENAME USER, GRANT/REVOKE role, GRANT PROXY, SET
// DEFAULT ROLE, SET PASSWORD FOR -- so all accept it where TiDB does.
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
// `u@localhost`, and it is not a host. v8.5.5, mysql:8.0 and mariadb:11.8 all answer 1064 for
// `CREATE USER 'u' `@localhost“.
func (p *Parser) atAtIdent() bool {
	if p.cur.Type != tokIDENT || !strings.HasPrefix(p.cur.Str, "@") {
		return false
	}
	return !strings.HasPrefix(p.inputText(p.cur.Loc, p.cur.End), "`")
}

// parseUserHostOpt consumes the optional @host half of a user name.
//
// The lexer folds `@word` and "@`word`" into a single identifier token, which
// is the grammar's `StringName singleAtIdentifier` (parser.y:11342), so a host
// written without quotes arrives attached to its '@'. A quoted host arrives as
// a separate string literal, the grammar's `StringName '@' StringName`
// (parser.y:11338), and there the '@' has to touch it: v8.5.5, mysql:8.0 and
// mariadb:11.8 all answer 1064 for CREATE USER u@ 'h' and run CREATE USER
// u @'h'. A '@' with nothing attached is an empty host, which all three run.
func (p *Parser) parseUserHostOpt() (string, bool, error) {
	if !p.atAtIdent() {
		return "", false, nil
	}
	// A bare `@` spans one source byte; a folded host, `@'host'` or even the empty
	// `@''`, spans more. Str is "@" for both, so the span, not Str, tells them
	// apart -- `'u'@''` is a valid empty host, and its next token (an auth option,
	// say) must not be read as the host.
	if p.cur.End-p.cur.Loc > 1 {
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
	//
	// The IF/NOT/EXISTS run must be complete: TiDB v8.5.5 parser.y IfNotExists
	// reads `/* empty */ | "IF" NotSym "EXISTS"`, so `CREATE USER IF EXISTS u`
	// and `CREATE USER IF NOT u` are 1064. parseIfNotExists no-ops when p.cur
	// is not IF, so `CREATE USER u` is unchanged.
	ifNotExists, err := p.parseIfNotExists()
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

	// [RESOURCE GROUP group_name]
	if err := p.parseUserResourceGroup(&stmt.ResourceGroup); err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseUserResourceGroup parses the RESOURCE GROUP clause CREATE USER and
// ALTER USER take last, after the comment and attribute options
// (ResourceGroupNameOption, TiDB v8.5.5 parser.y:14213). It is TiDB's and not
// MySQL's: v8.5.5 runs `CREATE USER u RESOURCE GROUP rg` and MySQL 8.0
// answers 1064 for it.
//
// The clause ends the statement, so `RESOURCE GROUP rg COMMENT 'c'` is 1064
// where `COMMENT 'c' RESOURCE GROUP rg` is not, and a second RESOURCE GROUP
// is 1064 too.
func (p *Parser) parseUserResourceGroup(out *string) error {
	if p.cur.Type != kwRESOURCE {
		return nil
	}
	p.advance() // consume RESOURCE
	if _, ok := p.match(kwGROUP); !ok {
		return &ParseError{
			Message:  "expected GROUP after RESOURCE",
			Position: p.cur.Loc,
		}
	}
	name, err := p.parseResourceGroupName()
	if err != nil {
		return err
	}
	*out = name
	return nil
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
	//
	// EXISTS is required after IF: TiDB v8.5.5 parser.y IfExists reads
	// `/* empty */ | "IF" "EXISTS"`, so `DROP USER IF u` is 1064.
	ifExists, err := p.parseIfExists()
	if err != nil {
		return nil, err
	}
	stmt.IfExists = ifExists

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
	//
	// EXISTS is required after IF: TiDB v8.5.5 parser.y IfExists reads
	// `/* empty */ | "IF" "EXISTS"`, so `ALTER USER IF u` is 1064.
	ifExists, err := p.parseIfExists()
	if err != nil {
		return nil, err
	}
	stmt.IfExists = ifExists

	// Check for ALTER USER USER() ... (user_func_auth_option or registration_option)
	if p.cur.Type == kwUSER && p.peekNext().Type == '(' {
		stmt.IsUserFunc = true
		userStart := p.pos()
		p.advance() // consume USER
		p.advance() // consume (
		// The current-user function spec is `"USER" '(' ')'` and the closing
		// paren is required (parser.y:13745), so TiDB v8.5.5 answers 1064 for
		// `ALTER USER USER( IDENTIFIED BY 's'`.
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}

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
			// ALTER USER USER() takes only `IDENTIFIED BY AuthString`, a plain
			// string (parser.y:13745): the general AuthOption's `BY PASSWORD`,
			// `WITH plugin` and `BY RANDOM PASSWORD` do not reach the USER() form,
			// so TiDB v8.5.5 answers 1064 to
			// `ALTER USER USER() IDENTIFIED BY PASSWORD '...'`.
			p.advance() // consume IDENTIFIED
			if _, err := p.expect(kwBY); err != nil {
				return nil, err
			}
			if p.cur.Type != tokSCONST {
				return nil, p.syntaxErrorAtCur()
			}
			userSpec.Password = p.cur.Str
			p.advance()
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

	// [RESOURCE GROUP group_name]
	if err := p.parseUserResourceGroup(&stmt.ResourceGroup); err != nil {
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
//	CREATE ROLE [IF NOT EXISTS] role [, role] ...
func (p *Parser) parseCreateRoleStmt(start int) (*nodes.CreateRoleStmt, error) {
	p.advance() // consume ROLE

	stmt := &nodes.CreateRoleStmt{Loc: nodes.Loc{Start: start}}

	// [IF NOT EXISTS]
	if p.cur.Type == kwIF {
		p.advance()
		if _, err := p.expect(kwNOT); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwEXISTS_KW); err != nil {
			return nil, err
		}
		stmt.IfNotExists = true
	}

	// role [, role] ...
	roles, err := p.parseRoleList()
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
		if _, err := p.expect(kwEXISTS_KW); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}

	// role [, role] ...
	roles, err := p.parseRoleList()
	if err != nil {
		return nil, err
	}
	stmt.Roles = roles

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseUserAuthOption parses AuthOption (TiDB v8.5.5 parser.y:14334), the
// IDENTIFIED portion of a user spec. Caller has verified
// p.cur.Type == kwIDENTIFIED.
//
//	IDENTIFIED BY auth_string
//	IDENTIFIED BY PASSWORD hash_string
//	IDENTIFIED WITH auth_plugin [BY auth_string | AS hash_string]
//
// RANDOM PASSWORD is MySQL 8.0's and not in this grammar, but it is left
// where it was: v8.5.5 answers 1064 for `IDENTIFIED BY RANDOM PASSWORD`, so
// taking it is a false accept, and one to retire on its own rather than
// alongside the clauses this change adds.
// A bare IDENTIFIED is not a clause: v8.5.5, mysql:8.0, mariadb:11.8 and oceanbase-ce:4.3.5 all answer
// 1064 for `CREATE USER u@'h' IDENTIFIED` and for the same tail on a GRANT
// grantee, so an arm has to follow.
func (p *Parser) parseUserAuthOption(spec *nodes.UserSpec) error {
	p.advance() // consume IDENTIFIED
	switch p.cur.Type {
	case kwBY:
		p.advance()
		// BY PASSWORD names the stored hash rather than the password
		// (parser.y:14367). MySQL 8.0 dropped the form; v8.5.5,
		// mariadb:11.8 and oceanbase-ce:4.3.5 all run it.
		if p.cur.Type == kwPASSWORD {
			p.advance()
			hash, ok := p.parseAuthHashString()
			if !ok {
				return p.syntaxErrorAtCur()
			}
			spec.AuthHash = hash
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
		switch p.cur.Type {
		case kwBY:
			p.advance()
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
		case kwAS:
			p.advance()
			hash, ok := p.parseAuthHashString()
			if !ok {
				return p.syntaxErrorAtCur()
			}
			spec.AuthHash = hash
		}
	default:
		return p.syntaxErrorAtCur()
	}
	return nil
}

// parseAuthPlugin parses AuthPlugin (parser.y:14377), which is StringName and
// so a string literal as well as an identifier: `IDENTIFIED WITH
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

// parseAuthHashString parses HashString (parser.y:14379): a string literal or
// a hex literal. A dumped caching_sha2_password account carries its hash as
// 0x..., which holds bytes no string literal spells.
func (p *Parser) parseAuthHashString() (string, bool) {
	if p.cur.Type == tokSCONST || p.cur.Type == tokXCONST {
		s := p.cur.Str
		p.advance()
		return s, true
	}
	return "", false
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
		roles, err := p.parseRoleList()
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
			roles, err := p.parseRoleList()
			if err != nil {
				return nil, err
			}
			stmt.AllExcept = roles
		}
	default:
		roles, err := p.parseRoleList()
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
		//
		// TO is required in every pair: TiDB v8.5.5 parser.y UserToUserList reads
		// `UserToUser | UserToUserList ',' UserToUser` and UserToUser is
		// `Username "TO" Username`, so `RENAME USER a b` and a list with any
		// pair missing TO are 1064.
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

// parseRequireClause parses RequireClause (TiDB v8.5.5 parser.y:14118).
//
//	REQUIRE NONE
//	REQUIRE SSL
//	REQUIRE X509
//	REQUIRE require_list_element [[AND] require_list_element] ...
//
//	require_list_element: {
//	    CIPHER 'cipher'
//	  | ISSUER 'issuer'
//	  | SUBJECT 'subject'
//	  | SAN 'san'
//	  | TOKEN_ISSUER 'issuer'
//	}
//
// NONE, SSL and X509 each stand alone: the grammar gives them their own
// productions and RequireList (parser.y:14145) holds only the value-carrying
// options, so `REQUIRE SSL AND X509`, `REQUIRE X509 SAN 'DNS:x'` and
// `REQUIRE NONE SSL` are 1064 on v8.5.5, mysql:8.0 and mariadb:11.8 alike.
//
// SAN and TOKEN_ISSUER are TiDB's own; mysql:8.0 and mariadb:11.8 answer 1064
// for both. They are read as bare words because the keyword table does not
// hold them, which is also what the server does: `REQUIRE `SAN` 'DNS:x'` is
// 1064 on v8.5.5, since a backquoted word is an identifier and not the
// keyword the rule names.
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
	// two elements the AND is optional (parser.y:14145).
	for required := true; ; {
		var out *string
		switch {
		case p.cur.Type == kwCIPHER:
			out = &req.Cipher
		case p.cur.Type == kwISSUER:
			out = &req.Issuer
		case p.cur.Type == kwSUBJECT:
			out = &req.Subject
		case p.isBareWord("SAN"):
			out = &req.SAN
		case p.isBareWord("TOKEN_ISSUER"):
			out = &req.TokenIssuer
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
//
// Each count is Int64Num (ConnectionOption in TiDB v8.5.5 pkg/parser/parser.y), so it is
// required and 2^63 is a syntax error.
func (p *Parser) parseResourceOptions() (*nodes.ResourceOption, error) {
	var res *nodes.ResourceOption
	for p.cur.Type == kwMAX_QUERIES_PER_HOUR || p.cur.Type == kwMAX_UPDATES_PER_HOUR ||
		p.cur.Type == kwMAX_CONNECTIONS_PER_HOUR || p.cur.Type == kwMAX_USER_CONNECTIONS {
		if res == nil {
			res = &nodes.ResourceOption{Loc: nodes.Loc{Start: p.pos()}}
		}
		switch p.cur.Type {
		case kwMAX_QUERIES_PER_HOUR:
			p.advance()
			n, err := p.parseAdminInt64()
			if err != nil {
				return nil, err
			}
			res.MaxQueriesPerHour = int(n)
			res.HasMaxQueries = true
		case kwMAX_UPDATES_PER_HOUR:
			p.advance()
			n, err := p.parseAdminInt64()
			if err != nil {
				return nil, err
			}
			res.MaxUpdatesPerHour = int(n)
			res.HasMaxUpdates = true
		case kwMAX_CONNECTIONS_PER_HOUR:
			p.advance()
			n, err := p.parseAdminInt64()
			if err != nil {
				return nil, err
			}
			res.MaxConnectionsPerHour = int(n)
			res.HasMaxConnections = true
		case kwMAX_USER_CONNECTIONS:
			p.advance()
			n, err := p.parseAdminInt64()
			if err != nil {
				return nil, err
			}
			res.MaxUserConnections = int(n)
			res.HasMaxUserConn = true
		}
	}
	if res != nil {
		res.Loc.End = p.prev.End
	}
	return res, nil
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
//
// Ref: PasswordOrLockOption in TiDB v8.5.5 parser.y:13993-14082. Each option
// requires its full operand: HISTORY {n|DEFAULT}, REUSE INTERVAL {n DAY|DEFAULT},
// REQUIRE CURRENT DEFAULT (the only REQUIRE form), ACCOUNT {LOCK|UNLOCK}.
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
					// Int64Num, so it is required and 2^63 is a syntax error.
					days, err := p.parseAdminInt64()
					if err != nil {
						return err
					}
					n := strconv.FormatInt(days, 10)
					if _, err := p.expect(kwDAY); err != nil {
						return err
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
				} else {
					return p.syntaxErrorAtCur()
				}
			case next.Type == kwREUSE:
				p.advance() // consume PASSWORD
				p.advance() // consume REUSE
				// INTERVAL is required, then DEFAULT or a count with DAY.
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
				// The only form is REQUIRE CURRENT DEFAULT; TiDB answers 1064 for
				// REQUIRE DEFAULT, REQUIRE CURRENT, and REQUIRE CURRENT OPTIONAL.
				if _, err := p.expect(kwCURRENT); err != nil {
					return err
				}
				if _, err := p.expect(kwDEFAULT); err != nil {
					return err
				}
				*passwordRequireCurrent = "DEFAULT"
			default:
				return nil
			}

		case p.cur.Type == kwFAILED_LOGIN_ATTEMPTS:
			p.advance()
			// Int64Num, so it is required and 2^63 is a syntax error.
			n, err := p.parseAdminInt64()
			if err != nil {
				return err
			}
			*failedLoginAttempts = int(n)
			*hasFailedLogin = true

		case p.cur.Type == kwPASSWORD_LOCK_TIME:
			p.advance()
			if p.cur.Type == tokICONST {
				// Int64Num, so 2^63 is a syntax error.
				n, err := p.parseAdminInt64()
				if err != nil {
					return err
				}
				*passwordLockTime = strconv.FormatInt(n, 10)
			} else if p.cur.Type == kwUNBOUNDED {
				*passwordLockTime = "UNBOUNDED"
				p.advance()
			} else {
				return p.syntaxErrorAtCur()
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

// parseSetResourceGroupStmt parses SET RESOURCE GROUP group_name.
// SET already consumed. p.cur is RESOURCE.
//
// Ref: TiDB v8.5.5 parser.y:11017
//
//	SET RESOURCE GROUP group_name
//
// MySQL's `FOR thread_id [, thread_id] ...` is not TiDB's: v8.5.5 answers
// 1064 for `SET RESOURCE GROUP x FOR 1`.
func (p *Parser) parseSetResourceGroupStmt(start int) (*nodes.SetResourceGroupStmt, error) {
	p.advance() // consume RESOURCE
	if _, ok := p.match(kwGROUP); !ok {
		return nil, &ParseError{
			Message:  "expected GROUP after RESOURCE",
			Position: p.cur.Loc,
		}
	}

	name, err := p.parseResourceGroupName()
	if err != nil {
		return nil, err
	}

	stmt := &nodes.SetResourceGroupStmt{
		Loc:  nodes.Loc{Start: start},
		Name: name,
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}
