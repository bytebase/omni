package parser

import (
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
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
			if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "admin") {
				p.advance()
				if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "option") {
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

	// [WITH GRANT OPTION]
	if p.cur.Type == kwWITH {
		p.advance()
		if p.cur.Type == kwGRANT {
			p.advance()
			if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "option") {
				p.advance()
			}
			stmt.WithGrant = true
		}
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseRevokeStmt parses a REVOKE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/revoke.html
//
//	REVOKE priv_type [, priv_type] ... ON [object_type] priv_level FROM user [, user] ...
//	REVOKE role [, role] ... FROM user [, user] ...
func (p *Parser) parseRevokeStmt() (nodes.Node, error) {
	start := p.pos()
	p.advance() // consume REVOKE

	// Parse the names (could be privileges or roles).
	privs, allPriv, err := p.parsePrivilegeList()
	if err != nil {
		return nil, err
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
			Loc:   nodes.Loc{Start: start},
			Roles: privs,
			From:  users,
		}
		stmt.Loc.End = p.pos()
		return stmt, nil
	}

	// Privilege revoke
	stmt := &nodes.RevokeStmt{Loc: nodes.Loc{Start: start}}
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

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parsePrivilegeList parses a comma-separated list of privilege names.
// Returns the list and whether ALL [PRIVILEGES] was specified.
func (p *Parser) parsePrivilegeList() ([]string, bool, error) {
	if p.cur.Type == kwALL {
		p.advance()
		// Optional PRIVILEGES keyword
		if p.cur.Type == kwPRIVILEGES || (p.cur.Type == tokIDENT && eqFold(p.cur.Str, "privileges")) {
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
func (p *Parser) parsePrivilegeName() (string, error) {
	switch p.cur.Type {
	case kwSELECT, kwINSERT, kwUPDATE, kwDELETE, kwCREATE, kwDROP, kwALTER,
		kwINDEX, kwEXECUTE, kwGRANT, kwREFERENCES, kwTRIGGER, kwEVENT:
		tok := p.advance()
		return strings.ToUpper(tok.Str), nil
	case tokIDENT:
		tok := p.advance()
		return strings.ToUpper(tok.Str), nil
	default:
		// Try accepting any keyword as a privilege name
		if p.cur.Type >= 700 {
			tok := p.advance()
			return strings.ToUpper(tok.Str), nil
		}
		return "", &ParseError{
			Message:  "expected privilege name",
			Position: p.cur.Loc,
		}
	}
}

// parseGrantTarget parses the ON clause target of a GRANT/REVOKE.
//
//	[object_type] priv_level
//	object_type: TABLE | FUNCTION | PROCEDURE
//	priv_level: *.* | db_name.* | db_name.tbl_name | tbl_name | *
func (p *Parser) parseGrantTarget() (*nodes.GrantTarget, error) {
	start := p.pos()
	target := &nodes.GrantTarget{Loc: nodes.Loc{Start: start}}

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
				Loc:  nodes.Loc{Start: start, End: p.pos()},
				Schema: "*",
				Name:   "*",
			}
			if p.cur.Type == '*' {
				p.advance()
				ref.Loc.End = p.pos()
			} else {
				// *.tbl_name (unusual but handle)
				name, _, err := p.parseIdentifier()
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
	name, _, err := p.parseIdentifier()
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
			name2, _, err := p.parseIdentifier()
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
		n, _, err := p.parseIdentifier()
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
			h, _, err := p.parseIdentifier()
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
//	CREATE USER [IF NOT EXISTS] user_spec [, user_spec] ...
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
func (p *Parser) parseAlterUserStmt() (*nodes.AlterUserStmt, error) {
	start := p.pos()
	p.advance() // consume USER

	stmt := &nodes.AlterUserStmt{Loc: nodes.Loc{Start: start}}

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

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseUserSpec parses a user specification.
//
//	user_spec: user [IDENTIFIED BY 'password' | IDENTIFIED WITH auth_plugin [BY 'password']]
func (p *Parser) parseUserSpec() (*nodes.UserSpec, error) {
	start := p.pos()
	spec := &nodes.UserSpec{Loc: nodes.Loc{Start: start}}

	// Parse user name
	if p.cur.Type == tokSCONST {
		spec.Name = p.cur.Str
		p.advance()
	} else if p.isIdentToken() {
		name, _, err := p.parseIdentifier()
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
			host, _, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			spec.Host = host
		}
	}

	// Optional IDENTIFIED BY/WITH
	if p.cur.Type == kwIDENTIFIED {
		p.advance()
		if p.cur.Type == kwBY {
			// IDENTIFIED BY 'password'
			p.advance()
			if p.cur.Type == tokSCONST {
				spec.Password = p.cur.Str
				p.advance()
			}
		} else if p.cur.Type == kwWITH {
			// IDENTIFIED WITH auth_plugin [BY 'password']
			p.advance()
			if p.isIdentToken() {
				plugin, _, err := p.parseIdentifier()
				if err != nil {
					return nil, err
				}
				spec.AuthPlugin = plugin
			}
			// Optional BY 'password'
			if p.cur.Type == kwBY {
				p.advance()
				if p.cur.Type == tokSCONST {
					spec.Password = p.cur.Str
					p.advance()
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
func (p *Parser) parseRenameUserStmt(start int) (*nodes.RenameUserStmt, error) {
	p.advance() // consume USER

	stmt := &nodes.RenameUserStmt{
		Loc: nodes.Loc{Start: start},
	}

	for {
		pairStart := p.pos()
		pair := &nodes.RenameUserPair{Loc: nodes.Loc{Start: pairStart}}

		// old_user
		pair.OldUser, _, _ = p.parseIdentifier()
		if p.cur.Type == tokIDENT && p.cur.Str == "@" {
			p.advance()
			pair.OldHost, _, _ = p.parseIdentifier()
		}

		// TO
		p.match(kwTO)

		// new_user
		pair.NewUser, _, _ = p.parseIdentifier()
		if p.cur.Type == tokIDENT && p.cur.Str == "@" {
			p.advance()
			pair.NewHost, _, _ = p.parseIdentifier()
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

// parseSetResourceGroupStmt parses SET RESOURCE GROUP group_name [FOR thread_id [, thread_id] ...].
// SET already consumed. p.cur is RESOURCE.
func (p *Parser) parseSetResourceGroupStmt(start int) (*nodes.SetResourceGroupStmt, error) {
	p.advance() // consume RESOURCE
	p.advance() // consume GROUP (identifier)

	name, _, err := p.parseIdentifier()
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
