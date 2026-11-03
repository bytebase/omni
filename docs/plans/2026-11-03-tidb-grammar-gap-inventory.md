# TiDB grammar gap inventory

Date: 2026-11-03

Goal: make the distance between `omni/tidb/parser` and the TiDB grammar a
number this repository measures, and group what is missing so the parser
PRs that close it have a list to work down.

Bytebase still carries `github.com/pingcap/tidb/pkg/parser` behind two
`replace` directives, and its TiDB adapter parses with omni first and falls
back to the pingcap parser per statement. Deleting that fallback needs two
different answers: the production fallback counter says whether real traffic
still needs it, and this inventory says whether the grammar is finished.

## Executive read

Measured on 2026-11-03 against TiDB v8.5.5, over every single-statement case
in pingcap's own parser test — 3,902 statements:

| Outcome | Cases |
|---|---:|
| both accept | 2090 |
| TiDB accepts, omni rejects (gap) | 1099 |
| TiDB rejects, omni accepts (over-acceptance) | 124 |
| both reject | 589 |
| no verdict | 0 |

- omni parses **2,090 of the 3,189 statements TiDB v8.5.5 accepts: 65.5%**.
- omni accepts **124 statements TiDB v8.5.5 rejects**, 5.6% of what it parses.
  An earlier pass put a floor of 32 on this; the floor was the measurement,
  not the engine.
- Nothing went unanswered. A statement the server never answered for is
  counted in its own column and fails the test, so the agreement numbers
  cannot absorb a dropped connection.

One thing the percentage is not: a parser test suite is written to press on
exotic grammar, so 65.5% is a completeness gate, not a model of customer
traffic. The decision to delete Bytebase's fallback rides on the production
counter; this number answers "is it finished".

## How it is measured

```bash
go test ./tidb/parser/ -run '^TestTiDBCorpusOmniVerdicts$' -count=1 -v        # no Docker
go test ./tidb/parser/ -run '^TestTiDBCorpusEngineVerdicts$' -count=1 -v      # TiDB v8.5.5
go test ./tidb/parser/ -run '^TestTiDBCorpusHarnessSurvivesALostServer$' -count=1 -v
```

The corpus is `tidb/parser/testdata/pingcap_parser_corpus.tsv`: one row per
statement holding the verdict each side gave it. `corpus_oracle_test.go`
reads it. `TestTiDBCorpusOmniVerdicts` holds `Parse` to the omni column and
needs no container, so it is the gate a grammar PR re-baselines with
`-update`. `TestTiDBCorpusEngineVerdicts` measures both sides live and holds
the engine column to what the pinned server says, so bumping the image shows
up as drift to read rather than as a silent change to these numbers. See
`tidb/parser/testdata/README.md` for provenance and for why the engine
verdict comes from `PREPARE` rather than from executing the statement.

## Target version

**TiDB v8.5.5**, the release `tidb/parser` already pins for its container
tests (`tidbImage`). `tidbCorpusVersion` names it beside the corpus and the
engine test refuses to run against any other release.

The 713 statements TiDB v8.5.5 rejects are recorded as expected rejections,
not as gaps: 589 omni also rejects, and 124 it should. 109 of them are
statements pingcap's own parser accepts at the extracted commit and this
release does not: grammar added after v8.5.5, or never released, and either
way not work for omni until the pin moves. A later bump re-records that column
and the diff is the list of what the new release added.

## How to work a row off this list

1. Pick a family below and read its rows in the corpus.
2. Add the grammar, with the usual engine-contract rules: every parse
   function returns an error, every node carries `Loc`, recursion goes
   through `enterNesting`, and a grammar decision cites upstream.
3. Check whether `mysql/parser` and `mariadb/parser` share the gap. A fix in
   a parent usually applies to its forks, in the same PR.
4. Rerun `TestTiDBCorpusOmniVerdicts -update` and read the diff: it should
   flip exactly the rows the change was for, from `reject` to `accept`.
5. Rerun `TestTiDBCorpusEngineVerdicts` so the live engine still agrees with
   the record.

## Gaps by stage

| Stage | Gaps |
|---|---:|
| Stage 1 — DDL the catalog walks | 336 |
| Stage 2 — the query and DML surface | 264 |
| Stage 3 — operational statements | 499 |
| total | 1099 |

## Stage 1 — DDL the catalog walks (336 gaps)

### ALTER TABLE clauses (163)

- **ALTER TABLE** — 10
  - `ALTER TABLE t ATTRIBUTES='str'`
  - `ALTER TABLE t ATTRIBUTES='str1,str2'`
- **ALTER TABLE ADD COLUMN INT** — 9
  - `alter table t add column a int partition by hash(a)`
  - `alter table t add column a int partition by hash(a) update indexes (idx_a global)`
- **ALTER TABLE DEFAULT** — 6
  - `ALTER TABLE t ATTRIBUTES=DEFAULT`
  - `ALTER TABLE t ATTRIBUTES=default`
- **ALTER TABLE ADD CONSTRAINT FOREIGN** — 5
  - `ALTER TABLE t ADD CONSTRAINT fk_t2_id FOREIGN KEY IF NOT EXISTS (t2_id) REFERENCES t(id)`
  - `alter table d_n.t_n add constraint foreign key ident (ident(1)) references d_n.t_n match full on delete…`
- **ALTER TABLE MODIFY COLUMN LONG** — 5
  - `alter table d_n.t_n modify column ident long after ident remove partitioning`
  - `alter table d_n.t_n modify column ident long char varying after ident remove partitioning`
- **ALTER TABLE OPTIMIZE PARTITION LOCAL** — 5
  - `ALTER TABLE t1 OPTIMIZE PARTITION LOCAL p0`
  - `ALTER TABLE t1 OPTIMIZE PARTITION LOCAL p0, p1`
- **ALTER TABLE PARTITION** — 5
  - `ALTER TABLE t PARTITION p ATTRIBUTES='str'`
  - `ALTER TABLE t PARTITION p ATTRIBUTES='str1,str2'`
- **ALTER TABLE REPAIR PARTITION LOCAL** — 5
  - `ALTER TABLE t1 REPAIR PARTITION LOCAL p0`
  - `ALTER TABLE t1 REPAIR PARTITION LOCAL p0, p1`
- **ALTER TABLE ADD INT CONSTRAINT** — 4
  - `ALTER TABLE t ADD (a1 int, CONSTRAINT PRIMARY KEY (a1))`
  - `ALTER TABLE t ADD a1 int CONSTRAINT ident CHECK ( a1 > 1 ) REFERENCES b ON DELETE CASCADE ON UPDATE CAS…`
- **ALTER TABLE REBUILD PARTITION LOCAL** — 4
  - `ALTER TABLE d_n.t_n REBUILD PARTITION LOCAL ALL`
  - `ALTER TABLE t_n REBUILD PARTITION LOCAL ident`
- **ALTER TABLE ADD COLUMN CONSTRAINT** — 3
  - `ALTER TABLE t ADD COLUMN ( CONSTRAINT ident CHECK ( 1>2 ) NOT ENFORCED )`
  - `ALTER TABLE t ADD COLUMN ( CONSTRAINT CHECK ( true ) )`
