# Doris Lexer Design

DAG node: F2 (lexer)

## Goal

Implement the lexer for `doris/parser` -- tokenize Doris SQL input into a stream of tokens that the parser (F4), statement splitter (F3), and all downstream consumers depend on. Must handle 522 keywords, MySQL-compatible syntax (backtick identifiers, `<=>` null-safe equal, conditional comments), and Doris-specific hint syntax.

## Scope

### In scope

- Token struct and TokenKind constants (EOF, Invalid, literals, operators, keywords)
- 522 keyword constants with reserved/non-reserved classification
- Lexer with character-by-character scanning
- String literals (single-quoted, double-quoted, backtick-quoted)
- Numeric literals (integer, decimal, float/exponent, hex `0x`, binary `0b`, suffixed `L`/`S`/`Y`/`BD`)
- All operators and punctuation from DorisLexer.g4
- Line comments (`--`, `//`), block comments (`/* */`)
- Hint tokens (`/*+ ... */`)
- Error recovery (unterminated strings/comments)
- Unit tests

### Out of scope

- Statement splitter (F3)
- Parser entry / lookahead (F4)
- AST node creation
- Conditional comments (`/*!NNNNN ... */`) -- not used by bytebase for Doris

## Reference

- Template: `snowflake/parser/lexer.go` (17KB), `snowflake/parser/tokens.go` (16KB), `snowflake/parser/keywords.go` (59KB)
- MySQL compat: `tidb/parser/lexer.go` (56KB) for backtick/operator patterns
- Grammar source: `/Users/h3n4l/OpenSource/parser/doris/DorisLexer.g4` (522 keywords, 709 lines)

## File Structure

| File | Responsibility | ~LOC |
|------|----------------|------|
| `doris/parser/tokens.go` | Token struct, TokenKind constants, LexError type | ~150 |
| `doris/parser/keywords.go` | 522 keyword constants, keyword map, reserved classification, lookup/name functions | ~2000 |
| `doris/parser/lexer.go` | Lexer struct, NewLexer, NextToken, all scanning functions | ~600 |
| `doris/parser/lexer_test.go` | Unit tests for all token categories | ~500 |

## Package: `doris/parser/`

### `tokens.go`

**Token struct:**

```go
type Token struct {
    Kind TokenKind   // token type identifier
    Str  string      // content for identifiers, strings, hints
    Ival int64       // parsed value for tokInt
    Loc  ast.Loc     // byte offset span {Start, End}
}
```

**TokenKind ranges** (following snowflake/tidb convention):

```go
type TokenKind int

// Sentinel tokens
tokEOF     = 0
tokInvalid = 1

// Literal tokens (600-699)
tokInt         = 600  // integer: 42, 0xFF, 0b101
tokFloat       = 601  // decimal/exponent: 3.14, 1e10
tokString      = 602  // single or double quoted string
tokIdent       = 603  // unquoted identifier
tokQuotedIdent = 604  // backtick-quoted identifier
tokHexLiteral  = 605  // X'...' or x'...'
tokBitLiteral  = 606  // B'...' or b'...'
tokPlaceholder = 607  // ?

// Multi-character operators (500-599)
tokLessEq      = 500  // <=
tokGreaterEq   = 501  // >=
tokNotEq       = 502  // <> or !=
tokNullSafeEq  = 503  // <=>
tokLogicalAnd  = 504  // &&
tokDoublePipes = 505  // ||
tokShiftLeft   = 506  // <<
tokShiftRight  = 507  // >>
tokArrow       = 508  // ->
tokDoubleAt    = 509  // @@
tokAssign      = 510  // :=
tokDotDotDot   = 511  // ...
tokHintStart   = 512  // /*+
tokHintEnd     = 513  // */

// Keywords start at 700 (defined in keywords.go)

// Single-character operators/punctuation use their ASCII byte value:
// '(' = 40, ')' = 41, '*' = 42, '+' = 43, ',' = 44, '-' = 45,
// '.' = 46, '/' = 47, ':' = 58, ';' = 59, '<' = 60, '=' = 61,
// '>' = 62, '?' = 63, '@' = 64, '[' = 91, ']' = 93, '^' = 94,
// '{' = 123, '|' = 124, '}' = 125, '~' = 126, '%' = 37, '!' = 33,
// '&' = 38
```

**LexError type:**

```go
type LexError struct {
    Msg string
    Loc ast.Loc
}
```

### `keywords.go`

**Keyword constants** (522 total, starting at 700):

```go
const (
    kwACCOUNT_LOCK TokenKind = 700 + iota
    kwACCOUNT_UNLOCK
    kwACTIONS
    kwADD
    // ... (522 constants in alphabetical order)
    kwYEAR
)
```

**Keyword map** (case-insensitive lookup):

```go
var keywordMap = map[string]TokenKind{
    "account_lock":    kwACCOUNT_LOCK,
    "account_unlock":  kwACCOUNT_UNLOCK,
    // ... all 522 entries, lowercase keys
}
```

**Non-reserved set** (172 keywords that can be used as identifiers):

```go
var nonReservedKeywords = map[TokenKind]bool{
    kwACTIONS:        true,
    kwAFTER:          true,
    // ... 172 entries from DorisParser.g4 nonReserved rule
}
```

**Public functions:**

```go
// KeywordToken returns the keyword TokenKind for a name, or (0, false)
// if it's not a keyword. Lookup is case-insensitive.
func KeywordToken(name string) (TokenKind, bool)

// IsReserved reports whether kind is a reserved keyword (cannot be
// used as an unquoted identifier).
func IsReserved(kind TokenKind) bool

// TokenName returns a human-readable name for a TokenKind. Used for
// error messages and debugging.
func TokenName(kind TokenKind) string
```

