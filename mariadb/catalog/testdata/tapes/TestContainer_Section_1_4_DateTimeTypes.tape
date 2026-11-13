-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestContainer_Section_1_4_DateTimeTypes default#1
> exec
| DROP TABLE IF EXISTS t_date
< ok 0 0
> exec
| CREATE TABLE t_date (a DATE)
< ok 0 0
> query
| SHOW CREATE TABLE t_date
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_date
: bytes
| CREATE TABLE `t_date` (
|   `a` date DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_time
< ok 0 0
> exec
| CREATE TABLE t_time (a TIME)
< ok 0 0
> query
| SHOW CREATE TABLE t_time
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_time
: bytes
| CREATE TABLE `t_time` (
|   `a` time DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_time3
< ok 0 0
> exec
| CREATE TABLE t_time3 (a TIME(3))
< ok 0 0
> query
| SHOW CREATE TABLE t_time3
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_time3
: bytes
| CREATE TABLE `t_time3` (
|   `a` time(3) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_datetime
< ok 0 0
> exec
| CREATE TABLE t_datetime (a DATETIME)
< ok 0 0
> query
| SHOW CREATE TABLE t_datetime
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_datetime
: bytes
| CREATE TABLE `t_datetime` (
|   `a` datetime DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_datetime6
< ok 0 0
> exec
| CREATE TABLE t_datetime6 (a DATETIME(6))
< ok 0 0
> query
| SHOW CREATE TABLE t_datetime6
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_datetime6
: bytes
| CREATE TABLE `t_datetime6` (
|   `a` datetime(6) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_timestamp
< ok 0 0
> exec
| CREATE TABLE t_timestamp (a TIMESTAMP)
< ok 0 0
> query
| SHOW CREATE TABLE t_timestamp
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_timestamp
: bytes
| CREATE TABLE `t_timestamp` (
|   `a` timestamp NULL DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_timestamp3
< ok 0 0
> exec
| CREATE TABLE t_timestamp3 (a TIMESTAMP(3))
< ok 0 0
> query
| SHOW CREATE TABLE t_timestamp3
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_timestamp3
: bytes
| CREATE TABLE `t_timestamp3` (
|   `a` timestamp(3) NULL DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_year
< ok 0 0
> exec
| CREATE TABLE t_year (a YEAR)
< ok 0 0
> query
| SHOW CREATE TABLE t_year
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_year
: bytes
| CREATE TABLE `t_year` (
|   `a` year(4) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_year4
< ok 0 0
> exec
| CREATE TABLE t_year4 (a YEAR(4))
< ok 0 0
> query
| SHOW CREATE TABLE t_year4
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_year4
: bytes
| CREATE TABLE `t_year4` (
|   `a` year(4) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
