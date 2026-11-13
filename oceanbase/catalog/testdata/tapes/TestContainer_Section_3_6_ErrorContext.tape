-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestContainer_Section_3_6_ErrorContext default#1
> exec
| DROP DATABASE IF EXISTS db_quote_test
< ok 0 0
> exec
| CREATE DATABASE db_quote_test
< ok 1 0
> exec ^
< error
| 1007 HY000 Can't create database 'db_quote_test'; database exists
> exec
| DROP TABLE IF EXISTS t_quote_test
< ok 0 0
> exec
| CREATE TABLE t_quote_test (id INT)
< ok 0 0
> exec ^
< error
| 1050 42S01 Table 't_quote_test' already exists
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
| 1091 42000 Can't DROP 'nokey'; check that column/key exists
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
| 1050 42S01 Table 't_pos1' already exists
> exec
| DROP TABLE IF EXISTS t_pos1
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_pos2
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_ifexists_none
< ok 0 0
> exec ^
< ok 0 0
> exec
| DROP DATABASE IF EXISTS db_ifexists_none
< ok 0 0
> exec ^
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_ifexists_none
< ok 0 0
> exec ^
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
< ok 1 0
> exec
| CREATE DATABASE IF NOT EXISTS db_ifne
< ok 1 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_ifne
: bytes
| CREATE TABLE `t_ifne` (
|   `id` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP TABLE IF EXISTS t_no_ifne
< ok 0 0
> exec
| CREATE TABLE t_no_ifne (id INT)
< ok 0 0
> exec ^
< error
| 1050 42S01 Table 't_no_ifne' already exists
> exec
| DROP TABLE IF EXISTS t_ifne
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_no_ifne
< ok 0 0
> exec
| DROP DATABASE IF EXISTS db_ifne
< ok 0 0
