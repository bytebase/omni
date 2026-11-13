-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestAdminIsNotARoutineStatementOnTiDB default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| CREATE PROCEDURE p() ADMIN SHOW DDL
< error
| 1064 42000 …line 1 column 26 near "ADMIN SHOW DDL" 
> exec ^
: string
| CREATE PROCEDURE p() ADMIN FLUSH BINDINGS
< error
| 1064 42000 …line 1 column 26 near "ADMIN FLUSH BINDINGS" 
> exec ^
: string
| CREATE PROCEDURE p() BEGIN ADMIN SHOW DDL; END
< error
| 1064 42000 …line 1 column 32 near "ADMIN SHOW DDL; END" 
> exec ^
: string
| CREATE PROCEDURE p() BEGIN ADMIN CHECK TABLE t; END
< error
| 1064 42000 …line 1 column 32 near "ADMIN CHECK TABLE t; END" 
> exec ^
: string
| CREATE PROCEDURE p() BEGIN WHILE 1 DO ADMIN SHOW DDL; END WHILE; END
< error
| 1064 42000 …line 1 column 43 near "ADMIN SHOW DDL; END WHILE; END" 
> exec ^
: string
| CREATE PROCEDURE p() BEGIN IF 1 THEN ADMIN SHOW DDL; END IF; END
< error
| 1064 42000 …line 1 column 42 near "ADMIN SHOW DDL; END IF; END" 
