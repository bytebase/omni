package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseIdentifiedClause parses an IDENTIFIED clause for CREATE/ALTER USER.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/CREATE-USER.html
//
//	IDENTIFIED
//	    { BY password
//	    | EXTERNALLY [ AS 'certificate_DN' | AS 'kerberos_principal_name' ]
//	    | GLOBALLY [ AS 'directory_DN' ]
//	    }
func (p *Parser) parseIdentifiedClause(allowReplace bool) (*nodes.IdentifiedClause, error) {
	start := p.pos()
	clause := &nodes.IdentifiedClause{
		Loc: nodes.Loc{Start: start},
	}

	// Already consumed IDENTIFIED keyword
	if p.cur.Type == kwBY {
		// IDENTIFIED BY password
		p.advance()
		clause.Type = nodes.IDENTIFIED_BY
		if p.cur.Type == tokSCONST || p.isIdentLike() {
			clause.Password = p.cur.Str
			p.advance()
		}
		// REPLACE old_password (ALTER USER only)
		if allowReplace && p.cur.Type == kwREPLACE {
			p.advance()
			if p.cur.Type == tokSCONST || p.isIdentLike() {
				clause.OldPass = p.cur.Str
				p.advance()
			}
		}
	} else if p.isIdentLikeStr("EXTERNALLY") {
		// IDENTIFIED EXTERNALLY [AS '...']
		p.advance()
		clause.Type = nodes.IDENTIFIED_EXTERNALLY
		if p.cur.Type == kwAS {
			p.advance()
			if p.cur.Type == tokSCONST {
				clause.ExternalAs = p.cur.Str
				p.advance()
			}
		}
	} else if p.isIdentLikeStr("GLOBALLY") {
		// IDENTIFIED GLOBALLY [AS '...']
		p.advance()
		clause.Type = nodes.IDENTIFIED_GLOBALLY
		if p.cur.Type == kwAS {
			p.advance()
			if p.cur.Type == tokSCONST {
				clause.ExternalAs = p.cur.Str
				p.advance()
			}
		}
	}

	clause.Loc.End = p.prev.End
	return clause, nil
}

// parseUserQuotaClause parses a QUOTA clause.
//
//	QUOTA { size_clause | UNLIMITED } ON tablespace
//
//	size_clause ::= integer [ K | M | G | T ]
func (p *Parser) parseUserQuotaClause() (*nodes.UserQuotaClause, error) {
	start := p.pos()
	// Already consumed QUOTA keyword
	clause := &nodes.UserQuotaClause{
		Loc: nodes.Loc{Start: start},
	}

	if p.isIdentLikeStr("UNLIMITED") {
		clause.Size = "UNLIMITED"
		p.advance()
	} else {
		// size_clause: integer [K|M|G|T]
		size := ""
		if p.cur.Type == tokICONST {
			size = p.cur.Str
			p.advance()
			// Optional size suffix
			if p.isIdentLike() {
				s := p.cur.Str
				if s == "K" || s == "M" || s == "G" || s == "T" {
					size += s
					p.advance()
				}
			}
		}
		clause.Size = size
	}

	// ON tablespace
	if p.cur.Type == kwON {
		p.advance()
		var parseErr900 error
		clause.Tablespace, parseErr900 = p.parseObjectName()
		if parseErr900 != nil {
			return nil, parseErr900
		}
	}

	clause.Loc.End = p.prev.End
	return clause, nil
}

// parseContainerClause parses CONTAINER = { ALL | CURRENT }.
// Returns: *bool (true=ALL, false=CURRENT), or nil if not present.
func (p *Parser) parseContainerClause() (*bool, error) {
	// Already consumed CONTAINER
	if p.cur.Type == '=' {
		p.advance()
	}
	if p.cur.Type == kwALL {
		p.advance()
		v := true
		return &v, nil
	} else if p.isIdentLikeStr("CURRENT") {
		p.advance()
		v := false
		return &v, nil
	}
	return nil, nil
}

// parseUserOptions parses common user option clauses for CREATE/ALTER USER.
// Returns when it encounters a token it doesn't recognize as a user option.
func (p *Parser) parseUserOptions(
	setTablespace func(string),
	setTempTablespace func(string, bool),
	addQuota func(*nodes.UserQuotaClause),
	setProfile func(string),
	setPasswordExpire func(),
	setAccountLock func(*bool),
	setEnableEditions func(),
	setCollation func(string),
	setContainer func(*bool),
) error {
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		switch {
		case p.cur.Type == kwDEFAULT:
			p.advance()
			if p.cur.Type == kwTABLESPACE {
				// DEFAULT TABLESPACE
				p.advance()
				if p.isIdentLike() {
					setTablespace(p.cur.Str)
					p.advance()
				}
			} else if p.isIdentLikeStr("COLLATION") {
				// DEFAULT COLLATION
				p.advance()
				if p.isIdentLike() {
					setCollation(p.cur.Str)
					p.advance()
				}
			} else if p.cur.Type == kwROLE {
				// DEFAULT ROLE — only for ALTER USER, handled by caller
				return nil
			} else {
				// Unknown DEFAULT clause, skip
				return nil
			}

		case p.isIdentLikeStr("LOCAL"):
			// LOCAL TEMPORARY TABLESPACE
			p.advance()
			if p.cur.Type == kwTEMPORARY {
				p.advance()
				if p.cur.Type == kwTABLESPACE {
					p.advance()
					if p.isIdentLike() {
						setTempTablespace(p.cur.Str, true)
						p.advance()
					}
				}
			}

		case p.cur.Type == kwTEMPORARY:
			// TEMPORARY TABLESPACE
			p.advance()
			if p.cur.Type == kwTABLESPACE {
				p.advance()
				if p.isIdentLike() {
					setTempTablespace(p.cur.Str, false)
					p.advance()
				}
			}

		case p.isIdentLikeStr("QUOTA"):
			// QUOTA clause
			p.advance()
			parseValue85, parseErr86 := p.parseUserQuotaClause()
			if parseErr86 != nil {
				return parseErr86
			}
			addQuota(parseValue85)

		case p.cur.Type == kwPROFILE:
			// PROFILE
			p.advance()
			if p.isIdentLike() {
				setProfile(p.cur.Str)
				p.advance()
			}

		case p.isIdentLikeStr("PASSWORD"):
			// PASSWORD EXPIRE
			p.advance()
			if p.isIdentLikeStr("EXPIRE") {
				p.advance()
				setPasswordExpire()
			} else {
				return nil
			}

		case p.isIdentLikeStr("ACCOUNT"):
			// ACCOUNT { LOCK | UNLOCK }
			p.advance()
			if p.cur.Type == kwLOCK {
				p.advance()
				v := true
				setAccountLock(&v)
			} else if p.isIdentLikeStr("UNLOCK") {
				p.advance()
				v := false
				setAccountLock(&v)
			}

		case p.cur.Type == kwENABLE:
			// ENABLE EDITIONS | ENABLE DICTIONARY PROTECTION
			p.advance()
			if p.isIdentLikeStr("EDITIONS") {
				p.advance()
				setEnableEditions()
				// [ FOR object_type [, object_type]... ] [ FORCE ]
				if p.cur.Type == kwFOR {
					p.advance()
					for p.isIdentLike() || p.cur.Type == ',' {
						p.advance()
					}
				}
				if p.cur.Type == kwFORCE {
					p.advance()
				}
			} else if p.isIdentLikeStr("DICTIONARY") {
				p.advance() // consume DICTIONARY
				if p.isIdentLikeStr("PROTECTION") {
					p.advance() // consume PROTECTION
				}
			} else {
				return nil
			}

		case p.cur.Type == kwDISABLE:
			// DISABLE DICTIONARY PROTECTION
			p.advance()
			if p.isIdentLikeStr("DICTIONARY") {
				p.advance() // consume DICTIONARY
				if p.isIdentLikeStr("PROTECTION") {
					p.advance() // consume PROTECTION
				}
			} else {
				return nil
			}

		case p.isIdentLikeStr("CONTAINER"):
			// CONTAINER = { ALL | CURRENT }
			p.advance()
			parseValue87, parseErr88 := p.parseContainerClause()
			if parseErr88 != nil {
				return parseErr88
			}
			setContainer(parseValue87)

		case p.cur.Type == kwREAD:
			// READ ONLY | READ WRITE
			p.advance()
			if p.isIdentLikeStr("ONLY") || p.cur.Type == kwWRITE {
				p.advance()
			}

		case p.isIdentLikeStr("HTTP"):
			// [HTTP] DIGEST { ENABLE | DISABLE }
			p.advance() // consume HTTP
			if p.isIdentLikeStr("DIGEST") {
				p.advance() // consume DIGEST
				if p.cur.Type == kwENABLE || p.cur.Type == kwDISABLE {
					p.advance()
				}
			}

		case p.isIdentLikeStr("DIGEST"):
			// DIGEST { ENABLE | DISABLE }
			p.advance() // consume DIGEST
			if p.cur.Type == kwENABLE || p.cur.Type == kwDISABLE {
				p.advance()
			}

		case p.isIdentLikeStr("EXPIRE"):
			// EXPIRE PASSWORD ROLLOVER PERIOD (ALTER USER)
			p.advance() // consume EXPIRE
			if p.isIdentLikeStr("PASSWORD") {
				p.advance() // consume PASSWORD
				if p.isIdentLikeStr("ROLLOVER") {
					p.advance() // consume ROLLOVER
					if p.isIdentLikeStr("PERIOD") {
						p.advance() // consume PERIOD
					}
				}
			}

		case p.cur.Type == kwSET || p.isIdentLikeStr("ADD") || p.isIdentLikeStr("REMOVE"):
			// container_data_clause: { SET | ADD | REMOVE } CONTAINER_DATA ...
			p.advance() // consume SET/ADD/REMOVE
			if p.isIdentLikeStr("CONTAINER_DATA") {
				p.advance() // consume CONTAINER_DATA
				// = or ( ... )
				if p.cur.Type == '=' {
					p.advance()
				}
				if p.cur.Type == '(' {
					p.advance()
					for p.cur.Type != ')' && p.cur.Type != tokEOF {
						p.advance()
					}
					if p.cur.Type == ')' {
						p.advance()
					}
				} else if p.cur.Type == kwALL || p.cur.Type == kwDEFAULT {
					p.advance()
				}
				// FOR [schema.]object_name
				if p.cur.Type == kwFOR {
					p.advance()
					parseDiscard902, parseErr901 := p.parseObjectName()
					_ = parseDiscard902
					if parseErr901 != nil {
						return parseErr901
					}
				}
			} else {
				return nil
			}

		default:
			return nil
		}
	}
	return nil
}

