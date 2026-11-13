-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestBlockCommentDoesNotNest_Oracle default#1
> exec
| CREATE DATABASE IF NOT EXISTS comment_nest_1
< ok 0 0
> exec
| USE comment_nest_1
< ok 0 0
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
| 		WHERE table_schema = 'comment_nest_1' AND table_name = 't2'
| 	
< result
< column "COUNT(*)" "BIGINT"
< row
: int64 0
> exec
| DROP DATABASE IF EXISTS comment_nest_1
< ok 0 0
