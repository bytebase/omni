-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

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
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int NOT NULL,
|   `b` int NOT NULL,
|   PRIMARY KEY (`a`,`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50500 PARTITION BY RANGE  COLUMNS(a,b)
| (PARTITION p0 VALUES LESS THAN (10,10) ENGINE = InnoDB,
|  PARTITION p1 VALUES LESS THAN (MAXVALUE,MAXVALUE) ENGINE = InnoDB) */
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
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int NOT NULL,
|   `b` int NOT NULL,
|   PRIMARY KEY (`a`,`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50500 PARTITION BY LIST  COLUMNS(a,b)
| (PARTITION p0 VALUES IN ((1,1),(2,2)) ENGINE = InnoDB,
|  PARTITION p1 VALUES IN ((3,3),(4,4)) ENGINE = InnoDB) */
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
| (PARTITION pa ENGINE = InnoDB,
|  PARTITION pb ENGINE = InnoDB,
|  PARTITION pc ENGINE = InnoDB) */
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
< ok 0 0
> query
| SHOW CREATE TABLE part_linear_hash.t
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
| /*!50100 PARTITION BY LINEAR HASH (`id`)
| PARTITIONS 4 */
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
| /*!50100 PARTITION BY KEY ()
| PARTITIONS 3 */
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
| /*!50100 PARTITION BY KEY (id)
| PARTITIONS 3 */
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
< ok 0 0
> query
| SHOW CREATE TABLE part_key_algorithm_2.t
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
| /*!50100 PARTITION BY KEY (id)
| PARTITIONS 3 */
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
< ok 0 0
> query
| SHOW CREATE TABLE part_key_algorithm_1.t
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
| /*!50100 PARTITION BY KEY */ /*!50611 ALGORITHM = 1 */ /*!50100 (id)
| PARTITIONS 3 */
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
< ok 0 0
> query
| SHOW CREATE TABLE part_linear_key.t
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
| DROP DATABASE IF EXISTS part_part_comment
< ok 0 0
> exec
| CREATE DATABASE part_part_comment
< ok 1 0
> exec
| USE part_part_comment
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, c INT NOT NULL, PRIMARY KEY (id,c)) PARTITION BY LIST (c) (PARTITION pa VALUES IN (1,2,3) COMMENT 'first', PARTITION pb VALUES IN (4,5,6) COMMENT 'second')
< ok 0 0
> query
| SHOW CREATE TABLE part_part_comment.t
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
| (PARTITION p0 VALUES LESS THAN (2010)
|  (SUBPARTITION s0a ENGINE = InnoDB,
|   SUBPARTITION s0b ENGINE = InnoDB),
|  PARTITION p1 VALUES LESS THAN MAXVALUE
|  (SUBPARTITION s1a ENGINE = InnoDB,
|   SUBPARTITION s1b ENGINE = InnoDB)) */
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
| SUBPARTITION BY KEY (id)
| SUBPARTITIONS 2
| (PARTITION p0 VALUES LESS THAN (2010) ENGINE = InnoDB,
|  PARTITION p1 VALUES LESS THAN MAXVALUE ENGINE = InnoDB) */
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
| (PARTITION p0 VALUES LESS THAN (33) ENGINE = InnoDB,
|  PARTITION p1 VALUES LESS THAN (100) ENGINE = InnoDB,
|  PARTITION pmax VALUES LESS THAN MAXVALUE ENGINE = InnoDB) */
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
|  PARTITION p1 VALUES LESS THAN (738156) ENGINE = InnoDB,
|  PARTITION pmax VALUES LESS THAN MAXVALUE ENGINE = InnoDB) */
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
| (PARTITION p0 VALUES LESS THAN (2020) ENGINE = InnoDB,
|  PARTITION pmax VALUES LESS THAN MAXVALUE ENGINE = InnoDB) */
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
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY LIST (`id`)
| (PARTITION p0 VALUES IN (2,4) ENGINE = InnoDB,
|  PARTITION p1 VALUES IN (10,20) ENGINE = InnoDB) */
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
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int NOT NULL,
|   `b` varchar(10) NOT NULL,
|   PRIMARY KEY (`a`,`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50500 PARTITION BY RANGE  COLUMNS(a,b)
| (PARTITION p0 VALUES LESS THAN (10,'m') ENGINE = InnoDB,
|  PARTITION p1 VALUES LESS THAN (MAXVALUE,MAXVALUE) ENGINE = InnoDB) */
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
| (PARTITION p0 VALUES LESS THAN (-10) ENGINE = InnoDB,
|  PARTITION p1 VALUES LESS THAN (0) ENGINE = InnoDB,
|  PARTITION p2 VALUES LESS THAN (100) ENGINE = InnoDB,
|  PARTITION pmax VALUES LESS THAN MAXVALUE ENGINE = InnoDB) */
