-- pin: default gvenzl/oracle-free:23.26.3-slim-faststart sha256:f5ff19033860d662c821cb04eb10483fa94f14f78eae252d054291ea07028093 23.26.3.0.0

== TestScanRulesAgainstOracle default#1
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

== TestScanRulesAgainstOracle/CREATE_TABLE_AS_SELECT,_which_copies_no_key default#1
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
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000009
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| CREATE TABLE c AS SELECT * FROM plain
< ok 1 0
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000009
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| C
: string
| 0
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000009
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| C
: string
| 0
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| CREATE TABLE d (id PRIMARY KEY) AS SELECT id FROM plain
< ok 1 0
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| D
: string
| P
: string
| SYS_C9000010
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000009
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| C
: string
| 0
< row
: string
| D
: string
| 1
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query ^
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| C
: string
| 0
< row
: string
| D
: string
| 1
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1

== TestScanRulesAgainstOracle/DROP_PRIMARY_KEY_and_DROP_UNIQUE default#1
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
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000004
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| ALTER TABLE plain DROP PRIMARY KEY
< ok 0 0
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000004
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 0
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000004
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 0
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| ALTER TABLE t DROP UNIQUE (code)
< ok 0 0
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| a'b
: string
| P
: string
| SYS_C9000004
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 0
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query ^
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 0
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1

== TestScanRulesAgainstOracle/a_column_dropped_with_its_constraint,_not_judged default#1
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
| SYS_C9000023
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
| SYS_C9000023
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
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000023
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| ALTER TABLE t DROP COLUMN n
< ok 0 0
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000023
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000023
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| ALTER TABLE t DROP CONSTRAINT uk_code
< ok 0 0
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| a'b
: string
| P
: string
| SYS_C9000023
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query ^
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1

== TestScanRulesAgainstOracle/a_constraint_the_change_defines,_then_drops default#1
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
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000005
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| ALTER TABLE nokey ADD CONSTRAINT ck_a CHECK (a > 0)
< ok 0 0
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| NOKEY
: string
| C
: string
| CK_A
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000005
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| NOKEY
: string
| C
: string
| CK_A
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000005
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| ALTER TABLE nokey DROP CONSTRAINT ck_a
< ok 0 0
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000005
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query ^
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1

== TestScanRulesAgainstOracle/a_key_a_table_without_one_gains_and_loses default#1
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
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000014
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| ALTER TABLE nokey ADD PRIMARY KEY (a)
< ok 0 0
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| NOKEY
: string
| P
: string
| SYS_C9000015
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000014
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 1
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| NOKEY
: string
| P
: string
| SYS_C9000015
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000014
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 1
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| ALTER TABLE nokey DROP PRIMARY KEY
< ok 0 0
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000014
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query ^
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1

== TestScanRulesAgainstOracle/a_key_added_later default#1
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
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000011
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| CREATE TABLE a (id NUMBER)
< ok 0 0
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000011
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| A
: string
| 0
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000011
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| A
: string
| 0
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| ALTER TABLE a ADD CONSTRAINT pk_a PRIMARY KEY (id)
< ok 0 0
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| A
: string
| P
: string
| PK_A
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000011
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| A
: string
| 1
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| A
: string
| P
: string
| PK_A
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000011
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| A
: string
| 1
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| CREATE TABLE b (id NUMBER)
< ok 0 0
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| A
: string
| P
: string
| PK_A
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000011
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| A
: string
| 1
< row
: string
| B
: string
| 0
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| A
: string
| P
: string
| PK_A
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000011
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| A
: string
| 1
< row
: string
| B
: string
| 0
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| ALTER TABLE b MODIFY (id PRIMARY KEY)
< ok 0 0
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| A
: string
| P
: string
| PK_A
< row
: string
| B
: string
| P
: string
| SYS_C9000012
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000011
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| A
: string
| 1
< row
: string
| B
: string
| 1
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query ^
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| A
: string
| 1
< row
: string
| B
: string
| 1
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1

== TestScanRulesAgainstOracle/a_key_the_change_makes,_then_drops default#1
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
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000016
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| CREATE TABLE a (id NUMBER CONSTRAINT pk_a PRIMARY KEY)
< ok 0 0
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| A
: string
| P
: string
| PK_A
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000016
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| A
: string
| 1
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| A
: string
| P
: string
| PK_A
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000016
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| A
: string
| 1
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| ALTER TABLE a DROP CONSTRAINT pk_a
< ok 0 0
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000016
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| A
: string
| 0
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query ^
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| A
: string
| 0
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1

