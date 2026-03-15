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
func (p *Parser) parseIdentifiedClause(allowReplace bool) *nodes.IdentifiedClause {
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

	clause.Loc.End = p.pos()
	return clause
}

// parseUserQuotaClause parses a QUOTA clause.
//
//	QUOTA { size_clause | UNLIMITED } ON tablespace
//
//	size_clause ::= integer [ K | M | G | T ]
func (p *Parser) parseUserQuotaClause() *nodes.UserQuotaClause {
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
		clause.Tablespace = p.parseObjectName()
	}

	clause.Loc.End = p.pos()
	return clause
}

// parseContainerClause parses CONTAINER = { ALL | CURRENT }.
// Returns: *bool (true=ALL, false=CURRENT), or nil if not present.
func (p *Parser) parseContainerClause() *bool {
	// Already consumed CONTAINER
	if p.cur.Type == '=' {
		p.advance()
	}
	if p.cur.Type == kwALL {
		p.advance()
		v := true
		return &v
	} else if p.isIdentLikeStr("CURRENT") {
		p.advance()
		v := false
		return &v
	}
	return nil
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
) {
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
				return
			} else {
				// Unknown DEFAULT clause, skip
				return
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
			addQuota(p.parseUserQuotaClause())

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
				return
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
				return
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
				return
			}

		case p.isIdentLikeStr("CONTAINER"):
			// CONTAINER = { ALL | CURRENT }
			p.advance()
			setContainer(p.parseContainerClause())

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
					p.parseObjectName()
				}
			} else {
				return
			}

		default:
			return
		}
	}
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
func (p *Parser) parseCreateUserStmt(start int) nodes.StmtNode {
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

	stmt.Name = p.parseObjectName()

	// IDENTIFIED clause or NO AUTHENTICATION
	if p.cur.Type == kwIDENTIFIED {
		p.advance()
		stmt.Identified = p.parseIdentifiedClause(false)
	} else if p.isIdentLikeStr("NO") {
		p.advance()
		if p.isIdentLikeStr("AUTHENTICATION") {
			p.advance()
			stmt.Identified = &nodes.IdentifiedClause{
				Type: nodes.IDENTIFIED_NO_AUTH,
				Loc:  nodes.Loc{Start: p.pos(), End: p.pos()},
			}
		}
	}

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

	stmt.Loc.End = p.pos()
	return stmt
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
func (p *Parser) parseAlterUserStmt(start int) nodes.StmtNode {
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

	stmt.Name = p.parseObjectName()

	// IDENTIFIED clause or NO AUTHENTICATION
	if p.cur.Type == kwIDENTIFIED {
		p.advance()
		stmt.Identified = p.parseIdentifiedClause(true)
	} else if p.isIdentLikeStr("NO") {
		p.advance()
		if p.isIdentLikeStr("AUTHENTICATION") {
			p.advance()
			stmt.Identified = &nodes.IdentifiedClause{
				Type: nodes.IDENTIFIED_NO_AUTH,
				Loc:  nodes.Loc{Start: p.pos(), End: p.pos()},
			}
		}
	}

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

	// DEFAULT ROLE clause (parseUserOptions returns when it sees DEFAULT ROLE)
	if p.cur.Type == kwROLE {
		p.advance()
		stmt.DefaultRole = p.parseDefaultRoleClause()

		// Continue parsing remaining options after DEFAULT ROLE
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
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseDefaultRoleClause parses the DEFAULT ROLE clause for ALTER USER.
//
//	DEFAULT ROLE { role [, role ]... | ALL [ EXCEPT role [, role ]... ] | NONE }
func (p *Parser) parseDefaultRoleClause() *nodes.DefaultRoleClause {
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
				clause.Roles = append(clause.Roles, p.parseObjectName())
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
			clause.Roles = append(clause.Roles, p.parseObjectName())
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}

	clause.Loc.End = p.pos()
	return clause
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
func (p *Parser) parseCreateRoleStmt(start int) nodes.StmtNode {
	stmt := &nodes.CreateRoleStmt{
		Loc: nodes.Loc{Start: start},
	}

	stmt.Name = p.parseObjectName()

	// Optional IDENTIFIED clause
	if p.cur.Type == kwIDENTIFIED {
		p.advance()
		stmt.HasIdentified = true
		p.parseRoleIdentifiedClause(
			func(t nodes.RoleIdentifiedType) { stmt.IdentifiedType = t },
			func(v string) { stmt.IdentifyBy = v },
			func(v string) { stmt.IdentifySchema = v },
		)
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
		stmt.ContainerAll = p.parseContainerClause()
	}

	stmt.Loc.End = p.pos()
	return stmt
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
) {
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
		name := p.parseObjectName()
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
func (p *Parser) parseAlterRoleStmt(start int) nodes.StmtNode {
	stmt := &nodes.AlterRoleStmt{
		Loc: nodes.Loc{Start: start},
	}

	stmt.Name = p.parseObjectName()

	// Required IDENTIFIED clause or NOT IDENTIFIED
	if p.cur.Type == kwIDENTIFIED {
		p.advance()
		p.parseRoleIdentifiedClause(
			func(t nodes.RoleIdentifiedType) { stmt.IdentifiedType = t },
			func(v string) { stmt.IdentifyBy = v },
			func(v string) { stmt.IdentifySchema = v },
		)
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
		stmt.ContainerAll = p.parseContainerClause()
	}

	stmt.Loc.End = p.pos()
	return stmt
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
func (p *Parser) parseAlterProfileStmt(start int) nodes.StmtNode {
	stmt := &nodes.AlterProfileStmt{
		Loc: nodes.Loc{Start: start},
	}

	stmt.Name = p.parseObjectName()

	// Parse LIMIT clauses (same logic as CREATE PROFILE)
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		if p.cur.Type == kwLIMIT {
			p.advance()
			for p.isProfileParam() {
				lim := p.parseProfileLimit()
				stmt.Limits = append(stmt.Limits, lim)
			}
		} else if p.isProfileParam() {
			lim := p.parseProfileLimit()
			stmt.Limits = append(stmt.Limits, lim)
		} else if p.isIdentLikeStr("CONTAINER") {
			p.advance()
			stmt.ContainerAll = p.parseContainerClause()
		} else {
			break
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
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
func (p *Parser) parseAlterResourceCostStmt(start int) nodes.StmtNode {
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
			entry.Loc.End = p.pos()
			stmt.Costs = append(stmt.Costs, entry)
		default:
			goto done
		}
	}
done:

	stmt.Loc.End = p.pos()
	return stmt
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
func (p *Parser) parseCreateProfileStmt(start int, mandatory bool) nodes.StmtNode {
	stmt := &nodes.CreateProfileStmt{
		Loc:       nodes.Loc{Start: start},
		Mandatory: mandatory,
	}

	stmt.Name = p.parseObjectName()

	// Parse LIMIT clauses
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		if p.cur.Type == kwLIMIT {
			p.advance()
			// Parse parameters after LIMIT
			for p.isProfileParam() {
				lim := p.parseProfileLimit()
				stmt.Limits = append(stmt.Limits, lim)
			}
		} else if p.isProfileParam() {
			// Parameters can appear without repeated LIMIT keyword
			lim := p.parseProfileLimit()
			stmt.Limits = append(stmt.Limits, lim)
		} else if p.isIdentLikeStr("CONTAINER") {
			p.advance()
			stmt.ContainerAll = p.parseContainerClause()
		} else {
			break
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseProfileLimit parses a single profile limit parameter and its value.
//
//	param_name { integer | size_clause | UNLIMITED | DEFAULT | NULL | function_name }
func (p *Parser) parseProfileLimit() *nodes.ProfileLimit {
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

	lim.Loc.End = p.pos()
	return lim
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
func (p *Parser) parseAdministerKeyManagementStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume ADMINISTER

	stmt := &nodes.AdminDDLStmt{
		Action:     "ADMINISTER",
		ObjectType: nodes.OBJECT_KEY_MANAGEMENT,
		Loc:        nodes.Loc{Start: start},
	}

	// Skip KEY MANAGEMENT and all remaining tokens
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt
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
func (p *Parser) parseCreateAuditPolicyStmt(start int) nodes.StmtNode {
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
			stmt.Privileges = p.parseAuditPrivilegeList()

		case p.isIdentLikeStr("ACTIONS"):
			// action_audit_clause
			p.advance()
			p.parseAuditActionsClause(&stmt.Actions, &stmt.ComponentActions)

		case p.isIdentLikeStr("ROLES"):
			// role_audit_clause
			p.advance()
			stmt.Roles = p.parseAuditIdentList()

		case p.cur.Type == kwWHEN:
			// WHEN 'audit_condition' EVALUATE PER { STATEMENT | SESSION | INSTANCE }
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
			stmt.ContainerAll = p.parseContainerClause()

		default:
			goto createDone
		}
	}
createDone:

	stmt.Loc.End = p.pos()
	return stmt
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
func (p *Parser) parseAlterAuditPolicyStmt(start int) nodes.StmtNode {
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
				stmt.Loc.End = p.pos()
				return stmt
			}
		}

		// Parse sub-clauses
		for p.cur.Type != ';' && p.cur.Type != tokEOF {
			switch {
			case p.cur.Type == kwPRIVILEGES:
				p.advance()
				stmt.Privileges = p.parseAuditPrivilegeList()
			case p.isIdentLikeStr("ACTIONS"):
				p.advance()
				p.parseAuditActionsClause(&stmt.Actions, &stmt.ComponentActions)
			case p.isIdentLikeStr("ROLES"):
				p.advance()
				stmt.Roles = p.parseAuditIdentList()
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

	stmt.Loc.End = p.pos()
	return stmt
}

// parseDropAuditPolicyStmt parses a DROP AUDIT POLICY statement (Unified Auditing).
// Called after DROP AUDIT POLICY has been consumed.
//
// BNF: oracle/parser/bnf/DROP-AUDIT-POLICY-Unified-Auditing.bnf
//
//	DROP AUDIT POLICY policy
func (p *Parser) parseDropAuditPolicyStmt(start int) nodes.StmtNode {
	stmt := &nodes.DropAuditPolicyStmt{
		Loc: nodes.Loc{Start: start},
	}

	if p.isIdentLike() {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseAuditActionsClause parses standard_actions and component_actions for audit policy.
// Called after ACTIONS keyword has been consumed.
func (p *Parser) parseAuditActionsClause(actions *[]*nodes.AuditActionEntry, compActions *[]*nodes.AuditComponentAction) {
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
			comp.Loc.End = p.pos()
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
		entry := p.parseAuditActionEntry()
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
}

// parseAuditActionEntry parses a single standard audit action entry.
func (p *Parser) parseAuditActionEntry() *nodes.AuditActionEntry {
	if !p.isIdentLike() && p.cur.Type != kwALL && p.cur.Type != kwSELECT &&
		p.cur.Type != kwINSERT && p.cur.Type != kwUPDATE && p.cur.Type != kwDELETE &&
		p.cur.Type != kwCREATE && p.cur.Type != kwALTER && p.cur.Type != kwDROP &&
		p.cur.Type != kwGRANT && p.cur.Type != kwEXECUTE {
		return nil
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
		p.parseAuditOnTarget(entry)
		entry.Loc.End = p.pos()
		return entry
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
		p.parseAuditOnTarget(entry)
	}

	entry.Loc.End = p.pos()
	return entry
}

// parseAuditOnTarget parses the ON target for an audit action entry.
// Called after ON has been consumed.
func (p *Parser) parseAuditOnTarget(entry *nodes.AuditActionEntry) {
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
		entry.Object = p.parseObjectName()
	} else {
		entry.Object = p.parseObjectName()
	}
}

// parsePrivilegeList parses a comma-separated list of privileges (may be multi-word).
// Called after PRIVILEGES keyword has been consumed.
func (p *Parser) parseAuditPrivilegeList() []string {
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
	return privs
}

// parseAuditIdentList parses a comma-separated list of simple identifiers.
func (p *Parser) parseAuditIdentList() []string {
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
	return list
}
