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
	"sort"
)

// Dialect selects dialect-specific syntax rules within GoogleSQL.
type Dialect int

const (
	// DialectDefault is the union GoogleSQL mode used by Parse and
	// ParseBestEffort (compatible with Cloud Spanner and shared GoogleSQL).
	DialectDefault Dialect = iota
	// DialectBigQuery applies BigQuery-specific syntax rules: enables
	// FROM-first queries, BigQuery ML/external/connection/privilege DDL,
	// untyped CREATE TABLE columns, and BigQuery reserved keywords, while
	// rejecting Spanner-only and disabled-in-BigQuery forms.
	DialectBigQuery
	// DialectSpanner applies Cloud Spanner GoogleSQL dialect rules:
	// rejects BigQuery-only procedural scripting statements, QUALIFY,
	// IS [NOT] DISTINCT FROM, parenthesized outer query aliases, and
	// BigQuery-only DDL statements/clauses while preserving all Spanner DDL
	// and shared GoogleSQL syntax.
	DialectSpanner
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
	dialect    Dialect      // dialect-specific syntax rules

	// stopAtDotStar tells parsePostfix to leave a trailing `.*` unconsumed so
	// parseSelectItem can turn an arbitrary expression_higher_prec_than_and
	// into a select_column_dot_star (`expr.*`). Cleared by nested parseExpr
	// calls so `.*` is never accepted inside parenthesized/argument expressions.
	stopAtDotStar bool

	// strictTrailing mirrors the parseSingle mode: the strict entry point
	// rejects a segment whose statement parsed but left tokens behind,
	// dropping the truncated node; best-effort keeps the parsed prefix.
	strictTrailing bool

	// inArrayColumnSchema is set while parsing a table column's type (a
	// column_schema_inner position) so parseType admits the Spanner ARRAY
	// vector-length parameter `ARRAY<FLOAT32>(vector_length => N)`. It propagates
	// into nested ARRAY elements (ARRAY<ARRAY<…>(vector_length=>…)> parses) but is
	// cleared when descending into a STRUCT/RANGE/MAP/FUNCTION template, where the
	// parameter is not a column schema (oracle: STRUCT<v ARRAY<…>(vector_length=>…)>
	// rejects). It is never set for a CAST target or any general `type` position.
	inArrayColumnSchema bool

	// parsedSubqueries holds, by the offset of its body, every expression
	// subquery parseSubqueryBodyRaw parsed in place. fillSubqueries attaches
	// these instead of reparsing the body text, which reparsed every nested
	// body once per level, so n nested subqueries cost O(n²).
	parsedSubqueries map[int]ast.Node

	// bracketGroupEnd records where the raw scan of each bracket group ended
	// (its closing bracket, or EOF), by the offset of the opening bracket. A
	// body whose in-place parse fails is scanned raw, and so is every
	// enclosing body whose parse fails after it; jumping over the groups
	// already scanned keeps that linear.
	bracketGroupEnd map[int]parserCheckpoint

	depth int // nesting depth of the productions being parsed (see enterNesting)

	// textSubqueryDepth holds, by the subquery node's Loc, the nesting depth
	// at which each subquery body captured as text (not parsed in place)
	// sits. reparseSubquery starts its nested parser at that depth, so the
	// nesting limit counts the levels around the body too.
	textSubqueryDepth map[ast.Loc]int
}

func (p *Parser) isBigQuery() bool {
	return p.dialect == DialectBigQuery
}

func (p *Parser) isSpanner() bool {
	return p.dialect == DialectSpanner
}

// parserCheckpoint is the parser and lexer position a speculative parse
// rewinds to. The lexer's own errors are not kept: lex diagnostics come from
// the separate whole-input pass in parseAll.
type parserCheckpoint struct {
	cur, prev, nextBuf Token
	hasNext            bool
	lexPos, lexStart   int
	errLen             int
}

// checkpoint captures the current position.
func (p *Parser) checkpoint() parserCheckpoint {
	return parserCheckpoint{
		cur: p.cur, prev: p.prev, nextBuf: p.nextBuf, hasNext: p.hasNext,
		lexPos: p.lexer.pos, lexStart: p.lexer.start,
		errLen: len(p.errors),
	}
}