== TestScanRulesAgainstOracle/a_name_already_taken_fails_its_statement default#1
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
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000019
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| ALTER TABLE nokey ADD CONSTRAINT ck_n CHECK (a > 0)
< error
| 2264 ORA-02264: name already used by an existing constraint
|
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1

== TestScanRulesAgainstOracle/a_name_the_change_frees_first default#1
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
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000020
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| ALTER TABLE t DROP CONSTRAINT ck_n
< ok 0 0
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000020
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000020
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| ALTER TABLE nokey ADD CONSTRAINT ck_n CHECK (a > 0)
< ok 0 0
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| NOKEY
: string
| C
: string
| CK_N
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000020
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| NOKEY
: string
| C
: string
| CK_N
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000020
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| ALTER TABLE nokey DROP CONSTRAINT ck_n
< ok 0 0
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000020
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query ^
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1

== TestScanRulesAgainstOracle/a_named_NOT_NULL_constraint default#1
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
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000002
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| ALTER TABLE t DROP CONSTRAINT nn_code
< ok 0 0
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000002
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query ^
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1

== TestScanRulesAgainstOracle/a_primary_key_dropped,_then_added_back default#1
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
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000013
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| ALTER TABLE plain DROP PRIMARY KEY
< ok 0 0
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000013
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 0
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000013
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 0
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| ALTER TABLE plain ADD CONSTRAINT pk_plain2 PRIMARY KEY (id)
< ok 0 0
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN2
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000013
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query ^
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1

== TestScanRulesAgainstOracle/a_primary_key_dropped_with_CASCADE default#1
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
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000018
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| ALTER TABLE p DROP PRIMARY KEY CASCADE
< ok 0 0
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000018
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 0
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000018
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 0
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| ALTER TABLE t DROP CONSTRAINT fk_p
< error
| 2443 ORA-02443: Cannot drop constraint  - nonexistent constraint
|
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 0
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1

== TestScanRulesAgainstOracle/a_renamed_table,_not_judged default#1
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
| SYS_C9000024
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
| SYS_C9000024
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
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000024
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| RENAME t TO t2
< ok 0 0
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T2
: string
| C
: string
| CK_N
< row
: string
| T2
: string
| R
: string
| FK_P
< row
: string
| T2
: string
| C
: string
| NN_CODE
< row
: string
| T2
: string
| P
: string
| PK_T
< row
: string
| T2
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000024
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T2
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T2
: string
| C
: string
| CK_N
< row
: string
| T2
: string
| R
: string
| FK_P
< row
: string
| T2
: string
| C
: string
| NN_CODE
< row
: string
| T2
: string
| P
: string
| PK_T
< row
: string
| T2
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000024
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T2
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| ALTER TABLE t2 DROP CONSTRAINT ck_n
< ok 0 0
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T2
: string
| R
: string
| FK_P
< row
: string
| T2
: string
| C
: string
| NN_CODE
< row
: string
| T2
: string
| P
: string
| PK_T
< row
: string
| T2
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000024
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T2
: string
| 1
< row
: string
| a'b
: string
| 1
> query ^
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T2
: string
| 1
< row
: string
| a'b
: string
| 1

== TestScanRulesAgainstOracle/a_synced_table_the_change_recreates default#1
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
| SYS_C9000022
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
| SYS_C9000022
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
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000022
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| DROP TABLE t
< ok 0 0
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| a'b
: string
| P
: string
| SYS_C9000022
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| a'b
: string
| P
: string
| SYS_C9000022
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| CREATE TABLE t (id NUMBER CONSTRAINT ck_n CHECK (id > 0))
< ok 0 0
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| a'b
: string
| P
: string
| SYS_C9000022
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 0
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| a'b
: string
| P
: string
| SYS_C9000022
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 0
< row
: string
| a'b
: string
| 1
> exec
| ALTER TABLE t DROP CONSTRAINT ck_n
< ok 0 0
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| a'b
: string
| P
: string
| SYS_C9000022
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 0
< row
: string
| a'b
: string
| 1
> query ^
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 0
< row
: string
| a'b
: string
| 1