- **ALTER TABLE ADD COLUMN IF** — 3
  - `ALTER TABLE t ADD COLUMN IF NOT EXISTS (a SMALLINT UNSIGNED)`
  - `ALTER TABLE t ADD COLUMN IF NOT EXISTS (a SMALLINT UNSIGNED, b varchar(255))`
- **ALTER TABLE ADD PARTITION LOCAL** — 3
  - `ALTER TABLE employees ADD PARTITION LOCAL`
  - `ALTER TABLE employees ADD PARTITION LOCAL PARTITIONS 10`
- **ALTER TABLE ADD PARTITION NO_WRITE_TO_BINLOG** — 3
  - `ALTER TABLE employees ADD PARTITION NO_WRITE_TO_BINLOG`
  - `ALTER TABLE employees ADD PARTITION NO_WRITE_TO_BINLOG PARTITIONS 10`
- **ALTER TABLE CONVERT TO CHARSET** — 3
  - `ALTER TABLE t CONVERT TO CHARSET UTF8;`
  - `ALTER TABLE t CONVERT TO CHARSET UTF8 COLLATE UTF8_BIN;`
- **ALTER TABLE OPTIMIZE PARTITION NO_WRITE_TO_BINLOG** — 3
  - `ALTER TABLE t1 OPTIMIZE PARTITION NO_WRITE_TO_BINLOG p0`
  - `ALTER TABLE t1 OPTIMIZE PARTITION NO_WRITE_TO_BINLOG p0, p1`
- 87 further statements in 64 smaller families; read them in the corpus

### CREATE TABLE clauses (96)

- **CREATE TABLE INT PARTITION BY** — 16
  - `CREATE TABLE t (a int) STATS_BUCKETS=1,STATS_TOPN=1 PARTITION BY RANGE (a) (PARTITION p1 VALUES LESS TH…`
  - `create table t1 (a int) partition by system_time (partition x, partition y)`
- **CREATE TABLE INT** — 5
  - `` create table `123` (123a1 int) ``
  - `CREATE TABLE t (a int) STATS_BUCKETS=1`
- **CREATE TABLE INT PRIMARY KEY** — 5
  - `create table t (a int primary key local)`
  - `create table t (a int primary key global)`
- **CREATE GLOBAL TEMPORARY TABLE INT** — 3
  - `create global temporary table t (a int, b varchar(255)) on commit delete rows`
  - `create global temporary table t (a int, b varchar(255)) partition by hash(a) partitions 10 on commit de…`
- **CREATE TABLE LONG CHAR VARYING** — 3
  - `create table t (a long char varying);`
  - `create table t (a long char varying, b long varbinary);`
- **CREATE TABLE CHARACTER VARYING** — 2
  - `create table t (a character varying(1));`
  - `create table t (a character varying(255));`
- **CREATE TABLE ENUM SET** — 2
  - `create table t (c1 enum(0x61, 'b'), c2 set(0x61, 'b'))`
  - `create table t (c1 enum(0b01100001, 'b'), c2 set(0b01100001, 'b'))`
- **CREATE TABLE IF NOT EXISTS** — 2
  - `CREATE TABLE IF NOT EXISTS table_ident (a SQL_TSI_YEAR(4), b SQL_TSI_YEAR);`
  - `CREATE TABLE IF NOT EXISTS table_ident (ident1 BOOL COMMENT "text_string" unique, ident2 SQL_TSI_YEAR(4…`
- **CREATE TABLE INT SERIAL DEFAULT** — 2
  - `create table t (a int serial default value);`
  - `create table t (a int serial default value null);`
- **CREATE TABLE INT TABLE_CHECKSUM** — 2
  - `create table t (c int) table_checksum = 0`
  - `create table t (c int) table_checksum 1`
- **CREATE TABLE INT UNIQUE KEY** — 2
  - `create table t (a int UNIQUE key local)`
  - `create table t (a int UNIQUE key global)`
- **CREATE TABLE INT VARCHAR PRIMARY** — 2
  - `create table t (a int, b varchar(255), primary key(b, a) using RTREE clustered nonclustered)`
  - `create table t (a int, b varchar(255), primary key(b, a) using RTREE nonclustered clustered)`
- **CREATE TABLE LONG CHARACTER VARYING** — 2
  - `create table t (a long character varying);`
  - `create table t (a long character varying character set utf8);`
- **CREATE TABLE VARCHARACTER** — 2
  - `create table t (a varcharacter(1));`
  - `create table t (a varcharacter(50));`
- **CREATE TABLE AS VALUES ROW** — 1
  - `CREATE TABLE ta AS VALUES ROW(1)`
- **CREATE TABLE BIGINT SERIAL DEFAULT** — 1
  - `create table t (a bigint serial default value);`
- 44 further statements in 44 smaller families; read them in the corpus

### other DDL: database, view, index, TTL, misc (40)

- **ALTER SCHEMA DEFAULT COLLATE** — 3
  - `ALTER SCHEMA t DEFAULT COLLATE = 'UTF8_BiN'`
  - `ALTER SCHEMA DEFAULT COLLATE = 'UTF8_BiN'`
- **CREATE UNIQUE INDEX ON ASC** — 3
  - `CREATE UNIQUE INDEX ident ON d_n.t_n ( ident , ident ASC ) TYPE BTREE`
  - `CREATE UNIQUE INDEX ident ON d_n.t_n ( ident , ident ASC ) TYPE HASH`
- **DROP INDEX IF EXISTS ON** — 3
  - `drop index if exists a on t`
  - `drop index if exists a on db.t`
- **DROP INDEX ON ALGORITHM DEFAULT** — 3
  - `drop index idx on t algorithm = default`
  - `drop index idx on t algorithm default`
- **DROP INDEX ON LOCK DEFAULT** — 3
  - `drop index idx on t lock = default`
  - `drop index idx on t lock default`
- **ALTER DATABASE CHARACTER SET** — 2
  - `ALTER DATABASE t CHARACTER SET = 'utf8'`
  - `ALTER DATABASE CHARACTER SET = 'utf8'`
- **ALTER DATABASE COLLATE** — 2
  - `ALTER DATABASE t COLLATE = 'utf8_bin'`
  - `ALTER DATABASE COLLATE = 'utf8_bin'`
- **ALTER DATABASE DEFAULT CHARSET DEFAULT** — 2
  - `ALTER DATABASE t DEFAULT CHARSET = 'utf8mb4' DEFAULT COLLATE = 'utf8mb4_general_ci' CHARACTER SET = 'ut…`
  - `ALTER DATABASE DEFAULT CHARSET = 'utf8mb4' DEFAULT COLLATE = 'utf8mb4_general_ci' CHARACTER SET = 'utf8…`
- **ALTER SCHEMA DEFAULT CHARACTER SET** — 2
  - `ALTER SCHEMA t DEFAULT CHARACTER SET = 'utf8'`
  - `ALTER SCHEMA DEFAULT CHARACTER SET = 'utf8'`
- **CREATE VIEW AS TABLE** — 2
  - `CREATE VIEW v AS TABLE t`
  - `CREATE VIEW v AS (TABLE t)`
- **DROP TABLES** — 2
  - `drop tables xxx`
  - `drop tables xxx, yyy`
- **ALTER DATABASE CHARSET BINARY COLLATE** — 1
  - `ALTER DATABASE t CHARSET=binary COLLATE = binary`
