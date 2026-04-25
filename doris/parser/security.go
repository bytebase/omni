package parser

import (
	"strings"

	"github.com/bytebase/omni/doris/ast"
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

// parseCreateDictionary parses:
//
//	CREATE DICTIONARY [IF NOT EXISTS] name
//	    USING table_name
//	    (col1 KEY, col2 VALUE, ...)
//	    LAYOUT(layout_type)
//	    PROPERTIES(...)
//
// On entry, CREATE has been consumed; cur is DICTIONARY.
func (p *Parser) parseCreateDictionary(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume DICTIONARY

	stmt := &ast.CreateDictionaryStmt{}

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

	// LAYOUT(layout_type)
	if p.cur.Kind == kwLAYOUT {
		p.advance()
		if _, err := p.expect(int('(')); err != nil {
			return nil, err
		}
		// Layout type is an identifier/keyword; normalize to lowercase.
		layoutName, layoutLoc, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Layout = strings.ToLower(layoutName)
		endLoc = layoutLoc
		rparen, err := p.expect(int(')'))
		if err != nil {
			return nil, err
		}
		endLoc = rparen.Loc
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

		if p.cur.Kind == int(',') {
			p.advance()
		}
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

// parseDropDictionary parses:
//
//	DROP DICTIONARY [IF EXISTS] name
//
// On entry, DROP has been consumed; cur is DICTIONARY.
func (p *Parser) parseDropDictionary(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume DICTIONARY

	stmt := &ast.DropDictionaryStmt{}

	if p.cur.Kind == kwIF {
		p.advance()
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}

	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	stmt.Loc = startLoc.Merge(ast.NodeLoc(name))
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

// parseCreateRole parses:
//
//	CREATE ROLE [IF NOT EXISTS] name [COMMENT 'text']
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

	name, nameLoc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
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

// parseAlterRole parses:
//
//	ALTER ROLE name COMMENT 'text'
//
// On entry, ALTER has been consumed; cur is ROLE.
func (p *Parser) parseAlterRole(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume ROLE

	stmt := &ast.AlterRoleStmt{}

	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	if _, err := p.expect(kwCOMMENT); err != nil {
		return nil, err
	}
	if p.cur.Kind != tokString {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Comment = p.cur.Str
	endLoc := p.cur.Loc
	p.advance()

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseDropRole parses:
//
//	DROP ROLE [IF EXISTS] name
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

	name, nameLoc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	stmt.Loc = startLoc.Merge(nameLoc)
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

	// Optional @'host'
	if p.cur.Kind == int('@') {
		p.advance() // consume @
		if p.cur.Kind != tokString && !isIdentifierToken(p.cur.Kind) {
			return nil, p.syntaxErrorAtCur()
		}
		var hostLoc ast.Loc
		host, hostLoc, err = p.parseIdentifierOrString()
		if err != nil {
			return nil, err
		}
		endLoc = hostLoc
	}

	return &ast.UserIdentity{
		Username: username,
		Host:     host,
		Loc:      ast.Loc{Start: startLoc.Start, End: endLoc.End},
	}, nil
}

// parseCreateUser parses:
//
//	CREATE USER [IF NOT EXISTS] 'user'@'host'
//	    [IDENTIFIED BY 'password' | IDENTIFIED BY PASSWORD 'hash']
//	    [DEFAULT ROLE 'role']
//	    [PASSWORD_EXPIRE [INTERVAL n DAY]]
//	    [FAILED_LOGIN_ATTEMPTS n]
//	    [PASSWORD_LOCK_TIME n DAY]
//	    [PASSWORD_HISTORY n]
//	    [COMMENT 'text']
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

	// Optional IDENTIFIED BY ['PASSWORD'] 'value'
	if p.cur.Kind == kwIDENTIFIED {
		p.advance() // consume IDENTIFIED
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		if p.cur.Kind == kwPASSWORD {
			// IDENTIFIED BY PASSWORD 'hash'
			p.advance()
			if p.cur.Kind != tokString {
				return nil, p.syntaxErrorAtCur()
			}
			stmt.PasswordHash = p.cur.Str
			endLoc = p.cur.Loc
			p.advance()
		} else {
			// IDENTIFIED BY 'password'
			if p.cur.Kind != tokString {
				return nil, p.syntaxErrorAtCur()
			}
			stmt.Password = p.cur.Str
			endLoc = p.cur.Loc
			p.advance()
		}
	}

	// Optional DEFAULT ROLE 'role'
	if p.cur.Kind == kwDEFAULT {
		p.advance()
		if _, err := p.expect(kwROLE); err != nil {
			return nil, err
		}
		roleName, roleLoc, err := p.parseIdentifierOrString()
		if err != nil {
			return nil, err
		}
		stmt.DefaultRole = roleName
		endLoc = roleLoc
	}

	// Optional password policy clauses
	endLoc = p.parsePasswordPolicyCreate(stmt, endLoc)

	// Optional COMMENT
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

// parsePasswordPolicyCreate parses optional password policy clauses for CREATE USER.
// Modifies stmt in place and returns the updated endLoc.
func (p *Parser) parsePasswordPolicyCreate(stmt *ast.CreateUserStmt, endLoc ast.Loc) ast.Loc {
	for {
		switch p.cur.Kind {
		case kwPASSWORD_EXPIRE:
			stmt.PasswordExpire = true
			endLoc = p.cur.Loc
			p.advance()
			// Optional INTERVAL n DAY
			if p.cur.Kind == kwINTERVAL {
				p.advance()
				if p.cur.Kind == tokInt {
					n := int(p.cur.Ival)
					stmt.PasswordExpireInterval = n
					endLoc = p.cur.Loc
					p.advance()
				}
				// Consume DAY/DAYS
				if p.cur.Kind == kwDAY || p.cur.Kind == kwDAYS {
					endLoc = p.cur.Loc
					p.advance()
				}
			}
		case kwFAILED_LOGIN_ATTEMPTS:
			endLoc = p.cur.Loc
			p.advance()
			if p.cur.Kind == tokInt {
				n := int(p.cur.Ival)
				stmt.FailedLoginAttempts = n
				endLoc = p.cur.Loc
				p.advance()
			}
		case kwPASSWORD_LOCK_TIME:
			endLoc = p.cur.Loc
			p.advance()
			if p.cur.Kind == tokInt {
				n := int(p.cur.Ival)
				stmt.PasswordLockTime = n
				endLoc = p.cur.Loc
				p.advance()
			}
			if p.cur.Kind == kwDAY || p.cur.Kind == kwDAYS {
				endLoc = p.cur.Loc
				p.advance()
			}
		case kwPASSWORD_HISTORY:
			endLoc = p.cur.Loc
			p.advance()
			if p.cur.Kind == tokInt {
				n := int(p.cur.Ival)
				stmt.PasswordHistory = n
				endLoc = p.cur.Loc
				p.advance()
			}
		default:
			return endLoc
		}
	}
}

// parseAlterUser parses:
//
//	ALTER USER [IF EXISTS] 'user'@'host'
//	    [IDENTIFIED BY 'password']
//	    [FAILED_LOGIN_ATTEMPTS n]
//	    [PASSWORD_LOCK_TIME n DAY]
//	    [ACCOUNT_LOCK | ACCOUNT_UNLOCK]
//	    [COMMENT 'text']
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

	// Optional IDENTIFIED BY ['PASSWORD'] 'value'
	if p.cur.Kind == kwIDENTIFIED {
		p.advance()
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		if p.cur.Kind == kwPASSWORD {
			p.advance()
			if p.cur.Kind != tokString {
				return nil, p.syntaxErrorAtCur()
			}
			stmt.PasswordHash = p.cur.Str
			endLoc = p.cur.Loc
			p.advance()
		} else {
			if p.cur.Kind != tokString {
				return nil, p.syntaxErrorAtCur()
			}
			stmt.Password = p.cur.Str
			endLoc = p.cur.Loc
			p.advance()
		}
	}

	// Optional password policy / account clauses
	for {
		switch p.cur.Kind {
		case kwFAILED_LOGIN_ATTEMPTS:
			endLoc = p.cur.Loc
			p.advance()
			if p.cur.Kind == tokInt {
				n := int(p.cur.Ival)
				stmt.FailedLoginAttempts = n
				endLoc = p.cur.Loc
				p.advance()
			}
		case kwPASSWORD_LOCK_TIME:
			endLoc = p.cur.Loc
			p.advance()
			if p.cur.Kind == tokInt {
				n := int(p.cur.Ival)
				stmt.PasswordLockTime = n
				endLoc = p.cur.Loc
				p.advance()
			}
			if p.cur.Kind == kwDAY || p.cur.Kind == kwDAYS {
				endLoc = p.cur.Loc
				p.advance()
			}
		case kwACCOUNT_LOCK:
			stmt.AccountLock = true
			endLoc = p.cur.Loc
			p.advance()
		case kwACCOUNT_UNLOCK:
			stmt.AccountUnlock = true
			endLoc = p.cur.Loc
			p.advance()
		case kwCOMMENT:
			p.advance()
			if p.cur.Kind != tokString {
				return nil, p.syntaxErrorAtCur()
			}
			stmt.Comment = p.cur.Str
			endLoc = p.cur.Loc
			p.advance()
		default:
			goto doneAlterUser
		}
	}
doneAlterUser:

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
				return p.input[startOff:endOff], ast.Loc{Start: startOff, End: endOff}, nil
			}
		}
		p.advance()
	}

	return "", ast.NoLoc(), p.syntaxErrorAtCur()
}
