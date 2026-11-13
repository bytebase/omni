-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestContainer_Section_2_6_DropTable default#1
> exec
| DROP TABLE IF EXISTS t_drop1
< ok 0 0
> exec
| CREATE TABLE t_drop1 (id INT PRIMARY KEY, name VARCHAR(100))
< ok 0 0
> exec
| DROP TABLE t_drop1
< ok 0 0
> query
| SHOW CREATE TABLE t_drop1
< error
| 1146 42S02 Table 'test.t_drop1' doesn't exist
> exec
| DROP TABLE IF EXISTS t_drop_ine
< ok 0 0
> exec ^
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_dm1
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_dm2
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_dm3
< ok 0 0
> exec
| CREATE TABLE t_dm1 (id INT)
< ok 0 0
> exec
| CREATE TABLE t_dm2 (id INT)
< ok 0 0
> exec
| CREATE TABLE t_dm3 (id INT)
< ok 0 0
> exec
| DROP TABLE t_dm1, t_dm2, t_dm3
< ok 0 0
> query
| SHOW CREATE TABLE t_dm1
< error
| 1146 42S02 Table 'test.t_dm1' doesn't exist
> query
| SHOW CREATE TABLE t_dm2
< error
| 1146 42S02 Table 'test.t_dm2' doesn't exist
> query
| SHOW CREATE TABLE t_dm3
< error
| 1146 42S02 Table 'test.t_dm3' doesn't exist
> exec
| DROP TABLE IF EXISTS t_noexist
< ok 0 0
> exec
| DROP TABLE t_noexist
< error
| 1051 42S02 Unknown table 'test.t_noexist'
> exec
| DROP TEMPORARY TABLE IF EXISTS t_temp_drop
< ok 0 0
> exec
| CREATE TEMPORARY TABLE t_temp_drop (id INT, val VARCHAR(50))
< ok 0 0
> exec
| DROP TEMPORARY TABLE t_temp_drop
< ok 0 0
> query
| SHOW CREATE TABLE t_temp_drop
< error
| 1146 42S02 Table 'test.t_temp_drop' doesn't exist
> exec
| DROP TABLE IF EXISTS t_fk_child
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_fk_parent
< ok 0 0
> exec
| CREATE TABLE t_fk_parent (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t_fk_child (id INT, parent_id INT, FOREIGN KEY (parent_id) REFERENCES t_fk_parent(id))
< ok 0 0
> exec
| DROP TABLE t_fk_parent
< error
| 3730 HY000 Cannot drop table 't_fk_parent' referenced by a foreign key constraint 't_fk_child_ibfk_1' on table 't_fk_child'.