- **ALTER DATABASE CHARSET COLLATE** — 1
  - `ALTER DATABASE t CHARSET = 'utf8mb4' COLLATE = 'utf8_bin'`
- **ALTER DATABASE COLLATE BINARY** — 1
  - `ALTER DATABASE t COLLATE = binary`
- **ALTER DATABASE DEFAULT CHARACTER SET** — 1
  - `ALTER DATABASE t DEFAULT CHARACTER SET = 'utf8'`
- **ALTER DATABASE DEFAULT CHARSET COLLATE** — 1
  - `ALTER DATABASE DEFAULT CHARSET = 'utf8mb4' COLLATE = 'utf8_bin'`
- 8 further statements in 8 smaller families; read them in the corpus

### vector and columnar index (29)

- **CREATE TABLE INT VECTOR VECTOR** — 9
  - `create table t(a int, b vector(3), vector index(b) USING HNSW);`
  - `create table t(a int, b vector(3), vector index(a, b) USING HNSW);`
- **ALTER TABLE ADD VECTOR INDEX** — 7
  - `ALTER TABLE t ADD VECTOR INDEX ((VEC_COSINE_DISTANCE(a, b))) USING HNSW COMMENT 'a'`
  - `ALTER TABLE t ADD VECTOR INDEX ((lower(a))) USING HNSW COMMENT 'a'`
- **CREATE VECTOR INDEX ON USING** — 5
  - `CREATE VECTOR INDEX idx ON t (a) USING HNSW`
  - `CREATE VECTOR INDEX idx ON t (a, b) USING HNSW`
- **CREATE VECTOR INDEX IF NOT** — 2
  - `CREATE VECTOR INDEX IF NOT EXISTS idx ON t ((VEC_COSINE_DISTANCE(a))) USING HNSW`
  - `CREATE VECTOR INDEX IF NOT EXISTS idx ON t ((VEC_COSINE_DISTANCE(a))) TYPE HNSW`
- **CREATE TABLE VECTOR FLOAT** — 1
  - `CREATE TABLE t (a VECTOR<FLOAT>)`
- **CREATE VECTOR INDEX ON** — 1
  - `CREATE VECTOR INDEX idx ON t ((VEC_COSINE_DISTANCE(a)))`
- **CREATE VECTOR INDEX ON ASC** — 1
  - `CREATE VECTOR INDEX ident ON d_n.t_n ( ident , ident ASC ) TYPE HNSW`
- **CREATE VECTOR INDEX ON TYPE** — 1
  - `CREATE VECTOR INDEX idx ON t ((VEC_COSINE_DISTANCE(a))) TYPE BTREE`
- **CREATE VECTOR INDEX TYPE ON** — 1
  - `CREATE VECTOR INDEX ident TYPE HNSW ON d_n.t_n ((VEC_COSINE_DISTANCE(a)))`
- **CREATE VECTOR INDEX USING ON** — 1
  - `CREATE VECTOR INDEX idx USING HNSW ON t ((VEC_COSINE_DISTANCE(a)))`

### placement policy (5)

- **ALTER RANGE GLOBAL PLACEMENT POLICY** — 2
  - `ALTER RANGE global PLACEMENT POLICY mypolicy`
  - `ALTER RANGE global PLACEMENT POLICY default`
- **ALTER RANGE PLACEMENT POLICY** — 1
  - `ALTER RANGE meta PLACEMENT POLICY mypolicy`
- **ALTER TABLE PARTITION PLACEMENT PLACEMENT** — 1
  - `alter table m partition t /*T![placement] placement policy="ww" */;`
- **ALTER TABLE PARTITION PLACEMENT POLICY** — 1
  - `alter table m partition t placement policy="ww";`

### sequences (3)

- **SELECT NEXT VALUE FOR** — 2
  - `select next value for seq`
  - `select NeXt vAluE for seQuEncE2`
- **SELECT NEXT VALUE FOR SEQUENCE** — 1
  - `select next value for sequence`

## Stage 2 — the query and DML surface (264 gaps)

### queries and DML (185)

- **SELECT** — 11
  - `select _utf8 0xD0B1;`
  - `select _utf8 X'D0B1';`
- **SELECT FROM TABLESAMPLE** — 4
  - `select * from tbl tablesample (50);`
  - `select * from tbl tablesample (0 percent);`
- **INSERT INTO VALUES** — 3
  - `INSERT INTO foo VALUES (1 || 2)`
  - `INSERT INTO foo VALUES ()`
- **REPLACE INTO VALUES** — 3
  - `REPLACE INTO foo VALUES (1 || 2)`
  - `REPLACE INTO foo VALUES ()`
- **SELECT FROM AS OF TIMESTAMP** — 3
  - `` SELECT * FROM `t` AS OF TIMESTAMP TIDB_BOUNDED_STALENESS(DATE_SUB(NOW(), INTERVAL 3 SECOND), NOW()); ``
  - `` select * from `t` as of timestamp '2021-04-15 00:00:00' ``
- **SELECT FROM TABLESAMPLE SYSTEM** — 3
  - `select * from tbl tablesample system (50);`
  - `select * from tbl tablesample system (50 percent);`
- **SELECT FROM WHERE** — 3
  - `select * from t where a > {ts '1989-09-10 11:11:11'}`
  - `select * from t where a > {ts {abc '1989-09-10 11:11:11'}}`
- **SELECT CAST DATA AS CHARACTER** — 2
  - `SELECT CAST(data AS CHARACTER);`
  - `SELECT CAST(data AS CHARACTER(10) CHARACTER SET utf8);`
- **SELECT FROM** — 2
  - `` select `t`.`1a`.1 from t; ``
  - `select SQL_CACHE * from t`
- **SELECT FROM FETCH FIRST ROW** — 2
  - `SELECT * FROM t FETCH FIRST 5 ROW ONLY`
  - `SELECT * FROM t FETCH FIRST ROW ONLY`
- **SELECT FROM FETCH NEXT ROW** — 2
  - `SELECT * FROM t FETCH NEXT 5 ROW ONLY`
  - `SELECT * FROM t FETCH NEXT ROW ONLY`
- **SELECT FROM LIMIT FOR UPDATE** — 2
  - `select * from t limit 1 for update wait 11`
  - `select * from t limit 1 for update of t wait 11`
- **SELECT FROM TABLESAMPLE TABLESAMPLE** — 2
  - `` select * from tbl `tablesample` tablesample (50); ``
  - `select * from tbl1 tablesample (10), tbl2 tablesample (20);`
- **SELECT INTERVAL INTERVAL** — 2
  - `SELECT (INTERVAL(1, 0, 1, 2)+5)*7+INTERVAL(1, 0, 1, 2)/2`
  - `SELECT INTERVAL(0, (1*5)/2)+INTERVAL(5, 4, 3)`
- **SELECT JSON_OBJECTAGG ALL FROM GROUP** — 2
  - `select json_objectagg(all c1, c2) from t group by c1`
  - `select json_objectagg(c1, all c2) from t group by c1`
- **TABLE EXCEPT TABLE** — 2
  - `table t1 except table t2`
  - `table t1 except (table t2)`
- 137 further statements in 132 smaller families; read them in the corpus

### SHOW variants (60)

