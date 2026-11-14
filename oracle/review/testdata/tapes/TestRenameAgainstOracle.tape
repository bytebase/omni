-- pin: default gvenzl/oracle-free:23.26.3-slim-faststart sha256:f5ff19033860d662c821cb04eb10483fa94f14f78eae252d054291ea07028093 23.26.3.0.0

== TestRenameAgainstOracle default#1
> exec
| CREATE TABLE system.pubt (id NUMBER, a NUMBER)
< ok 0 0
> exec
| GRANT SELECT ON system.pubt TO PUBLIC
< ok 0 0
> exec
| CREATE PUBLIC SYNONYM pubt FOR system.pubt
< ok 0 0
> exec
| BEGIN
|   FOR u IN (SELECT username FROM dba_users WHERE username IN ('SHOP', 'BBDATAARCHIVE')) LOOP
|     EXECUTE IMMEDIATE 'DROP USER ' || u.username || ' CASCADE';
|   END LOOP;
| END;
< ok 1 0
> exec
| CREATE USER shop IDENTIFIED BY omni_test QUOTA UNLIMITED ON users
< ok 0 0
> exec
| GRANT CREATE SESSION, CREATE TABLE, CREATE VIEW, CREATE MATERIALIZED VIEW, CREATE SEQUENCE, CREATE SYNONYM, CREATE PROCEDURE, CREATE TRIGGER, CREATE ANY TABLE, COMMENT ANY TABLE TO shop
< ok 0 0
> exec
| CREATE USER bbdataarchive IDENTIFIED BY omni_test QUOTA UNLIMITED ON users
< ok 0 0
> exec
| BEGIN
|   FOR u IN (SELECT username FROM dba_users WHERE username IN ('SHOP', 'BBDATAARCHIVE')) LOOP
|     EXECUTE IMMEDIATE 'DROP USER ' || u.username || ' CASCADE';
|   END LOOP;
| END;
< ok 1 0
> exec
| CREATE USER shop IDENTIFIED BY omni_test QUOTA UNLIMITED ON users
< ok 0 0
> exec
| GRANT CREATE SESSION, CREATE TABLE, CREATE VIEW, CREATE MATERIALIZED VIEW, CREATE SEQUENCE, CREATE SYNONYM, CREATE PROCEDURE, CREATE TRIGGER, CREATE ANY TABLE, COMMENT ANY TABLE TO shop
< ok 0 0
> exec
| CREATE USER bbdataarchive IDENTIFIED BY omni_test QUOTA UNLIMITED ON users
< ok 0 0
> exec
| BEGIN
|   FOR u IN (SELECT username FROM dba_users WHERE username IN ('SHOP', 'BBDATAARCHIVE')) LOOP
|     EXECUTE IMMEDIATE 'DROP USER ' || u.username || ' CASCADE';
|   END LOOP;
| END;
< ok 1 0
> exec
| CREATE USER shop IDENTIFIED BY omni_test QUOTA UNLIMITED ON users
< ok 0 0
> exec
| GRANT CREATE SESSION, CREATE TABLE, CREATE VIEW, CREATE MATERIALIZED VIEW, CREATE SEQUENCE, CREATE SYNONYM, CREATE PROCEDURE, CREATE TRIGGER, CREATE ANY TABLE, COMMENT ANY TABLE TO shop
< ok 0 0
> exec
| CREATE USER bbdataarchive IDENTIFIED BY omni_test QUOTA UNLIMITED ON users
< ok 0 0
> exec
| BEGIN
|   FOR u IN (SELECT username FROM dba_users WHERE username IN ('SHOP', 'BBDATAARCHIVE')) LOOP
|     EXECUTE IMMEDIATE 'DROP USER ' || u.username || ' CASCADE';
|   END LOOP;
| END;
< ok 1 0
> exec
| CREATE USER shop IDENTIFIED BY omni_test QUOTA UNLIMITED ON users
< ok 0 0
> exec
| GRANT CREATE SESSION, CREATE TABLE, CREATE VIEW, CREATE MATERIALIZED VIEW, CREATE SEQUENCE, CREATE SYNONYM, CREATE PROCEDURE, CREATE TRIGGER, CREATE ANY TABLE, COMMENT ANY TABLE TO shop
< ok 0 0
> exec
| CREATE USER bbdataarchive IDENTIFIED BY omni_test QUOTA UNLIMITED ON users
< ok 0 0
> exec
| BEGIN
|   FOR u IN (SELECT username FROM dba_users WHERE username IN ('SHOP', 'BBDATAARCHIVE')) LOOP
|     EXECUTE IMMEDIATE 'DROP USER ' || u.username || ' CASCADE';
|   END LOOP;
| END;
< ok 1 0
> exec
| CREATE USER shop IDENTIFIED BY omni_test QUOTA UNLIMITED ON users
< ok 0 0
> exec
| GRANT CREATE SESSION, CREATE TABLE, CREATE VIEW, CREATE MATERIALIZED VIEW, CREATE SEQUENCE, CREATE SYNONYM, CREATE PROCEDURE, CREATE TRIGGER, CREATE ANY TABLE, COMMENT ANY TABLE TO shop
< ok 0 0
> exec
| CREATE USER bbdataarchive IDENTIFIED BY omni_test QUOTA UNLIMITED ON users
< ok 0 0

