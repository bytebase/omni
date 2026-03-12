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
		p.advance() // consume RENAME
		return p.parseRenameTableStmt()

	case kwSET:
		return p.parseSetStmt()

	case kwSHOW:
		return p.parseShowStmt()

	case kwUSE:
		return p.parseUseStmt()

	case kwDESCRIBE, kwEXPLAIN:
		return p.parseExplainStmt()

	case kwBEGIN, kwSTART:
		return p.parseBeginStmt()

	case kwCOMMIT:
		return p.parseCommitStmt()

	case kwROLLBACK:
		return p.parseRollbackStmt()

	case kwSAVEPOINT:
		return p.parseSavepointStmt()

	case kwRELEASE:
		return p.parseReleaseSavepointStmt()

	case kwLOCK:
		p.advance() // consume LOCK
		return p.parseLockTablesStmt()

	case kwUNLOCK:
		p.advance() // consume UNLOCK
		return p.parseUnlockTablesStmt()

	case kwGRANT:
		return p.parseGrantStmt()

	case kwREVOKE:
		return p.parseRevokeStmt()

	case kwLOAD:
		return p.parseLoadDataStmt()

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
		return p.parseResetStmt()

	case kwKILL:
		return p.parseKillStmt()

	case kwDO:
		return p.parseDoStmt()

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

// parseCreateDispatch dispatches CREATE statements to the appropriate parser.
func (p *Parser) parseCreateDispatch() (nodes.Node, error) {
	p.advance() // consume CREATE

	orReplace := false
	temporary := false

	// OR REPLACE
	if p.cur.Type == kwOR {
		p.advance()
		p.match(kwREPLACE)
		orReplace = true
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

	default:
		return nil, &ParseError{
			Message:  "unexpected token after CREATE",
			Position: p.cur.Loc,
		}
	}
}

// parseAlterDispatch dispatches ALTER statements to the appropriate parser.
func (p *Parser) parseAlterDispatch() (nodes.Node, error) {
	p.advance() // consume ALTER

	switch p.cur.Type {
	case kwTABLE:
		return p.parseAlterTableStmt()

	case kwDATABASE, kwSCHEMA:
		return p.parseAlterDatabaseStmt()

	case kwUSER:
		return p.parseAlterUserStmt()

	default:
		return nil, &ParseError{
			Message:  "unexpected token after ALTER",
			Position: p.cur.Loc,
		}
	}
}

// parseDropDispatch dispatches DROP statements to the appropriate parser.
func (p *Parser) parseDropDispatch() (nodes.Node, error) {
	p.advance() // consume DROP

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

	case kwFUNCTION, kwPROCEDURE:
		// DROP FUNCTION/PROCEDURE — consume and skip remaining tokens
		p.advance()
		p.match(kwIF)
		p.match(kwEXISTS_KW)
		if p.isIdentToken() {
			p.parseIdentifier()
		}
		return &nodes.RawStmt{}, nil

	case kwTRIGGER:
		p.advance()
		p.match(kwIF)
		p.match(kwEXISTS_KW)
		if p.isIdentToken() {
			p.parseIdentifier()
		}
		return &nodes.RawStmt{}, nil

	case kwEVENT:
		p.advance()
		p.match(kwIF)
		p.match(kwEXISTS_KW)
		if p.isIdentToken() {
			p.parseIdentifier()
		}
		return &nodes.RawStmt{}, nil

	default:
		return nil, &ParseError{
			Message:  "unexpected token after DROP",
			Position: p.cur.Loc,
		}
	}
}
