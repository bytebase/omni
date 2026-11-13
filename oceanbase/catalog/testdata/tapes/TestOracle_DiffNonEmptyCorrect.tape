-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_DiffNonEmptyCorrect default#1
> exec
| DROP DATABASE IF EXISTS diff_nf
< ok 0 0
> exec
| CREATE DATABASE diff_nf
< ok 1 0
> exec
| USE diff_nf
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE diff_nf.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS diff_nt
< ok 0 0
> exec
| CREATE DATABASE diff_nt
< ok 1 0
> exec
| USE diff_nt
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(20))
< ok 0 0
> query
| SHOW CREATE TABLE diff_nt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS diff_nf
< ok 1 0
> exec
| CREATE DATABASE diff_nf
< ok 1 0
> exec
| USE diff_nf
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(20))
< ok 0 0
> query
| SHOW CREATE TABLE diff_nf.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS diff_nt
< ok 1 0
> exec
| CREATE DATABASE diff_nt
< ok 1 0
> exec
| USE diff_nt
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE diff_nt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS diff_nf
< ok 1 0
> exec
| CREATE DATABASE diff_nf
< ok 1 0
> exec
| USE diff_nf
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SHOW CREATE TABLE diff_nf.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `v` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS diff_nt
< ok 1 0
> exec
| CREATE DATABASE diff_nt
< ok 1 0
> exec
| USE diff_nt
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v BIGINT)
< ok 0 0
> query
| SHOW CREATE TABLE diff_nt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `v` bigint(20) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS diff_nf
< ok 1 0
> exec
| CREATE DATABASE diff_nf
< ok 1 0
> exec
| USE diff_nf
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v INT DEFAULT 0)
< ok 0 0
> query
| SHOW CREATE TABLE diff_nf.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `v` int(11) DEFAULT '0',
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS diff_nt
< ok 1 0
> exec
| CREATE DATABASE diff_nt
< ok 1 0
> exec
| USE diff_nt
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v INT DEFAULT 1)
< ok 0 0
> query
| SHOW CREATE TABLE diff_nt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `v` int(11) DEFAULT '1',
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS diff_nf
< ok 1 0
> exec
| CREATE DATABASE diff_nf
< ok 1 0
> exec
| USE diff_nf
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v INT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE diff_nf.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `v` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS diff_nt
< ok 1 0
> exec
| CREATE DATABASE diff_nt
< ok 1 0
> exec
| USE diff_nt
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v INT NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE diff_nt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `v` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
