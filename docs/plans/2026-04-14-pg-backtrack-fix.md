# pg/parser backtracking unification — findings, detection, plan

> Status: pre-implementation. This doc captures the audit findings, the
> detection method, and the proposed fix plan for review BEFORE any code
> changes. Not yet associated with any branch other than `pg-first-sets`
> (where the audit was performed).

## Background

After completing the FIRST-set consolidation refactor on `pg-first-sets`,
a follow-up question surfaced: the original report listed
`CREATE FUNCTION f(double precision)` as a known limitation tied to a
`pushBack` design issue in `create_function.go:parseFuncArg`. The user
asked: "is this the only place where omni's parser has botched ambiguity
disambiguation, or are there others?"

This document captures the audit performed in response, the bugs found
(plural), and the proposed fix plan.

---

## Section 1 — Findings

### Finding 1: `pushBack` in `create_function.go:parseFuncArg` is lossy (already known)

**Symptom**: `CREATE FUNCTION f(double precision) RETURNS int AS '...' LANGUAGE sql` fails with `syntax error at or near "precision"`.

**Mechanism**:
- `parseFuncArg` (line 211) needs to disambiguate `param_name func_type` vs bare `func_type` at the function-argument grammar position.
- It uses a "speculative consume + rollback" strategy: call `parseTypeFunctionName()` to consume the leading token, then check whether what follows can start a `func_type`.
- If the rollback decision fires, it calls `pushBack(savedName)` (4 sites: lines 259, 263, 283, 286) where `savedName` is a **string**.
- `pushBack` is defined at lines 325-335 with signature `func (p *Parser) pushBack(name string)`. It writes back a synthetic token: `Token{Type: IDENT, Str: name, Loc: p.prev.Loc}`.
- **The original token type is lost** — it is unconditionally rewritten to `IDENT`.

**Why only `double precision` triggers** (this was a key clarification during the audit):

The bug requires BOTH conditions:
1. The leading token has to be `UnreservedKeyword` (so `isTypeFunctionName` returns true and the speculative-consume branch is entered).
2. The type has multi-token continuation (so the lossy IDENT type breaks the downstream re-parse).

Of the 21 tokens in `simpleTypenameLeadTokens`, exactly **one** is `UnreservedKeyword`: `DOUBLE_P`. All other type leads (`INT_P`, `INTEGER`, `BIT`, `CHARACTER`, `VARCHAR`, `NATIONAL`, `NCHAR`, `BOOLEAN_P`, `JSON`, `TIMESTAMP`, `TIME`, `INTERVAL`, etc.) are `ColNameKeyword`, so `isTypeFunctionName(...)` returns false and they bypass the speculative-consume branch entirely — they go directly to `parseFuncType` → `parseSimpleTypename` where the dispatch switch correctly handles their multi-token forms (`bit varying`, `character varying(10)`, `national character varying`, `time with time zone`, etc.).

**Empirical blast radius for finding 1** (corrected per codex review on commit 66a8cbe):

The bug fires in **every** `parseFuncArg` grammar alternative that ends with `double precision`. Codex re-tested and reproduced the failure for all five forms below — each enters the same speculative-consume + lossy-pushBack code path:

| SQL | Result |
|---|---|
| `CREATE FUNCTION f(double precision) RETURNS int AS '...' LANGUAGE sql` (case 5: bare `func_type`) | ❌ FAIL |
| `CREATE FUNCTION f(IN double precision) RETURNS int AS '...' LANGUAGE sql` (case 4: `arg_class func_type`) | ❌ FAIL |
| `CREATE FUNCTION f(OUT double precision) RETURNS int AS '...' LANGUAGE sql` | ❌ FAIL |
| `CREATE FUNCTION f(INOUT double precision) RETURNS int AS '...' LANGUAGE sql` | ❌ FAIL |
| `CREATE FUNCTION f(VARIADIC double precision) RETURNS int AS '...' LANGUAGE sql` | ❌ FAIL |
| `CREATE FUNCTION f(double precision DEFAULT 1.0) RETURNS int AS '...' LANGUAGE sql` (case 5 with default) | ❌ FAIL |
| `CREATE FUNCTION f(p double precision) RETURNS int AS '...' LANGUAGE sql` (case 3: `param_name func_type` with explicit name) | ✅ OK (the IDENT param name keeps speculative-consume harmless) |

So the **single root cause** (lossy `pushBack`) manifests across **6+ visible SQL forms**. My initial report listed only the bare case 5; the corrected blast radius is "any `CREATE FUNCTION` argument that uses `double precision` without an explicit parameter name".

Note: the bug is still **one bug**, not six. Fixing the `pushBack` lossiness fixes all six simultaneously. The expanded list affects test-matrix coverage, not the fix.