// parseCreateUserStmt parses a CREATE USER statement.
//
// BNF: oracle/parser/bnf/CREATE-USER.bnf
//
//	CREATE USER [ IF NOT EXISTS ] user
//	    IDENTIFIED { BY password
//	               | EXTERNALLY [ AS 'certificate_DN' ]
//	               | GLOBALLY [ AS 'directory_DN' ]
//	               | NO AUTHENTICATION
//	               }
//	    [ AND FACTOR 'auth_method' AS 'external_name' ]
//	    [ [ HTTP ] DIGEST { ENABLE | DISABLE } ]
//	    [ DEFAULT COLLATION collation_name ]
//	    [ DEFAULT TABLESPACE tablespace ]
//	    [ [ LOCAL ] TEMPORARY TABLESPACE { tablespace | tablespace_group_name } ]
//	    [ { QUOTA { size_clause | UNLIMITED } ON tablespace } ]...
//	    [ PROFILE profile_name ]
//	    [ PASSWORD EXPIRE ]
//	    [ ACCOUNT { LOCK | UNLOCK } ]
//	    [ ENABLE EDITIONS ]
//	    [ CONTAINER = { ALL | CURRENT } ]
//	    [ READ ONLY | READ WRITE ] ;
//
//	size_clause:
//	    integer { K | M | G | T }
func (p *Parser) parseCreateUserStmt(start int) (nodes.StmtNode, error) {
	stmt := &nodes.CreateUserStmt{
		Loc: nodes.Loc{Start: start},
	}

	// IF NOT EXISTS
	if p.cur.Type == kwIF {
		p.advance()
		if p.cur.Type == kwNOT {
			p.advance()
			if p.cur.Type == kwEXISTS {
				p.advance()
				stmt.IfNotExists = true
			}
		}
	}
	var parseErr903 error

	stmt.Name, parseErr903 = p.parseObjectName()
	if parseErr903 !=

		// IDENTIFIED clause or NO AUTHENTICATION
		nil {
		return nil, parseErr903
	}

	if p.cur.Type == kwIDENTIFIED {
		p.advance()
		var parseErr904 error
		stmt.Identified, parseErr904 = p.parseIdentifiedClause(false)
		if parseErr904 != nil {
			return nil, parseErr904
		}
	} else if p.isIdentLikeStr("NO") {
		p.advance()
		if p.isIdentLikeStr("AUTHENTICATION") {
			p.advance()
			stmt.Identified = &nodes.IdentifiedClause{
				Type: nodes.IDENTIFIED_NO_AUTH,
				Loc:  nodes.Loc{Start: p.pos(), End: p.prev.End},
			}
		}
	}
	parseErr905 :=

		// Parse remaining options
		p.parseUserOptions(
			func(ts string) { stmt.DefaultTablespace = ts },
			func(ts string, local bool) { stmt.TempTablespace = ts; stmt.LocalTemp = local },
			func(q *nodes.UserQuotaClause) { stmt.Quotas = append(stmt.Quotas, q) },
			func(prof string) { stmt.Profile = prof },
			func() { stmt.PasswordExpire = true },
			func(v *bool) { stmt.AccountLock = v },
			func() { stmt.EnableEditions = true },
			func(c string) { stmt.DefaultCollation = c },
			func(v *bool) { stmt.ContainerAll = v },
		)
	if parseErr905 != nil {
		return nil, parseErr905
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterUserStmt parses an ALTER USER statement.
//
// BNF: oracle/parser/bnf/ALTER-USER.bnf
//
//	ALTER USER user
//	    [ IF EXISTS ]
//	    { alter_identified_clause
//	    | proxy_clause
//	    | alter_user_clauses
//	    } ;
//
//	alter_identified_clause:
//	    IDENTIFIED
//	    { BY password [ REPLACE old_password ]
//	    | EXTERNALLY [ AS 'certificate_DN' | AS 'kerberos_principal_name' ]
//	    | GLOBALLY [ AS '[ directory_DN ]' ]
//	    }
//	  | NO AUTHENTICATION
//
//	alter_user_clauses:
//	    [ DEFAULT TABLESPACE tablespace ]
//	    [ TEMPORARY TABLESPACE { tablespace | tablespace_group_name } ]
//	    [ LOCAL TEMPORARY TABLESPACE { tablespace | tablespace_group_name } ]
//	    [ quota_clause [, quota_clause ]... ]
//	    [ PROFILE profile ]
//	    [ DEFAULT ROLE { role [, role ]... | ALL [ EXCEPT role [, role ]... ] | NONE } ]
//	    [ PASSWORD EXPIRE ]
//	    [ EXPIRE PASSWORD ROLLOVER PERIOD ]
//	    [ ACCOUNT { LOCK | UNLOCK } ]
//	    [ ENABLE EDITIONS [ FOR object_type [, object_type ]... ] [ FORCE ] ]
//	    [ [ HTTP ] DIGEST { ENABLE | DISABLE } ]
//	    [ CONTAINER = { ALL | CURRENT } ]
//	    [ { ENABLE | DISABLE } DICTIONARY PROTECTION ]
//	    [ DEFAULT COLLATION collation_name ]
//	    [ container_data_clause ]
//
//	quota_clause:
//	    QUOTA { size_clause | UNLIMITED } ON tablespace
//
//	proxy_clause:
//	    { GRANT CONNECT THROUGH { ENTERPRISE USERS | db_user_proxy [ proxy_clause_options ] }
//	    | REVOKE CONNECT THROUGH { ENTERPRISE USERS | db_user_proxy }
//	    }
//
//	proxy_clause_options:
//	    [ WITH { ROLE { role [, role ]... | ALL EXCEPT role [, role ]... } | NO ROLES } ]
//	    [ AUTHENTICATION REQUIRED ]
//
//	container_data_clause:
//	    { SET | ADD | REMOVE } CONTAINER_DATA
//	    { ( container_name [, container_name ]... ) | ALL | DEFAULT }
//	    [ FOR [ schema. ] object_name ]
func (p *Parser) parseAlterUserStmt(start int) (nodes.StmtNode, error) {
	stmt := &nodes.AlterUserStmt{
		Loc: nodes.Loc{Start: start},
	}

	// IF EXISTS
	if p.cur.Type == kwIF {
		p.advance()
		if p.cur.Type == kwEXISTS {
			p.advance()
			stmt.IfExists = true
		}
	}
	var parseErr906 error

	stmt.Name, parseErr906 = p.parseObjectName()
	if parseErr906 !=

		// IDENTIFIED clause or NO AUTHENTICATION
		nil {
		return nil, parseErr906
	}

	if p.cur.Type == kwIDENTIFIED {
		p.advance()
		var parseErr907 error
		stmt.Identified, parseErr907 = p.parseIdentifiedClause(true)
		if parseErr907 != nil {
			return nil, parseErr907
		}
	} else if p.isIdentLikeStr("NO") {
		p.advance()
		if p.isIdentLikeStr("AUTHENTICATION") {
			p.advance()
			stmt.Identified = &nodes.IdentifiedClause{
				Type: nodes.IDENTIFIED_NO_AUTH,
				Loc:  nodes.Loc{Start: p.pos(), End: p.prev.End},
			}
		}
	}
	parseErr908 :=

		// Parse remaining options (may return early on DEFAULT ROLE)
		p.parseUserOptions(
			func(ts string) { stmt.DefaultTablespace = ts },
			func(ts string, local bool) { stmt.TempTablespace = ts; stmt.LocalTemp = local },
			func(q *nodes.UserQuotaClause) { stmt.Quotas = append(stmt.Quotas, q) },
			func(prof string) { stmt.Profile = prof },
			func() { stmt.PasswordExpire = true },
			func(v *bool) { stmt.AccountLock = v },
			func() { stmt.EnableEditions = true },
			func(c string) { stmt.DefaultCollation = c },
			func(v *bool) { stmt.ContainerAll = v },
		)
	if parseErr908 !=

		// DEFAULT ROLE clause (parseUserOptions returns when it sees DEFAULT ROLE)
		nil {
		return nil, parseErr908
	}

	if p.cur.Type == kwROLE {
		p.advance()
		var parseErr909 error
		stmt.DefaultRole, parseErr909 = p.parseDefaultRoleClause()
		if parseErr909 !=

			// Continue parsing remaining options after DEFAULT ROLE
			nil {
			return nil, parseErr909
		}
		parseErr910 := p.parseUserOptions(
			func(ts string) { stmt.DefaultTablespace = ts },
			func(ts string, local bool) { stmt.TempTablespace = ts; stmt.LocalTemp = local },
			func(q *nodes.UserQuotaClause) { stmt.Quotas = append(stmt.Quotas, q) },
			func(prof string) { stmt.Profile = prof },
			func() { stmt.PasswordExpire = true },
			func(v *bool) { stmt.AccountLock = v },
			func() { stmt.EnableEditions = true },
			func(c string) { stmt.DefaultCollation = c },
			func(v *bool) { stmt.ContainerAll = v },
		)
		if parseErr910 != nil {
			return nil, parseErr910
		}
	}
	if p.cur.Type == kwGRANT || p.cur.Type == kwREVOKE {
		parseErr911 := p.parseAlterUserProxyClause()
		if parseErr911 != nil {
			return nil, parseErr911
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterUserProxyClause consumes ALTER USER proxy clauses:
// GRANT/REVOKE CONNECT THROUGH { ENTERPRISE USERS | proxy_user ... }.
func (p *Parser) parseAlterUserProxyClause() error {
	p.advance() // GRANT or REVOKE
	if p.cur.Type == kwCONNECT {
		p.advance()
	}
	if p.isIdentLikeStr("THROUGH") {
		p.advance()
	}
	if p.isIdentLikeStr("ENTERPRISE") {
		p.advance()
		if p.cur.Type == kwUSER || p.isIdentLikeStr("USERS") {
			p.advance()
		}
	} else if p.isIdentLike() {
		parseDiscard913, parseErr912 := p.parseObjectName()
		_ = parseDiscard913
		if parseErr912 != nil {
			return parseErr912

			// parseDefaultRoleClause parses the DEFAULT ROLE clause for ALTER USER.
			//
			//	DEFAULT ROLE { role [, role ]... | ALL [ EXCEPT role [, role ]... ] | NONE }
		}
	}
	for !p.isStatementEnd() {
		p.advance()
	}
	return nil
}

func (p *Parser) parseDefaultRoleClause() (*nodes.DefaultRoleClause, error) {
	start := p.pos()
	clause := &nodes.DefaultRoleClause{
		Loc: nodes.Loc{Start: start},
	}

	if p.cur.Type == kwALL {
		p.advance()
		clause.AllRoles = true
		// ALL EXCEPT role [, role]...
		if p.cur.Type == kwEXCEPT {
			p.advance()
			clause.ExceptAll = true
			for {
				parseValue89, parseErr90 := p.parseObjectName()
				if parseErr90 != nil {
					return nil, parseErr90
				}
				clause.Roles = append(clause.Roles, parseValue89)
				if p.cur.Type != ',' {
					break
				}
				p.advance()
			}
		}
	} else if p.isIdentLikeStr("NONE") {
		p.advance()
		clause.NoneRole = true
	} else {
		// Specific role list
		for {
			parseValue91, parseErr92 := p.parseObjectName()
			if parseErr92 != nil {
				return nil, parseErr92
			}
			clause.Roles = append(clause.Roles, parseValue91)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}

	clause.Loc.End = p.prev.End
	return clause, nil
}

// parseCreateRoleStmt parses a CREATE ROLE statement.
//
// BNF: oracle/parser/bnf/CREATE-ROLE.bnf
//
//	CREATE ROLE role
//	    [ NOT IDENTIFIED
//	    | IDENTIFIED { BY password
//	                 | USING [ schema. ] package
//	                 | EXTERNALLY
//	                 | GLOBALLY [ AS 'directory_name' ]
//	                 }
//	    ]
//	    [ CONTAINER = { ALL | CURRENT } ] ;
func (p *Parser) parseCreateRoleStmt(start int) (nodes.StmtNode, error) {
	stmt := &nodes.CreateRoleStmt{
		Loc: nodes.Loc{Start: start},
	}
	var parseErr914 error

	stmt.Name, parseErr914 = p.parseObjectName()
	if parseErr914 !=

		// Optional IDENTIFIED clause
		nil {
		return nil, parseErr914
	}

	if p.cur.Type == kwIDENTIFIED {
		p.advance()
		stmt.HasIdentified = true
		parseErr915 := p.parseRoleIdentifiedClause(
			func(t nodes.RoleIdentifiedType) { stmt.IdentifiedType = t },
			func(v string) { stmt.IdentifyBy = v },
			func(v string) { stmt.IdentifySchema = v },
		)
		if parseErr915 != nil {
			return nil, parseErr915
		}
	} else if p.cur.Type == kwNOT {
		p.advance()
		if p.cur.Type == kwIDENTIFIED {
			p.advance()
		}
		stmt.HasIdentified = true
		stmt.IdentifiedType = nodes.ROLE_NOT_IDENTIFIED
	}

	// CONTAINER = { ALL | CURRENT }
	if p.isIdentLikeStr("CONTAINER") {
		p.advance()
		var parseErr916 error
		stmt.ContainerAll, parseErr916 = p.parseContainerClause()
		if parseErr916 != nil {
			return nil, parseErr916
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseRoleIdentifiedClause parses the IDENTIFIED clause for roles.
// Called after IDENTIFIED keyword has been consumed.
//
//	IDENTIFIED { BY password
//	           | USING [ schema. ] package
//	           | EXTERNALLY
//	           | GLOBALLY [ AS 'directory_name' ]
//	           }
func (p *Parser) parseRoleIdentifiedClause(
	setType func(nodes.RoleIdentifiedType),
	setValue func(string),
	setSchema func(string),
) error {
	switch {
	case p.cur.Type == kwBY:
		// IDENTIFIED BY password
		p.advance()
		setType(nodes.ROLE_IDENTIFIED_BY)
		if p.isIdentLike() || p.cur.Type == tokSCONST {
			setValue(p.cur.Str)
			p.advance()
		}
	case p.cur.Type == kwUSING:
		// IDENTIFIED USING [schema.]package
		p.advance()
		setType(nodes.ROLE_IDENTIFIED_USING)
		name, parseErr917 := p.parseObjectName()
		if parseErr917 != nil {
			return parseErr917
		}
		if name != nil {
			if name.Schema != "" {
				setSchema(name.Schema)
			}
			setValue(name.Name)
		}
	case p.isIdentLikeStr("EXTERNALLY"):
		// IDENTIFIED EXTERNALLY
		p.advance()
		setType(nodes.ROLE_IDENTIFIED_EXTERNALLY)
	case p.isIdentLikeStr("GLOBALLY"):
		// IDENTIFIED GLOBALLY [AS 'directory_name']
		p.advance()
		setType(nodes.ROLE_IDENTIFIED_GLOBALLY)
		if p.cur.Type == kwAS {
			p.advance()
			if p.cur.Type == tokSCONST {
				setValue(p.cur.Str)
				p.advance()
			}
		}
	}
	return nil
}

// parseAlterRoleStmt parses an ALTER ROLE statement.
//
// BNF: oracle/parser/bnf/ALTER-ROLE.bnf
//
//	ALTER ROLE role
//	    { NOT IDENTIFIED
//	    | IDENTIFIED BY password
//	    | IDENTIFIED EXTERNALLY
//	    | IDENTIFIED GLOBALLY AS 'domain_name'
//	    | IDENTIFIED USING [ schema. ] package_name
//	    }
//	    [ CONTAINER = { ALL | CURRENT } ]
func (p *Parser) parseAlterRoleStmt(start int) (nodes.StmtNode, error) {
	stmt := &nodes.AlterRoleStmt{
		Loc: nodes.Loc{Start: start},
	}
	var parseErr918 error

	stmt.Name, parseErr918 = p.parseObjectName()
	if parseErr918 !=

		// Required IDENTIFIED clause or NOT IDENTIFIED
		nil {
		return nil, parseErr918
	}

	if p.cur.Type == kwIDENTIFIED {
		p.advance()
		parseErr919 := p.parseRoleIdentifiedClause(
			func(t nodes.RoleIdentifiedType) { stmt.IdentifiedType = t },
			func(v string) { stmt.IdentifyBy = v },
			func(v string) { stmt.IdentifySchema = v },
		)
		if parseErr919 != nil {
			return nil, parseErr919
		}
	} else if p.cur.Type == kwNOT {
		p.advance()
		if p.cur.Type == kwIDENTIFIED {
			p.advance()
		}
		stmt.IdentifiedType = nodes.ROLE_NOT_IDENTIFIED
	}

	// CONTAINER = { ALL | CURRENT }
	if p.isIdentLikeStr("CONTAINER") {
		p.advance()
		var parseErr920 error
		stmt.ContainerAll, parseErr920 = p.parseContainerClause()
		if parseErr920 != nil {
			return nil, parseErr920
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterProfileStmt parses an ALTER PROFILE statement.
//
// BNF: oracle/parser/bnf/ALTER-PROFILE.bnf
//
//	ALTER PROFILE profile_name
//	    LIMIT { resource_parameters | password_parameters }
//	    [ CONTAINER = { ALL | CURRENT } ] ;
//
//	resource_parameters:
//	    { SESSIONS_PER_USER { integer | UNLIMITED | DEFAULT }
//	    | CPU_PER_SESSION { integer | UNLIMITED | DEFAULT }
//	    | CPU_PER_CALL { integer | UNLIMITED | DEFAULT }
//	    | CONNECT_TIME { integer | UNLIMITED | DEFAULT }
//	    | IDLE_TIME { integer | UNLIMITED | DEFAULT }
//	    | LOGICAL_READS_PER_SESSION { integer | UNLIMITED | DEFAULT }
//	    | LOGICAL_READS_PER_CALL { integer | UNLIMITED | DEFAULT }
//	    | PRIVATE_SGA { size_clause | UNLIMITED | DEFAULT }
//	    | COMPOSITE_LIMIT { integer | UNLIMITED | DEFAULT }
//	    }
//
//	password_parameters:
//	    { PASSWORD_LIFE_TIME { integer | UNLIMITED | DEFAULT }
//	    | PASSWORD_GRACE_TIME { integer | UNLIMITED | DEFAULT }
//	    | PASSWORD_REUSE_TIME { integer | UNLIMITED | DEFAULT }
//	    | PASSWORD_REUSE_MAX { integer | UNLIMITED | DEFAULT }
//	    | PASSWORD_LOCK_TIME { integer | UNLIMITED | DEFAULT }
//	    | FAILED_LOGIN_ATTEMPTS { integer | UNLIMITED | DEFAULT }
//	    | INACTIVE_ACCOUNT_TIME { integer | UNLIMITED | DEFAULT }
//	    | PASSWORD_ROLLOVER_TIME { integer | UNLIMITED | DEFAULT }
//	    | PASSWORD_VERIFY_FUNCTION { function_name | NULL | DEFAULT }
//	    }
func (p *Parser) parseAlterProfileStmt(start int) (nodes.StmtNode, error) {
	stmt := &nodes.AlterProfileStmt{
		Loc: nodes.Loc{Start: start},
	}
	var parseErr921 error

	stmt.Name, parseErr921 = p.parseObjectName()
	if parseErr921 !=

		// Parse LIMIT clauses (same logic as CREATE PROFILE)
		nil {
		return nil, parseErr921
	}

	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		if p.cur.Type == kwLIMIT {
			p.advance()
			for p.isProfileParam() {
				lim, parseErr922 := p.parseProfileLimit()
				if parseErr922 != nil {
					return nil, parseErr922
				}
				stmt.Limits = append(stmt.Limits, lim)
			}
		} else if p.isProfileParam() {
			lim, parseErr923 := p.parseProfileLimit()
			if parseErr923 != nil {
				return nil, parseErr923
			}
			stmt.Limits = append(stmt.Limits, lim)
		} else if p.isIdentLikeStr("CONTAINER") {
			p.advance()
			var parseErr924 error
			stmt.ContainerAll, parseErr924 = p.parseContainerClause()
			if parseErr924 != nil {
				return nil, parseErr924
			}
		} else {
			break
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterResourceCostStmt parses an ALTER RESOURCE COST statement.
//
// BNF: oracle/parser/bnf/ALTER-RESOURCE-COST.bnf
//
//	ALTER RESOURCE COST
//	    { CPU_PER_SESSION integer
//	    | CONNECT_TIME integer
//	    | LOGICAL_READS_PER_SESSION integer
//	    | PRIVATE_SGA integer
//	    } [ { CPU_PER_SESSION integer
//	        | CONNECT_TIME integer
//	        | LOGICAL_READS_PER_SESSION integer
//	        | PRIVATE_SGA integer
//	        } ]...
func (p *Parser) parseAlterResourceCostStmt(start int) (nodes.StmtNode, error) {
	stmt := &nodes.AlterResourceCostStmt{
		Loc: nodes.Loc{Start: start},
	}

	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		if !p.isIdentLike() {
			break
		}
		name := p.cur.Str
		switch name {
		case "CPU_PER_SESSION", "CONNECT_TIME", "LOGICAL_READS_PER_SESSION", "PRIVATE_SGA":
			entryStart := p.pos()
			p.advance()
			entry := &nodes.ResourceCostEntry{
				Name: name,
				Loc:  nodes.Loc{Start: entryStart},
			}
			if p.cur.Type == tokICONST {
				entry.Value = p.cur.Str
				p.advance()
			}
			entry.Loc.End = p.prev.End
			stmt.Costs = append(stmt.Costs, entry)
		default:
			goto done
		}
	}
done:

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// isProfileParam returns true if the current identifier is a known profile resource or password parameter.
func (p *Parser) isProfileParam() bool {
	if !p.isIdentLike() {
		return false
	}
	switch p.cur.Str {
	case "SESSIONS_PER_USER", "CPU_PER_SESSION", "CPU_PER_CALL",
		"CONNECT_TIME", "IDLE_TIME",
		"LOGICAL_READS_PER_SESSION", "LOGICAL_READS_PER_CALL",
		"PRIVATE_SGA", "COMPOSITE_LIMIT",
		"FAILED_LOGIN_ATTEMPTS",
		"PASSWORD_LIFE_TIME", "PASSWORD_REUSE_TIME", "PASSWORD_REUSE_MAX",
		"PASSWORD_LOCK_TIME", "PASSWORD_GRACE_TIME",
		"PASSWORD_VERIFY_FUNCTION", "PASSWORD_ROLLOVER_TIME",
		"INACTIVE_ACCOUNT_TIME":
		return true
	}
	return false
}

// parseCreateProfileStmt parses a CREATE PROFILE statement.
//
// BNF: oracle/parser/bnf/CREATE-PROFILE.bnf
//
//	CREATE [ MANDATORY ] PROFILE profile
//	    LIMIT { resource_parameters | password_parameters }
//	          [ { resource_parameters | password_parameters } ]...
//	    [ CONTAINER = { ALL | CURRENT } ] ;
//
//	resource_parameters:
//	    { SESSIONS_PER_USER { integer | UNLIMITED | DEFAULT }
//	    | CPU_PER_SESSION { integer | UNLIMITED | DEFAULT }
//	    | CPU_PER_CALL { integer | UNLIMITED | DEFAULT }
//	    | CONNECT_TIME { integer | UNLIMITED | DEFAULT }
//	    | IDLE_TIME { integer | UNLIMITED | DEFAULT }
//	    | LOGICAL_READS_PER_SESSION { integer | UNLIMITED | DEFAULT }
//	    | LOGICAL_READS_PER_CALL { integer | UNLIMITED | DEFAULT }
//	    | PRIVATE_SGA { size_clause | UNLIMITED | DEFAULT }
//	    | COMPOSITE_LIMIT { integer | UNLIMITED | DEFAULT }
//	    }
//
//	password_parameters:
//	    { FAILED_LOGIN_ATTEMPTS { integer | UNLIMITED | DEFAULT }
//	    | PASSWORD_LIFE_TIME { expr | UNLIMITED | DEFAULT }
//	    | PASSWORD_REUSE_TIME { expr | UNLIMITED | DEFAULT }
//	    | PASSWORD_REUSE_MAX { integer | UNLIMITED | DEFAULT }
//	    | PASSWORD_LOCK_TIME { expr | UNLIMITED | DEFAULT }
//	    | PASSWORD_GRACE_TIME { expr | UNLIMITED | DEFAULT }
//	    | INACTIVE_ACCOUNT_TIME { integer | UNLIMITED | DEFAULT }
//	    | PASSWORD_VERIFY_FUNCTION { function_name | NULL | DEFAULT }
//	    | PASSWORD_ROLLOVER_TIME { expr | UNLIMITED | DEFAULT }
//	    }
//
//	size_clause:
//	    integer [ K | M | G | T | P | E ]
func (p *Parser) parseCreateProfileStmt(start int, mandatory bool) (nodes.StmtNode, error) {
	stmt := &nodes.CreateProfileStmt{
		Loc:       nodes.Loc{Start: start},
		Mandatory: mandatory,
	}
	var parseErr925 error

	stmt.Name, parseErr925 = p.parseObjectName()
	if parseErr925 !=

		// Parse LIMIT clauses
		nil {
		return nil, parseErr925
	}

	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		if p.cur.Type == kwLIMIT {
			p.advance()
			// Parse parameters after LIMIT
			for p.isProfileParam() {
				lim, parseErr926 := p.parseProfileLimit()
				if parseErr926 != nil {
					return nil, parseErr926
				}
				stmt.Limits = append(stmt.Limits, lim)
			}
		} else if p.isProfileParam() {
			// Parameters can appear without repeated LIMIT keyword
			lim, parseErr927 := p.parseProfileLimit()
			if parseErr927 != nil {
				return nil, parseErr927
			}
			stmt.Limits = append(stmt.Limits, lim)
		} else if p.isIdentLikeStr("CONTAINER") {
			p.advance()
			var parseErr928 error
			stmt.ContainerAll, parseErr928 = p.parseContainerClause()
			if parseErr928 != nil {
				return nil, parseErr928
			}
		} else {
			break
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseProfileLimit parses a single profile limit parameter and its value.
//
//	param_name { integer | size_clause | UNLIMITED | DEFAULT | NULL | function_name }
func (p *Parser) parseProfileLimit() (*nodes.ProfileLimit, error) {
	start := p.pos()
	lim := &nodes.ProfileLimit{
		Loc: nodes.Loc{Start: start},
	}

	lim.Name = p.cur.Str
	p.advance()

	// Parse value
	switch {
	case p.isIdentLikeStr("UNLIMITED"):
		lim.Value = "UNLIMITED"
		p.advance()
	case p.cur.Type == kwDEFAULT:
		lim.Value = "DEFAULT"
		p.advance()
	case p.cur.Type == kwNULL:
		lim.Value = "NULL"
		p.advance()
	case p.cur.Type == tokICONST:
		val := p.cur.Str
		p.advance()
		// Optional size suffix (K, M, G)
		if p.isIdentLike() {
			s := p.cur.Str
			if s == "K" || s == "M" || s == "G" || s == "T" {
				val += s
				p.advance()
			}
		}
		lim.Value = val
	case p.isIdentLike():
		// function_name for PASSWORD_VERIFY_FUNCTION
		lim.Value = p.cur.Str
		p.advance()
	}

	lim.Loc.End = p.prev.End
	return lim, nil
}

// parseAdministerKeyManagementStmt parses an ADMINISTER KEY MANAGEMENT statement.
// The current token is "ADMINISTER" (identifier).
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/ADMINISTER-KEY-MANAGEMENT.html
//
//	ADMINISTER KEY MANAGEMENT
//	  { keystore_management_clauses
//	  | key_management_clauses
//	  | secret_management_clauses
//	  | zero_downtime_software_patching_clauses }
//
//	keystore_management_clauses:
//	    CREATE KEYSTORE 'keystore_location' IDENTIFIED BY { password | EXTERNAL STORE }
//	  | CREATE [ LOCAL ] AUTO_LOGIN KEYSTORE FROM KEYSTORE 'keystore_location'
//	        IDENTIFIED BY { password | EXTERNAL STORE }
//	  | ALTER KEYSTORE PASSWORD [ FORCE KEYSTORE ]
//	        IDENTIFIED BY old_password SET new_password [ WITH BACKUP [ USING 'description' ] ]
//	  | CLOSE KEYSTORE [ IDENTIFIED BY { password | EXTERNAL STORE } ]
//	  | BACKUP KEYSTORE [ USING 'description' ] [ FORCE KEYSTORE ]
//	        IDENTIFIED BY { password | EXTERNAL STORE } [ TO 'keystore_location' ]
//	  | MERGE KEYSTORE 'keystore_location1' [ AND 'keystore_location2' ]
//	        IDENTIFIED BY { password | EXTERNAL STORE }
//	        INTO [ NEW ] KEYSTORE 'keystore_location3'
//	        IDENTIFIED BY { password | EXTERNAL STORE }
//	  | FORCE KEYSTORE { ISOLATE KEYSTORE | UNITE KEYSTORE }
//	        IDENTIFIED BY { EXTERNAL STORE | isolated_password }
//	  | SET KEYSTORE OPEN [ FORCE KEYSTORE ]
//	        IDENTIFIED BY { password | EXTERNAL STORE }
//	        [ CONTAINER = { CURRENT | ALL } ]
//	  | SET KEY [ USING TAG 'tag' ]
//	        [ USING ALGORITHM 'algorithm' ]
//	        [ FORCE KEYSTORE ]
//	        IDENTIFIED BY { password | EXTERNAL STORE }
//	        [ WITH BACKUP [ USING 'description' ] ]
//	        [ CONTAINER = { CURRENT | ALL } ]
//	  | CREATE KEY [ USING TAG 'tag' ]
//	        [ USING ALGORITHM 'algorithm' ]
//	        [ FORCE KEYSTORE ]
//	        IDENTIFIED BY { password | EXTERNAL STORE }
//	        [ WITH BACKUP [ USING 'description' ] ]
//	        [ CONTAINER = { CURRENT | ALL } ]
//	  | USE KEY 'key_id'
//	        [ USING TAG 'tag' ]
//	        [ FORCE KEYSTORE ]
//	        IDENTIFIED BY { password | EXTERNAL STORE }
//	        [ WITH BACKUP [ USING 'description' ] ]
//	  | SET TAG 'tag' FOR 'key_id'
//	        [ FORCE KEYSTORE ]
//	        IDENTIFIED BY { password | EXTERNAL STORE }
//	        [ WITH BACKUP [ USING 'description' ] ]
//	  | EXPORT [ ENCRYPTION ] KEYS WITH SECRET 'secret'
//	        TO 'filename'
//	        [ FORCE KEYSTORE ]
//	        IDENTIFIED BY { password | EXTERNAL STORE }
//	        [ WITH IDENTIFIER IN ( key_id [, ...] | subquery )  ]
//	  | IMPORT [ ENCRYPTION ] KEYS WITH SECRET 'secret'
//	        FROM 'filename'
//	        [ FORCE KEYSTORE ]
//	        IDENTIFIED BY { password | EXTERNAL STORE }
//	        [ WITH BACKUP [ USING 'description' ] ]
//	  | MOVE [ ENCRYPTION ] KEYS
//	        TO NEW KEYSTORE 'keystore_location'
//	        IDENTIFIED BY keystore_password
//	        FROM [ FORCE ] KEYSTORE
//	        IDENTIFIED BY { password | EXTERNAL STORE }
//	        [ WITH IDENTIFIER IN ( key_id [, ...] | subquery ) ]
//	        [ WITH BACKUP [ USING 'description' ] ]
//
//	secret_management_clauses:
//	    ADD SECRET 'secret' FOR CLIENT 'client_identifier'
//	        [ USING TAG 'tag' ]
//	        [ FORCE KEYSTORE ]
//	        IDENTIFIED BY { password | EXTERNAL STORE }
//	        [ WITH BACKUP [ USING 'description' ] ]
//	  | UPDATE SECRET 'secret' FOR CLIENT 'client_identifier'
//	        [ USING TAG 'tag' ]
//	        [ FORCE KEYSTORE ]
//	        IDENTIFIED BY { password | EXTERNAL STORE }
//	        [ WITH BACKUP [ USING 'description' ] ]
//	  | DELETE SECRET FOR CLIENT 'client_identifier'
//	        [ FORCE KEYSTORE ]
//	        IDENTIFIED BY { password | EXTERNAL STORE }
//	        [ WITH BACKUP [ USING 'description' ] ]
func (p *Parser) parseAdministerKeyManagementStmt() (nodes.StmtNode, error) {
	start := p.pos()
	p.advance() // consume ADMINISTER

	stmt := &nodes.AdminDDLStmt{
		Action:     "ADMINISTER",
		ObjectType: nodes.OBJECT_KEY_MANAGEMENT,
		Loc:        nodes.Loc{Start: start},
	}

	opts := &nodes.List{}
	if p.cur.Type != kwKEY {
		return nil, p.syntaxErrorAtCur()
	}
	scopeStart := p.pos()
	p.advance() // KEY
	if !p.isIdentLikeStr("MANAGEMENT") {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // MANAGEMENT
	opts.Items = append(opts.Items, &nodes.DDLOption{
		Key:   "SCOPE",
		Value: "KEY MANAGEMENT",
		Loc:   nodes.Loc{Start: scopeStart, End: p.prev.End},
	})

	if p.cur.Type != ';' && p.cur.Type != tokEOF {
		cmdStart := p.pos()
		command := p.consumeAdministerKeyManagementCommand()
		if command == "" {
			return nil, p.syntaxErrorAtCur()
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key:   "COMMAND",
			Value: command,
			Loc:   nodes.Loc{Start: cmdStart, End: p.prev.End},
		})
	}

	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		optStart := p.pos()
		switch {
		case p.cur.Type == kwIDENTIFIED:
			p.advance()
			if p.cur.Type != kwBY {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
			value := p.consumeAdministerKeyManagementValue()
			if value == "" {
				return nil, p.syntaxErrorAtCur()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key:   "IDENTIFIED BY",
				Value: value,
				Loc:   nodes.Loc{Start: optStart, End: p.prev.End},
			})
		case p.cur.Type == kwWITH:
			p.advance()
			if p.isIdentLikeStr("SECRET") {
				p.advance()
				value := p.consumeAdministerKeyManagementValue()
				if value == "" {
					return nil, p.syntaxErrorAtCur()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{Key: "WITH SECRET", Value: value, Loc: nodes.Loc{Start: optStart, End: p.prev.End}})
				continue
			}
			if p.isIdentLikeStr("BACKUP") {
				p.advance()
				opts.Items = append(opts.Items, &nodes.DDLOption{Key: "WITH BACKUP", Loc: nodes.Loc{Start: optStart, End: p.prev.End}})
				continue
			}
			return nil, p.syntaxErrorAtCur()
		case p.cur.Type == kwUSING:
			p.advance()
			key := "USING"
			if p.isIdentLikeStr("TAG") {
				p.advance()
				key = "USING TAG"
			}
			value := p.consumeAdministerKeyManagementValue()
			if value == "" {
				return nil, p.syntaxErrorAtCur()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: key, Value: value, Loc: nodes.Loc{Start: optStart, End: p.prev.End}})
		case p.cur.Type == kwFOR:
			p.advance()
			key := "FOR"
			if p.isIdentLikeStr("CLIENT") {
				p.advance()
				key = "FOR CLIENT"
			}
			value := p.consumeAdministerKeyManagementValue()
			if value == "" {
				return nil, p.syntaxErrorAtCur()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: key, Value: value, Loc: nodes.Loc{Start: optStart, End: p.prev.End}})
		case p.cur.Type == kwTO || p.cur.Type == kwFROM || p.cur.Type == kwINTO:
			key := p.cur.Str
			p.advance()
			value := p.consumeAdministerKeyManagementValue()
			if value == "" {
				return nil, p.syntaxErrorAtCur()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: key, Value: value, Loc: nodes.Loc{Start: optStart, End: p.prev.End}})
		case p.isIdentLikeStr("CONTAINER"):
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			value := p.consumeAdministerKeyManagementValue()
			if value == "" {
				return nil, p.syntaxErrorAtCur()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "CONTAINER", Value: value, Loc: nodes.Loc{Start: optStart, End: p.prev.End}})
		case p.cur.Type == kwFORCE:
			p.advance()
			if p.isIdentLikeStr("KEYSTORE") {
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "FORCE KEYSTORE", Loc: nodes.Loc{Start: optStart, End: p.prev.End}})
		default:
			value := p.consumeAdministerKeyManagementValue()
			if value == "" {
				return nil, p.syntaxErrorAtCur()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "ARG", Value: value, Loc: nodes.Loc{Start: optStart, End: p.prev.End}})
		}
	}

	stmt.Options = opts
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

func (p *Parser) consumeAdministerKeyManagementCommand() string {
	switch {
	case p.cur.Type == kwCREATE:
		p.advance()
		switch {
		case p.isIdentLikeStr("LOCAL"):
			p.advance()
			if p.isIdentLikeStr("AUTO_LOGIN") {
				p.advance()
			}
			if p.isIdentLikeStr("KEYSTORE") {
				p.advance()
			}
			return "CREATE LOCAL AUTO_LOGIN KEYSTORE"
		case p.isIdentLikeStr("AUTO_LOGIN"):
			p.advance()
			if p.isIdentLikeStr("KEYSTORE") {
				p.advance()
			}
			return "CREATE AUTO_LOGIN KEYSTORE"
		case p.cur.Type == kwKEY:
			p.advance()
			return "CREATE KEY"
		case p.isIdentLikeStr("KEYSTORE"):
			p.advance()
			return "CREATE KEYSTORE"
		}
	case p.cur.Type == kwSET:
		p.advance()
		switch {
		case p.cur.Type == kwKEY:
			p.advance()
			return "SET KEY"
		case p.isIdentLikeStr("KEYSTORE"):
			p.advance()
			if p.isIdentLikeStr("OPEN") || p.isIdentLikeStr("CLOSE") {
				action := p.cur.Str
				p.advance()
				return "SET KEYSTORE " + action
			}
			return "SET KEYSTORE"
		case p.isIdentLikeStr("TAG"):
			p.advance()
			return "SET TAG"
		}
	case p.isIdentLikeStr("USE"):
		p.advance()
		if p.cur.Type == kwKEY {
			p.advance()
		}
		return "USE KEY"
	case p.isIdentLikeStr("BACKUP"):
		p.advance()
		if p.isIdentLikeStr("KEYSTORE") {
			p.advance()
		}
		return "BACKUP KEYSTORE"
	case p.cur.Type == kwALTER:
		p.advance()
		if p.isIdentLikeStr("KEYSTORE") {
			p.advance()
		}
		if p.isIdentLikeStr("PASSWORD") {
			p.advance()
		}
		return "ALTER KEYSTORE PASSWORD"
	case p.isIdentLikeStr("EXPORT"):
		p.advance()
		if p.cur.Type == kwKEY || p.isIdentLikeStr("KEYS") {
			p.advance()
		}
		return "EXPORT KEYS"
	case p.isIdentLikeStr("IMPORT"):
		p.advance()
		if p.cur.Type == kwKEY || p.isIdentLikeStr("KEYS") {
			p.advance()
		}
		return "IMPORT KEYS"
	case p.cur.Type == kwMERGE:
		p.advance()
		if p.isIdentLikeStr("KEYSTORE") {
			p.advance()
		}
		return "MERGE KEYSTORE"
	case p.cur.Type == kwADD:
		p.advance()
		if p.isIdentLikeStr("SECRET") {
			p.advance()
		}
		return "ADD SECRET"
	case p.cur.Type == kwUPDATE:
		p.advance()
		if p.isIdentLikeStr("SECRET") {
			p.advance()
		}
		return "UPDATE SECRET"
	case p.cur.Type == kwDELETE:
		p.advance()
		if p.isIdentLikeStr("SECRET") {
			p.advance()
		}
		return "DELETE SECRET"
	}
	return ""
}

func (p *Parser) consumeAdministerKeyManagementValue() string {
	if p.cur.Type == tokEOF || p.cur.Type == ';' {
		return ""
	}
	if p.isIdentLikeStr("EXTERNAL") {
		p.advance()
		if p.isIdentLikeStr("STORE") {
			p.advance()
			return "EXTERNAL STORE"
		}
		return "EXTERNAL"
	}
	value := p.cur.Str
	p.advance()
	return value
}

// parseCreateAuditPolicyStmt parses a CREATE AUDIT POLICY statement (Unified Auditing).
// Called after CREATE AUDIT POLICY has been consumed.
//
// BNF: oracle/parser/bnf/CREATE-AUDIT-POLICY-Unified-Auditing.bnf
//
//	CREATE AUDIT POLICY policy
//	    [ privilege_audit_clause ]
//	    [ action_audit_clause ]
//	    [ role_audit_clause ]
//	    [ WHEN 'audit_condition' EVALUATE PER { STATEMENT | SESSION | INSTANCE } ]
//	    [ ONLY TOPLEVEL ]
//	    [ CONTAINER = { ALL | CURRENT } ] ;
//
//	privilege_audit_clause:
//	    PRIVILEGES system_privilege [, system_privilege ]...
//
//	action_audit_clause:
//	    ACTIONS [ standard_actions ] [ component_actions ]
//
//	standard_actions:
//	    { object_action [ ( column [, column ]... ) ]
//	        ON { [ schema. ] object_name
//	           | DIRECTORY directory_name
//	           | MINING MODEL [ schema. ] object_name }
//	    | ALL ON { [ schema. ] object_name
//	             | DIRECTORY directory_name
//	             | MINING MODEL [ schema. ] object_name }
//	    | system_action
//	    | ALL
//	    } [, { object_action [ ( column [, column ]... ) ]
//	            ON { [ schema. ] object_name
//	               | DIRECTORY directory_name
//	               | MINING MODEL [ schema. ] object_name }
//	          | ALL ON { [ schema. ] object_name
//	                   | DIRECTORY directory_name
//	                   | MINING MODEL [ schema. ] object_name }
//	          | system_action
//	          | ALL
//	          } ]...
//
//	component_actions:
//	    COMPONENT = { DATAPUMP { component_action | ALL } [, { component_action | ALL } ]...
//	                | DIRECT_LOAD { component_action | ALL } [, { component_action | ALL } ]...
//	                | OLS { component_action | ALL } [, { component_action | ALL } ]...
//	                | XS { component_action | ALL } [, { component_action | ALL } ]...
//	                | DV { component_action [ ON object_name ] | ALL } [, { component_action [ ON object_name ] | ALL } ]...
//	                | SQL_FIREWALL
//	                | PROTOCOL { HTTP | FTP | AUTHENTICATION }
//	                }
//
//	role_audit_clause:
//	    ROLES role [, role ]...
func (p *Parser) parseCreateAuditPolicyStmt(start int) (nodes.StmtNode, error) {
	stmt := &nodes.CreateAuditPolicyStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Policy name
	if p.isIdentLike() {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// Parse clauses in any order
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		switch {
		case p.cur.Type == kwPRIVILEGES:
			// privilege_audit_clause
			p.advance()
			var parseErr929 error
			stmt.Privileges, parseErr929 = p.parseAuditPrivilegeList()
			if parseErr929 != nil {
				return nil, parseErr929
			}

		case p.isIdentLikeStr("ACTIONS"):
			// action_audit_clause
			p.advance()
			parseErr930 := p.parseAuditActionsClause(&stmt.Actions, &stmt.ComponentActions)
			if parseErr930 != nil {
				return nil, parseErr930

				// role_audit_clause
			}

		case p.isIdentLikeStr("ROLES"):

			p.advance()
			var parseErr931 error
			stmt.Roles, parseErr931 = p.parseAuditIdentList()
			if parseErr931 != nil {
				return nil, parseErr931

				// WHEN 'audit_condition' EVALUATE PER { STATEMENT | SESSION | INSTANCE }
			}

		case p.cur.Type == kwWHEN:

			p.advance()
			if p.cur.Type == tokSCONST {
				stmt.WhenCondition = p.cur.Str
				p.advance()
			}
			if p.isIdentLikeStr("EVALUATE") {
				p.advance()
				if p.isIdentLikeStr("PER") {
					p.advance()
				}
				if p.isIdentLike() {
					stmt.EvaluatePer = p.cur.Str
					p.advance()
				}
			}

		case p.isIdentLikeStr("ONLY"):
			// ONLY TOPLEVEL
			p.advance()
			if p.isIdentLikeStr("TOPLEVEL") {
				p.advance()
				stmt.OnlyToplevel = true
			}

		case p.isIdentLikeStr("CONTAINER"):
			// CONTAINER = { ALL | CURRENT }
			p.advance()
			var parseErr932 error
			stmt.ContainerAll, parseErr932 = p.parseContainerClause()
			if parseErr932 != nil {
				return nil, parseErr932
			}

		default:
			goto createDone
		}
	}
createDone:

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterAuditPolicyStmt parses an ALTER AUDIT POLICY statement (Unified Auditing).
// Called after ALTER AUDIT POLICY has been consumed.
//
// BNF: oracle/parser/bnf/ALTER-AUDIT-POLICY-Unified-Auditing.bnf
//
//	alter_audit_policy ::=
//	  ALTER AUDIT POLICY policy
//	    { ADD | DROP }
//	      [ privilege_audit_clause ] [ action_audit_clause ] [ role_audit_clause ]
//	  | ALTER AUDIT POLICY policy
//	      CONDITION { DROP | 'audit_condition' [ EVALUATE { PER STATEMENT | PER SESSION | PER INSTANCE } ] }
//	  | ALTER AUDIT POLICY policy
//	      { ADD | DROP } ONLY TOPLEVEL
//
//	privilege_audit_clause ::=
//	  PRIVILEGES privilege [, privilege ]...
//
//	action_audit_clause ::=
//	  ACTIONS { standard_actions | component_actions } [, { standard_actions | component_actions } ]...
//
//	standard_actions ::=
//	  action [ ON [ schema . ] object [ ( column [, column ]... ) ] ]
//
//	component_actions ::=
//	  COMPONENT = component_name action [, action ]...
//
//	role_audit_clause ::=
//	  ROLES role [, role ]...
func (p *Parser) parseAlterAuditPolicyStmt(start int) (nodes.StmtNode, error) {
	stmt := &nodes.AlterAuditPolicyStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Policy name
	if p.isIdentLike() {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// { ADD | DROP } or CONDITION
	if p.isIdentLikeStr("ADD") || p.cur.Type == kwDROP {
		stmt.AddDrop = p.cur.Str
		p.advance()

		// Check for ONLY TOPLEVEL
		if p.isIdentLikeStr("ONLY") {
			p.advance()
			if p.isIdentLikeStr("TOPLEVEL") {
				p.advance()
				if stmt.AddDrop == "ADD" {
					stmt.AddToplevel = true
				} else {
					stmt.DropToplevel = true
				}
				stmt.Loc.End = p.prev.End
				return stmt, nil
			}
		}

		// Parse sub-clauses
		for p.cur.Type != ';' && p.cur.Type != tokEOF {
			switch {
			case p.cur.Type == kwPRIVILEGES:
				p.advance()
				var parseErr933 error
				stmt.Privileges, parseErr933 = p.parseAuditPrivilegeList()
				if parseErr933 != nil {
					return nil, parseErr933
				}
			case p.isIdentLikeStr("ACTIONS"):
				p.advance()
				parseErr934 := p.parseAuditActionsClause(&stmt.Actions, &stmt.ComponentActions)
				if parseErr934 != nil {
					return nil, parseErr934
				}
			case p.isIdentLikeStr("ROLES"):
				p.advance()
				var parseErr935 error
				stmt.Roles, parseErr935 = p.parseAuditIdentList()
				if parseErr935 != nil {
					return nil, parseErr935
				}
			default:
				goto alterDone
			}
		}
	} else if p.isIdentLikeStr("CONDITION") {
		p.advance()
		if p.cur.Type == kwDROP {
			stmt.ConditionDrop = true
			p.advance()
		} else if p.cur.Type == tokSCONST {
			stmt.Condition = p.cur.Str
			p.advance()
			// Optional EVALUATE PER { STATEMENT | SESSION | INSTANCE }
			if p.isIdentLikeStr("EVALUATE") {
				p.advance()
				if p.isIdentLikeStr("PER") {
					p.advance()
				}
				if p.isIdentLike() {
					stmt.EvaluatePer = p.cur.Str
					p.advance()
				}
			}
		}
	}
alterDone:

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseDropAuditPolicyStmt parses a DROP AUDIT POLICY statement (Unified Auditing).
// Called after DROP AUDIT POLICY has been consumed.
//
// BNF: oracle/parser/bnf/DROP-AUDIT-POLICY-Unified-Auditing.bnf
//
//	DROP AUDIT POLICY policy
func (p *Parser) parseDropAuditPolicyStmt(start int) (nodes.StmtNode, error) {
	stmt := &nodes.DropAuditPolicyStmt{
		Loc: nodes.Loc{Start: start},
	}

	if p.isIdentLike() {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAuditActionsClause parses standard_actions and component_actions for audit policy.
// Called after ACTIONS keyword has been consumed.
func (p *Parser) parseAuditActionsClause(actions *[]*nodes.AuditActionEntry, compActions *[]*nodes.AuditComponentAction) error {
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		// Check for COMPONENT =
		if p.isIdentLikeStr("COMPONENT") {
			actStart := p.pos()
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			comp := &nodes.AuditComponentAction{
				Loc: nodes.Loc{Start: actStart},
			}
			if p.isIdentLike() {
				comp.Component = p.cur.Str
				p.advance()
			}
			// Parse component actions
			for {
				if p.isIdentLike() || p.cur.Type == kwALL {
					comp.Actions = append(comp.Actions, p.cur.Str)
					p.advance()
					// DV: optional ON object_name
					if p.cur.Type == kwON {
						p.advance()
						if p.isIdentLike() {
							comp.Object = p.cur.Str
							p.advance()
						}
					}
				}
				if p.cur.Type != ',' {
					break
				}
				p.advance()
			}
			comp.Loc.End = p.prev.End
			*compActions = append(*compActions, comp)
			continue
		}

		// Check if we hit a clause keyword that ends actions
		if p.cur.Type == kwPRIVILEGES || p.isIdentLikeStr("ROLES") || p.cur.Type == kwWHEN ||
			p.isIdentLikeStr("ONLY") || p.isIdentLikeStr("CONTAINER") ||
			p.isIdentLikeStr("CONDITION") {
			break
		}

		// Parse standard action entry
		entry, parseErr936 := p.parseAuditActionEntry()
		if parseErr936 != nil {
			return parseErr936
		}
		if entry == nil {
			break
		}
		*actions = append(*actions, entry)

		if p.cur.Type == ',' {
			p.advance()
		} else {
			break
		}
	}
	return nil
}

// parseAuditActionEntry parses a single standard audit action entry.
func (p *Parser) parseAuditActionEntry() (*nodes.AuditActionEntry, error) {
	if !p.isIdentLike() && p.cur.Type != kwALL && p.cur.Type != kwSELECT &&
		p.cur.Type != kwINSERT && p.cur.Type != kwUPDATE && p.cur.Type != kwDELETE &&
		p.cur.Type != kwCREATE && p.cur.Type != kwALTER && p.cur.Type != kwDROP &&
		p.cur.Type != kwGRANT && p.cur.Type != kwEXECUTE {
		return nil, nil
	}

	start := p.pos()
	entry := &nodes.AuditActionEntry{
		Loc: nodes.Loc{Start: start},
	}

	// Collect action name (may be multi-word like "CREATE TABLE")
	action := p.cur.Str
	p.advance()

	// Check for ON immediately after ALL (ALL ON = all actions on specific object)
	if action == "ALL" && p.cur.Type == kwON {
		entry.Action = "ALL"
		p.advance()
		parseErr937 := p.parseAuditOnTarget(entry)
		if parseErr937 != nil {
			return nil, parseErr937
		}
		entry.Loc.End = p.prev.End
		return entry, nil
	}

	// Multi-word action: keep collecting until we see ON, comma, semicolon, or clause keyword
	for p.isIdentLike() || p.cur.Type == kwALL {
		// Stop at ON (introduces object target) or clause keywords
		if p.cur.Type == kwON {
			break
		}
		str := p.cur.Str
		if str == "COMPONENT" || str == "ROLES" || str == "ONLY" || str == "CONTAINER" || str == "CONDITION" {
			break
		}
		action += " " + str
		p.advance()
	}
	entry.Action = action

	// Optional (column [, column]...)
	if p.cur.Type == '(' {
		p.advance()
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			if p.isIdentLike() {
				entry.Columns = append(entry.Columns, p.cur.Str)
				p.advance()
			}
			if p.cur.Type == ',' {
				p.advance()
			} else {
				break
			}
		}
		if p.cur.Type == ')' {
			p.advance()
		}
	}

	// Optional ON target
	if p.cur.Type == kwON {
		p.advance()
		parseErr938 := p.parseAuditOnTarget(entry)
		if parseErr938 != nil {
			return nil, parseErr938
		}
	}

	entry.Loc.End = p.prev.End
	return entry, nil
}

// parseAuditOnTarget parses the ON target for an audit action entry.
// Called after ON has been consumed.
func (p *Parser) parseAuditOnTarget(entry *nodes.AuditActionEntry) error {
	if p.isIdentLikeStr("DIRECTORY") {
		p.advance()
		if p.isIdentLike() {
			entry.Directory = p.cur.Str
			p.advance()
		}
	} else if p.isIdentLikeStr("MINING") {
		p.advance()
		if p.cur.Type == kwMODEL {
			p.advance()
		}
		entry.MiningModel = true
		var parseErr939 error
		entry.Object, parseErr939 = p.parseObjectName()
		if parseErr939 != nil {
			return parseErr939
		}
	} else {
		var parseErr940 error
		entry.Object, parseErr940 = p.parseObjectName()
		if parseErr940 !=

			// parsePrivilegeList parses a comma-separated list of privileges (may be multi-word).
			// Called after PRIVILEGES keyword has been consumed.
			nil {
			return parseErr940
		}
	}
	return nil
}

func (p *Parser) parseAuditPrivilegeList() ([]string, error) {
	var privs []string
	for {
		priv := ""
		for p.isIdentLike() || p.cur.Type == kwALL || p.cur.Type == kwCREATE ||
			p.cur.Type == kwALTER || p.cur.Type == kwDROP || p.cur.Type == kwSELECT ||
			p.cur.Type == kwINSERT || p.cur.Type == kwUPDATE || p.cur.Type == kwDELETE ||
			p.cur.Type == kwEXECUTE || p.cur.Type == kwGRANT || p.cur.Type == kwINDEX {
			// Stop if this is a clause keyword
			str := p.cur.Str
			if str == "ACTIONS" || str == "ROLES" || str == "COMPONENT" ||
				str == "ONLY" || str == "CONTAINER" || str == "CONDITION" {
				break
			}
			if p.cur.Type == kwWHEN {
				break
			}
			if priv != "" {
				priv += " "
			}
			priv += str
			p.advance()
		}
		if priv == "" {
			break
		}
		privs = append(privs, priv)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return privs, nil
}

// parseAuditIdentList parses a comma-separated list of simple identifiers.
func (p *Parser) parseAuditIdentList() ([]string, error) {
	var list []string
	for {
		if !p.isIdentLike() {
			break
		}
		list = append(list, p.cur.Str)
		p.advance()
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return list, nil
}
