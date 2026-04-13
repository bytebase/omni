// Package parser — F4 parser entry point and dispatch framework.
//
// The Parser struct, Parse function, and ParseBestEffort function defined
// in this file provide the public recursive-descent parser API for
// Snowflake SQL. Actual per-statement parsing is stubbed out — F4 ships
// only the framework. Each concrete statement type (SELECT, INSERT,
// CREATE, etc.) is added by Tier 1+ DAG nodes that replace specific
// dispatch cases in parseStmt below.
//
// See docs/superpowers/specs/2026-04-09-snowflake-parser-entry-design.md
// for the design rationale.
package parser

import (
	"github.com/bytebase/omni/snowflake/ast"
)

// Parser is a recursive-descent parser for Snowflake SQL. It operates
// on a single segment of input (one top-level statement produced by
// F3's Split) at a time; callers should use Parse or ParseBestEffort
// instead of constructing Parsers directly.
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

// match tries each given token type; if cur matches any, it is consumed
// and returned with ok=true. Used for optional tokens.
func (p *Parser) match(types ...int) (Token, bool) {
	for _, t := range types {
		if p.cur.Type == t {
			return p.advance(), true
		}
	}
	return Token{}, false
}

// expect consumes the current token if it matches the expected type.
// Otherwise returns a ParseError describing the mismatch.
func (p *Parser) expect(tokenType int) (Token, error) {
	if p.cur.Type == tokenType {
		return p.advance(), nil
	}
	return Token{}, p.syntaxErrorAtCur()
}

