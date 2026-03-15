package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseLockTableStmt parses a LOCK TABLE statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/LOCK-TABLE.html
//
//	LOCK TABLE [schema.]table IN lock_mode MODE [NOWAIT | WAIT integer]
//	lock_mode ::= ROW SHARE | ROW EXCLUSIVE | SHARE | SHARE UPDATE |
//	              SHARE ROW EXCLUSIVE | EXCLUSIVE
func (p *Parser) parseLockTableStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume LOCK

	if p.cur.Type == kwTABLE {
		p.advance()
	}

	stmt := &nodes.LockTableStmt{
		Loc: nodes.Loc{Start: start},
	}

	stmt.Table = p.parseObjectName()

	// IN
	if p.cur.Type == kwIN {
		p.advance()
	}

	// Lock mode: collect words until MODE
	mode := ""
	for p.cur.Type != kwMODE && p.cur.Type != tokEOF && p.cur.Type != ';' {
		if mode != "" {
			mode += " "
		}
		if p.cur.Type == kwSHARE {
			mode += "SHARE"
		} else if p.cur.Type == kwROW {
			mode += "ROW"
		} else if p.cur.Type == kwEXCLUSIVE {
			mode += "EXCLUSIVE"
		} else if p.isIdentLike() {
			mode += p.cur.Str
		}
		p.advance()
	}
	stmt.LockMode = mode

	// MODE
	if p.cur.Type == kwMODE {
		p.advance()
	}

	// NOWAIT or WAIT n
	if p.cur.Type == kwNOWAIT {
		stmt.Nowait = true
		p.advance()
	} else if p.cur.Type == kwWAIT {
		p.advance()
		stmt.Wait = p.parseExpr()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseCallStmt parses a CALL statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/CALL.html
//
//	CALL [schema.]routine_name ( [args] ) [INTO :bind_variable]
func (p *Parser) parseCallStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume CALL

	stmt := &nodes.CallStmt{
		Args: &nodes.List{},
		Loc:  nodes.Loc{Start: start},
	}

	stmt.Name = p.parseObjectName()

	// Arguments
	if p.cur.Type == '(' {
		p.advance()
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			arg := p.parseExpr()
			if arg != nil {
				stmt.Args.Items = append(stmt.Args.Items, arg)
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if p.cur.Type == ')' {
			p.advance()
		}
	}

	// INTO :bind_variable
	if p.cur.Type == kwINTO {
		p.advance()
		stmt.Into = p.parseExpr()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseRenameStmt parses a RENAME statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/RENAME.html
//
//	RENAME old_name TO new_name
func (p *Parser) parseRenameStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume RENAME

	stmt := &nodes.RenameStmt{
		Loc: nodes.Loc{Start: start},
	}

	stmt.OldName = p.parseObjectName()

	if p.cur.Type == kwTO {
		p.advance()
	}

	stmt.NewName = p.parseObjectName()

	stmt.Loc.End = p.pos()
	return stmt
}

// parseTruncateStmt parses a TRUNCATE TABLE statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/TRUNCATE-TABLE.html
//
//	TRUNCATE TABLE [schema.]table
//	    [{ PRESERVE | PURGE } MATERIALIZED VIEW LOG]
//	    [{ DROP | REUSE } STORAGE]
//	    [CASCADE]
func (p *Parser) parseTruncateStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume TRUNCATE

	stmt := &nodes.TruncateStmt{
		Loc: nodes.Loc{Start: start},
	}

	// TRUNCATE TABLE or TRUNCATE CLUSTER
	if p.cur.Type == kwTABLE {
		p.advance()
	} else if p.cur.Type == kwCLUSTER {
		stmt.Cluster = true
		p.advance()
	}

	// Parse table/cluster name
	stmt.Table = p.parseObjectName()

	// Parse optional clauses
	for {
		if p.cur.Type == kwPURGE {
			// PURGE MATERIALIZED VIEW LOG
			p.advance()
			if p.cur.Type == kwMATERIALIZED {
				stmt.PurgeMVLog = true
				p.advance() // consume MATERIALIZED
				if p.cur.Type == kwVIEW {
					p.advance()
				}
				if p.cur.Type == kwLOG {
					p.advance()
				}
			}
		} else if p.cur.Type == kwCASCADE {
			stmt.Cascade = true
			p.advance()
		} else if p.cur.Type == kwDROP || p.isIdentLike() && p.cur.Str == "REUSE" {
			// DROP STORAGE or REUSE STORAGE
			p.advance()
			if p.cur.Type == kwSTORAGE {
				p.advance()
			}
		} else if p.isIdentLike() && p.cur.Str == "PRESERVE" {
			// PRESERVE MATERIALIZED VIEW LOG
			p.advance()
			if p.cur.Type == kwMATERIALIZED {
				p.advance()
				if p.cur.Type == kwVIEW {
					p.advance()
				}
				if p.cur.Type == kwLOG {
					p.advance()
				}
			}
		} else {
			break
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseAnalyzeStmt parses an ANALYZE statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/ANALYZE.html
//
//	ANALYZE { TABLE | INDEX } [schema.]name
//	    { COMPUTE STATISTICS | ESTIMATE STATISTICS | DELETE STATISTICS | VALIDATE STRUCTURE }
func (p *Parser) parseAnalyzeStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume ANALYZE

	stmt := &nodes.AnalyzeStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Object type: TABLE or INDEX
	switch p.cur.Type {
	case kwTABLE:
		stmt.ObjectType = nodes.OBJECT_TABLE
		p.advance()
	case kwINDEX:
		stmt.ObjectType = nodes.OBJECT_INDEX
		p.advance()
	default:
		stmt.ObjectType = nodes.OBJECT_TABLE
	}

	// Object name
	stmt.Table = p.parseObjectName()

	// Action: COMPUTE STATISTICS, ESTIMATE STATISTICS, DELETE STATISTICS, VALIDATE STRUCTURE
	if p.isIdentLike() {
		action := p.cur.Str
		p.advance()
		// Second word of the action
		if p.isIdentLike() {
			action += " " + p.cur.Str
			p.advance()
		}
		stmt.Action = action
	} else if p.cur.Type == kwDELETE {
		p.advance() // consume DELETE
		action := "DELETE"
		if p.isIdentLike() {
			action += " " + p.cur.Str
			p.advance()
		}
		stmt.Action = action
	} else if p.cur.Type == kwVALIDATE {
		p.advance() // consume VALIDATE
		action := "VALIDATE"
		if p.isIdentLike() {
			action += " " + p.cur.Str
			p.advance()
		}
		stmt.Action = action
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseExplainPlanStmt parses an EXPLAIN PLAN statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/EXPLAIN-PLAN.html
//
//	EXPLAIN PLAN
//	    [SET STATEMENT_ID = 'id']
//	    [INTO [schema.]table[@dblink]]
//	    FOR statement
func (p *Parser) parseExplainPlanStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume EXPLAIN

	// Expect PLAN
	if p.cur.Type == kwPLAN {
		p.advance()
	}

	stmt := &nodes.ExplainPlanStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Optional SET STATEMENT_ID = 'id'
	if p.cur.Type == kwSET {
		p.advance() // consume SET
		// STATEMENT_ID is an identifier
		if p.isIdentLike() && p.cur.Str == "STATEMENT_ID" {
			p.advance() // consume STATEMENT_ID
			if p.cur.Type == '=' {
				p.advance() // consume =
			}
			if p.cur.Type == tokSCONST {
				stmt.StatementID = p.cur.Str
				p.advance()
			}
		}
	}

	// Optional INTO [schema.]table
	if p.cur.Type == kwINTO {
		p.advance()
		stmt.Into = p.parseObjectName()
	}

	// FOR statement
	if p.cur.Type == kwFOR {
		p.advance()
		stmt.Statement = p.parseStmt()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseFlashbackTableStmt parses a FLASHBACK TABLE statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/FLASHBACK-TABLE.html
//
//	FLASHBACK TABLE [schema.]table TO
//	    { SCN expr | TIMESTAMP expr | BEFORE DROP [RENAME TO name] }
func (p *Parser) parseFlashbackTableStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume FLASHBACK

	// Expect TABLE
	if p.cur.Type == kwTABLE {
		p.advance()
	}

	stmt := &nodes.FlashbackTableStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Table name
	stmt.Table = p.parseObjectName()

	// TO
	if p.cur.Type == kwTO {
		p.advance()
	}

	// SCN expr | TIMESTAMP expr | BEFORE DROP
	switch p.cur.Type {
	case kwSCN:
		p.advance()
		stmt.ToSCN = p.parseExpr()
	case kwTIMESTAMP:
		p.advance()
		stmt.ToTimestamp = p.parseExpr()
	case kwBEFORE:
		p.advance() // consume BEFORE
		if p.cur.Type == kwDROP {
			p.advance() // consume DROP
			stmt.ToBeforeDrop = true
		}
		// Optional RENAME TO name
		if p.cur.Type == kwRENAME {
			p.advance() // consume RENAME
			if p.cur.Type == kwTO {
				p.advance()
			}
			stmt.Rename = p.parseIdentifier()
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseFlashbackDatabaseStmt parses a FLASHBACK DATABASE statement.
//
// BNF: oracle/parser/bnf/FLASHBACK-DATABASE.bnf
//
//	FLASHBACK [ STANDBY | PLUGGABLE ] DATABASE [ database ]
//	    { TO SCN scn_number
//	    | TO BEFORE SCN scn_number
//	    | TO TIMESTAMP timestamp_expression
//	    | TO BEFORE TIMESTAMP timestamp_expression
//	    | TO RESTORE POINT restore_point_name
//	    | TO BEFORE RESETLOGS
//	    } ;
func (p *Parser) parseFlashbackDatabaseStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume FLASHBACK

	stmt := &nodes.FlashbackDatabaseStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Optional STANDBY | PLUGGABLE
	if p.isIdentLike() && p.cur.Str == "STANDBY" {
		stmt.Modifier = "STANDBY"
		p.advance()
	} else if p.isIdentLike() && p.cur.Str == "PLUGGABLE" {
		stmt.Modifier = "PLUGGABLE"
		p.advance()
	}

	// DATABASE
	if p.cur.Type == kwDATABASE {
		p.advance()
	}

	// Optional database name (not TO keyword)
	if p.isIdentLike() && p.cur.Str != "TO" {
		stmt.DatabaseName = p.parseObjectName()
	}

	// TO
	if p.cur.Type == kwTO {
		p.advance()
	}

	// Optional BEFORE
	if p.isIdentLike() && p.cur.Str == "BEFORE" {
		stmt.Before = true
		p.advance()
	}

	switch p.cur.Type {
	case kwSCN:
		p.advance()
		stmt.ToSCN = p.parseExpr()
	case kwTIMESTAMP:
		p.advance()
		stmt.ToTimestamp = p.parseExpr()
	default:
		if p.isIdentLike() && p.cur.Str == "RESTORE" {
			p.advance()
			if p.isIdentLike() && p.cur.Str == "POINT" {
				p.advance()
			}
			if p.isIdentLike() {
				stmt.ToRestorePoint = p.cur.Str
				p.advance()
			}
		} else if p.isIdentLike() && p.cur.Str == "RESETLOGS" {
			stmt.ToResetlogs = true
			p.advance()
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parsePurgeStmt parses a PURGE statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/PURGE.html
//
//	PURGE { TABLE name | INDEX name | RECYCLEBIN | DBA_RECYCLEBIN | TABLESPACE name }
func (p *Parser) parsePurgeStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume PURGE

	stmt := &nodes.PurgeStmt{
		Loc: nodes.Loc{Start: start},
	}

	switch p.cur.Type {
	case kwTABLE:
		stmt.ObjectType = nodes.OBJECT_TABLE
		p.advance()
		stmt.Name = p.parseObjectName()
	case kwINDEX:
		stmt.ObjectType = nodes.OBJECT_INDEX
		p.advance()
		stmt.Name = p.parseObjectName()
	case kwTABLESPACE:
		stmt.ObjectType = nodes.OBJECT_TABLESPACE
		p.advance()
		stmt.Name = p.parseObjectName()
	default:
		// RECYCLEBIN or DBA_RECYCLEBIN (parsed as identifiers)
		if p.isIdentLike() {
			ident := p.cur.Str
			p.advance()
			stmt.Name = &nodes.ObjectName{
				Name: ident,
				Loc:  nodes.Loc{Start: start, End: p.pos()},
			}
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseAuditStmt parses an AUDIT statement (Traditional + Unified Auditing).
//
// BNF: oracle/parser/bnf/AUDIT-Traditional-Auditing.bnf
//
//	AUDIT { audit_operation_clause [ auditing_by_clause | auditing_on_clause ]
//	      | audit_schema_object_clause
//	      }
//	    [ BY { SESSION | ACCESS } ]
//	    [ WHENEVER [ NOT ] SUCCESSFUL ]
//	    [ CONTAINER = { CURRENT | ALL } ]
//
//	audit_operation_clause:
//	    { sql_statement_shortcut
//	    | system_privilege
//	    | ALL
//	    | ALL STATEMENTS
//	    | ALL PRIVILEGES
//	    }
//	    [, { sql_statement_shortcut | system_privilege } ]...
//
//	auditing_by_clause:
//	    BY { user [, user ]...
//	       | SESSION CURRENT
//	       }
//
//	audit_schema_object_clause:
//	    sql_operation [, sql_operation ]...
//	    auditing_on_clause
//
//	auditing_on_clause:
//	    ON [ schema. ] object
//	  | ON DEFAULT
//	  | ON DIRECTORY directory_name
//	  | ON MINING MODEL model_name
//	  | ON SQL TRANSLATION PROFILE profile_name
//	  | NETWORK
//	  | DIRECT_PATH LOAD
//
// BNF: oracle/parser/bnf/AUDIT-Unified-Auditing.bnf
//
//	AUDIT
//	    { POLICY policy_name
//	        [ { BY | EXCEPT } user [, user ]... ]
//	        [ { BY | EXCEPT } USERS WITH ROLE role [, role ]... ]
//	        [ WHENEVER [ NOT ] SUCCESSFUL ]
//	    | CONTEXT NAMESPACE namespace
//	        ATTRIBUTES attribute [, attribute ]...
//	        [ BY user [, user ]... ]
//	    } ;
func (p *Parser) parseAuditStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume AUDIT

	stmt := &nodes.AuditStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Unified: AUDIT POLICY ...
	if p.isIdentLikeStr("POLICY") {
		p.advance()
		if p.isIdentLike() {
			stmt.Policy = p.cur.Str
			p.advance()
		}
		// [ { BY | EXCEPT } user [, user]... ]
		// [ { BY | EXCEPT } USERS WITH ROLE role [, role]... ]
		for p.cur.Type != ';' && p.cur.Type != tokEOF {
			if p.cur.Type == kwBY {
				p.advance()
				if p.isIdentLikeStr("USERS") {
					// BY USERS WITH ROLE role [, role]...
					p.advance()
					if p.cur.Type == kwWITH {
						p.advance()
					}
					if p.cur.Type == kwROLE {
						p.advance()
					}
					stmt.WithRoles = p.parseIdentListForAudit()
				} else {
					// BY user [, user]...
					stmt.ByUsers = p.parseIdentListForAudit()
				}
			} else if p.cur.Type == kwEXCEPT {
				p.advance()
				if p.isIdentLikeStr("USERS") {
					// EXCEPT USERS WITH ROLE role [, role]...
					p.advance()
					if p.cur.Type == kwWITH {
						p.advance()
					}
					if p.cur.Type == kwROLE {
						p.advance()
					}
					stmt.WithRoles = p.parseIdentListForAudit()
					stmt.WithRoleExcept = true
				} else {
					// EXCEPT user [, user]...
					stmt.ExceptUsers = p.parseIdentListForAudit()
				}
			} else if p.cur.Type == kwWHENEVER {
				stmt.When = p.parseWheneverClause()
			} else {
				break
			}
		}
		stmt.Loc.End = p.pos()
		return stmt
	}

	// Unified: AUDIT CONTEXT NAMESPACE ...
	if p.cur.Type == kwCONTEXT {
		p.advance()
		if p.isIdentLikeStr("NAMESPACE") {
			p.advance()
		}
		if p.isIdentLike() {
			stmt.ContextNS = p.cur.Str
			p.advance()
		}
		if p.isIdentLikeStr("ATTRIBUTES") {
			p.advance()
			stmt.ContextAttrs = p.parseIdentListForAudit()
		}
		if p.cur.Type == kwBY {
			p.advance()
			stmt.ByUsers = p.parseIdentListForAudit()
		}
		stmt.Loc.End = p.pos()
		return stmt
	}

	// Traditional: parse audit actions
	stmt.Actions = p.parseAuditActions()

	// ON clause or auditing_by_clause
	if p.cur.Type == kwON {
		p.advance()
		p.parseTraditionalAuditOnClause(stmt)
	} else if p.isIdentLikeStr("NETWORK") {
		stmt.OnNetwork = true
		p.advance()
	} else if p.isIdentLikeStr("DIRECT_PATH") {
		stmt.OnDirectPath = true
		p.advance()
		if p.isIdentLikeStr("LOAD") {
			p.advance()
		}
	}

	// BY { SESSION | ACCESS } or BY user [, user]...
	if p.cur.Type == kwBY {
		p.advance()
		tok := p.cur.Str
		if tok == "SESSION" || tok == "ACCESS" {
			stmt.By = tok
			p.advance()
			// Check for CURRENT after SESSION
			if tok == "SESSION" && p.isIdentLikeStr("CURRENT") {
				stmt.By = "SESSION CURRENT"
				p.advance()
			}
		} else {
			// BY user [, user]...
			stmt.ByUsers2 = p.parseIdentListForAudit()
		}
	}

	// WHENEVER [NOT] SUCCESSFUL
	if p.cur.Type == kwWHENEVER {
		stmt.When = p.parseWheneverClause()
	}

	// CONTAINER = { CURRENT | ALL }
	if p.isIdentLikeStr("CONTAINER") {
		p.advance()
		stmt.ContainerAll = p.parseContainerClause()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseNoauditStmt parses a NOAUDIT statement (Traditional + Unified Auditing).
//
// BNF: oracle/parser/bnf/NOAUDIT-Traditional-Auditing.bnf
//
//	NOAUDIT
//	    { audit_operation_clause [ auditing_by_clause ]
//	    | audit_schema_object_clause
//	    }
//	    [ WHENEVER [ NOT ] SUCCESSFUL ]
//	    [ CONTAINER = { CURRENT | ALL } ] ;
//
//	audit_operation_clause::=
//	    { statement_option [, statement_option ]...
//	    | ALL
//	    | ALL STATEMENTS
//	    | system_privilege [, system_privilege ]...
//	    | ALL PRIVILEGES
//	    }
//
//	auditing_by_clause::=
//	    BY user [, user ]...
//
//	audit_schema_object_clause::=
//	    { sql_operation [, sql_operation ]...
//	    | ALL
//	    }
//	    ON auditing_on_clause
//
//	auditing_on_clause::=
//	    [ schema. ] object
//	  | DIRECTORY directory_name
//	  | SQL TRANSLATION PROFILE [ schema. ] profile
//	  | DEFAULT
//	  | NETWORK
//	  | DIRECT_PATH LOAD
//
// BNF: oracle/parser/bnf/NOAUDIT-Unified-Auditing.bnf
//
//	NOAUDIT POLICY policy
//	    [ BY user [, user ]...
//	      [ WITH ROLE role [, role ]... ]
//	    ] ;
//
//	NOAUDIT CONTEXT NAMESPACE namespace
//	    ATTRIBUTES attribute [, attribute ]...
//	    [ BY user [, user ]...
//	      [ WITH ROLE role [, role ]... ]
//	    ] ;
func (p *Parser) parseNoauditStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume NOAUDIT

	stmt := &nodes.NoauditStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Unified: NOAUDIT POLICY ...
	if p.isIdentLikeStr("POLICY") {
		p.advance()
		if p.isIdentLike() {
			stmt.Policy = p.cur.Str
			p.advance()
		}
		if p.cur.Type == kwBY {
			p.advance()
			stmt.ByUsers = p.parseIdentListForAudit()
			if p.cur.Type == kwWITH {
				p.advance()
				if p.cur.Type == kwROLE {
					p.advance()
				}
				stmt.WithRoles = p.parseIdentListForAudit()
			}
		}
		stmt.Loc.End = p.pos()
		return stmt
	}

	// Unified: NOAUDIT CONTEXT NAMESPACE ...
	if p.cur.Type == kwCONTEXT {
		p.advance()
		if p.isIdentLikeStr("NAMESPACE") {
			p.advance()
		}
		if p.isIdentLike() {
			stmt.ContextNS = p.cur.Str
			p.advance()
		}
		if p.isIdentLikeStr("ATTRIBUTES") {
			p.advance()
			stmt.ContextAttrs = p.parseIdentListForAudit()
		}
		if p.cur.Type == kwBY {
			p.advance()
			stmt.ByUsers = p.parseIdentListForAudit()
			if p.cur.Type == kwWITH {
				p.advance()
				if p.cur.Type == kwROLE {
					p.advance()
				}
				stmt.WithRoles = p.parseIdentListForAudit()
			}
		}
		stmt.Loc.End = p.pos()
		return stmt
	}

	// Traditional: parse audit actions
	stmt.Actions = p.parseAuditActions()

	// ON clause
	if p.cur.Type == kwON {
		p.advance()
		p.parseTraditionalNoauditOnClause(stmt)
	} else if p.isIdentLikeStr("NETWORK") {
		stmt.OnNetwork = true
		p.advance()
	} else if p.isIdentLikeStr("DIRECT_PATH") {
		stmt.OnDirectPath = true
		p.advance()
		if p.isIdentLikeStr("LOAD") {
			p.advance()
		}
	}

	// BY user [, user]...
	if p.cur.Type == kwBY {
		p.advance()
		stmt.ByUsers2 = p.parseIdentListForAudit()
	}

	// WHENEVER [NOT] SUCCESSFUL
	if p.cur.Type == kwWHENEVER {
		stmt.When = p.parseWheneverClause()
	}

	// CONTAINER = { CURRENT | ALL }
	if p.isIdentLikeStr("CONTAINER") {
		p.advance()
		stmt.ContainerAll = p.parseContainerClause()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseTraditionalAuditOnClause parses the ON clause for traditional AUDIT.
// Called after ON has been consumed.
func (p *Parser) parseTraditionalAuditOnClause(stmt *nodes.AuditStmt) {
	if p.cur.Type == kwDEFAULT {
		stmt.OnDefault = true
		p.advance()
	} else if p.isIdentLikeStr("DIRECTORY") {
		p.advance()
		if p.isIdentLike() {
			stmt.OnDirectory = p.cur.Str
			p.advance()
		}
	} else if p.isIdentLikeStr("MINING") {
		p.advance()
		if p.cur.Type == kwMODEL {
			p.advance()
		}
		stmt.Object = p.parseObjectName()
	} else if p.isIdentLikeStr("SQL") {
		// SQL TRANSLATION PROFILE
		p.advance()
		if p.isIdentLikeStr("TRANSLATION") {
			p.advance()
		}
		if p.cur.Type == kwPROFILE {
			p.advance()
		}
		stmt.Object = p.parseObjectName()
	} else {
		stmt.Object = p.parseObjectName()
	}
}

// parseTraditionalNoauditOnClause parses the ON clause for traditional NOAUDIT.
// Called after ON has been consumed.
func (p *Parser) parseTraditionalNoauditOnClause(stmt *nodes.NoauditStmt) {
	if p.cur.Type == kwDEFAULT {
		stmt.OnDefault = true
		p.advance()
	} else if p.isIdentLikeStr("DIRECTORY") {
		p.advance()
		if p.isIdentLike() {
			stmt.OnDirectory = p.cur.Str
			p.advance()
		}
	} else if p.isIdentLikeStr("SQL") {
		// SQL TRANSLATION PROFILE
		p.advance()
		if p.isIdentLikeStr("TRANSLATION") {
			p.advance()
		}
		if p.cur.Type == kwPROFILE {
			p.advance()
		}
		stmt.Object = p.parseObjectName()
	} else if p.isIdentLikeStr("NETWORK") {
		stmt.OnNetwork = true
		p.advance()
	} else if p.isIdentLikeStr("DIRECT_PATH") {
		stmt.OnDirectPath = true
		p.advance()
		if p.isIdentLikeStr("LOAD") {
			p.advance()
		}
	} else {
		stmt.Object = p.parseObjectName()
	}
}

// parseAuditActions collects audit action identifiers separated by commas.
func (p *Parser) parseAuditActions() []string {
	var actions []string
	for {
		// Collect multi-word action (e.g., "CREATE TABLE", "ALTER SESSION")
		action := ""
		for p.isIdentLike() || p.cur.Type == kwSELECT || p.cur.Type == kwINSERT ||
			p.cur.Type == kwUPDATE || p.cur.Type == kwDELETE || p.cur.Type == kwCREATE ||
			p.cur.Type == kwALTER || p.cur.Type == kwDROP || p.cur.Type == kwGRANT ||
			p.cur.Type == kwEXECUTE || p.cur.Type == kwINDEX || p.cur.Type == kwALL {
			// Stop before special clause identifiers
			if p.isIdentLikeStr("NETWORK") || p.isIdentLikeStr("DIRECT_PATH") || p.isIdentLikeStr("CONTAINER") {
				break
			}
			if action != "" {
				action += " "
			}
			action += p.cur.Str
			p.advance()
			// Stop if we hit keywords that start a clause
			if p.cur.Type == kwON || p.cur.Type == kwBY || p.cur.Type == kwWHENEVER {
				break
			}
		}
		if action != "" {
			actions = append(actions, action)
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return actions
}

// parseWheneverClause parses WHENEVER [NOT] SUCCESSFUL.
func (p *Parser) parseWheneverClause() string {
	result := "WHENEVER"
	p.advance() // consume WHENEVER
	if p.cur.Type == kwNOT {
		result += " NOT"
		p.advance()
	}
	if p.cur.Type == kwSUCCESSFUL {
		result += " SUCCESSFUL"
		p.advance()
	}
	return result
}

// parseIdentListForAudit parses a comma-separated list of identifiers,
// stopping at clause keywords (BY, EXCEPT, WHENEVER, WITH, CONTAINER, ;, EOF).
func (p *Parser) parseIdentListForAudit() []string {
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

// parseAssociateStatisticsStmt parses an ASSOCIATE STATISTICS statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/ASSOCIATE-STATISTICS.html
//
//	ASSOCIATE STATISTICS WITH { COLUMNS | FUNCTIONS | PACKAGES | TYPES | INDEXES }
//	    object [,...] USING [schema.]statistics_type
func (p *Parser) parseAssociateStatisticsStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume ASSOCIATE
	if p.cur.Type == kwSTATISTICS {
		p.advance()
	}

	stmt := &nodes.AssociateStatisticsStmt{
		Loc: nodes.Loc{Start: start},
	}

	// WITH
	if p.cur.Type == kwWITH {
		p.advance()
	}

	// Object type: COLUMNS, FUNCTIONS, PACKAGES, TYPES, INDEXES
	if p.isIdentLike() || p.cur.Type == kwINDEX {
		if p.cur.Type == kwINDEX {
			stmt.ObjectType = "INDEXES"
		} else {
			stmt.ObjectType = p.cur.Str
		}
		p.advance()
	}

	// Object names
	for {
		name := p.parseObjectName()
		if name != nil {
			stmt.Objects = append(stmt.Objects, name)
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	// USING statistics_type
	if p.isIdentLike() && p.cur.Str == "USING" {
		p.advance()
		stmt.Using = p.parseObjectName()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseDisassociateStatisticsStmt parses a DISASSOCIATE STATISTICS statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/DISASSOCIATE-STATISTICS.html
//
//	DISASSOCIATE STATISTICS FROM { COLUMNS | FUNCTIONS | PACKAGES | TYPES | INDEXES }
//	    object [,...] [FORCE]
func (p *Parser) parseDisassociateStatisticsStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume DISASSOCIATE
	if p.cur.Type == kwSTATISTICS {
		p.advance()
	}

	stmt := &nodes.DisassociateStatisticsStmt{
		Loc: nodes.Loc{Start: start},
	}

	// FROM
	if p.cur.Type == kwFROM {
		p.advance()
	}

	// Object type
	if p.isIdentLike() || p.cur.Type == kwINDEX {
		if p.cur.Type == kwINDEX {
			stmt.ObjectType = "INDEXES"
		} else {
			stmt.ObjectType = p.cur.Str
		}
		p.advance()
	}

	// Object names
	for {
		name := p.parseObjectName()
		if name != nil {
			stmt.Objects = append(stmt.Objects, name)
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	// FORCE
	if p.cur.Type == kwFORCE {
		stmt.Force = true
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt
}
