-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

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
