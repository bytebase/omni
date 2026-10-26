# WALK_THROUGH for MariaDB and TiDB: v1 scope

Status: draft 2026-10-26, decisions in §6 pending. This is step 3 of [the engines design](2026-10-22-walk-through-engines-design.md) §7, which its §5.2 and §5.3 sketch. It is written as deltas from [the MySQL v1 scope](2026-10-22-walk-through-mysql.md): what that spec says holds here unless this one says otherwise. Every fact here was verified on `mariadb:11.8.8` and `pingcap/tidb:v8.5.5`, the images the packages' oracle tests already use (Appendix).

## 1. Shape

- **One copy per engine.** Each engine gets the MySQL walk as a copy in its own review package (`mariadb/review`, `tidb/review`), since an engine never imports another. The copies keep MySQL's file layout (`walk_through.go`, `walk_snapshot.go`, `walk_lower*.go`, and the three test files), so a fix in one ports to the others by diff.
- **Deltas only where the engine differs.** Each copy changes only what §2 or §3 lists, and every delta names its verified fact in a comment.
- **The compile delta is small:** MariaDB's DROP-form classifier has no `dropCheck`, and TiDB's `AlterTableCmd` has no `NewTable`.
- **Bytebase is not changed** (decided 2026-10-26). The walk runs on what `review.Target` carries today. `ForeignKeyChecks` and `OtherSpecs` stay unset until Bytebase fills them, and the walk behaves as it does on MySQL in that case.

## 2. MariaDB (11.8)

Built in `mariadb/review` (2026-10-26). The facts this section adds to the draft, and the two it corrects, were verified on `mariadb:11.8.8` while building it; §2.5 lists them.

### 2.1 Authority

- **Relation absence is untrusted** (engines design D5, decision N2). Before Bytebase #21665 the sync listed only `BASE TABLE`s and views, and sequences and system-versioned tables share their namespace. The target does not say which sync it comes from. So:
  - no claim that a synced database lacks a relation;
  - "already exists" from the sync is still trusted, for the sequences #21665 lists too;
  - so are the columns and indexes of synced tables, but for a system-versioned table's columns (below);
  - names the change itself creates or drops are followed as on MySQL.
- **Trigger absence is untrusted too.** The sync before #21665 lost the triggers of the tables it skipped, so `DROP TRIGGER` of one would be a false finding. Foreign key names stay trusted: the walk claims a missing one only with the table's own foreign keys, which the sync lists with the table.
- **Primary key absence is trusted.** MariaDB has no invisible primary key, so `DROP PRIMARY KEY` on a table the sync lists without one fails (1091).
- **System-versioned tables:** the sync since #21665 lists one as a table, with `WITH SYSTEM VERSIONING` in its create options. Its columns are Unknown: `information_schema` does not list the hidden `row_start` and `row_end`.

### 2.2 Semantics that differ from MySQL

- **Multi-object drops are not atomic.**
  - `DROP TABLE a, missing`, `DROP VIEW v, missing`, `DROP SEQUENCE s, missing`, and `DROP TEMPORARY TABLE a, missing` drop every object that exists, then fail on the missing one (1051, 4092, 4091). `DROP TABLE t, v` drops t and fails on the view (1965).
  - A failed `CREATE OR REPLACE TABLE t` may have dropped the old t: it has after a duplicate column (1060) and a foreign key name taken (1005, errno 121), and has not after a missing `LIKE` source (1146).
  - The claim that the statement fails still holds. The oracle stops judging the change there (decision N1).
