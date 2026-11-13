-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestParenOracleHarness default#1
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
| SELECT * FROM (T JOIN U ON TRUE)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (T JOIN U ON TRUE)
< ok 0 -
> rollback
< ok
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
| SELECT * FROM (SELECT 1)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM ()
< error
| 42601 syntax error at or near ")"
| position: 16
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM ()
< error
| 42601 syntax error at or near ")"
| position: 16
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (T)
< error
| 42601 syntax error at or near ")"
| position: 17
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (T)
< error
| 42601 syntax error at or near ")"
| position: 17
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM ((T JOIN U ON TRUE) JOIN V ON TRUE)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM ((T JOIN U ON TRUE) JOIN V ON TRUE)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM T
< ok 0 -
> rollback
< ok
