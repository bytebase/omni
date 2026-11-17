# Cutover open decisions: what the proofs left the pin bump

2026-11-17 · Decisions

Bytebase `345007290e`, omni `af2e1a130`. The cutover plan
([2026-11-17-cutover-dispatch.md](2026-11-17-cutover-dispatch.md)) says how each engine reaches its
driver; the contract ([2026-11-15-driver-contract-design.md](2026-11-15-driver-contract-design.md))
fixes what a driver must do. The cutover proofs then cut engines over for real: trino, pg, mysql,
cassandra and oracle have reported, and mssql was still running when this was written. Each report
ends with items that are not measurements but choices.

**Amended 2026-11-17.** Ten proofs have reported: redis, spanner and elasticsearch in the dispatch
plan's revision (#1390), then mssql, then doris. Only MariaDB and DynamoDB are still inside.

This note is for the person opening the Bytebase pin-bump PR. It says which choices are already
settled, which are theirs, and what a deferral costs. The measured costs themselves — adapter lines,
deleted lines, product-side files — live in the dispatch plan's §3 ("Measured cost") and are not
restated here.

## Which of these are yours

| # | Decision | State | Who acts |
| --- | --- | --- | --- |
| 1 | `QueryConn`'s `*sql.Conn` parameter | settled for the cutover wave; drop later | nobody in the pin bump |
| 2 | `db.SQLDB` optional handle | open, one shape to pick | the first pin bump that removes `GetDB` |
| 3 | Truncation cut point | settled by the contract | nobody; note the row change |
| 4 | NUL byte in the username | settled: accept | nobody; do not "fix" |
| 5 | plancheck's concrete assertions | open, timing only | the first cutover that lets a concrete driver type go (cockroachdb, since mssql did not need it) |
| 6 | Cassandra sync on a missing keyspace | open, one behavior to own | the cassandra pin bump |
| 7 | Cassandra's new task-run log entries | settled: keep | nobody |
| 8 | Cassandra's SSH, no-op `Dump`, dead read-only branch | settled: recorded in §9 | nobody |
| 9 | Oracle's private-link guard | open, must land before the oracle pin | the oracle cutover |
| 10 | Oracle's `GetVersion` | settled: the metadata is the fallback | nobody |
| 11 | `MaximumSQLResultSize` of 0 | open, one mapping to pick | the doris pin bump |
| 12 | A procedure created qualified syncs qualified | open, needs a re-recorded golden to change | the doris pin bump, if the product wants bare names |

A pin bump that changes nothing but the `go.mod` line and an adapter body takes 1, 3, 4 and 7 as
written, adopts 2 only where it removes `GetDB`, owns 6 for cassandra, 9 for oracle and 11 and 12
for doris, and leaves 5 to the PR that first lets a concrete driver type go.

## 1. `QueryConn(ctx, conn *sql.Conn, …)`'s parameter is dead for pg

### What the proofs measured

The pg adapter never reads `conn`: it asks `Querier.Session` for its own session, and `GetDB()`
returns nil. Trino's cutover did the same, and neither changed a caller. The cutovers found callers
already tolerate the handle being nil, because `GetDB()` is nil for nine of the twenty-three engine
implementations on main (cassandra, cosmosdb, dynamodb, elasticsearch, mongodb, redis, spanner,
bigquery, databricks).

On main the parameter is named and used by thirteen implementations (`pg`, `mysql`, `tidb`,
`oceanbase`, `doris`, `starrocks`, `clickhouse`, `mssql`, `oracle`, `snowflake`, `redshift`,
`cockroachdb`, `trino`); the other ten take `_ *sql.Conn`. Of the thirteen, eleven use it only as
the place to run one statement — the thing omni's `Session` replaces. Two need the connection
itself, not just somewhere to run a statement:

| Engine | Why the handle is real |
| --- | --- |
| `mysql` (also serves `MARIADB`) | `getConnectionID` before the query, so `StopConnectionByID` can cancel the query from another goroutine (`mysql/mysql.go:620-691`; `tidb`, `oceanbase`, `doris`, `starrocks` copy it) |
| `mssql` | `queryBatch` runs each statement of a multi-result script on the same connection, so the batch shares the session (`mssql/mssql.go:439-473`) |

