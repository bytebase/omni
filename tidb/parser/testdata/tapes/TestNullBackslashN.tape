-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestNullBackslashN default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| select \N
< ok 0 0
> exec ^
: string
| select \N from test
< error
| 1146 42S02 Table 'test.test' doesn't exist
> exec ^
: string
| select \N from test
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select (\N) from test
< error
| 1146 42S02 Table 'test.test' doesn't exist
> exec ^
: string
| select (\N) from test
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select \N + 1
< ok 0 0
> exec ^
: string
| select \N, 1 from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select \N, 1 from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select \Nx
< ok 0 0
> exec ^
: string
| insert into t values (\N)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values (\N)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select \n
< error
| 1064 42000 …line 1 column 8 near "\n" 
> exec ^
: string
| select \M
< error
| 1064 42000 …line 1 column 8 near "\M" 
