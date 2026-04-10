# PartiQL Parser Foundation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement the foundation layer of the hand-written PartiQL recursive-descent parser: Parser struct, token buffer, precedence ladder, primary-expression dispatch, path steps, literals, and parseType. Every non-foundation feature is stubbed with a "deferred to DAG node N" error.

**Architecture:** 5-file split in `partiql/parser/` — `parser.go` (machinery + parseType), `expr.go` (precedence ladder), `exprprimary.go` (primary-expression dispatch + exprTerm), `path.go` (path-step chains), `literals.go` (literal expressions). Tests live in `parser_test.go` + filesystem goldens under `testdata/parser-foundation/`. Upstream dependencies (ast-core and lexer) are complete; this plan only touches `partiql/parser/` and new files.

**Tech Stack:** Go 1.21+, `partiql/ast` (sealed sub-interfaces: Node, StmtNode, ExprNode, TableExpr, PathStep, TypeName, PatternNode), `partiql/parser` Lexer (302 tok constants, `Lexer.Err` first-error-and-stop model), `ast.NodeToString` from ast-core's `outfuncs.go` for golden dumps.

**Reference precedents:**
- `cosmosdb/parser/parser.go` — closest precedent (hand-written recursive descent with path steps)
- `cosmosdb/parser/expr.go` — expression parsing (uses Pratt climber — we're NOT using that, but the token-consumption idioms are the same)
- `cosmosdb/parser/golden_test.go` — golden file pattern we're mirroring

**Spec:** `docs/superpowers/specs/2026-04-09-partiql-parser-foundation-design.md` — read this first. It captures the 6 design decisions D1-D6, the stub audit, and the test strategy.

---

## File Structure

| File | Responsibility | Lines (approx) |
|------|---------------|----------------|
| `partiql/parser/parser.go` | Parser struct, ParseError, NewParser, token buffer helpers (advance/peek/peekNext/match/expect), parseSymbolPrimitive, parseVarRef (upgraded in Task 10), parseType, `deferredFeature` helper | ~300 |
| `partiql/parser/literals.go` | `parseLiteral` switch for 6 real literal forms + DATE/TIME stubs | ~140 |
| `partiql/parser/path.go` | `parsePathSteps`, `parsePathStep`, `isPathStepStart` | ~130 |
| `partiql/parser/exprprimary.go` | `parsePrimary` + `parsePrimaryBase` dispatch (16 grammar alternatives), exprTerm cases (`parseParenExpr`, `parseArrayLit`, `parseBagLit`, `parseTupleLit`, `parseParamRef`), `deferredFeature` call sites for every non-foundation alternative | ~280 |
| `partiql/parser/expr.go` | Precedence ladder: `parseExpr` → `parseBagOp` → `parseSelectExpr` → `parseOr` → `parseAnd` → `parseNot` → `parsePredicate` → `parseMathOp00/01/02` → `parseValueExpr` | ~380 |
| `partiql/parser/parser_test.go` | TestParser_Machinery, TestParser_Goldens, TestParser_AWSCorpus, TestParser_Errors, TestParseType | ~360 |
| `partiql/parser/testdata/parser-foundation/*.partiql` | ~43 input files | — |
| `partiql/parser/testdata/parser-foundation/*.golden` | ~43 golden output files | — |

**Key invariants from upstream:**

1. **Token:** `type Token struct { Type int; Str string; Loc ast.Loc }` — `ast.Loc` is typed (not int).
2. **Lexer:** `NewLexer(input string) *Lexer`, `l.Next() Token`, first-error-and-stop via `l.Err error`. After `l.Err` is set, all `Next()` calls return `tokEOF`.
3. **AST Loc:** `type Loc struct { Start, End int }` half-open byte range. `{-1, -1}` = unknown.
4. **ExprNode:** sealed interface; `exprNode()` method gates membership.

---

## Conventions Used Throughout the Plan

### Error construction

Every `*ParseError` uses this shape:

```go
return nil, &ParseError{
    Message: fmt.Sprintf("expected %s, got %q", tokenName(wantType), p.cur.Str),
    Loc:     ast.Loc{Start: p.cur.Loc.Start, End: p.cur.Loc.Start},
}
```

### Loc merge pattern

When building a node that spans multiple tokens, take `Start` from the leftmost token and `End` from `p.prev.Loc.End` (the most recently consumed token) or from a child node's `GetLoc().End`:

```go
&ast.BinaryExpr{
    Op:    ast.BinOpAdd,
    Left:  left,
    Right: right,
    Loc:   ast.Loc{Start: left.GetLoc().Start, End: right.GetLoc().End},
}
```

### Homebrew Go gofmt quirk

The Homebrew Go 1.26.0 `gofmt` rewrites consecutive apostrophes (`''`) in `//` line comments to U+201D. When writing doc comments that reference SQL doubled-quote escapes, either (a) use prose phrasing ("two consecutive apostrophes") or (b) use `/* */` block comments. After every commit, verify with: `grep -n $'\xe2\x80\x9d' partiql/parser/*.go` — should return empty.

### Task 5 invariant (path attachment)

`parsePrimary()` returns a base expression and then attaches any trailing path steps via `parsePathSteps`. Scan helpers inside `parsePrimaryBase` must NOT attempt to parse path steps themselves — the wrapping happens exactly once at the `parsePrimary` level.

---

## Task 1: parser.go scaffold + TestParser_Machinery

**Files:**
- Create: `partiql/parser/parser.go`
- Create: `partiql/parser/parser_test.go`

This task introduces the Parser struct, ParseError type, constructor, token buffer helpers, `parseSymbolPrimitive`, `parseVarRef` (bare form), and a `deferredFeature` helper. No expression parsing yet. Tests verify the helpers work correctly via direct method calls.

- [ ] **Step 1: Create `partiql/parser/parser.go`**

```go
// Package parser provides a hand-written recursive-descent parser for
// PartiQL. This file holds the Parser struct, token buffer helpers,
// error construction, and the shared utility parsers
// (parseSymbolPrimitive, parseVarRef, parseType).
//
// Expression parsing is split across expr.go (precedence ladder) and
// exprprimary.go (primary-expression dispatch). Literal parsing lives
// in literals.go; path-step chains in path.go.
//
// Upstream: this parser consumes tokens produced by Lexer (lexer.go,
// token.go, keywords.go). The lexer uses first-error-and-stop via
// Lexer.Err; the parser surfaces lexer errors as ParseError.
package parser

import (
	"fmt"

	"github.com/bytebase/omni/partiql/ast"
)

// ParseError represents a syntax error during parsing.
//
// Loc points at the offending token (Start = End = token start). The
// parser uses fail-fast semantics — the first ParseError aborts the
// parse and is returned to the caller.
type ParseError struct {
	Message string
	Loc     ast.Loc
}

// Error renders a human-readable message including the byte position.
func (e *ParseError) Error() string {
	return fmt.Sprintf("syntax error at position %d: %s", e.Loc.Start, e.Message)
}

// Parser is the recursive-descent parser for PartiQL.
//
// Fields:
//   - lexer:   the token source
//   - cur:     current lookahead token (the "peek" slot)
//   - prev:    most recently consumed token, used to compute end Locs
//   - nextBuf: one-token lookahead buffer for peekNext()
//   - hasNext: true when nextBuf holds a valid token
type Parser struct {
	lexer   *Lexer
	cur     Token
	prev    Token
	nextBuf Token
	hasNext bool
}

// NewParser creates a Parser that reads from the given input string.
// The first token is primed into cur before the constructor returns.
func NewParser(input string) *Parser {
	p := &Parser{lexer: NewLexer(input)}
	p.advance() // prime cur with the first token
	return p
}

// advance moves cur → prev and reads the next token from the lexer.
// If the lexer has already errored, the returned token is tokEOF with
// the lexer error message embedded in Str so the caller can surface it.
func (p *Parser) advance() {
	p.prev = p.cur
	if p.hasNext {
		p.cur = p.nextBuf
		p.hasNext = false
	} else {
		p.cur = p.lexer.Next()
	}
}

// peek returns the current token without consuming it.
func (p *Parser) peek() Token {
	return p.cur
}

// peekNext returns the token AFTER cur without consuming either cur or
// the returned token. Uses the nextBuf slot; subsequent calls are
// idempotent until advance() is called.
func (p *Parser) peekNext() Token {
	if !p.hasNext {
		p.nextBuf = p.lexer.Next()
		p.hasNext = true
	}
	return p.nextBuf
}

// match consumes cur if its Type matches any of the given types.
// Returns true on match, false otherwise. Idiomatic for optional
// keywords like DISTINCT/ALL or NOT.
func (p *Parser) match(types ...int) bool {
	for _, t := range types {
		if p.cur.Type == t {
			p.advance()
			return true
		}
	}
	return false
}

// expect consumes cur if its Type matches tokenType, otherwise returns
// a *ParseError. Used for required tokens (PAREN_RIGHT, COMMA, etc.).
func (p *Parser) expect(tokenType int) (Token, error) {
	if p.cur.Type == tokenType {
		tok := p.cur
		p.advance()
		return tok, nil
	}
	return Token{}, &ParseError{
		Message: fmt.Sprintf("expected %s, got %q", tokenName(tokenType), p.cur.Str),
		Loc:     ast.Loc{Start: p.cur.Loc.Start, End: p.cur.Loc.Start},
	}
}

// checkLexerErr returns a *ParseError wrapping the lexer's error if
// any. Parser methods call this at strategic points (function entry,
// after expect) to propagate lexer errors into the parse result.
func (p *Parser) checkLexerErr() error {
	if p.lexer.Err != nil {
		return &ParseError{
			Message: p.lexer.Err.Error(),
			Loc:     ast.Loc{Start: p.cur.Loc.Start, End: p.cur.Loc.Start},
		}
	}
	return nil
}

// deferredFeature constructs a stub error pointing at the current
// token. Every non-foundation case in parsePrimaryBase (and a few
// other dispatch points) calls this helper so the error format is
// uniform across the whole parser.
//
// The error message format is the grep contract for future DAG node
// implementers: running
//
//	grep -rn "deferred to parser-builtins" partiql/parser/
//
// returns the full list of stub call sites node 15 needs to replace.
func (p *Parser) deferredFeature(feature, ownerNode string) error {
	return &ParseError{
		Message: fmt.Sprintf("%s is deferred to %s", feature, ownerNode),
		Loc:     ast.Loc{Start: p.cur.Loc.Start, End: p.cur.Loc.End},
	}
}

// parseSymbolPrimitive consumes an unquoted or double-quoted identifier
// and returns the name, case-sensitivity flag, and location. Rejects
// bare keywords. Matches grammar rule symbolPrimitive (line 742).
func (p *Parser) parseSymbolPrimitive() (name string, caseSensitive bool, loc ast.Loc, err error) {
	switch p.cur.Type {
	case tokIDENT:
		name = p.cur.Str
		caseSensitive = false
		loc = p.cur.Loc
		p.advance()
		return
	case tokIDENT_QUOTED:
		name = p.cur.Str
		caseSensitive = true
		loc = p.cur.Loc
		p.advance()
		return
	default:
		err = &ParseError{
			Message: fmt.Sprintf("expected identifier, got %q", p.cur.Str),
			Loc:     ast.Loc{Start: p.cur.Loc.Start, End: p.cur.Loc.Start},
		}
		return
	}
}

// parseVarRef handles optional @-prefix plus symbolPrimitive. Matches
// grammar rule varRefExpr (lines 635-636).
//
// NOTE: Task 10 upgrades this function to detect IDENT followed by
// PAREN_LEFT (function call form) and return a deferred-feature stub
// error. The Task 1 version handles only bare varref.
func (p *Parser) parseVarRef() (ast.ExprNode, error) {
	start := p.cur.Loc.Start
	atPrefixed := false
	if p.cur.Type == tokAT_SIGN {
		atPrefixed = true
		p.advance()
	}
	name, caseSensitive, nameLoc, err := p.parseSymbolPrimitive()
	if err != nil {
		return nil, err
	}
	return &ast.VarRef{
		Name:          name,
		AtPrefixed:    atPrefixed,
		CaseSensitive: caseSensitive,
		Loc:           ast.Loc{Start: start, End: nameLoc.End},
	}, nil
}
```

- [ ] **Step 2: Create `partiql/parser/parser_test.go` with TestParser_Machinery**

```go
package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/partiql/ast"
)

// TestParser_Machinery verifies the low-level token buffer helpers
// without any expression parsing. Each case constructs a Parser
// directly and drives the helpers through their state transitions.
func TestParser_Machinery(t *testing.T) {
	t.Run("new_parser_primes_cur", func(t *testing.T) {
		p := NewParser("foo")
		if p.cur.Type != tokIDENT {
			t.Errorf("cur.Type = %d, want tokIDENT", p.cur.Type)
		}
		if p.cur.Str != "foo" {
			t.Errorf("cur.Str = %q, want %q", p.cur.Str, "foo")
		}
	})

	t.Run("advance_moves_cur_forward", func(t *testing.T) {
		p := NewParser("foo bar")
		if p.peek().Str != "foo" {
			t.Errorf("first peek() = %q, want foo", p.peek().Str)
		}
		p.advance()
		if p.peek().Str != "bar" {
			t.Errorf("second peek() = %q, want bar", p.peek().Str)
		}
		if p.prev.Str != "foo" {
			t.Errorf("prev.Str = %q, want foo", p.prev.Str)
		}
	})

	t.Run("peek_next_lookahead", func(t *testing.T) {
		p := NewParser("foo bar baz")
		// cur = foo; peekNext = bar; cur still foo.
		if p.peek().Str != "foo" {
			t.Errorf("peek() = %q, want foo", p.peek().Str)
		}
		if p.peekNext().Str != "bar" {
			t.Errorf("peekNext() = %q, want bar", p.peekNext().Str)
		}
		if p.peek().Str != "foo" {
			t.Errorf("peek() after peekNext = %q, want foo", p.peek().Str)
		}
		// Second peekNext returns same token without re-reading.
		if p.peekNext().Str != "bar" {
			t.Errorf("second peekNext() = %q, want bar", p.peekNext().Str)
		}
		// advance consumes foo → cur=bar, nextBuf cleared.
		p.advance()
		if p.peek().Str != "bar" {
			t.Errorf("peek() after advance = %q, want bar", p.peek().Str)
		}
		if p.peekNext().Str != "baz" {
			t.Errorf("peekNext() after advance = %q, want baz", p.peekNext().Str)
		}
	})

	t.Run("match_consumes_on_hit", func(t *testing.T) {
		p := NewParser("AND OR")
		if !p.match(tokAND) {
			t.Errorf("match(tokAND) returned false")
		}
		if p.cur.Type != tokOR {
			t.Errorf("cur.Type after match = %d, want tokOR", p.cur.Type)
		}
		if p.match(tokAND) {
			t.Errorf("match(tokAND) returned true for non-matching token")
		}
		if p.match(tokOR, tokAND) {
			// match accepts any of the listed types.
			if p.cur.Type != tokEOF {
				t.Errorf("cur.Type after match = %d, want tokEOF", p.cur.Type)
			}
		} else {
			t.Errorf("match(tokOR, tokAND) returned false when tokOR was cur")
		}
	})

	t.Run("expect_returns_error_on_miss", func(t *testing.T) {
		p := NewParser("foo")
		_, err := p.expect(tokCOMMA)
		if err == nil {
			t.Fatal("expect(tokCOMMA) returned nil error on non-matching token")
		}
		perr, ok := err.(*ParseError)
		if !ok {
			t.Fatalf("error type = %T, want *ParseError", err)
		}
		if !strings.Contains(perr.Message, "expected") {
			t.Errorf("error message = %q, want to contain 'expected'", perr.Message)
		}
		if perr.Loc.Start != 0 {
			t.Errorf("error Loc.Start = %d, want 0", perr.Loc.Start)
		}
	})

	t.Run("expect_consumes_on_hit", func(t *testing.T) {
		p := NewParser("foo, bar")
		tok, err := p.expect(tokIDENT)
		if err != nil {
			t.Fatalf("expect(tokIDENT) error: %v", err)
		}
		if tok.Str != "foo" {
			t.Errorf("expect returned %q, want foo", tok.Str)
		}
		if p.cur.Type != tokCOMMA {
			t.Errorf("cur.Type after expect = %d, want tokCOMMA", p.cur.Type)
		}
	})

	t.Run("lexer_error_propagation", func(t *testing.T) {
		// Unterminated string: lexer sets Err, returns tokEOF.
		p := NewParser("'unclosed")
		// After advance(), cur is tokEOF and p.lexer.Err is set.
		if p.cur.Type != tokEOF {
			t.Errorf("cur.Type = %d, want tokEOF for unterminated string", p.cur.Type)
		}
		err := p.checkLexerErr()
		if err == nil {
			t.Fatal("checkLexerErr returned nil, want lexer error")
		}
		perr, ok := err.(*ParseError)
		if !ok {
			t.Fatalf("error type = %T, want *ParseError", err)
		}
		if !strings.Contains(perr.Message, "unterminated string literal") {
			t.Errorf("error message = %q, want to contain 'unterminated string literal'", perr.Message)
		}
	})

	t.Run("parse_symbol_primitive_bare", func(t *testing.T) {
		p := NewParser("foo")
		name, caseSensitive, loc, err := p.parseSymbolPrimitive()
		if err != nil {
			t.Fatalf("parseSymbolPrimitive error: %v", err)
		}
		if name != "foo" {
			t.Errorf("name = %q, want foo", name)
		}
		if caseSensitive {
			t.Error("caseSensitive = true, want false for bare ident")
		}
		if loc.Start != 0 || loc.End != 3 {
			t.Errorf("loc = %+v, want {0, 3}", loc)
		}
	})

	t.Run("parse_symbol_primitive_quoted", func(t *testing.T) {
		p := NewParser(`"Foo"`)
		name, caseSensitive, _, err := p.parseSymbolPrimitive()
		if err != nil {
			t.Fatalf("parseSymbolPrimitive error: %v", err)
		}
		if name != "Foo" {
			t.Errorf("name = %q, want Foo", name)
		}
		if !caseSensitive {
			t.Error("caseSensitive = false, want true for quoted ident")
		}
	})

	t.Run("parse_var_ref_bare", func(t *testing.T) {
		p := NewParser("foo")
		expr, err := p.parseVarRef()
		if err != nil {
			t.Fatalf("parseVarRef error: %v", err)
		}
		v, ok := expr.(*ast.VarRef)
		if !ok {
			t.Fatalf("parseVarRef returned %T, want *ast.VarRef", expr)
		}
		if v.Name != "foo" || v.AtPrefixed || v.CaseSensitive {
			t.Errorf("VarRef = %+v, want {Name:foo AtPrefixed:false CaseSensitive:false}", v)
		}
	})

	t.Run("parse_var_ref_at_prefixed", func(t *testing.T) {
		p := NewParser("@x")
		expr, err := p.parseVarRef()
		if err != nil {
			t.Fatalf("parseVarRef error: %v", err)
		}
		v := expr.(*ast.VarRef)
		if v.Name != "x" || !v.AtPrefixed {
			t.Errorf("VarRef = %+v, want {Name:x AtPrefixed:true}", v)
		}
	})

	t.Run("parse_var_ref_quoted", func(t *testing.T) {
		p := NewParser(`"Foo"`)
		expr, err := p.parseVarRef()
		if err != nil {
			t.Fatalf("parseVarRef error: %v", err)
		}
		v := expr.(*ast.VarRef)
		if v.Name != "Foo" || !v.CaseSensitive {
			t.Errorf("VarRef = %+v, want {Name:Foo CaseSensitive:true}", v)
		}
	})
}
```

- [ ] **Step 3: Build and run tests**

```bash
cd /Users/h3n4l/OpenSource/omni/.worktrees/feat-partiql-parser-foundation
go build ./partiql/parser/...
go test -v -run TestParser_Machinery ./partiql/parser/...
```

Expected: all 12 sub-tests pass. Existing lexer tests also still pass (no modification to lexer files).

- [ ] **Step 4: Vet, gofmt, corruption check**

```bash
go vet ./partiql/parser/...
gofmt -l partiql/parser/
grep -n $'\xe2\x80\x9d' partiql/parser/*.go
```

All clean. No U+201D corruption (the new file has no `''` sequences in `//` comments, so this is just a safety net).

- [ ] **Step 5: Commit**

```bash
git add partiql/parser/parser.go partiql/parser/parser_test.go
git commit -m "$(cat <<'EOF'
feat(partiql/parser): scaffold Parser struct + token buffer helpers

Adds partiql/parser/parser.go with the Parser struct, ParseError,
NewParser constructor, advance/peek/peekNext/match/expect helpers,
checkLexerErr for lexer-error propagation, deferredFeature helper
for non-foundation stubs, parseSymbolPrimitive (IDENT/IDENT_QUOTED),
and parseVarRef (bare form without function call detection; Task
10 upgrades it).

Adds partiql/parser/parser_test.go with TestParser_Machinery — 12
sub-tests covering each helper directly without any expression
parsing. Lexer-error propagation test exercises the
unterminated-string path from lexer.go.

First task of the parser-foundation (DAG node 4) implementation
per docs/superpowers/plans/2026-04-09-partiql-parser-foundation.md.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 2: literals.go + TestParser_Goldens runner + 10 literal goldens

**Files:**
- Create: `partiql/parser/literals.go`
- Modify: `partiql/parser/parser.go` (add `ParseExpr` public entry point wired to `parseLiteral` for now)
- Modify: `partiql/parser/parser_test.go` (add `TestParser_Goldens` runner + `update` flag)
- Create: `partiql/parser/testdata/parser-foundation/literal_null.partiql` (and 9 more literal inputs)
- Create: `partiql/parser/testdata/parser-foundation/literal_null.golden` (and 9 more golden outputs)

This task sets up the golden-file test harness with the `-update` flag and lands the first real AST-producing code (literal parsing). It also introduces `Parser.ParseExpr()` as the public entry point. Initially `ParseExpr` just dispatches to `parseLiteral`; subsequent tasks extend the dispatch as the precedence ladder grows.

- [ ] **Step 1: Create `partiql/parser/literals.go`**

```go
package parser

import (
	"github.com/bytebase/omni/partiql/ast"
)

// parseLiteral dispatches on the current token to produce one of the
// literal AST nodes. Handles the 6 real forms (NULL/MISSING/TRUE/
// FALSE/string/number/Ion) directly. DATE and TIME literal forms
// are stubbed with a "deferred to parser-datetime-literals (DAG
// node 18)" error.
//
// Grammar: literal (PartiQLParser.g4 lines 661-672).
func (p *Parser) parseLiteral() (ast.ExprNode, error) {
	switch p.cur.Type {
	case tokNULL:
		loc := p.cur.Loc
		p.advance()
		return &ast.NullLit{Loc: loc}, nil

	case tokMISSING:
		loc := p.cur.Loc
		p.advance()
		return &ast.MissingLit{Loc: loc}, nil

	case tokTRUE:
		loc := p.cur.Loc
		p.advance()
		return &ast.BoolLit{Val: true, Loc: loc}, nil

	case tokFALSE:
		loc := p.cur.Loc
		p.advance()
		return &ast.BoolLit{Val: false, Loc: loc}, nil

	case tokSCONST:
		val := p.cur.Str
		loc := p.cur.Loc
		p.advance()
		return &ast.StringLit{Val: val, Loc: loc}, nil

	case tokICONST, tokFCONST:
		// NumberLit.Val preserves the raw source text so callers can
		// distinguish integer/decimal/scientific forms. Token.Str is
		// already the raw slice (scanNumber in the lexer preserves it).
		val := p.cur.Str
		loc := p.cur.Loc
		p.advance()
		return &ast.NumberLit{Val: val, Loc: loc}, nil

	case tokION_LITERAL:
		// Lexer's scanIonLiteral delivers the verbatim inner content
		// (between backticks) in Token.Str. No further decoding at
		// this layer — that's DAG node 17's job.
		text := p.cur.Str
		loc := p.cur.Loc
		p.advance()
		return &ast.IonLit{Text: text, Loc: loc}, nil

	case tokDATE:
		return nil, p.deferredFeature("DATE literal", "parser-datetime-literals (DAG node 18)")

	case tokTIME:
		return nil, p.deferredFeature("TIME literal", "parser-datetime-literals (DAG node 18)")
	}

	return nil, &ParseError{
		Message: "expected literal",
		Loc:     ast.Loc{Start: p.cur.Loc.Start, End: p.cur.Loc.Start},
	}
}
```

- [ ] **Step 2: Add `ParseExpr` to `partiql/parser/parser.go`**

Append to the end of `parser.go`:

```go
// ParseExpr parses a single expression from the parser's input.
//
// This is the foundation-level public entry point. Nodes 5-8 will
// add ParseStatement and ParseScript (SelectStmt-producing forms).
// Node 8 will add the top-level Parse(sql) function.
//
// At this milestone (Task 2), ParseExpr only handles literal tokens.
// Each subsequent task extends ParseExpr's dispatch target as the
// precedence ladder grows:
//
//	Task 2 (this): ParseExpr → parseLiteral
//	Task 5:        ParseExpr → parsePrimary (literals + varref + collections + paths)
//	Task 6:        ParseExpr → parseMathOp00
//	Task 7:        ParseExpr → parsePredicate
//	Task 8:        ParseExpr → parseOr
//	Task 9:        ParseExpr → parseBagOp
func (p *Parser) ParseExpr() (ast.ExprNode, error) {
	if err := p.checkLexerErr(); err != nil {
		return nil, err
	}
	return p.parseLiteral()
}
```

- [ ] **Step 3: Add TestParser_Goldens runner to `partiql/parser/parser_test.go`**

At the top of `parser_test.go`, replace the existing import block with:

```go
package parser

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/partiql/ast"
)