== TestRenameAgainstOracle/a_materialized_view,_which_RENAME_refuses default#1
> exec
| CREATE TABLE p (id NUMBER CONSTRAINT pk_p PRIMARY KEY, note VARCHAR2(20))
< ok 0 0
> exec
| CREATE TABLE t (id NUMBER CONSTRAINT pk_t PRIMARY KEY, code VARCHAR2(20) CONSTRAINT nn_code NOT NULL CONSTRAINT uk_code UNIQUE, p_id NUMBER CONSTRAINT fk_p REFERENCES p (id), n NUMBER CONSTRAINT ck_n CHECK (n > 0))
< ok 0 0
> exec
| CREATE INDEX ix_n ON t (n)
< ok 0 0
> exec
| CREATE TABLE nokey (a NUMBER)
< ok 0 0
> exec
| CREATE TABLE plain (id NUMBER CONSTRAINT pk_plain PRIMARY KEY, v NUMBER)
< ok 0 0
> exec
| CREATE UNIQUE INDEX ux_v ON plain (v)
< ok 0 0
> exec
| CREATE TABLE longt (id NUMBER CONSTRAINT pk_longt PRIMARY KEY, body LONG)
< ok 0 0
> exec
| CREATE TABLE "a'b" (id NUMBER PRIMARY KEY)
< ok 0 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a')
< ok 1 0
> exec
| INSERT INTO p VALUES (2, 'b')
< ok 1 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1)
< ok 1 0
> exec
| INSERT INTO t VALUES (2, 'y', 2, 2)
< ok 1 0
> exec
| INSERT INTO nokey VALUES (1)
< ok 1 0
> exec
| INSERT INTO plain VALUES (1, 1)
< ok 1 0
> exec
| INSERT INTO longt VALUES (1, 'x')
< ok 1 0
> exec
| INSERT INTO "a'b" VALUES (1)
< ok 1 0
> query
| SELECT table_name FROM all_tables WHERE owner = 'SHOP' AND table_name NOT IN (SELECT mview_name FROM all_mviews WHERE owner = 'SHOP') ORDER BY table_name
< result
< column "TABLE_NAME" "NCHAR"
< row
: string
| LONGT
< row
: string
| NOKEY
< row
: string
| P
< row
: string
| PLAIN
< row
: string
| T
< row
: string
| a'b
> query
| SELECT table_name, column_name, data_type FROM all_tab_columns WHERE owner = 'SHOP' AND column_id IS NOT NULL ORDER BY table_name, column_id
< result
< column "TABLE_NAME" "NCHAR"
< column "COLUMN_NAME" "NCHAR"
< column "DATA_TYPE" "NCHAR"
< row
: string
| LONGT
: string
| ID
: string
| NUMBER
< row
: string
| LONGT
: string
| BODY
: string
| LONG
< row
: string
| NOKEY
: string
| A
: string
| NUMBER
< row
: string
| P
: string
| ID
: string
| NUMBER
< row
: string
| P
: string
| NOTE
: string
| VARCHAR2
< row
: string
| PLAIN
: string
| ID
: string
| NUMBER
< row
: string
| PLAIN
: string
| V
: string
| NUMBER
< row
: string
| PV
: string
| ID
: string
| NUMBER
< row
: string
| PV
: string
| NOTE
: string
| VARCHAR2
< row
: string
| T
: string
| ID
: string
| NUMBER
< row
: string
| T
: string
| CODE
: string
| VARCHAR2
< row
: string
| T
: string
| P_ID
: string
| NUMBER
< row
: string
| T
: string
| N
: string
| NUMBER
< row
: string
| a'b
: string
| ID
: string
| NUMBER
> query
| SELECT table_name, constraint_name, constraint_type FROM all_constraints WHERE owner = 'SHOP' ORDER BY table_name, constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< row
: string
| LONGT
: string
| PK_LONGT
: string
| P
< row
: string
| P
: string
| PK_P
: string
| P
< row
: string
| PLAIN
: string
| PK_PLAIN
: string
| P
< row
: string
| T
: string
| CK_N
: string
| C
< row
: string
| T
: string
| FK_P
: string
| R
< row
: string
| T
: string
| NN_CODE
: string
| C
< row
: string
| T
: string
| PK_T
: string
| P
< row
: string
| T
: string
| UK_CODE
: string
| U
< row
: string
| a'b
: string
| SYS_C9000004
: string
| P
> query
| SELECT table_name, index_name, uniqueness FROM all_indexes WHERE owner = 'SHOP' ORDER BY table_name, index_name
< result
< column "TABLE_NAME" "NCHAR"
< column "INDEX_NAME" "NCHAR"
< column "UNIQUENESS" "NCHAR"
< row
: string
| LONGT
: string
| PK_LONGT
: string
| UNIQUE
< row
: string
| P
: string
| PK_P
: string
| UNIQUE
< row
: string
| PLAIN
: string
| PK_PLAIN
: string
| UNIQUE
< row
: string
| PLAIN
: string
| UX_V
: string
| UNIQUE
< row
: string
| T
: string
| IX_N
: string
| NONUNIQUE
< row
: string
| T
: string
| PK_T
: string
| UNIQUE
< row
: string
| T
: string
| UK_CODE
: string
| UNIQUE
< row
: string
| a'b
: string
| SYS_C9000004
: string
| UNIQUE
> query
| SELECT view_name FROM all_views WHERE owner = 'SHOP' ORDER BY view_name
< result
< column "VIEW_NAME" "NCHAR"
< row
: string
| PV
> exec
| CREATE MATERIALIZED VIEW mv AS SELECT * FROM p
< ok 0 0
> query
| SELECT object_type FROM user_objects WHERE object_name = :1 ORDER BY object_type
: string
| MV
< result
< column "OBJECT_TYPE" "NCHAR"
< row
: string
| MATERIALIZED VIEW
< row
: string
| TABLE
> exec
| RENAME mv TO mv2
< error
| 32318 ORA-32318: cannot rename a materialized view
|

