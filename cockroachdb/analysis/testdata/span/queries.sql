-- One query per line, run on the fixture and under the search_path the last
-- "-- fixture:" and "-- session:" lines name.
-- fixture: base
-- session: public

-- Plain columns, aliases and expressions.
SELECT id, name FROM customers
SELECT c.id AS customer, c.email FROM customers c
SELECT id + customer_id, upper(status) FROM orders
SELECT amount * 2 AS doubled, status || '!' AS shout FROM orders
SELECT CAST(amount AS text) AS a, amount::int FROM orders
SELECT 1, 'a' AS b
SELECT DISTINCT region FROM customers
SELECT name FROM customers ORDER BY region LIMIT 5
SELECT id FROM orders OFFSET 1 ROWS FETCH FIRST 2 ROWS ONLY
SELECT CASE WHEN amount > 10 THEN status ELSE 'small' END AS s FROM orders
SELECT coalesce(email, name) FROM customers
SELECT id FROM customers WHERE email LIKE '%@example.com'
SELECT public.customers.name FROM public.customers
TABLE customers

-- Star expansion.
SELECT (c).name, (c).* FROM customers c
SELECT row_to_json(c) FROM customers c
SELECT * FROM ONLY customers
SELECT a, b FROM customers AS c(a, b)
SELECT * FROM customers
SELECT c.* FROM customers c
SELECT *, id AS again FROM orders
SELECT o.*, c.name FROM orders o JOIN customers c ON c.id = o.customer_id
SELECT * FROM orders o, customers c WHERE o.customer_id = c.id
SELECT s.* FROM (SELECT id, email FROM customers) s

-- Hidden columns: a table with no primary key has a hidden rowid, which *
-- leaves out and a name reads.
SELECT * FROM items
SELECT rowid, sku FROM items
SELECT i.*, a.* FROM items i JOIN accounts a ON a.id = i.order_id
SELECT * FROM items NATURAL JOIN accounts
SELECT * FROM items JOIN accounts USING (rowid)
SELECT s.* FROM (SELECT * FROM hr.employees) s
TABLE items

-- Joins.
SELECT o.id, c.name FROM orders o JOIN customers c ON c.id = o.customer_id
SELECT o.id, c.name FROM orders o LEFT JOIN customers c ON c.id = o.customer_id
SELECT o.id, c.name FROM orders o RIGHT JOIN customers c ON c.id = o.customer_id
SELECT o.id, c.name FROM orders o FULL JOIN customers c ON c.id = o.customer_id
SELECT o.id, c.name FROM orders o CROSS JOIN customers c
SELECT * FROM orders JOIN items ON items.order_id = orders.id
SELECT * FROM customers JOIN accounts ON accounts.customer_id = customers.id
SELECT a.id, b.id FROM customers a JOIN customers b ON a.region = b.region
SELECT * FROM (orders o JOIN customers c ON c.id = o.customer_id) j
SELECT j.name FROM (orders o JOIN customers c ON c.id = o.customer_id) AS j

-- USING and NATURAL.
SELECT * FROM customers JOIN orders USING (id) JOIN items ON items.order_id = id
SELECT * FROM customers JOIN orders USING (id) JOIN accounts USING (id)
SELECT * FROM customers FULL JOIN orders USING (id) FULL JOIN accounts USING (id)
SELECT * FROM customers JOIN customer_emails USING (id, email)
SELECT * FROM customers NATURAL FULL JOIN customer_emails
SELECT * FROM (SELECT * FROM customers) s JOIN orders USING (id)
SELECT * FROM orders JOIN items ON true JOIN customers USING (id)
SELECT id, name, amount FROM customers JOIN orders USING (id)
SELECT * FROM customers JOIN orders USING (id)
SELECT * FROM customers LEFT JOIN orders USING (id)
SELECT * FROM customers RIGHT JOIN orders USING (id)
SELECT * FROM customers FULL JOIN orders USING (id)
SELECT * FROM customers JOIN accounts USING (id)
SELECT id FROM accounts JOIN customers USING (id)
SELECT customers.id, orders.id FROM customers JOIN orders USING (id)
SELECT * FROM customers NATURAL JOIN customer_emails
SELECT * FROM orders NATURAL JOIN items

