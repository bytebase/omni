# Snowflake Lexer (F2) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the hand-written Snowflake SQL tokenizer in `snowflake/parser/`. Mysql-style flat lexer with no state machine, full coverage of `SnowflakeLexer.g4`'s ~881 single-word keywords, error recovery via collected `LexError`s + synthetic `tokInvalid` tokens, and a 6-file test suite that includes the 27-file legacy corpus regression baseline.

**Architecture:** Hand-written `Lexer` struct holding input/pos/start/errors. `NextToken()` dispatches by first non-whitespace byte to scan helpers (`scanString`, `scanNumber`, `scanIdentOrKeyword`, `scanQuotedIdent`, `scanDollar`, `scanOperatorOrPunct`). Comments dropped silently. Token positions tracked via F1's `ast.Loc{Start, End}`. ~881 keyword constants and the keyword map are extracted from the legacy `SnowflakeLexer.g4` via a grep+awk pipeline, not hand-typed.

**Tech Stack:** Go 1.25, stdlib only (`strings`, `strconv`, `unicode`, `unicode/utf8`, plus `os`/`path/filepath`/`regexp` for the keyword-completeness test).

**Spec:** `docs/superpowers/specs/2026-04-07-snowflake-lexer-design.md` (commit `d4ad94f`)

**Worktree:** `/Users/h3n4l/OpenSource/omni/.worktrees/lexer` on branch `feat/snowflake/lexer`

**Module path:** `github.com/bytebase/omni`. The lexer's package is `github.com/bytebase/omni/snowflake/parser`. It imports `github.com/bytebase/omni/snowflake/ast` for `ast.Loc`.

**Commit policy:** No commits during implementation. The user reviews the full diff at the end; only after explicit approval do we commit.

---

## File Structure

| File | Responsibility | Approx LOC |
|------|----------------|-----------|
| `snowflake/parser/errors.go` | `LexError` type + error message constants | 35 |
| `snowflake/parser/tokens.go` | `tokXxx` operator/literal constants, ~881 `kwXxx` keyword constants, `Token` struct, `TokenName(int) string` | 1,000 |
| `snowflake/parser/keywords.go` | `keywordMap[string]int` (~881 entries), `keywordReserved[int]bool` (~99 entries), `KeywordToken`, `IsReservedKeyword` | 1,050 |
| `snowflake/parser/lexer.go` | `Lexer` struct, `NewLexer`, `NextToken`, `Tokenize`, `Errors`, helper functions, `skipWhitespaceAndComments`, `scanBlockComment`, all `scan*` helpers | 850 |
| `snowflake/parser/lexer_test.go` | Table-driven token recognition (~80 cases) | 400 |
| `snowflake/parser/position_test.go` | Position math for every token kind | 150 |
| `snowflake/parser/edge_cases_test.go` | EOF mid-string, unterminated comment, ambiguous operators, numeric edges, Unicode | 200 |
| `snowflake/parser/error_recovery_test.go` | Multi-error recovery | 100 |
| `snowflake/parser/legacy_corpus_test.go` | Smoke-tokenize all 27 `.sql` files in `testdata/legacy/` | 80 |
| `snowflake/parser/keyword_completeness_test.go` | Drift detector against `SnowflakeLexer.g4` and `build_id_contains_non_reserved_keywords.py` | 120 |

Total: ~3,985 LOC across 4 production files + 6 test files. The keyword constant block + keyword map dominate (~2,000 LOC combined).

---

## Task 1: Scaffold the package directory and verify empty compile

**Files:**
- Create: `snowflake/parser/errors.go`
- Create: `snowflake/parser/tokens.go`
- Create: `snowflake/parser/keywords.go`
- Create: `snowflake/parser/lexer.go`
- Create: `snowflake/parser/lexer_test.go`
- Create: `snowflake/parser/position_test.go`
- Create: `snowflake/parser/edge_cases_test.go`
- Create: `snowflake/parser/error_recovery_test.go`
- Create: `snowflake/parser/legacy_corpus_test.go`
- Create: `snowflake/parser/keyword_completeness_test.go`

- [ ] **Step 1: Confirm worktree state**

Run: `pwd && git rev-parse --abbrev-ref HEAD`
Expected:
```
/Users/h3n4l/OpenSource/omni/.worktrees/lexer
feat/snowflake/lexer
```

If either is wrong, stop. Do NOT proceed.

- [ ] **Step 2: Create the package directory**

Run: `mkdir -p snowflake/parser`
Expected: no output, exit 0.

- [ ] **Step 3: Stub `snowflake/parser/errors.go`**

Write `snowflake/parser/errors.go`:

```go
// Package parser implements a hand-written Snowflake SQL lexer and (in
// future DAG nodes) a recursive-descent parser. F2 ships only the lexer.
package parser
```

- [ ] **Step 4: Stub `snowflake/parser/tokens.go`**

Write `snowflake/parser/tokens.go`:

```go
package parser
```

- [ ] **Step 5: Stub `snowflake/parser/keywords.go`**

Write `snowflake/parser/keywords.go`:

```go
package parser
```

- [ ] **Step 6: Stub `snowflake/parser/lexer.go`**

Write `snowflake/parser/lexer.go`:

```go
package parser
```

- [ ] **Step 7: Stub all 6 test files**

Write each of the following with content `package parser`:

- `snowflake/parser/lexer_test.go`
- `snowflake/parser/position_test.go`
- `snowflake/parser/edge_cases_test.go`
- `snowflake/parser/error_recovery_test.go`
- `snowflake/parser/legacy_corpus_test.go`
- `snowflake/parser/keyword_completeness_test.go`

- [ ] **Step 8: Verify the package compiles empty**

Run: `go build ./snowflake/parser/...`
Expected: no output, exit 0.

If there is any error, stop and diagnose before continuing.

---

## Task 2: Implement `errors.go`

**Files:**
- Modify: `snowflake/parser/errors.go`

- [ ] **Step 1: Replace the stub with the full file**

Overwrite `snowflake/parser/errors.go` with this exact content:

```go
// Package parser implements a hand-written Snowflake SQL lexer and (in
// future DAG nodes) a recursive-descent parser. F2 ships only the lexer.
package parser

import "github.com/bytebase/omni/snowflake/ast"

// LexError describes a single lexing failure with its source location.
//
// Lex errors are non-fatal: the lexer collects them in a slice (Lexer.Errors)
// and emits a synthetic tokInvalid token at each failure site so consumers
// can choose to halt on the first error or proceed with best-effort parsing.
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

- [ ] **Step 2: Verify compile**

Run: `go build ./snowflake/parser/...`
Expected: no output, exit 0. The unused-import vet warning will fire on `ast` until tokens.go and lexer.go reference it; that's expected at this point — `go build` accepts it but `go vet` may not. Skip vet until Task 5.

If `go build` fails for any reason other than the expected unused-import warning, stop and diagnose.

---

## Task 3: Implement `tokens.go` — operator constants, Token struct, TokenName, and ~881 keyword constants

This is the bulkiest production file. The keyword constants are extracted from `SnowflakeLexer.g4` via a grep+awk pipeline.

**Files:**
- Modify: `snowflake/parser/tokens.go`

- [ ] **Step 1: Write the operator-and-Token portion**

Overwrite `snowflake/parser/tokens.go` with:

```go
package parser

import "github.com/bytebase/omni/snowflake/ast"

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

	tokDoubleColon // ::  (legacy COLON_COLON)
	tokConcat      // ||  (legacy PIPE_PIPE)
	tokArrow       // ->  (legacy ARROW)
	tokFlow        // ->> (legacy FLOW)
	tokAssoc       // =>  (legacy ASSOC)
	tokNotEq       // != or <> (legacy NE / LTGT — folded into one token)
	tokLessEq      // <=  (legacy LE)
	tokGreaterEq   // >=  (legacy GE)
)

// Token represents a single lexical token.
type Token struct {
	Type    int     // tok* / kw* / ASCII byte value
	Str     string  // identifier text, string content, raw operator text
	Ival    int64   // integer value for tokInt
	Loc     ast.Loc // {Start, End} byte offsets in source text
	XPrefix bool    // true if string was X'...' (binary literal); only set on tokString
}

// Keyword token constants. Values start at 700 and grow with iota.
// Generated from /Users/h3n4l/OpenSource/parser/snowflake/SnowflakeLexer.g4
// via the extraction pipeline in docs/superpowers/plans/2026-04-07-snowflake-lexer.md.
// Numeric values are NOT stable across edits — do not persist them.
const (
	// KEYWORD_CONSTANTS_BEGIN — block populated by Step 2 below
	// KEYWORD_CONSTANTS_END
)

// TokenName returns a human-readable name for a token type. Used by tests
// and debug output.
//
// Returns:
//   - "EOF" for tokEOF
//   - "INVALID" for tokInvalid
//   - the upper-case keyword name (e.g. "SELECT") for kw* constants
//   - the constant name minus the "tok" prefix (e.g. "INT", "STRING") for tok* constants
//   - the literal character in quotes (e.g. "'+'") for ASCII single-char tokens
func TokenName(t int) string {
	switch t {
	case tokEOF:
		return "EOF"
	case tokInvalid:
		return "INVALID"
	case tokInt:
		return "INT"
	case tokFloat:
		return "FLOAT"
	case tokReal:
		return "REAL"
	case tokString:
		return "STRING"
	case tokIdent:
		return "IDENT"
	case tokQuotedIdent:
		return "QUOTED_IDENT"
	case tokVariable:
		return "VARIABLE"
	case tokDoubleColon:
		return "::"
	case tokConcat:
		return "||"
	case tokArrow:
		return "->"
	case tokFlow:
		return "->>"
	case tokAssoc:
		return "=>"
	case tokNotEq:
		return "!="
	case tokLessEq:
		return "<="
	case tokGreaterEq:
		return ">="
	}
	if t >= 700 {
		// Keyword token — look up by reverse-mapping the keywordMap. The
		// inverse map is built lazily by keywords.go's keywordName helper.
		if name, ok := keywordName(t); ok {
			return name
		}
		return "UNKNOWN_KEYWORD"
	}
	if t > 0 && t < 128 {
		// ASCII single-char token
		return "'" + string(rune(t)) + "'"
	}
	return "UNKNOWN"
}
```

- [ ] **Step 2: Extract the keyword constants from `SnowflakeLexer.g4`**

Run this pipeline to extract the 881 single-word keyword tokens and emit them as Go iota constants. The first keyword gets `kwACCOUNT = 700 + iota` (assigning the explicit base), the rest get bare names so iota auto-increments.

Run:
```bash
grep -E "^[A-Z_][A-Z0-9_]*[[:space:]]*:[[:space:]]*'[A-Z][A-Z0-9_]*'[[:space:]]*;" \
    /Users/h3n4l/OpenSource/parser/snowflake/SnowflakeLexer.g4 \
  | awk -F"[[:space:]]*:[[:space:]]*'" '{rest=$2; sub(/\047.*$/, "", rest); print rest}' \
  | sort -u \
  > /tmp/snowflake-keywords.txt
