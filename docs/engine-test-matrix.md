# Engine test database matrix

This records the test backends currently selected by the repository. A container that speaks an engine's wire protocol is not necessarily that engine. MySQL-backed MariaDB or OceanBase tests do not establish native dialect compatibility.

The initial helper-layout change preserves database versions, startup behavior, test expectations, and existing test gates. Native backend migrations, image upgrades, and the compatibility repairs they require belong in a follow-up change; this table makes those gaps explicit.

| Engine | Current live test backend | Remaining work |
| --- | --- | --- |
| `cassandra` | No local database oracle | Add native Cassandra coverage. |
| `cosmosdb` | No local database oracle | Add the official Cosmos DB NoSQL emulator. |
| `doris` | `dyrnq/doris:3.1.4` | Upgrade the native baseline after checking dialect drift. |
| `elasticsearch` | No local database oracle | Add native Elasticsearch request coverage. |
| `googlesql` | Spanner emulator `1.5.54` | Covers the GoogleSQL subset implemented by Spanner. Shared helper: `googlesql/internal/testcontainer`. |
| `mariadb` | Native `mariadb:11.8.8` | Keep MariaDB native. |
| `mongo` | No local database oracle | Add native MongoDB shell coverage. |
| `mssql` | SQL Server `2022-CU27-GDR1-ubuntu-22.04`; ScriptDOM reference | Move inline container setup into an engine-owned helper. |
| `mysql` | `mysql:8.0` | Consolidate on the selected MySQL 8.4 patch and investigate version drift. Shared helper: `mysql/internal/testcontainer`. |
| `oceanbase` | Native OceanBase `4.3.5.6-106000012026040916`, which the parser and validate oracles ask, plus inherited `mysql:8.0` / `mysql:8.4` oracles in catalog and review | Replace the inherited catalog and review oracles with OceanBase and repair native compatibility. Shared helper: `oceanbase/internal/testcontainer`. `oceanbase/parser`, `oceanbase/validate` and the tests of `oceanbase/catalog` that ask OceanBase (80, every oracle test among them) replay `compat.tsv` and tapes and start no server; catalog's container tests and `review` still start theirs. |
| `oracle` | Oracle Free `23.26.3-slim-faststart`; optional reference DSN | Move inline container setup into an engine-owned helper. |
| `partiql` | No local database oracle | Add DynamoDB Local for its supported PartiQL subset. |
| `pg` | `postgres:17-alpine`; review also checks `postgres:18-alpine` | Consolidate on the selected PostgreSQL 18 patch and investigate version drift. |
| `redshift` | PostgreSQL `17-alpine` baseline; explicit Redshift reference checks | Temporary PostgreSQL substitute is accepted; align its version with `pg` when upgrading. |
| `snowflake` | No local native database oracle | A native test account or suitable licensed emulator is needed. |
| `starrocks` | `starrocks/allin1-ubuntu:3.4.10` | Upgrade the native baseline after checking dialect drift. |
| `tidb` | `pingcap/tidb:v8.5.5` | Keep TiDB native. Catalog parity reached and the quarantine removed in #803; the differential tests run by default. It does not run against MySQL. `tidb/catalog` replays tapes and starts no server. |
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
