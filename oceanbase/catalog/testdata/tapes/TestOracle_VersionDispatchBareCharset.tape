-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_VersionDispatchBareCharset default#1
> exec
| DROP DATABASE IF EXISTS vd_bare
< ok 0 0
> exec
| CREATE DATABASE vd_bare
< ok 1 0
> exec
| USE vd_bare
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(50)) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE vd_bare.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `name` varchar(50) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS vdb
< ok 0 0
> exec
| CREATE DATABASE vdb DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vdb
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(50)) DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `vdb`.`t` ADD COLUMN `isbn` varchar(20)
< ok 0 0
> query
| SHOW CREATE TABLE vdb.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `name` varchar(50) DEFAULT NULL,
|   `isbn` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