- **`CREATE TEMPORARY TABLE db.t` and `CREATE TEMPORARY SEQUENCE db.s` succeed in a missing database.** The walk checks no database for a temporary relation. With no database selected, they fail (1046).
- **`CREATE TABLE IF NOT EXISTS t LIKE missing` is a note when t exists.** MySQL checks the source first. Here the source is checked only when the table is created. A temporary table's source is checked first, as on MySQL (1146 even when the temporary table exists).
- **Trigger names compare with case.** `TR` can exist beside `tr`, as on 5.7. MySQL's spelling rule already covers this.
- **`DROP FOREIGN KEY` of a missing name on MyISAM fails** (1091), where MySQL accepts it. The walk still claims nothing of a non-InnoDB table's foreign keys, so this costs recall only.
- **`CREATE OR REPLACE DATABASE d` can fail when d holds a parent another of its tables references** (1451): it fails at such a parent while its child is still there, having dropped some of d's tables. `DROP DATABASE` drops them all. The walk does not claim it.
- **MySQL's `ALTER TABLE` forms that MariaDB lacks fail with 1064:** `ALTER INDEX ... [IN]VISIBLE`, `ALTER COLUMN ... SET [IN]VISIBLE` and `SET/DROP NOT NULL`, `ALTER CHECK ... ENFORCED`, and `SECONDARY_LOAD`. omni's MariaDB parser takes them, so the walk ends at them rather than claim a failure for another reason.
- **Error codes differ.** This changes only the oracle's code map:

  | Case | MariaDB code |
  | --- | --- |
  | FK name taken | 1005 (errno 121) |
  | FK to a missing table, a view, a temporary table, or a sequence | 1005 (errno 150) |
  | two FKs of one name in one statement | 1061 |
  | `DROP TABLE` of a view | 1965 |
  | `DROP VIEW` of a table, a sequence, or a missing view | 4092 |
  | database name too long | 1102 |
  | table, view, or sequence name too long | 1103 |
  | rename into a missing database | 1025 |

### 2.3 Statements added

| Statement | Checks | Effects |
| --- | --- | --- |
| `CREATE OR REPLACE TABLE t ...` | t not a view (1965); a sequence is replaced | t Present (table) with what it declares; what the old t held Absent, its foreign keys and triggers too, so a foreign key name only the old t had is free; a temporary table of the name stays |
| `CREATE OR REPLACE INDEX`, `TRIGGER`, `PROCEDURE`, `FUNCTION`, `EVENT` | none of the "already exists" checks | name Present |
| `CREATE OR REPLACE DATABASE d` | none | d Present, everything under it Absent but temporary tables; when d is the current database, none is selected after it (1046), as after `DROP DATABASE` |
| `ALTER TABLE t DROP {COLUMN, INDEX, FOREIGN KEY, CONSTRAINT} IF EXISTS x` | no presence check; a second such drop of x is a note, so the repeated-subcommand checks leave it out | x Absent; `DROP FOREIGN KEY IF EXISTS` of a name only another table has leaves that one |
| `CREATE [OR REPLACE] [TEMPORARY] SEQUENCE [IF NOT EXISTS] s` | s Absent (1050), skipped for OR REPLACE and IF NOT EXISTS; for OR REPLACE, s not a view (1965); s's database Present (1049), unless temporary | s Present, of kind sequence; OR REPLACE replaces a table or a temporary table |
| `DROP [TEMPORARY] SEQUENCE [IF EXISTS] s, ...` | each s Present (4091) and a sequence (4090; 1965 for a view). Without TEMPORARY, a temporary sequence of the name is found first and a temporary table of it passed over; with it, only a temporary sequence is found (4091) | s Absent |
| `DROP TABLE` | as on MySQL, with sequences taken as tables, temporary ones too | a sequence goes too |
| `CREATE TABLE t LIKE s` | s a table or a sequence (1347 for a view) | t is what s is: a sequence of a sequence, a system-versioned table of one, but never a temporary system-versioned table (4137) |

OR REPLACE with IF NOT EXISTS fails on the server (1221), and omni's parser rejects it too, so it is a Syntax finding.

Other kinds:
- `DROP VIEW` of a sequence fails (4092), and `CREATE OR REPLACE VIEW` of one (1347).
- A foreign key to a sequence fails (1005, errno 150).
- A known sequence in `ALTER TABLE`, `TRUNCATE`, `CREATE INDEX`, `DROP INDEX`, or `CREATE TRIGGER` ends the walk without a claim. Some `ALTER TABLE`s on a sequence succeed, so no claim is safe there.
- `ALTER SEQUENCE` is not in the table, so it ends the walk.

