// Package parser implements a hand-written recursive-descent parser for Trino
// SQL (Trino 481 grammar, derived from the official SqlBase.g4).
//
// This file is the parser-foundation node: it ships the Parser struct, the
// token-cursor primitives (advance / peek / peekNext / match / expect), the
// public Parse / ParseBestEffort entry points, the per-segment driver
// (parseSingle), and the top-level statement-dispatch switch (parseStmt).
//
// Per-statement parsing is STUBBED here — the foundation ships only the
// dispatch framework. Every dispatch case routes to a real parse function or,
// where that function does not yet exist, to the unsupported helper which
// records a "not yet supported" ParseError. Later DAG nodes (types,
// expressions, parser-select, parser-ddl, parser-dml, …) REPLACE individual
// dispatch-case bodies with concrete parsers that build real ast.Node trees.
//
// The dispatch switch enumerates every first-keyword of Trino's `statement`
// and `rootQuery`/`queryPrimary` grammar rules so that bytebase's Diagnose —
// which runs on every statement — never emits a false "unknown statement"
// diagnostic for syntactically valid Trino. Reaching a stubbed case is correct
// foundation behavior; reaching the default (unknown) branch for a valid
// statement is a bug.
//
// The package mirrors omni's doris/parser and snowflake/parser conventions: a
// stateless Lexer feeds the Parser one Token{Kind, Str, Ival, Loc} at a time
// with two-token lookahead, and source positions are byte offsets.
package parser

import (
	"github.com/bytebase/omni/trino/ast"
)

// Parser is a recursive-descent parser for Trino SQL. It operates on a single
// segment of input (one top-level statement produced by Split) at a time;
// callers should use Parse or ParseBestEffort instead of constructing Parsers
// directly.
type Parser struct {
	lexer      *Lexer
	input      string       // the segment text (used for error reporting)
	baseOffset int          // absolute offset of input within the original source
	cur        Token        // current token
	prev       Token        // previous token (for error context)
	nextBuf    Token        // buffered lookahead token
	hasNext    bool         // whether nextBuf is valid
	errors     []ParseError // collected errors for best-effort mode
}

// advance consumes the current token and moves to the next one. Returns the
// token that was just consumed (the new "previous" token).
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
// (one-token lookahead beyond cur). Used to disambiguate based on the token
// following the current position.
func (p *Parser) peekNext() Token {
	if !p.hasNext {
		p.nextBuf = p.lexer.NextToken()
		p.hasNext = true
	}
	return p.nextBuf
}

// match tries each given token kind; if cur matches any, it is consumed and
// returned with ok == true. Used for optional tokens.
func (p *Parser) match(kinds ...TokenKind) (Token, bool) {
	for _, k := range kinds {
		if p.cur.Kind == k {
			return p.advance(), true
		}
	}
	return Token{}, false
}

// expect consumes the current token if it matches the expected kind. Otherwise
// returns a ParseError describing the mismatch (without consuming anything).
func (p *Parser) expect(kind TokenKind) (Token, error) {
	if p.cur.Kind == kind {
		return p.advance(), nil
	}
	return Token{}, p.syntaxErrorAtCur()
}

// syntaxErrorAtCur returns a *ParseError describing a syntax error at the
// current token position. At EOF the message says "at end of input"; otherwise
// "at or near X", where X is the token's source text or, for value-less tokens
// (operators, EOF), its symbolic name.
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

// skipToNextStatement advances the parser past the current erroneous statement
// to the next ';' at bracket-depth 0 within the current segment, or to tokEOF.
// Always consumes at least one token to guarantee forward progress.
//
// Because Parse pre-segments the input with Split, each segment usually
// contains a single statement, so skipToNextStatement's typical behavior is to
// advance to EOF.
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

// unsupported emits a "<name> statement parsing is not yet supported"
// ParseError at the current token, skips to the next statement boundary, and
// returns (nil, err). Used by every stubbed dispatch case in the foundation.
//
// Later DAG nodes REPLACE individual dispatch cases with real parse functions.
// Those functions do NOT call unsupported — they consume tokens and build real
// ast.Node values.
func (p *Parser) unsupported(name string) (ast.Node, error) {
	err := &ParseError{
		Loc: p.cur.Loc,
		Msg: name + " statement parsing is not yet supported",
	}
	p.skipToNextStatement()
	return nil, err
}

