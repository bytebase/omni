# Conformance scoreboard: tidb

| meta | value |
|---|---|
| engine_version | v8.5.5 |
| omni_sha | 0df0e512aeff928e4bb9526e2b004142bbcf61f0 |
| corpus_tag | v8.5.5 |
| container_digest | - |
| classifier_version | 1 |

## Counts (upstream lane)

| class | statements |
|---|---|
| AGREE_ACCEPT | 2001 |
| AGREE_REJECT | 479 |
| GAP | 1141 |
| OVER | 159 |
| INDETERMINATE | 0 |
| SKIP | 1 |
| duplicates_dropped | 57 |
| duplicate_label_conflicts | 0 |
| total | 3781 |

GAP clusters: 71

OVER clusters: 101

Clusters are the work unit; statement counts are coverage context.

## GAP clusters (engine accepts, omni rejects) — the burn-down list

| n | family | count | divergence | exemplar | source |
|---|---|---|---|---|---|
| 1 | OTHER | 179 | syntax error at or near ? (line N, column N) | `trace begin` | corpus/tidb/pkg/parser/parser_test.go:5666 |
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
| 15 | DCL | 21 | syntax error at or near ? (line N, column N) | `CREATE USER uesr1@LOCALhost` | corpus/tidb/pkg/parser/parser_test.go:5024 |
| 16 | CREATE TABLE | 20 | syntax error at or near ? (line N, column N) | `CREATE TABLE t (a int) STATS_TOPN=1` | corpus/tidb/pkg/parser/parser_test.go:4057 |
| 17 | CREATE TABLE | 19 | expected data type (line N, column N) | `CREATE TABLE t (a VECTOR)` | corpus/tidb/pkg/parser/parser_test.go:7688 |
| 18 | LOAD | 17 | syntax error at or near ? (line N, column N) | `import into t from '/file.csv'` | corpus/tidb/pkg/parser/parser_test.go:781 |
| 19 | EXPLAIN | 16 | expected identifier or keyword (line N, column N) | `EXPLAIN FORMAT = 'ROW' SELECT 1` | corpus/tidb/pkg/parser/parser_test.go:5603 |
| 20 | SELECT | 16 | expected identifier (line N, column N) | `select * from 1db.1table;` | corpus/tidb/pkg/parser/parser_test.go:1126 |
| 21 | SET | 14 | unexpected token (line N, column N) | `SET SESSION_STATES 'x'` | corpus/tidb/pkg/parser/parser_test.go:1154 |
| 22 | TXN | 13 | syntax error at or near ? (line N, column N) | `begin optimistic` | corpus/tidb/pkg/parser/parser_test.go:1099 |
| 23 | ALTER DATABASE | 12 | expected identifier (line N, column N) | `ALTER DATABASE t COLLATE = binary` | corpus/tidb/pkg/parser/parser_test.go:3056 |
| 24 | OTHER | 11 | expected identifier (line N, column N) | `ALTER SCHEMA t DEFAULT CHARSET = 'UTF8'` | corpus/tidb/pkg/parser/parser_test.go:3054 |
| 25 | OTHER | 11 | unexpected token after ALTER (line N, column N) | `alter sequence seq restart` | corpus/tidb/pkg/parser/parser_test.go:3851 |
| 26 | ALTER TABLE | 10 | expected BY after PARTITION (line N, column N) | `ALTER TABLE t PARTITION p ATTRIBUTES='str'` | corpus/tidb/pkg/parser/parser_test.go:3357 |
| 27 | INSERT | 8 | syntax error at or near ? (line N, column N) | `INSERT INTO foo VALUE ()` | corpus/tidb/pkg/parser/parser_test.go:596 |
| 28 | SELECT | 8 | unexpected token (line N, column N) | `Select (1, 1) > (1, 1)` | corpus/tidb/pkg/parser/parser_test.go:1728 |
| 29 | STATS | 8 | expected identifier (line N, column N) | `analyze table t with 4 topn` | corpus/tidb/pkg/parser/parser_test.go:6091 |
| 30 | DROP | 6 | syntax error at or near ? (line N, column N) | `drop index idx on t lock default` | corpus/tidb/pkg/parser/parser_test.go:3445 |
| 31 | ALTER TABLE | 5 | expected data type (line N, column N) | `ALTER TABLE t ADD VECTOR INDEX ((lower(a))) USING HNSW COMMENT 'a'` | corpus/tidb/pkg/parser/parser_test.go:3222 |
| 32 | CREATE INDEX | 5 | syntax error at or near ? (line N, column N) | `create index i on t (a) local` | corpus/tidb/pkg/parser/parser_test.go:3045 |
| 33 | DCL | 5 | expected privilege name (line N, column N) | `GRANT 'u1' TO 'u1';` | corpus/tidb/pkg/parser/parser_test.go:5111 |
| 34 | DROP | 5 | expected identifier (line N, column N) | `drop index if exists a on t` | corpus/tidb/pkg/parser/parser_test.go:3437 |
| 35 | INSERT | 4 | expected identifier (line N, column N) | `INSERT INTO foo () VALUES ()` | corpus/tidb/pkg/parser/parser_test.go:595 |
| 36 | SET | 4 | expected identifier (line N, column N) | `set names binary` | corpus/tidb/pkg/parser/parser_test.go:1405 |
| 37 | ALTER TABLE | 3 | syntax error at end of input (line N, column N) | `ALTER TABLE employees ADD PARTITION` | corpus/tidb/pkg/parser/parser_test.go:3090 |
| 38 | CREATE TABLE | 3 | expected CHAR or VARCHAR after NATIONAL (line N, column N) | `create table t (a national character);` | corpus/tidb/pkg/parser/parser_test.go:3591 |
| 39 | CREATE TABLE | 3 | expected HASH, KEY, RANGE, or LIST after PARTITION BY (line N, column N) | `create table t1 (a int) partition by system_time (partition x history, partition y current)` | corpus/tidb/pkg/parser/parser_test.go:6391 |
| 40 | CREATE VIEW | 3 | expected SELECT or ? (line N, column N) | `CREATE VIEW v AS TABLE t` | corpus/tidb/pkg/parser/parser_test.go:666 |
| 41 | DELETE | 3 | expected identifier (line N, column N) | `delete delayed from t where a = 2` | corpus/tidb/pkg/parser/parser_test.go:5525 |
| 42 | INSERT | 3 | unexpected token (line N, column N) | `insert into t set a := default` | corpus/tidb/pkg/parser/parser_test.go:3724 |
| 43 | LOAD | 3 | expected SELECT or ? (line N, column N) | `import into t (a,b) from '/file.csv'` | corpus/tidb/pkg/parser/parser_test.go:782 |
| 44 | LOAD | 3 | unexpected token (line N, column N) | `load stats '/tmp/stats.json'` | corpus/tidb/pkg/parser/parser_test.go:1334 |
| 45 | SELECT | 3 | expected UPDATE or SHARE after FOR (line N, column N) | `select next value for seq` | corpus/tidb/pkg/parser/parser_test.go:2318 |
| 46 | SET | 3 | syntax error at or near ? (line N, column N) | `set names utf8, @@session.sql_mode=1;` | corpus/tidb/pkg/parser/parser_test.go:1419 |
| 47 | CREATE OTHER | 2 | expected identifier (line N, column N) | `create resource group default ru_per_sec=20, priority=LOW, burstable` | corpus/tidb/pkg/parser/parser_test.go:3949 |
| 48 | CREATE TABLE | 2 | expected string literal in value list (line N, column N) | `create table t (c1 enum(0x61, 'b'), c2 set(0x61, 'b'))` | corpus/tidb/pkg/parser/parser_test.go:4981 |
| 49 | DCL | 2 | unexpected token (line N, column N) | `grant grant option on *.* to u1` | corpus/tidb/pkg/parser/parser_test.go:5125 |
| 50 | EXPLAIN | 2 | expected identifier (line N, column N) | `EXPLAIN ALTER TABLE t1 ADD INDEX (a)` | corpus/tidb/pkg/parser/parser_test.go:5629 |
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
| 1 | CREATE TABLE | 26 | CREATE TABLE TN (A | `create table t1 (a int) partition by list (a)` | corpus/tidb/pkg/parser/parser_test.go:6397 |
| 2 | CREATE TABLE | 8 | CREATE TABLE T (A | `create table t (a national varchar);` | corpus/tidb/pkg/parser/parser_test.go:3595 |
| 3 | CREATE TABLE | 7 | CREATE TABLE T (B | `create table t (b int) partition by hash ( not b );` | corpus/tidb/pkg/parser/parser_test.go:2625 |
| 4 | CREATE TABLE | 5 | CREATE TABLE T (CREATED_AT | `create table t (created_at datetime) TTL_ENABLE = 'test_case'` | corpus/tidb/pkg/parser/parser_test.go:7619 |
| 5 | ALTER TABLE | 4 | ALTER TABLE T ADD | `alter table t add column a int, remove partitioning` | corpus/tidb/pkg/parser/parser_test.go:3565 |
| 6 | DROP | 4 | DROP INDEX IDX ON | `drop index idx on t lock lock_type` | corpus/tidb/pkg/parser/parser_test.go:3455 |
| 7 | ALTER TABLE | 3 | ALTER TABLE T ALTER | `ALTER TABLE t ALTER COLUMN a SET DEFAULT 1+1` | corpus/tidb/pkg/parser/parser_test.go:3179 |
| 8 | ALTER TABLE | 2 | ALTER TABLE T PARTITION | `alter table t partition by range(a)` | corpus/tidb/pkg/parser/parser_test.go:3309 |
| 9 | CREATE INDEX | 2 | CREATE INDEX IDX ON | `CREATE INDEX idx ON t ( a ) ALGORITHM ident` | corpus/tidb/pkg/parser/parser_test.go:3431 |
| 10 | CREATE TABLE | 2 | CREATE TABLE T (C | `create table t (c int) comment comment` | corpus/tidb/pkg/parser/parser_test.go:5155 |
| 11 | CREATE TABLE | 2 | CREATE TABLE T (CN | `CREATE TABLE t  (c1 integer ,c2 integer) PARTITION BY LINEAR KEY ALGORITHM = 0 (c1,c2) PARTITIONS 4` | corpus/tidb/pkg/parser/parser_test.go:7215 |
| 12 | DELETE | 2 | DELETE TN, TN FROM | `DELETE t1, t2 FROM t1 INNER JOIN t2 INNER JOIN t3 WHERE t1.id=t2.id AND t2.id=t3.id limit 10;` | corpus/tidb/pkg/parser/parser_test.go:907 |
| 13 | LOAD | 2 | LOAD DATA INFILE '/TMP/T.CSV' | `load data infile '/tmp/t.csv' into table t fields escaped by 'aa'` | corpus/tidb/pkg/parser/parser_test.go:756 |
| 14 | SELECT | 2 | SELECT * FROM TN | `select * from t1 right join t2 using (id) left join t3` | corpus/tidb/pkg/parser/parser_test.go:833 |
| 15 | UPDATE | 2 | UPDATE ITEMS,MONTH SET ITEMS.PRICE=MONTH.PRICE | `UPDATE items,month SET items.price=month.price WHERE items.id=month.id LIMIT 10;` | corpus/tidb/pkg/parser/parser_test.go:933 |
| 16 | ALTER TABLE | 1 | ALTER TABLE T /*T![TTL] | `alter table t /*T![ttl] TTL_ENABLE = 'test_case' */` | corpus/tidb/pkg/parser/parser_test.go:7621 |
| 17 | ALTER TABLE | 1 | ALTER TABLE T ENABLE | `alter table t enable keys, comment = 'cmt', partition by hash(a)` | corpus/tidb/pkg/parser/parser_test.go:3315 |
| 18 | ALTER TABLE | 1 | ALTER TABLE T MODIFY | `alter table t modify column f int as (a+1) default 55;` | corpus/tidb/pkg/parser/parser_test.go:3507 |
| 19 | ALTER TABLE | 1 | ALTER TABLE T REMOVE | `alter table t remove partitioning, add column a int` | corpus/tidb/pkg/parser/parser_test.go:3568 |
| 20 | ALTER TABLE | 1 | ALTER TABLE T_N ALTER | `ALTER TABLE t_n ALTER CONSTRAINT ident` | corpus/tidb/pkg/parser/parser_test.go:3294 |
| 21 | ALTER TABLE | 1 | ALTER TABLE T_N OPTIMIZE | `ALTER TABLE t_n OPTIMIZE PARTITION LOCAL` | corpus/tidb/pkg/parser/parser_test.go:6292 |
| 22 | ALTER TABLE | 1 | ALTER TABLE T_N REBUILD | `ALTER TABLE t_n REBUILD PARTITION LOCAL` | corpus/tidb/pkg/parser/parser_test.go:3113 |
| 23 | ALTER TABLE | 1 | ALTER TABLE T_N REPAIR | `ALTER TABLE t_n REPAIR PARTITION LOCAL` | corpus/tidb/pkg/parser/parser_test.go:6310 |
| 24 | CREATE INDEX | 1 | CREATE UNIQUE INDEX IDENT | `CREATE UNIQUE INDEX ident USING HNSW ON d_n.t_n ( ident , ident ASC )` | corpus/tidb/pkg/parser/parser_test.go:3417 |
| 25 | CREATE OTHER | 1 | CREATE ROLE RESOURCE | `CREATE ROLE RESOURCE` | corpus/tidb/pkg/parser/parser_test.go:4050 |
| 26 | CREATE TABLE | 1 | CREATE TABLE FOO () | `CREATE TABLE foo ()` | corpus/tidb/pkg/parser/parser_test.go:2515 |
| 27 | CREATE TABLE | 1 | CREATE TABLE FOO (); | `CREATE TABLE foo ();` | corpus/tidb/pkg/parser/parser_test.go:2516 |
| 28 | CREATE TABLE | 1 | CREATE TABLE T (ID | `create table t (id int) partition by range columns (id) (partition p0 values less than (1, 2))` | corpus/tidb/pkg/parser/parser_test.go:6405 |
| 29 | DELETE | 1 | DELETE FROM T WHERE | `delete from t where a = 7 or 1=1/*' and b = 'p'` | corpus/tidb/pkg/parser/parser_test.go:5161 |
| 30 | DELETE | 1 | DELETE FROM TN.* | `DELETE from t1.*` | corpus/tidb/pkg/parser/parser_test.go:851 |
| 31 | INSERT | 1 | INSERT INTO T (A) | `INSERT INTO t (a) SET a=1` | corpus/tidb/pkg/parser/parser_test.go:918 |
| 32 | SELECT | 1 | SELECT """; | `select """;` | corpus/tidb/pkg/parser/parser_test.go:5571 |
| 33 | SELECT | 1 | SELECT "ABC_" LIKE "ABC\\_" | `select "abc_" like "abc\\_" escape '\|\|'` | corpus/tidb/pkg/parser/parser_test.go:5443 |
| 34 | SELECT | 1 | SELECT '{}'->'$.A' FROM T | `SELECT '{}'->'$.a' FROM t` | corpus/tidb/pkg/parser/parser_test.go:2299 |
| 35 | SELECT | 1 | SELECT '{}'->>'$.A' FROM T | `SELECT '{}'->>'$.a' FROM t` | corpus/tidb/pkg/parser/parser_test.go:2300 |
| 36 | SELECT | 1 | SELECT + NOT EXISTS | `SELECT + NOT EXISTS (select 1)` | corpus/tidb/pkg/parser/parser_test.go:5232 |
| 37 | SELECT | 1 | SELECT - NOT EXISTS | `SELECT - NOT EXISTS (select 1)` | corpus/tidb/pkg/parser/parser_test.go:5233 |
| 38 | SELECT | 1 | SELECT A->>N FROM T | `SELECT a->>3 FROM t` | corpus/tidb/pkg/parser/parser_test.go:2302 |
| 39 | SELECT | 1 | SELECT A->N FROM T | `SELECT a->3 FROM t` | corpus/tidb/pkg/parser/parser_test.go:2301 |
| 40 | SELECT | 1 | SELECT AVG(), AVG(CN,CN) FROM | `select avg(), avg(c1,c2) from t;` | corpus/tidb/pkg/parser/parser_test.go:2170 |
| 41 | SELECT | 1 | SELECT BIT_AND(), BIT_AND(DISTINCT CN) | `select bit_and(), bit_and(distinct c1) from t;` | corpus/tidb/pkg/parser/parser_test.go:2182 |
| 42 | SELECT | 1 | SELECT BIT_AND(DISTINCT CN) FROM | `select bit_and(distinct c1) from t;` | corpus/tidb/pkg/parser/parser_test.go:2178 |
| 43 | SELECT | 1 | SELECT BIT_OR(), BIT_OR(DISTINCT CN) | `select bit_or(), bit_or(distinct c1) from t;` | corpus/tidb/pkg/parser/parser_test.go:2191 |
| 44 | SELECT | 1 | SELECT BIT_OR(DISTINCT CN) FROM | `select bit_or(distinct c1) from t;` | corpus/tidb/pkg/parser/parser_test.go:2187 |
| 45 | SELECT | 1 | SELECT BIT_XOR(), BIT_XOR(DISTINCT CN) | `select bit_xor(), bit_xor(distinct c1) from t;` | corpus/tidb/pkg/parser/parser_test.go:2199 |
| 46 | SELECT | 1 | SELECT BIT_XOR(DISTINCT CN) FROM | `select bit_xor(distinct c1) from t;` | corpus/tidb/pkg/parser/parser_test.go:2196 |
| 47 | SELECT | 1 | SELECT CAST(N AS FLOAT(N)); | `select cast(1 as float(54));` | corpus/tidb/pkg/parser/parser_test.go:1753 |
| 48 | SELECT | 1 | SELECT COUNT(CN, CN) FROM | `select count(c1, c2) from t;` | corpus/tidb/pkg/parser/parser_test.go:2226 |
| 49 | SELECT | 1 | SELECT COUNT(DISTINCT *) FROM | `select count(distinct *) from t;` | corpus/tidb/pkg/parser/parser_test.go:2221 |
| 50 | SELECT | 1 | SELECT CURRENT_DATE, CURRENT_DATE(), CURDATE(N) | `SELECT CURRENT_DATE, CURRENT_DATE(), CURDATE(1)` | corpus/tidb/pkg/parser/parser_test.go:1832 |
| 51 | SELECT | 1 | SELECT CURRENT_TIME('N') | `select current_time('1')` | corpus/tidb/pkg/parser/parser_test.go:1797 |
| 52 | SELECT | 1 | SELECT CURRENT_TIME(-N) | `select current_time(-1)` | corpus/tidb/pkg/parser/parser_test.go:1795 |
| 53 | SELECT | 1 | SELECT CURRENT_TIME(N.N) | `select current_time(1.0)` | corpus/tidb/pkg/parser/parser_test.go:1796 |
| 54 | SELECT | 1 | SELECT CURRENT_TIME(NULL) | `select current_time(null)` | corpus/tidb/pkg/parser/parser_test.go:1798 |
| 55 | SELECT | 1 | SELECT CURRENT_TIMESTAMP('N') | `select current_timestamp('2')` | corpus/tidb/pkg/parser/parser_test.go:1779 |
| 56 | SELECT | 1 | SELECT CURRENT_TIMESTAMP(-N) | `select current_timestamp(-1)` | corpus/tidb/pkg/parser/parser_test.go:1777 |
| 57 | SELECT | 1 | SELECT CURRENT_TIMESTAMP(N.N) | `select current_timestamp(1.0)` | corpus/tidb/pkg/parser/parser_test.go:1778 |
| 58 | SELECT | 1 | SELECT CURRENT_TIMESTAMP(NULL) | `select current_timestamp(null)` | corpus/tidb/pkg/parser/parser_test.go:1776 |
| 59 | SELECT | 1 | SELECT CURTIME('N') | `select curtime('1')` | corpus/tidb/pkg/parser/parser_test.go:1803 |
| 60 | SELECT | 1 | SELECT CURTIME(-N) | `select curtime(-1)` | corpus/tidb/pkg/parser/parser_test.go:1801 |
| 61 | SELECT | 1 | SELECT CURTIME(N.N) | `select curtime(1.0)` | corpus/tidb/pkg/parser/parser_test.go:1802 |
| 62 | SELECT | 1 | SELECT CURTIME(NULL) | `select curtime(null)` | corpus/tidb/pkg/parser/parser_test.go:1804 |
| 63 | SELECT | 1 | SELECT DATE_ADD("N-N-N N:N:N.N", "N,N") | `select date_add("2011-11-11 10:10:10.123456", "11,11")` | corpus/tidb/pkg/parser/parser_test.go:2068 |
| 64 | SELECT | 1 | SELECT DATE_ADD("N-N-N N:N:N.N", N) | `select date_add("2011-11-11 10:10:10.123456", 10)` | corpus/tidb/pkg/parser/parser_test.go:2066 |
| 65 | SELECT | 1 | SELECT DATE_ADD("N-N-N N:N:N.N", N.N) | `select date_add("2011-11-11 10:10:10.123456", 0.10)` | corpus/tidb/pkg/parser/parser_test.go:2067 |
| 66 | SELECT | 1 | SELECT DATE_SUB("N-N-N N:N:N.N", "N,N") | `select date_sub("2011-11-11 10:10:10.123456", "11,11")` | corpus/tidb/pkg/parser/parser_test.go:2133 |
| 67 | SELECT | 1 | SELECT DATE_SUB("N-N-N N:N:N.N", N) | `select date_sub("2011-11-11 10:10:10.123456", 10)` | corpus/tidb/pkg/parser/parser_test.go:2131 |
| 68 | SELECT | 1 | SELECT DATE_SUB("N-N-N N:N:N.N", N.N) | `select date_sub("2011-11-11 10:10:10.123456", 0.10)` | corpus/tidb/pkg/parser/parser_test.go:2132 |
| 69 | SELECT | 1 | SELECT DISTINCT ALL * | `SELECT DISTINCT ALL * FROM t` | corpus/tidb/pkg/parser/parser_test.go:587 |
| 70 | SELECT | 1 | SELECT DISTINCTROW ALL * | `SELECT DISTINCTROW ALL * FROM t` | corpus/tidb/pkg/parser/parser_test.go:588 |
| 71 | SELECT | 1 | SELECT JSON_ARRAYAGG(CN, CN) FROM | `select json_arrayagg(c1, c2) from t group by c1` | corpus/tidb/pkg/parser/parser_test.go:2256 |
| 72 | SELECT | 1 | SELECT JSON_ARRAYAGG(DISTINCT CN) FROM | `select json_arrayagg(distinct c2) from t group by c1` | corpus/tidb/pkg/parser/parser_test.go:2257 |
| 73 | SELECT | 1 | SELECT JSON_OBJECTAGG(CN, CN, CN) | `select json_objectagg(c1, c2, c3) from t group by c1` | corpus/tidb/pkg/parser/parser_test.go:2260 |
| 74 | SELECT | 1 | SELECT JSON_OBJECTAGG(DISTINCT CN, CN) | `select json_objectagg(distinct c1, c2) from t group by c1` | corpus/tidb/pkg/parser/parser_test.go:2261 |
| 75 | SELECT | 1 | SELECT MAX(CN,CN) FROM T; | `select max(c1,c2) from t;` | corpus/tidb/pkg/parser/parser_test.go:2202 |
| 76 | SELECT | 1 | SELECT MIN(CN,CN) FROM T; | `select min(c1,c2) from t;` | corpus/tidb/pkg/parser/parser_test.go:2208 |
| 77 | SELECT | 1 | SELECT N MEMBER OF | `SELECT 1 member of (1+1)` | corpus/tidb/pkg/parser/parser_test.go:2308 |
| 78 | SELECT | 1 | SELECT NTH_VALUE(VAL) OVER W | `SELECT NTH_VALUE(val) OVER w FROM t;` | corpus/tidb/pkg/parser/parser_test.go:6548 |
| 79 | SELECT | 1 | SELECT NXN | `select 0X11` | corpus/tidb/pkg/parser/parser_test.go:4968 |
| 80 | SELECT | 1 | SELECT ROW(N) | `select row(1)` | corpus/tidb/pkg/parser/parser_test.go:1724 |
| 81 | SELECT | 1 | SELECT SQL_NO_CACHE * FROM | `select SQL_NO_CACHE * from t` | corpus/tidb/pkg/parser/parser_test.go:5554 |
| 82 | SELECT | 1 | SELECT STD(CN, CN) FROM | `select std(c1, c2) from t` | corpus/tidb/pkg/parser/parser_test.go:2242 |
| 83 | SELECT | 1 | SELECT STDDEV(CN, CN) FROM | `select stddev(c1, c2) from t` | corpus/tidb/pkg/parser/parser_test.go:2244 |
| 84 | SELECT | 1 | SELECT STDDEV_POP(CN, CN) FROM | `select stddev_pop(c1, c2) from t` | corpus/tidb/pkg/parser/parser_test.go:2246 |
| 85 | SELECT | 1 | SELECT STDDEV_SAMP(CN, CN) FROM | `select stddev_samp(c1, c2) from t` | corpus/tidb/pkg/parser/parser_test.go:2248 |
| 86 | SELECT | 1 | SELECT SUM(CN,CN) FROM T; | `select sum(c1,c2) from t;` | corpus/tidb/pkg/parser/parser_test.go:2214 |
| 87 | SELECT | 1 | SELECT TIMESTAMPADD(SQL_TSI_MICROSECOND,N,'N-N-N'); | `SELECT TIMESTAMPADD(SQL_TSI_MICROSECOND,1,'2003-01-02');` | corpus/tidb/pkg/parser/parser_test.go:1914 |
| 88 | SELECT | 1 | SELECT UTC_TIME('N') | `select utc_time('1')` | corpus/tidb/pkg/parser/parser_test.go:1821 |
| 89 | SELECT | 1 | SELECT UTC_TIME(-N) | `select utc_time(-1)` | corpus/tidb/pkg/parser/parser_test.go:1819 |
| 90 | SELECT | 1 | SELECT UTC_TIME(N.N) | `select utc_time(1.0)` | corpus/tidb/pkg/parser/parser_test.go:1820 |
| 91 | SELECT | 1 | SELECT UTC_TIME(NULL) | `select utc_time(null)` | corpus/tidb/pkg/parser/parser_test.go:1822 |
| 92 | SELECT | 1 | SELECT UTC_TIMESTAMP('N') | `select utc_timestamp('1')` | corpus/tidb/pkg/parser/parser_test.go:1812 |
| 93 | SELECT | 1 | SELECT UTC_TIMESTAMP(-N) | `select utc_timestamp(-1)` | corpus/tidb/pkg/parser/parser_test.go:1810 |
| 94 | SELECT | 1 | SELECT UTC_TIMESTAMP(N.N) | `select utc_timestamp(1.0)` | corpus/tidb/pkg/parser/parser_test.go:1811 |
| 95 | SELECT | 1 | SELECT UTC_TIMESTAMP(NULL) | `select utc_timestamp(null)` | corpus/tidb/pkg/parser/parser_test.go:1813 |
| 96 | SELECT | 1 | SELECT VARIANCE(CN, CN) FROM | `select variance(c1, c2) from t` | corpus/tidb/pkg/parser/parser_test.go:2250 |
| 97 | SELECT | 1 | SELECT VAR_POP(CN, CN) FROM | `select var_pop(c1, c2) from t` | corpus/tidb/pkg/parser/parser_test.go:2252 |
| 98 | SELECT | 1 | SELECT VAR_SAMP(CN, CN) FROM | `select var_samp(c1, c2) from t` | corpus/tidb/pkg/parser/parser_test.go:2254 |
| 99 | SELECT | 1 | SELECT X'NXAA' | `select x'0xaa'` | corpus/tidb/pkg/parser/parser_test.go:4967 |
| 100 | TXN | 1 | COMMIT AND CHAIN RELEASE | `COMMIT AND CHAIN RELEASE` | corpus/tidb/pkg/parser/parser_test.go:626 |
| 101 | TXN | 1 | ROLLBACK AND CHAIN RELEASE | `ROLLBACK AND CHAIN RELEASE` | corpus/tidb/pkg/parser/parser_test.go:633 |

