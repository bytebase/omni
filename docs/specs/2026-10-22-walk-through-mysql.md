# WALK_THROUGH for MySQL: v1 scope

Status: decided 2026-10-23; §9 records the decisions. Step 2 of [the engines design](2026-10-22-walk-through-engines-design.md) §7, which its §5.1 sketches. Every fact here was verified on MySQL 5.7.44, 8.0.46 (with `lower_case_table_names` 0 and 1), and 8.4.13 (Appendix). Where the versions differ, the walk claims only what holds on all of them.

## 1. Environment

- A Bytebase database is one MySQL database. Its metadata holds one schema named `""`. The connection selects the database, and the change runs on one connection, so `USE` and temporary tables carry over between its statements.
- **Current database:** the target at the start. After `USE x`, unqualified names resolve in `x`, which the walk does not know unless the change made it, so they are Unknown until a `USE` returns to the target.
- **Other databases:** only their existence can be known, and only from the change itself. `CREATE DATABASE q` makes q Present; nothing makes another database Absent. The target is Present. After `DROP DATABASE` of the current database, an unqualified creation fails with 1046 "No database selected".
- **Temporary tables:**
  - The session starts with none.
  - A temporary table shadows a base table of its name for every statement but `CREATE TABLE`.
  - `DROP TEMPORARY TABLE` never drops a base table.
  - `DROP TABLE` drops the temporary table first.
- **Server settings the walk reads:** `Target.LowerCaseTableNames`, and `Target.ForeignKeyChecks` (§9 M1). It reads no other.

## 2. Spaces

| Space | Parent | Names compare | Claims |
| --- | --- | --- | --- |
| database | root | as tables do | exists (1007), missing (1049, 1008) |
| relation: tables and views, one space | database | `lower_case_table_names`: 0 exact; 1 and 2 ASCII-lowercase | exists (1050), missing (1146, 1051), wrong kind (1347) |
| temporary | database | as tables | shadows relation for resolution |
| column | table | ASCII case-insensitive | exists (1060), missing (1054, 1091) |
| index | table | ASCII case-insensitive; `PRIMARY` reserved (1280) | exists (1061), missing (1091) |
| foreign key | database | ASCII case-insensitive | exists (1826 on 8.0+, 1022 on 5.7), missing (1091) |
| trigger | database | case-sensitive on 5.7, insensitive on 8.0 | exists only on the exact spelling (1359), missing only when no ASCII case variant exists (1360) |
| procedure, function: two spaces | database | ASCII case-insensitive | exists (1304), missing (1305) |
| event | database | ASCII case-insensitive | exists (1537), missing (1539) |

- **Non-ASCII names:** they compare equal only when spelled identically. `é` and `É` are two names under every setting.
- **Length:** an identifier longer than 64 characters fails with 1059 before anything is looked up. The walk claims this for ASCII names, from the text alone.
- **Not modeled:**
  - CHECK constraint names: 5.7 parses CHECK and ignores it;
  - the primary key's absence (§3);
  - partitions.

## 3. Authority

Present is trusted from the sync. Absent is trusted for:
- **relations** of the target database: the sync lists every `BASE TABLE` and view;
- **columns and index names** of synced tables, apart from `PRIMARY`;
- **foreign key names,** triggers, routines, and events of the target database.

Never trusted Absent:
- **the primary key:** a server with `sql_generate_invisible_primary_key=ON` gives a table created without one a hidden key. With `show_gipk_in_create_table_and_information_schema=OFF`, the sync does not see it (verified: the hidden key then fails `ADD PRIMARY KEY` with 1068 and `DROP PRIMARY KEY` with 1235);
- **another database's objects;**
- **CHECK constraints.**

**Accepted assumption (engines design D6):** the admin user sees every object of its database, since `information_schema` lists only what it may see.

**OceanBase** (`ReviewOceanBase`) will use the same lowering, with less authority: its sync lists no CHECK constraints or partitions. Its oracle is DSN-gated (D4). Until it lands (§9 M3), WALK_THROUGH reports nothing for OceanBase.

## 4. Statement semantics the lowering follows

**`ALTER TABLE` runs drops first.**
- DROP, RENAME, CHANGE, and MODIFY of a column or index are judged against the table as it was. `ADD z, DROP z` fails with 1091, and `ADD a, DROP a` succeeds where a exists and leaves a new a.
- ADD COLUMN sees the table after them. Index columns resolve against the final columns, so `ADD INDEX k (q), ADD q` succeeds.
- This is `walk.Passes`, with the drops and renames lowered first.
- A second DROP of one column or index fails with 1091, and a second ADD with 1060 or 1061.

