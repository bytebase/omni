# WALK_THROUGH for OceanBase (MySQL mode): v1 scope

Status: draft 2026-10-27; O1 decided and O3 resolved 2026-10-28, the rest of §8 pending. This is the OceanBase part of step 2 of [the engines design](2026-10-22-walk-through-engines-design.md) §7, which the MySQL PR deferred (M3). It is written as deltas from [the MySQL v1 scope](2026-10-22-walk-through-mysql.md): what that spec says holds here unless this one says otherwise. Every fact here was verified on OceanBase CE 4.3.5.6 (`oceanbase/oceanbase-ce:4.3.5-lts`, `MODE=mini`, the MySQL-mode tenant `test` at its defaults), with `mysql:8.0` as the baseline (Appendix).

## 1. Shape

- **Same package, same lowering.** OceanBase is served by `mysql/review` (`ReviewOceanBase`). The walker gains an `oceanBase` flag from the dialect, and every delta below is gated on it, so MySQL's behavior and oracle are unchanged. Each delta names its verified fact in a comment.
- **Bytebase is not changed by this work.** Its sync gained OceanBase's external tables and materialized views in bytebase/bytebase#21685 and its events in bytebase/bytebase#21687 (both 2026-10-28), which §2 and §3 rely on. Its rule review already maps `OCEANBASE` to `ReviewOceanBase` and asks every engine for WALK_THROUGH, so the pin bump turns the walk on. It sets `OtherSpecs` (bytebase/bytebase#21693) and leaves `ForeignKeyChecks` unset, and the walk handles both as it does on MySQL.
- **The metadata is current** (decided 2026-10-28): every synced schema is taken as the output of Bytebase's current sync code, run against the schema as it is now. The walk does not account for an older sync or a failed one.
- **The version.** Bytebase names no supported OceanBase version. It detects OceanBase by the instance's engine (`OCEANBASE`, registered on the MySQL driver) and reads `SELECT VERSION()`, which 4.3.5 answers `5.7.25-OceanBase_CE-v4.3.5.6`. The oracle pins the 4.3.5 LTS line; `OB_IMAGE` in `scripts/oceanbase-oracle.sh` runs another.
- **Claims that rest on a restriction are left out.** OceanBase refuses many combinations as "not supported" (1235, 4179): several foreign key changes in one `ALTER TABLE`, a dropped and an added index in one, two primary key changes in one. Later versions lift such restrictions, and one image cannot show which, so the walk claims none of them (§8 O5).

## 2. What Bytebase syncs from OceanBase

The OceanBase driver is the MySQL one (`backend/plugin/db/mysql/mysql.go` registers `Engine_OCEANBASE` on `newDriver`), so `SyncDBSchema` is MySQL's (`sync.go`, Bytebase origin/main 43fbe0afdd), with these differences, all from the version string and the engine:

