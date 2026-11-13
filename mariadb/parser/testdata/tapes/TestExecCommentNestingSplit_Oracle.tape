-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestExecCommentNestingSplit_Oracle default#1
> exec
| DROP TABLE IF EXISTS exec_nest_t1, exec_nest_t2
< ok 0 0
> exec
| CREATE TABLE exec_nest_t1 (id INT)
< ok 0 0
> exec
| CREATE TABLE exec_nest_t2 (id INT)
< ok 0 0
> exec
| INSERT INTO exec_nest_t1 VALUES (1) /*! , (2) /*! , (3) */ ; DROP TABLE exec_nest_t2; -- */
< ok 0 0
> query
| SELECT COUNT(*) FROM exec_nest_t1
< result
< column "COUNT(*)" "BIGINT"
< row
: int64 3
> query
| SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 'exec_nest_t2'
< result
< column "COUNT(*)" "BIGINT"
< row
: int64 0
> exec
| DROP TABLE IF EXISTS exec_nest_t1, exec_nest_t2
< ok 0 0
