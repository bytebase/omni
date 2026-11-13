-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

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
< column "Create Table" "TEXT"
< row
: bytes
| t_int
: bytes
| CREATE TABLE `t_int` (
|   `a` int(11) DEFAULT NULL,
|   `b` bigint(20) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL,
|   `f` smallint(6) DEFAULT NULL,
|   `o` bigint(20) DEFAULT NULL,
|   `p` tinyint(4) DEFAULT NULL,
|   `q` smallint(6) DEFAULT NULL,
|   `r` mediumint(9) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_int
: bytes
| CREATE TABLE `t_int` (
|   `a` int(11) DEFAULT NULL,
|   `b` bigint(20) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL,
|   `f` smallint(6) DEFAULT NULL,
|   `o` bigint(20) DEFAULT NULL,
|   `p` tinyint(4) DEFAULT NULL,
|   `q` smallint(6) DEFAULT NULL,
|   `r` mediumint(9) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_int
: bytes
| CREATE TABLE `t_int` (
|   `a` int(11) DEFAULT NULL,
|   `b` bigint(20) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL,
|   `f` smallint(6) DEFAULT NULL,
|   `o` bigint(20) DEFAULT NULL,
|   `p` tinyint(4) DEFAULT NULL,
|   `q` smallint(6) DEFAULT NULL,
|   `r` mediumint(9) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_int
: bytes
| CREATE TABLE `t_int` (
|   `a` int(11) DEFAULT NULL,
|   `b` bigint(20) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL,
|   `f` smallint(6) DEFAULT NULL,
|   `o` bigint(20) DEFAULT NULL,
|   `p` tinyint(4) DEFAULT NULL,
|   `q` smallint(6) DEFAULT NULL,
|   `r` mediumint(9) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_int
: bytes
| CREATE TABLE `t_int` (
|   `a` int(11) DEFAULT NULL,
|   `b` bigint(20) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL,
|   `f` smallint(6) DEFAULT NULL,
|   `o` bigint(20) DEFAULT NULL,
|   `p` tinyint(4) DEFAULT NULL,
|   `q` smallint(6) DEFAULT NULL,
|   `r` mediumint(9) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_int
: bytes
| CREATE TABLE `t_int` (
|   `a` int(11) DEFAULT NULL,
|   `b` bigint(20) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL,
|   `f` smallint(6) DEFAULT NULL,
|   `o` bigint(20) DEFAULT NULL,
|   `p` tinyint(4) DEFAULT NULL,
|   `q` smallint(6) DEFAULT NULL,
|   `r` mediumint(9) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_int
: bytes
| CREATE TABLE `t_int` (
|   `a` int(11) DEFAULT NULL,
|   `b` bigint(20) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL,
|   `f` smallint(6) DEFAULT NULL,
|   `o` bigint(20) DEFAULT NULL,
|   `p` tinyint(4) DEFAULT NULL,
|   `q` smallint(6) DEFAULT NULL,
|   `r` mediumint(9) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_int
: bytes
| CREATE TABLE `t_int` (
|   `a` int(11) DEFAULT NULL,
|   `b` bigint(20) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL,
|   `f` smallint(6) DEFAULT NULL,
|   `o` bigint(20) DEFAULT NULL,
|   `p` tinyint(4) DEFAULT NULL,
|   `q` smallint(6) DEFAULT NULL,
|   `r` mediumint(9) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_bool
: bytes
| CREATE TABLE `t_bool` (
|   `d` tinyint(1) DEFAULT NULL,
|   `e` tinyint(4) DEFAULT NULL,
|   `m` tinyint(1) DEFAULT NULL,
|   `n` tinyint(1) DEFAULT NULL,
|   `p` tinyint(4) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_bool
: bytes
| CREATE TABLE `t_bool` (
|   `d` tinyint(1) DEFAULT NULL,
|   `e` tinyint(4) DEFAULT NULL,
|   `m` tinyint(1) DEFAULT NULL,
|   `n` tinyint(1) DEFAULT NULL,
|   `p` tinyint(4) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_bool
: bytes
| CREATE TABLE `t_bool` (
|   `d` tinyint(1) DEFAULT NULL,
|   `e` tinyint(4) DEFAULT NULL,
|   `m` tinyint(1) DEFAULT NULL,
|   `n` tinyint(1) DEFAULT NULL,
|   `p` tinyint(4) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_bool
: bytes
| CREATE TABLE `t_bool` (
|   `d` tinyint(1) DEFAULT NULL,
|   `e` tinyint(4) DEFAULT NULL,
|   `m` tinyint(1) DEFAULT NULL,
|   `n` tinyint(1) DEFAULT NULL,
|   `p` tinyint(4) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_bool
: bytes
| CREATE TABLE `t_bool` (
|   `d` tinyint(1) DEFAULT NULL,
|   `e` tinyint(4) DEFAULT NULL,
|   `m` tinyint(1) DEFAULT NULL,
|   `n` tinyint(1) DEFAULT NULL,
|   `p` tinyint(4) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_uns
: bytes
| CREATE TABLE `t_uns` (
|   `h` int(10) unsigned DEFAULT NULL,
|   `i` int(11) unsigned DEFAULT NULL,
|   `b` bigint(20) unsigned DEFAULT NULL,
|   `ti` tinyint(3) unsigned DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_uns
: bytes
| CREATE TABLE `t_uns` (
|   `h` int(10) unsigned DEFAULT NULL,
|   `i` int(11) unsigned DEFAULT NULL,
|   `b` bigint(20) unsigned DEFAULT NULL,
|   `ti` tinyint(3) unsigned DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_uns
: bytes
| CREATE TABLE `t_uns` (
|   `h` int(10) unsigned DEFAULT NULL,
|   `i` int(11) unsigned DEFAULT NULL,
|   `b` bigint(20) unsigned DEFAULT NULL,
|   `ti` tinyint(3) unsigned DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_uns
: bytes
| CREATE TABLE `t_uns` (
|   `h` int(10) unsigned DEFAULT NULL,
|   `i` int(11) unsigned DEFAULT NULL,
|   `b` bigint(20) unsigned DEFAULT NULL,
|   `ti` tinyint(3) unsigned DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_zf
: bytes
| CREATE TABLE `t_zf` (
|   `j` int(10) unsigned zerofill DEFAULT NULL,
|   `k` int(5) unsigned zerofill DEFAULT NULL,
|   `l` int(10) unsigned zerofill DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_zf
: bytes
| CREATE TABLE `t_zf` (
|   `j` int(10) unsigned zerofill DEFAULT NULL,
|   `k` int(5) unsigned zerofill DEFAULT NULL,
|   `l` int(10) unsigned zerofill DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_zf
: bytes
| CREATE TABLE `t_zf` (
|   `j` int(10) unsigned zerofill DEFAULT NULL,
|   `k` int(5) unsigned zerofill DEFAULT NULL,
|   `l` int(10) unsigned zerofill DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
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
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
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
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
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
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
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
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
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
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
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
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
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
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
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
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
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
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
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
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
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
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
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
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
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
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
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
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
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
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
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
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
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
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
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
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_binu
: bytes
| CREATE TABLE `t_binu` (
|   `k` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
|   `c` char(5) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   `t` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   `ti` tinytext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   `me` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   `lo` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS nrm_char_binary_attribute_utf8mb4
< ok 3 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_binu
: bytes
| CREATE TABLE `t_binu` (
|   `k` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
|   `c` char(5) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   `t` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   `ti` tinytext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   `me` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   `lo` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS nrm_char_binary_attribute_utf8mb4
< ok 3 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_binu
: bytes
| CREATE TABLE `t_binu` (
|   `k` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
|   `c` char(5) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   `t` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   `ti` tinytext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   `me` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   `lo` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS nrm_char_binary_attribute_utf8mb4
< ok 3 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_binu
: bytes
| CREATE TABLE `t_binu` (
|   `k` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
|   `c` char(5) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   `t` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   `ti` tinytext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   `me` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   `lo` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS nrm_char_binary_attribute_utf8mb4
< ok 3 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_binu
: bytes
| CREATE TABLE `t_binu` (
|   `k` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
|   `c` char(5) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   `t` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   `ti` tinytext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   `me` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   `lo` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS nrm_char_binary_attribute_utf8mb4
< ok 3 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_binu
: bytes
| CREATE TABLE `t_binu` (
|   `k` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
|   `c` char(5) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   `t` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   `ti` tinytext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   `me` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   `lo` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_binl
: bytes
| CREATE TABLE `t_binl` (
|   `a` varchar(50) CHARACTER SET latin1 COLLATE latin1_bin DEFAULT NULL,
|   `c` char(8) CHARACTER SET latin1 COLLATE latin1_bin DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = latin1 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_binl
: bytes
| CREATE TABLE `t_binl` (
|   `a` varchar(50) CHARACTER SET latin1 COLLATE latin1_bin DEFAULT NULL,
|   `c` char(8) CHARACTER SET latin1 COLLATE latin1_bin DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = latin1 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< error
| 1064 42000 …near 'BINARY, c CHAR(4) CHARACTER SET ascii BINARY) DEFAULT CHARSET=utf8mb4' at line 1
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
< error
| 1064 42000 …near 'BINARY, c CHAR(4) CHARACTER SET ascii BINARY) DEFAULT CHARSET=utf8mb4' at line 1
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
< error
| 1064 42000 …near 'BINARY, b VARCHAR(10) CHARACTER SET utf8mb3 BINARY) DEFAULT CHARSET=utf8mb4' at line 1
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
< error
| 1064 42000 …near 'BINARY, b VARCHAR(10) CHARACTER SET utf8mb3 BINARY) DEFAULT CHARSET=utf8mb4' at line 1
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
< column "Create Table" "TEXT"
< row
: bytes
| t_binp
: bytes
| CREATE TABLE `t_binp` (
|   `a` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_binr
: bytes
| CREATE TABLE `t_binr` (
|   `a` binary(8) DEFAULT NULL,
|   `b` varbinary(16) DEFAULT NULL,
|   `c` blob DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS nrm_char_binary_attribute_real_binary_type
< ok 3 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_binr
: bytes
| CREATE TABLE `t_binr` (
|   `a` binary(8) DEFAULT NULL,
|   `b` varbinary(16) DEFAULT NULL,
|   `c` blob DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS nrm_char_binary_attribute_real_binary_type
< ok 3 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_binr
: bytes
| CREATE TABLE `t_binr` (
|   `a` binary(8) DEFAULT NULL,
|   `b` varbinary(16) DEFAULT NULL,
|   `c` blob DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_tb
: bytes
| CREATE TABLE `t_tb` (
|   `g` text DEFAULT NULL,
|   `h` blob DEFAULT NULL,
|   `j` longtext DEFAULT NULL,
|   `k` tinytext DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS nrm_text_blob_no_default_null
< ok 3 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_tb
: bytes
| CREATE TABLE `t_tb` (
|   `g` text DEFAULT NULL,
|   `h` blob DEFAULT NULL,
|   `j` longtext DEFAULT NULL,
|   `k` tinytext DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS nrm_text_blob_no_default_null
< ok 3 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_tb
: bytes
| CREATE TABLE `t_tb` (
|   `g` text DEFAULT NULL,
|   `h` blob DEFAULT NULL,
|   `j` longtext DEFAULT NULL,
|   `k` tinytext DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS nrm_text_blob_no_default_null
< ok 3 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_tb
: bytes
| CREATE TABLE `t_tb` (
|   `g` text DEFAULT NULL,
|   `h` blob DEFAULT NULL,
|   `j` longtext DEFAULT NULL,
|   `k` tinytext DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_year
: bytes
| CREATE TABLE `t_year` (
|   `a` year(4) DEFAULT NULL,
|   `b` year(4) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_year
: bytes
| CREATE TABLE `t_year` (
|   `a` year(4) DEFAULT NULL,
|   `b` year(4) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_bit
: bytes
| CREATE TABLE `t_bit` (
|   `c` bit(1) DEFAULT NULL,
|   `d` bit(8) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_bit
: bytes
| CREATE TABLE `t_bit` (
|   `c` bit(1) DEFAULT NULL,
|   `d` bit(8) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_bitd
: bytes
| CREATE TABLE `t_bitd` (
|   `a` bit(8) NOT NULL DEFAULT b'0',
|   `b` bit(8) NOT NULL DEFAULT b'101',
|   `c` bit(8) NOT NULL DEFAULT b'101'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_bitd
: bytes
| CREATE TABLE `t_bitd` (
|   `a` bit(8) NOT NULL DEFAULT b'0',
|   `b` bit(8) NOT NULL DEFAULT b'101',
|   `c` bit(8) NOT NULL DEFAULT b'101'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_bitd
: bytes
| CREATE TABLE `t_bitd` (
|   `a` bit(8) NOT NULL DEFAULT b'0',
|   `b` bit(8) NOT NULL DEFAULT b'101',
|   `c` bit(8) NOT NULL DEFAULT b'101'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_bithex
: bytes
| CREATE TABLE `t_bithex` (
|   `a` bit(8) NOT NULL DEFAULT b'101',
|   `b` bit(8) NOT NULL DEFAULT b'101',
|   `c` bit(16) NOT NULL DEFAULT b'1010101111001101'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_bithex
: bytes
| CREATE TABLE `t_bithex` (
|   `a` bit(8) NOT NULL DEFAULT b'101',
|   `b` bit(8) NOT NULL DEFAULT b'101',
|   `c` bit(16) NOT NULL DEFAULT b'1010101111001101'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_bithex
: bytes
| CREATE TABLE `t_bithex` (
|   `a` bit(8) NOT NULL DEFAULT b'101',
|   `b` bit(8) NOT NULL DEFAULT b'101',
|   `c` bit(16) NOT NULL DEFAULT b'1010101111001101'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_bitns
: bytes
| CREATE TABLE `t_bitns` (
|   `a` bit(8) NOT NULL DEFAULT b'101',
|   `b` bit(8) NOT NULL DEFAULT b'101',
|   `c` bit(16) NOT NULL DEFAULT b'100000101000010'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_bitns
: bytes
| CREATE TABLE `t_bitns` (
|   `a` bit(8) NOT NULL DEFAULT b'101',
|   `b` bit(8) NOT NULL DEFAULT b'101',
|   `c` bit(16) NOT NULL DEFAULT b'100000101000010'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_bitns
: bytes
| CREATE TABLE `t_bitns` (
|   `a` bit(8) NOT NULL DEFAULT b'101',
|   `b` bit(8) NOT NULL DEFAULT b'101',
|   `c` bit(16) NOT NULL DEFAULT b'100000101000010'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_yeard
: bytes
| CREATE TABLE `t_yeard` (
|   `a` year(4) NOT NULL DEFAULT '2000',
|   `b` year(4) NOT NULL DEFAULT '1999',
|   `c` year(4) NOT NULL DEFAULT '0000',
|   `d` year(4) NOT NULL DEFAULT '1999',
|   `e` year(4) NOT NULL DEFAULT '2005'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_yeard
: bytes
| CREATE TABLE `t_yeard` (
|   `a` year(4) NOT NULL DEFAULT '2000',
|   `b` year(4) NOT NULL DEFAULT '1999',
|   `c` year(4) NOT NULL DEFAULT '0000',
|   `d` year(4) NOT NULL DEFAULT '1999',
|   `e` year(4) NOT NULL DEFAULT '2005'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_yeard
: bytes
| CREATE TABLE `t_yeard` (
|   `a` year(4) NOT NULL DEFAULT '2000',
|   `b` year(4) NOT NULL DEFAULT '1999',
|   `c` year(4) NOT NULL DEFAULT '0000',
|   `d` year(4) NOT NULL DEFAULT '1999',
|   `e` year(4) NOT NULL DEFAULT '2005'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_yeard
: bytes
| CREATE TABLE `t_yeard` (
|   `a` year(4) NOT NULL DEFAULT '2000',
|   `b` year(4) NOT NULL DEFAULT '1999',
|   `c` year(4) NOT NULL DEFAULT '0000',
|   `d` year(4) NOT NULL DEFAULT '1999',
|   `e` year(4) NOT NULL DEFAULT '2005'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_yeard
: bytes
| CREATE TABLE `t_yeard` (
|   `a` year(4) NOT NULL DEFAULT '2000',
|   `b` year(4) NOT NULL DEFAULT '1999',
|   `c` year(4) NOT NULL DEFAULT '0000',
|   `d` year(4) NOT NULL DEFAULT '1999',
|   `e` year(4) NOT NULL DEFAULT '2005'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_bigmax
: bytes
| CREATE TABLE `t_bigmax` (
|   `big` bigint(20) unsigned NOT NULL DEFAULT '18446744073709551615',
|   `mid` bigint(20) unsigned NOT NULL DEFAULT '9223372036854775808'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_bigmax
: bytes
| CREATE TABLE `t_bigmax` (
|   `big` bigint(20) unsigned NOT NULL DEFAULT '18446744073709551615',
|   `mid` bigint(20) unsigned NOT NULL DEFAULT '9223372036854775808'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_bigbnd
: bytes
| CREATE TABLE `t_bigbnd` (
|   `lo` bigint(20) NOT NULL DEFAULT '-9223372036854775808',
|   `hi` bigint(20) NOT NULL DEFAULT '9223372036854775807'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_bigbnd
: bytes
| CREATE TABLE `t_bigbnd` (
|   `lo` bigint(20) NOT NULL DEFAULT '-9223372036854775808',
|   `hi` bigint(20) NOT NULL DEFAULT '9223372036854775807'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_tjd
: bytes
| CREATE TABLE `t_tjd` (
|   `e` json DEFAULT NULL,
|   `f` date DEFAULT NULL,
|   `g` time DEFAULT NULL,
|   `h` time(3) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS nrm_time_json_date_stable
< ok 3 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_tjd
: bytes
| CREATE TABLE `t_tjd` (
|   `e` json DEFAULT NULL,
|   `f` date DEFAULT NULL,
|   `g` time DEFAULT NULL,
|   `h` time(3) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS nrm_time_json_date_stable
< ok 3 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_tjd
: bytes
| CREATE TABLE `t_tjd` (
|   `e` json DEFAULT NULL,
|   `f` date DEFAULT NULL,
|   `g` time DEFAULT NULL,
|   `h` time(3) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS nrm_time_json_date_stable
< ok 3 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_tjd
: bytes
| CREATE TABLE `t_tjd` (
|   `e` json DEFAULT NULL,
|   `f` date DEFAULT NULL,
|   `g` time DEFAULT NULL,
|   `h` time(3) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_def
: bytes
| CREATE TABLE `t_def` (
|   `a` int(11) DEFAULT '0',
|   `b` int(11) DEFAULT '5',
|   `h` int(11) DEFAULT '7',
|   `c` varchar(20) DEFAULT 'hello',
|   `d` varchar(20) DEFAULT '',
|   `e` varchar(20) DEFAULT 'it\'s'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_def
: bytes
| CREATE TABLE `t_def` (
|   `a` int(11) DEFAULT '0',
|   `b` int(11) DEFAULT '5',
|   `h` int(11) DEFAULT '7',
|   `c` varchar(20) DEFAULT 'hello',
|   `d` varchar(20) DEFAULT '',
|   `e` varchar(20) DEFAULT 'it\'s'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_def
: bytes
| CREATE TABLE `t_def` (
|   `a` int(11) DEFAULT '0',
|   `b` int(11) DEFAULT '5',
|   `h` int(11) DEFAULT '7',
|   `c` varchar(20) DEFAULT 'hello',
|   `d` varchar(20) DEFAULT '',
|   `e` varchar(20) DEFAULT 'it\'s'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_def
: bytes
| CREATE TABLE `t_def` (
|   `a` int(11) DEFAULT '0',
|   `b` int(11) DEFAULT '5',
|   `h` int(11) DEFAULT '7',
|   `c` varchar(20) DEFAULT 'hello',
|   `d` varchar(20) DEFAULT '',
|   `e` varchar(20) DEFAULT 'it\'s'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_def
: bytes
| CREATE TABLE `t_def` (
|   `a` int(11) DEFAULT '0',
|   `b` int(11) DEFAULT '5',
|   `h` int(11) DEFAULT '7',
|   `c` varchar(20) DEFAULT 'hello',
|   `d` varchar(20) DEFAULT '',
|   `e` varchar(20) DEFAULT 'it\'s'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_def
: bytes
| CREATE TABLE `t_def` (
|   `a` int(11) DEFAULT '0',
|   `b` int(11) DEFAULT '5',
|   `h` int(11) DEFAULT '7',
|   `c` varchar(20) DEFAULT 'hello',
|   `d` varchar(20) DEFAULT '',
|   `e` varchar(20) DEFAULT 'it\'s'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_dpad
: bytes
| CREATE TABLE `t_dpad` (
|   `f` decimal(10,2) DEFAULT '0.00',
|   `g` decimal(10,2) DEFAULT '0.00'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_dpad
: bytes
| CREATE TABLE `t_dpad` (
|   `f` decimal(10,2) DEFAULT '0.00',
|   `g` decimal(10,2) DEFAULT '0.00'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_bdef
: bytes
| CREATE TABLE `t_bdef` (
|   `j` tinyint(1) DEFAULT '1',
|   `k` tinyint(1) DEFAULT '0'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_bdef
: bytes
| CREATE TABLE `t_bdef` (
|   `j` tinyint(1) DEFAULT '1',
|   `k` tinyint(1) DEFAULT '0'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_null
: bytes
| CREATE TABLE `t_null` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` varchar(10) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_null
: bytes
| CREATE TABLE `t_null` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` varchar(10) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_null
: bytes
| CREATE TABLE `t_null` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` varchar(10) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_null
: bytes
| CREATE TABLE `t_null` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` varchar(10) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_enum
: bytes
| CREATE TABLE `t_enum` (
|   `a` enum('x','y','z') DEFAULT NULL,
|   `b` set('a','b','c') DEFAULT NULL,
|   `e` enum('dq1','dq2') DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_enum
: bytes
| CREATE TABLE `t_enum` (
|   `a` enum('x','y','z') DEFAULT NULL,
|   `b` set('a','b','c') DEFAULT NULL,
|   `e` enum('dq1','dq2') DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_enum
: bytes
| CREATE TABLE `t_enum` (
|   `a` enum('x','y','z') DEFAULT NULL,
|   `b` set('a','b','c') DEFAULT NULL,
|   `e` enum('dq1','dq2') DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_comment
: bytes
| CREATE TABLE `t_comment` (
|   `a` int(11) DEFAULT NULL COMMENT 'has \'quote\' inside'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_ce80
: bytes
| CREATE TABLE `t_ce80` (
|   `a` varchar(10) DEFAULT NULL,
|   `c` varchar(10) DEFAULT NULL,
|   `d` varchar(10) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_ce80
: bytes
| CREATE TABLE `t_ce80` (
|   `a` varchar(10) DEFAULT NULL,
|   `c` varchar(10) DEFAULT NULL,
|   `d` varchar(10) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_ce80
: bytes
| CREATE TABLE `t_ce80` (
|   `a` varchar(10) DEFAULT NULL,
|   `c` varchar(10) DEFAULT NULL,
|   `d` varchar(10) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_d1
: bytes
| CREATE TABLE `t_d1` (
|   `a` varchar(10) CHARACTER SET utf8mb4 DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_gen
: bytes
| CREATE TABLE `t_gen` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) GENERATED ALWAYS AS ((`a` + 1)) VIRTUAL,
|   `c` int(11) GENERATED ALWAYS AS ((`a` * 2)) STORED,
|   `d` int(11) GENERATED ALWAYS AS ((`a` + 1)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_gen
: bytes
| CREATE TABLE `t_gen` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) GENERATED ALWAYS AS ((`a` + 1)) VIRTUAL,
|   `c` int(11) GENERATED ALWAYS AS ((`a` * 2)) STORED,
|   `d` int(11) GENERATED ALWAYS AS ((`a` + 1)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_gen
: bytes
| CREATE TABLE `t_gen` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) GENERATED ALWAYS AS ((`a` + 1)) VIRTUAL,
|   `c` int(11) GENERATED ALWAYS AS ((`a` * 2)) STORED,
|   `d` int(11) GENERATED ALWAYS AS ((`a` + 1)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_gen2
: bytes
| CREATE TABLE `t_gen2` (
|   `a` varchar(20) DEFAULT NULL,
|   `e` varchar(20) GENERATED ALWAYS AS (CONCAT(`a`,'x')) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
| CREATE TABLE t_funcdef (a INT DEFAULT (1+1), b VARCHAR(36) DEFAULT (UUID()))
< ok 0 0
> query
| SHOW CREATE TABLE nrm_functional_default.t_funcdef
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t_funcdef
: bytes
| CREATE TABLE `t_funcdef` (
|   `a` int(11) DEFAULT ((1 + 1)),
|   `b` varchar(36) DEFAULT (UUID())
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
| CREATE TABLE t_funcdef (a INT DEFAULT (1+1), b VARCHAR(36) DEFAULT (UUID()))
< ok 0 0
> query
| SHOW CREATE TABLE nrm_functional_default.t_funcdef
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t_funcdef
: bytes
| CREATE TABLE `t_funcdef` (
|   `a` int(11) DEFAULT ((1 + 1)),
|   `b` varchar(36) DEFAULT (UUID())
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS nrm_functional_default_json
< ok 0 0
> exec
| CREATE DATABASE nrm_functional_default_json
< ok 1 0
> exec
| USE nrm_functional_default_json
< ok 0 0
> exec
| CREATE TABLE t_funcjson (c JSON DEFAULT (JSON_ARRAY()))
< error
| 1101 42000 BLOB, TEXT, GEOMETRY or JSON column 'c' can't have a default value
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
< column "Create Table" "TEXT"
< row
: bytes
| t_ts
: bytes
| CREATE TABLE `t_ts` (
|   `a` timestamp NULL DEFAULT NULL,
|   `b` timestamp NULL DEFAULT NULL,
|   `d` timestamp NULL DEFAULT CURRENT_TIMESTAMP
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_ts
: bytes
| CREATE TABLE `t_ts` (
|   `a` timestamp NULL DEFAULT NULL,
|   `b` timestamp NULL DEFAULT NULL,
|   `d` timestamp NULL DEFAULT CURRENT_TIMESTAMP
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_ts
: bytes
| CREATE TABLE `t_ts` (
|   `a` timestamp NULL DEFAULT NULL,
|   `b` timestamp NULL DEFAULT NULL,
|   `d` timestamp NULL DEFAULT CURRENT_TIMESTAMP
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_ts2
: bytes
| CREATE TABLE `t_ts2` (
|   `g` datetime DEFAULT CURRENT_TIMESTAMP,
|   `h` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   `i` timestamp(3) NULL DEFAULT CURRENT_TIMESTAMP(3)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_ts2
: bytes
| CREATE TABLE `t_ts2` (
|   `g` datetime DEFAULT CURRENT_TIMESTAMP,
|   `h` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   `i` timestamp(3) NULL DEFAULT CURRENT_TIMESTAMP(3)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t_ts2
: bytes
| CREATE TABLE `t_ts2` (
|   `g` datetime DEFAULT CURRENT_TIMESTAMP,
|   `h` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   `i` timestamp(3) NULL DEFAULT CURRENT_TIMESTAMP(3)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
