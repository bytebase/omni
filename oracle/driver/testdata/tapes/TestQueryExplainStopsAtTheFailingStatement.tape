-- pin: default gvenzl/oracle-free:23.26.3-slim-faststart sha256:f5ff19033860d662c821cb04eb10483fa94f14f78eae252d054291ea07028093 23.26.3.0.0

== TestQueryExplainStopsAtTheFailingStatement default#1
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
| CREATE TABLE OMNI_APP.Q (
| 	ID NUMBER(4) NOT NULL,
| 	LABEL VARCHAR2(20),
| 	AMOUNT NUMBER(8, 2),
| 	FLAG NUMBER(1),
| 	CREATED DATE,
| 	AT TIMESTAMP(3),
| 	AT_ZONE TIMESTAMP WITH TIME ZONE,
| 	DOC CLOB,
| 	BIN BLOB,
| 	CONSTRAINT Q_PK PRIMARY KEY (ID)
| )
< ok 0 0
> exec
|
| INSERT INTO OMNI_APP.Q (ID, LABEL, AMOUNT, FLAG, CREATED, AT, AT_ZONE, DOC, BIN)
| 	VALUES (1, 'a', 1.5, 1, DATE '2026-03-04', TIMESTAMP '2026-03-04 05:06:07.891', TIMESTAMP '2026-03-04 05:06:07.891 +02:00', 'doc', HEXTORAW('01'))
< ok 1 0
> exec
|
| INSERT INTO OMNI_APP.Q (ID, LABEL) VALUES (2, 'b')
< ok 1 0
> exec
|
| INSERT INTO OMNI_APP.Q (ID, LABEL) VALUES (3, NULL)
< ok 1 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_1' FOR SELECT ID FROM OMNI_APP.Q
< ok 0 0
> query
| SELECT LPAD(' ', LEVEL-1) || OPERATION || ' (' || OPTIONS || ')' "Operation", OBJECT_NAME "Object", OPTIMIZER "Optimizer", COST "Cost", CARDINALITY "Cardinality", BYTES "Bytes", PARTITION_START "Partition Start", PARTITION_ID "Partition ID", ACCESS_PREDICATES "Access Predicates" FROM PLAN_TABLE START WITH ID = 0 AND statement_id = 'omni_1' CONNECT BY PRIOR ID=PARENT_ID AND statement_id = 'omni_1' ORDER BY id
< result
< column "Operation" "NCHAR"
< column "Object" "NCHAR"
< column "Optimizer" "NCHAR"
< column "Cost" "NUMBER"
< column "Cardinality" "NUMBER"
< column "Bytes" "NUMBER"
< column "Partition Start" "NCHAR"
< column "Partition ID" "NUMBER"
< column "Access Predicates" "NCHAR"
< row
: string
| SELECT STATEMENT ()
: null
: string
| ALL_ROWS
: string
| 2
: string
| 3
: string
| 39
: null
: null
: null
< row
: string
|  INDEX (FAST FULL SCAN)
: string
| Q_PK
: null
: string
| 2
: string
| 3
: string
| 39
: null
: null
: null
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_2' FOR  SELECT NOSUCH FROM OMNI_APP.Q
< error
| 904 ORA-00904: "NOSUCH": invalid identifier
|
