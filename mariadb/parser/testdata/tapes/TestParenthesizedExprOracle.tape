-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestParenthesizedExprOracle default#1
> exec
| CREATE TABLE a (id INT, k VARCHAR(20))
< ok 0 0
> exec
| CREATE TABLE b (id INT, k VARCHAR(20))
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, d INT)
< ok 0 0
> exec
| CREATE TABLE u (x INT)
< ok 0 0
> exec
| SELECT count(*) FROM a WHERE ((a.id = 1) AND (a.k = 'x'))
< ok 0 0
> exec
| SELECT count(*) FROM a WHERE (((a.id = 1) AND (a.k = 'x')))
< ok 0 0
> exec
| SELECT count(*) FROM a t1 JOIN b t2 ON (((t1.id = t2.id) AND (t1.k = t2.k)))
< ok 0 0
> exec
| select `t1`.`id` AS `id`,`t2`.`k` AS `k` from (`a` `t1` join `b` `t2` on(((`t1`.`id` = `t2`.`id`) and (`t1`.`k` = `t2`.`k`) and (`t2`.`k` <> 'x') and (`t2`.`k` <> 'y'))))
< ok 0 0
> exec
| SELECT ((1))
< ok 0 0
> exec
| SELECT (((1 + 2)))
< ok 0 0
> exec
| SELECT ((a = b)) FROM t
< ok 0 0
> exec
| SELECT * FROM t WHERE ((((a = b))))
< ok 0 0
> exec
| SELECT ((a)) + ((b)) FROM t
< ok 0 0
> exec
| SELECT * FROM t WHERE a IN ((1), (2))
< ok 0 0
> exec
| SELECT (SELECT ((1)))
< ok 0 0
> exec
| SELECT * FROM (((a JOIN b ON a.id = b.id) JOIN t ON b.id = t.a) JOIN u ON t.a = u.x)
< ok 0 0
> exec
| SELECT ((SELECT 1))
< ok 0 0
> exec
| SELECT (((SELECT 1)))
< ok 0 0
> exec
| SELECT * FROM t WHERE a IN ((SELECT x FROM u))
< ok 0 0
> exec
| SELECT * FROM t WHERE EXISTS ((SELECT 1))
< ok 0 0
> exec
| SELECT * FROM t WHERE a > ALL ((SELECT x FROM u))
< ok 0 0
> exec
| SELECT * FROM ((SELECT 1)) x
< ok 0 0
> exec
| SELECT * FROM (((SELECT 1))) x
< ok 0 0
> exec
| SELECT * FROM ((SELECT 1) UNION (SELECT 2)) x
< ok 0 0
> exec
| SELECT * FROM t WHERE EXISTS ((1))
< error
| 1064 42000 …near '1))' at line 1
> exec
| SELECT * FROM t WHERE a > ALL ((1))
< error
| 1064 42000 …near '1))' at line 1
