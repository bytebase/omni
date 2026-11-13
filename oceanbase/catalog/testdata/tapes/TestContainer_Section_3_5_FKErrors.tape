-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestContainer_Section_3_5_FKErrors default#1
> exec
| CREATE TABLE t_fk_noref (id INT, pid INT, FOREIGN KEY (pid) REFERENCES t_nonexistent_parent(id))
< error
| 1146 42S02 Table doesn't exist
> exec
| CREATE TABLE t_fk_parent_nokey (id INT, val INT)
< ok 0 0
> exec
| CREATE TABLE t_fk_child_nokey (id INT, pid INT, FOREIGN KEY (pid) REFERENCES t_fk_parent_nokey(val))
< error
| 1215 HY000 Cannot add foreign key constraint
> exec
| CREATE TABLE t_fk_parent_mismatch (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t_fk_child_mismatch (id INT, pid VARCHAR(50), FOREIGN KEY (pid) REFERENCES t_fk_parent_mismatch(id))
< error
| 1215 HY000 Cannot add foreign key constraint
