-- pin: default gvenzl/oracle-free:23.26.3-slim-faststart sha256:f5ff19033860d662c821cb04eb10483fa94f14f78eae252d054291ea07028093 23.26.3.0.0

== TestQuotedNamesSurviveOracle default#1
> exec
| CREATE USER "OmniQuoted" IDENTIFIED BY "omni_test" QUOTA UNLIMITED ON USERS
< ok 0 0
> exec
| GRANT CREATE SESSION, CREATE TABLE, CREATE VIEW, CREATE SEQUENCE TO "OmniQuoted"
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "OmniQuoted"
< ok 0 0
> exec
| CREATE TABLE "Mixed Case" (
|   "CamelCase" NUMBER NOT NULL,
|   "lower" VARCHAR2(20 BYTE) NOT NULL
| )
< ok 0 0
> exec
|
|
| CREATE VIEW "Upper View" AS SELECT "CamelCase" FROM "Mixed Case"
< ok 0 0
> query
| SELECT object_type || '/' || object_name FROM all_objects
| 		WHERE owner = 'OmniQuoted' AND object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW', 'SEQUENCE', 'FUNCTION', 'PROCEDURE', 'INDEX', 'TRIGGER')
| 		ORDER BY object_type, object_name
< result
< column "OBJECT_TYPE||'/'||OBJECT_NAME" "NCHAR"
< row
: string
| TABLE/Mixed Case
< row
: string
| VIEW/Upper View
