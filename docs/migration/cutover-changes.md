# Cutover changes: what users see differently

2026-11-17 · Release note for the Bytebase pin bump

This is the user-facing change list for moving Bytebase's engine packages onto the omni drivers. It
is not a migration plan: the dispatch plan
([2026-11-17-cutover-dispatch.md](../specs/2026-11-17-cutover-dispatch.md)) says how each engine
reaches its driver, and the open decisions
([2026-11-17-cutover-open-decisions.md](../specs/2026-11-17-cutover-open-decisions.md)) say which
choices remain. This page answers one question: **what will a user see differently after the pin
moves?**

The measurements come from fourteen cutover proofs run against the pinned servers: trino, pg, mysql,
cassandra, oracle, mssql, redis, spanner, elasticsearch, doris, tidb, clickhouse, dynamodb and hive.
Each table is ordered fixes first, then changes that look the same but behave differently, then
items that need a product decision. The last class is where regression risk lives: not every row is
a regression, but each one needs an owner before the pin moves. "Pinned" means a test, fixture or PR
that holds the new behaviour; "proof only" means the behaviour was measured in a scratch tree and no
test in the repository pins it yet.

## At a glance

| Engine | Report | What to expect | Decisions left |
| --- | --- | --- | --- |
| hive | measured | mostly fixes; two pre-existing failures disappear | SSH refusal, materialized-view definition |
| mssql | measured | fixes; dump fidelity and verified TLS improve | none engine-specific |
| doris | measured | fixes; the old dump could not be restored | zero cap, qualified procedure name |
| tidb | measured | fixes; partition dumps become restorable | zero cap, parse-error text, advisor/ghost |
| clickhouse | measured | fixes; EXPLAIN ANALYZE and server DDL improve | zero cap, false ROLLBACK log |
| dynamodb | measured | fixes; limits and typed values work | TLS default, role-key lifetime |
| oracle | measured | fixes and capability gains | quoted key columns, private-link signature |
| trino | measured | fixes, plus a dump-format move from the pin | synthetic `Affected Rows 0`, `SHOW` row cap |
| cassandra | measured | same-looking changes and one new sync failure | missing keyspace, dead read-only branch |
| redis | measured | same-looking changes; timeouts and refusal semantics move | none; batching is a caveat |
| elasticsearch | measured | same-looking changes; bulk failures now fail | AWS error hint |
| spanner | measured | same-looking changes, plus a half-landed dialect capability | PostgreSQL-dialect listing |
| pg | measured | same-looking changes; row counts and explain shape move | truncation is settled; no new choice |
| mysql | measured | fixes, but the largest set of product-visible choices | gh-ost, masked CALL shape, zero cap |
| mariadb | proof running | partial; the mysql rows apply and `RETURNING`/`CALL` rows return | same as mysql |

## Engine tables

### Cassandra

| Behaviour | Before | After | Class | Where it is pinned |
| --- | --- | --- | --- | --- |
| Sync on a keyspace the server no longer has | `SyncDBSchema` returned an empty schema and the sync succeeded | `driver.ErrNotFound`: `not found: keyspace X does not exist`; the sync runner records a failed sync | needs a product decision | `cassandra/driver` `TestSyncDatabaseMissing`; the doc correction is omni#1387 |
| `Execute` task-run history | No command entries | One `COMMAND_EXECUTE`/`COMMAND_RESPONSE` per statement | same-looking-different | Proof only |
| SSH data sources | No SSH path | `Config.Dial` gives a working tunnel | fix | Proof only |
| The editor's read-only check | `ValidateSQLForEditor`'s result was ignored, so a write through the query path reached the server | The same dead branch is preserved; the real gate remains `requiresAdminDataSource` | needs a product decision | Proof only; this is pre-existing, not introduced by the cutover |

### ClickHouse

