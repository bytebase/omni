-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

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
| 1064 42000 …near ''a*/'' at line 2
> exec
| SELECT 1 /*!50000 + 1 -- c */
|
< error
| 1064 42000 …near 'SELECT 1 /*!50000 + 1 -- c */' at line 2
> exec
| /*!SELECT 'a'*/'b'
< ok 0 0
> exec
| SELECT /*!50000 1, */ 2
< ok 0 0
> query
| /*!SELECT 'a'*/'b'
< result
< column "'a'*/'b'" "VARCHAR"
< row
: bytes
| ab
