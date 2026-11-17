# Verifying databricks/driver

`databricks/driver` was verified against a scratch Databricks workspace on
2026-11-17. This file records how: the account and warehouse pinned, the tests
recorded in order, the tape harness at the SDK's HTTP seam, and the parity
harness that holds this package to Bytebase's driver.

The recordings and the parity run are real. The live tests passed on the
workspace, the six tapes replay with no credentials, two consecutive recordings
wrote byte-identical tapes, and the parity harness agreed with Bytebase's driver
on the sync, the instance, queries and Execute. One difference is deliberate:
the workspace states a materialized view's own `SHOW CREATE TABLE`, and
Bytebase's rendered `DESC FORMATTED` shape is not a statement the workspace
accepts. Nothing in this record runs in CI. Databricks is an account service:
omni cannot start a workspace, SQL warehouse or Unity Catalog in a container,
and no PostgreSQL or Hive stand-in answers the Unity Catalog or Statement
Execution APIs. `databricks/parser` does not exist either, so there is no
grammar corpus to judge.

## What is here, and what is missing

| capability | state |
| --- | --- |
| `SyncInstance` (SQL version, catalog names) | verified on the pinned workspace |
| `SyncDatabase` (schemas, tables, views, materialized views, external tables) | verified on the pinned workspace |
| `Session.Query` (types, NULL, limit, explain, `MaxResultBytes`) | verified on the pinned workspace |
| `Execute` (statement by statement, no transaction control) | verified on the pinned workspace |
| `Dump` (`SHOW CREATE TABLE` for every object) | verified on the pinned workspace |
| tapes (`testdata/tapes`) | six, recorded from the pinned workspace; replay without credentials |
| sync goldens (`testdata/sync`) | `workspace.json`, compared with `proto.Equal` |
| `driver` row in `databricks/testdata/pins.tsv` | workspace host, SQL release, warehouse ID and catalog |
| live tests (`live_test.go`) | run against the workspace; skip without `DATABRICKS_DSN` |
| recorded tests (`tape_test.go`) | `TestSyncInstance`, `TestSyncDatabase`, `TestQuery`, `TestExecute`, `TestDump`, `TestPing` |
| `databricks/catalog` | **does not exist**; there is no metadata → DDL renderer, so `Dump` states the server's own DDL |

The unit tests that do run cover the SDK configuration, the splitter, statement
classification, the value conversions, the sync conversions and order, and the
dump's pure helpers.

## Step 1 — the account and warehouse

Pin one dedicated, scratch-only workspace. The tests create a schema and a
table, `Execute` writes, and `Dump` reads every object in the workspace, so it
must be a workspace whose contents nobody else depends on.

The account must provide:

- a workspace URL, such as `https://dbc-00000000-0000.cloud.databricks.com`;
- a personal access token for a user or service principal that may read Unity
  Catalog (`USE CATALOG`, `USE SCHEMA`, `SELECT`), list catalogs, schemas and
  tables, and create and drop a scratch schema and managed table in the
  configured catalog;
- one SQL warehouse the token may run statements on, and its ID;
- one catalog to sync and dump. A catalog with a managed table, a view and
  (where the workspace supports it) a materialized view is enough. An external
  table needs a cloud storage location the token can read and is optional.

If the token is a service principal's, do not use a production workspace: the
dump prints every catalog it can list, and the tests create objects.

The DSN is one variable because every live test needs all four values. The token
is in the query string, so keep it out of shell history:

```bash
export DATABRICKS_DSN='databricks://dbc-00000000-0000.cloud.databricks.com?token=<personal-access-token>&warehouse_id=<warehouse-id>&catalog=<catalog>'
```

`databricks/driver.live_test.go` parses `databricks://` and `https://` hosts,
`warehouse_id` (or `warehouseId`) and `catalog`. `Open` reads the token from
`driver.Config.Password`, the warehouse from `Params["warehouseId"]` and the
catalog from `Config.Database`; nothing reads the environment in production.

Record what the workspace answers, because that is the pin:

```bash
# Ask the SQL Editor, a SQL client or the SDK; the driver records the same
# answer through TestLiveSyncInstance.
SELECT VERSION();
```

Write the two words the server answers into the pin row in Step 3. The SQL
release is the first word, which is the version `SyncInstance` keeps. The
workspace answered `4.2.0 0000000000000000000000000000000000000000`.

