# TiDB catalog parity quarantine inventory

Date: 2026-11-09. Triage that opens the `tidb/catalog` parity phase, the last
blocker before Bytebase can drop its `github.com/pingcap/tidb` dependency. The
parser rejection suite (§3.6) is closed; this measures the quarantined
differential catalog tests and clusters their failures so fixes can be
parallelized. **Triage only — no fixes in this pass.**

## The measurement

`TIDB_CATALOG_PARITY=1 go test ./tidb/catalog/... -count=1 -v`, one shared TiDB
v8.5.5 container (`startTiDBForCatalog`, cleaned up in `TestMain`). The gate is
`startContainer` in [tidb/catalog/container_test.go](../../tidb/catalog/container_test.go);
every differential test routes through it and is skipped unless the env var is set.

| | count |
|---|---:|
| failing top-level tests | **37** |
| passing | 759 |
| skipped (non-container) | 3 |

**37, not the 93 the quarantine comment records.** The gap narrowed on its own:
the comment predates both the repoint of the shared container from `mysql:8.0`
to TiDB v8.5.5 and the §3.6-era fixes. The `startContainer` comment should be
updated to 37 when the first fix lands.

**The task's "utf8mb4_bin SHOW CREATE ripple" is already fixed in the renderer.**
`ShowCreateTable` ([show.go](../../tidb/catalog/show.go)) emits
`COLLATE=utf8mb4_bin` and the per-column `COLLATE` clauses correctly — on both
sides of every diff. The residual collation failures are at the
`INFORMATION_SCHEMA`/assertion-literal layer, where the tests still hard-code
MySQL's `utf8mb4_0900_ai_ci` / `latin1_swedish_ci`.

**Most failures are MySQL-shaped test assumptions, not omni bugs.** These tests
were forked from `mysql/catalog` with MySQL fixtures and MySQL `want` literals.
A large majority fail because real TiDB v8.5.5 behaves differently from MySQL
(rejects spatial types, lacks session variables, doesn't retain CHECK
constraints, returns `*_bin` collations), and the test still asserts MySQL. A
minority are genuine omni catalog gaps — those are flagged below and are the
higher-risk work.

## Omni-side scope: what this unblocks in dropping pingcap/tidb

This parity work is one of a handful of omni-side deliverables that let Bytebase
drop `github.com/pingcap/tidb`. In `bytebase/bytebase@main` only 18 files still
import pingcap, all under `pkg/parser`/`pkg/types` (no planner/executor), behind
an omni-first / pingcap-fallback dispatcher
(`backend/plugin/parser/tidb/dispatcher.go`). The omni-side prerequisites, and
which Bytebase consumer each one unblocks:

| Omni deliverable | Package | Unblocks (Bytebase) | Status |
|---|---|---|---|
| **Catalog + walk-through parity** | `tidb/catalog`, `tidb/review/walk`, `metadata` | `schema/tidb/get_database_metadata.go`, `walk_through.go` (today fully pingcap: `parser.New()` + pingcap AST/format/mysql/types) | **this doc — 37 failures** |
| Parser acceptance completeness | `tidb/parser` | the dispatcher fallback — delete it once omni never falls back | rejection side (§3.6) closed; acceptance tracked by §3.3/§3.4 gap inventories |
| AST read-coverage | `tidb/ast` (257 nodes), `tidb/scope` | query-span extraction (`query.go`/`query_span_extractor.go`/`query_type.go`/`util.go`), `resource_change.go`, `statement_type.go` | ast already imported 57× by Bytebase; gaps surface when the (not-yet-started) TiDB query-span omni migration begins |
| Deparse / Restore | `tidb/deparse` | `parser/tidb/limit.go` (LIMIT rewrite, replacing pingcap `format.RestoreCtx` + `NewValueExpr`) | present; confirm coverage during that migration |
| Constraint AST nodes | `tidb/ast` | shared `parser/tokenizer/tokenizer.go` (pingcap `Constraint*` types) | present; mechanical swap |

