# Snowflake Parser Entry (F4) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the F4 parser-entry framework for `snowflake/parser` — public `Parse` / `ParseBestEffort` API, `Parser` struct with helpers, 37-case dispatch switch with stubs, `ParseError` type, `LineTable` helper for byte-offset → (line, col) conversion, and a small F2 lexer enhancement (`NewLexerWithOffset`) that unlocks mysql-style split-then-parse.

**Architecture:** Mysql-style split-then-parse. `ParseBestEffort` calls F3's `Split` to segment the input, then `parseSingle(segText, baseOffset)` constructs a `Parser{lexer: NewLexerWithOffset(segText, baseOffset)}` per segment and drives a recursive-descent loop. Each segment's `Loc` values are absolute in the original input thanks to the lexer's offset support. Statement parsing is stubbed — 37 dispatch cases each call `p.unsupported("X")` to emit a "not yet supported" `ParseError`. Tier 1+ DAG nodes replace specific stubs with real implementations.

**Tech Stack:** Go 1.25, stdlib only (`sort` for LineTable binary search). F4 depends only on F1's `snowflake/ast` and F2/F3's `snowflake/parser` (intra-package).

**Spec:** `docs/superpowers/specs/2026-04-09-snowflake-parser-entry-design.md` (commit `62f15f2`)

**Worktree:** `/Users/h3n4l/OpenSource/omni/.worktrees/parser-entry` on branch `feat/snowflake/parser-entry`

**Module path:** `github.com/bytebase/omni`. F4 lives in `github.com/bytebase/omni/snowflake/parser` (same package as F2 and F3).

**Commit policy:** No commits during implementation. The user reviews the full diff at the end; only after explicit approval do we commit.

---

## File Structure

### Created

| File | Responsibility | Approx LOC |
|------|----------------|-----------|
| `snowflake/parser/parser.go` | Parser struct, Parse, ParseBestEffort, parseSingle, dispatch switch, helpers, ParseResult | 360 |
| `snowflake/parser/parser_test.go` | 11 test categories for Parse / ParseBestEffort | 420 |
| `snowflake/parser/linetable.go` | LineTable struct, NewLineTable, Position | 90 |
| `snowflake/parser/linetable_test.go` | 7 LineTable round-trip tests | 140 |

### Modified

| File | Changes | Approx LOC delta |
|------|---------|-----|
| `snowflake/parser/lexer.go` | Add `baseOffset` field, `NewLexerWithOffset` constructor, rename `NextToken` → `nextTokenInner` (unexported) + thin `NextToken` wrapper that shifts, update `Errors()` to return shifted copies | +35 |
| `snowflake/parser/errors.go` | Add `ParseError` type + `Error() string` method alongside existing `LexError` | +20 |
| `snowflake/parser/lexer_test.go` | Add `TestLexer_NewLexerWithOffset`, `TestLexer_NewLexerWithOffset_LexErrors`, `TestLexer_NewLexer_UnchangedBehavior` | +60 |

**Total: ~1,125 LOC** across 4 new + 3 modified files.

---

## Task 1: F2 enhancement — add `NewLexerWithOffset` to lexer.go

This task modifies the F2 lexer to support a `baseOffset` that shifts all token and error `Loc` values. The approach: add an unexported `nextTokenInner` method containing the existing `NextToken` body, then make the public `NextToken` a thin wrapper that shifts `Loc` values by `baseOffset` before returning. Similarly, the public `Errors()` method returns offset-shifted copies. Scan helpers inside the lexer continue to use local positions — no changes to them.

**Files:**
- Modify: `snowflake/parser/lexer.go`
- Modify: `snowflake/parser/lexer_test.go`

- [ ] **Step 1: Confirm worktree state**

Run: `pwd && git rev-parse --abbrev-ref HEAD`
Expected:
```
/Users/h3n4l/OpenSource/omni/.worktrees/parser-entry
feat/snowflake/parser-entry
```

If either is wrong, stop.

- [ ] **Step 2: Modify the Lexer struct to add baseOffset**

In `snowflake/parser/lexer.go`, replace the current Lexer struct definition with:

```go
// Lexer is a Snowflake SQL tokenizer. Construct via NewLexer (or
// NewLexerWithOffset when tokenizing a substring of a larger document)
// and call NextToken until it returns Token{Type: tokEOF}. Lex errors
// are collected in a slice (Errors); each error is accompanied by a
// synthetic tokInvalid token at the failure site so consumers can
// choose to halt on the first error or proceed with best-effort parsing.
type Lexer struct {
	input      string
	pos        int // current byte offset (one past the last consumed byte)
	start      int // start byte of the token currently being scanned
	errors     []LexError
	baseOffset int // added to token Loc.Start/End and error Loc.Start/End when returned via public API
}
```

- [ ] **Step 3: Add the NewLexerWithOffset constructor**

In `snowflake/parser/lexer.go`, replace the current `NewLexer` function with:

```go
// NewLexer constructs a Lexer for the given input. Token positions are
// zero-based byte offsets into input.
func NewLexer(input string) *Lexer {
	return &Lexer{input: input}
}

// NewLexerWithOffset constructs a Lexer whose emitted token Loc values
// are shifted by baseOffset. Use this when tokenizing a substring of a
// larger document and you want Loc values to refer to positions in the
// larger document rather than the substring. The F4 parser-entry uses
// this to plumb absolute positions through F3's split-then-parse
// pipeline.
//
// The shift is applied lazily when tokens and errors leave the Lexer
// via NextToken() and Errors(). Internal scan helpers continue to use
// local (unshifted) positions against input.
func NewLexerWithOffset(input string, baseOffset int) *Lexer {
	return &Lexer{input: input, baseOffset: baseOffset}
}
```

- [ ] **Step 4: Rename NextToken body to nextTokenInner and add a shifting wrapper**

Find the existing `NextToken` function body in `snowflake/parser/lexer.go`. Rename it to `nextTokenInner` (lower-case, unexported). Then add a new public `NextToken` that calls it and applies the offset:

```go
// NextToken advances the lexer and returns the next token. At end of input
// it returns Token{Type: tokEOF}; subsequent calls continue to return EOF.
//
// If the Lexer was constructed with NewLexerWithOffset, the returned
// token's Loc values are shifted by baseOffset.
func (l *Lexer) NextToken() Token {
	tok := l.nextTokenInner()
	if l.baseOffset != 0 {
		tok.Loc.Start += l.baseOffset
		tok.Loc.End += l.baseOffset
	}
	return tok
}

// nextTokenInner returns the next token with LOCAL (unshifted) positions.
// It is the original body of the lexer dispatch — the public NextToken
// above is a thin wrapper that applies baseOffset before returning.
func (l *Lexer) nextTokenInner() Token {
	l.skipWhitespaceAndComments()
	if l.pos >= len(l.input) {
		return Token{Type: tokEOF, Loc: ast.Loc{Start: l.pos, End: l.pos}}
	}
	l.start = l.pos
	ch := l.input[l.pos]

	switch {
	case ch == '\'':
		return l.scanString()
	case ch == '"':
		return l.scanQuotedIdent()
	case ch == '$':
		return l.scanDollar()
	case ch >= '0' && ch <= '9':
		return l.scanNumber()
	case ch == '.' && l.peekDigit():
		return l.scanNumber()
	case (ch == 'X' || ch == 'x') && l.peek(1) == '\'':
		l.pos++ // consume X
		tok := l.scanString()
		tok.XPrefix = true
		tok.Loc.Start = l.start
		return tok
	case isIdentStart(ch):
		return l.scanIdentOrKeyword()
	}
	return l.scanOperatorOrPunct(ch)
}
```

The dispatch body is byte-identical to the current `NextToken` — the ONLY change is the function name from `NextToken` → `nextTokenInner`.

- [ ] **Step 5: Update Errors() to return shifted copies**

Replace the current `Errors()` method with:

```go
// Errors returns all lex errors collected so far. Errors are appended in
// source order; each error is accompanied by a tokInvalid token in the
// stream.
//
// If the Lexer was constructed with NewLexerWithOffset, the returned
// errors have Loc values shifted by baseOffset. The internal l.errors
// slice retains unshifted local positions.
func (l *Lexer) Errors() []LexError {
	if l.baseOffset == 0 {
		return l.errors
	}
	shifted := make([]LexError, len(l.errors))
	for i, e := range l.errors {
		shifted[i] = LexError{
			Loc: ast.Loc{
				Start: e.Loc.Start + l.baseOffset,
				End:   e.Loc.End + l.baseOffset,
			},
			Msg: e.Msg,
		}
	}
	return shifted
}
```

- [ ] **Step 6: Verify all existing F2 tests still pass**

Run: `go build ./snowflake/parser/...`
Expected: no output, exit 0.

Run: `go vet ./snowflake/parser/...`
Expected: no output, exit 0.

Run: `go test ./snowflake/parser/...`
Expected: all existing F2 + F3 tests pass.

If any existing test fails, the refactoring is incorrect. The most likely bug is a missing `NextToken` → `nextTokenInner` rename in one of the scan helpers. Inspect `git diff snowflake/parser/lexer.go` and re-check that only one `func (l *Lexer) nextTokenInner() Token` exists and one `func (l *Lexer) NextToken() Token` wrapper exists.

- [ ] **Step 7: Add offset tests to lexer_test.go**

Append the following test functions to the end of `snowflake/parser/lexer_test.go`:

```go

func TestLexer_NewLexerWithOffset(t *testing.T) {
	// Tokenize "   SELECT   " with baseOffset=100. The kwSELECT token
	// should report Loc.Start = 103, Loc.End = 109 (positions shifted
	// by 100 relative to the segment's local positions 3..9).
	tokens, errs := TokenizeWithOffset("   SELECT   ", 100)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %+v", errs)
	}
	// Find the kwSELECT token.
	var selectTok *Token
	for i, tok := range tokens {
		if tok.Type == kwSELECT {
			selectTok = &tokens[i]
			break
		}
	}
	if selectTok == nil {
		t.Fatalf("kwSELECT not found in tokens: %+v", tokens)
	}
	if selectTok.Loc.Start != 103 || selectTok.Loc.End != 109 {
		t.Errorf("kwSELECT Loc = %+v, want {103, 109}", selectTok.Loc)
	}
}

func TestLexer_NewLexerWithOffset_LexErrors(t *testing.T) {
	// An unterminated string should produce a LexError whose Loc is
	// shifted by baseOffset.
	_, errs := TokenizeWithOffset("'unterminated", 50)
	if len(errs) != 1 {
		t.Fatalf("expected 1 lex error, got %d: %+v", len(errs), errs)
	}
	// The error's Loc.Start should be 50 (the position of the opening
	// quote in the full document).
	if errs[0].Loc.Start != 50 {
		t.Errorf("lex error Loc.Start = %d, want 50", errs[0].Loc.Start)
	}
}

func TestLexer_NewLexer_UnchangedBehavior(t *testing.T) {
	// Regression check: the existing NewLexer should produce zero-based
	// positions, same as before the baseOffset refactor.
	tokens, errs := Tokenize("SELECT")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %+v", errs)
	}
	if len(tokens) < 1 || tokens[0].Type != kwSELECT {
		t.Fatalf("expected first token kwSELECT, got %+v", tokens)
	}
	if tokens[0].Loc.Start != 0 || tokens[0].Loc.End != 6 {
		t.Errorf("kwSELECT Loc = %+v, want {0, 6}", tokens[0].Loc)
	}
}

// TokenizeWithOffset is a test helper: one-shot tokenize with a base
// offset. Mirrors the existing Tokenize helper but constructs via
// NewLexerWithOffset.
func TokenizeWithOffset(input string, baseOffset int) (tokens []Token, errors []LexError) {
	l := NewLexerWithOffset(input, baseOffset)
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == tokEOF {
			break
		}
	}
	return tokens, l.Errors()
}
```

- [ ] **Step 8: Run the new tests**

Run: `go test -v ./snowflake/parser/... -run "TestLexer_NewLexer" 2>&1 | tail -20`
Expected: three new tests all PASS:
```
--- PASS: TestLexer_NewLexerWithOffset
--- PASS: TestLexer_NewLexerWithOffset_LexErrors
--- PASS: TestLexer_NewLexer_UnchangedBehavior
```

- [ ] **Step 9: Run the full parser package test suite**

Run: `go test ./snowflake/parser/...`
Expected: clean pass. All F2 + F3 + new offset tests green.

---

## Task 2: Add `ParseError` type to errors.go

**Files:**
- Modify: `snowflake/parser/errors.go`

- [ ] **Step 1: Append ParseError to errors.go**

Use Edit to append the following to the end of `snowflake/parser/errors.go`:

```go

// ParseError describes a single parse error with its source location.
//
// Loc uses the same shape as LexError so consumers can handle both
// uniformly. Line/column conversion is a caller-side concern via
// LineTable (defined in linetable.go).
type ParseError struct {
	Loc ast.Loc
	Msg string
}

// Error implements the error interface. Returns just the message — line
// and column are omitted here to keep ParseError a pure data carrier.
// Callers that want formatted "msg (line N, col M)" output should use a
// LineTable to convert Loc.Start into a (line, col) pair.
func (e *ParseError) Error() string {
	return e.Msg
}
```

- [ ] **Step 2: Verify compile**

Run: `go build ./snowflake/parser/...`
Expected: no output, exit 0.

---

## Task 3: Scaffold parser.go and parser_test.go

**Files:**
- Create: `snowflake/parser/parser.go`
- Create: `snowflake/parser/parser_test.go`

- [ ] **Step 1: Write parser.go stub**

Create `snowflake/parser/parser.go`:

```go
package parser
```

- [ ] **Step 2: Write parser_test.go stub**

Create `snowflake/parser/parser_test.go`:

```go
package parser
```

- [ ] **Step 3: Verify compile**

Run: `go build ./snowflake/parser/...`
Expected: no output, exit 0.

Run: `go test ./snowflake/parser/...`
Expected: clean pass (nothing new to test yet).

---

## Task 4: Implement Parser struct + helpers in parser.go

**Files:**
- Modify: `snowflake/parser/parser.go`

- [ ] **Step 1: Replace the stub with the Parser struct and helpers**

Overwrite `snowflake/parser/parser.go` with:

```go
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
```

- [ ] **Step 2: Verify compile**

Run: `go build ./snowflake/parser/...`
Expected: no output, exit 0.

Run: `go vet ./snowflake/parser/...`
Expected: no output, exit 0.

---

## Task 5: Implement `unsupported` and `unknownStatementError` helpers

**Files:**
- Modify: `snowflake/parser/parser.go`

- [ ] **Step 1: Append the two helpers to parser.go**

Use Edit to append the following to the end of `snowflake/parser/parser.go`:

```go

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
```

- [ ] **Step 2: Verify compile**

Run: `go build ./snowflake/parser/...`
Expected: no output, exit 0.

---

## Task 6: Implement the parseStmt dispatch switch

