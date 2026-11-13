-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

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
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY HASH (`id`)
| PARTITIONS 4 */
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
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY HASH (`id`)
| PARTITIONS 4 */
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
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY HASH (`id`)
| PARTITIONS 4 */
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
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY HASH (`id`)
| PARTITIONS 8 */
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
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY HASH (`id`)
| PARTITIONS 4 */
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
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY KEY (id)
| PARTITIONS 4 */
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
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `dt` date NOT NULL,
|   PRIMARY KEY (`id`,`dt`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY RANGE (year(`dt`))
| (PARTITION p0 VALUES LESS THAN (2010) ENGINE = InnoDB,
|  PARTITION pmax VALUES LESS THAN MAXVALUE ENGINE = InnoDB) */
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
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `dt` date NOT NULL,
|   PRIMARY KEY (`id`,`dt`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY RANGE (year(`dt`))
| (PARTITION p0 VALUES LESS THAN (2015) ENGINE = InnoDB,
|  PARTITION pmax VALUES LESS THAN MAXVALUE ENGINE = InnoDB) */
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
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `dt` date NOT NULL,
|   PRIMARY KEY (`id`,`dt`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY RANGE (year(`dt`))
| (PARTITION p0 VALUES LESS THAN (2010) ENGINE = InnoDB,
|  PARTITION pmax VALUES LESS THAN MAXVALUE ENGINE = InnoDB) */
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
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `dt` date NOT NULL,
|   PRIMARY KEY (`id`,`dt`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY RANGE (year(`dt`))
| (PARTITION p0 VALUES LESS THAN (2010) ENGINE = InnoDB,
|  PARTITION p1 VALUES LESS THAN (2020) ENGINE = InnoDB,
|  PARTITION pmax VALUES LESS THAN MAXVALUE ENGINE = InnoDB) */
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
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY KEY (id)
| PARTITIONS 4 */
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
< ok 0 0
> query
| SHOW CREATE TABLE pc_to_toggle_linear.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY LINEAR KEY (id)
| PARTITIONS 4 */
