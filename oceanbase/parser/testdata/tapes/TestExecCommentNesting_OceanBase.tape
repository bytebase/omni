-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestExecCommentNesting_OceanBase default#1
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
| 1064 42000 …near '*/' at line 1
> exec
| SELECT 1 /*M! + 1 /*! + 2 */ */
< error
| 1064 42000 …near '/' at line 1
> exec
| SELECT 1 /*! + 1 /*M! + 2 */
< error
| 1064 42000 …near 'SELECT 1 /*! + 1 /*M! + 2 */' at line 2
> exec
| SELECT 1 /*!99999 + 1 /*! + 2 */
< ok 0 0
> exec
| SELECT 1 /*!99999 + 1 /*! + 2 */ */
< error
| 1064 42000 …near '/' at line 1