**System-versioned tables:** `WITH SYSTEM VERSIONING`, of the table or of a column, or `ADD SYSTEM VERSIONING`, gives a table hidden `row_start`/`row_end` columns that `information_schema` does not list. By default it also refuses most ALTERs (4119). The hidden columns act as absent to every check the walk makes: only `ADD COLUMN` of their names fails (1060), which the walk never claims, and `CHANGE`/`MODIFY` of them fails with 4130. So the MySQL lowering would hold here, but would name a column where the server refuses the versioned table. As a precaution, the walk follows neither:
- such a table's columns are Unknown;
- `ALTER TABLE` on it ends the walk, and so do `ADD`/`DROP SYSTEM VERSIONING` and the `PERIOD` subcommands on any table;
- its index names are followed: `CREATE` and `DROP INDEX` on it check them as on any table (1061, 1091);
- `LIKE` copies the versioning, `RENAME` keeps it, and `CREATE OR REPLACE TABLE` without it ends it.

### 2.4 Parser gaps (Syntax findings today; out of scope)

The server accepts these, and omni's MariaDB parser rejects them:
- `ALTER TABLE IF EXISTS`;
- `ADD COLUMN`/`ADD INDEX`/`ADD FOREIGN KEY ... IF NOT EXISTS`;
- `DROP INDEX IF EXISTS ... ON t`;
- `RENAME`/`CHANGE`/`MODIFY COLUMN IF EXISTS`;
- `RENAME TABLE IF EXISTS`;
- `ALTER ONLINE TABLE`;
- `WAIT n`;
- `CREATE VIEW IF NOT EXISTS`;
- `SET STATEMENT ... FOR`;
- `ALTER INDEX ... [NOT] IGNORED`;
- `CREATE TABLE ... SEQUENCE=1`;
- `BEGIN NOT ATOMIC` and `EXECUTE IMMEDIATE`;
- qualified event and trigger names, as in omni's MySQL parser.

The walk lowers each of them once the parser takes it. The parser also takes the MySQL forms of §2.2 that MariaDB rejects.

### 2.5 Verified while building

On `mariadb:11.8.8`, 2026-10-26:
- **Corrections to the draft:**
  - `CREATE OR REPLACE TABLE` of a sequence replaces it with a table;
  - `CREATE OR REPLACE DATABASE` of the current database leaves none selected.
- **OR REPLACE:**
  - `CREATE OR REPLACE TABLE t` with a temporary t replaces the base t and leaves the temporary one;
  - the old t's foreign key names are free for the new t, and another table's still fail (1005, errno 121);
  - its triggers go with it;
  - `CREATE OR REPLACE TRIGGER tr` of tr on another table fails with 1360, and `TR` beside `tr` is a second trigger;
  - `CREATE OR REPLACE EVENT E` replaces e;
  - `CREATE OR REPLACE DATABASE` keeps the session's temporary tables of the database;
  - OR REPLACE of a temporary table and a temporary sequence replaces either kind.
- **IF EXISTS drops:**
  - two of one name are a note; one with and one without fail (1091);
  - `DROP INDEX IF EXISTS `` `PRIMARY` ``` drops the primary key;
  - `DROP PRIMARY KEY IF EXISTS` is a syntax error.
- **Sequences:**
  - `CREATE SEQUENCE db.s` in a missing database gives 1049;
  - `RENAME TABLE` and `ALTER TABLE ... RENAME` move a sequence, to another database too;
  - `CREATE TABLE ... AS SELECT` from a sequence makes a table;
  - a temporary sequence and a temporary table share their names (1050);
  - `DROP TEMPORARY TABLE` and `DROP TABLE` drop a temporary sequence;
  - `CREATE INDEX` on a sequence gives 4086, `DROP INDEX` 1091, `CREATE TRIGGER` 1361, `TRUNCATE` 1031, and `ALTER TABLE ... COMMENT` succeeds;
  - `INSERT` into a sequence succeeds.
