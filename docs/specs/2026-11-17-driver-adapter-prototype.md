# Driver adapter prototype: what the contract costs a caller

2026-11-17 · Audit

Bytebase `345007290e` and omni `4ca314ae6`. In a scratch Bytebase worktree (never committed) the omni
pin was replaced with `origin/main` and each engine's `backend/plugin/db/<engine>` was replaced by a
thin adapter over `<engine>/driver`, then compiled and vetted. Nineteen engines: pg, mysql, mariadb,
tidb, oceanbase, mssql, oracle, trino, mongo, spanner, doris, starrocks, cassandra, elasticsearch,
cosmosdb, dynamodb, redis, clickhouse, hive.

The question was whether the cutover is mechanical. The answer is yes for eighteen of them, with the
seams below; clickhouse is the one engine whose `Dumper` has no omni counterpart.

## What an adapter is

Every adapter is the same nine methods, and each one is one call:

| Bytebase method | Adapter body |
| --- | --- |
| `Open(ctx, engine, ConnectionConfig)` | resolve the credential, tunnel and TLS material, then `<engine>.Open(ctx, driver.Config)` |
| `Close(ctx)` | `Driver.Close()` plus the tunnel or cloud dialer Bytebase opened |
| `Ping(ctx)` | `Driver.Ping(ctx)` |
| `Execute(ctx, script, opts)` | parse `-- txn-mode`, then `Executor.Execute(ctx, script, trinoOptions(opts))` |
| `QueryConn(ctx, conn, script, qc)` | `Querier.Session` + `Session.Query` + `driver.Result` → `v1pb.QueryResult` |
| `SyncInstance(ctx)` | `Syncer.SyncInstance` → `db.InstanceMetadata` (roles, `lower_case_table_names`) |
| `SyncDBSchema(ctx)` | `Syncer.SyncDatabase` (both sides already hold `omni/metadata`) |
| `Dump(ctx, w, md)` | `Dumper.Dump` |
| `GetDB()` | nothing to return; see below |

## Size

Lines in the prototype adapter, one file per engine. `config` is `Open` and the helpers beside it,
`boiler` is the nine methods above, `extra` is the engine's own methods (`CountAffectedRows`,
`ExplainJSON`, `ColumnExtra`, `OwnsPrivateDatabaseLink`, …).

| Engine | Total | Config | Boilerplate | Extra |
| --- | --- | --- | --- | --- |
| pg | 226 | 87 | 50 | 3 (+25 pg error mapping) |
| mysql | 136 | 15 | 44 | 15 |
| mariadb | 136 | 15 | 44 | 15 |
| tidb | 125 | 15 | 44 | 9 |
| oceanbase | 130 | 15 | 44 | 12 |
| mssql | 188 | 83 | 39 | 3 |
| oracle | 143 | 44 | 36 | 6 |
| trino | 101 | 27 | 31 | 0 |
| mongo | 144 | 56 | 39 | 0 |
| spanner | 120 | 36 | 36 | 0 |
| doris | 100 | 15 | 39 | 0 |
| starrocks | 100 | 15 | 39 | 0 |
| cassandra | 126 | 39 | 39 | 0 |
| elasticsearch | 169 | 61 | 31 | 0 (+11 SigV4 transport) |
| cosmosdb | 143 | 45 | 41 | 0 |
| dynamodb | 133 | 54 | 31 | 0 |
| redis | 154 | 52 | 53 | 0 |
| clickhouse | 111 | 27 | 36 | 0 |
| hive | 124 | 46 | 31 | 0 |
| **Total** | **2609** | **747** | **747** | **110** |

Nineteen adapters replace ~29.5k lines of Bytebase driver code plus a 6.5k-line copy of the Cosmos DB
query engine that ends up imported by nothing. 747 of the 2,609 lines are the same nine methods
repeated; they exist only because Bytebase's own interface keeps `*sql.DB`, `*sql.Conn` and
`v1pb.QueryResult` in its signatures.

The caller-side conversion every adapter shares came to 542 lines besides tests: `driver.Result` →
`v1pb.QueryResult` 127, `driver.Instance`/`Role` → `db.InstanceMetadata`/`storepb.InstanceRole` 29,
the `ExecuteLogger` bridge to task-run logs 36, the TLS/SSH/dial/application-name and
`QueryOptions` helpers 48, the MySQL family's config builder 158 (82 of it is the RDS CA bundle
Bytebase already has), the AWS RDS token 15, and declarations and package doc for the rest. Six
engines (mysql, mariadb, tidb, oceanbase, doris, starrocks) reach that one config builder, which is
why it belongs in a shared package rather than in each adapter.

## Seams the prototype found

1. **`GetDB` has no counterpart, and a cutover compiles anyway.** Eighteen of nineteen adapters
   return a nil `*sql.DB` and the repository still builds: only three packages fail to compile
   (the MySQL/MariaDB advisors' package-level helpers and Oracle's `GetVersion`). The remaining
   `GetDB` users — the SQL editor's three session paths, rollout session listing, the ghost
   validator's five MySQL probes, gh-ost cleanup, and the 16 advisor files that take a `*sql.DB` as
   their probe connection — fail at run time, not at build time. That is the one place where a
   cutover cannot be validated by compiling.
