# Cutover dispatch: how each engine reaches its omni driver

2026-11-17 · Proposed · revised with measurements

Bytebase `345007290e`, omni `6dee22d57`. This note is the dispatch half of the driver move:
[2026-11-15-driver-contract-design.md](2026-11-15-driver-contract-design.md) fixes the contract and
[2026-11-17-driver-adapter-prototype.md](2026-11-17-driver-adapter-prototype.md) measures what an
adapter body costs. Neither says how the twenty-four real `storepb.Engine` values reach twenty
driver trees on main through Bytebase's registry, which engines share code today, or which engine
can be flipped first without a product change. That is this note.

Revised after the first cutovers were built: §3 replaces the adapter predictions with measured
costs, §5 records the two places a prototype could not see the shape of the work, and §6 reorders
the waves. The twenty-four rows in §2 keep the order the prototypes proposed, so a reader can see
what moved; each row carries a measured value or says it has none.

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
cutover changes one package and no product code. The rest need the product work in §6.

## 2. The picture today

`tapes` means the package and its recorded tests are on `origin/main`; `draft` means the package
exists only on the named open branch and has never been recorded against the engine. Cutover order
is the number in the first column; the waves are explained in §6. The last column carries the
measured adapter lines and product-side file count from §3; `unmeasured` means the prototype's
estimate is all there is, and the row is a prediction, not a cost.

