# Verifying databricks/driver

`databricks/driver` is a port, not a driver anybody has seen work. This file is
what turns a workspace into a day of recording instead of a week of porting: it
lists what is missing, the account and warehouse to pin, the tests to record in
order, the tape harness at the SDK's HTTP seam, and the parity harness that
holds this package to Bytebase's driver.

Nothing in this plan runs in CI. Databricks is an account service: omni cannot
start a workspace, SQL warehouse or Unity Catalog in a container, and no
PostgreSQL or Hive stand-in answers the Unity Catalog or Statement Execution
APIs. `databricks/parser` does not exist either, so there is no grammar corpus
to judge. The live tests in this package are written but have never been run;
they skip with `DATABRICKS_DSN` when it is unset.

## What is here, and what is missing

| capability | state |
| --- | --- |
| `SyncInstance` (SQL version, catalog names) | ported, needs a workspace |
| `SyncDatabase` (schemas, tables, views, materialized views, external tables) | ported, needs a workspace |
| `Session.Query` (types, NULL, limit, explain, `MaxResultBytes`) | ported, needs a workspace |
| `Execute` (statement by statement, no transaction control) | ported, needs a workspace |
| `Dump` (`SHOW CREATE TABLE`, `DESC FORMATTED`) | ported, needs a workspace |
| tapes (`testdata/tapes`) | **none**: nothing has been recorded |
| sync goldens (`testdata/sync`) | **none**: same reason |
| `driver` row in `databricks/testdata/pins.tsv` | **none**: no workspace to pin |
| live tests (`live_test.go`) | written, never run; skip without `DATABRICKS_DSN` |
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
release is the first word, which is the version `SyncInstance` keeps.

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
  SELECT id, amount FROM `<catalog>`.`omni_record`.`facts` WHERE id > 0;

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

`databricks/testdata/pins.tsv` does not exist yet. Add it with the columns the
other engine pins use (`role`, `image`, `digest`, `source tag`):

```
# columns: role, image, digest, source tag
driver	databricks:4.2.0	sha256:<sha256 of SELECT VERSION()'s answer>	2026-11-17
```

There is no image to pin, so `image` holds the first word the server answered and
`digest` holds the SHA-256 of the whole answer. The recording PR decides the
exact spelling; nothing in the driver depends on it today.

## Step 4 — the tape harness

The driver talks HTTP through the Databricks SDK, so the tape belongs at the
SDK's HTTP transport, not at a database/sql connection:

- Add `databricks/internal/testcontainer`, owned by this engine, with a
  `RoundTripper` that records one request (method, URL, headers the driver
  sets, body) and one answer (status, content type, body). `tape.OpenRaw` is the
  shape other SDK-backed engines use; use an HTTP request/answer encoding, not a
  database/sql one.
- Add `Open(t)` returning a `*databricks.WorkspaceClient` whose
  `Config.HTTPTransport` is the tape. The package's tests are in package
  `driver`, so a taped test constructs `&Driver{client: testcontainer.Client(t),
  warehouseID: ..., catalog: ...}` directly.
- `live` builds the same client against `DATABRICKS_DSN` and **fails** when it
  is unset. A recording without a workspace must not skip, or `OMNI_RECORD=1`
  would write an empty tape and call it a pass.
- Add `main_test.go` with `func TestMain(m *testing.M) { os.Exit(tape.Main(m)) }`.
- Add a codec round-trip test for the request/answer encoding, and a test that
  an unset DSN fails in a recording rather than skipping.

The request/answer tape must narrow what the service makes up:

- `statement_id` is random and names a later poll; replace it with a numbered
  stand-in and keep the poll's path consistent with the stand-in.
- `created_at`, timings and `warehouse` state move between runs; keep only what
  the driver reads, or replace them with stable values.
- `result.data_array` chunk indexes, `next_chunk_index`,
  `next_chunk_internal_link` and `total_*_count` move with the chunking; keep
  the first chunk and record what the driver actually reads.
- The `DESC FORMATTED` answer carries a table's location, owner and creation
  time; a recording of the same table on another day changes them. Either keep
  them out of the fixture's properties or use a stand-in.
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
| 5 | `TestDump` | `SHOW CREATE TABLE` for managed tables, external tables and views, `DESC FORMATTED` for a materialized view, and the dump text |
| 6 | `TestPing` | `Workspace.ListAll` as the connectivity check |

