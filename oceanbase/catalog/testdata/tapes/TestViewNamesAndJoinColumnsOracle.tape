-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestViewNamesAndJoinColumnsOracle default#1
> exec
| CREATE TABLE n1 (a INT, b INT)
< ok 0 0
> exec
| CREATE TABLE j1 (k1 INT, a INT)
< ok 0 0
> exec
| CREATE TABLE j2 (k1 INT, k2 INT, b INT)
< ok 0 0
> exec
| CREATE TABLE j3 (k2 INT, k3 INT, c INT)
< ok 0 0
> exec
| CREATE TABLE j4 (k3 INT, d INT)
< ok 0 0
> exec
| CREATE TABLE j5 (a INT, k1 INT, b INT, k2 INT)
< ok 0 0
> exec
| CREATE TABLE p (k INT, p1 INT)
< ok 0 0
> exec
| CREATE TABLE q (q1 INT, q2 INT)
< ok 0 0
> exec
| CREATE TABLE r (r1 INT, k INT)
< ok 0 0
> exec
| CREATE TABLE s (s1 INT, s2 INT)
< ok 0 0
> exec
| CREATE TABLE u (k INT, u1 INT)
< ok 0 0
> exec
| CREATE VIEW a01 AS SELECT (WITH c AS (SELECT 1 AS x) SELECT x FROM c UNION SELECT 2) FROM n1
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'a01' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| (WITH c AS (SELECT 1 AS x) SELECT x FROM c UNION SELECT 2)
> exec
| CREATE VIEW a02 AS select (with c as (select 1 as x) select x from   c union select 2) from n1
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'a02' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| (with c as (select 1 as x) select x from   c union select 2)
> exec
| CREATE VIEW a03 AS SELECT (SELECT 1 UNION SELECT 2) FROM n1
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'a03' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| (SELECT 1 UNION SELECT 2)
> exec
| CREATE VIEW a04 AS SELECT (SELECT a FROM n1 GROUP BY a HAVING a > 0 ORDER BY a LIMIT 1), (SELECT DISTINCT a FROM n1 LIMIT 1) FROM n1
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'a04' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| (SELECT a FROM n1 GROUP BY a HAVING a > 0 ORDER BY a LIMIT 1)
< row
: bytes
| (SELECT DISTINCT a FROM n1 LIMIT 1)
> exec
| CREATE VIEW a05 AS SELECT a+1, a  +  1, count(*), cast(a as char) FROM n1 GROUP BY a
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'a05' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| a+1
< row
: bytes
| a  +  1
< row
: bytes
| count(*)
< row
: bytes
| cast(a as char)
> exec
| CREATE VIEW a06 AS SELECT (a), ((b)), (1), 0004, - 2, x'41', B'10', true, date '2020-01-01', ('e'), (null) FROM n1
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'a06' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| a
< row
: bytes
| b
< row
: bytes
| (1)
< row
: bytes
| 0004
< row
: bytes
| - 2
< row
: bytes
| x'41'
< row
: bytes
| B'10'
< row
: bytes
| true
< row
: bytes
| date '2020-01-01'
< row
: bytes
| e
< row
: bytes
| NULL
> exec
| CREATE VIEW a07 AS SELECT a /* c */ + 1 FROM n1
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'a07' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| a /* c */ + 1
> exec
| CREATE VIEW a08 AS SELECT concat('éééééééééééééééééééééééééééééééééééééééééééééééééé') FROM n1
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'a08' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| Name_exp_1
> exec
| CREATE VIEW a09 AS SELECT concat('eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee') FROM n1
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'a09' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| Name_exp_1
> exec
| CREATE VIEW a10 AS SELECT * FROM (SELECT a+1, b FROM n1) d
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'a10' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| a+1
< row
: bytes
| b
> exec
| CREATE VIEW a11 AS SELECT (select   1), (SELECT a FROM n1 LIMIT 1) + 0 FROM n1
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'a11' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| (select   1)
< row
: bytes
| (SELECT a FROM n1 LIMIT 1) + 0
> exec
| CREATE VIEW a12 AS SELECT a IN (SELECT a FROM n1), a in (1,2) FROM n1
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'a12' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| a IN (SELECT a FROM n1)
< row
: bytes
| a in (1,2)
> exec
| CREATE VIEW a13 AS SELECT a + 1 AS `x + y`, count(*) AS `count(z)` FROM n1
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'a13' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| x + y
< row
: bytes
| count(z)
> exec
| CREATE VIEW a14 AS SELECT EXISTS (SELECT 1 UNION SELECT 2), a IN (SELECT 1 UNION SELECT 2) FROM n1
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'a14' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| EXISTS (SELECT 1 UNION SELECT 2)
< row
: bytes
| a IN (SELECT 1 UNION SELECT 2)
> exec
| CREATE VIEW a15 AS SELECT * FROM (SELECT `a`+1 FROM n1) d
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'a15' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| `a`+1
> exec
| CREATE VIEW a16 AS SELECT * FROM (SELECT a AS `x``y` FROM n1) d
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'a16' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| x`y
> exec
| CREATE VIEW b01 AS SELECT * FROM (j1 JOIN j2 USING (k1)) JOIN j3 USING (k2)
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b01' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k2
< row
: bytes
| k1
< row
: bytes
| a
< row
: bytes
| b
< row
: bytes
| k3
< row
: bytes
| c
> exec
| CREATE VIEW b02 AS SELECT * FROM j3 JOIN (j1 JOIN j2 USING (k1)) USING (k2)
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b02' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k2
< row
: bytes
| k3
< row
: bytes
| c
< row
: bytes
| k1
< row
: bytes
| a
< row
: bytes
| b
> exec
| CREATE VIEW b03 AS SELECT j1.a, k2, c FROM (j1 JOIN j2 ON j1.k1 = j2.k1) JOIN j3 USING (k2)
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b03' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| a
< row
: bytes
| k2
< row
: bytes
| c
> exec
| CREATE VIEW b04 AS SELECT * FROM (j1 NATURAL JOIN j2) NATURAL JOIN j3
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b04' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k2
< row
: bytes
| k1
< row
: bytes
| a
< row
: bytes
| b
< row
: bytes
| k3
< row
: bytes
| c
> exec
| CREATE VIEW b05 AS SELECT * FROM j3 NATURAL JOIN (j1 NATURAL JOIN j2)
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b05' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k2
< row
: bytes
| k3
< row
: bytes
| c
< row
: bytes
| k1
< row
: bytes
| a
< row
: bytes
| b
> exec
| CREATE VIEW b06 AS SELECT * FROM j1 RIGHT JOIN j2 USING (k1)
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b06' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k1
< row
: bytes
| k2
< row
: bytes
| b
< row
: bytes
| a
> exec
| CREATE VIEW b07 AS SELECT k1 FROM j1 RIGHT JOIN j2 USING (k1)
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b07' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k1
> exec
| CREATE VIEW b08 AS SELECT * FROM (j1 RIGHT JOIN j2 USING (k1)) JOIN j3 USING (k2)
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b08' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k2
< row
: bytes
| k1
< row
: bytes
| b
< row
: bytes
| a
< row
: bytes
| k3
< row
: bytes
| c
> exec
| CREATE VIEW b09 AS SELECT * FROM (j1 JOIN j2 USING (k1)) RIGHT JOIN j3 USING (k2)
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b09' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k2
< row
: bytes
| k3
< row
: bytes
| c
< row
: bytes
| k1
< row
: bytes
| a
< row
: bytes
| b
> exec
| CREATE VIEW b10 AS SELECT * FROM j1 NATURAL RIGHT JOIN j2
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b10' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k1
< row
: bytes
| k2
< row
: bytes
| b
< row
: bytes
| a
> exec
| CREATE VIEW b11 AS SELECT * FROM ((j1 JOIN j2 USING (k1)) JOIN j3 USING (k2)) JOIN j4 USING (k3)
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b11' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k3
< row
: bytes
| k2
< row
: bytes
| k1
< row
: bytes
| a
< row
: bytes
| b
< row
: bytes
| c
< row
: bytes
| d
> exec
| CREATE VIEW b12 AS SELECT * FROM (SELECT k1, a FROM j1) d JOIN j2 USING (k1)
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b12' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k1
< row
: bytes
| a
< row
: bytes
| k2
< row
: bytes
| b
> exec
| CREATE VIEW b13 AS SELECT * FROM j2 JOIN (SELECT k2, c FROM j3) d USING (k2)
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b13' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k2
< row
: bytes
| k1
< row
: bytes
| b
< row
: bytes
| c
> exec
| CREATE VIEW b14 AS SELECT * FROM (SELECT k1, a FROM j1) d NATURAL JOIN j2
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b14' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k1
< row
: bytes
| a
< row
: bytes
| k2
< row
: bytes
| b
> exec
| CREATE VIEW b15 AS SELECT * FROM j5 JOIN (SELECT k1, k2 FROM j2) d USING (k2, k1)
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b15' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k1
< row
: bytes
| k2
< row
: bytes
| a
< row
: bytes
| b
> exec
| CREATE VIEW b16 AS SELECT * FROM j1 RIGHT JOIN j2 USING (k1) WHERE k1 > 0 ORDER BY k1
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b16' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k1
< row
: bytes
| k2
< row
: bytes
| b
< row
: bytes
| a
> exec
| CREATE VIEW b17 AS SELECT * FROM u JOIN ((q RIGHT JOIN p ON 1=1) JOIN s ON 1=1) USING (k)
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b17' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k
< row
: bytes
| u1
< row
: bytes
| q1
< row
: bytes
| q2
< row
: bytes
| p1
< row
: bytes
| s1
< row
: bytes
| s2
> exec
| CREATE VIEW b18 AS SELECT * FROM (p JOIN r USING (k)) RIGHT JOIN q ON 1=1
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b18' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k
< row
: bytes
| p1
< row
: bytes
| r1
< row
: bytes
| q1
< row
: bytes
| q2
> exec
| CREATE VIEW b19 AS SELECT * FROM j1 JOIN j2 USING (K1)
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b19' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k1
< row
: bytes
| a
< row
: bytes
| k2
< row
: bytes
| b
> exec
| CREATE VIEW b20 AS SELECT * FROM j4, j1 JOIN j2 USING (k1)
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b20' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k3
< row
: bytes
| d
< row
: bytes
| k1
< row
: bytes
| a
< row
: bytes
| k2
< row
: bytes
| b
> exec
| CREATE VIEW b21 AS SELECT * FROM j4 JOIN j1 ON j4.d = j1.a JOIN j2 USING (k1)
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b21' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k1
< row
: bytes
| k3
< row
: bytes
| d
< row
: bytes
| a
< row
: bytes
| k2
< row
: bytes
| b
> exec
| CREATE VIEW b22 AS SELECT * FROM u JOIN (q RIGHT JOIN (s RIGHT JOIN p ON 1=1) ON 1=1) USING (k)
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b22' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k
< row
: bytes
| u1
< row
: bytes
| q1
< row
: bytes
| q2
< row
: bytes
| s1
< row
: bytes
| s2
< row
: bytes
| p1
> exec
| CREATE VIEW b23 AS SELECT * FROM p NATURAL RIGHT JOIN q
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b23' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k
< row
: bytes
| p1
< row
: bytes
| q1
< row
: bytes
| q2
> exec
| CREATE VIEW b24 AS SELECT * FROM u NATURAL JOIN (p NATURAL RIGHT JOIN q)
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b24' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k
< row
: bytes
| u1
< row
: bytes
| p1
< row
: bytes
| q1
< row
: bytes
| q2
> exec
| CREATE VIEW b25 AS SELECT * FROM (q NATURAL LEFT JOIN p) NATURAL JOIN u
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b25' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k
< row
: bytes
| q1
< row
: bytes
| q2
< row
: bytes
| p1
< row
: bytes
| u1
