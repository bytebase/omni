package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// userClauses collects the clauses CREATE USER and ALTER USER share.
type userClauses struct {
	defaultTablespace string
	tempTablespace    string
	localTemp         bool
	quotas            []*nodes.UserQuotaClause
	profile           string
	defaultRole       *nodes.DefaultRoleClause
	passwordExpire    bool
	accountLock       *bool
	enableEditions    bool
	defaultCollation  string
	containerAll      *bool
	options           int // clauses read other than CONTAINER
}

// parseUserClauses parses the user clauses after the user's name and
// IDENTIFIED clause, in any order, until one does not start.
//
// BNF: oracle/parser/bnf/CREATE-USER.bnf, oracle/parser/bnf/ALTER-USER.bnf
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-USER.html
//
//	  DEFAULT TABLESPACE tablespace
//	| DEFAULT ROLE { role [, role ]... | ALL [ EXCEPT role [, role ]... ] | NONE }  -- ALTER
//	| DEFAULT COLLATION collation_name
//	| [ LOCAL ] TEMPORARY TABLESPACE { tablespace | tablespace_group_name }
//	| QUOTA { integer [ K | M | G | T | P | E ] | UNLIMITED } ON tablespace
//	| PROFILE profile
//	| PASSWORD EXPIRE
//	| EXPIRE PASSWORD ROLLOVER PERIOD
//	| ACCOUNT { LOCK | UNLOCK }
//	| ENABLE EDITIONS [ FOR object_type [, object_type ]... ] [ FORCE ]
//	| { ENABLE | DISABLE } DICTIONARY PROTECTION                                  -- ALTER
//	| [ HTTP ] DIGEST { ENABLE | DISABLE }
//	| CONTAINER = { ALL | CURRENT }
//	| READ { ONLY | WRITE }
//	| { SET | ADD | REMOVE } CONTAINER_DATA =
//	      { ALL | DEFAULT | ( container_name [, container_name ]... ) }
//	      [ FOR [ schema. ] object_name ]                                         -- ALTER
func (p *Parser) parseUserClauses(alter bool, c *userClauses) error {
	for {
		var err error
		container := p.isKeywordStr("CONTAINER")
		switch {
		case p.cur.Type == kwDEFAULT:
			err = p.parseUserDefaultClause(alter, c)
		case p.isKeywordStr("LOCAL") || p.cur.Type == kwTEMPORARY:
			c.localTemp = p.isKeywordStr("LOCAL")
			if c.localTemp {
				p.advance()
				if p.cur.Type != kwTEMPORARY {
					return p.syntaxErrorAtCur()
				}
			}
			p.advance() // TEMPORARY
			if p.cur.Type != kwTABLESPACE {
				return p.syntaxErrorAtCur()
			}
			p.advance()
			c.tempTablespace, err = p.parseUserNameWord()
		case p.isKeywordStr("QUOTA"):
			err = p.parseUserQuota(c)
		case p.cur.Type == kwPROFILE:
			p.advance()
			c.profile, err = p.parseUserNameWord()
		case p.isKeywordStr("PASSWORD"):
			p.advance()
			err = p.userWords("EXPIRE")
			c.passwordExpire = true
		case p.isKeywordStr("EXPIRE"):
			p.advance()
			err = p.userWords("PASSWORD", "ROLLOVER", "PERIOD")
		case p.isKeywordStr("ACCOUNT"):
			p.advance()
			switch {
			case p.cur.Type == kwLOCK:
				lock := true
				c.accountLock = &lock
			case p.isKeywordStr("UNLOCK"):
				lock := false
				c.accountLock = &lock
			default:
				return p.syntaxErrorAtCur()
			}
			p.advance()
		case p.cur.Type == kwENABLE && p.isIdentLikeStrAt(p.peekNext(), "EDITIONS"):
			p.advance()
			p.advance()
			c.enableEditions = true
			err = p.parseUserEditionsFor()
		case alter && (p.cur.Type == kwENABLE || p.cur.Type == kwDISABLE):
			p.advance()
			err = p.userWords("DICTIONARY", "PROTECTION")
		case p.isKeywordStr("HTTP") || p.isKeywordStr("DIGEST"):
			if p.isKeywordStr("HTTP") {
				p.advance()
			}
			if err = p.userWords("DIGEST"); err != nil {
				return err
			}
			if p.cur.Type != kwENABLE && p.cur.Type != kwDISABLE {
				return p.syntaxErrorAtCur()
			}
			p.advance()
		case p.isKeywordStr("CONTAINER"):
			c.containerAll, err = p.parseSecurityContainer()
		case p.cur.Type == kwREAD:
			p.advance()
			if !p.isKeywordStr("ONLY") && p.cur.Type != kwWRITE {
				return p.syntaxErrorAtCur()
			}
			p.advance()
		case alter && (p.cur.Type == kwSET || p.isKeywordStr("ADD") || p.isKeywordStr("REMOVE")):
			err = p.parseUserContainerData()
		default:
			return nil
		}
		if err != nil {
			return err
		}
		if !container {
			c.options++
		}
	}
}

