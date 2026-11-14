// Package parser - security_principals.go implements T-SQL security principal
// statement parsing: CREATE/ALTER/DROP USER, LOGIN, ROLE, APPLICATION ROLE,
// and ADD/DROP ROLE MEMBER.
package parser

import (
	"fmt"
	"math"
	"slices"
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// SQL Server reads the words of these statements by their text and refuses
// them bracketed or quoted ([PASSWORD], [ENABLE], [MEMBER], [ON]), so the
// helpers below match words with prinIsWord and option values by the form the
// grammar gives each option (TSql160.g:9240-9500, 10829-10990, 11097-11150,
// 11938-12080, 12519-12600, 12776-12880), with the value forms and the option
// sets that the pinned SQL Server 2022 takes.

// prinIsWord reports whether the current token is word, written without
// brackets or quotes.
func (p *Parser) prinIsWord(word string) bool {
	if kw := lookupKeyword(word); kw != tokIDENT {
		return p.cur.Type == kw
	}
	return p.cur.Type == tokIDENT && !p.isDelimited(p.cur) && strings.ToUpper(p.cur.Str) == word
}

// prinExpectWord consumes word or returns a syntax error.
func (p *Parser) prinExpectWord(word string) error {
	if !p.prinIsWord(word) {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// prinName reads a principal, schema, database, credential, certificate or
// key name: an identifier, a reserved word only when delimited (TSql160.g
// identifier), so [default] but not default, 'a' or NULL.
func (p *Parser) prinName() (string, error) {
	if !p.isIdentLike() {
		return "", p.syntaxErrorAtCur()
	}
	name := p.cur.Str
	p.advance()
	return name, nil
}

// prinIsString reports whether the current token is a string literal.
func (p *Parser) prinIsString() bool {
	return p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST
}

// prinIsBinary reports whether the current token is a hex literal (TSql160.g
// binary: HexLiteral), 0x included.
func (p *Parser) prinIsBinary() bool {
	return (p.cur.Type == tokICONST || p.cur.Type == tokFCONST) && strings.HasPrefix(strings.ToLower(p.cur.Str), "0x")
}

// prinValueKind is the form of an option's value.
type prinValueKind int

const (
	prinIdent        prinValueKind = iota // identifier
	prinString                            // string literal
	prinBinary                            // hex literal
	prinOnOff                             // ON | OFF
	prinOnOffOrIdent                      // ON | OFF | identifier
	prinIdentOrInt                        // identifier | Int32 (a language name or LCID)
	prinIdentOrNull                       // identifier | NULL
	prinPassword                          // see prinPasswordForm
	prinNoCredential                      // NO CREDENTIAL, which has no "="
)

// prinPasswordForm is what a PASSWORD option takes in a statement.
type prinPasswordForm struct {
	binary      bool     // a hex password, which needs HASHED when HASHED is among flags
	oldPassword bool     // OLD_PASSWORD = 'old' may follow, instead of flags
	flags       []string // words that may follow, each once
}

// prinOpt is an option word and the form of its value.
type prinOpt struct {
	word     string
	kind     prinValueKind
	password prinPasswordForm
}

// prinValue reads a value of the given form.
func (p *Parser) prinValue(kind prinValueKind) (string, error) {
	switch kind {
	case prinString:
		if !p.prinIsString() {
			return "", p.syntaxErrorAtCur()
		}
	case prinBinary:
		if !p.prinIsBinary() {
			return "", p.syntaxErrorAtCur()
		}
	case prinOnOff, prinOnOffOrIdent:
		if p.prinIsWord("ON") || p.prinIsWord("OFF") {
			v := strings.ToUpper(p.cur.Str)
			p.advance()
			return v, nil
		}
		if kind == prinOnOff {
			return "", p.syntaxErrorAtCur()
		}
		return p.prinName()
	case prinIdentOrInt:
		if p.cur.Type == tokICONST && !p.prinIsBinary() && p.cur.Ival <= math.MaxInt32 {
			break
		}
		return p.prinName()
	case prinIdentOrNull:
		if p.cur.Type == kwNULL {
			p.advance()
			return "NULL", nil
		}
		return p.prinName()
	default:
		return p.prinName()
	}
	v := p.cur.Str
	p.advance()
	return v, nil
}

// prinPasswordValue reads what follows PASSWORD =: a string, or a hex literal
// where form allows one, then OLD_PASSWORD = 'old' or the flags form allows.
// SQL Server takes a hex password only HASHED where HASHED may follow, each
// flag once, and OLD_PASSWORD only without flags.
func (p *Parser) prinPasswordValue(opt *nodes.SecurityPrincipalOption, form prinPasswordForm) error {
	valLoc := p.pos()
	isBinary := p.prinIsBinary()
	if !p.prinIsString() && !(isBinary && form.binary) {
		return p.syntaxErrorAtCur()
	}
	opt.Value = p.cur.Str
	p.advance()
	if form.oldPassword && p.prinIsWord("OLD_PASSWORD") {
		p.advance()
		if _, err := p.expect('='); err != nil {
			return err
		}
		if !p.prinIsString() {
			return p.syntaxErrorAtCur()
		}
		opt.OldPassword = p.cur.Str
		p.advance()
	} else {
		for {
			flag := ""
			for _, f := range form.flags {
				if p.prinIsWord(f) {
					flag = f
					break
				}
			}
			if flag == "" {
				break
			}
			var set *bool
			switch flag {
			case "HASHED":
				set = &opt.Hashed
			case "MUST_CHANGE":
				set = &opt.MustChange
			default:
				set = &opt.Unlock
			}
			if *set {
				return p.syntaxErrorAtCur()
			}
			*set = true
			p.advance()
		}
	}
	if isBinary && slices.Contains(form.flags, "HASHED") && !opt.Hashed {
		return p.newParseError(valLoc, "a hex password must be HASHED")
	}
	return nil
}

// prinOptions reads opt [ , opt ]..., each option word one of opts, each at
// most once (SQL Server refuses a repeat).
func (p *Parser) prinOptions(opts []prinOpt) ([]nodes.Node, error) {
	var out []nodes.Node
	seen := map[string]bool{}
	for {
		optLoc := p.pos()
		var o *prinOpt
		for i := range opts {
			if p.prinIsWord(opts[i].word) {
				o = &opts[i]
				break
			}
		}
		if o == nil {
			return nil, p.syntaxErrorAtCur()
		}
		// NO CREDENTIAL and CREDENTIAL = c are one option to SQL Server.
		key := o.word
		if o.kind == prinNoCredential {
			key = "CREDENTIAL"
		}
		if seen[key] {
			return nil, p.newParseError(optLoc, fmt.Sprintf("option '%s' is specified more than once", key))
		}
		seen[key] = true
		p.advance()
		opt := &nodes.SecurityPrincipalOption{Name: o.word, Loc: nodes.Loc{Start: optLoc, End: -1}}
		switch o.kind {
		case prinNoCredential:
			if err := p.prinExpectWord("CREDENTIAL"); err != nil {
				return nil, err
			}
			opt.Name = "NO CREDENTIAL"
		case prinPassword:
			if _, err := p.expect('='); err != nil {
				return nil, err
			}
			if err := p.prinPasswordValue(opt, o.password); err != nil {
				return nil, err
			}
		default:
			if _, err := p.expect('='); err != nil {
				return nil, err
			}
			v, err := p.prinValue(o.kind)
			if err != nil {
				return nil, err
			}
			opt.Value = v
		}
		opt.Loc.End = p.prevEnd()
		out = append(out, opt)
		if p.cur.Type != ',' {
			return out, nil
		}
		p.advance()
	}
}

// prinIfExists consumes IF EXISTS, which DROP USER and DROP ROLE take
// (TSql160.g:15342, 15482).
func (p *Parser) prinIfExists() {
	if p.cur.Type == kwIF && p.peekNext().Type == kwEXISTS {
		p.advance()
		p.advance()
	}
}

// prinFlag returns a valueless option node that spans from start to the
// previous token.
func (p *Parser) prinFlag(start int, name, value string) *nodes.SecurityPrincipalOption {
	return &nodes.SecurityPrincipalOption{Name: name, Value: value, Loc: nodes.Loc{Start: start, End: p.prevEnd()}}
}

// The option sets SQL Server 2022 takes in each statement.
var (
	// CREATE USER (TSql160.g:10911). PASSWORD only for a user with no login
	// source; DEFAULT_LANGUAGE takes an LCID too; SQL Server reads
	// ALLOW_ENCRYPTED_VALUE_MODIFICATIONS as any identifier here.
	prinCreateUserOpts = []prinOpt{
		{word: "DEFAULT_SCHEMA", kind: prinIdent},
		{word: "DEFAULT_LANGUAGE", kind: prinIdentOrInt},
		{word: "SID", kind: prinBinary},
		{word: "TYPE", kind: prinIdent},
		{word: "ALLOW_ENCRYPTED_VALUE_MODIFICATIONS", kind: prinOnOffOrIdent},
	}
	prinCreateUserPassword = prinOpt{word: "PASSWORD", kind: prinPassword}

	// ALTER USER (TSql160.g:12776-12880).
	prinAlterUserOpts = []prinOpt{
		{word: "NAME", kind: prinIdent},
		{word: "LOGIN", kind: prinIdent},
		{word: "DEFAULT_SCHEMA", kind: prinIdentOrNull},
		{word: "DEFAULT_LANGUAGE", kind: prinIdentOrInt},
		{word: "PASSWORD", kind: prinPassword, password: prinPasswordForm{binary: true, oldPassword: true}},
		{word: "ALLOW_ENCRYPTED_VALUE_MODIFICATIONS", kind: prinOnOff},
	}

	// CREATE LOGIN ... WITH PASSWORD = ..., <these> (TSql160.g:9316).
	prinCreateLoginOpts = []prinOpt{
		{word: "DEFAULT_DATABASE", kind: prinIdent},
		{word: "DEFAULT_LANGUAGE", kind: prinIdent},
		{word: "CHECK_EXPIRATION", kind: prinOnOff},
		{word: "CHECK_POLICY", kind: prinOnOff},
		{word: "CREDENTIAL", kind: prinIdent},
		{word: "SID", kind: prinBinary},
	}
	prinCreateLoginPassword = prinPasswordForm{binary: true, flags: []string{"HASHED", "MUST_CHANGE"}}

	// CREATE LOGIN ... FROM WINDOWS WITH (TSql160.g:9383).
	prinWindowsLoginOpts = []prinOpt{
		{word: "DEFAULT_DATABASE", kind: prinIdent},
		{word: "DEFAULT_LANGUAGE", kind: prinIdent},
	}

	// CREATE LOGIN ... FROM EXTERNAL PROVIDER WITH (TSql160.g:9432), and
	// OBJECT_ID, which ScriptDOM lacks and SQL Server 2022 takes.
	prinExternalLoginOpts = []prinOpt{
		{word: "OBJECT_ID", kind: prinString},
		{word: "TYPE", kind: prinIdent},
		{word: "SID", kind: prinBinary},
		{word: "DEFAULT_DATABASE", kind: prinIdent},
		{word: "DEFAULT_LANGUAGE", kind: prinIdent},
	}

	// ALTER LOGIN ... WITH (TSql160.g:12016-12070).
	prinAlterLoginOpts = []prinOpt{
		{word: "PASSWORD", kind: prinPassword, password: prinPasswordForm{binary: true, oldPassword: true, flags: []string{"HASHED", "MUST_CHANGE", "UNLOCK"}}},
		{word: "DEFAULT_DATABASE", kind: prinIdent},
		{word: "DEFAULT_LANGUAGE", kind: prinIdent},
		{word: "NAME", kind: prinIdent},
		{word: "CHECK_POLICY", kind: prinOnOff},
		{word: "CHECK_EXPIRATION", kind: prinOnOff},
		{word: "CREDENTIAL", kind: prinIdent},
		{word: "NO", kind: prinNoCredential},
	}

	// CREATE APPLICATION ROLE (TSql160.g:11097): PASSWORD and DEFAULT_SCHEMA.
	prinCreateAppRoleOpts = []prinOpt{
		{word: "PASSWORD", kind: prinPassword},
		{word: "DEFAULT_SCHEMA", kind: prinIdent},
	}

	// ALTER APPLICATION ROLE: NAME and LOGIN as well.
	prinAlterAppRoleOpts = []prinOpt{
		{word: "PASSWORD", kind: prinPassword},
		{word: "DEFAULT_SCHEMA", kind: prinIdent},
		{word: "NAME", kind: prinIdent},
		{word: "LOGIN", kind: prinIdent},
	}
)

// parseSecurityUserStmt parses CREATE/ALTER/DROP USER.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-user-transact-sql
// Grammar: TSql160.g:10829-10990 (CREATE), 12776-12880 (ALTER), 15482 (DROP).
//
//	CREATE USER user_name
//	    [ { FOR | FROM } { LOGIN login_name | CERTIFICATE cert_name
//	                     | ASYMMETRIC KEY asym_key_name | EXTERNAL PROVIDER }
//	    | WITHOUT LOGIN ]
//	    [ WITH <option> [ ,...n ] ]
//	ALTER USER user_name WITH <option> [ ,...n ]
//	DROP USER [ IF EXISTS ] user_name
func (p *Parser) parseSecurityUserStmt(action string) (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     action,
		ObjectType: "USER",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	if action == "DROP" {
		p.prinIfExists()
	}
	name, err := p.prinName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	if action == "DROP" {
		stmt.Loc.End = p.prevEnd()
		return stmt, nil
	}

	var opts []nodes.Node
	hasSource := false
	if action == "CREATE" {
		optLoc := p.pos()
		switch {
		case p.cur.Type == kwFOR || p.cur.Type == kwFROM:
			hasSource = true
			p.advance()
			switch {
			case p.prinIsWord("LOGIN"), p.prinIsWord("CERTIFICATE"):
				kind := strings.ToUpper(p.cur.Str)
				p.advance()
				v, err := p.prinName()
				if err != nil {
					return nil, err
				}
				opts = append(opts, p.prinFlag(optLoc, kind, v))
			case p.prinIsWord("ASYMMETRIC"):
				p.advance()
				if err := p.prinExpectWord("KEY"); err != nil {
					return nil, err
				}
				v, err := p.prinName()
				if err != nil {
					return nil, err
				}
				opts = append(opts, p.prinFlag(optLoc, "ASYMMETRIC KEY", v))
			case p.prinIsWord("EXTERNAL"):
				p.advance()
				if err := p.prinExpectWord("PROVIDER"); err != nil {
					return nil, err
				}
				opts = append(opts, p.prinFlag(optLoc, "EXTERNAL PROVIDER", ""))
			default:
				return nil, p.syntaxErrorAtCur()
			}
		case p.prinIsWord("WITHOUT"):
			hasSource = true
			p.advance()
			if err := p.prinExpectWord("LOGIN"); err != nil {
				return nil, err
			}
			opts = append(opts, p.prinFlag(optLoc, "WITHOUT LOGIN", ""))
		}
	}

	if p.cur.Type == kwWITH {
		p.advance()
		set := prinAlterUserOpts
		if action == "CREATE" {
			set = prinCreateUserOpts
			if !hasSource {
				set = append([]prinOpt{prinCreateUserPassword}, prinCreateUserOpts...)
			}
		}
		more, err := p.prinOptions(set)
		if err != nil {
			return nil, err
		}
		opts = append(opts, more...)
	} else if action == "ALTER" {
		return nil, p.syntaxErrorAtCur()
	}

	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseSecurityLoginStmt parses CREATE/ALTER/DROP LOGIN.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-login-transact-sql
// Grammar: TSql160.g:9240-9500 (CREATE), 11938-12070 (ALTER), 15331 (DROP).
//
//	CREATE LOGIN login_name
//	    { WITH PASSWORD = { 'password' | hashed_password HASHED } [ MUST_CHANGE ]
//	        [ , <option> [ ,...n ] ]
//	    | FROM WINDOWS [ WITH { DEFAULT_DATABASE | DEFAULT_LANGUAGE } = name [ ,...n ] ]
//	    | FROM EXTERNAL PROVIDER [ WITH <external_option> [ ,...n ] ]
//	    | FROM CERTIFICATE certname [ WITH CREDENTIAL = credential_name ]
//	    | FROM ASYMMETRIC KEY asym_key_name [ WITH CREDENTIAL = credential_name ]
//	    }
//	ALTER LOGIN login_name
//	    { ENABLE | DISABLE
//	    | WITH <set_option> [ ,...n ]
//	    | { ADD | DROP } CREDENTIAL credential_name
//	    }
//	DROP LOGIN login_name
func (p *Parser) parseSecurityLoginStmt(action string) (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     action,
		ObjectType: "LOGIN",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	name, err := p.prinName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	var opts []nodes.Node
	switch action {
	case "CREATE":
		opts, err = p.prinCreateLoginSource()
	case "ALTER":
		opts, err = p.prinAlterLoginAction()
	}
	if err != nil {
		return nil, err
	}

	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// prinCreateLoginSource reads what follows CREATE LOGIN login_name: WITH
// PASSWORD, which comes first, or FROM and a source.
func (p *Parser) prinCreateLoginSource() ([]nodes.Node, error) {
	if p.cur.Type == kwWITH {
		p.advance()
		optLoc := p.pos()
		if err := p.prinExpectWord("PASSWORD"); err != nil {
			return nil, err
		}
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		pw := &nodes.SecurityPrincipalOption{Name: "PASSWORD", Loc: nodes.Loc{Start: optLoc, End: -1}}
		if err := p.prinPasswordValue(pw, prinCreateLoginPassword); err != nil {
			return nil, err
		}
		pw.Loc.End = p.prevEnd()
		opts := []nodes.Node{pw}
		if p.cur.Type == ',' {
			p.advance()
			more, err := p.prinOptions(prinCreateLoginOpts)
			if err != nil {
				return nil, err
			}
			opts = append(opts, more...)
		}
		return opts, nil
	}

	optLoc := p.pos()
	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}
	var source string
	var withOpts []prinOpt
	switch {
	case p.prinIsWord("WINDOWS"):
		p.advance()
		source = "WINDOWS"
		withOpts = prinWindowsLoginOpts
	case p.prinIsWord("EXTERNAL"):
		p.advance()
		if err := p.prinExpectWord("PROVIDER"); err != nil {
			return nil, err
		}
		source = "EXTERNAL PROVIDER"
		withOpts = prinExternalLoginOpts
	case p.prinIsWord("CERTIFICATE"), p.prinIsWord("ASYMMETRIC"):
		source = "CERTIFICATE"
		if p.prinIsWord("ASYMMETRIC") {
			p.advance()
			if !p.prinIsWord("KEY") {
				return nil, p.syntaxErrorAtCur()
			}
			source = "ASYMMETRIC KEY"
		}
		p.advance()
		keyName, err := p.prinName()
		if err != nil {
			return nil, err
		}
		source += " " + keyName
		withOpts = []prinOpt{{word: "CREDENTIAL", kind: prinIdent}}
	default:
		return nil, p.syntaxErrorAtCur()
	}
	opts := []nodes.Node{p.prinFlag(optLoc, "FROM", source)}
	if p.cur.Type == kwWITH {
		p.advance()
		more, err := p.prinOptions(withOpts)
		if err != nil {
			return nil, err
		}
		opts = append(opts, more...)
	}
	return opts, nil
}

// prinAlterLoginAction reads what follows ALTER LOGIN login_name.
func (p *Parser) prinAlterLoginAction() ([]nodes.Node, error) {
	optLoc := p.pos()
	switch {
	case p.prinIsWord("ENABLE"), p.prinIsWord("DISABLE"):
		name := strings.ToUpper(p.cur.Str)
		p.advance()
		return []nodes.Node{p.prinFlag(optLoc, name, "")}, nil
	case p.cur.Type == kwADD, p.cur.Type == kwDROP:
		name := strings.ToUpper(p.cur.Str) + " CREDENTIAL"
		p.advance()
		if err := p.prinExpectWord("CREDENTIAL"); err != nil {
			return nil, err
		}
		v, err := p.prinName()
		if err != nil {
			return nil, err
		}
		return []nodes.Node{p.prinFlag(optLoc, name, v)}, nil
	case p.cur.Type == kwWITH:
		p.advance()
		return p.prinOptions(prinAlterLoginOpts)
	}
	return nil, p.syntaxErrorAtCur()
}

// prinRoleAction reads the action of ALTER ROLE and ALTER SERVER ROLE
// (TSql160.g:12553-12590): ADD MEMBER name, DROP MEMBER name or WITH NAME =
// name. It returns the action as ADD MEMBER, DROP MEMBER or NAME.
func (p *Parser) prinRoleAction() (string, string, error) {
	var action string
	switch {
	case p.cur.Type == kwADD, p.cur.Type == kwDROP:
		action = strings.ToUpper(p.cur.Str) + " MEMBER"
		p.advance()
		if err := p.prinExpectWord("MEMBER"); err != nil {
			return "", "", err
		}
	case p.cur.Type == kwWITH:
		action = "NAME"
		p.advance()
		if err := p.prinExpectWord("NAME"); err != nil {
			return "", "", err
		}
		if _, err := p.expect('='); err != nil {
			return "", "", err
		}
	default:
		return "", "", p.syntaxErrorAtCur()
	}
	v, err := p.prinName()
	if err != nil {
		return "", "", err
	}
	return action, v, nil
}

// prinAuthorization reads an optional AUTHORIZATION owner_name (TSql160.g
// authorizationOpt) and returns the owner, or "" without one.
func (p *Parser) prinAuthorization() (string, bool, error) {
	if p.cur.Type != kwAUTHORIZATION {
		return "", false, nil
	}
	p.advance()
	owner, err := p.prinName()
	if err != nil {
		return "", false, err
	}
	return owner, true, nil
}

// parseSecurityRoleStmt parses CREATE/ALTER/DROP ROLE and ADD/DROP ROLE MEMBER.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-role-transact-sql
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-role-transact-sql
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-role-transact-sql
// Grammar: TSql160.g:10287, 12519, 15342.
//
//	CREATE ROLE role_name [ AUTHORIZATION owner_name ]
//	ALTER  ROLE role_name { ADD MEMBER member_name | DROP MEMBER member_name | WITH NAME = new_name }
//	DROP   ROLE [ IF EXISTS ] role_name
func (p *Parser) parseSecurityRoleStmt(action string) (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     action,
		ObjectType: "ROLE",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	if action == "DROP" {
		p.prinIfExists()
	}
	name, err := p.prinName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	optLoc := p.pos()
	switch action {
	case "CREATE":
		owner, ok, err := p.prinAuthorization()
		if err != nil {
			return nil, err
		}
		if ok {
			stmt.Options = &nodes.List{Items: []nodes.Node{p.prinFlag(optLoc, "AUTHORIZATION", owner)}}
		}
	case "ALTER":
		act, v, err := p.prinRoleAction()
		if err != nil {
			return nil, err
		}
		stmt.Options = &nodes.List{Items: []nodes.Node{p.prinFlag(optLoc, act, v)}}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseSecurityApplicationRoleStmt parses CREATE/ALTER/DROP APPLICATION ROLE.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-application-role-transact-sql
// Grammar: TSql160.g:11097-11150, 15284.
//
//	CREATE APPLICATION ROLE role_name WITH <option> [ ,...n ]
//	    <option> ::= PASSWORD = 'password' | DEFAULT_SCHEMA = schema_name
//	ALTER  APPLICATION ROLE role_name WITH <option> [ ,...n ]
//	    <option> ::= ... | NAME = new_name | LOGIN = login_name
//	DROP   APPLICATION ROLE role_name
func (p *Parser) parseSecurityApplicationRoleStmt(action string) (*nodes.SecurityStmt, error) {
	loc := p.pos()
	// APPLICATION keyword already consumed by caller
	if err := p.prinExpectWord("ROLE"); err != nil {
		return nil, err
	}

	stmt := &nodes.SecurityStmt{
		Action:     action,
		ObjectType: "APPLICATION ROLE",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	name, err := p.prinName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	if action != "DROP" {
		if _, err := p.expect(kwWITH); err != nil {
			return nil, err
		}
		set := prinCreateAppRoleOpts
		if action == "ALTER" {
			set = prinAlterAppRoleOpts
		}
		opts, err := p.prinOptions(set)
		if err != nil {
			return nil, err
		}
		stmt.Options = &nodes.List{Items: opts}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseExecuteAsStmt parses a standalone EXECUTE AS statement (context switching).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/execute-as-transact-sql
//
//	EXECUTE AS { LOGIN | USER | CALLER | SELF | OWNER } [ = 'name' ]
//	    [ WITH { NO REVERT | COOKIE INTO @cookie_variable } ]
func (p *Parser) parseExecuteAsStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	p.advance() // consume EXECUTE
	p.advance() // consume AS

	stmt := &nodes.SecurityStmt{
		Action: "EXECUTE AS",
		Loc:    nodes.Loc{Start: loc, End: -1},
	}

	// { LOGIN | USER | CALLER | SELF | OWNER }
	if p.cur.Type == kwLOGIN {
		stmt.ObjectType = "LOGIN"
		p.advance()
	} else if p.cur.Type == kwUSER {
		stmt.ObjectType = "USER"
		p.advance()
	} else if p.cur.Type == kwCALLER {
		stmt.ObjectType = "CALLER"
		p.advance()
		stmt.Loc.End = p.prevEnd()
		return stmt, nil
	} else if p.cur.Type == kwSELF {
		stmt.ObjectType = "SELF"
		p.advance()
		stmt.Loc.End = p.prevEnd()
		return stmt, nil
	} else if p.cur.Type == kwOWNER {
		stmt.ObjectType = "OWNER"
		p.advance()
		stmt.Loc.End = p.prevEnd()
		return stmt, nil
	}

	// = 'name'
	if p.cur.Type == '=' {
		p.advance()
		if p.cur.Type == tokSCONST || p.isIdentLike() {
			stmt.Name = p.cur.Str
			p.advance()
		}
	}

	// WITH { NO REVERT | COOKIE INTO @cookie_variable }
	if p.cur.Type == kwWITH {
		p.advance()
		var opts []nodes.Node
		if p.cur.Type == kwNO {
			p.advance()
			if p.cur.Type == kwREVERT {
				p.advance()
				opts = append(opts, &nodes.String{Str: "NO REVERT"})
			}
		} else if p.cur.Type == kwCOOKIE {
			p.advance()
			if p.cur.Type == kwINTO {
				p.advance()
			}
			if p.cur.Type == tokVARIABLE {
				opts = append(opts, &nodes.String{Str: "COOKIE=" + p.cur.Str})
				p.advance()
			}
		}
		if len(opts) > 0 {
			stmt.Options = &nodes.List{Items: opts}
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseRevertStmt parses the REVERT statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/revert-transact-sql
//
//	REVERT [ WITH COOKIE = @cookie_variable ]
func (p *Parser) parseRevertStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	p.advance() // consume REVERT

	stmt := &nodes.SecurityStmt{
		Action:     "REVERT",
		ObjectType: "CONTEXT",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// WITH COOKIE = expression (revertStatement, TSql160.g:32977): SQL
	// Server takes a constant there as well as a variable.
	if p.cur.Type == kwWITH {
		p.advance()
		if _, err := p.expect(kwCOOKIE); err != nil {
			return nil, err
		}
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		valLoc := p.pos()
		value, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if value == nil {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Options = &nodes.List{Items: []nodes.Node{&nodes.String{Str: "COOKIE=" + p.source[valLoc:p.prevEnd()]}}}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterAuthorizationStmt parses ALTER AUTHORIZATION.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-authorization-transact-sql
//
//	ALTER AUTHORIZATION
//	    ON [ <class_type> :: ] entity_name
//	    TO { principal_name | SCHEMA OWNER }
//
//	<class_type> ::=
//	    { OBJECT | ASSEMBLY | ASYMMETRIC KEY | CERTIFICATE | CONTRACT
//	    | DATABASE | ENDPOINT | FULLTEXT CATALOG | FULLTEXT STOPLIST
//	    | MESSAGE TYPE | REMOTE SERVICE BINDING | ROLE | ROUTE
//	    | SCHEMA | SEARCH PROPERTY LIST | SERVER ROLE | SERVICE
//	    | SYMMETRIC KEY | TYPE | XML SCHEMA COLLECTION }
func (p *Parser) parseAlterAuthorizationStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	// AUTHORIZATION keyword already consumed by caller

	stmt := &nodes.SecurityStmt{
		Action:     "ALTER AUTHORIZATION",
		ObjectType: "ALTER AUTHORIZATION",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// ON [ entity_type :: ] entity_name
	if p.cur.Type == kwON {
		p.advance()
		// Try to detect entity_type :: pattern
		// Entity types can be multi-word: XML SCHEMA COLLECTION, SEARCH PROPERTY LIST, etc.
		entityType := p.tryParseAlterAuthEntityType()
		if entityType != "" {
			// We have entity_type, now expect ::
			if p.cur.Type == tokCOLONCOLON {
				p.advance() // consume ::
			}
			stmt.Name = entityType + "::" + p.parseAlterAuthEntityName()
		} else {
			// No entity_type, just parse entity_name
			stmt.Name = p.parseAlterAuthEntityName()
		}
	}

	// TO { principal_name | SCHEMA OWNER }
	if p.cur.Type == kwTO {
		p.advance()
		var opts []nodes.Node
		if p.cur.Type == kwSCHEMA {
			p.advance()
			if p.cur.Type == kwOWNER {
				p.advance()
				opts = append(opts, &nodes.String{Str: "SCHEMA OWNER"})
			}
		} else if p.isIdentLike() || p.cur.Type == tokSCONST {
			opts = append(opts, &nodes.String{Str: p.cur.Str})
			p.advance()
		}
		if len(opts) > 0 {
			stmt.Options = &nodes.List{Items: opts}
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// tryParseAlterAuthEntityType checks if the current position has an entity_type :: pattern.
// If so, consumes the entity type tokens and returns the entity type string.
// Otherwise returns empty string without consuming anything.
//
// Known single-word entity types: OBJECT, ASSEMBLY, CONTRACT, CERTIFICATE, ENDPOINT,
// ROUTE, SERVICE, TYPE, SCHEMA, DATABASE, ROLE
// Multi-word: XML SCHEMA COLLECTION, MESSAGE TYPE, REMOTE SERVICE BINDING,
// FULLTEXT CATALOG, FULLTEXT STOPLIST, SEARCH PROPERTY LIST,
// SYMMETRIC KEY, ASYMMETRIC KEY, AVAILABILITY GROUP, SERVER ROLE
func (p *Parser) tryParseAlterAuthEntityType() string {
	// For single-word entity types, check if next token is ':'
	if (p.isIdentLike() || p.cur.Type == kwDATABASE || p.cur.Type == kwXML ||
		p.cur.Type == kwROLE || p.cur.Type == kwSCHEMA || p.cur.Type == kwTYPE) &&
		p.peekNext().Type == tokCOLONCOLON {
		word := strings.ToUpper(p.cur.Str)
		p.advance() // consume entity type word
		return word
	}
	// For multi-word entity types: check first word, then peek second
	if p.isIdentLike() || p.cur.Type == kwXML || p.cur.Type == kwSCHEMA {
		first := strings.ToUpper(p.cur.Str)
		switch first {
		case "XML":
			return p.tryParseMultiWordEntityType(first, []string{"SCHEMA", "COLLECTION"})
		case "MESSAGE":
			return p.tryParseMultiWordEntityType(first, []string{"TYPE"})
		case "REMOTE":
			return p.tryParseMultiWordEntityType(first, []string{"SERVICE", "BINDING"})
		case "FULLTEXT":
			next := p.peekNext()
			if next.Type == kwCATALOG || next.Type == kwSTOPLIST {
				return p.tryParseMultiWordEntityType(first, []string{strings.ToUpper(next.Str)})
			}
		case "SEARCH":
			return p.tryParseMultiWordEntityType(first, []string{"PROPERTY", "LIST"})
		case "SYMMETRIC", "ASYMMETRIC":
			return p.tryParseMultiWordEntityType(first, []string{"KEY"})
		case "AVAILABILITY":
			return p.tryParseMultiWordEntityType(first, []string{"GROUP"})
		case "SERVER":
			return p.tryParseMultiWordEntityType(first, []string{"ROLE"})
		}
	}
	return ""
}

// tryParseMultiWordEntityType consumes a multi-word entity type if the next word matches.
// Validates by peeking before consuming. Returns the entity type or "".
func (p *Parser) tryParseMultiWordEntityType(first string, remaining []string) string {
	if len(remaining) == 0 {
		return ""
	}
	// Peek to confirm next token matches first remaining word
	next := p.peekNext()
	if next.Type != lookupKeyword(remaining[0]) {
		return ""
	}

	// Pattern looks good, consume first word
	p.advance()
	parts := []string{first}
	for _, word := range remaining {
		if (p.isIdentLike() || p.cur.Type == kwSCHEMA || p.cur.Type == kwTYPE) &&
			p.cur.Type == lookupKeyword(word) {
			parts = append(parts, strings.ToUpper(p.cur.Str))
			p.advance()
		} else {
			break
		}
	}
	return strings.Join(parts, " ")
}

// parseAlterAuthEntityName parses the entity name (possibly qualified with dots) until TO keyword.
func (p *Parser) parseAlterAuthEntityName() string {
	var parts []string
	for p.cur.Type != tokEOF && p.cur.Type != ';' {
		if p.cur.Type == kwTO {
			break
		}
		if p.isIdentLike() || p.cur.Type == kwDATABASE || p.cur.Type == kwXML ||
			p.cur.Type == kwROLE || p.cur.Type == kwSCHEMA || p.cur.Type == kwTYPE {
			parts = append(parts, p.cur.Str)
		} else if p.cur.Type == '.' {
			parts = append(parts, ".")
		} else {
			break
		}
		p.advance()
	}
	return strings.Join(parts, " ")
}
