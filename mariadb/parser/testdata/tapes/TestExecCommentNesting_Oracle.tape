-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

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
| SELECT 1 /*!50000 + 1 /*!50100 + 2 */
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
| SELECT 1 /*! + 1 /* c */ + 2 */
< ok 0 0
> exec
| /*!50000 /*!50100 SELECT 1 */
< ok 0 0
> exec
| SELECT 1 /*! + 1 /*! + 2 */ */
< error
| 1064 42000 …near '/' at line 1
> exec
| SELECT 1 /*! + 1 /*! + 2 */ + 3 */
< error
| 1064 42000 …near '/' at line 1
> exec
| SELECT 1 /*!50000 + 1 /*!50000 + 2 */ */
< error
| 1064 42000 …near '/' at line 1
> exec
| /*!50000 /*!50100 SELECT 1 */ */
< error
| 1064 42000 …near '/' at line 1
> exec
| SELECT 1 /*M! + 1 /*! + 2 */ */
< error
| 1064 42000 …near '/' at line 1
> exec
| SELECT 1 /*!80023 + 1 /* c */ + 2 */ + 3
< ok 0 0
> exec
| SELECT 1 /*!80023 + 1 /*!80023 + 2 */ */ + 3
< ok 0 0
> exec
| SELECT 1 /*!80023 + 1 /*! + 2 */ */ + 3
< ok 0 0
> exec
| SELECT 1 /*!80023 'x */ + 3
< ok 0 0
> exec
| SELECT 1 /*! + 1 /*!80023 + 2 */ */
< ok 0 0
> exec
| SELECT 1 /*! + 1 /*!80023 + 2 /* c */ */ */ + 3
< ok 0 0
> exec
| SELECT 1 /*!50700 + 1 */ + 3
< ok 0 0
> exec
| SELECT 1 /*!080023 + 1 */ + 3
< ok 0 0
> exec
| SELECT 1 /*! + 1 /*!80023 + 2 */
< error
| 1064 42000 …near '' at line 1
> exec
| SELECT 1 /*!80023 + 1
< error
| 1064 42000 …near '/*!80023 + 1' at line 1
> exec
| SELECT 1 /*!80023 + 1 /*
< error
| 1064 42000 …near '/*!80023 + 1 (*' at line 1
> exec
| SELECT 1 /*!80023 + LENGTH('*/') */
< error
| 1064 42000 …near '') */' at line 1
> exec
| SELECT 1 /*! + 1 /*M! + 2 */
< ok 0 0
> exec
| SELECT 1 /*M! + 1 /*M! + 2 */
< ok 0 0
> exec
| SELECT 1 /*M! + 1 /*M! + 2 */ */
< error
| 1064 42000 …near '/' at line 1
