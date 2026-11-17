-- pin: default gvenzl/oracle-free:23.26.3-slim-faststart sha256:f5ff19033860d662c821cb04eb10483fa94f14f78eae252d054291ea07028093 23.26.3.0.0

== TestOwnsPrivateDatabaseLink default#1
> exec
|
| CREATE USER OMNI_APP NO AUTHENTICATION QUOTA UNLIMITED ON USERS
< ok 0 0
> query
| SELECT COUNT(*) FROM SYS.USER_DB_LINKS
< result
< column "COUNT(*)" "NUMBER" "float64" precision 38 scale 255
< row
: string
| 0
> exec
| CREATE DATABASE LINK omni_link CONNECT TO omni IDENTIFIED BY omni USING 'nosuch.example.com:1521/XE'
< ok 0 0
> query
| SELECT COUNT(*) FROM SYS.USER_DB_LINKS
< result
< column "COUNT(*)" "NUMBER" "float64" precision 38 scale 255
< row
: string
| 1
> exec
| DROP DATABASE LINK omni_link
< ok 0 0
