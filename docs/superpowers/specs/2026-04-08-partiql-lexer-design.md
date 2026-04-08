# PartiQL `lexer` Design Spec

**DAG node:** `lexer` (node 2 in `docs/migration/partiql/dag.md`)
**Priority:** P0 (foundational — `parser-foundation` and onwards depend on this)
**Branch:** `feat/partiql/lexer`
**Worktree:** `/Users/h3n4l/OpenSource/omni/.worktrees/feat-partiql-lexer`
**Linear umbrella:** BYT-9000

## Goal

Create the `partiql/parser` package with a hand-written tokenizer for PartiQL source code. This is the second node in the PartiQL migration DAG (after `ast-core`, which has been merged in PR bytebase/omni#13). The lexer is the foundation that `parser-foundation` (DAG node 4) and every subsequent parser DAG node will consume.

The scope is **full coverage of `bytebase/parser/partiql/PartiQLLexer.g4`** for the standard SQL token set plus PartiQL extensions, with two intentional deferrals:

- **Ion-mode body parsing** is deferred to DAG node 17 (`parser-ion-literals`). The base lexer recognizes a backtick-delimited region as a single `tokION_LITERAL` via a naive byte-to-byte scan; the refined Ion-mode-aware scanner lands later.
- **Date/time literal special-casing** is deferred to DAG node 18 (`parser-datetime-literals`). The base lexer emits the keyword tokens `tokDATE`, `tokTIME`, `tokTIMESTAMP`; the parser composes them into `ast.DateLit` / `ast.TimeLit` at parse time.

## Inputs

- **Authoritative grammar:** `bytebase/parser/partiql/PartiQLLexer.g4` (542 lines, ~200 keywords + ~28 operators + 5 literals/identifiers + Ion-mode rules)
- **Reference precedent:** `cosmosdb/parser/lexer.go` (most recent hand-written lexer in omni; closest model)
- **Looser precedent:** `mongo/parser/lexer.go` (compared but not adopted)
- **AST package:** `partiql/ast` (already merged) — provides `ast.Loc{Start, End}` for token position tracking
- **Test corpus:** `partiql/parser/testdata/aws-corpus/` (63 PartiQL examples crawled from the AWS DynamoDB PartiQL reference; already committed)

## Architecture decisions

### D1. Single-pass hand-written lexer following `cosmosdb/parser/lexer.go`

Approach: hand-written byte-at-a-time tokenizer with two-character lookahead for compound operators. The cosmosdb lexer is the precedent set by the most recent NoSQL engine in omni; the PartiQL lexer follows it almost verbatim with PartiQL-specific token types and keywords.

**Rationale:** the omni project already has two hand-written lexers (cosmosdb and mongo); a third one in the same shape minimizes the cognitive cost of cross-engine work. Generated lexers (ANTLR, goyacc, ragel) introduce build-time dependencies and obscure debug paths.

### D2. Token positions use `ast.Loc` directly

```go
type Token struct {
    Type int      // tok* constant from token.go
    Str  string   // raw source text (or decoded value for string literals)
    Loc  ast.Loc  // half-open byte range
}
```

The cosmosdb lexer uses flat `Loc int, End int`; the parser then converts to `ast.Loc{Start, End}` at every AST node construction site. PartiQL's lexer embeds `ast.Loc` directly to eliminate that boilerplate at the parser-foundation boundary. The `partiql/parser → partiql/ast` import is natural — the parser builds AST nodes and has to import `ast` anyway.

### D3. First-error-and-stop error model

Identical to cosmosdb. When the scanner encounters an unrecoverable lex error (unterminated string, unexpected character, unterminated block comment, unterminated Ion literal):

1. Set `Lexer.Err` to the first error encountered (subsequent errors are dropped)
2. Return `Token{Type: tokEOF, ...}` from this `Next()` call
3. All subsequent `Next()` calls return `tokEOF` unconditionally

The parser checks `l.Err` after consuming the token stream. Error recovery happens at the parser level, not the lexer level. Keeps the lexer simple and predictable.

### D4. Three-file split: `token.go` + `keywords.go` + `lexer.go`

```
partiql/parser/
├── token.go        # Token struct, ~227 token-type constants in iota groups, tokenName(int)
├── keywords.go     # ~200-entry keyword map (lowercase → token type)
├── lexer.go        # Lexer struct, NewLexer, Next(), scan helpers
└── lexer_test.go   # 3 test functions: TestLexer_Tokens, TestLexer_AWSCorpus, TestLexer_Errors
```

The DAG entry for this node listed `lexer.go, token.go` (2 files), but PartiQL's keyword count is ~7× cosmosdb's. With ~200 keywords, the keyword map alone is ~250 lines and would dominate `lexer.go`, burying the scan logic. Splitting `keywords.go` out keeps each file under ~400 lines and lets readers focus on data-vs-logic separately. The DAG was speculative; a 3-file split is the natural response to the keyword count.

### D5. Unexported token constants (`tok*`)

The lexer is internal to the `partiql/parser` package. The recursive-descent parser-foundation node will live in the same package and consume `tok*` constants directly. External consumers go through the future public `partiql.Parse()` API at `partiql/parse.go`, which returns `[]ast.StmtNode` — no token types crossed the package boundary.

Naming convention: `tokSELECT`, `tokFROM`, `tokANGLE_DOUBLE_LEFT`, etc. — lowercase `tok` prefix + UPPERCASE name from the grammar. Compound operator names follow `PartiQLLexer.g4` rule names verbatim for traceability.

### D6. Naive backtick scan for `tokION_LITERAL`

The `scanIonLiteral` helper does a byte-to-byte scan from `` ` `` to the next `` ` ``. The captured `Token.Str` is the raw inner content (no decoding); `Token.Loc` covers the entire `` `…` `` range including both backticks.

**Known limitation, documented inline and deferred to DAG node 17:** if the Ion content contains a backtick inside a `'…'` quoted symbol or `"…"` string (which Ion technically allows), this naive scan closes the literal prematurely. The full Ion-mode-aware implementation is in scope for DAG node 17 (`parser-ion-literals`).

**Why this is safe for the base lexer:** the AWS DynamoDB PartiQL corpus has zero real Ion literals. The only 2 backtick uses are in `select-001.partiql` and `insert-002.partiql` (syntax skeletons with placeholder backticks like `` `table` ``), and both are filtered out of the corpus smoke test by hard-coded skip list.

## Token type taxonomy

### File: `token.go`

```go
package parser

import "github.com/bytebase/omni/partiql/ast"

// Token is a single PartiQL lexer token.
type Token struct {
    Type int     // tok* constant below
    Str  string  // raw source text for most tokens; decoded value for string/quoted-identifier literals
    Loc  ast.Loc // half-open byte range
}

// ===== Special tokens =====
const (
    tokEOF     = 0  // end of input or after lex error
    tokInvalid = 1  // sentinel; never returned (Lexer.Err is the error channel)
)

// ===== Literal tokens — group 1000 =====
const (
    tokSCONST       = iota + 1000  // single-quoted string literal: 'hello'
    tokICONST                       // integer literal: 42
    tokFCONST                       // decimal/float literal: 3.14, 1e10, .5
    tokIDENT                        // unquoted identifier (case-insensitive lookup)
    tokIDENT_QUOTED                 // double-quoted identifier (case-sensitive): "Foo"
    tokION_LITERAL                  // backtick-delimited Ion blob (body deferred to node 17)
)

// ===== Operator and punctuation tokens — group 2000 (~28 entries) =====
const (
    tokPLUS = iota + 2000      // +
    tokMINUS                    // -
    tokASTERISK                 // *
    tokSLASH_FORWARD            // /
    tokPERCENT                  // %
    tokCARET                    // ^
    tokTILDE                    // ~
    tokAT_SIGN                  // @
    tokEQ                       // =
    tokNEQ                      // <> or !=
    tokLT                       // <
    tokGT                       // >
    tokLT_EQ                    // <=
    tokGT_EQ                    // >=
    tokCONCAT                   // ||
    tokANGLE_DOUBLE_LEFT        // <<  (PartiQL bag-literal start)
    tokANGLE_DOUBLE_RIGHT       // >>  (PartiQL bag-literal end)
    tokPAREN_LEFT               // (
    tokPAREN_RIGHT              // )
    tokBRACKET_LEFT             // [
    tokBRACKET_RIGHT            // ]
    tokBRACE_LEFT               // {
    tokBRACE_RIGHT              // }
    tokCOLON                    // :
    tokCOLON_SEMI               // ;
    tokCOMMA                    // ,
    tokPERIOD                   // .
    tokQUESTION_MARK            // ?
)

// ===== Keyword tokens — group 3000 (~200 entries, alphabetical) =====
const (
    tokABSOLUTE = iota + 3000
    tokACTION
    tokADD
    tokALL
    tokALLOCATE
    // ... 195+ more, one per uppercase keyword rule in PartiQLLexer.g4 lines 13–295
    // Including window keywords (LAG, LEAD, OVER, PARTITION),
    // PartiQL extension keywords (CAN_CAST, CAN_LOSSLESS_CAST, MISSING, PIVOT, UNPIVOT,
    // LIMIT, OFFSET, REMOVE, INDEX, LET, CONFLICT, DO, RETURNING, MODIFIED, NEW, OLD, NOTHING),
    // and data type keywords (TUPLE, INT2/4/8, INTEGER2/4/8, BIGINT, BOOL, BOOLEAN, STRING,
    // SYMBOL, CLOB, BLOB, STRUCT, LIST, SEXP, BAG)
)

// tokenName returns the canonical printable name for a token type constant.
// Used by error messages, test failure output, and future debugging.
func tokenName(t int) string {
    switch t {
    case tokEOF:
        return "EOF"
    case tokInvalid:
        return "INVALID"
    case tokSCONST:
        return "SCONST"
    // ... one arm per constant, ~227 cases total
    }
}
```

**Total token constants:** approximately 227 (2 specials + 6 literals + 28 operators + ~200 keywords). Exact count is determined during implementation by enumerating `PartiQLLexer.g4` keyword rules; documented in the implementation plan.

### File: `keywords.go`

```go
package parser

// keywords maps lowercase keyword strings to their tok* constants.
// PartiQL keywords are case-insensitive per PartiQLLexer.g4
// `options { caseInsensitive = true; }`. The lexer lowercases identifiers
// before lookup.
//
// Source: every uppercase rule in PartiQLLexer.g4 from line 13 (ABSOLUTE)
// through line 295 (BAG), plus the window/other keywords (lines 245–270),
// plus the data type keywords (lines 277–295).
var keywords = map[string]int{
    "absolute": tokABSOLUTE,
    "action":   tokACTION,
    "add":      tokADD,
    "all":      tokALL,
    // ... ~196 more entries, alphabetical
}
```

**Acceptance constraint:** `len(keywords)` must equal the count of `tok*` keyword constants in the group-3000 block. Verified by a one-line test: `if len(keywords) != expectedKeywordCount { t.Errorf(...) }`.

## Scanning behavior

### `Lexer` struct + public API

```go
// Lexer is a hand-written tokenizer for PartiQL source code.
type Lexer struct {
    input string // source text
    pos   int    // current read position (next byte to consume)
    start int    // byte offset of token currently being scanned
    Err   error  // first error encountered, nil if none
}

// NewLexer creates a Lexer for the given source string.
func NewLexer(input string) *Lexer {
    return &Lexer{input: input}
}

// Next returns the next token. At end of input or after a lex error,
// returns Token{Type: tokEOF, ...}. After Err is set, all subsequent
// calls return tokEOF.
func (l *Lexer) Next() Token { ... }
```

### `Next()` dispatch

```go
func (l *Lexer) Next() Token {
    if l.Err != nil {
        return Token{Type: tokEOF, Loc: ast.Loc{Start: l.pos, End: l.pos}}
    }
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
    case ch == '`':
        return l.scanIonLiteral()
    case ch >= '0' && ch <= '9':
        return l.scanNumber()
    case ch == '.' && l.pos+1 < len(l.input) && isDigit(l.input[l.pos+1]):
        return l.scanNumber() // .5 is a decimal literal per LITERAL_DECIMAL alt 2
    case isIdentStart(ch):
        return l.scanIdentOrKeyword()
    default:
        return l.scanOperator()
    }
}
```

Each scan helper sets `Token.Loc = ast.Loc{Start: l.start, End: l.pos}` after consuming.

### `skipWhitespaceAndComments`

Three skip targets (all on the HIDDEN channel per the grammar):

- **Whitespace:** `' '`, `'\t'`, `'\n'`, `'\r'`
- **Line comment:** `--` to end of line (`COMMENT_SINGLE` rule, terminates on `\n` or EOF)
- **Block comment:** `/* */` non-nested, greedy-shortest (`COMMENT_BLOCK` rule). Unterminated block comment sets `l.Err`.

Loops until none of the three match.

### `scanString` — single-quoted string literals

Grammar: `LITERAL_STRING : '\'' ( ('\'\'') | ~('\'') )* '\'';`

1. Skip opening `'`
2. Walk forward; on `''` (doubled single quote), append a single `'` to the value buffer and skip 2 bytes
3. On a single `'`, terminate — token is `tokSCONST` with `Token.Str` = decoded value
4. On any other byte, append to the value buffer and advance
5. On EOF before closing `'`, set `l.Err` to `"unterminated string literal at position N"`, return `tokEOF`

