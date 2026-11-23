-- The span.tsv corpus: one query per line, run with the database the last
-- "-- session:" line names selected, on the fixture the last "-- fixture:"
-- line names. See ../README.md.

-- fixture: base
-- session: shop

-- Columns, aliases and expressions.
SELECT id, amount FROM orders
SELECT o.id AS order_id, o.amount FROM orders o
SELECT id + customer_id, upper(status) FROM orders
SELECT CAST(id AS VARCHAR) AS a, concat(name, '@', region) AS b FROM customers
SELECT 1, 'a' AS b
SELECT count(*) FROM orders
SELECT DISTINCT region FROM customers
SELECT name FROM customers ORDER BY region LIMIT 2
SELECT CASE WHEN amount > 5 THEN status ELSE 'small' END AS label FROM orders
SELECT coalesce(email, name) AS contact FROM customers
SELECT amount AS id, id AS amount FROM orders
SELECT id   +   1, -amount, id IN (1, 2), NOT id FROM orders
SELECT 'x', 2.50, NULL, TRUE, "y"
SELECT (id), ((amount)) FROM orders
VALUES (1, 2 + 3), (4, 5)

-- Stars.
SELECT * FROM orders
SELECT o.* FROM orders o
SELECT *, amount FROM orders
SELECT c.*, o.amount FROM customers c JOIN orders o ON o.customer_id = c.id
SELECT * FROM (SELECT * FROM customers) c
SELECT * EXCEPT (amount, status) FROM orders

-- Qualified names, other databases, catalogs and case.
SELECT id FROM shop.orders
SELECT shop.orders.amount, orders.status FROM shop.orders
SELECT s.name, s.salary FROM hr.staff s
SELECT o.id, s.name FROM orders o JOIN hr.staff s ON s.id = o.customer_id
SELECT id, amount FROM internal.shop.orders
SELECT o.id FROM internal.shop.orders o JOIN internal.hr.staff s ON s.id = o.customer_id
SELECT internal.shop.orders.amount FROM internal.shop.orders
SELECT * FROM internal.hr.rich
SELECT shop.o.id, internal.shop.o.amount FROM orders o
SELECT ID, Amount FROM orders
SELECT `id`, `AMOUNT` FROM `orders`
SELECT O.ID FROM orders O
SELECT `Id`, `CamelCol`, lower_col FROM `MixedCase`
SELECT * FROM `MixedCase`
SELECT * FROM mixedcase
SELECT m.camelcol, m.LOWER_COL FROM `MixedCase` m
SELECT id FROM Orders

-- Aliases that shadow table names.
SELECT orders.name FROM customers AS orders
SELECT customers.amount FROM orders customers
SELECT o.id FROM customers o JOIN orders c ON c.customer_id = o.id
SELECT * FROM (SELECT amount AS id FROM orders) orders

-- Joins.
SELECT o.id, c.name FROM orders o JOIN customers c ON c.id = o.customer_id
SELECT o.id, c.name FROM orders o LEFT JOIN customers c ON c.id = o.customer_id
SELECT o.id, c.name FROM customers c RIGHT JOIN orders o ON c.id = o.customer_id
SELECT o.id, c.name FROM orders o FULL OUTER JOIN customers c ON c.id = o.customer_id
SELECT o.id, c.name FROM orders o CROSS JOIN customers c
SELECT o.id, c.name FROM orders o, customers c WHERE c.id = o.customer_id
SELECT o.id FROM orders o LEFT SEMI JOIN customers c ON c.id = o.customer_id
SELECT o.id FROM orders o LEFT ANTI JOIN customers c ON c.id = o.customer_id
SELECT o.id FROM orders o JOIN customers c USING (id)
SELECT id, name FROM orders JOIN customers USING (id)
SELECT * FROM orders o JOIN customers c USING (id)
SELECT * FROM orders o LEFT JOIN customers c USING (id)
SELECT * FROM orders NATURAL JOIN customers
SELECT i.sku, o.status, c.region FROM items i JOIN orders o ON o.id = i.order_id JOIN customers c ON c.id = o.customer_id
SELECT a.id, b.id FROM orders a JOIN orders b ON b.customer_id = a.customer_id
SELECT o.id FROM orders o JOIN [shuffle] customers c ON c.id = o.customer_id

-- Subqueries.
SELECT id, (SELECT name FROM customers c WHERE c.id = o.customer_id) AS who FROM orders o
SELECT id, (SELECT max(amount) FROM orders) AS top FROM customers
SELECT id FROM orders WHERE customer_id IN (SELECT id FROM customers WHERE region = 'east')
SELECT id FROM orders o WHERE EXISTS (SELECT 1 FROM items i WHERE i.order_id = o.id)
SELECT d.k, d.v FROM (SELECT id AS k, amount AS v FROM orders) d
SELECT d.k FROM (SELECT id AS k FROM orders LIMIT 2) d
SELECT x.total FROM (SELECT customer_id, sum(amount) AS total FROM orders GROUP BY customer_id) x
SELECT y.k FROM (SELECT x.k FROM (SELECT id AS k FROM customers) x) y
SELECT d.a FROM (SELECT 1 AS a, 2) d
SELECT o.id, e.x FROM orders o LATERAL VIEW explode(array(1, 2)) e AS x
SELECT o.id, l.n FROM orders o, LATERAL (SELECT count(*) AS n FROM items i WHERE i.order_id = o.id) l

