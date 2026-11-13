-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracle_IndexMigrationApplyCorrectness default#1

== TestOracle_IndexMigrationApplyCorrectness/8.0/add-column-and-index-on-it default#1
> exec
| DROP DATABASE IF EXISTS pa_add_column_and_index_on_it_1a536244_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_column_and_index_on_it_1a536244_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_column_and_index_on_it_1a536244_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ka (a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_column_and_index_on_it_1a536244_to.t
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
|   PRIMARY KEY (`id`),
|   KEY `ka` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_add_column_and_index_on_it_1a536244_from
< ok 0 0
> exec
| CREATE DATABASE pa_add_column_and_index_on_it_1a536244_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_column_and_index_on_it_1a536244_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_column_and_index_on_it_1a536244_from.t
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
| DROP DATABASE IF EXISTS pa_add_column_and_index_on_it_1a536244
< ok 0 0
> exec
| CREATE DATABASE pa_add_column_and_index_on_it_1a536244 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_column_and_index_on_it_1a536244
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE `pa_add_column_and_index_on_it_1a536244`.`t` ADD COLUMN `a` int DEFAULT NULL
< ok 0 0
> exec
| ALTER TABLE `pa_add_column_and_index_on_it_1a536244`.`t` ADD KEY `ka` (`a`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_column_and_index_on_it_1a536244.t
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
|   PRIMARY KEY (`id`),
|   KEY `ka` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_IndexMigrationApplyCorrectness/8.0/add-comment default#1
> exec
| DROP DATABASE IF EXISTS pa_add_comment_90bb330c_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_comment_90bb330c_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_comment_90bb330c_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ka (a) COMMENT 'c')
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_comment_90bb330c_to.t
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
|   PRIMARY KEY (`id`),
|   KEY `ka` (`a`) COMMENT 'c'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_add_comment_90bb330c_from
< ok 0 0
> exec
| CREATE DATABASE pa_add_comment_90bb330c_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_comment_90bb330c_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_comment_90bb330c_from.t
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
| DROP DATABASE IF EXISTS pa_add_comment_90bb330c
< ok 0 0
> exec
| CREATE DATABASE pa_add_comment_90bb330c DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_comment_90bb330c
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> exec
| ALTER TABLE `pa_add_comment_90bb330c`.`t` ADD KEY `ka` (`a`) COMMENT 'c'
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_comment_90bb330c.t
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
|   PRIMARY KEY (`id`),
|   KEY `ka` (`a`) COMMENT 'c'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_IndexMigrationApplyCorrectness/8.0/add-composite default#1
> exec
| DROP DATABASE IF EXISTS pa_add_composite_c4f9a1fe_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_composite_c4f9a1fe_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_composite_c4f9a1fe_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY kab (a,b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_composite_c4f9a1fe_to.t
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
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `kab` (`a`,`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_add_composite_c4f9a1fe_from
< ok 0 0
> exec
| CREATE DATABASE pa_add_composite_c4f9a1fe_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_composite_c4f9a1fe_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_composite_c4f9a1fe_from.t
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
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_add_composite_c4f9a1fe
< ok 0 0
> exec
| CREATE DATABASE pa_add_composite_c4f9a1fe DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_composite_c4f9a1fe
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT)
< ok 0 0
> exec
| ALTER TABLE `pa_add_composite_c4f9a1fe`.`t` ADD KEY `kab` (`a`,`b`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_composite_c4f9a1fe.t
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
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `kab` (`a`,`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_IndexMigrationApplyCorrectness/8.0/add-descending default#1
> exec
| DROP DATABASE IF EXISTS pa_add_descending_daaa5a5_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_descending_daaa5a5_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_descending_daaa5a5_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY ab (a, b DESC))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_descending_daaa5a5_to.t
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
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `ab` (`a`,`b` DESC)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_add_descending_daaa5a5_from
< ok 0 0
> exec
| CREATE DATABASE pa_add_descending_daaa5a5_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_descending_daaa5a5_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_descending_daaa5a5_from.t
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
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_add_descending_daaa5a5
< ok 0 0
> exec
| CREATE DATABASE pa_add_descending_daaa5a5 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_descending_daaa5a5
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT)
< ok 0 0
> exec
| ALTER TABLE `pa_add_descending_daaa5a5`.`t` ADD KEY `ab` (`a`,`b` DESC)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_descending_daaa5a5.t
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
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `ab` (`a`,`b` DESC)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_IndexMigrationApplyCorrectness/8.0/add-fulltext default#1
> exec
| DROP DATABASE IF EXISTS pa_add_fulltext_db344b55_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_fulltext_db344b55_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_fulltext_db344b55_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, body TEXT, FULLTEXT KEY ftb (body))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_fulltext_db344b55_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `body` text,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `ftb` (`body`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_add_fulltext_db344b55_from
< ok 0 0
> exec
| CREATE DATABASE pa_add_fulltext_db344b55_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_fulltext_db344b55_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, body TEXT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_fulltext_db344b55_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `body` text,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_add_fulltext_db344b55
< ok 0 0
> exec
| CREATE DATABASE pa_add_fulltext_db344b55 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_fulltext_db344b55
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, body TEXT)
< ok 0 0
> exec
| ALTER TABLE `pa_add_fulltext_db344b55`.`t` ADD FULLTEXT KEY `ftb` (`body`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_fulltext_db344b55.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `body` text,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `ftb` (`body`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_IndexMigrationApplyCorrectness/8.0/add-functional default#1
> exec
| DROP DATABASE IF EXISTS pa_add_functional_fe148686_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_functional_fe148686_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_functional_fe148686_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY fa ((a + 1)))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_functional_fe148686_to.t
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
|   PRIMARY KEY (`id`),
|   KEY `fa` (((`a` + 1)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_add_functional_fe148686_from
< ok 0 0
> exec
| CREATE DATABASE pa_add_functional_fe148686_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_functional_fe148686_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_functional_fe148686_from.t
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
| DROP DATABASE IF EXISTS pa_add_functional_fe148686
< ok 0 0
> exec
| CREATE DATABASE pa_add_functional_fe148686 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_functional_fe148686
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> exec
| ALTER TABLE `pa_add_functional_fe148686`.`t` ADD KEY `fa` (((`a` + 1)))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_functional_fe148686.t
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
|   PRIMARY KEY (`id`),
|   KEY `fa` (((`a` + 1)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_IndexMigrationApplyCorrectness/8.0/add-invisible default#1
> exec
| DROP DATABASE IF EXISTS pa_add_invisible_26b110a_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_invisible_26b110a_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_invisible_26b110a_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY va (a) /*!80000 INVISIBLE */)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_invisible_26b110a_to.t
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
|   PRIMARY KEY (`id`),
|   KEY `va` (`a`) /*!80000 INVISIBLE */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_add_invisible_26b110a_from
< ok 0 0
> exec
| CREATE DATABASE pa_add_invisible_26b110a_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_invisible_26b110a_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_invisible_26b110a_from.t
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
| DROP DATABASE IF EXISTS pa_add_invisible_26b110a
< ok 0 0
> exec
| CREATE DATABASE pa_add_invisible_26b110a DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_invisible_26b110a
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> exec
| ALTER TABLE `pa_add_invisible_26b110a`.`t` ADD KEY `va` (`a`) /*!80000 INVISIBLE */
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_invisible_26b110a.t
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
|   PRIMARY KEY (`id`),
|   KEY `va` (`a`) /*!80000 INVISIBLE */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_IndexMigrationApplyCorrectness/8.0/add-mv-index default#1
> exec
| DROP DATABASE IF EXISTS pa_add_mv_index_1429d879_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_mv_index_1429d879_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_mv_index_1429d879_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, j JSON, KEY mv ((CAST(j->'$.ids' AS UNSIGNED ARRAY))))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_mv_index_1429d879_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `j` json DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `mv` ((cast(json_extract(`j`,_utf8mb4'$.ids') as unsigned array)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_add_mv_index_1429d879_from
< ok 0 0
> exec
| CREATE DATABASE pa_add_mv_index_1429d879_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_mv_index_1429d879_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, j JSON)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_mv_index_1429d879_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `j` json DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_add_mv_index_1429d879
< ok 0 0
> exec
| CREATE DATABASE pa_add_mv_index_1429d879 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_mv_index_1429d879
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, j JSON)
< ok 0 0
> exec
| ALTER TABLE `pa_add_mv_index_1429d879`.`t` ADD KEY `mv` ((cast(json_extract(`j`,_utf8mb4'$.ids') as unsigned array)))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_mv_index_1429d879.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `j` json DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `mv` ((cast(json_extract(`j`,_utf8mb4'$.ids') as unsigned array)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_IndexMigrationApplyCorrectness/8.0/add-prefix default#1
> exec
| DROP DATABASE IF EXISTS pa_add_prefix_dc85778f_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_prefix_dc85778f_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_prefix_dc85778f_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, s VARCHAR(100), KEY sp (s(10)))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_prefix_dc85778f_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `s` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `sp` (`s`(10))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_add_prefix_dc85778f_from
< ok 0 0
> exec
| CREATE DATABASE pa_add_prefix_dc85778f_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_prefix_dc85778f_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, s VARCHAR(100))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_prefix_dc85778f_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `s` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_add_prefix_dc85778f
< ok 0 0
> exec
| CREATE DATABASE pa_add_prefix_dc85778f DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_prefix_dc85778f
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, s VARCHAR(100))
< ok 0 0
> exec
| ALTER TABLE `pa_add_prefix_dc85778f`.`t` ADD KEY `sp` (`s`(10))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_prefix_dc85778f.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `s` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `sp` (`s`(10))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_IndexMigrationApplyCorrectness/8.0/add-primary-key default#1
> exec
| DROP DATABASE IF EXISTS pa_add_primary_key_14ec0169_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_primary_key_14ec0169_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_primary_key_14ec0169_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_primary_key_14ec0169_to.t
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
| DROP DATABASE IF EXISTS pa_add_primary_key_14ec0169_from
< ok 0 0
> exec
| CREATE DATABASE pa_add_primary_key_14ec0169_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_primary_key_14ec0169_from
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_primary_key_14ec0169_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_add_primary_key_14ec0169
< ok 0 0
> exec
| CREATE DATABASE pa_add_primary_key_14ec0169 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_primary_key_14ec0169
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, a INT)
< ok 0 0
> exec
| ALTER TABLE `pa_add_primary_key_14ec0169`.`t` ADD PRIMARY KEY (`id`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_primary_key_14ec0169.t
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

== TestOracle_IndexMigrationApplyCorrectness/8.0/add-secondary default#1
> exec
| DROP DATABASE IF EXISTS pa_add_secondary_ac0a30bf_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_secondary_ac0a30bf_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_secondary_ac0a30bf_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ka (a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_secondary_ac0a30bf_to.t
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
|   PRIMARY KEY (`id`),
|   KEY `ka` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_add_secondary_ac0a30bf_from
< ok 0 0
> exec
| CREATE DATABASE pa_add_secondary_ac0a30bf_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_secondary_ac0a30bf_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_secondary_ac0a30bf_from.t
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
| DROP DATABASE IF EXISTS pa_add_secondary_ac0a30bf
< ok 0 0
> exec
| CREATE DATABASE pa_add_secondary_ac0a30bf DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_secondary_ac0a30bf
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> exec
| ALTER TABLE `pa_add_secondary_ac0a30bf`.`t` ADD KEY `ka` (`a`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_secondary_ac0a30bf.t
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
|   PRIMARY KEY (`id`),
|   KEY `ka` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_IndexMigrationApplyCorrectness/8.0/add-unique default#1
> exec
| DROP DATABASE IF EXISTS pa_add_unique_278ca5f8_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_unique_278ca5f8_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_unique_278ca5f8_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, UNIQUE KEY ua (a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_unique_278ca5f8_to.t
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
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `ua` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_add_unique_278ca5f8_from
< ok 0 0
> exec
| CREATE DATABASE pa_add_unique_278ca5f8_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_unique_278ca5f8_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_unique_278ca5f8_from.t
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
| DROP DATABASE IF EXISTS pa_add_unique_278ca5f8
< ok 0 0
> exec
| CREATE DATABASE pa_add_unique_278ca5f8 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_unique_278ca5f8
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> exec
| ALTER TABLE `pa_add_unique_278ca5f8`.`t` ADD UNIQUE KEY `ua` (`a`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_unique_278ca5f8.t
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
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `ua` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_IndexMigrationApplyCorrectness/8.0/create-with-comment default#1
> exec
| DROP DATABASE IF EXISTS pa_create_with_comment_1f7c67e8_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_comment_1f7c67e8_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_comment_1f7c67e8_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ka (a) COMMENT 'hello')
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_comment_1f7c67e8_to.t
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
|   PRIMARY KEY (`id`),
|   KEY `ka` (`a`) COMMENT 'hello'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_with_comment_1f7c67e8
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_comment_1f7c67e8 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_comment_1f7c67e8
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_comment_1f7c67e8`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_comment_1f7c67e8`.`t` ADD KEY `ka` (`a`) COMMENT 'hello'
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_comment_1f7c67e8.t
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
|   PRIMARY KEY (`id`),
|   KEY `ka` (`a`) COMMENT 'hello'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_IndexMigrationApplyCorrectness/8.0/create-with-composite-pk default#1
> exec
| DROP DATABASE IF EXISTS pa_create_with_composite_pk_d9a446a4_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_composite_pk_d9a446a4_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_composite_pk_d9a446a4_to
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, c INT, PRIMARY KEY (a,b), KEY kc (c))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_composite_pk_d9a446a4_to.t
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
|   `c` int DEFAULT NULL,
|   PRIMARY KEY (`a`,`b`),
|   KEY `kc` (`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_with_composite_pk_d9a446a4
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_composite_pk_d9a446a4 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_composite_pk_d9a446a4
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_composite_pk_d9a446a4`.`t` (
|   `a` int NOT NULL,
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   PRIMARY KEY (`a`,`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_composite_pk_d9a446a4`.`t` ADD KEY `kc` (`c`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_composite_pk_d9a446a4.t
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
|   `c` int DEFAULT NULL,
|   PRIMARY KEY (`a`,`b`),
|   KEY `kc` (`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_IndexMigrationApplyCorrectness/8.0/create-with-descending default#1
> exec
| DROP DATABASE IF EXISTS pa_create_with_descending_f8b1c381_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_descending_f8b1c381_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_descending_f8b1c381_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY ab (a, b DESC))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_descending_f8b1c381_to.t
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
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `ab` (`a`,`b` DESC)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_with_descending_f8b1c381
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_descending_f8b1c381 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_descending_f8b1c381
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_descending_f8b1c381`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_descending_f8b1c381`.`t` ADD KEY `ab` (`a`,`b` DESC)
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_descending_f8b1c381.t
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
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `ab` (`a`,`b` DESC)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_IndexMigrationApplyCorrectness/8.0/create-with-fulltext default#1
> exec
| DROP DATABASE IF EXISTS pa_create_with_fulltext_6f01e751_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_fulltext_6f01e751_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_fulltext_6f01e751_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, body TEXT, FULLTEXT KEY ftb (body))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_fulltext_6f01e751_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `body` text,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `ftb` (`body`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_with_fulltext_6f01e751
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_fulltext_6f01e751 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_fulltext_6f01e751
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_fulltext_6f01e751`.`t` (
|   `id` int NOT NULL,
|   `body` text,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_fulltext_6f01e751`.`t` ADD FULLTEXT KEY `ftb` (`body`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_fulltext_6f01e751.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `body` text,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `ftb` (`body`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_IndexMigrationApplyCorrectness/8.0/create-with-functional default#1
> exec
| DROP DATABASE IF EXISTS pa_create_with_functional_5d1778a2_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_functional_5d1778a2_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_functional_5d1778a2_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY fa ((a + 1)))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_functional_5d1778a2_to.t
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
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `fa` (((`a` + 1)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_with_functional_5d1778a2
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_functional_5d1778a2 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_functional_5d1778a2
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_functional_5d1778a2`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_functional_5d1778a2`.`t` ADD KEY `fa` (((`a` + 1)))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_functional_5d1778a2.t
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
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `fa` (((`a` + 1)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_IndexMigrationApplyCorrectness/8.0/create-with-invisible default#1
> exec
| DROP DATABASE IF EXISTS pa_create_with_invisible_cc42a0a6_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_invisible_cc42a0a6_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_invisible_cc42a0a6_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY va (a) /*!80000 INVISIBLE */)
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_invisible_cc42a0a6_to.t
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
|   PRIMARY KEY (`id`),
|   KEY `va` (`a`) /*!80000 INVISIBLE */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_with_invisible_cc42a0a6
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_invisible_cc42a0a6 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_invisible_cc42a0a6
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_invisible_cc42a0a6`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_invisible_cc42a0a6`.`t` ADD KEY `va` (`a`) /*!80000 INVISIBLE */
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_invisible_cc42a0a6.t
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
|   PRIMARY KEY (`id`),
|   KEY `va` (`a`) /*!80000 INVISIBLE */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_IndexMigrationApplyCorrectness/8.0/create-with-kbs default#1
> exec
| DROP DATABASE IF EXISTS pa_create_with_kbs_e6e3578b_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_kbs_e6e3578b_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_kbs_e6e3578b_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ka (a) KEY_BLOCK_SIZE=4) ENGINE=InnoDB
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_kbs_e6e3578b_to.t
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
|   PRIMARY KEY (`id`),
|   KEY `ka` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_with_kbs_e6e3578b
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_kbs_e6e3578b DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_kbs_e6e3578b
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_kbs_e6e3578b`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_kbs_e6e3578b`.`t` ADD KEY `ka` (`a`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_kbs_e6e3578b.t
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
|   PRIMARY KEY (`id`),
|   KEY `ka` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_IndexMigrationApplyCorrectness/8.0/create-with-mv-index default#1
> exec
| DROP DATABASE IF EXISTS pa_create_with_mv_index_c826891d_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_mv_index_c826891d_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_mv_index_c826891d_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, j JSON, KEY mv ((CAST(j->'$.ids' AS UNSIGNED ARRAY))))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_mv_index_c826891d_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `j` json DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `mv` ((cast(json_extract(`j`,_utf8mb4'$.ids') as unsigned array)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_with_mv_index_c826891d
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_mv_index_c826891d DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_mv_index_c826891d
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_mv_index_c826891d`.`t` (
|   `id` int NOT NULL,
|   `j` json DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_mv_index_c826891d`.`t` ADD KEY `mv` ((cast(json_extract(`j`,_utf8mb4'$.ids') as unsigned array)))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_mv_index_c826891d.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `j` json DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `mv` ((cast(json_extract(`j`,_utf8mb4'$.ids') as unsigned array)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_IndexMigrationApplyCorrectness/8.0/create-with-prefix default#1
> exec
| DROP DATABASE IF EXISTS pa_create_with_prefix_7fafef33_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_prefix_7fafef33_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_prefix_7fafef33_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, s VARCHAR(100), KEY sp (s(10)))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_prefix_7fafef33_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `s` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `sp` (`s`(10))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_with_prefix_7fafef33
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_prefix_7fafef33 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_prefix_7fafef33
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_prefix_7fafef33`.`t` (
|   `id` int NOT NULL,
|   `s` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_prefix_7fafef33`.`t` ADD KEY `sp` (`s`(10))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_prefix_7fafef33.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `s` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `sp` (`s`(10))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_IndexMigrationApplyCorrectness/8.0/create-with-secondary default#1
> exec
| DROP DATABASE IF EXISTS pa_create_with_secondary_f9fda65b_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_secondary_f9fda65b_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_secondary_f9fda65b_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY ka (a), KEY kab (a,b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_secondary_f9fda65b_to.t
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
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `ka` (`a`),
|   KEY `kab` (`a`,`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_with_secondary_f9fda65b
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_secondary_f9fda65b DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_secondary_f9fda65b
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_secondary_f9fda65b`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_secondary_f9fda65b`.`t` ADD KEY `ka` (`a`)
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_secondary_f9fda65b`.`t` ADD KEY `kab` (`a`,`b`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_secondary_f9fda65b.t
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
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `ka` (`a`),
|   KEY `kab` (`a`,`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_IndexMigrationApplyCorrectness/8.0/create-with-spatial default#1
> exec
| DROP DATABASE IF EXISTS pa_create_with_spatial_33a532ed_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_spatial_33a532ed_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_spatial_33a532ed_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL, SPATIAL KEY sg (g))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_spatial_33a532ed_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL,
|   PRIMARY KEY (`id`),
|   SPATIAL KEY `sg` (`g`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_with_spatial_33a532ed
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_spatial_33a532ed DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_spatial_33a532ed
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_spatial_33a532ed`.`t` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_spatial_33a532ed`.`t` ADD SPATIAL KEY `sg` (`g`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_spatial_33a532ed.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL,
|   PRIMARY KEY (`id`),
|   SPATIAL KEY `sg` (`g`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_IndexMigrationApplyCorrectness/8.0/create-with-unique default#1
> exec
| DROP DATABASE IF EXISTS pa_create_with_unique_1a28689c_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_unique_1a28689c_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_unique_1a28689c_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, UNIQUE KEY ua (a), UNIQUE KEY uab (a,b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_unique_1a28689c_to.t
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
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `ua` (`a`),
|   UNIQUE KEY `uab` (`a`,`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_with_unique_1a28689c
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_unique_1a28689c DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_unique_1a28689c
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_unique_1a28689c`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_unique_1a28689c`.`t` ADD UNIQUE KEY `ua` (`a`)
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_unique_1a28689c`.`t` ADD UNIQUE KEY `uab` (`a`,`b`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_unique_1a28689c.t
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
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `ua` (`a`),
|   UNIQUE KEY `uab` (`a`,`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_IndexMigrationApplyCorrectness/8.0/create-with-using-btree default#1
> exec
| DROP DATABASE IF EXISTS pa_create_with_using_btree_716e96ee_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_using_btree_716e96ee_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_using_btree_716e96ee_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ua (a) USING BTREE) ENGINE=InnoDB
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_using_btree_716e96ee_to.t
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
|   PRIMARY KEY (`id`),
|   KEY `ua` (`a`) USING BTREE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_with_using_btree_716e96ee
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_using_btree_716e96ee DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_using_btree_716e96ee
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_using_btree_716e96ee`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_using_btree_716e96ee`.`t` ADD KEY `ua` (`a`) USING BTREE
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_using_btree_716e96ee.t
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
|   PRIMARY KEY (`id`),
|   KEY `ua` (`a`) USING BTREE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_IndexMigrationApplyCorrectness/8.0/drop-column-and-its-single-col-index default#1
> exec
| DROP DATABASE IF EXISTS pa_drop_column_and_its_single_col_index_70b58e8f_to
< ok 0 0
> exec
| CREATE DATABASE pa_drop_column_and_its_single_col_index_70b58e8f_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_column_and_its_single_col_index_70b58e8f_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_column_and_its_single_col_index_70b58e8f_to.t
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
| DROP DATABASE IF EXISTS pa_drop_column_and_its_single_col_index_70b58e8f_from
< ok 0 0
> exec
| CREATE DATABASE pa_drop_column_and_its_single_col_index_70b58e8f_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_column_and_its_single_col_index_70b58e8f_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ka (a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_column_and_its_single_col_index_70b58e8f_from.t
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
|   PRIMARY KEY (`id`),
|   KEY `ka` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_drop_column_and_its_single_col_index_70b58e8f
< ok 0 0
> exec
| CREATE DATABASE pa_drop_column_and_its_single_col_index_70b58e8f DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_column_and_its_single_col_index_70b58e8f
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ka (a))
< ok 0 0
> exec
| ALTER TABLE `pa_drop_column_and_its_single_col_index_70b58e8f`.`t` DROP INDEX `ka`
< ok 0 0
> exec
| ALTER TABLE `pa_drop_column_and_its_single_col_index_70b58e8f`.`t` DROP COLUMN `a`
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_column_and_its_single_col_index_70b58e8f.t
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

== TestOracle_IndexMigrationApplyCorrectness/8.0/drop-columns-and-their-composite-index default#1
> exec
| DROP DATABASE IF EXISTS pa_drop_columns_and_their_composite_index_46f973ec_to
< ok 0 0
> exec
| CREATE DATABASE pa_drop_columns_and_their_composite_index_46f973ec_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_columns_and_their_composite_index_46f973ec_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_columns_and_their_composite_index_46f973ec_to.t
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
| DROP DATABASE IF EXISTS pa_drop_columns_and_their_composite_index_46f973ec_from
< ok 0 0
> exec
| CREATE DATABASE pa_drop_columns_and_their_composite_index_46f973ec_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_columns_and_their_composite_index_46f973ec_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY kab (a,b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_columns_and_their_composite_index_46f973ec_from.t
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
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `kab` (`a`,`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_drop_columns_and_their_composite_index_46f973ec
< ok 0 0
> exec
| CREATE DATABASE pa_drop_columns_and_their_composite_index_46f973ec DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_columns_and_their_composite_index_46f973ec
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY kab (a,b))
< ok 0 0
> exec
| ALTER TABLE `pa_drop_columns_and_their_composite_index_46f973ec`.`t` DROP INDEX `kab`
< ok 0 0
> exec
| ALTER TABLE `pa_drop_columns_and_their_composite_index_46f973ec`.`t` DROP COLUMN `a`
< ok 0 0
> exec
| ALTER TABLE `pa_drop_columns_and_their_composite_index_46f973ec`.`t` DROP COLUMN `b`
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_columns_and_their_composite_index_46f973ec.t
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

== TestOracle_IndexMigrationApplyCorrectness/8.0/drop-mv-index default#1
> exec
| DROP DATABASE IF EXISTS pa_drop_mv_index_a6df9201_to
< ok 0 0
> exec
| CREATE DATABASE pa_drop_mv_index_a6df9201_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_mv_index_a6df9201_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, j JSON)
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_mv_index_a6df9201_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `j` json DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_drop_mv_index_a6df9201_from
< ok 0 0
> exec
| CREATE DATABASE pa_drop_mv_index_a6df9201_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_mv_index_a6df9201_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, j JSON, KEY mv ((CAST(j->'$.ids' AS UNSIGNED ARRAY))))
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_mv_index_a6df9201_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `j` json DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `mv` ((cast(json_extract(`j`,_utf8mb4'$.ids') as unsigned array)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_drop_mv_index_a6df9201
< ok 0 0
> exec
| CREATE DATABASE pa_drop_mv_index_a6df9201 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_mv_index_a6df9201
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, j JSON, KEY mv ((CAST(j->'$.ids' AS UNSIGNED ARRAY))))
< ok 0 0
> exec
| ALTER TABLE `pa_drop_mv_index_a6df9201`.`t` DROP INDEX `mv`
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_mv_index_a6df9201.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `j` json DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_IndexMigrationApplyCorrectness/8.0/drop-one-of-many default#1
> exec
| DROP DATABASE IF EXISTS pa_drop_one_of_many_e1d1300f_to
< ok 0 0
> exec
| CREATE DATABASE pa_drop_one_of_many_e1d1300f_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_one_of_many_e1d1300f_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY kb (b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_one_of_many_e1d1300f_to.t
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
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `kb` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_drop_one_of_many_e1d1300f_from
< ok 0 0
> exec
| CREATE DATABASE pa_drop_one_of_many_e1d1300f_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_one_of_many_e1d1300f_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY ka (a), KEY kb (b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_one_of_many_e1d1300f_from.t
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
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `ka` (`a`),
|   KEY `kb` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_drop_one_of_many_e1d1300f
< ok 0 0
> exec
| CREATE DATABASE pa_drop_one_of_many_e1d1300f DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_one_of_many_e1d1300f
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY ka (a), KEY kb (b))
< ok 0 0
> exec
| ALTER TABLE `pa_drop_one_of_many_e1d1300f`.`t` DROP INDEX `ka`
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_one_of_many_e1d1300f.t
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
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `kb` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_IndexMigrationApplyCorrectness/8.0/drop-primary-key default#1
> exec
| DROP DATABASE IF EXISTS pa_drop_primary_key_59ff7af1_to
< ok 0 0
> exec
| CREATE DATABASE pa_drop_primary_key_59ff7af1_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_primary_key_59ff7af1_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_primary_key_59ff7af1_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_drop_primary_key_59ff7af1_from
< ok 0 0
> exec
| CREATE DATABASE pa_drop_primary_key_59ff7af1_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_primary_key_59ff7af1_from
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_primary_key_59ff7af1_from.t
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
| DROP DATABASE IF EXISTS pa_drop_primary_key_59ff7af1
< ok 0 0
> exec
| CREATE DATABASE pa_drop_primary_key_59ff7af1 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_primary_key_59ff7af1
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT)
< ok 0 0
> exec
| ALTER TABLE `pa_drop_primary_key_59ff7af1`.`t` DROP PRIMARY KEY
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_primary_key_59ff7af1.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_IndexMigrationApplyCorrectness/8.0/drop-secondary default#1
> exec
| DROP DATABASE IF EXISTS pa_drop_secondary_534bb337_to
< ok 0 0
> exec
| CREATE DATABASE pa_drop_secondary_534bb337_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_secondary_534bb337_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_secondary_534bb337_to.t
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
| DROP DATABASE IF EXISTS pa_drop_secondary_534bb337_from
< ok 0 0
> exec
| CREATE DATABASE pa_drop_secondary_534bb337_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_secondary_534bb337_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ka (a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_secondary_534bb337_from.t
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
|   PRIMARY KEY (`id`),
|   KEY `ka` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_drop_secondary_534bb337
< ok 0 0
> exec
| CREATE DATABASE pa_drop_secondary_534bb337 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_secondary_534bb337
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ka (a))
< ok 0 0
> exec
| ALTER TABLE `pa_drop_secondary_534bb337`.`t` DROP INDEX `ka`
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_secondary_534bb337.t
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

== TestOracle_IndexMigrationApplyCorrectness/8.0/drop-unique default#1
> exec
| DROP DATABASE IF EXISTS pa_drop_unique_6e0d0370_to
< ok 0 0
> exec
| CREATE DATABASE pa_drop_unique_6e0d0370_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_unique_6e0d0370_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_unique_6e0d0370_to.t
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
| DROP DATABASE IF EXISTS pa_drop_unique_6e0d0370_from
< ok 0 0
> exec
| CREATE DATABASE pa_drop_unique_6e0d0370_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_unique_6e0d0370_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, UNIQUE KEY ua (a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_unique_6e0d0370_from.t
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
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `ua` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_drop_unique_6e0d0370
< ok 0 0
> exec
| CREATE DATABASE pa_drop_unique_6e0d0370 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_unique_6e0d0370
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, UNIQUE KEY ua (a))
< ok 0 0
> exec
| ALTER TABLE `pa_drop_unique_6e0d0370`.`t` DROP INDEX `ua`
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_unique_6e0d0370.t
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

== TestOracle_IndexMigrationApplyCorrectness/8.0/modify-columns default#1
> exec
| DROP DATABASE IF EXISTS pa_modify_columns_9aa1af1_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_columns_9aa1af1_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_columns_9aa1af1_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY k (a,b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_columns_9aa1af1_to.t
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
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`a`,`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_columns_9aa1af1_from
< ok 0 0
> exec
| CREATE DATABASE pa_modify_columns_9aa1af1_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_columns_9aa1af1_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY k (a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_columns_9aa1af1_from.t
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
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_columns_9aa1af1
< ok 0 0
> exec
| CREATE DATABASE pa_modify_columns_9aa1af1 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_columns_9aa1af1
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY k (a))
< ok 0 0
> exec
| ALTER TABLE `pa_modify_columns_9aa1af1`.`t` DROP INDEX `k`
< ok 0 0
> exec
| ALTER TABLE `pa_modify_columns_9aa1af1`.`t` ADD KEY `k` (`a`,`b`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_columns_9aa1af1.t
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
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`a`,`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_IndexMigrationApplyCorrectness/8.0/modify-comment default#1
> exec
| DROP DATABASE IF EXISTS pa_modify_comment_5b40ff77_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_comment_5b40ff77_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_comment_5b40ff77_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a) COMMENT 'new')
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_comment_5b40ff77_to.t
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
|   PRIMARY KEY (`id`),
|   KEY `k` (`a`) COMMENT 'new'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_comment_5b40ff77_from
< ok 0 0
> exec
| CREATE DATABASE pa_modify_comment_5b40ff77_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_comment_5b40ff77_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a) COMMENT 'old')
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_comment_5b40ff77_from.t
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
|   PRIMARY KEY (`id`),
|   KEY `k` (`a`) COMMENT 'old'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_comment_5b40ff77
< ok 0 0
> exec
| CREATE DATABASE pa_modify_comment_5b40ff77 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_comment_5b40ff77
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a) COMMENT 'old')
< ok 0 0
> exec
| ALTER TABLE `pa_modify_comment_5b40ff77`.`t` DROP INDEX `k`
< ok 0 0
> exec
| ALTER TABLE `pa_modify_comment_5b40ff77`.`t` ADD KEY `k` (`a`) COMMENT 'new'
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_comment_5b40ff77.t
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
|   PRIMARY KEY (`id`),
|   KEY `k` (`a`) COMMENT 'new'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_IndexMigrationApplyCorrectness/8.0/modify-direction default#1
> exec
| DROP DATABASE IF EXISTS pa_modify_direction_ef73fc2f_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_direction_ef73fc2f_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_direction_ef73fc2f_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a DESC))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_direction_ef73fc2f_to.t
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
|   PRIMARY KEY (`id`),
|   KEY `k` (`a` DESC)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_direction_ef73fc2f_from
< ok 0 0
> exec
| CREATE DATABASE pa_modify_direction_ef73fc2f_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_direction_ef73fc2f_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_direction_ef73fc2f_from.t
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
|   PRIMARY KEY (`id`),
|   KEY `k` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_direction_ef73fc2f
< ok 0 0
> exec
| CREATE DATABASE pa_modify_direction_ef73fc2f DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_direction_ef73fc2f
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a))
< ok 0 0
> exec
| ALTER TABLE `pa_modify_direction_ef73fc2f`.`t` DROP INDEX `k`
< ok 0 0
> exec
| ALTER TABLE `pa_modify_direction_ef73fc2f`.`t` ADD KEY `k` (`a` DESC)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_direction_ef73fc2f.t
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
|   PRIMARY KEY (`id`),
|   KEY `k` (`a` DESC)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_IndexMigrationApplyCorrectness/8.0/modify-index-shrink-with-column-drop default#1
> exec
| DROP DATABASE IF EXISTS pa_modify_index_shrink_with_column_drop_b9513dee_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_index_shrink_with_column_drop_b9513dee_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_index_shrink_with_column_drop_b9513dee_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, b INT, KEY k (b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_index_shrink_with_column_drop_b9513dee_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_index_shrink_with_column_drop_b9513dee_from
< ok 0 0
> exec
| CREATE DATABASE pa_modify_index_shrink_with_column_drop_b9513dee_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_index_shrink_with_column_drop_b9513dee_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY k (a,b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_index_shrink_with_column_drop_b9513dee_from.t
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
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`a`,`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_index_shrink_with_column_drop_b9513dee
< ok 0 0
> exec
| CREATE DATABASE pa_modify_index_shrink_with_column_drop_b9513dee DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_index_shrink_with_column_drop_b9513dee
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY k (a,b))
< ok 0 0
> exec
| ALTER TABLE `pa_modify_index_shrink_with_column_drop_b9513dee`.`t` DROP INDEX `k`
< ok 0 0
> exec
| ALTER TABLE `pa_modify_index_shrink_with_column_drop_b9513dee`.`t` DROP COLUMN `a`
< ok 0 0
> exec
| ALTER TABLE `pa_modify_index_shrink_with_column_drop_b9513dee`.`t` ADD KEY `k` (`b`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_index_shrink_with_column_drop_b9513dee.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_IndexMigrationApplyCorrectness/8.0/modify-kind-to-unique default#1
> exec
| DROP DATABASE IF EXISTS pa_modify_kind_to_unique_832aa88a_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_kind_to_unique_832aa88a_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_kind_to_unique_832aa88a_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, UNIQUE KEY k (a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_kind_to_unique_832aa88a_to.t
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
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `k` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_kind_to_unique_832aa88a_from
< ok 0 0
> exec
| CREATE DATABASE pa_modify_kind_to_unique_832aa88a_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_kind_to_unique_832aa88a_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_kind_to_unique_832aa88a_from.t
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
|   PRIMARY KEY (`id`),
|   KEY `k` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_kind_to_unique_832aa88a
< ok 0 0
> exec
| CREATE DATABASE pa_modify_kind_to_unique_832aa88a DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_kind_to_unique_832aa88a
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a))
< ok 0 0
> exec
| ALTER TABLE `pa_modify_kind_to_unique_832aa88a`.`t` DROP INDEX `k`
< ok 0 0
> exec
| ALTER TABLE `pa_modify_kind_to_unique_832aa88a`.`t` ADD UNIQUE KEY `k` (`a`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_kind_to_unique_832aa88a.t
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
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `k` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_IndexMigrationApplyCorrectness/8.0/modify-pk-autoincrement default#1
> exec
| DROP DATABASE IF EXISTS pa_modify_pk_autoincrement_b84ccdaa_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_pk_autoincrement_b84ccdaa_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_pk_autoincrement_b84ccdaa_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, PRIMARY KEY (id,a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_pk_autoincrement_b84ccdaa_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `a` int NOT NULL,
|   PRIMARY KEY (`id`,`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_pk_autoincrement_b84ccdaa_from
< ok 0 0
> exec
| CREATE DATABASE pa_modify_pk_autoincrement_b84ccdaa_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_pk_autoincrement_b84ccdaa_from
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT PRIMARY KEY, a INT NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_pk_autoincrement_b84ccdaa_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `a` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_pk_autoincrement_b84ccdaa
< ok 0 0
> exec
| CREATE DATABASE pa_modify_pk_autoincrement_b84ccdaa DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_pk_autoincrement_b84ccdaa
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT PRIMARY KEY, a INT NOT NULL)
< ok 0 0
> exec
| ALTER TABLE `pa_modify_pk_autoincrement_b84ccdaa`.`t` DROP PRIMARY KEY, ADD PRIMARY KEY (`id`,`a`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_pk_autoincrement_b84ccdaa.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `a` int NOT NULL,
|   PRIMARY KEY (`id`,`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_IndexMigrationApplyCorrectness/8.0/modify-pk-columns default#1
> exec
| DROP DATABASE IF EXISTS pa_modify_pk_columns_9150244d_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_pk_columns_9150244d_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_pk_columns_9150244d_to
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, PRIMARY KEY (a,b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_pk_columns_9150244d_to.t
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
> exec
| DROP DATABASE IF EXISTS pa_modify_pk_columns_9150244d_from
< ok 0 0
> exec
| CREATE DATABASE pa_modify_pk_columns_9150244d_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_pk_columns_9150244d_from
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, PRIMARY KEY (a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_pk_columns_9150244d_from.t
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
|   PRIMARY KEY (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_pk_columns_9150244d
< ok 0 0
> exec
| CREATE DATABASE pa_modify_pk_columns_9150244d DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_pk_columns_9150244d
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, PRIMARY KEY (a))
< ok 0 0
> exec
| ALTER TABLE `pa_modify_pk_columns_9150244d`.`t` DROP PRIMARY KEY
< ok 0 0
> exec
| ALTER TABLE `pa_modify_pk_columns_9150244d`.`t` ADD PRIMARY KEY (`a`,`b`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_pk_columns_9150244d.t
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

== TestOracle_IndexMigrationApplyCorrectness/8.0/modify-pk-drop-old-member default#1
> exec
| DROP DATABASE IF EXISTS pa_modify_pk_drop_old_member_63f6c0f0_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_pk_drop_old_member_63f6c0f0_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_pk_drop_old_member_63f6c0f0_to
< ok 0 0
> exec
| CREATE TABLE t (b INT NOT NULL, c INT, PRIMARY KEY (b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_pk_drop_old_member_63f6c0f0_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_pk_drop_old_member_63f6c0f0_from
< ok 0 0
> exec
| CREATE DATABASE pa_modify_pk_drop_old_member_63f6c0f0_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_pk_drop_old_member_63f6c0f0_from
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, PRIMARY KEY (a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_pk_drop_old_member_63f6c0f0_from.t
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
|   PRIMARY KEY (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_pk_drop_old_member_63f6c0f0
< ok 0 0
> exec
| CREATE DATABASE pa_modify_pk_drop_old_member_63f6c0f0 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_pk_drop_old_member_63f6c0f0
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, PRIMARY KEY (a))
< ok 0 0
> exec
| ALTER TABLE `pa_modify_pk_drop_old_member_63f6c0f0`.`t` DROP PRIMARY KEY
< ok 0 0
> exec
| ALTER TABLE `pa_modify_pk_drop_old_member_63f6c0f0`.`t` DROP COLUMN `a`
< ok 0 0
> exec
| ALTER TABLE `pa_modify_pk_drop_old_member_63f6c0f0`.`t` ADD COLUMN `c` int DEFAULT NULL
< ok 0 0
> exec
| ALTER TABLE `pa_modify_pk_drop_old_member_63f6c0f0`.`t` ADD PRIMARY KEY (`b`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_pk_drop_old_member_63f6c0f0.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_IndexMigrationApplyCorrectness/8.0/modify-pk-shrink-demote-nullable default#1
> exec
| DROP DATABASE IF EXISTS pa_modify_pk_shrink_demote_nullable_28ffb648_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_pk_shrink_demote_nullable_28ffb648_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_pk_shrink_demote_nullable_28ffb648_to
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NULL, PRIMARY KEY (a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_pk_shrink_demote_nullable_28ffb648_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int NOT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_pk_shrink_demote_nullable_28ffb648_from
< ok 0 0
> exec
| CREATE DATABASE pa_modify_pk_shrink_demote_nullable_28ffb648_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_pk_shrink_demote_nullable_28ffb648_from
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, PRIMARY KEY (a,b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_pk_shrink_demote_nullable_28ffb648_from.t
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
> exec
| DROP DATABASE IF EXISTS pa_modify_pk_shrink_demote_nullable_28ffb648
< ok 0 0
> exec
| CREATE DATABASE pa_modify_pk_shrink_demote_nullable_28ffb648 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_pk_shrink_demote_nullable_28ffb648
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, PRIMARY KEY (a,b))
< ok 0 0
> exec
| ALTER TABLE `pa_modify_pk_shrink_demote_nullable_28ffb648`.`t` DROP PRIMARY KEY
< ok 0 0
> exec
| ALTER TABLE `pa_modify_pk_shrink_demote_nullable_28ffb648`.`t` MODIFY COLUMN `b` int DEFAULT NULL
< ok 0 0
> exec
| ALTER TABLE `pa_modify_pk_shrink_demote_nullable_28ffb648`.`t` ADD PRIMARY KEY (`a`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_pk_shrink_demote_nullable_28ffb648.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int NOT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_IndexMigrationApplyCorrectness/8.0/modify-pk-swap-autoincrement default#1
> exec
| DROP DATABASE IF EXISTS pa_modify_pk_swap_autoincrement_60311f00_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_pk_swap_autoincrement_60311f00_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_pk_swap_autoincrement_60311f00_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, b INT NOT NULL, PRIMARY KEY (id,b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_pk_swap_autoincrement_60311f00_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `a` int NOT NULL,
|   `b` int NOT NULL,
|   PRIMARY KEY (`id`,`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_pk_swap_autoincrement_60311f00_from
< ok 0 0
> exec
| CREATE DATABASE pa_modify_pk_swap_autoincrement_60311f00_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_pk_swap_autoincrement_60311f00_from
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, b INT NOT NULL, PRIMARY KEY (id,a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_pk_swap_autoincrement_60311f00_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `a` int NOT NULL,
|   `b` int NOT NULL,
|   PRIMARY KEY (`id`,`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_pk_swap_autoincrement_60311f00
< ok 0 0
> exec
| CREATE DATABASE pa_modify_pk_swap_autoincrement_60311f00 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_pk_swap_autoincrement_60311f00
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, b INT NOT NULL, PRIMARY KEY (id,a))
< ok 0 0
> exec
| ALTER TABLE `pa_modify_pk_swap_autoincrement_60311f00`.`t` DROP PRIMARY KEY, ADD PRIMARY KEY (`id`,`b`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_pk_swap_autoincrement_60311f00.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `a` int NOT NULL,
|   `b` int NOT NULL,
|   PRIMARY KEY (`id`,`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_IndexMigrationApplyCorrectness/8.0/modify-prefix-length default#1
> exec
| DROP DATABASE IF EXISTS pa_modify_prefix_length_6325d501_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_prefix_length_6325d501_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_prefix_length_6325d501_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, s VARCHAR(100), KEY k (s(20)))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_prefix_length_6325d501_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `s` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`s`(20))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_prefix_length_6325d501_from
< ok 0 0
> exec
| CREATE DATABASE pa_modify_prefix_length_6325d501_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_prefix_length_6325d501_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, s VARCHAR(100), KEY k (s(10)))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_prefix_length_6325d501_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `s` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`s`(10))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_prefix_length_6325d501
< ok 0 0
> exec
| CREATE DATABASE pa_modify_prefix_length_6325d501 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_prefix_length_6325d501
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, s VARCHAR(100), KEY k (s(10)))
< ok 0 0
> exec
| ALTER TABLE `pa_modify_prefix_length_6325d501`.`t` DROP INDEX `k`
< ok 0 0
> exec
| ALTER TABLE `pa_modify_prefix_length_6325d501`.`t` ADD KEY `k` (`s`(20))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_prefix_length_6325d501.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `s` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`s`(20))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_IndexMigrationApplyCorrectness/8.0/modify-visibility default#1
> exec
| DROP DATABASE IF EXISTS pa_modify_visibility_40157212_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_visibility_40157212_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_visibility_40157212_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a) /*!80000 INVISIBLE */)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_visibility_40157212_to.t
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
|   PRIMARY KEY (`id`),
|   KEY `k` (`a`) /*!80000 INVISIBLE */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_visibility_40157212_from
< ok 0 0
> exec
| CREATE DATABASE pa_modify_visibility_40157212_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_visibility_40157212_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_visibility_40157212_from.t
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
|   PRIMARY KEY (`id`),
|   KEY `k` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_visibility_40157212
< ok 0 0
> exec
| CREATE DATABASE pa_modify_visibility_40157212 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_visibility_40157212
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a))
< ok 0 0
> exec
| ALTER TABLE `pa_modify_visibility_40157212`.`t` DROP INDEX `k`
< ok 0 0
> exec
| ALTER TABLE `pa_modify_visibility_40157212`.`t` ADD KEY `k` (`a`) /*!80000 INVISIBLE */
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_visibility_40157212.t
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
|   PRIMARY KEY (`id`),
|   KEY `k` (`a`) /*!80000 INVISIBLE */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
