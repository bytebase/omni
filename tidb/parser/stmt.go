package parser

import (
	nodes "github.com/bytebase/omni/tidb/ast"
)

// parseStmt parses a single SQL statement.
func (p *Parser) parseStmt() (nodes.Node, error) {
	// Completion: check if cursor has been reached.
	p.checkCursor()

	if p.collectMode() {
		// Add all top-level statement-starting keywords as candidates.
		stmtTokens := []int{
			kwSELECT, kwINSERT, kwUPDATE, kwDELETE, kwBATCH,
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
			kwADMIN, // TiDB's own maintenance family
			kwSPLIT, kwDISTRIBUTE,
			kwCALIBRATE, kwQUERY,
			kwBACKUP, kwRESTORE, kwRECOVER, kwFLASHBACK,
			kwPLAN, // PLAN REPLAYER
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
	// query expression (container-verified: "table t1 union table t2",
	// "values row(1) union select 1").
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

	case kwBATCH:
		return p.parseBatchStmt()

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
		if p.isShowBRIE() {
			return p.parseShowBRIEStmt()
		}
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
		switch {
		case p.cur.Type == kwINSTANCE:
			return p.parseLockInstanceStmt(start)
		case p.isBareWord("STATS"):
			return p.parseLockStatsStmt(start)
		}
		return p.parseLockTablesStmt(start)

	case kwUNLOCK:
		start := p.pos()
		p.advance() // consume UNLOCK
		switch {
		case p.cur.Type == kwINSTANCE:
			return p.parseUnlockInstanceStmt(start)
		case p.isBareWord("STATS"):
			return p.parseUnlockStatsStmt(start)
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

	case kwSIGNAL, kwRESIGNAL:
		// TiDB v8.5.0 has no SIGNAL / RESIGNAL statement.
		return nil, p.syntaxErrorAtCur()

	case kwGET:
		return p.parseGetDiagnosticsStmt()

	case kwCHANGE:
		return p.parseChangeDispatch()

	case kwSTOP:
		return p.parseStopDispatch()

	case kwPURGE:
		return p.parsePurgeStmt()

	case kwRESUME:
		return p.parseResumeBackupLogsStmt()

	case kwIMPORT:
		return p.parseImportDispatch()

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

	case kwSPLIT:
		return p.parseSplitRegionStmt()

	case kwDISTRIBUTE:
		return p.parseDistributeTableStmt()

	case kwUNINSTALL:
		return p.parseUninstallStmt()

	case kwADMIN:
		return p.parseAdminStmt()

	case kwBACKUP:
		return p.parseBackupStmt()

	case kwRESTORE:
		return p.parseRestoreStmt()

	case kwRECOVER:
		return p.parseRecoverTableStmt()

	case kwFLASHBACK:
		return p.parseFlashbackStmt()

	case kwCALIBRATE:
		return p.parseCalibrateResourceStmt()

	case kwQUERY:
		// QUERY WATCH ADD / REMOVE, TiDB's runaway-query watch list. QUERY
		// starts no other statement, but it is an unreserved keyword and so
		// an ordinary identifier everywhere else; only WATCH after it makes
		// a statement.
		if p.peekNext().Type == kwWATCH {
			return p.parseQueryWatchStmt()
		}
		return nil, p.syntaxErrorAtCur()

	case kwPLAN:
		// PLAN REPLAYER, TiDB's plan-reproduction dump. PLAN starts no other
		// statement; it is a keyword here only so the SQL plan bindings can
		// read it, and stays an ordinary identifier everywhere else, so
		// REPLAYER after it is what makes a statement.
		return p.parsePlanReplayerStmt()

	case '(':
		// Parenthesized subquery / select
		return p.parseSelectStmt()

	default:
		// TRACE, TiDB's execution and optimizer tracer. It is not a keyword
		// in this lexer and needs none: it is read as a bare word here, so
		// `trace` stays usable as a table name, a column name and an
		// implicit alias, which registering it would have put at risk.
		if p.isBareWord("TRACE") {
			return p.parseTraceStmt()
		}
		// PAUSE, CANCEL and RECOMMEND are keywords upstream and start one
		// family of statements each. They are bare words here for the same
		// reason TRACE is.
		switch {
		case p.isBareWord("PAUSE"):
			return p.parsePauseBackupLogsStmt()
		case p.isBareWord("CANCEL"):
			return p.parseCancelStmt()
		case p.isBareWord("RECOMMEND"):
			return p.parseRecommendIndexStmt()
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

	return p.parseStartTransactionStmt(start)
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

	if p.cur.Type == kwBACKUP {
		return p.parseStopBackupLogsStmt(start)
	}

	return nil, &ParseError{Message: "expected REPLICA, SLAVE, GROUP_REPLICATION, or BACKUP after STOP", Position: p.cur.Loc}
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

	// LOAD STATS, which TiDB has and MySQL has not. STATS is read as a bare
	// word for the reason statistics.go gives: v8.5.5 answers 1064 for
	// `LOAD `STATS` 'x'`, where a backquoted name is no part of this
	// grammar at all.
	if p.isBareWord("STATS") {
		return p.parseLoadStatsStmt(start)
	}

	// Fall through to LOAD DATA / LOAD XML
	return p.parseLoadDataStmt(start)
}

// parseImportDispatch dispatches IMPORT statements. TiDB has one, IMPORT
// INTO; MySQL's `IMPORT TABLE FROM 'x.sdi'` is not TiDB's and v8.5.5 answers
// 1064 for it, where mysql:8.0.46 takes it (mariadb:11.8.8 and OceanBase
// 4.3.5 answer 1064 too, read on 2026-11-04). This parser used to read it,
// which also left `IMPORT INTO t` parsing as an IMPORT TABLE of no files.
func (p *Parser) parseImportDispatch() (nodes.Node, error) {
	start := p.pos()
	p.advance() // consume IMPORT

	if p.cur.Type != kwINTO {
		return nil, p.syntaxErrorAtCur()
	}
	return p.parseImportIntoStmt(start)
}

// parseCreateDispatch dispatches CREATE statements to the appropriate parser.
func (p *Parser) parseCreateDispatch() (nodes.Node, error) {
	start := p.pos()
	p.advance() // consume CREATE

	// Completion: after CREATE keyword, offer object type candidates.
	p.checkCursor()
	if p.collectMode() {
		for _, t := range []int{
			kwTABLE, kwINDEX, kwVIEW, kwDATABASE, kwPROCEDURE, // TiDB has no CREATE FUNCTION/TRIGGER/EVENT
			kwPLACEMENT, // TiDB: CREATE PLACEMENT POLICY
			kwSEQUENCE,
			kwBINDING, // TiDB: CREATE [GLOBAL|SESSION] BINDING
		} {
			p.addTokenCandidate(t)
		}
		return nil, &ParseError{Message: "collecting"}
	}

	if p.cur.Type == kwOR {
		return p.parseCreateOrReplace(start)
	}

	temporary := false

	// CREATE [GLOBAL | SESSION] BINDING, TiDB's SQL plan binding. The scope
	// word is read only when BINDING follows it, so `CREATE GLOBAL TEMPORARY
	// TABLE` below and `CREATE SESSION` on its own are untouched.
	if p.cur.Type == kwBINDING {
		return p.parseCreateBindingStmt(start, nodes.BindingScopeUnset)
	}
	if (p.cur.Type == kwGLOBAL || p.cur.Type == kwSESSION) && p.peekNext().Type == kwBINDING {
		scope := p.parseBindingScope()
		return p.parseCreateBindingStmt(start, scope)
	}

	// CREATE STATISTICS, TiDB's extended statistics. STATISTICS is a bare
	// word, so `CREATE TABLE statistics (a int)` is untouched.
	if p.isBareWord("STATISTICS") {
		return p.parseCreateStatisticsStmt(start)
	}

	// GLOBAL TEMPORARY, TiDB's own kind of temporary table (parser.y, rule
	// CreateTableStmt with OptTemporary "GLOBAL TEMPORARY"). GLOBAL on its own
	// opens nothing, so `CREATE GLOBAL TABLE t (a int)` falls through to the
	// switch below and is a 1064, as it is on the server.
	// Temporary stays false for it: only the rows of a global temporary table
	// are per-session, and the table itself is an ordinary schema object (see
	// ast.CreateTableStmt.Global).
	//
	// TABLE is the only target the prefix reaches, and it is required here
	// rather than left to the switch below: v8.5.5 answers 1064 for `CREATE
	// GLOBAL TEMPORARY VIEW`, `... INDEX`, `... DATABASE` and `... SEQUENCE`,
	// and falling through would have read each of them as an ordinary statement
	// of whatever followed, with the prefix dropped on the floor.
	global := false
	if p.cur.Type == kwGLOBAL && p.peekNext().Type == kwTEMPORARY {
		p.advance() // consume GLOBAL
		p.advance() // consume TEMPORARY
		global = true
		if p.cur.Type != kwTABLE && p.cur.Type != kwTABLES {
			return nil, p.syntaxErrorAtCur()
		}
	} else if p.cur.Type == kwTEMPORARY {
		// TEMPORARY
		p.advance()
		temporary = true
	}

	// CreateIndexStmt (parser.y) is `"CREATE" IndexKeyTypeOpt "INDEX" ...`, with
	// no TEMPORARY in front and INDEX the only keyword after UNIQUE, FULLTEXT,
	// SPATIAL or VECTOR: v8.5.5 answers 1064 for `CREATE FULLTEXT KEY` and for
	// `CREATE TEMPORARY INDEX`.
	if temporary {
		switch p.cur.Type {
		case kwINDEX, kwUNIQUE, kwFULLTEXT, kwSPATIAL:
			return nil, p.syntaxErrorAtCur()
		}
	}

	switch p.cur.Type {
	case kwTABLE:
		return p.parseCreateTableStmt(start, temporary, global)

	case kwINDEX:
		return p.parseCreateIndexStmt(start, false, false, false, false)

	case kwUNIQUE:
		p.advance()
		if p.cur.Type == kwINDEX {
			return p.parseCreateIndexStmt(start, true, false, false, false)
		}
		return nil, &ParseError{Message: "expected INDEX after UNIQUE", Position: p.cur.Loc}

	case kwFULLTEXT:
		p.advance()
		if p.cur.Type != kwINDEX {
			return nil, p.syntaxErrorAtCur()
		}
		return p.parseCreateIndexStmt(start, false, true, false, false)

	case kwVECTOR:
		// CREATE VECTOR INDEX, which IndexKeyTypeOpt reads beside UNIQUE, SPATIAL
		// and FULLTEXT and which is spelled INDEX only: v8.5.5 answers 1064 for
		// CREATE VECTOR KEY, for CREATE VECTOR with no keyword after it, and for TEMPORARY
		// in front of it.
		if temporary || global {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		if p.cur.Type != kwINDEX {
			return nil, p.syntaxErrorAtCur()
		}
		stmt, err := p.parseCreateIndexStmt(start, false, false, false, true)
		if err != nil {
			return nil, err
		}
		return stmt, nil

	case kwSPATIAL:
		// CREATE SPATIAL INDEX is the one spatial statement TiDB reads; it has
		// no CREATE SPATIAL REFERENCE SYSTEM, and v8.5.5 answers 1064 there.
		p.advance()
		if p.cur.Type != kwINDEX {
			return nil, p.syntaxErrorAtCur()
		}
		return p.parseCreateIndexStmt(start, false, false, true, false)

	case kwVIEW:
		return p.parseCreateViewStmt(start, false)

	case kwDATABASE, kwSCHEMA:
		return p.parseCreateDatabaseStmt(start)

	case kwPLACEMENT:
		// CREATE PLACEMENT POLICY [IF NOT EXISTS] name opts
		return p.parseCreatePlacementPolicyStmt(start, false)

	case kwSEQUENCE:
		// TiDB v8.5.5 has no CREATE TEMPORARY SEQUENCE (parser.y:15462); it
		// rejects at SEQUENCE. OR REPLACE is refused in parseCreateOrReplace.
		if temporary {
			return nil, p.syntaxErrorAtCur()
		}
		return p.parseCreateSequenceStmt(start)

	case kwAGGREGATE:
		// TiDB v8.5.0 has no CREATE [AGGREGATE] FUNCTION (stored or loadable
		// UDF); it rejects at the FUNCTION/AGGREGATE keyword.
		return nil, p.syntaxErrorAtCur()

	case kwFUNCTION:
		// TiDB v8.5.0 has no CREATE FUNCTION (stored or loadable UDF).
		return nil, p.syntaxErrorAtCur()

	case kwPROCEDURE:
		return p.parseCreateFunctionStmt(start, true)

	case kwTRIGGER:
		// TiDB v8.5.0 has no CREATE TRIGGER.
		return nil, p.syntaxErrorAtCur()

	case kwEVENT:
		// TiDB v8.5.0 has no CREATE EVENT (no event scheduler).
		return nil, p.syntaxErrorAtCur()

	case kwUSER:
		return p.parseCreateUserStmt(start)

	case kwROLE:
		return p.parseCreateRoleStmt(start)

	case kwALGORITHM:
		// Only CREATE VIEW uses ALGORITHM
		return p.parseCreateViewStmt(start, false)

	case kwSQL:
		// CREATE SQL SECURITY {DEFINER|INVOKER} VIEW — delegate to CREATE VIEW
		return p.parseCreateViewStmt(start, false)

	case kwDEFINER:
		// CREATE DEFINER=... {VIEW|FUNCTION|PROCEDURE|TRIGGER|EVENT}
		// Parse DEFINER = user[@host] and propagate to sub-parsers.
		definer := p.parseDefinerValue()
		// Optional SQL SECURITY {DEFINER|INVOKER}
		var sqlSecurity string
		if p.cur.Type == kwSQL {
			p.advance()
			p.match(kwSECURITY)
			if p.isIdentToken() {
				sqlSecurity, _, _ = p.parseIdentifier()
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
		case kwFUNCTION, kwPROCEDURE, kwTRIGGER, kwEVENT:
			// TiDB v8.5.0 accepts a DEFINER clause only on CREATE VIEW. It has
			// no FUNCTION/TRIGGER/EVENT, and CREATE PROCEDURE takes no DEFINER
			// clause — all four reject at the object keyword.
			return nil, p.syntaxErrorAtCur()
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
// is OR. TiDB v8.5.5 takes OR REPLACE only before a view and before
// PLACEMENT POLICY; on any other object, and before TEMPORARY, the statement
// fails at the keyword that follows OR REPLACE. Unlike MySQL 8.0, TiDB has no
// CREATE OR REPLACE SPATIAL REFERENCE SYSTEM.
//
// Ref: TiDB v8.5.5 parser.y:5099 (CreateViewStmt: "CREATE" OrReplace
// ViewAlgorithm ViewDefiner ViewSQLSecurity "VIEW" ...) and parser.y:15428
// (CreatePolicyStmt: "CREATE" OrReplace "PLACEMENT" "POLICY" IfNotExists ...),
// the only two rules that use OrReplace.
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
	case kwPLACEMENT:
		return p.parseCreatePlacementPolicyStmt(start, true)
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
		for _, t := range []int{
			kwTABLE, kwDATABASE, // TiDB has no ALTER VIEW/FUNCTION/PROCEDURE/EVENT
			kwPLACEMENT, // TiDB: ALTER PLACEMENT POLICY
			kwSEQUENCE,
		} {
			p.addTokenCandidate(t)
		}
		return nil, &ParseError{Message: "collecting"}
	}

	switch p.cur.Type {
	case kwTABLE:
		return p.parseAlterTableStmt(start)

	case kwDATABASE, kwSCHEMA:
		return p.parseAlterDatabaseStmt(start)

	case kwPLACEMENT:
		// ALTER PLACEMENT POLICY [IF EXISTS] name opts
		return p.parseAlterPlacementPolicyStmt(start)

	case kwSEQUENCE:
		return p.parseAlterSequenceStmt(start)

	case kwRANGE:
		return p.parseAlterRangeStmt(start)

	case kwUSER:
		return p.parseAlterUserStmt(start)

	case kwVIEW, kwALGORITHM, kwSQL, kwDEFINER:
		// TiDB has no ALTER VIEW in any spelling, and no ALTER EVENT for a
		// DEFINER to lead into: v8.5.5 answers 1064 for `ALTER VIEW`, `ALTER
		// ALGORITHM = ... VIEW`, `ALTER DEFINER = ...` and `ALTER SQL SECURITY
		// ...`. A view is redefined with CREATE OR REPLACE VIEW.
		//
		// Ref: TiDB v8.5.5 pkg/parser/parser.y, where no rule starts with
		// "ALTER" followed by "VIEW", ViewAlgorithm, ViewDefiner or
		// ViewSQLSecurity.
		return nil, p.syntaxErrorAtCur()

	case kwEVENT:
		// TiDB v8.5.0 has no events (no ALTER EVENT).
		return nil, p.syntaxErrorAtCur()

	case kwFUNCTION, kwPROCEDURE:
		// TiDB v8.5.0 has no ALTER FUNCTION / ALTER PROCEDURE.
		return nil, p.syntaxErrorAtCur()

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
		for _, t := range []int{
			kwTABLE, kwINDEX, kwVIEW, kwDATABASE, kwPROCEDURE, kwIF, // TiDB has no DROP FUNCTION/TRIGGER/EVENT
			kwPLACEMENT, // TiDB: DROP PLACEMENT POLICY
			kwSEQUENCE,
			kwBINDING, // TiDB: DROP [GLOBAL|SESSION] BINDING
		} {
			p.addTokenCandidate(t)
		}
		return nil, &ParseError{Message: "collecting"}
	}

	// DROP [GLOBAL | SESSION] BINDING, read before the GLOBAL TEMPORARY
	// prefix below for the same reason as on CREATE.
	if p.cur.Type == kwBINDING {
		return p.parseDropBindingStmt(start, nodes.BindingScopeUnset)
	}
	if (p.cur.Type == kwGLOBAL || p.cur.Type == kwSESSION) && p.peekNext().Type == kwBINDING {
		scope := p.parseBindingScope()
		return p.parseDropBindingStmt(start, scope)
	}

	// DROP STATS and DROP STATISTICS, both TiDB's own. Each word is read
	// bare, so `DROP TABLE stats` and a table named statistics are
	// untouched; v8.5.5 answers 1064 for `DROP `STATS` t`.
	if p.isBareWord("STATS") {
		return p.parseDropStatsStmt(start)
	}
	if p.isBareWord("STATISTICS") {
		return p.parseDropStatisticsStmt(start)
	}

	// DROP HYPO INDEX name ON tbl, TiDB's drop of a hypothetical index. HYPO
	// is read bare, so `DROP TABLE hypo` and a table named hypo are untouched;
	// INDEX must follow, so a bare `DROP hypo` falls through to the switch and
	// is a 1064. Unlike plain DROP INDEX this arm takes no lock or algorithm.
	//
	// Ref: `"DROP" "HYPO" "INDEX" IfExists Identifier "ON" TableName`
	// (DropIndexStmt) in TiDB v8.5.5 pkg/parser/parser.y.
	if p.isBareWord("HYPO") && p.peekNext().Type == kwINDEX {
		return p.parseDropHypoIndexStmt(start)
	}

	// DROP GLOBAL TEMPORARY TABLE, one of the two ways to drop a table CREATE
	// GLOBAL TEMPORARY TABLE made: a plain DROP TABLE drops it too, and DROP
	// TEMPORARY TABLE does not (1051, verified), so Temporary stays false here.
	//
	// As on CREATE, TABLE is the only target the prefix reaches and is required
	// here: v8.5.5 answers 1064 for `DROP GLOBAL TEMPORARY DATABASE`, `... VIEW`
	// and `... INDEX`, each of which the switch below would have read as the
	// ordinary drop of whatever followed. GLOBAL without TEMPORARY after it is
	// not consumed at all, so `DROP GLOBAL TABLE t` reaches a switch with no arm
	// for it and is a 1064 as well.
	temporary := false
	global := false
	if p.cur.Type == kwGLOBAL && p.peekNext().Type == kwTEMPORARY {
		p.advance() // consume GLOBAL
		p.advance() // consume TEMPORARY
		global = true
		if p.cur.Type != kwTABLE && p.cur.Type != kwTABLES {
			return nil, p.syntaxErrorAtCur()
		}
	} else if p.cur.Type == kwTEMPORARY {
		p.advance()
		temporary = true
	}

	switch p.cur.Type {
	case kwTABLE, kwTABLES:
		// TableOrTables (parser.y:5446): DROP accepts the plural TABLES too.
		return p.parseDropTableStmt(start, temporary, global)

	case kwINDEX:
		return p.parseDropIndexStmt(start)

	case kwVIEW:
		return p.parseDropViewStmt(start)

	case kwDATABASE, kwSCHEMA:
		return p.parseDropDatabaseStmt(start)

	case kwPLACEMENT:
		// DROP PLACEMENT POLICY [IF EXISTS] name
		return p.parseDropPlacementPolicyStmt(start)

	case kwSEQUENCE:
		// TiDB v8.5.5 has no DROP TEMPORARY SEQUENCE (parser.y:15575).
		if temporary {
			return nil, p.syntaxErrorAtCur()
		}
		return p.parseDropSequenceStmt(start)

	case kwUSER:
		return p.parseDropUserStmt(start)

	case kwROLE:
		return p.parseDropRoleStmt(start)

	case kwFUNCTION:
		// TiDB v8.5.0 has no stored functions — DROP FUNCTION rejects.
		return nil, p.syntaxErrorAtCur()

	case kwPROCEDURE:
		return p.parseDropRoutineStmt(start, true)

	case kwTRIGGER:
		// TiDB v8.5.0 has no triggers — DROP TRIGGER rejects.
		return nil, p.syntaxErrorAtCur()

	case kwEVENT:
		// TiDB v8.5.0 has no events — DROP EVENT rejects.
		return nil, p.syntaxErrorAtCur()

	// No DROP SPATIAL REFERENCE SYSTEM: TiDB has no spatial reference systems,
	// and v8.5.5 answers 1064 at SPATIAL.

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
func (p *Parser) parseDefinerValue() string {
	p.advance() // consume DEFINER
	p.match('=')

	// Completion: after DEFINER=, offer CURRENT_USER keyword.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwCURRENT_USER)
		return ""
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
		return ""
	}

	// Handle CURRENT_USER()
	if isCurUser && p.cur.Type == '(' {
		p.advance() // consume (
		p.match(')')
		return "CURRENT_USER()"
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
			return user + "@" + host
		}
		// @host as single token: strip @ prefix
		host := atHost[1:]
		return user + "@" + host
	}

	return user
}
