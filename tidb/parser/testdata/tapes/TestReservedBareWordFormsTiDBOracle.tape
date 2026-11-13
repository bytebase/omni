-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestReservedBareWordFormsTiDBOracle default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| EXECUTE TIDB_CURRENT_TSO
< error
| 1064 42000 …line 1 column 24 near "TIDB_CURRENT_TSO" 
> exec ^
: string
| EXECUTE tidb_current_tso
< error
| 1064 42000 …line 1 column 24 near "tidb_current_tso" 
> exec ^
: string
| EXECUTE TIDB_CURRENT_TSO USING @a
< error
| 1064 42000 …line 1 column 24 near "TIDB_CURRENT_TSO USING @a" 
> exec ^
: string
| PREPARE TIDB_CURRENT_TSO FROM 'select 1'
< error
| 1064 42000 …line 1 column 24 near "TIDB_CURRENT_TSO FROM 'select 1'" 
> exec ^
: string
| DEALLOCATE PREPARE TIDB_CURRENT_TSO
< error
| 1064 42000 …line 1 column 35 near "TIDB_CURRENT_TSO" 
> exec ^
: string
| DROP PREPARE TIDB_CURRENT_TSO
< error
| 1064 42000 …line 1 column 29 near "TIDB_CURRENT_TSO" 
> exec ^
: string
| EXECUTE `tidb_current_tso`
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| EXECUTE `tidb_current_tso`
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXECUTE `tidb_current_tso` USING @a
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| EXECUTE `tidb_current_tso` USING @a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| PREPARE `tidb_current_tso` FROM 'select 1'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| PREPARE `tidb_current_tso` FROM 'select 1'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXECUTE s1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| EXECUTE s1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXECUTE s1 USING @tidb_current_tso
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| EXECUTE s1 USING @tidb_current_tso
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE tidb_current_tso (a int)
< error
| 1064 42000 …line 1 column 29 near "tidb_current_tso (a int)" 
> exec ^
: string
| CREATE TABLE t9 (tidb_current_tso int)
< error
| 1064 42000 …line 1 column 33 near "tidb_current_tso int)" 
> exec ^
: string
| CREATE DATABASE tidb_current_tso
< error
| 1064 42000 …line 1 column 32 near "tidb_current_tso" 
> exec ^
: string
| CREATE VIEW tidb_current_tso AS SELECT 1
< error
| 1064 42000 …line 1 column 28 near "tidb_current_tso AS SELECT 1" 
> exec ^
: string
| CREATE SEQUENCE tidb_current_tso
< error
| 1064 42000 …line 1 column 32 near "tidb_current_tso" 
> exec ^
: string
| CREATE INDEX tidb_current_tso ON t (a)
< error
| 1064 42000 …line 1 column 29 near "tidb_current_tso ON t (a)" 
> exec ^
: string
| CREATE INDEX i9 ON t (tidb_current_tso)
< error
| 1064 42000 …line 1 column 38 near "tidb_current_tso)" 
> exec ^
: string
| CREATE USER tidb_current_tso
< error
| 1064 42000 …line 1 column 28 near "tidb_current_tso" 
> exec ^
: string
| DROP TABLE tidb_current_tso
< error
| 1064 42000 …line 1 column 27 near "tidb_current_tso" 
> exec ^
: string
| TRUNCATE TABLE tidb_current_tso
< error
| 1064 42000 …line 1 column 31 near "tidb_current_tso" 
> exec ^
: string
| ANALYZE TABLE tidb_current_tso
< error
| 1064 42000 …line 1 column 30 near "tidb_current_tso" 
> exec ^
: string
| SHOW CREATE TABLE tidb_current_tso
< error
| 1064 42000 …line 1 column 34 near "tidb_current_tso" 
> exec ^
: string
| DESC tidb_current_tso
< error
| 1064 42000 …line 1 column 21 near "tidb_current_tso" 
> exec ^
: string
| LOCK TABLES tidb_current_tso READ
< error
| 1064 42000 …line 1 column 28 near "tidb_current_tso READ" 
> exec ^
: string
| ALTER TABLE t RENAME TO tidb_current_tso
< error
| 1064 42000 …line 1 column 40 near "tidb_current_tso" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN tidb_current_tso INT
< error
| 1064 42000 …line 1 column 41 near "tidb_current_tso INT" 
> exec ^
: string
| ALTER TABLE t DROP COLUMN tidb_current_tso
< error
| 1064 42000 …line 1 column 42 near "tidb_current_tso" 
> exec ^
: string
| ALTER TABLE t ADD INDEX tidb_current_tso (a)
< error
| 1064 42000 …line 1 column 40 near "tidb_current_tso (a)" 
> exec ^
: string
| SELECT * FROM tidb_current_tso
< error
| 1064 42000 …line 1 column 30 near "tidb_current_tso" 
> exec ^
: string
| SELECT 1 FROM tidb_current_tso t
< error
| 1064 42000 …line 1 column 30 near "tidb_current_tso t" 
> exec ^
: string
| SELECT a FROM t tidb_current_tso
< error
| 1064 42000 …line 1 column 32 near "tidb_current_tso" 
> exec ^
: string
| SELECT a FROM t AS tidb_current_tso
< error
| 1064 42000 …line 1 column 35 near "tidb_current_tso" 
> exec ^
: string
| SELECT a FROM t JOIN t tidb_current_tso
< error
| 1064 42000 …line 1 column 39 near "tidb_current_tso" 
> exec ^
: string
| SELECT * FROM t PARTITION (tidb_current_tso)
< error
| 1064 42000 …line 1 column 43 near "tidb_current_tso)" 
> exec ^
: string
| SELECT * FROM t USE INDEX (tidb_current_tso)
< error
| 1064 42000 …line 1 column 43 near "tidb_current_tso)" 
> exec ^
: string
| SELECT 1 AS tidb_current_tso
< error
| 1064 42000 …line 1 column 28 near "tidb_current_tso" 
> exec ^
: string
| SELECT 1 tidb_current_tso
< error
| 1064 42000 …line 1 column 25 near "tidb_current_tso" 
> exec ^
: string
| SELECT a tidb_current_tso, b FROM t
< error
| 1064 42000 …line 1 column 25 near "tidb_current_tso, b FROM t" 
> exec ^
: string
| SELECT TIDB_CURRENT_TSO AS TIDB_CURRENT_TSO
< error
| 1064 42000 …line 1 column 43 near "TIDB_CURRENT_TSO" 
> exec ^
: string
| SELECT 1 FROM t WINDOW tidb_current_tso AS ()
< error
| 1064 42000 …line 1 column 39 near "tidb_current_tso AS ()" 
> exec ^
: string
| WITH tidb_current_tso AS (SELECT 1) SELECT 1
< error
| 1064 42000 …line 1 column 21 near "tidb_current_tso AS (SELECT 1) SELECT 1" 
> exec ^
: string
| INSERT INTO t (tidb_current_tso) VALUES (1)
< error
| 1064 42000 …line 1 column 31 near "tidb_current_tso) VALUES (1)" 
> exec ^
: string
| UPDATE t SET tidb_current_tso = 1
< error
| 1064 42000 …line 1 column 29 near "tidb_current_tso = 1" 
> exec ^
: string
| SET tidb_current_tso = 1
< error
| 1064 42000 …line 1 column 20 near "tidb_current_tso = 1" 
> exec ^
: string
| SET NAMES tidb_current_tso
< error
| 1064 42000 …line 1 column 26 near "tidb_current_tso" 
> exec ^
: string
| SAVEPOINT tidb_current_tso
< error
| 1064 42000 …line 1 column 26 near "tidb_current_tso" 
> exec ^
: string
| CALL tidb_current_tso()
< error
| 1064 42000 …line 1 column 21 near "tidb_current_tso()" 
> exec ^
: string
| CREATE TABLE `tidb_current_tso` (a int)
< ok 0 0
> exec ^
: string
| CREATE TABLE t9 (`tidb_current_tso` int)
< ok 0 0
> exec ^
: string
| SELECT a AS `tidb_current_tso` FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT a AS `tidb_current_tso` FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT tidb_current_tso.a FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT tidb_current_tso.a FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT tidb_current_tso. a FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT tidb_current_tso. a FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT t.tidb_current_tso FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT t.tidb_current_tso FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT t. tidb_current_tso FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT t. tidb_current_tso FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT tidb_current_tso.* FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT tidb_current_tso.* FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM tidb_current_tso.t
< error
| 1146 42S02 Table 'tidb_current_tso.t' doesn't exist
> exec ^
: string
| SELECT * FROM tidb_current_tso.t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t.tidb_current_tso
< error
| 1146 42S02 Table 't.tidb_current_tso' doesn't exist
> exec ^
: string
| SELECT * FROM t.tidb_current_tso
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a.tidb_current_tso() FROM t
< error
| 1305 42000 FUNCTION a.tidb_current_tso does not exist
> exec ^
: string
| SELECT a.tidb_current_tso() FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT tidb_current_tso.f() FROM t
< error
| 1305 42000 FUNCTION tidb_current_tso.f does not exist
> exec ^
: string
| SELECT tidb_current_tso.f() FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SET @tidb_current_tso = 1
< ok 0 0
> exec ^
: string
| SELECT @tidb_current_tso
< ok 0 0
> exec ^
: string
| SELECT @@tidb_current_tso
< error
| 1193 HY000 Unknown system variable 'tidb_current_tso'
> exec ^
: string
| SELECT TIDB_CURRENT_TSO
< ok 0 0
> exec ^
: string
| SELECT TIDB_CURRENT_TSO()
< ok 0 0
> exec ^
: string
| SELECT tidb_current_tso ()
< ok 0 0
> exec ^
: string
| SELECT tidb_current_tso(), 1
< ok 0 0
> exec ^
: string
| SELECT tidb_current_tso FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT tidb_current_tso FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT tidb_current_tso + 1
< ok 0 0
> exec ^
: string
| SELECT tidb_current_tso AS x
< ok 0 0
> exec ^
: string
| SELECT CAST(tidb_current_tso AS CHAR)
< ok 0 0
> exec ^
: string
| SELECT COUNT(tidb_current_tso) FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT COUNT(tidb_current_tso) FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t WHERE a = tidb_current_tso
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t WHERE a = tidb_current_tso
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t ORDER BY tidb_current_tso
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t ORDER BY tidb_current_tso
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| INSERT INTO t (a) VALUES (tidb_current_tso)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| INSERT INTO t (a) VALUES (tidb_current_tso)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| UPDATE t SET a = tidb_current_tso
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| UPDATE t SET a = tidb_current_tso
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SET @a = tidb_current_tso()
< ok 0 0
> exec ^
: string
| DO TIDB_CURRENT_TSO
< ok 0 0
