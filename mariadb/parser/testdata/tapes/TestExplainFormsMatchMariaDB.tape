-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestExplainFormsMatchMariaDB default#1
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| EXPLAIN FORMAT = JSON SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = TRADITIONAL SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = TREE SELECT 1
< error
| 1791 HY000 Unknown EXPLAIN/ANALYZE format name: 'TREE'
> exec ^
: string
| EXPLAIN FORMAT = TREE SELECT 1
| ; SELECT FROM
< error
| 1791 HY000 Unknown EXPLAIN/ANALYZE format name: 'TREE'
> exec ^
: string
| EXPLAIN FORMAT = 'json' SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = 'JSON' SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = "json" SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = `json` SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = 'traditional' SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = 'nosuchformat' SELECT 1
< error
| 1791 HY000 Unknown EXPLAIN/ANALYZE format name: 'nosuchformat'
> exec ^
: string
| EXPLAIN FORMAT = 'nosuchformat' SELECT 1
| ; SELECT FROM
< error
| 1791 HY000 Unknown EXPLAIN/ANALYZE format name: 'nosuchformat'
> exec ^
: string
| EXPLAIN FORMAT = nosuchformat SELECT 1
< error
| 1791 HY000 Unknown EXPLAIN/ANALYZE format name: 'nosuchformat'
> exec ^
: string
| EXPLAIN FORMAT = nosuchformat SELECT 1
| ; SELECT FROM
< error
| 1791 HY000 Unknown EXPLAIN/ANALYZE format name: 'nosuchformat'
> exec ^
: string
| EXPLAIN FORMAT = '' SELECT 1
< error
| 1791 HY000 Unknown EXPLAIN/ANALYZE format name: ''
> exec ^
: string
| EXPLAIN FORMAT = '' SELECT 1
| ; SELECT FROM
< error
| 1791 HY000 Unknown EXPLAIN/ANALYZE format name: ''
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
| EXPLAIN FORMAT 'json' SELECT 1
< error
| 1064 42000 …near 'SELECT 1' at line 1
> exec ^
: string
| EXPLAIN FORMAT = 'json' FORMAT = 'traditional' SELECT 1
< error
| 1064 42000 …near 'FORMAT = 'traditional' SELECT 1' at line 1
> exec ^
: string
| EXPLAIN FORMAT = 'traditional' 'traditional' SELECT 1
< error
| 1064 42000 …near ''traditional' SELECT 1' at line 1
> exec ^
: string
| EXPLAIN FORMAT = JSON FORMAT = TRADITIONAL SELECT 1
< error
| 1064 42000 …near 'FORMAT = TRADITIONAL SELECT 1' at line 1
> exec ^
: string
| EXPLAIN EXTENDED SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN PARTITIONS SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN EXTENDED SELECT * FROM t1
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| EXPLAIN EXTENDED SELECT * FROM t1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| EXPLAIN PARTITIONS SELECT * FROM t1
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| EXPLAIN PARTITIONS SELECT * FROM t1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| EXPLAIN EXTENDED FOR CONNECTION 1
< error
| 1064 42000 …near 'FOR CONNECTION 1' at line 1
> exec ^
: string
| EXPLAIN PARTITIONS FOR CONNECTION 1
< error
| 1064 42000 …near 'FOR CONNECTION 1' at line 1
> exec ^
: string
| EXPLAIN ANALYZE SELECT 1
< error
| 1064 42000 …near 'ANALYZE SELECT 1' at line 1
> exec ^
: string
| EXPLAIN ANALYZE FORMAT = JSON SELECT 1
< error
| 1064 42000 …near 'ANALYZE FORMAT = JSON SELECT 1' at line 1
> exec ^
: string
| EXPLAIN ANALYZE FORMAT = 'json' SELECT 1
< error
| 1064 42000 …near 'ANALYZE FORMAT = 'json' SELECT 1' at line 1
> exec ^
: string
| EXPLAIN ANALYZE FOR CONNECTION 1
< error
| 1064 42000 …near 'ANALYZE FOR CONNECTION 1' at line 1
> exec ^
: string
| EXPLAIN ANALYZE t1
< error
| 1064 42000 …near 'ANALYZE t1' at line 1
> exec ^
: string
| EXPLAIN ANALYZE UPDATE t1 SET a = 1
< error
| 1064 42000 …near 'ANALYZE UPDATE t1 SET a = 1' at line 1
> exec ^
: string
| DESC ANALYZE SELECT 1
< error
| 1064 42000 …near 'ANALYZE SELECT 1' at line 1
> exec ^
: string
| EXPLAIN FOR CONNECTION 1
< ok 0 0
> exec ^
: string
| EXPLAIN FOR connection 42
< ok 0 0
> exec ^
: string
| EXPLAIN FOR CONNECTION 0
< ok 0 0
> exec ^
: string
| EXPLAIN FOR CONNECTION 18446744073709551615
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = 'traditional' FOR CONNECTION 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = JSON FOR CONNECTION 1
< ok 0 0
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
| EXPLAIN FORMAT = JSON FOR CONNECTION
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| EXPLAIN SELECT 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| EXPLAIN VALUES (1)
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = JSON UPDATE t1 SET a = 1
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| EXPLAIN FORMAT = JSON UPDATE t1 SET a = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| EXPLAIN FORMAT = JSON DELETE FROM t1
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| EXPLAIN FORMAT = JSON DELETE FROM t1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| EXPLAIN FORMAT = JSON INSERT INTO t1 VALUES (1)
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| EXPLAIN FORMAT = JSON INSERT INTO t1 VALUES (1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| EXPLAIN FORMAT = JSON REPLACE INTO t1 VALUES (1)
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| EXPLAIN FORMAT = JSON REPLACE INTO t1 VALUES (1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| EXPLAIN FORMAT = JSON VALUES (1)
< ok 0 0
> exec ^
: string
| EXPLAIN EXTENDED UPDATE t1 SET a = 1
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| EXPLAIN EXTENDED UPDATE t1 SET a = 1
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
| EXPLAIN FORMAT = JSON (SELECT 1)
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = JSON WITH c AS (SELECT 1) SELECT * FROM c
< ok 0 0
> exec ^
: string
| EXPLAIN TABLE t1
< error
| 1064 42000 …near 'TABLE t1' at line 1
> exec ^
: string
| EXPLAIN VALUES ROW(1)
< error
| 1064 42000 …near 'ROW(1)' at line 1
> exec ^
: string
| EXPLAIN FORMAT = JSON TABLE t1
< error
| 1064 42000 …near 'TABLE t1' at line 1
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
| EXPLAIN FORMAT = JSON ALTER TABLE t1 ADD INDEX (a)
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
| EXPLAIN SET @a = 1
< error
| 1064 42000 …near 'SET @a = 1' at line 1
> exec ^
: string
| EXPLAIN FORMAT = JSON SET @a = 1
< error
| 1064 42000 …near 'SET @a = 1' at line 1
> exec ^
: string
| EXPLAIN FORMAT = JSON SHOW TABLES
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
< ok 0 0
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
| EXPLAIN FORMAT = JSON t1
< error
| 1064 42000 …near 't1' at line 1
> exec ^
: string
| EXPLAIN FORMAT = 'json' t1
< error
| 1064 42000 …near 't1' at line 1
> exec ^
: string
| EXPLAIN FORMAT = JSON t1 a
< error
| 1064 42000 …near 't1 a' at line 1
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
| EXPLAIN PARTITIONS `a`
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = JSON t1
< error
| 1064 42000 …near 't1' at line 1
> exec ^
: string
| EXPLAIN 'sqldigest'
< error
| 1064 42000 …near ''sqldigest'' at line 1
