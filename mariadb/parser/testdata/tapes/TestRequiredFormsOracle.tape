-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestRequiredFormsOracle default#1
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE tl (a check (a > 0))
< error
| 1064 42000 …near 'a > 0))' at line 1
> exec ^
: string
| CREATE TABLE tl (a)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| CREATE TABLE tl (a, b int)
< error
| 1064 42000 …near ' b int)' at line 1
> exec ^
: string
| CREATE TABLE tl (a) SELECT 1 AS a
< error
| 1064 42000 …near ') SELECT 1 AS a' at line 1
> exec ^
: string
| CREATE TABLE tl (a) AS SELECT 1 AS a
< error
| 1064 42000 …near ') AS SELECT 1 AS a' at line 1
> exec ^
: string
| CREATE TABLE tl (a, b) AS SELECT 1 AS a, 2 AS b
< error
| 1064 42000 …near ' b) AS SELECT 1 AS a, 2 AS b' at line 1
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
< ok 0 0
> exec ^
: string
| ALTER TABLE t MODIFY a int not null
< ok 0 0
> exec ^
: string
| ALTER TABLE t CHANGE a b int key
< ok 0 0
> exec ^
: string
| SELECT d.comment(1)
< error
| 1630 42000 FUNCTION d.comment does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT d.comment(1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT `comment`(1)
< error
| 1630 42000 FUNCTION test.comment does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT `comment`(1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT d.vector(1)
< error
| 1630 42000 FUNCTION d.vector does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT d.vector(1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT rank() OVER ()
< error
| 4018 HY000 No order list in window specification for 'rank'
> exec ^
: string
| SELECT rank() OVER ()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT ntile(1) OVER ()
< ok 0 0
> exec ^
: string
| SET @a = comment(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| DO comment(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT a FROM t ORDER BY comment(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT a FROM t WHERE comment(1) = 1
< error
| 1064 42000 …near '(1) = 1' at line 1
> exec ^
: string
| CREATE TABLE tf (a int DEFAULT (comment(1)))
< error
| 1064 42000 …near '(1)))' at line 1
> exec ^
: string
| CREATE TABLE tf (a int, b int AS (comment(1)))
< error
| 1064 42000 …near '(1)))' at line 1
> exec ^
: string
| CREATE VIEW vf AS SELECT comment(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| CREATE FUNCTION comment() RETURNS INT RETURN 1
< ok 0 0
> exec ^
: string
| CALL comment(1)
< ok 0 0
> exec ^
: string
| DROP FUNCTION comment
< ok 0 0
