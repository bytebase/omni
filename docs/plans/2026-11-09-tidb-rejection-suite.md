# TiDB rejection suite (§3.6)

Date: 2026-11-09. Tracks §7 step 4 of
`docs/specs/2026-11-06-tidb-drop-pingcap-design.md`: proving omni rejects what
TiDB v8.5.5 rejects. The invariant is zero over-acceptances (omni accepts, TiDB
rejects); each is fixed across `mysql`, `mariadb`, `oceanbase` in the same PR,
like a §3.3 gap. omni rejecting what TiDB accepts is a §3.3 gap, out of scope
here — noted, never fixed by widening `Parse`.

## Deliverables and PR order

- **(a) Syntax-mapping ledger (lint/audit, no gate on coverage).**
  `tidb/parser/cmd/refaudit` + its test. Validates every `// Ref: ... parser.y:N`
  citation resolves to a real v8.5.5 rule (gate), and reports uncited rules
  (guide). **Done.**
- **(b) Grammar-guided negatives (gate).** Extend `grammarcorpus` to derive, per
  alternative, near-misses that break the production in one place; emit to
  `testdata/parser_y_negatives_corpus.tsv` with omni/engine verdict columns, an
  oracle test with `-update`, measured against live TiDB. Expected reject/reject.
  **Done (baseline landed).**
- **(c) Close over-acceptances (b) surfaces**, one cluster per PR, across forks.
- **(d) Differential mutation fuzz sweep (gate).** Mutate §3.3 corpus and §3.4
  samples at token/byte level; `TIDB_REJECT_FUZZ_N`-sized, seeded, wired into CI
  and AGENTS.md. A failing seed is minimised and added to (b)'s TSV as a
  regression row.
- **(e) Close over-acceptances (d) surfaces.**

## Log

### (a) ledger — PR pending

- `refaudit` resolves 124 line-anchored citations against 690 v8.5.5 rules.
- Found 6 broken citations (3 distinct lines), each off by exactly 437 lines —
  written against master, not the vendored tag. Corrected in `ast/parsenodes.go`
  and `parser/resource_control.go` to the real rule heads: CalibrateResourceStmt
  16445, AddQueryWatchStmt 16543, DropQueryWatchStmt 16631.
- Baseline coverage: 209/690 rules cited. Coarse by design (per-group citations;
  option sub-rules handled inline). Uncited list is a lead set for (b)/grammar
  work, not a gap list.