## Step 2 — the fixtures

Load one corpus into the configured catalog before recording, because the sync
reads whatever is there and a tape of an empty catalog pins nothing:

```sql
CREATE SCHEMA IF NOT EXISTS `<catalog>`.`omni_record`;

CREATE TABLE `<catalog>`.`omni_record`.`lookup` (
  id INT,
  label STRING NOT NULL,
  amount DECIMAL(12, 2),
  flag BOOLEAN,
  seen TIMESTAMP,
  payload BINARY
) COMMENT 'a synced comment';

CREATE TABLE `<catalog>`.`omni_record`.`facts` (
  id BIGINT,
  lookup_id INT,
  notes STRING
) USING DELTA;

COMMENT ON COLUMN `<catalog>`.`omni_record`.`lookup`.`label` IS 'a column comment';

CREATE VIEW `<catalog>`.`omni_record`.`recent` AS
  SELECT id, lookup_id FROM `<catalog>`.`omni_record`.`facts` WHERE id > 0;

-- Optional: requires a workspace that supports materialized views.
CREATE MATERIALIZED VIEW `<catalog>`.`omni_record`.`totals` AS
  SELECT lookup_id, SUM(id) AS total
  FROM `<catalog>`.`omni_record`.`facts`
  GROUP BY lookup_id;
```

The sync keeps a managed table, a view, a materialized view and an external
table, and skips a streaming table. The corpus above reaches the first three. An
external table needs a `LOCATION` under cloud storage; record it only when the
workspace already has a location the token may read, and say so in the PR that
records.

The Dump test reads every catalog the token can list. If the workspace holds a
catalog the token must not see, use a dedicated workspace rather than teaching
the dump to narrow itself; the dump's contract is the server's DDL, not the
metadata argument.

## Step 3 — the pin row

`databricks/testdata/pins.tsv` holds:

```
# columns: role, image, digest, source tag
driver	dbc-3c689f3f-f754.cloud.databricks.com	sha256:21c3596967a931ba19b528877357d45cb7c2ea34abc68315d0c4ca7907d9f408	4.2.0 warehouse_id=45847ad2fbd0b316 catalog=workspace
```

There is no image to pin, so `image` holds the workspace host, `digest` the
SHA-256 of the whole `SELECT VERSION()` answer, and the source tag names the
SQL release, warehouse and catalog. The token is never part of a pin.

## Step 4 — the tape harness

The driver talks HTTP through the Databricks SDK, so the tape sits at the SDK's
HTTP transport, not at a database/sql connection.
`databricks/internal/testcontainer` holds it:

- a `RoundTripper` records one request (method, route and body) and one answer
  (status, content type, body). `tape.OpenRaw` is the shape other SDK-backed
  engines use; the encoding is HTTP request/answer, not database/sql.
- `Open(t)` returns a `*databricks.WorkspaceClient` whose
  `Config.HTTPTransport` is the tape. The driver tests construct
  `&Driver{client: testcontainer.Open(t), warehouseID: testcontainer.WarehouseID,
  catalog: testcontainer.Catalog}` directly, so a replay needs no environment.
- a recording reads `DATABRICKS_HOST` and `DATABRICKS_TOKEN`, refuses a host
  other than the pin's, and never records the Authorization header.
- `main_test.go` runs `tape.Main`.
- `container_test.go` covers the encoding, the narrowing, the statement-ID
  stand-in and the pin constants.

The request/answer tape must narrow what the service makes up:

- `statement_id` is random and names a later poll; replace it with a numbered
  stand-in and keep the poll's path consistent with the stand-in.
- `created_at`, timings and `warehouse` state move between runs; keep only what
  the driver reads, or replace them with stable values.
- `result.data_array` chunk indexes, `next_chunk_index`,
  `next_chunk_internal_link` and `total_*_count` move with the chunking; keep
  the first chunk and record what the driver actually reads.
- `DESC FORMATTED` is no longer part of the dump: the workspace states a
  materialized view's own `SHOW CREATE TABLE`, so the tape carries that answer.
- A plan's expression IDs (`n#20318`) move between runs; the harness replaces
  them with `#0`.
- Listing order is not stable. The driver sorts catalogs, schemas, tables,
  columns and dependencies; a replay must not depend on the service's order.

