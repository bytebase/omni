-- pin: default gvenzl/oracle-free:23.26.3-slim-faststart sha256:f5ff19033860d662c821cb04eb10483fa94f14f78eae252d054291ea07028093 23.26.3.0.0

== TestQuerySpanFidelity default#1
> exec
| CREATE USER shop IDENTIFIED BY "omni_test" QUOTA UNLIMITED ON USERS
< ok 0 0
> exec
| CREATE USER people IDENTIFIED BY "omni_test" QUOTA UNLIMITED ON USERS
< ok 0 0
> exec
| GRANT CREATE SESSION, CREATE TABLE TO shop, people
< ok 0 0
> exec
| GRANT SELECT ANY TABLE TO shop, people
< ok 0 0
> exec
| CREATE TABLE shop.customers (id NUMBER, name VARCHAR2(50), email VARCHAR2(100), region VARCHAR2(20))
< ok 0 0
> exec
| CREATE TABLE shop.orders (id NUMBER, customer_id NUMBER, amount NUMBER(10,2), status VARCHAR2(10), created DATE)
< ok 0 0
> exec
| CREATE TABLE shop.items (order_id NUMBER, sku VARCHAR2(20), qty NUMBER)
< ok 0 0
> exec
| CREATE TABLE shop.emp (empno NUMBER, ename VARCHAR2(30), mgr NUMBER, sal NUMBER, deptno NUMBER)
< ok 0 0
> exec
| CREATE TABLE shop.sales (region VARCHAR2(20), quarter VARCHAR2(2), amount NUMBER)
< ok 0 0
> exec
| CREATE TABLE shop.quarterly (region VARCHAR2(20), q1 NUMBER, q2 NUMBER)
< ok 0 0
> exec
| CREATE TABLE shop.secrets (id NUMBER, token VARCHAR2(40) INVISIBLE, note VARCHAR2(40))
< ok 0 0
> exec
| CREATE TABLE shop."Mixed Case" ("Id" NUMBER, "Full Name" VARCHAR2(50), plain VARCHAR2(10))
< ok 0 0
> exec
| CREATE TABLE shop."lower" ("id" NUMBER, "val" VARCHAR2(10))
< ok 0 0
> exec
| CREATE TABLE people.employees (id NUMBER, name VARCHAR2(50), salary NUMBER, dept_id NUMBER)
< ok 0 0
> exec
| CREATE TABLE people.departments (id NUMBER, name VARCHAR2(50))
< ok 0 0
> exec
| CREATE TABLE people.customers (id NUMBER, name VARCHAR2(50))
< ok 0 0
> exec
| CREATE VIEW shop.order_totals AS SELECT o.id, o.customer_id, c.name AS customer_name, o.amount * 2 AS gross FROM shop.orders o JOIN shop.customers c ON c.id = o.customer_id
< ok 0 0
> exec
| CREATE VIEW shop.big_orders AS SELECT id, customer_name, gross FROM shop.order_totals WHERE gross > 100
< ok 0 0
> exec
| CREATE VIEW shop.all_people AS SELECT id, name FROM shop.customers UNION ALL SELECT id, name FROM people.employees
< ok 0 0
> exec
| CREATE VIEW shop."Vip Names" AS SELECT name AS "VipName", region FROM shop.customers WHERE region = 'vip'
< ok 0 0
> exec
| CREATE VIEW people.staff AS SELECT e.id AS staff_id, e.name AS staff_name, d.name AS dept FROM people.employees e JOIN people.departments d ON d.id = e.dept_id
< ok 0 0
> exec
| CREATE MATERIALIZED VIEW shop.region_counts AS SELECT region, COUNT(*) AS n FROM shop.customers GROUP BY region
< ok 0 0
> exec
| CREATE SYNONYM shop.syn_orders FOR shop.orders
< ok 0 0
> exec
| CREATE SYNONYM shop.staff_syn FOR people.staff
< ok 0 0
> exec
| CREATE PUBLIC SYNONYM omni_pub_emp FOR people.employees
< ok 0 0
> exec
| CREATE PUBLIC DATABASE LINK omni_loop CONNECT TO people IDENTIFIED BY "omni_test" USING 'localhost:1521/FREEPDB1'
< ok 0 0
> exec
| CREATE SYNONYM shop.remote_emp FOR people.employees@omni_loop
< ok 0 0
> query
| SELECT username FROM all_users WHERE oracle_maintained = 'N' ORDER BY username
< result
< column "USERNAME" "NCHAR" "string"
< row
: string
| PDBADMIN
< row
: string
| PEOPLE
< row
: string
| SHOP
> query
| SELECT c.table_name, c.column_name, c.hidden_column
| FROM all_tab_cols c JOIN all_tables t ON t.owner = c.owner AND t.table_name = c.table_name
| WHERE c.owner = 'PDBADMIN' AND c.user_generated = 'YES'
| AND NOT EXISTS (SELECT 1 FROM all_mviews m WHERE m.owner = c.owner AND m.mview_name = c.table_name)
| ORDER BY c.table_name, c.column_id NULLS LAST, c.column_name
< result
< column "TABLE_NAME" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
> query
| SELECT view_name, text FROM all_views WHERE owner = 'PDBADMIN' ORDER BY view_name
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT" "LONG"
> query
| SELECT mview_name, query FROM all_mviews WHERE owner = 'PDBADMIN' ORDER BY mview_name
< result
< column "MVIEW_NAME" "NCHAR" "string"
< column "QUERY" "LONG"
> query
| SELECT synonym_name, table_owner, table_name, db_link FROM dba_synonyms WHERE owner = 'PDBADMIN' ORDER BY synonym_name
< result
< column "SYNONYM_NAME" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "DB_LINK" "NCHAR" "string"
> query
| SELECT c.table_name, c.column_name, c.hidden_column
| FROM all_tab_cols c JOIN all_tables t ON t.owner = c.owner AND t.table_name = c.table_name
| WHERE c.owner = 'PEOPLE' AND c.user_generated = 'YES'
| AND NOT EXISTS (SELECT 1 FROM all_mviews m WHERE m.owner = c.owner AND m.mview_name = c.table_name)
| ORDER BY c.table_name, c.column_id NULLS LAST, c.column_name
< result
< column "TABLE_NAME" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| CUSTOMERS
: string
| ID
: string
| NO
< row
: string
| CUSTOMERS
: string
| NAME
: string
| NO
< row
: string
| DEPARTMENTS
: string
| ID
: string
| NO
< row
: string
| DEPARTMENTS
: string
| NAME
: string
| NO
< row
: string
| EMPLOYEES
: string
| ID
: string
| NO
< row
: string
| EMPLOYEES
: string
| NAME
: string
| NO
< row
: string
| EMPLOYEES
: string
| SALARY
: string
| NO
< row
: string
| EMPLOYEES
: string
| DEPT_ID
: string
| NO
> query
| SELECT view_name, text FROM all_views WHERE owner = 'PEOPLE' ORDER BY view_name
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT" "LONG"
< row
: string
| STAFF
: string
| SELECT e.id AS staff_id, e.name AS staff_name, d.name AS dept FROM people.employees e JOIN people.departments d ON d.id = e.dept_id
> query
| SELECT mview_name, query FROM all_mviews WHERE owner = 'PEOPLE' ORDER BY mview_name
< result
< column "MVIEW_NAME" "NCHAR" "string"
< column "QUERY" "LONG"
> query
| SELECT synonym_name, table_owner, table_name, db_link FROM dba_synonyms WHERE owner = 'PEOPLE' ORDER BY synonym_name
< result
< column "SYNONYM_NAME" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "DB_LINK" "NCHAR" "string"
> query
| SELECT c.table_name, c.column_name, c.hidden_column
| FROM all_tab_cols c JOIN all_tables t ON t.owner = c.owner AND t.table_name = c.table_name
| WHERE c.owner = 'SHOP' AND c.user_generated = 'YES'
| AND NOT EXISTS (SELECT 1 FROM all_mviews m WHERE m.owner = c.owner AND m.mview_name = c.table_name)
| ORDER BY c.table_name, c.column_id NULLS LAST, c.column_name
< result
< column "TABLE_NAME" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| CUSTOMERS
: string
| ID
: string
| NO
< row
: string
| CUSTOMERS
: string
| NAME
: string
| NO
< row
: string
| CUSTOMERS
: string
| EMAIL
: string
| NO
< row
: string
| CUSTOMERS
: string
| REGION
: string
| NO
< row
: string
| EMP
: string
| EMPNO
: string
| NO
< row
: string
| EMP
: string
| ENAME
: string
| NO
< row
: string
| EMP
: string
| MGR
: string
| NO
< row
: string
| EMP
: string
| SAL
: string
| NO
< row
: string
| EMP
: string
| DEPTNO
: string
| NO
< row
: string
| ITEMS
: string
| ORDER_ID
: string
| NO
< row
: string
| ITEMS
: string
| SKU
: string
| NO
< row
: string
| ITEMS
: string
| QTY
: string
| NO
< row
: string
| Mixed Case
: string
| Id
: string
| NO
< row
: string
| Mixed Case
: string
| Full Name
: string
| NO
< row
: string
| Mixed Case
: string
| PLAIN
: string
| NO
< row
: string
| ORDERS
: string
| ID
: string
| NO
< row
: string
| ORDERS
: string
| CUSTOMER_ID
: string
| NO
< row
: string
| ORDERS
: string
| AMOUNT
: string
| NO
< row
: string
| ORDERS
: string
| STATUS
: string
| NO
< row
: string
| ORDERS
: string
| CREATED
: string
| NO
< row
: string
| QUARTERLY
: string
| REGION
: string
| NO
< row
: string
| QUARTERLY
: string
| Q1
: string
| NO
< row
: string
| QUARTERLY
: string
| Q2
: string
| NO
< row
: string
| SALES
: string
| REGION
: string
| NO
< row
: string
| SALES
: string
| QUARTER
: string
| NO
< row
: string
| SALES
: string
| AMOUNT
: string
| NO
< row
: string
| SECRETS
: string
| ID
: string
| NO
< row
: string
| SECRETS
: string
| NOTE
: string
| NO
< row
: string
| SECRETS
: string
| TOKEN
: string
| YES
< row
: string
| lower
: string
| id
: string
| NO
< row
: string
| lower
: string
| val
: string
| NO
> query
| SELECT view_name, text FROM all_views WHERE owner = 'SHOP' ORDER BY view_name
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT" "LONG"
< row
: string
| ALL_PEOPLE
: string
| SELECT id, name FROM shop.customers UNION ALL SELECT id, name FROM people.employees
< row
: string
| BIG_ORDERS
: string
| SELECT id, customer_name, gross FROM shop.order_totals WHERE gross > 100
< row
: string
| ORDER_TOTALS
: string
| SELECT o.id, o.customer_id, c.name AS customer_name, o.amount * 2 AS gross FROM shop.orders o JOIN shop.customers c ON c.id = o.customer_id
< row
: string
| Vip Names
: string
| SELECT name AS "VipName", region FROM shop.customers WHERE region = 'vip'
> query
| SELECT mview_name, query FROM all_mviews WHERE owner = 'SHOP' ORDER BY mview_name
< result
< column "MVIEW_NAME" "NCHAR" "string"
< column "QUERY" "LONG"
< row
: string
| REGION_COUNTS
: string
| SELECT region, COUNT(*) AS n FROM shop.customers GROUP BY region
> query
| SELECT synonym_name, table_owner, table_name, db_link FROM dba_synonyms WHERE owner = 'SHOP' ORDER BY synonym_name
< result
< column "SYNONYM_NAME" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "DB_LINK" "NCHAR" "string"
< row
: string
| REMOTE_EMP
: string
| PEOPLE
: string
| EMPLOYEES
: string
| OMNI_LOOP
< row
: string
| STAFF_SYN
: string
| PEOPLE
: string
| STAFF
: null
< row
: string
| SYN_ORDERS
: string
| SHOP
: string
| ORDERS
: null
> query
| SELECT synonym_name, table_owner, table_name, db_link FROM dba_synonyms WHERE owner = 'PUBLIC' AND table_owner IN ('PDBADMIN', 'PEOPLE', 'SHOP') ORDER BY synonym_name
< result
< column "SYNONYM_NAME" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "DB_LINK" "NCHAR" "string"
< row
: string
| OMNI_PUB_EMP
: string
| PEOPLE
: string
| EMPLOYEES
: null
> query
| SELECT c.table_name, c.column_name, c.data_type, c.data_precision, c.data_scale, c.char_length, c.char_used
| FROM all_tab_cols c JOIN all_tables t ON t.owner = c.owner AND t.table_name = c.table_name
| WHERE c.owner = 'PDBADMIN' AND c.user_generated = 'YES'
| AND NOT EXISTS (SELECT 1 FROM all_mviews m WHERE m.owner = c.owner AND m.mview_name = c.table_name)
| ORDER BY c.table_name, c.column_name
< result
< column "TABLE_NAME" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
> query
| SELECT c.table_name, c.column_name, c.data_type, c.data_precision, c.data_scale, c.char_length, c.char_used
| FROM all_tab_cols c JOIN all_tables t ON t.owner = c.owner AND t.table_name = c.table_name
| WHERE c.owner = 'PEOPLE' AND c.user_generated = 'YES'
| AND NOT EXISTS (SELECT 1 FROM all_mviews m WHERE m.owner = c.owner AND m.mview_name = c.table_name)
| ORDER BY c.table_name, c.column_name
< result
< column "TABLE_NAME" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< row
: string
| CUSTOMERS
: string
| ID
: string
| NUMBER
: null
: null
: string
| 0
: null
< row
: string
| CUSTOMERS
: string
| NAME
: string
| VARCHAR2
: null
: null
: string
| 50
: string
| B
< row
: string
| DEPARTMENTS
: string
| ID
: string
| NUMBER
: null
: null
: string
| 0
: null
< row
: string
| DEPARTMENTS
: string
| NAME
: string
| VARCHAR2
: null
: null
: string
| 50
: string
| B
< row
: string
| EMPLOYEES
: string
| DEPT_ID
: string
| NUMBER
: null
: null
: string
| 0
: null
< row
: string
| EMPLOYEES
: string
| ID
: string
| NUMBER
: null
: null
: string
| 0
: null
< row
: string
| EMPLOYEES
: string
| NAME
: string
| VARCHAR2
: null
: null
: string
| 50
: string
| B
< row
: string
| EMPLOYEES
: string
| SALARY
: string
| NUMBER
: null
: null
: string
| 0
: null
> exec
| CREATE DOMAIN "PEOPLE"."OMNI_SPAN_D1" AS NUMBER
< ok 0 0
> exec
| ALTER TABLE "PEOPLE"."CUSTOMERS" MODIFY ("ID" DOMAIN "PEOPLE"."OMNI_SPAN_D1")
< ok 0 0
> exec
| CREATE DOMAIN "PEOPLE"."OMNI_SPAN_D2" AS VARCHAR2(50 BYTE)
< ok 0 0
> exec
| ALTER TABLE "PEOPLE"."CUSTOMERS" MODIFY ("NAME" DOMAIN "PEOPLE"."OMNI_SPAN_D2")
< ok 0 0
> exec
| CREATE DOMAIN "PEOPLE"."OMNI_SPAN_D3" AS NUMBER
< ok 0 0
> exec
| ALTER TABLE "PEOPLE"."DEPARTMENTS" MODIFY ("ID" DOMAIN "PEOPLE"."OMNI_SPAN_D3")
< ok 0 0
> exec
| CREATE DOMAIN "PEOPLE"."OMNI_SPAN_D4" AS VARCHAR2(50 BYTE)
< ok 0 0
> exec
| ALTER TABLE "PEOPLE"."DEPARTMENTS" MODIFY ("NAME" DOMAIN "PEOPLE"."OMNI_SPAN_D4")
< ok 0 0
> exec
| CREATE DOMAIN "PEOPLE"."OMNI_SPAN_D5" AS NUMBER
< ok 0 0
> exec
| ALTER TABLE "PEOPLE"."EMPLOYEES" MODIFY ("DEPT_ID" DOMAIN "PEOPLE"."OMNI_SPAN_D5")
< ok 0 0
> exec
| CREATE DOMAIN "PEOPLE"."OMNI_SPAN_D6" AS NUMBER
< ok 0 0
> exec
| ALTER TABLE "PEOPLE"."EMPLOYEES" MODIFY ("ID" DOMAIN "PEOPLE"."OMNI_SPAN_D6")
< ok 0 0
> exec
| CREATE DOMAIN "PEOPLE"."OMNI_SPAN_D7" AS VARCHAR2(50 BYTE)
< ok 0 0
> exec
| ALTER TABLE "PEOPLE"."EMPLOYEES" MODIFY ("NAME" DOMAIN "PEOPLE"."OMNI_SPAN_D7")
< ok 0 0
> exec
| CREATE DOMAIN "PEOPLE"."OMNI_SPAN_D8" AS NUMBER
< ok 0 0
> exec
| ALTER TABLE "PEOPLE"."EMPLOYEES" MODIFY ("SALARY" DOMAIN "PEOPLE"."OMNI_SPAN_D8")
< ok 0 0
> query
| SELECT c.table_name, c.column_name, c.data_type, c.data_precision, c.data_scale, c.char_length, c.char_used
| FROM all_tab_cols c JOIN all_tables t ON t.owner = c.owner AND t.table_name = c.table_name
| WHERE c.owner = 'SHOP' AND c.user_generated = 'YES'
| AND NOT EXISTS (SELECT 1 FROM all_mviews m WHERE m.owner = c.owner AND m.mview_name = c.table_name)
| ORDER BY c.table_name, c.column_name
< result
< column "TABLE_NAME" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< row
: string
| CUSTOMERS
: string
| EMAIL
: string
| VARCHAR2
: null
: null
: string
| 100
: string
| B
< row
: string
| CUSTOMERS
: string
| ID
: string
| NUMBER
: null
: null
: string
| 0
: null
< row
: string
| CUSTOMERS
: string
| NAME
: string
| VARCHAR2
: null
: null
: string
| 50
: string
| B
< row
: string
| CUSTOMERS
: string
| REGION
: string
| VARCHAR2
: null
: null
: string
| 20
: string
| B
< row
: string
| EMP
: string
| DEPTNO
: string
| NUMBER
: null
: null
: string
| 0
: null
< row
: string
| EMP
: string
| EMPNO
: string
| NUMBER
: null
: null
: string
| 0
: null
< row
: string
| EMP
: string
| ENAME
: string
| VARCHAR2
: null
: null
: string
| 30
: string
| B
< row
: string
| EMP
: string
| MGR
: string
| NUMBER
: null
: null
: string
| 0
: null
< row
: string
| EMP
: string
| SAL
: string
| NUMBER
: null
: null
: string
| 0
: null
< row
: string
| ITEMS
: string
| ORDER_ID
: string
| NUMBER
: null
: null
: string
| 0
: null
< row
: string
| ITEMS
: string
| QTY
: string
| NUMBER
: null
: null
: string
| 0
: null
< row
: string
| ITEMS
: string
| SKU
: string
| VARCHAR2
: null
: null
: string
| 20
: string
| B
< row
: string
| Mixed Case
: string
| Full Name
: string
| VARCHAR2
: null
: null
: string
| 50
: string
| B
< row
: string
| Mixed Case
: string
| Id
: string
| NUMBER
: null
: null
: string
| 0
: null
< row
: string
| Mixed Case
: string
| PLAIN
: string
| VARCHAR2
: null
: null
: string
| 10
: string
| B
< row
: string
| ORDERS
: string
| AMOUNT
: string
| NUMBER
: string
| 10
: string
| 2
: string
| 0
: null
< row
: string
| ORDERS
: string
| CREATED
: string
| DATE
: null
: null
: string
| 0
: null
< row
: string
| ORDERS
: string
| CUSTOMER_ID
: string
| NUMBER
: null
: null
: string
| 0
: null
< row
: string
| ORDERS
: string
| ID
: string
| NUMBER
: null
: null
: string
| 0
: null
< row
: string
| ORDERS
: string
| STATUS
: string
| VARCHAR2
: null
: null
: string
| 10
: string
| B
< row
: string
| QUARTERLY
: string
| Q1
: string
| NUMBER
: null
: null
: string
| 0
: null
< row
: string
| QUARTERLY
: string
| Q2
: string
| NUMBER
: null
: null
: string
| 0
: null
< row
: string
| QUARTERLY
: string
| REGION
: string
| VARCHAR2
: null
: null
: string
| 20
: string
| B
< row
: string
| SALES
: string
| AMOUNT
: string
| NUMBER
: null
: null
: string
| 0
: null
< row
: string
| SALES
: string
| QUARTER
: string
| VARCHAR2
: null
: null
: string
| 2
: string
| B
< row
: string
| SALES
: string
| REGION
: string
| VARCHAR2
: null
: null
: string
| 20
: string
| B
< row
: string
| SECRETS
: string
| ID
: string
| NUMBER
: null
: null
: string
| 0
: null
< row
: string
| SECRETS
: string
| NOTE
: string
| VARCHAR2
: null
: null
: string
| 40
: string
| B
< row
: string
| SECRETS
: string
| TOKEN
: string
| VARCHAR2
: null
: null
: string
| 40
: string
| B
< row
: string
| lower
: string
| id
: string
| NUMBER
: null
: null
: string
| 0
: null
< row
: string
| lower
: string
| val
: string
| VARCHAR2
: null
: null
: string
| 10
: string
| B
> exec
| CREATE DOMAIN "SHOP"."OMNI_SPAN_D9" AS VARCHAR2(100 BYTE)
< ok 0 0
> exec
| ALTER TABLE "SHOP"."CUSTOMERS" MODIFY ("EMAIL" DOMAIN "SHOP"."OMNI_SPAN_D9")
< ok 0 0
> exec
| CREATE DOMAIN "SHOP"."OMNI_SPAN_D10" AS NUMBER
< ok 0 0
> exec
| ALTER TABLE "SHOP"."CUSTOMERS" MODIFY ("ID" DOMAIN "SHOP"."OMNI_SPAN_D10")
< ok 0 0
> exec
| CREATE DOMAIN "SHOP"."OMNI_SPAN_D11" AS VARCHAR2(50 BYTE)
< ok 0 0
> exec
| ALTER TABLE "SHOP"."CUSTOMERS" MODIFY ("NAME" DOMAIN "SHOP"."OMNI_SPAN_D11")
< ok 0 0
> exec
| CREATE DOMAIN "SHOP"."OMNI_SPAN_D12" AS VARCHAR2(20 BYTE)
< ok 0 0
> exec
| ALTER TABLE "SHOP"."CUSTOMERS" MODIFY ("REGION" DOMAIN "SHOP"."OMNI_SPAN_D12")
< ok 0 0
> exec
| CREATE DOMAIN "SHOP"."OMNI_SPAN_D13" AS NUMBER
< ok 0 0
> exec
| ALTER TABLE "SHOP"."EMP" MODIFY ("DEPTNO" DOMAIN "SHOP"."OMNI_SPAN_D13")
< ok 0 0
> exec
| CREATE DOMAIN "SHOP"."OMNI_SPAN_D14" AS NUMBER
< ok 0 0
> exec
| ALTER TABLE "SHOP"."EMP" MODIFY ("EMPNO" DOMAIN "SHOP"."OMNI_SPAN_D14")
< ok 0 0
> exec
| CREATE DOMAIN "SHOP"."OMNI_SPAN_D15" AS VARCHAR2(30 BYTE)
< ok 0 0
> exec
| ALTER TABLE "SHOP"."EMP" MODIFY ("ENAME" DOMAIN "SHOP"."OMNI_SPAN_D15")
< ok 0 0
> exec
| CREATE DOMAIN "SHOP"."OMNI_SPAN_D16" AS NUMBER
< ok 0 0
> exec
| ALTER TABLE "SHOP"."EMP" MODIFY ("MGR" DOMAIN "SHOP"."OMNI_SPAN_D16")
< ok 0 0
> exec
| CREATE DOMAIN "SHOP"."OMNI_SPAN_D17" AS NUMBER
< ok 0 0
> exec
| ALTER TABLE "SHOP"."EMP" MODIFY ("SAL" DOMAIN "SHOP"."OMNI_SPAN_D17")
< ok 0 0
> exec
| CREATE DOMAIN "SHOP"."OMNI_SPAN_D18" AS NUMBER
< ok 0 0
> exec
| ALTER TABLE "SHOP"."ITEMS" MODIFY ("ORDER_ID" DOMAIN "SHOP"."OMNI_SPAN_D18")
< ok 0 0
> exec
| CREATE DOMAIN "SHOP"."OMNI_SPAN_D19" AS NUMBER
< ok 0 0
> exec
| ALTER TABLE "SHOP"."ITEMS" MODIFY ("QTY" DOMAIN "SHOP"."OMNI_SPAN_D19")
< ok 0 0
> exec
| CREATE DOMAIN "SHOP"."OMNI_SPAN_D20" AS VARCHAR2(20 BYTE)
< ok 0 0
> exec
| ALTER TABLE "SHOP"."ITEMS" MODIFY ("SKU" DOMAIN "SHOP"."OMNI_SPAN_D20")
< ok 0 0
> exec
| CREATE DOMAIN "SHOP"."OMNI_SPAN_D21" AS VARCHAR2(50 BYTE)
< ok 0 0
> exec
| ALTER TABLE "SHOP"."Mixed Case" MODIFY ("Full Name" DOMAIN "SHOP"."OMNI_SPAN_D21")
< ok 0 0
> exec
| CREATE DOMAIN "SHOP"."OMNI_SPAN_D22" AS NUMBER
< ok 0 0
> exec
| ALTER TABLE "SHOP"."Mixed Case" MODIFY ("Id" DOMAIN "SHOP"."OMNI_SPAN_D22")
< ok 0 0
> exec
| CREATE DOMAIN "SHOP"."OMNI_SPAN_D23" AS VARCHAR2(10 BYTE)
< ok 0 0
> exec
| ALTER TABLE "SHOP"."Mixed Case" MODIFY ("PLAIN" DOMAIN "SHOP"."OMNI_SPAN_D23")
< ok 0 0
> exec
| CREATE DOMAIN "SHOP"."OMNI_SPAN_D24" AS NUMBER(10,2)
< ok 0 0
> exec
| ALTER TABLE "SHOP"."ORDERS" MODIFY ("AMOUNT" DOMAIN "SHOP"."OMNI_SPAN_D24")
< ok 0 0
> exec
| CREATE DOMAIN "SHOP"."OMNI_SPAN_D25" AS DATE
< ok 0 0
> exec
| ALTER TABLE "SHOP"."ORDERS" MODIFY ("CREATED" DOMAIN "SHOP"."OMNI_SPAN_D25")
< ok 0 0
> exec
| CREATE DOMAIN "SHOP"."OMNI_SPAN_D26" AS NUMBER
< ok 0 0
> exec
| ALTER TABLE "SHOP"."ORDERS" MODIFY ("CUSTOMER_ID" DOMAIN "SHOP"."OMNI_SPAN_D26")
< ok 0 0
> exec
| CREATE DOMAIN "SHOP"."OMNI_SPAN_D27" AS NUMBER
< ok 0 0
> exec
| ALTER TABLE "SHOP"."ORDERS" MODIFY ("ID" DOMAIN "SHOP"."OMNI_SPAN_D27")
< ok 0 0
> exec
| CREATE DOMAIN "SHOP"."OMNI_SPAN_D28" AS VARCHAR2(10 BYTE)
< ok 0 0
> exec
| ALTER TABLE "SHOP"."ORDERS" MODIFY ("STATUS" DOMAIN "SHOP"."OMNI_SPAN_D28")
< ok 0 0
> exec
| CREATE DOMAIN "SHOP"."OMNI_SPAN_D29" AS NUMBER
< ok 0 0
> exec
| ALTER TABLE "SHOP"."QUARTERLY" MODIFY ("Q1" DOMAIN "SHOP"."OMNI_SPAN_D29")
< ok 0 0
> exec
| CREATE DOMAIN "SHOP"."OMNI_SPAN_D30" AS NUMBER
< ok 0 0
> exec
| ALTER TABLE "SHOP"."QUARTERLY" MODIFY ("Q2" DOMAIN "SHOP"."OMNI_SPAN_D30")
< ok 0 0
> exec
| CREATE DOMAIN "SHOP"."OMNI_SPAN_D31" AS VARCHAR2(20 BYTE)
< ok 0 0
> exec
| ALTER TABLE "SHOP"."QUARTERLY" MODIFY ("REGION" DOMAIN "SHOP"."OMNI_SPAN_D31")
< ok 0 0
> exec
| CREATE DOMAIN "SHOP"."OMNI_SPAN_D32" AS NUMBER
< ok 0 0
> exec
| ALTER TABLE "SHOP"."SALES" MODIFY ("AMOUNT" DOMAIN "SHOP"."OMNI_SPAN_D32")
< ok 0 0
> exec
| CREATE DOMAIN "SHOP"."OMNI_SPAN_D33" AS VARCHAR2(2 BYTE)
< ok 0 0
> exec
| ALTER TABLE "SHOP"."SALES" MODIFY ("QUARTER" DOMAIN "SHOP"."OMNI_SPAN_D33")
< ok 0 0
> exec
| CREATE DOMAIN "SHOP"."OMNI_SPAN_D34" AS VARCHAR2(20 BYTE)
< ok 0 0
> exec
| ALTER TABLE "SHOP"."SALES" MODIFY ("REGION" DOMAIN "SHOP"."OMNI_SPAN_D34")
< ok 0 0
> exec
| CREATE DOMAIN "SHOP"."OMNI_SPAN_D35" AS NUMBER
< ok 0 0
> exec
| ALTER TABLE "SHOP"."SECRETS" MODIFY ("ID" DOMAIN "SHOP"."OMNI_SPAN_D35")
< ok 0 0
> exec
| CREATE DOMAIN "SHOP"."OMNI_SPAN_D36" AS VARCHAR2(40 BYTE)
< ok 0 0
> exec
| ALTER TABLE "SHOP"."SECRETS" MODIFY ("NOTE" DOMAIN "SHOP"."OMNI_SPAN_D36")
< ok 0 0
> exec
| CREATE DOMAIN "SHOP"."OMNI_SPAN_D37" AS VARCHAR2(40 BYTE)
< ok 0 0
> exec
| ALTER TABLE "SHOP"."SECRETS" MODIFY ("TOKEN" DOMAIN "SHOP"."OMNI_SPAN_D37")
< ok 0 0
> exec
| CREATE DOMAIN "SHOP"."OMNI_SPAN_D38" AS NUMBER
< ok 0 0
> exec
| ALTER TABLE "SHOP"."lower" MODIFY ("id" DOMAIN "SHOP"."OMNI_SPAN_D38")
< ok 0 0
> exec
| CREATE DOMAIN "SHOP"."OMNI_SPAN_D39" AS VARCHAR2(10 BYTE)
< ok 0 0
> exec
| ALTER TABLE "SHOP"."lower" MODIFY ("val" DOMAIN "SHOP"."OMNI_SPAN_D39")
< ok 0 0
> exec
| BEGIN DBMS_UTILITY.COMPILE_SCHEMA('PDBADMIN', FALSE); END;
< ok 1 0
> exec
| BEGIN DBMS_UTILITY.COMPILE_SCHEMA('PEOPLE', FALSE); END;
< ok 1 0
> exec
| BEGIN DBMS_UTILITY.COMPILE_SCHEMA('SHOP', FALSE); END;
< ok 1 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT id, amount FROM orders
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "AMOUNT" "NUMBER" "float64" precision 10 scale 2
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT id, amount FROM orders
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D27
< row
: string
| SHOP
: string
| OMNI_SPAN_D24
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT id, amount FROM orders
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT o.id, o.amount AS total FROM orders o
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "TOTAL" "NUMBER" "float64" precision 10 scale 2
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT o.id, o.amount AS total FROM orders o
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D27
< row
: string
| SHOP
: string
| OMNI_SPAN_D24
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT o.id, o.amount AS total FROM orders o
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT * FROM customers
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "NAME" "NCHAR" "string"
< column "EMAIL" "NCHAR" "string"
< column "REGION" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2, C3, C4) AS SELECT * FROM customers
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D10
< row
: string
| SHOP
: string
| OMNI_SPAN_D11
< row
: string
| SHOP
: string
| OMNI_SPAN_D9
< row
: string
| SHOP
: string
| OMNI_SPAN_D12
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT * FROM customers
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT c.* FROM customers c
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "NAME" "NCHAR" "string"
< column "EMAIL" "NCHAR" "string"
< column "REGION" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2, C3, C4) AS SELECT c.* FROM customers c
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D10
< row
: string
| SHOP
: string
| OMNI_SPAN_D11
< row
: string
| SHOP
: string
| OMNI_SPAN_D9
< row
: string
| SHOP
: string
| OMNI_SPAN_D12
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT c.* FROM customers c
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT customers.name FROM customers
< result
< column "NAME" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1) AS SELECT customers.name FROM customers
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D11
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT customers.name FROM customers
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT shop.customers.email FROM shop.customers
< result
< column "EMAIL" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1) AS SELECT shop.customers.email FROM shop.customers
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D9
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT shop.customers.email FROM shop.customers
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT o.*, c.name FROM orders o, customers c WHERE o.customer_id = c.id
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "CUSTOMER_ID" "NUMBER" "float64" precision 38 scale 255
< column "AMOUNT" "NUMBER" "float64" precision 10 scale 2
< column "STATUS" "NCHAR" "string"
< column "CREATED" "DATE" "time.Time"
< column "NAME" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2, C3, C4, C5, C6) AS SELECT o.*, c.name FROM orders o, customers c WHERE o.customer_id = c.id
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D27
< row
: string
| SHOP
: string
| OMNI_SPAN_D26
< row
: string
| SHOP
: string
| OMNI_SPAN_D24
< row
: string
| SHOP
: string
| OMNI_SPAN_D28
< row
: string
| SHOP
: string
| OMNI_SPAN_D25
< row
: string
| SHOP
: string
| OMNI_SPAN_D11
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT o.*, c.name FROM orders o, customers c WHERE o.customer_id = c.id
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| HASH JOIN
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 4 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT id AS "MixedAlias", name "lower alias" FROM customers
< result
< column "MixedAlias" "NUMBER" "float64" precision 38 scale 255
< column "lower alias" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT id AS "MixedAlias", name "lower alias" FROM customers
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D10
< row
: string
| SHOP
: string
| OMNI_SPAN_D11
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT id AS "MixedAlias", name "lower alias" FROM customers
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT id, id FROM orders
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "ID" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT id, id FROM orders
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D27
< row
: string
| SHOP
: string
| OMNI_SPAN_D27
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT id, id FROM orders
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT orders.id FROM customers orders
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1) AS SELECT orders.id FROM customers orders
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D10
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT orders.id FROM customers orders
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT id + 1, amount * 2 AS doubled FROM orders
< result
< column "ID+1" "NUMBER" "float64" precision 38 scale 255
< column "DOUBLED" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT id + 1, amount * 2 AS doubled FROM orders
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: null
: null
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT id + 1, amount * 2 AS doubled FROM orders
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT name || ' <' || email || '>' AS contact FROM customers
< result
< column "CONTACT" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1) AS SELECT name || ' <' || email || '>' AS contact FROM customers
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT name || ' <' || email || '>' AS contact FROM customers
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT UPPER(name), NVL(email, 'none') AS email FROM customers
< result
< column "UPPER(NAME)" "NCHAR" "string"
< column "EMAIL" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT UPPER(name), NVL(email, 'none') AS email FROM customers
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: null
: null
< row
: string
| SHOP
: string
| OMNI_SPAN_D9
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT UPPER(name), NVL(email, 'none') AS email FROM customers
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT NVL(email, name) AS reach FROM customers
< result
< column "REACH" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1) AS SELECT NVL(email, name) AS reach FROM customers
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D9
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT NVL(email, name) AS reach FROM customers
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT CASE WHEN amount > 100 THEN 'big' ELSE status END AS kind FROM orders
< result
< column "KIND" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1) AS SELECT CASE WHEN amount > 100 THEN 'big' ELSE status END AS kind FROM orders
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT CASE WHEN amount > 100 THEN 'big' ELSE status END AS kind FROM orders
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT CASE WHEN status = 'paid' THEN amount END AS paid FROM orders
< result
< column "PAID" "NUMBER" "float64" precision 10 scale 2
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1) AS SELECT CASE WHEN status = 'paid' THEN amount END AS paid FROM orders
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D24
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT CASE WHEN status = 'paid' THEN amount END AS paid FROM orders
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT DECODE(status, 'paid', amount, 0) AS paid FROM orders
< result
< column "PAID" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1) AS SELECT DECODE(status, 'paid', amount, 0) AS paid FROM orders
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D24
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT DECODE(status, 'paid', amount, 0) AS paid FROM orders
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT COALESCE(email, region, name) AS c FROM customers
< result
< column "C" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1) AS SELECT COALESCE(email, region, name) AS c FROM customers
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D9
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT COALESCE(email, region, name) AS c FROM customers
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT CAST(amount AS NUMBER(12,2)) AS a, TO_CHAR(created, 'YYYY') AS y FROM orders
< result
< column "A" "NUMBER" "float64" precision 12 scale 2
< column "Y" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT CAST(amount AS NUMBER(12,2)) AS a, TO_CHAR(created, 'YYYY') AS y FROM orders
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: null
: null
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT CAST(amount AS NUMBER(12,2)) AS a, TO_CHAR(created, 'YYYY') AS y FROM orders
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT 1 AS one, 'x' AS lit, NULL AS n FROM orders
< result
< column "ONE" "NUMBER" "float64" precision 38 scale 255
< column "LIT" "CHAR" "string"
< column "N" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2, C3) AS SELECT 1 AS one, 'x' AS lit, NULL AS n FROM orders
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: null
: null
< row
: null
: null
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT 1 AS one, 'x' AS lit, NULL AS n FROM orders
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT 'abc', 'a B', n'x', 1e3 FROM orders
< result
< column "'ABC'" "CHAR" "string"
< column "'AB'" "CHAR" "string"
< column "N'X'" "CHAR" "string"
< column "1E3" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2, C3, C4) AS SELECT 'abc', 'a B', n'x', 1e3 FROM orders
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: null
: null
< row
: null
: null
< row
: null
: null
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT 'abc', 'a B', n'x', 1e3 FROM orders
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT "Id" + 1, LOWER(plain), "Mixed Case".plain || 'Z', "Full Name" || 'y' FROM "Mixed Case"
< result
< column "\"ID\"+1" "NUMBER" "float64" precision 38 scale 255
< column "LOWER(PLAIN)" "NCHAR" "string"
< column "\"MIXEDCASE\".PLAIN||'Z'" "NCHAR" "string"
< column "\"FULLNAME\"||'Y'" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2, C3, C4) AS SELECT "Id" + 1, LOWER(plain), "Mixed Case".plain || 'Z', "Full Name" || 'y' FROM "Mixed Case"
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: null
: null
< row
: null
: null
< row
: null
: null
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| Mixed Case
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT "Id" + 1, LOWER(plain), "Mixed Case".plain || 'Z', "Full Name" || 'y' FROM "Mixed Case"
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| Mixed Case
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT amount /* note */ + 1, -amount, q'[it's]' FROM orders
< result
< column "AMOUNT/*NOTE*/+1" "NUMBER" "float64" precision 38 scale 255
< column "-AMOUNT" "NUMBER" "float64" precision 38 scale 255
< column "Q'[IT'S]'" "CHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2, C3) AS SELECT amount /* note */ + 1, -amount, q'[it's]' FROM orders
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: null
: null
< row
: null
: null
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT amount /* note */ + 1, -amount, q'[it's]' FROM orders
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT id, (SELECT name FROM customers c WHERE c.id = o.customer_id) AS who FROM orders o
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "WHO" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT id, (SELECT name FROM customers c WHERE c.id = o.customer_id) AS who FROM orders o
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D27
< row
: string
| SHOP
: string
| OMNI_SPAN_D11
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT id, (SELECT name FROM customers c WHERE c.id = o.customer_id) AS who FROM orders o
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 3 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT (SELECT MAX(amount) FROM orders) AS top FROM customers
< result
< column "TOP" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1) AS SELECT (SELECT MAX(amount) FROM orders) AS top FROM customers
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT (SELECT MAX(amount) FROM orders) AS top FROM customers
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| SORT
: string
| AGGREGATE
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 4 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT ROWNUM, id FROM orders
< result
< column "ROWNUM" "NUMBER" "float64" precision 38 scale 255
< column "ID" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT ROWNUM, id FROM orders
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: null
: null
< row
: string
| SHOP
: string
| OMNI_SPAN_D27
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT ROWNUM, id FROM orders
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| COUNT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 3 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT ROWID, name FROM customers
< result
< column "ROWID" "ROWID" "string"
< column "NAME" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT ROWID, name FROM customers
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: null
: null
< row
: string
| SHOP
: string
| OMNI_SPAN_D11
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT ROWID, name FROM customers
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT ROWNUM AS rn, o.* FROM orders o WHERE ROWNUM <= 10
< result
< column "RN" "NUMBER" "float64" precision 38 scale 255
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "CUSTOMER_ID" "NUMBER" "float64" precision 38 scale 255
< column "AMOUNT" "NUMBER" "float64" precision 10 scale 2
< column "STATUS" "NCHAR" "string"
< column "CREATED" "DATE" "time.Time"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2, C3, C4, C5, C6) AS SELECT ROWNUM AS rn, o.* FROM orders o WHERE ROWNUM <= 10
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: null
: null
< row
: string
| SHOP
: string
| OMNI_SPAN_D27
< row
: string
| SHOP
: string
| OMNI_SPAN_D26
< row
: string
| SHOP
: string
| OMNI_SPAN_D24
< row
: string
| SHOP
: string
| OMNI_SPAN_D28
< row
: string
| SHOP
: string
| OMNI_SPAN_D25
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT ROWNUM AS rn, o.* FROM orders o WHERE ROWNUM <= 10
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| COUNT
: string
| STOPKEY
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 3 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT * FROM (SELECT id, amount FROM orders ORDER BY amount DESC) WHERE ROWNUM <= 3
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "AMOUNT" "NUMBER" "float64" precision 10 scale 2
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT * FROM (SELECT id, amount FROM orders ORDER BY amount DESC) WHERE ROWNUM <= 3
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D27
< row
: string
| SHOP
: string
| OMNI_SPAN_D24
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT * FROM (SELECT id, amount FROM orders ORDER BY amount DESC) WHERE ROWNUM <= 3
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| COUNT
: string
| STOPKEY
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| VIEW
: null
: string
| SHOP
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| SORT
: string
| ORDER BY STOPKEY
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 5 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT ORA_ROWSCN, id FROM orders
< result
< column "ORA_ROWSCN" "NUMBER" "float64" precision 38 scale 255
< column "ID" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT ORA_ROWSCN, id FROM orders
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: null
: null
< row
: string
| SHOP
: string
| OMNI_SPAN_D27
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT ORA_ROWSCN, id FROM orders
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT USER, SYSDATE - created AS age FROM orders
< result
< column "USER" "NCHAR" "string"
< column "AGE" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT USER, SYSDATE - created AS age FROM orders
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: null
: null
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT USER, SYSDATE - created AS age FROM orders
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT 1 FROM DUAL
< result
< column "1" "NUMBER" "float64" precision 38 scale 255
< row
: string
| 1
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1) AS SELECT 1 FROM DUAL
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| PUBLIC
: string
| DUAL
: string
| SYNONYM
: null
> query
| SELECT table_owner, table_name, db_link FROM dba_synonyms WHERE owner = 'PUBLIC' AND synonym_name = 'DUAL'
< result
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "DB_LINK" "NCHAR" "string"
< row
: string
| SYS
: string
| DUAL
: null
> query
| SELECT object_type FROM all_objects WHERE owner = 'SYS' AND object_name = 'DUAL' AND object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW', 'SYNONYM') ORDER BY object_type
< result
< column "OBJECT_TYPE" "NCHAR" "string"
< row
: string
| TABLE
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT 1 FROM DUAL
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| FAST DUAL
: null
: null
: null
: null
: null
: null
: null
: string
| 0
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT 1 AS x, 'a' AS y FROM DUAL
< result
< column "X" "NUMBER" "float64" precision 38 scale 255
< column "Y" "CHAR" "string"
< row
: string
| 1
: string
| a
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT 1 AS x, 'a' AS y FROM DUAL
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: null
: null
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| PUBLIC
: string
| DUAL
: string
| SYNONYM
: null
> query
| SELECT table_owner, table_name, db_link FROM dba_synonyms WHERE owner = 'PUBLIC' AND synonym_name = 'DUAL'
< result
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "DB_LINK" "NCHAR" "string"
< row
: string
| SYS
: string
| DUAL
: null
> query
| SELECT object_type FROM all_objects WHERE owner = 'SYS' AND object_name = 'DUAL' AND object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW', 'SYNONYM') ORDER BY object_type
< result
< column "OBJECT_TYPE" "NCHAR" "string"
< row
: string
| TABLE
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT 1 AS x, 'a' AS y FROM DUAL
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| FAST DUAL
: null
: null
: null
: null
: null
: null
: null
: string
| 0
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT dummy FROM dual
< result
< column "DUMMY" "NCHAR" "string"
< row
: string
| X
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1) AS SELECT dummy FROM dual
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| PUBLIC
: string
| DUAL
: string
| SYNONYM
: null
> query
| SELECT table_owner, table_name, db_link FROM dba_synonyms WHERE owner = 'PUBLIC' AND synonym_name = 'DUAL'
< result
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "DB_LINK" "NCHAR" "string"
< row
: string
| SYS
: string
| DUAL
: null
> query
| SELECT object_type FROM all_objects WHERE owner = 'SYS' AND object_name = 'DUAL' AND object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW', 'SYNONYM') ORDER BY object_type
< result
< column "OBJECT_TYPE" "NCHAR" "string"
< row
: string
| TABLE
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT dummy FROM dual
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SYS
: string
| DUAL
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT (SELECT COUNT(*) FROM orders) AS n FROM DUAL
< result
< column "N" "NUMBER" "float64" precision 38 scale 255
< row
: string
| 0
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1) AS SELECT (SELECT COUNT(*) FROM orders) AS n FROM DUAL
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| PUBLIC
: string
| DUAL
: string
| SYNONYM
: null
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> query
| SELECT table_owner, table_name, db_link FROM dba_synonyms WHERE owner = 'PUBLIC' AND synonym_name = 'DUAL'
< result
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "DB_LINK" "NCHAR" "string"
< row
: string
| SYS
: string
| DUAL
: null
> query
| SELECT object_type FROM all_objects WHERE owner = 'SYS' AND object_name = 'DUAL' AND object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW', 'SYNONYM') ORDER BY object_type
< result
< column "OBJECT_TYPE" "NCHAR" "string"
< row
: string
| TABLE
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT (SELECT COUNT(*) FROM orders) AS n FROM DUAL
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| SORT
: string
| AGGREGATE
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| FAST DUAL
: null
: null
: null
: null
: null
: null
: null
: string
| 0
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 4 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT o.id FROM orders o, dual
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1) AS SELECT o.id FROM orders o, dual
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D27
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| PUBLIC
: string
| DUAL
: string
| SYNONYM
: null
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> query
| SELECT table_owner, table_name, db_link FROM dba_synonyms WHERE owner = 'PUBLIC' AND synonym_name = 'DUAL'
< result
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "DB_LINK" "NCHAR" "string"
< row
: string
| SYS
: string
| DUAL
: null
> query
| SELECT object_type FROM all_objects WHERE owner = 'SYS' AND object_name = 'DUAL' AND object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW', 'SYNONYM') ORDER BY object_type
< result
< column "OBJECT_TYPE" "NCHAR" "string"
< row
: string
| TABLE
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT o.id FROM orders o, dual
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| MERGE JOIN
: string
| CARTESIAN
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| BUFFER
: string
| SORT
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| FAST DUAL
: null
: null
: null
: null
: null
: null
: null
: string
| 0
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 5 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT o.id, c.name FROM orders o JOIN customers c ON c.id = o.customer_id
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "NAME" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT o.id, c.name FROM orders o JOIN customers c ON c.id = o.customer_id
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D27
< row
: string
| SHOP
: string
| OMNI_SPAN_D11
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT o.id, c.name FROM orders o JOIN customers c ON c.id = o.customer_id
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| HASH JOIN
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 4 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT o.id, c.name FROM orders o LEFT JOIN customers c ON c.id = o.customer_id
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "NAME" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT o.id, c.name FROM orders o LEFT JOIN customers c ON c.id = o.customer_id
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D27
< row
: string
| SHOP
: string
| OMNI_SPAN_D11
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT o.id, c.name FROM orders o LEFT JOIN customers c ON c.id = o.customer_id
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| HASH JOIN
: string
| OUTER
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 4 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT o.id, c.name FROM orders o RIGHT OUTER JOIN customers c ON c.id = o.customer_id
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "NAME" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT o.id, c.name FROM orders o RIGHT OUTER JOIN customers c ON c.id = o.customer_id
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D27
< row
: string
| SHOP
: string
| OMNI_SPAN_D11
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT o.id, c.name FROM orders o RIGHT OUTER JOIN customers c ON c.id = o.customer_id
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| HASH JOIN
: string
| OUTER
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 4 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT o.id, c.id AS cid FROM orders o FULL JOIN customers c ON c.id = o.customer_id
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "CID" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT o.id, c.id AS cid FROM orders o FULL JOIN customers c ON c.id = o.customer_id
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D27
< row
: string
| SHOP
: string
| OMNI_SPAN_D10
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT o.id, c.id AS cid FROM orders o FULL JOIN customers c ON c.id = o.customer_id
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| VIEW
: null
: string
| SYS
: string
| VW_FOJ_0
: string
| VIEW
: null
: null
: null
: string
| 0
< row
: string
| HASH JOIN
: string
| FULL OUTER
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 5 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT o.id, c.name FROM orders o, customers c WHERE c.id = o.customer_id(+)
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "NAME" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT o.id, c.name FROM orders o, customers c WHERE c.id = o.customer_id(+)
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D27
< row
: string
| SHOP
: string
| OMNI_SPAN_D11
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT o.id, c.name FROM orders o, customers c WHERE c.id = o.customer_id(+)
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| HASH JOIN
: string
| OUTER
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 4 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT id, name, amount FROM customers JOIN orders USING (id)
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "NAME" "NCHAR" "string"
< column "AMOUNT" "NUMBER" "float64" precision 10 scale 2
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2, C3) AS SELECT id, name, amount FROM customers JOIN orders USING (id)
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D27
< row
: string
| SHOP
: string
| OMNI_SPAN_D11
< row
: string
| SHOP
: string
| OMNI_SPAN_D24
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT id, name, amount FROM customers JOIN orders USING (id)
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| HASH JOIN
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 4 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT * FROM customers JOIN orders USING (id)
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "NAME" "NCHAR" "string"
< column "EMAIL" "NCHAR" "string"
< column "REGION" "NCHAR" "string"
< column "CUSTOMER_ID" "NUMBER" "float64" precision 38 scale 255
< column "AMOUNT" "NUMBER" "float64" precision 10 scale 2
< column "STATUS" "NCHAR" "string"
< column "CREATED" "DATE" "time.Time"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2, C3, C4, C5, C6, C7, C8) AS SELECT * FROM customers JOIN orders USING (id)
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D27
< row
: string
| SHOP
: string
| OMNI_SPAN_D11
< row
: string
| SHOP
: string
| OMNI_SPAN_D9
< row
: string
| SHOP
: string
| OMNI_SPAN_D12
< row
: string
| SHOP
: string
| OMNI_SPAN_D26
< row
: string
| SHOP
: string
| OMNI_SPAN_D24
< row
: string
| SHOP
: string
| OMNI_SPAN_D28
< row
: string
| SHOP
: string
| OMNI_SPAN_D25
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT * FROM customers JOIN orders USING (id)
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| HASH JOIN
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 4 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT * FROM customers NATURAL JOIN people.customers
< error
| 918 ORA-00918: : column ambiguously specified - appears in CUSTOMERS and CUSTOMERS
|
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT id FROM customers NATURAL JOIN people.customers
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1) AS SELECT id FROM customers NATURAL JOIN people.customers
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| PEOPLE
: string
| OMNI_SPAN_D1
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| PEOPLE
: string
| CUSTOMERS
: string
| TABLE
: null
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT id FROM customers NATURAL JOIN people.customers
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| MERGE JOIN
: string
| CARTESIAN
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| BUFFER
: string
| SORT
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| PEOPLE
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 5 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT o.id, i.sku FROM orders o CROSS JOIN items i
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "SKU" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT o.id, i.sku FROM orders o CROSS JOIN items i
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D27
< row
: string
| SHOP
: string
| OMNI_SPAN_D20
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| ITEMS
: string
| TABLE
: null
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT o.id, i.sku FROM orders o CROSS JOIN items i
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| MERGE JOIN
: string
| CARTESIAN
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| BUFFER
: string
| SORT
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ITEMS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 5 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT o.id, x.total FROM orders o CROSS APPLY (SELECT SUM(qty) AS total FROM items i WHERE i.order_id = o.id) x
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "TOTAL" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT o.id, x.total FROM orders o CROSS APPLY (SELECT SUM(qty) AS total FROM items i WHERE i.order_id = o.id) x
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D27
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| ITEMS
: string
| TABLE
: null
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT o.id, x.total FROM orders o CROSS APPLY (SELECT SUM(qty) AS total FROM items i WHERE i.order_id = o.id) x
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| HASH JOIN
: string
| OUTER
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| VIEW
: null
: string
| SHOP
: string
| VW_DCL_535DE542
: string
| VIEW
: null
: null
: null
: string
| 0
< row
: string
| HASH
: string
| GROUP BY
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ITEMS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 6 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT o.id, x.sku FROM orders o OUTER APPLY (SELECT sku FROM items i WHERE i.order_id = o.id) x
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "SKU" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT o.id, x.sku FROM orders o OUTER APPLY (SELECT sku FROM items i WHERE i.order_id = o.id) x
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D27
< row
: string
| SHOP
: string
| OMNI_SPAN_D20
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| ITEMS
: string
| TABLE
: null
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT o.id, x.sku FROM orders o OUTER APPLY (SELECT sku FROM items i WHERE i.order_id = o.id) x
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| HASH JOIN
: string
| OUTER
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ITEMS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 4 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT o.id, l.qty FROM orders o, LATERAL (SELECT qty FROM items i WHERE i.order_id = o.id) l
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "QTY" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT o.id, l.qty FROM orders o, LATERAL (SELECT qty FROM items i WHERE i.order_id = o.id) l
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D27
< row
: string
| SHOP
: string
| OMNI_SPAN_D19
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| ITEMS
: string
| TABLE
: null
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT o.id, l.qty FROM orders o, LATERAL (SELECT qty FROM items i WHERE i.order_id = o.id) l
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| HASH JOIN
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ITEMS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 4 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT e.ename, m.ename AS boss FROM emp e LEFT JOIN emp m ON m.empno = e.mgr
< result
< column "ENAME" "NCHAR" "string"
< column "BOSS" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT e.ename, m.ename AS boss FROM emp e LEFT JOIN emp m ON m.empno = e.mgr
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D15
< row
: string
| SHOP
: string
| OMNI_SPAN_D15
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| EMP
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT e.ename, m.ename AS boss FROM emp e LEFT JOIN emp m ON m.empno = e.mgr
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| HASH JOIN
: string
| OUTER
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| EMP
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| EMP
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 4 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT x.id, x.total FROM (SELECT id, amount AS total FROM orders) x
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "TOTAL" "NUMBER" "float64" precision 10 scale 2
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT x.id, x.total FROM (SELECT id, amount AS total FROM orders) x
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D27
< row
: string
| SHOP
: string
| OMNI_SPAN_D24
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT x.id, x.total FROM (SELECT id, amount AS total FROM orders) x
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT total FROM (SELECT amount + 1 AS total FROM orders)
< result
< column "TOTAL" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1) AS SELECT total FROM (SELECT amount + 1 AS total FROM orders)
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT total FROM (SELECT amount + 1 AS total FROM orders)
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT * FROM (SELECT id, name FROM customers) WHERE id > 1
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "NAME" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT * FROM (SELECT id, name FROM customers) WHERE id > 1
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D10
< row
: string
| SHOP
: string
| OMNI_SPAN_D11
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT * FROM (SELECT id, name FROM customers) WHERE id > 1
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT id FROM orders WHERE customer_id IN (SELECT id FROM customers WHERE region = 'eu')
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1) AS SELECT id FROM orders WHERE customer_id IN (SELECT id FROM customers WHERE region = 'eu')
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D27
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT id FROM orders WHERE customer_id IN (SELECT id FROM customers WHERE region = 'eu')
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| HASH JOIN
: string
| SEMI
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 4 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT id FROM customers c WHERE EXISTS (SELECT 1 FROM orders o WHERE o.customer_id = c.id)
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1) AS SELECT id FROM customers c WHERE EXISTS (SELECT 1 FROM orders o WHERE o.customer_id = c.id)
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D10
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT id FROM customers c WHERE EXISTS (SELECT 1 FROM orders o WHERE o.customer_id = c.id)
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| HASH JOIN
: string
| SEMI
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 4 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT name FROM customers WHERE id = (SELECT MAX(customer_id) FROM orders)
< result
< column "NAME" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1) AS SELECT name FROM customers WHERE id = (SELECT MAX(customer_id) FROM orders)
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D11
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT name FROM customers WHERE id = (SELECT MAX(customer_id) FROM orders)
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| SORT
: string
| AGGREGATE
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 4 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| WITH t AS (SELECT id, amount FROM orders) SELECT id, amount FROM t
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "AMOUNT" "NUMBER" "float64" precision 10 scale 2
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS WITH t AS (SELECT id, amount FROM orders) SELECT id, amount FROM t
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D27
< row
: string
| SHOP
: string
| OMNI_SPAN_D24
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR WITH t AS (SELECT id, amount FROM orders) SELECT id, amount FROM t
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| WITH t (k, v) AS (SELECT id, amount FROM orders) SELECT k, v FROM t
< result
< column "K" "NUMBER" "float64" precision 38 scale 255
< column "V" "NUMBER" "float64" precision 10 scale 2
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS WITH t (k, v) AS (SELECT id, amount FROM orders) SELECT k, v FROM t
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D27
< row
: string
| SHOP
: string
| OMNI_SPAN_D24
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR WITH t (k, v) AS (SELECT id, amount FROM orders) SELECT k, v FROM t
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| WITH a AS (SELECT id, name FROM customers), b AS (SELECT id FROM a) SELECT * FROM b
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1) AS WITH a AS (SELECT id, name FROM customers), b AS (SELECT id FROM a) SELECT * FROM b
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D10
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR WITH a AS (SELECT id, name FROM customers), b AS (SELECT id FROM a) SELECT * FROM b
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| WITH orders AS (SELECT id, name FROM customers) SELECT * FROM orders
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "NAME" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS WITH orders AS (SELECT id, name FROM customers) SELECT * FROM orders
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D10
< row
: string
| SHOP
: string
| OMNI_SPAN_D11
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR WITH orders AS (SELECT id, name FROM customers) SELECT * FROM orders
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| WITH orders AS (SELECT id FROM customers) SELECT o.amount FROM shop.orders o JOIN orders c ON c.id = o.customer_id
< result
< column "AMOUNT" "NUMBER" "float64" precision 10 scale 2
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1) AS WITH orders AS (SELECT id FROM customers) SELECT o.amount FROM shop.orders o JOIN orders c ON c.id = o.customer_id
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D24
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR WITH orders AS (SELECT id FROM customers) SELECT o.amount FROM shop.orders o JOIN orders c ON c.id = o.customer_id
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| HASH JOIN
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 4 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| WITH r (n, nm) AS (SELECT empno, ename FROM emp WHERE mgr IS NULL UNION ALL SELECT e.empno, e.ename FROM emp e JOIN r ON e.mgr = r.n) SELECT n, nm FROM r
< result
< column "N" "NUMBER" "float64" precision 38 scale 255
< column "NM" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS WITH r (n, nm) AS (SELECT empno, ename FROM emp WHERE mgr IS NULL UNION ALL SELECT e.empno, e.ename FROM emp e JOIN r ON e.mgr = r.n) SELECT n, nm FROM r
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D14
< row
: string
| SHOP
: string
| OMNI_SPAN_D15
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| EMP
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR WITH r (n, nm) AS (SELECT empno, ename FROM emp WHERE mgr IS NULL UNION ALL SELECT e.empno, e.ename FROM emp e JOIN r ON e.mgr = r.n) SELECT n, nm FROM r
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| VIEW
: null
: string
| SHOP
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| UNION ALL (RECURSIVE WITH)
: string
| BREADTH FIRST
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| EMP
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| HASH JOIN
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| BUFFER
: string
| SORT (REUSE)
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| EMP
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| RECURSIVE WITH PUMP
: null
: null
: null
: null
: null
: null
: null
: string
| 0
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 8 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT empno, ename, LEVEL FROM emp START WITH mgr IS NULL CONNECT BY PRIOR empno = mgr
< result
< column "EMPNO" "NUMBER" "float64" precision 38 scale 255
< column "ENAME" "NCHAR" "string"
< column "LEVEL" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2, C3) AS SELECT empno, ename, LEVEL FROM emp START WITH mgr IS NULL CONNECT BY PRIOR empno = mgr
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D14
< row
: string
| SHOP
: string
| OMNI_SPAN_D15
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| EMP
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT empno, ename, LEVEL FROM emp START WITH mgr IS NULL CONNECT BY PRIOR empno = mgr
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| CONNECT BY
: string
| NO FILTERING WITH START-WITH
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| EMP
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 3 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT ename, PRIOR ename AS boss FROM emp CONNECT BY PRIOR empno = mgr
< result
< column "ENAME" "NCHAR" "string"
< column "BOSS" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT ename, PRIOR ename AS boss FROM emp CONNECT BY PRIOR empno = mgr
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D15
< row
: string
| SHOP
: string
| OMNI_SPAN_D15
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| EMP
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT ename, PRIOR ename AS boss FROM emp CONNECT BY PRIOR empno = mgr
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| CONNECT BY
: string
| WITHOUT FILTERING
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| EMP
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 3 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT ename, SYS_CONNECT_BY_PATH(ename, '/') AS path, CONNECT_BY_ROOT ename AS root FROM emp START WITH mgr IS NULL CONNECT BY PRIOR empno = mgr
< result
< column "ENAME" "NCHAR" "string"
< column "PATH" "NCHAR" "string"
< column "ROOT" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2, C3) AS SELECT ename, SYS_CONNECT_BY_PATH(ename, '/') AS path, CONNECT_BY_ROOT ename AS root FROM emp START WITH mgr IS NULL CONNECT BY PRIOR empno = mgr
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D15
< row
: null
: null
< row
: string
| SHOP
: string
| OMNI_SPAN_D15
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| EMP
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT ename, SYS_CONNECT_BY_PATH(ename, '/') AS path, CONNECT_BY_ROOT ename AS root FROM emp START WITH mgr IS NULL CONNECT BY PRIOR empno = mgr
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| CONNECT BY
: string
| NO FILTERING WITH START-WITH
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| EMP
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 3 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT LEVEL AS n FROM DUAL CONNECT BY LEVEL <= 5
< result
< column "N" "NUMBER" "float64" precision 38 scale 255
< row
: string
| 1
< row
: string
| 2
< row
: string
| 3
< row
: string
| 4
< row
: string
| 5
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1) AS SELECT LEVEL AS n FROM DUAL CONNECT BY LEVEL <= 5
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| PUBLIC
: string
| DUAL
: string
| SYNONYM
: null
> query
| SELECT table_owner, table_name, db_link FROM dba_synonyms WHERE owner = 'PUBLIC' AND synonym_name = 'DUAL'
< result
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "DB_LINK" "NCHAR" "string"
< row
: string
| SYS
: string
| DUAL
: null
> query
| SELECT object_type FROM all_objects WHERE owner = 'SYS' AND object_name = 'DUAL' AND object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW', 'SYNONYM') ORDER BY object_type
< result
< column "OBJECT_TYPE" "NCHAR" "string"
< row
: string
| TABLE
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT LEVEL AS n FROM DUAL CONNECT BY LEVEL <= 5
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| CONNECT BY
: string
| WITHOUT FILTERING
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| FAST DUAL
: null
: null
: null
: null
: null
: null
: null
: string
| 0
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 3 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT id, name FROM customers UNION SELECT id, name FROM people.employees
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "NAME" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT id, name FROM customers UNION SELECT id, name FROM people.employees
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: null
: null
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| PEOPLE
: string
| EMPLOYEES
: string
| TABLE
: null
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT id, name FROM customers UNION SELECT id, name FROM people.employees
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| HASH
: string
| UNIQUE
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| UNION-ALL
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| PEOPLE
: string
| EMPLOYEES
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 5 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT id FROM customers UNION ALL SELECT id FROM orders
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1) AS SELECT id FROM customers UNION ALL SELECT id FROM orders
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT id FROM customers UNION ALL SELECT id FROM orders
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| UNION-ALL
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 4 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT id FROM customers INTERSECT SELECT customer_id FROM orders
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1) AS SELECT id FROM customers INTERSECT SELECT customer_id FROM orders
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT id FROM customers INTERSECT SELECT customer_id FROM orders
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| INTERSECTION
: string
| HASH
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 4 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT id FROM customers MINUS SELECT customer_id FROM orders
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1) AS SELECT id FROM customers MINUS SELECT customer_id FROM orders
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D10
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT id FROM customers MINUS SELECT customer_id FROM orders
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| MINUS
: string
| HASH
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 4 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT name FROM customers UNION SELECT name FROM customers
< result
< column "NAME" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1) AS SELECT name FROM customers UNION SELECT name FROM customers
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D11
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT name FROM customers UNION SELECT name FROM customers
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| HASH
: string
| UNIQUE
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| UNION-ALL
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 5 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT id, name FROM customers UNION ALL SELECT id, 'x' FROM orders
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "NAME" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT id, name FROM customers UNION ALL SELECT id, 'x' FROM orders
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: null
: null
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT id, name FROM customers UNION ALL SELECT id, 'x' FROM orders
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| UNION-ALL
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 4 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT COUNT(*) FROM orders
< result
< column "COUNT(*)" "NUMBER" "float64" precision 38 scale 255
< row
: string
| 0
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1) AS SELECT COUNT(*) FROM orders
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT COUNT(*) FROM orders
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| SORT
: string
| AGGREGATE
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 3 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT customer_id, SUM(amount) AS total, COUNT(*) AS n FROM orders GROUP BY customer_id
< result
< column "CUSTOMER_ID" "NUMBER" "float64" precision 38 scale 255
< column "TOTAL" "NUMBER" "float64" precision 38 scale 255
< column "N" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2, C3) AS SELECT customer_id, SUM(amount) AS total, COUNT(*) AS n FROM orders GROUP BY customer_id
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D26
< row
: null
: null
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT customer_id, SUM(amount) AS total, COUNT(*) AS n FROM orders GROUP BY customer_id
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| HASH
: string
| GROUP BY
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 3 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT status, MAX(amount) FROM orders GROUP BY status HAVING COUNT(*) > 1
< result
< column "STATUS" "NCHAR" "string"
< column "MAX(AMOUNT)" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT status, MAX(amount) FROM orders GROUP BY status HAVING COUNT(*) > 1
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D28
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT status, MAX(amount) FROM orders GROUP BY status HAVING COUNT(*) > 1
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| FILTER
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| HASH
: string
| GROUP BY
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 4 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT MIN(name) AS first FROM customers
< result
< column "FIRST" "NCHAR" "string"
< row
: null
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1) AS SELECT MIN(name) AS first FROM customers
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT MIN(name) AS first FROM customers
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| SORT
: string
| AGGREGATE
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 3 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT region, LISTAGG(name, ',') WITHIN GROUP (ORDER BY name) AS names FROM customers GROUP BY region
< result
< column "REGION" "NCHAR" "string"
< column "NAMES" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT region, LISTAGG(name, ',') WITHIN GROUP (ORDER BY name) AS names FROM customers GROUP BY region
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D12
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT region, LISTAGG(name, ',') WITHIN GROUP (ORDER BY name) AS names FROM customers GROUP BY region
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| SORT
: string
| GROUP BY
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 3 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT id, amount, SUM(amount) OVER (PARTITION BY customer_id ORDER BY id) AS running FROM orders
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "AMOUNT" "NUMBER" "float64" precision 10 scale 2
< column "RUNNING" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2, C3) AS SELECT id, amount, SUM(amount) OVER (PARTITION BY customer_id ORDER BY id) AS running FROM orders
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D27
< row
: string
| SHOP
: string
| OMNI_SPAN_D24
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT id, amount, SUM(amount) OVER (PARTITION BY customer_id ORDER BY id) AS running FROM orders
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| WINDOW
: string
| SORT
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 3 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT id, ROW_NUMBER() OVER (ORDER BY amount) AS rn FROM orders
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "RN" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT id, ROW_NUMBER() OVER (ORDER BY amount) AS rn FROM orders
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D27
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT id, ROW_NUMBER() OVER (ORDER BY amount) AS rn FROM orders
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| WINDOW
: string
| SORT
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 3 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT id, LAG(amount) OVER (ORDER BY id) AS prev FROM orders
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "PREV" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT id, LAG(amount) OVER (ORDER BY id) AS prev FROM orders
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D27
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT id, LAG(amount) OVER (ORDER BY id) AS prev FROM orders
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| WINDOW
: string
| SORT
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 3 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT DISTINCT region FROM customers
< result
< column "REGION" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1) AS SELECT DISTINCT region FROM customers
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D12
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT DISTINCT region FROM customers
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| HASH
: string
| UNIQUE
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 3 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT region, COUNT(DISTINCT id) FROM customers GROUP BY ROLLUP (region)
< result
< column "REGION" "NCHAR" "string"
< column "COUNT(DISTINCTID)" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT region, COUNT(DISTINCT id) FROM customers GROUP BY ROLLUP (region)
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D12
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT region, COUNT(DISTINCT id) FROM customers GROUP BY ROLLUP (region)
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| SORT
: string
| GROUP BY ROLLUP
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 3 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT * FROM sales PIVOT (SUM(amount) FOR quarter IN ('Q1' AS q1, 'Q2' AS q2))
< result
< column "REGION" "NCHAR" "string"
< column "Q1" "NUMBER" "float64" precision 38 scale 255
< column "Q2" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2, C3) AS SELECT * FROM sales PIVOT (SUM(amount) FOR quarter IN ('Q1' AS q1, 'Q2' AS q2))
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D34
< row
: null
: null
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| SALES
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT * FROM sales PIVOT (SUM(amount) FOR quarter IN ('Q1' AS q1, 'Q2' AS q2))
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| HASH
: string
| GROUP BY PIVOT
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| SALES
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 3 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT region, q1 FROM sales PIVOT (SUM(amount) FOR quarter IN ('Q1' AS q1, 'Q2' AS q2))
< result
< column "REGION" "NCHAR" "string"
< column "Q1" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT region, q1 FROM sales PIVOT (SUM(amount) FOR quarter IN ('Q1' AS q1, 'Q2' AS q2))
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D34
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| SALES
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT region, q1 FROM sales PIVOT (SUM(amount) FOR quarter IN ('Q1' AS q1, 'Q2' AS q2))
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| HASH
: string
| GROUP BY PIVOT
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| SALES
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 3 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT * FROM sales PIVOT (SUM(amount) AS total FOR quarter IN ('q1', 'Q2' AS b, 3))
< result
< column "REGION" "NCHAR" "string"
< column "'q1'_TOTAL" "NUMBER" "float64" precision 38 scale 255
< column "B_TOTAL" "NUMBER" "float64" precision 38 scale 255
< column "3_TOTAL" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2, C3, C4) AS SELECT * FROM sales PIVOT (SUM(amount) AS total FOR quarter IN ('q1', 'Q2' AS b, 3))
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D34
< row
: null
: null
< row
: null
: null
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| SALES
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT * FROM sales PIVOT (SUM(amount) AS total FOR quarter IN ('q1', 'Q2' AS b, 3))
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| HASH
: string
| GROUP BY PIVOT
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| SALES
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 3 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT * FROM sales PIVOT (SUM(amount) FOR quarter IN ('q1', 'a b'))
< result
< column "REGION" "NCHAR" "string"
< column "'q1'" "NUMBER" "float64" precision 38 scale 255
< column "'a b'" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2, C3) AS SELECT * FROM sales PIVOT (SUM(amount) FOR quarter IN ('q1', 'a b'))
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D34
< row
: null
: null
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| SALES
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT * FROM sales PIVOT (SUM(amount) FOR quarter IN ('q1', 'a b'))
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| HASH
: string
| GROUP BY PIVOT
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| SALES
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 3 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT * FROM quarterly UNPIVOT (amount FOR quarter IN (q1 AS 'Q1', q2 AS 'Q2'))
< result
< column "REGION" "NCHAR" "string"
< column "QUARTER" "CHAR" "string"
< column "AMOUNT" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2, C3) AS SELECT * FROM quarterly UNPIVOT (amount FOR quarter IN (q1 AS 'Q1', q2 AS 'Q2'))
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: null
: null
< row
: null
: null
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| PUBLIC
: string
| DUAL
: string
| SYNONYM
: null
< row
: string
| SHOP
: string
| QUARTERLY
: string
| TABLE
: null
> query
| SELECT table_owner, table_name, db_link FROM dba_synonyms WHERE owner = 'PUBLIC' AND synonym_name = 'DUAL'
< result
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "DB_LINK" "NCHAR" "string"
< row
: string
| SYS
: string
| DUAL
: null
> query
| SELECT object_type FROM all_objects WHERE owner = 'SYS' AND object_name = 'DUAL' AND object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW', 'SYNONYM') ORDER BY object_type
< result
< column "OBJECT_TYPE" "NCHAR" "string"
< row
: string
| TABLE
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT * FROM quarterly UNPIVOT (amount FOR quarter IN (q1 AS 'Q1', q2 AS 'Q2'))
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| VIEW
: null
: string
| SYS
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| UNPIVOT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| QUARTERLY
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 4 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT region, amount FROM quarterly UNPIVOT (amount FOR quarter IN (q1, q2))
< result
< column "REGION" "NCHAR" "string"
< column "AMOUNT" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT region, amount FROM quarterly UNPIVOT (amount FOR quarter IN (q1, q2))
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: null
: null
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| PUBLIC
: string
| DUAL
: string
| SYNONYM
: null
< row
: string
| SHOP
: string
| QUARTERLY
: string
| TABLE
: null
> query
| SELECT table_owner, table_name, db_link FROM dba_synonyms WHERE owner = 'PUBLIC' AND synonym_name = 'DUAL'
< result
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "DB_LINK" "NCHAR" "string"
< row
: string
| SYS
: string
| DUAL
: null
> query
| SELECT object_type FROM all_objects WHERE owner = 'SYS' AND object_name = 'DUAL' AND object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW', 'SYNONYM') ORDER BY object_type
< result
< column "OBJECT_TYPE" "NCHAR" "string"
< row
: string
| TABLE
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT region, amount FROM quarterly UNPIVOT (amount FOR quarter IN (q1, q2))
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| VIEW
: null
: string
| SYS
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| UNPIVOT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| QUARTERLY
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 4 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT region, quarter, amount FROM sales MODEL PARTITION BY (region) DIMENSION BY (quarter) MEASURES (amount) RULES (amount['Q9'] = amount['Q1'])
< result
< column "REGION" "NCHAR" "string"
< column "QUARTER" "NCHAR" "string"
< column "AMOUNT" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2, C3) AS SELECT region, quarter, amount FROM sales MODEL PARTITION BY (region) DIMENSION BY (quarter) MEASURES (amount) RULES (amount['Q9'] = amount['Q1'])
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D34
< row
: string
| SHOP
: string
| OMNI_SPAN_D33
< row
: string
| SHOP
: string
| OMNI_SPAN_D32
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| SALES
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT region, quarter, amount FROM sales MODEL PARTITION BY (region) DIMENSION BY (quarter) MEASURES (amount) RULES (amount['Q9'] = amount['Q1'])
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| SQL MODEL
: string
| ORDERED FAST
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| SALES
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 3 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT * FROM order_totals
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "CUSTOMER_ID" "NUMBER" "float64" precision 38 scale 255
< column "CUSTOMER_NAME" "NCHAR" "string"
< column "GROSS" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2, C3, C4) AS SELECT * FROM order_totals
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D27
< row
: string
| SHOP
: string
| OMNI_SPAN_D26
< row
: string
| SHOP
: string
| OMNI_SPAN_D11
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| ORDER_TOTALS
: string
| VIEW
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'ORDER_TOTALS' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT * FROM order_totals
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| HASH JOIN
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 4 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT id, gross FROM order_totals
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "GROSS" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT id, gross FROM order_totals
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D27
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| ORDER_TOTALS
: string
| VIEW
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'ORDER_TOTALS' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT id, gross FROM order_totals
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| HASH JOIN
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 4 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT customer_name FROM shop.order_totals WHERE customer_id = 1
< result
< column "CUSTOMER_NAME" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1) AS SELECT customer_name FROM shop.order_totals WHERE customer_id = 1
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D11
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| ORDER_TOTALS
: string
| VIEW
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'ORDER_TOTALS' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT customer_name FROM shop.order_totals WHERE customer_id = 1
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| HASH JOIN
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 4 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT * FROM big_orders
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "CUSTOMER_NAME" "NCHAR" "string"
< column "GROSS" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2, C3) AS SELECT * FROM big_orders
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D27
< row
: string
| SHOP
: string
| OMNI_SPAN_D11
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| BIG_ORDERS
: string
| VIEW
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'BIG_ORDERS' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| ORDER_TOTALS
: string
| VIEW
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'ORDER_TOTALS' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT * FROM big_orders
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| HASH JOIN
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 4 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT b.customer_name, o.status FROM big_orders b JOIN orders o ON o.id = b.id
< result
< column "CUSTOMER_NAME" "NCHAR" "string"
< column "STATUS" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT b.customer_name, o.status FROM big_orders b JOIN orders o ON o.id = b.id
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D11
< row
: string
| SHOP
: string
| OMNI_SPAN_D28
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| BIG_ORDERS
: string
| VIEW
: null
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'BIG_ORDERS' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| ORDER_TOTALS
: string
| VIEW
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'ORDER_TOTALS' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT b.customer_name, o.status FROM big_orders b JOIN orders o ON o.id = b.id
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| HASH JOIN
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| MERGE JOIN
: string
| CARTESIAN
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| BUFFER
: string
| SORT
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 7 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT * FROM all_people
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "NAME" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT * FROM all_people
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: null
: null
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| ALL_PEOPLE
: string
| VIEW
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'ALL_PEOPLE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| PEOPLE
: string
| EMPLOYEES
: string
| TABLE
: null
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT * FROM all_people
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| VIEW
: null
: string
| SHOP
: string
| ALL_PEOPLE
: string
| VIEW
: null
: null
: null
: string
| 1
< row
: string
| UNION-ALL
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| PEOPLE
: string
| EMPLOYEES
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 5 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT * FROM people.staff
< result
< column "STAFF_ID" "NUMBER" "float64" precision 38 scale 255
< column "STAFF_NAME" "NCHAR" "string"
< column "DEPT" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2, C3) AS SELECT * FROM people.staff
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| PEOPLE
: string
| OMNI_SPAN_D6
< row
: string
| PEOPLE
: string
| OMNI_SPAN_D7
< row
: string
| PEOPLE
: string
| OMNI_SPAN_D4
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| PEOPLE
: string
| STAFF
: string
| VIEW
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'PEOPLE' AND name = 'STAFF' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| PEOPLE
: string
| DEPARTMENTS
: string
| TABLE
: null
< row
: string
| PEOPLE
: string
| EMPLOYEES
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT * FROM people.staff
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| HASH JOIN
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| PEOPLE
: string
| EMPLOYEES
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| PEOPLE
: string
| DEPARTMENTS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 4 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT staff_name, dept FROM people.staff WHERE staff_id = 1
< result
< column "STAFF_NAME" "NCHAR" "string"
< column "DEPT" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT staff_name, dept FROM people.staff WHERE staff_id = 1
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| PEOPLE
: string
| OMNI_SPAN_D7
< row
: string
| PEOPLE
: string
| OMNI_SPAN_D4
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| PEOPLE
: string
| STAFF
: string
| VIEW
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'PEOPLE' AND name = 'STAFF' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| PEOPLE
: string
| DEPARTMENTS
: string
| TABLE
: null
< row
: string
| PEOPLE
: string
| EMPLOYEES
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT staff_name, dept FROM people.staff WHERE staff_id = 1
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| HASH JOIN
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| PEOPLE
: string
| EMPLOYEES
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| PEOPLE
: string
| DEPARTMENTS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 4 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT * FROM "Vip Names"
< result
< column "VipName" "NCHAR" "string"
< column "REGION" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT * FROM "Vip Names"
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D11
< row
: string
| SHOP
: string
| OMNI_SPAN_D12
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| Vip Names
: string
| VIEW
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'Vip Names' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT * FROM "Vip Names"
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT "VipName" FROM shop."Vip Names"
< result
< column "VipName" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1) AS SELECT "VipName" FROM shop."Vip Names"
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D11
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| Vip Names
: string
| VIEW
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'Vip Names' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT "VipName" FROM shop."Vip Names"
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT * FROM region_counts
< result
< column "REGION" "NCHAR" "string"
< column "N" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT * FROM region_counts
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: null
: null
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| REGION_COUNTS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT * FROM region_counts
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| MAT_VIEW ACCESS
: string
| FULL
: string
| SHOP
: string
| REGION_COUNTS
: string
| MAT_VIEW
: null
: null
: null
: string
| 2
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT region, n FROM region_counts
< result
< column "REGION" "NCHAR" "string"
< column "N" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT region, n FROM region_counts
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: null
: null
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| REGION_COUNTS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT region, n FROM region_counts
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| MAT_VIEW ACCESS
: string
| FULL
: string
| SHOP
: string
| REGION_COUNTS
: string
| MAT_VIEW
: null
: null
: null
: string
| 2
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT id, amount FROM syn_orders
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "AMOUNT" "NUMBER" "float64" precision 10 scale 2
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT id, amount FROM syn_orders
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D27
< row
: string
| SHOP
: string
| OMNI_SPAN_D24
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| SYN_ORDERS
: string
| SYNONYM
: null
> query
| SELECT table_owner, table_name, db_link FROM dba_synonyms WHERE owner = 'SHOP' AND synonym_name = 'SYN_ORDERS'
< result
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "DB_LINK" "NCHAR" "string"
< row
: string
| SHOP
: string
| ORDERS
: null
> query
| SELECT object_type FROM all_objects WHERE owner = 'SHOP' AND object_name = 'ORDERS' AND object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW', 'SYNONYM') ORDER BY object_type
< result
< column "OBJECT_TYPE" "NCHAR" "string"
< row
: string
| TABLE
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT id, amount FROM syn_orders
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT syn_orders.id, s2.status FROM syn_orders, shop.syn_orders s2 WHERE s2.id = syn_orders.id
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "STATUS" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT syn_orders.id, s2.status FROM syn_orders, shop.syn_orders s2 WHERE s2.id = syn_orders.id
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D27
< row
: string
| SHOP
: string
| OMNI_SPAN_D28
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| SYN_ORDERS
: string
| SYNONYM
: null
> query
| SELECT table_owner, table_name, db_link FROM dba_synonyms WHERE owner = 'SHOP' AND synonym_name = 'SYN_ORDERS'
< result
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "DB_LINK" "NCHAR" "string"
< row
: string
| SHOP
: string
| ORDERS
: null
> query
| SELECT object_type FROM all_objects WHERE owner = 'SHOP' AND object_name = 'ORDERS' AND object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW', 'SYNONYM') ORDER BY object_type
< result
< column "OBJECT_TYPE" "NCHAR" "string"
< row
: string
| TABLE
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT syn_orders.id, s2.status FROM syn_orders, shop.syn_orders s2 WHERE s2.id = syn_orders.id
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| HASH JOIN
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 4 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT name FROM remote_emp
< result
< column "NAME" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1) AS SELECT name FROM remote_emp
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| REMOTE_EMP
: string
| SYNONYM
: null
> query
| SELECT table_owner, table_name, db_link FROM dba_synonyms WHERE owner = 'SHOP' AND synonym_name = 'REMOTE_EMP'
< result
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "DB_LINK" "NCHAR" "string"
< row
: string
| PEOPLE
: string
| EMPLOYEES
: string
| OMNI_LOOP
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT name FROM remote_emp
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: string
| REMOTE
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| PEOPLE
: string
| EMPLOYEES
: string
| TABLE
: string
| FREEPDB1
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT * FROM staff_syn
< result
< column "STAFF_ID" "NUMBER" "float64" precision 38 scale 255
< column "STAFF_NAME" "NCHAR" "string"
< column "DEPT" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2, C3) AS SELECT * FROM staff_syn
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| PEOPLE
: string
| OMNI_SPAN_D6
< row
: string
| PEOPLE
: string
| OMNI_SPAN_D7
< row
: string
| PEOPLE
: string
| OMNI_SPAN_D4
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| STAFF_SYN
: string
| SYNONYM
: null
> query
| SELECT table_owner, table_name, db_link FROM dba_synonyms WHERE owner = 'SHOP' AND synonym_name = 'STAFF_SYN'
< result
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "DB_LINK" "NCHAR" "string"
< row
: string
| PEOPLE
: string
| STAFF
: null
> query
| SELECT object_type FROM all_objects WHERE owner = 'PEOPLE' AND object_name = 'STAFF' AND object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW', 'SYNONYM') ORDER BY object_type
< result
< column "OBJECT_TYPE" "NCHAR" "string"
< row
: string
| VIEW
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'PEOPLE' AND name = 'STAFF' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| PEOPLE
: string
| DEPARTMENTS
: string
| TABLE
: null
< row
: string
| PEOPLE
: string
| EMPLOYEES
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT * FROM staff_syn
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| HASH JOIN
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| PEOPLE
: string
| EMPLOYEES
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| PEOPLE
: string
| DEPARTMENTS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 4 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT name FROM omni_pub_emp
< result
< column "NAME" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1) AS SELECT name FROM omni_pub_emp
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| PEOPLE
: string
| OMNI_SPAN_D7
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| PUBLIC
: string
| OMNI_PUB_EMP
: string
| SYNONYM
: null
> query
| SELECT table_owner, table_name, db_link FROM dba_synonyms WHERE owner = 'PUBLIC' AND synonym_name = 'OMNI_PUB_EMP'
< result
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "DB_LINK" "NCHAR" "string"
< row
: string
| PEOPLE
: string
| EMPLOYEES
: null
> query
| SELECT object_type FROM all_objects WHERE owner = 'PEOPLE' AND object_name = 'EMPLOYEES' AND object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW', 'SYNONYM') ORDER BY object_type
< result
< column "OBJECT_TYPE" "NCHAR" "string"
< row
: string
| TABLE
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT name FROM omni_pub_emp
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| PEOPLE
: string
| EMPLOYEES
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT id, name FROM employees@omni_loop
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "NAME" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT id, name FROM employees@omni_loop
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: null
: null
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT id, name FROM employees@omni_loop
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: string
| REMOTE
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| PEOPLE
: string
| EMPLOYEES
: string
| TABLE
: string
| FREEPDB1
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT e.name, c.name AS customer FROM people.employees@omni_loop e JOIN customers c ON c.id = e.id
< result
< column "NAME" "NCHAR" "string"
< column "CUSTOMER" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT e.name, c.name AS customer FROM people.employees@omni_loop e JOIN customers c ON c.id = e.id
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: null
: null
< row
: string
| SHOP
: string
| OMNI_SPAN_D11
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT e.name, c.name AS customer FROM people.employees@omni_loop e JOIN customers c ON c.id = e.id
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| HASH JOIN
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| REMOTE
: null
: null
: string
| EMPLOYEES
: string
| REMOTE
: string
| OMNI_LOOP
: null
: null
: string
| 0
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 4 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT * FROM "Mixed Case"
< result
< column "Id" "NUMBER" "float64" precision 38 scale 255
< column "Full Name" "NCHAR" "string"
< column "PLAIN" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2, C3) AS SELECT * FROM "Mixed Case"
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D22
< row
: string
| SHOP
: string
| OMNI_SPAN_D21
< row
: string
| SHOP
: string
| OMNI_SPAN_D23
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| Mixed Case
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT * FROM "Mixed Case"
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| Mixed Case
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT "Id", "Full Name", PLAIN FROM "Mixed Case"
< result
< column "Id" "NUMBER" "float64" precision 38 scale 255
< column "Full Name" "NCHAR" "string"
< column "PLAIN" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2, C3) AS SELECT "Id", "Full Name", PLAIN FROM "Mixed Case"
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D22
< row
: string
| SHOP
: string
| OMNI_SPAN_D21
< row
: string
| SHOP
: string
| OMNI_SPAN_D23
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| Mixed Case
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT "Id", "Full Name", PLAIN FROM "Mixed Case"
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| Mixed Case
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT m."Full Name" AS fn FROM shop."Mixed Case" m
< result
< column "FN" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1) AS SELECT m."Full Name" AS fn FROM shop."Mixed Case" m
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D21
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| Mixed Case
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT m."Full Name" AS fn FROM shop."Mixed Case" m
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| Mixed Case
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT "id", "val" FROM "lower"
< result
< column "id" "NUMBER" "float64" precision 38 scale 255
< column "val" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT "id", "val" FROM "lower"
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D38
< row
: string
| SHOP
: string
| OMNI_SPAN_D39
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| lower
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT "id", "val" FROM "lower"
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| lower
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT * FROM "lower"
< result
< column "id" "NUMBER" "float64" precision 38 scale 255
< column "val" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT * FROM "lower"
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D38
< row
: string
| SHOP
: string
| OMNI_SPAN_D39
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| lower
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT * FROM "lower"
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| lower
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT ID, NAME FROM SHOP.CUSTOMERS
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "NAME" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT ID, NAME FROM SHOP.CUSTOMERS
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D10
< row
: string
| SHOP
: string
| OMNI_SPAN_D11
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT ID, NAME FROM SHOP.CUSTOMERS
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT "ID", "NAME" FROM "SHOP"."CUSTOMERS"
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "NAME" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT "ID", "NAME" FROM "SHOP"."CUSTOMERS"
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D10
< row
: string
| SHOP
: string
| OMNI_SPAN_D11
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT "ID", "NAME" FROM "SHOP"."CUSTOMERS"
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT * FROM secrets
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "NOTE" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT * FROM secrets
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D35
< row
: string
| SHOP
: string
| OMNI_SPAN_D36
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| SECRETS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT * FROM secrets
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| SECRETS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT id, token FROM secrets
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "TOKEN" "NCHAR" "string"
> exec
| CREATE VIEW "SHOP".OMNI_SPAN_PROBE (C1, C2) AS SELECT id, token FROM secrets
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'SHOP' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D35
< row
: string
| SHOP
: string
| OMNI_SPAN_D37
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| SECRETS
: string
| TABLE
: null
> exec
| DROP VIEW "SHOP".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT id, token FROM secrets
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| SECRETS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT nope FROM orders
< error
| 904 ORA-00904: "NOPE": invalid identifier
|
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT id FROM missing_table
< error
| 942 ORA-00942: table or view "SHOP"."MISSING_TABLE" does not exist
|
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT id FROM orders, customers
< error
| 918 ORA-00918: ID: column ambiguously specified - appears in CUSTOMERS and ORDERS
|
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "SHOP"
< ok 0 0
> query
| SELECT id FROM "orders"
< error
| 942 ORA-00942: table or view "SHOP"."orders" does not exist
|
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "PEOPLE"
< ok 0 0
> query
| SELECT * FROM employees
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "NAME" "NCHAR" "string"
< column "SALARY" "NUMBER" "float64" precision 38 scale 255
< column "DEPT_ID" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "PEOPLE".OMNI_SPAN_PROBE (C1, C2, C3, C4) AS SELECT * FROM employees
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'PEOPLE' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| PEOPLE
: string
| OMNI_SPAN_D6
< row
: string
| PEOPLE
: string
| OMNI_SPAN_D7
< row
: string
| PEOPLE
: string
| OMNI_SPAN_D8
< row
: string
| PEOPLE
: string
| OMNI_SPAN_D5
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'PEOPLE' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| PEOPLE
: string
| EMPLOYEES
: string
| TABLE
: null
> exec
| DROP VIEW "PEOPLE".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT * FROM employees
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| PEOPLE
: string
| EMPLOYEES
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "PEOPLE"
< ok 0 0
> query
| SELECT id, name FROM customers
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "NAME" "NCHAR" "string"
> exec
| CREATE VIEW "PEOPLE".OMNI_SPAN_PROBE (C1, C2) AS SELECT id, name FROM customers
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'PEOPLE' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| PEOPLE
: string
| OMNI_SPAN_D1
< row
: string
| PEOPLE
: string
| OMNI_SPAN_D2
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'PEOPLE' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| PEOPLE
: string
| CUSTOMERS
: string
| TABLE
: null
> exec
| DROP VIEW "PEOPLE".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT id, name FROM customers
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| PEOPLE
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 2 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "PEOPLE"
< ok 0 0
> query
| SELECT * FROM staff
< result
< column "STAFF_ID" "NUMBER" "float64" precision 38 scale 255
< column "STAFF_NAME" "NCHAR" "string"
< column "DEPT" "NCHAR" "string"
> exec
| CREATE VIEW "PEOPLE".OMNI_SPAN_PROBE (C1, C2, C3) AS SELECT * FROM staff
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'PEOPLE' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| PEOPLE
: string
| OMNI_SPAN_D6
< row
: string
| PEOPLE
: string
| OMNI_SPAN_D7
< row
: string
| PEOPLE
: string
| OMNI_SPAN_D4
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'PEOPLE' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| PEOPLE
: string
| STAFF
: string
| VIEW
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'PEOPLE' AND name = 'STAFF' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| PEOPLE
: string
| DEPARTMENTS
: string
| TABLE
: null
< row
: string
| PEOPLE
: string
| EMPLOYEES
: string
| TABLE
: null
> exec
| DROP VIEW "PEOPLE".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT * FROM staff
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| HASH JOIN
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| PEOPLE
: string
| EMPLOYEES
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| PEOPLE
: string
| DEPARTMENTS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 4 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "PEOPLE"
< ok 0 0
> query
| SELECT * FROM shop.order_totals
< result
< column "ID" "NUMBER" "float64" precision 38 scale 255
< column "CUSTOMER_ID" "NUMBER" "float64" precision 38 scale 255
< column "CUSTOMER_NAME" "NCHAR" "string"
< column "GROSS" "NUMBER" "float64" precision 38 scale 255
> exec
| CREATE VIEW "PEOPLE".OMNI_SPAN_PROBE (C1, C2, C3, C4) AS SELECT * FROM shop.order_totals
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'PEOPLE' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| OMNI_SPAN_D27
< row
: string
| SHOP
: string
| OMNI_SPAN_D26
< row
: string
| SHOP
: string
| OMNI_SPAN_D11
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'PEOPLE' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| ORDER_TOTALS
: string
| VIEW
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'SHOP' AND name = 'ORDER_TOTALS' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
< row
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
> exec
| DROP VIEW "PEOPLE".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT * FROM shop.order_totals
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| HASH JOIN
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| ORDERS
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 4 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "PEOPLE"
< ok 0 0
> query
| SELECT e.name, d.name AS dept FROM employees e JOIN departments d ON d.id = e.dept_id
< result
< column "NAME" "NCHAR" "string"
< column "DEPT" "NCHAR" "string"
> exec
| CREATE VIEW "PEOPLE".OMNI_SPAN_PROBE (C1, C2) AS SELECT e.name, d.name AS dept FROM employees e JOIN departments d ON d.id = e.dept_id
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'PEOPLE' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: string
| PEOPLE
: string
| OMNI_SPAN_D7
< row
: string
| PEOPLE
: string
| OMNI_SPAN_D4
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'PEOPLE' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| PEOPLE
: string
| DEPARTMENTS
: string
| TABLE
: null
< row
: string
| PEOPLE
: string
| EMPLOYEES
: string
| TABLE
: null
> exec
| DROP VIEW "PEOPLE".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT e.name, d.name AS dept FROM employees e JOIN departments d ON d.id = e.dept_id
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| HASH JOIN
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| PEOPLE
: string
| EMPLOYEES
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| PEOPLE
: string
| DEPARTMENTS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 4 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = "PEOPLE"
< ok 0 0
> query
| SELECT name FROM shop.customers UNION SELECT name FROM customers
< result
< column "NAME" "NCHAR" "string"
> exec
| CREATE VIEW "PEOPLE".OMNI_SPAN_PROBE (C1) AS SELECT name FROM shop.customers UNION SELECT name FROM customers
< ok 0 0
> query
| SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = 'PEOPLE' AND table_name = 'OMNI_SPAN_PROBE' ORDER BY column_id
< result
< column "DOMAIN_OWNER" "NCHAR" "string"
< column "DOMAIN_NAME" "NCHAR" "string"
< row
: null
: null
> query
| SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
| FROM all_dependencies WHERE owner = 'PEOPLE' AND name = 'OMNI_SPAN_PROBE' AND type = 'VIEW'
| ORDER BY 1, 2, 3
< result
< column "REFERENCED_OWNER" "NCHAR" "string"
< column "REFERENCED_NAME" "NCHAR" "string"
< column "REFERENCED_TYPE" "NCHAR" "string"
< column "REFERENCED_LINK_NAME" "NCHAR" "string"
< row
: string
| PEOPLE
: string
| CUSTOMERS
: string
| TABLE
: null
< row
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
> exec
| DROP VIEW "PEOPLE".OMNI_SPAN_PROBE
< ok 0 0
> exec
| EXPLAIN PLAN SET STATEMENT_ID = 'omni_span' FOR SELECT name FROM shop.customers UNION SELECT name FROM customers
< ok 0 0
> query
| SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
| 	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
| FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
| WHERE p.statement_id = 'omni_span' ORDER BY p.id
< result
< column "OPERATION" "NCHAR" "string"
< column "OPTIONS" "NCHAR" "string"
< column "OBJECT_OWNER" "NCHAR" "string"
< column "OBJECT_NAME" "NCHAR" "string"
< column "OBJECT_TYPE" "NCHAR" "string"
< column "OBJECT_NODE" "NCHAR" "string"
< column "TABLE_OWNER" "NCHAR" "string"
< column "TABLE_NAME" "NCHAR" "string"
< column "(SELECTCOUNT(*)FROMALL_OBJECTSOWHEREO.OWNER=P.OBJECT_OWNERANDO.OBJECT_NAME=P.OBJECT_NAMEANDO.OBJECT_TYPEIN('TABLE','VIEW','MATERIALIZEDVIEW'))" "NUMBER" "float64" precision 38 scale 255
< row
: string
| SELECT STATEMENT
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| HASH
: string
| UNIQUE
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| UNION-ALL
: null
: null
: null
: null
: null
: null
: null
: string
| 0
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| SHOP
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
< row
: string
| TABLE ACCESS
: string
| FULL
: string
| PEOPLE
: string
| CUSTOMERS
: string
| TABLE
: null
: null
: null
: string
| 1
> exec
| DELETE FROM plan_table WHERE statement_id = 'omni_span'
< ok 5 0
> exec
| ALTER SESSION SET CURRENT_SCHEMA = SYSTEM
< ok 0 0
