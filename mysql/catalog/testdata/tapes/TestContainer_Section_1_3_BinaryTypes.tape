-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestContainer_Section_1_3_BinaryTypes default#1
> exec
| DROP TABLE IF EXISTS t_binary16
< ok 0 0
> exec
| CREATE TABLE t_binary16 (a BINARY(16))
< ok 0 0
> query
| SHOW CREATE TABLE t_binary16
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_binary16
: bytes
| CREATE TABLE `t_binary16` (
|   `a` binary(16) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_binary1
< ok 0 0
> exec
| CREATE TABLE t_binary1 (a BINARY)
< ok 0 0
> query
| SHOW CREATE TABLE t_binary1
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_binary1
: bytes
| CREATE TABLE `t_binary1` (
|   `a` binary(1) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_varbinary255
< ok 0 0
> exec
| CREATE TABLE t_varbinary255 (a VARBINARY(255))
< ok 0 0
> query
| SHOW CREATE TABLE t_varbinary255
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_varbinary255
: bytes
| CREATE TABLE `t_varbinary255` (
|   `a` varbinary(255) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_tinyblob
< ok 0 0
> exec
| CREATE TABLE t_tinyblob (a TINYBLOB)
< ok 0 0
> query
| SHOW CREATE TABLE t_tinyblob
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_tinyblob
: bytes
| CREATE TABLE `t_tinyblob` (
|   `a` tinyblob
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_blob
< ok 0 0
> exec
| CREATE TABLE t_blob (a BLOB)
< ok 0 0
> query
| SHOW CREATE TABLE t_blob
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_blob
: bytes
| CREATE TABLE `t_blob` (
|   `a` blob
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_mediumblob
< ok 0 0
> exec
| CREATE TABLE t_mediumblob (a MEDIUMBLOB)
< ok 0 0
> query
| SHOW CREATE TABLE t_mediumblob
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_mediumblob
: bytes
| CREATE TABLE `t_mediumblob` (
|   `a` mediumblob
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_longblob
< ok 0 0
> exec
| CREATE TABLE t_longblob (a LONGBLOB)
< ok 0 0
> query
| SHOW CREATE TABLE t_longblob
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_longblob
: bytes
| CREATE TABLE `t_longblob` (
|   `a` longblob
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_blob1000
< ok 0 0
> exec
| CREATE TABLE t_blob1000 (a BLOB(1000))
< ok 0 0
> query
| SHOW CREATE TABLE t_blob1000
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_blob1000
: bytes
| CREATE TABLE `t_blob1000` (
|   `a` blob
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