2. **`QueryConn(*sql.Conn)` cannot honour its argument.** The contract hands out its own `Session`,
   so the adapter opens one per call and the caller loses the session it was reusing. Bytebase's
   editor holds a connection across statements; a cutover must hold a `driver.Session` instead, or
   pass `Query()` calls through one adapter-held session per editor session.
3. **`Command`'s byte range excludes the trailing semicolon.** omni's `mysql/parser.Split` (and the
   splitters built the same way) end a segment at the `;`, and `Segment.Text` says so; Bytebase's
   `NewStatementFromRange` extends `storepb.Range.End` past it. A task-run's command log therefore
   highlights one byte less per statement unless the caller extends the range itself. Nothing in
   `driver.Command` carries the script, so an adapter cannot fix this on its own.
4. **The transaction mode and isolation level are the caller's.** Every Bytebase driver parses its
   own `-- txn-mode`, with `transaction.DefaultMode()` when the script states none. `driver.
   ExecuteOptions` wants `Transaction` and `Isolation` filled in, and nothing in the contract parses
   the directive. The prototype has each adapter call `base.ParseTransactionConfig` itself, which
   keeps a parser dependency inside a driver package — the wrong side of the contract, and the reason
   `ExecuteOptions` needs the caller to do it.
5. **Engine-specific endpoints are `Params` the caller must know.** The contract has no field for
   any of these, and each is a fact only the caller has: spanner `project`/`instance`; dynamodb
   `region`/`sessionToken`/`accountId` (account ID is an STS call, which the caller must make);
   cosmosdb `auth` (`masterkey`/`clientsecret`/`default`), `tenantId` and `emulator`; hive
   `auth`/`service`; oracle `SERVICE_NAME` or `SID`; mssql `certificate` (a path — the caller writes
   the CA to a file); clickhouse settings unless they are the defaults; trino catalog via `Database`
   and the rest as query parameters. Redis is the only engine that needed a second `Open` argument
   for topology (`Options.Type`, `MasterName`, sentinel credentials).
6. **Cosmos DB's container is `QueryOptions.Schema`.** Bytebase's `QueryContext` has a separate
   `Container` field; the adapter remaps it. Any other caller needs to know the overload.
7. **Redis's all-nodes query is outside the interface.** Bytebase's editor offers "run on every
   node", which only `*redis/driver.Session.QueryAllNodes` answers; `driver.Session` does not declare
   it, so the adapter type-asserts the concrete session.
8. **`MaxResultBytes` moved the cut-off.** Bytebase measured `proto.Size` of the whole `QueryResult`
   at power-of-two row counts and answered an error message; omni measures the rows as it reads them
   and sets `Result.Truncated`, and the caller formats the message. The cap lands at a different row
   for the same setting. Every engine's package doc says so, but it is a visible product change.
9. **`clickhouse` has no `Dumper`.** Bytebase's `Dump` asks the server for `CREATE` statements and
   ignores the metadata it is given, so it does not need a `clickhouse/catalog`. omni's
   `clickhouse/driver` has no `Dump` at all, so this is the one method of the five interfaces a
   caller cannot delegate today. Hive, Doris and StarRocks dumps are the same shape (server-stated
   DDL, metadata ignored) and are implemented, so clickhouse is an omission rather than a
   consequence of the missing catalog. Note that the contract's wording — "a dump renders the
   metadata it is given and nothing else" — is not what those four drivers do.
10. **Two engines cannot take a caller's dialer.** `hive/driver` says `Config.Dial` does nothing
    (gohive dials the host itself) and `elasticsearch/driver` says so too (the caller's `Options.
    Transport` is the seam instead). A Hive data source behind an SSH tunnel cannot be cut over, and
    Elasticsearch's AWS SigV4 transport is 40 lines of caller code rather than a `Config` value.
11. **`oracle/driver` exports `Version` with a different shape.** Bytebase's `GetVersion()`
    (`*plsql.Version`, no error, from the recorded version) becomes `Version(ctx)`, which asks the
    server and falls back. The one call site in Bytebase's DML-backup transform becomes an interface
    assertion with an error to handle. Cheap, but not "one line of registration".

## Verdict per engine

Fits with a config mapping and nothing else: trino, doris, starrocks, cassandra, cosmosdb, dynamodb,
redis.

Fits with a documented limitation: hive (no caller dialer), elasticsearch (a transport instead of a
dialer, and the SigV4 signer is the caller's).

Fits, and the adapter must also expose the extra methods its callers use: pg, mysql, mariadb, tidb,
oceanbase, mssql, oracle, spanner. Spanner's `Open` reaches the server for the dialect, which its own
doc records.

Does not fit yet: clickhouse, whose `Dumper` is missing.

Pending rather than prototyped: cockroachdb (its driver is in flight), and redshift, snowflake,
bigquery and databricks (blocked).

## Refresh

Re-run from Bytebase's `backend/plugin/db/<engine>` against omni's `<engine>/driver` with the pin
replaced, one engine at a time, and compile `go build ./backend/plugin/db/<engine>/...`. The
adapter line counts above are the measurement of a caller's cost, not a target.