`Token.Str` holds the **decoded** value (with `''` collapsed to `'`). **No backslash escapes** — PartiQL grammar has none.

### `scanQuotedIdent` — double-quoted identifiers

Grammar: `IDENTIFIER_QUOTED : '"' ( ('""') | ~('"') )* '"';`

Same algorithm as `scanString` with `"` as the delimiter. Returns `tokIDENT_QUOTED`. `Token.Str` holds the decoded value (with `""` collapsed to `"`).

**Quoted identifiers are case-preserving** — `"Foo"` is distinct from `"foo"`. Unquoted identifiers are case-insensitive (the keyword map lookup is lowercased).

### `scanIdentOrKeyword`

1. Walk forward while `isIdentContinue(ch)` (letters, digits, `_`, `$`)
2. Take the source slice `raw := l.input[l.start:l.pos]`
3. Lowercase via `strings.ToLower(raw)` — PartiQL keywords are case-insensitive
4. Look up in `keywords` map; if found, return that keyword token with `Token.Str = raw` (preserving original case)
5. Otherwise return `tokIDENT` with `Token.Str = raw`

Helper functions:
- `isIdentStart(ch)`: `[a-zA-Z_$]` — letter, underscore, or dollar
- `isIdentContinue(ch)`: `isIdentStart(ch) || isDigit(ch)`

