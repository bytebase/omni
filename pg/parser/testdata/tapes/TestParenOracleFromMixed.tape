-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestParenOracleFromMixed default#1
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
| SELECT * FROM ((SELECT 1) x JOIN U ON TRUE)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM ((SELECT 1) x CROSS JOIN U)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM ((SELECT 1) x NATURAL JOIN U)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM ((SELECT 1 AS a) x JOIN (SELECT 2 AS a) y ON x.a = y.a)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM ((VALUES (1)) v(a) JOIN U ON TRUE)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM ((TABLE T) JOIN U ON TRUE)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM ((T) JOIN U ON TRUE)
< error
| 42601 syntax error at or near ")"
| position: 18
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (T JOIN (SELECT 1) s(x) ON TRUE)
< ok 0 -
> rollback
< ok