// restore rewinds to c, dropping the errors recorded since.
func (p *Parser) restore(c parserCheckpoint) {
	p.cur, p.prev, p.nextBuf, p.hasNext = c.cur, c.prev, c.nextBuf, c.hasNext
	p.lexer.pos, p.lexer.start = c.lexPos, c.lexStart
	p.errors = p.errors[:c.errLen]
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
		Position: p.cur.Loc.Start, End: p.cur.Loc.End,
		Message: msg,
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
		Position: p.cur.Loc.Start, End: p.cur.Loc.End,
		Message: name + " statement parsing is not yet supported",
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
			Position: p.cur.Loc.Start, End: p.cur.Loc.End,
			Message: "syntax error at end of input",
		}
	}
	text := p.cur.Str
	if text == "" {
		text = TokenName(p.cur.Type)
	}
	return &ParseError{
		Position: p.cur.Loc.Start, End: p.cur.Loc.End,
		Message: "unknown or unsupported statement starting with " + text,
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
		return p.parseHintBody(hintStart, "empty statement hint", "unterminated statement hint")
	case next.Type == int('['):
		// @[ int @] { ... }
		p.advance() // consume '@'
		p.advance() // consume '['
		// Skip up to the matching @] then the brace body.
		for p.cur.Type != tokEOF && p.cur.Type != int(']') {
			p.advance()
		}
		if p.cur.Type != int(']') {
			return &ParseError{Position: hintStart.Start, End: hintStart.End, Message: "unterminated statement hint"}
		}
		p.advance() // consume ']'
		if p.cur.Type == int('{') {
			return p.parseHintBody(hintStart, "empty statement hint", "unterminated statement hint")
		}
		// `@[ … @]` with no following `{` body: a hint preamble must carry a
		// brace body, so this is malformed.
		return &ParseError{Position: hintStart.Start, End: hintStart.End, Message: "unterminated statement hint"}
	case next.Type == tokInteger:
		p.advance() // consume '@'
		p.advance() // consume the int
		// Combined integer + brace hint: `@1@{ a = NULL }` (hint_with_body_prefix).
		if p.cur.Type == int('@') && p.peekNext().Type == int('{') {
			p.advance() // consume '@'
			return p.parseHintBody(hintStart, "empty statement hint", "unterminated statement hint")
		}
	}
	return nil
}