// unknownStatementError reports a statement that starts with a token the
// dispatch switch does not recognize. Called from the default branch of
// parseStmt. Distinct from unsupported so Diagnose can tell "valid Trino, not
// yet implemented" apart from "not valid Trino".
func (p *Parser) unknownStatementError() *ParseError {
	if p.cur.Kind == tokEOF {
		return &ParseError{
			Loc: p.cur.Loc,
			Msg: "syntax error at end of input",
		}
	}
	text := p.cur.Str
	if text == "" {
		text = TokenName(p.cur.Kind)
	}
	return &ParseError{
		Loc: p.cur.Loc,
		Msg: "unknown or unsupported statement starting with " + text,
	}
}

// parseStmt parses one top-level statement by dispatching on the leading
// keyword(s). The arms enumerate the first-token vocabulary of Trino's
// `statement` rule plus the query-starting tokens of `rootQuery` /
// `queryPrimary` (SELECT, WITH, TABLE, VALUES, '('). For the foundation every
// arm routes to the unsupported stub; later nodes replace the bodies.
//
// CREATE / DROP / ALTER / SET / RESET / COMMENT / DESCRIBE / SHOW are
// second-keyword-dispatched in the grammar; the foundation keeps the dispatch
// flat (one stub per leading keyword) because the per-object parsers do not
// exist yet. Later nodes will introduce the inner switches (compare
// doris/parser.go's CREATE/DROP/ALTER fan-out) as they add object parsers.
func (p *Parser) parseStmt() (ast.Node, error) {
	switch p.cur.Kind {
	// --- Query (rootQuery / queryPrimary leading tokens) ---
	// SELECT / WITH / TABLE / VALUES / '(' all begin a query; parser-select's
	// parseQueryStmt parses the whole query layer. (A leading `WITH FUNCTION`
	// inline routine is the parser-routines node; parseQuery reports it as
	// unsupported until that node lands.)
	case kwSELECT, kwWITH, kwTABLE, kwVALUES, int('('):
		return p.parseQueryStmt()

	// --- Session / catalog context ---
	case kwUSE:
		return p.parseUseStmt()

	// --- DDL ---
	case kwCREATE:
		// CREATE ROLE belongs to the parser-dcl-tcl node (grant_revoke.go);
		// every other CREATE object is a parser-ddl concern (still stubbed).
		if p.peekNext().Kind == kwROLE {
			createTok := p.advance() // consume CREATE
			p.advance()              // consume ROLE
			return p.parseCreateRoleStmt(createTok.Loc.Start)
		}
		return p.unsupported("CREATE")
	case kwDROP:
		// DROP ROLE belongs to the parser-dcl-tcl node (grant_revoke.go);
		// every other DROP object is a parser-ddl concern (still stubbed).
		if p.peekNext().Kind == kwROLE {
			dropTok := p.advance() // consume DROP
			p.advance()            // consume ROLE
			return p.parseDropRoleStmt(dropTok.Loc.Start)
		}
		return p.unsupported("DROP")
	case kwALTER:
		return p.unsupported("ALTER")
	case kwCOMMENT:
		return p.unsupported("COMMENT")
	case kwANALYZE:
		return p.unsupported("ANALYZE")
	case kwREFRESH:
		return p.unsupported("REFRESH")

	// --- DML --- [parser-dml node: insert.go, update_delete.go, merge.go]
	case kwINSERT:
		return p.parseInsertStmt()
	case kwDELETE:
		return p.parseDeleteStmt()
	case kwUPDATE:
		return p.parseUpdateStmt()
	case kwMERGE:
		return p.parseMergeStmt()
	case kwTRUNCATE:
		return p.parseTruncateStmt()

	// --- Procedures ---
	case kwCALL:
		return p.parseCallStmt()

	// --- DCL (roles / privileges) --- [parser-dcl-tcl node: grant_revoke.go]
	case kwGRANT:
		return p.parseGrantStmt()
	case kwREVOKE:
		return p.parseRevokeStmt()
	case kwDENY:
		return p.parseDenyStmt()

	// --- Session / path / time zone ---
	case kwSET:
		// SET fans out on the second keyword: SESSION (setSession /
		// setSessionAuthorization), ROLE (setRole), PATH (setPath), TIME
		// (setTimeZone). See session.go.
		return p.parseSetStmt()
	case kwRESET:
		// RESET SESSION { AUTHORIZATION | name }. See session.go.
		return p.parseResetStmt()

	// --- Inspection ---
	case kwSHOW:
		// SHOW family + DESCRIBE/DESC aliases. See show.go.
		return p.parseShowStmt()
	case kwEXPLAIN:
		// EXPLAIN / EXPLAIN ANALYZE. See explain_call.go.
		return p.parseExplainStmt()
	case kwDESCRIBE:
		// DESCRIBE INPUT/OUTPUT <name> are prepared-statement introspection
		// (parser-dcl-tcl: prepared.go), claimed only when INPUT/OUTPUT is
		// followed by a statement name. INPUT/OUTPUT are NON-RESERVED, so a
		// bare `DESCRIBE INPUT` or `DESCRIBE INPUT.col` is the SHOW COLUMNS
		// alias `DESCRIBE qualifiedName` (INPUT/OUTPUT as the table name).
		if (p.peekNext().Kind == kwINPUT || p.peekNext().Kind == kwOUTPUT) && p.describeIntrospectionFollows() {
			descTok := p.advance() // consume DESCRIBE
			kind := p.advance()    // consume INPUT / OUTPUT
			if kind.Kind == kwINPUT {
				return p.parseDescribeInputStmt(descTok.Loc.Start)
			}
			return p.parseDescribeOutputStmt(descTok.Loc.Start)
		}
		// DESCRIBE name == SHOW COLUMNS (parser-utility). See show.go.
		return p.parseDescribeStmt()
	case kwDESC:
		// DESC name == SHOW COLUMNS; DESC has no prepared form.
		return p.parseDescribeStmt()

	// --- Transactions --- [parser-dcl-tcl node: transaction.go]
	case kwSTART:
		return p.parseStartTransactionStmt()
	case kwCOMMIT:
		return p.parseCommitStmt()
	case kwROLLBACK:
		return p.parseRollbackStmt()

	// --- Prepared statements --- [parser-dcl-tcl node: prepared.go]
	case kwPREPARE:
		return p.parsePrepareStmt()
	case kwDEALLOCATE:
		return p.parseDeallocateStmt()
	case kwEXECUTE:
		return p.parseExecuteStmt()

	default:
		return nil, p.unknownStatementError()
	}
}