- **SHOW** — 14
  - `SHOW SESSION_STATES`
  - `show stats_extended`
- **SHOW WHERE TABLE_NAME** — 6
  - `show stats_extended where table_name = 't'`
  - `show stats_meta where table_name = 't'`
- **SHOW TABLE** — 3
  - `show table t1 distributions`
  - `show table t1.t1 next_row_id`
- **SHOW CREATE SEQUENCE** — 2
  - `show create sequence seq`
  - `show create sequence test.seq`
- **SHOW PLACEMENT FOR TABLE** — 2
  - `SHOW PLACEMENT FOR TABLE tb1`
  - `SHOW PLACEMENT FOR TABLE db1.tb1`
- **SHOW PLACEMENT FOR TABLE PARTITION** — 2
  - `SHOW PLACEMENT FOR TABLE tb1 PARTITION p1`
  - `SHOW PLACEMENT FOR TABLE db1.tb1 PARTITION p1`
- **SHOW TABLE PARTITION INDEX REGIONS** — 2
  - `show table t1 partition (p0) index idx1 regions`
  - `show table t1 partition (p0,p1) index idx1 regions where a=2`
- **SHOW WHERE** — 2
  - `show stats_histograms where col_name = 'a'`
  - `show stats_buckets where col_name = 'a'`
- **SHOW BINARY LOG STATUS** — 1
  - `SHOW BINARY LOG STATUS`
- **SHOW COLUMNS FROM** — 1
  - `SHOW COLUMNS FROM tv189.1_t_1_x;`
- **SHOW CREATE PLACEMENT POLICY** — 1
  - `show create placement policy x`
- **SHOW CREATE PLACEMENT POLICY PLACEMENT** — 1
  - `` show create placement policy `placement` ``
- **SHOW LIKE** — 1
  - `show restores like 'r0001'`
- **SHOW PLACEMENT** — 1
  - `SHOW PLACEMENT`
- **SHOW PLACEMENT FOR DATABASE** — 1
  - `SHOW PLACEMENT FOR DATABASE db1`
- **SHOW PLACEMENT FOR SCHEMA** — 1
  - `SHOW PLACEMENT FOR SCHEMA db1`
- 19 further statements in 19 smaller families; read them in the corpus

### EXPLAIN forms (19)

- **EXPLAIN FORMAT SELECT** — 6
  - `EXPLAIN FORMAT = 'dot' SELECT 1`
  - `EXPLAIN FORMAT = 'BRIEF' SELECT 1`
- **EXPLAIN ANALYZE FORMAT SELECT** — 3
  - `EXPLAIN ANALYZE format='VERBOSE' SELECT 1`
  - `EXPLAIN ANALYZE format='TRUE_CARD_COST' SELECT 1`
- **EXPLAIN FORMAT ROW SELECT** — 2
  - `EXPLAIN FORMAT = 'row' SELECT 1`
  - `EXPLAIN FORMAT = 'ROW' SELECT 1`
- **EXPLAIN ALTER TABLE ADD INDEX** — 1
  - `EXPLAIN ALTER TABLE t1 ADD INDEX (a)`
- **EXPLAIN ALTER TABLE ADD VARCHAR** — 1
  - `EXPLAIN ALTER TABLE t1 ADD a varchar(255)`
- **EXPLAIN ANALYZE FORMAT BINARY SELECT** — 1
  - `EXPLAIN ANALYZE FORMAT = 'binary' SELECT 1`
- **EXPLAIN FORMAT FOR CONNECTION** — 1
  - `EXPLAIN FORMAT = 'dot' FOR CONNECTION 1`
- **EXPLAIN FORMAT ROW FOR CONNECTION** — 1
  - `EXPLAIN FORMAT = 'row' FOR connection 1`
- **EXPLAIN FORMAT ROW SELECT FROM** — 1
  - `explain format = "row" select c1 from t1 union (select c2 from t2) limit 1, 1`
- **EXPLAIN FORMAT SELECT FROM** — 1
  - `explain format = 'brief' select * from t`
- **EXPLAIN REPLACE INTO VALUES** — 1
  - `explain replace into foo values (1 || 2)`

## Stage 3 — operational statements (499 gaps)

### statistics, analyze and SQL bindings (122)

- **ANALYZE TABLE WITH** — 6
  - `analyze table t with 4 topn`
  - `analyze table t with 4 cmsketch width`
- **ANALYZE TABLE INDEX** — 5
  - `analyze table t1 index`
  - `analyze table t1 index a`
- **CREATE ON** — 4
  - `create statistics stats1 (cardinality) on t(a,b,c)`
  - `create statistics stats2 (dependency) on t(a,b)`
- **DROP GLOBAL FOR DELETE FROM** — 4
  - `drop global binding for delete from t where a = 1`
  - `` DROP GLOBAL BINDING FOR DELETE FROM `t` WHERE `a`=1 USING DELETE /*+ USE_INDEX(`t` `a`)*/ FROM `t` WHER… ``
- **DROP GLOBAL FOR SELECT FROM** — 4
  - `drop global binding for select * from t`
  - `drop global binding for select * from t using select * from t use index(a)`
- **DROP SESSION FOR DELETE FROM** — 4
  - `drop session binding for delete from t where a = 1`
  - `` DROP SESSION BINDING FOR DELETE FROM `t` WHERE `a`=1 USING DELETE /*+ USE_INDEX(`t` `a`)*/ FROM `t` WHE… ``
- **CREATE IF NOT EXISTS ON** — 3
  - `create statistics if not exists stats1 (cardinality) on t(a,b,c)`
  - `create statistics if not exists stats2 (dependency) on t(a,b)`
- **DROP** — 3
  - `drop stats t`
  - `drop stats t1, t2, t3`
- **DROP SESSION FOR SELECT FROM** — 3
  - `drop session binding for select * from t`
  - `drop session binding for select * from t using select * from t use index(a)`
- **ADMIN** — 2
  - `admin capture bindings`
  - `admin evolve bindings`
- **ANALYZE TABLE COLUMNS** — 2
  - `analyze table t1,t2 predicate columns`
  - `analyze table t columns c1,c2`
- **ANALYZE TABLE COLUMNS WITH** — 2
  - `analyze table t1,t2 predicate columns with 4 topn`
  - `analyze table t columns c1,c2 with 4 topn`
- **ANALYZE TABLE PARTITION ALL COLUMNS** — 2
  - `analyze table t partition a all columns`
  - `analyze table t partition a all columns with 1024 buckets`
- **ANALYZE TABLE PARTITION COLUMNS** — 2
  - `analyze table t partition a predicate columns`
  - `analyze table t partition a columns c1,c2`
- **ANALYZE TABLE PARTITION COLUMNS WITH** — 2
  - `analyze table t partition a predicate columns with 1024 buckets`
  - `analyze table t partition a columns c1,c2 with 1024 buckets`
- **ANALYZE TABLE WITH BUCKETS** — 2
  - `analyze table t with 4 buckets`
  - `analyze table t with 4 buckets, 4 topn, 4 cmsketch width, 4 cmsketch depth, 4 samples`
- 72 further statements in 50 smaller families; read them in the corpus

### ADMIN, SPLIT and DISTRIBUTE (90)