Only the **first row is this phase's work** and the critical-path blocker; the
rest are prerequisites that are largely in place or are Bytebase-side rewrites
gated on omni capabilities that already exist. The largest un-started Bytebase
piece — TiDB query-span omni migration (peers exist for mssql/oracle/mysql/pg,
none for tidb) — needs no new omni package, only AST read-coverage above.

## Clusters

Count = failing top-level tests whose primary cause is this cluster (each test
assigned once; many tests touch several clusters). Locus **test** = re-baseline
expectations/fixtures to TiDB; **omni** = real catalog/deparse behavior gap.

| # | Cluster | Tests | Locus | Representative failing tests | Diff signature |
|---|---|---:|---|---|---|
| C1 | Collation/charset default metadata | 4 | test | `TestScenario_C4`, `TestScenario_C2`, `TestContainer_Section_2_11_CreateDropAlterDatabase` | `got "latin1_bin", want "latin1_swedish_ci"`; `got "utf8", want "utf8mb3"` — both omni and oracle return TiDB's value; `want` is a MySQL literal. TiDB also 1273-rejects `latin1_swedish_ci` at CREATE. |
| C2 | Numeric type-bound & zerofill validation | 3 | **omni** | `TestScenario_C25`, `TestScenario_C16`, `TestContainer_Section_1_7_DefaultValues` | `expected error for DECIMAL(66,0)/DECIMAL(40,31)/DECIMAL(5,6), got nil` — omni **over-accepts** out-of-range precision/scale TiDB rejects. Also `decimal(10,2) unsigned` missing `zerofill`. |
| C3 | Partition-expr validation + function casing | 2 | **omni** | `TestContainer_Section_4_1_Partitioning`, `TestScenario_C6` | `expected error for non-integer / misplaced-MAXVALUE / bare-TIMESTAMP partition expr, got nil` (over-accept); render `RANGE (year(...))` vs TiDB `RANGE (YEAR(...))`; `LINEAR HASH` normalized away by TiDB. |
| C4 | CHECK-constraint semantics | 4 | test | `TestScenario_C14`, `TestScenario_PS`, `TestScenario_C1`, `TestScenario_C5` | `expected 1 CHECK on t1, got 0`; `expected ER_CHECK_CONSTRAINT_DUP_NAME, got success`. TiDB here does not retain CHECKs / does not enforce dup-name; tests assert MySQL. |
| C5 | FK implicit naming & referential action | 1 | test | `TestScenario_C21` (also in C1/C4: `TestScenario_C1`, `TestScenario_PS`) | `got "...,fk_2", want "...,child_ibfk_2"` — omni's `fk_N` naming is TiDB-correct (`wt_3_3_test.go:453` cites pingcap `buildTableInfo`); `want child_ibfk_N` is the MySQL convention. `got "", want "NO ACTION"` — `showForeignKey` deliberately omits an unspecified action (matching TiDB SHOW CREATE). Re-baseline. One item to confirm: whether omni should surface `INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS.DELETE_RULE/UPDATE_RULE = "NO ACTION"` (TiDB does; omni's internal model holds `""`). |
| C6 | MySQL-only features in fixtures | 7 | test (+omni) | `TestContainer_Section_2_8_CreateDropIndex`, `TestScenario_C7`, `TestScenario_C20`, `TestContainer_Section_4_2_StoredRoutines` | Fixtures use `GEOMETRY/POINT/...`, `SPATIAL`/`FULLTEXT KEY`, stored routines (`8108 Unsupported type *resolve.NodeW`), `DEFAULT (FLOOR(...))`. TiDB 1064/8200/3770-rejects; omni parser also rejects spatial. **omni subset**: `DEFAULT (FLOOR(RAND()*100))` — omni accepts, TiDB 3770-rejects. |
| C7 | MySQL-only session vars + GIPK | 2 | test (+arch) | `TestScenario_C24`, `TestScenario_C3` | `1193 Unknown system variable` / `1238 read only variable` on MySQL `SET`; cascades to `1054 Unknown column`. GIPK (`sql_generate_invisible_primary_key`) has no omni session-var state — architectural. |
| C8 | TiDB DDL-semantics deltas (oracle) | 4 | test | `TestContainer_Section_2_2_AlterTableColumnOps`, `TestContainer_Section_2_3_AlterTableIndexOps`, `TestScenario_AX`, `TestScenario_C22` | `8200 can't drop column covered by composite index/PK`; `8200 Unsupported drop primary key ... clustered index`; `8200 Unsupported modify charset`; `ALGORITHM=/LOCK=` not validated; `CREATE INDEX IF NOT EXISTS` accepted. TiDB ≠ MySQL; tests assert MySQL. |
| C9 | Deparse / SHOW CREATE formatting | 7 | **omni** | `TestContainer_Section_3_6_ErrorContext`, `TestContainer_Section_1_9_GeneratedColumns`, `TestContainer_Section_3_3_ColumnErrors`, `TestScenario_C18` | `mismatch:` blocks and `got "year(2)", want "year"`; generated-column and error-context SHOW CREATE rendering deltas (function-name casing recurs here and in C3). **Not** `ROW_FORMAT`: C18's `18_6_row_format_elision` asserts both sides retain explicit `ROW_FORMAT=DYNAMIC`, but TiDB v8.5.5 omits it (`wt_13_2_test.go:36`) and omni already omits it — that subtest is a test re-baseline, not a renderer change. |
| — | Broad / cascade / misc | 3 | mixed | `TestSpotCheck_CatalogVerification`, `TestScenario_C15`, `TestScenario_C9` | Catch-all verification touching many clusters; `got false, want true` / `Unknown column` cascades downstream of a failed setup step. |

