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
		if p.isIdentToken() {
			name, _, err := p.parseIdentifier()
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

	case kwAGGREGATE:
		// CREATE AGGREGATE FUNCTION — loadable UDF
		p.advance() // consume AGGREGATE
		if p.cur.Type == kwFUNCTION {
			return p.parseCreateLoadableFunction(true)
		}
		return nil, &ParseError{Message: "expected FUNCTION after AGGREGATE", Position: p.cur.Loc}

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
		return p.parseCreateLogfileGroupStmt(start)

	case kwALGORITHM:
		// Only CREATE VIEW uses ALGORITHM
		return p.parseCreateViewStmt(orReplace)

	case kwSQL:
		// CREATE SQL SECURITY {DEFINER|INVOKER} VIEW — delegate to CREATE VIEW
		return p.parseCreateViewStmt(orReplace)

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
			stmt, err := p.parseCreateViewStmt(orReplace)
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
			stmt, err := p.parseCreateFunctionStmt(false)
			if err != nil {
				return nil, err
			}
			if stmt.Definer == "" {
				stmt.Definer = definer
			}
			return stmt, nil
		case kwPROCEDURE:
			stmt, err := p.parseCreateFunctionStmt(true)
			if err != nil {
				return nil, err
			}
			if stmt.Definer == "" {
				stmt.Definer = definer
			}
			return stmt, nil
		case kwTRIGGER:
			stmt, err := p.parseCreateTriggerStmt()
			if err != nil {
				return nil, err
			}
			if stmt.Definer == "" {
				stmt.Definer = definer
			}
			return stmt, nil
		case kwEVENT:
			stmt, err := p.parseCreateEventStmt()
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

	case kwSQL:
		// ALTER SQL SECURITY {DEFINER|INVOKER} VIEW — delegate to ALTER VIEW
		return p.parseAlterViewStmt()

	case kwDEFINER:
		// ALTER DEFINER = user {VIEW|EVENT}
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
		switch p.cur.Type {
		case kwVIEW:
			stmt, err := p.parseAlterViewStmt()
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
			stmt, err := p.parseAlterEventStmt()
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
		return p.parseAlterEventStmt()

	case kwFUNCTION:
		return p.parseAlterRoutineStmt(false)

	case kwPROCEDURE:
		return p.parseAlterRoutineStmt(true)

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
			Loc:  nodes.Loc{Start: start, End: p.pos()},
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

	// Parse user part: identifier, string literal, or CURRENT_USER
	var user string
	if p.cur.Type == tokSCONST {
		user = "'" + p.cur.Str + "'"
		p.advance()
	} else if p.isIdentToken() {
		user, _, _ = p.parseIdentifier()
	} else {
		return ""
	}

	// Handle CURRENT_USER()
	if eqFold(user, "current_user") && p.cur.Type == '(' {
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