Three `sql_service` sites open the handle from `GetDB()` and pass it
(`backend/api/v1/sql_service.go:117`, `:404`, `:1347`, reaching `:1218`). Each already branches on
`sqlDB != nil`.

### Options

- **Drop it in the cutover wave.** Mechanical: thirteen signatures of five to eight lines each,
  plus the three `sql_service` sites, 65–100 lines. But for the mysql family it also removes the
  handle their cancel path reads, unless the adapter opens one from `GetDB` itself — a session
  behavior change the proofs did not measure. And it forces the thirteen into a pg PR.
- **Keep it as a tolerated nil.** Every adapter ignores the argument; `db.Driver` keeps a parameter
  that no cut-over engine reads.
- **Keep it for the engines that use it and document the split.** That is a second interface (or a
  runtime branch), and it keeps two shapes alive for the length of the migration.

### Recommendation

Keep it as a tolerated nil for the cutover wave, and record the split: for a cut-over engine the
argument is dead; for the mysql family and mssql it is still the session. Drop it in one dedicated
change once those two have cut over to `Session` — then the same 65–100 lines and three
`sql_service` sites are purely mechanical, because no engine semantics ride on the handle and no
caller needs to preserve an editor session. Cost now: zero. Cost then: unchanged, in a change that
cannot regress a feature.

### If nobody decides

The pg pin bump may land with the parameter intact — nothing breaks today. The cost is deferred and
compounds: every adapter written from here on takes an argument it must ignore, and the "one
session per `QueryConn`" behavior change stays undocumented in the contract. The failure mode is a
future author who believes the handle works and writes a cancel path or a session-scoped `SET`
against it, which silently does nothing; or someone drops the parameter in an engine PR and takes
MySQL's kill-on-cancel with it.

## 2. Whether the optional `db.SQLDB` handle survives

### What the proofs measured

`db.SQLDB` does not exist on main; it is the pg and mysql proofs' proposal. The pg proof removed
`GetDB` from `db.Driver` and replaced it with an optional `SQLDBProvider` interface plus a
`db.SQLDB(driver)` helper; the mysql proof measured that edit to `backend/plugin/db/driver.go` at
+25/−1, which also adds the transaction options. Seven call sites read it, and on pg each resolves
to nil:

| Site | What it does with the handle | On pg |
| --- | --- | --- |
| `sql_service.go:117`, `:404`, `:1347` | open a `*sql.Conn` for the editor and export | nil, no conn; adapter session instead |
| `rollout_service.go:675` | task-run session listing | nil |
| `component/ghost/validator.go:54-136` | five binlog probes: `SHOW MASTER STATUS`, `SHOW BINARY LOG STATUS`, `@@log_bin`, `SHOW GRANTS`, `@@binlog_format` | unreachable (gh-ost is MySQL) |
| `runner/plancheck/ghost_sync_executor.go:155` | drop gh-ost temp tables | unreachable (MySQL) |
| `runner/taskrun/database_migrate_executor.go:369` | drop gh-ost temp tables | unreachable (MySQL) |

`GetDB` is already optional in practice: nine of twenty-three implementations return nil on main and
every guarded caller works. The mysql proof found the shim's limit: gh-ost validation now fails
closed — `binlog access validation needs a driver with a SQL handle` — so the cutover silently
turns off a shipped MySQL feature. The two cleanups are plain `DROP` statements and can go through
`driver.Executor`, which the contract already has; the five probes cannot, because the contract has
no binlog-status method.

### Options

- **Optional capability (`SQLDBProvider`) and retire it site by site.** `GetDB` leaves
  `db.Driver`, so a driver that has no handle says so at the type level instead of lying with nil.
  Each site names its fallback, which they all already have.
- **Retire it engine by engine, keeping `GetDB` on `db.Driver`.** Zero interface change; the
  method survives with nine nil implementations and grows one per cutover until it is dead.
- **Keep the handle everywhere the engine has a pool.** Rejects the contract's `GetDB` removal and
  keeps the SQL editor's `*sql.Conn` in the product forever.

### Recommendation

Adopt `SQLDBProvider` for the migration, and retire the seven sites as each has a driver
counterpart, not as each engine cuts over:

