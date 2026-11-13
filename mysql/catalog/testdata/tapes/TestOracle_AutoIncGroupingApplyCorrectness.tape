-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracle_AutoIncGroupingApplyCorrectness default#1

== TestOracle_AutoIncGroupingApplyCorrectness/8.0/aig-add-col-composite-backing default#1
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_composite_backing_1f3769cf_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_composite_backing_1f3769cf_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_composite_backing_1f3769cf_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, seq BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, d INT NOT NULL, UNIQUE KEY uk_sd (seq, d))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_composite_backing_1f3769cf_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `seq` bigint unsigned NOT NULL AUTO_INCREMENT,
|   `d` int NOT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_sd` (`seq`,`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_composite_backing_1f3769cf_from
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_composite_backing_1f3769cf_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_composite_backing_1f3769cf_from
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_composite_backing_1f3769cf_from.t
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
| DROP DATABASE IF EXISTS pa_aig_add_col_composite_backing_1f3769cf
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_composite_backing_1f3769cf DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_composite_backing_1f3769cf
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE `pa_aig_add_col_composite_backing_1f3769cf`.`t` ADD COLUMN `seq` bigint unsigned NOT NULL AUTO_INCREMENT, ADD COLUMN `d` int NOT NULL, ADD UNIQUE KEY `uk_sd` (`seq`,`d`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_composite_backing_1f3769cf.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `seq` bigint unsigned NOT NULL AUTO_INCREMENT,
|   `d` int NOT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_sd` (`seq`,`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_AutoIncGroupingApplyCorrectness/8.0/aig-add-col-desc-key default#1
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_desc_key_8ef498cb_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_desc_key_8ef498cb_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_desc_key_8ef498cb_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, seq INT NOT NULL AUTO_INCREMENT, KEY k_seq (seq DESC))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_desc_key_8ef498cb_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `seq` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k_seq` (`seq` DESC)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_desc_key_8ef498cb_from
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_desc_key_8ef498cb_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_desc_key_8ef498cb_from
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_desc_key_8ef498cb_from.t
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
| DROP DATABASE IF EXISTS pa_aig_add_col_desc_key_8ef498cb
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_desc_key_8ef498cb DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_desc_key_8ef498cb
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE `pa_aig_add_col_desc_key_8ef498cb`.`t` ADD COLUMN `seq` int NOT NULL AUTO_INCREMENT, ADD KEY `k_seq` (`seq` DESC)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_desc_key_8ef498cb.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `seq` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k_seq` (`seq` DESC)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_AutoIncGroupingApplyCorrectness/8.0/aig-add-col-invisible-key default#1
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_invisible_key_effc39b9_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_invisible_key_effc39b9_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_invisible_key_effc39b9_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, seq INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_seq (seq) INVISIBLE)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_invisible_key_effc39b9_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `seq` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_seq` (`seq`) /*!80000 INVISIBLE */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_invisible_key_effc39b9_from
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_invisible_key_effc39b9_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_invisible_key_effc39b9_from
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_invisible_key_effc39b9_from.t
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
| DROP DATABASE IF EXISTS pa_aig_add_col_invisible_key_effc39b9
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_invisible_key_effc39b9 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_invisible_key_effc39b9
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE `pa_aig_add_col_invisible_key_effc39b9`.`t` ADD COLUMN `seq` int NOT NULL AUTO_INCREMENT, ADD UNIQUE KEY `uk_seq` (`seq`) /*!80000 INVISIBLE */
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_invisible_key_effc39b9.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `seq` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_seq` (`seq`) /*!80000 INVISIBLE */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_AutoIncGroupingApplyCorrectness/8.0/aig-add-col-new-pk default#1
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_new_pk_7ae518e8_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_new_pk_7ae518e8_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_new_pk_7ae518e8_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, a INT, seq BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, PRIMARY KEY (seq))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_new_pk_7ae518e8_to.t
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
|   `seq` bigint unsigned NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`seq`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_new_pk_7ae518e8_from
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_new_pk_7ae518e8_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_new_pk_7ae518e8_from
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_new_pk_7ae518e8_from.t
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
| DROP DATABASE IF EXISTS pa_aig_add_col_new_pk_7ae518e8
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_new_pk_7ae518e8 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_new_pk_7ae518e8
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, a INT)
< ok 0 0
> exec
| ALTER TABLE `pa_aig_add_col_new_pk_7ae518e8`.`t` ADD COLUMN `seq` bigint unsigned NOT NULL AUTO_INCREMENT, ADD PRIMARY KEY (`seq`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_new_pk_7ae518e8.t
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
|   `seq` bigint unsigned NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`seq`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_AutoIncGroupingApplyCorrectness/8.0/aig-add-col-plain-key default#1
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_plain_key_4656419a_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_plain_key_4656419a_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_plain_key_4656419a_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT, seq BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, KEY k_seq (seq))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_plain_key_4656419a_to.t
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
|   `seq` bigint unsigned NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k_seq` (`seq`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_plain_key_4656419a_from
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_plain_key_4656419a_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_plain_key_4656419a_from
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_plain_key_4656419a_from.t
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
| DROP DATABASE IF EXISTS pa_aig_add_col_plain_key_4656419a
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_plain_key_4656419a DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_plain_key_4656419a
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT)
< ok 0 0
> exec
| ALTER TABLE `pa_aig_add_col_plain_key_4656419a`.`t` ADD COLUMN `seq` bigint unsigned NOT NULL AUTO_INCREMENT, ADD KEY `k_seq` (`seq`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_plain_key_4656419a.t
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
|   `seq` bigint unsigned NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k_seq` (`seq`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_AutoIncGroupingApplyCorrectness/8.0/aig-add-col-promote-pk default#1
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_promote_pk_fc2df8b6_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_promote_pk_fc2df8b6_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_promote_pk_fc2df8b6_to
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT, seq BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, PRIMARY KEY (seq))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_promote_pk_fc2df8b6_to.t
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
|   `seq` bigint unsigned NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`seq`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_promote_pk_fc2df8b6_from
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_promote_pk_fc2df8b6_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_promote_pk_fc2df8b6_from
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL PRIMARY KEY, b INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_promote_pk_fc2df8b6_from.t
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
| DROP DATABASE IF EXISTS pa_aig_add_col_promote_pk_fc2df8b6
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_promote_pk_fc2df8b6 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_promote_pk_fc2df8b6
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL PRIMARY KEY, b INT)
< ok 0 0
> exec
| ALTER TABLE `pa_aig_add_col_promote_pk_fc2df8b6`.`t` ADD COLUMN `seq` bigint unsigned NOT NULL AUTO_INCREMENT, DROP PRIMARY KEY, ADD PRIMARY KEY (`seq`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_promote_pk_fc2df8b6.t
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
|   `seq` bigint unsigned NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`seq`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_AutoIncGroupingApplyCorrectness/8.0/aig-add-col-promote-pk-drop-old default#1
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_promote_pk_drop_old_80fdc0f0_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_promote_pk_drop_old_80fdc0f0_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_promote_pk_drop_old_80fdc0f0_to
< ok 0 0
> exec
| CREATE TABLE t (b INT, seq BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, PRIMARY KEY (seq))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_promote_pk_drop_old_80fdc0f0_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int DEFAULT NULL,
|   `seq` bigint unsigned NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`seq`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_promote_pk_drop_old_80fdc0f0_from
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_promote_pk_drop_old_80fdc0f0_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_promote_pk_drop_old_80fdc0f0_from
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL PRIMARY KEY, b INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_promote_pk_drop_old_80fdc0f0_from.t
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
| DROP DATABASE IF EXISTS pa_aig_add_col_promote_pk_drop_old_80fdc0f0
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_promote_pk_drop_old_80fdc0f0 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_promote_pk_drop_old_80fdc0f0
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL PRIMARY KEY, b INT)
< ok 0 0
> exec
| ALTER TABLE `pa_aig_add_col_promote_pk_drop_old_80fdc0f0`.`t` DROP PRIMARY KEY
< ok 0 0
> exec
| ALTER TABLE `pa_aig_add_col_promote_pk_drop_old_80fdc0f0`.`t` DROP COLUMN `a`
< ok 0 0
> exec
| ALTER TABLE `pa_aig_add_col_promote_pk_drop_old_80fdc0f0`.`t` ADD COLUMN `seq` bigint unsigned NOT NULL AUTO_INCREMENT, ADD PRIMARY KEY (`seq`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_promote_pk_drop_old_80fdc0f0.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int DEFAULT NULL,
|   `seq` bigint unsigned NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`seq`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_AutoIncGroupingApplyCorrectness/8.0/aig-add-col-trailing-expr-key default#1
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_trailing_expr_key_13e9879e_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_trailing_expr_key_13e9879e_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_trailing_expr_key_13e9879e_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT, seq INT NOT NULL AUTO_INCREMENT, KEY k_se (seq, (a + 1)))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_trailing_expr_key_13e9879e_to.t
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
|   `seq` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k_se` (`seq`,((`a` + 1)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_trailing_expr_key_13e9879e_from
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_trailing_expr_key_13e9879e_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_trailing_expr_key_13e9879e_from
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_trailing_expr_key_13e9879e_from.t
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
| DROP DATABASE IF EXISTS pa_aig_add_col_trailing_expr_key_13e9879e
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_trailing_expr_key_13e9879e DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_trailing_expr_key_13e9879e
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT)
< ok 0 0
> exec
| ALTER TABLE `pa_aig_add_col_trailing_expr_key_13e9879e`.`t` ADD COLUMN `seq` int NOT NULL AUTO_INCREMENT, ADD KEY `k_se` (`seq`,((`a` + 1)))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_trailing_expr_key_13e9879e.t
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
|   `seq` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k_se` (`seq`,((`a` + 1)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_AutoIncGroupingApplyCorrectness/8.0/aig-add-col-uk default#1
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_uk_5e172a84_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_uk_5e172a84_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_uk_5e172a84_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT, seq BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_seq (seq))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_uk_5e172a84_to.t
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
|   `seq` bigint unsigned NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_seq` (`seq`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_uk_5e172a84_from
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_uk_5e172a84_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_uk_5e172a84_from
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_uk_5e172a84_from.t
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
| DROP DATABASE IF EXISTS pa_aig_add_col_uk_5e172a84
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_uk_5e172a84 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_uk_5e172a84
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT)
< ok 0 0
> exec
| ALTER TABLE `pa_aig_add_col_uk_5e172a84`.`t` ADD COLUMN `seq` bigint unsigned NOT NULL AUTO_INCREMENT, ADD UNIQUE KEY `uk_seq` (`seq`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_uk_5e172a84.t
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
|   `seq` bigint unsigned NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_seq` (`seq`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_AutoIncGroupingApplyCorrectness/8.0/aig-add-col-with-other-changes default#1
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_with_other_changes_bab13057_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_with_other_changes_bab13057_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_with_other_changes_bab13057_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT, e VARCHAR(20), f INT, seq BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_seq (seq))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_with_other_changes_bab13057_to.t
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
|   `e` varchar(20) DEFAULT NULL,
|   `f` int DEFAULT NULL,
|   `seq` bigint unsigned NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_seq` (`seq`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_with_other_changes_bab13057_from
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_with_other_changes_bab13057_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_with_other_changes_bab13057_from
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT, e VARCHAR(10))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_with_other_changes_bab13057_from.t
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
|   `e` varchar(10) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_with_other_changes_bab13057
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_with_other_changes_bab13057 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_with_other_changes_bab13057
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT, e VARCHAR(10))
< ok 0 0
> exec
| ALTER TABLE `pa_aig_add_col_with_other_changes_bab13057`.`t` MODIFY COLUMN `e` varchar(20) DEFAULT NULL
< ok 0 0
> exec
| ALTER TABLE `pa_aig_add_col_with_other_changes_bab13057`.`t` ADD COLUMN `f` int DEFAULT NULL
< ok 0 0
> exec
| ALTER TABLE `pa_aig_add_col_with_other_changes_bab13057`.`t` ADD COLUMN `seq` bigint unsigned NOT NULL AUTO_INCREMENT, ADD UNIQUE KEY `uk_seq` (`seq`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_with_other_changes_bab13057.t
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
|   `e` varchar(20) DEFAULT NULL,
|   `f` int DEFAULT NULL,
|   `seq` bigint unsigned NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_seq` (`seq`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_AutoIncGroupingApplyCorrectness/8.0/aig-ai-migrate-columns default#1
> exec
| DROP DATABASE IF EXISTS pa_aig_ai_migrate_columns_a6cb228b_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_ai_migrate_columns_a6cb228b_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_ai_migrate_columns_a6cb228b_to
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_a (a), UNIQUE KEY uk_b (b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_ai_migrate_columns_a6cb228b_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int NOT NULL,
|   `b` int NOT NULL AUTO_INCREMENT,
|   UNIQUE KEY `uk_a` (`a`),
|   UNIQUE KEY `uk_b` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_ai_migrate_columns_a6cb228b_from
< ok 0 0
> exec
| CREATE DATABASE pa_aig_ai_migrate_columns_a6cb228b_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_ai_migrate_columns_a6cb228b_from
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL AUTO_INCREMENT, b INT NOT NULL, UNIQUE KEY uk_a (a), UNIQUE KEY uk_b (b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_ai_migrate_columns_a6cb228b_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int NOT NULL AUTO_INCREMENT,
|   `b` int NOT NULL,
|   UNIQUE KEY `uk_a` (`a`),
|   UNIQUE KEY `uk_b` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_ai_migrate_columns_a6cb228b
< ok 0 0
> exec
| CREATE DATABASE pa_aig_ai_migrate_columns_a6cb228b DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_ai_migrate_columns_a6cb228b
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL AUTO_INCREMENT, b INT NOT NULL, UNIQUE KEY uk_a (a), UNIQUE KEY uk_b (b))
< ok 0 0
> exec
| ALTER TABLE `pa_aig_ai_migrate_columns_a6cb228b`.`t` MODIFY COLUMN `a` int NOT NULL, MODIFY COLUMN `b` int NOT NULL AUTO_INCREMENT
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_ai_migrate_columns_a6cb228b.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int NOT NULL,
|   `b` int NOT NULL AUTO_INCREMENT,
|   UNIQUE KEY `uk_a` (`a`),
|   UNIQUE KEY `uk_b` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_AutoIncGroupingApplyCorrectness/8.0/aig-ai-migrate-to-new-col default#1
> exec
| DROP DATABASE IF EXISTS pa_aig_ai_migrate_to_new_col_571ea041_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_ai_migrate_to_new_col_571ea041_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_ai_migrate_to_new_col_571ea041_to
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_a (a), UNIQUE KEY uk_b (b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_ai_migrate_to_new_col_571ea041_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int NOT NULL,
|   `b` int NOT NULL AUTO_INCREMENT,
|   UNIQUE KEY `uk_a` (`a`),
|   UNIQUE KEY `uk_b` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_ai_migrate_to_new_col_571ea041_from
< ok 0 0
> exec
| CREATE DATABASE pa_aig_ai_migrate_to_new_col_571ea041_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_ai_migrate_to_new_col_571ea041_from
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_a (a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_ai_migrate_to_new_col_571ea041_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int NOT NULL AUTO_INCREMENT,
|   UNIQUE KEY `uk_a` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_ai_migrate_to_new_col_571ea041
< ok 0 0
> exec
| CREATE DATABASE pa_aig_ai_migrate_to_new_col_571ea041 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_ai_migrate_to_new_col_571ea041
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_a (a))
< ok 0 0
> exec
| ALTER TABLE `pa_aig_ai_migrate_to_new_col_571ea041`.`t` MODIFY COLUMN `a` int NOT NULL, ADD COLUMN `b` int NOT NULL AUTO_INCREMENT, ADD UNIQUE KEY `uk_b` (`b`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_ai_migrate_to_new_col_571ea041.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int NOT NULL,
|   `b` int NOT NULL AUTO_INCREMENT,
|   UNIQUE KEY `uk_a` (`a`),
|   UNIQUE KEY `uk_b` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_AutoIncGroupingApplyCorrectness/8.0/aig-drop-composite-pk-and-column default#1
> exec
| DROP DATABASE IF EXISTS pa_aig_drop_composite_pk_and_column_44c4ceaf_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_drop_composite_pk_and_column_44c4ceaf_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_drop_composite_pk_and_column_44c4ceaf_to
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_drop_composite_pk_and_column_44c4ceaf_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_drop_composite_pk_and_column_44c4ceaf_from
< ok 0 0
> exec
| CREATE DATABASE pa_aig_drop_composite_pk_and_column_44c4ceaf_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_drop_composite_pk_and_column_44c4ceaf_from
< ok 0 0
> exec
| CREATE TABLE t (c INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, PRIMARY KEY (c, a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_drop_composite_pk_and_column_44c4ceaf_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `c` int NOT NULL AUTO_INCREMENT,
|   `a` int NOT NULL,
|   PRIMARY KEY (`c`,`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_drop_composite_pk_and_column_44c4ceaf
< ok 0 0
> exec
| CREATE DATABASE pa_aig_drop_composite_pk_and_column_44c4ceaf DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_drop_composite_pk_and_column_44c4ceaf
< ok 0 0
> exec
| CREATE TABLE t (c INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, PRIMARY KEY (c, a))
< ok 0 0
> exec
| ALTER TABLE `pa_aig_drop_composite_pk_and_column_44c4ceaf`.`t` DROP PRIMARY KEY, DROP COLUMN `c`
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_drop_composite_pk_and_column_44c4ceaf.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_AutoIncGroupingApplyCorrectness/8.0/aig-drop-covered-elsewhere default#1
> exec
| DROP DATABASE IF EXISTS pa_aig_drop_covered_elsewhere_a08f6790_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_drop_covered_elsewhere_a08f6790_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_drop_covered_elsewhere_a08f6790_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, KEY k2 (c))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_drop_covered_elsewhere_a08f6790_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `c` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k2` (`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_drop_covered_elsewhere_a08f6790_from
< ok 0 0
> exec
| CREATE DATABASE pa_aig_drop_covered_elsewhere_a08f6790_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_drop_covered_elsewhere_a08f6790_from
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk1 (c), KEY k2 (c))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_drop_covered_elsewhere_a08f6790_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `c` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk1` (`c`),
|   KEY `k2` (`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_drop_covered_elsewhere_a08f6790
< ok 0 0
> exec
| CREATE DATABASE pa_aig_drop_covered_elsewhere_a08f6790 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_drop_covered_elsewhere_a08f6790
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk1 (c), KEY k2 (c))
< ok 0 0
> exec
| ALTER TABLE `pa_aig_drop_covered_elsewhere_a08f6790`.`t` DROP INDEX `uk1`
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_drop_covered_elsewhere_a08f6790.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `c` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k2` (`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_AutoIncGroupingApplyCorrectness/8.0/aig-drop-key-and-column default#1
> exec
| DROP DATABASE IF EXISTS pa_aig_drop_key_and_column_87fffe21_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_drop_key_and_column_87fffe21_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_drop_key_and_column_87fffe21_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_drop_key_and_column_87fffe21_to.t
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
| DROP DATABASE IF EXISTS pa_aig_drop_key_and_column_87fffe21_from
< ok 0 0
> exec
| CREATE DATABASE pa_aig_drop_key_and_column_87fffe21_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_drop_key_and_column_87fffe21_from
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_c (c))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_drop_key_and_column_87fffe21_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `c` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_c` (`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_drop_key_and_column_87fffe21
< ok 0 0
> exec
| CREATE DATABASE pa_aig_drop_key_and_column_87fffe21 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_drop_key_and_column_87fffe21
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_c (c))
< ok 0 0
> exec
| ALTER TABLE `pa_aig_drop_key_and_column_87fffe21`.`t` DROP INDEX `uk_c`, DROP COLUMN `c`
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_drop_key_and_column_87fffe21.t
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

== TestOracle_AutoIncGroupingApplyCorrectness/8.0/aig-drop-key-deai default#1
> exec
| DROP DATABASE IF EXISTS pa_aig_drop_key_deai_38e4ec0e_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_drop_key_deai_38e4ec0e_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_drop_key_deai_38e4ec0e_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_drop_key_deai_38e4ec0e_to.t
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
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_drop_key_deai_38e4ec0e_from
< ok 0 0
> exec
| CREATE DATABASE pa_aig_drop_key_deai_38e4ec0e_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_drop_key_deai_38e4ec0e_from
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_c (c))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_drop_key_deai_38e4ec0e_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `c` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_c` (`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_drop_key_deai_38e4ec0e
< ok 0 0
> exec
| CREATE DATABASE pa_aig_drop_key_deai_38e4ec0e DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_drop_key_deai_38e4ec0e
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_c (c))
< ok 0 0
> exec
| ALTER TABLE `pa_aig_drop_key_deai_38e4ec0e`.`t` MODIFY COLUMN `c` int NOT NULL, DROP INDEX `uk_c`
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_drop_key_deai_38e4ec0e.t
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
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_AutoIncGroupingApplyCorrectness/8.0/aig-drop-pk-and-column default#1
> exec
| DROP DATABASE IF EXISTS pa_aig_drop_pk_and_column_95158c61_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_drop_pk_and_column_95158c61_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_drop_pk_and_column_95158c61_to
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_drop_pk_and_column_95158c61_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_drop_pk_and_column_95158c61_from
< ok 0 0
> exec
| CREATE DATABASE pa_aig_drop_pk_and_column_95158c61_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_drop_pk_and_column_95158c61_from
< ok 0 0
> exec
| CREATE TABLE t (c INT NOT NULL AUTO_INCREMENT PRIMARY KEY, a INT NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_drop_pk_and_column_95158c61_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `c` int NOT NULL AUTO_INCREMENT,
|   `a` int NOT NULL,
|   PRIMARY KEY (`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_drop_pk_and_column_95158c61
< ok 0 0
> exec
| CREATE DATABASE pa_aig_drop_pk_and_column_95158c61 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_drop_pk_and_column_95158c61
< ok 0 0
> exec
| CREATE TABLE t (c INT NOT NULL AUTO_INCREMENT PRIMARY KEY, a INT NOT NULL)
< ok 0 0
> exec
| ALTER TABLE `pa_aig_drop_pk_and_column_95158c61`.`t` DROP PRIMARY KEY, DROP COLUMN `c`
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_drop_pk_and_column_95158c61.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_AutoIncGroupingApplyCorrectness/8.0/aig-drop-pk-deai default#1
> exec
| DROP DATABASE IF EXISTS pa_aig_drop_pk_deai_31c4aace_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_drop_pk_deai_31c4aace_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_drop_pk_deai_31c4aace_to
< ok 0 0
> exec
| CREATE TABLE t (c INT NOT NULL, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_drop_pk_deai_31c4aace_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `c` int NOT NULL,
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_drop_pk_deai_31c4aace_from
< ok 0 0
> exec
| CREATE DATABASE pa_aig_drop_pk_deai_31c4aace_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_drop_pk_deai_31c4aace_from
< ok 0 0
> exec
| CREATE TABLE t (c INT NOT NULL AUTO_INCREMENT PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_drop_pk_deai_31c4aace_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `c` int NOT NULL AUTO_INCREMENT,
|   `a` int DEFAULT NULL,
|   PRIMARY KEY (`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_drop_pk_deai_31c4aace
< ok 0 0
> exec
| CREATE DATABASE pa_aig_drop_pk_deai_31c4aace DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_drop_pk_deai_31c4aace
< ok 0 0
> exec
| CREATE TABLE t (c INT NOT NULL AUTO_INCREMENT PRIMARY KEY, a INT)
< ok 0 0
> exec
| ALTER TABLE `pa_aig_drop_pk_deai_31c4aace`.`t` MODIFY COLUMN `c` int NOT NULL, DROP PRIMARY KEY
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_drop_pk_deai_31c4aace.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `c` int NOT NULL,
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_AutoIncGroupingApplyCorrectness/8.0/aig-funcpart-new-col default#1
> exec
| DROP DATABASE IF EXISTS pa_aig_funcpart_new_col_9ad57dac_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_funcpart_new_col_9ad57dac_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_funcpart_new_col_9ad57dac_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, seq INT NOT NULL AUTO_INCREMENT, z INT, KEY k_sz (seq, (z + 1)))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_funcpart_new_col_9ad57dac_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `seq` int NOT NULL AUTO_INCREMENT,
|   `z` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k_sz` (`seq`,((`z` + 1)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_funcpart_new_col_9ad57dac_from
< ok 0 0
> exec
| CREATE DATABASE pa_aig_funcpart_new_col_9ad57dac_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_funcpart_new_col_9ad57dac_from
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_funcpart_new_col_9ad57dac_from.t
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
| DROP DATABASE IF EXISTS pa_aig_funcpart_new_col_9ad57dac
< ok 0 0
> exec
| CREATE DATABASE pa_aig_funcpart_new_col_9ad57dac DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_funcpart_new_col_9ad57dac
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE `pa_aig_funcpart_new_col_9ad57dac`.`t` ADD COLUMN `seq` int NOT NULL AUTO_INCREMENT, ADD COLUMN `z` int DEFAULT NULL, ADD KEY `k_sz` (`seq`,((`z` + 1)))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_funcpart_new_col_9ad57dac.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `seq` int NOT NULL AUTO_INCREMENT,
|   `z` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k_sz` (`seq`,((`z` + 1)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_AutoIncGroupingApplyCorrectness/8.0/aig-generated-keypart-pullin default#1
> exec
| DROP DATABASE IF EXISTS pa_aig_generated_keypart_pullin_40dea58e_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_generated_keypart_pullin_40dea58e_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_generated_keypart_pullin_40dea58e_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, seq INT NOT NULL AUTO_INCREMENT, d INT NOT NULL, g INT GENERATED ALWAYS AS (d+1) VIRTUAL, UNIQUE KEY uk_sg (seq, g))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_generated_keypart_pullin_40dea58e_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `seq` int NOT NULL AUTO_INCREMENT,
|   `d` int NOT NULL,
|   `g` int GENERATED ALWAYS AS ((`d` + 1)) VIRTUAL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_sg` (`seq`,`g`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_generated_keypart_pullin_40dea58e_from
< ok 0 0
> exec
| CREATE DATABASE pa_aig_generated_keypart_pullin_40dea58e_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_generated_keypart_pullin_40dea58e_from
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_generated_keypart_pullin_40dea58e_from.t
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
| DROP DATABASE IF EXISTS pa_aig_generated_keypart_pullin_40dea58e
< ok 0 0
> exec
| CREATE DATABASE pa_aig_generated_keypart_pullin_40dea58e DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_generated_keypart_pullin_40dea58e
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE `pa_aig_generated_keypart_pullin_40dea58e`.`t` ADD COLUMN `seq` int NOT NULL AUTO_INCREMENT, ADD COLUMN `d` int NOT NULL, ADD COLUMN `g` int GENERATED ALWAYS AS (((`d` + 1))) VIRTUAL, ADD UNIQUE KEY `uk_sg` (`seq`,`g`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_generated_keypart_pullin_40dea58e.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `seq` int NOT NULL AUTO_INCREMENT,
|   `d` int NOT NULL,
|   `g` int GENERATED ALWAYS AS ((`d` + 1)) VIRTUAL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_sg` (`seq`,`g`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_AutoIncGroupingApplyCorrectness/8.0/aig-genpart-base-widened default#1
> exec
| DROP DATABASE IF EXISTS pa_aig_genpart_base_widened_1edcf863_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_genpart_base_widened_1edcf863_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_genpart_base_widened_1edcf863_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, z BIGINT, seq INT NOT NULL AUTO_INCREMENT, g INT GENERATED ALWAYS AS (z + 1) VIRTUAL, UNIQUE KEY uk_sg (seq, g))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_genpart_base_widened_1edcf863_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `z` bigint DEFAULT NULL,
|   `seq` int NOT NULL AUTO_INCREMENT,
|   `g` int GENERATED ALWAYS AS ((`z` + 1)) VIRTUAL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_sg` (`seq`,`g`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_genpart_base_widened_1edcf863_from
< ok 0 0
> exec
| CREATE DATABASE pa_aig_genpart_base_widened_1edcf863_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_genpart_base_widened_1edcf863_from
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, z INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_genpart_base_widened_1edcf863_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `z` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_genpart_base_widened_1edcf863
< ok 0 0
> exec
| CREATE DATABASE pa_aig_genpart_base_widened_1edcf863 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_genpart_base_widened_1edcf863
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, z INT)
< ok 0 0
> exec
| ALTER TABLE `pa_aig_genpart_base_widened_1edcf863`.`t` ADD COLUMN `seq` int NOT NULL AUTO_INCREMENT, ADD COLUMN `g` int GENERATED ALWAYS AS (((`z` + 1))) VIRTUAL, ADD UNIQUE KEY `uk_sg` (`seq`,`g`)
< ok 0 0
> exec
| ALTER TABLE `pa_aig_genpart_base_widened_1edcf863`.`t` MODIFY COLUMN `z` bigint DEFAULT NULL
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_genpart_base_widened_1edcf863.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `z` bigint DEFAULT NULL,
|   `seq` int NOT NULL AUTO_INCREMENT,
|   `g` int GENERATED ALWAYS AS ((`z` + 1)) VIRTUAL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_sg` (`seq`,`g`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_AutoIncGroupingApplyCorrectness/8.0/aig-modify-gain-ai-existing-key default#1
> exec
| DROP DATABASE IF EXISTS pa_aig_modify_gain_ai_existing_key_46160758_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_modify_gain_ai_existing_key_46160758_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_modify_gain_ai_existing_key_46160758_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_c (c))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_modify_gain_ai_existing_key_46160758_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `c` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_c` (`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_modify_gain_ai_existing_key_46160758_from
< ok 0 0
> exec
| CREATE DATABASE pa_aig_modify_gain_ai_existing_key_46160758_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_modify_gain_ai_existing_key_46160758_from
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL, UNIQUE KEY uk_c (c))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_modify_gain_ai_existing_key_46160758_from.t
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
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_c` (`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_modify_gain_ai_existing_key_46160758
< ok 0 0
> exec
| CREATE DATABASE pa_aig_modify_gain_ai_existing_key_46160758 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_modify_gain_ai_existing_key_46160758
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL, UNIQUE KEY uk_c (c))
< ok 0 0
> exec
| ALTER TABLE `pa_aig_modify_gain_ai_existing_key_46160758`.`t` MODIFY COLUMN `c` int NOT NULL AUTO_INCREMENT
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_modify_gain_ai_existing_key_46160758.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `c` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_c` (`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_AutoIncGroupingApplyCorrectness/8.0/aig-modify-gain-ai-new-key default#1
> exec
| DROP DATABASE IF EXISTS pa_aig_modify_gain_ai_new_key_b6e21499_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_modify_gain_ai_new_key_b6e21499_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_modify_gain_ai_new_key_b6e21499_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_c (c))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_modify_gain_ai_new_key_b6e21499_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `c` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_c` (`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_modify_gain_ai_new_key_b6e21499_from
< ok 0 0
> exec
| CREATE DATABASE pa_aig_modify_gain_ai_new_key_b6e21499_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_modify_gain_ai_new_key_b6e21499_from
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_modify_gain_ai_new_key_b6e21499_from.t
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
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_modify_gain_ai_new_key_b6e21499
< ok 0 0
> exec
| CREATE DATABASE pa_aig_modify_gain_ai_new_key_b6e21499 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_modify_gain_ai_new_key_b6e21499
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL)
< ok 0 0
> exec
| ALTER TABLE `pa_aig_modify_gain_ai_new_key_b6e21499`.`t` MODIFY COLUMN `c` int NOT NULL AUTO_INCREMENT, ADD UNIQUE KEY `uk_c` (`c`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_modify_gain_ai_new_key_b6e21499.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `c` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_c` (`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_AutoIncGroupingApplyCorrectness/8.0/aig-multi-drop-plus-add default#1
> exec
| DROP DATABASE IF EXISTS pa_aig_multi_drop_plus_add_116c0943_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_multi_drop_plus_add_116c0943_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_multi_drop_plus_add_116c0943_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk3 (c))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_multi_drop_plus_add_116c0943_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `c` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk3` (`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_multi_drop_plus_add_116c0943_from
< ok 0 0
> exec
| CREATE DATABASE pa_aig_multi_drop_plus_add_116c0943_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_multi_drop_plus_add_116c0943_from
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk1 (c), KEY k2 (c))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_multi_drop_plus_add_116c0943_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `c` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk1` (`c`),
|   KEY `k2` (`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_multi_drop_plus_add_116c0943
< ok 0 0
> exec
| CREATE DATABASE pa_aig_multi_drop_plus_add_116c0943 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_multi_drop_plus_add_116c0943
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk1 (c), KEY k2 (c))
< ok 0 0
> exec
| ALTER TABLE `pa_aig_multi_drop_plus_add_116c0943`.`t` DROP INDEX `k2`, DROP INDEX `uk1`, ADD UNIQUE KEY `uk3` (`c`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_multi_drop_plus_add_116c0943.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `c` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk3` (`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_AutoIncGroupingApplyCorrectness/8.0/aig-myisam-add-nonleftmost default#1
> exec
| DROP DATABASE IF EXISTS pa_aig_myisam_add_nonleftmost_9f4f801_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_myisam_add_nonleftmost_9f4f801_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_myisam_add_nonleftmost_9f4f801_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, seq INT NOT NULL AUTO_INCREMENT, KEY k_as (a, seq)) ENGINE=MyISAM
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_myisam_add_nonleftmost_9f4f801_to.t
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
|   `seq` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k_as` (`a`,`seq`)
| ) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_myisam_add_nonleftmost_9f4f801_from
< ok 0 0
> exec
| CREATE DATABASE pa_aig_myisam_add_nonleftmost_9f4f801_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_myisam_add_nonleftmost_9f4f801_from
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT NOT NULL) ENGINE=MyISAM
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_myisam_add_nonleftmost_9f4f801_from.t
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
|   PRIMARY KEY (`id`)
| ) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_myisam_add_nonleftmost_9f4f801
< ok 0 0
> exec
| CREATE DATABASE pa_aig_myisam_add_nonleftmost_9f4f801 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_myisam_add_nonleftmost_9f4f801
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT NOT NULL) ENGINE=MyISAM
< ok 0 0
> exec
| ALTER TABLE `pa_aig_myisam_add_nonleftmost_9f4f801`.`t` ADD COLUMN `seq` int NOT NULL AUTO_INCREMENT, ADD KEY `k_as` (`a`,`seq`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_myisam_add_nonleftmost_9f4f801.t
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
|   `seq` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k_as` (`a`,`seq`)
| ) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_AutoIncGroupingApplyCorrectness/8.0/aig-myisam-gain-ai-existing-nonleftmost default#1
> exec
| DROP DATABASE IF EXISTS pa_aig_myisam_gain_ai_existing_nonleftmost_280cdca4_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_myisam_gain_ai_existing_nonleftmost_280cdca4_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_myisam_gain_ai_existing_nonleftmost_280cdca4_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, c INT NOT NULL AUTO_INCREMENT, KEY k (a, c)) ENGINE=MyISAM
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_myisam_gain_ai_existing_nonleftmost_280cdca4_to.t
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
|   `c` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k` (`a`,`c`)
| ) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_myisam_gain_ai_existing_nonleftmost_280cdca4_from
< ok 0 0
> exec
| CREATE DATABASE pa_aig_myisam_gain_ai_existing_nonleftmost_280cdca4_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_myisam_gain_ai_existing_nonleftmost_280cdca4_from
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, c INT NOT NULL, KEY k (a, c)) ENGINE=MyISAM
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_myisam_gain_ai_existing_nonleftmost_280cdca4_from.t
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
|   `c` int NOT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`a`,`c`)
| ) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_myisam_gain_ai_existing_nonleftmost_280cdca4
< ok 0 0
> exec
| CREATE DATABASE pa_aig_myisam_gain_ai_existing_nonleftmost_280cdca4 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_myisam_gain_ai_existing_nonleftmost_280cdca4
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, c INT NOT NULL, KEY k (a, c)) ENGINE=MyISAM
< ok 0 0
> exec
| ALTER TABLE `pa_aig_myisam_gain_ai_existing_nonleftmost_280cdca4`.`t` MODIFY COLUMN `c` int NOT NULL AUTO_INCREMENT
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_myisam_gain_ai_existing_nonleftmost_280cdca4.t
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
|   `c` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k` (`a`,`c`)
| ) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_AutoIncGroupingApplyCorrectness/8.0/aig-myisam-shared-key-ai-migrate default#1
> exec
| DROP DATABASE IF EXISTS pa_aig_myisam_shared_key_ai_migrate_683fe99e_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_myisam_shared_key_ai_migrate_683fe99e_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_myisam_shared_key_ai_migrate_683fe99e_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c2 INT NOT NULL AUTO_INCREMENT, KEY k2 (c2)) ENGINE=MyISAM
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_myisam_shared_key_ai_migrate_683fe99e_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `c2` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k2` (`c2`)
| ) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_myisam_shared_key_ai_migrate_683fe99e_from
< ok 0 0
> exec
| CREATE DATABASE pa_aig_myisam_shared_key_ai_migrate_683fe99e_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_myisam_shared_key_ai_migrate_683fe99e_from
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c1 INT NOT NULL AUTO_INCREMENT, c2 INT NOT NULL, KEY k (c1, c2)) ENGINE=MyISAM
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_myisam_shared_key_ai_migrate_683fe99e_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `c1` int NOT NULL AUTO_INCREMENT,
|   `c2` int NOT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`c1`,`c2`)
| ) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_myisam_shared_key_ai_migrate_683fe99e
< ok 0 0
> exec
| CREATE DATABASE pa_aig_myisam_shared_key_ai_migrate_683fe99e DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_myisam_shared_key_ai_migrate_683fe99e
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c1 INT NOT NULL AUTO_INCREMENT, c2 INT NOT NULL, KEY k (c1, c2)) ENGINE=MyISAM
< ok 0 0
> exec
| ALTER TABLE `pa_aig_myisam_shared_key_ai_migrate_683fe99e`.`t` DROP INDEX `k`, DROP COLUMN `c1`
< ok 0 0
> exec
| ALTER TABLE `pa_aig_myisam_shared_key_ai_migrate_683fe99e`.`t` MODIFY COLUMN `c2` int NOT NULL AUTO_INCREMENT, ADD KEY `k2` (`c2`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_myisam_shared_key_ai_migrate_683fe99e.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `c2` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k2` (`c2`)
| ) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_AutoIncGroupingApplyCorrectness/8.0/aig-myisam-shared-key-ai-migrate-survive default#1
> exec
| DROP DATABASE IF EXISTS pa_aig_myisam_shared_key_ai_migrate_survive_1b251011_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_myisam_shared_key_ai_migrate_survive_1b251011_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_myisam_shared_key_ai_migrate_survive_1b251011_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL AUTO_INCREMENT, KEY k2 (b)) ENGINE=MyISAM
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_myisam_shared_key_ai_migrate_survive_1b251011_to.t
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
|   `b` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k2` (`b`)
| ) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_myisam_shared_key_ai_migrate_survive_1b251011_from
< ok 0 0
> exec
| CREATE DATABASE pa_aig_myisam_shared_key_ai_migrate_survive_1b251011_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_myisam_shared_key_ai_migrate_survive_1b251011_from
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT NOT NULL AUTO_INCREMENT, b INT NOT NULL, KEY k (a, b)) ENGINE=MyISAM
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_myisam_shared_key_ai_migrate_survive_1b251011_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int NOT NULL AUTO_INCREMENT,
|   `b` int NOT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`a`,`b`)
| ) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_myisam_shared_key_ai_migrate_survive_1b251011
< ok 0 0
> exec
| CREATE DATABASE pa_aig_myisam_shared_key_ai_migrate_survive_1b251011 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_myisam_shared_key_ai_migrate_survive_1b251011
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT NOT NULL AUTO_INCREMENT, b INT NOT NULL, KEY k (a, b)) ENGINE=MyISAM
< ok 0 0
> exec
| ALTER TABLE `pa_aig_myisam_shared_key_ai_migrate_survive_1b251011`.`t` MODIFY COLUMN `a` int NOT NULL, DROP INDEX `k`, MODIFY COLUMN `b` int NOT NULL AUTO_INCREMENT, ADD KEY `k2` (`b`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_myisam_shared_key_ai_migrate_survive_1b251011.t
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
|   `b` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k2` (`b`)
| ) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_AutoIncGroupingApplyCorrectness/8.0/aig-pk-drop-member-drop-newkey default#1
> exec
| DROP DATABASE IF EXISTS pa_aig_pk_drop_member_drop_newkey_2f8c1b05_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_pk_drop_member_drop_newkey_2f8c1b05_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_pk_drop_member_drop_newkey_2f8c1b05_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, KEY k_id (id))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_pk_drop_member_drop_newkey_2f8c1b05_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   KEY `k_id` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_pk_drop_member_drop_newkey_2f8c1b05_from
< ok 0 0
> exec
| CREATE DATABASE pa_aig_pk_drop_member_drop_newkey_2f8c1b05_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_pk_drop_member_drop_newkey_2f8c1b05_from
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, PRIMARY KEY (id, a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_pk_drop_member_drop_newkey_2f8c1b05_from.t
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
| DROP DATABASE IF EXISTS pa_aig_pk_drop_member_drop_newkey_2f8c1b05
< ok 0 0
> exec
| CREATE DATABASE pa_aig_pk_drop_member_drop_newkey_2f8c1b05 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_pk_drop_member_drop_newkey_2f8c1b05
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, PRIMARY KEY (id, a))
< ok 0 0
> exec
| ALTER TABLE `pa_aig_pk_drop_member_drop_newkey_2f8c1b05`.`t` DROP COLUMN `a`
< ok 0 0
> exec
| ALTER TABLE `pa_aig_pk_drop_member_drop_newkey_2f8c1b05`.`t` DROP PRIMARY KEY, ADD KEY `k_id` (`id`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_pk_drop_member_drop_newkey_2f8c1b05.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   KEY `k_id` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_AutoIncGroupingApplyCorrectness/8.0/aig-pk-keep-ai-demote-member default#1
> exec
| DROP DATABASE IF EXISTS pa_aig_pk_keep_ai_demote_member_1a517d31_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_pk_keep_ai_demote_member_1a517d31_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_pk_keep_ai_demote_member_1a517d31_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, a INT NULL, PRIMARY KEY (id))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_pk_keep_ai_demote_member_1a517d31_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `a` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_pk_keep_ai_demote_member_1a517d31_from
< ok 0 0
> exec
| CREATE DATABASE pa_aig_pk_keep_ai_demote_member_1a517d31_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_pk_keep_ai_demote_member_1a517d31_from
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, PRIMARY KEY (id, a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_pk_keep_ai_demote_member_1a517d31_from.t
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
| DROP DATABASE IF EXISTS pa_aig_pk_keep_ai_demote_member_1a517d31
< ok 0 0
> exec
| CREATE DATABASE pa_aig_pk_keep_ai_demote_member_1a517d31 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_pk_keep_ai_demote_member_1a517d31
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, PRIMARY KEY (id, a))
< ok 0 0
> exec
| ALTER TABLE `pa_aig_pk_keep_ai_demote_member_1a517d31`.`t` DROP PRIMARY KEY, ADD PRIMARY KEY (`id`), MODIFY COLUMN `a` int DEFAULT NULL
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_pk_keep_ai_demote_member_1a517d31.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `a` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_AutoIncGroupingApplyCorrectness/8.0/aig-pk-relocate-ai-to-secondary default#1
> exec
| DROP DATABASE IF EXISTS pa_aig_pk_relocate_ai_to_secondary_40b74550_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_pk_relocate_ai_to_secondary_40b74550_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_pk_relocate_ai_to_secondary_40b74550_to
< ok 0 0
> exec
| CREATE TABLE t (c INT NOT NULL AUTO_INCREMENT, b INT NOT NULL, x INT, PRIMARY KEY (b), KEY zz_c (c))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_pk_relocate_ai_to_secondary_40b74550_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `c` int NOT NULL AUTO_INCREMENT,
|   `b` int NOT NULL,
|   `x` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `zz_c` (`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_pk_relocate_ai_to_secondary_40b74550_from
< ok 0 0
> exec
| CREATE DATABASE pa_aig_pk_relocate_ai_to_secondary_40b74550_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_pk_relocate_ai_to_secondary_40b74550_from
< ok 0 0
> exec
| CREATE TABLE t (c INT NOT NULL AUTO_INCREMENT, b INT NOT NULL, x INT, PRIMARY KEY (c))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_pk_relocate_ai_to_secondary_40b74550_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `c` int NOT NULL AUTO_INCREMENT,
|   `b` int NOT NULL,
|   `x` int DEFAULT NULL,
|   PRIMARY KEY (`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_pk_relocate_ai_to_secondary_40b74550
< ok 0 0
> exec
| CREATE DATABASE pa_aig_pk_relocate_ai_to_secondary_40b74550 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_pk_relocate_ai_to_secondary_40b74550
< ok 0 0
> exec
| CREATE TABLE t (c INT NOT NULL AUTO_INCREMENT, b INT NOT NULL, x INT, PRIMARY KEY (c))
< ok 0 0
> exec
| ALTER TABLE `pa_aig_pk_relocate_ai_to_secondary_40b74550`.`t` DROP PRIMARY KEY, ADD KEY `zz_c` (`c`), ADD PRIMARY KEY (`b`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_pk_relocate_ai_to_secondary_40b74550.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `c` int NOT NULL AUTO_INCREMENT,
|   `b` int NOT NULL,
|   `x` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `zz_c` (`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_AutoIncGroupingApplyCorrectness/8.0/aig-pk-shrink-member-drop default#1
> exec
| DROP DATABASE IF EXISTS pa_aig_pk_shrink_member_drop_87079b29_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_pk_shrink_member_drop_87079b29_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_pk_shrink_member_drop_87079b29_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, PRIMARY KEY (id))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_pk_shrink_member_drop_87079b29_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_pk_shrink_member_drop_87079b29_from
< ok 0 0
> exec
| CREATE DATABASE pa_aig_pk_shrink_member_drop_87079b29_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_pk_shrink_member_drop_87079b29_from
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, PRIMARY KEY (id, a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_pk_shrink_member_drop_87079b29_from.t
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
| DROP DATABASE IF EXISTS pa_aig_pk_shrink_member_drop_87079b29
< ok 0 0
> exec
| CREATE DATABASE pa_aig_pk_shrink_member_drop_87079b29 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_pk_shrink_member_drop_87079b29
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, PRIMARY KEY (id, a))
< ok 0 0
> exec
| ALTER TABLE `pa_aig_pk_shrink_member_drop_87079b29`.`t` DROP COLUMN `a`
< ok 0 0
> exec
| ALTER TABLE `pa_aig_pk_shrink_member_drop_87079b29`.`t` DROP PRIMARY KEY, ADD PRIMARY KEY (`id`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_pk_shrink_member_drop_87079b29.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_AutoIncGroupingApplyCorrectness/8.0/aig-pk-swap-drop-secondary default#1
> exec
| DROP DATABASE IF EXISTS pa_aig_pk_swap_drop_secondary_486804bf_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_pk_swap_drop_secondary_486804bf_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_pk_swap_drop_secondary_486804bf_to
< ok 0 0
> exec
| CREATE TABLE t (x INT NOT NULL, c INT NOT NULL AUTO_INCREMENT, PRIMARY KEY (c))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_pk_swap_drop_secondary_486804bf_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `x` int NOT NULL,
|   `c` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_pk_swap_drop_secondary_486804bf_from
< ok 0 0
> exec
| CREATE DATABASE pa_aig_pk_swap_drop_secondary_486804bf_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_pk_swap_drop_secondary_486804bf_from
< ok 0 0
> exec
| CREATE TABLE t (x INT NOT NULL, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_c (c), PRIMARY KEY (x))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_pk_swap_drop_secondary_486804bf_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `x` int NOT NULL,
|   `c` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`x`),
|   UNIQUE KEY `uk_c` (`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_pk_swap_drop_secondary_486804bf
< ok 0 0
> exec
| CREATE DATABASE pa_aig_pk_swap_drop_secondary_486804bf DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_pk_swap_drop_secondary_486804bf
< ok 0 0
> exec
| CREATE TABLE t (x INT NOT NULL, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_c (c), PRIMARY KEY (x))
< ok 0 0
> exec
| ALTER TABLE `pa_aig_pk_swap_drop_secondary_486804bf`.`t` DROP PRIMARY KEY, ADD PRIMARY KEY (`c`), DROP INDEX `uk_c`
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_pk_swap_drop_secondary_486804bf.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `x` int NOT NULL,
|   `c` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_AutoIncGroupingApplyCorrectness/8.0/aig-prefix-over-widened default#1
> exec
| DROP DATABASE IF EXISTS pa_aig_prefix_over_widened_475331f_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_prefix_over_widened_475331f_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_prefix_over_widened_475331f_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, v VARCHAR(500), seq INT NOT NULL AUTO_INCREMENT, KEY k_sv (seq, v(300)))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_prefix_over_widened_475331f_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `v` varchar(500) DEFAULT NULL,
|   `seq` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k_sv` (`seq`,`v`(300))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_prefix_over_widened_475331f_from
< ok 0 0
> exec
| CREATE DATABASE pa_aig_prefix_over_widened_475331f_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_prefix_over_widened_475331f_from
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, v VARCHAR(100))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_prefix_over_widened_475331f_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `v` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_prefix_over_widened_475331f
< ok 0 0
> exec
| CREATE DATABASE pa_aig_prefix_over_widened_475331f DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_prefix_over_widened_475331f
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, v VARCHAR(100))
< ok 0 0
> exec
| ALTER TABLE `pa_aig_prefix_over_widened_475331f`.`t` ADD COLUMN `seq` int NOT NULL AUTO_INCREMENT, MODIFY COLUMN `v` varchar(500) DEFAULT NULL, ADD KEY `k_sv` (`seq`,`v`(300))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_prefix_over_widened_475331f.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `v` varchar(500) DEFAULT NULL,
|   `seq` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k_sv` (`seq`,`v`(300))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_AutoIncGroupingApplyCorrectness/8.0/aig-replace-backing-key default#1
> exec
| DROP DATABASE IF EXISTS pa_aig_replace_backing_key_c75fafdf_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_replace_backing_key_c75fafdf_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_replace_backing_key_c75fafdf_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, UNIQUE KEY uk_c (c, a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_replace_backing_key_c75fafdf_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `c` int NOT NULL AUTO_INCREMENT,
|   `a` int NOT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_c` (`c`,`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_replace_backing_key_c75fafdf_from
< ok 0 0
> exec
| CREATE DATABASE pa_aig_replace_backing_key_c75fafdf_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_replace_backing_key_c75fafdf_from
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, UNIQUE KEY uk_c (c))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_replace_backing_key_c75fafdf_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `c` int NOT NULL AUTO_INCREMENT,
|   `a` int NOT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_c` (`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_replace_backing_key_c75fafdf
< ok 0 0
> exec
| CREATE DATABASE pa_aig_replace_backing_key_c75fafdf DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_replace_backing_key_c75fafdf
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, UNIQUE KEY uk_c (c))
< ok 0 0
> exec
| ALTER TABLE `pa_aig_replace_backing_key_c75fafdf`.`t` DROP INDEX `uk_c`, ADD UNIQUE KEY `uk_c` (`c`,`a`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_replace_backing_key_c75fafdf.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `c` int NOT NULL AUTO_INCREMENT,
|   `a` int NOT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_c` (`c`,`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_AutoIncGroupingApplyCorrectness/8.0/aig-swap-backing-key default#1
> exec
| DROP DATABASE IF EXISTS pa_aig_swap_backing_key_8a8be4d2_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_swap_backing_key_8a8be4d2_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_swap_backing_key_8a8be4d2_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, KEY k_c (c))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_swap_backing_key_8a8be4d2_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `c` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k_c` (`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_swap_backing_key_8a8be4d2_from
< ok 0 0
> exec
| CREATE DATABASE pa_aig_swap_backing_key_8a8be4d2_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_swap_backing_key_8a8be4d2_from
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_c (c))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_swap_backing_key_8a8be4d2_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `c` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_c` (`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_swap_backing_key_8a8be4d2
< ok 0 0
> exec
| CREATE DATABASE pa_aig_swap_backing_key_8a8be4d2 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_swap_backing_key_8a8be4d2
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_c (c))
< ok 0 0
> exec
| ALTER TABLE `pa_aig_swap_backing_key_8a8be4d2`.`t` DROP INDEX `uk_c`, ADD KEY `k_c` (`c`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_swap_backing_key_8a8be4d2.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `c` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k_c` (`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_AutoIncGroupingApplyCorrectness/8.0/aig-widen-ai-key-replaced default#1
> exec
| DROP DATABASE IF EXISTS pa_aig_widen_ai_key_replaced_81695c70_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_widen_ai_key_replaced_81695c70_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_widen_ai_key_replaced_81695c70_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c BIGINT NOT NULL AUTO_INCREMENT, a INT NOT NULL, UNIQUE KEY uk_c (c, a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_widen_ai_key_replaced_81695c70_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `c` bigint NOT NULL AUTO_INCREMENT,
|   `a` int NOT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_c` (`c`,`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_widen_ai_key_replaced_81695c70_from
< ok 0 0
> exec
| CREATE DATABASE pa_aig_widen_ai_key_replaced_81695c70_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_widen_ai_key_replaced_81695c70_from
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, UNIQUE KEY uk_c (c))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_widen_ai_key_replaced_81695c70_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `c` int NOT NULL AUTO_INCREMENT,
|   `a` int NOT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_c` (`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_aig_widen_ai_key_replaced_81695c70
< ok 0 0
> exec
| CREATE DATABASE pa_aig_widen_ai_key_replaced_81695c70 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_widen_ai_key_replaced_81695c70
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, UNIQUE KEY uk_c (c))
< ok 0 0
> exec
| ALTER TABLE `pa_aig_widen_ai_key_replaced_81695c70`.`t` MODIFY COLUMN `c` bigint NOT NULL AUTO_INCREMENT
< ok 0 0
> exec
| ALTER TABLE `pa_aig_widen_ai_key_replaced_81695c70`.`t` DROP INDEX `uk_c`, ADD UNIQUE KEY `uk_c` (`c`,`a`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_widen_ai_key_replaced_81695c70.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `c` bigint NOT NULL AUTO_INCREMENT,
|   `a` int NOT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_c` (`c`,`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
