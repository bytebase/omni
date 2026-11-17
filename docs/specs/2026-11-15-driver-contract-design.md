# Driver contract

2026-11-15 · Proposed

Bytebase is moving its target-database drivers (`backend/plugin/db/<engine>`) into omni as
`<engine>/driver` packages. This note fixes the contract they implement, `driver/`, and maps every
part of Bytebase's `db.Driver` today to its home after the move. Bytebase's plan is
`docs/design/engine-code-to-omni.md` in the Bytebase repository, Phase 4.

## The cut

Omni speaks the protocol; Bytebase supplies the path. Bytebase resolves secrets, cloud IAM tokens,
SSH tunnels and TLS material, and passes the result as values and funcs in `driver.Config`. The
engine package does everything after the dial: session setup, sync queries, statement execution,
result decoding. Bytebase keeps the lifecycle around it: which driver to open (engine dispatch),
pooling and caching opened drivers, sync scheduling, task execution, and conversion of omni's
types into its API and store protos, the way `omni/metadata` is converted today.

## Package rules

- `driver/` holds types and interfaces only, and imports nothing but the standard library and
  `omni/metadata`.
- `<engine>/driver` is the one kind of non-test omni package that links an engine client (pgx,
  go-sql-driver, ...). No other omni package imports a driver package, so a caller of a parser or
  catalog does not link a client.
- An engine package exposes `Open(ctx, driver.Config)`; an engine with settings of its own (pg's
  tenant mode, Redshift's datashare) takes them as a second, engine-specific argument.
- A driver's tests replay tapes (`<engine>/internal/testcontainer`) recorded on one pinned server
  per supported major version, listed in `<engine>/testdata/pins.tsv`.

## Where each part of `db.Driver` goes

| Bytebase today | after the move |
| --- | --- |
| `Open(ctx, storepb.Engine, ConnectionConfig)` | engine dispatch and `DataSource` → `driver.Config` stay in Bytebase; the open is `<engine>/driver.Open` |
| `Close`, `Ping` | `driver.Driver`; Bytebase also closes the tunnel or cloud dialer it built |
| `GetDB() *sql.DB` | removed. Its seven uses outside `plugin/db` (SQL editor sessions, rollout session listing, advisor live facts, gh-ost cleanup) each become a driver method in the query slice |
| `Execute(ctx, stmt, ExecuteOptions)` | `driver.Executor` |
| `QueryConn(ctx, *sql.Conn, stmt, QueryContext)` | `driver.Querier` → `Session.Query`, returning `[]*driver.Result` |
| `SyncInstance` → `*db.InstanceMetadata` | `driver.Syncer.SyncInstance` → `*driver.Instance` |
| `SyncDBSchema` | `driver.Syncer.SyncDatabase`; `driver.ErrNotFound` for a missing database |
| `Dump(ctx, w, metadata)` | `driver.Dumper`; dropping Bytebase's backup schema and repairing legacy stored metadata stay in Bytebase, before the call |

| `ConnectionConfig` / `ConnectionContext` field | after the move |
| --- | --- |
| `DataSource` | converted by Bytebase into `Host`, `Port`, `User`, `TLS`, `Dial`, `Params` |
| `Password` (resolved secret or IAM token) | `Password` |
| `DatabaseName` | `Database` |
| `EngineVersion` | `EngineVersion`, a fallback only (see Decisions) |
| `ReadOnly` | `ReadOnly` |
| `TaskRunUID` | `ApplicationName`, formatted by Bytebase |
| `TenantMode`, `DataShare` | engine-specific `Open` settings (pg, Redshift) |
| `MessageBuffer` | `Result.Messages` |
| `EnvironmentID`, `InstanceID` | stay: they serve IAM and logging, which stay |

| `QueryContext` field | after the move |
| --- | --- |
| `Schema`, `Limit`, `Explain`, `Option.ExplainFormat` | `QueryOptions` |
| `MaximumSQLResultSize` | `QueryOptions.MaxResultBytes` (see Decisions) |
| `Timeout` | the caller's context deadline |
| `OperatorEmail`, `SkipMasking`, `DataSourceType` | stay: product decisions |
| `MaskingEnabled`, `Container` | open for the query slice: the first asks a driver to withhold untraced rows, the second is Cosmos DB's |

| `ExecuteOptions` | after the move |
| --- | --- |
| `CreateTaskRunLog` and the `Log*` methods drivers call | `driver.ExecuteLogger`; Bytebase's implementation writes `storepb.TaskRunLog` |
| `LogCommandStatement` | stays in Bytebase's logger: `Command` carries both the text and its range |
| `MaximumRetries` | `MaxRetries` |
| `Log{ComputeDiff,DatabaseSync,SchemaDump,GhostMigration}*` | stay: task runners emit them, not drivers |
| `CreateDatabase` | open for the query slice |

Methods Bytebase reaches on concrete drivers (pg's `GetCurrentDatabaseOwner`,
`CountAffectedRows`, `GetSearchPath`) move to the engine's driver package as its own methods.

## Decisions

- **TLS is a func of the host.** `TLS func(host string) (*tls.Config, error)` rather than one
  `*tls.Config`: pg dials one of several hosts, and Bytebase binds certificate verification to the
  host it verifies. Nil keeps the engine's default negotiation, so pg's default stays libpq's
  `sslmode=prefer` as today.
- **The version is detected; the passed one is a fallback.** Bytebase's pg sync already prefers
  the live `server_version_num` because the cached one can be empty or stale.
- **A federated token is a func, not the client's connector.**
  `TokenProvider func(ctx context.Context) (string, error)` sits beside `TLS` and `Dial`: the
  caller resolves the token, and the engine hands it to its client's native token path. For
  mssql, Bytebase passes an Entra ID access token for `https://database.windows.net/.default`,
  from its Azure client-secret credential or default credential chain, and the driver opens
  go-mssqldb's security-token connector. `User` may be empty for a token login and `Password` is
  unused. The engine's own connector or token provider cannot cross the contract: it would put
  the client type in `driver/`, and mssql's DSN-only Azure AD connector cannot also carry
  `Config.TLS` and `Config.Dial`.
- **Results are Go structs, not a proto.** They cross no process or storage boundary, so
  `driver.Result` needs no wire tolerance; Bytebase converts it field by field into
  `v1pb.QueryResult`. `Value` keeps every `RowValue` kind distinct (FLOAT vs DOUBLE, INT32 vs
  INT64) so the conversion loses nothing.
- **Engine errors stay engine types.** `Result.Err` wraps the client's error (pg:
  `*pgconn.PgError`); the Bytebase pg adapter reads its fields for `QueryResult.PostgresError`. A
  neutral error struct would have to grow every engine's fields.
- **The caller picks the transaction mode.** Bytebase's `-- txn-mode` directive and its default
  are product conventions; Bytebase parses them and passes `ExecuteOptions.Transaction`.
- **The result size cap is measured by omni.** Bytebase measures `proto.Size` of the API
  message today; a driver cannot. The query slice must show the cut-off moves by little.

## Order

1. This contract.
2. pg sync and dump (`pg/driver`), proved equal to Bytebase's `db/pg` on PostgreSQL 12, Bytebase's
   documented floor, through 18.
3. pg query and execute, which settles `Querier`, `Session`, `Executor` and the open rows above.
   Until then those types are a draft.
4. Each further engine: sync and dump first, then query and execute.