**Files:**
- Modify: `snowflake/parser/parser.go`

- [ ] **Step 1: Append parseStmt to parser.go**

Use Edit to append the following to the end of `snowflake/parser/parser.go`:

```go

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
		return p.unsupported("SELECT")
	case kwWITH:
		return p.unsupported("WITH")
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
```

- [ ] **Step 2: Verify compile**

Run: `go build ./snowflake/parser/...`
Expected: no output, exit 0.

Run: `go vet ./snowflake/parser/...`
Expected: no output, exit 0.

If the build fails with "undefined: kw*" for any of the 37 keywords, the keyword list has drifted from F2's tokens.go. Check `grep "^	kw<NAME>$" snowflake/parser/tokens.go` for each missing keyword and either add the missing constant to F2 (scope creep — probably don't) or remove the case from the dispatch.

---

## Task 7: Implement Parse, ParseBestEffort, parseSingle, ParseResult

**Files:**
- Modify: `snowflake/parser/parser.go`

- [ ] **Step 1: Append the public API to parser.go**

Use Edit to append the following to the end of `snowflake/parser/parser.go`:

```go

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
```

- [ ] **Step 2: Verify compile**

Run: `go build ./snowflake/parser/...`
Expected: no output, exit 0.

Run: `go vet ./snowflake/parser/...`
Expected: no output, exit 0.

Run: `go test ./snowflake/parser/...`
Expected: clean pass. No new tests yet, but existing F2/F3 tests should still be green.

- [ ] **Step 3: Smoke-test Parse by hand**

Run a quick sanity check via a temporary Go file:

```bash
cat > /tmp/parse-smoke.go <<'EOF'
package main

import (
	"fmt"

	"github.com/bytebase/omni/snowflake/parser"
)

func main() {
	cases := []string{
		"",
		"SELECT 1;",
		"SELECT 1; INSERT INTO t VALUES (1);",
		"FOO BAR;",
		"-- comment only\n",
	}
	for _, c := range cases {
		result := parser.ParseBestEffort(c)
		fmt.Printf("%q → %d stmts, %d errors\n", c, len(result.File.Stmts), len(result.Errors))
		for _, e := range result.Errors {
			fmt.Printf("  [%d,%d] %s\n", e.Loc.Start, e.Loc.End, e.Msg)
		}
	}
}
EOF
go run /tmp/parse-smoke.go
rm /tmp/parse-smoke.go
```

Expected output:
```
"" → 0 stmts, 0 errors
"SELECT 1;" → 0 stmts, 1 errors
  [0,6] SELECT statement parsing is not yet supported
"SELECT 1; INSERT INTO t VALUES (1);" → 0 stmts, 2 errors
  [0,6] SELECT statement parsing is not yet supported
  [10,16] INSERT statement parsing is not yet supported
"FOO BAR;" → 0 stmts, 1 errors
  [0,3] unknown or unsupported statement starting with FOO
"-- comment only\n" → 0 stmts, 0 errors
```

Key observations:
- Empty input and comment-only input produce zero statements and zero errors
- Valid but unsupported SELECT produces one ParseError at byte 0 with the friendly message
- Multi-statement input produces multiple errors, with the SECOND error at byte 10 (absolute position, NOT relative to the second segment)
- Unknown keyword FOO produces a "unknown or unsupported" error

If the output doesn't match, stop and diagnose. The most likely issues are:
- Wrong byte offsets → `parseSingle` is not using `NewLexerWithOffset` correctly
- Wrong error count → `parseStmt` is not dispatching correctly
- Wrong position for the SELECT/INSERT → the `unsupported` helper is pulling from `p.cur.Loc` AFTER some advance() that shouldn't have happened

---

## Task 8: Implement LineTable

**Files:**
- Create: `snowflake/parser/linetable.go`
- Create: `snowflake/parser/linetable_test.go`

- [ ] **Step 1: Write linetable.go**

Create `snowflake/parser/linetable.go`:

```go
package parser

// LineTable provides O(log n) byte-offset → (line, column) conversion
// for a single source string. Construct once per input, then call
// Position repeatedly.
//
// Handles LF, CRLF, and bare-CR line endings. Each of the following is
// treated as exactly ONE line break:
//   - "\n" (Unix / LF)
//   - "\r\n" (Windows / CRLF)
//   - "\r" (classic Mac / bare CR)
//
// Lines and columns are 1-based and measured in bytes (not runes).
// Position is a byte count, so a single multi-byte UTF-8 character
// occupies multiple column positions.
type LineTable struct {
	// lineStarts[i] is the byte offset of the start of line (i+1).
	// lineStarts[0] is always 0 (line 1 starts at the beginning of input).
	lineStarts []int
}

// NewLineTable builds a LineTable for the given input. O(n) time and
// O(lines) space.
func NewLineTable(input string) *LineTable {
	// Line 1 always starts at byte 0.
	starts := []int{0}
	for i := 0; i < len(input); i++ {
		switch input[i] {
		case '\n':
			// LF line break. Next line starts at i+1.
			starts = append(starts, i+1)
		case '\r':
			if i+1 < len(input) && input[i+1] == '\n' {
				// CRLF. Skip the \n and start the next line at i+2.
				starts = append(starts, i+2)
				i++ // consume the \n so we don't double-count it
			} else {
				// Bare CR. Next line starts at i+1.
				starts = append(starts, i+1)
			}
		}
	}
	return &LineTable{lineStarts: starts}
}

// Position returns the 1-based (line, column) for the given byte offset.
// If byteOffset is negative, returns (1, 1). If byteOffset exceeds the
// length of the input, it is clamped to the offset of the last known
// position.
//
// Column is measured in bytes: a byte offset at the start of a line
// returns column 1; the byte after returns column 2; and so on.
func (lt *LineTable) Position(byteOffset int) (line, col int) {
	if byteOffset < 0 {
		return 1, 1
	}
	if len(lt.lineStarts) == 0 {
		return 1, 1
	}
	// Binary search for the largest index i such that lineStarts[i] <= byteOffset.
	lo, hi := 0, len(lt.lineStarts)-1
	for lo < hi {
		mid := (lo + hi + 1) / 2
		if lt.lineStarts[mid] <= byteOffset {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	line = lo + 1                            // 1-based line number
	col = byteOffset - lt.lineStarts[lo] + 1 // 1-based column
	return line, col
}
```

- [ ] **Step 2: Write linetable_test.go**

Create `snowflake/parser/linetable_test.go`:

```go
package parser

import "testing"

func TestLineTable_Empty(t *testing.T) {
	lt := NewLineTable("")
	line, col := lt.Position(0)
	if line != 1 || col != 1 {
		t.Errorf("Position(0) on empty input = (%d, %d), want (1, 1)", line, col)
	}
}

func TestLineTable_SingleLine(t *testing.T) {
	lt := NewLineTable("abc")
	cases := []struct {
		offset       int
		wantLine     int
		wantCol      int
	}{
		{0, 1, 1}, // 'a'
		{1, 1, 2}, // 'b'
		{2, 1, 3}, // 'c'
		{3, 1, 4}, // just past 'c'
	}
	for _, c := range cases {
		line, col := lt.Position(c.offset)
		if line != c.wantLine || col != c.wantCol {
			t.Errorf("Position(%d) = (%d, %d), want (%d, %d)",
				c.offset, line, col, c.wantLine, c.wantCol)
		}
	}
}

func TestLineTable_LFLineEndings(t *testing.T) {
	// "a\nb\nc"  positions: a=0 \n=1 b=2 \n=3 c=4
	lt := NewLineTable("a\nb\nc")
	cases := []struct {
		offset       int
		wantLine     int
		wantCol      int
	}{
		{0, 1, 1}, // 'a'
		{1, 1, 2}, // '\n' (still line 1)
		{2, 2, 1}, // 'b'
		{3, 2, 2}, // '\n' (still line 2)
		{4, 3, 1}, // 'c'
	}
	for _, c := range cases {
		line, col := lt.Position(c.offset)
		if line != c.wantLine || col != c.wantCol {
			t.Errorf("Position(%d) = (%d, %d), want (%d, %d)",
				c.offset, line, col, c.wantLine, c.wantCol)
		}
	}
}

func TestLineTable_CRLFLineEndings(t *testing.T) {
	// "a\r\nb\r\nc"  positions: a=0 \r=1 \n=2 b=3 \r=4 \n=5 c=6
	// Line 1: starts at 0, contains a, \r, \n (through offset 2)
	// Line 2: starts at 3, contains b, \r, \n (through offset 5)
	// Line 3: starts at 6, contains c
	lt := NewLineTable("a\r\nb\r\nc")
	cases := []struct {
		offset       int
		wantLine     int
		wantCol      int
	}{
		{0, 1, 1}, // 'a'
		{1, 1, 2}, // '\r'
		{2, 1, 3}, // '\n' (still line 1 — CRLF is ONE break)
		{3, 2, 1}, // 'b'
		{4, 2, 2}, // '\r'
		{5, 2, 3}, // '\n'
		{6, 3, 1}, // 'c'
	}
	for _, c := range cases {
		line, col := lt.Position(c.offset)
		if line != c.wantLine || col != c.wantCol {
			t.Errorf("Position(%d) = (%d, %d), want (%d, %d)",
				c.offset, line, col, c.wantLine, c.wantCol)
		}
	}
}

func TestLineTable_CROnly(t *testing.T) {
	// "a\rb\rc"  positions: a=0 \r=1 b=2 \r=3 c=4
	lt := NewLineTable("a\rb\rc")
	cases := []struct {
		offset       int
		wantLine     int
		wantCol      int
	}{
		{0, 1, 1}, // 'a'
		{1, 1, 2}, // '\r' (still line 1)
		{2, 2, 1}, // 'b'
		{3, 2, 2}, // '\r'
		{4, 3, 1}, // 'c'
	}
	for _, c := range cases {
		line, col := lt.Position(c.offset)
		if line != c.wantLine || col != c.wantCol {
			t.Errorf("Position(%d) = (%d, %d), want (%d, %d)",
				c.offset, line, col, c.wantLine, c.wantCol)
		}
	}
}

func TestLineTable_OffsetPastEnd(t *testing.T) {
	lt := NewLineTable("abc")
	line, col := lt.Position(100)
	// Should clamp to the last known line (line 1), with col computed
	// from the last line start.
	if line != 1 {
		t.Errorf("Position(100) line = %d, want 1", line)
	}
	// col should be at least 4 (one past 'c')
	if col < 4 {
		t.Errorf("Position(100) col = %d, want >= 4", col)
	}
}

func TestLineTable_OffsetNegative(t *testing.T) {
	lt := NewLineTable("abc")
	line, col := lt.Position(-5)
	if line != 1 || col != 1 {
		t.Errorf("Position(-5) = (%d, %d), want (1, 1)", line, col)
	}
}
```

