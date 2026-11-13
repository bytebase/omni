-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestExecCommentQuotedDelimiters_Oracle default#1
> exec
| SELECT 1 /*! + LENGTH('/*') */
< ok 0 0
> exec
| SELECT 1 /*!50000 + LENGTH('/*') */
< ok 0 0
> exec
| SELECT 1 /*! + LENGTH('*/') */
< ok 0 0
> exec
| SELECT 1 /*! + LENGTH("/*") */
< ok 0 0
> exec
| SELECT 1 /*! + LENGTH(`a*/`) */
< error
| 1054 42S22 Unknown column 'a*/' in 'field list'
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| SELECT 1 /*! + LENGTH(`a*/`) */
< error
| 1054 42S22 Unknown column 'a*/' in 'field list'
> exec ^
: string
| SELECT 1 /*! + LENGTH(`a*/`) */
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT 1 /*! + LENGTH('a\'*/') */
< ok 0 0
> exec
| /*!50000 SELECT 'a*/' */
< ok 0 0
> exec
| SELECT 1 /*! /* inner */ + 1 */
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
> exec
| SELECT 1 /*! + 1 # c */
|
< error
| 1064 42000 …near '' at line 1
> exec
| SELECT 1 /*! + LENGTH('a--*/') -- */
|
< error
| 1064 42000 …near '' at line 1
