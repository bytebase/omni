# Cutover dispatch: how each engine reaches its omni driver

2026-11-17 · Proposed

Bytebase `345007290e`, omni `6dee22d57`. This note is the dispatch half of the driver move:
[2026-11-15-driver-contract-design.md](2026-11-15-driver-contract-design.md) fixes the contract and
[2026-11-17-driver-adapter-prototype.md](2026-11-17-driver-adapter-prototype.md) measures what an
adapter body costs. Neither says how the twenty-four real `storepb.Engine` values reach twenty
driver trees on main through Bytebase's registry, which engines share code today, or which engine
can be flipped first without a product change. That is this note.

The short answer: the registry does not move. A cutover replaces the body of one
`backend/plugin/db/<engine>` package and the caller-side conversions around it. There is no driver
alias to add: every engine already registers its own constructor, including the engines the design
calls aliases. The alias table describes language trees, not dispatch.

## 1. The seam

Bytebase's `db` package is a map, not a switch:

```go
// backend/plugin/db/driver.go
var drivers = make(map[storepb.Engine]driverFunc) // :68
func Register(dbType storepb.Engine, f driverFunc) // :174
func Open(ctx, dbType storepb.Engine, cfg ConnectionConfig) (Driver, error) // :187
```

Each engine package registers in `init`:

```go
func init() { db.Register(storepb.Engine_POSTGRES, newDriver) } // db/pg/pg.go:45-46
```

`backend/server/ultimate.go` blank-imports all twenty-three driver packages, one of which — `db/mysql` —
registers two engine values. That is the whole dispatch. `db.Open` looks up the enum, constructs the
package's `Driver`, and calls its `Open(ctx, dbType, config)`. `backend/component/dbfactory/dbfactory.go:63`
is the one caller that resolves the data source and reaches `db.Open`; everything else goes through
the factory.

