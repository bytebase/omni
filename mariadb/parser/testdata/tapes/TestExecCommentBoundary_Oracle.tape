-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

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
| SELECT /*!50000 1, */ 2
< ok 0 0
> exec
| SELECT 'x' /*! , 'a*/'
< error
| 1064 42000 …near '' at line 1
> exec
| SELECT 1 /*!50000 + 1 -- c */
|
< error
| 1064 42000 …near '' at line 1
