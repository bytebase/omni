-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestScenario_C25 default#1
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (d DECIMAL)
< ok 0 0
> query
| SELECT COLUMN_TYPE, NUMERIC_PRECISION, NUMERIC_SCALE
| 			 FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='d'
< result
< column "COLUMN_TYPE" "TEXT"
< column "NUMERIC_PRECISION" "UNSIGNED BIGINT"
< column "NUMERIC_SCALE" "UNSIGNED BIGINT"
< row
: bytes
| decimal(10,0)
: uint64 10
: uint64 0
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (
| 			d5  DECIMAL(5),
| 			d15 DECIMAL(15),
| 			d65 DECIMAL(65)
| 		)
< ok 0 0
> query
| SELECT COLUMN_NAME, COLUMN_TYPE, NUMERIC_PRECISION, NUMERIC_SCALE
| 			 FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
| 			 ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< column "COLUMN_TYPE" "TEXT"
< column "NUMERIC_PRECISION" "UNSIGNED BIGINT"
< column "NUMERIC_SCALE" "UNSIGNED BIGINT"
< row
: bytes
| d5
: bytes
| decimal(5,0)
: uint64 5
: uint64 0
< row
: bytes
| d15
: bytes
| decimal(15,0)
: uint64 15
: uint64 0
< row
: bytes
| d65
: bytes
| decimal(65,0)
: uint64 65
: uint64 0
> query
| SHOW CREATE TABLE t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `d5` decimal(5,0) DEFAULT NULL,
|   `d15` decimal(15,0) DEFAULT NULL,
|   `d65` decimal(65,0) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE ok_max_p (d DECIMAL(65,30))
< ok 0 0
> query
| SELECT COLUMN_TYPE FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='ok_max_p' AND COLUMN_NAME='d'
< result
< column "COLUMN_TYPE" "TEXT"
< row
: bytes
| decimal(65,30)
> exec
| CREATE TABLE err_p_gt_65 (d DECIMAL(66, 0))
< error
| 1426 42000 Too big precision 66 specified for column 'd'. Maximum is 65.
> exec
| CREATE TABLE err_s_gt_30 (d DECIMAL(40, 31))
< error
| 1425 42000 Too big scale 31 specified for column 'd'. Maximum is 30.
> exec
| CREATE TABLE err_s_gt_p (d DECIMAL(5, 6))
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column 'd').
> exec
| CREATE TABLE err_neg (d DECIMAL(-1, 0))
< error
| 1064 42000 …near '-1, 0))' at line 1
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (
| 			d1 DECIMAL(10,2) UNSIGNED,
| 			d2 DECIMAL(10,2) UNSIGNED ZEROFILL,
| 			d3 NUMERIC(10,2) UNSIGNED
| 		)
< ok 0 0
> query
| SELECT COLUMN_NAME, COLUMN_TYPE, DATA_TYPE, NUMERIC_PRECISION, NUMERIC_SCALE
| 			 FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
| 			 ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< column "COLUMN_TYPE" "TEXT"
< column "DATA_TYPE" "TEXT"
< column "NUMERIC_PRECISION" "UNSIGNED BIGINT"
< column "NUMERIC_SCALE" "UNSIGNED BIGINT"
< row
: bytes
| d1
: bytes
| decimal(10,2) unsigned
: bytes
| decimal
: uint64 10
: uint64 2
< row
: bytes
| d2
: bytes
| decimal(10,2) unsigned zerofill
: bytes
| decimal
: uint64 10
: uint64 2
< row
: bytes
| d3
: bytes
| decimal(10,2) unsigned
: bytes
| decimal
: uint64 10
: uint64 2
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (
| 			a DECIMAL,
| 			b DECIMAL(5),
| 			c DECIMAL(5,0),
| 			d DECIMAL(10,0)
| 		)
< ok 0 0
> query
| SELECT COLUMN_NAME, COLUMN_TYPE FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
| 			 ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< column "COLUMN_TYPE" "TEXT"
< row
: bytes
| a
: bytes
| decimal(10,0)
< row
: bytes
| b
: bytes
| decimal(5,0)
< row
: bytes
| c
: bytes
| decimal(5,0)
< row
: bytes
| d
: bytes
| decimal(10,0)
> query
| SHOW CREATE TABLE t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` decimal(10,0) DEFAULT NULL,
|   `b` decimal(5,0) DEFAULT NULL,
|   `c` decimal(5,0) DEFAULT NULL,
|   `d` decimal(10,0) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
