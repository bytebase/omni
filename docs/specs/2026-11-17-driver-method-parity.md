# Driver method parity beyond the five interfaces

2026-11-17 · Audit

Bytebase calls more on a `db.Driver` than the five contract interfaces (`Driver`, `Syncer`,
`Dumper`, `Querier`, `Executor`). This records, engine by engine, where those extra methods live
on omni's side, so a cutover does not fail on a method instead of on a behavior. The audit is
against Bytebase `345007290e` and omni `16a7f5c49`; the [dialer and probe
support](#dialer-and-probe-support) section is re-read at omni `c7b6c5a50`.

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
| starrocks, clickhouse, hive, doris | (no extra method) | Their `Dumper` has landed since this audit, clickhouse's in #1350, so no cell is left there; their extra-method surface was already covered. |

## Dialer and probe support

A cutover hands every engine the same `driver.Config`, so two things have to hold per engine: the
caller's dialer reaches the server, and every capability a driver answers outside the five
interfaces is reachable through a declared interface. `Config.Dial` replaces the socket the
engine's client opens, which is what an SSH tunnel or a cloud connector needs. A probe interface
lives in `driver/probe.go` and is type-asserted by name: on the driver where the answer needs no
session, on the session `Querier.Session` returns where it does.

| Engine | `Config.Dial` | Probe interfaces |
| --- | --- | --- |
| pg | yes — pgx `DialFunc` | — |
| mysql, mariadb, tidb, oceanbase | yes — go-sql-driver `DialFunc` | — |
| doris, starrocks | yes — go-sql-driver `DialFunc` (their mysql fork's) | — |
| mssql | yes — a protocol dialer registered with go-mssqldb | — |
| oracle | yes — go-ora `DialerContext` | `PrivateDatabaseLinkOwner` |
| trino | yes — `http.Transport.DialContext` | — |
| mongo | yes — `SetDialer` | — |
| spanner | yes — gRPC `WithContextDialer` | — |
| cassandra | yes — gocql `HostDialer`, one per host the cluster publishes | — |
| elasticsearch | yes — `http.Transport.DialContext` | — |
| cosmosdb, dynamodb | yes — `http.Transport.DialContext` | — |
| redis | yes — go-redis `Dialer`, on all three topologies | `AllNodesQuerier` on the session |
| clickhouse | yes — clickhouse-go `DialContext` | — |
| hive | **no** — gohive v2's connector has no field for one | — |

Two engines accept an empty `Config.Host` when `Dial` names the endpoint itself, as a cloud
connector does: oracle, and redis with no topology set.

**hive is the one engine a tunnel cannot reach through this contract.** gohive v2 keeps its dialer
in an unexported connect configuration that `HiveConnector` fills from a fixed list of exported
fields, and `gohive.Config` has no field for a dialer, so `hive/driver` opens the socket itself.
`driver.Config.Dial` says an engine in this position must say so in its package doc and name what a
caller does instead; `hive/driver`'s doc does, and the answer is a local port forward whose address
the caller names as `Config.Host` and `Config.Port`, with the server's own host in the
`tls.Config` its `Config.TLS` returns.

**redis's all-nodes query is the one probe a session answers.** Bytebase offers a per-request "run
on every node" choice for Redis; `Session.Query` routes a command the way the server does, so only
`*redis/driver.Session.QueryAllNodes` answers it, and `driver.AllNodesQuerier` is the interface a
caller asks the session for. Dropping the method was the alternative and it does not work: a
cluster's nodes each answer the command, and the caller shows them as one result.

Nothing else is outside the contract: no other driver has a method a caller would have to hold the
concrete type for.

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
only reach `GetDB`. Re-read [Dialer and probe support](#dialer-and-probe-support) from each
`<engine>/driver`'s `Open` (where `Config.Dial` reaches the client) and from the interfaces in
`driver/probe.go` that a driver or session declares.
