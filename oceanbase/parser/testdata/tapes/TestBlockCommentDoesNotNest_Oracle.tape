-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestBlockCommentDoesNotNest_Oracle default#1
> exec
| CREATE TABLE t1 (id INT)
< ok 0 0
> exec
| CREATE TABLE t2 (id INT)
< ok 0 0
> exec
| INSERT INTO t1 VALUES (1) /* /* */; DROP TABLE t2; -- */
< ok 0 0
> query
|
| 		SELECT COUNT(*) FROM information_schema.tables
| 		WHERE table_schema = DATABASE() AND table_name = 't2'
| 	
< result
< column "COUNT(*)" "BIGINT"
< row
: int64 0
