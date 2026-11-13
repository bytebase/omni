-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestSortByUsingOperator default#1
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
| SELECT 1 ORDER BY 1 USING >
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 ORDER BY 1 USING < NULLS LAST
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 ORDER BY 1 USING ~<~
< error
| 42883 operator does not exist: integer ~<~ integer
| position: 27
| hint: No operator matches the given name and argument types. You might need to add explicit type casts.
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 ORDER BY 1 USING OPERATOR(pg_catalog.>)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT array_agg(1 ORDER BY 1 USING >)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 ORDER BY 1 USING
< error
| 42601 syntax error at end of input
| position: 26
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 ORDER BY 1 USING NULLS FIRST
< error
| 42601 syntax error at or near "NULLS"
| position: 27
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 ORDER BY 1 USING OPERATOR(>
< error
| 42601 syntax error at end of input
| position: 37
> rollback
< ok
