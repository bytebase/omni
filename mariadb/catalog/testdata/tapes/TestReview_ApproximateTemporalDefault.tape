-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_ApproximateTemporalDefault default#1
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
| CREATE TABLE t (a DATETIME(6) DEFAULT (20200101123456.123456e0))
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
|   `a` datetime(6) DEFAULT '2020-01-01 12:34:56.125000'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT CAST(a AS CHAR) FROM t
< result
< column "CAST(a AS CHAR)" "VARCHAR"
< row
: bytes
| 2020-01-01 12:34:56.125000
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` datetime(6) DEFAULT '2020-01-01 12:34:56.125000'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT CAST(a AS CHAR) FROM t
< result
< column "CAST(a AS CHAR)" "VARCHAR"
< row
: bytes
| 2020-01-01 12:34:56.125000
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
| CREATE TABLE t (a DATETIME(6) DEFAULT (20200101123456.123456))
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
|   `a` datetime(6) DEFAULT '2020-01-01 12:34:56.123456'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT CAST(a AS CHAR) FROM t
< result
< column "CAST(a AS CHAR)" "VARCHAR"
< row
: bytes
| 2020-01-01 12:34:56.123456
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` datetime(6) DEFAULT '2020-01-01 12:34:56.123456'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT CAST(a AS CHAR) FROM t
< result
< column "CAST(a AS CHAR)" "VARCHAR"
< row
: bytes
| 2020-01-01 12:34:56.123456
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
| CREATE TABLE t (a DATETIME(3) DEFAULT (20200101123456.123e0))
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
|   `a` datetime(3) DEFAULT '2020-01-01 12:34:56.121'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT CAST(a AS CHAR) FROM t
< result
< column "CAST(a AS CHAR)" "VARCHAR"
< row
: bytes
| 2020-01-01 12:34:56.121
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` datetime(3) DEFAULT '2020-01-01 12:34:56.121'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT CAST(a AS CHAR) FROM t
< result
< column "CAST(a AS CHAR)" "VARCHAR"
< row
: bytes
| 2020-01-01 12:34:56.121
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
| CREATE TABLE t (a DATETIME(6) DEFAULT (20200101123456.1e0))
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
|   `a` datetime(6) DEFAULT '2020-01-01 12:34:56.101562'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT CAST(a AS CHAR) FROM t
< result
< column "CAST(a AS CHAR)" "VARCHAR"
< row
: bytes
| 2020-01-01 12:34:56.101562
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` datetime(6) DEFAULT '2020-01-01 12:34:56.101562'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT CAST(a AS CHAR) FROM t
< result
< column "CAST(a AS CHAR)" "VARCHAR"
< row
: bytes
| 2020-01-01 12:34:56.101562
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
| CREATE TABLE t (a DATETIME DEFAULT (2.0200101123456e13))
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
|   `a` datetime DEFAULT '2020-01-01 12:34:56'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT CAST(a AS CHAR) FROM t
< result
< column "CAST(a AS CHAR)" "VARCHAR"
< row
: bytes
| 2020-01-01 12:34:56
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` datetime DEFAULT '2020-01-01 12:34:56'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT CAST(a AS CHAR) FROM t
< result
< column "CAST(a AS CHAR)" "VARCHAR"
< row
: bytes
| 2020-01-01 12:34:56
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
| CREATE TABLE t (a DATETIME DEFAULT (20200101123456e0))
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
|   `a` datetime DEFAULT '2020-01-01 12:34:56'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT CAST(a AS CHAR) FROM t
< result
< column "CAST(a AS CHAR)" "VARCHAR"
< row
: bytes
| 2020-01-01 12:34:56
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` datetime DEFAULT '2020-01-01 12:34:56'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT CAST(a AS CHAR) FROM t
< result
< column "CAST(a AS CHAR)" "VARCHAR"
< row
: bytes
| 2020-01-01 12:34:56
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
| CREATE TABLE t (a DATE DEFAULT (2.0200101e7))
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
|   `a` date DEFAULT '2020-01-01'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT CAST(a AS CHAR) FROM t
< result
< column "CAST(a AS CHAR)" "VARCHAR"
< row
: bytes
| 2020-01-01
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` date DEFAULT '2020-01-01'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT CAST(a AS CHAR) FROM t
< result
< column "CAST(a AS CHAR)" "VARCHAR"
< row
: bytes
| 2020-01-01
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
| CREATE TABLE t (a DATE DEFAULT (20200101e0))
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
|   `a` date DEFAULT '2020-01-01'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT CAST(a AS CHAR) FROM t
< result
< column "CAST(a AS CHAR)" "VARCHAR"
< row
: bytes
| 2020-01-01
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` date DEFAULT '2020-01-01'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT CAST(a AS CHAR) FROM t
< result
< column "CAST(a AS CHAR)" "VARCHAR"
< row
: bytes
| 2020-01-01
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
| CREATE TABLE t (a TIME(6) DEFAULT (123456.123456e0))
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
|   `a` time(6) DEFAULT '12:34:56.123456'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT CAST(a AS CHAR) FROM t
< result
< column "CAST(a AS CHAR)" "VARCHAR"
< row
: bytes
| 12:34:56.123456
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` time(6) DEFAULT '12:34:56.123456'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT CAST(a AS CHAR) FROM t
< result
< column "CAST(a AS CHAR)" "VARCHAR"
< row
: bytes
| 12:34:56.123456
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
| CREATE TABLE t (a TIME(6) DEFAULT (123456.123456))
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
|   `a` time(6) DEFAULT '12:34:56.123456'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT CAST(a AS CHAR) FROM t
< result
< column "CAST(a AS CHAR)" "VARCHAR"
< row
: bytes
| 12:34:56.123456
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` time(6) DEFAULT '12:34:56.123456'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT CAST(a AS CHAR) FROM t
< result
< column "CAST(a AS CHAR)" "VARCHAR"
< row
: bytes
| 12:34:56.123456
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
| CREATE TABLE t (a TIME DEFAULT (1.005e2))
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
|   `a` time DEFAULT '00:01:00'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT CAST(a AS CHAR) FROM t
< result
< column "CAST(a AS CHAR)" "VARCHAR"
< row
: bytes
| 00:01:00
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` time DEFAULT '00:01:00'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT CAST(a AS CHAR) FROM t
< result
< column "CAST(a AS CHAR)" "VARCHAR"
< row
: bytes
| 00:01:00
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
| CREATE TABLE t (a TIME(3) DEFAULT (100.5e0))
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
|   `a` time(3) DEFAULT '00:01:00.500'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT CAST(a AS CHAR) FROM t
< result
< column "CAST(a AS CHAR)" "VARCHAR"
< row
: bytes
| 00:01:00.500
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` time(3) DEFAULT '00:01:00.500'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT CAST(a AS CHAR) FROM t
< result
< column "CAST(a AS CHAR)" "VARCHAR"
< row
: bytes
| 00:01:00.500
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
| CREATE TABLE t (a TIME(3) DEFAULT (-100.5e0))
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
|   `a` time(3) DEFAULT '-00:01:00.500'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT CAST(a AS CHAR) FROM t
< result
< column "CAST(a AS CHAR)" "VARCHAR"
< row
: bytes
| -00:01:00.500
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` time(3) DEFAULT '-00:01:00.500'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT CAST(a AS CHAR) FROM t
< result
< column "CAST(a AS CHAR)" "VARCHAR"
< row
: bytes
| -00:01:00.500
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
| CREATE TABLE t (a TIME(6) DEFAULT (8385959.999999e0))
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
|   `a` time(6) DEFAULT '838:59:59.999998'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT CAST(a AS CHAR) FROM t
< result
< column "CAST(a AS CHAR)" "VARCHAR"
< row
: bytes
| 838:59:59.999998
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` time(6) DEFAULT '838:59:59.999998'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT CAST(a AS CHAR) FROM t
< result
< column "CAST(a AS CHAR)" "VARCHAR"
< row
: bytes
| 838:59:59.999998
