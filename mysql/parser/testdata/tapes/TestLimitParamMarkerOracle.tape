-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestLimitParamMarkerOracle default#1
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| SELECT * FROM t LIMIT ?
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
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
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
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
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT * FROM t LIMIT ? OFFSET ?
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
