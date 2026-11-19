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

	// TO user_or_role, where a role is written TO ROLE name
	// (DorisParser.g4:175).
	if _, err := p.expect(kwTO); err != nil {
		return nil, err
	}
	if p.cur.Kind == kwROLE {
		p.advance() // consume ROLE
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
		// booleanExpression is not optional, and the engine refuses USING().
		if strings.TrimSpace(raw[1:len(raw)-1]) == "" {
			return nil, &ParseError{Position: loc.Start, End: loc.End, Message: "expected expression in USING"}
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

	if p.cur.Kind == kwIF {
		p.advance() // consume IF
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}

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
	endLoc := ast.NodeLoc(onTable)

	// Optional FOR (userIdentify | ROLE roleName)
	if p.cur.Kind == kwFOR {
		p.advance() // consume FOR
		if p.cur.Kind == kwROLE {
			p.advance() // consume ROLE
			stmt.ForRole = true
		}
		name, nameLoc, err := p.parseIdentifierOrString()
		if err != nil {
			return nil, err
		}
		stmt.For = name
		endLoc = nameLoc
	}

	stmt.Loc = startLoc.Merge(endLoc)
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

// parseCreateUser parses (DorisParser.g4:808-810; DorisLegacyParser.g4's
// create_stmt, sql_parser.cup:2100, take the same shape, and the engine runs
// both for a CREATE USER):
//
//	CREATE USER [IF NOT EXISTS] 'user'@'host'
//	    [IDENTIFIED BY 'password' | IDENTIFIED BY PASSWORD 'hash']
//	    [SUPERUSER | DEFAULT ROLE 'role']
//	    [passwordOption]
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

	// Optional IDENTIFIED BY ['PASSWORD'] 'value' (grantUserIdentify).
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
			p.advance()
		} else {
			// IDENTIFIED BY 'password'
			if p.cur.Kind != tokString {
				return nil, p.syntaxErrorAtCur()
			}
			stmt.Password = p.cur.Str
			p.advance()
		}
	}

	// Optional (SUPERUSER | DEFAULT ROLE 'role'). The grammar takes the role
	// as a STRING_LITERAL in DorisParser.g4:809 and in the legacy
	// opt_user_role (DorisLegacyParser.g4:2284), so an identifier is refused.
	switch p.cur.Kind {
	case kwSUPERUSER:
		stmt.Superuser = true
		p.advance()
	case kwDEFAULT:
		p.advance()
		if _, err := p.expect(kwROLE); err != nil {
			return nil, err
		}
		if p.cur.Kind != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.DefaultRole = p.cur.Str
		p.advance()
	}

	// Optional passwordOption.
	opt, err := p.parsePasswordOption()
	if err != nil {
		return nil, err
	}
	stmt.PasswordHistory = opt.history
	stmt.PasswordExpire = opt.expire
	stmt.PasswordReuse = opt.reuse
	stmt.FailedLoginAttempts = opt.failedLoginAttempts
	stmt.PasswordLockTime = opt.lockTime
	stmt.PasswordLockTimeUnit = opt.lockTimeUnit
	stmt.AccountLock = opt.accountLock
	stmt.AccountUnlock = opt.accountUnlock

	// Optional COMMENT
	if p.cur.Kind == kwCOMMENT {
		p.advance()
		if p.cur.Kind != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Comment = p.cur.Str
		p.advance()
	}

	stmt.Loc = startLoc.Merge(p.prev.Loc)
	return stmt, nil
}

// passwordOption is the value of one passwordOption clause (DorisParser.g4:861),
// which DorisLegacyParser.g4's opt_password_option (sql_parser.cup:2289)
// matches option for option. CREATE USER and ALTER USER both carry it.
type passwordOption struct {
	history             string
	expire              string
	reuse               string
	failedLoginAttempts int
	lockTime            int
	lockTimeUnit        string
	accountLock         bool
	accountUnlock       bool
}

