-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

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
