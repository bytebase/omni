-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestCollateConvertWeightFormsOracle default#1
> exec
| ALTER TABLE t ADD c JSON COLLATE BINARY
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| ALTER TABLE t ADD c JSON COLLATE BINARY
< ok 0 0
> exec
| ALTER TABLE t ADD c JSON COLLATE utf8mb4_bin
< error
| 1253 42000 COLLATION 'utf8mb4_bin' is not valid for CHARACTER SET 'binary'
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| ALTER TABLE t ADD c JSON COLLATE utf8mb4_bin
< error
| 1253 42000 COLLATION 'utf8mb4_bin' is not valid for CHARACTER SET 'binary'
> exec
| DO x COLLATE BINARY
< error
| 1064 42000 …near 'BINARY' at line 1
> exec
| DO x COLLATE utf8mb4_bin
< error
| 1054 42S22 Unknown column 'x' in 'field list'
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| DO x COLLATE utf8mb4_bin
< error
| 1054 42S22 Unknown column 'x' in 'field list'
> exec ^
: string
| DO x COLLATE utf8mb4_bin
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| DO CONVERT(x USING BINARY)
< error
| 1054 42S22 Unknown column 'x' in 'field list'
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| DO CONVERT(x USING BINARY)
< error
| 1054 42S22 Unknown column 'x' in 'field list'
> exec ^
: string
| DO CONVERT(x USING BINARY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| DO CONVERT(x USING utf8mb4)
< error
| 1054 42S22 Unknown column 'x' in 'field list'
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| DO CONVERT(x USING utf8mb4)
< error
| 1054 42S22 Unknown column 'x' in 'field list'
> exec ^
: string
| DO CONVERT(x USING utf8mb4)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| DO WEIGHT_STRING(x)
< error
| 1054 42S22 Unknown column 'x' in 'field list'
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| DO WEIGHT_STRING(x)
< error
| 1054 42S22 Unknown column 'x' in 'field list'
> exec ^
: string
| DO WEIGHT_STRING(x)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| DO WEIGHT_STRING(x AS CHAR(1))
< error
| 1054 42S22 Unknown column 'x' in 'field list'
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| DO WEIGHT_STRING(x AS CHAR(1))
< error
| 1054 42S22 Unknown column 'x' in 'field list'
> exec ^
: string
| DO WEIGHT_STRING(x AS CHAR(1))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| DO WEIGHT_STRING(x AS CHARACTER(1))
< error
| 1054 42S22 Unknown column 'x' in 'field list'
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| DO WEIGHT_STRING(x AS CHARACTER(1))
< error
| 1054 42S22 Unknown column 'x' in 'field list'
> exec ^
: string
| DO WEIGHT_STRING(x AS CHARACTER(1))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| DO WEIGHT_STRING(x AS BINARY(1))
< error
| 1054 42S22 Unknown column 'x' in 'field list'
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| DO WEIGHT_STRING(x AS BINARY(1))
< error
| 1054 42S22 Unknown column 'x' in 'field list'
> exec ^
: string
| DO WEIGHT_STRING(x AS BINARY(1))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| DO WEIGHT_STRING(x AS CHAR(0))
< error
| 1064 42000 …near '0))' at line 1
> exec
| DO WEIGHT_STRING(x AS CHAR)
< error
| 1064 42000 …near ')' at line 1
> exec
| DO WEIGHT_STRING(x AS DECIMAL(1))
< error
| 1064 42000 …near 'DECIMAL(1))' at line 1
> exec
| DO WEIGHT_STRING(x LEVEL 1)
< error
| 1064 42000 …near 'LEVEL 1)' at line 1
> exec
| DO WEIGHT_STRING(x LEVEL 0)
< error
| 1064 42000 …near 'LEVEL 0)' at line 1
> exec
| DO WEIGHT_STRING(x LEVEL 1 ASC)
< error
| 1064 42000 …near 'LEVEL 1 ASC)' at line 1
> exec
| DO WEIGHT_STRING(x LEVEL 1 DESC)
< error
| 1064 42000 …near 'LEVEL 1 DESC)' at line 1
> exec
| DO WEIGHT_STRING(x LEVEL 1 REVERSE)
< error
| 1064 42000 …near 'LEVEL 1 REVERSE)' at line 1
> exec
| DO WEIGHT_STRING(x LEVEL 1 DESC REVERSE)
< error
| 1064 42000 …near 'LEVEL 1 DESC REVERSE)' at line 1
> exec
| DO WEIGHT_STRING(x LEVEL 1 REVERSE DESC)
< error
| 1064 42000 …near 'LEVEL 1 REVERSE DESC)' at line 1
> exec
| DO WEIGHT_STRING(x LEVEL 1, 2)
< error
| 1064 42000 …near 'LEVEL 1, 2)' at line 1
> exec
| DO WEIGHT_STRING(x LEVEL 1 DESC, 2 REVERSE)
< error
| 1064 42000 …near 'LEVEL 1 DESC, 2 REVERSE)' at line 1
> exec
| DO WEIGHT_STRING(x LEVEL 1 - 2)
< error
| 1064 42000 …near 'LEVEL 1 - 2)' at line 1
> exec
| DO WEIGHT_STRING(x AS CHAR(4) LEVEL 1)
< error
| 1064 42000 …near 'LEVEL 1)' at line 1
> exec
| DO WEIGHT_STRING(x AS BINARY(4) LEVEL 1)
< error
| 1064 42000 …near 'LEVEL 1)' at line 1
> exec
| DO WEIGHT_STRING(x LEVEL)
< error
| 1064 42000 …near 'LEVEL)' at line 1
