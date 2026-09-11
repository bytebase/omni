# Doris Lexer Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement a hand-written lexer for Doris SQL that tokenizes all 522 keywords, identifiers, literals, operators, comments, and hints per the DorisLexer.g4 grammar.

**Architecture:** Character-by-character state machine following snowflake/parser/lexer.go patterns. Single dispatch function routes to specialized scanners (string, number, identifier, operator). Keywords recognized via case-insensitive map lookup. Errors accumulated for best-effort recovery.

**Tech Stack:** Go 1.25, module `github.com/bytebase/omni`, depends on `doris/ast` (Loc type)

**Worktree:** `/Users/h3n4l/OpenSource/omni/.worktrees/feat-doris-lexer` (branch `feat/doris/lexer`)

---

## File Structure

| File | Responsibility | ~LOC |
|------|----------------|------|
| `doris/parser/tokens.go` | Token struct, TokenKind type, sentinel/literal/operator constants, LexError type | ~120 |
| `doris/parser/keywords.go` | 522 keyword constants (kw*), keywordMap, nonReservedKeywords, KeywordToken(), IsReserved(), TokenName() | ~2100 |
| `doris/parser/lexer.go` | Lexer struct, NewLexer, NextToken, Tokenize, Errors, all scan* helpers | ~550 |
| `doris/parser/lexer_test.go` | Unit tests covering all token categories | ~500 |

---

### Task 1: Create `doris/parser/tokens.go`

**Files:**
- Create: `doris/parser/tokens.go`

- [ ] **Step 1: Create the file**

```go
package parser

import "github.com/bytebase/omni/doris/ast"

// Special tokens.
const (
	tokEOF     = 0
	tokInvalid = 1 // error-recovery token; always accompanied by a LexError
)

// TokenKind identifies a token type. Single-character operators and
// punctuation use their ASCII byte value directly as their TokenKind
// (e.g., ';' = 59, '(' = 40). Named constants below cover multi-char
// operators, literals, and keywords.
type TokenKind = int

// Multi-character operators (500-599).
const (
	tokLessEq      = 500 + iota // <=
	tokGreaterEq                // >=
	tokNotEq                    // <> or !=
	tokNullSafeEq               // <=>
	tokLogicalAnd               // &&
	tokDoublePipes              // ||
	tokShiftLeft                // <<
	tokShiftRight               // >>
	tokArrow                    // ->
	tokDoubleAt                 // @@
	tokAssign                   // :=
	tokDotDotDot                // ...
	tokHintStart                // /*+
	tokHintEnd                  // */ (only emitted after a hint-start)
)

// Literal tokens (600-699).
const (
	tokInt         = 600 + iota // integer: 42, 0xFF, 0b101
	tokFloat                    // decimal/exponent: 3.14, 1e10, 1.5e-10
	tokString                   // single or double quoted string
	tokIdent                    // unquoted identifier
	tokQuotedIdent              // backtick-quoted identifier
	tokHexLiteral               // X'...' or x'...'
	tokBitLiteral               // B'...' or b'...'
	tokPlaceholder              // ?
)

// Token represents a single lexical token.
type Token struct {
	Kind TokenKind // tok*/kw* constant or ASCII byte value
	Str  string    // content for identifiers, strings, hints
	Ival int64     // parsed integer value for tokInt
	Loc  ast.Loc   // byte offset span in source
}

// LexError records a lexing error with its source position.
type LexError struct {
	Msg string
	Loc ast.Loc
}

// Error messages.
const (
	errUnterminatedString  = "unterminated string literal"
	errUnterminatedQuoted  = "unterminated backtick-quoted identifier"
	errUnterminatedComment = "unterminated block comment"
	errUnknownChar         = "unknown character"
)
```

- [ ] **Step 2: Verify it compiles**

Run: `cd /Users/h3n4l/OpenSource/omni/.worktrees/feat-doris-lexer && go build ./doris/parser/...`

Expected: May fail because package has only one file. That's fine.

- [ ] **Step 3: Commit**

```bash
git add doris/parser/tokens.go
git commit -m "feat(doris/parser): add Token, TokenKind, and LexError types"
```

---

### Task 2: Create `doris/parser/keywords.go`

**Files:**
- Create: `doris/parser/keywords.go`

This is the largest file (~2100 LOC). It contains:
1. 522 keyword constants (kwACCOUNT_LOCK through kwYEAR) starting at 700
2. A `keywordMap` mapping lowercase strings to keyword constants
3. A `nonReservedKeywords` set (172 entries)
4. Public lookup functions

- [ ] **Step 1: Create the file with all 522 keywords**

The keyword list must be extracted from `/Users/h3n4l/OpenSource/parser/doris/DorisLexer.g4`. The file structure is:

