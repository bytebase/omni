package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseContainerClause parses CONTAINER = { ALL | CURRENT }.
// Returns: *bool (true=ALL, false=CURRENT), or nil if not present.
func (p *Parser) parseContainerClause() (*bool, error) {
	// Already consumed CONTAINER
	if p.cur.Type != '=' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	if p.cur.Type == kwALL {
		p.advance()
		v := true
		return &v, nil
	} else if p.isIdentLikeStr("CURRENT") {
		p.advance()
		v := false
		return &v, nil
	}
	return nil, p.syntaxErrorAtCur()
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
	// IF NOT EXISTS: Oracle reports an incomplete one at IF (ORA-11543).
	if p.cur.Type == kwIF {
		ifTok := p.advance()
		if p.cur.Type != kwNOT || p.peekNext().Type != kwEXISTS {
			return nil, p.syntaxErrorAtTok(ifTok)
		}
		p.advance()
		p.advance()
		stmt.IfNotExists = true
	}
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	switch {
	case p.cur.Type == kwIDENTIFIED:
		if stmt.Identified, err = p.parseUserIdentified(false); err != nil {
			return nil, err
		}
	case p.isKeywordStr("NO"):
		if stmt.Identified, err = p.parseUserNoAuthentication(); err != nil {
			return nil, err
		}
	}
	var c userClauses
	if err := p.parseUserClauses(false, &c); err != nil {
		return nil, err
	}
	stmt.DefaultTablespace, stmt.TempTablespace, stmt.LocalTemp = c.defaultTablespace, c.tempTablespace, c.localTemp
	stmt.Quotas, stmt.Profile, stmt.PasswordExpire = c.quotas, c.profile, c.passwordExpire
	stmt.AccountLock, stmt.EnableEditions = c.accountLock, c.enableEditions
	stmt.DefaultCollation, stmt.ContainerAll = c.defaultCollation, c.containerAll
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
	if p.cur.Type == kwIF {
		p.advance()
		if p.cur.Type != kwEXISTS {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		stmt.IfExists = true
	}
	if !p.isIdentLike() {
		return nil, p.syntaxErrorAtCur()
	}
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	if p.cur.Type == kwGRANT || p.cur.Type == kwREVOKE {
		if err := p.parseUserProxyClause(); err != nil {
			return nil, err
		}
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}
	switch {
	case p.cur.Type == kwIDENTIFIED:
		if stmt.Identified, err = p.parseUserIdentified(true); err != nil {
			return nil, err
		}
	case p.isKeywordStr("NO"):
		if stmt.Identified, err = p.parseUserNoAuthentication(); err != nil {
			return nil, err
		}
	}
	var c userClauses
	if err := p.parseUserClauses(true, &c); err != nil {
		return nil, err
	}
	// ALTER USER takes at least one clause, CONTAINER aside (ORA-02157).
	if stmt.Identified == nil && c.options == 0 {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.DefaultTablespace, stmt.TempTablespace, stmt.LocalTemp = c.defaultTablespace, c.tempTablespace, c.localTemp
	stmt.Quotas, stmt.Profile, stmt.DefaultRole, stmt.PasswordExpire = c.quotas, c.profile, c.defaultRole, c.passwordExpire
	stmt.AccountLock, stmt.EnableEditions = c.accountLock, c.enableEditions
	stmt.DefaultCollation, stmt.ContainerAll = c.defaultCollation, c.containerAll
	stmt.Loc.End = p.prev.End
	return stmt, nil
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

	stmt.Name, parseErr914 = p.parseRequiredObjectName()
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
		if p.cur.Type != kwIDENTIFIED {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
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
		if !p.isIdentLike() && p.cur.Type != tokSCONST {
			return p.syntaxErrorAtCur()
		}
		setValue(p.cur.Str)
		p.advance()
	case p.cur.Type == kwUSING:
		// IDENTIFIED USING [schema.]package
		p.advance()
		setType(nodes.ROLE_IDENTIFIED_USING)
		name, parseErr917 := p.parseRequiredObjectName()
		if parseErr917 != nil {
			return parseErr917
		}
		if name.Schema != "" {
			setSchema(name.Schema)
		}
		setValue(name.Name)
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
	default:
		return p.syntaxErrorAtCur()
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

	stmt.Name, parseErr918 = p.parseRequiredObjectName()
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
		if p.cur.Type != kwIDENTIFIED {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		stmt.IdentifiedType = nodes.ROLE_NOT_IDENTIFIED
	} else {
		return nil, p.syntaxErrorAtCur()
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

	if err := p.parseProfileLimits(&stmt.Limits, &stmt.ContainerAll); err != nil {
		return nil, err
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
			} else if len(stmt.Costs) == 0 {
				// The first cost takes its value (ORA-00923 at the end).
				// Oracle 23ai takes a name without one after it.
				return nil, p.syntaxErrorAtCur()
			}
			entry.Loc.End = p.prev.End
			stmt.Costs = append(stmt.Costs, entry)
		default:
			goto done
		}
	}
done:
	if len(stmt.Costs) == 0 {
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
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

	if err := p.parseProfileLimits(&stmt.Limits, &stmt.ContainerAll); err != nil {
		return nil, err
	}

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
