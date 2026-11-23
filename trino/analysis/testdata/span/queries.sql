-- The span.tsv corpus: one query per line, run in the session the last
-- "-- session:" line names on the fixture the last "-- fixture:" line names.
-- See ../README.md.

-- fixture: base
-- session: memory.sales

-- Columns, aliases and expressions.
SELECT id, amount FROM orders
SELECT o.id AS order_id, o.amount FROM orders o
SELECT id + customer_id, upper(status) FROM orders
SELECT CAST(amount AS varchar) AS a, try_cast(status AS integer) FROM orders
SELECT 1, 'a' AS b
SELECT count(*) FROM orders
SELECT DISTINCT region FROM customers
SELECT name FROM customers ORDER BY region LIMIT 5
SELECT id FROM orders OFFSET 1 ROWS FETCH FIRST 2 ROWS ONLY
TABLE customers

-- Stars.
SELECT * FROM orders
SELECT o.* FROM orders o
SELECT *, amount FROM orders
SELECT c.*, o.amount FROM customers c JOIN orders o ON o.customer_id = c.id
SELECT * FROM (SELECT * FROM customers) c

-- Qualified names and case.
SELECT id FROM sales.orders
SELECT id FROM memory.sales.orders
SELECT orders.id, customers.name FROM sales.orders JOIN memory.sales.customers ON orders.customer_id = customers.id
SELECT ID, Amount FROM ORDERS
SELECT "id", "AMOUNT" FROM "orders"
SELECT O.ID FROM Orders o
SELECT camelcol, "CamelCol", "LOWER_COL" FROM "MixedCase"
SELECT * FROM "MixedCase"
SELECT * FROM mixedcase
SELECT m.CAMELCOL FROM MIXEDCASE m
SELECT id AS Foo, "Foo"."AMOUNT" AS "Bar" FROM orders AS "Foo"

-- Joins.
SELECT o.id, c.name FROM orders o JOIN customers c ON o.customer_id = c.id
SELECT o.id, c.name FROM orders o LEFT JOIN customers c ON o.customer_id = c.id WHERE c.region = 'EU'
SELECT * FROM orders JOIN items ON orders.id = items.order_id
SELECT * FROM customers JOIN hr.customers USING (id)
SELECT id, h.vip FROM customers c JOIN hr.customers h USING (id)
SELECT c.name, count(*) FROM customers c CROSS JOIN orders o GROUP BY c.name
SELECT e.name, m.name AS manager FROM hr.employees e LEFT JOIN hr.employees m ON e.manager_id = m.id
SELECT c.name, h.name FROM customers c JOIN hr.customers h ON c.id = h.id
SELECT region, vip FROM customers JOIN hr.customers h ON customers.id = h.id

-- Subqueries.
SELECT x.id, x.total FROM (SELECT id, amount * 2 AS total FROM orders) x
SELECT * FROM (SELECT id, amount FROM orders) AS x(a, b)
SELECT name FROM customers WHERE id IN (SELECT customer_id FROM orders WHERE amount > 10)
SELECT name, (SELECT max(amount) FROM orders o WHERE o.customer_id = c.id) AS top FROM customers c
SELECT name FROM customers c WHERE EXISTS (SELECT 1 FROM orders o WHERE o.customer_id = c.id)
SELECT id, (SELECT name FROM customers c WHERE c.id = b.customer_id) FROM big_orders b
SELECT name FROM customers WHERE region IN (SELECT region FROM hr.depts)
SELECT (SELECT max(doubled) FROM big_orders) AS top, (SELECT count(*) FROM hr.employees) AS n
SELECT (SELECT name FROM hr.depts) AS n FROM customers AS depts

-- Common table expressions.
WITH big AS (SELECT id, amount FROM orders WHERE amount > 10) SELECT id, amount FROM big
WITH a AS (SELECT id, name FROM customers), b AS (SELECT a.name, o.amount FROM a JOIN orders o ON o.customer_id = a.id) SELECT * FROM b
WITH orders AS (SELECT id, name AS status FROM customers) SELECT status FROM orders
WITH t(x, y) AS (SELECT id, name FROM customers) SELECT y FROM t
WITH RECURSIVE r(n) AS (SELECT 1 UNION ALL SELECT n + 1 FROM r WHERE n < 3) SELECT n FROM r

