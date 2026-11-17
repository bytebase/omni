# Verifying redshift/driver

`redshift/driver` is a port, not a driver anybody has seen work. This file is
what turns a cluster into a day of recording instead of a week of porting: it
lists what is missing, the cluster to pin, the tests to record in order, the
parity harness that holds this package to Bytebase's, and what to do when one of
them disagrees.

Nothing in this plan runs in CI. Redshift has no image omni can start — it is
proprietary, its wire protocol is PostgreSQL 8.0.2, and none of its catalog
(`SVV_ALL_SCHEMAS`, `SVV_ALL_TABLES`, `SVV_ALL_COLUMNS`, `SVV_DATASHARES`) exists
on PostgreSQL. A PostgreSQL stand-in therefore fails the sync for the wrong
reason and proves nothing about Redshift, which is why `redshift/parser` is
judged by the AWS docs and why this package's tests skip.

## What is here, and what is missing

| capability | state |
| --- | --- |
| `SyncInstance`, `SyncDatabase` | ported from `backend/plugin/db/redshift`, needs a cluster |
| `Session.Query` (decode, limit, explain, `MaxResultBytes`, search path) | ported, needs a cluster |
| `Execute` (one transaction, autocommit, `CreateDatabase`) | ported, needs a cluster |
| `CountAffectedRows` | ported, needs a cluster |
| `Dump` | **not implemented**: `redshift/catalog` has no `DatabaseDefinition` |
| tapes (`testdata/tapes`) | **none**: nothing has been recorded |
| sync goldens (`testdata/sync`) | **none**: same reason |
| `driver` row in `redshift/testdata/pins.tsv` | **none**: no server to pin yet |
| live tests (`live_test.go`) | written, never run; skip without a DSN |

The unit tests that do run cover the connection configuration, statement
classification, the value conversions, the row estimate's plan reading, the
foreign-key and version helpers, and `redshift/analysis.SplitSearchPath`.

## Step 1 — the cluster

Pin one dedicated, scratch-only cluster. The tests create temp tables, `Execute`
writes, and a future dump will execute DDL, so it must be a cluster nobody else
uses. Either form works:

- a provisioned cluster (`ra3.xlplus` is enough) with a database named
  `omni_record`;
- a Redshift Serverless workgroup with base capacity 8.

Record what it answers, because that is what the pin is:

```bash
psql "$REDSHIFT_DSN" -c 'SELECT version(), current_database();'
```

Redshift's own version is the tail of that string, `Redshift 1.0.xxxxx`;
`redshift/driver`'s `getVersion` parses it. Write the release, the cluster kind
and the region into the pin row's image column in Step 3.

The login needs to read `pg_catalog`, `INFORMATION_SCHEMA` and the `SVV_*` views
and to create, write and drop in a schema of its own, and to create and drop a
database of its own for the creation step. If the sync's
`obj_description` and `pg_get_constraintdef` calls come back NULL, the account
lacks ownership of the objects and the sync is right to complain (`pg_views`
with a NULL definition fails on purpose).

Two environment variables reach a cluster. Use the same one for both:

```bash
export REDSHIFT_DSN='postgres://USER:PASSWORD@HOST:5439/omni_record?sslmode=require'
export REDSHIFT_COMPAT_DSN="$REDSHIFT_DSN"
```

`REDSHIFT_DSN` is what `live_test.go` reads; `REDSHIFT_COMPAT_DSN` is the
variable omni's compat suite already names. Pointing both at the same cluster is
the point: the statements the driver sends, and the statements the parser
accepts, are then judged by the same engine.

## Step 2 — the fixtures

Load one corpus into the scratch database before recording, because the sync
reads whatever is there and a tape of an empty database pins nothing:

```sql
CREATE SCHEMA app;
CREATE TABLE public.lookup (id int PRIMARY KEY, label varchar(20) NOT NULL DEFAULT 'x');
COMMENT ON TABLE public.lookup IS 'a synced comment';
COMMENT ON COLUMN public.lookup.label IS 'a column comment';
CREATE TABLE app.facts (
  id bigint PRIMARY KEY,
  lookup_id int REFERENCES public.lookup (id) ON DELETE CASCADE,
  amount numeric(12,2),
  flag boolean,
  seen timestamp,
  notes varchar(80)
);
CREATE INDEX facts_lookup_id ON app.facts (lookup_id);
CREATE VIEW app.recent AS SELECT id, amount FROM app.facts WHERE seen > '2026-01-01';
```

The datashare branch (`Params["datashare"]`) is the one path this corpus cannot
reach: it needs a second account or namespace that shares a database in, which
`SVV_DATASHARES` then reports. Record it last, on its own, and skip it if
provisioning one is unattractive — say so in the PR that records, rather than
leaving a tape that pretends. `TestLiveSyncDatabaseDatashare` reads the share
out of the cluster instead of being told its name: it skips, naming this
provisioning step, when `SVV_DATASHARES` reports no consumer database.

