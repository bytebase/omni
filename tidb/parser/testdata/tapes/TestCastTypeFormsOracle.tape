-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestCastTypeFormsOracle default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| SELECT CAST(a AS FLOAT(53))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS FLOAT(53))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
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
| SELECT CAST(a AS FLOAT(54,2))
< error
| 1426 42000 Too-big precision 54 specified for 'CAST'. Maximum is 53.
> exec ^
: string
| SELECT CAST(a AS FLOAT(54,2))
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
| SELECT CAST(a AS INT)
< error
| 1064 42000 …line 1 column 20 near "INT)" 
> exec ^
: string
| SELECT CAST(a AS INTEGER)
< error
| 1064 42000 …line 1 column 24 near "INTEGER)" 
> exec ^
: string
| SELECT CAST(a AS TINYINT)
< error
| 1064 42000 …line 1 column 24 near "TINYINT)" 
> exec ^
: string
| SELECT CAST(a AS SMALLINT)
< error
| 1064 42000 …line 1 column 25 near "SMALLINT)" 
> exec ^
: string
| SELECT CAST(a AS MEDIUMINT)
< error
| 1064 42000 …line 1 column 26 near "MEDIUMINT)" 
> exec ^
: string
| SELECT CAST(a AS BIGINT)
< error
| 1064 42000 …line 1 column 23 near "BIGINT)" 
> exec ^
: string
| SELECT CAST(a AS INT(3))
< error
| 1064 42000 …line 1 column 20 near "INT(3))" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS SIGNED(3))
< error
| 1064 42000 …line 1 column 24 near "(3))" 
> exec ^
: string
| SELECT CAST(a AS INT UNSIGNED)
< error
| 1064 42000 …line 1 column 20 near "INT UNSIGNED)" 
> exec ^
: string
| SELECT CAST(a AS INT SIGNED)
< error
| 1064 42000 …line 1 column 20 near "INT SIGNED)" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS FLOAT(3,2))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS FLOAT(3,2))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS DOUBLE(3,2))
< error
| 1064 42000 …line 1 column 24 near "(3,2))" 
> exec ^
: string
| SELECT CAST(a AS DOUBLE PRECISION)
< error
| 1064 42000 …line 1 column 33 near "PRECISION)" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS NUMERIC)
< error
| 1064 42000 …line 1 column 24 near "NUMERIC)" 
> exec ^
: string
| SELECT CAST(a AS NUMERIC(3,2))
< error
| 1064 42000 …line 1 column 24 near "NUMERIC(3,2))" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS FIXED)
< error
| 1064 42000 …line 1 column 22 near "FIXED)" 
> exec ^
: string
| SELECT CAST(a AS FIXED(3,2))
< error
| 1064 42000 …line 1 column 22 near "FIXED(3,2))" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL UNSIGNED)
< error
| 1064 42000 …line 1 column 33 near "UNSIGNED)" 
> exec ^
: string
| SELECT CAST(a AS DOUBLE UNSIGNED)
< error
| 1064 42000 …line 1 column 32 near "UNSIGNED)" 
> exec ^
: string
| SELECT CAST(a AS FLOAT UNSIGNED)
< error
| 1064 42000 …line 1 column 31 near "UNSIGNED)" 
> exec ^
: string
| SELECT CAST(a AS BOOL)
< error
| 1064 42000 …line 1 column 21 near "BOOL)" 
> exec ^
: string
| SELECT CAST(a AS BOOLEAN)
< error
| 1064 42000 …line 1 column 24 near "BOOLEAN)" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS NCHAR)
< error
| 1064 42000 …line 1 column 22 near "NCHAR)" 
> exec ^
: string
| SELECT CAST(a AS NCHAR(3))
< error
| 1064 42000 …line 1 column 22 near "NCHAR(3))" 
> exec ^
: string
| SELECT CAST(a AS NATIONAL CHAR)
< error
| 1064 42000 …line 1 column 25 near "NATIONAL CHAR)" 
> exec ^
: string
| SELECT CAST(a AS NATIONAL CHAR(3))
< error
| 1064 42000 …line 1 column 25 near "NATIONAL CHAR(3))" 
> exec ^
: string
| SELECT CAST(a AS CHAR VARYING(3))
< error
| 1064 42000 …line 1 column 29 near "VARYING(3))" 
> exec ^
: string
| SELECT CAST(a AS CHARACTER VARYING(3))
< error
| 1064 42000 …line 1 column 34 near "VARYING(3))" 
> exec ^
: string
| SELECT CAST(a AS VARCHAR(3))
< error
| 1064 42000 …line 1 column 24 near "VARCHAR(3))" 
> exec ^
: string
| SELECT CAST(a AS VARCHAR)
< error
| 1064 42000 …line 1 column 24 near "VARCHAR)" 
> exec ^
: string
| SELECT CAST(a AS VARCHARACTER(3))
< error
| 1064 42000 …line 1 column 29 near "VARCHARACTER(3))" 
> exec ^
: string
| SELECT CAST(a AS NVARCHAR(3))
< error
| 1064 42000 …line 1 column 25 near "NVARCHAR(3))" 
> exec ^
: string
| SELECT CAST(a AS NATIONAL VARCHAR(3))
< error
| 1064 42000 …line 1 column 25 near "NATIONAL VARCHAR(3))" 
> exec ^
: string
| SELECT CAST(a AS NCHAR VARCHAR(3))
< error
| 1064 42000 …line 1 column 22 near "NCHAR VARCHAR(3))" 
> exec ^
: string
| SELECT CAST(a AS TEXT)
< error
| 1064 42000 …line 1 column 21 near "TEXT)" 
> exec ^
: string
| SELECT CAST(a AS TEXT(3))
< error
| 1064 42000 …line 1 column 21 near "TEXT(3))" 
> exec ^
: string
| SELECT CAST(a AS TINYTEXT)
< error
| 1064 42000 …line 1 column 25 near "TINYTEXT)" 
> exec ^
: string
| SELECT CAST(a AS MEDIUMTEXT)
< error
| 1064 42000 …line 1 column 27 near "MEDIUMTEXT)" 
> exec ^
: string
| SELECT CAST(a AS LONGTEXT)
< error
| 1064 42000 …line 1 column 25 near "LONGTEXT)" 
> exec ^
: string
| SELECT CAST(a AS LONG)
< error
| 1064 42000 …line 1 column 21 near "LONG)" 
> exec ^
: string
| SELECT CAST(a AS LONG VARCHAR)
< error
| 1064 42000 …line 1 column 21 near "LONG VARCHAR)" 
> exec ^
: string
| SELECT CAST(a AS BLOB)
< error
| 1064 42000 …line 1 column 21 near "BLOB)" 
> exec ^
: string
| SELECT CAST(a AS BLOB(3))
< error
| 1064 42000 …line 1 column 21 near "BLOB(3))" 
> exec ^
: string
| SELECT CAST(a AS TINYBLOB)
< error
| 1064 42000 …line 1 column 25 near "TINYBLOB)" 
> exec ^
: string
| SELECT CAST(a AS MEDIUMBLOB)
< error
| 1064 42000 …line 1 column 27 near "MEDIUMBLOB)" 
> exec ^
: string
| SELECT CAST(a AS LONGBLOB)
< error
| 1064 42000 …line 1 column 25 near "LONGBLOB)" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS VARBINARY(3))
< error
| 1064 42000 …line 1 column 26 near "VARBINARY(3))" 
> exec ^
: string
| SELECT CAST(a AS VARBINARY)
< error
| 1064 42000 …line 1 column 26 near "VARBINARY)" 
> exec ^
: string
| SELECT CAST(a AS BIT)
< error
| 1064 42000 …line 1 column 20 near "BIT)" 
> exec ^
: string
| SELECT CAST(a AS BIT(1))
< error
| 1064 42000 …line 1 column 20 near "BIT(1))" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS DATE(3))
< error
| 1064 42000 …line 1 column 22 near "(3))" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS TIMESTAMP)
< error
| 1064 42000 …line 1 column 26 near "TIMESTAMP)" 
> exec ^
: string
| SELECT CAST(a AS TIMESTAMP(3))
< error
| 1064 42000 …line 1 column 26 near "TIMESTAMP(3))" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS YEAR(4))
< error
| 1064 42000 …line 1 column 22 near "(4))" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS JSON(3))
< error
| 1064 42000 …line 1 column 22 near "(3))" 
> exec ^
: string
| SELECT CAST(a AS ENUM('a'))
< error
| 1064 42000 …line 1 column 21 near "ENUM('a'))" 
> exec ^
: string
| SELECT CAST(a AS SET('a'))
< error
| 1064 42000 …line 1 column 20 near "SET('a'))" 
> exec ^
: string
| SELECT CAST(a AS GEOMETRY)
< error
| 1064 42000 …line 1 column 25 near "GEOMETRY)" 
> exec ^
: string
| SELECT CAST(a AS POINT)
< error
| 1064 42000 …line 1 column 22 near "POINT)" 
> exec ^
: string
| SELECT CAST(a AS LINESTRING)
< error
| 1064 42000 …line 1 column 27 near "LINESTRING)" 
> exec ^
: string
| SELECT CAST(a AS POLYGON)
< error
| 1064 42000 …line 1 column 24 near "POLYGON)" 
> exec ^
: string
| SELECT CAST(a AS UUID)
< error
| 1064 42000 …line 1 column 21 near "UUID)" 
> exec ^
: string
| SELECT CAST(a AS INET4)
< error
| 1064 42000 …line 1 column 22 near "INET4)" 
> exec ^
: string
| SELECT CAST(a AS INET6)
< error
| 1064 42000 …line 1 column 22 near "INET6)" 
> exec ^
: string
| SELECT CAST(a AS SERIAL)
< error
| 1064 42000 …line 1 column 23 near "SERIAL)" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS CHAR ASCII)
< error
| 1064 42000 …line 1 column 27 near "ASCII)" 
> exec ^
: string
| SELECT CAST(a AS CHAR UNICODE)
< error
| 1064 42000 …line 1 column 29 near "UNICODE)" 
> exec ^
: string
| SELECT CAST(a AS CHAR BYTE)
< error
| 1064 42000 …line 1 column 26 near "BYTE)" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS BINARY CHARACTER SET utf8mb4)
< error
| 1064 42000 …line 1 column 33 near "CHARACTER SET utf8mb4)" 
> exec ^
: string
| SELECT CAST(a AS CHAR COLLATE utf8mb4_bin)
< error
| 1064 42000 …line 1 column 29 near "COLLATE utf8mb4_bin)" 
> exec ^
: string
| SELECT CAST(a AS VARCHAR(3) CHARACTER SET utf8mb4)
< error
| 1064 42000 …line 1 column 24 near "VARCHAR(3) CHARACTER SET utf8mb4)" 
> exec ^
: string
| SELECT CAST(a AS INTERVAL DAY)
< error
| 1064 42000 …line 1 column 25 near "INTERVAL DAY)" 
> exec ^
: string
| SELECT CAST(a AS DAY)
< error
| 1064 42000 …line 1 column 20 near "DAY)" 
> exec ^
: string
| SELECT CAST(a AS VECTOR)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS VECTOR)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS VECTOR(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS VECTOR(3))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CONVERT(a, INT)
< error
| 1064 42000 …line 1 column 21 near "INT)" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CONVERT(a, BIGINT)
< error
| 1064 42000 …line 1 column 24 near "BIGINT)" 
> exec ^
: string
| SELECT CONVERT(a, TEXT)
< error
| 1064 42000 …line 1 column 22 near "TEXT)" 
> exec ^
: string
| SELECT CONVERT(a, VARCHAR(3))
< error
| 1064 42000 …line 1 column 25 near "VARCHAR(3))" 
> exec ^
: string
| SELECT CONVERT(a, BOOLEAN)
< error
| 1064 42000 …line 1 column 25 near "BOOLEAN)" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS INT ARRAY)
< error
| 1064 42000 …line 1 column 20 near "INT ARRAY)" 
> exec ^
: string
| SELECT CAST(a AS BIGINT ARRAY)
< error
| 1064 42000 …line 1 column 23 near "BIGINT ARRAY)" 
> exec ^
: string
| SELECT CAST(a AS SIGNED ARRAY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS SIGNED ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS UNSIGNED ARRAY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS UNSIGNED ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS SIGNED INTEGER ARRAY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS SIGNED INTEGER ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS CHAR ARRAY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS CHAR(3) ARRAY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR(3) ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS BINARY ARRAY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS BINARY ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS BINARY(3) ARRAY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS BINARY(3) ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS DATE ARRAY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DATE ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS DATETIME ARRAY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DATETIME ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS TIME ARRAY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS TIME ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS YEAR ARRAY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS YEAR ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS JSON ARRAY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS JSON ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL ARRAY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DECIMAL ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS DOUBLE ARRAY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DOUBLE ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS FLOAT ARRAY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS FLOAT ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS REAL ARRAY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS REAL ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CONVERT(a, SIGNED ARRAY)
< error
| 1064 42000 …line 1 column 30 near "ARRAY)" 
> exec ^
: string
| SELECT CONVERT(a, CHAR(3) ARRAY)
< error
| 1064 42000 …line 1 column 31 near "ARRAY)" 
> exec ^
: string
| SELECT CAST(a AS JSON ARRAY ARRAY)
< error
| 1064 42000 …line 1 column 33 near "ARRAY)" 
> exec ^
: string
| SELECT CAST(a AS REAL(3))
< error
| 1064 42000 …line 1 column 22 near "(3))" 
> exec ^
: string
| SELECT CAST(a AS REAL(3,2))
< error
| 1064 42000 …line 1 column 22 near "(3,2))" 
> exec ^
: string
| SELECT CAST(a AS DOUBLE(3))
< error
| 1064 42000 …line 1 column 24 near "(3))" 
> exec ^
: string
| SELECT CAST(a AS DOUBLE PRECISION(3))
< error
| 1064 42000 …line 1 column 33 near "PRECISION(3))" 
> exec ^
: string
| SELECT CAST(a AS DOUBLE PRECISION(3,2))
< error
| 1064 42000 …line 1 column 33 near "PRECISION(3,2))" 
> exec ^
: string
| SELECT CAST(a AS FLOAT())
< error
| 1064 42000 …line 1 column 24 near "))" 
> exec ^
: string
| SELECT CAST(a AS FLOAT(3,))
< error
| 1064 42000 …line 1 column 26 near "))" 
> exec ^
: string
| SELECT CAST(a AS FLOAT(1.5))
< error
| 1064 42000 …line 1 column 26 near "1.5))" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL())
< error
| 1064 42000 …line 1 column 26 near "))" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL(3,2,1))
< error
| 1064 42000 …line 1 column 29 near ",1))" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS NUMERIC(3))
< error
| 1064 42000 …line 1 column 24 near "NUMERIC(3))" 
> exec ^
: string
| SELECT CAST(a AS FIXED(3))
< error
| 1064 42000 …line 1 column 22 near "FIXED(3))" 
> exec ^
: string
| SELECT CAST(a AS BIGINT(3))
< error
| 1064 42000 …line 1 column 23 near "BIGINT(3))" 
> exec ^
: string
| SELECT CAST(a AS BIGINT UNSIGNED)
< error
| 1064 42000 …line 1 column 23 near "BIGINT UNSIGNED)" 
> exec ^
: string
| SELECT CAST(a AS BIGINT SIGNED)
< error
| 1064 42000 …line 1 column 23 near "BIGINT SIGNED)" 
> exec ^
: string
| SELECT CAST(a AS INTEGER(3))
< error
| 1064 42000 …line 1 column 24 near "INTEGER(3))" 
> exec ^
: string
| SELECT CAST(a AS INTEGER UNSIGNED)
< error
| 1064 42000 …line 1 column 24 near "INTEGER UNSIGNED)" 
> exec ^
: string
| SELECT CAST(a AS INT ZEROFILL)
< error
| 1064 42000 …line 1 column 20 near "INT ZEROFILL)" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL ZEROFILL)
< error
| 1064 42000 …line 1 column 33 near "ZEROFILL)" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL SIGNED)
< error
| 1064 42000 …line 1 column 31 near "SIGNED)" 
> exec ^
: string
| SELECT CAST(a AS FLOAT SIGNED)
< error
| 1064 42000 …line 1 column 29 near "SIGNED)" 
> exec ^
: string
| SELECT CAST(a AS DOUBLE SIGNED)
< error
| 1064 42000 …line 1 column 30 near "SIGNED)" 
> exec ^
: string
| SELECT CAST(a AS SIGNED INT(3))
< error
| 1064 42000 …line 1 column 28 near "(3))" 
> exec ^
: string
| SELECT CAST(a AS UNSIGNED BIGINT)
< error
| 1064 42000 …line 1 column 32 near "BIGINT)" 
> exec ^
: string
| SELECT CAST(a AS SIGNED BIGINT)
< error
| 1064 42000 …line 1 column 30 near "BIGINT)" 
> exec ^
: string
| SELECT CAST(a AS SIGNED UNSIGNED)
< error
| 1064 42000 …line 1 column 32 near "UNSIGNED)" 
> exec ^
: string
| SELECT CAST(a AS CHAR())
< error
| 1064 42000 …line 1 column 23 near "))" 
> exec ^
: string
| SELECT CAST(a AS CHAR(3,2))
< error
| 1064 42000 …line 1 column 24 near ",2))" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS CHAR(3) COLLATE utf8mb4_bin)
< error
| 1064 42000 …line 1 column 32 near "COLLATE utf8mb4_bin)" 
> exec ^
: string
| SELECT CAST(a AS CHAR CHARACTER SET utf8mb4 COLLATE utf8mb4_bin)
< error
| 1064 42000 …line 1 column 51 near "COLLATE utf8mb4_bin)" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS CHAR BINARY BINARY)
< error
| 1064 42000 …line 1 column 35 near "BINARY)" 
> exec ^
: string
| SELECT CAST(a AS CHAR(3) ASCII)
< error
| 1064 42000 …line 1 column 30 near "ASCII)" 
> exec ^
: string
| SELECT CAST(a AS CHAR(3) UNICODE)
< error
| 1064 42000 …line 1 column 32 near "UNICODE)" 
> exec ^
: string
| SELECT CAST(a AS CHAR(3) BYTE)
< error
| 1064 42000 …line 1 column 29 near "BYTE)" 
> exec ^
: string
| SELECT CAST(a AS CHAR ASCII BINARY)
< error
| 1064 42000 …line 1 column 27 near "ASCII BINARY)" 
> exec ^
: string
| SELECT CAST(a AS CHAR BINARY ASCII)
< error
| 1064 42000 …line 1 column 34 near "ASCII)" 
> exec ^
: string
| SELECT CAST(a AS CHAR UNICODE BINARY)
< error
| 1064 42000 …line 1 column 29 near "UNICODE BINARY)" 
> exec ^
: string
| SELECT CAST(a AS CHAR BINARY UNICODE)
< error
| 1064 42000 …line 1 column 36 near "UNICODE)" 
> exec ^
: string
| SELECT CAST(a AS CHAR ASCII UNICODE)
< error
| 1064 42000 …line 1 column 27 near "ASCII UNICODE)" 
> exec ^
: string
| SELECT CAST(a AS CHAR ASCII CHARACTER SET utf8mb4)
< error
| 1064 42000 …line 1 column 27 near "ASCII CHARACTER SET utf8mb4)" 
> exec ^
: string
| SELECT CAST(a AS CHAR BYTE BINARY)
< error
| 1064 42000 …line 1 column 26 near "BYTE BINARY)" 
> exec ^
: string
| SELECT CAST(a AS CHAR(3) ASCII BINARY)
< error
| 1064 42000 …line 1 column 30 near "ASCII BINARY)" 
> exec ^
: string
| SELECT CAST(a AS CHARACTER ASCII)
< error
| 1064 42000 …line 1 column 32 near "ASCII)" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS NCHAR BINARY)
< error
| 1064 42000 …line 1 column 22 near "NCHAR BINARY)" 
> exec ^
: string
| SELECT CAST(a AS NCHAR CHARACTER SET utf8mb4)
< error
| 1064 42000 …line 1 column 22 near "NCHAR CHARACTER SET utf8mb4)" 
> exec ^
: string
| SELECT CAST(a AS NCHAR ASCII)
< error
| 1064 42000 …line 1 column 22 near "NCHAR ASCII)" 
> exec ^
: string
| SELECT CAST(a AS NATIONAL CHAR BINARY)
< error
| 1064 42000 …line 1 column 25 near "NATIONAL CHAR BINARY)" 
> exec ^
: string
| SELECT CAST(a AS NATIONAL CHARACTER)
< error
| 1064 42000 …line 1 column 25 near "NATIONAL CHARACTER)" 
> exec ^
: string
| SELECT CAST(a AS NATIONAL CHARACTER(3))
< error
| 1064 42000 …line 1 column 25 near "NATIONAL CHARACTER(3))" 
> exec ^
: string
| SELECT CAST(a AS NATIONAL)
< error
| 1064 42000 …line 1 column 25 near "NATIONAL)" 
> exec ^
: string
| SELECT CAST(a AS VARCHAR(3) BINARY)
< error
| 1064 42000 …line 1 column 24 near "VARCHAR(3) BINARY)" 
> exec ^
: string
| SELECT CAST(a AS VARCHAR(3) COLLATE utf8mb4_bin)
< error
| 1064 42000 …line 1 column 24 near "VARCHAR(3) COLLATE utf8mb4_bin)" 
> exec ^
: string
| SELECT CAST(a AS VARCHAR(3) ASCII)
< error
| 1064 42000 …line 1 column 24 near "VARCHAR(3) ASCII)" 
> exec ^
: string
| SELECT CAST(a AS VARCHARACTER(3) CHARACTER SET utf8mb4)
< error
| 1064 42000 …line 1 column 29 near "VARCHARACTER(3) CHARACTER SET utf8mb4)" 
> exec ^
: string
| SELECT CAST(a AS VARCHAR())
< error
| 1064 42000 …line 1 column 24 near "VARCHAR())" 
> exec ^
: string
| SELECT CAST(a AS BINARY())
< error
| 1064 42000 …line 1 column 25 near "))" 
> exec ^
: string
| SELECT CAST(a AS BINARY(3,2))
< error
| 1064 42000 …line 1 column 26 near ",2))" 
> exec ^
: string
| SELECT CAST(a AS BINARY BINARY)
< error
| 1064 42000 …line 1 column 30 near "BINARY)" 
> exec ^
: string
| SELECT CAST(a AS BINARY ASCII)
< error
| 1064 42000 …line 1 column 29 near "ASCII)" 
> exec ^
: string
| SELECT CAST(a AS BINARY COLLATE utf8mb4_bin)
< error
| 1064 42000 …line 1 column 31 near "COLLATE utf8mb4_bin)" 
> exec ^
: string
| SELECT CAST(a AS DATE BINARY)
< error
| 1064 42000 …line 1 column 28 near "BINARY)" 
> exec ^
: string
| SELECT CAST(a AS TIME())
< error
| 1064 42000 …line 1 column 23 near "))" 
> exec ^
: string
| SELECT CAST(a AS TIME(3,2))
< error
| 1064 42000 …line 1 column 24 near ",2))" 
> exec ^
: string
| SELECT CAST(a AS TIME(1.5))
< error
| 1064 42000 …line 1 column 25 near "1.5))" 
> exec ^
: string
| SELECT CAST(a AS DATETIME())
< error
| 1064 42000 …line 1 column 27 near "))" 
> exec ^
: string
| SELECT CAST(a AS DATETIME(3,2))
< error
| 1064 42000 …line 1 column 28 near ",2))" 
> exec ^
: string
| SELECT CAST(a AS YEAR(2))
< error
| 1064 42000 …line 1 column 22 near "(2))" 
> exec ^
: string
| SELECT CAST(a AS YEAR())
< error
| 1064 42000 …line 1 column 22 near "())" 
> exec ^
: string
| SELECT CAST(a AS JSON BINARY)
< error
| 1064 42000 …line 1 column 28 near "BINARY)" 
> exec ^
: string
| SELECT CAST(a AS POINT(3))
< error
| 1064 42000 …line 1 column 22 near "POINT(3))" 
> exec ^
: string
| SELECT CAST(a AS MULTIPOINT)
< error
| 1064 42000 …line 1 column 27 near "MULTIPOINT)" 
> exec ^
: string
| SELECT CAST(a AS MULTILINESTRING)
< error
| 1064 42000 …line 1 column 32 near "MULTILINESTRING)" 
> exec ^
: string
| SELECT CAST(a AS MULTIPOLYGON)
< error
| 1064 42000 …line 1 column 29 near "MULTIPOLYGON)" 
> exec ^
: string
| SELECT CAST(a AS GEOMETRYCOLLECTION)
< error
| 1064 42000 …line 1 column 35 near "GEOMETRYCOLLECTION)" 
> exec ^
: string
| SELECT CAST(a AS GEOMCOLLECTION)
< error
| 1064 42000 …line 1 column 31 near "GEOMCOLLECTION)" 
> exec ^
: string
| SELECT CAST(a AS GEOMETRY(3))
< error
| 1064 42000 …line 1 column 25 near "GEOMETRY(3))" 
> exec ^
: string
| SELECT CAST(a AS CHAR SET utf8mb4)
< error
| 1064 42000 …line 1 column 25 near "SET utf8mb4)" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS CHAR CHARACTER SET)
< error
| 1064 42000 …line 1 column 36 near ")" 
> exec ^
: string
| SELECT CAST(a AS CHAR COLLATE)
< error
| 1064 42000 …line 1 column 29 near "COLLATE)" 
> exec ^
: string
| SELECT CAST(a AS `CHAR`)
< error
| 1064 42000 …line 1 column 23 near "`CHAR`)" 
> exec ^
: string
| SELECT CAST(a AS `SIGNED`)
< error
| 1064 42000 …line 1 column 25 near "`SIGNED`)" 
> exec ^
: string
| SELECT CAST(a AS `DATE`)
< error
| 1064 42000 …line 1 column 23 near "`DATE`)" 
> exec ^
: string
| SELECT CAST(a AS 'CHAR')
< error
| 1064 42000 …line 1 column 23 near "'CHAR')" 
> exec ^
: string
| SELECT CAST(a AS 1)
< error
| 1064 42000 …line 1 column 18 near "1)" 
> exec ^
: string
| SELECT CAST(a AS NULL)
< error
| 1064 42000 …line 1 column 21 near "NULL)" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CONVERT(a, NCHAR)
< error
| 1064 42000 …line 1 column 23 near "NCHAR)" 
> exec ^
: string
| SELECT CONVERT(a, POINT)
< error
| 1064 42000 …line 1 column 23 near "POINT)" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CONVERT(a, GEOMETRY)
< error
| 1064 42000 …line 1 column 26 near "GEOMETRY)" 
> exec ^
: string
| SELECT CONVERT(a, UUID)
< error
| 1064 42000 …line 1 column 22 near "UUID)" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CONVERT(a, )
< error
| 1064 42000 …line 1 column 19 near ")" 
> exec ^
: string
| SELECT CONVERT(a)
< error
| 1064 42000 …line 1 column 17 near ")" 
> exec ^
: string
| SELECT CAST(a AS)
< error
| 1064 42000 …line 1 column 17 near ")" 
> exec ^
: string
| SELECT CAST(a)
< error
| 1064 42000 …line 1 column 14 near ")" 
> exec ^
: string
| SELECT CAST(a AS CHAR, b)
< error
| 1064 42000 …line 1 column 22 near ", b)" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS NCHAR ARRAY)
< error
| 1064 42000 …line 1 column 22 near "NCHAR ARRAY)" 
> exec ^
: string
| SELECT CAST(a AS POINT ARRAY)
< error
| 1064 42000 …line 1 column 22 near "POINT ARRAY)" 
> exec ^
: string
| SELECT CAST(a AS GEOMETRY ARRAY)
< error
| 1064 42000 …line 1 column 25 near "GEOMETRY ARRAY)" 
> exec ^
: string
| SELECT CAST(a AS VECTOR ARRAY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS VECTOR ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS BIGINT UNSIGNED ARRAY)
< error
| 1064 42000 …line 1 column 23 near "BIGINT UNSIGNED ARRAY)" 
> exec ^
: string
| CREATE TABLE t9 (a int, KEY k ((CAST(a AS SIGNED ARRAY))))
< ok 0 0
> exec ^
: string
| CREATE TABLE t9 (a int, KEY k ((CAST(a AS INT ARRAY))))
< error
| 1064 42000 …line 1 column 45 near "INT ARRAY))))" 
> exec ^
: string
| CREATE TABLE t9 (a int, b int AS (CAST(a AS CHAR)))
< ok 0 0
> exec ^
: string
| CREATE TABLE t9 (a int, b int AS (CAST(a AS TEXT)))
< error
| 1064 42000 …line 1 column 48 near "TEXT)))" 
> exec ^
: string
| SELECT 1 FROM t WHERE CAST(a AS BOOLEAN)
< error
| 1064 42000 …line 1 column 39 near "BOOLEAN)" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS NUMBER)
< error
| 1064 42000 …line 1 column 23 near "NUMBER)" 
> exec ^
: string
| SELECT CAST(a AS NUMBER(3,2))
< error
| 1064 42000 …line 1 column 23 near "NUMBER(3,2))" 
> exec ^
: string
| SELECT CAST(a AS ROARINGBITMAP)
< error
| 1064 42000 …line 1 column 30 near "ROARINGBITMAP)" 
> exec ^
: string
| SELECT CAST(a AS ARRAY(INT))
< error
| 1064 42000 …line 1 column 22 near "ARRAY(INT))" 
> exec ^
: string
| SELECT CAST(a AS ARRAY<INT>)
< error
| 1064 42000 …line 1 column 22 near "ARRAY<INT>)" 
> exec ^
: string
| SELECT CAST(a AS CHAR(1.5))
< error
| 1064 42000 …line 1 column 25 near "1.5))" 
> exec ^
: string
| SELECT CAST(a AS BINARY(1.5))
< error
| 1064 42000 …line 1 column 27 near "1.5))" 
> exec ^
: string
| SELECT CAST(a AS NCHAR(1.5))
< error
| 1064 42000 …line 1 column 22 near "NCHAR(1.5))" 
> exec ^
: string
| SELECT CAST(a AS VARCHAR(1.5))
< error
| 1064 42000 …line 1 column 24 near "VARCHAR(1.5))" 
> exec ^
: string
| SELECT CAST(a AS DATETIME(1.5))
< error
| 1064 42000 …line 1 column 29 near "1.5))" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL(1.5))
< error
| 1064 42000 …line 1 column 28 near "1.5))" 
> exec ^
: string
| SELECT CAST(a AS DATE(1.5))
< error
| 1064 42000 …line 1 column 22 near "(1.5))" 
> exec ^
: string
| SELECT CAST(a AS FLOAT(1e5))
< error
| 1064 42000 …line 1 column 26 near "1e5))" 
> exec ^
: string
| SELECT CAST(a AS TIME(1e5))
< error
| 1064 42000 …line 1 column 25 near "1e5))" 
> exec ^
: string
| SELECT CAST(a AS VECTOR<FLOAT>(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS VECTOR<FLOAT>(3))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS VECTOR<FLOAT>)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS VECTOR<FLOAT>)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS VECTOR<DOUBLE>)
< error
| 1064 42000 …line 1 column 32 near ")"Only VECTOR is supported for now 
> exec ^
: string
| SELECT CAST(a AS VECTOR())
< error
| 1064 42000 …line 1 column 25 near "))" 
> exec ^
: string
| SELECT CAST(a AS VECTOR(3,2))
< error
| 1064 42000 …line 1 column 26 near ",2))" 
> exec ^
: string
| SELECT CAST(a AS SIGNED INT INTEGER)
< error
| 1064 42000 …line 1 column 35 near "INTEGER)" 
> exec ^
: string
| SELECT CAST(a AS SIGNED INTEGER INT)
< error
| 1064 42000 …line 1 column 35 near "INT)" 
> exec ^
: string
| SELECT CAST(a AS UNSIGNED INT INT)
< error
| 1064 42000 …line 1 column 33 near "INT)" 
> exec ^
: string
| SELECT CAST(a AS NATIONAL VARCHAR)
< error
| 1064 42000 …line 1 column 25 near "NATIONAL VARCHAR)" 
> exec ^
: string
| SELECT CAST(a AS NCHAR VARYING(3))
< error
| 1064 42000 …line 1 column 22 near "NCHAR VARYING(3))" 
> exec ^
: string
| SELECT CAST(a AS CHAR(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin BINARY)
< error
| 1064 42000 …line 1 column 54 near "COLLATE utf8mb4_bin BINARY)" 
> exec ^
: string
| SELECT CAST(a AS CHAR COLLATE utf8mb4_bin BINARY)
< error
| 1064 42000 …line 1 column 29 near "COLLATE utf8mb4_bin BINARY)" 
> exec ^
: string
| SELECT CAST(a AS CHAR BINARY COLLATE utf8mb4_bin)
< error
| 1064 42000 …line 1 column 36 near "COLLATE utf8mb4_bin)" 
> exec ^
: string
| SELECT CAST(a AS CHAR ASCII COLLATE utf8mb4_bin)
< error
| 1064 42000 …line 1 column 27 near "ASCII COLLATE utf8mb4_bin)" 
> exec ^
: string
| SELECT CAST(a AS BOOL ARRAY)
< error
| 1064 42000 …line 1 column 21 near "BOOL ARRAY)" 
> exec ^
: string
| SELECT * FROM JSON_TABLE('[]', '$[*]' COLUMNS (a INT PATH '$')) AS j
< error
| 1064 42000 …line 1 column 25 near "('[]', '$[*]' COLUMNS (a INT PATH '$')) AS j" 
> exec ^
: string
| SELECT * FROM `json_table`('[1]', '$[*]' COLUMNS (x INT PATH '$')) AS jt
< error
| 1064 42000 …line 1 column 27 near "('[1]', '$[*]' COLUMNS (x INT PATH '$')) AS jt" 
> exec ^
: string
| SELECT json_table.a FROM json_table
< error
| 1146 42S02 Table 'test.json_table' doesn't exist
> exec ^
: string
| SELECT json_table.a FROM json_table
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a FROM json_table
< error
| 1146 42S02 Table 'test.json_table' doesn't exist
> exec ^
: string
| SELECT a FROM json_table
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a FROM t json_table
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT a FROM t json_table
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE VIEW json_table AS SELECT 1
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD COLUMN json_table INT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD COLUMN json_table INT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DROP TABLE json_table
< ok 0 0
> exec ^
: string
| INSERT INTO json_table VALUES (1)
< error
| 1146 42S02 Table 'test.json_table' doesn't exist
> exec ^
: string
| INSERT INTO json_table VALUES (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT json_table(1)
< error
| 1305 42000 FUNCTION test.json_table does not exist
> exec ^
: string
| SELECT json_table(1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT `json_table`(1)
< error
| 1305 42000 FUNCTION test.json_table does not exist
> exec ^
: string
| SELECT `json_table`(1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM json_table('[]', '$' COLUMNS (a INT PATH '$')) j
< error
| 1064 42000 …line 1 column 25 near "('[]', '$' COLUMNS (a INT PATH '$')) j" 
> exec ^
: string
| SELECT * FROM t, json_table(a, '$' COLUMNS (b INT PATH '$')) j
< error
| 1064 42000 …line 1 column 28 near "(a, '$' COLUMNS (b INT PATH '$')) j" 
