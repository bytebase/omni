-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_YearColumnDefault default#1
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
| CREATE TABLE t (a YEAR DEFAULT 0)
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
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 0
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
| SELECT a FROM t
< result
< column "a" "YEAR"
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
| CREATE TABLE t (a YEAR DEFAULT 0000)
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
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 0
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
| SELECT a FROM t
< result
< column "a" "YEAR"
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
| CREATE TABLE t (a YEAR DEFAULT 00000)
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
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 0
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
| SELECT a FROM t
< result
< column "a" "YEAR"
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
| CREATE TABLE t (a YEAR DEFAULT +0)
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
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 0
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
| SELECT a FROM t
< result
< column "a" "YEAR"
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
| CREATE TABLE t (a YEAR DEFAULT 0.0)
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
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 0
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
| SELECT a FROM t
< result
< column "a" "YEAR"
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
| CREATE TABLE t (a YEAR DEFAULT (0))
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
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 0
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
| SELECT a FROM t
< result
< column "a" "YEAR"
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
| CREATE TABLE t (a YEAR DEFAULT '0')
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
|   `a` year(4) DEFAULT 2000
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 2000
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` year(4) DEFAULT 2000
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 2000
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
| CREATE TABLE t (a YEAR DEFAULT '00')
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
|   `a` year(4) DEFAULT 2000
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 2000
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` year(4) DEFAULT 2000
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 2000
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
| CREATE TABLE t (a YEAR DEFAULT '000')
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
|   `a` year(4) DEFAULT 2000
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 2000
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` year(4) DEFAULT 2000
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 2000
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
| CREATE TABLE t (a YEAR DEFAULT '0000')
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
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 0
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
| SELECT a FROM t
< result
< column "a" "YEAR"
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
| CREATE TABLE t (a YEAR DEFAULT '00000')
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
|   `a` year(4) DEFAULT 2000
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 2000
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` year(4) DEFAULT 2000
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 2000
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
| CREATE TABLE t (a YEAR DEFAULT '+000')
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
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 0
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
| SELECT a FROM t
< result
< column "a" "YEAR"
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
| CREATE TABLE t (a YEAR DEFAULT ' 0000 ')
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
|   `a` year(4) DEFAULT 2000
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 2000
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` year(4) DEFAULT 2000
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 2000
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
| CREATE TABLE t (a YEAR DEFAULT '-0')
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
|   `a` year(4) DEFAULT 2000
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 2000
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` year(4) DEFAULT 2000
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 2000
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
| CREATE TABLE t (a YEAR DEFAULT '0.0')
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
|   `a` year(4) DEFAULT 2000
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 2000
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` year(4) DEFAULT 2000
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 2000
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
| CREATE TABLE t (a YEAR DEFAULT ('0'))
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
|   `a` year(4) DEFAULT 2000
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 2000
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` year(4) DEFAULT 2000
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 2000
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
| CREATE TABLE t (a YEAR DEFAULT _utf8mb4'0')
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
|   `a` year(4) DEFAULT 2000
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 2000
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` year(4) DEFAULT 2000
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 2000
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
| CREATE TABLE t (a YEAR DEFAULT 1)
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
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 2001
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
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 2001
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
| CREATE TABLE t (a YEAR DEFAULT '1')
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
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 2001
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
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 2001
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
| CREATE TABLE t (a YEAR DEFAULT '01')
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
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 2001
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
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 2001
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
| CREATE TABLE t (a YEAR DEFAULT 1.0)
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
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 2001
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
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 2001
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
| CREATE TABLE t (a YEAR DEFAULT 69)
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
|   `a` year(4) DEFAULT 2069
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 2069
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` year(4) DEFAULT 2069
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 2069
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
| CREATE TABLE t (a YEAR DEFAULT 70)
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
|   `a` year(4) DEFAULT 1970
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 1970
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` year(4) DEFAULT 1970
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 1970
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
| CREATE TABLE t (a YEAR DEFAULT '70')
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
|   `a` year(4) DEFAULT 1970
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 1970
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` year(4) DEFAULT 1970
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 1970
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
| CREATE TABLE t (a YEAR DEFAULT 99)
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
|   `a` year(4) DEFAULT 1999
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 1999
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` year(4) DEFAULT 1999
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 1999
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
| CREATE TABLE t (a YEAR DEFAULT '99')
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
|   `a` year(4) DEFAULT 1999
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 1999
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` year(4) DEFAULT 1999
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 1999
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
| CREATE TABLE t (a YEAR DEFAULT ' 99 ')
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
|   `a` year(4) DEFAULT 1999
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 1999
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` year(4) DEFAULT 1999
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 1999
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
| CREATE TABLE t (a YEAR DEFAULT _utf8mb4'99')
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
|   `a` year(4) DEFAULT 1999
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 1999
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` year(4) DEFAULT 1999
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 1999
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
| CREATE TABLE t (a YEAR DEFAULT 1901)
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
|   `a` year(4) DEFAULT 1901
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 1901
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` year(4) DEFAULT 1901
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 1901
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
| CREATE TABLE t (a YEAR DEFAULT '1901')
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
|   `a` year(4) DEFAULT 1901
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 1901
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` year(4) DEFAULT 1901
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 1901
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
| CREATE TABLE t (a YEAR DEFAULT '01901')
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
|   `a` year(4) DEFAULT 1901
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 1901
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` year(4) DEFAULT 1901
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 1901
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
| CREATE TABLE t (a YEAR DEFAULT 2155)
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
|   `a` year(4) DEFAULT 2155
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 2155
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` year(4) DEFAULT 2155
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 2155
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
| CREATE TABLE t (a YEAR DEFAULT 2000)
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
|   `a` year(4) DEFAULT 2000
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 2000
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` year(4) DEFAULT 2000
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 2000
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
| CREATE TABLE t (a YEAR DEFAULT '2000')
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
|   `a` year(4) DEFAULT 2000
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 2000
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` year(4) DEFAULT 2000
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 2000
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
| CREATE TABLE t (a YEAR DEFAULT 0x41)
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
|   `a` year(4) DEFAULT 2065
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 2065
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` year(4) DEFAULT 2065
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 2065
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
| CREATE TABLE t (a YEAR DEFAULT 0x0)
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
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 0
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
| SELECT a FROM t
< result
< column "a" "YEAR"
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
| CREATE TABLE t (a YEAR DEFAULT 0x7E3)
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
|   `a` year(4) DEFAULT 2019
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 2019
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` year(4) DEFAULT 2019
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 2019
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
| CREATE TABLE t (a YEAR DEFAULT b'1100011')
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
|   `a` year(4) DEFAULT 1999
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 1999
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` year(4) DEFAULT 1999
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 1999
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
| CREATE TABLE t (a YEAR DEFAULT b'0')
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
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: int64 0
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
| SELECT a FROM t
< result
< column "a" "YEAR"
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
| CREATE TABLE t (a YEAR DEFAULT NULL)
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
|   `a` year(4) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: null
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` year(4) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT a FROM t
< result
< column "a" "YEAR"
< row
: null