-- Subqueries.
SELECT EXISTS (SELECT 1 FROM orders) AS any
SELECT id FROM customers WHERE id = ANY (SELECT customer_id FROM orders)
SELECT orders.name FROM (SELECT * FROM customers) AS orders
SELECT c.id, t.x FROM customers c, LATERAL unnest(ARRAY[c.name, c.email]) AS t(x)
SELECT * FROM generate_series(1, 3) AS g
SELECT g.n, c.name FROM generate_series(1, 3) AS g(n) JOIN customers c ON c.id = g.n
SELECT * FROM ROWS FROM (generate_series(1, 2), unnest(ARRAY['a', 'b'])) AS r(n, s)
SELECT s.id FROM (SELECT id FROM customers UNION SELECT id FROM orders ORDER BY 1 LIMIT 2) s
SELECT (SELECT max(amount) FROM orders) AS top
SELECT name, (SELECT count(*) FROM orders o WHERE o.customer_id = c.id) AS n FROM customers c
SELECT name FROM customers WHERE id IN (SELECT customer_id FROM orders)
SELECT name FROM customers c WHERE EXISTS (SELECT 1 FROM orders o WHERE o.customer_id = c.id)
SELECT id FROM (SELECT id, name FROM customers) AS s
SELECT x FROM (SELECT id AS x FROM customers) s
SELECT a, b FROM (SELECT id, email FROM customers) AS s(a, b)
SELECT t.n FROM (SELECT amount + 1 AS n FROM orders) t
SELECT * FROM (SELECT * FROM (SELECT id, region FROM customers) a) b
SELECT c.name, l.amount FROM customers c, LATERAL (SELECT amount FROM orders o WHERE o.customer_id = c.id) l
SELECT c.name, l.* FROM customers c CROSS JOIN LATERAL (SELECT max(amount) AS m, min(status) FROM orders o WHERE o.customer_id = c.id) l
SELECT c.name, l.amount FROM customers c LEFT JOIN LATERAL (SELECT amount FROM orders o WHERE o.customer_id = c.id LIMIT 1) l ON true

-- CTEs.
WITH x AS (SELECT id, email FROM customers) SELECT * FROM x a JOIN x b USING (id)
WITH "Cte" AS (SELECT id, name FROM customers) SELECT "Cte".name FROM "Cte"
WITH c AS (SELECT id, email FROM customers) SELECT * FROM (WITH c AS (SELECT id, status FROM orders) SELECT * FROM c) inner_q
WITH orders AS (SELECT * FROM customers) SELECT o.name, (SELECT count(*) FROM orders) FROM orders o
WITH a AS (SELECT id, status FROM orders), b AS (SELECT * FROM a UNION ALL SELECT id, name FROM customers) SELECT status FROM b
WITH c AS (SELECT id, email FROM customers) SELECT * FROM c
WITH c AS (SELECT id, email FROM customers) SELECT email FROM c WHERE id > 1
WITH c(k, v) AS (SELECT id, email FROM customers) SELECT k, v FROM c
WITH a AS (SELECT id, name FROM customers), b AS (SELECT a.name FROM a) SELECT * FROM b
WITH t AS (SELECT customer_id, sum(amount) AS s FROM orders GROUP BY customer_id) SELECT c.name, t.s FROM customers c JOIN t ON t.customer_id = c.id
WITH customers AS (SELECT id, status AS name FROM orders) SELECT name FROM customers
WITH c AS MATERIALIZED (SELECT id, email FROM customers) SELECT email FROM c
WITH unused AS (SELECT ssn FROM hr.employees) SELECT id FROM customers
WITH RECURSIVE r(n) AS (SELECT 1 UNION ALL SELECT n + 1 FROM r WHERE n < 5) SELECT n FROM r
WITH RECURSIVE chain AS (SELECT id, manager_id, name FROM hr.employees WHERE manager_id IS NULL UNION ALL SELECT e.id, e.manager_id, e.name FROM hr.employees e JOIN chain ON e.manager_id = chain.id) SELECT * FROM chain
WITH RECURSIVE chain AS (SELECT id, name FROM hr.employees UNION SELECT e.id, e.name FROM hr.employees e JOIN chain c ON c.id = e.manager_id) SELECT name FROM chain

