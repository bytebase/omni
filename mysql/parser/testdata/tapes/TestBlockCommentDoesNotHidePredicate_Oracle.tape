-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

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
