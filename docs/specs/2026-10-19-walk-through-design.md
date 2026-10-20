# WALK_THROUGH without a catalog: design

Date: 2026-10-19. Status: proposal. First engine: PostgreSQL (`pg/review`). The core is engine-neutral, so the other `<engine>/review` packages can adopt it.

## 1. Requirements

WALK_THROUGH (`review.WalkThrough`) reports that applying the change to the target fails. It is a P0 rule in Bytebase, so a wrong finding blocks a valid change. The requirements are:

1. **No catalog.** The rule does not load the target into `pg/catalog`, or into any other simulator, and replay the change. The rule has to reach every engine with a review package. MSSQL, Oracle, and Snowflake have no catalog, and the catalogs that do exist differ in depth.
2. **No false findings.** Every reported failure must be real. Missing a failure is acceptable, and completeness is not a goal.

Bytebase's v1 rule (`bytebase/backend/plugin/schema/pg/walk_through.go`) does the opposite. It loads the synced schema into `pg/catalog`, runs `catalog.Exec`, and reports the first error. Every point where the simulator and PostgreSQL disagree becomes a wrong result. #483, for example, fixed a `DROP COLUMN` the catalog refused because of the table's dependents rather than the column's. An engine without a catalog cannot use the approach at all. The closed PR #432 took the same approach in `pg/review`.

## 2. What a finding claims

> **Statement k fails, given that statements 1..k-1 succeeded.** This holds for every database state consistent with the target's synced schema (§4.2) and for every way those earlier statements could have succeeded.

- The claim is conditional on purpose. A statement can fail for reasons the walk does not look at, such as data, types, or permissions. In that case the change stops earlier than the finding, and the finding still holds: fix the earlier statement and statement k still fails.
- **At most one WALK_THROUGH finding per target:** the first failure the walk can establish. The change stops there when it runs. A statement that failed did nothing, and nothing says what the corrected statement will do, so nothing after it can be claimed. Bytebase v1 also reports only the first failure.
- The message uses PostgreSQL's wording. It names the object as the SQL wrote it, quoted the way `pg/review`'s other rules quote names (contract: one line, no target name), for example `relation orders already exists`. The anchor is the name in the statement.

## 3. Idea: an existence ledger, not a simulator

The walk asks one question: **"right before this statement runs, does the object it names certainly exist, or certainly not exist?"** It models nothing else: no types, expressions, dependency graph, or privileges.

Every fact is three-valued, `Present | Absent | Unknown`, and **Unknown never produces a finding**. Precision comes from three rules:

1. **Initial facts come from the snapshot, with explicit authority.** Present is trusted for whatever the synced schema lists. Absent is trusted only in the spaces where the engine knows Bytebase's sync is complete (§4.2). Everywhere else it is Unknown.
2. **Default deny.** A statement form the engine has not listed ends the walk. Coverage grows statement by statement. An unlisted statement can only cost recall, never precision.
3. **Effects are explicit or conservative.** A listed statement either sets exact facts or declares the scope it may change, through `MayAdd` or `MayRemove`. A statement with no bound on its effect ends the walk.

## 4. Engine-neutral core (`review/walk`)

The core imports nothing from any engine and does not read `metadata`. Engines share `review` already, and `review/walk` sits under it. The root AGENTS.md "Repository shape" line should name it.

### 4.1 Keys and facts

```go
type Fact uint8 // Unknown, Present, Absent

// Key names one object. Space is an engine-defined namespace
// ("relation", "column", ...). Parent is the object that holds it
// (a schema; a schema and table). Name is the stored form of the name.
type Key struct{ Space, Parent, Name string }
```

- Present facts carry an engine-defined kind ("table", "view", "index"), so the walk can report a wrong kind (`"v" is not a table`). A Present fact may have kind Unknown, for example after `CREATE TABLE IF NOT EXISTS` on a name that was Unknown.
- The stored form of a name belongs to the engine (§4.4). The core compares keys only by equality.

