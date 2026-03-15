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
func (p *Parser) parseAlterStmt() nodes.StmtNode {
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
			return p.parseAdminDDLStmt("ALTER", nodes.OBJECT_MATERIALIZED_ZONEMAP, start)
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
		return nil
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
		return p.parseAdminDDLStmt("ALTER", nodes.OBJECT_FLASHBACK_ARCHIVE, start)
	case kwUSER, kwROLE, kwPROFILE,
		kwTABLESPACE, kwCLUSTER, kwJAVA, kwLIBRARY:
		if adminStmt := p.parseAlterAdminObject(start); adminStmt != nil {
			return adminStmt
		}
		p.skipToSemicolon()
		return nil
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
				// Unknown ALTER SHARED target
				p.skipToSemicolon()
				return nil
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
			if adminStmt := p.parseAlterAdminObject(start); adminStmt != nil {
				return adminStmt
			}
		}
		// Unknown ALTER target — skip to semicolon or EOF.
		p.skipToSemicolon()
		return nil
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
func (p *Parser) parseAlterSessionStmt(start int) nodes.StmtNode {
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
		stmt.DBLink = p.parseIdentifier()

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
					stmt.Timeout = p.parseIntValue()
				}
			}
			// optional NAME 'string'
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
				stmt.ParallelDegree = p.parseIntValue()
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
		stmt.SetParams = p.parseSetParams()
	}

	stmt.Loc.End = p.pos()
	return stmt
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
func (p *Parser) parseAlterSystemStmt(start int) nodes.StmtNode {
	p.advance() // consume SYSTEM

	stmt := &nodes.AlterSystemStmt{
		Loc: nodes.Loc{Start: start},
	}

	switch {
	case p.cur.Type == kwSET:
		p.advance() // consume SET
		// Check for SET ENCRYPTION (security clause)
		if p.isIdentLikeStr("ENCRYPTION") {
			p.parseAlterSystemEncryption(stmt)
		} else {
			stmt.Action = "SET"
			stmt.SetParams = p.parseSetParams()
			// Parse optional SET modifiers: COMMENT, DEFERRED, SCOPE, SID, CONTAINER
			p.parseAlterSystemSetModifiers(stmt)
		}

	case p.isIdentLikeStr("RESET"):
		p.advance() // consume RESET
		stmt.Action = "RESET"
		stmt.ResetParam = p.parseIdentifier()
		// optional SCOPE
		if p.isIdentLikeStr("SCOPE") {
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			stmt.Scope = p.parseIdentifier()
		}
		// optional SID
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
				stmt.Timeout = p.parseIntValue()
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
		p.parseAlterSystemFlush(stmt)

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
		p.parseAlterSystemArchiveLog(stmt)

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
		p.skipToSemicolon()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseAlterSystemSetModifiers parses optional modifiers for ALTER SYSTEM SET:
// COMMENT, DEFERRED, SCOPE, SID, CONTAINER.
func (p *Parser) parseAlterSystemSetModifiers(stmt *nodes.AlterSystemStmt) {
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
			stmt.Scope = p.parseIdentifier()
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
			stmt.Container = p.parseIdentifier()
		default:
			return
		}
	}
}

// parseAlterSystemFlush parses the FLUSH sub-clause of ALTER SYSTEM.
func (p *Parser) parseAlterSystemFlush(stmt *nodes.AlterSystemStmt) {
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
		stmt.FlushRedoDB = p.parseIdentifier()
		// optional {NO CONFIRM APPLY | CONFIRM APPLY}
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
}

// parseAlterSystemArchiveLog parses the ARCHIVE LOG sub-clause of ALTER SYSTEM.
func (p *Parser) parseAlterSystemArchiveLog(stmt *nodes.AlterSystemStmt) {
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
			stmt.ArchiveThread = p.parseIntValue()
		}
	}

	// optional TO 'location'
	if p.cur.Type == kwTO {
		p.advance()
		if p.cur.Type == tokSCONST {
			stmt.ArchiveTo = p.cur.Str
			p.advance()
		}
	}
}

