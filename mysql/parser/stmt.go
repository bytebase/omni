package parser

import (
	nodes "github.com/bytebase/omni/mysql/ast"
)

// parseStmt parses a single SQL statement.
func (p *Parser) parseStmt() (nodes.Node, error) {
	// Completion: check if cursor has been reached.
	p.checkCursor()

	if p.collectMode() {
		// Add all top-level statement-starting keywords as candidates.
		stmtTokens := []int{
			kwSELECT, kwINSERT, kwUPDATE, kwDELETE,
			kwCREATE, kwALTER, kwDROP, kwTRUNCATE, kwRENAME,
			kwWITH, kwTABLE, kwVALUES, kwREPLACE,
			kwSET, kwSHOW, kwUSE,
			kwDESCRIBE, kwDESC, kwEXPLAIN,
			kwBEGIN, kwSTART, kwCOMMIT, kwROLLBACK, kwSAVEPOINT, kwRELEASE,
			kwLOCK, kwUNLOCK,
			kwGRANT, kwREVOKE,
			kwLOAD, kwPREPARE, kwEXECUTE, kwDEALLOCATE,
			kwANALYZE, kwOPTIMIZE, kwCHECK, kwREPAIR,
			kwFLUSH, kwRESET, kwKILL, kwDO,
			kwXA, kwCALL, kwHELP,
		}
		for _, t := range stmtTokens {
			p.addTokenCandidate(t)
		}
		return nil, &ParseError{Message: "collecting"}
	}

	switch p.cur.Type {
	case kwWITH:
		return p.parseWithStmt()

	// TABLE and VALUES are query primaries, not statements of their own, so a
	// statement that begins with one is the query expression it begins: it may
	// be combined with UNION/INTERSECT/EXCEPT and carry the clauses of a whole
	// query expression (container-verified: "TABLE t1 UNION TABLE t2",
	// "VALUES ROW(1) UNION SELECT 1").
	case kwSELECT, kwTABLE, kwVALUES:
		return p.parseSelectStmt()

	case kwINSERT:
		return p.parseInsertStmt()

	case kwREPLACE:
		return p.parseReplaceStmt()

	case kwUPDATE:
		return p.parseUpdateStmt()

	case kwDELETE:
		return p.parseDeleteStmt()

	case kwCREATE:
		return p.parseCreateDispatch()

	case kwALTER:
		return p.parseAlterDispatch()

	case kwDROP:
		return p.parseDropDispatch()

	case kwTRUNCATE:
		return p.parseTruncateStmt()

	case kwRENAME:
		start := p.pos()
		p.advance() // consume RENAME
		if p.cur.Type == kwUSER {
			return p.parseRenameUserStmt(start)
		}
		return p.parseRenameTableStmt(start)

	case kwSET:
		return p.parseSetStmt()

	case kwSHOW:
		return p.parseShowStmt()

	case kwUSE:
		return p.parseUseStmt()

	case kwDESCRIBE, kwDESC, kwEXPLAIN:
		return p.parseExplainStmt()

	case kwBEGIN:
		return p.parseBeginStmt()

	case kwSTART:
		return p.parseStartDispatch()

	case kwCOMMIT:
		return p.parseCommitStmt()

	case kwROLLBACK:
		return p.parseRollbackStmt()

	case kwSAVEPOINT:
		return p.parseSavepointStmt()

	case kwRELEASE:
		return p.parseReleaseSavepointStmt()

	case kwLOCK:
		start := p.pos()
		p.advance() // consume LOCK
		if p.cur.Type == kwINSTANCE {
			return p.parseLockInstanceStmt(start)
		}
		return p.parseLockTablesStmt(start)

	case kwUNLOCK:
		start := p.pos()
		p.advance() // consume UNLOCK
		if p.cur.Type == kwINSTANCE {
			return p.parseUnlockInstanceStmt(start)
		}
		return p.parseUnlockTablesStmt(start)

	case kwGRANT:
		return p.parseGrantStmt()

	case kwREVOKE:
		return p.parseRevokeStmt()

	case kwLOAD:
		return p.parseLoadDispatch()

	case kwPREPARE:
		return p.parsePrepareStmt()

	case kwEXECUTE:
		return p.parseExecuteStmt()

	case kwDEALLOCATE:
		return p.parseDeallocateStmt()

	case kwANALYZE:
		return p.parseAnalyzeTableStmt()

	case kwOPTIMIZE:
		return p.parseOptimizeTableStmt()

	case kwCHECK:
		return p.parseCheckTableStmt()

	case kwREPAIR:
		return p.parseRepairTableStmt()

	case kwFLUSH:
		return p.parseFlushStmt()

	case kwRESET:
		return p.parseResetDispatch()

	case kwKILL:
		return p.parseKillStmt()

	case kwDO:
		return p.parseDoStmt()

	case kwXA:
		return p.parseXAStmt()

	case kwCALL:
		return p.parseCallStmt()

	case kwHANDLER:
		return p.parseHandlerStmt()

	case kwSIGNAL:
		return p.parseSignalStmt()

	case kwRESIGNAL:
		return p.parseResignalStmt()

	case kwGET:
		return p.parseGetDiagnosticsStmt()

	case kwCHANGE:
		return p.parseChangeDispatch()

	case kwSTOP:
		return p.parseStopDispatch()

	case kwPURGE:
		return p.parsePurgeBinaryLogsStmt()

	case kwIMPORT:
		return p.parseImportTableStmt()

	case kwBINLOG:
		return p.parseBinlogStmt()

	case kwCACHE:
		return p.parseCacheIndexStmt()

	case kwHELP:
		return p.parseHelpStmt()

	case kwCHECKSUM:
		return p.parseChecksumTableStmt()

	case kwSHUTDOWN:
		return p.parseShutdownStmt()

	case kwRESTART:
		return p.parseRestartStmt()

	case kwCLONE:
		return p.parseCloneStmt()

	case kwINSTALL:
		return p.parseInstallStmt()

	case kwUNINSTALL:
		return p.parseUninstallStmt()

	case '(':
		// Parenthesized subquery / select
		return p.parseSelectStmt()

	default:
		return nil, p.syntaxErrorAtCur()
	}
}

