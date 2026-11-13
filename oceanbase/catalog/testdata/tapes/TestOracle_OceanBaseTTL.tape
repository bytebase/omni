-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_OceanBaseTTL/columns default#1
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, dt DATETIME, ts TIMESTAMP, d DATE)
< ok 0 0
> exec
| ALTER TABLE a TTL (c + INTERVAL 1 DAY)
< error
| 10503 HY000 Column type of 'c' is not supported for TTL definition
> exec
| ALTER TABLE a TTL (d + INTERVAL 1 DAY)
< error
| 10503 HY000 Column type of 'd' is not supported for TTL definition
> exec
| ALTER TABLE a TTL (zz + INTERVAL 1 DAY)
< error
| 10502 HY000 TTL column 'zz' not exists
> exec
| ALTER TABLE a TTL (c + INTERVAL 1 DAY, zz + INTERVAL 1 DAY)
< error
| 10503 HY000 Column type of 'c' is not supported for TTL definition
> exec
| ALTER TABLE a TTL (zz + INTERVAL 1 DAY, c + INTERVAL 1 DAY)
< error
| 10502 HY000 TTL column 'zz' not exists
> exec
| ALTER TABLE a TTL (dt + INTERVAL 1 DAY, ts + INTERVAL 1 SECOND)
< ok 0 0
> exec
| CREATE TABLE t2 (b INT PRIMARY KEY, c INT) TTL (c + INTERVAL 1 DAY)
< error
| 10503 HY000 Column type of 'c' is not supported for TTL definition
> exec
| CREATE TABLE t2 (b INT PRIMARY KEY, c DATETIME) TTL (zz + INTERVAL 1 DAY)
< error
| 10502 HY000 TTL column 'zz' not exists
> exec
| CREATE TABLE t2 (b INT PRIMARY KEY, c DATETIME) TTL (c + INTERVAL 1 DAY)
< ok 0 0
