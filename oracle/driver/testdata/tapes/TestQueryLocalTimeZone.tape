-- pin: default gvenzl/oracle-free:23.26.3-slim-faststart sha256:f5ff19033860d662c821cb04eb10483fa94f14f78eae252d054291ea07028093 23.26.3.0.0

== TestQueryLocalTimeZone default#1
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
| CREATE TABLE OMNI_APP.LTZ (ID NUMBER, AT TIMESTAMP WITH LOCAL TIME ZONE)
< ok 0 0
> exec
|
| ALTER SESSION SET TIME_ZONE = '+02:00'
< ok 0 0
> exec
|
| INSERT INTO OMNI_APP.LTZ (ID, AT) VALUES (1, TIMESTAMP '2026-03-04 05:06:07')
< ok 1 0
> exec
| ALTER SESSION SET TIME_ZONE = '+02:00'
< ok 0 0
> query
|  SELECT AT FROM OMNI_APP.LTZ
< result
< column "AT" "TimeStampLTZ_DTY"
< row
: time 2026-03-04T03:06:07Z
> query
| SELECT SESSIONTIMEZONE FROM DUAL
< result
< column "SESSIONTIMEZONE" "NCHAR"
< row
: string
| +02:00
> exec
| ALTER SESSION SET TIME_ZONE = '+00:00'
< ok 0 0
> query
|  SELECT AT FROM OMNI_APP.LTZ
< result
< column "AT" "TimeStampLTZ_DTY"
< row
: time 2026-03-04T03:06:07Z
> query
| SELECT SESSIONTIMEZONE FROM DUAL
< result
< column "SESSIONTIMEZONE" "NCHAR"
< row
: string
| +00:00
> exec
| ALTER SESSION SET TIME_ZONE = 'Europe/Berlin'
< ok 0 0
> query
|  SELECT AT FROM OMNI_APP.LTZ
< result
< column "AT" "TimeStampLTZ_DTY"
< row
: time 2026-03-04T03:06:07Z
> query
| SELECT SESSIONTIMEZONE FROM DUAL
< result
< column "SESSIONTIMEZONE" "NCHAR"
< row
: string
| Europe/Berlin