// parseChangeDispatch dispatches CHANGE statements to the appropriate parser.
// CHANGE is already the current token.
func (p *Parser) parseChangeDispatch() (nodes.Node, error) {
	// Don't consume CHANGE here; let the sub-parsers handle start position.
	// But we need to peek ahead: CHANGE REPLICATION {SOURCE|FILTER}
	// Since we have no Peek, consume CHANGE and REPLICATION, then check.
	start := p.pos()
	p.advance() // consume CHANGE

	if p.cur.Type == kwREPLICATION {
		p.advance() // consume REPLICATION
		switch p.cur.Type {
		case kwSOURCE:
			return p.parseChangeReplicationSourceStmtInner(start)
		case kwFILTER:
			return p.parseChangeReplicationFilterStmtInner(start)
		default:
			return nil, &ParseError{Message: "expected SOURCE or FILTER after CHANGE REPLICATION", Position: p.cur.Loc}
		}
	}

	// CHANGE MASTER TO (legacy alias for CHANGE REPLICATION SOURCE TO)
	if p.cur.Type == kwMASTER {
		p.advance() // consume MASTER
		return p.parseChangeMasterStmtInner(start)
	}

	return nil, &ParseError{Message: "unexpected token after CHANGE", Position: p.cur.Loc}
}

