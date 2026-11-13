-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracle_MigrationApplyCorrectness default#1

== TestOracle_MigrationApplyCorrectness/8.0/add-column default#1
> exec
| DROP DATABASE IF EXISTS pa_add_column_eddb4475_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_column_eddb4475_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_column_eddb4475_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(20))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_column_eddb4475_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_add_column_eddb4475_from
< ok 0 0
> exec
| CREATE DATABASE pa_add_column_eddb4475_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_column_eddb4475_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_column_eddb4475_from.t
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
| DROP DATABASE IF EXISTS pa_add_column_eddb4475
< ok 0 0
> exec
| CREATE DATABASE pa_add_column_eddb4475 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_column_eddb4475
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE `pa_add_column_eddb4475`.`t` ADD COLUMN `name` varchar(20) DEFAULT NULL
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_column_eddb4475.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationApplyCorrectness/8.0/add-column-binary-modifier default#1
> exec
| DROP DATABASE IF EXISTS pa_add_column_binary_modifier_9598d4d9_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_column_binary_modifier_9598d4d9_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_column_binary_modifier_9598d4d9_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, k VARCHAR(128) BINARY NOT NULL) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_column_binary_modifier_9598d4d9_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `k` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_add_column_binary_modifier_9598d4d9_from
< ok 0 0
> exec
| CREATE DATABASE pa_add_column_binary_modifier_9598d4d9_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_column_binary_modifier_9598d4d9_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_column_binary_modifier_9598d4d9_from.t
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
| DROP DATABASE IF EXISTS pa_add_column_binary_modifier_9598d4d9
< ok 0 0
> exec
| CREATE DATABASE pa_add_column_binary_modifier_9598d4d9 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_column_binary_modifier_9598d4d9
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_add_column_binary_modifier_9598d4d9`.`t` ADD COLUMN `k` varchar(128) COLLATE utf8mb4_bin NOT NULL
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_column_binary_modifier_9598d4d9.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `k` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationApplyCorrectness/8.0/add-column-notnull-default default#1
> exec
| DROP DATABASE IF EXISTS pa_add_column_notnull_default_e6fd9208_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_column_notnull_default_e6fd9208_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_column_notnull_default_e6fd9208_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, n INT NOT NULL DEFAULT 5)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_column_notnull_default_e6fd9208_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `n` int NOT NULL DEFAULT '5',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_add_column_notnull_default_e6fd9208_from
< ok 0 0
> exec
| CREATE DATABASE pa_add_column_notnull_default_e6fd9208_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_column_notnull_default_e6fd9208_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_column_notnull_default_e6fd9208_from.t
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
| DROP DATABASE IF EXISTS pa_add_column_notnull_default_e6fd9208
< ok 0 0
> exec
| CREATE DATABASE pa_add_column_notnull_default_e6fd9208 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_column_notnull_default_e6fd9208
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE `pa_add_column_notnull_default_e6fd9208`.`t` ADD COLUMN `n` int NOT NULL DEFAULT '5'
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_column_notnull_default_e6fd9208.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `n` int NOT NULL DEFAULT '5',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationApplyCorrectness/8.0/add-generated-column default#1
> exec
| DROP DATABASE IF EXISTS pa_add_generated_column_99c927b9_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_generated_column_99c927b9_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_generated_column_99c927b9_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT GENERATED ALWAYS AS (a+1) STORED)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_generated_column_99c927b9_to.t
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
|   `b` int GENERATED ALWAYS AS ((`a` + 1)) STORED,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_add_generated_column_99c927b9_from
< ok 0 0
> exec
| CREATE DATABASE pa_add_generated_column_99c927b9_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_generated_column_99c927b9_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_generated_column_99c927b9_from.t
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
| DROP DATABASE IF EXISTS pa_add_generated_column_99c927b9
< ok 0 0
> exec
| CREATE DATABASE pa_add_generated_column_99c927b9 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_generated_column_99c927b9
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> exec
| ALTER TABLE `pa_add_generated_column_99c927b9`.`t` ADD COLUMN `b` int GENERATED ALWAYS AS (((`a` + 1))) STORED
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_generated_column_99c927b9.t
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
|   `b` int GENERATED ALWAYS AS ((`a` + 1)) STORED,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationApplyCorrectness/8.0/add-invisible-column default#1
> exec
| DROP DATABASE IF EXISTS pa_add_invisible_column_ae82ece3_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_invisible_column_ae82ece3_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_invisible_column_ae82ece3_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, secret INT INVISIBLE)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_invisible_column_ae82ece3_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `secret` int DEFAULT NULL /*!80023 INVISIBLE */,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_add_invisible_column_ae82ece3_from
< ok 0 0
> exec
| CREATE DATABASE pa_add_invisible_column_ae82ece3_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_invisible_column_ae82ece3_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_invisible_column_ae82ece3_from.t
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
| DROP DATABASE IF EXISTS pa_add_invisible_column_ae82ece3
< ok 0 0
> exec
| CREATE DATABASE pa_add_invisible_column_ae82ece3 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_invisible_column_ae82ece3
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE `pa_add_invisible_column_ae82ece3`.`t` ADD COLUMN `secret` int DEFAULT NULL /*!80023 INVISIBLE */
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_invisible_column_ae82ece3.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `secret` int DEFAULT NULL /*!80023 INVISIBLE */,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationApplyCorrectness/8.0/add-timestamp-column default#1
> exec
| DROP DATABASE IF EXISTS pa_add_timestamp_column_86ad2852_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_timestamp_column_86ad2852_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_timestamp_column_86ad2852_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, ts TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_timestamp_column_86ad2852_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `ts` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_add_timestamp_column_86ad2852_from
< ok 0 0
> exec
| CREATE DATABASE pa_add_timestamp_column_86ad2852_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_timestamp_column_86ad2852_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_timestamp_column_86ad2852_from.t
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
| DROP DATABASE IF EXISTS pa_add_timestamp_column_86ad2852
< ok 0 0
> exec
| CREATE DATABASE pa_add_timestamp_column_86ad2852 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_timestamp_column_86ad2852
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE `pa_add_timestamp_column_86ad2852`.`t` ADD COLUMN `ts` timestamp NULL DEFAULT CURRENT_TIMESTAMP
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_timestamp_column_86ad2852.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `ts` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationApplyCorrectness/8.0/alter-charset default#1
> exec
| DROP DATABASE IF EXISTS pa_alter_charset_80df01ca_to
< ok 0 0
> exec
| CREATE DATABASE pa_alter_charset_80df01ca_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_alter_charset_80df01ca_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_alter_charset_80df01ca_to.t
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
| DROP DATABASE IF EXISTS pa_alter_charset_80df01ca_from
< ok 0 0
> exec
| CREATE DATABASE pa_alter_charset_80df01ca_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_alter_charset_80df01ca_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) DEFAULT CHARSET=latin1
< ok 0 0
> query
| SHOW CREATE TABLE pa_alter_charset_80df01ca_from.t
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
| ) ENGINE=InnoDB DEFAULT CHARSET=latin1
> exec
| DROP DATABASE IF EXISTS pa_alter_charset_80df01ca
< ok 0 0
> exec
| CREATE DATABASE pa_alter_charset_80df01ca DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_alter_charset_80df01ca
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) DEFAULT CHARSET=latin1
< ok 0 0
> exec
| ALTER TABLE `pa_alter_charset_80df01ca`.`t` DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_alter_charset_80df01ca.t
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

== TestOracle_MigrationApplyCorrectness/8.0/alter-comment default#1
> exec
| DROP DATABASE IF EXISTS pa_alter_comment_6372d06b_to
< ok 0 0
> exec
| CREATE DATABASE pa_alter_comment_6372d06b_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_alter_comment_6372d06b_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) COMMENT='new'
< ok 0 0
> query
| SHOW CREATE TABLE pa_alter_comment_6372d06b_to.t
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='new'
> exec
| DROP DATABASE IF EXISTS pa_alter_comment_6372d06b_from
< ok 0 0
> exec
| CREATE DATABASE pa_alter_comment_6372d06b_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_alter_comment_6372d06b_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) COMMENT='old'
< ok 0 0
> query
| SHOW CREATE TABLE pa_alter_comment_6372d06b_from.t
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='old'
> exec
| DROP DATABASE IF EXISTS pa_alter_comment_6372d06b
< ok 0 0
> exec
| CREATE DATABASE pa_alter_comment_6372d06b DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_alter_comment_6372d06b
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) COMMENT='old'
< ok 0 0
> exec
| ALTER TABLE `pa_alter_comment_6372d06b`.`t` COMMENT='new'
< ok 0 0
> query
| SHOW CREATE TABLE pa_alter_comment_6372d06b.t
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='new'

== TestOracle_MigrationApplyCorrectness/8.0/alter-engine default#1
> exec
| DROP DATABASE IF EXISTS pa_alter_engine_9436e474_to
< ok 0 0
> exec
| CREATE DATABASE pa_alter_engine_9436e474_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_alter_engine_9436e474_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) ENGINE=MyISAM
< ok 0 0
> query
| SHOW CREATE TABLE pa_alter_engine_9436e474_to.t
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
| ) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_alter_engine_9436e474_from
< ok 0 0
> exec
| CREATE DATABASE pa_alter_engine_9436e474_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_alter_engine_9436e474_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) ENGINE=InnoDB
< ok 0 0
> query
| SHOW CREATE TABLE pa_alter_engine_9436e474_from.t
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
| DROP DATABASE IF EXISTS pa_alter_engine_9436e474
< ok 0 0
> exec
| CREATE DATABASE pa_alter_engine_9436e474 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_alter_engine_9436e474
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) ENGINE=InnoDB
< ok 0 0
> exec
| ALTER TABLE `pa_alter_engine_9436e474`.`t` ENGINE=MyISAM
< ok 0 0
> query
| SHOW CREATE TABLE pa_alter_engine_9436e474.t
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
| ) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationApplyCorrectness/8.0/alter-rowformat-add default#1
> exec
| DROP DATABASE IF EXISTS pa_alter_rowformat_add_59de78e7_to
< ok 0 0
> exec
| CREATE DATABASE pa_alter_rowformat_add_59de78e7_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_alter_rowformat_add_59de78e7_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=8
< ok 0 0
> query
| SHOW CREATE TABLE pa_alter_rowformat_add_59de78e7_to.t
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=8
> exec
| DROP DATABASE IF EXISTS pa_alter_rowformat_add_59de78e7_from
< ok 0 0
> exec
| CREATE DATABASE pa_alter_rowformat_add_59de78e7_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_alter_rowformat_add_59de78e7_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_alter_rowformat_add_59de78e7_from.t
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
| DROP DATABASE IF EXISTS pa_alter_rowformat_add_59de78e7
< ok 0 0
> exec
| CREATE DATABASE pa_alter_rowformat_add_59de78e7 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_alter_rowformat_add_59de78e7
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE `pa_alter_rowformat_add_59de78e7`.`t` ROW_FORMAT=COMPRESSED
< ok 0 0
> query
| SHOW CREATE TABLE pa_alter_rowformat_add_59de78e7.t
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=COMPRESSED

== TestOracle_MigrationApplyCorrectness/8.0/alter-rowformat-change default#1
> exec
| DROP DATABASE IF EXISTS pa_alter_rowformat_change_8983aaac_to
< ok 0 0
> exec
| CREATE DATABASE pa_alter_rowformat_change_8983aaac_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_alter_rowformat_change_8983aaac_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) ROW_FORMAT=DYNAMIC
< ok 0 0
> query
| SHOW CREATE TABLE pa_alter_rowformat_change_8983aaac_to.t
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC
> exec
| DROP DATABASE IF EXISTS pa_alter_rowformat_change_8983aaac_from
< ok 0 0
> exec
| CREATE DATABASE pa_alter_rowformat_change_8983aaac_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_alter_rowformat_change_8983aaac_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) ROW_FORMAT=COMPACT
< ok 0 0
> query
| SHOW CREATE TABLE pa_alter_rowformat_change_8983aaac_from.t
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=COMPACT
> exec
| DROP DATABASE IF EXISTS pa_alter_rowformat_change_8983aaac
< ok 0 0
> exec
| CREATE DATABASE pa_alter_rowformat_change_8983aaac DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_alter_rowformat_change_8983aaac
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) ROW_FORMAT=COMPACT
< ok 0 0
> exec
| ALTER TABLE `pa_alter_rowformat_change_8983aaac`.`t` ROW_FORMAT=DYNAMIC
< ok 0 0
> query
| SHOW CREATE TABLE pa_alter_rowformat_change_8983aaac.t
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC

== TestOracle_MigrationApplyCorrectness/8.0/create-autoinc default#1
> exec
| DROP DATABASE IF EXISTS pa_create_autoinc_925a4441_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_autoinc_925a4441_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_autoinc_925a4441_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT PRIMARY KEY, a INT) AUTO_INCREMENT=100
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_autoinc_925a4441_to.t
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
| ) ENGINE=InnoDB AUTO_INCREMENT=100 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_autoinc_925a4441
< ok 0 0
> exec
| CREATE DATABASE pa_create_autoinc_925a4441 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_autoinc_925a4441
< ok 0 0
> exec
| CREATE TABLE `pa_create_autoinc_925a4441`.`t` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `a` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_autoinc_925a4441.t
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

== TestOracle_MigrationApplyCorrectness/8.0/create-boolean default#1
> exec
| DROP DATABASE IF EXISTS pa_create_boolean_3b873e88_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_boolean_3b873e88_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_boolean_3b873e88_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, m BOOLEAN, d TINYINT(1))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_boolean_3b873e88_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `m` tinyint(1) DEFAULT NULL,
|   `d` tinyint(1) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_boolean_3b873e88
< ok 0 0
> exec
| CREATE DATABASE pa_create_boolean_3b873e88 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_boolean_3b873e88
< ok 0 0
> exec
| CREATE TABLE `pa_create_boolean_3b873e88`.`t` (
|   `id` int NOT NULL,
|   `m` tinyint(1) DEFAULT NULL,
|   `d` tinyint(1) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_boolean_3b873e88.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `m` tinyint(1) DEFAULT NULL,
|   `d` tinyint(1) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationApplyCorrectness/8.0/create-char-bit-year default#1
> exec
| DROP DATABASE IF EXISTS pa_create_char_bit_year_f76ce9bc_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_char_bit_year_f76ce9bc_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_char_bit_year_f76ce9bc_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a CHAR, d BINARY, c BIT, y YEAR, v VARCHAR(40), bb VARBINARY(8))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_char_bit_year_f76ce9bc_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` char(1) DEFAULT NULL,
|   `d` binary(1) DEFAULT NULL,
|   `c` bit(1) DEFAULT NULL,
|   `y` year DEFAULT NULL,
|   `v` varchar(40) DEFAULT NULL,
|   `bb` varbinary(8) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_char_bit_year_f76ce9bc
< ok 0 0
> exec
| CREATE DATABASE pa_create_char_bit_year_f76ce9bc DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_char_bit_year_f76ce9bc
< ok 0 0
> exec
| CREATE TABLE `pa_create_char_bit_year_f76ce9bc`.`t` (
|   `id` int NOT NULL,
|   `a` char(1) DEFAULT NULL,
|   `d` binary(1) DEFAULT NULL,
|   `c` bit(1) DEFAULT NULL,
|   `y` year DEFAULT NULL,
|   `v` varchar(40) DEFAULT NULL,
|   `bb` varbinary(8) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_char_bit_year_f76ce9bc.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` char(1) DEFAULT NULL,
|   `d` binary(1) DEFAULT NULL,
|   `c` bit(1) DEFAULT NULL,
|   `y` year DEFAULT NULL,
|   `v` varchar(40) DEFAULT NULL,
|   `bb` varbinary(8) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationApplyCorrectness/8.0/create-charset-explicit default#1
> exec
| DROP DATABASE IF EXISTS pa_create_charset_explicit_a87fd03f_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_charset_explicit_a87fd03f_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_charset_explicit_a87fd03f_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a VARCHAR(10) CHARACTER SET utf8mb4, b VARCHAR(10)) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_charset_explicit_a87fd03f_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
|   `b` varchar(10) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_charset_explicit_a87fd03f
< ok 0 0
> exec
| CREATE DATABASE pa_create_charset_explicit_a87fd03f DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_charset_explicit_a87fd03f
< ok 0 0
> exec
| CREATE TABLE `pa_create_charset_explicit_a87fd03f`.`t` (
|   `id` int NOT NULL,
|   `a` varchar(10) DEFAULT NULL,
|   `b` varchar(10) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_charset_explicit_a87fd03f.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` varchar(10) DEFAULT NULL,
|   `b` varchar(10) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationApplyCorrectness/8.0/create-collate-only-column default#1
> exec
| DROP DATABASE IF EXISTS pa_create_collate_only_column_d3abd75c_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_collate_only_column_d3abd75c_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_collate_only_column_d3abd75c_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a VARCHAR(10) COLLATE utf8mb4_unicode_ci) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_collate_only_column_d3abd75c_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_collate_only_column_d3abd75c
< ok 0 0
> exec
| CREATE DATABASE pa_create_collate_only_column_d3abd75c DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_collate_only_column_d3abd75c
< ok 0 0
> exec
| CREATE TABLE `pa_create_collate_only_column_d3abd75c`.`t` (
|   `id` int NOT NULL,
|   `a` varchar(10) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_collate_only_column_d3abd75c.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationApplyCorrectness/8.0/create-comment default#1
> exec
| DROP DATABASE IF EXISTS pa_create_comment_faa12879_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_comment_faa12879_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_comment_faa12879_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT COMMENT 'has ''quote'' inside') COMMENT='tbl ''c'''
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_comment_faa12879_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL COMMENT 'has ''quote'' inside',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='tbl ''c'''
> exec
| DROP DATABASE IF EXISTS pa_create_comment_faa12879
< ok 0 0
> exec
| CREATE DATABASE pa_create_comment_faa12879 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_comment_faa12879
< ok 0 0
> exec
| CREATE TABLE `pa_create_comment_faa12879`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL COMMENT 'has ''quote'' inside',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='tbl ''c'''
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_comment_faa12879.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL COMMENT 'has ''quote'' inside',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='tbl ''c'''

