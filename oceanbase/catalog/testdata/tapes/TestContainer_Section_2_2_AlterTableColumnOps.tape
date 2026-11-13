-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestContainer_Section_2_2_AlterTableColumnOps default#1
> exec
| CREATE TABLE t_drop_pk (a INT, b INT, PRIMARY KEY (a, b))
< ok 0 0
> exec
| ALTER TABLE t_drop_pk DROP COLUMN a
< error
| 1235 0A000 drop rowkey column is not supported
> exec
| DROP TABLE IF EXISTS t_fk_child
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_fk_parent
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_fk_child
< ok 0 0
> exec
| CREATE TABLE t_fk_parent (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t_fk_child (id INT PRIMARY KEY, pid INT, FOREIGN KEY (pid) REFERENCES t_fk_parent(id))
< ok 0 0
> exec
| ALTER TABLE t_fk_child DROP COLUMN pid
< error
| 5320 HY000 'pid': used in a foreign key constraint
