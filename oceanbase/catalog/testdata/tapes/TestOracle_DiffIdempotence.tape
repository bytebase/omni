-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_DiffIdempotence default#1
> exec
| DROP DATABASE IF EXISTS diff_int_widths
< ok 0 0
> exec
| CREATE DATABASE diff_int_widths
< ok 1 0
> exec
| USE diff_int_widths
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT(11), b BIGINT(20), c INT, p TINYINT, h INT UNSIGNED)
< ok 0 0
> query
| SHOW CREATE TABLE diff_int_widths.t
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
|   `b` bigint(20) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL,
|   `p` tinyint(4) DEFAULT NULL,
|   `h` int(10) unsigned DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS diff_boolean
< ok 0 0
> exec
| CREATE DATABASE diff_boolean
< ok 1 0
> exec
| USE diff_boolean
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, m BOOLEAN, n BOOL, d TINYINT(1))
< ok 0 0
> query
| SHOW CREATE TABLE diff_boolean.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `m` tinyint(1) DEFAULT NULL,
|   `n` tinyint(1) DEFAULT NULL,
|   `d` tinyint(1) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS diff_numeric
< ok 0 0
> exec
| CREATE DATABASE diff_numeric
< ok 1 0
> exec
| USE diff_numeric
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a DECIMAL, b DECIMAL(10,2), d NUMERIC, j REAL, f FLOAT)
< ok 0 0
> query
| SHOW CREATE TABLE diff_numeric.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` decimal(10,0) DEFAULT NULL,
|   `b` decimal(10,2) DEFAULT NULL,
|   `d` decimal(10,0) DEFAULT NULL,
|   `j` double DEFAULT NULL,
|   `f` float DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS diff_char_bit_year
< ok 0 0
> exec
| CREATE DATABASE diff_char_bit_year
< ok 1 0
> exec
| USE diff_char_bit_year
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a CHAR, d BINARY, c BIT, y YEAR)
< ok 0 0
> query
| SHOW CREATE TABLE diff_char_bit_year.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` char(1) DEFAULT NULL,
|   `d` binary(1) DEFAULT NULL,
|   `c` bit(1) DEFAULT NULL,
|   `y` year(4) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS diff_defaults
< ok 0 0
> exec
| CREATE DATABASE diff_defaults
< ok 1 0
> exec
| USE diff_defaults
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT DEFAULT 0, b INT DEFAULT '7', f DECIMAL(10,2) DEFAULT 0, c VARCHAR(20) DEFAULT 'x', j BOOLEAN DEFAULT TRUE)
< ok 0 0
> query
| SHOW CREATE TABLE diff_defaults.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT '0',
|   `b` int(11) DEFAULT '7',
|   `f` decimal(10,2) DEFAULT '0.00',
|   `c` varchar(20) DEFAULT 'x',
|   `j` tinyint(1) DEFAULT '1',
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS diff_bigint_unsigned_max_default
< ok 0 0
> exec
| CREATE DATABASE diff_bigint_unsigned_max_default
< ok 1 0
> exec
| USE diff_bigint_unsigned_max_default
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, big BIGINT UNSIGNED NOT NULL DEFAULT 18446744073709551615, mid BIGINT UNSIGNED NOT NULL DEFAULT 9223372036854775808)
< ok 0 0
> query
| SHOW CREATE TABLE diff_bigint_unsigned_max_default.t
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
|   `mid` bigint(20) unsigned NOT NULL DEFAULT '9223372036854775808',
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS diff_bigint_signed_boundary_default
< ok 0 0
> exec
| CREATE DATABASE diff_bigint_signed_boundary_default
< ok 1 0
> exec
| USE diff_bigint_signed_boundary_default
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, lo BIGINT NOT NULL DEFAULT -9223372036854775808, hi BIGINT NOT NULL DEFAULT 9223372036854775807)
< ok 0 0
> query
| SHOW CREATE TABLE diff_bigint_signed_boundary_default.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `lo` bigint(20) NOT NULL DEFAULT '-9223372036854775808',
|   `hi` bigint(20) NOT NULL DEFAULT '9223372036854775807',
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS diff_bit_default_forms
< ok 0 0
> exec
| CREATE DATABASE diff_bit_default_forms
< ok 1 0
> exec
| USE diff_bit_default_forms
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a BIT(8) NOT NULL DEFAULT b'0', b BIT(8) NOT NULL DEFAULT 0x05, c BIT(8) NOT NULL DEFAULT 5, d BIT(16) NOT NULL DEFAULT 'AB')
< ok 0 0
> query
| SHOW CREATE TABLE diff_bit_default_forms.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` bit(8) NOT NULL DEFAULT b'0',
|   `b` bit(8) NOT NULL DEFAULT b'101',
|   `c` bit(8) NOT NULL DEFAULT b'101',
|   `d` bit(16) NOT NULL DEFAULT b'100000101000010',
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS diff_year_default
< ok 0 0
> exec
| CREATE DATABASE diff_year_default
< ok 1 0
> exec
| USE diff_year_default
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a YEAR NOT NULL DEFAULT 2000, b YEAR NOT NULL DEFAULT 1999, c YEAR NOT NULL DEFAULT 99, d YEAR NOT NULL DEFAULT 0)
< ok 0 0
> query
| SHOW CREATE TABLE diff_year_default.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` year(4) NOT NULL DEFAULT '2000',
|   `b` year(4) NOT NULL DEFAULT '1999',
|   `c` year(4) NOT NULL DEFAULT '1999',
|   `d` year(4) NOT NULL DEFAULT '0000',
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS diff_nullability
< ok 0 0
> exec
| CREATE DATABASE diff_nullability
< ok 1 0
> exec
| USE diff_nullability
< ok 0 0
> exec
| CREATE TABLE t (id INT, a INT NULL, b INT DEFAULT NULL, d VARCHAR(10) NULL, PRIMARY KEY (id))
< ok 0 0
> query
| SHOW CREATE TABLE diff_nullability.t
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
|   `b` int(11) DEFAULT NULL,
|   `d` varchar(10) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS diff_enum_set
< ok 0 0
> exec
| CREATE DATABASE diff_enum_set
< ok 1 0
> exec
| USE diff_enum_set
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a ENUM('x','y','z'), b SET('a','b','c'), e ENUM("dq1","dq2"))
< ok 0 0
> query
| SHOW CREATE TABLE diff_enum_set.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` enum('x','y','z') DEFAULT NULL,
|   `b` set('a','b','c') DEFAULT NULL,
|   `e` enum('dq1','dq2') DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS diff_generated
< ok 0 0
> exec
| CREATE DATABASE diff_generated
< ok 1 0
> exec
| USE diff_generated
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT GENERATED ALWAYS AS (a+1) VIRTUAL, c INT GENERATED ALWAYS AS ( a * 2 ) STORED, d INT AS (a+1))
< ok 0 0
> query
| SHOW CREATE TABLE diff_generated.t
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
|   `b` int(11) GENERATED ALWAYS AS ((`a` + 1)) VIRTUAL,
|   `c` int(11) GENERATED ALWAYS AS ((`a` * 2)) STORED,
|   `d` int(11) GENERATED ALWAYS AS ((`a` + 1)) VIRTUAL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS diff_generated_case
< ok 0 0
> exec
| CREATE DATABASE diff_generated_case
< ok 1 0
> exec
| USE diff_generated_case
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, s VARCHAR(10), b INT GENERATED ALWAYS AS (IF(a > 0, a, 0)) VIRTUAL, c INT GENERATED ALWAYS AS (CASE WHEN a > 0 THEN a ELSE 0 END) STORED, d INT AS (CASE a WHEN 1 THEN 10 ELSE 20 END), e INT AS (a IN (1, 2)), f INT AS (a IS NULL))
< ok 0 0
> query
| SHOW CREATE TABLE diff_generated_case.t
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
|   `s` varchar(10) DEFAULT NULL,
|   `b` int(11) GENERATED ALWAYS AS ((case when (`a` > 0) then `a` else 0 end)) VIRTUAL,
|   `c` int(11) GENERATED ALWAYS AS ((case when (`a` > 0) then `a` else 0 end)) STORED,
|   `d` int(11) GENERATED ALWAYS AS ((case `a` when 1 then 10 else 20 end)) VIRTUAL,
|   `e` int(11) GENERATED ALWAYS AS ((`a` in (1,2))) VIRTUAL,
|   `f` int(11) GENERATED ALWAYS AS ((`a` is null)) VIRTUAL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS diff_implicit_table_opts
< ok 0 0
> exec
| CREATE DATABASE diff_implicit_table_opts
< ok 1 0
> exec
| USE diff_implicit_table_opts
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a VARCHAR(10))
< ok 0 0
> query
| SHOW CREATE TABLE diff_implicit_table_opts.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` varchar(10) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS diff_charset_echo
< ok 0 0
> exec
| CREATE DATABASE diff_charset_echo
< ok 1 0
> exec
| USE diff_charset_echo
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a VARCHAR(10) CHARACTER SET utf8mb4, b VARCHAR(10)) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE diff_charset_echo.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` varchar(10) DEFAULT NULL,
|   `b` varchar(10) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS diff_binary_modifier_utf8mb4
< ok 0 0
> exec
| CREATE DATABASE diff_binary_modifier_utf8mb4
< ok 1 0
> exec
| USE diff_binary_modifier_utf8mb4
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, k VARCHAR(128) BINARY NOT NULL, c CHAR(5) BINARY, tx TEXT BINARY, ti TINYTEXT BINARY) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE diff_binary_modifier_utf8mb4.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `k` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
|   `c` char(5) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   `tx` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   `ti` tinytext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS diff_binary_modifier_before_charset
< ok 0 0
> exec
| CREATE DATABASE diff_binary_modifier_before_charset
< ok 1 0
> exec
| USE diff_binary_modifier_before_charset
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, ac VARCHAR(10) BINARY CHARACTER SET ascii) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE diff_binary_modifier_before_charset.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `ac` varchar(10) CHARACTER SET ascii COLLATE ascii_bin DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS diff_binary_modifier_after_charset
< ok 0 0
> exec
| CREATE DATABASE diff_binary_modifier_after_charset
< ok 1 0
> exec
| USE diff_binary_modifier_after_charset
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, ac VARCHAR(10) CHARACTER SET ascii BINARY) DEFAULT CHARSET=utf8mb4
< error
| 1064 42000 …near 'BINARY) DEFAULT CHARSET=utf8mb4' at line 1
> exec
| DROP DATABASE IF EXISTS diff_binary_modifier_latin1
< ok 0 0
> exec
| CREATE DATABASE diff_binary_modifier_latin1
< ok 1 0
> exec
| USE diff_binary_modifier_latin1
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a VARCHAR(50) BINARY) DEFAULT CHARSET=latin1
< ok 0 0
> query
| SHOW CREATE TABLE diff_binary_modifier_latin1.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` varchar(50) CHARACTER SET latin1 COLLATE latin1_bin DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = latin1 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS diff_comment
< ok 0 0
> exec
| CREATE DATABASE diff_comment
< ok 1 0
> exec
| USE diff_comment
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT COMMENT 'has ''quote'' inside') COMMENT='tbl ''c'''
< ok 0 0
> query
| SHOW CREATE TABLE diff_comment.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL COMMENT 'has \'quote\' inside',
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 COMMENT = 'tbl \'c\''
> exec
| DROP DATABASE IF EXISTS diff_autoinc
< ok 0 0
> exec
| CREATE DATABASE diff_autoinc
< ok 1 0
> exec
| USE diff_autoinc
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT PRIMARY KEY, a INT) AUTO_INCREMENT=100
< ok 0 0
> query
| SHOW CREATE TABLE diff_autoinc.t
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
| ) ORGANIZATION INDEX AUTO_INCREMENT = 100 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS diff_timestamp
< ok 0 0
> exec
| CREATE DATABASE diff_timestamp
< ok 1 0
> exec
| USE diff_timestamp
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a TIMESTAMP, b TIMESTAMP NULL, d TIMESTAMP DEFAULT CURRENT_TIMESTAMP)
< ok 0 0
> query
| SHOW CREATE TABLE diff_timestamp.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` timestamp NULL DEFAULT NULL,
|   `b` timestamp NULL DEFAULT NULL,
|   `d` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS diff_datetime_default
< ok 0 0
> exec
| CREATE DATABASE diff_datetime_default
< ok 1 0
> exec
| USE diff_datetime_default
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g DATETIME DEFAULT CURRENT_TIMESTAMP, h DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP)
< ok 0 0
> query
| SHOW CREATE TABLE diff_datetime_default.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `g` datetime DEFAULT CURRENT_TIMESTAMP,
|   `h` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS diff_functional_default
< ok 0 0
> exec
| CREATE DATABASE diff_functional_default
< ok 1 0
> exec
| USE diff_functional_default
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT DEFAULT (1+1), b VARCHAR(36) DEFAULT (UUID()))
< ok 0 0
> query
| SHOW CREATE TABLE diff_functional_default.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT ((1 + 1)),
|   `b` varchar(36) DEFAULT (UUID()),
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS diff_gen_introducer
< ok 0 0
> exec
| CREATE DATABASE diff_gen_introducer
< ok 1 0
> exec
| USE diff_gen_introducer
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a VARCHAR(20), e VARCHAR(20) GENERATED ALWAYS AS (CONCAT(a,'x')) VIRTUAL) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE diff_gen_introducer.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` varchar(20) DEFAULT NULL,
|   `e` varchar(20) GENERATED ALWAYS AS (CONCAT(`a`,'x')) VIRTUAL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS diff_multi_column
< ok 0 0
> exec
| CREATE DATABASE diff_multi_column
< ok 1 0
> exec
| USE diff_multi_column
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, c1 VARCHAR(10), c2 INT DEFAULT 0, c3 DATE, c4 DECIMAL(8,3), PRIMARY KEY (id))
< ok 0 0
> query
| SHOW CREATE TABLE diff_multi_column.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c1` varchar(10) DEFAULT NULL,
|   `c2` int(11) DEFAULT '0',
|   `c3` date DEFAULT NULL,
|   `c4` decimal(8,3) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS diff_row_format_explicit
< ok 0 0
> exec
| CREATE DATABASE diff_row_format_explicit
< ok 1 0
> exec
| USE diff_row_format_explicit
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=8
< ok 0 0
> query
| SHOW CREATE TABLE diff_row_format_explicit.t
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
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = COMPRESSED COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS diff_adjacent_literal_default
< ok 0 0
> exec
| CREATE DATABASE diff_adjacent_literal_default
< ok 1 0
> exec
| USE diff_adjacent_literal_default
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a VARCHAR(20) DEFAULT 'x' 'y', b VARCHAR(20) DEFAULT 'x' 'y' 'z')
< ok 0 0
> query
| SHOW CREATE TABLE diff_adjacent_literal_default.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` varchar(20) DEFAULT 'xy',
|   `b` varchar(20) DEFAULT 'xyz',
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS diff_adjacent_literal_default_double_quoted
< ok 0 0
> exec
| CREATE DATABASE diff_adjacent_literal_default_double_quoted
< ok 1 0
> exec
| USE diff_adjacent_literal_default_double_quoted
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(20) DEFAULT "x" 'y')
< error
| 1064 42000 …near 'FROM DUAL' at line 1
> exec
| DROP DATABASE IF EXISTS diff_adjacent_literal_partition
< ok 0 0
> exec
| CREATE DATABASE diff_adjacent_literal_partition
< ok 1 0
> exec
| USE diff_adjacent_literal_partition
< ok 0 0
> exec
| CREATE TABLE t (c VARCHAR(10) NOT NULL PRIMARY KEY) PARTITION BY LIST COLUMNS(c) (PARTITION p0 VALUES IN ('a' 'b'), PARTITION p1 VALUES IN ('z'))
< ok 0 0
> query
| SHOW CREATE TABLE diff_adjacent_literal_partition.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `c` varchar(10) NOT NULL,
|   PRIMARY KEY (`c`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by list columns(`c`)
| (partition `p0` values in ('ab'),
| partition `p1` values in ('z'))
> exec
| DROP DATABASE IF EXISTS diff_adjacent_literal_default_introducer
< ok 0 0
> exec
| CREATE DATABASE diff_adjacent_literal_default_introducer
< ok 1 0
> exec
| USE diff_adjacent_literal_default_introducer
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a VARCHAR(20) CHARACTER SET utf8mb4 DEFAULT _utf8mb4'x' 'y', b VARCHAR(20) CHARACTER SET utf8mb4 DEFAULT _utf8mb4'z') DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE diff_adjacent_literal_default_introducer.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` varchar(20) DEFAULT 'xy',
|   `b` varchar(20) DEFAULT 'z',
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