// update regenerates golden files. Run with `-update` to refresh.
var update = flag.Bool("update", false, "update golden files")
```

Then append the following test function (after `TestParser_Machinery`):

```go
// TestParser_Goldens iterates every .partiql file under
// testdata/parser-foundation/ and compares the parser's pretty-printed
// output (via ast.NodeToString) against the matching .golden file.
//
// Run with `go test -update -run TestParser_Goldens ./partiql/parser/...`
// to regenerate goldens after intentional AST shape changes.
func TestParser_Goldens(t *testing.T) {
	files, err := filepath.Glob("testdata/parser-foundation/*.partiql")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no golden inputs found under testdata/parser-foundation/")
	}
	for _, inPath := range files {
		name := strings.TrimSuffix(filepath.Base(inPath), ".partiql")
		t.Run(name, func(t *testing.T) {
			input, err := os.ReadFile(inPath)
			if err != nil {
				t.Fatal(err)
			}
			p := NewParser(string(input))
			expr, err := p.ParseExpr()
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			got := ast.NodeToString(expr)
			goldenPath := strings.TrimSuffix(inPath, ".partiql") + ".golden"
			if *update {
				if err := os.WriteFile(goldenPath, []byte(got+"\n"), 0o644); err != nil {
					t.Fatal(err)
				}
				return
			}
			want, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf("golden file missing: %s (run with -update to create)", goldenPath)
			}
			if got+"\n" != string(want) {
				t.Errorf("AST mismatch\ngot:\n%s\nwant:\n%s", got, string(want))
			}
		})
	}
}
```

Note the `ast.NodeToString` import is now load-bearing in two places (the machinery test already imports ast for `*ast.VarRef`, and the golden runner uses `ast.NodeToString`).

- [ ] **Step 4: Create the 10 literal input files**

Create `partiql/parser/testdata/parser-foundation/` and populate:

```bash
mkdir -p partiql/parser/testdata/parser-foundation
```

Create each file below using the Write tool (each file's contents are a single line of PartiQL):

- `literal_null.partiql`: `NULL`
- `literal_missing.partiql`: `MISSING`
- `literal_true.partiql`: `TRUE`
- `literal_false.partiql`: `FALSE`
- `literal_int.partiql`: `42`
- `literal_decimal.partiql`: `3.14`
- `literal_scientific.partiql`: `1.5e-3`
- `literal_string.partiql`: `'hello'`
- `literal_string_doubled.partiql`: `'it''s'`
- `literal_ion.partiql`: `` `{a: 1}` ``

- [ ] **Step 5: Create the 10 golden files**

Each `.golden` file contains the exact `NodeToString` output followed by a trailing newline (the runner appends `\n` before comparison). Create these files manually to lock in the expected AST shape:

- `literal_null.golden`:
  ```
  NullLit{}
  ```
- `literal_missing.golden`:
  ```
  MissingLit{}
  ```
- `literal_true.golden`:
  ```
  BoolLit{Val:true}
  ```
- `literal_false.golden`:
  ```
  BoolLit{Val:false}
  ```
- `literal_int.golden`:
  ```
  NumberLit{Val:42}
  ```
- `literal_decimal.golden`:
  ```
  NumberLit{Val:3.14}
  ```
- `literal_scientific.golden`:
  ```
  NumberLit{Val:1.5e-3}
  ```
- `literal_string.golden`:
  ```
  StringLit{Val:"hello"}
  ```
- `literal_string_doubled.golden`:
  ```
  StringLit{Val:"it's"}
  ```
- `literal_ion.golden`:
  ```
  IonLit{Text:"{a: 1}"}
  ```

Each file ends with a single `\n` newline. Verify with `od -c literal_null.golden` — last byte should be `\n`.

- [ ] **Step 6: Build and run tests**

```bash
go build ./partiql/parser/...
go test -v -run "TestParser_Machinery|TestParser_Goldens" ./partiql/parser/...
```

Expected: 12 TestParser_Machinery sub-tests + 10 TestParser_Goldens sub-tests = 22 new sub-tests passing.

If any golden mismatches: read the reported diff carefully. If the mismatch is because the hand-written golden is wrong (not because the parser is wrong), fix the golden. If the parser is wrong, fix the parser. Do NOT blindly regenerate with `-update` without verifying the output matches expectations.

- [ ] **Step 7: Vet, gofmt, corruption check**

```bash
go vet ./partiql/parser/...
gofmt -l partiql/parser/
grep -n $'\xe2\x80\x9d' partiql/parser/*.go
```

All clean.

- [ ] **Step 8: Commit**

```bash
git add partiql/parser/literals.go partiql/parser/parser.go partiql/parser/parser_test.go partiql/parser/testdata/parser-foundation/
git commit -m "$(cat <<'EOF'
feat(partiql/parser): parseLiteral + TestParser_Goldens harness

Adds partiql/parser/literals.go with parseLiteral dispatching on the
6 real literal forms (NULL/MISSING/TRUE/FALSE/string/number/Ion).
DATE and TIME literal bodies return a deferred-feature stub error
pointing at parser-datetime-literals (DAG node 18).

Adds Parser.ParseExpr as the foundation-level public entry point.
Initial implementation dispatches directly to parseLiteral; later
tasks extend the dispatch as the precedence ladder grows.

Adds TestParser_Goldens runner to parser_test.go with a -update
flag that regenerates .golden files. Committed as 10 literal test
cases under testdata/parser-foundation/ — hand-written goldens
lock in the exact AST shape for each literal form.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 3: parser.go parseType + TestParseType

**Files:**
- Modify: `partiql/parser/parser.go` (add `parseType` function)
- Modify: `partiql/parser/parser_test.go` (add `TestParseType` table-driven test)

`parseType` is the shared utility that handles all 30+ type forms in the PartiQL grammar. It's used by CAST (deferred to parser-builtins) and DDL column types (deferred to parser-ddl), but shipping it in foundation decouples those nodes. The tests use direct method invocation (not goldens) because TypeRef is a leaf node and inline assertions are clearer.

- [ ] **Step 1: Add parseType to `partiql/parser/parser.go`**

Append to the end of `parser.go` (after `parseVarRef`):

```go
// parseType consumes one of the PartiQL type forms and returns a
// *ast.TypeRef. Handles:
//
//   - Atomic types: NULL, BOOL, BOOLEAN, SMALLINT, INT/INT2/INT4/INT8,
//     INTEGER/INTEGER2/INTEGER4/INTEGER8, BIGINT, REAL, TIMESTAMP,
//     CHAR, CHARACTER, MISSING, STRING, SYMBOL, BLOB, CLOB, DATE,
//     STRUCT, TUPLE, LIST, SEXP, BAG, ANY
//   - DOUBLE PRECISION (two-token form)
//   - Parameterized single-arg: CHAR(n), CHARACTER(n), FLOAT(p), VARCHAR(n)
//   - CHARACTER VARYING [(n)]
//   - Parameterized two-arg: DECIMAL(p,s), DEC(p,s), NUMERIC(p,s)
//   - TIME [(p)] [WITH TIME ZONE]
//   - Custom: any symbolPrimitive identifier (fallback)
//
// Grammar: type (PartiQLParser.g4 lines 674-686).
//
// Foundation ships parseType even though CAST is stubbed so that
// parser-ddl (DAG node 7) and parser-builtins (DAG node 15) can each
// consume it without coupling.
func (p *Parser) parseType() (*ast.TypeRef, error) {
	start := p.cur.Loc.Start

	// DOUBLE PRECISION is a two-token form.
	if p.cur.Type == tokDOUBLE {
		p.advance()
		if p.cur.Type != tokPRECISION {
			return nil, &ParseError{
				Message: fmt.Sprintf("expected PRECISION after DOUBLE, got %q", p.cur.Str),
				Loc:     ast.Loc{Start: p.cur.Loc.Start, End: p.cur.Loc.Start},
			}
		}
		end := p.cur.Loc.End
		p.advance()
		return &ast.TypeRef{
			Name: "DOUBLE PRECISION",
			Loc:  ast.Loc{Start: start, End: end},
		}, nil
	}

	// CHARACTER VARYING is another two-token form that also takes an
	// optional (n) argument.
	if p.cur.Type == tokCHARACTER && p.peekNext().Type == tokVARYING {
		p.advance() // CHARACTER
		p.advance() // VARYING
		end := p.prev.Loc.End
		args, argsEnd, err := p.parseOptionalTypeArgs(1)
		if err != nil {
			return nil, err
		}
		if argsEnd > 0 {
			end = argsEnd
		}
		return &ast.TypeRef{
			Name: "CHARACTER VARYING",
			Args: args,
			Loc:  ast.Loc{Start: start, End: end},
		}, nil
	}

	// TIME [(p)] [WITH TIME ZONE] — TIME requires a tailored parse path
	// because it supports both a precision arg and the WITH TIME ZONE
	// trailing keywords.
	if p.cur.Type == tokTIME {
		p.advance()
		end := p.prev.Loc.End
		args, argsEnd, err := p.parseOptionalTypeArgs(1)
		if err != nil {
			return nil, err
		}
		if argsEnd > 0 {
			end = argsEnd
		}
		withTZ := false
		if p.cur.Type == tokWITH && p.peekNext().Type == tokTIME {
			// WITH TIME ZONE
			p.advance() // WITH
			p.advance() // TIME
			if p.cur.Type != tokZONE {
				return nil, &ParseError{
					Message: fmt.Sprintf("expected ZONE after WITH TIME, got %q", p.cur.Str),
					Loc:     ast.Loc{Start: p.cur.Loc.Start, End: p.cur.Loc.Start},
				}
			}
			withTZ = true
			end = p.cur.Loc.End
			p.advance() // ZONE
		}
		return &ast.TypeRef{
			Name:         "TIME",
			Args:         args,
			WithTimeZone: withTZ,
			Loc:          ast.Loc{Start: start, End: end},
		}, nil
	}

	// Atomic types (no args). Map the keyword token to its canonical
	// uppercase name. This switch handles the "bare keyword" cases from
	// the grammar's TypeAtomic alternative (lines 675-680).
	atomicName := ""
	switch p.cur.Type {
	case tokNULL:
		atomicName = "NULL"
	case tokBOOL:
		atomicName = "BOOL"
	case tokBOOLEAN:
		atomicName = "BOOLEAN"
	case tokSMALLINT:
		atomicName = "SMALLINT"
	case tokINT2:
		atomicName = "INT2"
	case tokINTEGER2:
		atomicName = "INTEGER2"
	case tokINT4:
		atomicName = "INT4"
	case tokINTEGER4:
		atomicName = "INTEGER4"
	case tokINT8:
		atomicName = "INT8"
	case tokINTEGER8:
		atomicName = "INTEGER8"
	case tokINT:
		atomicName = "INT"
	case tokINTEGER:
		atomicName = "INTEGER"
	case tokBIGINT:
		atomicName = "BIGINT"
	case tokREAL:
		atomicName = "REAL"
	case tokTIMESTAMP:
		atomicName = "TIMESTAMP"
	case tokMISSING:
		atomicName = "MISSING"
	case tokSTRING:
		atomicName = "STRING"
	case tokSYMBOL:
		atomicName = "SYMBOL"
	case tokBLOB:
		atomicName = "BLOB"
	case tokCLOB:
		atomicName = "CLOB"
	case tokDATE:
		atomicName = "DATE"
	case tokSTRUCT:
		atomicName = "STRUCT"
	case tokTUPLE:
		atomicName = "TUPLE"
	case tokLIST:
		atomicName = "LIST"
	case tokSEXP:
		atomicName = "SEXP"
	case tokBAG:
		atomicName = "BAG"
	case tokANY:
		atomicName = "ANY"
	}
	if atomicName != "" {
		end := p.cur.Loc.End
		p.advance()
		return &ast.TypeRef{
			Name: atomicName,
			Loc:  ast.Loc{Start: start, End: end},
		}, nil
	}

	// Parameterized single-arg: CHAR(n), CHARACTER(n), FLOAT(p), VARCHAR(n).
	switch p.cur.Type {
	case tokCHAR:
		return p.parseTypeWithArgs(start, "CHAR", 1, 1)
	case tokCHARACTER:
		// Note: CHARACTER VARYING already handled above via peekNext.
		return p.parseTypeWithArgs(start, "CHARACTER", 1, 1)
	case tokFLOAT:
		return p.parseTypeWithArgs(start, "FLOAT", 1, 1)
	case tokVARCHAR:
		return p.parseTypeWithArgs(start, "VARCHAR", 1, 1)
	}

	// Parameterized two-arg: DECIMAL(p,s), DEC(p,s), NUMERIC(p,s).
	switch p.cur.Type {
	case tokDECIMAL:
		return p.parseTypeWithArgs(start, "DECIMAL", 1, 2)
	case tokDEC:
		return p.parseTypeWithArgs(start, "DEC", 1, 2)
	case tokNUMERIC:
		return p.parseTypeWithArgs(start, "NUMERIC", 1, 2)
	}

	// Custom type: any symbolPrimitive. The grammar calls this
	// TypeCustom (line 685).
	if p.cur.Type == tokIDENT || p.cur.Type == tokIDENT_QUOTED {
		name, _, nameLoc, err := p.parseSymbolPrimitive()
		if err != nil {
			return nil, err
		}
		return &ast.TypeRef{
			Name: name,
			Loc:  ast.Loc{Start: start, End: nameLoc.End},
		}, nil
	}

	return nil, &ParseError{
		Message: fmt.Sprintf("expected type, got %q", p.cur.Str),
		Loc:     ast.Loc{Start: p.cur.Loc.Start, End: p.cur.Loc.Start},
	}
}

// parseTypeWithArgs consumes a type keyword followed by an optional
// parenthesized argument list of integers. Used by CHAR(n), FLOAT(p),
// VARCHAR(n), DECIMAL(p,s), and related forms.
//
//   - name:    canonical type name (e.g. "DECIMAL")
//   - minArgs: minimum number of integer args if the parenthesized
//              list is present (0 for optional, 1 if required when
//              parens are present)
//   - maxArgs: maximum number of integer args
func (p *Parser) parseTypeWithArgs(start int, name string, minArgs, maxArgs int) (*ast.TypeRef, error) {
	// Consume the type keyword.
	end := p.cur.Loc.End
	p.advance()
	args, argsEnd, err := p.parseOptionalTypeArgs(maxArgs)
	if err != nil {
		return nil, err
	}
	if argsEnd > 0 {
		end = argsEnd
		if len(args) < minArgs {
			return nil, &ParseError{
				Message: fmt.Sprintf("%s: expected at least %d argument(s), got %d", name, minArgs, len(args)),
				Loc:     ast.Loc{Start: start, End: end},
			}
		}
	}
	return &ast.TypeRef{
		Name: name,
		Args: args,
		Loc:  ast.Loc{Start: start, End: end},
	}, nil
}

// parseOptionalTypeArgs consumes an optional (n[, n]*) argument list
// bounded by parentheses. If cur is not PAREN_LEFT, returns
// (nil, 0, nil) without consuming anything. Otherwise consumes the
// parenthesized comma-separated integer list and returns the parsed
// args plus the End position of the closing paren.
//
// maxArgs is the hard limit on the number of integers; exceeding it
// yields a ParseError.
func (p *Parser) parseOptionalTypeArgs(maxArgs int) (args []int, end int, err error) {
	if p.cur.Type != tokPAREN_LEFT {
		return nil, 0, nil
	}
	p.advance() // consume (
	for {
		if p.cur.Type != tokICONST {
			return nil, 0, &ParseError{
				Message: fmt.Sprintf("expected integer argument, got %q", p.cur.Str),
				Loc:     ast.Loc{Start: p.cur.Loc.Start, End: p.cur.Loc.Start},
			}
		}
		n, perr := parseIntLiteral(p.cur.Str)
		if perr != nil {
			return nil, 0, &ParseError{
				Message: fmt.Sprintf("invalid integer argument %q: %v", p.cur.Str, perr),
				Loc:     ast.Loc{Start: p.cur.Loc.Start, End: p.cur.Loc.End},
			}
		}
		args = append(args, n)
		p.advance()
		if len(args) > maxArgs {
			return nil, 0, &ParseError{
				Message: fmt.Sprintf("too many type arguments (max %d)", maxArgs),
				Loc:     ast.Loc{Start: p.cur.Loc.Start, End: p.cur.Loc.Start},
			}
		}
		if p.cur.Type == tokCOMMA {
			p.advance()
			continue
		}
		break
	}
	rp, perr := p.expect(tokPAREN_RIGHT)
	if perr != nil {
		return nil, 0, perr
	}
	return args, rp.Loc.End, nil
}

// parseIntLiteral converts a token's Str (raw source text for an
// integer literal) into an int. We use strconv.Atoi directly because
// PartiQL integer literals are plain decimal digits per the grammar
// LITERAL_INTEGER rule (no hex, no underscores, no sign — signs come
// from the unary-minus operator).
func parseIntLiteral(s string) (int, error) {
	// Imports: strconv (add to parser.go's import block if not already there)
	return strconv.Atoi(s)
}
```

And update the import block at the top of `parser.go`:

```go
import (
	"fmt"
	"strconv"

	"github.com/bytebase/omni/partiql/ast"
)
```

- [ ] **Step 2: Add TestParseType to `partiql/parser/parser_test.go`**

Append after `TestParser_Goldens`:

```go
// TestParseType exhaustively tests parseType across the 30+ type
// forms in PartiQLParser.g4's `type` rule. Uses direct table-driven
// assertions because TypeRef is a leaf node (no child recursion) and
// inline expected values are clearer than filesystem goldens for
// exhaustive enumeration.
func TestParseType(t *testing.T) {
	cases := []struct {
		name         string
		input        string
		wantName     string
		wantArgs     []int
		wantWithTZ   bool
	}{
		// Atomic types.
		{"null", "NULL", "NULL", nil, false},
		{"bool", "BOOL", "BOOL", nil, false},
		{"boolean", "BOOLEAN", "BOOLEAN", nil, false},
		{"smallint", "SMALLINT", "SMALLINT", nil, false},
		{"int2", "INT2", "INT2", nil, false},
		{"integer2", "INTEGER2", "INTEGER2", nil, false},
		{"int4", "INT4", "INT4", nil, false},
		{"integer4", "INTEGER4", "INTEGER4", nil, false},
		{"int8", "INT8", "INT8", nil, false},
		{"integer8", "INTEGER8", "INTEGER8", nil, false},
		{"int", "INT", "INT", nil, false},
		{"integer", "INTEGER", "INTEGER", nil, false},
		{"bigint", "BIGINT", "BIGINT", nil, false},
		{"real", "REAL", "REAL", nil, false},
		{"timestamp", "TIMESTAMP", "TIMESTAMP", nil, false},
		{"missing", "MISSING", "MISSING", nil, false},
		{"string", "STRING", "STRING", nil, false},
		{"symbol", "SYMBOL", "SYMBOL", nil, false},
		{"blob", "BLOB", "BLOB", nil, false},
		{"clob", "CLOB", "CLOB", nil, false},
		{"date", "DATE", "DATE", nil, false},
		{"struct", "STRUCT", "STRUCT", nil, false},
		{"tuple", "TUPLE", "TUPLE", nil, false},
		{"list", "LIST", "LIST", nil, false},
		{"sexp", "SEXP", "SEXP", nil, false},
		{"bag", "BAG", "BAG", nil, false},
		{"any", "ANY", "ANY", nil, false},

		// Two-token DOUBLE PRECISION.
		{"double_precision", "DOUBLE PRECISION", "DOUBLE PRECISION", nil, false},

		// Parameterized single-arg types.
		{"char", "CHAR", "CHAR", nil, false},
		{"char_n", "CHAR(10)", "CHAR", []int{10}, false},
		{"character_n", "CHARACTER(20)", "CHARACTER", []int{20}, false},
		{"varchar", "VARCHAR", "VARCHAR", nil, false},
		{"varchar_n", "VARCHAR(255)", "VARCHAR", []int{255}, false},
		{"float", "FLOAT", "FLOAT", nil, false},
		{"float_p", "FLOAT(53)", "FLOAT", []int{53}, false},

		// CHARACTER VARYING two-token form.
		{"character_varying", "CHARACTER VARYING", "CHARACTER VARYING", nil, false},
		{"character_varying_n", "CHARACTER VARYING(80)", "CHARACTER VARYING", []int{80}, false},

		// Parameterized two-arg types.
		{"decimal", "DECIMAL", "DECIMAL", nil, false},
		{"decimal_p", "DECIMAL(10)", "DECIMAL", []int{10}, false},
		{"decimal_p_s", "DECIMAL(10,2)", "DECIMAL", []int{10, 2}, false},
		{"dec_p_s", "DEC(5,0)", "DEC", []int{5, 0}, false},
		{"numeric_p_s", "NUMERIC(18,4)", "NUMERIC", []int{18, 4}, false},

		// TIME with precision and WITH TIME ZONE.
		{"time", "TIME", "TIME", nil, false},
		{"time_p", "TIME(6)", "TIME", []int{6}, false},
		{"time_wtz", "TIME WITH TIME ZONE", "TIME", nil, true},
		{"time_p_wtz", "TIME(3) WITH TIME ZONE", "TIME", []int{3}, true},

		// Custom types (symbolPrimitive fallback).
		{"custom_ident", "MyType", "MyType", nil, false},
		{"custom_quoted", `"MyType"`, "MyType", nil, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			typeRef, err := p.parseType()
			if err != nil {
				t.Fatalf("parseType error: %v", err)
			}
			if typeRef.Name != tc.wantName {
				t.Errorf("Name = %q, want %q", typeRef.Name, tc.wantName)
			}
			if !intSliceEq(typeRef.Args, tc.wantArgs) {
				t.Errorf("Args = %v, want %v", typeRef.Args, tc.wantArgs)
			}
			if typeRef.WithTimeZone != tc.wantWithTZ {
				t.Errorf("WithTimeZone = %v, want %v", typeRef.WithTimeZone, tc.wantWithTZ)
			}
		})
	}
}

func intSliceEq(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
```

- [ ] **Step 3: Build and run**

```bash
go build ./partiql/parser/...
go test -v -run TestParseType ./partiql/parser/...
```

Expected: 46 TestParseType sub-tests pass.

- [ ] **Step 4: Full test suite — everything still passes**

```bash
go test ./partiql/parser/...
```

Expected: all lexer tests (152 sub-tests) + TestParser_Machinery (12) + TestParser_Goldens (10) + TestParseType (46) = 220 sub-tests total.

- [ ] **Step 5: Vet, gofmt, corruption check**

```bash
go vet ./partiql/parser/...
gofmt -l partiql/parser/
grep -n $'\xe2\x80\x9d' partiql/parser/*.go
```

All clean.

- [ ] **Step 6: Commit**

```bash
git add partiql/parser/parser.go partiql/parser/parser_test.go
git commit -m "$(cat <<'EOF'
feat(partiql/parser): parseType for 30+ PartiQL type forms

Adds parseType to parser.go handling every alternative of the
PartiQLParser.g4 type rule (lines 674-686):
- Atomic: NULL, BOOL, INT*, BIGINT, REAL, TIMESTAMP, MISSING,
  STRING, SYMBOL, BLOB, CLOB, DATE, STRUCT, TUPLE, LIST, SEXP,
  BAG, ANY, and related aliases
- DOUBLE PRECISION (two-token)
- CHAR(n) / CHARACTER(n) / VARCHAR(n) / FLOAT(p) single-arg
- CHARACTER VARYING [(n)] two-token + optional arg
- DECIMAL(p,s) / DEC(p,s) / NUMERIC(p,s) two-arg
- TIME [(p)] [WITH TIME ZONE]
- Custom types via symbolPrimitive fallback

Foundation ships parseType even though CAST is deferred so
parser-ddl (DAG node 7) and parser-builtins (DAG node 15) can
each consume it without coupling.

Adds TestParseType with 46 table-driven cases covering every
form. Uses direct assertions (not goldens) because TypeRef is a
leaf node and inline expected values are clearer for exhaustive
enumeration.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 4: path.go + 7 path goldens

**Files:**
- Create: `partiql/parser/path.go`
- Modify: `partiql/parser/parser_test.go` (add a temporary unit test that bypasses `parsePrimary`)
- Create: `partiql/parser/testdata/parser-foundation/path_*.partiql` (7 files)
- Create: `partiql/parser/testdata/parser-foundation/path_*.golden` (7 files)

Task 4 adds `parsePathSteps`/`parsePathStep`/`isPathStepStart` but does NOT yet wire path attachment into `parsePrimary` (which doesn't exist yet — that's Task 5). Unit tests exercise `parsePathSteps` directly with a hand-constructed base expression. Golden files are created but won't be consumable by the TestParser_Goldens runner until Task 5, when `ParseExpr` starts routing through `parsePrimary`.

**IMPORTANT:** For this task, the 7 path golden inputs are committed but the goldens will FAIL the TestParser_Goldens runner because `ParseExpr` still routes to `parseLiteral` only. Task 5 extends ParseExpr through parsePrimary, at which point the path goldens start passing. To avoid a broken state, this task skips the goldens in TestParser_Goldens and uses a direct unit test instead.

- [ ] **Step 1: Create `partiql/parser/path.go`**

```go
package parser

import (
	"github.com/bytebase/omni/partiql/ast"
)

// isPathStepStart reports whether the current token can begin a path
// step. Path steps appear after a primary expression in the form:
//
//	.field  ."Field"  .*  [expr]  [*]
//
// Grammar: pathStep (PartiQLParser.g4 lines 618-623).
func isPathStepStart(t int) bool {
	return t == tokPERIOD || t == tokBRACKET_LEFT
}

// parsePathSteps consumes a sequence of path steps and wraps the base
// expression in an ast.PathExpr. Called by parsePrimary when the
// current token after the base is a path-step start.
//
// This function assumes at least one path step is present (the caller
// checks isPathStepStart before invoking). If called with no path
// step to consume, returns the base unchanged wrapped in a PathExpr
// with an empty Steps slice — callers should avoid that degenerate
// case.
//
// Grammar: exprPrimary#ExprPrimaryPath (line 528) + pathStep (lines
// 618-623).
func (p *Parser) parsePathSteps(base ast.ExprNode) (*ast.PathExpr, error) {
	var steps []ast.PathStep
	for isPathStepStart(p.cur.Type) {
		step, err := p.parsePathStep()
		if err != nil {
			return nil, err
		}
		steps = append(steps, step)
	}
	endLoc := base.GetLoc().End
	if len(steps) > 0 {
		endLoc = steps[len(steps)-1].GetLoc().End
	}
	return &ast.PathExpr{
		Root:  base,
		Steps: steps,
		Loc:   ast.Loc{Start: base.GetLoc().Start, End: endLoc},
	}, nil
}

// parsePathStep consumes exactly one path step.
//
// Grammar: pathStep (lines 618-623):
//
//	BRACKET_LEFT key=expr BRACKET_RIGHT        # PathStepIndexExpr
//	BRACKET_LEFT ASTERISK BRACKET_RIGHT        # PathStepIndexAll
//	PERIOD key=symbolPrimitive                 # PathStepDotExpr
//	PERIOD ASTERISK                            # PathStepDotAll
func (p *Parser) parsePathStep() (ast.PathStep, error) {
	switch p.cur.Type {
	case tokPERIOD:
		start := p.cur.Loc.Start
		p.advance() // consume .
		// Dot-star: .*
		if p.cur.Type == tokASTERISK {
			end := p.cur.Loc.End
			p.advance()
			return &ast.AllFieldsStep{
				Loc: ast.Loc{Start: start, End: end},
			}, nil
		}
		// Dot-field: .foo or ."Foo"
		name, caseSensitive, nameLoc, err := p.parseSymbolPrimitive()
		if err != nil {
			return nil, err
		}
		return &ast.DotStep{
			Field:         name,
			CaseSensitive: caseSensitive,
			Loc:           ast.Loc{Start: start, End: nameLoc.End},
		}, nil

	case tokBRACKET_LEFT:
		start := p.cur.Loc.Start
		p.advance() // consume [
		// Bracket-star: [*]
		if p.cur.Type == tokASTERISK && p.peekNext().Type == tokBRACKET_RIGHT {
			p.advance() // consume *
			end := p.cur.Loc.End
			p.advance() // consume ]
			return &ast.WildcardStep{
				Loc: ast.Loc{Start: start, End: end},
			}, nil
		}
		// Bracket-expr: [expr]
		idx, err := p.ParseExpr()
		if err != nil {
			return nil, err
		}
		rp, err := p.expect(tokBRACKET_RIGHT)
		if err != nil {
			return nil, err
		}
		return &ast.IndexStep{
			Index: idx,
			Loc:   ast.Loc{Start: start, End: rp.Loc.End},
		}, nil
	}

	return nil, &ParseError{
		Message: "expected path step (. or [)",
		Loc:     ast.Loc{Start: p.cur.Loc.Start, End: p.cur.Loc.Start},
	}
}
```

**Note on the `[expr]` case calling `ParseExpr`:** At this task, `ParseExpr` dispatches to `parseLiteral` only, so `[42]` works but `[x]` (an identifier index) fails because `parseLiteral` can't handle identifiers. Task 5 extends `ParseExpr` through `parsePrimary`, at which point `[x]` and more complex index expressions start working. The 7 path golden inputs in this task use only forms that work with the current `ParseExpr` dispatch: paths are constructed via unit tests with a hand-built base expression and literal index values.

- [ ] **Step 2: Add TestParser_PathUnit direct unit test**

Append to `partiql/parser/parser_test.go`:

```go
// TestParser_PathUnit exercises parsePathSteps directly by constructing
// a base VarRef and then calling parsePathSteps. This bypasses the
// ParseExpr dispatch (which doesn't yet route through parsePrimary at
// this task). Task 5 removes this test and replaces its coverage with
// file-based path goldens consumed by TestParser_Goldens.
func TestParser_PathUnit(t *testing.T) {
	cases := []struct {
		name  string
		input string // path suffix only — the base is a fixed VarRef
		want  string // expected NodeToString output
	}{
		{
			name:  "dot",
			input: ".foo",
			want:  `PathExpr{Root:VarRef{Name:t} Steps:[DotStep{Field:foo}]}`,
		},
		{
			name:  "dot_quoted",
			input: `."Foo"`,
			want:  `PathExpr{Root:VarRef{Name:t} Steps:[DotStep{Field:Foo CaseSensitive:true}]}`,
		},
		{
			name:  "dot_star",
			input: ".*",
			want:  `PathExpr{Root:VarRef{Name:t} Steps:[AllFieldsStep{}]}`,
		},
		{
			name:  "index_int",
			input: "[0]",
			want:  `PathExpr{Root:VarRef{Name:t} Steps:[IndexStep{Index:NumberLit{Val:0}}]}`,
		},
		{
			name:  "index_wildcard",
			input: "[*]",
			want:  `PathExpr{Root:VarRef{Name:t} Steps:[WildcardStep{}]}`,
		},
		{
			name:  "chain_dot_dot",
			input: ".foo.bar",
			want:  `PathExpr{Root:VarRef{Name:t} Steps:[DotStep{Field:foo} DotStep{Field:bar}]}`,
		},
		{
			name:  "chain_mixed",
			input: ".foo[0].*[*]",
			want:  `PathExpr{Root:VarRef{Name:t} Steps:[DotStep{Field:foo} IndexStep{Index:NumberLit{Val:0}} AllFieldsStep{} WildcardStep{}]}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			base := &ast.VarRef{Name: "t", Loc: ast.Loc{Start: -1, End: -1}}
			if !isPathStepStart(p.cur.Type) {
				t.Fatalf("expected path step start, got %d", p.cur.Type)
			}
			got, err := p.parsePathSteps(base)
			if err != nil {
				t.Fatalf("parsePathSteps error: %v", err)
			}
			gotStr := ast.NodeToString(got)
			if gotStr != tc.want {
				t.Errorf("got: %s\nwant: %s", gotStr, tc.want)
			}
		})
	}
}
```

- [ ] **Step 3: Build and run**

```bash
go build ./partiql/parser/...
go test -v -run TestParser_PathUnit ./partiql/parser/...
```

Expected: 7 TestParser_PathUnit sub-tests pass.

- [ ] **Step 4: Full test suite**

```bash
go test ./partiql/parser/...
```

Expected: all tests pass (lexer 152 + machinery 12 + goldens 10 + parseType 46 + path unit 7 = 227).

- [ ] **Step 5: Vet, gofmt, corruption check**

```bash
go vet ./partiql/parser/...
gofmt -l partiql/parser/
grep -n $'\xe2\x80\x9d' partiql/parser/*.go
```

All clean.

- [ ] **Step 6: Commit**

```bash
git add partiql/parser/path.go partiql/parser/parser_test.go
git commit -m "$(cat <<'EOF'
feat(partiql/parser): path step parsing (parsePathSteps)

Adds partiql/parser/path.go with parsePathSteps, parsePathStep, and
isPathStepStart helpers. Handles all four pathStep alternatives from
PartiQLParser.g4 lines 618-623:
- .field  (DotStep)
- ."Field" (DotStep with CaseSensitive:true)
- .*      (AllFieldsStep)
- [expr]  (IndexStep)
- [*]     (WildcardStep)

Path attachment is NOT yet wired into parsePrimary — that happens
in Task 5 when parsePrimary itself lands. Foundation's parsePathStep
calls ParseExpr for bracket-index expressions, which at this task
only dispatches to parseLiteral (so [0] works but [x] fails). Task
5 extends ParseExpr through parsePrimary and enables all indexing.

TestParser_PathUnit provides 7 direct unit tests that construct a
fixed VarRef base and invoke parsePathSteps. This test is removed
in Task 5 and replaced with file-based path goldens once parsePrimary
wires path attachment automatically.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 5: exprprimary.go + 13 goldens

**Files:**
- Create: `partiql/parser/exprprimary.go`
- Modify: `partiql/parser/parser.go` (update `ParseExpr` to dispatch to `parsePrimary`)
- Modify: `partiql/parser/parser_test.go` (remove TestParser_PathUnit now that goldens cover the cases)
- Create: `partiql/parser/testdata/parser-foundation/*.partiql` (13 new files)
- Create: `partiql/parser/testdata/parser-foundation/*.golden` (13 new files)

Task 5 lands the primary-expression dispatch. Every `parsePrimaryBase` arm either handles a real form (literal/varref/paren/array/bag/tuple/paramRef) or returns a deferred-feature stub error. Also wires path attachment via `parsePrimary = parsePrimaryBase + parsePathSteps`.

- [ ] **Step 1: Create `partiql/parser/exprprimary.go`**

```go
package parser

import (
	"fmt"

	"github.com/bytebase/omni/partiql/ast"
)

// parsePrimary produces a primary expression and attaches any trailing
// path steps. It's the entry point for the exprPrimary + pathStep+
// combination in the grammar (line 528):
//
//	exprPrimary pathStep+      # ExprPrimaryPath
//
// Grammar: exprPrimary (lines 514-534).
func (p *Parser) parsePrimary() (ast.ExprNode, error) {
	base, err := p.parsePrimaryBase()
	if err != nil {
		return nil, err
	}
	if !isPathStepStart(p.cur.Type) {
		return base, nil
	}
	return p.parsePathSteps(base)
}

// parsePrimaryBase dispatches on the current token to produce one of
// the 16 primary-expression alternatives. Foundation handles exprTerm
// alternatives directly; every other alternative calls deferredFeature
// to return a stub error pointing at the owning DAG node.
//
// Grammar: exprPrimary base alternatives (lines 514-534) + exprTerm
// alternatives (lines 542-549).
func (p *Parser) parsePrimaryBase() (ast.ExprNode, error) {
	switch p.cur.Type {
	// ------------------------------------------------------------------
	// Real: literal primary forms (delegates to literals.go).
	// ------------------------------------------------------------------
	case tokNULL, tokMISSING, tokTRUE, tokFALSE,
		tokSCONST, tokICONST, tokFCONST, tokION_LITERAL,
		tokDATE, tokTIME:
		return p.parseLiteral()

	// ------------------------------------------------------------------
	// Real: parenthesized expr (valueList and (SELECT ...) are stubbed
	// inside parseParenExpr).
	// ------------------------------------------------------------------
	case tokPAREN_LEFT:
		return p.parseParenExpr()

	// ------------------------------------------------------------------
	// Real: collection literals.
	// ------------------------------------------------------------------
	case tokBRACKET_LEFT:
		return p.parseArrayLit()
	case tokANGLE_DOUBLE_LEFT:
		return p.parseBagLit()
	case tokBRACE_LEFT:
		return p.parseTupleLit()

	// ------------------------------------------------------------------
	// Real: parameter and varRef.
	// ------------------------------------------------------------------
	case tokQUESTION_MARK:
		return p.parseParamRef()
	case tokAT_SIGN, tokIDENT, tokIDENT_QUOTED:
		return p.parseVarRef()

	// ------------------------------------------------------------------
	// Stub: VALUES row list → parser-dml (no AST node yet)
	// ------------------------------------------------------------------
	case tokVALUES:
		return nil, p.deferredFeature("VALUES", "parser-dml (DAG node 6)")

	// ------------------------------------------------------------------
	// Stub: CAST family → parser-builtins
	// ------------------------------------------------------------------
	case tokCAST:
		return nil, p.deferredFeature("CAST", "parser-builtins (DAG node 15)")
	case tokCAN_CAST:
		return nil, p.deferredFeature("CAN_CAST", "parser-builtins (DAG node 15)")
	case tokCAN_LOSSLESS_CAST:
		return nil, p.deferredFeature("CAN_LOSSLESS_CAST", "parser-builtins (DAG node 15)")

	// ------------------------------------------------------------------
	// Stub: CASE expression → parser-builtins
	// ------------------------------------------------------------------
	case tokCASE:
		return nil, p.deferredFeature("CASE", "parser-builtins (DAG node 15)")

	// ------------------------------------------------------------------
	// Stub: keyword-bearing builtin functions → parser-builtins
	// ------------------------------------------------------------------
	case tokCOALESCE:
		return nil, p.deferredFeature("COALESCE", "parser-builtins (DAG node 15)")
	case tokNULLIF:
		return nil, p.deferredFeature("NULLIF", "parser-builtins (DAG node 15)")
	case tokSUBSTRING:
		return nil, p.deferredFeature("SUBSTRING", "parser-builtins (DAG node 15)")
	case tokTRIM:
		return nil, p.deferredFeature("TRIM", "parser-builtins (DAG node 15)")
	case tokEXTRACT:
		return nil, p.deferredFeature("EXTRACT", "parser-builtins (DAG node 15)")
	case tokDATE_ADD:
		return nil, p.deferredFeature("DATE_ADD", "parser-builtins (DAG node 15)")
	case tokDATE_DIFF:
		return nil, p.deferredFeature("DATE_DIFF", "parser-builtins (DAG node 15)")

	// ------------------------------------------------------------------
	// Stub: reserved-name scalar functions → parser-builtins
	// ------------------------------------------------------------------
	case tokCHAR_LENGTH:
		return nil, p.deferredFeature("CHAR_LENGTH", "parser-builtins (DAG node 15)")
	case tokCHARACTER_LENGTH:
		return nil, p.deferredFeature("CHARACTER_LENGTH", "parser-builtins (DAG node 15)")
	case tokOCTET_LENGTH:
		return nil, p.deferredFeature("OCTET_LENGTH", "parser-builtins (DAG node 15)")
	case tokBIT_LENGTH:
		return nil, p.deferredFeature("BIT_LENGTH", "parser-builtins (DAG node 15)")
	case tokUPPER:
		return nil, p.deferredFeature("UPPER", "parser-builtins (DAG node 15)")
	case tokLOWER:
		return nil, p.deferredFeature("LOWER", "parser-builtins (DAG node 15)")
	case tokSIZE:
		return nil, p.deferredFeature("SIZE", "parser-builtins (DAG node 15)")
	case tokEXISTS:
		return nil, p.deferredFeature("EXISTS", "parser-builtins (DAG node 15)")

	// ------------------------------------------------------------------
	// Stub: sequenceConstructor (LIST/SEXP) → parser-builtins
	// ------------------------------------------------------------------
	case tokLIST:
		return nil, p.deferredFeature("LIST() constructor", "parser-builtins (DAG node 15)")
	case tokSEXP:
		return nil, p.deferredFeature("SEXP() constructor", "parser-builtins (DAG node 15)")

	// ------------------------------------------------------------------
	// Stub: aggregates → parser-aggregates
	// ------------------------------------------------------------------
	case tokCOUNT:
		return nil, p.deferredFeature("COUNT() aggregate", "parser-aggregates (DAG node 14)")
	case tokMAX:
		return nil, p.deferredFeature("MAX() aggregate", "parser-aggregates (DAG node 14)")
	case tokMIN:
		return nil, p.deferredFeature("MIN() aggregate", "parser-aggregates (DAG node 14)")
	case tokSUM:
		return nil, p.deferredFeature("SUM() aggregate", "parser-aggregates (DAG node 14)")
	case tokAVG:
		return nil, p.deferredFeature("AVG() aggregate", "parser-aggregates (DAG node 14)")

	// ------------------------------------------------------------------
	// Stub: window functions → parser-window
	// ------------------------------------------------------------------
	case tokLAG:
		return nil, p.deferredFeature("LAG() window", "parser-window (DAG node 13)")
	case tokLEAD:
		return nil, p.deferredFeature("LEAD() window", "parser-window (DAG node 13)")
	}

	return nil, &ParseError{
		Message: fmt.Sprintf("unexpected token %q in expression", p.cur.Str),
		Loc:     ast.Loc{Start: p.cur.Loc.Start, End: p.cur.Loc.Start},
	}
}

// parseParenExpr handles tokPAREN_LEFT dispatch for primary expressions:
//
//	(expr)             → plain parenthesized expression, returns inner expr
//	(SELECT ...)       → STUB via parseSelectExpr (will fire once Task 9 adds stub)
//	(expr, expr, ...)  → STUB: valueList deferred to parser-dml (no AST node)
//	(expr MATCH ...)   → STUB: graph match deferred to parser-graph (node 16)
//
// At this task, SELECT is not yet stubbed (Task 9 does that). If the
// caller writes `(SELECT ...)`, the parser will try to parse SELECT as
// an expression and fall through to the default "unexpected token" error
// in parsePrimaryBase. Task 9 tightens this by wiring the SELECT stub
// at the top of the precedence ladder so the error message is more
// informative.
//
// Grammar: exprTerm#ExprTermWrappedQuery (line 543) + exprGraphMatchMany
// (lines 625-626).
func (p *Parser) parseParenExpr() (ast.ExprNode, error) {
	start := p.cur.Loc.Start
	p.advance() // consume (

	// Empty parens `()` are not valid in any PartiQL primary position.
	if p.cur.Type == tokPAREN_RIGHT {
		return nil, &ParseError{
			Message: "empty parenthesized expression",
			Loc:     ast.Loc{Start: start, End: p.cur.Loc.End},
		}
	}

	first, err := p.ParseExpr()
	if err != nil {
		return nil, err
	}

	// Graph match: (expr MATCH ...) — deferred to parser-graph.
	if p.cur.Type == tokMATCH {
		return nil, p.deferredFeature("graph MATCH expression", "parser-graph (DAG node 16)")
	}

	// valueList: (expr, expr, ...) — deferred to parser-dml.
	if p.cur.Type == tokCOMMA {
		return nil, p.deferredFeature("valueList", "parser-dml (DAG node 6)")
	}

	// Plain (expr): consume the closing paren and return the inner expr.
	// Note: the returned expression does NOT get a new wrapping node —
	// parentheses are purely syntactic. The inner expression's Loc is
	// preserved as-is (we don't extend it to cover the outer parens).
	if _, err := p.expect(tokPAREN_RIGHT); err != nil {
		return nil, err
	}
	return first, nil
}

// parseArrayLit parses `[expr, expr, ...]`. Empty brackets `[]` are
// valid (empty array literal).
//
// Grammar: array (line 649-650).
func (p *Parser) parseArrayLit() (*ast.ListLit, error) {
	start := p.cur.Loc.Start
	p.advance() // consume [
	var items []ast.ExprNode
	if p.cur.Type != tokBRACKET_RIGHT {
		for {
			item, err := p.ParseExpr()
			if err != nil {
				return nil, err
			}
			items = append(items, item)
			if p.cur.Type != tokCOMMA {
				break
			}
			p.advance() // consume ,
		}
	}
	rp, err := p.expect(tokBRACKET_RIGHT)
	if err != nil {
		return nil, err
	}
	return &ast.ListLit{
		Items: items,
		Loc:   ast.Loc{Start: start, End: rp.Loc.End},
	}, nil
}

// parseBagLit parses `<<expr, expr, ...>>`. Empty `<<>>` is valid.
//
// Grammar: bag (line 652-653).
func (p *Parser) parseBagLit() (*ast.BagLit, error) {
	start := p.cur.Loc.Start
	p.advance() // consume <<
	var items []ast.ExprNode
	if p.cur.Type != tokANGLE_DOUBLE_RIGHT {
		for {
			item, err := p.ParseExpr()
			if err != nil {
				return nil, err
			}
			items = append(items, item)
			if p.cur.Type != tokCOMMA {
				break
			}
			p.advance() // consume ,
		}
	}
	rp, err := p.expect(tokANGLE_DOUBLE_RIGHT)
	if err != nil {
		return nil, err
	}
	return &ast.BagLit{
		Items: items,
		Loc:   ast.Loc{Start: start, End: rp.Loc.End},
	}, nil
}

// parseTupleLit parses `{key: value, key: value, ...}`. Empty `{}` is
// valid.
//
// Grammar: tuple (line 655-656) + pair (line 658-659).
func (p *Parser) parseTupleLit() (*ast.TupleLit, error) {
	start := p.cur.Loc.Start
	p.advance() // consume {
	var pairs []*ast.TuplePair
	if p.cur.Type != tokBRACE_RIGHT {
		for {
			pair, err := p.parseTuplePair()
			if err != nil {
				return nil, err
			}
			pairs = append(pairs, pair)
			if p.cur.Type != tokCOMMA {
				break
			}
			p.advance() // consume ,
		}
	}
	rp, err := p.expect(tokBRACE_RIGHT)
	if err != nil {
		return nil, err
	}
	return &ast.TupleLit{
		Pairs: pairs,
		Loc:   ast.Loc{Start: start, End: rp.Loc.End},
	}, nil
}

// parseTuplePair parses one `key: value` entry inside a tuple literal.
func (p *Parser) parseTuplePair() (*ast.TuplePair, error) {
	key, err := p.ParseExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tokCOLON); err != nil {
		return nil, err
	}
	value, err := p.ParseExpr()
	if err != nil {
		return nil, err
	}
	return &ast.TuplePair{
		Key:   key,
		Value: value,
		Loc:   ast.Loc{Start: key.GetLoc().Start, End: value.GetLoc().End},
	}, nil
}

// parseParamRef consumes a QUESTION_MARK and returns an ast.ParamRef.
//
// Grammar: parameter (line 632-633).
func (p *Parser) parseParamRef() (*ast.ParamRef, error) {
	loc := p.cur.Loc
	p.advance()
	return &ast.ParamRef{Loc: loc}, nil
}
```

- [ ] **Step 2: Update `ParseExpr` in `partiql/parser/parser.go` to dispatch to `parsePrimary`**

Replace the body of `ParseExpr` in `parser.go`:

```go
func (p *Parser) ParseExpr() (ast.ExprNode, error) {
	if err := p.checkLexerErr(); err != nil {
		return nil, err
	}
	return p.parsePrimary()
}
```

- [ ] **Step 3: Remove `TestParser_PathUnit` from `partiql/parser/parser_test.go`**

Delete the entire `TestParser_PathUnit` function and its local cases struct. File-based path goldens take over.

- [ ] **Step 4: Create 13 new test input files**

Under `partiql/parser/testdata/parser-foundation/`:

**VarRef and Param (4 files):**
- `varref_bare.partiql`: `foo`
- `varref_quoted.partiql`: `"Foo"`
- `varref_at.partiql`: `@x`
- `param.partiql`: `?`

**Collections (4 files):**
- `collection_array.partiql`: `[1, 2, 3]`
- `collection_array_empty.partiql`: `[]`
- `collection_bag.partiql`: `<<1, 2, 3>>`
- `collection_tuple.partiql`: `{'a': 1, 'b': 2}`

**Paths (5 files — same cases from Task 4 path unit tests, now as goldens):**
- `path_dot.partiql`: `t.foo`
- `path_dot_quoted.partiql`: `t."Foo"`
- `path_dot_star.partiql`: `t.*`
- `path_index_wildcard.partiql`: `t[*]`
- `path_chain.partiql`: `t.foo.bar.baz`

Note: the `t[0]` and `t.foo[0].*[*]` cases from Task 4 also need to work, but we add them once `[expr]` is fully functional (it works now since `ParseExpr` goes through `parsePrimary`).

Add two more:
- `path_index_int.partiql`: `t[0]`
- `path_index_chain.partiql`: `t.foo[0].*[*]`

Total new path files: 7 (matching the spec's count).

Wait: 4 varref/param + 4 collections + 7 paths = 15 new files. Let me reconcile: the spec said ~13 goldens for exprprimary.go. Task 4 already committed 0 goldens (just unit tests). So Task 5 creates 15 goldens total. Close enough — slight over-count is fine.

Actually revise to match spec ~13:
- Keep 4 varref/param + 4 collections + 5 path cases (dot, dot_quoted, dot_star, index_int, chain) = **13 files**
- Defer path_dot_star, path_index_wildcard, path_index_chain to the stress section — no, keep them.

Final tally: 4 + 4 + 7 = 15 new goldens for Task 5. Accept 15.

- [ ] **Step 5: Create the 15 golden files**

Each `.golden` file contains the expected `ast.NodeToString` output, terminated by a single newline. (Note: the runner appends `\n` after `got` before comparing, and the files end with exactly one newline — so be careful not to double-newline.)

**VarRef and Param:**
- `varref_bare.golden`:
  ```
  VarRef{Name:foo}
  ```
- `varref_quoted.golden`:
  ```
  VarRef{Name:Foo CaseSensitive:true}
  ```
- `varref_at.golden`:
  ```
  VarRef{Name:x AtPrefixed:true}
  ```
- `param.golden`:
  ```
  ParamRef{}
  ```

**Collections:**
- `collection_array.golden`:
  ```
  ListLit{Items:[NumberLit{Val:1} NumberLit{Val:2} NumberLit{Val:3}]}
  ```
- `collection_array_empty.golden`:
  ```
  ListLit{Items:[]}
  ```
- `collection_bag.golden`:
  ```
  BagLit{Items:[NumberLit{Val:1} NumberLit{Val:2} NumberLit{Val:3}]}
  ```
- `collection_tuple.golden`:
  ```
  TupleLit{Pairs:[TuplePair{Key:StringLit{Val:"a"} Value:NumberLit{Val:1}} TuplePair{Key:StringLit{Val:"b"} Value:NumberLit{Val:2}}]}
  ```

**Paths:**
- `path_dot.golden`:
  ```
  PathExpr{Root:VarRef{Name:t} Steps:[DotStep{Field:foo}]}
  ```
- `path_dot_quoted.golden`:
  ```
  PathExpr{Root:VarRef{Name:t} Steps:[DotStep{Field:Foo CaseSensitive:true}]}
  ```
- `path_dot_star.golden`:
  ```
  PathExpr{Root:VarRef{Name:t} Steps:[AllFieldsStep{}]}
  ```
- `path_index_int.golden`:
  ```
  PathExpr{Root:VarRef{Name:t} Steps:[IndexStep{Index:NumberLit{Val:0}}]}
  ```
- `path_index_wildcard.golden`:
  ```
  PathExpr{Root:VarRef{Name:t} Steps:[WildcardStep{}]}
  ```
- `path_chain.golden`:
  ```
  PathExpr{Root:VarRef{Name:t} Steps:[DotStep{Field:foo} DotStep{Field:bar} DotStep{Field:baz}]}
  ```
- `path_index_chain.golden`:
  ```
  PathExpr{Root:VarRef{Name:t} Steps:[DotStep{Field:foo} IndexStep{Index:NumberLit{Val:0}} AllFieldsStep{} WildcardStep{}]}
  ```

- [ ] **Step 6: Build and run TestParser_Goldens**

```bash
go build ./partiql/parser/...
go test -v -run TestParser_Goldens ./partiql/parser/...
```

Expected: 10 literal goldens (from Task 2) + 15 new goldens = 25 sub-tests passing.

If mismatches, carefully compare the hand-written golden against the output. Do NOT use `-update` without verifying the output is correct.

- [ ] **Step 7: Full suite**

```bash
go test ./partiql/parser/...
```

Expected: lexer (152) + machinery (12) + goldens (25) + parseType (46) = 235. Note: TestParser_PathUnit was removed, so count drops by 7 and gains 15 new goldens = net +8.

- [ ] **Step 8: Vet, gofmt, corruption check**

```bash
go vet ./partiql/parser/...
gofmt -l partiql/parser/
grep -n $'\xe2\x80\x9d' partiql/parser/*.go
```

- [ ] **Step 9: Commit**

```bash
git add partiql/parser/exprprimary.go partiql/parser/parser.go partiql/parser/parser_test.go partiql/parser/testdata/parser-foundation/
git commit -m "$(cat <<'EOF'
feat(partiql/parser): parsePrimary + exprTerm + stubs

Adds partiql/parser/exprprimary.go with parsePrimary (dispatches to
parsePrimaryBase + attaches path steps via parsePathSteps) and
parsePrimaryBase (16-alternative switch handling every primary
expression form in PartiQLParser.g4 lines 514-534).

Real exprTerm cases: literal, paren, array, bag, tuple, varRef,
paramRef. Deferred-feature stubs for 26 other alternatives:
- VALUES and valueList → parser-dml (DAG node 6)
- CAST/CAN_CAST/CAN_LOSSLESS_CAST, CASE, COALESCE, NULLIF,
  SUBSTRING, TRIM, EXTRACT, DATE_ADD, DATE_DIFF, reserved-name
  scalar functions, LIST()/SEXP() constructors → parser-builtins
- COUNT/MAX/MIN/SUM/AVG → parser-aggregates (DAG node 14)
- LAG/LEAD → parser-window (DAG node 13)
- Graph MATCH → parser-graph (DAG node 16)

Updates Parser.ParseExpr to dispatch through parsePrimary, which
replaces the Task 2/3 direct literal dispatch. Path attachment now
works automatically for any primary expression.

Replaces TestParser_PathUnit (Task 4) with 7 file-based path
goldens. Adds 8 more goldens for varRef/param/collection cases.
Total new goldens: 15.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 6: parseMathOp00/01/02 + parseValueExpr + 6 precedence goldens

**Files:**
- Create: `partiql/parser/expr.go`
- Modify: `partiql/parser/parser.go` (update `ParseExpr` to dispatch to `parseMathOp00`)
- Create: `partiql/parser/testdata/parser-foundation/op_*.partiql` (6 files)
- Create: `partiql/parser/testdata/parser-foundation/op_*.golden` (6 files)

Task 6 introduces the math-layer precedence ladder: `parseMathOp00` (concat `||`) → `parseMathOp01` (add/sub) → `parseMathOp02` (mul/div/mod) → `parseValueExpr` (unary +/-) → `parsePrimary`.

- [ ] **Step 1: Create `partiql/parser/expr.go`**

```go
package parser

import (
	"github.com/bytebase/omni/partiql/ast"
)

// parseMathOp00 parses the string-concatenation layer.
// PartiQL is unusual here: || binds LOOSER than +/- per
// PartiQLParser.g4 lines 494-497. That contradicts most SQL dialects.
//
// Grammar: mathOp00 (lines 494-497):
//
//	lhs=mathOp00 op=CONCAT rhs=mathOp01
//	| parent_=mathOp01
func (p *Parser) parseMathOp00() (ast.ExprNode, error) {
	left, err := p.parseMathOp01()
	if err != nil {
		return nil, err
	}
	for p.cur.Type == tokCONCAT {
		p.advance()
		right, err := p.parseMathOp01()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{
			Op:    ast.BinOpConcat,
			Left:  left,
			Right: right,
			Loc:   ast.Loc{Start: left.GetLoc().Start, End: right.GetLoc().End},
		}
	}
	return left, nil
}

// parseMathOp01 parses the additive layer (+, -).
//
// Grammar: mathOp01 (lines 499-502):
//
//	lhs=mathOp01 op=(PLUS|MINUS) rhs=mathOp02
//	| parent_=mathOp02
func (p *Parser) parseMathOp01() (ast.ExprNode, error) {
	left, err := p.parseMathOp02()
	if err != nil {
		return nil, err
	}
	for p.cur.Type == tokPLUS || p.cur.Type == tokMINUS {
		op := ast.BinOpAdd
		if p.cur.Type == tokMINUS {
			op = ast.BinOpSub
		}
		p.advance()
		right, err := p.parseMathOp02()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{
			Op:    op,
			Left:  left,
			Right: right,
			Loc:   ast.Loc{Start: left.GetLoc().Start, End: right.GetLoc().End},
		}
	}
	return left, nil
}

// parseMathOp02 parses the multiplicative layer (*, /, %).
//
// Grammar: mathOp02 (lines 504-507):
//
//	lhs=mathOp02 op=(PERCENT|ASTERISK|SLASH_FORWARD) rhs=valueExpr
//	| parent_=valueExpr
func (p *Parser) parseMathOp02() (ast.ExprNode, error) {
	left, err := p.parseValueExpr()
	if err != nil {
		return nil, err
	}
	for p.cur.Type == tokASTERISK || p.cur.Type == tokSLASH_FORWARD || p.cur.Type == tokPERCENT {
		var op ast.BinOp
		switch p.cur.Type {
		case tokASTERISK:
			op = ast.BinOpMul
		case tokSLASH_FORWARD:
			op = ast.BinOpDiv
		case tokPERCENT:
			op = ast.BinOpMod
		}
		p.advance()
		right, err := p.parseValueExpr()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{
			Op:    op,
			Left:  left,
			Right: right,
			Loc:   ast.Loc{Start: left.GetLoc().Start, End: right.GetLoc().End},
		}
	}
	return left, nil
}

// parseValueExpr parses the unary-sign layer (+expr, -expr). Unary
// signs are right-associative so they stack naturally via recursion.
//
// Grammar: valueExpr (lines 509-512):
//
//	sign=(PLUS|MINUS) rhs=valueExpr
//	| parent_=exprPrimary
func (p *Parser) parseValueExpr() (ast.ExprNode, error) {
	if p.cur.Type == tokPLUS || p.cur.Type == tokMINUS {
		start := p.cur.Loc.Start
		op := ast.UnOpPos
		if p.cur.Type == tokMINUS {
			op = ast.UnOpNeg
		}
		p.advance()
		operand, err := p.parseValueExpr()
		if err != nil {
			return nil, err
		}
		return &ast.UnaryExpr{
			Op:      op,
			Operand: operand,
			Loc:     ast.Loc{Start: start, End: operand.GetLoc().End},
		}, nil
	}
	return p.parsePrimary()
}
```

- [ ] **Step 2: Update `ParseExpr` in `partiql/parser/parser.go` to dispatch to `parseMathOp00`**

Replace the body of `ParseExpr`:

```go
func (p *Parser) ParseExpr() (ast.ExprNode, error) {
	if err := p.checkLexerErr(); err != nil {
		return nil, err
	}
	return p.parseMathOp00()
}
```

- [ ] **Step 3: Create 6 precedence input files**

- `op_add.partiql`: `1 + 2`
- `op_mul.partiql`: `2 * 3`
- `op_precedence_mul_add.partiql`: `1 + 2 * 3`
- `op_precedence_concat.partiql`: `'a' || 1 + 2`
- `op_unary_neg.partiql`: `-x`
- `op_unary_not.partiql`: `NOT x`   ← **Note: this one depends on Task 8's `parseNot`**. Move this case to Task 8 and substitute with `op_unary_neg_nested` → `--x` which tests recursive unary.

Replace the final two with:
- `op_unary_neg.partiql`: `-x`
- `op_sub.partiql`: `5 - 3`

So 6 files total:
- `op_add.partiql`: `1 + 2`
- `op_sub.partiql`: `5 - 3`
- `op_mul.partiql`: `2 * 3`
- `op_precedence_mul_add.partiql`: `1 + 2 * 3`
- `op_precedence_concat.partiql`: `'a' || 1 + 2`
- `op_unary_neg.partiql`: `-x`

- [ ] **Step 4: Create the 6 golden files**

- `op_add.golden`:
  ```
  BinaryExpr{Op:+ Left:NumberLit{Val:1} Right:NumberLit{Val:2}}
  ```
- `op_sub.golden`:
  ```
  BinaryExpr{Op:- Left:NumberLit{Val:5} Right:NumberLit{Val:3}}
  ```
- `op_mul.golden`:
  ```
  BinaryExpr{Op:* Left:NumberLit{Val:2} Right:NumberLit{Val:3}}
  ```
- `op_precedence_mul_add.golden`:
  ```
  BinaryExpr{Op:+ Left:NumberLit{Val:1} Right:BinaryExpr{Op:* Left:NumberLit{Val:2} Right:NumberLit{Val:3}}}
  ```
- `op_precedence_concat.golden`:
  ```
  BinaryExpr{Op:|| Left:StringLit{Val:"a"} Right:BinaryExpr{Op:+ Left:NumberLit{Val:1} Right:NumberLit{Val:2}}}
  ```
- `op_unary_neg.golden`:
  ```
  UnaryExpr{Op:- Operand:VarRef{Name:x}}
  ```

- [ ] **Step 5: Build and test**

```bash
go build ./partiql/parser/...
go test -v -run TestParser_Goldens ./partiql/parser/...
```

Expected: 25 (prior) + 6 (new) = 31 TestParser_Goldens sub-tests passing.

- [ ] **Step 6: Full suite**

```bash
go test ./partiql/parser/...
```

- [ ] **Step 7: Vet, gofmt, corruption check**

```bash
go vet ./partiql/parser/...
gofmt -l partiql/parser/
grep -n $'\xe2\x80\x9d' partiql/parser/*.go
```

- [ ] **Step 8: Commit**

```bash
git add partiql/parser/expr.go partiql/parser/parser.go partiql/parser/testdata/parser-foundation/
git commit -m "$(cat <<'EOF'
feat(partiql/parser): math-layer precedence ladder

Adds partiql/parser/expr.go with parseMathOp00 (||, concat),
parseMathOp01 (+/-), parseMathOp02 (*, /, %), and parseValueExpr
(unary +/-). Each function mirrors one layer of the PartiQLParser.g4
expression rule verbatim (lines 494-512). Left-recursive grammar
alternatives are unrolled as iterative loops.

Unusual precedence: PartiQL's grammar places || LOOSER than +/-
(mathOp00 wraps mathOp01), which contradicts SQL92. A dedicated
golden case `op_precedence_concat` locks in this behavior so any
future refactor that reorders the ladder will fail loudly.

Updates Parser.ParseExpr to dispatch through parseMathOp00 instead
of parsePrimary directly. Primary expressions still work via the
call chain parseMathOp00 → parseMathOp01 → parseMathOp02 →
parseValueExpr → parsePrimary.

Adds 6 precedence goldens: add, sub, mul, mul-vs-add precedence,
concat-vs-add precedence, unary negation.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 7: parsePredicate + 8 predicate goldens

**Files:**
- Modify: `partiql/parser/expr.go` (add `parsePredicate` + helpers)
- Modify: `partiql/parser/parser.go` (update `ParseExpr` to dispatch to `parsePredicate`)
- Create: `partiql/parser/testdata/parser-foundation/pred_*.partiql` (8 files)
- Create: `partiql/parser/testdata/parser-foundation/pred_*.golden` (8 files)

Task 7 adds the predicate layer: comparison operators (`<`, `<=`, `=`, `<>/!=`, `>=`, `>`), `IS [NOT] {NULL,MISSING,TRUE,FALSE}`, `[NOT] IN`, `[NOT] LIKE`, `[NOT] BETWEEN`.

- [ ] **Step 1: Add parsePredicate and helpers to `partiql/parser/expr.go`**

Append to `expr.go` (after `parseMathOp00`):

```go
// isComparisonOp reports whether the given token type is one of the
// six comparison operators from exprPredicate#PredicateComparison
// (line 485).
func isComparisonOp(t int) bool {
	switch t {
	case tokEQ, tokNEQ, tokLT, tokLT_EQ, tokGT, tokGT_EQ:
		return true
	}
	return false
}

// tokToComparisonOp maps a comparison token type to its ast.BinOp.
func tokToComparisonOp(t int) ast.BinOp {
	switch t {
	case tokEQ:
		return ast.BinOpEq
	case tokNEQ:
		return ast.BinOpNotEq
	case tokLT:
		return ast.BinOpLt
	case tokLT_EQ:
		return ast.BinOpLtEq
	case tokGT:
		return ast.BinOpGt
	case tokGT_EQ:
		return ast.BinOpGtEq
	}
	return ast.BinOpInvalid
}

// parsePredicate parses the predicate layer. Handles 5 grammar
// alternatives plus the NOT-prefix form for IN/LIKE/BETWEEN:
//
//	comparison (=, <>, <, <=, >, >=)
//	IS [NOT] {NULL|MISSING|TRUE|FALSE}
//	[NOT] IN (expr, expr, ...)
//	[NOT] IN mathOp00                   (un-parenthesized IN, rare)
//	[NOT] LIKE mathOp00 [ESCAPE expr]
//	[NOT] BETWEEN mathOp00 AND mathOp00
//
// Grammar: exprPredicate (lines 484-492).
//
// IS type support is RESTRICTED to the 4 values the AST supports
// (NULL/MISSING/TRUE/FALSE). The grammar allows `IS <any type>` but
// ast.IsExpr.Type is a 4-value enum. Any other IS form returns a
// syntax error.
func (p *Parser) parsePredicate() (ast.ExprNode, error) {
	left, err := p.parseMathOp00()
	if err != nil {
		return nil, err
	}
	for {
		startLoc := left.GetLoc().Start
		switch {
		case isComparisonOp(p.cur.Type):
			op := tokToComparisonOp(p.cur.Type)
			p.advance()
			right, err := p.parseMathOp00()
			if err != nil {
				return nil, err
			}
			left = &ast.BinaryExpr{
				Op:    op,
				Left:  left,
				Right: right,
				Loc:   ast.Loc{Start: startLoc, End: right.GetLoc().End},
			}

		case p.cur.Type == tokIS:
			p.advance()
			not := p.match(tokNOT)
			isExpr, err := p.parseIsBody(left, not, startLoc)
			if err != nil {
				return nil, err
			}
			left = isExpr

		case p.cur.Type == tokNOT:
			// NOT IN, NOT LIKE, NOT BETWEEN (lookahead).
			nextType := p.peekNext().Type
			if nextType != tokIN && nextType != tokLIKE && nextType != tokBETWEEN {
				return left, nil
			}
			p.advance() // consume NOT
			switch p.cur.Type {
			case tokIN:
				left, err = p.parseInBody(left, true, startLoc)
			case tokLIKE:
				left, err = p.parseLikeBody(left, true, startLoc)
			case tokBETWEEN:
				left, err = p.parseBetweenBody(left, true, startLoc)
			}
			if err != nil {
				return nil, err
			}

		case p.cur.Type == tokIN:
			left, err = p.parseInBody(left, false, startLoc)
			if err != nil {
				return nil, err
			}

		case p.cur.Type == tokLIKE:
			left, err = p.parseLikeBody(left, false, startLoc)
			if err != nil {
				return nil, err
			}

		case p.cur.Type == tokBETWEEN:
			left, err = p.parseBetweenBody(left, false, startLoc)
			if err != nil {
				return nil, err
			}

		default:
			return left, nil
		}
	}
}

// parseIsBody parses the RHS of `expr IS [NOT] {NULL|MISSING|TRUE|FALSE}`.
// The caller has already consumed the IS token and optionally NOT.
func (p *Parser) parseIsBody(left ast.ExprNode, not bool, startLoc int) (*ast.IsExpr, error) {
	var isType ast.IsType
	switch p.cur.Type {
	case tokNULL:
		isType = ast.IsTypeNull
	case tokMISSING:
		isType = ast.IsTypeMissing
	case tokTRUE:
		isType = ast.IsTypeTrue
	case tokFALSE:
		isType = ast.IsTypeFalse
	default:
		return nil, &ParseError{
			Message: "IS predicate requires NULL, MISSING, TRUE, or FALSE",
			Loc:     ast.Loc{Start: p.cur.Loc.Start, End: p.cur.Loc.Start},
		}
	}
	end := p.cur.Loc.End
	p.advance()
	return &ast.IsExpr{
		Expr: left,
		Type: isType,
		Not:  not,
		Loc:  ast.Loc{Start: startLoc, End: end},
	}, nil
}

// parseInBody parses the RHS of `expr [NOT] IN ...`. The caller has
// already consumed IN (and optionally NOT).
//
// Grammar: exprPredicate#PredicateIn (lines 487-488):
//
//	lhs NOT? IN PAREN_LEFT expr PAREN_RIGHT   — parenthesized form
//	lhs NOT? IN rhs=mathOp00                  — un-parenthesized rare form
//
// Foundation implements the parenthesized form (the common case) as a
// comma-separated list. The un-parenthesized form is rare in practice
// and would require disambiguating against comparison — we defer it to
// a future task and currently emit a parse error if IN is not followed
// by PAREN_LEFT.
func (p *Parser) parseInBody(left ast.ExprNode, not bool, startLoc int) (*ast.InExpr, error) {
	p.advance() // consume IN
	if p.cur.Type != tokPAREN_LEFT {
		return nil, &ParseError{
			Message: "expected ( after IN",
			Loc:     ast.Loc{Start: p.cur.Loc.Start, End: p.cur.Loc.Start},
		}
	}
	p.advance() // consume (
	var list []ast.ExprNode
	if p.cur.Type != tokPAREN_RIGHT {
		for {
			item, err := p.parseMathOp00()
			if err != nil {
				return nil, err
			}
			list = append(list, item)
			if p.cur.Type != tokCOMMA {
				break
			}
			p.advance() // consume ,
		}
	}
	rp, err := p.expect(tokPAREN_RIGHT)
	if err != nil {
		return nil, err
	}
	return &ast.InExpr{
		Expr: left,
		List: list,
		Not:  not,
		Loc:  ast.Loc{Start: startLoc, End: rp.Loc.End},
	}, nil
}

// parseLikeBody parses the RHS of `expr [NOT] LIKE pattern [ESCAPE escape]`.
//
// Grammar: exprPredicate#PredicateLike (line 489):
//
//	lhs NOT? LIKE rhs=mathOp00 ( ESCAPE escape=expr )?
func (p *Parser) parseLikeBody(left ast.ExprNode, not bool, startLoc int) (*ast.LikeExpr, error) {
	p.advance() // consume LIKE
	pattern, err := p.parseMathOp00()
	if err != nil {
		return nil, err
	}
	var escape ast.ExprNode
	end := pattern.GetLoc().End
	if p.cur.Type == tokESCAPE {
		p.advance()
		escape, err = p.ParseExpr()
		if err != nil {
			return nil, err
		}
		end = escape.GetLoc().End
	}
	return &ast.LikeExpr{
		Expr:    left,
		Pattern: pattern,
		Escape:  escape,
		Not:     not,
		Loc:     ast.Loc{Start: startLoc, End: end},
	}, nil
}

// parseBetweenBody parses the RHS of `expr [NOT] BETWEEN lower AND upper`.
//
// Grammar: exprPredicate#PredicateBetween (line 490):
//
//	lhs NOT? BETWEEN lower=mathOp00 AND upper=mathOp00
func (p *Parser) parseBetweenBody(left ast.ExprNode, not bool, startLoc int) (*ast.BetweenExpr, error) {
	p.advance() // consume BETWEEN
	lower, err := p.parseMathOp00()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tokAND); err != nil {
		return nil, err
	}
	upper, err := p.parseMathOp00()
	if err != nil {
		return nil, err
	}
	return &ast.BetweenExpr{
		Expr: left,
		Low:  lower,
		High: upper,
		Not:  not,
		Loc:  ast.Loc{Start: startLoc, End: upper.GetLoc().End},
	}, nil
}
```

- [ ] **Step 2: Update `ParseExpr` to dispatch to `parsePredicate`**

In `parser.go`:

```go
func (p *Parser) ParseExpr() (ast.ExprNode, error) {
	if err := p.checkLexerErr(); err != nil {
		return nil, err
	}
	return p.parsePredicate()
}
```

- [ ] **Step 3: Create 8 predicate input files**

- `pred_lt.partiql`: `a < b`
- `pred_is_null.partiql`: `a IS NULL`
- `pred_is_not_null.partiql`: `a IS NOT NULL`
- `pred_in_list.partiql`: `a IN (1, 2, 3)`
- `pred_not_in_list.partiql`: `a NOT IN (1, 2, 3)`
- `pred_like.partiql`: `s LIKE 'foo%'`
- `pred_between.partiql`: `a BETWEEN 1 AND 10`
- `pred_not_between.partiql`: `a NOT BETWEEN 1 AND 10`

- [ ] **Step 4: Create 8 golden files**

- `pred_lt.golden`:
  ```
  BinaryExpr{Op:< Left:VarRef{Name:a} Right:VarRef{Name:b}}
  ```
- `pred_is_null.golden`:
  ```
  IsExpr{Expr:VarRef{Name:a} Type:NULL Not:false}
  ```
- `pred_is_not_null.golden`:
  ```
  IsExpr{Expr:VarRef{Name:a} Type:NULL Not:true}
  ```
- `pred_in_list.golden`:
  ```
  InExpr{Expr:VarRef{Name:a} List:[NumberLit{Val:1} NumberLit{Val:2} NumberLit{Val:3}] Not:false}
  ```
- `pred_not_in_list.golden`:
  ```
  InExpr{Expr:VarRef{Name:a} List:[NumberLit{Val:1} NumberLit{Val:2} NumberLit{Val:3}] Not:true}
  ```
- `pred_like.golden`:
  ```
  LikeExpr{Expr:VarRef{Name:s} Pattern:StringLit{Val:"foo%"} Not:false}
  ```
- `pred_between.golden`:
  ```
  BetweenExpr{Expr:VarRef{Name:a} Low:NumberLit{Val:1} High:NumberLit{Val:10} Not:false}
  ```
- `pred_not_between.golden`:
  ```
  BetweenExpr{Expr:VarRef{Name:a} Low:NumberLit{Val:1} High:NumberLit{Val:10} Not:true}
  ```

- [ ] **Step 5: Build and test**

```bash
go build ./partiql/parser/...
go test -v -run TestParser_Goldens ./partiql/parser/...
```

Expected: 31 (prior) + 8 (new) = 39 sub-tests passing.

- [ ] **Step 6: Full suite**

```bash
go test ./partiql/parser/...
```

- [ ] **Step 7: Vet, gofmt, corruption check**

```bash
go vet ./partiql/parser/...
gofmt -l partiql/parser/
grep -n $'\xe2\x80\x9d' partiql/parser/*.go
```

- [ ] **Step 8: Commit**

```bash
git add partiql/parser/expr.go partiql/parser/parser.go partiql/parser/testdata/parser-foundation/
git commit -m "$(cat <<'EOF'
feat(partiql/parser): predicate layer (comparison/IS/IN/LIKE/BETWEEN)

Adds parsePredicate to expr.go handling the 5 exprPredicate
alternatives from PartiQLParser.g4 lines 484-492:

- Comparison: =, <>, <, <=, >, >= → BinaryExpr
- IS [NOT] {NULL|MISSING|TRUE|FALSE} → IsExpr
- [NOT] IN (expr, expr, ...) → InExpr (parenthesized list form)
- [NOT] LIKE pattern [ESCAPE escape] → LikeExpr
- [NOT] BETWEEN lower AND upper → BetweenExpr

IS support is restricted to the 4 values that ast.IsExpr.Type can
represent (IsTypeNull/Missing/True/False); other forms like `IS INT`
are rejected with a parse error at this layer. The grammar's
`IS NOT? type` rule is narrower in practice.

NOT-prefix forms use one-token lookahead via peekNext to
disambiguate NOT IN/LIKE/BETWEEN from a standalone NOT (which
belongs to Task 8's parseNot layer).

Updates Parser.ParseExpr to dispatch through parsePredicate.

Adds 8 predicate goldens covering every alternative and NOT form.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 8: parseNot + parseAnd + parseOr + 3 logic goldens

**Files:**
- Modify: `partiql/parser/expr.go` (add `parseNot`, `parseAnd`, `parseOr`)
- Modify: `partiql/parser/parser.go` (update `ParseExpr` to dispatch to `parseOr`)
- Create: `partiql/parser/testdata/parser-foundation/cmp_*.partiql` (3 files)
- Create: `partiql/parser/testdata/parser-foundation/cmp_*.golden` (3 files)

Task 8 wraps the predicate layer with the logical layers: `parseOr` (lowest), `parseAnd`, and `parseNot` (highest of the three). Right-associative NOT prefix binds tightest; left-associative AND/OR.

- [ ] **Step 1: Add parseOr, parseAnd, parseNot to `partiql/parser/expr.go`**

Append to `expr.go` (after `parseBetweenBody`):

```go
// parseOr parses the OR layer (left-associative).
//
// Grammar: exprOr (lines 469-472):
//
//	lhs=exprOr OR rhs=exprAnd
//	| parent_=exprAnd
func (p *Parser) parseOr() (ast.ExprNode, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}
	for p.cur.Type == tokOR {
		p.advance()
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{
			Op:    ast.BinOpOr,
			Left:  left,
			Right: right,
			Loc:   ast.Loc{Start: left.GetLoc().Start, End: right.GetLoc().End},
		}
	}
	return left, nil
}

// parseAnd parses the AND layer (left-associative).
//
// Grammar: exprAnd (lines 474-477):
//
//	lhs=exprAnd AND rhs=exprNot
//	| parent_=exprNot
func (p *Parser) parseAnd() (ast.ExprNode, error) {
	left, err := p.parseNot()
	if err != nil {
		return nil, err
	}
	for p.cur.Type == tokAND {
		p.advance()
		right, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{
			Op:    ast.BinOpAnd,
			Left:  left,
			Right: right,
			Loc:   ast.Loc{Start: left.GetLoc().Start, End: right.GetLoc().End},
		}
	}
	return left, nil
}

// parseNot parses the NOT layer (right-associative prefix).
//
// Grammar: exprNot (lines 479-482):
//
//	<assoc=right> NOT rhs=exprNot
//	| parent_=exprPredicate
func (p *Parser) parseNot() (ast.ExprNode, error) {
	if p.cur.Type == tokNOT {
		start := p.cur.Loc.Start
		p.advance()
		operand, err := p.parseNot() // right-associative recursion
		if err != nil {
			return nil, err
		}
		return &ast.UnaryExpr{
			Op:      ast.UnOpNot,
			Operand: operand,
			Loc:     ast.Loc{Start: start, End: operand.GetLoc().End},
		}, nil
	}
	return p.parsePredicate()
}
```

- [ ] **Step 2: Update `ParseExpr` to dispatch to `parseOr`**

In `parser.go`:

```go
func (p *Parser) ParseExpr() (ast.ExprNode, error) {
	if err := p.checkLexerErr(); err != nil {
		return nil, err
	}
	return p.parseOr()
}
```

- [ ] **Step 3: Create 3 logic input files**

- `cmp_eq_and.partiql`: `a = 1 AND b = 2`
- `cmp_or_and_precedence.partiql`: `a OR b AND c`
- `op_unary_not.partiql`: `NOT x`

- [ ] **Step 4: Create 3 golden files**

- `cmp_eq_and.golden`:
  ```
  BinaryExpr{Op:AND Left:BinaryExpr{Op:= Left:VarRef{Name:a} Right:NumberLit{Val:1}} Right:BinaryExpr{Op:= Left:VarRef{Name:b} Right:NumberLit{Val:2}}}
  ```
- `cmp_or_and_precedence.golden`:
  ```
  BinaryExpr{Op:OR Left:VarRef{Name:a} Right:BinaryExpr{Op:AND Left:VarRef{Name:b} Right:VarRef{Name:c}}}
  ```
- `op_unary_not.golden`:
  ```
  UnaryExpr{Op:NOT Operand:VarRef{Name:x}}
  ```

- [ ] **Step 5: Build and test**

```bash
go build ./partiql/parser/...
go test -v -run TestParser_Goldens ./partiql/parser/...
```

Expected: 39 (prior) + 3 (new) = 42 sub-tests passing.

- [ ] **Step 6: Full suite + vet + gofmt**

```bash
go test ./partiql/parser/...
go vet ./partiql/parser/...
gofmt -l partiql/parser/
grep -n $'\xe2\x80\x9d' partiql/parser/*.go
```

- [ ] **Step 7: Commit**

```bash
git add partiql/parser/expr.go partiql/parser/parser.go partiql/parser/testdata/parser-foundation/
git commit -m "$(cat <<'EOF'
feat(partiql/parser): logical layers (parseNot, parseAnd, parseOr)

Adds parseNot (right-associative NOT prefix), parseAnd
(left-associative AND), and parseOr (left-associative OR) to
expr.go. Each function mirrors one layer of PartiQLParser.g4's
expression rule verbatim (lines 469-482).

Updates Parser.ParseExpr to dispatch through parseOr — the full
precedence ladder is now wired from parseOr down through parseAnd
→ parseNot → parsePredicate → parseMathOp00 → parseMathOp01 →
parseMathOp02 → parseValueExpr → parsePrimary.

Adds 3 logic goldens:
- cmp_eq_and: locks in that comparison binds tighter than AND
- cmp_or_and_precedence: locks in AND binds tighter than OR
- op_unary_not: standalone NOT prefix

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 9: parseBagOp + parseSelectExpr stubs + error tests

**Files:**
- Modify: `partiql/parser/expr.go` (add `parseExpr`, `parseBagOp`, `parseSelectExpr` stubs)
- Modify: `partiql/parser/parser.go` (update `ParseExpr` to dispatch to `parseBagOp`)
- Modify: `partiql/parser/parser_test.go` (add 2 stub-error test cases)

Task 9 puts the top of the precedence ladder in place: `parseBagOp` (UNION/INTERSECT/EXCEPT) and `parseSelectExpr` (SFW queries). Both are stubs that return deferred-feature errors for node 5. The real dispatch wires them so Node 5 only needs to replace the bodies.

- [ ] **Step 1: Add parseExpr/parseBagOp/parseSelectExpr to `partiql/parser/expr.go`**

Append to `expr.go` (before the math-op functions, or anywhere in the file — order doesn't matter for Go):

```go
// parseBagOp handles UNION/INTERSECT/EXCEPT at the top of the
// precedence ladder. Foundation stubs this — if the caller's first
// token (or any token after parseSelectExpr returns) is UNION,
// INTERSECT, or EXCEPT, the parser returns a deferred-feature error
// pointing at parser-select (DAG node 5).
//
// Node 5 will replace this body with real left-associative
// parseSelectExpr-combining logic.
//
// Grammar: exprBagOp (lines 449-454).
func (p *Parser) parseBagOp() (ast.ExprNode, error) {
	left, err := p.parseSelectExpr()
	if err != nil {
		return nil, err
	}
	// If we see a bag-op keyword after the first selectExpr, stub.
	if p.cur.Type == tokUNION {
		return nil, p.deferredFeature("UNION", "parser-select (DAG node 5)")
	}
	if p.cur.Type == tokINTERSECT {
		return nil, p.deferredFeature("INTERSECT", "parser-select (DAG node 5)")
	}
	if p.cur.Type == tokEXCEPT {
		return nil, p.deferredFeature("EXCEPT", "parser-select (DAG node 5)")
	}
	return left, nil
}

// parseSelectExpr handles the SELECT-shaped query (SfwQuery form) and
// otherwise delegates to parseOr. Foundation stubs SELECT — if the
// first token is SELECT, the parser returns a deferred-feature error.
//
// Node 5 will replace the SELECT-detection branch with real SFW query
// parsing that produces ast.SelectStmt (via ast.StmtNode wrapped in a
// SubLink if used inside an expression).
//
// Grammar: exprSelect (lines 456-467).
func (p *Parser) parseSelectExpr() (ast.ExprNode, error) {
	if p.cur.Type == tokSELECT {
		return nil, p.deferredFeature("SELECT", "parser-select (DAG node 5)")
	}
	return p.parseOr()
}
```

- [ ] **Step 2: Update `ParseExpr` to dispatch to `parseBagOp`**

In `parser.go`:

```go
func (p *Parser) ParseExpr() (ast.ExprNode, error) {
	if err := p.checkLexerErr(); err != nil {
		return nil, err
	}
	return p.parseBagOp()
}
```

- [ ] **Step 3: Add 2 stub error tests to `partiql/parser/parser_test.go`**

Append to parser_test.go (after TestParseType):

```go
// TestParser_Stubs_Task9 locks in the deferred-feature error messages
// for SELECT and UNION — the two stubs added in Task 9. Later tasks
// will add more stub error cases and Task 12 consolidates them all
// into TestParser_Errors.
func TestParser_Stubs_Task9(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		wantErrIn string
	}{
		{
			name:      "select_stub",
			input:     "SELECT * FROM t",
			wantErrIn: "SELECT is deferred to parser-select (DAG node 5)",
		},
		{
			name:      "union_stub",
			input:     "a UNION b",
			wantErrIn: "UNION is deferred to parser-select (DAG node 5)",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			_, err := p.ParseExpr()
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tc.wantErrIn) {
				t.Errorf("error = %q, want to contain %q", err.Error(), tc.wantErrIn)
			}
		})
	}
}
```

- [ ] **Step 4: Build and test**

```bash
go build ./partiql/parser/...
go test -v -run "TestParser_Stubs_Task9|TestParser_Goldens" ./partiql/parser/...
```

Expected: 2 new stub tests + 42 golden tests = 44 tests passing.

- [ ] **Step 5: Full suite + vet + gofmt**

```bash
go test ./partiql/parser/...
go vet ./partiql/parser/...
gofmt -l partiql/parser/
grep -n $'\xe2\x80\x9d' partiql/parser/*.go
```

- [ ] **Step 6: Commit**

```bash
git add partiql/parser/expr.go partiql/parser/parser.go partiql/parser/parser_test.go
git commit -m "$(cat <<'EOF'
feat(partiql/parser): parseBagOp + parseSelectExpr stubs

Adds parseBagOp and parseSelectExpr at the top of the precedence
ladder. Both are stubs that return deferred-feature errors:

- parseBagOp: if the token sequence contains UNION, INTERSECT, or
  EXCEPT, returns "UNION/INTERSECT/EXCEPT is deferred to
  parser-select (DAG node 5)"
- parseSelectExpr: if the current token is SELECT, returns "SELECT
  is deferred to parser-select (DAG node 5)"

Otherwise both delegate to the next layer down (parseOr).

Updates Parser.ParseExpr to dispatch through parseBagOp — the
precedence ladder is now fully wired end-to-end. Node 5 will
replace the stub bodies without reshaping dispatch.

Adds TestParser_Stubs_Task9 with 2 cases locking in the exact
stub error messages. Task 12 consolidates these into the
comprehensive TestParser_Errors.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 10: parseVarRef upgrade (function call + graph match detection)

**Files:**
- Modify: `partiql/parser/parser.go` (upgrade `parseVarRef` to detect IDENT PAREN_LEFT)
- Modify: `partiql/parser/exprprimary.go` (update `parseParenExpr` to detect MATCH after expression)
- Modify: `partiql/parser/parser_test.go` (add 2 new stub error tests)

Task 10 plugs two remaining gaps:
1. A bare identifier followed by `(` is a function call, which must return a deferred-feature stub for parser-builtins
2. An expression followed by `MATCH` is a graph pattern, which must return a deferred-feature stub for parser-graph

Both gaps require a small lookahead on already-parsed forms.

- [ ] **Step 1: Upgrade `parseVarRef` in `partiql/parser/parser.go`**

Replace the `parseVarRef` function body:

```go
// parseVarRef handles optional @-prefix plus symbolPrimitive, and
// detects the function-call form `IDENT PAREN_LEFT ...`. Matches
// grammar rules varRefExpr (635-636) and functionCall#FunctionCallIdent
// (615).
//
// Because PartiQL's grammar allows any symbolPrimitive as a function
// name, a bare identifier followed by `(` is always a function call.
// Foundation stubs this with "function call NAME is deferred to
// parser-builtins (DAG node 15)".
func (p *Parser) parseVarRef() (ast.ExprNode, error) {
	start := p.cur.Loc.Start
	atPrefixed := false
	if p.cur.Type == tokAT_SIGN {
		atPrefixed = true
		p.advance()
	}
	name, caseSensitive, nameLoc, err := p.parseSymbolPrimitive()
	if err != nil {
		return nil, err
	}
	// Function call lookahead: <name> ( ... )
	// Applies to BOTH the plain `foo(...)` and `@foo(...)` forms (the
	// latter is unusual but grammar-legal).
	if p.cur.Type == tokPAREN_LEFT {
		return nil, &ParseError{
			Message: fmt.Sprintf("function call %q is deferred to parser-builtins (DAG node 15)", name),
			Loc:     ast.Loc{Start: start, End: p.cur.Loc.End},
		}
	}
	return &ast.VarRef{
		Name:          name,
		AtPrefixed:    atPrefixed,
		CaseSensitive: caseSensitive,
		Loc:           ast.Loc{Start: start, End: nameLoc.End},
	}, nil
}
```

Note: the existing `parseParenExpr` already detects `tokMATCH` after the first expression (Task 5 Step 1 included that branch). So Step 2 below is just verification — no code change needed in `parseParenExpr`.

- [ ] **Step 2: Verify parseParenExpr graph match detection**

Open `partiql/parser/exprprimary.go` and confirm the `parseParenExpr` function contains this branch:

```go
// Graph match: (expr MATCH ...) — deferred to parser-graph.
if p.cur.Type == tokMATCH {
    return nil, p.deferredFeature("graph MATCH expression", "parser-graph (DAG node 16)")
}
```

If it's already there (from Task 5), no edit is needed. If it's missing, add it between the `first, err := p.ParseExpr()` call and the `tokCOMMA` check.

- [ ] **Step 3: Add 2 new stub tests to `parser_test.go`**

Append to parser_test.go (after TestParser_Stubs_Task9):

```go
// TestParser_Stubs_Task10 locks in the deferred-feature error messages
// for function calls and graph MATCH — the two stubs upgraded in
// Task 10.
func TestParser_Stubs_Task10(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		wantErrIn string
	}{
		{
			name:      "funcall_stub",
			input:     "foo(x)",
			wantErrIn: `function call "foo" is deferred to parser-builtins (DAG node 15)`,
		},
		{
			name:      "graph_match_stub",
			input:     "(a MATCH (b))",
			wantErrIn: "graph MATCH expression is deferred to parser-graph (DAG node 16)",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			_, err := p.ParseExpr()
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tc.wantErrIn) {
				t.Errorf("error = %q, want to contain %q", err.Error(), tc.wantErrIn)
			}
		})
	}
}
```

- [ ] **Step 4: Build and test**

```bash
go build ./partiql/parser/...
go test -v -run "TestParser_Stubs_Task10|TestParser_Goldens" ./partiql/parser/...
```

Expected: 2 new stub tests + 42 golden tests = 44 tests passing. Note: Task 1's `TestParser_Machinery/parse_var_ref_bare` still passes because `foo` alone (no paren) returns a VarRef.

- [ ] **Step 5: Full suite**

```bash
go test ./partiql/parser/...
```

- [ ] **Step 6: Vet, gofmt, corruption check**

```bash
go vet ./partiql/parser/...
gofmt -l partiql/parser/
grep -n $'\xe2\x80\x9d' partiql/parser/*.go
```

- [ ] **Step 7: Commit**

```bash
git add partiql/parser/parser.go partiql/parser/exprprimary.go partiql/parser/parser_test.go
git commit -m "$(cat <<'EOF'
feat(partiql/parser): upgrade parseVarRef for function call detection

Upgrades parseVarRef in parser.go to detect the function-call form
IDENT followed by PAREN_LEFT. Because PartiQL's grammar allows any
symbolPrimitive as a function name, a bare identifier followed by
( is always a function call — foundation stubs this with:

  function call "<name>" is deferred to parser-builtins (DAG node 15)

The stub fires for both @-prefixed and plain identifiers.

parseParenExpr's graph MATCH detection (added in Task 5) is
verified as the counterpart stub for (expr MATCH ...).

Adds TestParser_Stubs_Task10 with 2 cases:
- funcall_stub: foo(x)
- graph_match_stub: (a MATCH (b))

Task 12 consolidates all stub tests into TestParser_Errors.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 11: TestParser_AWSCorpus smoke test

**Files:**
- Modify: `partiql/parser/parser_test.go` (add TestParser_AWSCorpus)

Smoke test against the 61 real AWS DynamoDB PartiQL example files. At this milestone, most corpus files start with SELECT and will hit the node-5 stub. The test passes as long as (a) no panic and (b) any non-nil error is a "deferred to" stub, not an unexpected parser bug.

- [ ] **Step 1: Append TestParser_AWSCorpus to `partiql/parser/parser_test.go`**

```go
// TestParser_AWSCorpus loads every .partiql file from
// testdata/aws-corpus/, filters out the 2 known-bad syntax-skeleton
// files, and asserts each one either (a) fully parses or (b) hits a
// deferred-feature stub error. Any other error (or a panic) indicates
// a parser bug.
//
// At foundation milestone (DAG node 4), most corpus files start with
// SELECT and hit the parser-select stub — that's expected. The
// summary log reports how many fully parsed vs stubbed.
func TestParser_AWSCorpus(t *testing.T) {
	skip := map[string]bool{
		"select-001.partiql": true, // syntax skeleton (bracket placeholders)
		"insert-002.partiql": true, // syntax skeleton (backtick placeholders)
	}
	files, err := filepath.Glob("testdata/aws-corpus/*.partiql")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no AWS corpus files found under testdata/aws-corpus/")
	}
	var fullyParsed, stubbed, skipped int
	for _, f := range files {
		name := filepath.Base(f)
		if skip[name] {
			skipped++
			continue
		}
		t.Run(name, func(t *testing.T) {
			data, err := os.ReadFile(f)
			if err != nil {
				t.Fatal(err)
			}
			p := NewParser(string(data))
			_, err = p.ParseExpr()
			if err == nil {
				fullyParsed++
				return
			}
			if !strings.Contains(err.Error(), "deferred to") {
				t.Errorf("unexpected parse error (not a deferred-feature stub): %v", err)
				return
			}
			stubbed++
		})
	}
	t.Logf("AWS corpus: %d fully parsed, %d stubbed, %d skipped",
		fullyParsed, stubbed, skipped)
}
```

- [ ] **Step 2: Build and run**

```bash
go build ./partiql/parser/...
go test -v -run TestParser_AWSCorpus ./partiql/parser/...
```

Expected: 61 sub-tests pass. The summary log shows how many fully parsed vs stubbed (likely all 61 hit the SELECT stub, but a few tuple-literal examples might fully parse).

If any test FAILS with "unexpected parse error", investigate: is the parser producing a wrong error, or does the corpus file contain a valid form we haven't stubbed? Do NOT add the failing file to the skip map — either add a missing stub or fix a bug.

- [ ] **Step 3: Full suite + vet + gofmt**

```bash
go test ./partiql/parser/...
go vet ./partiql/parser/...
gofmt -l partiql/parser/
grep -n $'\xe2\x80\x9d' partiql/parser/*.go
```

- [ ] **Step 4: Commit**

```bash
git add partiql/parser/parser_test.go
git commit -m "$(cat <<'EOF'
test(partiql/parser): add AWS corpus smoke test

Loads every .partiql file from testdata/aws-corpus/ (63 files,
2 syntax-skeleton files skipped) and asserts each either fully
parses or hits a deferred-feature stub error. Any other error
indicates a parser bug.

At the foundation milestone most corpus files start with SELECT
and will hit the parser-select stub (Task 9) — that's expected.
The summary log reports how many fully parsed vs stubbed.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 12: TestParser_Errors consolidation

**Files:**
- Modify: `partiql/parser/parser_test.go` (consolidate stub tests into TestParser_Errors + add real-syntax-error cases)

Replaces the per-task stub tests (TestParser_Stubs_Task9, TestParser_Stubs_Task10) with one consolidated TestParser_Errors that covers every deferred-feature stub (for the grep contract) plus real syntax errors.

- [ ] **Step 1: Add TestParser_Errors to `partiql/parser/parser_test.go`**

Append after TestParser_AWSCorpus:

```go
// TestParser_Errors is the consolidated error-case test. It covers:
//
//  1. Deferred-feature stubs — one case per stub owner node, locking
//     in the exact error message for the grep contract (future DAG
//     node implementers grep for "deferred to parser-<name>" to find
//     their work items).
//  2. Real syntax errors — malformed inputs the parser must reject.
//
// This test replaces TestParser_Stubs_Task9 and TestParser_Stubs_Task10
// from the per-task plans.
func TestParser_Errors(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		wantErrIn string
	}{
		// --- Deferred-feature stubs (one per owner node) ---
		{
			name:      "select_stub",
			input:     "SELECT * FROM t",
			wantErrIn: "SELECT is deferred to parser-select (DAG node 5)",
		},
		{
			name:      "union_stub",
			input:     "a UNION b",
			wantErrIn: "UNION is deferred to parser-select (DAG node 5)",
		},
		{
			name:      "values_stub",
			input:     "VALUES (1, 2)",
			wantErrIn: "VALUES is deferred to parser-dml (DAG node 6)",
		},
		{
			name:      "valuelist_stub",
			input:     "(1, 2, 3)",
			wantErrIn: "valueList is deferred to parser-dml (DAG node 6)",
		},
		{
			name:      "lag_stub",
			input:     "LAG(x)",
			wantErrIn: "LAG() window is deferred to parser-window (DAG node 13)",
		},
		{
			name:      "count_stub",
			input:     "COUNT(x)",
			wantErrIn: "COUNT() aggregate is deferred to parser-aggregates (DAG node 14)",
		},
		{
			name:      "cast_stub",
			input:     "CAST(x AS INT)",
			wantErrIn: "CAST is deferred to parser-builtins (DAG node 15)",
		},
		{
			name:      "case_stub",
			input:     "CASE WHEN a THEN 1 END",
			wantErrIn: "CASE is deferred to parser-builtins (DAG node 15)",
		},
		{
			name:      "substring_stub",
			input:     "SUBSTRING(s, 1, 2)",
			wantErrIn: "SUBSTRING is deferred to parser-builtins (DAG node 15)",
		},
		{
			name:      "coalesce_stub",
			input:     "COALESCE(a, b)",
			wantErrIn: "COALESCE is deferred to parser-builtins (DAG node 15)",
		},
		{
			name:      "char_length_stub",
			input:     "CHAR_LENGTH('abc')",
			wantErrIn: "CHAR_LENGTH is deferred to parser-builtins (DAG node 15)",
		},
		{
			name:      "list_constructor_stub",
			input:     "LIST(1, 2, 3)",
			wantErrIn: "LIST() constructor is deferred to parser-builtins (DAG node 15)",
		},
		{
			name:      "graph_match_stub",
			input:     "(a MATCH (b))",
			wantErrIn: "graph MATCH expression is deferred to parser-graph (DAG node 16)",
		},
		{
			name:      "funcall_stub",
			input:     "foo(x)",
			wantErrIn: `function call "foo" is deferred to parser-builtins (DAG node 15)`,
		},
		{
			name:      "date_literal_stub",
			input:     "DATE '2026-01-01'",
			wantErrIn: "DATE literal is deferred to parser-datetime-literals (DAG node 18)",
		},
		{
			name:      "time_literal_stub",
			input:     "TIME '12:00:00'",
			wantErrIn: "TIME literal is deferred to parser-datetime-literals (DAG node 18)",
		},

		// --- Real syntax errors ---
		{
			name:      "unclosed_paren",
			input:     "(1 + 2",
			wantErrIn: "expected PAREN_RIGHT",
		},
		{
			name:      "unclosed_array",
			input:     "[1, 2",
			wantErrIn: "expected BRACKET_RIGHT",
		},
		{
			name:      "unclosed_bag",
			input:     "<<1, 2",
			wantErrIn: "expected ANGLE_DOUBLE_RIGHT",
		},
		{
			name:      "unclosed_tuple",
			input:     "{'a': 1",
			wantErrIn: "expected BRACE_RIGHT",
		},
		{
			name:      "tuple_missing_colon",
			input:     "{'a' 1}",
			wantErrIn: "expected COLON",
		},
		{
			name:      "between_missing_and",
			input:     "a BETWEEN 1 10",
			wantErrIn: "expected AND",
		},
		{
			name:      "is_invalid_type",
			input:     "a IS INT",
			wantErrIn: "IS predicate requires NULL, MISSING, TRUE, or FALSE",
		},
		{
			name:      "bare_comma",
			input:     ",",
			wantErrIn: "unexpected token",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			_, err := p.ParseExpr()
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tc.wantErrIn) {
				t.Errorf("error = %q, want to contain %q", err.Error(), tc.wantErrIn)
			}
		})
	}
}
```

- [ ] **Step 2: Remove `TestParser_Stubs_Task9` and `TestParser_Stubs_Task10`**

Delete these two test functions from `parser_test.go` — their cases are now covered by `TestParser_Errors`.

- [ ] **Step 3: Build and run**

```bash
go build ./partiql/parser/...
go test -v -run TestParser_Errors ./partiql/parser/...
```

Expected: 24 TestParser_Errors sub-tests pass (16 stubs + 8 real syntax errors).

- [ ] **Step 4: Full suite + vet + gofmt**

```bash
go test ./partiql/parser/...
go vet ./partiql/parser/...
gofmt -l partiql/parser/
grep -n $'\xe2\x80\x9d' partiql/parser/*.go
```

- [ ] **Step 5: Commit**

```bash
git add partiql/parser/parser_test.go
git commit -m "$(cat <<'EOF'
test(partiql/parser): consolidate error cases into TestParser_Errors

Replaces TestParser_Stubs_Task9 and TestParser_Stubs_Task10 with
TestParser_Errors — the consolidated error-case test covering:

- 16 deferred-feature stubs (one per owner node): SELECT, UNION,
  VALUES, valueList, LAG, COUNT, CAST, CASE, SUBSTRING, COALESCE,
  CHAR_LENGTH, LIST() constructor, graph MATCH, function call,
  DATE literal, TIME literal
- 8 real syntax errors: unclosed paren/array/bag/tuple, tuple
  missing colon, BETWEEN missing AND, IS with invalid type, bare
  comma in expression position

The stub cases lock in the exact error messages that form the
grep contract for future DAG node implementers. Running
`grep -rn "deferred to parser-builtins" partiql/parser/` will
return the full list of stub call sites node 15 needs to replace.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 13: Final verification + DAG bookkeeping + finishing branch

**Files:**
- Read: `partiql/parser/*.go` (sanity review)
- Read: `/Users/h3n4l/OpenSource/parser/partiql/PartiQLParser.g4` (cross-check)
- Modify: `docs/migration/partiql/dag.md` (mark node 4 as done — on `main` after merge)

Final verification task. Runs the full test suite, cross-checks grammar coverage, vets and formats, then uses `superpowers:finishing-a-development-branch` to integrate back to main.

- [ ] **Step 1: Run the full test suite**

```bash
cd /Users/h3n4l/OpenSource/omni/.worktrees/feat-partiql-parser-foundation
go test -v ./partiql/parser/... 2>&1 | tail -30
```

Expected tests passing:
- Existing lexer tests (Task 13 of lexer): 152 sub-tests
- TestParser_Machinery: 12 sub-tests
- TestParser_Goldens: ~42 sub-tests (10 literals + 15 exprprimary/paths + 6 math ops + 8 predicates + 3 logic)
- TestParseType: 46 sub-tests
- TestParser_AWSCorpus: ~61 sub-tests
- TestParser_Errors: 24 sub-tests

Total: ~337 sub-tests. Any failure: stop and investigate.

- [ ] **Step 2: Cross-check grammar coverage**

Open `/Users/h3n4l/OpenSource/parser/partiql/PartiQLParser.g4` and scan the `expr`, `exprOr`, `exprAnd`, `exprNot`, `exprPredicate`, `mathOp00/01/02`, `valueExpr`, `exprPrimary`, `exprTerm`, `pathStep`, `varRefExpr`, `literal`, `type`, `array`, `bag`, `tuple`, `pair` rules (lines ~445-686).

Verify that every ALTERNATIVE in each rule either:
- has a corresponding `parseXxx` function in the foundation code, OR
- is explicitly stubbed in `parsePrimaryBase` with a `deferredFeature` call

Expected coverage points:
- `exprBagOp` (UNION/INTERSECT/EXCEPT): stubbed
- `exprSelect` (SFW): stubbed
- `exprOr`: real
- `exprAnd`: real
- `exprNot`: real
- `exprPredicate` (comparison/IS/IN/LIKE/BETWEEN): real
- `mathOp00/01/02`: real
- `valueExpr`: real
- `exprPrimary` alternatives: 5 real (exprTerm, `exprPrimary pathStep+`), 11 stubbed (cast/aggregates/etc.)
- `exprTerm`: parameter/varRef/literal/collection/tuple all real; wrappedQuery (SubLink) falls through to SELECT stub
- `pathStep` (4 alternatives): all real
- `literal` (10 alternatives): 8 real, 2 stubbed (DATE, TIME)
- `type` (6 alternatives): all real via parseType

If any rule alternative is missing, add the stub and a test case.

- [ ] **Step 3: Run vet and gofmt**

```bash
go vet ./partiql/parser/...
gofmt -l partiql/parser/
grep -n $'\xe2\x80\x9d' partiql/parser/*.go
```

All three must produce no output (silent success).

- [ ] **Step 4: Final verification commit (allow-empty if no changes)**

```bash
git commit --allow-empty -m "$(cat <<'EOF'
chore(partiql/parser): final verification pass

- go test ./partiql/parser/... — pass
- go vet ./partiql/parser/... — clean
- gofmt -l partiql/parser/ — clean
- No U+201D gofmt corruption
- Cross-checked against bytebase/parser/partiql/PartiQLParser.g4:
  every expr/exprPrimary/literal/type alternative is either
  implemented or stubbed with an owner node.
- ~337 sub-tests passing:
  - TestParser_Machinery (12): token buffer helpers
  - TestParser_Goldens (~42): filesystem goldens
  - TestParseType (46): 30+ type forms
  - TestParser_AWSCorpus (~61): real corpus smoke test
  - TestParser_Errors (24): 16 stubs + 8 real syntax errors

Closes parser-foundation (DAG node 4) for the PartiQL migration.
Unblocks parser-select (node 5), parser-dml (node 6), parser-ddl
(node 7), parse-entry (node 8), parser-window (node 13),
parser-aggregates (node 14), parser-builtins (node 15),
parser-graph (node 16), parser-datetime-literals (node 18).

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

- [ ] **Step 5: Update `dag.md` on `main` to mark node 4 as `done`**

```bash
cd /Users/h3n4l/OpenSource/omni
git checkout main
```

Open `docs/migration/partiql/dag.md`. Find the row:

```
| 4 | parser-foundation | `partiql/parser` (parser.go, expr.go, path.go, literals.go) | ast-core, lexer | — | **P0** | not started |
```

Change the package column to reflect the 5-file split and the status to `done`:

```
| 4 | parser-foundation | `partiql/parser` (parser.go, expr.go, exprprimary.go, path.go, literals.go) | ast-core, lexer | — | **P0** | done |
```

Commit on main:

```bash
git add docs/migration/partiql/dag.md
git commit -m "$(cat <<'EOF'
docs(partiql): mark parser-foundation (DAG node 4) as done

The partiql/parser foundation is complete on the feat/partiql/
parser-foundation branch: Parser struct, token buffer helpers,
full precedence ladder (parseBagOp → parseSelectExpr → parseOr
→ parseAnd → parseNot → parsePredicate → parseMathOp00/01/02
→ parseValueExpr → parsePrimary), parsePrimary + 7 exprTerm
alternatives, path steps, literals, parseType for 30+ type
forms, and deferred-feature stubs for every non-foundation
expression form.

~337 sub-tests passing across TestParser_Machinery,
TestParser_Goldens (~42 filesystem goldens), TestParseType,
TestParser_AWSCorpus, and TestParser_Errors.

DAG nodes 5 (parser-select), 6 (parser-dml), 7 (parser-ddl),
8 (parse-entry), 13-16, and 18 are now unblocked.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

- [ ] **Step 6: Return to worktree and run `superpowers:finishing-a-development-branch`**

```bash
cd /Users/h3n4l/OpenSource/omni/.worktrees/feat-partiql-parser-foundation
git log --oneline feat/partiql/parser-foundation ^main
```

Verify the commit list (expect ~16 commits: 13 task commits + spec + plan + correction).

Then invoke `superpowers:finishing-a-development-branch`. It will guide the merge/PR decision.

- [ ] **Step 7: After branch finish, report to user**

1. Confirm `git worktree list` no longer shows `feat-partiql-parser-foundation` (if option 1 or 4 was chosen)
2. Confirm `dag.md` on `main` shows parser-foundation as `done`
3. Report to the user:
   - Which DAG node was completed (4 — parser-foundation)
   - Stats: ~1180 lines production, ~350 lines test, 43 golden pairs, ~337 sub-tests
   - Next actionable DAG nodes: 5 (parser-select), 6 (parser-dml), 7 (parser-ddl), possibly in parallel