- **System-versioned tables:**
  - `CREATE INDEX` and `ADD INDEX` on one succeed, `DROP INDEX` of a missing index gives 1091, and `TRUNCATE` 4137;
  - with `system_versioning_alter_history=KEEP`, `AFTER row_start`, `RENAME COLUMN row_start`, and `ALTER COLUMN row_start SET DEFAULT` give 1054, `DROP COLUMN row_end` 1091, `ADD INDEX (row_start)` 1072, `ADD COLUMN row_start` 1060, and `CHANGE` or `MODIFY row_start` 4130;
  - a column's `WITH SYSTEM VERSIONING` makes the table system-versioned, and on a table that is not, in `ALTER TABLE`, fails (4124);
  - `PERIOD FOR SYSTEM_TIME` without `WITH SYSTEM VERSIONING` gives 4125, and `PARTITION BY SYSTEM_TIME` 4124;
  - `CREATE TEMPORARY TABLE ... WITH SYSTEM VERSIONING` gives 4137, and `LIKE` one strips the versioning.
- **Other:**
  - names over 64: databases 1102, tables, views, and sequences 1103, the rest 1059;
  - `CREATE TABLE x LIKE v` of a view gives 1347;
  - with no database selected, `CREATE TEMPORARY TABLE`, `CREATE SEQUENCE`, and `DROP SEQUENCE IF EXISTS` give 1046;
  - `DELETE HISTORY` of a missing table gives 1146, and of a table that is not system-versioned 4124.

## 3. TiDB (8.5)

### 3.1 Names

`lower_case_table_names` is 2 and fixed. Databases and relations fold, whatever `Target.LowerCaseTableNames` says. The other spaces fold as on MySQL.

### 3.2 Objects TiDB does not have

- Triggers, functions, and events are syntax errors.
- `CREATE PROCEDURE` parses and fails (8108).
- None is lowered: such a statement ends the walk, and the snapshot lists none.

### 3.3 ALTER TABLE

- **Every subcommand is judged against the table as it was,** and nothing runs first. So TiDB's ALTER uses `walk.Lookups`, with no touched rule. Verified:
  - `ADD a, DROP a` fails with 1060 where a exists;
  - a swap of two columns fails with 1060;
  - `RENAME INDEX k1 TO k2, DROP INDEX k2` fails with 1061;
  - `AFTER` a column the statement adds fails with 1054;
  - an index on a column the statement adds fails with 1072 (not claimed: index columns are not modeled).
- **Two subcommands on one column or index fail** with 8200 ("Unsupported operate same column"). MySQL's repeated-subcommand text checks hold, with 8200 among their codes.
- **`RENAME TO` with any other subcommand fails** (8200), from the text alone.
- **ALTER or RENAME of a local temporary table fails** (8200). The walk ends there, as on MySQL.
- **`ALTER TABLE t RENAME TO db.t`** is a syntax error in omni's TiDB parser (gap). A quoted `` `db.t` `` is one name.

### 3.4 Foreign keys

- **Names are unique per table, not per database.** Two foreign keys of one name on one table fail (1826, either case), and a name another table has is accepted. TiDB keeps only the table's foreign key key, not the database-wide one.
- **Generated names are `fk_<n>`,** a May-effect on the table. RENAME TABLE does not rename them.
- **TiDB ignores `ENGINE`,** so every table keeps its foreign keys.

### 3.5 Other

- Multi-table `DROP` and `RENAME` are atomic.
- `DROP VIEW IF EXISTS t` of a table fails (1347): IF EXISTS does not cover the wrong kind. The walk claims it.
- `DROP TABLE` of a sequence fails (1051). Sequences are synced, so they are relations of kind sequence that the walk trusts. `CREATE`/`DROP SEQUENCE` are lowered as in §2.3, with MySQL's codes.
- `CREATE TABLE ... SELECT` fails on 8.5 (1105, not implemented). It is lowered as on MySQL, so a server that takes it is followed, and 8.5's rejection is missed.
- `_tidb_rowid`, the hidden row id of a table without a clustered integer key, is Unknown, as MySQL's `my_row_id` is.
- `CREATE GLOBAL TEMPORARY TABLE` is a syntax error in omni's TiDB parser (gap). The server lists such tables as `BASE TABLE`, so they sync as tables.

