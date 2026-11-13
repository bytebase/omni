-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_PartitionApplyCorrectness default#1
> exec
| DROP DATABASE IF EXISTS pto_create_range
< ok 0 0
> exec
| CREATE DATABASE pto_create_range DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pto_create_range
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt)) PARTITION BY RANGE (YEAR(dt)) (PARTITION p0 VALUES LESS THAN (2010), PARTITION p1 VALUES LESS THAN (2020), PARTITION pmax VALUES LESS THAN MAXVALUE)
< ok 0 0
> query
| SHOW CREATE TABLE pto_create_range.t
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
| DROP DATABASE IF EXISTS pa_create_range
< ok 0 0
> exec
| CREATE DATABASE pa_create_range DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_range
< ok 0 0
> exec
| CREATE TABLE `pa_create_range`.`t` (
|   `id` int NOT NULL,
|   `dt` date NOT NULL,
|   PRIMARY KEY (`id`,`dt`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_range`.`t` /*!50100 PARTITION BY RANGE (year(`dt`))
| (PARTITION p0 VALUES LESS THAN (2010) ENGINE = innodb,
|  PARTITION p1 VALUES LESS THAN (2020) ENGINE = innodb,
|  PARTITION pmax VALUES LESS THAN MAXVALUE ENGINE = innodb) */
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_range.t
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
| DROP DATABASE IF EXISTS pto_create_range_columns
< ok 0 0
> exec
| CREATE DATABASE pto_create_range_columns DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pto_create_range_columns
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, a INT NOT NULL, PRIMARY KEY (id,a)) PARTITION BY RANGE COLUMNS(a) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (MAXVALUE))
< ok 0 0
> query
| SHOW CREATE TABLE pto_create_range_columns.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) NOT NULL,
|   PRIMARY KEY (`id`, `a`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by range columns(`a`)
| (partition `p0` values less than (10),
| partition `p1` values less than (MAXVALUE))
> exec
| DROP DATABASE IF EXISTS pa_create_range_columns
< ok 0 0
> exec
| CREATE DATABASE pa_create_range_columns DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_range_columns
< ok 0 0
> exec
| CREATE TABLE `pa_create_range_columns`.`t` (
|   `id` int NOT NULL,
|   `a` int NOT NULL,
|   PRIMARY KEY (`id`,`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_range_columns`.`t` /*!50500 PARTITION BY RANGE  COLUMNS(a)
| (PARTITION p0 VALUES LESS THAN (10) ENGINE = innodb,
|  PARTITION p1 VALUES LESS THAN (MAXVALUE) ENGINE = innodb) */
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_range_columns.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) NOT NULL,
|   PRIMARY KEY (`id`, `a`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by range columns(`a`)
| (partition `p0` values less than (10),
| partition `p1` values less than (MAXVALUE))
> exec
| DROP DATABASE IF EXISTS pto_create_list
< ok 0 0
> exec
| CREATE DATABASE pto_create_list DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pto_create_list
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, c INT NOT NULL, PRIMARY KEY (id,c)) PARTITION BY LIST (c) (PARTITION pa VALUES IN (1,2,3), PARTITION pb VALUES IN (4,5,6))
< ok 0 0
> query
| SHOW CREATE TABLE pto_create_list.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` int(11) NOT NULL,
|   PRIMARY KEY (`id`, `c`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by list(`c`)
| (partition `pa` values in (1,2,3),
| partition `pb` values in (4,5,6))
> exec
| DROP DATABASE IF EXISTS pa_create_list
< ok 0 0
> exec
| CREATE DATABASE pa_create_list DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_list
< ok 0 0
> exec
| CREATE TABLE `pa_create_list`.`t` (
|   `id` int NOT NULL,
|   `c` int NOT NULL,
|   PRIMARY KEY (`id`,`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_list`.`t` /*!50100 PARTITION BY LIST (`c`)
| (PARTITION pa VALUES IN (1,2,3) ENGINE = innodb,
|  PARTITION pb VALUES IN (4,5,6) ENGINE = innodb) */
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_list.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` int(11) NOT NULL,
|   PRIMARY KEY (`id`, `c`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by list(`c`)
| (partition `pa` values in (1,2,3),
| partition `pb` values in (4,5,6))
> exec
| DROP DATABASE IF EXISTS pto_create_hash_n
< ok 0 0
> exec
| CREATE DATABASE pto_create_hash_n DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pto_create_hash_n
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT) PARTITION BY HASH (id) PARTITIONS 4
< ok 0 0
> query
| SHOW CREATE TABLE pto_create_hash_n.t
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
| DROP DATABASE IF EXISTS pa_create_hash_n
< ok 0 0
> exec
| CREATE DATABASE pa_create_hash_n DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_hash_n
< ok 0 0
> exec
| CREATE TABLE `pa_create_hash_n`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_hash_n`.`t` /*!50100 PARTITION BY HASH (`id`)
| (PARTITION p0 ENGINE = innodb,
|  PARTITION p1 ENGINE = innodb,
|  PARTITION p2 ENGINE = innodb,
|  PARTITION p3 ENGINE = innodb) */
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_hash_n.t
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
| DROP DATABASE IF EXISTS pto_create_key
< ok 0 0
> exec
| CREATE DATABASE pto_create_key DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pto_create_key
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY KEY (id) PARTITIONS 3
< ok 0 0
> query
| SHOW CREATE TABLE pto_create_key.t
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
| partition `p2`)
> exec
| DROP DATABASE IF EXISTS pa_create_key
< ok 0 0
> exec
| CREATE DATABASE pa_create_key DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_key
< ok 0 0
> exec
| CREATE TABLE `pa_create_key`.`t` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_key`.`t` /*!50100 PARTITION BY KEY (id)
| (PARTITION p0 ENGINE = innodb,
|  PARTITION p1 ENGINE = innodb,
|  PARTITION p2 ENGINE = innodb) */
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_key.t
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
| partition `p2`)
> exec
| DROP DATABASE IF EXISTS pto_create_linear_key
< ok 0 0
> exec
| CREATE DATABASE pto_create_linear_key DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pto_create_linear_key
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY LINEAR KEY (id) PARTITIONS 2
< error
| 1064 42000 …near 'LINEAR KEY (id) PARTITIONS 2' at line 1
> exec
| DROP DATABASE IF EXISTS pto_create_subpart_n
< ok 0 0
> exec
| CREATE DATABASE pto_create_subpart_n DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pto_create_subpart_n
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, d DATE NOT NULL, PRIMARY KEY (id,d)) PARTITION BY RANGE (YEAR(d)) SUBPARTITION BY HASH (id) SUBPARTITIONS 2 (PARTITION p0 VALUES LESS THAN (2010), PARTITION p1 VALUES LESS THAN MAXVALUE)
< ok 0 0
> query
| SHOW CREATE TABLE pto_create_subpart_n.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `d` date NOT NULL,
|   PRIMARY KEY (`id`, `d`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by range(year(`d`)) subpartition by hash(`id`) subpartition template (
| subpartition `p0`,
| subpartition `p1`)
| (partition `p0` values less than (2010),
| partition `p1` values less than (MAXVALUE))
> exec
| DROP DATABASE IF EXISTS pa_create_subpart_n
< ok 0 0
> exec
| CREATE DATABASE pa_create_subpart_n DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_subpart_n
< ok 0 0
> exec
| CREATE TABLE `pa_create_subpart_n`.`t` (
|   `id` int NOT NULL,
|   `d` date NOT NULL,
|   PRIMARY KEY (`id`,`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_subpart_n`.`t` /*!50100 PARTITION BY RANGE (year(`d`))
| SUBPARTITION BY HASH (`id`)
| (PARTITION p0 VALUES LESS THAN (2010)
|  (SUBPARTITION p0sp0 ENGINE = innodb,
|   SUBPARTITION p0sp1 ENGINE = innodb),
|  PARTITION p1 VALUES LESS THAN MAXVALUE
|  (SUBPARTITION p1sp0 ENGINE = innodb,
|   SUBPARTITION p1sp1 ENGINE = innodb)) */
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_subpart_n.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `d` date NOT NULL,
|   PRIMARY KEY (`id`, `d`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by range(year(`d`)) subpartition by hash(`id`)
| (partition `p0` values less than (2010) (
| subpartition `p0sp0`,
| subpartition `p0sp1`),
| partition `p1` values less than (MAXVALUE) (
| subpartition `p1sp0`,
| subpartition `p1sp1`))
> exec
| DROP DATABASE IF EXISTS pto_add_partitioning_hash
< ok 0 0
> exec
| CREATE DATABASE pto_add_partitioning_hash DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pto_add_partitioning_hash
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT) PARTITION BY HASH (id) PARTITIONS 4
< ok 0 0
> query
| SHOW CREATE TABLE pto_add_partitioning_hash.t
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
| DROP DATABASE IF EXISTS pfrom_add_partitioning_hash
< ok 0 0
> exec
| CREATE DATABASE pfrom_add_partitioning_hash DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pfrom_add_partitioning_hash
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pfrom_add_partitioning_hash.t
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
| DROP DATABASE IF EXISTS pa_add_partitioning_hash
< ok 0 0
> exec
| CREATE DATABASE pa_add_partitioning_hash DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_partitioning_hash
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT)
< ok 0 0
> exec
| ALTER TABLE `pa_add_partitioning_hash`.`t` /*!50100 PARTITION BY HASH (`id`)
| (PARTITION p0 ENGINE = innodb,
|  PARTITION p1 ENGINE = innodb,
|  PARTITION p2 ENGINE = innodb,
|  PARTITION p3 ENGINE = innodb) */
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_partitioning_hash.t
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
| DROP DATABASE IF EXISTS pto_add_partitioning_range
< ok 0 0
> exec
| CREATE DATABASE pto_add_partitioning_range DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pto_add_partitioning_range
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt)) PARTITION BY RANGE (YEAR(dt)) (PARTITION p0 VALUES LESS THAN (2010), PARTITION pmax VALUES LESS THAN MAXVALUE)
< ok 0 0
> query
| SHOW CREATE TABLE pto_add_partitioning_range.t
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
| DROP DATABASE IF EXISTS pfrom_add_partitioning_range
< ok 0 0
> exec
| CREATE DATABASE pfrom_add_partitioning_range DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pfrom_add_partitioning_range
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt))
< ok 0 0
> query
| SHOW CREATE TABLE pfrom_add_partitioning_range.t
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
> exec
| DROP DATABASE IF EXISTS pa_add_partitioning_range
< ok 0 0
> exec
| CREATE DATABASE pa_add_partitioning_range DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_partitioning_range
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt))
< ok 0 0
> exec
| ALTER TABLE `pa_add_partitioning_range`.`t` /*!50100 PARTITION BY RANGE (year(`dt`))
| (PARTITION p0 VALUES LESS THAN (2010) ENGINE = innodb,
|  PARTITION pmax VALUES LESS THAN MAXVALUE ENGINE = innodb) */
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_partitioning_range.t
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
| DROP DATABASE IF EXISTS pto_remove_partitioning
< ok 0 0
> exec
| CREATE DATABASE pto_remove_partitioning DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pto_remove_partitioning
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pto_remove_partitioning.t
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
| DROP DATABASE IF EXISTS pfrom_remove_partitioning
< ok 0 0
> exec
| CREATE DATABASE pfrom_remove_partitioning DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pfrom_remove_partitioning
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT) PARTITION BY HASH (id) PARTITIONS 4
< ok 0 0
> query
| SHOW CREATE TABLE pfrom_remove_partitioning.t
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
| DROP DATABASE IF EXISTS pa_remove_partitioning
< ok 0 0
> exec
| CREATE DATABASE pa_remove_partitioning DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_remove_partitioning
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT) PARTITION BY HASH (id) PARTITIONS 4
< ok 0 0
> exec
| ALTER TABLE `pa_remove_partitioning`.`t` REMOVE PARTITIONING
< error
| 1064 42000 …near 'PARTITIONING' at line 1
> exec
| DROP DATABASE IF EXISTS pto_remove_partitioning_range
< ok 0 0
> exec
| CREATE DATABASE pto_remove_partitioning_range DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pto_remove_partitioning_range
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt))
< ok 0 0
> query
| SHOW CREATE TABLE pto_remove_partitioning_range.t
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
> exec
| DROP DATABASE IF EXISTS pfrom_remove_partitioning_range
< ok 0 0
> exec
| CREATE DATABASE pfrom_remove_partitioning_range DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pfrom_remove_partitioning_range
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt)) PARTITION BY RANGE (YEAR(dt)) (PARTITION p0 VALUES LESS THAN (2010), PARTITION pmax VALUES LESS THAN MAXVALUE)
< ok 0 0
> query
| SHOW CREATE TABLE pfrom_remove_partitioning_range.t
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
| DROP DATABASE IF EXISTS pa_remove_partitioning_range
< ok 0 0
> exec
| CREATE DATABASE pa_remove_partitioning_range DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_remove_partitioning_range
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt)) PARTITION BY RANGE (YEAR(dt)) (PARTITION p0 VALUES LESS THAN (2010), PARTITION pmax VALUES LESS THAN MAXVALUE)
< ok 0 0
> exec
| ALTER TABLE `pa_remove_partitioning_range`.`t` REMOVE PARTITIONING
< error
| 1064 42000 …near 'PARTITIONING' at line 1
> exec
| DROP DATABASE IF EXISTS pto_repartition_count
< ok 0 0
> exec
| CREATE DATABASE pto_repartition_count DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pto_repartition_count
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY HASH (id) PARTITIONS 8
< ok 0 0
> query
| SHOW CREATE TABLE pto_repartition_count.t
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
| DROP DATABASE IF EXISTS pfrom_repartition_count
< ok 0 0
> exec
| CREATE DATABASE pfrom_repartition_count DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pfrom_repartition_count
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY HASH (id) PARTITIONS 4
< ok 0 0
> query
| SHOW CREATE TABLE pfrom_repartition_count.t
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
| DROP DATABASE IF EXISTS pa_repartition_count
< ok 0 0
> exec
| CREATE DATABASE pa_repartition_count DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_repartition_count
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY HASH (id) PARTITIONS 4
< ok 0 0
> exec
| ALTER TABLE `pa_repartition_count`.`t` /*!50100 PARTITION BY HASH (`id`)
| (PARTITION p0 ENGINE = innodb,
|  PARTITION p1 ENGINE = innodb,
|  PARTITION p2 ENGINE = innodb,
|  PARTITION p3 ENGINE = innodb,
|  PARTITION p4 ENGINE = innodb,
|  PARTITION p5 ENGINE = innodb,
|  PARTITION p6 ENGINE = innodb,
|  PARTITION p7 ENGINE = innodb) */
< ok 0 0
> query
| SHOW CREATE TABLE pa_repartition_count.t
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
| DROP DATABASE IF EXISTS pto_repartition_type
< ok 0 0
> exec
| CREATE DATABASE pto_repartition_type DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pto_repartition_type
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY KEY (id) PARTITIONS 4
< ok 0 0
> query
| SHOW CREATE TABLE pto_repartition_type.t
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
| DROP DATABASE IF EXISTS pfrom_repartition_type
< ok 0 0
> exec
| CREATE DATABASE pfrom_repartition_type DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pfrom_repartition_type
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY HASH (id) PARTITIONS 4
< ok 0 0
> query
| SHOW CREATE TABLE pfrom_repartition_type.t
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
| DROP DATABASE IF EXISTS pa_repartition_type
< ok 0 0
> exec
| CREATE DATABASE pa_repartition_type DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_repartition_type
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY HASH (id) PARTITIONS 4
< ok 0 0
> exec
| ALTER TABLE `pa_repartition_type`.`t` /*!50100 PARTITION BY KEY (id)
| (PARTITION p0 ENGINE = innodb,
|  PARTITION p1 ENGINE = innodb,
|  PARTITION p2 ENGINE = innodb,
|  PARTITION p3 ENGINE = innodb) */
< ok 0 0
> query
| SHOW CREATE TABLE pa_repartition_type.t
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
| DROP DATABASE IF EXISTS pto_repartition_range_bound
< ok 0 0
> exec
| CREATE DATABASE pto_repartition_range_bound DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pto_repartition_range_bound
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt)) PARTITION BY RANGE (YEAR(dt)) (PARTITION p0 VALUES LESS THAN (2015), PARTITION pmax VALUES LESS THAN MAXVALUE)
< ok 0 0
> query
| SHOW CREATE TABLE pto_repartition_range_bound.t
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
| DROP DATABASE IF EXISTS pfrom_repartition_range_bound
< ok 0 0
> exec
| CREATE DATABASE pfrom_repartition_range_bound DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pfrom_repartition_range_bound
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt)) PARTITION BY RANGE (YEAR(dt)) (PARTITION p0 VALUES LESS THAN (2010), PARTITION pmax VALUES LESS THAN MAXVALUE)
< ok 0 0
> query
| SHOW CREATE TABLE pfrom_repartition_range_bound.t
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
| DROP DATABASE IF EXISTS pa_repartition_range_bound
< ok 0 0
> exec
| CREATE DATABASE pa_repartition_range_bound DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_repartition_range_bound
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt)) PARTITION BY RANGE (YEAR(dt)) (PARTITION p0 VALUES LESS THAN (2010), PARTITION pmax VALUES LESS THAN MAXVALUE)
< ok 0 0
> exec
| ALTER TABLE `pa_repartition_range_bound`.`t` /*!50100 PARTITION BY RANGE (year(`dt`))
| (PARTITION p0 VALUES LESS THAN (2015) ENGINE = innodb,
|  PARTITION pmax VALUES LESS THAN MAXVALUE ENGINE = innodb) */
< ok 0 0
> query
| SHOW CREATE TABLE pa_repartition_range_bound.t
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
| DROP DATABASE IF EXISTS pto_repartition_add_range
< ok 0 0
> exec
| CREATE DATABASE pto_repartition_add_range DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pto_repartition_add_range
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt)) PARTITION BY RANGE (YEAR(dt)) (PARTITION p0 VALUES LESS THAN (2010), PARTITION p1 VALUES LESS THAN (2020), PARTITION pmax VALUES LESS THAN MAXVALUE)
< ok 0 0
> query
| SHOW CREATE TABLE pto_repartition_add_range.t
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
| DROP DATABASE IF EXISTS pfrom_repartition_add_range
< ok 0 0
> exec
| CREATE DATABASE pfrom_repartition_add_range DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pfrom_repartition_add_range
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt)) PARTITION BY RANGE (YEAR(dt)) (PARTITION p0 VALUES LESS THAN (2010), PARTITION pmax VALUES LESS THAN MAXVALUE)
< ok 0 0
> query
| SHOW CREATE TABLE pfrom_repartition_add_range.t
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
| DROP DATABASE IF EXISTS pa_repartition_add_range
< ok 0 0
> exec
| CREATE DATABASE pa_repartition_add_range DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_repartition_add_range
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt)) PARTITION BY RANGE (YEAR(dt)) (PARTITION p0 VALUES LESS THAN (2010), PARTITION pmax VALUES LESS THAN MAXVALUE)
< ok 0 0
> exec
| ALTER TABLE `pa_repartition_add_range`.`t` /*!50100 PARTITION BY RANGE (year(`dt`))
| (PARTITION p0 VALUES LESS THAN (2010) ENGINE = innodb,
|  PARTITION p1 VALUES LESS THAN (2020) ENGINE = innodb,
|  PARTITION pmax VALUES LESS THAN MAXVALUE ENGINE = innodb) */
< ok 0 0
> query
| SHOW CREATE TABLE pa_repartition_add_range.t
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
| DROP DATABASE IF EXISTS pto_repartition_toggle_linear
< ok 0 0
> exec
| CREATE DATABASE pto_repartition_toggle_linear DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pto_repartition_toggle_linear
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY LINEAR KEY (id) PARTITIONS 4
< error
| 1064 42000 …near 'LINEAR KEY (id) PARTITIONS 4' at line 1
> exec
| DROP DATABASE IF EXISTS pto_add_column_keep_partitioning
< ok 0 0
> exec
| CREATE DATABASE pto_add_column_keep_partitioning DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pto_add_column_keep_partitioning
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT) PARTITION BY HASH (id) PARTITIONS 4
< ok 0 0
> query
| SHOW CREATE TABLE pto_add_column_keep_partitioning.t
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
| DROP DATABASE IF EXISTS pfrom_add_column_keep_partitioning
< ok 0 0
> exec
| CREATE DATABASE pfrom_add_column_keep_partitioning DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pfrom_add_column_keep_partitioning
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY HASH (id) PARTITIONS 4
< ok 0 0
> query
| SHOW CREATE TABLE pfrom_add_column_keep_partitioning.t
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
| DROP DATABASE IF EXISTS pa_add_column_keep_partitioning
< ok 0 0
> exec
| CREATE DATABASE pa_add_column_keep_partitioning DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_column_keep_partitioning
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY HASH (id) PARTITIONS 4
< ok 0 0
> exec
| ALTER TABLE `pa_add_column_keep_partitioning`.`t` ADD COLUMN `a` int DEFAULT NULL
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_column_keep_partitioning.t
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
| DROP DATABASE IF EXISTS pto_remove_partitioning_drop_key_column
< ok 0 0
> exec
| CREATE DATABASE pto_remove_partitioning_drop_key_column DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pto_remove_partitioning_drop_key_column
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pto_remove_partitioning_drop_key_column.t
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
| DROP DATABASE IF EXISTS pfrom_remove_partitioning_drop_key_column
< ok 0 0
> exec
| CREATE DATABASE pfrom_remove_partitioning_drop_key_column DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pfrom_remove_partitioning_drop_key_column
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt)) PARTITION BY RANGE (YEAR(dt)) (PARTITION p0 VALUES LESS THAN (2010), PARTITION pmax VALUES LESS THAN MAXVALUE)
< ok 0 0
> query
| SHOW CREATE TABLE pfrom_remove_partitioning_drop_key_column.t
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
| DROP DATABASE IF EXISTS pa_remove_partitioning_drop_key_column
< ok 0 0
> exec
| CREATE DATABASE pa_remove_partitioning_drop_key_column DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_remove_partitioning_drop_key_column
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt)) PARTITION BY RANGE (YEAR(dt)) (PARTITION p0 VALUES LESS THAN (2010), PARTITION pmax VALUES LESS THAN MAXVALUE)
< ok 0 0
> exec
| ALTER TABLE `pa_remove_partitioning_drop_key_column`.`t` REMOVE PARTITIONING
< error
| 1064 42000 …near 'PARTITIONING' at line 1
> exec
| DROP DATABASE IF EXISTS pto_create_key_algorithm_2
< ok 0 0
> exec
| CREATE DATABASE pto_create_key_algorithm_2 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pto_create_key_algorithm_2
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY KEY ALGORITHM=2 (id) PARTITIONS 3
< error
| 1064 42000 …near 'ALGORITHM=2 (id) PARTITIONS 3' at line 1
> exec
| DROP DATABASE IF EXISTS pto_repartition_drop_old_key_column
< ok 0 0
> exec
| CREATE DATABASE pto_repartition_drop_old_key_column DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pto_repartition_drop_old_key_column
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY HASH (id) PARTITIONS 4
< ok 0 0
> query
| SHOW CREATE TABLE pto_repartition_drop_old_key_column.t
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
| DROP DATABASE IF EXISTS pfrom_repartition_drop_old_key_column
< ok 0 0
> exec
| CREATE DATABASE pfrom_repartition_drop_old_key_column DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pfrom_repartition_drop_old_key_column
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt)) PARTITION BY RANGE (YEAR(dt)) (PARTITION p0 VALUES LESS THAN (2010), PARTITION pmax VALUES LESS THAN MAXVALUE)
< ok 0 0
> query
| SHOW CREATE TABLE pfrom_repartition_drop_old_key_column.t
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
| DROP DATABASE IF EXISTS pa_repartition_drop_old_key_column
< ok 0 0
> exec
| CREATE DATABASE pa_repartition_drop_old_key_column DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_repartition_drop_old_key_column
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt)) PARTITION BY RANGE (YEAR(dt)) (PARTITION p0 VALUES LESS THAN (2010), PARTITION pmax VALUES LESS THAN MAXVALUE)
< ok 0 0
> exec
| ALTER TABLE `pa_repartition_drop_old_key_column`.`t` /*!50100 PARTITION BY HASH (`id`)
| (PARTITION p0 ENGINE = innodb,
|  PARTITION p1 ENGINE = innodb,
|  PARTITION p2 ENGINE = innodb,
|  PARTITION p3 ENGINE = innodb) */
< ok 0 0
> exec
| ALTER TABLE `pa_repartition_drop_old_key_column`.`t` DROP PRIMARY KEY, ADD PRIMARY KEY (`id`)
< ok 0 0
> exec
| ALTER TABLE `pa_repartition_drop_old_key_column`.`t` DROP COLUMN `dt`
< ok 0 0
> query
| SHOW CREATE TABLE pa_repartition_drop_old_key_column.t
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
| DROP DATABASE IF EXISTS pto_repartition_drop_special_col
< ok 0 0
> exec
| CREATE DATABASE pto_repartition_drop_special_col DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pto_repartition_drop_special_col
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY HASH (id) PARTITIONS 2
< ok 0 0
> query
| SHOW CREATE TABLE pto_repartition_drop_special_col.t
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
| partition `p1`)
> exec
| DROP DATABASE IF EXISTS pfrom_repartition_drop_special_col
< ok 0 0
> exec
| CREATE DATABASE pfrom_repartition_drop_special_col DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pfrom_repartition_drop_special_col
< ok 0 0
> exec
| CREATE TABLE t (`a-b` INT NOT NULL, id INT NOT NULL, PRIMARY KEY (`a-b`, id)) PARTITION BY HASH (`a-b`) PARTITIONS 4
< ok 0 0
> query
| SHOW CREATE TABLE pfrom_repartition_drop_special_col.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a-b` int(11) NOT NULL,
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`a-b`, `id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by hash(`a-b`)
| (partition `p0`,
| partition `p1`,
| partition `p2`,
| partition `p3`)
> exec
| DROP DATABASE IF EXISTS pa_repartition_drop_special_col
< ok 0 0
> exec
| CREATE DATABASE pa_repartition_drop_special_col DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_repartition_drop_special_col
< ok 0 0
> exec
| CREATE TABLE t (`a-b` INT NOT NULL, id INT NOT NULL, PRIMARY KEY (`a-b`, id)) PARTITION BY HASH (`a-b`) PARTITIONS 4
< ok 0 0
> exec
| ALTER TABLE `pa_repartition_drop_special_col`.`t` /*!50100 PARTITION BY HASH (`id`)
| (PARTITION p0 ENGINE = innodb,
|  PARTITION p1 ENGINE = innodb) */
< ok 0 0
> exec
| ALTER TABLE `pa_repartition_drop_special_col`.`t` DROP PRIMARY KEY, ADD PRIMARY KEY (`id`)
< ok 0 0
> exec
| ALTER TABLE `pa_repartition_drop_special_col`.`t` DROP COLUMN `a-b`
< ok 0 0
> query
| SHOW CREATE TABLE pa_repartition_drop_special_col.t
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
| partition `p1`)
> exec
| DROP DATABASE IF EXISTS pto_create_range_bound_arith
< ok 0 0
> exec
| CREATE DATABASE pto_create_range_bound_arith DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pto_create_range_bound_arith
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (5+5), PARTITION p1 VALUES LESS THAN (10*2), PARTITION pmax VALUES LESS THAN MAXVALUE)
< ok 0 0
> query
| SHOW CREATE TABLE pto_create_range_bound_arith.t
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
|  partition by range(`id`)
| (partition `p0` values less than (10),
| partition `p1` values less than (20),
| partition `pmax` values less than (MAXVALUE))
> exec
| DROP DATABASE IF EXISTS pa_create_range_bound_arith
< ok 0 0
> exec
| CREATE DATABASE pa_create_range_bound_arith DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_range_bound_arith
< ok 0 0
> exec
| CREATE TABLE `pa_create_range_bound_arith`.`t` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_range_bound_arith`.`t` /*!50100 PARTITION BY RANGE (`id`)
| (PARTITION p0 VALUES LESS THAN (10) ENGINE = innodb,
|  PARTITION p1 VALUES LESS THAN (20) ENGINE = innodb,
|  PARTITION pmax VALUES LESS THAN MAXVALUE ENGINE = innodb) */
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_range_bound_arith.t
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
|  partition by range(`id`)
| (partition `p0` values less than (10),
| partition `p1` values less than (20),
| partition `pmax` values less than (MAXVALUE))
> exec
| DROP DATABASE IF EXISTS pto_create_range_bound_todays
< ok 0 0
> exec
| CREATE DATABASE pto_create_range_bound_todays DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pto_create_range_bound_todays
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt)) PARTITION BY RANGE (TO_DAYS(dt)) (PARTITION p0 VALUES LESS THAN (TO_DAYS('2020-01-01')), PARTITION pmax VALUES LESS THAN MAXVALUE)
< ok 0 0
> query
| SHOW CREATE TABLE pto_create_range_bound_todays.t
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
|  partition by range(TO_DAYS(`dt`))
| (partition `p0` values less than (737790),
| partition `pmax` values less than (MAXVALUE))
> exec
| DROP DATABASE IF EXISTS pa_create_range_bound_todays
< ok 0 0
> exec
| CREATE DATABASE pa_create_range_bound_todays DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_range_bound_todays
< ok 0 0
> exec
| CREATE TABLE `pa_create_range_bound_todays`.`t` (
|   `id` int NOT NULL,
|   `dt` date NOT NULL,
|   PRIMARY KEY (`id`,`dt`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_range_bound_todays`.`t` /*!50100 PARTITION BY RANGE (to_days(`dt`))
| (PARTITION p0 VALUES LESS THAN (737790) ENGINE = innodb,
|  PARTITION pmax VALUES LESS THAN MAXVALUE ENGINE = innodb) */
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_range_bound_todays.t
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
|  partition by range(to_days(`dt`))
| (partition `p0` values less than (737790),
| partition `pmax` values less than (MAXVALUE))
