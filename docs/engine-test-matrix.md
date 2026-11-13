# Engine test database matrix

This records the test backends currently selected by the repository. A container that speaks an engine's wire protocol is not necessarily that engine. MySQL-backed MariaDB tests do not establish native dialect compatibility.

The initial helper-layout change preserves database versions, startup behavior, test expectations, and existing test gates. Native backend migrations, image upgrades, and the compatibility repairs they require belong in a follow-up change; this table makes those gaps explicit.

| Engine | Current live test backend | Remaining work |
| --- | --- | --- |
| `cassandra` | No local database oracle | Add native Cassandra coverage. |
| `cosmosdb` | No local database oracle | Add the official Cosmos DB NoSQL emulator. |
| `doris` | `dyrnq/doris:3.1.4` | Upgrade the native baseline after checking dialect drift. |
| `elasticsearch` | No local database oracle | Add native Elasticsearch request coverage. |
| `googlesql` | Spanner emulator `1.5.54` | Covers the GoogleSQL subset implemented by Spanner. Shared helper: `googlesql/internal/testcontainer`. |
| `mariadb` | Native `mariadb:11.8.8`, pinned by digest in `mariadb/testdata/pins.tsv`, started only by a recording (`OMNI_RECORD=1`) | Work the disagreements down: `compat.tsv` (12,960 rows, 12,082 agreeing; 404 only MariaDB accepts, 474 only omni) and `simulate.tsv` (1,452 rows on one fixture, 986 agreeing). Every package replays `compat.tsv` and tapes through `mariadb/internal/testcontainer` and starts no server: `mariadb/parser`, `mariadb/validate`, `mariadb/catalog` and `mariadb/review`. |
| `mongo` | No local database oracle | Add native MongoDB shell coverage. |
| `mssql` | SQL Server `2022-CU27-GDR1-ubuntu-22.04`; ScriptDOM reference | Move inline container setup into an engine-owned helper. |
| `mysql` | `mysql:8.0`, pinned as `mysql:8.0.46` in `mysql/testdata/pins.tsv`; review also checks `mysql:8.4` | `mysql/parser` replays `compat.tsv` (12,282 rows, 1,123 disagreeing) and tapes, and starts no server; so do `mysql/validate` and three tests of `mysql/catalog`. Move the rest of `mysql/catalog` and `mysql/review` to tapes, then `simulate.tsv`. Shared helper: `mysql/internal/testcontainer`. |
| `oceanbase` | Native OceanBase `4.3.5.6-106000012026040916`, started only by a recording (`OMNI_RECORD=1`) | Repair native compatibility: the gaps catalog records (`obGaps`), the disagreeing `compat.tsv` rows, and the 434 of `simulate.tsv`'s 2,276 rows (one fixture, four namings) where the catalog does not do what OceanBase did, 381 of them statements the parser rejects (Appendix A of the compatibility spec). Shared helper: `oceanbase/internal/testcontainer`. Every package replays `compat.tsv`, `simulate.tsv` and tapes and starts no server: `oceanbase/parser`, `oceanbase/validate`, `oceanbase/catalog` (162 tests) and `oceanbase/review` (the walk, its sweep, the scan rules and the column drops). No test under `oceanbase/` asks MySQL. |
| `oracle` | Oracle Free `23.26.3-slim-faststart`; optional reference DSN | Move inline container setup into an engine-owned helper. |
| `partiql` | No local database oracle | Add DynamoDB Local for its supported PartiQL subset. |
| `pg` | `postgres:17-alpine`; review also checks `postgres:18-alpine` | Consolidate on the selected PostgreSQL 18 patch and investigate version drift. |
| `redshift` | PostgreSQL `17-alpine` baseline; explicit Redshift reference checks | Temporary PostgreSQL substitute is accepted; align its version with `pg` when upgrading. |
| `snowflake` | No local native database oracle | A native test account or suitable licensed emulator is needed. |
| `starrocks` | `starrocks/allin1-ubuntu:3.4.10` | Upgrade the native baseline after checking dialect drift. |
| `tidb` | `pingcap/tidb:v8.5.5` | Keep TiDB native. Catalog parity reached and the quarantine removed in #803; the differential tests run by default. It does not run against MySQL. `tidb/catalog`, `tidb/parser` and `tidb/review` replay tapes and start no server, and the parser's corpora are `suite:` rows of its `compat.tsv`; `tidb/catalog/testdata/simulate.tsv` holds one fixture in four namings, and 573 of its 1,724 rows disagree. |
| `trino` | `trinodb/trino:482` | Shared helper: `trino/internal/testcontainer`. |

## Shared helper layout

Existing shared helpers use the same package path and filenames:

```text
<engine>/internal/testcontainer/
  container.go
  container_test.go  # only when the helper has tests
```

`container.go` contains image selection, startup, connections, protocol support, and cleanup. Tests of that helper belong in `container_test.go`. Do not add empty files or share helper implementations across engines. Import the engine-owned package without an alias and call it as `testcontainer.*`.

Inline setup in other engines has not yet been extracted by this initial change. The nested Spanner harness also remains an independent Go module.

## Selected upgrade targets

The follow-up targets are MySQL `8.4.13` only, PostgreSQL `18.6-alpine` for both PostgreSQL and the temporary Redshift baseline, Doris `4.0.8`, and StarRocks `4.1.7`. These are selected targets, not claims about the baselines currently exercised above. Pin a complete patch release or a dated vendor build when switching, then rerun the corresponding native oracle suites and investigate drift.

TiDB must use TiDB throughout the migration. Doris's reported `5.7.99` is its MySQL wire-protocol compatibility version, not a Doris release number.
