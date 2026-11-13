-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracle_PartitionApplyCorrectness default#1

== TestOracle_PartitionApplyCorrectness/8.0/add-column-keep-partitioning default#1
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

== TestOracle_PartitionApplyCorrectness/8.0/add-partitioning-hash default#1
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
| PARTITIONS 4 */
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_partitioning_hash.t
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

== TestOracle_PartitionApplyCorrectness/8.0/add-partitioning-range default#1
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

== TestOracle_PartitionApplyCorrectness/8.0/create-hash-n default#1
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
| PARTITIONS 4 */
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_hash_n.t
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

== TestOracle_PartitionApplyCorrectness/8.0/create-key default#1
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
| PARTITIONS 3 */
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
| PARTITIONS 3 */
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_key.t
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
| PARTITIONS 3 */

== TestOracle_PartitionApplyCorrectness/8.0/create-key-algorithm-2 default#1
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
< ok 0 0
> query
| SHOW CREATE TABLE pto_create_key_algorithm_2.t
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
| PARTITIONS 3 */
> exec
| DROP DATABASE IF EXISTS pa_create_key_algorithm_2
< ok 0 0
> exec
| CREATE DATABASE pa_create_key_algorithm_2 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_key_algorithm_2
< ok 0 0
> exec
| CREATE TABLE `pa_create_key_algorithm_2`.`t` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_key_algorithm_2`.`t` /*!50100 PARTITION BY KEY (id)
| PARTITIONS 3 */
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_key_algorithm_2.t
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
| PARTITIONS 3 */

== TestOracle_PartitionApplyCorrectness/8.0/create-linear-key default#1
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
< ok 0 0
> query
| SHOW CREATE TABLE pto_create_linear_key.t
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
| PARTITIONS 2 */
> exec
| DROP DATABASE IF EXISTS pa_create_linear_key
< ok 0 0
> exec
| CREATE DATABASE pa_create_linear_key DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_linear_key
< ok 0 0
> exec
| CREATE TABLE `pa_create_linear_key`.`t` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_linear_key`.`t` /*!50100 PARTITION BY LINEAR KEY (id)
| PARTITIONS 2 */
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_linear_key.t
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
| PARTITIONS 2 */

