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

## What the recording found (2026-11-17)

Recorded from a single-node `rg.large` cluster in us-west-2, `Redshift
1.0.539507` on PostgreSQL 8.0.2, pinned as `driver aws:redshift:1.0.539507` in
`redshift/testdata/pins.tsv` by the SHA-256 of `SELECT version()`. The cluster
was scratch and is deleted; the corpus is `redshift/driver/testdata/corpus.sql`.

defects this package had, each found only by asking the server:

- **The sync failed on every non-superuser login.** `getViews` and
  `getTableColumns` read Redshift's own schemas (`pg_auto_copy`, `pg_internal`),
  and `obj_description` on `pg_auto_copy.copy_job_detail` is `permission denied
  for schema pg_auto_copy` (SQLSTATE 42501). Only `pg_catalog` and
  `information_schema` were excluded. The predicate is now
  `substr(name, 1, 3) <> 'pg_' AND name <> 'information_schema'`: the server
  refuses the literal `pg_` prefix in `CREATE SCHEMA` (`unacceptable schema
  name ... The prefix "pg_" is reserved for system schemas`), while `pgdata` is
  a legal user schema, so a `pg_%` LIKE would hide it. `getSchemas` applied the
  same filter, because `SVV_ALL_SCHEMAS` reports `pg_catalog` and
  `information_schema` and they were being emitted as empty user schemas.
- **Three sync queries had no order.** `getDatabases`, `getTables` and the
  `SVV_DATASHARES` read were unordered, and `getIndexes` iterated a Go map into
  each table's index list. A tape recorded twice writes different bytes. All
  now order, and the map is sorted.
- **`ANALYZE` was refused in a transactional script.** It shares VACUUM's parse
  node (`ast.VacuumStmt`), so `refuseImplicitCommit` rejected it. A cluster runs
  `BEGIN; CREATE TABLE ...; ANALYZE ...; ROLLBACK` and the table is gone, so
  ANALYZE does not commit; only `IsVacuumCmd` is refused now. `VACUUM` inside a
  transaction really is refused by the server (`VACUUM cannot run inside a
  transaction block`).
- **The row-estimate test made a temporary table** and read its estimate on a
  pooled connection that need not see it. It uses a permanent table now.
- **The creation test named its database from the clock**, which no tape can
  hold. It uses a fixed name and drops it first.

what the server disagreed with in these documents:

- **`ON DELETE CASCADE` is not supported.** Step 2's corpus failed with
  `ERROR: ON DELETE is not supported`; the constraint states no action, and the
  corpus says so.
- **A primary key is an index.** "Redshift has no CREATE INDEX, so ... the
  sync's index read records an empty list" is wrong: `pg_indexes` and
  `pg_index` report `lookup_pkey` and `facts_pkey`, and the golden holds them.
  Redshift has no *secondary* index, which is what the corpus writes none of.
- **`PG_VIEWS.definition` is not truncated at 256 here.** A 546-character view
  is stated whole, byte-for-byte as `SHOW VIEW` states it. `SHOW VIEW` remains
  correct, but the truncation it was ported to work around does not reproduce
  on 1.0.539507; the golden pins a long definition coming back whole.
- **`version()` ends with a NUL byte.** The server's `length(version())` is 123
  and `ascii(substr(version(), 123, 1))` is 0; a driver receives 122 bytes. The
  pin digest hashes what the driver reads.
- **`SHOW server_version` needs a superuser** (`must be superuser to examine
  "server_version"`); the master user is not one. `getVersion`'s choice of
  `SELECT version()` is not optional.
- **A foreign key's `confmatchtype` is `u`** (unspecified) for a constraint that
  states no match type, where PostgreSQL says `s` (simple). The sync passes the
  server's letter through; the contract calls the field PostgreSQL-specific.
- A public endpoint accepted TCP 5439 a few minutes *after* `describe_clusters`
  said `available`; the first `psql` connection timed out.

## What is here, and what is missing

