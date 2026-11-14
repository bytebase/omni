package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
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
			kwWITH, kwVALUES, kwREPLACE,
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

	case kwSELECT, kwVALUES:
		return p.parseQueryStatement()

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
		if p.analyzeTakesStatement() {
			return p.parseAnalyzeStatement()
		}
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

	case kwCLONE:
		return p.parseCloneStmt()

	case kwINSTALL:
		return p.parseInstallStmt()

	case kwUNINSTALL:
		return p.parseUninstallStmt()

	case '(':
		// Parenthesized subquery / select
		return p.parseQueryStatement()

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
		return p.parseStartReplicaStmt(start, false)
	}
	if p.atAllSlaves() {
		return p.parseStartReplicaStmt(start, true)
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
			// READ ONLY with READ WRITE is the action's syntax error (start, MariaDB 11.8 sql/sql_yacc.yy).
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

	if p.cur.Type == kwREPLICA || p.cur.Type == kwSLAVE {
		return p.parseResetReplicaStmt(start)
	}
	if p.cur.Type == kwMASTER {
		return p.parseResetMasterStmt(start)
	}
	if p.cur.Type == kwPERSIST {
		return p.parseResetPersistStmt(start)
	}

	// reset_option's last arm is QUERY CACHE, both words, and reset_options
	// is a comma list of reset_option (reset in MariaDB 11.8 sql/sql_yacc.yy).
	// A list after QUERY CACHE keeps each later option's text.
	if p.cur.Type != kwQUERY {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	if _, err := p.expect(kwCACHE); err != nil {
		return nil, err
	}
	stmt := &nodes.FlushStmt{Loc: nodes.Loc{Start: start}, Options: []string{"QUERY CACHE"}}
	for {
		if _, ok := p.match(','); !ok {
			break
		}
		optStart := p.pos()
		switch p.cur.Type {
		case kwQUERY:
			p.advance()
			if _, err := p.expect(kwCACHE); err != nil {
				return nil, err
			}
		case kwMASTER:
			if _, err := p.parseResetMasterStmt(optStart); err != nil {
				return nil, err
			}
		case kwREPLICA, kwSLAVE:
			if _, err := p.parseResetReplicaStmt(optStart); err != nil {
				return nil, err
			}
		default:
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Options = append(stmt.Options, p.inputText(optStart, p.prev.End))
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseStopDispatch dispatches STOP statements.
// STOP REPLICA/SLAVE -> replication; STOP GROUP_REPLICATION -> group replication.
func (p *Parser) parseStopDispatch() (nodes.Node, error) {
	start := p.pos()
	p.advance() // consume STOP

	if p.cur.Type == kwREPLICA || p.cur.Type == kwSLAVE {
		p.advance() // consume REPLICA or SLAVE
		return p.parseStopReplicaStmt(start, false)
	}
	if p.atAllSlaves() {
		return p.parseStopReplicaStmt(start, true)
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
		for _, t := range []int{kwTABLE, kwSEQUENCE, kwINDEX, kwVIEW, kwDATABASE, kwFUNCTION, kwPROCEDURE, kwTRIGGER, kwEVENT} {
			p.addTokenCandidate(t)
		}
		return nil, &ParseError{Message: "collecting"}
	}

	orReplace := false
	temporary := false

	// OR REPLACE — OR must be followed by REPLACE; MariaDB rejects a bare
	// CREATE OR <object> (without REPLACE) at parse (1064). In sql_yacc.yy,
	// create_or_replace opens the TABLE, SEQUENCE, INDEX, DATABASE, VIEW,
	// FUNCTION, PROCEDURE, TRIGGER, EVENT, USER, ROLE, SERVER, and PACKAGE
	// (not parsed here) alternatives of the create rule; each node below
	// records it. MariaDB 11.8 has no TABLESPACE, UNDO TABLESPACE, LOGFILE
	// GROUP, RESOURCE GROUP, or SPATIAL REFERENCE SYSTEM statement; this parser
	// keeps their plain MySQL forms but rejects OR REPLACE on them at the
	// object keyword, as MariaDB does (1064).
	if p.cur.Type == kwOR {
		p.advance()
		if _, err := p.expect(kwREPLACE); err != nil {
			return nil, err
		}
		orReplace = true
	}

	// TEMPORARY — opt_temporary precedes only TABLE and SEQUENCE in the create
	// rule; MariaDB rejects it before any other object at parse (1064).
	if p.cur.Type == kwTEMPORARY {
		p.advance()
		if p.cur.Type != kwTABLE && p.cur.Type != kwSEQUENCE {
			return nil, p.syntaxErrorAtCur()
		}
		temporary = true
	}

	switch p.cur.Type {
	case kwTABLE:
		return p.parseCreateTableStmt(start, orReplace, temporary)

	case kwSEQUENCE:
		return p.parseCreateSequenceStmt(start, orReplace, temporary)

	case kwINDEX:
		return p.parseCreateIndexStmt(start, orReplace, false, false, false)

	case kwUNIQUE:
		p.advance()
		if p.cur.Type == kwINDEX {
			return p.parseCreateIndexStmt(start, orReplace, true, false, false)
		}
		return nil, &ParseError{Message: "expected INDEX after UNIQUE", Position: p.cur.Loc}

	case kwFULLTEXT:
		p.advance()
		// p.cur is now INDEX or KEY; parseCreateIndexStmt consumes it
		return p.parseCreateIndexStmt(start, orReplace, false, true, false)

	case kwVECTOR:
		return p.parseCreateVectorIndexStmt(start, orReplace)

	case kwSPATIAL:
		p.advance()
		// Distinguish SPATIAL REFERENCE SYSTEM from SPATIAL INDEX
		if p.cur.Type == kwREFERENCE {
			if orReplace {
				return nil, p.syntaxErrorAtCur()
			}
			return p.parseCreateSpatialRefSysStmt(start)
		}
		return p.parseCreateIndexStmt(start, orReplace, false, false, true)

	case kwVIEW:
		return p.parseCreateViewStmt(start, orReplace)

	case kwDATABASE, kwSCHEMA:
		return p.parseCreateDatabaseStmt(start, orReplace)

	case kwAGGREGATE:
		// CREATE AGGREGATE FUNCTION — loadable UDF
		p.advance() // consume AGGREGATE
		if p.cur.Type == kwFUNCTION {
			return p.parseCreateLoadableFunction(start, orReplace, true)
		}
		return nil, &ParseError{Message: "expected FUNCTION after AGGREGATE", Position: p.cur.Loc}

	case kwFUNCTION:
		return p.parseCreateFunctionStmt(start, orReplace, false)

	case kwPROCEDURE:
		return p.parseCreateFunctionStmt(start, orReplace, true)

	case kwTRIGGER:
		return p.parseCreateTriggerStmt(start, orReplace)

	case kwEVENT:
		return p.parseCreateEventStmt(start, orReplace)

	case kwUSER:
		return p.parseCreateUserStmt(start, orReplace)

	case kwROLE:
		return p.parseCreateRoleStmt(start, orReplace)

	case kwSERVER:
		return p.parseCreateServerStmt(start, orReplace)

	case kwALGORITHM:
		// Only CREATE VIEW uses ALGORITHM
		return p.parseCreateViewStmt(start, orReplace)

	case kwSQL:
		// CREATE SQL SECURITY {DEFINER|INVOKER} VIEW — delegate to CREATE VIEW
		return p.parseCreateViewStmt(start, orReplace)

	case kwDEFINER:
		// CREATE DEFINER=... {VIEW|FUNCTION|PROCEDURE|TRIGGER|EVENT}
		// Parse DEFINER = user[@host] and propagate to sub-parsers.
		definer, err := p.parseDefinerValue()
		if err != nil {
			return nil, err
		}
		// Optional SQL SECURITY {DEFINER|INVOKER} (view_suid)
		var sqlSecurity string
		sawSQLSecurity := p.cur.Type == kwSQL
		if p.cur.Type == kwSQL {
			if sqlSecurity, err = p.parseViewSQLSecurity(); err != nil {
				return nil, err
			}
		}
		// Now dispatch on the object keyword
		switch p.cur.Type {
		case kwVIEW:
			stmt, err := p.parseCreateViewStmt(start, orReplace)
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
			stmt, err := p.parseCreateFunctionStmt(start, orReplace, false)
			if err != nil {
				return nil, err
			}
			if stmt.Definer == "" {
				stmt.Definer = definer
			}
			return stmt, nil
		case kwPROCEDURE:
			stmt, err := p.parseCreateFunctionStmt(start, orReplace, true)
			if err != nil {
				return nil, err
			}
			if stmt.Definer == "" {
				stmt.Definer = definer
			}
			return stmt, nil
		case kwTRIGGER:
			stmt, err := p.parseCreateTriggerStmt(start, orReplace)
			if err != nil {
				return nil, err
			}
			if stmt.Definer == "" {
				stmt.Definer = definer
			}
			return stmt, nil
		case kwEVENT:
			stmt, err := p.parseCreateEventStmt(start, orReplace)
			if err != nil {
				return nil, err
			}
			if stmt.Definer == "" {
				stmt.Definer = definer
			}
			return stmt, nil
		default:
			// definer_opt sp_handler_package_spec: no SQL SECURITY before PACKAGE.
			if p.atWord("PACKAGE") && !sawSQLSecurity {
				return p.parseCreatePackageStmt(start, orReplace, definer)
			}
			return nil, &ParseError{Message: "unexpected token after DEFINER clause", Position: p.cur.Loc}
		}

	default:
		if p.atWord("PACKAGE") {
			return p.parseCreatePackageStmt(start, orReplace, "")
		}
		return nil, &ParseError{
			Message:  "unexpected token after CREATE",
			Position: p.cur.Loc,
		}
	}
}

// parseCreateIfNotExists parses an optional IF NOT EXISTS after a CREATE
// object keyword and reports whether it was present. MariaDB rejects the
// clause together with OR REPLACE while parsing: the create rule's actions in
// sql_yacc.yy call LEX::check_create_options (sql/sql_lex.h), which raises
// ER_WRONG_USAGE (1221). The error is reported at IF.
func (p *Parser) parseCreateIfNotExists(orReplace bool) (bool, error) {
	if p.cur.Type != kwIF {
		return false, nil
	}
	ifPos := p.pos()
	p.advance()
	if _, err := p.expect(kwNOT); err != nil {
		return false, err
	}
	if _, err := p.expect(kwEXISTS_KW); err != nil {
		return false, err
	}
	if orReplace {
		return false, &ParseError{Message: "Incorrect usage of OR REPLACE and IF NOT EXISTS", Position: ifPos}
	}
	return true, nil
}

// parseAlterDispatch dispatches ALTER statements to the appropriate parser.
func (p *Parser) parseAlterDispatch() (nodes.Node, error) {
	start := p.pos()
	p.advance() // consume ALTER

	// Completion: after ALTER keyword, offer object type candidates.
	p.checkCursor()
	if p.collectMode() {
		for _, t := range []int{kwTABLE, kwSEQUENCE, kwDATABASE, kwVIEW, kwFUNCTION, kwPROCEDURE, kwEVENT} {
			p.addTokenCandidate(t)
		}
		return nil, &ParseError{Message: "collecting"}
	}

	// Optional IGNORE keyword between ALTER and TABLE: legacy copy-ALTER
	// syntax MariaDB still accepts (MySQL removed it in 5.7). It is valid
	// only immediately before TABLE and carries no AST state.
	// Ref: MariaDB ALTER TABLE grammar, "ALTER [ONLINE] [IGNORE] TABLE ...";
	// measured accept on MariaDB 11.8.8.
	// ALTER [ONLINE] [IGNORE] TABLE: alter_option_list, ONLINE and IGNORE in any
	// order and number, before TABLE only.
	// Ref: alter and alter_option_list in MariaDB 11.8 sql/sql_yacc.yy.
	var online, ignore bool
	for p.cur.Type == kwIGNORE || p.atBareWord("ONLINE") {
		if p.cur.Type == kwIGNORE {
			ignore = true
		} else {
			online = true
		}
		p.advance()
		if p.cur.Type != kwIGNORE && !p.atBareWord("ONLINE") && p.cur.Type != kwTABLE {
			return nil, p.syntaxErrorAtCur()
		}
	}
	if online || ignore {
		stmt, err := p.parseAlterTableStmt(start)
		if err != nil {
			return nil, err
		}
		stmt.Online, stmt.Ignore = online, ignore
		return stmt, nil
	}

	switch p.cur.Type {
	case kwTABLE:
		return p.parseAlterTableStmt(start)

	case kwSEQUENCE:
		return p.parseAlterSequenceStmt(start)

	case kwDATABASE, kwSCHEMA:
		return p.parseAlterDatabaseStmt(start)

	case kwUSER:
		return p.parseAlterUserStmt(start)

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
		for _, t := range []int{kwTABLE, kwSEQUENCE, kwINDEX, kwVIEW, kwDATABASE, kwFUNCTION, kwPROCEDURE, kwTRIGGER, kwEVENT, kwIF} {
			p.addTokenCandidate(t)
		}
		return nil, &ParseError{Message: "collecting"}
	}

	temporary := false
	if p.cur.Type == kwTEMPORARY {
		p.advance()
		temporary = true
	}

	switch p.cur.Type {
	case kwTABLE:
		return p.parseDropTableStmt(start, temporary)

	case kwSEQUENCE:
		return p.parseDropSequenceStmt(start, temporary)

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

	case kwSERVER:
		return p.parseDropServerStmt(start)

	case kwSPATIAL:
		p.advance() // consume SPATIAL
		return p.parseDropSpatialRefSysStmt(start)

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
		if p.atWord("PACKAGE") && !temporary {
			return p.parseDropPackageStmt(start)
		}
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
// definer: DEFINER_SYM '=' user_or_role (MariaDB 11.8 sql/sql_yacc.yy), so `=` and the
// user are required; CURRENT_USER and CURRENT_ROLE take optional_braces.
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

	// Parse user part: identifier, string literal, CURRENT_USER or CURRENT_ROLE
	var user string
	isCurUser := p.cur.Type == kwCURRENT_USER
	// CURRENT_ROLE is no keyword here; the source text keeps `CURRENT_ROLE`, a name, apart.
	isCurRole := p.cur.Type == tokIDENT && strings.EqualFold(p.inputText(p.cur.Loc, p.cur.End), "CURRENT_ROLE")
	if p.cur.Type == tokSCONST {
		user = "'" + p.cur.Str + "'"
		p.advance()
	} else if isCurUser {
		user = "CURRENT_USER"
		p.advance()
	} else if isCurRole {
		user = "CURRENT_ROLE"
		p.advance()
	} else if p.isIdentToken() {
		var err error
		if user, _, err = p.parseIdentifier(); err != nil {
			return "", err
		}
	} else {
		return "", p.syntaxErrorAtCur()
	}

	// Handle CURRENT_USER() and CURRENT_ROLE()
	if (isCurUser || isCurRole) && p.cur.Type == '(' {
		p.advance() // consume (
		if _, err := p.expect(')'); err != nil {
			return "", err
		}
		return user + "()", nil
	}
	if isCurRole {
		return user, nil
	}

	// Handle @host: the lexer scans @host as a single tokIDENT with Str="@host",
	// or for quoted form 'user'@'host' the @ is also part of a variable token.
	// Check for tokIDENT starting with "@".
	if p.cur.Type == tokIDENT && len(p.cur.Str) > 0 && p.cur.Str[0] == '@' {
		atHost := p.cur.Str // e.g. "@localhost" or "@"
		p.advance()
		if atHost == "@" {
			// Standalone @ followed by a host string or name that touches it. With
			// space after it the host is empty: mariadb:11.8.8 reads `root@ VIEW` as
			// root@'' and refuses `root@ 'h'` and `root@ h`.
			var host string
			adjacent := p.cur.Loc == p.prev.End
			if adjacent && p.cur.Type == tokSCONST {
				host = "'" + p.cur.Str + "'"
				p.advance()
			} else if adjacent && p.isIdentToken() {
				var err error
				if host, _, err = p.parseIdentifier(); err != nil {
					return "", err
				}
			}
			return user + "@" + host, nil
		}
		// @host as single token: strip @ prefix
		host := atHost[1:]
		return user + "@" + host, nil
	}

	return user, nil
}
