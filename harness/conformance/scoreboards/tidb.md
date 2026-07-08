# Conformance scoreboard: tidb

| meta | value |
|---|---|
| engine_version | v8.5.5 |
| omni_sha | f7991e553a27d880dcd70b6d5bcfb007981f3be8 |
| corpus_tag | v8.5.5 |
| container_digest | pingcap/tidb@sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 |
| classifier_version | 1 |

## Counts (upstream lane)

| class | statements |
|---|---|
| AGREE_ACCEPT | 2004 |
| AGREE_REJECT | 479 |
| GAP | 1093 |
| OVER | 160 |
| INDETERMINATE | 49 |
| SKIP | 13 |
| duplicates_dropped | 58 |
| duplicate_label_conflicts | 0 |
| total | 3798 |

GAP clusters: 71

OVER clusters: 46

Clusters are the work unit; statement counts are coverage context.

## GAP clusters (engine accepts, omni rejects) — the burn-down list

| n | family | count | divergence | exemplar | source |
|---|---|---|---|---|---|
| 1 | OTHER | 156 | syntax error at or near ? (line N, column N) | `trace begin` | corpus/tidb/pkg/parser/parser_test.go:5666 |
| 2 | SELECT | 112 | syntax error at or near ? (line N, column N) | `SELECT SCHEMA();` | corpus/tidb/pkg/parser/parser_test.go:1678 |
| 3 | ADMIN | 76 | syntax error at or near ? (line N, column N) | `admin show ddl;` | corpus/tidb/pkg/parser/parser_test.go:489 |
| 4 | CREATE OTHER | 73 | unexpected token after CREATE (line N, column N) | `create sequence seq` | corpus/tidb/pkg/parser/parser_test.go:3762 |
| 5 | ALTER TABLE | 68 | syntax error at or near ? (line N, column N) | `ALTER TABLE db.t RENAME db1.t1` | corpus/tidb/pkg/parser/parser_test.go:3144 |
| 6 | SHOW | 66 | syntax error at or near ? (line N, column N) | `show config` | corpus/tidb/pkg/parser/parser_test.go:1429 |
| 7 | DROP | 49 | unexpected token after DROP (line N, column N) | `drop stats t` | corpus/tidb/pkg/parser/parser_test.go:2836 |
| 8 | CREATE TABLE | 46 | unexpected token (line N, column N) | `create table t (a text byte)` | corpus/tidb/pkg/parser/parser_test.go:3735 |
| 9 | ALTER TABLE | 33 | expected identifier (line N, column N) | `ALTER TABLE t RENAME = t1` | corpus/tidb/pkg/parser/parser_test.go:3149 |
| 10 | CREATE OTHER | 33 | syntax error at or near ? (line N, column N) | `create resource group x ru_per_sec=2000` | corpus/tidb/pkg/parser/parser_test.go:3936 |
| 11 | ALTER TABLE | 30 | expected ALTER TABLE operation (line N, column N) | `ALTER TABLE tmp CACHE` | corpus/tidb/pkg/parser/parser_test.go:2546 |
| 12 | ALTER TABLE | 27 | unexpected token (line N, column N) | `ALTER TABLE d_n.t_n ADD PARTITION LOCAL` | corpus/tidb/pkg/parser/parser_test.go:3332 |
| 13 | SELECT | 27 | expected SELECT or ? (line N, column N) | `(TABLE t)` | corpus/tidb/pkg/parser/parser_test.go:650 |
| 14 | STATS | 22 | syntax error at or near ? (line N, column N) | `analyze table t1 index` | corpus/tidb/pkg/parser/parser_test.go:6087 |
| 15 | CREATE TABLE | 20 | syntax error at or near ? (line N, column N) | `CREATE TABLE t (a int) STATS_TOPN=1` | corpus/tidb/pkg/parser/parser_test.go:4057 |
| 16 | CREATE TABLE | 19 | expected data type (line N, column N) | `CREATE TABLE t (a VECTOR)` | corpus/tidb/pkg/parser/parser_test.go:7688 |
| 17 | LOAD | 17 | syntax error at or near ? (line N, column N) | `import into t from '/file.csv'` | corpus/tidb/pkg/parser/parser_test.go:781 |
| 18 | EXPLAIN | 16 | expected identifier or keyword (line N, column N) | `EXPLAIN FORMAT = 'ROW' SELECT 1` | corpus/tidb/pkg/parser/parser_test.go:5603 |
| 19 | SELECT | 16 | expected identifier (line N, column N) | `select * from 1db.1table;` | corpus/tidb/pkg/parser/parser_test.go:1126 |
| 20 | TXN | 13 | syntax error at or near ? (line N, column N) | `begin optimistic` | corpus/tidb/pkg/parser/parser_test.go:1099 |
| 21 | ALTER DATABASE | 12 | expected identifier (line N, column N) | `ALTER DATABASE t COLLATE = binary` | corpus/tidb/pkg/parser/parser_test.go:3056 |
| 22 | OTHER | 11 | unexpected token after ALTER (line N, column N) | `alter sequence seq restart` | corpus/tidb/pkg/parser/parser_test.go:3851 |
| 23 | ALTER TABLE | 10 | expected BY after PARTITION (line N, column N) | `ALTER TABLE t PARTITION p ATTRIBUTES='str'` | corpus/tidb/pkg/parser/parser_test.go:3357 |
| 24 | OTHER | 10 | expected identifier (line N, column N) | `ALTER SCHEMA t DEFAULT CHARSET = 'UTF8'` | corpus/tidb/pkg/parser/parser_test.go:3054 |
| 25 | INSERT | 8 | syntax error at or near ? (line N, column N) | `INSERT INTO foo VALUE ()` | corpus/tidb/pkg/parser/parser_test.go:596 |
| 26 | SELECT | 8 | unexpected token (line N, column N) | `Select (1, 1) > (1, 1)` | corpus/tidb/pkg/parser/parser_test.go:1728 |
| 27 | SET | 8 | unexpected token (line N, column N) | `SET SESSION_STATES 'x'` | corpus/tidb/pkg/parser/parser_test.go:1154 |
| 28 | STATS | 8 | expected identifier (line N, column N) | `analyze table t with 4 topn` | corpus/tidb/pkg/parser/parser_test.go:6091 |
| 29 | DROP | 6 | syntax error at or near ? (line N, column N) | `drop index idx on t lock default` | corpus/tidb/pkg/parser/parser_test.go:3445 |
| 30 | ALTER TABLE | 5 | expected data type (line N, column N) | `ALTER TABLE t ADD VECTOR INDEX ((lower(a))) USING HNSW COMMENT 'a'` | corpus/tidb/pkg/parser/parser_test.go:3222 |
| 31 | CREATE INDEX | 5 | syntax error at or near ? (line N, column N) | `create index i on t (a) local` | corpus/tidb/pkg/parser/parser_test.go:3045 |
| 32 | DCL | 5 | expected privilege name (line N, column N) | `GRANT 'u1' TO 'u1';` | corpus/tidb/pkg/parser/parser_test.go:5111 |
| 33 | DROP | 5 | expected identifier (line N, column N) | `drop index if exists a on t` | corpus/tidb/pkg/parser/parser_test.go:3437 |
| 34 | DCL | 4 | syntax error at or near ? (line N, column N) | `grant all privileges on zabbix.* to 'zabbix'@'localhost' identified by 'password';` | corpus/tidb/pkg/parser/parser_test.go:5106 |
| 35 | INSERT | 4 | expected identifier (line N, column N) | `INSERT INTO foo () VALUES ()` | corpus/tidb/pkg/parser/parser_test.go:595 |
| 36 | SET | 4 | expected identifier (line N, column N) | `set names binary` | corpus/tidb/pkg/parser/parser_test.go:1405 |
| 37 | ALTER TABLE | 3 | syntax error at end of input (line N, column N) | `ALTER TABLE employees ADD PARTITION` | corpus/tidb/pkg/parser/parser_test.go:3090 |
| 38 | CREATE TABLE | 3 | expected CHAR or VARCHAR after NATIONAL (line N, column N) | `create table t (a national character);` | corpus/tidb/pkg/parser/parser_test.go:3591 |
| 39 | CREATE TABLE | 3 | expected HASH, KEY, RANGE, or LIST after PARTITION BY (line N, column N) | `create table t1 (a int) partition by system_time (partition x history, partition y current)` | corpus/tidb/pkg/parser/parser_test.go:6391 |
| 40 | CREATE VIEW | 3 | expected SELECT or ? (line N, column N) | `CREATE VIEW v AS TABLE t` | corpus/tidb/pkg/parser/parser_test.go:666 |
| 41 | DELETE | 3 | expected identifier (line N, column N) | `delete delayed from t where a = 2` | corpus/tidb/pkg/parser/parser_test.go:5525 |
| 42 | INSERT | 3 | unexpected token (line N, column N) | `insert into t set a := default` | corpus/tidb/pkg/parser/parser_test.go:3724 |
| 43 | LOAD | 3 | expected SELECT or ? (line N, column N) | `import into t (a,b) from '/file.csv'` | corpus/tidb/pkg/parser/parser_test.go:782 |
| 44 | SELECT | 3 | expected UPDATE or SHARE after FOR (line N, column N) | `select next value for seq` | corpus/tidb/pkg/parser/parser_test.go:2318 |
| 45 | SET | 3 | syntax error at or near ? (line N, column N) | `set names utf8, @@session.sql_mode=1;` | corpus/tidb/pkg/parser/parser_test.go:1419 |
| 46 | CREATE OTHER | 2 | expected identifier (line N, column N) | `create resource group default ru_per_sec=20, priority=LOW, burstable` | corpus/tidb/pkg/parser/parser_test.go:3949 |
| 47 | CREATE TABLE | 2 | expected string literal in value list (line N, column N) | `create table t (c1 enum(0x61, 'b'), c2 set(0x61, 'b'))` | corpus/tidb/pkg/parser/parser_test.go:4981 |
| 48 | DCL | 2 | unexpected token (line N, column N) | `grant grant option on *.* to u1` | corpus/tidb/pkg/parser/parser_test.go:5125 |
| 49 | EXPLAIN | 2 | expected identifier (line N, column N) | `EXPLAIN ALTER TABLE t1 ADD INDEX (a)` | corpus/tidb/pkg/parser/parser_test.go:5629 |
| 50 | LOAD | 2 | unexpected token (line N, column N) | `load data infile '/tmp/t.csv' into table 't' with threads=10` | corpus/tidb/pkg/parser/parser_test.go:777 |
| 51 | SELECT | 2 | expected data type (line N, column N) | `SELECT CAST(data AS CHARACTER);` | corpus/tidb/pkg/parser/parser_test.go:1734 |
| 52 | SET | 2 | expected variable reference (line N, column N) | `SET @ = 1` | corpus/tidb/pkg/parser/parser_test.go:1345 |
| 53 | UPDATE | 2 | expected identifier (line N, column N) | `update delayed t set a = 2` | corpus/tidb/pkg/parser/parser_test.go:5522 |
| 54 | ALTER TABLE | 1 | expected CHAR or VARCHAR after NATIONAL (line N, column N) | `alter table t_n storage disk , modify ident national varcharacter(12) column_format fixed first;` | corpus/tidb/pkg/parser/parser_test.go:3603 |
| 55 | ALTER TABLE | 1 | expected ROLE after SET DEFAULT (line N, column N) | `alter table d_n.t_n convert to char set default` | corpus/tidb/pkg/parser/parser_test.go:3263 |
| 56 | CREATE DATABASE | 1 | expected identifier (line N, column N) | `create database 123test` | corpus/tidb/pkg/parser/parser_test.go:2374 |
| 57 | CREATE INDEX | 1 | unexpected token (line N, column N) | `CREATE UNIQUE INDEX ident TYPE BTREE ON d_n.t_n ( ident , ident ASC )` | corpus/tidb/pkg/parser/parser_test.go:3377 |
| 58 | CREATE TABLE | 1 | expected SELECT or ? (line N, column N) | `CREATE TABLE ta AS VALUES ROW(1)` | corpus/tidb/pkg/parser/parser_test.go:684 |
| 59 | CREATE TABLE | 1 | expected identifier (line N, column N) | `create table '123' (123a1 int)` | corpus/tidb/pkg/parser/parser_test.go:2378 |
| 60 | EXPLAIN | 1 | syntax error at or near ? (line N, column N) | `explain replace into foo values (1 \|\| 2)` | corpus/tidb/pkg/parser/parser_test.go:5587 |
| 61 | LOAD | 1 | syntax error at end of input (line N, column N) | `load data infile '/tmp/t.csv' into table t with detached` | corpus/tidb/pkg/parser/parser_test.go:775 |
| 62 | SELECT | 1 | expected INDEX or KEY (line N, column N) | `SELECT NTH_VALUE(val, 233) FROM LAST IGNORE NULLS OVER w FROM t;` | corpus/tidb/pkg/parser/parser_test.go:6547 |
| 63 | SELECT | 1 | invalid INTERVAL unit: sql_tsi_day (line N, column N) | `select date_add("2011-11-11 10:10:10.123456", interval 10 sql_tsi_day)` | corpus/tidb/pkg/parser/parser_test.go:2074 |
| 64 | SELECT | 1 | invalid INTERVAL unit: sql_tsi_hour (line N, column N) | `select date_add("2011-11-11 10:10:10.123456", interval 10 sql_tsi_hour)` | corpus/tidb/pkg/parser/parser_test.go:2073 |
| 65 | SELECT | 1 | invalid INTERVAL unit: sql_tsi_minute (line N, column N) | `select date_add("2011-11-11 10:10:10.123456", interval 10 sql_tsi_minute)` | corpus/tidb/pkg/parser/parser_test.go:2072 |
| 66 | SELECT | 1 | invalid INTERVAL unit: sql_tsi_month (line N, column N) | `select date_add("2011-11-11 10:10:10.123456", interval 1 sql_tsi_month)` | corpus/tidb/pkg/parser/parser_test.go:2076 |
| 67 | SELECT | 1 | invalid INTERVAL unit: sql_tsi_quarter (line N, column N) | `select date_add("2011-11-11 10:10:10.123456", interval 1 sql_tsi_quarter)` | corpus/tidb/pkg/parser/parser_test.go:2077 |
| 68 | SELECT | 1 | invalid INTERVAL unit: sql_tsi_second (line N, column N) | `select date_add("2011-11-11 10:10:10.123456", interval 10 sql_tsi_second)` | corpus/tidb/pkg/parser/parser_test.go:2071 |
| 69 | SELECT | 1 | invalid INTERVAL unit: sql_tsi_week (line N, column N) | `select date_add("2011-11-11 10:10:10.123456", interval 1 sql_tsi_week)` | corpus/tidb/pkg/parser/parser_test.go:2075 |
| 70 | SELECT | 1 | invalid INTERVAL unit: sql_tsi_year (line N, column N) | `select date_add("2011-11-11 10:10:10.123456", interval 1 sql_tsi_year)` | corpus/tidb/pkg/parser/parser_test.go:2078 |
| 71 | SET | 1 | syntax error at end of input (line N, column N) | `SET @@character_set_results = binary` | corpus/tidb/pkg/parser/parser_test.go:1381 |

