-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_NegativeZeroStringDefault default#1
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
| CREATE TABLE t (a VARCHAR(12) DEFAULT (-0))
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
|   `a` varchar(12) DEFAULT '0'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "VARCHAR"
< row
: bytes
| 30
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` varchar(12) DEFAULT '0'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "VARCHAR"
< row
: bytes
| 30
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
| CREATE TABLE t (a VARCHAR(12) DEFAULT (-0.0))
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
|   `a` varchar(12) DEFAULT '0.0'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "VARCHAR"
< row
: bytes
| 302E30
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` varchar(12) DEFAULT '0.0'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "VARCHAR"
< row
: bytes
| 302E30
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
| CREATE TABLE t (a VARCHAR(12) DEFAULT (-0.00))
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
|   `a` varchar(12) DEFAULT '0.00'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "VARCHAR"
< row
: bytes
| 302E3030
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` varchar(12) DEFAULT '0.00'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "VARCHAR"
< row
: bytes
| 302E3030
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
| CREATE TABLE t (a VARCHAR(12) DEFAULT (-00))
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
|   `a` varchar(12) DEFAULT '0'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "VARCHAR"
< row
: bytes
| 30
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` varchar(12) DEFAULT '0'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "VARCHAR"
< row
: bytes
| 30
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
| CREATE TABLE t (a VARCHAR(12) DEFAULT (-0e0))
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
|   `a` varchar(12) DEFAULT '0'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "VARCHAR"
< row
: bytes
| 30
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` varchar(12) DEFAULT '0'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "VARCHAR"
< row
: bytes
| 30
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
| CREATE TABLE t (a VARCHAR(12) DEFAULT (-0.0e0))
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
|   `a` varchar(12) DEFAULT '0'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "VARCHAR"
< row
: bytes
| 30
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` varchar(12) DEFAULT '0'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "VARCHAR"
< row
: bytes
| 30
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
| CREATE TABLE t (a VARCHAR(12) DEFAULT (+0))
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
|   `a` varchar(12) DEFAULT '0'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "VARCHAR"
< row
: bytes
| 30
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` varchar(12) DEFAULT '0'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "VARCHAR"
< row
: bytes
| 30
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
| CREATE TABLE t (a VARCHAR(12) DEFAULT (-0.5))
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
|   `a` varchar(12) DEFAULT '-0.5'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "VARCHAR"
< row
: bytes
| 2D302E35
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` varchar(12) DEFAULT '-0.5'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "VARCHAR"
< row
: bytes
| 2D302E35
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
| CREATE TABLE t (a VARCHAR(12) DEFAULT (-0.5e0))
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
|   `a` varchar(12) DEFAULT '-0.5'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "VARCHAR"
< row
: bytes
| 2D302E35
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` varchar(12) DEFAULT '-0.5'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "VARCHAR"
< row
: bytes
| 2D302E35
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
| CREATE TABLE t (a VARCHAR(12) DEFAULT (-1))
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
|   `a` varchar(12) DEFAULT '-1'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "VARCHAR"
< row
: bytes
| 2D31
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` varchar(12) DEFAULT '-1'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "VARCHAR"
< row
: bytes
| 2D31
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
| CREATE TABLE t (a VARBINARY(12) DEFAULT (-0))
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
|   `a` varbinary(12) DEFAULT x'30'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "VARCHAR"
< row
: bytes
| 30
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` varbinary(12) DEFAULT x'30'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "VARCHAR"
< row
: bytes
| 30
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
| CREATE TABLE t (a VARBINARY(12) DEFAULT (-0.0))
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
|   `a` varbinary(12) DEFAULT x'302e30'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "VARCHAR"
< row
: bytes
| 302E30
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` varbinary(12) DEFAULT x'302e30'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "VARCHAR"
< row
: bytes
| 302E30
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
| CREATE TABLE t (a VARBINARY(12) DEFAULT (-0e0))
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
|   `a` varbinary(12) DEFAULT x'30'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "VARCHAR"
< row
: bytes
| 30
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` varbinary(12) DEFAULT x'30'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "VARCHAR"
< row
: bytes
| 30
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
| CREATE TABLE t (a BINARY(4) DEFAULT (-0))
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
|   `a` binary(4) DEFAULT x'30000000'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "VARCHAR"
< row
: bytes
| 30000000
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` binary(4) DEFAULT x'30000000'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "VARCHAR"
< row
: bytes
| 30000000
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
| CREATE TABLE t (a TEXT DEFAULT (-0))
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
|   `a` text DEFAULT 0
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "MEDIUMTEXT"
< row
: bytes
| 30
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` text DEFAULT 0
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "MEDIUMTEXT"
< row
: bytes
| 30
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
| CREATE TABLE t (a TEXT DEFAULT (-0.0))
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
|   `a` text DEFAULT 0.0
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "MEDIUMTEXT"
< row
: bytes
| 302E30
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` text DEFAULT 0.0
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "MEDIUMTEXT"
< row
: bytes
| 302E30
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
| CREATE TABLE t (a TEXT DEFAULT (-0.000))
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
|   `a` text DEFAULT 0.000
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "MEDIUMTEXT"
< row
: bytes
| 302E303030
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` text DEFAULT 0.000
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "MEDIUMTEXT"
< row
: bytes
| 302E303030
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
| CREATE TABLE t (a TEXT DEFAULT (-0e0))
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
|   `a` text DEFAULT -0e0
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "MEDIUMTEXT"
< row
: bytes
| 30
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` text DEFAULT -0e0
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "MEDIUMTEXT"
< row
: bytes
| 30
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
| CREATE TABLE t (a TEXT DEFAULT (007))
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
|   `a` text DEFAULT 7
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "MEDIUMTEXT"
< row
: bytes
| 37
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` text DEFAULT 7
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "MEDIUMTEXT"
< row
: bytes
| 37
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
| CREATE TABLE t (a TEXT DEFAULT (00.100))
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
|   `a` text DEFAULT 0.100
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "MEDIUMTEXT"
< row
: bytes
| 302E313030
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` text DEFAULT 0.100
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "MEDIUMTEXT"
< row
: bytes
| 302E313030
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
| CREATE TABLE t (a TEXT DEFAULT (+1))
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
|   `a` text DEFAULT 1
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "MEDIUMTEXT"
< row
: bytes
| 31
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` text DEFAULT 1
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "MEDIUMTEXT"
< row
: bytes
| 31
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
| CREATE TABLE t (a TEXT DEFAULT (+0))
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
|   `a` text DEFAULT 0
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "MEDIUMTEXT"
< row
: bytes
| 30
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` text DEFAULT 0
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "MEDIUMTEXT"
< row
: bytes
| 30
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
| CREATE TABLE t (a TEXT DEFAULT (0.10))
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
|   `a` text DEFAULT 0.10
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "MEDIUMTEXT"
< row
: bytes
| 302E3130
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` text DEFAULT 0.10
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "MEDIUMTEXT"
< row
: bytes
| 302E3130
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
| CREATE TABLE t (a TEXT DEFAULT (1.0000000000000000001))
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
|   `a` text DEFAULT 1.0000000000000000001
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "MEDIUMTEXT"
< row
: bytes
| 312E30303030303030303030303030303030303031
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` text DEFAULT 1.0000000000000000001
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "MEDIUMTEXT"
< row
: bytes
| 312E30303030303030303030303030303030303031
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
| CREATE TABLE t (a TEXT DEFAULT (1e2))
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
|   `a` text DEFAULT 1e2
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "MEDIUMTEXT"
< row
: bytes
| 313030
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` text DEFAULT 1e2
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "MEDIUMTEXT"
< row
: bytes
| 313030
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
| CREATE TABLE t (a TEXT DEFAULT (0e0))
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
|   `a` text DEFAULT 0e0
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "MEDIUMTEXT"
< row
: bytes
| 30
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` text DEFAULT 0e0
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "MEDIUMTEXT"
< row
: bytes
| 30
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
| CREATE TABLE t (a TEXT DEFAULT (0x41))
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
|   `a` text DEFAULT 0x41
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "MEDIUMTEXT"
< row
: bytes
| 41
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` text DEFAULT 0x41
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "MEDIUMTEXT"
< row
: bytes
| 41
