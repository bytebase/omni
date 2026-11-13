-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestOptimizeTableOrTables default#1
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| OPTIMIZE TABLE t1
< ok 0 0
> exec ^
: string
| OPTIMIZE TABLES t1, t2
< ok 0 0
> exec ^
: string
| OPTIMIZE NO_WRITE_TO_BINLOG TABLES t1
< ok 0 0
> exec ^
: string
| OPTIMIZE LOCAL TABLE t1
< ok 0 0
> exec ^
: string
| OPTIMIZE t1
< error
| 1064 42000 …near 't1' at line 1
> exec ^
: string
| OPTIMIZE
< error
| 1064 42000 …near '' at line 1