| # | Engine | Bytebase package | Shares / aliases today | omni driver | Recommendation | Blockers | Measured |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 1 | `CASSANDRA` | `db/cassandra` | own | `cassandra/driver` (tapes) | thin adapter | measured: none. `GetDB` returns nil, `QueryConn` ignores its `*sql.Conn`, no product file names the driver type | 312 / 0 callers |
| 2 | `COSMOSDB` | `db/cosmosdb` | own | `cosmosdb/driver` (tapes) | thin adapter; `Params` auth/tenant/emulator, `QueryOptions.Schema` is the container | none | unmeasured |
| 3 | `DYNAMODB` | `db/dynamodb` | own; PartiQL is its language tree | `dynamodb/driver` (tapes) | thin adapter; `Params` region/sessionToken/accountId | none | unmeasured |
| 4 | `MONGODB` | `db/mongodb` | own; omni tree is named `mongo` | `mongo/driver` (tapes) | thin adapter; `Open` takes `Options{SRV}` | none | unmeasured |
| 5 | `REDIS` | `db/redis` | own | `redis/driver` (tapes) | thin adapter; `Options` topology, and the session probe for all-nodes | none | unmeasured |
| 6 | `SPANNER` | `db/spanner` | own; parser handled by `googlesql` | `spanner/driver` (tapes) | thin adapter; `Params` project/instance | none; `Open` reads the dialect | unmeasured |
| 7 | `HIVE` | `db/hive` | own; parser is `standard` | `hive/driver` (tapes) | thin adapter | no caller `Dial` (forward a local port) | unmeasured |
| 8 | `ELASTICSEARCH` | `db/elasticsearch` | own | `elasticsearch/driver` (tapes) | thin adapter; `Options{OpenSearch, Transport}` | SigV4 is a caller `Transport`; `Dial` also works | unmeasured |
| 9 | `TRINO` | `db/trino` | own | `trino/driver` (tapes) | thin adapter | measured: none. No caller changed; `QueryConn` ignores its `*sql.Conn`, so a session is per call | 234 / 0 callers |
| 10 | `CLICKHOUSE` | `db/clickhouse` | own; parser is `standard`/`mysql` | `clickhouse/driver` (tapes) | thin adapter; `Dumper` now exists | session, as trino's (per call, not a caller change) | unmeasured |
| 11 | `DORIS` | `db/doris` | own; imports `db/mysql` helpers | `doris/driver` (tapes) | thin adapter + shared MySQL config builder | session; the shared builder is a precondition | unmeasured |
| 12 | `STARROCKS` | `db/starrocks` | own; imports `db/mysql` helpers | `starrocks/driver` (tapes) | thin adapter + shared MySQL config builder | session; the shared builder is a precondition | unmeasured |
| 13 | `MSSQL` | `db/mssql` | own | `mssql/driver` (tapes) | thin adapter + `CountAffectedRows` | session, plancheck, masking consumers | unmeasured (proof running) |
| 14 | `COCKROACHDB` | `db/cockroachdb` | own; parser handlers are pg's, design says alias pg | `cockroachdb/driver` (tapes) | thin adapter; do not alias pg | plancheck only | unmeasured |
| 15 | `POSTGRES` | `db/pg` | own; parser also serves CockroachDB | `pg/driver` (tapes) + `Options{TenantMode}` | thin adapter | measured: the advisor probe is forced and `GetDB` retires in 8 files. plancheck was **not** forced; `QueryConn`'s `*sql.Conn` is dead but tolerated | 658 / 16 files |
| 16 | `MYSQL` | `db/mysql` | registers `MYSQL` and `MARIADB` | `mysql/driver` (tapes) | thin adapter; keep the shared MySQL config builder | measured: an engine-specific `DriverHandle` for the advisors, plancheck's 7 assertions plus 9 blank parser imports, `GetDB`, and gh-ost failing closed | 706 / 17 files |
| 17 | `MARIADB` | `db/mysql` | served inside `db/mysql`; DiffSchema aliases it to `MYSQL` | `mariadb/driver` (tapes) | branch inside `db/mysql` on `dbType` | as `MYSQL`, and it pays half the advisor rules twice | unmeasured |
| 18 | `TIDB` | `db/tidb` | own; imports `db/mysql` helpers | `tidb/driver` (tapes) | thin adapter + shared MySQL config builder | session, plancheck, advisor | unmeasured |
| 19 | `OCEANBASE` | `db/oceanbase` | own; DiffSchema aliases it to `MYSQL` | `oceanbase/driver` (tapes) | thin adapter | session, plancheck, advisor | unmeasured |
| 20 | `ORACLE` | `db/oracle` | own | `oracle/driver` (tapes) | thin adapter + `CountAffectedRows` | measured: the `Prober` seam (reused verbatim from pg), and the private-link guard loses its `*sql.Conn`. plancheck untouched; `CountAffectedRows` needed no product change; the version blocker dissolved into an 8-line deletion | 390 / 12 files |
| 21 | `REDSHIFT` | `db/redshift` | own | `redshift/driver` draft, PR #1371 | thin adapter after recording; `Dump` stays a no-op like Bytebase's | session, plancheck, unverified | unmeasured (draft driver) |
| 22 | `SNOWFLAKE` | `db/snowflake` | own | `snowflake/driver` draft, PR #1378 | thin adapter after recording | session, unverified | unmeasured (draft driver) |
| 23 | `BIGQUERY` | `db/bigquery` | own; parser handled by `googlesql` | `bigquery/driver` draft, PR #1377 | thin adapter over `bigquery/driver`, not the `googlesql` alias | unverified; no tape seam yet | unmeasured (draft driver) |
| 24 | `DATABRICKS` | `db/databricks` | own; parser is `standard` | `databricks/driver` draft, PR #1379 | thin adapter after recording | unverified; no `catalog`, so its `Dump` is server DDL | unmeasured (draft driver) |

Blockers:

- **session**: Bytebase's editor opens a `*sql.Conn` from `GetDB` and reuses it across statements
  (`api/v1/sql_service.go:117`, `:404`, `:1347`). The adapter has no `*sql.DB` to hand back, so a
  cutover must hold a `driver.Session` per editor session, or accept one session per call. Engines
  whose `QueryConn` ignores `conn` already behave as the second.
- **plancheck**: `runner/plancheck/statement_report_executor.go:168-220` type-asserts
  `*pgdriver.Driver`, `*redshiftdriver.Driver`, `*cockroachdbdriver.Driver`,
  `*mysqldriver.Driver` (for MySQL and MariaDB), `*oceanbasedriver.Driver`, `*tidbdriver.Driver`,
  `*oracledriver.Driver`, `*mssqldriver.Driver` to reach `CountAffectedRows`. It must name an
  interface instead. **Measured:** no engine has needed this yet. mysql and oracle both landed with
  the assertions intact, because an adapter that keeps its concrete type keeps `CountAffectedRows`
  reachable; the change buys the removal of the type pin, not a compile. When someone does collapse
  them, the nine parser imports in §5 are part of the price.
