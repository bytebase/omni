-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestContainer_Section_1_8_OnUpdate default#1
> exec
| DROP TABLE IF EXISTS t_ou1
< ok 0 0
> exec
| CREATE TABLE t_ou1 (a TIMESTAMP ON UPDATE CURRENT_TIMESTAMP)
< ok 0 0
> query
| SHOW CREATE TABLE t_ou1
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_ou1
: bytes
| CREATE TABLE `t_ou1` (
|   `a` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp()
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_ou2
< ok 0 0
> exec
| CREATE TABLE t_ou2 (a DATETIME(3) ON UPDATE CURRENT_TIMESTAMP(3))
< ok 0 0
> query
| SHOW CREATE TABLE t_ou2
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_ou2
: bytes
| CREATE TABLE `t_ou2` (
|   `a` datetime(3) DEFAULT NULL ON UPDATE current_timestamp(3)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_ou3
< ok 0 0
> exec
| CREATE TABLE t_ou3 (a TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP)
< ok 0 0
> query
| SHOW CREATE TABLE t_ou3
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_ou3
: bytes
| CREATE TABLE `t_ou3` (
|   `a` timestamp NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_ou4
< ok 0 0
> exec
| CREATE TABLE t_ou4 (a DATETIME(6) DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6))
< ok 0 0
> query
| SHOW CREATE TABLE t_ou4
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_ou4
: bytes
| CREATE TABLE `t_ou4` (
|   `a` datetime(6) DEFAULT current_timestamp(6) ON UPDATE current_timestamp(6)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
