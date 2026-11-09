# TiDB catalog parity residual (merged-tree)

Date: 2026-11-09. Follow-up to the quarantine inventory
([2026-11-09-tidb-catalog-parity-inventory.md](2026-11-09-tidb-catalog-parity-inventory.md)),
measured after the nine cluster PRs (#791–#799) merged to `main`. The inventory
clustered the quarantined failures on separate branches; this remeasures on the
fully-merged tree so the cumulative effect is realized, and maps what is left so
the cleanup can be finished and the quarantine removed.

## The measurement

`TIDB_CATALOG_PARITY=1 go test ./tidb/catalog/... -count=1 -v`, one shared TiDB
v8.5.5 container (`startTiDBForCatalog`, cleaned up in `TestMain`). The gate is
`startContainer` in [tidb/catalog/container_test.go](../../tidb/catalog/container_test.go);
every differential test routes through it and is skipped unless the env var is set.

| | count |
|---|---:|
| failing top-level tests | **22** |
| passing | 2721 |
| skipped (non-container) | 5 |

Down from 37 at inventory time: many top-level tests failed on several clusters
at once and only flip green once every one of their clusters is fixed, so the
merged drop is larger than any per-PR delta. It is not zero — the quarantine
stays until these 22 are cleared.

The shared reset (`resetSharedContainer`) drops every user database between
tests, so the recurring `child_ibfk_5,fk_2` signatures are **not** cross-test
state leakage: the live-TiDB oracle and omni return the same value, so omni is
at parity and the hard-coded `want` is MySQL-shaped. That stale-expectation
pattern is the bulk of the residual; a minority are genuine omni catalog gaps,
flagged below as the higher-risk work.

## The 22, by cluster

### Bucket 1 — Stale MySQL-shaped expectations (oracle == omni; re-baseline)

omni already matches live TiDB; the test `want` is MySQL's answer.

- `TestScenario_C1` — `1_1`/`1_2` FK auto-names (TiDB `child_ibfk_5,fk_2`), `1_10` unique-primary fallback (`PRIMARY` vs `PRIMARY_2`), `1_11`/`1_12` functional-index name/hidden col
- `TestScenario_PS` — `PS_3`/`PS_4` FK counter, `PS_7` FK name collision (TiDB allows the duplicate name)
- `TestScenario_C2` — `2_8` `zerofill` dropped from `COLUMN_TYPE`, `2_14` FLOAT→`double`, `2_17`/`2_18`/`2_23`/`2_25`/`2_26` type naming/promotion
- `TestScenario_C8` — `8_2` `ROW_FORMAT` `compact` vs `dynamic`
- `TestScenario_C15` — `15_5` multi-`ADD ... AFTER` left-to-right (TiDB resolves `AFTER` against the pre-ALTER table)
- `TestSpotCheck_CatalogVerification` — `C1_1`/`PS1` FK/CHECK counter, `C6_2` subpartition count

### Bucket 2 — Container-environment / TiDB config (oracle-side setup errors)

Tests assume a MySQL environment the TiDB container does not provide.

- `TestScenario_C3` — `3_1`/`3_8` `SET SESSION explicit_defaults_for_timestamp=0` → Error 1238 (read-only variable in TiDB)
- `TestScenario_C4` — `CREATE DATABASE ... latin1_swedish_ci` → Error 1273 (unsupported under `new_collations_enabled`); `utf8`/`utf8mb3` naming (`4_4`/`4_8`/`4_9`/`4_10`/`4_11`)
- `TestScenario_C8` — `8_4`/`8_5` latin1 collations unsupported; `8_10` `information_schema.INNODB_TABLES` absent (Error 1146)
- `TestSpotCheck_CatalogVerification` — `C3_1` read-only var, `C4_1` latin1 unsupported

### Bucket 3 — Genuine omni over-acceptance (real fixes; omni accepts, TiDB rejects)

- `TestScenario_C1` — `1_8` index explicitly named `PRIMARY` not rejected (TiDB 1280)
- `TestScenario_C5` — `5_2` FK `SET NULL` on a `NOT NULL` column; `5_7` FK on a VIRTUAL generated column (TiDB 3104)
- `TestScenario_C9` — `9_2` FK on VIRTUAL gcol, `9_3` VIRTUAL gcol as PRIMARY KEY, `9_4` non-deterministic gcol expression (`NOW()`)
- `TestScenario_AX` — `AX_3` DROP of the last remaining column not rejected (TiDB 1090)

### Bucket 4 — Documented known bugs (datetime precision / default functions)

Tracked historically in `mysql/catalog/scenarios_bug_queue/c16.md`; self-labeled
"KNOWN BUG" in the tests. Need real omni validation, not a skip.

- `TestScenario_C16` — `16_2`/`16_4`/`16_5`/`16_6`/`16_8`/`16_9`/`16_10`/`16_11`/`16_12`: fsp > 6 not rejected, `DEFAULT SYSDATE` / `DEFAULT UTC_TIMESTAMP` / fsp-mismatch `DEFAULT NOW()` not errored
- `TestScenario_PS` — `PS_5` datetime fsp-mismatch default

### Bucket 5 — Error-message text mismatches (omni message vs TiDB wording)

- `TestContainer_Section_3_2_TableErrors` — `1050_table_already_exists`
- `TestContainer_Section_3_3_ColumnErrors` — `1054_unknown_column`
- `TestContainer_Section_3_4_IndexKeyErrors` — `1091_cant_drop_key`
- `TestContainer_Section_3_6_ErrorContext` — identifier quoting, `IF NOT EXISTS` suppression
- `TestScenario_C19` — expression-index rejection wording (3762/3753/3757) and `19_1`/`19_2`/`19_3` functional-index `COLUMN_NAME` metadata

### Bucket 6 — Feature / partition metadata

- `TestScenario_C6` — omni under-reports partition metadata: default HASH/KEY partition names empty (want `p0`), subpartition count 0 (`6_1`/`6_2`/`6_5`/`6_6`/`6_15`); reconcile against the stated v8.5.5 "drops LINEAR/subpartition info" assumption in `tablecmds.go` before changing
- `TestContainer_Section_2_5_AlterTableTableLevel` — ALTER charset 8200, FULLTEXT 8200, GEOMETRY 1064 assertion shapes
- `TestContainer_Section_4_1_Partitioning` — `partition_by_range`, `partition_subpartition`
- `TestContainer_Section_6_1_ShowCreateTableIntegration` — `all_data_types`/`all_defaults`/`all_index_types` container CREATE + partitioned SHOW CREATE
- `TestContainer_Section_6_2_ShowCreateOtherObjects` — `CREATE PROCEDURE` (8108), deparse reserved word `select`, JSON-extract type errors

## Sequencing

Re-baselines read omni's final output, and the quarantine deletion needs every
test green, so the work is sequential, not parallel; a server boot per package
and the shared-devbox OOM budget also argue against many TiDB containers at
once. Suggested order: land the real omni fixes first (buckets 3, 4, and the C6
partition-metadata gap), then the mechanical re-baselines (buckets 1, 2, 5 and
the remaining bucket 6), then remove the `startContainer` quarantine so the
differential suite runs by default.
