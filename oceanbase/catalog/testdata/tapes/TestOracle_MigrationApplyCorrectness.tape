-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_MigrationApplyCorrectness default#1
> exec
| DROP DATABASE IF EXISTS pa_create_simple_72990483_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_simple_72990483_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_simple_72990483_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_simple_72990483_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_simple_72990483
< ok 0 0
> exec
| CREATE DATABASE pa_create_simple_72990483 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_simple_72990483
< ok 0 0
> exec
| CREATE TABLE `pa_create_simple_72990483`.`t` (
|   `id` int NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_simple_72990483.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_int_widths_15b15bdc_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_int_widths_15b15bdc_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_int_widths_15b15bdc_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT(11), b BIGINT(20), p TINYINT, h INT UNSIGNED, z INT ZEROFILL)
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_int_widths_15b15bdc_to.t
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
|   `p` tinyint(4) DEFAULT NULL,
|   `h` int(10) unsigned DEFAULT NULL,
|   `z` int(10) unsigned zerofill DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_int_widths_15b15bdc
< ok 0 0
> exec
| CREATE DATABASE pa_create_int_widths_15b15bdc DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_int_widths_15b15bdc
< ok 0 0
> exec
| CREATE TABLE `pa_create_int_widths_15b15bdc`.`t` (
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
| SHOW CREATE TABLE pa_create_int_widths_15b15bdc.t
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
|   `p` tinyint(4) DEFAULT NULL,
|   `h` int(10) unsigned DEFAULT NULL,
|   `z` int(10) unsigned zerofill DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_boolean_9beff95b_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_boolean_9beff95b_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_boolean_9beff95b_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, m BOOLEAN, d TINYINT(1))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_boolean_9beff95b_to.t
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
|   `d` tinyint(1) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_boolean_9beff95b
< ok 0 0
> exec
| CREATE DATABASE pa_create_boolean_9beff95b DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_boolean_9beff95b
< ok 0 0
> exec
| CREATE TABLE `pa_create_boolean_9beff95b`.`t` (
|   `id` int NOT NULL,
|   `m` tinyint(1) DEFAULT NULL,
|   `d` tinyint(1) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_boolean_9beff95b.t
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
|   `d` tinyint(1) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_numeric_5b66789c_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_numeric_5b66789c_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_numeric_5b66789c_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a DECIMAL, b DECIMAL(10,2), j REAL, f FLOAT, g DOUBLE(15,4))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_numeric_5b66789c_to.t
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
|   `j` double DEFAULT NULL,
|   `f` float DEFAULT NULL,
|   `g` double(15,4) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_numeric_5b66789c
< ok 0 0
> exec
| CREATE DATABASE pa_create_numeric_5b66789c DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_numeric_5b66789c
< ok 0 0
> exec
| CREATE TABLE `pa_create_numeric_5b66789c`.`t` (
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
| SHOW CREATE TABLE pa_create_numeric_5b66789c.t
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
|   `j` double DEFAULT NULL,
|   `f` float DEFAULT NULL,
|   `g` double(15,4) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_char_bit_year_b40f2ab3_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_char_bit_year_b40f2ab3_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_char_bit_year_b40f2ab3_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a CHAR, d BINARY, c BIT, y YEAR, v VARCHAR(40), bb VARBINARY(8))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_char_bit_year_b40f2ab3_to.t
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
|   `v` varchar(40) DEFAULT NULL,
|   `bb` varbinary(8) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_char_bit_year_b40f2ab3
< ok 0 0
> exec
| CREATE DATABASE pa_create_char_bit_year_b40f2ab3 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_char_bit_year_b40f2ab3
< ok 0 0
> exec
| CREATE TABLE `pa_create_char_bit_year_b40f2ab3`.`t` (
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
| SHOW CREATE TABLE pa_create_char_bit_year_b40f2ab3.t
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
|   `v` varchar(40) DEFAULT NULL,
|   `bb` varbinary(8) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_enum_set_ff5749b1_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_enum_set_ff5749b1_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_enum_set_ff5749b1_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a ENUM('x','y','z'), b SET('a','b','c'))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_enum_set_ff5749b1_to.t
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
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_enum_set_ff5749b1
< ok 0 0
> exec
| CREATE DATABASE pa_create_enum_set_ff5749b1 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_enum_set_ff5749b1
< ok 0 0
> exec
| CREATE TABLE `pa_create_enum_set_ff5749b1`.`t` (
|   `id` int NOT NULL,
|   `a` enum('x','y','z') DEFAULT NULL,
|   `b` set('a','b','c') DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_enum_set_ff5749b1.t
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
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_defaults_b9d0a193_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_defaults_b9d0a193_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_defaults_b9d0a193_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT DEFAULT 0, b INT DEFAULT '7', f DECIMAL(10,2) DEFAULT 0, c VARCHAR(20) DEFAULT 'x', j BOOLEAN DEFAULT TRUE)
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_defaults_b9d0a193_to.t
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
| DROP DATABASE IF EXISTS pa_create_defaults_b9d0a193
< ok 0 0
> exec
| CREATE DATABASE pa_create_defaults_b9d0a193 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_defaults_b9d0a193
< ok 0 0
> exec
| CREATE TABLE `pa_create_defaults_b9d0a193`.`t` (
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
| SHOW CREATE TABLE pa_create_defaults_b9d0a193.t
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
| DROP DATABASE IF EXISTS pa_create_nullability_c12f35f6_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_nullability_c12f35f6_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_nullability_c12f35f6_to
< ok 0 0
> exec
| CREATE TABLE t (id INT, a INT NULL, b INT DEFAULT NULL, d VARCHAR(10) NOT NULL, PRIMARY KEY (id))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_nullability_c12f35f6_to.t
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
|   `d` varchar(10) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_nullability_c12f35f6
< ok 0 0
> exec
| CREATE DATABASE pa_create_nullability_c12f35f6 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_nullability_c12f35f6
< ok 0 0
> exec
| CREATE TABLE `pa_create_nullability_c12f35f6`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   `d` varchar(10) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_nullability_c12f35f6.t
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
|   `d` varchar(10) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_comment_6e225e62_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_comment_6e225e62_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_comment_6e225e62_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT COMMENT 'has ''quote'' inside') COMMENT='tbl ''c'''
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_comment_6e225e62_to.t
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
| DROP DATABASE IF EXISTS pa_create_comment_6e225e62
< ok 0 0
> exec
| CREATE DATABASE pa_create_comment_6e225e62 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_comment_6e225e62
< ok 0 0
> exec
| CREATE TABLE `pa_create_comment_6e225e62`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL COMMENT 'has ''quote'' inside',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='tbl ''c'''
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_comment_6e225e62.t
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
| DROP DATABASE IF EXISTS pa_create_autoinc_1f928fa2_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_autoinc_1f928fa2_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_autoinc_1f928fa2_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT PRIMARY KEY, a INT) AUTO_INCREMENT=100
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_autoinc_1f928fa2_to.t
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
| DROP DATABASE IF EXISTS pa_create_autoinc_1f928fa2
< ok 0 0
> exec
| CREATE DATABASE pa_create_autoinc_1f928fa2 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_autoinc_1f928fa2
< ok 0 0
> exec
| CREATE TABLE `pa_create_autoinc_1f928fa2`.`t` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `a` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_autoinc_1f928fa2.t
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
| DROP DATABASE IF EXISTS pa_create_generated_54d6a56e_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_generated_54d6a56e_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_generated_54d6a56e_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT GENERATED ALWAYS AS (a+1) VIRTUAL, c INT GENERATED ALWAYS AS (a*2) STORED)
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_generated_54d6a56e_to.t
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
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_generated_54d6a56e
< ok 0 0
> exec
| CREATE DATABASE pa_create_generated_54d6a56e DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_generated_54d6a56e
< ok 0 0
> exec
| CREATE TABLE `pa_create_generated_54d6a56e`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   `b` int GENERATED ALWAYS AS (((`a` + 1))) VIRTUAL,
|   `c` int GENERATED ALWAYS AS (((`a` * 2))) STORED,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_generated_54d6a56e.t
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
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_timestamp_2efbb6bf_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_timestamp_2efbb6bf_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_timestamp_2efbb6bf_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a TIMESTAMP, b TIMESTAMP NULL, d TIMESTAMP DEFAULT CURRENT_TIMESTAMP)
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_timestamp_2efbb6bf_to.t
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
| DROP DATABASE IF EXISTS pa_create_timestamp_2efbb6bf
< ok 0 0
> exec
| CREATE DATABASE pa_create_timestamp_2efbb6bf DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_timestamp_2efbb6bf
< ok 0 0
> exec
| CREATE TABLE `pa_create_timestamp_2efbb6bf`.`t` (
|   `id` int NOT NULL,
|   `a` timestamp NULL DEFAULT NULL,
|   `b` timestamp NULL DEFAULT NULL,
|   `d` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_timestamp_2efbb6bf.t
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
| DROP DATABASE IF EXISTS pa_create_datetime_default_81d50ab4_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_datetime_default_81d50ab4_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_datetime_default_81d50ab4_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g DATETIME DEFAULT CURRENT_TIMESTAMP, h DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP)
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_datetime_default_81d50ab4_to.t
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
| DROP DATABASE IF EXISTS pa_create_datetime_default_81d50ab4
< ok 0 0
> exec
| CREATE DATABASE pa_create_datetime_default_81d50ab4 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_datetime_default_81d50ab4
< ok 0 0
> exec
| CREATE TABLE `pa_create_datetime_default_81d50ab4`.`t` (
|   `id` int NOT NULL,
|   `g` datetime DEFAULT CURRENT_TIMESTAMP,
|   `h` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_datetime_default_81d50ab4.t
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
| DROP DATABASE IF EXISTS pa_create_charset_explicit_6d19189e_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_charset_explicit_6d19189e_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_charset_explicit_6d19189e_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a VARCHAR(10) CHARACTER SET utf8mb4, b VARCHAR(10)) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_charset_explicit_6d19189e_to.t
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
| DROP DATABASE IF EXISTS pa_create_charset_explicit_6d19189e
< ok 0 0
> exec
| CREATE DATABASE pa_create_charset_explicit_6d19189e DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_charset_explicit_6d19189e
< ok 0 0
> exec
| CREATE TABLE `pa_create_charset_explicit_6d19189e`.`t` (
|   `id` int NOT NULL,
|   `a` varchar(10) DEFAULT NULL,
|   `b` varchar(10) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_charset_explicit_6d19189e.t
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
| DROP DATABASE IF EXISTS pa_create_engine_myisam_a87060fa_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_engine_myisam_a87060fa_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_engine_myisam_a87060fa_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT) ENGINE=MyISAM
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_engine_myisam_a87060fa_to.t
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
| DROP DATABASE IF EXISTS pa_create_engine_myisam_a87060fa
< ok 0 0
> exec
| CREATE DATABASE pa_create_engine_myisam_a87060fa DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_engine_myisam_a87060fa
< ok 0 0
> exec
| CREATE TABLE `pa_create_engine_myisam_a87060fa`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_engine_myisam_a87060fa.t
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
| DROP DATABASE IF EXISTS pa_create_multi_column_7f6c769b_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_multi_column_7f6c769b_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_multi_column_7f6c769b_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, c1 VARCHAR(10), c2 INT DEFAULT 0, c3 DATE, c4 DECIMAL(8,3), PRIMARY KEY (id))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_multi_column_7f6c769b_to.t
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
| DROP DATABASE IF EXISTS pa_create_multi_column_7f6c769b
< ok 0 0
> exec
| CREATE DATABASE pa_create_multi_column_7f6c769b DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_multi_column_7f6c769b
< ok 0 0
> exec
| CREATE TABLE `pa_create_multi_column_7f6c769b`.`t` (
|   `id` int NOT NULL,
|   `c1` varchar(10) DEFAULT NULL,
|   `c2` int DEFAULT '0',
|   `c3` date DEFAULT NULL,
|   `c4` decimal(8,3) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_multi_column_7f6c769b.t
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
| DROP DATABASE IF EXISTS pa_create_functional_default_b1005140_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_functional_default_b1005140_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_functional_default_b1005140_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT DEFAULT (1+1), b VARCHAR(36) DEFAULT (UUID()))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_functional_default_b1005140_to.t
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
| DROP DATABASE IF EXISTS pa_create_functional_default_b1005140
< ok 0 0
> exec
| CREATE DATABASE pa_create_functional_default_b1005140 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_functional_default_b1005140
< ok 0 0
> exec
| CREATE TABLE `pa_create_functional_default_b1005140`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT ((1 + 1)),
|   `b` varchar(36) DEFAULT (uuid()),
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_functional_default_b1005140.t
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
|   `b` varchar(36) DEFAULT (uuid()),
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_row_format_485c0745_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_row_format_485c0745_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_row_format_485c0745_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=8
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_row_format_485c0745_to.t
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
| DROP DATABASE IF EXISTS pa_create_row_format_485c0745
< ok 0 0
> exec
| CREATE DATABASE pa_create_row_format_485c0745 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_row_format_485c0745
< ok 0 0
> exec
| CREATE TABLE `pa_create_row_format_485c0745`.`t` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPRESSED
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_row_format_485c0745.t
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
| DROP DATABASE IF EXISTS pa_drop_table_f5d5b5c2
< ok 0 0
> exec
| CREATE DATABASE pa_drop_table_f5d5b5c2 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_table_f5d5b5c2
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_table_f5d5b5c2.t
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
| DROP TABLE `pa_drop_table_f5d5b5c2`.`t`
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_table_f5d5b5c2.t
< error
| 1146 42S02 Table 'pa_drop_table_f5d5b5c2.t' doesn't exist
> exec
| DROP DATABASE IF EXISTS pa_add_column_3f217ece_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_column_3f217ece_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_column_3f217ece_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(20))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_column_3f217ece_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_column_3f217ece
< ok 0 0
> exec
| CREATE DATABASE pa_add_column_3f217ece DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_column_3f217ece
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_column_3f217ece.t
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
| ALTER TABLE `pa_add_column_3f217ece`.`t` ADD COLUMN `name` varchar(20) DEFAULT NULL
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_column_3f217ece.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_column_notnull_default_390c49cb_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_column_notnull_default_390c49cb_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_column_notnull_default_390c49cb_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, n INT NOT NULL DEFAULT 5)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_column_notnull_default_390c49cb_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `n` int(11) NOT NULL DEFAULT '5',
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_column_notnull_default_390c49cb
< ok 0 0
> exec
| CREATE DATABASE pa_add_column_notnull_default_390c49cb DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_column_notnull_default_390c49cb
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_column_notnull_default_390c49cb.t
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
| ALTER TABLE `pa_add_column_notnull_default_390c49cb`.`t` ADD COLUMN `n` int NOT NULL DEFAULT '5'
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_column_notnull_default_390c49cb.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `n` int(11) NOT NULL DEFAULT '5',
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_generated_column_eb62683a_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_generated_column_eb62683a_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_generated_column_eb62683a_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT GENERATED ALWAYS AS (a+1) STORED)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_generated_column_eb62683a_to.t
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
|   `b` int(11) GENERATED ALWAYS AS ((`a` + 1)) STORED,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_generated_column_eb62683a
< ok 0 0
> exec
| CREATE DATABASE pa_add_generated_column_eb62683a DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_generated_column_eb62683a
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_generated_column_eb62683a.t
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
| ALTER TABLE `pa_add_generated_column_eb62683a`.`t` ADD COLUMN `b` int GENERATED ALWAYS AS (((`a` + 1))) STORED
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_generated_column_eb62683a.t
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
|   `b` int(11) GENERATED ALWAYS AS ((`a` + 1)) STORED,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_timestamp_column_9a349bc9_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_timestamp_column_9a349bc9_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_timestamp_column_9a349bc9_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, ts TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_timestamp_column_9a349bc9_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `ts` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_timestamp_column_9a349bc9
< ok 0 0
> exec
| CREATE DATABASE pa_add_timestamp_column_9a349bc9 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_timestamp_column_9a349bc9
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_timestamp_column_9a349bc9.t
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
| ALTER TABLE `pa_add_timestamp_column_9a349bc9`.`t` ADD COLUMN `ts` timestamp NULL DEFAULT CURRENT_TIMESTAMP
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_timestamp_column_9a349bc9.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `ts` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_drop_column_4bb3b7c4_to
< ok 0 0
> exec
| CREATE DATABASE pa_drop_column_4bb3b7c4_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_column_4bb3b7c4_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_column_4bb3b7c4_to.t
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
| DROP DATABASE IF EXISTS pa_drop_column_4bb3b7c4
< ok 0 0
> exec
| CREATE DATABASE pa_drop_column_4bb3b7c4 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_column_4bb3b7c4
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(20))
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_column_4bb3b7c4.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_drop_column_4bb3b7c4`.`t` DROP COLUMN `name`
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_column_4bb3b7c4.t
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
| DROP DATABASE IF EXISTS pa_modify_type_f46a5ff7_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_type_f46a5ff7_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_type_f46a5ff7_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v BIGINT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_type_f46a5ff7_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `v` bigint(20) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_type_f46a5ff7
< ok 0 0
> exec
| CREATE DATABASE pa_modify_type_f46a5ff7 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_type_f46a5ff7
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_type_f46a5ff7.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `v` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_modify_type_f46a5ff7`.`t` MODIFY COLUMN `v` bigint DEFAULT NULL
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_type_f46a5ff7.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `v` bigint(20) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_type_width_57_vs_80_965cfd1a_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_type_width_57_vs_80_965cfd1a_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_type_width_57_vs_80_965cfd1a_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v BIGINT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_type_width_57_vs_80_965cfd1a_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `v` bigint(20) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_type_width_57_vs_80_965cfd1a
< ok 0 0
> exec
| CREATE DATABASE pa_modify_type_width_57_vs_80_965cfd1a DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_type_width_57_vs_80_965cfd1a
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v SMALLINT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_type_width_57_vs_80_965cfd1a.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `v` smallint(6) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_modify_type_width_57_vs_80_965cfd1a`.`t` MODIFY COLUMN `v` bigint DEFAULT NULL
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_type_width_57_vs_80_965cfd1a.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `v` bigint(20) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_nullability_13e71f48_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_nullability_13e71f48_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_nullability_13e71f48_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v INT NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_nullability_13e71f48_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `v` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_nullability_13e71f48
< ok 0 0
> exec
| CREATE DATABASE pa_modify_nullability_13e71f48 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_nullability_13e71f48
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v INT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_nullability_13e71f48.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `v` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_modify_nullability_13e71f48`.`t` MODIFY COLUMN `v` int NOT NULL
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_nullability_13e71f48.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `v` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_default_4b16c9d8_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_default_4b16c9d8_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_default_4b16c9d8_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v INT DEFAULT 1)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_default_4b16c9d8_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `v` int(11) DEFAULT '1',
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_default_4b16c9d8
< ok 0 0
> exec
| CREATE DATABASE pa_modify_default_4b16c9d8 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_default_4b16c9d8
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v INT DEFAULT 0)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_default_4b16c9d8.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `v` int(11) DEFAULT '0',
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_modify_default_4b16c9d8`.`t` MODIFY COLUMN `v` int DEFAULT '1'
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_default_4b16c9d8.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `v` int(11) DEFAULT '1',
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_default_bigint_unsigned_max_18cd4f07_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_default_bigint_unsigned_max_18cd4f07_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_default_bigint_unsigned_max_18cd4f07_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, big BIGINT UNSIGNED NOT NULL DEFAULT 18446744073709551615)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_default_bigint_unsigned_max_18cd4f07_to.t
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
| DROP DATABASE IF EXISTS pa_modify_default_bigint_unsigned_max_18cd4f07
< ok 0 0
> exec
| CREATE DATABASE pa_modify_default_bigint_unsigned_max_18cd4f07 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_default_bigint_unsigned_max_18cd4f07
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, big BIGINT UNSIGNED NOT NULL DEFAULT 0)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_default_bigint_unsigned_max_18cd4f07.t
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
| ALTER TABLE `pa_modify_default_bigint_unsigned_max_18cd4f07`.`t` MODIFY COLUMN `big` bigint unsigned NOT NULL DEFAULT '18446744073709551615'
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_default_bigint_unsigned_max_18cd4f07.t
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
| DROP DATABASE IF EXISTS pa_modify_default_bit_hex_bcf9526a_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_default_bit_hex_bcf9526a_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_default_bit_hex_bcf9526a_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, f BIT(8) NOT NULL DEFAULT 0x05)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_default_bit_hex_bcf9526a_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `f` bit(8) NOT NULL DEFAULT b'101',
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_default_bit_hex_bcf9526a
< ok 0 0
> exec
| CREATE DATABASE pa_modify_default_bit_hex_bcf9526a DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_default_bit_hex_bcf9526a
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, f BIT(8) NOT NULL DEFAULT b'0')
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_default_bit_hex_bcf9526a.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `f` bit(8) NOT NULL DEFAULT b'0',
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_modify_default_bit_hex_bcf9526a`.`t` MODIFY COLUMN `f` bit(8) NOT NULL DEFAULT b'101'
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_default_bit_hex_bcf9526a.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `f` bit(8) NOT NULL DEFAULT b'101',
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_default_year_e8d89aba_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_default_year_e8d89aba_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_default_year_e8d89aba_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, y YEAR NOT NULL DEFAULT 2000)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_default_year_e8d89aba_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `y` year(4) NOT NULL DEFAULT '2000',
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_default_year_e8d89aba
< ok 0 0
> exec
| CREATE DATABASE pa_modify_default_year_e8d89aba DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_default_year_e8d89aba
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, y YEAR NOT NULL DEFAULT 1999)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_default_year_e8d89aba.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `y` year(4) NOT NULL DEFAULT '1999',
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_modify_default_year_e8d89aba`.`t` MODIFY COLUMN `y` year NOT NULL DEFAULT '2000'
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_default_year_e8d89aba.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `y` year(4) NOT NULL DEFAULT '2000',
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_charset_e9e266ed_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_charset_e9e266ed_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_charset_e9e266ed_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v VARCHAR(10) CHARACTER SET utf8mb4) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_charset_e9e266ed_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `v` varchar(10) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_charset_e9e266ed
< ok 0 0
> exec
| CREATE DATABASE pa_modify_charset_e9e266ed DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_charset_e9e266ed
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v VARCHAR(10) CHARACTER SET latin1) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_charset_e9e266ed.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `v` varchar(10) CHARACTER SET latin1 DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_modify_charset_e9e266ed`.`t` MODIFY COLUMN `v` varchar(10) DEFAULT NULL
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_charset_e9e266ed.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `v` varchar(10) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_column_binary_modifier_cbcae4b6_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_column_binary_modifier_cbcae4b6_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_column_binary_modifier_cbcae4b6_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, k VARCHAR(128) BINARY NOT NULL) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_column_binary_modifier_cbcae4b6_to.t
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
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_column_binary_modifier_cbcae4b6
< ok 0 0
> exec
| CREATE DATABASE pa_add_column_binary_modifier_cbcae4b6 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_column_binary_modifier_cbcae4b6
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_column_binary_modifier_cbcae4b6.t
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
| ALTER TABLE `pa_add_column_binary_modifier_cbcae4b6`.`t` ADD COLUMN `k` varchar(128) COLLATE utf8mb4_bin NOT NULL
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_column_binary_modifier_cbcae4b6.t
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
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_column_to_binary_modifier_89fc16ed_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_column_to_binary_modifier_89fc16ed_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_column_to_binary_modifier_89fc16ed_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v VARCHAR(10) BINARY) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_column_to_binary_modifier_89fc16ed_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `v` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_column_to_binary_modifier_89fc16ed
< ok 0 0
> exec
| CREATE DATABASE pa_modify_column_to_binary_modifier_89fc16ed DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_column_to_binary_modifier_89fc16ed
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v VARCHAR(10)) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_column_to_binary_modifier_89fc16ed.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `v` varchar(10) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_modify_column_to_binary_modifier_89fc16ed`.`t` MODIFY COLUMN `v` varchar(10) COLLATE utf8mb4_bin DEFAULT NULL
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_column_to_binary_modifier_89fc16ed.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `v` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_comment_16325158_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_comment_16325158_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_comment_16325158_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v INT COMMENT 'b')
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_comment_16325158_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `v` int(11) DEFAULT NULL COMMENT 'b',
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_comment_16325158
< ok 0 0
> exec
| CREATE DATABASE pa_modify_comment_16325158 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_comment_16325158
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v INT COMMENT 'a')
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_comment_16325158.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `v` int(11) DEFAULT NULL COMMENT 'a',
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_modify_comment_16325158`.`t` MODIFY COLUMN `v` int DEFAULT NULL COMMENT 'b'
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_comment_16325158.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `v` int(11) DEFAULT NULL COMMENT 'b',
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_add_autoinc_67791cae_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_add_autoinc_67791cae_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_add_autoinc_67791cae_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_add_autoinc_67791cae_to.t
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
| DROP DATABASE IF EXISTS pa_modify_add_autoinc_67791cae
< ok 0 0
> exec
| CREATE DATABASE pa_modify_add_autoinc_67791cae DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_add_autoinc_67791cae
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_add_autoinc_67791cae.t
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
| ALTER TABLE `pa_modify_add_autoinc_67791cae`.`t` MODIFY COLUMN `id` int NOT NULL AUTO_INCREMENT
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_add_autoinc_67791cae.t
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
| DROP DATABASE IF EXISTS pa_modify_generated_expr_3ab763b0_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_generated_expr_3ab763b0_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_generated_expr_3ab763b0_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, g INT GENERATED ALWAYS AS (a+2) STORED)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_generated_expr_3ab763b0_to.t
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
|   `g` int(11) GENERATED ALWAYS AS ((`a` + 2)) STORED,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_generated_expr_3ab763b0
< ok 0 0
> exec
| CREATE DATABASE pa_modify_generated_expr_3ab763b0 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_generated_expr_3ab763b0
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, g INT GENERATED ALWAYS AS (a+1) STORED)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_generated_expr_3ab763b0.t
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
|   `g` int(11) GENERATED ALWAYS AS ((`a` + 1)) STORED,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_modify_generated_expr_3ab763b0`.`t` DROP COLUMN `g`
< ok 0 0
> exec
| ALTER TABLE `pa_modify_generated_expr_3ab763b0`.`t` ADD COLUMN `g` int GENERATED ALWAYS AS (((`a` + 2))) STORED
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_generated_expr_3ab763b0.t
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
|   `g` int(11) GENERATED ALWAYS AS ((`a` + 2)) STORED,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_generated_expr_indexed_12220534_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_generated_expr_indexed_12220534_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_generated_expr_indexed_12220534_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, g INT GENERATED ALWAYS AS (a+2) VIRTUAL, KEY kg (g), KEY kga (g, a), KEY ka (a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_generated_expr_indexed_12220534_to.t
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
|   `g` int(11) GENERATED ALWAYS AS ((`a` + 2)) VIRTUAL,
|   PRIMARY KEY (`id`),
|   KEY `kg` (`g`) BLOCK_SIZE 16384 LOCAL,
|   KEY `kga` (`g`, `a`) BLOCK_SIZE 16384 LOCAL,
|   KEY `ka` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_generated_expr_indexed_12220534
< ok 0 0
> exec
| CREATE DATABASE pa_modify_generated_expr_indexed_12220534 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_generated_expr_indexed_12220534
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, g INT GENERATED ALWAYS AS (a+1) VIRTUAL, KEY kg (g), KEY kga (g, a), KEY ka (a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_generated_expr_indexed_12220534.t
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
|   `g` int(11) GENERATED ALWAYS AS ((`a` + 1)) VIRTUAL,
|   PRIMARY KEY (`id`),
|   KEY `kg` (`g`) BLOCK_SIZE 16384 LOCAL,
|   KEY `kga` (`g`, `a`) BLOCK_SIZE 16384 LOCAL,
|   KEY `ka` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_modify_generated_expr_indexed_12220534`.`t` DROP INDEX `kg`
< ok 0 0
> exec
| ALTER TABLE `pa_modify_generated_expr_indexed_12220534`.`t` DROP INDEX `kga`
< ok 0 0
> exec
| ALTER TABLE `pa_modify_generated_expr_indexed_12220534`.`t` DROP COLUMN `g`
< ok 0 0
> exec
| ALTER TABLE `pa_modify_generated_expr_indexed_12220534`.`t` ADD COLUMN `g` int GENERATED ALWAYS AS (((`a` + 2))) VIRTUAL
< ok 0 0
> exec
| ALTER TABLE `pa_modify_generated_expr_indexed_12220534`.`t` ADD KEY `kg` (`g`) LOCAL
< ok 0 0
> exec
| ALTER TABLE `pa_modify_generated_expr_indexed_12220534`.`t` ADD KEY `kga` (`g`,`a`) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_generated_expr_indexed_12220534.t
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
|   `g` int(11) GENERATED ALWAYS AS ((`a` + 2)) VIRTUAL,
|   PRIMARY KEY (`id`),
|   KEY `ka` (`a`) BLOCK_SIZE 16384 LOCAL,
|   KEY `kg` (`g`) BLOCK_SIZE 16384 LOCAL,
|   KEY `kga` (`g`, `a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_generated_expr_under_check_760670ca_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_generated_expr_under_check_760670ca_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_generated_expr_under_check_760670ca_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, g INT GENERATED ALWAYS AS (a+2) STORED, CONSTRAINT c_g CHECK (g > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_generated_expr_under_check_760670ca_to.t
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
|   `g` int(11) GENERATED ALWAYS AS ((`a` + 2)) STORED,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `c_g` CHECK ((`g` > 0))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_generated_expr_under_check_760670ca
< ok 0 0
> exec
| CREATE DATABASE pa_modify_generated_expr_under_check_760670ca DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_generated_expr_under_check_760670ca
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, g INT GENERATED ALWAYS AS (a+1) STORED, CONSTRAINT c_g CHECK (g > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_generated_expr_under_check_760670ca.t
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
|   `g` int(11) GENERATED ALWAYS AS ((`a` + 1)) STORED,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `c_g` CHECK ((`g` > 0))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_modify_generated_expr_under_check_760670ca`.`t` MODIFY COLUMN `g` int GENERATED ALWAYS AS (((`a` + 2))) STORED
< error
| 1235 0A000 Modify generated column definition not supported
> exec
| DROP DATABASE IF EXISTS pa_alter_engine_dff1830b_to
< ok 0 0
> exec
| CREATE DATABASE pa_alter_engine_dff1830b_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_alter_engine_dff1830b_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) ENGINE=MyISAM
< ok 0 0
> query
| SHOW CREATE TABLE pa_alter_engine_dff1830b_to.t
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
| DROP DATABASE IF EXISTS pa_alter_engine_dff1830b
< ok 0 0
> exec
| CREATE DATABASE pa_alter_engine_dff1830b DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_alter_engine_dff1830b
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) ENGINE=InnoDB
< ok 0 0
> query
| SHOW CREATE TABLE pa_alter_engine_dff1830b.t
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
> query ^
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
| DROP DATABASE IF EXISTS pa_alter_comment_484732ae_to
< ok 0 0
> exec
| CREATE DATABASE pa_alter_comment_484732ae_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_alter_comment_484732ae_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) COMMENT='new'
< ok 0 0
> query
| SHOW CREATE TABLE pa_alter_comment_484732ae_to.t
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
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 COMMENT = 'new'
> exec
| DROP DATABASE IF EXISTS pa_alter_comment_484732ae
< ok 0 0
> exec
| CREATE DATABASE pa_alter_comment_484732ae DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_alter_comment_484732ae
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) COMMENT='old'
< ok 0 0
> query
| SHOW CREATE TABLE pa_alter_comment_484732ae.t
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
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 COMMENT = 'old'
> exec
| ALTER TABLE `pa_alter_comment_484732ae`.`t` COMMENT='new'
< ok 0 0
> query
| SHOW CREATE TABLE pa_alter_comment_484732ae.t
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
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 COMMENT = 'new'
> exec
| DROP DATABASE IF EXISTS pa_alter_charset_fb6d7c5b_to
< ok 0 0
> exec
| CREATE DATABASE pa_alter_charset_fb6d7c5b_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_alter_charset_fb6d7c5b_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_alter_charset_fb6d7c5b_to.t
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
| DROP DATABASE IF EXISTS pa_alter_charset_fb6d7c5b
< ok 0 0
> exec
| CREATE DATABASE pa_alter_charset_fb6d7c5b DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_alter_charset_fb6d7c5b
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) DEFAULT CHARSET=latin1
< ok 0 0
> query
| SHOW CREATE TABLE pa_alter_charset_fb6d7c5b.t
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
| ) ORGANIZATION INDEX DEFAULT CHARSET = latin1 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_alter_charset_fb6d7c5b`.`t` DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_alter_charset_fb6d7c5b.t
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
| DROP DATABASE IF EXISTS pa_alter_rowformat_add_7f4528ae_to
< ok 0 0
> exec
| CREATE DATABASE pa_alter_rowformat_add_7f4528ae_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_alter_rowformat_add_7f4528ae_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=8
< ok 0 0
> query
| SHOW CREATE TABLE pa_alter_rowformat_add_7f4528ae_to.t
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
| DROP DATABASE IF EXISTS pa_alter_rowformat_add_7f4528ae
< ok 0 0
> exec
| CREATE DATABASE pa_alter_rowformat_add_7f4528ae DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_alter_rowformat_add_7f4528ae
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_alter_rowformat_add_7f4528ae.t
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
| ALTER TABLE `pa_alter_rowformat_add_7f4528ae`.`t` ROW_FORMAT=COMPRESSED
< ok 0 0
> query
| SHOW CREATE TABLE pa_alter_rowformat_add_7f4528ae.t
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
| DROP DATABASE IF EXISTS pa_alter_rowformat_change_4c74315f_to
< ok 0 0
> exec
| CREATE DATABASE pa_alter_rowformat_change_4c74315f_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_alter_rowformat_change_4c74315f_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) ROW_FORMAT=DYNAMIC
< ok 0 0
> query
| SHOW CREATE TABLE pa_alter_rowformat_change_4c74315f_to.t
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
| DROP DATABASE IF EXISTS pa_alter_rowformat_change_4c74315f
< ok 0 0
> exec
| CREATE DATABASE pa_alter_rowformat_change_4c74315f DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_alter_rowformat_change_4c74315f
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) ROW_FORMAT=COMPACT
< ok 0 0
> query
| SHOW CREATE TABLE pa_alter_rowformat_change_4c74315f.t
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
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = COMPACT COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_alter_rowformat_change_4c74315f`.`t` ROW_FORMAT=DEFAULT
< ok 0 0
> query
| SHOW CREATE TABLE pa_alter_rowformat_change_4c74315f.t
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
| DROP DATABASE IF EXISTS pa_gen_storage_flip_e1491f60_to
< ok 0 0
> exec
| CREATE DATABASE pa_gen_storage_flip_e1491f60_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_gen_storage_flip_e1491f60_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, g INT GENERATED ALWAYS AS (a+1) STORED)
< ok 0 0
> query
| SHOW CREATE TABLE pa_gen_storage_flip_e1491f60_to.t
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
|   `g` int(11) GENERATED ALWAYS AS ((`a` + 1)) STORED,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_gen_storage_flip_e1491f60
< ok 0 0
> exec
| CREATE DATABASE pa_gen_storage_flip_e1491f60 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_gen_storage_flip_e1491f60
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, g INT GENERATED ALWAYS AS (a+1) VIRTUAL)
< ok 0 0
> query
| SHOW CREATE TABLE pa_gen_storage_flip_e1491f60.t
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
|   `g` int(11) GENERATED ALWAYS AS ((`a` + 1)) VIRTUAL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_gen_storage_flip_e1491f60`.`t` DROP COLUMN `g`
< ok 0 0
> exec
| ALTER TABLE `pa_gen_storage_flip_e1491f60`.`t` ADD COLUMN `g` int GENERATED ALWAYS AS (((`a` + 1))) STORED
< ok 0 0
> query
| SHOW CREATE TABLE pa_gen_storage_flip_e1491f60.t
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
|   `g` int(11) GENERATED ALWAYS AS ((`a` + 1)) STORED,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_gen_drop_with_base_3206fcc5_to
< ok 0 0
> exec
| CREATE DATABASE pa_gen_drop_with_base_3206fcc5_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_gen_drop_with_base_3206fcc5_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_gen_drop_with_base_3206fcc5_to.t
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
| DROP DATABASE IF EXISTS pa_gen_drop_with_base_3206fcc5
< ok 0 0
> exec
| CREATE DATABASE pa_gen_drop_with_base_3206fcc5 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_gen_drop_with_base_3206fcc5
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, z INT GENERATED ALWAYS AS (a+1) STORED)
< ok 0 0
> query
| SHOW CREATE TABLE pa_gen_drop_with_base_3206fcc5.t
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
|   `z` int(11) GENERATED ALWAYS AS ((`a` + 1)) STORED,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_gen_drop_with_base_3206fcc5`.`t` DROP COLUMN `z`
< ok 0 0
> exec
| ALTER TABLE `pa_gen_drop_with_base_3206fcc5`.`t` DROP COLUMN `a`
< ok 0 0
> query
| SHOW CREATE TABLE pa_gen_drop_with_base_3206fcc5.t
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
| DROP DATABASE IF EXISTS pa_gen_flip_with_base_modify_98b424d4_to
< ok 0 0
> exec
| CREATE DATABASE pa_gen_flip_with_base_modify_98b424d4_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_gen_flip_with_base_modify_98b424d4_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g INT GENERATED ALWAYS AS (z+1) STORED, z BIGINT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_gen_flip_with_base_modify_98b424d4_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `g` int(11) GENERATED ALWAYS AS ((`z` + 1)) STORED,
|   `z` bigint(20) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_gen_flip_with_base_modify_98b424d4
< ok 0 0
> exec
| CREATE DATABASE pa_gen_flip_with_base_modify_98b424d4 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_gen_flip_with_base_modify_98b424d4
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g INT GENERATED ALWAYS AS (z+1) VIRTUAL, z SMALLINT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_gen_flip_with_base_modify_98b424d4.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `g` int(11) GENERATED ALWAYS AS ((`z` + 1)) VIRTUAL,
|   `z` smallint(6) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_gen_flip_with_base_modify_98b424d4`.`t` DROP COLUMN `g`
< ok 0 0
> exec
| ALTER TABLE `pa_gen_flip_with_base_modify_98b424d4`.`t` MODIFY COLUMN `z` bigint DEFAULT NULL
< ok 0 0
> exec
| ALTER TABLE `pa_gen_flip_with_base_modify_98b424d4`.`t` ADD COLUMN `g` int GENERATED ALWAYS AS (((`z` + 1))) STORED
< ok 0 0
> query
| SHOW CREATE TABLE pa_gen_flip_with_base_modify_98b424d4.t
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
|   `g` int(11) GENERATED ALWAYS AS ((`z` + 1)) STORED,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_invisible_column_34ae6591_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_invisible_column_34ae6591_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_invisible_column_34ae6591_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, secret INT INVISIBLE)
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_invisible_column_34ae6591_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `secret` int(11) INVISIBLE DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_invisible_column_34ae6591
< ok 0 0
> exec
| CREATE DATABASE pa_create_invisible_column_34ae6591 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_invisible_column_34ae6591
< ok 0 0
> exec
| CREATE TABLE `pa_create_invisible_column_34ae6591`.`t` (
|   `id` int NOT NULL,
|   `secret` int DEFAULT NULL /*!80023 INVISIBLE */,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_invisible_column_34ae6591.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `secret` int(11) INVISIBLE DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_invisible_column_6ae27004_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_invisible_column_6ae27004_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_invisible_column_6ae27004_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, secret INT INVISIBLE)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_invisible_column_6ae27004_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `secret` int(11) INVISIBLE DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_invisible_column_6ae27004
< ok 0 0
> exec
| CREATE DATABASE pa_add_invisible_column_6ae27004 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_invisible_column_6ae27004
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_invisible_column_6ae27004.t
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
| ALTER TABLE `pa_add_invisible_column_6ae27004`.`t` ADD COLUMN `secret` int DEFAULT NULL /*!80023 INVISIBLE */
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_invisible_column_6ae27004.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `secret` int(11) INVISIBLE DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_collate_only_column_c0047233_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_collate_only_column_c0047233_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_collate_only_column_c0047233_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a VARCHAR(10) COLLATE utf8mb4_unicode_ci) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_collate_only_column_c0047233_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_collate_only_column_c0047233
< ok 0 0
> exec
| CREATE DATABASE pa_create_collate_only_column_c0047233 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_collate_only_column_c0047233
< ok 0 0
> exec
| CREATE TABLE `pa_create_collate_only_column_c0047233`.`t` (
|   `id` int NOT NULL,
|   `a` varchar(10) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_collate_only_column_c0047233.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
