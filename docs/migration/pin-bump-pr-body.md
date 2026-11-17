# Pin bump: omni `8851de182729` → `ad3d0f433`

2026-11-17 · PR body for Bytebase · written for Bytebase `345007290e` and omni `ad3d0f433`

This is the body of the Bytebase pull request that moves `github.com/bytebase/omni` from the current
pin (`8851de182729`, 2026-11-15) to omni `ad3d0f433` (2026-11-17): 143 commits, 1,314 files,
+3,744,722/−23,225. It says what the range brings, what this PR (or one before it) must carry, what a
user will see differently, what is deliberately outside the range, and how the range was audited.

Everything below is measured at the named SHAs. Where a number comes from a cutover proof, the proof's
own record is cited instead of being re-derived: [dispatch plan §3](../specs/2026-11-17-cutover-dispatch.md#3-measured-cost)
(adapter sizes, deleted lines, product-side diffs), [open decisions](../specs/2026-11-17-cutover-open-decisions.md)
(which choice belongs to whom), [release note](cutover-changes.md) (per-engine user-visible rows).

## 0. The short version

- **19 driver trees enter the dependency.** Only `pg/driver` exists at the pin, and the pin predates
  `redis/driver` entirely, so "cut over redis" is this bump rather than an adapter.
- **The bump compiles.** A scratch Bytebase tree at `345007290e` with only the pin moved passes
  `go build ./...` and `go vet ./backend/...`; nothing in the range is an API break for a package
  Bytebase already imports.
- **It is not a mechanical change.** Five parser expectations and the whole `backend/plugin/schema/pg`
  golden set move, because the range fixes rendering the old goldens pinned. §2.1 lists each one.
- **Five contract surfaces land** and one you may have expected does not: `Config.TokenProvider`,
  `Driver.PrivateDatabaseLinkOwner`, `driver.AllNodesQuerier`, `ExecuteOptions.Isolation` are in;
  `SQLDBProvider`/`db.SQLDB` is Bytebase-side and is **not** in omni (§1.6).
- **The adapters are not in this bump.** omni holds the drivers; Bytebase still runs
  `backend/plugin/db/*`. Landing the pin changes which code is available, not which code executes,
  and the one place where it changes a byte of user output today is the fixtures in §2.1.

## 1. What the bump brings

### 1.1 Drivers: 19 new engine trees

| Engine (`storepb`) | Tree | What the range adds | PRs |
| --- | --- | --- | --- |
| `MYSQL` | `mysql/driver` | sync + dump, then query and change scripts | #1235, #1344 |
| `MARIADB` | `mariadb/driver` | full engine: sync, dump, query, execute, row estimates, live probes | #1237 |
| `TIDB` | `tidb/driver` | as MariaDB | #1237 |
| `OCEANBASE` | `oceanbase/driver` | as MariaDB | #1237 |
| `POSTGRES` | `pg/driver` | query and execute (sync + dump predate the pin, #1228) | #1233 |
| `MSSQL` | `mssql/driver` | sync + run, then `Dump` | #1239, #1310 |
| `ORACLE` | `oracle/driver` | sync, query, execute, then `Dump` | #1241, #1332 |
| `TRINO` | `trino/driver` | sync, dump, query, execute | #1246 |
| `MONGODB` | `mongo/driver` | read, sync, execute, dump | #1257 |
| `CASSANDRA` | `cassandra/driver` | sync, query, execute, then `Dumper` | #1263, #1341 |
| `ELASTICSEARCH` | `elasticsearch/driver` | sync, query, execute, then `Dumper`; caller `Dial` and `Transport` | #1263, #1341, #1365 |
| `COSMOSDB` | `cosmosdb/driver` | sync, query, then `Dumper` | #1266, #1341 |
| `DYNAMODB` | `dynamodb/driver` | sync, query, execute, then `Dumper` | #1272, #1341 |
| `REDIS` | `redis/driver` | sync, query, execute, `Dumper`, `AllNodesQuerier` probe | #1275, #1341, #1365 |
| `HIVE` | `hive/driver` | sync, query, execute, then server-DDL `Dump` | #1279, #1343 |
| `CLICKHOUSE` | `clickhouse/driver` | new tree, then server-DDL `Dump` | #1274, #1350 |
| `SPANNER` | `spanner/driver` | sync, query, execute, dump; PostgreSQL-dialect databases | #1288, #1338 |
| `DORIS` | `doris/driver` | sync, query, execute, then server-DDL `Dump` | #1252, #1340 |
| `STARROCKS` | `starrocks/driver` | as Doris, its own `Dump` | #1252, #1351 |
| `COCKROACHDB` | `cockroachdb/driver` | sync, query, execute, dump, decimal fidelity | #1348, #1372 |

Every new tree carries a `dump.go`; `pg`'s is the only one that predates the pin. Dumper coverage
is therefore complete for all 20 trees on main, which is the fact the cutover needed: no engine is
still "one that cannot delegate `Dump`" except the four drafts in §4.

Engine-specific driver fixes in the same range, all with a failing test in omni: mssql TLS
`ServerName`, statement text, batch ordering, `EXEC` limits and masking (#1385); mssql transient
`ReturnMessage` handling (#1313); clickhouse `EXPLAIN ANALYZE` rows (#1394); oceanbase dump fidelity
for generated columns, temporal defaults, triggers, materialized views, events, and recreated
events (#1374, #1396); mysql/mariadb query answers (#1392); cockroachdb decimal precision (#1372);
spanner `GRANT`/`REVOKE` routing (#1289); trino `information_schema` and dump restore (#1292, #1280);
elasticsearch error matching (#1327); cassandra NULL vs zero (#1305); mariadb/oceanbase/tidb fork
decode fixes (#1354, #1362).

### 1.2 Analysis moves

| What moved | Where it now lives | PRs |
| --- | --- | --- |
| Redshift and Trino query spans | `redshift/analysis`, `trino/analysis` | #1256, #1259 |
| Redshift's remaining judgment (diagnose, explain, limit, query, resource change, statement ranges + type, system objects) | `redshift/analysis` | #1329 |
| Cassandra and GoogleSQL query spans; BigQuery/Spanner corpora | `cassandra/analysis`, `googlesql/analysis` | #1261 |
| Doris and StarRocks analysis (split, query type, `EXPLAIN` unwrap, system databases, limit wrap) | `doris/analysis`, `starrocks/analysis` | #1262, #1265 |
| The MySQL family's result-limit wrap chosen in analysis | `{mysql,mariadb,oceanbase,tidb}/analysis` | #1268 |
| The MySQL family's query-span corpora (447 recorded cases) | four `analysis` packages | #1267 |
| GoogleSQL `StatementTypeOf`, Mongo `ExplainStatement`/`DescribeExplain`, CosmosDB `ValidateQuery`, Trino `QueryTypeOf`/`ValidateQuery`/`Diagnose` | their `analysis` packages | #1346 |
| CockroachDB analysis | `cockroachdb/analysis` | #1361 |

Bytebase-side consequence: `parser/{trino,redshift,cassandra,googlesql,doris,starrocks}` are now thin
adapters over these packages, so the pin makes their local copies redundant rather than wrong — which
is what §2.2 deletes.

### 1.3 V1 review ports

The V1 rules enter omni as engine-owned implementations, ported from Bytebase:

| Rule | Engines | PRs |
| --- | --- | --- |
| V1 prior-backup check | mysql, mariadb | #1242 |
| V1 online-migration advice | mysql, mariadb | #1243 |
| V1 prior-backup check | oracle, mssql | #1250 |
| V1 prior-backup check | pg, tidb | #1253 |
| Prior-backup script refusal (pg) | pg | #1254 |
| ALTER judged only against its own database's rows (fix) | mysql, mariadb | #1245 |

These are omni packages. Bytebase's six V1 rule adapters are product work and are **not** in this
bump (§4).

### 1.4 Metadata-proto additions

`proto/omni/metadata/database.proto` gained fields four times in the range. Bytebase does not generate
code from this proto — it consumes the Go package — so the effect is in `backend/plugin/schema/*`:

| Field(s) | PR |
| --- | --- |
| `CheckConstraintMetadata.comment`, `ForeignKeyMetadata.comment` | #1308 |
| A partition's `schema` and `comment` | #1355 |
| User-defined collation and range type | #1358 |
| Oracle's constraint-backed index and sequence fidelity | #1323 |

`backend/plugin/schema/mssql` in particular must learn to write and read the constraint comment
(`@level2type = 'CONSTRAINT'`), or dumps and migration diffs silently drop it. The same class of work
follows any field Bytebase renders itself.

### 1.5 Contract changes

| Surface | Commit / PR | In omni? | What a Bytebase caller must do |
| --- | --- | --- | --- |
| `driver.Config.TokenProvider func(ctx) (string, error)`, called per connection | #1389 | yes (`driver/driver.go`) | mssql's adapter passes a closure over `azidentity`; `User` may be empty and `Password` unused. `Config.TLS`/`Config.Dial` still apply. |
| `driver.PrivateDatabaseLinkOwner` = `OwnsPrivateDatabaseLink(ctx) (bool, error)` | #1241 | yes (`driver/probe.go`) | Oracle's private-link guard loses its `*sql.Conn` (`sql_service.go:697`); leaving it as-is refuses every linked read once `GetDB` is nil. |
| `driver.AllNodesQuerier` = `QueryAllNodes(ctx, script, opts) ([]*Result, error)` on a `Session` | #1365 | yes (`driver/probe.go`) | Redis's adapter answers `QueryOption_ALL_NODES` from the session; keep the refusal on a standalone server (it is the test that matters). |
| `Config.Dial` reachable for Elasticsearch, refused for Hive | #1365 | yes | Elasticsearch adapter passes the tunnel; Hive forwards a local port and names it as `Host`/`Port`. |
| `ExecuteOptions.Isolation` + `IsolationLevel` (`READ UNCOMMITTED` … `SERIALIZABLE`) | #1237 | yes (`driver/execute.go`) | Adopt when the adapter reads `db.ExecuteOptions` (§2.6); `Transaction` already existed. |
| `SQLDBProvider` / `db.SQLDB(driver)` | — | **no** | It is Bytebase-side wave 0 (`db.Driver` losing `GetDB`, +25/−1). Nothing in omni to adopt; §2.3 is where it lands. |

## 2. What Bytebase must change, in this PR or before it

### 2.1 Fixtures the pin invalidates — measured, not predicted

Measured in a scratch worktree of Bytebase `345007290e` with only
`go mod edit -replace github.com/bytebase/omni=<this SHA>`: `go test ./backend/plugin/parser/...` and
`go test ./backend/plugin/schema/pg/`. Five parser expectations fail, in three files; each one is a
fixture that pinned behaviour the engine never had, so the fixture moves and the new text is the one
to accept.

| File:line | Statement | Now | Class |
| --- | --- | --- | --- |
| `backend/plugin/parser/mariadb/query_span_extractor_test.go:688` (assert at `:696`) | `SELECT product_info->'$.id' AS id FROM products` | `syntax error at or near ">"`: MariaDB has no `->` operator | fix; use `JSON_EXTRACT`, or mark the case rejected |
| `backend/plugin/parser/mariadb/query_span_predicate_test.go:72` (assert at `:218`) | `SELECT id FROM j WHERE doc->>'$.k' = 'v'` | `syntax error at or near ">>"` | fix; the case is `rejected` for MariaDB |
| `backend/plugin/parser/doris/limit_test.go:147` (assert at `:197`) | `SELECT * FROM t # hash comment` | `SELECT * FROM (SELECT * FROM t # hash comment) result LIMIT 10;` — Doris does not take `#` as a comment | fix; adopt the wrap |
| `backend/plugin/parser/starrocks/limit_test.go:165` (assert at `:274`) | `SELECT * FROM t # hash comment` | same subquery wrap; StarRocks drops `#` rather than commenting | fix; adopt the wrap |
| `backend/plugin/parser/starrocks/limit_test.go` `WITH RECURSIVE …` case | `WITH RECURSIVE cte_count (n) AS ( … ) SELECT n FROM cte_count;` | `SELECT * FROM (WITH RECURSIVE …) result LIMIT 10;` — the wrap decision is now omni's and it does not splice `LIMIT` here | fix, but read the analysis record first: the in-place expectation was never true of the engine |

The records called this set "the five fixtures", listed as mariadb `:688`/`:72` plus `limit_test.go`
`:147`/`:165` in doris **and** starrocks. Two of those six sites now pass (the doris `INTERSECT` and
starrocks `INTO OUTFILE` cases are untouched by the pin), and the starrocks `WITH RECURSIVE` case —
not on the list — now fails. Treat the table above as the list; re-run the command at bump time rather
than trusting a count.

Three more fixture classes move, from the release note's action list ([cutover-changes.md](cutover-changes.md#product-side-actions-at-the-pin-bump)):

- `backend/plugin/schema/trino/get_database_definition_test.go`: four dump strings gain
  `CREATE SCHEMA IF NOT EXISTS`, and the diff SDL gains views (omni #1280).
- `backend/plugin/aireview/catalog_test.go`: a MySQL secondary-index comment is inline in
  `CREATE TABLE` now, not a separate `COMMENT ON INDEX` (omni `58f05cf7a`).
- Result-shape expectations: trino `DECIMAL` as text and real timestamp precision, tidb/mysql
  `uint64` and `[]byte`, clickhouse/doris/trino no synthetic `Affected Rows 0`.

**The pg schema goldens are a blocker until they are adjudicated.** At the same measurement,
`backend/plugin/schema/pg` fails as a block: 14 `TestGetDatabaseDefinition` subtests plus
`TestGetMultiFileDatabaseDefinition_LegacyIndexExpressionsParse`. The dump text gained
`COMMENT ON SCHEMA "public" IS …`, `CACHE 1` on sequences and `ALTER TABLE … OWNER TO`, because omni's
pg catalog now renders a dump that restores to the metadata it was synced from (omni #1232) — accept
those goldens, they are a fidelity fix. The second failure is not a fixture: the SDL renders
`) PARTITION BY ;` when the snapshot carries a partition with no expression, so omni emits DDL omni's
own parser rejects. That needs the omni side (open PR **#1398**, "render no partition clause a
snapshot does not report"); do not pin the broken text.

### 2.2 The duplicated corpora to delete

Bytebase still carries copies of corpora omni now owns, and every cutover the sweep wrote adds drift
to them. Delete, in the same PR or a preceding one:

`test-data/query-span/*`, `test_split.yaml`, `test_statement_type.yaml`, `test_backup/restore.yaml`,
`statement-ranges/standard.yaml`, `query_span_parity_test.go`, `query_span_probe_test.go`.

Verify the list against the tree at bump time — some of these paths exist under several engine
packages, and the copies that must stay are the ones asserting Bytebase's conversion layer (a
`base.QuerySpan` shape), not the engine's analysis result.

### 2.3 `GetDB` retirement and its call sites

`db.Driver` loses `GetDB() *sql.DB` for an optional `SQLDBProvider` plus the `db.SQLDB(driver)`
helper: +25/−1 in `backend/plugin/db/driver.go`, measured by the pg and mysql proofs. `GetDB` is
already nil for nine of twenty-three implementations on main, so the change is mostly honest typing.
Seven call sites, each with a fallback it can already name:

| Site | What it does with the handle | Fallback |
| --- | --- | --- |
| `backend/api/v1/sql_service.go:117`, `:404`, `:1347` | opens the editor's `*sql.Conn` | already guards `sqlDB != nil`; a driver session instead |
| `backend/runner/taskrun/database_migrate_executor.go:369` | drops gh-ost temp tables | plain `driver.Execute` (no session needed) |
| `backend/runner/plancheck/ghost_sync_executor.go:155` | drops gh-ost temp tables | plain `driver.Execute` |
| `backend/component/ghost/validator.go:54-136` | five binlog probes | none in the contract: fails closed today; §2.7 entry 3 decides it |
| `backend/api/v1/rollout_service.go:675` | task-run session listing | `driver.Activity()` (the pg proof's +63/−3 rewrite) |

Per-engine product-side file counts from the proofs, so the reviewer can size the first adapter:
mssql 8 files (+50/−16), pg 8 for `GetDB` plus 7 for the advisor (+199/−45 across 14 product and 4
test files), oracle 12 files (+243/−36), mysql 17 files (+363/−104 plus 9 blank parser imports),
cassandra/redis/spanner/elasticsearch/doris/trino 0 files.

### 2.4 The advisor prober — two seams, not one

- **`advisor.Prober` and `Context.Prober`** (+46, +11/−1): the generic session probe. pg's five live
  rule files change by +22/−16, and oracle reuses it verbatim.
- **`advisor.DriverHandle`** (+11/−1) for mysql and mariadb: their rules call the driver's AST-aware
  `ExplainJSON`/`EstimateAffectedRows`, which no session probe reaches. Four mysql/mariadb rule files
  grow +73/−20, and the MariaDB rules pay 29 of those lines twice because both engine values live in
  `db/mysql`.
- mssql needs neither (its advisors are static AST rules). TiDB's `STATEMENT_DML_DRY_RUN` needs the
  prober **and** an autocommit probe before the seam is complete.

### 2.5 `*sql.Conn` handling

Keep the parameter as a tolerated nil for this wave; do not drop it engine-by-engine. Thirteen
implementations name it and eleven use it only as somewhere to run one statement, which omni's
session replaces. Two need the connection itself: `mysql` (also `MARIADB`), whose
`getConnectionID` + `StopConnectionByID` is the editor's cancel path, and `mssql`, whose `queryBatch`
runs a script's statements on one connection. The three `sql_service` sites already tolerate nil. The
removal is 65–100 lines across thirteen signatures plus three call sites, and it belongs in one
dedicated change after the mysql family has cut over — dropping it now removes the cancel path with
no replacement (open decision #1).

### 2.6 `ExecuteOptions` adoption

`Transaction` and `Isolation` exist in omni (`driver/execute.go`, `Isolation` since #1237), but every
adapter still parses `-- txn-mode` / `-- txn-isolation` out of the script itself, so no task-run
caller is forced to change. Keep that fallback and adopt the options when the adapter starts reading
`db.ExecuteOptions`; the doris proof is the reference, where `-- txn-isolation` now reaches the
driver and an engine-refused level fails at `Execute`.

### 2.7 The decisions the bump inherits

Each of these is a contract decision already made in omni's records, not a judgement for this PR. The
PR must state which way it went, or the behaviour changes with no record ([open decisions](../specs/2026-11-17-cutover-open-decisions.md)):

| # | Decision | For this bump |
| --- | --- | --- |
| 1 | `QueryConn`'s `*sql.Conn` | keep as tolerated nil (§2.5) |
| 2 | the optional `db.SQLDB` handle | adopt only where it retires `GetDB` |
| 3 | the truncation cut point (pg: 16,384 → 9,448 rows at 1 MB, same message) | accept; the contract measures the cap |
| 4 | a NUL byte in a username now fails at first use, not at `Open` | accept |
| 6 | cassandra's missing keyspace is now `ErrNotFound` | accept; name it in the release note |
| 9 | oracle's private-link guard takes no connection | land before the oracle adapter |
| 11 | `MaximumSQLResultSize` of 0 means "no cap" in omni, "zero-byte cap" in Bytebase | refuse 0 in the adapter |
| 12 | a procedure created qualified syncs qualified (doris) | leave until a consumer needs bare names |
| 13 | dynamodb keeps verifying TLS whatever `verify_tls_certificate` says | verify until the product can express an explicit opt-out |

Plus the product-side items the release note owns: gh-ost's five binlog probes (fails closed today —
either add the probe to `mysql/driver` or record that validation is off), plancheck's eight concrete
assertions (optional for every measured engine; if taken, the nine blank parser imports are
load-bearing), Spanner's PostgreSQL-dialect listing (do not ship it alone), and the three engines with
no server-backed test at Bytebase's seam (cassandra, hive, elasticsearch).

### 2.8 What order

The wave order the measurements produced, if the adapters follow in later PRs: wave 0 the shared
caller layer (`db/util/omnidriver.go`, `db.SQLDB`/`SQLDBProvider`, `advisor.Prober` — not a gate,
since cassandra needed none of it); wave 1 cassandra, trino, redis, spanner, elasticsearch, then
cosmosdb/dynamodb/mongodb, then clickhouse/doris/starrocks and hive; wave 2 mssql, cockroachdb
(`GetDB` only); wave 3 pg, oracle, tidb/oceanbase, mysql/mariadb; wave 4 the drafts. This PR is only
the pin, §2.1's fixtures, §2.2's deletion, and the decisions it must record.

## 3. What users will see differently

The full per-engine list is the release note: [docs/migration/cutover-changes.md](cutover-changes.md)
(omni #1397). It is built from fourteen cutover proofs against pinned servers — trino, pg, mysql,
cassandra, oracle, mssql, redis, spanner, elasticsearch, doris, tidb, clickhouse, dynamodb and hive —
and every row says whether a test pins it or only a proof measured it.

Three classes, and the last one is where the risk sits:

1. **Fixes.** Dumps that could not be restored now restore (doris, clickhouse, mssql's filtered index,
   oracle's sequence definitions, tidb's partition order, hive's materialized views); verified TLS
   connects (mssql); typed values replace strings and NULLs (mysql/tidb unsigned and binary,
   clickhouse timestamps, hive floats); a partly applied bulk fails the task run (elasticsearch);
   timeouts interrupt the command (redis); `GRANT`/`REVOKE` route as DDL (spanner).
2. **Same looking, different behaviour.** Result caps cut at a real byte boundary (pg 16,384 → 9,448
   rows at 1 MB; spanner 32 → 25 at 300 bytes; doris 4,096 → 3,973 at 1 MB); one session per call, so
   a `SET SESSION` does not survive a query (trino, pg); a whole script is parsed before anything is
   sent, so a parse error fails the call (mysql, mariadb); scripts pay one round trip per command
   (redis: 200 `PING`s 2.4 ms → 13.8–23.4 ms); task-run history gains per-statement entries on every
   engine that had none.
3. **Needs a product decision.** A reviewer has to accept each of these consciously; none is a bug in
   the driver:

| Engine | What changes | Owner |
| --- | --- | --- |
| mysql/mariadb | gh-ost binlog validation fails closed: `binlog access validation needs a driver with a SQL handle` | add the probe to `mysql/driver`, or record that the feature is off |
| mysql/mariadb | a masked `CALL p()` returns the count shape (`Affected Rows: 0`) with no rows | confirm the grid and the support copy |
| mysql, mariadb, doris, tidb, clickhouse | `MaximumSQLResultSize == 0` becomes "every row" instead of "no rows" | refuse 0 in the adapter |
| mariadb | a masked untraced result set is withheld with the count shape | as mysql |
| tidb | invalid SQL now reads `expected identifier (line 1, column 51)` instead of `Syntax error at line 1:51: … related text: …` | wrap with `convertParseError` if the console copy depends on the old text |
| tidb | `STATEMENT_DML_DRY_RUN` returns zero advice when there is no SQL handle | land the prober + autocommit probe |
| tidb | `-- ghost` cleanup would panic unless ghost is gated to the MySQL family | engine gate |
| clickhouse | a script that cannot roll back still logs `BEGIN`/`ROLLBACK` | fix at the log layer, not the adapter |
| cassandra | a deleted keyspace is a failed sync (`not found: keyspace X does not exist`) instead of an empty success | accept, name in the release note |
| cassandra | the dead read-only branch stays dead; `requiresAdminDataSource` is the real gate | product gate, not the adapter |
| dynamodb | TLS keeps verifying even with *Verify server certificate* off, as today | product tri-state before honouring the flag |
| dynamodb | an assumed role's keys are frozen per `Open`, not refreshed for the driver's life | engine-neutral credential source if a long run needs it |
| elasticsearch | AWS ping/sync errors lose the role-ARN hint (`ping failed: [403 Forbidden] <body>`) | preserve the hint in the adapter if support relies on it |
| spanner | PostgreSQL-dialect databases are listed and dumped, but split/query-span is still GoogleSQL-only | keep the old filter until the product routes by dialect |
| spanner | `CREATE DATABASE "name"` creates a PostgreSQL-dialect database | unreachable today (`SPANNER` is absent from `EngineSupportCreateDatabase`) |
| hive | an SSH data source is refused with a port-forward instruction instead of silently ignoring the jump host | local forwarder in `backend/plugin/db/util` restores tunnels |
| oracle | the catalogue renders a synced key's columns quoted: `PRIMARY KEY (""ID"")` | decide the normalisation, no test crosses the two catalogues |
| doris, mariadb | a procedure created as `probe.p1` syncs qualified (mysql records the bare name) | re-recorded golden if the product wants bare names |
| trino | DDL through the query path returns an empty result instead of `Affected Rows: 0`; `SHOW`/`DESCRIBE` are no longer row-capped | 5 adapter lines restore the row; the real cap belongs in `sql_service` |

## 4. What is NOT in this bump

So that nobody infers it from a neighbouring engine:

- **The adapters.** No `backend/plugin/db/*` package changes into an omni driver here; the pin only
  makes the drivers available. Everything in §2.3–§2.6 is the prerequisite for that later PR.
- **The wave-0 layer.** `db/util/omnidriver.go`, `db.SQLDB`/`SQLDBProvider`, `advisor.Prober` and the
  shared MySQL config builder are Bytebase code. Only `ExecuteOptions.Isolation`, `Config.Dial`, the
  two probe interfaces and `Config.TokenProvider` are omni's (§1.5).
- **The four unverified cloud-driver drafts.** `redshift/driver` (#1371), `bigquery/driver` (#1377),
  `snowflake/driver` (#1378) and `databricks/driver` (#1379) are open drafts with no tapes, no sync
  goldens and no pin rows. None is on main, so `REDSHIFT`, `BIGQUERY`, `SNOWFLAKE` and `DATABRICKS`
  gain nothing from this bump.
- **The endpoint-blocked engines.** Redshift and Snowflake cannot be verified without a cluster and an
  account. Their drivers are drafted and blocked on that, not on this pin.
- **The parked completion branches.** `completion-mysql` (26 commits ahead of main),
  `completion-spike` (19), `completion-mariadb` (11), `completion-oceanbase` (10),
  `completion-tidb` (10) are not merged, and the Mongo completion cut
  ([design](../specs/2026-11-17-mongo-completion-cut-design.md)) is design only. Bytebase's
  `backend/plugin/parser/mongodb/completion.go` keeps working exactly as it does today.
- **Open omni work at writing.** #1398 (pg partition clause), #1386 (hive managed-table restore),
  #1376 (foreign servers and external table locations — a metadata-proto change), #1264 (GoogleSQL
  predicate columns, a parked draft). If any of them merges first, the body of this PR should name the
  new SHA and re-run §2.1.
- **Everything the proofs left unpinned.** The release note marks each row `Proof only` where no test
  in the repository holds the behaviour: most of doris, clickhouse, hive, tidb, elasticsearch,
  dynamodb, redis, spanner and cassandra's user-visible deltas, including the pg truncation row count.
  A green omni suite does not mean these are held; it means they are measured.
- **The Bytebase test-container pins** that move with the cutovers (redis `7.4.12-alpine`, tidb
  `v8.5.5`, mssql `2022-CU27-GDR1-ubuntu-22.04`) — they belong to the adapter PRs, not the pin.

## 5. How this bump was audited

**Range.** `8851de182729` (2026-11-15 01:04 −0800, `feat(pg): move Bytebase's SDL style … #1226`) →
`ad3d0f433` (2026-11-17 05:13 −0800, `fix(mysql,mariadb): answer the queries the MySQL driver already
answers correctly #1392`), 143 commits. Every commit is a reviewable PR; nothing was pushed directly.

**Gates at `ad3d0f433`**, in a worktree of omni's `origin/main`, Docker pointed at a nonexistent socket
so a container attempt fails rather than starts, `OMNI_RECORD` unset, CI's sweep values
(`SPLITTEST_N=1000000`, `WALK_SWEEP_N=500`):

| Gate | Result |
| --- | --- |
| `go build ./...` | pass |
| `go vet ./...` | pass |
| `go test ./... -count=1` (replay-only, everything) | **pass**, 182 packages ok, 29 with no test files, 0 fail, 79 s |
| `go generate` over the 10 AST walkers, `git diff --exit-code` | clean, no generated change |
| `go test ./enginetest/ -run '^TestDocumentedRecordingCommands$'` (the `VERIFY.md` lint) | pass |
| `go run ./enginetest/cmd/compatcheck origin/main` | pass — commits change engineered answers only by re-recording |
| `scripts/affected-packages.sh origin/main` | empty, as expected at main |
| proto `format`/`lint`/`generate`/`breaking` | last measured at `46dc54bfd`; `git log 46dc54bfd..ad3d0f433 -- proto/` is empty, so the measured tree and this one are byte-identical under `proto/` |

An earlier full protocol at `bf1a3396e` — the same gates plus `go vet`, two consecutive full runs
(211 packages, 178,241 tests, 18 skips, 0 fail, identical outcomes), the nested `googlesql/spanner-oracle`
module, and a flake report that found no flakes — is the reason `go vet` is on this list at all: it was
red with four `unreachable code` findings in `oracle/parser` until #1369. The 17 commits between that
measurement and `ad3d0f433` are 8 docs, 3 test-only, 5 engine fixes and one contract feature (#1389),
each merged on the usual gate (CI green plus a Codex review with no unresolved thread). The table
above is the first measurement that covers the exact SHA this PR names.

**Per-engine measurements.** The cost figures in §2.3–§2.4 are not estimates: ten cutovers were built
end to end in scratch Bytebase trees against live servers (trino, pg, mysql, cassandra, oracle, mssql,
redis, spanner, elasticsearch, doris), and four more proofs (tidb, clickhouse, dynamodb, hive) report
per-behaviour deltas. Their per-file diffs and verbatim failures are in the dispatch plan's §3 and §6
and in the release note's tables; this body restates only the totals it needs.

**The fixture measurement.** A scratch Bytebase worktree at `345007290e` with only the pin moved
(`go mod edit -replace`, no other change) ran `go test ./backend/plugin/parser/...` and
`go test ./backend/plugin/schema/pg/`; §2.1 is that output, including the fact that the record's
six-site list had two sites that no longer differ and one failing site it did not list. The parser
fixtures' own container (`mariadb:11.8.8`) was the only container started, and it was removed.

**What this audit did not cover.** Bytebase's full suite is not green-or-red evidence for this pin:
`backend/plugin/advisor/...` was run separately and passes, `backend/plugin/db/...` was still running
at writing, and no API or end-to-end Bytebase test was run. `go build ./...` and
`go vet ./backend/...` for that tree are green. The `backend/plugin/schema/pg` result is
adjudicated in §2.1, not fixed here. The four draft drivers, the parked completion branches and the
open omni PRs in §4 were not exercised at all.

## 6. Reviewer's checklist

- [ ] The pin in `go.mod` is `ad3d0f433…` (full SHA), `go mod tidy` run, no `replace` left behind.
- [ ] §2.1's five parser expectations move, and the pg schema goldens move with #1398's fix in place.
- [ ] §2.2's corpora are deleted; the conversion-layer assertions stay.
- [ ] Every decision in §2.7 either lands or is named in the PR body as deferred, with the reason.
- [ ] No new `GetDB` call site, no adapter returning a fake handle.
- [ ] The release note ([cutover-changes.md](cutover-changes.md)) is linked so a support engineer can
      answer "is this a regression?" without reading a driver.
