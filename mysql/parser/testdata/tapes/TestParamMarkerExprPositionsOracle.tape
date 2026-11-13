-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestParamMarkerExprPositionsOracle default#1
> exec
| PREPARE omni_db_option_stmt FROM ?
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
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
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
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT ROW_NUMBER() OVER (ROWS ? PRECEDING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT ROW_NUMBER() OVER (ROWS BETWEEN ? FOLLOWING AND UNBOUNDED PRECEDING) FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT ROW_NUMBER() OVER (ROWS BETWEEN ? FOLLOWING AND UNBOUNDED PRECEDING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