// parsePasswordOption reads passwordOption, whose clauses are each optional,
// at most once, and in this order:
//
//	[PASSWORD_HISTORY (DEFAULT | n)]
//	[PASSWORD_EXPIRE (DEFAULT | NEVER | INTERVAL n DAY|HOUR|SECOND)]
//	[PASSWORD_REUSE INTERVAL (DEFAULT | n DAY)]
//	[FAILED_LOGIN_ATTEMPTS n]
//	[PASSWORD_LOCK_TIME (UNBOUNDED | n DAY|HOUR|SECOND)]
//	[ACCOUNT_LOCK | ACCOUNT_UNLOCK]
//
// A value keeps the text after its keyword, as written: "DEFAULT", "NEVER",
// "INTERVAL 1 DAY", "1". Nothing is read when cur does not start one.
func (p *Parser) parsePasswordOption() (passwordOption, error) {
	var opt passwordOption
	timeUnit := func() (string, error) {
		switch p.cur.Kind {
		case kwDAY, kwHOUR, kwSECOND:
			return strings.ToUpper(p.advance().Str), nil
		}
		return "", p.syntaxErrorAtCur()
	}
	if p.cur.Kind == kwPASSWORD_HISTORY {
		p.advance()
		switch p.cur.Kind {
		case kwDEFAULT, tokInt:
			opt.history = strings.ToUpper(p.advance().Str)
		default:
			return opt, p.syntaxErrorAtCur()
		}
	}
	if p.cur.Kind == kwPASSWORD_EXPIRE {
		p.advance()
		switch p.cur.Kind {
		case kwDEFAULT, kwNEVER:
			opt.expire = strings.ToUpper(p.advance().Str)
		case kwINTERVAL:
			p.advance()
			n, err := p.alterExpectInt()
			if err != nil {
				return opt, err
			}
			unit, err := timeUnit()
			if err != nil {
				return opt, err
			}
			opt.expire = "INTERVAL " + n.Str + " " + unit
		default:
			return opt, p.syntaxErrorAtCur()
		}
	}
	if p.cur.Kind == kwPASSWORD_REUSE {
		p.advance()
		if _, err := p.expect(kwINTERVAL); err != nil {
			return opt, err
		}
		switch p.cur.Kind {
		case kwDEFAULT:
			opt.reuse = "DEFAULT"
			p.advance()
		case tokInt:
			n := p.advance()
			if _, err := p.expect(kwDAY); err != nil {
				return opt, err
			}
			opt.reuse = n.Str + " DAY"
		default:
			return opt, p.syntaxErrorAtCur()
		}
	}
	if p.cur.Kind == kwFAILED_LOGIN_ATTEMPTS {
		p.advance()
		n, err := p.alterExpectInt()
		if err != nil {
			return opt, err
		}
		opt.failedLoginAttempts = int(n.Ival)
	}
	if p.cur.Kind == kwPASSWORD_LOCK_TIME {
		p.advance()
		if p.cur.Kind == kwUNBOUNDED {
			opt.lockTimeUnit = "UNBOUNDED"
			p.advance()
		} else {
			n, err := p.alterExpectInt()
			if err != nil {
				return opt, err
			}
			opt.lockTime = int(n.Ival)
			if opt.lockTimeUnit, err = timeUnit(); err != nil {
				return opt, err
			}
		}
	}
	switch p.cur.Kind {
	case kwACCOUNT_LOCK:
		opt.accountLock = true
		p.advance()
	case kwACCOUNT_UNLOCK:
		opt.accountUnlock = true
		p.advance()
	}
	return opt, nil
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

	// passwordOption (DorisParser.g4:861; DorisLegacyParser.g4
	// opt_password_option), each option at most once and in this order, then
	// COMMENT.
	if err := p.parseAlterUserPasswordOptions(stmt); err != nil {
		return nil, err
	}
	if p.cur.Kind == kwCOMMENT {
		p.advance()
		if p.cur.Kind != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Comment = p.cur.Str
		p.advance()
	}
	if p.prev.Loc.End > endLoc.End {
		endLoc = p.prev.Loc
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

// parseAlterUserPasswordOptions reads passwordOption for ALTER USER, which
// DorisParser.g4:861 gives and DorisLegacyParser.g4's opt_password_option
// matches option for option, and stores it on stmt.
func (p *Parser) parseAlterUserPasswordOptions(stmt *ast.AlterUserStmt) error {
	opt, err := p.parsePasswordOption()
	if err != nil {
		return err
	}
	stmt.PasswordHistory = opt.history
	stmt.PasswordExpire = opt.expire
	stmt.PasswordReuse = opt.reuse
	stmt.FailedLoginAttempts = opt.failedLoginAttempts
	stmt.PasswordLockTime = opt.lockTime
	stmt.PasswordLockTimeUnit = opt.lockTimeUnit
	stmt.AccountLock = opt.accountLock
	stmt.AccountUnlock = opt.accountUnlock
	return nil
}