### Finding 2: `parseFuncType` speculative branch in `type.go:758-789` does partial state save (NEW, not previously known)

**Symptom**: any `parseFuncType` caller — which includes function parameters, function return types, `RETURNS TABLE` columns, and `CREATE OPERATOR` LEFTARG/RIGHTARG — that uses a 2-component schema-qualified type fails to parse.

**parseFuncType call sites (verified by grep + manual trace)**:

| File:line | Caller | Grammar context |
|---|---|---|
| `create_function.go:104` | `parseCreateFunctionStmt` | `RETURNS func_type` |
| `create_function.go:291` | `parseFuncArg` | function parameter type |
| `create_function.go:359` | `parseTableFuncColumn` | `RETURNS TABLE (col type, ...)` column type |
| `define.go:329` | `parseDefArg` | `CREATE OPERATOR ... LEFTARG = type`, `CREATE AGGREGATE ... STYPE = type`, etc. |

`parseDefArg` is itself called from `define.go:284`, `define.go:403`, `define.go:1442`, `create_table.go:1583`, `create_table.go:1600` — all within DefList/option-list parsing contexts. So the broken branch surfaces in **at least four user-visible grammar positions**, not just `CREATE FUNCTION`.

**Empirical blast radius (corrected per codex review on commit 66a8cbe)**:

| SQL | Result |
|---|---|
| `CREATE FUNCTION f(x pg_catalog.int4) RETURNS int AS '...' LANGUAGE sql` | ❌ FAIL |
| `CREATE FUNCTION f(x schema.mytype) RETURNS int AS '...' LANGUAGE sql` | ❌ FAIL |
| `CREATE FUNCTION f() RETURNS pg_catalog.int4 AS '...' LANGUAGE sql` | ❌ FAIL |
| `CREATE FUNCTION f() RETURNS TABLE (x pg_catalog.int4) AS '...' LANGUAGE sql` | ❌ FAIL (parseTableFuncColumn → parseFuncType) |
| `CREATE OPERATOR === (FUNCTION = int4eq, LEFTARG = pg_catalog.int4, RIGHTARG = int4)` | ❌ FAIL (parseDefArg → parseFuncType) |
| `CREATE OPERATOR === (FUNCTION = int4eq, LEFTARG = int4, RIGHTARG = int4)` | ✅ OK (unqualified — never enters the broken branch) |
| `CREATE FUNCTION f() RETURNS SETOF pg_catalog.int4 AS '...' LANGUAGE sql` | ✅ OK (SETOF takes a different code path that bypasses the broken branch) |
| `CREATE FUNCTION f() RETURNS schema.tab.col%TYPE AS '...' LANGUAGE sql` | ✅ OK (success path of the speculative branch — no rollback needed) |
| `CREATE TABLE t (c pg_catalog.int4)` | ✅ OK (uses `parseTypename` directly, not `parseFuncType`) |
| `SELECT CAST(NULL AS pg_catalog.int4)` | ✅ OK (uses `parseTypename` via expr) |
| `SELECT 1::pg_catalog.int4` | ✅ OK (TYPECAST) |
| `ALTER TABLE t ALTER COLUMN c TYPE pg_catalog.int4` | ✅ OK (uses `parseTypename`) |

**Out of evidence (not actually proof of finding 2 — see "Pre-existing 3-component name limitation" below)**:

| SQL | Status | Why excluded |
|---|---|---|
| `CREATE FUNCTION f() RETURNS db.schema.mytype AS '...' LANGUAGE sql` | ❌ FAIL | But `CREATE TABLE t (c db.schema.mytype)` ALSO fails. So 3-component names fail in **all** contexts, not just `parseFuncType`. This is a separate `parseGenericType` limitation, not the backtracking bug. |
| `SELECT CAST(NULL AS db.schema.mytype)` | ❌ FAIL | Same. |
| `ALTER TABLE t ALTER COLUMN c TYPE db.schema.mytype` | ❌ FAIL | Same. |

#### Pre-existing 3-component name limitation (out of scope for this fix)

Codex's audit revealed that omni's `parseGenericType` only supports **at most one dot** in qualified names. So `db.schema.mytype` fails in `CREATE TABLE`, `CAST`, `ALTER TABLE`, and `CREATE FUNCTION` alike — all because the underlying `parseGenericType` can't handle 3-component qualified names. This is a **separate independent bug**, not finding 2. It is **out of scope** for this PR.

If you fix the backtracking bug per this plan, `db.schema.mytype` will still fail in `CREATE FUNCTION` — but it will fail **for the same reason it fails everywhere else**, not because of the backtracking restore. To actually accept `db.schema.mytype`, `parseGenericType` needs to be extended to consume more dots. Track as a separate follow-up.

**Mechanism**:

`parseFuncType` at `type.go:720-792` handles the grammar:
```
func_type: Typename
         | type_function_name attrs '%' TYPE_P
         | SETOF type_function_name attrs '%' TYPE_P
```

The non-SETOF branch needs to disambiguate "qualified Typename" (e.g. `pg_catalog.int4`) from "%TYPE reference" (e.g. `schema.tab.col%TYPE`). Both start with `name '.' name`, so omni speculatively parses the `%TYPE` form first and rolls back if the `%` doesn't appear:

```go
// type.go:756-789
if next.Type == '.' {
    // Save state for backtracking.
    savedCur := p.cur
    savedPrev := p.prev
    savedNext := p.nextBuf
    savedHasNext := p.hasNext
    savedLexerErr := p.lexer.Err
    // ⚠️ NOT saved: lexer.pos, lexer.start, lexer.state

    name, _ := p.parseTypeFunctionName()
    if p.cur.Type == '.' {
        attrs, _ := p.parseAttrs()  // ← this advances the lexer
        if p.cur.Type == '%' {
            p.advance()
            if p.cur.Type == TYPE_P {
                p.advance()
                // success path — return without restoring
                return ...
            }
        }
    }

    // Not %TYPE pattern — restore and parse as Typename
    p.cur = savedCur
    p.prev = savedPrev
    p.nextBuf = savedNext
    p.hasNext = savedHasNext
    p.lexer.Err = savedLexerErr
    // ⚠️ lexer.pos, lexer.start, lexer.state NOT restored
}
return p.parseTypename()
```

**Why this corrupts state**:
- After saving state, `parseTypeFunctionName` consumes one token (cur becomes nextBuf=`.`, hasNext=false).
- `parseAttrs` then calls `advance()` to consume `.`, then expects an identifier → fresh `lexer.NextToken()` call → **lexer.pos advances past savedNext**.
- For deeper qualified names (3+ components), parseAttrs calls advance() multiple times, each potentially triggering fresh lexer reads.
- On the failure path, restore puts cur/prev/nextBuf back, but `lexer.pos` is wherever parseAttrs left it.
- After restore, the next `advance()` that exhausts cur+nextBuf calls `lexer.NextToken()` at the wrong position → reads the wrong token → syntax error downstream.

**Comparison with the correct pattern in the same codebase**:

`define.go:21-56` `parseDefineStmtAggregate` does an analogous speculative parse and saves the **complete** state:

```go
savedCur := p.cur
savedPrev := p.prev
savedNext := p.nextBuf
savedHasNext := p.hasNext
savedLexerErr := p.lexer.Err
savedLexerPos := p.lexer.pos       // ← present
savedLexerStart := p.lexer.start   // ← present
savedLexerState := p.lexer.state   // ← present
```

The `define.go` author understood that lexer state must be saved. The `type.go` author of the speculative branch did not. The two patterns coexist in the same codebase, written by different authors, with no shared abstraction.

### Finding 3 (informational, not a bug): CTAS rollback in `parser.go:894-918` is the correct pattern

`parseCreateStmt` needs to disambiguate `CREATE TABLE t (col_def, ...)` from `CREATE TABLE t (col_name, col_name) AS SELECT ...` (CTAS). Its rollback uses the proper `Token` struct + `nextBuf` pushback:

```go
savedPrev := p.prev
firstName := p.advance() // consume ColId; firstName is the full Token

if p.cur.Type == ',' || p.cur.Type == ')' {
    // CTAS column list — no rollback
}

// CREATE TABLE column def — rollback
p.nextBuf = p.cur          // push the post-advance token into lookahead buffer
p.hasNext = true
p.cur = firstName          // restore the original Token (with type AND str)
p.prev = savedPrev
```

This is **lossless**:
- `firstName` is a full `Token` value (struct), so `Type`, `Str`, `Loc` are all preserved
- The advance() that was performed is undone by pushing the post-advance token into nextBuf
- Stream-equivalence holds: the parser will re-emit `firstName`, then the post-advance token, then resume from the lexer at its current position

**This is exactly how `pushBack` should have been written** — same codebase, same author pool, but never extracted as a reusable helper.

There is **one** theoretical edge case: if `advance()` did `_LA` reclassification (NOT/WITH/WITHOUT/NULLS_P) and pre-populated `nextBuf` via `peekNext()`, then the rollback's `p.nextBuf = p.cur` would overwrite the `_LA`-peeked token. But the trigger condition is "CREATE TABLE `(` followed immediately by `NOT`/`WITH`/`WITHOUT`/`NULLS`", which is not a valid column definition lead in any PG version. **Theoretically lossy, practically unreachable**. Not classified as a bug, but worth a one-line comment.

### Summary of all 5 backtracking patterns in pg/parser

