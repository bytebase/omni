# WALK_THROUGH for Oracle: v1 scope

Status: decided 2026-10-28; §9 records the decisions. Step 5 of [the engines design](2026-10-22-walk-through-engines-design.md) §7, which §5.5 there sketches. The user chose to run every engine at once. Every fact here was verified on `gvenzl/oracle-free:23.26.3-slim-faststart` (Oracle AI Database 26ai Free 23.26.3), the image `oracle/parser` pins (Appendix). Bytebase was read at 36190dc616, after its Oracle sync lists object tables and invisible columns (#21674) and its driver sets the current schema on every session (#21672).

## 1. Environment

- **A Bytebase database is one schema.** Its metadata holds one `SchemaMetadata` named `""`, with the objects of the schema the database is named after (`backend/plugin/db/oracle/sync.go`, `SyncDBSchema`).
- **Session.** Bytebase connects as the data source's user, and its connector runs `ALTER SESSION SET CURRENT_SCHEMA = "<database>"` on every session it opens (`oracle.go`, `connector.Connect`). The change then runs statement by statement on one connection (`Execute`), and stops at the first error.
  - That user may or may not own the schema. It matters only to `RENAME`, which works only on the session user's own objects (ORA-03001 otherwise).
- **Current schema:** the target at the start. `ALTER SESSION SET CURRENT_SCHEMA = s` moves unqualified names to `s`. A value the walk cannot read makes them Unknown until a readable one returns them. Other schemas are not synced: their objects are Unknown, apart from what the change itself does to them.
- **Implicit commits.** DDL commits the open transaction before it runs, even when it then fails. No `ROLLBACK` undoes DDL (verified), so a `ROLLBACK` changes nothing the walk keeps.
- **Server settings:** the walk reads none. The server version is not an input. A claim that depends on a 23ai behaviour (`IF [NOT] EXISTS`) is about a statement that older versions reject as a syntax error. That statement then fails, and a later claim is conditional on it succeeding anyway.

## 2. Spaces

Unquoted names fold to upper case, which the parser does. Quoted names keep their case. Names compare exactly.

| Space | Parent | Holds | Claims |
| --- | --- | --- | --- |
| object | schema | the schema object namespace: tables, views, materialized views, sequences, private synonyms, procedures, functions, packages, types | taken (00955); a statement's kind (02289, 01434, 04043, 01702, 02201, 04044, 03001, 32318) |
| table | schema | tables of any sort, and materialized views | missing (00942), a materialized view (12083) |
| view | schema | views of any sort: relational, object, JSON duality, editioning | missing (00942) |
| sequence | schema | sequences | missing (02289) |
| column | schema.table | columns | exists (01430, 00957), missing (00904) |
| primary key | schema.table | whether the table has one | has one (02260), has none (02441) |
| table constraint | schema.table | the constraints `DROP CONSTRAINT` finds on the table | missing (02443, 23292, 02430, 02431) |
| constraint | schema | constraint names, unique in the schema | taken (02264) |
| index | schema | index names, unique in the schema | taken (00955), missing (01418) |
| trigger | schema | trigger names, unique in the schema | taken (04081, 04095, 00955), missing (04080) |

- **Kinds.** An object key's kind is what the object is. An index, constraint, or trigger key's kind is its table. A `MayRemove` narrowed to that kind then forgets only one table's names when the table goes.
- **Derived keys.** The table, view, and sequence keys repeat what the object key holds for those kinds. Creating an object of another kind makes the name Absent there, since the namespace is shared.
- **Names the walk never claims:**
  - generated names: `SYS_` (constraints, indexes, LOB segments, IOT tables, hidden columns), `ISEQ$$_` (identity sequences), `BIN$` (recycle bin), `ORA$` (private temporary tables); and among columns, the hidden columns Oracle names itself, `ORA_ARCHIVE_STATE` and `ORABCTAB_...$`, on which DDL fails with an error of its own (38398, 05715);
  - a name with a byte outside ASCII. The parser folds it with Go's rules, which need not be the database's, so such a name ends the walk;
  - a name with a database link. DDL on a remote object fails, and the walk does not follow it.
- **Length:** an identifier longer than 128 bytes fails with ORA-00972 before anything is looked up, in `CREATE`, `DROP`, and `ALTER` alike. 128 is the limit from 12.2, and 30 before it, so the claim holds on every version. Only ASCII names are judged.
- **Not modeled:** partitions, LOB segments (their namespace is separate: a table may share an LOB segment's name), private database links, clusters, and dimensions.

## 3. Authority

Present is trusted from the sync, with one exception: an index the sync lists as a constraint. The sync lists `P` and `U` constraints as indexes named after the constraint (`getConstraints`). The index behind such a constraint can have another name (`USING INDEX`, or an existing index the constraint reuses) or none (a disabled constraint). Its name is therefore Unknown in the index space.

Absent is trusted only where the sync lists every object of the space. Bytebase's metadata is taken to be what its current sync produces (decision O1):

| Space | Absent trusted? | Why |
| --- | --- | --- |
| object | **No** | The sync lists no synonyms, types, operators, libraries, Java, or property graphs |
| table | Yes, but generated names | The sync lists `ALL_TABLES` and, with `SkipDump`, `ALL_OBJECT_TABLES` (object and XMLType tables), less materialized views' containers, and lists materialized views apart. `ALL_TABLES` lists JSON collection tables and nested tables' storage tables (verified). No other object answers `ALTER TABLE`, `DROP TABLE`, or `TRUNCATE`: they fail with 00942 or an error of their own on a view, a synonym, a sequence, or a routine (§4) |
| view | Yes | `ALL_VIEWS` lists relational, object, and JSON duality views (verified). No other object answers `DROP VIEW` or `ALTER VIEW`: on a table, materialized view, sequence, synonym, type, or routine they fail with 00942 |
| sequence | Yes, except `ISEQ$$_` | `ALL_SEQUENCES`, less the identity sequences Bytebase skips |
| column | Yes, for a table whose columns the sync lists in full, but generated names | The sync lists `ALL_TAB_COLS` with `USER_GENERATED = 'YES'`, invisible columns included. Not in full: an object or XMLType table, whose hidden `XMLDATA` `DROP COLUMN` drops (verified); a table it lists no column of, as a nested table's storage table, whose columns `ALL_TAB_COLS` does not show (verified); and a table with a check constraint naming a column it does not list, which a valid-time period adds with its hidden columns (verified: `PERIOD FOR vt` adds `VT`, `VT_START`, `VT_END`, and a check `"VT_START" < "VT_END"`) |
| primary key | Yes | Every `P` constraint is listed, whatever its name |
| table constraint | **No** | A supplemental log group is a constraint: `ADD CONSTRAINT lg` fails with 02264, and `DROP CONSTRAINT lg` drops it. `ALL_CONSTRAINTS` does not list it (verified). Unnamed `CHECK` and `NOT NULL` constraints are `SYS_C`. A constraint name the sync lists on another table is not on this one: constraint names are unique in the schema |
| constraint | **No** | Views' constraints (`WITH CHECK OPTION CONSTRAINT c`, `WITH READ ONLY CONSTRAINT c`) and log groups share the namespace (02264, verified), and the sync does not list them |
| index | **No** | The sync lists the indexes of its tables and materialized views. It loses those of clusters and of other schemas' tables an index of this schema is on, and an index whose columns are all hidden `SYS_NC` columns, such as one on an object attribute |
| trigger | **No** | The sync keeps the triggers of this schema on its own tables and views. A schema or database trigger (`AFTER LOGON ON SCHEMA`) takes a name in the same namespace (04081, verified), and so does one on another schema's table |

So the snapshot yields:
- every "already exists" claim;
- every wrong-kind claim;
- "does not exist" for tables, views, sequences, the columns of a table listed in full, and the primary key.

Every other "does not exist" comes from the change itself: an object it created, dropped, or renamed. Under `Target.OtherSpecs` the snapshot is all Unknown, and DML ends the walk (engines design §2).

Accepted assumptions, decision O3:
- the admin user sees every object of the schema in the `ALL_*` views (engines design D6, extended to Oracle);
- the change runs in the edition the sync read;
- no schema-level DDL trigger (D6);
- user code reached only implicitly changes no schema: VPD policies, virtual columns, function-based indexes, auditing handlers, and the triggers of O2;
- no hidden column of a name of the user's but a valid-time period's: Oracle Label Security's policy column, which may be hidden, is one the walk would take for missing.

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
- `DROP TABLE` takes the table's indexes, constraints, and triggers, and a nested table column's storage table, which the sync lists as a table of no column without its parent. `DROP COLUMN` of the nested table column takes it too, and `SET UNUSED` keeps it until `DROP UNUSED COLUMNS`. So these forget every storage table the sync lists in the schema, and `ALTER TABLE ... ADD (c ...) NESTED TABLE c STORE AS s` ends the walk.
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
| `CREATE TABLE` with a nested table, an external organization, or a supplemental log group; `CREATE PRIVATE TEMPORARY TABLE` | | ends the walk |
| `ALTER TABLE t ...` | t Present in the table space (00942 when Absent; a materialized view's container is accepted); text: §4's repeats | per subcommand: |
| — `ADD (c ..., constraints)` | c Absent (01430); named constraints Absent in the schema (02264); a primary key t has (02260) | c Present; constraints as in CREATE TABLE |
| — `MODIFY (c ...)` | c Present (00904); inline constraints as in ADD | inline constraints |
| — `DROP COLUMN c`, `DROP (c, ...)`, `SET UNUSED COLUMN c`, `SET UNUSED (c, ...)` [`CASCADE CONSTRAINTS`] | each c Present (00904) | c Absent; t's constraints, primary key, and indexes May-removed; with CASCADE CONSTRAINTS, every constraint of the schema |
| — `RENAME COLUMN a TO b` | a Present (00904); b Absent (00957, b = a too) | a Absent, b Present |
| — `ADD [CONSTRAINT c] PRIMARY KEY/UNIQUE/FOREIGN KEY/CHECK` | c Absent in the schema (02264); a primary key t has (02260) | c Present on t and in the schema; the primary key Present; index c Unknown |
| — `DROP CONSTRAINT c [CASCADE]` | c Present on t, and not another table's in the schema (02443) | c Absent on t and in the schema; t's primary key and indexes May-removed; with CASCADE, every constraint |
| — `DROP PRIMARY KEY [CASCADE]` | t has a primary key (02441) | the primary key Absent; t's constraints and indexes May-removed |
| — `DROP UNIQUE (...)` | none | t's constraints and indexes May-removed |
| — `RENAME CONSTRAINT a TO b` | a Present on t, and not another table's (23292); b Absent in the schema (02264) | a Absent, b Present, on t and in the schema |
| — `ENABLE`, `DISABLE`, `MODIFY CONSTRAINT c`; `ENABLE`/`DISABLE PRIMARY KEY`/`UNIQUE (...)` | c Present on t, and not another table's (02430, 02431) | index names May-added and May-removed for t |
| — `RENAME TO u` | u Absent in the object space (00955, u = t too) | t moves to u; t's index, constraint, and trigger names Unknown (their keys name the old table) |
| — table properties, partition maintenance, `MOVE`, `SHRINK`, `ENABLE`/`DISABLE ALL TRIGGERS`/`TABLE LOCK`, `DROP UNUSED COLUMNS` | none | none |
| — `ADD SUPPLEMENTAL LOG`, `ADD`/`DROP PERIOD`, anything else | | ends the walk |
| `DROP TABLE t [CASCADE CONSTRAINTS] [PURGE]`, `DROP TABLE IF EXISTS t` | t Present in the table space (00942 when Absent; 12083 for a materialized view, IF EXISTS too) | t Absent in every space; its columns, constraints, and key cleared; its indexes, constraints, and triggers May-removed in the schema; with CASCADE CONSTRAINTS, every constraint. IF EXISTS: nothing on another kind, Optional when Unknown |
| `TRUNCATE TABLE t` | t Absent from the table space (00942), a view (01702), a sequence (02201), a type or routine (04044) | none |
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
  - When: any function call that is qualified, or whose name is not an Oracle SQL builtin; or a name that may call a function without parentheses. In a statement of one query block, a name is no call when it is a column the walk knows a table in scope has, a column a table or alias in scope qualifies, a sequence's `NEXTVAL` or `CURRVAL`, a builtin Oracle takes first, or a name an index lists alone. Any other name, or any name in a statement of more query blocks, may be a call.
  - In SQL a builtin takes precedence over a schema function of the same name, while a qualified call runs the schema's (verified with a schema `UPPER` and `NVL`). An autonomous function called from `SELECT` or `CREATE TABLE ... AS SELECT` can run DDL (verified).
- **After its checks, a statement that may run an indextype's code:** `CREATE INDEX ... INDEXTYPE IS`, and DDL or DML on a table with a domain index or on such an index (the sync gives the index type `DOMAIN`). An indextype's `ODCIIndex*` routines are user code, and Oracle Text's make and drop `DR$` tables of their own. Oracle Free has no Oracle Text (29833), so this is pinned by unit tests.
- **After its checks, DML (`INSERT`, `UPDATE`, `DELETE`, `MERGE`) where a trigger may fire:**
  - The cases:
    - the snapshot or the change has any trigger;
    - the target is not a table the walk knows in its schema: it may be a view, or a public synonym for another schema's table;
    - `Target.OtherSpecs`.
- **A name of an `ALTER SESSION SET CURRENT_SCHEMA` the walk cannot read** ends the walk at the next statement that names an object unqualified.
  - Verified: a trigger can run DDL through an autonomous procedure, and a cascaded delete fires the child table's trigger.
  - O2 covers triggers the sync cannot show.
- **Accepted assumptions:** O3.

## 7. Not covered (declined by scope)

Review asks for these are declined with this list:

- "does not exist" from the snapshot for indexes, constraints, triggers, synonyms, types, and routines, and for the columns of a table the sync may not list in full (§3);
- missing targets of statements that follow synonyms: `CREATE INDEX ... ON`, `CREATE TRIGGER ... ON`, `COMMENT`, `LOCK TABLE`, `GRANT`, `REFERENCES`, DML;
- dependency failures (§4);
- the index a `PRIMARY KEY` or `UNIQUE` constraint adds or reuses, and its name collisions (00955 on `ADD CONSTRAINT` or `ENABLE`);
- columns named in constraints, indexes, views, and queries; a view's query (`CREATE VIEW` failing with 00904 or 00942); `ORA-01408`;
- data-dependent failures: `NOT NULL` on rows with nulls, 01442, type changes (01439, 01440), constraint validation;
- privileges, editions, and `RENAME` by a user who does not own the schema (03001);
- partitions, LOBs, nested tables, clusters, index-organized and external tables, materialized views, private temporary tables, database links, users, tablespaces;
- PL/SQL bodies and their compilation;
- names outside ASCII;
- 23ai syntax the parser rejects (`ALTER TABLE IF EXISTS`, `CREATE TABLE ... OF`, `CREATE JSON COLLECTION TABLE`, `CREATE BLOCKCHAIN TABLE`, `PERIOD FOR`, `ROW ARCHIVAL`, `IF NOT EXISTS` on synonyms, routines, types, and triggers);
- the recycle bin's contents (`FLASHBACK`, `PURGE` of a name).

## 8. Budget and tests

- **Code:** about 1.3k lines of lowering, 350 of names and snapshot, and 300 of walker and plumbing. `targets.go` gains `walks` and `OtherSpecs` in its group key. Growth past this is a decision for a person.
- **Oracle:** `walktest.Judge` against `gvenzl/oracle-free:23.26.3-slim-faststart`, in the container `oracle/review`'s tests already start.
  - Oracle commits DDL as it runs, so the engine implements `walktest.Digester` over `ALL_OBJECTS`, `ALL_TAB_COLS`, `ALL_CONSTRAINTS`, `ALL_INDEXES`, and `ALL_TRIGGERS`.
  - The change runs as an admin user with `CURRENT_SCHEMA = SHOP`, as Bytebase runs it. `RENAME` cases run as SHOP, which owns the objects.
  - The connect, the container's start, and every statement run under a deadline, so a hung container fails the test instead of stalling CI (it hung once, in #493's 4829c2fb run).
  - The test sync mirrors Bytebase's for every space the walk reads: object tables with `SkipDump`, invisible columns listed, `SYS_C` checks dropped, indexes of listed tables and materialized views only, and triggers of the schema on its own tables and views.
- **Curated cases:** one per row of §5 and per fact in the Appendix, including the wrong kinds, the synonym-following forms, and `IF [NOT] EXISTS`.
- **Digest:** the dictionary's views take hundreds of milliseconds a query (`DBA_CONSTRAINTS` most), so the digest reads `SYS.OBJ$`, `COL$`, and `CON$` by owner number, in a few. It leaves out columns of generated names: a rejected `ADD (b NUMBER CONSTRAINT nn NOT NULL)` of a taken `nn` fails with 02264 and leaves `b` behind as an unused hidden `SYS_C..._$` column (verified), which no claim names.
- **Sweep:** templates over a colliding name pool, mixing case and quoted names:
  - tables, views, sequences, synonyms, indexes named like tables, and constraint names;
  - `IF [NOT] EXISTS`, `RENAME`, and triggers;
  - `CURRENT_SCHEMA` switches.

  CI runs it at `WALK_SWEEP_N`. Precision must hold, and recall is logged. The sync reads the views Bytebase reads, a second or so, so each base is synced once and serves ten changes, each on the base made again: made again, it is the same but for the names Oracle generates, which the walk never claims. Two workers run at once on schemas of their own. At `WALK_SWEEP_N=500` it takes about 150 s and finds about half of the statements the server rejects (2860 of 5584 at seed 20261028), with no false claim.
- **Coverage:** `walktest.Coverage` over every row of §5. **Unit tests:** the snapshot's authority per space, and `OtherSpecs` in the group key.

## 9. Decisions

Decided 2026-10-28:

- **O1. Table and column absence from the sync: trusted.** Bytebase's sync lists object and XMLType tables and invisible columns (#21674), and the user decided that metadata Bytebase syncs is taken as what its current sync produces, with no signal for metadata synced before. The columns of a table the sync may not list in full stay Unknown (§3).
- **O2. DML and triggers the sync cannot show: (a).** DML goes on as §6 says, with the accepted assumption that no other schema's trigger fires on a change to the schema's tables, directly or through a cascade. (b), ending the walk at every DML statement, would end it at the first data fix in most changes.
- **O3. The assumptions of §3: accepted.** The first and third are D6. Without the first, a view the admin cannot see is one it cannot drop either: `DROP VIEW` then fails with 00942, as the walk claims.
- **O4. `CURRENT_SCHEMA` on the execution connection: closed.** Bytebase's connector sets it on every session (#21672).
- **O5. Versions before 23ai: (a).** The curated cases ran once on `gvenzl/oracle-xe:21-full-faststart` for precision, on 2026-10-28: no claim failed, and the walk found 145 of the 178 statements 21c rejects, which counts the `IF [NOT] EXISTS` forms 21c refuses as syntax errors. CI keeps 23ai. Every code here is long-standing, and the length claim holds from 11g.

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
  - `ADD (b NUMBER CONSTRAINT nn NOT NULL)` with `nn` taken fails with 02264 and leaves `b` behind as an unused hidden column `SYS_C..._$`. The same clash in a `CHECK`, a key, a table constraint, `MODIFY`, or `CREATE TABLE` leaves nothing.
- **IF [NOT] EXISTS:**
  - `CREATE TABLE IF NOT EXISTS` of a table, view, sequence, or synonym name: a no-op. `CREATE TABLE IF NOT EXISTS t (a, a)` with t present: 00957. `CREATE TABLE IF NOT EXISTS t (... CONSTRAINT pk_t ...)` with t and pk_t present: a no-op.
  - `DROP TABLE IF EXISTS` of a view or sequence, `DROP VIEW IF EXISTS` of a table or object table, `DROP SEQUENCE IF EXISTS` of a table, `ALTER TABLE IF EXISTS` of a view: no-ops. `DROP TABLE IF EXISTS` of a materialized view: 12083. `DROP TABLE IF EXISTS` of an object table drops it.
  - `CREATE INDEX IF NOT EXISTS` of a taken index name: a no-op.
  - `CREATE VIEW` / `SEQUENCE` / `SYNONYM` / `PROCEDURE` / `TRIGGER IF NOT EXISTS`, and `DROP TRIGGER` / `PROCEDURE` / `SYNONYM` / `TYPE` / `INDEX IF EXISTS`: accepted.
  - `OR REPLACE` with `IF NOT EXISTS`: 11541. `ALTER INDEX IF EXISTS ... RENAME`: 11557. `ALTER TABLE ... ADD IF NOT EXISTS`, `DROP COLUMN IF EXISTS`, `DROP CONSTRAINT IF EXISTS`: syntax errors.
- **Sync gaps:**
  - `CREATE TABLE ot OF ty` and `CREATE TABLE xt OF XMLTYPE` are absent from `ALL_TABLES` and present in `ALL_OBJECT_TABLES`. `CREATE JSON COLLECTION TABLE jct`, a nested table's storage table, and blockchain, immutable, row-archival, IOT, and temporary tables are in `ALL_TABLES`. `DROP TABLE "jct"` of the lower-case synonym a JSON collection table adds: 00942.
  - An `INVISIBLE` column has a null `COLUMN_ID` and `USER_GENERATED = 'YES'`.
  - Hidden columns with `USER_GENERATED = 'NO'`: `SYS_NC...$` (function-based indexes, object and XMLType columns), `SYS_C..._$` (unused columns), `SYS_IME_OSON_...` (JSON), `SYS_NC_OID$` and `SYS_NC_ROWINFO$` (object tables), `XMLDATA` (XMLType tables), `ORA_ARCHIVE_STATE` (row archival), `ORABCTAB_...$` (blockchain and immutable tables), and a valid-time period's `VT`, `VT_START`, `VT_END` (`PERIOD FOR vt`) or `VT2` (`PERIOD FOR vt2 (s, e)`). DDL on them: `DROP COLUMN vt_start`, `RENAME COLUMN vt`, `MODIFY`, `SET UNUSED` fail with 55603; `ADD (vt_start ...)` with 01430; `RENAME COLUMN vt2 TO zz` renames it; `ORA_ARCHIVE_STATE` fails with 38398; `ORABCTAB_HASH$` with 05715; `DROP COLUMN xmldata` of an XMLType table succeeds.
  - A valid-time period adds a check constraint named after it (`VT8cf288`), whose condition names its hidden columns (`"VT_START" < "VT_END"`, `"S"<"E" AND "ZZ">0`).
  - `ALL_TAB_COLS` shows no column of a nested table's storage table. On it, `DROP COLUMN` fails with 12982, `ADD` with 22915, `TRUNCATE` with 22812, and `DROP TABLE` with 22914. On an object table, `DROP COLUMN` of any name fails with 12981 and `ADD` with 22856. On a JSON collection table, `DROP COLUMN etag` succeeds, `DROP COLUMN data` fails with 00904, and `ADD` with 03001.
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
  - A function called without parentheses runs: `VALUES (f)`, `SET a = f`, `INSERT ... SELECT f FROM dual`, `WHERE ... AND f = 1`, `pkg.g`, and `shop.f`. A column wins over a function of its name (`WHERE f = 1` on a table with a column `F`). `CURRENT_DATE`, `CURRENT_TIMESTAMP`, `LOCALTIMESTAMP`, `DBTIMEZONE`, `SESSIONTIMEZONE`, and `ORA_ROWSCN` win over schema functions of their names; a function named `UID` cannot be created (04050). `sq.NEXTVAL` reads the sequence.
  - In `CREATE INDEX`, a name listed alone is a column: `(fx)`, `(shop.fx)`, `(pkg.g)` fail with 00904 and call nothing, while `(a + fx)` calls `fx`.
  - `CREATE FUNCTION IF NOT EXISTS` of a function's name, or of a table's, does nothing.
  - A quoted name of a function of a syntax of its own is a plain call: `"EXTRACT"(1)` calls a schema function `"EXTRACT"`, while `"UPPER"(1)` still takes the builtin.
- **Nested tables:** `DROP TABLE parent` and `ALTER TABLE parent DROP COLUMN c` drop the storage table `CHILD` of `NESTED TABLE c STORE AS child`; `SET UNUSED COLUMN c` keeps it.
- **PL/SQL with errors:** through go-ora, `CREATE OR REPLACE PROCEDURE` with an unknown call, `CREATE OR REPLACE PACKAGE BODY` without a spec, a trigger with an unknown call, and `CREATE OR REPLACE FORCE VIEW` on a missing table return no error and leave `INVALID` objects. `CREATE OR REPLACE VIEW` with an unknown column fails with 00904.
- **Namespaces:** a LOB segment `STORE AS lseg` beside a table `LSEG`: accepted. Quoted `"t"` and unquoted `T` are two tables; `DROP TABLE t2x` does not find `"T2x"`.
