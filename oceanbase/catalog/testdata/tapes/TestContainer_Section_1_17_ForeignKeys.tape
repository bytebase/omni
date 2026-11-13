-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestContainer_Section_1_17_ForeignKeys default#1
> exec
| CREATE TABLE t_parent5 (id INT NOT NULL, PRIMARY KEY (id)) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE t_fk5 (id INT, pid INT, FOREIGN KEY (pid) REFERENCES t_parent5(id) ON DELETE SET DEFAULT)
< error
| 1215 HY000 Cannot add foreign key constraint
