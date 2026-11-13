-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

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
| 		WHERE table_schema = 'test' AND table_name = 't2'
| 	
< result
< column "COUNT(*)" "BIGINT"
< row
: int64 0
