-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestParenOracleFromDegenerate default#1
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
| SELECT * FROM ()
< error
| 42601 syntax error at or near ")"
| position: 16
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (
< error
| 42601 syntax error at end of input
| position: 16
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM )
< error
| 42601 syntax error at or near ")"
| position: 15
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (SELECT 1
< error
| 42601 syntax error at end of input
| position: 24
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (T JOIN
< error
| 42601 syntax error at end of input
| position: 22
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (T JOIN U)
< error
| 42601 syntax error at or near ")"
| position: 24
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (T CROSS JOIN U ON TRUE)
< error
| 42601 syntax error at or near "ON"
| position: 31
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (T NATURAL JOIN U ON TRUE)
< error
| 42601 syntax error at or near "ON"
| position: 33
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (SELECT)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (SELECT 1))
< error
| 42601 syntax error at or near ")"
| position: 25
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM ((SELECT 1)
< error
| 42601 syntax error at end of input
| position: 26
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (SELECT 1 FROM)
< error
| 42601 syntax error at or near ")"
| position: 29
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM ( ) SELECT 1
< error
| 42601 syntax error at or near ")"
| position: 17
> rollback
< ok
