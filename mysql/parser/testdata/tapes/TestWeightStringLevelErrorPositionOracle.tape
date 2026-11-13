-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestWeightStringLevelErrorPositionOracle default#1
> exec
| DO WEIGHT_STRING(x LEVEL 1)
< error
| 1064 42000 …near 'LEVEL 1)' at line 1
> exec
| DO WEIGHT_STRING(x LEVEL 1 DESC REVERSE)
< error
| 1064 42000 …near 'LEVEL 1 DESC REVERSE)' at line 1
> exec
| DO WEIGHT_STRING(x LEVEL 1, 2)
< error
| 1064 42000 …near 'LEVEL 1, 2)' at line 1
> exec
| DO WEIGHT_STRING(x LEVEL 1 - 2)
< error
| 1064 42000 …near 'LEVEL 1 - 2)' at line 1
> exec
| DO WEIGHT_STRING(x AS CHAR(4) LEVEL 1)
< error
| 1064 42000 …near 'LEVEL 1)' at line 1
> exec
| DO WEIGHT_STRING(x AS BINARY(4) LEVEL 1)
< error
| 1064 42000 …near 'LEVEL 1)' at line 1
