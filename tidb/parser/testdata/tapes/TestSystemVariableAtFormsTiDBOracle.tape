-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestSystemVariableAtFormsTiDBOracle default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| select @@@tidb_index_lookup_pushdown_policy
< error
| 1064 42000 …line 1 column 9 near "@@@tidb_index_lookup_pushdown_policy" 
> exec ^
: string
| select @@@autocommit
< error
| 1064 42000 …line 1 column 9 near "@@@autocommit" 
> exec ^
: string
| select @@@a
< error
| 1064 42000 …line 1 column 9 near "@@@a" 
> exec ^
: string
| select @@@
< error
| 1064 42000 …line 1 column 9 near "@@@" 
> exec ^
: string
| select @@@a.b
< error
| 1064 42000 …line 1 column 9 near "@@@a.b" 
> exec ^
: string
| select @@@`autocommit`
< error
| 1064 42000 …line 1 column 9 near "@@@`autocommit`" 
> exec ^
: string
| select @@@'autocommit'
< error
| 1064 42000 …line 1 column 9 near "@@@'autocommit'" 
> exec ^
: string
| select @@@a from t
< error
| 1064 42000 …line 1 column 9 near "@@@a from t" 
> exec ^
: string
| select 1 + @@@a
< error
| 1064 42000 …line 1 column 13 near "@@@a" 
> exec ^
: string
| select @@autocommit, @@@a
< error
| 1064 42000 …line 1 column 23 near "@@@a" 
> exec ^
: string
| select @@global.@autocommit
< error
| 1064 42000 …line 1 column 16 near "@@global.@autocommit" 
> exec ^
: string
| select @@session.@a
< error
| 1064 42000 …line 1 column 17 near "@@session.@a" 
> exec ^
: string
| select @@local.@a
< error
| 1064 42000 …line 1 column 15 near "@@local.@a" 
> exec ^
: string
| select @@global.@`a`
< error
| 1064 42000 …line 1 column 16 near "@@global.@`a`" 
> exec ^
: string
| select @@@global.autocommit
< error
| 1064 42000 …line 1 column 9 near "@@@global.autocommit" 
> exec ^
: string
| select @@global.@@autocommit
< error
| 1064 42000 …line 1 column 16 near "@@global.@@autocommit" 
> exec ^
: string
| set @@@autocommit = 1
< error
| 1064 42000 …line 1 column 6 near "@@@autocommit = 1" 
> exec ^
: string
| set @@global.@autocommit = 1
< error
| 1064 42000 …line 1 column 13 near "@@global.@autocommit = 1" 
> exec ^
: string
| set @@session.@a = 1
< error
| 1064 42000 …line 1 column 14 near "@@session.@a = 1" 
> exec ^
: string
| set @@@a := 1
< error
| 1064 42000 …line 1 column 6 near "@@@a := 1" 
> exec ^
: string
| select @@autocommit
< ok 0 0
> exec ^
: string
| select @@`autocommit`
< ok 0 0
> exec ^
: string
| select @@'autocommit'
< ok 0 0
> exec ^
: string
| select @@global.autocommit
< ok 0 0
> exec ^
: string
| select @@global.`autocommit`
< ok 0 0
> exec ^
: string
| select @@global.'autocommit'
< ok 0 0
> exec ^
: string
| set @@autocommit = 1
< ok 0 0
> exec ^
: string
| select @a
< ok 0 0
> exec ^
: string
| select @@autocommit, @a
< ok 0 0
> exec ^
: string
| show variables like '@@@a'
< ok 0 0
