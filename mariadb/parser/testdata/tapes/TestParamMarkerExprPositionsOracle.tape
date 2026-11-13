-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestParamMarkerExprPositionsOracle default#1
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT ?
< ok 0 0
> exec ^
: string
| SELECT 1 + ?
< ok 0 0
> exec ^
: string
| SELECT LEAD(a, ?) OVER () FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT LEAD(a, ?) OVER () FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT ROW_NUMBER() OVER (ROWS ? PRECEDING) FROM t
< error
| 1064 42000 …near '? PRECEDING) FROM t' at line 1
> exec ^
: string
| SELECT ROW_NUMBER() OVER (ROWS BETWEEN ? FOLLOWING AND UNBOUNDED PRECEDING) FROM t
< error
| 1064 42000 …near '? FOLLOWING AND UNBOUNDED PRECEDING) FROM t' at line 1
