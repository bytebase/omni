-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

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
| 1054 42S22 Unknown column 'a*/' in 'SELECT'
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT 1 /*! + LENGTH(`a*/`) */
< error
| 1054 42S22 Unknown column 'a*/' in 'SELECT'
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
