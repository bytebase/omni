-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

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
| SELECT 1 /*! + 1 /*M! + 2 */
< error
| 1064 42000 …near '' at line 1