// ParseResult holds the outcome of a best-effort parse. File contains every
// statement that parsed successfully (empty while statement bodies are
// stubbed); Errors contains every ParseError encountered, including LexErrors
// promoted from the underlying Lexer.
type ParseResult struct {
	File   *ast.File
	Errors []ParseError
}

// Parse is the public entry point. It returns the parsed File plus every error
// encountered. The File always reflects whatever statements parsed
// successfully — even in the error case it may be non-empty.
//
// The signature returns all errors (matching doris/parser.Parse) rather than a
// single error: bytebase's Diagnose needs the complete diagnostic set, and a
// multi-statement script can fail in several places at once.
func Parse(input string) (*ast.File, []ParseError) {
	result := ParseBestEffort(input)
	return result.File, result.Errors
}

// ParseBestEffort runs Split to segment the input, then parses each segment via
// parseSingle. Per-segment errors are collected; every successfully-parsed
// statement is appended to the result File. This is the canonical entry point
// for the bytebase consumers (Diagnose, query-type classification, query-span
// extraction) that need partial results plus diagnostics.
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
// ast.Node. Returns (node, errors) where node may be nil if the segment failed
// to parse and errors lists every ParseError encountered, including LexErrors
// promoted from the underlying Lexer.
//
// segText is the statement text without the trailing ';' (from Segment.Text).
// baseOffset is segText's byte offset within the original input; it is passed
// to NewLexerWithOffset so token and error Loc values stay absolute.
func parseSingle(segText string, baseOffset int) (ast.Node, []ParseError) {
	p := &Parser{
		lexer:      NewLexerWithOffset(segText, baseOffset),
		input:      segText,
		baseOffset: baseOffset,
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
		} else if p.cur.Kind != tokEOF {
			// A statement parsed cleanly but did not consume the whole segment:
			// the leftover tokens are a syntax error (e.g. `SELECT a a a`,
			// `SELECT 1 garbage`). Each segment holds exactly one top-level
			// statement (Split cut on ';'), so any trailing token is invalid here.
			p.errors = append(p.errors, *p.syntaxErrorAtCur())
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
