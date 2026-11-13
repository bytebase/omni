-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracle_SRIDPresenceIdempotence default#1
> exec
| DROP DATABASE IF EXISTS srid_srid_zero
< ok 0 0
> exec
| CREATE DATABASE srid_srid_zero
< ok 1 0
> exec
| USE srid_srid_zero
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0)
< ok 0 0
> query
| SHOW CREATE TABLE srid_srid_zero.t
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
| DROP DATABASE IF EXISTS srid_srid_zero_idem
< ok 0 0
> exec
| CREATE DATABASE srid_srid_zero_idem
< ok 1 0
> exec
| USE srid_srid_zero_idem
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0)
< ok 0 0
> query
| SHOW CREATE TABLE srid_srid_zero_idem.t
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
| DROP DATABASE IF EXISTS srid_srid_4326
< ok 0 0
> exec
| CREATE DATABASE srid_srid_4326
< ok 1 0
> exec
| USE srid_srid_4326
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 4326)
< ok 0 0
> query
| SHOW CREATE TABLE srid_srid_4326.t
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
| DROP DATABASE IF EXISTS srid_srid_4326_idem
< ok 0 0
> exec
| CREATE DATABASE srid_srid_4326_idem
< ok 1 0
> exec
| USE srid_srid_4326_idem
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 4326)
< ok 0 0
> query
| SHOW CREATE TABLE srid_srid_4326_idem.t
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
| DROP DATABASE IF EXISTS srid_srid_none
< ok 0 0
> exec
| CREATE DATABASE srid_srid_none
< ok 1 0
> exec
| USE srid_srid_none
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE srid_srid_none.t
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
| DROP DATABASE IF EXISTS srid_srid_none_idem
< ok 0 0
> exec
| CREATE DATABASE srid_srid_none_idem
< ok 1 0
> exec
| USE srid_srid_none_idem
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE srid_srid_none_idem.t
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
| DROP DATABASE IF EXISTS srid_srid_zero_nullable
< ok 0 0
> exec
| CREATE DATABASE srid_srid_zero_nullable
< ok 1 0
> exec
| USE srid_srid_zero_nullable
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY SRID 0)
< ok 0 0
> query
| SHOW CREATE TABLE srid_srid_zero_nullable.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `g` geometry /*!80003 SRID 0 */ DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS srid_srid_zero_nullable_idem
< ok 0 0
> exec
| CREATE DATABASE srid_srid_zero_nullable_idem
< ok 1 0
> exec
| USE srid_srid_zero_nullable_idem
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY SRID 0)
< ok 0 0
> query
| SHOW CREATE TABLE srid_srid_zero_nullable_idem.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `g` geometry /*!80003 SRID 0 */ DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS srid_srid_zero_multi_type
< ok 0 0
> exec
| CREATE DATABASE srid_srid_zero_multi_type
< ok 1 0
> exec
| USE srid_srid_zero_multi_type
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, p POINT NOT NULL SRID 0, l LINESTRING NOT NULL SRID 0, y POLYGON NOT NULL SRID 0)
< ok 0 0
> query
| SHOW CREATE TABLE srid_srid_zero_multi_type.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `p` point NOT NULL /*!80003 SRID 0 */,
|   `l` linestring NOT NULL /*!80003 SRID 0 */,
|   `y` polygon NOT NULL /*!80003 SRID 0 */,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS srid_srid_zero_multi_type_idem
< ok 0 0
> exec
| CREATE DATABASE srid_srid_zero_multi_type_idem
< ok 1 0
> exec
| USE srid_srid_zero_multi_type_idem
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, p POINT NOT NULL SRID 0, l LINESTRING NOT NULL SRID 0, y POLYGON NOT NULL SRID 0)
< ok 0 0
> query
| SHOW CREATE TABLE srid_srid_zero_multi_type_idem.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `p` point NOT NULL /*!80003 SRID 0 */,
|   `l` linestring NOT NULL /*!80003 SRID 0 */,
|   `y` polygon NOT NULL /*!80003 SRID 0 */,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS srid_srid_zero_spatial_index
< ok 0 0
> exec
| CREATE DATABASE srid_srid_zero_spatial_index
< ok 1 0
> exec
| USE srid_srid_zero_spatial_index
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0, SPATIAL INDEX sp_g (g))
< ok 0 0
> query
| SHOW CREATE TABLE srid_srid_zero_spatial_index.t
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
|   PRIMARY KEY (`id`),
|   SPATIAL KEY `sp_g` (`g`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS srid_srid_zero_spatial_index_idem
< ok 0 0
> exec
| CREATE DATABASE srid_srid_zero_spatial_index_idem
< ok 1 0
> exec
| USE srid_srid_zero_spatial_index_idem
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0, SPATIAL INDEX sp_g (g))
< ok 0 0
> query
| SHOW CREATE TABLE srid_srid_zero_spatial_index_idem.t
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
|   PRIMARY KEY (`id`),
|   SPATIAL KEY `sp_g` (`g`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS srid_srid_4326_spatial_index
< ok 0 0
> exec
| CREATE DATABASE srid_srid_4326_spatial_index
< ok 1 0
> exec
| USE srid_srid_4326_spatial_index
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 4326, SPATIAL INDEX sp_g (g))
< ok 0 0
> query
| SHOW CREATE TABLE srid_srid_4326_spatial_index.t
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
|   PRIMARY KEY (`id`),
|   SPATIAL KEY `sp_g` (`g`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS srid_srid_4326_spatial_index_idem
< ok 0 0
> exec
| CREATE DATABASE srid_srid_4326_spatial_index_idem
< ok 1 0
> exec
| USE srid_srid_4326_spatial_index_idem
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 4326, SPATIAL INDEX sp_g (g))
< ok 0 0
> query
| SHOW CREATE TABLE srid_srid_4326_spatial_index_idem.t
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
|   PRIMARY KEY (`id`),
|   SPATIAL KEY `sp_g` (`g`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS srid_srid_mixed_none_and_zero
< ok 0 0
> exec
| CREATE DATABASE srid_srid_mixed_none_and_zero
< ok 1 0
> exec
| USE srid_srid_mixed_none_and_zero
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a GEOMETRY, b GEOMETRY NOT NULL SRID 0)
< ok 0 0
> query
| SHOW CREATE TABLE srid_srid_mixed_none_and_zero.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` geometry DEFAULT NULL,
|   `b` geometry NOT NULL /*!80003 SRID 0 */,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS srid_srid_mixed_none_and_zero_idem
< ok 0 0
> exec
| CREATE DATABASE srid_srid_mixed_none_and_zero_idem
< ok 1 0
> exec
| USE srid_srid_mixed_none_and_zero_idem
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a GEOMETRY, b GEOMETRY NOT NULL SRID 0)
< ok 0 0
> query
| SHOW CREATE TABLE srid_srid_mixed_none_and_zero_idem.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` geometry DEFAULT NULL,
|   `b` geometry NOT NULL /*!80003 SRID 0 */,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
