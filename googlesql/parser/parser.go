// Package parser implements a hand-written recursive-descent parser for
// GoogleSQL — the SQL dialect shared by Google BigQuery and Google Cloud
// Spanner. One omni parser serves both engines; the grammar is a hand-port of
// Google's open-source ZetaSQL reference (the same lineage as the legacy ANTLR
// GoogleSQLParser.g4, 625 rules).
//
// This file is the parser-foundation node. It ships:
//   - the Parser struct + token-cursor primitives (advance / peek / peekNext /
//     match / expect),
//   - the public Parse / ParseBestEffort entry points and the per-segment
//     driver (parseSingle),
//   - the leading statement-level-hint skip (@{...} / @N), and
//   - the top-level statement-dispatch switch (parseStmt) enumerating every
//     first-keyword of GoogleSQL's sql_statement_body.
//
// Per-statement parsing is STUBBED here — the foundation ships only the
// dispatch framework. Every dispatch case routes to a real parse function or,
// where that function does not yet exist, to the unsupported helper which
// records a "not yet supported" ParseError. Later DAG nodes (types,
// expressions, parser-select, parser-ddl, parser-dml, …) REPLACE individual
// dispatch-case bodies with concrete parsers that build real ast.Node trees.
//
// The dispatch switch enumerates every first-keyword of the grammar's
// sql_statement_body so bytebase's Diagnose — which runs on every statement —
// never emits a false "unknown statement" diagnostic for syntactically valid
// GoogleSQL. Reaching a stubbed case is correct foundation behavior; reaching
// the default (unknown) branch for a valid statement is a bug.
//
// The package mirrors omni's snowflake/parser and trino/parser conventions: a
// stateless Lexer feeds the Parser one Token{Type, Str, Ival, Loc} at a time
// with two-token lookahead, and source positions are byte offsets.
package parser

import (
	"github.com/bytebase/omni/googlesql/ast"
)

// Parser is a recursive-descent parser for GoogleSQL. It operates on a single
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

// match tries each given token type; if cur matches any, it is consumed and
// returned with ok == true. Used for optional tokens.
func (p *Parser) match(types ...int) (Token, bool) {
	for _, t := range types {
		if p.cur.Type == t {
			return p.advance(), true
		}
	}
	return Token{}, false
}

