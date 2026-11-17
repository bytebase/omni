# Cutover dispatch: how each engine reaches its omni driver

2026-11-17 · Proposed · revised with measurements

Bytebase `345007290e`, omni `f0cbe6dc4`. This note is the dispatch half of the driver move:
[2026-11-15-driver-contract-design.md](2026-11-15-driver-contract-design.md) fixes the contract and
[2026-11-17-driver-adapter-prototype.md](2026-11-17-driver-adapter-prototype.md) measures what an
adapter body costs. Neither says how the twenty-four real `storepb.Engine` values reach twenty
driver trees on main through Bytebase's registry, which engines share code today, or which engine
can be flipped first without a product change. That is this note.

Revised after nine cutovers were built: §3 replaces the adapter predictions with measured costs, §5
records the two places a prototype could not see the shape of the work, §6 separates what a live
server found from what a compile found, and §7 reorders the waves. The twenty-four rows in §2 keep
the order the prototypes proposed, so a reader can see what moved; each row carries a measured value
or says it has none.

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
cutover changes one package and no product code. The rest need the product work in §7.

## 2. The picture today

`tapes` means the package and its recorded tests are on `origin/main`; `draft` means the package
exists only on the named open branch and has never been recorded against the engine. Cutover order
is the number in the first column; the waves are explained in §7. The last column carries the
measured adapter lines and product-side file count from §3; `unmeasured` means the prototype's
estimate is all there is, and the row is a prediction, not a cost.

| # | Engine | Bytebase package | Shares / aliases today | omni driver | Recommendation | Blockers | Measured |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 1 | `CASSANDRA` | `db/cassandra` | own | `cassandra/driver` (tapes) | thin adapter | measured: none. `GetDB` returns nil, `QueryConn` ignores its `*sql.Conn`, no product file names the driver type | 312 / 0 callers, 0 files |
| 2 | `COSMOSDB` | `db/cosmosdb` | own | `cosmosdb/driver` (tapes) | thin adapter; `Params` auth/tenant/emulator, `QueryOptions.Schema` is the container | none | unmeasured |
| 3 | `DYNAMODB` | `db/dynamodb` | own; PartiQL is its language tree | `dynamodb/driver` (tapes) | thin adapter; `Params` region/sessionToken/accountId | none | unmeasured |
| 4 | `MONGODB` | `db/mongodb` | own; omni tree is named `mongo` | `mongo/driver` (tapes) | thin adapter; `Open` takes `Options{SRV}` | none | unmeasured |
| 5 | `REDIS` | `db/redis` | own | `redis/driver` (tapes) | thin adapter; `Options` topology, and the session probe for all-nodes | measured: none. Both seams it might have wanted are inapplicable — `SQLDB` because `GetDB` was already nil, `Prober` because the all-nodes choice is inside `QueryConn` — and the probe arrived as `driver.AllNodesQuerier`. The pin predates `redis/driver` entirely, so nothing lands without a version bump | 283 / 0 callers, 0 files |
| 6 | `SPANNER` | `db/spanner` | own; parser handled by `googlesql` | `spanner/driver` (tapes) | thin adapter; `Params` project/instance | measured: none. `Prober`, `SQLDB`, `Activity` and the private-link guard are each inapplicable, for the reasons §6 records, and the server-DDL dump is byte-identical for both dialects. One product decision is open: the driver now lists PostgreSQL-dialect databases, which the product's parser does not route (§4) | 333 / 0 callers, 0 files |
| 7 | `HIVE` | `db/hive` | own; parser is `standard` | `hive/driver` (tapes) | thin adapter | no caller `Dial` (forward a local port) | unmeasured |
| 8 | `ELASTICSEARCH` | `db/elasticsearch` | own | `elasticsearch/driver` (tapes) | thin adapter; `Options{OpenSearch, Transport}` | measured: neither seam blocks — a caller `Dial` reaches the client and `Options.Transport` carries SigV4, both by test. Ten behaviour changes, and one of them (AWS ping hints) is a product decision | 347 / 0 callers, 0 files |
| 9 | `TRINO` | `db/trino` | own | `trino/driver` (tapes) | thin adapter | measured: none. No caller changed; `QueryConn` ignores its `*sql.Conn`, so a session is per call | 234 / 0 callers |
| 10 | `CLICKHOUSE` | `db/clickhouse` | own; parser is `standard`/`mysql` | `clickhouse/driver` (tapes) | thin adapter; `Dumper` now exists | session, as trino's (per call, not a caller change) | unmeasured |
| 11 | `DORIS` | `db/doris` | own; imports `db/mysql` helpers | `doris/driver` (tapes) | thin adapter + shared MySQL config builder | session; the shared builder is a precondition | unmeasured |
| 12 | `STARROCKS` | `db/starrocks` | own; imports `db/mysql` helpers | `starrocks/driver` (tapes) | thin adapter + shared MySQL config builder | session; the shared builder is a precondition | unmeasured |
| 13 | `MSSQL` | `db/mssql` | own | `mssql/driver` (tapes) | thin adapter + `CountAffectedRows` | measured: the shared `GetDB` retirement only. No advisor probe, no plancheck change; masking and batch delegate. Azure AD IAM is refused loudly | 497 / 8 files |
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
  interface instead. **Measured:** no engine has needed this yet. mssql, mysql and oracle all landed
  with the assertions intact, because an adapter that keeps its concrete type keeps
  `CountAffectedRows` reachable; the change buys the removal of the type pin, not a compile. When
  someone does collapse them, the nine parser imports in §5 are part of the price.
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

