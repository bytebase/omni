-- pin: default gvenzl/oracle-free:23.26.3-slim-faststart sha256:f5ff19033860d662c821cb04eb10483fa94f14f78eae252d054291ea07028093 23.26.3.0.0

== TestCountAffectedRows default#1
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
| CREATE TABLE OMNI_APP.BIG AS SELECT LEVEL AS ID, LEVEL AS V FROM DUAL CONNECT BY LEVEL <= 1000
< ok 1000 0
> exec
|
| CREATE TABLE OMNI_APP.MERGE_TARGET (ID NUMBER, V NUMBER)
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_1' FOR UPDATE /*+ OPT_PARAM('optimizer_dynamic_sampling' 11) */ OMNI_APP.BIG SET V = 1 WHERE ID < 10
< ok 0 0
> query
| SELECT ID, PARENT_ID, OPERATION, CARDINALITY FROM PLAN_TABLE WHERE STATEMENT_ID = 'omni_1' ORDER BY ID
< result
< column "ID" "NUMBER"
< column "PARENT_ID" "NUMBER"
< column "OPERATION" "NCHAR"
< column "CARDINALITY" "NUMBER"
< row
: string
| 0
: null
: string
| UPDATE STATEMENT
: string
| 9
< row
: string
| 1
: string
| 0
: string
| UPDATE
: null
< row
: string
| 2
: string
| 1
: string
| TABLE ACCESS
: string
| 9
> exec
| DELETE FROM PLAN_TABLE WHERE STATEMENT_ID = 'omni_1'
< ok 3 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_2' FOR MERGE /*+ OPT_PARAM('optimizer_dynamic_sampling' 11) */ INTO OMNI_APP.MERGE_TARGET t USING (SELECT ID, V FROM OMNI_APP.BIG WHERE ID < 10) s ON (t.ID = s.ID) WHEN NOT MATCHED THEN INSERT (ID, V) VALUES (s.ID, s.V)
< ok 0 0
> query
| SELECT ID, PARENT_ID, OPERATION, CARDINALITY FROM PLAN_TABLE WHERE STATEMENT_ID = 'omni_2' ORDER BY ID
< result
< column "ID" "NUMBER"
< column "PARENT_ID" "NUMBER"
< column "OPERATION" "NCHAR"
< column "CARDINALITY" "NUMBER"
< row
: string
| 0
: null
: string
| MERGE STATEMENT
: string
| 1
< row
: string
| 1
: string
| 0
: string
| MERGE
: null
< row
: string
| 2
: string
| 1
: string
| VIEW
: null
< row
: string
| 3
: string
| 2
: string
| HASH JOIN
: string
| 9
< row
: string
| 4
: string
| 3
: string
| TABLE ACCESS
: string
| 9
< row
: string
| 5
: string
| 3
: string
| TABLE ACCESS
: string
| 1
> exec
| DELETE FROM PLAN_TABLE WHERE STATEMENT_ID = 'omni_2'
< ok 6 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_3' FOR INSERT /*+ OPT_PARAM('optimizer_dynamic_sampling' 11) */ ALL INTO OMNI_APP.BIG (ID) VALUES (1) INTO OMNI_APP.BIG (ID) VALUES (2) SELECT ID FROM OMNI_APP.BIG
< ok 0 0
> query
| SELECT ID, PARENT_ID, OPERATION, CARDINALITY FROM PLAN_TABLE WHERE STATEMENT_ID = 'omni_3' ORDER BY ID
< result
< column "ID" "NUMBER"
< column "PARENT_ID" "NUMBER"
< column "OPERATION" "NCHAR"
< column "CARDINALITY" "NUMBER"
< row
: string
| 0
: null
: string
| INSERT STATEMENT
: string
| 1000
< row
: string
| 1
: string
| 0
: string
| MULTI-TABLE INSERT
: null
< row
: string
| 2
: string
| 1
: string
| TABLE ACCESS
: string
| 1000
< row
: string
| 3
: string
| 1
: string
| INTO
: null
< row
: string
| 4
: string
| 1
: string
| INTO
: null
> exec
| DELETE FROM PLAN_TABLE WHERE STATEMENT_ID = 'omni_3'
< ok 5 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_4' FOR INSERT /*+ OPT_PARAM('optimizer_dynamic_sampling' 11) */ FIRST WHEN ID < 5 THEN INTO OMNI_APP.BIG (ID) VALUES (ID) ELSE INTO OMNI_APP.BIG (ID) VALUES (ID) SELECT ID FROM OMNI_APP.BIG
< ok 0 0
> query
| SELECT ID, PARENT_ID, OPERATION, CARDINALITY FROM PLAN_TABLE WHERE STATEMENT_ID = 'omni_4' ORDER BY ID
< result
< column "ID" "NUMBER"
< column "PARENT_ID" "NUMBER"
< column "OPERATION" "NCHAR"
< column "CARDINALITY" "NUMBER"
< row
: string
| 0
: null
: string
| INSERT STATEMENT
: string
| 1000
< row
: string
| 1
: string
| 0
: string
| MULTI-TABLE INSERT
: null
< row
: string
| 2
: string
| 1
: string
| TABLE ACCESS
: string
| 1000
< row
: string
| 3
: string
| 1
: string
| INTO
: null
< row
: string
| 4
: string
| 1
: string
| INTO
: null
> exec
| DELETE FROM PLAN_TABLE WHERE STATEMENT_ID = 'omni_4'
< ok 5 0
