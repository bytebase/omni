-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestAdjacentStringFormsOracle default#1
> exec
| PREPARE omni_parses_stmt FROM ?
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near '_utf8mb4'b'' at line 1
> exec ^
: string
| select 'a' x'62'
< error
| 1064 42000 …near 'x'62'' at line 1
> exec ^
: string
| select * from t limit 1 offset 'a' 'b'
< error
| 1064 42000 …near ''a' 'b'' at line 1
> exec ^
: string
| create table t (a varchar(10) default 'a' 'b')
< ok 0 0
> exec ^
: string
| create table t (a int comment 'x' 'y')
< error
| 1064 42000 …near ''y')' at line 1
> exec ^
: string
| create table t (a enum('x' 'y'))
< error
| 1064 42000 …near ''y'))' at line 1
> exec ^
: string
| create table t (a int) comment 'x' 'y'
< error
| 1064 42000 …near ''y'' at line 1
> exec ^
: string
| select 1 from t where 'a' like 'b' escape 'x' 'y'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select 1 from t where 'a' like 'b' escape 'x' 'y'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select 'a' like 'b' escape 'x' 'y'
< error
| 1210 HY000 Incorrect arguments to ESCAPE
> exec ^
: string
| select 'a' like 'b' escape 'x' 'y'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
