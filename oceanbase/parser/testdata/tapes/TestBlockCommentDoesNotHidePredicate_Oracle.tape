-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestBlockCommentDoesNotHidePredicate_Oracle default#1
> exec
| CREATE TABLE t3 (id INT)
< ok 0 0
> exec
| INSERT INTO t3 VALUES (1), (2), (3)
< ok 3 0
> exec
| DELETE FROM t3 WHERE id = 100 /* /* */ OR id = 1 -- */
< ok 1 0