-- Set operations.
SELECT id, name FROM customers UNION SELECT id, name FROM hr.customers
SELECT id FROM orders INTERSECT SELECT customer_id FROM orders
SELECT name FROM customers EXCEPT SELECT name FROM hr.employees
SELECT id FROM orders UNION ALL SELECT id FROM customers UNION ALL SELECT order_id FROM items
SELECT * FROM (SELECT id AS k FROM orders UNION SELECT id FROM customers) u

-- UNNEST and VALUES.
SELECT tag FROM items CROSS JOIN UNNEST(tags) AS t(tag)
SELECT o.id, t.tag FROM items i CROSS JOIN UNNEST(i.tags) AS t(tag) JOIN orders o ON o.id = i.order_id
SELECT tag, n FROM items CROSS JOIN UNNEST(tags) WITH ORDINALITY AS t(tag, n)
SELECT * FROM UNNEST(ARRAY[1, 2]) AS t(x)
SELECT * FROM (VALUES (1, 'a'), (2, 'b')) AS v(x, y)
SELECT transform(tags, x -> upper(x)) FROM items

-- Lateral.
SELECT c.name, l.amount FROM customers c, LATERAL (SELECT amount FROM orders o WHERE o.customer_id = c.id) l
SELECT c.name, l.total FROM customers c CROSS JOIN LATERAL (SELECT sum(amount) AS total FROM orders o WHERE o.customer_id = c.id) l
SELECT c.name, l.n FROM customers c CROSS JOIN LATERAL (SELECT c.region AS n) l
SELECT * FROM customers c, LATERAL (SELECT amount, status FROM orders o WHERE o.customer_id = c.id) AS l(a, s)

-- Windows.
SELECT id, row_number() OVER (PARTITION BY customer_id ORDER BY placed_at) AS rn FROM orders
SELECT id, sum(amount) OVER w FROM orders WINDOW w AS (PARTITION BY status)
SELECT id, rank() OVER (ORDER BY amount DESC) FROM orders

-- Grouping.
SELECT status, count(*) AS n, sum(amount) FROM orders GROUP BY status
SELECT status, customer_id, sum(amount) FROM orders GROUP BY ROLLUP (status, customer_id)
SELECT region, count(*) FROM customers GROUP BY CUBE (region)
SELECT status, grouping(status) FROM orders GROUP BY GROUPING SETS ((status), ())
SELECT status FROM orders GROUP BY status HAVING sum(amount) > 100
SELECT count(*) FILTER (WHERE amount > 10) FROM orders
SELECT max_by(id, amount) FROM orders

-- CASE and conditionals.
SELECT CASE WHEN amount > 100 THEN 'big' ELSE status END AS kind FROM orders
SELECT CASE status WHEN 'x' THEN id END FROM orders
SELECT coalesce(c.email, c.name) AS contact, if(o.amount > 0, o.id) FROM orders o JOIN customers c ON o.customer_id = c.id

-- Views, over views, across schemas.
SELECT * FROM big_orders
SELECT id, doubled FROM big_orders WHERE customer_id > 0
SELECT * FROM big_order_names
SELECT b.name FROM big_order_names b JOIN orders o ON b.id = o.id
SELECT * FROM hr.dept_pay
SELECT * FROM hr.sales_customers
SELECT * FROM hr."Vip_Names"
SELECT vipname FROM hr.vip_names

-- session: memory.hr

-- The same names resolve in the session's schema.
SELECT id, name FROM customers
SELECT * FROM customers
SELECT e.name, d.name FROM employees e JOIN depts d ON e.dept_id = d.id
SELECT * FROM dept_pay
SELECT name FROM sales.customers
SELECT c.name, h.vip FROM sales.customers c JOIN customers h ON c.id = h.id
SELECT * FROM sales_customers
SELECT * FROM sales.big_orders

-- session: memory.default

SELECT id FROM sales.orders
SELECT * FROM hr.employees e JOIN hr.depts d ON e.dept_id = d.id
SELECT name FROM memory.hr.customers
SELECT * FROM sales.big_order_names
WITH x AS (SELECT * FROM sales.customers) SELECT x.name, d.name FROM x JOIN hr.depts d ON x.region = d.region
