-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracle_CheckMigrationApplyCorrectness default#1

== TestOracle_CheckMigrationApplyCorrectness/8.0/add-column-and-check-on-it default#1
> exec
| DROP DATABASE IF EXISTS pa_add_column_and_check_on_it_f79d83d4_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_column_and_check_on_it_f79d83d4_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_column_and_check_on_it_f79d83d4_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_column_and_check_on_it_f79d83d4_to.t
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
|   CONSTRAINT `ck` CHECK ((`a` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_add_column_and_check_on_it_f79d83d4_from
< ok 0 0
> exec
| CREATE DATABASE pa_add_column_and_check_on_it_f79d83d4_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_column_and_check_on_it_f79d83d4_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_column_and_check_on_it_f79d83d4_from.t
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
| DROP DATABASE IF EXISTS pa_add_column_and_check_on_it_f79d83d4
< ok 0 0
> exec
| CREATE DATABASE pa_add_column_and_check_on_it_f79d83d4 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_column_and_check_on_it_f79d83d4
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE `pa_add_column_and_check_on_it_f79d83d4`.`t` ADD COLUMN `a` int DEFAULT NULL
< ok 0 0
> exec
| ALTER TABLE `pa_add_column_and_check_on_it_f79d83d4`.`t` ADD CONSTRAINT `ck` CHECK ((`a` > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_column_and_check_on_it_f79d83d4.t
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
|   CONSTRAINT `ck` CHECK ((`a` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_CheckMigrationApplyCorrectness/8.0/add-compound default#1
> exec
| DROP DATABASE IF EXISTS pa_add_compound_cee99ee_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_compound_cee99ee_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_compound_cee99ee_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CHECK (a >= 0 AND b <= 100))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_compound_cee99ee_to.t
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
|   CONSTRAINT `t_chk_1` CHECK (((`a` >= 0) and (`b` <= 100)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_add_compound_cee99ee_from
< ok 0 0
> exec
| CREATE DATABASE pa_add_compound_cee99ee_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_compound_cee99ee_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_compound_cee99ee_from.t
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
| DROP DATABASE IF EXISTS pa_add_compound_cee99ee
< ok 0 0
> exec
| CREATE DATABASE pa_add_compound_cee99ee DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_compound_cee99ee
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT)
< ok 0 0
> exec
| ALTER TABLE `pa_add_compound_cee99ee`.`t` ADD CONSTRAINT `t_chk_1` CHECK (((`a` >= 0) and (`b` <= 100)))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_compound_cee99ee.t
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
|   CONSTRAINT `t_chk_1` CHECK (((`a` >= 0) and (`b` <= 100)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_CheckMigrationApplyCorrectness/8.0/add-named-check default#1
> exec
| DROP DATABASE IF EXISTS pa_add_named_check_742a6701_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_named_check_742a6701_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_named_check_742a6701_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_named_check_742a6701_to.t
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
|   CONSTRAINT `ck` CHECK ((`a` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_add_named_check_742a6701_from
< ok 0 0
> exec
| CREATE DATABASE pa_add_named_check_742a6701_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_named_check_742a6701_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_named_check_742a6701_from.t
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
| DROP DATABASE IF EXISTS pa_add_named_check_742a6701
< ok 0 0
> exec
| CREATE DATABASE pa_add_named_check_742a6701 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_named_check_742a6701
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> exec
| ALTER TABLE `pa_add_named_check_742a6701`.`t` ADD CONSTRAINT `ck` CHECK ((`a` > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_named_check_742a6701.t
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
|   CONSTRAINT `ck` CHECK ((`a` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_CheckMigrationApplyCorrectness/8.0/add-not-enforced default#1
> exec
| DROP DATABASE IF EXISTS pa_add_not_enforced_ff99c12d_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_not_enforced_ff99c12d_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_not_enforced_ff99c12d_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT ne CHECK (a < b) NOT ENFORCED)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_not_enforced_ff99c12d_to.t
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
|   CONSTRAINT `ne` CHECK ((`a` < `b`)) /*!80016 NOT ENFORCED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_add_not_enforced_ff99c12d_from
< ok 0 0
> exec
| CREATE DATABASE pa_add_not_enforced_ff99c12d_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_not_enforced_ff99c12d_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_not_enforced_ff99c12d_from.t
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
| DROP DATABASE IF EXISTS pa_add_not_enforced_ff99c12d
< ok 0 0
> exec
| CREATE DATABASE pa_add_not_enforced_ff99c12d DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_not_enforced_ff99c12d
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT)
< ok 0 0
> exec
| ALTER TABLE `pa_add_not_enforced_ff99c12d`.`t` ADD CONSTRAINT `ne` CHECK ((`a` < `b`)) /*!80016 NOT ENFORCED */
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_not_enforced_ff99c12d.t
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
|   CONSTRAINT `ne` CHECK ((`a` < `b`)) /*!80016 NOT ENFORCED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_CheckMigrationApplyCorrectness/8.0/add-second-check default#1
> exec
| DROP DATABASE IF EXISTS pa_add_second_check_ee3207e_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_second_check_ee3207e_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_second_check_ee3207e_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT c1 CHECK (a > 0), CONSTRAINT c2 CHECK (b < 100))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_second_check_ee3207e_to.t
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
|   CONSTRAINT `c1` CHECK ((`a` > 0)),
|   CONSTRAINT `c2` CHECK ((`b` < 100))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_add_second_check_ee3207e_from
< ok 0 0
> exec
| CREATE DATABASE pa_add_second_check_ee3207e_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_second_check_ee3207e_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT c1 CHECK (a > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_second_check_ee3207e_from.t
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
|   CONSTRAINT `c1` CHECK ((`a` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_add_second_check_ee3207e
< ok 0 0
> exec
| CREATE DATABASE pa_add_second_check_ee3207e DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_second_check_ee3207e
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT c1 CHECK (a > 0))
< ok 0 0
> exec
| ALTER TABLE `pa_add_second_check_ee3207e`.`t` ADD CONSTRAINT `c2` CHECK ((`b` < 100))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_second_check_ee3207e.t
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
|   CONSTRAINT `c1` CHECK ((`a` > 0)),
|   CONSTRAINT `c2` CHECK ((`b` < 100))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_CheckMigrationApplyCorrectness/8.0/add-unnamed-check default#1
> exec
| DROP DATABASE IF EXISTS pa_add_unnamed_check_b26a9f80_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_unnamed_check_b26a9f80_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_unnamed_check_b26a9f80_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CHECK (a > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_unnamed_check_b26a9f80_to.t
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
|   CONSTRAINT `t_chk_1` CHECK ((`a` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_add_unnamed_check_b26a9f80_from
< ok 0 0
> exec
| CREATE DATABASE pa_add_unnamed_check_b26a9f80_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_unnamed_check_b26a9f80_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_unnamed_check_b26a9f80_from.t
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
| DROP DATABASE IF EXISTS pa_add_unnamed_check_b26a9f80
< ok 0 0
> exec
| CREATE DATABASE pa_add_unnamed_check_b26a9f80 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_unnamed_check_b26a9f80
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> exec
| ALTER TABLE `pa_add_unnamed_check_b26a9f80`.`t` ADD CONSTRAINT `t_chk_1` CHECK ((`a` > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_unnamed_check_b26a9f80.t
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
|   CONSTRAINT `t_chk_1` CHECK ((`a` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_CheckMigrationApplyCorrectness/8.0/create-with-column-level-check default#1
> exec
| DROP DATABASE IF EXISTS pa_create_with_column_level_check_559cbacf_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_column_level_check_559cbacf_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_column_level_check_559cbacf_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT CHECK (a > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_column_level_check_559cbacf_to.t
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
|   CONSTRAINT `t_chk_1` CHECK ((`a` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_with_column_level_check_559cbacf
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_column_level_check_559cbacf DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_column_level_check_559cbacf
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_column_level_check_559cbacf`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_column_level_check_559cbacf`.`t` ADD CONSTRAINT `t_chk_1` CHECK ((`a` > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_column_level_check_559cbacf.t
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
|   CONSTRAINT `t_chk_1` CHECK ((`a` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_CheckMigrationApplyCorrectness/8.0/create-with-compound default#1
> exec
| DROP DATABASE IF EXISTS pa_create_with_compound_acc19902_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_compound_acc19902_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_compound_acc19902_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CHECK (a >= 0 AND b <= 100))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_compound_acc19902_to.t
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
|   CONSTRAINT `t_chk_1` CHECK (((`a` >= 0) and (`b` <= 100)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_with_compound_acc19902
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_compound_acc19902 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_compound_acc19902
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_compound_acc19902`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_compound_acc19902`.`t` ADD CONSTRAINT `t_chk_1` CHECK (((`a` >= 0) and (`b` <= 100)))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_compound_acc19902.t
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
|   CONSTRAINT `t_chk_1` CHECK (((`a` >= 0) and (`b` <= 100)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_CheckMigrationApplyCorrectness/8.0/create-with-function default#1
> exec
| DROP DATABASE IF EXISTS pa_create_with_function_298eb693_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_function_298eb693_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_function_298eb693_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, s VARCHAR(20), CHECK (CHAR_LENGTH(s) > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_function_298eb693_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `s` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `t_chk_1` CHECK ((char_length(`s`) > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_with_function_298eb693
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_function_298eb693 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_function_298eb693
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_function_298eb693`.`t` (
|   `id` int NOT NULL,
|   `s` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_function_298eb693`.`t` ADD CONSTRAINT `t_chk_1` CHECK ((char_length(`s`) > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_function_298eb693.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `s` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `t_chk_1` CHECK ((char_length(`s`) > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_CheckMigrationApplyCorrectness/8.0/create-with-multiple-checks default#1
> exec
| DROP DATABASE IF EXISTS pa_create_with_multiple_checks_896211f_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_multiple_checks_896211f_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_multiple_checks_896211f_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CHECK (a > 0), CONSTRAINT cb CHECK (b < 100), CHECK (a < b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_multiple_checks_896211f_to.t
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
|   CONSTRAINT `cb` CHECK ((`b` < 100)),
|   CONSTRAINT `t_chk_1` CHECK ((`a` > 0)),
|   CONSTRAINT `t_chk_2` CHECK ((`a` < `b`))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_with_multiple_checks_896211f
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_multiple_checks_896211f DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_multiple_checks_896211f
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_multiple_checks_896211f`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_multiple_checks_896211f`.`t` ADD CONSTRAINT `cb` CHECK ((`b` < 100))
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_multiple_checks_896211f`.`t` ADD CONSTRAINT `t_chk_1` CHECK ((`a` > 0))
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_multiple_checks_896211f`.`t` ADD CONSTRAINT `t_chk_2` CHECK ((`a` < `b`))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_multiple_checks_896211f.t
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
|   CONSTRAINT `cb` CHECK ((`b` < 100)),
|   CONSTRAINT `t_chk_1` CHECK ((`a` > 0)),
|   CONSTRAINT `t_chk_2` CHECK ((`a` < `b`))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_CheckMigrationApplyCorrectness/8.0/create-with-named-check default#1
> exec
| DROP DATABASE IF EXISTS pa_create_with_named_check_65998e6d_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_named_check_65998e6d_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_named_check_65998e6d_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_named_check_65998e6d_to.t
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
|   CONSTRAINT `ck` CHECK ((`a` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_with_named_check_65998e6d
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_named_check_65998e6d DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_named_check_65998e6d
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_named_check_65998e6d`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_named_check_65998e6d`.`t` ADD CONSTRAINT `ck` CHECK ((`a` > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_named_check_65998e6d.t
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
|   CONSTRAINT `ck` CHECK ((`a` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_CheckMigrationApplyCorrectness/8.0/create-with-not-enforced default#1
> exec
| DROP DATABASE IF EXISTS pa_create_with_not_enforced_cb97b301_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_not_enforced_cb97b301_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_not_enforced_cb97b301_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT ne CHECK (a < b) NOT ENFORCED)
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_not_enforced_cb97b301_to.t
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
|   CONSTRAINT `ne` CHECK ((`a` < `b`)) /*!80016 NOT ENFORCED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_with_not_enforced_cb97b301
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_not_enforced_cb97b301 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_not_enforced_cb97b301
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_not_enforced_cb97b301`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_not_enforced_cb97b301`.`t` ADD CONSTRAINT `ne` CHECK ((`a` < `b`)) /*!80016 NOT ENFORCED */
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_not_enforced_cb97b301.t
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
|   CONSTRAINT `ne` CHECK ((`a` < `b`)) /*!80016 NOT ENFORCED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_CheckMigrationApplyCorrectness/8.0/create-with-string-literal default#1
> exec
| DROP DATABASE IF EXISTS pa_create_with_string_literal_dd023428_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_string_literal_dd023428_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_string_literal_dd023428_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, s VARCHAR(20), CHECK (s = 'hello')) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_string_literal_dd023428_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `s` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `t_chk_1` CHECK ((`s` = _utf8mb4'hello'))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_with_string_literal_dd023428
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_string_literal_dd023428 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_string_literal_dd023428
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_string_literal_dd023428`.`t` (
|   `id` int NOT NULL,
|   `s` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_string_literal_dd023428`.`t` ADD CONSTRAINT `t_chk_1` CHECK ((`s` = _utf8mb4'hello'))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_string_literal_dd023428.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `s` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `t_chk_1` CHECK ((`s` = _utf8mb4'hello'))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_CheckMigrationApplyCorrectness/8.0/create-with-unnamed-check default#1
> exec
| DROP DATABASE IF EXISTS pa_create_with_unnamed_check_6759801c_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_unnamed_check_6759801c_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_unnamed_check_6759801c_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CHECK (a > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_unnamed_check_6759801c_to.t
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
|   CONSTRAINT `t_chk_1` CHECK ((`a` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_with_unnamed_check_6759801c
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_unnamed_check_6759801c DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_unnamed_check_6759801c
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_unnamed_check_6759801c`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_unnamed_check_6759801c`.`t` ADD CONSTRAINT `t_chk_1` CHECK ((`a` > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_unnamed_check_6759801c.t
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
|   CONSTRAINT `t_chk_1` CHECK ((`a` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_CheckMigrationApplyCorrectness/8.0/drop-check-and-its-column default#1
> exec
| DROP DATABASE IF EXISTS pa_drop_check_and_its_column_b054c31f_to
< ok 0 0
> exec
| CREATE DATABASE pa_drop_check_and_its_column_b054c31f_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_check_and_its_column_b054c31f_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_check_and_its_column_b054c31f_to.t
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
| DROP DATABASE IF EXISTS pa_drop_check_and_its_column_b054c31f_from
< ok 0 0
> exec
| CREATE DATABASE pa_drop_check_and_its_column_b054c31f_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_check_and_its_column_b054c31f_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_check_and_its_column_b054c31f_from.t
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
|   CONSTRAINT `ck` CHECK ((`a` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_drop_check_and_its_column_b054c31f
< ok 0 0
> exec
| CREATE DATABASE pa_drop_check_and_its_column_b054c31f DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_check_and_its_column_b054c31f
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| ALTER TABLE `pa_drop_check_and_its_column_b054c31f`.`t` DROP CHECK `ck`
< ok 0 0
> exec
| ALTER TABLE `pa_drop_check_and_its_column_b054c31f`.`t` DROP COLUMN `a`
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_check_and_its_column_b054c31f.t
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

== TestOracle_CheckMigrationApplyCorrectness/8.0/drop-named-check default#1
> exec
| DROP DATABASE IF EXISTS pa_drop_named_check_f87bf771_to
< ok 0 0
> exec
| CREATE DATABASE pa_drop_named_check_f87bf771_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_named_check_f87bf771_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_named_check_f87bf771_to.t
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
| DROP DATABASE IF EXISTS pa_drop_named_check_f87bf771_from
< ok 0 0
> exec
| CREATE DATABASE pa_drop_named_check_f87bf771_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_named_check_f87bf771_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_named_check_f87bf771_from.t
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
|   CONSTRAINT `ck` CHECK ((`a` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_drop_named_check_f87bf771
< ok 0 0
> exec
| CREATE DATABASE pa_drop_named_check_f87bf771 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_named_check_f87bf771
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| ALTER TABLE `pa_drop_named_check_f87bf771`.`t` DROP CHECK `ck`
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_named_check_f87bf771.t
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

== TestOracle_CheckMigrationApplyCorrectness/8.0/drop-one-of-many default#1
> exec
| DROP DATABASE IF EXISTS pa_drop_one_of_many_e774b8b_to
< ok 0 0
> exec
| CREATE DATABASE pa_drop_one_of_many_e774b8b_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_one_of_many_e774b8b_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT c2 CHECK (b < 100))
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_one_of_many_e774b8b_to.t
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
|   CONSTRAINT `c2` CHECK ((`b` < 100))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_drop_one_of_many_e774b8b_from
< ok 0 0
> exec
| CREATE DATABASE pa_drop_one_of_many_e774b8b_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_one_of_many_e774b8b_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT c1 CHECK (a > 0), CONSTRAINT c2 CHECK (b < 100))
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_one_of_many_e774b8b_from.t
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
|   CONSTRAINT `c1` CHECK ((`a` > 0)),
|   CONSTRAINT `c2` CHECK ((`b` < 100))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_drop_one_of_many_e774b8b
< ok 0 0
> exec
| CREATE DATABASE pa_drop_one_of_many_e774b8b DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_one_of_many_e774b8b
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT c1 CHECK (a > 0), CONSTRAINT c2 CHECK (b < 100))
< ok 0 0
> exec
| ALTER TABLE `pa_drop_one_of_many_e774b8b`.`t` DROP CHECK `c1`
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_one_of_many_e774b8b.t
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
|   CONSTRAINT `c2` CHECK ((`b` < 100))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_CheckMigrationApplyCorrectness/8.0/drop-unnamed-check default#1
> exec
| DROP DATABASE IF EXISTS pa_drop_unnamed_check_58051f50_to
< ok 0 0
> exec
| CREATE DATABASE pa_drop_unnamed_check_58051f50_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_unnamed_check_58051f50_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_unnamed_check_58051f50_to.t
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
| DROP DATABASE IF EXISTS pa_drop_unnamed_check_58051f50_from
< ok 0 0
> exec
| CREATE DATABASE pa_drop_unnamed_check_58051f50_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_unnamed_check_58051f50_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CHECK (a > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_unnamed_check_58051f50_from.t
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
|   CONSTRAINT `t_chk_1` CHECK ((`a` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_drop_unnamed_check_58051f50
< ok 0 0
> exec
| CREATE DATABASE pa_drop_unnamed_check_58051f50 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_unnamed_check_58051f50
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CHECK (a > 0))
< ok 0 0
> exec
| ALTER TABLE `pa_drop_unnamed_check_58051f50`.`t` DROP CHECK `t_chk_1`
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_unnamed_check_58051f50.t
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

== TestOracle_CheckMigrationApplyCorrectness/8.0/modify-enforced-to-not default#1
> exec
| DROP DATABASE IF EXISTS pa_modify_enforced_to_not_f3d92732_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_enforced_to_not_f3d92732_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_enforced_to_not_f3d92732_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT ck CHECK (a < b) NOT ENFORCED)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_enforced_to_not_f3d92732_to.t
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
|   CONSTRAINT `ck` CHECK ((`a` < `b`)) /*!80016 NOT ENFORCED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_enforced_to_not_f3d92732_from
< ok 0 0
> exec
| CREATE DATABASE pa_modify_enforced_to_not_f3d92732_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_enforced_to_not_f3d92732_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT ck CHECK (a < b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_enforced_to_not_f3d92732_from.t
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
|   CONSTRAINT `ck` CHECK ((`a` < `b`))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_enforced_to_not_f3d92732
< ok 0 0
> exec
| CREATE DATABASE pa_modify_enforced_to_not_f3d92732 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_enforced_to_not_f3d92732
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT ck CHECK (a < b))
< ok 0 0
> exec
| ALTER TABLE `pa_modify_enforced_to_not_f3d92732`.`t` DROP CHECK `ck`
< ok 0 0
> exec
| ALTER TABLE `pa_modify_enforced_to_not_f3d92732`.`t` ADD CONSTRAINT `ck` CHECK ((`a` < `b`)) /*!80016 NOT ENFORCED */
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_enforced_to_not_f3d92732.t
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
|   CONSTRAINT `ck` CHECK ((`a` < `b`)) /*!80016 NOT ENFORCED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_CheckMigrationApplyCorrectness/8.0/modify-expr default#1
> exec
| DROP DATABASE IF EXISTS pa_modify_expr_9badc381_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_expr_9badc381_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_expr_9badc381_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT ck CHECK (a > 10))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_expr_9badc381_to.t
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
|   CONSTRAINT `ck` CHECK ((`a` > 10))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_expr_9badc381_from
< ok 0 0
> exec
| CREATE DATABASE pa_modify_expr_9badc381_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_expr_9badc381_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_expr_9badc381_from.t
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
|   CONSTRAINT `ck` CHECK ((`a` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_expr_9badc381
< ok 0 0
> exec
| CREATE DATABASE pa_modify_expr_9badc381 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_expr_9badc381
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| ALTER TABLE `pa_modify_expr_9badc381`.`t` DROP CHECK `ck`
< ok 0 0
> exec
| ALTER TABLE `pa_modify_expr_9badc381`.`t` ADD CONSTRAINT `ck` CHECK ((`a` > 10))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_expr_9badc381.t
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
|   CONSTRAINT `ck` CHECK ((`a` > 10))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_CheckMigrationApplyCorrectness/8.0/modify-not-to-enforced default#1
> exec
| DROP DATABASE IF EXISTS pa_modify_not_to_enforced_aa7daef0_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_not_to_enforced_aa7daef0_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_not_to_enforced_aa7daef0_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT ck CHECK (a < b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_not_to_enforced_aa7daef0_to.t
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
|   CONSTRAINT `ck` CHECK ((`a` < `b`))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_not_to_enforced_aa7daef0_from
< ok 0 0
> exec
| CREATE DATABASE pa_modify_not_to_enforced_aa7daef0_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_not_to_enforced_aa7daef0_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT ck CHECK (a < b) NOT ENFORCED)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_not_to_enforced_aa7daef0_from.t
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
|   CONSTRAINT `ck` CHECK ((`a` < `b`)) /*!80016 NOT ENFORCED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_not_to_enforced_aa7daef0
< ok 0 0
> exec
| CREATE DATABASE pa_modify_not_to_enforced_aa7daef0 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_not_to_enforced_aa7daef0
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT ck CHECK (a < b) NOT ENFORCED)
< ok 0 0
> exec
| ALTER TABLE `pa_modify_not_to_enforced_aa7daef0`.`t` DROP CHECK `ck`
< ok 0 0
> exec
| ALTER TABLE `pa_modify_not_to_enforced_aa7daef0`.`t` ADD CONSTRAINT `ck` CHECK ((`a` < `b`))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_not_to_enforced_aa7daef0.t
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
|   CONSTRAINT `ck` CHECK ((`a` < `b`))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