**`DROP TABLE a, b` looks every name up first** (`walk.Lookups`).
- 8.0 and 8.4 are atomic; 5.7 drops `a` before failing on a missing `b`. The claim that the statement fails holds on both.
- `DROP TABLE` on a view fails with 1051, and `DROP VIEW` on a table with 1347.

**`RENAME TABLE a TO b, c TO d` runs its pairs in order,** so a swap through a temporary name succeeds.
- Each pair is a part, its effects applied in order (`walk.Passes`). A pair's check on a name another pair changes is left out.
- A missing source fails, and the rename is undone on every version.

**`CREATE TABLE IF NOT EXISTS t`** is a notice when t exists, as a table or a view. Its `LIKE` source is checked first, so a missing source fails even then (1146).

**Kinds:**
- `ALTER TABLE`, `CREATE INDEX`, and `CREATE OR REPLACE VIEW` on the wrong kind fail with 1347.
- `TRUNCATE` of a view fails with 1146.

**Foreign keys:**
- A foreign key to a missing table or column fails only while `foreign_key_checks` is on (1824, 3734; 1215 on 5.7).
- The walk tracks `SET foreign_key_checks` in the change; see decision M1 for the session's start.

**Generated names:**
- an unnamed index takes its first column's name, then `_2`, `_3`, ...;
- an unnamed foreign key is `<table>_ibfk_<n>`, and it adds an index named after itself or its first column when none serves it;
- an unnamed CHECK is `<table>_chk_<n>`.

These are May-effects on the table's index space and on the database's foreign key prefix.

**Dependency failures happen, but the walk makes no claim about them:**
- `DROP TABLE` of a table a foreign key references (3730);
- `DROP INDEX` of an index a foreign key needs (1553).

Views are not dependencies: dropping a table a view reads succeeds.

## 5. Statements in v1

