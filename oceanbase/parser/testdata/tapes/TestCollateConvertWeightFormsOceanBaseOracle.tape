-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestCollateConvertWeightFormsOceanBaseOracle default#1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| ALTER TABLE t ADD c JSON COLLATE BINARY
< error
| 1064 42000 …near 'BINARY' at line 1
> exec ^
: string
| ALTER TABLE t ADD c JSON COLLATE utf8mb4_bin
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT x COLLATE BINARY
< error
| 1064 42000 …near 'BINARY' at line 1
> exec ^
: string
| SELECT x COLLATE utf8mb4_bin
< error
| 1054 42S22 Unknown column 'x' in 'field list'
> exec ^
: string
| SELECT CONVERT(x USING BINARY)
< error
| 1054 42S22 Unknown column 'x' in 'field list'
> exec ^
: string
| SELECT CONVERT(x USING utf8mb4)
< error
| 1054 42S22 Unknown column 'x' in 'field list'
> exec ^
: string
| SELECT WEIGHT_STRING(x)
< error
| 1054 42S22 Unknown column 'x' in 'field list'
> exec ^
: string
| SELECT WEIGHT_STRING(x AS CHAR(1))
< error
| 1054 42S22 Unknown column 'x' in 'field list'
> exec ^
: string
| SELECT WEIGHT_STRING(x AS CHARACTER(1))
< error
| 1054 42S22 Unknown column 'x' in 'field list'
> exec ^
: string
| SELECT WEIGHT_STRING(x AS BINARY(1))
< error
| 1054 42S22 Unknown column 'x' in 'field list'
> exec ^
: string
| SELECT WEIGHT_STRING(x AS CHAR(0))
< error
| 1064 42000 …near '(0))' at line 1
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
> exec ^
: string
| SELECT WEIGHT_STRING(x LEVEL 1)
< error
| 1054 42S22 Unknown column 'x' in 'field list'
> exec ^
: string
| SELECT WEIGHT_STRING(x LEVEL 0)
< error
| 1054 42S22 Unknown column 'x' in 'field list'
> exec ^
: string
| SELECT WEIGHT_STRING(x LEVEL 7)
< error
| 1054 42S22 Unknown column 'x' in 'field list'
> exec ^
: string
| SELECT WEIGHT_STRING(x LEVEL 1 ASC)
< error
| 1054 42S22 Unknown column 'x' in 'field list'
> exec ^
: string
| SELECT WEIGHT_STRING(x LEVEL 1 DESC)
< error
| 1054 42S22 Unknown column 'x' in 'field list'
> exec ^
: string
| SELECT WEIGHT_STRING(x LEVEL 1 REVERSE)
< error
| 1054 42S22 Unknown column 'x' in 'field list'
> exec ^
: string
| SELECT WEIGHT_STRING(x LEVEL 1 ASC REVERSE)
< error
| 1054 42S22 Unknown column 'x' in 'field list'
> exec ^
: string
| SELECT WEIGHT_STRING(x LEVEL 1 DESC REVERSE)
< error
| 1054 42S22 Unknown column 'x' in 'field list'
> exec ^
: string
| SELECT WEIGHT_STRING(x LEVEL 1 REVERSE DESC)
< error
| 1064 42000 …near 'DESC)' at line 1
> exec ^
: string
| SELECT WEIGHT_STRING(x LEVEL 1, 2)
< error
| 1054 42S22 Unknown column 'x' in 'field list'
> exec ^
: string
| SELECT WEIGHT_STRING(x LEVEL 1 DESC, 2 REVERSE)
< error
| 1054 42S22 Unknown column 'x' in 'field list'
> exec ^
: string
| SELECT WEIGHT_STRING(x LEVEL 1 - 2)
< error
| 1054 42S22 Unknown column 'x' in 'field list'
> exec ^
: string
| SELECT WEIGHT_STRING(x LEVEL 2 - 1)
< error
| 1054 42S22 Unknown column 'x' in 'field list'
> exec ^
: string
| SELECT WEIGHT_STRING(x LEVEL 1 ASC - 2)
< error
| 1064 42000 …near '- 2)' at line 1
> exec ^
: string
| SELECT WEIGHT_STRING(x LEVEL 1 DESC - 2)
< error
| 1064 42000 …near '- 2)' at line 1
> exec ^
: string
| SELECT WEIGHT_STRING(x LEVEL 1 - 2 DESC)
< error
| 1064 42000 …near 'DESC)' at line 1
> exec ^
: string
| SELECT WEIGHT_STRING(x LEVEL 1, 2 - 3)
< error
| 1064 42000 …near '- 3)' at line 1
> exec ^
: string
| SELECT WEIGHT_STRING(x LEVEL 1 - 2, 3)
< error
| 1064 42000 …near ', 3)' at line 1
> exec ^
: string
| SELECT WEIGHT_STRING(x LEVEL)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT WEIGHT_STRING(x AS CHAR(4) LEVEL 1)
< error
| 1054 42S22 Unknown column 'x' in 'field list'
> exec ^
: string
| SELECT WEIGHT_STRING(x AS CHARACTER(4) LEVEL 1 - 2)
< error
| 1054 42S22 Unknown column 'x' in 'field list'
> exec ^
: string
| SELECT WEIGHT_STRING(x AS BINARY(4) LEVEL 1)
< error
| 1064 42000 …near 'LEVEL 1)' at line 1
> exec ^
: string
| SELECT WEIGHT_STRING(x AS BINARY(4) LEVEL 1 DESC)
< error
| 1064 42000 …near 'LEVEL 1 DESC)' at line 1