== TestOracle_MigrationApplyCorrectness/8.0/create-datetime-default default#1
> exec
| DROP DATABASE IF EXISTS pa_create_datetime_default_2aa5b909_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_datetime_default_2aa5b909_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_datetime_default_2aa5b909_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g DATETIME DEFAULT CURRENT_TIMESTAMP, h DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP)
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_datetime_default_2aa5b909_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `g` datetime DEFAULT CURRENT_TIMESTAMP,
|   `h` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_datetime_default_2aa5b909
< ok 0 0
> exec
| CREATE DATABASE pa_create_datetime_default_2aa5b909 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_datetime_default_2aa5b909
< ok 0 0
> exec
| CREATE TABLE `pa_create_datetime_default_2aa5b909`.`t` (
|   `id` int NOT NULL,
|   `g` datetime DEFAULT CURRENT_TIMESTAMP,
|   `h` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_datetime_default_2aa5b909.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `g` datetime DEFAULT CURRENT_TIMESTAMP,
|   `h` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationApplyCorrectness/8.0/create-defaults default#1
> exec
| DROP DATABASE IF EXISTS pa_create_defaults_1b01d526_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_defaults_1b01d526_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_defaults_1b01d526_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT DEFAULT 0, b INT DEFAULT '7', f DECIMAL(10,2) DEFAULT 0, c VARCHAR(20) DEFAULT 'x', j BOOLEAN DEFAULT TRUE)
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_defaults_1b01d526_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT '0',
|   `b` int DEFAULT '7',
|   `f` decimal(10,2) DEFAULT '0.00',
|   `c` varchar(20) DEFAULT 'x',
|   `j` tinyint(1) DEFAULT '1',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_defaults_1b01d526
< ok 0 0
> exec
| CREATE DATABASE pa_create_defaults_1b01d526 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_defaults_1b01d526
< ok 0 0
> exec
| CREATE TABLE `pa_create_defaults_1b01d526`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT '0',
|   `b` int DEFAULT '7',
|   `f` decimal(10,2) DEFAULT '0.00',
|   `c` varchar(20) DEFAULT 'x',
|   `j` tinyint(1) DEFAULT '1',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_defaults_1b01d526.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT '0',
|   `b` int DEFAULT '7',
|   `f` decimal(10,2) DEFAULT '0.00',
|   `c` varchar(20) DEFAULT 'x',
|   `j` tinyint(1) DEFAULT '1',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationApplyCorrectness/8.0/create-engine-myisam default#1
> exec
| DROP DATABASE IF EXISTS pa_create_engine_myisam_26fea2f9_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_engine_myisam_26fea2f9_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_engine_myisam_26fea2f9_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT) ENGINE=MyISAM
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_engine_myisam_26fea2f9_to.t
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
| ) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_engine_myisam_26fea2f9
< ok 0 0
> exec
| CREATE DATABASE pa_create_engine_myisam_26fea2f9 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_engine_myisam_26fea2f9
< ok 0 0
> exec
| CREATE TABLE `pa_create_engine_myisam_26fea2f9`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_engine_myisam_26fea2f9.t
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
| ) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationApplyCorrectness/8.0/create-enum-set default#1
> exec
| DROP DATABASE IF EXISTS pa_create_enum_set_6dab39d4_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_enum_set_6dab39d4_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_enum_set_6dab39d4_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a ENUM('x','y','z'), b SET('a','b','c'))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_enum_set_6dab39d4_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` enum('x','y','z') DEFAULT NULL,
|   `b` set('a','b','c') DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_enum_set_6dab39d4
< ok 0 0
> exec
| CREATE DATABASE pa_create_enum_set_6dab39d4 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_enum_set_6dab39d4
< ok 0 0
> exec
| CREATE TABLE `pa_create_enum_set_6dab39d4`.`t` (
|   `id` int NOT NULL,
|   `a` enum('x','y','z') DEFAULT NULL,
|   `b` set('a','b','c') DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_enum_set_6dab39d4.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` enum('x','y','z') DEFAULT NULL,
|   `b` set('a','b','c') DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationApplyCorrectness/8.0/create-functional-default default#1
> exec
| DROP DATABASE IF EXISTS pa_create_functional_default_24d1fb21_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_functional_default_24d1fb21_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_functional_default_24d1fb21_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT DEFAULT (1+1), b VARCHAR(36) DEFAULT (UUID()))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_functional_default_24d1fb21_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT ((1 + 1)),
|   `b` varchar(36) DEFAULT (uuid()),
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_functional_default_24d1fb21
< ok 0 0
> exec
| CREATE DATABASE pa_create_functional_default_24d1fb21 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_functional_default_24d1fb21
< ok 0 0
> exec
| CREATE TABLE `pa_create_functional_default_24d1fb21`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT ((1 + 1)),
|   `b` varchar(36) DEFAULT (uuid()),
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_functional_default_24d1fb21.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT ((1 + 1)),
|   `b` varchar(36) DEFAULT (uuid()),
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationApplyCorrectness/8.0/create-generated default#1
> exec
| DROP DATABASE IF EXISTS pa_create_generated_3c18219d_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_generated_3c18219d_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_generated_3c18219d_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT GENERATED ALWAYS AS (a+1) VIRTUAL, c INT GENERATED ALWAYS AS (a*2) STORED)
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_generated_3c18219d_to.t
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
|   `b` int GENERATED ALWAYS AS ((`a` + 1)) VIRTUAL,
|   `c` int GENERATED ALWAYS AS ((`a` * 2)) STORED,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_generated_3c18219d
< ok 0 0
> exec
| CREATE DATABASE pa_create_generated_3c18219d DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_generated_3c18219d
< ok 0 0
> exec
| CREATE TABLE `pa_create_generated_3c18219d`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   `b` int GENERATED ALWAYS AS (((`a` + 1))) VIRTUAL,
|   `c` int GENERATED ALWAYS AS (((`a` * 2))) STORED,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_generated_3c18219d.t
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
|   `b` int GENERATED ALWAYS AS ((`a` + 1)) VIRTUAL,
|   `c` int GENERATED ALWAYS AS ((`a` * 2)) STORED,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationApplyCorrectness/8.0/create-int-widths default#1
> exec
| DROP DATABASE IF EXISTS pa_create_int_widths_6c756a45_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_int_widths_6c756a45_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_int_widths_6c756a45_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT(11), b BIGINT(20), p TINYINT, h INT UNSIGNED, z INT ZEROFILL)
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_int_widths_6c756a45_to.t
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
|   `b` bigint DEFAULT NULL,
|   `p` tinyint DEFAULT NULL,
|   `h` int unsigned DEFAULT NULL,
|   `z` int(10) unsigned zerofill DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_int_widths_6c756a45
< ok 0 0
> exec
| CREATE DATABASE pa_create_int_widths_6c756a45 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_int_widths_6c756a45
< ok 0 0
> exec
| CREATE TABLE `pa_create_int_widths_6c756a45`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   `b` bigint DEFAULT NULL,
|   `p` tinyint DEFAULT NULL,
|   `h` int unsigned DEFAULT NULL,
|   `z` int(10) unsigned zerofill DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_int_widths_6c756a45.t
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
|   `b` bigint DEFAULT NULL,
|   `p` tinyint DEFAULT NULL,
|   `h` int unsigned DEFAULT NULL,
|   `z` int(10) unsigned zerofill DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationApplyCorrectness/8.0/create-invisible-column default#1
> exec
| DROP DATABASE IF EXISTS pa_create_invisible_column_c34f34c4_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_invisible_column_c34f34c4_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_invisible_column_c34f34c4_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, secret INT INVISIBLE)
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_invisible_column_c34f34c4_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `secret` int DEFAULT NULL /*!80023 INVISIBLE */,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_invisible_column_c34f34c4
< ok 0 0
> exec
| CREATE DATABASE pa_create_invisible_column_c34f34c4 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_invisible_column_c34f34c4
< ok 0 0
> exec
| CREATE TABLE `pa_create_invisible_column_c34f34c4`.`t` (
|   `id` int NOT NULL,
|   `secret` int DEFAULT NULL /*!80023 INVISIBLE */,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_invisible_column_c34f34c4.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `secret` int DEFAULT NULL /*!80023 INVISIBLE */,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationApplyCorrectness/8.0/create-multi-column default#1
> exec
| DROP DATABASE IF EXISTS pa_create_multi_column_e3de953e_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_multi_column_e3de953e_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_multi_column_e3de953e_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, c1 VARCHAR(10), c2 INT DEFAULT 0, c3 DATE, c4 DECIMAL(8,3), PRIMARY KEY (id))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_multi_column_e3de953e_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `c1` varchar(10) DEFAULT NULL,
|   `c2` int DEFAULT '0',
|   `c3` date DEFAULT NULL,
|   `c4` decimal(8,3) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_multi_column_e3de953e
< ok 0 0
> exec
| CREATE DATABASE pa_create_multi_column_e3de953e DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_multi_column_e3de953e
< ok 0 0
> exec
| CREATE TABLE `pa_create_multi_column_e3de953e`.`t` (
|   `id` int NOT NULL,
|   `c1` varchar(10) DEFAULT NULL,
|   `c2` int DEFAULT '0',
|   `c3` date DEFAULT NULL,
|   `c4` decimal(8,3) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_multi_column_e3de953e.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `c1` varchar(10) DEFAULT NULL,
|   `c2` int DEFAULT '0',
|   `c3` date DEFAULT NULL,
|   `c4` decimal(8,3) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationApplyCorrectness/8.0/create-nullability default#1
> exec
| DROP DATABASE IF EXISTS pa_create_nullability_39e33d6d_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_nullability_39e33d6d_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_nullability_39e33d6d_to
< ok 0 0
> exec
| CREATE TABLE t (id INT, a INT NULL, b INT DEFAULT NULL, d VARCHAR(10) NOT NULL, PRIMARY KEY (id))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_nullability_39e33d6d_to.t
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
|   `d` varchar(10) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_nullability_39e33d6d
< ok 0 0
> exec
| CREATE DATABASE pa_create_nullability_39e33d6d DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_nullability_39e33d6d
< ok 0 0
> exec
| CREATE TABLE `pa_create_nullability_39e33d6d`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   `d` varchar(10) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_nullability_39e33d6d.t
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
|   `d` varchar(10) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationApplyCorrectness/8.0/create-numeric default#1
> exec
| DROP DATABASE IF EXISTS pa_create_numeric_8953d01f_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_numeric_8953d01f_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_numeric_8953d01f_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a DECIMAL, b DECIMAL(10,2), j REAL, f FLOAT, g DOUBLE(15,4))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_numeric_8953d01f_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` decimal(10,0) DEFAULT NULL,
|   `b` decimal(10,2) DEFAULT NULL,
|   `j` double DEFAULT NULL,
|   `f` float DEFAULT NULL,
|   `g` double(15,4) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_numeric_8953d01f
< ok 0 0
> exec
| CREATE DATABASE pa_create_numeric_8953d01f DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_numeric_8953d01f
< ok 0 0
> exec
| CREATE TABLE `pa_create_numeric_8953d01f`.`t` (
|   `id` int NOT NULL,
|   `a` decimal(10,0) DEFAULT NULL,
|   `b` decimal(10,2) DEFAULT NULL,
|   `j` double DEFAULT NULL,
|   `f` float DEFAULT NULL,
|   `g` double(15,4) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_numeric_8953d01f.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` decimal(10,0) DEFAULT NULL,
|   `b` decimal(10,2) DEFAULT NULL,
|   `j` double DEFAULT NULL,
|   `f` float DEFAULT NULL,
|   `g` double(15,4) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationApplyCorrectness/8.0/create-row-format default#1
> exec
| DROP DATABASE IF EXISTS pa_create_row_format_54b3878_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_row_format_54b3878_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_row_format_54b3878_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=8
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_row_format_54b3878_to.t
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=8
> exec
| DROP DATABASE IF EXISTS pa_create_row_format_54b3878
< ok 0 0
> exec
| CREATE DATABASE pa_create_row_format_54b3878 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_row_format_54b3878
< ok 0 0
> exec
| CREATE TABLE `pa_create_row_format_54b3878`.`t` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPRESSED
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_row_format_54b3878.t
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=COMPRESSED

