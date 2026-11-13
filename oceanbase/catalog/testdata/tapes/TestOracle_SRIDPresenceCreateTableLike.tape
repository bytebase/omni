-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_SRIDPresenceCreateTableLike default#1
> exec
| DROP DATABASE IF EXISTS srid_like_clone_srid_zero
< ok 0 0
> exec
| CREATE DATABASE srid_like_clone_srid_zero
< ok 1 0
> exec
| USE srid_like_clone_srid_zero
< ok 0 0
> exec
| CREATE TABLE src (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0);
| CREATE TABLE clone LIKE src
< ok 0 0
> query
| SHOW CREATE TABLE srid_like_clone_srid_zero.clone
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| clone
: bytes
| CREATE TABLE `clone` (
|   `id` int(11) NOT NULL,
|   `g` geometry NOT NULL /*!80003 SRID 0 */,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS srid_like_clone_srid_4326
< ok 0 0
> exec
| CREATE DATABASE srid_like_clone_srid_4326
< ok 1 0
> exec
| USE srid_like_clone_srid_4326
< ok 0 0
> exec
| CREATE TABLE src (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 4326);
| CREATE TABLE clone LIKE src
< ok 0 0
> query
| SHOW CREATE TABLE srid_like_clone_srid_4326.clone
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| clone
: bytes
| CREATE TABLE `clone` (
|   `id` int(11) NOT NULL,
|   `g` geometry NOT NULL /*!80003 SRID 4326 */,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS srid_like_clone_srid_none
< ok 0 0
> exec
| CREATE DATABASE srid_like_clone_srid_none
< ok 1 0
> exec
| USE srid_like_clone_srid_none
< ok 0 0
> exec
| CREATE TABLE src (id INT PRIMARY KEY, g GEOMETRY NOT NULL);
| CREATE TABLE clone LIKE src
< ok 0 0
> query
| SHOW CREATE TABLE srid_like_clone_srid_none.clone
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| clone
: bytes
| CREATE TABLE `clone` (
|   `id` int(11) NOT NULL,
|   `g` geometry NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