== TestRenameAgainstOracle/a_name_dropped_and_made_again_as_another_kind default#1
> exec
| CREATE TABLE p (id NUMBER CONSTRAINT pk_p PRIMARY KEY, note VARCHAR2(20))
< ok 0 0
> exec
| CREATE TABLE t (id NUMBER CONSTRAINT pk_t PRIMARY KEY, code VARCHAR2(20) CONSTRAINT nn_code NOT NULL CONSTRAINT uk_code UNIQUE, p_id NUMBER CONSTRAINT fk_p REFERENCES p (id), n NUMBER CONSTRAINT ck_n CHECK (n > 0))
< ok 0 0
> exec
| CREATE INDEX ix_n ON t (n)
< ok 0 0
> exec
| CREATE TABLE nokey (a NUMBER)
< ok 0 0
> exec
| CREATE TABLE plain (id NUMBER CONSTRAINT pk_plain PRIMARY KEY, v NUMBER)
< ok 0 0
> exec
| CREATE UNIQUE INDEX ux_v ON plain (v)
< ok 0 0
> exec
| CREATE TABLE longt (id NUMBER CONSTRAINT pk_longt PRIMARY KEY, body LONG)
< ok 0 0
> exec
| CREATE TABLE "a'b" (id NUMBER PRIMARY KEY)
< ok 0 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a')
< ok 1 0
> exec
| INSERT INTO p VALUES (2, 'b')
< ok 1 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1)
< ok 1 0
> exec
| INSERT INTO t VALUES (2, 'y', 2, 2)
< ok 1 0
> exec
| INSERT INTO nokey VALUES (1)
< ok 1 0
> exec
| INSERT INTO plain VALUES (1, 1)
< ok 1 0
> exec
| INSERT INTO longt VALUES (1, 'x')
< ok 1 0
> exec
| INSERT INTO "a'b" VALUES (1)
< ok 1 0
> query
| SELECT table_name FROM all_tables WHERE owner = 'SHOP' AND table_name NOT IN (SELECT mview_name FROM all_mviews WHERE owner = 'SHOP') ORDER BY table_name
< result
< column "TABLE_NAME" "NCHAR"
< row
: string
| LONGT
< row
: string
| NOKEY
< row
: string
| P
< row
: string
| PLAIN
< row
: string
| T
< row
: string
| a'b
> query
| SELECT table_name, column_name, data_type FROM all_tab_columns WHERE owner = 'SHOP' AND column_id IS NOT NULL ORDER BY table_name, column_id
< result
< column "TABLE_NAME" "NCHAR"
< column "COLUMN_NAME" "NCHAR"
< column "DATA_TYPE" "NCHAR"
< row
: string
| LONGT
: string
| ID
: string
| NUMBER
< row
: string
| LONGT
: string
| BODY
: string
| LONG
< row
: string
| NOKEY
: string
| A
: string
| NUMBER
< row
: string
| P
: string
| ID
: string
| NUMBER
< row
: string
| P
: string
| NOTE
: string
| VARCHAR2
< row
: string
| PLAIN
: string
| ID
: string
| NUMBER
< row
: string
| PLAIN
: string
| V
: string
| NUMBER
< row
: string
| PV
: string
| ID
: string
| NUMBER
< row
: string
| PV
: string
| NOTE
: string
| VARCHAR2
< row
: string
| T
: string
| ID
: string
| NUMBER
< row
: string
| T
: string
| CODE
: string
| VARCHAR2
< row
: string
| T
: string
| P_ID
: string
| NUMBER
< row
: string
| T
: string
| N
: string
| NUMBER
< row
: string
| a'b
: string
| ID
: string
| NUMBER
> query
| SELECT table_name, constraint_name, constraint_type FROM all_constraints WHERE owner = 'SHOP' ORDER BY table_name, constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< row
: string
| LONGT
: string
| PK_LONGT
: string
| P
< row
: string
| P
: string
| PK_P
: string
| P
< row
: string
| PLAIN
: string
| PK_PLAIN
: string
| P
< row
: string
| T
: string
| CK_N
: string
| C
< row
: string
| T
: string
| FK_P
: string
| R
< row
: string
| T
: string
| NN_CODE
: string
| C
< row
: string
| T
: string
| PK_T
: string
| P
< row
: string
| T
: string
| UK_CODE
: string
| U
< row
: string
| a'b
: string
| SYS_C9000003
: string
| P
> query
| SELECT table_name, index_name, uniqueness FROM all_indexes WHERE owner = 'SHOP' ORDER BY table_name, index_name
< result
< column "TABLE_NAME" "NCHAR"
< column "INDEX_NAME" "NCHAR"
< column "UNIQUENESS" "NCHAR"
< row
: string
| LONGT
: string
| PK_LONGT
: string
| UNIQUE
< row
: string
| P
: string
| PK_P
: string
| UNIQUE
< row
: string
| PLAIN
: string
| PK_PLAIN
: string
| UNIQUE
< row
: string
| PLAIN
: string
| UX_V
: string
| UNIQUE
< row
: string
| T
: string
| IX_N
: string
| NONUNIQUE
< row
: string
| T
: string
| PK_T
: string
| UNIQUE
< row
: string
| T
: string
| UK_CODE
: string
| UNIQUE
< row
: string
| a'b
: string
| SYS_C9000003
: string
| UNIQUE
> query
| SELECT view_name FROM all_views WHERE owner = 'SHOP' ORDER BY view_name
< result
< column "VIEW_NAME" "NCHAR"
< row
: string
| PV
> exec
| DROP TABLE nokey
< ok 0 0
> exec
| CREATE VIEW nokey AS SELECT 1 a FROM dual
< ok 0 0
> query
| SELECT object_type FROM user_objects WHERE object_name = :1 ORDER BY object_type
: string
| NOKEY
< result
< column "OBJECT_TYPE" "NCHAR"
< row
: string
| VIEW
> exec
| RENAME nokey TO nk
< ok 0 0
> exec
| DROP VIEW pv
< ok 0 0
> exec
| CREATE SEQUENCE pv
< ok 0 0
> query
| SELECT object_type FROM user_objects WHERE object_name = :1 ORDER BY object_type
: string
| PV
< result
< column "OBJECT_TYPE" "NCHAR"
< row
: string
| SEQUENCE
> exec
| RENAME pv TO pv2
< ok 0 0