| Statement | Checks | Effects |
| --- | --- | --- |
| `CREATE DATABASE [IF NOT EXISTS] d` | d Absent (1007) | d Present, everything under it Absent |
| `DROP DATABASE [IF EXISTS] d` | d Present (1008) | d Absent, everything under it Absent; dropping the current database leaves none selected |
| `USE d` | d Present (1049) | current database |
| `CREATE [TEMPORARY] TABLE [IF NOT EXISTS] t (...)` | database Present (1049) or selected (1046); t Absent (1050), skipped for a temporary table; within the statement: duplicate column (1060), duplicate index name (1061), two primary keys (1068), index named `PRIMARY` (1280), identifier too long (1059); a foreign key's target Present while `foreign_key_checks` is on (1824) | t Present (table) with exact columns, its named indexes and foreign keys Present; generated names as May-effects |
| `CREATE TABLE t LIKE s`, `CREATE TABLE t AS SELECT` | s Present (1146), even with IF NOT EXISTS; t Absent (1050) | t Present with unknown columns and indexes |
| `ALTER TABLE t ...` | t Present (1146) and a table (1347); per subcommand: | per subcommand, drops first: |
| — `ADD COLUMN c [FIRST, AFTER x]` | c Absent (1060); x Present (1054) | c Present; inline keys as in CREATE TABLE |
| — `DROP COLUMN c` | c Present (1091) | c Absent; indexes on it May-removed |
| — `CHANGE a b`, `MODIFY a`, `RENAME COLUMN a TO b` | a Present (1054); b Absent unless b is a (1060) | a moves to b |
| — `ADD INDEX/KEY/UNIQUE/FULLTEXT/SPATIAL [i]`, `ADD PRIMARY KEY`, `ADD FOREIGN KEY [f]` | i Absent (1061); a primary key the table certainly has (1068); f Absent in the database (1826); the target Present while `foreign_key_checks` is on | i, f Present; generated names May-added |
| — `DROP INDEX i`, `DROP KEY i`, `DROP PRIMARY KEY`, `DROP FOREIGN KEY f` | i Present (1091); f Present (1091) | i, f Absent |
| — `RENAME INDEX a TO b` | a Present (1091); b Absent (1061) | a moves to b |
| — `RENAME [TO] t2` | t2 Absent (1050); its database Present (1049; 1025 on 5.7) | t moves to t2 |
| — table options, `ALGORITHM`, `LOCK`, `ENGINE`, `COMMENT`, `ALTER COLUMN SET/DROP DEFAULT`, `ALTER INDEX VISIBLE`, `CONVERT TO`, `FORCE` | none beyond t | none |
| — partitions, `DISCARD/IMPORT TABLESPACE`, `ADD/DROP/ALTER CHECK`, `ORDER BY`, `ENABLE/DISABLE KEYS` | none beyond t | none: they change no name the walk keeps (§7) |
| — `DROP CONSTRAINT c` | none | c's foreign key and index Unknown: c may name either, or a CHECK |
| — anything else | | ends the walk |
| `RENAME TABLE a TO b, ...` | per pair: a Present (1146), b Absent (1050), b's database Present (1049) | each pair moves in order |
| `DROP [TEMPORARY] TABLE [IF EXISTS] a, ...` | each a Present unless IF EXISTS, and a table (1051) | a and what it holds Absent |
| `TRUNCATE [TABLE] t` | t Present and a table (1146) | none |
| `CREATE [UNIQUE] INDEX i ON t`, `DROP INDEX i ON t` | t Present (1146) and a table (1347); i Absent (1061) / Present (1091) | i Present / Absent |
| `CREATE [OR REPLACE] VIEW v`, `ALTER VIEW v` | v Absent (1050), or a view for OR REPLACE (1347); ALTER: v Present (1146) | v Present (view) |
| `DROP VIEW [IF EXISTS] v, ...` | each v Present (1051) and a view (1347) | v Absent |
| `CREATE/DROP TRIGGER`, `PROCEDURE`, `FUNCTION`, `EVENT [IF [NOT] EXISTS]` | name Absent / Present per §2 (1359/1360, 1304/1305, 1537/1539); a trigger's table Present (1146) | name Present / Absent; a trigger marks its table (§6) |
| `INSERT`, `REPLACE`, `UPDATE`, `DELETE`, `LOAD DATA` | each target relation Present (1146) | none; see §6 |
| `SELECT`, `DO`, `VALUES`, `TABLE`, `EXPLAIN`, `SET` (other than `foreign_key_checks`), `SHOW`, `ANALYZE/OPTIMIZE/CHECK/CHECKSUM/REPAIR TABLE`, `FLUSH`, `GRANT/REVOKE` and the user and role statements, `BEGIN/COMMIT/SAVEPOINT/RELEASE`, `UNLOCK TABLES` | none | none, unless it may call a stored function (§6) |
| `SET foreign_key_checks = 0/1` | none | session state |

## 6. Where the walk stops

- `CALL`, `PREPARE`, `EXECUTE`, `DEALLOCATE`, `HANDLER`, `XA`, `ROLLBACK [TO SAVEPOINT]`, `LOCK TABLES`, and anything not in §5.
- **Temporary tables are forgotten, and the walk goes on,** after a statement that may run code which creates or drops one:
  - DML, after its checks, on a table that has a trigger (synced or created by the change), on a view or a table the walk does not know, or on any table while triggers are unknown (`OtherSpecs`). Cascaded foreign key actions fire no triggers (verified), so only the statement's own target tables count.
  - A statement that may call a stored function: a qualified call, a call to a function of the target database, or any call while functions are unknown (`OtherSpecs`).
  - Neither a trigger nor a stored function may run a statement that commits: `CREATE FUNCTION` with `CREATE TABLE` in its body fails with 1422, while one that creates or drops a temporary table runs (verified on 8.0). So nothing else is forgotten.
- A routine or event body is not walked; its `CREATE` is lowered for the name alone.
- **Accepted assumption (D6):** events do not change the schema while the change runs.

## 7. Not covered (declined by scope)

Review asks for these are declined with this list:

- CHECK constraints: names, existence, and violations;
- dependency failures: dropping a referenced table (3730), an index a foreign key needs (1553), a column a foreign key uses;
- columns named in index or foreign key definitions (1072, 3734), and in views, defaults, generated columns, and queries;
- relation references inside queries and view bodies;
- data-dependent failures: duplicate keys, NULL into NOT NULL, conversions under the SQL mode;
- privileges, `DEFINER`, and `sql_mode`;
- partitions and tablespaces;
- MySQL 5.7's partial effects of a failed multi-table DROP;
- non-ASCII identifiers longer than 64 characters.

## 8. Budget and tests