| # | Location | Pattern | Status |
|---|---|---|---|
| 1 | `create_function.go:325` `pushBack` | Token type hardcoded to IDENT | 🔴 LOSSY — finding 1 |
| 2 | `type.go:758-789` `parseFuncType` speculative | Partial snapshot (no `lexer.pos`) | 🔴 LOSSY — finding 2 |
| 3 | `define.go:21-56` `parseDefineStmtAggregate` | Full snapshot incl. lexer state | ✅ LOSSLESS |
| 4 | `parser.go:894-918` `parseCreateStmt` CTAS | Token struct + nextBuf push | ✅ LOSSLESS (with documented unreachable edge) |
| 5 | `lexer.go:786-790` / `1177-1214` `savedPos` | Lexer-internal pos rollback for read-only peeks | ✅ LOSSLESS |

**Two real bugs found, both reproducible.** Three correct backtrack patterns coexist with the two broken ones.

---

## Section 2 — Detection method

### Why the audit was scoped to "backtracking" rather than "ambiguity"

The user asked: "are there other ambiguity disambiguations done wrong?". The naive interpretation would search for "places where the parser disambiguates two grammar productions". That space is huge — every conditional in a parse function disambiguates something. The narrower, tractable framing is: "places where the parser SPECULATIVELY CONSUMES tokens and might need to ROLL BACK". This is a finite, greppable pattern, and it's exactly the failure mode we already knew about.

So: the audit hunts for speculative-consume + rollback machinery, not for grammar ambiguity in general.

### Search pattern (5 greps)

```bash
# 1. Direct uses of pushBack (the known lossy mechanism)
grep -n "pushBack" pg/parser/*.go | grep -v _test.go

# 2. saved* variables suggesting rollback intent
grep -nE "saved(Name|Tok|Cur|Pos|Prev)|savedLoc|savedNext|savedHasNext|savedLexer" \
    pg/parser/*.go | grep -v _test.go

# 3. Direct manipulation of nextBuf (bypassing peekNext)
grep -n "nextBuf\b" pg/parser/*.go | grep -v _test.go

# 4. Manual reassignment of p.cur outside advance()
grep -n "p\.cur = " pg/parser/*.go | grep -v _test.go | grep -v "p\.cur = p\."

# 5. parseTypeFunctionName call sites (the standard speculative-consume entry point)
grep -n "parseTypeFunctionName()" pg/parser/*.go | grep -v _test.go
```

Each grep returned a small enough result set (under 30 lines) to read by hand. The union of the first 4 greps captured every location where parser state was saved and restored. Grep 5 narrowed it to the speculative-consume sites specifically.

### Classification heuristic

For each candidate site, classify as:

- **LOSSY**: the rollback restores some state but not all. The bug fires when the speculative parse advances the lexer or modifies state that isn't covered by the restore.
- **LOSSLESS**: the rollback restores all state that the speculative parse could have modified. The proof is: list what the speculative parse calls; for each call, list what state it touches; verify each is in the restore set.
- **CORRECT-BY-CONSTRUCTION**: the "rollback" is just a `Token` struct copy + `nextBuf` push. No state ever needs restoring because the parser only consumed 1 token and the post-consume token is preserved in `nextBuf`.

### Empirical verification

For each LOSSY classification, write a small Go program that calls `parser.Parse` on the suspected SQL and check whether it returns an error. This was the only reliable way to confirm the bugs — code-reading alone could be wrong about whether the lexer state mismatch actually manifests.

Test programs used (saved as `/tmp/test_*.go`, not committed):

```go
// /tmp/test_paramless.go — verify finding 1 blast radius
sqls := []string{
    `CREATE FUNCTION f(double precision) RETURNS int AS 'select 1' LANGUAGE sql`,
    `CREATE FUNCTION f(bit varying(10)) RETURNS int AS 'select 1' LANGUAGE sql`,
    `CREATE FUNCTION f(time with time zone) RETURNS int AS 'select 1' LANGUAGE sql`,
    // ... 11 cases total covering every multi-word type variant
}
// Result: only `double precision` failed.

// /tmp/test_qualified_type.go — verify finding 2 blast radius
sqls := []string{
    `CREATE FUNCTION f() RETURNS pg_catalog.int4 AS 'select 1' LANGUAGE sql`,
    `CREATE FUNCTION f() RETURNS schema.mytype AS 'select 1' LANGUAGE sql`,
    `CREATE FUNCTION f() RETURNS db.schema.mytype AS 'select 1' LANGUAGE sql`,
    `CREATE FUNCTION f() RETURNS schema.tab.col%TYPE AS 'select 1' LANGUAGE sql`,
    `CREATE FUNCTION f() RETURNS int AS 'select 1' LANGUAGE sql`,
}
// Result: 3 fail, 2 pass — confirmed `func_type` non-%TYPE qualified branch is broken.

// /tmp/test_qualified_type2.go — verify cross-grammar isolation
sqls := []string{
    `CREATE FUNCTION f(x pg_catalog.int4) RETURNS int AS '...' LANGUAGE sql`,  // FAIL
    `CREATE FUNCTION f() RETURNS pg_catalog.int4 AS '...' LANGUAGE sql`,        // FAIL
    `CREATE FUNCTION f() RETURNS SETOF pg_catalog.int4 AS '...' LANGUAGE sql`,  // OK (SETOF bypass)
    `SELECT CAST(NULL AS pg_catalog.int4)`,                                     // OK (parseTypename)
    `CREATE TABLE t (c pg_catalog.int4)`,                                       // OK (parseTypename)
    `ALTER TABLE t ALTER COLUMN c TYPE pg_catalog.int4`,                        // OK (parseTypename)
    `SELECT 1::pg_catalog.int4`,                                                // OK (TYPECAST)
}
// Confirms the bug is exactly localized to parseFuncType's broken speculative branch.
```