== TestOracle_PartitionApplyCorrectness/8.0/create-list default#1
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
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `c` int NOT NULL,
|   PRIMARY KEY (`id`,`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY LIST (`c`)
| (PARTITION pa VALUES IN (1,2,3) ENGINE = InnoDB,
|  PARTITION pb VALUES IN (4,5,6) ENGINE = InnoDB) */
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
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `c` int NOT NULL,
|   PRIMARY KEY (`id`,`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY LIST (`c`)
| (PARTITION pa VALUES IN (1,2,3) ENGINE = InnoDB,
|  PARTITION pb VALUES IN (4,5,6) ENGINE = InnoDB) */

== TestOracle_PartitionApplyCorrectness/8.0/create-part-comment default#1
> exec
| DROP DATABASE IF EXISTS pto_create_part_comment
< ok 0 0
> exec
| CREATE DATABASE pto_create_part_comment DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pto_create_part_comment
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, c INT NOT NULL, PRIMARY KEY (id,c)) PARTITION BY LIST (c) (PARTITION pa VALUES IN (1,2,3) COMMENT 'first', PARTITION pb VALUES IN (4,5,6) COMMENT 'second')
< ok 0 0
> query
| SHOW CREATE TABLE pto_create_part_comment.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `c` int NOT NULL,
|   PRIMARY KEY (`id`,`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY LIST (`c`)
| (PARTITION pa VALUES IN (1,2,3) COMMENT = 'first' ENGINE = InnoDB,
|  PARTITION pb VALUES IN (4,5,6) COMMENT = 'second' ENGINE = InnoDB) */
> exec
| DROP DATABASE IF EXISTS pa_create_part_comment
< ok 0 0
> exec
| CREATE DATABASE pa_create_part_comment DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_part_comment
< ok 0 0
> exec
| CREATE TABLE `pa_create_part_comment`.`t` (
|   `id` int NOT NULL,
|   `c` int NOT NULL,
|   PRIMARY KEY (`id`,`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_part_comment`.`t` /*!50100 PARTITION BY LIST (`c`)
| (PARTITION pa VALUES IN (1,2,3) ENGINE = innodb COMMENT = 'first',
|  PARTITION pb VALUES IN (4,5,6) ENGINE = innodb COMMENT = 'second') */
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_part_comment.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `c` int NOT NULL,
|   PRIMARY KEY (`id`,`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY LIST (`c`)
| (PARTITION pa VALUES IN (1,2,3) COMMENT = 'first' ENGINE = InnoDB,
|  PARTITION pb VALUES IN (4,5,6) COMMENT = 'second' ENGINE = InnoDB) */

== TestOracle_PartitionApplyCorrectness/8.0/create-range default#1
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

== TestOracle_PartitionApplyCorrectness/8.0/create-range-bound-arith default#1
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
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY RANGE (`id`)
| (PARTITION p0 VALUES LESS THAN (10) ENGINE = InnoDB,
|  PARTITION p1 VALUES LESS THAN (20) ENGINE = InnoDB,
|  PARTITION pmax VALUES LESS THAN MAXVALUE ENGINE = InnoDB) */
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
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY RANGE (`id`)
| (PARTITION p0 VALUES LESS THAN (10) ENGINE = InnoDB,
|  PARTITION p1 VALUES LESS THAN (20) ENGINE = InnoDB,
|  PARTITION pmax VALUES LESS THAN MAXVALUE ENGINE = InnoDB) */

== TestOracle_PartitionApplyCorrectness/8.0/create-range-bound-todays default#1
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
| /*!50100 PARTITION BY RANGE (to_days(`dt`))
| (PARTITION p0 VALUES LESS THAN (737790) ENGINE = InnoDB,
|  PARTITION pmax VALUES LESS THAN MAXVALUE ENGINE = InnoDB) */
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
| /*!50100 PARTITION BY RANGE (to_days(`dt`))
| (PARTITION p0 VALUES LESS THAN (737790) ENGINE = InnoDB,
|  PARTITION pmax VALUES LESS THAN MAXVALUE ENGINE = InnoDB) */

== TestOracle_PartitionApplyCorrectness/8.0/create-range-columns default#1
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
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int NOT NULL,
|   PRIMARY KEY (`id`,`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50500 PARTITION BY RANGE  COLUMNS(a)
| (PARTITION p0 VALUES LESS THAN (10) ENGINE = InnoDB,
|  PARTITION p1 VALUES LESS THAN (MAXVALUE) ENGINE = InnoDB) */
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
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int NOT NULL,
|   PRIMARY KEY (`id`,`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50500 PARTITION BY RANGE  COLUMNS(a)
| (PARTITION p0 VALUES LESS THAN (10) ENGINE = InnoDB,
|  PARTITION p1 VALUES LESS THAN (MAXVALUE) ENGINE = InnoDB) */

== TestOracle_PartitionApplyCorrectness/8.0/create-subpart-n default#1
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
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `d` date NOT NULL,
|   PRIMARY KEY (`id`,`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY RANGE (year(`d`))
| SUBPARTITION BY HASH (`id`)
| SUBPARTITIONS 2
| (PARTITION p0 VALUES LESS THAN (2010) ENGINE = InnoDB,
|  PARTITION p1 VALUES LESS THAN MAXVALUE ENGINE = InnoDB) */
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
| SUBPARTITIONS 2
| (PARTITION p0 VALUES LESS THAN (2010) ENGINE = innodb,
|  PARTITION p1 VALUES LESS THAN MAXVALUE ENGINE = innodb) */
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_subpart_n.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `d` date NOT NULL,
|   PRIMARY KEY (`id`,`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY RANGE (year(`d`))
| SUBPARTITION BY HASH (`id`)
| SUBPARTITIONS 2
| (PARTITION p0 VALUES LESS THAN (2010) ENGINE = InnoDB,
|  PARTITION p1 VALUES LESS THAN MAXVALUE ENGINE = InnoDB) */

== TestOracle_PartitionApplyCorrectness/8.0/remove-partitioning default#1
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
< ok 0 0
> query
| SHOW CREATE TABLE pa_remove_partitioning.t
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

== TestOracle_PartitionApplyCorrectness/8.0/remove-partitioning-drop-key-column default#1
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
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
< ok 0 0
> exec
| ALTER TABLE `pa_remove_partitioning_drop_key_column`.`t` DROP PRIMARY KEY
< ok 0 0
> exec
| ALTER TABLE `pa_remove_partitioning_drop_key_column`.`t` DROP COLUMN `dt`
< ok 0 0
> exec
| ALTER TABLE `pa_remove_partitioning_drop_key_column`.`t` ADD PRIMARY KEY (`id`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_remove_partitioning_drop_key_column.t
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

== TestOracle_PartitionApplyCorrectness/8.0/remove-partitioning-range default#1
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
< ok 0 0
> query
| SHOW CREATE TABLE pa_remove_partitioning_range.t
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

== TestOracle_PartitionApplyCorrectness/8.0/repartition-add-range default#1
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

== TestOracle_PartitionApplyCorrectness/8.0/repartition-count default#1
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
| PARTITIONS 8 */
< ok 0 0
> query
| SHOW CREATE TABLE pa_repartition_count.t
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

== TestOracle_PartitionApplyCorrectness/8.0/repartition-drop-old-key-column default#1
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
| ALTER TABLE `pa_repartition_drop_old_key_column`.`t` REMOVE PARTITIONING
< ok 0 0
> exec
| ALTER TABLE `pa_repartition_drop_old_key_column`.`t` DROP PRIMARY KEY
< ok 0 0
> exec
| ALTER TABLE `pa_repartition_drop_old_key_column`.`t` DROP COLUMN `dt`
< ok 0 0
> exec
| ALTER TABLE `pa_repartition_drop_old_key_column`.`t` ADD PRIMARY KEY (`id`)
< ok 0 0
> exec
| ALTER TABLE `pa_repartition_drop_old_key_column`.`t` /*!50100 PARTITION BY HASH (`id`)
| PARTITIONS 4 */
< ok 0 0
> query
| SHOW CREATE TABLE pa_repartition_drop_old_key_column.t
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

== TestOracle_PartitionApplyCorrectness/8.0/repartition-drop-special-col default#1
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
| PARTITIONS 2 */
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
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a-b` int NOT NULL,
|   `id` int NOT NULL,
|   PRIMARY KEY (`a-b`,`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY HASH (`a-b`)
| PARTITIONS 4 */
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
| ALTER TABLE `pa_repartition_drop_special_col`.`t` REMOVE PARTITIONING
< ok 0 0
> exec
| ALTER TABLE `pa_repartition_drop_special_col`.`t` DROP PRIMARY KEY
< ok 0 0
> exec
| ALTER TABLE `pa_repartition_drop_special_col`.`t` DROP COLUMN `a-b`
< ok 0 0
> exec
| ALTER TABLE `pa_repartition_drop_special_col`.`t` ADD PRIMARY KEY (`id`)
< ok 0 0
> exec
| ALTER TABLE `pa_repartition_drop_special_col`.`t` /*!50100 PARTITION BY HASH (`id`)
| PARTITIONS 2 */
< ok 0 0
> query
| SHOW CREATE TABLE pa_repartition_drop_special_col.t
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
| PARTITIONS 2 */

== TestOracle_PartitionApplyCorrectness/8.0/repartition-range-bound default#1
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

== TestOracle_PartitionApplyCorrectness/8.0/repartition-toggle-linear default#1
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
< ok 0 0
> query
| SHOW CREATE TABLE pto_repartition_toggle_linear.t
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
> exec
| DROP DATABASE IF EXISTS pfrom_repartition_toggle_linear
< ok 0 0
> exec
| CREATE DATABASE pfrom_repartition_toggle_linear DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pfrom_repartition_toggle_linear
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY KEY (id) PARTITIONS 4
< ok 0 0
> query
| SHOW CREATE TABLE pfrom_repartition_toggle_linear.t
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
| DROP DATABASE IF EXISTS pa_repartition_toggle_linear
< ok 0 0
> exec
| CREATE DATABASE pa_repartition_toggle_linear DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_repartition_toggle_linear
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY KEY (id) PARTITIONS 4
< ok 0 0
> exec
| ALTER TABLE `pa_repartition_toggle_linear`.`t` /*!50100 PARTITION BY LINEAR KEY (id)
| PARTITIONS 4 */
< ok 0 0
> query
| SHOW CREATE TABLE pa_repartition_toggle_linear.t
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

== TestOracle_PartitionApplyCorrectness/8.0/repartition-type default#1
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
| PARTITIONS 4 */
< ok 0 0
> query
| SHOW CREATE TABLE pa_repartition_type.t
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