// parseAlterSystemEncryption parses the SET ENCRYPTION sub-clause of ALTER SYSTEM.
func (p *Parser) parseAlterSystemEncryption(stmt *nodes.AlterSystemStmt) {
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
				p.parseIdentifier() // consume password (not stored for security)
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
				p.parseIdentifier()
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
			p.parseIdentifier()
		}
	}
}

// parseSetParams parses one or more param = value pairs.
func (p *Parser) parseSetParams() *nodes.List {
	params := &nodes.List{}
	for {
		param := p.parseSetParam()
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
	return params
}

// parseSetParam parses a single name = value parameter setting.
func (p *Parser) parseSetParam() *nodes.SetParam {
	if !p.isIdentLike() {
		return nil
	}
	start := p.pos()
	name := p.parseIdentifier()
	if name == "" {
		return nil
	}

	// Expect '='
	if p.cur.Type != '=' {
		return &nodes.SetParam{
			Name: name,
			Loc:  nodes.Loc{Start: start, End: p.pos()},
		}
	}
	p.advance() // consume '='

	value := p.parseExpr()

	return &nodes.SetParam{
		Name:  name,
		Value: value,
		Loc:   nodes.Loc{Start: start, End: p.pos()},
	}
}

// parseAlterMaterializedViewStmt parses an ALTER MATERIALIZED VIEW statement.
// Called after ALTER MATERIALIZED VIEW has been consumed.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/ALTER-MATERIALIZED-VIEW.html
//
//	ALTER MATERIALIZED VIEW [IF EXISTS] [schema.]materialized_view
//	  { alter_mv_refresh
//	  | ENABLE QUERY REWRITE
//	  | DISABLE QUERY REWRITE
//	  | COMPILE
//	  | CONSIDER FRESH
//	  | { ENABLE | DISABLE } CONCURRENT REFRESH
//	  | SHRINK SPACE [ COMPACT | CASCADE ]
//	  | { CACHE | NOCACHE }
//	  | { PARALLEL [integer] | NOPARALLEL }
//	  | { LOGGING | NOLOGGING }
//	  | ... }
//
//	alter_mv_refresh:
//	  REFRESH
//	    [ FAST | COMPLETE | FORCE ]
//	    [ ON { COMMIT | DEMAND } ]
//	    [ START WITH date ]
//	    [ NEXT date ]
//	    [ WITH PRIMARY KEY ]
//	    [ USING ROLLBACK SEGMENT rollback_segment ]
//	    [ USING { ENFORCED | TRUSTED } CONSTRAINTS ]
//	    [ { ENABLE | DISABLE } ON QUERY COMPUTATION ]
func (p *Parser) parseAlterMaterializedViewStmt(start int) nodes.StmtNode {
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

	stmt.Name = p.parseObjectName()

	// Parse action clause
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
		p.advance() // consume REFRESH
		p.parseAlterMViewRefreshClause(stmt)

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
		} else {
			stmt.Action = "ENABLE_QUERY_REWRITE"
			p.skipToSemicolon()
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
		} else {
			stmt.Action = "DISABLE_QUERY_REWRITE"
			p.skipToSemicolon()
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

	default:
		// Fallback for unrecognized clauses
		p.skipToSemicolon()
	}

	stmt.Loc.End = p.pos()
	return stmt
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
func (p *Parser) parseAlterMViewRefreshClause(stmt *nodes.AlterMaterializedViewStmt) {
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
			stmt.StartWith = p.parseExpr()
		case p.cur.Type == kwNEXT:
			p.advance() // consume NEXT
			stmt.Next = p.parseExpr()
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
				return // not part of REFRESH clause, back out
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
				return // not part of REFRESH clause, back out
			}
		default:
			return
		}
	}
}

