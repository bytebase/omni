-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestContainer_Section_1_6_JSONType default#1
> exec
| DROP TABLE IF EXISTS t_json
< ok 0 0
> exec
| CREATE TABLE t_json (a JSON)
< ok 0 0
> query
| SHOW CREATE TABLE t_json
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_json
: bytes
| CREATE TABLE `t_json` (
|   `a` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL CHECK (json_valid(`a`))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_json_dn
< ok 0 0
> exec
| CREATE TABLE t_json_dn (a JSON DEFAULT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE t_json_dn
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_json_dn
: bytes
| CREATE TABLE `t_json_dn` (
|   `a` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL CHECK (json_valid(`a`))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
