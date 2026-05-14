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
		case kwTABLE, kwEXTERNAL, kwTEMPORARY:
			return p.parseCreateTable()
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
		default:
			return p.unsupported("DROP")
		}
	case kwTRUNCATE:
		return p.unsupported("TRUNCATE")

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
		return p.unsupported("LOAD")
	case kwEXPORT:
		return p.unsupported("EXPORT")
	case kwCOPY:
		return p.unsupported("COPY")

	// Transaction
	case kwBEGIN:
		return p.unsupported("BEGIN")
	case kwCOMMIT:
		return p.unsupported("COMMIT")
	case kwROLLBACK:
		return p.unsupported("ROLLBACK")

	// DCL
	case kwGRANT:
		return p.unsupported("GRANT")
	case kwREVOKE:
		return p.unsupported("REVOKE")

	// Show / Info
	case kwSHOW:
		return p.unsupported("SHOW")
	case kwDESCRIBE:
		return p.unsupported("DESCRIBE")
	case kwDESC:
		return p.unsupported("DESC")
	case kwEXPLAIN:
		return p.unsupported("EXPLAIN")
	case kwUSE:
		return p.unsupported("USE")
	case kwHELP:
		return p.unsupported("HELP")

	// Set / Unset
	case kwSET:
		return p.unsupported("SET")
	case kwUNSET:
		return p.unsupported("UNSET")

	// Admin / System
	case kwADMIN:
		return p.unsupported("ADMIN")
	case kwKILL:
		return p.unsupported("KILL")
	case kwLOCK:
		return p.unsupported("LOCK")
	case kwUNLOCK:
		return p.unsupported("UNLOCK")
	case kwINSTALL:
		return p.unsupported("INSTALL")
	case kwUNINSTALL:
		return p.unsupported("UNINSTALL")

	// Backup / Restore / Recovery
	case kwBACKUP:
		return p.unsupported("BACKUP")
	case kwRESTORE:
		return p.unsupported("RESTORE")
	case kwRECOVER:
		return p.unsupported("RECOVER")

	// Materialized View / Refresh
	case kwREFRESH:
		refreshTok := p.advance() // consume REFRESH
		if p.cur.Kind == kwMATERIALIZED {
			return p.parseRefreshMTMV(refreshTok.Loc)
		}
		return p.unsupported("REFRESH")

	// Job control
	case kwCANCEL:
		cancelTok := p.advance() // consume CANCEL
		if p.cur.Kind == kwMATERIALIZED {
			return p.parseCancelMTMVTask(cancelTok.Loc)
		}
		return p.unsupported("CANCEL")
	case kwPAUSE:
		pauseTok := p.advance() // consume PAUSE
		if p.cur.Kind == kwMATERIALIZED {
			return p.parsePauseMTMVJob(pauseTok.Loc)
		}
		return p.unsupported("PAUSE")
	case kwRESUME:
		resumeTok := p.advance() // consume RESUME
		if p.cur.Kind == kwMATERIALIZED {
			return p.parseResumeMTMVJob(resumeTok.Loc)
		}
		return p.unsupported("RESUME")

	// Analyze / Sync / Warm
	case kwANALYZE:
		return p.unsupported("ANALYZE")
	case kwSYNC:
		return p.unsupported("SYNC")
	case kwWARM:
		return p.unsupported("WARM")

	// Clean
	case kwCLEAN:
		return p.unsupported("CLEAN")

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
