// Package parser — F4 parser entry point and dispatch framework.
//
// The Parser struct, Parse function, and ParseBestEffort function defined
// in this file provide the public recursive-descent parser API for Doris SQL.
// Actual per-statement parsing is stubbed out — F4 ships only the framework.
// Each concrete statement type (SELECT, INSERT, CREATE, etc.) is added by
// Tier 1+ DAG nodes that replace specific dispatch cases in parseStmt below.
package parser

import (
	"github.com/bytebase/omni/doris/ast"
)

// Parser is a recursive-descent parser for Doris SQL. It operates on a
// single segment of input (one top-level statement produced by Split) at a
// time; callers should use Parse or ParseBestEffort instead of constructing
// Parsers directly.
type Parser struct {
	lexer   *Lexer
	input   string       // the segment text (used for error reporting)
	cur     Token        // current token
	prev    Token        // previous token (for error context)
	nextBuf Token        // buffered lookahead token
	hasNext bool         // whether nextBuf is valid
	errors  []ParseError // collected errors for best-effort mode
}

// advance consumes the current token and moves to the next one.
// Returns the token that was just consumed (the new "previous" token).
func (p *Parser) advance() Token {
	p.prev = p.cur
	if p.hasNext {
		p.cur = p.nextBuf
		p.hasNext = false
	} else {
		p.cur = p.lexer.NextToken()
	}
	return p.prev
}

// peek returns the current token without consuming it.
func (p *Parser) peek() Token {
	return p.cur
}

// peekNext returns the token AFTER the current one without consuming it
// (one-token lookahead beyond cur). Used to disambiguate based on the
// token following the current position.
func (p *Parser) peekNext() Token {
	if !p.hasNext {
		p.nextBuf = p.lexer.NextToken()
		p.hasNext = true
	}
	return p.nextBuf
}

// match tries each given token kind; if cur matches any, it is consumed
// and returned with ok=true. Used for optional tokens.
func (p *Parser) match(types ...int) (Token, bool) {
	for _, t := range types {
		if p.cur.Kind == t {
			return p.advance(), true
		}
	}
	return Token{}, false
}

// expect consumes the current token if it matches the expected kind.
// Otherwise returns a ParseError describing the mismatch.
func (p *Parser) expect(tokenType int) (Token, error) {
	if p.cur.Kind == tokenType {
		return p.advance(), nil
	}
	return Token{}, p.syntaxErrorAtCur()
}

// syntaxErrorAtCur returns a *ParseError describing a syntax error at
// the current token position. If cur is tokEOF, the message says "at
// end of input"; otherwise "at or near X".
func (p *Parser) syntaxErrorAtCur() *ParseError {
	var msg string
	if p.cur.Kind == tokEOF {
		msg = "syntax error at end of input"
	} else {
		text := p.cur.Str
		if text == "" {
			text = TokenName(p.cur.Kind)
		}
		msg = "syntax error at or near " + text
	}
	return &ParseError{
		Loc: p.cur.Loc,
		Msg: msg,
	}
}

// skipToNextStatement advances the parser past the current erroneous
// statement to the next ; at depth 0 within the current segment, or to
// tokEOF. Always consumes at least one token to guarantee progress.
//
// Because Parse uses Split to pre-segment the input, each segment usually
// contains one statement — skipToNextStatement's typical behavior inside
// a segment is to advance to EOF.
func (p *Parser) skipToNextStatement() {
	// Always consume at least one token to avoid infinite loops.
	if p.cur.Kind != tokEOF {
		p.advance()
	}
	depth := 0
	for p.cur.Kind != tokEOF {
		switch p.cur.Kind {
		case int('('), int('['), int('{'):
			depth++
		case int(')'), int(']'), int('}'):
			if depth > 0 {
				depth--
			}
		case int(';'):
			if depth == 0 {
				p.advance()
				return
			}
		}
		p.advance()
	}
}

