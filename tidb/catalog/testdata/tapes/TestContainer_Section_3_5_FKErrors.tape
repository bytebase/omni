-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

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
| 1822 HY000 Failed to add the foreign key constraint. Missing index for constraint 'fk_1' in the referenced table 't_fk_parent_nokey'
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
| 3780 HY000 Referencing column 'pid' and referenced column 'id' in foreign key constraint 'fk_1' are incompatible.
> exec
| DROP TABLE IF EXISTS t_fk_child_mismatch
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_fk_parent_mismatch
< ok 0 0
