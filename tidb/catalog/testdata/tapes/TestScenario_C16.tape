-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestScenario_C16 default#1
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (ts DATETIME DEFAULT NOW())
< ok 0 0
> query
| SHOW CREATE TABLE t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `ts` datetime DEFAULT CURRENT_TIMESTAMP
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT DATETIME_PRECISION FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME=?
: string
| ts
< result
< column "DATETIME_PRECISION" "BIGINT"
< row
: int64 0
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> query
| SELECT LENGTH(NOW(0))
< result
< column "LENGTH(NOW(0))" "BIGINT"
< row
: int64 19
> query
| SELECT LENGTH(NOW(1))
< result
< column "LENGTH(NOW(1))" "BIGINT"
< row
: int64 21
> query
| SELECT LENGTH(NOW(2))
< result
< column "LENGTH(NOW(2))" "BIGINT"
< row
: int64 22
> query
| SELECT LENGTH(NOW(3))
< result
< column "LENGTH(NOW(3))" "BIGINT"
< row
: int64 23
> query
| SELECT LENGTH(NOW(4))
< result
< column "LENGTH(NOW(4))" "BIGINT"
< row
: int64 24
> query
| SELECT LENGTH(NOW(5))
< result
< column "LENGTH(NOW(5))" "BIGINT"
< row
: int64 25
> query
| SELECT LENGTH(NOW(6))
< result
< column "LENGTH(NOW(6))" "BIGINT"
< row
: int64 26
> exec
| DO NOW(7)
< error
| 1426 42000 Too-big precision 7 specified for 'now'. Maximum is 6.
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| DO CURDATE()
< ok 0 0
> exec
| DO CURDATE(6)
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 12 near "6)" 
> exec
| CREATE TABLE t (d DATE DEFAULT (CURDATE()))
< ok 0 0
> query
| SELECT DATETIME_PRECISION FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME=?
: string
| d
< result
< column "DATETIME_PRECISION" "BIGINT"
< row
: null
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> query
| SELECT LENGTH(CURTIME()), LENGTH(CURTIME(6))
< result
< column "LENGTH(CURTIME())" "BIGINT"
< column "LENGTH(CURTIME(6))" "BIGINT"
< row
: int64 8
: int64 15
> exec
| DO CURTIME(7)
< error
| 1426 42000 Too-big precision 7 specified for 'curtime'. Maximum is 6.
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a DATETIME DEFAULT SYSDATE())
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 42 near "SYSDATE())" 
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> query
| SELECT LENGTH(UTC_TIMESTAMP()), LENGTH(UTC_TIMESTAMP(3))
< result
< column "LENGTH(UTC_TIMESTAMP())" "BIGINT"
< column "LENGTH(UTC_TIMESTAMP(3))" "BIGINT"
< row
: int64 19
: int64 23
> exec
| CREATE TABLE t (a DATETIME DEFAULT UTC_TIMESTAMP)
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 48 near "UTC_TIMESTAMP)" 
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE VIEW v AS SELECT UNIX_TIMESTAMP() AS u0, UNIX_TIMESTAMP(NOW(6)) AS u6
< ok 0 0
> query
| SELECT DATA_TYPE, NUMERIC_SCALE FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v' AND COLUMN_NAME='u6'
< result
< column "DATA_TYPE" "VARCHAR"
< column "NUMERIC_SCALE" "BIGINT"
< row
: bytes
| decimal
: int64 6
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v' AND COLUMN_NAME='u0'
< result
< column "DATA_TYPE" "VARCHAR"
< row
: bytes
| bigint
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a DATETIME(6) DEFAULT NOW())
< error
| 1067 42000 Invalid default value for 'a'
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a DATETIME(3) DEFAULT NOW(6))
< error
| 1067 42000 Invalid default value for 'a'
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a DATETIME(6) DEFAULT NOW(6))
< ok 0 0
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a DATETIME(6) DEFAULT NOW(6) ON UPDATE NOW())
< error
| 1294 HY000 Invalid ON UPDATE clause for 'a' column
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a DATE ON UPDATE NOW())
< error
| 1294 HY000 Invalid ON UPDATE clause for 'a' column
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a DATETIME NULL)
< ok 0 0
> query
| SHOW CREATE TABLE t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` datetime DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT DATETIME_PRECISION FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME=?
: string
| a
< result
< column "DATETIME_PRECISION" "BIGINT"
< row
: int64 0
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a DATETIME(0) NULL)
< ok 0 0
> query
| SHOW CREATE TABLE t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` datetime DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT DATETIME_PRECISION FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME=?
: string
| a
< result
< column "DATETIME_PRECISION" "BIGINT"
< row
: int64 0
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a DATETIME(3) NULL)
< ok 0 0
> query
| SHOW CREATE TABLE t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` datetime(3) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT DATETIME_PRECISION FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME=?
: string
| a
< result
< column "DATETIME_PRECISION" "BIGINT"
< row
: int64 3
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a DATETIME(6) NULL)
< ok 0 0
> query
| SHOW CREATE TABLE t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` datetime(6) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT DATETIME_PRECISION FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME=?
: string
| a
< result
< column "DATETIME_PRECISION" "BIGINT"
< row
: int64 6
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a TIMESTAMP(6) NULL)
< ok 0 0
> query
| SHOW CREATE TABLE t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` timestamp(6) NULL DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT DATETIME_PRECISION FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME=?
: string
| a
< result
< column "DATETIME_PRECISION" "BIGINT"
< row
: int64 6
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (y YEAR(2))
< error
| 1818 HY000 Supports only YEAR or YEAR(4) column
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (y YEAR(3))
< error
| 1818 HY000 Supports only YEAR or YEAR(4) column
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (y YEAR(5))
< error
| 1818 HY000 Supports only YEAR or YEAR(4) column
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (y YEAR(4))
< ok 0 0
> query
| SHOW CREATE TABLE t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `y` year(4) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (ts TIMESTAMP(3) NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `ts` timestamp(3) NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
