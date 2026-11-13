-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestViewExprFormsOracle default#1
> exec
| PREPARE omni_db_option_stmt FROM ?
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
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select N 'x' from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select N 'x' as y from t
< error
| 1064 42000 …near 'as y from t' at line 1
> exec ^
: string
| select `N`'x' from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select `N`'x' from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select N"x" from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select N"x" from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select N"x" as y from t
< error
| 1064 42000 …near 'as y from t' at line 1
> exec ^
: string
| select trim(both from c) from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select trim(both from c) from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select trim(leading from c) from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select trim(leading from c) from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select trim(trailing from c) from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select trim(trailing from c) from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select trim(both 'x' from c) from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select trim(both 'x' from c) from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select trim('x' from c) from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select trim('x' from c) from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select trim(c) from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select trim(c) from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'using utf8mb4)' at line 1
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
| 1064 42000 …near 'over ()' at line 1
> exec ^
: string
| select `char`(65 using utf8mb4)
< error
| 1064 42000 …near 'using utf8mb4)' at line 1