// parseHintBody parses `{ hint_entry (, hint_entry)* }`, where
// `hint_entry: identifier_in_hints [. identifier_in_hints] = expression`
// (googlesql.tm hint_entry). It walks the body once: a body that ends before
// its closing '}' is unterminated, and one whose entry is malformed reports
// the syntax error at the token that ended it.
func (p *Parser) parseHintBody(hintStart ast.Loc, emptyMsg, untermMsg string) *ParseError {
	if p.cur.Type != int('{') {
		return nil
	}
	p.advance() // consume '{'
	if p.cur.Type == int('}') {
		end := p.advance().Loc.End
		return &ParseError{Position: hintStart.Start, End: end, Message: emptyMsg}
	}
	unterminated := func() *ParseError {
		return &ParseError{Position: hintStart.Start, End: hintStart.End, Message: untermMsg}
	}
	// fail reports the syntax error at cur. A body that runs out of input is
	// unterminated instead; that is decided where cur reaches EOF, so the body is
	// walked once however deeply hints nest.
	fail := func() *ParseError {
		if p.cur.Type == tokEOF {
			return unterminated()
		}
		return p.syntaxErrorAtCur()
	}
	for {
		if p.cur.Type == tokEOF {
			return unterminated()
		}
		if !isIdentifierInHints(p.cur.Type) {
			return fail()
		}
		p.advance() // consume identifier_in_hints
		if p.cur.Type == int('.') {
			p.advance() // consume '.'
			if !isIdentifierInHints(p.cur.Type) {
				return fail()
			}
			p.advance() // consume second identifier_in_hints
		}
		if p.cur.Type != int('=') {
			return fail()
		}
		p.advance() // consume '='
		if p.cur.Type == tokEOF {
			return unterminated()
		}
		if _, err := p.parseExpr(); err != nil {
			return fail()
		}
		switch p.cur.Type {
		case int(','):
			p.advance() // consume ','
		case int('}'):
			p.advance() // consume '}'
			return nil
		case tokEOF:
			return unterminated()
		default:
			return fail()
		}
	}
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
	// A labeled procedural statement (`label: BEGIN|WHILE|LOOP|REPEAT|FOR …`) leads
	// with an identifier followed by ':'. It is the only statement form that starts
	// with a bare identifier, so check it before the keyword dispatch (owned by the
	// parser-scripting node). Reachable both at the top level and inside a
	// statement_list.
	if p.isScriptLabelStart() {
		return p.parseLabeledStmt()
	}
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
	case kwTABLE:
		// `TABLE t` is a query (query_primary's table_clause_reserved), a read
		// equivalent to `SELECT * FROM t`; it takes the query's trailing
		// ORDER BY / LIMIT / pipe suffix like any other query.
		return p.parseQueryStatement()
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
	case tokPipe:
		// A statement that is only a pipe suffix (subpipeline_statement). It is
		// a complete statement in the grammar, not a fragment: `|> AGGREGATE
		// COUNT(*) |> WHERE n > 1` runs against the pipeline's input. Its
		// operators' expressions may embed subqueries, so fill them like the
		// other statements.
		return p.parseStmtWithSubqueries(p.parseSubpipelineStmt)

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
		// UNDROP SCHEMA [IF NOT EXISTS] <path> [FOR SYSTEM_TIME AS OF expr]
		// [OPTIONS(…)] (undrop_statement). The grammar's action allows only
		// SCHEMA; every other kind is a syntax error, as in ZetaSQL. The
		// AS OF expression and the OPTIONS values may embed subqueries, so wrap
		// for fillSubqueries.
		return p.parseStmtWithSubqueries(p.parseUndropStmt)
	case kwTRUNCATE:
		// TRUNCATE TABLE (a BigQuery DML statement; the .g4 groups it under DDL).
		if p.isSpanner() {
			return nil, &ParseError{
				Position: p.cur.Loc.Start, End: p.cur.Loc.End,
				Message: "syntax error: TRUNCATE is not supported in Spanner",
			}
		}
		return p.parseStmtWithSubqueries(p.parseTruncateStmt)
	case kwDEFINE:
		// DEFINE TABLE (DEFINE MACRO remains unsupported).
		return p.parseDefineStmt()

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
		return p.parseBeginEndBlock()
	case kwSTART:
		// START TRANSACTION (begin_statement) | START BATCH (start_batch_statement).
		return p.parseStartStmt()
	case kwCOMMIT:
		return p.parseCommitStmt()
	case kwROLLBACK:
		return p.parseRollbackStmt()
	case kwSET:
		// SET TRANSACTION / SET variable / SET @p / SET @@v / SET (tuple)
		// (set_statement) — owned by the parser-scripting node. The RHS / tuple
		// values are full expressions that may embed subqueries, filled by
		// parseSetStmt itself.
		return p.parseSetStmt()
	case kwRUN:
		// RUN BATCH is a sql_statement_body_no_pipe_suffix form: the
		// pipe-capable rule holds the distinct run_statement, which this
		// parser does not read, so a suffix after RUN BATCH is not accepted.
		return p.parseRunBatchStmt()
	case kwABORT:
		return p.parseAbortBatchStmt()

	// --- Utility / metadata ---
	case kwEXPLAIN:
		return p.parseExplainStmt()
	case kwDESCRIBE:
		// DESCRIBE may carry a pipe suffix in shared GoogleSQL; BigQuery
		// rejects pipe operators on DESCRIBE.
		if p.isBigQuery() {
			return p.parseStmtWithSubqueries(p.parseDescribeStmt)
		}
		return p.parseMaybePipeSuffixStmt(p.parseDescribeStmt)
	case kwDESC:
		if p.isBigQuery() {
			return p.parseStmtWithSubqueries(p.parseDescribeStmt)
		}
		return p.parseMaybePipeSuffixStmt(p.parseDescribeStmt)
	case kwSHOW:
		// SHOW show_target [FROM path] [LIKE 'str'] (show_statement), one of the
		// statements that may carry a pipe suffix
		// (sql_statement_body_maybe_pipe_suffix). Owned by the parser-utility
		// node (utility.go).
		return p.parseMaybePipeSuffixStmt(p.parseShowStmt)
	case kwANALYZE:
		return p.parseAnalyzeStmt()
	case kwASSERT:
		// ASSERT carries a full expression that may embed subqueries; wrap so they
		// are re-parsed (fillSubqueries) for the query-span / lineage extractor.
		return p.parseStmtWithSubqueries(p.parseAssertStmt)
	case kwCALL:
		// CALL arguments are full expressions that may embed subqueries; wrap so
		// they are re-parsed (fillSubqueries) like the query / DML paths. CALL is
		// one of the statements that may carry a pipe suffix
		// (sql_statement_body_maybe_pipe_suffix).
		return p.parseMaybePipeSuffixStmt(p.parseCallStmt)
	case kwEXECUTE:
		// EXECUTE IMMEDIATE <sql> [INTO ...] [USING ...] (execute_immediate) — owned
		// by the parser-scripting node. The SQL / USING expressions may embed
		// subqueries, filled by parseExecuteImmediateStmt itself. It may carry a
		// pipe suffix like CALL.
		return p.parseMaybePipeSuffixStmt(p.parseExecuteImmediateStmt)
	case kwIMPORT:
		// IMPORT MODULE|PROTO <path|string> [AS|INTO alias] [OPTIONS(…)]
		// (import_statement, parser-utility). The OPTIONS values are full
		// expressions that may embed subqueries, so wrap for fillSubqueries.
		return p.parseStmtWithSubqueries(p.parseImportStmt)
	case kwMODULE:
		// MODULE <path> [OPTIONS(…)] (module_statement, parser-utility). Same
		// OPTIONS subquery filling.
		return p.parseStmtWithSubqueries(p.parseModuleStmt)
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

	// --- Procedural / scripting (parser-scripting node) ---
	//
	// These are legal inside a procedural body (statement_list) and, for the
	// BEGIN…END block / SET / EXECUTE IMMEDIATE forms, also as the top-level
	// statement of a script (oracle: the Spanner emulator parses a top-level
	// BEGIN…END / SET / EXECUTE IMMEDIATE). The control-flow forms are
	// BigQuery-only at the union level (Spanner syntax-rejects them standalone);
	// omni accepts them on the authority of the legacy .g4 + the BigQuery truth1
	// corpus so a script fragment parses rather than drawing a false syntax error.
	case kwIF:
		if p.isSpanner() {
			return nil, p.syntaxErrorAtCur()
		}
		return p.parseIfStmt()
	case kwCASE:
		if p.isSpanner() {
			return nil, p.syntaxErrorAtCur()
		}
		return p.parseCaseStmt()
	case kwWHILE:
		if p.isSpanner() {
			return nil, p.syntaxErrorAtCur()
		}
		return p.parseWhileStmt()
	case kwLOOP:
		if p.isSpanner() {
			return nil, p.syntaxErrorAtCur()
		}
		return p.parseLoopStmt()
	case kwREPEAT:
		if p.isSpanner() {
			return nil, p.syntaxErrorAtCur()
		}
		return p.parseRepeatStmt()
	case kwFOR:
		if p.isSpanner() {
			return nil, p.syntaxErrorAtCur()
		}
		return p.parseForInStmt()
	case kwDECLARE:
		if p.isSpanner() {
			return nil, p.syntaxErrorAtCur()
		}
		return p.parseDeclareStmt()
	case kwBREAK, kwLEAVE:
		if p.isSpanner() {
			return nil, p.syntaxErrorAtCur()
		}
		return p.parseBreakStmt()
	case kwCONTINUE, kwITERATE:
		if p.isSpanner() {
			return nil, p.syntaxErrorAtCur()
		}
		return p.parseContinueStmt()
	case kwRETURN:
		if p.isSpanner() {
			return nil, p.syntaxErrorAtCur()
		}
		return p.parseReturnStmt()
	case kwRAISE:
		if p.isSpanner() {
			return nil, p.syntaxErrorAtCur()
		}
		return p.parseRaiseStmt()

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
// The subquery seam: the expressions node captures `(SELECT …)`, `EXISTS(…)`,
// and `ARRAY(…)` subqueries as RawText with Query nil. After the top-level
// query parses, fillSubqueries walks the produced tree and gives each one a
// real *QueryStmt in its Query field: the parse parseSubqueryBodyRaw made of
// the body in place, or else a re-parse of RawText. This is a post-pass so
// that only statements that fill their subqueries get them, and it yields a
// complete tree for the downstream query-span extractor. A re-parse error is
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

// parseMaybePipeSuffixStmt parses one of the statements that may carry a pipe
// suffix (sql_statement_body_maybe_pipe_suffix) and wraps it when a suffix
// follows, then fills every expression-embedded subquery in the result — in
// the statement and in the pipe operators' expressions alike.
func (p *Parser) parseMaybePipeSuffixStmt(parse func() (ast.Node, error)) (ast.Node, error) {
	start := p.cur.Loc.Start
	n, err := parse()
	if err != nil {
		return nil, err
	}
	n, err = p.parseMaybePipeSuffix(n, start)
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
//
// The walk runs before parseSingle measures the tree's height (the tree is
// not complete until the bodies are attached), so it keeps its own stack
// instead of recursing, and stops past maxTreeHeight levels: parseSingle
// rejects a tree that deep.
func (p *Parser) fillSubqueries(node ast.Node) {
	if node == nil {
		return
	}
	type level struct {
		node  ast.Node
		depth int
	}
	stack := []level{{node, 1}}
	var kids childLister
	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if top.depth > maxTreeHeight {
			return
		}
		switch sq := top.node.(type) {
		case *ast.SubqueryExpr:
			if sq.Query == nil {
				sq.Query = p.fillSubquery(sq.RawText, sq.TextStart, sq.Loc)
			}
		case *ast.ExistsExpr:
			if sq.Query == nil {
				sq.Query = p.fillSubquery(sq.RawText, sq.TextStart, sq.Loc)
			}
		case *ast.ArraySubqueryExpr:
			if sq.Query == nil {
				sq.Query = p.fillSubquery(sq.RawText, sq.TextStart, sq.Loc)
			}
		}
		// Push the children last first, so they are visited in order, as a
		// recursive pre-order walk visits them.
		kids.list(top.node)
		for i := len(kids.children) - 1; i >= 0; i-- {
			stack = append(stack, level{kids.children[i], top.depth + 1})
		}
	}
}

// fillSubquery re-parses one captured body, or attaches the parse
// parseSubqueryBodyRaw already made of it in place. An empty or comment-only
// body (`EXISTS()`, `ARRAY(/* c */)`) is not a query; the grammar requires
// one, so it is recorded as a syntax error at the node instead of silently
// skipped.
func (p *Parser) fillSubquery(raw string, textStart int, loc ast.Loc) ast.Node {
	if q, ok := p.parsedSubqueries[textStart]; ok && raw != "" {
		return q
	}
	if raw == "" {
		p.errors = append(p.errors, ParseError{Position: loc.Start, End: loc.End, Message: "syntax error: expected a query inside the parentheses"})
		return nil
	}
	return p.reparseSubquery(raw, textStart, p.textSubqueryDepth[loc])
}

// reparseSubquery parses a captured subquery body (the RawText between the outer
// parens) into a *QueryStmt. It runs a nested Parser over the raw text with
// textStart, the absolute offset of RawText's first character, as the lexer
// base, so every position and error it produces is in the outer text's
// coordinates. A parse failure records the error and returns nil (the node
// keeps RawText, so round-trip still works). The returned QueryStmt's own
// embedded subqueries are filled by the fillSubqueries walk that attaches it.
//
// The nested parser starts at depth, the nesting depth the body sits at in
// the outer statement. A fresh count would let each fill descend another
// maxNestingDepth levels, and the fills chain: every body attached is walked
// for the bodies inside it.
func (p *Parser) reparseSubquery(raw string, textStart, depth int) ast.Node {
	base := textStart
	if base < 0 {
		base = 0
	}
	sub := &Parser{
		lexer:      NewLexerWithOffset(raw, base),
		input:      raw,
		baseOffset: base,
		dialect:    p.dialect,
		depth:      depth,
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
			p.errors = append(p.errors, ParseError{Position: base, End: base + len(raw), Message: err.Error()})
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
	// fillSubqueries goes on to walk q once it is attached and fills q's own
	// subqueries, using the ones the nested parser made in place. Filling them
	// here as well walked q twice and reported a failing nested body once per
	// enclosing level.
	for off, nested := range sub.parsedSubqueries {
		p.recordParsedSubquery(off, nested)
	}
	for loc, depth := range sub.textSubqueryDepth {
		if p.textSubqueryDepth == nil {
			p.textSubqueryDepth = make(map[ast.Loc]int)
		}
		p.textSubqueryDepth[loc] = depth
	}
	return q
}

// recordParsedSubquery keeps q, a subquery body parsed in place at absolute
// offset off, for fillSubquery.
func (p *Parser) recordParsedSubquery(off int, q ast.Node) {
	if p.parsedSubqueries == nil {
		p.parsedSubqueries = make(map[int]ast.Node)
	}
	p.parsedSubqueries[off] = q
}

// ParseResult holds the outcome of a best-effort parse. File contains every
// statement that parsed successfully (empty while statement bodies are
// stubbed); Errors contains every ParseError encountered, including LexErrors
// promoted from the underlying Lexer.
type ParseResult struct {
	File   *ast.File
	Errors []ParseError
}

// Parse is the strict entry point. It returns the parsed File and, when any
// segment fails, a ParseErrors listing every error; the File holds only the
// statements that parsed completely, so a caller that reports errors can
// still show what parsed. A
// segment whose statement parsed but left tokens behind is a syntax error at
// the first leftover token and contributes no node.
func Parse(input string) (*ast.File, error) {
	return ParseWithDialect(input, DialectDefault)
}

// ParseWithDialect is the strict entry point configured with a specific
// GoogleSQL dialect (e.g. DialectBigQuery).
func ParseWithDialect(input string, dialect Dialect) (*ast.File, error) {
	result := parseAll(input, true, dialect)
	if len(result.Errors) == 0 {
		return result.File, nil
	}
	return result.File, ParseErrors(result.Errors)
}

// ParseBestEffort is the tolerant entry point for partial or in-progress
// input. Unlike Parse, it keeps a statement whose prefix parsed even when
// tokens follow it, and reports no error for them; every other error is still
// collected. Nothing in the repository consumes this today: diagnostics, query
// type classification, and query span extraction all read Parse. Do not add
// the strict trailing-token check here.
func ParseBestEffort(input string) *ParseResult {
	return ParseBestEffortWithDialect(input, DialectDefault)
}

// ParseBestEffortWithDialect is the tolerant entry point configured with a
// specific GoogleSQL dialect.
func ParseBestEffortWithDialect(input string, dialect Dialect) *ParseResult {
	return parseAll(input, false, dialect)
}

// parseAll is the shared implementation behind Parse and ParseBestEffort:
// Split the input, parse each segment, collect nodes and errors.
// strictTrailing selects whether parseSingle rejects unconsumed trailing
// tokens.
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
// must still surface). The result is in source order.
func parseAll(input string, strictTrailing bool, dialect Dialect) *ParseResult {
	file := &ast.File{Loc: ast.Loc{Start: 0, End: len(input)}}
	result := &ParseResult{File: file}

	lexErrs := collectLexErrors(input)
	for _, seg := range Split(input) {
		node, errs := parseSingleWithDialect(seg.Text, seg.ByteStart, strictTrailing, dialect)
		if node != nil && strictTrailing && hasLexErrorIn(lexErrs, seg.ByteStart, seg.ByteEnd) {
			// The segment lexed with an error (an unterminated comment after
			// a complete statement, say), which the whole-input pass reports
			// below. Strict mode returns no node for a segment that carries
			// any error; best-effort keeps the prefix.
			node = nil
		}
		if node != nil {
			file.Stmts = append(file.Stmts, node)
		}
		result.Errors = append(result.Errors, errs...)
	}

	// Merge the lex errors from the one authoritative full-input pass
	// (absolute offsets) into source order. The sort is stable, so at one
	// position a statement's syntactic complaint still reads before its
	// lexical one.
	result.Errors = append(result.Errors, lexErrs...)
	sort.SliceStable(result.Errors, func(i, j int) bool {
		return result.Errors[i].Position < result.Errors[j].Position
	})

	return result
}

// hasLexErrorIn reports whether any lex error starts inside [start, end).
func hasLexErrorIn(errs []ParseError, start, end int) bool {
	for _, e := range errs {
		if e.Position >= start && e.Position < end {
			return true
		}
	}
	return false
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
		out[i] = ParseError{Position: le.Loc.Start, End: le.Loc.End, Message: le.Msg}
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
// strictTrailing selects whether a statement that parsed without consuming the
// whole segment is rejected (node dropped) or kept as the parsed prefix.
func parseSingle(segText string, baseOffset int, strictTrailing bool) (ast.Node, []ParseError) {
	return parseSingleWithDialect(segText, baseOffset, strictTrailing, DialectDefault)
}

func parseSingleWithDialect(segText string, baseOffset int, strictTrailing bool, dialect Dialect) (ast.Node, []ParseError) {
	p := &Parser{
		lexer:          NewLexerWithOffset(segText, baseOffset),
		input:          segText,
		baseOffset:     baseOffset,
		dialect:        dialect,
		strictTrailing: strictTrailing,
	}
	p.advance() // prime cur with the first token
	first := p.cur

	var result ast.Node
	hadHint := false
	if p.cur.Type != tokEOF {
		beforeHint := p.cur.Loc
		if hintErr := p.skipStatementLevelHint(); hintErr != nil {
			// An unterminated hint consumed the rest of the segment; record the
			// diagnostic and stop (cur is at EOF, so dispatch is skipped below).
			p.errors = append(p.errors, *hintErr)
		} else if p.cur.Loc != beforeHint {
			hadHint = true
		}
	}
	if p.cur.Type != tokEOF {
		node, err := p.parseStmt()
		if err != nil {
			if pe, ok := err.(*ParseError); ok {
				p.errors = append(p.errors, *pe)
			} else {
				p.errors = append(p.errors, ParseError{
					Position: p.cur.Loc.Start, End: p.cur.Loc.End,
					Message: err.Error(),
				})
			}
		} else if treeTooHigh(node, len(segText)) {
			// A left-deep chain parsed by a loop is shallow to parse but deep
			// as a tree, and every consumer that walks the tree recursively
			// would overflow the stack on it.
			p.errors = append(p.errors, *nestedTooDeeply(first.Loc))
			node = nil
		} else if strictTrailing && p.cur.Type != tokEOF {
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
		} else if strictTrailing && len(p.errors) > 0 {
			// The segment recorded an error somewhere: a malformed statement
			// hint before parseStmt (`@[5@] SELECT 1`) or an embedded query
			// that does not parse (`SELECT (SELECT 1 FROM t a b)`) while
			// parseStmt still returned a node. It did not parse completely,
			// so strict mode drops the node; ParseBestEffort keeps the
			// partial tree.
			node = nil
		}
		result = node
	} else if hadHint && len(p.errors) == 0 {
		p.errors = append(p.errors, *p.syntaxErrorAtCur())
	}

	return result, p.errors
}
