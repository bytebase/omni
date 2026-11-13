-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracle_PhantomDiffElimination default#1
> exec
| DROP DATABASE IF EXISTS nrm_int_display_width
< ok 0 0
> exec
| CREATE DATABASE nrm_int_display_width
< ok 1 0
> exec
| USE nrm_int_display_width
< ok 0 0
> exec
| CREATE TABLE t_int (a INT(11), b BIGINT(20), c INT, f SMALLINT(6), o BIGINT, p TINYINT, q SMALLINT, r MEDIUMINT)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_int_display_width.t_int
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_int
: bytes
| CREATE TABLE `t_int` (
|   `a` int DEFAULT NULL,
|   `b` bigint DEFAULT NULL,
|   `c` int DEFAULT NULL,
|   `f` smallint DEFAULT NULL,
|   `o` bigint DEFAULT NULL,
|   `p` tinyint DEFAULT NULL,
|   `q` smallint DEFAULT NULL,
|   `r` mediumint DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_int_display_width
< ok 1 0
> exec
| CREATE DATABASE nrm_int_display_width
< ok 1 0
> exec
| USE nrm_int_display_width
< ok 0 0
> exec
| CREATE TABLE t_int (a INT(11), b BIGINT(20), c INT, f SMALLINT(6), o BIGINT, p TINYINT, q SMALLINT, r MEDIUMINT)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_int_display_width.t_int
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_int
: bytes
| CREATE TABLE `t_int` (
|   `a` int DEFAULT NULL,
|   `b` bigint DEFAULT NULL,
|   `c` int DEFAULT NULL,
|   `f` smallint DEFAULT NULL,
|   `o` bigint DEFAULT NULL,
|   `p` tinyint DEFAULT NULL,
|   `q` smallint DEFAULT NULL,
|   `r` mediumint DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_int_display_width
< ok 1 0
> exec
| CREATE DATABASE nrm_int_display_width
< ok 1 0
> exec
| USE nrm_int_display_width
< ok 0 0
> exec
| CREATE TABLE t_int (a INT(11), b BIGINT(20), c INT, f SMALLINT(6), o BIGINT, p TINYINT, q SMALLINT, r MEDIUMINT)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_int_display_width.t_int
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_int
: bytes
| CREATE TABLE `t_int` (
|   `a` int DEFAULT NULL,
|   `b` bigint DEFAULT NULL,
|   `c` int DEFAULT NULL,
|   `f` smallint DEFAULT NULL,
|   `o` bigint DEFAULT NULL,
|   `p` tinyint DEFAULT NULL,
|   `q` smallint DEFAULT NULL,
|   `r` mediumint DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_int_display_width
< ok 1 0
> exec
| CREATE DATABASE nrm_int_display_width
< ok 1 0
> exec
| USE nrm_int_display_width
< ok 0 0
> exec
| CREATE TABLE t_int (a INT(11), b BIGINT(20), c INT, f SMALLINT(6), o BIGINT, p TINYINT, q SMALLINT, r MEDIUMINT)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_int_display_width.t_int
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_int
: bytes
| CREATE TABLE `t_int` (
|   `a` int DEFAULT NULL,
|   `b` bigint DEFAULT NULL,
|   `c` int DEFAULT NULL,
|   `f` smallint DEFAULT NULL,
|   `o` bigint DEFAULT NULL,
|   `p` tinyint DEFAULT NULL,
|   `q` smallint DEFAULT NULL,
|   `r` mediumint DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_int_display_width
< ok 1 0
> exec
| CREATE DATABASE nrm_int_display_width
< ok 1 0
> exec
| USE nrm_int_display_width
< ok 0 0
> exec
| CREATE TABLE t_int (a INT(11), b BIGINT(20), c INT, f SMALLINT(6), o BIGINT, p TINYINT, q SMALLINT, r MEDIUMINT)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_int_display_width.t_int
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_int
: bytes
| CREATE TABLE `t_int` (
|   `a` int DEFAULT NULL,
|   `b` bigint DEFAULT NULL,
|   `c` int DEFAULT NULL,
|   `f` smallint DEFAULT NULL,
|   `o` bigint DEFAULT NULL,
|   `p` tinyint DEFAULT NULL,
|   `q` smallint DEFAULT NULL,
|   `r` mediumint DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_int_display_width
< ok 1 0
> exec
| CREATE DATABASE nrm_int_display_width
< ok 1 0
> exec
| USE nrm_int_display_width
< ok 0 0
> exec
| CREATE TABLE t_int (a INT(11), b BIGINT(20), c INT, f SMALLINT(6), o BIGINT, p TINYINT, q SMALLINT, r MEDIUMINT)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_int_display_width.t_int
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_int
: bytes
| CREATE TABLE `t_int` (
|   `a` int DEFAULT NULL,
|   `b` bigint DEFAULT NULL,
|   `c` int DEFAULT NULL,
|   `f` smallint DEFAULT NULL,
|   `o` bigint DEFAULT NULL,
|   `p` tinyint DEFAULT NULL,
|   `q` smallint DEFAULT NULL,
|   `r` mediumint DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_int_display_width
< ok 1 0
> exec
| CREATE DATABASE nrm_int_display_width
< ok 1 0
> exec
| USE nrm_int_display_width
< ok 0 0
> exec
| CREATE TABLE t_int (a INT(11), b BIGINT(20), c INT, f SMALLINT(6), o BIGINT, p TINYINT, q SMALLINT, r MEDIUMINT)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_int_display_width.t_int
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_int
: bytes
| CREATE TABLE `t_int` (
|   `a` int DEFAULT NULL,
|   `b` bigint DEFAULT NULL,
|   `c` int DEFAULT NULL,
|   `f` smallint DEFAULT NULL,
|   `o` bigint DEFAULT NULL,
|   `p` tinyint DEFAULT NULL,
|   `q` smallint DEFAULT NULL,
|   `r` mediumint DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_int_display_width
< ok 1 0
> exec
| CREATE DATABASE nrm_int_display_width
< ok 1 0
> exec
| USE nrm_int_display_width
< ok 0 0
> exec
| CREATE TABLE t_int (a INT(11), b BIGINT(20), c INT, f SMALLINT(6), o BIGINT, p TINYINT, q SMALLINT, r MEDIUMINT)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_int_display_width.t_int
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_int
: bytes
| CREATE TABLE `t_int` (
|   `a` int DEFAULT NULL,
|   `b` bigint DEFAULT NULL,
|   `c` int DEFAULT NULL,
|   `f` smallint DEFAULT NULL,
|   `o` bigint DEFAULT NULL,
|   `p` tinyint DEFAULT NULL,
|   `q` smallint DEFAULT NULL,
|   `r` mediumint DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_tinyint1_boolean
< ok 0 0
> exec
| CREATE DATABASE nrm_tinyint1_boolean
< ok 1 0
> exec
| USE nrm_tinyint1_boolean
< ok 0 0
> exec
| CREATE TABLE t_bool (d TINYINT(1), e TINYINT(4), m BOOLEAN, n BOOL, p TINYINT)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_tinyint1_boolean.t_bool
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_bool
: bytes
| CREATE TABLE `t_bool` (
|   `d` tinyint(1) DEFAULT NULL,
|   `e` tinyint DEFAULT NULL,
|   `m` tinyint(1) DEFAULT NULL,
|   `n` tinyint(1) DEFAULT NULL,
|   `p` tinyint DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_tinyint1_boolean
< ok 1 0
> exec
| CREATE DATABASE nrm_tinyint1_boolean
< ok 1 0
> exec
| USE nrm_tinyint1_boolean
< ok 0 0
> exec
| CREATE TABLE t_bool (d TINYINT(1), e TINYINT(4), m BOOLEAN, n BOOL, p TINYINT)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_tinyint1_boolean.t_bool
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_bool
: bytes
| CREATE TABLE `t_bool` (
|   `d` tinyint(1) DEFAULT NULL,
|   `e` tinyint DEFAULT NULL,
|   `m` tinyint(1) DEFAULT NULL,
|   `n` tinyint(1) DEFAULT NULL,
|   `p` tinyint DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_tinyint1_boolean
< ok 1 0
> exec
| CREATE DATABASE nrm_tinyint1_boolean
< ok 1 0
> exec
| USE nrm_tinyint1_boolean
< ok 0 0
> exec
| CREATE TABLE t_bool (d TINYINT(1), e TINYINT(4), m BOOLEAN, n BOOL, p TINYINT)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_tinyint1_boolean.t_bool
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_bool
: bytes
| CREATE TABLE `t_bool` (
|   `d` tinyint(1) DEFAULT NULL,
|   `e` tinyint DEFAULT NULL,
|   `m` tinyint(1) DEFAULT NULL,
|   `n` tinyint(1) DEFAULT NULL,
|   `p` tinyint DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_tinyint1_boolean
< ok 1 0
> exec
| CREATE DATABASE nrm_tinyint1_boolean
< ok 1 0
> exec
| USE nrm_tinyint1_boolean
< ok 0 0
> exec
| CREATE TABLE t_bool (d TINYINT(1), e TINYINT(4), m BOOLEAN, n BOOL, p TINYINT)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_tinyint1_boolean.t_bool
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_bool
: bytes
| CREATE TABLE `t_bool` (
|   `d` tinyint(1) DEFAULT NULL,
|   `e` tinyint DEFAULT NULL,
|   `m` tinyint(1) DEFAULT NULL,
|   `n` tinyint(1) DEFAULT NULL,
|   `p` tinyint DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_tinyint1_boolean
< ok 1 0
> exec
| CREATE DATABASE nrm_tinyint1_boolean
< ok 1 0
> exec
| USE nrm_tinyint1_boolean
< ok 0 0
> exec
| CREATE TABLE t_bool (d TINYINT(1), e TINYINT(4), m BOOLEAN, n BOOL, p TINYINT)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_tinyint1_boolean.t_bool
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_bool
: bytes
| CREATE TABLE `t_bool` (
|   `d` tinyint(1) DEFAULT NULL,
|   `e` tinyint DEFAULT NULL,
|   `m` tinyint(1) DEFAULT NULL,
|   `n` tinyint(1) DEFAULT NULL,
|   `p` tinyint DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_int_unsigned_width
< ok 0 0
> exec
| CREATE DATABASE nrm_int_unsigned_width
< ok 1 0
> exec
| USE nrm_int_unsigned_width
< ok 0 0
> exec
| CREATE TABLE t_uns (h INT UNSIGNED, i INT(11) UNSIGNED, b BIGINT UNSIGNED, ti TINYINT UNSIGNED)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_int_unsigned_width.t_uns
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_uns
: bytes
| CREATE TABLE `t_uns` (
|   `h` int unsigned DEFAULT NULL,
|   `i` int unsigned DEFAULT NULL,
|   `b` bigint unsigned DEFAULT NULL,
|   `ti` tinyint unsigned DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_int_unsigned_width
< ok 1 0
> exec
| CREATE DATABASE nrm_int_unsigned_width
< ok 1 0
> exec
| USE nrm_int_unsigned_width
< ok 0 0
> exec
| CREATE TABLE t_uns (h INT UNSIGNED, i INT(11) UNSIGNED, b BIGINT UNSIGNED, ti TINYINT UNSIGNED)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_int_unsigned_width.t_uns
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_uns
: bytes
| CREATE TABLE `t_uns` (
|   `h` int unsigned DEFAULT NULL,
|   `i` int unsigned DEFAULT NULL,
|   `b` bigint unsigned DEFAULT NULL,
|   `ti` tinyint unsigned DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_int_unsigned_width
< ok 1 0
> exec
| CREATE DATABASE nrm_int_unsigned_width
< ok 1 0
> exec
| USE nrm_int_unsigned_width
< ok 0 0
> exec
| CREATE TABLE t_uns (h INT UNSIGNED, i INT(11) UNSIGNED, b BIGINT UNSIGNED, ti TINYINT UNSIGNED)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_int_unsigned_width.t_uns
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_uns
: bytes
| CREATE TABLE `t_uns` (
|   `h` int unsigned DEFAULT NULL,
|   `i` int unsigned DEFAULT NULL,
|   `b` bigint unsigned DEFAULT NULL,
|   `ti` tinyint unsigned DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_int_unsigned_width
< ok 1 0
> exec
| CREATE DATABASE nrm_int_unsigned_width
< ok 1 0
> exec
| USE nrm_int_unsigned_width
< ok 0 0
> exec
| CREATE TABLE t_uns (h INT UNSIGNED, i INT(11) UNSIGNED, b BIGINT UNSIGNED, ti TINYINT UNSIGNED)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_int_unsigned_width.t_uns
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_uns
: bytes
| CREATE TABLE `t_uns` (
|   `h` int unsigned DEFAULT NULL,
|   `i` int unsigned DEFAULT NULL,
|   `b` bigint unsigned DEFAULT NULL,
|   `ti` tinyint unsigned DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_int_zerofill
< ok 0 0
> exec
| CREATE DATABASE nrm_int_zerofill
< ok 1 0
> exec
| USE nrm_int_zerofill
< ok 0 0
> exec
| CREATE TABLE t_zf (j INT ZEROFILL, k INT(5) ZEROFILL, l INT UNSIGNED ZEROFILL)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_int_zerofill.t_zf
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_zf
: bytes
| CREATE TABLE `t_zf` (
|   `j` int(10) unsigned zerofill DEFAULT NULL,
|   `k` int(5) unsigned zerofill DEFAULT NULL,
|   `l` int(10) unsigned zerofill DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_int_zerofill
< ok 1 0
> exec
| CREATE DATABASE nrm_int_zerofill
< ok 1 0
> exec
| USE nrm_int_zerofill
< ok 0 0
> exec
| CREATE TABLE t_zf (j INT ZEROFILL, k INT(5) ZEROFILL, l INT UNSIGNED ZEROFILL)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_int_zerofill.t_zf
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_zf
: bytes
| CREATE TABLE `t_zf` (
|   `j` int(10) unsigned zerofill DEFAULT NULL,
|   `k` int(5) unsigned zerofill DEFAULT NULL,
|   `l` int(10) unsigned zerofill DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_int_zerofill
< ok 1 0
> exec
| CREATE DATABASE nrm_int_zerofill
< ok 1 0
> exec
| USE nrm_int_zerofill
< ok 0 0
> exec
| CREATE TABLE t_zf (j INT ZEROFILL, k INT(5) ZEROFILL, l INT UNSIGNED ZEROFILL)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_int_zerofill.t_zf
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_zf
: bytes
| CREATE TABLE `t_zf` (
|   `j` int(10) unsigned zerofill DEFAULT NULL,
|   `k` int(5) unsigned zerofill DEFAULT NULL,
|   `l` int(10) unsigned zerofill DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_decimal_precision_scale
< ok 0 0
> exec
| CREATE DATABASE nrm_decimal_precision_scale
< ok 1 0
> exec
| USE nrm_decimal_precision_scale
< ok 0 0
> exec
| CREATE TABLE t_num (a DECIMAL, b DECIMAL(10), c DECIMAL(10,2), d NUMERIC, e NUMERIC(8,3), k DEC, l FIXED)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_decimal_precision_scale.t_num
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_num
: bytes
| CREATE TABLE `t_num` (
|   `a` decimal(10,0) DEFAULT NULL,
|   `b` decimal(10,0) DEFAULT NULL,
|   `c` decimal(10,2) DEFAULT NULL,
|   `d` decimal(10,0) DEFAULT NULL,
|   `e` decimal(8,3) DEFAULT NULL,
|   `k` decimal(10,0) DEFAULT NULL,
|   `l` decimal(10,0) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_decimal_precision_scale
< ok 1 0
> exec
| CREATE DATABASE nrm_decimal_precision_scale
< ok 1 0
> exec
| USE nrm_decimal_precision_scale
< ok 0 0
> exec
| CREATE TABLE t_num (a DECIMAL, b DECIMAL(10), c DECIMAL(10,2), d NUMERIC, e NUMERIC(8,3), k DEC, l FIXED)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_decimal_precision_scale.t_num
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_num
: bytes
| CREATE TABLE `t_num` (
|   `a` decimal(10,0) DEFAULT NULL,
|   `b` decimal(10,0) DEFAULT NULL,
|   `c` decimal(10,2) DEFAULT NULL,
|   `d` decimal(10,0) DEFAULT NULL,
|   `e` decimal(8,3) DEFAULT NULL,
|   `k` decimal(10,0) DEFAULT NULL,
|   `l` decimal(10,0) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_decimal_precision_scale
< ok 1 0
> exec
| CREATE DATABASE nrm_decimal_precision_scale
< ok 1 0
> exec
| USE nrm_decimal_precision_scale
< ok 0 0
> exec
| CREATE TABLE t_num (a DECIMAL, b DECIMAL(10), c DECIMAL(10,2), d NUMERIC, e NUMERIC(8,3), k DEC, l FIXED)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_decimal_precision_scale.t_num
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_num
: bytes
| CREATE TABLE `t_num` (
|   `a` decimal(10,0) DEFAULT NULL,
|   `b` decimal(10,0) DEFAULT NULL,
|   `c` decimal(10,2) DEFAULT NULL,
|   `d` decimal(10,0) DEFAULT NULL,
|   `e` decimal(8,3) DEFAULT NULL,
|   `k` decimal(10,0) DEFAULT NULL,
|   `l` decimal(10,0) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_decimal_precision_scale
< ok 1 0
> exec
| CREATE DATABASE nrm_decimal_precision_scale
< ok 1 0
> exec
| USE nrm_decimal_precision_scale
< ok 0 0
> exec
| CREATE TABLE t_num (a DECIMAL, b DECIMAL(10), c DECIMAL(10,2), d NUMERIC, e NUMERIC(8,3), k DEC, l FIXED)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_decimal_precision_scale.t_num
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_num
: bytes
| CREATE TABLE `t_num` (
|   `a` decimal(10,0) DEFAULT NULL,
|   `b` decimal(10,0) DEFAULT NULL,
|   `c` decimal(10,2) DEFAULT NULL,
|   `d` decimal(10,0) DEFAULT NULL,
|   `e` decimal(8,3) DEFAULT NULL,
|   `k` decimal(10,0) DEFAULT NULL,
|   `l` decimal(10,0) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_decimal_precision_scale
< ok 1 0
> exec
| CREATE DATABASE nrm_decimal_precision_scale
< ok 1 0
> exec
| USE nrm_decimal_precision_scale
< ok 0 0
> exec
| CREATE TABLE t_num (a DECIMAL, b DECIMAL(10), c DECIMAL(10,2), d NUMERIC, e NUMERIC(8,3), k DEC, l FIXED)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_decimal_precision_scale.t_num
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_num
: bytes
| CREATE TABLE `t_num` (
|   `a` decimal(10,0) DEFAULT NULL,
|   `b` decimal(10,0) DEFAULT NULL,
|   `c` decimal(10,2) DEFAULT NULL,
|   `d` decimal(10,0) DEFAULT NULL,
|   `e` decimal(8,3) DEFAULT NULL,
|   `k` decimal(10,0) DEFAULT NULL,
|   `l` decimal(10,0) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_decimal_precision_scale
< ok 1 0
> exec
| CREATE DATABASE nrm_decimal_precision_scale
< ok 1 0
> exec
| USE nrm_decimal_precision_scale
< ok 0 0
> exec
| CREATE TABLE t_num (a DECIMAL, b DECIMAL(10), c DECIMAL(10,2), d NUMERIC, e NUMERIC(8,3), k DEC, l FIXED)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_decimal_precision_scale.t_num
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_num
: bytes
| CREATE TABLE `t_num` (
|   `a` decimal(10,0) DEFAULT NULL,
|   `b` decimal(10,0) DEFAULT NULL,
|   `c` decimal(10,2) DEFAULT NULL,
|   `d` decimal(10,0) DEFAULT NULL,
|   `e` decimal(8,3) DEFAULT NULL,
|   `k` decimal(10,0) DEFAULT NULL,
|   `l` decimal(10,0) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_decimal_precision_scale
< ok 1 0
> exec
| CREATE DATABASE nrm_decimal_precision_scale
< ok 1 0
> exec
| USE nrm_decimal_precision_scale
< ok 0 0
> exec
| CREATE TABLE t_num (a DECIMAL, b DECIMAL(10), c DECIMAL(10,2), d NUMERIC, e NUMERIC(8,3), k DEC, l FIXED)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_decimal_precision_scale.t_num
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_num
: bytes
| CREATE TABLE `t_num` (
|   `a` decimal(10,0) DEFAULT NULL,
|   `b` decimal(10,0) DEFAULT NULL,
|   `c` decimal(10,2) DEFAULT NULL,
|   `d` decimal(10,0) DEFAULT NULL,
|   `e` decimal(8,3) DEFAULT NULL,
|   `k` decimal(10,0) DEFAULT NULL,
|   `l` decimal(10,0) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_float_double_aliasing
< ok 0 0
> exec
| CREATE DATABASE nrm_float_double_aliasing
< ok 1 0
> exec
| USE nrm_float_double_aliasing
< ok 0 0
> exec
| CREATE TABLE t_fd (f FLOAT, g FLOAT(10,2), h DOUBLE, i DOUBLE(15,4), j REAL, m FLOAT(5))
< ok 0 0
> query
| SHOW CREATE TABLE nrm_float_double_aliasing.t_fd
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_fd
: bytes
| CREATE TABLE `t_fd` (
|   `f` float DEFAULT NULL,
|   `g` float(10,2) DEFAULT NULL,
|   `h` double DEFAULT NULL,
|   `i` double(15,4) DEFAULT NULL,
|   `j` double DEFAULT NULL,
|   `m` float DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_float_double_aliasing
< ok 1 0
> exec
| CREATE DATABASE nrm_float_double_aliasing
< ok 1 0
> exec
| USE nrm_float_double_aliasing
< ok 0 0
> exec
| CREATE TABLE t_fd (f FLOAT, g FLOAT(10,2), h DOUBLE, i DOUBLE(15,4), j REAL, m FLOAT(5))
< ok 0 0
> query
| SHOW CREATE TABLE nrm_float_double_aliasing.t_fd
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_fd
: bytes
| CREATE TABLE `t_fd` (
|   `f` float DEFAULT NULL,
|   `g` float(10,2) DEFAULT NULL,
|   `h` double DEFAULT NULL,
|   `i` double(15,4) DEFAULT NULL,
|   `j` double DEFAULT NULL,
|   `m` float DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_float_double_aliasing
< ok 1 0
> exec
| CREATE DATABASE nrm_float_double_aliasing
< ok 1 0
> exec
| USE nrm_float_double_aliasing
< ok 0 0
> exec
| CREATE TABLE t_fd (f FLOAT, g FLOAT(10,2), h DOUBLE, i DOUBLE(15,4), j REAL, m FLOAT(5))
< ok 0 0
> query
| SHOW CREATE TABLE nrm_float_double_aliasing.t_fd
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_fd
: bytes
| CREATE TABLE `t_fd` (
|   `f` float DEFAULT NULL,
|   `g` float(10,2) DEFAULT NULL,
|   `h` double DEFAULT NULL,
|   `i` double(15,4) DEFAULT NULL,
|   `j` double DEFAULT NULL,
|   `m` float DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_float_double_aliasing
< ok 1 0
> exec
| CREATE DATABASE nrm_float_double_aliasing
< ok 1 0
> exec
| USE nrm_float_double_aliasing
< ok 0 0
> exec
| CREATE TABLE t_fd (f FLOAT, g FLOAT(10,2), h DOUBLE, i DOUBLE(15,4), j REAL, m FLOAT(5))
< ok 0 0
> query
| SHOW CREATE TABLE nrm_float_double_aliasing.t_fd
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_fd
: bytes
| CREATE TABLE `t_fd` (
|   `f` float DEFAULT NULL,
|   `g` float(10,2) DEFAULT NULL,
|   `h` double DEFAULT NULL,
|   `i` double(15,4) DEFAULT NULL,
|   `j` double DEFAULT NULL,
|   `m` float DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_float_double_aliasing
< ok 1 0
> exec
| CREATE DATABASE nrm_float_double_aliasing
< ok 1 0
> exec
| USE nrm_float_double_aliasing
< ok 0 0
> exec
| CREATE TABLE t_fd (f FLOAT, g FLOAT(10,2), h DOUBLE, i DOUBLE(15,4), j REAL, m FLOAT(5))
< ok 0 0
> query
| SHOW CREATE TABLE nrm_float_double_aliasing.t_fd
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_fd
: bytes
| CREATE TABLE `t_fd` (
|   `f` float DEFAULT NULL,
|   `g` float(10,2) DEFAULT NULL,
|   `h` double DEFAULT NULL,
|   `i` double(15,4) DEFAULT NULL,
|   `j` double DEFAULT NULL,
|   `m` float DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_float_double_aliasing
< ok 1 0
> exec
| CREATE DATABASE nrm_float_double_aliasing
< ok 1 0
> exec
| USE nrm_float_double_aliasing
< ok 0 0
> exec
| CREATE TABLE t_fd (f FLOAT, g FLOAT(10,2), h DOUBLE, i DOUBLE(15,4), j REAL, m FLOAT(5))
< ok 0 0
> query
| SHOW CREATE TABLE nrm_float_double_aliasing.t_fd
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_fd
: bytes
| CREATE TABLE `t_fd` (
|   `f` float DEFAULT NULL,
|   `g` float(10,2) DEFAULT NULL,
|   `h` double DEFAULT NULL,
|   `i` double(15,4) DEFAULT NULL,
|   `j` double DEFAULT NULL,
|   `m` float DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_char_binary_length_default
< ok 0 0
> exec
| CREATE DATABASE nrm_char_binary_length_default
< ok 1 0
> exec
| USE nrm_char_binary_length_default
< ok 0 0
> exec
| CREATE TABLE t_char (a CHAR, b CHAR(10), d BINARY, e BINARY(16), f VARBINARY(32))
< ok 0 0
> query
| SHOW CREATE TABLE nrm_char_binary_length_default.t_char
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_char
: bytes
| CREATE TABLE `t_char` (
|   `a` char(1) DEFAULT NULL,
|   `b` char(10) DEFAULT NULL,
|   `d` binary(1) DEFAULT NULL,
|   `e` binary(16) DEFAULT NULL,
|   `f` varbinary(32) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_char_binary_length_default
< ok 1 0
> exec
| CREATE DATABASE nrm_char_binary_length_default
< ok 1 0
> exec
| USE nrm_char_binary_length_default
< ok 0 0
> exec
| CREATE TABLE t_char (a CHAR, b CHAR(10), d BINARY, e BINARY(16), f VARBINARY(32))
< ok 0 0
> query
| SHOW CREATE TABLE nrm_char_binary_length_default.t_char
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_char
: bytes
| CREATE TABLE `t_char` (
|   `a` char(1) DEFAULT NULL,
|   `b` char(10) DEFAULT NULL,
|   `d` binary(1) DEFAULT NULL,
|   `e` binary(16) DEFAULT NULL,
|   `f` varbinary(32) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_char_binary_length_default
< ok 1 0
> exec
| CREATE DATABASE nrm_char_binary_length_default
< ok 1 0
> exec
| USE nrm_char_binary_length_default
< ok 0 0
> exec
| CREATE TABLE t_char (a CHAR, b CHAR(10), d BINARY, e BINARY(16), f VARBINARY(32))
< ok 0 0
> query
| SHOW CREATE TABLE nrm_char_binary_length_default.t_char
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_char
: bytes
| CREATE TABLE `t_char` (
|   `a` char(1) DEFAULT NULL,
|   `b` char(10) DEFAULT NULL,
|   `d` binary(1) DEFAULT NULL,
|   `e` binary(16) DEFAULT NULL,
|   `f` varbinary(32) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_char_binary_length_default
< ok 1 0
> exec
| CREATE DATABASE nrm_char_binary_length_default
< ok 1 0
> exec
| USE nrm_char_binary_length_default
< ok 0 0
> exec
| CREATE TABLE t_char (a CHAR, b CHAR(10), d BINARY, e BINARY(16), f VARBINARY(32))
< ok 0 0
> query
| SHOW CREATE TABLE nrm_char_binary_length_default.t_char
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_char
: bytes
| CREATE TABLE `t_char` (
|   `a` char(1) DEFAULT NULL,
|   `b` char(10) DEFAULT NULL,
|   `d` binary(1) DEFAULT NULL,
|   `e` binary(16) DEFAULT NULL,
|   `f` varbinary(32) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_char_binary_length_default
< ok 1 0
> exec
| CREATE DATABASE nrm_char_binary_length_default
< ok 1 0
> exec
| USE nrm_char_binary_length_default
< ok 0 0
> exec
| CREATE TABLE t_char (a CHAR, b CHAR(10), d BINARY, e BINARY(16), f VARBINARY(32))
< ok 0 0
> query
| SHOW CREATE TABLE nrm_char_binary_length_default.t_char
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_char
: bytes
| CREATE TABLE `t_char` (
|   `a` char(1) DEFAULT NULL,
|   `b` char(10) DEFAULT NULL,
|   `d` binary(1) DEFAULT NULL,
|   `e` binary(16) DEFAULT NULL,
|   `f` varbinary(32) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_char_binary_attribute_utf8mb4
< ok 0 0
> exec
| CREATE DATABASE nrm_char_binary_attribute_utf8mb4
< ok 1 0
> exec
| USE nrm_char_binary_attribute_utf8mb4
< ok 0 0
> exec
| CREATE TABLE t_binu (k VARCHAR(128) BINARY NOT NULL, c CHAR(5) BINARY, t TEXT BINARY, ti TINYTEXT BINARY, me MEDIUMTEXT BINARY, lo LONGTEXT BINARY) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE nrm_char_binary_attribute_utf8mb4.t_binu
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_binu
: bytes
| CREATE TABLE `t_binu` (
|   `k` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
|   `c` char(5) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   `t` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin,
|   `ti` tinytext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin,
|   `me` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin,
|   `lo` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_char_binary_attribute_utf8mb4
< ok 1 0
> exec
| CREATE DATABASE nrm_char_binary_attribute_utf8mb4
< ok 1 0
> exec
| USE nrm_char_binary_attribute_utf8mb4
< ok 0 0
> exec
| CREATE TABLE t_binu (k VARCHAR(128) BINARY NOT NULL, c CHAR(5) BINARY, t TEXT BINARY, ti TINYTEXT BINARY, me MEDIUMTEXT BINARY, lo LONGTEXT BINARY) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE nrm_char_binary_attribute_utf8mb4.t_binu
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_binu
: bytes
| CREATE TABLE `t_binu` (
|   `k` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
|   `c` char(5) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   `t` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin,
|   `ti` tinytext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin,
|   `me` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin,
|   `lo` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_char_binary_attribute_utf8mb4
< ok 1 0
> exec
| CREATE DATABASE nrm_char_binary_attribute_utf8mb4
< ok 1 0
> exec
| USE nrm_char_binary_attribute_utf8mb4
< ok 0 0
> exec
| CREATE TABLE t_binu (k VARCHAR(128) BINARY NOT NULL, c CHAR(5) BINARY, t TEXT BINARY, ti TINYTEXT BINARY, me MEDIUMTEXT BINARY, lo LONGTEXT BINARY) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE nrm_char_binary_attribute_utf8mb4.t_binu
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_binu
: bytes
| CREATE TABLE `t_binu` (
|   `k` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
|   `c` char(5) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   `t` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin,
|   `ti` tinytext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin,
|   `me` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin,
|   `lo` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_char_binary_attribute_utf8mb4
< ok 1 0
> exec
| CREATE DATABASE nrm_char_binary_attribute_utf8mb4
< ok 1 0
> exec
| USE nrm_char_binary_attribute_utf8mb4
< ok 0 0
> exec
| CREATE TABLE t_binu (k VARCHAR(128) BINARY NOT NULL, c CHAR(5) BINARY, t TEXT BINARY, ti TINYTEXT BINARY, me MEDIUMTEXT BINARY, lo LONGTEXT BINARY) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE nrm_char_binary_attribute_utf8mb4.t_binu
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_binu
: bytes
| CREATE TABLE `t_binu` (
|   `k` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
|   `c` char(5) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   `t` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin,
|   `ti` tinytext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin,
|   `me` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin,
|   `lo` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_char_binary_attribute_utf8mb4
< ok 1 0
> exec
| CREATE DATABASE nrm_char_binary_attribute_utf8mb4
< ok 1 0
> exec
| USE nrm_char_binary_attribute_utf8mb4
< ok 0 0
> exec
| CREATE TABLE t_binu (k VARCHAR(128) BINARY NOT NULL, c CHAR(5) BINARY, t TEXT BINARY, ti TINYTEXT BINARY, me MEDIUMTEXT BINARY, lo LONGTEXT BINARY) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE nrm_char_binary_attribute_utf8mb4.t_binu
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_binu
: bytes
| CREATE TABLE `t_binu` (
|   `k` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
|   `c` char(5) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   `t` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin,
|   `ti` tinytext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin,
|   `me` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin,
|   `lo` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_char_binary_attribute_utf8mb4
< ok 1 0
> exec
| CREATE DATABASE nrm_char_binary_attribute_utf8mb4
< ok 1 0
> exec
| USE nrm_char_binary_attribute_utf8mb4
< ok 0 0
> exec
| CREATE TABLE t_binu (k VARCHAR(128) BINARY NOT NULL, c CHAR(5) BINARY, t TEXT BINARY, ti TINYTEXT BINARY, me MEDIUMTEXT BINARY, lo LONGTEXT BINARY) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE nrm_char_binary_attribute_utf8mb4.t_binu
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_binu
: bytes
| CREATE TABLE `t_binu` (
|   `k` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
|   `c` char(5) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   `t` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin,
|   `ti` tinytext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin,
|   `me` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin,
|   `lo` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_char_binary_attribute_latin1
< ok 0 0
> exec
| CREATE DATABASE nrm_char_binary_attribute_latin1
< ok 1 0
> exec
| USE nrm_char_binary_attribute_latin1
< ok 0 0
> exec
| CREATE TABLE t_binl (a VARCHAR(50) BINARY, c CHAR(8) BINARY) DEFAULT CHARSET=latin1
< ok 0 0
> query
| SHOW CREATE TABLE nrm_char_binary_attribute_latin1.t_binl
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_binl
: bytes
| CREATE TABLE `t_binl` (
|   `a` varchar(50) CHARACTER SET latin1 COLLATE latin1_bin DEFAULT NULL,
|   `c` char(8) CHARACTER SET latin1 COLLATE latin1_bin DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=latin1
> exec
| DROP DATABASE IF EXISTS nrm_char_binary_attribute_latin1
< ok 1 0
> exec
| CREATE DATABASE nrm_char_binary_attribute_latin1
< ok 1 0
> exec
| USE nrm_char_binary_attribute_latin1
< ok 0 0
> exec
| CREATE TABLE t_binl (a VARCHAR(50) BINARY, c CHAR(8) BINARY) DEFAULT CHARSET=latin1
< ok 0 0
> query
| SHOW CREATE TABLE nrm_char_binary_attribute_latin1.t_binl
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_binl
: bytes
| CREATE TABLE `t_binl` (
|   `a` varchar(50) CHARACTER SET latin1 COLLATE latin1_bin DEFAULT NULL,
|   `c` char(8) CHARACTER SET latin1 COLLATE latin1_bin DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=latin1
> exec
| DROP DATABASE IF EXISTS nrm_char_binary_attribute_ascii_explicit
< ok 0 0
> exec
| CREATE DATABASE nrm_char_binary_attribute_ascii_explicit
< ok 1 0
> exec
| USE nrm_char_binary_attribute_ascii_explicit
< ok 0 0
> exec
| CREATE TABLE t_bina (a VARCHAR(10) CHARACTER SET ascii BINARY, c CHAR(4) CHARACTER SET ascii BINARY) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE nrm_char_binary_attribute_ascii_explicit.t_bina
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_bina
: bytes
| CREATE TABLE `t_bina` (
|   `a` varchar(10) CHARACTER SET ascii COLLATE ascii_bin DEFAULT NULL,
|   `c` char(4) CHARACTER SET ascii COLLATE ascii_bin DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_char_binary_attribute_ascii_explicit
< ok 1 0
> exec
| CREATE DATABASE nrm_char_binary_attribute_ascii_explicit
< ok 1 0
> exec
| USE nrm_char_binary_attribute_ascii_explicit
< ok 0 0
> exec
| CREATE TABLE t_bina (a VARCHAR(10) CHARACTER SET ascii BINARY, c CHAR(4) CHARACTER SET ascii BINARY) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE nrm_char_binary_attribute_ascii_explicit.t_bina
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_bina
: bytes
| CREATE TABLE `t_bina` (
|   `a` varchar(10) CHARACTER SET ascii COLLATE ascii_bin DEFAULT NULL,
|   `c` char(4) CHARACTER SET ascii COLLATE ascii_bin DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_char_binary_attribute_utf8_alias
< ok 0 0
> exec
| CREATE DATABASE nrm_char_binary_attribute_utf8_alias
< ok 1 0
> exec
| USE nrm_char_binary_attribute_utf8_alias
< ok 0 0
> exec
| CREATE TABLE t_binm (a VARCHAR(10) CHARACTER SET utf8 BINARY, b VARCHAR(10) CHARACTER SET utf8mb3 BINARY) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE nrm_char_binary_attribute_utf8_alias.t_binm
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_binm
: bytes
| CREATE TABLE `t_binm` (
|   `a` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_bin DEFAULT NULL,
|   `b` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_bin DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_char_binary_attribute_utf8_alias
< ok 1 0
> exec
| CREATE DATABASE nrm_char_binary_attribute_utf8_alias
< ok 1 0
> exec
| USE nrm_char_binary_attribute_utf8_alias
< ok 0 0
> exec
| CREATE TABLE t_binm (a VARCHAR(10) CHARACTER SET utf8 BINARY, b VARCHAR(10) CHARACTER SET utf8mb3 BINARY) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE nrm_char_binary_attribute_utf8_alias.t_binm
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_binm
: bytes
| CREATE TABLE `t_binm` (
|   `a` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_bin DEFAULT NULL,
|   `b` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_bin DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_char_binary_attribute_collate_precedence
< ok 0 0
> exec
| CREATE DATABASE nrm_char_binary_attribute_collate_precedence
< ok 1 0
> exec
| USE nrm_char_binary_attribute_collate_precedence
< ok 0 0
> exec
| CREATE TABLE t_binp (a VARCHAR(10) BINARY COLLATE utf8mb4_unicode_ci) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE nrm_char_binary_attribute_collate_precedence.t_binp
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_binp
: bytes
| CREATE TABLE `t_binp` (
|   `a` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_char_binary_attribute_real_binary_type
< ok 0 0
> exec
| CREATE DATABASE nrm_char_binary_attribute_real_binary_type
< ok 1 0
> exec
| USE nrm_char_binary_attribute_real_binary_type
< ok 0 0
> exec
| CREATE TABLE t_binr (a BINARY(8), b VARBINARY(16), c BLOB) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE nrm_char_binary_attribute_real_binary_type.t_binr
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_binr
: bytes
| CREATE TABLE `t_binr` (
|   `a` binary(8) DEFAULT NULL,
|   `b` varbinary(16) DEFAULT NULL,
|   `c` blob
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_char_binary_attribute_real_binary_type
< ok 1 0
> exec
| CREATE DATABASE nrm_char_binary_attribute_real_binary_type
< ok 1 0
> exec
| USE nrm_char_binary_attribute_real_binary_type
< ok 0 0
> exec
| CREATE TABLE t_binr (a BINARY(8), b VARBINARY(16), c BLOB) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE nrm_char_binary_attribute_real_binary_type.t_binr
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_binr
: bytes
| CREATE TABLE `t_binr` (
|   `a` binary(8) DEFAULT NULL,
|   `b` varbinary(16) DEFAULT NULL,
|   `c` blob
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_char_binary_attribute_real_binary_type
< ok 1 0
> exec
| CREATE DATABASE nrm_char_binary_attribute_real_binary_type
< ok 1 0
> exec
| USE nrm_char_binary_attribute_real_binary_type
< ok 0 0
> exec
| CREATE TABLE t_binr (a BINARY(8), b VARBINARY(16), c BLOB) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE nrm_char_binary_attribute_real_binary_type.t_binr
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_binr
: bytes
| CREATE TABLE `t_binr` (
|   `a` binary(8) DEFAULT NULL,
|   `b` varbinary(16) DEFAULT NULL,
|   `c` blob
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_text_blob_no_default_null
< ok 0 0
> exec
| CREATE DATABASE nrm_text_blob_no_default_null
< ok 1 0
> exec
| USE nrm_text_blob_no_default_null
< ok 0 0
> exec
| CREATE TABLE t_tb (g TEXT, h BLOB, j LONGTEXT, k TINYTEXT)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_text_blob_no_default_null.t_tb
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_tb
: bytes
| CREATE TABLE `t_tb` (
|   `g` text,
|   `h` blob,
|   `j` longtext,
|   `k` tinytext
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_text_blob_no_default_null
< ok 1 0
> exec
| CREATE DATABASE nrm_text_blob_no_default_null
< ok 1 0
> exec
| USE nrm_text_blob_no_default_null
< ok 0 0
> exec
| CREATE TABLE t_tb (g TEXT, h BLOB, j LONGTEXT, k TINYTEXT)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_text_blob_no_default_null.t_tb
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_tb
: bytes
| CREATE TABLE `t_tb` (
|   `g` text,
|   `h` blob,
|   `j` longtext,
|   `k` tinytext
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_text_blob_no_default_null
< ok 1 0
> exec
| CREATE DATABASE nrm_text_blob_no_default_null
< ok 1 0
> exec
| USE nrm_text_blob_no_default_null
< ok 0 0
> exec
| CREATE TABLE t_tb (g TEXT, h BLOB, j LONGTEXT, k TINYTEXT)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_text_blob_no_default_null.t_tb
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_tb
: bytes
| CREATE TABLE `t_tb` (
|   `g` text,
|   `h` blob,
|   `j` longtext,
|   `k` tinytext
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_text_blob_no_default_null
< ok 1 0
> exec
| CREATE DATABASE nrm_text_blob_no_default_null
< ok 1 0
> exec
| USE nrm_text_blob_no_default_null
< ok 0 0
> exec
| CREATE TABLE t_tb (g TEXT, h BLOB, j LONGTEXT, k TINYTEXT)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_text_blob_no_default_null.t_tb
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_tb
: bytes
| CREATE TABLE `t_tb` (
|   `g` text,
|   `h` blob,
|   `j` longtext,
|   `k` tinytext
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_year_width
< ok 0 0
> exec
| CREATE DATABASE nrm_year_width
< ok 1 0
> exec
| USE nrm_year_width
< ok 0 0
> exec
| CREATE TABLE t_year (a YEAR, b YEAR(4))
< ok 0 0
> query
| SHOW CREATE TABLE nrm_year_width.t_year
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_year
: bytes
| CREATE TABLE `t_year` (
|   `a` year DEFAULT NULL,
|   `b` year DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_year_width
< ok 1 0
> exec
| CREATE DATABASE nrm_year_width
< ok 1 0
> exec
| USE nrm_year_width
< ok 0 0
> exec
| CREATE TABLE t_year (a YEAR, b YEAR(4))
< ok 0 0
> query
| SHOW CREATE TABLE nrm_year_width.t_year
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_year
: bytes
| CREATE TABLE `t_year` (
|   `a` year DEFAULT NULL,
|   `b` year DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_bit_length_default
< ok 0 0
> exec
| CREATE DATABASE nrm_bit_length_default
< ok 1 0
> exec
| USE nrm_bit_length_default
< ok 0 0
> exec
| CREATE TABLE t_bit (c BIT, d BIT(8))
< ok 0 0
> query
| SHOW CREATE TABLE nrm_bit_length_default.t_bit
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_bit
: bytes
| CREATE TABLE `t_bit` (
|   `c` bit(1) DEFAULT NULL,
|   `d` bit(8) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_bit_length_default
< ok 1 0
> exec
| CREATE DATABASE nrm_bit_length_default
< ok 1 0
> exec
| USE nrm_bit_length_default
< ok 0 0
> exec
| CREATE TABLE t_bit (c BIT, d BIT(8))
< ok 0 0
> query
| SHOW CREATE TABLE nrm_bit_length_default.t_bit
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_bit
: bytes
| CREATE TABLE `t_bit` (
|   `c` bit(1) DEFAULT NULL,
|   `d` bit(8) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_bit_default_bitliteral
< ok 0 0
> exec
| CREATE DATABASE nrm_bit_default_bitliteral
< ok 1 0
> exec
| USE nrm_bit_default_bitliteral
< ok 0 0
> exec
| CREATE TABLE t_bitd (a BIT(8) NOT NULL DEFAULT b'0', b BIT(8) NOT NULL DEFAULT b'101', c BIT(8) NOT NULL DEFAULT b'00000101')
< ok 0 0
> query
| SHOW CREATE TABLE nrm_bit_default_bitliteral.t_bitd
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_bitd
: bytes
| CREATE TABLE `t_bitd` (
|   `a` bit(8) NOT NULL DEFAULT b'0',
|   `b` bit(8) NOT NULL DEFAULT b'101',
|   `c` bit(8) NOT NULL DEFAULT b'101'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_bit_default_bitliteral
< ok 1 0
> exec
| CREATE DATABASE nrm_bit_default_bitliteral
< ok 1 0
> exec
| USE nrm_bit_default_bitliteral
< ok 0 0
> exec
| CREATE TABLE t_bitd (a BIT(8) NOT NULL DEFAULT b'0', b BIT(8) NOT NULL DEFAULT b'101', c BIT(8) NOT NULL DEFAULT b'00000101')
< ok 0 0
> query
| SHOW CREATE TABLE nrm_bit_default_bitliteral.t_bitd
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_bitd
: bytes
| CREATE TABLE `t_bitd` (
|   `a` bit(8) NOT NULL DEFAULT b'0',
|   `b` bit(8) NOT NULL DEFAULT b'101',
|   `c` bit(8) NOT NULL DEFAULT b'101'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_bit_default_bitliteral
< ok 1 0
> exec
| CREATE DATABASE nrm_bit_default_bitliteral
< ok 1 0
> exec
| USE nrm_bit_default_bitliteral
< ok 0 0
> exec
| CREATE TABLE t_bitd (a BIT(8) NOT NULL DEFAULT b'0', b BIT(8) NOT NULL DEFAULT b'101', c BIT(8) NOT NULL DEFAULT b'00000101')
< ok 0 0
> query
| SHOW CREATE TABLE nrm_bit_default_bitliteral.t_bitd
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_bitd
: bytes
| CREATE TABLE `t_bitd` (
|   `a` bit(8) NOT NULL DEFAULT b'0',
|   `b` bit(8) NOT NULL DEFAULT b'101',
|   `c` bit(8) NOT NULL DEFAULT b'101'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_bit_default_hex
< ok 0 0
> exec
| CREATE DATABASE nrm_bit_default_hex
< ok 1 0
> exec
| USE nrm_bit_default_hex
< ok 0 0
> exec
| CREATE TABLE t_bithex (a BIT(8) NOT NULL DEFAULT 0x05, b BIT(8) NOT NULL DEFAULT x'05', c BIT(16) NOT NULL DEFAULT 0xABCD)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_bit_default_hex.t_bithex
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_bithex
: bytes
| CREATE TABLE `t_bithex` (
|   `a` bit(8) NOT NULL DEFAULT b'101',
|   `b` bit(8) NOT NULL DEFAULT b'101',
|   `c` bit(16) NOT NULL DEFAULT b'1010101111001101'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_bit_default_hex
< ok 1 0
> exec
| CREATE DATABASE nrm_bit_default_hex
< ok 1 0
> exec
| USE nrm_bit_default_hex
< ok 0 0
> exec
| CREATE TABLE t_bithex (a BIT(8) NOT NULL DEFAULT 0x05, b BIT(8) NOT NULL DEFAULT x'05', c BIT(16) NOT NULL DEFAULT 0xABCD)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_bit_default_hex.t_bithex
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_bithex
: bytes
| CREATE TABLE `t_bithex` (
|   `a` bit(8) NOT NULL DEFAULT b'101',
|   `b` bit(8) NOT NULL DEFAULT b'101',
|   `c` bit(16) NOT NULL DEFAULT b'1010101111001101'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_bit_default_hex
< ok 1 0
> exec
| CREATE DATABASE nrm_bit_default_hex
< ok 1 0
> exec
| USE nrm_bit_default_hex
< ok 0 0
> exec
| CREATE TABLE t_bithex (a BIT(8) NOT NULL DEFAULT 0x05, b BIT(8) NOT NULL DEFAULT x'05', c BIT(16) NOT NULL DEFAULT 0xABCD)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_bit_default_hex.t_bithex
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_bithex
: bytes
| CREATE TABLE `t_bithex` (
|   `a` bit(8) NOT NULL DEFAULT b'101',
|   `b` bit(8) NOT NULL DEFAULT b'101',
|   `c` bit(16) NOT NULL DEFAULT b'1010101111001101'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_bit_default_numeric_and_string
< ok 0 0
> exec
| CREATE DATABASE nrm_bit_default_numeric_and_string
< ok 1 0
> exec
| USE nrm_bit_default_numeric_and_string
< ok 0 0
> exec
| CREATE TABLE t_bitns (a BIT(8) NOT NULL DEFAULT 5, b BIT(8) NOT NULL DEFAULT 0b101, c BIT(16) NOT NULL DEFAULT 'AB')
< ok 0 0
> query
| SHOW CREATE TABLE nrm_bit_default_numeric_and_string.t_bitns
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_bitns
: bytes
| CREATE TABLE `t_bitns` (
|   `a` bit(8) NOT NULL DEFAULT b'101',
|   `b` bit(8) NOT NULL DEFAULT b'101',
|   `c` bit(16) NOT NULL DEFAULT b'100000101000010'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_bit_default_numeric_and_string
< ok 1 0
> exec
| CREATE DATABASE nrm_bit_default_numeric_and_string
< ok 1 0
> exec
| USE nrm_bit_default_numeric_and_string
< ok 0 0
> exec
| CREATE TABLE t_bitns (a BIT(8) NOT NULL DEFAULT 5, b BIT(8) NOT NULL DEFAULT 0b101, c BIT(16) NOT NULL DEFAULT 'AB')
< ok 0 0
> query
| SHOW CREATE TABLE nrm_bit_default_numeric_and_string.t_bitns
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_bitns
: bytes
| CREATE TABLE `t_bitns` (
|   `a` bit(8) NOT NULL DEFAULT b'101',
|   `b` bit(8) NOT NULL DEFAULT b'101',
|   `c` bit(16) NOT NULL DEFAULT b'100000101000010'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_bit_default_numeric_and_string
< ok 1 0
> exec
| CREATE DATABASE nrm_bit_default_numeric_and_string
< ok 1 0
> exec
| USE nrm_bit_default_numeric_and_string
< ok 0 0
> exec
| CREATE TABLE t_bitns (a BIT(8) NOT NULL DEFAULT 5, b BIT(8) NOT NULL DEFAULT 0b101, c BIT(16) NOT NULL DEFAULT 'AB')
< ok 0 0
> query
| SHOW CREATE TABLE nrm_bit_default_numeric_and_string.t_bitns
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_bitns
: bytes
| CREATE TABLE `t_bitns` (
|   `a` bit(8) NOT NULL DEFAULT b'101',
|   `b` bit(8) NOT NULL DEFAULT b'101',
|   `c` bit(16) NOT NULL DEFAULT b'100000101000010'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_year_default_numeric
< ok 0 0
> exec
| CREATE DATABASE nrm_year_default_numeric
< ok 1 0
> exec
| USE nrm_year_default_numeric
< ok 0 0
> exec
| CREATE TABLE t_yeard (a YEAR NOT NULL DEFAULT 2000, b YEAR NOT NULL DEFAULT 1999, c YEAR NOT NULL DEFAULT 0, d YEAR NOT NULL DEFAULT 99, e YEAR NOT NULL DEFAULT 5)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_year_default_numeric.t_yeard
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_yeard
: bytes
| CREATE TABLE `t_yeard` (
|   `a` year NOT NULL DEFAULT '2000',
|   `b` year NOT NULL DEFAULT '1999',
|   `c` year NOT NULL DEFAULT '0000',
|   `d` year NOT NULL DEFAULT '1999',
|   `e` year NOT NULL DEFAULT '2005'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_year_default_numeric
< ok 1 0
> exec
| CREATE DATABASE nrm_year_default_numeric
< ok 1 0
> exec
| USE nrm_year_default_numeric
< ok 0 0
> exec
| CREATE TABLE t_yeard (a YEAR NOT NULL DEFAULT 2000, b YEAR NOT NULL DEFAULT 1999, c YEAR NOT NULL DEFAULT 0, d YEAR NOT NULL DEFAULT 99, e YEAR NOT NULL DEFAULT 5)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_year_default_numeric.t_yeard
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_yeard
: bytes
| CREATE TABLE `t_yeard` (
|   `a` year NOT NULL DEFAULT '2000',
|   `b` year NOT NULL DEFAULT '1999',
|   `c` year NOT NULL DEFAULT '0000',
|   `d` year NOT NULL DEFAULT '1999',
|   `e` year NOT NULL DEFAULT '2005'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_year_default_numeric
< ok 1 0
> exec
| CREATE DATABASE nrm_year_default_numeric
< ok 1 0
> exec
| USE nrm_year_default_numeric
< ok 0 0
> exec
| CREATE TABLE t_yeard (a YEAR NOT NULL DEFAULT 2000, b YEAR NOT NULL DEFAULT 1999, c YEAR NOT NULL DEFAULT 0, d YEAR NOT NULL DEFAULT 99, e YEAR NOT NULL DEFAULT 5)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_year_default_numeric.t_yeard
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_yeard
: bytes
| CREATE TABLE `t_yeard` (
|   `a` year NOT NULL DEFAULT '2000',
|   `b` year NOT NULL DEFAULT '1999',
|   `c` year NOT NULL DEFAULT '0000',
|   `d` year NOT NULL DEFAULT '1999',
|   `e` year NOT NULL DEFAULT '2005'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_year_default_numeric
< ok 1 0
> exec
| CREATE DATABASE nrm_year_default_numeric
< ok 1 0
> exec
| USE nrm_year_default_numeric
< ok 0 0
> exec
| CREATE TABLE t_yeard (a YEAR NOT NULL DEFAULT 2000, b YEAR NOT NULL DEFAULT 1999, c YEAR NOT NULL DEFAULT 0, d YEAR NOT NULL DEFAULT 99, e YEAR NOT NULL DEFAULT 5)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_year_default_numeric.t_yeard
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_yeard
: bytes
| CREATE TABLE `t_yeard` (
|   `a` year NOT NULL DEFAULT '2000',
|   `b` year NOT NULL DEFAULT '1999',
|   `c` year NOT NULL DEFAULT '0000',
|   `d` year NOT NULL DEFAULT '1999',
|   `e` year NOT NULL DEFAULT '2005'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_year_default_numeric
< ok 1 0
> exec
| CREATE DATABASE nrm_year_default_numeric
< ok 1 0
> exec
| USE nrm_year_default_numeric
< ok 0 0
> exec
| CREATE TABLE t_yeard (a YEAR NOT NULL DEFAULT 2000, b YEAR NOT NULL DEFAULT 1999, c YEAR NOT NULL DEFAULT 0, d YEAR NOT NULL DEFAULT 99, e YEAR NOT NULL DEFAULT 5)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_year_default_numeric.t_yeard
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_yeard
: bytes
| CREATE TABLE `t_yeard` (
|   `a` year NOT NULL DEFAULT '2000',
|   `b` year NOT NULL DEFAULT '1999',
|   `c` year NOT NULL DEFAULT '0000',
|   `d` year NOT NULL DEFAULT '1999',
|   `e` year NOT NULL DEFAULT '2005'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_bigint_unsigned_max_default
< ok 0 0
> exec
| CREATE DATABASE nrm_bigint_unsigned_max_default
< ok 1 0
> exec
| USE nrm_bigint_unsigned_max_default
< ok 0 0
> exec
| CREATE TABLE t_bigmax (big BIGINT UNSIGNED NOT NULL DEFAULT 18446744073709551615, mid BIGINT UNSIGNED NOT NULL DEFAULT 9223372036854775808)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_bigint_unsigned_max_default.t_bigmax
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_bigmax
: bytes
| CREATE TABLE `t_bigmax` (
|   `big` bigint unsigned NOT NULL DEFAULT '18446744073709551615',
|   `mid` bigint unsigned NOT NULL DEFAULT '9223372036854775808'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_bigint_unsigned_max_default
< ok 1 0
> exec
| CREATE DATABASE nrm_bigint_unsigned_max_default
< ok 1 0
> exec
| USE nrm_bigint_unsigned_max_default
< ok 0 0
> exec
| CREATE TABLE t_bigmax (big BIGINT UNSIGNED NOT NULL DEFAULT 18446744073709551615, mid BIGINT UNSIGNED NOT NULL DEFAULT 9223372036854775808)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_bigint_unsigned_max_default.t_bigmax
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_bigmax
: bytes
| CREATE TABLE `t_bigmax` (
|   `big` bigint unsigned NOT NULL DEFAULT '18446744073709551615',
|   `mid` bigint unsigned NOT NULL DEFAULT '9223372036854775808'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_bigint_signed_boundary_default
< ok 0 0
> exec
| CREATE DATABASE nrm_bigint_signed_boundary_default
< ok 1 0
> exec
| USE nrm_bigint_signed_boundary_default
< ok 0 0
> exec
| CREATE TABLE t_bigbnd (lo BIGINT NOT NULL DEFAULT -9223372036854775808, hi BIGINT NOT NULL DEFAULT 9223372036854775807)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_bigint_signed_boundary_default.t_bigbnd
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_bigbnd
: bytes
| CREATE TABLE `t_bigbnd` (
|   `lo` bigint NOT NULL DEFAULT '-9223372036854775808',
|   `hi` bigint NOT NULL DEFAULT '9223372036854775807'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_bigint_signed_boundary_default
< ok 1 0
> exec
| CREATE DATABASE nrm_bigint_signed_boundary_default
< ok 1 0
> exec
| USE nrm_bigint_signed_boundary_default
< ok 0 0
> exec
| CREATE TABLE t_bigbnd (lo BIGINT NOT NULL DEFAULT -9223372036854775808, hi BIGINT NOT NULL DEFAULT 9223372036854775807)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_bigint_signed_boundary_default.t_bigbnd
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_bigbnd
: bytes
| CREATE TABLE `t_bigbnd` (
|   `lo` bigint NOT NULL DEFAULT '-9223372036854775808',
|   `hi` bigint NOT NULL DEFAULT '9223372036854775807'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_time_json_date_stable
< ok 0 0
> exec
| CREATE DATABASE nrm_time_json_date_stable
< ok 1 0
> exec
| USE nrm_time_json_date_stable
< ok 0 0
> exec
| CREATE TABLE t_tjd (e JSON, f DATE, g TIME, h TIME(3))
< ok 0 0
> query
| SHOW CREATE TABLE nrm_time_json_date_stable.t_tjd
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_tjd
: bytes
| CREATE TABLE `t_tjd` (
|   `e` json DEFAULT NULL,
|   `f` date DEFAULT NULL,
|   `g` time DEFAULT NULL,
|   `h` time(3) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_time_json_date_stable
< ok 1 0
> exec
| CREATE DATABASE nrm_time_json_date_stable
< ok 1 0
> exec
| USE nrm_time_json_date_stable
< ok 0 0
> exec
| CREATE TABLE t_tjd (e JSON, f DATE, g TIME, h TIME(3))
< ok 0 0
> query
| SHOW CREATE TABLE nrm_time_json_date_stable.t_tjd
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_tjd
: bytes
| CREATE TABLE `t_tjd` (
|   `e` json DEFAULT NULL,
|   `f` date DEFAULT NULL,
|   `g` time DEFAULT NULL,
|   `h` time(3) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_time_json_date_stable
< ok 1 0
> exec
| CREATE DATABASE nrm_time_json_date_stable
< ok 1 0
> exec
| USE nrm_time_json_date_stable
< ok 0 0
> exec
| CREATE TABLE t_tjd (e JSON, f DATE, g TIME, h TIME(3))
< ok 0 0
> query
| SHOW CREATE TABLE nrm_time_json_date_stable.t_tjd
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_tjd
: bytes
| CREATE TABLE `t_tjd` (
|   `e` json DEFAULT NULL,
|   `f` date DEFAULT NULL,
|   `g` time DEFAULT NULL,
|   `h` time(3) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_time_json_date_stable
< ok 1 0
> exec
| CREATE DATABASE nrm_time_json_date_stable
< ok 1 0
> exec
| USE nrm_time_json_date_stable
< ok 0 0
> exec
| CREATE TABLE t_tjd (e JSON, f DATE, g TIME, h TIME(3))
< ok 0 0
> query
| SHOW CREATE TABLE nrm_time_json_date_stable.t_tjd
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_tjd
: bytes
| CREATE TABLE `t_tjd` (
|   `e` json DEFAULT NULL,
|   `f` date DEFAULT NULL,
|   `g` time DEFAULT NULL,
|   `h` time(3) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_default_literal_quoting
< ok 0 0
> exec
| CREATE DATABASE nrm_default_literal_quoting
< ok 1 0
> exec
| USE nrm_default_literal_quoting
< ok 0 0
> exec
| CREATE TABLE t_def (a INT DEFAULT 0, b INT DEFAULT 5, h INT DEFAULT '7', c VARCHAR(20) DEFAULT 'hello', d VARCHAR(20) DEFAULT '', e VARCHAR(20) DEFAULT 'it''s')
< ok 0 0
> query
| SHOW CREATE TABLE nrm_default_literal_quoting.t_def
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_def
: bytes
| CREATE TABLE `t_def` (
|   `a` int DEFAULT '0',
|   `b` int DEFAULT '5',
|   `h` int DEFAULT '7',
|   `c` varchar(20) DEFAULT 'hello',
|   `d` varchar(20) DEFAULT '',
|   `e` varchar(20) DEFAULT 'it''s'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_default_literal_quoting
< ok 1 0
> exec
| CREATE DATABASE nrm_default_literal_quoting
< ok 1 0
> exec
| USE nrm_default_literal_quoting
< ok 0 0
> exec
| CREATE TABLE t_def (a INT DEFAULT 0, b INT DEFAULT 5, h INT DEFAULT '7', c VARCHAR(20) DEFAULT 'hello', d VARCHAR(20) DEFAULT '', e VARCHAR(20) DEFAULT 'it''s')
< ok 0 0
> query
| SHOW CREATE TABLE nrm_default_literal_quoting.t_def
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_def
: bytes
| CREATE TABLE `t_def` (
|   `a` int DEFAULT '0',
|   `b` int DEFAULT '5',
|   `h` int DEFAULT '7',
|   `c` varchar(20) DEFAULT 'hello',
|   `d` varchar(20) DEFAULT '',
|   `e` varchar(20) DEFAULT 'it''s'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_default_literal_quoting
< ok 1 0
> exec
| CREATE DATABASE nrm_default_literal_quoting
< ok 1 0
> exec
| USE nrm_default_literal_quoting
< ok 0 0
> exec
| CREATE TABLE t_def (a INT DEFAULT 0, b INT DEFAULT 5, h INT DEFAULT '7', c VARCHAR(20) DEFAULT 'hello', d VARCHAR(20) DEFAULT '', e VARCHAR(20) DEFAULT 'it''s')
< ok 0 0
> query
| SHOW CREATE TABLE nrm_default_literal_quoting.t_def
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_def
: bytes
| CREATE TABLE `t_def` (
|   `a` int DEFAULT '0',
|   `b` int DEFAULT '5',
|   `h` int DEFAULT '7',
|   `c` varchar(20) DEFAULT 'hello',
|   `d` varchar(20) DEFAULT '',
|   `e` varchar(20) DEFAULT 'it''s'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_default_literal_quoting
< ok 1 0
> exec
| CREATE DATABASE nrm_default_literal_quoting
< ok 1 0
> exec
| USE nrm_default_literal_quoting
< ok 0 0
> exec
| CREATE TABLE t_def (a INT DEFAULT 0, b INT DEFAULT 5, h INT DEFAULT '7', c VARCHAR(20) DEFAULT 'hello', d VARCHAR(20) DEFAULT '', e VARCHAR(20) DEFAULT 'it''s')
< ok 0 0
> query
| SHOW CREATE TABLE nrm_default_literal_quoting.t_def
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_def
: bytes
| CREATE TABLE `t_def` (
|   `a` int DEFAULT '0',
|   `b` int DEFAULT '5',
|   `h` int DEFAULT '7',
|   `c` varchar(20) DEFAULT 'hello',
|   `d` varchar(20) DEFAULT '',
|   `e` varchar(20) DEFAULT 'it''s'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_default_literal_quoting
< ok 1 0
> exec
| CREATE DATABASE nrm_default_literal_quoting
< ok 1 0
> exec
| USE nrm_default_literal_quoting
< ok 0 0
> exec
| CREATE TABLE t_def (a INT DEFAULT 0, b INT DEFAULT 5, h INT DEFAULT '7', c VARCHAR(20) DEFAULT 'hello', d VARCHAR(20) DEFAULT '', e VARCHAR(20) DEFAULT 'it''s')
< ok 0 0
> query
| SHOW CREATE TABLE nrm_default_literal_quoting.t_def
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_def
: bytes
| CREATE TABLE `t_def` (
|   `a` int DEFAULT '0',
|   `b` int DEFAULT '5',
|   `h` int DEFAULT '7',
|   `c` varchar(20) DEFAULT 'hello',
|   `d` varchar(20) DEFAULT '',
|   `e` varchar(20) DEFAULT 'it''s'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_default_literal_quoting
< ok 1 0
> exec
| CREATE DATABASE nrm_default_literal_quoting
< ok 1 0
> exec
| USE nrm_default_literal_quoting
< ok 0 0
> exec
| CREATE TABLE t_def (a INT DEFAULT 0, b INT DEFAULT 5, h INT DEFAULT '7', c VARCHAR(20) DEFAULT 'hello', d VARCHAR(20) DEFAULT '', e VARCHAR(20) DEFAULT 'it''s')
< ok 0 0
> query
| SHOW CREATE TABLE nrm_default_literal_quoting.t_def
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_def
: bytes
| CREATE TABLE `t_def` (
|   `a` int DEFAULT '0',
|   `b` int DEFAULT '5',
|   `h` int DEFAULT '7',
|   `c` varchar(20) DEFAULT 'hello',
|   `d` varchar(20) DEFAULT '',
|   `e` varchar(20) DEFAULT 'it''s'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_decimal_default_padding
< ok 0 0
> exec
| CREATE DATABASE nrm_decimal_default_padding
< ok 1 0
> exec
| USE nrm_decimal_default_padding
< ok 0 0
> exec
| CREATE TABLE t_dpad (f DECIMAL(10,2) DEFAULT 0, g DECIMAL(10,2) DEFAULT 0.00)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_decimal_default_padding.t_dpad
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_dpad
: bytes
| CREATE TABLE `t_dpad` (
|   `f` decimal(10,2) DEFAULT '0.00',
|   `g` decimal(10,2) DEFAULT '0.00'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_decimal_default_padding
< ok 1 0
> exec
| CREATE DATABASE nrm_decimal_default_padding
< ok 1 0
> exec
| USE nrm_decimal_default_padding
< ok 0 0
> exec
| CREATE TABLE t_dpad (f DECIMAL(10,2) DEFAULT 0, g DECIMAL(10,2) DEFAULT 0.00)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_decimal_default_padding.t_dpad
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_dpad
: bytes
| CREATE TABLE `t_dpad` (
|   `f` decimal(10,2) DEFAULT '0.00',
|   `g` decimal(10,2) DEFAULT '0.00'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_boolean_default
< ok 0 0
> exec
| CREATE DATABASE nrm_boolean_default
< ok 1 0
> exec
| USE nrm_boolean_default
< ok 0 0
> exec
| CREATE TABLE t_bdef (j BOOLEAN DEFAULT TRUE, k BOOLEAN DEFAULT FALSE)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_boolean_default.t_bdef
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_bdef
: bytes
| CREATE TABLE `t_bdef` (
|   `j` tinyint(1) DEFAULT '1',
|   `k` tinyint(1) DEFAULT '0'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_boolean_default
< ok 1 0
> exec
| CREATE DATABASE nrm_boolean_default
< ok 1 0
> exec
| USE nrm_boolean_default
< ok 0 0
> exec
| CREATE TABLE t_bdef (j BOOLEAN DEFAULT TRUE, k BOOLEAN DEFAULT FALSE)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_boolean_default.t_bdef
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_bdef
: bytes
| CREATE TABLE `t_bdef` (
|   `j` tinyint(1) DEFAULT '1',
|   `k` tinyint(1) DEFAULT '0'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_nullable_default_null
< ok 0 0
> exec
| CREATE DATABASE nrm_nullable_default_null
< ok 1 0
> exec
| USE nrm_nullable_default_null
< ok 0 0
> exec
| CREATE TABLE t_null (a INT NULL, b INT DEFAULT NULL, c INT NULL DEFAULT NULL, d VARCHAR(10) NULL)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_nullable_default_null.t_null
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_null
: bytes
| CREATE TABLE `t_null` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   `c` int DEFAULT NULL,
|   `d` varchar(10) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_nullable_default_null
< ok 1 0
> exec
| CREATE DATABASE nrm_nullable_default_null
< ok 1 0
> exec
| USE nrm_nullable_default_null
< ok 0 0
> exec
| CREATE TABLE t_null (a INT NULL, b INT DEFAULT NULL, c INT NULL DEFAULT NULL, d VARCHAR(10) NULL)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_nullable_default_null.t_null
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_null
: bytes
| CREATE TABLE `t_null` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   `c` int DEFAULT NULL,
|   `d` varchar(10) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_nullable_default_null
< ok 1 0
> exec
| CREATE DATABASE nrm_nullable_default_null
< ok 1 0
> exec
| USE nrm_nullable_default_null
< ok 0 0
> exec
| CREATE TABLE t_null (a INT NULL, b INT DEFAULT NULL, c INT NULL DEFAULT NULL, d VARCHAR(10) NULL)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_nullable_default_null.t_null
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_null
: bytes
| CREATE TABLE `t_null` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   `c` int DEFAULT NULL,
|   `d` varchar(10) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_nullable_default_null
< ok 1 0
> exec
| CREATE DATABASE nrm_nullable_default_null
< ok 1 0
> exec
| USE nrm_nullable_default_null
< ok 0 0
> exec
| CREATE TABLE t_null (a INT NULL, b INT DEFAULT NULL, c INT NULL DEFAULT NULL, d VARCHAR(10) NULL)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_nullable_default_null.t_null
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_null
: bytes
| CREATE TABLE `t_null` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   `c` int DEFAULT NULL,
|   `d` varchar(10) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_enum_set_quoting
< ok 0 0
> exec
| CREATE DATABASE nrm_enum_set_quoting
< ok 1 0
> exec
| USE nrm_enum_set_quoting
< ok 0 0
> exec
| CREATE TABLE t_enum (a ENUM('x','y','z'), b SET('a','b','c'), e ENUM("dq1","dq2"))
< ok 0 0
> query
| SHOW CREATE TABLE nrm_enum_set_quoting.t_enum
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_enum
: bytes
| CREATE TABLE `t_enum` (
|   `a` enum('x','y','z') DEFAULT NULL,
|   `b` set('a','b','c') DEFAULT NULL,
|   `e` enum('dq1','dq2') DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_enum_set_quoting
< ok 1 0
> exec
| CREATE DATABASE nrm_enum_set_quoting
< ok 1 0
> exec
| USE nrm_enum_set_quoting
< ok 0 0
> exec
| CREATE TABLE t_enum (a ENUM('x','y','z'), b SET('a','b','c'), e ENUM("dq1","dq2"))
< ok 0 0
> query
| SHOW CREATE TABLE nrm_enum_set_quoting.t_enum
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_enum
: bytes
| CREATE TABLE `t_enum` (
|   `a` enum('x','y','z') DEFAULT NULL,
|   `b` set('a','b','c') DEFAULT NULL,
|   `e` enum('dq1','dq2') DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_enum_set_quoting
< ok 1 0
> exec
| CREATE DATABASE nrm_enum_set_quoting
< ok 1 0
> exec
| USE nrm_enum_set_quoting
< ok 0 0
> exec
| CREATE TABLE t_enum (a ENUM('x','y','z'), b SET('a','b','c'), e ENUM("dq1","dq2"))
< ok 0 0
> query
| SHOW CREATE TABLE nrm_enum_set_quoting.t_enum
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_enum
: bytes
| CREATE TABLE `t_enum` (
|   `a` enum('x','y','z') DEFAULT NULL,
|   `b` set('a','b','c') DEFAULT NULL,
|   `e` enum('dq1','dq2') DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_comment_escaping
< ok 0 0
> exec
| CREATE DATABASE nrm_comment_escaping
< ok 1 0
> exec
| USE nrm_comment_escaping
< ok 0 0
> exec
| CREATE TABLE t_comment (a INT COMMENT 'has ''quote'' inside')
< ok 0 0
> query
| SHOW CREATE TABLE nrm_comment_escaping.t_comment
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_comment
: bytes
| CREATE TABLE `t_comment` (
|   `a` int DEFAULT NULL COMMENT 'has ''quote'' inside'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_column_charset_echo_80
< ok 0 0
> exec
| CREATE DATABASE nrm_column_charset_echo_80
< ok 1 0
> exec
| USE nrm_column_charset_echo_80
< ok 0 0
> exec
| CREATE TABLE t_ce80 (a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci, c VARCHAR(10) CHARACTER SET utf8mb4, d VARCHAR(10)) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE nrm_column_charset_echo_80.t_ce80
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_ce80
: bytes
| CREATE TABLE `t_ce80` (
|   `a` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
|   `c` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
|   `d` varchar(10) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_column_charset_echo_80
< ok 1 0
> exec
| CREATE DATABASE nrm_column_charset_echo_80
< ok 1 0
> exec
| USE nrm_column_charset_echo_80
< ok 0 0
> exec
| CREATE TABLE t_ce80 (a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci, c VARCHAR(10) CHARACTER SET utf8mb4, d VARCHAR(10)) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE nrm_column_charset_echo_80.t_ce80
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_ce80
: bytes
| CREATE TABLE `t_ce80` (
|   `a` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
|   `c` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
|   `d` varchar(10) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_column_charset_echo_80
< ok 1 0
> exec
| CREATE DATABASE nrm_column_charset_echo_80
< ok 1 0
> exec
| USE nrm_column_charset_echo_80
< ok 0 0
> exec
| CREATE TABLE t_ce80 (a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci, c VARCHAR(10) CHARACTER SET utf8mb4, d VARCHAR(10)) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE nrm_column_charset_echo_80.t_ce80
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_ce80
: bytes
| CREATE TABLE `t_ce80` (
|   `a` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
|   `c` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
|   `d` varchar(10) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_column_charset_only_collation_resolution
< ok 0 0
> exec
| CREATE DATABASE nrm_column_charset_only_collation_resolution
< ok 1 0
> exec
| USE nrm_column_charset_only_collation_resolution
< ok 0 0
> exec
| CREATE TABLE t_d1 (a VARCHAR(10) CHARACTER SET utf8mb4) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
< ok 0 0
> query
| SHOW CREATE TABLE nrm_column_charset_only_collation_resolution.t_d1
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_d1
: bytes
| CREATE TABLE `t_d1` (
|   `a` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
> exec
| DROP DATABASE IF EXISTS nrm_generated_expr_normalization
< ok 0 0
> exec
| CREATE DATABASE nrm_generated_expr_normalization
< ok 1 0
> exec
| USE nrm_generated_expr_normalization
< ok 0 0
> exec
| CREATE TABLE t_gen (a INT, b INT GENERATED ALWAYS AS (a+1) VIRTUAL, c INT GENERATED ALWAYS AS ( a * 2 ) STORED, d INT AS (a+1))
< ok 0 0
> query
| SHOW CREATE TABLE nrm_generated_expr_normalization.t_gen
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_gen
: bytes
| CREATE TABLE `t_gen` (
|   `a` int DEFAULT NULL,
|   `b` int GENERATED ALWAYS AS ((`a` + 1)) VIRTUAL,
|   `c` int GENERATED ALWAYS AS ((`a` * 2)) STORED,
|   `d` int GENERATED ALWAYS AS ((`a` + 1)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_generated_expr_normalization
< ok 1 0
> exec
| CREATE DATABASE nrm_generated_expr_normalization
< ok 1 0
> exec
| USE nrm_generated_expr_normalization
< ok 0 0
> exec
| CREATE TABLE t_gen (a INT, b INT GENERATED ALWAYS AS (a+1) VIRTUAL, c INT GENERATED ALWAYS AS ( a * 2 ) STORED, d INT AS (a+1))
< ok 0 0
> query
| SHOW CREATE TABLE nrm_generated_expr_normalization.t_gen
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_gen
: bytes
| CREATE TABLE `t_gen` (
|   `a` int DEFAULT NULL,
|   `b` int GENERATED ALWAYS AS ((`a` + 1)) VIRTUAL,
|   `c` int GENERATED ALWAYS AS ((`a` * 2)) STORED,
|   `d` int GENERATED ALWAYS AS ((`a` + 1)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_generated_expr_normalization
< ok 1 0
> exec
| CREATE DATABASE nrm_generated_expr_normalization
< ok 1 0
> exec
| USE nrm_generated_expr_normalization
< ok 0 0
> exec
| CREATE TABLE t_gen (a INT, b INT GENERATED ALWAYS AS (a+1) VIRTUAL, c INT GENERATED ALWAYS AS ( a * 2 ) STORED, d INT AS (a+1))
< ok 0 0
> query
| SHOW CREATE TABLE nrm_generated_expr_normalization.t_gen
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_gen
: bytes
| CREATE TABLE `t_gen` (
|   `a` int DEFAULT NULL,
|   `b` int GENERATED ALWAYS AS ((`a` + 1)) VIRTUAL,
|   `c` int GENERATED ALWAYS AS ((`a` * 2)) STORED,
|   `d` int GENERATED ALWAYS AS ((`a` + 1)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_generated_expr_string_introducer
< ok 0 0
> exec
| CREATE DATABASE nrm_generated_expr_string_introducer
< ok 1 0
> exec
| USE nrm_generated_expr_string_introducer
< ok 0 0
> exec
| CREATE TABLE t_gen2 (a VARCHAR(20), e VARCHAR(20) GENERATED ALWAYS AS (CONCAT(a,'x')) VIRTUAL) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE nrm_generated_expr_string_introducer.t_gen2
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_gen2
: bytes
| CREATE TABLE `t_gen2` (
|   `a` varchar(20) DEFAULT NULL,
|   `e` varchar(20) GENERATED ALWAYS AS (concat(`a`,_utf8mb4'x')) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_functional_default
< ok 0 0
> exec
| CREATE DATABASE nrm_functional_default
< ok 1 0
> exec
| USE nrm_functional_default
< ok 0 0
> exec
| CREATE TABLE t_funcdef (a INT DEFAULT (1+1), b VARCHAR(36) DEFAULT (UUID()), c JSON DEFAULT (JSON_ARRAY()))
< ok 0 0
> query
| SHOW CREATE TABLE nrm_functional_default.t_funcdef
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_funcdef
: bytes
| CREATE TABLE `t_funcdef` (
|   `a` int DEFAULT ((1 + 1)),
|   `b` varchar(36) DEFAULT (uuid()),
|   `c` json DEFAULT (json_array())
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_functional_default
< ok 1 0
> exec
| CREATE DATABASE nrm_functional_default
< ok 1 0
> exec
| USE nrm_functional_default
< ok 0 0
> exec
| CREATE TABLE t_funcdef (a INT DEFAULT (1+1), b VARCHAR(36) DEFAULT (UUID()), c JSON DEFAULT (JSON_ARRAY()))
< ok 0 0
> query
| SHOW CREATE TABLE nrm_functional_default.t_funcdef
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_funcdef
: bytes
| CREATE TABLE `t_funcdef` (
|   `a` int DEFAULT ((1 + 1)),
|   `b` varchar(36) DEFAULT (uuid()),
|   `c` json DEFAULT (json_array())
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_functional_default
< ok 1 0
> exec
| CREATE DATABASE nrm_functional_default
< ok 1 0
> exec
| USE nrm_functional_default
< ok 0 0
> exec
| CREATE TABLE t_funcdef (a INT DEFAULT (1+1), b VARCHAR(36) DEFAULT (UUID()), c JSON DEFAULT (JSON_ARRAY()))
< ok 0 0
> query
| SHOW CREATE TABLE nrm_functional_default.t_funcdef
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_funcdef
: bytes
| CREATE TABLE `t_funcdef` (
|   `a` int DEFAULT ((1 + 1)),
|   `b` varchar(36) DEFAULT (uuid()),
|   `c` json DEFAULT (json_array())
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_timestamp_magic
< ok 0 0
> exec
| CREATE DATABASE nrm_timestamp_magic
< ok 1 0
> exec
| USE nrm_timestamp_magic
< ok 0 0
> exec
| CREATE TABLE t_ts (a TIMESTAMP, b TIMESTAMP NULL, d TIMESTAMP DEFAULT CURRENT_TIMESTAMP)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_timestamp_magic.t_ts
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_ts
: bytes
| CREATE TABLE `t_ts` (
|   `a` timestamp NULL DEFAULT NULL,
|   `b` timestamp NULL DEFAULT NULL,
|   `d` timestamp NULL DEFAULT CURRENT_TIMESTAMP
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_timestamp_magic
< ok 1 0
> exec
| CREATE DATABASE nrm_timestamp_magic
< ok 1 0
> exec
| USE nrm_timestamp_magic
< ok 0 0
> exec
| CREATE TABLE t_ts (a TIMESTAMP, b TIMESTAMP NULL, d TIMESTAMP DEFAULT CURRENT_TIMESTAMP)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_timestamp_magic.t_ts
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_ts
: bytes
| CREATE TABLE `t_ts` (
|   `a` timestamp NULL DEFAULT NULL,
|   `b` timestamp NULL DEFAULT NULL,
|   `d` timestamp NULL DEFAULT CURRENT_TIMESTAMP
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_timestamp_magic
< ok 1 0
> exec
| CREATE DATABASE nrm_timestamp_magic
< ok 1 0
> exec
| USE nrm_timestamp_magic
< ok 0 0
> exec
| CREATE TABLE t_ts (a TIMESTAMP, b TIMESTAMP NULL, d TIMESTAMP DEFAULT CURRENT_TIMESTAMP)
< ok 0 0
> query
| SHOW CREATE TABLE nrm_timestamp_magic.t_ts
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_ts
: bytes
| CREATE TABLE `t_ts` (
|   `a` timestamp NULL DEFAULT NULL,
|   `b` timestamp NULL DEFAULT NULL,
|   `d` timestamp NULL DEFAULT CURRENT_TIMESTAMP
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_timestamp_datetime_default_expr
< ok 0 0
> exec
| CREATE DATABASE nrm_timestamp_datetime_default_expr
< ok 1 0
> exec
| USE nrm_timestamp_datetime_default_expr
< ok 0 0
> exec
| CREATE TABLE t_ts2 (g DATETIME DEFAULT CURRENT_TIMESTAMP, h DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, i TIMESTAMP(3) NULL DEFAULT CURRENT_TIMESTAMP(3))
< ok 0 0
> query
| SHOW CREATE TABLE nrm_timestamp_datetime_default_expr.t_ts2
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_ts2
: bytes
| CREATE TABLE `t_ts2` (
|   `g` datetime DEFAULT CURRENT_TIMESTAMP,
|   `h` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   `i` timestamp(3) NULL DEFAULT CURRENT_TIMESTAMP(3)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_timestamp_datetime_default_expr
< ok 1 0
> exec
| CREATE DATABASE nrm_timestamp_datetime_default_expr
< ok 1 0
> exec
| USE nrm_timestamp_datetime_default_expr
< ok 0 0
> exec
| CREATE TABLE t_ts2 (g DATETIME DEFAULT CURRENT_TIMESTAMP, h DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, i TIMESTAMP(3) NULL DEFAULT CURRENT_TIMESTAMP(3))
< ok 0 0
> query
| SHOW CREATE TABLE nrm_timestamp_datetime_default_expr.t_ts2
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_ts2
: bytes
| CREATE TABLE `t_ts2` (
|   `g` datetime DEFAULT CURRENT_TIMESTAMP,
|   `h` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   `i` timestamp(3) NULL DEFAULT CURRENT_TIMESTAMP(3)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS nrm_timestamp_datetime_default_expr
< ok 1 0
> exec
| CREATE DATABASE nrm_timestamp_datetime_default_expr
< ok 1 0
> exec
| USE nrm_timestamp_datetime_default_expr
< ok 0 0
> exec
| CREATE TABLE t_ts2 (g DATETIME DEFAULT CURRENT_TIMESTAMP, h DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, i TIMESTAMP(3) NULL DEFAULT CURRENT_TIMESTAMP(3))
< ok 0 0
> query
| SHOW CREATE TABLE nrm_timestamp_datetime_default_expr.t_ts2
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_ts2
: bytes
| CREATE TABLE `t_ts2` (
|   `g` datetime DEFAULT CURRENT_TIMESTAMP,
|   `h` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   `i` timestamp(3) NULL DEFAULT CURRENT_TIMESTAMP(3)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