== TestRenameAgainstOracle/a_table_renamed_by_ALTER_TABLE,_then_by_RENAME default#1
> exec
| CREATE TABLE p (id NUMBER CONSTRAINT pk_p PRIMARY KEY, note VARCHAR2(20))
< ok 0 0
> exec
| CREATE TABLE t (id NUMBER CONSTRAINT pk_t PRIMARY KEY, code VARCHAR2(20) CONSTRAINT nn_code NOT NULL CONSTRAINT uk_code UNIQUE, p_id NUMBER CONSTRAINT fk_p REFERENCES p (id), n NUMBER CONSTRAINT ck_n CHECK (n > 0))
< ok 0 0
> exec
| CREATE INDEX ix_n ON t (n)
< ok 0 0
> exec
| CREATE TABLE nokey (a NUMBER)
< ok 0 0
> exec
| CREATE TABLE plain (id NUMBER CONSTRAINT pk_plain PRIMARY KEY, v NUMBER)
< ok 0 0
> exec
| CREATE UNIQUE INDEX ux_v ON plain (v)
< ok 0 0
> exec
| CREATE TABLE longt (id NUMBER CONSTRAINT pk_longt PRIMARY KEY, body LONG)
< ok 0 0
> exec
| CREATE TABLE "a'b" (id NUMBER PRIMARY KEY)
< ok 0 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a')
< ok 1 0
> exec
| INSERT INTO p VALUES (2, 'b')
< ok 1 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1)
< ok 1 0
> exec
| INSERT INTO t VALUES (2, 'y', 2, 2)
< ok 1 0
> exec
| INSERT INTO nokey VALUES (1)
< ok 1 0
> exec
| INSERT INTO plain VALUES (1, 1)
< ok 1 0
> exec
| INSERT INTO longt VALUES (1, 'x')
< ok 1 0
> exec
| INSERT INTO "a'b" VALUES (1)
< ok 1 0
> query
| SELECT table_name FROM all_tables WHERE owner = 'SHOP' AND table_name NOT IN (SELECT mview_name FROM all_mviews WHERE owner = 'SHOP') ORDER BY table_name
< result
< column "TABLE_NAME" "NCHAR"
< row
: string
| LONGT
< row
: string
| NOKEY
< row
: string
| P
< row
: string
| PLAIN
< row
: string
| T
< row
: string
| a'b
> query
| SELECT table_name, column_name, data_type FROM all_tab_columns WHERE owner = 'SHOP' AND column_id IS NOT NULL ORDER BY table_name, column_id
< result
< column "TABLE_NAME" "NCHAR"
< column "COLUMN_NAME" "NCHAR"
< column "DATA_TYPE" "NCHAR"
< row
: string
| LONGT
: string
| ID
: string
| NUMBER
< row
: string
| LONGT
: string
| BODY
: string
| LONG
< row
: string
| NOKEY
: string
| A
: string
| NUMBER
< row
: string
| P
: string
| ID
: string
| NUMBER
< row
: string
| P
: string
| NOTE
: string
| VARCHAR2
< row
: string
| PLAIN
: string
| ID
: string
| NUMBER
< row
: string
| PLAIN
: string
| V
: string
| NUMBER
< row
: string
| PV
: string
| ID
: string
| NUMBER
< row
: string
| PV
: string
| NOTE
: string
| VARCHAR2
< row
: string
| T
: string
| ID
: string
| NUMBER
< row
: string
| T
: string
| CODE
: string
| VARCHAR2
< row
: string
| T
: string
| P_ID
: string
| NUMBER
< row
: string
| T
: string
| N
: string
| NUMBER
< row
: string
| a'b
: string
| ID
: string
| NUMBER
> query
| SELECT table_name, constraint_name, constraint_type FROM all_constraints WHERE owner = 'SHOP' ORDER BY table_name, constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< row
: string
| LONGT
: string
| PK_LONGT
: string
| P
< row
: string
| P
: string
| PK_P
: string
| P
< row
: string
| PLAIN
: string
| PK_PLAIN
: string
| P
< row
: string
| T
: string
| CK_N
: string
| C
< row
: string
| T
: string
| FK_P
: string
| R
< row
: string
| T
: string
| NN_CODE
: string
| C
< row
: string
| T
: string
| PK_T
: string
| P
< row
: string
| T
: string
| UK_CODE
: string
| U
< row
: string
| a'b
: string
| SYS_C9000005
: string
| P
> query
| SELECT table_name, index_name, uniqueness FROM all_indexes WHERE owner = 'SHOP' ORDER BY table_name, index_name
< result
< column "TABLE_NAME" "NCHAR"
< column "INDEX_NAME" "NCHAR"
< column "UNIQUENESS" "NCHAR"
< row
: string
| LONGT
: string
| PK_LONGT
: string
| UNIQUE
< row
: string
| P
: string
| PK_P
: string
| UNIQUE
< row
: string
| PLAIN
: string
| PK_PLAIN
: string
| UNIQUE
< row
: string
| PLAIN
: string
| UX_V
: string
| UNIQUE
< row
: string
| T
: string
| IX_N
: string
| NONUNIQUE
< row
: string
| T
: string
| PK_T
: string
| UNIQUE
< row
: string
| T
: string
| UK_CODE
: string
| UNIQUE
< row
: string
| a'b
: string
| SYS_C9000005
: string
| UNIQUE
> query
| SELECT view_name FROM all_views WHERE owner = 'SHOP' ORDER BY view_name
< result
< column "VIEW_NAME" "NCHAR"
< row
: string
| PV
> exec
| ALTER TABLE plain RENAME TO plain2
< ok 0 0
> query
| SELECT object_type FROM user_objects WHERE object_name = :1 ORDER BY object_type
: string
| PLAIN2
< result
< column "OBJECT_TYPE" "NCHAR"
< row
: string
| TABLE
> exec
| RENAME plain2 TO plain3
< ok 0 0

