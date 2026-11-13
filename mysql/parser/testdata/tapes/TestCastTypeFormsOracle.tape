-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestCastTypeFormsOracle default#1
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| SELECT CAST(a AS CHAR(4294967295))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR(4294967295))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHAR(4294967296))
< error
| 1439 42000 Display width out of range for column 'cast as char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS CHAR(4294967296))
| ; SELECT FROM
< error
| 1439 42000 Display width out of range for column 'cast as char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS BINARY(4294967295))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS BINARY(4294967295))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS BINARY(4294967296))
< error
| 1439 42000 Display width out of range for column 'cast as char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS BINARY(4294967296))
| ; SELECT FROM
< error
| 1439 42000 Display width out of range for column 'cast as char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS NCHAR(4294967295))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS NCHAR(4294967295))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS NCHAR(4294967296))
< error
| 1439 42000 Display width out of range for column 'cast as char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS NCHAR(4294967296))
| ; SELECT FROM
< error
| 1439 42000 Display width out of range for column 'cast as char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS NATIONAL CHAR(4294967296))
< error
| 1439 42000 Display width out of range for column 'cast as char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS NATIONAL CHAR(4294967296))
| ; SELECT FROM
< error
| 1439 42000 Display width out of range for column 'cast as char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS CHAR(4294967296.5))
< error
| 1439 42000 Display width out of range for column 'cast as char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS CHAR(4294967296.5))
| ; SELECT FROM
< error
| 1439 42000 Display width out of range for column 'cast as char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS CHAR(4294967296) BYTE)
< error
| 1439 42000 Display width out of range for column 'cast as char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS CHAR(4294967296) BYTE)
| ; SELECT FROM
< error
| 1439 42000 Display width out of range for column 'cast as char' (max = 4294967295)
> exec ^
: string
| SELECT CONVERT(a, CHAR(4294967296))
< error
| 1439 42000 Display width out of range for column 'cast as char' (max = 4294967295)
> exec ^
: string
| SELECT CONVERT(a, CHAR(4294967296))
| ; SELECT FROM
< error
| 1439 42000 Display width out of range for column 'cast as char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS DECIMAL(65,30))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DECIMAL(65,30))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS DECIMAL(65,31))
< error
| 1425 42000 Too big scale 31 specified for column 'a'. Maximum is 30.
> exec ^
: string
| SELECT CAST(a AS DECIMAL(65,31))
| ; SELECT FROM
< error
| 1425 42000 Too big scale 31 specified for column 'a'. Maximum is 30.
> exec ^
: string
| SELECT CAST(a AS DECIMAL(31,31))
< error
| 1425 42000 Too big scale 31 specified for column 'a'. Maximum is 30.
> exec ^
: string
| SELECT CAST(a AS DECIMAL(31,31))
| ; SELECT FROM
< error
| 1425 42000 Too big scale 31 specified for column 'a'. Maximum is 30.
> exec ^
: string
| SELECT CAST(a AS DECIMAL(10,30))
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column '').
> exec ^
: string
| SELECT CAST(a AS DECIMAL(10,30))
| ; SELECT FROM
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column '').
> exec ^
: string
| SELECT CAST(a AS DECIMAL(5,6))
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column '').
> exec ^
: string
| SELECT CAST(a AS DECIMAL(5,6))
| ; SELECT FROM
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column '').
> exec ^
: string
| SELECT CAST(a AS DECIMAL(1,1))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DECIMAL(1,1))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS DECIMAL(1,2))
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column '').
> exec ^
: string
| SELECT CAST(a AS DECIMAL(1,2))
| ; SELECT FROM
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column '').
> exec ^
: string
| SELECT CAST(a AS DECIMAL(0,1))
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column '').
> exec ^
: string
| SELECT CAST(a AS DECIMAL(0,1))
| ; SELECT FROM
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column '').
> exec ^
: string
| SELECT CAST(a AS DECIMAL(0,0))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DECIMAL(0,0))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS DECIMAL(65))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DECIMAL(65))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS DECIMAL(66))
< error
| 1426 42000 Too-big precision 66 specified for 'a'. Maximum is 65.
> exec ^
: string
| SELECT CAST(a AS DECIMAL(66))
| ; SELECT FROM
< error
| 1426 42000 Too-big precision 66 specified for 'a'. Maximum is 65.
> exec ^
: string
| SELECT CAST(a AS DECIMAL(66,1))
< error
| 1426 42000 Too-big precision 66 specified for 'a'. Maximum is 65.
> exec ^
: string
| SELECT CAST(a AS DECIMAL(66,1))
| ; SELECT FROM
< error
| 1426 42000 Too-big precision 66 specified for 'a'. Maximum is 65.
> exec ^
: string
| SELECT CAST(a AS DECIMAL(65.5))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DECIMAL(65.5))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS DECIMAL(66.5))
< error
| 1426 42000 Too-big precision 66 specified for 'a'. Maximum is 65.
> exec ^
: string
| SELECT CAST(a AS DECIMAL(66.5))
| ; SELECT FROM
< error
| 1426 42000 Too-big precision 66 specified for 'a'. Maximum is 65.
> exec ^
: string
| SELECT CAST(a AS DEC(66))
< error
| 1426 42000 Too-big precision 66 specified for 'a'. Maximum is 65.
> exec ^
: string
| SELECT CAST(a AS DEC(66))
| ; SELECT FROM
< error
| 1426 42000 Too-big precision 66 specified for 'a'. Maximum is 65.
> exec ^
: string
| SELECT CONVERT(a, DECIMAL(66))
< error
| 1426 42000 Too-big precision 66 specified for 'a'. Maximum is 65.
> exec ^
: string
| SELECT CONVERT(a, DECIMAL(66))
| ; SELECT FROM
< error
| 1426 42000 Too-big precision 66 specified for 'a'. Maximum is 65.
> exec ^
: string
| SELECT CAST(a AS FLOAT(53))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS FLOAT(53))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS FLOAT(53.5))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS FLOAT(53.5))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS FLOAT(54))
< error
| 1426 42000 Too-big precision 54 specified for 'CAST'. Maximum is 53.
> exec ^
: string
| SELECT CAST(a AS FLOAT(54))
| ; SELECT FROM
< error
| 1426 42000 Too-big precision 54 specified for 'CAST'. Maximum is 53.
> exec ^
: string
| SELECT CAST(a AS FLOAT(54.5))
< error
| 1426 42000 Too-big precision 54 specified for 'CAST'. Maximum is 53.
> exec ^
: string
| SELECT CAST(a AS FLOAT(54.5))
| ; SELECT FROM
< error
| 1426 42000 Too-big precision 54 specified for 'CAST'. Maximum is 53.
> exec ^
: string
| SELECT CONVERT(a, FLOAT(54))
< error
| 1426 42000 Too-big precision 54 specified for 'CAST'. Maximum is 53.
> exec ^
: string
| SELECT CONVERT(a, FLOAT(54))
| ; SELECT FROM
< error
| 1426 42000 Too-big precision 54 specified for 'CAST'. Maximum is 53.
> exec ^
: string
| SELECT CAST(a AS TIME(6))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS TIME(6))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS TIME(7))
< error
| 1426 42000 Too-big precision 7 specified for 'CAST'. Maximum is 6.
> exec ^
: string
| SELECT CAST(a AS TIME(7))
| ; SELECT FROM
< error
| 1426 42000 Too-big precision 7 specified for 'CAST'. Maximum is 6.
> exec ^
: string
| SELECT CAST(a AS DATETIME(6))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DATETIME(6))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS DATETIME(7))
< error
| 1426 42000 Too-big precision 7 specified for 'CAST'. Maximum is 6.
> exec ^
: string
| SELECT CAST(a AS DATETIME(7))
| ; SELECT FROM
< error
| 1426 42000 Too-big precision 7 specified for 'CAST'. Maximum is 6.
> exec ^
: string
| SELECT CONVERT(a, TIME(7))
< error
| 1426 42000 Too-big precision 7 specified for 'CAST'. Maximum is 6.
> exec ^
: string
| SELECT CONVERT(a, TIME(7))
| ; SELECT FROM
< error
| 1426 42000 Too-big precision 7 specified for 'CAST'. Maximum is 6.
> exec ^
: string
| CREATE TABLE t (a INT, c INT AS (CAST(a AS FLOAT(54))))
< error
| 1426 42000 Too-big precision 54 specified for 'CAST'. Maximum is 53.
> exec ^
: string
| CREATE TABLE t (a INT, c INT AS (CAST(a AS FLOAT(54))))
| ; SELECT FROM
< error
| 1426 42000 Too-big precision 54 specified for 'CAST'. Maximum is 53.
> exec ^
: string
| SELECT CAST(a AS CHAR(0) ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec ^
: string
| SELECT CAST(a AS CHAR(0) ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHARACTER(3) ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec ^
: string
| SELECT CAST(a AS CHARACTER(3) ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHAR(3) BYTE ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec ^
: string
| SELECT CAST(a AS CHAR(3) BYTE ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS DATETIME(3) ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec ^
: string
| SELECT CAST(a AS DATETIME(3) ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS TIME(3) ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec ^
: string
| SELECT CAST(a AS TIME(3) ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS DECIMAL(5,2) ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec ^
: string
| SELECT CAST(a AS DECIMAL(5,2) ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS UNSIGNED INT ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec ^
: string
| SELECT CAST(a AS UNSIGNED INT ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHARACTER ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of char/binary BLOBs'
> exec ^
: string
| SELECT CAST(a AS CHARACTER ARRAY)
| ; SELECT FROM
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of char/binary BLOBs'
> exec ^
: string
| SELECT CAST(a AS CHAR BYTE ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of char/binary BLOBs'
> exec ^
: string
| SELECT CAST(a AS CHAR BYTE ARRAY)
| ; SELECT FROM
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of char/binary BLOBs'
> exec ^
: string
| SELECT CAST(a AS CHAR CHARACTER SET utf8mb4 ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'specifying charset for multi-valued index'
> exec ^
: string
| SELECT CAST(a AS CHAR CHARACTER SET utf8mb4 ARRAY)
| ; SELECT FROM
< error
| 1235 42000 This version of MySQL doesn't yet support 'specifying charset for multi-valued index'
> exec ^
: string
| SELECT CAST(a AS CHAR(3) CHARACTER SET utf8mb4 ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'specifying charset for multi-valued index'
> exec ^
: string
| SELECT CAST(a AS CHAR(3) CHARACTER SET utf8mb4 ARRAY)
| ; SELECT FROM
< error
| 1235 42000 This version of MySQL doesn't yet support 'specifying charset for multi-valued index'
> exec ^
: string
| SELECT CAST(a AS CHAR(3) CHARACTER SET binary ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec ^
: string
| SELECT CAST(a AS CHAR(3) CHARACTER SET binary ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHAR(3) CHARACTER SET BINARY ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec ^
: string
| SELECT CAST(a AS CHAR(3) CHARACTER SET BINARY ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHAR(3) CHARSET BiNaRy ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec ^
: string
| SELECT CAST(a AS CHAR(3) CHARSET BiNaRy ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHAR /* ( */ CHARACTER SET binary ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of char/binary BLOBs'
> exec ^
: string
| SELECT CAST(a AS CHAR /* ( */ CHARACTER SET binary ARRAY)
| ; SELECT FROM
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of char/binary BLOBs'
> exec ^
: string
| SELECT CAST(a AS BINARY /* ( */ ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of char/binary BLOBs'
> exec ^
: string
| SELECT CAST(a AS BINARY /* ( */ ARRAY)
| ; SELECT FROM
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of char/binary BLOBs'
> exec ^
: string
| SELECT CAST(a AS BINARY(3)), CAST(a AS BINARY ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of char/binary BLOBs'
> exec ^
: string
| SELECT CAST(a AS BINARY(3)), CAST(a AS BINARY ARRAY)
| ; SELECT FROM
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of char/binary BLOBs'
> exec ^
: string
| SELECT CAST(a AS CHAR(3) BINARY CHARACTER SET binary ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec ^
: string
| SELECT CAST(a AS CHAR(3) BINARY CHARACTER SET binary ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHAR(3) ASCII ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'specifying charset for multi-valued index'
> exec ^
: string
| SELECT CAST(a AS CHAR(3) ASCII ARRAY)
| ; SELECT FROM
< error
| 1235 42000 This version of MySQL doesn't yet support 'specifying charset for multi-valued index'
> exec ^
: string
| SELECT CAST(a AS CHAR(3) UNICODE ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'specifying charset for multi-valued index'
> exec ^
: string
| SELECT CAST(a AS CHAR(3) UNICODE ARRAY)
| ; SELECT FROM
< error
| 1235 42000 This version of MySQL doesn't yet support 'specifying charset for multi-valued index'
> exec ^
: string
| SELECT CAST(a AS CHAR(3) BINARY ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'specifying charset for multi-valued index'
> exec ^
: string
| SELECT CAST(a AS CHAR(3) BINARY ARRAY)
| ; SELECT FROM
< error
| 1235 42000 This version of MySQL doesn't yet support 'specifying charset for multi-valued index'
> exec ^
: string
| SELECT CAST(a AS NCHAR(3) ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'specifying charset for multi-valued index'
> exec ^
: string
| SELECT CAST(a AS NCHAR(3) ARRAY)
| ; SELECT FROM
< error
| 1235 42000 This version of MySQL doesn't yet support 'specifying charset for multi-valued index'
> exec ^
: string
| SELECT CAST(a AS NATIONAL CHAR(3) ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'specifying charset for multi-valued index'
> exec ^
: string
| SELECT CAST(a AS NATIONAL CHAR(3) ARRAY)
| ; SELECT FROM
< error
| 1235 42000 This version of MySQL doesn't yet support 'specifying charset for multi-valued index'
> exec ^
: string
| SELECT CAST(a AS NATIONAL CHAR ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'specifying charset for multi-valued index'
> exec ^
: string
| SELECT CAST(a AS NATIONAL CHAR ARRAY)
| ; SELECT FROM
< error
| 1235 42000 This version of MySQL doesn't yet support 'specifying charset for multi-valued index'
> exec ^
: string
| SELECT CAST(a AS DOUBLE PRECISION ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of DOUBLE'
> exec ^
: string
| SELECT CAST(a AS DOUBLE PRECISION ARRAY)
| ; SELECT FROM
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of DOUBLE'
> exec ^
: string
| SELECT CAST(a AS FLOAT(3) ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of FLOAT'
> exec ^
: string
| SELECT CAST(a AS FLOAT(3) ARRAY)
| ; SELECT FROM
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of FLOAT'
> exec ^
: string
| SELECT CAST(a AS LINESTRING ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of LINESTRING'
> exec ^
: string
| SELECT CAST(a AS LINESTRING ARRAY)
| ; SELECT FROM
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of LINESTRING'
> exec ^
: string
| SELECT CAST(a AS POLYGON ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of POLYGON'
> exec ^
: string
| SELECT CAST(a AS POLYGON ARRAY)
| ; SELECT FROM
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of POLYGON'
> exec ^
: string
| SELECT CAST(a AS MULTIPOINT ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of MULTIPOINT'
> exec ^
: string
| SELECT CAST(a AS MULTIPOINT ARRAY)
| ; SELECT FROM
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of MULTIPOINT'
> exec ^
: string
| SELECT CAST(a AS MULTILINESTRING ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of MULTILINESTRING>'
> exec ^
: string
| SELECT CAST(a AS MULTILINESTRING ARRAY)
| ; SELECT FROM
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of MULTILINESTRING>'
> exec ^
: string
| SELECT CAST(a AS MULTIPOLYGON ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of MULTIPOLYGON'
> exec ^
: string
| SELECT CAST(a AS MULTIPOLYGON ARRAY)
| ; SELECT FROM
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of MULTIPOLYGON'
> exec ^
: string
| SELECT CAST(a AS GEOMCOLLECTION ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of GEOMETRYCOLLECTION'
> exec ^
: string
| SELECT CAST(a AS GEOMCOLLECTION ARRAY)
| ; SELECT FROM
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of GEOMETRYCOLLECTION'
> exec ^
: string
| SELECT CAST(a AS GEOMETRYCOLLECTION ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of GEOMETRYCOLLECTION'
> exec ^
: string
| SELECT CAST(a AS GEOMETRYCOLLECTION ARRAY)
| ; SELECT FROM
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of GEOMETRYCOLLECTION'
> exec ^
: string
| CREATE TABLE t (a INT, KEY k ((CAST(a AS CHAR ARRAY))))
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of char/binary BLOBs'
> exec ^
: string
| CREATE TABLE t (a INT, KEY k ((CAST(a AS CHAR ARRAY))))
| ; SELECT FROM
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of char/binary BLOBs'
> exec ^
: string
| CREATE TABLE t (a INT, KEY k ((CAST(a AS JSON ARRAY))))
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of JSON'
> exec ^
: string
| CREATE TABLE t (a INT, KEY k ((CAST(a AS JSON ARRAY))))
| ; SELECT FROM
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of JSON'
> exec ^
: string
| CREATE TABLE t (a INT, KEY k ((CAST(a AS CHAR(3) ARRAY))))
< ok 0 0
> exec ^
: string
| SELECT CAST(a AS INT)
< error
| 1064 42000 …near 'INT)' at line 1
> exec ^
: string
| SELECT CAST(a AS INTEGER)
< error
| 1064 42000 …near 'INTEGER)' at line 1
> exec ^
: string
| SELECT CAST(a AS TINYINT)
< error
| 1064 42000 …near 'TINYINT)' at line 1
> exec ^
: string
| SELECT CAST(a AS SMALLINT)
< error
| 1064 42000 …near 'SMALLINT)' at line 1
> exec ^
: string
| SELECT CAST(a AS MEDIUMINT)
< error
| 1064 42000 …near 'MEDIUMINT)' at line 1
> exec ^
: string
| SELECT CAST(a AS BIGINT)
< error
| 1064 42000 …near 'BIGINT)' at line 1
> exec ^
: string
| SELECT CAST(a AS INT(3))
< error
| 1064 42000 …near 'INT(3))' at line 1
> exec ^
: string
| SELECT CAST(a AS SIGNED)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS SIGNED)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS SIGNED INT)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS SIGNED INT)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS SIGNED INTEGER)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS SIGNED INTEGER)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS UNSIGNED)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS UNSIGNED)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS UNSIGNED INT)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS UNSIGNED INT)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS UNSIGNED INTEGER)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS UNSIGNED INTEGER)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS SIGNED(3))
< error
| 1064 42000 …near '(3))' at line 1
> exec ^
: string
| SELECT CAST(a AS INT UNSIGNED)
< error
| 1064 42000 …near 'INT UNSIGNED)' at line 1
> exec ^
: string
| SELECT CAST(a AS INT SIGNED)
< error
| 1064 42000 …near 'INT SIGNED)' at line 1
> exec ^
: string
| SELECT CAST(a AS FLOAT)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS FLOAT)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS FLOAT(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS FLOAT(3))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS FLOAT(3,2))
< error
| 1064 42000 …near ',2))' at line 1
> exec ^
: string
| SELECT CAST(a AS DOUBLE)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DOUBLE)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS DOUBLE(3,2))
< error
| 1064 42000 …near '(3,2))' at line 1
> exec ^
: string
| SELECT CAST(a AS DOUBLE PRECISION)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DOUBLE PRECISION)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS REAL)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS REAL)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS DECIMAL)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DECIMAL)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS DECIMAL(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DECIMAL(3))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS DECIMAL(3,2))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DECIMAL(3,2))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS NUMERIC)
< error
| 1064 42000 …near 'NUMERIC)' at line 1
> exec ^
: string
| SELECT CAST(a AS NUMERIC(3,2))
< error
| 1064 42000 …near 'NUMERIC(3,2))' at line 1
> exec ^
: string
| SELECT CAST(a AS DEC)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DEC)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS DEC(3,2))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DEC(3,2))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS FIXED)
< error
| 1064 42000 …near 'FIXED)' at line 1
> exec ^
: string
| SELECT CAST(a AS FIXED(3,2))
< error
| 1064 42000 …near 'FIXED(3,2))' at line 1
> exec ^
: string
| SELECT CAST(a AS DECIMAL UNSIGNED)
< error
| 1064 42000 …near 'UNSIGNED)' at line 1
> exec ^
: string
| SELECT CAST(a AS DOUBLE UNSIGNED)
< error
| 1064 42000 …near 'UNSIGNED)' at line 1
> exec ^
: string
| SELECT CAST(a AS FLOAT UNSIGNED)
< error
| 1064 42000 …near 'UNSIGNED)' at line 1
> exec ^
: string
| SELECT CAST(a AS BOOL)
< error
| 1064 42000 …near 'BOOL)' at line 1
> exec ^
: string
| SELECT CAST(a AS BOOLEAN)
< error
| 1064 42000 …near 'BOOLEAN)' at line 1
> exec ^
: string
| SELECT CAST(a AS CHAR)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHAR(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR(3))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHARACTER)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHARACTER)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHARACTER(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHARACTER(3))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS NCHAR)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS NCHAR)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS NCHAR(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS NCHAR(3))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS NATIONAL CHAR)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS NATIONAL CHAR)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS NATIONAL CHAR(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS NATIONAL CHAR(3))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHAR VARYING(3))
< error
| 1064 42000 …near 'VARYING(3))' at line 1
> exec ^
: string
| SELECT CAST(a AS CHARACTER VARYING(3))
< error
| 1064 42000 …near 'VARYING(3))' at line 1
> exec ^
: string
| SELECT CAST(a AS VARCHAR(3))
< error
| 1064 42000 …near 'VARCHAR(3))' at line 1
> exec ^
: string
| SELECT CAST(a AS VARCHAR)
< error
| 1064 42000 …near 'VARCHAR)' at line 1
> exec ^
: string
| SELECT CAST(a AS VARCHARACTER(3))
< error
| 1064 42000 …near 'VARCHARACTER(3))' at line 1
> exec ^
: string
| SELECT CAST(a AS NVARCHAR(3))
< error
| 1064 42000 …near 'NVARCHAR(3))' at line 1
> exec ^
: string
| SELECT CAST(a AS NATIONAL VARCHAR(3))
< error
| 1064 42000 …near 'VARCHAR(3))' at line 1
> exec ^
: string
| SELECT CAST(a AS NCHAR VARCHAR(3))
< error
| 1064 42000 …near 'VARCHAR(3))' at line 1
> exec ^
: string
| SELECT CAST(a AS TEXT)
< error
| 1064 42000 …near 'TEXT)' at line 1
> exec ^
: string
| SELECT CAST(a AS TEXT(3))
< error
| 1064 42000 …near 'TEXT(3))' at line 1
> exec ^
: string
| SELECT CAST(a AS TINYTEXT)
< error
| 1064 42000 …near 'TINYTEXT)' at line 1
> exec ^
: string
| SELECT CAST(a AS MEDIUMTEXT)
< error
| 1064 42000 …near 'MEDIUMTEXT)' at line 1
> exec ^
: string
| SELECT CAST(a AS LONGTEXT)
< error
| 1064 42000 …near 'LONGTEXT)' at line 1
> exec ^
: string
| SELECT CAST(a AS LONG)
< error
| 1064 42000 …near 'LONG)' at line 1
> exec ^
: string
| SELECT CAST(a AS LONG VARCHAR)
< error
| 1064 42000 …near 'LONG VARCHAR)' at line 1
> exec ^
: string
| SELECT CAST(a AS BLOB)
< error
| 1064 42000 …near 'BLOB)' at line 1
> exec ^
: string
| SELECT CAST(a AS BLOB(3))
< error
| 1064 42000 …near 'BLOB(3))' at line 1
> exec ^
: string
| SELECT CAST(a AS TINYBLOB)
< error
| 1064 42000 …near 'TINYBLOB)' at line 1
> exec ^
: string
| SELECT CAST(a AS MEDIUMBLOB)
< error
| 1064 42000 …near 'MEDIUMBLOB)' at line 1
> exec ^
: string
| SELECT CAST(a AS LONGBLOB)
< error
| 1064 42000 …near 'LONGBLOB)' at line 1
> exec ^
: string
| SELECT CAST(a AS BINARY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS BINARY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS BINARY(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS BINARY(3))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS VARBINARY(3))
< error
| 1064 42000 …near 'VARBINARY(3))' at line 1
> exec ^
: string
| SELECT CAST(a AS VARBINARY)
< error
| 1064 42000 …near 'VARBINARY)' at line 1
> exec ^
: string
| SELECT CAST(a AS BIT)
< error
| 1064 42000 …near 'BIT)' at line 1
> exec ^
: string
| SELECT CAST(a AS BIT(1))
< error
| 1064 42000 …near 'BIT(1))' at line 1
> exec ^
: string
| SELECT CAST(a AS DATE)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DATE)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS DATE(3))
< error
| 1064 42000 …near '(3))' at line 1
> exec ^
: string
| SELECT CAST(a AS TIME)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS TIME)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS TIME(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS TIME(3))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS DATETIME)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DATETIME)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS DATETIME(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DATETIME(3))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS TIMESTAMP)
< error
| 1064 42000 …near 'TIMESTAMP)' at line 1
> exec ^
: string
| SELECT CAST(a AS TIMESTAMP(3))
< error
| 1064 42000 …near 'TIMESTAMP(3))' at line 1
> exec ^
: string
| SELECT CAST(a AS YEAR)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS YEAR)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS YEAR(4))
< error
| 1064 42000 …near '(4))' at line 1
> exec ^
: string
| SELECT CAST(a AS JSON)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS JSON)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS JSON(3))
< error
| 1064 42000 …near '(3))' at line 1
> exec ^
: string
| SELECT CAST(a AS ENUM('a'))
< error
| 1064 42000 …near 'ENUM('a'))' at line 1
> exec ^
: string
| SELECT CAST(a AS SET('a'))
< error
| 1064 42000 …near 'SET('a'))' at line 1
> exec ^
: string
| SELECT CAST(a AS GEOMETRY)
< error
| 1064 42000 …near 'GEOMETRY)' at line 1
> exec ^
: string
| SELECT CAST(a AS POINT)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS POINT)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS LINESTRING)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS LINESTRING)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS POLYGON)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS POLYGON)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS UUID)
< error
| 1064 42000 …near 'UUID)' at line 1
> exec ^
: string
| SELECT CAST(a AS INET4)
< error
| 1064 42000 …near 'INET4)' at line 1
> exec ^
: string
| SELECT CAST(a AS INET6)
< error
| 1064 42000 …near 'INET6)' at line 1
> exec ^
: string
| SELECT CAST(a AS SERIAL)
< error
| 1064 42000 …near 'SERIAL)' at line 1
> exec ^
: string
| SELECT CAST(a AS CHAR CHARACTER SET utf8mb4)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR CHARACTER SET utf8mb4)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHAR CHARSET utf8mb4)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR CHARSET utf8mb4)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHAR(3) CHARACTER SET utf8mb4)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR(3) CHARACTER SET utf8mb4)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHAR ASCII)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR ASCII)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHAR UNICODE)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR UNICODE)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHAR BYTE)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR BYTE)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHAR BINARY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR BINARY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHAR(3) BINARY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR(3) BINARY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS BINARY CHARACTER SET utf8mb4)
< error
| 1064 42000 …near 'CHARACTER SET utf8mb4)' at line 1
> exec ^
: string
| SELECT CAST(a AS CHAR COLLATE utf8mb4_bin)
< error
| 1064 42000 …near 'COLLATE utf8mb4_bin)' at line 1
> exec ^
: string
| SELECT CAST(a AS VARCHAR(3) CHARACTER SET utf8mb4)
< error
| 1064 42000 …near 'VARCHAR(3) CHARACTER SET utf8mb4)' at line 1
> exec ^
: string
| SELECT CAST(a AS INTERVAL DAY)
< error
| 1064 42000 …near 'INTERVAL DAY)' at line 1
> exec ^
: string
| SELECT CAST(a AS DAY)
< error
| 1064 42000 …near 'DAY)' at line 1
> exec ^
: string
| SELECT CAST(a AS VECTOR)
< error
| 1064 42000 …near 'VECTOR)' at line 1
> exec ^
: string
| SELECT CAST(a AS VECTOR(3))
< error
| 1064 42000 …near 'VECTOR(3))' at line 1
> exec ^
: string
| SELECT CONVERT(a, INT)
< error
| 1064 42000 …near 'INT)' at line 1
> exec ^
: string
| SELECT CONVERT(a, SIGNED)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CONVERT(a, SIGNED)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CONVERT(a, CHAR(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CONVERT(a, CHAR(3))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CONVERT(a, DATE)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CONVERT(a, DATE)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CONVERT(a, DECIMAL(3,2))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CONVERT(a, DECIMAL(3,2))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CONVERT(a, BIGINT)
< error
| 1064 42000 …near 'BIGINT)' at line 1
> exec ^
: string
| SELECT CONVERT(a, TEXT)
< error
| 1064 42000 …near 'TEXT)' at line 1
> exec ^
: string
| SELECT CONVERT(a, VARCHAR(3))
< error
| 1064 42000 …near 'VARCHAR(3))' at line 1
> exec ^
: string
| SELECT CONVERT(a, BOOLEAN)
< error
| 1064 42000 …near 'BOOLEAN)' at line 1
> exec ^
: string
| SELECT CONVERT(a, JSON)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CONVERT(a, JSON)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CONVERT(a, DOUBLE)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CONVERT(a, DOUBLE)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CONVERT(a, FLOAT)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CONVERT(a, FLOAT)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CONVERT(a, TIME(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CONVERT(a, TIME(3))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS INT ARRAY)
< error
| 1064 42000 …near 'INT ARRAY)' at line 1
> exec ^
: string
| SELECT CAST(a AS BIGINT ARRAY)
< error
| 1064 42000 …near 'BIGINT ARRAY)' at line 1
> exec ^
: string
| SELECT CAST(a AS SIGNED ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec ^
: string
| SELECT CAST(a AS SIGNED ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS UNSIGNED ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec ^
: string
| SELECT CAST(a AS UNSIGNED ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS SIGNED INTEGER ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec ^
: string
| SELECT CAST(a AS SIGNED INTEGER ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHAR ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of char/binary BLOBs'
> exec ^
: string
| SELECT CAST(a AS CHAR ARRAY)
| ; SELECT FROM
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of char/binary BLOBs'
> exec ^
: string
| SELECT CAST(a AS CHAR(3) ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec ^
: string
| SELECT CAST(a AS CHAR(3) ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS BINARY ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of char/binary BLOBs'
> exec ^
: string
| SELECT CAST(a AS BINARY ARRAY)
| ; SELECT FROM
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of char/binary BLOBs'
> exec ^
: string
| SELECT CAST(a AS BINARY(3) ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec ^
: string
| SELECT CAST(a AS BINARY(3) ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS DATE ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec ^
: string
| SELECT CAST(a AS DATE ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS DATETIME ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec ^
: string
| SELECT CAST(a AS DATETIME ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS TIME ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec ^
: string
| SELECT CAST(a AS TIME ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS YEAR ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of YEAR'
> exec ^
: string
| SELECT CAST(a AS YEAR ARRAY)
| ; SELECT FROM
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of YEAR'
> exec ^
: string
| SELECT CAST(a AS JSON ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of JSON'
> exec ^
: string
| SELECT CAST(a AS JSON ARRAY)
| ; SELECT FROM
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of JSON'
> exec ^
: string
| SELECT CAST(a AS DECIMAL ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec ^
: string
| SELECT CAST(a AS DECIMAL ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS DOUBLE ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of DOUBLE'
> exec ^
: string
| SELECT CAST(a AS DOUBLE ARRAY)
| ; SELECT FROM
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of DOUBLE'
> exec ^
: string
| SELECT CAST(a AS FLOAT ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of FLOAT'
> exec ^
: string
| SELECT CAST(a AS FLOAT ARRAY)
| ; SELECT FROM
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of FLOAT'
> exec ^
: string
| SELECT CAST(a AS REAL ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of DOUBLE'
> exec ^
: string
| SELECT CAST(a AS REAL ARRAY)
| ; SELECT FROM
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of DOUBLE'
> exec ^
: string
| SELECT CONVERT(a, SIGNED ARRAY)
< error
| 1064 42000 …near 'ARRAY)' at line 1
> exec ^
: string
| SELECT CONVERT(a, CHAR(3) ARRAY)
< error
| 1064 42000 …near 'ARRAY)' at line 1
> exec ^
: string
| SELECT CAST(a AS JSON ARRAY ARRAY)
< error
| 1064 42000 …near 'ARRAY)' at line 1
> exec ^
: string
| SELECT CAST(a AS REAL(3))
< error
| 1064 42000 …near '(3))' at line 1
> exec ^
: string
| SELECT CAST(a AS REAL(3,2))
< error
| 1064 42000 …near '(3,2))' at line 1
> exec ^
: string
| SELECT CAST(a AS DOUBLE(3))
< error
| 1064 42000 …near '(3))' at line 1
> exec ^
: string
| SELECT CAST(a AS DOUBLE PRECISION(3))
< error
| 1064 42000 …near '(3))' at line 1
> exec ^
: string
| SELECT CAST(a AS DOUBLE PRECISION(3,2))
< error
| 1064 42000 …near '(3,2))' at line 1
> exec ^
: string
| SELECT CAST(a AS FLOAT())
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| SELECT CAST(a AS FLOAT(3,))
< error
| 1064 42000 …near ',))' at line 1
> exec ^
: string
| SELECT CAST(a AS FLOAT(1.5))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS FLOAT(1.5))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS DECIMAL())
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| SELECT CAST(a AS DECIMAL(3,2,1))
< error
| 1064 42000 …near ',1))' at line 1
> exec ^
: string
| SELECT CAST(a AS DEC(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DEC(3))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS NUMERIC(3))
< error
| 1064 42000 …near 'NUMERIC(3))' at line 1
> exec ^
: string
| SELECT CAST(a AS FIXED(3))
< error
| 1064 42000 …near 'FIXED(3))' at line 1
> exec ^
: string
| SELECT CAST(a AS BIGINT(3))
< error
| 1064 42000 …near 'BIGINT(3))' at line 1
> exec ^
: string
| SELECT CAST(a AS BIGINT UNSIGNED)
< error
| 1064 42000 …near 'BIGINT UNSIGNED)' at line 1
> exec ^
: string
| SELECT CAST(a AS BIGINT SIGNED)
< error
| 1064 42000 …near 'BIGINT SIGNED)' at line 1
> exec ^
: string
| SELECT CAST(a AS INTEGER(3))
< error
| 1064 42000 …near 'INTEGER(3))' at line 1
> exec ^
: string
| SELECT CAST(a AS INTEGER UNSIGNED)
< error
| 1064 42000 …near 'INTEGER UNSIGNED)' at line 1
> exec ^
: string
| SELECT CAST(a AS INT ZEROFILL)
< error
| 1064 42000 …near 'INT ZEROFILL)' at line 1
> exec ^
: string
| SELECT CAST(a AS DECIMAL ZEROFILL)
< error
| 1064 42000 …near 'ZEROFILL)' at line 1
> exec ^
: string
| SELECT CAST(a AS DECIMAL SIGNED)
< error
| 1064 42000 …near 'SIGNED)' at line 1
> exec ^
: string
| SELECT CAST(a AS FLOAT SIGNED)
< error
| 1064 42000 …near 'SIGNED)' at line 1
> exec ^
: string
| SELECT CAST(a AS DOUBLE SIGNED)
< error
| 1064 42000 …near 'SIGNED)' at line 1
> exec ^
: string
| SELECT CAST(a AS SIGNED INT(3))
< error
| 1064 42000 …near '(3))' at line 1
> exec ^
: string
| SELECT CAST(a AS UNSIGNED BIGINT)
< error
| 1064 42000 …near 'BIGINT)' at line 1
> exec ^
: string
| SELECT CAST(a AS SIGNED BIGINT)
< error
| 1064 42000 …near 'BIGINT)' at line 1
> exec ^
: string
| SELECT CAST(a AS SIGNED UNSIGNED)
< error
| 1064 42000 …near 'UNSIGNED)' at line 1
> exec ^
: string
| SELECT CAST(a AS CHAR())
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| SELECT CAST(a AS CHAR(3,2))
< error
| 1064 42000 …near ',2))' at line 1
> exec ^
: string
| SELECT CAST(a AS CHAR(0))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR(0))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHAR(3) COLLATE utf8mb4_bin)
< error
| 1064 42000 …near 'COLLATE utf8mb4_bin)' at line 1
> exec ^
: string
| SELECT CAST(a AS CHAR CHARACTER SET utf8mb4 COLLATE utf8mb4_bin)
< error
| 1064 42000 …near 'COLLATE utf8mb4_bin)' at line 1
> exec ^
: string
| SELECT CAST(a AS CHAR(3) CHARSET utf8mb4 BINARY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR(3) CHARSET utf8mb4 BINARY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHAR CHARACTER SET utf8mb4 BINARY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR CHARACTER SET utf8mb4 BINARY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHAR BINARY CHARACTER SET utf8mb4)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR BINARY CHARACTER SET utf8mb4)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHAR BINARY BINARY)
< error
| 1064 42000 …near 'BINARY)' at line 1
> exec ^
: string
| SELECT CAST(a AS CHAR(3) ASCII)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR(3) ASCII)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHAR(3) UNICODE)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR(3) UNICODE)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHAR(3) BYTE)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR(3) BYTE)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHAR ASCII BINARY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR ASCII BINARY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHAR BINARY ASCII)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR BINARY ASCII)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHAR UNICODE BINARY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR UNICODE BINARY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHAR BINARY UNICODE)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR BINARY UNICODE)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHAR ASCII UNICODE)
< error
| 1064 42000 …near 'UNICODE)' at line 1
> exec ^
: string
| SELECT CAST(a AS CHAR ASCII CHARACTER SET utf8mb4)
< error
| 1064 42000 …near 'CHARACTER SET utf8mb4)' at line 1
> exec ^
: string
| SELECT CAST(a AS CHAR BYTE BINARY)
< error
| 1064 42000 …near 'BINARY)' at line 1
> exec ^
: string
| SELECT CAST(a AS CHAR BINARY BYTE)
< error
| 1064 42000 …near 'BYTE)' at line 1
> exec ^
: string
| SELECT CAST(a AS CHAR(3) ASCII BINARY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR(3) ASCII BINARY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHARACTER ASCII)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHARACTER ASCII)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHARACTER BINARY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHARACTER BINARY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHARACTER CHARACTER SET utf8mb4)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHARACTER CHARACTER SET utf8mb4)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS NCHAR BINARY)
< error
| 1064 42000 …near 'BINARY)' at line 1
> exec ^
: string
| SELECT CAST(a AS NCHAR CHARACTER SET utf8mb4)
< error
| 1064 42000 …near 'CHARACTER SET utf8mb4)' at line 1
> exec ^
: string
| SELECT CAST(a AS NCHAR ASCII)
< error
| 1064 42000 …near 'ASCII)' at line 1
> exec ^
: string
| SELECT CAST(a AS NATIONAL CHAR BINARY)
< error
| 1064 42000 …near 'BINARY)' at line 1
> exec ^
: string
| SELECT CAST(a AS NATIONAL CHARACTER)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS NATIONAL CHARACTER)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS NATIONAL CHARACTER(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS NATIONAL CHARACTER(3))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS NATIONAL)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT CAST(a AS VARCHAR(3) BINARY)
< error
| 1064 42000 …near 'VARCHAR(3) BINARY)' at line 1
> exec ^
: string
| SELECT CAST(a AS VARCHAR(3) COLLATE utf8mb4_bin)
< error
| 1064 42000 …near 'VARCHAR(3) COLLATE utf8mb4_bin)' at line 1
> exec ^
: string
| SELECT CAST(a AS VARCHAR(3) ASCII)
< error
| 1064 42000 …near 'VARCHAR(3) ASCII)' at line 1
> exec ^
: string
| SELECT CAST(a AS VARCHARACTER(3) CHARACTER SET utf8mb4)
< error
| 1064 42000 …near 'VARCHARACTER(3) CHARACTER SET utf8mb4)' at line 1
> exec ^
: string
| SELECT CAST(a AS VARCHAR())
< error
| 1064 42000 …near 'VARCHAR())' at line 1
> exec ^
: string
| SELECT CAST(a AS BINARY())
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| SELECT CAST(a AS BINARY(3,2))
< error
| 1064 42000 …near ',2))' at line 1
> exec ^
: string
| SELECT CAST(a AS BINARY BINARY)
< error
| 1064 42000 …near 'BINARY)' at line 1
> exec ^
: string
| SELECT CAST(a AS BINARY ASCII)
< error
| 1064 42000 …near 'ASCII)' at line 1
> exec ^
: string
| SELECT CAST(a AS BINARY COLLATE utf8mb4_bin)
< error
| 1064 42000 …near 'COLLATE utf8mb4_bin)' at line 1
> exec ^
: string
| SELECT CAST(a AS DATE BINARY)
< error
| 1064 42000 …near 'BINARY)' at line 1
> exec ^
: string
| SELECT CAST(a AS TIME())
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| SELECT CAST(a AS TIME(3,2))
< error
| 1064 42000 …near ',2))' at line 1
> exec ^
: string
| SELECT CAST(a AS TIME(1.5))
< error
| 1064 42000 …near '1.5))' at line 1
> exec ^
: string
| SELECT CAST(a AS DATETIME())
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| SELECT CAST(a AS DATETIME(3,2))
< error
| 1064 42000 …near ',2))' at line 1
> exec ^
: string
| SELECT CAST(a AS YEAR(2))
< error
| 1064 42000 …near '(2))' at line 1
> exec ^
: string
| SELECT CAST(a AS YEAR())
< error
| 1064 42000 …near '())' at line 1
> exec ^
: string
| SELECT CAST(a AS JSON BINARY)
< error
| 1064 42000 …near 'BINARY)' at line 1
> exec ^
: string
| SELECT CAST(a AS POINT(3))
< error
| 1064 42000 …near '(3))' at line 1
> exec ^
: string
| SELECT CAST(a AS MULTIPOINT)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS MULTIPOINT)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS MULTILINESTRING)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS MULTILINESTRING)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS MULTIPOLYGON)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS MULTIPOLYGON)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS GEOMETRYCOLLECTION)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS GEOMETRYCOLLECTION)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS GEOMCOLLECTION)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS GEOMCOLLECTION)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS GEOMETRY(3))
< error
| 1064 42000 …near 'GEOMETRY(3))' at line 1
> exec ^
: string
| SELECT CAST(a AS CHAR SET utf8mb4)
< error
| 1064 42000 …near 'SET utf8mb4)' at line 1
> exec ^
: string
| SELECT CAST(a AS CHAR CHAR SET utf8mb4)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR CHAR SET utf8mb4)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHAR CHARACTER SET binary)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR CHARACTER SET binary)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHAR CHARACTER SET 'utf8mb4')
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR CHARACTER SET 'utf8mb4')
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHAR CHARACTER SET)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT CAST(a AS CHAR COLLATE)
< error
| 1064 42000 …near 'COLLATE)' at line 1
> exec ^
: string
| SELECT CAST(a AS `CHAR`)
< error
| 1064 42000 …near '`CHAR`)' at line 1
> exec ^
: string
| SELECT CAST(a AS `SIGNED`)
< error
| 1064 42000 …near '`SIGNED`)' at line 1
> exec ^
: string
| SELECT CAST(a AS `DATE`)
< error
| 1064 42000 …near '`DATE`)' at line 1
> exec ^
: string
| SELECT CAST(a AS 'CHAR')
< error
| 1064 42000 …near ''CHAR')' at line 1
> exec ^
: string
| SELECT CAST(a AS 1)
< error
| 1064 42000 …near '1)' at line 1
> exec ^
: string
| SELECT CAST(a AS NULL)
< error
| 1064 42000 …near 'NULL)' at line 1
> exec ^
: string
| SELECT CONVERT(a, SIGNED INT)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CONVERT(a, SIGNED INT)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CONVERT(a, UNSIGNED INTEGER)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CONVERT(a, UNSIGNED INTEGER)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CONVERT(a, CHAR BINARY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CONVERT(a, CHAR BINARY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CONVERT(a, CHAR CHARACTER SET utf8mb4)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CONVERT(a, CHAR CHARACTER SET utf8mb4)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CONVERT(a, FLOAT(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CONVERT(a, FLOAT(3))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CONVERT(a, REAL)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CONVERT(a, REAL)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CONVERT(a, NCHAR)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CONVERT(a, NCHAR)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CONVERT(a, POINT)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CONVERT(a, POINT)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CONVERT(a, YEAR)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CONVERT(a, YEAR)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CONVERT(a, GEOMETRY)
< error
| 1064 42000 …near 'GEOMETRY)' at line 1
> exec ^
: string
| SELECT CONVERT(a, UUID)
< error
| 1064 42000 …near 'UUID)' at line 1
> exec ^
: string
| SELECT CONVERT(a USING utf8mb4)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CONVERT(a USING utf8mb4)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CONVERT(a, )
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT CONVERT(a)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT CAST(a AS)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT CAST(a)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT CAST(a AS CHAR, b)
< error
| 1064 42000 …near ', b)' at line 1
> exec ^
: string
| SELECT CAST(a AS CHAR) AS x
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR) AS x
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS NCHAR ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'specifying charset for multi-valued index'
> exec ^
: string
| SELECT CAST(a AS NCHAR ARRAY)
| ; SELECT FROM
< error
| 1235 42000 This version of MySQL doesn't yet support 'specifying charset for multi-valued index'
> exec ^
: string
| SELECT CAST(a AS POINT ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of POINT'
> exec ^
: string
| SELECT CAST(a AS POINT ARRAY)
| ; SELECT FROM
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of POINT'
> exec ^
: string
| SELECT CAST(a AS GEOMETRY ARRAY)
< error
| 1064 42000 …near 'GEOMETRY ARRAY)' at line 1
> exec ^
: string
| SELECT CAST(a AS VECTOR ARRAY)
< error
| 1064 42000 …near 'VECTOR ARRAY)' at line 1
> exec ^
: string
| SELECT CAST(a AS BIGINT UNSIGNED ARRAY)
< error
| 1064 42000 …near 'BIGINT UNSIGNED ARRAY)' at line 1
> exec ^
: string
| CREATE TABLE t9 (a int, KEY k ((CAST(a AS SIGNED ARRAY))))
< ok 0 0
> exec ^
: string
| CREATE TABLE t9 (a int, KEY k ((CAST(a AS INT ARRAY))))
< error
| 1064 42000 …near 'INT ARRAY))))' at line 1
> exec ^
: string
| CREATE TABLE t9 (a int, b int AS (CAST(a AS CHAR)))
< ok 0 0
> exec ^
: string
| CREATE TABLE t9 (a int, b int AS (CAST(a AS TEXT)))
< error
| 1064 42000 …near 'TEXT)))' at line 1
> exec ^
: string
| SELECT 1 FROM t WHERE CAST(a AS BOOLEAN)
< error
| 1064 42000 …near 'BOOLEAN)' at line 1
> exec ^
: string
| SELECT BINARY a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT BINARY a
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHAR) COLLATE utf8mb4_bin
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR) COLLATE utf8mb4_bin
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS NUMBER)
< error
| 1064 42000 …near 'NUMBER)' at line 1
> exec ^
: string
| SELECT CAST(a AS NUMBER(3,2))
< error
| 1064 42000 …near 'NUMBER(3,2))' at line 1
> exec ^
: string
| SELECT CAST(a AS ROARINGBITMAP)
< error
| 1064 42000 …near 'ROARINGBITMAP)' at line 1
> exec ^
: string
| SELECT CAST(a AS ARRAY(INT))
< error
| 1064 42000 …near 'ARRAY(INT))' at line 1
> exec ^
: string
| SELECT CAST(a AS ARRAY<INT>)
< error
| 1064 42000 …near 'ARRAY<INT>)' at line 1
> exec ^
: string
| SELECT CAST(a AS CHAR(1.5))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR(1.5))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS BINARY(1.5))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS BINARY(1.5))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS NCHAR(1.5))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS NCHAR(1.5))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS VARCHAR(1.5))
< error
| 1064 42000 …near 'VARCHAR(1.5))' at line 1
> exec ^
: string
| SELECT CAST(a AS DATETIME(1.5))
< error
| 1064 42000 …near '1.5))' at line 1
> exec ^
: string
| SELECT CAST(a AS DECIMAL(1.5))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DECIMAL(1.5))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS DATE(1.5))
< error
| 1064 42000 …near '(1.5))' at line 1
> exec ^
: string
| SELECT CAST(a AS FLOAT(1e5))
< error
| 1064 42000 …near '1e5))' at line 1
> exec ^
: string
| SELECT CAST(a AS TIME(1e5))
< error
| 1064 42000 …near '1e5))' at line 1
> exec ^
: string
| SELECT CAST(a AS VECTOR<FLOAT>(3))
< error
| 1064 42000 …near 'VECTOR<FLOAT>(3))' at line 1
> exec ^
: string
| SELECT CAST(a AS VECTOR<FLOAT>)
< error
| 1064 42000 …near 'VECTOR<FLOAT>)' at line 1
> exec ^
: string
| SELECT CAST(a AS VECTOR<DOUBLE>)
< error
| 1064 42000 …near 'VECTOR<DOUBLE>)' at line 1
> exec ^
: string
| SELECT CAST(a AS VECTOR())
< error
| 1064 42000 …near 'VECTOR())' at line 1
> exec ^
: string
| SELECT CAST(a AS VECTOR(3,2))
< error
| 1064 42000 …near 'VECTOR(3,2))' at line 1
> exec ^
: string
| SELECT CAST(a AS SIGNED INT INTEGER)
< error
| 1064 42000 …near 'INTEGER)' at line 1
> exec ^
: string
| SELECT CAST(a AS SIGNED INTEGER INT)
< error
| 1064 42000 …near 'INT)' at line 1
> exec ^
: string
| SELECT CAST(a AS UNSIGNED INT INT)
< error
| 1064 42000 …near 'INT)' at line 1
> exec ^
: string
| SELECT CAST(a AS NATIONAL VARCHAR)
< error
| 1064 42000 …near 'VARCHAR)' at line 1
> exec ^
: string
| SELECT CAST(a AS NCHAR VARYING(3))
< error
| 1064 42000 …near 'VARYING(3))' at line 1
> exec ^
: string
| SELECT CAST(a AS CHAR(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin BINARY)
< error
| 1064 42000 …near 'COLLATE utf8mb4_bin BINARY)' at line 1
> exec ^
: string
| SELECT CAST(a AS CHAR COLLATE utf8mb4_bin BINARY)
< error
| 1064 42000 …near 'COLLATE utf8mb4_bin BINARY)' at line 1
> exec ^
: string
| SELECT CAST(a AS CHAR BINARY COLLATE utf8mb4_bin)
< error
| 1064 42000 …near 'COLLATE utf8mb4_bin)' at line 1
> exec ^
: string
| SELECT CAST(a AS CHAR ASCII COLLATE utf8mb4_bin)
< error
| 1064 42000 …near 'COLLATE utf8mb4_bin)' at line 1
> exec ^
: string
| SELECT CAST(a AS BOOL ARRAY)
< error
| 1064 42000 …near 'BOOL ARRAY)' at line 1
> exec ^
: string
| SELECT CAST(a AS FLOAT4)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS FLOAT4)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS FLOAT4(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS FLOAT4(3))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS FLOAT8)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS FLOAT8)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CONVERT(a, FLOAT8)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CONVERT(a, FLOAT8)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS FLOAT8 PRECISION)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS FLOAT8 PRECISION)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS FLOAT8(10,2))
< error
| 1064 42000 …near '(10,2))' at line 1
> exec ^
: string
| SELECT CAST(a AS INT4)
< error
| 1064 42000 …near 'INT4)' at line 1
> exec ^
: string
| SELECT CAST(a AS INT8)
< error
| 1064 42000 …near 'INT8)' at line 1
> exec ^
: string
| SELECT CAST(a AS INT1)
< error
| 1064 42000 …near 'INT1)' at line 1
> exec ^
: string
| SELECT CAST(a AS SIGNED INT4)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS SIGNED INT4)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS UNSIGNED INT4)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS UNSIGNED INT4)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS SIGNED INT8)
< error
| 1064 42000 …near 'INT8)' at line 1