| Behaviour | Before | After | Class | Where it is pinned |
| --- | --- | --- | --- | --- |
| `EXPLAIN ANALYZE` through the query path | Plan came back but the rows the server returned were discarded | Returns the plan and its rows as a typed `QueryPlan` | fix | omni#1394 (`clickhouse/driver/query_test.go`) |
| Date/time and address values | `DateTime*` was an RFC3339 string, `Time*` was NULL, `IPv4`/`IPv6`/`FixedString` were NULL | Typed timestamps with their real accuracy, clock strings, and the actual address/string values | fix | Proof only |
| `SHOW`/`DESCRIBE` in the editor | Fell through as a write and showed `Affected Rows: 0` | Runs as a query and returns rows | fix | Proof only |
| Stored schema dump | Written from `system.tables.create_table_query`: own-db prefixes stripped, `.inner_id.*` storage tables written, MVs before their targets | `SHOW CREATE` text: database-qualified, dependency-ordered, storage tables excluded, dictionaries and materialized views labelled | fix | Proof only; round-trip `dump → restore → dump` is byte-identical |
| Extra parameters and SSH | Ignored | Reach ClickHouse as settings; tunnelled data sources work | fix | Proof only |
| A statement with no rows | A synthetic `Affected Rows: 0` row was added | Empty result; `Execute` still reports 0 | same-looking-different | Proof only |
| Result-size cut | `proto.Size` measured at power-of-two row counts | Byte-accurate cut; same message, possibly fewer rows | same-looking-different | Proof only |
| Task-run log on a script that cannot roll back | Logged `BEGIN`/`ROLLBACK` even though ClickHouse does not roll back | Unchanged: the adapter keeps the same false `ROLLBACK` entry | needs a product decision | Proof only; fix the log wording at the contract/log layer |
| `MaximumSQLResultSize == 0` | A cap of zero bytes, so one row plus an error | No cap, so every row | needs a product decision | Contract reads 0 as no cap; see the action list |

### Doris

| Behaviour | Before | After | Class | Where it is pinned |
| --- | --- | --- | --- | --- |
| Stored schema dump | `SET @OLD_UNIQUE_CHECKS …` prologue, placeholder views, no procedures; the text was not replayable | No prologue/epilogue; dependency-ordered definitions; procedures, functions and sync materialized views included; restores | fix | `doris/driver/dump_test.go` (`TestDumpRestores`) |
| `SyncInstance` version | `5.7.99` from the MySQL-compatible handshake | `3.1.4` from `@@version_comment` | fix | Proof only |
| `SyncInstance` databases | `mysql`, `probe`, `__internal_schema` | Only the user database | fix | Proof only |
| `DELIMITER` scripts | The directive reached the server and was refused | The script is split and the procedure is created | fix | Proof only |
| `-- txn-isolation` | Ignored | Passed through; an engine-refused level fails at `Execute` | fix | Proof only |
| `Execute` affected rows | The last statement's count | Summed over all statements, with per-result counts | fix | Proof only |
| The `#` comment limit case | The parser spliced `LIMIT` before a `#` comment Doris does not accept | Wraps the statement in a subquery instead | fix | `backend/plugin/parser/doris/limit_test.go:197` needs its expected text updated |
| Task-run history | Transaction control only | One command entry per statement, with byte ranges | same-looking-different | Proof only |
| Empty script | Logged `BEGIN` + `COMMIT` | Logs nothing | same-looking-different | Proof only |
| Result-size cut | Power-of-two row counts | Byte-accurate cut; 4,096 → 3,973 rows at 1 MB, same message | same-looking-different | Proof only |
| `MaximumSQLResultSize == 0` | A zero-byte cap (one row plus `0MB` error) | No cap | needs a product decision | Contract reads 0 as no cap; open decision #11 |
| A procedure created as `probe.p1` | Synced with the server's qualified name | Still synced qualified, while `mysql/driver` records the bare name | needs a product decision | Open decision #12; changing it needs a re-recorded sync golden |

### DynamoDB

