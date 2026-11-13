-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_OceanBaseNameOrder/CREATE_TABLE default#1
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, KEY idx_c (c))
< ok 0 0
> exec
| CREATE TABLE a (x INT, x INT)
< error
| 1060 42S21 Duplicate column name 'x'
> exec
| CREATE TABLE IF NOT EXISTS a (x INT, x INT)
< error
| 1060 42S21 Duplicate column name 'x'
> exec
| CREATE TABLE a (x INT, KEY k (x), KEY k (x))
< error
| 1061 42000 Duplicate key name 'k'
> exec
| CREATE TABLE a (x VARCHAR(300000))
< error
| 1074 42000 Column length too big for column 'x' (max = 262144)
> exec
| CREATE TABLE IF NOT EXISTS a (x INT)
< ok 0 0
> exec
| CREATE TABLE a LIKE a
< error
| 1066 42000 Not unique table/alias: 'a'
> exec
| CREATE TABLE IF NOT EXISTS a LIKE a
< error
| 1066 42000 Not unique table/alias: 'a'
> exec
| CREATE TABLE zz LIKE zz
< error
| 1066 42000 Not unique table/alias: 'zz'
> exec
| CREATE TABLE a LIKE qq
< error
| 1146 42S02 Table 'omni_t_oceanbase_catalog_name_order_00.qq' doesn't exist

== TestOracle_OceanBaseNameOrder/DROP_and_RENAME default#1
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, KEY idx_c (c))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT 1 AS x
< ok 0 0
> exec
| DROP TABLE zz, a, zz
< error
| 1066 42000 Not unique table/alias: 'zz'
> exec
| DROP TABLE IF EXISTS a, a
< error
| 1066 42000 Not unique table/alias: 'a'
> exec
| DROP VIEW a
< error
| 1347 HY000 'omni_t_oceanbase_catalog_name_order_01.a' is not VIEW
> exec
| DROP VIEW IF EXISTS v, a
< error
| 1347 HY000 'omni_t_oceanbase_catalog_name_order_01.a' is not VIEW
> exec
| RENAME TABLE zz TO qq
< error
| 1017 HY000 File not exist
> exec
| RENAME TABLE zz TO a
< error
| 1050 42S01 Table 'a' already exists

== TestOracle_OceanBaseNameOrder/keys default#1
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, KEY idx_c (c))
< ok 0 0
> exec
| ALTER TABLE a ALTER INDEX zz VISIBLE
< error
| 1176 42000 Key 'zz' doesn't exist in table 'a'
> exec
| ALTER TABLE a RENAME KEY zz TO idx_c
< error
| 1176 42000 Key 'zz' doesn't exist in table 'a'
> exec
| ALTER TABLE a RENAME KEY idx_c TO idx_c
< error
| 1061 42000 Duplicate key name 'idx_c'
> exec
| ALTER TABLE a RENAME KEY idx_c TO idx_d
< ok 0 0
