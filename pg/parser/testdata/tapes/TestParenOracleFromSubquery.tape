-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestParenOracleFromSubquery default#1
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
| SELECT * FROM (SELECT 1)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (SELECT 1) AS s
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (SELECT 1) s(x)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM ((SELECT 1))
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (((SELECT 1)))
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM ((((SELECT 1))))
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (VALUES (1))
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (VALUES (1)) AS v(a)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM ((VALUES (1)) AS v(a))
< error
| 42601 syntax error at or near ")"
| position: 36
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (WITH cte AS (SELECT 1) SELECT * FROM cte)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM ((WITH cte AS (SELECT 1) SELECT * FROM cte))
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (TABLE T)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (SELECT 1 UNION SELECT 2)
< ok 2 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3)
< ok 3 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (SELECT 1 INTERSECT SELECT 2)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (SELECT 1 EXCEPT SELECT 2)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM ((SELECT 1) UNION (SELECT 2))
< ok 2 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (((SELECT 1) UNION (SELECT 2)) UNION (SELECT 3))
< ok 3 -
> rollback
< ok