| Behaviour | Before | After | Class | Where it is pinned |
| --- | --- | --- | --- | --- |
| `Limit` | Treated as a per-page limit: `Limit=1` returned 2 rows | Caps the statement's rows | fix | Proof only |
| `MaximumSQLResultSize` | Not enforced | Enforced; `Truncated` carries the usual message | fix | Proof only |
| Binary values | Raw bytes in a proto string, which could break protojson | `bytesValue`, with base64 inside documents | fix | omni#1273 |
| TLS validation | An `https` endpoint always verified; `verify_tls_certificate` was ignored | The adapter must keep verifying until the product can send an explicit opt-out | needs a product decision | omni#1395 (`dynamodb/driver/client_test.go`) |
| Synced metadata | Secondary-index key attributes were missing; a hash-only table gained a phantom empty `RANGE` index | Index key attributes and `GLOBAL`/`LOCAL SECONDARY` entries appear; the phantom index is gone | same-looking-different | Proof only; existing databases gain columns/index rows at the next sync |
| A missing attribute | A nil row-value entry | Explicit `NULL_VALUE` | same-looking-different | Proof only; confirm the grid renders it as a null cell |
| STS lookup | On `Open`/each editor query in the old flow | Once per `SyncInstance`, never on query/execute/schema sync | same-looking-different | Proof only; avoids an STS round trip and permission on every query |
| An assumed role's keys | A refreshing cache for the driver's life | Frozen per `Open`, bounded by the sync deadline and the role session | same-looking-different | Proof only; a run longer than the role session would need an engine-neutral credential source |

### Elasticsearch

| Behaviour | Before | After | Class | Where it is pinned |
| --- | --- | --- | --- | --- |
| A partially applied bulk | 200 with `error=""`, so the task run succeeded | `1 of 2 bulk items failed: …`; the task run fails | fix | Proof only; add a release note |
| `HEAD` or an empty answer | `failed to parse json body: unexpected end of JSON input` | Empty result with no columns or rows | fix | Proof only |
| Result-size cap | Ignored | Enforced at the cap, with the legacy message | fix | Proof only |
| Explaining a document write | Profiled and ran the write | Refused: only `_search` has a plan to explain | fix | Proof only |
| Contradictory TLS config | `use_ssl` plus an explicit `http://` host silently used plain HTTP | Refused | fix | Proof only |
| Sync ordering | Indices in server order, roles in map order | Both sorted by name | fix | Proof only |
| Object-answer columns | Map order, different per answer | Sorted | same-looking-different | Proof only; export headers become deterministic |
| `Execute` task-run log | No command entries | One `COMMAND_EXECUTE`/`COMMAND_RESPONSE` per request | same-looking-different | Proof only |
| AWS ping/sync error text | `authentication failed: unable to assume role …` and a path hint | `ping failed: [403 Forbidden] <body>` | needs a product decision | Proof only; preserve the hint in the adapter if support relies on it |

### Hive

| Behaviour | Before | After | Class | Where it is pinned |
| --- | --- | --- | --- | --- |
| Opening a data source with a username | SASL PLAIN ran with an empty login: `Bad SASL negotiation status: 3 (Error validating the login)` | The configured user/password are sent and the server accepts them | fix | Proof only; Bytebase has no Hive tests |
| `SyncDBSchema` | Failed on the first table: `sql: Scan error … converting NULL to string` on `data_type` | Tables, views, external tables and materialized views sync | fix | Proof only |
| Materialized view in a dump | Un-replayable `CREATE MATERIALIZED VIEW … AS` with an empty body | Restores, but as an external table; the materialized-view definition is still the documented `TestDumpGaps` gap | fix | Proof only; the definition read is on the unmerged `gap-audit` branch |
| Dead dump paths | Wrote `CREATE INDEX`, which sync never filled, and synthesised MVs | Removed; server `SHOW CREATE` only | fix | Proof only |
| `DOUBLE`/`FLOAT` cells | 20-decimal text | `doubleValue` with the value, so CSV/export text changes from `1.500…` to `1.5` | fix | Proof only |
| `QueryContext.Schema` | Ignored | `USE` before the script | fix | Proof only |
| Extra parameters, SSL, Kerberos | Not all reached the session | Reach it; Kerberos still maps to `auth=KERBEROS`/`service=<Primary>` | fix | Proof only |
| Dump text | Unqualified names after a `USE`; `-- TABLE structure`; semicolon on its own line | Qualified names; `-- Table structure`; semicolon on the last line; blank line between blocks | same-looking-different | Proof only; statements are otherwise identical |
| A failing statement in a script | The whole RPC failed | Earlier results come back and the failure rides on its `QueryResult` | same-looking-different | Proof only |
| `Execute` task-run history | No command entries | One entry per statement | same-looking-different | Proof only |
| `GetDB` | A handle | nil; no Hive caller dereferences it | same-looking-different | Proof only |
| A script against a database dropped under an open driver | Worked from the old connection state | Fails at `USE <database>` | same-looking-different | Proof only; the product always opens on an existing database |
| Result-size cut | `proto.Size` at power-of-two counts | Estimated per-row bytes; at 64 bytes, 2 rows became 3 | same-looking-different | Proof only; the cap is approximate either way |
| An SSH data source | Opened and pinged the direct host, silently ignoring the jump host | Refused with the port-forward instruction | needs a product decision | Proof only; a local forwarder in `backend/plugin/db/util` would restore tunnels |

