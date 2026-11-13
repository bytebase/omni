-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

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
| 1005 HY000 Can't create table `catalog_shared`.`t_fk_noref` (errno: 150 "Foreign key constraint is incorrectly formed")
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
| 1005 HY000 Can't create table `catalog_shared`.`t_fk_child_nokey` (errno: 150 "Foreign key constraint is incorrectly formed")
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
| 1005 HY000 Can't create table `catalog_shared`.`t_fk_child_mismatch` (errno: 150 "Foreign key constraint is incorrectly formed")
> exec
| DROP TABLE IF EXISTS t_fk_child_mismatch
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_fk_parent_mismatch
< ok 0 0
