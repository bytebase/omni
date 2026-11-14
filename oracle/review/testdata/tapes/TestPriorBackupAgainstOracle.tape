-- pin: default gvenzl/oracle-free:23.26.3-slim-faststart sha256:f5ff19033860d662c821cb04eb10483fa94f14f78eae252d054291ea07028093 23.26.3.0.0

== TestPriorBackupAgainstOracle default#1
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
> query
| SELECT COUNT(*) FROM dba_users WHERE username = 'BBDATAARCHIVE'
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 1
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
> query
| SELECT COUNT(*) FROM dba_users WHERE username = 'BBDATAARCHIVE'
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 1
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
> query
| SELECT COUNT(*) FROM dba_users WHERE username = 'BBDATAARCHIVE'
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 1
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
> query
| SELECT COUNT(*) FROM dba_users WHERE username = 'BBDATAARCHIVE'
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 1
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
> query
| SELECT COUNT(*) FROM dba_users WHERE username = 'BBDATAARCHIVE'
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 1
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
> query
| SELECT COUNT(*) FROM dba_users WHERE username = 'BBDATAARCHIVE'
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 1
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
> query
| SELECT COUNT(*) FROM dba_users WHERE username = 'BBDATAARCHIVE'
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 1
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
> query
| SELECT COUNT(*) FROM dba_users WHERE username = 'BBDATAARCHIVE'
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 1
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
> query
| SELECT COUNT(*) FROM dba_users WHERE username = 'BBDATAARCHIVE'
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 0
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
> query
| SELECT COUNT(*) FROM dba_users WHERE username = 'BBDATAARCHIVE'
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 1
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
> query
| SELECT COUNT(*) FROM dba_users WHERE username = 'BBDATAARCHIVE'
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 1
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
> query
| SELECT COUNT(*) FROM dba_users WHERE username = 'BBDATAARCHIVE'
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 1
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
> query
| SELECT COUNT(*) FROM dba_users WHERE username = 'BBDATAARCHIVE'
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 1
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
> query
| SELECT COUNT(*) FROM dba_users WHERE username = 'BBDATAARCHIVE'
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 1
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
> query
| SELECT COUNT(*) FROM dba_users WHERE username = 'BBDATAARCHIVE'
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 1
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
> query
| SELECT COUNT(*) FROM dba_users WHERE username = 'BBDATAARCHIVE'
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 1
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
> query
| SELECT COUNT(*) FROM dba_users WHERE username = 'BBDATAARCHIVE'
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 1
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
> query
| SELECT COUNT(*) FROM dba_users WHERE username = 'BBDATAARCHIVE'
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 1
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
> query
| SELECT COUNT(*) FROM dba_users WHERE username = 'BBDATAARCHIVE'
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 1
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
> query
| SELECT COUNT(*) FROM dba_users WHERE username = 'BBDATAARCHIVE'
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 1
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
> query
| SELECT COUNT(*) FROM dba_users WHERE username = 'BBDATAARCHIVE'
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 1

== TestPriorBackupAgainstOracle/UPDATE_and_DELETE_on_one_table default#1
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

== TestPriorBackupAgainstOracle/a_LONG_column default#1
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
| CREATE TABLE "BBDATAARCHIVE"."_0_LONGT_SHOP" AS
| SELECT "LONGT".* FROM longt WHERE id = 1
< error
| 997 ORA-00997: illegal use of LONG datatype
|

== TestPriorBackupAgainstOracle/a_column_in_a_subquery,_not_judged default#1
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
| SYS_C9000020
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
| SYS_C9000020
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
| CREATE TABLE "BBDATAARCHIVE"."_0_PLAIN_SHOP" AS
| SELECT "PLAIN".* FROM plain WHERE id IN (SELECT gone FROM t)
< error
| 904 ORA-00904: "GONE": invalid identifier
|

== TestPriorBackupAgainstOracle/a_column_of_a_table_an_UPDATE's_FROM_list_joins default#1
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
| SYS_C9000017
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
| SYS_C9000017
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
| CREATE TABLE "BBDATAARCHIVE"."_0_T_SHOP" AS
| SELECT "T".* FROM t, p WHERE t.p_id = p.id AND p.memo = 1
< error
| 904 ORA-00904: "P"."MEMO": invalid identifier
|

