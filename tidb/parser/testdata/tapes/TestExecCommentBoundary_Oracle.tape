-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestExecCommentBoundary_Oracle default#1
> exec
| /*!SELECT 1*/2
< error
| 1064 42000 …line 1 column 14 near "2" 
> exec
| /*!SELECT 'v' x*/y
< error
| 1064 42000 …line 1 column 18 near "y" 
> exec
| /*T!SELECT 1*/2
< error
| 1064 42000 …line 1 column 15 near "2" 
> exec
| SELECT /*!50000 1, */ 2
< ok 0 0
> exec
| SELECT /*T! 1, */ 2
< ok 0 0
> exec
| SELECT 'x' /*! , 'a*/'
< ok 0 0
> exec
| SELECT 1 /*!50000 + 1 -- c */
|
< ok 0 0
