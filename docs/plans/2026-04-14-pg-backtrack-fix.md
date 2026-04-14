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

**Empirical blast radius for finding 1**: exactly 1 SQL pattern, `CREATE FUNCTION f(double precision) RETURNS int AS '...' LANGUAGE sql`. Verified by direct test against omni at HEAD of `pg-first-sets`.

### Finding 2: `parseFuncType` speculative branch in `type.go:758-789` does partial state save (NEW, not previously known)

**Symptom**: any `CREATE FUNCTION` whose parameter or return type is schema-qualified (`pg_catalog.int4`, `schema.mytype`, `db.schema.mytype`, ...) fails to parse.

**Empirical blast radius**:

| SQL | Result |
|---|---|
| `CREATE FUNCTION f(x pg_catalog.int4) RETURNS int AS '...' LANGUAGE sql` | ❌ FAIL |
| `CREATE FUNCTION f(x schema.mytype) RETURNS int AS '...' LANGUAGE sql` | ❌ FAIL |
| `CREATE FUNCTION f() RETURNS pg_catalog.int4 AS '...' LANGUAGE sql` | ❌ FAIL |
| `CREATE FUNCTION f() RETURNS db.schema.mytype AS '...' LANGUAGE sql` | ❌ FAIL |
| `CREATE FUNCTION f() RETURNS SETOF pg_catalog.int4 AS '...' LANGUAGE sql` | ✅ OK (SETOF takes a different code path that bypasses the broken branch) |
| `CREATE FUNCTION f() RETURNS schema.tab.col%TYPE AS '...' LANGUAGE sql` | ✅ OK (success path of the speculative branch — no rollback needed) |
| `CREATE TABLE t (c pg_catalog.int4)` | ✅ OK (uses `parseTypename` directly, not `parseFuncType`) |
| `SELECT CAST(NULL AS pg_catalog.int4)` | ✅ OK (uses `parseTypename` via expr) |
| `SELECT 1::pg_catalog.int4` | ✅ OK (TYPECAST) |
| `ALTER TABLE t ALTER COLUMN c TYPE pg_catalog.int4` | ✅ OK (uses `parseTypename`) |

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

```go
// pg/parser/backtrack.go (new file)

// parserState is a complete snapshot of parser + lexer state. It is the
// minimum information needed to rewind the parser to a previous position
// after a speculative parse. Use snapshot/restore (see below) rather than
// hand-rolling state save/restore — incomplete snapshots are a known
// source of bugs (see commit history for the create_function pushBack
// and type.go parseFuncType incidents).
type parserState struct {
    cur, prev, nextBuf Token
    hasNext            bool
    lexerErr           error
    lexerPos           int
    lexerStart         int
    lexerState         lexerStateMode  // whatever type pg/parser/lexer.go uses
}

// snapshot captures the current parser+lexer state for later restoration.
func (p *Parser) snapshot() parserState {
    return parserState{
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

// restore rewinds parser+lexer state to a previously captured snapshot.
// After restore, the next advance() will emit the same token as it would
// have at the moment snapshot() was called.
func (p *Parser) restore(s parserState) {
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

The exact field types depend on what `Lexer` exports — the implementer should verify `lexer.state`'s type name in `lexer.go` before writing this.

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

func TestParseCreateFunctionParamlessDoublePrecision(t *testing.T) {
    cases := []struct {
        sql       string
        paramName string  // "" for paramless
        typeLast  string  // last component of parsed TypeName
    }{
        {
            sql:      `CREATE FUNCTION f(double precision) RETURNS int AS 'select 1' LANGUAGE sql`,
            typeLast: "float8",
        },
        // sanity: with explicit param name (must keep working)
        {
            sql:       `CREATE FUNCTION f(p double precision) RETURNS int AS 'select 1' LANGUAGE sql`,
            paramName: "p",
            typeLast:  "float8",
        },
        // sanity: paramless single-token type (must keep working)
        {
            sql:      `CREATE FUNCTION f(int) RETURNS int AS 'select 1' LANGUAGE sql`,
            typeLast: "int4",
        },
    }
    // ... unwrap RawStmt → CreateFunctionStmt → FunctionParameter, assert AST shape ...
}

func TestParseCreateFunctionQualifiedTypes(t *testing.T) {
    cases := []struct {
        sql        string
        paramName  string
        typeNames  []string  // full Names list of parsed TypeName
    }{
        // bug 2 — param qualified
        {
            sql:       `CREATE FUNCTION f(x pg_catalog.int4) RETURNS int AS 'select 1' LANGUAGE sql`,
            paramName: "x",
            typeNames: []string{"pg_catalog", "int4"},
        },
        {
            sql:       `CREATE FUNCTION f(x schema.mytype) RETURNS int AS 'select 1' LANGUAGE sql`,
            paramName: "x",
            typeNames: []string{"schema", "mytype"},
        },
        // bug 2 — return qualified
        {
            sql:       `CREATE FUNCTION f() RETURNS pg_catalog.int4 AS 'select 1' LANGUAGE sql`,
            typeNames: []string{"pg_catalog", "int4"},
        },
        // bug 2 — deep qualified
        {
            sql:       `CREATE FUNCTION f() RETURNS db.schema.mytype AS 'select 1' LANGUAGE sql`,
            typeNames: []string{"db", "schema", "mytype"},
        },
        // sanity: SETOF qualified (works today, must keep working)
        {
            sql:       `CREATE FUNCTION f() RETURNS SETOF pg_catalog.int4 AS 'select 1' LANGUAGE sql`,
            typeNames: []string{"pg_catalog", "int4"},
        },
        // sanity: %TYPE form (success path of speculative branch, must keep working)
        {
            sql: `CREATE FUNCTION f() RETURNS schema.tab.col%TYPE AS 'select 1' LANGUAGE sql`,
            // ... pct_type assertion ...
        },
    }
    // ... AST shape assertions ...
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
