-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestParenOracleFromLateral default#1
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
| SELECT * FROM T, LATERAL (SELECT 1)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T, LATERAL (SELECT 1) x
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T, LATERAL ((SELECT 1))
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T, LATERAL ((SELECT 1) UNION (SELECT 2))
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T, LATERAL (a JOIN b ON t.x = a.x)
< error
| 42601 syntax error at or near "a"
| position: 27
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T, LATERAL XMLTABLE('/root' PASSING T.doc COLUMNS a int PATH 'a')
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T, LATERAL JSON_TABLE(T.doc, '$' COLUMNS(a int PATH '$.a'))
< error
| 42846 cannot cast type xml to jsonb
| position: 37
> rollback
< ok
