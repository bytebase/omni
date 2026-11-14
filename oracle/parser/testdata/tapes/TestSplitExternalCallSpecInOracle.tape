-- pin: default gvenzl/oracle-free:23.26.3-slim-faststart sha256:f5ff19033860d662c821cb04eb10483fa94f14f78eae252d054291ea07028093 23.26.3.0.0

== TestSplitExternalCallSpecInOracle default#1
> exec
| CREATE USER omni_split IDENTIFIED BY omni_test QUOTA UNLIMITED ON users
< ok 0 0
> exec
| GRANT CREATE SESSION, CREATE TABLE, CREATE VIEW, CREATE SEQUENCE, CREATE SYNONYM, CREATE TYPE, CREATE PROCEDURE, CREATE TRIGGER TO omni_split
< ok 0 0

== TestSplitExternalCallSpecInOracle default#2
> exec
| CREATE OR REPLACE PROCEDURE omni_ext_X1 AS EXTERNAL LANGUAGE C NAME "c_p" LIBRARY omni_nolib_X1;
< ok 0 0
> exec
|
| CREATE OR REPLACE PROCEDURE omni_after_X1 AS BEGIN NULL; END;
< ok 0 0
> query
| SELECT COUNT(*) FROM user_errors WHERE name = :1 AND text LIKE 'PLS-00103%'
: string
| OMNI_EXT_X1
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 0
> query
| SELECT status FROM user_objects WHERE object_name = :1 AND object_type = 'PROCEDURE'
: string
| OMNI_AFTER_X1
< result
< column "STATUS" "NCHAR"
< row
: string
| VALID