wc -l /tmp/snowflake-keywords.txt
```

Expected: `881 /tmp/snowflake-keywords.txt`. If the count is off by more than 5 in either direction, stop and investigate — the legacy grammar may have changed since the spec was written.

- [ ] **Step 3: Format the extracted keywords as Go iota constants**

Run:
```bash
{
  head -1 /tmp/snowflake-keywords.txt | awk '{print "\tkw" $1 " = 700 + iota"}'
  tail -n +2 /tmp/snowflake-keywords.txt | awk '{print "\tkw" $1}'
} > /tmp/snowflake-keywords-go.txt
head -3 /tmp/snowflake-keywords-go.txt
tail -3 /tmp/snowflake-keywords-go.txt
wc -l /tmp/snowflake-keywords-go.txt
```

Expected: head should show something like:
```
	kwAAD_PROVISIONER = 700 + iota
	kwABORT
	kwABORT_DETACHED_QUERY
```
Tail should show 3 keyword constants (last alphabetically). Line count should be 881.

- [ ] **Step 4: Splice the keyword constants into `tokens.go`**

The file currently has the marker block:
```go
const (
	// KEYWORD_CONSTANTS_BEGIN — block populated by Step 2 below
	// KEYWORD_CONSTANTS_END
)
```

Replace those two marker comment lines with the contents of `/tmp/snowflake-keywords-go.txt`. Use `Read` on the file first to get the exact byte context, then `Edit` to replace the marker block. The result should look like:

```go
const (
	// KEYWORD_CONSTANTS_BEGIN — populated from SnowflakeLexer.g4
	kwAAD_PROVISIONER = 700 + iota
	kwABORT
	kwABORT_DETACHED_QUERY
	// ... 878 more entries ...
	kwZSTD
	// KEYWORD_CONSTANTS_END
)
```

- [ ] **Step 5: Verify the file compiles**

Run: `go build ./snowflake/parser/...`
Expected: a single error from `tokens.go` referencing `keywordName` (called from `TokenName`'s `t >= 700` branch). This is expected — `keywordName` is defined in `keywords.go` (Task 4). Note the error and proceed.

If you see any error OTHER than `undefined: keywordName`, stop and investigate.

- [ ] **Step 6: Clean up the extraction temp files**

Run: `rm /tmp/snowflake-keywords.txt /tmp/snowflake-keywords-go.txt`
Expected: no output, exit 0.

---

## Task 4: Implement `keywords.go` — keyword map, reserved set, lookup helpers

**Files:**
- Modify: `snowflake/parser/keywords.go`

- [ ] **Step 1: Write the helper-and-reserved-set portion**

Overwrite `snowflake/parser/keywords.go` with:

```go
package parser

import "strings"

// KeywordToken returns the kw* constant for name (case-insensitive) and
// true if name is a recognized keyword. Returns (0, false) otherwise.
func KeywordToken(name string) (int, bool) {
	t, ok := keywordMap[strings.ToLower(name)]
	return t, ok
}

// IsReservedKeyword reports whether name (case-insensitive) is a reserved
// keyword that cannot be used as an unquoted identifier without quoting.
//
// The set of reserved keywords is seeded from
// /Users/h3n4l/OpenSource/parser/snowflake/build_id_contains_non_reserved_keywords.py
// (the snowflake_reserved_keyword dict literal).
func IsReservedKeyword(name string) bool {
	t, ok := KeywordToken(name)
	return ok && keywordReserved[t]
}

// keywordName returns the upper-case name of a keyword token. Used by
// TokenName for debug output. Builds a reverse-lookup map on first use.
//
// The map is constructed lazily and cached. It does not need to be
// thread-safe — TokenName is only called from tests and debug paths.
var keywordNameCache map[int]string

func keywordName(t int) (string, bool) {
	if keywordNameCache == nil {
		keywordNameCache = make(map[int]string, len(keywordMap))
		for name, tok := range keywordMap {
			keywordNameCache[tok] = strings.ToUpper(name)
		}
	}
	name, ok := keywordNameCache[t]
	return name, ok
}

// keywordReserved marks the 82 reserved keywords seeded from the legacy
// build_id_contains_non_reserved_keywords.py script. These cannot be used
// as unquoted identifiers without quoting.
//
// The legacy Python file declares 91 reserved keywords, but 9 of those
// (CURRENT_USER, FOLLOWING, GSCLUSTER, ISSUE, LOCALTIME, LOCALTIMESTAMP,
// REGEXP, TRIGGER, WHENEVER) do NOT appear as token rules in the current
// SnowflakeLexer.g4 — two are commented out (CURRENT_USER, TRIGGER) and
// the other seven were never added. The omni lexer can only mark a
// keyword as reserved if it has a corresponding kw* token, so these 9
// phantoms are dropped here. Their absence is a deliberate divergence
// from the stale Python source. The keyword_completeness_test.go drift
// detector knows about these 9 names and skips them.
var keywordReserved = map[int]bool{
	kwACCOUNT:           true,
	kwALL:               true,
	kwALTER:             true,
	kwAND:               true,
	kwANY:               true,
	kwAS:                true,
	kwBETWEEN:           true,
	kwBY:                true,
	kwCASE:              true,
	kwCAST:              true,
	kwCHECK:             true,
	kwCOLUMN:            true,
	kwCONNECT:           true,
	kwCONNECTION:        true,
	kwCONSTRAINT:        true,
	kwCREATE:            true,
	kwCROSS:             true,
	kwCURRENT:           true,
	kwCURRENT_DATE:      true,
	kwCURRENT_TIME:      true,
	kwCURRENT_TIMESTAMP: true,
	kwDATABASE:          true,
	kwDELETE:            true,
	kwDISTINCT:          true,
	kwDROP:              true,
	kwELSE:              true,
	kwEXISTS:            true,
	kwFALSE:             true,
	kwFOR:               true,
	kwFROM:              true,
	kwFULL:              true,
	kwGRANT:             true,
	kwGROUP:             true,
	kwHAVING:            true,
	kwILIKE:             true,
	kwIN:                true,
	kwINCREMENT:         true,
	kwINNER:             true,
	kwINSERT:            true,
	kwINTERSECT:         true,
	kwINTO:              true,
	kwIS:                true,
	kwJOIN:              true,
	kwLATERAL:           true,
	kwLEFT:              true,
	kwLIKE:              true,
	kwMINUS:             true,
	kwNATURAL:           true,
	kwNOT:               true,
	kwNULL:              true,
	kwOF:                true,
	kwON:                true,
	kwOR:                true,
	kwORDER:             true,
	kwORGANIZATION:      true,
	kwQUALIFY:           true,
	kwREVOKE:            true,
	kwRIGHT:             true,
	kwRLIKE:             true,
	kwROW:               true,
	kwROWS:              true,
	kwSAMPLE:            true,
	kwSCHEMA:            true,
	kwSELECT:            true,
	kwSET:               true,
	kwSOME:              true,
	kwSTART:             true,
	kwTABLE:             true,
	kwTABLESAMPLE:       true,
	kwTHEN:              true,
	kwTO:                true,
	kwTRUE:              true,
	kwTRY_CAST:          true,
	kwUNION:             true,
	kwUNIQUE:            true,
	kwUPDATE:            true,
	kwUSING:             true,
	kwVALUES:            true,
	kwVIEW:              true,
	kwWHEN:              true,
	kwWHERE:             true,
	kwWITH:              true,
}

// keywordMap maps lowercased keyword text to its kw* constant.
// Generated from /Users/h3n4l/OpenSource/parser/snowflake/SnowflakeLexer.g4
// via the extraction pipeline in docs/superpowers/plans/2026-04-07-snowflake-lexer.md.
var keywordMap = map[string]int{
	// KEYWORD_MAP_BEGIN — block populated by Step 2 below
	// KEYWORD_MAP_END
}
```

- [ ] **Step 2: Generate the keyword map entries**

Run the same extraction pipeline used in Task 3, but format the output as map entries instead of iota constants:

```bash
grep -E "^[A-Z_][A-Z0-9_]*[[:space:]]*:[[:space:]]*'[A-Z][A-Z0-9_]*'[[:space:]]*;" \
    /Users/h3n4l/OpenSource/parser/snowflake/SnowflakeLexer.g4 \
  | awk -F"[[:space:]]*:[[:space:]]*'" '{rest=$2; sub(/\047.*$/, "", rest); print rest}' \
  | sort -u \
  | awk '{print "\t\"" tolower($1) "\": kw" $1 ","}' \
  > /tmp/snowflake-keyword-map.txt
wc -l /tmp/snowflake-keyword-map.txt
head -3 /tmp/snowflake-keyword-map.txt
tail -3 /tmp/snowflake-keyword-map.txt
```

Expected: line count 881. Head should show:
```
	"aad_provisioner": kwAAD_PROVISIONER,
	"abort": kwABORT,
	"abort_detached_query": kwABORT_DETACHED_QUERY,
```

- [ ] **Step 3: Splice the keyword map into `keywords.go`**

Replace the two marker comment lines (`// KEYWORD_MAP_BEGIN` and `// KEYWORD_MAP_END`) with the contents of `/tmp/snowflake-keyword-map.txt`. Use Read+Edit, same as Task 3 Step 4.

- [ ] **Step 4: Verify compile and vet**

Run: `go build ./snowflake/parser/...`
Expected: no output, exit 0. The earlier `undefined: keywordName` error from Task 3 is now resolved because keywords.go provides it.

Run: `go vet ./snowflake/parser/...`
Expected: no output, exit 0.

If either fails, stop and diagnose.

- [ ] **Step 5: Clean up the extraction temp file**

Run: `rm /tmp/snowflake-keyword-map.txt`
Expected: no output, exit 0.

---

## Task 5: Implement `lexer.go` core — Lexer struct, dispatch, helpers, comments

This task writes the lexer body WITHOUT the scan* helpers. Each scan* function is added in Tasks 6–11. To keep the file compiling between tasks, this step writes stub scan* functions that return a placeholder token; subsequent tasks replace each stub with its real implementation.