- **ADMIN REPAIR TABLE CREATE TABLE** — 6
  - `ADMIN REPAIR TABLE t CREATE TABLE t (a int)`
  - `ADMIN REPAIR TABLE t CREATE TABLE t (a char(1))`
- **ADMIN SHOW** — 5
  - `admin show ddl;`
  - `admin show ddl jobs;`
- **ADMIN** — 4
  - `admin cancel ddl jobs 1`
  - `admin cancel ddl jobs 1, 2`
- **ADMIN ALTER** — 4
  - `admin alter ddl jobs 1 thread = 2`
  - `admin alter ddl jobs 1 batch_size = 3`
- **ADMIN RELOAD** — 3
  - `admin reload expr_pushdown_blacklist`
  - `admin reload statistics`
- **ADMIN SHOW SLOW** — 3
  - `admin show slow top 3`
  - `admin show slow top internal 7`
- **TABLE BETWEEN AND REGIONS** — 3
  - `split table t1 between ('a') and ('z') regions 10`
  - `split table t1 between ('a',1) and ('z',2) regions 10`
- **TABLE BY** — 3
  - `split table t1 by ('a'),('b'),('c')`
  - `split table t1 by (1)`
- **TABLE INDEX BETWEEN AND REGIONS** — 3
  - `split table t1 index idx1 between ('a') and ('z') regions 10`
  - `split table t1 index idx1 between ('a',1) and ('z',2) regions 10`
- **TABLE INDEX BY** — 3
  - `split table t1 index idx1 by ('a'),('b'),('c')`
  - `split table t1 index idx1 by (1)`
- **ADMIN CHECK INDEX** — 2
  - `admin check index tableName idxName;`
  - `admin check index tableName idxName (1, 2), (4, 5);`
- **ADMIN RESUME** — 2
  - `admin resume ddl jobs 1, 2`
  - `admin resume ddl jobs 3`
- **ADMIN SHOW LIMIT** — 2
  - `admin show ddl job queries limit 5`
  - `admin show ddl job queries limit 5, 10`
- **ADMIN SHOW LIMIT OFFSET** — 2
  - `admin show ddl job queries limit 3 offset 2`
  - `admin show ddl job queries limit 22 offset 0`
- **ADMIN SHOW WHERE** — 2
  - `admin show ddl jobs where id > 0;`
  - `admin show ddl jobs 20 where id=0;`
- **ADMIN TABLE LOCK** — 2
  - `ADMIN CLEANUP TABLE LOCK t`
  - `ADMIN CLEANUP TABLE LOCK t1,t2`
- 41 further statements in 37 smaller families; read them in the corpus

### resource groups and query watch (82)

- **CREATE RESOURCE GROUP** — 19
  - `create resource group x ru_per_sec=2000`
  - `create resource group x ru_per_sec=200000`
- **ALTER RESOURCE GROUP** — 13
  - `alter resource group x ru_per_sec=1000`
  - `alter resource group x ru_per_sec=2000, BURSTABLE`
- **CREATE RESOURCE GROUP ACTION** — 12
  - `create resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '10s' ACTION DRYRUN)`
  - `create resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '10m' ACTION COOLDOWN)`
- **ALTER RESOURCE GROUP ACTION** — 7
  - `alter resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '10s' RU 100 ACTION DRYRUN)`
  - `alter resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '10s' PROCESSED_KEYS 100 ACTION DRYRUN)`
- **QUERY ADD SQL TEXT TO** — 4
  - `query watch add SQL TEXT SIMILAR to 'select 1'`
  - `query watch add SQL TEXT EXACT to 'select 1'`
- **ALTER RESOURCE GROUP NULL** — 2
  - `alter resource group x QUERY_LIMIT=NULL`
  - `alter resource group x background NULL`
- **CREATE RESOURCE GROUP DEFAULT** — 2
  - `create resource group default ru_per_sec=20, priority=LOW, burstable`
  - `create resource group default ru_per_sec=UNLIMITED, priority=LOW, burstable`
- **QUERY ADD** — 2
  - `` query watch add PLAN DIGEST `5e3ddd388f6012e328233dbcdda5d48f404e0536c6c54d9618233210f3d5762a` ``
  - `query watch add PLAN DIGEST @digest1`
- **QUERY ADD RESOURCE GROUP SQL** — 2
  - `query watch add resource group @rg SQL TEXT SIMILAR to @sql1`
  - `query watch add resource group rg1 SQL TEXT SIMILAR to 'select 1'`
- **QUERY ADD SQL** — 2
  - `query watch add SQL DIGEST b13858789fce00208f9a262c99621b7045f4869807cd4e6568008ae7ca19a377`
  - `query watch add SQL DIGEST 'b13858789fce00208f9a262c99621b7045f4869807cd4e6568008ae7ca19a377'`
- **ALTER RESOURCE GROUP ACTION KILL** — 1
  - `alter resource group x ru_per_sec=1000 QUERY_LIMIT=( ACTION KILL EXEC_ELAPSED '10m')`
- **ALTER RESOURCE GROUP DEFAULT** — 1
  - `alter resource group default priority = high`
- **ALTER RESOURCE GROUP DEFAULT TTL** — 1
  - `alter resource group default priority=low background = ( task_types "ttl" )`
- **ALTER RESOURCE GROUP FALSE** — 1
  - `alter resource group x burstable=false`
- **ALTER RESOURCE GROUP TRUE** — 1
  - `alter resource group x burstable=true`
- **ALTER USER RESOURCE GROUP** — 1
  - `ALTER USER 'root'@'localhost' RESOURCE GROUP rg2`
- 11 further statements in 11 smaller families; read them in the corpus

### diagnostics: TRACE, PLAN REPLAYER, RECOMMEND INDEX, TRAFFIC (64)

- **RESOURCE** — 13
  - `calibrate resource`
  - `calibrate resource START_TIME '2021-04-15 00:00:00'`
- **SELECT FROM** — 4
  - `trace select c1 from t1`
  - `trace plan select c1 from t1`
- **INDEX RUN FOR SELECT FROM** — 3
  - `recommend index run for 'select * from t where a=1'`
  - `recommend index run for 'select * from t where a=1' with A = 1`
- **INDEX SET** — 3
  - `recommend index set A = 1`
  - `recommend index set A = 1, B = 2`
- **EXPLAIN ANALYZE SLOW QUERY WHERE** — 2
  - `PLAN REPLAYER DUMP EXPLAIN ANALYZE SLOW QUERY WHERE a > 10 and t < 1 ORDER BY t LIMIT 10`
  - `PLAN REPLAYER DUMP EXPLAIN ANALYZE SLOW QUERY WHERE a > 10 and t < 1 LIMIT 10`
- **EXPLAIN SLOW QUERY WHERE AND** — 2
  - `PLAN REPLAYER DUMP EXPLAIN SLOW QUERY WHERE a > 10 and t < 1 ORDER BY t LIMIT 10`
  - `PLAN REPLAYER DUMP EXPLAIN SLOW QUERY WHERE a > 10 and t < 1 LIMIT 10`
- **INDEX RUN WITH** — 2
  - `recommend index run with A = 1`
  - `recommend index run with A = 1, B = 2`
- **RESOURCE NOW** — 2
  - `calibrate resource END_TIME now()`
  - `calibrate resource START_TIME now()`
- **BEGIN** — 1
  - `trace begin`
