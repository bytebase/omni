-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracle_SRIDPresenceEmitsClause default#1
> exec
| DROP DATABASE IF EXISTS srid_emit_f
< ok 0 0
> exec
| CREATE DATABASE srid_emit_f
< ok 1 0
> exec
| USE srid_emit_f
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE srid_emit_f.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS srid_emit_t
< ok 0 0
> exec
| CREATE DATABASE srid_emit_t
< ok 1 0
> exec
| USE srid_emit_t
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0)
< ok 0 0
> query
| SHOW CREATE TABLE srid_emit_t.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL /*!80003 SRID 0 */,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS srid_emit_f
< ok 1 0
> exec
| CREATE DATABASE srid_emit_f
< ok 1 0
> exec
| USE srid_emit_f
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE srid_emit_f.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS srid_emit_t
< ok 1 0
> exec
| CREATE DATABASE srid_emit_t
< ok 1 0
> exec
| USE srid_emit_t
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 4326)
< ok 0 0
> query
| SHOW CREATE TABLE srid_emit_t.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL /*!80003 SRID 4326 */,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS srid_emit_f
< ok 1 0
> exec
| CREATE DATABASE srid_emit_f
< ok 1 0
> exec
| USE srid_emit_f
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0)
< ok 0 0
> query
| SHOW CREATE TABLE srid_emit_f.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL /*!80003 SRID 0 */,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS srid_emit_t
< ok 1 0
> exec
| CREATE DATABASE srid_emit_t
< ok 1 0
> exec
| USE srid_emit_t
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE srid_emit_t.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS srid_emit_f
< ok 1 0
> exec
| CREATE DATABASE srid_emit_f
< ok 1 0
> exec
| USE srid_emit_f
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 4326)
< ok 0 0
> query
| SHOW CREATE TABLE srid_emit_f.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL /*!80003 SRID 4326 */,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS srid_emit_t
< ok 1 0
> exec
| CREATE DATABASE srid_emit_t
< ok 1 0
> exec
| USE srid_emit_t
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0)
< ok 0 0
> query
| SHOW CREATE TABLE srid_emit_t.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL /*!80003 SRID 0 */,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