**Files:**
- Modify: `snowflake/parser/lexer.go`

- [ ] **Step 1: Write the full lexer.go**

Overwrite `snowflake/parser/lexer.go` with:

```go
package parser

import "github.com/bytebase/omni/snowflake/ast"

// Lexer is a Snowflake SQL tokenizer. Construct via NewLexer and call
// NextToken until it returns Token{Type: tokEOF}. Lex errors are collected
// in a slice (Errors); each error is accompanied by a synthetic tokInvalid
// token at the failure site so consumers can choose to halt on the first
// error or proceed with best-effort parsing.
type Lexer struct {
	input  string
	pos    int // current byte offset (one past the last consumed byte)
	start  int // start byte of the token currently being scanned
	errors []LexError
}

// NewLexer constructs a Lexer for the given input.
func NewLexer(input string) *Lexer {
	return &Lexer{input: input}
}

// Errors returns all lex errors collected so far. Errors are appended in
// source order; each error is accompanied by a tokInvalid token in the
// stream.
func (l *Lexer) Errors() []LexError {
	return l.errors
}

// Tokenize is a one-shot convenience that runs a lexer to EOF and returns
// the full token stream and error list. Useful for tests and for callers
// that don't need streaming.
func Tokenize(input string) (tokens []Token, errors []LexError) {
	l := NewLexer(input)
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == tokEOF {
			break
		}
	}
	return tokens, l.Errors()
}

// NextToken advances the lexer and returns the next token. At end of input
// it returns Token{Type: tokEOF}; subsequent calls continue to return EOF.
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

// peek returns the byte at l.pos+offset, or 0 if past end of input.
func (l *Lexer) peek(offset int) byte {
	if l.pos+offset >= len(l.input) {
		return 0
	}
	return l.input[l.pos+offset]
}

// peekDigit reports whether the byte at l.pos+1 is an ASCII digit.
func (l *Lexer) peekDigit() bool {
	c := l.peek(1)
	return c >= '0' && c <= '9'
}

// isIdentStart reports whether ch may begin an unquoted identifier.
// Snowflake's ID rule is [A-Za-z_]; the case-insensitive flag in the
// legacy grammar already covers both cases.
func isIdentStart(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

// isIdentCont reports whether ch may continue an unquoted identifier.
// Snowflake's ID rule allows letters, digits, underscore, @, and $.
func isIdentCont(ch byte) bool {
	return isIdentStart(ch) || (ch >= '0' && ch <= '9') || ch == '@' || ch == '$'
}

// skipWhitespaceAndComments advances l.pos past any whitespace, line
// comments (-- or //), and block comments (/* ... */ with nesting).
// Comments are dropped silently; whitespace produces no token.
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

// scanBlockComment advances past a /* ... */ block comment. Block comments
// nest in Snowflake (the legacy grammar uses recursive
// '/*' (SQL_COMMENT | .)*? '*/'). On unterminated comment, appends a
// LexError but does NOT emit a token (comments are channel-HIDDEN).
func (l *Lexer) scanBlockComment() {
	start := l.pos
	l.pos += 2 // consume /*
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

// scanString is implemented in Task 6.
func (l *Lexer) scanString() Token {
	l.pos = len(l.input) // stub: consume rest of input to make tests fail loudly
	return Token{Type: tokInvalid, Loc: ast.Loc{Start: l.start, End: l.pos}}
}

// scanQuotedIdent is implemented in Task 7.
func (l *Lexer) scanQuotedIdent() Token {
	l.pos = len(l.input)
	return Token{Type: tokInvalid, Loc: ast.Loc{Start: l.start, End: l.pos}}
}

// scanDollar is implemented in Task 8.
func (l *Lexer) scanDollar() Token {
	l.pos = len(l.input)
	return Token{Type: tokInvalid, Loc: ast.Loc{Start: l.start, End: l.pos}}
}

// scanNumber is implemented in Task 9.
func (l *Lexer) scanNumber() Token {
	l.pos = len(l.input)
	return Token{Type: tokInvalid, Loc: ast.Loc{Start: l.start, End: l.pos}}
}

// scanIdentOrKeyword is implemented in Task 10.
func (l *Lexer) scanIdentOrKeyword() Token {
	l.pos = len(l.input)
	return Token{Type: tokInvalid, Loc: ast.Loc{Start: l.start, End: l.pos}}
}

// scanOperatorOrPunct is implemented in Task 11.
func (l *Lexer) scanOperatorOrPunct(ch byte) Token {
	l.pos++
	return Token{Type: int(ch), Loc: ast.Loc{Start: l.start, End: l.pos}}
}
```

Note: the stub `scanString`, `scanQuotedIdent`, `scanDollar`, `scanNumber`, and `scanIdentOrKeyword` advance to end-of-input and return `tokInvalid`. This makes any test that exercises them fail loudly until the real implementation lands. The stub `scanOperatorOrPunct` is partially correct — it just emits the single byte as a token — and Task 11 will replace it with the lookahead version.

- [ ] **Step 2: Verify compile and vet**

Run: `go build ./snowflake/parser/...`
Expected: no output, exit 0.

Run: `go vet ./snowflake/parser/...`
Expected: no output, exit 0.

---

## Task 6: Implement `scanString`

**Files:**
- Modify: `snowflake/parser/lexer.go`

- [ ] **Step 1: Replace the scanString stub**

In `snowflake/parser/lexer.go`, replace the `scanString` stub function with the full implementation:

```go
// scanString reads a single-quoted string literal. Snowflake supports two
// escape mechanisms inside '...':
//   - backslash escapes: \n \t \r \0 \\ \' \"  (and \<other> → <other>)
//   - doubled-quote escape: '' inside a string is a literal '
//
// Token.Str contains the unescaped content (no surrounding quotes).
// On unterminated string, appends an unterminated-string LexError, emits
// a tokInvalid token covering the bad span, and advances to the next
// newline or EOF.
func (l *Lexer) scanString() Token {
	start := l.start // l.start is the byte offset of the opening quote
	l.pos++          // consume opening quote
	var sb strings.Builder
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == '\'' {
			if l.peek(1) == '\'' {
				// Doubled-quote escape: '' → literal '
				sb.WriteByte('\'')
				l.pos += 2
				continue
			}
			// Closing quote.
			l.pos++
			return Token{
				Type: tokString,
				Str:  sb.String(),
				Loc:  ast.Loc{Start: start, End: l.pos},
			}
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
			default:
				sb.WriteByte(esc)
			}
			l.pos++
			continue
		}
		if ch == '\n' {
			// Unterminated single-quoted string — strings cannot span newlines.
			l.errors = append(l.errors, LexError{
				Loc: ast.Loc{Start: start, End: l.pos},
				Msg: errUnterminatedString,
			})
			return Token{Type: tokInvalid, Loc: ast.Loc{Start: start, End: l.pos}}
		}
		sb.WriteByte(ch)
		l.pos++
	}
	// EOF before closing quote.
	l.errors = append(l.errors, LexError{
		Loc: ast.Loc{Start: start, End: l.pos},
		Msg: errUnterminatedString,
	})
	return Token{Type: tokInvalid, Loc: ast.Loc{Start: start, End: l.pos}}
}
```

- [ ] **Step 2: Add the `strings` import to lexer.go**

The file currently imports only `github.com/bytebase/omni/snowflake/ast`. Replace the import block with:

```go
import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)
```

- [ ] **Step 3: Verify compile and vet**

Run: `go build ./snowflake/parser/...`
Expected: no output, exit 0.

Run: `go vet ./snowflake/parser/...`
Expected: no output, exit 0.

---

## Task 7: Implement `scanQuotedIdent`

**Files:**
- Modify: `snowflake/parser/lexer.go`

- [ ] **Step 1: Replace the scanQuotedIdent stub**

In `snowflake/parser/lexer.go`, replace the `scanQuotedIdent` stub function with:

```go
// scanQuotedIdent reads a "quoted identifier". Snowflake's quoted identifiers
// preserve case and can contain any character except an unescaped ". The
// only escape is "" → literal ". The empty form "" is also a valid quoted
// identifier (with Str = "").
//
// Token.Str contains the unescaped content (no surrounding quotes).
// On unterminated quoted identifier, appends an unterminated-quoted-identifier
// LexError and emits a tokInvalid token.
func (l *Lexer) scanQuotedIdent() Token {
	start := l.start
	l.pos++ // consume opening "
	var sb strings.Builder
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == '"' {
			if l.peek(1) == '"' {
				// Doubled-quote escape: "" → literal "
				sb.WriteByte('"')
				l.pos += 2
				continue
			}
			// Closing quote.
			l.pos++
			return Token{
				Type: tokQuotedIdent,
				Str:  sb.String(),
				Loc:  ast.Loc{Start: start, End: l.pos},
			}
		}
		if ch == '\n' {
			l.errors = append(l.errors, LexError{
				Loc: ast.Loc{Start: start, End: l.pos},
				Msg: errUnterminatedQuoted,
			})
			return Token{Type: tokInvalid, Loc: ast.Loc{Start: start, End: l.pos}}
		}
		sb.WriteByte(ch)
		l.pos++
	}
	l.errors = append(l.errors, LexError{
		Loc: ast.Loc{Start: start, End: l.pos},
		Msg: errUnterminatedQuoted,
	})
	return Token{Type: tokInvalid, Loc: ast.Loc{Start: start, End: l.pos}}
}
```

- [ ] **Step 2: Verify compile and vet**

Run: `go build ./snowflake/parser/... && go vet ./snowflake/parser/...`
Expected: no output, exit 0.

---

## Task 8: Implement `scanDollar`

**Files:**
- Modify: `snowflake/parser/lexer.go`

- [ ] **Step 1: Replace the scanDollar stub**

In `snowflake/parser/lexer.go`, replace the `scanDollar` stub function with:

```go
// scanDollar handles tokens that begin with '$':
//   - $$...$$ → tokString with raw content (no escapes)
//   - $NAME   → tokVariable with Str = "NAME" (no leading $)
//   - $       → ASCII single-char token (Type = int('$'))
//
// The $$ form must be checked first so $ followed by $... is treated as a
// dollar string opener rather than a variable.
func (l *Lexer) scanDollar() Token {
	start := l.start

	// $$...$$ form
	if l.peek(1) == '$' {
		l.pos += 2 // consume opening $$
		contentStart := l.pos
		for l.pos < len(l.input)-1 {
			if l.input[l.pos] == '$' && l.input[l.pos+1] == '$' {
				content := l.input[contentStart:l.pos]
				l.pos += 2 // consume closing $$
				return Token{
					Type: tokString,
					Str:  content,
					Loc:  ast.Loc{Start: start, End: l.pos},
				}
			}
			l.pos++
		}
		// Unterminated $$ string. Advance to EOF.
		l.pos = len(l.input)
		l.errors = append(l.errors, LexError{
			Loc: ast.Loc{Start: start, End: l.pos},
			Msg: errUnterminatedDollar,
		})
		return Token{Type: tokInvalid, Loc: ast.Loc{Start: start, End: l.pos}}
	}

	// $NAME form: $ followed by [A-Z0-9_]+
	// Snowflake's ID2 rule: '$' [A-Z0-9_]*. Note that an empty $ (with no
	// follow-on identifier characters) falls through to the bare-$ case below.
	l.pos++ // consume $
	nameStart := l.pos
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '_' {
			l.pos++
			continue
		}
		break
	}
	if l.pos > nameStart {
		return Token{
			Type: tokVariable,
			Str:  l.input[nameStart:l.pos],
			Loc:  ast.Loc{Start: start, End: l.pos},
		}
	}

	// Bare $ — emit as ASCII single-char token.
	return Token{
		Type: int('$'),
		Loc:  ast.Loc{Start: start, End: l.pos},
	}
}
```

- [ ] **Step 2: Verify compile and vet**

Run: `go build ./snowflake/parser/... && go vet ./snowflake/parser/...`
Expected: no output, exit 0.

---

## Task 9: Implement `scanNumber`

**Files:**
- Modify: `snowflake/parser/lexer.go`

- [ ] **Step 1: Replace the scanNumber stub**

In `snowflake/parser/lexer.go`, replace the `scanNumber` stub function with:

```go
// scanNumber reads a numeric literal and returns one of three token kinds:
//   - tokInt   for digits-only (e.g. 42, 0)
//   - tokFloat for digits with a dot (e.g. 1.5, 1., .5)
//   - tokReal  for digits with an E exponent (e.g. 1e10, 1.5e-10)
//
// All three kinds populate Token.Str with the verbatim source text. tokInt
// additionally populates Token.Ival via strconv.ParseInt; if the integer
// overflows int64, the token is downgraded to tokFloat (still preserving
// Str) — this matches Snowflake's NUMBER(38, 0) arbitrary-precision behavior.
func (l *Lexer) scanNumber() Token {
	start := l.start
	isFloat := false

	// Leading dot form (.5).
	if l.input[l.pos] == '.' {
		isFloat = true
		l.pos++
		for l.pos < len(l.input) && l.input[l.pos] >= '0' && l.input[l.pos] <= '9' {
			l.pos++
		}
	} else {
		// Integer part.
		for l.pos < len(l.input) && l.input[l.pos] >= '0' && l.input[l.pos] <= '9' {
			l.pos++
		}
		// Optional decimal point.
		if l.pos < len(l.input) && l.input[l.pos] == '.' {
			isFloat = true
			l.pos++
			for l.pos < len(l.input) && l.input[l.pos] >= '0' && l.input[l.pos] <= '9' {
				l.pos++
			}
		}
	}

	// Optional exponent.
	isReal := false
	if l.pos < len(l.input) && (l.input[l.pos] == 'e' || l.input[l.pos] == 'E') {
		isReal = true
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

	if isReal {
		return Token{Type: tokReal, Str: text, Loc: loc}
	}
	if isFloat {
		return Token{Type: tokFloat, Str: text, Loc: loc}
	}

	// Integer. Try strconv.ParseInt; on overflow, downgrade to tokFloat.
	ival, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		// Overflow or other parse failure. Treat as tokFloat (Snowflake
		// NUMBER(38, 0) accepts arbitrary precision).
		return Token{Type: tokFloat, Str: text, Loc: loc}
	}
	return Token{Type: tokInt, Str: text, Ival: ival, Loc: loc}
}
```

- [ ] **Step 2: Add the `strconv` import**

In `snowflake/parser/lexer.go`, update the import block:

```go
import (
	"strconv"
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)
```

- [ ] **Step 3: Verify compile and vet**

Run: `go build ./snowflake/parser/... && go vet ./snowflake/parser/...`
Expected: no output, exit 0.

---

## Task 10: Implement `scanIdentOrKeyword`

**Files:**
- Modify: `snowflake/parser/lexer.go`

- [ ] **Step 1: Replace the scanIdentOrKeyword stub**

In `snowflake/parser/lexer.go`, replace the `scanIdentOrKeyword` stub function with:

```go
// scanIdentOrKeyword reads an identifier-shaped run and looks it up in the
// keyword map. Returns kw* if found (with Str = source text, preserving
// case) or tokIdent otherwise.
//
// Snowflake's ID rule is [A-Za-z_][A-Za-z0-9_@$]*. The legacy grammar's
// case-insensitive flag is honored by KeywordToken (which lowercases for
// lookup).
func (l *Lexer) scanIdentOrKeyword() Token {
	start := l.start
	for l.pos < len(l.input) && isIdentCont(l.input[l.pos]) {
		l.pos++
	}
	text := l.input[start:l.pos]
	loc := ast.Loc{Start: start, End: l.pos}
	if t, ok := KeywordToken(text); ok {
		return Token{Type: t, Str: text, Loc: loc}
	}
	return Token{Type: tokIdent, Str: text, Loc: loc}
}
```

- [ ] **Step 2: Verify compile and vet**

Run: `go build ./snowflake/parser/... && go vet ./snowflake/parser/...`
Expected: no output, exit 0.

---

## Task 11: Implement `scanOperatorOrPunct`

**Files:**
- Modify: `snowflake/parser/lexer.go`

- [ ] **Step 1: Replace the scanOperatorOrPunct stub**

In `snowflake/parser/lexer.go`, replace the `scanOperatorOrPunct` stub function with:

```go
// scanOperatorOrPunct handles single- and multi-char operators and
// punctuation. Multi-char operators (with their first byte and lookahead):
//
//   :: → tokDoubleColon
//   || → tokConcat
//   -> → tokArrow
//   ->>  → tokFlow (must be checked before -> via 2-byte lookahead)
//   => → tokAssoc
//   != → tokNotEq
//   <> → tokNotEq
//   <= → tokLessEq
//   >= → tokGreaterEq
//
// Single-char tokens use the ASCII byte value as their Type. Bytes that
// are not valid Snowflake operators or punctuation produce a tokInvalid
// token and an invalid-byte LexError.
func (l *Lexer) scanOperatorOrPunct(ch byte) Token {
	start := l.start
	switch ch {
	case ':':
		if l.peek(1) == ':' {
			l.pos += 2
			return Token{Type: tokDoubleColon, Loc: ast.Loc{Start: start, End: l.pos}}
		}
		l.pos++
		return Token{Type: int(':'), Loc: ast.Loc{Start: start, End: l.pos}}
	case '|':
		if l.peek(1) == '|' {
			l.pos += 2
			return Token{Type: tokConcat, Loc: ast.Loc{Start: start, End: l.pos}}
		}
		// Bare | is not a valid Snowflake operator.
		l.errors = append(l.errors, LexError{
			Loc: ast.Loc{Start: start, End: start + 1},
			Msg: errInvalidByte,
		})
		l.pos++
		return Token{Type: tokInvalid, Loc: ast.Loc{Start: start, End: l.pos}}
	case '-':
		if l.peek(1) == '>' {
			if l.peek(2) == '>' {
				l.pos += 3
				return Token{Type: tokFlow, Loc: ast.Loc{Start: start, End: l.pos}}
			}
			l.pos += 2
			return Token{Type: tokArrow, Loc: ast.Loc{Start: start, End: l.pos}}
		}
		l.pos++
		return Token{Type: int('-'), Loc: ast.Loc{Start: start, End: l.pos}}
	case '=':
		if l.peek(1) == '>' {
			l.pos += 2
			return Token{Type: tokAssoc, Loc: ast.Loc{Start: start, End: l.pos}}
		}
		l.pos++
		return Token{Type: int('='), Loc: ast.Loc{Start: start, End: l.pos}}
	case '!':
		if l.peek(1) == '=' {
			l.pos += 2
			return Token{Type: tokNotEq, Loc: ast.Loc{Start: start, End: l.pos}}
		}
		l.errors = append(l.errors, LexError{
			Loc: ast.Loc{Start: start, End: start + 1},
			Msg: errInvalidByte,
		})
		l.pos++
		return Token{Type: tokInvalid, Loc: ast.Loc{Start: start, End: l.pos}}
	case '<':
		if l.peek(1) == '=' {
			l.pos += 2
			return Token{Type: tokLessEq, Loc: ast.Loc{Start: start, End: l.pos}}
		}
		if l.peek(1) == '>' {
			l.pos += 2
			return Token{Type: tokNotEq, Loc: ast.Loc{Start: start, End: l.pos}}
		}
		l.pos++
		return Token{Type: int('<'), Loc: ast.Loc{Start: start, End: l.pos}}
	case '>':
		if l.peek(1) == '=' {
			l.pos += 2
			return Token{Type: tokGreaterEq, Loc: ast.Loc{Start: start, End: l.pos}}
		}
		l.pos++
		return Token{Type: int('>'), Loc: ast.Loc{Start: start, End: l.pos}}
	case '+', '*', '/', '%', '~', '(', ')', '[', ']', '{', '}', ',', ';', '.', '@':
		l.pos++
		return Token{Type: int(ch), Loc: ast.Loc{Start: start, End: l.pos}}
	}
	// Unknown byte.
	l.errors = append(l.errors, LexError{
		Loc: ast.Loc{Start: start, End: start + 1},
		Msg: errInvalidByte,
	})
	l.pos++
	return Token{Type: tokInvalid, Loc: ast.Loc{Start: start, End: l.pos}}
}
```

- [ ] **Step 2: Verify compile and vet**

Run: `go build ./snowflake/parser/... && go vet ./snowflake/parser/...`
Expected: no output, exit 0.

- [ ] **Step 3: Smoke-tokenize a tiny SQL string by hand**

Run a quick sanity check via a temporary file:

```bash
cat > /tmp/snowflake-smoke-test.go <<'EOF'
package main

import (
	"fmt"

	"github.com/bytebase/omni/snowflake/parser"
)

func main() {
	tokens, errs := parser.Tokenize("SELECT 1 + 2 FROM t WHERE x = 'hello';")
	for _, tok := range tokens {
		fmt.Printf("%d %q\n", tok.Type, tok.Str)
	}
	fmt.Printf("errors: %d\n", len(errs))
}
EOF
go run /tmp/snowflake-smoke-test.go
rm /tmp/snowflake-smoke-test.go
```