### What the audit would have missed

1. **Backtracking that doesn't use `saved*` naming convention.** If a future contributor wrote `oldCur := p.cur; ... p.cur = oldCur`, grep 2 would miss it. Mitigation: greps 3 and 4 (nextBuf manipulation and direct p.cur reassignment) are wider nets that catch any rollback shape.

2. **Rollback via panic/recover or via returning a sentinel.** None observed in pg/parser, but worth noting for completeness.

3. **Speculative parse that doesn't roll back at all** — i.e., the parser blindly commits and produces wrong AST. Grep 5 would catch the consume sites, but grep 5 returned 6 hits and only 2 were speculative; the other 4 were genuine first-position parses. Manual review of each of the 4 confirmed they're not rollback candidates.

4. **Backtracking inside non-pg parsers.** This audit is scoped to `pg/parser/`. mysql/mssql/oracle parsers may have analogous bugs and warrant separate audits. **Out of scope for this fix.**

---

## Section 3 — Plan

### Goal

Fix both bugs by introducing a single shared backtracking helper, then migrating the two broken sites and the one verbose-but-correct site (`define.go`) to use it. After the fix:

- `pushBack(string)` is deleted.
- `parseFuncArg` either uses the new helper for backtracking, OR is rewritten to peek-then-commit (to be decided during implementation; see "Open question" below).
- `parseFuncType`'s speculative branch uses the new helper for complete state save/restore.
- `define.go`'s 8-line manual snapshot is replaced by 1 helper call.
- A new test file covers both bug classes against PG 17 via the existing oracle infrastructure where applicable, plus AST-shape unit tests for both regression sets.

### The helper

**Important scope note (per codex review)**: this helper is a **token-stream snapshot**, NOT a complete parser/lexer snapshot. It captures exactly the state needed to rewind the token stream to a previous **token boundary** during ordinary parsing. It does NOT cover:

- Lexer mid-token-content state: `stateBeforeStrStop`, `literalbuf`, `dolqstart`, `utf16FirstPart`, `xcdepth`, lexer warning flags. These fields are reset to known values at every token boundary, so they don't need saving for token-stream rollback. They WOULD need saving if a future caller wanted to roll back from inside a string literal or dollar-quoted block — that's not a use case omni currently has.
- Completion-mode state: `candidates`, `collecting`, etc. used during IDE completion mode. The current speculative parses don't run in completion mode.
- Any state in `Parser` that's not part of the token stream view (e.g., error accumulators if any).

Naming convention: the type and methods are deliberately named `tokenStreamState` / `snapshotTokenStream` / `restoreTokenStream` to make the limited scope visible at every call site. A future "complete parser snapshot" would be a different, larger struct with different methods.