```go
package parser

import "strings"

// Keyword token constants. Values start at 700 and grow with iota.
// Extracted from DorisLexer.g4 (522 keywords).
// Numeric values are NOT stable — do not persist them.
const (
	kwACCOUNT_LOCK TokenKind = 700 + iota
	kwACCOUNT_UNLOCK
	kwACTIONS
	kwADD
	kwADMIN
	kwAFTER
	kwAGG_STATE
	kwAGGREGATE
	kwALIAS
	kwALL
	kwALTER
	kwALWAYS
	kwANALYZE
	kwANALYZED
	kwANALYZER
	kwAND
	kwANN
	kwANTI
	kwAPPEND
	kwARRAY
	kwAS
	kwASC
	kwAT
	kwAUTHORS
	kwAUTO
	kwAUTO_INCREMENT
	// ... all 522 keywords in alphabetical order through kwYEAR
)
```

To build this file correctly, the implementer MUST:
1. Read `/Users/h3n4l/OpenSource/parser/doris/DorisLexer.g4` to get the complete keyword list
2. Generate all 522 `kw*` constants in alphabetical order starting at `700 + iota`
3. Build the `keywordMap` with lowercase keys mapping to constants
4. Build the `nonReservedKeywords` set from the `nonReserved` rule in `/Users/h3n4l/OpenSource/parser/doris/DorisParser.g4` (lines 1906-2257)
5. Add public functions: `KeywordToken`, `IsReserved`, `TokenName`

The public functions:

```go
// KeywordToken returns the keyword TokenKind for a name, or (0, false)
// if name is not a keyword. Lookup is case-insensitive.
func KeywordToken(name string) (TokenKind, bool) {
	kind, ok := keywordMap[strings.ToLower(name)]
	return kind, ok
}

// IsReserved reports whether kind is a reserved keyword that cannot be
// used as an unquoted identifier.
func IsReserved(kind TokenKind) bool {
	return kind >= 700 && !nonReservedKeywords[kind]
}

// tokenNames is lazily initialized by TokenName.
var tokenNames map[TokenKind]string

// TokenName returns a human-readable name for a TokenKind.
func TokenName(kind TokenKind) string {
	if tokenNames == nil {
		tokenNames = make(map[TokenKind]string, len(keywordMap)+20)
		for name, kind := range keywordMap {
			tokenNames[kind] = strings.ToUpper(name)
		}
		tokenNames[tokEOF] = "EOF"
		tokenNames[tokInvalid] = "INVALID"
		tokenNames[tokInt] = "INT"
		tokenNames[tokFloat] = "FLOAT"
		tokenNames[tokString] = "STRING"
		tokenNames[tokIdent] = "IDENT"
		tokenNames[tokQuotedIdent] = "QUOTED_IDENT"
		tokenNames[tokHexLiteral] = "HEX_LITERAL"
		tokenNames[tokBitLiteral] = "BIT_LITERAL"
		tokenNames[tokPlaceholder] = "PLACEHOLDER"
		tokenNames[tokLessEq] = "<="
		tokenNames[tokGreaterEq] = ">="
		tokenNames[tokNotEq] = "<>"
		tokenNames[tokNullSafeEq] = "<=>"
		tokenNames[tokLogicalAnd] = "&&"
		tokenNames[tokDoublePipes] = "||"
		tokenNames[tokShiftLeft] = "<<"
		tokenNames[tokShiftRight] = ">>"
		tokenNames[tokArrow] = "->"
		tokenNames[tokDoubleAt] = "@@"
		tokenNames[tokAssign] = ":="
		tokenNames[tokDotDotDot] = "..."
		tokenNames[tokHintStart] = "HINT_START"
		tokenNames[tokHintEnd] = "HINT_END"
	}
	if name, ok := tokenNames[kind]; ok {
		return name
	}
	if kind > 0 && kind < 128 {
		return string(rune(kind))
	}
	return "UNKNOWN"
}
```

- [ ] **Step 2: Verify it compiles**

Run: `cd /Users/h3n4l/OpenSource/omni/.worktrees/feat-doris-lexer && go build ./doris/parser/...`

Expected: Success (tokens.go + keywords.go together should compile).

- [ ] **Step 3: Commit**

```bash
git add doris/parser/keywords.go
git commit -m "feat(doris/parser): add 522 keyword constants and lookup functions"
```

---

### Task 3: Create `doris/parser/lexer.go`

**Files:**
- Create: `doris/parser/lexer.go`

This is the core scanning engine. Follow the snowflake/parser/lexer.go pattern exactly.

- [ ] **Step 1: Create the file**