Expected output (token type numbers will vary by alphabetical position of keywords; the key check is `errors: 0`):
```
<some-int> "SELECT"
604 "1"          # tokInt = 604
43 ""           # '+' = 43
604 "2"
<some-int> "FROM"
604 ""          # tokIdent for 't'? actually 604+5=609 etc — verify via TokenName output if confused
...
errors: 0
```

The exact token types depend on the keyword alphabetical order, so don't strictly assert them. The key signals are:
- The program runs without panic
- `errors: 0`
- Some recognizable tokens like `"SELECT"`, `"FROM"`, `"WHERE"`, `"hello"` appear in the output

If `errors: > 0` or the program panics, stop and investigate.

---

## Task 12: Write `lexer_test.go` and `position_test.go`

**Files:**
- Modify: `snowflake/parser/lexer_test.go`
- Modify: `snowflake/parser/position_test.go`

Both files use the same fixture pattern: a table-driven helper that runs `Tokenize` and asserts the resulting token stream against an expected list of `(type, str)` pairs (`lexer_test.go`) or `(type, locStart, locEnd)` triples (`position_test.go`).

- [ ] **Step 1: Write `lexer_test.go` skeleton + helper**

Overwrite `snowflake/parser/lexer_test.go` with:

```go
package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// expect describes one expected token in the output stream. Loc is omitted
// here — position assertions live in position_test.go.
type expect struct {
	Type int
	Str  string
}

// runLexerCases drives a table-driven test. For each row, it tokenizes
// row.input and asserts the resulting non-EOF tokens match row.want by
// (Type, Str). Errors must be empty.
func runLexerCases(t *testing.T, cases []struct {
	name  string
	input string
	want  []expect
}) {
	t.Helper()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tokens, errs := Tokenize(c.input)
			if len(errs) != 0 {
				t.Fatalf("unexpected errors: %+v", errs)
			}
			// Strip the trailing EOF token before comparison.
			if len(tokens) == 0 || tokens[len(tokens)-1].Type != tokEOF {
				t.Fatalf("expected stream to end with tokEOF, got %+v", tokens)
			}
			tokens = tokens[:len(tokens)-1]
			if len(tokens) != len(c.want) {
				t.Fatalf("token count mismatch: got %d, want %d\n got=%+v\nwant=%+v",
					len(tokens), len(c.want), tokens, c.want)
			}
			for i, tok := range tokens {
				if tok.Type != c.want[i].Type {
					t.Errorf("token[%d] Type = %d (%s), want %d (%s)",
						i, tok.Type, TokenName(tok.Type), c.want[i].Type, TokenName(c.want[i].Type))
				}
				if tok.Str != c.want[i].Str {
					t.Errorf("token[%d] Str = %q, want %q", i, tok.Str, c.want[i].Str)
				}
			}
		})
	}
}

// silenceUnused keeps ast import live for tests that don't otherwise use it.
var _ = ast.NoLoc

func TestLexer_SingleCharOperators(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []expect
	}{
		{"plus", "+", []expect{{int('+'), ""}}},
		{"minus", "-", []expect{{int('-'), ""}}},
		{"star", "*", []expect{{int('*'), ""}}},
		{"slash", "/", []expect{{int('/'), ""}}},
		{"percent", "%", []expect{{int('%'), ""}}},
		{"tilde", "~", []expect{{int('~'), ""}}},
		{"lparen", "(", []expect{{int('('), ""}}},
		{"rparen", ")", []expect{{int(')'), ""}}},
		{"lbracket", "[", []expect{{int('['), ""}}},
		{"rbracket", "]", []expect{{int(']'), ""}}},
		{"lbrace", "{", []expect{{int('{'), ""}}},
		{"rbrace", "}", []expect{{int('}'), ""}}},
		{"comma", ",", []expect{{int(','), ""}}},
		{"semi", ";", []expect{{int(';'), ""}}},
		{"dot", ".", []expect{{int('.'), ""}}},
		{"at", "@", []expect{{int('@'), ""}}},
		{"colon", ":", []expect{{int(':'), ""}}},
		{"lt", "<", []expect{{int('<'), ""}}},
		{"gt", ">", []expect{{int('>'), ""}}},
		{"eq", "=", []expect{{int('='), ""}}},
	}
	runLexerCases(t, cases)
}

func TestLexer_MultiCharOperators(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []expect
	}{
		{"double_colon", "::", []expect{{tokDoubleColon, ""}}},
		{"concat", "||", []expect{{tokConcat, ""}}},
		{"arrow", "->", []expect{{tokArrow, ""}}},
		{"flow", "->>", []expect{{tokFlow, ""}}},
		{"assoc", "=>", []expect{{tokAssoc, ""}}},
		{"not_eq_bang", "!=", []expect{{tokNotEq, ""}}},
		{"not_eq_angle", "<>", []expect{{tokNotEq, ""}}},
		{"less_eq", "<=", []expect{{tokLessEq, ""}}},
		{"greater_eq", ">=", []expect{{tokGreaterEq, ""}}},
	}
	runLexerCases(t, cases)
}

func TestLexer_Identifiers(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []expect
	}{
		{"simple_lower", "abc", []expect{{tokIdent, "abc"}}},
		{"simple_upper", "ABC", []expect{{tokIdent, "ABC"}}},
		{"underscore_start", "_x", []expect{{tokIdent, "_x"}}},
		{"with_digits", "x123", []expect{{tokIdent, "x123"}}},
		{"with_at", "x@y", []expect{{tokIdent, "x@y"}}},
		{"with_dollar", "x$y", []expect{{tokIdent, "x$y"}}},
		{"quoted_simple", `"my_col"`, []expect{{tokQuotedIdent, "my_col"}}},
		{"quoted_with_space", `"my col"`, []expect{{tokQuotedIdent, "my col"}}},
		{"quoted_doubled_quote", `"a""b"`, []expect{{tokQuotedIdent, `a"b`}}},
		{"quoted_empty", `""`, []expect{{tokQuotedIdent, ""}}},
		{"variable_simple", "$x", []expect{{tokVariable, "x"}}},
		{"variable_numeric", "$1", []expect{{tokVariable, "1"}}},
		{"variable_long", "$ABC_123", []expect{{tokVariable, "ABC_123"}}},
		{"bare_dollar", "$", []expect{{int('$'), ""}}},
	}
	runLexerCases(t, cases)
}

func TestLexer_StringLiterals(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []expect
	}{
		{"simple", `'hello'`, []expect{{tokString, "hello"}}},
		{"empty", `''`, []expect{{tokString, ""}}},
		{"escape_n", `'a\nb'`, []expect{{tokString, "a\nb"}}},
		{"escape_t", `'a\tb'`, []expect{{tokString, "a\tb"}}},
		{"escape_backslash", `'a\\b'`, []expect{{tokString, `a\b`}}},
		{"escape_quote", `'a\'b'`, []expect{{tokString, "a'b"}}},
		{"doubled_quote", `'a''b'`, []expect{{tokString, "a'b"}}},
		{"dollar_simple", `$$hello$$`, []expect{{tokString, "hello"}}},
		{"dollar_empty", `$$$$`, []expect{{tokString, ""}}},
		{"dollar_with_quotes", `$$can't$$`, []expect{{tokString, "can't"}}},
	}
	runLexerCases(t, cases)
}

func TestLexer_HexBinaryLiterals(t *testing.T) {
	// X'...' is tokString with XPrefix=true. We can't use runLexerCases
	// for this since it doesn't compare XPrefix.
	tokens, errs := Tokenize("X'48656C6C6F'")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %+v", errs)
	}
	if len(tokens) != 2 || tokens[0].Type != tokString || !tokens[0].XPrefix {
		t.Fatalf("expected single tokString with XPrefix=true, got %+v", tokens)
	}
	if tokens[0].Str != "48656C6C6F" {
		t.Errorf("Str = %q, want %q", tokens[0].Str, "48656C6C6F")
	}

	// Lowercase x prefix should also work.
	tokens, errs = Tokenize("x'DEADBEEF'")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %+v", errs)
	}
	if !tokens[0].XPrefix {
		t.Errorf("lowercase x prefix not detected")
	}
}

func TestLexer_NumericLiterals(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []expect
	}{
		{"int_zero", "0", []expect{{tokInt, "0"}}},
		{"int_small", "42", []expect{{tokInt, "42"}}},
		{"int_large", "1234567890", []expect{{tokInt, "1234567890"}}},
		{"float_simple", "1.5", []expect{{tokFloat, "1.5"}}},
		{"float_trailing_dot", "1.", []expect{{tokFloat, "1."}}},
		{"float_leading_dot", ".5", []expect{{tokFloat, ".5"}}},
		{"real_simple", "1e10", []expect{{tokReal, "1e10"}}},
		{"real_pos_exp", "1e+10", []expect{{tokReal, "1e+10"}}},
		{"real_neg_exp", "1e-10", []expect{{tokReal, "1e-10"}}},
		{"real_float_base", "1.5e10", []expect{{tokReal, "1.5e10"}}},
		{"real_dot_base", ".5e10", []expect{{tokReal, ".5e10"}}},
	}
	runLexerCases(t, cases)
}

func TestLexer_NumericInt64Overflow(t *testing.T) {
	// Snowflake's NUMBER(38, 0) accepts arbitrary-precision integers.
	// When a literal exceeds int64, scanNumber downgrades to tokFloat.
	tokens, errs := Tokenize("99999999999999999999999999999999999999")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %+v", errs)
	}
	if len(tokens) != 2 {
		t.Fatalf("expected 2 tokens (literal + EOF), got %d: %+v", len(tokens), tokens)
	}
	if tokens[0].Type != tokFloat {
		t.Errorf("overflow literal Type = %d (%s), want tokFloat",
			tokens[0].Type, TokenName(tokens[0].Type))
	}
	if tokens[0].Str != "99999999999999999999999999999999999999" {
		t.Errorf("overflow literal Str = %q, want full source text", tokens[0].Str)
	}
	if tokens[0].Ival != 0 {
		t.Errorf("overflow literal Ival = %d, want 0 (unset)", tokens[0].Ival)
	}
}

