-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestWithinGroupClause default#1
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
| SELECT percentile_cont(0.5) WITHIN GROUP (ORDER BY 1)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT percentile_cont(0.5) WITHIN GROUP (ORDER BY 1) FILTER (WHERE true)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT percentile_cont(0.5) WITHIN GROUP ()
< error
| 42601 syntax error at or near ")"
| position: 43
> rollback
< ok
> begin
< ok
> exec
| SELECT percentile_cont(DISTINCT 0.5) WITHIN GROUP (ORDER BY 1)
< error
| 42601 cannot use DISTINCT with WITHIN GROUP
| position: 38
> rollback
< ok
> begin
< ok
> exec
| SELECT percentile_cont(VARIADIC ARRAY[0.5]) WITHIN GROUP (ORDER BY 1)
< error
| 42601 cannot use VARIADIC with WITHIN GROUP
| position: 45
> rollback
< ok
> begin
< ok
> exec
| SELECT percentile_cont(0.5 ORDER BY 1) WITHIN GROUP (ORDER BY 1)
< error
| 42601 cannot use multiple ORDER BY clauses with WITHIN GROUP
| position: 40
> rollback
< ok