## Step 5 — record, in this order

Sync first: it is the read-only half and the half Bytebase never tested. Then
the reads, then the writes, then the dump.

| order | test | what it records |
| --- | --- | --- |
| 1 | `TestSyncDatabase` | `Catalogs.ListAll`, `Schemas.ListAll`, `Tables.ListAll`, and the metadata they build: schemas, tables, columns, views, materialized views, external tables, comments |
| 2 | `TestSyncInstance` | `SELECT VERSION()` and the catalog names |
| 3 | `TestQuery` | statement classification, every value kind, NULL versus empty STRING, the row limit, the EXPLAIN text, `MaxResultBytes`, a zero-row query, a failing statement |
| 4 | `TestExecute` | `CREATE SCHEMA`, `CREATE TABLE`, `INSERT`, each committing itself, and the drop in cleanup |
| 5 | `TestDump` | `SHOW CREATE TABLE` for managed tables, external tables, views and a materialized view, and the dump text |
| 6 | `TestPing` | `Workspace.ListAll` as the connectivity check |

The recorded tests are `TestSyncDatabase`, `TestSyncInstance`, `TestQuery`,
`TestExecute`, `TestDump` and `TestPing`; the live tests `TestLive*` remain for
a direct run against `DATABRICKS_DSN`. The recordings were made one test at a
time so a failure was attributable, then the compiled test binary ran twice
from the package directory and the tapes hashed identically.

```bash
export DATABRICKS_DSN='databricks://...'
OMNI_RECORD=1 DATABRICKS_HOST="$DATABRICKS_HOST" DATABRICKS_TOKEN="$DATABRICKS_TOKEN" go test ./databricks/driver/ -run '^TestSyncDatabase$' -count=1 -v
OMNI_RECORD=1 DATABRICKS_HOST="$DATABRICKS_HOST" DATABRICKS_TOKEN="$DATABRICKS_TOKEN" go test ./databricks/driver/ -run '^TestSyncInstance$' -count=1 -v
OMNI_RECORD=1 DATABRICKS_HOST="$DATABRICKS_HOST" DATABRICKS_TOKEN="$DATABRICKS_TOKEN" go test ./databricks/driver/ -run '^TestQuery$' -count=1 -v
OMNI_RECORD=1 DATABRICKS_HOST="$DATABRICKS_HOST" DATABRICKS_TOKEN="$DATABRICKS_TOKEN" go test ./databricks/driver/ -run '^TestExecute$' -count=1 -v
OMNI_RECORD=1 DATABRICKS_HOST="$DATABRICKS_HOST" DATABRICKS_TOKEN="$DATABRICKS_TOKEN" go test ./databricks/driver/ -run '^TestDump$' -count=1 -v
OMNI_RECORD=1 DATABRICKS_HOST="$DATABRICKS_HOST" DATABRICKS_TOKEN="$DATABRICKS_TOKEN" go test ./databricks/driver/ -run '^TestPing$' -count=1 -v
```

Then record everything twice and require byte-identical tapes. Two runs of one
compiled test binary are deterministic; two `go test` runs are byte-identical in
the tapes but not in a protojson golden, so compare the sync goldens as protos:

```bash
go test -c -o /tmp/databricks-driver.test ./databricks/driver/
(cd databricks/driver && OMNI_RECORD=1 DATABRICKS_HOST="$DATABRICKS_HOST" DATABRICKS_TOKEN="$DATABRICKS_TOKEN" /tmp/databricks-driver.test -test.count=1)
sha256sum databricks/driver/testdata/tapes/*.tape > /tmp/databricks-tapes.a
(cd databricks/driver && OMNI_RECORD=1 DATABRICKS_HOST="$DATABRICKS_HOST" DATABRICKS_TOKEN="$DATABRICKS_TOKEN" /tmp/databricks-driver.test -test.count=1)
sha256sum databricks/driver/testdata/tapes/*.tape | diff /tmp/databricks-tapes.a -
```

