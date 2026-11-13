-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestLimitParamMarkerOracle default#1
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT * FROM t LIMIT ?
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t LIMIT ?
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT * FROM t LIMIT ?, ?
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t LIMIT ?, ?
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT * FROM t LIMIT ? OFFSET ?
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t LIMIT ? OFFSET ?
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
