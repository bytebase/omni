-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestPrepareParsesControls default#1
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT 1 FROM nosuch
< error
| 1146 42S02 Table 'test.nosuch' doesn't exist
> exec ^
: string
| SELECT 1 FROM nosuch
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT nosuch
< error
| 1054 42S22 Unknown column 'nosuch' in 'SELECT'
> exec ^
: string
| SELECT nosuch
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 1 FROM nosuch -- a comment
< error
| 1146 42S02 Table 'test.nosuch' doesn't exist
> exec ^
: string
| SELECT 1 FROM nosuch -- a comment
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELEKT 1
< error
| 1064 42000 …near 'SELEKT 1' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(0))
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a VECTOR(0))
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a VECTOR(16384))
< error
| 1074 42000 Column length too big for column 'a' (max = 16383); use BLOB or TEXT instead
> exec ^
: string
| CREATE TABLE t (a VECTOR(16384))
| ; SELECT FROM
< error
| 1074 42000 Column length too big for column 'a' (max = 16383); use BLOB or TEXT instead
> exec ^
: string
| SELECT CAST(1 AS POINT)
< error
| 4162 HY000 Operator does not exist: 'CAST(expr AS point)'
> exec ^
: string
| SELECT CAST(1 AS POINT)
| ; SELECT FROM
< error
| 4162 HY000 Operator does not exist: 'CAST(expr AS point)'
> exec ^
: string
| SELECT @@nosuch
< error
| 1193 HY000 Unknown system variable 'nosuch'
