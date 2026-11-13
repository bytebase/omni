-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_PartitionChangeDetected default#1
> exec
| DROP DATABASE IF EXISTS pc_from_add_partitioning
< ok 0 0
> exec
| CREATE DATABASE pc_from_add_partitioning
< ok 1 0
> exec
| USE pc_from_add_partitioning
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pc_from_add_partitioning.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pc_to_add_partitioning
< ok 0 0
> exec
| CREATE DATABASE pc_to_add_partitioning
< ok 1 0
> exec
| USE pc_to_add_partitioning
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT) PARTITION BY HASH (id) PARTITIONS 4
< ok 0 0
> query
| SHOW CREATE TABLE pc_to_add_partitioning.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by hash(`id`)
| (partition `p0`,
| partition `p1`,
| partition `p2`,
| partition `p3`)
> exec
| DROP DATABASE IF EXISTS pc_from_remove_partitioning
< ok 0 0
> exec
| CREATE DATABASE pc_from_remove_partitioning
< ok 1 0
> exec
| USE pc_from_remove_partitioning
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT) PARTITION BY HASH (id) PARTITIONS 4
< ok 0 0
> query
| SHOW CREATE TABLE pc_from_remove_partitioning.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by hash(`id`)
| (partition `p0`,
| partition `p1`,
| partition `p2`,
| partition `p3`)
> exec
| DROP DATABASE IF EXISTS pc_to_remove_partitioning
< ok 0 0
> exec
| CREATE DATABASE pc_to_remove_partitioning
< ok 1 0
> exec
| USE pc_to_remove_partitioning
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pc_to_remove_partitioning.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pc_from_change_count
< ok 0 0
> exec
| CREATE DATABASE pc_from_change_count
< ok 1 0
> exec
| USE pc_from_change_count
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY HASH (id) PARTITIONS 4
< ok 0 0
> query
| SHOW CREATE TABLE pc_from_change_count.t
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
|  partition by hash(`id`)
| (partition `p0`,
| partition `p1`,
| partition `p2`,
| partition `p3`)
> exec
| DROP DATABASE IF EXISTS pc_to_change_count
< ok 0 0
> exec
| CREATE DATABASE pc_to_change_count
< ok 1 0
> exec
| USE pc_to_change_count
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY HASH (id) PARTITIONS 8
< ok 0 0
> query
| SHOW CREATE TABLE pc_to_change_count.t
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
|  partition by hash(`id`)
| (partition `p0`,
| partition `p1`,
| partition `p2`,
| partition `p3`,
| partition `p4`,
| partition `p5`,
| partition `p6`,
| partition `p7`)
> exec
| DROP DATABASE IF EXISTS pc_from_change_type
< ok 0 0
> exec
| CREATE DATABASE pc_from_change_type
< ok 1 0
> exec
| USE pc_from_change_type
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY HASH (id) PARTITIONS 4
< ok 0 0
> query
| SHOW CREATE TABLE pc_from_change_type.t
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
|  partition by hash(`id`)
| (partition `p0`,
| partition `p1`,
| partition `p2`,
| partition `p3`)
> exec
| DROP DATABASE IF EXISTS pc_to_change_type
< ok 0 0
> exec
| CREATE DATABASE pc_to_change_type
< ok 1 0
> exec
| USE pc_to_change_type
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY KEY (id) PARTITIONS 4
< ok 0 0
> query
| SHOW CREATE TABLE pc_to_change_type.t
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
|  partition by key(`id`)
| (partition `p0`,
| partition `p1`,
| partition `p2`,
| partition `p3`)
> exec
| DROP DATABASE IF EXISTS pc_from_change_range_bound
< ok 0 0
> exec
| CREATE DATABASE pc_from_change_range_bound
< ok 1 0
> exec
| USE pc_from_change_range_bound
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt)) PARTITION BY RANGE (YEAR(dt)) (PARTITION p0 VALUES LESS THAN (2010), PARTITION pmax VALUES LESS THAN MAXVALUE)
< ok 0 0
> query
| SHOW CREATE TABLE pc_from_change_range_bound.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `dt` date NOT NULL,
|   PRIMARY KEY (`id`, `dt`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by range(year(`dt`))
| (partition `p0` values less than (2010),
| partition `pmax` values less than (MAXVALUE))
> exec
| DROP DATABASE IF EXISTS pc_to_change_range_bound
< ok 0 0
> exec
| CREATE DATABASE pc_to_change_range_bound
< ok 1 0
> exec
| USE pc_to_change_range_bound
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt)) PARTITION BY RANGE (YEAR(dt)) (PARTITION p0 VALUES LESS THAN (2015), PARTITION pmax VALUES LESS THAN MAXVALUE)
< ok 0 0
> query
| SHOW CREATE TABLE pc_to_change_range_bound.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `dt` date NOT NULL,
|   PRIMARY KEY (`id`, `dt`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by range(year(`dt`))
| (partition `p0` values less than (2015),
| partition `pmax` values less than (MAXVALUE))
> exec
| DROP DATABASE IF EXISTS pc_from_add_range_partition
< ok 0 0
> exec
| CREATE DATABASE pc_from_add_range_partition
< ok 1 0
> exec
| USE pc_from_add_range_partition
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt)) PARTITION BY RANGE (YEAR(dt)) (PARTITION p0 VALUES LESS THAN (2010), PARTITION pmax VALUES LESS THAN MAXVALUE)
< ok 0 0
> query
| SHOW CREATE TABLE pc_from_add_range_partition.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `dt` date NOT NULL,
|   PRIMARY KEY (`id`, `dt`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by range(year(`dt`))
| (partition `p0` values less than (2010),
| partition `pmax` values less than (MAXVALUE))
> exec
| DROP DATABASE IF EXISTS pc_to_add_range_partition
< ok 0 0
> exec
| CREATE DATABASE pc_to_add_range_partition
< ok 1 0
> exec
| USE pc_to_add_range_partition
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt)) PARTITION BY RANGE (YEAR(dt)) (PARTITION p0 VALUES LESS THAN (2010), PARTITION p1 VALUES LESS THAN (2020), PARTITION pmax VALUES LESS THAN MAXVALUE)
< ok 0 0
> query
| SHOW CREATE TABLE pc_to_add_range_partition.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `dt` date NOT NULL,
|   PRIMARY KEY (`id`, `dt`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by range(year(`dt`))
| (partition `p0` values less than (2010),
| partition `p1` values less than (2020),
| partition `pmax` values less than (MAXVALUE))
> exec
| DROP DATABASE IF EXISTS pc_from_toggle_linear
< ok 0 0
> exec
| CREATE DATABASE pc_from_toggle_linear
< ok 1 0
> exec
| USE pc_from_toggle_linear
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY KEY (id) PARTITIONS 4
< ok 0 0
> query
| SHOW CREATE TABLE pc_from_toggle_linear.t
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
|  partition by key(`id`)
| (partition `p0`,
| partition `p1`,
| partition `p2`,
| partition `p3`)
> exec
| DROP DATABASE IF EXISTS pc_to_toggle_linear
< ok 0 0
> exec
| CREATE DATABASE pc_to_toggle_linear
< ok 1 0
> exec
| USE pc_to_toggle_linear
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY LINEAR KEY (id) PARTITIONS 4
< error
| 1064 42000 …near 'LINEAR KEY (id) PARTITIONS 4' at line 1