### `lexer.go`

**Lexer struct:**

```go
type Lexer struct {
    input      string
    pos        int        // current byte position
    start      int        // start of current token
    errors     []LexError
    baseOffset int        // added to all Loc positions on return
}
```

**Public API:**

```go
// NewLexer creates a lexer for the given input.
func NewLexer(input string) *Lexer

// NewLexerWithOffset creates a lexer with a base offset added to all
// token positions. Used for multi-statement parsing where each statement
// is lexed independently but positions are relative to the full input.
func NewLexerWithOffset(input string, baseOffset int) *Lexer

// NextToken returns the next token from the input. Returns tokEOF when
// the input is exhausted. Errors (unterminated strings, unknown chars)
// produce tokInvalid tokens and are accumulated in Errors().
func (l *Lexer) NextToken() Token

// Tokenize is a convenience function that tokenizes the entire input
// and returns all tokens (including tokEOF) and any errors.
func Tokenize(input string) ([]Token, []LexError)

// Errors returns all lexing errors accumulated during scanning.
func (l *Lexer) Errors() []LexError
```

**Character dispatch** in `nextToken()`:

| First byte(s) | Action |
|---------------|--------|
| `' '`, `\t`, `\n`, `\r` | Skip whitespace |
| `'` or `"` | Scan string literal (backslash escapes + doubled-quote escapes) |
| `` ` `` | Scan backtick-quoted identifier (doubled-backtick escape) |
| `[0-9]` | Scan numeric (int/float/hex/binary/exponent/suffixed) |
| `.` + digit | Scan decimal literal (`.5`); otherwise return `'.'` |
| `[A-Za-z_$]` or non-ASCII | Scan identifier, check keyword map |
| `X/x` + `'` | Scan hex literal `X'...'` |
| `B/b` + `'` | Scan bit literal `B'...'` |
| `-` | `--` + space: line comment (skip); `->`: tokArrow; else `'-'` |
| `/` | `/*+`: tokHintStart; `/*`: block comment (skip); `//`: line comment (skip); else `'/'` |
| `*` | `*/`: tokHintEnd; else `'*'` |
| `<` | `<=`: tokLessEq; `<=>`: tokNullSafeEq; `<>`: tokNotEq; `<<`: tokShiftLeft; else `'<'` |
| `>` | `>=`: tokGreaterEq; `>>`: tokShiftRight; else `'>'` |
| `!` | `!=`: tokNotEq; `!<`: tokGreaterEq; `!>`: tokLessEq; else `'!'` |
| `&` | `&&`: tokLogicalAnd; else `'&'` |
| `\|` | `\|\|`: tokDoublePipes; else `'\|'` |
| `@` | `@@`: tokDoubleAt; else `'@'` |
| `:` | `:=`: tokAssign; else `':'` |
| `.` | `...`: tokDotDotDot; else `'.'` |
| Other single-char | Return ASCII value as TokenKind |

**Comment handling:**
- `--` followed by space/tab/newline/EOF: skip to end of line (MySQL-compatible strict)
- `//`: skip to end of line
- `/* ... */`: skip block comment (non-nesting)
- `/*+ ... */`: emit tokHintStart; the parser will handle hint content

**Hint strategy:**
- When `/*+` is encountered, emit `tokHintStart` token
- When `*/` is encountered (outside a block comment context), emit `tokHintEnd` token
- The parser is responsible for collecting hint body tokens between these markers

**String scanning:**
- Single-quoted: backslash escapes (`\n`, `\t`, `\r`, `\0`, `\\`, `\'`, `\"`, `\b`, `\Z`) and doubled-quote escape (`''`)
- Double-quoted: same escaping rules
- Backtick-quoted: doubled-backtick escape only (``` `` ``` -> `` ` ``), no backslash escapes
- Unterminated string: emit tokInvalid + LexError, consume to EOF

**Numeric scanning:**
- Leading `0x`/`0X`: hex integer
- Leading `0b`/`0B`: binary integer
- Digits optionally followed by `.` digits: decimal
- Optional exponent `E/e [+-] digits`: float
- Optional suffix `L` (bigint), `S` (smallint), `Y` (tinyint), `BD` (bigdecimal)
- Overflow on `strconv.ParseInt`: downgrade to tokFloat
- Leading `.` + digit: decimal (handled in main dispatch, not here)

**Error recovery:**
- Unterminated string: emit tokInvalid, record LexError, continue from EOF
- Unterminated block comment: emit tokInvalid, record LexError, continue from EOF
- Unknown character: emit tokInvalid for that byte, record LexError, advance one byte

### Testing

Tests in `doris/parser/lexer_test.go`:

- **Keywords**: Verify all 522 keywords tokenize correctly (case-insensitive)
- **Reserved/non-reserved**: Spot-check classification
- **Identifiers**: Unquoted, backtick-quoted with escaping
- **String literals**: Single-quoted, double-quoted, escape sequences, unterminated
- **Numeric literals**: Integer, decimal, float, hex, binary, suffixed, overflow
- **Operators**: All single-char and multi-char operators
- **Comments**: Line (`--`, `//`), block (`/* */`), hints (`/*+ ... */`)
- **Error recovery**: Unterminated string, unterminated comment, unknown character
- **Position tracking**: Verify Loc.Start/End byte offsets, baseOffset behavior
- **Edge cases**: Empty input, whitespace-only, EOF

Run with: `go test ./doris/parser/... -v`
