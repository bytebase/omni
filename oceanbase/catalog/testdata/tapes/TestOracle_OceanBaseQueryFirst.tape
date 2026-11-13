-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_OceanBaseQueryFirst/before_the_name default#1
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT)
< ok 0 0
> exec
| CREATE VIEW a AS SELECT 1
< error
| 1050 42S01 Table 'a' already exists
> exec
| CREATE VIEW a AS SELECT zz FROM a
< error
| 1054 42S22 Unknown column 'zz' in 'field list'
> exec
| CREATE VIEW a AS SELECT * FROM zz
< error
| 1146 42S02 Table 'omni_t_oceanbase_catalog_query_first_00.zz' doesn't exist
> exec
| CREATE OR REPLACE VIEW a AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| CREATE TABLE a AS SELECT zz FROM a
< error
| 1054 42S22 Unknown column 'zz' in 'field list'
> exec
| CREATE TABLE IF NOT EXISTS a AS SELECT zz
< error
| 1054 42S22 Unknown column 'zz' in 'field list'
> exec
| CREATE TABLE IF NOT EXISTS a AS SELECT 1
< ok 0 0
> exec
| ALTER VIEW zz AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| ALTER VIEW zz AS SELECT * FROM qq
< error
| 1146 42S02 Table 'omni_t_oceanbase_catalog_query_first_00.qq' doesn't exist

== TestOracle_OceanBaseQueryFirst/names_in_the_query default#1
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT)
< ok 0 0
> exec
| CREATE VIEW v1 AS SELECT b FROM a ORDER BY zz
< error
| 1054 42S22 Unknown column 'zz' in 'order clause'
> exec
| CREATE VIEW v2 AS SELECT b FROM a GROUP BY zz
< error
| 1054 42S22 Unknown column 'zz' in 'group statement'
> exec
| CREATE VIEW v3 AS SELECT (SELECT zz) FROM a
< error
| 1054 42S22 Unknown column 'zz' in 'field list'
> exec
| CREATE VIEW v4 AS SELECT b AS q FROM a ORDER BY q
< ok 0 0
> exec
| CREATE VIEW v5 AS SELECT b + 1 AS q, count(*) AS n FROM a GROUP BY q HAVING n > 0
< ok 0 0
> exec
| CREATE VIEW v6 AS (SELECT b AS q FROM a) ORDER BY q
< ok 0 0
> exec
| CREATE VIEW v7 AS SELECT b AS q FROM a WHERE q > 0
< error
| 1054 42S22 Unknown column 'q' in 'where clause'
> exec
| CREATE VIEW v8 AS SELECT b AS q FROM a ORDER BY (SELECT q)
< ok 0 0
> exec
| CREATE VIEW v9 AS SELECT b AS q FROM a GROUP BY q HAVING (SELECT q) > 0
< ok 0 0
> exec
| CREATE VIEW v10 AS SELECT b AS q FROM a ORDER BY (SELECT zz)
< error
| 1054 42S22 Unknown column 'zz' in 'field list'
> exec
| CREATE VIEW v11 AS SELECT b AS q FROM a WHERE (SELECT q) > 0
< error
| 1054 42S22 Unknown column 'q' in 'field list'
> exec
| CREATE VIEW v12 AS WITH c AS (SELECT 1 AS x) (SELECT x FROM c) ORDER BY x
< ok 0 0
> exec
| CREATE VIEW v13 AS WITH c AS (SELECT 1 AS x) (SELECT x FROM c)
< ok 0 0
> exec
| CREATE VIEW v14 AS WITH c AS (SELECT 1 AS x) (SELECT x FROM c) ORDER BY zz
< error
| 1054 42S22 Unknown column 'zz' in 'order clause'
> exec
| CREATE VIEW v15 AS WITH c AS (SELECT 1 AS x) (SELECT zz FROM c)
< error
| 1054 42S22 Unknown column 'zz' in 'field list'
