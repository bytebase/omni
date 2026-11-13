-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestServerParsesControls default#1
> exec
| PREPARE omni_name_form_stmt FROM ?
: string
| SELECT 1 FROM nosuch
< error
| 1146 42S02 Table 'omni_t_mysql_parser.nosuch' doesn't exist
> exec ^
: string
| SELECT 1 FROM nosuch
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 1 FROM nosuch -- a comment
< error
| 1146 42S02 Table 'omni_t_mysql_parser.nosuch' doesn't exist
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
| FLUSH RELAY LOGS FOR CHANNEL 'a
| b'
< error
| 1525 HY000 Incorrect argument contains not-allowed LF value: 'a
| b'
> exec ^
: string
| FLUSH RELAY LOGS FOR CHANNEL 'a
| b'
| ; SELECT FROM
< error
| 1525 HY000 Incorrect argument contains not-allowed LF value: 'a
| b'
> exec ^
: string
| SELECT @@nosuch
< error
| 1193 HY000 Unknown system variable 'nosuch'
> exec ^
: string
| FLUSH RELAY LOGS FOR CHANNEL 'c' FOR CHANNEL
< error
| 1064 42000 …near 'FOR CHANNEL' at line 1
