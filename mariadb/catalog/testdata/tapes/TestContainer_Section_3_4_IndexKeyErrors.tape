-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestContainer_Section_3_4_IndexKeyErrors default#1
> exec
| DROP TABLE IF EXISTS t_dup_key
< ok 0 0
> exec
| CREATE TABLE t_dup_key (a INT, b INT, KEY idx_a (a))
< ok 0 0
> exec
| ALTER TABLE t_dup_key ADD INDEX idx_a (b)
< error
| 1061 42000 Duplicate key name 'idx_a'
> exec
| DROP TABLE IF EXISTS t_dup_key
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_drop_key
< ok 0 0
> exec
| CREATE TABLE t_drop_key (a INT)
< ok 0 0
> exec
| ALTER TABLE t_drop_key DROP INDEX idx_nonexistent
< error
| 1091 42000 Can't DROP INDEX `idx_nonexistent`; check that it exists
> exec
| DROP TABLE IF EXISTS t_drop_key
< ok 0 0
