-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestScenario_C16 default#1
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT DATETIME_PRECISION FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME=?
: string
| ts
< result
< column "DATETIME_PRECISION" "UNSIGNED INT"
< row
: int64 0
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
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
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| DO CURDATE()
< ok 0 0
> exec
| DO CURDATE(6)
< error
| 1064 42000 …near '6)' at line 1
> exec
| CREATE TABLE t (d DATE DEFAULT (CURDATE()))
< ok 0 0
> query
| SELECT DATETIME_PRECISION FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME=?
: string
| d
< result
< column "DATETIME_PRECISION" "UNSIGNED INT"
< row
: null
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
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
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a DATETIME DEFAULT SYSDATE())
< error
| 1064 42000 …near 'SYSDATE())' at line 1
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 1 0
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
| 1064 42000 …near 'UTC_TIMESTAMP)' at line 1
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
| CREATE VIEW v AS SELECT UNIX_TIMESTAMP() AS u0, UNIX_TIMESTAMP(NOW(6)) AS u6
< ok 0 0
> query
| SELECT DATA_TYPE, NUMERIC_SCALE FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v' AND COLUMN_NAME='u6'
< result
< column "DATA_TYPE" "LONGTEXT"
< column "NUMERIC_SCALE" "UNSIGNED INT"
< row
: bytes
| decimal
: int64 6
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v' AND COLUMN_NAME='u0'
< result
< column "DATA_TYPE" "LONGTEXT"
< row
: bytes
| bigint
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
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 1 0
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
< ok 1 0
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
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a DATETIME(6) DEFAULT NOW(6))
< ok 0 0
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
| CREATE TABLE t (a DATETIME(6) DEFAULT NOW(6) ON UPDATE NOW())
< error
| 1294 HY000 Invalid ON UPDATE clause for 'a' column
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
| CREATE TABLE t (a DATE ON UPDATE NOW())
< error
| 1294 HY000 Invalid ON UPDATE clause for 'a' column
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
| SET SESSION explicit_defaults_for_timestamp=1
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT DATETIME_PRECISION FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME=?
: string
| a
< result
< column "DATETIME_PRECISION" "UNSIGNED INT"
< row
: int64 0
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
| SET SESSION explicit_defaults_for_timestamp=1
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT DATETIME_PRECISION FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME=?
: string
| a
< result
< column "DATETIME_PRECISION" "UNSIGNED INT"
< row
: int64 0
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
| SET SESSION explicit_defaults_for_timestamp=1
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT DATETIME_PRECISION FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME=?
: string
| a
< result
< column "DATETIME_PRECISION" "UNSIGNED INT"
< row
: int64 3
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
| SET SESSION explicit_defaults_for_timestamp=1
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT DATETIME_PRECISION FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME=?
: string
| a
< result
< column "DATETIME_PRECISION" "UNSIGNED INT"
< row
: int64 6
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
| SET SESSION explicit_defaults_for_timestamp=1
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT DATETIME_PRECISION FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME=?
: string
| a
< result
< column "DATETIME_PRECISION" "UNSIGNED INT"
< row
: int64 6
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
| CREATE TABLE t (y YEAR(2))
< error
| 1818 HY000 Invalid display width. Use YEAR instead.
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
| CREATE TABLE t (y YEAR(3))
< error
| 1818 HY000 Invalid display width. Use YEAR instead.
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
| CREATE TABLE t (y YEAR(5))
< error
| 1818 HY000 Invalid display width. Use YEAR instead.
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
|   `y` year DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
| SET SESSION explicit_defaults_for_timestamp=0
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
|   `ts` timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