== TestPriorBackupAgainstOracle/a_column_of_a_table_the_change_recreates default#1
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
| SYS_C9000018
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
| SYS_C9000018
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
| CREATE TABLE "BBDATAARCHIVE"."_0_NOKEY_SHOP" AS
| SELECT "NOKEY".* FROM nokey WHERE b = 1
< error
| 904 ORA-00904: "B": invalid identifier
|

== TestPriorBackupAgainstOracle/a_column_the_SET_clause_names,_not_copied default#1
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
| SYS_C9000019
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
| SYS_C9000019
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
| CREATE TABLE "BBDATAARCHIVE"."_0_T_SHOP" AS
| SELECT "T".* FROM t WHERE id = 1
< ok 1 0
> exec
| COMMENT ON TABLE "BBDATAARCHIVE"."_0_T_SHOP" IS 'task 1, source table (SHOP, T)'
< ok 0 0

== TestPriorBackupAgainstOracle/a_column_the_change_adds,_unqualified default#1
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
| SYS_C9000015
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
| SYS_C9000015
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
| CREATE TABLE "BBDATAARCHIVE"."_0_T_SHOP" AS
| SELECT "T".* FROM t WHERE gone = 1
< error
| 904 ORA-00904: "GONE": invalid identifier
|

== TestPriorBackupAgainstOracle/a_column_the_change_renames_to,_through_an_alias default#1
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
| SYS_C9000016
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
| SYS_C9000016
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
| CREATE TABLE "BBDATAARCHIVE"."_0_P_SHOP" AS
| SELECT "X".* FROM p x WHERE x.memo = 'a'
< error
| 904 ORA-00904: "X"."MEMO": invalid identifier
|

== TestPriorBackupAgainstOracle/a_database_link default#1
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
| SYS_C9000007
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
| SYS_C9000007
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
| CREATE TABLE "BBDATAARCHIVE"."_1_NOKEY_SHOP" AS
| SELECT "NOKEY".* FROM nokey WHERE a = 1
< ok 1 0
> exec
| COMMENT ON TABLE "BBDATAARCHIVE"."_1_NOKEY_SHOP" IS 'task 1, source table (SHOP, NOKEY)'
< ok 0 0

== TestPriorBackupAgainstOracle/a_plain_backup default#1
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
> exec
| CREATE TABLE "BBDATAARCHIVE"."_0_T_SHOP" AS
| SELECT "T".* FROM t WHERE id = 1
< ok 1 0
> exec
| COMMENT ON TABLE "BBDATAARCHIVE"."_0_T_SHOP" IS 'task 1, source table (SHOP, T)'
< ok 0 0
> exec
| CREATE TABLE "BBDATAARCHIVE"."_1_NOKEY_SHOP" AS
| SELECT "NOKEY".* FROM nokey WHERE a = 1
< ok 1 0
> exec
| COMMENT ON TABLE "BBDATAARCHIVE"."_1_NOKEY_SHOP" IS 'task 1, source table (SHOP, NOKEY)'
< ok 0 0
> exec
| CREATE TABLE "BBDATAARCHIVE"."_2_PLAIN_SHOP" AS
| SELECT "X".* FROM plain x
< ok 1 0
> exec
| COMMENT ON TABLE "BBDATAARCHIVE"."_2_PLAIN_SHOP" IS 'task 1, source table (SHOP, PLAIN)'
< ok 0 0

== TestPriorBackupAgainstOracle/a_quote_in_a_table_name default#1
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
| SYS_C9000006
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
| SYS_C9000006
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
| CREATE TABLE "BBDATAARCHIVE"."_0_a'b_SHOP" AS
| SELECT "a'b".* FROM "a'b" WHERE id = 1
< ok 1 0
> exec
| COMMENT ON TABLE "BBDATAARCHIVE"."_0_a'b_SHOP" IS 'task 1, source table (SHOP, a'b)'
< error
| 3049 ORA-03049: SQL keyword 'B' is not syntactically valid following '...<string>'
|

== TestPriorBackupAgainstOracle/a_table_the_change_creates,_in_a_group_of_two default#1
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
| SYS_C9000012
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
| SYS_C9000012
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
| CREATE TABLE "BBDATAARCHIVE"."_0_N2_SHOP" AS
| SELECT "SHOP"."N2".* FROM "SHOP"."N2" WHERE ROWID IN (SELECT "N2".ROWID FROM n2 WHERE id = 1 UNION SELECT "N2".ROWID FROM n2 WHERE id = 2)
< error
| 942 ORA-00942: table or view "SHOP"."N2" does not exist
|

