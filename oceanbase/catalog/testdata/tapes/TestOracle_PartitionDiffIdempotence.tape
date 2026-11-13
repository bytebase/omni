-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_PartitionDiffIdempotence default#1
> exec
| DROP DATABASE IF EXISTS part_range_expr
< ok 0 0
> exec
| CREATE DATABASE part_range_expr
< ok 1 0
> exec
| USE part_range_expr
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt)) PARTITION BY RANGE (YEAR(dt)) (PARTITION p0 VALUES LESS THAN (2010), PARTITION p1 VALUES LESS THAN (2020), PARTITION pmax VALUES LESS THAN MAXVALUE)
< ok 0 0
> query
| SHOW CREATE TABLE part_range_expr.t
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
| DROP DATABASE IF EXISTS part_range_columns
< ok 0 0
> exec
| CREATE DATABASE part_range_columns
< ok 1 0
> exec
| USE part_range_columns
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, a INT NOT NULL, PRIMARY KEY (id,a)) PARTITION BY RANGE COLUMNS(a) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (MAXVALUE))
< ok 0 0
> query
| SHOW CREATE TABLE part_range_columns.t
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
| DROP DATABASE IF EXISTS part_range_columns_multi
< ok 0 0
> exec
| CREATE DATABASE part_range_columns_multi
< ok 1 0
> exec
| USE part_range_columns_multi
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, PRIMARY KEY (a,b)) PARTITION BY RANGE COLUMNS(a,b) (PARTITION p0 VALUES LESS THAN (10,10), PARTITION p1 VALUES LESS THAN (MAXVALUE,MAXVALUE))
< ok 0 0
> query
| SHOW CREATE TABLE part_range_columns_multi.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) NOT NULL,
|   `b` int(11) NOT NULL,
|   PRIMARY KEY (`a`, `b`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by range columns(`a`, `b`)
| (partition `p0` values less than (10,10),
| partition `p1` values less than (MAXVALUE,MAXVALUE))
> exec
| DROP DATABASE IF EXISTS part_list_expr
< ok 0 0
> exec
| CREATE DATABASE part_list_expr
< ok 1 0
> exec
| USE part_list_expr
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, c INT NOT NULL, PRIMARY KEY (id,c)) PARTITION BY LIST (c) (PARTITION pa VALUES IN (1,2,3), PARTITION pb VALUES IN (4,5,6))
< ok 0 0
> query
| SHOW CREATE TABLE part_list_expr.t
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
| DROP DATABASE IF EXISTS part_list_columns_multi
< ok 0 0
> exec
| CREATE DATABASE part_list_columns_multi
< ok 1 0
> exec
| USE part_list_columns_multi
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, PRIMARY KEY (a,b)) PARTITION BY LIST COLUMNS(a,b) (PARTITION p0 VALUES IN ((1,1),(2,2)), PARTITION p1 VALUES IN ((3,3),(4,4)))
< ok 0 0
> query
| SHOW CREATE TABLE part_list_columns_multi.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) NOT NULL,
|   `b` int(11) NOT NULL,
|   PRIMARY KEY (`a`, `b`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by list columns(`a`, `b`)
| (partition `p0` values in ((1,1),(2,2)),
| partition `p1` values in ((3,3),(4,4)))
> exec
| DROP DATABASE IF EXISTS part_hash_n
< ok 0 0
> exec
| CREATE DATABASE part_hash_n
< ok 1 0
> exec
| USE part_hash_n
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT) PARTITION BY HASH (id) PARTITIONS 4
< ok 0 0
> query
| SHOW CREATE TABLE part_hash_n.t
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
| DROP DATABASE IF EXISTS part_hash_explicit
< ok 0 0
> exec
| CREATE DATABASE part_hash_explicit
< ok 1 0
> exec
| USE part_hash_explicit
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY HASH (id) (PARTITION pa, PARTITION pb, PARTITION pc)
< ok 0 0
> query
| SHOW CREATE TABLE part_hash_explicit.t
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
| (partition `pa`,
| partition `pb`,
| partition `pc`)
> exec
| DROP DATABASE IF EXISTS part_linear_hash
< ok 0 0
> exec
| CREATE DATABASE part_linear_hash
< ok 1 0
> exec
| USE part_linear_hash
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY LINEAR HASH (id) PARTITIONS 4
< error
| 1064 42000 …near 'LINEAR HASH (id) PARTITIONS 4' at line 1
> exec
| DROP DATABASE IF EXISTS part_key_implicit
< ok 0 0
> exec
| CREATE DATABASE part_key_implicit
< ok 1 0
> exec
| USE part_key_implicit
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT) PARTITION BY KEY () PARTITIONS 3
< ok 0 0
> query
| SHOW CREATE TABLE part_key_implicit.t
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
|  partition by key()
| (partition `p0`,
| partition `p1`,
| partition `p2`)
> exec
| DROP DATABASE IF EXISTS part_key_explicit
< ok 0 0
> exec
| CREATE DATABASE part_key_explicit
< ok 1 0
> exec
| USE part_key_explicit
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT) PARTITION BY KEY (id) PARTITIONS 3
< ok 0 0
> query
| SHOW CREATE TABLE part_key_explicit.t
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
|  partition by key(`id`)
| (partition `p0`,
| partition `p1`,
| partition `p2`)
> exec
| DROP DATABASE IF EXISTS part_key_algorithm_2
< ok 0 0
> exec
| CREATE DATABASE part_key_algorithm_2
< ok 1 0
> exec
| USE part_key_algorithm_2
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT) PARTITION BY KEY ALGORITHM=2 (id) PARTITIONS 3
< error
| 1064 42000 …near 'ALGORITHM=2 (id) PARTITIONS 3' at line 1
> exec
| DROP DATABASE IF EXISTS part_key_algorithm_1
< ok 0 0
> exec
| CREATE DATABASE part_key_algorithm_1
< ok 1 0
> exec
| USE part_key_algorithm_1
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT) PARTITION BY KEY ALGORITHM=1 (id) PARTITIONS 3
< error
| 1064 42000 …near 'ALGORITHM=1 (id) PARTITIONS 3' at line 1
> exec
| DROP DATABASE IF EXISTS part_linear_key
< ok 0 0
> exec
| CREATE DATABASE part_linear_key
< ok 1 0
> exec
| USE part_linear_key
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY LINEAR KEY (id) PARTITIONS 2
< error
| 1064 42000 …near 'LINEAR KEY (id) PARTITIONS 2' at line 1
> exec
| DROP DATABASE IF EXISTS part_subpart_n
< ok 0 0
> exec
| CREATE DATABASE part_subpart_n
< ok 1 0
> exec
| USE part_subpart_n
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, d DATE NOT NULL, PRIMARY KEY (id,d)) PARTITION BY RANGE (YEAR(d)) SUBPARTITION BY HASH (id) SUBPARTITIONS 2 (PARTITION p0 VALUES LESS THAN (2010), PARTITION p1 VALUES LESS THAN MAXVALUE)
< ok 0 0
> query
| SHOW CREATE TABLE part_subpart_n.t
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
| DROP DATABASE IF EXISTS part_subpart_explicit
< ok 0 0
> exec
| CREATE DATABASE part_subpart_explicit
< ok 1 0
> exec
| USE part_subpart_explicit
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, d DATE NOT NULL, PRIMARY KEY (id,d)) PARTITION BY RANGE (YEAR(d)) SUBPARTITION BY HASH (id) (PARTITION p0 VALUES LESS THAN (2010) (SUBPARTITION s0a, SUBPARTITION s0b), PARTITION p1 VALUES LESS THAN MAXVALUE (SUBPARTITION s1a, SUBPARTITION s1b))
< ok 0 0
> query
| SHOW CREATE TABLE part_subpart_explicit.t
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
| subpartition `s0a`,
| subpartition `s0b`),
| partition `p1` values less than (MAXVALUE) (
| subpartition `s1a`,
| subpartition `s1b`))
> exec
| DROP DATABASE IF EXISTS part_subpart_key
< ok 0 0
> exec
| CREATE DATABASE part_subpart_key
< ok 1 0
> exec
| USE part_subpart_key
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, d DATE NOT NULL, PRIMARY KEY (id,d)) PARTITION BY RANGE (YEAR(d)) SUBPARTITION BY KEY (id) SUBPARTITIONS 2 (PARTITION p0 VALUES LESS THAN (2010), PARTITION p1 VALUES LESS THAN MAXVALUE)
< ok 0 0
> query
| SHOW CREATE TABLE part_subpart_key.t
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
|  partition by range(year(`d`)) subpartition by key(`id`) subpartition template (
| subpartition `p0`,
| subpartition `p1`)
| (partition `p0` values less than (2010),
| partition `p1` values less than (MAXVALUE))
> exec
| DROP DATABASE IF EXISTS part_range_bound_arith
< ok 0 0
> exec
| CREATE DATABASE part_range_bound_arith
< ok 1 0
> exec
| USE part_range_bound_arith
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (5+5), PARTITION p1 VALUES LESS THAN (10*2), PARTITION pmax VALUES LESS THAN MAXVALUE)
< ok 0 0
> query
| SHOW CREATE TABLE part_range_bound_arith.t
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
| DROP DATABASE IF EXISTS part_range_bound_divmod
< ok 0 0
> exec
| CREATE DATABASE part_range_bound_divmod
< ok 1 0
> exec
| USE part_range_bound_divmod
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (100 DIV 3), PARTITION p1 VALUES LESS THAN ((20+5)*4), PARTITION pmax VALUES LESS THAN MAXVALUE)
< ok 0 0
> query
| SHOW CREATE TABLE part_range_bound_divmod.t
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
| (partition `p0` values less than (33),
| partition `p1` values less than (100),
| partition `pmax` values less than (MAXVALUE))
> exec
| DROP DATABASE IF EXISTS part_range_bound_todays
< ok 0 0
> exec
| CREATE DATABASE part_range_bound_todays
< ok 1 0
> exec
| USE part_range_bound_todays
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt)) PARTITION BY RANGE (TO_DAYS(dt)) (PARTITION p0 VALUES LESS THAN (TO_DAYS('2020-01-01')), PARTITION p1 VALUES LESS THAN (TO_DAYS('2021-01-01')), PARTITION pmax VALUES LESS THAN MAXVALUE)
< ok 0 0
> query
| SHOW CREATE TABLE part_range_bound_todays.t
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
| partition `p1` values less than (738156),
| partition `pmax` values less than (MAXVALUE))
> exec
| DROP DATABASE IF EXISTS part_range_bound_year_fn
< ok 0 0
> exec
| CREATE DATABASE part_range_bound_year_fn
< ok 1 0
> exec
| USE part_range_bound_year_fn
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt)) PARTITION BY RANGE (YEAR(dt)) (PARTITION p0 VALUES LESS THAN (YEAR('2020-06-15')), PARTITION pmax VALUES LESS THAN MAXVALUE)
< ok 0 0
> query
| SHOW CREATE TABLE part_range_bound_year_fn.t
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
| (partition `p0` values less than (2020),
| partition `pmax` values less than (MAXVALUE))
> exec
| DROP DATABASE IF EXISTS part_list_bound_arith
< ok 0 0
> exec
| CREATE DATABASE part_list_bound_arith
< ok 1 0
> exec
| USE part_list_bound_arith
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY LIST (id) (PARTITION p0 VALUES IN (1+1, 2+2), PARTITION p1 VALUES IN (10, 20))
< ok 0 0
> query
| SHOW CREATE TABLE part_list_bound_arith.t
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
|  partition by list(`id`)
| (partition `p0` values in (2,4),
| partition `p1` values in (10,20))
> exec
| DROP DATABASE IF EXISTS part_range_columns_bound_arith
< ok 0 0
> exec
| CREATE DATABASE part_range_columns_bound_arith
< ok 1 0
> exec
| USE part_range_columns_bound_arith
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b VARCHAR(10) NOT NULL, PRIMARY KEY (a,b)) PARTITION BY RANGE COLUMNS(a,b) (PARTITION p0 VALUES LESS THAN (5+5, 'm'), PARTITION p1 VALUES LESS THAN (MAXVALUE, MAXVALUE))
< ok 0 0
> query
| SHOW CREATE TABLE part_range_columns_bound_arith.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) NOT NULL,
|   `b` varchar(10) NOT NULL,
|   PRIMARY KEY (`a`, `b`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by range columns(`a`, `b`)
| (partition `p0` values less than (10,'m'),
| partition `p1` values less than (MAXVALUE,MAXVALUE))
> exec
| DROP DATABASE IF EXISTS part_range_bound_literal_guard
< ok 0 0
> exec
| CREATE DATABASE part_range_bound_literal_guard
< ok 1 0
> exec
| USE part_range_bound_literal_guard
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (-10), PARTITION p1 VALUES LESS THAN (0), PARTITION p2 VALUES LESS THAN (100), PARTITION pmax VALUES LESS THAN MAXVALUE)
< ok 0 0
> query
| SHOW CREATE TABLE part_range_bound_literal_guard.t
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
| (partition `p0` values less than (-10),
| partition `p1` values less than (0),
| partition `p2` values less than (100),
| partition `pmax` values less than (MAXVALUE))