### Microsoft SQL Server

| Behaviour | Before | After | Class | Where it is pinned |
| --- | --- | --- | --- | --- |
| Verified TLS | `x509: certificate is valid for localhost, not placeholder` because the TLS config verified the DSN placeholder | Verifies the real server name; verified-TLS instances connect | fix | omni#1385 (`mssql/driver` live suite) |
| `Dump` of a filtered index | Dropped the `WHERE` predicate | Keeps `WHERE ([amount] IS NOT NULL)` | fix | omni#1391 (`mssql/driver/dump_test.go`) |
| `Dump` after a view | Wrote a doubled `;` | Exactly one `;`; 921 → 948 bytes on the same schema | fix | omni#1391 |
| `SET`/`DECLARE` statement text | `results[0].GetStatement()` was empty | The statement text is kept | fix | omni#1385 |
| Procedure result sets after a batch | Returned after their batch instead of after the script | Returned in script order | fix | omni#1385 |
| `EXEC` limit and a cut result | `Limit` unenforced; a `MaxResultBytes` cut swallowed the rest of the batch | Limit enforced; the cut result is read and the batch still matches its statements | fix | omni#1385 |
| Masking-withheld rows | Set `Result.Err`, so later batches never ran | Withheld rows no longer look like a failed statement; later batches run | fix | omni#1385 |
| Azure AD logins | Refused loudly: the driver always built the plain connector | `Config.TokenProvider` carries the token-fetching login, so the adapter can pass an `azidentity` closure | fix | omni#1389; contract decisions |

### MySQL

| Behaviour | Before | After | Class | Where it is pinned |
| --- | --- | --- | --- | --- |
| Unsigned integers | Returned as strings, e.g. `"18446744073709551615"` | `uint64Value` | fix | omni `mysql/driver` tests |
| Procedure result sets | A `CALL` that answers with rows could lose them | The procedure's result set is returned | fix | Proof only |
| Secondary-index comments in a dump | Written as a separate `COMMENT ON INDEX` line | Inline in `CREATE TABLE`, which is valid MySQL | fix | omni `58f05cf7a`; `backend/plugin/aireview/catalog_test.go` expects the old line and needs updating |
| RDS IAM with cleartext passwords | The extra parameter was not passed through | `allowCleartextPasswords=true` reaches the client | fix | Proof only; no RDS endpoint tested |
| Script parsing | Bytebase split and ran one statement at a time, so a parse error arrived on that statement's result | omni parses the whole script before sending anything; a parse error fails the call | same-looking-different | omni `mysql/driver/query.go` tests |
| Result-size cut | A size error | Partial rows with `Truncated` and the same message | same-looking-different | Proof only |
| Session state after `Close`/`Execute` | A pooled connection could carry the selected database, variables or temporary tables to the next caller | The physical connection is discarded; no session state leaks | same-looking-different | omni `mysql/driver` `TestSessionCloseDiscardsTheConnection` |
| A masked `CALL p()` | Rows were withheld, but the result looked like an ordinary empty-rows result | The contract's count shape, `Affected Rows: 0`, with no rows | needs a product decision | Proof only; confirm the grid and support copy |
| `MaximumSQLResultSize == 0` | A zero-byte cap | No cap | needs a product decision | Contract reads 0 as no cap; see the action list |
| gh-ost binlog validation | Five probes ran through `GetDB` | Fails closed: `binlog access validation needs a driver with a SQL handle` | needs a product decision | Proof only; open decision #2; the contract has no binlog probe |

