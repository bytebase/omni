-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestParenOracleArray default#1
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
| SELECT ARRAY[1, 2, 3] FROM T
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT ARRAY(SELECT 1) FROM T
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT ARRAY[ARRAY[1, 2], ARRAY[3, 4]] FROM T
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT ARRAY[ARRAY[ARRAY[1, 2]]] FROM T
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT ARRAY[[1, 2], [3, 4]] FROM T
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT ARRAY[ARRAY[1, 2], ARRAY[3]] FROM T
< error
| 2202E multidimensional arrays must have array expressions with matching dimensions
> rollback
< ok
> begin
< ok
> exec
| SELECT ARRAY(SELECT ARRAY(SELECT 1)) FROM T
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT ARRAY[1, 2]::int[] FROM T
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT ARRAY['a', 'b']::text[] FROM T
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(ARRAY[1, 2] AS int[]) FROM T
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT ARRAY(SELECT 1)::int[] FROM T
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT (ARRAY[1, 2])::int[] FROM T
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT ARRAY[1::text, 2::text] FROM T
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT ARRAY[]::int[] FROM T
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT ARRAY[[1, 2], [3, 4]]::int[][] FROM T
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT ARRAY(SELECT 1) FROM T
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT ARRAY(VALUES (1), (2)) FROM T
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT ARRAY(TABLE foo) FROM T
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT ARRAY(WITH cte AS (SELECT 1) SELECT * FROM cte) FROM T
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT ARRAY(SELECT 1 UNION SELECT 2) FROM T
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT ARRAY(SELECT a FROM U ORDER BY a) FROM T
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT ARRAY(SELECT a FROM U LIMIT 10) FROM T
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT ARRAY(SELECT DISTINCT a FROM U) FROM T
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT ARRAY() FROM T
< error
| 42601 syntax error at or near ")"
| position: 14
> rollback
< ok
> begin
< ok
> exec
| SELECT ARRAY[SELECT 1] FROM T
< error
| 42601 syntax error at or near "SELECT"
| position: 14
> rollback
< ok
> begin
< ok
> exec
| SELECT ARRAY FROM T
< error
| 42601 syntax error at or near "FROM"
| position: 14
> rollback
< ok
> begin
< ok
> exec
| SELECT ARRAY(1) FROM T
< error
| 42601 syntax error at or near "1"
| position: 14
> rollback
< ok
> begin
< ok
> exec
| SELECT ARRAY(1, 2) FROM T
< error
| 42601 syntax error at or near "1"
| position: 14
> rollback
< ok
> begin
< ok
> exec
| SELECT ARRAY[,] FROM T
< error
| 42601 syntax error at or near ","
| position: 14
> rollback
< ok
> begin
< ok
> exec
| SELECT ARRAY[1 2] FROM T
< error
| 42601 syntax error at or near "2"
| position: 16
> rollback
< ok
> begin
< ok
> exec
| SELECT ARRAY(ROWS FROM (probe_f(1))) FROM T
< error
| 42601 syntax error at or near "ROWS"
| position: 14
> rollback
< ok
> begin
< ok
> exec
| SELECT ARRAY[] FROM T
< error
| 42P18 cannot determine type of empty array
| position: 8
| hint: Explicitly cast to the desired type, for example ARRAY[]::integer[].
> rollback
< ok
> begin
< ok
> exec
| SELECT ARRAY[1, 2, 3]
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT ARRAY(SELECT 1)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT ARRAY[]::int[]
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT ARRAY[ARRAY[1, 2]]
< ok 1 -
> rollback
< ok
