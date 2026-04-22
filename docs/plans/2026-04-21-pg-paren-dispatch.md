# PG Parser `(` Dispatch Alignment Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Section 0 is done via manual validation on branch `junyi/byt-9315-fix-postgresql-grammar-to-accept-multi-join-syntax-in-sql`; subsequent sections dispatch via the `pg-paren-dispatch` starmap skill (to be created).

**Goal:** Align omni's hand-written recursive-descent PG parser with PG's yacc grammar at every `(` dispatch point, so that each ambiguous `(` (where `(` could start >1 grammar production) is resolved by one of the standard techniques in §3 rather than a 1-token peek that picks the wrong branch.

**North star metric:** `pg/pgregress/known_failures.json` count. Currently 293 (down from 298 after §0). Every section targets -N entries; the starmap is "done" when residual failures are all provably non-syntax issues (missing PG features, catalog-semantic mismatches, etc.) rather than `(` mis-dispatches.

---

## 1. Background

### 1.1 Triggering bug: BYT-9315

A bytebase customer's view had a `pg_get_viewdef()` output of shape `SELECT * FROM ((a JOIN b ON TRUE) JOIN c ON TRUE)`. omni's parser rejected it, which cascaded through `catalog.Exec` and killed SQL review for every statement on the instance (originally surfaced as BYT-9308, root cause filed as BYT-9315).

### 1.2 Root cause: eager peek-dispatch at `(`

omni's `parseParenTableRef` (select.go, pre-fix):

```go
next := p.peekNext()
if next.Type == SELECT || VALUES || WITH || TABLE || next.Type == '(' {
    // Treat as subquery
} else {
    // Treat as joined_table wrap
}
```

The 1-token peek cannot distinguish:
- `((SELECT ...))` — nested select_with_parens
- `((a JOIN b))` — nested joined_table
- `((SELECT 1) JOIN ...)` — outer joined_table, subquery is first operand

All three start with `((`, so the peek routes uniformly to one branch and silently mishandles the other two cases.

### 1.3 Architectural root cause: LALR vs. recursive descent

PG's yacc grammar (`postgres/src/backend/parser/gram.y:13537-13545`) explicitly notes this ambiguity:

> "It may seem silly to separate joined_table from table_ref, but there is method in SQL's madness: if you don't do it this way you get reduce-reduce conflicts, because it's not clear to the parser generator whether to expect alias_clause after `')'` or not."

PG resolves `(` ambiguities at **two layers**:
1. **Grammar structure.** The `select_with_parens` / `select_no_parens` split (`gram.y:12682-12759`) and the `joined_table` / `table_ref` split (`gram.y:13446-13534`, `13554-13600`) are themselves the disambiguation: different productions live in different nonterminals so LALR doesn't have to decide at `(`.
2. **LALR(1) stack lookahead.** Where the grammar split isn't enough (e.g. `'(' joined_table ')'` with optional alias after `)`), LALR(1) shifts both alternatives and defers reduction until a later token discriminates.

omni's recursive descent has no parse stack holding multiple alternatives. At every `(`, a single function must commit to one production. Some ambiguities disappear once omni mirrors PG's grammar-structure split (layer 1 → technique T6 below). The rest need bounded lookahead or speculative parse (layer 2 → T3/T4). The bugs this plan addresses are where omni took a shortcut — sharing one parser across multiple PG nonterminals, or a 1-token peek where the grammar itself demands more.

---

## 2. What "aligned with PG" means concretely

For each grammar nonterminal N that has productions starting with `(`:

1. There is exactly one omni function responsible for parsing N.
2. That function matches PG's set of productions for N and nothing else:
   - **Not stricter:** rejects only inputs PG rejects.
   - **Not looser:** accepts only inputs PG accepts.
3. The disambiguation technique used is recorded in the function's doc comment (see §3).
4. The function is covered by:
   - Handwritten tests pinning the accept/reject boundary.
   - `pg/pgregress` corpus (indirect — fixing this function reduces known_failures count).
   - Where feasible, an oracle test against a real PG container (see first-set consolidation plan 2026-04-14-pg-first-sets.md for the pattern).

---

## 3. Disambiguation technique catalogue

omni's existing toolkit, in order of preference. Each section in §5 must declare which technique it uses and why.