- [ ] **Step 3: Run linetable tests**

Run: `go test -v ./snowflake/parser/... -run TestLineTable 2>&1 | tail -30`
Expected: all 7 test functions reported as PASS.

- [ ] **Step 4: Verify full package still green**

Run: `go test ./snowflake/parser/...`
Expected: clean pass.

---

## Task 9: Write parser_test.go

**Files:**
- Modify: `snowflake/parser/parser_test.go`

- [ ] **Step 1: Replace the stub with the test helper and full test suite**

Overwrite `snowflake/parser/parser_test.go` with:

```go
package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// runParseBestEffortCase is a table-driven helper that asserts the shape
// of ParseBestEffort's result: expected statement count and expected
// number of errors (with optional per-error content matchers).
type parseCase struct {
	name         string
	input        string
	wantStmtCnt  int              // expected len(result.File.Stmts)
	wantErrCnt   int              // expected len(result.Errors)
	wantErrMsgs  []string         // expected substring in each error's Msg (len must match wantErrCnt)
	wantErrLocs  []int            // expected Loc.Start per error (len must match wantErrCnt); -1 to skip
}

func runParseBestEffortCase(t *testing.T, c parseCase) {
	t.Helper()
	result := ParseBestEffort(c.input)
	if got := len(result.File.Stmts); got != c.wantStmtCnt {
		t.Errorf("%s: stmt count = %d, want %d", c.name, got, c.wantStmtCnt)
	}
	if got := len(result.Errors); got != c.wantErrCnt {
		t.Errorf("%s: error count = %d, want %d", c.name, got, c.wantErrCnt)
		for i, e := range result.Errors {
			t.Errorf("  [%d] [%d,%d] %s", i, e.Loc.Start, e.Loc.End, e.Msg)
		}
		return
	}
	for i, want := range c.wantErrMsgs {
		if !strings.Contains(result.Errors[i].Msg, want) {
			t.Errorf("%s: error[%d] Msg = %q, want to contain %q",
				c.name, i, result.Errors[i].Msg, want)
		}
	}
	for i, wantLoc := range c.wantErrLocs {
		if wantLoc < 0 {
			continue
		}
		if result.Errors[i].Loc.Start != wantLoc {
			t.Errorf("%s: error[%d] Loc.Start = %d, want %d",
				c.name, i, result.Errors[i].Loc.Start, wantLoc)
		}
	}
}

func TestParse_EmptyInput(t *testing.T) {
	runParseBestEffortCase(t, parseCase{
		name:        "empty",
		input:       "",
		wantStmtCnt: 0,
		wantErrCnt:  0,
	})
}

func TestParse_WhitespaceOnly(t *testing.T) {
	runParseBestEffortCase(t, parseCase{
		name:        "whitespace only",
		input:       "   \n\t  ",
		wantStmtCnt: 0,
		wantErrCnt:  0,
	})
}

func TestParse_CommentOnly(t *testing.T) {
	cases := []parseCase{
		{name: "line comment", input: "-- comment\n", wantStmtCnt: 0, wantErrCnt: 0},
		{name: "block comment", input: "/* comment */", wantStmtCnt: 0, wantErrCnt: 0},
		{name: "slash line comment", input: "// comment\n", wantStmtCnt: 0, wantErrCnt: 0},
		{name: "nested block", input: "/* outer /* inner */ still */", wantStmtCnt: 0, wantErrCnt: 0},
	}
	for _, c := range cases {
		runParseBestEffortCase(t, c)
	}
}

func TestParse_SingleUnsupportedSelect(t *testing.T) {
	runParseBestEffortCase(t, parseCase{
		name:        "single SELECT",
		input:       "SELECT 1;",
		wantStmtCnt: 0,
		wantErrCnt:  1,
		wantErrMsgs: []string{"SELECT statement parsing is not yet supported"},
		wantErrLocs: []int{0},
	})
}

func TestParse_MultiUnsupported(t *testing.T) {
	// "SELECT 1; INSERT INTO t VALUES (1);" has SELECT at byte 0 and
	// INSERT at byte 10 (after "SELECT 1; ").
	runParseBestEffortCase(t, parseCase{
		name:        "SELECT then INSERT",
		input:       "SELECT 1; INSERT INTO t VALUES (1);",
		wantStmtCnt: 0,
		wantErrCnt:  2,
		wantErrMsgs: []string{
			"SELECT statement parsing is not yet supported",
			"INSERT statement parsing is not yet supported",
		},
		wantErrLocs: []int{0, 10},
	})
}

func TestParse_UnknownStatement(t *testing.T) {
	runParseBestEffortCase(t, parseCase{
		name:        "unknown keyword FOO",
		input:       "FOO BAR;",
		wantStmtCnt: 0,
		wantErrCnt:  1,
		wantErrMsgs: []string{"unknown or unsupported statement starting with FOO"},
		wantErrLocs: []int{0},
	})
}

func TestParse_LexErrorPropagated(t *testing.T) {
	// "'unterminated" triggers an unterminated-string LexError. It should
	// appear in result.Errors as a ParseError (lex errors are promoted).
	result := ParseBestEffort("'unterminated")
	if len(result.Errors) == 0 {
		t.Fatalf("expected at least one error, got none")
	}
	// Check that at least one error mentions "unterminated".
	found := false
	for _, e := range result.Errors {
		if strings.Contains(e.Msg, "unterminated") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected a 'unterminated' error, got: %+v", result.Errors)
	}
}

func TestParse_BeginEndBlockOneError(t *testing.T) {
	// F3's Split treats "BEGIN SELECT 1; END;" as ONE segment, so parseSingle
	// sees the whole block as one statement starting with BEGIN. The BEGIN
	// stub emits ONE ParseError for the whole block.
	runParseBestEffortCase(t, parseCase{
		name:        "BEGIN..END block",
		input:       "BEGIN SELECT 1; END;",
		wantStmtCnt: 0,
		wantErrCnt:  1,
		wantErrMsgs: []string{"BEGIN statement parsing is not yet supported"},
		wantErrLocs: []int{0},
	})
}

func TestParse_StrictVsBestEffort(t *testing.T) {
	// Parse returns the first error; ParseBestEffort returns all errors.
	input := "SELECT 1; INSERT INTO t VALUES (1);"

	file, err := Parse(input)
	if err == nil {
		t.Errorf("Parse: expected error, got nil")
	} else {
		pe, ok := err.(*ParseError)
		if !ok {
			t.Errorf("Parse: expected *ParseError, got %T", err)
		} else if !strings.Contains(pe.Msg, "SELECT") {
			t.Errorf("Parse: first error Msg = %q, want to contain SELECT", pe.Msg)
		}
	}
	if file == nil {
		t.Errorf("Parse: file is nil, want non-nil")
	}

	result := ParseBestEffort(input)
	if len(result.Errors) != 2 {
		t.Errorf("ParseBestEffort: got %d errors, want 2", len(result.Errors))
	}
}

func TestParse_StrictNoErrors(t *testing.T) {
	// Empty input has no errors; Parse should return (file, nil).
	file, err := Parse("")
	if err != nil {
		t.Errorf("Parse(\"\") error = %v, want nil", err)
	}
	if file == nil {
		t.Errorf("Parse(\"\") file = nil, want non-nil")
	}
}

func TestParse_AbsoluteSegmentPositions(t *testing.T) {
	// Given "SELECT 1; SELECT 2;", the second SELECT error's Loc.Start
	// should be 10 (the absolute byte position), NOT 0 (which would be
	// relative to the second segment).
	result := ParseBestEffort("SELECT 1; SELECT 2;")
	if len(result.Errors) != 2 {
		t.Fatalf("expected 2 errors, got %d: %+v", len(result.Errors), result.Errors)
	}
	if result.Errors[0].Loc.Start != 0 {
		t.Errorf("first error Loc.Start = %d, want 0", result.Errors[0].Loc.Start)
	}
	if result.Errors[1].Loc.Start != 10 {
		t.Errorf("second error Loc.Start = %d, want 10", result.Errors[1].Loc.Start)
	}
}

func TestParse_LegacyCorpus(t *testing.T) {
	// Smoke test: run ParseBestEffort against every .sql file in the
	// legacy corpus. Assert no panic; log the ParseError count per file
	// (most files will have many "X not supported" errors — that's fine
	// for stubs-only F4).
	corpusDir := "testdata/legacy"
	entries, err := os.ReadDir(corpusDir)
	if err != nil {
		t.Fatalf("failed to read corpus dir %s: %v", corpusDir, err)
	}

	fileCount := 0
	totalErrors := 0
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		fileCount++

		t.Run(entry.Name(), func(t *testing.T) {
			path := filepath.Join(corpusDir, entry.Name())
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read %s: %v", path, err)
			}
			// Must not panic.
			result := ParseBestEffort(string(data))
			totalErrors += len(result.Errors)
			t.Logf("%s: %d stmts, %d errors", entry.Name(), len(result.File.Stmts), len(result.Errors))
		})
	}

	if fileCount == 0 {
		t.Errorf("found 0 .sql files in %s — corpus missing?", corpusDir)
	}
	t.Logf("legacy corpus: %d files, %d total parse errors", fileCount, totalErrors)
}
```