### Oracle

| Behaviour | Before | After | Class | Where it is pinned |
| --- | --- | --- | --- | --- |
| `Execute` task-run history | No command entries | `COMMAND_EXECUTE`/`COMMAND_RESPONSE` per statement | fix | Proof only |
| SSH data sources | No SSH path | `Config.Dial` gives a working tunnel | fix | Proof only |
| Sequence definitions in a dump | Omitted `MINVALUE 1 CACHE 20` | The server's full sequence definition | fix | `TestDumpStaysPairedWithOmniSync` in the proof |
| A dropped schema | An empty schema | `driver.ErrNotFound`; the sync records a failure instead of overwriting the schema with nothing | same-looking-different | Proof only; same sync-runner effect as Cassandra |
| A synced key's columns | Unquoted in Bytebase's catalogue | omni sync stores them quoted and prints them as-is, so Bytebase's catalogue renders `PRIMARY KEY (""ID"")` | needs a product decision | omni#1388 records the quoted key; no test crosses the two catalogues |
| The private-link guard | Took a `*sql.Conn`; the cutover makes that nil, which would refuse every linked read | Takes no connection; the driver answers from its own session | same-looking-different | Proof guard tests; the signature change is product-side |
| `GetVersion` | Returned the stored engine version | Deleted; sync reads `version_full` with a stored-version fallback, and the prior-backup transform parses the metadata | fix | Proof only; no new interface |

### PostgreSQL

| Behaviour | Before | After | Class | Where it is pinned |
| --- | --- | --- | --- | --- |
| Result-size cut | 16,384 rows at a 1 MB cap | 9,448 rows at the same cap, same message | same-looking-different | Open decision #3, settled by the contract; no test pins the row count |
| `QueryPlan` on an explain request | A caller's own `EXPLAIN` could switch the plan renderer's format | The adapter suppresses `Result.Plan` when `QueryContext.Explain` is set | same-looking-different | Proof `TestQueryConnExplainFormat` |
| NUL byte in a username | `Open` failed inside DSN parsing | `Open` succeeds and the first use fails | same-looking-different | Open decision #4, settled: accept |

### Redis

| Behaviour | Before | After | Class | Where it is pinned |
| --- | --- | --- | --- | --- |
| Context deadline | Only began to interrupt a command inside the client: a 200 ms deadline took 5.0 s | Interrupts the command: 201 ms | fix | Proof only |
| Cluster database selection | Accepted a non-zero or negative database and silently ran on db 0 | Refuses: `a Redis Cluster has database 0 only` / `a Redis database index cannot be negative` | fix | Proof only |
| TLS per node | One config for the connection | Each dialled host gets its own `ServerName`, so cluster/sentinel nodes are verified individually | fix | Proof only; no TLS fixture |
| Empty host/port | Built `host:` | Port defaults to 6379; an empty host without a dialer is refused | fix | Proof only |
| A refused command in a query | The whole call failed and earlier results were lost | Earlier results stand; the refused command's result carries the error; later commands are not sent | same-looking-different | Proof only |
| `Execute` and a refused command | Later commands in the script still ran | Stops at the failing command | same-looking-different | Proof only |
| One pipeline per script | 200 `PING`s in 2.4 ms | One request per command: 13.8–23.4 ms on localhost | same-looking-different | Proof only; long scripts pay N round trips |
| Per-result latency | Time since the call began | The command's own duration | same-looking-different | Proof only |
| Result-size cut | `proto.Size` at power-of-two counts | Byte-accurate cut; same message | same-looking-different | Proof only |