## Step 3 — the pin row

`redshift/testdata/pins.tsv` holds `role, image, digest, source tag`. The
grammar keeps role `default`; add a second row for the cluster:

```
# columns: role, image, digest, source tag
driver	aws:redshift:1.0.xxxxx	sha256:<digest of the recorded version string>	2026-11-17
```

The digest column has no image digest to hold, so it holds the SHA-256 of
`SELECT version()`'s answer as recorded. That is the same idea as the grammar's
`docs` row: the pin is whatever the answers were recorded from, and a cluster of
another release fails the pin instead of quietly replaying. Decide this in the
recording PR; nothing here depends on the format today.

## Step 4 — the tape harness

Tapes need somewhere to record through. Add `redshift/internal/testcontainer`,
modeled on `pg/internal/testcontainer` and owned by this engine (omni keeps no
helper shared across engines):

- `container.go`: `var server = &tape.Server{Pins: Pins, Live: live, Codec: Codec{}}`;
  `Pins` is `redshift/testdata/pins.tsv` resolved from `runtime.Caller(0)`, a
  `const Driver = "driver"` names the pin row (`cassandra`, `elasticsearch` and
  `redis` each keep that const in their helper), and `Open(t)` returns
  `tape.Open(t, server, Driver)`;
  `Codec` encodes and decodes `*pgconn.PgError` (error fields, not the message
  alone) so a replayed test's `errors.As` finds the same type;
  `live` connects to `REDSHIFT_DSN` and **fails** when it is unset — a recording
  without a cluster must not skip, or `OMNI_RECORD=1` would produce an empty
  tape and call it a pass.
- `main_test.go`: `func TestMain(m *testing.M) { os.Exit(tape.Main(m)) }`.
- `container_test.go`: a codec round-trip, and a test that `Open` caps the pool
  at one connection.

Then change `live_test.go`'s scenarios into recorded tests that call
`testcontainer.Open(t)`: the driver under test is `&Driver{db: testcontainer.Open(t), databaseName: ...}`,
exactly as `pg/driver/query_test.go` opens its driver.

## Step 5 — record, in this order

Sync first: it is the read-only half and the half Bytebase never tested. Then
the reads, then the writes, then the estimate.

| order | test | what it records |
| --- | --- | --- |
| 1 | `TestSyncDatabase` | the catalog and `pg_catalog` queries and the metadata they build: schemas, tables, columns, indexes, foreign keys, views, comments, row estimates |
| 2 | `TestSyncInstance` | `SELECT version()`, `pg_user`, `pg_database`, `SVV_DATASHARES`; the excluded databases |
| 3 | `TestQuery` | statement classification (SELECT/SHOW/EXPLAIN vs commands), every value kind, the LIMIT rewrite, the EXPLAIN text, `MaxResultBytes` truncation, `QueryOptions.Schema`'s search path, error-halting |
| 4 | `TestExecute` | one transaction with its BEGIN/COMMIT, the rollback of a failing script, autocommit leaving earlier statements committed |
| 5 | `TestExecuteCreateDatabase` | `CREATE DATABASE` and the second run's no-op |
| 6 | `TestCountAffectedRows` | `EXPLAIN` after `ANALYZE`, and a statement with no estimate |
| 7 | `TestSyncDatabaseDatashare` | the `SVV_ALL_*` branch, on a shared-in database |

The table names each test by what it records, which is the name it takes once it
replays a tape. `live_test.go` holds `TestLiveSyncDatabase`,
`TestLiveSyncInstance`, `TestLiveQuery`, `TestLiveExecute`,
`TestLiveExecuteCreateDatabase`, `TestLiveCountAffectedRows` and
`TestLiveSyncDatabaseDatashare` today, so the commands below name those, and
step 4's conversion renames them. Every row has a test.

Recording commands, one test at a time so a failure is attributable. They name
the tests that exist today: a `-run` pattern that matches nothing exits 0 and
records nothing, which reads as a pass:

```bash
OMNI_RECORD=1 REDSHIFT_DSN="$REDSHIFT_DSN" go test ./redshift/driver/ -run '^TestLiveSyncDatabase$' -count=1 -v
OMNI_RECORD=1 REDSHIFT_DSN="$REDSHIFT_DSN" go test ./redshift/driver/ -run '^TestLiveSyncInstance$' -count=1 -v
OMNI_RECORD=1 REDSHIFT_DSN="$REDSHIFT_DSN" go test ./redshift/driver/ -run '^TestLiveQuery$' -count=1 -v
OMNI_RECORD=1 REDSHIFT_DSN="$REDSHIFT_DSN" go test ./redshift/driver/ -run '^TestLiveExecute$' -count=1 -v
OMNI_RECORD=1 REDSHIFT_DSN="$REDSHIFT_DSN" go test ./redshift/driver/ -run '^TestLiveExecuteCreateDatabase$' -count=1 -v
OMNI_RECORD=1 REDSHIFT_DSN="$REDSHIFT_DSN" go test ./redshift/driver/ -run '^TestLiveCountAffectedRows$' -count=1 -v
OMNI_RECORD=1 REDSHIFT_DSN="$REDSHIFT_DSN" go test ./redshift/driver/ -run '^TestLiveSyncDatabaseDatashare$' -count=1 -v
```

Then record everything twice and require byte-identical tapes. Two runs of one
compiled test binary are deterministic; two `go test` runs are byte-identical in
the tapes but not in a protojson golden, so compare the sync goldens as protos:

```bash
go test -c -o /tmp/redshift-driver.test ./redshift/driver/
OMNI_RECORD=1 REDSHIFT_DSN="$REDSHIFT_DSN" /tmp/redshift-driver.test -test.count=1
sha256sum redshift/driver/testdata/tapes/*.tape > /tmp/redshift-tapes.a
OMNI_RECORD=1 REDSHIFT_DSN="$REDSHIFT_DSN" /tmp/redshift-driver.test -test.count=1
sha256sum redshift/driver/testdata/tapes/*.tape | diff /tmp/redshift-tapes.a -
```