== TestOracle_MigrationApplyCorrectness/8.0/create-simple default#1
> exec
| DROP DATABASE IF EXISTS pa_create_simple_7cb4d00e_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_simple_7cb4d00e_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_simple_7cb4d00e_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_simple_7cb4d00e_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_simple_7cb4d00e
< ok 0 0
> exec
| CREATE DATABASE pa_create_simple_7cb4d00e DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_simple_7cb4d00e
< ok 0 0
> exec
| CREATE TABLE `pa_create_simple_7cb4d00e`.`t` (
|   `id` int NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_simple_7cb4d00e.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationApplyCorrectness/8.0/create-timestamp default#1
> exec
| DROP DATABASE IF EXISTS pa_create_timestamp_ab0524fc_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_timestamp_ab0524fc_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_timestamp_ab0524fc_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a TIMESTAMP, b TIMESTAMP NULL, d TIMESTAMP DEFAULT CURRENT_TIMESTAMP)
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_timestamp_ab0524fc_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` timestamp NULL DEFAULT NULL,
|   `b` timestamp NULL DEFAULT NULL,
|   `d` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_timestamp_ab0524fc
< ok 0 0
> exec
| CREATE DATABASE pa_create_timestamp_ab0524fc DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_timestamp_ab0524fc
< ok 0 0
> exec
| CREATE TABLE `pa_create_timestamp_ab0524fc`.`t` (
|   `id` int NOT NULL,
|   `a` timestamp NULL DEFAULT NULL,
|   `b` timestamp NULL DEFAULT NULL,
|   `d` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_timestamp_ab0524fc.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` timestamp NULL DEFAULT NULL,
|   `b` timestamp NULL DEFAULT NULL,
|   `d` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationApplyCorrectness/8.0/drop-column default#1
> exec
| DROP DATABASE IF EXISTS pa_drop_column_513780e5_to
< ok 0 0
> exec
| CREATE DATABASE pa_drop_column_513780e5_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_column_513780e5_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_column_513780e5_to.t
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
| DROP DATABASE IF EXISTS pa_drop_column_513780e5_from
< ok 0 0
> exec
| CREATE DATABASE pa_drop_column_513780e5_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_column_513780e5_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(20))
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_column_513780e5_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_drop_column_513780e5
< ok 0 0
> exec
| CREATE DATABASE pa_drop_column_513780e5 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_column_513780e5
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(20))
< ok 0 0
> exec
| ALTER TABLE `pa_drop_column_513780e5`.`t` DROP COLUMN `name`
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_column_513780e5.t
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