// parseAlterGeneric parses ALTER INDEX/VIEW/SEQUENCE/TABLE by consuming the
// object name and skipping the rest (simplified). Returns an AlterSessionStmt
// as a placeholder — in practice these would have their own AST types, but for
// now we skip the body to avoid blocking other work.
func (p *Parser) parseAlterGeneric(start int, objType nodes.ObjectType) nodes.StmtNode {
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

	// Parse the object name.
	stmt.Name = p.parseObjectName()

	// Skip remainder of the statement (clauses vary greatly by object type).
	p.skipToSemicolon()

	stmt.Loc.End = p.pos()
	return stmt
}

// parseAlterDatabaseLinkStmt parses an ALTER DATABASE LINK statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/ALTER-DATABASE-LINK.html
//
//	ALTER [ SHARED ] [ PUBLIC ] DATABASE LINK dblink_name
//	  CONNECT TO user IDENTIFIED BY password
//	  [ AUTHENTICATED BY user IDENTIFIED BY password ]
func (p *Parser) parseAlterDatabaseLinkStmt(start int, shared bool, public bool) nodes.StmtNode {
	p.advance() // consume DATABASE
	if p.cur.Type == kwLINK {
		p.advance() // consume LINK
	}

	stmt := &nodes.AlterDatabaseLinkStmt{
		Shared: shared,
		Public: public,
		Loc:    nodes.Loc{Start: start},
	}

	// Parse the database link name.
	stmt.Name = p.parseObjectName()

	// CONNECT TO user IDENTIFIED BY password
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

	stmt.Loc.End = p.pos()
	return stmt
}

