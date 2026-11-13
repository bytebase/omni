-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestPrepareParsesControls default#1
> exec
| PREPARE omni_db_option_stmt FROM ?
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
| SELECT nosuch
< error
| 1054 42S22 Unknown column 'nosuch' in 'field list'
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
| SELECT CAST(a AS JSON ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of JSON'
> exec ^
: string
| SELECT CAST(a AS JSON ARRAY)
| ; SELECT FROM
< error
| 1235 42000 This version of MySQL doesn't yet support 'CAST-ing data to array of JSON'
> exec ^
: string
| SELECT CAST(a AS SIGNED ARRAY)
< error
| 1235 42000 This version of MySQL doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec ^
: string
| SELECT CAST(a AS SIGNED ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS FLOAT(54))
< error
| 1426 42000 Too-big precision 54 specified for 'CAST'. Maximum is 53.
> exec ^
: string
| SELECT CAST(a AS FLOAT(54))
| ; SELECT FROM
< error
| 1426 42000 Too-big precision 54 specified for 'CAST'. Maximum is 53.
> exec ^
: string
| SELECT @@nosuch
< error
| 1193 HY000 Unknown system variable 'nosuch'