== TestOracle_MigrationApplyCorrectness/8.0/drop-table default#1
> exec
| DROP DATABASE IF EXISTS pa_drop_table_775be169_from
< ok 0 0
> exec
| CREATE DATABASE pa_drop_table_775be169_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_table_775be169_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_table_775be169_from.t
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
| DROP DATABASE IF EXISTS pa_drop_table_775be169
< ok 0 0
> exec
| CREATE DATABASE pa_drop_table_775be169 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_table_775be169
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> exec
| DROP TABLE `pa_drop_table_775be169`.`t`
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_table_775be169.t
< error
| 1146 42S02 Table 'pa_drop_table_775be169.t' doesn't exist

== TestOracle_MigrationApplyCorrectness/8.0/gen-drop-with-base default#1
> exec
| DROP DATABASE IF EXISTS pa_gen_drop_with_base_c17d416_to
< ok 0 0
> exec
| CREATE DATABASE pa_gen_drop_with_base_c17d416_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_gen_drop_with_base_c17d416_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_gen_drop_with_base_c17d416_to.t
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
| DROP DATABASE IF EXISTS pa_gen_drop_with_base_c17d416_from
< ok 0 0
> exec
| CREATE DATABASE pa_gen_drop_with_base_c17d416_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_gen_drop_with_base_c17d416_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, z INT GENERATED ALWAYS AS (a+1) STORED)
< ok 0 0
> query
| SHOW CREATE TABLE pa_gen_drop_with_base_c17d416_from.t
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
|   `z` int GENERATED ALWAYS AS ((`a` + 1)) STORED,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_gen_drop_with_base_c17d416
< ok 0 0
> exec
| CREATE DATABASE pa_gen_drop_with_base_c17d416 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_gen_drop_with_base_c17d416
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, z INT GENERATED ALWAYS AS (a+1) STORED)
< ok 0 0
> exec
| ALTER TABLE `pa_gen_drop_with_base_c17d416`.`t` DROP COLUMN `z`
< ok 0 0
> exec
| ALTER TABLE `pa_gen_drop_with_base_c17d416`.`t` DROP COLUMN `a`
< ok 0 0
> query
| SHOW CREATE TABLE pa_gen_drop_with_base_c17d416.t
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

