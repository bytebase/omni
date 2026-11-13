-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_ApproximateEnumNumberDefault default#1
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
| CREATE TABLE t (a ENUM('2','1') DEFAULT (1.99999999999999999e0))
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
|   `a` enum('2','1') DEFAULT '1'
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
| 1
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` enum('2','1') DEFAULT '1'
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
| 1
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
| CREATE TABLE t (a ENUM('2','1') DEFAULT (0.99999999999999999e0))
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
|   `a` enum('2','1') DEFAULT '2'
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
| 2
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` enum('2','1') DEFAULT '2'
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
| 2
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
| CREATE TABLE t (a ENUM('2','1') DEFAULT (1.0000000000000001e0))
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
|   `a` enum('2','1') DEFAULT '2'
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
| 2
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` enum('2','1') DEFAULT '2'
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
| 2
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
| CREATE TABLE t (a ENUM('2','1') DEFAULT (2.0e0))
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
|   `a` enum('2','1') DEFAULT '1'
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
| 1
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` enum('2','1') DEFAULT '1'
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
| 1
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
| CREATE TABLE t (a SET('0','1','2') DEFAULT (0.99999999999999999e0))
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
|   `a` set('0','1','2') DEFAULT '0'
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
| 0
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` set('0','1','2') DEFAULT '0'
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
| 0
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
| CREATE TABLE t (a SET('a','1') DEFAULT (0.99999999999999999e0))
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
|   `a` set('a','1') DEFAULT 'a'
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
| a
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` set('a','1') DEFAULT 'a'
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
| a
