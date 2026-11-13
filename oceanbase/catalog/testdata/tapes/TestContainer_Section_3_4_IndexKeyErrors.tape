-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

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
| 1091 42000 Can't DROP 'idx_nonexistent'; check that column/key exists
> exec
| DROP TABLE IF EXISTS t_drop_key
< ok 0 0
