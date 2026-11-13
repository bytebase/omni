-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestContainer_Section_3_5_FKErrors default#1
> exec
| DROP TABLE IF EXISTS t_fk_noref
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_nonexistent_parent
< ok 0 0
> exec
| CREATE TABLE t_fk_noref (id INT, pid INT, FOREIGN KEY (pid) REFERENCES t_nonexistent_parent(id))
< error
| 1824 HY000 Failed to open the referenced table 't_nonexistent_parent'
> exec
| DROP TABLE IF EXISTS t_fk_child_nokey
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_fk_parent_nokey
< ok 0 0
> exec
| CREATE TABLE t_fk_parent_nokey (id INT, val INT)
< ok 0 0
> exec
| CREATE TABLE t_fk_child_nokey (id INT, pid INT, FOREIGN KEY (pid) REFERENCES t_fk_parent_nokey(val))
< error
| 1822 HY000 Failed to add the foreign key constraint. Missing index for constraint 't_fk_child_nokey_ibfk_1' in the referenced table 't_fk_parent_nokey'
> exec
| DROP TABLE IF EXISTS t_fk_child_nokey
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_fk_parent_nokey
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_fk_child_mismatch
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_fk_parent_mismatch
< ok 0 0
> exec
| CREATE TABLE t_fk_parent_mismatch (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t_fk_child_mismatch (id INT, pid VARCHAR(50), FOREIGN KEY (pid) REFERENCES t_fk_parent_mismatch(id))
< error
| 3780 HY000 Referencing column 'pid' and referenced column 'id' in foreign key constraint 't_fk_child_mismatch_ibfk_1' are incompatible.
> exec
| DROP TABLE IF EXISTS t_fk_child_mismatch
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_fk_parent_mismatch
< ok 0 0