### Spanner

| Behaviour | Before | After | Class | Where it is pinned |
| --- | --- | --- | --- | --- |
| Extra parameters and TLS | Unknown names ignored; `use_ssl`/TLS material ignored | `emulator` is read; TLS material and `ServerName` are honoured | fix | Proof `TestSpannerConfig` |
| An instance with a stored Cloud credential pointed at the emulator | The Google client refused `WithoutAuthentication` plus credentials | The adapter drops the credential when `emulator` is set | fix | Proof `TestOpenReadsTheCredentialDocument` |
| `GRANT`/`REVOKE` | Failed as "not valid DML" | Routed as DDL through `UpdateDatabaseDdl` | fix | omni#1289 |
| `Close` | Leaked the data client | Closes the data client too | fix | Proof only |
| Statement errors | The whole call failed with `code = "Unknown"` | The error rides on its result; the code is `InvalidArgument` | same-looking-different | Proof only |
| Result-size cut | Power-of-two counts; 32 rows at a 300-byte cap | Byte-accurate; 25 rows, same message | same-looking-different | Proof only |
| PostgreSQL-dialect databases | Filtered from `SyncInstance`; `SyncDBSchema` failed with `text = boolean` | Listed, synced and dumped, but the product's split/query-span path is still GoogleSQL-only | needs a product decision | Proof only; keep the old filter until the product routes by database dialect |
| `CREATE DATABASE "name"` | Created a GoogleSQL database with a quoted name | Creates a PostgreSQL-dialect database | needs a product decision | Proof only; not reachable because `SPANNER` is absent from `EngineSupportCreateDatabase` |

### TiDB

| Behaviour | Before | After | Class | Where it is pinned |
| --- | --- | --- | --- | --- |
| Partition order | Name order (`p1, p10, p2`), so the dump could not be restored (`VALUES LESS THAN value must be strictly increasing`) | Ordinal order (`p10, p2, p1`); restore and re-dump are byte-identical | fix | Proof only; add the `p10,p2,p1` regression test |
| Unsigned integers | Strings | `uint64Value` | fix | omni#1354 |
| BINARY/VARBINARY/BLOB | Raw bytes in a string; an empty binary became NULL | `bytesValue`; an empty binary is empty bytes | fix | Proof only |
| Role grants | Interpolated account text | Properly quoted `'u'@'h'` | fix | Proof only |
| Explain requests | `QueryPlan` was nil | `QueryPlan{format}` is set; a written `EXPLAIN` is unchanged | same-looking-different | Proof only |
| Result-size cut | 32 → 20 rows at 4 KB; 128 → 77 at 16 KB | Same direction, same message | same-looking-different | Proof only |
| `CountAffectedRows` on a non-DML statement | Server-style error | Client-side refusal with a different message | same-looking-different | Proof only; plancheck asks only for DML |
| Invalid-SQL text | `Syntax error at line 1:51: … related text: …` | `expected identifier (line 1, column 51)` | needs a product decision | Proof only; wrap with `convertParseError` if the console copy depends on the old text |
| `MaximumSQLResultSize == 0` | One row plus a `0MB` error | No cap | needs a product decision | Contract reads 0 as no cap |
| `STATEMENT_DML_DRY_RUN` advising | Ran one `EXPLAIN` and produced one advice through `GetDB` | With a nil handle it returns early: zero queries, zero advice | needs a product decision | Proof `TestDMLDryRunNeedsASQLDB`; needs the advisor prober seam and an autocommit probe |
| `-- ghost` on a TiDB task | The cleanup called `GetDB` and would panic | Still would unless ghost is gated to MySQL/MariaDB or the SQLDB seam is added | needs a product decision | Proof call-site audit; not exercised |

### Trino

