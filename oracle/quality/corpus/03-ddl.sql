-- @name: create table simple
-- @valid: true
-- @source: Oracle SQL Reference
CREATE TABLE test_table (
  id NUMBER PRIMARY KEY,
  name VARCHAR2(100) NOT NULL,
  created_date DATE DEFAULT SYSDATE
)

-- @name: create table with constraints
-- @valid: true
-- @source: Oracle SQL Reference
CREATE TABLE orders (
  order_id NUMBER,
  customer_id NUMBER NOT NULL,
  order_date DATE,
  amount NUMBER(10,2),
  CONSTRAINT pk_orders PRIMARY KEY (order_id),
  CONSTRAINT fk_customer FOREIGN KEY (customer_id) REFERENCES customers(customer_id),
  CONSTRAINT chk_amount CHECK (amount > 0)
)

-- @name: create table as select
-- @valid: true
-- @source: Oracle SQL Reference
CREATE TABLE emp_backup AS SELECT * FROM employees WHERE department_id = 10

-- @name: create global temporary table
-- @valid: true
-- @source: Oracle SQL Reference
CREATE GLOBAL TEMPORARY TABLE temp_results (
  id NUMBER,
  result VARCHAR2(200)
) ON COMMIT DELETE ROWS

-- @name: create index
-- @valid: true
-- @source: Oracle SQL Reference
CREATE INDEX idx_emp_dept ON employees(department_id)

-- @name: create unique index
-- @valid: true
-- @source: Oracle SQL Reference
CREATE UNIQUE INDEX uk_emp_email ON employees(email)

-- @name: create bitmap index
-- @valid: true
-- @source: Oracle SQL Reference
CREATE BITMAP INDEX bx_emp_dept ON employees(department_id)

-- @name: create view
-- @valid: true
-- @source: Oracle SQL Reference
CREATE VIEW emp_view AS
SELECT employee_id, first_name, last_name, department_id FROM employees

-- @name: create or replace view
-- @valid: true
-- @source: Oracle SQL Reference
CREATE OR REPLACE VIEW dept_summary AS
SELECT d.department_id, d.department_name, COUNT(e.employee_id) AS emp_count
FROM departments d LEFT JOIN employees e ON d.department_id = e.department_id
GROUP BY d.department_id, d.department_name

-- @name: create materialized view
-- @valid: true
-- @source: Oracle SQL Reference
CREATE MATERIALIZED VIEW mv_dept_salary
BUILD IMMEDIATE
REFRESH COMPLETE ON DEMAND
AS SELECT department_id, SUM(salary) AS total_sal FROM employees GROUP BY department_id

-- @name: create sequence
-- @valid: true
-- @source: Oracle SQL Reference
CREATE SEQUENCE emp_seq START WITH 1000 INCREMENT BY 1 NOCACHE NOCYCLE

-- @name: create synonym
-- @valid: true
-- @source: Oracle SQL Reference
CREATE SYNONYM emp FOR hr.employees

-- @name: create public synonym
-- @valid: true
-- @source: Oracle SQL Reference
CREATE PUBLIC SYNONYM all_employees FOR hr.employees

-- @name: create type object
-- @valid: true
-- @source: Oracle SQL Reference
CREATE TYPE address_type AS OBJECT (
  street VARCHAR2(200),
  city VARCHAR2(100),
  state VARCHAR2(50),
  zip VARCHAR2(20)
)

-- @name: create type collection
-- @valid: true
-- @source: Oracle SQL Reference
CREATE TYPE number_list AS TABLE OF NUMBER

-- @name: create type varray
-- @valid: true
-- @source: Oracle SQL Reference
CREATE TYPE phone_list AS VARRAY(5) OF VARCHAR2(20)

-- @name: alter table add column
-- @valid: true
-- @source: Oracle SQL Reference
ALTER TABLE employees ADD (middle_name VARCHAR2(50))

-- @name: alter table modify column
-- @valid: true
-- @source: Oracle SQL Reference
ALTER TABLE employees MODIFY (first_name VARCHAR2(200))

-- @name: alter table drop column
-- @valid: true
-- @source: Oracle SQL Reference
ALTER TABLE employees DROP COLUMN middle_name

-- @name: alter table add constraint
-- @valid: true
-- @source: Oracle SQL Reference
ALTER TABLE employees ADD CONSTRAINT uk_email UNIQUE (email)

-- @name: alter table drop constraint
-- @valid: true
-- @source: Oracle SQL Reference
ALTER TABLE employees DROP CONSTRAINT uk_email

-- @name: alter table rename column
-- @valid: true
-- @source: Oracle SQL Reference
ALTER TABLE employees RENAME COLUMN first_name TO given_name

-- @name: alter table rename
-- @valid: true
-- @source: Oracle SQL Reference
ALTER TABLE old_table RENAME TO new_table

-- @name: alter index rebuild
-- @valid: true
-- @source: Oracle SQL Reference
ALTER INDEX idx_emp_dept REBUILD

-- @name: drop table
-- @valid: true
-- @source: Oracle SQL Reference
DROP TABLE test_table

-- @name: drop table cascade constraints
-- @valid: true
-- @source: Oracle SQL Reference
DROP TABLE orders CASCADE CONSTRAINTS

-- @name: drop table purge
-- @valid: true
-- @source: Oracle SQL Reference
DROP TABLE temp_data PURGE

-- @name: drop index
-- @valid: true
-- @source: Oracle SQL Reference
DROP INDEX idx_emp_dept

-- @name: drop view
-- @valid: true
-- @source: Oracle SQL Reference
DROP VIEW emp_view

-- @name: drop sequence
-- @valid: true
-- @source: Oracle SQL Reference
DROP SEQUENCE emp_seq

-- @name: drop synonym
-- @valid: true
-- @source: Oracle SQL Reference
DROP SYNONYM emp

-- @name: drop type
-- @valid: true
-- @source: Oracle SQL Reference
DROP TYPE address_type

-- @name: truncate table
-- @valid: true
-- @source: Oracle SQL Reference
TRUNCATE TABLE temp_results

-- @name: comment on table
-- @valid: true
-- @source: Oracle SQL Reference
COMMENT ON TABLE employees IS 'Employee master table'

-- @name: comment on column
-- @valid: true
-- @source: Oracle SQL Reference
COMMENT ON COLUMN employees.salary IS 'Annual salary in USD'
