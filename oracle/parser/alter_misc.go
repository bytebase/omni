package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseAlterStmt dispatches ALTER statements based on the next keyword.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/SQL-Statements-ALTER-ANALYTIC-VIEW-to-ALTER-SYSTEM.html
//
//	ALTER SESSION SET param = value [ param = value ... ]
//	ALTER SYSTEM  SET param = value [ param = value ... ]
//	ALTER SYSTEM  KILL SESSION 'sid,serial#'
//	ALTER INDEX   name ...
//	ALTER VIEW    name ...
//	ALTER SEQUENCE name ...
func (p *Parser) parseAlterStmt() (nodes.StmtNode, error) {
	start := p.pos()
	p.advance() // consume ALTER

	switch p.cur.Type {
	case kwSESSION:
		return p.parseAlterSessionStmt(start)
	case kwSYSTEM:
		return p.parseAlterSystemStmt(start)
	case kwINDEX:
		return p.parseAlterIndexStmt(start)
	case kwVIEW:
		return p.parseAlterViewStmt(start)
	case kwSEQUENCE:
		return p.parseAlterSequenceStmt(start)
	case kwTABLE:
		return p.parseAlterTableStmt(start)
	case kwPROCEDURE:
		return p.parseAlterProcedureStmt(start)
	case kwFUNCTION:
		return p.parseAlterFunctionStmt(start)
	case kwTRIGGER:
		return p.parseAlterTriggerStmt(start)
	case kwTYPE:
		return p.parseAlterTypeStmt(start)
	case kwPACKAGE:
		return p.parseAlterPackageStmt(start)
	case kwMATERIALIZED:
		// Check for MATERIALIZED ZONEMAP vs MATERIALIZED VIEW
		p.advance() // consume MATERIALIZED
		if p.isIdentLike() && p.cur.Str == "ZONEMAP" {
			p.advance() // consume ZONEMAP
			return p.parseAlterMaterializedZonemapStmt(start)
		}
		// MATERIALIZED VIEW - consume VIEW and check if LOG follows
		if p.cur.Type == kwVIEW {
			p.advance() // consume VIEW
		}
		// Check for MATERIALIZED VIEW LOG
		if p.cur.Type == kwLOG {
			p.advance() // consume LOG
			return p.parseAlterMviewLogStmt(start)
		}
		return p.parseAlterMaterializedViewStmt(start)
	case kwDATABASE:
		// Distinguish ALTER DATABASE LINK, ALTER DATABASE DICTIONARY, ALTER DATABASE
		next := p.peekNext()
		if next.Type == kwLINK {
			return p.parseAlterDatabaseLinkStmt(start, false, false)
		}
		p.advance() // consume DATABASE
		if p.isIdentLikeStr("DICTIONARY") {
			p.advance() // consume DICTIONARY
			return p.parseAlterDatabaseDictionaryStmt(start)
		}
		return p.parseAlterDatabaseStmt(start)
	case kwSYNONYM:
		return p.parseAlterSynonymStmt(start, false)
	case kwPUBLIC:
		// ALTER PUBLIC DATABASE LINK or ALTER PUBLIC SYNONYM
		p.advance() // consume PUBLIC
		if p.cur.Type == kwDATABASE {
			return p.parseAlterDatabaseLinkStmt(start, false, true)
		}
		if p.cur.Type == kwSYNONYM {
			return p.parseAlterSynonymStmt(start, true)
		}
		// Unknown ALTER PUBLIC target
		p.skipToSemicolon()
		return nil, nil
	case kwAUDIT:
		// ALTER AUDIT POLICY
		p.advance() // consume AUDIT
		if p.isIdentLikeStr("POLICY") {
			p.advance() // consume POLICY
		}
		return p.parseAlterAuditPolicyStmt(start)
	case kwJSON:
		// ALTER JSON RELATIONAL DUALITY VIEW
		p.advance() // consume JSON
		if p.isIdentLike() && p.cur.Str == "RELATIONAL" {
			p.advance() // consume RELATIONAL
		}
		if p.isIdentLike() && p.cur.Str == "DUALITY" {
			p.advance() // consume DUALITY
		}
		if p.cur.Type == kwVIEW {
			p.advance() // consume VIEW
		}
		return p.parseAlterJsonDualityViewStmt(start)
	case kwFLASHBACK:
		// ALTER FLASHBACK ARCHIVE
		p.advance() // consume FLASHBACK
		if p.isIdentLike() && p.cur.Str == "ARCHIVE" {
			p.advance() // consume ARCHIVE
		}
		return p.parseAlterFlashbackArchiveStmt(start)
	case kwUSER, kwROLE, kwPROFILE,
		kwTABLESPACE, kwCLUSTER, kwJAVA, kwLIBRARY:
		adminStmt, err := p.parseAlterAdminObject(start)
		if err != nil {
			return nil, err
		}
		if adminStmt != nil {
			return adminStmt, nil
		}
		p.skipToSemicolon()
		return nil, nil
	default:
		if p.isIdentLike() {
			// ALTER SHARED [PUBLIC] DATABASE LINK
			if p.cur.Str == "SHARED" {
				p.advance() // consume SHARED
				isPublic := false
				if p.cur.Type == kwPUBLIC {
					isPublic = true
					p.advance() // consume PUBLIC
				}
				if p.cur.Type == kwDATABASE {
					return p.parseAlterDatabaseLinkStmt(start, true, isPublic)
				}
				if p.cur.Type == kwSYNONYM {
					return p.parseAlterSynonymStmt(start, isPublic)
				}
				// Unknown ALTER SHARED target
				p.skipToSemicolon()
				return nil, nil
			}
			// ALTER USECASE DOMAIN
			if p.cur.Str == "USECASE" {
				p.advance() // consume USECASE
				if p.isIdentLike() && p.cur.Str == "DOMAIN" {
					p.advance() // consume DOMAIN
					return p.parseAlterDomainStmt(start, true)
				}
			}
			// Check for DIMENSION and other identifier-based objects
			adminStmt, err := p.parseAlterAdminObject(start)
			if err != nil {
				return nil, err
			}
			if adminStmt != nil {
				return adminStmt, nil
			}
		}
		// Unknown ALTER target — skip to semicolon or EOF.
		p.skipToSemicolon()
		return nil, nil
	}
}