PartiQL allows `$` in identifiers (consistent with cosmosdb and standard SQL dialects).

### `scanNumber` — integer and decimal literals

Grammar:
```
LITERAL_INTEGER : DIGIT+;
LITERAL_DECIMAL :
    DIGIT+ '.' DIGIT* ([e] [+-]? DIGIT+)?
  | '.' DIGIT+ ([e] [+-]? DIGIT+)?
  | DIGIT+ ([e] [+-]? DIGIT+)?
  ;
```

(`[e]` is case-insensitive per the grammar option — both `e` and `E`.)

Algorithm:
1. `isFloat := false`
2. If first char is `.`, consume it, set `isFloat = true`, then consume one or more digits
3. Otherwise consume one or more digits, then optionally `.` followed by zero or more digits (set `isFloat = true` if `.` present)
4. Optionally consume `[eE]`, then optional `[+-]`, then one or more digits (set `isFloat = true`)
5. Return `tokFCONST` if `isFloat`, else `tokICONST`. `Token.Str` is the raw source slice.

**No support for `0x`/`0X` hex literals** — PartiQL grammar has no rule for them.

### `scanOperator` — one and two-character operators

Two-character lookahead first; fall through to single-character. Required two-char operators: `<=`, `>=`, `<>`, `<<`, `>>`, `||`, `!=`.