```go
// pg/parser/backtrack.go (new file)

// tokenStreamState captures the parser + lexer state needed to rewind the
// token stream to a previous token boundary. This is sufficient for the
// "speculative parse, then rollback if it doesn't match" pattern at every
// site in pg/parser as of this writing.
//
// SCOPE: this is a TOKEN-STREAM snapshot, not a complete parser/lexer
// snapshot. It does not cover mid-token-content lexer state (literalbuf,
// dolqstart, utf16FirstPart, xcdepth, etc.) or completion-mode state
// (candidates, collecting). Those fields are either reset at token
// boundaries (lexer internals) or not used by speculative parses
// (completion mode), so they don't need to be saved here.
//
// If a future caller needs to roll back from INSIDE a token (e.g., from
// inside a string literal or dollar-quoted block), this struct is
// insufficient — extend it carefully.
//
// Why this exists: see commit history for the create_function pushBack
// incident and the type.go parseFuncType partial-snapshot incident.
// Both bugs were caused by hand-rolled rollback machinery that captured
// only a subset of the necessary state.
type tokenStreamState struct {
    cur, prev, nextBuf Token
    hasNext            bool
    lexerErr           error
    lexerPos           int
    lexerStart         int
    lexerState         lexerStateMode  // whatever type pg/parser/lexer.go uses
}

// snapshotTokenStream captures the current token-stream position for
// later restoration via restoreTokenStream. See tokenStreamState for
// scope and limitations.
func (p *Parser) snapshotTokenStream() tokenStreamState {
    return tokenStreamState{
        cur:        p.cur,
        prev:       p.prev,
        nextBuf:    p.nextBuf,
        hasNext:    p.hasNext,
        lexerErr:   p.lexer.Err,
        lexerPos:   p.lexer.pos,
        lexerStart: p.lexer.start,
        lexerState: p.lexer.state,
    }
}

// restoreTokenStream rewinds parser+lexer state to a previously captured
// snapshot. After restore, the next advance() will emit the same token as
// it would have at the moment snapshotTokenStream() was called.
//
// Caller responsibility: do not interleave restore with completion-mode
// queries or with any operation that mutates lexer state outside the
// token stream (string literal scanning, etc.).
func (p *Parser) restoreTokenStream(s tokenStreamState) {
    p.cur = s.cur
    p.prev = s.prev
    p.nextBuf = s.nextBuf
    p.hasNext = s.hasNext
    p.lexer.Err = s.lexerErr
    p.lexer.pos = s.lexerPos
    p.lexer.start = s.lexerStart
    p.lexer.state = s.lexerState
}
```

The exact field types depend on what `Lexer` exports — the implementer should verify `lexer.state`'s type name in `lexer.go` before writing this. If the type is unexported, the helper file lives in the same package so it's accessible directly.

### Commit breakdown

**Commit 1**: Introduce `pg/parser/backtrack.go` with `parserState`, `snapshot()`, `restore()`. No call site changes. Add a basic test that snapshot/restore is a no-op when no state changes between save and restore.

**Commit 2**: Fix finding 2 (`parseFuncType` qualified type bug) by replacing the partial snapshot at `type.go:758-789` with `snapshot()`/`restore()`. Add `TestParseCreateFunctionQualifiedTypes` covering 6+ cases (param qualified, return qualified, deep qualified, sanity SETOF, sanity %TYPE, sanity unqualified) with AST-shape assertions.

**Commit 3**: Fix finding 1 (`parseFuncArg` `double precision` bug) by either (a) using `snapshot()`/`restore()` to replace the lossy `pushBack`, or (b) rewriting `parseFuncArg` to peek-then-commit using `peekNext()`. Decision criterion: if option (b) is < 80 lines and behavior-preserving, prefer it (eliminates the speculative-consume anti-pattern entirely). Otherwise option (a). Either way, delete `pushBack(string)`. Add `TestParseCreateFunctionParamlessDoublePrecision` with AST-shape assertions.

**Commit 4** (optional, can be deferred): Migrate `define.go`'s manual snapshot to the helper. Pure refactor, no behavior change. Reduces 8 lines of save + 8 lines of restore to 1 + 1.

### Open question: option (a) vs option (b) for parseFuncArg

The deeper anti-pattern in `parseFuncArg` is "speculative consume + rollback" where peek-then-commit would suffice. omni already has 2-token lookahead via `peekNext()` (used in 9 other places, including `_LA` token reclassification and `parseCreateDispatch`). A peek-then-commit rewrite would:

- Look at `cur` and `peekNext()` together
- If `cur` is a `type_function_name` AND `peekNext()` is a token in `FIRST(func_type)`: `cur` is `param_name`, consume it, then parse type
- If `cur` is a `type_function_name` AND `peekNext()` is NOT in `FIRST(func_type)`: `cur` IS the type, don't consume, fall through to `parseFuncType`
- Eliminates the need for `pushBack` entirely

This is the bison-style approach. The risk is that `parseFuncArg` has 5 grammar alternatives and 4 existing pushBack sites — the rewrite needs to handle each correctly. The implementer should:

1. Read the current parseFuncArg body (~80 lines)
2. Sketch the peek-then-commit version
3. If it's ≤ 100 lines and clearly behavior-preserving, choose option (b)
4. Otherwise choose option (a) — same fix at the rollback level, deferred backtracking-elimination

Both options fix bug 1. Option (b) is more elegant; option (a) is more conservative.

### Test matrix