```go
package parser

import (
	"strconv"
	"strings"

	"github.com/bytebase/omni/doris/ast"
)

// Lexer is a Doris SQL tokenizer. Construct via NewLexer (or
// NewLexerWithOffset when tokenizing a substring) and call NextToken
// until it returns Token{Kind: tokEOF}. Lex errors are collected in
// Errors; each error is accompanied by a tokInvalid token so consumers
// can halt or proceed with best-effort parsing.
type Lexer struct {
	input      string
	pos        int // current byte offset
	start      int // start byte of current token
	errors     []LexError
	baseOffset int // added to Loc positions on return
	inHint     bool // true between /*+ and */
}

// NewLexer creates a lexer for the given input.
func NewLexer(input string) *Lexer {
	return &Lexer{input: input}
}

// NewLexerWithOffset creates a lexer whose emitted Loc values are shifted
// by baseOffset. Used for multi-statement parsing.
func NewLexerWithOffset(input string, baseOffset int) *Lexer {
	return &Lexer{input: input, baseOffset: baseOffset}
}

// Errors returns all lex errors collected so far, with positions shifted
// by baseOffset if applicable.
func (l *Lexer) Errors() []LexError {
	if l.baseOffset == 0 {
		return l.errors
	}
	shifted := make([]LexError, len(l.errors))
	for i, e := range l.errors {
		shifted[i] = LexError{
			Loc: ast.Loc{Start: e.Loc.Start + l.baseOffset, End: e.Loc.End + l.baseOffset},
			Msg: e.Msg,
		}
	}
	return shifted
}

// Tokenize is a one-shot convenience that lexes the entire input.
func Tokenize(input string) ([]Token, []LexError) {
	l := NewLexer(input)
	var tokens []Token
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Kind == tokEOF {
			break
		}
	}
	return tokens, l.Errors()
}

// NextToken returns the next token with baseOffset-shifted positions.
func (l *Lexer) NextToken() Token {
	tok := l.nextTokenInner()
	if l.baseOffset != 0 {
		tok.Loc.Start += l.baseOffset
		tok.Loc.End += l.baseOffset
	}
	return tok
}

func (l *Lexer) nextTokenInner() Token {
	l.skipWhitespaceAndComments()
	if l.pos >= len(l.input) {
		return Token{Kind: tokEOF, Loc: ast.Loc{Start: l.pos, End: l.pos}}
	}
	l.start = l.pos
	ch := l.input[l.pos]

	switch {
	case ch == '\'':
		return l.scanString('\'')
	case ch == '"':
		return l.scanString('"')
	case ch == '`':
		return l.scanBacktickIdent()
	case ch >= '0' && ch <= '9':
		return l.scanNumber()
	case ch == '.' && l.peekDigit():
		return l.scanNumber()
	case (ch == 'X' || ch == 'x') && l.peek(1) == '\'':
		return l.scanHexLiteral()
	case (ch == 'B' || ch == 'b') && l.peek(1) == '\'':
		return l.scanBitLiteral()
	case isIdentStart(ch):
		return l.scanIdentOrKeyword()
	}
	return l.scanOperator(ch)
}

// --- Helpers ---------------------------------------------------------------

func (l *Lexer) peek(offset int) byte {
	if l.pos+offset >= len(l.input) {
		return 0
	}
	return l.input[l.pos+offset]
}

func (l *Lexer) peekDigit() bool {
	c := l.peek(1)
	return c >= '0' && c <= '9'
}

func isIdentStart(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_' || ch == '$' || ch > 127
}

func isIdentCont(ch byte) bool {
	return isIdentStart(ch) || (ch >= '0' && ch <= '9')
}

// --- Whitespace & Comments -------------------------------------------------

func (l *Lexer) skipWhitespaceAndComments() {
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		switch {
		case ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r':
			l.pos++
		case ch == '-' && l.peek(1) == '-' && l.isLineCommentAfterDash():
			// MySQL-compatible: -- must be followed by space/tab/newline/EOF
			l.pos += 2
			for l.pos < len(l.input) && l.input[l.pos] != '\n' {
				l.pos++
			}
		case ch == '/' && l.peek(1) == '/':
			l.pos += 2
			for l.pos < len(l.input) && l.input[l.pos] != '\n' {
				l.pos++
			}
		case ch == '#':
			l.pos++
			for l.pos < len(l.input) && l.input[l.pos] != '\n' {
				l.pos++
			}
		case ch == '/' && l.peek(1) == '*' && l.peek(2) != '+':
			// Block comment (but NOT hint /*+)
			l.scanBlockComment()
		default:
			return
		}
	}
}

// isLineCommentAfterDash checks that byte after -- is space/tab/newline/EOF.
func (l *Lexer) isLineCommentAfterDash() bool {
	c := l.peek(2)
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == 0
}

func (l *Lexer) scanBlockComment() {
	start := l.pos
	l.pos += 2 // consume /*
	for l.pos < len(l.input)-1 {
		if l.input[l.pos] == '*' && l.input[l.pos+1] == '/' {
			l.pos += 2
			return
		}
		l.pos++
	}
	l.pos = len(l.input)
	l.errors = append(l.errors, LexError{
		Loc: ast.Loc{Start: start, End: l.pos},
		Msg: errUnterminatedComment,
	})
}

// --- String Scanning -------------------------------------------------------

// scanString reads a single-quoted or double-quoted string literal.
// Supports backslash escapes and doubled-quote escapes.
func (l *Lexer) scanString(quote byte) Token {
	start := l.start
	l.pos++ // consume opening quote
	var sb strings.Builder
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == quote {
			if l.peek(1) == quote {
				sb.WriteByte(quote)
				l.pos += 2
				continue
			}
			l.pos++
			return Token{Kind: tokString, Str: sb.String(), Loc: ast.Loc{Start: start, End: l.pos}}
		}
		if ch == '\\' {
			l.pos++
			if l.pos >= len(l.input) {
				break
			}
			esc := l.input[l.pos]
			switch esc {
			case 'n':
				sb.WriteByte('\n')
			case 't':
				sb.WriteByte('\t')
			case 'r':
				sb.WriteByte('\r')
			case '0':
				sb.WriteByte(0)
			case '\\':
				sb.WriteByte('\\')
			case '\'':
				sb.WriteByte('\'')
			case '"':
				sb.WriteByte('"')
			case 'b':
				sb.WriteByte(0x08)
			case 'Z':
				sb.WriteByte(0x1A)
			default:
				sb.WriteByte(esc)
			}
			l.pos++
			continue
		}
		sb.WriteByte(ch)
		l.pos++
	}
	l.errors = append(l.errors, LexError{Loc: ast.Loc{Start: start, End: l.pos}, Msg: errUnterminatedString})
	return Token{Kind: tokInvalid, Loc: ast.Loc{Start: start, End: l.pos}}
}

