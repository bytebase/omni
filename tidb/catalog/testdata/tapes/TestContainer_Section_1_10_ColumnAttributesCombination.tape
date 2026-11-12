-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

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
|   `a` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`a`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
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
|   `a` bigint unsigned NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`a`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
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
|   `a` varchar(100) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
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
|   `a` int NOT NULL COMMENT 'user id'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
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
|   `a` int unsigned NOT NULL DEFAULT '0' COMMENT 'count'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
