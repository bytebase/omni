-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_SRIDPresenceDistinctness default#1
> exec
| DROP DATABASE IF EXISTS srid_distinct
< ok 0 0
> exec
| CREATE DATABASE srid_distinct
< ok 1 0
> exec
| USE srid_distinct
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, gz GEOMETRY NOT NULL SRID 0, gn GEOMETRY NOT NULL, gf GEOMETRY NOT NULL SRID 4326)
< ok 0 0
> query
| SHOW CREATE TABLE srid_distinct.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `gz` geometry NOT NULL /*!80003 SRID 0 */,
|   `gn` geometry NOT NULL,
|   `gf` geometry NOT NULL /*!80003 SRID 4326 */,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