// parseAlterSessionStmt parses ALTER SESSION SET param = value [, ...].
// parseAlterSessionStmt parses an ALTER SESSION statement.
// Called after ALTER has been consumed.
//
// BNF: oracle/parser/bnf/ALTER-SESSION.bnf
//
//	ALTER SESSION
//	    { ADVISE { COMMIT | ROLLBACK | NOTHING }
//	    | CLOSE DATABASE LINK dblink
//	    | { ENABLE | DISABLE } COMMIT IN PROCEDURE
//	    | { ENABLE | DISABLE } GUARD
//	    | { ENABLE | DISABLE } PARALLEL { DML | DDL | QUERY }
//	    | FORCE PARALLEL { DML | DDL | QUERY } [ PARALLEL integer ]
//	    | { ENABLE | DISABLE } RESUMABLE [ TIMEOUT integer ] [ NAME 'string' ]
//	    | { ENABLE | DISABLE } SHARD DDL
//	    | SYNC WITH PRIMARY
//	    | alter_session_set_clause
//	    }
//
//	alter_session_set_clause:
//	    SET { parameter_name = parameter_value [, parameter_name = parameter_value ]...
//	        | EDITION = edition_name
//	        | CONTAINER = container_name [ SERVICE = service_name ]
//	        | ROW ARCHIVAL VISIBILITY = { ACTIVE | ALL }
//	        | DEFAULT_COLLATION = { collation_name | NONE }
//	        | CONSTRAINT[S] = { IMMEDIATE | DEFERRED | DEFAULT }
//	        | CURRENT_SCHEMA = schema
//	        | ERROR_ON_OVERLAP_TIME = { TRUE | FALSE }
//	        | FLAGGER = { ENTRY | OFF }
//	        | INSTANCE = integer
//	        | ISOLATION_LEVEL = { SERIALIZABLE | READ COMMITTED }
//	        | STANDBY_MAX_DATA_DELAY = { integer | NONE }
//	        | TIME_ZONE = { '{ + | - } hh:mi' | LOCAL | DBTIMEZONE | 'time_zone_region' }
//	        | USE_PRIVATE_OUTLINES = { TRUE | FALSE | category_name }
//	        | USE_STORED_OUTLINES = { TRUE | FALSE | category_name }
//	        }
func (p *Parser) parseAlterSessionStmt(start int) (nodes.StmtNode, error) {
	p.advance() // consume SESSION

	stmt := &nodes.AlterSessionStmt{
		Loc: nodes.Loc{Start: start},
	}

	switch {
	case p.isIdentLikeStr("ADVISE"):
		p.advance() // consume ADVISE
		stmt.Action = "ADVISE"
		if p.cur.Type == kwCOMMIT {
			stmt.AdviseAction = "COMMIT"
			p.advance()
		} else if p.isIdentLikeStr("ROLLBACK") {
			stmt.AdviseAction = "ROLLBACK"
			p.advance()
		} else if p.isIdentLikeStr("NOTHING") {
			stmt.AdviseAction = "NOTHING"
			p.advance()
		}

	case p.isIdentLikeStr("CLOSE"):
		p.advance() // consume CLOSE
		stmt.Action = "CLOSE_DATABASE_LINK"
		// DATABASE
		if p.cur.Type == kwDATABASE {
			p.advance()
		}
		// LINK
		if p.isIdentLikeStr("LINK") {
			p.advance()
		}
		var parseErr1 error
		stmt.DBLink, parseErr1 = p.parseIdentifier()
		if parseErr1 != nil {
			return nil, parseErr1
		}

	case p.cur.Type == kwENABLE || p.cur.Type == kwDISABLE:
		action := "ENABLE"
		if p.cur.Type == kwDISABLE {
			action = "DISABLE"
		}
		stmt.Action = action
		p.advance() // consume ENABLE/DISABLE

		switch {
		case p.cur.Type == kwCOMMIT:
			// COMMIT IN PROCEDURE
			stmt.Feature = "COMMIT_IN_PROCEDURE"
			p.advance() // consume COMMIT
			if p.cur.Type == kwIN {
				p.advance() // consume IN
			}
			if p.isIdentLikeStr("PROCEDURE") {
				p.advance()
			}
		case p.isIdentLikeStr("GUARD"):
			stmt.Feature = "GUARD"
			p.advance()
		case p.cur.Type == kwPARALLEL:
			p.advance() // consume PARALLEL
			if p.isIdentLikeStr("DML") {
				stmt.Feature = "PARALLEL_DML"
				p.advance()
			} else if p.isIdentLikeStr("DDL") {
				stmt.Feature = "PARALLEL_DDL"
				p.advance()
			} else if p.isIdentLikeStr("QUERY") {
				stmt.Feature = "PARALLEL_QUERY"
				p.advance()
			}
		case p.isIdentLikeStr("RESUMABLE"):
			stmt.Feature = "RESUMABLE"
			p.advance()
			// optional TIMEOUT integer
			if p.isIdentLikeStr("TIMEOUT") {
				p.advance()
				if p.cur.Type == tokICONST {
					var parseErr2 error
					stmt.Timeout, parseErr2 = p.parseIntValue()
					if parseErr2 != nil {

						// optional NAME 'string'
						return nil, parseErr2
					}
				}
			}

			if p.cur.Type == kwNAME {
				p.advance()
				if p.cur.Type == tokSCONST {
					stmt.ResumableName = p.cur.Str
					p.advance()
				}
			}
		case p.isIdentLikeStr("SHARD"):
			stmt.Feature = "SHARD_DDL"
			p.advance() // consume SHARD
			if p.isIdentLikeStr("DDL") {
				p.advance()
			}
		}

	case p.cur.Type == kwFORCE:
		p.advance() // consume FORCE
		stmt.Action = "FORCE_PARALLEL"
		// PARALLEL
		if p.cur.Type == kwPARALLEL {
			p.advance()
		}
		// {DML | DDL | QUERY}
		if p.isIdentLikeStr("DML") {
			stmt.Feature = "PARALLEL_DML"
			p.advance()
		} else if p.isIdentLikeStr("DDL") {
			stmt.Feature = "PARALLEL_DDL"
			p.advance()
		} else if p.isIdentLikeStr("QUERY") {
			stmt.Feature = "PARALLEL_QUERY"
			p.advance()
		}
		// optional PARALLEL integer
		if p.cur.Type == kwPARALLEL {
			p.advance()
			if p.cur.Type == tokICONST {
				var parseErr3 error
				stmt.ParallelDegree, parseErr3 = p.parseIntValue()
				if parseErr3 != nil {
					return nil, parseErr3
				}
			}
		}

	case p.isIdentLikeStr("SYNC"):
		p.advance() // consume SYNC
		stmt.Action = "SYNC_WITH_PRIMARY"
		if p.cur.Type == kwWITH {
			p.advance()
		}
		if p.isIdentLikeStr("PRIMARY") {
			p.advance()
		}

	case p.cur.Type == kwSET:
		p.advance() // consume SET
		stmt.Action = "SET"
		var parseErr4 error
		stmt.SetParams, parseErr4 = p.parseSetParams()
		if parseErr4 != nil {
			return nil, parseErr4
		}
	}
	if stmt.Action == "" {
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterSystemStmt parses an ALTER SYSTEM statement.
// Called after ALTER has been consumed.
//
// BNF: oracle/parser/bnf/ALTER-SYSTEM.bnf
//
//	ALTER SYSTEM
//	    { archive_log_clause
//	    | checkpoint_clause
//	    | check_datafiles_clause
//	    | distributed_recov_clauses
//	    | SWITCH LOGFILE
//	    | { SUSPEND | RESUME }
//	    | quiesce_clauses
//	    | rolling_migration_clauses
//	    | rolling_patch_clauses
//	    | security_clauses
//	    | shutdown_dispatcher_clause
//	    | REGISTER
//	    | alter_system_set_clause
//	    | alter_system_reset_clause
//	    | cancel_sql_clause
//	    | flush_clause
//	    | RELOCATE CLIENT 'client_id'
//	    | end_session_clauses
//	    }
//
//	archive_log_clause:
//	    ARCHIVE LOG [ INSTANCE 'instance_name' ]
//	        { SEQUENCE integer
//	        | CHANGE integer
//	        | CURRENT [ NOSWITCH ]
//	        | GROUP integer
//	        | LOGFILE 'filename' [ USING BACKUP CONTROLFILE ]
//	        | NEXT
//	        | ALL
//	        }
//	        [ THREAD integer ]
//	        [ TO 'location' ]
//
//	checkpoint_clause:
//	    CHECKPOINT [ { GLOBAL | LOCAL } ]
//
//	check_datafiles_clause:
//	    CHECK DATAFILES [ { GLOBAL | LOCAL } ]
//
//	distributed_recov_clauses:
//	    { ENABLE | DISABLE } DISTRIBUTED RECOVERY
//
//	end_session_clauses:
//	    { DISCONNECT SESSION 'session_id, serial_number'
//	        [ POST_TRANSACTION ] [ IMMEDIATE ]
//	    | KILL SESSION 'session_id, serial_number' [ , @instance_id ]
//	        [ IMMEDIATE | FORCE ]
//	        [ NOREPLAY ]
//	        [ TIMEOUT integer ]
//	    }
//
//	quiesce_clauses:
//	    { QUIESCE RESTRICTED | UNQUIESCE }
//
//	rolling_migration_clauses:
//	    { START ROLLING MIGRATION TO 'ASM_version'
//	    | STOP ROLLING MIGRATION
//	    }
//
//	rolling_patch_clauses:
//	    { START ROLLING PATCH | STOP ROLLING PATCH }
//
//	security_clauses:
//	    { ENABLE RESTRICTED SESSION
//	    | DISABLE RESTRICTED SESSION
//	    | SET ENCRYPTION WALLET OPEN IDENTIFIED BY password
//	    | SET ENCRYPTION WALLET CLOSE [ IDENTIFIED BY password ]
//	    | SET ENCRYPTION KEY [ IDENTIFIED BY password ]
//	    }
//
//	shutdown_dispatcher_clause:
//	    SHUTDOWN [ IMMEDIATE ] 'dispatcher_name'
//
//	alter_system_set_clause:
//	    SET parameter_name = parameter_value [, parameter_value ]...
//	        [ COMMENT = 'comment' ]
//	        [ DEFERRED ]
//	        [ SCOPE = { MEMORY | SPFILE | BOTH } ]
//	        [ SID = { 'sid' | '*' } ]
//	        [ CONTAINER = { ALL | CURRENT } ]
//
//	alter_system_reset_clause:
//	    RESET parameter_name
//	        [ SCOPE = { MEMORY | SPFILE | BOTH } ]
//	        [ SID = { 'sid' | '*' } ]
//
//	cancel_sql_clause:
//	    CANCEL SQL 'session_id, serial_number' [ , @instance_id ] [ SQL_ID 'sql_id' ]
//
//	flush_clause:
//	    FLUSH
//	        { SHARED_POOL
//	        | GLOBAL CONTEXT
//	        | BUFFER_CACHE [ { GLOBAL | LOCAL } ]
//	        | FLASH_CACHE [ { GLOBAL | LOCAL } ]
//	        | REDO TO target_db_name [ { NO CONFIRM APPLY | CONFIRM APPLY } ]
//	        | PASSWORDFILE_METADATA_CACHE
//	        }
func (p *Parser) parseAlterSystemStmt(start int) (nodes.StmtNode, error) {
	p.advance() // consume SYSTEM

	stmt := &nodes.AlterSystemStmt{
		Loc: nodes.Loc{Start: start},
	}

	switch {
	case p.cur.Type == kwSET:
		p.advance() // consume SET
		// Check for SET ENCRYPTION (security clause)
		if p.isIdentLikeStr("ENCRYPTION") {
			parseErr5 := p.parseAlterSystemEncryption(stmt)
			if parseErr5 != nil {
				return nil, parseErr5
			}
		} else {
			stmt.Action = "SET"
			var parseErr6 error
			stmt.SetParams, parseErr6 = p.parseSetParams()
			if parseErr6 !=
				// Parse optional SET modifiers: COMMENT, DEFERRED, SCOPE, SID, CONTAINER
				nil {
				return nil, parseErr6
			}
			parseErr7 := p.parseAlterSystemSetModifiers(stmt)
			if parseErr7 != nil {
				return nil, parseErr7
			}
		}

	case p.isIdentLikeStr("RESET"):
		p.advance() // consume RESET
		stmt.Action = "RESET"
		var parseErr8 error
		stmt.ResetParam, parseErr8 = p.parseIdentifier()
		if parseErr8 !=
			// optional SCOPE
			nil {
			return nil, parseErr8
		}

		if p.isIdentLikeStr("SCOPE") {
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			var parseErr9 error
			stmt.Scope, parseErr9 = p.parseIdentifier()
			if parseErr9 !=

				// optional SID
				nil {
				return nil, parseErr9
			}
		}

		if p.isIdentLikeStr("SID") {
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			if p.cur.Type == tokSCONST {
				stmt.SID = p.cur.Str
				p.advance()
			} else if p.cur.Type == '*' {
				stmt.SID = "*"
				p.advance()
			}
		}

	case p.isIdentLikeStr("KILL"):
		p.advance() // consume KILL
		stmt.Action = "KILL_SESSION"
		if p.cur.Type == kwSESSION {
			p.advance()
		}
		if p.cur.Type == tokSCONST {
			stmt.SessionID = p.cur.Str
			p.advance()
		}
		// optional , @instance_id
		if p.cur.Type == ',' {
			p.advance()
			if p.cur.Type == '@' {
				p.advance()
				if p.cur.Type == tokICONST {
					stmt.InstanceID = p.cur.Str
					p.advance()
				}
			}
		}
		// optional IMMEDIATE | FORCE
		if p.cur.Type == kwIMMEDIATE {
			stmt.Immediate = true
			p.advance()
		} else if p.cur.Type == kwFORCE {
			stmt.Force = true
			p.advance()
		}
		// optional NOREPLAY
		if p.isIdentLikeStr("NOREPLAY") {
			stmt.NoReplay = true
			p.advance()
		}
		// optional TIMEOUT integer
		if p.isIdentLikeStr("TIMEOUT") {
			p.advance()
			if p.cur.Type == tokICONST {
				var parseErr10 error
				stmt.Timeout, parseErr10 = p.parseIntValue()
				if parseErr10 != nil {
					return nil, parseErr10
				}
			}
		}

	case p.isIdentLikeStr("DISCONNECT"):
		p.advance() // consume DISCONNECT
		stmt.Action = "DISCONNECT_SESSION"
		if p.cur.Type == kwSESSION {
			p.advance()
		}
		if p.cur.Type == tokSCONST {
			stmt.SessionID = p.cur.Str
			p.advance()
		}
		// optional POST_TRANSACTION
		if p.isIdentLikeStr("POST_TRANSACTION") {
			stmt.PostTransaction = true
			p.advance()
		}
		// optional IMMEDIATE
		if p.cur.Type == kwIMMEDIATE {
			stmt.Immediate = true
			p.advance()
		}

	case p.isIdentLikeStr("FLUSH"):
		p.advance() // consume FLUSH
		stmt.Action = "FLUSH"
		parseErr11 := p.parseAlterSystemFlush(stmt)
		if parseErr11 != nil {
			return nil, parseErr11
		}

	case p.isIdentLikeStr("CHECKPOINT"):
		p.advance() // consume CHECKPOINT
		stmt.Action = "CHECKPOINT"
		if p.cur.Type == kwGLOBAL {
			stmt.CheckScope = "GLOBAL"
			p.advance()
		} else if p.isIdentLikeStr("LOCAL") {
			stmt.CheckScope = "LOCAL"
			p.advance()
		}

	case p.isIdentLikeStr("CHECK"):
		p.advance() // consume CHECK
		stmt.Action = "CHECK_DATAFILES"
		if p.isIdentLikeStr("DATAFILES") {
			p.advance()
		}
		if p.cur.Type == kwGLOBAL {
			stmt.CheckScope = "GLOBAL"
			p.advance()
		} else if p.isIdentLikeStr("LOCAL") {
			stmt.CheckScope = "LOCAL"
			p.advance()
		}

	case p.isIdentLikeStr("SWITCH"):
		p.advance() // consume SWITCH
		stmt.Action = "SWITCH_LOGFILE"
		if p.isIdentLikeStr("LOGFILE") {
			p.advance()
		}

	case p.isIdentLikeStr("ARCHIVE"):
		p.advance() // consume ARCHIVE
		stmt.Action = "ARCHIVE_LOG"
		if p.isIdentLikeStr("LOG") {
			p.advance()
		}
		parseErr12 := p.parseAlterSystemArchiveLog(stmt)
		if parseErr12 != nil {
			return nil, parseErr12
		}

	case p.isIdentLikeStr("SUSPEND"):
		stmt.Action = "SUSPEND"
		p.advance()

	case p.isIdentLikeStr("RESUME"):
		stmt.Action = "RESUME"
		p.advance()

	case p.isIdentLikeStr("QUIESCE"):
		p.advance() // consume QUIESCE
		stmt.Action = "QUIESCE"
		if p.isIdentLikeStr("RESTRICTED") {
			p.advance()
		}

	case p.isIdentLikeStr("UNQUIESCE"):
		stmt.Action = "UNQUIESCE"
		p.advance()

	case p.cur.Type == kwENABLE:
		p.advance() // consume ENABLE
		stmt.Action = "ENABLE"
		if p.isIdentLikeStr("DISTRIBUTED") {
			stmt.Feature = "DISTRIBUTED_RECOVERY"
			p.advance()
			if p.isIdentLikeStr("RECOVERY") {
				p.advance()
			}
		} else if p.isIdentLikeStr("RESTRICTED") {
			stmt.Feature = "RESTRICTED_SESSION"
			p.advance()
			if p.cur.Type == kwSESSION {
				p.advance()
			}
		}

	case p.cur.Type == kwDISABLE:
		p.advance() // consume DISABLE
		stmt.Action = "DISABLE"
		if p.isIdentLikeStr("DISTRIBUTED") {
			stmt.Feature = "DISTRIBUTED_RECOVERY"
			p.advance()
			if p.isIdentLikeStr("RECOVERY") {
				p.advance()
			}
		} else if p.isIdentLikeStr("RESTRICTED") {
			stmt.Feature = "RESTRICTED_SESSION"
			p.advance()
			if p.cur.Type == kwSESSION {
				p.advance()
			}
		}

	case p.isIdentLikeStr("REGISTER"):
		stmt.Action = "REGISTER"
		p.advance()

	case p.isIdentLikeStr("CANCEL"):
		p.advance() // consume CANCEL
		stmt.Action = "CANCEL_SQL"
		if p.isIdentLikeStr("SQL") {
			p.advance()
		}
		if p.cur.Type == tokSCONST {
			stmt.SessionID = p.cur.Str
			p.advance()
		}
		// optional , @instance_id
		if p.cur.Type == ',' {
			p.advance()
			if p.cur.Type == '@' {
				p.advance()
				if p.cur.Type == tokICONST {
					stmt.InstanceID = p.cur.Str
					p.advance()
				}
			}
		}
		// optional SQL_ID 'sql_id'
		if p.isIdentLikeStr("SQL_ID") {
			p.advance()
			if p.cur.Type == tokSCONST {
				stmt.SqlID = p.cur.Str
				p.advance()
			}
		}

	case p.isIdentLikeStr("SHUTDOWN"):
		p.advance() // consume SHUTDOWN
		stmt.Action = "SHUTDOWN"
		if p.cur.Type == kwIMMEDIATE {
			stmt.Immediate = true
			p.advance()
		}
		if p.cur.Type == tokSCONST {
			stmt.ShutdownDisp = p.cur.Str
			p.advance()
		}

	case p.isIdentLikeStr("RELOCATE"):
		p.advance() // consume RELOCATE
		stmt.Action = "RELOCATE_CLIENT"
		if p.isIdentLikeStr("CLIENT") {
			p.advance()
		}
		if p.cur.Type == tokSCONST {
			stmt.RelocateClient = p.cur.Str
			p.advance()
		}

	case p.cur.Type == kwSTART:
		p.advance() // consume START
		if p.isIdentLikeStr("ROLLING") {
			p.advance() // consume ROLLING
			if p.isIdentLikeStr("MIGRATION") {
				stmt.Action = "START_ROLLING_MIGRATION"
				p.advance()
				if p.cur.Type == kwTO {
					p.advance()
				}
				if p.cur.Type == tokSCONST {
					stmt.RollingVersion = p.cur.Str
					p.advance()
				}
			} else if p.isIdentLikeStr("PATCH") {
				stmt.Action = "START_ROLLING_PATCH"
				p.advance()
			}
		}

	case p.isIdentLikeStr("STOP"):
		p.advance() // consume STOP
		if p.isIdentLikeStr("ROLLING") {
			p.advance()
			if p.isIdentLikeStr("MIGRATION") {
				stmt.Action = "STOP_ROLLING_MIGRATION"
				p.advance()
			} else if p.isIdentLikeStr("PATCH") {
				stmt.Action = "STOP_ROLLING_PATCH"
				p.advance()
			}
		}

	default:
		// All BNF branches are covered above; this handles truly unrecognized tokens.
		p.skipToSemicolon()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterSystemSetModifiers parses optional modifiers for ALTER SYSTEM SET:
// COMMENT, DEFERRED, SCOPE, SID, CONTAINER.
func (p *Parser) parseAlterSystemSetModifiers(stmt *nodes.AlterSystemStmt) error {
	for {
		switch {
		case p.isIdentLikeStr("COMMENT"):
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			if p.cur.Type == tokSCONST {
				stmt.Comment = p.cur.Str
				p.advance()
			}
		case p.isIdentLikeStr("DEFERRED"):
			stmt.Deferred = true
			p.advance()
		case p.isIdentLikeStr("SCOPE"):
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			var parseErr13 error
			stmt.Scope, parseErr13 = p.parseIdentifier()
			if parseErr13 != nil {
				return parseErr13
			}
		case p.isIdentLikeStr("SID"):
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			if p.cur.Type == tokSCONST {
				stmt.SID = p.cur.Str
				p.advance()
			} else if p.cur.Type == '*' {
				stmt.SID = "*"
				p.advance()
			}
		case p.isIdentLikeStr("CONTAINER"):
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			var parseErr14 error
			stmt.Container, parseErr14 = p.parseIdentifier()
			if parseErr14 != nil {
				return parseErr14

				// parseAlterSystemFlush parses the FLUSH sub-clause of ALTER SYSTEM.
			}
		default:
			return nil
		}
	}
	return nil
}

func (p *Parser) parseAlterSystemFlush(stmt *nodes.AlterSystemStmt) error {
	switch {
	case p.isIdentLikeStr("SHARED_POOL"):
		stmt.FlushTarget = "SHARED_POOL"
		p.advance()
	case p.cur.Type == kwGLOBAL:
		// GLOBAL CONTEXT
		stmt.FlushTarget = "GLOBAL_CONTEXT"
		p.advance()
		if p.isIdentLikeStr("CONTEXT") {
			p.advance()
		}
	case p.isIdentLikeStr("BUFFER_CACHE"):
		stmt.FlushTarget = "BUFFER_CACHE"
		p.advance()
		if p.cur.Type == kwGLOBAL {
			stmt.FlushScope = "GLOBAL"
			p.advance()
		} else if p.isIdentLikeStr("LOCAL") {
			stmt.FlushScope = "LOCAL"
			p.advance()
		}
	case p.isIdentLikeStr("FLASH_CACHE"):
		stmt.FlushTarget = "FLASH_CACHE"
		p.advance()
		if p.cur.Type == kwGLOBAL {
			stmt.FlushScope = "GLOBAL"
			p.advance()
		} else if p.isIdentLikeStr("LOCAL") {
			stmt.FlushScope = "LOCAL"
			p.advance()
		}
	case p.isIdentLikeStr("REDO"):
		stmt.FlushTarget = "REDO"
		p.advance()
		if p.cur.Type == kwTO {
			p.advance()
		}
		var parseErr15 error
		stmt.FlushRedoDB, parseErr15 = p.parseIdentifier()
		if parseErr15 !=
			// optional {NO CONFIRM APPLY | CONFIRM APPLY}
			nil {
			return parseErr15
		}

		if p.isIdentLikeStr("NO") {
			p.advance() // consume NO
			if p.isIdentLikeStr("CONFIRM") {
				p.advance()
			}
			if p.isIdentLikeStr("APPLY") {
				p.advance()
			}
			stmt.FlushRedoConfirm = "NO_CONFIRM_APPLY"
		} else if p.isIdentLikeStr("CONFIRM") {
			p.advance()
			if p.isIdentLikeStr("APPLY") {
				p.advance()
			}
			stmt.FlushRedoConfirm = "CONFIRM_APPLY"
		}
	case p.isIdentLikeStr("PASSWORDFILE_METADATA_CACHE"):
		stmt.FlushTarget = "PASSWORDFILE_METADATA_CACHE"
		p.advance()
	}
	return nil
}

// parseAlterSystemArchiveLog parses the ARCHIVE LOG sub-clause of ALTER SYSTEM.
func (p *Parser) parseAlterSystemArchiveLog(stmt *nodes.AlterSystemStmt) error {
	// optional INSTANCE 'instance_name'
	if p.isIdentLikeStr("INSTANCE") {
		p.advance()
		if p.cur.Type == tokSCONST {
			stmt.ArchiveInstance = p.cur.Str
			p.advance()
		}
	}

	// archive log spec
	switch {
	case p.isIdentLikeStr("SEQUENCE"):
		stmt.ArchiveLogSpec = "SEQUENCE"
		p.advance()
		if p.cur.Type == tokICONST {
			stmt.ArchiveLogValue = p.cur.Str
			p.advance()
		}
	case p.isIdentLikeStr("CHANGE"):
		stmt.ArchiveLogSpec = "CHANGE"
		p.advance()
		if p.cur.Type == tokICONST {
			stmt.ArchiveLogValue = p.cur.Str
			p.advance()
		}
	case p.isIdentLikeStr("CURRENT"):
		stmt.ArchiveLogSpec = "CURRENT"
		p.advance()
		if p.isIdentLikeStr("NOSWITCH") {
			stmt.ArchiveNoSwitch = true
			p.advance()
		}
	case p.cur.Type == kwGROUP:
		stmt.ArchiveLogSpec = "GROUP"
		p.advance()
		if p.cur.Type == tokICONST {
			stmt.ArchiveLogValue = p.cur.Str
			p.advance()
		}
	case p.isIdentLikeStr("LOGFILE"):
		stmt.ArchiveLogSpec = "LOGFILE"
		p.advance()
		if p.cur.Type == tokSCONST {
			stmt.ArchiveLogValue = p.cur.Str
			p.advance()
		}
		// optional USING BACKUP CONTROLFILE
		if p.cur.Type == kwUSING {
			p.advance()
			if p.isIdentLikeStr("BACKUP") {
				p.advance()
			}
			if p.isIdentLikeStr("CONTROLFILE") {
				p.advance()
			}
			stmt.ArchiveBackupCF = true
		}
	case p.cur.Type == kwNEXT:
		stmt.ArchiveLogSpec = "NEXT"
		p.advance()
	case p.cur.Type == kwALL:
		stmt.ArchiveLogSpec = "ALL"
		p.advance()
	}

	// optional THREAD integer
	if p.isIdentLikeStr("THREAD") {
		p.advance()
		if p.cur.Type == tokICONST {
			var parseErr16 error
			stmt.ArchiveThread, parseErr16 = p.parseIntValue()
			if parseErr16 !=

				// optional TO 'location'
				nil {
				return parseErr16
			}
		}
	}

	if p.cur.Type == kwTO {
		p.advance()
		if p.cur.Type == tokSCONST {
			stmt.ArchiveTo = p.cur.Str
			p.advance()
		}
	}
	return nil
}

// parseAlterSystemEncryption parses the SET ENCRYPTION sub-clause of ALTER SYSTEM.
func (p *Parser) parseAlterSystemEncryption(stmt *nodes.AlterSystemStmt) error {
	stmt.Action = "SET_ENCRYPTION"
	p.advance() // consume ENCRYPTION
	if p.isIdentLikeStr("WALLET") {
		p.advance() // consume WALLET
		if p.cur.Type == kwOPEN {
			stmt.EncryptionAction = "OPEN"
			p.advance()
			// IDENTIFIED BY password
			if p.isIdentLikeStr("IDENTIFIED") {
				p.advance()
				if p.cur.Type == kwBY {
					p.advance()
				}
				parseDiscard18, parseErr17 := p.parseIdentifier()
				_ = // consume password (not stored for security)
					parseDiscard18
				if parseErr17 != nil {
					return parseErr17
				}
			}
		} else if p.isIdentLikeStr("CLOSE") {
			stmt.EncryptionAction = "CLOSE"
			p.advance()
			// optional IDENTIFIED BY password
			if p.isIdentLikeStr("IDENTIFIED") {
				p.advance()
				if p.cur.Type == kwBY {
					p.advance()
				}
				parseDiscard20, parseErr19 := p.parseIdentifier()
				_ = parseDiscard20
				if parseErr19 != nil {
					return parseErr19
				}
			}
		}
	} else if p.cur.Type == kwKEY {
		stmt.EncryptionAction = "SET_KEY"
		p.advance()
		// optional IDENTIFIED BY password
		if p.isIdentLikeStr("IDENTIFIED") {
			p.advance()
			if p.cur.Type == kwBY {
				p.advance()
			}
			parseDiscard22, parseErr21 := p.parseIdentifier()
			_ = parseDiscard22

			// parseSetParams parses one or more param = value pairs.
			if parseErr21 != nil {
				return parseErr21
			}
		}
	}
	return nil
}

func (p *Parser) parseSetParams() (*nodes.List, error) {
	params := &nodes.List{}
	for {
		param, parseErr23 := p.parseSetParam()
		if parseErr23 != nil {
			return nil, parseErr23
		}
		if param == nil {
			break
		}
		params.Items = append(params.Items, param)
		// Some Oracle ALTER SESSION SET supports multiple params without commas;
		// but also handle comma separation.
		if p.cur.Type == ',' {
			p.advance()
		}
		// Stop if we hit end of statement.
		if !p.isIdentLike() {
			break
		}
	}
	return params, nil
}

// parseSetParam parses a single name = value parameter setting.
func (p *Parser) parseSetParam() (*nodes.SetParam, error) {
	if !p.isIdentLike() {
		return nil, nil
	}
	start := p.pos()
	name, parseErr24 := p.parseIdentifier()
	if parseErr24 != nil {
		return nil, parseErr24

		// Expect '='
	}
	if name == "" {
		return nil, nil
	}

	if p.cur.Type != '=' {
		return &nodes.SetParam{
			Name: name,
			Loc:  nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	}
	p.advance() // consume '='

	value, parseErr25 := p.parseExpr()
	if parseErr25 != nil {
		return nil, parseErr25
	}

	return &nodes.SetParam{
		Name:  name,
		Value: value,
		Loc:   nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseAlterMaterializedViewStmt parses an ALTER MATERIALIZED VIEW statement.
// Called after ALTER MATERIALIZED VIEW has been consumed.
//
// BNF: oracle/parser/bnf/ALTER-MATERIALIZED-VIEW.bnf
//
//	ALTER MATERIALIZED VIEW [ IF EXISTS ] [ schema. ] materialized_view
//	    { physical_attributes_clause
//	    | modify_mv_column_clause
//	    | table_compression
//	    | inmemory_table_clause
//	    | LOB_storage_clause
//	    | modify_LOB_storage_clause
//	    | alter_table_partitioning
//	    | parallel_clause
//	    | logging_clause
//	    | allocate_extent_clause
//	    | deallocate_unused_clause
//	    | shrink_clause
//	    | { CACHE | NOCACHE }
//	    | annotations_clause
//	    | alter_iot_clauses
//	    | MODIFY scoped_table_ref_constraint
//	    | alter_mv_refresh
//	    | evaluation_edition_clause
//	    | { ENABLE | DISABLE } ON QUERY COMPUTATION
//	    | alter_query_rewrite_clause
//	    | COMPILE
//	    | CONSIDER FRESH
//	    } ;
//
//	physical_attributes_clause:
//	    [ PCTFREE integer ]
//	    [ PCTUSED integer ]
//	    [ INITRANS integer ]
//	    [ storage_clause ]
//
//	alter_mv_refresh:
//	    REFRESH
//	    { FAST | COMPLETE | FORCE }
//	    [ { ON COMMIT | ON DEMAND } ]
//	    [ START WITH date_expression ]
//	    [ NEXT date_expression ]
//	    [ USING ROLLBACK SEGMENT rollback_segment_name ]
//	    [ WITH PRIMARY KEY ]
//	    [ { ENABLE | DISABLE } CONCURRENT REFRESH ]
//
//	evaluation_edition_clause:
//	    EVALUATE USING EDITION edition_name
//	    [ CONSIDER FRESH ]
//
//	alter_query_rewrite_clause:
//	    QUERY REWRITE
//	    { ENABLE [ unusable_editions_clause ]
//	    | DISABLE
//	    | unusable_editions_clause
//	    }
//
//	allocate_extent_clause:
//	    ALLOCATE EXTENT [ ( [ SIZE size_clause ] [ DATAFILE 'filename' ] [ INSTANCE integer ] ) ]
//
//	deallocate_unused_clause:
//	    DEALLOCATE UNUSED [ KEEP size_clause ]
//
//	shrink_clause:
//	    SHRINK SPACE [ COMPACT ] [ CASCADE ]
func (p *Parser) parseAlterMaterializedViewStmt(start int) (nodes.StmtNode, error) {
	stmt := &nodes.AlterMaterializedViewStmt{
		Loc: nodes.Loc{Start: start},
	}

	// IF EXISTS
	if p.cur.Type == kwIF {
		if p.peekNext().Type == kwEXISTS {
			stmt.IfExists = true
			p.advance() // consume IF
			p.advance() // consume EXISTS
		}
	}
	var parseErr26 error

	stmt.Name, parseErr26 = p.parseObjectName()
	if parseErr26 !=

		// Parse action clause
		nil {
		return nil, parseErr26
	}

	switch {
	case p.isIdentLikeStr("COMPILE"):
		stmt.Action = "COMPILE"
		p.advance() // consume COMPILE

	case p.isIdentLikeStr("CONSIDER"):
		stmt.Action = "CONSIDER_FRESH"
		p.advance() // consume CONSIDER
		if p.isIdentLikeStr("FRESH") {
			p.advance() // consume FRESH
		}

	case p.cur.Type == kwREFRESH:
		stmt.Action = "REFRESH"
		p.advance()
		parseErr27 := // consume REFRESH
			p.parseAlterMViewRefreshClause(stmt)
		if parseErr27 != nil {
			return nil, parseErr27
		}

	case p.cur.Type == kwENABLE:
		p.advance() // consume ENABLE
		if p.isIdentLikeStr("QUERY") {
			stmt.Action = "ENABLE_QUERY_REWRITE"
			p.advance() // consume QUERY
			if p.cur.Type == kwREWRITE {
				p.advance() // consume REWRITE
			}
		} else if p.isIdentLikeStr("CONCURRENT") {
			stmt.Action = "ENABLE_CONCURRENT_REFRESH"
			p.advance() // consume CONCURRENT
			if p.cur.Type == kwREFRESH {
				p.advance() // consume REFRESH
			}
		} else if p.cur.Type == kwON {
			// ENABLE ON QUERY COMPUTATION
			stmt.Action = "ENABLE_ON_QUERY_COMPUTATION"
			p.advance() // consume ON
			if p.isIdentLikeStr("QUERY") {
				p.advance() // consume QUERY
			}
			if p.isIdentLikeStr("COMPUTATION") {
				p.advance() // consume COMPUTATION
			}
		} else {
			stmt.Action = "ENABLE_QUERY_REWRITE"
		}

	case p.cur.Type == kwDISABLE:
		p.advance() // consume DISABLE
		if p.isIdentLikeStr("QUERY") {
			stmt.Action = "DISABLE_QUERY_REWRITE"
			p.advance() // consume QUERY
			if p.cur.Type == kwREWRITE {
				p.advance() // consume REWRITE
			}
		} else if p.isIdentLikeStr("CONCURRENT") {
			stmt.Action = "DISABLE_CONCURRENT_REFRESH"
			p.advance() // consume CONCURRENT
			if p.cur.Type == kwREFRESH {
				p.advance() // consume REFRESH
			}
		} else if p.cur.Type == kwON {
			// DISABLE ON QUERY COMPUTATION
			stmt.Action = "DISABLE_ON_QUERY_COMPUTATION"
			p.advance() // consume ON
			if p.isIdentLikeStr("QUERY") {
				p.advance() // consume QUERY
			}
			if p.isIdentLikeStr("COMPUTATION") {
				p.advance() // consume COMPUTATION
			}
		} else {
			stmt.Action = "DISABLE_QUERY_REWRITE"
		}

	case p.isIdentLikeStr("SHRINK"):
		stmt.Action = "SHRINK"
		p.advance() // consume SHRINK
		if p.isIdentLikeStr("SPACE") {
			p.advance() // consume SPACE
		}
		if p.isIdentLikeStr("COMPACT") {
			stmt.Compact = true
			p.advance()
		} else if p.cur.Type == kwCASCADE {
			stmt.Cascade = true
			p.advance()
		}

	case p.cur.Type == kwCACHE:
		stmt.Action = "CACHE"
		p.advance()

	case p.cur.Type == kwNOCACHE:
		stmt.Action = "NOCACHE"
		p.advance()

	case p.cur.Type == kwPARALLEL:
		stmt.Action = "PARALLEL"
		p.advance() // consume PARALLEL
		if p.cur.Type == tokICONST {
			stmt.ParallelDegree = p.cur.Str
			p.advance()
		}

	case p.cur.Type == kwNOPARALLEL:
		stmt.Action = "NOPARALLEL"
		p.advance()

	case p.cur.Type == kwLOGGING:
		stmt.Action = "LOGGING"
		p.advance()

	case p.cur.Type == kwNOLOGGING:
		stmt.Action = "NOLOGGING"
		p.advance()

	case p.isIdentLikeStr("ALLOCATE"):
		stmt.Action = "ALLOCATE_EXTENT"
		p.advance() // consume ALLOCATE
		if p.isIdentLikeStr("EXTENT") {
			p.advance() // consume EXTENT
		}
		// optional ( SIZE size DATAFILE 'filename' INSTANCE integer )
		if p.cur.Type == '(' {
			p.advance()
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				switch {
				case p.isIdentLikeStr("SIZE"):
					p.advance()
					var parseErr28 error
					stmt.AllocateSize, parseErr28 = p.parseIdentifier()
					if parseErr28 != nil {
						return nil, parseErr28
					}
				case p.isIdentLikeStr("DATAFILE"):
					p.advance()
					if p.cur.Type == tokSCONST {
						stmt.AllocateDatafile = p.cur.Str
						p.advance()
					}
				case p.isIdentLikeStr("INSTANCE"):
					p.advance()
					if p.cur.Type == tokICONST {
						var parseErr29 error
						stmt.AllocateInstance, parseErr29 = p.parseIntValue()
						if parseErr29 != nil {
							return nil, parseErr29
						}
					}
				default:
					p.advance()
				}
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}

	case p.isIdentLikeStr("DEALLOCATE"):
		stmt.Action = "DEALLOCATE_UNUSED"
		p.advance() // consume DEALLOCATE
		if p.isIdentLikeStr("UNUSED") {
			p.advance()
		}
		if p.cur.Type == kwKEEP {
			p.advance()
			// size_clause: integer [K|M|G|T|P|E]
			if p.cur.Type == tokICONST {
				val := p.cur.Str
				p.advance()
				// optional size unit
				if p.isIdentLikeStr("K") || p.isIdentLikeStr("M") || p.isIdentLikeStr("G") || p.isIdentLikeStr("T") || p.isIdentLikeStr("P") || p.isIdentLikeStr("E") {
					val += p.cur.Str
					p.advance()
				}
				stmt.DeallocateKeep = val
			} else {
				var parseErr30 error
				stmt.DeallocateKeep, parseErr30 = p.parseIdentifier()
				if parseErr30 != nil {
					return nil, parseErr30
				}
			}
		}

	case p.isIdentLikeStr("PCTFREE"):
		stmt.Action = "PCTFREE"
		p.advance()
		if p.cur.Type == tokICONST {
			var parseErr31 error
			stmt.PctFree, parseErr31 = p.parseIntValue()
			if parseErr31 != nil {
				return nil, parseErr31
			}
		}

	case p.isIdentLikeStr("PCTUSED"):
		stmt.Action = "PCTUSED"
		p.advance()
		if p.cur.Type == tokICONST {
			var parseErr32 error
			stmt.PctUsed, parseErr32 = p.parseIntValue()
			if parseErr32 != nil {
				return nil, parseErr32
			}
		}

	case p.isIdentLikeStr("INITRANS"):
		stmt.Action = "INITRANS"
		p.advance()
		if p.cur.Type == tokICONST {
			var parseErr33 error
			stmt.IniTrans, parseErr33 = p.parseIntValue()
			if parseErr33 != nil {
				return nil, parseErr33
			}
		}

	case p.isIdentLikeStr("EVALUATE"):
		stmt.Action = "EVALUATE_USING_EDITION"
		p.advance() // consume EVALUATE
		if p.cur.Type == kwUSING {
			p.advance() // consume USING
		}
		if p.isIdentLikeStr("EDITION") {
			p.advance() // consume EDITION
		}
		var parseErr34 error
		stmt.EditionName, parseErr34 = p.parseIdentifier()
		if parseErr34 != nil {
			return nil, parseErr34
		}

	case p.cur.Type == kwCOMPRESS:
		stmt.Action = "COMPRESS"
		p.advance()
		// optional BASIC or FOR { OLTP | QUERY | ARCHIVE }
		if p.isIdentLikeStr("BASIC") {
			p.advance()
		} else if p.cur.Type == kwFOR {
			p.advance()
			parseDiscard36, parseErr35 := p.parseIdentifier()
			_ = // OLTP, QUERY, ARCHIVE
				parseDiscard36
			if parseErr35 !=
				// optional LOW | HIGH
				nil {
				return nil, parseErr35
			}

			if p.isIdentLikeStr("LOW") || p.isIdentLikeStr("HIGH") {
				p.advance()
			}
		}

	case p.isIdentLikeStr("NOCOMPRESS"):
		stmt.Action = "NOCOMPRESS"
		p.advance()

	case p.isIdentLikeStr("INMEMORY"):
		stmt.Action = "INMEMORY"
		p.advance()
		parseErr37 :=
			// skip remaining in-memory attributes
			p.parseAlterMViewInMemoryAttrs()
		if parseErr37 != nil {
			return nil, parseErr37

			// NO INMEMORY
		}

	case p.isIdentLikeStr("NO"):

		p.advance() // consume NO
		if p.isIdentLikeStr("INMEMORY") {
			stmt.Action = "NO_INMEMORY"
			p.advance()
		}

	case p.cur.Type == kwMODIFY:
		p.advance() // consume MODIFY
		// MODIFY scoped_table_ref_constraint: SCOPE FOR ( ref_col ) IS [schema.]table
		if p.isIdentLikeStr("SCOPE") {
			stmt.Action = "MODIFY_SCOPE"
			p.advance() // consume SCOPE
			if p.cur.Type == kwFOR {
				p.advance() // consume FOR
			}
			if p.cur.Type == '(' {
				p.advance()
				var // consume (
				parseErr38 error
				stmt.ScopeColumn, parseErr38 = p.parseIdentifier()
				if parseErr38 != nil {
					return nil, parseErr38
				}
				if p.cur.Type == ')' {
					p.advance() // consume )
				}
			}
			if p.cur.Type == kwIS {
				p.advance() // consume IS
			}
			var parseErr39 error
			stmt.ScopeTable, parseErr39 = p.parseObjectName()
			if parseErr39 != nil {

				// modify_mv_column_clause: MODIFY ( column_name encryption_spec [, ...] )
				return nil, parseErr39
			}
		} else {

			stmt.Action = "MODIFY"
			if p.cur.Type == '(' {
				p.advance() // consume (
				for p.cur.Type != ')' && p.cur.Type != tokEOF {
					// column_name
					columnName, parseErr40 := p.parseIdentifier()
					_ = columnName
					if parseErr40 != nil {
						return nil, parseErr40
					}
					if p.isIdentLikeStr("ENCRYPT") {
						p.advance() // consume ENCRYPT
						// parse encryption_spec options
						for {
							if p.cur.Type == kwUSING {
								p.advance() // consume USING
								if p.cur.Type == tokSCONST {
									p.advance() // consume algorithm string
								}
							} else if p.isIdentLikeStr("IDENTIFIED") {
								p.advance() // consume IDENTIFIED
								if p.cur.Type == kwBY {
									p.advance()
									parseDiscard42, // consume BY
										parseErr41 := p.parseIdentifier()
									_ = // consume password
										parseDiscard42
									if parseErr41 != nil {
										return nil, parseErr41
									}
								}
							} else if p.isIdentLikeStr("SALT") {
								p.advance()
							} else if p.isIdentLikeStr("NO") {
								next := p.peekNext()
								if (next.Type == tokIDENT || next.Type >= 2000) && next.Str == "SALT" {
									p.advance() // consume NO
									p.advance() // consume SALT
								} else {
									break
								}
							} else {
								break
							}
						}
					}
					if p.cur.Type == ',' {
						p.advance() // consume ,
					} else {
						break
					}
				}
				if p.cur.Type == ')' {
					p.advance() // consume )
				}
			}
		}

	case p.isIdentLikeStr("QUERY"):
		// alter_query_rewrite_clause: QUERY REWRITE { ENABLE | DISABLE | ... }
		stmt.Action = "ENABLE_QUERY_REWRITE"
		p.advance() // consume QUERY
		if p.cur.Type == kwREWRITE {
			p.advance() // consume REWRITE
		}
		if p.cur.Type == kwENABLE {
			stmt.Action = "ENABLE_QUERY_REWRITE"
			p.advance()
		} else if p.cur.Type == kwDISABLE {
			stmt.Action = "DISABLE_QUERY_REWRITE"
			p.advance()
		}

	default:
		// Fallback for truly unrecognized clauses
		p.skipToSemicolon()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterMViewRefreshClause parses the alter_mv_refresh clause.
// Called after REFRESH has been consumed.
//
//	[ FAST | COMPLETE | FORCE ]
//	[ ON { COMMIT | DEMAND } ]
//	[ START WITH date ]
//	[ NEXT date ]
//	[ WITH PRIMARY KEY ]
//	[ USING ROLLBACK SEGMENT rollback_segment ]
//	[ USING { ENFORCED | TRUSTED } CONSTRAINTS ]
//	[ { ENABLE | DISABLE } ON QUERY COMPUTATION ]
func (p *Parser) parseAlterMViewRefreshClause(stmt *nodes.AlterMaterializedViewStmt) error {
	for {
		switch {
		case p.isIdentLikeStr("FAST"):
			stmt.RefreshMethod = "FAST"
			p.advance()
		case p.isIdentLikeStr("COMPLETE"):
			stmt.RefreshMethod = "COMPLETE"
			p.advance()
		case p.cur.Type == kwFORCE:
			stmt.RefreshMethod = "FORCE"
			p.advance()
		case p.cur.Type == kwON:
			p.advance() // consume ON
			if p.cur.Type == kwCOMMIT {
				stmt.RefreshMode = "ON_COMMIT"
				p.advance()
			} else if p.isIdentLikeStr("DEMAND") {
				stmt.RefreshMode = "ON_DEMAND"
				p.advance()
			} else if p.isIdentLikeStr("QUERY") {
				// ON QUERY COMPUTATION - part of ENABLE/DISABLE ON QUERY COMPUTATION
				// This shouldn't happen here, but handle gracefully
				p.advance() // consume QUERY
				if p.isIdentLikeStr("COMPUTATION") {
					p.advance()
				}
			}
		case p.cur.Type == kwSTART:
			p.advance() // consume START
			if p.cur.Type == kwWITH {
				p.advance() // consume WITH
			}
			var parseErr43 error
			stmt.StartWith, parseErr43 = p.parseExpr()
			if parseErr43 != nil {
				return parseErr43
			}
		case p.cur.Type == kwNEXT:
			p.advance()
			var // consume NEXT
			parseErr44 error
			stmt.Next, parseErr44 = p.parseExpr()
			if parseErr44 != nil {
				return parseErr44
			}
		case p.cur.Type == kwWITH:
			p.advance() // consume WITH
			if p.cur.Type == kwPRIMARY {
				p.advance() // consume PRIMARY
				if p.cur.Type == kwKEY {
					p.advance() // consume KEY
				}
				stmt.WithPrimaryKey = true
			}
		case p.cur.Type == kwUSING:
			p.advance() // consume USING
			if p.cur.Type == kwROLLBACK {
				p.advance() // consume ROLLBACK
				if p.isIdentLikeStr("SEGMENT") {
					p.advance() // consume SEGMENT
				}
				if p.isIdentLike() {
					stmt.UsingRollbackSegment = p.cur.Str
					p.advance()
				}
			} else if p.isIdentLikeStr("ENFORCED") {
				stmt.UsingConstraints = "ENFORCED"
				p.advance() // consume ENFORCED
				if p.cur.Type == kwCONSTRAINTS {
					p.advance() // consume CONSTRAINTS
				}
			} else if p.isIdentLikeStr("TRUSTED") {
				stmt.UsingConstraints = "TRUSTED"
				p.advance() // consume TRUSTED
				if p.cur.Type == kwCONSTRAINTS {
					p.advance() // consume CONSTRAINTS
				}
			}
		case p.cur.Type == kwENABLE:
			p.advance() // consume ENABLE
			if p.cur.Type == kwON {
				p.advance() // consume ON
				if p.isIdentLikeStr("QUERY") {
					p.advance() // consume QUERY
				}
				if p.isIdentLikeStr("COMPUTATION") {
					p.advance() // consume COMPUTATION
				}
				stmt.EnableOnQueryComputation = true
			} else {
				return nil // not part of REFRESH clause, back out
			}
		case p.cur.Type == kwDISABLE:
			p.advance() // consume DISABLE
			if p.cur.Type == kwON {
				p.advance() // consume ON
				if p.isIdentLikeStr("QUERY") {
					p.advance() // consume QUERY
				}
				if p.isIdentLikeStr("COMPUTATION") {
					p.advance() // consume COMPUTATION
				}
				stmt.DisableOnQueryComputation = true
			} else {
				return nil // not part of REFRESH clause, back out
			}
		default:
			return nil
		}
	}
	return nil
}

// parseAlterMViewInMemoryAttrs consumes optional INMEMORY attributes:
// MEMCOMPRESS, PRIORITY, DISTRIBUTE, DUPLICATE keywords until we hit
// something that doesn't look like an inmemory attribute.
func (p *Parser) parseAlterMViewInMemoryAttrs() error {
	for {
		switch {
		case p.isIdentLikeStr("MEMCOMPRESS"):
			p.advance()
			if p.cur.Type == kwFOR {
				p.advance()
				parseDiscard46, parseErr45 := p.parseIdentifier()
				_ = // QUERY, OLTP, ARCHIVE
					parseDiscard46
				if parseErr45 != nil {
					return parseErr45
				}
				if p.isIdentLikeStr("LOW") || p.isIdentLikeStr("HIGH") {
					p.advance()
				}
			}
		case p.isIdentLikeStr("PRIORITY"):
			p.advance()
			parseDiscard48, parseErr47 := p.parseIdentifier()
			_ = // LOW, MEDIUM, HIGH, CRITICAL
				parseDiscard48
			if parseErr47 != nil {
				return parseErr47
			}
		case p.isIdentLikeStr("DISTRIBUTE"):
			p.advance()
			if p.cur.Type == kwBY {
				p.advance()
				parseDiscard50, parseErr49 := p.parseIdentifier()
				_ = // ROWID, PARTITION, SUBPARTITION
					parseDiscard50
				if parseErr49 != nil {
					return parseErr49
				}
			}
		case p.isIdentLikeStr("DUPLICATE"):
			p.advance()
			if p.cur.Type == kwALL {
				p.advance()
			} else if p.isIdentLikeStr("NONE") {
				p.advance()
			}
		default:
			return nil
		}
	}
	return nil
}

// parseAlterGeneric parses ALTER INDEX/VIEW/SEQUENCE/TABLE by consuming the
// object name and skipping the rest (simplified). Returns an AlterSessionStmt
// as a placeholder — in practice these would have their own AST types, but for
// now we skip the body to avoid blocking other work.
func (p *Parser) parseAlterGeneric(start int, objType nodes.ObjectType) (nodes.StmtNode, error) {
	p.advance() // consume INDEX/VIEW/SEQUENCE/etc.

	// For MATERIALIZED VIEW, consume VIEW too
	if objType == nodes.OBJECT_MATERIALIZED_VIEW && p.cur.Type == kwVIEW {
		p.advance()
	}
	// For DATABASE LINK, consume LINK too
	if objType == nodes.OBJECT_DATABASE_LINK && p.cur.Type == kwLINK {
		p.advance()
	}

	stmt := &nodes.AdminDDLStmt{
		Action:     "ALTER",
		ObjectType: objType,
		Loc:        nodes.Loc{Start: start},
	}
	var parseErr51 error

	// Parse the object name.
	stmt.Name, parseErr51 = p.parseObjectName()
	if parseErr51 !=

		// Skip remainder of the statement (clauses vary greatly by object type).
		nil {
		return nil, parseErr51
	}

	p.skipToSemicolon()

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterDatabaseLinkStmt parses an ALTER DATABASE LINK statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/ALTER-DATABASE-LINK.html
//
//	ALTER [ SHARED ] [ PUBLIC ] DATABASE LINK dblink_name
//	  CONNECT TO user IDENTIFIED BY password
//	  [ AUTHENTICATED BY user IDENTIFIED BY password ]
func (p *Parser) parseAlterDatabaseLinkStmt(start int, shared bool, public bool) (nodes.StmtNode, error) {
	p.advance() // consume DATABASE
	if p.cur.Type == kwLINK {
		p.advance() // consume LINK
	}

	stmt := &nodes.AlterDatabaseLinkStmt{
		Shared: shared,
		Public: public,
		Loc:    nodes.Loc{Start: start},
	}
	var parseErr52 error

	// Parse the database link name.
	stmt.Name, parseErr52 = p.parseObjectName()
	if parseErr52 !=

		// CONNECT TO user IDENTIFIED BY password
		nil {
		return nil, parseErr52
	}

	if p.cur.Type == kwCONNECT {
		p.advance() // consume CONNECT
		if p.cur.Type == kwTO {
			p.advance() // consume TO
		}
		if p.isIdentLike() || p.cur.Type == tokIDENT {
			stmt.ConnectUser = p.cur.Str
			p.advance() // consume user
		}
		if p.cur.Type == kwIDENTIFIED {
			p.advance() // consume IDENTIFIED
			if p.cur.Type == kwBY {
				p.advance() // consume BY
			}
			if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokSCONST {
				stmt.ConnectPassword = p.cur.Str
				p.advance() // consume password
			}
		}
	}

	// [ AUTHENTICATED BY user IDENTIFIED BY password ]
	if p.isIdentLikeStr("AUTHENTICATED") {
		p.advance() // consume AUTHENTICATED
		if p.cur.Type == kwBY {
			p.advance() // consume BY
		}
		if p.isIdentLike() || p.cur.Type == tokIDENT {
			stmt.AuthenticatedUser = p.cur.Str
			p.advance() // consume user
		}
		if p.cur.Type == kwIDENTIFIED {
			p.advance() // consume IDENTIFIED
			if p.cur.Type == kwBY {
				p.advance() // consume BY
			}
			if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokSCONST {
				stmt.AuthenticatedPass = p.cur.Str
				p.advance() // consume password
			}
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterSynonymStmt parses an ALTER SYNONYM statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/ALTER-SYNONYM.html
//
//	ALTER [ PUBLIC ] SYNONYM [ IF EXISTS ] [ schema. ] synonym
//	  { EDITIONABLE | NONEDITIONABLE | COMPILE }
func (p *Parser) parseAlterSynonymStmt(start int, public bool) (nodes.StmtNode, error) {
	p.advance() // consume SYNONYM

	stmt := &nodes.AlterSynonymStmt{
		Public: public,
		Loc:    nodes.Loc{Start: start},
	}

	// IF EXISTS
	if p.cur.Type == kwIF {
		next := p.peekNext()
		if next.Type == kwEXISTS {
			stmt.IfExists = true
			p.advance() // consume IF
			p.advance() // consume EXISTS
		}
	}
	var parseErr53 error

	// Parse the synonym name.
	stmt.Name, parseErr53 = p.parseObjectName()
	if parseErr53 !=

		// { EDITIONABLE | NONEDITIONABLE | COMPILE }
		nil {
		return nil, parseErr53
	}

	if p.isIdentLikeStr("EDITIONABLE") {
		stmt.Action = "EDITIONABLE"
		p.advance()
	} else if p.isIdentLikeStr("NONEDITIONABLE") {
		stmt.Action = "NONEDITIONABLE"
		p.advance()
	} else if p.isIdentLikeStr("COMPILE") {
		stmt.Action = "COMPILE"
		p.advance()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterIndexStmt parses an ALTER INDEX statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-INDEX.html
//
//	ALTER INDEX [ IF EXISTS ] [ schema. ] index_name
//	    { deallocate_unused_clause
//	    | allocate_extent_clause
//	    | shrink_clause
//	    | parallel_clause
//	    | physical_attributes_clause
//	    | logging_clause
//	    | partial_index_clause
//	    | rebuild_clause
//	    | alter_index_partitioning
//	    | PARAMETERS ( 'odci_parameters' )
//	    | { DEFERRED | IMMEDIATE } INVALIDATION
//	    | COMPILE
//	    | ENABLE
//	    | DISABLE
//	    | { USABLE | UNUSABLE } [ ONLINE ]
//	    | { VISIBLE | INVISIBLE }
//	    | RENAME TO new_index_name
//	    | COALESCE [ CLEANUP [ ONLY ] ] [ parallel_clause ]
//	    | MONITORING USAGE
//	    | NOMONITORING USAGE
//	    | UPDATE BLOCK REFERENCES
//	    | annotations_clause
//	    }
func (p *Parser) parseAlterIndexStmt(start int) (nodes.StmtNode, error) {
	p.advance() // consume INDEX

	stmt := &nodes.AlterIndexStmt{
		Loc: nodes.Loc{Start: start},
	}

	// IF EXISTS
	if p.cur.Type == kwIF {
		next := p.peekNext()
		if next.Type == kwEXISTS {
			stmt.IfExists = true
			p.advance() // consume IF
			p.advance() // consume EXISTS
		}
	}
	var parseErr54 error

	// Parse index name
	stmt.Name, parseErr54 = p.parseObjectName()
	if parseErr54 !=

		// Parse action
		nil {
		return nil, parseErr54
	}

	switch {
	case p.isIdentLikeStr("REBUILD"):
		stmt.Action = "REBUILD"
		p.advance() // consume REBUILD
		// Optional PARTITION/SUBPARTITION and rebuild options
		for p.cur.Type != ';' && p.cur.Type != tokEOF {
			switch {
			case p.cur.Type == kwPARTITION:
				p.advance()
				var // consume PARTITION
				parseErr55 error
				stmt.Partition, parseErr55 = p.parseIdentifier()
				if parseErr55 != nil {
					return nil, parseErr55
				}
			case p.cur.Type == kwSUBPARTITION:
				p.advance()
				var // consume SUBPARTITION
				parseErr56 error
				stmt.Subpartition, parseErr56 = p.parseIdentifier()
				if parseErr56 != nil {
					return nil, parseErr56
				}
			case p.cur.Type == kwTABLESPACE:
				p.advance()
				var // consume TABLESPACE
				parseErr57 error
				stmt.Tablespace, parseErr57 = p.parseIdentifier()
				if parseErr57 != nil {
					return nil, parseErr57
				}
			case p.cur.Type == kwONLINE:
				stmt.Online = true
				p.advance()
			case p.cur.Type == kwREVERSE:
				stmt.Reverse = true
				p.advance()
			case p.isIdentLikeStr("NOREVERSE"):
				stmt.NoReverse = true
				p.advance()
			case p.cur.Type == kwPARALLEL:
				p.advance() // consume PARALLEL
				if p.cur.Type == tokICONST {
					stmt.Parallel = p.cur.Str
					p.advance()
				}
			case p.cur.Type == kwNOPARALLEL:
				stmt.NoParallel = true
				p.advance()
			case p.cur.Type == kwCOMPRESS:
				p.advance() // consume COMPRESS
				if p.isIdentLikeStr("ADVANCED") {
					p.advance()
					if p.isIdentLikeStr("LOW") {
						stmt.Compress = "ADVANCED LOW"
						p.advance()
					} else if p.isIdentLikeStr("HIGH") {
						stmt.Compress = "ADVANCED HIGH"
						p.advance()
					} else {
						stmt.Compress = "ADVANCED"
					}
				} else if p.cur.Type == tokICONST {
					stmt.Compress = p.cur.Str
					p.advance()
				} else {
					stmt.Compress = "1"
				}
			case p.cur.Type == kwNOCOMPRESS:
				stmt.NoCompress = true
				p.advance()
			case p.cur.Type == kwLOGGING:
				stmt.Logging = true
				p.advance()
			case p.cur.Type == kwNOLOGGING:
				stmt.NoLogging = true
				p.advance()
			case p.isIdentLikeStr("PARAMETERS"):
				p.advance() // consume PARAMETERS
				if p.cur.Type == '(' {
					p.advance()
					if p.cur.Type == tokSCONST {
						stmt.Parameters = p.cur.Str
						p.advance()
					}
					if p.cur.Type == ')' {
						p.advance()
					}
				}
			case p.cur.Type == kwPCTFREE:
				p.advance()
				if p.cur.Type == tokICONST {
					stmt.PctFree = p.cur.Str
					p.advance()
				}
			case p.isIdentLikeStr("INITRANS"):
				p.advance()
				if p.cur.Type == tokICONST {
					stmt.InitTrans = p.cur.Str
					p.advance()
				}
			case p.isIdentLikeStr("INDEXING"):
				p.advance()
				if p.isIdentLikeStr("FULL") {
					stmt.IndexingFull = true
					p.advance()
				} else if p.isIdentLikeStr("PARTIAL") {
					stmt.IndexingPartial = true
					p.advance()
				}
			case p.cur.Type == kwDEFERRED:
				stmt.Invalidation = "DEFERRED"
				p.advance()
				if p.isIdentLikeStr("INVALIDATION") {
					p.advance()
				}
			case p.cur.Type == kwIMMEDIATE:
				stmt.Invalidation = "IMMEDIATE"
				p.advance()
				if p.isIdentLikeStr("INVALIDATION") {
					p.advance()
				}
			default:
				goto done
			}
		}

	case p.cur.Type == kwRENAME:
		p.advance() // consume RENAME
		// RENAME TO new_name or RENAME PARTITION/SUBPARTITION name TO new_name
		if p.cur.Type == kwPARTITION {
			stmt.Action = "RENAME_PARTITION"
			p.advance()
			var // consume PARTITION
			parseErr58 error
			stmt.Partition, parseErr58 = p.parseIdentifier()
			if parseErr58 != nil {
				return nil, parseErr58
			}
			if p.cur.Type == kwTO {
				p.advance() // consume TO
			}
			var parseErr59 error
			stmt.NewName, parseErr59 = p.parseIdentifier()
			if parseErr59 != nil {
				return nil, parseErr59
			}
		} else if p.cur.Type == kwSUBPARTITION {
			stmt.Action = "RENAME_SUBPARTITION"
			p.advance()
			var // consume SUBPARTITION
			parseErr60 error
			stmt.Subpartition, parseErr60 = p.parseIdentifier()
			if parseErr60 != nil {
				return nil, parseErr60
			}
			if p.cur.Type == kwTO {
				p.advance() // consume TO
			}
			var parseErr61 error
			stmt.NewName, parseErr61 = p.parseIdentifier()
			if parseErr61 != nil {
				return nil, parseErr61
			}
		} else {
			stmt.Action = "RENAME"
			if p.cur.Type == kwTO {
				p.advance() // consume TO
			}
			var parseErr62 error
			stmt.NewName, parseErr62 = p.parseIdentifier()
			if parseErr62 != nil {
				return nil, parseErr62
			}
		}

	case p.isIdentLikeStr("COALESCE"):
		stmt.Action = "COALESCE"
		p.advance() // consume COALESCE
		// COALESCE PARTITION [parallel_clause]
		if p.cur.Type == kwPARTITION {
			stmt.Action = "COALESCE_PARTITION"
			p.advance() // consume PARTITION
		} else if p.isIdentLikeStr("CLEANUP") {
			stmt.Cleanup = true
			p.advance() // consume CLEANUP
			if p.isIdentLikeStr("ONLY") {
				stmt.CleanupOnly = true
				p.advance() // consume ONLY
			}
		}
		// Optional parallel clause
		if p.cur.Type == kwPARALLEL {
			p.advance()
			if p.cur.Type == tokICONST {
				stmt.Parallel = p.cur.Str
				p.advance()
			}
		} else if p.cur.Type == kwNOPARALLEL {
			stmt.NoParallel = true
			p.advance()
		}

	case p.isIdentLikeStr("MONITORING"):
		stmt.Action = "MONITORING_USAGE"
		p.advance() // consume MONITORING
		if p.isIdentLikeStr("USAGE") {
			p.advance() // consume USAGE
		}

	case p.isIdentLikeStr("NOMONITORING"):
		stmt.Action = "NOMONITORING_USAGE"
		p.advance() // consume NOMONITORING
		if p.isIdentLikeStr("USAGE") {
			p.advance() // consume USAGE
		}

	case p.isIdentLikeStr("UNUSABLE"):
		stmt.Action = "UNUSABLE"
		p.advance() // consume UNUSABLE
		if p.cur.Type == kwONLINE {
			stmt.Online = true
			p.advance()
		}

	case p.isIdentLikeStr("USABLE"):
		stmt.Action = "USABLE"
		p.advance() // consume USABLE

	case p.isIdentLikeStr("VISIBLE"):
		stmt.Action = "VISIBLE"
		p.advance() // consume VISIBLE

	case p.cur.Type == kwINVISIBLE:
		stmt.Action = "INVISIBLE"
		p.advance() // consume INVISIBLE

	case p.cur.Type == kwENABLE:
		stmt.Action = "ENABLE"
		p.advance() // consume ENABLE

	case p.cur.Type == kwDISABLE:
		stmt.Action = "DISABLE"
		p.advance() // consume DISABLE

	case p.isIdentLikeStr("COMPILE"):
		stmt.Action = "COMPILE"
		p.advance() // consume COMPILE

	case p.isIdentLikeStr("SHRINK"):
		stmt.Action = "SHRINK_SPACE"
		p.advance() // consume SHRINK
		if p.isIdentLikeStr("SPACE") {
			p.advance() // consume SPACE
		}
		if p.isIdentLikeStr("COMPACT") {
			stmt.Compact = true
			p.advance()
		}
		if p.cur.Type == kwCASCADE {
			stmt.Cascade = true
			p.advance()
		}

	case p.cur.Type == kwPARALLEL:
		stmt.Action = "PARALLEL"
		p.advance() // consume PARALLEL
		if p.cur.Type == tokICONST {
			stmt.Parallel = p.cur.Str
			p.advance()
		}

	case p.cur.Type == kwNOPARALLEL:
		stmt.Action = "NOPARALLEL"
		p.advance() // consume NOPARALLEL

	case p.cur.Type == kwLOGGING:
		stmt.Action = "LOGGING"
		p.advance() // consume LOGGING

	case p.cur.Type == kwNOLOGGING:
		stmt.Action = "NOLOGGING"
		p.advance() // consume NOLOGGING

	case p.isIdentLikeStr("DEALLOCATE"):
		// deallocate_unused_clause: DEALLOCATE UNUSED [KEEP size_clause]
		stmt.Action = "DEALLOCATE_UNUSED"
		p.advance() // consume DEALLOCATE
		if p.isIdentLikeStr("UNUSED") {
			p.advance() // consume UNUSED
		}
		if p.cur.Type == kwKEEP {
			p.advance()
			var // consume KEEP
			parseErr63 error
			stmt.DeallocateKeep, parseErr63 = p.parseSizeClause()
			if parseErr63 != nil {
				return nil, parseErr63
			}
		}

	case p.isIdentLikeStr("ALLOCATE"):
		// allocate_extent_clause: ALLOCATE EXTENT [(...)]
		stmt.Action = "ALLOCATE_EXTENT"
		p.advance() // consume ALLOCATE
		if p.isIdentLikeStr("EXTENT") {
			p.advance() // consume EXTENT
		}
		p.skipParenthesizedBlock()

	case p.isIdentLikeStr("PARAMETERS"):
		stmt.Action = "PARAMETERS"
		p.advance() // consume PARAMETERS
		if p.cur.Type == '(' {
			p.advance()
			if p.cur.Type == tokSCONST {
				stmt.Parameters = p.cur.Str
				p.advance()
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}

	case p.cur.Type == kwDEFERRED:
		stmt.Action = "INVALIDATION"
		stmt.Invalidation = "DEFERRED"
		p.advance() // consume DEFERRED
		if p.isIdentLikeStr("INVALIDATION") {
			p.advance()
		}

	case p.cur.Type == kwIMMEDIATE:
		stmt.Action = "INVALIDATION"
		stmt.Invalidation = "IMMEDIATE"
		p.advance() // consume IMMEDIATE
		if p.isIdentLikeStr("INVALIDATION") {
			p.advance()
		}

	case p.isIdentLikeStr("UPDATE"):
		stmt.Action = "UPDATE_BLOCK_REFERENCES"
		p.advance() // consume UPDATE
		if p.isIdentLikeStr("BLOCK") {
			p.advance() // consume BLOCK
		}
		if p.isIdentLikeStr("REFERENCES") {
			p.advance() // consume REFERENCES
		}

	case p.isIdentLikeStr("INDEXING"):
		p.advance() // consume INDEXING
		if p.isIdentLikeStr("FULL") {
			stmt.Action = "INDEXING"
			stmt.IndexingFull = true
			p.advance()
		} else if p.isIdentLikeStr("PARTIAL") {
			stmt.Action = "INDEXING"
			stmt.IndexingPartial = true
			p.advance()
		}

	case p.cur.Type == kwPCTFREE:
		stmt.Action = "PHYSICAL_ATTRIBUTES"
		p.advance()
		if p.cur.Type == tokICONST {
			stmt.PctFree = p.cur.Str
			p.advance()
		}
		parseErr64 :=
			// May have more physical attributes
			p.parseAlterIndexPhysicalAttrs(stmt)
		if parseErr64 != nil {
			return nil, parseErr64
		}

	case p.isIdentLikeStr("PCTUSED"):
		stmt.Action = "PHYSICAL_ATTRIBUTES"
		p.advance()
		if p.cur.Type == tokICONST {
			stmt.PctUsed = p.cur.Str
			p.advance()
		}
		parseErr65 := p.parseAlterIndexPhysicalAttrs(stmt)
		if parseErr65 != nil {
			return nil, parseErr65
		}

	case p.isIdentLikeStr("INITRANS"):
		stmt.Action = "PHYSICAL_ATTRIBUTES"
		p.advance()
		if p.cur.Type == tokICONST {
			stmt.InitTrans = p.cur.Str
			p.advance()
		}
		parseErr66 := p.parseAlterIndexPhysicalAttrs(stmt)
		if parseErr66 != nil {
			return nil, parseErr66
		}

	case p.isIdentLikeStr("MAXTRANS"):
		stmt.Action = "PHYSICAL_ATTRIBUTES"
		p.advance()
		if p.cur.Type == tokICONST {
			stmt.MaxTrans = p.cur.Str
			p.advance()
		}
		parseErr67 := p.parseAlterIndexPhysicalAttrs(stmt)
		if parseErr67 != nil {
			return nil, parseErr67
		}

	case p.isIdentLikeStr("STORAGE"):
		stmt.Action = "PHYSICAL_ATTRIBUTES"
		p.advance()
		p.skipParenthesizedBlock()

	case p.cur.Type == kwMODIFY:
		p.advance() // consume MODIFY
		if p.cur.Type == kwDEFAULT {
			// modify_index_default_attrs
			stmt.Action = "MODIFY_DEFAULT_ATTRIBUTES"
			p.advance() // consume DEFAULT
			if p.isIdentLikeStr("ATTRIBUTES") {
				p.advance() // consume ATTRIBUTES
			}
			// Optional FOR PARTITION partition_name
			if p.cur.Type == kwFOR {
				p.advance() // consume FOR
				if p.cur.Type == kwPARTITION {
					p.advance() // consume PARTITION
				}
				var parseErr68 error
				stmt.ModifyDefaultFor, parseErr68 = p.parseIdentifier()
				if parseErr68 !=

					// Skip remaining options (physical_attributes, TABLESPACE, logging)
					nil {
					return nil, parseErr68
				}
			}
			parseErr69 := p.parseAlterIndexPhysicalAttrs(stmt)
			if parseErr69 != nil {
				return nil, parseErr69
			}
		} else if p.cur.Type == kwPARTITION {
			// modify_index_partition
			stmt.Action = "MODIFY_PARTITION"
			p.advance()
			var // consume PARTITION
			parseErr70 error
			stmt.Partition, parseErr70 = p.parseIdentifier()
			if parseErr70 !=
				// Sub-actions: skip remaining tokens
				nil {
				return nil, parseErr70
			}

			for p.cur.Type != ';' && p.cur.Type != tokEOF {
				switch {
				case p.isIdentLikeStr("UNUSABLE"):
					stmt.ModifyPartAction = "UNUSABLE"
					p.advance()
				case p.isIdentLikeStr("COALESCE"):
					stmt.ModifyPartAction = "COALESCE"
					p.advance()
					if p.isIdentLikeStr("CLEANUP") {
						p.advance()
						if p.isIdentLikeStr("ONLY") {
							p.advance()
						}
					}
				case p.isIdentLikeStr("UPDATE"):
					stmt.ModifyPartAction = "UPDATE_BLOCK_REFERENCES"
					p.advance()
					if p.isIdentLikeStr("BLOCK") {
						p.advance()
					}
					if p.isIdentLikeStr("REFERENCES") {
						p.advance()
					}
				case p.isIdentLikeStr("PARAMETERS"):
					stmt.ModifyPartAction = "PARAMETERS"
					p.advance()
					if p.cur.Type == '(' {
						p.advance()
						if p.cur.Type == tokSCONST {
							stmt.Parameters = p.cur.Str
							p.advance()
						}
						if p.cur.Type == ')' {
							p.advance()
						}
					}
				case p.isIdentLikeStr("DEALLOCATE"):
					stmt.ModifyPartAction = "DEALLOCATE"
					p.advance()
					if p.isIdentLikeStr("UNUSED") {
						p.advance()
					}
					if p.cur.Type == kwKEEP {
						p.advance()
						var parseErr71 error
						stmt.DeallocateKeep, parseErr71 = p.parseSizeClause()
						if parseErr71 != nil {
							return nil, parseErr71
						}
					}
				case p.isIdentLikeStr("ALLOCATE"):
					stmt.ModifyPartAction = "ALLOCATE"
					p.advance()
					if p.isIdentLikeStr("EXTENT") {
						p.advance()
					}
					p.skipParenthesizedBlock()
				case p.cur.Type == kwPCTFREE:
					stmt.ModifyPartAction = "PHYSICAL"
					p.advance()
					if p.cur.Type == tokICONST {
						stmt.PctFree = p.cur.Str
						p.advance()
					}
				case p.cur.Type == kwLOGGING:
					stmt.ModifyPartAction = "LOGGING"
					stmt.Logging = true
					p.advance()
				case p.cur.Type == kwNOLOGGING:
					stmt.ModifyPartAction = "NOLOGGING"
					stmt.NoLogging = true
					p.advance()
				case p.cur.Type == kwCOMPRESS:
					stmt.ModifyPartAction = "COMPRESS"
					p.advance()
					if p.cur.Type == tokICONST {
						stmt.Compress = p.cur.Str
						p.advance()
					} else {
						stmt.Compress = "1"
					}
				case p.cur.Type == kwNOCOMPRESS:
					stmt.ModifyPartAction = "NOCOMPRESS"
					stmt.NoCompress = true
					p.advance()
				default:
					goto done
				}
			}
		} else if p.cur.Type == kwSUBPARTITION {
			// modify_index_subpartition
			stmt.Action = "MODIFY_SUBPARTITION"
			p.advance()
			var // consume SUBPARTITION
			parseErr72 error
			stmt.Subpartition, parseErr72 = p.parseIdentifier()
			if parseErr72 !=
				// Sub-actions: UNUSABLE, allocate_extent, deallocate_unused
				nil {
				return nil, parseErr72
			}

			if p.isIdentLikeStr("UNUSABLE") {
				stmt.ModifyPartAction = "UNUSABLE"
				p.advance()
			} else if p.isIdentLikeStr("ALLOCATE") {
				stmt.ModifyPartAction = "ALLOCATE"
				p.advance()
				if p.isIdentLikeStr("EXTENT") {
					p.advance()
				}
				p.skipParenthesizedBlock()
			} else if p.isIdentLikeStr("DEALLOCATE") {
				stmt.ModifyPartAction = "DEALLOCATE"
				p.advance()
				if p.isIdentLikeStr("UNUSED") {
					p.advance()
				}
				if p.cur.Type == kwKEEP {
					p.advance()
					var parseErr73 error
					stmt.DeallocateKeep, parseErr73 = p.parseSizeClause()
					if parseErr73 != nil {
						return nil, parseErr73
					}
				}
			}
		}

	case p.cur.Type == kwADD:
		// add_hash_index_partition
		stmt.Action = "ADD_PARTITION"
		p.advance() // consume ADD
		if p.cur.Type == kwPARTITION {
			p.advance() // consume PARTITION
		}
		// Optional partition name
		if p.isIdentLike() && p.cur.Type != kwTABLESPACE && p.cur.Type != kwCOMPRESS &&
			p.cur.Type != kwNOCOMPRESS && p.cur.Type != kwPARALLEL && p.cur.Type != kwNOPARALLEL &&
			p.cur.Type != ';' && p.cur.Type != tokEOF {
			var parseErr74 error
			stmt.AddPartitionName, parseErr74 = p.parseIdentifier()
			if parseErr74 !=

				// Optional TABLESPACE, index_compression, parallel_clause
				nil {
				return nil, parseErr74
			}
		}

		for p.cur.Type != ';' && p.cur.Type != tokEOF {
			switch {
			case p.cur.Type == kwTABLESPACE:
				p.advance()
				var parseErr75 error
				stmt.Tablespace, parseErr75 = p.parseIdentifier()
				if parseErr75 != nil {
					return nil, parseErr75
				}
			case p.cur.Type == kwCOMPRESS:
				p.advance()
				if p.cur.Type == tokICONST {
					stmt.Compress = p.cur.Str
					p.advance()
				} else {
					stmt.Compress = "1"
				}
			case p.cur.Type == kwNOCOMPRESS:
				stmt.NoCompress = true
				p.advance()
			case p.cur.Type == kwPARALLEL:
				p.advance()
				if p.cur.Type == tokICONST {
					stmt.Parallel = p.cur.Str
					p.advance()
				}
			case p.cur.Type == kwNOPARALLEL:
				stmt.NoParallel = true
				p.advance()
			default:
				goto done
			}
		}

	case p.cur.Type == kwDROP:
		// drop_index_partition
		stmt.Action = "DROP_PARTITION"
		p.advance() // consume DROP
		if p.cur.Type == kwPARTITION {
			p.advance() // consume PARTITION
		}
		var parseErr76 error
		stmt.Partition, parseErr76 = p.parseIdentifier()
		if parseErr76 != nil {
			return nil, parseErr76

			// split_index_partition
		}

	case p.isIdentLikeStr("SPLIT"):

		stmt.Action = "SPLIT_PARTITION"
		p.advance() // consume SPLIT
		if p.cur.Type == kwPARTITION {
			p.advance() // consume PARTITION
		}
		var parseErr77 error
		stmt.SplitPartition, parseErr77 = p.parseIdentifier()
		if parseErr77 !=
			// AT ( literal [, literal ]... )
			nil {
			return nil, parseErr77
		}

		if p.isIdentLikeStr("AT") {
			p.advance() // consume AT
			if p.cur.Type == '(' {
				p.advance()
				stmt.SplitValues = &nodes.List{}
				for {
					val, parseErr78 := p.parseExpr()
					if parseErr78 != nil {
						return nil, parseErr78
					}
					if val != nil {
						stmt.SplitValues.Items = append(stmt.SplitValues.Items, val)
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
		}
		// Optional INTO ( ... ) and parallel_clause - skip
		if p.cur.Type == kwINTO {
			p.advance()
			p.skipParenthesizedBlock()
		}
		if p.cur.Type == kwPARALLEL {
			p.advance()
			if p.cur.Type == tokICONST {
				stmt.Parallel = p.cur.Str
				p.advance()
			}
		} else if p.cur.Type == kwNOPARALLEL {
			stmt.NoParallel = true
			p.advance()
		}

	case p.isIdentLikeStr("ANNOTATIONS"):
		stmt.Action = "ANNOTATIONS"
		p.advance()
		p.skipParenthesizedBlock()

	default:
		// Consume remaining tokens for unknown actions
		for p.cur.Type != ';' && p.cur.Type != tokEOF {
			p.advance()
		}
	}

done:
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterIndexPhysicalAttrs parses remaining physical_attributes_clause
// options for ALTER INDEX (PCTFREE, PCTUSED, INITRANS, MAXTRANS, STORAGE, TABLESPACE, logging).
func (p *Parser) parseAlterIndexPhysicalAttrs(stmt *nodes.AlterIndexStmt) error {
	for {
		switch {
		case p.cur.Type == kwPCTFREE:
			p.advance()
			if p.cur.Type == tokICONST {
				stmt.PctFree = p.cur.Str
				p.advance()
			}
		case p.isIdentLikeStr("PCTUSED"):
			p.advance()
			if p.cur.Type == tokICONST {
				stmt.PctUsed = p.cur.Str
				p.advance()
			}
		case p.isIdentLikeStr("INITRANS"):
			p.advance()
			if p.cur.Type == tokICONST {
				stmt.InitTrans = p.cur.Str
				p.advance()
			}
		case p.isIdentLikeStr("MAXTRANS"):
			p.advance()
			if p.cur.Type == tokICONST {
				stmt.MaxTrans = p.cur.Str
				p.advance()
			}
		case p.isIdentLikeStr("STORAGE"):
			p.advance()
			p.skipParenthesizedBlock()
		case p.cur.Type == kwTABLESPACE:
			p.advance()
			var parseErr79 error
			stmt.Tablespace, parseErr79 = p.parseIdentifier()
			if parseErr79 != nil {
				return parseErr79
			}
		case p.cur.Type == kwLOGGING:
			stmt.Logging = true
			p.advance()
		case p.cur.Type == kwNOLOGGING:
			stmt.NoLogging = true
			p.advance()
		default:
			return nil
		}
	}
	return nil
}

// parseSizeClause parses a size clause like "10M", "100K", "1G", etc.
// Returns the combined string.
func (p *Parser) parseSizeClause() (string, error) {
	if p.cur.Type == tokICONST {
		size := p.cur.Str
		p.advance()
		// Optional unit suffix (K, M, G, T)
		if p.isIdentLike() {
			u := p.cur.Str
			if u == "K" || u == "M" || u == "G" || u == "T" {
				size += u
				p.advance()
			}
		}
		return size, nil
	}
	return "", nil
}

// parseAlterViewStmt parses an ALTER VIEW statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/ALTER-VIEW.html
//
//	ALTER VIEW [IF EXISTS] [schema.]view
//	{   COMPILE
//	  | ADD out_of_line_constraint
//	  | MODIFY CONSTRAINT constraint_name { RELY | NORELY }
//	  | DROP CONSTRAINT constraint_name
//	  | { READ ONLY | READ WRITE }
//	  | { EDITIONABLE | NONEDITIONABLE }
//	}
func (p *Parser) parseAlterViewStmt(start int) (nodes.StmtNode, error) {
	p.advance() // consume VIEW

	stmt := &nodes.AlterViewStmt{
		Loc: nodes.Loc{Start: start},
	}

	// IF EXISTS
	if p.cur.Type == kwIF {
		next := p.peekNext()
		if next.Type == kwEXISTS {
			stmt.IfExists = true
			p.advance() // consume IF
			p.advance() // consume EXISTS
		}
	}
	var parseErr80 error

	// Parse view name
	stmt.Name, parseErr80 = p.parseObjectName()
	if parseErr80 !=

		// Parse action
		nil {
		return nil, parseErr80
	}

	switch {
	case p.isIdentLikeStr("COMPILE"), p.isIdentLikeStr("RECOMPILE"):
		stmt.Action = "COMPILE"
		p.advance()

	case p.cur.Type == kwADD:
		stmt.Action = "ADD_CONSTRAINT"
		p.advance()
		var // consume ADD
		parseErr81 error
		stmt.Constraint, parseErr81 = p.parseTableConstraint()
		if parseErr81 != nil {
			return nil, parseErr81
		}

	case p.cur.Type == kwMODIFY:
		p.advance() // consume MODIFY
		if p.cur.Type == kwCONSTRAINT {
			stmt.Action = "MODIFY_CONSTRAINT"
			p.advance()
			var // consume CONSTRAINT
			parseErr82 error
			stmt.ConstraintName, parseErr82 = p.parseIdentifier()
			if parseErr82 != nil {
				return nil, parseErr82
			}
			if p.cur.Type == kwRELY {
				stmt.Rely = true
				p.advance()
			} else if p.isIdentLikeStr("NORELY") {
				stmt.NoRely = true
				p.advance()
			}
		} else {
			// skip unrecognized MODIFY clause
			for p.cur.Type != ';' && p.cur.Type != tokEOF {
				p.advance()
			}
		}

	case p.cur.Type == kwDROP:
		p.advance() // consume DROP
		if p.cur.Type == kwCONSTRAINT {
			stmt.Action = "DROP_CONSTRAINT"
			p.advance()
			var // consume CONSTRAINT
			parseErr83 error
			stmt.ConstraintName, parseErr83 = p.parseIdentifier()
			if parseErr83 != nil {

				// skip unrecognized DROP clause
				return nil, parseErr83
			}
		} else {

			for p.cur.Type != ';' && p.cur.Type != tokEOF {
				p.advance()
			}
		}

	case p.isIdentLike() && p.cur.Str == "ANNOTATIONS":
		stmt.Action = "ANNOTATIONS"
		p.advance() // consume ANNOTATIONS
		// skip annotations list if present
		if p.cur.Type == '(' {
			p.advance()
			depth := 1
			stmt.Annotations = &nodes.List{}
			for depth > 0 && p.cur.Type != tokEOF {
				if p.cur.Type == '(' {
					depth++
				} else if p.cur.Type == ')' {
					depth--
					if depth == 0 {
						break
					}
				}
				p.advance()
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}

	case p.cur.Type == kwREAD:
		p.advance() // consume READ
		if p.isIdentLikeStr("ONLY") {
			stmt.Action = "READ_ONLY"
			p.advance()
		} else if p.cur.Type == kwWRITE {
			stmt.Action = "READ_WRITE"
			p.advance()
		}

	case p.isIdentLikeStr("EDITIONABLE"):
		stmt.Action = "EDITIONABLE"
		p.advance()

	case p.isIdentLikeStr("NONEDITIONABLE"):
		stmt.Action = "NONEDITIONABLE"
		p.advance()

	default:
		p.skipToSemicolon()
	}

	// Skip optional trailing clauses (DISABLE NOVALIDATE on constraints, etc.)
	if stmt.Action == "ADD_CONSTRAINT" {
		// Consume DISABLE/ENABLE NOVALIDATE/VALIDATE after constraint
		for p.cur.Type != ';' && p.cur.Type != tokEOF {
			if p.cur.Type == kwDISABLE || p.cur.Type == kwENABLE || p.isIdentLikeStr("NOVALIDATE") || p.isIdentLikeStr("VALIDATE") {
				p.advance()
			} else {
				break
			}
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterSequenceStmt parses an ALTER SEQUENCE statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/ALTER-SEQUENCE.html
//
//	ALTER SEQUENCE [IF EXISTS] [schema.]sequence_name
//	  [ INCREMENT BY integer ]
//	  [ MAXVALUE integer | NOMAXVALUE ]
//	  [ MINVALUE integer | NOMINVALUE ]
//	  [ CYCLE | NOCYCLE ]
//	  [ CACHE integer | NOCACHE ]
//	  [ ORDER | NOORDER ]
//	  [ KEEP | NOKEEP ]
//	  [ RESTART [ WITH integer ] ]
//	  [ SCALE [ EXTEND | NOEXTEND ] | NOSCALE ]
//	  [ SHARD [ EXTEND | NOEXTEND ] | NOSHARD ]
//	  [ GLOBAL | SESSION ]
func (p *Parser) parseAlterSequenceStmt(start int) (nodes.StmtNode, error) {
	p.advance() // consume SEQUENCE

	stmt := &nodes.AlterSequenceStmt{
		Loc: nodes.Loc{Start: start},
	}

	// IF EXISTS
	if p.cur.Type == kwIF {
		next := p.peekNext()
		if next.Type == kwEXISTS {
			stmt.IfExists = true
			p.advance() // consume IF
			p.advance() // consume EXISTS
		}
	}
	var parseErr84 error

	// Parse sequence name
	stmt.Name, parseErr84 = p.parseObjectName()
	if parseErr84 !=

		// Parse sequence options (loop, multiple may be specified)
		nil {
		return nil, parseErr84
	}

	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		switch {
		case p.cur.Type == kwINCREMENT:
			p.advance() // consume INCREMENT
			if p.cur.Type == kwBY {
				p.advance() // consume BY
			}
			var parseErr85 error
			stmt.IncrementBy, parseErr85 = p.parseExpr()
			if parseErr85 != nil {
				return nil, parseErr85
			}

		case p.cur.Type == kwMAXVALUE:
			p.advance()
			var // consume MAXVALUE
			parseErr86 error
			stmt.MaxValue, parseErr86 = p.parseExpr()
			if parseErr86 != nil {
				return nil, parseErr86
			}

		case p.cur.Type == kwNOMAXVALUE:
			stmt.NoMaxValue = true
			p.advance()

		case p.cur.Type == kwMINVALUE:
			p.advance()
			var // consume MINVALUE
			parseErr87 error
			stmt.MinValue, parseErr87 = p.parseExpr()
			if parseErr87 != nil {
				return nil, parseErr87
			}

		case p.cur.Type == kwNOMINVALUE:
			stmt.NoMinValue = true
			p.advance()

		case p.cur.Type == kwCYCLE:
			stmt.Cycle = true
			p.advance()

		case p.cur.Type == kwNOCYCLE:
			stmt.NoCycle = true
			p.advance()

		case p.cur.Type == kwCACHE:
			p.advance()
			var // consume CACHE
			parseErr88 error
			stmt.Cache, parseErr88 = p.parseExpr()
			if parseErr88 != nil {
				return nil, parseErr88
			}

		case p.cur.Type == kwNOCACHE:
			stmt.NoCache = true
			p.advance()

		case p.cur.Type == kwORDER:
			stmt.Order = true
			p.advance()

		case p.cur.Type == kwNOORDER:
			stmt.NoOrder = true
			p.advance()

		case p.cur.Type == kwKEEP:
			stmt.Keep = true
			p.advance()

		case p.isIdentLikeStr("NOKEEP"):
			stmt.NoKeep = true
			p.advance()

		case p.isIdentLikeStr("RESTART"):
			stmt.Restart = true
			p.advance() // consume RESTART
			if p.cur.Type == kwWITH {
				p.advance()
				var // consume WITH
				parseErr89 error
				stmt.RestartWith, parseErr89 = p.parseExpr()
				if parseErr89 != nil {
					return nil, parseErr89
				}
			}

		case p.isIdentLikeStr("SCALE"):
			stmt.Scale = true
			p.advance() // consume SCALE
			if p.isIdentLikeStr("EXTEND") {
				stmt.ScaleExtend = true
				p.advance()
			} else if p.isIdentLikeStr("NOEXTEND") {
				stmt.ScaleNoExtend = true
				p.advance()
			}

		case p.isIdentLikeStr("NOSCALE"):
			stmt.NoScale = true
			p.advance()

		case p.isIdentLikeStr("SHARD"):
			stmt.Shard = true
			p.advance() // consume SHARD
			if p.isIdentLikeStr("EXTEND") {
				stmt.ShardExtend = true
				p.advance()
			} else if p.isIdentLikeStr("NOEXTEND") {
				stmt.ShardNoExtend = true
				p.advance()
			}

		case p.isIdentLikeStr("NOSHARD"):
			stmt.NoShard = true
			p.advance()

		case p.cur.Type == kwGLOBAL:
			stmt.Global = true
			p.advance()

		case p.cur.Type == kwSESSION:
			stmt.Session = true
			p.advance()

		default:
			goto done
		}
	}

done:
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseCompileClause parses COMPILE [DEBUG] [compiler_parameters_clause ...] [REUSE SETTINGS].
// Returns (debug, reuseSettings, compilerParams).
func (p *Parser) parseCompileClause() (bool, bool, []*nodes.SetParam, error) {
	var debug bool
	var reuseSettings bool
	var compilerParams []*nodes.SetParam

	// Optional DEBUG
	if p.isIdentLikeStr("DEBUG") {
		debug = true
		p.advance()
	}

	// Optional compiler_parameters_clause (name = value pairs) and REUSE SETTINGS
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		if p.isIdentLikeStr("REUSE") {
			p.advance() // consume REUSE
			if p.isIdentLikeStr("SETTINGS") {
				p.advance() // consume SETTINGS
			}
			reuseSettings = true
			break
		}
		// Check for compiler parameter: identifier = value
		if p.isIdentLike() && p.peekNext().Type == '=' {
			name, parseErr90 := p.parseIdentifier()
			if parseErr90 != nil {
				// consume '='
				return false, false, nil, parseErr90
			}
			p.advance()
			value, parseErr91 := p.parseExpr()
			if parseErr91 != nil {
				return false, false, nil, parseErr91
			}
			compilerParams = append(compilerParams, &nodes.SetParam{
				Name:  name,
				Value: value,
				Loc:   nodes.Loc{Start: p.pos(), End: p.prev.End},
			})
			continue
		}
		break
	}

	return debug, reuseSettings, compilerParams, nil
}

// parseAlterProcedureStmt parses an ALTER PROCEDURE statement.
//
// BNF: oracle/parser/bnf/ALTER-PROCEDURE.bnf
//
//	ALTER PROCEDURE [ IF EXISTS ] [ schema. ] procedure_name
//	    [ procedure_compile_clause ]
//	    [ { EDITIONABLE | NONEDITIONABLE } ] ;
//
//	procedure_compile_clause:
//	    COMPILE [ DEBUG ]
//	    [ compiler_parameters_clause ]
//	    [ REUSE SETTINGS ]
//
//	compiler_parameters_clause:
//	    parameter_name = parameter_value
func (p *Parser) parseAlterProcedureStmt(start int) (nodes.StmtNode, error) {
	p.advance() // consume PROCEDURE

	stmt := &nodes.AlterProcedureStmt{
		Loc: nodes.Loc{Start: start},
	}

	// IF EXISTS
	if p.cur.Type == kwIF {
		if p.peekNext().Type == kwEXISTS {
			stmt.IfExists = true
			p.advance() // consume IF
			p.advance() // consume EXISTS
		}
	}
	var parseErr92 error

	stmt.Name, parseErr92 = p.parseObjectName()
	if parseErr92 !=

		// Parse action
		nil {
		return nil, parseErr92
	}

	if p.isIdentLikeStr("COMPILE") {
		stmt.Compile = true
		p.advance()
		var // consume COMPILE
		parseErr93 error
		stmt.Debug, stmt.ReuseSettings, stmt.CompilerParams, parseErr93 = p.parseCompileClause()
		if parseErr93 !=

			// Trailing EDITIONABLE | NONEDITIONABLE
			nil {
			return nil, parseErr93
		}
	}

	if p.isIdentLikeStr("EDITIONABLE") {
		stmt.Editionable = true
		p.advance()
	} else if p.isIdentLikeStr("NONEDITIONABLE") {
		stmt.NonEditionable = true
		p.advance()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterFunctionStmt parses an ALTER FUNCTION statement.
//
// BNF: oracle/parser/bnf/ALTER-FUNCTION.bnf
//
//	ALTER FUNCTION [ IF EXISTS ] [ schema. ] function_name
//	    { function_compile_clause }
//	    [ EDITIONABLE | NONEDITIONABLE ] ;
//
//	function_compile_clause:
//	    COMPILE [ DEBUG ] [ compiler_parameters_clause ] [ REUSE SETTINGS ]
//
//	compiler_parameters_clause:
//	    parameter_name = parameter_value
func (p *Parser) parseAlterFunctionStmt(start int) (nodes.StmtNode, error) {
	p.advance() // consume FUNCTION

	stmt := &nodes.AlterFunctionStmt{
		Loc: nodes.Loc{Start: start},
	}

	// IF EXISTS
	if p.cur.Type == kwIF {
		if p.peekNext().Type == kwEXISTS {
			stmt.IfExists = true
			p.advance() // consume IF
			p.advance() // consume EXISTS
		}
	}
	var parseErr94 error

	stmt.Name, parseErr94 = p.parseObjectName()
	if parseErr94 !=

		// Parse action
		nil {
		return nil, parseErr94
	}

	if p.isIdentLikeStr("COMPILE") {
		stmt.Compile = true
		p.advance()
		var // consume COMPILE
		parseErr95 error
		stmt.Debug, stmt.ReuseSettings, stmt.CompilerParams, parseErr95 = p.parseCompileClause()
		if parseErr95 !=

			// Trailing EDITIONABLE | NONEDITIONABLE
			nil {
			return nil, parseErr95
		}
	}

	if p.isIdentLikeStr("EDITIONABLE") {
		stmt.Editionable = true
		p.advance()
	} else if p.isIdentLikeStr("NONEDITIONABLE") {
		stmt.NonEditionable = true
		p.advance()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterPackageStmt parses an ALTER PACKAGE statement.
//
// BNF: oracle/parser/bnf/ALTER-PACKAGE.bnf
//
//	ALTER PACKAGE [ IF EXISTS ] [ schema. ] package_name
//	    [ package_compile_clause ]
//	    [ { EDITIONABLE | NONEDITIONABLE } ] ;
//
//	package_compile_clause:
//	    COMPILE [ DEBUG ] [ PACKAGE | SPECIFICATION | BODY ]
//	    [ compiler_parameters_clause ]
//	    [ REUSE SETTINGS ]
//
//	compiler_parameters_clause:
//	    parameter_name = parameter_value
func (p *Parser) parseAlterPackageStmt(start int) (nodes.StmtNode, error) {
	p.advance() // consume PACKAGE

	stmt := &nodes.AlterPackageStmt{
		Loc: nodes.Loc{Start: start},
	}

	// IF EXISTS
	if p.cur.Type == kwIF {
		if p.peekNext().Type == kwEXISTS {
			stmt.IfExists = true
			p.advance() // consume IF
			p.advance() // consume EXISTS
		}
	}
	var parseErr96 error

	stmt.Name, parseErr96 = p.parseObjectName()
	if parseErr96 !=

		// Parse action
		nil {
		return nil, parseErr96
	}

	if p.isIdentLikeStr("COMPILE") {
		stmt.Compile = true
		p.advance() // consume COMPILE

		// Optional DEBUG (comes before target per BNF)
		if p.isIdentLikeStr("DEBUG") {
			stmt.Debug = true
			p.advance()
		}

		// Optional PACKAGE | BODY | SPECIFICATION
		switch {
		case p.cur.Type == kwPACKAGE:
			stmt.CompileTarget = "PACKAGE"
			p.advance()
		case p.cur.Type == kwBODY:
			stmt.CompileTarget = "BODY"
			p.advance()
		case p.isIdentLikeStr("SPECIFICATION"):
			stmt.CompileTarget = "SPECIFICATION"
			p.advance()
		}

		// Remaining compile clause (DEBUG if not yet consumed, compiler params, REUSE SETTINGS)
		debug, reuseSettings, compilerParams, parseErr97 := p.parseCompileClause()
		if parseErr97 != nil {
			return nil, parseErr97
		}
		if debug {
			stmt.Debug = true
		}
		stmt.ReuseSettings = reuseSettings
		stmt.CompilerParams = compilerParams
	}

	// Trailing EDITIONABLE | NONEDITIONABLE
	if p.isIdentLikeStr("EDITIONABLE") {
		stmt.Editionable = true
		p.advance()
	} else if p.isIdentLikeStr("NONEDITIONABLE") {
		stmt.NonEditionable = true
		p.advance()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterTriggerStmt parses an ALTER TRIGGER statement.
//
// BNF: oracle/parser/bnf/ALTER-TRIGGER.bnf
//
//	ALTER TRIGGER [ IF EXISTS ] [ schema. ] trigger_name
//	    { trigger_compile_clause
//	    | ENABLE
//	    | DISABLE
//	    | RENAME TO new_name
//	    | EDITIONABLE
//	    | NONEDITIONABLE
//	    } ;
//
//	trigger_compile_clause:
//	    COMPILE [ DEBUG ] [ compiler_parameters_clause [ compiler_parameters_clause ]... ] [ REUSE SETTINGS ]
//
//	compiler_parameters_clause:
//	    parameter_name = parameter_value
func (p *Parser) parseAlterTriggerStmt(start int) (nodes.StmtNode, error) {
	p.advance() // consume TRIGGER

	stmt := &nodes.AlterTriggerStmt{
		Loc: nodes.Loc{Start: start},
	}

	// IF EXISTS
	if p.cur.Type == kwIF {
		if p.peekNext().Type == kwEXISTS {
			stmt.IfExists = true
			p.advance() // consume IF
			p.advance() // consume EXISTS
		}
	}
	var parseErr98 error

	stmt.Name, parseErr98 = p.parseObjectName()
	if parseErr98 !=

		// Parse action
		nil {
		return nil, parseErr98
	}

	switch {
	case p.cur.Type == kwENABLE:
		stmt.Action = "ENABLE"
		p.advance()
	case p.cur.Type == kwDISABLE:
		stmt.Action = "DISABLE"
		p.advance()
	case p.isIdentLikeStr("COMPILE"):
		stmt.Action = "COMPILE"
		p.advance()
		var // consume COMPILE
		parseErr99 error
		stmt.Debug, stmt.ReuseSettings, stmt.CompilerParams, parseErr99 = p.parseCompileClause()
		if parseErr99 != nil {
			return nil, parseErr99
		}
	case p.cur.Type == kwRENAME:
		stmt.Action = "RENAME"
		p.advance() // consume RENAME
		if p.cur.Type == kwTO {
			p.advance() // consume TO
		}
		var parseErr100 error
		stmt.NewName, parseErr100 = p.parseIdentifier()
		if parseErr100 != nil {
			return nil, parseErr100
		}
	case p.isIdentLikeStr("EDITIONABLE"):
		stmt.Action = "EDITIONABLE"
		p.advance()
	case p.isIdentLikeStr("NONEDITIONABLE"):
		stmt.Action = "NONEDITIONABLE"
		p.advance()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterTypeStmt parses an ALTER TYPE statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/lnpls/ALTER-TYPE-statement.html
//
//	ALTER TYPE [IF EXISTS] [schema.]type_name
//	  { alter_type_clause | type_compile_clause }
//	  [ EDITIONABLE | NONEDITIONABLE ]
//
//	alter_type_clause:
//	    RESET
//	  | [NOT] INSTANTIABLE
//	  | [NOT] FINAL
//	  | ADD ATTRIBUTE ( attribute datatype [, ...] )
//	  | DROP ATTRIBUTE ( attribute [, ...] )
//	  | MODIFY ATTRIBUTE ( attribute datatype [, ...] )
//	  | ADD { MAP | ORDER } MEMBER FUNCTION ...
//	  | ADD { MEMBER | STATIC } { FUNCTION | PROCEDURE } ...
//	  | ADD CONSTRUCTOR FUNCTION ...
//	  | DROP { MAP | ORDER } MEMBER FUNCTION ...
//	  | DROP { MEMBER | STATIC } { FUNCTION | PROCEDURE } ...
//	  | MODIFY LIMIT integer
//	  | MODIFY ELEMENT TYPE datatype
//	  | dependent_handling_clause
//
//	type_compile_clause:
//	    COMPILE [SPECIFICATION | BODY] [DEBUG] [compiler_parameters_clause ...] [REUSE SETTINGS]
//
//	dependent_handling_clause:
//	    INVALIDATE
//	  | CASCADE [INCLUDING TABLE DATA | NOT INCLUDING TABLE DATA | CONVERT TO SUBSTITUTABLE]
//	    [FORCE]
func (p *Parser) parseAlterTypeStmt(start int) (nodes.StmtNode, error) {
	p.advance() // consume TYPE

	stmt := &nodes.AlterTypeStmt{
		Loc: nodes.Loc{Start: start},
	}

	// IF EXISTS
	if p.cur.Type == kwIF {
		if p.peekNext().Type == kwEXISTS {
			stmt.IfExists = true
			p.advance() // consume IF
			p.advance() // consume EXISTS
		}
	}
	var parseErr101 error

	stmt.Name, parseErr101 = p.parseObjectName()
	if parseErr101 !=

		// Parse action
		nil {
		return nil, parseErr101
	}

	switch {
	case p.isIdentLikeStr("COMPILE"):
		stmt.Action = "COMPILE"
		p.advance() // consume COMPILE
		// Optional DEBUG (comes before SPECIFICATION/BODY per BNF)
		if p.isIdentLikeStr("DEBUG") {
			stmt.Debug = true
			p.advance()
		}
		// Optional SPECIFICATION | BODY
		if p.isIdentLikeStr("SPECIFICATION") {
			stmt.CompileTarget = "SPECIFICATION"
			p.advance()
		} else if p.isIdentLikeStr("BODY") {
			stmt.CompileTarget = "BODY"
			p.advance()
		}
		// Remaining compile clause (DEBUG if not yet consumed, compiler params, REUSE SETTINGS)
		debug, reuseSettings, compilerParams, parseErr102 := p.parseCompileClause()
		if parseErr102 != nil {
			return nil, parseErr102
		}
		if debug {
			stmt.Debug = true
		}
		stmt.ReuseSettings = reuseSettings
		stmt.CompilerParams = compilerParams

	case p.isIdentLikeStr("REPLACE"):
		stmt.Action = "REPLACE"
		p.advance()
		parseErr103 := // consume REPLACE
			// REPLACE alter_type_replace_clause — parse the inline type spec
			p.parseAlterTypeReplaceClause(stmt)
		if parseErr103 != nil {
			return nil, parseErr103
		}

	case p.isIdentLikeStr("RESET"):
		stmt.Action = "RESET"
		p.advance()

	case p.cur.Type == kwNOT:
		p.advance() // consume NOT
		if p.isIdentLikeStr("INSTANTIABLE") {
			stmt.Action = "NOT_INSTANTIABLE"
			p.advance()
		} else if p.isIdentLikeStr("FINAL") {
			stmt.Action = "NOT_FINAL"
			p.advance()
		}
		parseErr104 := p.parseAlterTypeDependentHandling(stmt)
		if parseErr104 != nil {
			return nil, parseErr104
		}

	case p.isIdentLikeStr("INSTANTIABLE"):
		stmt.Action = "INSTANTIABLE"
		p.advance()
		parseErr105 := p.parseAlterTypeDependentHandling(stmt)
		if parseErr105 != nil {
			return nil, parseErr105
		}

	case p.isIdentLikeStr("FINAL"):
		stmt.Action = "FINAL"
		p.advance()
		parseErr106 := p.parseAlterTypeDependentHandling(stmt)
		if parseErr106 != nil {
			return nil, parseErr106
		}

	case p.cur.Type == kwADD:
		p.advance()
		parseErr107 := // consume ADD
			p.parseAlterTypeAddDrop(stmt, "ADD")
		if parseErr107 != nil {
			return nil, parseErr107
		}

	case p.cur.Type == kwDROP:
		p.advance()
		parseErr108 := // consume DROP
			p.parseAlterTypeAddDrop(stmt, "DROP")
		if parseErr108 != nil {
			return nil, parseErr108
		}

	case p.isIdentLikeStr("MODIFY"):
		p.advance() // consume MODIFY
		if p.isIdentLikeStr("ATTRIBUTE") {
			stmt.Action = "MODIFY_ATTRIBUTE"
			p.advance()
			var // consume ATTRIBUTE
			parseErr109 error
			stmt.Attributes, parseErr109 = p.parseAlterTypeAttributes(true)
			if parseErr109 != nil {
				return nil, parseErr109
			}
			parseErr110 := p.parseAlterTypeDependentHandling(stmt)
			if parseErr110 != nil {
				return nil, parseErr110
			}
		} else if p.isIdentLikeStr("LIMIT") {
			stmt.Action = "MODIFY_LIMIT"
			p.advance()
			var // consume LIMIT
			parseErr111 error
			stmt.LimitValue, parseErr111 = p.parseExpr()
			if parseErr111 != nil {
				return nil, parseErr111
			}
			parseErr112 := p.parseAlterTypeDependentHandling(stmt)
			if parseErr112 != nil {
				return nil, parseErr112
			}
		} else if p.isIdentLikeStr("ELEMENT") {
			stmt.Action = "MODIFY_ELEMENT_TYPE"
			p.advance() // consume ELEMENT
			if p.cur.Type == kwTYPE {
				p.advance() // consume TYPE
			}
			var parseErr113 error
			stmt.ElementType, parseErr113 = p.parseTypeName()
			if parseErr113 != nil {
				return nil, parseErr113
			}
			parseErr114 := p.parseAlterTypeDependentHandling(stmt)
			if parseErr114 != nil {

				// modifier_clause: NOT INSTANTIABLE | INSTANTIABLE | NOT FINAL | FINAL
				return nil, parseErr114
			}
		} else {

			if p.cur.Type == kwNOT {
				p.advance()
				if p.isIdentLikeStr("INSTANTIABLE") {
					stmt.Action = "MODIFY_NOT_INSTANTIABLE"
					p.advance()
				} else if p.isIdentLikeStr("FINAL") {
					stmt.Action = "MODIFY_NOT_FINAL"
					p.advance()
				}
			} else if p.isIdentLikeStr("INSTANTIABLE") {
				stmt.Action = "MODIFY_INSTANTIABLE"
				p.advance()
			} else if p.isIdentLikeStr("FINAL") {
				stmt.Action = "MODIFY_FINAL"
				p.advance()
			}
			parseErr115 := p.parseAlterTypeDependentHandling(stmt)
			if parseErr115 != nil {
				return nil, parseErr115
			}
		}

	case p.isIdentLikeStr("INVALIDATE"):
		stmt.Invalidate = true
		p.advance()

	case p.isIdentLikeStr("CASCADE"):
		parseErr116 := p.parseAlterTypeCascade(stmt)
		if parseErr116 != nil {
			return nil, parseErr116
		}

	case p.isIdentLikeStr("EDITIONABLE"):
		stmt.Action = "EDITIONABLE"
		stmt.Editionable = true
		p.advance()

	case p.isIdentLikeStr("NONEDITIONABLE"):
		stmt.Action = "NONEDITIONABLE"
		stmt.NonEditionable = true
		p.advance()
	}

	// Trailing EDITIONABLE / NONEDITIONABLE (after other clauses)
	if stmt.Action != "EDITIONABLE" && stmt.Action != "NONEDITIONABLE" {
		if p.isIdentLikeStr("EDITIONABLE") {
			stmt.Editionable = true
			p.advance()
		} else if p.isIdentLikeStr("NONEDITIONABLE") {
			stmt.NonEditionable = true
			p.advance()
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterTypeAddDrop handles ADD/DROP ATTRIBUTE or ADD/DROP method.
func (p *Parser) parseAlterTypeAddDrop(stmt *nodes.AlterTypeStmt, action string) error {
	switch {
	case p.isIdentLikeStr("ATTRIBUTE"):
		if action == "ADD" {
			stmt.Action = "ADD_ATTRIBUTE"
		} else {
			stmt.Action = "DROP_ATTRIBUTE"
		}
		p.advance() // consume ATTRIBUTE
		needDatatype := action == "ADD"
		var parseErr117 error
		stmt.Attributes, parseErr117 = p.parseAlterTypeAttributes(needDatatype)
		if parseErr117 != nil {
			return parseErr117
		}
		parseErr118 := p.parseAlterTypeDependentHandling(stmt)
		if parseErr118 != nil {
			return parseErr118
		}

	case p.isIdentLikeStr("MEMBER"), p.isIdentLikeStr("STATIC"):
		kind := p.cur.Str
		p.advance() // consume MEMBER/STATIC
		if action == "ADD" {
			stmt.Action = "ADD_METHOD"
		} else {
			stmt.Action = "DROP_METHOD"
		}
		stmt.MethodKind = kind
		parseErr119 := p.parseAlterTypeMethodSig(stmt)
		if parseErr119 != nil {
			return parseErr119
		}
		parseErr120 := p.parseAlterTypeDependentHandling(stmt)
		if parseErr120 != nil {
			return parseErr120
		}

	case p.isIdentLikeStr("MAP"), p.isIdentLikeStr("ORDER"):
		kind := p.cur.Str
		p.advance() // consume MAP/ORDER
		if p.isIdentLikeStr("MEMBER") {
			kind += " MEMBER"
			p.advance() // consume MEMBER
		}
		if action == "ADD" {
			stmt.Action = "ADD_METHOD"
		} else {
			stmt.Action = "DROP_METHOD"
		}
		stmt.MethodKind = kind
		parseErr121 := p.parseAlterTypeMethodSig(stmt)
		if parseErr121 != nil {
			return parseErr121
		}
		parseErr122 := p.parseAlterTypeDependentHandling(stmt)
		if parseErr122 != nil {
			return parseErr122
		}

	case p.isIdentLikeStr("CONSTRUCTOR"):
		p.advance() // consume CONSTRUCTOR
		if action == "ADD" {
			stmt.Action = "ADD_METHOD"
		} else {
			stmt.Action = "DROP_METHOD"
		}
		stmt.MethodKind = "CONSTRUCTOR"
		parseErr123 := p.parseAlterTypeMethodSig(stmt)
		if parseErr123 != nil {
			return parseErr123
		}
		parseErr124 := p.parseAlterTypeDependentHandling(stmt)
		if parseErr124 !=

			// parseAlterTypeMethodSig parses a method signature: FUNCTION/PROCEDURE name [(params)] [RETURN type].
			nil {
			return parseErr124
		}
	}
	return nil
}

func (p *Parser) parseAlterTypeMethodSig(stmt *nodes.AlterTypeStmt) error {
	// FUNCTION or PROCEDURE
	if p.cur.Type == kwFUNCTION {
		stmt.MethodType = "FUNCTION"
		p.advance()
	} else if p.cur.Type == kwPROCEDURE {
		stmt.MethodType = "PROCEDURE"
		p.advance()
	}
	var parseErr125 error

	// Method name
	stmt.MethodName, parseErr125 = p.parseIdentifier()
	if parseErr125 !=

		// Optional parameter list
		nil {
		return parseErr125
	}

	if p.cur.Type == '(' {
		params, parseErr126 := p.parseParameterList()
		if parseErr126 != nil {
			return parseErr126
		}
		if params != nil {
			for _, item := range params.Items {
				if param, ok := item.(*nodes.Parameter); ok {
					stmt.MethodParams = append(stmt.MethodParams, param)
				}
			}
		}
	}

	// Optional RETURN type
	if p.cur.Type == kwRETURN {
		p.advance() // consume RETURN
		// SELF AS RESULT
		if p.isIdentLikeStr("SELF") {
			p.advance() // consume SELF
			if p.cur.Type == kwAS {
				p.advance() // consume AS
			}
			if p.isIdentLikeStr("RESULT") {
				p.advance() // consume RESULT
			}
			stmt.MethodReturn = &nodes.TypeName{
				Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "SELF AS RESULT"}}},
			}
		} else {
			var parseErr127 error
			stmt.MethodReturn, parseErr127 = p.parseTypeName()
			if parseErr127 !=

				// parseAlterTypeAttributes parses ( attribute [datatype] [, ...] ).
				nil {
				return parseErr127
			}
		}
	}
	return nil
}

func (p *Parser) parseAlterTypeAttributes(withDatatype bool) ([]*nodes.TypeAttribute, error) {
	var attrs []*nodes.TypeAttribute

	if p.cur.Type == '(' {
		p.advance() // consume (
		for {
			attrStart := p.pos()
			parseValue1, parseErr2 := p.parseIdentifier()
			if parseErr2 != nil {
				return nil, parseErr2
			}
			attr := &nodes.TypeAttribute{
				Name: parseValue1,
				Loc:  nodes.Loc{Start: attrStart},
			}
			if withDatatype {
				var parseErr128 error
				attr.DataType, parseErr128 = p.parseTypeName()
				if parseErr128 != nil {
					return nil, parseErr128
				}
			}
			attr.Loc.End = p.prev.End
			attrs = append(attrs, attr)
			if p.cur.Type == ',' {
				p.advance() // consume ,
				continue
			}
			break
		}
		if p.cur.Type == ')' {
			p.advance() // consume )
		}
	}

	return attrs, nil
}

// parseAlterTypeDependentHandling parses optional INVALIDATE / CASCADE / FORCE.
func (p *Parser) parseAlterTypeDependentHandling(stmt *nodes.AlterTypeStmt) error {
	if p.isIdentLikeStr("INVALIDATE") {
		stmt.Invalidate = true
		p.advance()
		return nil
	}
	if p.isIdentLikeStr("CASCADE") {
		parseErr129 := p.parseAlterTypeCascade(stmt)
		if parseErr129 !=

			// parseAlterTypeCascade parses CASCADE [INCLUDING TABLE DATA | NOT INCLUDING TABLE DATA | CONVERT TO SUBSTITUTABLE] [FORCE].
			nil {
			return parseErr129
		}
	}
	return nil
}

func (p *Parser) parseAlterTypeCascade(stmt *nodes.AlterTypeStmt) error {
	stmt.Cascade = true
	p.advance() // consume CASCADE

	if p.isIdentLikeStr("INCLUDING") {
		p.advance() // consume INCLUDING
		if p.cur.Type == kwTABLE {
			p.advance() // consume TABLE
		}
		if p.isIdentLikeStr("DATA") {
			p.advance() // consume DATA
		}
		t := true
		stmt.IncludeData = &t
	} else if p.cur.Type == kwNOT {
		p.advance() // consume NOT
		if p.isIdentLikeStr("INCLUDING") {
			p.advance() // consume INCLUDING
		}
		if p.cur.Type == kwTABLE {
			p.advance() // consume TABLE
		}
		if p.isIdentLikeStr("DATA") {
			p.advance() // consume DATA
		}
		f := false
		stmt.IncludeData = &f
	} else if p.isIdentLikeStr("CONVERT") {
		p.advance() // consume CONVERT
		if p.cur.Type == kwTO {
			p.advance() // consume TO
		}
		if p.isIdentLikeStr("SUBSTITUTABLE") {
			p.advance() // consume SUBSTITUTABLE
		}
		stmt.ConvertToSubst = true
	}

	// Optional FORCE
	if p.isIdentLikeStr("FORCE") {
		stmt.Force = true
		p.advance()
	}
	return nil
}

// parseAlterTypeReplaceClause parses the REPLACE clause of ALTER TYPE.
// The REPLACE keyword has already been consumed.
//
// BNF: oracle/parser/bnf/ALTER-TYPE.bnf
//
//	REPLACE
//	    [ AUTHID { CURRENT_USER | DEFINER } ]
//	    [ ACCESSIBLE BY ( accessor [, accessor ]... ) ]
//	    AS OBJECT (
//	        attribute datatype [, attribute datatype ]...
//	        [, { MEMBER | STATIC } { FUNCTION | PROCEDURE } spec ]...
//	        [, { MAP | ORDER } MEMBER FUNCTION spec ]
//	        [, CONSTRUCTOR FUNCTION spec ]
//	    )
//	    [ [ NOT ] INSTANTIABLE ] [ [ NOT ] FINAL ]
func (p *Parser) parseAlterTypeReplaceClause(stmt *nodes.AlterTypeStmt) error {
	// Optional AUTHID clause
	if p.isIdentLikeStr("AUTHID") {
		p.advance() // AUTHID
		if p.isIdentLike() {
			p.advance() // CURRENT_USER or DEFINER
		}
	}

	// Optional ACCESSIBLE BY
	if p.isIdentLikeStr("ACCESSIBLE") {
		p.advance() // ACCESSIBLE
		if p.cur.Type == kwBY {
			p.advance() // BY
		}
		if p.cur.Type == '(' {
			p.advance()
			depth := 1
			for depth > 0 && p.cur.Type != tokEOF {
				if p.cur.Type == '(' {
					depth++
				} else if p.cur.Type == ')' {
					depth--
					if depth == 0 {
						break
					}
				}
				p.advance()
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	}

	// AS OBJECT ( ... ) or IS OBJECT ( ... )
	if p.cur.Type == kwAS || p.cur.Type == kwIS {
		p.advance() // AS/IS
	}

	if p.isIdentLikeStr("OBJECT") {
		p.advance() // OBJECT
	}

	if p.cur.Type == '(' {
		p.advance()
		// Parse attributes and method specs inside parentheses.
		// Attributes are: name datatype
		// Methods start with MEMBER, STATIC, MAP, ORDER, or CONSTRUCTOR
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			switch {
			case p.isIdentLikeStr("MEMBER") || p.isIdentLikeStr("STATIC") ||
				p.isIdentLikeStr("MAP") || p.isIdentLikeStr("ORDER") ||
				p.isIdentLikeStr("CONSTRUCTOR"):
				// Method spec — skip to next comma or closing paren at depth 0
				depth := 0
				for p.cur.Type != tokEOF {
					if p.cur.Type == '(' {
						depth++
					} else if p.cur.Type == ')' {
						if depth == 0 {
							break
						}
						depth--
					} else if p.cur.Type == ',' && depth == 0 {
						break
					}
					p.advance()
				}
			default:
				// Attribute: name datatype
				attrStart := p.pos()
				name, parseErr130 := p.parseIdentifier()
				if parseErr130 != nil {
					return parseErr130
				}
				if name == "" {
					p.advance()
					continue
				}
				typeName, parseErr131 := p.parseTypeName()
				if parseErr131 != nil {
					return parseErr131
				}
				attr := &nodes.TypeAttribute{
					Name:     name,
					DataType: typeName,
					Loc:      nodes.Loc{Start: attrStart, End: p.prev.End},
				}
				stmt.Attributes = append(stmt.Attributes, attr)
			}

			if p.cur.Type == ',' {
				p.advance()
			}
		}
		if p.cur.Type == ')' {
			p.advance()
		}
	}

	// Optional trailing modifiers: [NOT] INSTANTIABLE, [NOT] FINAL
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		switch {
		case p.cur.Type == kwNOT:
			p.advance()
			if p.isIdentLikeStr("INSTANTIABLE") || p.isIdentLikeStr("FINAL") {
				p.advance()
			}
		case p.isIdentLikeStr("INSTANTIABLE") || p.isIdentLikeStr("FINAL"):
			p.advance()
		default:
			return nil
		}
	}
	return nil
}

// skipToSemicolon advances until a semicolon or EOF is found.
// It does NOT consume the semicolon.
func (p *Parser) skipToSemicolon() {
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		p.advance()
	}
}
