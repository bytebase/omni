-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestCollateConvertWeightFormsOracle default#1
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| ALTER TABLE t ADD c JSON COLLATE BINARY
< error
| 1064 42000 …near 'BINARY' at line 1
> exec ^
: string
| ALTER TABLE t ADD c JSON COLLATE utf8mb4_bin
< ok 0 0
> exec ^
: string
| SELECT x COLLATE BINARY
< error
| 1064 42000 …near 'BINARY' at line 1
> exec ^
: string
| SELECT x COLLATE utf8mb4_bin
< error
| 1054 42S22 Unknown column 'x' in 'SELECT'
> exec ^
: string
| SELECT x COLLATE utf8mb4_bin
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CONVERT(x USING BINARY)
< error
| 1054 42S22 Unknown column 'x' in 'SELECT'
> exec ^
: string
| SELECT CONVERT(x USING BINARY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CONVERT(x USING utf8mb4)
< error
| 1054 42S22 Unknown column 'x' in 'SELECT'
> exec ^
: string
| SELECT CONVERT(x USING utf8mb4)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT WEIGHT_STRING(x)
< error
| 1054 42S22 Unknown column 'x' in 'SELECT'
> exec ^
: string
| SELECT WEIGHT_STRING(x)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT WEIGHT_STRING(x AS CHAR(1))
< error
| 1054 42S22 Unknown column 'x' in 'SELECT'
> exec ^
: string
| SELECT WEIGHT_STRING(x AS CHAR(1))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT WEIGHT_STRING(x AS CHARACTER(1))
< error
| 1054 42S22 Unknown column 'x' in 'SELECT'
> exec ^
: string
| SELECT WEIGHT_STRING(x AS CHARACTER(1))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT WEIGHT_STRING(x AS BINARY(1))
< error
| 1054 42S22 Unknown column 'x' in 'SELECT'
> exec ^
: string
| SELECT WEIGHT_STRING(x AS BINARY(1))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT WEIGHT_STRING(x AS CHAR(0))
< error
| 1064 42000 …near '0))' at line 1
> exec ^
: string
| SELECT WEIGHT_STRING(x AS CHAR)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT WEIGHT_STRING(x AS DECIMAL(1))
< error
| 1064 42000 …near 'DECIMAL(1))' at line 1
