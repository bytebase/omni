-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOptimizeTableOrTables default#1
> exec
| OPTIMIZE TABLE t1
< ok 0 0
> exec
| OPTIMIZE TABLES t1, t2
< ok 0 0
> exec
| OPTIMIZE NO_WRITE_TO_BINLOG TABLES t1
< ok 0 0
> exec
| OPTIMIZE LOCAL TABLE t1
< ok 0 0
> exec
| OPTIMIZE t1
< error
| 1064 42000 …near 't1' at line 1
> exec
| OPTIMIZE
< error
| 1064 42000 …near '' at line 1