- **CHECK constraints are not read:** the query runs only from 8.0.16, and OceanBase reports 5.7.25. The walk never trusts them anyway (MySQL spec §3).
- **Index expressions are not read:** the `EXPRESSION` column is selected only from 8.0.13. The walk reads index names, not their parts.
- **Partitions are not read:** `listPartitionTables` runs only for `Engine_MYSQL`. The walk does not model partitions.
- **Tables, views, and external tables** come from `information_schema.TABLES`: `TABLE_TYPE` `BASE TABLE` as tables, `VIEW` as views, and `EXTERNAL TABLE` as external tables with their columns (bytebase/bytebase#21685). They come from one query, so the sync lists all of them or fails.
- **Materialized views** are in neither `TABLES` nor `VIEWS`. Since #21685 the sync reads them from `oceanbase.DBA_MVIEWS` into `materialized_views` (`skip_dump`, with the comment and indexes of the view's storage table). The storage table `__mv_container_<id>` and a materialized view log `mlog$_<table>` stay listed as tables. Reading `DBA_MVIEWS` takes `SELECT` on the `oceanbase` database; without it the sync succeeds with no materialized views. Every materialized view has exactly one storage table, which `TABLES` lists even then, so a storage table that no listed materialized view accounts for shows the list is incomplete (verified in #21685 on 4.3.5.6).
- **`lower_case_table_names`** is read with `SHOW VARIABLES` at instance sync and reaches `Target.LowerCaseTableNames`. A MySQL-mode tenant defaults to 1, fixed when the tenant is created.
- **Triggers, routines, and foreign keys** come from the same `information_schema` queries, which OceanBase answers. Every table reports `ENGINE` `InnoDB`.
- **Events** are listed by name only since bytebase/bytebase#21687: `SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS`, with no definition, because 4.3.5 rejects `SHOW CREATE EVENT` (1064), which made the sync of a database with an event fail before. `DROP DATABASE` leaves the dropped database's events in `information_schema.EVENTS` under its name, so a database created again under that name lists them as its own, and `CREATE EVENT` of such a name succeeds there (verified in #21687 on 4.3.5.6).

## 3. Authority

What differs from the MySQL spec's §3:

- **Relations of two more kinds.** An external table and a materialized view share the table namespace: `CREATE TABLE` and `CREATE VIEW` of their name fail with 1050. `RENAME TABLE` and `DROP TABLE` of an external table succeed; `RENAME TABLE`, `CREATE INDEX`, and `DROP INDEX` on a materialized view succeed; most other DDL on either fails with 1235 (Appendix). The snapshot lists each synced one as Present of unknown kind, which satisfies every kind a check asks for, so the walk claims no wrong kind of them. Their columns and indexes are Unknown, and an `ALTER TABLE` of one checks that it exists and ends the walk, as on a view (§4.1): what succeeds on them was probed only as far as the Appendix.
- **Relation absence** is trusted, as on MySQL, in a database where the sync lists no more `__mv_container_` tables than materialized views (§8 O1). Where it lists more, the `DBA_MVIEWS` read failed (§2), and relation absence is untrusted in that database, as on MariaDB. A table a user names `__mv_container_x` only makes that check fail, which costs recall.
- **The primary key.**
  - OceanBase has no invisible primary key: `sql_generate_invisible_primary_key` does not exist, and `DROP PRIMARY KEY` on a table without one fails (1091). So a table the sync lists without an index named `PRIMARY` has no primary key, which MySQL does not trust.
  - A secondary index may be named `PRIMARY` (`CREATE INDEX`, `ALTER TABLE ... ADD INDEX`; only `CREATE TABLE` refuses it, 1280), on a table with or without a primary key, and Bytebase's sync marks any index named `PRIMARY` as the primary key. So an index the sync lists as `PRIMARY` leaves the primary key Unknown (§8 O2).
  - A table without a primary key has the hidden column `__pk_increment`, which `information_schema` does not list. `ADD COLUMN x AFTER __pk_increment` and `ALTER COLUMN __pk_increment SET DEFAULT` succeed, and `ADD COLUMN __pk_increment` fails (1060). The column is Unknown on every table, in place of MySQL's `my_row_id`.
- **Foreign keys.** OceanBase ignores `ENGINE` (`ENGINE=MyISAM` is a warning, 1286), so every table keeps its foreign keys, synced or created; MySQL's engine rule does not apply.
- **Names outside ASCII.** OceanBase compares every name with `utf8mb4_general_ci`, accents included: `é` is `e` and `É` for tables, columns, indexes, databases, routines, and triggers, and `DROP TABLE É` drops `e`. MySQL's walk already follows no such name (bytebase/omni#515): a statement with a name outside ASCII ends the walk, and a synced space that lists one has no name known to be absent. That covers OceanBase's folding, so no delta is needed.
- **Events.** Absence is trusted, as on MySQL: the sync lists every event of the database. A listed event is Unknown: it may be a dropped database's, whose name is free (§2), so a `CREATE EVENT` of its name is not claimed (§8 O8).

## 4. Semantics that differ from MySQL

### 4.1 ALTER TABLE

OceanBase judges every subcommand against the table as it was, with exceptions, and refuses most pairs of subcommands on one name:

- `ADD a, DROP a` and `DROP a, ADD a` fail with 1060 where a exists; `RENAME a TO b, DROP b` and `DROP b, RENAME a TO b` fail with 1060; a swap of two indexes fails (1061).
- But a rename frees its old name (`RENAME a TO c, ADD a` succeeds), a swap of two columns succeeds, an `AFTER` sees a column the statement adds or renames to, and `MODIFY a`, `ALTER COLUMN a SET DEFAULT`, or `AFTER a` beside `DROP a` succeed, and a is dropped.

MySQL's lowering (drops and renames first, a check on a key another part changes left out) still holds: each check it keeps is on a key no other part touches, which OceanBase also checks against the table as it was, and for every combination OceanBase accepts, its effects are the walk's (Appendix). Two of MySQL's text claims are false here and are left out:

- `DROP FOREIGN KEY f, DROP FOREIGN KEY f` succeeds and drops f;
- `` ADD INDEX `PRIMARY` `` succeeds (§4.3).

Repeated forms still fail: a column added twice (1060), dropped twice (1091), an index added twice (1061), dropped twice (1072), `DROP PRIMARY KEY` twice (1235), two primary keys (1068, or 1235 with `CHANGE ... PRIMARY KEY`). An `ALTER TABLE` of a view fails (1347) for every subcommand probed, `ENGINE=` among them, but `COMMENT`: `ALTER TABLE v COMMENT 'x'` succeeds. So on a view the walk checks only that it exists and then ends (§8 O7).

### 4.2 Foreign keys

- **Names** are unique per database and compare without case, as on MySQL, but a taken name fails with 1022, and two of one name in a statement with 1215.
- **Generated names** are `<table>_OBFK_<timestamp>`, a May-effect on the prefix `<table>_obfk_`. `RENAME TABLE` and `ALTER TABLE ... RENAME` keep every foreign key's name, generated or not; MySQL's `_ibfk_` renaming does not apply.
- **An unqualified parent is in the current database,** not the child's: `CREATE TABLE q.c (... REFERENCES r)` with shop current references shop.r. With no database selected, the walk does not check it.
- **Targets** (while `foreign_key_checks` is on, as on MySQL): a missing parent fails with 1146, a view or a temporary table with 1215. A view fails even with the checks off, and a temporary table that shadows a base table fails too (1215); both only cost recall.
- `CREATE TABLE ... LIKE` copies no foreign key, and a foreign key adds no index, as the walk already assumes or over-approximates.
- `DROP CONSTRAINT c` drops a foreign key or a CHECK, never an index (3940). MySQL's lowering, which leaves the index Unknown too, is weaker and stays.

### 4.3 Indexes named PRIMARY

`` CREATE INDEX `PRIMARY` `` and `` ALTER TABLE ... ADD INDEX `PRIMARY` `` make a secondary index of that name, `` RENAME INDEX `PRIMARY` TO k `` renames it, and `` DROP INDEX `PRIMARY` `` drops only it: with only a primary key, it fails (1091). So `PRIMARY` as an index name in `CREATE`, `DROP`, `RENAME`, or `ADD INDEX` ends the walk. `ADD PRIMARY KEY` and `DROP PRIMARY KEY` name the primary key, as on MySQL, and `CREATE TABLE` refuses an index named `PRIMARY` (1280), which the walk still claims.

### 4.4 Kinds and views

- `RENAME TABLE v TO q.v` moves a view to another database; MySQL's 1450 claim is left out.
- `DROP VIEW IF EXISTS t` of a table fails (1347), as on TiDB. The walk claims it.
- `DROP TABLE` of a view (1051), `DROP VIEW` of a table (1347), `TRUNCATE` of a view (1146), `CREATE INDEX` and `CREATE OR REPLACE VIEW` on the wrong kind (1347), and `LIKE` a view (1347) are as on MySQL.
- Sequences are a namespace of their own: a table and a sequence may share a name, and `DROP TABLE` does not reach a sequence. The walk does not model them.

### 4.5 Temporary tables

As on MySQL, but **`DROP DATABASE` drops the session's temporary tables in it**: they do not come back with the database (1146, 1051). `CREATE TEMPORARY TABLE IF NOT EXISTS` of an existing one checks its text (1060), which costs recall only.

Neither a trigger nor a stored function can create or drop a temporary table, even through `CALL` (1422). Cascaded foreign key actions do fire the child's triggers on OceanBase, unlike MySQL, but since no trigger can reach a temporary table, MySQL's rule for forgetting them stays, as a conservative one.

### 4.6 Routines, triggers, and events

- **Loadable functions cannot be created** (1235), and `mysql.func` is empty, so an unqualified `DROP FUNCTION f` names the stored function: the walk claims 1305 for it. With no database selected, `DROP FUNCTION IF EXISTS f` succeeds (a note about a loadable function), so then the walk claims nothing.
- **Trigger names** are per database and compare without case: one of the name on the same table fails with 5734, on another table with 5735, even with `IF NOT EXISTS`. MySQL's spelling rule stays, as a conservative one.
- **Events:** a quoted event name with a period is split: `` CREATE EVENT `shop.e3` `` makes e3 in shop. The walk ends at an event name with a period. `ALTER EVENT ... RENAME` of a running event fails (1235); that costs recall only.

### 4.7 Names and lengths

| Name | Longest |
| --- | --- |
| database | 128 |
| column | 128 |
| foreign key | 128 (1406 beyond) |
| table, view, index, CHECK, trigger, routine, event | 64 |

Longer names fail with 1059, before anything is looked up, as on MySQL. `lower_case_table_names` behaves as MySQL's 1: names are stored and compared in lower case.

### 4.8 Error codes

Only the oracle's code map changes:

| Case | OceanBase code |
| --- | --- |
| foreign key name taken | 1022 |
| two foreign keys of one name in a statement | 1215 |
| foreign key to a missing table / a view or temporary table | 1146 / 1215 |
| two primary keys in `CREATE TABLE` | 1062 or 1068 |
| an index dropped twice in one `ALTER TABLE` | 1072 |
| `DROP PRIMARY KEY` twice | 1235 |
| `RENAME TABLE` of a missing table | 1017 |
| trigger name taken, on the same / another table | 5734 / 5735 |
| `ALTER TABLE ... RENAME TO` a missing database, `ALTER DATABASE` with none selected | 1210 |
| a trigger or `ALTER TABLE` on a view | 1235 / 1347 |

## 5. Statements in v1

The MySQL table (MySQL spec §5) with these changes, all gated on OceanBase:

| Statement | Change |
| --- | --- |
| a synced external table or materialized view | Present of unknown kind, columns and indexes Unknown (§3) |
| a synced database with more `__mv_container_` tables than materialized views | relation absence Unknown in it (§3) |
| `ALTER TABLE t ...` on a relation of unknown kind | t Present (1146), then the walk ends, as on a view |
| a synced event | Unknown (§3) |
| `CREATE TABLE` | no engine rule (§3); generated foreign keys `<t>_obfk_`; an unqualified parent in the current database; name limits per §4.7; no primary key means none, with `__pk_increment` Unknown |
| `ALTER TABLE t ...` on a view | t Present (1146), then the walk ends |
| — `DROP FOREIGN KEY f` twice | no claim |
| — `` ADD INDEX `PRIMARY` ``, `` DROP INDEX `PRIMARY` ``, `RENAME INDEX` from or to `PRIMARY` | ends the walk |
| — `DROP PRIMARY KEY` | `__pk_increment` Unknown |
| — `ENGINE=` | foreign keys still kept |
| `RENAME TABLE`, `ALTER TABLE ... RENAME` | foreign key names stay; a view may change database |
| `` CREATE INDEX `PRIMARY` ``, `` DROP INDEX `PRIMARY` `` | end the walk |
| `DROP DATABASE d` | d's temporary tables Absent too |
| `DROP VIEW IF EXISTS v` | v a view when Present (1347) |
| `DROP FUNCTION f`, unqualified | as qualified (1305), but no claim with no database selected |
| `CREATE`/`ALTER`/`DROP EVENT` of a name with a period | ends the walk |

## 6. Not covered (beyond the MySQL spec's §7)

Review asks for these are declined with this list:

- restrictions OceanBase reports as not supported (1235, 4179): several foreign key or primary key changes in one `ALTER TABLE`, a foreign key added beside any other subcommand, an index dropped and one added together, `ALTER EVENT ... RENAME` of a started event, DDL on a materialized view;
- the pairs of `ALTER TABLE` subcommands on one name that OceanBase refuses and MySQL accepts (1060, 1061 above);
- `CREATE TABLE IF NOT EXISTS` and `CREATE TEMPORARY TABLE IF NOT EXISTS` of an existing table with a bad definition (1060);
- a foreign key to a view with the checks off, or to a base table a temporary table shadows (1215);
- `CREATE TRIGGER IF NOT EXISTS` of a name another table's trigger has (5735);
- the statements of materialized views, materialized view logs, external tables, sequences, table groups, and the recycle bin (`CREATE EXTERNAL TABLE`, `CREATE MATERIALIZED VIEW`, `FLASHBACK`, `PURGE`, ...): they are syntax errors in omni's parser, or, as `PURGE`, end the walk;
- the kind of a synced external table or materialized view: the walk claims no wrong kind of one (§3);
- `CREATE EVENT`, or `ALTER EVENT ... RENAME TO`, of a name the sync lists (1537): the listed event may be a dropped database's (§3);
- identifiers outside ASCII, as on MySQL (§3);
- `lower_case_table_names` 0 and 2, which a tenant can be created with but the oracle does not run (§8 O4).

## 7. Budget and tests

- **Code:** about 265 lines of deltas in `mysql/review` (lowering 170, snapshot 65, plumbing 30).
- **Script:** `scripts/oceanbase-oracle.sh` boots the image under the heavy-engine lock and the memory gate, sets `secure_file_priv` over the server's unix socket (the only way it can be set) so cases can create external tables, exports `OCEANBASE_WALK_DSN`, runs one command, and removes the container. The image needs `--ulimit nofile` above 20000, and boots in about a minute.
- **Oracle:** `TestWalkThroughAgainstOceanBase` reads `OCEANBASE_WALK_DSN` and skips without it; CI does not run it (engines design D4). It runs `walktest.Judge` with `Digester` (multi-object `DROP` and `RENAME` are atomic) on curated cases: the MySQL cases that apply, and one per row of §5 and per fact in the Appendix. The target is synced as Bytebase syncs OceanBase since bytebase/bytebase#21685 (§2): no CHECK constraints or partitions, external tables from `TABLES`, materialized views from `DBA_MVIEWS`, and since bytebase/bytebase#21687 events by name from `information_schema.EVENTS`, each once. The base holds an external table and a materialized view, and cases drop, rename, and alter them, so they show the walk claims no wrong kind of either. One case syncs the materialized views as a failed read does (none listed beside a synced storage table) and shows no relation absence is claimed in that database. One case drops and creates again a database that held an event and then creates the event again, so the sync lists the dropped one's and the walk must claim nothing. The MySQL cases for each assignment's scope in `SET` (bytebase/omni#515) run too: if a `SET GLOBAL foreign_key_checks` changes the session's checks on OceanBase, that becomes a delta.
  - The tenant's log stream lost its leader once in a probe run (4767, "log stream is not leader") under a burst of DDL; the oracle retries a statement on 4767, as the probe runner does.
- **Sweep:** the MySQL templates with OceanBase's names (`_obfk_`), at `WALK_SWEEP_N` (default 40). OceanBase's DDL is about ten times slower than MySQL's (100 ms for a table, 1 to 3 s for an index), so a sweep of 200 changes takes about 40 minutes; the PR's oracle run uses 200 (§8 O6).
- **Unit and coverage tests:** `TestWalkThroughOceanBaseWaits` becomes `TestWalkThroughOceanBase`, one case per delta, on the unit schema; `walktest.Coverage` over the MySQL forms with the OceanBase dialect.

## 8. Decisions

O1 is decided and O3 resolved; the rest are pending, each with a recommendation:

- **O1. Relation absence.** Decided 2026-10-28: trusted, since bytebase/bytebase#21685 syncs external tables and materialized views, except in a database whose synced storage tables outnumber its materialized views (§3). No `review.Target` switch: the metadata is current (§1). The options set aside were (a) never trusting it, as on MariaDB, and (b) trusting it while assuming no external table exists.
- **O2. A synced index named `PRIMARY`.**
  - (a) The primary key is Unknown: `ADD PRIMARY KEY` on a synced table is not claimed (1068).
  - (b) Trust it as the primary key, at the cost of a false finding on a keyless table with a secondary index named `PRIMARY`, which `CREATE INDEX` allows.

  Recommended: (a).
- **O3. A database with an event does not sync** (§2). Resolved 2026-10-28 by bytebase/bytebase#21687, which lists events by name; what it lists is O8.
- **O4. Oracle settings.** 4.3.5 with `lower_case_table_names` 1, the tenant default, only. A tenant with 0 or 2 needs a second tenant in the 8 GB budget, and the folding is MySQL's code, which the MySQL oracle runs with 0 and 1. Recommended: so.
- **O5. Restrictions** (§1). Recommended: claim none. The alternative claims `DROP INDEX` with `ADD INDEX` (4179), common when replacing an index, which 4.3.5 refuses.
- **O6. Sweep size.** Recommended: 200 changes in the PR's oracle run, 40 by default.
- **O7. `ALTER TABLE` of a view.** Recommended: check that it exists and end the walk. The alternative claims 1347 for the subcommands that fail on 4.3.5 (`ADD`/`DROP COLUMN`, `ADD INDEX`, `RENAME TO`, `ENGINE=`, `CONVERT TO`, `DROP FOREIGN KEY`).
- **O8. A synced event** (§3).
  - (a) Unknown, with absence trusted: a duplicate `CREATE EVENT` of a listed name is not claimed.
  - (b) Trusted as on MySQL, at the cost of a false finding on `CREATE EVENT e` in a database created again under the name of a dropped one that held e.

  Recommended: (a). It costs one rare finding, and nothing in the metadata tells the two apart.

## Appendix: verified facts

OceanBase CE 4.3.5.6, 2026-10-27; "MySQL" is `mysql:8.0` run on the same probes.

- **Settings:** `lower_case_table_names` 1, `foreign_key_checks` ON, `recyclebin` OFF, `default_storage_engine` OceanBase, `event_scheduler` OFF; no `sql_generate_invisible_primary_key`; `SHOW ENGINES` lists only InnoDB; `mysql.func` is empty.
- **ALTER TABLE, columns:**
  - fail: `ADD a, DROP a` and `DROP a, ADD a` (1060) where a exists; `ADD z, DROP z` and `DROP z, ADD z` (1091); `ADD x, RENAME x TO y` (1054); `DROP a, ADD INDEX (a)` (1072); `RENAME a TO b, DROP b` and `DROP b, RENAME a TO b` (1060); `RENAME a TO c, CHANGE a d` (1054); `CHANGE a c, MODIFY c` (1054); `DROP a, CHANGE a c` and `CHANGE a c, DROP a` (1091); `RENAME a TO c, DROP c` (1091); `DROP a, RENAME b TO a` (1060); `ADD c, MODIFY c` and `ADD y, ALTER COLUMN y SET DEFAULT` (1054); `RENAME a TO c, ADD c` (1060); two renames or changes to one name, `ADD z, RENAME a TO z`, `CHANGE a z, ADD Z` (1060); `CHANGE a a, CHANGE a a` and `MODIFY a, MODIFY a` (1054); a column added twice (1060) or dropped twice (1091);
  - succeed: `RENAME a TO c, ADD a`; `CHANGE b d, ADD b`; `ADD INDEX k (q), ADD q`; `ADD y1, ADD y2 AFTER y1`; `CHANGE a a2, ADD y4 AFTER a2`; `DROP a, ADD y5 AFTER a` (y5 takes a's place); swaps by `RENAME COLUMN` and by `CHANGE`; `RENAME a TO b, RENAME b TO y` (leaves b and y); `CHANGE b b2, ADD b`; `CHANGE a A, ADD d`; `ADD y1, CHANGE a a2`; `MODIFY a, DROP a`, `DROP a, MODIFY a`, and `DROP a, ALTER COLUMN a SET DEFAULT` (a dropped); dropping a column drops it from its indexes, and an index with it alone.
- **ALTER TABLE, indexes and keys:**
  - fail: `ADD INDEX i, DROP INDEX i` and `DROP INDEX j, ADD INDEX jj` and `DROP INDEX i, ADD UNIQUE`/`ADD CONSTRAINT ... UNIQUE` (4179); `ADD INDEX z, DROP INDEX z` (1091); `RENAME INDEX i TO j, DROP INDEX j` and an index swap (1061); `RENAME INDEX i TO k, RENAME INDEX k TO m` (1176); `RENAME INDEX i TO I` (1061); an index added twice (1061), dropped twice (1072); `DROP PRIMARY KEY` twice, or with `` DROP INDEX `PRIMARY` `` or `ADD INDEX` (1235); `DROP INDEX i, ADD PRIMARY KEY` (1235); `ADD PRIMARY KEY` twice or two inline primary keys (1068), `CHANGE ... PRIMARY KEY, ADD PRIMARY KEY` (1235); `ALTER INDEX nosuch` (1176);
  - succeed: `RENAME INDEX i TO k, ADD INDEX i`; `DROP PRIMARY KEY, ADD PRIMARY KEY`; `ALTER INDEX i INVISIBLE, DROP INDEX i`; `DROP INDEX i, ADD COLUMN`; `DROP INDEX i, RENAME INDEX j TO k`; `DROP INDEX i, DROP INDEX j`; `ADD INDEX k, ADD INDEX m`; `ADD PRIMARY KEY` on a keyless table; `RENAME TO` with other subcommands;
  - `ORDER BY`, `ENABLE`/`DISABLE KEYS`, `ALTER COLUMN ... SET [IN]VISIBLE`, and `DISCARD TABLESPACE` are syntax errors (1064); `ENGINE=` any engine is a warning (1286); `ALTER TABLE ... RENAME TO` a missing database fails with 1210.
- **The primary key and PRIMARY:** `DROP PRIMARY KEY` and `` DROP INDEX `PRIMARY` `` of a keyless table fail (1091); `CREATE TABLE` with an index named `PRIMARY` fails (1280); `` CREATE INDEX `PRIMARY` `` on a keyless table and `` ADD INDEX `PRIMARY` `` on a keyed one succeed, `` RENAME INDEX `PRIMARY` TO k `` renames it, and `` DROP INDEX `PRIMARY` `` drops it, but fails with only a primary key (1091). On a keyless table, `__pk_increment` is not listed; `DROP COLUMN` (1235), `MODIFY`, `CHANGE` (1054) of it fail; `AFTER` it, `ALTER COLUMN ... SET DEFAULT` on it, and an index on it succeed; `ADD COLUMN __pk_increment` and a keyless `CREATE TABLE` with that column fail (1060).
- **Foreign keys:**
  - a name another table has fails with 1022, either case; the same name in another database succeeds; `DROP FOREIGN KEY` of another table's name fails (1091), of a case variant drops it;
  - `DROP FOREIGN KEY f, DROP FOREIGN KEY f` succeeds; two of one name in `CREATE TABLE` or `ALTER TABLE` fail (1215); two in one `ALTER TABLE`, or one beside `ADD INDEX`, `DROP FOREIGN KEY`, or `ADD COLUMN` fail (1235, 5031);
  - generated names are `<table>_OBFK_<16 digits>`; `RENAME TABLE` within and across databases, and `ALTER TABLE ... RENAME`, keep every name; a foreign key adds no index; `LIKE` copies none; inline `REFERENCES` makes none;
  - an unqualified parent resolves in the current database, for `CREATE TABLE` and `ALTER TABLE`, in both directions (1146 where the current database lacks it);
  - a missing parent fails with 1146 and is accepted with `foreign_key_checks=0`, and enforced once created; a view or temporary table parent fails with 1215, a view even with the checks off; a missing column with 5031;
  - `ENGINE=MyISAM` keeps foreign keys (1452 on a bad row); `DROP TABLE` of a referenced parent fails (3730);
  - `DROP CONSTRAINT` drops a foreign key or a CHECK and fails for a unique key, an index, or `PRIMARY` (3940).
- **Generated names:** unnamed indexes `a`, `a_2`, `a_3`; CHECK `<table>_OBCHECK_<16 digits>`, unique per database (3822).
- **Drops, renames, and kinds:**
  - `DROP TABLE t1, nosuch` and `DROP TABLE t5, nosuch, t6` fail (1051) and keep every table; `DROP TABLE t2, v` and `DROP VIEW v, nosuch` and `DROP VIEW v, t` keep every object; `DROP TEMPORARY TABLE tt1, nosuch` keeps tt1;
  - `RENAME TABLE t1 TO t4, nosuch TO t5` fails (1017) and keeps t1; a swap through a third name succeeds; `RENAME TABLE t TO t` fails (1050); a view moves to another database; a table with a trigger moves to another database, its trigger staying in the old one;
  - `DROP TABLE` of a view (1051), `DROP VIEW` of a table (1347), `DROP VIEW IF EXISTS` of a table (1347), `DROP TABLE IF EXISTS` of a view (a note), `TRUNCATE` of a view (1146), `ALTER TABLE` of a view (1347; `COMMENT` succeeds), `CREATE INDEX`, `CREATE OR REPLACE VIEW`, `ALTER VIEW`, and `LIKE` on the wrong kind (1347), a trigger on a view (1235);
  - `ANALYZE TABLE` of a missing table fails (1146); `OPTIMIZE` and `CHECK` return.
- **CREATE TABLE:** a failed `CREATE TABLE ... SELECT` leaves no table; `IF NOT EXISTS ... LIKE nosuch` fails with 1146 whether the table exists or not; `IF NOT EXISTS` of an existing table or view is a note, but with a bad definition fails (1060); a temporary table of the name does not block `CREATE TABLE`.
- **Temporary tables:** shadow a base table for reads, `ALTER TABLE`, `CREATE INDEX`, and `DROP TABLE`; are not in `information_schema`; `DROP TEMPORARY TABLE` of a base table fails (1051); `RENAME TABLE` and `ALTER TABLE ... RENAME` of one succeed; `DROP DATABASE` drops them; `CREATE TEMPORARY TABLE` in a missing database fails (1049), with none selected 1046.
- **Triggers, routines, events:**
  - a trigger of a taken name fails with 5734 on its table, 5735 on another, also with `IF NOT EXISTS` on another; `CREATE OR REPLACE TRIGGER` is a syntax error; `DROP TRIGGER TR` drops tr; a trigger in another database than its table's fails (1435);
  - a trigger or function with `CREATE TEMPORARY TABLE` or `DROP TEMPORARY TABLE` fails to create (1422), and one that `CALL`s a procedure that creates one fails when it runs (1422); a procedure may run DDL;
  - a cascaded `ON DELETE CASCADE` delete fires the child's trigger;
  - procedure and function names compare without case (1304); `IF NOT EXISTS` is a warning; `DROP FUNCTION nosuch` fails with 1305; loadable functions cannot be created (1235); with no database selected, `DROP FUNCTION f` fails with 1305 and `DROP FUNCTION IF EXISTS f` succeeds;
  - event names compare without case (1537); `` `shop.e3` `` is e3 in shop; events in another database and renames of a running one are not supported (1235); `SHOW CREATE EVENT` and `SHOW EVENTS` are syntax errors.
- **Names:** with `lower_case_table_names=1`, table and database names are stored in lower case and compare without case; every name compares accent-insensitively (`é` = `e` = `É` for tables, columns, indexes, databases, procedures, triggers), `ß` and `ss` differ; lengths as §4.7, counted in characters; a table name with a trailing space fails (1103).
- **No database selected:** after `DROP DATABASE` of the current database, and after creating it again, unqualified `CREATE`, `DROP`, `ALTER`, DML, and `RENAME` fail with 1046, `IF EXISTS` too (but `DROP FUNCTION IF EXISTS`); `ALTER DATABASE` naming none fails with 1210.
- **Other kinds:** an external table (with `secure_file_priv` set over the unix socket) is listed as `EXTERNAL TABLE`; `CREATE TABLE` and `CREATE VIEW` of its name fail (1050); `RENAME TABLE` and `DROP TABLE` of it succeed; `ALTER TABLE`, `CREATE INDEX`, and `INSERT` fail (1235). A materialized view is in neither `information_schema.TABLES` nor `VIEWS`; its container `__mv_container_<id>` and a materialized view log `mlog$_<table>` are `BASE TABLE`s; `CREATE TABLE mv` and `CREATE VIEW mv` fail (1050); `RENAME TABLE mv`, `CREATE INDEX` and `DROP INDEX` on it, and `CREATE TABLE ... AS SELECT` from it succeed; `ALTER TABLE`, `TRUNCATE`, DML, and triggers on it fail (1235). Sequences are a namespace of their own (4297 for a taken sequence name, 4298 for a missing one). With `recyclebin` ON, a dropped table's name is free and `FLASHBACK` restores it.
- **Sync:** Bytebase's queries run as written, but `SHOW CREATE EVENT` (1064), which bytebase/bytebase#21687 replaced with a name-only read of `information_schema.EVENTS`. That view keeps a dropped database's events under its name: a database created again under it lists them, lists an event created again twice, and accepts `CREATE EVENT` of such a name (verified in #21687). bytebase/bytebase#21685 verified on 4.3.5.6: every object `oceanbase.DBA_OBJECTS` lists in a database, but indexes, makes `CREATE TABLE` of its name fail with 1050, and the sync lists each; `CREATE SYNONYM` is a syntax error in MySQL mode; each materialized view has exactly one storage table; a user without `SELECT` on `oceanbase` cannot read `DBA_MVIEWS` (1142) but sees the storage table in `TABLES`.
