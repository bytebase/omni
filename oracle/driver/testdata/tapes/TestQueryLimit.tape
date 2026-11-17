-- pin: default gvenzl/oracle-free:23.26.3-slim-faststart sha256:f5ff19033860d662c821cb04eb10483fa94f14f78eae252d054291ea07028093 23.26.3.0.0

== TestQueryLimit default#1
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
> query
| SELECT ID FROM OMNI_APP.Q ORDER BY ID FETCH NEXT 2 ROWS ONLY
< result
< column "ID" "NUMBER" "float64" precision 4 scale 0
< row
: string
| 1
< row
: string
| 2
> query
| -- the rows in order
| SELECT ID FROM OMNI_APP.Q ORDER BY ID FETCH NEXT 2 ROWS ONLY
< result
< column "ID" "NUMBER" "float64" precision 4 scale 0
< row
: string
| 1
< row
: string
| 2
> query
| SELECT ID, LABEL FROM OMNI_APP.Q ORDER BY ID
< result
< column "ID" "NUMBER" "float64" precision 4 scale 0
< column "LABEL" "NCHAR" "string"
< row
: string
| 1
: string
| a
< row
: string
| 2
: string
| b
< row
: string
| 3
: null
