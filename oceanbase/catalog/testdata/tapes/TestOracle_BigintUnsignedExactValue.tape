-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_BigintUnsignedExactValue default#1
> exec
| DROP DATABASE IF EXISTS pa_bigexact_uint64_max_103f6768_from
< ok 0 0
> exec
| CREATE DATABASE pa_bigexact_uint64_max_103f6768_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_bigexact_uint64_max_103f6768_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, big BIGINT UNSIGNED NOT NULL DEFAULT 0)
< ok 0 0
> query
| SHOW CREATE TABLE pa_bigexact_uint64_max_103f6768_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `big` bigint(20) unsigned NOT NULL DEFAULT '0',
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_bigexact_uint64_max_103f6768_to
< ok 0 0
> exec
| CREATE DATABASE pa_bigexact_uint64_max_103f6768_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_bigexact_uint64_max_103f6768_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, big BIGINT UNSIGNED NOT NULL DEFAULT 18446744073709551615)
< ok 0 0
> query
| SHOW CREATE TABLE pa_bigexact_uint64_max_103f6768_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `big` bigint(20) unsigned NOT NULL DEFAULT '18446744073709551615',
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_bigexact_uint64_max_103f6768
< ok 0 0
> exec
| CREATE DATABASE pa_bigexact_uint64_max_103f6768 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_bigexact_uint64_max_103f6768
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, big BIGINT UNSIGNED NOT NULL DEFAULT 0)
< ok 0 0
> exec
| ALTER TABLE `pa_bigexact_uint64_max_103f6768`.`t` MODIFY COLUMN `big` bigint unsigned NOT NULL DEFAULT '18446744073709551615'
< ok 0 0
> query
| SELECT COALESCE(COLUMN_DEFAULT, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=? AND TABLE_NAME='t' AND COLUMN_NAME='big'
: string
| pa_bigexact_uint64_max_103f6768
< result
< column "COALESCE(COLUMN_DEFAULT, '')" "TEXT"
< row
: bytes
| 18446744073709551615
> exec
| DROP DATABASE IF EXISTS pa_bigexact_int64_max_plus_1_229f49dc_from
< ok 0 0
> exec
| CREATE DATABASE pa_bigexact_int64_max_plus_1_229f49dc_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_bigexact_int64_max_plus_1_229f49dc_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, big BIGINT UNSIGNED NOT NULL DEFAULT 0)
< ok 0 0
> query
| SHOW CREATE TABLE pa_bigexact_int64_max_plus_1_229f49dc_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `big` bigint(20) unsigned NOT NULL DEFAULT '0',
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_bigexact_int64_max_plus_1_229f49dc_to
< ok 0 0
> exec
| CREATE DATABASE pa_bigexact_int64_max_plus_1_229f49dc_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_bigexact_int64_max_plus_1_229f49dc_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, big BIGINT UNSIGNED NOT NULL DEFAULT 9223372036854775808)
< ok 0 0
> query
| SHOW CREATE TABLE pa_bigexact_int64_max_plus_1_229f49dc_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `big` bigint(20) unsigned NOT NULL DEFAULT '9223372036854775808',
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_bigexact_int64_max_plus_1_229f49dc
< ok 0 0
> exec
| CREATE DATABASE pa_bigexact_int64_max_plus_1_229f49dc DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_bigexact_int64_max_plus_1_229f49dc
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, big BIGINT UNSIGNED NOT NULL DEFAULT 0)
< ok 0 0
> exec
| ALTER TABLE `pa_bigexact_int64_max_plus_1_229f49dc`.`t` MODIFY COLUMN `big` bigint unsigned NOT NULL DEFAULT '9223372036854775808'
< ok 0 0
> query
| SELECT COALESCE(COLUMN_DEFAULT, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=? AND TABLE_NAME='t' AND COLUMN_NAME='big'
: string
| pa_bigexact_int64_max_plus_1_229f49dc
< result
< column "COALESCE(COLUMN_DEFAULT, '')" "TEXT"
< row
: bytes
| 9223372036854775808
> exec
| DROP DATABASE IF EXISTS pa_bigexact_int64_max_87d7069f_from
< ok 0 0
> exec
| CREATE DATABASE pa_bigexact_int64_max_87d7069f_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_bigexact_int64_max_87d7069f_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, big BIGINT NOT NULL DEFAULT 0)
< ok 0 0
> query
| SHOW CREATE TABLE pa_bigexact_int64_max_87d7069f_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `big` bigint(20) NOT NULL DEFAULT '0',
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_bigexact_int64_max_87d7069f_to
< ok 0 0
> exec
| CREATE DATABASE pa_bigexact_int64_max_87d7069f_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_bigexact_int64_max_87d7069f_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, big BIGINT NOT NULL DEFAULT 9223372036854775807)
< ok 0 0
> query
| SHOW CREATE TABLE pa_bigexact_int64_max_87d7069f_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `big` bigint(20) NOT NULL DEFAULT '9223372036854775807',
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_bigexact_int64_max_87d7069f
< ok 0 0
> exec
| CREATE DATABASE pa_bigexact_int64_max_87d7069f DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_bigexact_int64_max_87d7069f
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, big BIGINT NOT NULL DEFAULT 0)
< ok 0 0
> exec
| ALTER TABLE `pa_bigexact_int64_max_87d7069f`.`t` MODIFY COLUMN `big` bigint NOT NULL DEFAULT '9223372036854775807'
< ok 0 0
> query
| SELECT COALESCE(COLUMN_DEFAULT, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=? AND TABLE_NAME='t' AND COLUMN_NAME='big'
: string
| pa_bigexact_int64_max_87d7069f
< result
< column "COALESCE(COLUMN_DEFAULT, '')" "TEXT"
< row
: bytes
| 9223372036854775807
> exec
| DROP DATABASE IF EXISTS pa_bigexact_int64_min_75c36291_from
< ok 0 0
> exec
| CREATE DATABASE pa_bigexact_int64_min_75c36291_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_bigexact_int64_min_75c36291_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, big BIGINT NOT NULL DEFAULT 0)
< ok 0 0
> query
| SHOW CREATE TABLE pa_bigexact_int64_min_75c36291_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `big` bigint(20) NOT NULL DEFAULT '0',
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_bigexact_int64_min_75c36291_to
< ok 0 0
> exec
| CREATE DATABASE pa_bigexact_int64_min_75c36291_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_bigexact_int64_min_75c36291_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, big BIGINT NOT NULL DEFAULT -9223372036854775808)
< ok 0 0
> query
| SHOW CREATE TABLE pa_bigexact_int64_min_75c36291_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `big` bigint(20) NOT NULL DEFAULT '-9223372036854775808',
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_bigexact_int64_min_75c36291
< ok 0 0
> exec
| CREATE DATABASE pa_bigexact_int64_min_75c36291 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_bigexact_int64_min_75c36291
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, big BIGINT NOT NULL DEFAULT 0)
< ok 0 0
> exec
| ALTER TABLE `pa_bigexact_int64_min_75c36291`.`t` MODIFY COLUMN `big` bigint NOT NULL DEFAULT '-9223372036854775808'
< ok 0 0
> query
| SELECT COALESCE(COLUMN_DEFAULT, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=? AND TABLE_NAME='t' AND COLUMN_NAME='big'
: string
| pa_bigexact_int64_min_75c36291
< result
< column "COALESCE(COLUMN_DEFAULT, '')" "TEXT"
< row
: bytes
| -9223372036854775808
