-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestContainer_Section_2_8_CreateDropIndex default#1
> exec
| DROP TABLE IF EXISTS t_ci5
< ok 0 0
> exec
| CREATE TABLE t_ci5 (id INT PRIMARY KEY, val INT)
< ok 0 0
> exec
| CREATE INDEX idx_val ON t_ci5 (val)
< ok 0 0
> exec
| CREATE INDEX IF NOT EXISTS idx_val ON t_ci5 (val)
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_ci6
< ok 0 0
> exec
| CREATE TABLE t_ci6 (id INT PRIMARY KEY, a INT, b INT)
< ok 0 0
> exec
| CREATE INDEX idx_a ON t_ci6 (a)
< ok 0 0
> exec
| CREATE INDEX idx_a ON t_ci6 (b)
< error
| 1061 42000 Duplicate key name 'idx_a'
> exec
| DROP TABLE IF EXISTS t_di2
< ok 0 0
> exec
| CREATE TABLE t_di2 (id INT PRIMARY KEY)
< ok 0 0
> exec
| DROP INDEX idx_noexist ON t_di2
< error
| 1091 42000 Can't DROP 'idx_noexist'; check that column/key exists
