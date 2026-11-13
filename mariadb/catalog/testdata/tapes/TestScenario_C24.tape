-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestScenario_C24 default#1
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = ON
< error
| 1193 HY000 Unknown system variable 'sql_generate_invisible_primary_key'
> exec
| SET SESSION show_gipk_in_create_table_and_information_schema = ON
< error
| 1193 HY000 Unknown system variable 'show_gipk_in_create_table_and_information_schema'
> exec
| CREATE TABLE t_no_pk (a INT)
< ok 0 0
> query
| SHOW CREATE TABLE t_no_pk
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_no_pk
: bytes
| CREATE TABLE `t_no_pk` (
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
