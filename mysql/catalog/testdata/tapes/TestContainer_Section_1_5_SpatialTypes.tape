-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestContainer_Section_1_5_SpatialTypes default#1
> exec
| DROP TABLE IF EXISTS t_geometry
< ok 0 0
> exec
| CREATE TABLE t_geometry (a GEOMETRY)
< ok 0 0
> query
| SHOW CREATE TABLE t_geometry
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_geometry
: bytes
| CREATE TABLE `t_geometry` (
|   `a` geometry DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_point
< ok 0 0
> exec
| CREATE TABLE t_point (a POINT)
< ok 0 0
> query
| SHOW CREATE TABLE t_point
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_point
: bytes
| CREATE TABLE `t_point` (
|   `a` point DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_linestring
< ok 0 0
> exec
| CREATE TABLE t_linestring (a LINESTRING)
< ok 0 0
> query
| SHOW CREATE TABLE t_linestring
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_linestring
: bytes
| CREATE TABLE `t_linestring` (
|   `a` linestring DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_polygon
< ok 0 0
> exec
| CREATE TABLE t_polygon (a POLYGON)
< ok 0 0
> query
| SHOW CREATE TABLE t_polygon
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_polygon
: bytes
| CREATE TABLE `t_polygon` (
|   `a` polygon DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_multipoint
< ok 0 0
> exec
| CREATE TABLE t_multipoint (a MULTIPOINT)
< ok 0 0
> query
| SHOW CREATE TABLE t_multipoint
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_multipoint
: bytes
| CREATE TABLE `t_multipoint` (
|   `a` multipoint DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_multilinestring
< ok 0 0
> exec
| CREATE TABLE t_multilinestring (a MULTILINESTRING)
< ok 0 0
> query
| SHOW CREATE TABLE t_multilinestring
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_multilinestring
: bytes
| CREATE TABLE `t_multilinestring` (
|   `a` multilinestring DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_multipolygon
< ok 0 0
> exec
| CREATE TABLE t_multipolygon (a MULTIPOLYGON)
< ok 0 0
> query
| SHOW CREATE TABLE t_multipolygon
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_multipolygon
: bytes
| CREATE TABLE `t_multipolygon` (
|   `a` multipolygon DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_geomcoll
< ok 0 0
> exec
| CREATE TABLE t_geomcoll (a GEOMETRYCOLLECTION)
< ok 0 0
> query
| SHOW CREATE TABLE t_geomcoll
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_geomcoll
: bytes
| CREATE TABLE `t_geomcoll` (
|   `a` geomcollection DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_point_srid
< ok 0 0
> exec
| CREATE TABLE t_point_srid (a POINT NOT NULL SRID 4326)
< ok 0 0
> query
| SHOW CREATE TABLE t_point_srid
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_point_srid
: bytes
| CREATE TABLE `t_point_srid` (
|   `a` point NOT NULL /*!80003 SRID 4326 */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_ls_srid
< ok 0 0
> exec
| CREATE TABLE t_ls_srid (a LINESTRING SRID 4326)
< ok 0 0
> query
| SHOW CREATE TABLE t_ls_srid
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_ls_srid
: bytes
| CREATE TABLE `t_ls_srid` (
|   `a` linestring /*!80003 SRID 4326 */ DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_poly_srid
< ok 0 0
> exec
| CREATE TABLE t_poly_srid (a POLYGON NOT NULL SRID 4326)
< ok 0 0
> query
| SHOW CREATE TABLE t_poly_srid
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_poly_srid
: bytes
| CREATE TABLE `t_poly_srid` (
|   `a` polygon NOT NULL /*!80003 SRID 4326 */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
