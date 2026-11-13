-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestCollateConvertWeightFormsMatchTiDB default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| ALTER TABLE t ADD c JSON COLLATE BINARY
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD c JSON COLLATE BINARY
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD c JSON COLLATE utf8mb4_bin
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD c JSON COLLATE utf8mb4_bin
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a INT COLLATE BINARY)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT COLLATE utf8mb4_bin)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT COLLATE 'utf8mb4_bin')
< ok 0 0
> exec ^
: string
| DO x COLLATE BINARY
< error
| 1054 42S22 Unknown column 'x' in ''
> exec ^
: string
| DO x COLLATE BINARY
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DO x COLLATE utf8mb4_bin
< error
| 1054 42S22 Unknown column 'x' in ''
> exec ^
: string
| DO x COLLATE utf8mb4_bin
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DO x COLLATE 'utf8mb4_bin'
< error
| 1054 42S22 Unknown column 'x' in ''
> exec ^
: string
| DO x COLLATE 'utf8mb4_bin'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DO x COLLATE
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| DO CONVERT(x USING BINARY)
< error
| 1054 42S22 Unknown column 'x' in ''
> exec ^
: string
| DO CONVERT(x USING BINARY)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DO CONVERT(x USING utf8mb4)
< error
| 1054 42S22 Unknown column 'x' in ''
> exec ^
: string
| DO CONVERT(x USING utf8mb4)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DO CONVERT(x USING 'utf8mb4')
< error
| 1054 42S22 Unknown column 'x' in ''
> exec ^
: string
| DO CONVERT(x USING 'utf8mb4')
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DO CONVERT(x USING)
< error
| 1064 42000 …line 1 column 19 near ")" 
> exec ^
: string
| DO WEIGHT_STRING(x)
< error
| 1054 42S22 Unknown column 'x' in ''
> exec ^
: string
| DO WEIGHT_STRING(x)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DO WEIGHT_STRING(x AS CHAR(1))
< error
| 1054 42S22 Unknown column 'x' in ''
> exec ^
: string
| DO WEIGHT_STRING(x AS CHAR(1))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DO WEIGHT_STRING(x AS CHARACTER(1))
< error
| 1054 42S22 Unknown column 'x' in ''
> exec ^
: string
| DO WEIGHT_STRING(x AS CHARACTER(1))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DO WEIGHT_STRING(x AS BINARY(1))
< error
| 1054 42S22 Unknown column 'x' in ''
> exec ^
: string
| DO WEIGHT_STRING(x AS BINARY(1))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DO WEIGHT_STRING(x AS CHAR(0))
< error
| 1054 42S22 Unknown column 'x' in ''
> exec ^
: string
| DO WEIGHT_STRING(x AS CHAR(0))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DO WEIGHT_STRING(x AS CHAR)
< error
| 1064 42000 …line 1 column 27 near ")" 
> exec ^
: string
| DO WEIGHT_STRING(x AS BINARY)
< error
| 1064 42000 …line 1 column 29 near ")" 
> exec ^
: string
| DO WEIGHT_STRING(x AS DECIMAL(1))
< error
| 1064 42000 …line 1 column 29 near "DECIMAL(1))" 
> exec ^
: string
| DO WEIGHT_STRING(x AS)
< error
| 1064 42000 …line 1 column 22 near ")" 
> exec ^
: string
| DO WEIGHT_STRING(x LEVEL 1)
< error
| 1064 42000 …line 1 column 24 near "LEVEL 1)" 
> exec ^
: string
| DO WEIGHT_STRING(x LEVEL 1 DESC REVERSE)
< error
| 1064 42000 …line 1 column 24 near "LEVEL 1 DESC REVERSE)" 
> exec ^
: string
| DO WEIGHT_STRING(x LEVEL 1, 2)
< error
| 1064 42000 …line 1 column 24 near "LEVEL 1, 2)" 
> exec ^
: string
| DO WEIGHT_STRING(x LEVEL 1 - 2)
< error
| 1064 42000 …line 1 column 24 near "LEVEL 1 - 2)" 
> exec ^
: string
| DO WEIGHT_STRING(x AS CHAR(4) LEVEL 1)
< error
| 1064 42000 …line 1 column 35 near "LEVEL 1)" 
> exec ^
: string
| DO WEIGHT_STRING(x AS BINARY(4) LEVEL 1)
< error
| 1064 42000 …line 1 column 37 near "LEVEL 1)" 