== TestScanRulesAgainstOracle/a_table_created_without_a_key default#1
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
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000007
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| CREATE TABLE a (id NUMBER)
< ok 0 0
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000007
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| A
: string
| 0
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000007
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| A
: string
| 0
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| CREATE TABLE b (id NUMBER PRIMARY KEY)
< ok 0 0
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| B
: string
| P
: string
| SYS_C9000008
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000007
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| A
: string
| 0
< row
: string
| B
: string
| 1
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| B
: string
| P
: string
| SYS_C9000008
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000007
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| A
: string
| 0
< row
: string
| B
: string
| 1
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| CREATE GLOBAL TEMPORARY TABLE g (id NUMBER)
< ok 0 0
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| B
: string
| P
: string
| SYS_C9000008
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000007
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| A
: string
| 0
< row
: string
| B
: string
| 1
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query ^
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| A
: string
| 0
< row
: string
| B
: string
| 1
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1

== TestScanRulesAgainstOracle/a_table_renamed_after_losing_its_key default#1
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
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000017
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| ALTER TABLE plain DROP PRIMARY KEY
< ok 0 0
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000017
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 0
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000017
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 0
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| RENAME plain TO plain2
< ok 0 0
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000017
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN2
: string
| 0
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query ^
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN2
: string
| 0
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1

== TestScanRulesAgainstOracle/a_table_the_change_recreates default#1
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
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000021
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| CREATE TABLE a (id NUMBER CONSTRAINT ck_a CHECK (id > 0))
< ok 0 0
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| A
: string
| C
: string
| CK_A
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000021
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| A
: string
| 0
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| A
: string
| C
: string
| CK_A
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000021
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| A
: string
| 0
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| DROP TABLE a
< ok 0 0
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000021
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000021
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| CREATE TABLE a (id NUMBER)
< ok 0 0
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000021
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| A
: string
| 0
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000021
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| A
: string
| 0
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| ALTER TABLE a DROP CONSTRAINT ck_a
< error
| 2443 ORA-02443: Cannot drop constraint  - nonexistent constraint
|
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| A
: string
| 0
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1

== TestScanRulesAgainstOracle/a_unique_index_that_backs_no_constraint default#1
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
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000006
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| ALTER TABLE plain DROP CONSTRAINT ux_v
< error
| 2443 ORA-02443: Cannot drop constraint  - nonexistent constraint
|
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1

== TestScanRulesAgainstOracle/constraints_of_the_synced_schema default#1
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
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000001
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| ALTER TABLE t DROP CONSTRAINT fk_p DROP CONSTRAINT ck_n
< ok 0 0
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000001
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000001
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| ALTER TABLE t DROP CONSTRAINT uk_code
< ok 0 0
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| a'b
: string
| P
: string
| SYS_C9000001
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query ^
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1

== TestScanRulesAgainstOracle/the_primary_key_by_name default#1
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
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| PLAIN
: string
| P
: string
| PK_PLAIN
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000003
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 1
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> exec
| ALTER TABLE plain DROP CONSTRAINT pk_plain
< ok 0 0
> query
| SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO' ORDER BY c.table_name, c.constraint_name
< result
< column "TABLE_NAME" "NCHAR"
< column "CONSTRAINT_TYPE" "NCHAR"
< column "CONSTRAINT_NAME" "NCHAR"
< row
: string
| LONGT
: string
| P
: string
| PK_LONGT
< row
: string
| P
: string
| P
: string
| PK_P
< row
: string
| T
: string
| C
: string
| CK_N
< row
: string
| T
: string
| R
: string
| FK_P
< row
: string
| T
: string
| C
: string
| NN_CODE
< row
: string
| T
: string
| P
: string
| PK_T
< row
: string
| T
: string
| U
: string
| UK_CODE
< row
: string
| a'b
: string
| P
: string
| SYS_C9000003
> query
| SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N' ORDER BY t.table_name
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 0
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
> query ^
< result
< column "TABLE_NAME" "NCHAR"
< column "(SELECTCOUNT(*)FROMUSER_CONSTRAINTSCWHEREC.TABLE_NAME=T.TABLE_NAMEANDC.CONSTRAINT_TYPE='P')" "NUMBER"
< row
: string
| LONGT
: string
| 1
< row
: string
| NOKEY
: string
| 0
< row
: string
| P
: string
| 1
< row
: string
| PLAIN
: string
| 0
< row
: string
| T
: string
| 1
< row
: string
| a'b
: string
| 1
