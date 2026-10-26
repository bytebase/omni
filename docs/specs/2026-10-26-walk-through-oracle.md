# WALK_THROUGH for Oracle: v1 scope

Status: proposal, 2026-10-26; §9 lists the open decisions. Step 5 of [the engines design](2026-10-22-walk-through-engines-design.md) §7, which §5.5 there sketches. The user chose to run every engine at once. Every fact here was verified on `gvenzl/oracle-free:23.26.3-slim-faststart` (Oracle AI Database 26ai Free 23.26.3), the image `oracle/parser` pins (Appendix). Bytebase was read at bdebba8604.

## 1. Environment

- **A Bytebase database is one schema.** Its metadata holds one `SchemaMetadata` named `""`, with the objects of the schema the database is named after (`backend/plugin/db/oracle/sync.go`, `SyncDBSchema`).
- **Session.** Bytebase connects as the data source's user and runs `ALTER SESSION SET CURRENT_SCHEMA = "<database>"` in `Open` (`oracle.go`). The change then runs statement by statement on one connection (`Execute`), and stops at the first error.
  - The `ALTER SESSION` runs on the pool. It reaches the connection `Execute` takes only because `database/sql` hands back the one idle connection. The existing rules already rely on this (`schema.go`, `db()`), and the walk does too (§9 O4).
  - That user may or may not own the schema. It matters only to `RENAME`, which works only on the session user's own objects (ORA-03001 otherwise).
- **Current schema:** the target at the start. `ALTER SESSION SET CURRENT_SCHEMA = s` moves unqualified names to `s`. A value the walk cannot read makes them Unknown until a readable one returns them. Other schemas are not synced: their objects are Unknown, apart from what the change itself does to them.
- **Implicit commits.** DDL commits the open transaction before it runs, even when it then fails. No `ROLLBACK` undoes DDL (verified), so a `ROLLBACK` changes nothing the walk keeps.
- **Server settings:** the walk reads none. The server version is not an input. A claim that depends on a 23ai behaviour (`IF [NOT] EXISTS`) is about a statement that older versions reject as a syntax error. That statement then fails, and a later claim is conditional on it succeeding anyway.

## 2. Spaces

Unquoted names fold to upper case, which the parser does. Quoted names keep their case. Names compare exactly.

| Space | Parent | Holds | Claims |
| --- | --- | --- | --- |
| object | schema | the schema object namespace: tables, views, materialized views, sequences, private synonyms, procedures, functions, packages, types | taken (00955); a statement's kind (00942, 02289, 01434, 04043, 12083, 01702, 02201, 04044, 03001) |
| view | schema | views of any sort: relational, object, JSON duality, editioning | missing (00942) |
| sequence | schema | sequences | missing (02289) |
| column | schema.table | columns | exists (01430, 00957), missing (00904) |
| primary key | schema.table | whether the table has one | has one (02260), has none (02441) |
| table constraint | schema.table | the constraints `DROP CONSTRAINT` finds on the table | missing (02443, 23292, 02430, 02431) |
| constraint | schema | constraint names, unique in the schema | taken (02264) |
| index | schema | index names, unique in the schema | taken (00955), missing (01418) |
| trigger | schema | trigger names, unique in the schema | taken (04081, 04095, 00955), missing (04080) |

- **Kinds.** An object key's kind is what the object is. An index, constraint, or trigger key's kind is its table. A `MayRemove` narrowed to that kind then forgets only one table's names when the table goes.
- **Derived keys.** The view and sequence keys repeat what the object key holds for those kinds. Creating an object of another kind makes the name Absent there, since the namespace is shared.
- **Names the walk never claims:**
  - generated names: `SYS_` (constraints, indexes, LOB segments, IOT tables), `ISEQ$$_` (identity sequences), `BIN$` (recycle bin);
  - a name with a byte outside ASCII. The parser folds it with Go's rules, which need not be the database's, so such a name ends the walk;
  - a name with a database link. DDL on a remote object fails, and the walk does not follow it.