| Behaviour | Before | After | Class | Where it is pinned |
| --- | --- | --- | --- | --- |
| `information_schema` in sync | Listed as a database; the dump emitted a bogus `CREATE SCHEMA "information_schema"` | Gone; only the real database remains | fix | Proof only |
| Views in the diff SDL | Dropped, although the driver wrote them in a dump | Carried, with a test that the driver dump equals `schema.GetDatabaseDefinition` | fix | `backend/plugin/schema/trino/get_database_definition.go` (+14 lines in the proof) |
| Timestamps | Fixed accuracy of 6, and zoned timestamps in the plain timestamp field | The server's own precision (`0`, `3`, …), and `TIMESTAMP WITH TIME ZONE` in `timestamp_tz_value` | fix | `backend/plugin/db/trino` expectation tests |
| Dump format | Table/view definitions only | `CREATE SCHEMA IF NOT EXISTS` first, then tables, then topologically ordered views | same-looking-different | omni#1280; four fixtures in `backend/plugin/schema/trino/get_database_definition_test.go` gain the prefix |
| Exact numerics | `DECIMAL` arrived as a double | Arrives as text, because `driver.Value` has no exact-decimal kind | same-looking-different | `backend/plugin/db/trino` expectation test |
| Session across `QueryConn` calls | The caller's `*sql.Conn` kept session state | One session per call; a `SET SESSION` in one editor run does not reach the next | same-looking-different | Proof only; the three API sites already tolerate a nil handle |
| Task-run statement text | Comments stripped; `Range` was nil | Comments kept; the adapter supplies the real byte range | same-looking-different | Proof only |
| DDL through the query path | Synthesised `Affected Rows: 0` | Empty result | needs a product decision | Proof only; 5 adapter lines restore the old row |
| `SHOW`/`DESCRIBE` row cap | Appended `LIMIT` even where the engine refuses it | The cap is not enforced; `SHOW CATALOGS` with `Limit: 1` returns all rows | needs a product decision | Proof only; real fix is post-truncation in `sql_service` |

### MariaDB (proof still running)

MariaDB is served inside `db/mysql`, so the MySQL rows above apply. The rows below are the
differences the in-flight proof has already measured; this is not yet a complete delta list.

| Behaviour | Before | After | Class | Where it is pinned |
| --- | --- | --- | --- | --- |
| `CALL` and `RETURNING` result sets | Rows could be dropped when the statement's columns were not traced | Unmasked result sets are returned; `INSERT`/`DELETE … RETURNING` returns its rows | fix | omni `mariadb/driver` query tests and tapes |
| Script parsing | Per-statement fallback | Whole-script parse before sending; a parse error fails the call | same-looking-different | omni `mariadb/driver/query.go` |
| Result-size cut | A size error | `Truncated` with partial rows | same-looking-different | Proof only |
| Session state after `Close`/`Execute` | A pooled connection could carry state | The physical connection is discarded | same-looking-different | omni `mariadb/driver` tests |
| Unsigned integers | Strings | `uint64Value` | fix | omni#1354 |
| A masked untraced result set | Returned rows | Withheld with the count shape | needs a product decision | omni `mariadb/driver` masking tests |
| The MySQL advisor and gh-ost changes | N/A | The mysql rows apply, including the `DriverHandle` probe and gh-ost failing closed | needs a product decision | Same as MySQL |

## Product-side actions at the pin bump

These are not release-note lines. They are changes the Bytebase PR must carry, or the pin bump will
compile and behave wrongly.

1. **Fixtures that move with the pin.**
   - `backend/plugin/parser/doris/limit_test.go:197`: the `SELECT * FROM t # hash comment` expectation becomes the subquery wrap.
   - `backend/plugin/schema/trino/get_database_definition_test.go`: four dump strings gain `CREATE SCHEMA IF NOT EXISTS`; add views to the diff SDL.
   - `backend/plugin/aireview/catalog_test.go`: the MySQL secondary-index comment is inline in `CREATE TABLE`, not a separate `COMMENT ON INDEX`.
   - Result-shape expectations in the engine packages: trino `DECIMAL` as text, real timestamp accuracy, `timestamp_tz`; tidb/mysql unsigned and binary values; clickhouse/doris no synthetic `Affected Rows 0`.
   - Any stored `raw_dump`/golden for doris, clickhouse, hive, mssql, trino and tidb changes on the first sync; the new text should be accepted, not pinned back to the old renderer.