### 4.2 Snapshot and authority

```go
// Snapshot is the synced schema as the engine reads it.
type Snapshot interface{ Lookup(Key) (Fact, kind string) }
```

The engine implements it over `metadata.DatabaseSchemaMetadata`:

- It answers Present for every object the schema lists.
- It answers Absent only for a space it declares authoritative. A space is authoritative only when the engine has read **Bytebase's sync code for that engine** and found that the sync lists every object of the space, regardless of privileges, filters, or extensions. The engine cites that code next to the declaration.
- It answers Unknown everywhere else.

The rule is asymmetric because absence is the dangerous claim. The PostgreSQL sync's real gaps are listed in §5.1.

### 4.3 Ledger and steps

```go
type Scope struct {
	Space, Parent string
	AnyParent     bool
	Prefix        string // only names with this prefix, "" for all
}

func (*Ledger) Lookup(Key) (Fact, string)
func (*Ledger) Set(Key, Fact, kind string)
func (*Ledger) MayAdd(Scope)    // Absent -> Unknown in scope; Present survives
func (*Ledger) MayRemove(Scope) // Present -> Unknown in scope; Absent survives
func (*Ledger) Move(space, from, to string) // children follow a renamed or moved parent

type Check struct {
	Key     Key
	Want    Fact     // Present or Absent
	Kinds   []string // for Present: acceptable kinds; nil = any
	Anchor  any      // engine data: range and message for the finding
}

type Step struct {
	Checks    []Check          // judged against the state before the statement
	Apply     func(*Ledger)    // effects of success
	StopAfter bool             // effect unbounded: end the walk after the checks
}

// Run lowers and judges statement by statement and returns the first check
// that is certainly violated.
func Run(n int, snap Snapshot, lower func(i int, l *Ledger) Step) (stmt int, failed Check, ok bool)
```

A check is violated only when the ledger shows the opposite fact for certain: Absent for `Want: Present`, Present for `Want: Absent`, or Present with a kind outside `Kinds`.

`MayAdd` and `MayRemove` are what let the walk keep going after statements it cannot follow exactly:

- An unnamed index adds a name the walk does not derive. That is a `MayAdd` over the names it could take.
- `DROP ... CASCADE` removes dependents the walk cannot list. That is a `MayRemove` over every space.
- A statement whose effect depends on a fact that is Unknown, such as `IF [NOT] EXISTS` on an Unknown name, applies its effects as May-effects.

### 4.4 Rules every engine follows

- **Name form is exact or uncertain, never approximate.** The engine folds and truncates names exactly as the server stores them. PostgreSQL lowercases unquoted names and truncates at 63 bytes (verified: `DROP TABLE <64 a's>b` drops the 63-byte table). Where equality depends on something the walk does not know (an MSSQL collation, MySQL `lower_case_table_names=2`), two names that differ only in case look up as Unknown.
- **Within one statement, no claim about a name the statement itself creates, drops, or renames.** Servers apply subcommands in their own order, not the order written. PostgreSQL runs `ALTER TABLE` in passes: `ADD COLUMN a, DROP COLUMN a` fails with "column does not exist", while `DROP COLUMN a, ADD COLUMN a` succeeds (verified on PG 17). Checks are judged against the state before the statement, and the engine leaves out checks on names the statement touches twice.
- **Name resolution is the engine's job**: search path, current database or schema, synonyms. Resolution answers either a key, a certain "not found", or Unknown. The first candidate schema that answers Unknown makes the whole resolution Unknown.
- **The target's server version is not an input.** Every fact the walk uses must hold on every version Bytebase supports for that engine. For PostgreSQL that is 12 to 18.
- **Event triggers and procedural code end the walk** (§5.3).

## 5. PostgreSQL v1

### 5.1 Spaces and authority

