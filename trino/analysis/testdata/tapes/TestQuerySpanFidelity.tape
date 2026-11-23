-- pin: default trinodb/trino:482 sha256:90b35b7c603eaa1f889bf03981a62b75f998ee6c0f851d9f4e341b49a57022b6 482

== TestQuerySpanFidelity default#1
> raw
| fixture
| CREATE SCHEMA memory.sales;
| CREATE SCHEMA memory.hr;
| USE memory.sales;
| CREATE TABLE orders (id bigint, customer_id bigint, amount decimal(10, 2), status varchar, placed_at date);
| CREATE TABLE customers (id bigint, name varchar, region varchar, email varchar);
| CREATE TABLE items (order_id bigint, sku varchar, qty integer, tags array(varchar));
| CREATE TABLE "MixedCase" ("CamelCol" integer, "lower_col" varchar);
| CREATE VIEW big_orders AS SELECT id, customer_id, amount * 2 AS doubled FROM orders WHERE amount > 100;
| CREATE VIEW big_order_names AS SELECT b.id, c.name, b.doubled FROM big_orders b JOIN customers c ON b.customer_id = c.id;
| USE memory.hr;
| CREATE TABLE employees (id bigint, name varchar, dept_id integer, salary decimal(10, 2), manager_id bigint);
| CREATE TABLE depts (id integer, name varchar, region varchar);
| CREATE TABLE customers (id bigint, name varchar, vip boolean);
| CREATE VIEW dept_pay AS SELECT d.name AS dept, sum(e.salary) AS total FROM employees e JOIN depts d ON e.dept_id = d.id GROUP BY d.name;
| CREATE VIEW sales_customers AS SELECT name, region FROM memory.sales.customers;
| CREATE VIEW "Vip_Names" AS SELECT c.name AS "VipName" FROM customers c WHERE c.vip;
< raw
| {"schemas":[{"name":"default","tables":null,"views":null},{"name":"hr","tables":[{"name":"customers","columns":[{"name":"id","type":"bigint"},{"name":"name","type":"varchar"},{"name":"vip","type":"boolean"}]},{"name":"depts","columns":[{"name":"id","type":"integer"},{"name":"name","type":"varchar"},{"name":"region","type":"varchar"}]},{"name":"employees","columns":[{"name":"id","type":"bigint"},{"name":"name","type":"varchar"},{"name":"dept_id","type":"integer"},{"name":"salary","type":"decimal(10,2)"},{"name":"manager_id","type":"bigint"}]}],"views":[{"name":"dept_pay","columns":[{"name":"dept","type":"varchar"},{"name":"total","type":"decimal(38,2)"}],"definition":"SELECT\n  d.name dept\n, sum(e.salary) total\nFROM\n  (employees e\nINNER JOIN depts d ON (e.dept_id = d.id))\nGROUP BY d.name\n"},{"name":"sales_customers","columns":[{"name":"name","type":"varchar"},{"name":"region","type":"varchar"}],"definition":"SELECT\n  name\n, region\nFROM\n  memory.sales.customers\n"},{"name":"vip_names","columns":[{"name":"vipname","type":"varchar"}],"definition":"SELECT c.name \"VipName\"\nFROM\n  customers c\nWHERE c.vip\n"}]},{"name":"sales","tables":[{"name":"customers","columns":[{"name":"id","type":"bigint"},{"name":"name","type":"varchar"},{"name":"region","type":"varchar"},{"name":"email","type":"varchar"}]},{"name":"items","columns":[{"name":"order_id","type":"bigint"},{"name":"sku","type":"varchar"},{"name":"qty","type":"integer"},{"name":"tags","type":"array(varchar)"}]},{"name":"mixedcase","columns":[{"name":"camelcol","type":"integer"},{"name":"lower_col","type":"varchar"}]},{"name":"orders","columns":[{"name":"id","type":"bigint"},{"name":"customer_id","type":"bigint"},{"name":"amount","type":"decimal(10,2)"},{"name":"status","type":"varchar"},{"name":"placed_at","type":"date"}]}],"views":[{"name":"big_order_names","columns":[{"name":"id","type":"bigint"},{"name":"name","type":"varchar"},{"name":"doubled","type":"decimal(20,2)"}],"definition":"SELECT\n  b.id\n, c.name\n, b.doubled\nFROM\n  (big_orders b\nINNER JOIN customers c ON (b.customer_id = c.id))\n"},{"name":"big_orders","columns":[{"name":"id","type":"bigint"},{"name":"customer_id","type":"bigint"},{"name":"doubled","type":"decimal(20,2)"}],"definition":"SELECT\n  id\n, customer_id\n, (amount * 2) doubled\nFROM\n  orders\nWHERE (amount \u003e 100)\n"}]}],"views":{"memory.hr.dept_pay":[["memory.hr.depts.name"],["memory.hr.employees.salary"]],"memory.hr.sales_customers":[["memory.sales.customers.name"],["memory.sales.customers.region"]],"memory.hr.vip_names":[["memory.hr.customers.name"]],"memory.sales.big_order_names":[["memory.sales.big_orders.id"],["memory.sales.customers.name"],["memory.sales.big_orders.doubled"]],"memory.sales.big_orders":[["memory.sales.orders.id"],["memory.sales.orders.customer_id"],["memory.sales.orders.amount"]]}}
> raw
| span memory.sales
| SELECT id, amount FROM orders
< raw
| {"columns":["id","amount"],"tables":[{"name":"memory.sales.orders"}],"lineage":[["memory.sales.orders.id"],["memory.sales.orders.amount"]]}
> raw
| span memory.sales
| SELECT o.id AS order_id, o.amount FROM orders o
< raw
| {"columns":["order_id","amount"],"tables":[{"name":"memory.sales.orders"}],"lineage":[["memory.sales.orders.id"],["memory.sales.orders.amount"]]}
> raw
| span memory.sales
| SELECT id + customer_id, upper(status) FROM orders
< raw
| {"columns":["_col0","_col1"],"tables":[{"name":"memory.sales.orders"}],"lineage":[["memory.sales.orders.id","memory.sales.orders.customer_id"],["memory.sales.orders.status"]]}
> raw
| span memory.sales
| SELECT CAST(amount AS varchar) AS a, try_cast(status AS integer) FROM orders
< raw
| {"columns":["a","_col1"],"tables":[{"name":"memory.sales.orders"}],"lineage":[["memory.sales.orders.amount"],["memory.sales.orders.status"]]}
> raw
| span memory.sales
| SELECT 1, 'a' AS b
< raw
| {"columns":["_col0","b"],"lineage":[[],[]]}
> raw
| span memory.sales
| SELECT count(*) FROM orders
< raw
| {"columns":["_col0"],"tables":[{"name":"memory.sales.orders"}],"lineage":[[]]}
> raw
| span memory.sales
| SELECT DISTINCT region FROM customers
< raw
| {"columns":["region"],"tables":[{"name":"memory.sales.customers"}],"lineage":[["memory.sales.customers.region"]]}
> raw
| span memory.sales
| SELECT name FROM customers ORDER BY region LIMIT 5
< raw
| {"columns":["name"],"tables":[{"name":"memory.sales.customers"}],"lineage":[["memory.sales.customers.name"]]}
> raw
| span memory.sales
| SELECT id FROM orders OFFSET 1 ROWS FETCH FIRST 2 ROWS ONLY
< raw
| {"columns":["id"],"tables":[{"name":"memory.sales.orders"}],"lineage":[["memory.sales.orders.id"]]}
> raw
| span memory.sales
| TABLE customers
< raw
| {"columns":["id","name","region","email"],"tables":[{"name":"memory.sales.customers"}],"lineage":[["memory.sales.customers.id"],["memory.sales.customers.name"],["memory.sales.customers.region"],["memory.sales.customers.email"]]}
> raw
| span memory.sales
| SELECT * FROM orders
< raw
| {"columns":["id","customer_id","amount","status","placed_at"],"tables":[{"name":"memory.sales.orders"}],"lineage":[["memory.sales.orders.id"],["memory.sales.orders.customer_id"],["memory.sales.orders.amount"],["memory.sales.orders.status"],["memory.sales.orders.placed_at"]]}
> raw
| span memory.sales
| SELECT o.* FROM orders o
< raw
| {"columns":["id","customer_id","amount","status","placed_at"],"tables":[{"name":"memory.sales.orders"}],"lineage":[["memory.sales.orders.id"],["memory.sales.orders.customer_id"],["memory.sales.orders.amount"],["memory.sales.orders.status"],["memory.sales.orders.placed_at"]]}
> raw
| span memory.sales
| SELECT *, amount FROM orders
< raw
| {"columns":["id","customer_id","amount","status","placed_at","amount"],"tables":[{"name":"memory.sales.orders"}],"lineage":[["memory.sales.orders.id"],["memory.sales.orders.customer_id"],["memory.sales.orders.amount"],["memory.sales.orders.status"],["memory.sales.orders.placed_at"],["memory.sales.orders.amount"]]}
> raw
| span memory.sales
| SELECT c.*, o.amount FROM customers c JOIN orders o ON o.customer_id = c.id
< raw
| {"columns":["id","name","region","email","amount"],"tables":[{"name":"memory.sales.customers"},{"name":"memory.sales.orders"}],"lineage":[["memory.sales.customers.id"],["memory.sales.customers.name"],["memory.sales.customers.region"],["memory.sales.customers.email"],["memory.sales.orders.amount"]]}
> raw
| span memory.sales
| SELECT * FROM (SELECT * FROM customers) c
< raw
| {"columns":["id","name","region","email"],"tables":[{"name":"memory.sales.customers"}],"lineage":[["memory.sales.customers.id"],["memory.sales.customers.name"],["memory.sales.customers.region"],["memory.sales.customers.email"]]}
> raw
| span memory.sales
| SELECT id FROM sales.orders
< raw
| {"columns":["id"],"tables":[{"name":"memory.sales.orders"}],"lineage":[["memory.sales.orders.id"]]}
> raw
| span memory.sales
| SELECT id FROM memory.sales.orders
< raw
| {"columns":["id"],"tables":[{"name":"memory.sales.orders"}],"lineage":[["memory.sales.orders.id"]]}
> raw
| span memory.sales
| SELECT orders.id, customers.name FROM sales.orders JOIN memory.sales.customers ON orders.customer_id = customers.id
< raw
| {"columns":["id","name"],"tables":[{"name":"memory.sales.orders"},{"name":"memory.sales.customers"}],"lineage":[["memory.sales.orders.id"],["memory.sales.customers.name"]]}
> raw
| span memory.sales
| SELECT ID, Amount FROM ORDERS
< raw
| {"columns":["ID","Amount"],"tables":[{"name":"memory.sales.orders"}],"lineage":[["memory.sales.orders.id"],["memory.sales.orders.amount"]]}
> raw
| span memory.sales
| SELECT "id", "AMOUNT" FROM "orders"
< raw
| {"columns":["id","AMOUNT"],"tables":[{"name":"memory.sales.orders"}],"lineage":[["memory.sales.orders.id"],["memory.sales.orders.amount"]]}
> raw
| span memory.sales
| SELECT O.ID FROM Orders o
< raw
| {"columns":["ID"],"tables":[{"name":"memory.sales.orders"}],"lineage":[["memory.sales.orders.id"]]}
> raw
| span memory.sales
| SELECT camelcol, "CamelCol", "LOWER_COL" FROM "MixedCase"
< raw
| {"columns":["camelcol","CamelCol","LOWER_COL"],"tables":[{"name":"memory.sales.mixedcase"}],"lineage":[["memory.sales.mixedcase.camelcol"],["memory.sales.mixedcase.camelcol"],["memory.sales.mixedcase.lower_col"]]}
> raw
| span memory.sales
| SELECT * FROM "MixedCase"
< raw
| {"columns":["camelcol","lower_col"],"tables":[{"name":"memory.sales.mixedcase"}],"lineage":[["memory.sales.mixedcase.camelcol"],["memory.sales.mixedcase.lower_col"]]}
> raw
| span memory.sales
| SELECT * FROM mixedcase
< raw
| {"columns":["camelcol","lower_col"],"tables":[{"name":"memory.sales.mixedcase"}],"lineage":[["memory.sales.mixedcase.camelcol"],["memory.sales.mixedcase.lower_col"]]}
> raw
| span memory.sales
| SELECT m.CAMELCOL FROM MIXEDCASE m
< raw
| {"columns":["CAMELCOL"],"tables":[{"name":"memory.sales.mixedcase"}],"lineage":[["memory.sales.mixedcase.camelcol"]]}
> raw
| span memory.sales
| SELECT id AS Foo, "Foo"."AMOUNT" AS "Bar" FROM orders AS "Foo"
< raw
| {"columns":["Foo","Bar"],"tables":[{"name":"memory.sales.orders"}],"lineage":[["memory.sales.orders.id"],["memory.sales.orders.amount"]]}
> raw
| span memory.sales
| SELECT o.id, c.name FROM orders o JOIN customers c ON o.customer_id = c.id
< raw
| {"columns":["id","name"],"tables":[{"name":"memory.sales.orders"},{"name":"memory.sales.customers"}],"lineage":[["memory.sales.orders.id"],["memory.sales.customers.name"]]}
> raw
| span memory.sales
| SELECT o.id, c.name FROM orders o LEFT JOIN customers c ON o.customer_id = c.id WHERE c.region = 'EU'
< raw
| {"columns":["id","name"],"tables":[{"name":"memory.sales.orders"},{"name":"memory.sales.customers"}],"lineage":[["memory.sales.orders.id"],["memory.sales.customers.name"]]}
> raw
| span memory.sales
| SELECT * FROM orders JOIN items ON orders.id = items.order_id
< raw
| {"columns":["id","customer_id","amount","status","placed_at","order_id","sku","qty","tags"],"tables":[{"name":"memory.sales.orders"},{"name":"memory.sales.items"}],"lineage":[["memory.sales.orders.id"],["memory.sales.orders.customer_id"],["memory.sales.orders.amount"],["memory.sales.orders.status"],["memory.sales.orders.placed_at"],["memory.sales.items.order_id"],["memory.sales.items.sku"],["memory.sales.items.qty"],["memory.sales.items.tags"]]}
> raw
| span memory.sales
| SELECT * FROM customers JOIN hr.customers USING (id)
< raw
| {"columns":["id","name","region","email","name","vip"],"tables":[{"name":"memory.sales.customers"},{"name":"memory.hr.customers"}],"lineage":[[],["memory.sales.customers.name"],["memory.sales.customers.region"],["memory.sales.customers.email"],["memory.hr.customers.name"],["memory.hr.customers.vip"]]}
> raw
| span memory.sales
| SELECT id, h.vip FROM customers c JOIN hr.customers h USING (id)
< raw
| {"columns":["id","vip"],"tables":[{"name":"memory.sales.customers"},{"name":"memory.hr.customers"}],"lineage":[[],["memory.hr.customers.vip"]]}
> raw
| span memory.sales
| SELECT c.name, count(*) FROM customers c CROSS JOIN orders o GROUP BY c.name
< raw
| {"columns":["name","_col1"],"tables":[{"name":"memory.sales.customers"},{"name":"memory.sales.orders"}],"lineage":[["memory.sales.customers.name"],[]]}
> raw
| span memory.sales
| SELECT e.name, m.name AS manager FROM hr.employees e LEFT JOIN hr.employees m ON e.manager_id = m.id
< raw
| {"columns":["name","manager"],"tables":[{"name":"memory.hr.employees"},{"name":"memory.hr.employees"}],"lineage":[["memory.hr.employees.name"],["memory.hr.employees.name"]]}
> raw
| span memory.sales
| SELECT c.name, h.name FROM customers c JOIN hr.customers h ON c.id = h.id
< raw
| {"columns":["name","name"],"tables":[{"name":"memory.sales.customers"},{"name":"memory.hr.customers"}],"lineage":[["memory.sales.customers.name"],["memory.hr.customers.name"]]}
> raw
| span memory.sales
| SELECT region, vip FROM customers JOIN hr.customers h ON customers.id = h.id
< raw
| {"columns":["region","vip"],"tables":[{"name":"memory.sales.customers"},{"name":"memory.hr.customers"}],"lineage":[["memory.sales.customers.region"],["memory.hr.customers.vip"]]}
> raw
| span memory.sales
| SELECT x.id, x.total FROM (SELECT id, amount * 2 AS total FROM orders) x
< raw
| {"columns":["id","total"],"tables":[{"name":"memory.sales.orders"}],"lineage":[["memory.sales.orders.id"],["memory.sales.orders.amount"]]}
> raw
| span memory.sales
| SELECT * FROM (SELECT id, amount FROM orders) AS x(a, b)
< raw
| {"columns":["a","b"],"tables":[{"name":"memory.sales.orders"}],"lineage":[["memory.sales.orders.id"],["memory.sales.orders.amount"]]}
> raw
| span memory.sales
| SELECT name FROM customers WHERE id IN (SELECT customer_id FROM orders WHERE amount > 10)
< raw
| {"columns":["name"],"tables":[{"name":"memory.sales.customers"},{"name":"memory.sales.orders"}],"lineage":[["memory.sales.customers.name"]]}
> raw
| span memory.sales
| SELECT name, (SELECT max(amount) FROM orders o WHERE o.customer_id = c.id) AS top FROM customers c
< raw
| {"columns":["name","top"],"tables":[{"name":"memory.sales.customers"},{"name":"memory.sales.orders"}],"lineage":[["memory.sales.customers.name"],["memory.sales.orders.amount"]]}
> raw
| span memory.sales
| SELECT name FROM customers c WHERE EXISTS (SELECT 1 FROM orders o WHERE o.customer_id = c.id)
< raw
| {"columns":["name"],"tables":[{"name":"memory.sales.customers"},{"name":"memory.sales.orders"}],"lineage":[["memory.sales.customers.name"]]}
> raw
| span memory.sales
| SELECT id, (SELECT name FROM customers c WHERE c.id = b.customer_id) FROM big_orders b
< raw
| {"columns":["id","_col1"],"tables":[{"name":"memory.sales.orders"},{"name":"memory.sales.big_orders","view":true},{"name":"memory.sales.customers"}],"lineage":[["memory.sales.big_orders.id"],["memory.sales.customers.name"]]}
> raw
| span memory.sales
| SELECT name FROM customers WHERE region IN (SELECT region FROM hr.depts)
< raw
| {"columns":["name"],"tables":[{"name":"memory.sales.customers"},{"name":"memory.hr.depts"}],"lineage":[["memory.sales.customers.name"]]}
> raw
| span memory.sales
| SELECT (SELECT max(doubled) FROM big_orders) AS top, (SELECT count(*) FROM hr.employees) AS n
< raw
| {"columns":["top","n"],"tables":[{"name":"memory.sales.orders"},{"name":"memory.sales.big_orders","view":true},{"name":"memory.hr.employees"}],"lineage":[["memory.sales.big_orders.doubled"],[]]}
> raw
| span memory.sales
| SELECT (SELECT name FROM hr.depts) AS n FROM customers AS depts
< raw
| {"columns":["n"],"tables":[{"name":"memory.sales.customers"},{"name":"memory.hr.depts"}],"lineage":[["memory.hr.depts.name"]]}
> raw
| span memory.sales
| WITH big AS (SELECT id, amount FROM orders WHERE amount > 10) SELECT id, amount FROM big
< raw
| {"columns":["id","amount"],"tables":[{"name":"memory.sales.orders"}],"lineage":[["memory.sales.orders.id"],["memory.sales.orders.amount"]]}
> raw
| span memory.sales
| WITH a AS (SELECT id, name FROM customers), b AS (SELECT a.name, o.amount FROM a JOIN orders o ON o.customer_id = a.id) SELECT * FROM b
< raw
| {"columns":["name","amount"],"tables":[{"name":"memory.sales.customers"},{"name":"memory.sales.orders"}],"lineage":[["memory.sales.customers.name"],["memory.sales.orders.amount"]]}
> raw
| span memory.sales
| WITH orders AS (SELECT id, name AS status FROM customers) SELECT status FROM orders
< raw
| {"columns":["status"],"tables":[{"name":"memory.sales.customers"}],"lineage":[["memory.sales.customers.name"]]}
> raw
| span memory.sales
| WITH t(x, y) AS (SELECT id, name FROM customers) SELECT y FROM t
< raw
| {"columns":["y"],"tables":[{"name":"memory.sales.customers"}],"lineage":[["memory.sales.customers.name"]]}
> raw
| span memory.sales
| WITH RECURSIVE r(n) AS (SELECT 1 UNION ALL SELECT n + 1 FROM r WHERE n < 3) SELECT n FROM r
< raw
| {"columns":["n"],"lineage":[[]]}
> raw
| span memory.sales
| SELECT id, name FROM customers UNION SELECT id, name FROM hr.customers
< raw
| {"columns":["id","name"],"tables":[{"name":"memory.sales.customers"},{"name":"memory.hr.customers"}],"lineage":[["memory.sales.customers.id","memory.hr.customers.id"],["memory.sales.customers.name","memory.hr.customers.name"]]}
> raw
| span memory.sales
| SELECT id FROM orders INTERSECT SELECT customer_id FROM orders
< raw
| {"columns":["id"],"tables":[{"name":"memory.sales.orders"},{"name":"memory.sales.orders"}],"lineage":[["memory.sales.orders.id","memory.sales.orders.customer_id"]]}
> raw
| span memory.sales
| SELECT name FROM customers EXCEPT SELECT name FROM hr.employees
< raw
| {"columns":["name"],"tables":[{"name":"memory.sales.customers"},{"name":"memory.hr.employees"}],"lineage":[["memory.sales.customers.name","memory.hr.employees.name"]]}
> raw
| span memory.sales
| SELECT id FROM orders UNION ALL SELECT id FROM customers UNION ALL SELECT order_id FROM items
< raw
| {"columns":["id"],"tables":[{"name":"memory.sales.orders"},{"name":"memory.sales.customers"},{"name":"memory.sales.items"}],"lineage":[["memory.sales.orders.id","memory.sales.customers.id","memory.sales.items.order_id"]]}
> raw
| span memory.sales
| SELECT * FROM (SELECT id AS k FROM orders UNION SELECT id FROM customers) u
< raw
| {"columns":["k"],"tables":[{"name":"memory.sales.orders"},{"name":"memory.sales.customers"}],"lineage":[["memory.sales.orders.id","memory.sales.customers.id"]]}
> raw
| span memory.sales
| SELECT tag FROM items CROSS JOIN UNNEST(tags) AS t(tag)
< raw
| {"columns":["tag"],"tables":[{"name":"memory.sales.items"}],"lineage":[["memory.sales.items.tags"]]}
> raw
| span memory.sales
| SELECT o.id, t.tag FROM items i CROSS JOIN UNNEST(i.tags) AS t(tag) JOIN orders o ON o.id = i.order_id
< raw
| {"columns":["id","tag"],"tables":[{"name":"memory.sales.items"},{"name":"memory.sales.orders"}],"lineage":[["memory.sales.orders.id"],["memory.sales.items.tags"]]}
> raw
| span memory.sales
| SELECT tag, n FROM items CROSS JOIN UNNEST(tags) WITH ORDINALITY AS t(tag, n)
< raw
| {"columns":["tag","n"],"tables":[{"name":"memory.sales.items"}],"lineage":[["memory.sales.items.tags"],[]]}
> raw
| span memory.sales
| SELECT * FROM UNNEST(ARRAY[1, 2]) AS t(x)
< raw
| {"columns":["x"],"lineage":[[]]}
> raw
| span memory.sales
| SELECT * FROM (VALUES (1, 'a'), (2, 'b')) AS v(x, y)
< raw
| {"columns":["x","y"],"lineage":[[],[]]}
> raw
| span memory.sales
| SELECT transform(tags, x -> upper(x)) FROM items
< raw
| {"columns":["_col0"],"tables":[{"name":"memory.sales.items"}],"lineage":[["memory.sales.items.tags"]]}
> raw
| span memory.sales
| SELECT c.name, l.amount FROM customers c, LATERAL (SELECT amount FROM orders o WHERE o.customer_id = c.id) l
< raw
| {"columns":["name","amount"],"tables":[{"name":"memory.sales.customers"},{"name":"memory.sales.orders"}],"lineage":[["memory.sales.customers.name"],["memory.sales.orders.amount"]]}
> raw
| span memory.sales
| SELECT c.name, l.total FROM customers c CROSS JOIN LATERAL (SELECT sum(amount) AS total FROM orders o WHERE o.customer_id = c.id) l
< raw
| {"columns":["name","total"],"tables":[{"name":"memory.sales.customers"},{"name":"memory.sales.orders"}],"lineage":[["memory.sales.customers.name"],["memory.sales.orders.amount"]]}
> raw
| span memory.sales
| SELECT c.name, l.n FROM customers c CROSS JOIN LATERAL (SELECT c.region AS n) l
< raw
| {"columns":["name","n"],"tables":[{"name":"memory.sales.customers"}],"lineage":[["memory.sales.customers.name"],["memory.sales.customers.region"]]}
> raw
| span memory.sales
| SELECT * FROM customers c, LATERAL (SELECT amount, status FROM orders o WHERE o.customer_id = c.id) AS l(a, s)
< raw
| {"columns":["id","name","region","email","a","s"],"tables":[{"name":"memory.sales.customers"},{"name":"memory.sales.orders"}],"lineage":[["memory.sales.customers.id"],["memory.sales.customers.name"],["memory.sales.customers.region"],["memory.sales.customers.email"],["memory.sales.orders.amount"],["memory.sales.orders.status"]]}
> raw
| span memory.sales
| SELECT id, row_number() OVER (PARTITION BY customer_id ORDER BY placed_at) AS rn FROM orders
< raw
| {"columns":["id","rn"],"tables":[{"name":"memory.sales.orders"}],"lineage":[["memory.sales.orders.id"],["memory.sales.orders.customer_id","memory.sales.orders.placed_at"]]}
> raw
| span memory.sales
| SELECT id, sum(amount) OVER w FROM orders WINDOW w AS (PARTITION BY status)
< raw
| {"columns":["id","_col1"],"tables":[{"name":"memory.sales.orders"}],"lineage":[["memory.sales.orders.id"],["memory.sales.orders.amount"]]}
> raw
| span memory.sales
| SELECT id, rank() OVER (ORDER BY amount DESC) FROM orders
< raw
| {"columns":["id","_col1"],"tables":[{"name":"memory.sales.orders"}],"lineage":[["memory.sales.orders.id"],["memory.sales.orders.amount"]]}
> raw
| span memory.sales
| SELECT status, count(*) AS n, sum(amount) FROM orders GROUP BY status
< raw
| {"columns":["status","n","_col2"],"tables":[{"name":"memory.sales.orders"}],"lineage":[["memory.sales.orders.status"],[],["memory.sales.orders.amount"]]}
> raw
| span memory.sales
| SELECT status, customer_id, sum(amount) FROM orders GROUP BY ROLLUP (status, customer_id)
< raw
| {"columns":["status","customer_id","_col2"],"tables":[{"name":"memory.sales.orders"}],"lineage":[["memory.sales.orders.status"],["memory.sales.orders.customer_id"],["memory.sales.orders.amount"]]}
> raw
| span memory.sales
| SELECT region, count(*) FROM customers GROUP BY CUBE (region)
< raw
| {"columns":["region","_col1"],"tables":[{"name":"memory.sales.customers"}],"lineage":[["memory.sales.customers.region"],[]]}
> raw
| span memory.sales
| SELECT status, grouping(status) FROM orders GROUP BY GROUPING SETS ((status), ())
< raw
| {"columns":["status","_col1"],"tables":[{"name":"memory.sales.orders"}],"lineage":[["memory.sales.orders.status"],["memory.sales.orders.status"]]}
> raw
| span memory.sales
| SELECT status FROM orders GROUP BY status HAVING sum(amount) > 100
< raw
| {"columns":["status"],"tables":[{"name":"memory.sales.orders"}],"lineage":[["memory.sales.orders.status"]]}
> raw
| span memory.sales
| SELECT count(*) FILTER (WHERE amount > 10) FROM orders
< raw
| {"columns":["_col0"],"tables":[{"name":"memory.sales.orders"}],"lineage":[["memory.sales.orders.amount"]]}
> raw
| span memory.sales
| SELECT max_by(id, amount) FROM orders
< raw
| {"columns":["_col0"],"tables":[{"name":"memory.sales.orders"}],"lineage":[["memory.sales.orders.id","memory.sales.orders.amount"]]}
> raw
| span memory.sales
| SELECT CASE WHEN amount > 100 THEN 'big' ELSE status END AS kind FROM orders
< raw
| {"columns":["kind"],"tables":[{"name":"memory.sales.orders"}],"lineage":[["memory.sales.orders.amount","memory.sales.orders.status"]]}
> raw
| span memory.sales
| SELECT CASE status WHEN 'x' THEN id END FROM orders
< raw
| {"columns":["_col0"],"tables":[{"name":"memory.sales.orders"}],"lineage":[["memory.sales.orders.status","memory.sales.orders.id"]]}
> raw
| span memory.sales
| SELECT coalesce(c.email, c.name) AS contact, if(o.amount > 0, o.id) FROM orders o JOIN customers c ON o.customer_id = c.id
< raw
| {"columns":["contact","_col1"],"tables":[{"name":"memory.sales.orders"},{"name":"memory.sales.customers"}],"lineage":[["memory.sales.customers.email","memory.sales.customers.name"],["memory.sales.orders.amount","memory.sales.orders.id"]]}
> raw
| span memory.sales
| SELECT * FROM big_orders
< raw
| {"columns":["id","customer_id","doubled"],"tables":[{"name":"memory.sales.orders"},{"name":"memory.sales.big_orders","view":true}],"lineage":[["memory.sales.big_orders.id"],["memory.sales.big_orders.customer_id"],["memory.sales.big_orders.doubled"]]}
> raw
| span memory.sales
| SELECT id, doubled FROM big_orders WHERE customer_id > 0
< raw
| {"columns":["id","doubled"],"tables":[{"name":"memory.sales.orders"},{"name":"memory.sales.big_orders","view":true}],"lineage":[["memory.sales.big_orders.id"],["memory.sales.big_orders.doubled"]]}
> raw
| span memory.sales
| SELECT * FROM big_order_names
< raw
| {"columns":["id","name","doubled"],"tables":[{"name":"memory.sales.orders"},{"name":"memory.sales.big_orders","view":true},{"name":"memory.sales.customers"},{"name":"memory.sales.big_order_names","view":true}],"lineage":[["memory.sales.big_order_names.id"],["memory.sales.big_order_names.name"],["memory.sales.big_order_names.doubled"]]}
> raw
| span memory.sales
| SELECT b.name FROM big_order_names b JOIN orders o ON b.id = o.id
< raw
| {"columns":["name"],"tables":[{"name":"memory.sales.orders"},{"name":"memory.sales.big_orders","view":true},{"name":"memory.sales.customers"},{"name":"memory.sales.big_order_names","view":true},{"name":"memory.sales.orders"}],"lineage":[["memory.sales.big_order_names.name"]]}
> raw
| span memory.sales
| SELECT * FROM hr.dept_pay
< raw
| {"columns":["dept","total"],"tables":[{"name":"memory.hr.employees"},{"name":"memory.hr.depts"},{"name":"memory.hr.dept_pay","view":true}],"lineage":[["memory.hr.dept_pay.dept"],["memory.hr.dept_pay.total"]]}
> raw
| span memory.sales
| SELECT * FROM hr.sales_customers
< raw
| {"columns":["name","region"],"tables":[{"name":"memory.sales.customers"},{"name":"memory.hr.sales_customers","view":true}],"lineage":[["memory.hr.sales_customers.name"],["memory.hr.sales_customers.region"]]}
> raw
| span memory.sales
| SELECT * FROM hr."Vip_Names"
< raw
| {"columns":["VipName"],"tables":[{"name":"memory.hr.customers"},{"name":"memory.hr.vip_names","view":true}],"lineage":[["memory.hr.vip_names.VipName"]]}
> raw
| span memory.sales
| SELECT vipname FROM hr.vip_names
< raw
| {"columns":["vipname"],"tables":[{"name":"memory.hr.customers"},{"name":"memory.hr.vip_names","view":true}],"lineage":[["memory.hr.vip_names.VipName"]]}
> raw
| span memory.hr
| SELECT id, name FROM customers
< raw
| {"columns":["id","name"],"tables":[{"name":"memory.hr.customers"}],"lineage":[["memory.hr.customers.id"],["memory.hr.customers.name"]]}
> raw
| span memory.hr
| SELECT * FROM customers
< raw
| {"columns":["id","name","vip"],"tables":[{"name":"memory.hr.customers"}],"lineage":[["memory.hr.customers.id"],["memory.hr.customers.name"],["memory.hr.customers.vip"]]}
> raw
| span memory.hr
| SELECT e.name, d.name FROM employees e JOIN depts d ON e.dept_id = d.id
< raw
| {"columns":["name","name"],"tables":[{"name":"memory.hr.employees"},{"name":"memory.hr.depts"}],"lineage":[["memory.hr.employees.name"],["memory.hr.depts.name"]]}
> raw
| span memory.hr
| SELECT * FROM dept_pay
< raw
| {"columns":["dept","total"],"tables":[{"name":"memory.hr.employees"},{"name":"memory.hr.depts"},{"name":"memory.hr.dept_pay","view":true}],"lineage":[["memory.hr.dept_pay.dept"],["memory.hr.dept_pay.total"]]}
> raw
| span memory.hr
| SELECT name FROM sales.customers
< raw
| {"columns":["name"],"tables":[{"name":"memory.sales.customers"}],"lineage":[["memory.sales.customers.name"]]}
> raw
| span memory.hr
| SELECT c.name, h.vip FROM sales.customers c JOIN customers h ON c.id = h.id
< raw
| {"columns":["name","vip"],"tables":[{"name":"memory.sales.customers"},{"name":"memory.hr.customers"}],"lineage":[["memory.sales.customers.name"],["memory.hr.customers.vip"]]}
> raw
| span memory.hr
| SELECT * FROM sales_customers
< raw
| {"columns":["name","region"],"tables":[{"name":"memory.sales.customers"},{"name":"memory.hr.sales_customers","view":true}],"lineage":[["memory.hr.sales_customers.name"],["memory.hr.sales_customers.region"]]}
> raw
| span memory.hr
| SELECT * FROM sales.big_orders
< raw
| {"columns":["id","customer_id","doubled"],"tables":[{"name":"memory.sales.orders"},{"name":"memory.sales.big_orders","view":true}],"lineage":[["memory.sales.big_orders.id"],["memory.sales.big_orders.customer_id"],["memory.sales.big_orders.doubled"]]}
> raw
| span memory.default
| SELECT id FROM sales.orders
< raw
| {"columns":["id"],"tables":[{"name":"memory.sales.orders"}],"lineage":[["memory.sales.orders.id"]]}
> raw
| span memory.default
| SELECT * FROM hr.employees e JOIN hr.depts d ON e.dept_id = d.id
< raw
| {"columns":["id","name","dept_id","salary","manager_id","id","name","region"],"tables":[{"name":"memory.hr.employees"},{"name":"memory.hr.depts"}],"lineage":[["memory.hr.employees.id"],["memory.hr.employees.name"],["memory.hr.employees.dept_id"],["memory.hr.employees.salary"],["memory.hr.employees.manager_id"],["memory.hr.depts.id"],["memory.hr.depts.name"],["memory.hr.depts.region"]]}
> raw
| span memory.default
| SELECT name FROM memory.hr.customers
< raw
| {"columns":["name"],"tables":[{"name":"memory.hr.customers"}],"lineage":[["memory.hr.customers.name"]]}
> raw
| span memory.default
| SELECT * FROM sales.big_order_names
< raw
| {"columns":["id","name","doubled"],"tables":[{"name":"memory.sales.orders"},{"name":"memory.sales.big_orders","view":true},{"name":"memory.sales.customers"},{"name":"memory.sales.big_order_names","view":true}],"lineage":[["memory.sales.big_order_names.id"],["memory.sales.big_order_names.name"],["memory.sales.big_order_names.doubled"]]}
> raw
| span memory.default
| WITH x AS (SELECT * FROM sales.customers) SELECT x.name, d.name FROM x JOIN hr.depts d ON x.region = d.region
< raw
| {"columns":["name","name"],"tables":[{"name":"memory.sales.customers"},{"name":"memory.hr.depts"}],"lineage":[["memory.sales.customers.name"],["memory.hr.depts.name"]]}