```go
ch := l.input[l.pos]
l.pos++
if l.pos < len(l.input) {
    next := l.input[l.pos]
    switch {
    case ch == '<' && next == '=':  l.pos++; return tok(tokLT_EQ, "<=")
    case ch == '<' && next == '>':  l.pos++; return tok(tokNEQ, "<>")
    case ch == '<' && next == '<':  l.pos++; return tok(tokANGLE_DOUBLE_LEFT, "<<")
    case ch == '>' && next == '=':  l.pos++; return tok(tokGT_EQ, ">=")
    case ch == '>' && next == '>':  l.pos++; return tok(tokANGLE_DOUBLE_RIGHT, ">>")
    case ch == '|' && next == '|':  l.pos++; return tok(tokCONCAT, "||")
    case ch == '!' && next == '=':  l.pos++; return tok(tokNEQ, "!=")
    }
}
// Fall through to single-char switch
switch ch {
case '+': return tok(tokPLUS, "+")
case '-': return tok(tokMINUS, "-")
// ... 26 more single-char cases (asterisk, slash, percent, caret, tilde, at_sign,
// eq, lt, gt, paren_left/right, bracket_left/right, brace_left/right, colon,
// colon_semi, comma, period, question_mark)
default:
    l.Err = fmt.Errorf("unexpected character %q at position %d", ch, l.start)
    return Token{Type: tokEOF, Loc: ast.Loc{Start: l.start, End: l.start}}
}
```

