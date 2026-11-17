# Driver method parity beyond the five interfaces

2026-11-17 · Audit

Bytebase calls more on a `db.Driver` than the five contract interfaces (`Driver`, `Syncer`,
`Dumper`, `Querier`, `Executor`). This records, engine by engine, where those extra methods live
on omni's side, so a cutover does not fail on a method instead of on a behavior. The audit is
against Bytebase `345007290e` and omni `16a7f5c49`.

Result: every extra method Bytebase calls on pg, mysql/mariadb, oceanbase, tidb, mssql and Oracle
has an omni counterpart. Oracle's version accessor is exported, and the four MySQL-family drivers
agree on which statements an affected-row estimate may describe. The only missing cells belong to
cockroachdb and redshift, whose drivers are in flight or blocked on an endpoint.

| Bytebase method | Where Bytebase uses it | Omni counterpart |
| --- | --- | --- |
| `CountAffectedRows` | plancheck's SQL summary report, through a concrete driver type | `<engine>/driver.CountAffectedRows` |
| `Explain` / `ExplainJSON` / `ExplainTabularEstimate` / `EstimateAffectedRows` | MySQL-family advisors; pg advisors' live EXPLAIN; OceanBase and TiDB advisors | same-named driver methods |
| `GetCurrentDatabaseOwner` / `GetSearchPath` / current user | pg tenant mode and sync; advisor object-owner check; Redshift owner | `CurrentDatabaseOwner`, `SearchPath`, `CurrentUser` |
| `OwnsPrivateDatabaseLink` | Oracle linked-read refusal in `sql_service.go` | `driver.PrivateDatabaseLinkOwner` |
| `GetVersion` | Oracle `SyncDBSchema`'s release-dependent column queries | `oracle/driver.Version` |
| `GetDB` | SQL editor, export, rollout session listing, advisor live facts, gh-ost cleanup and validation | `Session`/`Query`, `Activity`, the probe methods above, and `Execute` |
| `StopConnectionByID` | internal to the MySQL-family `Execute`/`QueryConn` paths only | unexported kill-on-cancel inside `Session`/`Execute` |
| `PushAndClearMessages` | internal to pg `QueryConn` | `Result.Messages` |
| role listing | `SyncInstance`'s roles, never a separate driver call | `Syncer.SyncInstance`'s `Roles` |

## Matrix

`same` means the method exists with Bytebase's behavior. `extra` means omni exposes something
Bytebase has no method for. `n/a` means Bytebase does not call a driver method for it, with the
reason in the next section. Empty means neither side has the method.

| Engine | `CountAffectedRows` | explain probes | owner / search path / user | `OwnsPrivateDatabaseLink` | `GetVersion` | `StopConnectionByID` | `GetDB` |
| --- | --- | --- | --- | --- | --- | --- | --- |
| pg | same | extra | same | | | n/a | n/a |
| mysql | same | extra | | | | n/a | n/a |
| mariadb | extra (Bytebase uses the mysql driver) | extra | | | | n/a | n/a |
| oceanbase | same | extra | | | | n/a | n/a |
| tidb | same | extra | | | | n/a | n/a |
| mssql | same | | | | | | n/a |
| oracle | same | | | same | same | | n/a |
| cockroachdb | **omni missing** | | **omni missing** | | | | n/a |
| redshift | **omni missing** | | **omni missing** | | | | n/a |
| doris | | | | | | n/a | n/a |
| starrocks | | | | | | n/a | n/a |
| clickhouse | | | | | | | n/a |
| hive | | | | | | | n/a |
| snowflake, bigquery, databricks, spanner, mongo, redis, cassandra, cosmosdb, dynamodb, elasticsearch, trino | | | | | | | n/a |

## Cells left for another session or blocked

| Engine | Method | Why it is left |
| --- | --- | --- |
| cockroachdb | `CountAffectedRows`, `GetCurrentDatabaseOwner` | The driver is in flight in another session. |
| redshift | `CountAffectedRows`, `GetCurrentDatabaseOwner` | There is no `redshift/driver` and no Redshift endpoint on this machine; the driver is blocked on a real cluster. |
| starrocks, clickhouse, hive | (no extra method) | Their `Dumper` is in flight; their extra-method surface is already covered. |
| doris | (no extra method) | `Dumper` is in flight; `StopConnectionByID` is internal and omni kills on cancel. |

## `GetDB` replacements

`GetDB` is removed by the contract, so every use has a replacement; none is a missing driver
method.

| Use | Replacement |
| --- | --- |
| SQL editor, export, interactive queries | `Session` + `Session.Query` |
| rollout session listing | pg `Activity`; the endpoint only supports PostgreSQL and CockroachDB |
| advisor live facts (pg, mysql, mariadb, oceanbase, tidb) | `Explain`, `ExplainJSON`, `EstimateAffectedRows`, `ColumnExtra`, `DryRunAutoCommit`, `CurrentUser` |
| gh-ost cleanup | `Execute` |
| gh-ost binlog validation | `Session.Query`; the four reads are product SQL, not engine judgment |
| Oracle linked-read and DML dry-run probes | the private-link probe and `Session.Query`; Oracle is another session's engine |

## Harmless `extra` and divergences

- **mariadb's `CountAffectedRows`.** Bytebase has no `db/mariadb`; its MySQL driver branches on
  `Engine_MARIADB` and calls the MariaDB parser. Omni's separate `mariadb/driver` implements the
  same behavior, including the tabular-EXPLAIN fallback, `RETURNING` syntax and the same refusal of
  a statement it cannot estimate, so the cell is harmless.
- **MySQL-family explain probes.** Bytebase exposes `ExplainJSON`, `ExplainTabularEstimate` and
  `EstimateAffectedRows` as package functions over `*sql.DB`; omni makes them driver methods.
  Same behavior. Omni also adds `Explain`, which Bytebase uses only through raw SQL.
- **The MySQL family refuses what `CountAffectedRows` cannot describe.** MySQL, MariaDB,
  OceanBase and TiDB each refuse anything but exactly one INSERT, UPDATE or DELETE before they
  send an EXPLAIN; TiDB also still refuses a plan whose root is not that DML. Bytebase's old driver
  rejected a non-DML statement only in MySQL, but plancheck passes DML only, so the stricter
  MariaDB, OceanBase and TiDB paths change no production result.
- **Oracle's `Version`.** Bytebase's `GetVersion` parsed the engine version it had already recorded.
  Omni's `Version` asks the server and falls back to `Config.EngineVersion`, so a cutover reads a
  live release when it can and the same recorded one when it cannot.
- **pg's `Explain`, `ExplainJSON`, `CurrentUser` and `Activity`.** These replace Bytebase's raw
  `*sql.DB` advisor and session-listing queries; `PushAndClearMessages` is internal and becomes
  `Result.Messages`.
- **oceanbase's `ColumnExtra` and tidb's `DryRunAutoCommit`.** Extra driver surface that matches
  Bytebase's raw SQL; neither is called by the cutover as a missing method.

## Refresh

Re-run the comparison from Bytebase's `backend/plugin/db/<engine>` and `backend/plugin/db/driver.go`
against omni's `*/driver`. Bytebase engines with an extra method today are pg, mysql (covering
mariadb), oceanbase, tidb, mssql, oracle, cockroachdb, redshift, doris and starrocks; the rest
only reach `GetDB`.
