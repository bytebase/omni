package parser

import (
	nodes "github.com/bytebase/omni/mysql/ast"
)

// parseStmt parses a single SQL statement.
func (p *Parser) parseStmt() (nodes.Node, error) {
	switch p.cur.Type {
	case kwWITH:
		return p.parseWithStmt()

	case kwSELECT:
		return p.parseSelectStmt()

	case kwTABLE:
		return p.parseTableStmt()

	case kwVALUES:
		return p.parseValuesStmt()

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
		return p.parseRenameTableStmt()

	case kwSET:
		return p.parseSetStmt()

	case kwSHOW:
		return p.parseShowStmt()

	case kwUSE:
		return p.parseUseStmt()

	case kwDESCRIBE, kwEXPLAIN:
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
		return p.parseLockTablesStmt()

	case kwUNLOCK:
		start := p.pos()
		p.advance() // consume UNLOCK
		if p.cur.Type == kwINSTANCE {
			return p.parseUnlockInstanceStmt(start)
		}
		return p.parseUnlockTablesStmt()

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
		return nil, &ParseError{
			Message:  "unexpected token at statement start",
			Position: p.cur.Loc,
		}
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

	if p.isIdentToken() && eqFold(p.cur.Str, "GROUP_REPLICATION") {
		p.advance() // consume GROUP_REPLICATION
		return p.parseStartGroupReplicationStmt(start)
	}

	// Fall through to START TRANSACTION (parseBeginStmt expects p.cur is START or BEGIN)
	// We already consumed START, so we need to handle TRANSACTION here
	p.match(kwTRANSACTION)
	stmt := &nodes.BeginStmt{Loc: nodes.Loc{Start: start}}
	for {
		if p.cur.Type == kwWITH {
			p.advance()
			p.match(kwCONSISTENT)
			p.match(kwSNAPSHOT)
			stmt.WithConsistentSnapshot = true
		} else if p.cur.Type == kwREAD {
			p.advance()
			if _, ok := p.match(kwONLY); ok {
				stmt.ReadOnly = true
			} else if _, ok := p.match(kwWRITE); ok {
				stmt.ReadWrite = true
			}
		} else {
			break
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	stmt.Loc.End = p.pos()
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

	// Fall through to generic RESET (re-assemble the FlushStmt from utility.go)
	stmt := &nodes.FlushStmt{Loc: nodes.Loc{Start: start}}
	for {
		if p.cur.Type == tokEOF || p.cur.Type == ';' {
			break
		}
		if p.isIdentToken() {
			name, _, _ := p.parseIdentifier()
			stmt.Options = append(stmt.Options, name)
		} else {
			break
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	stmt.Loc.End = p.pos()
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

	if p.isIdentToken() && eqFold(p.cur.Str, "GROUP_REPLICATION") {
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

	orReplace := false
	temporary := false

	// OR REPLACE
	if p.cur.Type == kwOR {
		p.advance()
		p.match(kwREPLACE)
		orReplace = true
	}

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
		return p.parseCreateTableStmt(temporary)

	case kwINDEX:
		return p.parseCreateIndexStmt(false, false, false)

	case kwUNIQUE:
		p.advance()
		if p.cur.Type == kwINDEX {
			return p.parseCreateIndexStmt(true, false, false)
		}
		return nil, &ParseError{Message: "expected INDEX after UNIQUE", Position: p.cur.Loc}

	case kwFULLTEXT:
		p.advance()
		// p.cur is now INDEX or KEY; parseCreateIndexStmt consumes it
		return p.parseCreateIndexStmt(false, true, false)

	case kwSPATIAL:
		p.advance()
		// Distinguish SPATIAL REFERENCE SYSTEM from SPATIAL INDEX
		if p.isIdentToken() && eqFold(p.cur.Str, "reference") {
			return p.parseCreateSpatialRefSysStmt(start, orReplace)
		}
		return p.parseCreateIndexStmt(false, false, true)

	case kwVIEW:
		return p.parseCreateViewStmt(orReplace)

	case kwDATABASE, kwSCHEMA:
		return p.parseCreateDatabaseStmt()

	case kwFUNCTION:
		return p.parseCreateFunctionStmt(false)

	case kwPROCEDURE:
		return p.parseCreateFunctionStmt(true)

	case kwTRIGGER:
		return p.parseCreateTriggerStmt()

	case kwEVENT:
		return p.parseCreateEventStmt()

	case kwUSER:
		return p.parseCreateUserStmt()

	case kwROLE:
		return p.parseCreateRoleStmt()

	case kwTABLESPACE:
		return p.parseCreateTablespaceStmt(start, false)

	case kwSERVER:
		return p.parseCreateServerStmt(start)

	case kwLOGFILE:
		// CREATE LOGFILE GROUP — NDB-specific, use RawStmt
		for p.cur.Type != tokEOF && p.cur.Type != ';' {
			p.advance()
		}
		return &nodes.RawStmt{}, nil

	case kwALGORITHM:
		// Only CREATE VIEW uses ALGORITHM
		return p.parseCreateViewStmt(orReplace)

	case kwDEFINER:
		// CREATE DEFINER=... {VIEW|FUNCTION|PROCEDURE|TRIGGER|EVENT}
		// Skip DEFINER = user[@host] to find the object keyword, then delegate.
		// Sub-parsers won't see DEFINER since we consume it here.
		p.advance() // consume DEFINER
		p.match('=')
		// user[@host] or CURRENT_USER[()]
		if p.isIdentToken() || p.cur.Type == tokSCONST {
			p.advance()
		}
		if p.cur.Type == tokIDENT && p.cur.Str == "@" {
			p.advance()
			if p.isIdentToken() || p.cur.Type == tokSCONST {
				p.advance()
			}
		}
		// Optional SQL SECURITY {DEFINER|INVOKER}
		if p.cur.Type == kwSQL {
			p.advance()
			p.match(kwSECURITY)
			if p.isIdentToken() {
				p.advance()
			}
		}
		// Now dispatch on the object keyword
		switch p.cur.Type {
		case kwVIEW:
			return p.parseCreateViewStmt(orReplace)
		case kwFUNCTION:
			return p.parseCreateFunctionStmt(false)
		case kwPROCEDURE:
			return p.parseCreateFunctionStmt(true)
		case kwTRIGGER:
			return p.parseCreateTriggerStmt()
		case kwEVENT:
			return p.parseCreateEventStmt()
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

// parseAlterDispatch dispatches ALTER statements to the appropriate parser.
func (p *Parser) parseAlterDispatch() (nodes.Node, error) {
	start := p.pos()
	p.advance() // consume ALTER

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
		return p.parseAlterTableStmt()

	case kwDATABASE, kwSCHEMA:
		return p.parseAlterDatabaseStmt()

	case kwUSER:
		return p.parseAlterUserStmt()

	case kwTABLESPACE:
		return p.parseAlterTablespaceStmt(start, false)

	case kwSERVER:
		return p.parseAlterServerStmt(start)

	case kwVIEW:
		return p.parseAlterViewStmt()

	case kwALGORITHM:
		// ALTER ALGORITHM = ... VIEW — delegate to ALTER VIEW parser
		return p.parseAlterViewStmt()

	case kwDEFINER:
		// ALTER DEFINER = user {VIEW|EVENT|FUNCTION|PROCEDURE}
		// Skip DEFINER = user[@host] to find the object keyword, then delegate.
		p.advance() // consume DEFINER
		p.match('=')
		if p.isIdentToken() || p.cur.Type == tokSCONST {
			p.advance()
		}
		if p.cur.Type == tokIDENT && p.cur.Str == "@" {
			p.advance()
			if p.isIdentToken() || p.cur.Type == tokSCONST {
				p.advance()
			}
		}
		// Optional SQL SECURITY {DEFINER|INVOKER}
		if p.cur.Type == kwSQL {
			p.advance()
			p.match(kwSECURITY)
			if p.isIdentToken() {
				p.advance()
			}
		}
		switch p.cur.Type {
		case kwVIEW:
			return p.parseAlterViewStmt()
		case kwEVENT:
			return p.parseAlterEventStmt()
		default:
			return nil, &ParseError{Message: "unexpected token after DEFINER clause in ALTER", Position: p.cur.Loc}
		}

	case kwEVENT:
		return p.parseAlterEventStmt()

	case kwFUNCTION:
		return p.parseAlterRoutineStmt(false)

	case kwPROCEDURE:
		return p.parseAlterRoutineStmt(true)

	case kwLOGFILE:
		// ALTER LOGFILE GROUP — NDB-specific, use RawStmt
		for p.cur.Type != tokEOF && p.cur.Type != ';' {
			p.advance()
		}
		return &nodes.RawStmt{}, nil

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
		return p.parseDropTableStmt(temporary)

	case kwINDEX:
		return p.parseDropIndexStmt()

	case kwVIEW:
		return p.parseDropViewStmt()

	case kwDATABASE, kwSCHEMA:
		return p.parseDropDatabaseStmt()

	case kwUSER:
		return p.parseDropUserStmt()

	case kwROLE:
		return p.parseDropRoleStmt()

	case kwFUNCTION:
		return p.parseDropRoutineStmt(false)

	case kwPROCEDURE:
		return p.parseDropRoutineStmt(true)

	case kwTRIGGER:
		return p.parseDropTriggerStmt()

	case kwEVENT:
		return p.parseDropEventStmt()

	case kwTABLESPACE:
		return p.parseDropTablespaceStmt(start, false)

	case kwSERVER:
		return p.parseDropServerStmt(start)

	case kwLOGFILE:
		// DROP LOGFILE GROUP — NDB-specific, use RawStmt
		for p.cur.Type != tokEOF && p.cur.Type != ';' {
			p.advance()
		}
		return &nodes.RawStmt{}, nil

	case kwSPATIAL:
		p.advance() // consume SPATIAL
		return p.parseDropSpatialRefSysStmt(start)

	case kwRESOURCE:
		p.advance() // consume RESOURCE
		return p.parseDropResourceGroupStmt(start)

	default:
		return nil, &ParseError{
			Message:  "unexpected token after DROP",
			Position: p.cur.Loc,
		}
	}
}
