-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracle_MigrationIdempotence default#1
> exec
| DROP DATABASE IF EXISTS gen_idem_int_widths
< ok 0 0
> exec
| CREATE DATABASE gen_idem_int_widths
< ok 1 0
> exec
| USE gen_idem_int_widths
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT(11), b BIGINT(20), c INT, p TINYINT, h INT UNSIGNED)
< ok 0 0
> query
| SHOW CREATE TABLE gen_idem_int_widths.t
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
|   `c` int DEFAULT NULL,
|   `p` tinyint DEFAULT NULL,
|   `h` int unsigned DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS gen_idem_boolean
< ok 0 0
> exec
| CREATE DATABASE gen_idem_boolean
< ok 1 0
> exec
| USE gen_idem_boolean
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, m BOOLEAN, n BOOL, d TINYINT(1))
< ok 0 0
> query
| SHOW CREATE TABLE gen_idem_boolean.t
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
|   `n` tinyint(1) DEFAULT NULL,
|   `d` tinyint(1) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS gen_idem_numeric
< ok 0 0
> exec
| CREATE DATABASE gen_idem_numeric
< ok 1 0
> exec
| USE gen_idem_numeric
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a DECIMAL, b DECIMAL(10,2), d NUMERIC, j REAL, f FLOAT)
< ok 0 0
> query
| SHOW CREATE TABLE gen_idem_numeric.t
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
|   `d` decimal(10,0) DEFAULT NULL,
|   `j` double DEFAULT NULL,
|   `f` float DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS gen_idem_char_bit_year
< ok 0 0
> exec
| CREATE DATABASE gen_idem_char_bit_year
< ok 1 0
> exec
| USE gen_idem_char_bit_year
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a CHAR, d BINARY, c BIT, y YEAR)
< ok 0 0
> query
| SHOW CREATE TABLE gen_idem_char_bit_year.t
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
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS gen_idem_defaults
< ok 0 0
> exec
| CREATE DATABASE gen_idem_defaults
< ok 1 0
> exec
| USE gen_idem_defaults
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT DEFAULT 0, b INT DEFAULT '7', f DECIMAL(10,2) DEFAULT 0, c VARCHAR(20) DEFAULT 'x', j BOOLEAN DEFAULT TRUE)
< ok 0 0
> query
| SHOW CREATE TABLE gen_idem_defaults.t
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
| DROP DATABASE IF EXISTS gen_idem_bigint_unsigned_max_default
< ok 0 0
> exec
| CREATE DATABASE gen_idem_bigint_unsigned_max_default
< ok 1 0
> exec
| USE gen_idem_bigint_unsigned_max_default
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, big BIGINT UNSIGNED NOT NULL DEFAULT 18446744073709551615, mid BIGINT UNSIGNED NOT NULL DEFAULT 9223372036854775808)
< ok 0 0
> query
| SHOW CREATE TABLE gen_idem_bigint_unsigned_max_default.t
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
|   `mid` bigint unsigned NOT NULL DEFAULT '9223372036854775808',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS gen_idem_bigint_signed_boundary_default
< ok 0 0
> exec
| CREATE DATABASE gen_idem_bigint_signed_boundary_default
< ok 1 0
> exec
| USE gen_idem_bigint_signed_boundary_default
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, lo BIGINT NOT NULL DEFAULT -9223372036854775808, hi BIGINT NOT NULL DEFAULT 9223372036854775807)
< ok 0 0
> query
| SHOW CREATE TABLE gen_idem_bigint_signed_boundary_default.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `lo` bigint NOT NULL DEFAULT '-9223372036854775808',
|   `hi` bigint NOT NULL DEFAULT '9223372036854775807',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS gen_idem_bit_default_forms
< ok 0 0
> exec
| CREATE DATABASE gen_idem_bit_default_forms
< ok 1 0
> exec
| USE gen_idem_bit_default_forms
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a BIT(8) NOT NULL DEFAULT b'0', b BIT(8) NOT NULL DEFAULT 0x05, c BIT(8) NOT NULL DEFAULT 5, d BIT(16) NOT NULL DEFAULT 'AB')
< ok 0 0
> query
| SHOW CREATE TABLE gen_idem_bit_default_forms.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` bit(8) NOT NULL DEFAULT b'0',
|   `b` bit(8) NOT NULL DEFAULT b'101',
|   `c` bit(8) NOT NULL DEFAULT b'101',
|   `d` bit(16) NOT NULL DEFAULT b'100000101000010',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS gen_idem_year_default
< ok 0 0
> exec
| CREATE DATABASE gen_idem_year_default
< ok 1 0
> exec
| USE gen_idem_year_default
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a YEAR NOT NULL DEFAULT 2000, b YEAR NOT NULL DEFAULT 1999, c YEAR NOT NULL DEFAULT 99, d YEAR NOT NULL DEFAULT 0)
< ok 0 0
> query
| SHOW CREATE TABLE gen_idem_year_default.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` year NOT NULL DEFAULT '2000',
|   `b` year NOT NULL DEFAULT '1999',
|   `c` year NOT NULL DEFAULT '1999',
|   `d` year NOT NULL DEFAULT '0000',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS gen_idem_nullability
< ok 0 0
> exec
| CREATE DATABASE gen_idem_nullability
< ok 1 0
> exec
| USE gen_idem_nullability
< ok 0 0
> exec
| CREATE TABLE t (id INT, a INT NULL, b INT DEFAULT NULL, d VARCHAR(10) NULL, PRIMARY KEY (id))
< ok 0 0
> query
| SHOW CREATE TABLE gen_idem_nullability.t
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
|   `d` varchar(10) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS gen_idem_enum_set
< ok 0 0
> exec
| CREATE DATABASE gen_idem_enum_set
< ok 1 0
> exec
| USE gen_idem_enum_set
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a ENUM('x','y','z'), b SET('a','b','c'), e ENUM("dq1","dq2"))
< ok 0 0
> query
| SHOW CREATE TABLE gen_idem_enum_set.t
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
|   `e` enum('dq1','dq2') DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS gen_idem_generated
< ok 0 0
> exec
| CREATE DATABASE gen_idem_generated
< ok 1 0
> exec
| USE gen_idem_generated
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT GENERATED ALWAYS AS (a+1) VIRTUAL, c INT GENERATED ALWAYS AS ( a * 2 ) STORED, d INT AS (a+1))
< ok 0 0
> query
| SHOW CREATE TABLE gen_idem_generated.t
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
|   `d` int GENERATED ALWAYS AS ((`a` + 1)) VIRTUAL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS gen_idem_implicit_table_opts
< ok 0 0
> exec
| CREATE DATABASE gen_idem_implicit_table_opts
< ok 1 0
> exec
| USE gen_idem_implicit_table_opts
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a VARCHAR(10))
< ok 0 0
> query
| SHOW CREATE TABLE gen_idem_implicit_table_opts.t
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
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS gen_idem_charset_echo
< ok 0 0
> exec
| CREATE DATABASE gen_idem_charset_echo
< ok 1 0
> exec
| USE gen_idem_charset_echo
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a VARCHAR(10) CHARACTER SET utf8mb4, b VARCHAR(10)) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE gen_idem_charset_echo.t
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
| DROP DATABASE IF EXISTS gen_idem_binary_modifier_utf8mb4
< ok 0 0
> exec
| CREATE DATABASE gen_idem_binary_modifier_utf8mb4
< ok 1 0
> exec
| USE gen_idem_binary_modifier_utf8mb4
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, k VARCHAR(128) BINARY NOT NULL, c CHAR(5) BINARY, tx TEXT BINARY, ti TINYTEXT BINARY, ac VARCHAR(10) CHARACTER SET ascii BINARY) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE gen_idem_binary_modifier_utf8mb4.t
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
|   `c` char(5) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   `tx` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin,
|   `ti` tinytext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin,
|   `ac` varchar(10) CHARACTER SET ascii COLLATE ascii_bin DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS gen_idem_binary_modifier_latin1
< ok 0 0
> exec
| CREATE DATABASE gen_idem_binary_modifier_latin1
< ok 1 0
> exec
| USE gen_idem_binary_modifier_latin1
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a VARCHAR(50) BINARY) DEFAULT CHARSET=latin1
< ok 0 0
> query
| SHOW CREATE TABLE gen_idem_binary_modifier_latin1.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` varchar(50) CHARACTER SET latin1 COLLATE latin1_bin DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=latin1
> exec
| DROP DATABASE IF EXISTS gen_idem_comment
< ok 0 0
> exec
| CREATE DATABASE gen_idem_comment
< ok 1 0
> exec
| USE gen_idem_comment
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT COMMENT 'has ''quote'' inside') COMMENT='tbl ''c'''
< ok 0 0
> query
| SHOW CREATE TABLE gen_idem_comment.t
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
| DROP DATABASE IF EXISTS gen_idem_autoinc
< ok 0 0
> exec
| CREATE DATABASE gen_idem_autoinc
< ok 1 0
> exec
| USE gen_idem_autoinc
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT PRIMARY KEY, a INT) AUTO_INCREMENT=100
< ok 0 0
> query
| SHOW CREATE TABLE gen_idem_autoinc.t
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
| DROP DATABASE IF EXISTS gen_idem_timestamp
< ok 0 0
> exec
| CREATE DATABASE gen_idem_timestamp
< ok 1 0
> exec
| USE gen_idem_timestamp
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a TIMESTAMP, b TIMESTAMP NULL, d TIMESTAMP DEFAULT CURRENT_TIMESTAMP)
< ok 0 0
> query
| SHOW CREATE TABLE gen_idem_timestamp.t
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
| DROP DATABASE IF EXISTS gen_idem_datetime_default
< ok 0 0
> exec
| CREATE DATABASE gen_idem_datetime_default
< ok 1 0
> exec
| USE gen_idem_datetime_default
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g DATETIME DEFAULT CURRENT_TIMESTAMP, h DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP)
< ok 0 0
> query
| SHOW CREATE TABLE gen_idem_datetime_default.t
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
| DROP DATABASE IF EXISTS gen_idem_functional_default
< ok 0 0
> exec
| CREATE DATABASE gen_idem_functional_default
< ok 1 0
> exec
| USE gen_idem_functional_default
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT DEFAULT (1+1), b VARCHAR(36) DEFAULT (UUID()))
< ok 0 0
> query
| SHOW CREATE TABLE gen_idem_functional_default.t
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
| DROP DATABASE IF EXISTS gen_idem_gen_introducer
< ok 0 0
> exec
| CREATE DATABASE gen_idem_gen_introducer
< ok 1 0
> exec
| USE gen_idem_gen_introducer
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a VARCHAR(20), e VARCHAR(20) GENERATED ALWAYS AS (CONCAT(a,'x')) VIRTUAL) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE gen_idem_gen_introducer.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` varchar(20) DEFAULT NULL,
|   `e` varchar(20) GENERATED ALWAYS AS (concat(`a`,_utf8mb4'x')) VIRTUAL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS gen_idem_multi_column
< ok 0 0
> exec
| CREATE DATABASE gen_idem_multi_column
< ok 1 0
> exec
| USE gen_idem_multi_column
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, c1 VARCHAR(10), c2 INT DEFAULT 0, c3 DATE, c4 DECIMAL(8,3), PRIMARY KEY (id))
< ok 0 0
> query
| SHOW CREATE TABLE gen_idem_multi_column.t
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
| DROP DATABASE IF EXISTS gen_idem_row_format_explicit
< ok 0 0
> exec
| CREATE DATABASE gen_idem_row_format_explicit
< ok 1 0
> exec
| USE gen_idem_row_format_explicit
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=8
< ok 0 0
> query
| SHOW CREATE TABLE gen_idem_row_format_explicit.t
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
| DROP DATABASE IF EXISTS gen_idem_adjacent_literal_default
< ok 0 0
> exec
| CREATE DATABASE gen_idem_adjacent_literal_default
< ok 1 0
> exec
| USE gen_idem_adjacent_literal_default
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a VARCHAR(20) DEFAULT 'x' 'y', b VARCHAR(20) DEFAULT 'x' 'y' 'z', c VARCHAR(20) DEFAULT "x" 'y')
< ok 0 0
> query
| SHOW CREATE TABLE gen_idem_adjacent_literal_default.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` varchar(20) DEFAULT 'xy',
|   `b` varchar(20) DEFAULT 'xyz',
|   `c` varchar(20) DEFAULT 'xy',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS gen_idem_adjacent_literal_partition
< ok 0 0
> exec
| CREATE DATABASE gen_idem_adjacent_literal_partition
< ok 1 0
> exec
| USE gen_idem_adjacent_literal_partition
< ok 0 0
> exec
| CREATE TABLE t (c VARCHAR(10) NOT NULL PRIMARY KEY) PARTITION BY LIST COLUMNS(c) (PARTITION p0 VALUES IN ('a' 'b'), PARTITION p1 VALUES IN ('z'))
< ok 0 0
> query
| SHOW CREATE TABLE gen_idem_adjacent_literal_partition.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `c` varchar(10) NOT NULL,
|   PRIMARY KEY (`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50500 PARTITION BY LIST  COLUMNS(c)
| (PARTITION p0 VALUES IN ('ab') ENGINE = InnoDB,
|  PARTITION p1 VALUES IN ('z') ENGINE = InnoDB) */
> exec
| DROP DATABASE IF EXISTS gen_idem_adjacent_literal_default_introducer
< ok 0 0
> exec
| CREATE DATABASE gen_idem_adjacent_literal_default_introducer
< ok 1 0
> exec
| USE gen_idem_adjacent_literal_default_introducer
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a VARCHAR(20) CHARACTER SET utf8mb4 DEFAULT _utf8mb4'x' 'y', b VARCHAR(20) CHARACTER SET utf8mb4 DEFAULT _utf8mb4'z') DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE gen_idem_adjacent_literal_default_introducer.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'xy',
|   `b` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'z',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