## 4. Not covered (beyond the MySQL spec's §7)

- **MariaDB:**
  - system-versioned history and its ALTER rules (4119), `TRUNCATE` of a system-versioned table (4137), `DELETE HISTORY` of a table that is not (4124), and application-time periods;
  - the rejections of `ALTER TABLE` and `TRUNCATE` on a sequence, and `ALTER SEQUENCE`;
  - `CREATE OR REPLACE TRIGGER` of a name on another table (1360);
  - `CREATE OR REPLACE TABLE` of a referenced parent (1451), and `CREATE OR REPLACE DATABASE` of a database holding one (1451);
  - `DROP TRIGGER` of a case variant of a trigger's name (1360), and of a trigger the sync does not list;
  - the MySQL `ALTER TABLE` forms MariaDB rejects (1064), at which the walk ends.
- **TiDB:**
  - clustered-index rules (`DROP PRIMARY KEY` of a clustered key, 8200);
  - index columns added in the same ALTER (1072);
  - placement policies and TiFlash replicas;
  - stored procedures (8108).

## 5. Budget and tests

- **Code:** per engine, the MySQL walk copied (about 1.9k lines) plus its deltas: about 250 lines for MariaDB, about 150 for TiDB.
- **Oracle:** `walktest.Judge` against `mariadb:11.8.8` and `pingcap/tidb:v8.5.5`, with curated cases for every delta row and every Appendix fact.
- **Sweep:** one per engine at `WALK_SWEEP_N`. The MySQL sweep takes about 3 minutes at 500, so each engine adds about 1.5 minutes.
- **Unit and coverage tests:** as on MySQL.
- **MariaDB as built:**
  - the walk is 2462 lines against MySQL's 2172, so the deltas are about 290 lines, over the 250 budgeted: sequences took 140;
  - the oracle runs every curated case on the sync since Bytebase #21665 and on the one before it; the second holds precision alone, since it hides the base's sequence and system-versioned table. It found 194 of 230 rejections on the first and 181 of 228 on the second, with no false finding. On the plain MySQL copy, wired the same way, the delta cases report false findings or miss what they pin;
  - the sweep runs the two syncs by turns. At 500 (seed 20261026) it found 1238 of 6811 rejections, with no false finding, in under 3 minutes. Recall is low by design: most misses are names the walk does not know, since relation and trigger absence is untrusted, and databases other than the target;
  - `walktest` gained `PartialRejecter` (decision N1).

## 6. Decisions

Decided 2026-10-26 for MariaDB: N1, N2, and N3 as recommended, and N4 as one PR per engine, since the two are built in parallel. TiDB's parts are pending. The options as they were weighed:

- **N1. MariaDB's non-atomic drops in the oracle.**
  - Recommended: `walktest` gains an opt-in for an engine whose rejected statement can change the database. `Judge` then stops judging a change at such a statement, after counting it, instead of failing.
  - Alternative: keep `Judge` strict, and keep multi-object drops of a missing name out of MariaDB's cases and sweep.
  - As built: `walktest.PartialRejecter`. A finding on such a statement must name the server's error code, since the catalog no longer shows the facts before it. A failed `CREATE OR REPLACE TABLE` that dropped the old table is judged the same way.
- **N2. MariaDB relation absence stays untrusted (D5)** while Bytebase is not changed.
  - Recommended: so, at the cost of every "table does not exist" claim about a synced MariaDB database.
  - The alternative, trusting it, has a false finding for every sequence or system-versioned table the sync drops.
  - As decided: Bytebase #21665 now syncs both, but the target does not say which sync it comes from. A `review.Target` switch, off by default, may later let Bytebase say so. Trigger absence waits for the same switch (§2.1).
