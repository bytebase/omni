# PartiQL Parser Foundation — Design Spec

**DAG node:** 4 (parser-foundation), P0
**Depends on:** ast-core (node 1, merged #13), lexer (node 2, merged #16)
**Unblocks:** parser-select (5), parser-dml (6), parser-ddl (7), parse-entry (8), parser-window (13), parser-aggregates (14), parser-builtins (15), parser-graph (16), parser-datetime-literals (18)
**Package:** `partiql/parser`
**Files added:** `parser.go`, `expr.go`, `exprprimary.go`, `path.go`, `literals.go`, `parser_test.go`, `testdata/parser-foundation/*.partiql|*.golden`
**Files modified:** none
**Grammar source:** `/Users/h3n4l/OpenSource/parser/partiql/PartiQLParser.g4`

---

## 1. Goal

Land the foundation layer of a hand-written recursive-descent parser for PartiQL. This PR ships every piece of parser machinery that future DAG nodes need to share: the `Parser` struct, token buffer, precedence ladder, primary-expression dispatch, path-step chains, literal parsing, and `parseType` (the type reference used by both CAST and DDL).

Every non-foundation feature is stubbed with a "deferred to DAG node N" error. Future nodes replace stub bodies without reshaping dispatch.

---

## 2. Architecture

### 2.1 File layout

```
partiql/parser/
├── parser.go       (~280 lines)  Parser struct, ParseError, token buffer helpers,
│                                 parseSymbolPrimitive, parseVarRef, parseType
├── expr.go         (~400 lines)  Precedence ladder: parseExpr → parseBagOp →
│                                 parseSelectExpr → parseOr → parseAnd → parseNot →
│                                 parsePredicate → parseMathOp00/01/02 → parseValueExpr
├── exprprimary.go  (~250 lines)  parsePrimary dispatch, exprTerm alternatives
│                                 (parenthesis, varRef, array, bag, tuple, values,
│                                 paramRef), deferred-feature stubs
├── path.go         (~140 lines)  parsePathSteps (.field, .*, [expr], [*] chain)
├── literals.go     (~170 lines)  parseLiteral switch for the 6 base literal forms
├── parser_test.go  (~350 lines)  TestParser_Goldens, TestParser_AWSCorpus,
│                                 TestParser_Errors, TestParser_Machinery, TestParseType
└── testdata/parser-foundation/
    ├── <45 .partiql input files>
    └── <45 matching .golden files>
```

File decomposition principle: each file holds one clear responsibility; `expr.go` (precedence ladder) and `exprprimary.go` (primary-expression dispatch) are split to keep each under 400 lines and focused on one concern.

### 2.2 Public API surface

Exports from `partiql/parser`:
- `type Parser struct` — the parser state machine
- `type ParseError struct` — syntax error with Message and `ast.Loc` position
- `func NewParser(input string) *Parser` — constructor wrapping a Lexer
- `func (p *Parser) ParseExpr() (ast.ExprNode, error)` — the foundation test entry point

**Not exported by this PR:**
- `Parse(sql string) (*ast.List, error)` — public entry point belongs to node 8 (parse-entry)
- `ParseStatement`, `ParseScript` — statement-producing entry points belong to nodes 5/6/7/8

### 2.3 Invariants inherited from upstream

- **Token struct**: `{Type int, Str string, Loc ast.Loc}` from `partiql/parser/token.go`
- **Lexer contract**: first-error-and-stop via `Lexer.Err`; after Err is set, all `Next()` calls return `tokEOF`
- **ast.Loc contract**: `{Start, End int}` half-open byte range; `{-1, -1}` means synthetic/unknown
- **AST sealed interfaces**: Node / StmtNode / ExprNode / TableExpr / PathStep / TypeName / PatternNode. This PR constructs `ExprNode`, `PathStep`, and `TypeName` values; no `StmtNode` construction.

---

## 3. Parser Machinery (`parser.go`)

### 3.1 Parser struct

```go
type Parser struct {
    lexer   *Lexer
    cur     Token   // current token (peek position)
    prev    Token   // most recently consumed token, for end-loc computation
    nextBuf Token   // one-token lookahead buffer for peekNext
    hasNext bool
}

type ParseError struct {
    Message string
    Loc     ast.Loc // points at the offending token (Start = End = token start)
}

func (e *ParseError) Error() string {
    return fmt.Sprintf("syntax error at position %d: %s", e.Loc.Start, e.Message)
}

func NewParser(input string) *Parser {
    p := &Parser{lexer: NewLexer(input)}
    p.advance() // prime cur with the first token
    return p
}
```

### 3.2 Token buffer helpers

| Helper | Purpose |
|--------|---------|
| `advance()` | Moves cur → prev, reads next token from lexer, propagates lexer errors as tokEOF-with-embedded-Err |
| `peek() Token` | Returns cur without consuming |
| `peekNext() Token` | Returns the token after cur without consuming (uses nextBuf slot) |
| `match(types ...int) bool` | Consumes cur if its Type is in the set; idiomatic for optional keywords |
| `expect(tokenType int) (Token, error)` | Consumes cur if its Type matches, else returns a `*ParseError` |

**Lexer error propagation**: Parser calls `advance()`, then checks `p.lexer.Err` at strategic points (function entry, after `expect`). When Err is set, it's surfaced as a `*ParseError` wrapping the lexer message. Caller code only handles one error type (`*ParseError`).

### 3.3 Error model

**Fail-fast.** The first syntax error aborts the parse and returns one `*ParseError`. Matches `cosmosdb/parser`. Error recovery is NOT in scope for this node — if DAG node 10 (completion) needs partial-parse recovery later, that layer will be added informed by completion's actual needs.

**Error format** (consistent across all error sites):
```go
&ParseError{
    Message: fmt.Sprintf("expected %s, got %q", tokenName(want), p.cur.Str),
    Loc:     ast.Loc{Start: p.cur.Loc.Start, End: p.cur.Loc.Start},
}
```

Position is the START of the offending token — readable error, points where the user should look.

### 3.4 Identifier and varref helpers

```go
// parseSymbolPrimitive consumes an IDENTIFIER or IDENTIFIER_QUOTED, returning
// the name and CaseSensitive flag (true for quoted).
// Grammar: symbolPrimitive (line 742)
func (p *Parser) parseSymbolPrimitive() (name string, caseSensitive bool, loc ast.Loc, err error)

// parseVarRef handles @-prefix + symbolPrimitive. Also detects IDENT followed
// by PAREN_LEFT (function call) and returns a deferred-feature stub error.
// Grammar: varRefExpr (line 635-636)
func (p *Parser) parseVarRef() (ast.ExprNode, error)
```

### 3.5 parseType (shared utility)

```go
// parseType consumes a PartiQL type reference. Used by CAST (deferred to
// parser-builtins, but parseType ships now so both parser-builtins and
// parser-ddl can import it without coupling).
// Grammar: type (lines 674-686)
func (p *Parser) parseType() (*ast.TypeRef, error)
```

Handles:
- Atomic types: NULL, BOOL, BOOLEAN, SMALLINT, INT/INT2/INT4/INT8, INTEGER/INTEGER2/INTEGER4/INTEGER8, BIGINT, REAL, TIMESTAMP, CHAR, CHARACTER, MISSING, STRING, SYMBOL, BLOB, CLOB, DATE, STRUCT, TUPLE, LIST, SEXP, BAG, ANY
- `DOUBLE PRECISION` (two-token form)
- Parameterized single-arg: `CHAR(n) / CHARACTER(n) / FLOAT(p) / VARCHAR(n)`
- `CHARACTER VARYING [ (n) ]`
- Parameterized two-arg: `DECIMAL(p,s) / DEC(p,s) / NUMERIC(p,s)`
- `TIME [(p)] [WITH TIME ZONE]`
- Custom: any `symbolPrimitive` fallback

---

## 4. Expression Layer (`expr.go`)

### 4.1 Precedence ladder — one function per grammar layer

The PartiQL grammar's expression precedence is explicitly layered across 8 rules (lines 445-512). Foundation mirrors each rule as a dedicated function, not a Pratt-style precedence table. Rationale:

1. Trivial cross-check against the spec — each function cites its grammar rule by name and line number
2. No precedence table to maintain — the order is encoded in the call chain
3. Left-recursive grammar rules translate cleanly to iterative loops
4. Matches how future maintainers read the grammar file

```go
func (p *Parser) parseExpr() (ast.ExprNode, error)        // → parseBagOp
func (p *Parser) parseBagOp() (ast.ExprNode, error)       // UNION/INTERSECT/EXCEPT [STUB]
func (p *Parser) parseSelectExpr() (ast.ExprNode, error)  // SfwQuery shape [STUB]
func (p *Parser) parseOr() (ast.ExprNode, error)          // OR (left-assoc)
func (p *Parser) parseAnd() (ast.ExprNode, error)         // AND (left-assoc)
func (p *Parser) parseNot() (ast.ExprNode, error)         // NOT (right-assoc prefix)
func (p *Parser) parsePredicate() (ast.ExprNode, error)   // comparison/IS/IN/LIKE/BETWEEN
func (p *Parser) parseMathOp00() (ast.ExprNode, error)    // || concat
func (p *Parser) parseMathOp01() (ast.ExprNode, error)    // + -
func (p *Parser) parseMathOp02() (ast.ExprNode, error)    // * / %
func (p *Parser) parseValueExpr() (ast.ExprNode, error)   // unary +/- prefix
```

### 4.2 Iterative left-recursion pattern

Grammar: `mathOp01 : lhs=mathOp01 op=(PLUS|MINUS) rhs=mathOp02 | mathOp02`

Implementation:
```go
func (p *Parser) parseMathOp01() (ast.ExprNode, error) {
    left, err := p.parseMathOp02()
    if err != nil {
        return nil, err
    }
    for p.cur.Type == tokPLUS || p.cur.Type == tokMINUS {
        op := p.cur.Str
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
```

Same shape for `parseMathOp00` (concat) and `parseMathOp02` (mult).

### 4.3 parsePredicate

Handles 5 grammar alternatives: comparison (`< <= = != >= >`), `IS [NOT] type`, `[NOT] IN ...`, `[NOT] LIKE ... [ESCAPE ...]`, `[NOT] BETWEEN lower AND upper`.

```go
func (p *Parser) parsePredicate() (ast.ExprNode, error) {
    left, err := p.parseMathOp00()
    if err != nil {
        return nil, err
    }
    for {
        switch {
        case isComparisonOp(p.cur.Type):
            // build BinaryExpr
        case p.cur.Type == tokIS:
            // IS [NOT] type  →  IsExpr
        case p.cur.Type == tokNOT:
            next := p.peekNext().Type
            if next == tokIN || next == tokLIKE || next == tokBETWEEN {
                p.advance() // consume NOT
                // dispatch with Negated=true
            } else {
                return left, nil
            }
        case p.cur.Type == tokIN:
            // parseIn(left, false, startLoc)
        case p.cur.Type == tokLIKE:
            // parseLike(left, false, startLoc)
        case p.cur.Type == tokBETWEEN:
            // parseBetween(left, false, startLoc)
        default:
            return left, nil
        }
    }
}
```

Private helpers:
- `parseIn(left ExprNode, not bool, startLoc int) (*ast.InExpr, error)` — handles both `IN (...)` parenthesized form and `IN rhs=mathOp00`
- `parseBetween(left ExprNode, not bool, startLoc int) (*ast.BetweenExpr, error)`
- `parseLike(left ExprNode, not bool, startLoc int) (*ast.LikeExpr, error)` — optional `ESCAPE escape=expr`
- `isComparisonOp(t int) bool`

### 4.4 Quirky precedence: `||` is LOOSER than `+`

The PartiQL grammar at lines 494-497 defines string concatenation (`||`, `mathOp00`) as binding LOOSER than addition (`+/-`, `mathOp01`). This is opposite to SQL92 and most other dialects. Foundation mirrors the grammar literally — `parseMathOp00` (concat) calls `parseMathOp01` (add) calls `parseMathOp02` (mult).

A dedicated golden test case (`op_precedence_concat.partiql`) locks this in:
```
input:  'a' || 1 + 2
AST:    BinaryExpr{Op: '||', Left: 'a', Right: BinaryExpr{Op: '+', Left: 1, Right: 2}}
```

---

## 5. Primary Expressions and Path Steps

### 5.1 exprprimary.go — parsePrimary dispatch

```go
func (p *Parser) parsePrimary() (ast.ExprNode, error) {
    base, err := p.parsePrimaryBase()
    if err != nil {
        return nil, err
    }
    if isPathStepStart(p.cur.Type) {
        return p.parsePathSteps(base)
    }
    return base, nil
}
```

`parsePrimaryBase` dispatches on the current token across 16 grammar alternatives. Real cases (foundation-owned):

- Literals: TRUE, FALSE, NULL, MISSING, SCONST, ICONST, FCONST, ION_LITERAL, DATE, TIME → `parseLiteral()` (date/time are stubbed inside parseLiteral itself)
- `tokPAREN_LEFT` → `parseParenExpr()` — plain parenthesized expression; `(SELECT …)` subquery routes to the node-5 stub; `(expr, …)` valueList routes to the node-6 stub (no AST support yet)
- `tokBRACKET_LEFT` → `parseArrayLit()`
- `tokANGLE_DOUBLE_LEFT` → `parseBagLit()`
- `tokBRACE_LEFT` → `parseTupleLit()`
- `tokQUESTION_MARK` → `parseParamRef()`
- `tokAT_SIGN`, `tokIDENT`, `tokIDENT_QUOTED` → `parseVarRef()`
- `tokVALUES` → STUB: deferred to parser-dml (DAG node 6) — no AST node for VALUES row lists exists yet in ast-core

All other primary-expression alternatives are stubbed with `p.deferredFeature(...)` errors pointing at the owning DAG node.

### 5.2 exprTerm alternatives (foundation-owned)

```go
// parseParenExpr disambiguates parenthesis-wrapped forms:
//   (expr)            → parenthesized expression (returns the inner expr)
//   (SELECT ...)      → STUB: deferred to parser-select (DAG node 5)
//   (expr, expr, ...) → STUB: deferred to parser-dml (DAG node 6) — valueList has
//                       no AST node yet
//   (expr MATCH ...)  → STUB: deferred to parser-graph (DAG node 16)
func (p *Parser) parseParenExpr() (ast.ExprNode, error)

func (p *Parser) parseArrayLit() (*ast.ListLit, error)   // [expr, ...]
func (p *Parser) parseBagLit() (*ast.BagLit, error)      // <<expr, ...>>
func (p *Parser) parseTupleLit() (*ast.TupleLit, error)  // {key: value, ...}
func (p *Parser) parseParamRef() (*ast.ParamRef, error)  // ?
```

### 5.3 path.go — path step chain

```go
func (p *Parser) parsePathSteps(base ast.ExprNode) (*ast.PathExpr, error) {
    steps := []ast.PathStep{}
    for isPathStepStart(p.cur.Type) {
        step, err := p.parsePathStep()
        if err != nil {
            return nil, err
        }
        steps = append(steps, step)
    }
    return &ast.PathExpr{
        Root:  base,
        Steps: steps,
        Loc:   ast.Loc{Start: base.GetLoc().Start, End: steps[len(steps)-1].GetLoc().End},
    }, nil
}

func (p *Parser) parsePathStep() (ast.PathStep, error)

func isPathStepStart(t int) bool { return t == tokPERIOD || t == tokBRACKET_LEFT }
```

Four path-step flavors per the grammar (lines 618-623):
| Source | AST node | Notes |
|--------|----------|-------|
| `.field` | `ast.DotStep{Field, CaseSensitive: false}` | bare identifier |
| `."Field"` | `ast.DotStep{Field, CaseSensitive: true}` | quoted identifier |
| `.*` | `ast.AllFieldsStep{}` | all-fields wildcard |
| `[expr]` | `ast.IndexStep{Index}` | key/index by expression |
| `[*]` | `ast.WildcardStep{}` | all-elements wildcard |

The `[expr]` vs `[*]` disambiguation: after consuming `BRACKET_LEFT`, peek at the next token. If `ASTERISK` immediately followed by `BRACKET_RIGHT`, emit `WildcardStep`; otherwise parse a full expression for `IndexStep`.

---

## 6. Literals (`literals.go`)

```go
func (p *Parser) parseLiteral() (ast.ExprNode, error) {
    switch p.cur.Type {
    case tokNULL:
        // NullLit{Loc: p.cur.Loc}
    case tokMISSING:
        // MissingLit{Loc: p.cur.Loc}
    case tokTRUE:
        // BoolLit{Val: true, Loc: p.cur.Loc}
    case tokFALSE:
        // BoolLit{Val: false, Loc: p.cur.Loc}
    case tokSCONST:
        // StringLit{Val: p.cur.Str, Loc: p.cur.Loc}
    case tokICONST, tokFCONST:
        // NumberLit{Val: p.cur.Str, Loc: p.cur.Loc}
    case tokION_LITERAL:
        // IonLit{Text: p.cur.Str, Loc: p.cur.Loc}
    case tokDATE:
        // STUB: deferred to parser-datetime-literals (DAG node 18)
    case tokTIME:
        // STUB: deferred to parser-datetime-literals (DAG node 18)
    }
}
```

Token.Str arrives pre-decoded from the lexer (scanString already collapsed doubled-quote escapes in Task 6 of the lexer plan). NumberLit.Val preserves the raw source text to distinguish integer / decimal / scientific forms; consumers needing a typed value call `strconv.ParseFloat` or `shopspring/decimal` at the call site. This matches the ast-core spec's `NumberLit` contract.

---

## 7. Deferred-Feature Stubs

### 7.1 Stub error format

All stubs use a uniform error format produced by a shared helper:

```go
func (p *Parser) deferredFeature(feature, ownerNode string) error {
    return &ParseError{
        Message: fmt.Sprintf("%s is deferred to %s", feature, ownerNode),
        Loc:     ast.Loc{Start: p.cur.Loc.Start, End: p.cur.Loc.End},
    }
}
```

### 7.2 Complete stub audit

**Node 5 (parser-select):**
- `parseBagOp` — UNION/INTERSECT/EXCEPT dispatch (expr.go)
- `parseSelectExpr` — SFW query (expr.go)
- `SubLink` case in `parseParenExpr` — `(SELECT ...)` subquery

**Node 6 (parser-dml):**
- `tokVALUES` case in parsePrimaryBase — no AST node for VALUES row lists yet
- `(expr, expr, ...)` valueList case in `parseParenExpr` — same rationale
- `tokINSERT` case in `parseSelectExpr` — INSERT statement top-level dispatch
- `tokUPDATE` case in `parseSelectExpr` — UPDATE statement top-level dispatch
- `tokDELETE` case in `parseSelectExpr` — DELETE statement top-level dispatch

**Node 13 (parser-window):**
- `tokLAG`, `tokLEAD` in parsePrimaryBase

**Node 14 (parser-aggregates):**
- `tokCOUNT`, `tokMAX`, `tokMIN`, `tokSUM`, `tokAVG` in parsePrimaryBase

**Node 15 (parser-builtins):**
- `tokCAST`, `tokCAN_CAST`, `tokCAN_LOSSLESS_CAST`
- `tokCASE`
- `tokCOALESCE`, `tokNULLIF`
- `tokSUBSTRING`, `tokTRIM`, `tokEXTRACT`
- `tokDATE_ADD`, `tokDATE_DIFF`
- `tokCHAR_LENGTH`, `tokCHARACTER_LENGTH`, `tokOCTET_LENGTH`, `tokBIT_LENGTH`
- `tokUPPER`, `tokLOWER`, `tokSIZE`, `tokEXISTS`
- `tokLIST`, `tokSEXP` (sequenceConstructor)
- Generic `functionCall` — triggered by `IDENT PAREN_LEFT` lookahead in parseVarRef

**Node 16 (parser-graph):**
- `exprGraphMatchMany` — detected via `MATCH` after an expression in `parseParenExpr`

**Node 18 (parser-datetime-literals):**
- `tokDATE` literal body in parseLiteral
- `tokTIME` literal body in parseLiteral

### 7.3 Stub-grep contract

Future DAG node implementers can find their work by running:
```
grep -rn "deferred to parser-builtins" partiql/parser/
grep -rn "deferred to parser-select" partiql/parser/
```
and get an exact list of call sites they need to replace. Each stub's doc comment cites the grammar line number to speed the swap.

### 7.4 Example stub-produced errors

| Input | Error |
|-------|-------|
| `SELECT * FROM t` | `SELECT is deferred to parser-select (DAG node 5)` |
| `a UNION b` | `UNION is deferred to parser-select (DAG node 5)` |
| `CAST(x AS INT)` | `CAST is deferred to parser-builtins (DAG node 15)` |
| `COUNT(*)` | `COUNT() aggregate is deferred to parser-aggregates (DAG node 14)` |
| `LAG(x) OVER (...)` | `LAG() window is deferred to parser-window (DAG node 13)` |
| `SUBSTRING(x, 1, 2)` | `SUBSTRING is deferred to parser-builtins (DAG node 15)` |
| `CASE WHEN a THEN 1 END` | `CASE is deferred to parser-builtins (DAG node 15)` |
| `DATE '2026-01-01'` | `DATE literal is deferred to parser-datetime-literals (DAG node 18)` |
| `foo(x)` | `function call "foo" is deferred to parser-builtins (DAG node 15)` |
| `LIST(1,2,3)` | `LIST() constructor is deferred to parser-builtins (DAG node 15)` |
| `VALUES (1,2)` | `VALUES is deferred to parser-dml (DAG node 6)` |
| `(1, 2, 3)` | `valueList is deferred to parser-dml (DAG node 6)` |

---

## 8. Testing

### 8.1 TestParser_Goldens — filesystem goldens

Matches `cosmosdb/parser/golden_test.go`. Pairs `.partiql` input files with `.golden` pretty-printed AST output files.

```go
var update = flag.Bool("update", false, "update golden files")

func TestParser_Goldens(t *testing.T) {
    files, _ := filepath.Glob("testdata/parser-foundation/*.partiql")
    for _, inPath := range files {
        name := strings.TrimSuffix(filepath.Base(inPath), ".partiql")
        t.Run(name, func(t *testing.T) {
            input, _ := os.ReadFile(inPath)
            p := NewParser(string(input))
            expr, err := p.ParseExpr()
            if err != nil {
                t.Fatalf("parse error: %v", err)
            }
            got := ast.NodeToString(expr)
            goldenPath := strings.TrimSuffix(inPath, ".partiql") + ".golden"
            if *update {
                os.WriteFile(goldenPath, []byte(got), 0644)
                return
            }
            want, _ := os.ReadFile(goldenPath)
            if got != string(want) {
                t.Errorf("mismatch\ngot:\n%s\nwant:\n%s", got, string(want))
            }
        })
    }
}
```

### 8.2 Golden corpus (~43 inputs)

Committed under `testdata/parser-foundation/`:

**Literals (10 cases):** null, missing, true, false, int, decimal, scientific, string, string with doubled quote, ion

**VarRef and Param (4 cases):** bare, quoted, at-prefixed, parameter `?`

**Collections (4 cases):** array, empty array, bag, tuple

**Path expressions (7 cases):** dot, dot-quoted, dot-star, index-expr, index-wildcard, multi-step chain, path on parenthesized base

**Arithmetic and precedence (6 cases):** add, mul, `1 + 2 * 3`, quirky `'a' || 1 + 2`, unary neg, unary NOT

**Comparison and logic (3 cases):** `a < b`, `a = 1 AND b = 2`, `a OR b AND c` (precedence)

**Predicates (8 cases):** IS NULL, IS NOT NULL, IN list, NOT IN, LIKE, LIKE ESCAPE, BETWEEN, NOT BETWEEN

**Parentheses (1 case):** `(1+2)` plain parenthesized expression (VALUES and valueList are stubbed, not tested as goldens)

**Stress (1 case):** `a.b.c + (d.e[0] * f[*]) - g`

### 8.3 TestParser_AWSCorpus — smoke test

```go
func TestParser_AWSCorpus(t *testing.T) {
    skip := map[string]bool{
        "select-001.partiql": true,
        "insert-002.partiql": true,
    }
    files, _ := filepath.Glob("testdata/aws-corpus/*.partiql")
    var ok, stubbed int
    for _, f := range files {
        name := filepath.Base(f)
        if skip[name] {
            continue
        }
        t.Run(name, func(t *testing.T) {
            data, _ := os.ReadFile(f)
            p := NewParser(string(data))
            _, err := p.ParseExpr()
            if err == nil {
                ok++
                return
            }
            if !strings.Contains(err.Error(), "deferred to") {
                t.Errorf("unexpected parse error (not a deferred-feature stub): %v", err)
                return
            }
            stubbed++
        })
    }
    t.Logf("AWS corpus: %d fully parsed, %d hit deferred-feature stubs", ok, stubbed)
}
```

At this milestone, most corpus files start with `SELECT` and will hit the node-5 stub — that's expected. The test only fails on panics or non-deferred errors.

### 8.4 TestParser_Errors — ~18 cases

Covers deferred-feature stubs (locks in their exact messages for the grep contract) and real syntax errors (unclosed paren, missing rhs, unclosed array/bag/tuple, missing colon in pair, BETWEEN missing AND, LIKE missing pattern).

### 8.5 TestParser_Machinery — ~6 cases

Unit tests for match/expect/peek/peekNext/lexer error propagation.

### 8.6 TestParseType — ~20 cases

Table-driven tests for every type form in `parseType`: atomic types, parameterized types, DOUBLE PRECISION, CHARACTER VARYING, TIME WITH TIME ZONE, custom types.

### 8.7 What's deliberately NOT tested

- **No tokenName-style parity test.** The parser has no constant table that can drift; the golden tests are the de facto parity mechanism.
- **No error recovery tests.** Foundation is fail-fast by design (Q2 decision).
- **No statement-level parsing.** Nodes 5/6/7/8 own that.

---

## 9. Implementation Order (13 tasks)

1. **parser.go scaffold + TestParser_Machinery** — Parser struct, ParseError, constructor, token buffer helpers, parseSymbolPrimitive, parseVarRef (bare, no function call detection yet). Tests: TestParser_Machinery (~6 cases).
2. **literals.go + TestParser_Goldens runner + 10 literal goldens** — parseLiteral switch for the 6 base forms + date/time stubs. Sets up the golden harness with `-update` flag so all subsequent tasks can use it.
3. **parser.go parseType + TestParseType** — all 30+ type forms. Table-driven tests (not goldens; TypeRef is a leaf node and table-driven is more explicit).
4. **path.go + 7 path goldens** — parsePathSteps, parsePathStep, isPathStepStart. Does NOT wire into parsePrimary yet (no parsePrimary exists at this point).
5. **exprprimary.go parsePrimaryBase + exprTerm cases + ~13 goldens** — paren (plain expr), paramRef, varRef (bare), array, bag, tuple. All deferred-feature stubs present (including VALUES/valueList/SubLink/graph-match). Wires parsePathSteps into parsePrimary.
6. **expr.go parseMathOp02/01/00 and parseValueExpr + 6 precedence goldens** — mult, add, concat, unary sign. Includes the quirky `||` precedence case.
7. **expr.go parsePredicate + 8 predicate goldens** — comparison, IS, IN, LIKE, BETWEEN.
8. **expr.go parseNot, parseAnd, parseOr + 3 logic goldens** — logical layers.
9. **expr.go parseBagOp and parseSelectExpr stubs + error tests** — top-of-ladder dispatch.
10. **parseVarRef upgrade for function call detection + graph match stub** — `IDENT PAREN_LEFT` and `MATCH` detection.
11. **TestParser_AWSCorpus** — corpus smoke test against the 61 AWS corpus files.
12. **TestParser_Errors** — consolidated error-case table covering deferred-feature stubs and real syntax errors.
13. **Final verification + DAG bookkeeping + finishing branch** — full test run, grammar cross-check, gofmt/vet, mark node 4 done, invoke finishing-a-development-branch.

Estimated delta: ~1180 lines production + ~350 lines test + ~43 golden pairs.

---

## 10. Design Decisions Log

**D1. Scope boundary — stub exprBagOp/exprSelect instead of stopping at exprOr.** Foundation wires the full precedence ladder dispatch but the SELECT-shaped layers return deferred-feature errors. Node 5 replaces stub bodies without reshaping dispatch. Rationale: cleaner ownership, future-proofs the ladder.

**D2. Error model — fail-fast, no recovery.** Matches `cosmosdb/parser`. Error recovery (if needed by node 10 completion) will be designed then, informed by real completion requirements. Avoids speculative complexity.

**D3. File decomposition — 5 files, not 4.** Added `exprprimary.go` to hold primary-expression dispatch and exprTerm alternatives, separating the precedence ladder (expr.go) from the primary-expression dispatch. Keeps each file under ~400 lines and focused on one concern.

**D4. Test strategy — filesystem goldens + corpus smoke + errors.** Matches `cosmosdb/parser/golden_test.go` pattern. Leverages `ast.NodeToString` from ast-core. Regeneration via `go test -update` flag. Avoids hand-written Go struct literal goldens (verbose, painful to maintain across AST shape changes).

**D5. parseType in foundation; valueList/values deferred to parser-dml.** `parseType` is ~80 lines of grammar-driven switch; keeping it in foundation decouples node 7 (DDL) and node 15 (builtins) — either can ship without the other. `valueList`/`values` were originally slated for foundation during brainstorming, but plan-time review found no AST representation for them (no `ValuesExpr`, no `ValueList`, no `RowExpr` in `partiql/ast`), so they are stubbed to parser-dml (DAG node 6) which owns VALUES as an INSERT source. Adding the AST nodes is out of scope for parser-foundation.

**D6. Precedence ladder — one function per grammar layer, not a Pratt table.** PartiQL's grammar is explicitly layered (lines 445-512); mirroring it 1:1 makes the code easier to cross-check against the spec and matches how maintainers read the grammar file. No hidden precedence table to drift.
