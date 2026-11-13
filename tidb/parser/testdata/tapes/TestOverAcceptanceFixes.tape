-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestOverAcceptanceFixes default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| CREATE TABLE COUNT(a bigint)
< error
| 1064 42000 …line 1 column 18 near "COUNT(a bigint)" 
> exec ^
: string
| CREATE TABLE BIT_AND(a bigint)
< error
| 1064 42000 …line 1 column 20 near "BIT_AND(a bigint)" 
> exec ^
: string
| CREATE TABLE NOW(a bigint)
< error
| 1064 42000 …line 1 column 16 near "NOW(a bigint)" 
> exec ^
: string
| CREATE TABLE SUM(a int)
< error
| 1064 42000 …line 1 column 16 near "SUM(a int)" 
> exec ^
: string
| CREATE TABLE GROUP_CONCAT(a int)
< error
| 1064 42000 …line 1 column 25 near "GROUP_CONCAT(a int)" 
> exec ^
: string
| CREATE TABLE COUNT (a bigint)
< ok 0 0
> exec ^
: string
| CREATE TABLE `COUNT` (a bigint)
< ok 0 0
> exec ^
: string
| SELECT COUNT(*) FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT COUNT(*) FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM COUNT
< error
| 1146 42S02 Table 'test.COUNT' doesn't exist
> exec ^
: string
| SELECT * FROM COUNT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (COUNT int)
< ok 0 0
> exec ^
: string
| ALTER USER USER() IDENTIFIED BY PASSWORD '*B50FBDB37F1256824274912F2A1CE648082C3F1F'
< error
| 1064 42000 …line 1 column 40 near "PASSWORD '*B50FBDB37F1256824274912F2A1CE648082C3F1F'" 
> exec ^
: string
| ALTER USER USER() IDENTIFIED BY 'plain'
< ok 0 0
> exec ^
: string
| CREATE USER u IDENTIFIED BY PASSWORD '*hash'
< ok 0 0
> exec ^
: string
| ALTER USER u IDENTIFIED BY PASSWORD '*hash'
< ok 0 0
> exec ^
: string
| select * from t limit 18446744073709551616 offset 3
< error
| 1064 42000 …line 1 column 42 near "18446744073709551616 offset 3" 
> exec ^
: string
| select * from t limit 18446744073709551615 offset 3
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t limit 18446744073709551615 offset 3
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select 18446744073709551616
< ok 0 0
> exec ^
: string
| insert into tt1 values(4556414e723532)
< error
| 1367 22007 Illegal double '4556414e723532' value found during parsing
> exec ^
: string
| insert into tt1 values(4556414e723532)
| ; SELECT FROM
< error
| 1367 22007 Illegal double '4556414e723532' value found during parsing
> exec ^
: string
| select 1e400
< error
| 1367 22007 Illegal double '1e400' value found during parsing
> exec ^
: string
| select 1e400
| ; SELECT FROM
< error
| 1367 22007 Illegal double '1e400' value found during parsing
> exec ^
: string
| select 1e-400
< ok 0 0
> exec ^
: string
| select 1.5e10
< ok 0 0
