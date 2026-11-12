-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
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
|   `a` tinyblob DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
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
|   `a` blob DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
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
|   `a` mediumblob DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
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
|   `a` longblob DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
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
|   `a` blob DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