Recording commands, one test at a time so a failure is attributable. The live
tests are named `TestLive*` today; rename them to the recorded names above and
point them at `testcontainer.Open` when the harness lands.

```bash
export DATABRICKS_DSN='databricks://...'
OMNI_RECORD=1 DATABRICKS_DSN="$DATABRICKS_DSN" go test ./databricks/driver/ -run '^TestLiveSyncDatabase$' -count=1 -v
OMNI_RECORD=1 DATABRICKS_DSN="$DATABRICKS_DSN" go test ./databricks/driver/ -run '^TestLiveSyncInstance$' -count=1 -v
OMNI_RECORD=1 DATABRICKS_DSN="$DATABRICKS_DSN" go test ./databricks/driver/ -run '^TestLiveQuery$' -count=1 -v
OMNI_RECORD=1 DATABRICKS_DSN="$DATABRICKS_DSN" go test ./databricks/driver/ -run '^TestLiveExecute$' -count=1 -v
OMNI_RECORD=1 DATABRICKS_DSN="$DATABRICKS_DSN" go test ./databricks/driver/ -run '^TestLiveDump$' -count=1 -v
```

Then record everything twice and require byte-identical tapes. Two runs of one
compiled test binary are deterministic; two `go test` runs are byte-identical in
the tapes but not in a protojson golden, so compare the sync goldens as protos:

```bash
go test -c -o /tmp/databricks-driver.test ./databricks/driver/
OMNI_RECORD=1 DATABRICKS_DSN="$DATABRICKS_DSN" /tmp/databricks-driver.test -test.count=1
sha256sum databricks/driver/testdata/tapes/*.tape > /tmp/databricks-tapes.a
OMNI_RECORD=1 DATABRICKS_DSN="$DATABRICKS_DSN" /tmp/databricks-driver.test -test.count=1
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

- `SHOW CREATE TABLE` for a managed table, an external table and a view.
- `DESC FORMATTED` for a materialized view; the driver writes a
  `CREATE MATERIALIZED VIEW` with the comment and table properties Bytebase
  wrote. Bytebase's property list is not parenthesized or quoted, and this port
  kept that shape. If the workspace refuses it, the dump's materialized-view
  branch is wrong and must be fixed against the server's `SHOW CREATE` (or
  against `DESC FORMATTED` fields that state the view's query).
- The dump writes `CREATE CATALOG` for every catalog and `CREATE SCHEMA` for
  every schema but `default`, and skips the `system` catalog and each catalog's
  `information_schema`. Confirm that a replayed dump creates the objects in the
  order a catalog accepts them.
- Add `TestDumpRestores` only after the dump text is known to apply: drop the
  scratch catalog's objects and run the dump back. Do not add a `testdata/dump`
  fixture before then.

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

## Step 10 — after the recording lands

- `docs/engine-test-matrix.md`: add the `databricks` row's pinned workspace
  release and name the tapes, the sync goldens and the parity result.
- Delete the package doc's "Unverified against a workspace" section, not soften
  it, in the PR that records.
- `AGENTS.md`'s list of engines whose tests replay a tape gains `databricks`, and
  its driver-client sentence gains the Databricks SDK.
- Add `databricks/testdata/pins.tsv` and `databricks/internal/testcontainer` in
  the recording PR.

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
| `Dump` sorts and quotes identifiers, and sorts table properties | Bytebase's map iteration and raw interpolation were not stable; the recording decides whether the quoted form is what the workspace accepts |
| a 404 joins `driver.ErrNotFound` | Bytebase returned a plain error; the contract gives callers a not-found marker |
| `Config.Dial` and `Config.TLS` reach the SDK's HTTP transport | Bytebase had no caller dialer or per-host TLS; the contract does |

## What would unblock it

A workspace: one dedicated scratch Databricks workspace with a personal access
token, a SQL warehouse and a catalog, and `DATABRICKS_DSN` set. The six
recordings above take about a day; everything else in this package is already
written and gated. If no workspace is going to exist, the alternative is to drop
Databricks from the first wave and leave `backend/plugin/db/databricks` where it
is — it has run against real workspaces, which is more than this package can
claim.