A tape that is not reproducible is a nondeterminism in what the test sends or
reads. The usual causes: a listing without a sort (the service's order becomes
the tape's order), a map iterated into a statement, a made-up statement ID or
timestamp, and a result chunk boundary. Fix the driver, never the tape.

## Step 6 — the sync goldens

After the tapes, commit `databricks/driver/testdata/sync/<catalog>.json`: the
protojson of `SyncDatabase`'s output, compared with `proto.Equal` in the test
rather than byte-for-byte, because protojson spacing moves between builds. One
golden per catalog synced by the fixture. `SyncInstance`'s database list is
catalog names only, so it needs no separate golden beyond the names the test
asserts.

## Step 7 — the dump

`Dump` is the server's own DDL, not a rendering of the synced metadata, because
there is no `databricks/catalog` and no metadata → Databricks DDL renderer. The
recording settles the shape:

- `SHOW CREATE TABLE` for a managed table, an external table, a view and a
  materialized view.
- The workspace states a materialized view's own statement, including
  `AS SELECT ...`. Bytebase's `DESC FORMATTED` rendering has no query and its
  `TBLPROPERTIES` list carries internal UUIDs, so it is not restorable. This
  port now uses the server's `SHOW CREATE TABLE`, and the parity harness records
  the difference.
- The dump writes `CREATE CATALOG` for every catalog and `CREATE SCHEMA` for
  every schema but `default`, and skips the `system` catalog and each catalog's
  `information_schema`. The server accepted the quoted identifiers and the
  sorted object order.
- A `TestDumpRestores` is still future work: the dump text is known to apply to
  the objects it states, but the corpus has not been dropped and replayed.

## Step 8 — the parity harness

The gate that a tape cannot be: Bytebase's driver and this one, on the same
workspace, compared cell by cell. It lives in a throwaway Bytebase worktree and
is never committed.

Start a worktree of the omni branch this file landed on, and a scratch Bytebase
worktree beside it (never the shared checkout):

```bash
git -C <omni repo root> worktree add ~/omni-worktrees/databricks-verify -b databricks-verify origin/databricks-prep
git -C <bytebase repo root> worktree add ~/bb-worktrees/databricks-parity -b databricks-parity origin/main
cd ~/bb-worktrees/databricks-parity
go mod edit -replace github.com/bytebase/omni=$HOME/omni-worktrees/databricks-verify
go mod tidy
```

Records happen in `~/omni-worktrees/databricks-verify`; the parity harness is
`backend/plugin/db/databricks/parity_omni_test.go` in
`~/bb-worktrees/databricks-parity`, and neither worktree is committed. Remove
both when the recording lands.

Add `backend/plugin/db/databricks/parity_omni_test.go` (package `databricks_test`,
skipping without `DATABRICKS_DSN`) and compare:

| Bytebase | omni | comparison |
| --- | --- | --- |
| `(&Driver{}).Open(ctx, storepb.Engine_DATABRICKS, db.ConnectionConfig{DataSource: &storepb.DataSource{Host: ..., AuthenticationPrivateKey: token, WarehouseId: ...}, ConnectionContext: db.ConnectionContext{DatabaseName: catalog}})` | `databricksdriver.Open(ctx, driver.Config{Host: ..., Password: token, Database: catalog, Params: map[string]string{"warehouseId": warehouseID}})` | both open, neither reaches the service |
| `SyncDBSchema(ctx)` | `SyncDatabase(ctx)` | same `*metadatapb.DatabaseSchemaMetadata`, `proto.Equal` |
| `SyncInstance(ctx)` | `SyncInstance(ctx)` | version text; catalog names and order. Bytebase's `Databases` carry schemas and tables, omni's carry names only: compare names, and record the difference rather than "fixing" it |
| `QueryConn(ctx, conn, sql, queryContext)` | `Session.Query(ctx, sql, queryOptions)` | one result per statement; column names and `TypeName`s; each cell's kind and value; `AffectedRows`; the statement text; the error's text. Bytebase drops a zero-row query's result, omni keeps its columns |
| `Execute(ctx, script, executeOptions)` | `Execute(ctx, script, executeOptions)` | total affected rows (both 0), and Bytebase's `TaskRunLog` sequence against a recording `driver.ExecuteLogger` |
| `Dump(ctx, writer)` | `Dump(ctx, writer, nil)` | the same statements. Bytebase iterates maps, so compare the statements as a set unless the recording shows one order |

Then mutation-check the harness, which is the only proof it can fail. Apply each
mutation in the omni worktree, run the parity test, and require a failure; then
`git checkout -- databricks/`:

1. stop treating an empty non-STRING cell as NULL in `toValue`;
2. read a `LONG` as `Int32` in `toValue`;
3. classify `INSERT` as returning rows in `classifyStatement`;
4. prepend `EXPLAIN FORMATTED` to a statement that already starts with
   `EXPLAIN`;
5. drop the `MaxResultBytes` truncation in `resultFromResponse`;
6. return the whole `SELECT VERSION()` answer from `versionFromResult`.

## Step 9 — what a failure means

- **The service rejects a request this package sends.** The SDK request is
  wrong (catalog, schema, warehouse, statement). Ask the workspace with a SQL
  client and fix `query.go`, `sync.go` or `dump.go`; do not add a retry.
- **The service accepts, the metadata differs from Bytebase's.** Bytebase's
  mapping is what customers see today; record the difference, decide which is
  right against the workspace, and note it in the PR. A change in what a synced
  column, view or materialized view says is a behavior change for every
  Databricks instance.
- **The parity harness and the tape disagree.** The tape holds an answer
  recorded from a corpus that has moved. Re-record; never edit a tape.
- **The live test fails before it sends anything** (a host, a token, a
  warehouse, a catalog). That is the account or the DSN, not the driver.
- **The Dump's materialized-view branch is refused.** Bytebase's shape is
  unconfirmed; fix it against the server's own statement and record the
  difference. Do not delete the branch to make the dump pass.

## Step 10 — after the recording landed

- `docs/engine-test-matrix.md` names the pinned workspace, SQL release, tapes,
  sync golden and parity result.
- The package doc and this file no longer call the driver unverified.
- `AGENTS.md` lists `databricks` among the engines whose tests replay a tape and
  records that its client is HTTP.
- `databricks/testdata/pins.tsv` and `databricks/internal/testcontainer` are in
  the tree.

## Differences from Bytebase's driver, each to confirm on a workspace

The port is close to `backend/plugin/db/databricks` and diverges where the
contract or determinism requires it. Every one is a decision the recording
settles.

| difference | why |
| --- | --- |
| `SyncInstance` returns catalog names only | the contract says databases without schemas; Bytebase walked every schema and table of every catalog to produce a list the caller only uses for names |
| listings are sorted | Bytebase iterated maps, so its metadata and dump order changed between runs; a tape and a golden need one order |
| `QueryOptions.Schema` is sent with each statement | the SQL Statement Execution API takes a schema per request; Bytebase sent only the catalog |
| a zero-row query keeps its columns | Bytebase dropped the result when `DataArray` was empty, losing the column list; the contract keeps the result set |
| `QueryOptions.MaxResultBytes` stops reading | Bytebase had no cap; the contract asks for one, counted while the rows are read |
| `QueryOptions.Explain` is `EXPLAIN FORMATTED` and leaves an existing plan alone | Databricks' plan is Spark's FORMATTED tree; Bytebase's standard rewrite also avoided planning a plan |
| `Execute` returns 0 and logs each command | the API reports no affected count on this path; Bytebase returned 0 and logged nothing |
| `Dump` sorts and quotes identifiers | Bytebase's map iteration and raw interpolation were not stable; the workspace accepted the quoted form |
| a materialized view is the server's `SHOW CREATE TABLE` | the workspace states the query and columns; Bytebase's `DESC FORMATTED` rendering had no `AS SELECT` and was not restorable |
| a 404 joins `driver.ErrNotFound` | Bytebase returned a plain error; the contract gives callers a not-found marker |
| `Config.Dial` and `Config.TLS` reach the SDK's HTTP transport | Bytebase had no caller dialer or per-host TLS; the contract does |

## Result

The six live tests passed against the scratch workspace (`SELECT VERSION()`
answered `4.2.0`), the six tapes replay without credentials, two consecutive
recordings wrote the same bytes, and `TestSyncDatabase` compares the synced
catalog with `testdata/sync/workspace.json`. The parity harness passed on the
sync, the instance, queries, `Execute` and the dump statements; it records the
materialized-view difference above. Five mutations of the driver made the
parity harness fail (`toValue`'s NULLs, a `LONG` read as `Int32`, explaining a
plan twice, dropping the result cap, and returning the whole version text). The
sixth, classifying `INSERT` as returning rows, is not observable on this
workspace because the service states `num_affected_rows` columns for an
`INSERT`; the package's `TestClassifyStatement` is what fails under it.
