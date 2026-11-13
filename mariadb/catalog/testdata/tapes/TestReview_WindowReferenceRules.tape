-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_WindowReferenceRules default#1
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE VIEW wv0 AS SELECT ROW_NUMBER() OVER (w PARTITION BY b) FROM t WINDOW w AS (ORDER BY a)
< error
| 4011 HY000 Window specification referencing another one 'w' cannot contain partition list
> exec
| CREATE VIEW wv1 AS SELECT ROW_NUMBER() OVER (w PARTITION BY b) FROM t WINDOW w AS ()
< error
| 4011 HY000 Window specification referencing another one 'w' cannot contain partition list
> exec
| CREATE VIEW wv2 AS SELECT ROW_NUMBER() OVER (w ORDER BY b) FROM t WINDOW w AS (ORDER BY a)
< error
| 4012 HY000 Referenced window specification 'w' already contains order list
> exec
| CREATE VIEW wv3 AS SELECT ROW_NUMBER() OVER (w ORDER BY b) FROM t WINDOW w AS (PARTITION BY a)
< ok 0 0
> exec
| CREATE VIEW wv4 AS SELECT SUM(b) OVER (w ROWS UNBOUNDED PRECEDING) FROM t WINDOW w AS (ORDER BY a ROWS 2 PRECEDING)
< error
| 4013 HY000 Referenced window specification 'w' cannot contain window frame
> exec
| CREATE VIEW wv5 AS SELECT SUM(b) OVER (w) FROM t WINDOW w AS (ORDER BY a ROWS 2 PRECEDING)
< error
| 4013 HY000 Referenced window specification 'w' cannot contain window frame
> exec
| CREATE VIEW wv6 AS SELECT SUM(b) OVER w FROM t WINDOW w AS (ORDER BY a ROWS 2 PRECEDING)
< ok 0 0
> exec
| CREATE VIEW wv7 AS SELECT SUM(b) OVER (w ORDER BY b) FROM t WINDOW w AS (ROWS 2 PRECEDING)
< error
| 4013 HY000 Referenced window specification 'w' cannot contain window frame
> exec
| CREATE VIEW wv8 AS SELECT SUM(b) OVER (w ROWS UNBOUNDED PRECEDING) FROM t WINDOW w AS (ORDER BY a)
< ok 0 0
> exec
| CREATE VIEW wv9 AS SELECT ROW_NUMBER() OVER w2 FROM t WINDOW w1 AS (ORDER BY a), w2 AS (w1)
< ok 0 0
> exec
| CREATE VIEW wv10 AS SELECT ROW_NUMBER() OVER w2 FROM t WINDOW w1 AS (ORDER BY a), w2 AS (w1 PARTITION BY b)
< error
| 4011 HY000 Window specification referencing another one 'w1' cannot contain partition list
> exec
| CREATE VIEW wv11 AS SELECT ROW_NUMBER() OVER w2 FROM t WINDOW w1 AS (PARTITION BY a), w2 AS (w1 ORDER BY b)
< ok 0 0
> exec
| CREATE VIEW wv12 AS SELECT ROW_NUMBER() OVER w2 FROM t WINDOW w1 AS (ORDER BY a ROWS 2 PRECEDING), w2 AS (w1)
< error
| 4013 HY000 Referenced window specification 'w1' cannot contain window frame
> exec
| CREATE VIEW wv13 AS SELECT ROW_NUMBER() OVER w2 FROM t WINDOW w1 AS (ORDER BY a ROWS 2 PRECEDING), w2 AS (w1 ORDER BY b)
< error
| 4012 HY000 Referenced window specification 'w1' already contains order list
> exec
| CREATE VIEW wv14 AS SELECT ROW_NUMBER() OVER w1 FROM t WINDOW w1 AS (w2), w2 AS (ORDER BY a)
< error
| 4009 HY000 Window specification with name 'w2' is not defined
> exec
| CREATE VIEW wv15 AS SELECT ROW_NUMBER() OVER w FROM t WINDOW w AS (ORDER BY a), w AS (ORDER BY b)
< error
| 4010 HY000 Multiple window specifications with the same name 'w'
> exec
| CREATE VIEW wv16 AS SELECT ROW_NUMBER() OVER (nope ORDER BY b) FROM t
< error
| 4009 HY000 Window specification with name 'nope' is not defined
> exec
| CREATE VIEW wv17 AS SELECT ROW_NUMBER() OVER nope FROM t
< error
| 4009 HY000 Window specification with name 'nope' is not defined
> exec
| CREATE VIEW wv18 AS SELECT (SELECT SUM(b) OVER w FROM t WINDOW w AS (ORDER BY b) LIMIT 1) FROM t WINDOW w AS (ORDER BY a)
< ok 0 0
> exec
| CREATE VIEW wv19 AS SELECT (SELECT SUM(b) OVER nope FROM t LIMIT 1) FROM t WINDOW w AS (ORDER BY a)
< error
| 4009 HY000 Window specification with name 'nope' is not defined