| capability | state |
| --- | --- |
| `SyncInstance`, `SyncDatabase` | ported from `backend/plugin/db/redshift`, recorded from a cluster |
| `Session.Query` (decode, limit, explain, `MaxResultBytes`, search path) | ported, recorded from a cluster |
| `Execute` (one transaction, autocommit, `CreateDatabase`) | ported, recorded from a cluster |
| `CountAffectedRows` | ported, recorded from a cluster |
| `Dump` | **not implemented**: `redshift/catalog` has no `DatabaseDefinition` |
| tapes (`testdata/tapes`) | six, recorded on 2026-11-17 and byte-identical across two recordings |
| sync goldens (`testdata/sync`) | `dev.json`, compared with `proto.Equal` |
| `driver` row in `redshift/testdata/pins.tsv` | `aws:redshift:1.0.539507`, digest of `SELECT version()` |
| live tests (`live_test.go`) | replay tapes with no credentials; record with `OMNI_RECORD=1` and `REDSHIFT_DSN` |

The unit tests that do run without a server cover the connection configuration,
statement classification, the value conversions, the row estimate's plan
reading, the foreign-key and version helpers, and
`redshift/analysis.SplitSearchPath`. The recorded tests replay the tapes and
need nothing.

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
reads whatever is there and a tape of an empty database pins nothing. The
corpus is committed as `redshift/driver/testdata/corpus.sql`; load it into an
empty database before recording:

```bash
psql "$REDSHIFT_DSN" -v ON_ERROR_STOP=1 -f redshift/driver/testdata/corpus.sql
```

The corpus writes no `ON DELETE` action, because Redshift refuses one
(`ON DELETE is not supported`). A primary key is the only index a Redshift
table reports, and the sync reads it: `lookup_pkey` and `facts_pkey` are in the
golden. What a table also declares is its `DISTKEY` and `SORTKEY`, which the
sync does not read yet.

The datashare branch (`Params["datashare"]`) is the one path this corpus cannot
reach: it needs a second account or namespace that shares a database in, which
`SVV_DATASHARES` then reports. `SVV_ALL_TABLES` lists a shared view with
`table_type` `VIEW`, and the sync leaves it out: a datashare has no view
definition to report, so the recording is what decides whether a shared view
should be a view with a definition here or not synced at all. **The 2026-11-17
recording skipped it**: one scratch cluster in one account has no share, so
`TestSyncDatabaseDatashare` skips and has no tape, and the `SVV_ALL_*` queries
were read only against an ordinary database. `TestSyncDatabaseDatashare` reads
the share out of the cluster instead of being told its name: it skips, naming
this provisioning step, when `SVV_DATASHARES` reports no consumer database.
It is the one test that still reads `REDSHIFT_DSN` rather than a tape.

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
| 1 | `TestSyncDatabase` | the catalog and `pg_catalog` queries and the metadata they build: schemas, tables, columns, indexes, foreign keys, views with their `SHOW VIEW` definitions, comments, row estimates |
| 2 | `TestSyncInstance` | `SELECT version()`, `pg_user`, `pg_database`, `SVV_DATASHARES`; the excluded databases |
| 3 | `TestQuery` | statement classification (SELECT/SHOW/EXPLAIN vs commands), every value kind, the LIMIT rewrite, the EXPLAIN text, `MaxResultBytes` truncation, `QueryOptions.Schema`'s search path, error-halting |
| 4 | `TestExecute` | one transaction with its BEGIN/COMMIT, the rollback of a failing script, autocommit leaving earlier statements committed |
| 5 | `TestExecuteCreateDatabase` | `CREATE DATABASE` and the second run's no-op |
| 6 | `TestCountAffectedRows` | `EXPLAIN` after `ANALYZE`, and a statement with no estimate |
| 7 | `TestSyncDatabaseDatashare` | the `SVV_ALL_*` branch, on a shared-in database |

The table names each test by what it records, which is the name it takes once it
replays a tape. `live_test.go` holds `TestSyncDatabase`, `TestSyncInstance`,
`TestQuery`, `TestExecute`, `TestExecuteCreateDatabase`, `TestCountAffectedRows`
and `TestSyncDatabaseDatashare`, and the tapes are named for them.

Recording commands, one test at a time so a failure is attributable. A `-run`
pattern that matches nothing exits 0 and records nothing, which reads as a pass:

```bash
OMNI_RECORD=1 REDSHIFT_DSN="$REDSHIFT_DSN" go test ./redshift/driver/ -run '^TestSyncDatabase$' -count=1 -v
OMNI_RECORD=1 REDSHIFT_DSN="$REDSHIFT_DSN" go test ./redshift/driver/ -run '^TestSyncInstance$' -count=1 -v
OMNI_RECORD=1 REDSHIFT_DSN="$REDSHIFT_DSN" go test ./redshift/driver/ -run '^TestQuery$' -count=1 -v
OMNI_RECORD=1 REDSHIFT_DSN="$REDSHIFT_DSN" go test ./redshift/driver/ -run '^TestExecute$' -count=1 -v
OMNI_RECORD=1 REDSHIFT_DSN="$REDSHIFT_DSN" go test ./redshift/driver/ -run '^TestExecuteCreateDatabase$' -count=1 -v
OMNI_RECORD=1 REDSHIFT_DSN="$REDSHIFT_DSN" go test ./redshift/driver/ -run '^TestCountAffectedRows$' -count=1 -v
OMNI_RECORD=1 REDSHIFT_DSN="$REDSHIFT_DSN" go test ./redshift/driver/ -run '^TestSyncDatabaseDatashare$' -count=1 -v
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

Committed as `redshift/driver/testdata/sync/dev.json`: the protojson of
`SyncDatabase`'s output, written by `checkSyncGolden` in a recording and
compared with `proto.Equal` in a replay rather than byte-for-byte, because
protojson spacing moves between builds. The datashare golden does not exist,
because Step 2 produced no share.

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

## What the parity harness found (2026-11-17)

Run against the same cluster, with Bytebase's `backend/plugin/db/redshift` at
`origin/main` (`db28e088d6`) and `go.mod` replacing omni with the recording
worktree:

| area | result |
| --- | --- |
| `SyncInstance` | the same version text (`Redshift 1.0.539507 based on PostgreSQL 8.0.2`), the same user databases (`dev`) and the same roles and attributes |
| `SyncDBSchema` vs `SyncDatabase` | Bytebase's sync **succeeded**; the metadata agreed on every table, column, index, foreign key, view definition, comment and character set except below |
| `QueryConn` vs `Session.Query` | column names and `TypeName`s, every cell, the rewritten statement (LIMIT included), the EXPLAIN text and the error text agreed; the two differences below are the contract's |
| `Execute` | the same affected rows (3), and the same count (6) after both rolled a failing script back |
| `CountAffectedRows` | 100 from both, off the same plan |

What differed, and why it is not a defect:

- **A synced column keeps precision and scale in omni.** Bytebase syncs
  `numeric`; omni syncs `numeric(12,2)` and `numeric(14,2)`. The server's
  `data_type` says `numeric` and states the precision and scale beside it, so
  Bytebase's value loses what the column declares. This is one of the port's
differences, confirmed by the cluster.
- **Bytebase emits two empty schemas omni omits.** `pg_catalog` and
  `information_schema` are in `SVV_ALL_SCHEMAS`, and Bytebase reports both with
  no tables or views; omni excludes them. A sync reports user schemas.
- **An affected-row result is a column in Bytebase and a count in omni.**
  Bytebase answers `CREATE TEMP TABLE` with one `Affected Rows` INT column and a
  row of `0`; omni sets `Result.AffectedRows` and returns no columns, as the
  contract states.
- **A plan is marked in omni.** Bytebase returns EXPLAIN's rows with
  `QueryPlan` unset; omni sets `Result.Plan.Format = "TEXT"`. The rows are the
  same plan lines.

What the harness did **not** exercise: an inbound datashare (no share existed,
so the rewrite and `SVV_ALL_*` are uncompared), `SHOW VIEW` against a view the
login does not own, `MaxResultBytes` truncation, a search path across two
schemas, and the six mutation checks above. Bytebase's driver also never failed
on the corpus, so the permission bug the sync had is fixed by omni's exclusion
rather than observed in both.

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

The port is close to `backend/plugin/db/redshift` and diverges in twenty places.
Every one is a decision the recording settles.

| difference | why |
| --- | --- |
| `QueryOptions.Schema` sets the search path before each statement, not once per call | a `Session` outlives a call, so a path set for one schema must not carry into the next |
| `Execute` always splits the script; Bytebase collapsed a script over 2 MB or 200 commands into one request | the sheet-size bounds are the caller's; the splitter is lexical and never fails |
| the database `CREATE DATABASE` names is read from the parse tree; Bytebase trimmed the literal prefix `CREATE DATABASE`, so `create database x` named a database called `create` | a statement Redshift accepts must not be misread by case |
| `padZeroes` treats a leading `-` as the value's sign, not as a zone offset | Bytebase's version panicked on `INTERVAL '-00:02:45.25'` (`slice bounds out of range`); `pg/driver` carries the guard |
| a datashare database is `Params["datashare"]`, not a field on the connection config | the contract carries engine-specific settings in `Params`, as `cosmosdb/driver` does for credentials; the flag is read by this driver and is not sent to the server as a startup parameter |
| autocommit `Execute` runs the whole script on one connection | each call would borrow a pooled connection of its own, so a `SET search_path` or a temporary table would not hold for the statements after it |
| a datashare statement loses the database's name from the relations its `FROM` and `JOIN` clauses name, and from nothing else | Bytebase replaced the text everywhere, including inside a literal and in a column qualifier; a qualified name loses its first component only, so `shared.shared.facts` keeps the schema |
| a foreign key's table names are the `pg_class.relname` the constraint joins, not the tail of a quoted `schema.table` | a schema holding a period broke that split, so the constraint was filed under a table that does not exist and dropped from the sync |
| a foreign key's column lists are found by their parentheses and split between identifiers | a column named `a,b` was read as two columns, `a"b` kept its doubled quote, and a name holding a parenthesis ended the list early |
| a datashare's columns are read in ordinal order | `SVV_ALL_COLUMNS` answers in no order a tape can hold, and the catalog consumes the slice's order |
| an expression's `FROM` is not a relation list, and the list ends at the clause that follows it, `QUALIFY` included | `EXTRACT(x FROM db.t)` holds a `FROM` that opens an argument, and `QUALIFY` after a `FROM` holds column qualifiers, so the datashare rewrite leaves both alone |
| a synced column keeps its numeric precision and scale | `data_type` states `numeric` alone, so `numeric(12,2)` and `numeric(8,4)` synced alike |
| a view's definition is read from `SHOW VIEW`, not `PG_VIEWS.definition` | `PG_VIEWS` states it in a varchar(256) and truncates a longer one, and a truncated definition is not the view; the recording settles what `SHOW VIEW` answers and that a long view's text comes back whole |
| a procedure call is read as a query though it is not read-only | `CALL` answers its OUT and INOUT parameters as one row, which the classification now says independently of whether the statement may run read-only |
| a procedure call's rows are withheld from a caller that masks columns | an OUT or INOUT value has no column lineage to trace, so a masked caller is given the affected-row count instead |
| a reference that begins with a relation's alias or name is left alone | `FROM app.orders AS shared, shared.products` reads a SUPER column of the aliased relation, so its text is not the database's name |
| a database-creation script may hold nothing after `CREATE DATABASE` | the step creates the database, and a statement that followed could not run in it on a retry that cannot tell a created database from an initialized one; initialization is a script of its own |
| a transactional script refuses TRUNCATE, VACUUM, a database statement and its own transaction control | Redshift commits the first three on their own and the last would end the transaction the call owns, so the script cannot be rolled back as a whole; the caller runs it in autocommit mode |
| a `dbname` in Params is the database the connection uses | Params are applied last, so their value overrides Config.Database and Redshift's default |
| a derived table's alias is collected like a base table's | the alias after a subquery names the derived table, so a reference beginning with it is a column, not the database |
| a `dbname` in Params is read as pgx parsed it, quotes included in the libpq value | the parameter is connection-string text, so `'analytics'` names `analytics` |
| a transactional script refuses its own transaction control | a `BEGIN`, `COMMIT`, `END` or `ROLLBACK` in the script would decide the transaction the call promised to own |
| a `user` or `password` in Params is the login the connection uses | Params are applied last, so their values override Config's, and require no Config user |
| a `FETCH` is read as a query and a `MOVE` is not | a fetch answers a cursor's rows; a move only advances the cursor |
| a cursor's rows are withheld from a caller that masks columns | the cursor's column lineage is not the statement's, so a masked caller is given the affected-row count |
| the search-path prefix is split at the statement boundary the splitter finds | a schema's quoted name may hold the delimiter text, which a search for it would cut inside the identifier |

## What would unblock it

A cluster: one dedicated scratch Redshift with `REDSHIFT_DSN` set. The seven
recordings above take about a day; everything else in this package is already
written and gated. If no cluster is going to exist, the alternative is to drop
Redshift from the first wave and leave `backend/plugin/db/redshift` where it is —
it has run for years against real instances, which is more than this package can
claim.