- **COMMIT** — 1
  - `trace commit`
- **DELETE FROM INNER JOIN INNER** — 1
  - `trace delete t1, t2 from t1 inner join t2 inner join t3 where t1.id=t2.id and t2.id=t3.id;`
- **EXPLAIN ANALYZE SELECT FROM WHERE** — 1
  - `PLAN REPLAYER DUMP EXPLAIN ANALYZE SELECT * FROM t WHERE a > 10`
- **EXPLAIN ANALYZE SLOW QUERY** — 1
  - `PLAN REPLAYER DUMP EXPLAIN ANALYZE SLOW QUERY`
- **EXPLAIN ANALYZE SLOW QUERY LIMIT** — 1
  - `PLAN REPLAYER DUMP EXPLAIN ANALYZE SLOW QUERY LIMIT 10`
- **EXPLAIN ANALYZE SQL** — 1
  - `PLAN REPLAYER DUMP EXPLAIN ANALYZE 'sql.txt'`
- **EXPLAIN SELECT FROM** — 1
  - `PLAN REPLAYER DUMP EXPLAIN SELECT a FROM t`
- 25 further statements in 25 smaller families; read them in the corpus

### BACKUP, RESTORE, IMPORT and LOAD (59)

- **LOAD DATA LOCAL INFILE INTO** — 8
  - `` load data local infile '~/1.csv' into table `t_ascii` fields terminated by X'6B6B'; ``
  - `` load data local infile '~/1.csv' into table `t_ascii` fields terminated by X'6B6B' enclosed by X'0D'; ``
- **LOAD DATA INFILE INTO TABLE** — 6
  - `load data infile '/tmp/t.csv' into table t fields defined null by 'nil'`
  - `load data infile '/tmp/t.csv' into table t fields defined null by X'00'`
- **BACKUP DATABASE TO** — 4
  - `BACKUP DATABASE a,b,c TO 'noop://'`
  - `BACKUP DATABASE * TO 'noop://'`
- **RESTORE POINT FROM** — 4
  - `restore point from 'noop://log_backup'`
  - `restore point from 'noop://log_backup' full_backup_storage='noop://full_log'`
- **BACKUP TABLE TO** — 3
  - `BACKUP TABLE a TO 'noop://' checksum_concurrency 4 compression_level 4 ignore_stats 1 compression_type …`
  - `BACKUP TABLE a.b TO 'noop://'`
- **IMPORT INTO FROM FILE** — 3
  - `import into t from '/file.csv'`
  - `import into t (a,b) from '/file.csv'`
- **IMPORT INTO FROM FILE WITH** — 3
  - `import into t from '/file.csv' with detached`
  - `` import into `t` from '/file.csv' with thread=1 ``
- **BACKUP DATABASE TO SNAPSHOT** — 2
  - `backup database * to 'noop://' snapshot = '2020-03-18 18:13:54'`
  - `backup database * to 'noop://' snapshot = 1234567890`
- **BACKUP LOGS** — 2
  - `pause backup logs`
  - `pause backup logs gc_ttl='20220304'`
- **BACKUP LOGS TO** — 2
  - `backup logs to 'noop://'`
  - `backup logs to 'noop://' start_ts='20220304'`
- **PURGE BACKUP LOGS FROM** — 2
  - `purge backup logs from 'noop://'`
  - `purge backup logs from 'noop://' until_ts='2012122304'`
- **BACKUP DATABASE TO LOCAL** — 1
  - `BACKUP DATABASE a TO 'local:///tmp/archive01/'`
- **BACKUP DATABASE TO S3 PATH** — 1
  - `backup database * to 's3://bucket/path/?provider=alibaba&region=us-west-9&storage-class=glacier&sse=AES…`
- **BACKUP DATABASE TO SECOND SNAPSHOT** — 1
  - `backup database * to 'noop://' rate_limit 500 MB/second snapshot 5 minute ago`
- **BACKUP SCHEMA TO LOCAL** — 1
  - `BACKUP SCHEMA a TO 'local:///tmp/archive01/'`
- **BACKUP TABLE TO S3 PATH** — 1
  - `backup table x to 's3://bucket/path/?endpoint=https://test-cluster-s3.local&access-key=aaaaaaaaa&secret…`
- 15 further statements in 15 smaller families; read them in the corpus

### users, grants and roles (30)

- **CREATE USER** — 6
  - `CREATE USER uesr1@LOCALhost`
  - `` CREATE USER `uesr1`@localhost ``
- **GRANT TO** — 4
  - `GRANT 'role1', 'role2' TO 'user1'@'LOCalhost', 'user2'@'LOcalhost';`
  - `GRANT 'u1' TO 'u1';`
- **CREATE USER IDENTIFIED WITH AS** — 3
  - `CREATE USER 'u1'@'%' IDENTIFIED WITH 'mysql_native_password' AS '' REQUIRE NONE PASSWORD EXPIRE DEFAULT…`
  - `CREATE USER 'sha_test3'@'localhost' IDENTIFIED WITH 'caching_sha2_password' AS 0x24412430303524255B0349…`
- **CREATE USER IDENTIFIED WITH** — 2
  - `CREATE USER 'nopwd_native'@'localhost' IDENTIFIED WITH 'mysql_native_password'`
  - `CREATE USER 'nopwd_sha'@'localhost' IDENTIFIED WITH 'caching_sha2_password'`
- **ALTER USER IDENTIFIED BY NEW** — 1
  - `ALTER USER 'root'@'localhost' IDENTIFIED BY 'new-password', 'root'@'127.0.0.1' IDENTIFIED BY PASSWORD '…`
- **ALTER USER IDENTIFIED BY PASSWORD** — 1
  - `ALTER USER 'root'@'localhost' IDENTIFIED BY PASSWORD 'hashstring'`
- **CREATE USER IDENTIFIED BY NEW** — 1
  - `CREATE USER 'root'@'localhost' IDENTIFIED BY 'new-password', 'root'@'127.0.0.1' IDENTIFIED BY PASSWORD …`
- **CREATE USER IDENTIFIED BY PASSWORD** — 1
  - `CREATE USER 'root'@'localhost' IDENTIFIED BY PASSWORD 'hashstring'`
- **CREATE USER IDENTIFIED WITH BY** — 1
  - `CREATE USER 'sha_test'@'localhost' IDENTIFIED WITH 'caching_sha2_password' BY 'sha_test'`
- **CREATE USER REQUIRE USER URI** — 1
  - `CREATE USER 'ttt' REQUIRE SAN 'DNS:mysql-user, URI:spiffe://example.org/myservice'`
- **CREATE USER USER IDENTIFIED WITH** — 1
  - `` CREATE USER `user@pingcap.com`@'localhost' IDENTIFIED WITH 'tidb_auth_token' REQUIRE token_issuer 'issu… ``
- **GRANT ALL PRIVILEGES ON TO** — 1
  - `grant all privileges on zabbix.* to 'zabbix'@'localhost' identified by 'password';`
- **GRANT GRANT OPTION ON TO** — 1
  - `grant grant option on *.* to u1`
- **GRANT LOAD FROM S3 SELECT** — 1
  - `GRANT LOAD FROM S3, SELECT INTO S3, INVOKE LAMBDA, INVOKE SAGEMAKER, INVOKE COMPREHEND ON *.* TO 'root'…`