`storepb.Engine` reaches the package as the second argument of `Driver.Open`. Most packages ignore
it. `db/mysql` stores it and branches three times: `mysql/sync.go:548` and `:588` (MySQL vs MariaDB
sync), and `mysql/stats.go:28` (MariaDB's `EXPLAIN`). `db/snowflake` stores it and never reads it.
So the enum's only dispatch use inside a driver is the MySQL/MariaDB split; a cutover moves that
branch into the adapter's choice of `mysql/driver` or `mariadb/driver`, not into the registry.

The map is keyed by enum, so a driver alias is not a lookup rewrite: it is a second `Register` call
for the same constructor, which is exactly what `db/mysql` does for MariaDB. `Register` panics on a
duplicate key, so a new package taking over `MARIADB` must remove that second call in the same
change.

What a cutover replaces:

| Stays | Replaced |
| --- | --- |
| `db.Register`, `db.Open`, the map, the blank import in `ultimate.go` | the engine package's `Open`, `Close`, `Ping`, `Execute`, `QueryConn`, `SyncInstance`, `SyncDBSchema`, `Dump` |
| `storepb.DataSource` → `driver.Config` conversion, secret resolution, IAM tokens, SSH and TLS material | the per-engine config assembly inside `Open` |
| Product switches that name engines: `db/explain_format.go`, `parser/base` and `schema` registries, `common/engine.go`, plancheck, `convertToParserEngine` | the concrete driver type those switches may assert (`*pg.Driver`, …) |

The last row is why a cutover is not "change the `Open` line". The registry is one lookup; the
`db.Driver` method set is `GetDB`, `QueryConn(*sql.Conn)`, `ExecuteOptions.Transaction` and
`[]*v1pb.QueryResult`, and those force caller changes that the prototype recorded. A clean engine
cutover changes one package and no product code. The rest need the product work in §4.

## 2. The picture today

`tapes` means the package and its recorded tests are on `origin/main`; `draft` means the package
exists only on the named open branch and has never been recorded against the engine. Cutover order
is the number in the first column; the waves are explained in §4.

| # | Engine | Bytebase package | Shares / aliases today | omni driver | Recommendation | Blockers |
| --- | --- | --- | --- | --- | --- | --- |
| 1 | `CASSANDRA` | `db/cassandra` | own | `cassandra/driver` (tapes) | thin adapter | none |
| 2 | `COSMOSDB` | `db/cosmosdb` | own | `cosmosdb/driver` (tapes) | thin adapter; `Params` auth/tenant/emulator, `QueryOptions.Schema` is the container | none |
| 3 | `DYNAMODB` | `db/dynamodb` | own; PartiQL is its language tree | `dynamodb/driver` (tapes) | thin adapter; `Params` region/sessionToken/accountId | none |
| 4 | `MONGODB` | `db/mongodb` | own; omni tree is named `mongo` | `mongo/driver` (tapes) | thin adapter; `Open` takes `Options{SRV}` | none |
| 5 | `REDIS` | `db/redis` | own | `redis/driver` (tapes) | thin adapter; `Options` topology, and the session probe for all-nodes | none |
| 6 | `SPANNER` | `db/spanner` | own; parser handled by `googlesql` | `spanner/driver` (tapes) | thin adapter; `Params` project/instance | none; `Open` reads the dialect |
| 7 | `HIVE` | `db/hive` | own; parser is `standard` | `hive/driver` (tapes) | thin adapter | no caller `Dial` (forward a local port) |
| 8 | `ELASTICSEARCH` | `db/elasticsearch` | own | `elasticsearch/driver` (tapes) | thin adapter; `Options{OpenSearch, Transport}` | SigV4 is a caller `Transport`; `Dial` also works |
| 9 | `TRINO` | `db/trino` | own | `trino/driver` (tapes) | thin adapter | session |
| 10 | `CLICKHOUSE` | `db/clickhouse` | own; parser is `standard`/`mysql` | `clickhouse/driver` (tapes) | thin adapter; `Dumper` now exists | session |
| 11 | `DORIS` | `db/doris` | own; imports `db/mysql` helpers | `doris/driver` (tapes) | thin adapter + shared MySQL config builder | session |
| 12 | `STARROCKS` | `db/starrocks` | own; imports `db/mysql` helpers | `starrocks/driver` (tapes) | thin adapter + shared MySQL config builder | session |
| 13 | `MSSQL` | `db/mssql` | own | `mssql/driver` (tapes) | thin adapter + `CountAffectedRows` | session, plancheck |
| 14 | `COCKROACHDB` | `db/cockroachdb` | own; parser handlers are pg's, design says alias pg | `cockroachdb/driver` (tapes) | thin adapter; do not alias pg | session, plancheck |
| 15 | `POSTGRES` | `db/pg` | own; parser also serves CockroachDB | `pg/driver` (tapes) + `Options{TenantMode}` | thin adapter | session, plancheck, advisor |
| 16 | `MYSQL` | `db/mysql` | registers `MYSQL` and `MARIADB` | `mysql/driver` (tapes) | thin adapter; keep the shared MySQL config builder | session, plancheck, advisor, ghost |
| 17 | `MARIADB` | `db/mysql` | served inside `db/mysql`; DiffSchema aliases it to `MYSQL` | `mariadb/driver` (tapes) | branch inside `db/mysql` on `dbType` | session, plancheck, advisor, ghost |
| 18 | `TIDB` | `db/tidb` | own; imports `db/mysql` helpers | `tidb/driver` (tapes) | thin adapter + shared MySQL config builder | session, plancheck, advisor |
| 19 | `OCEANBASE` | `db/oceanbase` | own; DiffSchema aliases it to `MYSQL` | `oceanbase/driver` (tapes) | thin adapter | session, plancheck, advisor |
| 20 | `ORACLE` | `db/oracle` | own | `oracle/driver` (tapes) | thin adapter + `Version`, `CountAffectedRows` | session, plancheck, advisor, version |
| 21 | `REDSHIFT` | `db/redshift` | own | `redshift/driver` draft, PR #1371 | thin adapter after recording; `Dump` stays a no-op like Bytebase's | session, plancheck, unverified |
| 22 | `SNOWFLAKE` | `db/snowflake` | own | `snowflake/driver` draft, PR #1378 | thin adapter after recording | session, unverified |
| 23 | `BIGQUERY` | `db/bigquery` | own; parser handled by `googlesql` | `bigquery/driver` draft, PR #1377 | thin adapter over `bigquery/driver`, not the `googlesql` alias | unverified; no tape seam yet |
| 24 | `DATABRICKS` | `db/databricks` | own; parser is `standard` | `databricks/driver` draft, PR #1379 | thin adapter after recording | unverified; no `catalog`, so its `Dump` is server DDL |

Blockers:

- **session**: Bytebase's editor opens a `*sql.Conn` from `GetDB` and reuses it across statements
  (`api/v1/sql_service.go:117`, `:404`, `:1347`). The adapter has no `*sql.DB` to hand back, so a
  cutover must hold a `driver.Session` per editor session, or accept one session per call. Engines
  whose `QueryConn` ignores `conn` already behave as the second.
- **plancheck**: `runner/plancheck/statement_report_executor.go:168-220` type-asserts
  `*pgdriver.Driver`, `*redshiftdriver.Driver`, `*cockroachdbdriver.Driver`,
  `*mysqldriver.Driver` (for MySQL and MariaDB), `*oceanbasedriver.Driver`, `*tidbdriver.Driver`,
  `*oracledriver.Driver`, `*mssqldriver.Driver` to reach `CountAffectedRows`. It must name an
  interface instead.
- **advisor**: the V1 live rules take `advisor.Context.Driver *sql.DB`
  (`plugin/advisor/advisor.go:51`). `pg`, `mysql`, `mariadb`, `oceanbase`, `tidb` and `oracle` have
  rules that EXPLAIN or query through it. A nil `GetDB` makes them silently skip, not fail.
- **ghost**: MySQL's gh-ost validator and cleanup reach `GetDB` directly
  (`component/ghost/validator.go:54`, `runner/taskrun/database_migrate_executor.go:369`).
- **version**: `runner/taskrun/database_migrate_executor.go:827` casts the Oracle driver and calls
  `GetVersion() *plsql.Version`; `oracle/driver` exposes `Version(ctx)`.

## 3. Where the two pictures disagree

### MariaDB is the only multiplexed driver

`db/mysql/mysql.go:52-53` registers both `MYSQL` and `MARIADB` on one constructor. The shared
package is real: the two engines use the same client and config builder, and three sync/stats
branches read `dbType`. `db/tidb`, `db/doris` and `db/starrocks` are separate packages that import
`db/mysql` helpers (`UnquoteMySQLString`, `IsCurrentTimestampLike`, `UnescapeExpressionDefault`,
`RegisterDialContext`, `RegisterTLSConfig`), so the MySQL family shares *code* in four packages but
only `db/mysql` shares a *registration*.

omni has separate `mysql/driver` and `mariadb/driver` trees. Recommendation: keep `db/mysql` as the
registered package for both engine values for the first cutover and choose the driver from
`dbType` inside its adapter. This preserves the concrete type the plancheck assertion expects and
keeps the six engines that use the MySQL config builder on one path. Splitting `db/mariadb` into
its own package is cleaner eventually, but it is a second product change (the plancheck assertion
and its import) and `db/mysql` would have to drop its `MARIADB` registration in the same change.
**Judgement call.**

### CockroachDB is not a pg driver

The design's alias table says `cockroachdb→pg`. Bytebase has `db/cockroachdb` and omni has
`cockroachdb/driver`, whose package doc says why: the wire protocol is PostgreSQL's, but the
catalog, system databases, indexes, EXPLAIN options and plan decoding are CockroachDB's. Bytebase's
parser does borrow pg handlers for `COCKROACHDB` (`parser/pg/split.go`, `parser/pg/query.go`,
`parser/pg/query_span.go`, `parser/pg/diagnose.go`, plus limit, completion and explain), and
`db/cockroachdb` already copies pg's connection setup rather than calling it. Recommendation: thin
adapter over `cockroachdb/driver`; do not alias the pg driver. The alias is a language-tree fact.

### BigQuery and Spanner are not googlesql drivers

The design aliases `bigquery/spanner→googlesql`. That is true of parsing: `parser/bigquery/plugin.go`
and `parser/spanner/plugin.go` both call `googlesql.Register`. It is not true of the wire:
`googlesql` has no driver, `spanner/driver` is its own tree, and `bigquery/driver` is a draft beside
it. Recommendation: one thin adapter per service, over its own driver. Spanner's driver is on main
with tapes; BigQuery's is unverified.

### DynamoDB is not a PartiQL driver

`parser/partiql` registers `DYNAMODB` for parse, split, validate and completion, so the language
tree is shared. The wire is `dynamodb/driver`, which builds an AWS client and decodes DynamoDB
values. Recommendation: thin adapter over `dynamodb/driver`; do not route the engine through a
hypothetical PartiQL driver.

### Mongo is a name difference, not an alias

Bytebase's package is `db/mongodb`; omni's tree is `mongo`. The engine enum is `MONGODB`. Keep the
Bytebase path and registration, and import `omni/mongo/driver`. No alias table entry is needed; a
rename of either side would be churn.

### ClickHouse's prototype gap is closed

The adapter prototype says `clickhouse/driver` has no `Dumper`. Current main has one
(`clickhouse/driver/dump.go`, #1350): it asks the server for `SHOW CREATE` and ignores the metadata,
as Bytebase's `Dump` did. ClickHouse is no longer the one engine that cannot delegate `Dump`. Its
remaining cutover blocker is the session one, because Bytebase's `QueryConn` runs through the
editor's `*sql.Conn`.

### Redis is clean, but not through `QueryConn` alone

`redis/driver` implements the contract and an `AllNodesQuerier` probe on its session. Bytebase's
editor offers "run on every node" via `QueryOption_ALL_NODES`, which has no field in
`driver.QueryOptions`. The adapter must ask the session for the probe when that option is set. This
is adapter code, not product code.

### Hive cannot take a caller's dialer; Elasticsearch takes one or a transport

`hive/driver` cannot accept `Config.Dial`: gohive v2 keeps its dialer in unexported fields and has
no field for one, so a caller forwards a local port and names it as `Config.Host` and `Config.Port`.
`elasticsearch/driver` does take `Config.Dial`; `Options.Transport` is the separate seam for a
managed OpenSearch domain's AWS SigV4 signer, and the driver refuses both together. Both engines
can cut over; only Hive behind an SSH tunnel needs the local-port workaround.

### The alias table is spread across registries

There is no single engine→alias table in Bytebase. `convertToParserEngine`
(`api/v1/database_service.go:1309`) maps MariaDB and OceanBase to MySQL for DiffSchema and leaves
CockroachDB as itself. Parser aliases are in package `init`s: `parser/bigquery` and `parser/spanner`
register through `googlesql`, `parser/partiql` registers `DYNAMODB`, `parser/pg` also registers
`COCKROACHDB` handlers. Capability lists are switches in `common/engine.go`. None of these is the
driver registry, and a cutover does not change them.

## 4. Order

Waves are by the blocker, not by engine importance.

**Wave 1 — no product change (1-8).** These engines already return a nil `GetDB` or ignore the
`*sql.Conn`, and no product package names their concrete driver type. `cassandra`, `cosmosdb`,
`dynamodb`, `mongodb`, `redis`, `spanner`, `hive`, `elasticsearch`. Start with `cassandra`: it has
no extra `Open` options, no engine-specific `Params`, and tapes on main. Then take the other seven
in any order; `hive` carries the no-dialer limitation above, while `elasticsearch`'s AWS SigV4
signer is caller code rather than a driver limitation.

**Wave 2 — the session change (9-12).** `trino`, `clickhouse`, `doris`, `starrocks` run queries
through the editor's `*sql.Conn`, but their engines appear in no plancheck, advisor or ghost switch.
Do the product change once: let the SQL editor hold a `driver.Session` instead of a `*sql.Conn` (or
open one session per call and record the behavior change), then flip these four. `doris` and
`starrocks` also need the shared MySQL config builder from wave 3's family work, but not the
plancheck or advisor changes.

**Wave 3 — the product changes (13-20).** First replace plancheck's concrete assertions with a
`CountAffectedRows` interface, which unblocks `mssql` and `cockroachdb` without touching advisors.
Then replace the advisor's `*sql.DB` with a probe interface and the ghost path with driver methods,
which unblocks `pg`, `mysql`, `mariadb`, `tidb` and `oceanbase`. Oracle also needs its
`Version(ctx)` call site changed. Within this wave, `mssql` and `cockroachdb` are the cheap ones;
`mysql`/`mariadb` are the expensive pair because they carry ghost and the live advisors; `pg` is
the reference adapter and should be built on a branch early even though its product flip waits.

**Wave 4 — the unverified drivers (21-24).** `redshift`, `snowflake`, `bigquery` and `databricks`
have draft drivers on open branches, no tapes, no sync goldens and no pin rows. A cutover cannot be
validated by compiling. Redshift has no `Dump`, but Bytebase's wrote nothing either, so the adapter
keeps returning nil; Snowflake has no `CountAffectedRows`, which Bytebase never answered; BigQuery's
`Dump` is deliberately a no-op, as Bytebase's was, and its client has no tape seam; Databricks has
no `catalog`, so its dump writes server DDL. Record each on a scratch account or cluster, land its
pin, then cut over. None of these should be the first cutover.

The contract's own order — pg sync and dump, then pg query and execute, to settle the interfaces —
is not this order. That order is about proving the omni contract; this one is about not breaking a
running Bytebase. Build the pg adapter early as the reference, but flip clean engines first.

## 5. Start the first cutover

1. Pick `CASSANDRA` (row 1). It registers one engine in `db/cassandra/cassandra.go`, its `GetDB`
   already returns nil, its `QueryConn` ignores the `*sql.Conn`, and no product package names its
   concrete type.
2. Keep the package path and `db.Register(storepb.Engine_CASSANDRA, …)`. Replace the body of
   `Driver` with an adapter over `omni/cassandra/driver.Open`. Do not touch `db.Open`, the map, or
   `backend/server/ultimate.go`.
3. Convert `storepb.DataSource` and `ConnectionContext` into `driver.Config` in the adapter:
   host/port, database, user/password, `TLS func(host)`, `Dial` for the SSH tunnel, `Params`,
   application name, read-only, engine version. Cache nothing the old driver did not cache.
4. Convert the results at the seam: `driver.Result` → `v1pb.QueryResult`, `driver.Instance`/`Role`
   → `db.InstanceMetadata`, `driver.Message` → `QueryResult_Message`, and the task-run logger to
   `driver.ExecuteLogger`. These conversions are the same for every engine and belong in a shared
   caller package, not in each adapter.
5. Keep `GetDB` returning nil. For cassandra nothing calls it; for the rest of wave 1 nothing calls
   it in a way that needs a `*sql.DB`.
6. Run the existing Bytebase smoke tests for the engine and omni's `cassandra/driver` tape suite.
   Then take `cosmosdb`, `dynamodb`, `mongodb`, `redis`, `spanner`, `hive`, `elasticsearch` one
   package at a time, with the same shape.

## 6. What stays product-side

The cutover does not move these, and a proposal that does is a different change:

- The `storepb.Engine` enum, `db.Register`/`db.Open`, the blank imports, and `dbfactory`.
- Data-source resolution, external secret replacement, AWS/GCP IAM tokens, SSH tunnels, TLS
  material, and the Cloud SQL / RDS dialers.
- The parser and schema registries, `convertToParserEngine`, `common/engine.go` capability lists,
  `db/explain_format.go`, advisory registration, plancheck statement typing, and the SQL editor's
  session policy.
- The task-run logger renderer. `driver.ExecuteLogger` is the bridge; `storepb.TaskRunLog` and the
  command range are product.

The shared MySQL-family config builder (the prototype measures 158 lines, 82 of them the RDS CA
bundle) stays in Bytebase, because six adapters reach it and no engine owns it. Put it in one
package beside the adapters rather than copying it six times.
