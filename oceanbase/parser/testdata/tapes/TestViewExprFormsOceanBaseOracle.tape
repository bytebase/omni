-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestViewExprFormsOceanBaseOracle default#1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select N'nat'
< ok 0 0
> exec ^
: string
| select n'nat'
< ok 0 0
> exec ^
: string
| select N''
< ok 0 0
> exec ^
: string
| select N'x' 'y'
< ok 0 0
> exec ^
: string
| select N'x' 'y' 'z'
< ok 0 0
> exec ^
: string
| select N'x' as y
< ok 0 0
> exec ^
: string
| select N 'x' from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select N 'x' as y from t
< error
| 1064 42000 …near 'as y from t' at line 1
> exec ^
: string
| select `N`'x' from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select N"x" from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select N"x" as y from t
< error
| 1064 42000 …near 'as y from t' at line 1
> exec ^
: string
| select trim(both from c) from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select trim(leading from c) from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select trim(trailing from c) from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select trim(both 'x' from c) from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select trim('x' from c) from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select trim(c) from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select trim(leading 'x')
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select trim(both)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select trim(both from)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select trim(both from 'x' from 'y')
< error
| 1064 42000 …near 'from 'y')' at line 1
> exec ^
: string
| select char(65 using utf8mb4)
< ok 0 0
> exec ^
: string
| select char(65, 66 using 'latin1')
< ok 0 0
> exec ^
: string
| select char(65 using binary)
< ok 0 0
> exec ^
: string
| select char(65 using `utf8mb4`)
< ok 0 0
> exec ^
: string
| select char(65)
< ok 0 0
> exec ^
: string
| select char()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select char(using utf8mb4)
< error
| 1064 42000 …near 'utf8mb4)' at line 1
> exec ^
: string
| select char(65 using)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select char(65 using utf8mb4 collate utf8mb4_bin)
< error
| 1064 42000 …near 'collate utf8mb4_bin)' at line 1
> exec ^
: string
| select char(65 using utf8mb4, 66)
< error
| 1064 42000 …near ', 66)' at line 1
> exec ^
: string
| select char(65 using utf8mb4) over ()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| select `char`(65 using utf8mb4)
< error
| 1064 42000 …near 'using utf8mb4)' at line 1