- **GRANT PROCESS USAGE REPLICATION SLAVE** — 1
  - `grant PROCESS,usage, REPLICATION SLAVE, REPLICATION CLIENT on *.* to 'xxxxxxxxxx'@'%' identified by pas…`
- **GRANT SELECT INSERT UPDATE DELETE** — 1
  - `/* rds internal mark */ GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, REFERENCES, RELOAD, PROCESS…`
- 3 further statements in 3 smaller families; read them in the corpus

### SET and session state (31)

- **SET LOG LEVEL** — 4
  - `set config TIKV LOG.LEVEL='info'`
  - `set config PD LOG.LEVEL='info'`
- **START TRANSACTION READ ONLY AS** — 4
  - `START TRANSACTION READ ONLY AS OF TIMESTAMP '2015-09-21 00:07:01'`
  - `START TRANSACTION READ ONLY AS OF TIMESTAMP TIDB_BOUNDED_STALENESS(_UTF8MB4'2015-09-21 00:07:01', NOW())`
- **SET** — 3
  - `SET SESSION_STATES 'x'`
  - `SET @ = 1`
- **BEGIN** — 2
  - `begin pessimistic`
  - `begin optimistic`
- **FLUSH PLUGINS** — 2
  - `flush tidb plugins plugin1`
  - `flush tidb plugins plugin1, plugin2`
- **KILL** — 1
  - `kill tidb 23123`
- **KILL CONNECTION** — 1
  - `kill tidb connection 23123`
- **KILL QUERY** — 1
  - `kill tidb query 23123`
- **LOCK TABLE WRITE LOCAL** — 1
  - `LOCK TABLE t2 WRITE LOCAL;`
- **LOCK TABLES WRITE LOCAL** — 1
  - `LOCK TABLES t2 WRITE LOCAL;`
- **LOCK TABLES WRITE LOCAL READ** — 1
  - `LOCK TABLES t1 WRITE LOCAL, t2 READ LOCAL;`
- **SET AUTO MODE TRUE** — 1
  - `set config '127.0.0.1:3306' AUTO-COMPACTION-MODE=TRUE`
- **SET BINARY** — 1
  - `SET @@character_set_results = binary`
- **SET CHAR SET DEFAULT** — 1
  - `set char set default`
- **SET CHARACTER SET** — 1
  - `SET CHARACTER SET 'utf8mb4';`
- **SET KEY ZONE** — 1
  - `set config '127.0.0.1:3306' LABEL-PROPERTY.REJECT-LEADER.KEY='zone'`
- 5 further statements in 5 smaller families; read them in the corpus

### FLASHBACK and RECOVER (21)

- **TABLE TO** — 3
  - `flashback table t TO t1`
  - `flashback table t to tso 445494955052105722`
- **TABLE TO TIMESTAMP** — 3
  - `flashback table t TO timestamp`
  - `flashback table t to timestamp '2021-05-26 16:45:26'`
- **DATABASE TO** — 2
  - `flashback database db1 to db2`
  - `flashback database test to tso 445494955052105724`
- **RECOVER TABLE** — 2
  - `recover table t1`
  - `recover table t1 100`
- **RECOVER TABLE BY** — 2
  - `recover table by job 11`
  - `recover table by job 0`
- **SCHEMA TO** — 2
  - `flashback schema db1 to db2`
  - `flashback schema test to tso 445494955052105725`
- **DATABASE** — 1
  - `flashback database db1`
- **DATABASE TO TIMESTAMP** — 1
  - `flashback database test to timestamp '2021-05-26 16:45:26'`
- **SCHEMA** — 1
  - `flashback schema db1`
- **SCHEMA TO TIMESTAMP** — 1
  - `flashback schema test to timestamp '2021-05-26 16:45:26'`
- **TABLE** — 1
  - `flashback table t`
- **TO** — 1
  - `flashback cluster to tso 445494955052105721`
- **TO TIMESTAMP** — 1
  - `flashback cluster to timestamp '2021-05-26 16:45:26'`

## Over-acceptance (124)

Statements `omni/tidb/parser` accepts and TiDB v8.5.5 rejects. These are
not features to add but strictness to restore, and they matter more than
their count: a parser that accepts what the engine will not lets a review
pass a statement that then fails on the server.

### ALTER TABLE clauses (14)

- **ALTER TABLE ADD COLUMN** — 3
  - `alter table t add column a int, partition by range(a) (partition x values less than (75))`
  - `alter table t add column a int, remove partitioning`
  - `alter table t add column a int, add index (c), remove partitioning`
- **ALTER TABLE ALTER COLUMN** — 3
  - `ALTER TABLE t ALTER COLUMN a SET DEFAULT CURRENT_TIMESTAMP`
  - `ALTER TABLE t ALTER COLUMN a SET DEFAULT NOW()`
  - `ALTER TABLE t ALTER COLUMN a SET DEFAULT 1+1`
- **ALTER TABLE ALTER CONSTRAINT** — 1
  - `ALTER TABLE t_n ALTER CONSTRAINT ident`
- **ALTER TABLE ENABLE KEYS** — 1
  - `alter table t enable keys, comment = 'cmt', partition by hash(a)`
- **ALTER TABLE OPTIMIZE PARTITION** — 1
  - `ALTER TABLE t_n OPTIMIZE PARTITION LOCAL`
- **ALTER TABLE PARTITION BY** — 1
  - `alter table t partition by hash(a), enable keys`
- **ALTER TABLE REBUILD PARTITION** — 1
  - `ALTER TABLE t_n REBUILD PARTITION LOCAL`
- **ALTER TABLE REMOVE PARTITIONING** — 1
  - `alter table t remove partitioning, add column a int`
- **ALTER TABLE REPAIR PARTITION** — 1
  - `ALTER TABLE t_n REPAIR PARTITION LOCAL`
- **ALTER TABLE TTL TTL_ENABLE** — 1
  - `alter table t /*T![ttl] TTL_ENABLE = 'test_case' */`

### CREATE TABLE clauses (32)

- **CREATE TABLE INT PARTITION** — 8
  - `create table t (b int) partition by range columns (b) (partition p0 values less than (not 3), partition…`
  - `create table t (b int) partition by range columns (b) (partition p0 values less than (1 or 3), partitio…`
  - `create table t (b int) partition by range columns (b) (partition p0 values less than (3 is null), parti…`
- **CREATE TABLE INT SECONDARY_ENGINE_ATTRIBUTE** — 4
  - `CREATE TABLE t (id INT) SECONDARY_ENGINE_ATTRIBUTE='{"key":"value"}'`
  - `CREATE TABLE t (id INT SECONDARY_ENGINE_ATTRIBUTE='{"key":"value"}')`
  - `CREATE TABLE t (id INT SECONDARY_ENGINE_ATTRIBUTE=)`
- **CREATE TABLE DATETIME TTL_JOB_INTERVAL** — 3
  - `create table t (created_at datetime) TTL_JOB_INTERVAL = '@monthly'`
  - `create table t (created_at datetime) TTL_JOB_INTERVAL = '10hourxx'`
  - `create table t (created_at datetime) TTL_JOB_INTERVAL = '10.10.255h'`