| # | Technique | When to use | omni primitive |
|---|---|---|---|
| T1 | **Single peekNext()** | Disambiguation lives in 1 token beyond cur | `p.peekNext()` in parser.go |
| T2 | **Lexer-level `_LA` reclassification** | Disambiguation is a fixed 2-token combination that PG itself merges (e.g. `WITH_LA`, `NULLS_LA`, `NOT_LA`) | `advance()` reclassification in parser.go |
| T3 | **Snapshot + token-stream scan** | Disambiguation is "is there a keyword X inside this balanced paren block" — bounded lookahead, no partial AST construction | `snapshotTokenStream` / `restoreTokenStream` in backtrack.go |
| T4 | **Snapshot + speculative parse** | Disambiguation needs to actually run a parse production, and rollback is cheaper than building an ambiguity-aware grammar | Same primitive as T3, but the speculative block calls parse\* helpers |
| T5 | **Left factoring** | The ambiguous productions share a prefix; consume the prefix unconditionally and let downstream (set-op / alias / join) disambiguate | Hand refactor: unify the entry-point, delegate to the shared downstream |
| T6 | **Dedicated nonterminal function** | omni currently shares a generic parser across multiple PG nonterminals (e.g. `parseTableRef` used where only `joined_table` is valid). Introduce a dedicated function that only accepts the target nonterminal's productions | Hand refactor: new `parseJoinedTable`-style function |
| T7 | **Parse-then-validate** | Accept the over-broad production, then assert invariants after parse. Use when grammar refactor is disproportionate to the leniency being closed | Post-parse type assertion + syntax error |
| T8 | **Centralized FIRST-set predicate** | The ambiguity is really "which production does this lead token start?" and that decision is duplicated in multiple probe sites. Replace with one `isXStart()` predicate backed by one token slice (see 2026-04-14-pg-first-sets.md) | `is*Start()` helpers in `first_sets.go` |

Not used in omni (for reference): PEG/packrat, ANTLR LL(*), yacc GLR, parser generator rewrite. These would replace the parser rather than align it.

**Precedence** (revised after Codex review):
1. **T5/T6 first** when the ambiguity maps cleanly to a PG nonterminal split. PG solved this at the grammar-structure level; mirroring that is the durable fix. A scan-based probe (T3) is tactical and accumulates bespoke helpers over time.
2. **T1/T2** for ambiguities that PG itself resolves with a fixed N-token combination or `_LA` reclassification.
3. **T8** when the ambiguity is a FIRST-set question duplicated across sites. Coordinate with the pg-first-sets starmap.
4. **T3** when the decision needs bounded intra-paren scan and no clean nonterminal split exists.
5. **T4** only when T3 can't encode the decision (the speculative parse needs to actually construct AST to decide).
6. **T7** as last resort when the grammar refactor is disproportionate. Document as a hybrid (T7 behind a T6 boundary) rather than pretending it's pure structural enforcement. **Phase 0's `parseJoinedTable` is T6/T7 hybrid** — the outer function is a dedicated nonterminal (T6) but internally delegates to `parseTableRef` and asserts `*JoinExpr` (T7). Hardening to pure T6 (refuse to return without consuming JOIN) is tracked as a Phase 2 follow-up.

---

## 4. Scope

### 4.1 In scope (this starmap: `pg-paren-dispatch`)

