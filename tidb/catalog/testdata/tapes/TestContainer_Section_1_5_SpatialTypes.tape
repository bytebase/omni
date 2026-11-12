-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestContainer_Section_1_5_SpatialTypes default#1
> exec
| DROP TABLE IF EXISTS t_geometry
< ok 0 0
> exec
| CREATE TABLE t_geometry (a GEOMETRY)
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 35 near "GEOMETRY)" 
> exec
| DROP TABLE IF EXISTS t_point
< ok 0 0
> exec
| CREATE TABLE t_point (a POINT)
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 29 near "POINT)" 
> exec
| DROP TABLE IF EXISTS t_linestring
< ok 0 0
> exec
| CREATE TABLE t_linestring (a LINESTRING)
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 39 near "LINESTRING)" 
> exec
| DROP TABLE IF EXISTS t_polygon
< ok 0 0
> exec
| CREATE TABLE t_polygon (a POLYGON)
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 33 near "POLYGON)" 
> exec
| DROP TABLE IF EXISTS t_multipoint
< ok 0 0
> exec
| CREATE TABLE t_multipoint (a MULTIPOINT)
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 39 near "MULTIPOINT)" 
> exec
| DROP TABLE IF EXISTS t_multilinestring
< ok 0 0
> exec
| CREATE TABLE t_multilinestring (a MULTILINESTRING)
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 49 near "MULTILINESTRING)" 
> exec
| DROP TABLE IF EXISTS t_multipolygon
< ok 0 0
> exec
| CREATE TABLE t_multipolygon (a MULTIPOLYGON)
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 43 near "MULTIPOLYGON)" 
> exec
| DROP TABLE IF EXISTS t_geomcoll
< ok 0 0
> exec
| CREATE TABLE t_geomcoll (a GEOMETRYCOLLECTION)
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 45 near "GEOMETRYCOLLECTION)" 
> exec
| DROP TABLE IF EXISTS t_point_srid
< ok 0 0
> exec
| CREATE TABLE t_point_srid (a POINT NOT NULL SRID 4326)
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 34 near "POINT NOT NULL SRID 4326)" 
> exec
| DROP TABLE IF EXISTS t_ls_srid
< ok 0 0
> exec
| CREATE TABLE t_ls_srid (a LINESTRING SRID 4326)
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 36 near "LINESTRING SRID 4326)" 
> exec
| DROP TABLE IF EXISTS t_poly_srid
< ok 0 0
> exec
| CREATE TABLE t_poly_srid (a POLYGON NOT NULL SRID 4326)
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 35 near "POLYGON NOT NULL SRID 4326)" 