// expect consumes the current token if it matches the expected type. Otherwise
// returns a ParseError describing the mismatch (without consuming anything).
func (p *Parser) expect(tokenType int) (Token, error) {
	if p.cur.Type == tokenType {
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

// skipToNextStatement advances the parser past the current erroneous statement
// to the next ';' at bracket-depth 0 within the current segment, or to tokEOF.
// Always consumes at least one token to guarantee forward progress.
//
// Because Parse pre-segments the input with Split, each segment usually
// contains a single statement, so skipToNextStatement's typical behavior is to
// advance to EOF.
func (p *Parser) skipToNextStatement() {
	// Always consume at least one token to avoid infinite loops.
	if p.cur.Type != tokEOF {
		p.advance()
	}
	depth := 0
	for p.cur.Type != tokEOF {
		switch p.cur.Type {
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
// parseStmt. Distinct from unsupported so Diagnose can tell "valid GoogleSQL,
// not yet implemented" apart from "not valid GoogleSQL".
func (p *Parser) unknownStatementError() *ParseError {
	if p.cur.Type == tokEOF {
		return &ParseError{
			Loc: p.cur.Loc,
			Msg: "syntax error at end of input",
		}
	}
	text := p.cur.Str
	if text == "" {
		text = TokenName(p.cur.Type)
	}
	return &ParseError{
		Loc: p.cur.Loc,
		Msg: "unknown or unsupported statement starting with " + text,
	}
}

// skipStatementLevelHint consumes a leading statement-level hint, if present,
// so dispatch sees the real statement keyword. GoogleSQL allows an optional
// statement_level_hint before sql_statement_body:
//
//	@<int>                         simple hint
//	@{ entry, ... }                hint body
//	@[ <int> @] { entry, ... }     hint body with a leading int
//
// (The lexer emits '@' as the single-char AT_SYMBOL token; '@@' — the
// system-variable prefix — is a different token, tokAtAt, and is NOT a hint.)
//
// The foundation only SKIPS the hint to reach the statement keyword; building
// the hint AST is owned by a later node. Skipping is brace-balanced so a ';'
// inside the hint body never leaks (Split already keeps the hint with its
// statement). skipStatementLevelHint only advances over a well-formed-looking
// prefix and otherwise leaves cur untouched.
//
// It returns a *ParseError when the hint body is UNTERMINATED — an `@{` (or
// `@[…@]{`) whose closing '}' never arrives before EOF. Without this, an
// unterminated hint would be silently consumed to EOF and the (invalid)
// statement would draw no diagnostic at all, which bytebase's Diagnose must not
// allow (oracle: Spanner rejects a malformed hint such as `@{}` / `@` with a
// syntax error). Returns nil when there is no hint or the hint body closes
// cleanly.
func (p *Parser) skipStatementLevelHint() *ParseError {
	if p.cur.Type != int('@') {
		return nil
	}
	hintStart := p.cur.Loc
	next := p.peekNext()
	switch {
	case next.Type == int('{'):
		p.advance() // consume '@'
		return p.skipBalancedBraces(hintStart)
	case next.Type == int('['):
		// @[ int @] { ... }
		p.advance() // consume '@'
		p.advance() // consume '['
		// Skip up to the matching @] then the brace body.
		for p.cur.Type != tokEOF && p.cur.Type != int(']') {
			p.advance()
		}
		if p.cur.Type != int(']') {
			return &ParseError{Loc: hintStart, Msg: "unterminated statement hint"}
		}
		p.advance() // consume ']'
		if p.cur.Type == int('{') {
			return p.skipBalancedBraces(hintStart)
		}
		// `@[ … @]` with no following `{` body: a hint preamble must carry a
		// brace body, so this is malformed.
		return &ParseError{Loc: hintStart, Msg: "unterminated statement hint"}
	case next.Type == tokInteger:
		p.advance() // consume '@'
		p.advance() // consume the int
	}
	return nil
}

// skipBalancedBraces consumes a '{' ... '}' run with balanced nesting starting
// at the current token (which must be '{'). Returns a *ParseError anchored at
// hintStart if EOF is reached before the matching '}' (unterminated hint);
// returns nil on a clean close or when cur is not '{'.
func (p *Parser) skipBalancedBraces(hintStart ast.Loc) *ParseError {
	if p.cur.Type != int('{') {
		return nil
	}
	depth := 0
	for p.cur.Type != tokEOF {
		switch p.cur.Type {
		case int('{'):
			depth++
		case int('}'):
			depth--
			if depth == 0 {
				p.advance()
				return nil
			}
		}
		p.advance()
	}
	return &ParseError{Loc: hintStart, Msg: "unterminated statement hint"}
}

// parseStmt parses one top-level statement by dispatching on the leading
// keyword(s). The arms enumerate the first-token vocabulary of GoogleSQL's
// sql_statement_body (antlr_rules.md §4): query (SELECT / WITH / GRAPH / '(' /
// statement-hinted), DDL, DML, DCL, transactions, batch, and utility/metadata
// statements. For the foundation every arm routes to the unsupported stub;
// later nodes replace the bodies.
//
// CREATE / DROP / ALTER / SET / SHOW / DESCRIBE / EXPORT etc. are
// second-keyword-dispatched in the grammar; the foundation keeps the dispatch
// flat (one stub per leading keyword) because the per-object parsers do not
// exist yet. Later nodes introduce the inner switches (compare doris/parser.go
// CREATE/DROP/ALTER fan-out) as they add object parsers.
//
// Procedural-script statements (IF / CASE / WHILE / LOOP / REPEAT / FOR /
// DECLARE / BREAK / LEAVE / CONTINUE / ITERATE / RETURN / RAISE, BEGIN…END
// block) are legal only inside a procedural body (statement_list), not as a
// top-level sql_statement_body — EXCEPT the BEGIN…END block, which is itself a
// script statement that can be the body of a script. The dispatcher recognizes
// them at top level so a standalone script fragment fed to Diagnose is reported
// as "not yet supported" rather than "unknown statement".
func (p *Parser) parseStmt() (ast.Node, error) {
	switch p.cur.Type {
	// --- Query (query_statement / GQL) ---
	case kwSELECT:
		return p.unsupported("SELECT")
	case kwWITH:
		// WITH cte... query  (also WITH RECURSIVE, and the WITH-pipe error alt).
		return p.unsupported("WITH")
	case kwGRAPH:
		// GQL graph query: GRAPH <name> <gql ops>.
		return p.unsupported("GRAPH")
	case int('('):
		// Parenthesized query at top level, e.g. (SELECT 1) UNION ALL (SELECT 2).
		return p.unsupported("query")
	case kwFROM:
		// FROM-first query head. The legacy grammar recognizes a leading FROM as
		// a query (query_without_pipe_operators' `with_clause? from_clause` alt)
		// and then emits "Syntax error: Unexpected FROM" — i.e. FROM is a query-
		// statement head that the query rule itself rejects (oracle-confirmed:
		// Spanner rejects `FROM t` as a syntax error). Recognizing it here (vs the
		// unknown branch) keeps the head with the query family; the precise
		// "Unexpected FROM" rejection is produced by the query node that replaces
		// this case. See the parser-foundation divergence ledger.
		return p.unsupported("FROM")

	// --- DDL ---
	case kwCREATE:
		return p.unsupported("CREATE")
	case kwALTER:
		return p.unsupported("ALTER")
	case kwDROP:
		return p.unsupported("DROP")
	case kwRENAME:
		return p.unsupported("RENAME")
	case kwUNDROP:
		return p.unsupported("UNDROP")
	case kwTRUNCATE:
		return p.unsupported("TRUNCATE")
	case kwDEFINE:
		// DEFINE TABLE (DEFINE MACRO is intentionally unimplemented — a
		// reserved error alt in the legacy grammar).
		return p.unsupported("DEFINE")

	// --- DML ---
	case kwINSERT:
		return p.unsupported("INSERT")
	case kwUPDATE:
		return p.unsupported("UPDATE")
	case kwDELETE:
		return p.unsupported("DELETE")
	case kwMERGE:
		return p.unsupported("MERGE")

	// --- DCL ---
	case kwGRANT:
		return p.unsupported("GRANT")
	case kwREVOKE:
		return p.unsupported("REVOKE")

	// --- Transactions / batch / session ---
	case kwBEGIN:
		// BEGIN [TRANSACTION] (TCL) OR a procedural BEGIN…END block.
		return p.unsupported("BEGIN")
	case kwSTART:
		// START TRANSACTION | START BATCH.
		return p.unsupported("START")
	case kwCOMMIT:
		return p.unsupported("COMMIT")
	case kwROLLBACK:
		return p.unsupported("ROLLBACK")
	case kwSET:
		return p.unsupported("SET")
	case kwRUN:
		return p.unsupported("RUN BATCH")
	case kwABORT:
		return p.unsupported("ABORT BATCH")

	// --- Utility / metadata ---
	case kwEXPLAIN:
		return p.unsupported("EXPLAIN")
	case kwDESCRIBE:
		return p.unsupported("DESCRIBE")
	case kwDESC:
		return p.unsupported("DESC")
	case kwSHOW:
		return p.unsupported("SHOW")
	case kwANALYZE:
		return p.unsupported("ANALYZE")
	case kwASSERT:
		return p.unsupported("ASSERT")
	case kwCALL:
		return p.unsupported("CALL")
	case kwEXECUTE:
		// EXECUTE IMMEDIATE.
		return p.unsupported("EXECUTE")
	case kwIMPORT:
		return p.unsupported("IMPORT")
	case kwMODULE:
		return p.unsupported("MODULE")
	case kwEXPORT:
		// EXPORT DATA | EXPORT MODEL | EXPORT … METADATA.
		return p.unsupported("EXPORT")
	case kwLOAD:
		// LOAD DATA FROM FILES.
		return p.unsupported("LOAD")
	case kwCLONE:
		// CLONE DATA.
		return p.unsupported("CLONE")

	// --- Procedural / scripting (legal inside a script body; recognized at
	// top level so a script fragment is reported as unsupported, not unknown) ---
	case kwIF:
		return p.unsupported("IF")
	case kwCASE:
		return p.unsupported("CASE")
	case kwWHILE:
		return p.unsupported("WHILE")
	case kwLOOP:
		return p.unsupported("LOOP")
	case kwREPEAT:
		return p.unsupported("REPEAT")
	case kwFOR:
		return p.unsupported("FOR")
	case kwDECLARE:
		return p.unsupported("DECLARE")
	case kwBREAK:
		return p.unsupported("BREAK")
	case kwLEAVE:
		return p.unsupported("LEAVE")
	case kwCONTINUE:
		return p.unsupported("CONTINUE")
	case kwITERATE:
		return p.unsupported("ITERATE")
	case kwRETURN:
		return p.unsupported("RETURN")
	case kwRAISE:
		return p.unsupported("RAISE")

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
// The signature returns all errors (matching snowflake/trino parser.Parse)
// rather than a single error: bytebase's Diagnose needs the complete diagnostic
// set, and a multi-statement script can fail in several places at once.
func Parse(input string) (*ast.File, []ParseError) {
	result := ParseBestEffort(input)
	return result.File, result.Errors
}

// ParseBestEffort runs Split to segment the input, then parses each segment via
// parseSingle. Per-segment parse errors are collected; every successfully-parsed
// statement is appended to the result File. This is the canonical entry point
// for the bytebase consumers (Diagnose, query-type classification, query-span
// extraction) that need partial results plus diagnostics.
//
// Split (the block-aware variant) is used so a procedural BEGIN/END body is fed
// to parseSingle whole. The BigQuery lexer-split semantics are available via
// SplitFlat for the bytebase-switch node when it wires the splitter map; the
// parse driver itself uses the block-aware split.
//
// Lex errors are collected from a SINGLE full-input lexer pass (collectLexErrors)
// rather than per segment. Lexing is a whole-input property and the full pass is
// complete: it reports every unterminated string/bytes/identifier/comment and
// invalid byte regardless of (a) which statement it lands in, (b) whether the
// per-statement parser stopped early on the first error-recovery boundary, or
// (c) whether Split dropped the containing chunk as "empty" (an unterminated
// block comment lexes to EOF, so its segment is filtered — yet the lex error
// must still surface). Parse errors precede lex errors in the result.
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

	// Append lex errors from one authoritative full-input pass (absolute
	// offsets). Done after parse errors so a statement's syntactic complaint
	// reads before its lexical one.
	result.Errors = append(result.Errors, collectLexErrors(input)...)

	return result
}

// collectLexErrors runs the lexer over the entire input to EOF and returns every
// lex error as a ParseError with absolute byte offsets. It is the single source
// of lex diagnostics for ParseBestEffort: a full pass cannot miss an error the
// way the per-segment parse path can (early recovery-stop, empty-segment
// filtering). The lexer hides string/bytes/identifier/comment content, so this
// pass sees exactly the same token stream the per-segment lexers would, only
// complete.
func collectLexErrors(input string) []ParseError {
	_, lexErrs := Tokenize(input)
	if len(lexErrs) == 0 {
		return nil
	}
	out := make([]ParseError, len(lexErrs))
	for i, le := range lexErrs {
		out[i] = ParseError{Loc: le.Loc, Msg: le.Msg}
	}
	return out
}

// parseSingle parses one segment (one top-level statement) into a single
// ast.Node. Returns (node, errors) where node may be nil if the segment failed
// to parse and errors lists the ParseErrors from PARSING this segment.
//
// Lex errors are NOT promoted here — they are collected once over the whole
// input by ParseBestEffort (see collectLexErrors) so none is lost to an early
// recovery-stop within this segment or to Split dropping an error-only chunk.
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
	if p.cur.Type != tokEOF {
		if hintErr := p.skipStatementLevelHint(); hintErr != nil {
			// An unterminated hint consumed the rest of the segment; record the
			// diagnostic and stop (cur is at EOF, so dispatch is skipped below).
			p.errors = append(p.errors, *hintErr)
		}
	}
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

	return result, p.errors
}
