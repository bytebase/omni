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
// It returns a *ParseError when the hint body is malformed in a way that would
// otherwise be silently swallowed to EOF, leaving the (invalid) statement with
// no diagnostic at all — which bytebase's Diagnose must not allow:
//   - UNTERMINATED — an `@{` (or `@[…@]{`) whose closing '}' never arrives
//     before EOF.
//   - EMPTY — a balanced but entry-less `@{}` (or `@[…@]{}`). GoogleSQL requires
//     at least one hint entry (oracle: the Spanner emulator rejects `@{}` with
//     `Syntax error: Unexpected "}"`).
//
// Returns nil when there is no hint or the hint body is non-empty and closes
// cleanly. (The body's internal `key=value` entry shape is NOT validated here —
// that is the later hint-parsing node's job; the foundation only enforces
// termination and non-emptiness so dispatch can safely reach the statement.)
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
// at the current token (which must be '{'). It returns a *ParseError anchored at
// hintStart when:
//   - EOF is reached before the matching '}' (an unterminated hint), or
//   - the body is EMPTY — the closing '}' immediately follows the opening '{'
//     with nothing but whitespace/comments between (the lexer drops both). A
//     GoogleSQL hint requires at least one entry, so `@{}` is a syntax error
//     (oracle: the Spanner emulator rejects `@{} SELECT 1` with
//     `Syntax error: Unexpected "}"`). Without this guard the empty hint is
//     skipped, dispatch reaches EOF, and the invalid input draws no diagnostic.
//
// It returns nil on a clean close of a non-empty body, or when cur is not '{'.
// Validating the body's internal `key=value` entry shape is deferred to the
// later hint-parsing node; this only enforces non-emptiness and termination.
func (p *Parser) skipBalancedBraces(hintStart ast.Loc) *ParseError {
	if p.cur.Type != int('{') {
		return nil
	}
	// Empty body: the closing '}' immediately follows the opening '{' (the lexer
	// has already dropped any whitespace/comments between them). Detect it by the
	// token right after '{' rather than by counting interior tokens, so a body
	// that merely STARTS with a brace (e.g. `@{{…}}`) is treated as non-empty —
	// its malformed `key=value` shape is the later hint node's concern, not an
	// emptiness error.
	if p.peekNext().Type == int('}') {
		return &ParseError{Loc: hintStart, Msg: "empty statement hint"}
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
		return p.parseQueryStatement()
	case kwWITH:
		// WITH cte... query  (also WITH RECURSIVE). The query-leading WITH is an
		// aliased-query CTE list, distinct from the inline WITH(...) expression
		// (owned by the expressions node in expression position).
		return p.parseQueryStatement()
	case kwGRAPH:
		// GQL graph query: GRAPH <name> <gql ops> (gql_statement). Owned by the
		// parser-gql node (graph_query.go).
		return p.parseGQLStmt()
	case int('('):
		// Parenthesized query at top level, e.g. (SELECT 1) UNION ALL (SELECT 2).
		return p.parseQueryStatement()
	case kwFROM:
		// FROM-first query head. The legacy grammar recognizes a leading FROM as
		// a query (query_without_pipe_operators' `with_clause? from_clause` alt)
		// and then emits "Syntax error: Unexpected FROM" — i.e. FROM is a query-
		// statement head that the query rule itself rejects (oracle-confirmed:
		// Spanner rejects `FROM t` as a syntax error). parseQuery routes a leading
		// FROM to parseQueryPrimary, which produces exactly that rejection.
		return p.parseQueryStatement()

	// --- DDL ---
	case kwCREATE:
		// CREATE TABLE/VIEW/INDEX/SCHEMA/DATABASE are owned by the parser-ddl node
		// (create_table.go / view.go / create_index.go / database_schema.go); other
		// object kinds (FUNCTION/PROCEDURE/MODEL/MATERIALIZED VIEW/…) route to the
		// unsupported stub from parseCreateStmt.
		return p.parseCreateStmt()
	case kwALTER:
		return p.parseAlterStmt()
	case kwDROP:
		return p.parseDropStmt()
	case kwRENAME:
		return p.parseRenameStmt()
	case kwUNDROP:
		return p.unsupported("UNDROP")
	case kwTRUNCATE:
		// TRUNCATE TABLE (a BigQuery DML statement; the .g4 groups it under DDL).
		return p.parseStmtWithSubqueries(p.parseTruncateStmt)
	case kwDEFINE:
		// DEFINE TABLE (DEFINE MACRO is intentionally unimplemented — a
		// reserved error alt in the legacy grammar).
		return p.unsupported("DEFINE")

	// --- DML ---
	case kwINSERT:
		return p.parseStmtWithSubqueries(p.parseInsertStmt)
	case kwUPDATE:
		return p.parseStmtWithSubqueries(p.parseUpdateStmt)
	case kwDELETE:
		return p.parseStmtWithSubqueries(p.parseDeleteStmt)
	case kwMERGE:
		return p.parseStmtWithSubqueries(p.parseMergeStmt)

	// --- DCL ---
	case kwGRANT:
		return p.parseGrantStmt()
	case kwREVOKE:
		return p.parseRevokeStmt()

	// --- Transactions / batch / session ---
	case kwBEGIN:
		// A top-level BEGIN is ambiguous: `BEGIN [TRANSACTION] [modes]` is a TCL
		// transaction (begin_statement); `BEGIN … END` is a procedural block (owned
		// by parser-scripting). Disambiguate with the SAME predicate the splitter
		// uses (isTCLBeginFollower): a TCL follower (';' / EOF / TRANSACTION / READ /
		// ISOLATION) means a transaction; anything else means a BEGIN…END block.
		if isTCLBeginFollower(p.peekNext()) {
			return p.parseBeginStmt()
		}
		return p.unsupported("BEGIN...END block")
	case kwSTART:
		// START TRANSACTION (begin_statement) | START BATCH (start_batch_statement).
		return p.parseStartStmt()
	case kwCOMMIT:
		return p.parseCommitStmt()
	case kwROLLBACK:
		return p.parseRollbackStmt()
	case kwSET:
		// SET TRANSACTION / SET variable / SET system-var (set_statement) — owned by
		// a separate node (not parser-utility).
		return p.unsupported("SET")
	case kwRUN:
		return p.parseRunBatchStmt()
	case kwABORT:
		return p.parseAbortBatchStmt()

	// --- Utility / metadata ---
	case kwEXPLAIN:
		return p.unsupported("EXPLAIN")
	case kwDESCRIBE:
		return p.parseDescribeStmt()
	case kwDESC:
		return p.parseDescribeStmt()
	case kwSHOW:
		return p.unsupported("SHOW")
	case kwANALYZE:
		return p.parseAnalyzeStmt()
	case kwASSERT:
		// ASSERT carries a full expression that may embed subqueries; wrap so they
		// are re-parsed (fillSubqueries) for the query-span / lineage extractor.
		return p.parseStmtWithSubqueries(p.parseAssertStmt)
	case kwCALL:
		// CALL arguments are full expressions that may embed subqueries; wrap so
		// they are re-parsed (fillSubqueries) like the query / DML paths.
		return p.parseStmtWithSubqueries(p.parseCallStmt)
	case kwEXECUTE:
		// EXECUTE IMMEDIATE.
		return p.unsupported("EXECUTE")
	case kwIMPORT:
		return p.unsupported("IMPORT")
	case kwMODULE:
		return p.unsupported("MODULE")
	case kwEXPORT:
		// EXPORT DATA (AS query) | EXPORT MODEL | EXPORT … METADATA (parser-utility,
		// still a stub). DATA/MODEL carry expressions (the AS query; OPTIONS values)
		// that may embed subqueries, so parseExportStmt wraps them for fillSubqueries.
		return p.parseExportStmt()
	case kwLOAD:
		// LOAD DATA (INTO|OVERWRITE) … FROM FILES (…). The OPTIONS / FROM FILES /
		// PARTITIONS / PARTITION BY / CLUSTER BY clauses parse full expressions that
		// can embed subqueries, so wrap for fillSubqueries like the DML family.
		return p.parseStmtWithSubqueries(p.parseLoadData)
	case kwCLONE:
		// CLONE DATA INTO … FROM source (UNION ALL source)*. Each source can carry a
		// FOR SYSTEM_TIME AS OF expr and a WHERE expr that may embed subqueries, so
		// wrap for fillSubqueries.
		return p.parseStmtWithSubqueries(p.parseCloneData)

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

// parseQueryStatement is the dispatch entry for a query_statement (SELECT /
// WITH / parenthesized-query / FROM-first head). It parses the full query via
// parseQuery (select.go + companions), then fills in the inner Query of every
// expression-embedded subquery node (SubqueryExpr / ExistsExpr /
// ArraySubqueryExpr) that the expressions node left as RawText-only.
//
// The subquery seam: the expressions node (a dependency, frozen) captures
// `(SELECT …)`, `EXISTS(…)`, and `ARRAY(…)` subqueries as RawText with Query nil
// (it does not own the query grammar). parser-select owns the query grammar, so
// after the top-level query parses, fillSubqueries walks the produced tree and
// re-parses each subquery's RawText into a real *QueryStmt, attaching it to the
// node's Query field. This is done as a post-pass (rather than inline) because
// the expressions parser already consumed those tokens; re-parsing RawText is
// the only seam available without editing the frozen expressions node, and it
// yields a complete tree for the downstream query-span extractor. Subquery
// RawText is a balanced parenthesized query (the expressions node scanned it
// with bracket balancing), so it always re-parses cleanly; a re-parse error is
// surfaced as a diagnostic but does not discard the outer statement.
func (p *Parser) parseQueryStatement() (ast.Node, error) {
	q, err := p.parseQuery()
	if err != nil {
		return nil, err
	}
	p.fillSubqueries(q)
	return q, nil
}

// parseStmtWithSubqueries runs a statement parse function and then fills every
// expression-embedded subquery the way parseQueryStatement does. It wraps every
// non-query statement that can carry full expressions: the DML family (VALUES
// rows, SET / merge-action values, WHERE / ON / merge conditions,
// ASSERT_ROWS_MODIFIED, THEN RETURN items) AND the utility statements that embed
// expressions — ASSERT (`ASSERT expression`) and CALL (`CALL p(expr, …)`). Those
// expressions may contain SubqueryExpr / ExistsExpr / ArraySubqueryExpr captured
// as RawText with Query==nil by the frozen expressions node. fillSubqueries
// re-parses each into a real *QueryStmt — which both (a) completes the tree for
// the downstream query-span / lineage extractor (which walks the AST and needs
// the inner *QueryStmt to resolve a subquery's tables/columns) and (b) surfaces
// a malformed embedded subquery (e.g. `UPDATE t SET x = (SELECT 1 FROM s a b)` /
// `ASSERT (SELECT 1 FROM s a b) = 1`, which the oracle rejects on the stray `b`)
// as a diagnostic, matching the query-statement path.
func (p *Parser) parseStmtWithSubqueries(parse func() (ast.Node, error)) (ast.Node, error) {
	n, err := parse()
	if err != nil {
		return nil, err
	}
	p.fillSubqueries(n)
	return n, nil
}

// fillSubqueries walks node and re-parses the RawText of every SubqueryExpr /
// ExistsExpr / ArraySubqueryExpr whose Query is still nil, attaching the parsed
// *QueryStmt. Re-parse errors are appended to the parser's error list (so
// Diagnose still surfaces a malformed embedded subquery) but never abort the
// outer parse. The walk descends into the filled subqueries too, so arbitrarily
// nested subqueries are completed.
func (p *Parser) fillSubqueries(node ast.Node) {
	ast.Inspect(node, func(n ast.Node) bool {
		switch sq := n.(type) {
		case *ast.SubqueryExpr:
			if sq.Query == nil && sq.RawText != "" {
				sq.Query = p.reparseSubquery(sq.RawText, sq.Loc)
			}
		case *ast.ExistsExpr:
			if sq.Query == nil && sq.RawText != "" {
				sq.Query = p.reparseSubquery(sq.RawText, sq.Loc)
			}
		case *ast.ArraySubqueryExpr:
			if sq.Query == nil && sq.RawText != "" {
				sq.Query = p.reparseSubquery(sq.RawText, sq.Loc)
			}
		}
		return true
	})
}

// reparseSubquery parses a captured subquery body (the RawText between the outer
// parens) into a *QueryStmt. It runs a nested Parser over the raw text; lexing
// uses a best-effort base offset (the node's start) so positions stay close to
// the original source. A parse failure records the error and returns nil (the
// node keeps RawText, so round-trip still works). The returned QueryStmt's own
// embedded subqueries are filled recursively.
func (p *Parser) reparseSubquery(raw string, loc ast.Loc) ast.Node {
	base := 0
	if loc.Start >= 0 {
		// Anchor near the original site. The exact column of RawText within the
		// parens is approximate (leading '(' + whitespace was trimmed), but this
		// keeps offsets monotonic and close, which is all the query-span / Diagnose
		// consumers need for an embedded subquery.
		base = loc.Start
	}
	sub := &Parser{
		lexer:      NewLexerWithOffset(raw, base),
		input:      raw,
		baseOffset: base,
	}
	sub.advance()
	if sub.cur.Type == tokEOF {
		return nil
	}
	q, err := sub.parseQuery()
	if err != nil {
		if pe, ok := err.(*ParseError); ok {
			p.errors = append(p.errors, *pe)
		} else {
			p.errors = append(p.errors, ParseError{Loc: loc, Msg: err.Error()})
		}
		return nil
	}
	// The subquery body must consume its entire RawText: a subquery is a complete
	// `query` and any token left over is trailing junk that makes the embedding
	// statement invalid (oracle: `SELECT (SELECT 1 FROM t a b)` rejects on the
	// stray `b`). The expressions node captured the balanced parens but did not
	// validate the interior; we surface that reject here so the outer statement is
	// diagnosed rather than silently accepting a partial subquery parse.
	if sub.cur.Type != tokEOF {
		p.errors = append(p.errors, *sub.syntaxErrorAtCur())
		return nil
	}
	// Recurse into the freshly-parsed subquery to fill its own nested subqueries.
	p.fillSubqueries(q)
	return q
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
		} else if p.cur.Type != tokEOF {
			// The grammar root is `stmts EOF` (antlr_rules.md §1): a complete
			// statement must be followed by end-of-input. parseStmt succeeded but
			// left tokens behind — e.g. `GRANT a ON foo TO 'x' garbage` — so the
			// segment is NOT a single well-formed statement. Without this check the
			// trailing junk is silently dropped (stmts=1, errs=0), a false-accept
			// that yields zero diagnostics for the Diagnose consumer. Report a
			// syntax error at the first unconsumed token and drop the node: the
			// segment as a whole does not parse, matching every other reject path
			// (which all return a nil node). We assert EOF only on the success
			// path — a parse that already errored left cur mid-statement, so
			// asserting EOF there would emit a spurious second diagnostic.
			p.errors = append(p.errors, *p.syntaxErrorAtCur())
			node = nil
		}
		result = node
	}

	return result, p.errors
}
