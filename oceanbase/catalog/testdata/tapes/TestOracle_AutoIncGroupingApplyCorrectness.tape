-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_AutoIncGroupingApplyCorrectness default#1
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_uk_951d8d3f_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_uk_951d8d3f_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_uk_951d8d3f_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT, seq BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_seq (seq))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_uk_951d8d3f_to.t
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
|   `seq` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_seq` (`seq`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_uk_951d8d3f
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_uk_951d8d3f DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_uk_951d8d3f
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_uk_951d8d3f.t
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
| ALTER TABLE `pa_aig_add_col_uk_951d8d3f`.`t` ADD COLUMN `seq` bigint unsigned NOT NULL AUTO_INCREMENT
< ok 0 0
> exec
| ALTER TABLE `pa_aig_add_col_uk_951d8d3f`.`t` ADD UNIQUE KEY `uk_seq` (`seq`) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_uk_951d8d3f.t
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
|   `seq` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_seq` (`seq`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_plain_key_7d8da3cb_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_plain_key_7d8da3cb_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_plain_key_7d8da3cb_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT, seq BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, KEY k_seq (seq))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_plain_key_7d8da3cb_to.t
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
|   `seq` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k_seq` (`seq`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_plain_key_7d8da3cb
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_plain_key_7d8da3cb DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_plain_key_7d8da3cb
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_plain_key_7d8da3cb.t
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
| ALTER TABLE `pa_aig_add_col_plain_key_7d8da3cb`.`t` ADD COLUMN `seq` bigint unsigned NOT NULL AUTO_INCREMENT
< ok 0 0
> exec
| ALTER TABLE `pa_aig_add_col_plain_key_7d8da3cb`.`t` ADD KEY `k_seq` (`seq`) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_plain_key_7d8da3cb.t
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
|   `seq` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k_seq` (`seq`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_new_pk_cff5f1cb_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_new_pk_cff5f1cb_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_new_pk_cff5f1cb_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, a INT, seq BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, PRIMARY KEY (seq))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_new_pk_cff5f1cb_to.t
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
|   `seq` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`seq`)
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_new_pk_cff5f1cb
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_new_pk_cff5f1cb DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_new_pk_cff5f1cb
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_new_pk_cff5f1cb.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_aig_add_col_new_pk_cff5f1cb`.`t` ADD COLUMN `seq` bigint unsigned NOT NULL AUTO_INCREMENT
< ok 0 0
> exec
| ALTER TABLE `pa_aig_add_col_new_pk_cff5f1cb`.`t` ADD PRIMARY KEY (`seq`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_new_pk_cff5f1cb.t
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
|   `seq` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`seq`)
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_promote_pk_e88e031_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_promote_pk_e88e031_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_promote_pk_e88e031_to
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT, seq BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, PRIMARY KEY (seq))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_promote_pk_e88e031_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) NOT NULL,
|   `b` int(11) DEFAULT NULL,
|   `seq` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`seq`)
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_promote_pk_e88e031
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_promote_pk_e88e031 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_promote_pk_e88e031
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL PRIMARY KEY, b INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_promote_pk_e88e031.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) NOT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`a`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_aig_add_col_promote_pk_e88e031`.`t` ADD COLUMN `seq` bigint unsigned NOT NULL AUTO_INCREMENT
< ok 0 0
> exec
| ALTER TABLE `pa_aig_add_col_promote_pk_e88e031`.`t` DROP PRIMARY KEY, ADD PRIMARY KEY (`seq`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_promote_pk_e88e031.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) NOT NULL,
|   `b` int(11) DEFAULT NULL,
|   `seq` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`seq`)
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_promote_pk_drop_old_50e32e19_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_promote_pk_drop_old_50e32e19_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_promote_pk_drop_old_50e32e19_to
< ok 0 0
> exec
| CREATE TABLE t (b INT, seq BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, PRIMARY KEY (seq))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_promote_pk_drop_old_50e32e19_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) DEFAULT NULL,
|   `seq` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`seq`)
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_promote_pk_drop_old_50e32e19
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_promote_pk_drop_old_50e32e19 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_promote_pk_drop_old_50e32e19
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL PRIMARY KEY, b INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_promote_pk_drop_old_50e32e19.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) NOT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`a`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_aig_add_col_promote_pk_drop_old_50e32e19`.`t` ADD COLUMN `seq` bigint unsigned NOT NULL AUTO_INCREMENT
< ok 0 0
> exec
| ALTER TABLE `pa_aig_add_col_promote_pk_drop_old_50e32e19`.`t` DROP PRIMARY KEY, ADD PRIMARY KEY (`seq`)
< ok 0 0
> exec
| ALTER TABLE `pa_aig_add_col_promote_pk_drop_old_50e32e19`.`t` DROP COLUMN `a`
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_promote_pk_drop_old_50e32e19.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) DEFAULT NULL,
|   `seq` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`seq`)
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_modify_gain_ai_new_key_da2d751e_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_modify_gain_ai_new_key_da2d751e_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_modify_gain_ai_new_key_da2d751e_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_c (c))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_modify_gain_ai_new_key_da2d751e_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_c` (`c`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_modify_gain_ai_new_key_da2d751e
< ok 0 0
> exec
| CREATE DATABASE pa_aig_modify_gain_ai_new_key_da2d751e DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_modify_gain_ai_new_key_da2d751e
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_modify_gain_ai_new_key_da2d751e.t
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
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_aig_modify_gain_ai_new_key_da2d751e`.`t` MODIFY COLUMN `c` int NOT NULL AUTO_INCREMENT
< ok 0 0
> exec
| ALTER TABLE `pa_aig_modify_gain_ai_new_key_da2d751e`.`t` ADD UNIQUE KEY `uk_c` (`c`) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_modify_gain_ai_new_key_da2d751e.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_c` (`c`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_modify_gain_ai_existing_key_883014d5_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_modify_gain_ai_existing_key_883014d5_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_modify_gain_ai_existing_key_883014d5_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_c (c))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_modify_gain_ai_existing_key_883014d5_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_c` (`c`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_modify_gain_ai_existing_key_883014d5
< ok 0 0
> exec
| CREATE DATABASE pa_aig_modify_gain_ai_existing_key_883014d5 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_modify_gain_ai_existing_key_883014d5
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL, UNIQUE KEY uk_c (c))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_modify_gain_ai_existing_key_883014d5.t
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
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_c` (`c`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_aig_modify_gain_ai_existing_key_883014d5`.`t` MODIFY COLUMN `c` int NOT NULL AUTO_INCREMENT
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_modify_gain_ai_existing_key_883014d5.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_c` (`c`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_composite_backing_49d278b6_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_composite_backing_49d278b6_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_composite_backing_49d278b6_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, seq BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, d INT NOT NULL, UNIQUE KEY uk_sd (seq, d))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_composite_backing_49d278b6_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `seq` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
|   `d` int(11) NOT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_sd` (`seq`, `d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_composite_backing_49d278b6
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_composite_backing_49d278b6 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_composite_backing_49d278b6
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_composite_backing_49d278b6.t
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
> exec
| ALTER TABLE `pa_aig_add_col_composite_backing_49d278b6`.`t` ADD COLUMN `d` int NOT NULL
< ok 0 0
> exec
| ALTER TABLE `pa_aig_add_col_composite_backing_49d278b6`.`t` ADD COLUMN `seq` bigint unsigned NOT NULL AUTO_INCREMENT
< ok 0 0
> exec
| ALTER TABLE `pa_aig_add_col_composite_backing_49d278b6`.`t` ADD UNIQUE KEY `uk_sd` (`seq`,`d`) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_composite_backing_49d278b6.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `d` int(11) NOT NULL,
|   `seq` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_sd` (`seq`, `d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_ai_migrate_columns_6b9f75b4_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_ai_migrate_columns_6b9f75b4_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_ai_migrate_columns_6b9f75b4_to
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_a (a), UNIQUE KEY uk_b (b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_ai_migrate_columns_6b9f75b4_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) NOT NULL,
|   `b` int(11) NOT NULL AUTO_INCREMENT,
|   UNIQUE KEY `uk_a` (`a`) BLOCK_SIZE 16384 LOCAL,
|   UNIQUE KEY `uk_b` (`b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_ai_migrate_columns_6b9f75b4
< ok 0 0
> exec
| CREATE DATABASE pa_aig_ai_migrate_columns_6b9f75b4 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_ai_migrate_columns_6b9f75b4
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL AUTO_INCREMENT, b INT NOT NULL, UNIQUE KEY uk_a (a), UNIQUE KEY uk_b (b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_ai_migrate_columns_6b9f75b4.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) NOT NULL AUTO_INCREMENT,
|   `b` int(11) NOT NULL,
|   UNIQUE KEY `uk_a` (`a`) BLOCK_SIZE 16384 LOCAL,
|   UNIQUE KEY `uk_b` (`b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_aig_ai_migrate_columns_6b9f75b4`.`t` MODIFY COLUMN `a` int NOT NULL
< ok 0 0
> exec
| ALTER TABLE `pa_aig_ai_migrate_columns_6b9f75b4`.`t` MODIFY COLUMN `b` int NOT NULL AUTO_INCREMENT
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_ai_migrate_columns_6b9f75b4.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) NOT NULL,
|   `b` int(11) NOT NULL AUTO_INCREMENT,
|   UNIQUE KEY `uk_a` (`a`) BLOCK_SIZE 16384 LOCAL,
|   UNIQUE KEY `uk_b` (`b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_ai_migrate_to_new_col_649c6790_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_ai_migrate_to_new_col_649c6790_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_ai_migrate_to_new_col_649c6790_to
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_a (a), UNIQUE KEY uk_b (b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_ai_migrate_to_new_col_649c6790_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) NOT NULL,
|   `b` int(11) NOT NULL AUTO_INCREMENT,
|   UNIQUE KEY `uk_a` (`a`) BLOCK_SIZE 16384 LOCAL,
|   UNIQUE KEY `uk_b` (`b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_ai_migrate_to_new_col_649c6790
< ok 0 0
> exec
| CREATE DATABASE pa_aig_ai_migrate_to_new_col_649c6790 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_ai_migrate_to_new_col_649c6790
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_a (a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_ai_migrate_to_new_col_649c6790.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) NOT NULL AUTO_INCREMENT,
|   UNIQUE KEY `uk_a` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_aig_ai_migrate_to_new_col_649c6790`.`t` MODIFY COLUMN `a` int NOT NULL
< ok 0 0
> exec
| ALTER TABLE `pa_aig_ai_migrate_to_new_col_649c6790`.`t` ADD COLUMN `b` int NOT NULL AUTO_INCREMENT
< ok 0 0
> exec
| ALTER TABLE `pa_aig_ai_migrate_to_new_col_649c6790`.`t` ADD UNIQUE KEY `uk_b` (`b`) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_ai_migrate_to_new_col_649c6790.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) NOT NULL,
|   `b` int(11) NOT NULL AUTO_INCREMENT,
|   UNIQUE KEY `uk_a` (`a`) BLOCK_SIZE 16384 LOCAL,
|   UNIQUE KEY `uk_b` (`b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_with_other_changes_7dc04c8c_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_with_other_changes_7dc04c8c_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_with_other_changes_7dc04c8c_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT, e VARCHAR(20), f INT, seq BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_seq (seq))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_with_other_changes_7dc04c8c_to.t
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
|   `e` varchar(20) DEFAULT NULL,
|   `f` int(11) DEFAULT NULL,
|   `seq` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_seq` (`seq`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_with_other_changes_7dc04c8c
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_with_other_changes_7dc04c8c DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_with_other_changes_7dc04c8c
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT, e VARCHAR(10))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_with_other_changes_7dc04c8c.t
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
|   `e` varchar(10) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_aig_add_col_with_other_changes_7dc04c8c`.`t` MODIFY COLUMN `e` varchar(20) DEFAULT NULL
< ok 0 0
> exec
| ALTER TABLE `pa_aig_add_col_with_other_changes_7dc04c8c`.`t` ADD COLUMN `f` int DEFAULT NULL
< ok 0 0
> exec
| ALTER TABLE `pa_aig_add_col_with_other_changes_7dc04c8c`.`t` ADD COLUMN `seq` bigint unsigned NOT NULL AUTO_INCREMENT
< ok 0 0
> exec
| ALTER TABLE `pa_aig_add_col_with_other_changes_7dc04c8c`.`t` ADD UNIQUE KEY `uk_seq` (`seq`) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_with_other_changes_7dc04c8c.t
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
|   `e` varchar(20) DEFAULT NULL,
|   `f` int(11) DEFAULT NULL,
|   `seq` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_seq` (`seq`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_myisam_add_nonleftmost_da0317da_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_myisam_add_nonleftmost_da0317da_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_myisam_add_nonleftmost_da0317da_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, seq INT NOT NULL AUTO_INCREMENT, KEY k_as (a, seq)) ENGINE=MyISAM
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_myisam_add_nonleftmost_da0317da_to.t
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
|   `seq` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k_as` (`a`, `seq`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_myisam_add_nonleftmost_da0317da
< ok 0 0
> exec
| CREATE DATABASE pa_aig_myisam_add_nonleftmost_da0317da DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_myisam_add_nonleftmost_da0317da
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT NOT NULL) ENGINE=MyISAM
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_myisam_add_nonleftmost_da0317da.t
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
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_aig_myisam_add_nonleftmost_da0317da`.`t` ADD COLUMN `seq` int NOT NULL AUTO_INCREMENT
< ok 0 0
> exec
| ALTER TABLE `pa_aig_myisam_add_nonleftmost_da0317da`.`t` ADD KEY `k_as` (`a`,`seq`) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_myisam_add_nonleftmost_da0317da.t
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
|   `seq` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k_as` (`a`, `seq`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_myisam_gain_ai_existing_nonleftmost_4ba19dd5_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_myisam_gain_ai_existing_nonleftmost_4ba19dd5_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_myisam_gain_ai_existing_nonleftmost_4ba19dd5_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, c INT NOT NULL AUTO_INCREMENT, KEY k (a, c)) ENGINE=MyISAM
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_myisam_gain_ai_existing_nonleftmost_4ba19dd5_to.t
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
|   `c` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k` (`a`, `c`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_myisam_gain_ai_existing_nonleftmost_4ba19dd5
< ok 0 0
> exec
| CREATE DATABASE pa_aig_myisam_gain_ai_existing_nonleftmost_4ba19dd5 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_myisam_gain_ai_existing_nonleftmost_4ba19dd5
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, c INT NOT NULL, KEY k (a, c)) ENGINE=MyISAM
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_myisam_gain_ai_existing_nonleftmost_4ba19dd5.t
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
|   `c` int(11) NOT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`a`, `c`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_aig_myisam_gain_ai_existing_nonleftmost_4ba19dd5`.`t` MODIFY COLUMN `c` int NOT NULL AUTO_INCREMENT
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_myisam_gain_ai_existing_nonleftmost_4ba19dd5.t
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
|   `c` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k` (`a`, `c`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_myisam_shared_key_ai_migrate_14cd9661_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_myisam_shared_key_ai_migrate_14cd9661_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_myisam_shared_key_ai_migrate_14cd9661_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c2 INT NOT NULL AUTO_INCREMENT, KEY k2 (c2)) ENGINE=MyISAM
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_myisam_shared_key_ai_migrate_14cd9661_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c2` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k2` (`c2`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_myisam_shared_key_ai_migrate_14cd9661
< ok 0 0
> exec
| CREATE DATABASE pa_aig_myisam_shared_key_ai_migrate_14cd9661 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_myisam_shared_key_ai_migrate_14cd9661
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c1 INT NOT NULL AUTO_INCREMENT, c2 INT NOT NULL, KEY k (c1, c2)) ENGINE=MyISAM
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_myisam_shared_key_ai_migrate_14cd9661.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c1` int(11) NOT NULL AUTO_INCREMENT,
|   `c2` int(11) NOT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`c1`, `c2`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_aig_myisam_shared_key_ai_migrate_14cd9661`.`t` DROP INDEX `k`
< ok 0 0
> exec
| ALTER TABLE `pa_aig_myisam_shared_key_ai_migrate_14cd9661`.`t` DROP COLUMN `c1`
< ok 0 0
> exec
| ALTER TABLE `pa_aig_myisam_shared_key_ai_migrate_14cd9661`.`t` MODIFY COLUMN `c2` int NOT NULL AUTO_INCREMENT
< ok 0 0
> exec
| ALTER TABLE `pa_aig_myisam_shared_key_ai_migrate_14cd9661`.`t` ADD KEY `k2` (`c2`) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_myisam_shared_key_ai_migrate_14cd9661.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c2` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k2` (`c2`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_myisam_shared_key_ai_migrate_survive_37600f46_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_myisam_shared_key_ai_migrate_survive_37600f46_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_myisam_shared_key_ai_migrate_survive_37600f46_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL AUTO_INCREMENT, KEY k2 (b)) ENGINE=MyISAM
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_myisam_shared_key_ai_migrate_survive_37600f46_to.t
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
|   `b` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k2` (`b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_myisam_shared_key_ai_migrate_survive_37600f46
< ok 0 0
> exec
| CREATE DATABASE pa_aig_myisam_shared_key_ai_migrate_survive_37600f46 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_myisam_shared_key_ai_migrate_survive_37600f46
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT NOT NULL AUTO_INCREMENT, b INT NOT NULL, KEY k (a, b)) ENGINE=MyISAM
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_myisam_shared_key_ai_migrate_survive_37600f46.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) NOT NULL AUTO_INCREMENT,
|   `b` int(11) NOT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`a`, `b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_aig_myisam_shared_key_ai_migrate_survive_37600f46`.`t` DROP INDEX `k`
< ok 0 0
> exec
| ALTER TABLE `pa_aig_myisam_shared_key_ai_migrate_survive_37600f46`.`t` MODIFY COLUMN `a` int NOT NULL
< ok 0 0
> exec
| ALTER TABLE `pa_aig_myisam_shared_key_ai_migrate_survive_37600f46`.`t` MODIFY COLUMN `b` int NOT NULL AUTO_INCREMENT
< ok 0 0
> exec
| ALTER TABLE `pa_aig_myisam_shared_key_ai_migrate_survive_37600f46`.`t` ADD KEY `k2` (`b`) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_myisam_shared_key_ai_migrate_survive_37600f46.t
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
|   `b` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k2` (`b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_desc_key_ba13ecbc_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_desc_key_ba13ecbc_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_desc_key_ba13ecbc_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, seq INT NOT NULL AUTO_INCREMENT, KEY k_seq (seq DESC))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_desc_key_ba13ecbc_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `seq` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k_seq` (`seq`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_desc_key_ba13ecbc
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_desc_key_ba13ecbc DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_desc_key_ba13ecbc
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_desc_key_ba13ecbc.t
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
> exec
| ALTER TABLE `pa_aig_add_col_desc_key_ba13ecbc`.`t` ADD COLUMN `seq` int NOT NULL AUTO_INCREMENT
< ok 0 0
> exec
| ALTER TABLE `pa_aig_add_col_desc_key_ba13ecbc`.`t` ADD KEY `k_seq` (`seq`) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_desc_key_ba13ecbc.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `seq` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k_seq` (`seq`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_invisible_key_2f4e3648_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_invisible_key_2f4e3648_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_invisible_key_2f4e3648_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, seq INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_seq (seq) INVISIBLE)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_invisible_key_2f4e3648_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `seq` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_seq` (`seq`) BLOCK_SIZE 16384 LOCAL /*!80000 INVISIBLE */
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_invisible_key_2f4e3648
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_invisible_key_2f4e3648 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_invisible_key_2f4e3648
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_invisible_key_2f4e3648.t
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
> exec
| ALTER TABLE `pa_aig_add_col_invisible_key_2f4e3648`.`t` ADD COLUMN `seq` int NOT NULL AUTO_INCREMENT
< ok 0 0
> exec
| ALTER TABLE `pa_aig_add_col_invisible_key_2f4e3648`.`t` ADD UNIQUE KEY `uk_seq` (`seq`) /*!80000 INVISIBLE */ LOCAL
< ok 0 0
> exec
| ALTER TABLE `pa_aig_add_col_invisible_key_2f4e3648`.`t` ALTER INDEX `uk_seq` INVISIBLE
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_invisible_key_2f4e3648.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `seq` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_seq` (`seq`) BLOCK_SIZE 16384 LOCAL /*!80000 INVISIBLE */
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_trailing_expr_key_a11cd10f_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_trailing_expr_key_a11cd10f_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_trailing_expr_key_a11cd10f_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT, seq INT NOT NULL AUTO_INCREMENT, KEY k_se (seq, (a + 1)))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_trailing_expr_key_a11cd10f_to.t
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
|   `seq` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k_se` (`seq`, (`a` + 1)) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_add_col_trailing_expr_key_a11cd10f
< ok 0 0
> exec
| CREATE DATABASE pa_aig_add_col_trailing_expr_key_a11cd10f DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_add_col_trailing_expr_key_a11cd10f
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_trailing_expr_key_a11cd10f.t
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
| ALTER TABLE `pa_aig_add_col_trailing_expr_key_a11cd10f`.`t` ADD COLUMN `seq` int NOT NULL AUTO_INCREMENT
< ok 0 0
> exec
| ALTER TABLE `pa_aig_add_col_trailing_expr_key_a11cd10f`.`t` ADD KEY `k_se` (`seq`,((`a` + 1))) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_add_col_trailing_expr_key_a11cd10f.t
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
|   `seq` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k_se` (`seq`, (`a` + 1)) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_drop_key_and_column_b9d91a44_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_drop_key_and_column_b9d91a44_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_drop_key_and_column_b9d91a44_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_drop_key_and_column_b9d91a44_to.t
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
> exec
| DROP DATABASE IF EXISTS pa_aig_drop_key_and_column_b9d91a44
< ok 0 0
> exec
| CREATE DATABASE pa_aig_drop_key_and_column_b9d91a44 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_drop_key_and_column_b9d91a44
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_c (c))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_drop_key_and_column_b9d91a44.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_c` (`c`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_aig_drop_key_and_column_b9d91a44`.`t` DROP INDEX `uk_c`
< ok 0 0
> exec
| ALTER TABLE `pa_aig_drop_key_and_column_b9d91a44`.`t` DROP COLUMN `c`
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_drop_key_and_column_b9d91a44.t
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
> exec
| DROP DATABASE IF EXISTS pa_aig_drop_key_deai_62647d4b_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_drop_key_deai_62647d4b_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_drop_key_deai_62647d4b_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_drop_key_deai_62647d4b_to.t
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
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_drop_key_deai_62647d4b
< ok 0 0
> exec
| CREATE DATABASE pa_aig_drop_key_deai_62647d4b DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_drop_key_deai_62647d4b
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_c (c))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_drop_key_deai_62647d4b.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_c` (`c`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_aig_drop_key_deai_62647d4b`.`t` DROP INDEX `uk_c`
< ok 0 0
> exec
| ALTER TABLE `pa_aig_drop_key_deai_62647d4b`.`t` MODIFY COLUMN `c` int NOT NULL
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_drop_key_deai_62647d4b.t
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
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_replace_backing_key_c0765476_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_replace_backing_key_c0765476_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_replace_backing_key_c0765476_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, UNIQUE KEY uk_c (c, a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_replace_backing_key_c0765476_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` int(11) NOT NULL AUTO_INCREMENT,
|   `a` int(11) NOT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_c` (`c`, `a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_replace_backing_key_c0765476
< ok 0 0
> exec
| CREATE DATABASE pa_aig_replace_backing_key_c0765476 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_replace_backing_key_c0765476
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, UNIQUE KEY uk_c (c))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_replace_backing_key_c0765476.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` int(11) NOT NULL AUTO_INCREMENT,
|   `a` int(11) NOT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_c` (`c`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_aig_replace_backing_key_c0765476`.`t` DROP INDEX `uk_c`
< ok 0 0
> exec
| ALTER TABLE `pa_aig_replace_backing_key_c0765476`.`t` ADD UNIQUE KEY `uk_c` (`c`,`a`) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_replace_backing_key_c0765476.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` int(11) NOT NULL AUTO_INCREMENT,
|   `a` int(11) NOT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_c` (`c`, `a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_swap_backing_key_15f0ecd_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_swap_backing_key_15f0ecd_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_swap_backing_key_15f0ecd_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, KEY k_c (c))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_swap_backing_key_15f0ecd_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k_c` (`c`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_swap_backing_key_15f0ecd
< ok 0 0
> exec
| CREATE DATABASE pa_aig_swap_backing_key_15f0ecd DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_swap_backing_key_15f0ecd
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_c (c))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_swap_backing_key_15f0ecd.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_c` (`c`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_aig_swap_backing_key_15f0ecd`.`t` DROP INDEX `uk_c`
< ok 0 0
> exec
| ALTER TABLE `pa_aig_swap_backing_key_15f0ecd`.`t` ADD KEY `k_c` (`c`) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_swap_backing_key_15f0ecd.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k_c` (`c`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_multi_drop_plus_add_3ee9f702_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_multi_drop_plus_add_3ee9f702_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_multi_drop_plus_add_3ee9f702_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk3 (c))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_multi_drop_plus_add_3ee9f702_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk3` (`c`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_multi_drop_plus_add_3ee9f702
< ok 0 0
> exec
| CREATE DATABASE pa_aig_multi_drop_plus_add_3ee9f702 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_multi_drop_plus_add_3ee9f702
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk1 (c), KEY k2 (c))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_multi_drop_plus_add_3ee9f702.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk1` (`c`) BLOCK_SIZE 16384 LOCAL,
|   KEY `k2` (`c`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_aig_multi_drop_plus_add_3ee9f702`.`t` DROP INDEX `k2`
< ok 0 0
> exec
| ALTER TABLE `pa_aig_multi_drop_plus_add_3ee9f702`.`t` DROP INDEX `uk1`
< ok 0 0
> exec
| ALTER TABLE `pa_aig_multi_drop_plus_add_3ee9f702`.`t` ADD UNIQUE KEY `uk3` (`c`) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_multi_drop_plus_add_3ee9f702.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk3` (`c`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_drop_covered_elsewhere_ecbfef83_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_drop_covered_elsewhere_ecbfef83_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_drop_covered_elsewhere_ecbfef83_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, KEY k2 (c))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_drop_covered_elsewhere_ecbfef83_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k2` (`c`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_drop_covered_elsewhere_ecbfef83
< ok 0 0
> exec
| CREATE DATABASE pa_aig_drop_covered_elsewhere_ecbfef83 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_drop_covered_elsewhere_ecbfef83
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk1 (c), KEY k2 (c))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_drop_covered_elsewhere_ecbfef83.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk1` (`c`) BLOCK_SIZE 16384 LOCAL,
|   KEY `k2` (`c`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_aig_drop_covered_elsewhere_ecbfef83`.`t` DROP INDEX `uk1`
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_drop_covered_elsewhere_ecbfef83.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k2` (`c`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_drop_pk_and_column_ce8865da_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_drop_pk_and_column_ce8865da_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_drop_pk_and_column_ce8865da_to
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_drop_pk_and_column_ce8865da_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) NOT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_drop_pk_and_column_ce8865da
< ok 0 0
> exec
| CREATE DATABASE pa_aig_drop_pk_and_column_ce8865da DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_drop_pk_and_column_ce8865da
< ok 0 0
> exec
| CREATE TABLE t (c INT NOT NULL AUTO_INCREMENT PRIMARY KEY, a INT NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_drop_pk_and_column_ce8865da.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `c` int(11) NOT NULL AUTO_INCREMENT,
|   `a` int(11) NOT NULL,
|   PRIMARY KEY (`c`)
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_aig_drop_pk_and_column_ce8865da`.`t` DROP PRIMARY KEY
< ok 0 0
> exec
| ALTER TABLE `pa_aig_drop_pk_and_column_ce8865da`.`t` DROP COLUMN `c`
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_drop_pk_and_column_ce8865da.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) NOT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_drop_pk_deai_c594c6d1_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_drop_pk_deai_c594c6d1_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_drop_pk_deai_c594c6d1_to
< ok 0 0
> exec
| CREATE TABLE t (c INT NOT NULL, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_drop_pk_deai_c594c6d1_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `c` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_drop_pk_deai_c594c6d1
< ok 0 0
> exec
| CREATE DATABASE pa_aig_drop_pk_deai_c594c6d1 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_drop_pk_deai_c594c6d1
< ok 0 0
> exec
| CREATE TABLE t (c INT NOT NULL AUTO_INCREMENT PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_drop_pk_deai_c594c6d1.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `c` int(11) NOT NULL AUTO_INCREMENT,
|   `a` int(11) DEFAULT NULL,
|   PRIMARY KEY (`c`)
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_aig_drop_pk_deai_c594c6d1`.`t` DROP PRIMARY KEY
< ok 0 0
> exec
| ALTER TABLE `pa_aig_drop_pk_deai_c594c6d1`.`t` MODIFY COLUMN `c` int NOT NULL
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_drop_pk_deai_c594c6d1.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `c` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_drop_composite_pk_and_column_629e9684_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_drop_composite_pk_and_column_629e9684_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_drop_composite_pk_and_column_629e9684_to
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_drop_composite_pk_and_column_629e9684_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) NOT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_drop_composite_pk_and_column_629e9684
< ok 0 0
> exec
| CREATE DATABASE pa_aig_drop_composite_pk_and_column_629e9684 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_drop_composite_pk_and_column_629e9684
< ok 0 0
> exec
| CREATE TABLE t (c INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, PRIMARY KEY (c, a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_drop_composite_pk_and_column_629e9684.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `c` int(11) NOT NULL AUTO_INCREMENT,
|   `a` int(11) NOT NULL,
|   PRIMARY KEY (`c`, `a`)
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_aig_drop_composite_pk_and_column_629e9684`.`t` DROP PRIMARY KEY
< ok 0 0
> exec
| ALTER TABLE `pa_aig_drop_composite_pk_and_column_629e9684`.`t` DROP COLUMN `c`
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_drop_composite_pk_and_column_629e9684.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) NOT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_pk_swap_drop_secondary_c64a03c4_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_pk_swap_drop_secondary_c64a03c4_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_pk_swap_drop_secondary_c64a03c4_to
< ok 0 0
> exec
| CREATE TABLE t (x INT NOT NULL, c INT NOT NULL AUTO_INCREMENT, PRIMARY KEY (c))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_pk_swap_drop_secondary_c64a03c4_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `x` int(11) NOT NULL,
|   `c` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`c`)
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_pk_swap_drop_secondary_c64a03c4
< ok 0 0
> exec
| CREATE DATABASE pa_aig_pk_swap_drop_secondary_c64a03c4 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_pk_swap_drop_secondary_c64a03c4
< ok 0 0
> exec
| CREATE TABLE t (x INT NOT NULL, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_c (c), PRIMARY KEY (x))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_pk_swap_drop_secondary_c64a03c4.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `x` int(11) NOT NULL,
|   `c` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`x`),
|   UNIQUE KEY `uk_c` (`c`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_aig_pk_swap_drop_secondary_c64a03c4`.`t` DROP INDEX `uk_c`
< ok 0 0
> exec
| ALTER TABLE `pa_aig_pk_swap_drop_secondary_c64a03c4`.`t` DROP PRIMARY KEY, ADD PRIMARY KEY (`c`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_pk_swap_drop_secondary_c64a03c4.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `x` int(11) NOT NULL,
|   `c` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`c`)
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_pk_relocate_ai_to_secondary_fb6f4a2d_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_pk_relocate_ai_to_secondary_fb6f4a2d_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_pk_relocate_ai_to_secondary_fb6f4a2d_to
< ok 0 0
> exec
| CREATE TABLE t (c INT NOT NULL AUTO_INCREMENT, b INT NOT NULL, x INT, PRIMARY KEY (b), KEY zz_c (c))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_pk_relocate_ai_to_secondary_fb6f4a2d_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `c` int(11) NOT NULL AUTO_INCREMENT,
|   `b` int(11) NOT NULL,
|   `x` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `zz_c` (`c`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_pk_relocate_ai_to_secondary_fb6f4a2d
< ok 0 0
> exec
| CREATE DATABASE pa_aig_pk_relocate_ai_to_secondary_fb6f4a2d DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_pk_relocate_ai_to_secondary_fb6f4a2d
< ok 0 0
> exec
| CREATE TABLE t (c INT NOT NULL AUTO_INCREMENT, b INT NOT NULL, x INT, PRIMARY KEY (c))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_pk_relocate_ai_to_secondary_fb6f4a2d.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `c` int(11) NOT NULL AUTO_INCREMENT,
|   `b` int(11) NOT NULL,
|   `x` int(11) DEFAULT NULL,
|   PRIMARY KEY (`c`)
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_aig_pk_relocate_ai_to_secondary_fb6f4a2d`.`t` DROP PRIMARY KEY, ADD PRIMARY KEY (`b`)
< ok 0 0
> exec
| ALTER TABLE `pa_aig_pk_relocate_ai_to_secondary_fb6f4a2d`.`t` ADD KEY `zz_c` (`c`) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_pk_relocate_ai_to_secondary_fb6f4a2d.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `c` int(11) NOT NULL AUTO_INCREMENT,
|   `b` int(11) NOT NULL,
|   `x` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `zz_c` (`c`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_pk_keep_ai_demote_member_ae8d6d12_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_pk_keep_ai_demote_member_ae8d6d12_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_pk_keep_ai_demote_member_ae8d6d12_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, a INT NULL, PRIMARY KEY (id))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_pk_keep_ai_demote_member_ae8d6d12_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `a` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_pk_keep_ai_demote_member_ae8d6d12
< ok 0 0
> exec
| CREATE DATABASE pa_aig_pk_keep_ai_demote_member_ae8d6d12 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_pk_keep_ai_demote_member_ae8d6d12
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, PRIMARY KEY (id, a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_pk_keep_ai_demote_member_ae8d6d12.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `a` int(11) NOT NULL,
|   PRIMARY KEY (`id`, `a`)
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_aig_pk_keep_ai_demote_member_ae8d6d12`.`t` DROP PRIMARY KEY, ADD PRIMARY KEY (`id`)
< ok 0 0
> exec
| ALTER TABLE `pa_aig_pk_keep_ai_demote_member_ae8d6d12`.`t` MODIFY COLUMN `a` int DEFAULT NULL
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_pk_keep_ai_demote_member_ae8d6d12.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `a` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_generated_keypart_pullin_e9c71dcd_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_generated_keypart_pullin_e9c71dcd_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_generated_keypart_pullin_e9c71dcd_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, seq INT NOT NULL AUTO_INCREMENT, d INT NOT NULL, g INT GENERATED ALWAYS AS (d+1) VIRTUAL, UNIQUE KEY uk_sg (seq, g))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_generated_keypart_pullin_e9c71dcd_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `seq` int(11) NOT NULL AUTO_INCREMENT,
|   `d` int(11) NOT NULL,
|   `g` int(11) GENERATED ALWAYS AS ((`d` + 1)) VIRTUAL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_sg` (`seq`, `g`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_generated_keypart_pullin_e9c71dcd
< ok 0 0
> exec
| CREATE DATABASE pa_aig_generated_keypart_pullin_e9c71dcd DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_generated_keypart_pullin_e9c71dcd
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_generated_keypart_pullin_e9c71dcd.t
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
> exec
| ALTER TABLE `pa_aig_generated_keypart_pullin_e9c71dcd`.`t` ADD COLUMN `d` int NOT NULL
< ok 0 0
> exec
| ALTER TABLE `pa_aig_generated_keypart_pullin_e9c71dcd`.`t` ADD COLUMN `seq` int NOT NULL AUTO_INCREMENT
< ok 0 0
> exec
| ALTER TABLE `pa_aig_generated_keypart_pullin_e9c71dcd`.`t` ADD COLUMN `g` int GENERATED ALWAYS AS ((`d` + 1)) VIRTUAL
< ok 0 0
> exec
| ALTER TABLE `pa_aig_generated_keypart_pullin_e9c71dcd`.`t` ADD UNIQUE KEY `uk_sg` (`seq`,`g`) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_generated_keypart_pullin_e9c71dcd.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `d` int(11) NOT NULL,
|   `seq` int(11) NOT NULL AUTO_INCREMENT,
|   `g` int(11) GENERATED ALWAYS AS ((`d` + 1)) VIRTUAL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_sg` (`seq`, `g`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_funcpart_new_col_54fac567_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_funcpart_new_col_54fac567_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_funcpart_new_col_54fac567_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, seq INT NOT NULL AUTO_INCREMENT, z INT, KEY k_sz (seq, (z + 1)))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_funcpart_new_col_54fac567_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `seq` int(11) NOT NULL AUTO_INCREMENT,
|   `z` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k_sz` (`seq`, (`z` + 1)) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_funcpart_new_col_54fac567
< ok 0 0
> exec
| CREATE DATABASE pa_aig_funcpart_new_col_54fac567 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_funcpart_new_col_54fac567
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_funcpart_new_col_54fac567.t
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
> exec
| ALTER TABLE `pa_aig_funcpart_new_col_54fac567`.`t` ADD COLUMN `seq` int NOT NULL AUTO_INCREMENT
< ok 0 0
> exec
| ALTER TABLE `pa_aig_funcpart_new_col_54fac567`.`t` ADD COLUMN `z` int DEFAULT NULL
< ok 0 0
> exec
| ALTER TABLE `pa_aig_funcpart_new_col_54fac567`.`t` ADD KEY `k_sz` (`seq`,((`z` + 1))) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_funcpart_new_col_54fac567.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `seq` int(11) NOT NULL AUTO_INCREMENT,
|   `z` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k_sz` (`seq`, (`z` + 1)) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_prefix_over_widened_3c760716_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_prefix_over_widened_3c760716_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_prefix_over_widened_3c760716_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, v VARCHAR(500), seq INT NOT NULL AUTO_INCREMENT, KEY k_sv (seq, v(300)))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_prefix_over_widened_3c760716_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `v` varchar(500) DEFAULT NULL,
|   `seq` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k_sv` (`seq`, `v`(300)) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_prefix_over_widened_3c760716
< ok 0 0
> exec
| CREATE DATABASE pa_aig_prefix_over_widened_3c760716 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_prefix_over_widened_3c760716
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, v VARCHAR(100))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_prefix_over_widened_3c760716.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `v` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_aig_prefix_over_widened_3c760716`.`t` ADD COLUMN `seq` int NOT NULL AUTO_INCREMENT
< ok 0 0
> exec
| ALTER TABLE `pa_aig_prefix_over_widened_3c760716`.`t` MODIFY COLUMN `v` varchar(500) DEFAULT NULL
< ok 0 0
> exec
| ALTER TABLE `pa_aig_prefix_over_widened_3c760716`.`t` ADD KEY `k_sv` (`seq`,`v`(300)) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_prefix_over_widened_3c760716.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `v` varchar(500) DEFAULT NULL,
|   `seq` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `k_sv` (`seq`, `v`(300)) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_pk_drop_member_drop_newkey_d251c96e_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_pk_drop_member_drop_newkey_d251c96e_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_pk_drop_member_drop_newkey_d251c96e_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, KEY k_id (id))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_pk_drop_member_drop_newkey_d251c96e_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   KEY `k_id` (`id`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_pk_drop_member_drop_newkey_d251c96e
< ok 0 0
> exec
| CREATE DATABASE pa_aig_pk_drop_member_drop_newkey_d251c96e DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_pk_drop_member_drop_newkey_d251c96e
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, PRIMARY KEY (id, a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_pk_drop_member_drop_newkey_d251c96e.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `a` int(11) NOT NULL,
|   PRIMARY KEY (`id`, `a`)
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_aig_pk_drop_member_drop_newkey_d251c96e`.`t` DROP PRIMARY KEY
< ok 0 0
> exec
| ALTER TABLE `pa_aig_pk_drop_member_drop_newkey_d251c96e`.`t` DROP COLUMN `a`
< ok 0 0
> exec
| ALTER TABLE `pa_aig_pk_drop_member_drop_newkey_d251c96e`.`t` ADD KEY `k_id` (`id`) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_pk_drop_member_drop_newkey_d251c96e.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   KEY `k_id` (`id`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_pk_shrink_member_drop_fccc9f7c_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_pk_shrink_member_drop_fccc9f7c_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_pk_shrink_member_drop_fccc9f7c_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, PRIMARY KEY (id))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_pk_shrink_member_drop_fccc9f7c_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_pk_shrink_member_drop_fccc9f7c
< ok 0 0
> exec
| CREATE DATABASE pa_aig_pk_shrink_member_drop_fccc9f7c DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_pk_shrink_member_drop_fccc9f7c
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, PRIMARY KEY (id, a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_pk_shrink_member_drop_fccc9f7c.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `a` int(11) NOT NULL,
|   PRIMARY KEY (`id`, `a`)
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_aig_pk_shrink_member_drop_fccc9f7c`.`t` DROP PRIMARY KEY, ADD PRIMARY KEY (`id`)
< ok 0 0
> exec
| ALTER TABLE `pa_aig_pk_shrink_member_drop_fccc9f7c`.`t` DROP COLUMN `a`
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_pk_shrink_member_drop_fccc9f7c.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_genpart_base_widened_ebfcd6dc_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_genpart_base_widened_ebfcd6dc_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_genpart_base_widened_ebfcd6dc_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, z BIGINT, seq INT NOT NULL AUTO_INCREMENT, g INT GENERATED ALWAYS AS (z + 1) VIRTUAL, UNIQUE KEY uk_sg (seq, g))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_genpart_base_widened_ebfcd6dc_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `z` bigint(20) DEFAULT NULL,
|   `seq` int(11) NOT NULL AUTO_INCREMENT,
|   `g` int(11) GENERATED ALWAYS AS ((`z` + 1)) VIRTUAL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_sg` (`seq`, `g`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_genpart_base_widened_ebfcd6dc
< ok 0 0
> exec
| CREATE DATABASE pa_aig_genpart_base_widened_ebfcd6dc DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_genpart_base_widened_ebfcd6dc
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, z INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_genpart_base_widened_ebfcd6dc.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `z` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_aig_genpart_base_widened_ebfcd6dc`.`t` ADD COLUMN `seq` int NOT NULL AUTO_INCREMENT
< ok 0 0
> exec
| ALTER TABLE `pa_aig_genpart_base_widened_ebfcd6dc`.`t` MODIFY COLUMN `z` bigint DEFAULT NULL
< ok 0 0
> exec
| ALTER TABLE `pa_aig_genpart_base_widened_ebfcd6dc`.`t` ADD COLUMN `g` int GENERATED ALWAYS AS ((`z` + 1)) VIRTUAL
< ok 0 0
> exec
| ALTER TABLE `pa_aig_genpart_base_widened_ebfcd6dc`.`t` ADD UNIQUE KEY `uk_sg` (`seq`,`g`) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_genpart_base_widened_ebfcd6dc.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `z` bigint(20) DEFAULT NULL,
|   `seq` int(11) NOT NULL AUTO_INCREMENT,
|   `g` int(11) GENERATED ALWAYS AS ((`z` + 1)) VIRTUAL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_sg` (`seq`, `g`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_widen_ai_key_replaced_34c77b19_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_widen_ai_key_replaced_34c77b19_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_widen_ai_key_replaced_34c77b19_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c BIGINT NOT NULL AUTO_INCREMENT, a INT NOT NULL, UNIQUE KEY uk_c (c, a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_widen_ai_key_replaced_34c77b19_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` bigint(20) NOT NULL AUTO_INCREMENT,
|   `a` int(11) NOT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_c` (`c`, `a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_widen_ai_key_replaced_34c77b19
< ok 0 0
> exec
| CREATE DATABASE pa_aig_widen_ai_key_replaced_34c77b19 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_widen_ai_key_replaced_34c77b19
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, UNIQUE KEY uk_c (c))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_widen_ai_key_replaced_34c77b19.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` int(11) NOT NULL AUTO_INCREMENT,
|   `a` int(11) NOT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_c` (`c`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_aig_widen_ai_key_replaced_34c77b19`.`t` DROP INDEX `uk_c`
< ok 0 0
> exec
| ALTER TABLE `pa_aig_widen_ai_key_replaced_34c77b19`.`t` MODIFY COLUMN `c` bigint NOT NULL AUTO_INCREMENT
< ok 0 0
> exec
| ALTER TABLE `pa_aig_widen_ai_key_replaced_34c77b19`.`t` ADD UNIQUE KEY `uk_c` (`c`,`a`) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_widen_ai_key_replaced_34c77b19.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` bigint(20) NOT NULL AUTO_INCREMENT,
|   `a` int(11) NOT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_c` (`c`, `a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_ai_migrate_to_earlier_name_e5fbc263_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_ai_migrate_to_earlier_name_e5fbc263_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_ai_migrate_to_earlier_name_e5fbc263_to
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL AUTO_INCREMENT, b INT NOT NULL, UNIQUE KEY uk_a (a), UNIQUE KEY uk_b (b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_ai_migrate_to_earlier_name_e5fbc263_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) NOT NULL AUTO_INCREMENT,
|   `b` int(11) NOT NULL,
|   UNIQUE KEY `uk_a` (`a`) BLOCK_SIZE 16384 LOCAL,
|   UNIQUE KEY `uk_b` (`b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_ai_migrate_to_earlier_name_e5fbc263
< ok 0 0
> exec
| CREATE DATABASE pa_aig_ai_migrate_to_earlier_name_e5fbc263 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_ai_migrate_to_earlier_name_e5fbc263
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_a (a), UNIQUE KEY uk_b (b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_ai_migrate_to_earlier_name_e5fbc263.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) NOT NULL,
|   `b` int(11) NOT NULL AUTO_INCREMENT,
|   UNIQUE KEY `uk_a` (`a`) BLOCK_SIZE 16384 LOCAL,
|   UNIQUE KEY `uk_b` (`b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_aig_ai_migrate_to_earlier_name_e5fbc263`.`t` MODIFY COLUMN `b` int NOT NULL
< ok 0 0
> exec
| ALTER TABLE `pa_aig_ai_migrate_to_earlier_name_e5fbc263`.`t` MODIFY COLUMN `a` int NOT NULL AUTO_INCREMENT
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_ai_migrate_to_earlier_name_e5fbc263.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) NOT NULL AUTO_INCREMENT,
|   `b` int(11) NOT NULL,
|   UNIQUE KEY `uk_a` (`a`) BLOCK_SIZE 16384 LOCAL,
|   UNIQUE KEY `uk_b` (`b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_pk_ai_member_hands_on_and_drops_a8cd695d_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_pk_ai_member_hands_on_and_drops_a8cd695d_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_pk_ai_member_hands_on_and_drops_a8cd695d_to
< ok 0 0
> exec
| CREATE TABLE t (b INT NOT NULL AUTO_INCREMENT, PRIMARY KEY (b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_pk_ai_member_hands_on_and_drops_a8cd695d_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`b`)
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_pk_ai_member_hands_on_and_drops_a8cd695d
< ok 0 0
> exec
| CREATE DATABASE pa_aig_pk_ai_member_hands_on_and_drops_a8cd695d DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_pk_ai_member_hands_on_and_drops_a8cd695d
< ok 0 0
> exec
| CREATE TABLE t (c INT NOT NULL AUTO_INCREMENT, b INT NOT NULL, PRIMARY KEY (c))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_pk_ai_member_hands_on_and_drops_a8cd695d.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `c` int(11) NOT NULL AUTO_INCREMENT,
|   `b` int(11) NOT NULL,
|   PRIMARY KEY (`c`)
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_aig_pk_ai_member_hands_on_and_drops_a8cd695d`.`t` MODIFY COLUMN `c` int NOT NULL
< ok 0 0
> exec
| ALTER TABLE `pa_aig_pk_ai_member_hands_on_and_drops_a8cd695d`.`t` MODIFY COLUMN `b` int NOT NULL AUTO_INCREMENT
< ok 0 0
> exec
| ALTER TABLE `pa_aig_pk_ai_member_hands_on_and_drops_a8cd695d`.`t` DROP PRIMARY KEY, ADD PRIMARY KEY (`b`)
< ok 0 0
> exec
| ALTER TABLE `pa_aig_pk_ai_member_hands_on_and_drops_a8cd695d`.`t` DROP COLUMN `c`
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_pk_ai_member_hands_on_and_drops_a8cd695d.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`b`)
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_pk_ai_member_hands_on_and_nulls_a69cb05_to
< ok 0 0
> exec
| CREATE DATABASE pa_aig_pk_ai_member_hands_on_and_nulls_a69cb05_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_pk_ai_member_hands_on_and_nulls_a69cb05_to
< ok 0 0
> exec
| CREATE TABLE t (c INT, b INT NOT NULL AUTO_INCREMENT, PRIMARY KEY (b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_pk_ai_member_hands_on_and_nulls_a69cb05_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `c` int(11) DEFAULT NULL,
|   `b` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`b`)
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_aig_pk_ai_member_hands_on_and_nulls_a69cb05
< ok 0 0
> exec
| CREATE DATABASE pa_aig_pk_ai_member_hands_on_and_nulls_a69cb05 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_aig_pk_ai_member_hands_on_and_nulls_a69cb05
< ok 0 0
> exec
| CREATE TABLE t (c INT NOT NULL AUTO_INCREMENT, b INT NOT NULL, PRIMARY KEY (c))
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_pk_ai_member_hands_on_and_nulls_a69cb05.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `c` int(11) NOT NULL AUTO_INCREMENT,
|   `b` int(11) NOT NULL,
|   PRIMARY KEY (`c`)
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_aig_pk_ai_member_hands_on_and_nulls_a69cb05`.`t` MODIFY COLUMN `c` int NOT NULL
< ok 0 0
> exec
| ALTER TABLE `pa_aig_pk_ai_member_hands_on_and_nulls_a69cb05`.`t` MODIFY COLUMN `b` int NOT NULL AUTO_INCREMENT
< ok 0 0
> exec
| ALTER TABLE `pa_aig_pk_ai_member_hands_on_and_nulls_a69cb05`.`t` DROP PRIMARY KEY, ADD PRIMARY KEY (`b`)
< ok 0 0
> exec
| ALTER TABLE `pa_aig_pk_ai_member_hands_on_and_nulls_a69cb05`.`t` MODIFY COLUMN `c` int DEFAULT NULL
< ok 0 0
> query
| SHOW CREATE TABLE pa_aig_pk_ai_member_hands_on_and_nulls_a69cb05.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `c` int(11) DEFAULT NULL,
|   `b` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`b`)
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
