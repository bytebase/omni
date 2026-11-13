-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracle_BigintUnsignedExactValue default#1

== TestOracle_BigintUnsignedExactValue/8.0/int64-max default#1
> exec
| DROP DATABASE IF EXISTS pa_bigexact_int64_max_320fd92_from
< ok 0 0
> exec
| CREATE DATABASE pa_bigexact_int64_max_320fd92_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_bigexact_int64_max_320fd92_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, big BIGINT NOT NULL DEFAULT 0)
< ok 0 0
> query
| SHOW CREATE TABLE pa_bigexact_int64_max_320fd92_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `big` bigint NOT NULL DEFAULT '0',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_bigexact_int64_max_320fd92_to
< ok 0 0
> exec
| CREATE DATABASE pa_bigexact_int64_max_320fd92_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_bigexact_int64_max_320fd92_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, big BIGINT NOT NULL DEFAULT 9223372036854775807)
< ok 0 0
> query
| SHOW CREATE TABLE pa_bigexact_int64_max_320fd92_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `big` bigint NOT NULL DEFAULT '9223372036854775807',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_bigexact_int64_max_320fd92
< ok 0 0
> exec
| CREATE DATABASE pa_bigexact_int64_max_320fd92 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_bigexact_int64_max_320fd92
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, big BIGINT NOT NULL DEFAULT 0)
< ok 0 0
> exec
| ALTER TABLE `pa_bigexact_int64_max_320fd92`.`t` MODIFY COLUMN `big` bigint NOT NULL DEFAULT '9223372036854775807'
< ok 0 0
> query
| SELECT COALESCE(COLUMN_DEFAULT, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=? AND TABLE_NAME='t' AND COLUMN_NAME='big'
: string
| pa_bigexact_int64_max_320fd92
< result
< column "COALESCE(COLUMN_DEFAULT, '')" "TEXT"
< row
: bytes
| 9223372036854775807

== TestOracle_BigintUnsignedExactValue/8.0/int64-max-plus-1 default#1
> exec
| DROP DATABASE IF EXISTS pa_bigexact_int64_max_plus_1_1852c303_from
< ok 0 0
> exec
| CREATE DATABASE pa_bigexact_int64_max_plus_1_1852c303_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_bigexact_int64_max_plus_1_1852c303_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, big BIGINT UNSIGNED NOT NULL DEFAULT 0)
< ok 0 0
> query
| SHOW CREATE TABLE pa_bigexact_int64_max_plus_1_1852c303_from.t
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
| DROP DATABASE IF EXISTS pa_bigexact_int64_max_plus_1_1852c303_to
< ok 0 0
> exec
| CREATE DATABASE pa_bigexact_int64_max_plus_1_1852c303_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_bigexact_int64_max_plus_1_1852c303_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, big BIGINT UNSIGNED NOT NULL DEFAULT 9223372036854775808)
< ok 0 0
> query
| SHOW CREATE TABLE pa_bigexact_int64_max_plus_1_1852c303_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `big` bigint unsigned NOT NULL DEFAULT '9223372036854775808',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_bigexact_int64_max_plus_1_1852c303
< ok 0 0
> exec
| CREATE DATABASE pa_bigexact_int64_max_plus_1_1852c303 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_bigexact_int64_max_plus_1_1852c303
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, big BIGINT UNSIGNED NOT NULL DEFAULT 0)
< ok 0 0
> exec
| ALTER TABLE `pa_bigexact_int64_max_plus_1_1852c303`.`t` MODIFY COLUMN `big` bigint unsigned NOT NULL DEFAULT '9223372036854775808'
< ok 0 0
> query
| SELECT COALESCE(COLUMN_DEFAULT, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=? AND TABLE_NAME='t' AND COLUMN_NAME='big'
: string
| pa_bigexact_int64_max_plus_1_1852c303
< result
< column "COALESCE(COLUMN_DEFAULT, '')" "TEXT"
< row
: bytes
| 9223372036854775808

== TestOracle_BigintUnsignedExactValue/8.0/int64-min default#1
> exec
| DROP DATABASE IF EXISTS pa_bigexact_int64_min_ed0d5338_from
< ok 0 0
> exec
| CREATE DATABASE pa_bigexact_int64_min_ed0d5338_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_bigexact_int64_min_ed0d5338_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, big BIGINT NOT NULL DEFAULT 0)
< ok 0 0
> query
| SHOW CREATE TABLE pa_bigexact_int64_min_ed0d5338_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `big` bigint NOT NULL DEFAULT '0',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_bigexact_int64_min_ed0d5338_to
< ok 0 0
> exec
| CREATE DATABASE pa_bigexact_int64_min_ed0d5338_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_bigexact_int64_min_ed0d5338_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, big BIGINT NOT NULL DEFAULT -9223372036854775808)
< ok 0 0
> query
| SHOW CREATE TABLE pa_bigexact_int64_min_ed0d5338_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `big` bigint NOT NULL DEFAULT '-9223372036854775808',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_bigexact_int64_min_ed0d5338
< ok 0 0
> exec
| CREATE DATABASE pa_bigexact_int64_min_ed0d5338 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_bigexact_int64_min_ed0d5338
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, big BIGINT NOT NULL DEFAULT 0)
< ok 0 0
> exec
| ALTER TABLE `pa_bigexact_int64_min_ed0d5338`.`t` MODIFY COLUMN `big` bigint NOT NULL DEFAULT '-9223372036854775808'
< ok 0 0
> query
| SELECT COALESCE(COLUMN_DEFAULT, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=? AND TABLE_NAME='t' AND COLUMN_NAME='big'
: string
| pa_bigexact_int64_min_ed0d5338
< result
< column "COALESCE(COLUMN_DEFAULT, '')" "TEXT"
< row
: bytes
| -9223372036854775808

== TestOracle_BigintUnsignedExactValue/8.0/uint64-max default#1
> exec
| DROP DATABASE IF EXISTS pa_bigexact_uint64_max_b7e5ca1b_from
< ok 0 0
> exec
| CREATE DATABASE pa_bigexact_uint64_max_b7e5ca1b_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_bigexact_uint64_max_b7e5ca1b_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, big BIGINT UNSIGNED NOT NULL DEFAULT 0)
< ok 0 0
> query
| SHOW CREATE TABLE pa_bigexact_uint64_max_b7e5ca1b_from.t
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
| DROP DATABASE IF EXISTS pa_bigexact_uint64_max_b7e5ca1b_to
< ok 0 0
> exec
| CREATE DATABASE pa_bigexact_uint64_max_b7e5ca1b_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_bigexact_uint64_max_b7e5ca1b_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, big BIGINT UNSIGNED NOT NULL DEFAULT 18446744073709551615)
< ok 0 0
> query
| SHOW CREATE TABLE pa_bigexact_uint64_max_b7e5ca1b_to.t
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
| DROP DATABASE IF EXISTS pa_bigexact_uint64_max_b7e5ca1b
< ok 0 0
> exec
| CREATE DATABASE pa_bigexact_uint64_max_b7e5ca1b DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_bigexact_uint64_max_b7e5ca1b
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, big BIGINT UNSIGNED NOT NULL DEFAULT 0)
< ok 0 0
> exec
| ALTER TABLE `pa_bigexact_uint64_max_b7e5ca1b`.`t` MODIFY COLUMN `big` bigint unsigned NOT NULL DEFAULT '18446744073709551615'
< ok 0 0
> query
| SELECT COALESCE(COLUMN_DEFAULT, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=? AND TABLE_NAME='t' AND COLUMN_NAME='big'
: string
| pa_bigexact_uint64_max_b7e5ca1b
< result
< column "COALESCE(COLUMN_DEFAULT, '')" "TEXT"
< row
: bytes
| 18446744073709551615
