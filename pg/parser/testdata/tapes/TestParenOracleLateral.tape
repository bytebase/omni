-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestParenOracleLateral default#1
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
| SELECT * FROM T, LATERAL (SELECT 1) s
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T, LATERAL XMLTABLE('/root' PASSING T.doc COLUMNS a int PATH 'a') AS xt
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T, LATERAL JSON_TABLE(T.doc, '$' COLUMNS(a int PATH '$.a')) AS jt
< error
| 42846 cannot cast type xml to jsonb
| position: 37
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T, LATERAL (SELECT 1) AS s(x)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T, LATERAL (SELECT 1, 2) AS s(x, y)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T, LATERAL (SELECT 1, 2) s(x, y)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T, LATERAL XMLTABLE('/root' PASSING T.doc COLUMNS a int PATH 'a') AS xt
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T, LATERAL JSON_TABLE(T.doc, '$' COLUMNS(a int PATH '$.a')) AS jt
< error
| 42846 cannot cast type xml to jsonb
| position: 37
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T, LATERAL probe_f(T.x) AS fa(x)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T, LATERAL (SELECT * FROM U WHERE U.x = T.x) y
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T, LATERAL (SELECT T.x + 1) AS s(n)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T, LATERAL probe_f(T.x)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T, LATERAL ROWS FROM (probe_f(T.x), probe_g(T.y))
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T, LATERAL XMLTABLE('/root' PASSING T.doc COLUMNS a int PATH 'a') xt
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T, LATERAL (V JOIN W ON TRUE)
< error
| 42601 syntax error at or near "V"
| position: 27
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T, LATERAL ROWS FROM probe_f(T.x)
< error
| 42601 syntax error at or near "probe_f"
| position: 36
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T, LATERAL U
< error
| 42601 syntax error at end of input
| position: 27
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T, LATERAL ()
< error
| 42601 syntax error at or near ")"
| position: 27
> rollback
< ok
