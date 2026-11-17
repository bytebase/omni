# Verifying bigquery/driver

`bigquery/driver` is a port, not a driver anybody has seen work. This file is what
turns a project into a day of recording instead of a week of porting: it lists what
is missing, the project and account to pin, the corpus, the tests to record in
order, the parity harness that holds this package to Bytebase's, and what to do
when one of them disagrees.

Nothing in this plan runs in CI. BigQuery has no image omni can start and no
emulator that answers its API — there is no local BigQuery — so every live test
skips with `BIGQUERY_CREDENTIALS`' name, and no tape, golden, `testdata` fixture or
pin row exists beside this package. Unlike `redshift/driver`, the obstacle is not
proprietary catalog views but the service itself: a recording needs a Google Cloud
project, a service account and a dataset. One run on 2026-11-17 had the first two
and not the third; what it answered, and what it left blocked, is in [the project
this run used](#the-project-this-run-used-on-2026-11-17) below.

## What is here, and what is missing

| capability | state |
| --- | --- |
| `SyncInstance` (the project's datasets) | ported from `backend/plugin/db/bigquery`, needs a project |
| `SyncDatabase` (INFORMATION_SCHEMA.COLUMNS and each table's metadata) | ported, needs a dataset |
| `Session.Query` (decode, limit wrap, dry-run explain, `MaxResultBytes`, plan) | ported, needs a project |
| `Execute` (one job for the whole script) | ported, needs a project |
| `Dump` | ported as a **no-op**: `googlesql` has no catalog, and Bytebase's BigQuery `Dump` wrote nothing |
| row-estimate (`CountAffectedRows`) | **not ported**: BigQuery has no EXPLAIN to read one from, and Bytebase counted none |
| client seam for tapes | **absent**: `Open` builds the real `*bigquery.Client`; step 4 adds the seam |
| tapes (`testdata/tapes`) | **none**: nothing has been recorded |
| sync goldens (`testdata/sync`) | **none**: same reason |
| `driver` row in `bigquery/testdata/pins.tsv` | **none**: no project to pin yet |
| live tests (`live_test.go`) | run against `bytebase-dev` on 2026-11-17: the instance listing, the reads, the dry run, the result cap, the plans and the not-found mapping ask a real project and agree with the port apart from the three defects that run found and fixed. The dataset sync and the writes stay blocked on a dataset grant |

The unit tests that do run cover the client configuration and the four credential
document types, statement classification, the result-limit wrap, the cell
conversions, the query-plan rendering, the `INFORMATION_SCHEMA` column mapping, the
not-found mapping, and the no-op dump.

## Step 1 — the project, the account and the dataset

Use a dedicated project, or a dedicated dataset in a project nobody else writes to.
BigQuery's free tier is a project tier, not a local emulator, so a trial or a
sandbox project is the whole prerequisite.

The account is a service account whose key is the recording's credential, with:

- `roles/bigquery.jobUser` on the project: run jobs, which every statement is;
- `roles/bigquery.dataViewer` on the dataset: read `INFORMATION_SCHEMA` and table
  metadata;
- `roles/bigquery.dataEditor` on the dataset: the live tests create and drop
  `omni_query_probe` and `omni_execute_probe`, and `Execute` writes.

A dataset in the corpus's location (`US` is enough), named `omni_record`. Do not
point a recording at a dataset that holds anything else: `SyncDatabase` reports
every table in it, so the golden is the dataset's whole contents.

Keep the key outside the repository. The live tests read three settings:

```bash
export BIGQUERY_CREDENTIALS=/path/to/omni-record-sa.json   # a path, or the JSON itself
export BIGQUERY_PROJECT="$(jq -r .project_id "$BIGQUERY_CREDENTIALS")"   # optional: the key names it
export BIGQUERY_DATASET=omni_record
```

`BIGQUERY_CREDENTIALS` is the gate: unset, every live test skips with its name.
`BIGQUERY_PROJECT` defaults to the key's `project_id` and `BIGQUERY_DATASET` has no
default, because a recording against a dataset nobody named is a golden nobody can
read. The credential reaches the driver as `driver.Config.Password` (the JSON
document) and the project as `Config.Params["project"]`, which is the shape
`spanner/driver` takes the same credential in: credential material crosses as a
value, and the driver hands it to the client library.

A machine whose credential is the Google authentication library's own default
chain — an attached service account read from the metadata server, a
`GOOGLE_APPLICATION_CREDENTIALS` file, the well-known file — has no key to name and
no key file to point at. There the value is `default`, the one value reserved for
it: it carries no document, so `Config.Password` stays empty and the driver leaves
the chain to the library, and it names no project, so `BIGQUERY_PROJECT` has to:

```bash
export BIGQUERY_CREDENTIALS=default
export BIGQUERY_PROJECT=bytebase-dev
export BIGQUERY_DATASET=omni_record
```

A key file cannot be called `default`. An unset variable still skips; `default` is
a credential of some kind, so a machine that named it and has none fails the test
that asked (`Ping`, the first call an open driver makes) rather than skipping it.

## The project this run used, on 2026-11-17

`bytebase-dev` is a Google Cloud project with a service account attached to the
machine that ran this, `devbox@bytebase-dev.iam.gserviceaccount.com`. It has no
dataset, so the run could not record anything; it asked everything that needs no
dataset, and its answers are below.

What the account holds, measured with `projects.testIamPermissions`:

- `bigquery.jobs.create` on the project. Every statement is a job, so every
  statement runs, and `bigquery-public-data` — which grants reads to all
  authenticated users — is readable through it.
- Nothing else: no `bigquery.datasets.create`, no `bigquery.tables.getData`, no
  `bigquery.tables.create`, no `bigquery.datasets.get`. `datasets.list` answers 200
  with no datasets, `datasets.get` answers 404 for every name in the project, and
  `bq mk --dataset` is refused.

The grant that unblocks the rest, on a dataset named `omni_record` in `US`:

- `roles/bigquery.dataEditor` on the dataset. It carries the
  `bigquery.tables.getData` the corpus's reads and `INFORMATION_SCHEMA` need and
  the `bigquery.tables.create`, `bigquery.tables.updateData` and
  `bigquery.tables.delete` the probe tables, the corpus and `Execute`'s writes
  need, so it is the one role to grant. `roles/bigquery.dataViewer` beside it is
  the read half alone, for a grantor who prefers the two apart.
- `roles/bigquery.dataOwner` on the project instead, if the recording is to create
  and drop `omni_record` itself: `bigquery.datasets.create` alone creates the
  dataset without granting anything inside it, and leaves it behind, which this
  plan forbids.

`roles/bigquery.jobUser` is already held and needs no change.

### What ran, and what it answered

| test | result |
| --- | --- |
| `TestLiveSyncInstance` | **ran**: the project lists zero datasets. The same call lists 382 from `bigquery-public-data`, in name order, and returns the service's `404 notFound` for a project that does not exist, so the empty answer is the service's and not a swallowed error. The test still fails, because it requires the project to hold the dataset the driver is open on |
| `TestLiveQuery` | **ran to its first write**: `SELECT 1`, a pipe query's rows, a pipe query under `Limit`, the CTE wrap of two statements, `MaxResultBytes` truncation, a query that returns no rows and the dry run all pass. It fails at `CREATE TABLE omni_query_probe`, which needs a dataset, and the dry run after it fails with it |
| `TestLiveSyncDatabaseMissing` | **ran**: a dataset the project does not have is `driver.ErrNotFound`. It failed before `fix(bigquery): report a missing dataset as not found` (see below) |
| `TestLiveSyncDatabase` | **blocked**: `omni_record` does not exist, which the sync now reports as `404 notFound` wrapped in `driver.ErrNotFound` |
| `TestLiveExecute` | **blocked**: `CREATE TABLE omni_execute_probe` needs a dataset |

Nothing was recorded: no tape, no sync golden, no pin row. Step 4's seam and
`bigquery/internal/testcontainer` do not exist yet, and a tape needs a dataset to
record.

### What the service answers, kept here until a tape holds it

A query's job statistics state **no result schema**; only a script's last result
leaves one there. `getQueryResults` has the schema either way, and the job's
`statementType` is `SELECT` for a plain query and `SCRIPT` for a script — which is
what `jobReturnsRows` reads.

A query that returns no rows reports its columns and no rows (recorded with
``SELECT name FROM `bigquery-public-data.usa_names.usa_1910_2013` WHERE 1 = 0``),
and a read that passes `MaxResultBytes` sets `Truncated` and keeps the row that
went over.

A job the service refuses arrives as `*googleapi.Error` on the result's `Err`, with
the call's own error nil, and never as `driver.ErrNotFound`:

| statement | answer |
| --- | --- |
| ``SELECT * FROM `bigquery-public-data.usa_names.no_such_table` `` | `404 Not found: Table bigquery-public-data:usa_names.no_such_table was not found in location US, notFound` |
| ``SELECT x FROM `bigquery-public-data.usa_names.usa_1910_2013` `` | `400 Unrecognized name: x at [1:8], invalidQuery` |
| `SELEKT 1` | `400 Syntax error: Unexpected identifier "SELEKT" at [1:1], invalidQuery` |
| `CREATE TABLE omni_query_probe (a INT64)` with no default dataset | `400 Table "omni_query_probe" must be qualified with a dataset (e.g. dataset.table)., invalid` |
| `INSERT INTO bytebase-dev.omni_record.customers (id) VALUES (1)` | `404 Not found: Dataset bytebase-dev:omni_record was not found in location US, notFound` |

The same answers come back from a dry run of a statement the service refuses,
which then ends the dry run's script.

A dataset that is not there answers two different things: `datasets.get` states
`404 Not found: Dataset <project>:<dataset>, notFound`, and a read of its
`INFORMATION_SCHEMA` states `403 Access Denied: Table
<project>:<dataset>.INFORMATION_SCHEMA.COLUMNS: User does not have permission to
query table ..., or perhaps it does not exist., accessDenied`. The sync therefore
asks `datasets.get` first.

A whole script is one job: `statementType` is `SCRIPT` and `numDmlAffectedRows` is
zero even when its statements inserted rows, and a script whose second statement
fails fails the job with `400 Query error: ... , invalidQuery`.

A cell of each type, from `SELECT NUMERIC '12.50' AS num, BIGNUMERIC '-7.25' AS bn,
NUMERIC '3' AS whole, DATE '2026-01-02' AS d, DATETIME '2026-01-02 03:04:05.123' AS
dt, TIME '03:04:05' AS t, TIMESTAMP '2026-01-02 03:04:05.120000 UTC' AS ts, JSON
'{"a":1}' AS j, [1,2] AS arr, STRUCT(1 AS a, 'x' AS b) AS st, ST_GEOGPOINT(1,2) AS
g, INTERVAL 1 DAY AS iv, RANGE<DATE> '[2026-01-01, 2026-02-01)' AS r`: the client
hands over a `*big.Rat`, a `civil.Date`, a `civil.DateTime`, a `civil.Time`, a
`time.Time`, a `string` for JSON and GEOGRAPHY, `[]bigquery.Value` for an ARRAY and
a RECORD, `*bigquery.IntervalValue` and `*bigquery.RangeValue`. The cells the driver
maps:

| column | kind | value |
| --- | --- | --- |
| `num`, `bn`, `whole` | `STRING` | `12.5`, `-7.25`, `3` |
| `dt` | `TIMESTAMP` | `2026-01-02 03:04:05.123 UTC`, precision 6 |
| `ts` | `TIMESTAMP_TZ` | `2026-01-02 03:04:05.120 UTC`, precision 6 |
| `d`, `t` | `STRING` | `2026-01-02`, `03:04:05` — the literals |
| `j`, `g` | `STRING` | `{"a":1}`, `POINT(1 2)` |
| `arr`, `st` | `STRING` | `[1 2]`, `[1 x]` — Go's rendering |
| `iv`, `r` | `STRING` | `0-0 1 0:0:0`, `&{2026-01-01 2026-02-01}` — Go's rendering |

A plan arrives only for a query that ran as more than a constant: `SELECT 1`
states none at all, and every query here stated one stage, so a plan's
`inputStages` is still unexercised against the service.

### The defects that run found

Each is fixed, tested and committed on the branch that recorded this.

1. **A statement the analysis could not classify lost its rows.** `command` read a
   job's rows only when its statistics carried a result schema, which a plain query
   never does, so a pipe query answered zero rows with an affected count of zero.
   `SELECT 1 AS n |> WHERE n = 1` returned `n = 1` through the service and nothing
   through the driver. Fixed in `fix(bigquery): keep the rows of a statement the
   analysis cannot classify`.
2. **A NUMERIC read as a ratio.** `NUMERIC '12.50'` is a `*big.Rat` whose own
   rendering is `25/2`; a NUMERIC and a BIGNUMERIC are now the exact decimal text
   the literal states. A TIMESTAMP and a DATETIME were Go's text as well, and are
   now the contract's `TimestampTZ` and `Timestamp`. Fixed in `fix(bigquery): read a
   numeric and a timestamp as their kinds`.
3. **A missing dataset was not not-found.** The sync mapped a not-found answer from
   the columns query, but the service answers that query with `403 accessDenied` for
   a dataset that is not there. Fixed in `fix(bigquery): report a missing dataset as
   not found`.

A fourth commit, `test(bigquery): let the live gate name the default credential
chain`, is what lets a machine with no key run any of this: `BIGQUERY_CREDENTIALS`
now accepts `default`.

### What is still unverified

- `SyncDatabase`'s column mapping, its table metadata read and the sync golden: a
  dataset is needed.
- `Execute`'s `numDmlAffectedRows`, the write half of the corpus, and a script
  whose earlier statements stay committed: a dataset is needed.
- A dataset that exists and cannot be read: the sync should report the service's
  `403`, not `driver.ErrNotFound`. Unverified, because no dataset exists to try it
  on.
- A plan with more than one stage, so an `inputStages` list and the ratio fields
  the port leaves out: unverified.
- `Dump`, which writes nothing: unit-tested, and nothing about it needs a project.

## Step 2 — the corpus

Load this into the scratch dataset before recording, because the sync reads
whatever is there and a tape of an empty dataset pins nothing. The live sync test
expects `customers` with the comment and the STRING `city` column named here:

```sql
CREATE TABLE `omni_record.customers` (
  id INT64 NOT NULL,
  city STRING,
  signed_up TIMESTAMP,
  score NUMERIC,
  tags ARRAY<STRING>
);
ALTER TABLE `omni_record.customers` SET OPTIONS (description = 'a synced comment');

INSERT INTO `omni_record.customers` (id, city, signed_up, score, tags) VALUES
  (1, 'berlin', TIMESTAMP '2026-01-02 03:04:05 UTC', NUMERIC '12.50', ['a']),
  (2, 'berlin', TIMESTAMP '2026-01-03 03:04:05 UTC', NUMERIC '7.25', ['a', 'b']),
  (3, 'lisbon', TIMESTAMP '2026-01-04 03:04:05 UTC', NULL, []);

CREATE VIEW `omni_record.recent` AS SELECT id, city FROM `omni_record.customers` WHERE id > 1;

-- Ingestion-time partitioning lists _PARTITIONTIME and _PARTITIONDATE in
-- INFORMATION_SCHEMA.COLUMNS with no ordinal position, so this table is what
-- proves the sync leaves them out instead of failing on the NULL.
CREATE TABLE `omni_record.partitioned` (id INT64)
PARTITION BY DATE(_PARTITIONTIME);
```

The corpus is also where the conversions this port does *not* map get recorded:
`signed_up` (TIMESTAMP), `score` (NUMERIC), `tags` (ARRAY<STRING>) and a
`SELECT STRUCT(1 AS a, 'x' AS b)` all fall back to Go's rendering of the client's
value today. Record what the client hands over, decide in the recording PR whether
that is good enough for a viewer, and pin whichever answer holds.

## Step 3 — the pin row

BigQuery states no version, so the pin is the project and dataset the answers were
recorded from. Add `bigquery/testdata/pins.tsv` with one row:

```
# columns: role, image, digest, source tag
default	gcp:bigquery:<project-id>:<dataset-location>	sha256:<digest of `datasets.get`'s id, location and creationTime>	<yyyy-mm-dd>
```

The digest column has no image digest to hold, and `datasets.get` is the closest
thing to a server identity: a dataset recreated elsewhere changes `creationTime`
and fails the pin instead of quietly replaying. Decide the exact spelling in the
recording PR; nothing in this package depends on it today.

## Step 4 — the tape harness, which needs a seam first

`Open` builds the real client, so there is nothing to tape through. The recording
PR therefore adds what `spanner/driver`, `cosmosdb/driver` and `dynamodb/driver`
each added: a `client` interface in this package holding the calls the driver
makes, with the live implementation the one call per method it takes to reach
`*bigquery.Client`, and `bigquery/internal/testcontainer` holding
`tape.OpenRaw`'s server, the pin, and a `Live` that **fails** when
`BIGQUERY_CREDENTIALS` is unset — a recording with no project must not skip, or
`OMNI_RECORD=1` would write an empty tape and call it a pass.

The calls the seam must carry, because they are what the driver asks of the
service:

| request | what the answer holds |
| --- | --- |
| `Ping` | the `SELECT 1` job's rows |
| `ListDatasets` | dataset ids, in name order |
| `DatasetMetadata(dataset)` | the dataset's own identity: what tells a missing dataset from an empty one, because the columns query answers `403` for both |
| `Read(statement)` | a job's columns and rows, every cell with the schema field type beside it, and the job's identity |
| `Run(statement)` / `DryRun(statement)` | a job, its finished status, and the `QueryStatistics` block (`statementType`, DML row count, result schema, bytes processed, query plan) |
| `TableMetadata(table)` | the table's description, row count and size |
| `Status(job)` | the plan of the job that returned the rows |

Two things a tape cannot hold, documented rather than faked: the credential path
(the API's OAuth handshake happens in the transport below any seam, so a
recording exercises what the driver builds from `Config`, not the service's
identity answer), and a real project's latency, which stays the driver's own
measurement.

## Step 5 — record, in this order

Sync first: it is the read-only half. Then the reads, then the writes.

| order | test | what it records |
| --- | --- | --- |
| 1 | `TestSyncDatabase` | the dataset's own metadata read, the `INFORMATION_SCHEMA.COLUMNS` query (the pseudo-columns of `partitioned` left out), the dataset's table listing, each table's metadata read, and the metadata it builds |
| 2 | `TestSyncInstance` | `datasets.list` and the order the sync reports |
| 3 | `TestQuery` | statement classification, every cell kind, the CTE limit wrap, a pipe query's rows and its own `Limit` cap, the dry run, `MaxResultBytes` truncation, the plan of a finished read, a statement the service refuses |
| 4 | `TestExecute` | one job for a whole script, its DML row count, and a script whose second statement fails |
| 5 | `TestLiveSyncDatabaseMissing` | the not-found mapping, on a dataset name the project does not have |

Commands, one test at a time so a failure is attributable. The table names each
test by what it records, which is the name it takes once it replays a tape;
`live_test.go` holds `TestLiveSyncDatabase`, `TestLiveSyncInstance`,
`TestLiveQuery`, `TestLiveExecute` and `TestLiveSyncDatabaseMissing` today, so
the commands below name those and step 4's seam conversion renames them. They
name the tests that exist: a `-run` pattern that matches nothing exits 0 and
records nothing, which reads as a pass.

```bash
OMNI_RECORD=1 BIGQUERY_CREDENTIALS="$BIGQUERY_CREDENTIALS" BIGQUERY_DATASET="$BIGQUERY_DATASET" go test ./bigquery/driver/ -run '^TestLiveSyncDatabase$' -count=1 -v
OMNI_RECORD=1 BIGQUERY_CREDENTIALS="$BIGQUERY_CREDENTIALS" BIGQUERY_DATASET="$BIGQUERY_DATASET" go test ./bigquery/driver/ -run '^TestLiveSyncInstance$' -count=1 -v
OMNI_RECORD=1 BIGQUERY_CREDENTIALS="$BIGQUERY_CREDENTIALS" BIGQUERY_DATASET="$BIGQUERY_DATASET" go test ./bigquery/driver/ -run '^TestLiveQuery$' -count=1 -v
OMNI_RECORD=1 BIGQUERY_CREDENTIALS="$BIGQUERY_CREDENTIALS" BIGQUERY_DATASET="$BIGQUERY_DATASET" go test ./bigquery/driver/ -run '^TestLiveExecute$' -count=1 -v
OMNI_RECORD=1 BIGQUERY_CREDENTIALS="$BIGQUERY_CREDENTIALS" BIGQUERY_DATASET="$BIGQUERY_DATASET" go test ./bigquery/driver/ -run '^TestLiveSyncDatabaseMissing$' -count=1 -v
```

`liveConfig` derives `BIGQUERY_PROJECT` from the key's `project_id`; export it
too when the recording's project is not the one the key was made in, and export
`BIGQUERY_DATASET`, which has no default and fails the test rather than skip it.

Then record everything twice from one compiled test binary and require
byte-identical tapes. Two runs of one binary are deterministic; two `go test` runs
are byte-identical in the tapes but not in a protojson golden, so compare the sync
goldens as protos:

```bash
go test -c -o /tmp/bigquery-driver.test ./bigquery/driver/
OMNI_RECORD=1 BIGQUERY_CREDENTIALS="$BIGQUERY_CREDENTIALS" /tmp/bigquery-driver.test -test.count=1
sha256sum bigquery/driver/testdata/tapes/*.tape > /tmp/bigquery-tapes.a
OMNI_RECORD=1 BIGQUERY_CREDENTIALS="$BIGQUERY_CREDENTIALS" /tmp/bigquery-driver.test -test.count=1
sha256sum bigquery/driver/testdata/tapes/*.tape | diff /tmp/bigquery-tapes.a -
```

A tape that is not reproducible is a nondeterminism in what the test sends or in
what the driver reads. The usual causes, each already met on another engine: a
listing without an order (this package sorts both listings for exactly that
reason), a value the service makes up (a job id, a creation time, a
`creationTime` in a dataset, a plan's timing), and a `SELECT` with no `ORDER BY`
whose rows arrive in a moving order. Fix the driver or narrow the answer, never
edit a tape.

## Step 6 — the sync golden

After the tapes, commit `bigquery/driver/testdata/sync/<dataset>.json`: the
protojson of `SyncDatabase`'s output, compared with `proto.Equal` in the test
rather than byte for byte, because protojson spacing moves between builds.

## Step 7 — the dump

`Dump` writes nothing and must not be declared until a renderer exists. The cheap
path — `SyncDatabase` → `Diff` → `GenerateMigration` — is not open: `googlesql`
holds no catalog at all, and the synced metadata is too thin to render even if it
did. What a dump needs, in order:

1. the sync must carry what a `CREATE TABLE` states and it does not read today:
   partitioning (`PARTITION BY` and its column), clustering, labels, the table's
   expiration, and the table type;
2. the sync must read a view's query (`INFORMATION_SCHEMA.VIEWS.view_definition`)
   and the routines, materialized views and snapshots the dataset holds;
3. a `bigquery/catalog` that renders those back into DDL.

Bytebase's BigQuery `Dump` wrote nothing either, so nothing is lost yet, and a
caller that dumps a BigQuery database writes an empty dump exactly as the product
does today.

## Step 8 — the parity harness

The gate that a tape cannot be: Bytebase's driver and this one, on the same
project, compared cell by cell. It lives in a throwaway Bytebase worktree and is
never committed.

```bash
git -C <omni repo root> worktree add ~/omni-worktrees/bigquery-verify -b bigquery-verify origin/bigquery-prep
git -C <bytebase repo root> worktree add ~/bb-worktrees/bigquery-parity -b bigquery-parity origin/main
cd ~/bb-worktrees/bigquery-parity
go mod edit -replace github.com/bytebase/omni=$HOME/omni-worktrees/bigquery-verify
go mod tidy
```

Records happen in `~/omni-worktrees/bigquery-verify`; the parity harness is
`backend/plugin/db/bigquery/parity_omni_test.go` in `~/bb-worktrees/bigquery-parity`,
and neither worktree is committed. Remove both when the recording lands.

Add `backend/plugin/db/bigquery/parity_omni_test.go` (package `bigquery_test`,
skipping without `BIGQUERY_CREDENTIALS`) and compare, for the corpus's statements
plus the sync's objects:

| Bytebase | omni | comparison |
| --- | --- | --- |
| `(&Driver{}).Open(ctx, storepb.Engine_BIGQUERY, db.ConnectionConfig{DataSource: &storepb.DataSource{ProjectId: …, GcpCredential: …}})` | `bigquerydriver.Open(ctx, driver.Config{Database: …, Password: …, Params: {"project": …}})` | both open, neither reaches the service |
| `SyncDBSchema(ctx)` | `SyncDatabase(ctx)` | the same tables, columns, types, nullability, collations, defaults, comments, row counts and sizes; compare as sorted, since this package sorts its listing and Bytebase does not |
| `SyncInstance(ctx)` | `SyncInstance(ctx)` | the dataset names as sets |
| `QueryConn(ctx, conn, sql, queryContext)` | `Session.Query(ctx, sql, queryOptions)` | one result per statement; column names and type names; each cell's kind and value; `AffectedRows`; the rewritten `Statement` text; the error's text |
| `Execute(ctx, script, executeOptions)` | `Execute(ctx, script, executeOptions)` | total affected rows, and Bytebase's `TaskRunLog` against a recording `driver.ExecuteLogger` (one command for the whole script, its rows, and the error) |
| `Dump(ctx, out, dbSchema)` | `Dump(ctx, out, metadata)` | both write nothing, and neither returns an error |

What is expected to differ, and must not be "fixed" away without a decision:

- the dataset and table orders (this package sorts both);
- the columns of a query that returned no rows (this package reports the schema it
  read; Bytebase reported none);
- a `(SELECT …)` or a CTE-prefixed DML statement (this package classifies by the
  AST; Bytebase tested the first keyword);
- `MaxResultBytes` (this package measures every row and sets `Truncated`;
  Bytebase measured powers of two and put a message in the error field);
- no `operator_email` job label (the contract carries no operator identity);
- the dry-run branch stops at the first failing statement, where Bytebase
  continued;
- an unclosed literal is the service's error on that statement's result, where
  Bytebase's own splitter returned an error from the call.

Then mutation-check the harness, which is the only proof it can fail. Apply each
mutation in the omni worktree, run the parity test (and, once it exists, the replay
test), and require a failure; then `git checkout -- bigquery/`:

1. invert the `is_nullable` mapping in `columnMetadata`;
2. drop the `"NULL"` default a column with no `column_default` reports;
3. classify `INSERT` as returning rows in `returnsRows`;
4. replace `limitStatement`'s CTE wrap with a bare `... LIMIT n` append;
5. report the plan's last stage instead of the first, or `Executed: false`;
6. return the API's 404 unchanged instead of `driver.ErrNotFound`;
7. stop sorting the dataset listing.

## Step 9 — what a failure means

- **The service rejects a statement this package sends.** A sync query is wrong.
  Run it in the BigQuery console and fix `sync.go`; do not add a retry.
- **The service accepts, the metadata differs from Bytebase's.** Bytebase's
  mapping is what customers see today; record the difference, decide which is
  right against the service, and note it in the PR. A change in what a synced
  column or table says is a behavior change for every BigQuery instance.
- **The parity harness and the tape disagree.** The tape holds an answer recorded
  from a corpus that has moved. Re-record; never edit a tape.
- **The live test fails before it sends anything** (a credential, a permission, a
  dataset that is not there). That is the account or the environment, not the
  driver. `SyncDatabase` reporting `driver.ErrNotFound` means the dataset name is
  wrong, not that the sync is broken.

## Step 10 — after the recording lands

- `docs/engine-test-matrix.md`: the `bigquery` row says "no server" today; it
  becomes the pinned project and dataset, and the tapes, the golden and the parity
  result are named there.
- The package doc's "Unverified against a project" section is deleted, not
  softened, in the PR that records.
- `AGENTS.md`'s list of engines whose tests replay a tape gains `bigquery`, and the
  sentence naming `tape.OpenRaw`'s engines gains `bigquery/driver`.
- The conversion gap that is left is pinned against the statements in [what the
  service answers](#what-the-service-answers-kept-here-until-a-tape-holds-it): a
  DATE, TIME, JSON, GEOGRAPHY, INTERVAL, RANGE, RECORD and ARRAY are still Go's
  rendering of the client's value, and a RANGE reads as `&{2026-01-01 2026-02-01}`
  today. TIMESTAMP, DATETIME, NUMERIC and BIGNUMERIC no longer do.

## Differences from Bytebase's driver, each to confirm on a project

The port is close to `backend/plugin/db/bigquery` and diverges in fifteen places.
Every one is a decision the recording settles; the three marked *(recorded)* were
settled by the run on 2026-11-17.

| difference | why |
| --- | --- |
| the dataset and table listings are sorted by name | the service's order is not part of an answer, and a tape records whatever order arrives |
| a query that returns no rows still reports its columns | Bytebase learned the schema from the first row, so it reported none |
| `returnsRows` is `googlesql/analysis.ClassifySQL`, not a first-keyword test | a parenthesized query returns rows and Bytebase sent it down the command path; a CTE-prefixed `INSERT` does not and Bytebase sent it down the query path |
| `MaxResultBytes` is measured on every row and reports `Truncated` | Bytebase measured only on power-of-two row counts and put the cap's message in the result's error, where the contract has `Truncated` |
| no `operator_email` job label | the contract's `QueryOptions` carries no operator identity; adding one is a contract change, not a driver's guess |
| the dry run stops at the first failing statement | the contract says a failing statement ends a script; Bytebase continued to the next statement |
| an unclosed literal is the statement's error, not the call's | the splitter is `googlesql/parser.SplitFlat`, which never fails; Bytebase's hand-rolled splitter refused the whole script |
| a command whose job reports no statistics is zero rows affected, not an error | the script ran; Bytebase's `QueryConn` returned an error there, but its `Execute` did not |
| `Config.TLS` and `Config.Dial` do nothing | the client library owns its transport and holds no dialer; a caller that must tunnel forwards a local port and passes it as `Config.Host`/`Config.Port`, as hive's driver asks |
| a synced column with no ordinal position is left out | an ingestion-time partitioned table's `_PARTITIONTIME` and `_PARTITIONDATE` are pseudo-columns the view states without a position, and the client cannot decode that NULL into a position |
| the dataset the sync reads is quoted as a path, and the metadata reads are sorted | a dataset may be named a reserved keyword, which a path has to quote; the service's order is not part of an answer |
| a statement `googlesql/analysis` cannot read keeps the rows the service answered *(recorded)* | omni's parser does not read a pipe query yet, and Bytebase's classification of the same statement by its first word is not available here. A plain query's job statistics state no result schema at all — only a script's last result does — so the job's `statementType` is what says whether the service returned rows |
| a missing dataset is `driver.ErrNotFound`, from a `datasets.get` the sync asks for first *(recorded)* | the service answers a read of a missing dataset's `INFORMATION_SCHEMA` with `403 accessDenied`, so the mapping could never fire from the columns query alone; a caller that drops a database on `ErrNotFound` never dropped one |
| a NUMERIC and a BIGNUMERIC cell is its exact decimal text, and a TIMESTAMP and a DATETIME are the contract's `TimestampTZ` and `Timestamp` *(recorded)* | the client renders a `*big.Rat` as a ratio (`25/2` for `NUMERIC '12.50'`) and a timestamp as Go's text (`2026-01-02 03:04:05 +0000 UTC`), which a viewer cannot read as either |

## What would unblock it

A dataset: `omni_record` in `bytebase-dev`, `US`, with `roles/bigquery.dataEditor`
for `devbox@bytebase-dev.iam.gserviceaccount.com` — the exact permissions and the
role that would do it instead are in [the project this run
used](#the-project-this-run-used-on-2026-11-17). The four recordings and the
harness above take about a day; everything else in this package is written and
gated. Without it the alternative is to drop BigQuery from the first wave and
leave `backend/plugin/db/bigquery` where it is — it has run against real
projects, which is more than this package can claim for its sync and its writes.
