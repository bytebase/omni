-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestNaturalJoinFormsOracle default#1
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT * FROM a NATURAL JOIN b
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| SELECT * FROM a NATURAL JOIN b
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT * FROM a NATURAL INNER JOIN b
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| SELECT * FROM a NATURAL INNER JOIN b
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT * FROM a NATURAL STRAIGHT_JOIN b
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| SELECT * FROM a NATURAL STRAIGHT_JOIN b
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT * FROM a NATURAL LEFT OUTER JOIN b
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| SELECT * FROM a NATURAL LEFT OUTER JOIN b
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DELETE a FROM a NATURAL INNER JOIN b
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| DELETE a FROM a NATURAL INNER JOIN b
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT * FROM a NATURAL CROSS JOIN b
< error
| 1064 42000 …near 'CROSS JOIN b' at line 1
> exec ^
: string
| SELECT * FROM a NATURAL INNER b
< error
| 1064 42000 …near 'b' at line 1
> exec ^
: string
| SELECT * FROM a NATURAL INNER STRAIGHT_JOIN b
< error
| 1064 42000 …near 'STRAIGHT_JOIN b' at line 1
> exec ^
: string
| SELECT * FROM a NATURAL INNER JOIN b ON a.x = b.x
< error
| 1064 42000 …near 'ON a.x = b.x' at line 1