- **advisor**: the V1 live rules take `advisor.Context.Driver *sql.DB`
  (`plugin/advisor/advisor.go:51`). `pg`, `mysql`, `mariadb`, `oceanbase`, `tidb` and `oracle` have
  rules that EXPLAIN or query through it. A nil `GetDB` makes them silently skip, not fail.
- **ghost**: MySQL's gh-ost validator and cleanup reach `GetDB` directly
  (`component/ghost/validator.go:54`, `runner/taskrun/database_migrate_executor.go:369`).
- **version**: `runner/taskrun/database_migrate_executor.go:827` casts the Oracle driver and calls
  `GetVersion() *plsql.Version`; `oracle/driver` exposes `Version(ctx)`. **Measured:** the cutover
  deleted the problem instead of solving it. omni reads `version_full` and falls back to
  `Config.EngineVersion`, so nothing calls a version method any more; the one call site parses
  `instance.Metadata.GetVersion()` directly and the `*oracle.Driver` assertion goes away (−8/+9).
  A product change, but a deletion of an engine-specific poke, not a new seam.

## 3. Measured cost

Five cutovers have been built end to end in scratch Bytebase trees, against live servers, with no
Bytebase commit: trino at omni `bf1a3396e`, pg at `8664fa7b1`, cassandra at `b13c2e19b`, and
mysql and oracle at `6dee22d57`. Each replaces `backend/plugin/db/<engine>` with an adapter, then
makes the product-side change the adapter forces in the scratch tree only, so every failure is
classified as an adapter bug, a product-side change, or an omni gap. **No proof found a contract
gap**: every forced change was on Bytebase's side of the seam. pg was run against six majors, mysql
against two pinned images and oracle against `gvenzl/oracle-free:23.26.3-slim-faststart`, one
container at a time; each proof's per-file diff and verbatim failures are in its session's report,
and this section keeps the totals and the decisions they force.

The one omni change any proof asked for is a package doc: `cassandra/driver`'s `SyncDatabase`
promised an empty schema where the code returns `driver.ErrNotFound`. That proof fixed the doc on
its own branch (`5fcf51383`), and it is not part of this note.

mssql is still being proved in a scratch tree, so its row stays `unmeasured`; so do rows 2-8,
10-12, 14, 17-19 and 21-24, whose §2 estimate is the adapter prototype's, not a cutover's.

| Engine | Prototype adapter | Measured adapter, non-test | Deleted | Callers changed | Product-side diff |
| --- | --- | --- | --- | --- | --- |
| trino | 101 | 234 (166 code) + 120 shared conversion | 665 driver + 169 tests | **0** | 2 result-shape expectations in the package's own tests; `plugin/schema/trino`: 4 dump strings, +14 lines for views |
| pg | 226 | 658 (the delegation layer inside it is ~219) | package 5,376 → 1,741; `sync.go`, `role.go`, `stats.go` gone; 2,043 test lines across 3 files | 8 for `GetDB`, 7 for the advisor probe | +199 / −45 (net +154) across 14 product files and 4 test files |
| mysql | 136 | 706 (450 core, 68 retained helpers, 155 legacy `*sql.DB` probe bridge, 33 synonyms) | 2,952 → 706 non-test; 5 driver files (2,196 lines); 3 test files, 2,163 lines, no longer compile | 17 files | +363 / −104 (net +259), plus 9 blank parser imports |
| cassandra | 126 | 312 (one file) | 802 driver + 126 tests (928 lines, net −507) | **0** | **0 files.** `go.mod`'s pin replace and `ultimate.go`'s blank import are the only two lines outside the package |
| mssql | 188 | unmeasured | unmeasured | unmeasured | unmeasured |
| oracle | 143 | 390 (212 `oracle.go`, 139 `query.go`, 26 `sync.go`, 13 `dump.go`) + the same 120 shared conversion | package 2,672 → 390; 2,526 lines deleted; tests 1,435 → 843 | 12 files (+ 6 test files touched only to compile) | +243 / −36 (net +207) |

