-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestContainer_Section_1_10_ColumnAttributesCombination default#1
> exec
| DROP TABLE IF EXISTS t_ai1
< ok 0 0
> exec
| CREATE TABLE t_ai1 (a INT NOT NULL AUTO_INCREMENT, PRIMARY KEY (a))
< ok 0 0
> query
| SHOW CREATE TABLE t_ai1
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_ai1
: bytes
| CREATE TABLE `t_ai1` (
|   `a` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_ai2
< ok 0 0
> exec
| CREATE TABLE t_ai2 (a BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, PRIMARY KEY (a))
< ok 0 0
> query
| SHOW CREATE TABLE t_ai2
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_ai2
: bytes
| CREATE TABLE `t_ai2` (
|   `a` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_vnde
< ok 0 0
> exec
| CREATE TABLE t_vnde (a VARCHAR(100) NOT NULL DEFAULT '')
< ok 0 0
> query
| SHOW CREATE TABLE t_vnde
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_vnde
: bytes
| CREATE TABLE `t_vnde` (
|   `a` varchar(100) NOT NULL DEFAULT ''
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_vccnn
< ok 0 0
> exec
| CREATE TABLE t_vccnn (a VARCHAR(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_bin NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE t_vccnn
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_vccnn
: bytes
| CREATE TABLE `t_vccnn` (
|   `a` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_bin NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_innc
< ok 0 0
> exec
| CREATE TABLE t_innc (a INT NOT NULL COMMENT 'user id')
< ok 0 0
> query
| SHOW CREATE TABLE t_innc
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_innc
: bytes
| CREATE TABLE `t_innc` (
|   `a` int(11) NOT NULL COMMENT 'user id'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_vinv
< ok 0 0
> exec
| CREATE TABLE t_vinv (a INT, b VARCHAR(255) INVISIBLE)
< ok 0 0
> query
| SHOW CREATE TABLE t_vinv
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_vinv
: bytes
| CREATE TABLE `t_vinv` (
|   `a` int(11) DEFAULT NULL,
|   `b` varchar(255) DEFAULT NULL INVISIBLE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_all
< ok 0 0
> exec
| CREATE TABLE t_all (a INT UNSIGNED NOT NULL DEFAULT '0' COMMENT 'count')
< ok 0 0
> query
| SHOW CREATE TABLE t_all
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_all
: bytes
| CREATE TABLE `t_all` (
|   `a` int(10) unsigned NOT NULL DEFAULT 0 COMMENT 'count'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
