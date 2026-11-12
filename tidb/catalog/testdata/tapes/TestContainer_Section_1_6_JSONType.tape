-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