// scanBacktickIdent reads a backtick-quoted identifier. Only doubled-backtick
// escape is supported (`` → `). No backslash escapes.
func (l *Lexer) scanBacktickIdent() Token {
	start := l.start
	l.pos++ // consume opening backtick
	var sb strings.Builder
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == '`' {
			if l.peek(1) == '`' {
				sb.WriteByte('`')
				l.pos += 2
				continue
			}
			l.pos++
			return Token{Kind: tokQuotedIdent, Str: sb.String(), Loc: ast.Loc{Start: start, End: l.pos}}
		}
		sb.WriteByte(ch)
		l.pos++
	}
	l.errors = append(l.errors, LexError{Loc: ast.Loc{Start: start, End: l.pos}, Msg: errUnterminatedQuoted})
	return Token{Kind: tokInvalid, Loc: ast.Loc{Start: start, End: l.pos}}
}

// --- Hex/Bit Literal Scanning -----------------------------------------------

func (l *Lexer) scanHexLiteral() Token {
	start := l.start
	l.pos += 2 // consume X and opening quote
	contentStart := l.pos
	for l.pos < len(l.input) && l.input[l.pos] != '\'' {
		l.pos++
	}
	if l.pos >= len(l.input) {
		l.errors = append(l.errors, LexError{Loc: ast.Loc{Start: start, End: l.pos}, Msg: errUnterminatedString})
		return Token{Kind: tokInvalid, Loc: ast.Loc{Start: start, End: l.pos}}
	}
	content := l.input[contentStart:l.pos]
	l.pos++ // consume closing quote
	return Token{Kind: tokHexLiteral, Str: content, Loc: ast.Loc{Start: start, End: l.pos}}
}

func (l *Lexer) scanBitLiteral() Token {
	start := l.start
	l.pos += 2 // consume B and opening quote
	contentStart := l.pos
	for l.pos < len(l.input) && l.input[l.pos] != '\'' {
		l.pos++
	}
	if l.pos >= len(l.input) {
		l.errors = append(l.errors, LexError{Loc: ast.Loc{Start: start, End: l.pos}, Msg: errUnterminatedString})
		return Token{Kind: tokInvalid, Loc: ast.Loc{Start: start, End: l.pos}}
	}
	content := l.input[contentStart:l.pos]
	l.pos++ // consume closing quote
	return Token{Kind: tokBitLiteral, Str: content, Loc: ast.Loc{Start: start, End: l.pos}}
}

// --- Number Scanning -------------------------------------------------------

func (l *Lexer) scanNumber() Token {
	start := l.start

	// Hex: 0x or 0X
	if l.input[l.pos] == '0' && (l.peek(1) == 'x' || l.peek(1) == 'X') {
		l.pos += 2
		for l.pos < len(l.input) && isHexDigit(l.input[l.pos]) {
			l.pos++
		}
		text := l.input[start:l.pos]
		val, err := strconv.ParseInt(text[2:], 16, 64)
		if err != nil {
			return Token{Kind: tokFloat, Str: text, Loc: ast.Loc{Start: start, End: l.pos}}
		}
		return Token{Kind: tokInt, Str: text, Ival: val, Loc: ast.Loc{Start: start, End: l.pos}}
	}

	// Binary: 0b or 0B
	if l.input[l.pos] == '0' && (l.peek(1) == 'b' || l.peek(1) == 'B') && l.peek(2) != '\'' {
		l.pos += 2
		for l.pos < len(l.input) && (l.input[l.pos] == '0' || l.input[l.pos] == '1') {
			l.pos++
		}
		text := l.input[start:l.pos]
		val, err := strconv.ParseInt(text[2:], 2, 64)
		if err != nil {
			return Token{Kind: tokFloat, Str: text, Loc: ast.Loc{Start: start, End: l.pos}}
		}
		return Token{Kind: tokInt, Str: text, Ival: val, Loc: ast.Loc{Start: start, End: l.pos}}
	}

	isFloat := false

	// Leading dot (.5)
	if l.input[l.pos] == '.' {
		isFloat = true
		l.pos++
		for l.pos < len(l.input) && l.input[l.pos] >= '0' && l.input[l.pos] <= '9' {
			l.pos++
		}
	} else {
		// Integer part
		for l.pos < len(l.input) && l.input[l.pos] >= '0' && l.input[l.pos] <= '9' {
			l.pos++
		}
		// Optional decimal
		if l.pos < len(l.input) && l.input[l.pos] == '.' {
			isFloat = true
			l.pos++
			for l.pos < len(l.input) && l.input[l.pos] >= '0' && l.input[l.pos] <= '9' {
				l.pos++
			}
		}
	}

	// Optional exponent
	if l.pos < len(l.input) && (l.input[l.pos] == 'e' || l.input[l.pos] == 'E') {
		isFloat = true
		l.pos++
		if l.pos < len(l.input) && (l.input[l.pos] == '+' || l.input[l.pos] == '-') {
			l.pos++
		}
		for l.pos < len(l.input) && l.input[l.pos] >= '0' && l.input[l.pos] <= '9' {
			l.pos++
		}
	}

	text := l.input[start:l.pos]
	loc := ast.Loc{Start: start, End: l.pos}

	if isFloat {
		return Token{Kind: tokFloat, Str: text, Loc: loc}
	}
	val, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		return Token{Kind: tokFloat, Str: text, Loc: loc}
	}
	return Token{Kind: tokInt, Str: text, Ival: val, Loc: loc}
}

