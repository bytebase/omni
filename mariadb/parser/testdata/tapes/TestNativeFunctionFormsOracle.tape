-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestNativeFunctionFormsOracle default#1
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT JSON_SET(1,1,1)
< ok 0 0
> exec ^
: string
| SELECT JSON_SET(1,1,1,1)
< error
| 1582 42000 Incorrect parameter count in the call to native function 'JSON_SET'
> exec ^
: string
| SELECT JSON_SET(1,1,1,1)
| ; SELECT FROM
< error
| 1582 42000 Incorrect parameter count in the call to native function 'JSON_SET'
> exec ^
: string
| SELECT JSON_SET(1,1,1,1,1)
< ok 0 0
> exec ^
: string
| SELECT JSON_OBJECT()
< ok 0 0
> exec ^
: string
| SELECT JSON_OBJECT(1)
< error
| 1582 42000 Incorrect parameter count in the call to native function 'JSON_OBJECT'
> exec ^
: string
| SELECT JSON_OBJECT(1)
| ; SELECT FROM
< error
| 1582 42000 Incorrect parameter count in the call to native function 'JSON_OBJECT'
> exec ^
: string
| SELECT JSON_OBJECT(1,2)
< ok 0 0
> exec ^
: string
| SELECT JSON_OBJECT(1,2,3)
< error
| 1582 42000 Incorrect parameter count in the call to native function 'JSON_OBJECT'
> exec ^
: string
| SELECT JSON_OBJECT(1,2,3)
| ; SELECT FROM
< error
| 1582 42000 Incorrect parameter count in the call to native function 'JSON_OBJECT'
> exec ^
: string
| SELECT `abs`(1,2)
< error
| 1582 42000 Incorrect parameter count in the call to native function 'abs'
> exec ^
: string
| SELECT `abs`(1,2)
| ; SELECT FROM
< error
| 1582 42000 Incorrect parameter count in the call to native function 'abs'
> exec ^
: string
| SELECT d.abs(1,2)
< error
| 1305 42000 FUNCTION d.abs does not exist
> exec ^
: string
| SELECT d.abs(1,2)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
|q "SELECT 1 AS `\xff`"
< error
| 1300 HY000 Invalid utf8mb4 character string: '\xFF'
> exec ^
: string
|q "SELECT 1 AS `\xff`\n; SELECT FROM"
< error
| 1300 HY000 Invalid utf8mb4 character string: '\xFF'
