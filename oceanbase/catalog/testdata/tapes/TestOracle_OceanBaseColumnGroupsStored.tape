-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_OceanBaseColumnGroupsStored default#1
> exec
| DROP DATABASE IF EXISTS obcg_a
< ok 0 0
> exec
| CREATE DATABASE obcg_a
< ok 1 0
> exec
| USE obcg_a
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT) WITH COLUMN GROUP(all columns)
< ok 0 0
> query
| SHOW CREATE TABLE obcg_a.t
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
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 WITH COLUMN GROUP(all columns)
> exec
| DROP DATABASE IF EXISTS obcg_b
< ok 0 0
> exec
| CREATE DATABASE obcg_b
< ok 1 0
> exec
| USE obcg_b
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT) WITH COLUMN GROUP(EACH COLUMN)
< ok 0 0
> query
| SHOW CREATE TABLE obcg_b.t
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
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 WITH COLUMN GROUP(each column)
> exec
| DROP DATABASE IF EXISTS obcg_c
< ok 0 0
> exec
| CREATE DATABASE obcg_c
< ok 1 0
> exec
| USE obcg_c
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT) WITH COLUMN GROUP(each column, all columns)
< ok 0 0
> query
| SHOW CREATE TABLE obcg_c.t
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
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 WITH COLUMN GROUP(all columns, each column)
> exec
| DROP DATABASE IF EXISTS obcg_d
< ok 0 0
> exec
| CREATE DATABASE obcg_d
< ok 1 0
> exec
| USE obcg_d
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT) PARTITION BY HASH(id) PARTITIONS 2 WITH COLUMN GROUP(each column)
< ok 0 0
> query
| SHOW CREATE TABLE obcg_d.t
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
| partition `p1`) WITH COLUMN GROUP(each column)
> exec
| DROP DATABASE IF EXISTS pa_cg_create_de0a47e8_to
< ok 0 0
> exec
| CREATE DATABASE pa_cg_create_de0a47e8_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_cg_create_de0a47e8_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT) WITH COLUMN GROUP(all columns, each column)
< ok 0 0
> query
| SHOW CREATE TABLE pa_cg_create_de0a47e8_to.t
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
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 WITH COLUMN GROUP(all columns, each column)
> exec
| DROP DATABASE IF EXISTS pa_cg_create_de0a47e8
< ok 0 0
> exec
| CREATE DATABASE pa_cg_create_de0a47e8 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_cg_create_de0a47e8
< ok 0 0
> exec
| CREATE TABLE `pa_cg_create_de0a47e8`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 WITH COLUMN GROUP(all columns, each column)
< ok 0 0
> query
| SHOW CREATE TABLE pa_cg_create_de0a47e8.t
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
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 WITH COLUMN GROUP(all columns, each column)
> exec
| DROP DATABASE IF EXISTS pa_cg_add_all_45cd944d_to
< ok 0 0
> exec
| CREATE DATABASE pa_cg_add_all_45cd944d_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_cg_add_all_45cd944d_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT) WITH COLUMN GROUP(all columns)
< ok 0 0
> query
| SHOW CREATE TABLE pa_cg_add_all_45cd944d_to.t
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
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 WITH COLUMN GROUP(all columns)
> exec
| DROP DATABASE IF EXISTS pa_cg_add_all_45cd944d
< ok 0 0
> exec
| CREATE DATABASE pa_cg_add_all_45cd944d DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_cg_add_all_45cd944d
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_cg_add_all_45cd944d.t
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
| ALTER TABLE `pa_cg_add_all_45cd944d`.`t` ADD COLUMN GROUP(all columns)
< ok 0 0
> query
| SHOW CREATE TABLE pa_cg_add_all_45cd944d.t
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
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 WITH COLUMN GROUP(all columns)
> exec
| DROP DATABASE IF EXISTS pa_cg_add_both_66e574e5_to
< ok 0 0
> exec
| CREATE DATABASE pa_cg_add_both_66e574e5_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_cg_add_both_66e574e5_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT) WITH COLUMN GROUP(all columns, each column)
< ok 0 0
> query
| SHOW CREATE TABLE pa_cg_add_both_66e574e5_to.t
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
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 WITH COLUMN GROUP(all columns, each column)
> exec
| DROP DATABASE IF EXISTS pa_cg_add_both_66e574e5
< ok 0 0
> exec
| CREATE DATABASE pa_cg_add_both_66e574e5 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_cg_add_both_66e574e5
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_cg_add_both_66e574e5.t
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
| ALTER TABLE `pa_cg_add_both_66e574e5`.`t` ADD COLUMN GROUP(all columns, each column)
< ok 0 0
> query
| SHOW CREATE TABLE pa_cg_add_both_66e574e5.t
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
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 WITH COLUMN GROUP(all columns, each column)
> exec
| DROP DATABASE IF EXISTS pa_cg_all_to_each_51047dc7_to
< ok 0 0
> exec
| CREATE DATABASE pa_cg_all_to_each_51047dc7_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_cg_all_to_each_51047dc7_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT) WITH COLUMN GROUP(each column)
< ok 0 0
> query
| SHOW CREATE TABLE pa_cg_all_to_each_51047dc7_to.t
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
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 WITH COLUMN GROUP(each column)
> exec
| DROP DATABASE IF EXISTS pa_cg_all_to_each_51047dc7
< ok 0 0
> exec
| CREATE DATABASE pa_cg_all_to_each_51047dc7 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_cg_all_to_each_51047dc7
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT) WITH COLUMN GROUP(all columns)
< ok 0 0
> query
| SHOW CREATE TABLE pa_cg_all_to_each_51047dc7.t
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
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 WITH COLUMN GROUP(all columns)
> exec
| ALTER TABLE `pa_cg_all_to_each_51047dc7`.`t` DROP COLUMN GROUP(all columns)
< ok 0 0
> exec
| ALTER TABLE `pa_cg_all_to_each_51047dc7`.`t` ADD COLUMN GROUP(each column)
< ok 0 0
> query
| SHOW CREATE TABLE pa_cg_all_to_each_51047dc7.t
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
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 WITH COLUMN GROUP(each column)
> exec
| DROP DATABASE IF EXISTS pa_cg_drop_both_37663eb_to
< ok 0 0
> exec
| CREATE DATABASE pa_cg_drop_both_37663eb_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_cg_drop_both_37663eb_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_cg_drop_both_37663eb_to.t
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
| DROP DATABASE IF EXISTS pa_cg_drop_both_37663eb
< ok 0 0
> exec
| CREATE DATABASE pa_cg_drop_both_37663eb DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_cg_drop_both_37663eb
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT) WITH COLUMN GROUP(all columns, each column)
< ok 0 0
> query
| SHOW CREATE TABLE pa_cg_drop_both_37663eb.t
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
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 WITH COLUMN GROUP(all columns, each column)
> exec
| ALTER TABLE `pa_cg_drop_both_37663eb`.`t` DROP COLUMN GROUP(all columns, each column)
< ok 0 0
> query
| SHOW CREATE TABLE pa_cg_drop_both_37663eb.t
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
| DROP DATABASE IF EXISTS pa_cg_and_a_column_dcdf1058_to
< ok 0 0
> exec
| CREATE DATABASE pa_cg_and_a_column_dcdf1058_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_cg_and_a_column_dcdf1058_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT) WITH COLUMN GROUP(all columns)
< ok 0 0
> query
| SHOW CREATE TABLE pa_cg_and_a_column_dcdf1058_to.t
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
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 WITH COLUMN GROUP(all columns)
> exec
| DROP DATABASE IF EXISTS pa_cg_and_a_column_dcdf1058
< ok 0 0
> exec
| CREATE DATABASE pa_cg_and_a_column_dcdf1058 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_cg_and_a_column_dcdf1058
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT) WITH COLUMN GROUP(each column)
< ok 0 0
> query
| SHOW CREATE TABLE pa_cg_and_a_column_dcdf1058.t
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
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 WITH COLUMN GROUP(each column)
> exec
| ALTER TABLE `pa_cg_and_a_column_dcdf1058`.`t` DROP COLUMN GROUP(each column)
< ok 0 0
> exec
| ALTER TABLE `pa_cg_and_a_column_dcdf1058`.`t` ADD COLUMN GROUP(all columns)
< ok 0 0
> exec
| ALTER TABLE `pa_cg_and_a_column_dcdf1058`.`t` ADD COLUMN `b` int DEFAULT NULL
< ok 0 0
> query
| SHOW CREATE TABLE pa_cg_and_a_column_dcdf1058.t
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
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 WITH COLUMN GROUP(all columns)
