# pg/parser maintenance conventions

This document captures conventions specific to omni's pg parser package. It assumes you are familiar with the package's overall layout and the upstream PostgreSQL grammar (`postgres/src/backend/parser/gram.y`).

## FIRST-set discipline

omni's pg parser is hand-written recursive descent. Every disambiguation point that asks "could the current token start grammar production X?" is a **FIRST-set probe**.

**Rule:** every FIRST-set probe MUST go through a function named `isXStart(...)` defined in `first_sets.go`. Token enumerations may NOT be repeated at probe call sites.

**Why this rule exists:** before April 2026, FIRST sets for the type productions were duplicated across three places:
- `parseSimpleTypename`'s switch (the actual dispatch in `type.go`)
- `isAExprConstTypeCast`'s whitelist (`expr.go`, used as a typed-literal probe)
- `isBuiltinType`'s whitelist (`create_function.go`, used in `parseFuncArg` to disambiguate `param_name func_type`)

The third copy drifted: it was missing `JSON`, `DEC`, `NATIONAL`, `NCHAR`. A bytebase customer hit this when a generated `CREATE FUNCTION f(data json) ...` failed with `syntax error at or near "json"`. The fix was to introduce a single source of truth (`isSimpleTypenameStart` backed by `simpleTypenameLeadTokens`) and make every probe go through it. The rule above prevents the same drift from recurring.

### Enforcement

Drift is caught at three layers:

1. **PG testcontainer oracle tests** (`first_set_oracle_test.go`) probe each predicate against a real PG 17 instance and assert the predicate's accept set matches PG's. Any drift between omni's slice and PG's grammar fails CI loudly. Not all productions can be oracle-tested — see "Why no PG oracle for X" below.

