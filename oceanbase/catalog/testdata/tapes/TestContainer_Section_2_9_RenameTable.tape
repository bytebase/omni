-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestContainer_Section_2_9_RenameTable default#1
> exec
| DROP TABLE IF EXISTS t_noexist_ren
< ok 0 0
> exec
| RENAME TABLE t_noexist_ren TO t_noexist_ren2
< error
| 1017 HY000 File not exist
> exec
| DROP TABLE IF EXISTS t_ren_exist1
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_ren_exist2
< ok 0 0
> exec
| CREATE TABLE t_ren_exist1 (id INT)
< ok 0 0
> exec
| CREATE TABLE t_ren_exist2 (id INT)
< ok 0 0
> exec
| RENAME TABLE t_ren_exist1 TO t_ren_exist2
< error
| 1050 42S01 Table 't_ren_exist2' already exists