func TestLexer_Keywords(t *testing.T) {
	// Spot-check a few keywords from each category. The full coverage is
	// in keyword_completeness_test.go (which checks every keyword from the
	// legacy .g4).
	cases := []struct {
		name  string
		input string
		want  []expect
	}{
		{"select_lower", "select", []expect{{kwSELECT, "select"}}},
		{"select_upper", "SELECT", []expect{{kwSELECT, "SELECT"}}},
		{"select_mixed", "SeLeCt", []expect{{kwSELECT, "SeLeCt"}}},
		{"from", "FROM", []expect{{kwFROM, "FROM"}}},
		{"where", "WHERE", []expect{{kwWHERE, "WHERE"}}},
		{"create", "CREATE", []expect{{kwCREATE, "CREATE"}}},
		{"table", "TABLE", []expect{{kwTABLE, "TABLE"}}},
		{"varchar", "VARCHAR", []expect{{kwVARCHAR, "VARCHAR"}}},
		{"variant", "VARIANT", []expect{{kwVARIANT, "VARIANT"}}},
		{"declare", "DECLARE", []expect{{kwDECLARE, "DECLARE"}}},
	}
	runLexerCases(t, cases)
}

func TestLexer_KeywordReservedClassification(t *testing.T) {
	// Spot-check IsReservedKeyword against known reserved/non-reserved entries.
	reserved := []string{"SELECT", "FROM", "WHERE", "CREATE", "DROP", "TABLE", "JOIN", "WITH"}
	for _, kw := range reserved {
		if !IsReservedKeyword(kw) {
			t.Errorf("IsReservedKeyword(%q) = false, want true", kw)
		}
	}
	nonReserved := []string{"ALERT", "VARCHAR", "VARIANT", "STAGE", "DECLARE"}
	for _, kw := range nonReserved {
		if IsReservedKeyword(kw) {
			t.Errorf("IsReservedKeyword(%q) = true, want false", kw)
		}
	}
}

func TestLexer_Comments(t *testing.T) {
	// All three comment forms should produce empty token streams (just EOF).
	cases := []string{
		"-- line comment\n",
		"// line comment\n",
		"/* block comment */",
		"/* nested /* inner */ outer */",
		"/* a */ /* b */",
		"-- comment 1\n-- comment 2\n",
	}
	for _, input := range cases {
		t.Run(input, func(t *testing.T) {
			tokens, errs := Tokenize(input)
			if len(errs) != 0 {
				t.Fatalf("unexpected errors: %+v", errs)
			}
			if len(tokens) != 1 || tokens[0].Type != tokEOF {
				t.Errorf("expected only EOF token, got %+v", tokens)
			}
		})
	}
}

func TestLexer_CommentsBetweenTokens(t *testing.T) {
	// Comments between two real tokens should be invisible.
	tokens, errs := Tokenize("SELECT /* hidden */ FROM")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %+v", errs)
	}
	if len(tokens) != 3 {
		t.Fatalf("expected 3 tokens (SELECT FROM EOF), got %d: %+v", len(tokens), tokens)
	}
	if tokens[0].Type != kwSELECT || tokens[1].Type != kwFROM {
		t.Errorf("expected [SELECT, FROM, EOF], got %+v", tokens)
	}
}
```

- [ ] **Step 2: Run `lexer_test.go`**

Run: `go test ./snowflake/parser/...`
Expected:
```
ok  	github.com/bytebase/omni/snowflake/parser	(some duration)
```

If any test fails, read the failure output, fix the implementation OR the test, and re-run until green.

- [ ] **Step 3: Write `position_test.go`**

Overwrite `snowflake/parser/position_test.go` with:

```go
package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

func TestPosition_TokenAtStart(t *testing.T) {
	tokens, _ := Tokenize("SELECT")
	if got := tokens[0].Loc; got != (ast.Loc{Start: 0, End: 6}) {
		t.Errorf("Loc = %+v, want {0, 6}", got)
	}
}

func TestPosition_TokenAfterWhitespace(t *testing.T) {
	tokens, _ := Tokenize("   SELECT")
	if got := tokens[0].Loc; got != (ast.Loc{Start: 3, End: 9}) {
		t.Errorf("Loc = %+v, want {3, 9}", got)
	}
}

func TestPosition_TokenAtEnd(t *testing.T) {
	tokens, _ := Tokenize("SELECT")
	if got := tokens[0].Loc.End; got != 6 {
		t.Errorf("Loc.End = %d, want 6", got)
	}
}

func TestPosition_MultiCharOperator(t *testing.T) {
	tokens, _ := Tokenize("a::b")
	// expect: ident, ::, ident
	if tokens[1].Type != tokDoubleColon {
		t.Fatalf("token[1] = %+v, expected ::", tokens[1])
	}
	if tokens[1].Loc != (ast.Loc{Start: 1, End: 3}) {
		t.Errorf("Loc = %+v, want {1, 3}", tokens[1].Loc)
	}
}

func TestPosition_StringLiteral(t *testing.T) {
	tokens, _ := Tokenize("'hello'")
	// Loc.Start is the opening quote, Loc.End is one past the closing quote.
	if tokens[0].Loc != (ast.Loc{Start: 0, End: 7}) {
		t.Errorf("Loc = %+v, want {0, 7}", tokens[0].Loc)
	}
}

func TestPosition_StringWithEscape(t *testing.T) {
	tokens, _ := Tokenize(`'a\nb'`)
	if tokens[0].Loc != (ast.Loc{Start: 0, End: 6}) {
		t.Errorf("Loc = %+v, want {0, 6}", tokens[0].Loc)
	}
}

func TestPosition_QuotedIdent(t *testing.T) {
	tokens, _ := Tokenize(`"my col"`)
	if tokens[0].Loc != (ast.Loc{Start: 0, End: 8}) {
		t.Errorf("Loc = %+v, want {0, 8}", tokens[0].Loc)
	}
}

func TestPosition_DollarString(t *testing.T) {
	tokens, _ := Tokenize("$$hello$$")
	if tokens[0].Loc != (ast.Loc{Start: 0, End: 9}) {
		t.Errorf("Loc = %+v, want {0, 9}", tokens[0].Loc)
	}
}

func TestPosition_Variable(t *testing.T) {
	tokens, _ := Tokenize("$abc")
	if tokens[0].Loc != (ast.Loc{Start: 0, End: 4}) {
		t.Errorf("Loc = %+v, want {0, 4}", tokens[0].Loc)
	}
}

func TestPosition_Number(t *testing.T) {
	cases := []struct {
		input string
		want  ast.Loc
	}{
		{"42", ast.Loc{Start: 0, End: 2}},
		{"1.5", ast.Loc{Start: 0, End: 3}},
		{".5", ast.Loc{Start: 0, End: 2}},
		{"1e10", ast.Loc{Start: 0, End: 4}},
		{"1.5e-10", ast.Loc{Start: 0, End: 7}},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			tokens, _ := Tokenize(c.input)
			if got := tokens[0].Loc; got != c.want {
				t.Errorf("Loc = %+v, want %+v", got, c.want)
			}
		})
	}
}

func TestPosition_HexLiteral(t *testing.T) {
	// X'48656C6C6F' is 14 bytes. Loc.Start should be the X (offset 0).
	tokens, _ := Tokenize("X'48656C6C6F'")
	if tokens[0].Loc != (ast.Loc{Start: 0, End: 13}) {
		t.Errorf("Loc = %+v, want {0, 13}", tokens[0].Loc)
	}
}

func TestPosition_EOFToken(t *testing.T) {
	tokens, _ := Tokenize("SELECT")
	last := tokens[len(tokens)-1]
	if last.Type != tokEOF {
		t.Fatalf("last token Type = %d, want tokEOF", last.Type)
	}
	if last.Loc != (ast.Loc{Start: 6, End: 6}) {
		t.Errorf("EOF Loc = %+v, want {6, 6}", last.Loc)
	}
}

func TestPosition_EOFTokenEmptyInput(t *testing.T) {
	tokens, _ := Tokenize("")
	if len(tokens) != 1 {
		t.Fatalf("expected 1 token (just EOF), got %d", len(tokens))
	}
	if tokens[0].Loc != (ast.Loc{Start: 0, End: 0}) {
		t.Errorf("EOF Loc = %+v, want {0, 0}", tokens[0].Loc)
	}
}
```

- [ ] **Step 4: Run all tests**

Run: `go test ./snowflake/parser/...`
Expected:
```
ok  	github.com/bytebase/omni/snowflake/parser	(some duration)
```

---

## Task 13: Write `edge_cases_test.go` and `error_recovery_test.go`

**Files:**
- Modify: `snowflake/parser/edge_cases_test.go`
- Modify: `snowflake/parser/error_recovery_test.go`

- [ ] **Step 1: Write `edge_cases_test.go`**

Overwrite `snowflake/parser/edge_cases_test.go` with:

```go
package parser

import (
	"strings"
	"testing"
)

func TestEdge_NestedBlockComments(t *testing.T) {
	// 3-deep nesting should produce no tokens and no errors.
	input := "/* a /* b /* c */ */ */"
	tokens, errs := Tokenize(input)
	if len(errs) != 0 {
		t.Errorf("unexpected errors: %+v", errs)
	}
	if len(tokens) != 1 || tokens[0].Type != tokEOF {
		t.Errorf("expected only EOF, got %+v", tokens)
	}
}

func TestEdge_UnterminatedBlockComment(t *testing.T) {
	tokens, errs := Tokenize("/* unterminated")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %+v", len(errs), errs)
	}
	if !strings.Contains(errs[0].Msg, "unterminated block comment") {
		t.Errorf("error message = %q, want substring 'unterminated block comment'", errs[0].Msg)
	}
	if len(tokens) != 1 || tokens[0].Type != tokEOF {
		t.Errorf("expected only EOF token (comment is hidden), got %+v", tokens)
	}
}

func TestEdge_MismatchedBlockCommentDepth(t *testing.T) {
	// /* /* */  has depth 2 → 1 → 1 (not back to 0). Unterminated.
	tokens, errs := Tokenize("/* /* */")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %+v", len(errs), errs)
	}
	if len(tokens) != 1 || tokens[0].Type != tokEOF {
		t.Errorf("expected only EOF token, got %+v", tokens)
	}
}

func TestEdge_UnterminatedString(t *testing.T) {
	tokens, errs := Tokenize("'unterminated")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %+v", len(errs), errs)
	}
	if !strings.Contains(errs[0].Msg, "unterminated string literal") {
		t.Errorf("error = %q, want unterminated string", errs[0].Msg)
	}
	// Token stream: tokInvalid, EOF
	if len(tokens) != 2 || tokens[0].Type != tokInvalid {
		t.Errorf("expected [tokInvalid, EOF], got %+v", tokens)
	}
}