// userWords consumes the unquoted words ws in order.
func (p *Parser) userWords(ws ...string) error {
	for _, w := range ws {
		if !p.isKeywordStr(w) {
			return p.syntaxErrorAtCur()
		}
		p.advance()
	}
	return nil
}

// parseUserNameWord parses the name of a tablespace, profile or collation.
func (p *Parser) parseUserNameWord() (string, error) {
	if !p.isIdentLike() {
		return "", p.syntaxErrorAtCur()
	}
	return p.advance().Str, nil
}

// parseUserDefaultClause parses DEFAULT TABLESPACE, DEFAULT ROLE and
// DEFAULT COLLATION.
func (p *Parser) parseUserDefaultClause(alter bool, c *userClauses) error {
	p.advance() // DEFAULT
	var err error
	switch {
	case p.cur.Type == kwTABLESPACE:
		p.advance()
		c.defaultTablespace, err = p.parseUserNameWord()
	case p.isKeywordStr("COLLATION"):
		p.advance()
		c.defaultCollation, err = p.parseUserNameWord()
	case alter && p.cur.Type == kwROLE:
		c.defaultRole, err = p.parseUserDefaultRole()
	default:
		return p.syntaxErrorAtCur()
	}
	return err
}

// parseUserDefaultRole parses DEFAULT ROLE { role [, role ]... | ALL [
// EXCEPT role [, role ]... ] | NONE }. The current token is ROLE.
func (p *Parser) parseUserDefaultRole() (*nodes.DefaultRoleClause, error) {
	clause := &nodes.DefaultRoleClause{Loc: nodes.Loc{Start: p.pos()}}
	p.advance() // ROLE
	switch {
	case p.cur.Type == kwALL:
		p.advance()
		clause.AllRoles = true
		if p.cur.Type == kwEXCEPT {
			p.advance()
			clause.ExceptAll = true
			roles, err := p.parseUserRoleList()
			if err != nil {
				return nil, err
			}
			clause.Roles = roles
		}
	case p.isKeywordStr("NONE"):
		p.advance()
		clause.NoneRole = true
	default:
		roles, err := p.parseUserRoleList()
		if err != nil {
			return nil, err
		}
		clause.Roles = roles
	}
	clause.Loc.End = p.prev.End
	return clause, nil
}

// parseUserRoleList parses role [, role ]...
func (p *Parser) parseUserRoleList() ([]*nodes.ObjectName, error) {
	var roles []*nodes.ObjectName
	for {
		if !p.isIdentLike() {
			return nil, p.syntaxErrorAtCur()
		}
		start := p.pos()
		roles = append(roles, &nodes.ObjectName{Name: p.advance().Str, Loc: nodes.Loc{Start: start, End: p.prev.End}})
		if p.cur.Type != ',' {
			return roles, nil
		}
		p.advance()
	}
}

