-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestContainer_Section_2_3_AlterTableIndexOps default#1
> exec
| DROP TABLE IF EXISTS t_dup_pk
< ok 0 0
> exec
| CREATE TABLE t_dup_pk (id INT PRIMARY KEY, val INT NOT NULL)
< ok 0 0
> exec
| ALTER TABLE t_dup_pk ADD PRIMARY KEY (val)
< error
| 1068 42000 Multiple primary key defined
