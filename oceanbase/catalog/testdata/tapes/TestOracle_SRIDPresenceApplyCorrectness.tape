-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_SRIDPresenceApplyCorrectness default#1
> exec
| DROP DATABASE IF EXISTS pa_add_srid_none_to_zero_dc8b9c3b_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_srid_none_to_zero_dc8b9c3b_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_srid_none_to_zero_dc8b9c3b_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_srid_none_to_zero_dc8b9c3b_to.t
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
| DROP DATABASE IF EXISTS pa_add_srid_none_to_zero_dc8b9c3b
< ok 0 0
> exec
| CREATE DATABASE pa_add_srid_none_to_zero_dc8b9c3b DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_srid_none_to_zero_dc8b9c3b
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_srid_none_to_zero_dc8b9c3b.t
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
| ALTER TABLE `pa_add_srid_none_to_zero_dc8b9c3b`.`t` MODIFY COLUMN `g` geometry NOT NULL /*!80003 SRID 0 */
< error
| 1235 0A000 alter geometry srid not supported
> exec
| DROP DATABASE IF EXISTS pa_add_srid_none_to_4326_264a95a8_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_srid_none_to_4326_264a95a8_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_srid_none_to_4326_264a95a8_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 4326)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_srid_none_to_4326_264a95a8_to.t
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
| DROP DATABASE IF EXISTS pa_add_srid_none_to_4326_264a95a8
< ok 0 0
> exec
| CREATE DATABASE pa_add_srid_none_to_4326_264a95a8 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_srid_none_to_4326_264a95a8
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_srid_none_to_4326_264a95a8.t
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
| ALTER TABLE `pa_add_srid_none_to_4326_264a95a8`.`t` MODIFY COLUMN `g` geometry NOT NULL /*!80003 SRID 4326 */
< error
| 1235 0A000 alter geometry srid not supported
> exec
| DROP DATABASE IF EXISTS pa_remove_srid_zero_to_none_a0a7c44e_to
< ok 0 0
> exec
| CREATE DATABASE pa_remove_srid_zero_to_none_a0a7c44e_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_remove_srid_zero_to_none_a0a7c44e_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE pa_remove_srid_zero_to_none_a0a7c44e_to.t
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
| DROP DATABASE IF EXISTS pa_remove_srid_zero_to_none_a0a7c44e
< ok 0 0
> exec
| CREATE DATABASE pa_remove_srid_zero_to_none_a0a7c44e DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_remove_srid_zero_to_none_a0a7c44e
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0)
< ok 0 0
> query
| SHOW CREATE TABLE pa_remove_srid_zero_to_none_a0a7c44e.t
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
| ALTER TABLE `pa_remove_srid_zero_to_none_a0a7c44e`.`t` MODIFY COLUMN `g` geometry NOT NULL
< error
| 1235 0A000 alter geometry srid not supported
> exec
| DROP DATABASE IF EXISTS pa_remove_srid_4326_to_none_b6944bbd_to
< ok 0 0
> exec
| CREATE DATABASE pa_remove_srid_4326_to_none_b6944bbd_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_remove_srid_4326_to_none_b6944bbd_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE pa_remove_srid_4326_to_none_b6944bbd_to.t
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
| DROP DATABASE IF EXISTS pa_remove_srid_4326_to_none_b6944bbd
< ok 0 0
> exec
| CREATE DATABASE pa_remove_srid_4326_to_none_b6944bbd DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_remove_srid_4326_to_none_b6944bbd
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 4326)
< ok 0 0
> query
| SHOW CREATE TABLE pa_remove_srid_4326_to_none_b6944bbd.t
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
| ALTER TABLE `pa_remove_srid_4326_to_none_b6944bbd`.`t` MODIFY COLUMN `g` geometry NOT NULL
< error
| 1235 0A000 alter geometry srid not supported
> exec
| DROP DATABASE IF EXISTS pa_change_srid_4326_to_zero_a185a42d_to
< ok 0 0
> exec
| CREATE DATABASE pa_change_srid_4326_to_zero_a185a42d_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_change_srid_4326_to_zero_a185a42d_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0)
< ok 0 0
> query
| SHOW CREATE TABLE pa_change_srid_4326_to_zero_a185a42d_to.t
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
| DROP DATABASE IF EXISTS pa_change_srid_4326_to_zero_a185a42d
< ok 0 0
> exec
| CREATE DATABASE pa_change_srid_4326_to_zero_a185a42d DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_change_srid_4326_to_zero_a185a42d
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 4326)
< ok 0 0
> query
| SHOW CREATE TABLE pa_change_srid_4326_to_zero_a185a42d.t
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
| ALTER TABLE `pa_change_srid_4326_to_zero_a185a42d`.`t` MODIFY COLUMN `g` geometry NOT NULL /*!80003 SRID 0 */
< error
| 1235 0A000 alter geometry srid not supported
> exec
| DROP DATABASE IF EXISTS pa_change_srid_zero_to_4326_601cfccd_to
< ok 0 0
> exec
| CREATE DATABASE pa_change_srid_zero_to_4326_601cfccd_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_change_srid_zero_to_4326_601cfccd_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 4326)
< ok 0 0
> query
| SHOW CREATE TABLE pa_change_srid_zero_to_4326_601cfccd_to.t
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
| DROP DATABASE IF EXISTS pa_change_srid_zero_to_4326_601cfccd
< ok 0 0
> exec
| CREATE DATABASE pa_change_srid_zero_to_4326_601cfccd DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_change_srid_zero_to_4326_601cfccd
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0)
< ok 0 0
> query
| SHOW CREATE TABLE pa_change_srid_zero_to_4326_601cfccd.t
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
| ALTER TABLE `pa_change_srid_zero_to_4326_601cfccd`.`t` MODIFY COLUMN `g` geometry NOT NULL /*!80003 SRID 4326 */
< error
| 1235 0A000 alter geometry srid not supported
> exec
| DROP DATABASE IF EXISTS pa_add_column_with_srid_zero_6a9e5080_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_column_with_srid_zero_6a9e5080_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_column_with_srid_zero_6a9e5080_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_column_with_srid_zero_6a9e5080_to.t
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
| DROP DATABASE IF EXISTS pa_add_column_with_srid_zero_6a9e5080
< ok 0 0
> exec
| CREATE DATABASE pa_add_column_with_srid_zero_6a9e5080 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_column_with_srid_zero_6a9e5080
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_column_with_srid_zero_6a9e5080.t
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
| ALTER TABLE `pa_add_column_with_srid_zero_6a9e5080`.`t` ADD COLUMN `g` geometry NOT NULL /*!80003 SRID 0 */
< error
| 1138 22004 Invalid use of NULL value
> exec
| DROP DATABASE IF EXISTS pa_create_srid_zero_spatial_index_8aee60fd_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_srid_zero_spatial_index_8aee60fd_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_srid_zero_spatial_index_8aee60fd_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0, SPATIAL INDEX sp_g (g))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_srid_zero_spatial_index_8aee60fd_to.t
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
| DROP DATABASE IF EXISTS pa_create_srid_zero_spatial_index_8aee60fd
< ok 0 0
> exec
| CREATE DATABASE pa_create_srid_zero_spatial_index_8aee60fd DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_srid_zero_spatial_index_8aee60fd
< ok 0 0
> exec
| CREATE TABLE `pa_create_srid_zero_spatial_index_8aee60fd`.`t` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL /*!80003 SRID 0 */,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_srid_zero_spatial_index_8aee60fd`.`t` ADD SPATIAL KEY `sp_g` (`g`) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_srid_zero_spatial_index_8aee60fd.t
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
