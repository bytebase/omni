-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestExecCommentNesting_Oracle default#1
> exec
| SELECT 1 /*! + 1 /*! + 2 /* x */ + 3 */
< ok 0 0
> exec
| SELECT 1 /*! + 1 /*! + 2 /* x */ */
< ok 0 0
> exec
| SELECT 1 /*! + 1 /*! + 2 */
< ok 0 0
> exec
| SELECT 1 /*T! + 1 /*T! + 2 /* x */ + 3 */
< ok 0 0
> exec
| SELECT 1 /*! + 1 /*T![clustered_index] + 2 */
< ok 0 0
> exec
| SELECT 1 /*! + 1 /*T![no_such_feature] + 2 */ + 3 */
< ok 0 0
> exec
| SELECT 1 /*T![clustered_index] + 1 /*T![clustered_index] + 2 */
< ok 0 0
> exec
| SELECT 1 /*! + 1 /*!50000 + 2 */
< ok 0 0
> exec
| SELECT 1 /*! + 1 /*T! + 2 */
< ok 0 0
> exec
| SELECT 1 /*! + 1 /*! + 2 '*/' */
< ok 0 0
> exec
| SELECT 1 /*! + 1 /*! /*! + 2 */
< ok 0 0
> exec
| SELECT 1 /*! + 1 */ /*! + 2 */
< ok 0 0
> exec
| SELECT 1 /*! + 1 /*! + 2 */ */
< error
| 1064 42000 …line 1 column 30 near "/" 
> exec
| SELECT 1 /*! + 1 /*! + 2 */ + 3 */
< error
| 1064 42000 …line 1 column 34 near "/" 
> exec
| SELECT 1 /*!50000 + 1 /*!50000 + 2 */ */
< error
| 1064 42000 …line 1 column 40 near "/" 
> exec
| SELECT 1 /*T! + 1 /*! + 2 */ */
< error
| 1064 42000 …line 1 column 31 near "/" 
> exec
| SELECT 1 /*T![clustered_index] + 1 /*T![clustered_index] + 2 */ */
< error
| 1064 42000 …line 1 column 66 near "/" 
> exec
| SELECT 1 /*! + 1 /*T![clustered_index + 2 */ + 3 */
< error
| 1064 42000 …line 1 column 22 near "[clustered_index + 2 */ + 3 */" 
> exec
| SELECT 1 /*T![clustered_index,auto_rand] + 2 */
< ok 0 0
> exec
| SELECT 1 /*T![clustered_index]+ 2 */
< ok 0 0
> exec
| SELECT 1 /*T![no_such_feature] + 2 */ + 3
< ok 0 0
> exec
| SELECT 1 /*T![clustered_index,no_such_feature] + 2 */ + 3
< ok 0 0
> exec
| SELECT 1 /*T![Clustered_Index] + 2 */
< ok 0 0
> exec
| SELECT 1 /*T![clustered_index$] + 2 */
< ok 0 0
> exec
| SELECT 1 /*T![9] + 2 */
< ok 0 0
> exec
| SELECT 1 /*T![clustered_indexé] + 2 */
< ok 0 0
> exec
| SELECT 1 /*T![é] + 2 */
< ok 0 0
> exec
| SELECT 1 /*T![clustered_index] , 2
< ok 0 0
> exec
| /*T![clustered_index] WITH c AS (SELECT 1 AS a) SELECT a FROM c
< ok 0 0
> exec
| SELECT 1 /*T![clustered_index + 2 */
< error
| 1064 42000 …line 1 column 14 near "[clustered_index + 2 */" 
> exec
| SELECT 1 /*T![clustered_index + 2 */ + 3
< error
| 1064 42000 …line 1 column 14 near "[clustered_index + 2 */ + 3" 
> exec
| SELECT 1 /*T![ */ + 3
< error
| 1064 42000 …line 1 column 14 near "[ */ + 3" 
> exec
| SELECT 1 /*T![no_such_feature + 2 */ + 3
< error
| 1064 42000 …line 1 column 14 near "[no_such_feature + 2 */ + 3" 
> exec
| SELECT 1 /*T![clustered_index ] + 2 */
< error
| 1064 42000 …line 1 column 14 near "[clustered_index ] + 2 */" 
> exec
| SELECT 1 /*T![ clustered_index] + 2 */
< error
| 1064 42000 …line 1 column 14 near "[ clustered_index] + 2 */" 
> exec
| SELECT 1 /*T![] + 2 */
< error
| 1064 42000 …line 1 column 14 near "[] + 2 */" 
> exec
| SELECT 1 /*T![clustered_index,] + 2 */
< error
| 1064 42000 …line 1 column 14 near "[clustered_index,] + 2 */" 
> exec
| SELECT 1 /*T![clustered_index, auto_rand] + 2 */
< error
| 1064 42000 …line 1 column 14 near "[clustered_index, auto_rand] + 2 */" 
> exec
| SELECT 1 /*T![clustered-index] + 2 */
< error
| 1064 42000 …line 1 column 14 near "[clustered-index] + 2 */" 
