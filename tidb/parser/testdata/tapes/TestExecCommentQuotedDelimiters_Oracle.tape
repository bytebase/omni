-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

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
< ok 0 0
> exec
| SELECT 1 /*!50000 + 1 -- c */
|
< ok 0 0
> exec
| SELECT 1 /*! + 1 # c */
|
< ok 0 0
> exec
| SELECT 1 /*! + LENGTH('a--*/') -- */
|
< ok 0 0