- **Length:** an identifier longer than 128 bytes fails with ORA-00972 before anything is looked up, in `CREATE`, `DROP`, and `ALTER` alike. 128 is the limit from 12.2, and 30 before it, so the claim holds on every version. Only ASCII names are judged.
- **Not modeled:** partitions, LOB segments (their namespace is separate: a table may share an LOB segment's name), private database links, clusters, and dimensions.

## 3. Authority

Present is trusted from the sync, with one exception: an index the sync lists as a constraint. The sync lists `P` and `U` constraints as indexes named after the constraint (`getConstraints`). The index behind such a constraint can have another name (`USING INDEX`, or an existing index the constraint reuses) or none (a disabled constraint). Its name is therefore Unknown in the index space.

Absent is trusted only where the sync lists every object of the space:

| Space | Absent trusted? | Why |
| --- | --- | --- |
| object | **No** | The sync lists no synonyms, types, operators, libraries, or Java. `ALL_TABLES` omits object tables and XMLType tables (verified): `CREATE TABLE ot OF ty` is listed by `ALL_OBJECT_TABLES` only, and `DROP TABLE ot` succeeds. See O1 |
| view | Yes | `ALL_VIEWS` lists relational, object, and JSON duality views (verified). No other object answers `DROP VIEW` or `ALTER VIEW`: on a table, materialized view, sequence, synonym, type, or routine they fail with 00942 |
| sequence | Yes, except `ISEQ$$_` | `ALL_SEQUENCES`, less the identity sequences Bytebase skips |
| column | **No** | The sync keeps `COLUMN_ID IS NOT NULL`, which drops invisible columns (verified: `ADD (h)` of an invisible `h` fails with 01430, and `DROP COLUMN h` succeeds) |
| primary key | Yes | Every `P` constraint is listed, whatever its name |
| table constraint | **No** | A supplemental log group is a constraint: `ADD CONSTRAINT lg` fails with 02264, and `DROP CONSTRAINT lg` drops it. `ALL_CONSTRAINTS` does not list it (verified). Unnamed `CHECK` and `NOT NULL` constraints are `SYS_C` |
| constraint | **No** | Views' constraints (`WITH CHECK OPTION CONSTRAINT c`, `WITH READ ONLY CONSTRAINT c`) and log groups share the namespace (02264, verified), and the sync does not list them |
| index | **No** | The sync attaches indexes only to the tables it lists. So it loses those of materialized views, clusters, object tables, and other schemas' tables an index of this schema is on. It also drops an index whose columns are all hidden `SYS_NC` columns, such as one on an object attribute |
| trigger | **No** | The sync keeps triggers with a `TABLE_NAME`. A schema or database trigger (`AFTER LOGON ON SCHEMA`) takes a name in the same namespace (04081, verified) |

So the snapshot yields:
- every "already exists" claim;
- every wrong-kind claim;
- "does not exist" only for views, sequences, and the primary key.

Every other "does not exist" comes from the change itself: an object it created, dropped, or renamed. Under `Target.OtherSpecs` the snapshot is all Unknown, and DML ends the walk (engines design §2).

Accepted assumptions, decision O3:
- the admin user sees every object of the schema in the `ALL_*` views (engines design D6, extended to Oracle);
- the change runs in the edition the sync read;
- no schema-level DDL trigger (D6);
- user code reached only implicitly changes no schema: VPD policies, virtual columns, function-based indexes, auditing handlers, and the triggers of O2.

## 4. Statement semantics the lowering follows

**Resolution.**
- An unqualified name resolves in the current schema first, and the schema's own object wins over a public synonym. So a Present object of the wrong kind is claimed, whatever synonyms exist.
- Some forms do not follow synonyms, public or private: `ALTER TABLE`, `DROP TABLE`, `TRUNCATE`, `DROP VIEW`, `ALTER VIEW`, `DROP SEQUENCE`, `ALTER SEQUENCE`, and `RENAME`. They fail with 00942, 02289, or 04043 on a synonym (verified), so a name absent from the schema is a failure for them.
- Other forms follow both kinds of synonym: `CREATE INDEX ... ON`, `CREATE TRIGGER ... ON`, `COMMENT ON`, `LOCK TABLE`, `GRANT`, `REFERENCES`, and DML (verified). A missing target is therefore never claimed for them.

**`ALTER TABLE` is atomic, and its clauses see the table as it was.**
- `DROP COLUMN`, `RENAME COLUMN`, `RENAME TO`, and `DROP CONSTRAINT` do not combine with the other clause kinds (ORA-12987, ORA-23290, ORA-14047, or a syntax error).
- `ADD` and `MODIFY` of columns combine with `ADD CONSTRAINT`.
- Within one statement:
  - a column may appear once across `ADD`, `MODIFY`, and a `DROP (a, b)` list (00957);
  - a constraint name once (02264);
  - a dropped constraint once (02443);
  - a primary key once (02260).
- So every other check holds against the table before the statement (`walk.Lookups`), and the repeats are text claims. A rejected statement leaves the table as it was (verified for `ADD (n1, b)`, `DROP (a, nosuch)`, and `DROP CONSTRAINT ck DROP CONSTRAINT nosuch`).

**Other statements.**
- `DROP` names one object.
- `CREATE TABLE` and `CREATE TABLE ... AS SELECT` leave nothing when they fail (verified, the second through the index a `UNIQUE` constraint adds).

**`IF [NOT] EXISTS` (23ai).**
- `CREATE ... IF NOT EXISTS` does nothing when an object of any kind holds the name, as with `CREATE TABLE IF NOT EXISTS v` on a view.
  - A column named twice still fails first (00957, verified).
  - A constraint name taken elsewhere does not fail it.
  - On an Unknown name, the creation is a `walk.Builder` `Optional` part.
- `DROP` or `ALTER ... IF EXISTS` does nothing on a missing name or on another kind, as with `DROP TABLE IF EXISTS v`.
  - The exception: `DROP TABLE IF EXISTS mv` of a materialized view fails with 12083.
- `ALTER TABLE IF EXISTS` and `ALTER TABLE ... ADD IF NOT EXISTS` exist on 23ai, but the omni parser rejects the first and Oracle the second.

**Recycle bin.**
- `DROP TABLE` renames the table, its indexes, constraints, and triggers to `BIN$` names. Their names are free at once (verified: `CREATE INDEX`, `ADD CONSTRAINT`, and `CREATE TRIGGER` of the old names succeed).
- `ALL_TABLES`, `ALL_INDEXES`, and `ALL_TRIGGERS` hide dropped objects. `ALL_CONSTRAINTS` keeps them under `BIN$` table names, which the sync attaches to no table.
- `FLASHBACK TABLE` restores what the walk cannot list, so it ends the walk.

**What goes with what.**
- `DROP TABLE` takes the table's indexes, constraints, and triggers.
- `DROP COLUMN` takes the column's single-column constraints and indexes.
- `DROP VIEW` and `CREATE OR REPLACE VIEW` take the view's `INSTEAD OF` triggers and its constraints.
- `CASCADE CONSTRAINTS` takes other tables' foreign keys.
- A renamed table keeps the names of its indexes, constraints, and triggers.
- `DISABLE CONSTRAINT` drops the index the constraint made. `ENABLE` makes it again and can fail on its name (00955).

All verified.

**PL/SQL created with errors succeeds.** `CREATE PROCEDURE`, `FUNCTION`, `PACKAGE [BODY]`, `TRIGGER`, and `FORCE VIEW` with compilation errors leave an `INVALID` object, and go-ora reports no error (verified through go-ora, the driver Bytebase uses). This includes a package body without a spec. So these statements make their names whatever their bodies hold.

**Wrong kinds** (verified):

| Statement | On a view | On a sequence | On a synonym | On a type or routine | On a materialized view |
| --- | --- | --- | --- | --- | --- |
| `DROP TABLE` | 00942 | 00942 | 00942 | 00942 | 12083 |
| `ALTER TABLE` | 00942 | 00942 | 00942 | 00942 | accepted (container table) |
| `TRUNCATE` | 01702 | 02201 | 00942 | 04044 | accepted |
| `CREATE INDEX ... ON` | 01702 | 02201 | follows it | 04044 | accepted |
| `DROP VIEW`, `ALTER VIEW` | — | 00942 | 00942 | 00942 | 00942 |
| `DROP SEQUENCE`, `ALTER SEQUENCE` | 02289 | — | 02289 | 02289 | 02289 |
| `DROP SYNONYM` | 01434 | 01434 | — | 01434 | 01434 |
| `DROP PROCEDURE/FUNCTION/PACKAGE/TYPE` | 04043 | 04043 | 04043 | 04043 for another kind | 04043 |
| `RENAME` | renamed | renamed | renamed | 03001 | 32318 |
| `CREATE` of any kind, `OR REPLACE` of another kind | 00955 | 00955 | 00955 | 00955 | 00955 |

**Renames to the same name fail:**
- `RENAME t TO t`, `ALTER TABLE t RENAME TO t`, `ALTER INDEX i RENAME TO i`, `ALTER TRIGGER tr RENAME TO tr`: 00955;
- `RENAME COLUMN a TO a`: 00957;
- `RENAME CONSTRAINT c TO c`: 02264.

**Dependency failures happen, and the walk makes no claim about them:**
- 02449: dropping a referenced table;
- 02429: dropping a constraint's index;
- 12991, 12992: a column in a multi-column or referenced constraint;
- 02273: a referenced key.

## 5. Statements in v1

"Absent" and "Present" checks claim only what §3 trusts or the change established. Every check of `ALTER TABLE` holds against the table before the statement.

| Statement | Checks | Effects |
| --- | --- | --- |
| `CREATE [GLOBAL TEMPORARY] TABLE [IF NOT EXISTS] t (...)` | t Absent in the object space (00955); text: a column twice (00957, also with IF NOT EXISTS), a constraint name twice (02264), two primary keys (02260), an identifier over 128 bytes (00972); each named constraint Absent in the schema (02264) | t Present (table), Absent as view and sequence; for a column list, exactly its columns, its named constraints Present on t and in the schema, and its primary key exact; index names of named `PRIMARY KEY`/`UNIQUE` Unknown. IF NOT EXISTS: nothing when t is Present, Optional when Unknown |
| `CREATE TABLE t AS SELECT` | as above | t Present with Unknown columns, constraints, and key |
| `ALTER TABLE t ...` | t Present as a table (00942 when Absent or another kind; a materialized view's container is accepted); text: §4's repeats | per subcommand: |
| — `ADD (c ..., constraints)` | c Absent (01430); named constraints Absent in the schema (02264); a primary key t has (02260) | c Present; constraints as in CREATE TABLE |
| — `MODIFY (c ...)` | c Present (00904); inline constraints as in ADD | inline constraints |
| — `DROP COLUMN c`, `DROP (c, ...)`, `SET UNUSED COLUMN c`, `SET UNUSED (c, ...)` [`CASCADE CONSTRAINTS`] | each c Present (00904) | c Absent; t's constraints, primary key, and indexes May-removed; with CASCADE CONSTRAINTS, every constraint of the schema |
| — `RENAME COLUMN a TO b` | a Present (00904); b Absent (00957, b = a too) | a Absent, b Present |
| — `ADD [CONSTRAINT c] PRIMARY KEY/UNIQUE/FOREIGN KEY/CHECK` | c Absent in the schema (02264); a primary key t has (02260) | c Present on t and in the schema; the primary key Present; index c Unknown |
| — `DROP CONSTRAINT c [CASCADE]` | c Present on t (02443) | c Absent on t and in the schema; t's primary key and indexes May-removed; with CASCADE, every constraint |
| — `DROP PRIMARY KEY [CASCADE]` | t has a primary key (02441) | the primary key Absent; t's constraints and indexes May-removed |
| — `DROP UNIQUE (...)` | none | t's constraints and indexes May-removed |
| — `RENAME CONSTRAINT a TO b` | a Present on t (23292); b Absent in the schema (02264) | a Absent, b Present, on t and in the schema |
| — `ENABLE`, `DISABLE`, `MODIFY CONSTRAINT c`; `ENABLE`/`DISABLE PRIMARY KEY`/`UNIQUE (...)` | c Present on t (02430, 02431) | index names May-added and May-removed for t |
| — `RENAME TO u` | u Absent in the object space (00955, u = t too) | t moves to u; t's index, constraint, and trigger names Unknown (their keys name the old table) |
| — table properties, partition maintenance, `MOVE`, `SHRINK`, `ENABLE`/`DISABLE ALL TRIGGERS`/`TABLE LOCK`, `DROP UNUSED COLUMNS` | none | none |
| — `ADD SUPPLEMENTAL LOG`, anything else | | ends the walk |
| `DROP TABLE t [CASCADE CONSTRAINTS] [PURGE]`, `DROP TABLE IF EXISTS t` | t Present as a table (00942 when Absent or another kind; 12083 for a materialized view, IF EXISTS too) | t Absent in every space; its columns, constraints, and key cleared; its indexes, constraints, and triggers May-removed in the schema; with CASCADE CONSTRAINTS, every constraint. IF EXISTS: nothing on another kind, Optional when Unknown |
| `TRUNCATE TABLE t` | t Absent (00942), a view (01702), a sequence (02201), a type or routine (04044) | none |
| `CREATE [UNIQUE\|BITMAP] INDEX [IF NOT EXISTS] i ON t (...)` | i Absent in the index space (00955); t a view (01702), a sequence (02201), a type or routine (04044); a missing t claims nothing | i Present (of t) |
| `DROP INDEX i [IF EXISTS]` | i Present (01418) | i Absent |
| `ALTER INDEX i RENAME TO j` | i Present (01418); j Absent (00955, j = i too) | i Absent, j Present (i's table) |
| `ALTER INDEX i` other actions | i Present (01418) | none |
| `CREATE [OR REPLACE] [FORCE] [EDITIONING] VIEW v [IF NOT EXISTS]` | v Absent in the object space (00955); OR REPLACE: v no other kind (00955); a named view constraint Absent in the schema (02264) | v Present (view); OR REPLACE of a view: its triggers and constraints May-removed; its new constraint names Present |
| `ALTER VIEW v ... ` | v Present as a view (00942) | `ADD`/`DROP CONSTRAINT`: the name Unknown |
| `DROP VIEW v [CASCADE CONSTRAINTS] [IF EXISTS]` | v Present as a view (00942) | v Absent; its triggers and constraints May-removed |
| `CREATE SEQUENCE s [IF NOT EXISTS]` | s Absent (00955) | s Present (sequence) |
| `ALTER SEQUENCE s`, `DROP SEQUENCE s [IF EXISTS]` | s Present as a sequence (02289) | none / s Absent |
| `CREATE [OR REPLACE] SYNONYM s [IF NOT EXISTS] FOR x` | s Absent (00955); OR REPLACE: no other kind (00955) | s Present (synonym) |
| `DROP SYNONYM s [FORCE] [IF EXISTS]` | s Present as a synonym (01434) | s Absent |
| `CREATE`/`DROP PUBLIC SYNONYM` | none | none: the walk never relies on a public synonym's absence |
| `CREATE [OR REPLACE] TYPE`, `PROCEDURE`, `FUNCTION`, `PACKAGE [IF NOT EXISTS]` | the name Absent (00955); OR REPLACE: no other kind (00955) | the name Present (its kind); the body is not walked |
| `CREATE [OR REPLACE] TYPE BODY`, `PACKAGE BODY` | none | none: bodies have a namespace of their own |
| `DROP TYPE`, `PROCEDURE`, `FUNCTION`, `PACKAGE [IF EXISTS]` | the name Present as that kind (04043) | the name Absent |
| `DROP TYPE BODY`, `PACKAGE BODY` | none | none |
| `CREATE [OR REPLACE] TRIGGER tr [IF NOT EXISTS] ... ON t`, a DML or `INSTEAD OF` trigger | tr Absent (04081); OR REPLACE: tr not on another table (04095); t a view for a DML trigger is not claimed | tr Present (of t); the walk records that a trigger exists (§6) |
| `DROP TRIGGER tr [IF EXISTS]`; `ALTER TRIGGER tr ENABLE/DISABLE/COMPILE` | tr Present (04080) | tr Absent / none |
| `ALTER TRIGGER tr RENAME TO tr2` | tr Present (04080); tr2 Absent (00955, tr2 = tr too) | tr Absent, tr2 Present |
| `RENAME a TO b` | a Present (04043), and a table, view, sequence, or synonym (03001; 32318 for a materialized view); b Absent (00955, b = a too) | as `ALTER TABLE RENAME TO` for a table; the name moves for the other kinds |
| `COMMENT`, `GRANT`, `REVOKE`, `LOCK TABLE`, `ANALYZE`, `EXPLAIN PLAN`, `PURGE`, `AUDIT`, `NOAUDIT`, `COMMIT`, `ROLLBACK [TO SAVEPOINT]`, `SAVEPOINT`, `SET TRANSACTION`, `SET CONSTRAINTS` | none | none |
| `ALTER SESSION SET CURRENT_SCHEMA = s`, and parameters on an allowlist (`NLS_*`, `TIME_ZONE`, `DDL_LOCK_TIMEOUT`, `PLSQL_*`, `OPTIMIZER_*`, ...) | none | the current schema |
| `INSERT`, `UPDATE`, `DELETE`, `MERGE`, `SELECT` | none: the target may be a public synonym | none; §6 says when they end the walk |

The parser's shapes, read before writing this table:
- `DROP (a, b)` and `SET UNUSED (a, b)` record only the first column, so the lowering reads the list from the tokens.
- `DROP SYNONYM` does not record `PUBLIC`, which the tokens give.
- A DDL trigger's table is the bare word `SCHEMA` or `DATABASE`.
- `ADD SUPPLEMENTAL LOG` is a table property, so the property allowlist goes by name.

## 6. Where the walk stops

- **Before the statement:**
  - anonymous PL/SQL blocks and `CALL` (`EXECUTE IMMEDIATE` lives only in PL/SQL);
  - `FLASHBACK TABLE` and `FLASHBACK DATABASE`;
  - materialized view and materialized view log DDL;
  - `CREATE PRIVATE TEMPORARY TABLE`;
  - `CREATE TABLE` with a nested table or an external organization (the parser rejects `OF type` and `CLUSTER`);
  - DDL triggers `ON SCHEMA` or `ON DATABASE`;
  - `ALTER SESSION SET EDITION`, `CONTAINER`, or any parameter off the allowlist;
  - `ALTER SYSTEM`;
  - user, role, profile, and tablespace DDL;
  - a name outside ASCII or with a database link;
  - anything not in §5.
- **After its checks, for a statement that may run user code:**
  - Which statements: DML, `SELECT`, `CREATE TABLE ... AS SELECT`, and `CREATE INDEX` on an expression.
  - When: any function call that is qualified, or whose name is not an Oracle SQL builtin.
  - In SQL a builtin takes precedence over a schema function of the same name, while a qualified call runs the schema's (verified with a schema `UPPER` and `NVL`). An autonomous function called from `SELECT` or `CREATE TABLE ... AS SELECT` can run DDL (verified).
- **After its checks, DML (`INSERT`, `UPDATE`, `DELETE`, `MERGE`) where a trigger may fire:**
  - The cases:
    - the snapshot or the change has any trigger;
    - the target is not a table the walk knows in a schema: it may be a view, or a public synonym for another schema's table;
    - `Target.OtherSpecs`.
  - Verified: a trigger can run DDL through an autonomous procedure, and a cascaded delete fires the child table's trigger.
  - O2 covers triggers the sync cannot show.
- **Accepted assumptions:** O3.

## 7. Not covered (declined by scope)

Review asks for these are declined with this list:

- "does not exist" from the snapshot for tables (O1), columns, indexes, constraints, triggers, synonyms, types, and routines (§3);
- missing targets of statements that follow synonyms: `CREATE INDEX ... ON`, `CREATE TRIGGER ... ON`, `COMMENT`, `LOCK TABLE`, `GRANT`, `REFERENCES`, DML;
- dependency failures (§4);
- the index a `PRIMARY KEY` or `UNIQUE` constraint adds or reuses, and its name collisions (00955 on `ADD CONSTRAINT` or `ENABLE`);
- columns named in constraints, indexes, views, and queries; a view's query (`CREATE VIEW` failing with 00904 or 00942); `ORA-01408`;
- data-dependent failures: `NOT NULL` on rows with nulls, 01442, type changes (01439, 01440), constraint validation;
- privileges, editions, and `RENAME` by a user who does not own the schema (03001);
- partitions, LOBs, nested tables, object tables, clusters, index-organized and external tables, materialized views, private temporary tables, database links, users, tablespaces;
- PL/SQL bodies and their compilation;
- names outside ASCII;
- 23ai syntax the parser rejects (`ALTER TABLE IF EXISTS`, `CREATE TABLE ... OF`);
- the recycle bin's contents (`FLASHBACK`, `PURGE` of a name).

## 8. Budget and tests

- **Code:** about 1.3k lines of lowering, 350 of names and snapshot, and 300 of walker and plumbing. `targets.go` gains `walks` and `OtherSpecs` in its group key. Growth past this is a decision for a person.
- **Oracle:** `walktest.Judge` against `gvenzl/oracle-free:23.26.3-slim-faststart`, in the container `oracle/review`'s tests already start.
  - Oracle commits DDL as it runs, so the engine implements `walktest.Digester` over `ALL_OBJECTS`, `ALL_TAB_COLS`, `ALL_CONSTRAINTS`, `ALL_INDEXES`, and `ALL_TRIGGERS`.
  - The change runs as an admin user with `CURRENT_SCHEMA = SHOP`, as Bytebase runs it. `RENAME` cases run as SHOP, which owns the objects.
  - The connect, the container's start, and every statement run under a deadline, so a hung container fails the test instead of stalling CI (it hung once, in #493's 4829c2fb run).
  - The test sync mirrors Bytebase's for every space the walk reads: invisible columns dropped, `SYS_C` checks dropped, indexes of listed tables only, and triggers with a table.
- **Curated cases:** one per row of §5 and per fact in the Appendix, including the wrong kinds, the synonym-following forms, and `IF [NOT] EXISTS`.
- **Sweep:** templates over a colliding name pool, mixing case and quoted names:
  - tables, views, sequences, synonyms, indexes named like tables, and constraint names;
  - `IF [NOT] EXISTS`, `RENAME`, and triggers;
  - `CURRENT_SCHEMA` switches.

  CI runs it at `WALK_SWEEP_N`. Precision must hold, and recall is logged.
- **Coverage:** `walktest.Coverage` over every row of §5. **Unit tests:** the snapshot's authority per space, and `OtherSpecs` in the group key.

## 9. Decisions

Open, each with a recommendation:

- **O1. Table absence from the sync.**
  - The problem: `ALL_TABLES` omits object and XMLType tables, which the sync therefore never lists (verified).
  - Options:
    - (a) Untrusted until Bytebase syncs them (`ALL_ALL_TABLES`, or a flag in the metadata). The change's own drops and renames still give "does not exist".
    - (b) Trusted, with an accepted assumption that the schema holds no object table.
  - Recommended: (a), as D5 decided for MariaDB's unsynced relation kinds. It costs the ORA-00942 claims on tables the sync does not list. The Bytebase fix is to be filed.
- **O2. DML and triggers the sync cannot show.**
  - The problem: the sync lists only triggers the schema owns. A trigger another schema owns can be on the schema's tables, or on a table in another schema that a cascaded delete reaches.
  - Options:
    - (a) DML goes on as §6 says, with the accepted assumption that no other schema's trigger fires on a change to the schema's tables, directly or through a cascade.
    - (b) Every DML statement ends the walk after its checks.
  - Recommended: (a). (b) would end the walk at the first data fix in most changes.
- **O3. The assumptions of §3:**
  - the admin user sees every object of the schema;
  - the change runs in the edition the sync read;
  - no schema-level DDL trigger;
  - implicit user code changes no schema.

  Recommended: accept. The first and third are D6. Without the first, a view the admin cannot see is one it cannot drop either: `DROP VIEW` then fails with 00942, as the walk claims.
- **O4. `CURRENT_SCHEMA` reaches the execution connection only through `database/sql`'s reuse of the one idle connection (§1).**
  - Options:
    - (a) Rely on it, as the existing rules do, and file a Bytebase follow-up to set the schema on the connection `Execute` uses.
    - (b) Treat unqualified names as Unknown.
  - Recommended: (a). (b) would leave the walk almost nothing, since changes rarely qualify names.
- **O5. Versions before 23ai.**
  - Every code here is long-standing, and the length claim holds from 11g. The wrong-kind table was verified on 23ai only.
  - Options:
    - (a) Spot-check the curated cases once on `gvenzl/oracle-xe:21-full-faststart`, an image already on the devbox, and record any differences. CI keeps 23ai.
    - (b) 23ai only.
  - Recommended: (a). Where 21c gives another code for a statement that still fails, the code map lists both. Where it accepts the statement, the claim goes.

## Appendix: verified facts

Oracle AI Database 26ai Free 23.26.3.0.0 (`gvenzl/oracle-free:23.26.3-slim-faststart`). The change ran as an admin user (DBA) with `CURRENT_SCHEMA = SHOP` unless noted "as SHOP". Codes are `ORA-`.

- **Already exists:**
  - `CREATE TABLE` / `VIEW` / `SEQUENCE` / `SYNONYM` / `TYPE` / `PROCEDURE` / `FUNCTION` / `PACKAGE` of a name any of them holds: 00955. `OR REPLACE` of another kind: 00955. `CREATE OR REPLACE FUNCTION` of a function: accepted.
  - `CREATE INDEX` of an index name: 00955. An index may share a table's name.
  - A constraint name another table's constraint, or a view's, has (`ADD CONSTRAINT`, inline in `CREATE TABLE`, `WITH CHECK OPTION CONSTRAINT`): 02264.
  - `ADD (a)` of a column: 01430; `RENAME COLUMN a TO b` with b present: 00957.
  - `ADD PRIMARY KEY` with one: 02260.
  - `CREATE TRIGGER` of a trigger name, on any table: 04081. `CREATE OR REPLACE TRIGGER` of another table's trigger: 04095. A trigger may share a table's name.
  - `RENAME`, `ALTER TABLE RENAME TO`, `ALTER INDEX RENAME TO`, `ALTER TRIGGER RENAME TO` onto a taken name, or onto the same one: 00955. `RENAME COLUMN a TO a`: 00957. `RENAME CONSTRAINT c TO c`: 02264.
- **Missing:**
  - `ALTER TABLE` / `DROP TABLE` / `TRUNCATE` / `DROP VIEW` / `ALTER VIEW`: 00942, and 23ai names the object (`table or view "SHOP"."NOSUCH" does not exist`).
  - `DROP COLUMN` / `MODIFY` / `RENAME COLUMN` of a missing column, and one in a `DROP (a, nosuch)` list: 00904.
  - `DROP CONSTRAINT` of a missing constraint, or of another table's: 02443. `RENAME CONSTRAINT`: 23292. `ENABLE CONSTRAINT`: 02430. `DISABLE` or `MODIFY CONSTRAINT`: 02431.
  - `DROP PRIMARY KEY` without one: 02441. `DROP UNIQUE (a)` without one: 02442.
  - `DROP INDEX` / `ALTER INDEX`: 01418.
  - `DROP SEQUENCE` / `ALTER SEQUENCE`: 02289.
  - `DROP SYNONYM`: 01434.
  - `DROP PROCEDURE` / `FUNCTION` / `PACKAGE` / `TYPE`: 04043.
  - `DROP TRIGGER` / `ALTER TRIGGER`: 04080.
  - `RENAME` of a missing name: 04043.
  - `COMMENT ON COLUMN t.nosuch`: 00904.
- **Wrong kinds:** the table in §4, each cell probed. `RENAME` of a procedure or type fails with 03001, and of a materialized view with 32318, as SHOP too.
- **Text:**
  - In `CREATE TABLE`: a column twice, 00957; a constraint name twice, 02264; two primary keys, 02260.
  - In `ALTER TABLE`:
    - a column in `ADD` and `MODIFY`, twice in `ADD`, or twice in `DROP (...)`: 00957;
    - a constraint name twice across `ADD` clauses or inline: 02264;
    - two primary keys across clauses: 02260;
    - `DROP CONSTRAINT c DROP CONSTRAINT c`: 02443.
  - An identifier of 129 bytes in `CREATE TABLE`, `DROP TABLE`, or `ADD`: 00972; 128 is accepted.
- **Combinations:** `ADD ... DROP COLUMN`: 12987. `RENAME COLUMN` with anything: 23290. `RENAME TO shop.t2`: 14047. `DROP COLUMN` / `DROP CONSTRAINT` / `DROP PRIMARY KEY` / `MODIFY` followed by another clause kind: a syntax error (03048). `ADD (col) ADD CONSTRAINT ... MODIFY (...)`: accepted. `ADD (n2) ADD CONSTRAINT pk_t` with pk_t taken: 02264.
- **Synonyms:**
  - Not followed (00942 or 02289, public and private): `ALTER TABLE`, `TRUNCATE`, `DROP TABLE`, `DROP VIEW`, `ALTER VIEW`, `DROP SEQUENCE`, `ALTER SEQUENCE`. `RENAME` of a public synonym's name: 04043.
  - Followed: `CREATE INDEX ... ON`, `COMMENT ON TABLE` and `COLUMN`, `LOCK TABLE`, `CREATE TRIGGER ... ON`, `GRANT`, `REFERENCES` (a private synonym resolves; a public one reaches the other schema's table).
- **Atomicity:**
  - A rejected `ALTER TABLE ADD (n1, b)` adds neither column. `DROP (a, nosuch)` keeps a. `DROP CONSTRAINT ck_c DROP CONSTRAINT nosuch` keeps ck_c.
  - A `CREATE TABLE` refused for its `UNIQUE` index's name leaves neither table nor constraint. A failed `CREATE TABLE ... AS SELECT` leaves no table.
  - DDL commits a pending `INSERT` whether it succeeds or fails.
- **IF [NOT] EXISTS:**
  - `CREATE TABLE IF NOT EXISTS` of a table, view, sequence, or synonym name: a no-op. `CREATE TABLE IF NOT EXISTS t (a, a)` with t present: 00957. `CREATE TABLE IF NOT EXISTS t (... CONSTRAINT pk_t ...)` with t and pk_t present: a no-op.
  - `DROP TABLE IF EXISTS` of a view or sequence, `DROP VIEW IF EXISTS` of a table or object table, `DROP SEQUENCE IF EXISTS` of a table, `ALTER TABLE IF EXISTS` of a view: no-ops. `DROP TABLE IF EXISTS` of a materialized view: 12083. `DROP TABLE IF EXISTS` of an object table drops it.
  - `CREATE INDEX IF NOT EXISTS` of a taken index name: a no-op.
  - `CREATE VIEW` / `SEQUENCE` / `SYNONYM` / `PROCEDURE` / `TRIGGER IF NOT EXISTS`, and `DROP TRIGGER` / `PROCEDURE` / `SYNONYM` / `TYPE` / `INDEX IF EXISTS`: accepted.
  - `OR REPLACE` with `IF NOT EXISTS`: 11541. `ALTER INDEX IF EXISTS ... RENAME`: 11557. `ALTER TABLE ... ADD IF NOT EXISTS`, `DROP COLUMN IF EXISTS`, `DROP CONSTRAINT IF EXISTS`: syntax errors.
- **Sync gaps:**
  - `CREATE TABLE ot OF ty` and `CREATE TABLE xt OF XMLTYPE` are absent from `ALL_TABLES` and present in `ALL_OBJECT_TABLES`.
  - An `INVISIBLE` column has a null `COLUMN_ID`.
  - A supplemental log group occupies the constraint namespace, is absent from `ALL_CONSTRAINTS`, and `DROP CONSTRAINT` drops it.
  - View constraints are types `V` and `O`.
  - A schema `LOGON` trigger has no `TABLE_NAME` and occupies the trigger namespace.
  - `ALL_VIEWS` lists object views and JSON duality views.
  - An identity column adds `ISEQ$$_<n>`. A materialized view's index can be `SYS_C_SNAP$_...`.
- **Recycle bin:**
  - After `DROP TABLE t`, `ALL_TABLES`, `ALL_INDEXES`, `ALL_TRIGGERS`, and `ALL_TAB_COLUMNS` hide t's objects. `ALL_CONSTRAINTS` lists its constraints under `BIN$` names, as does `DBA_RECYCLEBIN`.
  - The old index, constraint, and trigger names can be created again.
  - `FLASHBACK TABLE t TO BEFORE DROP` with a new t: 38312. `PURGE TABLE t` of a live t: 38307.
- **Effects:**
  - `DROP COLUMN` drops the column's check and its single-column primary key. `DROP TABLE p` referenced by a foreign key: 02449. With `CASCADE CONSTRAINTS`, the child's foreign key name is free.
  - `CREATE OR REPLACE VIEW` drops the view's `INSTEAD OF` triggers and a check-option constraint it no longer declares. `DROP VIEW` frees both.
  - `RENAME` (as SHOP) and `ALTER TABLE RENAME TO` keep the index, constraint, and trigger names, and the trigger follows the table.
  - `ADD CONSTRAINT pk PRIMARY KEY` reuses an index on the columns and adds no index named pk. With `USING INDEX (CREATE INDEX ix ...)` the index is ix.
  - `DISABLE CONSTRAINT pk_q` drops index PK_Q. After `CREATE INDEX pk_q` elsewhere, `ENABLE` fails with 00955.
- **Session user:** `RENAME` of a table in SHOP by the admin: 03001; as SHOP it renames tables, views, sequences, and private synonyms.
- **User code:**
  - A schema `UPPER` and `NVL` lose to the builtins in SQL; `shop.upper(...)` runs the schema's.
  - An autonomous function called by `SELECT` created a table. `CREATE TABLE ... AS SELECT` ran the function.
  - A `BEFORE DELETE` trigger on a child table, fired by an `ON DELETE CASCADE` from its parent, created a table through an autonomous procedure.
- **PL/SQL with errors:** through go-ora, `CREATE OR REPLACE PROCEDURE` with an unknown call, `CREATE OR REPLACE PACKAGE BODY` without a spec, a trigger with an unknown call, and `CREATE OR REPLACE FORCE VIEW` on a missing table return no error and leave `INVALID` objects. `CREATE OR REPLACE VIEW` with an unknown column fails with 00904.
- **Namespaces:** a LOB segment `STORE AS lseg` beside a table `LSEG`: accepted. Quoted `"t"` and unquoted `T` are two tables; `DROP TABLE t2x` does not find `"T2x"`.