- the three `sql_service` sites go with wave 2's `GetDB` retirement;
- `rollout_service` goes when it asks the driver for its activity instead of holding a conn;
- the two cleanups move to `driver.Execute` in the change that cuts mysql over — they are plain
  SQL, need no session, and then no engine needs a handle to clean up gh-ost;
- the validator keeps `SQLDBProvider` on mysql until omni exposes a binlog probe; a mysql adapter
  must return non-nil there, and a non-mysql driver reaching that code is a bug, so it should
  return an error rather than fail closed.

Cost: the ~25-line interface and helper plus seven small edits. The only real work is the binlog
probe, which is an omni mysql method, not a contract change.

### If nobody decides

`GetDB` stays on `db.Driver` and every adapter implements a method it cannot honour, so the nine
nil engines keep looking like defects instead of a category. The pin bump for mysql either keeps
the handle (the shim becomes permanent by accident) or drops it and gh-ost validation stays off
with nobody owning that regression. The cleanup paths either nil-deref or skip silently, leaving
`_gho`/`_ghc` tables behind.

## 3. The truncation cut point moved

### What the proofs measured

At a 1 MB cap, the same pg result now keeps 9,448 rows where Bytebase kept 16,384, with the same
message (`Output of query exceeds max allowed output size of 1MB`).

The cause is the measurement, not the cap. Bytebase builds the `v1pb.QueryResult` and compares
`proto.Size` of the whole message only at power-of-two row counts
(`backend/plugin/db/util/driverutil.go:137`), so the cap fires at the next power after the real
crossing and overshoots by up to ~2×. omni accumulates the wire-framed header and row sizes as it
reads (`driver.HeaderByteSize`, `driver.RowByteSize`, used in `pg/driver/query.go:266-287`) and
stops at the first row over the cap, setting `Result.Truncated`; the caller formats the same
message. The contract already decided that "the result size cap is measured by omni"
(driver-contract, Decisions).

### Options

- **Accept the row-accurate cut.** The cap now means what it says; a result can hold fewer rows
  than before when the cap binds.
- **Add caller-measured sizing to the contract.** That is not a method change: the two sides
  measure different objects (the whole API message vs header-plus-row wire bytes), so matching the
  old row counts needs the caller to supply both the accounting and the stop point — for example a
  `QueryOptions` hook that gets each row and answers stop/continue, or a caller-supplied size
  function. A `Result.Size()` method cannot fix it, because the driver decides when to stop.
  Cost: a contract addition, a per-driver read-loop change, and API knowledge inside `driver/`,
  which imports only the standard library and `metadata`.
- **Raise Bytebase's cap to compensate.** Ships more bytes than the policy allows; not a fix.

### Recommendation

Accept. Cost: zero code; one line in the release note and a test expectation where the row count is
asserted. The old number was the cap plus up to another cap; the new number is the cap. If a
customer-visible row count must be preserved for a specific report, cap that report's rows in the
API, not the transport.

### If nobody decides

Someone reads 9,448 against 16,384 as a regression and either reopens the contract's measurement
decision — the expensive option above — or raises the limit in Bytebase, which delivers more than
the query-data policy promises. Either way the pin bump lands a visible behavior change with no
record of why.

## 4. A NUL byte in the username fails at first use, not at `Open`

### What the proofs measured

Bytebase's CloudSQL IAM cleanup test (`TestCloudSQLIAMDialerClosesWhenOpenFails`) used a NUL byte in
the username so pgx's DSN parser rejected the connection string during `Open`, after the dialer had
been built. omni builds the DSN without the user and assigns it afterwards
(`pg/driver/driver.go:133-137`), so the NUL is never parsed: `Open` succeeds and the failure arrives
at first use. The pg proof changed the test to an empty user, which still fails after the connector
is built, so the cleanup-on-failed-open path stays covered.

### Options

- **Accept.** No legitimate credential contains a NUL; the error still surfaces on the first
  connection attempt, one layer later.
- **Validate in `pg/driver.Open`.** Restores the old timing for one pathological value and adds a
  caller-input rule to the driver — and, to be consistent, to every driver.
- **Validate in Bytebase's config assembly.** Catches it before `db.Open`, which keeps the rule
  product-side, but only for the paths that build `driver.Config` there.

### Recommendation