- [ ] **Step 2: Run the tests**

Run: `go test ./snowflake/parser/...`
Expected: clean pass. All F2 + F3 + F4 tests pass.

- [ ] **Step 3: Verbose check for each F4 test function**

Run: `go test -v ./snowflake/parser/... -run TestParse 2>&1 | grep -E "^(=== RUN|--- (PASS|FAIL))" | tail -60`
Expected: every `TestParse_*` function (11 of them, plus the legacy corpus sub-tests) reported PASS.

If any test fails:
- For position errors, re-trace the byte offsets by hand
- For "X not supported" message mismatches, verify the dispatch switch's stmtName matches the case keyword
- For legacy corpus panics, find the failing file and re-run Parse on its content in isolation to debug

---

## Task 10: Final acceptance criteria sweep

- [ ] **Step 1: Build**

Run: `go build ./snowflake/parser/...`
Expected: no output, exit 0.

- [ ] **Step 2: Vet**

Run: `go vet ./snowflake/parser/...`
Expected: no output, exit 0.

- [ ] **Step 3: Gofmt**

Run: `gofmt -l snowflake/parser/`
Expected: no output.

If any file is listed, apply `gofmt -w snowflake/parser/` and re-run the check.

- [ ] **Step 4: Test**

Run: `go test ./snowflake/parser/...`
Expected:
```
ok  	github.com/bytebase/omni/snowflake/parser	(some duration)
```

