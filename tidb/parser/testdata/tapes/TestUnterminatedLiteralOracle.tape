-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestUnterminatedLiteralOracle default#1
> exec
| SELECT 'oops;
< error
| 1064 42000 …line 1 column 13 near "'oops;" 
> exec
| SELECT "oops;
< error
| 1064 42000 …line 1 column 13 near ""oops;" 
> exec
| SELECT _utf8'oops;
< error
| 1064 42000 …line 1 column 18 near "'oops;" 
> exec
| SELECT `oops;
< error
| 1064 42000 …line 1 column 13 near "`oops;" 
> exec
| SELECT @`oops;
< error
| 1064 42000 …line 1 column 14 near "@`oops;" 
> exec
| SELECT 1 /* oops;
< error
| 1064 42000 …near '/* oops;' at line 1
> exec
| /* oops
< error
| 1064 42000 …near '/* oops' at line 1
> exec
| SELECT 1 /*T![nonexistent_feature] oops
< error
| 1064 42000 …near '/*T![nonexistent_feature] oops' at line 1
> exec
| SELECT 1 /*T![oops
< error
| 1064 42000 …line 1 column 14 near "[oops" 
> exec
| SELECT 'oops';
< ok 0 0
> exec
| SELECT "oops";
< ok 0 0
> exec
| SELECT _utf8'oops';
< ok 0 0
> exec
| SELECT `oops`;
< error
| 1054 42S22 Unknown column 'oops' in 'field list'
> exec
| SELECT @`oops`;
< ok 0 0
> exec
| SELECT 1 /* oops */;
< ok 0 0
> exec
| /* oops */
< ok 0 0
> exec
| SELECT 1 /*!40000 + 1 */
< ok 0 0
> exec
| SELECT 1 /*T! + 1 */
< ok 0 0
> exec
| SELECT 1 /*T![clustered_index] + 1 */
< ok 0 0
> exec
| SELECT 1 /*T![nonexistent_feature] not sql at all */
< ok 0 0
> exec
| SELECT 1 /*!40000 oops
< ok 0 0
> exec
| SELECT 1 /*T! oops
< ok 0 0
> exec
| SELECT 1 /*T![clustered_index] oops
< ok 0 0
> exec
| /*T! SELECT 1
< ok 0 0
> exec
| SELECT /*T! 1
< ok 0 0
> exec
| SELECT /*T![ttl] 1
< ok 0 0
> exec
| SELECT /*T![auto_rand] col
< error
| 1054 42S22 Unknown column 'col' in 'field list'
> exec
| SELECT 1 /*! 
< ok 0 0
