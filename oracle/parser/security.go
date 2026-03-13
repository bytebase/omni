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
			// ENABLE EDITIONS
			p.advance()
			if p.isIdentLikeStr("EDITIONS") {
				p.advance()
				setEnableEditions()
			} else {
				return
			}

		case p.isIdentLikeStr("CONTAINER"):
			// CONTAINER = { ALL | CURRENT }
			p.advance()
			setContainer(p.parseContainerClause())

		default:
			return
		}
	}
}

// parseCreateUserStmt parses a CREATE USER statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/CREATE-USER.html
//
//	CREATE USER [ IF NOT EXISTS ] user
//	    IDENTIFIED
//	        { BY password
//	        | EXTERNALLY [ AS 'certificate_DN' | AS 'kerberos_principal_name' ]
//	        | GLOBALLY [ AS directory_DN | AS 'AZURE_USER=...' | AS 'AZURE_ROLE=...' ]
//	        | NO AUTHENTICATION
//	        }
//	    [ DEFAULT COLLATION collation_name ]
//	    [ DEFAULT TABLESPACE tablespace ]
//	    [ [ LOCAL ] TEMPORARY TABLESPACE { tablespace | tablespace_group_name } ]
//	    [ QUOTA { size_clause | UNLIMITED } ON tablespace ] ...
//	    [ PROFILE profile_name ]
//	    [ PASSWORD EXPIRE ]
//	    [ ACCOUNT { LOCK | UNLOCK } ]
//	    [ ENABLE EDITIONS ]
//	    [ CONTAINER = { ALL | CURRENT } ]
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
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/ALTER-USER.html
//
//	ALTER USER [ IF EXISTS ] user
//	    [ IDENTIFIED { BY password [ REPLACE old_password ]
//	                 | GLOBALLY AS 'distinguished_name'
//	                 | EXTERNALLY
//	                 | NO AUTHENTICATION
//	                 } ]
//	    [ DEFAULT COLLATION collation_name ]
//	    [ DEFAULT TABLESPACE tablespace ]
//	    [ [ LOCAL ] TEMPORARY TABLESPACE { tablespace | tablespace_group_name } ]
//	    [ QUOTA { size_clause | UNLIMITED } ON tablespace ] ...
//	    [ PROFILE profile_name ]
//	    [ DEFAULT ROLE { role [, role ]...
//	                   | ALL [ EXCEPT role [, role ]... ]
//	                   | NONE
//	                   } ]
//	    [ PASSWORD EXPIRE ]
//	    [ ACCOUNT { LOCK | UNLOCK } ]
//	    [ ENABLE EDITIONS ]
//	    [ CONTAINER = { ALL | CURRENT } ]
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
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/CREATE-ROLE.html
//
//	CREATE ROLE role [NOT IDENTIFIED | IDENTIFIED BY password | IDENTIFIED USING package | IDENTIFIED EXTERNALLY | IDENTIFIED GLOBALLY]
func (p *Parser) parseCreateRoleStmt(start int) nodes.StmtNode {
	stmt := &nodes.CreateRoleStmt{
		Loc: nodes.Loc{Start: start},
	}

	stmt.Name = p.parseObjectName()

	// Optional IDENTIFIED clause
	if p.cur.Type == kwIDENTIFIED {
		p.advance()
		if p.cur.Type == kwBY {
			p.advance()
			if p.isIdentLike() || p.cur.Type == tokSCONST {
				stmt.IdentifyBy = p.cur.Str
				p.advance()
			}
		} else if p.isIdentLike() {
			stmt.IdentifyBy = p.cur.Str
			p.advance()
		}
	} else if p.cur.Type == kwNOT {
		p.advance()
		if p.cur.Type == kwIDENTIFIED {
			p.advance()
		}
		stmt.IdentifyBy = "NOT IDENTIFIED"
	}

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
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/CREATE-PROFILE.html
//
//	CREATE [ MANDATORY ] PROFILE profile_name
//	    LIMIT { resource_parameters | password_parameters } ...
//	    [ CONTAINER = { ALL | CURRENT } ]
//
//	resource_parameters ::=
//	    SESSIONS_PER_USER { integer | UNLIMITED | DEFAULT }
//	    CPU_PER_SESSION { integer | UNLIMITED | DEFAULT }
//	    CPU_PER_CALL { integer | UNLIMITED | DEFAULT }
//	    CONNECT_TIME { integer | UNLIMITED | DEFAULT }
//	    IDLE_TIME { integer | UNLIMITED | DEFAULT }
//	    LOGICAL_READS_PER_SESSION { integer | UNLIMITED | DEFAULT }
//	    LOGICAL_READS_PER_CALL { integer | UNLIMITED | DEFAULT }
//	    PRIVATE_SGA { size_clause | UNLIMITED | DEFAULT }
//	    COMPOSITE_LIMIT { integer | UNLIMITED | DEFAULT }
//
//	password_parameters ::=
//	    FAILED_LOGIN_ATTEMPTS { integer | UNLIMITED | DEFAULT }
//	    PASSWORD_LIFE_TIME { expr | UNLIMITED | DEFAULT }
//	    PASSWORD_REUSE_TIME { expr | UNLIMITED | DEFAULT }
//	    PASSWORD_REUSE_MAX { integer | UNLIMITED | DEFAULT }
//	    PASSWORD_LOCK_TIME { expr | UNLIMITED | DEFAULT }
//	    PASSWORD_GRACE_TIME { expr | UNLIMITED | DEFAULT }
//	    INACTIVE_ACCOUNT_TIME { expr | UNLIMITED | DEFAULT }
//	    PASSWORD_VERIFY_FUNCTION { function_name | NULL }
//	    PASSWORD_ROLLOVER_TIME { expr | UNLIMITED | DEFAULT }
//
//	size_clause ::= integer [ K | M | G ]
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