func TestEdge_UnterminatedDollarString(t *testing.T) {
	tokens, errs := Tokenize("$$ unterminated")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %+v", len(errs), errs)
	}
	if !strings.Contains(errs[0].Msg, "unterminated $$ string literal") {
		t.Errorf("error = %q, want unterminated $$", errs[0].Msg)
	}
	if len(tokens) != 2 || tokens[0].Type != tokInvalid {
		t.Errorf("expected [tokInvalid, EOF], got %+v", tokens)
	}
}

func TestEdge_UnterminatedQuotedIdent(t *testing.T) {
	tokens, errs := Tokenize(`"unterminated`)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %+v", len(errs), errs)
	}
	if !strings.Contains(errs[0].Msg, "unterminated quoted identifier") {
		t.Errorf("error = %q, want unterminated quoted identifier", errs[0].Msg)
	}
	if len(tokens) != 2 || tokens[0].Type != tokInvalid {
		t.Errorf("expected [tokInvalid, EOF], got %+v", tokens)
	}
}

func TestEdge_NumericLeadingDot(t *testing.T) {
	tokens, _ := Tokenize(".5")
	if tokens[0].Type != tokFloat || tokens[0].Str != ".5" {
		t.Errorf("got %+v, want tokFloat .5", tokens[0])
	}
}

func TestEdge_NumericTrailingDot(t *testing.T) {
	tokens, _ := Tokenize("1.")
	if tokens[0].Type != tokFloat || tokens[0].Str != "1." {
		t.Errorf("got %+v, want tokFloat 1.", tokens[0])
	}
}

func TestEdge_NumericExponentVariants(t *testing.T) {
	for _, input := range []string{"1e10", "1E10", "1e+10", "1e-10", ".5e10", "1.5e10"} {
		t.Run(input, func(t *testing.T) {
			tokens, _ := Tokenize(input)
			if tokens[0].Type != tokReal {
				t.Errorf("got %+v, want tokReal", tokens[0])
			}
		})
	}
}

func TestEdge_AmbiguousOperators(t *testing.T) {
	cases := []struct {
		input    string
		opType   int
	}{
		{"a::b", tokDoubleColon},
		{"a||b", tokConcat},
		{"a->b", tokArrow},
		{"a->>b", tokFlow},
		{"a=>b", tokAssoc},
		{"a!=b", tokNotEq},
		{"a<>b", tokNotEq},
		{"a<=b", tokLessEq},
		{"a>=b", tokGreaterEq},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			tokens, errs := Tokenize(c.input)
			if len(errs) != 0 {
				t.Fatalf("unexpected errors: %+v", errs)
			}
			// Expected: [ident, op, ident, EOF]
			if len(tokens) != 4 {
				t.Fatalf("expected 4 tokens, got %d: %+v", len(tokens), tokens)
			}
			if tokens[1].Type != c.opType {
				t.Errorf("op token = %d (%s), want %d (%s)",
					tokens[1].Type, TokenName(tokens[1].Type),
					c.opType, TokenName(c.opType))
			}
		})
	}
}

func TestEdge_UnicodeInQuotedIdent(t *testing.T) {
	tokens, _ := Tokenize(`"héllo"`)
	if tokens[0].Type != tokQuotedIdent || tokens[0].Str != "héllo" {
		t.Errorf("got %+v, want tokQuotedIdent héllo", tokens[0])
	}
}

func TestEdge_UnicodeInString(t *testing.T) {
	tokens, _ := Tokenize(`'café'`)
	if tokens[0].Type != tokString || tokens[0].Str != "café" {
		t.Errorf("got %+v, want tokString café", tokens[0])
	}
}

func TestEdge_DoubledQuoteInQuotedIdent(t *testing.T) {
	tokens, _ := Tokenize(`"a""b"`)
	if tokens[0].Type != tokQuotedIdent || tokens[0].Str != `a"b` {
		t.Errorf("got %+v, want tokQuotedIdent a\"b", tokens[0])
	}
}

func TestEdge_DoubledQuoteInString(t *testing.T) {
	tokens, _ := Tokenize(`'a''b'`)
	if tokens[0].Type != tokString || tokens[0].Str != "a'b" {
		t.Errorf("got %+v, want tokString a'b", tokens[0])
	}
}

func TestEdge_EmptyQuotedIdent(t *testing.T) {
	tokens, _ := Tokenize(`""`)
	if tokens[0].Type != tokQuotedIdent || tokens[0].Str != "" {
		t.Errorf("got %+v, want tokQuotedIdent (empty)", tokens[0])
	}
}

func TestEdge_BareDollar(t *testing.T) {
	tokens, _ := Tokenize("$")
	if tokens[0].Type != int('$') {
		t.Errorf("got %+v, want int('$')", tokens[0])
	}
}

func TestEdge_BarePipe(t *testing.T) {
	// Bare | is invalid in Snowflake (only || is a token).
	tokens, errs := Tokenize("|")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %+v", len(errs), errs)
	}
	if tokens[0].Type != tokInvalid {
		t.Errorf("got %+v, want tokInvalid", tokens[0])
	}
}
```

- [ ] **Step 2: Write `error_recovery_test.go`**

Overwrite `snowflake/parser/error_recovery_test.go` with:

```go
package parser

import (
	"testing"
)

func TestRecovery_TwoUnterminatedStrings(t *testing.T) {
	// 'unterm1' is actually terminated; 'unterm2 is not. So this only
	// produces one real error. Use a multiline form for two errors.
	input := "'unterm1\n'unterm2"
	tokens, errs := Tokenize(input)
	if len(errs) != 2 {
		t.Errorf("expected 2 errors, got %d: %+v", len(errs), errs)
	}
	if len(tokens) < 1 {
		t.Errorf("expected at least 1 token, got 0")
	}
}

func TestRecovery_KeywordBetweenErrors(t *testing.T) {
	// SELECT preserved between two invalid bytes.
	input := "| SELECT |"
	tokens, errs := Tokenize(input)
	if len(errs) != 2 {
		t.Errorf("expected 2 errors, got %d: %+v", len(errs), errs)
	}
	// Stream should contain at least: tokInvalid, kwSELECT, tokInvalid, EOF
	hasSelect := false
	for _, tok := range tokens {
		if tok.Type == kwSELECT {
			hasSelect = true
			break
		}
	}
	if !hasSelect {
		t.Errorf("expected kwSELECT in stream, got %+v", tokens)
	}
}

func TestRecovery_ContinuesPastInvalidByte(t *testing.T) {
	// A NUL byte in the middle of input should be reported as one error
	// and lexing should continue past it.
	input := "SELECT \x00 FROM"
	tokens, errs := Tokenize(input)
	if len(errs) != 1 {
		t.Errorf("expected 1 error, got %d: %+v", len(errs), errs)
	}
	hasSelect := false
	hasFrom := false
	for _, tok := range tokens {
		if tok.Type == kwSELECT {
			hasSelect = true
		}
		if tok.Type == kwFROM {
			hasFrom = true
		}
	}
	if !hasSelect || !hasFrom {
		t.Errorf("expected SELECT and FROM preserved, got %+v", tokens)
	}
}

func TestRecovery_ReachesEOF(t *testing.T) {
	// Lexer must always reach EOF, even with errors.
	input := "| | | |"
	tokens, _ := Tokenize(input)
	if tokens[len(tokens)-1].Type != tokEOF {
		t.Errorf("did not reach EOF; last token = %+v", tokens[len(tokens)-1])
	}
}
```

- [ ] **Step 3: Run all tests**

Run: `go test ./snowflake/parser/...`
Expected:
```
ok  	github.com/bytebase/omni/snowflake/parser	(some duration)
```

If any test fails, fix the implementation OR the test, re-run until green.

---

## Task 14: Write `legacy_corpus_test.go` and `keyword_completeness_test.go`

**Files:**
- Modify: `snowflake/parser/legacy_corpus_test.go`
- Modify: `snowflake/parser/keyword_completeness_test.go`

- [ ] **Step 1: Write `legacy_corpus_test.go`**

Overwrite `snowflake/parser/legacy_corpus_test.go` with:

```go
package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestLegacyCorpus tokenizes every .sql file in testdata/legacy/ and asserts
// no lex errors and no tokInvalid tokens. This is the regression baseline:
// any failure here means the omni snowflake lexer cannot tokenize something
// the legacy ANTLR4 parser previously handled.
func TestLegacyCorpus(t *testing.T) {
	corpusDir := "testdata/legacy"
	entries, err := os.ReadDir(corpusDir)
	if err != nil {
		t.Fatalf("failed to read corpus dir %s: %v", corpusDir, err)
	}

	sqlCount := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		sqlCount++

		t.Run(entry.Name(), func(t *testing.T) {
			path := filepath.Join(corpusDir, entry.Name())
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read %s: %v", path, err)
			}
			tokens, errs := Tokenize(string(data))

			if len(errs) > 0 {
				t.Errorf("%d lex errors:", len(errs))
				for _, e := range errs {
					t.Errorf("  loc=%+v msg=%q", e.Loc, e.Msg)
				}
			}

			invalidCount := 0
			for _, tok := range tokens {
				if tok.Type == tokInvalid {
					invalidCount++
				}
			}
			if invalidCount > 0 {
				t.Errorf("%d tokInvalid tokens emitted", invalidCount)
			}

			// Sanity check: stream must end with EOF.
			if len(tokens) == 0 || tokens[len(tokens)-1].Type != tokEOF {
				t.Errorf("stream did not end with tokEOF")
			}

			// Sanity check: the last non-EOF token's End should be at or
			// beyond the file's last semicolon byte position. This catches
			// the lexer halting prematurely on some content.
			lastSemi := strings.LastIndexByte(string(data), ';')
			if lastSemi >= 0 && len(tokens) >= 2 {
				lastNonEOF := tokens[len(tokens)-2]
				if lastNonEOF.Loc.End < lastSemi {
					t.Errorf("last non-EOF token Loc.End=%d < last semicolon position %d (lexer stopped early)",
						lastNonEOF.Loc.End, lastSemi)
				}
			}
		})
	}

	if sqlCount == 0 {
		t.Errorf("found 0 .sql files in %s — corpus missing?", corpusDir)
	}
	if sqlCount != 27 {
		t.Logf("note: found %d .sql files (expected 27 from the legacy lift)", sqlCount)
	}
}
```

- [ ] **Step 2: Write `keyword_completeness_test.go`**

Overwrite `snowflake/parser/keyword_completeness_test.go` with:

```go
package parser

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