-- CTEs.
WITH c AS (SELECT id, name FROM customers) SELECT id, name FROM c
WITH c (k, v) AS (SELECT id, amount FROM orders) SELECT k, v FROM c
WITH a AS (SELECT id, customer_id FROM orders), b AS (SELECT a.id, c.name FROM a JOIN customers c ON c.id = a.customer_id) SELECT * FROM b
WITH c AS (SELECT id FROM customers) SELECT x.id, y.id FROM c x JOIN c y ON y.id = x.id
WITH orders AS (SELECT id, name FROM customers) SELECT id, name FROM orders
WITH c AS (SELECT id, amount FROM orders) SELECT id FROM c WHERE amount > (SELECT avg(amount) FROM c)
WITH RECURSIVE r (n) AS (SELECT 1 UNION ALL SELECT n + 1 FROM r WHERE n < 3) SELECT n FROM r

-- Set operations.
SELECT id FROM orders UNION SELECT id FROM customers
SELECT id, status FROM orders UNION ALL SELECT id, name FROM customers
SELECT name FROM customers INTERSECT SELECT status FROM orders
SELECT id FROM orders EXCEPT SELECT customer_id FROM orders
SELECT id FROM orders MINUS SELECT customer_id FROM orders
SELECT id AS k FROM orders UNION SELECT 1 ORDER BY k
(SELECT id FROM orders LIMIT 1) UNION (SELECT id FROM customers LIMIT 1)
SELECT u.id FROM (SELECT id FROM orders UNION SELECT id FROM customers) u

-- Views and materialized views.
SELECT id, doubled FROM big_orders
SELECT b.status, b.doubled FROM big_orders b
SELECT * FROM big_orders
SELECT customer_id, total, n FROM customer_totals
SELECT * FROM order_names
SELECT o.name, c.region FROM order_names o JOIN customers c ON c.name = o.name
SELECT a, b FROM renamed
SELECT id FROM all_ids
SELECT * FROM staff_pay
SELECT name FROM hr.rich
SELECT VipName FROM hr.`Rich_Names`
SELECT b.id, t.total FROM big_orders b JOIN customer_totals t ON t.customer_id = b.customer_id
SELECT d.doubled FROM (SELECT doubled FROM big_orders) d
WITH v AS (SELECT * FROM big_orders) SELECT doubled FROM v
SELECT customer_id, total FROM order_sums
SELECT * FROM internal.shop.order_sums

-- Window functions.
SELECT id, row_number() OVER (ORDER BY amount) AS rn FROM orders
SELECT id, sum(amount) OVER (PARTITION BY customer_id) AS running FROM orders
SELECT id, lag(amount) OVER (ORDER BY id) AS prev FROM orders
SELECT id, rank() OVER (PARTITION BY status ORDER BY amount) FROM orders
SELECT id FROM orders QUALIFY row_number() OVER (PARTITION BY customer_id ORDER BY amount) = 1

-- Aggregates.
SELECT customer_id, sum(amount), count(*) FROM orders GROUP BY customer_id
SELECT status, avg(amount) AS a FROM orders GROUP BY status HAVING count(*) > 1
SELECT region, group_concat(name ORDER BY name) AS names FROM customers GROUP BY region
SELECT max(amount) FROM orders WHERE status = 'open'
SELECT count(DISTINCT customer_id) AS buyers FROM orders
SELECT c.region, sum(o.amount) FROM customers c JOIN orders o ON o.customer_id = c.id GROUP BY c.region
SELECT customer_id, status, sum(amount) FROM orders GROUP BY ROLLUP (customer_id, status)
SELECT count(*) FROM orders WHERE 1 = 0

-- Table functions and system tables.
SELECT number FROM numbers("number" = "3")
SELECT TABLE_NAME FROM information_schema.tables WHERE TABLE_SCHEMA = 'shop'

-- Refused.
SELECT nope FROM orders
SELECT id FROM no_such_table
SELECT id FROM orders JOIN customers ON customers.id = orders.customer_id
SELECT id FROM no_such_db.orders

-- fixture: base
-- session: hr

SELECT name, salary FROM staff
SELECT s.name, d.name AS dept FROM staff s JOIN depts d ON d.id = s.dept_id
SELECT o.amount FROM shop.orders o
SELECT * FROM rich
SELECT * FROM shop.staff_pay
SELECT name FROM staff UNION SELECT name FROM shop.customers
SELECT id FROM internal.shop.big_orders