- **Code:** about 1.2k lines of lowering, 400 of names and snapshot, and 250 of walker and plumbing (`targets.go` gains `walks` and `OtherSpecs` in its key). As built: 1330 lines of lowering (`walk_lower*.go`), 225 of snapshot, 280 of walker.
- **Oracle:** `walktest.Judge` against `mysql:8.0` with `lower_case_table_names` 0 and 1, and `mysql:8.4`. MySQL commits DDL as it runs, so the engine implements `walktest.Digester`.
- **Curated cases:** one per row of §5 and per fact in the Appendix, the 5.7/8.0 differences included as text-only claims.
- **Sweep:** templates over a colliding name pool with mixed case, temporary tables, `USE`, and `SET foreign_key_checks`. CI runs it at `WALK_SWEEP_N` per `lower_case_table_names` setting, the two at once. At 500 each (seed 20261023, 8.0) it found 5247 of 7787 rejections with lctn 0 and 5119 of 7747 with lctn 1, with no false finding. Most misses are by design: statements on a database other than the target, and statements after a walk ends at a `RENAME` through a temporary or Unknown table.
- **Coverage:** `walktest.Coverage` over every row of §5.

## 9. Decisions

Decided 2026-10-23, each as recommended:

- **M1:** Bytebase's MySQL sync reads `foreign_key_checks` over the admin connection and passes it as `Target.ForeignKeyChecks` (`review.Switch`: unknown, on, off). Until Bytebase fills it, the field is unknown, and the walk checks a foreign key's target only after the change itself sets `foreign_key_checks = 1`.
- **M2:** the CI oracle runs on `mysql:8.0` with `lower_case_table_names` 0 and 1, and on `mysql:8.4`. 5.7 stays in the Appendix.
- **M3:** OceanBase comes in a follow-up PR.

The options as they were weighed:


- **M1. `foreign_key_checks` at the session's start.** What Bytebase does (origin/main b555065095):
  - **Execution** sets no session variable itself, so the session starts with the server's global default.
  - **The plain MySQL connection** (not RDS IAM or Cloud SQL) appends a data source's extra connection parameters to the DSN (`plugin/db/mysql/mysql.go`). go-sql-driver/mysql sends a parameter it does not know as `SET <name>=<value>` on connect. A data source configured with `foreign_key_checks=0` therefore runs every migration with the checks off, and so does a global default a DBA turned off.
  - **The schema dump** (`GetDatabaseDefinition`) wraps its DDL in `SET FOREIGN_KEY_CHECKS=0` and restores it at the end. That is in the change's text, which the walk follows.
  - **v1's walk-through** runs its catalog with `foreign_key_checks = 0` (`plugin/schema/mysql/walk_through.go`), so v1 never reports a missing foreign key target.
  - **Options:**
    - (a) Assume on: this would report a false failure on a data source configured with the checks off.
    - (b) Check foreign key targets only after the change sets `foreign_key_checks = 1`.
    - (c) Bytebase's MySQL sync reads `foreign_key_checks` over the admin connection, as it reads `lower_case_table_names` (`getServerVariable`, `SHOW VARIABLES`, which gives the session value, the extra parameters included). It is stored beside `MysqlLowerCaseTableNames` and passed as a new `Target.ForeignKeyChecks`.

  Recommended: (c), with (b) until the field lands. This is the pattern of D2: a session setting the target cannot know is read where Bytebase connects, not assumed.
- **M2. MySQL 5.7 in the CI oracle.**
  - (a) No: 5.7 facts stay in the Appendix as probes, with the claims that differ kept text-only.
  - (b) Yes: a third container, without `Digester`, because 5.7 is not atomic.

  Recommended: (a).
- **M3. OceanBase in this PR or the next.**
  - (a) This PR, with its authority deltas, and the DSN-gated oracle skipped until a DSN exists.
  - (b) A follow-up PR.

  Recommended: (a) if an OceanBase instance is available, else (b).

## Appendix: verified facts

MySQL 5.7.44, 8.0.46 (`lower_case_table_names` 0 and 1), and 8.4.13. "All" means every one of them.

- **`ALTER TABLE`:**
  - `ADD a, DROP a` and `DROP a, ADD a` succeed where a exists; `ADD z, DROP z` and `DROP z, ADD z` fail with 1091 where it does not (all). The same holds for indexes.
  - `RENAME COLUMN a TO c, ADD a` and `CHANGE b d, ADD b` succeed. `ADD x, RENAME COLUMN x TO y` fails with 1054. RENAME COLUMN is 8.0+.
  - `ADD INDEX k (q), ADD q` succeeds; `DROP a, ADD INDEX (a)` fails with 1072 (all).
  - Dropping the only column of an index drops the index (all).