```go
// pg/parser/create_function_extra_test.go (new file)

func TestParseCreateFunctionDoublePrecisionAllArgClasses(t *testing.T) {
    // Per codex review: bug 1 fires in EVERY parseFuncArg grammar
    // alternative that ends with `double precision`, not just the bare
    // case 5. Test all 6 visible forms.
    type testCase struct {
        sql       string
        paramName string                       // "" for paramless
        mode      nodes.FunctionParameterMode  // FUNC_PARAM_DEFAULT / IN / OUT / INOUT / VARIADIC
        typeLast  string                       // last component of parsed TypeName
    }
    cases := []testCase{
        // case 5: bare func_type
        {
            sql:      `CREATE FUNCTION f(double precision) RETURNS int AS 'select 1' LANGUAGE sql`,
            mode:     nodes.FUNC_PARAM_DEFAULT,
            typeLast: "float8",
        },
        // case 4: arg_class func_type
        {
            sql:      `CREATE FUNCTION f(IN double precision) RETURNS int AS 'select 1' LANGUAGE sql`,
            mode:     nodes.FUNC_PARAM_IN,
            typeLast: "float8",
        },
        {
            sql:      `CREATE FUNCTION f(OUT double precision) RETURNS int AS 'select 1' LANGUAGE sql`,
            mode:     nodes.FUNC_PARAM_OUT,
            typeLast: "float8",
        },
        {
            sql:      `CREATE FUNCTION f(INOUT double precision) RETURNS int AS 'select 1' LANGUAGE sql`,
            mode:     nodes.FUNC_PARAM_INOUT,
            typeLast: "float8",
        },
        {
            sql:      `CREATE FUNCTION f(VARIADIC double precision) RETURNS int AS 'select 1' LANGUAGE sql`,
            mode:     nodes.FUNC_PARAM_VARIADIC,
            typeLast: "float8",
        },
        // case 5 with default value
        {
            sql:      `CREATE FUNCTION f(double precision DEFAULT 1.0) RETURNS int AS 'select 1' LANGUAGE sql`,
            mode:     nodes.FUNC_PARAM_DEFAULT,
            typeLast: "float8",
        },
        // sanity: case 3 (param_name func_type) with explicit name (must keep working)
        {
            sql:       `CREATE FUNCTION f(p double precision) RETURNS int AS 'select 1' LANGUAGE sql`,
            paramName: "p",
            mode:      nodes.FUNC_PARAM_DEFAULT,
            typeLast:  "float8",
        },
        // sanity: paramless single-token type (must keep working)
        {
            sql:      `CREATE FUNCTION f(int) RETURNS int AS 'select 1' LANGUAGE sql`,
            mode:     nodes.FUNC_PARAM_DEFAULT,
            typeLast: "int4",
        },
    }
    // ... unwrap RawStmt → CreateFunctionStmt → FunctionParameter, assert
    //     Name == paramName, Mode == mode, ArgType.Names[-1].Str == typeLast ...
}

func TestParseQualifiedTypeInFuncTypePositions(t *testing.T) {
    // Bug 2 surfaces in EVERY parseFuncType caller, not just CREATE FUNCTION
    // params/returns. Test all four caller paths to lock down the full
    // blast radius.
    //
    // Excluded: db.schema.mytype (3-component qualified names fail in
    // parseGenericType in ALL contexts, not just parseFuncType — separate
    // bug, separate PR).
    cases := []struct {
        sql        string
        callerKind string  // diagnostic label
    }{
        // create_function.go:291 — function parameter type
        {
            sql:        `CREATE FUNCTION f(x pg_catalog.int4) RETURNS int AS 'select 1' LANGUAGE sql`,
            callerKind: "parseFuncArg argType",
        },
        {
            sql:        `CREATE FUNCTION f(x schema.mytype) RETURNS int AS 'select 1' LANGUAGE sql`,
            callerKind: "parseFuncArg argType (custom schema)",
        },
        // create_function.go:104 — function return type
        {
            sql:        `CREATE FUNCTION f() RETURNS pg_catalog.int4 AS 'select 1' LANGUAGE sql`,
            callerKind: "parseCreateFunctionStmt RETURNS",
        },
        // create_function.go:359 — RETURNS TABLE column type
        {
            sql:        `CREATE FUNCTION f() RETURNS TABLE (x pg_catalog.int4) AS 'select 1' LANGUAGE sql`,
            callerKind: "parseTableFuncColumn",
        },
        // define.go:329 — DefArg via CREATE OPERATOR
        {
            sql:        `CREATE OPERATOR === (FUNCTION = int4eq, LEFTARG = pg_catalog.int4, RIGHTARG = int4)`,
            callerKind: "parseDefArg LEFTARG",
        },
        {
            sql:        `CREATE OPERATOR === (FUNCTION = int4eq, LEFTARG = int4, RIGHTARG = pg_catalog.int4)`,
            callerKind: "parseDefArg RIGHTARG",
        },
        // sanity: SETOF qualified — bypasses speculative branch via SETOF prefix
        {
            sql:        `CREATE FUNCTION f() RETURNS SETOF pg_catalog.int4 AS 'select 1' LANGUAGE sql`,
            callerKind: "SETOF qualified (sanity)",
        },
        // sanity: %TYPE form — success path of speculative branch
        {
            sql:        `CREATE FUNCTION f() RETURNS schema.tab.col%TYPE AS 'select 1' LANGUAGE sql`,
            callerKind: "%TYPE form (sanity)",
        },
        // sanity: unqualified type — never enters speculative branch
        {
            sql:        `CREATE FUNCTION f() RETURNS int AS 'select 1' LANGUAGE sql`,
            callerKind: "unqualified (sanity)",
        },
    }
    for _, tc := range cases {
        t.Run(tc.callerKind, func(t *testing.T) {
            stmts, err := Parse(tc.sql)
            if err != nil {
                t.Fatalf("parse failed: %v", err)
            }
            // ... AST assertions per callerKind ...
            _ = stmts
        })
    }
}
```