== TestRenameAgainstOracle/synced_objects default#1
> exec
| CREATE TABLE p (id NUMBER CONSTRAINT pk_p PRIMARY KEY, note VARCHAR2(20))
< ok 0 0
> exec
| CREATE TABLE t (id NUMBER CONSTRAINT pk_t PRIMARY KEY, code VARCHAR2(20) CONSTRAINT nn_code NOT NULL CONSTRAINT uk_code UNIQUE, p_id NUMBER CONSTRAINT fk_p REFERENCES p (id), n NUMBER CONSTRAINT ck_n CHECK (n > 0))
< ok 0 0
> exec
| CREATE INDEX ix_n ON t (n)
< ok 0 0
> exec
| CREATE TABLE nokey (a NUMBER)
< ok 0 0
> exec
| CREATE TABLE plain (id NUMBER CONSTRAINT pk_plain PRIMARY KEY, v NUMBER)
< ok 0 0
> exec
| CREATE UNIQUE INDEX ux_v ON plain (v)
< ok 0 0
> exec
| CREATE TABLE longt (id NUMBER CONSTRAINT pk_longt PRIMARY KEY, body LONG)
< ok 0 0
> exec
| CREATE TABLE "a'b" (id NUMBER PRIMARY KEY)
< ok 0 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a')
< ok 1 0
> exec
| INSERT INTO p VALUES (2, 'b')
< ok 1 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1)
< ok 1 0
> exec
| INSERT INTO t VALUES (2, 'y', 2, 2)
< ok 1 0
> exec
| INSERT INTO nokey VALUES (1)
< ok 1 0
> exec
| INSERT INTO plain VALUES (1, 1)
< ok 1 0
> exec
| INSERT INTO longt VALUES (1, 'x')
< ok 1 0
> exec
| INSERT INTO "a'b" VALUES (1)
< ok 1 0
> query
| SELECT table_name FROM all_tables WHERE owner = 'SHOP' AND table_name NOT IN (SELECT mview_name FROM all_mviews WHERE owner = 'SHOP') ORDER BY table_name
< result
< column "TABLE_NAME" "NCHAR"
< row
: string
| LONGT
< row
: string
| NOKEY
< row
: string
| P
< row
: string
| PLAIN
< row
: string
| T
< row
: string
| a'b
> query
| SELECT table_name, column_name, data_type FROM all_tab_columns WHERE owner = 'SHOP' AND column_id IS NOT NULL ORDER BY table_name, column_id
< result
< column "TABLE_NAME" "NCHAR"
< column "COLUMN_NAME" "NCHAR"
< column "DATA_TYPE" "NCHAR"
< row
: string
| LONGT
: string
| ID
: string
| NUMBER
< row
: string
| LONGT
: string
| BODY
: string
| LONG
< row
: string
| NOKEY
: string
| A
: string
| NUMBER
< row
: string
| P
: string
| ID
: string
| NUMBER
< row
: string
| P
: string
| NOTE
: string
| VARCHAR2
< row
: string
| PLAIN
: string
| ID
: string
| NUMBER
< row
: string
| PLAIN
: string
| V
: string
| NUMBER
< row
: string
| PV
: string
| ID
: string
| NUMBER
< row
: string
| PV
: string
| NOTE
: string
| VARCHAR2
< row
: string
| T
: string
| ID
: string
| NUMBER
< row
: string
| T
: string
| CODE
: string
| VARCHAR2
< row
: string
| T
: string
| P_ID
: string
| NUMBER
< row
: string
| T
: string
| N
: string
| NUMBER
< row
: string
| a'b
: string
| ID
: string
| NUMBER
> query
| SELECT table_name, constraint_name, constraint_type FROM all_constraints WHERE owner = 'SHOP' ORDER BY table_name, constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< row
: string
| LONGT
: string
| PK_LONGT
: string
| P
< row
: string
| P
: string
| PK_P
: string
| P
< row
: string
| PLAIN
: string
| PK_PLAIN
: string
| P
< row
: string
| T
: string
| CK_N
: string
| C
< row
: string
| T
: string
| FK_P
: string
| R
< row
: string
| T
: string
| NN_CODE
: string
| C
< row
: string
| T
: string
| PK_T
: string
| P
< row
: string
| T
: string
| UK_CODE
: string
| U
< row
: string
| a'b
: string
| SYS_C9000001
: string
| P
> query
| SELECT table_name, index_name, uniqueness FROM all_indexes WHERE owner = 'SHOP' ORDER BY table_name, index_name
< result
< column "TABLE_NAME" "NCHAR"
< column "INDEX_NAME" "NCHAR"
< column "UNIQUENESS" "NCHAR"
< row
: string
| LONGT
: string
| PK_LONGT
: string
| UNIQUE
< row
: string
| P
: string
| PK_P
: string
| UNIQUE
< row
: string
| PLAIN
: string
| PK_PLAIN
: string
| UNIQUE
< row
: string
| PLAIN
: string
| UX_V
: string
| UNIQUE
< row
: string
| T
: string
| IX_N
: string
| NONUNIQUE
< row
: string
| T
: string
| PK_T
: string
| UNIQUE
< row
: string
| T
: string
| UK_CODE
: string
| UNIQUE
< row
: string
| a'b
: string
| SYS_C9000001
: string
| UNIQUE
> query
| SELECT view_name FROM all_views WHERE owner = 'SHOP' ORDER BY view_name
< result
< column "VIEW_NAME" "NCHAR"
< row
: string
| PV
> query
| SELECT object_type FROM user_objects WHERE object_name = :1 ORDER BY object_type
: string
| T
< result
< column "OBJECT_TYPE" "NCHAR"
< row
: string
| TABLE
> exec
| RENAME t TO t2
< ok 0 0
> query
| SELECT object_type FROM user_objects WHERE object_name = :1 ORDER BY object_type
: string
| PV
< result
< column "OBJECT_TYPE" "NCHAR"
< row
: string
| VIEW
> exec
| RENAME pv TO pv2
< ok 0 0
> exec
| CREATE SEQUENCE sq
< ok 0 0
> query
| SELECT object_type FROM user_objects WHERE object_name = :1 ORDER BY object_type
: string
| SQ
< result
< column "OBJECT_TYPE" "NCHAR"
< row
: string
| SEQUENCE
> exec
| RENAME sq TO sq2
< ok 0 0