Accept, with the reason recorded here. If the pg driver is touched next, one sentence beside
`connConfig.User = cfg.User` saying the assignment deliberately skips DSN validation is worth more
than an added check: it stops the next reader from "restoring" the old failure order.

### If nobody decides

A future pin bump re-adds a NUL-user test and sees it fail, or adds byte validation to the pg
driver, which reorders `Open`'s failure modes and invalidates the CloudSQL cleanup proof. The change
looks like a bug fix, so it will not be recognized as a decision.

## 5. plancheck's eight concrete type assertions

### What the proofs measured

`runner/plancheck/statement_report_executor.go:167-226` has eight concrete assertions over nine
engine values (`MYSQL` and `MARIADB` share one) to reach `CountAffectedRows`. The mysql proof
replaced them with one `CountAffectedRows` interface: +24/−62, net about −38 lines, and it exposed
nine blank parser imports that the concrete driver imports had been providing transitively — without
them, the summary report fails with `engine POSTGRES is not supported for ParseStatements`.

The pg proof found the *forced* interface landed elsewhere: `advisor.Context.Driver` (a `*sql.DB`
that pg's live rules EXPLAIN through), not plancheck. The oracle proof then kept plancheck's
`*oracle.Driver` assertion with zero product diff, so the collapse is optional for every cutover so
far. Because an adapter keeps the engine package's `Driver` type and its `CountAffectedRows`
method, the old assertions still compile.

### Options

- **One interface (`db.RowsEstimator`) before the first wave-2 adapter lands.** A single change,
  reused by every remaining engine.
- **Later, engine by engine.** Each engine's proof repeats the exercise, and each rediscovers the
  parser-import coupling.
- **Leave it.** The type pin survives, so the concrete `Driver` type can never be removed from the
  engine packages, and the eight imports stay load-bearing invisibly.

### Recommendation

Do it once, in the opening change of wave 2 — before the first engine whose plancheck case exists
(mssql or cockroachdb) cuts over, not in the clean wave-1 cutovers. Cost: the measured
+24/−62 in `statement_report_executor.go`, about −38 net lines, plus the nine blank parser imports,
plus one interface declaration. Reuse a single name (`db.RowsEstimator`) rather than adding a second
near-interface beside the advisor's probe. The parser imports must be checked by running
`go test ./backend/runner/plancheck/`, not by compiling.

### If nobody decides

Every remaining wave-2 cutover keeps its concrete type alive only because plancheck names it, so the
"thin adapter" has a product-side tail no measurement counts. The parser-registration coupling will
surface once per engine as a runtime error in the SQL summary report, and the eight assertions keep
eight packages linked into plancheck for as long as the type pin exists.

## 6. Cassandra reports a missing keyspace as a sync failure

### What the proof measured

Bytebase's cassandra `SyncDBSchema` returned an empty schema for a keyspace the server did not have.
`cassandra/driver.SyncDatabase` answers `driver.ErrNotFound` — `not found: keyspace … does not
exist` — and a driver opened on a missing keyspace fails even earlier, at `Ping`. The sync runner
wraps the error (`backend/runner/schemasync/syncer.go:575-578`), so the database records a failed
sync whose message names the keyspace; the old driver recorded an empty success and overwrote the
stored schema with nothing.

The omni doc had not caught up: `cassandra/driver/sync.go` said a missing keyspace "answers an empty
schema, as Bytebase's did", while the code returns `ErrNotFound`. The proof fixed the comment in
`5fcf51383` on the unmerged `cassandra-cutover` branch; it landed on its own before the pin bump, so
the driver's promise matches what the adapter receives.

### Options

- **Keep `ErrNotFound`.** The sync fails, the message names the keyspace, and a later successful
  sync clears it. Someone deleted the keyspace, and the product says so.
- **Map it back to an empty schema in the adapter.** Restores parity in about three lines; the sync
  succeeds and the schema page shows a database with no tables, indistinguishable from one emptied
  on purpose.
- **Treat it as a missing database in the product.** Larger than a cutover: mark the database
  missing, stop syncing it, and decide what a diff or a task against it does.

### Recommendation

Keep `ErrNotFound`; cost zero adapter lines and one release note. A user sees `failed to sync
database schema for database X: not found: keyspace X does not exist` instead of a silently empty
schema. The one-line doc fix (`5fcf51383`'s change) landed before the pin bump, so the driver's
promise matches what the adapter receives. Map the error only if the product wants the old silence,
and record that a vanished keyspace will then look like an empty database.

### If nobody decides

The cassandra pin bump ships a new sync failure with no record of why, and the first support report
reads it as a regression. Or the next reader "restores parity" in the adapter, which hides a
deleted keyspace, and the omni comment keeps promising that behavior, so nobody can tell which side
is intended.

## 7. Cassandra's `Execute` now writes task-run log entries

### What the proof measured

The old cassandra `Execute` emitted no task-run entries. The adapter maps `db.ExecuteOptions` to
`driver.ExecuteLogger`, so a task run now records `COMMAND_EXECUTE`/`COMMAND_RESPONSE` for each
statement, as every other engine does. The change is visible in task-run history, not just in logs.

### Options

- **Keep it.** Cassandra stops being the one engine whose task-run history has no commands.
- **Suppress it for parity.** Pass a discarding logger; the adapter grows a special case and the
  history stays empty.

### Recommendation

Keep it. Cost zero; the old silence was the outlier, and suppressing it removes the only per-command
record a Cassandra task run could show.

### If nobody decides

A Cassandra task run still shows nothing, or an adapter suppresses the logger to match the old
behavior while every later engine disagrees — and the task-run log stops being a cross-engine
record.

## 8. What the Cassandra proof settled without a decision

| Item | Record |
| --- | --- |
| SSH through `Config.Dial` | A gain, not parity: the old driver had no SSH path. A tunnel-backed Cassandra source is intended; do not read it as a regression. |
| The no-op `Dump` | Delegated, not synthesized: both sides write nothing (`cassandra/driver/dump.go`, `TestDumpWritesNothing`). The contract's `Dumper` wording does not promise text, so an engine whose schema cannot be restored keeps returning nil. |
| The ignored read-only boolean | Knowingly preserved. The old `QueryConn` called `ValidateSQLForEditor` and ignored its read-only result, so the `"support Cassandra SELECT statement only"` branch is dead; the adapter keeps it. The real gate is `requiresAdminDataSource` (`backend/api/v1/sql_service.go:2464`, `:2502`), which routes a non-read-only statement to the admin data source. A fix belongs to that product gate, not to the adapter. |
| Wave order | Cassandra is the cleanest engine, not the smallest diff: 312 adapter lines, −490 non-test, zero product files, and zero engine-specific `Open` options (no `Params`, no `Options`, nil `GetDB`, ignored conn, no-op `Dump` on both sides) — while deleting the most code of any wave-1 engine, because it carried its own sync. It needed none of the wave-0 seams, which is why wave 0 does not gate wave 1 (dispatch §7). Wave 1 orders by fewest engine-specific decisions, not by adapter line count, so it stays first. |
| Verification | Bytebase's cassandra package had no server-backed test; the proof's live exercise was temporary and is not in the measured test lines. A cassandra pin bump should add one server-backed test at Bytebase's seam, or say that omni's tape suite is the only coverage. |
| The shared conversion helper | The first cutover inlined it, as pg did. Three proofs then wrote the same 120 lines, which is why the revised plan makes `db/util/omnidriver.go` wave 0; cassandra's inline copy is the fallback if wave 0 has not landed (dispatch §7, §8). A second adapter copying it is what measured the extraction. |

## 9. Oracle's private-link guard fails closed without a connection

### What the proof measured

`refuseLinkedReadsShadowedByPrivateLinks` (`backend/api/v1/sql_service.go:697`) type-asserts a
driver to `OwnsPrivateDatabaseLink(ctx, conn *sql.Conn)` (`:684-685`) and refuses when the assertion
fails or `conn == nil` (`:705-709`). After the cutover `db.SQLDB(adapter)` is nil for Oracle, so
keeping the signature refuses **every linked read** — a fail-closed regression that compiles and has
no error naming the cause.

omni already models the probe without a connection: `driver.PrivateDatabaseLinkOwner` is
`OwnsPrivateDatabaseLink(ctx) (bool, error)` (`driver/probe.go:7-15`). The oracle proof dropped
`conn` from the guard, delegated the probe in three adapter lines, and tested it live: the probe
reported false, then true after a private link was created, then false after it was dropped, while
the driver's session stayed in the database's schema — so it reads `SYS.USER_DB_LINKS` for the
session's account, not for `CURRENT_SCHEMA`. When the driver cannot answer, the refusal becomes
"the driver exposes no private database link probe" (code 13) rather than silence.

The probe runs on a session from the driver's pool, not the executing connection. That is safe
because Oracle commits DDL, so a link an earlier statement of the same request created is visible.

### Options

- **Take the handle off the guard now.** The interface becomes `OwnsPrivateDatabaseLink(ctx)`; the
  oracle driver opens its own session for `SELECT COUNT(*) FROM SYS.USER_DB_LINKS` until its
  adapter lands. No omni dependency, so it can land before the pin.
- **Pass nil and keep the signature.** Every Oracle linked read is refused after cutover.
- **Keep a `*sql.DB` on the oracle adapter just for this probe.** Reopens the `GetDB` retirement
  for one caller and keeps the shim alive for Oracle alone.

### Recommendation

Take the handle off the guard, as the proof did. Cost: one interface signature, one caller (the
`sql_service` edit the oracle proof measured), a self-session method on today's `backend/plugin/db/oracle`
driver, and two guard tests — one table-driven refusal test and a package-level assertion that the
driver implements the interface. It is the only product change unique to Oracle, and it removes the
last way the guard can fail closed silently when a driver has no SQL handle.

### If nobody decides

After the oracle pin, reads on any database with a private link are refused with no message that
points at the guard, which reads as an outage of linked queries. Or someone keeps `GetDB` alive on
the oracle adapter only, and the optional-handle decision in §2 is undone one engine at a time.

## 10. Oracle's `GetVersion` is not needed (settled)

Bytebase's `GetVersion` returned `connectionCtx.EngineVersion`, which `dbfactory.go:56` set from
`instance.Metadata.GetVersion()` — the same string. omni's `SyncDatabase` asks the server
(`version_full`) and falls back to `Config.EngineVersion`, so the adapter needs no version method.
The one caller, the prior-backup DML transform, parses `instance.Metadata.GetVersion()` with
`plsqlparser.ParseVersion` instead of asserting `*oracle.Driver`; no interface and no new method.

Recorded so nobody re-adds a driver version accessor for that caller. If the stored version is ever
empty, the contract's fallback already covers sync; the transform's `(\d+)\.(\d+)` parse matches the
metadata string and the old banner identically.

## 11. `MaximumSQLResultSize` of 0 means the opposite on each side

### What the proof measured

Bytebase's `MaximumSQLResultSize` becomes `driver.QueryOptions.MaxResultBytes`, and the two sides
read zero differently. Bytebase's driver compared the running size against the cap as given, so 0
meant a cap of zero bytes and the result came back truncated. The contract reads 0 as no cap
(`driver/query.go:38-40`, and every reader tests `maxResultBytes > 0`: `mysql/driver/query.go:269`,
`doris/driver/query.go:187`), so the same configuration now returns every row.

Nothing breaks today: the result policy always sets a positive value. The difference is latent, and
its direction is the bad one — a policy row that loses its size turns a zero-byte cap into an
unbounded read, so a policy bypass surfaces as a query that suddenly returns more rows.

### Options

- **Refuse a non-positive value in the adapter.** `maximum SQL result size must be positive`: the
  product never sends 0, so nothing that works today stops working, and the fail-closed direction is
  kept.
- **Pass it through.** The contract's answer is honest, and a caller that meant "no rows" has no way
  to say it. Invisible until a policy value goes missing.
- **Change the contract.** A separate `NoCap` field or a pointer distinguishes the two meanings, at
  the price of a contract change for a value no caller produces.
- **Map 0 to the old cap in the adapter.** Pretends the driver can cap at zero; the closest it can
  come is refusing to read the statement, which is the first option with an extra fiction.

### Recommendation

Refuse it in the adapter, three lines and one error string. It is the only option that neither
changes the contract for a case that does not occur nor turns a missing cap into no cap, and it is
the mirror of the direction the mssql proof took for an unsupported mode. Record the contract's
0-means-no-cap as deliberate: the difference in §3 is about *where* a positive cap cuts, and this
one is about what 0 asks for, so neither subsumes the other.

### If nobody decides

The adapter passes 0 through, and the first workspace whose result policy has no size returns every
row of a query while the operator believes a cap is set. No test catches it, because every fixture
sets a positive value, so it reads as a regression in the policy rather than in the mapping.

## 12. Doris's `Dump` was not replayable, and one parser fixture moves

### What the proof measured

Two changes the doris pin bump owes, neither of them a choice:

- **The `Dump` output.** The `SET @OLD_UNIQUE_CHECKS …` prologue and epilogue are gone, the
  placeholder `CREATE VIEW … 1 AS col` that stood in for a view definition is gone, and procedures,
  functions and sync materialized views now appear. The old text could not be replayed at all, so
  this is a fix and not a format move, and it belongs in the release note beside the pin rather than
  in a support case.
- **One parser expectation.** `backend/plugin/parser/doris/limit_test.go:197` holds
  `SELECT * FROM t # hash comment`. Doris refuses `#`, so omni's limit rewrite wraps the statement
  instead of splicing the limit in front of the comment, and the expected text changes with the pin.

### Still open: a procedure created qualified syncs qualified

A procedure created as `probe.p1` syncs as `probe.p1`, the server's own text, where `mysql/driver`
records the bare name. Changing it is a `doris/driver` edit plus a re-recorded sync golden
(`OMNI_RECORD=1`), so it is a driver decision rather than adapter code: the doris pin bump either
normalizes the name or records that a synced procedure keeps its qualification. Normalize it when a
consumer needs the bare name — a diff against another engine's schema, or a lookup by procedure
name; until then the server's text is the truer record.

## Settled, not yours to choose

These were decided by the contract or the dispatch plan and the proofs confirmed them; a pin bump
that reopens them is a different change:

- The registry, `db.Register`/`db.Open`, the blank imports and `dbfactory` do not move; a cutover
  replaces one engine package and the caller conversions around it.
- The five interfaces, `driver.Config`, and the caller-side conversions in a shared `db/util`
  helper stay; trino and mysql both measured about 120 lines for the result conversion.
- The pin bump changes product output, not just plumbing: trino's dump gains a `CREATE SCHEMA IF NOT
  EXISTS` prefix and views (omni `#1280`), and its `DECIMAL` arrives as text because `driver.Value`
  has no exact-decimal kind. Neither is a cutover choice. Doris's dump is the same category with the
  sign reversed: the pin bump replaces text that could not be replayed (§12).
- No omni driver behavior gap surfaced in the five landed proofs. The one omni defect is the
  `cassandra/driver` comment that promises the opposite of its code (§6); the gh-ost binlog probe
  and the plancheck interface are Bytebase product work, and the truncation difference is a
  contract decision already made. **Superseded 2026-11-17 by #1387:** that comment was corrected on
  main (`3f1e34d8d`), so main no longer carries the false claim §6 describes.

## Still in flight

**Superseded 2026-11-17.** The mssql proof has reported: the dispatch plan's §3 carries its measured
cost (497 adapter lines, 8 caller files, all of them the shared `GetDB` retirement) and §6 its five
driver defects and the moved `Dump` text. Its Azure AD half is no longer a contract question — #1389
added `driver.Config.TokenProvider`, which `mssql/driver` hands to the client's security-token
connector, and the adapter passes a closure over `azidentity` returning the
`database.windows.net` token; the
[contract's Decisions](2026-11-15-driver-contract-design.md#decisions) and the
[parity note](2026-11-17-driver-method-parity.md#dialer-and-probe-support) own that wording.
The masking question below is still unrecorded: the contract's Decisions still list `MaskingEnabled`
as open for the query slice.

The mssql proof had not reported when this was written, so none of its choices are recorded here.
The dispatch plan's measured section already carries two of its facts: masking and batch delegation
need no caller change, and `mssql/driver` cannot take Azure AD's token-fetching connector, so the
adapter refuses that mode loudly. The second is a contract decision, not an adapter shortcut; the
report owes its options and cost. Its session names one more question of this kind: which consumers
set `MaskingEnabled`, and what a procedure call whose rows cannot be traced to columns returns to
each of them. Do not read the dispatch plan's summary as a decision. If the report lands before the
pin bump, it belongs here, one section, in the same shape.