== TestOracle_MigrationApplyCorrectness/8.0/gen-flip-with-base-modify default#1
> exec
| DROP DATABASE IF EXISTS pa_gen_flip_with_base_modify_ba61a425_to
< ok 0 0
> exec
| CREATE DATABASE pa_gen_flip_with_base_modify_ba61a425_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_gen_flip_with_base_modify_ba61a425_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g INT GENERATED ALWAYS AS (z+1) STORED, z BIGINT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_gen_flip_with_base_modify_ba61a425_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `g` int GENERATED ALWAYS AS ((`z` + 1)) STORED,
|   `z` bigint DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_gen_flip_with_base_modify_ba61a425_from
< ok 0 0
> exec
| CREATE DATABASE pa_gen_flip_with_base_modify_ba61a425_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_gen_flip_with_base_modify_ba61a425_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g INT GENERATED ALWAYS AS (z+1) VIRTUAL, z SMALLINT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_gen_flip_with_base_modify_ba61a425_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `g` int GENERATED ALWAYS AS ((`z` + 1)) VIRTUAL,
|   `z` smallint DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_gen_flip_with_base_modify_ba61a425
< ok 0 0
> exec
| CREATE DATABASE pa_gen_flip_with_base_modify_ba61a425 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_gen_flip_with_base_modify_ba61a425
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g INT GENERATED ALWAYS AS (z+1) VIRTUAL, z SMALLINT)
< ok 0 0
> exec
| ALTER TABLE `pa_gen_flip_with_base_modify_ba61a425`.`t` DROP COLUMN `g`
< ok 0 0
> exec
| ALTER TABLE `pa_gen_flip_with_base_modify_ba61a425`.`t` MODIFY COLUMN `z` bigint DEFAULT NULL
< ok 0 0
> exec
| ALTER TABLE `pa_gen_flip_with_base_modify_ba61a425`.`t` ADD COLUMN `g` int GENERATED ALWAYS AS (((`z` + 1))) STORED
< ok 0 0
> query
| SHOW CREATE TABLE pa_gen_flip_with_base_modify_ba61a425.t
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
|   `g` int GENERATED ALWAYS AS ((`z` + 1)) STORED,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationApplyCorrectness/8.0/gen-storage-flip default#1
> exec
| DROP DATABASE IF EXISTS pa_gen_storage_flip_76fab813_to
< ok 0 0
> exec
| CREATE DATABASE pa_gen_storage_flip_76fab813_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_gen_storage_flip_76fab813_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, g INT GENERATED ALWAYS AS (a+1) STORED)
< ok 0 0
> query
| SHOW CREATE TABLE pa_gen_storage_flip_76fab813_to.t
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
|   `g` int GENERATED ALWAYS AS ((`a` + 1)) STORED,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_gen_storage_flip_76fab813_from
< ok 0 0
> exec
| CREATE DATABASE pa_gen_storage_flip_76fab813_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_gen_storage_flip_76fab813_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, g INT GENERATED ALWAYS AS (a+1) VIRTUAL)
< ok 0 0
> query
| SHOW CREATE TABLE pa_gen_storage_flip_76fab813_from.t
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
|   `g` int GENERATED ALWAYS AS ((`a` + 1)) VIRTUAL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_gen_storage_flip_76fab813
< ok 0 0
> exec
| CREATE DATABASE pa_gen_storage_flip_76fab813 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_gen_storage_flip_76fab813
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, g INT GENERATED ALWAYS AS (a+1) VIRTUAL)
< ok 0 0
> exec
| ALTER TABLE `pa_gen_storage_flip_76fab813`.`t` DROP COLUMN `g`
< ok 0 0
> exec
| ALTER TABLE `pa_gen_storage_flip_76fab813`.`t` ADD COLUMN `g` int GENERATED ALWAYS AS (((`a` + 1))) STORED
< ok 0 0
> query
| SHOW CREATE TABLE pa_gen_storage_flip_76fab813.t
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
|   `g` int GENERATED ALWAYS AS ((`a` + 1)) STORED,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationApplyCorrectness/8.0/modify-add-autoinc default#1
> exec
| DROP DATABASE IF EXISTS pa_modify_add_autoinc_d963cd2d_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_add_autoinc_d963cd2d_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_add_autoinc_d963cd2d_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_add_autoinc_d963cd2d_to.t
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
| DROP DATABASE IF EXISTS pa_modify_add_autoinc_d963cd2d_from
< ok 0 0
> exec
| CREATE DATABASE pa_modify_add_autoinc_d963cd2d_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_add_autoinc_d963cd2d_from
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_add_autoinc_d963cd2d_from.t
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
| DROP DATABASE IF EXISTS pa_modify_add_autoinc_d963cd2d
< ok 0 0
> exec
| CREATE DATABASE pa_modify_add_autoinc_d963cd2d DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_add_autoinc_d963cd2d
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE `pa_modify_add_autoinc_d963cd2d`.`t` MODIFY COLUMN `id` int NOT NULL AUTO_INCREMENT
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_add_autoinc_d963cd2d.t
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

