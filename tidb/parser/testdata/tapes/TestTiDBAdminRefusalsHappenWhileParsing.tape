-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestTiDBAdminRefusalsHappenWhileParsing default#1
> exec
| SET @omni_admin = ?
: string
| recommend index apply 9223372036854775808
< ok 0 0
> exec
| PREPARE omni_admin_stmt FROM @omni_admin
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec
| SET @omni_admin = ?
: string
| recommend index apply 9223372036854775808
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_admin_stmt FROM @omni_admin
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec
| SET @omni_admin = ?
: string
| recommend index apply 18446744073709551615
< ok 0 0
> exec
| PREPARE omni_admin_stmt FROM @omni_admin
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec
| SET @omni_admin = ?
: string
| recommend index apply 18446744073709551615
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_admin_stmt FROM @omni_admin
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec
| SET @omni_admin = ?
: string
| recommend index ignore 18446744073709551615
< ok 0 0
> exec
| PREPARE omni_admin_stmt FROM @omni_admin
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec
| SET @omni_admin = ?
: string
| recommend index ignore 18446744073709551615
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_admin_stmt FROM @omni_admin
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec
| SET @omni_admin = ?
: string
| RECOMMEND INDEX APPLY 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_admin_stmt FROM @omni_admin
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec
| SET @omni_admin = ?
: string
| RECOMMEND INDEX IGNORE 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_admin_stmt FROM @omni_admin
< error
| 1064 42000 …line 2 column 14 near "FROM" 