- **Kinds and namespaces:**
  - Tables and views share a namespace (1050); the wrong kind gives 1051, 1347, and 1146 (TRUNCATE) (all).
  - Foreign key names are unique per database: 1826 on 8.0/8.4, 1022 on 5.7.
  - CHECK names are unique per database on 8.0/8.4 (3822); 5.7 accepts duplicates.
  - An index named like a foreign key's own index fails (1061). A CHECK may share an index's name.
  - Trigger names are unique per database (1359). `TR` beside `tr` is accepted on 5.7 and refused on 8.0/8.4.
  - A procedure and a function may share a name. Routine and event names ignore case (1304, 1537).
  - `PRIMARY` as an index name gives 1280. Two primary keys give 1068. `DROP PRIMARY KEY` without one gives 1091 (all).
- **Case:**
  - With `lower_case_table_names=0`, `T` and `t` are two tables and two databases. With 1, `CREATE TABLE T` fails with 1050 when t exists, `DROP TABLE T` drops t, and the same holds for views and databases.
  - Column, index, and foreign key names ignore case everywhere.
  - `é` and `É` are distinct column and table names under every setting.
- **Names:**
  - Identifiers over 64 characters give 1059, even for `DROP TABLE` of a missing one. A table name with a trailing space gives 1103.
  - Two indexes of one name in `CREATE TABLE` give 1061, and an index named `PRIMARY` gives 1280.
  - `ADD FOREIGN KEY` with a name another table's foreign key has fails (1826; 1022 on 5.7).
- **Atomicity:**
  - `DROP TABLE t1, missing` fails with 1051. 8.0/8.4 keep t1; 5.7 drops it.
  - `RENAME TABLE a TO b, missing TO c` fails (1146; 1017 on 5.7) and keeps a, everywhere.
  - A swap through a temporary name succeeds.
- **`CREATE TABLE`:**
  - A failed `CREATE TABLE ... AS SELECT` leaves no table.
  - `CREATE TABLE IF NOT EXISTS v` on a view is a notice.
  - `CREATE TABLE IF NOT EXISTS v LIKE missing` fails with 1146.
- **Foreign keys:** to a missing table they fail with `foreign_key_checks=1` (1824; 1215 on 5.7) and are accepted with 0. A missing referenced column fails (3734; 1215).
- **Temporary tables and `USE`:**
  - A temporary table shadows a base table. `DROP TEMPORARY TABLE` of a name with only a base table fails with 1051.
  - After `USE q`, unqualified names resolve in q.
- **Databases:**
  - `CREATE TABLE`/`CREATE VIEW` in a missing database fail with 1049; `DROP TABLE` there fails with 1051, and IF EXISTS makes it a note.
  - After `DROP DATABASE` of the current database, `CREATE TABLE t` fails with 1046.
  - `RENAME TABLE`, and `ALTER TABLE ... RENAME TO`, into a missing database fail (1049; 1025 on 5.7).
  - `USE` of a missing database fails with 1049; `DROP VIEW` of a missing view with 1051.
- **Generated names:** `INDEX (a)` twice gives `a` and `a_2`. An unnamed foreign key is `t_ibfk_1` and adds an index named `a` when none serves it. An unnamed CHECK is `t_chk_1`.
- **Missing objects:**
  - `DROP TRIGGER`/`PROCEDURE`/`FUNCTION`/`EVENT` of a missing one fail with 1360, 1305, 1305, and 1539.
  - `CREATE TRIGGER` on a missing table fails with 1146.
  - `CREATE TRIGGER/PROCEDURE IF NOT EXISTS` are 8.0+.
- **Invisible primary key:** with it on, a table created without a key gets one, and `show_gipk_in_create_table_and_information_schema` (on by default) decides whether `information_schema` lists it. On the hidden key, `DROP PRIMARY KEY` fails with 1235, `ADD PRIMARY KEY` with 1068, and `ADD COLUMN my_row_id` with 1060.
- **Triggers:** a cascaded `ON DELETE CASCADE` delete does not fire the child table's trigger.
- **No failure:**
  - `ANALYZE`, `OPTIMIZE`, and `CHECK TABLE` of a missing table return a row rather than fail.
  - After `LOCK TABLES t WRITE`, DML on another table fails with 1100.
