-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

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
| 1091 42000 index idx_nonexistent doesn't exist
> exec
| DROP TABLE IF EXISTS t_drop_key
< ok 0 0