== TestPriorBackupAgainstOracle/a_table_the_change_creates,_in_a_subquery default#1
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
| SYS_C9000011
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
| SYS_C9000011
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
| CREATE TABLE "BBDATAARCHIVE"."_0_T_SHOP" AS
| SELECT "T".* FROM t WHERE id IN (SELECT id FROM shop.n2)
< error
| 942 ORA-00942: table or view "SHOP"."N2" does not exist
|

== TestPriorBackupAgainstOracle/a_table_the_change_creates,_named_with_its_schema default#1
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
| SYS_C9000010
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
| SYS_C9000010
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
| CREATE TABLE "BBDATAARCHIVE"."_0_N2_SHOP" AS
| SELECT "SHOP"."N2".* FROM shop.n2 WHERE id = 1
< error
| 942 ORA-00942: table or view "SHOP"."N2" does not exist
|

== TestPriorBackupAgainstOracle/a_table_the_change_creates,_unqualified,_not_judged default#1
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
| SYS_C9000013
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
| SYS_C9000013
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
| CREATE TABLE "BBDATAARCHIVE"."_0_N2_SHOP" AS
| SELECT "N2".* FROM n2 WHERE id = 1
< error
| 942 ORA-00942: table or view "SHOP"."N2" does not exist
|

== TestPriorBackupAgainstOracle/a_table_the_change_creates,_which_a_public_synonym_names_before default#1
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
| SYS_C9000014
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
| SYS_C9000014
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
| CREATE TABLE "BBDATAARCHIVE"."_0_PUBT_SHOP" AS
| SELECT "PUBT".* FROM pubt WHERE id = 1
< ok 0 0
> exec
| COMMENT ON TABLE "BBDATAARCHIVE"."_0_PUBT_SHOP" IS 'task 1, source table (SHOP, PUBT)'
< ok 0 0

== TestPriorBackupAgainstOracle/a_view_as_the_target default#1
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
| CREATE TABLE "BBDATAARCHIVE"."_0_PV_SHOP" AS
| SELECT "SHOP"."PV".* FROM "SHOP"."PV" WHERE ROWID IN (SELECT "PV".ROWID FROM pv WHERE id = 1 UNION SELECT "PV".ROWID FROM pv WHERE id = 2)
< ok 2 0
> exec
| COMMENT ON TABLE "BBDATAARCHIVE"."_0_PV_SHOP" IS 'task 1, source table (SHOP, PV)'
< ok 0 0

== TestPriorBackupAgainstOracle/after_the_session_moves,_not_judged default#1
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
| SYS_C9000021
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
| SYS_C9000021
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
| CREATE TABLE "BBDATAARCHIVE"."_0_NOKEY_SHOP" AS
| SELECT "SHOP"."NOKEY".* FROM shop.nokey WHERE gone = 1
< error
| 904 ORA-00904: "GONE": invalid identifier
|

== TestPriorBackupAgainstOracle/no_BBDATAARCHIVE default#1
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
| SYS_C9000009
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
| SYS_C9000009
: string
| UNIQUE
> query
| SELECT view_name FROM all_views WHERE owner = 'SHOP' ORDER BY view_name
< result
< column "VIEW_NAME" "NCHAR"
< row
: string
| PV

== TestPriorBackupAgainstOracle/nothing_to_back_up default#1
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
| SYS_C9000008
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
| SYS_C9000008
: string
| UNIQUE
> query
| SELECT view_name FROM all_views WHERE owner = 'SHOP' ORDER BY view_name
< result
< column "VIEW_NAME" "NCHAR"
< row
: string
| PV

== TestPriorBackupAgainstOracle/two_statements_on_one_table default#1
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
| CREATE TABLE "BBDATAARCHIVE"."_0_T_SHOP" AS
| SELECT "SHOP"."T".* FROM "SHOP"."T" WHERE ROWID IN (SELECT "T".ROWID FROM t WHERE id = 1 UNION SELECT "X".ROWID FROM t x WHERE x.id = 2 UNION SELECT "X".ROWID FROM t x WHERE x.id = 1)
< ok 2 0
> exec
| COMMENT ON TABLE "BBDATAARCHIVE"."_0_T_SHOP" IS 'task 1, source table (SHOP, T)'
< ok 0 0
