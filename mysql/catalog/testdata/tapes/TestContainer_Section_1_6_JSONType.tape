-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

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
|   `a` json DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
|   `a` json DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
