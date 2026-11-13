-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_ApproximateBitAndYearDefault default#1
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
| CREATE TABLE t (a BIT(8) DEFAULT (1.9e0))
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
|   `a` bit(8) DEFAULT b'1'
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
|q "\x01"
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` bit(8) DEFAULT b'1'
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
|q "\x01"
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
| CREATE TABLE t (a BIT(8) DEFAULT (1.5e0))
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
|   `a` bit(8) DEFAULT b'1'
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
|q "\x01"
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` bit(8) DEFAULT b'1'
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
|q "\x01"
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
| CREATE TABLE t (a BIT(8) DEFAULT (2.5e0))
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
|   `a` bit(8) DEFAULT b'10'
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
|q "\x02"
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` bit(8) DEFAULT b'10'
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
|q "\x02"
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
| CREATE TABLE t (a BIT(8) DEFAULT (3.7e0))
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
|   `a` bit(8) DEFAULT b'11'
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
|q "\x03"
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` bit(8) DEFAULT b'11'
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
|q "\x03"
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
| CREATE TABLE t (a BIT(8) DEFAULT (0.9e0))
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
|   `a` bit(8) DEFAULT b'0'
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
|q "\x00"
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` bit(8) DEFAULT b'0'
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
|q "\x00"
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
| CREATE TABLE t (a BIT(8) DEFAULT (-0.9e0))
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
|   `a` bit(8) DEFAULT b'0'
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
|q "\x00"
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` bit(8) DEFAULT b'0'
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
|q "\x00"
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
| CREATE TABLE t (a YEAR DEFAULT (1.9e0))
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
|   `a` year(4) DEFAULT 2001
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
| 2001
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` year(4) DEFAULT 2001
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
| 2001
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
| CREATE TABLE t (a YEAR DEFAULT (1.5e0))
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
|   `a` year(4) DEFAULT 2001
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
| 2001
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` year(4) DEFAULT 2001
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
| 2001
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
| CREATE TABLE t (a YEAR DEFAULT (2.5e0))
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
|   `a` year(4) DEFAULT 2002
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
| 2002
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` year(4) DEFAULT 2002
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
| 2002
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
| CREATE TABLE t (a YEAR DEFAULT (2001.9e0))
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
|   `a` year(4) DEFAULT 2001
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
| 2001
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` year(4) DEFAULT 2001
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
| 2001
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
| CREATE TABLE t (a YEAR DEFAULT (0.9e0))
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
|   `a` year(4) DEFAULT 0000
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
| 0000
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` year(4) DEFAULT 0000
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
| 0000
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
| CREATE TABLE t (a BIT(8) DEFAULT (1.99999999999999999e0))
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
|   `a` bit(8) DEFAULT b'10'
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
|q "\x02"
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` bit(8) DEFAULT b'10'
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
|q "\x02"
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
| CREATE TABLE t (a BIT(8) DEFAULT (0.99999999999999999e0))
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
|   `a` bit(8) DEFAULT b'1'
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
|q "\x01"
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` bit(8) DEFAULT b'1'
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
|q "\x01"
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
| CREATE TABLE t (a BIT(8) DEFAULT (2.0000000000000001e0))
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
|   `a` bit(8) DEFAULT b'10'
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
|q "\x02"
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` bit(8) DEFAULT b'10'
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
|q "\x02"
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
| CREATE TABLE t (a BIT(8) DEFAULT (1.9999999999e0))
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
|   `a` bit(8) DEFAULT b'1'
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
|q "\x01"
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` bit(8) DEFAULT b'1'
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
|q "\x01"
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
| CREATE TABLE t (a BIT(64) DEFAULT (9007199254740993e0))
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
|   `a` bit(64) DEFAULT b'100000000000000000000000000000000000000000000000000000'
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
|q "\x00 \x00\x00\x00\x00\x00\x00"
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` bit(64) DEFAULT b'100000000000000000000000000000000000000000000000000000'
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
|q "\x00 \x00\x00\x00\x00\x00\x00"
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
| CREATE TABLE t (a BIT(64) DEFAULT (9007199254740993))
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
|   `a` bit(64) DEFAULT b'100000000000000000000000000000000000000000000000000001'
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
|q "\x00 \x00\x00\x00\x00\x00\x01"
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` bit(64) DEFAULT b'100000000000000000000000000000000000000000000000000001'
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
|q "\x00 \x00\x00\x00\x00\x00\x01"
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
| CREATE TABLE t (a YEAR DEFAULT (1.99999999999999999e0))
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
|   `a` year(4) DEFAULT 2002
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
| 2002
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` year(4) DEFAULT 2002
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
| 2002
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
| CREATE TABLE t (a YEAR DEFAULT (0.99999999999999999e0))
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
|   `a` year(4) DEFAULT 2001
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
| 2001
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` year(4) DEFAULT 2001
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
| 2001
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
| CREATE TABLE t (a YEAR DEFAULT (2000.9999999999999999e0))
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
|   `a` year(4) DEFAULT 2001
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
| 2001
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` year(4) DEFAULT 2001
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
| 2001
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
| CREATE TABLE t (a BIT(8) DEFAULT (1.9))
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
|   `a` bit(8) DEFAULT b'10'
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
|q "\x02"
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` bit(8) DEFAULT b'10'
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
|q "\x02"
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
| CREATE TABLE t (a BIT(8) DEFAULT (1.5))
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
|   `a` bit(8) DEFAULT b'10'
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
|q "\x02"
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` bit(8) DEFAULT b'10'
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
|q "\x02"
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
| CREATE TABLE t (a BIT(8) DEFAULT (1.7))
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
|   `a` bit(8) DEFAULT b'10'
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
|q "\x02"
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` bit(8) DEFAULT b'10'
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
|q "\x02"
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
| CREATE TABLE t (a INT DEFAULT (1.9e0))
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
|   `a` int(11) DEFAULT 2
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
|   `a` int(11) DEFAULT 2
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
| CREATE TABLE t (a INT DEFAULT (2.5e0))
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
|   `a` int(11) DEFAULT 2
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
|   `a` int(11) DEFAULT 2
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
| CREATE TABLE t (a INT DEFAULT (3.5e0))
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
|   `a` int(11) DEFAULT 4
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
| 4
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT 4
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
| 4
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
| CREATE TABLE t (a INT DEFAULT (-2.5e0))
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
|   `a` int(11) DEFAULT -2
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
| -2
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT -2
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
| -2
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
| CREATE TABLE t (a INT DEFAULT (2.5))
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
|   `a` int(11) DEFAULT 3
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
| 3
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT 3
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
| 3
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
| CREATE TABLE t (a DECIMAL(5,0) DEFAULT (2.5e0))
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
|   `a` decimal(5,0) DEFAULT 3
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
| 3
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` decimal(5,0) DEFAULT 3
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
| 3
