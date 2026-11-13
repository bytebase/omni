-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestParenOracleFromSimple default#1
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
| SELECT * FROM (T)
< error
| 42601 syntax error at or near ")"
| position: 17
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM ((T))
< error
| 42601 syntax error at or near ")"
| position: 18
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
| SELECT * FROM (T CROSS JOIN U)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (T LEFT JOIN U ON TRUE)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (T RIGHT JOIN U ON TRUE)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (T FULL JOIN U ON TRUE)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (T INNER JOIN U ON TRUE)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (T NATURAL JOIN U)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (T NATURAL LEFT JOIN U)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (T NATURAL FULL OUTER JOIN U)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (T JOIN U USING (a))
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (T JOIN U USING (a) AS alias_clause)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (T LEFT OUTER JOIN U ON T.a = U.a)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM LATERAL (T)
< error
| 42601 syntax error at or near "T"
| position: 24
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (T AS alias)
< error
| 42601 syntax error at or near ")"
| position: 26
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (T, U)
< error
| 42601 syntax error at or near ","
| position: 17
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (ONLY T)
< error
| 42601 syntax error at or near ")"
| position: 22
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (T TABLESAMPLE BERNOULLI(10))
< error
| 42601 syntax error at or near ")"
| position: 43
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (T WITH ORDINALITY)
< error
| 42601 syntax error at or near "WITH"
| position: 18
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (T FOR UPDATE)
< error
| 42601 syntax error at or near "FOR"
| position: 18
> rollback
< ok
