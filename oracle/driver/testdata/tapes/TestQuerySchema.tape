-- pin: default gvenzl/oracle-free:23.26.3-slim-faststart sha256:f5ff19033860d662c821cb04eb10483fa94f14f78eae252d054291ea07028093 23.26.3.0.0

== TestQuerySchema default#1
> exec
|
| CREATE USER OMNI_APP NO AUTHENTICATION QUOTA UNLIMITED ON USERS
< ok 0 0
> exec
|
| GRANT CREATE SESSION, CREATE TABLE TO OMNI_APP
< ok 0 0
> exec
|
| CREATE TABLE OMNI_APP.S (ID NUMBER)
< ok 0 0
> exec
|
| INSERT INTO OMNI_APP.S (ID) VALUES (1)
< ok 1 0
> query
| SELECT ID FROM S
< error
| 942 ORA-00942: table or view "SYSTEM"."S" does not exist
|
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "OMNI_APP"
< ok 0 0
> query
| SELECT ID FROM S
< result
< column "ID" "NUMBER"
< row
: string
| 1
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "OMNI_APP"
< ok 0 0
> query
| SELECT ID FROM S
< result
< column "ID" "NUMBER"
< row
: string
| 1
