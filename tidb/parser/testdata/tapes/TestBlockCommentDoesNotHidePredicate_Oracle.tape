-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestBlockCommentDoesNotHidePredicate_Oracle default#1
> exec
| CREATE DATABASE IF NOT EXISTS comment_nest_2
< ok 0 0
> exec
| USE comment_nest_2
< ok 0 0
> exec
| CREATE TABLE t3 (id INT)
< ok 0 0
> exec
| INSERT INTO t3 VALUES (1), (2), (3)
< ok 3 0
> exec
| DELETE FROM t3 WHERE id = 100 /* /* */ OR id = 1 -- */
< ok 1 0
> exec
| DROP DATABASE IF EXISTS comment_nest_2
< ok 0 0