// syntaxErrorAtCur returns a *ParseError describing a syntax error at
// the current token position. If cur is tokEOF, the message says "at
// end of input"; otherwise "at or near X".
func (p *Parser) syntaxErrorAtCur() *ParseError {
	var msg string
	if p.cur.Type == tokEOF {
		msg = "syntax error at end of input"
	} else {
		text := p.cur.Str
		if text == "" {
			text = TokenName(p.cur.Type)
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
// Because Parse uses F3.Split to pre-segment the input, each segment
// usually contains one statement — skipToNextStatement's typical behavior
// inside a segment is to advance to EOF.
func (p *Parser) skipToNextStatement() {
	// Always consume at least one token to avoid infinite loops.
	if p.cur.Type != tokEOF {
		p.advance()
	}
	depth := 0
	for p.cur.Type != tokEOF {
		switch p.cur.Type {
		case '(', '[', '{':
			depth++
		case ')', ']', '}':
			if depth > 0 {
				depth--
			}
		case ';':
			if depth == 0 {
				p.advance()
				return
			}
		}
		p.advance()
	}
}

// unsupported emits a "statement type X not yet supported" ParseError at
// the current token, advances to the next statement boundary, and
// returns (nil, err). Used by every stubbed dispatch case in F4.
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

// unknownStatementError reports a statement that starts with a token
// the dispatch switch doesn't recognize. It is called from the default
// branch of parseStmt.
func (p *Parser) unknownStatementError() *ParseError {
	if p.cur.Type == tokEOF {
		return &ParseError{
			Loc: p.cur.Loc,
			Msg: "syntax error at end of input",
		}
	}
	tokText := p.cur.Str
	if tokText == "" {
		tokText = TokenName(p.cur.Type)
	}
	return &ParseError{
		Loc: p.cur.Loc,
		Msg: "unknown or unsupported statement starting with " + tokText,
	}
}

// parseStmt parses one top-level statement. It dispatches on the first
// keyword. For stubs-only F4, every dispatch case calls the unsupported
// helper that appends a "not yet supported" ParseError and returns
// (nil, err).
//
// Tier 1+ DAG nodes REPLACE specific cases with real implementations.
// The dispatch switch itself is not expected to change — only the
// bodies of individual case branches.
//
// Keywords NOT in this switch (SAVEPOINT, EXECUTE, WHILE, REPEAT, LET,
// LOOP, BREAK) are currently absent from F2's keyword table because
// they're either missing from the legacy grammar or defined via
// non-standard ANTLR rule forms. When Tier 7 Snowflake Scripting support
// lands, F2 gains those constants and this switch gains matching cases.
func (p *Parser) parseStmt() (ast.Node, error) {
	switch p.cur.Type {
	// DDL (6 cases)
	case kwCREATE:
		return p.unsupported("CREATE")
	case kwALTER:
		return p.unsupported("ALTER")
	case kwDROP:
		return p.unsupported("DROP")
	case kwUNDROP:
		return p.unsupported("UNDROP")
	case kwTRUNCATE:
		return p.unsupported("TRUNCATE")
	case kwCOMMENT:
		return p.unsupported("COMMENT")

	// DML (12 cases)
	case kwSELECT:
		return p.parseSelectStmt()
	case kwWITH:
		return p.parseWithSelect()
	case kwINSERT:
		return p.unsupported("INSERT")
	case kwUPDATE:
		return p.unsupported("UPDATE")
	case kwDELETE:
		return p.unsupported("DELETE")
	case kwMERGE:
		return p.unsupported("MERGE")
	case kwCOPY:
		return p.unsupported("COPY")
	case kwPUT:
		return p.unsupported("PUT")
	case kwGET:
		return p.unsupported("GET")
	case kwLIST:
		return p.unsupported("LIST")
	case kwREMOVE:
		return p.unsupported("REMOVE")
	case kwCALL:
		return p.unsupported("CALL")

	// DCL (2 cases)
	case kwGRANT:
		return p.unsupported("GRANT")
	case kwREVOKE:
		return p.unsupported("REVOKE")

	// TCL (4 cases)
	case kwBEGIN:
		return p.unsupported("BEGIN")
	case kwSTART:
		return p.unsupported("START")
	case kwCOMMIT:
		return p.unsupported("COMMIT")
	case kwROLLBACK:
		return p.unsupported("ROLLBACK")

	// Info / inspection (5 cases)
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

	// Session (2 cases)
	case kwSET:
		return p.unsupported("SET")
	case kwUNSET:
		return p.unsupported("UNSET")

	// Snowflake Scripting (subset available as F2 keywords — 6 cases)
	case kwDECLARE:
		return p.unsupported("DECLARE")
	case kwIF:
		return p.unsupported("IF")
	case kwCASE:
		return p.unsupported("CASE")
	case kwFOR:
		return p.unsupported("FOR")
	case kwRETURN:
		return p.unsupported("RETURN")
	case kwCONTINUE:
		return p.unsupported("CONTINUE")

	default:
		return nil, p.unknownStatementError()
	}
}

// ParseResult holds the outcome of a best-effort parse. File contains
// all successfully-parsed statements (empty for stubs-only F4); Errors
// contains every ParseError encountered plus any LexErrors promoted
// from the underlying Lexer.
type ParseResult struct {
	File   *ast.File
	Errors []ParseError
}

// Parse is the strict entry point: returns the first error if any
// statement fails. Callers that want partial results should use
// ParseBestEffort instead.
//
// The returned *ast.File always reflects whatever statements parsed
// successfully before the first error — even in the error case, the
// File may be non-empty.
func Parse(input string) (*ast.File, error) {
	result := ParseBestEffort(input)
	if len(result.Errors) > 0 {
		return result.File, &result.Errors[0]
	}
	return result.File, nil
}

// ParseBestEffort runs F3's Split to segment the input, then parses each
// segment via parseSingle. Errors from individual segments are collected;
// all successfully-parsed statements are appended to the result File.
//
// This is the canonical entry point for bytebase consumers (diagnose.go,
// query_type.go, query_span_extractor.go) that need partial results plus
// error diagnostics.
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
// segText is the statement text (without the trailing ; delimiter); it
// comes from F3.Segment.Text. baseOffset is the byte offset of segText
// within the original input — passed to NewLexerWithOffset so token Loc
// values are absolute.
func parseSingle(segText string, baseOffset int) (ast.Node, []ParseError) {
	p := &Parser{
		lexer: NewLexerWithOffset(segText, baseOffset),
		input: segText,
	}
	p.advance() // prime cur with the first token

	var result ast.Node
	if p.cur.Type != tokEOF {
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

	// Promote any lex errors into ParseErrors. The Lexer's Errors()
	// getter returns positions already shifted by baseOffset.
	for _, le := range p.lexer.Errors() {
		p.errors = append(p.errors, ParseError{Loc: le.Loc, Msg: le.Msg})
	}

	return result, p.errors
}
