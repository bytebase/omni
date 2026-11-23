-- Two schemas, SHOP and PEOPLE, built by SYSTEM. Tables have no keys or
-- indexes, so the optimizer drops no table from a plan. PEOPLE is not HR:
-- Bytebase counts Oracle's sample schemas, HR among them, as system
-- schemas, and omni gives a query on one no lineage.
CREATE USER shop IDENTIFIED BY "omni_test" QUOTA UNLIMITED ON USERS;
CREATE USER people IDENTIFIED BY "omni_test" QUOTA UNLIMITED ON USERS;
GRANT CREATE SESSION, CREATE TABLE TO shop, people;
-- A view over another schema's table needs its owner to hold the grant
-- directly, and so does the view the test makes over each query.
GRANT SELECT ANY TABLE TO shop, people;
CREATE TABLE shop.customers (id NUMBER, name VARCHAR2(50), email VARCHAR2(100), region VARCHAR2(20));
CREATE TABLE shop.orders (id NUMBER, customer_id NUMBER, amount NUMBER(10,2), status VARCHAR2(10), created DATE);
CREATE TABLE shop.items (order_id NUMBER, sku VARCHAR2(20), qty NUMBER);
CREATE TABLE shop.emp (empno NUMBER, ename VARCHAR2(30), mgr NUMBER, sal NUMBER, deptno NUMBER);
CREATE TABLE shop.sales (region VARCHAR2(20), quarter VARCHAR2(2), amount NUMBER);
CREATE TABLE shop.quarterly (region VARCHAR2(20), q1 NUMBER, q2 NUMBER);
CREATE TABLE shop.secrets (id NUMBER, token VARCHAR2(40) INVISIBLE, note VARCHAR2(40));
CREATE TABLE shop."Mixed Case" ("Id" NUMBER, "Full Name" VARCHAR2(50), plain VARCHAR2(10));
CREATE TABLE shop."lower" ("id" NUMBER, "val" VARCHAR2(10));
CREATE TABLE people.employees (id NUMBER, name VARCHAR2(50), salary NUMBER, dept_id NUMBER);
CREATE TABLE people.departments (id NUMBER, name VARCHAR2(50));
CREATE TABLE people.customers (id NUMBER, name VARCHAR2(50));
CREATE VIEW shop.order_totals AS SELECT o.id, o.customer_id, c.name AS customer_name, o.amount * 2 AS gross FROM shop.orders o JOIN shop.customers c ON c.id = o.customer_id;
CREATE VIEW shop.big_orders AS SELECT id, customer_name, gross FROM shop.order_totals WHERE gross > 100;
CREATE VIEW shop.all_people AS SELECT id, name FROM shop.customers UNION ALL SELECT id, name FROM people.employees;
CREATE VIEW shop."Vip Names" AS SELECT name AS "VipName", region FROM shop.customers WHERE region = 'vip';
CREATE VIEW people.staff AS SELECT e.id AS staff_id, e.name AS staff_name, d.name AS dept FROM people.employees e JOIN people.departments d ON d.id = e.dept_id;
CREATE MATERIALIZED VIEW shop.region_counts AS SELECT region, COUNT(*) AS n FROM shop.customers GROUP BY region;
CREATE SYNONYM shop.syn_orders FOR shop.orders;
CREATE SYNONYM shop.staff_syn FOR people.staff;
CREATE PUBLIC SYNONYM omni_pub_emp FOR people.employees;
CREATE PUBLIC DATABASE LINK omni_loop CONNECT TO people IDENTIFIED BY "omni_test" USING 'localhost:1521/FREEPDB1';
CREATE SYNONYM shop.remote_emp FOR people.employees@omni_loop;
