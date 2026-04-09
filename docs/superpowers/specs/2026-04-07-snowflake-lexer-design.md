# Snowflake Lexer (F2) — Design

**DAG node:** F2 — `snowflake/parser` lexer
**Migration:** `docs/migration/snowflake/dag.md`
**Branch:** `feat/snowflake/lexer`
**Status:** Approved, ready for plan + implementation.
**Depends on:** F1 (ast-core, merged via PR #14) — uses `ast.Loc` for token position tracking.
**Unblocks:** F3 (statement-splitter), F4 (parser-entry).

## Purpose

F2 is the hand-written tokenizer for the Snowflake parser. It is the next critical-path node after F1 and unblocks F3 (statement-splitter) and F4 (parser-entry). It must:

1. Tokenize Snowflake SQL with full coverage of the legacy ANTLR4 `bytebase/parser/snowflake/SnowflakeLexer.g4` (1,275 lines, 967 token rules).
2. Track byte positions for every token using F1's `ast.Loc{Start, End}` so F4 can attach Locs to AST nodes and bytebase's `statement_ranges.go` consumer can extract statement boundaries.
3. Be consumable by both F3 (which needs only token kinds and positions) and F4 (which needs full token semantics including keyword classification and literal values).

F2 is intentionally lexer-only. Statement splitting, recursive-descent parsing, and AST construction are out of scope.

## Architecture

Hand-written, mysql-style flat lexer. No state machine. A single `Lexer` struct holding the input string, current position, and an error slice. `NextToken()` is the streaming primary API; `Tokenize()` is a one-shot convenience.

The structural template is `mysql/parser/lexer.go`, not `pg/parser/lexer.go`. PG's lexer is a near-direct port of PostgreSQL's `scan.l` flex grammar with multiple `LexerState` modes for handling extended quote variants, Unicode escape continuations, and dollar-quoted strings. Snowflake doesn't need any of that complexity — its string forms are simpler (`'...'` with `''`/`\\` escapes plus `$$...$$`), and its operator surface is small enough that flat character lookups + 1-or-2-byte lookahead handles every disambiguation.

The architectural decision to drop pg's state machine is documented in `docs/migration/snowflake/analysis.md` ("Reference engine") and informed by reading both `pg/parser/lexer.go` (1,312 lines) and `mysql/parser/lexer.go` (2,105 lines).

## Package Layout

```
snowflake/parser/
  lexer.go                       # Lexer struct, NewLexer, NextToken, Tokenize, scan* helpers, scanString, scanNumber, scanIdentOrKeyword, scanQuotedIdent, scanDollar, scanOperatorOrPunct, skipWhitespaceAndComments, scanBlockComment
  tokens.go                      # tokXxx constants (600 block), kwXxx keyword constants (700+ block), Token struct, TokenName(int) string for debug
  keywords.go                    # keywordMap[string]int, keywordReserved[int]bool, IsReservedKeyword, KeywordToken
  errors.go                      # LexError type and constant error message strings
  lexer_test.go                  # Table-driven coverage of every token kind
  position_test.go               # Position math edge cases
  edge_cases_test.go             # EOF mid-string, unterminated comment, ambiguous operators, numeric edge cases, Unicode
  error_recovery_test.go         # Multi-error recovery
  legacy_corpus_test.go          # Smoke-tokenize all 27 .sql files in testdata/legacy/
  keyword_completeness_test.go   # Drift detector against legacy SnowflakeLexer.g4
```

The keyword table is split into its own file because at ~970 entries it would dwarf the lexer body if inlined. mysql crams everything into a single 2,105-line `lexer.go`; Snowflake's keyword count makes that approach unmaintainable.

Estimated total: **~3,500 LOC** across 4 production files + 6 test files. The keyword table is ~970 lines on its own, the lexer body ~800 LOC, the rest is tests.

## Type Definitions

### `tokens.go`

```go
package parser

// Special tokens.
const (
    tokEOF     = 0
    tokInvalid = 1 // synthetic error-recovery token; always accompanied by a LexError
)

// Multi-character operators and literal token kinds. Single-character tokens
// (+, -, *, /, %, ~, (, ), [, ], {, }, ,, ;, ., @, :, <, >, =) use the ASCII
// byte value directly as their Type — no constant needed.
const (
    tokInt         = 600 + iota // numeric literal: integer (legacy DECIMAL)
    tokFloat                    // numeric literal: 1.5, 1., .5  (legacy FLOAT)
    tokReal                     // numeric literal: 1.5e10, 1e-10 (legacy REAL)
    tokString                   // 'string' or $$string$$  (legacy STRING / DBL_DOLLAR)
    tokIdent                    // unquoted identifier (preserves source case)
    tokQuotedIdent              // "quoted identifier" (preserves source bytes verbatim)
    tokVariable                 // $name Snowflake script variable (legacy ID2)

    tokDoubleColon              // ::  (legacy COLON_COLON)
    tokConcat                   // ||  (legacy PIPE_PIPE)
    tokArrow                    // ->  (legacy ARROW)
    tokFlow                     // ->> (legacy FLOW)
    tokAssoc                    // =>  (legacy ASSOC)
    tokNotEq                    // != or <> (legacy NE / LTGT — folded into one token)
    tokLessEq                   // <=  (legacy LE)
    tokGreaterEq                // >=  (legacy GE)
)

// Keyword token constants. Values start at 700 and grow with iota. Each
// constant corresponds to one keyword from SnowflakeLexer.g4. Order is
// alphabetical for human readability; numeric values are NOT stable across
// edits — do not persist them.
const (
    kwACCOUNT = 700 + iota
    kwACCOUNTS
    kwACCOUNTADMIN
    // ... approximately 970 entries
)

// Token represents a single lexical token.
type Token struct {
    Type    int      // tok* / kw* / ASCII byte value
    Str     string   // identifier text, string content, raw operator text
    Ival    int64    // integer value for tokInt
    Loc     ast.Loc  // {Start, End} byte offsets in source text
    XPrefix bool     // true if string was X'...' (binary literal); only set on tokString
}

// TokenName returns a human-readable name for a token type. Used by tests
// and debug output. Returns:
//   - "EOF" for tokEOF
//   - "INVALID" for tokInvalid
//   - the upper-case keyword name (e.g. "SELECT") for kw* constants
//   - the constant name minus the "tok" prefix (e.g. "INT", "STRING") for tok* constants
//   - the literal character in quotes (e.g. "'+'") for ASCII single-char tokens
func TokenName(t int) string
```

### `keywords.go`

```go
package parser

// keywordMap maps lowercased keyword text to its kw* constant.
// Snowflake matches keywords case-insensitively, so all lookups go through
// strings.ToLower of the source identifier.
var keywordMap = map[string]int{
    "account": kwACCOUNT,
    "accounts": kwACCOUNTS,
    "accountadmin": kwACCOUNTADMIN,
    // ... approximately 970 entries
}

// keywordReserved marks the ~99 reserved keywords that cannot be used as
// unquoted identifiers without quoting. Seeded from
// build_id_contains_non_reserved_keywords.py in the legacy parser.
var keywordReserved = map[int]bool{
    kwACCOUNT: true,
    kwALL: true,
    kwALTER: true,
    // ... 99 entries
}

// KeywordToken returns the kw* constant for name (case-insensitive) and
// true if name is a recognized keyword.
func KeywordToken(name string) (int, bool)

// IsReservedKeyword reports whether name (case-insensitive) is a reserved
// keyword that cannot be used as an unquoted identifier.
func IsReservedKeyword(name string) bool
```

### `lexer.go`

```go
package parser

import "github.com/bytebase/omni/snowflake/ast"

// Lexer is a Snowflake SQL tokenizer.
type Lexer struct {
    input  string
    pos    int        // current byte offset (one past the last consumed byte)
    start  int        // start byte of the token currently being scanned
    errors []LexError
}

// NewLexer constructs a Lexer for the given input.
func NewLexer(input string) *Lexer

// NextToken advances the lexer and returns the next token. At end of input
// it returns Token{Type: tokEOF}; subsequent calls continue to return EOF.
func (l *Lexer) NextToken() Token

// Errors returns all lex errors collected so far. Errors are appended in
// source order; each error is accompanied by a tokInvalid token in the
// stream so consumers can choose to halt or proceed.
func (l *Lexer) Errors() []LexError

// Tokenize is a one-shot convenience that runs a lexer to EOF and returns
// the full token stream and error list. Useful for tests and for callers
// that don't need streaming.
func Tokenize(input string) (tokens []Token, errors []LexError)
```

### `errors.go`

```go
package parser

import "github.com/bytebase/omni/snowflake/ast"

// LexError describes a single lexing failure with its source location.
type LexError struct {
    Loc ast.Loc
    Msg string
}

// Standard lex error messages. Tests assert against these constants so a
// reword in one place propagates everywhere.
const (
    errUnterminatedString  = "unterminated string literal"
    errUnterminatedDollar  = "unterminated $$ string literal"
    errUnterminatedComment = "unterminated block comment"
    errUnterminatedQuoted  = "unterminated quoted identifier"
    errInvalidByte         = "invalid byte"
)
```

## Scanning Loop

`NextToken` calls `skipWhitespaceAndComments` first, then dispatches based on the first non-whitespace byte:

```go
func (l *Lexer) NextToken() Token {
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
        return l.scanDollar() // $$ → string; $NAME → variable; $ → ASCII byte
    case ch >= '0' && ch <= '9':
        return l.scanNumber()
    case ch == '.' && l.peekDigit():
        return l.scanNumber() // .5 form
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

`scanOperatorOrPunct` does the multi-char lookahead for `::`, `||`, `->`, `->>`, `=>`, `!=`, `<>`, `<=`, `>=` and falls back to `Token{Type: int(ch), ...}` for everything else.

## Comments

```go
func (l *Lexer) skipWhitespaceAndComments() {
    for l.pos < len(l.input) {
        ch := l.input[l.pos]
        switch {
        case ch == ' ', ch == '\t', ch == '\n', ch == '\r':
            l.pos++
        case ch == '-' && l.peek(1) == '-':
            l.pos += 2
            for l.pos < len(l.input) && l.input[l.pos] != '\n' {
                l.pos++
            }
        case ch == '/' && l.peek(1) == '/':
            l.pos += 2
            for l.pos < len(l.input) && l.input[l.pos] != '\n' {
                l.pos++
            }
        case ch == '/' && l.peek(1) == '*':
            l.scanBlockComment()
        default:
            return
        }
    }
}

func (l *Lexer) scanBlockComment() {
    start := l.pos
    l.pos += 2
    depth := 1
    for l.pos < len(l.input) && depth > 0 {
        if l.input[l.pos] == '/' && l.peek(1) == '*' {
            depth++
            l.pos += 2
        } else if l.input[l.pos] == '*' && l.peek(1) == '/' {
            depth--
            l.pos += 2
        } else {
            l.pos++
        }
    }
    if depth > 0 {
        l.errors = append(l.errors, LexError{
            Loc: ast.Loc{Start: start, End: l.pos},
            Msg: errUnterminatedComment,
        })
    }
}
```

Snowflake's three comment forms — `--`, `//`, `/* ... */` — are dropped silently. Block comments nest; the legacy grammar uses recursive `'/*' (SQL_COMMENT | .)*? '*/'` which the depth counter implements.

## String Literals

`scanString` handles both single-quoted strings and `$$...$$` double-dollar strings (the latter reached via `scanDollar`):

- `'...'` form: read until matching `'`, treating `\\<char>` as an escape and `''` as a doubled-quote escape (literal single quote). Recognized escape sequences: `\n` (newline), `\t` (tab), `\r` (carriage return), `\0` (NUL), `\b` (backspace, 0x08), `\\` (backslash), `\'` (single quote), `\"` (double quote). Any other `\<char>` sequence is reduced to the literal `<char>`. The `Token.Str` field contains the unescaped content.
- `$$...$$` form: read until next `$$`, no escapes. The `Token.Str` field contains the raw content between the delimiters.

Unterminated strings produce a `LexError` and a `tokInvalid` token covering the bad span; the lexer advances to the end of input or to a newline (whichever comes first) and continues.

## Identifiers and Keywords

```go
func (l *Lexer) scanIdentOrKeyword() Token {
    for l.pos < len(l.input) && isIdentCont(l.input[l.pos]) {
        l.pos++
    }
    text := l.input[l.start:l.pos]
    if t, ok := KeywordToken(text); ok {
        return Token{
            Type: t,
            Str:  text,
            Loc:  ast.Loc{Start: l.start, End: l.pos},
        }
    }
    return Token{
        Type: tokIdent,
        Str:  text,
        Loc:  ast.Loc{Start: l.start, End: l.pos},
    }
}
```

`isIdentStart` accepts `[A-Za-z_]`. `isIdentCont` accepts `[A-Za-z0-9_@$]`. The case folding for keyword lookup happens inside `KeywordToken` via `strings.ToLower`. The `Token.Str` field always preserves source-text bytes (no case folding) so the resolver can decide what to do with case-sensitive contexts later.

`scanQuotedIdent` reads `"..."` content verbatim, handling `""` as a doubled-quote escape. The result is `tokQuotedIdent` with `Str` set to the unescaped content (not including the surrounding quotes). The empty form `""` (zero content) is also tokenized as `tokQuotedIdent` with `Str = ""`.

**Deviation from legacy.** The legacy grammar splits quoted identifiers into two token kinds: `DOUBLE_QUOTE_ID : '"' ~'"'+ '"'` (one or more non-quote characters) and `DOUBLE_QUOTE_BLANK : '""'` (the empty form). F2 folds both into a single `tokQuotedIdent` kind because the parser does not need to distinguish them — the empty case is a degenerate value of the same shape. The drift detector test (`keyword_completeness_test.go`) does NOT check `DOUBLE_QUOTE_BLANK` since it's a deliberate fold.

## Numeric Literals

Three token kinds match the legacy grammar:

- `tokInt`: digits only (e.g. `42`, `0`)
- `tokFloat`: digits with a dot (e.g. `1.5`, `1.`, `.5`)
- `tokReal`: digits with an `E`/`e` exponent (e.g. `1e10`, `1.5e-10`, `.5e+10`)

`scanNumber` decides which kind based on the consumed pattern.

Field population:
- All three kinds populate `Token.Str` with the verbatim source text (so the parser can deparse losslessly and emit precise error messages).
- `tokInt` additionally populates `Token.Ival` via `strconv.ParseInt(Str, 10, 64)`. If the integer overflows int64, `scanNumber` records a `LexError` and downgrades the token to `tokFloat` (preserving the source text in `Str`) — this matches the way Snowflake's NUMBER(38, 0) accepts arbitrary-precision integers that don't fit in 64 bits.
- `tokFloat` and `tokReal` do NOT populate `Token.Ival` — the parser converts via `strconv.ParseFloat(Str, 64)` if it needs a numeric value.

## Hex and Binary Literals

`X'...'` and `x'...'` are tokenized as `tokString` with `XPrefix = true`. The lexer does not interpret the hex content — that's a parser-level concern. This matches the legacy grammar where `X'...'` has no separate token kind.

## Variables

`$NAME` (where NAME is `[A-Z0-9_]+`) is tokenized as `tokVariable` with `Str = "NAME"` (no leading `$`). Bare `$` (not followed by an identifier-cont character or another `$`) is the ASCII single-char token `int('$')`.

`$$...$$` is handled before the variable case in `scanDollar`.

## Error Recovery

On any lex error:
1. Append a `LexError{Loc, Msg}` to `l.errors`.
2. Emit a `tokInvalid` token covering the bad span.
3. Advance `l.pos` past the error site according to the error type:
   - **Invalid byte**: advance one byte.
   - **Unterminated single-quoted string**: advance to the next newline or EOF, whichever comes first.
   - **Unterminated `$$` string**: advance to EOF.
   - **Unterminated quoted identifier**: advance to the next newline or EOF, whichever comes first.
   - **Unterminated block comment**: advance to EOF.
4. Continue tokenizing from the new position.

This mirrors mongo's recent error-recovery direction (commit 27fef3f). Bytebase's `diagnose.go` consumer iterates `Lexer.Errors()` to report syntax errors with positions; F3/F4 can choose to halt on the first lex error or proceed.

## Public API Summary

```go
// Construction
func NewLexer(input string) *Lexer
func Tokenize(input string) (tokens []Token, errors []LexError)

// Streaming
func (l *Lexer) NextToken() Token
func (l *Lexer) Errors() []LexError

// Keyword classification (used by F4)
func KeywordToken(name string) (int, bool)
func IsReservedKeyword(name string) bool

// Debug
func TokenName(t int) string
```

No `baseOffset` field. No `PreserveTrivia` mode. No line/column tracking. All deferred unless a future consumer demands them.

## Testing

Test file layout matches the package layout. Run with `go test ./snowflake/parser/...` (scoped — never global).

### `lexer_test.go` — token recognition

Table-driven coverage of every token kind. Each row: `{name string, input string, want []Token}`. ~80 cases minimum. Categories:

- Single-char operators: `+`, `-`, `*`, `/`, `%`, `~`, `(`, `)`, `[`, `]`, `{`, `}`, `,`, `;`, `.`, `@`, `:`, `<`, `>`, `=`
- Multi-char operators: `::`, `||`, `->`, `->>`, `=>`, `!=`, `<>`, `<=`, `>=`
- Identifiers: unquoted (`abc`, `_x`, `x123`, `x@y`, `x$y`), quoted (`"my col"`, `"a""b"` doubled-quote), variables (`$x`, `$1`, `$ABC_123`)
- String literals: simple (`'hello'`), with escapes (`'a\nb'`), doubled-quote (`'a''b'`), dollar (`$$hello$$`), hex prefix (`X'48656C6C6F'`)
- Numeric literals: int (`42`), float (`1.5`, `1.`, `.5`), real (`1e10`, `1.5e-10`, `.5e+10`)
- Keywords: at least 5 from each of {DDL, DML, type names, functions, scripting} categories, plus 5 spot-checks of reserved-vs-non-reserved
- Comments: line `--`, line `//`, block `/* */`, nested block `/* /* */ */` — assert all are dropped from the token stream

### `position_test.go` — position math

For each token kind, assert `Loc.Start`/`Loc.End` are byte-correct. Inputs include:
- Token at start of input: `Loc.Start == 0`
- Token after whitespace: `Loc.Start == len(whitespace)`
- Token at end of input: `Loc.End == len(input)`
- Multi-byte tokens (`::`, `->>`, identifier `hello`): `Loc.End - Loc.Start == len(token)`
- String literals: `Loc.Start` is the opening quote, `Loc.End` is one past the closing quote
- EOF token: `Loc.Start == Loc.End == len(input)`

### `edge_cases_test.go` — corner cases

- EOF mid-string: `'unterminated` → 1 error, 1 `tokInvalid`, EOF
- EOF mid-block-comment: `/* unterminated` → 1 error, EOF (comment is dropped, error recorded)
- Nested block comment 3 deep: `/* a /* b /* c */ */ */` → no errors, no tokens emitted
- Mismatched depth: `/* /* */` → 1 unterminated error
- `1.` numeric edge: emits `tokFloat`
- `.1` numeric edge: emits `tokFloat`
- `1e10` numeric edge: emits `tokReal`
- `1e+10`, `1e-10`: emit `tokReal`
- Ambiguous operator chains: `a::b` → `[ident, ::, ident]`; `a->b` → `[ident, ->, ident]`; `a->>b` → `[ident, ->>, ident]`; `a||b` → `[ident, ||, ident]`; `a=>b` → `[ident, =>, ident]`; `a!=b` → `[ident, !=, ident]`; `a<=b`, `a>=b`, `a<>b` similar
- Unicode in quoted identifier: `"héllo"` → `tokQuotedIdent` with `Str = "héllo"`
- Unicode in string literal: `'café'` → `tokString` with `Str = "café"`
- Doubled-quote in quoted ident: `"a""b"` → `tokQuotedIdent` with `Str = "a\"b"`
- Doubled-quote in string: `'a''b'` → `tokString` with `Str = "a'b"`
- Empty quoted ident: `""` → `tokQuotedIdent` with `Str = ""`
- Bare `$`: emits `int('$')` single-char token

### `error_recovery_test.go` — multi-error recovery

Inputs with 2+ deliberate errors. Assert all errors are collected and the lexer reaches EOF without halting:

- `'unterm1' SELECT 'unterm2` → 2 errors, reaches EOF
- `/* a SELECT /* b` → 1 error, reaches EOF (nested unterminated)
- `\x00 SELECT \x01` → 2 invalid-byte errors, reaches EOF, SELECT token preserved between them

### `legacy_corpus_test.go` — regression baseline

Iterate every `.sql` file in `snowflake/parser/testdata/legacy/` (27 files). For each file, tokenize and assert:

1. `len(errors) == 0` — no lex errors
2. No `tokInvalid` emitted
3. The last non-EOF token's `Loc.End` is at or beyond the file's last semicolon byte position (sanity check that we made it through the whole file)

This is the regression baseline. Any failure here is a parser bug.

### `keyword_completeness_test.go` — drift detector

Read `/Users/h3n4l/OpenSource/parser/snowflake/SnowflakeLexer.g4`. For every line matching `^[A-Z_][A-Z0-9_]* : '[^']+';`, extract the token name and assert it has a corresponding entry in `keywordMap`. Read the legacy `build_id_contains_non_reserved_keywords.py`'s `snowflake_reserved_keyword` dict (parse via simple regex on the dict literal) and assert every entry is present in `keywordReserved`.

The drift detector path is hardcoded to `/Users/h3n4l/OpenSource/parser/snowflake/SnowflakeLexer.g4`. The test skips with `t.Skip` if the file is missing (so CI on machines without the legacy parser checkout doesn't fail).

## Out of Scope

F2 does NOT include:

| Feature | Where it lives |
|---|---|
| Statement splitter | F3 |
| Recursive-descent parser | F4 |
| AST construction | F4 |
| Concrete statement / expression / type nodes | T1.x, T2.x, T4.x, T5.x |
| Snowflake Scripting `:=` (verify in F4 brainstorming) | F4 or future amendment |
| Trivia / comment preservation | Future amendment |
| Line/column tracking | F4 (LineTable helper) |
| Cloud-storage URL token kinds (`S3_PATH`, etc.) | Detected at parser level via STRING content |
| Charset introducers (`_utf8'...'`) | Not in Snowflake |
| `baseOffset` substring positioning | F3 if needed |

## Acceptance Criteria

F2 is complete when:

1. `go build ./snowflake/parser/...` succeeds.
2. `go vet ./snowflake/parser/...` clean.
3. `gofmt -l snowflake/parser/` clean.
4. `go test ./snowflake/parser/...` passes — all token-recognition, position, edge-case, error-recovery, legacy-corpus, and keyword-completeness tests pass.
5. All 27 files in `snowflake/parser/testdata/legacy/` tokenize without errors and without emitting any `tokInvalid`.
6. The keyword completeness test passes — every `[A-Z_]+ : '...';` rule in the legacy `.g4` has a `kwXxx` entry, and every reserved keyword from the Python script is marked reserved.
7. `go build ./snowflake/...` (isolation check) — F2 builds without depending on any DAG node beyond F1.
8. After merge, `docs/migration/snowflake/dag.md` F2 status is flipped to `done`.

## Files Created

```
snowflake/parser/lexer.go
snowflake/parser/tokens.go
snowflake/parser/keywords.go
snowflake/parser/errors.go
snowflake/parser/lexer_test.go
snowflake/parser/position_test.go
snowflake/parser/edge_cases_test.go
snowflake/parser/error_recovery_test.go
snowflake/parser/legacy_corpus_test.go
snowflake/parser/keyword_completeness_test.go
docs/superpowers/specs/2026-04-07-snowflake-lexer-design.md  (this file)
```

Estimated total: ~3,500 LOC across 10 code files plus this design document.
