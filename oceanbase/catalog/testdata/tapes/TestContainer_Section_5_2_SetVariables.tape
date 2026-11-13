-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestContainer_Section_5_2_SetVariables default#1
> exec
| DROP TABLE IF EXISTS t_fkoff_child
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t_fkoff_child (id INT, parent_id INT, FOREIGN KEY (parent_id) REFERENCES nonexistent_parent(id))
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_fkdrop_child
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_fkdrop_parent
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE t_fkdrop_parent (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t_fkdrop_child (id INT, parent_id INT, FOREIGN KEY (parent_id) REFERENCES t_fkdrop_parent(id))
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| DROP TABLE t_fkdrop_parent
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_fkon_child
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE t_fkon_child (id INT, parent_id INT, FOREIGN KEY (parent_id) REFERENCES nonexistent_on_parent(id))
< error
| 1146 42S02 Table doesn't exist
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET CHARACTER SET utf8mb4
< ok 0 0
> exec
| SET sql_mode = 'STRICT_TRANS_TABLES,NO_ENGINE_SUBSTITUTION'
< ok 0 0