2. **`GetDB` retirement (`SQLDBProvider` / `db.SQLDB`).** The pg proof removed `GetDB` from `db.Driver`; mssql (8 files), oracle (12 files) and mysql (17 files) then reused the seam. The three `sql_service` sites, `rollout_service`, `statement_advise_executor`, `release_service_check`, the gh-ost validator and the two gh-ost cleanups must each name their fallback. Do not leave an adapter returning a fake handle.
3. **gh-ost.** MySQL's five binlog probes have no contract counterpart, so validation now fails closed. Either add the probe to `omni/mysql/driver` or record that gh-ost validation is off; do not let it disappear silently. TiDB's `-- ghost` path also needs an engine gate because its cleanup calls `GetDB`.
4. **`ExecuteOptions`.** The mysql proof added `Transaction` and `Isolation`, but the adapters still parse `-- txn-mode` themselves, so no task-run caller was forced to change. Keep that fallback; do not remove the directive parse until every engine's adapter reads the options.
5. **Advisor prober.** The pg proof's `advisor.Prober`/`Context.Prober` is the generic seam; MySQL and MariaDB need the engine-specific `DriverHandle` instead, because their rules call `ExplainJSON`/`EstimateAffectedRows`, not a session query. TiDB's dry-run rule and Oracle's live rule need the prober; mssql's static rules do not. Add the autocommit probe for TiDB's `BATCH … DRY RUN` before calling the seam complete.
6. **plancheck.** The concrete `CountAffectedRows` assertions survive the adapters, so the interface change is optional for every engine measured. It carries nine blank parser imports that the concrete driver imports used to provide; if it lands, run `go test ./backend/runner/plancheck/`, not only a build, and never tidy those imports.
7. **TLS default.** DynamoDB's old driver always verified and ignored `verify_tls_certificate`; mapping the data source's TLS fields naively would make an unset checkbox mean "skip verification". The adapter must verify with the system roots (or `ssl_ca`) until the product can distinguish an explicit opt-out from an unset default. The same rule holds for the other cloud drivers: a nil `Config.TLS` means the client library verifies.
8. **Oracle's private-link guard.** Take `*sql.Conn` off the interface before the Oracle pin. Otherwise `db.SQLDB(adapter)` is nil and every linked read is refused with no message naming the cause.
9. **Spanner's dialect listing.** Do not ship the PG-dialect listing alone. Either keep the old `SyncInstance` filter until the product routes `SplitMultiSQL`/`GetQuerySpan` by database dialect, or ship both halves together; masking and access checks are keyed off the spans.
10. **Engine test pins and coverage.** Redis moves to `redis:7.4.12-alpine`; TiDB moves to `v8.5.5`; mssql moves to `2022-CU27-GDR1-ubuntu-22.04`. Cassandra, Hive and Elasticsearch have no server-backed test at Bytebase's seam; add one or say omni's tape suite is the only coverage.

## Where the records are incomplete

- **StarRocks, CockroachDB and OceanBase** have proofs still running. Their deltas are not in this
  page; do not infer them from Doris, MySQL or PostgreSQL. StarRocks is the one fork whose server-DDL
  dump is known to differ from Doris's (no procedures, its own object order), so it needs its own
  table before its pin moves.
- **CosmosDB, MongoDB, Redshift, Snowflake, BigQuery and Databricks** have no completed cutover report
  in this set. The last four also have unverified drivers.
- **Unclassified from the records:** Oracle's `Session.Query` has no transaction control, so a
  future advisor that probes with DML could not roll back; MySQL's RDS IAM mapping was not tested
  against a real endpoint; Spanner's real-project authentication was not tested; Redis's per-command
  round trip is accepted without a long-script benchmark; the Oracle catalogue's quoted key and the
  Doris procedure's qualified name are recorded but no consumer has yet forced a normalization
  decision.
