package parser

import (
	"strings"

	"github.com/bytebase/omni/starrocks/ast"
)

// ---------------------------------------------------------------------------
// ROW POLICY
// ---------------------------------------------------------------------------

// parseCreateRowPolicy parses:
//
//	CREATE ROW POLICY [IF NOT EXISTS] name ON table_name
//	    [AS {RESTRICTIVE | PERMISSIVE}]
//	    TO user_or_role
//	    USING (expr)
//
// On entry, CREATE and ROW have been consumed; cur is POLICY.
func (p *Parser) parseCreateRowPolicy(startLoc ast.Loc) (ast.Node, error) {
	if _, err := p.expect(kwPOLICY); err != nil {
		return nil, err
	}

	stmt := &ast.CreateRowPolicyStmt{}

	// Optional IF NOT EXISTS
	if p.cur.Kind == kwIF {
		p.advance()
		if _, err := p.expect(kwNOT); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfNotExists = true
	}

	// Policy name
	name, nameLoc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	endLoc := nameLoc

	// ON table_name
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	onTable, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.On = onTable
	endLoc = ast.NodeLoc(onTable)

	// Optional AS {RESTRICTIVE | PERMISSIVE}
	if p.cur.Kind == kwAS {
		p.advance()
		switch p.cur.Kind {
		case kwRESTRICTIVE:
			stmt.Type = "RESTRICTIVE"
			endLoc = p.cur.Loc
			p.advance()
		case kwPERMISSIVE:
			stmt.Type = "PERMISSIVE"
			endLoc = p.cur.Loc
			p.advance()
		default:
			return nil, p.syntaxErrorAtCur()
		}
	}

	// TO user_or_role
	if _, err := p.expect(kwTO); err != nil {
		return nil, err
	}
	toName, toLoc, err := p.parseIdentifierOrString()
	if err != nil {
		return nil, err
	}
	stmt.To = toName
	endLoc = toLoc

	// USING (expr) — consume as raw text until end of paren group
	if p.cur.Kind == kwUSING {
		p.advance()
		raw, loc, err := p.consumeParenGroup()
		if err != nil {
			return nil, err
		}
		stmt.Using = raw
		endLoc = loc
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseDropRowPolicy parses:
//
//	DROP ROW POLICY name ON table_name
//
// On entry, DROP and ROW have been consumed; cur is POLICY.
func (p *Parser) parseDropRowPolicy(startLoc ast.Loc) (ast.Node, error) {
	if _, err := p.expect(kwPOLICY); err != nil {
		return nil, err
	}

	stmt := &ast.DropRowPolicyStmt{}

	// Policy name
	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// ON table_name
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	onTable, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.On = onTable

	stmt.Loc = startLoc.Merge(ast.NodeLoc(onTable))
	return stmt, nil
}

// ---------------------------------------------------------------------------
// ENCRYPTION KEY
// ---------------------------------------------------------------------------

// parseCreateEncryptKey parses:
//
//	CREATE ENCRYPTKEY [IF NOT EXISTS] name AS 'key_value'
//
// On entry, CREATE has been consumed; cur is ENCRYPTKEY.
func (p *Parser) parseCreateEncryptKey(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume ENCRYPTKEY

	stmt := &ast.CreateEncryptKeyStmt{}

	// Optional IF NOT EXISTS
	if p.cur.Kind == kwIF {
		p.advance()
		if _, err := p.expect(kwNOT); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfNotExists = true
	}

	// Key name (can be qualified: db.key_name)
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// AS 'key_value'
	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}
	if p.cur.Kind != tokString {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Key = p.cur.Str
	endLoc := p.cur.Loc
	p.advance()

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseDropEncryptKey parses:
//
//	DROP ENCRYPTKEY [IF EXISTS] name
//
// On entry, DROP has been consumed; cur is ENCRYPTKEY.
func (p *Parser) parseDropEncryptKey(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume ENCRYPTKEY

	stmt := &ast.DropEncryptKeyStmt{}

	// Optional IF EXISTS
	if p.cur.Kind == kwIF {
		p.advance()
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}

	// Key name (can be qualified: db.key_name)
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	stmt.Loc = startLoc.Merge(ast.NodeLoc(name))
	return stmt, nil
}

// ---------------------------------------------------------------------------
// DICTIONARY
// ---------------------------------------------------------------------------

// parseCreateDictionary parses createDictionaryStatement (StarRocks.g4:842):
//
//	CREATE DICTIONARY name
//	    USING table_name
//	    (col1 KEY, col2 VALUE, ...)
//	    PROPERTIES(...)
//
// Doris's IF NOT EXISTS and LAYOUT(layout_type) are not in the rule, and the
// pinned server refuses both.
//
// On entry, CREATE has been consumed; cur is DICTIONARY.
func (p *Parser) parseCreateDictionary(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume DICTIONARY

	stmt := &ast.CreateDictionaryStmt{}

	// Dictionary name
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	endLoc := ast.NodeLoc(name)

	// USING table_name
	if _, err := p.expect(kwUSING); err != nil {
		return nil, err
	}
	usingTable, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.UsingTable = usingTable
	endLoc = ast.NodeLoc(usingTable)

	// (col1 KEY, col2 VALUE, ...)
	if p.cur.Kind == int('(') {
		cols, loc, err := p.parseDictionaryColumns()
		if err != nil {
			return nil, err
		}
		stmt.Columns = cols
		endLoc = loc
	}

	// PROPERTIES(...)
	if p.cur.Kind == kwPROPERTIES {
		props, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		stmt.Properties = props
		if len(props) > 0 {
			endLoc = ast.NodeLoc(props[len(props)-1])
		}
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseDictionaryColumns parses: (col1 KEY, col2 VALUE, ...)
// Returns the column list and the location of the closing ')'.
func (p *Parser) parseDictionaryColumns() ([]*ast.DictionaryColumn, ast.Loc, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, ast.NoLoc(), err
	}

	var cols []*ast.DictionaryColumn

	for p.cur.Kind != int(')') && p.cur.Kind != tokEOF {
		colStart := p.cur.Loc

		colName, _, err := p.parseIdentifier()
		if err != nil {
			return nil, ast.NoLoc(), err
		}

		col := &ast.DictionaryColumn{
			Name: colName,
			Loc:  colStart,
		}

		switch p.cur.Kind {
		case kwKEY:
			col.Role = "KEY"
			col.Loc.End = p.cur.Loc.End
			p.advance()
		case kwVALUE:
			col.Role = "VALUE"
			col.Loc.End = p.cur.Loc.End
			p.advance()
		default:
			return nil, ast.NoLoc(), p.syntaxErrorAtCur()
		}

		cols = append(cols, col)

		// dictionaryColumnDesc (',' dictionaryColumnDesc)* (StarRocks.g4:845):
		// the comma is required, so the pinned server refuses
		// `CREATE DICTIONARY a USING a(a KEY a KEY)`.
		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','
	}

	rparen, err := p.expect(int(')'))
	if err != nil {
		return nil, ast.NoLoc(), err
	}

	return cols, rparen.Loc, nil
}

// parseAlterDictionary parses:
//
//	ALTER DICTIONARY name PROPERTIES(...)
//
// On entry, ALTER has been consumed; cur is DICTIONARY.
func (p *Parser) parseAlterDictionary(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume DICTIONARY

	stmt := &ast.AlterDictionaryStmt{}

	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	endLoc := ast.NodeLoc(name)

	if p.cur.Kind == kwPROPERTIES {
		props, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		stmt.Properties = props
		if len(props) > 0 {
			endLoc = ast.NodeLoc(props[len(props)-1])
		}
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseDropDictionary parses dropDictionaryStatement (StarRocks.g4:848):
//
//	DROP DICTIONARY qualifiedName CACHE?
//
// The fork's IF EXISTS is not in the rule and the pinned server refuses it.
//
// On entry, DROP has been consumed; cur is DICTIONARY.
func (p *Parser) parseDropDictionary(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume DICTIONARY

	stmt := &ast.DropDictionaryStmt{}

	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	endLoc := ast.NodeLoc(name)

	if p.cur.Kind == kwCACHE {
		stmt.Cache = true
		endLoc = p.cur.Loc
		p.advance()
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseRefreshDictionary parses:
//
//	REFRESH DICTIONARY name
//
// On entry, REFRESH has been consumed; cur is DICTIONARY.
func (p *Parser) parseRefreshDictionary(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume DICTIONARY

	stmt := &ast.RefreshDictionaryStmt{}

	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	stmt.Loc = startLoc.Merge(ast.NodeLoc(name))
	return stmt, nil
}

// ---------------------------------------------------------------------------
// ROLE
// ---------------------------------------------------------------------------

// parseCreateRole parses createRoleStatement (StarRocks.g4:1651):
//
//	CREATE ROLE [IF NOT EXISTS] roleList [comment]
//
// On entry, CREATE has been consumed; cur is ROLE.
func (p *Parser) parseCreateRole(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume ROLE

	stmt := &ast.CreateRoleStmt{}

	if p.cur.Kind == kwIF {
		p.advance()
		if _, err := p.expect(kwNOT); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfNotExists = true
	}

	names, nameLoc, err := p.parseRoleList()
	if err != nil {
		return nil, err
	}
	stmt.Names = names
	endLoc := nameLoc

	if p.cur.Kind == kwCOMMENT {
		p.advance()
		if p.cur.Kind != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Comment = p.cur.Str
		endLoc = p.cur.Loc
		p.advance()
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseRoleList parses roleList (StarRocks.g4:1648):
// identifierOrString (',' identifierOrString)*.
func (p *Parser) parseRoleList() ([]string, ast.Loc, error) {
	var names []string
	var endLoc ast.Loc
	for {
		name, loc, err := p.parseIdentifierOrString()
		if err != nil {
			return nil, ast.NoLoc(), err
		}
		names = append(names, name)
		endLoc = loc
		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','
	}
	return names, endLoc, nil
}

// parseAlterRole parses alterRoleStatement (StarRocks.g4:1654):
//
//	ALTER ROLE [IF EXISTS] roleList SET COMMENT '=' string
//
// On entry, ALTER has been consumed; cur is ROLE.
func (p *Parser) parseAlterRole(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume ROLE

	stmt := &ast.AlterRoleStmt{}

	if p.cur.Kind == kwIF {
		p.advance()
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}

	names, endLoc, err := p.parseRoleList()
	if err != nil {
		return nil, err
	}
	stmt.Names = names

	if _, err := p.expect(kwSET); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwCOMMENT); err != nil {
		return nil, err
	}
	if _, err := p.expect(int('=')); err != nil {
		return nil, err
	}
	if p.cur.Kind != tokString {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Comment = p.cur.Str
	endLoc = p.cur.Loc
	p.advance()

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseDropRole parses dropRoleStatement (StarRocks.g4:1658):
//
//	DROP ROLE [IF EXISTS] roleList
//
// On entry, DROP has been consumed; cur is ROLE.
func (p *Parser) parseDropRole(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume ROLE

	stmt := &ast.DropRoleStmt{}

	if p.cur.Kind == kwIF {
		p.advance()
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}

	names, endLoc, err := p.parseRoleList()
	if err != nil {
		return nil, err
	}
	stmt.Names = names

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// ---------------------------------------------------------------------------
// USER
// ---------------------------------------------------------------------------

// parseUserIdentity parses: 'user'@'host' or user@'host' or 'user' or user
//
// Returns a *ast.UserIdentity. Host defaults to '%' when no @'host' suffix.
func (p *Parser) parseUserIdentity() (*ast.UserIdentity, error) {
	startLoc := p.cur.Loc

	// Username — string literal or bare identifier
	username, userLoc, err := p.parseIdentifierOrString()
	if err != nil {
		return nil, err
	}
	endLoc := userLoc

	host := "%"
	domain := false

	// Optional @'host', or @['host'] (StarRocks.g4:2892, user's third form,
	// which brackets the host as `userWithHostAndBlanket`).
	if p.cur.Kind == int('@') {
		p.advance() // consume @
		if p.cur.Kind == int('[') {
			p.advance() // consume '['
			bracketHost, _, err := p.parseIdentifierOrString()
			if err != nil {
				return nil, err
			}
			closeTok, err := p.expect(int(']'))
			if err != nil {
				return nil, err
			}
			host = bracketHost
			domain = true
			endLoc = closeTok.Loc
		} else {
			if p.cur.Kind != tokString && !isIdentifierToken(p.cur.Kind) {
				return nil, p.syntaxErrorAtCur()
			}
			plainHost, hostLoc, err := p.parseIdentifierOrString()
			if err != nil {
				return nil, err
			}
			host = plainHost
			endLoc = hostLoc
		}
	}

	return &ast.UserIdentity{
		Username: username,
		Host:     host,
		Domain:   domain,
		Loc:      ast.Loc{Start: startLoc.Start, End: endLoc.End},
	}, nil
}

// parseCreateUser parses createUserStatement (StarRocks.g4:1623):
//
//	CREATE USER [IF NOT EXISTS] user [authOption] [DEFAULT ROLE roleList] properties?
//
// On entry, CREATE has been consumed; cur is USER.
func (p *Parser) parseCreateUser(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume USER

	stmt := &ast.CreateUserStmt{}

	if p.cur.Kind == kwIF {
		p.advance()
		if _, err := p.expect(kwNOT); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfNotExists = true
	}

	identity, err := p.parseUserIdentity()
	if err != nil {
		return nil, err
	}
	stmt.Name = identity
	endLoc := identity.Loc

	if p.cur.Kind == kwIDENTIFIED {
		auth, end, err := p.parseAuthOption()
		if err != nil {
			return nil, err
		}
		auth.intoCreateUser(stmt)
		endLoc = end
	}

	if p.cur.Kind == kwDEFAULT {
		p.advance() // consume DEFAULT
		if _, err := p.expect(kwROLE); err != nil {
			return nil, err
		}
		names, loc, err := p.parseRoleList()
		if err != nil {
			return nil, err
		}
		stmt.DefaultRoles = names
		endLoc = loc
	}

	if p.cur.Kind == kwPROPERTIES {
		props, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		stmt.Properties = props
		endLoc = p.prev.Loc // the closing ')'
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// authOption holds the pieces of authOption (StarRocks.g4:1627):
//
//	IDENTIFIED BY PASSWORD? string
//	IDENTIFIED WITH identifierOrString ((BY | AS) string)?
type authOption struct {
	password string
	hash     string
	plugin   string
	by       string
	authStr  string
}

// intoCreateUser copies the option into a CREATE USER.
func (a authOption) intoCreateUser(stmt *ast.CreateUserStmt) {
	stmt.Password, stmt.PasswordHash = a.password, a.hash
	stmt.Plugin, stmt.AuthBy, stmt.AuthString = a.plugin, a.by, a.authStr
}

// parseAuthOption parses authOption (StarRocks.g4:1627) after IDENTIFIED.
func (p *Parser) parseAuthOption() (authOption, ast.Loc, error) {
	var a authOption
	p.advance() // consume IDENTIFIED

	switch p.cur.Kind {
	case kwBY:
		p.advance() // consume BY
		if p.cur.Kind == kwPASSWORD {
			// IDENTIFIED BY PASSWORD 'hash'
			p.advance()
			if p.cur.Kind != tokString {
				return a, ast.NoLoc(), p.syntaxErrorAtCur()
			}
			a.hash = p.cur.Str
			loc := p.cur.Loc
			p.advance()
			return a, loc, nil
		}
		if p.cur.Kind != tokString {
			return a, ast.NoLoc(), p.syntaxErrorAtCur()
		}
		a.password = p.cur.Str
		loc := p.cur.Loc
		p.advance()
		return a, loc, nil
	case kwWITH:
		p.advance() // consume WITH
		plugin, pluginLoc, err := p.parseIdentifierOrString()
		if err != nil {
			return a, ast.NoLoc(), err
		}
		a.plugin = plugin
		endLoc := pluginLoc
		if p.cur.Kind == kwBY || p.cur.Kind == kwAS {
			a.by = strings.ToUpper(p.cur.Str)
			p.advance()
			if p.cur.Kind != tokString {
				return a, ast.NoLoc(), p.syntaxErrorAtCur()
			}
			a.authStr = p.cur.Str
			endLoc = p.cur.Loc
			p.advance()
		}
		return a, endLoc, nil
	}
	return a, ast.NoLoc(), p.syntaxErrorAtCur()
}

// parseAlterUser parses alterUserStatement (StarRocks.g4:1631):
//
//	ALTER USER [IF EXISTS] user authOption
//	ALTER USER [IF EXISTS] user DEFAULT ROLE (NONE | ALL | roleList)
//	ALTER USER [IF EXISTS] user SET properties
//
// On entry, ALTER has been consumed; cur is USER.
func (p *Parser) parseAlterUser(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume USER

	stmt := &ast.AlterUserStmt{}

	if p.cur.Kind == kwIF {
		p.advance()
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}

	identity, err := p.parseUserIdentity()
	if err != nil {
		return nil, err
	}
	stmt.Name = identity
	endLoc := identity.Loc

	switch p.cur.Kind {
	case kwIDENTIFIED:
		auth, end, err := p.parseAuthOption()
		if err != nil {
			return nil, err
		}
		stmt.Password, stmt.PasswordHash = auth.password, auth.hash
		stmt.Plugin, stmt.AuthBy, stmt.AuthString = auth.plugin, auth.by, auth.authStr
		endLoc = end
	case kwDEFAULT:
		p.advance() // consume DEFAULT
		if _, err := p.expect(kwROLE); err != nil {
			return nil, err
		}
		switch p.cur.Kind {
		case kwNONE, kwALL:
			stmt.DefaultRole = strings.ToUpper(p.cur.Str)
			endLoc = p.cur.Loc
			p.advance()
		default:
			names, loc, err := p.parseRoleList()
			if err != nil {
				return nil, err
			}
			stmt.DefaultRoles = names
			endLoc = loc
		}
	case kwSET:
		p.advance() // consume SET
		props, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		stmt.Properties = props
		endLoc = p.prev.Loc // the closing ')'
	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseDropUser parses:
//
//	DROP USER [IF EXISTS] 'user'@'host'
//
// On entry, DROP has been consumed; cur is USER.
func (p *Parser) parseDropUser(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume USER

	stmt := &ast.DropUserStmt{}

	if p.cur.Kind == kwIF {
		p.advance()
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}

	identity, err := p.parseUserIdentity()
	if err != nil {
		return nil, err
	}
	stmt.Name = identity

	stmt.Loc = startLoc.Merge(identity.Loc)
	return stmt, nil
}

// parseSetPassword parses:
//
//	SET PASSWORD [FOR 'user'@'host'] = 'hash'
//	SET PASSWORD [FOR 'user'@'host'] = PASSWORD('cleartext')
//
// On entry, SET has been consumed; cur is PASSWORD.
func (p *Parser) parseSetPassword(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume PASSWORD

	stmt := &ast.SetPasswordStmt{}

	// Optional FOR 'user'@'host'
	if p.cur.Kind == kwFOR {
		p.advance()
		identity, err := p.parseUserIdentity()
		if err != nil {
			return nil, err
		}
		stmt.For = identity
	}

	// '='
	if _, err := p.expect(int('=')); err != nil {
		return nil, err
	}

	endLoc := p.cur.Loc

	if p.cur.Kind == kwPASSWORD {
		// PASSWORD('cleartext')
		p.advance()
		if _, err := p.expect(int('(')); err != nil {
			return nil, err
		}
		if p.cur.Kind != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Password = p.cur.Str
		stmt.IsHash = false
		p.advance()
		rparen, err := p.expect(int(')'))
		if err != nil {
			return nil, err
		}
		endLoc = rparen.Loc
	} else if p.cur.Kind == tokString {
		// bare hash string
		stmt.Password = p.cur.Str
		stmt.IsHash = true
		endLoc = p.cur.Loc
		p.advance()
	} else {
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// consumeParenGroup consumes a parenthesized group (including nested parens)
// and returns its raw text and location.
func (p *Parser) consumeParenGroup() (string, ast.Loc, error) {
	startTok, err := p.expect(int('('))
	if err != nil {
		return "", ast.NoLoc(), err
	}

	depth := 1
	startOff := startTok.Loc.Start
	endOff := startTok.Loc.End

	for depth > 0 && p.cur.Kind != tokEOF {
		endOff = p.cur.Loc.End
		switch p.cur.Kind {
		case int('('):
			depth++
		case int(')'):
			depth--
			if depth == 0 {
				endOff = p.cur.Loc.End
				p.advance()
				return p.input[startOff-p.baseOffset : endOff-p.baseOffset], ast.Loc{Start: startOff, End: endOff}, nil
			}
		}
		p.advance()
	}

	return "", ast.NoLoc(), p.syntaxErrorAtCur()
}