func isHexDigit(ch byte) bool {
	return (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}

// --- Identifier / Keyword Scanning -----------------------------------------

func (l *Lexer) scanIdentOrKeyword() Token {
	start := l.start
	for l.pos < len(l.input) && isIdentCont(l.input[l.pos]) {
		l.pos++
	}
	word := l.input[start:l.pos]
	loc := ast.Loc{Start: start, End: l.pos}
	if kw, ok := KeywordToken(word); ok {
		return Token{Kind: kw, Str: word, Loc: loc}
	}
	return Token{Kind: tokIdent, Str: word, Loc: loc}
}

// --- Operator Scanning -----------------------------------------------------

func (l *Lexer) scanOperator(ch byte) Token {
	start := l.start
	l.pos++

	switch ch {
	case '<':
		if l.pos < len(l.input) {
			switch l.input[l.pos] {
			case '=':
				if l.peek(1) == '>' {
					l.pos += 2
					return Token{Kind: tokNullSafeEq, Loc: ast.Loc{Start: start, End: l.pos}}
				}
				l.pos++
				return Token{Kind: tokLessEq, Loc: ast.Loc{Start: start, End: l.pos}}
			case '>':
				l.pos++
				return Token{Kind: tokNotEq, Loc: ast.Loc{Start: start, End: l.pos}}
			case '<':
				l.pos++
				return Token{Kind: tokShiftLeft, Loc: ast.Loc{Start: start, End: l.pos}}
			}
		}
	case '>':
		if l.pos < len(l.input) {
			switch l.input[l.pos] {
			case '=':
				l.pos++
				return Token{Kind: tokGreaterEq, Loc: ast.Loc{Start: start, End: l.pos}}
			case '>':
				l.pos++
				return Token{Kind: tokShiftRight, Loc: ast.Loc{Start: start, End: l.pos}}
			}
		}
	case '!':
		if l.pos < len(l.input) {
			switch l.input[l.pos] {
			case '=':
				l.pos++
				return Token{Kind: tokNotEq, Loc: ast.Loc{Start: start, End: l.pos}}
			case '<':
				l.pos++
				return Token{Kind: tokGreaterEq, Loc: ast.Loc{Start: start, End: l.pos}}
			case '>':
				l.pos++
				return Token{Kind: tokLessEq, Loc: ast.Loc{Start: start, End: l.pos}}
			}
		}
	case '&':
		if l.pos < len(l.input) && l.input[l.pos] == '&' {
			l.pos++
			return Token{Kind: tokLogicalAnd, Loc: ast.Loc{Start: start, End: l.pos}}
		}
	case '|':
		if l.pos < len(l.input) && l.input[l.pos] == '|' {
			l.pos++
			return Token{Kind: tokDoublePipes, Loc: ast.Loc{Start: start, End: l.pos}}
		}
	case '@':
		if l.pos < len(l.input) && l.input[l.pos] == '@' {
			l.pos++
			return Token{Kind: tokDoubleAt, Loc: ast.Loc{Start: start, End: l.pos}}
		}
	case ':':
		if l.pos < len(l.input) && l.input[l.pos] == '=' {
			l.pos++
			return Token{Kind: tokAssign, Loc: ast.Loc{Start: start, End: l.pos}}
		}
	case '.':
		if l.pos+1 < len(l.input) && l.input[l.pos] == '.' && l.input[l.pos+1] == '.' {
			l.pos += 2
			return Token{Kind: tokDotDotDot, Loc: ast.Loc{Start: start, End: l.pos}}
		}
	case '-':
		if l.pos < len(l.input) && l.input[l.pos] == '>' {
			l.pos++
			return Token{Kind: tokArrow, Loc: ast.Loc{Start: start, End: l.pos}}
		}
	case '/':
		if l.pos+1 < len(l.input) && l.input[l.pos] == '*' && l.input[l.pos+1] == '+' {
			l.pos += 2
			return Token{Kind: tokHintStart, Loc: ast.Loc{Start: start, End: l.pos}}
		}
	case '*':
		if l.pos < len(l.input) && l.input[l.pos] == '/' {
			l.pos++
			return Token{Kind: tokHintEnd, Loc: ast.Loc{Start: start, End: l.pos}}
		}
	case '?':
		return Token{Kind: tokPlaceholder, Loc: ast.Loc{Start: start, End: l.pos}}
	case '=':
		// == is equivalent to =
		if l.pos < len(l.input) && l.input[l.pos] == '=' {
			l.pos++
		}
		return Token{Kind: int('='), Loc: ast.Loc{Start: start, End: l.pos}}
	}

	// Valid single-char operators/punctuation
	if ch == '(' || ch == ')' || ch == ',' || ch == ';' || ch == '+' || ch == '-' ||
		ch == '*' || ch == '/' || ch == '%' || ch == '~' || ch == '^' || ch == '&' ||
		ch == '|' || ch == '<' || ch == '>' || ch == '=' || ch == '.' || ch == '@' ||
		ch == ':' || ch == '!' || ch == '[' || ch == ']' || ch == '{' || ch == '}' {
		return Token{Kind: int(ch), Loc: ast.Loc{Start: start, End: l.pos}}
	}

	// Unknown character
	l.errors = append(l.errors, LexError{
		Loc: ast.Loc{Start: start, End: l.pos},
		Msg: errUnknownChar,
	})
	return Token{Kind: tokInvalid, Loc: ast.Loc{Start: start, End: l.pos}}
}
```

- [ ] **Step 2: Verify compilation**

Run: `cd /Users/h3n4l/OpenSource/omni/.worktrees/feat-doris-lexer && go build ./doris/parser/... && go vet ./doris/parser/...`

Expected: Success.

- [ ] **Step 3: Commit**

```bash
git add doris/parser/lexer.go
git commit -m "feat(doris/parser): add hand-written lexer with full Doris SQL support"
```

---

### Task 4: Write tests in `doris/parser/lexer_test.go`

**Files:**
- Create: `doris/parser/lexer_test.go`

- [ ] **Step 1: Create the test file**

```go
package parser

import (
	"testing"

	"github.com/bytebase/omni/doris/ast"
)

func TestTokenizeKeywords(t *testing.T) {
	// Spot-check a representative set of keywords, case-insensitive.
	tests := []struct {
		input string
		want  TokenKind
	}{
		{"SELECT", kwSELECT},
		{"select", kwSELECT},
		{"SeLeCt", kwSELECT},
		{"FROM", kwFROM},
		{"CREATE", kwCREATE},
		{"DISTRIBUTED", kwDISTRIBUTED},
		{"AGGREGATE", kwAGGREGATE},
		{"MATERIALIZED", kwMATERIALIZED},
		{"BUCKETS", kwBUCKETS},
	}
	for _, tt := range tests {
		tokens, errs := Tokenize(tt.input)
		if len(errs) > 0 {
			t.Errorf("Tokenize(%q) errors: %v", tt.input, errs)
			continue
		}
		if len(tokens) < 2 || tokens[0].Kind != tt.want {
			t.Errorf("Tokenize(%q)[0].Kind = %d, want %d", tt.input, tokens[0].Kind, tt.want)
		}
	}
}

func TestTokenizeIdentifiers(t *testing.T) {
	tests := []struct {
		input    string
		wantKind TokenKind
		wantStr  string
	}{
		{"myTable", tokIdent, "myTable"},
		{"`my table`", tokQuotedIdent, "my table"},
		{"``backtick``", tokQuotedIdent, "backtick"},  // escaped backtick: `` -> `... wait
	}
	for _, tt := range tests {
		tokens, errs := Tokenize(tt.input)
		if len(errs) > 0 {
			t.Errorf("Tokenize(%q) errors: %v", tt.input, errs)
			continue
		}
		if tokens[0].Kind != tt.wantKind {
			t.Errorf("Tokenize(%q)[0].Kind = %d, want %d", tt.input, tokens[0].Kind, tt.wantKind)
		}
		if tokens[0].Str != tt.wantStr {
			t.Errorf("Tokenize(%q)[0].Str = %q, want %q", tt.input, tokens[0].Str, tt.wantStr)
		}
	}
}

func TestTokenizeBacktickEscape(t *testing.T) {
	// Doubled backtick `` inside backtick-quoted ident -> single backtick
	tokens, errs := Tokenize("`col``name`")
	if len(errs) > 0 {
		t.Fatalf("errors: %v", errs)
	}
	if tokens[0].Kind != tokQuotedIdent || tokens[0].Str != "col`name" {
		t.Errorf("got Kind=%d Str=%q, want tokQuotedIdent %q", tokens[0].Kind, tokens[0].Str, "col`name")
	}
}

func TestTokenizeStrings(t *testing.T) {
	tests := []struct {
		input   string
		wantStr string
	}{
		{`'hello'`, "hello"},
		{`'it''s'`, "it's"},        // doubled-quote escape
		{`'line\nbreak'`, "line\nbreak"}, // backslash escape
		{`"double"`, "double"},
		{`'tab\there'`, "tab\there"},
		{`'null\0byte'`, "null\x00byte"},
	}
	for _, tt := range tests {
		tokens, errs := Tokenize(tt.input)
		if len(errs) > 0 {
			t.Errorf("Tokenize(%q) errors: %v", tt.input, errs)
			continue
		}
		if tokens[0].Kind != tokString {
			t.Errorf("Tokenize(%q)[0].Kind = %d, want tokString", tt.input, tokens[0].Kind)
			continue
		}
		if tokens[0].Str != tt.wantStr {
			t.Errorf("Tokenize(%q)[0].Str = %q, want %q", tt.input, tokens[0].Str, tt.wantStr)
		}
	}
}

func TestTokenizeNumbers(t *testing.T) {
	tests := []struct {
		input    string
		wantKind TokenKind
		wantIval int64
		wantStr  string
	}{
		{"42", tokInt, 42, "42"},
		{"0", tokInt, 0, "0"},
		{"3.14", tokFloat, 0, "3.14"},
		{".5", tokFloat, 0, ".5"},
		{"1e10", tokFloat, 0, "1e10"},
		{"1.5e-10", tokFloat, 0, "1.5e-10"},
		{"0xFF", tokInt, 255, "0xFF"},
		{"0b101", tokInt, 5, "0b101"},
	}
	for _, tt := range tests {
		tokens, errs := Tokenize(tt.input)
		if len(errs) > 0 {
			t.Errorf("Tokenize(%q) errors: %v", tt.input, errs)
			continue
		}
		tok := tokens[0]
		if tok.Kind != tt.wantKind {
			t.Errorf("Tokenize(%q).Kind = %d, want %d", tt.input, tok.Kind, tt.wantKind)
		}
		if tok.Kind == tokInt && tok.Ival != tt.wantIval {
			t.Errorf("Tokenize(%q).Ival = %d, want %d", tt.input, tok.Ival, tt.wantIval)
		}
		if tok.Str != tt.wantStr {
			t.Errorf("Tokenize(%q).Str = %q, want %q", tt.input, tok.Str, tt.wantStr)
		}
	}
}

func TestTokenizeOperators(t *testing.T) {
	tests := []struct {
		input string
		want  TokenKind
	}{
		{"<=", tokLessEq},
		{">=", tokGreaterEq},
		{"<>", tokNotEq},
		{"!=", tokNotEq},
		{"<=>", tokNullSafeEq},
		{"&&", tokLogicalAnd},
		{"||", tokDoublePipes},
		{"<<", tokShiftLeft},
		{">>", tokShiftRight},
		{"->", tokArrow},
		{"@@", tokDoubleAt},
		{":=", tokAssign},
		{"...", tokDotDotDot},
		{"+", int('+')},
		{"-", int('-')},
		{"*", int('*')},
		{"/", int('/')},
		{"%", int('%')},
		{"~", int('~')},
		{"^", int('^')},
		{"(", int('(')},
		{")", int(')')},
		{",", int(',')},
		{";", int(';')},
		{".", int('.')},
		{"==", int('=')}, // == folds to =
	}
	for _, tt := range tests {
		tokens, errs := Tokenize(tt.input)
		if len(errs) > 0 {
			t.Errorf("Tokenize(%q) errors: %v", tt.input, errs)
			continue
		}
		if tokens[0].Kind != tt.want {
			t.Errorf("Tokenize(%q)[0].Kind = %d, want %d (%s)", tt.input, tokens[0].Kind, tt.want, TokenName(tt.want))
		}
	}
}

func TestTokenizeComments(t *testing.T) {
	tests := []struct {
		input string
		want  string // expected first non-EOF token Str
	}{
		{"-- comment\nSELECT", "SELECT"},
		{"// comment\nSELECT", "SELECT"},
		{"# comment\nSELECT", "SELECT"},
		{"/* block */SELECT", "SELECT"},
	}
	for _, tt := range tests {
		tokens, errs := Tokenize(tt.input)
		if len(errs) > 0 {
			t.Errorf("Tokenize(%q) errors: %v", tt.input, errs)
			continue
		}
		if tokens[0].Str != tt.want {
			t.Errorf("Tokenize(%q)[0].Str = %q, want %q", tt.input, tokens[0].Str, tt.want)
		}
	}
}

func TestTokenizeDashDashRequiresSpace(t *testing.T) {
	// -- without space after is two minus tokens, not a comment.
	tokens, _ := Tokenize("--5")
	if tokens[0].Kind != int('-') {
		t.Errorf("got Kind=%d, want '-' (%d)", tokens[0].Kind, int('-'))
	}
}

func TestTokenizeHints(t *testing.T) {
	tokens, errs := Tokenize("/*+ SET_VAR(k=v) */")
	if len(errs) > 0 {
		t.Fatalf("errors: %v", errs)
	}
	if tokens[0].Kind != tokHintStart {
		t.Errorf("tokens[0].Kind = %d, want tokHintStart", tokens[0].Kind)
	}
}

func TestTokenizeHexBitLiterals(t *testing.T) {
	tokens, errs := Tokenize("X'FF' B'101'")
	if len(errs) > 0 {
		t.Fatalf("errors: %v", errs)
	}
	if tokens[0].Kind != tokHexLiteral || tokens[0].Str != "FF" {
		t.Errorf("hex: Kind=%d Str=%q", tokens[0].Kind, tokens[0].Str)
	}
	if tokens[1].Kind != tokBitLiteral || tokens[1].Str != "101" {
		t.Errorf("bit: Kind=%d Str=%q", tokens[1].Kind, tokens[1].Str)
	}
}

func TestTokenizePositions(t *testing.T) {
	tokens, _ := Tokenize("SELECT 1")
	// "SELECT" is bytes 0-6, "1" is bytes 7-8
	if tokens[0].Loc != (ast.Loc{Start: 0, End: 6}) {
		t.Errorf("SELECT loc = %v, want {0,6}", tokens[0].Loc)
	}
	if tokens[1].Loc != (ast.Loc{Start: 7, End: 8}) {
		t.Errorf("1 loc = %v, want {7,8}", tokens[1].Loc)
	}
}

func TestTokenizeBaseOffset(t *testing.T) {
	l := NewLexerWithOffset("SELECT", 100)
	tok := l.NextToken()
	if tok.Loc.Start != 100 || tok.Loc.End != 106 {
		t.Errorf("offset loc = %v, want {100,106}", tok.Loc)
	}
}

func TestTokenizeUnterminatedString(t *testing.T) {
	tokens, errs := Tokenize("'unterminated")
	if len(errs) == 0 {
		t.Fatal("expected error for unterminated string")
	}
	if tokens[0].Kind != tokInvalid {
		t.Errorf("Kind = %d, want tokInvalid", tokens[0].Kind)
	}
}

func TestTokenizeUnterminatedComment(t *testing.T) {
	_, errs := Tokenize("/* unterminated")
	if len(errs) == 0 {
		t.Fatal("expected error for unterminated comment")
	}
}

func TestTokenizeUnknownChar(t *testing.T) {
	tokens, errs := Tokenize("\\")
	if len(errs) == 0 {
		t.Fatal("expected error for unknown char")
	}
	if tokens[0].Kind != tokInvalid {
		t.Errorf("Kind = %d, want tokInvalid", tokens[0].Kind)
	}
}

func TestTokenizeEmpty(t *testing.T) {
	tokens, errs := Tokenize("")
	if len(errs) > 0 {
		t.Errorf("errors: %v", errs)
	}
	if len(tokens) != 1 || tokens[0].Kind != tokEOF {
		t.Errorf("expected single EOF token, got %v", tokens)
	}
}

func TestTokenizeMultiStatement(t *testing.T) {
	tokens, errs := Tokenize("SELECT 1; SELECT 2")
	if len(errs) > 0 {
		t.Fatalf("errors: %v", errs)
	}
	// Expect: SELECT, 1, ;, SELECT, 2, EOF = 6 tokens
	if len(tokens) != 6 {
		t.Fatalf("got %d tokens, want 6", len(tokens))
	}
	if tokens[2].Kind != int(';') {
		t.Errorf("tokens[2].Kind = %d, want ';'", tokens[2].Kind)
	}
}

func TestIsReserved(t *testing.T) {
	if !IsReserved(kwSELECT) {
		t.Error("SELECT should be reserved")
	}
	if IsReserved(kwCOMMENT) {
		t.Error("COMMENT should NOT be reserved (it's in nonReserved)")
	}
}

func TestTokenName(t *testing.T) {
	if got := TokenName(kwSELECT); got != "SELECT" {
		t.Errorf("TokenName(kwSELECT) = %q, want %q", got, "SELECT")
	}
	if got := TokenName(tokEOF); got != "EOF" {
		t.Errorf("TokenName(tokEOF) = %q, want %q", got, "EOF")
	}
	if got := TokenName(int(';')); got != ";" {
		t.Errorf("TokenName(';') = %q, want %q", got, ";")
	}
}

func TestTokenizePlaceholder(t *testing.T) {
	tokens, errs := Tokenize("SELECT ? FROM t")
	if len(errs) > 0 {
		t.Fatalf("errors: %v", errs)
	}
	if tokens[1].Kind != tokPlaceholder {
		t.Errorf("tokens[1].Kind = %d, want tokPlaceholder", tokens[1].Kind)
	}
}
```

- [ ] **Step 2: Run tests**

Run: `cd /Users/h3n4l/OpenSource/omni/.worktrees/feat-doris-lexer && go test ./doris/parser/... -v`

Expected: All tests PASS.

- [ ] **Step 3: Commit**

```bash
git add doris/parser/lexer_test.go
git commit -m "test(doris/parser): add lexer unit tests for all token categories"
```

---

### Task 5: Final verification

- [ ] **Step 1: Run full build and vet**

Run: `cd /Users/h3n4l/OpenSource/omni/.worktrees/feat-doris-lexer && go build ./doris/... && go vet ./doris/...`

Expected: Success.

- [ ] **Step 2: Run all doris tests**

Run: `cd /Users/h3n4l/OpenSource/omni/.worktrees/feat-doris-lexer && go test ./doris/... -v -count=1`

Expected: All tests PASS (both ast and parser packages).

- [ ] **Step 3: Verify no regressions**

Run: `cd /Users/h3n4l/OpenSource/omni/.worktrees/feat-doris-lexer && go build ./...`

Expected: Success.
