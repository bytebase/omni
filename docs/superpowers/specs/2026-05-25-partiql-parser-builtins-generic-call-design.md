# PartiQL parser-builtins-generic-call (15a) — Design Spec

**DAG node:** 15a (parser-builtins-generic-call), **P0** — post-cutover production regression fix
**Depends on:** parser-foundation (node 4, merged)
**Unblocks:** parser-builtins-typed (15b), and unblocks every DynamoDB-native predicate (`attribute_exists`, `attribute_not_exists`, `attribute_type`, `begins_with`, `contains`, plus any user-defined `IDENT(...)` call) in the bytebase SQL editor
**Package:** `partiql/parser`
**Files modified:** `parser.go`, `parser_test.go`
**Files added:** none
**Grammar source:** `/Users/h3n4l/OpenSource/parser/partiql/PartiQLParser.g4` rule `functionCall#FunctionCallIdent` (lines 611-616)

---

## 1. Goal

After the bytebase cutover (omni PR series #19965–#19968), omni's parser gates every DynamoDB statement through the editor's syntax check. The legacy ANTLR grammar accepted any `IDENT(...)` call — including the DynamoDB-native predicates customers use constantly in `WHERE` clauses — but omni today rejects them with:

```
function call "X" is deferred to parser-builtins (DAG node 15)
```

The editor surfaces that as a syntax error. **Real customer queries like `SELECT * FROM Orders WHERE contains(Address, 'Kirkland')` cannot run today.**

Node 15a stops the bleeding by parsing generic `IDENT(args)` calls into the existing `*ast.FuncCall` node. Typed keyword builtins (CAST, CASE, SUBSTRING, ...) remain stubbed and ship under node 15b.

---

## 2. Architecture

### 2.1 What changes

One function in `partiql/parser/parser.go` (`parseVarRef`) and a new small helper (`parseFuncCallArgs`).

```
partiql/parser/
├── parser.go          MODIFY parseVarRef; ADD parseFuncCallArgs helper
└── parser_test.go     MODIFY funcall_stub case; ADD positive + negative cases
```

The AWS corpus index (`testdata/aws-corpus/index.json`) needs no edit: `TestParser_AWSCorpus` (parser_test.go:376) classifies each corpus entry as "stubbed" vs "fully parsed" at runtime by inspecting the parser output for a `"deferred to"` substring. Once `parseVarRef` returns `*ast.FuncCall` instead of the deferred error, the three function-* entries automatically reclassify from stubbed to fully-parsed.

### 2.2 Parser change — `parseVarRef`

**Current behavior** (parser.go:239-265):

```go
func (p *Parser) parseVarRef() (ast.ExprNode, error) {
    start := p.cur.Loc.Start
    atPrefixed := false
    if p.cur.Type == tokAT_SIGN { atPrefixed = true; p.advance() }
    name, caseSensitive, nameLoc, err := p.parseSymbolPrimitive()
    if err != nil { return nil, err }

    // Function call lookahead: <name> ( ... )
    if p.cur.Type == tokPAREN_LEFT {
        return nil, &ParseError{
            Message: fmt.Sprintf("function call %q is deferred to parser-builtins (DAG node 15)", name),
            Loc:     ast.Loc{Start: start, End: p.cur.Loc.End},
        }
    }
    return &ast.VarRef{...}, nil
}
```

**After 15a:**

```go
func (p *Parser) parseVarRef() (ast.ExprNode, error) {
    start := p.cur.Loc.Start
    atPrefixed := false
    if p.cur.Type == tokAT_SIGN { atPrefixed = true; p.advance() }
    name, caseSensitive, nameLoc, err := p.parseSymbolPrimitive()
    if err != nil { return nil, err }

    // Function call: name=symbolPrimitive PAREN_LEFT (expr (COMMA expr)*)? PAREN_RIGHT
    // (PartiQLParser.g4:615 — FunctionCallIdent)
    if p.cur.Type == tokPAREN_LEFT {
        if atPrefixed {
            // ANTLR grammar's functionCall rule does NOT permit @-prefix
            // before the name; only varRefExpr does. Reject for parity.
            return nil, &ParseError{
                Message: "@-prefix is not allowed before a function call",
                Loc:     ast.Loc{Start: start, End: p.cur.Loc.End},
            }
        }
        args, endLoc, err := p.parseFuncCallArgs()
        if err != nil { return nil, err }
        return &ast.FuncCall{
            Name: name,
            Args: args,
            Loc:  ast.Loc{Start: start, End: endLoc},
        }, nil
    }
    return &ast.VarRef{...}, nil
}
```

### 2.3 New helper — `parseFuncCallArgs`

Extracted because **15b reuses it** for the ordinary-comma-list typed builtins (`DATE_ADD`, `DATE_DIFF`, `COALESCE`, `NULLIF`, `SIZE`, `EXISTS` — per the docstring on `*ast.FuncCall`).

```go
// parseFuncCallArgs consumes the parenthesized argument list of a function
// call: PAREN_LEFT (expr (COMMA expr)*)? PAREN_RIGHT. The opening paren must
// be the current token. Returns the args and the End offset of PAREN_RIGHT.
func (p *Parser) parseFuncCallArgs() ([]ast.ExprNode, int, error) {
    if _, err := p.expect(tokPAREN_LEFT); err != nil {
        return nil, 0, err
    }
    var args []ast.ExprNode
    if p.cur.Type != tokPAREN_RIGHT {
        for {
            arg, err := p.parseExprTop()
            if err != nil { return nil, 0, err }
            args = append(args, arg)
            if p.cur.Type != tokCOMMA { break }
            p.advance() // consume COMMA
        }
    }
    endTok := p.cur
    if _, err := p.expect(tokPAREN_RIGHT); err != nil {
        return nil, 0, err
    }
    return args, endTok.Loc.End, nil
}
```

Trailing comma is rejected because after consuming `,` the loop loops and `parseExprTop` errors on whatever follows (typically `)`).

### 2.4 AST — no change

`*ast.FuncCall` at `partiql/ast/exprs.go:336` already has the right shape:

```go
type FuncCall struct {
    Name       string
    Args       []ExprNode
    Quantifier QuantifierKind // zero value: QuantifierNone — set later by node 14
    Star       bool           // zero value: false       — set later by node 14
    Over       *WindowSpec    // zero value: nil         — set later by node 13
    Loc        Loc
}
```

15a populates `Name`, `Args`, `Loc`. The other fields stay at zero values.

### 2.5 What stays out of scope (deferred to 15b)

The token-type-dispatched cases in `exprprimary.go:60-158` are **untouched** by this change because `parseVarRef` only runs for `tokIDENT` / `tokIDENT_QUOTED`. The following token types continue to dispatch to their existing deferred-feature stubs:

- `tokCAST`, `tokCAN_CAST`, `tokCAN_LOSSLESS_CAST` (CAST family — non-comma syntax with `AS`)
- `tokCASE` (`CASE … WHEN … THEN …`)
- `tokCOALESCE`, `tokNULLIF` (defined as keyword tokens; will be reusable through `parseFuncCallArgs` in 15b)
- `tokSUBSTRING` (`SUBSTRING(s FROM start FOR len)` — non-comma syntax)
- `tokTRIM` (`TRIM(LEADING 'x' FROM s)` — non-comma syntax)
- `tokEXTRACT` (`EXTRACT(year FROM dt)` — non-comma syntax)
- `tokDATE_ADD`, `tokDATE_DIFF` (comma syntax — reuse `parseFuncCallArgs` later)
- `tokCHAR_LENGTH`, `tokCHARACTER_LENGTH`, `tokOCTET_LENGTH`, `tokBIT_LENGTH`, `tokUPPER`, `tokLOWER`, `tokSIZE`, `tokEXISTS` (reserved-name scalars — comma syntax, reuse helper later)
- `tokLIST`, `tokSEXP` (sequence constructors)
- `tokCOUNT`, `tokMAX`, `tokMIN`, `tokSUM`, `tokAVG` (aggregates — node 14)
- `tokLAG`, `tokLEAD` (window — node 13)

These stay deferred. 15b/13/14 will replace them one at a time.

### 2.6 Path-step interaction

`parsePrimary` (exprprimary.go:16) calls `parsePrimaryBase` (which calls `parseVarRef`), then `parsePathSteps` (path.go:29). After 15a, `parseVarRef` can return either `*ast.VarRef` or `*ast.FuncCall`; both implement `exprNode()`. The downstream `parsePathSteps` wraps any `ExprNode` in a `*ast.PathExpr` when followed by `.field`/`[idx]`/etc. So `foo(x).bar` yields `*ast.PathExpr{Root: *ast.FuncCall{...}, Steps: [...]}` without any other code change.

---

## 3. Test Plan

### 3.1 Unit tests — `parser_test.go`

**Flip the existing `funcall_stub` case** (line 670). It currently asserts the deferred-feature error; replace with a positive assertion that `foo(x)` produces `*ast.FuncCall{Name: "foo", Args: [*ast.VarRef{Name: "x"}]}`.

**Add positive cases:**

| Input | Expected shape |
|---|---|
| `foo()` | `*FuncCall{Name: "foo", Args: nil}` (Go zero value — `parseFuncCallArgs` returns `nil` when no args, never an empty slice) |
| `foo(1, 2, 3)` | `*FuncCall{Name: "foo", Args: 3× *Literal}` |
| `f(g(x))` | `*FuncCall{Name: "f", Args: [*FuncCall{Name: "g", Args: [*VarRef]}]}` |
| `f(x).bar` | `*PathExpr{Root: *FuncCall{Name: "f"}, Steps: [.bar]}` |
| `attribute_exists("a")` | `*FuncCall{Name: "attribute_exists", Args: [*VarRef{Name: "a", CaseSensitive: true}]}` |
| `begins_with(addr, '7834')` | `*FuncCall{Name: "begins_with", Args: [*VarRef, *Literal]}` |

**Add negative cases:**

| Input | Expected error contains |
|---|---|
| `@foo(x)` | `"@-prefix is not allowed before a function call"` |
| `foo(x,)` | parse error on `)` (whatever `parseExprTop` produces for empty input) |
| `foo(x` | `"expected PAREN_RIGHT"` |
| `foo(,)` | parse error on `,` |

### 3.2 Corpus — `testdata/aws-corpus/`

No JSON change. `TestParser_AWSCorpus` (parser_test.go:376) automatically reclassifies corpus entries when their parse outcome shifts. Three entries that currently log as stubbed will silently flip to fully-parsed once `parseVarRef` produces `*ast.FuncCall`:

- `function-begins-with-001` — `SELECT * FROM "Orders" WHERE "OrderID"=1 AND begins_with("Address", '7834 24th')`
- `function-contains-001` — `SELECT * FROM "Orders" WHERE "OrderID"=1 AND contains("Address", 'Kirkland')`
- `function-attribute-type-001` — `SELECT * FROM "Music" WHERE attribute_type("Artist", 'S')`

The summary log line at parser_test.go:413 (`AWS corpus: N fully parsed, M stubbed, K skipped`) is the regression marker — `N` should increase by 3 and `M` should decrease by 3 after the change lands.

### 3.3 Regression net

- `go test ./partiql/parser/...` — all unit + corpus tests pass.
- `go test ./partiql/...` — downstream analysis/completion packages don't crash on a real `*ast.FuncCall` (they previously never saw one outside the deferred-stub error path).

---

## 4. Risk

1. **Latent assumptions in analysis/completion.** Those packages parse via the same `parser.Parse` and walk the AST; they may not handle `*ast.FuncCall` in every visit path. Mitigation: run the broader `go test ./partiql/...` suite, not just the parser package.
2. **Reserved-name shadowing.** A user query like `SIZE(x)` looks like a generic call but tokenizes as `tokSIZE` — it routes through the keyword-dispatched stub in `exprprimary.go`, **not** through `parseVarRef`. 15a does NOT unblock `SIZE(x)` / `EXISTS(...)` / `UPPER(...)` etc.; those wait for 15b. Document this in the spec so it's not framed as a complete unblocker.
3. **`@foo(x)` rejection is a behavior change.** Pre-15a the omni parser returned the deferred-feature error; pre-cutover the legacy ANTLR parser would reject `@foo(x)` outright. We're matching ANTLR. No customer query is expected to use `@foo(x)`, but this is worth a sentence in the release notes / commit message.

---

## 5. Out of Scope

- All typed keyword builtins (15b).
- Aggregates (node 14) — `COUNT/SUM/AVG/MIN/MAX` keyword dispatch.
- Window functions (node 13) — `LAG/LEAD` keyword dispatch.
- The `caseSensitive` flag on the FuncCall name. SQL function names are case-insensitive in standard PartiQL; if 15b needs to distinguish quoted call names later, extend the AST then.
- Tests for `analysis.ValidateQuery` behavior on FuncCall — that path already returns no-error for DQL containing function calls; nothing to verify beyond the regression net.

---

## 6. Files Touched Summary

| File | Change |
|---|---|
| `partiql/parser/parser.go` | Modify `parseVarRef`; add `parseFuncCallArgs` helper |
| `partiql/parser/parser_test.go` | Flip `funcall_stub` case; add positive (6) + negative (4) cases |

The AWS corpus index needs no edit — its stubbed/fully-parsed counters are runtime classifications, not stored state.
