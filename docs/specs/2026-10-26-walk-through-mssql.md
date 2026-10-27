# WALK_THROUGH for MSSQL: v1 scope

Status: decided 2026-10-27; §10 records the decisions. Step 4 of [the engines design](2026-10-22-walk-through-engines-design.md) §7, which its §5.4 sketches; [the MySQL v1 scope](2026-10-22-walk-through-mysql.md) is the model. Every fact here was verified on SQL Server 2022 CU27 (16.0.4305.5), the image `mssql/review` already uses (`mcr.microsoft.com/mssql/server:2022-latest`), and through go-mssqldb, Bytebase's fork included (Appendix). The Bytebase facts are from origin/main 6589af9dcb.

## 1. Environment

**How Bytebase runs a change** (`backend/plugin/db/mssql/mssql.go`, `Execute`):
- It splits the text into batches at `GO` lines, as sqlcmd does (`plugin/parser/tsql/batch`), and sends each batch to the server. `GO n` sends the batch n times.
- **Transaction mode** is the default: one transaction on one connection, so the session (temporary tables, `USE`, `SET` options) carries over from batch to batch. A `-- txn-mode = off` line at the top switches to **auto-commit mode**: each batch goes through the connection pool, and go-mssqldb resets the session before reusing a connection, so every batch starts with no temporary tables, on the target database, with default options (verified with Bytebase's fork).
- The change stops at the first batch that returns an error. Inside a batch, SQL Server goes on after many errors (`DROP TABLE missing` is level 11 and the next statement runs), but the change fails all the same.
- The session starts on the target database (the DSN names it), with no temporary tables, `QUOTED_IDENTIFIER` and `ANSI_NULLS` on, and the login's default schema.

**Batches:**
- A batch compiles as a unit. A compile error, such as a column a statement of the same batch adds, fails the whole batch before any of it runs. This only adds failures.
- `CREATE`, `ALTER`, and `CREATE OR ALTER` of a view, procedure, function, or trigger, and `CREATE SCHEMA`, must be the first statement of their batch (111). Comments and an empty `;` before them are fine; any statement, `SET NOCOUNT ON` included, is not.
- **A procedure's or trigger's body runs to the end of its batch.** `CREATE PROCEDURE p AS SELECT 1; CREATE TABLE t (a int)` creates p, whose body holds the `CREATE TABLE`, and creates no table, with or without `BEGIN ... END`. omni's parser ends the body at the first statement and lists the rest as statements of their own, so the walk must skip them (§6).
- A view, function, or `CREATE SCHEMA` followed by another statement in its batch fails the batch (156).
- A `GO` is a separator for Bytebase only alone on its line: blanks, `GO` in any case, and optionally blanks and a positive count, outside a string or block comment. omni also parses `SELECT 1 GO` and `GO -- note` as `GO`, which Bytebase sends to the server as SQL.

**Database and schema:**
- A Bytebase database is one SQL Server database with many schemas.
- `USE d` switches the database for the rest of the batch, and in transaction mode for the rest of the change.
- An unqualified creation goes to the default schema. An unqualified reference finds the object in the default schema, else in `dbo` (verified with a user whose default schema is s: `CREATE TABLE x` lands in s; `DROP TABLE`, `ALTER TABLE`, `CREATE INDEX`, `TRUNCATE`, `INSERT`, and `sp_rename` take dbo.x when s has none, and s.x when both have it).
- Bytebase neither sets nor syncs the default schema (D2). It is `dbo` for a sysadmin login, which maps to the dbo user.
- `CREATE/ALTER/DROP DATABASE` cannot run in a transaction (226, 574), so in Bytebase's default mode they always fail.

**Temporary tables:**
- A name that starts with `#` is a session's temporary table, even with a schema in front (`dbo.[#t]` is `#t`).
- A `##` name is a global temporary table, which other sessions can create and drop.
- Temporary table names compare under tempdb's collation, which is the server's: `#T` collides with `#t` in a case-sensitive database on a case-insensitive server.
- A foreign key on a temporary table is skipped with a warning, and its constraint names live in tempdb, where other sessions' names collide.

**Server settings the walk reads:** none. The collation is unknown (§2.1) and the version is not an input.

## 2. Spaces

| Space | Parent | Holds | Claims |
| --- | --- | --- | --- |
| schema | root | schemas | exists (2714), missing (2760 for a creation in it, 15151 for `DROP`/`ALTER SCHEMA`), not empty (3729) |
| object | schema | every `sys.objects` name: tables, views, procedures, functions, sequences, synonyms, DML triggers, and PK, UQ, FK, CK, and DF constraints; also CLR modules, legacy `CREATE DEFAULT`/`RULE` objects, queues, and edge constraints | exists (2714; 15335 for `sp_rename`), of another kind (3705, 2010), missing (208; 15225 for `sp_rename`) |
| relation | schema | tables, views, and sequences: the kinds the sync lists completely | missing or of another kind (3701/3705 for `DROP`, 4902/4909 `ALTER TABLE`, 1088 `CREATE`/`ALTER INDEX`, 4701/4708 `TRUNCATE`, 1767/1768 a foreign key's target, 8197 `CREATE TRIGGER ... ON`, 15151 `ALTER SEQUENCE`, 208/2010 `ALTER VIEW`) |
| column | schema.table | columns | exists (2705; 15335), missing (4924, 1752, 15248) |
| index | schema.table | indexes, the PK's and UQ's included | exists (1913; 15335), missing (3701, 7999, 2727, 15248), backs a constraint (3723) |
| constraint | schema.table | the table's PK, UQ, FK, CK, and DF constraints | missing (3728, 3733, 4917) |
| primary | schema.table | whether the table has a primary key | exists (1779) |
| temporary | root | the session's `#` tables | exists (2714), missing (3701, 4902, 1088, 208) |
| type | schema | alias and table types | exists (219), missing (218) |

- **Kinds** in the object space: table, view, sequence, procedure, function, trigger, synonym, pk, uq, fk, ck, df.
- **Constraint names are schema objects:** `CREATE TABLE pk_t` fails when another table's primary key is `pk_t`, and so does `ADD CONSTRAINT t` when a table is named t. A PK or UQ constraint also takes its name in its table's index space (`CREATE INDEX pk_t ON t` fails with 1913).
- **Index names are per table and not objects:** `CREATE TABLE ix` and `ADD CONSTRAINT ix CHECK` succeed beside an index ix; `CREATE INDEX pk_t ON r` succeeds beside t's pk_t.
- **Statistics share the index namespace** (1913) and are not synced or modeled.
- **Types are a space of their own:** `CREATE TABLE ty` succeeds beside a type ty. A table type's internal `TT_<name>_<hex>` object is in `sys`, apart from user schemas.
- **A trigger's schema is its table's** (2103).
- **Length:** an identifier longer than 128 characters fails with 103 before anything is looked up, `DROP ... IF EXISTS` included; a temporary table name longer than 116, with 193. The walk claims this for ASCII names, from the text alone.

### 2.1 Name comparison: the collation is unknown

Names of one space compare under the database's collation (verified on four), temporary table names under the server's, and database names under the server's. Bytebase records none of them: `DatabaseSchemaMetadata.Collation` stays empty for MSSQL.

The engines design's §4.3 rule assumes that two ASCII names a collation equates differ only in case. **That is false on SQL Server.** A sweep of every collation (Appendix) found these other equalities:
- `v` = `w`: Estonian and Finnish_Swedish, accent-insensitive;
- `i` = `y`: Frisian_100, accent-insensitive;
- `f` = `g`: Welsh_100, even case- and accent-sensitive;
- Hungarian doubled digraphs: `ccs` = `cscs`, `ddz` = `dzdz`, `ddzs` = `dzsdzs`, `ggy` = `gygy`, `lly` = `lyly`, `nny` = `nyny`, `ssz` = `szsz`, `tty` = `tyty`, `zzs` = `zszs`. This holds in Hungarian_CI_AS and CS_AS, all four Hungarian_100 collations, and both SQL_Hungarian_CP1250 ones, case-sensitive included. It does not hold in Hungarian_CI_AI, Hungarian_CS_AI, or the Hungarian_Technical collations.

So a database in Hungarian_CI_AS refuses `CREATE TABLE szsz` beside a table ssz, and `DROP TABLE SZSZ` drops ssz (verified).

The walk therefore keys every name so that any two names some collation may equate share a key:
- **A plain name** has only ASCII letters, digits, and `_`, apart from a temporary table's leading `#`. Its key is built in three steps:
  1. lowercase it and map `w` to `v`, `y` to `i`, `g` to `f`;
  2. keep the first occurrence of each character;
  3. append the last character.

  For example, `users` keys as `user|s`, `user` as `user|r`, and `ssz` and `szsz` both as `sz|z`. Lowercasing and the three letter maps cover the letter equalities, and doubling a digraph adds no new first occurrence and keeps the last character, so the key is the same under every equality the sweep found.
- **A marked name** is any other name: one with a character outside ASCII, a space, punctuation, or a bracket. Under an accent-insensitive collation `café` = `cafe`; in Maori_100 `-` is ignorable, and in Mohawk_100 `:`. Its key is its exact spelling, marked.

What the walk claims under these keys:
- **Exists, or is of another kind:** only when the key's Present fact has the statement's exact spelling. A Present fact keeps the spelling it was synced or created with.
- **Missing:** only when the key is Absent, and its parent holds no marked name, synced or created by the change. A key never claims a marked name missing.
- **Two spellings of one key:** the key is Unknown after the change creates, drops, or renames either, as MySQL's walk does for 5.7 triggers.
- **Effects across the mark:**
  - creating a marked name may add any name of its parent's space;
  - dropping or renaming a marked name may remove any;
  - dropping or renaming a plain name may remove the parent's marked names: under an accent-insensitive collation, `DROP TABLE cafe` drops `café`.
- **Within one statement:** two names collide (2705, 8168) only when spelled identically.
- **Database names:** a qualifier names the target only when spelled as Bytebase spells the database.

Costs:
- Names that differ only in repeated letters or in v/w, i/y, f/g share a key: `col1` and `col11`, `a` and `aa`. A lookup of one beside the other is Unknown.
- Bytebase recording the collation would remove this cost (decision M2).

## 3. Authority

Present is trusted for every name the sync lists. The sync (`backend/plugin/db/mssql/sync.go`, `SyncDBSchema`) lists the schemas of `INFORMATION_SCHEMA.SCHEMATA`, except the fixed-role schemas (`db_owner`, `db_accessadmin`, ...), `guest`, `INFORMATION_SCHEMA`, and `sys`.

**Absent is trusted for:**
- **schemas:** every name but the excluded ones, which stay Unknown. guest and db_owner can hold tables (verified).
- **relations** (tables, views, sequences) in a synced schema. `sys.tables`, `sys.views`, and `sys.sequences` are read unfiltered; every view has a `sys.sql_modules` row, an encrypted one with a NULL definition (verified). A schema that certainly does not exist holds none.
- **columns** of a synced table. `sys.columns` is read unfiltered, so hidden period and graph columns are listed too.
- **index names** of a synced table:
  - the PK and UQ constraints' indexes (`getKeys`);
  - every other index with a key column (`getIndexes`). Clustered and nonclustered columnstore, XML, and hypothetical indexes all have `sys.index_columns` rows (verified).

  Not a view's indexes: the sync attaches indexes to tables only.
- **constraint names** of a synced table: PK, UQ, FK, CK, and DF (`ColumnMetadata.DefaultConstraintName`). The exception is a graph edge table (a column named `$edge_id_...`), whose edge constraints the sync skips.
- **whether a synced table has a primary key.**
- **temporary tables** at the session's start.

Each "missing" claim also needs the parent to hold no marked name (§2.1).

**Never trusted Absent:**
- **the object space.** The sync skips objects that each hold a name in it (each one probed made `CREATE TABLE` of its name fail with 2714):
  - synonyms;
  - CLR procedures, functions, and triggers;
  - disabled triggers;
  - procedures and functions marked `microsoft_database_tools_support` (the SSMS diagram objects);
  - legacy defaults and rules;
  - queues;
  - edge constraints.

  So "already exists" is claimed only on Present. DML targets and procedures, functions, triggers, and synonyms are claimed missing only from the change's own facts, as after a `DROP` of the name.
- **types and statistics,** which are not synced;
- **the excluded schemas and other databases;**
- **keys the name rule leaves Unknown.**

**Accepted assumptions:**
- **A1** (engines design D6): the admin login sees every object of its database; catalog views list only what the login may see.
- **A2:** no database DDL trigger creates, drops, or renames objects. The sync lists none. They fire on the change's DDL (verified), and one that rolls back only adds failures.
- **A3:** no CLR DML trigger changes the schema. The sync lists only T-SQL triggers; §6 handles those.

## 4. Statement semantics the lowering follows

**Multi-object `DROP`** (`TABLE`, `VIEW`, `PROCEDURE`, `INDEX`, ...) runs its items in order and is not atomic:
- A missing item fails the statement (3701), and the others are still dropped: `DROP TABLE r, missing` and `DROP TABLE missing, r` both drop r.
- `DROP TABLE r, r` drops r and fails on the second.
- An item of another kind fails the statement and drops nothing (3705).
- `IF EXISTS` skips a missing item but not one of another kind (3705). `DROP INDEX IF EXISTS` of a PK's index fails with 3723.
- The walk uses `walk.Passes`: each item is checked against the state before the statement, and a check on a name another item changes is left out. A name written twice in the same spelling fails (3701).

**`ALTER TABLE` takes one action:**
- an `ADD` list of columns and constraints, whose constraints see the list's columns in any order;
- a `DROP` list of `COLUMN` and `CONSTRAINT` items, which is atomic: `DROP COLUMN b, COLUMN missing` drops nothing;
- `ALTER COLUMN`, or another single action.

`ADD` and `DROP` do not mix (a syntax error). A name twice in one statement fails (8168; 2705 for columns), and so do two primary keys (8110).

**Generated constraint names:** `PK__<table>__<hex>`, `UQ__...`, `FK__<table>__<column>__<hex>`, `CK__...`, `DF__...`, with the parts truncated. They become prefix May-effects (`pk__`, `uq__`, `fk__`, `ck__`, `df__`) in the object space, and for PK and UQ in the table's index space.

**`sp_rename`:**
- **Types:** OBJECT (the default) renames an object, COLUMN a column (`'[schema.]table.column'`), and INDEX an index (`'[schema.]table.index'`). `@objname` resolves like any reference.
- **The new name is taken literally:** `'dbo.t2'` makes a table named `dbo.t2` in dbo, and `'[t9]'` keeps the brackets.
- **Missing:** 15225 for OBJECT, 15248 for COLUMN and INDEX, which an object of another kind also gets.
- **In use:** 15335, in the object, column, or index space.
- **Renaming a table** keeps the names of its constraints, indexes, and triggers.
- **A PK's index and its constraint** rename together: renaming the PK's index renames the constraint, and the other way round.
- **A case-only rename** of a name to itself succeeds.
- **Errors:** all three are level 11, so the change fails.

**`ALTER SCHEMA s TRANSFER x`** moves x's constraints with it (verified), so their names must be free in s too. It is not covered (§7).

**Dependency failures happen, but the walk makes no claim about them:**
- `DROP TABLE` of a table a foreign key references (3726);
- `DROP COLUMN` of a column an index, default, or constraint uses (5074);
- `DROP SCHEMA` of a schema that holds only a type;
- `TRUNCATE` of a referenced table.

Views are not dependencies: dropping a table a view reads succeeds.

**Triggers run DDL:**
- An `AFTER INSERT` trigger that runs `CREATE TABLE` creates the table.
- A cascaded `ON DELETE CASCADE` fires the child table's `AFTER DELETE` trigger.
- An `INSTEAD OF` trigger on a view fires on DML to the view.
- A function cannot run DDL (443).

## 5. Statements in v1

**Notation:**
- "x Present" or "x Absent" is a check that yields a finding only when the ledger certainly says the opposite, under §2.1's name rule.
- A reference resolves as in §1: in the schema it names, or unqualified in `Target.DefaultSchema` and then dbo.
- **Unqualified with the default schema unset** (D2): the reference is Unknown. A creation of that name, or a drop or rename of it, may then have happened in any schema, so its key becomes Unknown across schemas. So it is when the default schema is set but may not exist.
- **Lookup with a default schema other than dbo:** the default schema when it certainly holds an object of the name, dbo when it certainly holds none, else Unknown.
- **A three-part name** is judged only when its database is the target, spelled as Bytebase spells it; a four-part name never.

| Statement | Checks | Effects |
| --- | --- | --- |
| a module statement (§1) not first in its batch | fails (111) (M4) | |
| `GO` alone on its line | none | ends the batch (§6) |
| `USE d` | none: databases are not synced | current database: the target for its spelling, another for a name of another class, else unknown. Schema statements follow it too: in another database they change nothing the walk keeps, in an unknown one they forget the schema |
| `CREATE SCHEMA s [AUTHORIZATION o]` | s Absent (2714) | s Present, every space under it Absent |
| `DROP SCHEMA [IF EXISTS] s` | s Present (15151); s certainly holds an object, synced or made or moved there by the change (3729) | s and everything under it Absent |
| `CREATE TABLE t (...)` | schema Present (2760); t Absent in the object space (2714), or as a temporary table (2714); identifiers not too long (103, 193); within the statement: a column twice (2705), a constraint or index name twice (8168), two primary keys (8110), a constraint named like the table (2714); each constraint name Absent (2714); a foreign key's target a table (1767, 1768), its own table included, and not for a temporary table | t Present (table), its columns exact, its named indexes and constraints Present, its primary key known; generated names May-added |
| `SELECT ... INTO t` | schema Present (2760); t Absent (2714), or `#t` as a temporary table (2714) | t Present (table), its columns Unknown |
| `ALTER TABLE t ...` | t Present and a table, or a temporary table (4902, 4909) | per action: |
| — `ADD` list | each column Absent (2705) and not twice; each constraint name Absent (2714) and not twice (8168); a primary key when t certainly has one (1779), or two (8110); a PK or UQ name Absent in t's index space (1913); a foreign key's target a table (1767, 1768) | as in `CREATE TABLE` |
| — `DROP` list of `[COLUMN] c [IF EXISTS]`, `[CONSTRAINT] n [IF EXISTS]` | each c Present (4924) and each n a constraint of t (3728, 3733), unless `IF EXISTS`; no name twice (8168) | c Absent; n Absent, with its index and the primary key when it is the PK; a constraint of unknown kind May-removes t's index of its name and its primary key |
| — `ALTER COLUMN c ...` | c Present (4924) | none |
| — `[WITH CHECK] CHECK`/`NOCHECK CONSTRAINT n` | n a constraint of t (4917), unless `ALL` | none |
| — `ENABLE`/`DISABLE TRIGGER` | none | `ENABLE`: a trigger may fire on t (§6) |
| — `SET (...)`, `REBUILD` | none | none |
| — `SWITCH`, `ADD`/`DROP PERIOD`, `SET (SYSTEM_VERSIONING ...)`, anything else | | ends the walk |
| `CREATE TABLE` with `AS NODE`/`EDGE`/`FILETABLE`, a period, `SYSTEM_VERSIONING`, `LEDGER`, or `MEMORY_OPTIMIZED` | | ends the walk |
| `ENABLE`/`DISABLE TRIGGER ... ON t` | none | `ENABLE`: a trigger may fire (§6); on a database or the server, ends the walk |
| `DROP TABLE [IF EXISTS] a, ...` | each a Present (3701), and a table under `IF EXISTS` too (3705); a name twice (3701) | a Absent with its columns, indexes, constraints, and primary key; its constraints' and triggers' names leave the object space: the known ones Absent, the rest May-removed |
| `TRUNCATE TABLE t` | t Present and a table (4701, 4708), or a temporary table | none |
| `CREATE [UNIQUE] [CLUSTERED \| NONCLUSTERED] [COLUMNSTORE] INDEX i ON t` | t Present and a table or view (1088, 1914), or a temporary table; i Absent in t's index space (1913); with `DROP_EXISTING = ON`, i Present (7999) | i Present |
| `DROP INDEX [IF EXISTS] i ON t, ...`, `DROP INDEX t.i` | each t Present (3701); i Present (3701), unless `IF EXISTS`; i not a constraint's index (3723), under `IF EXISTS` too | i Absent |
| `ALTER INDEX i \| ALL ON t ...` | t Present (1088); i Present (2727) | none |
| `CREATE VIEW v` / `CREATE OR ALTER VIEW v` / `ALTER VIEW v` | schema Present (2760); `CREATE`: v Absent (2714); `CREATE OR ALTER`: v a view if it exists (2010); `ALTER`: v Present and a view (208, 2010) | v Present (view) |
| `DROP VIEW [IF EXISTS] v, ...` | each v Present (3701), and a view under `IF EXISTS` too (3705); a name twice (3701) | v Absent with its indexes; its triggers May-removed |
| `CREATE` / `CREATE OR ALTER` / `ALTER` `PROCEDURE` / `FUNCTION p` | schema Present (2760); `CREATE`: p Absent (2714); `CREATE OR ALTER`: p of its kind if it exists (2010); `ALTER`: p Present (208, from the change's facts alone) and of its kind (2010) | p Present of its kind; a procedure's body runs to the end of its batch, so the walk skips to the next `GO` (§6) |
| `CREATE` / `CREATE OR ALTER` / `ALTER TRIGGER tr ON t` (DML) | t Present and a table or view (8197); tr's schema t's (2103); `CREATE`: tr Absent (2714); `ALTER`: tr Present (208, from the change's facts alone) and a trigger (2010) | tr Present (trigger); a trigger may fire on t; the walk skips to the next `GO` |
| `DROP PROCEDURE` / `FUNCTION` / `TRIGGER` / `SEQUENCE` / `SYNONYM [IF EXISTS] x, ...` | each x Present (3701), from the change's facts alone except for sequences, unless `IF EXISTS`; of the statement's kind, under `IF EXISTS` too (3705) | x Absent |
| `CREATE SEQUENCE s`, `CREATE SYNONYM s FOR ...` | schema Present (2760); s Absent (2714) | s Present of its kind |
| `ALTER SEQUENCE s` | s Present and a sequence (15151) | none |
| `CREATE TYPE ty`, `DROP TYPE [IF EXISTS] ty` | schema Present (2760); ty Absent (219) / Present (218), from the change's facts alone | ty Present / Absent |
| `EXEC sp_rename 'x', 'n' [, 'OBJECT']` of a table, view, procedure, function, sequence, or synonym | x Present (15225), from the change's facts alone; n Absent in the object space (15335) | x and its columns and indexes move to n; a constraint, a trigger, or an object of unknown kind ends the walk |
| `EXEC sp_rename 't.c', 'n', 'COLUMN'` | t Present and a table (15248); c Present (15248); n Absent (15335) | c moves to n |
| `EXEC sp_rename 't.i', 'n', 'INDEX'` | t Present (15248); i Present (15248); n Absent (15335) | i moves to n; an index that backs a constraint ends the walk |
| `INSERT`, `UPDATE`, `DELETE`, `MERGE` | the target Present (208), from the change's facts alone; never a CTE, an alias, or a table variable | none; ends the walk when a trigger may fire (§6) |
| `IF cond s [ELSE s]` (M5) | none | each branch's statements as May-effects (`walk.Builder.Optional`) |
| `SELECT` without `INTO`, `PRINT`, `DECLARE`, `SET @v`, `SET` options other than §6's, `BEGIN TRAN`, `COMMIT`, `SAVE TRAN`, `GRANT`/`DENY`/`REVOKE`, `UPDATE STATISTICS`, `RAISERROR`, `THROW`, `WAITFOR`, `EXEC sp_addextendedproperty`/`sp_updateextendedproperty`/`sp_dropextendedproperty` | none | none |

Notes:
- `sp_rename` is lowered only when unqualified or qualified with `sys`, with literal arguments, positional or named. Its new name is used as written, so a name with a dot or bracket is a marked name.
- `CREATE TABLE` has no `IF NOT EXISTS` in T-SQL, and `CreateTableStmt.IfNotExists` is never set.
- `ALTER` and `CREATE OR ALTER` of a module share an AST node (`OrAlter`). The lowering reads the first keyword at `Loc.Start`.

## 6. Where the walk stops, and what it skips

**Ends the walk:**
- any statement not in §5;
- `EXEC` of anything else, `EXEC (@sql)`, `sp_executesql`, and `INSERT ... EXEC`;
- `WHILE`, `BEGIN ... END` outside an `IF`, `TRY ... CATCH` (a caught failure fails nothing), `GOTO`, and `RETURN`;
- `ROLLBACK`;
- `EXECUTE AS`, `REVERT`, and `SETUSER`, which change the default schema;
- `SET NOEXEC`, `PARSEONLY`, `FMTONLY`, and `SHOWPLAN_*` on, after which statements no longer run (verified);
- `SET QUOTED_IDENTIFIER OFF`, after which the server reads the text otherwise than omni;
- `CREATE`/`ALTER`/`DROP DATABASE`, which fail in Bytebase's default mode;
- `CREATE TRIGGER ... ON DATABASE` or `ALL SERVER`;
- `GO n` with n > 1, at the `GO`: the batch before it was judged for its first run;
- a `GO` that is not alone on its line (§1), and a statement after a view, function, or `CREATE SCHEMA` in the same batch (156).

**DML ends the walk after its checks when a trigger may fire:**
- a synced trigger exists on any table or view;
- the change created a trigger or enabled one;
- `OtherSpecs` is set;
- or the target may lead to another database, whose triggers are not synced: a name in another database or one the walk cannot place, a CTE or a FROM item the walk does not resolve, and an object the sync does not list as a table or view, which may be a synonym.

A trigger can run DDL, and a cascaded action fires the triggers of other tables, so the trigger need not be on the target. Without any trigger, DML changes no name. With the default schema unset, every unqualified DML target is one the walk cannot place, so the walk ends there.

**Skipped:** after a `CREATE`, `ALTER`, or `CREATE OR ALTER` of a procedure or trigger, every statement up to the next `GO` is part of its body. The walk neither judges those statements nor applies them, and goes on after the `GO`.

**At each `GO`:**
- The session's temporary tables become Unknown.
- A current database other than the target becomes unknown.

Both hold in either of Bytebase's modes (M3).

**`OtherSpecs`:** an all-Unknown snapshot, and DML ends the walk, as on the other engines.

## 7. Not covered (declined by scope)

Review asks for these are declined with this list:

- **dependency failures:** 3726, 5074, 3729 for a schema that holds only a type, and `TRUNCATE` of a referenced table;
- **column references** in index, constraint, and foreign key definitions (1911, 1776, ...), and in views, defaults, computed columns, and queries;
- **object references** inside queries and module bodies;
- **`DEFAULT ... FOR c`:** a missing c (1752), since the parser does not keep the column; c with a default (1781), since legacy bound defaults are not synced;
- **a second clustered index;** **a system-versioned table** dropped (13552) or altered;
- **`ALTER SCHEMA ... TRANSFER`;** **`sp_rename`** of a constraint, a trigger, a constraint's index, or a database;
- **`ALTER TABLE ... SWITCH`,** periods, partitions, filegroups, and memory-optimized options;
- **statistics,** full-text, XML schema collections, and Service Broker objects;
- **types and expressions;** **data-dependent failures;** **privileges and ownership** (`AUTHORIZATION` users, `EXECUTE AS` clauses);
- **other databases;** **the global `##` temporary tables;**
- **non-ASCII identifiers longer than 128 characters;**
- **`CREATE OR ALTER TRIGGER`** of another table's trigger (2110);
- **SQL Server versions before 2022**, which were not probed.

## 8. Contract changes

`review.Target.DefaultSchema string`: the default schema of the login that runs the change, as `SELECT SCHEMA_NAME()` returns it over Bytebase's admin connection, which also runs the change.
- Empty means unknown.
- MSSQL reads it. The engines design's D2 has Snowflake read a current schema the same way later.
- The field goes in a commit of its own (`feat(review): ...`), with no change to the other engines.

Bytebase follow-ups, outside this PR:
- read `SCHEMA_NAME()` and the database's collation in the MSSQL sync. `DatabaseSchemaMetadata.Collation` exists but stays empty. For Azure SQL Database and contained databases, the collation names compare under is the catalog collation;
- set `Target.DefaultSchema` and `Target.OtherSpecs`;
- bump the omni pin.

`groupTargets` gains `DefaultSchema` and `OtherSpecs` in its key, each pinned by a test.

## 9. Budget and tests

**Code:**
- about 1.5k lines of lowering (`walk_lower*.go`). The table is as broad as MySQL's, plus `sp_rename` and `IF`;
- 450 of names and snapshot (`walk_names.go`, `walk_snapshot.go`);
- 350 of walker, batches, and plumbing (`walk_through.go`, `targets.go`).

As built: 1,717 lines of lowering, 434 of names and snapshot, and 700 of walker, batches, reference resolution, and the class and spelling bookkeeping of §2.1 (`walk_through.go`), about 25% over the budget in all. The overrun is mostly that bookkeeping, which the budget did not foresee, and the per-kind DROP and module rows.

**Oracle** (`walktest.Judge`), modeled on Bytebase's transaction mode:
- One connection runs one transaction. Each statement runs as its own batch after a `SAVE TRANSACTION`, and a rejected one is rolled back to it.
- A rejection that ends the transaction is handled by replaying the accepted prefix in a new one. 2714, 1913, and 3728 do (`XACT_STATE()` 0).
- `Holds` reads the catalog on that connection, so it sees the uncommitted state.
- `Exec` of a `GO` runs nothing.
- The test sync mirrors Bytebase's queries and exclusions, citing `sync.go`.

**Batch tests:** the two batch facts (111 and a procedure body that swallows statements) run whole batches on the server.

**Curated cases:**
- one per row of §5 and per fact in the Appendix;
- each on two databases: `SQL_Latin1_General_CP1_CI_AS` (the server's default) and `Latin1_General_100_CS_AS`;
- one with `DefaultSchema` s, run as a database user whose default schema is s (`EXECUTE AS USER`), and one with it unset.

**Sweep:**
- templates over a pool of colliding names: case variants, `ssz`/`szsz`, marked names (`[a b]`, `[é]`), temporary tables, two schemas, `GO`, `USE` of another database, `sp_rename`, and `IF`;
- on the CI and CS databases, and on a Hungarian_CI_AS one, one after another, with the default schema unset in a third of the changes;
- `WALK_SWEEP_N` changes in all, spread over the three;
- precision must be 100%, and recall is logged.

At 1,500 changes (seed 7) the sweep found 2,410 of 5,450 rejections on CI, 2,505 of 5,554 on CS, and 2,434 of 5,619 on Hungarian_CI_AS, with no false finding. Most misses are by design: unqualified names in the third of the changes without a default schema, names in another database, object-space absence the sync cannot vouch for, and the walk's end at DML that may fire a trigger.

**Coverage:** `walktest.Coverage` over every row of §5. **Unit tests** for the name key and the snapshot.

## 10. Decisions

Decided 2026-10-27, each as recommended: M1 (a), M2 (a), M3 (a), M4 (a), M5 (b), M6 accepted, M7 as stated.

The options as they were weighed:

- **M1. `Target.DefaultSchema`** (the D2 follow-through):
  - (a) add the field in this PR, unset meaning unknown;
  - (b) wait for Bytebase.

  Recommended: (a). Without it, every unqualified name is Unknown even after Bytebase syncs the schema.
- **M2. Collation.**
  - (a) Always key names by §2.1's rule, which holds under every equality the sweep found in any collation.
  - (b) Also read `Target.Schema.Collation` once Bytebase fills it, and key case-only (CI) or exactly (CS, BIN) for the collations whose only ASCII equality is case. That is every family but Estonian, Finnish_Swedish, Frisian_100, Welsh_100, Hungarian and SQL_Hungarian, Maori_100, and Mohawk_100.

  Recommended: (a) now. Do (b) with the Bytebase sync change, since the field is empty today and (b) would be untested code.
- **M3. Session state across `GO`.**
  - (a) At each `GO`, temporary tables become Unknown, and so does a current database other than the target. This holds in both of Bytebase's modes.
  - (b) Read the `-- txn-mode = off` directive as `ParseTransactionConfig` does, and keep the session across `GO` in transaction mode.

  Recommended: (a). It costs only claims about a temporary table across batches, and does not tie the walk to Bytebase's executor.
- **M4. A module statement not first in its batch** (111).
  - (a) Report it as a WALK_THROUGH finding, verified by the batch tests.
  - (b) End the walk there, and leave it to a parser follow-up that makes it a Syntax finding. omni's `Parse` accepts what SQL Server rejects here, which AGENTS.md says it should not; but changing the parser moves statement splitting Bytebase relies on.

  Recommended: (a). It is a certain failure that users hit when they leave out a `GO`, and it is found from the text.
- **M5. `IF`.**
  - (a) End the walk, as the engines design's §5.4 says for control flow.
  - (b) Lower `IF cond s [ELSE s]`: each branch's statements, alone or in a `BEGIN ... END`, are optional parts, with no checks and their effects weakened. A branch with a statement the walk does not lower, or with a session effect (`USE`, `SET`), still ends the walk. `WHILE` and `TRY` end it either way.

  Recommended: (b). SSMS-generated scripts guard nearly every object with `IF [NOT] EXISTS (...)`, so under (a) the walk would end at their first statement; `walk.Builder.Optional` was built for this.
- **M6. Triggers.** DML ends the walk after its checks whenever a trigger may fire, under assumptions A2 and A3. Recommended: accept. The stricter alternative, ending the walk at all DML, costs every data migration its findings after the first DML; the looser one, ignoring triggers, makes false findings (verified: a trigger creates a table).
- **M7. Oracle databases:**
  - CI and CS for the curated cases;
  - the same plus Hungarian_CI_AS for the sweep;
  - all on the one container.

  Recommended as stated. Accent-insensitive collations stay covered by the name rule's sweep (Appendix), not by the walk's oracle.

Parser follow-ups found here, outside this PR, which the walk works around:
- omni ends a procedure's or trigger's body at its first statement, while SQL Server runs it to the end of the batch. Other rules judge a body's later statements as top-level ones too: `DisallowDropObject` reports the `DROP TABLE` in `CREATE PROCEDURE p AS SELECT 1; DROP TABLE t`.
- omni parses `GO` that is not alone on its line.

## Appendix: verified facts

SQL Server 2022 CU27 (16.0.4305.5), server collation `SQL_Latin1_General_CP1_CI_AS`, through sqlcmd and go-mssqldb, 2026-10-26.

**Namespaces:**
- `CREATE TABLE` of the name of a view, procedure, function, sequence, synonym, trigger, a disabled trigger, or a PK, DF, CK, or UQ constraint fails with 2714. So does the name of a legacy rule or a queue.
- `CREATE VIEW`, `PROCEDURE`, `FUNCTION`, `SEQUENCE`, `SYNONYM`, and `TRIGGER` of a table's name fail with 2714. So do `ADD CONSTRAINT` of another table's constraint's name or of a table's name (with 1750), `SELECT ... INTO` an existing table, and `CREATE SCHEMA` of an existing schema (with 2759).
- `CREATE TABLE ty` beside a type, `CREATE TABLE ix_a` and `ADD CONSTRAINT ix_a CHECK` beside an index, and `CREATE INDEX pk_t ON r` beside t's pk_t succeed. `CREATE TYPE t` beside a table succeeds; `CREATE TYPE ty` twice fails with 219; `DROP TYPE` of a missing type fails with 218.
- `CREATE INDEX pk_t ON t`, `CREATE INDEX IX_A` beside ix_a, `ADD CONSTRAINT ix_a UNIQUE` beside index ix_a, and `CREATE INDEX st1` beside statistics st1 fail with 1913.
- A table type adds an object named `TT_<name>_<hex>` in `sys`; `CREATE TABLE dbo.TT_<name>_<hex>` succeeds. A legacy default and an edge constraint hold their names (2714).
- `CREATE TRIGGER s.tr ON dbo.r` fails with 2103; `CREATE TRIGGER` on a missing table or on a synonym fails with 8197.
- Database DDL triggers are not in `sys.objects`, their names collide among themselves (2714), and they fire on the change's DDL.

**Missing and wrong kind:**
- `DROP TABLE` / `VIEW` / `PROCEDURE` / `FUNCTION` / `SEQUENCE` / `SYNONYM` / `TRIGGER` of a missing name fails with 3701 (level 11). Of another kind it fails with 3705, which names the kind (a view, synonym, procedure, sequence, constraint, table, function). With `IF EXISTS`, a missing name is skipped, and another kind still fails with 3705.
- `DROP SCHEMA` of a missing schema fails with 15151, and of one holding a table or a type with 3729.
- `ALTER TABLE` of a missing table fails with 4902, and of a view, synonym, or procedure with 4909. `ALTER TABLE #missing` fails with 4902.
- `TRUNCATE` fails with 4701 for a missing table and 4708 for a view or synonym.
- `CREATE INDEX` fails with 1088 for a missing table, 1914 for a synonym, and 1939 for a view that is not schema-bound. `ALTER INDEX` fails with 1088 (missing table) and 2727 (missing index, level 11). `DROP INDEX` of a missing table or index fails with 3701, of a PK's or UQ's index with 3723, and so does `IF EXISTS` of a PK's index. `DROP_EXISTING = ON` of a missing index fails with 7999. `DROP INDEX t.i` works.
- `ALTER VIEW`, `PROCEDURE`, `FUNCTION`, and `TRIGGER` of a missing name fail with 208, and of another kind (a table, a procedure, a function) with 2010; `CREATE OR ALTER` of another kind fails with 2010. `CREATE OR ALTER TRIGGER` of another table's trigger fails with 2110. `ALTER SEQUENCE` of a missing name or a table fails with 15151.
- DML on a missing table fails with 208, and through a synonym it succeeds. `UPDATE` of a CTE or of an alias succeeds.
- A foreign key to a missing table fails with 1767, and to a view or synonym with 1768 (both with 1750). A self-reference in `CREATE TABLE` works.
- A creation in a missing schema (`TABLE`, `VIEW`, `PROCEDURE`, `SYNONYM`, `SEQUENCE`, `TYPE`, `SELECT ... INTO`) fails with 2760. So does one in `sys` or `INFORMATION_SCHEMA`, while guest and db_owner accept tables.

**Columns, constraints, keys:**
- `ADD a` beside a, or `ADD A` (CI), fails with 2705; `ADD x, x` and `ADD x, X` fail with 2705; `CREATE TABLE (a, a)` and `(a, A)` fail with 2705.
- `DROP COLUMN` of a missing column, and `ALTER COLUMN` of one, fail with 4924. `DROP COLUMN d, d` fails with 8168. `DROP COLUMN IF EXISTS` of a missing column is skipped. `DROP COLUMN b, COLUMN missing` fails and keeps b.
- `DROP CONSTRAINT` of a missing name, an index, or a view fails with 3728 (with 3727). Of another table's constraint it fails with 3733. A bare name (`DROP ck_b`) drops the constraint. `DROP CONSTRAINT CK_B` drops ck_b (CI). `DROP CONSTRAINT uq_c, CONSTRAINT missing` fails and keeps uq_c.
- A name twice in one `CREATE` or `ALTER` (constraints, indexes, or a constraint and an index) fails with 8168, case variants included (CI).
- Two primary keys in one statement fail with 8110. `ADD` of a primary key where one exists fails with 1779. A constraint named like its own table fails with 2714.
- `DEFAULT ... FOR` a column with a default fails with 1781, and for a missing column with 1752.
- `NOCHECK CONSTRAINT` of a missing one fails with 4917, and `ENABLE TRIGGER` of a missing one with 4920.
- `ADD CONSTRAINT ckx CHECK (x > 0), x int` succeeds: the list's columns come first.
- Generated names: `PK__c1__3213E83F3D93D9A5`, `UQ__c1__...`, `DF__c1__b__...`, `CK__c1__c__...`, `FK__c1__d__...`. For a long table: `PK__averyver__...`, `DF__averyveryvery__b__...`.
- After `DROP TABLE t`, `CREATE TABLE pk_t` and a trigger of t's trigger's name succeed.

**Multi-object `DROP`:**
- `DROP TABLE r, missing` and `DROP TABLE missing, r` fail with 3701, and r is gone.
- `DROP TABLE r, r` drops r and fails on the second; `DROP TABLE t, T` does the same (CI).
- `DROP TABLE r, v` fails with 3705 and keeps r.
- `DROP VIEW v, missing`, `DROP PROCEDURE p, missing`, and `DROP INDEX i ON t, missing ON t` drop the first and fail with 3701. `DROP INDEX i ON t, i ON t` fails on the second.

**`sp_rename`:**
- Missing: 15225 for an object (`'dbo.nosuch'`, `'nosuch.t'`, and an index named as an object); 15248 for a missing column or index, and for COLUMN or INDEX of a table name.
- In use: 15335, for the name of a table, a constraint, a column, or an index (the PK's).
- `'dbo.t2'` makes `dbo.t2`, and `'[t9]'` makes `[t9]`; `@objname` takes brackets.
- A table rename keeps its constraints' and indexes' names. `INDEX` of a PK's index renames the constraint too, and OBJECT of the PK renames its index. `'dbo.t', 'T'` succeeds.
- `sp_rename` and `sys.sp_rename` both run the system procedure.
- A missing object given to `sp_addextendedproperty` fails with 15135.

**Batches and sessions:**
- 111 for `CREATE VIEW`, `CREATE/ALTER PROCEDURE`, `CREATE FUNCTION`, `ALTER FUNCTION`, `CREATE TRIGGER`, `ALTER TRIGGER`, `CREATE SCHEMA`, `ALTER VIEW`, and `CREATE OR ALTER VIEW` after another statement, `SET NOCOUNT ON` included, and the batch's other statements do not run. A leading `;` or comment is fine.
- `CREATE PROCEDURE p1 AS SELECT 1; CREATE TABLE dbo.t (a int)` and `... AS BEGIN SELECT 1 END; CREATE TABLE ...` keep the `CREATE TABLE` in p1's definition and create no table. So does `CREATE TRIGGER`.
- `CREATE FUNCTION`, `CREATE VIEW`, and `CREATE SCHEMA` followed by a statement fail with 156; `IF ... CREATE VIEW` fails with 156.
- `ALTER TABLE t ADD z; UPDATE t SET z = 1` in one batch fails with 207 and adds no z.
- In a batch, `CREATE TABLE x4; INSERT INTO missing` creates x4 and fails; `DROP TABLE missing; CREATE TABLE x3` fails and creates x3 (go-mssqldb returns the error).
- In a transaction, 2714 (`CREATE TABLE`, `CREATE VIEW`), 1913, and 3728 end the transaction (`XACT_STATE()` 0). 3701, 208, 2705, 4902, and 15225 leave it open.
- Latin1_General_100_CS_AS keeps tables t and T, and columns a and A, apart; `DROP TABLE t` keeps T.
- `CREATE DATABASE` and `DROP DATABASE` in a transaction fail with 226 and 574; `USE missing` fails with 911.
- `SET NOEXEC ON` and `SET PARSEONLY ON` make the following batches run nothing.
- `GO 2` after `CREATE TABLE g1` fails with 2714 on the second run.
- Through go-mssqldb (v1.9.8 and Bytebase's fork `bytebase/go-mssqldb@3ff3ca07d898`):
  - Pooled `ExecContext` calls see neither a temporary table nor a `USE` of the call before, and the session starts in the DSN's database.
  - A transaction keeps both across batches.
  - `QUOTED_IDENTIFIER` and `ANSI_NULLS` are on, and `SCHEMA_NAME()` is dbo as sa.

**Temporary tables:**
- `CREATE TABLE #tmp` twice fails with 2714, as do `SELECT ... INTO #tmp` twice and `##g` twice.
- `DROP TABLE #missing` fails with 3701; `IF EXISTS` skips it.
- A temporary table's columns and indexes are checked as a table's (2705, 1913, 3701).
- A foreign key on one is skipped with a warning, and its named constraints go to tempdb.
- `CREATE TABLE dbo.[#t]` is `#t`.
- Names longer than 116 characters fail with 193.
- In a case-sensitive database on the CI server, `#TMP` collides with `#tmp`.

**Names:**
- Identifiers longer than 128 characters fail with 103, `DROP TABLE IF EXISTS` included; 128 work.
- `[u ]` and `u` are one name: trailing blanks do not count, while leading ones do.
- `DROP TABLE dbo.T` drops t (CI).
- An unqualified `sysobjects` names the system view in DML (259) and the user's table in DDL.
- Hungarian_CI_AS refuses `CREATE TABLE szsz` beside ssz, `DROP TABLE SZSZ` drops ssz, and `ADD cscs` beside column ccs fails with 2705. Welsh_100_CS_AS refuses g beside f and accepts F. Finnish_Swedish_CI_AI refuses w beside v and cafe beside café. `USE CCS` reaches ccs.
- **The collation sweep:**
  - All 5,508 collations, narrowed to 602 classes by dropping the BIN collations and the `_KS`, `_WS`, `_SC`, and `_UTF8` variants, which only make equality stricter or concern non-ASCII. The first run covered the variants too, up to the Chinese collations, with the same result.
  - Each class grouped, as `nvarchar` (identifiers are `sysname`), every name of up to 3 characters over `[a-z0-9_]` and of exactly 4 letters (509,035), every pair of printable ASCII characters, and `a`+c+`b` for every code point c from 1 to 127.
  - Every equality is explained by ASCII case, `v`=`w`, `i`=`y`, `f`=`g`, and the nine Hungarian doubled digraphs.
  - The only other characters are `-` and `:`: ignorable, and equal to a blank, in Maori_100 and Mohawk_100.

**Triggers and functions:**
- An `AFTER INSERT` trigger that creates a table creates it on `INSERT`.
- `ON DELETE CASCADE` fires the child table's `AFTER DELETE` trigger, which creates a table.
- An `INSTEAD OF INSERT` trigger on a view fires on `INSERT` to the view.
- A function with `CREATE TABLE` fails to compile (443).

**Codex round 1 on #514**, verified on the same image:
- In a Latin1_General_100_CS_AI database, tables `cafe` and `[CAFÉ]` coexist; `ALTER TABLE dbo.CAFE ADD a int` adds `a` to `CAFÉ`; `DROP TABLE dbo.cafe` and then `DROP TABLE dbo.CAFE` drop both. A class with one object stands for another spelling only while its space holds no marked name, and a class the change emptied leaves a spelling missing only then.
- A CHECK, DEFAULT, or FOREIGN KEY constraint and an inline `INDEX` of one name in `CREATE TABLE` fail with 8168, as do a UNIQUE constraint and an index.
- `ALTER TABLE dbo.t DROP CONSTRAINT IF EXISTS nokey`, of a table's name, and `... IF EXISTS pk_r`, of another table's key, succeed and keep the object.
- `CREATE UNIQUE CLUSTERED INDEX pk_t ON dbo.t (id) WITH (DROP_EXISTING = ON)` keeps `pk_t` a primary key; `sp_rename` of its index then renames the constraint, and `CREATE TABLE dbo.pk_t` succeeds.

**Sync** (`backend/plugin/db/mssql/sync.go`), probed on its own queries:
- `INFORMATION_SCHEMA.SCHEMATA` lists every schema as sa.
- Clustered and nonclustered columnstore, XML, and hypothetical indexes have `sys.index_columns` rows.
- An edge table's constraint is type EC, and its columns include `$edge_id_...` and hidden graph columns. A system-versioned table's history table is `MSSQL_TemporalHistoryFor_<id>` in the table's schema.
- Legacy defaults (type D, no parent) and rules (R), and queues (SQ, with an internal table), are in `sys.objects`.
- Disabled triggers hold their names.
