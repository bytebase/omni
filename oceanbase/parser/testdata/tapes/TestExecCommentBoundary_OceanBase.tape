-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestExecCommentBoundary_OceanBase default#1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| /*!SELECT 1*/2
< error
| 1064 42000 …near '2' at line 1
> exec ^
: string
| /*!SELECT 'v' x*/y
< error
| 1064 42000 …near 'y' at line 1
> exec ^
: string
| SELECT 'x' /*! , 'a*/'
< error
| 1064 42000 …near ''a*/'' at line 2
> exec ^
: string
| SELECT 1 /*!50000 + 1 -- c */
|
< error
| 1064 42000 …near 'SELECT 1 /*!50000 + 1 -- c */' at line 2
> exec ^
: string
| /*!SELECT 'a'*/'b'
< ok 0 0
> exec ^
: string
| SELECT /*!50000 1, */ 2
< ok 0 0
> exec ^
: string
| SELECT 1 /*! + LENGTH('/*') */
< ok 0 0
> exec ^
: string
| SELECT 1 /*!50000 + LENGTH('/*') */
< ok 0 0
> exec ^
: string
| SELECT 1 /*! + LENGTH('*/') */
< ok 0 0
> exec ^
: string
| SELECT 1 /*! + LENGTH("/*") */
< ok 0 0
> exec ^
: string
| SELECT 1 /*! + LENGTH(`a*/`) */
< error
| 1054 42S22 Unknown column 'a*/' in 'field list'
> exec ^
: string
| SELECT 1 /*! + LENGTH('a\'*/') */
< ok 0 0
> exec ^
: string
| /*!50000 SELECT 'a*/' */
< ok 0 0
> exec ^
: string
| SELECT 1 /*! /* inner */ + 1 */
< ok 0 0
> exec ^
: string
| SELECT 'x' /*! , 'a*/'
< error
| 1064 42000 …near ''a*/'' at line 2
> exec ^
: string
| SELECT 1 /*!50000 + 1 -- c */
|
< error
| 1064 42000 …near 'SELECT 1 /*!50000 + 1 -- c */' at line 2
> exec ^
: string
| SELECT 1 /*! + 1 # c */
|
< error
| 1064 42000 …near 'SELECT 1 /*! + 1 # c */' at line 2
> exec ^
: string
| SELECT 1 /*! + LENGTH('a--*/') -- */
|
< error
| 1064 42000 …near ''a--*/') -- */' at line 2
> exec ^
: string
| SELECT 'x' /*! , 'a*/'
< error
| 1064 42000 …near ''a*/'' at line 2
> exec ^
: string
| SELECT 1 /*!50000 + 1 -- c */
|
< error
| 1064 42000 …near 'SELECT 1 /*!50000 + 1 -- c */' at line 2
> exec ^
: string
| SELECT 1 /*! + 1 # c */
|
< error
| 1064 42000 …near 'SELECT 1 /*! + 1 # c */' at line 2
> exec ^
: string
| SELECT 1 /*! + LENGTH('a--*/') -- */
|
< error
| 1064 42000 …near ''a--*/') -- */' at line 2
