package parser

import (
	nodes "github.com/bytebase/omni/oceanbase/ast"
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
		return p.parseQueryStatement()

	default:
		// AUDIT and NOAUDIT are no keywords of this lexer's.
		if p.isBareWord("AUDIT") || p.isBareWord("NOAUDIT") {
			return p.parseAuditStmt()
		}
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

	return p.parseStartTransactionTail(start)
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

	// Fall through to generic RESET (re-assemble the FlushStmt from utility.go)
	// Options may be multi-word (e.g. RESET QUERY CACHE) or comma-separated.
	// Require at least one reset option.
	if p.cur.Type == tokEOF || p.cur.Type == ';' {
		return nil, p.syntaxErrorAtCur()
	}
	stmt := &nodes.FlushStmt{Loc: nodes.Loc{Start: start}}
	for {
		if p.cur.Type == tokEOF || p.cur.Type == ';' {
			break
		}
		// Reset option words must accept any keyword (e.g., QUERY CACHE).
		if p.cur.Type == tokIDENT || p.cur.Type >= 700 {
			name, _, err := p.parseKeywordOrIdent()
			if err != nil {
				return nil, err
			}
			stmt.Options = append(stmt.Options, name)
		} else {
			break
		}
		// Consume optional comma separator between options
		if p.cur.Type == ',' {
			p.advance()
		}
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

	// Fall through to LOAD DATA / LOAD XML
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

	if stmt, ok, err := p.parseClusterDDL(start, "CREATE"); ok {
		return stmt, err
	}

	temporary := false

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

	case kwVECTOR:
		// CREATE VECTOR INDEX, spelled INDEX only: 4.3.5.6 answers 1064 for CREATE VECTOR
		// KEY and CREATE VECTOR, and for OR REPLACE and TEMPORARY in front of it.
		if temporary {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		if p.cur.Type != kwINDEX {
			return nil, p.syntaxErrorAtCur()
		}
		stmt, err := p.parseCreateIndexStmt(start, false, false, false)
		if err != nil {
			return nil, err
		}
		stmt.Vector = true
		return stmt, nil

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
		return p.parseCreateTablespaceStmt(start)

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
		// CREATE DEFINER=... {FUNCTION|PROCEDURE|TRIGGER|EVENT}. A view's DEFINER
		// follows ALGORITHM on OceanBase, so here DEFINER heads a stored program,
		// which takes no SQL SECURITY before its keyword: 4.3.5.6 answers 1064 at
		// VIEW for CREATE DEFINER = u VIEW v AS SELECT 1, and at SQL for CREATE
		// DEFINER = u SQL SECURITY INVOKER VIEW (or PROCEDURE). See
		// parseCreateViewStmt.
		definer, err := p.parseDefinerValue(true)
		if err != nil {
			return nil, err
		}
		switch p.cur.Type {
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

	default:
		return nil, &ParseError{
			Message:  "unexpected token after CREATE",
			Position: p.cur.Loc,
		}
	}
}

// parseCreateOrReplace parses CREATE OR REPLACE. CREATE is consumed and p.cur
// is OR. OceanBase takes OR REPLACE only before a view; on any other object,
// and before TEMPORARY, the statement fails at the keyword that follows OR
// REPLACE.
//
// Ref: create_view_stmt in OceanBase's src/sql/parser/sql_parser_mysql_mode.y
// (v4.3.5_CE_BP6, line 9717): create_with_opt_hint opt_replace opt_algorithm
// opt_definer opt_sql_security VIEW. The engine answers 1064 at DEFINER unless
// ALGORITHM comes first (see parseCreateViewStmt), and has no CREATE OR REPLACE
// SPATIAL REFERENCE SYSTEM, which MySQL 8.0 has.
func (p *Parser) parseCreateOrReplace(start int) (nodes.Node, error) {
	p.advance() // consume OR
	if p.cur.Type != kwREPLACE {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume REPLACE

	if stmt, ok, err := p.parseClusterDDL(start, "CREATE OR REPLACE"); ok {
		return stmt, err
	}

	switch p.cur.Type {
	case kwVIEW, kwALGORITHM, kwSQL:
		// parseCreateViewStmt reads ALGORITHM, DEFINER, and SQL SECURITY in
		// that order and fails at whatever stands where VIEW belongs.
		return p.parseCreateViewStmt(start, true)
	default:
		// DEFINER, which OceanBase takes on a view only after ALGORITHM, and
		// SPATIAL, since OceanBase has no SPATIAL REFERENCE SYSTEM, fail here like
		// every other object: 4.3.5.6 answers 1064 at the word after OR REPLACE.
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

	if stmt, ok, err := p.parseClusterDDL(start, "ALTER"); ok {
		return stmt, err
	}

	// ALTER DISKGROUP is one of alter_system_stmt's alternatives.
	if p.isBareWord("DISKGROUP") {
		return p.parseAlterDiskgroupStmt(start)
	}

	switch p.cur.Type {
	case kwTABLE:
		return p.parseAlterTableStmt(start)

	case kwSYSTEM:
		return p.parseAlterSystemStmt(start)

	case kwDATABASE, kwSCHEMA:
		return p.parseAlterDatabaseStmt(start)

	case kwUSER:
		return p.parseAlterUserStmt(start)

	case kwTABLESPACE:
		return p.parseAlterTablespaceStmt(start)

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
		// ALTER DEFINER = user EVENT. A view's DEFINER follows ALGORITHM on
		// OceanBase: 4.3.5.6 answers 1064 at VIEW for ALTER DEFINER = u VIEW v AS
		// SELECT 1, and at SQL when SQL SECURITY follows the definer.
		definer, err := p.parseDefinerValue(true)
		if err != nil {
			return nil, err
		}
		switch p.cur.Type {
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

	if stmt, ok, err := p.parseClusterDDL(start, "DROP"); ok {
		return stmt, err
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
		return p.parseDropTablespaceStmt(start)

	case kwSERVER:
		return p.parseDropServerStmt(start)

	case kwLOGFILE:
		return p.parseDropLogfileGroupStmt(start)

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
		return nil, &ParseError{
			Message:  "unexpected token after DROP",
			Position: p.cur.Loc,
		}
	}
}

// parseDefinerValue parses DEFINER = user, where the = and the user are both
// required, as is the ) of CURRENT_USER(): 5.7.25-OceanBase_CE-v4.3.5.6 answers
// 1064 for DEFINER root, DEFINER = followed by the object keyword, and
// DEFINER = CURRENT_USER( (measured 2026-11-06). It takes user@ with no host,
// as the server does.
//
// program says the DEFINER heads a stored program and not a view. OceanBase reads a
// function, procedure, trigger or event with another lexer, whose user_variable has an
// alternative the SQL lexer's lacks: `(@%)` (pl_parser_mysql_mode.l, v4.3.5_CE_BP6, line
// 55, against sql_parser_mysql_mode.l, line 68). So 4.3.5.6 takes DEFINER = root@%
// before a stored program, which is how SHOW CREATE TRIGGER prints the definer, and
// answers 1064 at the % before a view. The % is the whole host and follows its @ at
// once: root@ %, root@%example and root@%.example.com are 1064s.
// TestDefinerHostFormsOceanBaseOracle holds each to the engine.
func (p *Parser) parseDefinerValue(program bool) (string, error) {
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
		var err error
		if user, _, err = p.parseIdentifier(); err != nil {
			return "", err
		}
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
				var err error
				if host, _, err = p.parseIdentifier(); err != nil {
					return "", err
				}
			} else if program && !isCurUser && p.cur.Type == '%' && p.cur.Loc == p.prev.End {
				host = "%"
				p.advance()
			}
			return user + "@" + host, nil
		}
		// @host as single token: strip @ prefix
		host := atHost[1:]
		return user + "@" + host, nil
	}

	return user, nil
}
