-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_AutoIncGroupingIdempotence default#1
> exec
| DROP DATABASE IF EXISTS aig_idem_aig_add_col_uk
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_add_col_uk
< ok 1 0
> exec
| USE aig_idem_aig_add_col_uk
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT, seq BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_seq (seq))
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_add_col_uk.t
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
| DROP DATABASE IF EXISTS aig_idem_aig_add_col_plain_key
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_add_col_plain_key
< ok 1 0
> exec
| USE aig_idem_aig_add_col_plain_key
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT, seq BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, KEY k_seq (seq))
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_add_col_plain_key.t
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
| DROP DATABASE IF EXISTS aig_idem_aig_add_col_new_pk
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_add_col_new_pk
< ok 1 0
> exec
| USE aig_idem_aig_add_col_new_pk
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, a INT, seq BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, PRIMARY KEY (seq))
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_add_col_new_pk.t
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
| DROP DATABASE IF EXISTS aig_idem_aig_add_col_promote_pk
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_add_col_promote_pk
< ok 1 0
> exec
| USE aig_idem_aig_add_col_promote_pk
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT, seq BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, PRIMARY KEY (seq))
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_add_col_promote_pk.t
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
| DROP DATABASE IF EXISTS aig_idem_aig_add_col_promote_pk_drop_old
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_add_col_promote_pk_drop_old
< ok 1 0
> exec
| USE aig_idem_aig_add_col_promote_pk_drop_old
< ok 0 0
> exec
| CREATE TABLE t (b INT, seq BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, PRIMARY KEY (seq))
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_add_col_promote_pk_drop_old.t
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
| DROP DATABASE IF EXISTS aig_idem_aig_modify_gain_ai_new_key
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_modify_gain_ai_new_key
< ok 1 0
> exec
| USE aig_idem_aig_modify_gain_ai_new_key
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_c (c))
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_modify_gain_ai_new_key.t
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
| DROP DATABASE IF EXISTS aig_idem_aig_modify_gain_ai_existing_key
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_modify_gain_ai_existing_key
< ok 1 0
> exec
| USE aig_idem_aig_modify_gain_ai_existing_key
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_c (c))
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_modify_gain_ai_existing_key.t
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
| DROP DATABASE IF EXISTS aig_idem_aig_add_col_composite_backing
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_add_col_composite_backing
< ok 1 0
> exec
| USE aig_idem_aig_add_col_composite_backing
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, seq BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, d INT NOT NULL, UNIQUE KEY uk_sd (seq, d))
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_add_col_composite_backing.t
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
| DROP DATABASE IF EXISTS aig_idem_aig_ai_migrate_columns
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_ai_migrate_columns
< ok 1 0
> exec
| USE aig_idem_aig_ai_migrate_columns
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_a (a), UNIQUE KEY uk_b (b))
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_ai_migrate_columns.t
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
| DROP DATABASE IF EXISTS aig_idem_aig_ai_migrate_to_new_col
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_ai_migrate_to_new_col
< ok 1 0
> exec
| USE aig_idem_aig_ai_migrate_to_new_col
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_a (a), UNIQUE KEY uk_b (b))
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_ai_migrate_to_new_col.t
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
| DROP DATABASE IF EXISTS aig_idem_aig_add_col_with_other_changes
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_add_col_with_other_changes
< ok 1 0
> exec
| USE aig_idem_aig_add_col_with_other_changes
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT, e VARCHAR(20), f INT, seq BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_seq (seq))
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_add_col_with_other_changes.t
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
| DROP DATABASE IF EXISTS aig_idem_aig_myisam_add_nonleftmost
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_myisam_add_nonleftmost
< ok 1 0
> exec
| USE aig_idem_aig_myisam_add_nonleftmost
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, seq INT NOT NULL AUTO_INCREMENT, KEY k_as (a, seq)) ENGINE=MyISAM
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_myisam_add_nonleftmost.t
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
| DROP DATABASE IF EXISTS aig_idem_aig_myisam_gain_ai_existing_nonleftmost
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_myisam_gain_ai_existing_nonleftmost
< ok 1 0
> exec
| USE aig_idem_aig_myisam_gain_ai_existing_nonleftmost
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, c INT NOT NULL AUTO_INCREMENT, KEY k (a, c)) ENGINE=MyISAM
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_myisam_gain_ai_existing_nonleftmost.t
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
| DROP DATABASE IF EXISTS aig_idem_aig_myisam_shared_key_ai_migrate
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_myisam_shared_key_ai_migrate
< ok 1 0
> exec
| USE aig_idem_aig_myisam_shared_key_ai_migrate
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c2 INT NOT NULL AUTO_INCREMENT, KEY k2 (c2)) ENGINE=MyISAM
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_myisam_shared_key_ai_migrate.t
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
| DROP DATABASE IF EXISTS aig_idem_aig_myisam_shared_key_ai_migrate_survive
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_myisam_shared_key_ai_migrate_survive
< ok 1 0
> exec
| USE aig_idem_aig_myisam_shared_key_ai_migrate_survive
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL AUTO_INCREMENT, KEY k2 (b)) ENGINE=MyISAM
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_myisam_shared_key_ai_migrate_survive.t
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
| DROP DATABASE IF EXISTS aig_idem_aig_add_col_desc_key
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_add_col_desc_key
< ok 1 0
> exec
| USE aig_idem_aig_add_col_desc_key
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, seq INT NOT NULL AUTO_INCREMENT, KEY k_seq (seq DESC))
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_add_col_desc_key.t
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
| DROP DATABASE IF EXISTS aig_idem_aig_add_col_invisible_key
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_add_col_invisible_key
< ok 1 0
> exec
| USE aig_idem_aig_add_col_invisible_key
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, seq INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_seq (seq) INVISIBLE)
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_add_col_invisible_key.t
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
| DROP DATABASE IF EXISTS aig_idem_aig_add_col_trailing_expr_key
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_add_col_trailing_expr_key
< ok 1 0
> exec
| USE aig_idem_aig_add_col_trailing_expr_key
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT, seq INT NOT NULL AUTO_INCREMENT, KEY k_se (seq, (a + 1)))
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_add_col_trailing_expr_key.t
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
| DROP DATABASE IF EXISTS aig_idem_aig_drop_key_and_column
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_drop_key_and_column
< ok 1 0
> exec
| USE aig_idem_aig_drop_key_and_column
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_drop_key_and_column.t
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
| DROP DATABASE IF EXISTS aig_idem_aig_drop_key_deai
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_drop_key_deai
< ok 1 0
> exec
| USE aig_idem_aig_drop_key_deai
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_drop_key_deai.t
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
| DROP DATABASE IF EXISTS aig_idem_aig_replace_backing_key
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_replace_backing_key
< ok 1 0
> exec
| USE aig_idem_aig_replace_backing_key
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, UNIQUE KEY uk_c (c, a))
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_replace_backing_key.t
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
| DROP DATABASE IF EXISTS aig_idem_aig_swap_backing_key
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_swap_backing_key
< ok 1 0
> exec
| USE aig_idem_aig_swap_backing_key
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, KEY k_c (c))
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_swap_backing_key.t
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
| DROP DATABASE IF EXISTS aig_idem_aig_multi_drop_plus_add
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_multi_drop_plus_add
< ok 1 0
> exec
| USE aig_idem_aig_multi_drop_plus_add
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk3 (c))
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_multi_drop_plus_add.t
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
| DROP DATABASE IF EXISTS aig_idem_aig_drop_covered_elsewhere
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_drop_covered_elsewhere
< ok 1 0
> exec
| USE aig_idem_aig_drop_covered_elsewhere
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, KEY k2 (c))
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_drop_covered_elsewhere.t
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
| DROP DATABASE IF EXISTS aig_idem_aig_drop_pk_and_column
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_drop_pk_and_column
< ok 1 0
> exec
| USE aig_idem_aig_drop_pk_and_column
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_drop_pk_and_column.t
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
| DROP DATABASE IF EXISTS aig_idem_aig_drop_pk_deai
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_drop_pk_deai
< ok 1 0
> exec
| USE aig_idem_aig_drop_pk_deai
< ok 0 0
> exec
| CREATE TABLE t (c INT NOT NULL, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_drop_pk_deai.t
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
| DROP DATABASE IF EXISTS aig_idem_aig_drop_composite_pk_and_column
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_drop_composite_pk_and_column
< ok 1 0
> exec
| USE aig_idem_aig_drop_composite_pk_and_column
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_drop_composite_pk_and_column.t
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
| DROP DATABASE IF EXISTS aig_idem_aig_pk_swap_drop_secondary
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_pk_swap_drop_secondary
< ok 1 0
> exec
| USE aig_idem_aig_pk_swap_drop_secondary
< ok 0 0
> exec
| CREATE TABLE t (x INT NOT NULL, c INT NOT NULL AUTO_INCREMENT, PRIMARY KEY (c))
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_pk_swap_drop_secondary.t
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
| DROP DATABASE IF EXISTS aig_idem_aig_pk_relocate_ai_to_secondary
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_pk_relocate_ai_to_secondary
< ok 1 0
> exec
| USE aig_idem_aig_pk_relocate_ai_to_secondary
< ok 0 0
> exec
| CREATE TABLE t (c INT NOT NULL AUTO_INCREMENT, b INT NOT NULL, x INT, PRIMARY KEY (b), KEY zz_c (c))
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_pk_relocate_ai_to_secondary.t
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
| DROP DATABASE IF EXISTS aig_idem_aig_pk_keep_ai_demote_member
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_pk_keep_ai_demote_member
< ok 1 0
> exec
| USE aig_idem_aig_pk_keep_ai_demote_member
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, a INT NULL, PRIMARY KEY (id))
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_pk_keep_ai_demote_member.t
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
| DROP DATABASE IF EXISTS aig_idem_aig_generated_keypart_pullin
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_generated_keypart_pullin
< ok 1 0
> exec
| USE aig_idem_aig_generated_keypart_pullin
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, seq INT NOT NULL AUTO_INCREMENT, d INT NOT NULL, g INT GENERATED ALWAYS AS (d+1) VIRTUAL, UNIQUE KEY uk_sg (seq, g))
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_generated_keypart_pullin.t
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
| DROP DATABASE IF EXISTS aig_idem_aig_funcpart_new_col
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_funcpart_new_col
< ok 1 0
> exec
| USE aig_idem_aig_funcpart_new_col
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, seq INT NOT NULL AUTO_INCREMENT, z INT, KEY k_sz (seq, (z + 1)))
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_funcpart_new_col.t
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
| DROP DATABASE IF EXISTS aig_idem_aig_prefix_over_widened
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_prefix_over_widened
< ok 1 0
> exec
| USE aig_idem_aig_prefix_over_widened
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, v VARCHAR(500), seq INT NOT NULL AUTO_INCREMENT, KEY k_sv (seq, v(300)))
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_prefix_over_widened.t
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
| DROP DATABASE IF EXISTS aig_idem_aig_pk_drop_member_drop_newkey
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_pk_drop_member_drop_newkey
< ok 1 0
> exec
| USE aig_idem_aig_pk_drop_member_drop_newkey
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, KEY k_id (id))
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_pk_drop_member_drop_newkey.t
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
| DROP DATABASE IF EXISTS aig_idem_aig_pk_shrink_member_drop
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_pk_shrink_member_drop
< ok 1 0
> exec
| USE aig_idem_aig_pk_shrink_member_drop
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, PRIMARY KEY (id))
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_pk_shrink_member_drop.t
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
| DROP DATABASE IF EXISTS aig_idem_aig_genpart_base_widened
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_genpart_base_widened
< ok 1 0
> exec
| USE aig_idem_aig_genpart_base_widened
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, z BIGINT, seq INT NOT NULL AUTO_INCREMENT, g INT GENERATED ALWAYS AS (z + 1) VIRTUAL, UNIQUE KEY uk_sg (seq, g))
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_genpart_base_widened.t
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
| DROP DATABASE IF EXISTS aig_idem_aig_widen_ai_key_replaced
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_widen_ai_key_replaced
< ok 1 0
> exec
| USE aig_idem_aig_widen_ai_key_replaced
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c BIGINT NOT NULL AUTO_INCREMENT, a INT NOT NULL, UNIQUE KEY uk_c (c, a))
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_widen_ai_key_replaced.t
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
| DROP DATABASE IF EXISTS aig_idem_aig_ai_migrate_to_earlier_name
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_ai_migrate_to_earlier_name
< ok 1 0
> exec
| USE aig_idem_aig_ai_migrate_to_earlier_name
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL AUTO_INCREMENT, b INT NOT NULL, UNIQUE KEY uk_a (a), UNIQUE KEY uk_b (b))
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_ai_migrate_to_earlier_name.t
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
| DROP DATABASE IF EXISTS aig_idem_aig_pk_ai_member_hands_on_and_drops
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_pk_ai_member_hands_on_and_drops
< ok 1 0
> exec
| USE aig_idem_aig_pk_ai_member_hands_on_and_drops
< ok 0 0
> exec
| CREATE TABLE t (b INT NOT NULL AUTO_INCREMENT, PRIMARY KEY (b))
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_pk_ai_member_hands_on_and_drops.t
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
| DROP DATABASE IF EXISTS aig_idem_aig_pk_ai_member_hands_on_and_nulls
< ok 0 0
> exec
| CREATE DATABASE aig_idem_aig_pk_ai_member_hands_on_and_nulls
< ok 1 0
> exec
| USE aig_idem_aig_pk_ai_member_hands_on_and_nulls
< ok 0 0
> exec
| CREATE TABLE t (c INT, b INT NOT NULL AUTO_INCREMENT, PRIMARY KEY (b))
< ok 0 0
> query
| SHOW CREATE TABLE aig_idem_aig_pk_ai_member_hands_on_and_nulls.t
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
