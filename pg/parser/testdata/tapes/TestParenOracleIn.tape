-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestParenOracleIn default#1
> exec
| CREATE TABLE T (a int, b int, x int, y int, doc xml)
< ok 0 -
> exec
| CREATE TABLE U (a int, b int, x int, y int)
< ok 0 -
> exec
| CREATE TABLE V (a int, b int)
< ok 0 -
> exec
| CREATE TABLE W (a int, b int)
< ok 0 -
> exec
| CREATE TABLE foo (a int)
< ok 0 -
> exec
| CREATE FUNCTION probe_f(i int) RETURNS int LANGUAGE sql AS 'SELECT $1'
< ok 0 -
> exec
| CREATE FUNCTION probe_g(i int) RETURNS int LANGUAGE sql AS 'SELECT $1'
< ok 0 -
> begin
< ok
> exec
| SELECT * FROM T WHERE a IN (1, 2, 3)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T WHERE a IN (SELECT 1)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T WHERE a IN (0)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T WHERE a IN (0, 1)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T WHERE a IN (0, 1, 2, 3, 4, 5, 6, 7, 8, 9)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T WHERE a IN (0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T WHERE a IN (1, 2, 3)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T WHERE a IN (1.5, 2.5, 3.25)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T WHERE a IN ('a', 'b', 'c')
< error
| 22P02 invalid input syntax for type integer: "a"
| position: 29
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T WHERE a IN (true, false)
< error
| 42883 operator does not exist: integer = boolean
| position: 25
| hint: No operator matches the given name and argument types. You might need to add explicit type casts.
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T WHERE a IN (NULL, NULL)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T WHERE a IN (SELECT 1)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T WHERE a IN (VALUES (1), (2))
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T WHERE a IN (WITH cte AS (SELECT 1) SELECT * FROM cte)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T WHERE a IN (TABLE foo)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T WHERE a IN (SELECT 1 UNION SELECT 2)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T WHERE (a, b) IN ((1, 2), (3, 4))
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T WHERE (a, b) IN (SELECT a, b FROM U)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T WHERE a NOT IN (1, 2, 3)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T WHERE a NOT IN (SELECT a FROM U)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT count(*) FROM T GROUP BY a HAVING count(*) IN (1, 2, 3)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT CASE WHEN a IN (1, 2) THEN 'y' ELSE 'n' END FROM T
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE TABLE t_check (x int CHECK (x IN (1, 2, 3)))
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT (SELECT 1 WHERE U.x IN (SELECT a FROM V)) FROM U
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T JOIN U ON T.a IN (SELECT a FROM V)
< ok 0 -
> rollback
< ok