**Note:** bare `!` (without trailing `=`) and bare `|` (without trailing `|`) trigger the default arm. PartiQL only uses these characters as parts of `!=` and `||` respectively.

### `scanIonLiteral` — backtick-delimited Ion blob

Naive base behavior, refined later by DAG node 17:

1. Skip opening `` ` ``
2. Walk forward byte by byte until the next `` ` `` is found
3. On the closing `` ` ``, advance past it and return `tokION_LITERAL` with `Token.Str` = raw inner content (no decoding) and `Token.Loc` covering the entire `` `…` `` range
4. On EOF before closing `` ` ``, set `l.Err` to `"unterminated Ion literal at position N"`, return `tokEOF`

**Inline doc comment** must spell out the known limitation: backticks inside Ion strings/symbols will close the literal prematurely. Refinement deferred to DAG node 17 (`parser-ion-literals`).

### Error reporting

First-error-and-stop. Five error triggers:

| Trigger | Error message | Location |
|---------|--------------|----------|
| Unrecognized character in `scanOperator` default arm | `"unexpected character %q at position %d"` | `scanOperator` |
| Unterminated `'…'` string | `"unterminated string literal at position %d"` | `scanString` |
| Unterminated `"…"` quoted identifier | `"unterminated quoted identifier at position %d"` | `scanQuotedIdent` |
| Unterminated `` `…` `` Ion literal | `"unterminated Ion literal at position %d"` | `scanIonLiteral` |
| Unterminated `/* … */` block comment | `"unterminated block comment at position %d"` | `skipWhitespaceAndComments` |

Position numbers are byte offsets from the start of input.

### Position tracking

- `l.start` is set to `l.pos` at the beginning of every token scan (after `skipWhitespaceAndComments`)
- `l.pos` advances byte by byte during scanning
- `Token.Loc = ast.Loc{Start: l.start, End: l.pos}` after the scan completes
- Whitespace and comments do **not** appear in tokens; they're consumed silently

## Test plan

All tests in a single `partiql/parser/lexer_test.go` file. Three test functions plus one helper test.

### 1. `TestLexer_Tokens` — token-stream golden tests

Table-driven, ~55–60 cases. Each case provides an input string and the expected token stream (excluding the trailing `tokEOF`):

```go
type tokenStreamCase struct {
    name   string
    input  string
    tokens []Token // expected, excluding trailing EOF
}
```

**Test categories** (the implementation plan enumerates the full case list; representative samples below):

- **Empty / whitespace / comments** (4 cases): empty string, whitespace-only, line-comment-only, block-comment-only — all produce empty token streams
- **Identifiers** (5 cases): `foo`, `_foo`, `x$y`, `"Foo"` (case-preserved), `"a""b"` (doubled-quote escape)
- **Keywords** (3+ cases): `select` / `SELECT` / `Select` (case-insensitive lookup, raw text preserved), plus a few from each category (DML, DDL, types, window, PartiQL extensions)
- **Number literals** (6 cases): integer `42`, decimal `3.14`, leading-dot `.5`, scientific `1e10`, uppercase exponent `1E10`, negative exponent `1.5e-3`
- **String literals** (4 cases): `'hello'`, `''` (empty), `'it''s'` (doubled quote → `it's`), `'  '` (whitespace preserved)
- **Single-char operators** (~28 cases): one per single-char op
- **Two-char operators** (7 cases): `<=`, `>=`, `<>`, `!=`, `<<`, `>>`, `||`
- **Ion literals** (2 cases): `` `{a: 1}` ``, `` `` `` (empty)
- **Multi-token sequences** (3+ cases): `SELECT * FROM Music`, `WHERE Artist='Pink Floyd'`, `<<1, 2, 3>>`
- **Position tracking with comments + whitespace** (2+ cases): identifier after line comment, identifier after block comment

Each row asserts `reflect.DeepEqual(got, want)` on the full token slice and `l.Err == nil`.

### 2. `TestLexer_AWSCorpus` — corpus smoke test