- [ ] **Step 5: Verbose test run**

Run: `go test -v ./snowflake/parser/... 2>&1 | tail -80`
Expected: all F2, F3, F4 tests PASS. The legacy corpus sub-tests each report a stmt/error count via `t.Logf`.

- [ ] **Step 6: Package isolation**

Run: `go build ./snowflake/...`
Expected: no output, exit 0. Confirms `snowflake/parser` builds without depending on any package beyond F1's `snowflake/ast`.

- [ ] **Step 7: List all files for review**

Run: `git status snowflake/parser docs/superpowers`
Expected:
```
On branch feat/snowflake/parser-entry
Changes not staged for commit:
	modified:   snowflake/parser/errors.go
	modified:   snowflake/parser/lexer.go
	modified:   snowflake/parser/lexer_test.go
Untracked files:
	snowflake/parser/linetable.go
	snowflake/parser/linetable_test.go
	snowflake/parser/parser.go
	snowflake/parser/parser_test.go
```

Plus the already-committed spec and plan docs.

- [ ] **Step 8: STOP and present for review**

Do NOT commit. Output a summary:

> F4 (snowflake/parser parser-entry) implementation complete.
>
> All 10 tasks done. Acceptance criteria green:
> - `go build ./snowflake/parser/...` ✓
> - `go vet ./snowflake/parser/...` ✓
> - `gofmt -l snowflake/parser/` clean ✓
> - `go test ./snowflake/parser/...` ✓
> - F2 regression: all existing lexer tests still pass after the NewLexerWithOffset refactor ✓
> - Legacy corpus smoke test: all 27 files parse without panic ✓
> - `go build ./snowflake/...` (isolation) ✓
>
> Files created / modified:
> - snowflake/parser/parser.go (NEW, ~360 LOC)
> - snowflake/parser/parser_test.go (NEW, ~420 LOC)
> - snowflake/parser/linetable.go (NEW, ~90 LOC)
> - snowflake/parser/linetable_test.go (NEW, ~140 LOC)
> - snowflake/parser/lexer.go (MODIFIED: +baseOffset field, +NewLexerWithOffset, NextToken wrapper, Errors() shifting)
> - snowflake/parser/errors.go (MODIFIED: +ParseError type + Error() method)
> - snowflake/parser/lexer_test.go (MODIFIED: +3 offset tests)
> - docs/superpowers/plans/2026-04-09-snowflake-parser-entry.md (this plan)
>
> Please review the diff. After your approval I will commit and run finishing-a-development-branch.

---

## Spec Coverage Checklist

| Spec section | Covered by |
|---|---|
| Purpose (unblocks Tier 1+) | Implicit across all tasks |
| Architecture (split-then-parse) | Task 7 (ParseBestEffort + parseSingle) |
| F2 NewLexerWithOffset enhancement | Task 1 |
| Public API (Parse, ParseBestEffort, ParseResult, ParseError) | Tasks 2, 7 |
| ParseError with Error() method | Task 2 |
| Parser struct shape | Task 4 |
| Helpers (advance, peek, peekNext, match, expect, syntaxErrorAtCur, skipToNextStatement) | Task 4 |
| unsupported + unknownStatementError | Task 5 |
| 37-case dispatch switch | Task 6 |
| parseSingle | Task 7 |
| LineTable | Task 8 |
| Test suite for Parse/ParseBestEffort (11 categories) | Task 9 |
| LineTable tests (7 categories) | Task 8 |
| lexer_test.go offset additions (3 categories) | Task 1 |
| Legacy corpus smoke test | Task 9 (TestParse_LegacyCorpus) |
| Acceptance criteria | Task 10 |
