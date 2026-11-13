-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

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
| SET SESSION sql_generate_invisible_primary_key = 1
< error
| 1193 HY000 Unknown system variable 'sql_generate_invisible_primary_key'
> exec
| SET SESSION show_gipk_in_create_table_and_information_schema = 1
< error
| 1193 HY000 Unknown system variable 'show_gipk_in_create_table_and_information_schema'
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> query
| SELECT COUNT(*) FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
< result
< column "COUNT(*)" "BIGINT"
< row
: int64 2
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (my_row_id INT, a INT)
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='my_row_id'
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| my_row_id