Loads every `.partiql` file from `partiql/parser/testdata/aws-corpus/`, filters out `select-001.partiql` and `insert-002.partiql` (syntax skeletons with bracket/backtick placeholders, per `index.json`), runs the lexer, and asserts:

- `l.Err == nil` after draining the token stream
- The token count is `> 0`
- The token stream terminates within 100,000 tokens (sanity bound; corpus files are short)

**Expected outcome:** 61 sub-tests pass (63 corpus files − 2 filtered skeletons). The skip list is hard-coded for clarity; `index.json` is not parsed at runtime.

### 3. `TestLexer_Errors` — error-case golden tests

Five error triggers, one test case each:

| Test name | Input | Expected error substring |
|-----------|-------|--------------------------|
| `unterminated_string` | `'hello` | `"unterminated string literal"` |
| `unterminated_quoted_ident` | `"foo` | `"unterminated quoted identifier"` |
| `unterminated_ion_literal` | `` `abc `` | `"unterminated Ion literal"` |
| `unterminated_block_comment` | `/* nope` | `"unterminated block comment"` |
| `lone_bang_then_space` | `! 1` | `"unexpected character"` |

Each case drains `Next()` until `tokEOF`, then asserts `l.Err != nil` and `strings.Contains(l.Err.Error(), wantSubstring)`.

### 4. `TestTokenName_AllCovered` — exhaustive `tokenName` coverage

Walks every `tok*` constant declared in `token.go` (via a hard-coded list in the test) and asserts `tokenName(t)` returns a non-empty, non-`"INVALID"` string. If a future contributor adds a new `tok*` constant without wiring it into `tokenName`, this test fails.

```go
func TestTokenName_AllCovered(t *testing.T) {
    all := []int{
        tokEOF, tokInvalid,
        tokSCONST, tokICONST, tokFCONST, tokIDENT, tokIDENT_QUOTED, tokION_LITERAL,
        tokPLUS, tokMINUS, /* ... 26 more operators */,
        tokABSOLUTE, tokACTION, /* ... ~198 more keywords */,
    }
    for _, tt := range all {
        if name := tokenName(tt); name == "" || name == "INVALID" {
            t.Errorf("tokenName(%d) returned %q — missing switch arm?", tt, name)
        }
    }
}
```

### 5. `TestKeywords_LenMatchesConstants` — keyword count parity

```go
func TestKeywords_LenMatchesConstants(t *testing.T) {
    const expectedKeywordCount = 200 // exact count, set during implementation
    if got := len(keywords); got != expectedKeywordCount {
        t.Errorf("len(keywords) = %d, want %d — did a tok* keyword constant get added without a map entry?", got, expectedKeywordCount)
    }
}
```

The exact `expectedKeywordCount` is set during implementation by counting `tok*` keyword constants and the corresponding `keywords` map entries; both must agree.

### Test scope

- Run with `go test ./partiql/parser/...` (no global test runs, per the implementing skill)
- Zero external dependencies — standard library only (`testing`, `os`, `path/filepath`, `reflect`, `strings`)
- Total expected sub-tests: ~55 (TestLexer_Tokens) + 61 (TestLexer_AWSCorpus) + 5 (TestLexer_Errors) + 1 (TestTokenName_AllCovered) + 1 (TestKeywords_LenMatchesConstants) = **~123 sub-tests across 5 test functions**

## Acceptance criteria

The `lexer` DAG node is **done** when:

1. Three source files exist: `partiql/parser/token.go`, `partiql/parser/keywords.go`, `partiql/parser/lexer.go` (plus `partiql/parser/lexer_test.go`)
2. ~227 unexported `tok*` constants defined across the 4 iota groups (specials at 0/1, literals at 1000+, operators at 2000+, keywords at 3000+) — exact count documented in the implementation plan and asserted by `TestKeywords_LenMatchesConstants`
3. `tokenName(int) string` returns a non-empty canonical name for every constant — verified by `TestTokenName_AllCovered`
4. `keywords` map has one entry per `tok*` keyword constant — verified by `TestKeywords_LenMatchesConstants`
5. `Lexer.Next()` produces correct tokens for every test case in `TestLexer_Tokens` (~55 cases)
6. `TestLexer_AWSCorpus` passes (all 61 lexable files produce non-error token streams ending in EOF)
7. `TestLexer_Errors` passes (all 5 error triggers fire with the expected message substring)
8. `go test ./partiql/parser/...` passes
9. `go vet ./partiql/parser/...` clean
10. `gofmt -l partiql/parser/` clean
11. **Grammar cross-check:** every uppercase token rule in `bytebase/parser/partiql/PartiQLLexer.g4` lines 13–356 (excluding fragments and Ion-mode rules) maps to exactly one `tok*` constant. Verified during implementation by reading the `.g4` file line by line.

## Non-goals

- **No Ion-mode-aware backtick scanning.** The base `scanIonLiteral` does a naive `` ` `` to `` ` `` byte scan. Refinement to handle backticks-inside-Ion-strings is deferred to DAG node 17 (`parser-ion-literals`). Documented inline in `scanIonLiteral`.
- **No date/time literal special-casing.** The grammar treats `DATE 'YYYY-MM-DD'` as the keyword `DATE` followed by a string literal — the parser composes them into `ast.DateLit`. Same for `TIME` and `TIMESTAMP`. The lexer just emits the keyword tokens. Date/time literal recognition at parse time is deferred to DAG node 18 (`parser-datetime-literals`).
- **No comment preservation.** Comments are silently skipped (HIDDEN channel per grammar). If a future task needs round-trip-preserving formatter output, the lexer can be extended; for parser-foundation onwards, comments are noise.
- **No `Lexer.Tokenize() []Token` API.** The recursive-descent parser pulls tokens via `Next()`; no tokenize-all-at-once API is needed.
- **No position-to-line/column conversion.** The `ast.Loc` is byte-offset only; line/column conversion happens at the public `partiql/parse.go` boundary (a future task).
- **No cross-comparison test against the legacy ANTLR lexer.** Overkill for a P0 lexer; the golden tests + corpus smoke + grammar cross-check are sufficient. Could be added later if regressions appear.
- **No hex/binary/octal number literals.** PartiQL grammar has no rule for them.
- **No backslash escapes in string literals.** PartiQL grammar uses only `''` doubling for embedded single quotes. No `\n`, `\t`, `\u0000`, etc.
- **No semicolon-as-statement-terminator handling.** Statement splitting on `;` is the splitter's job (a future task), not the lexer's. The lexer just emits `tokCOLON_SEMI` and lets the parser handle it.
- **No utf8.RuneError handling.** The lexer operates on bytes, not runes. UTF-8-encoded source is supported because identifier rules use ASCII-only character classes; non-ASCII bytes inside string literals or quoted identifiers are passed through verbatim.

## Risks & mitigations

| Risk | Mitigation |
|------|-----------|
| Missing a keyword from the grammar | Acceptance criterion 11 forces a manual line-by-line cross-check; `TestKeywords_LenMatchesConstants` enforces parity between the constants block and the map |
| Naive Ion scan misparses real Ion content | Documented as known limitation in `scanIonLiteral` doc comment; AWS corpus has zero real Ion examples (the only 2 backtick uses are in skeleton files filtered out); refinement deferred to DAG node 17 |
| Two-character operator lookahead order matters (e.g., `<` vs `<<` vs `<=` vs `<>`) | Golden tests cover all 7 two-char operators explicitly; the lookahead order is documented in `scanOperator` |
| Position tracking off-by-one errors | Every golden test row includes `ast.Loc{Start, End}` — any drift fails the test |
| Identifier `$` start vs `$` mid-name | Test case `unquoted_ident_with_dollar` (`x$y`) covers mid-name; the grammar permits `[A-Z$_]` as start so `$y` is also a valid identifier (covered by spec but not tested explicitly) |
| `len(keywords)` and `tok*` constants drifting apart | `TestKeywords_LenMatchesConstants` fails the build if they diverge |
| `tokenName` arms going stale when new constants are added | `TestTokenName_AllCovered` walks every constant explicitly |

## References

- `bytebase/parser/partiql/PartiQLLexer.g4` — authoritative grammar (lines 13–356 for the standard mode; lines 406–541 for Ion mode, deferred to node 17)
- `cosmosdb/parser/lexer.go` — closest precedent in omni
- `mongo/parser/lexer.go` — looser precedent (compared but not adopted)
- `partiql/ast/node.go` — `Loc` struct definition
- `docs/migration/partiql/analysis.md` — full grammar coverage
- `docs/migration/partiql/dag.md` — migration node ordering
- `partiql/parser/testdata/aws-corpus/index.json` — corpus manifest with skeleton-file flags
