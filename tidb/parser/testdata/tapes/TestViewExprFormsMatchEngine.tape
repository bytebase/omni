-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestViewExprFormsMatchEngine default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| select a sounds like b from t
< error
| 1064 42000 …line 1 column 20 near "like b from t" 
> exec ^
: string
| select a not sounds like b from t
< error
| 1064 42000 …line 1 column 19 near "sounds like b from t" 
> exec ^
: string
| select a sounds from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a sounds from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
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
< error
| 1064 42000 …line 1 column 19 near "'z'" 
> exec ^
: string
| select N'x' as y
< ok 0 0
> exec ^
: string
| select N'x' collate utf8_bin
< ok 0 0
> exec ^
: string
| select N 'x' from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select N 'x' from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select N 'x' as y from t
< error
| 1064 42000 …line 1 column 15 near "as y from t" 
> exec ^
: string
| select `N`'x' from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select `N`'x' from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select N"x" from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select N"x" from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select N"x" as y from t
< error
| 1064 42000 …line 1 column 14 near "as y from t" 
> exec ^
: string
| select _utf8'x' 'y'
< ok 0 0
> exec ^
: string
| select _utf8'x' 'y' as z
< error
| 1064 42000 …line 1 column 22 near "as z" 
> exec ^
: string
| select trim(both from c) from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select trim(both from c) from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select trim(leading from c) from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select trim(leading from c) from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select trim(trailing from c) from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select trim(trailing from c) from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select trim(both 'x' from c) from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select trim(both 'x' from c) from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select trim('x' from c) from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select trim('x' from c) from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select trim(c) from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select trim(c) from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select trim(leading 'x')
< error
| 1064 42000 …line 1 column 24 near ")" 
> exec ^
: string
| select trim(both)
< error
| 1064 42000 …line 1 column 17 near ")" 
> exec ^
: string
| select trim(both from)
< error
| 1064 42000 …line 1 column 22 near ")" 
> exec ^
: string
| select trim(both from 'x' from 'y')
< error
| 1064 42000 …line 1 column 30 near "from 'y')" 
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
| select char(65 using `utf8`)
< ok 0 0
> exec ^
: string
| select char(65 using utf8mb3)
< ok 0 0
> exec ^
: string
| select char(65)
< ok 0 0
> exec ^
: string
| select char()
< error
| 1064 42000 …line 1 column 13 near ")" 
> exec ^
: string
| select char(using utf8)
< error
| 1064 42000 …line 1 column 17 near "using utf8)" 
> exec ^
: string
| select char(65 using)
< error
| 1064 42000 …line 1 column 21 near ")" 
> exec ^
: string
| select char(65 using utf8mb4 collate utf8mb4_bin)
< error
| 1064 42000 …line 1 column 36 near "collate utf8mb4_bin)" 
> exec ^
: string
| select char(65 using utf8mb4, 66)
< error
| 1064 42000 …line 1 column 29 near ", 66)" 
> exec ^
: string
| select (a, b) in ((1, 2)) from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select (a, b) in ((1, 2)) from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select (a, b) = (1, 2) from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select (a, b) = (1, 2) from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select row(a, b) < (1, 2) from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select row(a, b) < (1, 2) from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
