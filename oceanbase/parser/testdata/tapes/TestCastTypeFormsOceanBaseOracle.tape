-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestCastTypeFormsOceanBaseOracle default#1
> exec
| DROP DATABASE IF EXISTS obcast
< ok 0 0
> exec
| CREATE DATABASE obcast
< ok 1 0
> exec
| USE obcast
< ok 0 0
> exec
| SELECT CAST(a AS INT)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS INTEGER)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS TINYINT)
< error
| 1064 42000 …near 'TINYINT)' at line 1
> exec
| SELECT CAST(a AS SMALLINT)
< error
| 1064 42000 …near 'SMALLINT)' at line 1
> exec
| SELECT CAST(a AS MEDIUMINT)
< error
| 1064 42000 …near 'MEDIUMINT)' at line 1
> exec
| SELECT CAST(a AS BIGINT)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS INT(3))
< error
| 1064 42000 …near '(3))' at line 1
> exec
| SELECT CAST(a AS SIGNED)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS SIGNED INT)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS SIGNED INTEGER)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS UNSIGNED)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS UNSIGNED INT)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS UNSIGNED INTEGER)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS SIGNED(3))
< error
| 1064 42000 …near '(3))' at line 1
> exec
| SELECT CAST(a AS INT UNSIGNED)
< error
| 1064 42000 …near 'UNSIGNED)' at line 1
> exec
| SELECT CAST(a AS INT SIGNED)
< error
| 1064 42000 …near 'SIGNED)' at line 1
> exec
| SELECT CAST(a AS FLOAT)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS FLOAT(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS FLOAT(3,2))
< error
| 1064 42000 …near ',2))' at line 1
> exec
| SELECT CAST(a AS DOUBLE)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS DOUBLE(3,2))
< error
| 1064 42000 …near '(3,2))' at line 1
> exec
| SELECT CAST(a AS DOUBLE PRECISION)
< error
| 1064 42000 …near 'PRECISION)' at line 1
> exec
| SELECT CAST(a AS REAL)
< error
| 1064 42000 …near 'REAL)' at line 1
> exec
| SELECT CAST(a AS DECIMAL)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS DECIMAL(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS DECIMAL(3,2))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS NUMERIC)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS NUMERIC(3,2))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS DEC)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS DEC(3,2))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS FIXED)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS FIXED(3,2))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS DECIMAL UNSIGNED)
< error
| 1064 42000 …near 'UNSIGNED)' at line 1
> exec
| SELECT CAST(a AS NUMBER)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS NUMBER(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS NUMBER(3,2))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS NUMBER UNSIGNED)
< error
| 1064 42000 …near 'UNSIGNED)' at line 1
> exec
| SELECT CAST(a AS DOUBLE UNSIGNED)
< error
| 1064 42000 …near 'UNSIGNED)' at line 1
> exec
| SELECT CAST(a AS FLOAT UNSIGNED)
< error
| 1064 42000 …near 'UNSIGNED)' at line 1
> exec
| SELECT CAST(a AS BOOL)
< error
| 1064 42000 …near 'BOOL)' at line 1
> exec
| SELECT CAST(a AS BOOLEAN)
< error
| 1064 42000 …near 'BOOLEAN)' at line 1
> exec
| SELECT CAST(a AS CHAR)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS CHAR(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS CHARACTER)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS CHARACTER(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS NCHAR)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS NCHAR(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS NATIONAL CHAR)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS NATIONAL CHAR(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS CHAR VARYING(3))
< error
| 1064 42000 …near 'VARYING(3))' at line 1
> exec
| SELECT CAST(a AS CHARACTER VARYING(3))
< error
| 1064 42000 …near 'VARYING(3))' at line 1
> exec
| SELECT CAST(a AS VARCHAR(3))
< error
| 1064 42000 …near 'VARCHAR(3))' at line 1
> exec
| SELECT CAST(a AS VARCHAR)
< error
| 1064 42000 …near 'VARCHAR)' at line 1
> exec
| SELECT CAST(a AS VARCHARACTER(3))
< error
| 1064 42000 …near 'VARCHARACTER(3))' at line 1
> exec
| SELECT CAST(a AS NVARCHAR(3))
< error
| 1064 42000 …near 'NVARCHAR(3))' at line 1
> exec
| SELECT CAST(a AS NATIONAL VARCHAR(3))
< error
| 1064 42000 …near 'VARCHAR(3))' at line 1
> exec
| SELECT CAST(a AS NCHAR VARCHAR(3))
< error
| 1064 42000 …near 'VARCHAR(3))' at line 1
> exec
| SELECT CAST(a AS TEXT)
< error
| 1064 42000 …near 'TEXT)' at line 1
> exec
| SELECT CAST(a AS TEXT(3))
< error
| 1064 42000 …near 'TEXT(3))' at line 1
> exec
| SELECT CAST(a AS TINYTEXT)
< error
| 1064 42000 …near 'TINYTEXT)' at line 1
> exec
| SELECT CAST(a AS MEDIUMTEXT)
< error
| 1064 42000 …near 'MEDIUMTEXT)' at line 1
> exec
| SELECT CAST(a AS LONGTEXT)
< error
| 1064 42000 …near 'LONGTEXT)' at line 1
> exec
| SELECT CAST(a AS LONG)
< error
| 1064 42000 …near 'LONG)' at line 1
> exec
| SELECT CAST(a AS LONG VARCHAR)
< error
| 1064 42000 …near 'LONG VARCHAR)' at line 1
> exec
| SELECT CAST(a AS BLOB)
< error
| 1064 42000 …near 'BLOB)' at line 1
> exec
| SELECT CAST(a AS BLOB(3))
< error
| 1064 42000 …near 'BLOB(3))' at line 1
> exec
| SELECT CAST(a AS TINYBLOB)
< error
| 1064 42000 …near 'TINYBLOB)' at line 1
> exec
| SELECT CAST(a AS MEDIUMBLOB)
< error
| 1064 42000 …near 'MEDIUMBLOB)' at line 1
> exec
| SELECT CAST(a AS LONGBLOB)
< error
| 1064 42000 …near 'LONGBLOB)' at line 1
> exec
| SELECT CAST(a AS BINARY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS BINARY(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS VARBINARY(3))
< error
| 1064 42000 …near 'VARBINARY(3))' at line 1
> exec
| SELECT CAST(a AS VARBINARY)
< error
| 1064 42000 …near 'VARBINARY)' at line 1
> exec
| SELECT CAST(a AS BIT)
< error
| 1064 42000 …near 'BIT)' at line 1
> exec
| SELECT CAST(a AS BIT(1))
< error
| 1064 42000 …near 'BIT(1))' at line 1
> exec
| SELECT CAST(a AS DATE)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS DATE(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS TIME)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS TIME(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS DATETIME)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS DATETIME(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS TIMESTAMP)
< error
| 1064 42000 …near 'TIMESTAMP)' at line 1
> exec
| SELECT CAST(a AS TIMESTAMP(3))
< error
| 1064 42000 …near 'TIMESTAMP(3))' at line 1
> exec
| SELECT CAST(a AS YEAR)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS YEAR(4))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS JSON)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS JSON(3))
< error
| 1064 42000 …near '(3))' at line 1
> exec
| SELECT CAST(a AS ENUM('a'))
< error
| 1064 42000 …near 'ENUM('a'))' at line 1
> exec
| SELECT CAST(a AS SET('a'))
< error
| 1064 42000 …near 'SET('a'))' at line 1
> exec
| SELECT CAST(a AS GEOMETRY)
< error
| 1064 42000 …near 'GEOMETRY)' at line 1
> exec
| SELECT CAST(a AS POINT)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS LINESTRING)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS POLYGON)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS UUID)
< error
| 1064 42000 …near 'UUID)' at line 1
> exec
| SELECT CAST(a AS INET4)
< error
| 1064 42000 …near 'INET4)' at line 1
> exec
| SELECT CAST(a AS INET6)
< error
| 1064 42000 …near 'INET6)' at line 1
> exec
| SELECT CAST(a AS SERIAL)
< error
| 1064 42000 …near 'SERIAL)' at line 1
> exec
| SELECT CAST(a AS CHAR CHARACTER SET utf8mb4)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS CHAR CHARSET utf8mb4)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS CHAR(3) CHARACTER SET utf8mb4)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS CHAR ASCII)
< error
| 1064 42000 …near 'ASCII)' at line 1
> exec
| SELECT CAST(a AS CHAR UNICODE)
< error
| 1064 42000 …near 'UNICODE)' at line 1
> exec
| SELECT CAST(a AS CHAR BYTE)
< error
| 1064 42000 …near 'BYTE)' at line 1
> exec
| SELECT CAST(a AS CHAR BINARY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS CHAR(3) BINARY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS BINARY CHARACTER SET utf8mb4)
< error
| 1064 42000 …near 'CHARACTER SET utf8mb4)' at line 1
> exec
| SELECT CAST(a AS CHAR COLLATE utf8mb4_bin)
< error
| 1064 42000 …near 'COLLATE utf8mb4_bin)' at line 1
> exec
| SELECT CAST(a AS VARCHAR(3) CHARACTER SET utf8mb4)
< error
| 1064 42000 …near 'VARCHAR(3) CHARACTER SET utf8mb4)' at line 1
> exec
| SELECT CAST(a AS INTERVAL DAY)
< error
| 1064 42000 …near 'INTERVAL DAY)' at line 1
> exec
| SELECT CAST(a AS DAY)
< error
| 1064 42000 …near 'DAY)' at line 1
> exec
| SELECT CAST(a AS VECTOR)
< error
| 1064 42000 …near 'VECTOR)' at line 1
> exec
| SELECT CAST(a AS VECTOR(3))
< error
| 1064 42000 …near 'VECTOR(3))' at line 1
> exec
| SELECT CONVERT(a, INT)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CONVERT(a, SIGNED)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CONVERT(a, CHAR(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CONVERT(a, DATE)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CONVERT(a, DECIMAL(3,2))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CONVERT(a, BIGINT)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CONVERT(a, TEXT)
< error
| 1064 42000 …near 'TEXT)' at line 1
> exec
| SELECT CONVERT(a, VARCHAR(3))
< error
| 1064 42000 …near 'VARCHAR(3))' at line 1
> exec
| SELECT CONVERT(a, BOOLEAN)
< error
| 1064 42000 …near 'BOOLEAN)' at line 1
> exec
| SELECT CONVERT(a, JSON)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CONVERT(a, DOUBLE)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CONVERT(a, FLOAT)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CONVERT(a, TIME(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS INT ARRAY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS BIGINT ARRAY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS SIGNED ARRAY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS UNSIGNED ARRAY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS SIGNED INTEGER ARRAY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS CHAR ARRAY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS CHAR(3) ARRAY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS BINARY ARRAY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS BINARY(3) ARRAY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS DATE ARRAY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS DATETIME ARRAY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS TIME ARRAY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS YEAR ARRAY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS JSON ARRAY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS DECIMAL ARRAY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS DOUBLE ARRAY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS FLOAT ARRAY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS REAL ARRAY)
< error
| 1064 42000 …near 'REAL ARRAY)' at line 1
> exec
| SELECT CONVERT(a, SIGNED ARRAY)
< error
| 1064 42000 …near 'ARRAY)' at line 1
> exec
| SELECT CONVERT(a, CHAR(3) ARRAY)
< error
| 1064 42000 …near 'ARRAY)' at line 1
> exec
| SELECT CAST(a AS JSON ARRAY ARRAY)
< error
| 1064 42000 …near 'ARRAY)' at line 1
> exec
| SELECT CAST(a AS REAL(3))
< error
| 1064 42000 …near 'REAL(3))' at line 1
> exec
| SELECT CAST(a AS REAL(3,2))
< error
| 1064 42000 …near 'REAL(3,2))' at line 1
> exec
| SELECT CAST(a AS DOUBLE(3))
< error
| 1064 42000 …near '(3))' at line 1
> exec
| SELECT CAST(a AS DOUBLE PRECISION(3))
< error
| 1064 42000 …near 'PRECISION(3))' at line 1
> exec
| SELECT CAST(a AS DOUBLE PRECISION(3,2))
< error
| 1064 42000 …near 'PRECISION(3,2))' at line 1
> exec
| SELECT CAST(a AS FLOAT())
< error
| 1064 42000 …near '))' at line 1
> exec
| SELECT CAST(a AS FLOAT(3,))
< error
| 1064 42000 …near ',))' at line 1
> exec
| SELECT CAST(a AS FLOAT(1.5))
< error
| 1064 42000 …near '1.5))' at line 1
> exec
| SELECT CAST(a AS DECIMAL())
< error
| 1064 42000 …near '))' at line 1
> exec
| SELECT CAST(a AS DECIMAL(3,2,1))
< error
| 1064 42000 …near ',1))' at line 1
> exec
| SELECT CAST(a AS DEC(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS NUMERIC(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS FIXED(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS BIGINT(3))
< error
| 1064 42000 …near '(3))' at line 1
> exec
| SELECT CAST(a AS BIGINT UNSIGNED)
< error
| 1064 42000 …near 'UNSIGNED)' at line 1
> exec
| SELECT CAST(a AS BIGINT SIGNED)
< error
| 1064 42000 …near 'SIGNED)' at line 1
> exec
| SELECT CAST(a AS INTEGER(3))
< error
| 1064 42000 …near '(3))' at line 1
> exec
| SELECT CAST(a AS INTEGER UNSIGNED)
< error
| 1064 42000 …near 'UNSIGNED)' at line 1
> exec
| SELECT CAST(a AS INT ZEROFILL)
< error
| 1064 42000 …near 'ZEROFILL)' at line 1
> exec
| SELECT CAST(a AS DECIMAL ZEROFILL)
< error
| 1064 42000 …near 'ZEROFILL)' at line 1
> exec
| SELECT CAST(a AS DECIMAL SIGNED)
< error
| 1064 42000 …near 'SIGNED)' at line 1
> exec
| SELECT CAST(a AS FLOAT SIGNED)
< error
| 1064 42000 …near 'SIGNED)' at line 1
> exec
| SELECT CAST(a AS DOUBLE SIGNED)
< error
| 1064 42000 …near 'SIGNED)' at line 1
> exec
| SELECT CAST(a AS SIGNED INT(3))
< error
| 1064 42000 …near '(3))' at line 1
> exec
| SELECT CAST(a AS UNSIGNED BIGINT)
< error
| 1064 42000 …near 'BIGINT)' at line 1
> exec
| SELECT CAST(a AS SIGNED BIGINT)
< error
| 1064 42000 …near 'BIGINT)' at line 1
> exec
| SELECT CAST(a AS SIGNED UNSIGNED)
< error
| 1064 42000 …near 'UNSIGNED)' at line 1
> exec
| SELECT CAST(a AS CHAR())
< error
| 1064 42000 …near '))' at line 1
> exec
| SELECT CAST(a AS CHAR(3,2))
< error
| 1064 42000 …near ',2))' at line 1
> exec
| SELECT CAST(a AS CHAR(0))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS CHAR(3) COLLATE utf8mb4_bin)
< error
| 1064 42000 …near 'COLLATE utf8mb4_bin)' at line 1
> exec
| SELECT CAST(a AS CHAR CHARACTER SET utf8mb4 COLLATE utf8mb4_bin)
< error
| 1064 42000 …near 'COLLATE utf8mb4_bin)' at line 1
> exec
| SELECT CAST(a AS CHAR(3) CHARSET utf8mb4 BINARY)
< error
| 1064 42000 …near 'BINARY)' at line 1
> exec
| SELECT CAST(a AS CHAR CHARACTER SET utf8mb4 BINARY)
< error
| 1064 42000 …near 'BINARY)' at line 1
> exec
| SELECT CAST(a AS CHAR BINARY CHARACTER SET utf8mb4)
< error
| 1064 42000 …near 'CHARACTER SET utf8mb4)' at line 1
> exec
| SELECT CAST(a AS CHAR BINARY BINARY)
< error
| 1064 42000 …near 'BINARY)' at line 1
> exec
| SELECT CAST(a AS CHAR(3) ASCII)
< error
| 1064 42000 …near 'ASCII)' at line 1
> exec
| SELECT CAST(a AS CHAR(3) UNICODE)
< error
| 1064 42000 …near 'UNICODE)' at line 1
> exec
| SELECT CAST(a AS CHAR(3) BYTE)
< error
| 1064 42000 …near 'BYTE)' at line 1
> exec
| SELECT CAST(a AS CHAR ASCII BINARY)
< error
| 1064 42000 …near 'ASCII BINARY)' at line 1
> exec
| SELECT CAST(a AS CHAR BINARY ASCII)
< error
| 1064 42000 …near 'ASCII)' at line 1
> exec
| SELECT CAST(a AS CHAR UNICODE BINARY)
< error
| 1064 42000 …near 'UNICODE BINARY)' at line 1
> exec
| SELECT CAST(a AS CHAR BINARY UNICODE)
< error
| 1064 42000 …near 'UNICODE)' at line 1
> exec
| SELECT CAST(a AS CHAR ASCII UNICODE)
< error
| 1064 42000 …near 'ASCII UNICODE)' at line 1
> exec
| SELECT CAST(a AS CHAR ASCII CHARACTER SET utf8mb4)
< error
| 1064 42000 …near 'ASCII CHARACTER SET utf8mb4)' at line 1
> exec
| SELECT CAST(a AS CHAR BYTE BINARY)
< error
| 1064 42000 …near 'BYTE BINARY)' at line 1
> exec
| SELECT CAST(a AS CHAR(3) ASCII BINARY)
< error
| 1064 42000 …near 'ASCII BINARY)' at line 1
> exec
| SELECT CAST(a AS CHARACTER ASCII)
< error
| 1064 42000 …near 'ASCII)' at line 1
> exec
| SELECT CAST(a AS CHARACTER BINARY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS CHARACTER CHARACTER SET utf8mb4)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS NCHAR BINARY)
< error
| 1064 42000 …near 'BINARY)' at line 1
> exec
| SELECT CAST(a AS NCHAR CHARACTER SET utf8mb4)
< error
| 1064 42000 …near 'CHARACTER SET utf8mb4)' at line 1
> exec
| SELECT CAST(a AS NCHAR ASCII)
< error
| 1064 42000 …near 'ASCII)' at line 1
> exec
| SELECT CAST(a AS NATIONAL CHAR BINARY)
< error
| 1064 42000 …near 'BINARY)' at line 1
> exec
| SELECT CAST(a AS NATIONAL CHARACTER)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS NATIONAL CHARACTER(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS NATIONAL)
< error
| 1064 42000 …near ')' at line 1
> exec
| SELECT CAST(a AS VARCHAR(3) BINARY)
< error
| 1064 42000 …near 'VARCHAR(3) BINARY)' at line 1
> exec
| SELECT CAST(a AS VARCHAR(3) COLLATE utf8mb4_bin)
< error
| 1064 42000 …near 'VARCHAR(3) COLLATE utf8mb4_bin)' at line 1
> exec
| SELECT CAST(a AS VARCHAR(3) ASCII)
< error
| 1064 42000 …near 'VARCHAR(3) ASCII)' at line 1
> exec
| SELECT CAST(a AS VARCHARACTER(3) CHARACTER SET utf8mb4)
< error
| 1064 42000 …near 'VARCHARACTER(3) CHARACTER SET utf8mb4)' at line 1
> exec
| SELECT CAST(a AS VARCHAR())
< error
| 1064 42000 …near 'VARCHAR())' at line 1
> exec
| SELECT CAST(a AS BINARY())
< error
| 1064 42000 …near '))' at line 1
> exec
| SELECT CAST(a AS BINARY(3,2))
< error
| 1064 42000 …near ',2))' at line 1
> exec
| SELECT CAST(a AS BINARY BINARY)
< error
| 1064 42000 …near 'BINARY)' at line 1
> exec
| SELECT CAST(a AS BINARY ASCII)
< error
| 1064 42000 …near 'ASCII)' at line 1
> exec
| SELECT CAST(a AS BINARY COLLATE utf8mb4_bin)
< error
| 1064 42000 …near 'COLLATE utf8mb4_bin)' at line 1
> exec
| SELECT CAST(a AS DATE BINARY)
< error
| 1064 42000 …near 'BINARY)' at line 1
> exec
| SELECT CAST(a AS TIME())
< error
| 1064 42000 …near '))' at line 1
> exec
| SELECT CAST(a AS TIME(3,2))
< error
| 1064 42000 …near ',2))' at line 1
> exec
| SELECT CAST(a AS TIME(1.5))
< error
| 1064 42000 …near '1.5))' at line 1
> exec
| SELECT CAST(a AS DATETIME())
< error
| 1064 42000 …near '))' at line 1
> exec
| SELECT CAST(a AS DATETIME(3,2))
< error
| 1064 42000 …near ',2))' at line 1
> exec
| SELECT CAST(a AS YEAR(2))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS YEAR())
< error
| 1064 42000 …near '))' at line 1
> exec
| SELECT CAST(a AS JSON BINARY)
< error
| 1064 42000 …near 'BINARY)' at line 1
> exec
| SELECT CAST(a AS POINT(3))
< error
| 1064 42000 …near '(3))' at line 1
> exec
| SELECT CAST(a AS MULTIPOINT)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS MULTILINESTRING)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS MULTIPOLYGON)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS GEOMETRYCOLLECTION)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS GEOMCOLLECTION)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS GEOMETRY(3))
< error
| 1064 42000 …near 'GEOMETRY(3))' at line 1
> exec
| SELECT CAST(a AS CHAR SET utf8mb4)
< error
| 1064 42000 …near 'SET utf8mb4)' at line 1
> exec
| SELECT CAST(a AS CHAR CHAR SET utf8mb4)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS CHAR CHARACTER SET binary)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS CHAR CHARACTER SET 'utf8mb4')
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS CHAR CHARACTER SET)
< error
| 1064 42000 …near ')' at line 1
> exec
| SELECT CAST(a AS CHAR COLLATE)
< error
| 1064 42000 …near 'COLLATE)' at line 1
> exec
| SELECT CAST(a AS `CHAR`)
< error
| 1064 42000 …near '`CHAR`)' at line 1
> exec
| SELECT CAST(a AS `SIGNED`)
< error
| 1064 42000 …near '`SIGNED`)' at line 1
> exec
| SELECT CAST(a AS `DATE`)
< error
| 1064 42000 …near '`DATE`)' at line 1
> exec
| SELECT CAST(a AS 'CHAR')
< error
| 1064 42000 …near ''CHAR')' at line 1
> exec
| SELECT CAST(a AS 1)
< error
| 1064 42000 …near '1)' at line 1
> exec
| SELECT CAST(a AS NULL)
< error
| 1064 42000 …near 'NULL)' at line 1
> exec
| SELECT CONVERT(a, SIGNED INT)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CONVERT(a, UNSIGNED INTEGER)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CONVERT(a, CHAR BINARY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CONVERT(a, CHAR CHARACTER SET utf8mb4)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CONVERT(a, FLOAT(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CONVERT(a, REAL)
< error
| 1064 42000 …near 'REAL)' at line 1
> exec
| SELECT CONVERT(a, NCHAR)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CONVERT(a, POINT)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CONVERT(a, YEAR)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CONVERT(a, GEOMETRY)
< error
| 1064 42000 …near 'GEOMETRY)' at line 1
> exec
| SELECT CONVERT(a, UUID)
< error
| 1064 42000 …near 'UUID)' at line 1
> exec
| SELECT CONVERT(a USING utf8mb4)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CONVERT(a, )
< error
| 1064 42000 …near ')' at line 1
> exec
| SELECT CONVERT(a)
< error
| 1064 42000 …near ')' at line 1
> exec
| SELECT CAST(a AS)
< error
| 1064 42000 …near ')' at line 1
> exec
| SELECT CAST(a)
< error
| 1064 42000 …near ')' at line 1
> exec
| SELECT CAST(a AS CHAR, b)
< error
| 1064 42000 …near ', b)' at line 1
> exec
| SELECT CAST(a AS CHAR) AS x
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS NCHAR ARRAY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS POINT ARRAY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS GEOMETRY ARRAY)
< error
| 1064 42000 …near 'GEOMETRY ARRAY)' at line 1
> exec
| SELECT CAST(a AS VECTOR ARRAY)
< error
| 1064 42000 …near 'VECTOR ARRAY)' at line 1
> exec
| SELECT CAST(a AS BIGINT UNSIGNED ARRAY)
< error
| 1064 42000 …near 'UNSIGNED ARRAY)' at line 1
> exec
| CREATE TABLE t9 (a int, KEY k ((CAST(a AS SIGNED ARRAY))))
< error
| 3146 22032 Invalid data type for JSON data in argument 1 to function json_value; a JSON string or JSON type is required.
> exec
| CREATE TABLE t9 (a int, KEY k ((CAST(a AS INT ARRAY))))
< error
| 3146 22032 Invalid data type for JSON data in argument 1 to function json_value; a JSON string or JSON type is required.
> exec
| CREATE TABLE t9 (a int, b int AS (CAST(a AS CHAR)))
< ok 0 0
> exec
| CREATE TABLE t9 (a int, b int AS (CAST(a AS TEXT)))
< error
| 1064 42000 …near 'TEXT)))' at line 1
> exec
| SELECT 1 FROM t WHERE CAST(a AS BOOLEAN)
< error
| 1064 42000 …near 'BOOLEAN)' at line 1
> exec
| SELECT BINARY a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS CHAR) COLLATE utf8mb4_bin
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS NUMBER)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS NUMBER(3,2))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS ROARINGBITMAP)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS ARRAY(INT))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS ARRAY<INT>)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS ARRAY<ARRAY<INT>>)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS ARRAY(ARRAY(INT)))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS ARRAY(VARCHAR(3) COLLATE utf8mb4_bin))
< error
| 1064 42000 …near 'COLLATE utf8mb4_bin))' at line 1
> exec
| SELECT CAST(a AS ARRAY(VARCHAR(3)))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS ARRAY<ARRAY(INT)>)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS ARRAY(ARRAY<INT>))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS ARRAY<ARRAY<ARRAY<INT>>>)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS ARRAY<ARRAY<INT>)
< error
| 1064 42000 …near ')' at line 1
> exec
| SELECT CAST(a AS ARRAY<INT>>)
< error
| 1064 42000 …near '>>)' at line 1
> exec
| SELECT CAST(a AS ARRAY(ARRAY<INT))
< error
| 1064 42000 …near '))' at line 1
> exec
| SELECT CAST(a AS ARRAY())
< error
| 1064 42000 …near '))' at line 1
> exec
| SELECT CAST(a AS ARRAY)
< error
| 1064 42000 …near ')' at line 1
> exec
| SELECT CAST(a AS ARRAY<>)
< error
| 1064 42000 …near '<>)' at line 1
> exec
| SELECT CAST(a AS ARRAY(INT)
< error
| 1064 42000 …near ')' at line 1
> exec
| SELECT CAST(a AS ARRAY<INT)
< error
| 1064 42000 …near ')' at line 1
> exec
| SELECT CAST(a AS CHAR(1.5))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS BINARY(1.5))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS NCHAR(1.5))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS VARCHAR(1.5))
< error
| 1064 42000 …near 'VARCHAR(1.5))' at line 1
> exec
| SELECT CAST(a AS DATETIME(1.5))
< error
| 1064 42000 …near '1.5))' at line 1
> exec
| SELECT CAST(a AS DECIMAL(1.5))
< error
| 1064 42000 …near '1.5))' at line 1
> exec
| SELECT CAST(a AS DATE(1.5))
< error
| 1064 42000 …near '1.5))' at line 1
> exec
| SELECT CAST(a AS FLOAT(1e5))
< error
| 1064 42000 …near '1e5))' at line 1
> exec
| SELECT CAST(a AS TIME(1e5))
< error
| 1064 42000 …near '1e5))' at line 1
> exec
| SELECT CAST(a AS VECTOR<FLOAT>(3))
< error
| 1064 42000 …near 'VECTOR<FLOAT>(3))' at line 1
> exec
| SELECT CAST(a AS VECTOR<FLOAT>)
< error
| 1064 42000 …near 'VECTOR<FLOAT>)' at line 1
> exec
| SELECT CAST(a AS VECTOR<DOUBLE>)
< error
| 1064 42000 …near 'VECTOR<DOUBLE>)' at line 1
> exec
| SELECT CAST(a AS VECTOR())
< error
| 1064 42000 …near 'VECTOR())' at line 1
> exec
| SELECT CAST(a AS VECTOR(3,2))
< error
| 1064 42000 …near 'VECTOR(3,2))' at line 1
> exec
| SELECT CAST(a AS SIGNED INT INTEGER)
< error
| 1064 42000 …near 'INTEGER)' at line 1
> exec
| SELECT CAST(a AS SIGNED INTEGER INT)
< error
| 1064 42000 …near 'INT)' at line 1
> exec
| SELECT CAST(a AS UNSIGNED INT INT)
< error
| 1064 42000 …near 'INT)' at line 1
> exec
| SELECT CAST(a AS NATIONAL VARCHAR)
< error
| 1064 42000 …near 'VARCHAR)' at line 1
> exec
| SELECT CAST(a AS NCHAR VARYING(3))
< error
| 1064 42000 …near 'VARYING(3))' at line 1
> exec
| SELECT CAST(a AS CHAR(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin BINARY)
< error
| 1064 42000 …near 'COLLATE utf8mb4_bin BINARY)' at line 1
> exec
| SELECT CAST(a AS CHAR COLLATE utf8mb4_bin BINARY)
< error
| 1064 42000 …near 'COLLATE utf8mb4_bin BINARY)' at line 1
> exec
| SELECT CAST(a AS CHAR BINARY COLLATE utf8mb4_bin)
< error
| 1064 42000 …near 'COLLATE utf8mb4_bin)' at line 1
> exec
| SELECT CAST(a AS CHAR ASCII COLLATE utf8mb4_bin)
< error
| 1064 42000 …near 'ASCII COLLATE utf8mb4_bin)' at line 1
> exec
| SELECT CAST(a AS BOOL ARRAY)
< error
| 1064 42000 …near 'BOOL ARRAY)' at line 1
> exec
| SELECT CAST(a AS FLOAT4)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS FLOAT4(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS FLOAT8)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CONVERT(a, FLOAT8)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS FLOAT8 PRECISION)
< error
| 1064 42000 …near 'PRECISION)' at line 1
> exec
| SELECT CAST(a AS INT4)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS INT8)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS INT1)
< error
| 1064 42000 …near 'INT1)' at line 1
> exec
| SELECT CAST(a AS SIGNED INT4)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SELECT CAST(a AS UNSIGNED INT4)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS obcast
< ok 1 0
