-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestRequiredFormsOracle default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| CREATE TABLE tl (a key)
< error
| 1064 42000 …line 1 column 22 near "key)" 
> exec ^
: string
| CREATE TABLE tl (a primary key)
< error
| 1064 42000 …line 1 column 26 near "primary key)" 
> exec ^
: string
| CREATE TABLE tl (a unique)
< error
| 1064 42000 …line 1 column 25 near "unique)" 
> exec ^
: string
| CREATE TABLE tl (a unique key)
< error
| 1064 42000 …line 1 column 25 near "unique key)" 
> exec ^
: string
| CREATE TABLE tl (a comment 'x')
< error
| 1064 42000 …line 1 column 26 near "comment 'x')" 
> exec ^
: string
| CREATE TABLE tl (a not null)
< error
| 1064 42000 …line 1 column 22 near "not null)" 
> exec ^
: string
| CREATE TABLE tl (a null)
< error
| 1064 42000 …line 1 column 23 near "null)" 
> exec ^
: string
| CREATE TABLE tl (a default 1)
< error
| 1064 42000 …line 1 column 26 near "default 1)" 
> exec ^
: string
| CREATE TABLE tl (a auto_increment)
< error
| 1064 42000 …line 1 column 33 near "auto_increment)" 
> exec ^
: string
| CREATE TABLE tl (a collate utf8mb4_bin)
< error
| 1064 42000 …line 1 column 26 near "collate utf8mb4_bin)" 
> exec ^
: string
| CREATE TABLE tl (a references r (id))
< error
| 1064 42000 …line 1 column 29 near "references r (id))" 
> exec ^
: string
| CREATE TABLE tl (a check (a > 0))
< error
| 1064 42000 …line 1 column 24 near "check (a > 0))" 
> exec ^
: string
| CREATE TABLE tl (a generated always as (1))
< error
| 1064 42000 …line 1 column 28 near "generated always as (1))" 
> exec ^
: string
| CREATE TABLE tl (a generated always as (1) stored)
< error
| 1064 42000 …line 1 column 28 near "generated always as (1) stored)" 
> exec ^
: string
| CREATE TABLE tl (a generated always as (1) virtual)
< error
| 1064 42000 …line 1 column 28 near "generated always as (1) virtual)" 
> exec ^
: string
| CREATE TABLE tl (a as (1))
< error
| 1064 42000 …line 1 column 21 near "as (1))" 
> exec ^
: string
| CREATE TABLE tl (a as (1) stored)
< error
| 1064 42000 …line 1 column 21 near "as (1) stored)" 
> exec ^
: string
| CREATE TABLE tl (a on update now())
< error
| 1064 42000 …line 1 column 21 near "on update now())" 
> exec ^
: string
| CREATE TABLE tl (a column_format fixed)
< error
| 1064 42000 …line 1 column 32 near "column_format fixed)" 
> exec ^
: string
| CREATE TABLE tl (a storage disk)
< error
| 1064 42000 …line 1 column 26 near "storage disk)" 
> exec ^
: string
| CREATE TABLE tl (a srid 0)
< error
| 1064 42000 …line 1 column 23 near "srid 0)" 
> exec ^
: string
| CREATE TABLE tl (a invisible)
< error
| 1064 42000 …line 1 column 28 near "invisible)" 
> exec ^
: string
| CREATE TABLE tl (a visible)
< error
| 1064 42000 …line 1 column 26 near "visible)" 
> exec ^
: string
| CREATE TABLE tl (a auto_random)
< error
| 1064 42000 …line 1 column 30 near "auto_random)" 
> exec ^
: string
| CREATE TABLE tl (a)
< error
| 1064 42000 …line 1 column 19 near ")" 
> exec ^
: string
| CREATE TABLE tl (a, b int)
< error
| 1064 42000 …line 1 column 19 near ", b int)" 
> exec ^
: string
| CREATE TABLE tl (a) SELECT 1 AS a
< error
| 1064 42000 …line 1 column 19 near ") SELECT 1 AS a" 
> exec ^
: string
| CREATE TABLE tl (a) AS SELECT 1 AS a
< error
| 1064 42000 …line 1 column 19 near ") AS SELECT 1 AS a" 
> exec ^
: string
| CREATE TABLE tl (a, b) AS SELECT 1 AS a, 2 AS b
< error
| 1064 42000 …line 1 column 19 near ", b) AS SELECT 1 AS a, 2 AS b" 
> exec ^
: string
| CREATE TABLE tl (a key) SELECT 1 AS a
< error
| 1064 42000 …line 1 column 22 near "key) SELECT 1 AS a" 
> exec ^
: string
| CREATE TABLE tl (a default 1) SELECT 1 AS a
< error
| 1064 42000 …line 1 column 26 near "default 1) SELECT 1 AS a" 
> exec ^
: string
| CREATE TABLE tl (a not null) AS SELECT 1 AS a
< error
| 1064 42000 …line 1 column 22 near "not null) AS SELECT 1 AS a" 
> exec ^
: string
| CREATE TABLE tl (a comment 'x') AS SELECT 1 AS a
< error
| 1064 42000 …line 1 column 26 near "comment 'x') AS SELECT 1 AS a" 
> exec ^
: string
| CREATE TEMPORARY TABLE tl (a key)
< error
| 1064 42000 …line 1 column 32 near "key)" 
> exec ^
: string
| ALTER TABLE t ADD a key
< error
| 1064 42000 …line 1 column 23 near "key" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN a comment 'x'
< error
| 1064 42000 …line 1 column 34 near "comment 'x'" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN a generated always as (1)
< error
| 1064 42000 …line 1 column 36 near "generated always as (1)" 
> exec ^
: string
| ALTER TABLE t ADD (a key)
< error
| 1064 42000 …line 1 column 24 near "key)" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN (a not null, b int)
< error
| 1064 42000 …line 1 column 31 near "not null, b int)" 
> exec ^
: string
| ALTER TABLE t MODIFY a comment 'x'
< error
| 1064 42000 …line 1 column 30 near "comment 'x'" 
> exec ^
: string
| ALTER TABLE t MODIFY COLUMN a not null
< error
| 1064 42000 …line 1 column 33 near "not null" 
> exec ^
: string
| ALTER TABLE t MODIFY a generated always as (1)
< error
| 1064 42000 …line 1 column 32 near "generated always as (1)" 
> exec ^
: string
| ALTER TABLE t CHANGE a b key
< error
| 1064 42000 …line 1 column 28 near "key" 
> exec ^
: string
| ALTER TABLE t CHANGE COLUMN a b comment 'x'
< error
| 1064 42000 …line 1 column 39 near "comment 'x'" 
> exec ^
: string
| ALTER TABLE t CHANGE a b generated always as (1)
< error
| 1064 42000 …line 1 column 34 near "generated always as (1)" 
> exec ^
: string
| CREATE TABLE tl (a int key)
< ok 0 0
> exec ^
: string
| CREATE TABLE tl (a int generated always as (1))
< ok 0 0
> exec ^
: string
| CREATE TABLE tl (a int as (1) stored)
< ok 0 0
> exec ^
: string
| CREATE TABLE tl (a serial)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD COLUMN a int comment 'x'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD COLUMN a int comment 'x'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t MODIFY a int not null
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t MODIFY a int not null
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t CHANGE a b int key
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t CHANGE a b int key
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT d.comment(1)
< error
| 1305 42000 FUNCTION d.comment does not exist
> exec ^
: string
| SELECT d.comment(1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT `comment`(1)
< error
| 1305 42000 FUNCTION test.comment does not exist
> exec ^
: string
| SELECT `comment`(1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT d.vector(1)
< error
| 1305 42000 FUNCTION d.vector does not exist
> exec ^
: string
| SELECT d.vector(1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT comment FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT comment FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT 1 FROM comment
< error
| 1146 42S02 Table 'test.comment' doesn't exist
> exec ^
: string
| SELECT 1 FROM comment
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| INSERT INTO t VALUES (1) ON DUPLICATE KEY UPDATE a = VALUES(a)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| INSERT INTO t VALUES (1) ON DUPLICATE KEY UPDATE a = VALUES(a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT 1 = ANY (SELECT 1)
< ok 0 0
> exec ^
: string
| SELECT 1 FROM t WHERE a = SOME (SELECT 1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT 1 FROM t WHERE a = SOME (SELECT 1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT rank() OVER ()
< ok 0 0
> exec ^
: string
| SELECT lag(a) OVER () FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT lag(a) OVER () FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT ntile(1) OVER ()
< ok 0 0
> exec ^
: string
| SET @a = comment(1)
< error
| 1064 42000 …line 1 column 17 near "(1)" 
> exec ^
: string
| DO comment(1)
< error
| 1064 42000 …line 1 column 11 near "(1)" 
> exec ^
: string
| SELECT a FROM t ORDER BY comment(1)
< error
| 1064 42000 …line 1 column 33 near "(1)" 
> exec ^
: string
| SELECT a FROM t WHERE comment(1) = 1
< error
| 1064 42000 …line 1 column 30 near "(1) = 1" 
> exec ^
: string
| CREATE TABLE tf (a int DEFAULT (comment(1)))
< error
| 1064 42000 …line 1 column 39 near "comment(1)))" 
> exec ^
: string
| CREATE TABLE tf (a int, b int AS (comment(1)))
< error
| 1064 42000 …line 1 column 42 near "(1)))" 
> exec ^
: string
| CREATE VIEW vf AS SELECT comment(1)
< error
| 1064 42000 …line 1 column 33 near "(1)" 
> exec ^
: string
| CREATE FUNCTION comment() RETURNS INT RETURN 1
< error
| 1064 42000 …line 1 column 15 near "FUNCTION comment() RETURNS INT RETURN 1" 
> exec ^
: string
| CALL comment(1)
< error
| 1064 42000 …line 1 column 13 near "(1)" 
> exec ^
: string
| DROP FUNCTION comment
< error
| 1064 42000 …line 1 column 13 near "FUNCTION comment" 
> exec ^
: string
| CALL comment
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| CALL comment(1)
< error
| 1064 42000 …line 1 column 13 near "(1)" 
> exec ^
: string
| CALL tpcc(1)
< error
| 1064 42000 …line 1 column 10 near "(1)" 
> exec ^
: string
| CALL foo(1)
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| CALL foo(1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CALL `comment`(1)
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| CALL `comment`(1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CALL d.comment(1)
< error
| 1305 42000 FUNCTION d.comment does not exist
> exec ^
: string
| CALL d.comment(1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CALL d.comment
< error
| 1305 42000 FUNCTION d.comment does not exist
> exec ^
: string
| CALL d.comment
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CALL comment.foo(1)
< error
| 1305 42000 FUNCTION comment.foo does not exist
> exec ^
: string
| CALL comment.foo(1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
