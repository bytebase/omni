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