The assertion functions should follow the pattern established by `TestParseCreateFunctionWithColNameKeywordParam` in `create_function_json_test.go` (commit `549a1c0`): unwrap `RawStmt` → `CreateFunctionStmt` → `FunctionParameter` → `ArgType.Names`, assert each level individually so a regression at any level produces a clear error.

### Risks

1. **Lexer state typing.** The helper assumes `lexer.state` has a stable type. If it's an unexported type or a function pointer, the snapshot may need a different shape. Implementer must verify before writing the helper.

2. **define.go migration may fight with existing tests.** The migration is supposed to be behavior-equivalent, but if the existing manual snapshot was actually relying on some quirk (e.g., NOT restoring something on purpose), the helper migration could regress. Read the existing test coverage for `parseDefineStmtAggregate` before migrating. If unclear, defer commit 4 and ship commits 1-3 only.

3. **parseFuncArg peek-then-commit rewrite scope creep.** If option (b) is chosen and the rewrite turns out to be more complex than expected (e.g., 200 lines), abandon it and use option (a). Don't get stuck.

4. **Bug 2 fix may surface OTHER pre-existing bugs.** If `parseFuncType`'s speculative branch was previously masking some other issue (e.g., a `parseAttrs` bug only caught by the broken restore), fixing the restore may surface that issue. Run the full pg/parser suite + the new tests after commit 2 and triage any new failures before continuing.

5. **The helper itself could be wrong.** The exact set of fields to save is currently inferred from `define.go`'s pattern, which is the most complete one in the codebase. But there could be lexer state I haven't seen. Mitigation: write a `TestSnapshotRestoreIsIdentity` test that calls snapshot, advances 5 tokens, restores, and asserts that re-advancing 5 tokens produces the same sequence as before. Run it on a corpus of representative SQL to flush out any missed state.

### Time estimate

| Phase | Time |
|---|---|
| Commit 1 (helper + identity test) | 1 hour |
| Commit 2 (parseFuncType fix + qualified-type tests) | 2 hours |
| Commit 3 (parseFuncArg fix + double-precision tests, option (a) or (b)) | 2-4 hours |
| Commit 4 (define.go migration, optional) | 1 hour |
| Final smoke + codex review | 1 hour |
| **Total** | **6-9 hours** (1 day) |

Significantly more than the original "30-minute 1-line workaround" estimate, but covers two bugs and eliminates the underlying anti-pattern.

### Deferred / out of scope

- mysql/parser, mssql/parser, oracle/parser audits for the same anti-pattern. Each may need its own audit.
- Removing the `parser.go:894-918` CTAS-rollback edge case (NOT/WITH/WITHOUT/NULLS_P first-token after `(` overwriting `_LA`-peeked nextBuf). Not migrated to the helper because the CTAS form is more compact and the edge case is unreachable in practice. Document with a one-line comment on the rollback site.
- Phase 3 audit follow-ups (3 `isColId() || isTypeFunctionName()` call sites flagged in commit `fc43201`). Tracked separately.

---

## Review questions for codex

1. **Are findings 1 and 2 actually independent**, or could they share a deeper root cause I missed?
2. **Is the audit detection method complete enough**? What backtracking patterns could it miss that I should grep for?
3. **Is the helper signature correct**? Specifically, is there any parser state I'm missing that should be in `parserState`?
4. **Is option (b) — peek-then-commit rewrite of parseFuncArg — a sound idea**, or are there grammar alternatives in `func_arg` that genuinely need backtracking and can't be peeked?
5. **Are the test cases sufficient**? Specifically, is there a category of qualified type or paramless multi-word type I'm not testing?
6. **Is there a simpler fix I'm not seeing**? E.g., could finding 2 be fixed by NOT speculative-parsing at all and instead deferring `%TYPE` detection to after parseTypename returns?
7. **Should commit 4 (define.go migration) be done**, or is the risk of regression higher than the value of the cleanup?