// parseAlterSynonymStmt parses an ALTER SYNONYM statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/ALTER-SYNONYM.html
//
//	ALTER [ PUBLIC ] SYNONYM [ IF EXISTS ] [ schema. ] synonym
//	  { EDITIONABLE | NONEDITIONABLE | COMPILE }
func (p *Parser) parseAlterSynonymStmt(start int, public bool) nodes.StmtNode {
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

	// Parse the synonym name.
	stmt.Name = p.parseObjectName()

	// { EDITIONABLE | NONEDITIONABLE | COMPILE }
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

	stmt.Loc.End = p.pos()
	return stmt
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
func (p *Parser) parseAlterIndexStmt(start int) nodes.StmtNode {
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

	// Parse index name
	stmt.Name = p.parseObjectName()

	// Parse action
	switch {
	case p.isIdentLikeStr("REBUILD"):
		stmt.Action = "REBUILD"
		p.advance() // consume REBUILD
		// Optional PARTITION/SUBPARTITION and rebuild options
		for p.cur.Type != ';' && p.cur.Type != tokEOF {
			switch {
			case p.cur.Type == kwPARTITION:
				p.advance() // consume PARTITION
				stmt.Partition = p.parseIdentifier()
			case p.cur.Type == kwSUBPARTITION:
				p.advance() // consume SUBPARTITION
				stmt.Subpartition = p.parseIdentifier()
			case p.cur.Type == kwTABLESPACE:
				p.advance() // consume TABLESPACE
				stmt.Tablespace = p.parseIdentifier()
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
			p.advance() // consume PARTITION
			stmt.Partition = p.parseIdentifier()
			if p.cur.Type == kwTO {
				p.advance() // consume TO
			}
			stmt.NewName = p.parseIdentifier()
		} else if p.cur.Type == kwSUBPARTITION {
			stmt.Action = "RENAME_SUBPARTITION"
			p.advance() // consume SUBPARTITION
			stmt.Subpartition = p.parseIdentifier()
			if p.cur.Type == kwTO {
				p.advance() // consume TO
			}
			stmt.NewName = p.parseIdentifier()
		} else {
			stmt.Action = "RENAME"
			if p.cur.Type == kwTO {
				p.advance() // consume TO
			}
			stmt.NewName = p.parseIdentifier()
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
			p.advance() // consume KEEP
			stmt.DeallocateKeep = p.parseSizeClause()
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
		// May have more physical attributes
		p.parseAlterIndexPhysicalAttrs(stmt)

	case p.isIdentLikeStr("PCTUSED"):
		stmt.Action = "PHYSICAL_ATTRIBUTES"
		p.advance()
		if p.cur.Type == tokICONST {
			stmt.PctUsed = p.cur.Str
			p.advance()
		}
		p.parseAlterIndexPhysicalAttrs(stmt)

	case p.isIdentLikeStr("INITRANS"):
		stmt.Action = "PHYSICAL_ATTRIBUTES"
		p.advance()
		if p.cur.Type == tokICONST {
			stmt.InitTrans = p.cur.Str
			p.advance()
		}
		p.parseAlterIndexPhysicalAttrs(stmt)

	case p.isIdentLikeStr("MAXTRANS"):
		stmt.Action = "PHYSICAL_ATTRIBUTES"
		p.advance()
		if p.cur.Type == tokICONST {
			stmt.MaxTrans = p.cur.Str
			p.advance()
		}
		p.parseAlterIndexPhysicalAttrs(stmt)

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
				stmt.ModifyDefaultFor = p.parseIdentifier()
			}
			// Skip remaining options (physical_attributes, TABLESPACE, logging)
			p.parseAlterIndexPhysicalAttrs(stmt)
		} else if p.cur.Type == kwPARTITION {
			// modify_index_partition
			stmt.Action = "MODIFY_PARTITION"
			p.advance() // consume PARTITION
			stmt.Partition = p.parseIdentifier()
			// Sub-actions: skip remaining tokens
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
						stmt.DeallocateKeep = p.parseSizeClause()
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
			p.advance() // consume SUBPARTITION
			stmt.Subpartition = p.parseIdentifier()
			// Sub-actions: UNUSABLE, allocate_extent, deallocate_unused
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
					stmt.DeallocateKeep = p.parseSizeClause()
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
			stmt.AddPartitionName = p.parseIdentifier()
		}
		// Optional TABLESPACE, index_compression, parallel_clause
		for p.cur.Type != ';' && p.cur.Type != tokEOF {
			switch {
			case p.cur.Type == kwTABLESPACE:
				p.advance()
				stmt.Tablespace = p.parseIdentifier()
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
		stmt.Partition = p.parseIdentifier()

	case p.isIdentLikeStr("SPLIT"):
		// split_index_partition
		stmt.Action = "SPLIT_PARTITION"
		p.advance() // consume SPLIT
		if p.cur.Type == kwPARTITION {
			p.advance() // consume PARTITION
		}
		stmt.SplitPartition = p.parseIdentifier()
		// AT ( literal [, literal ]... )
		if p.isIdentLikeStr("AT") {
			p.advance() // consume AT
			if p.cur.Type == '(' {
				p.advance()
				stmt.SplitValues = &nodes.List{}
				for {
					val := p.parseExpr()
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
	stmt.Loc.End = p.pos()
	return stmt
}

// parseAlterIndexPhysicalAttrs parses remaining physical_attributes_clause
// options for ALTER INDEX (PCTFREE, PCTUSED, INITRANS, MAXTRANS, STORAGE, TABLESPACE, logging).
func (p *Parser) parseAlterIndexPhysicalAttrs(stmt *nodes.AlterIndexStmt) {
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
			stmt.Tablespace = p.parseIdentifier()
		case p.cur.Type == kwLOGGING:
			stmt.Logging = true
			p.advance()
		case p.cur.Type == kwNOLOGGING:
			stmt.NoLogging = true
			p.advance()
		default:
			return
		}
	}
}

// parseSizeClause parses a size clause like "10M", "100K", "1G", etc.
// Returns the combined string.
func (p *Parser) parseSizeClause() string {
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
		return size
	}
	return ""
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
func (p *Parser) parseAlterViewStmt(start int) nodes.StmtNode {
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

	// Parse view name
	stmt.Name = p.parseObjectName()

	// Parse action
	switch {
	case p.isIdentLikeStr("COMPILE"), p.isIdentLikeStr("RECOMPILE"):
		stmt.Action = "COMPILE"
		p.advance()

	case p.cur.Type == kwADD:
		stmt.Action = "ADD_CONSTRAINT"
		p.advance() // consume ADD
		stmt.Constraint = p.parseTableConstraint()

	case p.cur.Type == kwMODIFY:
		p.advance() // consume MODIFY
		if p.cur.Type == kwCONSTRAINT {
			stmt.Action = "MODIFY_CONSTRAINT"
			p.advance() // consume CONSTRAINT
			stmt.ConstraintName = p.parseIdentifier()
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
			p.advance() // consume CONSTRAINT
			stmt.ConstraintName = p.parseIdentifier()
		} else {
			// skip unrecognized DROP clause
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

	stmt.Loc.End = p.pos()
	return stmt
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
func (p *Parser) parseAlterSequenceStmt(start int) nodes.StmtNode {
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

	// Parse sequence name
	stmt.Name = p.parseObjectName()

	// Parse sequence options (loop, multiple may be specified)
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		switch {
		case p.cur.Type == kwINCREMENT:
			p.advance() // consume INCREMENT
			if p.cur.Type == kwBY {
				p.advance() // consume BY
			}
			stmt.IncrementBy = p.parseExpr()

		case p.cur.Type == kwMAXVALUE:
			p.advance() // consume MAXVALUE
			stmt.MaxValue = p.parseExpr()

		case p.cur.Type == kwNOMAXVALUE:
			stmt.NoMaxValue = true
			p.advance()

		case p.cur.Type == kwMINVALUE:
			p.advance() // consume MINVALUE
			stmt.MinValue = p.parseExpr()

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
			p.advance() // consume CACHE
			stmt.Cache = p.parseExpr()

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
				p.advance() // consume WITH
				stmt.RestartWith = p.parseExpr()
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
	stmt.Loc.End = p.pos()
	return stmt
}

// parseCompileClause parses COMPILE [DEBUG] [compiler_parameters_clause ...] [REUSE SETTINGS].
// Returns (debug, reuseSettings, compilerParams).
func (p *Parser) parseCompileClause() (bool, bool, []*nodes.SetParam) {
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
			name := p.parseIdentifier()
			p.advance() // consume '='
			value := p.parseExpr()
			compilerParams = append(compilerParams, &nodes.SetParam{
				Name:  name,
				Value: value,
				Loc:   nodes.Loc{Start: p.pos(), End: p.pos()},
			})
			continue
		}
		break
	}

	return debug, reuseSettings, compilerParams
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
func (p *Parser) parseAlterProcedureStmt(start int) nodes.StmtNode {
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

	stmt.Name = p.parseObjectName()

	// Parse action
	if p.isIdentLikeStr("COMPILE") {
		stmt.Compile = true
		p.advance() // consume COMPILE
		stmt.Debug, stmt.ReuseSettings, stmt.CompilerParams = p.parseCompileClause()
	}

	// Trailing EDITIONABLE | NONEDITIONABLE
	if p.isIdentLikeStr("EDITIONABLE") {
		stmt.Editionable = true
		p.advance()
	} else if p.isIdentLikeStr("NONEDITIONABLE") {
		stmt.NonEditionable = true
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt
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
func (p *Parser) parseAlterFunctionStmt(start int) nodes.StmtNode {
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

	stmt.Name = p.parseObjectName()

	// Parse action
	if p.isIdentLikeStr("COMPILE") {
		stmt.Compile = true
		p.advance() // consume COMPILE
		stmt.Debug, stmt.ReuseSettings, stmt.CompilerParams = p.parseCompileClause()
	}

	// Trailing EDITIONABLE | NONEDITIONABLE
	if p.isIdentLikeStr("EDITIONABLE") {
		stmt.Editionable = true
		p.advance()
	} else if p.isIdentLikeStr("NONEDITIONABLE") {
		stmt.NonEditionable = true
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt
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
func (p *Parser) parseAlterPackageStmt(start int) nodes.StmtNode {
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

	stmt.Name = p.parseObjectName()

	// Parse action
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
		debug, reuseSettings, compilerParams := p.parseCompileClause()
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

	stmt.Loc.End = p.pos()
	return stmt
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
func (p *Parser) parseAlterTriggerStmt(start int) nodes.StmtNode {
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

	stmt.Name = p.parseObjectName()

	// Parse action
	switch {
	case p.cur.Type == kwENABLE:
		stmt.Action = "ENABLE"
		p.advance()
	case p.cur.Type == kwDISABLE:
		stmt.Action = "DISABLE"
		p.advance()
	case p.isIdentLikeStr("COMPILE"):
		stmt.Action = "COMPILE"
		p.advance() // consume COMPILE
		stmt.Debug, stmt.ReuseSettings, stmt.CompilerParams = p.parseCompileClause()
	case p.cur.Type == kwRENAME:
		stmt.Action = "RENAME"
		p.advance() // consume RENAME
		if p.cur.Type == kwTO {
			p.advance() // consume TO
		}
		stmt.NewName = p.parseIdentifier()
	case p.isIdentLikeStr("EDITIONABLE"):
		stmt.Action = "EDITIONABLE"
		p.advance()
	case p.isIdentLikeStr("NONEDITIONABLE"):
		stmt.Action = "NONEDITIONABLE"
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt
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
func (p *Parser) parseAlterTypeStmt(start int) nodes.StmtNode {
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

	stmt.Name = p.parseObjectName()

	// Parse action
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
		debug, reuseSettings, compilerParams := p.parseCompileClause()
		if debug {
			stmt.Debug = true
		}
		stmt.ReuseSettings = reuseSettings
		stmt.CompilerParams = compilerParams

	case p.isIdentLikeStr("REPLACE"):
		stmt.Action = "REPLACE"
		p.advance() // consume REPLACE
		// REPLACE alter_type_replace_clause — skip the type body spec to semicolon
		// (the replace clause redefines the type spec inline)
		p.skipToSemicolon()

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
		p.parseAlterTypeDependentHandling(stmt)

	case p.isIdentLikeStr("INSTANTIABLE"):
		stmt.Action = "INSTANTIABLE"
		p.advance()
		p.parseAlterTypeDependentHandling(stmt)

	case p.isIdentLikeStr("FINAL"):
		stmt.Action = "FINAL"
		p.advance()
		p.parseAlterTypeDependentHandling(stmt)

	case p.cur.Type == kwADD:
		p.advance() // consume ADD
		p.parseAlterTypeAddDrop(stmt, "ADD")

	case p.cur.Type == kwDROP:
		p.advance() // consume DROP
		p.parseAlterTypeAddDrop(stmt, "DROP")

	case p.isIdentLikeStr("MODIFY"):
		p.advance() // consume MODIFY
		if p.isIdentLikeStr("ATTRIBUTE") {
			stmt.Action = "MODIFY_ATTRIBUTE"
			p.advance() // consume ATTRIBUTE
			stmt.Attributes = p.parseAlterTypeAttributes(true)
			p.parseAlterTypeDependentHandling(stmt)
		} else if p.isIdentLikeStr("LIMIT") {
			stmt.Action = "MODIFY_LIMIT"
			p.advance() // consume LIMIT
			stmt.LimitValue = p.parseExpr()
			p.parseAlterTypeDependentHandling(stmt)
		} else if p.isIdentLikeStr("ELEMENT") {
			stmt.Action = "MODIFY_ELEMENT_TYPE"
			p.advance() // consume ELEMENT
			if p.cur.Type == kwTYPE {
				p.advance() // consume TYPE
			}
			stmt.ElementType = p.parseTypeName()
			p.parseAlterTypeDependentHandling(stmt)
		} else {
			// modifier_clause: NOT INSTANTIABLE | INSTANTIABLE | NOT FINAL | FINAL
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
			p.parseAlterTypeDependentHandling(stmt)
		}

	case p.isIdentLikeStr("INVALIDATE"):
		stmt.Invalidate = true
		p.advance()

	case p.isIdentLikeStr("CASCADE"):
		p.parseAlterTypeCascade(stmt)

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

	stmt.Loc.End = p.pos()
	return stmt
}

// parseAlterTypeAddDrop handles ADD/DROP ATTRIBUTE or ADD/DROP method.
func (p *Parser) parseAlterTypeAddDrop(stmt *nodes.AlterTypeStmt, action string) {
	switch {
	case p.isIdentLikeStr("ATTRIBUTE"):
		if action == "ADD" {
			stmt.Action = "ADD_ATTRIBUTE"
		} else {
			stmt.Action = "DROP_ATTRIBUTE"
		}
		p.advance() // consume ATTRIBUTE
		needDatatype := action == "ADD"
		stmt.Attributes = p.parseAlterTypeAttributes(needDatatype)
		p.parseAlterTypeDependentHandling(stmt)

	case p.isIdentLikeStr("MEMBER"), p.isIdentLikeStr("STATIC"):
		kind := p.cur.Str
		p.advance() // consume MEMBER/STATIC
		if action == "ADD" {
			stmt.Action = "ADD_METHOD"
		} else {
			stmt.Action = "DROP_METHOD"
		}
		stmt.MethodKind = kind
		p.parseAlterTypeMethodSig(stmt)
		p.parseAlterTypeDependentHandling(stmt)

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
		p.parseAlterTypeMethodSig(stmt)
		p.parseAlterTypeDependentHandling(stmt)

	case p.isIdentLikeStr("CONSTRUCTOR"):
		p.advance() // consume CONSTRUCTOR
		if action == "ADD" {
			stmt.Action = "ADD_METHOD"
		} else {
			stmt.Action = "DROP_METHOD"
		}
		stmt.MethodKind = "CONSTRUCTOR"
		p.parseAlterTypeMethodSig(stmt)
		p.parseAlterTypeDependentHandling(stmt)
	}
}

// parseAlterTypeMethodSig parses a method signature: FUNCTION/PROCEDURE name [(params)] [RETURN type].
func (p *Parser) parseAlterTypeMethodSig(stmt *nodes.AlterTypeStmt) {
	// FUNCTION or PROCEDURE
	if p.cur.Type == kwFUNCTION {
		stmt.MethodType = "FUNCTION"
		p.advance()
	} else if p.cur.Type == kwPROCEDURE {
		stmt.MethodType = "PROCEDURE"
		p.advance()
	}

	// Method name
	stmt.MethodName = p.parseIdentifier()

	// Optional parameter list
	if p.cur.Type == '(' {
		params := p.parseParameterList()
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
			stmt.MethodReturn = p.parseTypeName()
		}
	}
}

// parseAlterTypeAttributes parses ( attribute [datatype] [, ...] ).
func (p *Parser) parseAlterTypeAttributes(withDatatype bool) []*nodes.TypeAttribute {
	var attrs []*nodes.TypeAttribute

	if p.cur.Type == '(' {
		p.advance() // consume (
		for {
			attrStart := p.pos()
			attr := &nodes.TypeAttribute{
				Name: p.parseIdentifier(),
				Loc:  nodes.Loc{Start: attrStart},
			}
			if withDatatype {
				attr.DataType = p.parseTypeName()
			}
			attr.Loc.End = p.pos()
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

	return attrs
}

// parseAlterTypeDependentHandling parses optional INVALIDATE / CASCADE / FORCE.
func (p *Parser) parseAlterTypeDependentHandling(stmt *nodes.AlterTypeStmt) {
	if p.isIdentLikeStr("INVALIDATE") {
		stmt.Invalidate = true
		p.advance()
		return
	}
	if p.isIdentLikeStr("CASCADE") {
		p.parseAlterTypeCascade(stmt)
	}
}

// parseAlterTypeCascade parses CASCADE [INCLUDING TABLE DATA | NOT INCLUDING TABLE DATA | CONVERT TO SUBSTITUTABLE] [FORCE].
func (p *Parser) parseAlterTypeCascade(stmt *nodes.AlterTypeStmt) {
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
}

// skipToSemicolon advances until a semicolon or EOF is found.
// It does NOT consume the semicolon.
func (p *Parser) skipToSemicolon() {
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		p.advance()
	}
}
