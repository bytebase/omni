-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestContainer_Section_2_7_TruncateTable default#1
> exec
| CREATE TABLE t_trunc_ai (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(100)) AUTO_INCREMENT=1000
< ok 0 0
> query
| SHOW CREATE TABLE t_trunc_ai
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t_trunc_ai
: bytes
| CREATE TABLE `t_trunc_ai` (
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1000 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| TRUNCATE TABLE t_trunc_ai
< ok 0 0
> query
| SHOW CREATE TABLE t_trunc_ai
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t_trunc_ai
: bytes
| CREATE TABLE `t_trunc_ai` (
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP TABLE IF EXISTS t_trunc_noexist
< ok 0 0
> exec
| TRUNCATE TABLE t_trunc_noexist
< error
| 1146 42S02 Table 'test.t_trunc_noexist' doesn't exist