// parseUserQuota parses QUOTA { integer [ K | M | G | T | P | E ] | UNLIMITED } ON
// tablespace. The current token is QUOTA.
func (p *Parser) parseUserQuota(c *userClauses) error {
	clause := &nodes.UserQuotaClause{Loc: nodes.Loc{Start: p.pos()}}
	p.advance() // QUOTA
	switch {
	case p.isKeywordStr("UNLIMITED"):
		clause.Size = p.advance().Str
	case p.cur.Type == tokICONST:
		clause.Size = p.advance().Str
		if p.isKeywordStr("K") || p.isKeywordStr("M") || p.isKeywordStr("G") || p.isKeywordStr("T") ||
			p.isKeywordStr("P") || p.isKeywordStr("E") {
			clause.Size += p.advance().Str
		}
	default:
		return p.syntaxErrorAtCur()
	}
	if p.cur.Type != kwON {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	if !p.isIdentLike() {
		return p.syntaxErrorAtCur()
	}
	start := p.pos()
	clause.Tablespace = &nodes.ObjectName{Name: p.advance().Str, Loc: nodes.Loc{Start: start, End: p.prev.End}}
	clause.Loc.End = p.prev.End
	c.quotas = append(c.quotas, clause)
	return nil
}

// userEditionableTypes are the object types ENABLE EDITIONS FOR takes.
var userEditionableTypes = newAuditNames("FUNCTION", "LIBRARY", "PACKAGE", "PROCEDURE",
	"SQL TRANSLATION PROFILE", "SYNONYM", "TRIGGER", "TYPE", "VIEW")

// parseUserEditionsFor parses [ FOR object_type [, object_type ]... ] [
// FORCE ] after ENABLE EDITIONS.
func (p *Parser) parseUserEditionsFor() error {
	if p.cur.Type == kwFOR {
		p.advance()
		for {
			_, n := p.matchAuditName(userEditionableTypes)
			if n == 0 {
				return p.syntaxErrorAtCur()
			}
			for range n {
				p.advance()
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}
	if p.cur.Type == kwFORCE {
		p.advance()
	}
	return nil
}

// parseUserContainerData parses the container_data_clause. The current
// token is SET, ADD or REMOVE.
func (p *Parser) parseUserContainerData() error {
	p.advance() // SET, ADD or REMOVE
	if err := p.userWords("CONTAINER_DATA"); err != nil {
		return err
	}
	if p.cur.Type != '=' {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	switch {
	case p.cur.Type == kwALL || p.cur.Type == kwDEFAULT:
		p.advance()
	case p.cur.Type == '(':
		p.advance()
		for {
			if !p.isIdentLike() {
				return p.syntaxErrorAtCur()
			}
			p.advance()
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if p.cur.Type != ')' {
			return p.syntaxErrorAtCur()
		}
		p.advance()
	default:
		return p.syntaxErrorAtCur()
	}
	if p.cur.Type == kwFOR {
		p.advance()
		// The AST records no container data object.
		obj, err := p.parseAuditObjectName()
		if err != nil {
			return err
		}
		_ = obj
	}
	return nil
}

// parseUserIdentified parses the IDENTIFIED clause of CREATE and ALTER
// USER. The current token is IDENTIFIED.
//
//	IDENTIFIED { BY password [ REPLACE old_password ]
//	           | EXTERNALLY [ AS 'certificate_DN' | AS 'kerberos_principal_name' ]
//	           | GLOBALLY [ AS 'directory_DN' ] }
//
// A password is a word, a quoted identifier or an integer; Oracle refuses a
// string literal (ORA-00988).
func (p *Parser) parseUserIdentified(alter bool) (*nodes.IdentifiedClause, error) {
	clause := &nodes.IdentifiedClause{Loc: nodes.Loc{Start: p.pos()}}
	p.advance() // IDENTIFIED
	password := func() (string, error) {
		if p.cur.Type != tokICONST && !p.isIdentLike() {
			return "", p.syntaxErrorAtCur()
		}
		return p.advance().Str, nil
	}
	var err error
	switch {
	case p.cur.Type == kwBY:
		p.advance()
		clause.Type = nodes.IDENTIFIED_BY
		if clause.Password, err = password(); err != nil {
			return nil, err
		}
		if alter && p.cur.Type == kwREPLACE {
			p.advance()
			if clause.OldPass, err = password(); err != nil {
				return nil, err
			}
		}
	case p.isKeywordStr("EXTERNALLY"), p.isKeywordStr("GLOBALLY"):
		clause.Type = nodes.IDENTIFIED_EXTERNALLY
		if p.isKeywordStr("GLOBALLY") {
			clause.Type = nodes.IDENTIFIED_GLOBALLY
		}
		p.advance()
		if p.cur.Type == kwAS {
			p.advance()
			if p.cur.Type != tokSCONST {
				return nil, p.syntaxErrorAtCur()
			}
			clause.ExternalAs = p.advance().Str
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}
	clause.Loc.End = p.prev.End
	return clause, nil
}

// parseUserNoAuthentication parses NO AUTHENTICATION. The current token is
// NO.
func (p *Parser) parseUserNoAuthentication() (*nodes.IdentifiedClause, error) {
	start := p.pos()
	p.advance() // NO
	if err := p.userWords("AUTHENTICATION"); err != nil {
		return nil, err
	}
	return &nodes.IdentifiedClause{Type: nodes.IDENTIFIED_NO_AUTH, Loc: nodes.Loc{Start: start, End: p.prev.End}}, nil
}

// parseUserProxyClause parses the proxy_clause of ALTER USER. The current
// token is GRANT or REVOKE.
//
//	{ GRANT CONNECT THROUGH { ENTERPRISE USERS | db_user_proxy [ proxy_clause_options ] }
//	| REVOKE CONNECT THROUGH { ENTERPRISE USERS | db_user_proxy } }
//
//	proxy_clause_options:
//	    [ WITH { ROLE { role [, role ]... | ALL EXCEPT role [, role ]... } | NO ROLES } ]
//	    [ AUTHENTICATION REQUIRED ]
//
// Oracle takes ROLES for ROLE.
func (p *Parser) parseUserProxyClause() error {
	grant := p.advance().Type == kwGRANT
	if p.cur.Type != kwCONNECT {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	if err := p.userWords("THROUGH"); err != nil {
		return err
	}
	if p.isKeywordStr("ENTERPRISE") && p.isIdentLikeStrAt(p.peekNext(), "USERS") {
		p.advance()
		p.advance()
		return nil
	}
	if !p.isIdentLike() {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	if !grant {
		return nil
	}
	if p.cur.Type == kwWITH {
		p.advance()
		switch {
		case p.isKeywordStr("NO"):
			p.advance()
			if err := p.userWords("ROLES"); err != nil {
				return err
			}
		case p.cur.Type == kwROLE || p.isKeywordStr("ROLES"):
			p.advance()
			if p.cur.Type == kwALL {
				p.advance()
				if p.cur.Type != kwEXCEPT {
					return p.syntaxErrorAtCur()
				}
				p.advance()
			}
			for {
				if !p.isAuditRoleName() {
					return p.syntaxErrorAtCur()
				}
				p.advance()
				if p.cur.Type != ',' {
					break
				}
				p.advance()
			}
		default:
			return p.syntaxErrorAtCur()
		}
	}
	if p.isKeywordStr("AUTHENTICATION") {
		p.advance()
		if err := p.userWords("REQUIRED"); err != nil {
			return err
		}
	}
	return nil
}
