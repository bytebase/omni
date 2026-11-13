-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestContainer_Section_4_1_Partitioning default#1
> exec
| CREATE TABLE t_part_linear (id INT NOT NULL) PARTITION BY LINEAR HASH (id) PARTITIONS 4
< error
| 1064 42000 …near 'LINEAR HASH (id) PARTITIONS 4' at line 1
> exec
| CREATE TABLE t_alter_reorgp (
| 				id INT NOT NULL,
| 				val INT NOT NULL
| 			) PARTITION BY RANGE (val) (
| 				PARTITION p0 VALUES LESS THAN (10),
| 				PARTITION p1 VALUES LESS THAN (30)
| 			)
< ok 0 0
> exec
| ALTER TABLE t_alter_reorgp REORGANIZE PARTITION p1 INTO (
| 				PARTITION p1a VALUES LESS THAN (20),
| 				PARTITION p1b VALUES LESS THAN (30)
| 			)
< error
| 1235 0A000 reorganizing table without primary key(s) is not supported
> exec
| CREATE TABLE t_alter_coalp (id INT NOT NULL) PARTITION BY HASH (id) PARTITIONS 4
< ok 0 0
> exec
| ALTER TABLE t_alter_coalp COALESCE PARTITION 2
< error
| 1064 42000 …near 'COALESCE PARTITION 2' at line 1
> exec
| CREATE TABLE t_alter_exchp (
| 				id INT NOT NULL,
| 				val INT NOT NULL
| 			) PARTITION BY RANGE (val) (
| 				PARTITION p0 VALUES LESS THAN (10),
| 				PARTITION p1 VALUES LESS THAN MAXVALUE
| 			)
< ok 0 0
> exec
| CREATE TABLE t_alter_exchp_swap (
| 				id INT NOT NULL,
| 				val INT NOT NULL
| 			)
< ok 0 0
> exec
| ALTER TABLE t_alter_exchp EXCHANGE PARTITION p0 WITH TABLE t_alter_exchp_swap
< error
| 1064 42000 …near 't_alter_exchp_swap' at line 1