// unsupported emits a "statement type X not yet supported" ParseError at
// the current token, advances to the next statement boundary, and returns
// (nil, err). Used by every stubbed dispatch case in F4.
//
// Tier 1+ DAG nodes REPLACE individual dispatch cases with real parse
// functions. Those real functions do NOT call unsupported — they consume
// tokens themselves and produce real AST nodes.
func (p *Parser) unsupported(stmtName string) (ast.Node, error) {
	err := &ParseError{
		Loc: p.cur.Loc,
		Msg: stmtName + " statement parsing is not yet supported",
	}
	p.skipToNextStatement()
	return nil, err
}

// unknownStatementError reports a statement that starts with a token the
// dispatch switch doesn't recognize. It is called from the default branch
// of parseStmt.
func (p *Parser) unknownStatementError() *ParseError {
	if p.cur.Kind == tokEOF {
		return &ParseError{
			Loc: p.cur.Loc,
			Msg: "syntax error at end of input",
		}
	}
	tokText := p.cur.Str
	if tokText == "" {
		tokText = TokenName(p.cur.Kind)
	}
	return &ParseError{
		Loc: p.cur.Loc,
		Msg: "unknown or unsupported statement starting with " + tokText,
	}
}

// parseStmt parses one top-level statement by dispatching on the first
// keyword. For stubs-only F4, every dispatch case calls the unsupported
// helper that appends a "not yet supported" ParseError and returns
// (nil, err).
//
// Tier 1+ DAG nodes REPLACE specific cases with real implementations.
// The dispatch switch itself is not expected to change — only the bodies
// of individual case branches.
func (p *Parser) parseStmt() (ast.Node, error) {
	switch p.cur.Kind {
	// DDL
	case kwCREATE:
		createTok := p.advance() // consume CREATE; cur is now the object type keyword
		switch p.cur.Kind {
		case kwINDEX:
			return p.parseCreateIndex(createTok.Loc)
		case kwDATABASE, kwSCHEMA:
			return p.parseCreateDatabase()
		case kwTABLE, kwTEMPORARY:
			return p.parseCreateTable()
		case kwEXTERNAL:
			// Peek past EXTERNAL to distinguish CREATE EXTERNAL CATALOG,
			// CREATE EXTERNAL RESOURCE, from CREATE EXTERNAL TABLE.
			switch p.peekNext().Kind {
			case kwCATALOG:
				return p.parseCreateCatalog()
			case kwRESOURCE:
				p.advance() // consume EXTERNAL
				return p.parseCreateResource(createTok.Loc, true)
			}
			return p.parseCreateTable()
		case kwCATALOG:
			return p.parseCreateCatalog()
		case kwWORKLOAD:
			// CREATE WORKLOAD GROUP ... or CREATE WORKLOAD POLICY ...
			next := p.peekNext()
			if next.Kind == kwGROUP {
				return p.parseCreateWorkloadGroup(createTok.Loc)
			}
			if next.Kind == kwPOLICY {
				return p.parseCreateWorkloadPolicy(createTok.Loc)
			}
			return p.unsupported("CREATE WORKLOAD")
		case kwRESOURCE:
			return p.parseCreateResource(createTok.Loc, false)
		case kwSQL_BLOCK_RULE:
			return p.parseCreateSQLBlockRule(createTok.Loc)
		case kwVIEW:
			return p.parseCreateView(createTok.Loc, false)
		case kwMATERIALIZED:
			return p.parseCreateMTMV(createTok.Loc)
		case kwOR:
			// CREATE OR REPLACE VIEW ...
			p.advance() // consume OR
			if _, err := p.expect(kwREPLACE); err != nil {
				return nil, err
			}
			return p.parseCreateView(createTok.Loc, true)
		case kwSTORAGE:
			p.advance() // consume STORAGE
			return p.parseCreateStorage(createTok.Loc)
		case kwREPOSITORY:
			return p.parseCreateRepository(createTok.Loc, false)
		case kwREAD:
			// CREATE READ ONLY REPOSITORY ...
			p.advance() // consume READ
			if _, err := p.expect(kwONLY); err != nil {
				return nil, err
			}
			return p.parseCreateRepository(createTok.Loc, true)
		case kwSTAGE:
			return p.parseCreateStage(createTok.Loc)
		case kwFILE:
			return p.parseCreateFile(createTok.Loc)
		case kwROW:
			p.advance() // consume ROW
			return p.parseCreateRowPolicy(createTok.Loc)
		case kwENCRYPTKEY:
			return p.parseCreateEncryptKey(createTok.Loc)
		case kwDICTIONARY:
			return p.parseCreateDictionary(createTok.Loc)
		case kwROLE:
			return p.parseCreateRole(createTok.Loc)
		case kwUSER:
			return p.parseCreateUser(createTok.Loc)
		case kwROUTINE:
			return p.parseCreateRoutineLoad(createTok.Loc)
		case kwJOB:
			return p.parseCreateJob(createTok.Loc)
		default:
			return p.unsupported("CREATE")
		}
	case kwALTER:
		alterTok := p.advance() // consume ALTER; cur is now the object type keyword
		switch p.cur.Kind {
		case kwDATABASE, kwSCHEMA:
			return p.parseAlterDatabase()
		case kwTABLE:
			return p.parseAlterTable()
		case kwVIEW:
			return p.parseAlterView()
		case kwMATERIALIZED:
			return p.parseAlterMTMV(alterTok.Loc)
		case kwCATALOG:
			return p.parseAlterCatalog()
		case kwSTORAGE:
			p.advance() // consume STORAGE
			return p.parseAlterStorage(alterTok.Loc)
		case kwREPOSITORY:
			return p.parseAlterRepository(alterTok.Loc)
		case kwWORKLOAD:
			// ALTER WORKLOAD GROUP ... or ALTER WORKLOAD POLICY ...
			alterLoc := p.prev.Loc
			next := p.peekNext()
			if next.Kind == kwGROUP {
				return p.parseAlterWorkloadGroup(alterLoc)
			}
			if next.Kind == kwPOLICY {
				return p.parseAlterWorkloadPolicy(alterLoc)
			}
			return p.unsupported("ALTER WORKLOAD")
		case kwRESOURCE:
			return p.parseAlterResource(p.prev.Loc)
		case kwSQL_BLOCK_RULE:
			return p.parseAlterSQLBlockRule(p.prev.Loc)
		case kwDICTIONARY:
			return p.parseAlterDictionary(p.prev.Loc)
		case kwROLE:
			return p.parseAlterRole(p.prev.Loc)
		case kwUSER:
			return p.parseAlterUser(p.prev.Loc)
		case kwROUTINE:
			return p.parseAlterRoutineLoad(alterTok.Loc)
		case kwSYSTEM:
			return p.parseSystemAlter(alterTok.Loc)
		case kwJOB:
			return p.parseAlterJob(p.prev.Loc)
		default:
			return p.unsupported("ALTER")
		}
	case kwDROP:
		dropTok := p.advance() // consume DROP; cur is now the object type keyword
		switch p.cur.Kind {
		case kwINDEX:
			return p.parseDropIndex(dropTok.Loc)
		case kwDATABASE, kwSCHEMA:
			return p.parseDropDatabase()
		case kwVIEW:
			return p.parseDropView(dropTok.Loc)
		case kwMATERIALIZED:
			return p.parseDropMTMV(dropTok.Loc)
		case kwCATALOG:
			return p.parseDropCatalog()
		case kwSTORAGE:
			p.advance() // consume STORAGE
			return p.parseDropStorage(dropTok.Loc)
		case kwREPOSITORY:
			return p.parseDropRepository(dropTok.Loc)
		case kwSTAGE:
			return p.parseDropStage(dropTok.Loc)
		case kwFILE:
			return p.parseDropFile(dropTok.Loc)
		case kwWORKLOAD:
			// DROP WORKLOAD GROUP ... or DROP WORKLOAD POLICY ...
			next := p.peekNext()
			if next.Kind == kwGROUP {
				return p.parseDropWorkloadGroup(dropTok.Loc)
			}
			if next.Kind == kwPOLICY {
				return p.parseDropWorkloadPolicy(dropTok.Loc)
			}
			return p.unsupported("DROP WORKLOAD")
		case kwRESOURCE:
			return p.parseDropResource(dropTok.Loc)
		case kwSQL_BLOCK_RULE:
			return p.parseDropSQLBlockRule(dropTok.Loc)
		case kwROW:
			p.advance() // consume ROW
			return p.parseDropRowPolicy(dropTok.Loc)
		case kwENCRYPTKEY:
			return p.parseDropEncryptKey(dropTok.Loc)
		case kwDICTIONARY:
			return p.parseDropDictionary(dropTok.Loc)
		case kwROLE:
			return p.parseDropRole(dropTok.Loc)
		case kwUSER:
			return p.parseDropUser(dropTok.Loc)
		case kwJOB:
			return p.parseDropJob(dropTok.Loc)
		case kwSTATS:
			p.advance() // consume STATS
			return p.parseDropStats(dropTok.Loc, "")
		case kwEXPIRED:
			p.advance() // consume EXPIRED
			if _, err := p.expect(kwSTATS); err != nil {
				return nil, err
			}
			return p.parseDropStats(dropTok.Loc, "EXPIRED")
		case kwCACHED:
			p.advance() // consume CACHED
			if _, err := p.expect(kwSTATS); err != nil {
				return nil, err
			}
			return p.parseDropStats(dropTok.Loc, "CACHED")
		default:
			return p.unsupported("DROP")
		}
	case kwTRUNCATE:
		return p.parseTruncateTable()

	// DML
	case kwSELECT:
		left, err := p.parseSelectStmt()
		if err != nil {
			return nil, err
		}
		return p.parseSetOpTail(left)
	case kwWITH:
		return p.parseWithSelect()
	case kwINSERT:
		return p.parseInsert()
	case kwUPDATE:
		updateTok := p.advance() // consume UPDATE
		return p.parseUpdateStmt(updateTok.Loc.Start)
	case kwDELETE:
		deleteTok := p.advance() // consume DELETE
		return p.parseDeleteStmt(deleteTok.Loc.Start)
	case kwMERGE:
		mergeTok := p.advance() // consume MERGE; cur is now INTO
		return p.parseMergeStmt(mergeTok.Loc)

	// Load / Export / Copy
	case kwLOAD:
		return p.parseLoad()
	case kwEXPORT:
		return p.parseExport()
	case kwCOPY:
		return p.parseCopyInto()

	// Transaction
	case kwBEGIN:
		return p.parseBeginStmt()
	case kwCOMMIT:
		return p.parseCommitStmt()
	case kwROLLBACK:
		return p.parseRollbackStmt()

	// DCL
	case kwGRANT:
		startLoc := p.cur.Loc
		p.advance()
		return p.parseGrant(startLoc)
	case kwREVOKE:
		startLoc := p.cur.Loc
		p.advance()
		return p.parseRevoke(startLoc)

	// Show / Info
	case kwSHOW:
		return p.parseShow()
		showTok := p.advance() // consume SHOW
		if p.cur.Kind == kwJOB {
			return p.parseShowJob(showTok.Loc)
		}
		// Return unsupported, but we already consumed SHOW so we need to emit the error at cur.
		return p.unsupported("SHOW")
		p.advance() // consume SHOW; dispatch inside parseShow
		return p.parseShow()
	case kwDESCRIBE:
		return p.parseDescribe()
	case kwDESC:
		return p.parseDescribe()
	case kwEXPLAIN:
		return p.parseExplain()
	case kwUSE:
		return p.parseUse()
	case kwHELP:
		return p.parseHelp()

	// Set / Unset
	case kwSET:
		setTok := p.advance() // consume SET
		// SET DEFAULT STORAGE VAULT name
		if p.cur.Kind == kwDEFAULT && p.peekNext().Kind == kwSTORAGE {
			return p.parseSetDefaultStorageVault(setTok.Loc)
		}
		// SET PASSWORD = '...'
		if p.cur.Kind == kwPASSWORD {
			return p.parseSetPassword(setTok.Loc)
		}
		// Generic variable assignment / NAMES / CHARSET / TRANSACTION
		return p.parseGenericSet(setTok.Loc)
	case kwUNSET:
		unsetTok := p.advance() // consume UNSET
		// UNSET DEFAULT STORAGE VAULT
		if p.cur.Kind == kwDEFAULT && p.peekNext().Kind == kwSTORAGE {
			return p.parseUnsetDefaultStorageVault(unsetTok.Loc)
		}
		return p.parseGenericUnset(unsetTok.Loc)

	// Admin / System
	case kwADMIN:
		adminTok := p.advance() // consume ADMIN
		return p.parseAdminStmt(adminTok.Loc)
	case kwKILL:
		killTok := p.advance() // consume KILL
		if p.cur.Kind == kwANALYZE {
			return p.parseKillAnalyze(killTok.Loc)
		}
		return p.parseKill(killTok.Loc)
	case kwLOCK:
		lockTok := p.advance() // consume LOCK
		return p.parseLockTables(lockTok.Loc)
	case kwUNLOCK:
		unlockTok := p.advance() // consume UNLOCK
		return p.parseUnlockTables(unlockTok.Loc)
	case kwINSTALL:
		installTok := p.advance() // consume INSTALL
		return p.parseInstallPlugin(installTok.Loc)
	case kwUNINSTALL:
		uninstallTok := p.advance() // consume UNINSTALL
		return p.parseUninstallPlugin(uninstallTok.Loc)

	// Backup / Restore / Recovery
	case kwBACKUP:
		backupTok := p.advance() // consume BACKUP
		return p.parseBackup(backupTok.Loc)
	case kwRESTORE:
		restoreTok := p.advance() // consume RESTORE
		return p.parseRestore(restoreTok.Loc)
	case kwRECOVER:
		recoverTok := p.advance() // consume RECOVER
		return p.parseRecover(recoverTok.Loc)

	// Materialized View / Refresh
	case kwREFRESH:
		refreshTok := p.advance() // consume REFRESH; cur is now the object type keyword
		switch p.cur.Kind {
		case kwMATERIALIZED:
			return p.parseRefreshMTMV(refreshTok.Loc)
		case kwCATALOG:
			return p.parseRefreshCatalog()
		case kwDICTIONARY:
			return p.parseRefreshDictionary(refreshTok.Loc)
		}
		return p.unsupported("REFRESH")

	// Job control
	case kwCANCEL:
		cancelTok := p.advance() // consume CANCEL
		if p.cur.Kind == kwMATERIALIZED {
			return p.parseCancelMTMVTask(cancelTok.Loc)
		}
		if p.cur.Kind == kwDECOMMISSION {
			return p.parseCancelDecommission(cancelTok.Loc)
		}
		if p.cur.Kind == kwTASK {
			return p.parseCancelTask(cancelTok.Loc)
		}
		return p.parseCancelGeneric(cancelTok.Loc)
	case kwPAUSE:
		pauseTok := p.advance() // consume PAUSE
		if p.cur.Kind == kwMATERIALIZED {
			return p.parsePauseMTMVJob(pauseTok.Loc)
		}
		if p.cur.Kind == kwROUTINE {
			return p.parsePauseRoutineLoad(pauseTok.Loc)
		}
		if p.cur.Kind == kwALL {
			p.advance() // consume ALL
			if p.cur.Kind == kwROUTINE {
				return p.parsePauseAllRoutineLoad(pauseTok.Loc)
			}
			return p.unsupported("PAUSE ALL")
		}
		if p.cur.Kind == kwJOB {
			return p.parsePauseJob(pauseTok.Loc)
		}
		return p.unsupported("PAUSE")
	case kwRESUME:
		resumeTok := p.advance() // consume RESUME
		if p.cur.Kind == kwMATERIALIZED {
			return p.parseResumeMTMVJob(resumeTok.Loc)
		}
		if p.cur.Kind == kwROUTINE {
			return p.parseResumeRoutineLoad(resumeTok.Loc)
		}
		if p.cur.Kind == kwALL {
			p.advance() // consume ALL
			if p.cur.Kind == kwROUTINE {
				return p.parseResumeAllRoutineLoad(resumeTok.Loc)
			}
			return p.unsupported("RESUME ALL")
		}
		if p.cur.Kind == kwJOB {
			return p.parseResumeJob(resumeTok.Loc)
		}
		return p.unsupported("RESUME")

	case kwSTOP:
		stopTok := p.advance() // consume STOP
		if p.cur.Kind == kwROUTINE {
			return p.parseStopRoutineLoad(stopTok.Loc)
		}
		return p.unsupported("STOP")

	// Analyze / Sync / Warm
	case kwANALYZE:
		p.advance() // consume ANALYZE
		return p.parseAnalyze()
	case kwSYNC:
		syncTok := p.advance() // consume SYNC
		return p.parseSyncStmt(syncTok.Loc)
	case kwWARM:
		warmTok := p.advance() // consume WARM
		return p.parseWarmUp(warmTok.Loc)

	// Clean
	case kwCLEAN:
		cleanTok := p.advance() // consume CLEAN
		return p.parseClean(cleanTok.Loc)

	// Index async build
	case kwBUILD:
		buildTok := p.advance() // consume BUILD
		if p.cur.Kind == kwINDEX {
			return p.parseBuildIndex(buildTok.Loc)
		}
		return p.unsupported("BUILD")

	default:
		return nil, p.unknownStatementError()
	}
}

