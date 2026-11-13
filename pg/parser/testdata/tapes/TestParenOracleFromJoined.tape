-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestParenOracleFromJoined default#1
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
| SELECT * FROM ((T JOIN U ON TRUE))
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (((T JOIN U ON TRUE)))
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
| SELECT * FROM (T JOIN (U JOIN V ON TRUE) ON TRUE)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM ((T JOIN U ON TRUE) JOIN (V JOIN W ON TRUE) ON TRUE)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (((T JOIN U ON TRUE) JOIN V ON TRUE) JOIN W ON TRUE)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (T JOIN (U JOIN (V JOIN W ON TRUE) ON TRUE) ON TRUE)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (T CROSS JOIN U JOIN V ON TRUE)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (T LEFT JOIN U ON T.a = U.a RIGHT JOIN V ON U.b = V.b)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM ((T JOIN U ON TRUE) JOIN V ON TRUE) AS jt
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM ((T JOIN U ON TRUE) JOIN V ON TRUE) AS jt(col1)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (T JOIN U ON TRUE) AS jt JOIN V ON TRUE
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (T NATURAL JOIN U) JOIN V USING (a)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (T FULL OUTER JOIN U USING (a))
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (T LEFT OUTER JOIN U USING (a, b))
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM ((T JOIN U ON ((T.a = U.a))) JOIN V ON ((U.b = V.b)))
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM (T JOIN LATERAL (SELECT U.a FROM U WHERE U.x = T.x) v ON TRUE)
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| SELECT * FROM ((T JOIN U ON T.a = U.a) JOIN LATERAL (SELECT 1) l ON TRUE)
< ok 0 -
> rollback
< ok