The sync facts come from `bytebase/backend/plugin/db/pg/sync.go` and `bytebase/backend/plugin/parser/pg/system_objects.go`, read on 2026-10-19.

| Space | Parent | Holds | Absent trusted? |
| --- | --- | --- | --- |
| schema | database | schemas | Yes, except Bytebase's system schemas (`pg_catalog`, `information_schema`, `pg_toast*`, `pg_temp*`, the TimescaleDB schemas, ...) and names starting with `pg_` |
| relation | schema | tables, partitioned tables, views, materialized views, sequences, foreign tables, indexes, composite types (all of `pg_class`) | Yes, in a synced schema, except the names the sync drops (its `systemTables` and `systemViews` lists and prefixes) and the names of partitions and their indexes. Those are Unknown in **every** schema, because the metadata keeps a partition under its parent and does not record the partition's own schema. |
| type | schema | enum and composite types, plus the row type of every table, view, materialized view, and foreign table | **No.** Domains, ranges, and base types are not synced. Sequences and indexes have no row type (verified). |
| column | schema.table | columns of synced tables and foreign tables (from `pg_attribute`, with no privilege filter), plus the system columns `tableoid, xmin, cmin, xmax, cmax, ctid` | Yes for those tables. No for partitions: their columns are not synced. |
| constraint | schema.table | PK, UNIQUE, FK, CHECK, EXCLUDE names | **No.** PostgreSQL 18 names NOT NULL constraints (`t_a_not_null`), and the metadata does not list them. Verified: `DROP CONSTRAINT t_a_not_null` succeeds on 18 and fails on 17. |

Assumptions about how Bytebase runs the change, which the walk relies on and which must be verified in Bytebase's PostgreSQL executor before release:

- Each task runs on a fresh session, so `pg_temp` starts empty. Unqualified relation names look in `pg_temp` first, so this matters even for "already exists" claims.
- Execution stops at the first failed statement.

### 5.2 Statements in v1

"Present" means the name resolves certainly to an object, which is checked against the ledger. A check named in the table produces a finding only when the ledger contradicts it for certain.

