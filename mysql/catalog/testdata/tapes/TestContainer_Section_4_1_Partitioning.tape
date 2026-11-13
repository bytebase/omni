-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestContainer_Section_4_1_Partitioning default#1
> exec
| DROP TABLE IF EXISTS t_part_range
< ok 0 0
> exec
| CREATE TABLE t_part_range (
| 				id INT NOT NULL,
| 				created DATE NOT NULL
| 			) PARTITION BY RANGE (YEAR(created)) (
| 				PARTITION p0 VALUES LESS THAN (2020),
| 				PARTITION p1 VALUES LESS THAN (2025),
| 				PARTITION pmax VALUES LESS THAN MAXVALUE
| 			)
< ok 0 0
> query
| SHOW CREATE TABLE t_part_range
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_part_range
: bytes
| CREATE TABLE `t_part_range` (
|   `id` int NOT NULL,
|   `created` date NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY RANGE (year(`created`))
| (PARTITION p0 VALUES LESS THAN (2020) ENGINE = InnoDB,
|  PARTITION p1 VALUES LESS THAN (2025) ENGINE = InnoDB,
|  PARTITION pmax VALUES LESS THAN MAXVALUE ENGINE = InnoDB) */
> exec
| DROP TABLE IF EXISTS t_part_range_cols
< ok 0 0
> exec
| CREATE TABLE t_part_range_cols (
| 				id INT NOT NULL,
| 				city VARCHAR(50) NOT NULL,
| 				name VARCHAR(50) NOT NULL
| 			) PARTITION BY RANGE COLUMNS(city) (
| 				PARTITION p0 VALUES LESS THAN ('M'),
| 				PARTITION p1 VALUES LESS THAN MAXVALUE
| 			)
< ok 0 0
> query
| SHOW CREATE TABLE t_part_range_cols
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_part_range_cols
: bytes
| CREATE TABLE `t_part_range_cols` (
|   `id` int NOT NULL,
|   `city` varchar(50) NOT NULL,
|   `name` varchar(50) NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50500 PARTITION BY RANGE  COLUMNS(city)
| (PARTITION p0 VALUES LESS THAN ('M') ENGINE = InnoDB,
|  PARTITION p1 VALUES LESS THAN (MAXVALUE) ENGINE = InnoDB) */
> exec
| DROP TABLE IF EXISTS t_part_list
< ok 0 0
> exec
| CREATE TABLE t_part_list (
| 				id INT NOT NULL,
| 				region INT NOT NULL
| 			) PARTITION BY LIST (region) (
| 				PARTITION pNorth VALUES IN (1,2,3),
| 				PARTITION pSouth VALUES IN (4,5,6),
| 				PARTITION pWest VALUES IN (7,8,9)
| 			)
< ok 0 0
> query
| SHOW CREATE TABLE t_part_list
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_part_list
: bytes
| CREATE TABLE `t_part_list` (
|   `id` int NOT NULL,
|   `region` int NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY LIST (`region`)
| (PARTITION pNorth VALUES IN (1,2,3) ENGINE = InnoDB,
|  PARTITION pSouth VALUES IN (4,5,6) ENGINE = InnoDB,
|  PARTITION pWest VALUES IN (7,8,9) ENGINE = InnoDB) */
> exec
| DROP TABLE IF EXISTS t_part_hash
< ok 0 0
> exec
| CREATE TABLE t_part_hash (
| 				id INT NOT NULL,
| 				name VARCHAR(100)
| 			) PARTITION BY HASH (id) PARTITIONS 4
< ok 0 0
> query
| SHOW CREATE TABLE t_part_hash
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_part_hash
: bytes
| CREATE TABLE `t_part_hash` (
|   `id` int NOT NULL,
|   `name` varchar(100) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY HASH (`id`)
| PARTITIONS 4 */
> exec
| DROP TABLE IF EXISTS t_part_key
< ok 0 0
> exec
| CREATE TABLE t_part_key (
| 				id INT NOT NULL,
| 				name VARCHAR(100)
| 			) PARTITION BY KEY (id) PARTITIONS 3
< ok 0 0
> query
| SHOW CREATE TABLE t_part_key
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_part_key
: bytes
| CREATE TABLE `t_part_key` (
|   `id` int NOT NULL,
|   `name` varchar(100) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY KEY (id)
| PARTITIONS 3 */
> exec
| DROP TABLE IF EXISTS t_part_linear
< ok 0 0
> exec
| CREATE TABLE t_part_linear (
| 				id INT NOT NULL
| 			) PARTITION BY LINEAR HASH (id) PARTITIONS 4
< ok 0 0
> query
| SHOW CREATE TABLE t_part_linear
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_part_linear
: bytes
| CREATE TABLE `t_part_linear` (
|   `id` int NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY LINEAR HASH (`id`)
| PARTITIONS 4 */
> exec
| DROP TABLE IF EXISTS t_part_sub
< ok 0 0
> exec
| CREATE TABLE t_part_sub (
| 				id INT NOT NULL,
| 				purchased DATE NOT NULL
| 			) PARTITION BY RANGE (YEAR(purchased))
| 			SUBPARTITION BY HASH (id)
| 			SUBPARTITIONS 2 (
| 				PARTITION p0 VALUES LESS THAN (2020),
| 				PARTITION p1 VALUES LESS THAN MAXVALUE
| 			)
< ok 0 0
> query
| SHOW CREATE TABLE t_part_sub
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_part_sub
: bytes
| CREATE TABLE `t_part_sub` (
|   `id` int NOT NULL,
|   `purchased` date NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY RANGE (year(`purchased`))
| SUBPARTITION BY HASH (`id`)
| SUBPARTITIONS 2
| (PARTITION p0 VALUES LESS THAN (2020) ENGINE = InnoDB,
|  PARTITION p1 VALUES LESS THAN MAXVALUE ENGINE = InnoDB) */
> exec
| DROP TABLE IF EXISTS t_part_show
< ok 0 0
> exec
| CREATE TABLE t_part_show (
| 				id INT NOT NULL,
| 				val INT NOT NULL
| 			) PARTITION BY RANGE (val) (
| 				PARTITION p0 VALUES LESS THAN (10),
| 				PARTITION p1 VALUES LESS THAN (20),
| 				PARTITION p2 VALUES LESS THAN MAXVALUE
| 			)
< ok 0 0
> query
| SHOW CREATE TABLE t_part_show
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_part_show
: bytes
| CREATE TABLE `t_part_show` (
|   `id` int NOT NULL,
|   `val` int NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY RANGE (`val`)
| (PARTITION p0 VALUES LESS THAN (10) ENGINE = InnoDB,
|  PARTITION p1 VALUES LESS THAN (20) ENGINE = InnoDB,
|  PARTITION p2 VALUES LESS THAN MAXVALUE ENGINE = InnoDB) */
> exec
| DROP TABLE IF EXISTS t_alter_addp
< ok 0 0
> exec
| CREATE TABLE t_alter_addp (
| 			id INT NOT NULL,
| 			val INT NOT NULL
| 		) PARTITION BY RANGE (val) (
| 			PARTITION p0 VALUES LESS THAN (10),
| 			PARTITION p1 VALUES LESS THAN (20)
| 		)
< ok 0 0
> exec
| ALTER TABLE t_alter_addp ADD PARTITION (PARTITION p2 VALUES LESS THAN (30))
< ok 0 0
> query
| SHOW CREATE TABLE t_alter_addp
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_alter_addp
: bytes
| CREATE TABLE `t_alter_addp` (
|   `id` int NOT NULL,
|   `val` int NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY RANGE (`val`)
| (PARTITION p0 VALUES LESS THAN (10) ENGINE = InnoDB,
|  PARTITION p1 VALUES LESS THAN (20) ENGINE = InnoDB,
|  PARTITION p2 VALUES LESS THAN (30) ENGINE = InnoDB) */
> exec
| DROP TABLE IF EXISTS t_alter_dropp
< ok 0 0
> exec
| CREATE TABLE t_alter_dropp (
| 			id INT NOT NULL,
| 			val INT NOT NULL
| 		) PARTITION BY RANGE (val) (
| 			PARTITION p0 VALUES LESS THAN (10),
| 			PARTITION p1 VALUES LESS THAN (20),
| 			PARTITION p2 VALUES LESS THAN MAXVALUE
| 		)
< ok 0 0
> exec
| ALTER TABLE t_alter_dropp DROP PARTITION p1
< ok 0 0
> query
| SHOW CREATE TABLE t_alter_dropp
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_alter_dropp
: bytes
| CREATE TABLE `t_alter_dropp` (
|   `id` int NOT NULL,
|   `val` int NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY RANGE (`val`)
| (PARTITION p0 VALUES LESS THAN (10) ENGINE = InnoDB,
|  PARTITION p2 VALUES LESS THAN MAXVALUE ENGINE = InnoDB) */
> exec
| DROP TABLE IF EXISTS t_alter_reorgp
< ok 0 0
> exec
| CREATE TABLE t_alter_reorgp (
| 			id INT NOT NULL,
| 			val INT NOT NULL
| 		) PARTITION BY RANGE (val) (
| 			PARTITION p0 VALUES LESS THAN (10),
| 			PARTITION p1 VALUES LESS THAN (30)
| 		)
< ok 0 0
> exec
| ALTER TABLE t_alter_reorgp REORGANIZE PARTITION p1 INTO (
| 			PARTITION p1a VALUES LESS THAN (20),
| 			PARTITION p1b VALUES LESS THAN (30)
| 		)
< ok 0 0
> query
| SHOW CREATE TABLE t_alter_reorgp
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_alter_reorgp
: bytes
| CREATE TABLE `t_alter_reorgp` (
|   `id` int NOT NULL,
|   `val` int NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY RANGE (`val`)
| (PARTITION p0 VALUES LESS THAN (10) ENGINE = InnoDB,
|  PARTITION p1a VALUES LESS THAN (20) ENGINE = InnoDB,
|  PARTITION p1b VALUES LESS THAN (30) ENGINE = InnoDB) */
> exec
| DROP TABLE IF EXISTS t_alter_truncp
< ok 0 0
> exec
| CREATE TABLE t_alter_truncp (
| 			id INT NOT NULL,
| 			val INT NOT NULL
| 		) PARTITION BY RANGE (val) (
| 			PARTITION p0 VALUES LESS THAN (10),
| 			PARTITION p1 VALUES LESS THAN MAXVALUE
| 		)
< ok 0 0
> exec
| ALTER TABLE t_alter_truncp TRUNCATE PARTITION p0
< ok 0 0
> query
| SHOW CREATE TABLE t_alter_truncp
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_alter_truncp
: bytes
| CREATE TABLE `t_alter_truncp` (
|   `id` int NOT NULL,
|   `val` int NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY RANGE (`val`)
| (PARTITION p0 VALUES LESS THAN (10) ENGINE = InnoDB,
|  PARTITION p1 VALUES LESS THAN MAXVALUE ENGINE = InnoDB) */
> exec
| DROP TABLE IF EXISTS t_alter_coalp
< ok 0 0
> exec
| CREATE TABLE t_alter_coalp (
| 			id INT NOT NULL
| 		) PARTITION BY HASH (id) PARTITIONS 4
< ok 0 0
> exec
| ALTER TABLE t_alter_coalp COALESCE PARTITION 2
< ok 0 0
> query
| SHOW CREATE TABLE t_alter_coalp
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_alter_coalp
: bytes
| CREATE TABLE `t_alter_coalp` (
|   `id` int NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY HASH (`id`)
| PARTITIONS 2 */
> exec
| DROP TABLE IF EXISTS t_alter_exchp
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_alter_exchp_swap
< ok 0 0
> exec
| CREATE TABLE t_alter_exchp (
| 			id INT NOT NULL,
| 			val INT NOT NULL
| 		) PARTITION BY RANGE (val) (
| 			PARTITION p0 VALUES LESS THAN (10),
| 			PARTITION p1 VALUES LESS THAN MAXVALUE
| 		)
< ok 0 0
> exec
| CREATE TABLE t_alter_exchp_swap (
| 			id INT NOT NULL,
| 			val INT NOT NULL
| 		)
< ok 0 0
> exec
| ALTER TABLE t_alter_exchp EXCHANGE PARTITION p0 WITH TABLE t_alter_exchp_swap
< ok 0 0
> query
| SHOW CREATE TABLE t_alter_exchp
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_alter_exchp
: bytes
| CREATE TABLE `t_alter_exchp` (
|   `id` int NOT NULL,
|   `val` int NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY RANGE (`val`)
| (PARTITION p0 VALUES LESS THAN (10) ENGINE = InnoDB,
|  PARTITION p1 VALUES LESS THAN MAXVALUE ENGINE = InnoDB) */