A tape that is not reproducible is a nondeterminism in what the test sends. The
usual causes, all met before on other engines: a sync query without `ORDER BY`
(heap order becomes the tape's order), a map iterated into a query, a value the
cluster makes up (a plan's cost, a timestamp, an OID, `pg_database` order).
Fix the driver, never the tape.

## Step 6 — the sync goldens

After the tapes, commit `redshift/driver/testdata/sync/<database>.json`: the
protojson of `SyncDatabase`'s output, compared with `proto.Equal` in the test
rather than byte-for-byte, because protojson spacing moves between builds. One
golden per database kind: the corpus database, and the datashare one if Step 2
produced it.

## Step 7 — the dump

`Dump` is not implemented and must not be declared until a metadata → DDL
renderer exists. `redshift/catalog` loads a definition from SQL (`LoadSQL`,
`LoadSDL`), executes DDL, diffs two catalogs and generates a migration, so the
cheap path looks open — `LoadMetadata`, `Diff`, `GenerateMigration` — and is
not: `LoadMetadata` installs a column's *coarse* type because completion reads
only the shape, so a dump built on it would write types the sync never read.
What is missing is the renderer `pg/catalog/definition.go` provides for the
parent engine. Two ways forward, in order of preference:

1. port `pg/catalog/definition.go` and the table/sequence/index/constraint
   renderers it calls into `redshift/catalog`, re-verified per engine (Redshift
   has no `serial`, its `IDENTITY` and `DISTKEY`/`SORTKEY` are its own, and its
   deparse of a query is `ruleutils.go`'s), then `Dump` becomes
   `catalog.DatabaseDefinition(db, catalog.DefinitionOptions{PrintHeader: true})`
   with `TestDump`, `TestDumpRestores` and `TestDumpGaps` beside it;
2. leave it out and say so, which is today's state: Bytebase's Redshift `Dump`
   wrote nothing either, so nothing is lost yet.

## Step 8 — the parity harness

The gate that a tape cannot be: Bytebase's driver and this one, on the same
cluster, compared cell by cell. It lives in a throwaway Bytebase worktree and is
never committed.

Start a worktree of the omni branch this file landed on, and a scratch Bytebase
worktree beside it (never the shared checkout):

```bash
git -C <omni repo root> worktree add ~/omni-worktrees/redshift-verify -b redshift-verify origin/redshift-prep
git -C <bytebase repo root> worktree add ~/bb-worktrees/redshift-parity -b redshift-parity origin/main
cd ~/bb-worktrees/redshift-parity
go mod edit -replace github.com/bytebase/omni=$HOME/omni-worktrees/redshift-verify
go mod tidy
```

Records happen in `~/omni-worktrees/redshift-verify`; the parity harness is
`backend/plugin/db/redshift/parity_omni_test.go` in `~/bb-worktrees/redshift-parity`,
and neither worktree is committed. Remove both when the recording lands.

Add `backend/plugin/db/redshift/parity_omni_test.go` (package `redshift_test`,
skipping without `REDSHIFT_DSN`) and compare, for every case in a corpus that
repeats the unit tests' statements plus the sync's objects:

| Bytebase | omni | comparison |
| --- | --- | --- |
| `(&Driver{}).Open(ctx, storepb.Engine_REDSHIFT, db.ConnectionConfig{...})` | `redshiftdriver.Open(ctx, driver.Config{...})` | both open, neither reaches the server |
| `SyncDBSchema(ctx)` | `SyncDatabase(ctx)` | same `*metadatapb.DatabaseSchemaMetadata`, `proto.Equal` |
| `SyncInstance(ctx)` | `SyncInstance(ctx)` | version text, database names and order, role names and attributes |
| `QueryConn(ctx, conn, sql, queryContext)` | `Session.Query(ctx, sql, queryOptions)` | one result per statement; column names and `TypeName`s; each cell's kind and value; `AffectedRows`; the rewritten `Statement` text; the error's text |
| `Execute(ctx, script, executeOptions)` | `Execute(ctx, script, executeOptions)` | total affected rows, and Bytebase's `TaskRunLog` sequence against a recording `driver.ExecuteLogger` (BEGIN, each command, COMMIT/ROLLBACK) |
| `CountAffectedRows(ctx, statement)` | `CountAffectedRows(ctx, statement)` | the estimate |

What is expected to differ, and must not be "fixed" away: omni's `Session` is
one connection per session so a search path does not leak between calls;
Bytebase's `QueryConn` sets it once per call; `Execute` also splits differently
(see the table below).

Then mutation-check the harness, which is the only proof it can fail. Apply each
mutation in the omni worktree, run the parity test, and require a failure; then
`git checkout -- redshift/`:

1. drop the `GREATEST(pc.reltuples::bigint, 0::bigint)` clamp in `getTables`;
2. report an empty `BIT` value as `Bytes` instead of NULL in `convertValue`;
3. classify `INSERT` as returning rows in `classifyQuery`;
4. insert the LIMIT after `FOR UPDATE` instead of before it in
   `redshift/analysis.RewriteResultLimit`;
5. read the table estimate from the last plan line instead of the first;
6. map a foreign key's `confdeltype` `c` to `RESTRICT` instead of `CASCADE`.

## Step 9 — what a failure means

- **The server rejects a statement this package sends.** A sync query is wrong.
  Ask the cluster with `psql` and fix `sync.go`; do not add a retry.
- **The server accepts, the metadata differs from Bytebase's.** Bytebase's
  mapping is what customers see today; record the difference, decide which is
  right against the cluster, and note it in the PR. A change in what a synced
  column or index says is a behavior change for every Redshift instance.
- **The parity harness and the tape disagree.** The tape holds an answer recorded
  from a corpus that has moved. Re-record; never edit a tape.
- **The live test fails before it sends anything** (a connection, a TLS
  handshake). That is the cluster or the DSN, not the driver.

## Step 10 — after the recording lands

- `docs/engine-test-matrix.md`: the `redshift` row says "no server" today; it
  becomes the pinned cluster's release, and the tapes, the sync goldens and the
  parity result are named there.
- The package doc's "Unverified against a cluster" section is deleted, not
  softened, in the PR that records.
- `AGENTS.md`'s list of engines whose tests replay a tape gains `redshift`, and
  the `REDSHIFT_COMPAT_DSN` sentence there gains `redshift/driver`'s live tests.

## Differences from Bytebase's driver, each to confirm on a cluster

The port is close to `backend/plugin/db/redshift` and diverges in five places.
Every one is a decision the recording settles.

| difference | why |
| --- | --- |
| `QueryOptions.Schema` sets the search path before each statement, not once per call | a `Session` outlives a call, so a path set for one schema must not carry into the next |
| `Execute` always splits the script; Bytebase collapsed a script over 2 MB or 200 commands into one request | the sheet-size bounds are the caller's; the splitter is lexical and never fails |
| the database `CREATE DATABASE` names is read from the parse tree; Bytebase trimmed the literal prefix `CREATE DATABASE`, so `create database x` named a database called `create` | a statement Redshift accepts must not be misread by case |
| `padZeroes` treats a leading `-` as the value's sign, not as a zone offset | Bytebase's version panicked on `INTERVAL '-00:02:45.25'` (`slice bounds out of range`); `pg/driver` carries the guard |
| a datashare database is `Params["datashare"]`, not a field on the connection config | the contract carries engine-specific settings in `Params`, as `cosmosdb/driver` does for credentials |

## What would unblock it

A cluster: one dedicated scratch Redshift with `REDSHIFT_DSN` set. The seven
recordings above take about a day; everything else in this package is already
written and gated. If no cluster is going to exist, the alternative is to drop
Redshift from the first wave and leave `backend/plugin/db/redshift` where it is —
it has run for years against real instances, which is more than this package can
claim.
