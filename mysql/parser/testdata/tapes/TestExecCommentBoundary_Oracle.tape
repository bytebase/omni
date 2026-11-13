-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestExecCommentBoundary_Oracle default#1
> exec
| /*!SELECT 1*/2
< error
| 1064 42000 …near '2' at line 1
> exec
| /*!SELECT 'v' x*/y
< error
| 1064 42000 …near 'y' at line 1
> exec
| SELECT 'x' /*! , 'a*/'
< error
| 1064 42000 …near '' at line 1
> exec
| SELECT 1 /*!50000 + 1 -- c */
|
< error
| 1064 42000 …near '' at line 1
> exec
| /*!SELECT 'a'*/'b'
< ok 0 0
> exec
| SELECT /*!50000 1, */ 2
< ok 0 0
> query
| /*!SELECT 'a'*/'b'
< result
< column "a" "VARCHAR"
< row
: bytes
| ab