2. **In-Go bidirectional unit tests** for productions that cannot be oracle-tested (currently only `ConstTypename`, because PG's `AexprConst` grammar position is ambiguous with `func_name Sconst`). These iterate the entire `Keywords` table and assert the predicate matches the slice exactly.

3. **Code review:** any PR introducing a bare `case TOK_A, TOK_B, TOK_C:` cluster outside a dispatch switch should be flagged. If you're tempted to write such a cluster, ask: "is this the FIRST set of some production?" If yes, route it through `isXStart`. If no, document in a comment why it's not (operator categorization, error recovery, etc.).

## The six canonical predicates

All defined in `first_sets.go` and documented inline. Read the doc comment on each before adding a new caller. (`isSelectStart` lives in `expr.go` and predates this refactor; it was reused at the migrated `copy.go` and `publication.go` call sites but is not part of the `first_sets.go` family.)

| Predicate | Backing slice | Composes from | Has PG oracle test? |
|---|---|---|---|
| `isSimpleTypenameStart` | `simpleTypenameLeadTokens` (21) | set + `isTypeFunctionName` (GenericType path) | Yes — `TestSimpleTypenameLeadTokensMatchPG` via `CAST(NULL AS X)` |
| `isConstTypenameStart` | reuses `simpleTypenameLeadSet` | set only (no fallthrough) | No — see "Why no PG oracle for ConstTypename" below; in-Go bidirectional test instead |
| `isTypenameStart` | reuses `simpleTypenameLeadSet` | `SETOF \|\| isSimpleTypenameStart` | Yes — `TestTypenameLeadTokensMatchPG` via `RETURNS X` |
| `isFuncTypeStart` | (delegate) | thin alias for `isTypenameStart` | Yes — `TestFuncTypeLeadTokensMatchPG` via param type position |
| `isAExprStart` | `aExprLeadTokens` (32) | set + `isConstTypenameStart` + `isColId` | Yes — `TestAExprLeadTokensMatchPG` via `SELECT X` |
| `isTableConstraintStart` | `tableConstraintLeadTokens` (6) | set only | No — TableConstraint is only reachable via CREATE/ALTER TABLE; in-Go positive/negative test instead |

## Oracle probe templates

Each oracle test exercises its production at a chosen grammar position and uses PG's response code to classify "syntactically accepted" vs "syntactically rejected". The classifier is in `first_set_oracle_test.go`:

- SQLSTATE `42601` (`syntax_error`) → `probeReject`
- Anything else (success, undefined object, type mismatch, etc.) → `probeAccept`

The point is to ask "does PG's parser get past the lead token?", not "does the statement execute". Semantic errors are still "accept" for our purposes.

| Production | Probe template | Grammar position |
|---|---|---|
| SimpleTypename | `SELECT CAST(NULL AS %s)` | gram.y:14130 (CAST → Typename) |
| ConstTypename | (none — ambiguous with `func_name Sconst`) | n/a |
| Typename | `CREATE FUNCTION __omni_probe() RETURNS %s AS $$ SELECT 1 $$ LANGUAGE sql` | gram.y:8488 (RETURNS func_type, FIRST sets coincide with Typename) |
| func_type | `CREATE FUNCTION __omni_probe(dummy_name %s) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql` | gram.y:8405 (func_arg → param_name func_type) |
| a_expr | `SELECT %s` | gram.y:14780 (target_el) |

### Subtraction lists

Some probes need to subtract tokens that PG accepts via a SURROUNDING grammar rule rather than the production under test. These are documented inline in each test:

- **SimpleTypename via CAST** subtracts `SETOF` — CAST exercises `Typename`, which has SETOF in its FIRST set, but `SimpleTypename` proper does not.
- **func_type via func arg** subtracts `IN_P, OUT_P, INOUT, VARIADIC, DEFAULT, ARRAY` — the surrounding `func_arg` grammar absorbs `arg_class`, the default-value clause, and the array-bounds suffix before they reach `func_type`.
- **a_expr via SELECT** subtracts `ALL` — `SELECT ALL` parses as the duplicate-elimination quantifier from `select_clause`, not as an `a_expr` lead. (DISTINCT is NOT subtracted: bare `SELECT distinct` returns 42601 and is already correctly classified.)

When adding a new probe template, expect to discover a similar list. The procedure is: run the test without subtractions, observe "PG accepts but predicate rejects" drift, verify each token is genuinely absorbed by surrounding grammar (not a real bug), then add it to the subtraction map with an inline comment explaining why.

## Multi-token type rendering (`renderTypeCandidate`)

Several PG type keywords only form valid types in combination with follow-on tokens:
- `DOUBLE` → `DOUBLE PRECISION`
- `NATIONAL` → `NATIONAL CHARACTER` / `NATIONAL CHARACTER VARYING(n)`
- `BIT` → `BIT(n)` / `BIT VARYING(n)` (bare `BIT` parses with default length)
- `CHARACTER` → `CHARACTER(n)` / `CHARACTER VARYING(n)`

`renderTypeCandidate` in `first_set_oracle_test.go` enumerates all plausible forms (bare, `(1)`, `(1,1)`, `precision`, `character`, `character(1)`, `varying`, `varying(1)`, `character varying`, `character varying(1)`) and the runner accepts the candidate if ANY rendering parses. `SETOF` has its own special case (`setof int`, `setof varchar(10)`, `setof json`).

If you add a new multi-word type lead, extend `renderTypeCandidate` accordingly. The convention is:
- Bare-token types stay in the default 10-form list.
- Tokens that REQUIRE a continuation (like SETOF) get a dedicated branch.

## Expression continuations (`renderExpression`)

`a_expr` has many keyword leads that are syntax errors bare:
- `NOT` → needs an operand: `NOT TRUE`
- `EXISTS` → needs a subquery: `EXISTS (SELECT 1)`
- `CASE` → needs WHEN clauses: `CASE WHEN TRUE THEN 1 END`
- `CAST` → needs a target type: `CAST(NULL AS int)`
- `ARRAY` → needs `[...]`: `ARRAY[1]`
- `ROW` → needs a tuple: `ROW(1)`
- `NULLIF` / `COALESCE` / `GREATEST` / `LEAST` → need argument lists
- `INTERVAL` → needs a string literal: `INTERVAL '1 day'`

These are documented in `renderExpression`'s expansion map. **`COLLATE` is intentionally absent** — it is a postfix operator (`a_expr COLLATE any_name`), not an expression starter. Adding it would produce a misleading test where the rendering doesn't actually start with COLLATE.

## Why no PG oracle for ConstTypename

PG's grammar for `AexprConst` admits BOTH `ConstTypename Sconst` AND `func_name Sconst` at the same position. `func_name` starts on any `IDENT`, `UnreservedKeyword`, `TypeFuncNameKeyword`, or `ColNameKeyword`, so `SELECT between '1'`, `SELECT text '1'`, `SELECT foo '1'` all parse via the `func_name` branch — regardless of whether the lead is a real ConstTypename token. PG's accept signal at this position is therefore not specific to ConstTypename.

Solution: in-Go bidirectional unit test. `isConstTypenameStart` reuses `simpleTypenameLeadSet` (the two FIRST sets coincide in omni), and `TestIsConstTypenameStartRejectsAllOtherKeywords` iterates the entire `Keywords` table asserting the predicate matches the slice exactly. The hand-curated `TestIsConstTypenameStartRejectsNonTypeStarters` is kept alongside as readable documentation of the negative cases.

## Why no PG oracle for TableConstraint

`TableConstraint` is only reachable through `CREATE TABLE` and `ALTER TABLE` column-definition rules. The surrounding grammar admits column names (which can be any IDENT) at the same position as constraint leads, so PG's accept signal isn't specific to the constraint production. Coverage is provided by:
- `TestIsTableConstraintStartCoverage` — positive (6 lead tokens) and negative (IDENT, INT_P, SELECT, LIKE) checks
- The existing pg/parser parse tests for `CREATE TABLE` / `ALTER TABLE` that exercise the migrated call sites

## Adding a new production predicate

When you discover a new FIRST set that needs a single source of truth (typically when Phase 3 audit-style work surfaces a new probe site or a fresh production lands in PG):

1. **Add the lead-token slice** to `first_sets.go`. Include a doc comment with: a `gram.y` reference, a "DO NOT extend without ..." invariant, and a one-line justification for any subset/superset relationship to other slices.
2. **Add the predicate** as `isXStart`. Compose from existing predicates where possible (e.g., `isFuncTypeStart` is a thin alias for `isTypenameStart`). Avoid creating a new map if an existing one already represents the right set.
3. **Add the oracle test** in `first_set_oracle_test.go`. Choose a probe template that exercises the production at the cleanest possible grammar position. Run the test, observe drift, decide whether to fix the slice, fix the renderer, or add a subtraction (with documented justification).
4. **Migrate call sites** that previously hand-wrote the cluster. Use grep to find them. **Do not mass-replace** — read each call site, confirm the substitution is behavior-preserving, and adjust the surrounding control flow if necessary.
5. **Update this CLAUDE.md** — add the new predicate to the canonical list and the new probe template to the table.

If the production cannot be oracle-tested (grammar ambiguity), use the in-Go bidirectional pattern from `TestIsConstTypenameStartRejectsAllOtherKeywords`: iterate the full `Keywords` table and assert the predicate's accept set matches the slice's membership exactly.

## Probe template authoring rules

Each probe template MUST:

1. **Cite the grammar production** it exercises in a doc comment (file:line in gram.y is fine).
2. **Justify why the surrounding production doesn't pollute the FIRST-set measurement.** If the surrounding rule has a broader FIRST set, document the subtractions needed.
3. **Use the lowest-overhead SQL** that exercises the position. Prefer `SELECT %s` over `CREATE FUNCTION ...` when both work.
4. **Not depend on a specific schema or session state.** Probes run inside a transaction that gets rolled back, so they should be self-contained.

## CI vs local execution

Oracle tests use `testcontainers-go` to spin up `postgres:17-alpine`. The helper distinguishes:

- **Local dev**: if Docker is unavailable (or testcontainers-go can't reach it), the test calls `t.Skipf` and the run reports SKIP.
- **CI**: if `$CI` is set to `true` or `1`, the test calls `t.Fatalf` instead. This prevents the silent-disable failure mode where a misconfigured CI runner skips every FIRST-set test and hides drift.

The CI marker is the standard `CI` env var set by GitHub Actions, CircleCI, GitLab CI, etc. If omni adopts a project-specific marker, update `isCI()` in `first_set_oracle_test.go`.

### Known operational quirk: Ryuk reaper race

When tests are re-run rapidly in the same shell session (e.g., during local iteration), the testcontainer Ryuk reaper sometimes terminates the previous container while a new test is initializing, producing a `connection refused` ping failure. **This is a flake, not a test bug.** Wait a second and retry; or set `TESTCONTAINERS_RYUK_DISABLED=true` in the local environment if it recurs.

CI runners with isolated processes per job should not be affected, since Ryuk's tracked container lifetime aligns with the test process.

## PG version upgrade flow

The oracle pins `postgres:17-alpine` to match omni's pg parser grammar target (see `pg/parser/type.go:9` and similar references throughout the package). When omni eventually upgrades its parser target to PG 18 or later:

1. Bump the image tag in `startFirstSetOracle` (`first_set_oracle_test.go`).
2. Run the oracle tests. Expect drift on any production where PG added new lead tokens.
3. For each drift entry: read the relevant `gram.y` change in the upstream PG release notes, then either (a) add the new token to the omni slice if omni's parser handles it, or (b) leave it on the gap list if omni doesn't yet support that PG-N+1 feature.
4. Update this CLAUDE.md's "PG version pinning" note.

The oracle tests are the upgrade signal. **Do not** try to anticipate PG grammar changes without running the tests against the new image.

## PR review checklist

When reviewing a pg/parser PR, ask:

- Does any new file or function add a `case TOK_A, TOK_B, TOK_C:` cluster? If yes, is it a FIRST-set probe? If yes, does it route through `isXStart` or define a new `isXStart`?
- If the PR adds a new `isXStart`, is there a corresponding oracle test (or in-Go bidirectional test for ambiguous productions)?
- If the PR modifies `parseSimpleTypename`'s case list (or any other dispatch switch with a paired predicate), does the corresponding `xxxLeadTokens` slice change in the same commit?
- Does the PR update `renderTypeCandidate` or `renderExpression` if it adds a multi-word type or a continuation-requiring keyword?
- Does the PR's commit message describe the FIRST set being touched, or just say "fix parser bug"?

PRs that introduce hand-written multi-token clusters in probe positions without an `isXStart` indirection should be sent back for refactor — the whole point of the discipline is that the next reader doesn't have to remember it.

## File map (for quick navigation)

**FIRST-set core:**
- `first_sets.go` — all 6 canonical FIRST-set slices and predicates
- `first_set_oracle_test.go` — testcontainer infrastructure, candidate enumeration, probe driver, bidirectional assertion helper, and all FIRST-set tests
- `create_function_json_test.go` — JSON / DEC / NATIONAL CHARACTER / NCHAR regression test with AST-shape assertions

**Reference data:**
- `keywords.go` — the source-of-truth `Keywords []Keyword` table mapping keyword strings to token IDs and category
- `tokens.go` — the omni-side token constants (IDENT, ICONST, etc.)
- `name.go` — `isColId`, `isColLabel`, `isTypeFunctionName` (category-level predicates that the FIRST-set predicates compose with)

**Migrated call sites (FIRST-set probes use `isXStart` here):**
- `create_function.go` — `parseFuncArg` uses `isFuncTypeStart` (canonical example)
- `create_table.go` — `parseTableElement` and `parseTypedTableElement` use `isTableConstraintStart`
- `alter_table.go` — `parseAlterTableAdd` uses `isTableConstraintStart`
- `parser.go` — `isCreateTableElement` delegates to `isTableConstraintStart`
- `copy.go` — `parsePreparableStmt` uses `isSelectStart`
- `publication.go` — `parseRuleActionStmt` uses `isSelectStart` with WITH preserved
- `expr.go` — `isAExprConstTypeCast` is a one-line wrapper around `isConstTypenameStart`; also defines `isSelectStart`

**Dispatch switches (lockstep with FIRST-set slices):**
- `type.go` — `parseSimpleTypename` and friends (the dispatch switch that `isSimpleTypenameStart` mirrors)

**Documentation:**
- `docs/first-set-audit.md` — the Phase 3.1 audit output (one-time snapshot; re-run if a future audit is needed)
- `docs/plans/2026-04-14-pg-first-sets.md` (in repo root) — original implementation plan with full rationale and grammar references

## Backtracking discipline

omni's pg parser is hand-written recursive descent. When a grammar position has local ambiguity that requires lookahead, the canonical approaches in pg/parser are:

1. **Peek-then-commit using `peekNext()`** (preferred when the lookahead is a single token). Examples: `parseFuncArg`, `parseStmt`'s CREATE dispatch, `_LA` token reclassification in `advance()`.
2. **`snapshotTokenStream` + `restoreTokenStream`** from `backtrack.go` (when peek isn't enough — e.g. the speculative parse needs to walk multiple tokens before deciding). Example: `parseFuncType`'s `%TYPE` branch.
3. **`Token` struct snapshot + manual `nextBuf` push** (when the speculative parse only consumes 1 token and rolling back via `nextBuf` push is more compact). Example: `parseCreateStmt`'s CTAS detection at `parser.go:894-918`.

**Forbidden:** synthesizing a "rolled back" token by hand with a hardcoded `Type` field. Two prior bugs were caused by this anti-pattern (the deleted `pushBack(string)` in `create_function.go` and the partial state save in `parseFuncType`'s speculative branch). If you need to roll back, use `snapshotTokenStream` / `restoreTokenStream` — never reconstruct a token from a string.

The `tokenStreamState` snapshot covers cur, prev, nextBuf, hasNext, and the lexer's token-boundary state (Err, pos, start, state). It does **not** cover mid-token lexer state (literalbuf, dolqstart, utf16FirstPart, xcdepth, stateBeforeStrStop) or completion-mode state. If a future caller needs to roll back from inside a string literal or dollar-quoted block, extend `tokenStreamState` carefully and update its doc comment.

## Known limitations / follow-ups

### Three-component qualified type names (separate `parseGenericType` limitation)

`db.schema.mytype` (3-component qualified names) fail to parse in **all** type contexts — `CREATE TABLE`, `CAST`, `ALTER TABLE`, `CREATE FUNCTION`, etc. This is **not** the backtracking bug fixed by this refactor; `parseGenericType` only consumes at most one dot. To accept 3-component names, `parseGenericType` needs to be extended to consume more dots, and the AST needs to handle the additional name component.

Tracked as a separate follow-up — no client has reported it yet.

### Ryuk reaper flake (testcontainer infrastructure quirk)

See "Known operational quirk: Ryuk reaper race" above. Not a parser bug; a testcontainers-go race condition when tests are re-run rapidly in the same shell session.
