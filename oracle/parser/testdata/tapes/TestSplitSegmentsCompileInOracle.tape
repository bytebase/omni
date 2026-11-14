-- pin: default gvenzl/oracle-free:23.26.3-slim-faststart sha256:f5ff19033860d662c821cb04eb10483fa94f14f78eae252d054291ea07028093 23.26.3.0.0

== TestSplitSegmentsCompileInOracle default#1
> exec
| CREATE USER omni_split IDENTIFIED BY omni_test QUOTA UNLIMITED ON users
< ok 0 0
> exec
| GRANT CREATE SESSION, CREATE TABLE, CREATE VIEW, CREATE SEQUENCE, CREATE SYNONYM, CREATE TYPE, CREATE PROCEDURE, CREATE TRIGGER TO omni_split
< ok 0 0

== TestSplitSegmentsCompileInOracle default#2
> exec
| CREATE TABLE t (id NUMBER)
< ok 0 0
> query
| SELECT DBMS_DDL.WRAP(:1) FROM dual
: string
| CREATE OR REPLACE PROCEDURE omni_wrap_S1 AS BEGIN NULL; END;
< result
< column "DBMS_DDL.WRAP(:1)" "NCHAR"
< row
: string
| CREATE OR REPLACE PROCEDURE omni_wrap_S1 wrapped 
| a000000
| 369
| abcd
| abcd
| abcd
| abcd
| abcd
| abcd
| abcd
| abcd
| abcd
| abcd
| abcd
| abcd
| abcd
| abcd
| abcd
| 7
| 2b 61
| 4yxX7sFQFm75Mc8eGhFgdM0dNcUwg5nnm7+fMr2ywFyFctfZX5aW8lbjKN/DuHQruMAy/tJe
| pZmBCC0ruIHHLcmmpkkBdCM=
|
> exec
| CREATE OR REPLACE PROCEDURE omni_target_S1 AS BEGIN NULL; END;
< ok 0 0
> exec
|
|   CREATE OR REPLACE EDITIONABLE FUNCTION omni_jfn_S1 (p IN BLOB, q IN VARCHAR2) RETURN VARCHAR2
| AS LANGUAGE JAVA
| NAME 'Signer.check(oracle.sql.BLOB, java.lang.String) return java.lang.String';
< ok 0 0
> exec
|
|
| CREATE OR REPLACE PROCEDURE omni_jproc_S1 (a IN VARCHAR2) AS LANGUAGE JAVA NAME 'X.y(java.lang.String)';
< ok 0 0
> exec
|
| CREATE OR REPLACE PACKAGE omni_pk_S1 AS
|   FUNCTION f RETURN NUMBER AS LANGUAGE JAVA NAME 'X.f() return int';
|   PROCEDURE p;
| END;
< ok 0 0
> exec
|
| CREATE OR REPLACE PACKAGE BODY omni_pk_S1 AS
|   FUNCTION g RETURN NUMBER AS LANGUAGE JAVA NAME 'X.g() return int';
|   PROCEDURE p IS BEGIN NULL; END;
| END;
< ok 0 0
> exec
|
| CREATE OR REPLACE TRIGGER omni_trg_S1 BEFORE INSERT ON t FOR EACH ROW CALL omni_target_S1
< ok 0 0
> exec
|
| CREATE OR REPLACE PROCEDURE omni_wrap_S1 wrapped 
| a000000
| 369
| abcd
| abcd
| abcd
| abcd
| abcd
| abcd
| abcd
| abcd
| abcd
| abcd
| abcd
| abcd
| abcd
| abcd
| abcd
| 7
| 2b 61
| 4yxX7sFQFm75Mc8eGhFgdM0dNcUwg5nnm7+fMr2ywFyFctfZX5aW8lbjKN/DuHQruMAy/tJe
| pZmBCC0ruIHHLcmmpkkBdCM=
< ok 0 0
> query
|
| SELECT o.object_name, o.object_type, o.status,
|        (SELECT LISTAGG(e.text, ' / ') WITHIN GROUP (ORDER BY e.sequence)
|           FROM user_errors e WHERE e.name = o.object_name AND e.type = o.object_type)
|   FROM user_objects o
|  WHERE o.object_name LIKE '%\_' || :1 ESCAPE '\'
|  ORDER BY o.object_name, o.object_type
: string
| S1
< result
< column "OBJECT_NAME" "NCHAR"
< column "OBJECT_TYPE" "NCHAR"
< column "STATUS" "NCHAR"
< column "(SELECTLISTAGG(E.TEXT,'/')WITHINGROUP(ORDERBYE.SEQUENCE)FROMUSER_ERRORSEWHEREE.NAME=O.OBJECT_NAMEANDE.TYPE=O.OBJECT_TYPE)" "NCHAR"
< row
: string
| OMNI_JFN_S1
: string
| FUNCTION
: string
| VALID
: null
< row
: string
| OMNI_JPROC_S1
: string
| PROCEDURE
: string
| VALID
: null
< row
: string
| OMNI_PK_S1
: string
| PACKAGE
: string
| VALID
: null
< row
: string
| OMNI_PK_S1
: string
| PACKAGE BODY
: string
| VALID
: null
< row
: string
| OMNI_TARGET_S1
: string
| PROCEDURE
: string
| VALID
: null
< row
: string
| OMNI_TRG_S1
: string
| TRIGGER
: string
| VALID
: null
< row
: string
| OMNI_WRAP_S1
: string
| PROCEDURE
: string
| VALID
: null
