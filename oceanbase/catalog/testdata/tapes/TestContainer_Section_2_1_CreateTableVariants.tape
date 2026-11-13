-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestContainer_Section_2_1_CreateTableVariants default#1
> exec
| DROP TABLE IF EXISTS t_ine
< ok 0 0
> exec
| CREATE TABLE t_ine (id INT)
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS t_ine (id INT)
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_view_conflict
< ok 0 0
> exec
| DROP VIEW IF EXISTS t_view_conflict
< ok 0 0
> exec
| CREATE VIEW t_view_conflict AS SELECT 1 AS a
< ok 0 0
> exec
| CREATE TABLE t_view_conflict (id INT)
< error
| 1050 42S01 Table 't_view_conflict' already exists
