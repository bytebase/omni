-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestContainer_Section_3_6_ErrorContext default#1
> exec
| DROP DATABASE IF EXISTS db_quote_test
< ok 0 0
> exec
| CREATE DATABASE db_quote_test
< ok 0 0
> exec
| CREATE DATABASE db_quote_test
< error
| 1007 HY000 Can't create database 'db_quote_test'; database exists
> exec
| DROP TABLE IF EXISTS t_quote_test
< ok 0 0
> exec
| CREATE TABLE t_quote_test (id INT)
< ok 0 0
> exec
| CREATE TABLE t_quote_test (id INT)
< error
| 1050 42S01 Table 'test.t_quote_test' already exists
> exec
| DROP TABLE IF EXISTS t_col_quote
< ok 0 0
> exec
| CREATE TABLE t_col_quote (id INT)
< ok 0 0
> exec
| ALTER TABLE t_col_quote DROP COLUMN nonexistent
< error
| 1091 42000 Can't DROP 'nonexistent'; check that column/key exists
> exec
| DROP TABLE IF EXISTS t_key_quote
< ok 0 0
> exec
| CREATE TABLE t_key_quote (id INT, val INT, KEY idx_val (val))
< ok 0 0
> exec
| ALTER TABLE t_key_quote ADD INDEX idx_val (id)
< error
| 1061 42000 Duplicate key name 'idx_val'
> exec
| DROP TABLE IF EXISTS t_dropkey_quote
< ok 0 0
> exec
| CREATE TABLE t_dropkey_quote (id INT)
< ok 0 0
> exec
| ALTER TABLE t_dropkey_quote DROP INDEX nokey
< error
| 1091 42000 index nokey doesn't exist
> exec
| DROP TABLE IF EXISTS t_unknown_quote
< ok 0 0
> exec
| DROP TABLE t_unknown_quote
< error
| 1051 42S02 Unknown table 'test.t_unknown_quote'
> exec
| DROP TABLE IF EXISTS t_quote_test
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_col_quote
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_key_quote
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_dropkey_quote
< ok 0 0
> exec
| DROP DATABASE IF EXISTS db_quote_test
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_pos1
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_pos2
< ok 0 0
> exec
| CREATE TABLE t_pos1 (id INT)
< ok 0 0
> exec
| CREATE TABLE t_pos2 (id INT)
< ok 0 0
> exec
| CREATE TABLE t_pos1 (id INT)
< error
| 1050 42S01 Table 'test.t_pos1' already exists
> exec
| DROP TABLE IF EXISTS t_pos1
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_pos2
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_ifexists_none
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_ifexists_none
< ok 0 0
> exec
| DROP DATABASE IF EXISTS db_ifexists_none
< ok 0 0
> exec
| DROP DATABASE IF EXISTS db_ifexists_none
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_ifexists_none
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_ifexists_none
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_ifexists_real
< ok 0 0
> exec
| CREATE TABLE t_ifexists_real (id INT)
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_ifexists_real
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_noifexists
< ok 0 0
> exec
| DROP TABLE t_noifexists
< error
| 1051 42S02 Unknown table 'test.t_noifexists'
> exec
| DROP DATABASE IF EXISTS db_ifne
< ok 0 0
> exec
| CREATE DATABASE db_ifne
< ok 0 0
> exec
| CREATE DATABASE IF NOT EXISTS db_ifne
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_ifne
< ok 0 0
> exec
| CREATE TABLE t_ifne (id INT)
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS t_ifne (val VARCHAR(100))
< ok 0 0
> query
| SHOW CREATE TABLE t_ifne
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_ifne
: bytes
| CREATE TABLE `t_ifne` (
|   `id` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t_no_ifne
< ok 0 0
> exec
| CREATE TABLE t_no_ifne (id INT)
< ok 0 0
> exec
| CREATE TABLE t_no_ifne (id INT)
< error
| 1050 42S01 Table 'test.t_no_ifne' already exists
> exec
| DROP TABLE IF EXISTS t_ifne
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_no_ifne
< ok 0 0
> exec
| DROP DATABASE IF EXISTS db_ifne
< ok 0 0