- **CREATE TABLE INT DEFAULT** — 3
  - `create table t (i int default (0))`
  - `create table t (i int default (-1))`
  - `create table t (i int default (+1))`
- **CREATE TABLE INT INDEX** — 3
  - `CREATE TABLE t (id INT,INDEX idx (id) INVISIBLE SECONDARY_ENGINE_ATTRIBUTE='{"key1":"value1"}')`
  - `CREATE TABLE t (id INT, INDEX idx (id) SECONDARY_ENGINE_ATTRIBUTE=)`
  - `CREATE TABLE t (id INT, INDEX idx (id) SECONDARY_ENGINE_ATTRIBUTE)`
- **CREATE TABLE INT STATS_AUTO_RECALC** — 3
  - `create table t (a int) stats_auto_recalc 2;`
  - `create table t (a int) stats_auto_recalc = 10;`
  - `CREATE TABLE t (a int) STATS_AUTO_RECALC='abc'`
- **CREATE TABLE INT INT** — 2
  - `create table t (a int, b int, c char(33) default (b))`
  - `` create table t (a int, b int, c char(33) default `b`) ``
- **CREATE TABLE CHAR DEFAULT** — 1
  - `create table t (c char(33) default ('foobar'))`
- **CREATE TABLE DATETIME TTL** — 1
  - `create table t (created_at datetime) /*T![ttl] TTL_ENABLE = 'test_case' */`
- **CREATE TABLE DATETIME TTL_ENABLE** — 1
  - `create table t (created_at datetime) TTL_ENABLE = 'test_case'`
- **CREATE TABLE INT TABLESPACE** — 1
  - `CREATE TABLE t (id INT) TABLESPACE ts1 SECONDARY_ENGINE_ATTRIBUTE='{"key":"value"}'`
- **CREATE TABLE NATIONAL VARCHAR** — 1
  - `create table t (a national varchar);`
- 1 further statements in 1 smaller families; read them in the corpus

### other DDL: database, view, index, TTL, misc (2)

- **CREATE INDEX ON SECONDARY_ENGINE_ATTRIBUTE** — 2
  - `CREATE INDEX i ON t (a) SECONDARY_ENGINE_ATTRIBUTE = '{}'`
  - `CREATE INDEX i ON t (a) SECONDARY_ENGINE_ATTRIBUTE '{}'`

### queries and DML (70)

- **SELECT FROM** — 4
  - `SELECT '{}'->'$.a' FROM t`
  - `SELECT '{}'->>'$.a' FROM t`
  - `SELECT a->3 FROM t`
- **SELECT CURRENT_TIME** — 3
  - `select current_time(-1)`
  - `select current_time(1.0)`
  - `select current_time('1')`
- **SELECT CURRENT_TIMESTAMP** — 3
  - `select current_timestamp(-1)`
  - `select current_timestamp(1.0)`
  - `select current_timestamp('2')`
- **SELECT CURTIME** — 3
  - `select curtime(-1)`
  - `select curtime(1.0)`
  - `select curtime('1')`
- **SELECT DATE_ADD** — 3
  - `select date_add("2011-11-11 10:10:10.123456", 10)`
  - `select date_add("2011-11-11 10:10:10.123456", 0.10)`
  - `select date_add("2011-11-11 10:10:10.123456", "11,11")`
- **SELECT DATE_SUB** — 3
  - `select date_sub("2011-11-11 10:10:10.123456", 10)`
  - `select date_sub("2011-11-11 10:10:10.123456", 0.10)`
  - `select date_sub("2011-11-11 10:10:10.123456", "11,11")`
- **SELECT UTC_TIME** — 3
  - `select utc_time(-1)`
  - `select utc_time(1.0)`
  - `select utc_time('1')`
- **SELECT UTC_TIMESTAMP** — 3
  - `select utc_timestamp(-1)`
  - `select utc_timestamp(1.0)`
  - `select utc_timestamp('1')`
- **DELETE FROM INNER JOIN** — 2
  - `DELETE t1, t2 FROM t1 INNER JOIN t2 INNER JOIN t3 WHERE t1.id=t2.id AND t2.id=t3.id limit 10;`
  - `DELETE t1, t2 FROM t1 INNER JOIN t2 INNER JOIN t3 WHERE t1.id=t2.id AND t2.id=t3.id order by t1.id;`
- **SELECT** — 2
  - `select x'0xaa'`
  - `select 0X11`
- **SELECT FROM RIGHT JOIN** — 2
  - `select * from t1 right join t2 on t1.id = t2.id left join t3`
  - `select * from t1 right join t2 using (id) left join t3`
- **SELECT NOT EXISTS SELECT** — 2
  - `SELECT + NOT EXISTS (select 1)`
  - `SELECT - NOT EXISTS (select 1)`
- 37 further statements in 36 smaller families; read them in the corpus

### EXPLAIN forms (3)

- **EXPLAIN** — 2
  - `explain explore '23adc8e6f62'`
  - `explain explore 'digestxxx'`
- **EXPLAIN SELECT FROM** — 1
  - `explain explore 'select a from t'`

### users, grants and roles (1)

- **CREATE ROLE RESOURCE** — 1
  - `CREATE ROLE RESOURCE`

### SET and session state (2)

- **COMMIT AND CHAIN RELEASE** — 1
  - `COMMIT AND CHAIN RELEASE`
- **ROLLBACK AND CHAIN RELEASE** — 1
  - `ROLLBACK AND CHAIN RELEASE`

## What an executing measurement loses, and why this one does not

An earlier pass measured this by executing the corpus and stopped getting
verdicts partway through, with 2,775 of 3,902 statements unanswered. The
cause was never diagnosed. Replaying the corpus by execution reproduces
three separate ways it destroys its own run:

- **`shutdown`** succeeds and stops the server. Nothing after it gets a
  verdict, and no reconnect reaches a process that has exited. It sits at
  corpus position 3,386 of 3,902.
- **`SET PASSWORD = 'password';`** succeeds and changes the connecting
  user's password, at position 705. The connection that is already
  authenticated survives, so the run looks healthy; every new connection is
  refused with 1045, so the run dies the moment anything replaces it —
  which `kill CONNECTION_ID()`, at position 3,378, then forces.
- **`load stats '/tmp/stats.json'`** and the three `PLAN REPLAYER` file
  forms make the server ask the client for a local file. The driver refuses
  with an error that is not a server error at all, so a pass that reads any
  non-1064 server error as acceptance gets no verdict for them. These are
  the four cases the earlier pass reported as unanswered.

The corpus harness takes the engine's verdict from `PREPARE`, which parses
the statement and stops, so it runs none of them. Validated against direct
execution over 3,897 of the 3,902 statements, the two methods agree on
every one except the two corpus rows that are themselves `PREPARE`
statements, where nesting returns 1295 (not carried by the
prepared-statement protocol) instead of a verdict on the inner text; for
those two `PREPARE` gives the more accurate syntax verdict, since direct
execution's 1064 is about the runtime value of the user variable. Execution
lost five verdicts over the same range; `PREPARE` lost none.

Resilience is still pinned rather than assumed, because a container can die
for its own reasons: `TestTiDBCorpusHarnessSurvivesALostServer` kills the
connection and then the server, and the harness must come back with a
verdict, name the statement that preceded the loss, and report that it had
to recover.
