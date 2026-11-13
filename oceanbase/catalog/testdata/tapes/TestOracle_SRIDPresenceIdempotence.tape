-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_SRIDPresenceIdempotence default#1
> exec
| DROP DATABASE IF EXISTS srid_srid_zero
< ok 0 0
> exec
| CREATE DATABASE srid_srid_zero
< ok 1 0
> exec
| USE srid_srid_zero
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0)
< ok 0 0
> query
| SHOW CREATE TABLE srid_srid_zero.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `g` geometry NOT NULL /*!80003 SRID 0 */,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS srid_srid_zero_idem
< ok 0 0
> exec
| CREATE DATABASE srid_srid_zero_idem
< ok 1 0
> exec
| USE srid_srid_zero_idem
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0)
< ok 0 0
> query
| SHOW CREATE TABLE srid_srid_zero_idem.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `g` geometry NOT NULL /*!80003 SRID 0 */,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS srid_srid_4326
< ok 0 0
> exec
| CREATE DATABASE srid_srid_4326
< ok 1 0
> exec
| USE srid_srid_4326
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 4326)
< ok 0 0
> query
| SHOW CREATE TABLE srid_srid_4326.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `g` geometry NOT NULL /*!80003 SRID 4326 */,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS srid_srid_4326_idem
< ok 0 0
> exec
| CREATE DATABASE srid_srid_4326_idem
< ok 1 0
> exec
| USE srid_srid_4326_idem
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 4326)
< ok 0 0
> query
| SHOW CREATE TABLE srid_srid_4326_idem.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `g` geometry NOT NULL /*!80003 SRID 4326 */,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS srid_srid_none
< ok 0 0
> exec
| CREATE DATABASE srid_srid_none
< ok 1 0
> exec
| USE srid_srid_none
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE srid_srid_none.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `g` geometry NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS srid_srid_none_idem
< ok 0 0
> exec
| CREATE DATABASE srid_srid_none_idem
< ok 1 0
> exec
| USE srid_srid_none_idem
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE srid_srid_none_idem.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `g` geometry NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS srid_srid_zero_nullable
< ok 0 0
> exec
| CREATE DATABASE srid_srid_zero_nullable
< ok 1 0
> exec
| USE srid_srid_zero_nullable
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY SRID 0)
< ok 0 0
> query
| SHOW CREATE TABLE srid_srid_zero_nullable.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `g` geometry /*!80003 SRID 0 */ DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS srid_srid_zero_nullable_idem
< ok 0 0
> exec
| CREATE DATABASE srid_srid_zero_nullable_idem
< ok 1 0
> exec
| USE srid_srid_zero_nullable_idem
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY SRID 0)
< ok 0 0
> query
| SHOW CREATE TABLE srid_srid_zero_nullable_idem.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `g` geometry /*!80003 SRID 0 */ DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS srid_srid_zero_multi_type
< ok 0 0
> exec
| CREATE DATABASE srid_srid_zero_multi_type
< ok 1 0
> exec
| USE srid_srid_zero_multi_type
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, p POINT NOT NULL SRID 0, l LINESTRING NOT NULL SRID 0, y POLYGON NOT NULL SRID 0)
< ok 0 0
> query
| SHOW CREATE TABLE srid_srid_zero_multi_type.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `p` point NOT NULL /*!80003 SRID 0 */,
|   `l` linestring NOT NULL /*!80003 SRID 0 */,
|   `y` polygon NOT NULL /*!80003 SRID 0 */,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS srid_srid_zero_multi_type_idem
< ok 0 0
> exec
| CREATE DATABASE srid_srid_zero_multi_type_idem
< ok 1 0
> exec
| USE srid_srid_zero_multi_type_idem
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, p POINT NOT NULL SRID 0, l LINESTRING NOT NULL SRID 0, y POLYGON NOT NULL SRID 0)
< ok 0 0
> query
| SHOW CREATE TABLE srid_srid_zero_multi_type_idem.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `p` point NOT NULL /*!80003 SRID 0 */,
|   `l` linestring NOT NULL /*!80003 SRID 0 */,
|   `y` polygon NOT NULL /*!80003 SRID 0 */,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS srid_srid_zero_spatial_index
< ok 0 0
> exec
| CREATE DATABASE srid_srid_zero_spatial_index
< ok 1 0
> exec
| USE srid_srid_zero_spatial_index
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0, SPATIAL INDEX sp_g (g))
< ok 0 0
> query
| SHOW CREATE TABLE srid_srid_zero_spatial_index.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `g` geometry NOT NULL /*!80003 SRID 0 */,
|   PRIMARY KEY (`id`),
|   SPATIAL KEY `sp_g` (`g`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS srid_srid_zero_spatial_index_idem
< ok 0 0
> exec
| CREATE DATABASE srid_srid_zero_spatial_index_idem
< ok 1 0
> exec
| USE srid_srid_zero_spatial_index_idem
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0, SPATIAL INDEX sp_g (g))
< ok 0 0
> query
| SHOW CREATE TABLE srid_srid_zero_spatial_index_idem.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `g` geometry NOT NULL /*!80003 SRID 0 */,
|   PRIMARY KEY (`id`),
|   SPATIAL KEY `sp_g` (`g`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS srid_srid_4326_spatial_index
< ok 0 0
> exec
| CREATE DATABASE srid_srid_4326_spatial_index
< ok 1 0
> exec
| USE srid_srid_4326_spatial_index
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 4326, SPATIAL INDEX sp_g (g))
< ok 0 0
> query
| SHOW CREATE TABLE srid_srid_4326_spatial_index.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `g` geometry NOT NULL /*!80003 SRID 4326 */,
|   PRIMARY KEY (`id`),
|   SPATIAL KEY `sp_g` (`g`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS srid_srid_4326_spatial_index_idem
< ok 0 0
> exec
| CREATE DATABASE srid_srid_4326_spatial_index_idem
< ok 1 0
> exec
| USE srid_srid_4326_spatial_index_idem
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 4326, SPATIAL INDEX sp_g (g))
< ok 0 0
> query
| SHOW CREATE TABLE srid_srid_4326_spatial_index_idem.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `g` geometry NOT NULL /*!80003 SRID 4326 */,
|   PRIMARY KEY (`id`),
|   SPATIAL KEY `sp_g` (`g`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS srid_srid_mixed_none_and_zero
< ok 0 0
> exec
| CREATE DATABASE srid_srid_mixed_none_and_zero
< ok 1 0
> exec
| USE srid_srid_mixed_none_and_zero
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a GEOMETRY, b GEOMETRY NOT NULL SRID 0)
< ok 0 0
> query
| SHOW CREATE TABLE srid_srid_mixed_none_and_zero.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` geometry DEFAULT NULL,
|   `b` geometry NOT NULL /*!80003 SRID 0 */,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS srid_srid_mixed_none_and_zero_idem
< ok 0 0
> exec
| CREATE DATABASE srid_srid_mixed_none_and_zero_idem
< ok 1 0
> exec
| USE srid_srid_mixed_none_and_zero_idem
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a GEOMETRY, b GEOMETRY NOT NULL SRID 0)
< ok 0 0
> query
| SHOW CREATE TABLE srid_srid_mixed_none_and_zero_idem.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` geometry DEFAULT NULL,
|   `b` geometry NOT NULL /*!80003 SRID 0 */,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