1. **Every `(` dispatch site** where `(` starts ≥2 PG productions and omni's code makes the decision with ≤1 token of lookahead. Preliminary count: 56 `p.cur.Type == '('` occurrences across 18 files.
2. **Every `)` decision site** — `)` ambiguity is the symmetric half of `(` ambiguity (PG's comment at `gram.y:13538-13544` explicitly hinges on "alias_clause after `)` or not"). The audit catalogs `)` dispatch alongside `(`.
3. **Shared-helper reach.** If a fix requires splitting a shared helper (like the Phase 0 `parseTableRef` → `parseJoinedTable` split), the split is in scope even though the visible callsite was the `(`.

### 4.2 Related but explicitly delegated

The following overlap this starmap's domain but are owned by other work streams. Phase 1 audit rows that depend on these are **flagged blocked** rather than deferred silently:

| Concern | Owner | Handshake |
|---|---|---|
| Nonterminal boundary mismatches outside `(` / `)` | follow-up starmap `pg-nonterminal-alignment` (not yet started) | Phase 1 rows that need caller-context changes outside the `(` site flag `blocked_by: pg-nonterminal-alignment` |
| FIRST-set duplication | 2026-04-14-pg-first-sets.md | T8 scenarios coordinate with that plan; predicate moves during its Phase ≥2 may invalidate "aligned" rows here |
| Dual-return migration | pg-dual-return starmap | Parse functions we edit must land in the already-migrated (Node, error) form |
| Soft-fail cleanup | pg-soft-fail starmap | Any new error we emit from a `(` dispatch site must follow the soft-fail rules (advance-then-parse guard) |
| Error message wording | pg-error-align starmap | We emit generic "syntax error at or near" from new reject paths; wording polish is deferred to pg-error-align |
| Transform-layer semantic checks | analyze.go + downstream | Out of scope entirely |

### 4.3 Non-goals

- Becoming byte-for-byte AST-identical with PG. Different node structures are fine where omni has deliberate deviations documented in `pg/ast/`.
- Covering SQL that PG rejects. If PG rejects it, we reject it; we don't accept something stricter than PG.
- Rewriting the parser to use a generator. The alignment is done **in-place** on the recursive-descent code.

---

## 5. Phases

### 5.0 Phase 0: BYT-9315 core + Sections 1.1 / 1.2 ✅ DONE (manual validation)

Delivered on branch `junyi/byt-9315-fix-postgresql-grammar-to-accept-multi-join-syntax-in-sql`:

| Section | Change | Technique | `pgregress` delta |
|---|---|---|---|
| 0.0 | `parenBeginsSubquery` in `select.go` — scans matched paren block, classifies outer as subquery vs joined_table based on depth-1 contents and set-op / join continuation after nested close | T3 | -12 |
| 1.1 | `parseJoinedTable` nonterminal — separates "parse joined_table inside `'(' ... ')'`" from generic `parseTableRef`, enforces JoinExpr root | T6/T7 hybrid (structural nonterminal boundary + post-parse `*JoinExpr` assertion; pure-T6 hardening tracked as Phase 2 follow-up) | -3 (with 1.2) |
| 1.2 | `parseSelectWithParens` left-factored — removed `cur == '(' { recurse }` short-circuit, always delegates to `parseSelectNoParens`, letting select_clause set-op precedence handle nested paren operand | T5 | -2 (with 1.1) |

Net: **-17 pgregress known_failures**, 0 new failures, 4 previously unaligned FROM-clause edge cases now match PG:

```
( (a JOIN b) JOIN c )          — was fail, now pass
(a)                            — was accept (PG reject), now reject
((a))                          — was accept (post 0.0 regression), now reject
((SELECT 1) x)                 — was accept (post 0.0 regression), now reject
((SELECT 1) UNION (SELECT 2))  — was reject, now pass
```

**Tests landed:** `pg/parser/paren_multi_join_test.go` (4 test functions, 18 cases).

### 5.1 Phase 1: Cluster-based audit + fix (revised after Codex review)

The initial plan had "Phase 1 = audit everything, Phase 2 = fix". Codex flagged the failure mode: a helper-level fix (like Phase 0's `parseTableRef` split) invalidates earlier classifications and mis-attributes pgregress impact. The sister `pg-first-sets` plan phases by **production family** instead of "inventory-first" for this reason.

Revised model: **audit by cluster, fix within cluster, move to next cluster**.

A cluster is one PG nonterminal family whose `(` / `)` dispatches are coupled (share helpers, share grammar comment, likely to move together). Preliminary clusters:

| Cluster | Files | Nonterminals | pgregress suspects |
|---|---|---|---|
| C1: table_ref / joined_table / select_with_parens | `select.go` | `table_ref`, `joined_table`, `select_with_parens`, `select_no_parens` | High — Phase 0 already closed 17 failures here |
| C2: expression parens | `expr.go` | `a_expr`, `b_expr`, row constructor, subquery in expr, `IN (...)` | Medium |
| C3: type parens | `type.go` | `SimpleTypename` parens, `Typename` modifiers, `func_type` | Medium (overlaps pg-first-sets) |
| C4: DDL element lists | `create_table.go`, `create_index.go`, `define.go` | `TableElementList`, `IndexElem` parens | Low (mostly T1 sites) |
| C5: utility statements (tracked as **unbounded tail**, see note) | `utility.go`, `maintenance.go`, `copy.go`, `grant.go`, `extension.go`, `fdw.go`, `publication.go`, `prepare.go`, `schema.go`, `database.go`, `insert.go` | Per-command (heterogeneous) | Low, long tail |

**Note on C5:** Unlike C1-C4 which share grammar-family coupling, C5 is a catch-all of heterogeneous per-command dispatch sites. It will not close as a unit. Treat C5 as an **unbounded tail**: each file gets its own sub-cluster (C5.a `utility.go`, C5.b `copy.go`, …). PAREN_PROGRESS.json tracks them as independent entries; Phase 3's "starmap done" bar requires every C5.* sub-cluster closed or flagged blocked, not C5 as a whole.

Per cluster, one worker:
1. **Audit sub-step (C-specific rows):** populate PAREN_AUDIT.md with rows only for sites in this cluster.
2. **Fix sub-step:** implement section(s) in SCENARIOS-pg-paren-dispatch.md for that cluster.
3. **Close sub-step:** regenerate `pgregress -update`, report `fixed=K, new=0`, update `PAREN_PROGRESS.json`.

Cluster ordering: C1 (highest pgregress density, already validated in Phase 0) → C2 → C3 (coordinate with pg-first-sets) → C4 → C5.

**Exit criterion per cluster:** every site in the cluster has an audit row with either `aligned = yes` + 2 proofs (§5.3) or a scenario number tracking the fix.

### 5.2 Phase 2: Hardening pass on `parenBeginsSubquery`

The Phase 0 `parenBeginsSubquery` / `consumeMatchedParenIsSubquery` scanner (`select.go`) is a bespoke recursive primitive that 3 of the 4 BYT-9315 cases depend on. Current coverage: 18 handwritten cases (`paren_multi_join_test.go`). That's enough evidence for the fix, but thin for a foundational scanner.

Dedicated hardening pass (runs in parallel with C2+ cluster work):
1. Oracle-style test against PG testcontainer (pattern: see 2026-04-14-pg-first-sets.md §oracle design). Build a minimal corpus of 200+ `FROM (...)` SQL variants, compare omni's routing decision (subquery vs joined_table) against PG's acceptance.
2. Fuzz test: random balanced-paren SQL with interleaved SELECT / JOIN / set-op keywords, compare omni vs PG.
3. If fuzz surfaces a class of mis-routing, either fix `parenBeginsSubquery` or replace it with T5/T6 (the principled preference per §3).

**Delivered scope (post-Phase 2 acknowledgment):** the landed corpus is
N=188 probes — 100 PRNG-generated via `fuzzCorpusSize=100` in
`paren_oracle_fuzz_test.go` (`fuzzSeed=0xBADC0DE1`, deterministic) + 3
active seed entries in `testdata/paren-fuzz-corpus/seed-cases.txt` + 85
hand-curated across §2.2–§2.7 (simple/subquery/joined/mixed/LATERAL/
degenerate). The original "200+ targeting `parenBeginsSubquery`
specifically" was a rough estimate; in practice the oracle harness
covers the whole FROM-clause `(` dispatch surface: `parenBeginsSubquery`
plus LATERAL variants (select_with_parens / XMLTABLE / JSON_TABLE /
func_table / ROWS FROM), VALUES / TABLE / WITH subqueries, set-op
operand paren-wrapping, column-list aliases, and the obvious-reject
perimeter. This is broader than strictly needed for
`parenBeginsSubquery` alone; the extra coverage is kept for
defense-in-depth — it's the single cheapest regression fence for every
Phase-1 fix site (1.1–1.4) that routes through a paren in FROM context.

### 5.3 The "aligned without code change" bar (answers §8 Q3)

A site can be marked `aligned = yes` without changing code only if **both** of these hold:

1. **Caller-context proof.** A written argument that the ambiguous alternative is not actually reachable given how the function is called. References must be to specific callsites and PG grammar nonterminals, not hand-wave.
2. **Empirical check.** At least 5 SQL inputs exercising the nearby grammar surface are tested (accepted-by-PG inputs parse; rejected-by-PG inputs reject). Pinned in a test file.

A markdown assertion alone is below bar.

### 5.4 Phase 3: Close

Starmap is done when:
- Every row in PAREN_AUDIT.md has `aligned = yes` (with the two proofs) OR links to a landed fix section.
- pgregress known_failures: no remaining entries attributable to `(` / `)` mis-dispatch (per-row attribution recorded during audit).
- `parenBeginsSubquery` oracle/fuzz harness is landed and green.
- PAREN_AUDIT.md is committed as the permanent reference.

---

## 6. Execution model

### 6.1 Artifacts

```
docs/plans/
  2026-04-21-pg-paren-dispatch.md       ← this file
pg/parser/
  PAREN_AUDIT.md                        ← audit rows (grows per cluster)
  PAREN_AUDIT.json                      ← machine-readable mirror of AUDIT.md
  SCENARIOS-pg-paren-dispatch.md           ← per-section fix scenarios
  PAREN_PROGRESS.json                   ← cluster/section state + history[]
  paren_*_test.go                       ← per-section tests
  paren_oracle_test.go                  ← Phase 2 PG-container oracle (once landed)
```

PAREN_AUDIT.json schema (one array of row objects) mirrors the markdown; kept in sync by the worker skill. **Canonical schema doc:** `pg/parser/PAREN_AUDIT_SCHEMA.md` — enforced by `TestPARENAuditLint` on every CI run (SCENARIOS §5.3). Live fields: `site` (file:line, stable audit coordinate), `function` (enclosing Go function), `nonterminals` (array), `ambiguity_present` (bool), `current_technique` (T1..T8 or null), `pg_reference` (gram.y:line), `aligned` (enum: yes / no / blocked / unclear), `blocked_by` (nullable — e.g. "pg-nonterminal-alignment", "pg-first-sets"), `cluster` (C1..C5 with optional subcluster suffix), `priority` (high/med/low), `proof_notes` (free-form caller-context + empirical test citations; required non-empty when aligned=yes), `suspicion_notes` (nullable).

### 6.2 Skills (to be created after plan approval)

- `pg-paren-dispatch-driver` — modeled on `pg-error-align` / `pg-migration-driver`. Reads `PAREN_PROGRESS.json`, spawns workers, collects pgregress diffs, updates progress.
- `pg-paren-dispatch-worker` — executes one section: writes/edits code, adds tests, runs `./pg/parser/` + `./pg/pgregress/ -update`, reports "FIXED X, NEW Y".

Neither skill is built until after Codex review of this plan. Phase 1 (audit) is light enough to run manually without the driver.

### 6.3 Review gates

- **Gate A: this document.** Codex reviews the plan for scope, technique selection, and execution model. Any structural concerns fixed here before work starts.
- **Gate B: Phase 1 audit.** Human reviews PAREN_AUDIT.md for completeness and priority ranking.
- **Gate C: first 2 Phase 2 sections.** Human reviews to confirm worker velocity/quality before batching.

### 6.4 Observability

Each worker run reports:
```
cluster: C1  section: 1.3 <name>
files touched: ...
tests added: ...
pgregress before: failed=N known=N new=0
pgregress after:  failed=M known=M new=0 fixed=K
AUDIT rows added: X  rows modified: Y
```
Driver aggregates into `PAREN_PROGRESS.json.history[]`.

### 6.5 CI cost policy

Each section triggers `./pg/...` + `./pg/pgregress/`. Empirically `pgregress` takes ~1s cached, ~10-30s fresh (measured on the BYT-9315 branch on 2026-04-21, darwin-arm64, local `go test` cache warm vs cold — **single-machine estimate**, not CI). Oracle tests (once landed in §5.2) projected at ~2 min for 200 probes based on testcontainer startup + per-probe round-trip (not yet measured). Driver validates these numbers during the first section and updates this document if actuals differ by >2×.

- Worker runs: full `./pg/...` + `./pg/pgregress/` locally. Non-negotiable.
- PR CI: same, plus fuzz corpus (when landed) sampled at 1k inputs.
- Driver may batch multiple small sections into a single PR if they touch disjoint files, to amortize CI. Target: ≤5 section-PRs per day per cluster to keep review latency sane.

### 6.6 Metric: count is rollup, alignment is done-ness

`known_failures.json` delta is reported every section as a **rollup signal** — useful for trend visibility but not the completion bar. The completion bar is:

> Every audit row has `aligned = yes` with the two proofs, or `aligned = blocked` with a named downstream starmap.

Codex (§8 review) correctly flagged that `known_failures.json` biases work toward high-density clusters like `join.sql` and can starve subtle single-line bugs. The count is therefore a floor (never regresses) but not a ceiling (we do subtle sites even if their individual pgregress impact is small).

---

## 7. Risks & mitigations

| Risk | Likelihood | Mitigation |
|---|---|---|
| A fix for site X introduces leniency at site Y (cross-coupling through shared helpers like `parseTableRef`) | Medium | Always run full `./pg/...` + `./pg/pgregress/` after each section. Tests pin the accept/reject boundary, not just the happy path. |
| An "aligned" site regresses when PG grammar changes (future PG version) | Low | Document the PG version in each section. pg-first-sets oracle (separate plan) catches keyword-category drift; `(` dispatch drift would be caught by PG testcontainer regression once landed in §5.2. |
| Scope creep into non-`(` nonterminal boundary work | Medium | Strict §4.2 list + explicit `blocked_by` flag in audit rows. Anything genuinely outside goes into the follow-up starmap `pg-nonterminal-alignment`. |
| Technique T4 (speculative parse) slows common paths | Low | Prefer T5/T6 per revised §3 precedence. T4 only when unavoidable. Benchmark via existing `pg/parser/cache_test.go` if used. |
| pgregress `known_failures.json` churn masks real new failures | Low | Worker reports explicit `new=K` count. Any non-zero `new` blocks the PR. |
| **Completion-mode (IDE autocompletion) regression** — `parenBeginsSubquery` is a snapshot scanner (`select.go`), and `pg/parser/CLAUDE.md:197-207` warns that snapshot-vs-peek affects candidate collection | Medium | Each section adds a `complete_test.go` assertion when the touched function has a `collectMode()` branch. Phase 0 already runs `./pg/completion/` green; regression guard kept. |
| **Mis-attributed `pgregress` deltas** — a single helper refactor can collapse multiple audit rows, making "impact per row" misleading for prioritization | Medium | Priority column in AUDIT rows is informational, not load-bearing. Driver attributes fixed rows to the most specific audit row it can, or to the helper itself with a list. |
| **`parenBeginsSubquery` is a load-bearing new primitive** with thin coverage | Medium | §5.2 dedicated hardening pass (oracle + fuzz) before expanding use beyond Phase 0. Do not promote to general-purpose primitive (Codex §8 Q2). |
| **Other starmaps move the ground under us** (pg-first-sets, pg-dual-return, pg-soft-fail, pg-error-align are active concurrently) | Medium | §4.2 handshake table documents coordination. Driver pins the upstream starmap's commit hash in `PAREN_PROGRESS.json.dependencies[]` at section start; if upstream lands during a section, worker rebases and re-runs pgregress before reporting. |
| **CI saturation** — pgregress per section × N workers | Low-Medium | §6.5 CI policy + section batching when files are disjoint. |

---

## 8. Resolved design questions (Codex review pass, 2026-04-21)

Codex reviewed the v1 draft. The five open questions below are folded back into the plan as stated; this section preserves the resolutions for audit trail.

1. **T7 vs. T6 for `parseJoinedTable`.** Resolution: current implementation is **T6/T7 hybrid**, documented as such in §3 precedence table and §5.0. Pure-T6 hardening (refuse to return without consuming JOIN in the function's own loop) is a Phase 2 follow-up, not a blocker for shipping BYT-9315.
2. **Promote `parenBeginsSubquery` to a general primitive?** Resolution: **No.** Its contract is domain-specific to SELECT/set-op/JOIN routing. Each ambiguity site gets its own bespoke scan where T3 is unavoidable. §5.2 hardening pass is scoped to `parenBeginsSubquery` specifically.
3. **Bar for "aligned without code change".** Resolution: §5.3 requires **two proofs** — caller-context argument + ≥5 pinned empirical tests. Markdown assertion alone is rejected.
4. **Audit output format.** Resolution: both Markdown (PAREN_AUDIT.md, for human review) and JSON (PAREN_AUDIT.json, for tooling). §6.1 documents the schema.
5. **Catalog `)` too?** Resolution: **Yes**, now in-scope per §4.1. The `(` / `)` audit is one document, two axes.

---

## 9. Glossary

- **Ambiguity site:** a `(` in the grammar where PG has ≥2 productions whose FIRST sets both include `(`.
- **Alignment:** accept/reject set matches PG for all inputs; AST shape matches PG for accepted inputs.
- **nonterminal:** a named grammar rule (e.g. `table_ref`, `joined_table`, `select_with_parens`).
- **Left factoring:** refactoring grammar / code so the common prefix of two productions is parsed once, and divergence happens after the common prefix.
- **LALR(1):** lookahead-LR parsing with 1 token of lookahead — what yacc/bison produces. Keeps multiple reduction candidates on the parse stack until a discriminating token appears.
