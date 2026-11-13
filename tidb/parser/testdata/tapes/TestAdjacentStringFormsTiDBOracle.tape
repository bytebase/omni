-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestAdjacentStringFormsTiDBOracle default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| select 'a' 'b'
< ok 0 0
> exec ^
: string
| select 1 from t where a = 'a' 'b'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select 1 from t where a = 'a' 'b'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select ('a' 'b')
< ok 0 0
> exec ^
: string
| select 'a' 'b' 'c'
< ok 0 0
> exec ^
: string
| select 'a'
| 'b'
< ok 0 0
> exec ^
: string
| select 'a' /* c */ 'b'
< ok 0 0
> exec ^
: string
| select 'a' "b"
< ok 0 0
> exec ^
: string
| select '' 'b'
< ok 0 0
> exec ^
: string
| select 'a' 'b' as c
< ok 0 0
> exec ^
: string
| select 'a' 'b' c
< ok 0 0
> exec ^
: string
| select concat('a' 'b', 'c')
< ok 0 0
> exec ^
: string
| insert into t values ('a' 'b')
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values ('a' 'b')
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select _utf8mb4'a' 'b'
< ok 0 0
> exec ^
: string
| select N'a' 'b'
< ok 0 0
> exec ^
: string
| select 'a' _utf8mb4'b'
< error
| 1064 42000 …line 1 column 19 near "_utf8mb4'b'" 
> exec ^
: string
| select 'a' x'62'
< error
| 1064 42000 …line 1 column 16 near "x'62'" 
> exec ^
: string
| select * from t limit 1 offset 'a' 'b'
< error
| 1064 42000 …line 1 column 34 near "'a' 'b'" 
> exec ^
: string
| create table t (a varchar(10) default 'a' 'b')
< ok 0 0
> exec ^
: string
| create table t (a int comment 'x' 'y')
< error
| 1064 42000 …line 1 column 37 near "'y')" 
> exec ^
: string
| create table t (a enum('x' 'y'))
< error
| 1064 42000 …line 1 column 30 near "'y'))" 
> exec ^
: string
| create table t (a int) comment 'x' 'y'
< error
| 1064 42000 …line 1 column 38 near "'y'" 
> exec ^
: string
| select 1 from t where 'a' like 'b' escape 'x' 'y'
< error
| 1064 42000 …line 1 column 49 near "'y'" 
> exec ^
: string
| select 1 from t where 'a' ilike 'b' escape 'x' 'y'
< error
| 1064 42000 …line 1 column 50 near "'y'" 
> exec ^
: string
| select 'a' like 'b' escape 'x' 'y'
< ok 0 0