// ParseResult holds the outcome of a best-effort parse. File contains
// all successfully-parsed statements (empty for stubs-only F4); Errors
// contains every ParseError encountered plus any LexErrors promoted from
// the underlying Lexer.
type ParseResult struct {
	File   *ast.File
	Errors []ParseError
}

// Parse is the strict entry point: returns all errors encountered while
// parsing the full input. The returned *ast.File always reflects whatever
// statements parsed successfully — even in the error case, the File may
// be non-empty.
func Parse(input string) (*ast.File, []ParseError) {
	result := ParseBestEffort(input)
	return result.File, result.Errors
}

// ParseBestEffort runs Split to segment the input, then parses each segment
// via parseSingle. Errors from individual segments are collected; all
// successfully-parsed statements are appended to the result File.
func ParseBestEffort(input string) *ParseResult {
	file := &ast.File{Loc: ast.Loc{Start: 0, End: len(input)}}
	result := &ParseResult{File: file}

	for _, seg := range Split(input) {
		node, errs := parseSingle(seg.Text, seg.ByteStart)
		if node != nil {
			file.Stmts = append(file.Stmts, node)
		}
		result.Errors = append(result.Errors, errs...)
	}

	return result
}

// parseSingle parses one segment (one top-level statement) into a single
// ast.Node. Returns (node, errors) where node may be nil if the segment
// failed to parse and errors is the list of ParseErrors encountered,
// including any LexErrors promoted from the underlying Lexer.
//
// segText is the statement text (without the trailing ; delimiter); it comes
// from Segment.Text. baseOffset is the byte offset of segText within the
// original input — passed to NewLexerWithOffset so token Loc values are
// absolute.
func parseSingle(segText string, baseOffset int) (ast.Node, []ParseError) {
	p := &Parser{
		lexer: NewLexerWithOffset(segText, baseOffset),
		input: segText,
	}
	p.advance() // prime cur with the first token

	var result ast.Node
	if p.cur.Kind != tokEOF {
		node, err := p.parseStmt()
		if err != nil {
			if pe, ok := err.(*ParseError); ok {
				p.errors = append(p.errors, *pe)
			} else {
				p.errors = append(p.errors, ParseError{
					Loc: p.cur.Loc,
					Msg: err.Error(),
				})
			}
		}
		result = node
	}

	// Promote any lex errors into ParseErrors. The Lexer's Errors() getter
	// returns positions already shifted by baseOffset.
	for _, le := range p.lexer.Errors() {
		p.errors = append(p.errors, ParseError{Loc: le.Loc, Msg: le.Msg})
	}

	return result, p.errors
}