## OVER clusters (engine rejects, omni accepts) — triage: structural vs leniency

| n | family | count | divergence | exemplar | source |
|---|---|---|---|---|---|
| 1 | SELECT | 61 | You have an error in your SQL syntax; check the manual that corresponds to your ... | `select row(1)` | corpus/tidb/pkg/parser/parser_test.go:1724 |
| 2 | CREATE TABLE | 15 | You have an error in your SQL syntax; check the manual that corresponds to your ... | `CREATE TABLE foo ()` | corpus/tidb/pkg/parser/parser_test.go:2515 |
| 3 | ALTER TABLE | 13 | You have an error in your SQL syntax; check the manual that corresponds to your ... | `ALTER TABLE t_n ALTER CONSTRAINT ident` | corpus/tidb/pkg/parser/parser_test.go:3294 |
| 4 | CREATE TABLE | 7 | Inconsistency in usage of column lists for partitioning | `create table t (id int) partition by range columns (id) (partition p0 values less than (1, 2))` | corpus/tidb/pkg/parser/parser_test.go:6405 |
| 5 | CREATE TABLE | 4 | Only RANGE PARTITIONING can use VALUES LESS THAN in partition definition | `create table t1 (a int) partition by key (a) (partition x values less than (10))` | corpus/tidb/pkg/parser/parser_test.go:6370 |
| 6 | CREATE TABLE | 3 | Only LIST PARTITIONING can use VALUES IN in partition definition | `create table t1 (a int) partition by key (a) (partition x values in (10))` | corpus/tidb/pkg/parser/parser_test.go:6379 |
| 7 | DELETE | 3 | You have an error in your SQL syntax; check the manual that corresponds to your ... | `DELETE from t1.*` | corpus/tidb/pkg/parser/parser_test.go:851 |
| 8 | ALTER TABLE | 2 | Incorrect usage of DEFAULT and generated column | `alter table t modify column f int as (a+1) default 55;` | corpus/tidb/pkg/parser/parser_test.go:3507 |
| 9 | CREATE INDEX | 2 | Unknown ALGORITHM ? | `CREATE INDEX idx ON t ( a ) ALGORITHM ident` | corpus/tidb/pkg/parser/parser_test.go:3431 |
| 10 | CREATE TABLE | 2 | Wrong number of subpartitions defined, mismatch with previous setting | `create table t1 (a int, b int) partition by range (a) subpartition by hash (b) subpartitions 2 (part...` | corpus/tidb/pkg/parser/parser_test.go:6431 |
| 11 | CREATE TABLE | 2 | You have an error in your SQL syntax; check the manual that corresponds to your ... | `CREATE TABLE t  (c1 integer ,c2 integer) PARTITION BY LINEAR KEY ALGORITHM = 0 (c1,c2) PARTITIONS 4` | corpus/tidb/pkg/parser/parser_test.go:7215 |
| 12 | CREATE TABLE | 2 | You have an error in your SQL syntax; check the manual that corresponds to your ... | `create table t (created_at datetime) TTL_ENABLE = 'test_case'` | corpus/tidb/pkg/parser/parser_test.go:7619 |
| 13 | CREATE TABLE | 2 | You have an error in your SQL syntax; check the manual that corresponds to your ... | `create table t (created_at datetime) TTL_JOB_INTERVAL = '10hourxx'` | corpus/tidb/pkg/parser/parser_test.go:7625 |
| 14 | CREATE TABLE | 2 | You have an error in your SQL syntax; check the manual that corresponds to your ... | `create table t (a int) stats_auto_recalc 2;` | corpus/tidb/pkg/parser/parser_test.go:3636 |
| 15 | DELETE | 2 | You have an error in your SQL syntax; check the manual that corresponds to your ... | `delete from t where a = 7 or 1=1/*' and b = 'p'` | corpus/tidb/pkg/parser/parser_test.go:5161 |
| 16 | DROP | 2 | Unknown ALGORITHM ? | `drop index idx on t algorithm algorithm_type` | corpus/tidb/pkg/parser/parser_test.go:3453 |
| 17 | DROP | 2 | Unknown LOCK type ? | `drop index idx on t lock lock_type` | corpus/tidb/pkg/parser/parser_test.go:3455 |
| 18 | LOAD | 2 | Field separator argument is not what is expected; check the manual | `load data infile '/tmp/t.csv' into table t fields escaped by 'aa'` | corpus/tidb/pkg/parser/parser_test.go:756 |
| 19 | SELECT | 2 | Incorrect usage of ALL and DISTINCT | `SELECT DISTINCT ALL * FROM t` | corpus/tidb/pkg/parser/parser_test.go:587 |
| 20 | SELECT | 2 | You have an error in your SQL syntax; check the manual that corresponds to your ... | `select date_add("2011-11-11 10:10:10.123456", "11,11")` | corpus/tidb/pkg/parser/parser_test.go:2068 |
| 21 | TXN | 2 | You have an error in your SQL syntax; check the manual that corresponds to your ... | `COMMIT AND CHAIN RELEASE` | corpus/tidb/pkg/parser/parser_test.go:626 |
| 22 | UPDATE | 2 | You have an error in your SQL syntax; check the manual that corresponds to your ... | `UPDATE items,month SET items.price=month.price WHERE items.id=month.id LIMIT 10;` | corpus/tidb/pkg/parser/parser_test.go:933 |
| 23 | ALTER TABLE | 1 | For RANGE partitions each partition must be defined | `alter table t partition by range(a)` | corpus/tidb/pkg/parser/parser_test.go:3309 |
| 24 | ALTER TABLE | 1 | You have an error in your SQL syntax; check the manual that corresponds to your ... | `alter table t /*T![ttl] TTL_ENABLE = 'test_case' */` | corpus/tidb/pkg/parser/parser_test.go:7621 |
| 25 | CREATE INDEX | 1 | You have an error in your SQL syntax; check the manual that corresponds to your ... | `CREATE UNIQUE INDEX ident USING HNSW ON d_n.t_n ( ident , ident ASC )` | corpus/tidb/pkg/parser/parser_test.go:3417 |
| 26 | CREATE OTHER | 1 | You have an error in your SQL syntax; check the manual that corresponds to your ... | `CREATE ROLE RESOURCE` | corpus/tidb/pkg/parser/parser_test.go:4050 |
| 27 | CREATE TABLE | 1 | Cannot have more than one value for this type of RANGE partitioning | `create table t1 (a int, b int) partition by range (a) (partition x values less than (10, 20))` | corpus/tidb/pkg/parser/parser_test.go:6404 |
| 28 | CREATE TABLE | 1 | For LIST partitions each partition must be defined | `create table t1 (a int) partition by list (a)` | corpus/tidb/pkg/parser/parser_test.go:6397 |
| 29 | CREATE TABLE | 1 | For RANGE partitions each partition must be defined | `create table t1 (a int) partition by range (a)` | corpus/tidb/pkg/parser/parser_test.go:6396 |
| 30 | CREATE TABLE | 1 | Incorrect usage of AUTO_INCREMENT and generated column | `create table t (a bigint, b bigint as (a) primary key auto_increment);` | corpus/tidb/pkg/parser/parser_test.go:3500 |
| 31 | CREATE TABLE | 1 | Incorrect usage of DEFAULT and generated column | `create table t (a bigint, b bigint as (a) not null default 10);` | corpus/tidb/pkg/parser/parser_test.go:3501 |
| 32 | CREATE TABLE | 1 | Incorrect usage of ON UPDATE and generated column | `create table t (a timestamp, b timestamp as (a) not null on update current_timestamp);` | corpus/tidb/pkg/parser/parser_test.go:3499 |
| 33 | CREATE TABLE | 1 | It is only possible to mix RANGE/LIST partitioning with HASH/KEY partitioning fo... | `create table t1 (a int, b int) partition by range (a) (partition x values less than (10) (subpartiti...` | corpus/tidb/pkg/parser/parser_test.go:6441 |
| 34 | CREATE TABLE | 1 | Number of partitions = N is not an allowed value | `create table t1 (a int) partition by hash (a) partitions 0` | corpus/tidb/pkg/parser/parser_test.go:6442 |
| 35 | CREATE TABLE | 1 | Number of subpartitions = N is not an allowed value | `create table t1 (a int, b int) partition by range (a) subpartition by hash (b) subpartitions 0 (part...` | corpus/tidb/pkg/parser/parser_test.go:6443 |
| 36 | CREATE TABLE | 1 | Row expressions in VALUES IN only allowed for multi-field column partitioning | `create table t1 (a int, b int) partition by list (a) (partition x values in ((10, 20)))` | corpus/tidb/pkg/parser/parser_test.go:6409 |
| 37 | CREATE TABLE | 1 | Syntax : LIST PARTITIONING requires definition of VALUES IN for each partition | `create table t1 (a int) partition by list (a) (partition x, partition y)` | corpus/tidb/pkg/parser/parser_test.go:6366 |
| 38 | CREATE TABLE | 1 | Syntax : RANGE PARTITIONING requires definition of VALUES LESS THAN for each par... | `create table t1 (a int) partition by range (a) (partition x, partition y)` | corpus/tidb/pkg/parser/parser_test.go:6365 |
| 39 | CREATE TABLE | 1 | Wrong number of partitions defined, mismatch with previous setting | `create table t1 (a int) partition by hash (a) partitions 2 (partition x)` | corpus/tidb/pkg/parser/parser_test.go:6429 |
| 40 | CREATE TABLE | 1 | You have an error in your SQL syntax; check the manual that corresponds to your ... | `create table t (created_at datetime) TTL_JOB_INTERVAL = '10.10.255h'` | corpus/tidb/pkg/parser/parser_test.go:7626 |
| 41 | INSERT | 1 | You have an error in your SQL syntax; check the manual that corresponds to your ... | `INSERT INTO t (a) SET a=1` | corpus/tidb/pkg/parser/parser_test.go:918 |
| 42 | SELECT | 1 | Incorrect arguments to ESCAPE | `select "abc_" like "abc\\_" escape '\|\|'` | corpus/tidb/pkg/parser/parser_test.go:5443 |
| 43 | SELECT | 1 | Too-big precision N specified for ?. Maximum is N. | `select cast(1 as float(54));` | corpus/tidb/pkg/parser/parser_test.go:1753 |
| 44 | SELECT | 1 | You have an error in your SQL syntax; check the manual that corresponds to your ... | `select """;` | corpus/tidb/pkg/parser/parser_test.go:5571 |
| 45 | SELECT | 1 | You have an error in your SQL syntax; check the manual that corresponds to your ... | `select 0X11` | corpus/tidb/pkg/parser/parser_test.go:4968 |
| 46 | SELECT | 1 | You have an error in your SQL syntax; check the manual that corresponds to your ... | `select 1/*` | corpus/tidb/pkg/parser/parser_test.go:5188 |

