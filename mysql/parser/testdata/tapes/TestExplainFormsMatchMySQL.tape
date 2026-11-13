-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestExplainFormsMatchMySQL default#1
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| EXPLAIN FORMAT = TRADITIONAL SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = JSON SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = TREE SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = tree SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = 'traditional' SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = 'TRADITIONAL' SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = "traditional" SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = `traditional` SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = 'json' SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = 'row' SELECT 1
< error
| 1791 HY000 Unknown EXPLAIN format name: 'row'
> exec ^
: string
| EXPLAIN FORMAT = 'row' SELECT 1
| ; SELECT FROM
< error
| 1791 HY000 Unknown EXPLAIN format name: 'row'
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
| 1791 HY000 Unknown EXPLAIN format name: 'nosuchformat'
> exec ^
: string
| EXPLAIN FORMAT = nosuchformat SELECT 1
< error
| 1791 HY000 Unknown EXPLAIN format name: 'nosuchformat'
> exec ^
: string
| EXPLAIN FORMAT = nosuchformat SELECT 1
| ; SELECT FROM
< error
| 1791 HY000 Unknown EXPLAIN format name: 'nosuchformat'
> exec ^
: string
| EXPLAIN FORMAT = `row` SELECT 1
< error
| 1791 HY000 Unknown EXPLAIN format name: 'row'
> exec ^
: string
| EXPLAIN FORMAT = `row` SELECT 1
| ; SELECT FROM
< error
| 1791 HY000 Unknown EXPLAIN format name: 'row'
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
| 1791 HY000 Unknown EXPLAIN format name: ''
> exec ^
: string
| EXPLAIN ANALYZE FORMAT = TREE SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN ANALYZE FORMAT = 'tree' SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN ANALYZE FORMAT = 'traditional' SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = binary SELECT 1
< error
| 1064 42000 …near 'binary SELECT 1' at line 1
> exec ^
: string
| EXPLAIN FORMAT = select SELECT 1
< error
| 1064 42000 …near 'select SELECT 1' at line 1
> exec ^
: string
| EXPLAIN FORMAT = 1 SELECT 1
< error
| 1064 42000 …near '1 SELECT 1' at line 1
> exec ^
: string
| EXPLAIN FORMAT = SELECT 1
< error
| 1064 42000 …near 'SELECT 1' at line 1
> exec ^
: string
| EXPLAIN FORMAT SELECT 1
< error
| 1064 42000 …near 'SELECT 1' at line 1
> exec ^
: string
| EXPLAIN FORMAT 'row' SELECT 1
< error
| 1064 42000 …near 'SELECT 1' at line 1
> exec ^
: string
| EXPLAIN FORMAT = 'traditional' FORMAT = 'json' SELECT 1
< error
| 1064 42000 …near 'FORMAT = 'json' SELECT 1' at line 1
> exec ^
: string
| EXPLAIN FORMAT = 'traditional' 'traditional' SELECT 1
< error
| 1064 42000 …near ''traditional' SELECT 1' at line 1
> exec ^
: string
| EXPLAIN FORMAT = TRADITIONAL FORMAT = JSON SELECT 1
< error
| 1064 42000 …near 'FORMAT = JSON SELECT 1' at line 1
> exec ^
: string
| EXPLAIN FORMAT = 'traditional' SELECT 1 FORMAT = 'json'
< error
| 1064 42000 …near '= 'json'' at line 1
> exec ^
: string
| EXPLAIN FOR CONNECTION 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| EXPLAIN FOR CONNECTION 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| EXPLAIN FOR connection 42
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| EXPLAIN FOR connection 42
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| EXPLAIN FOR CONNECTION 0
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| EXPLAIN FOR CONNECTION 0
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| EXPLAIN FOR CONNECTION 18446744073709551615
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| EXPLAIN FOR CONNECTION 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| EXPLAIN FORMAT = 'traditional' FOR CONNECTION 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| EXPLAIN FORMAT = 'traditional' FOR CONNECTION 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| EXPLAIN FORMAT = TREE FOR CONNECTION 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| EXPLAIN FORMAT = TREE FOR CONNECTION 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| EXPLAIN ANALYZE FOR CONNECTION 1
< error
| 1235 42000 This version of MySQL doesn't yet support 'EXPLAIN ANALYZE FOR CONNECTION'
> exec ^
: string
| EXPLAIN ANALYZE FOR CONNECTION 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| EXPLAIN ANALYZE FORMAT = TREE FOR CONNECTION 1
< error
| 1235 42000 This version of MySQL doesn't yet support 'EXPLAIN ANALYZE FOR CONNECTION'
> exec ^
: string
| EXPLAIN ANALYZE FORMAT = TREE FOR CONNECTION 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| EXPLAIN FOR CONNECTION
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| EXPLAIN FOR 1
< error
| 1064 42000 …near '1' at line 1
> exec ^
: string
| EXPLAIN FOR CONNECTION 'x'
< error
| 1064 42000 …near ''x'' at line 1
> exec ^
: string
| EXPLAIN FOR CONNECTION 1.5
< error
| 1064 42000 Only integers allowed as number here near '1.5' at line 1
> exec ^
: string
| EXPLAIN FOR CONNECTION @a
< error
| 1064 42000 …near '@a' at line 1
> exec ^
: string
| EXPLAIN FOR CONNECTION -1
< error
| 1064 42000 …near '-1' at line 1
> exec ^
: string
| EXPLAIN FOR CONNECTION 1+1
< error
| 1064 42000 …near '+1' at line 1
> exec ^
: string
| EXPLAIN FOR CONNECTION (1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| EXPLAIN FOR CONNECTION a
< error
| 1064 42000 …near 'a' at line 1
> exec ^
: string
| EXPLAIN FOR CONNECTION t.a
< error
| 1064 42000 …near 't.a' at line 1
> exec ^
: string
| EXPLAIN FORMAT = TRADITIONAL FOR CONNECTION
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| EXPLAIN FOR DATABASE d SELECT 1
< error
| 1064 42000 …near 'DATABASE d SELECT 1' at line 1
> exec ^
: string
| EXPLAIN FOR SCHEMA d SELECT 1
< error
| 1064 42000 …near 'SCHEMA d SELECT 1' at line 1
> exec ^
: string
| EXPLAIN FORMAT = JSON INTO @v SELECT 1
< error
| 1064 42000 …near 'INTO @v SELECT 1' at line 1
> exec ^
: string
| EXPLAIN INTO @v SELECT 1
< error
| 1064 42000 …near 'INTO @v SELECT 1' at line 1
> exec ^
: string
| EXPLAIN EXTENDED SELECT 1
< error
| 1064 42000 …near 'SELECT 1' at line 1
> exec ^
: string
| EXPLAIN PARTITIONS SELECT 1
< error
| 1064 42000 …near 'SELECT 1' at line 1
> exec ^
: string
| EXPLAIN ANALYZE EXTENDED SELECT 1
< error
| 1064 42000 …near 'EXTENDED SELECT 1' at line 1
> exec ^
: string
| EXPLAIN EXTENDED FORMAT = JSON SELECT 1
< error
| 1064 42000 …near '= JSON SELECT 1' at line 1
> exec ^
: string
| EXPLAIN FORMAT = JSON EXTENDED SELECT 1
< error
| 1064 42000 …near 'EXTENDED SELECT 1' at line 1
> exec ^
: string
| DESC EXTENDED SELECT 1
< error
| 1064 42000 …near 'SELECT 1' at line 1
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
| 1146 42S02 Table 'omni_t_mysql_parser.t1' doesn't exist
> exec ^
: string
| EXPLAIN INSERT INTO t1 VALUES (1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| EXPLAIN REPLACE INTO t1 VALUES (1)
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t1' doesn't exist
> exec ^
: string
| EXPLAIN REPLACE INTO t1 VALUES (1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| EXPLAIN UPDATE t1 SET a = 1
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t1' doesn't exist
> exec ^
: string
| EXPLAIN UPDATE t1 SET a = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| EXPLAIN DELETE FROM t1
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t1' doesn't exist
> exec ^
: string
| EXPLAIN DELETE FROM t1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| EXPLAIN TABLE t1
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t1' doesn't exist
> exec ^
: string
| EXPLAIN TABLE t1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| EXPLAIN VALUES ROW(1)
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = 'traditional' UPDATE t1 SET a = 1
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t1' doesn't exist
> exec ^
: string
| EXPLAIN FORMAT = 'traditional' UPDATE t1 SET a = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| EXPLAIN FORMAT = 'traditional' DELETE FROM t1
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t1' doesn't exist
> exec ^
: string
| EXPLAIN FORMAT = 'traditional' DELETE FROM t1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| EXPLAIN FORMAT = 'traditional' INSERT INTO t1 VALUES (1)
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t1' doesn't exist
> exec ^
: string
| EXPLAIN FORMAT = 'traditional' INSERT INTO t1 VALUES (1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| EXPLAIN FORMAT = 'traditional' REPLACE INTO t1 VALUES (1)
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t1' doesn't exist
> exec ^
: string
| EXPLAIN FORMAT = 'traditional' REPLACE INTO t1 VALUES (1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| EXPLAIN FORMAT = 'traditional' TABLE t1
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t1' doesn't exist
> exec ^
: string
| EXPLAIN FORMAT = 'traditional' TABLE t1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| EXPLAIN FORMAT = 'traditional' VALUES ROW(1)
< ok 0 0
> exec ^
: string
| EXPLAIN ANALYZE TABLE t1
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t1' doesn't exist
> exec ^
: string
| EXPLAIN ANALYZE TABLE t1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| EXPLAIN ANALYZE VALUES ROW(1)
< ok 0 0
> exec ^
: string
| EXPLAIN ANALYZE UPDATE t1 SET a = 1
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t1' doesn't exist
> exec ^
: string
| EXPLAIN ANALYZE UPDATE t1 SET a = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| EXPLAIN FORMAT = 'traditional' (SELECT 1)
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = 'traditional' WITH c AS (SELECT 1) SELECT * FROM c
< ok 0 0
> exec ^
: string
| EXPLAIN ANALYZE (SELECT 1)
< ok 0 0
> exec ^
: string
| EXPLAIN ANALYZE WITH c AS (SELECT 1) SELECT * FROM c
< ok 0 0
> exec ^
: string
| EXPLAIN ALTER TABLE t1 ADD INDEX (a)
< error
| 1064 42000 …near 'ALTER TABLE t1 ADD INDEX (a)' at line 1
> exec ^
: string
| EXPLAIN ALTER TABLE t1 ADD a varchar(255)
< error
| 1064 42000 …near 'ALTER TABLE t1 ADD a varchar(255)' at line 1
> exec ^
: string
| EXPLAIN ANALYZE ALTER TABLE t1 ADD a int
< error
| 1064 42000 …near 'ALTER TABLE t1 ADD a int' at line 1
> exec ^
: string
| EXPLAIN FORMAT = 'traditional' ALTER TABLE t1 ADD INDEX (a)
< error
| 1064 42000 …near 'ALTER TABLE t1 ADD INDEX (a)' at line 1
> exec ^
: string
| EXPLAIN CREATE TABLE t2 (a int)
< error
| 1064 42000 …near 'CREATE TABLE t2 (a int)' at line 1
> exec ^
: string
| EXPLAIN CREATE INDEX i ON t1 (a)
< error
| 1064 42000 …near 'CREATE INDEX i ON t1 (a)' at line 1
> exec ^
: string
| EXPLAIN DROP TABLE t1
< error
| 1064 42000 …near 'DROP TABLE t1' at line 1
> exec ^
: string
| EXPLAIN TRUNCATE TABLE t1
< error
| 1064 42000 …near 'TABLE t1' at line 1
> exec ^
: string
| EXPLAIN RENAME TABLE t1 TO t2
< error
| 1064 42000 …near 'RENAME TABLE t1 TO t2' at line 1
> exec ^
: string
| EXPLAIN LOAD DATA INFILE 'x' INTO TABLE t1
< error
| 1064 42000 …near 'LOAD DATA INFILE 'x' INTO TABLE t1' at line 1
> exec ^
: string
| EXPLAIN SET @a = 1
< error
| 1064 42000 …near 'SET @a = 1' at line 1
> exec ^
: string
| EXPLAIN ANALYZE SET @a = 1
< error
| 1064 42000 …near 'SET @a = 1' at line 1
> exec ^
: string
| EXPLAIN FORMAT = 'traditional' SET @a = 1
< error
| 1064 42000 …near 'SET @a = 1' at line 1
> exec ^
: string
| EXPLAIN FORMAT = 'traditional' SHOW TABLES
< error
| 1064 42000 …near 'SHOW TABLES' at line 1
> exec ^
: string
| EXPLAIN t1
< ok 0 0
> exec ^
: string
| EXPLAIN t1 a
< ok 0 0
> exec ^
: string
| EXPLAIN t1 `a`
< ok 0 0
> exec ^
: string
| EXPLAIN t1 'abc'
< ok 0 0
> exec ^
: string
| EXPLAIN t1 '%a%'
< ok 0 0
> exec ^
: string
| EXPLAIN d1.t1 a
< error
| 1142 42000 SELECT command denied to user 'omni_test'@'172.17.0.1' for table 't1'
> exec ^
: string
| EXPLAIN d1.t1 a
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DESC t1 a
< ok 0 0
> exec ^
: string
| DESC t1 'abc'
< ok 0 0
> exec ^
: string
| DESCRIBE t1 'abc'
< ok 0 0
> exec ^
: string
| EXPLAIN BEGIN
< ok 0 0
> exec ^
: string
| EXPLAIN EXTENDED t1
< ok 0 0
> exec ^
: string
| EXPLAIN PARTITIONS t1
< ok 0 0
> exec ^
: string
| EXPLAIN EXTENDED
< ok 0 0
> exec ^
: string
| EXPLAIN t1 a.b
< error
| 1064 42000 …near '.b' at line 1
> exec ^
: string
| EXPLAIN t1 123
< error
| 1064 42000 …near '123' at line 1
> exec ^
: string
| EXPLAIN t1 select
< error
| 1064 42000 …near 'select' at line 1
> exec ^
: string
| EXPLAIN t1 a b
< error
| 1064 42000 …near 'b' at line 1
> exec ^
: string
| EXPLAIN FORMAT = 'traditional' t1
< error
| 1064 42000 …near 't1' at line 1
> exec ^
: string
| EXPLAIN FORMAT = 'traditional' t1 a
< error
| 1064 42000 …near 't1 a' at line 1
> exec ^
: string
| EXPLAIN FORMAT = TRADITIONAL t1
< error
| 1064 42000 …near 't1' at line 1
> exec ^
: string
| EXPLAIN ANALYZE t1
< error
| 1064 42000 …near 't1' at line 1
> exec ^
: string
| EXPLAIN ANALYZE t1 a
< error
| 1064 42000 …near 't1 a' at line 1
> exec ^
: string
| EXPLAIN 'sqldigest'
< error
| 1064 42000 …near ''sqldigest'' at line 1
> exec ^
: string
| EXPLAIN ANALYZE 'sqldigest'
< error
| 1064 42000 …near ''sqldigest'' at line 1