// parseStartDispatch dispatches START statements.
// START REPLICA/SLAVE -> replication; START TRANSACTION -> transaction.
func (p *Parser) parseStartDispatch() (nodes.Node, error) {
	start := p.pos()
	p.advance() // consume START

	if p.cur.Type == kwREPLICA || p.cur.Type == kwSLAVE {
		p.advance() // consume REPLICA or SLAVE
		return p.parseStartReplicaStmt(start)
	}

	if p.cur.Type == kwGROUP_REPLICATION {
		p.advance() // consume GROUP_REPLICATION
		return p.parseStartGroupReplicationStmt(start)
	}

	// Fall through to START TRANSACTION (parseBeginStmt expects p.cur is START or BEGIN)
	// We already consumed START, so we need to handle TRANSACTION here. start:
	// START_SYM TRANSACTION_SYM opt_start_transaction_option_list, TRANSACTION required.
	if _, err := p.expect(kwTRANSACTION); err != nil {
		return nil, err
	}

	// Completion: after START TRANSACTION, offer transaction characteristic keywords.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwWITH)
		p.addTokenCandidate(kwREAD)
		return nil, &ParseError{Message: "collecting"}
	}

	stmt := &nodes.BeginStmt{Loc: nodes.Loc{Start: start}}
	for {
		if p.cur.Type == kwWITH {
			p.advance()
			for _, kw := range []int{kwCONSISTENT, kwSNAPSHOT} {
				if _, err := p.expect(kw); err != nil {
					return nil, err
				}
			}
			stmt.WithConsistentSnapshot = true
		} else if p.cur.Type == kwREAD {
			p.advance()
			if _, ok := p.match(kwONLY); ok {
				stmt.ReadOnly = true
			} else if _, ok := p.match(kwWRITE); ok {
				stmt.ReadWrite = true
			} else {
				return nil, p.syntaxErrorAtCur()
			}
			// READ ONLY with READ WRITE is the action's syntax error (start, MySQL 8.0 sql/sql_yacc.yy).
			if stmt.ReadOnly && stmt.ReadWrite {
				return nil, p.syntaxErrorAtTok(p.prev)
			}
		} else {
			break
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseResetDispatch dispatches RESET statements.
// RESET REPLICA/SLAVE -> replication; RESET MASTER -> replication; else -> generic RESET.
func (p *Parser) parseResetDispatch() (nodes.Node, error) {
	start := p.pos()
	p.advance() // consume RESET

	if p.cur.Type == kwPERSIST {
		return p.parseResetPersistStmt(start)
	}

	// reset_options is a comma list of reset_option: SLAVE, REPLICA or MASTER, in any
	// order and repeated. One option is its own statement; a list is a FlushStmt that
	// keeps each option's text. mysql:8.0 runs `RESET SLAVE, SLAVE`.
	//
	// Ref: reset, reset_options and reset_option in MySQL 8.0 sql/sql_yacc.yy.
	var first nodes.Node
	var options []string
	for {
		optStart := p.pos()
		nodeStart := optStart
		if len(options) == 0 {
			nodeStart = start
		}
		var node nodes.Node
		var err error
		switch p.cur.Type {
		case kwREPLICA, kwSLAVE:
			node, err = p.parseResetReplicaStmt(nodeStart)
		case kwMASTER:
			node, err = p.parseResetMasterStmt(nodeStart)
		default:
			return nil, p.syntaxErrorAtCur()
		}
		if err != nil {
			return nil, err
		}
		if len(options) == 0 {
			first = node
		}
		options = append(options, p.inputText(optStart, p.prev.End))
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if len(options) == 1 {
		return first, nil
	}
	return &nodes.FlushStmt{Loc: nodes.Loc{Start: start, End: p.prev.End}, Options: options}, nil
}

// parseStopDispatch dispatches STOP statements.
// STOP REPLICA/SLAVE -> replication; STOP GROUP_REPLICATION -> group replication.
func (p *Parser) parseStopDispatch() (nodes.Node, error) {
	start := p.pos()
	p.advance() // consume STOP

	if p.cur.Type == kwREPLICA || p.cur.Type == kwSLAVE {
		p.advance() // consume REPLICA or SLAVE
		return p.parseStopReplicaStmt(start)
	}

	if p.cur.Type == kwGROUP_REPLICATION {
		p.advance() // consume GROUP_REPLICATION
		return p.parseStopGroupReplicationStmt(start)
	}

	return nil, &ParseError{Message: "expected REPLICA, SLAVE, or GROUP_REPLICATION after STOP", Position: p.cur.Loc}
}

// parseLoadDispatch dispatches LOAD statements.
// LOAD DATA/XML -> load data; LOAD INDEX INTO CACHE -> load index.
func (p *Parser) parseLoadDispatch() (nodes.Node, error) {
	start := p.pos()
	p.advance() // consume LOAD

	if p.cur.Type == kwINDEX {
		p.advance() // consume INDEX
		return p.parseLoadIndexIntoCacheStmt(start)
	}

	// LOAD DATA or LOAD XML (data_or_xml); any other word is a 1064.
	if p.cur.Type != kwDATA && p.cur.Type != kwXML {
		return nil, p.syntaxErrorAtCur()
	}
	return p.parseLoadDataStmt(start)
}

// parseCreateDispatch dispatches CREATE statements to the appropriate parser.
func (p *Parser) parseCreateDispatch() (nodes.Node, error) {
	start := p.pos()
	p.advance() // consume CREATE

	// Completion: after CREATE keyword, offer object type candidates.
	p.checkCursor()
	if p.collectMode() {
		for _, t := range []int{kwTABLE, kwINDEX, kwVIEW, kwDATABASE, kwFUNCTION, kwPROCEDURE, kwTRIGGER, kwEVENT} {
			p.addTokenCandidate(t)
		}
		return nil, &ParseError{Message: "collecting"}
	}

	if p.cur.Type == kwOR {
		return p.parseCreateOrReplace(start)
	}

	temporary := false

	// UNDO TABLESPACE
	if p.cur.Type == kwUNDO {
		p.advance()
		if p.cur.Type == kwTABLESPACE {
			return p.parseCreateTablespaceStmt(start, true)
		}
		return nil, &ParseError{Message: "expected TABLESPACE after UNDO", Position: p.cur.Loc}
	}

	// TEMPORARY
	if p.cur.Type == kwTEMPORARY {
		p.advance()
		temporary = true
	}

	switch p.cur.Type {
	case kwTABLE:
		return p.parseCreateTableStmt(start, temporary)

	case kwINDEX:
		return p.parseCreateIndexStmt(start, false, false, false)

	case kwUNIQUE:
		p.advance()
		if p.cur.Type == kwINDEX {
			return p.parseCreateIndexStmt(start, true, false, false)
		}
		return nil, &ParseError{Message: "expected INDEX after UNIQUE", Position: p.cur.Loc}

	case kwFULLTEXT:
		p.advance()
		// p.cur is now INDEX or KEY; parseCreateIndexStmt consumes it
		return p.parseCreateIndexStmt(start, false, true, false)

	case kwSPATIAL:
		p.advance()
		// Distinguish SPATIAL REFERENCE SYSTEM from SPATIAL INDEX
		if p.cur.Type == kwREFERENCE {
			return p.parseCreateSpatialRefSysStmt(start, false)
		}
		return p.parseCreateIndexStmt(start, false, false, true)

	case kwVIEW:
		return p.parseCreateViewStmt(start, false)

	case kwDATABASE, kwSCHEMA:
		return p.parseCreateDatabaseStmt(start)

	case kwAGGREGATE:
		// CREATE AGGREGATE FUNCTION — loadable UDF
		p.advance() // consume AGGREGATE
		if p.cur.Type == kwFUNCTION {
			return p.parseCreateLoadableFunction(start, true)
		}
		return nil, &ParseError{Message: "expected FUNCTION after AGGREGATE", Position: p.cur.Loc}

	case kwFUNCTION:
		return p.parseCreateFunctionStmt(start, false)

	case kwPROCEDURE:
		return p.parseCreateFunctionStmt(start, true)

	case kwTRIGGER:
		return p.parseCreateTriggerStmt(start)

	case kwEVENT:
		return p.parseCreateEventStmt(start)

	case kwUSER:
		return p.parseCreateUserStmt(start)

	case kwROLE:
		return p.parseCreateRoleStmt(start)

	case kwTABLESPACE:
		return p.parseCreateTablespaceStmt(start, false)

	case kwSERVER:
		return p.parseCreateServerStmt(start)

	case kwLOGFILE:
		return p.parseCreateLogfileGroupStmt(start)

	case kwALGORITHM:
		// Only CREATE VIEW uses ALGORITHM
		return p.parseCreateViewStmt(start, false)

	case kwSQL:
		// CREATE SQL SECURITY {DEFINER|INVOKER} VIEW — delegate to CREATE VIEW
		return p.parseCreateViewStmt(start, false)

	case kwDEFINER:
		// CREATE DEFINER=... {VIEW|FUNCTION|PROCEDURE|TRIGGER|EVENT}
		// Parse DEFINER = user[@host] and propagate to sub-parsers.
		definer, err := p.parseDefinerValue()
		if err != nil {
			return nil, err
		}
		// Optional SQL SECURITY {DEFINER|INVOKER} (view_suid)
		var sqlSecurity string
		if p.cur.Type == kwSQL {
			if sqlSecurity, err = p.parseViewSQLSecurity(); err != nil {
				return nil, err
			}
		}
		// Now dispatch on the object keyword
		switch p.cur.Type {
		case kwVIEW:
			stmt, err := p.parseCreateViewStmt(start, false)
			if err != nil {
				return nil, err
			}
			if stmt.Definer == "" {
				stmt.Definer = definer
			}
			if stmt.SqlSecurity == "" {
				stmt.SqlSecurity = sqlSecurity
			}
			return stmt, nil
		case kwFUNCTION:
			stmt, err := p.parseCreateFunctionStmt(start, false)
			if err != nil {
				return nil, err
			}
			if stmt.Definer == "" {
				stmt.Definer = definer
			}
			return stmt, nil
		case kwPROCEDURE:
			stmt, err := p.parseCreateFunctionStmt(start, true)
			if err != nil {
				return nil, err
			}
			if stmt.Definer == "" {
				stmt.Definer = definer
			}
			return stmt, nil
		case kwTRIGGER:
			stmt, err := p.parseCreateTriggerStmt(start)
			if err != nil {
				return nil, err
			}
			if stmt.Definer == "" {
				stmt.Definer = definer
			}
			return stmt, nil
		case kwEVENT:
			stmt, err := p.parseCreateEventStmt(start)
			if err != nil {
				return nil, err
			}
			if stmt.Definer == "" {
				stmt.Definer = definer
			}
			return stmt, nil
		default:
			return nil, &ParseError{Message: "unexpected token after DEFINER clause", Position: p.cur.Loc}
		}

	case kwRESOURCE:
		p.advance() // consume RESOURCE
		return p.parseCreateResourceGroupStmt(start)

	default:
		return nil, &ParseError{
			Message:  "unexpected token after CREATE",
			Position: p.cur.Loc,
		}
	}
}

// parseCreateOrReplace parses CREATE OR REPLACE. CREATE is consumed and p.cur
// is OR. MySQL 8.0 takes OR REPLACE only before a view and before SPATIAL
// REFERENCE SYSTEM; on any other object, and before TEMPORARY, the statement
// fails at the keyword that follows OR REPLACE.
//
// Ref: mysql-server 8.0 sql/sql_yacc.yy — view_replace_or_algorithm (the view
// head: OR REPLACE [ALGORITHM = ...] [DEFINER = ...] [SQL SECURITY ...] VIEW)
// and the CREATE OR_SYM REPLACE_SYM SPATIAL_SYM REFERENCE_SYM SYSTEM_SYM
// branch of create, which takes no IF NOT EXISTS.
// https://dev.mysql.com/doc/refman/8.0/en/create-view.html
// https://dev.mysql.com/doc/refman/8.0/en/create-spatial-reference-system.html
func (p *Parser) parseCreateOrReplace(start int) (nodes.Node, error) {
	p.advance() // consume OR
	if p.cur.Type != kwREPLACE {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume REPLACE

	switch p.cur.Type {
	case kwVIEW, kwALGORITHM, kwDEFINER, kwSQL:
		// parseCreateViewStmt reads ALGORITHM, DEFINER, and SQL SECURITY in
		// that order and fails at whatever stands where VIEW belongs.
		return p.parseCreateViewStmt(start, true)
	case kwSPATIAL:
		p.advance() // consume SPATIAL
		if p.cur.Type != kwREFERENCE {
			return nil, p.syntaxErrorAtCur()
		}
		return p.parseCreateSpatialRefSysStmt(start, true)
	default:
		return nil, p.syntaxErrorAtCur()
	}
}

// parseAlterDispatch dispatches ALTER statements to the appropriate parser.
func (p *Parser) parseAlterDispatch() (nodes.Node, error) {
	start := p.pos()
	p.advance() // consume ALTER

	// Completion: after ALTER keyword, offer object type candidates.
	p.checkCursor()
	if p.collectMode() {
		for _, t := range []int{kwTABLE, kwDATABASE, kwVIEW, kwFUNCTION, kwPROCEDURE, kwEVENT} {
			p.addTokenCandidate(t)
		}
		return nil, &ParseError{Message: "collecting"}
	}

	// UNDO TABLESPACE
	if p.cur.Type == kwUNDO {
		p.advance()
		if p.cur.Type == kwTABLESPACE {
			return p.parseAlterTablespaceStmt(start, true)
		}
		return nil, &ParseError{Message: "expected TABLESPACE after UNDO", Position: p.cur.Loc}
	}

	switch p.cur.Type {
	case kwTABLE:
		return p.parseAlterTableStmt(start)

	case kwDATABASE, kwSCHEMA:
		return p.parseAlterDatabaseStmt(start)

	case kwUSER:
		return p.parseAlterUserStmt(start)

	case kwTABLESPACE:
		return p.parseAlterTablespaceStmt(start, false)

	case kwSERVER:
		return p.parseAlterServerStmt(start)

	case kwVIEW:
		return p.parseAlterViewStmt(start)

	case kwALGORITHM:
		// ALTER ALGORITHM = ... VIEW — delegate to ALTER VIEW parser
		return p.parseAlterViewStmt(start)

	case kwSQL:
		// ALTER SQL SECURITY {DEFINER|INVOKER} VIEW — delegate to ALTER VIEW
		return p.parseAlterViewStmt(start)

	case kwDEFINER:
		// ALTER DEFINER = user {VIEW|EVENT}
		// Parse DEFINER = user[@host] and propagate to sub-parsers.
		definer, err := p.parseDefinerValue()
		if err != nil {
			return nil, err
		}
		// Optional SQL SECURITY {DEFINER|INVOKER} (view_suid)
		var sqlSecurity string
		if p.cur.Type == kwSQL {
			if sqlSecurity, err = p.parseViewSQLSecurity(); err != nil {
				return nil, err
			}
		}
		switch p.cur.Type {
		case kwVIEW:
			stmt, err := p.parseAlterViewStmt(start)
			if err != nil {
				return nil, err
			}
			if stmt.Definer == "" {
				stmt.Definer = definer
			}
			if stmt.SqlSecurity == "" {
				stmt.SqlSecurity = sqlSecurity
			}
			return stmt, nil
		case kwEVENT:
			stmt, err := p.parseAlterEventStmt(start)
			if err != nil {
				return nil, err
			}
			if stmt.Definer == "" {
				stmt.Definer = definer
			}
			return stmt, nil
		default:
			return nil, &ParseError{Message: "unexpected token after DEFINER clause in ALTER", Position: p.cur.Loc}
		}

	case kwEVENT:
		return p.parseAlterEventStmt(start)

	case kwFUNCTION:
		return p.parseAlterRoutineStmt(start, false)

	case kwPROCEDURE:
		return p.parseAlterRoutineStmt(start, true)

	case kwLOGFILE:
		return p.parseAlterLogfileGroupStmt(start)

	case kwRESOURCE:
		p.advance() // consume RESOURCE
		return p.parseAlterResourceGroupStmt(start)

	case kwINSTANCE:
		return p.parseAlterInstanceStmt(start)

	default:
		return nil, &ParseError{
			Message:  "unexpected token after ALTER",
			Position: p.cur.Loc,
		}
	}
}

// parseDropDispatch dispatches DROP statements to the appropriate parser.
func (p *Parser) parseDropDispatch() (nodes.Node, error) {
	start := p.pos()
	p.advance() // consume DROP

	// Completion: after DROP keyword, offer object type candidates.
	p.checkCursor()
	if p.collectMode() {
		for _, t := range []int{kwTABLE, kwINDEX, kwVIEW, kwDATABASE, kwFUNCTION, kwPROCEDURE, kwTRIGGER, kwEVENT, kwIF} {
			p.addTokenCandidate(t)
		}
		return nil, &ParseError{Message: "collecting"}
	}

	// UNDO TABLESPACE
	if p.cur.Type == kwUNDO {
		p.advance()
		if p.cur.Type == kwTABLESPACE {
			return p.parseDropTablespaceStmt(start, true)
		}
		return nil, &ParseError{Message: "expected TABLESPACE after UNDO", Position: p.cur.Loc}
	}

	temporary := false
	if p.cur.Type == kwTEMPORARY {
		p.advance()
		temporary = true
	}

	switch p.cur.Type {
	case kwTABLE:
		return p.parseDropTableStmt(start, temporary)

	case kwINDEX:
		return p.parseDropIndexStmt(start)

	case kwVIEW:
		return p.parseDropViewStmt(start)

	case kwDATABASE, kwSCHEMA:
		return p.parseDropDatabaseStmt(start)

	case kwUSER:
		return p.parseDropUserStmt(start)

	case kwROLE:
		return p.parseDropRoleStmt(start)

	case kwFUNCTION:
		return p.parseDropRoutineStmt(start, false)

	case kwPROCEDURE:
		return p.parseDropRoutineStmt(start, true)

	case kwTRIGGER:
		return p.parseDropTriggerStmt(start)

	case kwEVENT:
		return p.parseDropEventStmt(start)

	case kwTABLESPACE:
		return p.parseDropTablespaceStmt(start, false)

	case kwSERVER:
		return p.parseDropServerStmt(start)

	case kwLOGFILE:
		return p.parseDropLogfileGroupStmt(start)

	case kwSPATIAL:
		p.advance() // consume SPATIAL
		return p.parseDropSpatialRefSysStmt(start)

	case kwRESOURCE:
		p.advance() // consume RESOURCE
		return p.parseDropResourceGroupStmt(start)

	case kwPREPARE:
		// DROP PREPARE stmt_name (alias for DEALLOCATE PREPARE)
		p.advance() // consume PREPARE
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		return &nodes.DeallocateStmt{
			Loc:  nodes.Loc{Start: start, End: p.prev.End},
			Name: name,
		}, nil

	default:
		return nil, &ParseError{
			Message:  "unexpected token after DROP",
			Position: p.cur.Loc,
		}
	}
}

// parseDefinerValue parses the DEFINER = user[@host] clause and returns the
// definer string. Handles: 'user'@'host', user@host, CURRENT_USER, CURRENT_USER().
// The current token must be DEFINER; on return, the current token is the one
// after the definer clause.
//
// definer: DEFINER_SYM EQ user (sql_yacc.yy), so `=` and the user are required.
func (p *Parser) parseDefinerValue() (string, error) {
	p.advance() // consume DEFINER
	if _, err := p.expect('='); err != nil {
		return "", err
	}

	// Completion: after DEFINER=, offer CURRENT_USER keyword.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwCURRENT_USER)
		return "", &ParseError{Message: "collecting"}
	}

	// Parse user part: identifier, string literal, or CURRENT_USER
	var user string
	isCurUser := p.cur.Type == kwCURRENT_USER
	if p.cur.Type == tokSCONST {
		user = "'" + p.cur.Str + "'"
		p.advance()
	} else if isCurUser {
		user = "CURRENT_USER"
		p.advance()
	} else if p.isIdentToken() {
		user, _, _ = p.parseIdentifier()
	} else {
		return "", p.syntaxErrorAtCur()
	}

	// Handle CURRENT_USER()
	if isCurUser && p.cur.Type == '(' {
		p.advance() // consume (
		if _, err := p.expect(')'); err != nil {
			return "", err
		}
		return "CURRENT_USER()", nil
	}

	// Handle @host: the lexer scans @host as a single tokIDENT with Str="@host",
	// or for quoted form 'user'@'host' the @ is also part of a variable token.
	// Check for tokIDENT starting with "@".
	if p.cur.Type == tokIDENT && len(p.cur.Str) > 0 && p.cur.Str[0] == '@' {
		atHost := p.cur.Str // e.g. "@localhost" or "@"
		p.advance()
		if atHost == "@" {
			// Standalone @ followed by host identifier or string
			var host string
			if p.cur.Type == tokSCONST {
				host = "'" + p.cur.Str + "'"
				p.advance()
			} else if p.isIdentToken() {
				host, _, _ = p.parseIdentifier()
			}
			return user + "@" + host, nil
		}
		// @host as single token: strip @ prefix
		host := atHost[1:]
		return user + "@" + host, nil
	}

	return user, nil
}