- Codex review (PR #755) caught that name-only references ("Ref: AnalyzeOption
  in ... parser.y." with no line number) were dropped before coverage crediting,
  so cited rules read as uncited. Fixed: rule names on any Ref line now credit
  coverage independently of numeric ranges (114 → 209 cited).
- Codex review 2 (PR #755) caught that multi-line Ref blocks were scanned a line
  at a time, so a parser.y: range or anchor on a continuation line was dropped —
  both a missed stale citation and under-credited coverage. Fixed: whole Ref
  comment blocks are read as one unit (combine + offset→line map), with range
  lists separated by "," or "and". This surfaced 4 more stale master-based
  citations in resource_control.go (16882/16980/17068 at the file-level doc
  block, 16938 for DURATION); corrected to 16445/16543/16631 and 16501. Now 144
  citations, 290 cited, 0 broken.

### (b) grammar-guided negatives — PR pending

- `grammarcorpus -mode negatives` mutates each reachable alternative's shortest
  valid sentence: drop each body token, truncate, insert a stray identifier,
  swap the first keyword for its neighbour. Deduplicated by statement; mutations
  that reproduce a valid sentence are dropped. 8151 rows.
- Baseline measured live against TiDB v8.5.5 (0 no-verdict): 7795 both-reject,
  153 both-accept (valid near-misses), **201 over-acceptances** (omni accepts,
  TiDB rejects), 2 gaps (omni rejects, TiDB accepts).
- The 201 over-acceptances are the PR (c) work list. They cluster by shared
  root-cause sub-productions (alias-after-`*`; trailing stray token; dropped
  required keyword such as `IF NOT EXISTS` → `IF EXISTS`, `RENAME USER ... TO`,
  `PRIMARY KEY` → `PRIMARY`), not one-per-rule, so a cluster fix flips many rows.
- The 2 gaps (out of scope §3.6, for §3.3): `DELETE FROM t PARTITION(p) zzz` and
  `ALTER DATABASE COLLATE 's' zzz` — TiDB accepts a trailing identifier omni
  rejects. Noted, not fixed; do not widen `Parse`.

### (c) close over-acceptances — cluster 1: wildcard-field alias

- **Cluster: an alias after a `*` / `tbl.*` select field.** omni's
  `parseSelectExpr` parsed an alias after any expression, but in the grammar the
  bare `*`, `tbl.*` and `db.tbl.*` are their own alias-less `Field` alternatives
  (parser.y:6464/6468/6473); only `Expression FieldAsNameOpt` (6478) takes one.
  TiDB answers 1064 for `SELECT * x`, `SELECT t.* x`, `SELECT * AS x`.
- Fixed in all four engines (tidb, mysql, mariadb, oceanbase) — identical early
  return for a StarExpr or a `ColumnRef{Star:true}` before alias parsing, each
  cited to its own grammar (TiDB parser.y; MySQL sql_yacc.yy select_item /
  table_wild for the forks).
- Measured: negatives corpus over-acceptances **201 → 137** (64 rows flipped to
  reject/reject, live-confirmed on TiDB v8.5.5, 0 no-verdict). No new gaps
  (both-accept unchanged at 153); tidb positive corpora show no regression.
- Remaining 137 over-acceptances are later (c) clusters (trailing stray tokens in
  other productions, dropped required keywords, etc.).

### (c) close over-acceptances — cluster 2: foreign-key reference action

- **Cluster: the ON DELETE / ON UPDATE reference action.** `parseReferenceAction`
  returned `RefActNone` with no error for a bare `SET` (no NULL/DEFAULT), a bare
  `NO` (no ACTION), or any non-action token, so omni accepted `ON DELETE SET`,
  `ON DELETE NO`, `ON DELETE NULL`, `ON DELETE` with nothing, and `ON DELETE ON
  UPDATE …`. ReferOpt is always complete and required after ON DELETE/ON UPDATE
  (parser.y:4017-4037, 3983-3990).
- Fixed in all four engines: `parseReferenceAction` now returns an error for an
  incomplete/missing action; both FK caller sites (inline column REFERENCES and
  table-level FOREIGN KEY) propagate it. Cited to TiDB parser.y / MySQL
  sql_yacc.yy `delete_option`.
- Measured live: over-acceptances **137 → 128** (9 rows to reject/reject, 0
  no-verdict), no new gaps (both-accept 153), positive corpora unregressed.

### (c) close over-acceptances — cluster 3: CREATE/ALTER USER password options

- **Cluster: PasswordOrLockOption operand requirements.** `parseUserAccountOptions`
  accepted incomplete options: `PASSWORD HISTORY` (no n/DEFAULT), `PASSWORD REUSE`
  without INTERVAL / without DAY, `PASSWORD REQUIRE DEFAULT` / `REQUIRE CURRENT`
  (TiDB's only form is CURRENT DEFAULT), and bare `ACCOUNT`.
- Divergence found: `mysql` and `mariadb` were **already strict** here; only
  `tidb` and `oceanbase` were loose. Fixed those two. tidb uses TiDB's stricter
  REQUIRE (CURRENT DEFAULT, measured on live TiDB); oceanbase matches the
  already-correct mysql/mariadb (REQUIRE CURRENT [DEFAULT|OPTIONAL]). Measured a
  throwaway set of password forms on live TiDB first to pin each verdict.
- Cited to TiDB parser.y:13993-14082 / MySQL sql_yacc.yy password_option.
- Measured live: over-acceptances **128 → 118** (10 rows to reject/reject, 0
  no-verdict), no new gaps, positive corpora unregressed; oceanbase grant/user
  container tests pass. (SET PASSWORD `=` / `PASSWORD()` forms are a separate
  later cluster.)

### (c) close over-acceptances — cluster 4: ALTER INSTANCE

- **Cluster: ALTER INSTANCE.** `parseAlterInstanceStmt` collected arbitrary
  action words and blindly consumed the NO ROLLBACK ON ERROR tokens, so omni
  accepted `ALTER INSTANCE TLS`, `RELOAD`, `REMOVE TLS`, `RELOAD TLS zzz`, every
  partial NO-ROLLBACK phrase, and even MySQL-only actions.
- TiDB-only: TiDB's InstanceOption is just `RELOAD TLS` and `RELOAD TLS NO
  ROLLBACK ON ERROR` (parser.y:13758-13784); the forks' richer MySQL ALTER
  INSTANCE grammar makes the other actions valid there, so this fix is tidb-only.
  Rewrote the parser to accept exactly those two forms.
- Measured live: over-acceptances **118 → 106** (12 rows to reject/reject, 0
  no-verdict), no new gaps, positive corpora unregressed.
- Note: the forks' loose word-collector also over-accepts trailing junk relative
  to MySQL's grammar (`RELOAD TLS zzz`, malformed NO clause); fixing them to
  MySQL's full instance_action is separate, measured against MySQL.

### (c) close over-acceptances — cluster 5: ALTER TABLE spec required keywords

- **Cluster: ALTER TABLE spec trailing keywords.** The alter-spec parser used
  `p.match` (optional) where the grammar requires a keyword, so omni accepted
  `ALTER TABLE t ENABLE` / `DISABLE` (no KEYS), `WITH` / `WITHOUT` (no VALIDATION,
  incl. EXCHANGE PARTITION … WITH/WITHOUT), and `ORDER c` (no BY).
- Shared across all four engines (MySQL-family). Changed each `p.match` to
  `p.expect` for KEYS, VALIDATION, and BY in tidb/mysql/mariadb/oceanbase.
- Measured live: over-acceptances **106 → 99** (7 rows to reject/reject, 0
  no-verdict), no new gaps, positive corpora unregressed; broad parse unit tests
  pass on all four engines.

### (c) close over-acceptances — cluster 6: ALTER TABLE IMPORT/DISCARD TABLESPACE

- **Cluster: IMPORT/DISCARD TABLESPACE.** The parser used `p.match(kwTABLESPACE)`
  (optional) where the grammar requires it, so omni accepted `ALTER TABLE t
  IMPORT` / `DISCARD` (no TABLESPACE), `IMPORT/DISCARD PARTITION ALL` (no
  TABLESPACE), and `IMPORT/DISCARD PARTITION TABLESPACE` (no name-list). Required
  TABLESPACE with `p.expect` at all four spots in each engine.
- Shared across all four engines. Measured live: over-acceptances **99 → 93**
  (6 rows to reject/reject, 0 no-verdict), no new gaps, positives unregressed;
  broad parse unit tests pass on all four.

### (d) differential mutation fuzz sweep — harness + minimiser (CI wiring pending)

- `tidb/parser/reject_fuzz_test.go` (`TestTiDBRejectFuzz`): mutates the valid
  (accept/accept) statements of the §3.3 corpus and §3.4 samples at token and
  byte level, runs omni + live TiDB v8.5.5, flags over-acceptances. Mirrors
  pg's TestParenOracleFuzz: fixed seed (`rejectFuzzSeed`), size by
  `TIDB_REJECT_FUZZ_N`, known-mismatches allowlist with set-equality gate,
  `TIDB_REJECT_FUZZ_DEFER` log-and-persist mode, token-level ddmin minimiser.
- Validated at N=300: found 2 over-acceptances the grammar-guided corpus missed,
  e.g. `select ( * )` (omni accepts a parenthesized `*`; TiDB rejects). Baselined
  in `testdata/reject-fuzz-corpus/known-mismatches.txt`.
- **CI budget pending orchestrator confirmation** before wiring `.github/workflows/ci.yml`
  and the AGENTS.md sweep list: env `TIDB_REJECT_FUZZ_N`, seed `rejectFuzzSeed`,
  proposed CI size + defer mode.