Nine cutovers have been built end to end in scratch Bytebase trees, against live servers, with no
Bytebase commit: trino at omni `bf1a3396e`, pg at `8664fa7b1`, cassandra at `b13c2e19b`, redis at
`b515e82e2`, spanner at `b85cd8c9b`, elasticsearch at `af2e1a130`, and mysql, oracle and mssql at
`6dee22d57` (mssql on top of the driver fixes in #1385). Each replaces
`backend/plugin/db/<engine>` with an adapter, then makes the product-side change the adapter forces
in the scratch tree only, so every failure is classified as an adapter bug, a product-side change,
or an omni gap. **No caller change was omni's to make**: every product-side change landed on
Bytebase's side of the seam. pg was run against six majors, mysql against two pinned images, oracle
against `gvenzl/oracle-free:23.26.3-slim-faststart`, redis against `redis:7.4.12-alpine`, spanner
against the emulator digest in its `pins.tsv` row, and elasticsearch against 8.19.25 with security
on, one container at a time; each proof's per-file diff and verbatim failures are in its session's
report, and this section keeps the totals and the decisions they force.

The omni side needed work twice, none of it contract-shaped. `cassandra/driver`'s package doc
promised an empty schema where the code returns `driver.ErrNotFound` (fixed on that proof's branch,
`5fcf51383`). mssql's live suite found five driver defects, enumerated in §6, all fixed in #1385.
Redis, spanner and elasticsearch asked nothing of omni at all — no (c)-class item and no PR between
the three, which is the first time three engines in a row tested the driver and left it unedited.
One capability resisted: `mssql/driver` always builds the plain sqlserver connector, so a caller
cannot hand it Azure AD's token-fetching connector, and the adapter refuses that mode loudly. That
needs a contract decision, and it is the only shipped feature any proof could not carry.

The remaining rows stay `unmeasured`: rows 2-4, 7, 10-12, 14, 17-19 and 21-24, whose §2 estimate is
the adapter prototype's, not a cutover's. Doris, MariaDB and DynamoDB were still running when this
revision was written, so their rows keep the prototype's number and nothing here is inferred from a
neighbouring engine — not from doris to starrocks, not from mysql to mariadb, not from a driver's
tapes to a cutover.

| Engine | Prototype adapter | Measured adapter, non-test | Deleted | Callers changed | Product-side diff |
| --- | --- | --- | --- | --- | --- |
| trino | 101 | 234 (166 code) + 120 shared conversion | 665 driver + 169 tests | **0** | 2 result-shape expectations in the package's own tests; `plugin/schema/trino`: 4 dump strings, +14 lines for views |
| pg | 226 | 658 (the delegation layer inside it is ~219) | package 5,376 → 1,741; `sync.go`, `role.go`, `stats.go` gone; 2,043 test lines across 3 files | 8 for `GetDB`, 7 for the advisor probe | +199 / −45 (net +154) across 14 product files and 4 test files |
| mysql | 136 | 706 (450 core, 68 retained helpers, 155 legacy `*sql.DB` probe bridge, 33 synonyms) | 2,952 → 706 non-test; 5 driver files (2,196 lines); 3 test files, 2,163 lines, no longer compile | 17 files | +363 / −104 (net +259), plus 9 blank parser imports |
| cassandra | 126 | 312 (one file) | 802 driver + 126 tests (928 lines, net −507) | **0** | **0 files.** `go.mod`'s pin replace and `ultimate.go`'s blank import are the only two lines outside the package |
| mssql | 188 | 497 (359 `mssql.go`, 124 `query.go`, 14 `dump.go`) | 2,952 non-test (5 files); 40 test lines deleted, `mssql_test.go` 535 → 500 | 8 files (all `GetDB`; 3 test files touched only to compile) | +50 / −16 (net +34). No advisor seam, no plancheck change. Its `Dump` now renders through omni's catalog, so a filtered index keeps `WHERE ([amount] IS NOT NULL)` and the doubled `;` after a view definition is gone — 921 → 948 bytes on the same synced schema, a user-visible change no test pins yet |
| oracle | 143 | 390 (212 `oracle.go`, 139 `query.go`, 26 `sync.go`, 13 `dump.go`) + the same 120 shared conversion | package 2,672 → 390; 2,526 lines deleted; tests 1,435 → 843 | 12 files (+ 6 test files touched only to compile) | +243 / −36 (net +207) |
| redis | 154 | 283 (one file, ~120 of it result conversion and probe dispatch) | `sync.go` 148 deleted, `redis.go` 373 → 283; package 521 → 283 non-test; tests 191 → 205 | **0** | **0 files.** The pin bump *is* the cutover: Bytebase's pin (`8851de182`) predates `redis/driver` (`9be754a6`), so the driver enters the dependency for the first time; the testcontainer image moves to `redis:7.4.12-alpine` |
| spanner | 120 | 333 (one file, `spanner.go`: config mapping, result conversion, `executeLogger`) | `sync.go` 371 and `dump.go` 28 deleted; driver 1,127 → 333 (−70.5%), package 1,553 → 809; 314 test lines deleted, 2 live tests ported | **0** | **0 files.** `go.mod`'s replace is the only line outside the package. The server-DDL dump is byte-identical for both dialects, with and without the sync runner's metadata |
| elasticsearch | 169 (+11 SigV4) | 347 (one file: `Open`/dial, result mapping, execute logger, sync mapping, dump) | `sync.go` 429, `query_plan.go` 113, `role.go` 57, `dump.go` 13 deleted; driver 1,180 → 347 (−833); the package's own tests 405 → 541, plus a shared container fixture (+60) | **0** | **0 product files**, and two dependencies leave Bytebase: `elastic/go-elasticsearch/v7` and `rs/dnscache` (go.sum −16) |

| Engine | Did the prototype's verdict hold? | Therefore, do this first |
| --- | --- | --- |
| trino | Half. "Clean fit, no caller change" held exactly: `go build ./backend/...` green with no caller edited. The size did not: 234 against 101, because the prototype omitted the `ExecuteLogger` bridge and the command-range fix, and it did not count the shared conversion trino, mysql and oracle each wrote (120 lines). It also could not see that the pin bump, not the adapter, moves the dump text. | Land the shared `db/util/omnidriver.go` conversion once, then run the engine's dump/diff fixture suite before writing adapter code. The pin bump is the part that breaks the product, and it is cheaper to discover from a fixture than from an integration test. |
| pg | Half. 226 was the delegation layer (measured ~219); the adapter is 658 because connection resolution Bytebase still owns (~99), `driver.Result` → `v1pb` conversion (~168) and the advisor probe (~75) are part of it. The forced interface was the **advisor probe**, not plancheck, contrary to the prototype's framing. | `advisor.Prober` and `Context.Prober` first (+77/−19 across the new `prober.go`, `advisor.go` and five pg rule files), and the rest of the cutover then compiles. `GetDB` retirement and the 3 `sql_service` sites follow; `rollout_service` (+63/−3) is the one rewrite, onto omni's `Activity()`. |
| mysql | No. 136 → 706, a 5× miss, and the seam list held in kind but not in shape: a generic probe could not answer MySQL's rules (§5), and 2,163 lines of package tests stopped compiling — a cost no compile can show. | `DriverHandle` for the mysql/mariadb rules first, because the generic prober does not fit them; then gh-ost's fail-closed validator, the only silent feature loss in any proof. If wave 2 has not replaced plancheck's assertions by then, do it here (+24/−62 plus nine imports) — the mysql proof did, and would again, because it removes the pin that keeps `*mysqldriver.Driver` alive. |
| cassandra | Yes on the claim that matters, no on size. Every wave-1 precondition held exactly — `GetDB` nil before and after, `QueryConn` ignoring its `*sql.Conn`, no product package naming `*cassandra.Driver` — and it was the only measured cutover with zero product files when it reported; redis, spanner and elasticsearch have since matched that. 312 against 126, and it found four behavior deltas the plan does not name. | Land it first: this is the one cutover a reviewer can approve as a driver swap plus a deletion. Decide its four deltas as part of the change — a dropped keyspace is now `driver.ErrNotFound` instead of an empty schema, `Execute` now emits task-run logs, the adapter gains an SSH `Dial`, and the dead read-only branch is preserved. |
| mssql | Yes on cost, no on the seam list. The predicted blockers mostly were not: no `Prober` (mssql's advisors are static AST rules), no plancheck change (the `*mssql.Driver` assertion still holds), and all 8 caller files are the shared `GetDB` retirement, none mssql-specific. What broke was driver-side — five defects that only a live suite finds — and its `Dump` text moved without a test noticing. 497 against 188. | Reuse wave 0's `db.SQLDB` and change nothing engine-specific on the caller side. #1385's five driver fixes are on main (`81b91221d`) and are a precondition, not a follow-up, because the adapter is correct only on top of them; record the Azure AD IAM refusal as a contract decision rather than an adapter shortcut. |
| oracle | Mostly, in the opposite direction. 143 → 390 adapter lines, but the two blockers the plan named as product work mostly were not: `CountAffectedRows` needed **zero** product change, plancheck was never touched, and the `Version` seam disappeared because omni reads `version_full` itself, leaving an 8-line deletion. What did hold is the advisor `Prober`, reused verbatim from pg, and one blocker no one predicted: the private-link guard's interface loses its `*sql.Conn`. | Reuse pg's `prober.go` and `Context.Prober` without re-deriving them, delete the `*oracle.Driver` version assertion instead of porting it, and budget the private-link guard's signature change — it is the one oracle-specific caller edit. |
| redis | Yes on the seams, no on size — the reverse of the usual miss. 154 → 283, but both seams it might have wanted (`SQLDB`, `Prober`) were inapplicable, and the all-nodes probe came through `driver.AllNodesQuerier` with a mutation check. What no estimate could see is the dependency: the pin predates `redis/driver`, so this cutover is also a first install. | Bump the pin first and keep the all-nodes branch. A silent fallback to one session's `Query` passes a single-node cluster, so the test that matters is the refusal on a standalone server, not the row count. Accept one round trip per command (200 `PING`s: 2.4 ms → 13.8–23.4 ms on localhost) as the price of fail-fast `Execute` and per-command errors. |
| spanner | Half. 120 → 333, but zero seams, zero product files, and the dump byte-identical with and without the sync runner's metadata. What it could not see is that a *capability* arrives with the driver — PostgreSQL-dialect databases, which the old driver filtered — so the expensive part is a product decision, not adapter code. | Decide the dialect listing before writing the adapter. omni lists, syncs and dumps PG-dialect databases; `SPANNER`'s `SplitMultiSQL`/`GetQuerySpan` is GoogleSQL-only, so keep the old filter in `SyncInstance` until the product routes by database dialect. Ship the two halves together or not at all. |
| elasticsearch | Mostly. 169 → 347, 0 callers, 0 seams, and both seam questions answered by test: a caller `Dial` reaches the client, `Transport` carries SigV4. What it could not see is how many behaviours move at once (ten), and that one of them is support-facing: the AWS ping/sync error text loses its role-ARN hint. | Run the old driver beside the adapter on one pinned server before deciding anything; the ten changes only exist as a side-by-side diff, and three shapes expected to move (`hits`, `_cat`, refused-request text) came back identical. Decide the AWS hint first — it is the only one a support case would notice. |

Every measured row pays these, so every unmeasured row should budget them too:

- **The shared conversion.** 120 lines, written separately by trino, mysql and oracle (`db/util/omnidriver.go`,
  ~102 code), and written again inside redis (about 120 of its 283 lines, probe dispatch included),
  spanner (`toQueryResults`/`planOf`/`toRowValue`) and elasticsearch (result mapping). It still does
  not exist on Bytebase `345007290e`, so six proofs have paid for it separately. It is not engine
  cost; landing it once (wave 0) removes it from every adapter that does not inline it, and
  cassandra's 312-line single file is what inlining it looks like.
- **Test rewiring.** trino lost 169 test lines, pg 2,043, mysql 2,163, oracle 857 (592 net). A
  cutover that compiles is not a cutover: the tests that construct `&Driver{db: *sql.DB}` are part of
  the price. Spanner is the one cutover where test debt fell — 314 lines deleted because omni's
  `spanner/driver_test.go` already held the same three tests under the same names. Elasticsearch
  went the other way, 405 → 541, once its pinned server became a shared container fixture.
- **The pin bump.** The omni pin moves when the adapter lands, and the pin bump changes product output:
  trino's dumps gained `CREATE SCHEMA IF NOT EXISTS`, views entered the diff SDL (+14 lines), and six
  result shapes changed (DECIMAL as text, timestamp accuracy, `timestamp_tz`, no synthetic
  `Affected Rows 0`, `information_schema` gone from sync, no row cap where the engine cannot rewrite).
  Each engine needs its own fixture pass; none of it is visible from `go build`. Two ends of the
  range: redis's pin predates `redis/driver` entirely, so its bump introduces the dependency; and
  where an engine renders no schema the bump moves nothing — spanner's dump is byte-identical for
  both dialects, redis's and elasticsearch's are empty before and after, and the syncer's history
  says so.
- **The driver, not the adapter, is where behavior breaks.** mssql's five defects and spanner's
  credential/emulator conflict (§6) are findings no compile can produce, and neither can a review of
  the adapter body. The one exception is oracle's private-link guard, which is a product compile
  error and therefore the good kind. Budget a live run per engine, or the cutover ships on a build.
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
with tapes and was cut over for real: the thin adapter holds for sync, dump and query, with
byte-identical dumps for both dialects. What it also lands is a capability the product cannot route
— `spanner/driver` lists and dumps PostgreSQL-dialect databases, while `SPANNER`'s
`SplitMultiSQL`/`GetQuerySpan` is still GoogleSQL-only — so the executed statements could be split
differently from the analysed ones, and masking and access checks are keyed off the spans. Keep the
old filter until the parser routes by database dialect. BigQuery's driver is unverified.

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
is adapter code, not product code. **Measured:** the branch is the whole cost, and it survives —
proved against the pinned cluster and by mutation, since replacing the probe with one session's
`Query` makes the dispatch and refusal tests fail while the single-node cluster still passes. The
test that matters is therefore the refusal on a standalone server, not a row count.

### Hive cannot take a caller's dialer; Elasticsearch takes one or a transport

`hive/driver` cannot accept `Config.Dial`: gohive v2 keeps its dialer in unexported fields and has
no field for one, so a caller forwards a local port and names it as `Config.Host` and `Config.Port`.
`elasticsearch/driver` does take `Config.Dial`; `Options.Transport` is the separate seam for a
managed OpenSearch domain's AWS SigV4 signer, and the driver refuses both together. Both engines
can cut over; only Hive behind an SSH tunnel needs the local-port workaround. **Measured for
elasticsearch:** both halves hold by test — an unresolvable `elasticsearch.invalid` host reaches a
stub endpoint through an in-process bastion, and the signed transport carries `DialContext` with the
tunnel asked for `opensearch.invalid:443`. Hive stays unmeasured.

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
layer: no engine measured so far needed it — mssql, mysql and oracle all landed with the assertions
intact — so doing it early spends the nine imports before anything needs the pin gone. And once it
is done, never tidy those blank imports: they are load-bearing, and omni's own imports cannot
replace them, because the assertion table is over Bytebase's parser registry.

### Four smaller corrections

- **The session blocker is not a wave gate.** trino changed zero callers and pg's `*sql.Conn`
  parameter is dead but tolerated, because `QueryConn` keys its own omni session per driver. What
  remains is the `db.SQLDB` shim at 7 call sites, and it disappears one engine at a time.
- **gh-ost loses a probe.** `component/ghost/validator.go` (+13/−5) now fails closed: the contract
  has no binlog-status or grant probe, so a shipped MySQL feature turns off instead of erroring at
  cutover. It is the only silent regression any proof found, and a gap the next contract revision
  should name.
- **Compile is not evidence.** pg's 2,043 and mysql's 2,163 lines of test that no longer compile are
  the largest single line item outside the adapters, and no `go build` sees them. §6 is the same
  lesson from the other side: what compiles is not what runs.

## 6. What a live server finds, and what a compile finds

Every proof produced findings in two categories, and the split explains more than the line counts
do. A compile and a review of the adapter body find the first; only a server finds the second, and
every expensive surprise in the series is in the second.

| | What it finds | How a proof finds it | Cost of missing it |
| --- | --- | --- | --- |
| **A compile finds** | the adapter's shape: its size, its `Config` mapping, which seams it needs, whether any product file names the concrete type, whether the package's own tests still build | `go build ./backend/...`, the adapter's unit tests, grep for `*<engine>.Driver` | a wrong estimate, and a seam that was budgeted and never needed |
| **A live server finds** | what the driver does with real statements and a real handshake: statement text and ordering, per-item verdicts, masking, limits and truncation, TLS, dialect, error text | the old driver and the adapter in one tree, one pinned server, the same script through both, outputs diffed | a cutover that ships broken behind a green build |

### What a live server found

**mssql found five**, all fixed in #1385, and they are the sharpest example because its compile-side
result was the cheapest of the nine: 8 caller files, all of them the shared `GetDB` retirement, no
advisor seam, no plancheck change.

| # | Defect | What it looked like | Why no compile sees it |
| --- | --- | --- | --- |
| 1 | TLS `ServerName` | `TLS Handshake failed: x509: certificate is valid for localhost, not placeholder` — the config is parsed from a `placeholder` DSN host and the client's TLS config verified that name, so every instance with *Verify server certificate* on fails | the certificate is presented at handshake; neither the type system nor the adapter knows the host it should be checked against |
| 2 | `SET`/`DECLARE` statement text | `results[0].GetStatement()` empty: a statement with no result set lost its text | a well-formed result either way; only the product's statement-to-result matching and its logs read the field |
| 3 | Procedure extras after the batch | `EXEC …;GO;DECLARE …` returned the procedure's later result sets after their batch instead of after the script, so a batch stopped matching its statements | where a batch ends is a runtime fact about the script; the adapter hands back a flat list either way |
| 4 | `EXEC` limit and a cut set | `Limit` unenforced for `EXEC` (the product caps at read time), and a `MaxResultBytes` cut set left unread, swallowing the rest of the batch: `cannot match the batch's result sets to its statements` | a cap only fires on data large enough to cut, and the failure surfaces as a statement mismatch far from the cause |
| 5 | Masking-withheld rows | withheld rows set `Result.Err`, which the script loop read as a failed statement, so every batch after a masked procedure silently never ran | a masked batch is an ordinary result to the adapter; only the script loop's error handling distinguishes withheld from failed |

Number 1 is the one to order first, and not because it is first. It is a live defect in
`mssql/driver` today, independent of any cutover: it disables a security posture the operator
explicitly chose and reports itself as a certificate problem, so the natural next step is to debug
the server. The other four degrade results that a test catches. That is why the proof's instruction
is to land the `ServerName` fix before the adapter, not after it.

**mssql's `Dump` moved too, in the same category and for a different reason.** Rendering through
omni's catalog keeps a filtered index's `WHERE ([amount] IS NOT NULL)`, which Bytebase's renderer
dropped, and loses the doubled `;` after a view definition — 921 → 948 bytes on the same synced
schema. Nobody wrote a rule that changed; the renderer did, so the only way to see it is to dump the
same schema through both. No test pins it yet, which is worth saying plainly: a behaviour change
that a user sees in a stored `raw_dump` and that the diff suite does not compare is the gap the next
mssql pass should close first.

**spanner found one, of a second flavour: a library refusing a combination.** Pointing an instance
that carries a stored Cloud credential at the emulator made the driver build a credential option and
`WithoutAuthentication` together:

```text
failed to create the Spanner database admin client: options.WithoutAuthentication is incompatible
with any option that provides credentials
```

The adapter drops the credential when `emulator` is named, and a test holds it. Neither an omni gap
nor a product change — the third class the proofs kept producing, an interaction that only a test
endpoint creates. Real-project authentication stays unverified, because an emulator cannot answer an
IAM question.

**elasticsearch found the category in its most expensive form: not a bug, a contract.** The proof's
first assertion expected the bulk verdict's error to contain `document_missing_exception`; the
driver's real text is `1 of 2 bulk items failed: [missing]: document missing`. The string was the
small part. The behaviour behind it is that a partly applied bulk used to answer 200 with
`error=""`, and now fails the task run — a change no reviewer can see in either implementation,
because both are internally consistent and only the pair disagree. Its other live finding is the
one still open: AWS ping and sync errors lose their hint (`authentication failed: unable to assume
role …` becomes `ping failed: [403 Forbidden] <body>`), and no compile says whether support relies
on the old text. That is elasticsearch's decision #8 in §3's table.

The same category, smaller, from the earlier proofs: cassandra's dropped keyspace answering
`driver.ErrNotFound` where the package doc promised an empty schema; redis's context deadline only
beginning to interrupt a command once it is inside the client (5.0 s against a 200 ms deadline, now
201 ms), and its `Execute` stopping at a refused command where the old pipeline wrote the rest of
the script. Each is a behaviour, not a shape.

### What a compile found, and how it kept over-predicting

The compile-side findings are real — the seams, the type pins, the test rewrites — but the
predictions ran wrong in one consistent direction: the adapter is bigger than the prototype, and the
seam list is both smaller and elsewhere. Every measured adapter came in above its estimate, by
between 1.8× (redis, 154 → 283) and 5.2× (mysql, 136 → 706), with seven of the nine between 2×
and 3×.

Seams were wrong in both directions, and the net is what to remember:

- **Budgeted, not needed.** mssql's advisor probe and plancheck change; spanner's `Prober`, `SQLDB`,
  `Activity` and private-link guard, each inapplicable for a reason (no `SPANNER` rule exists,
  `GetDB` was already nil, the activity read is `pg_stat_activity`, the private-link guard is
  Oracle's); redis's `SQLDB` and `Prober`; oracle's `CountAffectedRows` and `Version`.
- **Not budgeted, needed.** pg's advisor probe, which was the one forced interface; mysql's
  `DriverHandle`; oracle's private-link guard signature.
- **Budgeted, and it holds.** The type pins. plancheck's `*mssql.Driver`, `*mysqldriver.Driver` and
  `*oracledriver.Driver` assertions all survive, because an adapter that keeps its concrete type
  keeps `CountAffectedRows` reachable.

Three rules follow, and they are the series' output rather than any one engine's:

1. Budget the adapter at the measured band, not the prototype's, and expect the seams to be fewer
   than predicted and in different places.
2. Treat "it compiles and no caller changed" as the proof's *start*. Five wave-1 engines reached
   that point and none of them was finished there: each still owed a live pass, and elasticsearch
   owed ten decisions.
3. Run the old driver beside the new one on one pinned server. That is what turns behaviour into an
   accept/reject list — eleven items for redis, thirteen for spanner, ten for elasticsearch —
   instead of a review opinion, and it is the only pass that surfaced mssql's five.

No proof could have been cleared by reading its diff, and no proof could have been cleared by
compiling it either. The two passes answer different questions: the compile answers "did every
caller survive", the server answers "does every statement still mean what it meant".

## 7. Order

Waves are by measured cost and by forced change, not by engine importance. The numbers in the
parentheses are the §2 rows, kept so a reader can see which prediction moved and which did not.

**Wave 0 — the shared caller layer (no engine, and not a gate).** Measured across the proofs and
identical wherever it appears: `db/util/omnidriver.go` (120 lines, ~102 code), which trino, mysql
and oracle each paid separately and which redis, spanner and elasticsearch each wrote again inside
their adapter, so budget it once; `db.Driver` losing `GetDB` for `SQLDBProvider` + `db.SQLDB`
(+25/−1); `advisor.Prober` + `Context.Prober` (+46 and +11/−1); and the MySQL-family config builder
the prototype measured at 158 lines (82 of them the RDS CA bundle). About +202/−2 plus the builder.
None of it exists on Bytebase `345007290e`, and the proofs say that is survivable: every wave-1
engine measured so far — cassandra, trino, redis, spanner, elasticsearch — needed **no** seam and
changed **zero** callers, and four of the five changed no product file either (trino's were the pin
bump's dump fixtures, not the adapter). So wave 0 gates neither wave 1's first adapter nor its
fifth. It
gates wave 2's `GetDB` retirement (mssql's 8 files are all that seam) and wave 3, where pg's
`Prober` and oracle's reuse of it verbatim are the whole point. Land it before the fifth wave-1
adapter if you want the conversion lines once instead of six times; do not let it gate the first.

**Wave 1 — adapters with no forced product change (1-12).** `cassandra` first: 312 adapter lines, 0
callers, 0 product files, 928 lines deleted, a live pass on `cassandra:5.0.9` — a driver swap plus a
deletion, and nothing else for a reviewer to read. `trino` second: 234 lines, 0 callers, and the
first engine to exercise the fixture path (4 dump strings, +14 lines for views, 2 result-shape
expectations) that every later engine repeats. `redis` third: 283 lines, measured, with no seam and
no decision left open; its cost is the pin, not the adapter. `spanner` fourth: 333 lines, whose one
open decision (the dialect listing) has an answer that needs no product change. `elasticsearch`
fifth: 347 lines, ten behaviour changes, and the one decision a support case could notice. Then
`cosmosdb`, `dynamodb` and `mongodb` in any order, each with the §3 fixture pass. `clickhouse`,
`doris` and `starrocks` follow once wave 0's shared builder exists: they appear in no plancheck,
advisor or ghost switch, and trino showed the session seam costs no caller. `hive` keeps its
documented limit (forward a local port); `elasticsearch` no longer has one, because both of its
seams were proven by test.

Cassandra-first survives, on one clause. The plan's reason was size — 126 lines, the smallest
prototype — and the size was wrong: 312, a 2.5× miss, and redis now measures smaller at 283. What
holds is the clause the plan put second: cassandra is the only wave-1 engine that needs no
engine-specific `Open` options, no `Params` and no `Options`. The other half of that sentence — zero
product files — no longer distinguishes it: redis, spanner and elasticsearch measured zero product
files too, and four of the five measured rows now do, all with zero callers. `cosmosdb` (510
non-test lines in Bytebase) is smaller to delete but needs `Params` auth/tenantId/emulator, the
emulator TLS path and `QueryOptions.Schema` as the container; `redis` needs topology `Options` and
the `AllNodesQuerier` probe; `spanner` needs `Params` and the emulator flag; `elasticsearch` needs
`Options{OpenSearch, Transport}`. So the order does not change and the reason narrows: cassandra is
first because it needs nothing outside the package, not because it is small and not because it
alone has no product side.

The width of the band is the finding. Five measured wave-1 adapters span 234 to 347 lines — 113
apart, against prototypes that spanned 101 to 169 — so adapter size does not order this wave. What
does is how many behaviour changes each engine hands back for a decision: cassandra four, redis
eleven (all accept), spanner thirteen plus one product decision, elasticsearch ten plus one
support-facing decision. That is why spanner precedes elasticsearch despite the smaller adapter:
spanner's decision has an answer that needs no product change (keep the filter), elasticsearch's may
need the adapter to preserve a hint it currently drops. Both can be decided before their adapter is
written, and then they would move earlier.

**Wave 2 — the `GetDB`-retirement engines (13, 14).** `mssql` first: measured at 8 caller files
+50/−16 (net +34), no advisor seam and no plancheck change — the cheapest product side any proof has
produced, and it needs only wave 0's `db.SQLDB`. #1385's five driver fixes are already on main
(`81b91221d`), and its TLS `ServerName` one is worth treating as its own defect rather than a
cutover's leftover, since it breaks verified TLS for every instance today (§6). Then `cockroachdb`,
unmeasured: it is the only remaining engine likely to need the plancheck interface, and if it does,
that change carries the nine blank parser imports (§5). `cockroachdb` stays its own adapter rather
than aliasing pg (§4).

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

## 8. Start the first cutover

1. Land wave 0 as its own change, but not as a gate on the first adapter: the shared conversion,
   the `db.SQLDB` seam and the advisor `Prober`. plancheck's interface and its nine imports belong
   to wave 2, when cockroachdb is the engine that needs the pin gone — mssql, mysql and oracle have
   now landed with the assertions intact (§2). Cassandra's proof needed none of it; every adapter
   after it wants the 120 conversion lines and, from wave 3 on, the `Prober`.
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
   shapes to change with the pin bump, and decide each one before the cutover lands. Where an
   engine renders no schema this pass is cheap and still worth running: spanner's dump came through
   byte-identical for both dialects, and redis's and elasticsearch's no-op dumps round-trip empty
   through the sync runner before and after.
7. Then take `TRINO` (row 9) with the same shape, which is where the fixture path is exercised for
   the first time, then the rest of wave 1 in §7's measured order, one package at a time. The other
   engines pay the same fixture pass; none of them should be the one that first discovers it. The
   live pass, not the fixture pass, is what finds the expensive defects — mssql's five and
   elasticsearch's bulk verdict all arrived that way (§6) — so run the old driver beside the adapter
   on the same server rather than only the new one.

## 9. What stays product-side

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