Soft-logged, not a hard failure (noted so fixers don't chase it): several
`TestContainer_Section_6_*` checks print `[~] partial: catalog internal state
matches, but INFORMATION_SCHEMA SQL queries require query engine`. omni's catalog
has no SQL query engine over `INFORMATION_SCHEMA`; those checks `t.Log` and pass.

## Real catalog-behavior bugs (higher risk — not mere formatting)

Flagged for careful handling; these change what omni accepts or stores:

- **C2** — catalog accepts out-of-range `DECIMAL(p,s)` (p>65, s>30, s>p) TiDB rejects. Over-acceptance. Locus: column/type validation in [define.go](../../tidb/catalog/define.go) / [table.go](../../tidb/catalog/table.go).
- **C3** — catalog accepts invalid partition expressions (non-integer, misplaced `MAXVALUE`, bare `TIMESTAMP`) TiDB rejects. Over-acceptance. Locus: [partition_interval.go](../../tidb/catalog/partition_interval.go) / [partition_key_coverage.go](../../tidb/catalog/partition_key_coverage.go).
- **C6 subset** — catalog accepts `DEFAULT (<disallowed function>)` TiDB 3770-rejects. Over-acceptance. Locus: default-value validation.

Reclassified into this list by wave 2 (PR #792) — the triage had marked them "test-shape / omni matches TiDB", but measuring omni directly showed real over-acceptance:

- **C8 DROP semantics (`Section_2_2`, `Section_2_3`, `TestScenario_AX`)** — catalog accepts `ALTER TABLE ... DROP COLUMN` covered by a composite index/PK, and `DROP PRIMARY KEY` on a clustered-index PK; TiDB v8.5.5 rejects both (8200). Over-acceptance.
- **C6 stored routines (`Section_4_2`)** — catalog accepts `CREATE PROCEDURE`; TiDB 8108-rejects stored routines entirely. Over-acceptance.
- **C7 index deltas (`TestScenario_C7` 7.3/7.7)** — catalog accepts prefix-less `TEXT KEY` (TiDB 1170-rejects) and keeps `IndexType="HASH"` for `KEY ... USING HASH` (TiDB coerces to BTREE). Over-acceptance.
- **C3 INT NULL PK (`TestScenario_C3` 3.4)** — catalog accepts `INT NULL PRIMARY KEY`; TiDB 1171-rejects. Over-acceptance.
- **C24 GIPK / C21** — architectural, not a point fix: GIPK has no omni session-var state; C21 needs the FK unspecified-action `NO ACTION` model plus the AST tri-state `ORDER BY` gap. Defer.

Reclassified out (omni already matches TiDB v8.5.5 — do **not** "fix" these or you reintroduce MySQL output):

- **C4 CHECK constraints** — TiDB v8.5.5 default (`tidb_enable_check_constraint` OFF) parses then silently discards every CHECK: no `CHECK_CONSTRAINTS`/`TABLE_CONSTRAINTS` row (TABLE_CONSTRAINTS has no `ENFORCED` column), SHOW CREATE omits it, predicate never evaluated, "forbidden" constructs and dup names accepted. omni's catalog discards CHECK identically. Test re-baseline (this wave): `TestScenario_C14` cleared; the CHECK subtests of `TestScenario_C1`/`C5`/`PS` cleared (those three stay red on unrelated non-CHECK subtests). Latent divergence, not test-exercised: omni discards CHECK unconditionally, so if a user sets `tidb_enable_check_constraint=ON` omni would under-retain — architectural follow-up, not a default-session failure.

- **C5 FK naming/action** — `fk_N` naming and omitted-vs-explicit referential action are deliberate TiDB behavior ([wt_3_3_test.go](../../tidb/catalog/wt_3_3_test.go), [show.go:483](../../tidb/catalog/show.go)). Test re-baseline. Sole open item: the `INFORMATION_SCHEMA` `DELETE_RULE/UPDATE_RULE = "NO ACTION"` mapping noted in C5 above.
- **C18 `ROW_FORMAT` elision** — omni omits `ROW_FORMAT` as TiDB v8.5.5 does ([wt_13_2_test.go](../../tidb/catalog/wt_13_2_test.go)); the C18 subtest asserting retention is the stale side. Test re-baseline, not a C9 renderer change.

## Recommended tackling order

Cheapest / highest-coverage first; the real-code gaps (higher risk) last.

1. **C1 — collation/charset metadata re-baseline** (test-only, mechanical). The renderer already emits TiDB values; only the `want` literals and `utf8`/`utf8mb3` naming need updating. Clears ~4 tests and the single largest assertion block.
2. **C9 — deparse/SHOW CREATE formatting** (omni, focused). Function-name casing (`YEAR`), `ROW_FORMAT`, generated-column/error-context rendering. One renderer fix clears ~7 tests and retires the casing half of C3.
3. **C5 + C8 + C7 + C6 — test re-baseline** (mechanical, parallelizable per file): re-baseline FK naming/action expectations to TiDB, adopt TiDB's DDL semantics, guard/remove MySQL-only `SET`s, drop MySQL-only fixtures (spatial, FULLTEXT/SPATIAL, stored routines). Also the C18 `ROW_FORMAT` subtest. ~14 tests; split across sessions by file.
4. **C4 — CHECK-constraint semantics** — DONE (this wave). Confirmed TiDB v8.5.5 default `tidb_enable_check_constraint=OFF` discards CHECK; omni matches. Re-baselined; cleared `TestScenario_C14`. The CHECK subtests of C1/C5/PS also re-baselined but those tests remain red on non-CHECK subtests (FK naming, index naming, datetime fsp — other clusters / real-bugs phase).
5. **C2 — numeric-bound validation** (omni, higher risk): add `DECIMAL(p,s)` bound and `zerofill` handling. ~3 tests.
6. **C3 — partition-expr validation** (omni, highest risk): add partition-expression rejection (non-integer, misplaced `MAXVALUE`, bare `TIMESTAMP`). ~2 tests; do last, with oracle confirmation per form.

After each cluster closes, drop its tests from the `startContainer` skip; delete
the gate when the count reaches zero.

## Appendix

Full `-v` log captured to the run scratchpad; regenerate with the measurement
command above. Container confirmed torn down (`docker ps` clean) after the run.
