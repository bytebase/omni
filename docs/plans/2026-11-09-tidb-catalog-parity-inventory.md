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
| C5 | FK implicit naming & referential action | 1 | **omni** | `TestScenario_C21` (also in C1/C4: `TestScenario_C1`, `TestScenario_PS`) | `got "child_ibfk_1,fk_2", want "child_ibfk_1,child_ibfk_2"`; `got "", want "NO ACTION"` — omni's implicit FK naming and ON DELETE/UPDATE rendering diverge. |
| C6 | MySQL-only features in fixtures | 7 | test (+omni) | `TestContainer_Section_2_8_CreateDropIndex`, `TestScenario_C7`, `TestScenario_C20`, `TestContainer_Section_4_2_StoredRoutines` | Fixtures use `GEOMETRY/POINT/...`, `SPATIAL`/`FULLTEXT KEY`, stored routines (`8108 Unsupported type *resolve.NodeW`), `DEFAULT (FLOOR(...))`. TiDB 1064/8200/3770-rejects; omni parser also rejects spatial. **omni subset**: `DEFAULT (FLOOR(RAND()*100))` — omni accepts, TiDB 3770-rejects. |
| C7 | MySQL-only session vars + GIPK | 2 | test (+arch) | `TestScenario_C24`, `TestScenario_C3` | `1193 Unknown system variable` / `1238 read only variable` on MySQL `SET`; cascades to `1054 Unknown column`. GIPK (`sql_generate_invisible_primary_key`) has no omni session-var state — architectural. |
| C8 | TiDB DDL-semantics deltas (oracle) | 4 | test | `TestContainer_Section_2_2_AlterTableColumnOps`, `TestContainer_Section_2_3_AlterTableIndexOps`, `TestScenario_AX`, `TestScenario_C22` | `8200 can't drop column covered by composite index/PK`; `8200 Unsupported drop primary key ... clustered index`; `8200 Unsupported modify charset`; `ALGORITHM=/LOCK=` not validated; `CREATE INDEX IF NOT EXISTS` accepted. TiDB ≠ MySQL; tests assert MySQL. |
| C9 | Deparse / SHOW CREATE formatting | 7 | **omni** | `TestContainer_Section_3_6_ErrorContext`, `TestContainer_Section_1_9_GeneratedColumns`, `TestContainer_Section_3_3_ColumnErrors`, `TestScenario_C18` | `mismatch:` blocks and `got "year(2)", want "year"`; `ROW_FORMAT compact` vs `dynamic`; generated-column and error-context SHOW CREATE rendering deltas (function-name casing recurs here and in C3). |
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
- **C5** — implicit FK naming (`fk_N` vs `child_ibfk_N`) and ON DELETE/UPDATE action rendering diverge. Locus: [constraint.go](../../tidb/catalog/constraint.go) + FK rendering in [show.go](../../tidb/catalog/show.go).

## Recommended tackling order

Cheapest / highest-coverage first; the real-code gaps (higher risk) last.

1. **C1 — collation/charset metadata re-baseline** (test-only, mechanical). The renderer already emits TiDB values; only the `want` literals and `utf8`/`utf8mb3` naming need updating. Clears ~4 tests and the single largest assertion block.
2. **C9 — deparse/SHOW CREATE formatting** (omni, focused). Function-name casing (`YEAR`), `ROW_FORMAT`, generated-column/error-context rendering. One renderer fix clears ~7 tests and retires the casing half of C3.
3. **C8 + C7 + C6 — test re-baseline** (mechanical, parallelizable per file): adopt TiDB's DDL semantics, guard/remove MySQL-only `SET`s, drop MySQL-only fixtures (spatial, FULLTEXT/SPATIAL, stored routines). ~13 tests; split across sessions by file.
4. **C4 — CHECK-constraint semantics**: confirm TiDB's retention/dup-name behavior (`tidb_enable_check_constraint`) once, then re-baseline. ~4 tests.
5. **C2 — numeric-bound validation** (omni, higher risk): add `DECIMAL(p,s)` bound and `zerofill` handling. ~3 tests.
6. **C3 partition-expr validation + C5 FK naming/action** (omni, highest risk): add partition-expression rejection and align implicit FK naming/action. ~3 tests; do last, with oracle confirmation per form.

After each cluster closes, drop its tests from the `startContainer` skip; delete
the gate when the count reaches zero.

## Appendix

Full `-v` log captured to the run scratchpad; regenerate with the measurement
command above. Container confirmed torn down (`docker ps` clean) after the run.
