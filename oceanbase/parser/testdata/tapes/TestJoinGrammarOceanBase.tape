-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestJoinGrammarOceanBase default#1
> exec
| CREATE DATABASE IF NOT EXISTS objoin
< ok 1 0
> exec
| USE objoin
< ok 0 0
> exec
| DROP TABLE IF EXISTS t1, t2, t3
< ok 0 0
> exec
| CREATE TABLE t1 (id INT, a INT, b INT)
< ok 0 0
> exec
| CREATE TABLE t2 (id INT, a INT, b INT)
< ok 0 0
> exec
| CREATE TABLE t3 (id INT, a INT, b INT)
< ok 0 0
> exec
| SELECT * FROM t1 RIGHT JOIN t2 ON t1.id = t2.id LEFT JOIN t3
< error
| 1064 42000 …near 't3' at line 1
> exec
| SELECT * FROM t1 RIGHT JOIN t2 USING (id) LEFT JOIN t3
< error
| 1064 42000 …near 't3' at line 1
> exec
| SELECT * FROM t3 JOIN t1 JOIN t2 ON t1.a=t2.a ON t3.b=t2.b
< ok 0 0
> exec
| SELECT * FROM t1 LEFT JOIN t2 JOIN t3 ON t2.a=t3.a ON t1.b=t2.b
< ok 0 0
> exec
| SELECT * FROM t1 NATURAL JOIN t2 ON t1.a=t2.a
< error
| 1064 42000 …near 'ON t1.a=t2.a' at line 1
> exec
| SELECT * FROM t1 NATURAL JOIN t2 JOIN t3 ON t2.a=t3.a ON t1.b=t2.b
< error
| 1064 42000 …near 'ON t1.b=t2.b' at line 1
> exec
| SELECT * FROM t1 STRAIGHT_JOIN t2 JOIN t3 ON t2.a=t3.a ON t1.b=t2.b
< ok 0 0
> exec
| SELECT * FROM t1 STRAIGHT_JOIN t2 STRAIGHT_JOIN t3 ON t2.a=t3.a ON t1.b=t2.b
< ok 0 0
> exec
| DROP DATABASE IF EXISTS objoin
< ok 3 0
