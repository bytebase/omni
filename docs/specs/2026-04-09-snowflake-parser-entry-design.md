# Snowflake Parser Entry + Walker (F4) — Design

**DAG node:** F4 — `snowflake/parser` parser-entry + walker
**Migration:** `docs/migration/snowflake/dag.md`
**Branch:** `feat/snowflake/parser-entry`
**Status:** Approved, ready for plan + implementation.
**Depends on:** F1 (ast-core, merged via PR #14), F2 (lexer, merged via PR #17), F3 (statement splitter, merged via PR #18).
**Unblocks:** Every Tier 1+ DAG node (T1.1 identifiers, T1.2 data types, T1.3 expressions, T1.4 SELECT core, T2.x DDL, T5.x DML, T7.1 Snowflake Scripting bodies).

## Purpose

F4 is the recursive-descent parser entry point for the omni Snowflake parser. It turns tokens produced by F2's Lexer into AST nodes defined in F1's `ast` package. Because concrete statement/expression nodes belong to Tier 1+ (T1.x, T2.x, T4.x, T5.x, T7.x), F4 ships **no actual statement parsing**. Instead, it ships:

1. The public `Parse` and `ParseBestEffort` API — the canonical way to invoke the snowflake parser from callers (including bytebase)
2. The `Parser` struct with internal state (lexer, token lookahead buffer, error list)
3. Helper primitives (`advance`, `peek`, `peekNext`, `match`, `expect`, `syntaxErrorAtCur`, `skipToNextStatement`)
4. A top-level statement dispatch switch targeting 37 Snowflake keywords, with one stub per keyword that returns a "statement type X not yet supported" `ParseError`
5. Error recovery via `skipToNextStatement` (sync on `;` at depth 0)
6. `ParseError` type
7. `LineTable` helper for byte-offset → (line, col) conversion

Once F4 lands, Tier 1+ contributors plug in concrete statement parsers by **replacing specific stubs** — the framework itself stays untouched.

## Architecture

F4 mirrors `mysql/parser/parser.go`'s split-then-parse pattern:

1. `Parse(input)` and `ParseBestEffort(input)` call F3's `Split(input)` to get `[]Segment`
2. For each segment, call `parseSingle(seg.Text, seg.ByteStart)` which constructs a `Parser{lexer: NewLexerWithOffset(seg.Text, seg.ByteStart)}` and drives a recursive-descent loop
3. Each statement's `Loc` values are absolute in the original input (thanks to `NewLexerWithOffset`)
4. Per-segment errors are collected and returned alongside the `*ast.File`

This is **deliberately different from mongo's one-pass approach** because mongo's `skipToNextStatement` would incorrectly fragment `BEGIN..END` blocks in the stubs-only F4. F3 already correctly groups `BEGIN..END` as a single segment, so reusing F3's Split gets the right behavior for free.

## Required F2 Enhancement: `NewLexerWithOffset`

F4 needs every token's `Loc.Start`/`Loc.End` to be absolute in the original input, not relative to a segment. The cleanest way is to add a `NewLexerWithOffset(input string, baseOffset int) *Lexer` constructor to F2's lexer that shifts all emitted token `Loc` values by `baseOffset`.

**Changes to `snowflake/parser/lexer.go`:**

1. Add an unexported `baseOffset int` field to the `Lexer` struct
2. Add the public constructor:

```go
// NewLexerWithOffset creates a Lexer whose emitted token Loc values are
// shifted by baseOffset. Use this when tokenizing a substring of a larger
// document and you want Loc values to refer to positions in the larger
// document rather than the substring.
func NewLexerWithOffset(input string, baseOffset int) *Lexer {
    return &Lexer{input: input, baseOffset: baseOffset}
}
```

3. In `NextToken`, before returning the token, add the offset to `Loc.Start` and `Loc.End`:

```go
func (l *Lexer) NextToken() Token {
    tok := l.nextTokenInner() // the existing body, refactored into an unexported helper
    if l.baseOffset != 0 && tok.Loc.Start >= 0 {
        tok.Loc.Start += l.baseOffset
        tok.Loc.End += l.baseOffset
    }
    return tok
}
```

The existing `NewLexer(input string) *Lexer` stays as a thin wrapper: `return NewLexerWithOffset(input, 0)`. All existing F2 tests remain valid because baseOffset defaults to 0.

Error recovery and internal helpers (`scanBlockComment`, `scanString`, etc.) continue to use `l.pos` as a local index into `l.input` — the offset is applied ONLY when building the final `Token.Loc` that's returned to the caller.

The `LexError.Loc` values that `scanString` and similar helpers append to `l.errors` also need to be shifted by `baseOffset`. Since those helpers already write the offset-less positions, the fix is to shift them at the `Lexer.Errors()` call site — OR to wrap every error construction in a helper that adds the offset. We'll use the latter: add a private `errorAt(start, end int, msg string)` helper that applies the offset, and replace all current `l.errors = append(l.errors, LexError{...})` call sites with `l.errorAt(...)`.

**Test additions to `lexer_test.go`:**

A new `TestLexer_NewLexerWithOffset` that tokenizes `" SELECT "` with baseOffset=100 and asserts the `kwSELECT` token has `Loc.Start = 101, Loc.End = 107`, and that errors from `"'unterminated"` with baseOffset=50 report `LexError.Loc.Start = 50`.

## Public API

```go
package parser

import (
    "github.com/bytebase/omni/snowflake/ast"
)

// Parse parses input as a complete Snowflake SQL file and returns the
// resulting *ast.File. If any parse errors occur, the first error is
// returned; the *ast.File may still contain any successfully-parsed
// statements that precede the error.
//
// Most callers should use ParseBestEffort instead — it always returns
// partial results alongside the full error list.
func Parse(input string) (*ast.File, error)

// ParseBestEffort parses input as a complete Snowflake SQL file, recovering
// from per-statement errors and always returning a ParseResult containing
// every successfully-parsed statement and every encountered error.
//
// This is the canonical entry point for bytebase consumers (diagnose.go,
// query_type.go, query_span_extractor.go) that need partial results plus
// error diagnostics.
func ParseBestEffort(input string) *ParseResult

// ParseResult holds the outcome of a best-effort parse.
type ParseResult struct {
    File   *ast.File
    Errors []ParseError
}
```

`ParseError` lives in `errors.go` alongside the existing `LexError`:

```go
// ParseError describes a single parse error with its source location.
//
// Loc uses the same shape as LexError so consumers can handle both uniformly.
// Line/column conversion is a caller-side concern via LineTable.
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

`LineTable` lives in `linetable.go`:

```go
// LineTable provides O(log n) byte-offset → (line, column) conversion for
// a single source string. Construct once per input, then call Position
// repeatedly.
//
// Handles LF, CRLF, and bare-CR line endings. Lines and columns are 1-based
// and measured in bytes (not runes).
type LineTable struct {
    // unexported: cumulative byte offset of the start of each line
    lineStarts []int
}

// NewLineTable builds a LineTable for the given input. O(n) time.
func NewLineTable(input string) *LineTable

// Position returns the 1-based (line, column) for the given byte offset.
// Offsets past the end of the input are clamped to the last position.
func (lt *LineTable) Position(byteOffset int) (line, col int)
```

## Parser Struct

```go
// Parser is a recursive-descent parser for Snowflake SQL. It operates on
// a single segment of input (one top-level statement) at a time; callers
// should use Parse or ParseBestEffort instead of constructing Parsers
// directly.
type Parser struct {
    lexer   *Lexer
    input   string    // the segment text (used for line/col computation in error messages)
    cur     Token     // current token
    prev    Token     // previous token (for error context)
    nextBuf Token     // buffered lookahead
    hasNext bool      // whether nextBuf is valid
    errors  []ParseError // collected errors for best-effort mode
}
```

The `errors` field is used by `ParseBestEffort`; `Parse` returns the first error and halts its segment.

Unlike `mysql/parser.Parser`, F4's Parser has **no completion-mode fields**. Auto-completion lives in a separate package (`snowflake/completion`) that wraps or extends the Parser. F4 keeps the Parser lean.

## Helpers

Direct port of mysql/mongo patterns:

```go
// advance consumes the current token and moves to the next.
func (p *Parser) advance() Token

// peek returns the current token without consuming it.
func (p *Parser) peek() Token

// peekNext returns the token AFTER the current one without consuming it
// (one-token lookahead beyond cur).
func (p *Parser) peekNext() Token

// match tries each given token type; if cur matches any, it is consumed
// and returned with ok=true.
func (p *Parser) match(types ...int) (Token, bool)

// expect consumes the current token if it matches the expected type.
// Otherwise returns a ParseError at the current token position.
func (p *Parser) expect(tokenType int) (Token, error)

// syntaxErrorAtCur constructs a ParseError describing a syntax error at
// the current token position. If cur is tokEOF, the message says "at end
// of input"; otherwise "at or near X".
func (p *Parser) syntaxErrorAtCur() *ParseError

// skipToNextStatement advances past the current erroneous statement to
// the next ; at depth 0, or EOF. It always consumes at least one token
// to guarantee progress.
//
// Because Parse uses F3.Split to pre-segment the input, each segment
// contains at most one top-level statement. skipToNextStatement therefore
// typically advances to EOF within the segment, which is fine — the outer
// ParseBestEffort loop moves to the next segment.
func (p *Parser) skipToNextStatement()
```

## Top-Level Dispatch

```go
// parseStmt parses one top-level statement. It dispatches on the first
// keyword. For stubs-only F4, every dispatch case calls an unsupported
// helper that appends a ParseError and returns (nil, err).
//
// Tier 1+ DAG nodes REPLACE specific cases with real implementations.
// The dispatch switch itself is not expected to change — only the bodies
// of individual case branches.
func (p *Parser) parseStmt() (ast.Node, error) {
    switch p.cur.Type {
    // DDL
    case kwCREATE:   return p.unsupported("CREATE")
    case kwALTER:    return p.unsupported("ALTER")
    case kwDROP:     return p.unsupported("DROP")
    case kwUNDROP:   return p.unsupported("UNDROP")
    case kwTRUNCATE: return p.unsupported("TRUNCATE")
    case kwCOMMENT:  return p.unsupported("COMMENT")

    // DML
    case kwSELECT: return p.unsupported("SELECT")
    case kwWITH:   return p.unsupported("WITH")
    case kwINSERT: return p.unsupported("INSERT")
    case kwUPDATE: return p.unsupported("UPDATE")
    case kwDELETE: return p.unsupported("DELETE")
    case kwMERGE:  return p.unsupported("MERGE")
    case kwCOPY:   return p.unsupported("COPY")
    case kwPUT:    return p.unsupported("PUT")
    case kwGET:    return p.unsupported("GET")
    case kwLIST:   return p.unsupported("LIST")
    case kwREMOVE: return p.unsupported("REMOVE")
    case kwCALL:   return p.unsupported("CALL")

    // DCL
    case kwGRANT:  return p.unsupported("GRANT")
    case kwREVOKE: return p.unsupported("REVOKE")

    // TCL
    case kwBEGIN:    return p.unsupported("BEGIN")
    case kwSTART:    return p.unsupported("START")
    case kwCOMMIT:   return p.unsupported("COMMIT")
    case kwROLLBACK: return p.unsupported("ROLLBACK")

    // Info / session
    case kwSHOW:     return p.unsupported("SHOW")
    case kwDESCRIBE: return p.unsupported("DESCRIBE")
    case kwDESC:     return p.unsupported("DESC")
    case kwEXPLAIN:  return p.unsupported("EXPLAIN")
    case kwUSE:      return p.unsupported("USE")
    case kwSET:      return p.unsupported("SET")
    case kwUNSET:    return p.unsupported("UNSET")

    // Snowflake Scripting (subset available as F2 keywords)
    case kwDECLARE:  return p.unsupported("DECLARE")
    case kwIF:       return p.unsupported("IF")
    case kwCASE:     return p.unsupported("CASE")
    case kwFOR:      return p.unsupported("FOR")
    case kwRETURN:   return p.unsupported("RETURN")
    case kwCONTINUE: return p.unsupported("CONTINUE")

    default:
        return nil, p.unknownStatementError()
    }
}
```

**37 keywords** are in the switch. The following 7 Snowflake Scripting keywords are NOT included because they are missing from F2's keyword table: `SAVEPOINT`, `EXECUTE`, `WHILE`, `REPEAT`, `LET`, `LOOP`, `BREAK`. Five of these are absent from the legacy grammar entirely; the other two (`EXECUTE`, with its unusual `'EXEC' 'UTE'?` rule form, and `SAVEPOINT`) were not extracted by F2's grep-based pipeline. When Tier 7 Snowflake Scripting support lands, contributors will add the missing `kw*` constants to F2 and append the corresponding dispatch cases here.

The `default` branch catches any token that isn't one of the 37. For an unknown keyword-like token, `unknownStatementError` reports "unknown or unsupported statement starting with X". For a non-keyword token (like a bare identifier or operator), it reports "expected statement, got X".

## Per-Statement Stub Helper

```go
// unsupported emits a "statement type X not yet supported" ParseError at
// the current token, advances to the next statement boundary, and returns
// (nil, err). Used by every stubbed dispatch case in F4.
//
// Tier 1+ DAG nodes replace individual dispatch cases with real parse
// functions. Those real functions do NOT call unsupported — they consume
// tokens themselves.
func (p *Parser) unsupported(stmtName string) (ast.Node, error) {
    err := &ParseError{
        Loc: p.cur.Loc,
        Msg: stmtName + " statement parsing is not yet supported",
    }
    p.skipToNextStatement()
    return nil, err
}

// unknownStatementError reports a statement that starts with a token the
// dispatch switch doesn't recognize.
func (p *Parser) unknownStatementError() *ParseError {
    if p.cur.Type == tokEOF {
        return &ParseError{
            Loc: p.cur.Loc,
            Msg: "syntax error at end of input",
        }
    }
    tokText := p.cur.Str
    if tokText == "" {
        tokText = "token " + TokenName(p.cur.Type)
    }
    return &ParseError{
        Loc: p.cur.Loc,
        Msg: "unknown or unsupported statement starting with " + tokText,
    }
}
```

## Parse / ParseBestEffort Implementation

```go
// Parse is the strict entry point: returns the first error if any statement
// fails. Callers that want partial results should use ParseBestEffort.
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
```

**Lexer error surfacing:** Each `parseSingle` invocation has its own `Parser` with its own `Lexer` (constructed via `NewLexerWithOffset` so positions are absolute in the original input). After the parse loop for a segment finishes (either by reaching EOF within the segment or by error recovery), `parseSingle` drains `p.lexer.Errors()` and promotes each `LexError` to a `ParseError`, appending it to the segment's error list. This keeps error handling local to each segment and avoids re-lexing the whole input.

```go
// parseSingle parses one segment into a single ast.Node. Returns (node, errors)
// where node may be nil if the segment failed to parse and errors is the list
// of ParseErrors (including any LexErrors promoted from the underlying Lexer).
func parseSingle(segText string, baseOffset int) (ast.Node, []ParseError) {
    p := &Parser{
        lexer: NewLexerWithOffset(segText, baseOffset),
        input: segText,
    }
    p.advance()

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

    // Promote any lex errors into ParseErrors.
    for _, le := range p.lexer.Errors() {
        p.errors = append(p.errors, ParseError{Loc: le.Loc, Msg: le.Msg})
    }

    return result, p.errors
}
```

Note: F2's `Lexer.Errors()` returns the error slice accumulated up to that point. Since `parseSingle` drives the lexer via `advance`/`NextToken` until EOF (the main loop in parseStmt should consume the whole segment or fail), calling `Errors()` after the parse captures everything.

## Error Recovery

The `skipToNextStatement` helper is simpler than mongo's because F3's Split has already grouped statements correctly:

```go
// skipToNextStatement advances the parser past the current erroneous
// statement to the next ; at depth 0 within the current segment, or to
// tokEOF. Always consumes at least one token.
//
// Because Parse uses F3.Split to pre-segment the input, each segment
// usually contains one statement — skipToNextStatement's typical behavior
// is to advance to EOF within the segment.
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

## Out of Scope

F4 does NOT include:

| Feature | Where it lives |
|---|---|
| Concrete statement parsing (SELECT, INSERT, CREATE, etc.) | Tier 1+ DAG nodes |
| Expression parsing | T1.3 |
| Data type parsing | T1.2 |
| Identifier / qualified-name helpers beyond basic tokIdent reading | T1.1 |
| Snowflake Scripting body parsing (BEGIN..END bodies, IF/FOR/WHILE/LOOP/REPEAT) | T7.1 |
| Completion mode (LSP-like candidate collection) | `snowflake/completion` (separate package, not a DAG node yet) |
| UTF-16 byte-offset conversion | Bytebase consumer handles this on top of LineTable |

## Testing

Test files: `snowflake/parser/parser_test.go`, `snowflake/parser/linetable_test.go`, and additions to `snowflake/parser/lexer_test.go`. Run with `go test ./snowflake/parser/...`.

### parser_test.go

1. **TestParse_EmptyInput** — `Parse("")` returns `*ast.File{Stmts: nil, Loc: {0,0}}`, no errors
2. **TestParse_WhitespaceOnly** — same for `"   \n\t"`
3. **TestParse_CommentOnly** — same for `"-- comment\n"` and `"/* block */"`
4. **TestParse_SingleUnsupportedSelect** — `Parse("SELECT 1;")` returns a File with no Stmts and one ParseError with `Msg: "SELECT statement parsing is not yet supported"` and `Loc.Start: 0`
5. **TestParse_MultiUnsupported** — `Parse("SELECT 1; INSERT INTO t VALUES (1);")` returns two ParseErrors, each with segment-absolute positions (SELECT at 0, INSERT at 10)
6. **TestParse_UnknownStatement** — `Parse("FOO BAR;")` returns one ParseError with `Msg: "unknown or unsupported statement starting with FOO"`
7. **TestParse_LexErrorPropagated** — `Parse("'unterminated")` returns the LexError promoted to ParseError
8. **TestParse_BeginEndBlock** — `Parse("BEGIN SELECT 1; END;")` returns one ParseError (for `BEGIN`) because F3's Split treats this as one segment
9. **TestParse_StrictVsBestEffort** — same input produces:
   - `Parse(input)` returns `(file, firstError)` if errors exist, else `(file, nil)`
   - `ParseBestEffort(input)` returns `*ParseResult{File: file, Errors: [...all errors...]}`
10. **TestParse_SegmentPositionsAbsolute** — for `Parse("BEGIN_STMT; SELECT 1;")`, the SELECT error's `Loc.Start` is 12, NOT 0 (verifies baseOffset plumbing)
11. **TestParse_LegacyCorpus** — iterate every `.sql` in `testdata/legacy/` and call Parse; assert no panic; log the ParseError count per file (most files will have many "X not supported" errors)

### linetable_test.go

1. **TestLineTable_Empty** — empty input returns (1, 1) for offset 0
2. **TestLineTable_SingleLine** — `"abc"` returns (1, 1), (1, 2), (1, 3), (1, 4) for offsets 0, 1, 2, 3
3. **TestLineTable_LFLineEndings** — `"a\nb\nc"` returns (1,1), (1,2), (2,1), (2,2), (3,1) for offsets 0..4
4. **TestLineTable_CRLFLineEndings** — `"a\r\nb\r\nc"` — CRLF counts as one line break
5. **TestLineTable_CROnly** — `"a\rb\rc"` — bare CR counts as one line break (rare but valid)
6. **TestLineTable_OffsetPastEnd** — offset > len(input) is clamped to the last position
7. **TestLineTable_OffsetNegative** — negative offset returns (1, 1) (defensive)

### lexer_test.go additions

1. **TestLexer_NewLexerWithOffset** — tokenize `"   SELECT   "` with baseOffset=100; assert the `kwSELECT` token has `Loc.Start = 103, Loc.End = 109`; assert all other tokens (if any) are similarly shifted
2. **TestLexer_NewLexerWithOffset_LexErrors** — tokenize `"'unterminated"` with baseOffset=50; assert `LexError.Loc.Start = 50, LexError.Loc.End = 63`
3. **TestLexer_NewLexer_UnchangedBehavior** — the existing `NewLexer` continues to produce tokens with zero-based offsets (regression check)

## Acceptance Criteria

F4 is complete when:

1. `go build ./snowflake/parser/...` succeeds.
2. `go vet ./snowflake/parser/...` clean.
3. `gofmt -l snowflake/parser/` clean.
4. `go test ./snowflake/parser/...` passes — all F2/F3 tests still green, plus all F4 tests (parser, linetable, lexer-offset).
5. `Parse` and `ParseBestEffort` are exported from `package parser`.
6. All 37 listed dispatch keywords route to a stub that emits a "not yet supported" ParseError with the correct position.
7. The legacy corpus smoke test runs without panic for all 27 files.
8. `LineTable` round-trip tests pass for LF, CRLF, and CR line endings.
9. F2's `NewLexerWithOffset` is exported and all existing F2 tests still pass with zero regressions.
10. `go build ./snowflake/...` (isolation check) — F4 builds without depending on any package beyond F1+F2+F3 (all intra-package).
11. After merge, `docs/migration/snowflake/dag.md` F4 status is flipped to `done`.

## Files Created / Modified

### Created

- `snowflake/parser/parser.go` — Parser struct, Parse, ParseBestEffort, parseSingle, helpers, dispatch switch, unsupported/unknownStatementError
- `snowflake/parser/linetable.go` — LineTable, NewLineTable, Position
- `snowflake/parser/parser_test.go` — Parse/ParseBestEffort tests
- `snowflake/parser/linetable_test.go` — LineTable round-trip tests

### Modified

- `snowflake/parser/lexer.go` — Add `baseOffset` field, `NewLexerWithOffset` constructor, offset-shifting in NextToken and all error-construction call sites
- `snowflake/parser/errors.go` — Add `ParseError` type alongside existing `LexError`
- `snowflake/parser/lexer_test.go` — Add offset-based test cases

Estimated total: ~900 LOC of new code + ~40 LOC modified in lexer.go + errors.go.
