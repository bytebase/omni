-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestContainer_Section_1_1_NumericTypes default#1
> exec
| DROP TABLE IF EXISTS t_int
< ok 0 0
> exec
| CREATE TABLE t_int (a INT)
< ok 0 0
> query
| SHOW CREATE TABLE t_int
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_int
: bytes
| CREATE TABLE `t_int` (
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_int_dw
< ok 0 0
> exec
| CREATE TABLE t_int_dw (a INT(11))
< ok 0 0
> query
| SHOW CREATE TABLE t_int_dw
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_int_dw
: bytes
| CREATE TABLE `t_int_dw` (
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_int_u
< ok 0 0
> exec
| CREATE TABLE t_int_u (a INT UNSIGNED)
< ok 0 0
> query
| SHOW CREATE TABLE t_int_u
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_int_u
: bytes
| CREATE TABLE `t_int_u` (
|   `a` int(10) unsigned DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_int_uz
< ok 0 0
> exec
| CREATE TABLE t_int_uz (a INT UNSIGNED ZEROFILL)
< ok 0 0
> query
| SHOW CREATE TABLE t_int_uz
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_int_uz
: bytes
| CREATE TABLE `t_int_uz` (
|   `a` int(10) unsigned zerofill DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_tinyint
< ok 0 0
> exec
| CREATE TABLE t_tinyint (a TINYINT)
< ok 0 0
> query
| SHOW CREATE TABLE t_tinyint
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_tinyint
: bytes
| CREATE TABLE `t_tinyint` (
|   `a` tinyint(4) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_smallint
< ok 0 0
> exec
| CREATE TABLE t_smallint (a SMALLINT)
< ok 0 0
> query
| SHOW CREATE TABLE t_smallint
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_smallint
: bytes
| CREATE TABLE `t_smallint` (
|   `a` smallint(6) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_mediumint
< ok 0 0
> exec
| CREATE TABLE t_mediumint (a MEDIUMINT)
< ok 0 0
> query
| SHOW CREATE TABLE t_mediumint
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_mediumint
: bytes
| CREATE TABLE `t_mediumint` (
|   `a` mediumint(9) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_bigint
< ok 0 0
> exec
| CREATE TABLE t_bigint (a BIGINT)
< ok 0 0
> query
| SHOW CREATE TABLE t_bigint
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_bigint
: bytes
| CREATE TABLE `t_bigint` (
|   `a` bigint(20) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_bigint_u
< ok 0 0
> exec
| CREATE TABLE t_bigint_u (a BIGINT UNSIGNED)
< ok 0 0
> query
| SHOW CREATE TABLE t_bigint_u
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_bigint_u
: bytes
| CREATE TABLE `t_bigint_u` (
|   `a` bigint(20) unsigned DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_float
< ok 0 0
> exec
| CREATE TABLE t_float (a FLOAT)
< ok 0 0
> query
| SHOW CREATE TABLE t_float
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_float
: bytes
| CREATE TABLE `t_float` (
|   `a` float DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_float_p
< ok 0 0
> exec
| CREATE TABLE t_float_p (a FLOAT(7,3))
< ok 0 0
> query
| SHOW CREATE TABLE t_float_p
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_float_p
: bytes
| CREATE TABLE `t_float_p` (
|   `a` float(7,3) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_float_u
< ok 0 0
> exec
| CREATE TABLE t_float_u (a FLOAT UNSIGNED)
< ok 0 0
> query
| SHOW CREATE TABLE t_float_u
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_float_u
: bytes
| CREATE TABLE `t_float_u` (
|   `a` float unsigned DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_double
< ok 0 0
> exec
| CREATE TABLE t_double (a DOUBLE)
< ok 0 0
> query
| SHOW CREATE TABLE t_double
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_double
: bytes
| CREATE TABLE `t_double` (
|   `a` double DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_double_p
< ok 0 0
> exec
| CREATE TABLE t_double_p (a DOUBLE PRECISION)
< ok 0 0
> query
| SHOW CREATE TABLE t_double_p
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_double_p
: bytes
| CREATE TABLE `t_double_p` (
|   `a` double DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_double_wp
< ok 0 0
> exec
| CREATE TABLE t_double_wp (a DOUBLE(15,5))
< ok 0 0
> query
| SHOW CREATE TABLE t_double_wp
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_double_wp
: bytes
| CREATE TABLE `t_double_wp` (
|   `a` double(15,5) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_decimal
< ok 0 0
> exec
| CREATE TABLE t_decimal (a DECIMAL(10,2))
< ok 0 0
> query
| SHOW CREATE TABLE t_decimal
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_decimal
: bytes
| CREATE TABLE `t_decimal` (
|   `a` decimal(10,2) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_numeric
< ok 0 0
> exec
| CREATE TABLE t_numeric (a NUMERIC(10,2))
< ok 0 0
> query
| SHOW CREATE TABLE t_numeric
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_numeric
: bytes
| CREATE TABLE `t_numeric` (
|   `a` decimal(10,2) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_decimal_np
< ok 0 0
> exec
| CREATE TABLE t_decimal_np (a DECIMAL)
< ok 0 0
> query
| SHOW CREATE TABLE t_decimal_np
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_decimal_np
: bytes
| CREATE TABLE `t_decimal_np` (
|   `a` decimal(10,0) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_bool
< ok 0 0
> exec
| CREATE TABLE t_bool (a BOOLEAN)
< ok 0 0
> query
| SHOW CREATE TABLE t_bool
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_bool
: bytes
| CREATE TABLE `t_bool` (
|   `a` tinyint(1) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_bool2
< ok 0 0
> exec
| CREATE TABLE t_bool2 (a BOOL)
< ok 0 0
> query
| SHOW CREATE TABLE t_bool2
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_bool2
: bytes
| CREATE TABLE `t_bool2` (
|   `a` tinyint(1) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_bit1
< ok 0 0
> exec
| CREATE TABLE t_bit1 (a BIT(1))
< ok 0 0
> query
| SHOW CREATE TABLE t_bit1
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_bit1
: bytes
| CREATE TABLE `t_bit1` (
|   `a` bit(1) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_bit8
< ok 0 0
> exec
| CREATE TABLE t_bit8 (a BIT(8))
< ok 0 0
> query
| SHOW CREATE TABLE t_bit8
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_bit8
: bytes
| CREATE TABLE `t_bit8` (
|   `a` bit(8) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_bit64
< ok 0 0
> exec
| CREATE TABLE t_bit64 (a BIT(64))
< ok 0 0
> query
| SHOW CREATE TABLE t_bit64
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_bit64
: bytes
| CREATE TABLE `t_bit64` (
|   `a` bit(64) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_serial
< ok 0 0
> exec
| CREATE TABLE t_serial (a SERIAL)
< ok 0 0
> query
| SHOW CREATE TABLE t_serial
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_serial
: bytes
| CREATE TABLE `t_serial` (
|   `a` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
|   UNIQUE KEY `a` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