- **N3. Oracle settings.**
  - Recommended: MariaDB with `lower_case_table_names` 0 only, since the folding is MySQL's code, which runs on mysql 8.0 with 1.
  - TiDB as it ships (2).
- **N4. Delivery.**
  - Recommended: one PR for both engines, as D7's step 3, with a commit each.
  - The parser gaps of §2.4 and §3.3/§3.5 go to their own PRs afterwards.

## Appendix: verified facts

MariaDB 11.8.8 and TiDB v8.5.5, 2026-10-26.

- **MariaDB, ALTER TABLE.** It runs drops first, as MySQL does:
  - `ADD z, DROP z` fails with 1091;
  - `ADD a, DROP a` succeeds where a exists;
  - a column swap keeps both;
  - `RENAME INDEX k1 TO k2, DROP INDEX k2` leaves k2;
  - `ADD INDEX kq (q), ADD q` succeeds.

  Repeated forms fail:
  - a second DROP COLUMN, DROP FOREIGN KEY, or DROP PRIMARY KEY with 1091;
  - a second ADD COLUMN with 1060;
  - two renames to one name with 1060;
  - `PRIMARY` as an index name with 1280;
  - a second primary key with 1068.
- **MariaDB, foreign keys.**
  - A name another table's FK has fails with 1005 (errno 121), either case.
  - Two of one name in a CREATE fail with 1061.
  - An FK to a missing table, a view, or a temporary table fails with 1005 (errno 150), and `foreign_key_checks=0` accepts the missing one.
  - `DROP FOREIGN KEY` of another table's name fails with 1091.
  - `RENAME TABLE child TO kid` renames `child_ibfk_1` to `kid_ibfk_1`.
  - MyISAM and Aria accept an FK to a missing table and keep none. MyISAM fails `DROP FOREIGN KEY` of a missing name with 1091.
- **MariaDB, drops and kinds.**
  - `DROP TABLE d1, nosuch` fails with 1051 and drops d1. `DROP TABLE d5, nosuch, d6` drops d5 and d6.
  - `DROP VIEW v, nosuch` fails with 4092 and drops v.
  - `DROP TEMPORARY TABLE tt1, nosuch` drops tt1.
  - `RENAME TABLE d2 TO d3, nosuch TO d4` fails with 1146 and keeps d2.
  - `DROP TABLE` of a view gives 1965, `DROP VIEW` of a table 4092, and `TRUNCATE` of a view 1146. `ALTER TABLE` and `CREATE INDEX` on a view give 1347.
  - `DROP VIEW IF EXISTS t` of a table and `DROP TABLE IF EXISTS v` of a view are notes.
- **MariaDB, other behavior.**
  - Case: `CREATE TRIGGER TR` beside `tr` succeeds. `P` beside `p` fails with 1304, and `E` beside `e` with 1537.
  - Temporary tables: one survives `DROP DATABASE` of its database. `CREATE TEMPORARY TABLE gone.tmp` in a missing database succeeds. One shadows a base table as on MySQL.
  - `OR REPLACE`:
    - `CREATE OR REPLACE TABLE` replaces a table; on a view it fails with 1965, on a referenced parent with 1451;
    - with IF NOT EXISTS it fails with 1221, and with `LIKE` a missing table with 1146;
    - `OR REPLACE` of an index, trigger, procedure, function, event, database, or sequence succeeds, except a trigger of a name on another table (1360).
  - Sequences:
    - a table, sequence, or view of a sequence's name fails with 1050, and `CREATE SEQUENCE IF NOT EXISTS` is a note;
    - `DROP SEQUENCE` of a table gives 4090, and of a missing one 4091;
    - `DROP TABLE` drops a sequence;
    - `DROP VIEW` of one gives 4092, `TRUNCATE` 1031, and `ALTER TABLE ... ADD COLUMN` 4086;
    - information_schema lists `TABLE_TYPE = 'SEQUENCE'`.
  - System versioning: information_schema lists `SYSTEM VERSIONED` and no `row_start`/`row_end`. ADD or DROP COLUMN fails with 4119.
  - IF EXISTS forms: `DROP COLUMN`/`DROP INDEX`/`DROP FOREIGN KEY ... IF EXISTS` of a missing one, and `ALTER TABLE IF EXISTS` of a missing table, are notes.
  - Names and databases:
    - a column or index name over 64 gives 1059, a table name 1103;
    - moving a view to another database gives 1450, and a table with a trigger 1435;
    - a rename into a missing database gives 1025;
    - with no database selected, CREATE TABLE, ALTER DATABASE, and CREATE PROCEDURE give 1046.
