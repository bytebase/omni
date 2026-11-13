-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestSpatialFormsOracle default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| CREATE TABLE t (a GEOMETRY)
< error
| 1064 42000 …line 1 column 26 near "GEOMETRY)" 
> exec ^
: string
| CREATE TABLE t (a POINT)
< error
| 1064 42000 …line 1 column 23 near "POINT)" 
> exec ^
: string
| CREATE TABLE t (a LINESTRING)
< error
| 1064 42000 …line 1 column 28 near "LINESTRING)" 
> exec ^
: string
| CREATE TABLE t (a POLYGON)
< error
| 1064 42000 …line 1 column 25 near "POLYGON)" 
> exec ^
: string
| CREATE TABLE t (a MULTIPOINT)
< error
| 1064 42000 …line 1 column 28 near "MULTIPOINT)" 
> exec ^
: string
| CREATE TABLE t (a MULTILINESTRING)
< error
| 1064 42000 …line 1 column 33 near "MULTILINESTRING)" 
> exec ^
: string
| CREATE TABLE t (a MULTIPOLYGON)
< error
| 1064 42000 …line 1 column 30 near "MULTIPOLYGON)" 
> exec ^
: string
| CREATE TABLE t (a GEOMETRYCOLLECTION)
< error
| 1064 42000 …line 1 column 36 near "GEOMETRYCOLLECTION)" 
> exec ^
: string
| CREATE TABLE t (a GEOMCOLLECTION)
< error
| 1064 42000 …line 1 column 32 near "GEOMCOLLECTION)" 
> exec ^
: string
| CREATE TABLE t (a POINT NOT NULL SRID 4326)
< error
| 1064 42000 …line 1 column 23 near "POINT NOT NULL SRID 4326)" 
> exec ^
: string
| CREATE TABLE t (a INT SRID 4326)
< error
| 1064 42000 …line 1 column 26 near "SRID 4326)" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN a GEOMETRY
< error
| 1064 42000 …line 1 column 35 near "GEOMETRY" 
> exec ^
: string
| ALTER TABLE t MODIFY a POINT
< error
| 1064 42000 …line 1 column 28 near "POINT" 
> exec ^
: string
| SELECT CAST(a AS GEOMETRY) FROM t
< error
| 1064 42000 …line 1 column 25 near "GEOMETRY) FROM t" 
> exec ^
: string
| CREATE TABLE t (geometry INT, point INT, linestring INT, polygon INT)
< ok 0 0
> exec ^
: string
| CREATE TABLE geometry (multipoint INT, multilinestring INT, multipolygon INT)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (geometrycollection INT, geomcollection INT, srid INT)
< ok 0 0
> exec ^
: string
| SELECT point, polygon FROM geometry
< error
| 1146 42S02 Table 'test.geometry' doesn't exist
> exec ^
: string
| SELECT point, polygon FROM geometry
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE SPATIAL INDEX i ON t (a)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE SPATIAL INDEX i ON t (a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE SPATIAL INDEX IF NOT EXISTS i ON t (a) USING BTREE
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE SPATIAL INDEX IF NOT EXISTS i ON t (a) USING BTREE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE SPATIAL KEY i ON t (a)
< error
| 1064 42000 …line 1 column 18 near "KEY i ON t (a)" 
> exec ^
: string
| CREATE SPATIAL i ON t (a)
< error
| 1064 42000 …line 1 column 16 near "i ON t (a)" 
> exec ^
: string
| CREATE TEMPORARY SPATIAL INDEX i ON t (a)
< error
| 1064 42000 …line 1 column 24 near "SPATIAL INDEX i ON t (a)" 
> exec ^
: string
| CREATE FULLTEXT INDEX i ON t (a)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE FULLTEXT INDEX i ON t (a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE FULLTEXT KEY i ON t (a)
< error
| 1064 42000 …line 1 column 19 near "KEY i ON t (a)" 
> exec ^
: string
| CREATE TEMPORARY FULLTEXT INDEX i ON t (a)
< error
| 1064 42000 …line 1 column 25 near "FULLTEXT INDEX i ON t (a)" 
> exec ^
: string
| CREATE TEMPORARY UNIQUE INDEX i ON t (a)
< error
| 1064 42000 …line 1 column 23 near "UNIQUE INDEX i ON t (a)" 
> exec ^
: string
| CREATE TEMPORARY INDEX i ON t (a)
< error
| 1064 42000 …line 1 column 22 near "INDEX i ON t (a)" 
> exec ^
: string
| CREATE TABLE t (a INT, SPATIAL INDEX i (a))
< error
| 1064 42000 …line 1 column 30 near "SPATIAL INDEX i (a))" 
> exec ^
: string
| CREATE TABLE t (a INT, SPATIAL KEY i (a))
< error
| 1064 42000 …line 1 column 30 near "SPATIAL KEY i (a))" 
> exec ^
: string
| CREATE TABLE t (a INT, SPATIAL (a))
< error
| 1064 42000 …line 1 column 30 near "SPATIAL (a))" 
> exec ^
: string
| ALTER TABLE t ADD SPATIAL INDEX i (a)
< error
| 1064 42000 …line 1 column 25 near "SPATIAL INDEX i (a)" 
> exec ^
: string
| ALTER TABLE t ADD SPATIAL KEY i (a)
< error
| 1064 42000 …line 1 column 25 near "SPATIAL KEY i (a)" 
> exec ^
: string
| ALTER TABLE t ADD SPATIAL (a)
< error
| 1064 42000 …line 1 column 25 near "SPATIAL (a)" 
> exec ^
: string
| CREATE SPATIAL REFERENCE SYSTEM 4120 NAME 'x' DEFINITION 'y'
< error
| 1064 42000 …line 1 column 24 near "REFERENCE SYSTEM 4120 NAME 'x' DEFINITION 'y'" 
> exec ^
: string
| CREATE OR REPLACE SPATIAL REFERENCE SYSTEM 4120 NAME 'x' DEFINITION 'y'
< error
| 1064 42000 …line 1 column 25 near "SPATIAL REFERENCE SYSTEM 4120 NAME 'x' DEFINITION 'y'" 
> exec ^
: string
| DROP SPATIAL REFERENCE SYSTEM 4120
< error
| 1064 42000 …line 1 column 12 near "SPATIAL REFERENCE SYSTEM 4120" 
> exec ^
: string
| DROP SPATIAL REFERENCE SYSTEM IF EXISTS 4120
< error
| 1064 42000 …line 1 column 12 near "SPATIAL REFERENCE SYSTEM IF EXISTS 4120" 
> exec ^
: string
| SELECT ST_GeomFromText('POINT(1 1)')
< error
| 1305 42000 FUNCTION test.st_geomfromtext does not exist
> exec ^
: string
| SELECT ST_GeomFromText('POINT(1 1)')
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT ST_Distance(a, b) FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT ST_Distance(a, b) FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
