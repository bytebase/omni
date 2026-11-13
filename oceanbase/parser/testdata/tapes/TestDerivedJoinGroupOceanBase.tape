-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestDerivedJoinGroupOceanBase default#1
> exec
| CREATE DATABASE IF NOT EXISTS obderived
< ok 1 0
> exec
| USE obderived
< ok 0 0
> exec
| DROP TABLE IF EXISTS j1, j2, j3
< ok 0 0
> exec
| CREATE TABLE j1 (k1 INT, a INT, KEY i (k1))
< ok 0 0
> exec
| CREATE TABLE j2 (k1 INT, b INT)
< ok 0 0
> exec
| CREATE TABLE j3 (k1 INT, c INT)
< ok 0 0
> exec
| select `d`.`k1` AS `k1` from ((select `j1`.`k1` AS `k1` from `j1`) `d` join `j2` on((`d`.`k1` = `j2`.`k1`)))
< ok 0 0
> exec
| select d.k1 from ((select j1.k1 from j1) d join j2 on d.k1 = j2.k1)
< ok 0 0
> exec
| select * from ((select k1 from j1) as d join j2 on d.k1 = j2.k1)
< ok 0 0
> exec
| select * from ((select k1 from j1) d (x) join j2 on d.x = j2.k1)
< ok 0 0
> exec
| select * from ((select k1 from j1) as d (x) join j2 on d.x = j2.k1)
< ok 0 0
> exec
| select * from ((select k1 from j1) join j2 on 1)
< ok 0 0
> exec
| select * from ((select k1 from j1) d left join j2 on d.k1 = j2.k1)
< ok 0 0
> exec
| select * from ((select k1 from j1) d left outer join j2 on d.k1 = j2.k1)
< ok 0 0
> exec
| select * from ((select k1 from j1) d right join j2 on d.k1 = j2.k1)
< ok 0 0
> exec
| select * from ((select k1 from j1) d inner join j2 on d.k1 = j2.k1)
< ok 0 0
> exec
| select * from ((select k1 from j1) d cross join j2)
< ok 0 0
> exec
| select * from ((select k1 from j1) d join j2)
< ok 0 0
> exec
| select * from ((select k1 from j1) d join j2 using (k1))
< ok 0 0
> exec
| select * from ((select k1 from j1) d natural join j2)
< ok 0 0
> exec
| select * from ((select k1 from j1) d natural left join j2)
< ok 0 0
> exec
| select * from ((select k1 from j1) d straight_join j2)
< ok 0 0
> exec
| select * from ((select k1 from j1) d straight_join j2 on d.k1 = j2.k1)
< ok 0 0
> exec
| select * from ((select k1 from j1) d, j2)
< ok 0 0
> exec
| select * from ((select k1 from j1) d, (select k1 from j2) e)
< ok 0 0
> exec
| select * from (((select k1 from j1) d, j3) join j2 on 1)
< ok 0 0
> exec
| select * from (((select k1 from j1) d join j2 on d.k1 = j2.k1) join j3 on j3.k1 = d.k1)
< ok 0 0
> exec
| select * from ((((select k1 from j1) d join j2 on d.k1 = j2.k1) join j3 on j3.k1 = d.k1))
< ok 0 0
> exec
| select * from ((select k1 from j1) d join (j2 join j3 on j2.k1 = j3.k1) on d.k1 = j2.k1)
< ok 0 0
> exec
| select * from ((select k1 from j1) d join ((select k1 from j2) e join j3 on e.k1 = j3.k1) on d.k1 = e.k1)
< ok 0 0
> exec
| select * from (j1 join ((select 1) d join j2 on 1) on 1)
< ok 0 0
> exec
| select * from j1 join ((select 1) d join j2 on 1) on 1
< ok 0 0
> exec
| select * from ((select 1) d join j2 on 1) join j3 on 1
< ok 0 0
> exec
| select * from ((select 1) d join j2 on 1), j3
< ok 0 0
> exec
| select * from (((select k1 from j1)) d join j2 on d.k1 = j2.k1)
< ok 0 0
> exec
| select * from ((((select k1 from j1))) d join j2 on d.k1 = j2.k1)
< ok 0 0
> exec
| select * from ((select k1 from j1 union select k1 from j2) d join j3 on d.k1 = j3.k1)
< ok 0 0
> exec
| select * from (((select k1 from j1) union (select k1 from j2)) d join j3 on d.k1 = j3.k1)
< ok 0 0
> exec
| select * from (((select k1 from j1) union (select k1 from j2) order by 1 limit 1) d join j3 on d.k1 = j3.k1)
< ok 0 0
> exec
| select * from (((select k1 from j1) except (select k1 from j2)) d join j3 on 1)
< ok 0 0
> exec
| select * from (((select k1 from j1) intersect (select k1 from j2)) d join j3 on 1)
< ok 0 0
> exec
| select * from ((with c as (select 1 k1) select k1 from c) d join j2 on 1)
< ok 0 0
> exec
| select * from ((values row(1)) d join j2 on 1)
< ok 0 0
> exec
| select * from (lateral (select 1) d join j2 on 1)
< ok 0 0
> exec
| select * from (j1 join lateral (select j1.k1) d on 1)
< ok 0 0
> exec
| select * from (j1, lateral (select j1.k1) d)
< ok 0 0
> exec
| update ((select k1 from j1) d join j2 on d.k1 = j2.k1) set j2.b = 1
< ok 0 0
> exec
| delete j2 from ((select k1 from j1) d join j2 on d.k1 = j2.k1)
< ok 0 0
> exec
| insert into j3 select d.k1, 1 from ((select k1 from j1) d join j2 on d.k1 = j2.k1)
< ok 0 0
> exec
| with c as (select 1 k1) select * from ((select k1 from c) d join j2 on 1)
< ok 0 0
> exec
| select (select 1 from ((select 1) d join j2 on 1) limit 1)
< ok 0 0
> exec
| select 1 from j1 where exists (select 1 from ((select 1) d join j2 on 1))
< ok 0 0
> exec
| select * from ((select k1 from j1)) d
< ok 0 0
> exec
| select * from (((select k1 from j1))) d
< ok 0 0
> exec
| select * from ((select k1 from j1) union (select k1 from j2)) d
< ok 0 0
> exec
| select * from ((select k1 from j1) union select k1 from j2) d
< ok 0 0
> exec
| select * from ((select k1 from j1) order by 1 limit 1) d
< ok 0 0
> exec
| select * from ((select 1) order by 1 for update) d
< ok 0 0
> exec
| select * from ((select 1) limit 1 lock in share mode) d
< ok 0 0
> exec
| select * from ((select k1 from j1) d)
< ok 0 0
> exec
| select * from (((select k1 from j1) d))
< ok 0 0
> exec
| select * from ((select k1 from j1) as d)
< ok 0 0
> exec
| select * from ((select k1 from j1) d (x))
< ok 0 0
> exec
| select * from ((select k1 from j1) d) join j2 on 1
< ok 0 0
> exec
| select * from (((select 1) d) join j2 on 1)
< ok 0 0
> exec
| select * from (((select 1) d), j2)
< ok 0 0
> exec
| select * from j1, (lateral (select j1.k1) d)
< ok 0 0
> exec
| select * from j1, ((lateral (select j1.k1) as d))
< ok 0 0
> exec
| select * from j1 join (lateral (select j1.k1) d) on 1
< ok 0 0
> exec
| select * from (json_table('[1]', '$[*]' columns (x int path '$')) jt)
< ok 0 0
> exec
| select * from ((json_table('[1]', '$[*]' columns (x int path '$')) as jt))
< ok 0 0
> exec
| select * from (json_table('[1]', '$[*]' columns (x int path '$')) jt) join j2 on 1
< ok 0 0
> exec
| select * from ((select k1 from j1) d join j2 on 1
< error
| 1064 42000 …near '1' at line 1
> exec
| select * from ((select k1 from j1) d join)
< error
| 1064 42000 …near ')' at line 1
> exec
| select * from ((select k1 from j1) d d2 join j2 on 1)
< error
| 1064 42000 …near 'd2 join j2 on 1)' at line 1
> exec
| select * from ((select k1 from j1) d join j2 on 1) d2
< error
| 1064 42000 …near 'd2' at line 1
> exec
| select * from ((select k1 from j1) d (x join j2 on 1)
< error
| 1064 42000 …near 'join j2 on 1)' at line 1
> exec
| select * from ((select k1 from j1) d use index (i) join j2 on 1)
< error
| 1064 42000 …near 'use index (i) join j2 on 1)' at line 1
> exec
| select * from ((select 1) union (select 2) d join j2 on 1)
< error
| 1064 42000 …near 'd join j2 on 1)' at line 1
> exec
| select * from ((select 1) order by 1 d join j2 on 1)
< error
| 1064 42000 …near 'd join j2 on 1)' at line 1
> exec
| select * from ((table j1) d join j2 on 1)
< error
| 1064 42000 …near 'j1) d join j2 on 1)' at line 1
> exec
| select * from ((select 1) for update) d
< error
| 1064 42000 …near 'for update) d' at line 1
> exec
| select * from ((select 1) lock in share mode) d
< error
| 1064 42000 …near 'lock in share mode) d' at line 1
> exec
| DROP DATABASE IF EXISTS obderived
< ok 4 0