// legacyLexerPath is the absolute path to the legacy SnowflakeLexer.g4
// grammar file used as the source of truth for keyword drift detection.
// If the file is missing (e.g. on a CI machine without the legacy parser
// checkout), the tests in this file are skipped.
const legacyLexerPath = "/Users/h3n4l/OpenSource/parser/snowflake/SnowflakeLexer.g4"

// legacyReservedScript is the absolute path to the Python script that
// defines the reserved keyword set in the legacy parser.
const legacyReservedScript = "/Users/h3n4l/OpenSource/parser/snowflake/build_id_contains_non_reserved_keywords.py"

// TestKeywordCompleteness asserts that every single-word keyword token
// defined in SnowflakeLexer.g4 has a corresponding entry in keywordMap.
//
// Skips with t.Skip if the legacy grammar file is missing.
func TestKeywordCompleteness(t *testing.T) {
	data, err := os.ReadFile(legacyLexerPath)
	if err != nil {
		t.Skipf("legacy grammar not available at %s: %v", legacyLexerPath, err)
	}

	// Match lines like:  ABORT : 'ABORT';
	// Captures the literal between the single quotes.
	re := regexp.MustCompile(`(?m)^[A-Z_][A-Z0-9_]*\s*:\s*'([A-Z][A-Z0-9_]*)'\s*;`)
	matches := re.FindAllStringSubmatch(string(data), -1)
	if len(matches) == 0 {
		t.Fatalf("no keyword tokens extracted from %s — regex broken?", legacyLexerPath)
	}

	missing := []string{}
	for _, m := range matches {
		literal := m[1] // the captured keyword text
		lower := strings.ToLower(literal)
		if _, ok := keywordMap[lower]; !ok {
			missing = append(missing, literal)
		}
	}

	if len(missing) > 0 {
		t.Errorf("%d keywords missing from keywordMap (sample of first 10): %v",
			len(missing), missing[:min(len(missing), 10)])
	}

	t.Logf("checked %d keywords from legacy grammar; %d missing", len(matches), len(missing))
}

// TestReservedKeywordsCompleteness asserts that every keyword in the legacy
// build_id_contains_non_reserved_keywords.py snowflake_reserved_keyword dict
// is present in keywordReserved, EXCEPT for a known set of phantoms that
// the legacy Python script reserves but the legacy SnowflakeLexer.g4 does
// not actually emit as token rules. The omni lexer cannot mark a keyword
// as reserved if it has no kw* token, so these phantoms are deliberately
// excluded from keywordReserved (see the comment in keywords.go).
//
// Skips with t.Skip if the script is missing.
func TestReservedKeywordsCompleteness(t *testing.T) {
	data, err := os.ReadFile(legacyReservedScript)
	if err != nil {
		t.Skipf("legacy reserved-keyword script not available at %s: %v", legacyReservedScript, err)
	}

	// knownPhantoms are keyword names that the Python reserved-set file
	// declares but the legacy SnowflakeLexer.g4 does not emit as a token
	// rule (two are commented out, seven are entirely absent). Verified
	// 2026-04-07. If a future grammar update adds any of these, drop the
	// entry from this set and add it to keywordReserved in keywords.go.
	knownPhantoms := map[string]bool{
		"CURRENT_USER":   true, // commented out in lexer grammar
		"FOLLOWING":      true, // not in lexer grammar
		"GSCLUSTER":      true, // not in lexer grammar
		"ISSUE":          true, // not in lexer grammar
		"LOCALTIME":      true, // not in lexer grammar
		"LOCALTIMESTAMP": true, // not in lexer grammar
		"REGEXP":         true, // not in lexer grammar
		"TRIGGER":        true, // commented out in lexer grammar
		"WHENEVER":       true, // not in lexer grammar
	}

	// Match lines like:    "ACCOUNT": True,
	// inside the snowflake_reserved_keyword dict.
	re := regexp.MustCompile(`(?m)^\s*"([A-Z_][A-Z0-9_]*)"\s*:\s*True\s*,?`)
	matches := re.FindAllStringSubmatch(string(data), -1)
	if len(matches) == 0 {
		t.Fatalf("no reserved keywords extracted from %s — regex broken?", legacyReservedScript)
	}

	missing := []string{}
	skipped := 0
	for _, m := range matches {
		kw := m[1]
		if knownPhantoms[kw] {
			skipped++
			continue
		}
		if !IsReservedKeyword(kw) {
			missing = append(missing, kw)
		}
	}

	if len(missing) > 0 {
		t.Errorf("%d reserved keywords missing from keywordReserved (not in known-phantom set): %v",
			len(missing), missing)
	}

	t.Logf("checked %d reserved keywords from legacy script; %d skipped as known phantoms; %d missing",
		len(matches), skipped, len(missing))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
```

- [ ] **Step 3: Run all tests**

Run: `go test ./snowflake/parser/...`
Expected:
```
ok  	github.com/bytebase/omni/snowflake/parser	(some duration)
```

If `TestKeywordCompleteness` reports missing keywords, the extraction in Task 3 was incomplete — investigate and fix `tokens.go` + `keywords.go`. If `TestReservedKeywordsCompleteness` reports missing reserved keywords, the hand-typed `keywordReserved` block in `keywords.go` is missing entries — add them.

If `TestLegacyCorpus` reports failures on specific files, the lexer is not handling some legacy syntax correctly. Read the failing file, find the un-tokenizable substring, and patch the lexer.

- [ ] **Step 4: Run with verbose output to see test counts**

Run: `go test -v ./snowflake/parser/... 2>&1 | tail -30`
Expected: All test functions reported as PASS. The corpus test will list 27 sub-tests (one per .sql file), all PASS. The keyword completeness test will log its check counts.

---

## Task 15: Final acceptance criteria sweep

This task runs every acceptance criterion from the spec back-to-back to confirm F2 is complete.

- [ ] **Step 1: Build**

Run: `go build ./snowflake/parser/...`
Expected: no output, exit 0.

- [ ] **Step 2: Vet**

Run: `go vet ./snowflake/parser/...`
Expected: no output, exit 0.

- [ ] **Step 3: Gofmt**

Run: `gofmt -l snowflake/parser/`
Expected: no output (no files need reformatting).

- [ ] **Step 4: Test**

Run: `go test ./snowflake/parser/...`
Expected:
```
ok  	github.com/bytebase/omni/snowflake/parser	(some duration)
```

- [ ] **Step 5: Test with verbose output**

Run: `go test -v ./snowflake/parser/... 2>&1 | tail -50`
Expected: see test summaries — the corpus test should list all 27 sub-tests as PASS, the keyword completeness test should log "checked N keywords from legacy grammar; 0 missing" and "checked 99 reserved keywords from legacy script; 0 missing".

- [ ] **Step 6: Package isolation**

Run: `go build ./snowflake/...`
Expected: no output, exit 0. Confirms `snowflake/parser` builds without depending on any Tier 1+ DAG node.

- [ ] **Step 7: List all created files**

Run: `git status snowflake/parser docs/superpowers/plans`
Expected output should show these untracked items (plus the design doc which was committed in d4ad94f):
```
?? snowflake/parser/edge_cases_test.go
?? snowflake/parser/error_recovery_test.go
?? snowflake/parser/errors.go
?? snowflake/parser/keyword_completeness_test.go
?? snowflake/parser/keywords.go
?? snowflake/parser/legacy_corpus_test.go
?? snowflake/parser/lexer.go
?? snowflake/parser/lexer_test.go
?? snowflake/parser/position_test.go
?? snowflake/parser/tokens.go
?? docs/superpowers/plans/2026-04-07-snowflake-lexer.md
```

- [ ] **Step 8: STOP and ask the user for review**

Do NOT commit. Output a summary message:

> F2 (snowflake/parser lexer) implementation complete.
>
> All 15 tasks done. Acceptance criteria green:
> - `go build ./snowflake/parser/...` ✓
> - `go vet ./snowflake/parser/...` ✓
> - `gofmt -l snowflake/parser/` clean ✓
> - `go test ./snowflake/parser/...` ✓
> - 27 legacy corpus files tokenize cleanly ✓
> - keyword completeness test green ✓
> - `go build ./snowflake/...` (isolation) ✓
>
> Files created (10 code + 1 plan):
> - snowflake/parser/errors.go
> - snowflake/parser/tokens.go
> - snowflake/parser/keywords.go
> - snowflake/parser/lexer.go
> - snowflake/parser/lexer_test.go
> - snowflake/parser/position_test.go
> - snowflake/parser/edge_cases_test.go
> - snowflake/parser/error_recovery_test.go
> - snowflake/parser/legacy_corpus_test.go
> - snowflake/parser/keyword_completeness_test.go
> - docs/superpowers/plans/2026-04-07-snowflake-lexer.md
>
> Please review the diff. After your approval I will commit and then run finishing-a-development-branch.

---

## Spec Coverage Checklist

After running through the plan once, verify each spec requirement is covered by a task:

| Spec section | Covered by |
|--------------|-----------|
| Architecture (mysql-style flat lexer, drop pg state machine) | Task 5 (lexer.go core) |
| File layout (4 production + 6 test files) | Task 1 (scaffold) + Tasks 2–14 (fill) |
| `errors.go` types | Task 2 |
| `tokens.go` operator constants | Task 3 step 1 |
| `tokens.go` keyword constants (~881) | Task 3 steps 2–4 |
| `tokens.go` Token struct | Task 3 step 1 |
| `tokens.go` TokenName function | Task 3 step 1 |
| `keywords.go` keywordMap | Task 4 steps 2–3 |
| `keywords.go` keywordReserved | Task 4 step 1 |
| `keywords.go` KeywordToken / IsReservedKeyword | Task 4 step 1 |
| `lexer.go` Lexer struct + NextToken dispatch | Task 5 |
| `lexer.go` skipWhitespaceAndComments + scanBlockComment | Task 5 |
| `lexer.go` scanString | Task 6 |
| `lexer.go` scanQuotedIdent | Task 7 |
| `lexer.go` scanDollar | Task 8 |
| `lexer.go` scanNumber (with int64 overflow downgrade) | Task 9 |
| `lexer.go` scanIdentOrKeyword | Task 10 |
| `lexer.go` scanOperatorOrPunct | Task 11 |
| Token recognition tests | Task 12 (lexer_test.go) |
| Position math tests | Task 12 (position_test.go) |
| Edge case tests | Task 13 (edge_cases_test.go) |
| Error recovery tests | Task 13 (error_recovery_test.go) |
| Legacy corpus regression | Task 14 (legacy_corpus_test.go) |
| Keyword completeness drift detector | Task 14 (keyword_completeness_test.go) |
| Acceptance criteria 1–7 | Task 15 |
