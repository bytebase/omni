CREATE DATABASE shop
CREATE DATABASE hr
CREATE TABLE shop.customers (id UInt32, name String, email String, country String, created Date) ENGINE = MergeTree ORDER BY id
CREATE TABLE shop.orders (id UInt32, customer_id UInt32, amount Decimal(10, 2), status String, tags Array(String), created DateTime) ENGINE = MergeTree ORDER BY id
CREATE TABLE shop.items (order_id UInt32, sku String, qty UInt32, price Float64) ENGINE = ReplacingMergeTree ORDER BY (order_id, sku)
CREATE VIEW shop.big_orders AS SELECT id, customer_id, amount FROM shop.orders WHERE amount > 100
CREATE VIEW shop.customer_totals AS SELECT c.id AS customer_id, c.name AS name, sum(o.amount) AS total FROM shop.customers AS c INNER JOIN shop.orders AS o ON o.customer_id = c.id GROUP BY c.id, c.name
CREATE VIEW shop.vip AS SELECT customer_id, upper(name) AS name FROM shop.customer_totals WHERE total > 1000
CREATE TABLE hr.employees (id UInt32, name String, salary UInt64, dept_id UInt32, manager_id UInt32) ENGINE = MergeTree ORDER BY id
CREATE TABLE hr.depts (id UInt32, title String) ENGINE = MergeTree ORDER BY id
CREATE VIEW hr.payroll AS SELECT e.name AS name, d.title AS dept, e.salary * 12 AS yearly FROM hr.employees AS e LEFT JOIN hr.depts AS d ON d.id = e.dept_id