| Engine | Did the prototype's verdict hold? | Therefore, do this first |
| --- | --- | --- |
| trino | Half. "Clean fit, no caller change" held exactly: `go build ./backend/...` green with no caller edited. The size did not: 234 against 101, because the prototype omitted the `ExecuteLogger` bridge and the command-range fix, and it did not count the shared conversion trino, mysql and oracle each wrote (120 lines). It also could not see that the pin bump, not the adapter, moves the dump text. | Land the shared `db/util/omnidriver.go` conversion once, then run the engine's dump/diff fixture suite before writing adapter code. The pin bump is the part that breaks the product, and it is cheaper to discover from a fixture than from an integration test. |
| pg | Half. 226 was the delegation layer (measured ~219); the adapter is 658 because connection resolution Bytebase still owns (~99), `driver.Result` → `v1pb` conversion (~168) and the advisor probe (~75) are part of it. The forced interface was the **advisor probe**, not plancheck, contrary to the prototype's framing. | `advisor.Prober` and `Context.Prober` first (+77/−19 across the new `prober.go`, `advisor.go` and five pg rule files), and the rest of the cutover then compiles. `GetDB` retirement and the 3 `sql_service` sites follow; `rollout_service` (+63/−3) is the one rewrite, onto omni's `Activity()`. |
| mysql | No. 136 → 706, a 5× miss, and the seam list held in kind but not in shape: a generic probe could not answer MySQL's rules (§5), and 2,163 lines of package tests stopped compiling — a cost no compile can show. | `DriverHandle` for the mysql/mariadb rules first, because the generic prober does not fit them; then gh-ost's fail-closed validator, the only silent feature loss in any proof. If wave 2 has not replaced plancheck's assertions by then, do it here (+24/−62 plus nine imports) — the mysql proof did, and would again, because it removes the pin that keeps `*mysqldriver.Driver` alive. |
| cassandra | Yes on the claim that matters, no on size. Every wave-1 precondition held exactly — `GetDB` nil before and after, `QueryConn` ignoring its `*sql.Conn`, no product package naming `*cassandra.Driver` — and it is the only measured cutover with zero product files. 312 against 126, and it found four behavior deltas the plan does not name. | Land it first: this is the one cutover a reviewer can approve as a driver swap plus a deletion. Decide its four deltas as part of the change — a dropped keyspace is now `driver.ErrNotFound` instead of an empty schema, `Execute` now emits task-run logs, the adapter gains an SSH `Dial`, and the dead read-only branch is preserved. |
| mssql | unmeasured | Budget the plancheck interface and its nine imports (wave 2's entry) as its structural blocker, plus whatever its masking consumers need, which is what its proof is measuring. |
| oracle | Mostly, in the opposite direction. 143 → 390 adapter lines, but the two blockers the plan named as product work mostly were not: `CountAffectedRows` needed **zero** product change, plancheck was never touched, and the `Version` seam disappeared because omni reads `version_full` itself, leaving an 8-line deletion. What did hold is the advisor `Prober`, reused verbatim from pg, and one blocker no one predicted: the private-link guard's interface loses its `*sql.Conn`. | Reuse pg's `prober.go` and `Context.Prober` without re-deriving them, delete the `*oracle.Driver` version assertion instead of porting it, and budget the private-link guard's signature change — it is the one oracle-specific caller edit. |

Every measured row pays these, so every unmeasured row should budget them too:

- **The shared conversion.** 120 lines, written separately by trino, mysql and oracle (`db/util/omnidriver.go`,
  ~102 code). It is not engine cost; landing it once (wave 0) removes it from every adapter that
  does not inline it, and cassandra's 312-line single file is what inlining it looks like.
- **Test rewiring.** trino lost 169 test lines, pg 2,043, mysql 2,163, oracle 857 (592 net). A
  cutover that compiles is not a cutover: the tests that construct `&Driver{db: *sql.DB}` are part of
  the price.
- **The pin bump.** The omni pin moves when the adapter lands, and the pin bump changes product output:
  trino's dumps gained `CREATE SCHEMA IF NOT EXISTS`, views entered the diff SDL (+14 lines), and six
  result shapes changed (DECIMAL as text, timestamp accuracy, `timestamp_tz`, no synthetic
  `Affected Rows 0`, `information_schema` gone from sync, no row cap where the engine cannot rewrite).
  Each engine needs its own fixture pass; none of it is visible from `go build`.
- **The session change is not a caller change.** trino changed 0 callers because `db.SQLDB` returns nil
  and the three `sql_service` sites already guard it; pg keeps a session per driver through omni's own
  `Activity()`. The behavior that changes is one session per call, not a compile error.

## 4. Where the two pictures disagree

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
**Judgement call.** Measured since: no cutover has needed the plancheck change yet (§2), so the
split's cost is the import and the assertion, not a compile.

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

## 5. What the measurements say the plan got wrong

The §3 numbers correct the plan's degrees. Two things correct its shape, and neither was visible
from a prototype that compiled.

### A generic `Prober` is not enough

The plan said the advisor's `*sql.DB` becomes "a probe interface" and the prototype measured a
46-line `Prober`. That holds for pg: `ProberOf(driver)` hands the constant `Prober` an omni session,
the probe runs `EXPLAIN (ANALYZE …)` and rolls back, and five pg rule files change by +22/−16,
including +12/−5 for the object-owner rule.

MySQL is where it breaks. Its affected-row and insert-row rules do not ask a session a question;
they ask the driver's AST-aware `ExplainJSON`/`EstimateAffectedRows`, which is engine code with a
different signature and a different answer. A generic `Prober` cannot reach it, so the cutover grew
a second seam: an engine-specific `DriverHandle` on `advisor.Context` (`advisor.go` +11/−1) with the
legacy `*sql.DB` still reachable, and the four mysql/mariadb rule files grew +73/−20. The generic
prober stayed in the tree for pg; the mysql path does not use it.

So the plan budgets **two** advisor seams and stops treating the advisor work as one change: the
generic `Prober` arrives with pg, the `DriverHandle` with mysql/mariadb, and because both engines
live in one package the MariaDB rules pay the same edit twice (29 of the 73 lines).

### plancheck's assertions were hiding parser registration

`runner/plancheck/statement_report_executor.go` type-asserted 7 concrete drivers covering 8 engine
values to reach `CountAffectedRows`. Replacing them with one interface is +24/−62 — but those
imports were also what registered parsers. Each driver package's `init` chain pulls in the engine's
parser, and `TestSummaryStatementTypesCoverEveryReportEngine` started failing with `engine POSTGRES
is not supported for ParseStatements`, and the same for `COCKROACHDB`. Nine explicit blank parser
imports restore it.

The interface change is therefore +24/−62 **plus nine imports**, and the coupling is invisible from
plancheck alone. The measurements draw two consequences. It belongs with wave 2, not with the shared
layer: no engine measured so far needed it — mysql and oracle both landed with the assertions intact
— so doing it early spends the nine imports before anything needs the pin gone. And once it is done,
never tidy those blank imports: they are load-bearing, and omni's own imports cannot replace them,
because the assertion table is over Bytebase's parser registry.

### Three smaller corrections

- **The session blocker is not a wave gate.** trino changed zero callers and pg's `*sql.Conn`
  parameter is dead but tolerated, because `QueryConn` keys its own omni session per driver. What
  remains is the `db.SQLDB` shim at 7 call sites, and it disappears one engine at a time.
- **gh-ost loses a probe.** `component/ghost/validator.go` (+13/−5) now fails closed: the contract
  has no binlog-status or grant probe, so a shipped MySQL feature turns off instead of erroring at
  cutover. It is the only silent regression any proof found, and a gap the next contract revision
  should name.
- **Compile is not evidence.** pg's 2,043 and mysql's 2,163 lines of test that no longer compile are
  the largest single line item outside the adapters, and no `go build` sees them.

## 6. Order

Waves are by measured cost and by forced change, not by engine importance. The numbers in the
parentheses are the §2 rows, kept so a reader can see which prediction moved and which did not.

**Wave 0 — the shared caller layer (no engine, and not a gate).** Measured across the proofs and
identical wherever it appears: `db/util/omnidriver.go` (120 lines, ~102 code), which trino, mysql
and oracle each paid separately, so budget it once; `db.Driver` losing `GetDB` for `SQLDBProvider`
+ `db.SQLDB` (+25/−1); `advisor.Prober` + `Context.Prober` (+46 and +11/−1); and the MySQL-family
config builder the prototype measured at 158 lines (82 of them the RDS CA bundle). About +202/−2
plus the builder. Cassandra's proof needed none of it, which is the useful measurement: wave 0 does
not gate wave 1, and the first adapter can land before it. Land it before the second wave-1 adapter
— trino and mysql wrote the same 120 lines twice — and before wave 3, where pg's `Prober` and
oracle's reuse of it verbatim are the whole point.

**Wave 1 — adapters with no forced product change (1-12).** `cassandra` first: 312 adapter lines, 0
callers, 0 product files, 928 lines deleted, a live pass on `cassandra:5.0.9` — a driver swap plus a
deletion, and nothing else for a reviewer to read. `trino` second: 234 lines, 0 callers, and the
first engine to exercise the fixture path (4 dump strings, +14 lines for views, 2 result-shape
expectations) that every later engine repeats. Then `cosmosdb`, `dynamodb`, `mongodb`, `redis`,
`spanner` in any order, each with the §3 fixture pass. `clickhouse`, `doris` and `starrocks` follow
once wave 0's shared builder exists: they appear in no plancheck, advisor or ghost switch, and
trino showed the session seam costs no caller. `hive` and `elasticsearch` keep their documented
limits (local port; SigV4 transport).

Cassandra-first survives, on an argument the plan did not make. The plan's reason was size — 126
lines, the smallest prototype — and the size was wrong: 312, a 2.5× miss. What holds is the clause
the plan put second: cassandra is the only wave-1 engine that needs no engine-specific `Open`
options, no `Params` and no `Options`, and it measured zero product files, which is why 126 was the
least important number in its row. `cosmosdb` (510 non-test lines in Bytebase) is smaller to delete
but needs `Params` auth/tenantId/emulator, the emulator TLS path and `QueryOptions.Schema` as the
container; `redis` (521) needs topology `Options` and the `AllNodesQuerier` probe. So the order does
not change and the reason does: cassandra is first because it is clean, not because it is small.

**Wave 2 — plancheck-only engines (13, 14).** `mssql` and `cockroachdb`: the plancheck interface is
their entry cost, and it is the only place the nine blank parser imports appear (+24/−62 plus the
imports). No measured engine has needed it, so it is wave 2's work rather than shared-layer work.
mssql also carries its masking consumers, which its proof is measuring; `cockroachdb` stays its own
adapter rather than aliasing pg (§4). Both are unmeasured, so budget below pg's 658 adapter lines
and above trino's 234.

**Wave 3 — the SQL family with advisor coupling (15-20), in this order: `pg`, `oracle`,
`tidb`/`oceanbase`, `mysql`/`mariadb`.** pg first: measured, 658 adapter lines, 14 product files,
net +154, and its `Prober` seam is the one the other engines reuse. Oracle second because it is
measured and nearly the same shape — 390 adapter lines, 12 files, net +207 — with two blockers gone
(no version seam, no plancheck) and one gnarly call site (the private-link guard's signature), so it
is the wave's cheap follow. tidb and oceanbase sit next: they carry the advisors but not ghost, and
they reach the shared builder. mysql and mariadb last because they are measured worst: 706 adapter
lines, 17 files, net +259, the second advisor seam (`DriverHandle`), and gh-ost failing closed. Do
not start mariadb before mysql: they are one package and pay the same rules twice.

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

## 7. Start the first cutover

1. Land wave 0 as its own change, but not as a gate on the first adapter: the shared conversion,
   the `db.SQLDB` seam and the advisor `Prober`. plancheck's interface and its nine imports belong
   to wave 2, when mssql and cockroachdb are the engines that need the pin gone. Cassandra's proof
   needed none of it; every adapter after it wants the 120 conversion lines and, from wave 3 on,
   the `Prober`.
2. Pick `CASSANDRA` (row 1). It is the first cutover measured at zero product files: one adapter
   file, 312 non-test lines, 928 deleted, 0 callers, and a live pass on `cassandra:5.0.9` (Ping,
   SyncInstance, SyncDBSchema, Execute, QueryConn, Dump) plus omni's tape suite. It registers one
   engine in `db/cassandra/cassandra.go`, its `GetDB` returns nil, its `QueryConn` ignores the
   `*sql.Conn`, and no product package names its concrete type.
3. Keep the package path and `db.Register(storepb.Engine_CASSANDRA, …)`. Replace the body of
   `Driver` with an adapter over `omni/cassandra/driver.Open`. Do not touch `db.Open`, the map, or
   `backend/server/ultimate.go`. Decide the four measured deltas in the same change: a dropped
   keyspace answers `driver.ErrNotFound` instead of an empty schema, `Execute` now emits task-run
   logs, the adapter gains an SSH `Dial`, and the dead read-only branch is preserved.
4. Convert `storepb.DataSource` and `ConnectionContext` into `driver.Config` in the adapter:
   host/port, database, user/password, `TLS func(host)`, `Dial` for the SSH tunnel, `Params`,
   application name, read-only, engine version. Cache nothing the old driver did not cache.
5. Reuse wave 0's conversion at the seam if it has landed; cassandra inlined it inside its 312
   lines, and that is the fallback if it has not. Keep the `ExecuteLogger` bridge per engine: the
   trino proof needed it (with the command-range fix) and the prototype did not count it.
6. Run the fixture pass before the integration pass: the engine's dump and diff expectations under
   `plugin/schema/<engine>`, then Bytebase's engine tests against the pinned image (one container at
   a time, removed before the next), then omni's tape suite. Expect the dump text and some result
   shapes to change with the pin bump, and decide each one before the cutover lands.
7. Then take `TRINO` (row 9) with the same shape, which is where the fixture path is exercised for
   the first time, then the rest of wave 1 one package at a time. The other engines pay the same
   fixture pass; none of them should be the one that first discovers it.

## 8. What stays product-side

The cutover does not move these, and a proposal that does is a different change:

- The `storepb.Engine` enum, `db.Register`/`db.Open`, the blank imports, and `dbfactory`.
- Data-source resolution, external secret replacement, AWS/GCP IAM tokens, SSH tunnels, TLS
  material, and the Cloud SQL / RDS dialers.
- The parser and schema registries, `convertToParserEngine`, `common/engine.go` capability lists,
  `db/explain_format.go`, advisory registration, plancheck statement typing, and the SQL editor's
  session policy.
- The task-run logger renderer. `driver.ExecuteLogger` is the bridge; `storepb.TaskRunLog` and the
  command range are product.

The measurements add three entries to that list, and each is where a future cutover will look for
them:

- `db.SQLDB` / `SQLDBProvider` and plancheck's nine blank parser imports. Both are shape, not
  plumbing: one is how an engine that still has a `*sql.DB` answers, the other is how the assertion
  table keeps its parsers.
- `advisor.Prober` and `advisor.DriverHandle`. Two seams, because a generic session probe answers
  pg's question and not mysql's (§5).
- gh-ost's binlog probe. The contract has no counterpart, so `component/ghost/validator.go` fails
  closed; a contract change that adds the probe is the only way to keep the feature.

The shared MySQL-family config builder (the prototype measures 158 lines, 82 of them the RDS CA
bundle) stays in Bytebase, because six adapters reach it and no engine owns it. Put it in one
package beside the adapters rather than copying it six times; wave 0 is where that happens.
