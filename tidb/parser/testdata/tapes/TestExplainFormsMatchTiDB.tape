-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestExplainFormsMatchTiDB default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| EXPLAIN FORMAT = TRADITIONAL SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = JSON SELECT 1
< error
| 1105 HY000 explain format 'json' is not supported now
> exec ^
: string
| EXPLAIN FORMAT = JSON SELECT 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN FORMAT = ROW SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = DOT SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = BRIEF SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = VERBOSE SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = TRUE_CARD_COST SELECT 1
< error
| 1105 HY000 'explain format=true_card_cost' cannot work without 'analyze', please use 'explain analyze format=true_card_cost'
> exec ^
: string
| EXPLAIN FORMAT = TRUE_CARD_COST SELECT 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN FORMAT = TIDB_JSON SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = tidb_json SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN ANALYZE FORMAT = VERBOSE SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN ANALYZE FORMAT = TIDB_JSON SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = 'row' SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = 'ROW' SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = "row" SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = 'dot' SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = 'brief' SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = 'verbose' SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = 'hint' SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = 'nosuchformat' SELECT 1
< error
| 1791 HY000 Unknown EXPLAIN format name: 'nosuchformat'
> exec ^
: string
| EXPLAIN FORMAT = 'nosuchformat' SELECT 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN FORMAT = '' SELECT 1
< error
| 1791 HY000 Unknown EXPLAIN format name: ''
> exec ^
: string
| EXPLAIN FORMAT = '' SELECT 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN ANALYZE FORMAT = 'verbose' SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN ANALYZE FORMAT = 'binary' SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN ANALYZE format='TRUE_CARD_COST' SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = tree SELECT 1
< error
| 1064 42000 …line 1 column 21 near "tree SELECT 1" 
> exec ^
: string
| EXPLAIN FORMAT = hint SELECT 1
< error
| 1064 42000 …line 1 column 21 near "hint SELECT 1" 
> exec ^
: string
| EXPLAIN FORMAT = binary SELECT 1
< error
| 1064 42000 …line 1 column 23 near "binary SELECT 1" 
> exec ^
: string
| EXPLAIN FORMAT = cost_trace SELECT 1
< error
| 1064 42000 …line 1 column 27 near "cost_trace SELECT 1" 
> exec ^
: string
| EXPLAIN FORMAT = plan_cache SELECT 1
< error
| 1064 42000 …line 1 column 27 near "plan_cache SELECT 1" 
> exec ^
: string
| EXPLAIN FORMAT = nosuchformat SELECT 1
< error
| 1064 42000 …line 1 column 29 near "nosuchformat SELECT 1" 
> exec ^
: string
| EXPLAIN FORMAT = `row` SELECT 1
< error
| 1064 42000 …line 1 column 22 near "`row` SELECT 1" 
> exec ^
: string
| EXPLAIN FORMAT = select SELECT 1
< error
| 1064 42000 …line 1 column 23 near "select SELECT 1" 
> exec ^
: string
| EXPLAIN FORMAT = 1 SELECT 1
< error
| 1064 42000 …line 1 column 18 near "1 SELECT 1" 
> exec ^
: string
| EXPLAIN ANALYZE FORMAT = tree SELECT 1
< error
| 1064 42000 …line 1 column 29 near "tree SELECT 1" 
> exec ^
: string
| EXPLAIN ANALYZE FORMAT = nosuchformat SELECT 1
< error
| 1064 42000 …line 1 column 37 near "nosuchformat SELECT 1" 
> exec ^
: string
| EXPLAIN FORMAT = SELECT 1
< error
| 1064 42000 …line 1 column 23 near "SELECT 1" 
> exec ^
: string
| EXPLAIN FORMAT SELECT 1
< error
| 1064 42000 …line 1 column 21 near "SELECT 1" 
> exec ^
: string
| EXPLAIN FORMAT 'row' SELECT 1
< error
| 1064 42000 …line 1 column 20 near "'row' SELECT 1" 
> exec ^
: string
| EXPLAIN FORMAT = 'row' FORMAT = 'json' SELECT 1
< error
| 1064 42000 …line 1 column 29 near "FORMAT = 'json' SELECT 1" 
> exec ^
: string
| EXPLAIN FORMAT = 'row' 'row' SELECT 1
< error
| 1064 42000 …line 1 column 28 near "'row' SELECT 1" 
> exec ^
: string
| EXPLAIN FORMAT = ROW ROW SELECT 1
< error
| 1064 42000 …line 1 column 24 near "ROW SELECT 1" 
> exec ^
: string
| EXPLAIN FOR CONNECTION 1
< error
| 1094 HY000 Unknown thread id: 1
> exec ^
: string
| EXPLAIN FOR CONNECTION 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN FOR connection 42
< error
| 1094 HY000 Unknown thread id: 42
> exec ^
: string
| EXPLAIN FOR connection 42
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN FOR CONNECTION 0
< error
| 1094 HY000 Unknown thread id: 0
> exec ^
: string
| EXPLAIN FOR CONNECTION 0
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN FOR CONNECTION 18446744073709551615
< error
| 1094 HY000 Unknown thread id: 18446744073709551615
> exec ^
: string
| EXPLAIN FOR CONNECTION 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN FORMAT = 'dot' FOR CONNECTION 1
< error
| 1094 HY000 Unknown thread id: 1
> exec ^
: string
| EXPLAIN FORMAT = 'dot' FOR CONNECTION 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN FORMAT = DOT FOR CONNECTION 1
< error
| 1094 HY000 Unknown thread id: 1
> exec ^
: string
| EXPLAIN FORMAT = DOT FOR CONNECTION 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN FORMAT = 'row' FOR connection 1
< error
| 1094 HY000 Unknown thread id: 1
> exec ^
: string
| EXPLAIN FORMAT = 'row' FOR connection 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN FOR CONNECTION
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| EXPLAIN FOR 1
< error
| 1064 42000 …line 1 column 13 near "1" 
> exec ^
: string
| EXPLAIN FOR CONNECTION 'x'
< error
| 1064 42000 …line 1 column 26 near "'x'" 
> exec ^
: string
| EXPLAIN FOR CONNECTION 1.5
< error
| 1064 42000 …line 1 column 26 near "1.5" 
> exec ^
: string
| EXPLAIN FOR CONNECTION @a
< error
| 1064 42000 …line 1 column 25 near "@a" 
> exec ^
: string
| EXPLAIN FOR CONNECTION -1
< error
| 1064 42000 …line 1 column 24 near "-1" 
> exec ^
: string
| EXPLAIN FOR CONNECTION 1+1
< error
| 1064 42000 …line 1 column 25 near "+1" 
> exec ^
: string
| EXPLAIN FOR CONNECTION (1)
< error
| 1064 42000 …line 1 column 24 near "(1)" 
> exec ^
: string
| EXPLAIN FOR CONNECTION a
< error
| 1064 42000 …line 1 column 24 near "a" 
> exec ^
: string
| EXPLAIN FORMAT = 'row' FOR CONNECTION
< error
| 1064 42000 …line 1 column 37 near "" 
> exec ^
: string
| EXPLAIN ANALYZE FOR CONNECTION 1
< error
| 1064 42000 …line 1 column 19 near "FOR CONNECTION 1" 
> exec ^
: string
| EXPLAIN ANALYZE FORMAT = 'row' FOR CONNECTION 1
< error
| 1064 42000 …line 1 column 34 near "FOR CONNECTION 1" 
> exec ^
: string
| EXPLAIN EXTENDED SELECT 1
< error
| 1064 42000 …line 1 column 23 near "SELECT 1" 
> exec ^
: string
| EXPLAIN PARTITIONS SELECT 1
< error
| 1064 42000 …line 1 column 25 near "SELECT 1" 
> exec ^
: string
| EXPLAIN ANALYZE EXTENDED SELECT 1
< error
| 1064 42000 …line 1 column 24 near "EXTENDED SELECT 1" 
> exec ^
: string
| EXPLAIN EXTENDED FORMAT = ROW SELECT 1
< error
| 1064 42000 …line 1 column 25 near "= ROW SELECT 1" 
> exec ^
: string
| EXPLAIN SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN ANALYZE SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN INSERT INTO t1 VALUES (1)
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| EXPLAIN INSERT INTO t1 VALUES (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN REPLACE INTO t1 VALUES (1)
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| EXPLAIN REPLACE INTO t1 VALUES (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN UPDATE t1 SET a = 1
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| EXPLAIN UPDATE t1 SET a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN DELETE FROM t1
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| EXPLAIN DELETE FROM t1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN TABLE t1
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| EXPLAIN TABLE t1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN VALUES ROW(1)
< error
| 1051 42S02 Unknown table ''
> exec ^
: string
| EXPLAIN VALUES ROW(1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN FORMAT = 'row' UPDATE t1 SET a = 1
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| EXPLAIN FORMAT = 'row' UPDATE t1 SET a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN FORMAT = 'row' DELETE FROM t1
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| EXPLAIN FORMAT = 'row' DELETE FROM t1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN FORMAT = 'row' INSERT INTO t1 VALUES (1)
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| EXPLAIN FORMAT = 'row' INSERT INTO t1 VALUES (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN FORMAT = 'row' REPLACE INTO t1 VALUES (1)
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| EXPLAIN FORMAT = 'row' REPLACE INTO t1 VALUES (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN FORMAT = 'row' TABLE t1
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| EXPLAIN FORMAT = 'row' TABLE t1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN FORMAT = 'row' VALUES ROW(1)
< error
| 1051 42S02 Unknown table ''
> exec ^
: string
| EXPLAIN FORMAT = 'row' VALUES ROW(1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN ANALYZE TABLE t1
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| EXPLAIN ANALYZE TABLE t1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN ANALYZE UPDATE t1 SET a = 1
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| EXPLAIN ANALYZE UPDATE t1 SET a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN ALTER TABLE t1 ADD INDEX (a)
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| EXPLAIN ALTER TABLE t1 ADD INDEX (a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN ALTER TABLE t1 ADD a varchar(255)
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| EXPLAIN ALTER TABLE t1 ADD a varchar(255)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN ALTER TABLE t1 MODIFY a bigint
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| EXPLAIN ALTER TABLE t1 MODIFY a bigint
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN ALTER TABLE t1 DROP INDEX a
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| EXPLAIN ALTER TABLE t1 DROP INDEX a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN ALTER TABLE t1 ADD INDEX (a), ADD INDEX (b)
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| EXPLAIN ALTER TABLE t1 ADD INDEX (a), ADD INDEX (b)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN FORMAT = 'brief' ALTER TABLE t1 ADD INDEX (a)
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| EXPLAIN FORMAT = 'brief' ALTER TABLE t1 ADD INDEX (a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN FORMAT = DOT ALTER TABLE t1 ADD INDEX (a)
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| EXPLAIN FORMAT = DOT ALTER TABLE t1 ADD INDEX (a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN ANALYZE ALTER TABLE t1 ADD INDEX (a)
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| EXPLAIN ANALYZE ALTER TABLE t1 ADD INDEX (a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN IMPORT INTO t FROM '/tmp/a.csv'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| EXPLAIN IMPORT INTO t FROM '/tmp/a.csv'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN ANALYZE IMPORT INTO t FROM '/tmp/a.csv'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| EXPLAIN ANALYZE IMPORT INTO t FROM '/tmp/a.csv'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN IMPORT INTO t FROM TABLE t1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| EXPLAIN IMPORT INTO t FROM TABLE t1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN ANALYZE FORMAT = 'brief' ALTER TABLE t1 ADD INDEX (a)
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| EXPLAIN ANALYZE FORMAT = 'brief' ALTER TABLE t1 ADD INDEX (a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN (SELECT 1)
< ok 0 0
> exec ^
: string
| EXPLAIN WITH c AS (SELECT 1) SELECT * FROM c
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = 'row' (SELECT 1)
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = 'row' WITH c AS (SELECT 1) SELECT * FROM c
< ok 0 0
> exec ^
: string
| explain format = 'row' select c1 from t1 union (select c2 from t2) limit 1, 1
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| explain format = 'row' select c1 from t1 union (select c2 from t2) limit 1, 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN ALTER DATABASE d CHARSET = utf8
< error
| 1064 42000 …line 1 column 22 near "DATABASE d CHARSET = utf8" 
> exec ^
: string
| EXPLAIN CREATE TABLE t2 (a int)
< error
| 1064 42000 …line 1 column 14 near "CREATE TABLE t2 (a int)" 
> exec ^
: string
| EXPLAIN CREATE INDEX i ON t1 (a)
< error
| 1064 42000 …line 1 column 14 near "CREATE INDEX i ON t1 (a)" 
> exec ^
: string
| EXPLAIN DROP TABLE t1
< error
| 1064 42000 …line 1 column 12 near "DROP TABLE t1" 
> exec ^
: string
| EXPLAIN TRUNCATE TABLE t1
< error
| 1064 42000 …line 1 column 22 near "TABLE t1" 
> exec ^
: string
| EXPLAIN RENAME TABLE t1 TO t2
< error
| 1064 42000 …line 1 column 14 near "RENAME TABLE t1 TO t2" 
> exec ^
: string
| EXPLAIN LOAD DATA INFILE 'x' INTO TABLE t1
< error
| 1064 42000 …line 1 column 12 near "LOAD DATA INFILE 'x' INTO TABLE t1" 
> exec ^
: string
| EXPLAIN SET @a = 1
< error
| 1064 42000 …line 1 column 11 near "SET @a = 1" 
> exec ^
: string
| EXPLAIN ANALYZE SET @a = 1
< error
| 1064 42000 …line 1 column 19 near "SET @a = 1" 
> exec ^
: string
| EXPLAIN FORMAT = 'row' SET @a = 1
< error
| 1064 42000 …line 1 column 26 near "SET @a = 1" 
> exec ^
: string
| EXPLAIN FORMAT = 'row' SHOW TABLES
< error
| 1064 42000 …line 1 column 27 near "SHOW TABLES" 
> exec ^
: string
| EXPLAIN FORMAT = ROW SHOW TABLES
< error
| 1064 42000 …line 1 column 25 near "SHOW TABLES" 
> exec ^
: string
| EXPLAIN t1
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| EXPLAIN t1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN t1 a
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| EXPLAIN t1 a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN t1 `a`
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| EXPLAIN t1 `a`
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN t1 a.b
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| EXPLAIN t1 a.b
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN d1.t1 a
< error
| 1146 42S02 Table 'd1.t1' doesn't exist
> exec ^
: string
| EXPLAIN d1.t1 a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DESC t1 a
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| DESC t1 a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DESC SCHE.TABL COLUM
< error
| 1146 42S02 Table 'sche.TABL' doesn't exist
> exec ^
: string
| DESC SCHE.TABL COLUM
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN BEGIN
< error
| 1146 42S02 Table 'test.BEGIN' doesn't exist
> exec ^
: string
| EXPLAIN BEGIN
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN EXTENDED t1
< error
| 1146 42S02 Table 'test.EXTENDED' doesn't exist
> exec ^
: string
| EXPLAIN EXTENDED t1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN PARTITIONS t1
< error
| 1146 42S02 Table 'test.PARTITIONS' doesn't exist
> exec ^
: string
| EXPLAIN PARTITIONS t1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN EXTENDED
< error
| 1146 42S02 Table 'test.EXTENDED' doesn't exist
> exec ^
: string
| EXPLAIN EXTENDED
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN t1 'abc'
< error
| 1064 42000 …line 1 column 16 near "'abc'" 
> exec ^
: string
| EXPLAIN t1 '%a%'
< error
| 1064 42000 …line 1 column 16 near "'%a%'" 
> exec ^
: string
| EXPLAIN t1 123
< error
| 1064 42000 …line 1 column 14 near "123" 
> exec ^
: string
| EXPLAIN t1 select
< error
| 1064 42000 …line 1 column 17 near "select" 
> exec ^
: string
| EXPLAIN t1 a b
< error
| 1064 42000 …line 1 column 14 near "b" 
> exec ^
: string
| DESC t1 'abc'
< error
| 1064 42000 …line 1 column 13 near "'abc'" 
> exec ^
: string
| DESCRIBE t1 'abc'
< error
| 1064 42000 …line 1 column 17 near "'abc'" 
> exec ^
: string
| explain explore 'digestxxx'
< error
| 1064 42000 …line 1 column 27 near "'digestxxx'" 
> exec ^
: string
| explain explore '23adc8e6f62'
< error
| 1064 42000 …line 1 column 29 near "'23adc8e6f62'" 
> exec ^
: string
| explain explore 'select a from t'
< error
| 1064 42000 …line 1 column 33 near "'select a from t'" 
> exec ^
: string
| EXPLAIN FORMAT = 'row' t1
< error
| 1064 42000 …line 1 column 25 near "t1" 
> exec ^
: string
| EXPLAIN FORMAT = 'row' t1 a
< error
| 1064 42000 …line 1 column 25 near "t1 a" 
> exec ^
: string
| EXPLAIN FORMAT = '' t1
< error
| 1064 42000 …line 1 column 22 near "t1" 
> exec ^
: string
| EXPLAIN FORMAT = '' t1 a
< error
| 1064 42000 …line 1 column 22 near "t1 a" 
> exec ^
: string
| EXPLAIN FORMAT = ROW t1
< error
| 1064 42000 …line 1 column 23 near "t1" 
> exec ^
: string
| EXPLAIN ANALYZE t1
< error
| 1064 42000 …line 1 column 18 near "t1" 
> exec ^
: string
| EXPLAIN ANALYZE t1 a
< error
| 1064 42000 …line 1 column 18 near "t1 a" 
> exec ^
: string
| EXPLAIN 'sqldigest'
< error
| 1064 42000 …line 1 column 19 near "'sqldigest'" 
> exec ^
: string
| EXPLAIN ANALYZE 'sqldigest'
< error
| 1064 42000 …line 1 column 27 near "'sqldigest'" 
> exec ^
: string
| EXPLAIN format='json' 'sqldigest'
< error
| 1064 42000 …line 1 column 33 near "'sqldigest'" 
