# Verifying snowflake/driver

`snowflake/driver` is a port, not a driver anybody has seen work. This file is
what turns an account into a day of recording instead of a week of porting: it
lists what is missing, the account to pin, the tests to record in order, the
parity harness that holds this package to Bytebase's, and what to do when one of
them disagrees.

Nothing in this plan runs in CI. Snowflake has no image omni can start: the
service is proprietary, a trial instance dies after its credit runs out, and no
stand-in answers `ACCOUNT_USAGE`, `SHOW SHARES`, `SHOW STREAMS`, `GET_DDL` or
`CURRENT_VERSION()` — the five statements this package's sync and dump are made
of. That is why the parser is judged by the docs and why this package's tests
skip.

## What is here, and what is missing

| capability | state |
| --- | --- |
| `SyncInstance` (version, databases, users) | ported from `backend/plugin/db/snowflake`, needs an account |
| `SyncDatabase` (schemas, tables, views, streams, tasks) | ported, needs an account |
| `Session.Query` (decode, limit, explain, `MaxResultBytes`, `USE SCHEMA`) | ported, needs an account |
| `Execute` (one transaction, autocommit) | ported, needs an account |
| `Dump` (the server's own `GET_DDL`) | ported, needs an account |
| `CountAffectedRows` | **not implemented**: Snowflake's `EXPLAIN` states operations and no cardinality, and Bytebase's driver never answered one |
| tapes (`testdata/tapes`) | **none**: nothing has been recorded |
| sync goldens (`testdata/sync`) | **none**: same reason |
| `driver` row in `snowflake/testdata/pins.tsv` | **none**: no account to pin yet |
| live tests (`live_test.go`) | written, never run; skip without `SNOWFLAKE_DSN` |

The unit tests that do run cover the client configuration built from a
`driver.Config`, the private key decoder, statement classification,
the value conversions, the dump's text transforms, and the pure helpers
(the system-schema predicate, `YES`/`NO`, the stream/task vocabularies,
`SHOW TASKS`'s predecessor array).

## Step 1 — the account

Pin one dedicated, scratch-only account. The tests create tables and a stream
and a task, `Execute` writes, and the dump reads DDL, so it must be an account
nobody else uses. Either form works:

- a Snowflake trial (30 days, a credit allowance, Enterprise-equivalent, so
  `ACCOUNT_USAGE`, streams, tasks and shares are all present);
- a dedicated account in an existing organization, with a database of its own.

Record what it answers, because that is what the pin is:

```bash
snowsql -a <account> -u OMNI_RECORD -d OMNI_RECORD -q 'SELECT CURRENT_VERSION(), CURRENT_REGION();'
```

`CURRENT_VERSION()` is what `SyncInstance` reports, and it is the pin's identity:
write the release into the pin row in Step 3.

The account must provide, and a recording fails without them:

**A login.** `OMNI_RECORD`, with a password (the live tests read a DSN, so a
key-pair login is unit-tested only) and a default warehouse and role. MFA must be
off for this user, or the login fails before any statement. A network policy must
admit the recording host.

**Privileges.** A role `OMNI_ROLE` with:

```sql
CREATE WAREHOUSE OMNI_WH WAREHOUSE_SIZE = XSMALL AUTO_SUSPEND = 60;
CREATE DATABASE OMNI_RECORD;
CREATE ROLE OMNI_ROLE;
GRANT USAGE ON WAREHOUSE OMNI_WH TO ROLE OMNI_ROLE;
GRANT OWNERSHIP ON DATABASE OMNI_RECORD TO ROLE OMNI_ROLE;
-- the instance sync reads both these views, and a role without this grant
-- fails SyncInstance the way it failed Bytebase's
GRANT IMPORTED PRIVILEGES ON DATABASE SNOWFLAKE TO ROLE OMNI_ROLE;
GRANT ROLE OMNI_ROLE TO USER OMNI_RECORD;
ALTER USER OMNI_RECORD SET DEFAULT_ROLE = OMNI_ROLE DEFAULT_WAREHOUSE = OMNI_WH;
```

`GRANT IMPORTED PRIVILEGES ON DATABASE SNOWFLAKE` is what makes
`SNOWFLAKE.ACCOUNT_USAGE.USERS` and `GRANTS_TO_USERS` readable. That view lags
the account by up to an hour, so a user created just now may not be listed, and
a recording taken minutes after provisioning pins a shorter list than the account
holds. Wait an hour after creating the login, or expect the difference and say so
in the PR. `SHOW SHARES` needs no grant on a trial with no shares; an account
with a share that the role cannot see fails `getDatabaseNames`.

**Session parameters left alone.** `TIMEZONE` decides how a `TIMESTAMP_LTZ` cell
renders, and the driver records that rendering. Leave the account's parameter at
its default for the whole recording, and write the release's default into the PR
so a replay reads as the same instant.

Two environment variables reach the account. Use the same one for both:

```bash
export SNOWFLAKE_DSN='OMNI_RECORD:PASSWORD@<account>/OMNI_RECORD?warehouse=OMNI_WH&role=OMNI_ROLE'
export SNOWFLAKE_COMPAT_DSN="$SNOWFLAKE_DSN"
```

`SNOWFLAKE_DSN` is what `live_test.go` reads; `SNOWFLAKE_COMPAT_DSN` is the
variable omni's compatibility work already names for this engine. Pointing both
at the same account is the point: the statements the driver sends, and the
statements the docs' grammar accepts, are then judged by the same server.

Run the live tests first. They are the cheapest way to learn that the account is
not what Step 1 promised:

```bash
go test ./snowflake/driver/ -run '^TestLive' -count=1 -v
```

Each of the five must pass before anything is recorded. A failure here is the
account, not the driver: `SyncInstance` failing names the `ACCOUNT_USAGE` grant,
`Dump` failing names the database, and `Ping` failing is the login, the network
policy or MFA.

## Step 2 — the fixtures

Load one corpus into the scratch database before recording, because the sync
reads whatever is there and a tape of an empty database pins nothing:

```sql
CREATE SCHEMA APP;

CREATE TABLE PUBLIC.LOOKUP (
  ID NUMBER(38,0) NOT NULL,
  LABEL VARCHAR(20) NOT NULL DEFAULT 'x',
  constraint LOOKUP_PK primary key (ID)
);
COMMENT ON TABLE PUBLIC.LOOKUP IS 'a synced comment';
COMMENT ON COLUMN PUBLIC.LOOKUP.LABEL IS 'a column comment';

CREATE TABLE APP.FACTS (
  ID NUMBER(38,0) NOT NULL,
  LOOKUP_ID NUMBER(38,0),
  AMOUNT NUMBER(12,2),
  FLAG BOOLEAN,
  SEEN TIMESTAMP_NTZ(3),
  NOTES VARCHAR(80),
  constraint FACTS_PK primary key (ID),
  constraint FACTS_FK foreign key (LOOKUP_ID) references PUBLIC.LOOKUP (ID)
);

CREATE VIEW APP.RECENT AS SELECT ID, AMOUNT FROM APP.FACTS WHERE SEEN > '2026-01-01';

CREATE STREAM APP.FACTS_STREAM ON TABLE APP.FACTS;
CREATE TASK APP.FACTS_TASK WAREHOUSE = OMNI_WH SCHEDULE = '60 MINUTE'
  AS INSERT INTO APP.FACTS (ID) SELECT 1;
```

The stream and the task are what make the two `GET_DDL` reads in the sync
reachable, and the task's `SHOW TASKS` row is the only place a nullable
`warehouse`/`schedule`/`condition` column and the JSON `predecessors` array are
read. A `MODE = APPEND_ONLY` stream and a task with a predecessor each add a row
to the same tape; add them if the account allows, and say so if it does not.

The inbound-share branch of `getDatabaseNames` is the one path this corpus cannot
reach: it needs a second account that shares a database in, which `SHOW SHARES`
then reports and the sync leaves out. Record it last, on its own, and skip it if
provisioning one is unattractive — say so in the PR that records, rather than
leaving a tape that pretends.

## Step 3 — the pin row

`snowflake/testdata/pins.tsv` holds `role, image, digest, source tag`. The
grammar keeps role `default`; add a second row for the account:

```
# columns: role, image, digest, source tag
driver	snowflake:9.12.1	sha256:<SHA-256 of CURRENT_VERSION()'s answer>	2026-11-17
```

There is no image and no digest to hold, so the digest column holds the SHA-256 of
the recorded `SELECT CURRENT_VERSION()` answer. That is the same idea as the
grammar's `docs` row: the pin is whatever the answers were recorded from, and an
account of another release fails the pin instead of quietly replaying. Decide the
spelling in the recording PR; nothing here depends on the format today.

The pin row says nothing about the account name, the region, the login or the
warehouse. Those stay in `SNOWFLAKE_DSN`: they are where the account is, not what
it answered.

## Step 4 — the tape harness

Tapes need somewhere to record through. Add `snowflake/internal/testcontainer`,
owned by this engine (omni keeps no helper shared across engines), modeled on
`oracle/internal/testcontainer` and `pg/internal/testcontainer`:

- `container.go`:
  `var server = &tape.Server{Pins: Pins, Live: live, Codec: Codec{}}`;
  `func Open(t testing.TB) *sql.DB { return tape.Open(t, server, RoleDriver) }`;
  `Pins` is `snowflake/testdata/pins.tsv` resolved from `runtime.Caller`, and
  `RoleDriver = "driver"` is the row Step 3 adds.
  `live` reads `SNOWFLAKE_DSN`, opens `sql.OpenDB(gosnowflake.NewConnector(...))`
  from the same `driver.Config` the package's `liveConfig` builds, and **fails**
  when the variable is unset — a recording without an account must not skip, or
  `OMNI_RECORD=1` would write an empty tape and call it a pass.
  `Codec` encodes and decodes `*gosnowflake.SnowflakeError`: its number, SQL
  state, message and query ID, so a replayed test's `errors.As` finds the same
  type with the same fields. The query ID is a value the account makes up, so the
  codec writes a stand-in for it, as DynamoDB's codec does for its `RequestID`.
- `main_test.go`: `func TestMain(m *testing.M) { os.Exit(tape.Main(m)) }`.
- `container_test.go`: a codec round-trip of a `SnowflakeError` with and without
  a SQL state, and a test that `Open` caps the pool at one connection (the tape
  holds one session).

Then turn `live_test.go`'s scenarios into recorded tests that call
`testcontainer.Open(t)` and build their driver with
`&Driver{db: testcontainer.Open(t), databaseName: ...}`, exactly as
`pg/driver/query_test.go` opens its driver. The login is not on a tape: the tape
records `database/sql` calls, so the client's REST login happens only when a
recording opens the connection. A replay never authenticates, which is why a tape
can be committed at all.

## Step 5 — record, in this order

Sync first: it is the read-only half and the half Bytebase never tested. Then the
reads, then the writes, then the dump.

| order | test | what it records |
| --- | --- | --- |
| 1 | `TestSyncDatabase` | `SHOW SHARES`, `INFORMATION_SCHEMA.DATABASES`, `SCHEMATA`, `COLUMNS`, `TABLES`, `VIEWS`, `SHOW STREAMS`, `SHOW TASKS` and the two `GET_DDL` reads |
| 2 | `TestSyncInstance` | `CURRENT_VERSION()`, `ACCOUNT_USAGE.USERS`, `GRANTS_TO_USERS`, the database list, and the excluded system schemas |
| 3 | `TestQuery` | statement classification (SELECT/SHOW/DESCRIBE/EXPLAIN vs commands), every value kind, the LIMIT rewrite and its wrapper, the EXPLAIN rows and `Plan`, `MaxResultBytes` truncation, `QueryOptions.Schema`'s `USE SCHEMA`, error halting |
| 4 | `TestExecute` | one transaction with its BEGIN/COMMIT, the rollback of a failing script, autocommit leaving earlier statements committed |
| 5 | `TestDump` | `GET_DDL('DATABASE', ...)` and the four transforms |
| 6 | `TestSyncDatabaseShare` | the inbound-share branch, on a database shared in from a second account |

Recording commands, one test at a time so a failure is attributable:

```bash
OMNI_RECORD=1 SNOWFLAKE_DSN="$SNOWFLAKE_DSN" go test ./snowflake/driver/ -run '^TestSyncDatabase$' -count=1 -v
OMNI_RECORD=1 SNOWFLAKE_DSN="$SNOWFLAKE_DSN" go test ./snowflake/driver/ -run '^TestSyncInstance$' -count=1 -v
OMNI_RECORD=1 SNOWFLAKE_DSN="$SNOWFLAKE_DSN" go test ./snowflake/driver/ -run '^TestQuery$' -count=1 -v
OMNI_RECORD=1 SNOWFLAKE_DSN="$SNOWFLAKE_DSN" go test ./snowflake/driver/ -run '^TestExecute$' -count=1 -v
OMNI_RECORD=1 SNOWFLAKE_DSN="$SNOWFLAKE_DSN" go test ./snowflake/driver/ -run '^TestDump$' -count=1 -v
```

Then record everything twice and require byte-identical tapes. Two runs of one
compiled test binary are deterministic; two `go test` runs are byte-identical in
the tapes but not in a protojson golden, so compare the sync goldens as protos:

```bash
go test -c -o /tmp/snowflake-driver.test ./snowflake/driver/
OMNI_RECORD=1 SNOWFLAKE_DSN="$SNOWFLAKE_DSN" /tmp/snowflake-driver.test -test.count=1
sha256sum snowflake/driver/testdata/tapes/*.tape > /tmp/snowflake-tapes.a
OMNI_RECORD=1 SNOWFLAKE_DSN="$SNOWFLAKE_DSN" /tmp/snowflake-driver.test -test.count=1
sha256sum snowflake/driver/testdata/tapes/*.tape | diff /tmp/snowflake-tapes.a -
```

A tape that is not reproducible is a nondeterminism in what the test sends. The
usual causes, all met before on other engines: a query without `ORDER BY` (the
views answer in whatever order they like), a map iterated into a query, and a
value the account makes up (a query ID, a created-on timestamp, a task ID, a
version string). Two of those are already handled in the driver — the database
list and the grants are ordered — and the stream/task lists are sorted before
they are returned. Fix the driver, never the tape.

Snowflake-specific hazards:

- `SHOW SHARES`, `SHOW STREAMS` and `SHOW TASKS` answer a fixed column order that
  the client does not order; the driver reads them by name, so a tape holds the
  order the account answered and a replay reads the same columns.
- A `TIMESTAMP_LTZ` cell is rendered in the session's `TIMEZONE`, so the tape
  holds that rendering. Do not change the account's parameter between the two
  recordings.
- `GET_DDL` states the account's own normalization: a comment, a `NUMBER(38,0)`
  and a `constraint` line come back as the server writes them, not as the corpus
  wrote them. That is what makes the dump worth recording and what makes its
  golden a *server* answer rather than omni's opinion.

## Step 6 — the sync goldens

After the tapes, commit `snowflake/driver/testdata/sync/<database>.json`: the
protojson of `SyncDatabase`'s output, compared with `proto.Equal` in the test
rather than byte-for-byte, because protojson spacing moves between builds. One
golden per database kind: the corpus database, and the shared one if Step 2
produced it. `SyncInstance`'s golden is beside it, holding the version, the
database names in sorted order and the users with their roles.

## Step 7 — the dump

`Dump` is a server-DDL dump: it writes `GET_DDL('DATABASE', ...)` with four
passages taken out — the `CREATE DATABASE` it opens with, the default `PUBLIC`
schema's `CREATE`, the database qualifier on a `CREATE` line, and `or replace`.
`omni` has no `snowflake/catalog` and no metadata→DDL renderer, and this dump
needs neither.

Two things the recording settles, both visible in the golden:

1. **The shape of the answer.** The transforms assume `GET_DDL` opens with the
   `CREATE DATABASE` statement, writes a blank line, then the `PUBLIC` schema,
   and states every object's `CREATE` on a line of its own. If the release writes
   it differently, the golden holds `create or replace` or a database qualifier
   and the transforms are wrong; fix them against the recorded text.
2. **The qualifier inside a statement's body.** Only a `CREATE` statement's own
   head is unqualified, so a foreign key's `REFERENCES OMNI_RECORD.PUBLIC.LOOKUP`
   keeps the database name, as it did in Bytebase's driver. Whether such a script
   restores into the same database is what the recorded text shows; if it does
   not, the transform needs to unqualify a body line too, and that is a change to
   make against the recording rather than a guess.

`TestDump` records the answer; `TestDumpRestores` (dump, execute into a second
database, dump again, require the same text) is worth adding on the day, because
it is the only test that says the script is executable.

## Step 8 — the parity harness

The gate that a tape cannot be: Bytebase's driver and this one, on the same
account, compared cell by cell. It lives in a throwaway Bytebase worktree and is
never committed.

Start a worktree of the omni branch this file landed on, and a scratch Bytebase
worktree beside it (never the shared checkout):

```bash
git -C <omni repo root> worktree add ~/omni-worktrees/snowflake-verify -b snowflake-verify origin/snowflake-prep
git -C <bytebase repo root> worktree add ~/bb-worktrees/snowflake-parity -b snowflake-parity origin/main
cd ~/bb-worktrees/snowflake-parity
go mod edit -replace github.com/bytebase/omni=$HOME/omni-worktrees/snowflake-verify
go mod tidy
```

Records happen in `~/omni-worktrees/snowflake-verify`; the parity harness is
`backend/plugin/db/snowflake/parity_omni_test.go` in
`~/bb-worktrees/snowflake-parity`, and neither worktree is committed. Remove both
when the recording lands.

Add `backend/plugin/db/snowflake/parity_omni_test.go` (package `snowflake_test`,
skipping without `SNOWFLAKE_DSN`) and compare, for a corpus that repeats the unit
tests' statements plus the sync's objects:

| Bytebase | omni | comparison |
| --- | --- | --- |
| `(&Driver{}).Open(ctx, storepb.Engine_SNOWFLAKE, db.ConnectionConfig{...})` | `snowflakedriver.Open(ctx, driver.Config{...})` | both open, neither reaches the account |
| `SyncDBSchema(ctx)` | `SyncDatabase(ctx)` | same `*metadatapb.DatabaseSchemaMetadata`, `proto.Equal` |
| `SyncInstance(ctx)` | `SyncInstance(ctx)` | version text, the database names **as a set**, role names, and each role's attribute text |
| `QueryConn(ctx, conn, sql, queryContext)` | `Session.Query(ctx, sql, queryOptions)` | one result per statement; column names and `TypeName`s; each cell's **text**; `AffectedRows`; the rewritten `Statement` text; the error's text and number |
| `Execute(ctx, script, executeOptions)` | `Execute(ctx, script, executeOptions)` | total affected rows, and Bytebase's `TaskRunLog` against a recording `driver.ExecuteLogger` |
| `Dump(ctx, w, schema)` | `Dump(ctx, w, schema)` | the two texts, byte for byte |

What is expected to differ, and must not be "fixed" away:

- **Cell kinds.** Bytebase's Snowflake path scans every column into a string
  because its scan table names MySQL's and PostgreSQL's type names (`INT`,
  `BOOL`, `TIMESTAMP`), which gosnowflake never reports (`FIXED`, `BOOLEAN`,
  `TIMESTAMP_NTZ`). This driver reads the names the client reports, so an integer
  is an `Int64`, a boolean a `Bool`, a date a calendar day, a timestamp a
  `Timestamp`/`TimestampTZ` and a `BINARY` bytes. The harness compares each cell's
  *text* — `strconv.FormatInt`, `strconv.FormatBool`, `RFC3339Nano`, the day — and
  a difference in kind is expected. A difference in text is a defect.
- **Execute's log.** Bytebase logs the whole script as one command
  (`TaskRunLog_CommandExecute` with `Range{0, len}`); this driver logs one command
  per statement with its byte offsets. The harness compares the transaction
  control steps and the total affected rows, and compares the command list after
  collapsing Bytebase's single entry.
- **The database order.** This driver sorts the database list; Bytebase reported
  the order the view answered in. Compare as a set, and note that this is the
  difference a tape needs.
- **`CreateDatabase`.** Bytebase's create-database step runs `CREATE DATABASE`
  under the default transaction mode, which is what this driver does. Confirm on
  the day that the `COMMIT` after a `CREATE DATABASE` is accepted (Snowflake
  commits DDL implicitly); if the account refuses it, `Execute` needs the
  creation step outside a transaction, and that is a fix with a recording behind
  it.

Then mutation-check the harness, which is the only proof it can fail. Apply each
mutation in the omni worktree, run the parity test, and require a failure; then
`git checkout -- snowflake/`:

1. drop `keepCreate`'s `or replace` replacement, so the dump restores with
   `create or replace`;
2. read a `FIXED` cell as a `String` even when it declares no scale, so every
   integer becomes text;
3. return `false` from `returnsData` for a `SELECT`, so the statement runs as a
   command and reports no rows;
4. map a `DELTA` stream to `TYPE_UNSPECIFIED`;
5. invert `convertYesNo`, so every column's nullability flips;
6. drop `dropPublicSchema`, so the dump creates `PUBLIC` again.

Each must fail both the replay and the harness. A mutation that fails only one
means the harness is comparing something narrower than the driver does.

## Step 9 — what a failure means

- **The account rejects a statement this package sends.** A sync or dump query is
  wrong. Ask the account with `snowsql` and fix `sync.go` or `dump.go`; do not add
  a retry.
- **The account accepts, the metadata differs from Bytebase's.** Bytebase's
  mapping is what customers see today; record the difference, decide which is
  right against the account, and note it in the PR. A change in what a synced
  column, stream or task says is a behavior change for every Snowflake instance.
- **The parity harness and the tape disagree.** The tape holds an answer recorded
  from a corpus that has moved. Re-record; never edit a tape.
- **The live test fails before it sends anything** (a login, a network policy, an
  MFA challenge). That is the account, not the driver.
- **`ACCOUNT_USAGE` answers nothing.** The lag, not the driver: wait an hour and
  re-record. A tape recorded inside the lag pins a shorter user list, and the next
  recording will disagree with it.

## Step 10 — after the recording lands

- `docs/engine-test-matrix.md`: the `snowflake` row says "no server" today; it
  becomes the pinned release, and the tapes, the sync goldens and the parity
  result are named there.
- The package doc's "Unverified against an account" section is deleted, not
  softened, in the PR that records.
- `AGENTS.md`'s list of engines whose tests replay a tape gains `snowflake`, and
  the sentence that names `SNOWFLAKE_COMPAT_DSN` gains `snowflake/driver`'s live
  tests.

## Differences from Bytebase's driver, each to confirm on an account

| difference | why |
| --- | --- |
| `Execute` ships one statement per request; Bytebase sent the whole script as one multi-statement request (`MULTI_STATEMENT_COUNT=0`) | the execution log a caller reads is per statement, and the splitter is lexical and never fails; what a multi-statement request does differently with a failing statement is what the recording shows |
| a cell carries the contract's kind; Bytebase wrote every cell as text | its scan table names MySQL's and PostgreSQL's type names, which the client never reports, so the text fallback caught every column; a decimal still stays exact text |
| an empty `BINARY` value is `Bytes` with no bytes, not NULL | Bytebase's shared converter turned an empty `[]byte` into NULL, which drops a value that is there |
| `Config.Port` is honored; Bytebase read the data source's port nowhere and always dialed 443 | the contract has a port field, and a proxy endpoint may answer on another one |
| a database name is escaped when it is quoted as an identifier; Bytebase interpolated it as written | a name holding `"` is still one identifier |
| the database list is sorted by name, and a user's roles are read in order | the views answer in no order a tape can hold |
| `SHOW SHARES`, `SHOW STREAMS` and `SHOW TASKS` are read by column name, and a column the answer does not state leaves its field empty | Bytebase read fixed positions, so a server that states fewer columns filled the first field with the wrong one |
| a NULL in a sync answer is its zero value: `TABLES.ROW_COUNT`, a stream's comment, a task's schedule and its predecessor array | Bytebase scanned those columns into plain values, so the whole sync failed on `converting NULL to int64 is unsupported` for a table no scan had measured |
| a rewritten `EXPLAIN` result is marked as a plan; Bytebase marked only a statement the user wrote as `EXPLAIN` | the rows are a plan either way |
| a private key is `Params[ParamPrivateKey]` with its passphrase beside it, not a data source field | the contract carries engine-specific settings in `Params`, as `cosmosdb/driver` does for credentials |
| a caller's `Dial` or `TLS` is installed as the client's HTTP transport | gosnowflake takes a transport rather than a dialer; through one, it checks no certificate revocation, so OCSP fail-open applies only when this driver builds the transport |
| a missing database is `driver.ErrNotFound`; Bytebase wrapped its own `NotFound` error | the contract's error, which `errors.Is` reads |

## What would unblock it

An account: a Snowflake trial or a dedicated account with `SNOWFLAKE_DSN` set,
`GRANT IMPORTED PRIVILEGES ON DATABASE SNOWFLAKE` for the login, and one scratch
database. The recordings above are about a day; everything else in this package is
already written and gated. If no account is going to exist, the alternative is to
drop Snowflake from the first wave and leave `backend/plugin/db/snowflake` where
it is — it has run against real accounts, which is more than this package can
claim.