| Statement | Checks | Effects of success |
| --- | --- | --- |
| `CREATE SCHEMA [IF NOT EXISTS] s` (no elements) | s Absent (42P06); name not `pg_*` (42939) | s Present; every space under s Absent |
| `CREATE [TEMP] TABLE [IF NOT EXISTS] t (...)` | creation schema Present (3F000); relation t Absent (42P07); type t Absent (42710); within the statement: duplicate column (42701), two primary keys (42P16), two constraints of one name (42P07 for two keys, else 42710; NOT NULL left out, as only PostgreSQL 18 names it); a database other than the current one (0A000, for any statement that names one) | t Present(table). For a plain column list, its columns become exact and the column space is authoritative. With `LIKE`/`INHERITS`/`OF`/`PARTITION OF`, the columns are Unknown. Named constraints become Present; an index-backed one also takes its name in the relation space. Generated names (serial sequence, key indexes, unnamed constraints): `MayAdd(relation of schema, prefix = t[:16])` and `MayAdd(constraint of t)`. |
| `CREATE [UNIQUE] INDEX [CONCURRENTLY] [IF NOT EXISTS] [i] ON t` | t Present with kind table, partitioned, or matview (42P01 / 42809); i Absent in t's schema (42P07) | i Present(index), or a prefix `MayAdd` when unnamed |
| `CREATE [MATERIALIZED] VIEW v`, `CREATE SEQUENCE q`, `CREATE TYPE` (enum, composite), `CREATE DOMAIN` | creation schema Present; relation Absent (not for enum or domain); type Absent (not for sequences); none for `CREATE OR REPLACE VIEW` | Present. A view's columns stay Unknown. A view of a temporary relation lands in `pg_temp`, so once the change has made one, an unqualified view is placed nowhere. |
| `ALTER TABLE [IF EXISTS] [ONLY] t ...` | t Present (42P01) | per subcommand: |
| — `ADD COLUMN [IF NOT EXISTS] c` | c Absent (42701), including system column names (42701) | c Present; inline constraints as in CREATE TABLE |
| — `DROP COLUMN [IF EXISTS] c [CASCADE]` | c Present (42703) | c Absent; t's constraints and indexes `MayRemove`; CASCADE is a `MayRemove` over every space |
| — `ALTER COLUMN c ...`, `RENAME COLUMN a TO b` | c / a Present (42703); b Absent (42701) | rename moves the column |
| — `RENAME TO n`, `SET SCHEMA s2` | n Absent (42P07, 42710); s2 Present (3F000) | `Move` t and its children; the indexes and owned sequences that move with t become May-effects in both schemas |
| — `ADD [CONSTRAINT n] ...` | n Absent (42710; 42P07 when index-backed); PK when t certainly has one (42P16); FK target relation Present (42P01) | n Present, or May-effects when unnamed |
| — `DROP CONSTRAINT [IF EXISTS] n` | n Present, from facts **the change itself** established only (42704) | n Absent; CASCADE is a `MayRemove` over every space |
| — `SET/RESET (...)`, `OWNER TO`, `SET/DROP DEFAULT`, `ENABLE/DISABLE TRIGGER`, `REPLICA IDENTITY`, `SET [UN]LOGGED`, `VALIDATE CONSTRAINT`, `CLUSTER ON`, `SET TABLESPACE` | none beyond t | none |
| — any other subcommand (`ATTACH/DETACH PARTITION`, `INHERIT`, `OF`, ...) | | stop |
| `ALTER INDEX/SEQUENCE/VIEW/MATERIALIZED VIEW [IF EXISTS] x RENAME TO n / SET SCHEMA s` | x Present; n Absent | `Move` |
| `DROP TABLE/VIEW/MATERIALIZED VIEW/INDEX/SEQUENCE/FOREIGN TABLE/TYPE [IF EXISTS] x, ... [CASCADE]` | each x Present (42P01 / 42704, in PG's per-kind wording) unless IF EXISTS; **a different kind fails even with IF EXISTS** (42809, verified) | x and its columns, constraints, and synced indexes Absent; owned sequences (`SequenceMetadata` owner) Absent; a change-created table's generated names become a prefix `MayRemove`; CASCADE is a `MayRemove` over every space |
| `DROP SCHEMA [IF EXISTS] s [CASCADE]` | s Present (3F000); without CASCADE, s certainly holds an object (2BP01, verified) | s and everything under it Absent; CASCADE is a `MayRemove` over every space |
| `ALTER SCHEMA s RENAME TO n` | s Present; n Absent (42P06) | `Move` |
| `COMMENT ON`, `GRANT/REVOKE ... ON TABLE`, `TRUNCATE` | object Present (42P01) | none |
| `INSERT/UPDATE/DELETE/MERGE` | the target Present (42P01). A CTE name is never a target (verified). | none; see §5.3 |
| `SELECT` | none | none; see §5.3 |
| `SET/RESET search_path`, `SET ROLE`, `SET SESSION AUTHORIZATION` | none | update the resolution context (reuse `searchpath.go`) |
| `BEGIN/START/COMMIT/SAVEPOINT/RELEASE` | none | none |
| `CREATE [OR REPLACE] FUNCTION/PROCEDURE` | none | record the name as user code (§5.3) |
| `CREATE TRIGGER ... ON t` | t Present | t has triggers (§5.3) |
| everything else: `DO`, `CALL`, `ROLLBACK [TO]`, `COPY`, `CREATE EXTENSION`, `ALTER EXTENSION`, `CREATE RULE`, `CREATE EVENT TRIGGER`, `CREATE SCHEMA` with elements, `IMPORT FOREIGN SCHEMA`, `DROP OWNED`, `REASSIGN OWNED`, `CREATE OR REPLACE` of anything but a view or a routine, ... | | stop |

Name resolution, from PostgreSQL's namespace rules:

- **Unqualified references** search `pg_temp`, then `pg_catalog` unless the path lists it, then the path. A name starting with `pg_` resolves as Unknown, since a sync gap or a catalog name could match. Verified: `DROP TABLE pg_class` reaches `pg_catalog.pg_class`, and every `pg_catalog` relation name starts with `pg_`.
- **Assumed: `allow_system_table_mods` is off.** That `pg_catalog` holds no other name rests on PostgreSQL refusing to create a relation there, to a superuser and an extension script alike (verified on PostgreSQL 17). Only a superuser who sets that developer setting can. Treating the implicitly searched `pg_catalog` as Unknown instead would turn every unqualified "does not exist" into Unknown. A name the SQL qualifies with `pg_catalog` is never judged.
- **Unqualified creations** go to the first path schema that exists. When the path is known and none of its schemas exists, that is a certain 3F000.
- **`$user`** expands to the current user. An unknown user makes it Unknown.

The 16-byte prefix for generated names comes from `makeObjectName`: it truncates the longer of the two parts first, so the table part always keeps at least 25 bytes when the table name is longer than that. Verified with two 60-byte parts: the table part kept 29 bytes. Generated names never fail. A colliding one gets a number (`t_a_idx1`), which is why the walk cannot derive them.

### 5.3 Where the walk stops

- **Event triggers.** If the target has an enabled event trigger, the walk stops at the first statement one may fire on. This reuses `scan.triggerFires`.
- **Procedural code.** `DO` and `CALL` stop the walk.
- **User functions.** A statement that calls a function the snapshot or the change defines, or calls `set_config`, stops the walk after its checks.
- **DML triggers.** A DML statement or `TRUNCATE` stops the walk after its checks if any table in the snapshot or the change has a trigger. Triggers that create partitions on INSERT are a real pattern.
- **What may hold user code the sync does not show.** Bytebase syncs no triggers of partitions or foreign tables, so a database with either counts as having triggers. A rule on a table or view counts as a trigger, since its action may call a function.
- **pg_dump's header.** `SELECT pg_catalog.set_config('search_path', '', false)` sets a known path, possibly empty, as the scan already reads it. Any other `set_config` call makes the path unknown.
- **Another spec of the plan** (`OtherSpecs`) may have added a DML trigger, so DML ends the walk then (§8.1).
- **Accepted assumption.** User code reached only implicitly (column defaults, CHECK expressions, policies, user-defined operators, domain checks) is assumed not to change the schema. This is the line `pg/review`'s scan already draws. The stricter alternative is to stop after any DML whenever a user function exists, which would end the walk at the first DML in most real databases.

### 5.4 Not covered (declined by scope)

Each of these needs the change simulated, or facts the metadata does not carry. Review comments asking for them are recall asks and get declined, citing this list:

- types and expressions;
- column references inside queries, defaults, and checks;
- relation references inside queries and view bodies (a v2 candidate);
- function and overload resolution;
- privileges and ownership. Superusers and role membership make "must be owner" uncertain, and Bytebase has no role graph;
- data-dependent failures (NOT NULL, UNIQUE, casts, rewrites);
- inheritance and partition columns;
- `OR REPLACE` compatibility;
- `CREATE SCHEMA` elements;
- dependency failures (2BP01 for `DROP TABLE` with FK or view dependents, `DROP COLUMN` with view dependents). These are v2 candidates, since the FK and view-dependency metadata could support them as Present facts.

**Size budget:** the core within about 400 lines and the PostgreSQL v1 lowering within about 1,500. A review loop that pushes past that is a scope decision for a person, not a fix. See `pg/review/rule_prior_backup.go`, which grew to 2,882 lines and was cut back to about 1,175.

As first written, v1 overshoots this budget: the core is about 440 lines, and the PostgreSQL side about 2,800 (lowering about 1,800, names and resolution about 320, snapshot about 290, walker about 430). Most of the excess is breadth of the §5.2 table rather than depth: each form is short. Trimming forms costs only recall, so the budget now stands at the delivered size, and growth beyond it is a person's decision.

Two inferences keep the walk useful where the snapshot is silent:

- A creation in a schema the SQL names implies, once the statement succeeds, that the schema exists.
- A relation the change itself made resolves even when its schema's existence is otherwise unknown.

The second matters under `OtherSpecs` and for schemas the sync skips. A creation the search path places stays unplaced while the path's schemas are unknown.

## 6. Verification: precision is tested, not argued

1. **Unit tests in `review/walk`** for the ledger semantics: `MayAdd` keeps Present, `MayRemove` keeps Absent, `Move` carries children, and check verdicts on every fact and kind combination.
2. **A curated PostgreSQL oracle** (`TestWalkThroughAgainstPostgres`, on `postgres:17-alpine`):
   - Each case loads a base schema and syncs it with a test sync that **mirrors Bytebase's exclusions** for every space the walk trusts as Absent. The existing `syncSchema` helper is more complete than Bytebase's sync, so it would hide authority bugs.
   - The harness runs the change one statement at a time, each in a savepoint. It keeps the statements PostgreSQL accepts and skips the ones it rejects.
   - For each rejected statement r, the case is the accepted prefix plus r. The walk must report nothing, or report r. A report on r must be true: either PostgreSQL rejects r with the SQLSTATE the message stands for, or the fact the finding claims holds in the server's catalog as it was right before r. The second covers a statement that fails for two reasons, where PostgreSQL checks the other one first.
   - The walk over the whole accepted sequence must report nothing.
   - Cases marked `complete` must be found.
3. **A differential sweep** (`TestWalkThroughSweep`, size `WALK_SWEEP_N`, larger in CI as for the other sweeps):
   - It generates random DDL over a small name pool: two schemas, tables `a/b/c`, names that collide with generated ones (`a_id_seq`, `a_pkey`), columns, `IF [NOT] EXISTS`, serial and keys, unnamed indexes, renames, `SET SCHEMA`, `CASCADE`, `SET search_path`, `TEMP`, and types.
   - It applies the same oracle as the curated test. Precision must be 100%. Recall is logged, not asserted.
   - This is what earns "never wrong" and keeps it earned.
4. **Version-sensitive facts.** The curated cases also run on `postgres:18-alpine` (`TestWalkThroughAgainstPostgres18`), whose NOT NULL constraints have names, and `WALK_SWEEP_IMAGE` points the sweep at any image. Facts for older versions are spot-checked when they are added.

## 7. Other engines

An engine adopts WALK_THROUGH by supplying five things. The core and the oracle method stay the same.

1. **Spaces and kinds**, following the engine's namespace sharing:
   - MySQL: tables and views share a space per database. Index names are per table. FK and CHECK constraint names are per database.
   - MSSQL: `sys.objects` names (tables, views, procedures, functions, constraints, sequences, synonyms) share a space per schema.
   - Oracle: tables, views, sequences, synonyms, procedures, packages, and types share a space per schema. Indexes and constraints each have their own space.
   - Snowflake: tables and views share a space per schema.
2. **Authority per space**, read from that engine's Bytebase sync. Known questions to answer first:
   - Is the MySQL sync's `information_schema` privilege-filtered?
   - Unqualified Oracle names can resolve to a public synonym, which the sync may not list.
3. **Name form:**
   - MySQL: `Target.LowerCaseTableNames`; column and index names are always case-insensitive.
   - MSSQL: the database collation.
   - Oracle and Snowflake: uppercase folding of unquoted names.
4. **Resolution:** the current database or schema and `USE`, the MSSQL default schema, Oracle synonyms.
5. **The statement allowlist and stop list**, each row pinned by an oracle case against the engine (testcontainers for mysql, mariadb, tidb, mssql; the DSN-gated harness for Oracle).

## 8. Decisions

### 8.1 Other changes to the same database

Bytebase reviews each spec against the same synced schema (`reviewrun.groupReviewUnits`). If one spec creates `t` and another alters it, the walk would report `relation "t" does not exist` for the second. That would be a false finding. Whether the walk can account for the other spec depends on whether Bytebase orders the two. The Bytebase code (main, 2026-10-19) shows that the UI and GitOps flows differ.

**UI workflow: no order.**
- `GetPipelineCreate` creates tasks in spec order, but nothing enforces that order when they run.
- `isSequentialTask` (`backend/runner/taskrun/running_scheduler.go`) returns true only for a release task. Its comment says: "Non-release migrations can run concurrently on the same database."
- So two specs of a plan on the same database may run in either order, or at the same time.

**GitOps (release): ordered.**
- A release task runs the release's unapplied files one after another, in the release's file order. Applied versions are skipped (`runVersionedRelease`), and release tasks on one database exclude each other.
- Bytebase v1's `CheckRelease` (`backend/api/v1/release_service_check.go`) already follows that order. It clones the schema once per target and passes the same `FinalMetadata` to every file, and the v1 walk-through writes each file's result back into it. So file k is checked against the schema after files 1..k-1.
- Note: `deriveTargets` skips release specs, so the V2 review run sees only UI specs today. Releases are checked only by `CheckRelease`.

Decision:
- **UI specs (v1):** add a flag the caller sets when another spec of the same plan changes the same database, e.g. `Target.OtherSpecs bool`. The caller has every spec of the plan in `deriveTargets`, so it can compute this per target. With the flag set, the walk runs on an all-Unknown snapshot, and only the change's own facts count (the same table created twice; a drop followed by a use). Replay is impossible, because there is no order to replay in. The other spec may also have added a DML trigger, so DML ends the walk under the flag. That it added no event trigger is an accepted assumption, the same line as user code reached only implicitly (§5.3): stopping before every DDL statement would leave such a change nothing to report.
- **GitOps, when V2 takes over release checks:** pass the release's earlier unapplied files in order, e.g. `Target.Before []string`. The walk replays them without reporting and then walks the file. This matches what v1 does with `FinalMetadata`. If a replayed file reaches a stop statement, the walk falls back to an all-Unknown snapshot.
- **Possible refinement for UI specs:** pass the other specs' SQL and apply their effects as May-effects in any order: every key they may create becomes `MayAdd`, and every key they may drop becomes `MayRemove`. That keeps the snapshot facts for names they do not touch. This is not needed for v1.

### 8.2 Implicit user code

Decided: accept the assumption in §5.3.

### 8.3 Where the core lives

Decided: `review/walk` from the first PR.

Related contract change, small and optional: with `Target.Schema == nil`, run the walk on an all-Unknown snapshot instead of skipping it, so the change's own conflicts are still reported.

## 9. Delivery

One PR, since the ledger is only exercised through its first engine:

1. **Commit `feat(review): add the walk ledger and Target.OtherSpecs`**: `review/walk` with its unit tests, and the contract field.
2. **Commit `feat(pg/review): add the WALK_THROUGH rule`**:
   - snapshot and authority (`walk_snapshot.go`);
   - names and resolution (`walk_names.go`, reusing `searchpath.go` and `names.go`);
   - the lowering of §5.2 (`walk_lower.go`);
   - wiring through `targets.go` (deduplication, workers, merging);
   - the curated oracle and the sweep (`WALK_SWEEP_N=500` in CI).
3. Bytebase: set `Target.OtherSpecs` in `reviewrun`, register `pgreview.Review` for WALK_THROUGH and bump the pin. This waits for the user.
4. v2 candidates (§5.4), then engines in the order of the multi-engine plan.

## Appendix: engine facts this design rests on

Probed on `postgres:17-alpine` (and `postgres:18-alpine` where noted) on 2026-10-19.

| Statement | Result |
| --- | --- |
| `ALTER TABLE t ADD COLUMN a int, DROP COLUMN a` | ERROR column "a" of relation "t" does not exist |
| `ALTER TABLE t DROP COLUMN a, ADD COLUMN a text` | ok |
| `DROP TABLE <63 a's>b` with table `<63 a's>` | ok: NOTICE truncated, table dropped |
| `CREATE TABLE mood` with enum `mood` | ERROR type "mood" already exists |
| `CREATE TYPE t` with table `t` | ERROR type "t" already exists |
| `CREATE TYPE q` with sequence `q`; `CREATE TYPE ix` with index `ix` | ok, both |
| `CREATE TABLE ix` with index `ix` | ERROR relation "ix" already exists |
| `DROP TABLE v`, `DROP TABLE IF EXISTS v` with view `v`; same on index `ix` | ERROR "v" is not a table (all four) |
| `CREATE TABLE IF NOT EXISTS v`, `CREATE INDEX IF NOT EXISTS u` on existing relations | NOTICE, skipped |
| `ADD COLUMN xmin` / `DROP COLUMN ctid` / `ADD COLUMN oid` | ERROR conflicts with a system column / ERROR cannot drop system column / ok |
| `WITH x AS (SELECT 1) UPDATE x ...` | ERROR relation "x" does not exist |
| `DROP SCHEMA s` holding a table | ERROR cannot drop schema s because other objects depend on it |
| `CREATE SCHEMA pg_x` | ERROR unacceptable schema name "pg_x" |
| `pg_catalog` relations not named `pg_*` | 0 |
| `DROP TABLE pg_class` with `public.pg_class` | ERROR permission denied: "pg_class" is a system catalog |
| `CREATE TEMP TABLE t` over `public.t`, then `ALTER TABLE t ADD COLUMN id` | ok: the temp table shadows |
| two primary keys, in CREATE TABLE or by `ADD PRIMARY KEY` | ERROR multiple primary keys for table "t" are not allowed |
| `ADD CONSTRAINT c CHECK` / `UNIQUE` with CHECK `c`; `ADD CONSTRAINT c UNIQUE` with table `c` | ERROR constraint "c" for relation "t" already exists (both) / ERROR relation "c" already exists |
| generated names, 60-byte table and column | `<29>_<29>_idx`, `<56>_id_seq`, `<58>_pkey`, ... |
| unnamed index whose generated name is taken | ok, named `t_a_idx1` |
| PG 18: `DROP CONSTRAINT t_a_not_null` | ok (PG 17: ERROR does not exist) |
| PG 18: `CREATE TABLE t (a int CONSTRAINT foo NOT NULL); ALTER TABLE t DROP CONSTRAINT foo` | ok (PG 17: ERROR does not exist) |
| PG 18: `CONSTRAINT c NOT NULL` and `CONSTRAINT c CHECK` in one CREATE TABLE | ERROR 23505 (PG 17: ok) |
| two CHECKs named `c` / two UNIQUEs, or UNIQUE and PRIMARY KEY / CHECK and UNIQUE, FK and CHECK | ERROR 42710 check constraint "c" already exists / 42P07 relation "c" already exists / 42710 constraint "c" for relation "x" already exists |
| `CREATE TABLE pg_catalog.foo` as a superuser | ERROR permission denied, system catalog modifications are currently disallowed; ok after `SET allow_system_table_mods = on` |
| `CREATE TABLE other.public.x`, `ALTER TABLE other.public.x`, `DROP TABLE IF EXISTS other.public.x` | ERROR 0A000 cross-database references are not implemented (the current database's name is accepted) |
