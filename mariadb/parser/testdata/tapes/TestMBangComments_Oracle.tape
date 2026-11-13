-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestMBangComments_Oracle default#1
> exec
| /*M!999999\- enable the sandbox mode */ SELECT 1
< ok 0 0
> exec
| /*M!999999 + 1 */ SELECT 1
< ok 0 0
> exec
| SELECT 1 /*M! + 1 */
< ok 0 0
> exec
| SELECT 1 /*M!100100 + 1 */
< ok 0 0
> exec
| SELECT 1 /*M!+ 1 */
< ok 0 0
> exec
| SELECT 1 /*m! + 1 */
< ok 0 0
> exec
| SELECT 1 /*M! + 1 /*M! + 2 */
< ok 0 0
> exec
| SELECT 1 /*M! + 1 /*! + 2 */
< ok 0 0
> exec
| SELECT 1 /*! + 1 /*M! + 2 */
< ok 0 0
> exec
| SELECT 1 /*M! + 1 '*/' */
< ok 0 0
> exec
| /*M! SELECT 1 */
< ok 0 0
> exec
| SELECT 1 /*M! + 1 /*M! + 2 */ */
< error
| 1064 42000 …near '/' at line 1
> exec
| SELECT 1 /*M! + 1
< error
| 1064 42000 …near '' at line 1
> exec
| SELECT 1 /*M! , */
< error
| 1064 42000 …near '' at line 1
> exec
| SELECT 1 /*M!999999 + 1 /*! + 2 */ */
< ok 0 0
> exec
| SELECT 1 /*!999999 + 1 /*! + 2 */ */ + 3
< ok 0 0
> exec
| SELECT 1 /*M!999999 + LENGTH('*/') */
< error
| 1064 42000 …near '') */' at line 1
> exec
| SELECT 1 /*!999999 + LENGTH('*/') */
< error
| 1064 42000 …near '') */' at line 1