== TestOracle_MigrationApplyCorrectness/8.0/modify-charset default#1
> exec
| DROP DATABASE IF EXISTS pa_modify_charset_c1e54a02_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_charset_c1e54a02_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_charset_c1e54a02_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v VARCHAR(10) CHARACTER SET utf8mb4) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_charset_c1e54a02_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `v` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_charset_c1e54a02_from
< ok 0 0
> exec
| CREATE DATABASE pa_modify_charset_c1e54a02_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_charset_c1e54a02_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v VARCHAR(10) CHARACTER SET latin1) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_charset_c1e54a02_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `v` varchar(10) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_charset_c1e54a02
< ok 0 0
> exec
| CREATE DATABASE pa_modify_charset_c1e54a02 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_charset_c1e54a02
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v VARCHAR(10) CHARACTER SET latin1) DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_modify_charset_c1e54a02`.`t` MODIFY COLUMN `v` varchar(10) DEFAULT NULL
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_charset_c1e54a02.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `v` varchar(10) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationApplyCorrectness/8.0/modify-column-to-binary-modifier default#1
> exec
| DROP DATABASE IF EXISTS pa_modify_column_to_binary_modifier_f642bbce_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_column_to_binary_modifier_f642bbce_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_column_to_binary_modifier_f642bbce_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v VARCHAR(10) BINARY) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_column_to_binary_modifier_f642bbce_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `v` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_column_to_binary_modifier_f642bbce_from
< ok 0 0
> exec
| CREATE DATABASE pa_modify_column_to_binary_modifier_f642bbce_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_column_to_binary_modifier_f642bbce_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v VARCHAR(10)) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_column_to_binary_modifier_f642bbce_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `v` varchar(10) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_column_to_binary_modifier_f642bbce
< ok 0 0
> exec
| CREATE DATABASE pa_modify_column_to_binary_modifier_f642bbce DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_column_to_binary_modifier_f642bbce
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v VARCHAR(10)) DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_modify_column_to_binary_modifier_f642bbce`.`t` MODIFY COLUMN `v` varchar(10) COLLATE utf8mb4_bin DEFAULT NULL
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_column_to_binary_modifier_f642bbce.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `v` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationApplyCorrectness/8.0/modify-comment default#1
> exec
| DROP DATABASE IF EXISTS pa_modify_comment_5a789853_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_comment_5a789853_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_comment_5a789853_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v INT COMMENT 'b')
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_comment_5a789853_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `v` int DEFAULT NULL COMMENT 'b',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_comment_5a789853_from
< ok 0 0
> exec
| CREATE DATABASE pa_modify_comment_5a789853_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_comment_5a789853_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v INT COMMENT 'a')
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_comment_5a789853_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `v` int DEFAULT NULL COMMENT 'a',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_comment_5a789853
< ok 0 0
> exec
| CREATE DATABASE pa_modify_comment_5a789853 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_comment_5a789853
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v INT COMMENT 'a')
< ok 0 0
> exec
| ALTER TABLE `pa_modify_comment_5a789853`.`t` MODIFY COLUMN `v` int DEFAULT NULL COMMENT 'b'
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_comment_5a789853.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `v` int DEFAULT NULL COMMENT 'b',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationApplyCorrectness/8.0/modify-default default#1
> exec
| DROP DATABASE IF EXISTS pa_modify_default_a02318c3_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_default_a02318c3_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_default_a02318c3_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v INT DEFAULT 1)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_default_a02318c3_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `v` int DEFAULT '1',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_default_a02318c3_from
< ok 0 0
> exec
| CREATE DATABASE pa_modify_default_a02318c3_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_default_a02318c3_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v INT DEFAULT 0)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_default_a02318c3_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `v` int DEFAULT '0',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_default_a02318c3
< ok 0 0
> exec
| CREATE DATABASE pa_modify_default_a02318c3 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_default_a02318c3
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v INT DEFAULT 0)
< ok 0 0
> exec
| ALTER TABLE `pa_modify_default_a02318c3`.`t` MODIFY COLUMN `v` int DEFAULT '1'
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_default_a02318c3.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `v` int DEFAULT '1',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationApplyCorrectness/8.0/modify-default-bigint-unsigned-max default#1
> exec
| DROP DATABASE IF EXISTS pa_modify_default_bigint_unsigned_max_983e5e30_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_default_bigint_unsigned_max_983e5e30_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_default_bigint_unsigned_max_983e5e30_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, big BIGINT UNSIGNED NOT NULL DEFAULT 18446744073709551615)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_default_bigint_unsigned_max_983e5e30_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `big` bigint unsigned NOT NULL DEFAULT '18446744073709551615',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_default_bigint_unsigned_max_983e5e30_from
< ok 0 0
> exec
| CREATE DATABASE pa_modify_default_bigint_unsigned_max_983e5e30_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_default_bigint_unsigned_max_983e5e30_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, big BIGINT UNSIGNED NOT NULL DEFAULT 0)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_default_bigint_unsigned_max_983e5e30_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `big` bigint unsigned NOT NULL DEFAULT '0',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_default_bigint_unsigned_max_983e5e30
< ok 0 0
> exec
| CREATE DATABASE pa_modify_default_bigint_unsigned_max_983e5e30 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_default_bigint_unsigned_max_983e5e30
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, big BIGINT UNSIGNED NOT NULL DEFAULT 0)
< ok 0 0
> exec
| ALTER TABLE `pa_modify_default_bigint_unsigned_max_983e5e30`.`t` MODIFY COLUMN `big` bigint unsigned NOT NULL DEFAULT '18446744073709551615'
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_default_bigint_unsigned_max_983e5e30.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `big` bigint unsigned NOT NULL DEFAULT '18446744073709551615',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationApplyCorrectness/8.0/modify-default-bit-hex default#1
> exec
| DROP DATABASE IF EXISTS pa_modify_default_bit_hex_ee7dda25_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_default_bit_hex_ee7dda25_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_default_bit_hex_ee7dda25_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, f BIT(8) NOT NULL DEFAULT 0x05)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_default_bit_hex_ee7dda25_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `f` bit(8) NOT NULL DEFAULT b'101',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_default_bit_hex_ee7dda25_from
< ok 0 0
> exec
| CREATE DATABASE pa_modify_default_bit_hex_ee7dda25_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_default_bit_hex_ee7dda25_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, f BIT(8) NOT NULL DEFAULT b'0')
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_default_bit_hex_ee7dda25_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `f` bit(8) NOT NULL DEFAULT b'0',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_default_bit_hex_ee7dda25
< ok 0 0
> exec
| CREATE DATABASE pa_modify_default_bit_hex_ee7dda25 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_default_bit_hex_ee7dda25
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, f BIT(8) NOT NULL DEFAULT b'0')
< ok 0 0
> exec
| ALTER TABLE `pa_modify_default_bit_hex_ee7dda25`.`t` MODIFY COLUMN `f` bit(8) NOT NULL DEFAULT b'101'
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_default_bit_hex_ee7dda25.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `f` bit(8) NOT NULL DEFAULT b'101',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationApplyCorrectness/8.0/modify-default-year default#1
> exec
| DROP DATABASE IF EXISTS pa_modify_default_year_b2dade97_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_default_year_b2dade97_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_default_year_b2dade97_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, y YEAR NOT NULL DEFAULT 2000)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_default_year_b2dade97_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `y` year NOT NULL DEFAULT '2000',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_default_year_b2dade97_from
< ok 0 0
> exec
| CREATE DATABASE pa_modify_default_year_b2dade97_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_default_year_b2dade97_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, y YEAR NOT NULL DEFAULT 1999)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_default_year_b2dade97_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `y` year NOT NULL DEFAULT '1999',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_default_year_b2dade97
< ok 0 0
> exec
| CREATE DATABASE pa_modify_default_year_b2dade97 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_default_year_b2dade97
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, y YEAR NOT NULL DEFAULT 1999)
< ok 0 0
> exec
| ALTER TABLE `pa_modify_default_year_b2dade97`.`t` MODIFY COLUMN `y` year NOT NULL DEFAULT '2000'
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_default_year_b2dade97.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `y` year NOT NULL DEFAULT '2000',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationApplyCorrectness/8.0/modify-generated-expr default#1
> exec
| DROP DATABASE IF EXISTS pa_modify_generated_expr_15bcfac1_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_generated_expr_15bcfac1_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_generated_expr_15bcfac1_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, g INT GENERATED ALWAYS AS (a+2) STORED)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_generated_expr_15bcfac1_to.t
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
|   `g` int GENERATED ALWAYS AS ((`a` + 2)) STORED,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_generated_expr_15bcfac1_from
< ok 0 0
> exec
| CREATE DATABASE pa_modify_generated_expr_15bcfac1_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_generated_expr_15bcfac1_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, g INT GENERATED ALWAYS AS (a+1) STORED)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_generated_expr_15bcfac1_from.t
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
|   `g` int GENERATED ALWAYS AS ((`a` + 1)) STORED,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_generated_expr_15bcfac1
< ok 0 0
> exec
| CREATE DATABASE pa_modify_generated_expr_15bcfac1 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_generated_expr_15bcfac1
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, g INT GENERATED ALWAYS AS (a+1) STORED)
< ok 0 0
> exec
| ALTER TABLE `pa_modify_generated_expr_15bcfac1`.`t` MODIFY COLUMN `g` int GENERATED ALWAYS AS (((`a` + 2))) STORED
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_generated_expr_15bcfac1.t
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
|   `g` int GENERATED ALWAYS AS ((`a` + 2)) STORED,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationApplyCorrectness/8.0/modify-nullability default#1
> exec
| DROP DATABASE IF EXISTS pa_modify_nullability_fc28fda3_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_nullability_fc28fda3_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_nullability_fc28fda3_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v INT NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_nullability_fc28fda3_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `v` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_nullability_fc28fda3_from
< ok 0 0
> exec
| CREATE DATABASE pa_modify_nullability_fc28fda3_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_nullability_fc28fda3_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v INT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_nullability_fc28fda3_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `v` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_nullability_fc28fda3
< ok 0 0
> exec
| CREATE DATABASE pa_modify_nullability_fc28fda3 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_nullability_fc28fda3
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v INT NULL)
< ok 0 0
> exec
| ALTER TABLE `pa_modify_nullability_fc28fda3`.`t` MODIFY COLUMN `v` int NOT NULL
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_nullability_fc28fda3.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `v` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationApplyCorrectness/8.0/modify-type default#1
> exec
| DROP DATABASE IF EXISTS pa_modify_type_e026e43e_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_type_e026e43e_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_type_e026e43e_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v BIGINT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_type_e026e43e_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `v` bigint DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_type_e026e43e_from
< ok 0 0
> exec
| CREATE DATABASE pa_modify_type_e026e43e_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_type_e026e43e_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_type_e026e43e_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `v` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_type_e026e43e
< ok 0 0
> exec
| CREATE DATABASE pa_modify_type_e026e43e DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_type_e026e43e
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v INT)
< ok 0 0
> exec
| ALTER TABLE `pa_modify_type_e026e43e`.`t` MODIFY COLUMN `v` bigint DEFAULT NULL
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_type_e026e43e.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `v` bigint DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationApplyCorrectness/8.0/modify-type-width-57-vs-80 default#1
> exec
| DROP DATABASE IF EXISTS pa_modify_type_width_57_vs_80_6c45b0b5_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_type_width_57_vs_80_6c45b0b5_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_type_width_57_vs_80_6c45b0b5_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v BIGINT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_type_width_57_vs_80_6c45b0b5_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `v` bigint DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_type_width_57_vs_80_6c45b0b5_from
< ok 0 0
> exec
| CREATE DATABASE pa_modify_type_width_57_vs_80_6c45b0b5_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_type_width_57_vs_80_6c45b0b5_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v SMALLINT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_type_width_57_vs_80_6c45b0b5_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `v` smallint DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_modify_type_width_57_vs_80_6c45b0b5
< ok 0 0
> exec
| CREATE DATABASE pa_modify_type_width_57_vs_80_6c45b0b5 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_type_width_57_vs_80_6c45b0b5
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v SMALLINT)
< ok 0 0
> exec
| ALTER TABLE `pa_modify_type_width_57_vs_80_6c45b0b5`.`t` MODIFY COLUMN `v` bigint DEFAULT NULL
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_type_width_57_vs_80_6c45b0b5.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `v` bigint DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