== TestRenameAgainstOracle/what_the_change_makes default#1
> exec
| CREATE TABLE p (id NUMBER CONSTRAINT pk_p PRIMARY KEY, note VARCHAR2(20))
< ok 0 0
> exec
| CREATE TABLE t (id NUMBER CONSTRAINT pk_t PRIMARY KEY, code VARCHAR2(20) CONSTRAINT nn_code NOT NULL CONSTRAINT uk_code UNIQUE, p_id NUMBER CONSTRAINT fk_p REFERENCES p (id), n NUMBER CONSTRAINT ck_n CHECK (n > 0))
< ok 0 0
> exec
| CREATE INDEX ix_n ON t (n)
< ok 0 0
> exec
| CREATE TABLE nokey (a NUMBER)
< ok 0 0
> exec
| CREATE TABLE plain (id NUMBER CONSTRAINT pk_plain PRIMARY KEY, v NUMBER)
< ok 0 0
> exec
| CREATE UNIQUE INDEX ux_v ON plain (v)
< ok 0 0
> exec
| CREATE TABLE longt (id NUMBER CONSTRAINT pk_longt PRIMARY KEY, body LONG)
< ok 0 0
> exec
| CREATE TABLE "a'b" (id NUMBER PRIMARY KEY)
< ok 0 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a')
< ok 1 0
> exec
| INSERT INTO p VALUES (2, 'b')
< ok 1 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1)
< ok 1 0
> exec
| INSERT INTO t VALUES (2, 'y', 2, 2)
< ok 1 0
> exec
| INSERT INTO nokey VALUES (1)
< ok 1 0
> exec
| INSERT INTO plain VALUES (1, 1)
< ok 1 0
> exec
| INSERT INTO longt VALUES (1, 'x')
< ok 1 0
> exec
| INSERT INTO "a'b" VALUES (1)
< ok 1 0
> query
| SELECT table_name FROM all_tables WHERE owner = 'SHOP' AND table_name NOT IN (SELECT mview_name FROM all_mviews WHERE owner = 'SHOP') ORDER BY table_name
< result
< column "TABLE_NAME" "NCHAR"
< row
: string
| LONGT
< row
: string
| NOKEY
< row
: string
| P
< row
: string
| PLAIN
< row
: string
| T
< row
: string
| a'b
> query
| SELECT table_name, column_name, data_type FROM all_tab_columns WHERE owner = 'SHOP' AND column_id IS NOT NULL ORDER BY table_name, column_id
< result
< column "TABLE_NAME" "NCHAR"
< column "COLUMN_NAME" "NCHAR"
< column "DATA_TYPE" "NCHAR"
< row
: string
| LONGT
: string
| ID
: string
| NUMBER
< row
: string
| LONGT
: string
| BODY
: string
| LONG
< row
: string
| NOKEY
: string
| A
: string
| NUMBER
< row
: string
| P
: string
| ID
: string
| NUMBER
< row
: string
| P
: string
| NOTE
: string
| VARCHAR2
< row
: string
| PLAIN
: string
| ID
: string
| NUMBER
< row
: string
| PLAIN
: string
| V
: string
| NUMBER
< row
: string
| PV
: string
| ID
: string
| NUMBER
< row
: string
| PV
: string
| NOTE
: string
| VARCHAR2
< row
: string
| T
: string
| ID
: string
| NUMBER
< row
: string
| T
: string
| CODE
: string
| VARCHAR2
< row
: string
| T
: string
| P_ID
: string
| NUMBER
< row
: string
| T
: string
| N
: string
| NUMBER
< row
: string
| a'b
: string
| ID
: string
| NUMBER
> query
| SELECT table_name, constraint_name, constraint_type FROM all_constraints WHERE owner = 'SHOP' ORDER BY table_name, constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< row
: string
| LONGT
: string
| PK_LONGT
: string
| P
< row
: string
| P
: string
| PK_P
: string
| P
< row
: string
| PLAIN
: string
| PK_PLAIN
: string
| P
< row
: string
| T
: string
| CK_N
: string
| C
< row
: string
| T
: string
| FK_P
: string
| R
< row
: string
| T
: string
| NN_CODE
: string
| C
< row
: string
| T
: string
| PK_T
: string
| P
< row
: string
| T
: string
| UK_CODE
: string
| U
< row
: string
| a'b
: string
| SYS_C9000002
: string
| P
> query
| SELECT table_name, index_name, uniqueness FROM all_indexes WHERE owner = 'SHOP' ORDER BY table_name, index_name
< result
< column "TABLE_NAME" "NCHAR"
< column "INDEX_NAME" "NCHAR"
< column "UNIQUENESS" "NCHAR"
< row
: string
| LONGT
: string
| PK_LONGT
: string
| UNIQUE
< row
: string
| P
: string
| PK_P
: string
| UNIQUE
< row
: string
| PLAIN
: string
| PK_PLAIN
: string
| UNIQUE
< row
: string
| PLAIN
: string
| UX_V
: string
| UNIQUE
< row
: string
| T
: string
| IX_N
: string
| NONUNIQUE
< row
: string
| T
: string
| PK_T
: string
| UNIQUE
< row
: string
| T
: string
| UK_CODE
: string
| UNIQUE
< row
: string
| a'b
: string
| SYS_C9000002
: string
| UNIQUE
> query
| SELECT view_name FROM all_views WHERE owner = 'SHOP' ORDER BY view_name
< result
< column "VIEW_NAME" "NCHAR"
< row
: string
| PV
> exec
| CREATE TABLE n1 (id NUMBER)
< ok 0 0
> query
| SELECT object_type FROM user_objects WHERE object_name = :1 ORDER BY object_type
: string
| N1
< result
< column "OBJECT_TYPE" "NCHAR"
< row
: string
| TABLE
> exec
| RENAME n1 TO n2
< ok 0 0
> query
| SELECT object_type FROM user_objects WHERE object_name = :1 ORDER BY object_type
: string
| N2
< result
< column "OBJECT_TYPE" "NCHAR"
< row
: string
| TABLE
> exec
| RENAME n2 TO n3
< ok 0 0
> exec
| CREATE VIEW w AS SELECT 1 a FROM dual
< ok 0 0
> query
| SELECT object_type FROM user_objects WHERE object_name = :1 ORDER BY object_type
: string
| W
< result
< column "OBJECT_TYPE" "NCHAR"
< row
: string
| VIEW
> exec
| RENAME w TO w2
< ok 0 0
> exec
| CREATE SYNONYM sy FOR p
< ok 0 0
> query
| SELECT object_type FROM user_objects WHERE object_name = :1 ORDER BY object_type
: string
| SY
< result
< column "OBJECT_TYPE" "NCHAR"
< row
: string
| SYNONYM
> exec
| RENAME sy TO sy2
< ok 0 0