- **TiDB, settings and objects.**
  - `foreign_key_checks` is 1 by default, `lower_case_table_names` 2, and clustered indexes on.
  - `CREATE TABLE T` beside t, `CREATE DATABASE SHOP` beside shop, `ADD COLUMN A` beside a, and `CREATE VIEW V` beside v all fail.
  - Triggers, functions, and events are 1064; procedures 8108.
- **TiDB, ALTER TABLE.**
  - What fails:
    - `ADD z, DROP z` gives 1091 and `ADD a, DROP a` 1060;
    - a column swap and `CHANGE b b2, ADD b` give 1060;
    - `DROP c, ADD c` gives 1060, and `DROP INDEX i, ADD INDEX i` 1061;
    - two renames to one name, a second DROP, a second ADD, and `MODIFY a, MODIFY a` give 8200;
    - `RENAME INDEX k1 TO k2, DROP INDEX k2` gives 1061;
    - `ADD y1, ADD y2 AFTER y1` and `CHANGE a a2, ADD y4 AFTER a2` give 1054;
    - `ADD INDEX kq (q), ADD q` and `ADD e, ADD INDEX ke (e)` give 1072;
    - `RENAME TO` with another subcommand gives 8200;
    - `DROP PRIMARY KEY` of a clustered key gives 8200;
    - a second primary key gives 1068;
    - ALTER or RENAME of a local temporary table gives 8200.
  - What succeeds: `DROP INDEX j, ADD INDEX jj (b)`; `ADD y1, CHANGE a a2`; `ADD PRIMARY KEY` on a keyless table.
- **TiDB, foreign keys.**
  - `ADD CONSTRAINT fk1` on a table whose fk1 another table has succeeds, as does `CREATE TABLE x (CONSTRAINT fk1 ...)`.
  - These fail with 1826: `FK1` on the table that has fk1; two of one name in one statement; `DROP FOREIGN KEY fk_1, ADD CONSTRAINT fk_1`.
  - Generated names are `fk_1`, and `RENAME TABLE child TO kid` keeps them.
  - An FK to a missing table or a temporary table gives 1824, and `foreign_key_checks=0` accepts the missing one.
  - `ENGINE=MyISAM` is ignored, so the FK is enforced (1824).
- **TiDB, drops, temporary tables, and the rest.**
  - Atomicity: `DROP TABLE d1, nosuch` and `RENAME TABLE d2 TO d3, nosuch TO d4` fail and keep everything. A swap through a temporary name succeeds.
  - Views: `DROP TABLE` of a view gives 1051, `DROP VIEW` of a table 1347, and `DROP VIEW IF EXISTS` of a table 1347. `ALTER TABLE` and `CREATE OR REPLACE VIEW` of the wrong kind give 1347.
  - Temporary tables: one shadows a base table and survives `DROP DATABASE`. `CREATE TEMPORARY TABLE gone.tmp` gives 1049. A `GLOBAL TEMPORARY` table is a `BASE TABLE`, and `CREATE TABLE` of its name gives 1050.
  - Sequences: a table or sequence of a sequence's name gives 1050. `DROP SEQUENCE` of a table gives 1347, and `DROP TABLE` of a sequence 1051.
  - Other statements:
    - `CREATE TABLE ... SELECT` gives 1105;
    - `CREATE TABLE IF NOT EXISTS t LIKE missing` gives 1146, as on MySQL;
    - names over 64 give 1059;
    - a rename into a missing database gives 1025, and moving a view to another database 1450;
    - with no database selected, CREATE TABLE and ALTER DATABASE give 1046.