-- Set operations.
SELECT id, name FROM customers UNION SELECT id, status FROM orders
SELECT id, name FROM customers UNION ALL SELECT id, status FROM orders
SELECT email FROM customers INTERSECT SELECT status FROM orders
SELECT region FROM customers EXCEPT SELECT status FROM orders
SELECT id FROM customers UNION SELECT id FROM orders UNION SELECT order_id FROM items
SELECT * FROM (SELECT id, email FROM customers UNION SELECT id, status FROM orders) u
SELECT name FROM customers UNION SELECT 'none' ORDER BY 1
VALUES (1, 'a'), (2, 'b')
SELECT * FROM (VALUES (1, 'a')) AS v(n, s)

-- Views, nested views and views over other schemas.
SELECT * FROM renamed
SELECT mail FROM renamed WHERE cid > 0
SELECT * FROM "Mixed"."Upper View"
SELECT "ID", "Lower" FROM "Mixed"."Upper View" u
SELECT * FROM customer_emails
SELECT email FROM customer_emails WHERE id = 1
SELECT * FROM order_totals
SELECT customer_id, total FROM order_totals
SELECT * FROM vip
SELECT v.email, o.status FROM vip v JOIN orders o ON o.customer_id = v.id
SELECT label FROM labeled
SELECT * FROM staff_orders
SELECT * FROM hr.directory
SELECT d.name, e.ssn FROM hr.directory d JOIN hr.employees e USING (id)
SELECT * FROM region_counts
SELECT r.region, r.n, c.name FROM region_counts r JOIN customers c USING (region)

-- Quoted and mixed-case names.
SELECT "Id", "MixedCol", lower_col FROM "Mixed"."CamelTable"
SELECT * FROM "Mixed"."CamelTable"
SELECT t."MixedCol" AS "Out" FROM "Mixed"."CamelTable" t
SELECT * FROM "Mixed"."CamelView"
SELECT "Renamed" FROM "Mixed"."CamelView"
SELECT ID, NAME FROM CUSTOMERS
SELECT "id", "name" FROM "customers"
SELECT id AS "Weird Name", name AS "x.y" FROM customers

-- Aliases that shadow table names.
SELECT orders.id FROM customers AS orders
SELECT customers.amount FROM orders AS customers
SELECT o.id FROM orders o JOIN customers orders ON orders.id = o.customer_id
SELECT customers.name FROM customers JOIN (SELECT id, status AS name FROM orders) AS customers2 ON customers2.id = customers.id

-- Window functions.
SELECT id, row_number() OVER (ORDER BY amount) FROM orders
SELECT id, sum(amount) OVER (PARTITION BY customer_id) AS running FROM orders
SELECT name, rank() OVER w FROM customers WINDOW w AS (PARTITION BY region ORDER BY name)
SELECT id, lag(status) OVER (ORDER BY created_at) AS prev FROM orders

-- Aggregates.
SELECT DISTINCT ON (region) region, name FROM customers ORDER BY region, name
SELECT max(id) AS id FROM customers
SELECT count(*) FROM orders
SELECT customer_id, sum(amount), max(status) FROM orders GROUP BY customer_id
SELECT region, count(DISTINCT email) FROM customers GROUP BY region HAVING count(*) > 1
SELECT string_agg(name, ',' ORDER BY name) FROM customers
SELECT customer_id FROM orders GROUP BY customer_id
SELECT count(*) FILTER (WHERE status = 'paid') AS paid FROM orders

-- Stars beside other targets.
SELECT *, c.name FROM customers c CROSS JOIN unnest(ARRAY['a']) AS tags
SELECT c.name, upper(c.email) AS e, * FROM customers c CROSS JOIN unnest(ARRAY['a']) AS tags
SELECT id, *, region FROM customers, unnest(ARRAY['a']) AS tags
SELECT *, s.email FROM (SELECT id, upper(email) AS email FROM customers) s CROSS JOIN unnest(ARRAY['a']) AS tags
SELECT *, customer_id FROM orders JOIN items ON items.order_id = orders.id CROSS JOIN unnest(ARRAY['a']) AS tags

-- session: sales,public
SELECT id, region FROM orders
SELECT * FROM orders
SELECT o.total, t.goal FROM orders o JOIN targets t USING (region)
SELECT id, amount FROM public.orders
SELECT name FROM customers

-- session: hr,public
SELECT name, ssn FROM employees
SELECT * FROM directory
SELECT e.name, m.name AS manager FROM employees e LEFT JOIN employees m ON m.id = e.manager_id
SELECT d.name, c.email FROM directory d JOIN customers c ON c.id = d.id
