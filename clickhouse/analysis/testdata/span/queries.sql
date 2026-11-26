-- The span corpus: one query per line, run in the session's database on the
-- fixture. span_test.go reads it; testdata/README.md says how ClickHouse
-- answers each query.
-- fixture: shop
-- session: shop

-- Columns, stars and qualifiers.
SELECT id, name FROM customers
SELECT * FROM customers
SELECT customers.* FROM customers
SELECT c.* FROM customers AS c
SELECT shop.customers.name FROM customers
SELECT c.name, c.email FROM customers AS c
SELECT name AS customer, email AS mail FROM customers
SELECT id FROM shop.orders
SELECT amount, status FROM shop.orders WHERE status = 'paid'
SELECT DISTINCT country FROM customers
SELECT id, name FROM customers ORDER BY name LIMIT 10
SELECT id FROM customers LIMIT 5 OFFSET 2
SELECT name FROM customers WHERE id IN (1, 2, 3)
SELECT * FROM items
SELECT o.* FROM shop.orders o
SELECT `name`, `email` FROM `customers`
SELECT 1
SELECT 1 AS one, 'x' AS letter
SELECT now() AS ts
SELECT count() FROM customers
SELECT count(*) FROM orders

-- Expressions and the names ClickHouse gives them.
SELECT amount + 1 FROM orders
SELECT amount * 2, amount / 3, amount - 1 FROM orders
SELECT id % 2 FROM orders
SELECT -amount FROM orders
SELECT NOT id FROM orders
SELECT id = 1 AND status != 'x' FROM orders
SELECT id = 1 OR id = 2 OR id = 3 FROM orders
SELECT id BETWEEN 1 AND 10 FROM orders
SELECT id NOT BETWEEN 1 AND 10 FROM orders
SELECT CASE WHEN amount > 100 THEN 'big' ELSE 'small' END FROM orders
SELECT CASE status WHEN 'paid' THEN 1 ELSE 0 END FROM orders
SELECT amount > 100 ? 'big' : 'small' FROM orders
SELECT if(amount > 100, id, customer_id) FROM orders
SELECT [id, customer_id] FROM orders
SELECT (id, status) FROM orders
SELECT tags[1] FROM orders
SELECT CAST(id AS String) FROM orders
SELECT id::Int64 FROM orders
SELECT toString(id) || status FROM orders
SELECT concat(name, ' <', email, '>') AS contact FROM customers
SELECT lower(name), upper(country) FROM customers
SELECT length(name) > 3 FROM customers
SELECT name LIKE 'A%' FROM customers
SELECT name NOT LIKE 'A%' FROM customers
SELECT name ILIKE 'a%' FROM customers
SELECT email IS NULL, email IS NOT NULL FROM customers
SELECT id IN (1, 2) FROM customers
SELECT toYear(created), toMonth(created) FROM customers
SELECT EXTRACT(YEAR FROM created) FROM customers
SELECT created + INTERVAL 1 DAY FROM customers
SELECT substring(name, 1, 3) FROM customers
SELECT 1.5, 'it''s', NULL, true FROM customers
SELECT -1, 0x1F, 1e3 FROM customers
SELECT coalesce(email, name) FROM customers
SELECT multiIf(qty > 10, 'many', qty > 1, 'some', 'one') FROM items
SELECT round(price * qty, 2) AS total FROM items
SELECT (price + 1) * qty FROM items
SELECT tuple(id, sku).1 FROM items

-- Aggregates and grouping.
SELECT count(), sum(amount) FROM orders
SELECT COUNT(id) FROM orders
SELECT count(DISTINCT customer_id) FROM orders
SELECT customer_id, sum(amount) FROM orders GROUP BY customer_id
SELECT customer_id, sum(amount) AS total FROM orders GROUP BY customer_id HAVING total > 10
SELECT status, max(amount), min(amount), avg(amount) FROM orders GROUP BY status
SELECT uniqExact(customer_id) FROM orders
SELECT sumIf(amount, status = 'paid') FROM orders
SELECT quantile(0.5)(amount) FROM orders
SELECT groupArray(sku) FROM items
SELECT argMax(sku, price) FROM items
SELECT order_id, sum(price * qty) AS total FROM items GROUP BY order_id ORDER BY total DESC
SELECT country, count() FROM customers GROUP BY country WITH TOTALS
SELECT country, count() FROM customers GROUP BY ROLLUP(country)
SELECT sum(amount) * 2 FROM orders
SELECT status, count() FROM orders GROUP BY ALL

-- Joins.
SELECT c.name, o.amount FROM customers AS c JOIN orders AS o ON o.customer_id = c.id
SELECT c.name, o.amount FROM customers c INNER JOIN orders o ON o.customer_id = c.id WHERE o.amount > 10
SELECT c.name, o.amount FROM customers AS c LEFT JOIN orders AS o ON o.customer_id = c.id
SELECT c.name, o.amount FROM customers AS c RIGHT JOIN orders AS o ON o.customer_id = c.id
SELECT c.name, o.amount FROM customers AS c FULL JOIN orders AS o ON o.customer_id = c.id
SELECT c.name, o.amount FROM customers AS c CROSS JOIN orders AS o
SELECT c.name, o.amount FROM customers AS c, orders AS o WHERE o.customer_id = c.id
SELECT * FROM customers AS c JOIN orders AS o ON o.customer_id = c.id
SELECT c.*, o.amount FROM customers AS c JOIN orders AS o ON o.customer_id = c.id
SELECT id FROM customers AS c JOIN orders AS o ON o.customer_id = c.id
SELECT c.id, o.id FROM customers AS c JOIN orders AS o ON o.customer_id = c.id
SELECT o.id, c.id FROM customers AS c JOIN orders AS o ON o.customer_id = c.id
SELECT name, amount FROM customers JOIN orders ON orders.customer_id = customers.id
SELECT * FROM orders JOIN items ON items.order_id = orders.id
SELECT o.id, i.sku, c.name FROM orders AS o JOIN items AS i ON i.order_id = o.id JOIN customers AS c ON c.id = o.customer_id
SELECT * FROM orders AS o JOIN items AS i ON i.order_id = o.id JOIN customers AS c ON c.id = o.customer_id
SELECT id FROM orders AS o JOIN items AS i ON i.order_id = o.id JOIN customers AS c ON c.id = o.customer_id
SELECT sku FROM orders AS o JOIN items AS i ON i.order_id = o.id JOIN customers AS c ON c.id = o.customer_id
SELECT * FROM orders AS o JOIN (SELECT order_id AS id, sku FROM items) AS i USING (id)
SELECT id, sku FROM orders AS o JOIN (SELECT order_id AS id, sku FROM items) AS i USING id
SELECT a.name, b.name FROM customers AS a JOIN customers AS b ON a.id = b.id
SELECT * FROM customers AS a, customers AS b
SELECT b.id FROM customers AS a JOIN customers AS b USING (id)
SELECT c.name FROM customers AS c LEFT SEMI JOIN orders AS o ON o.customer_id = c.id
SELECT c.name FROM customers AS c LEFT ANTI JOIN orders AS o ON o.customer_id = c.id
SELECT c.name, o.amount FROM customers AS c ANY LEFT JOIN orders AS o ON o.customer_id = c.id
SELECT o.id, i.qty FROM orders AS o ASOF JOIN items AS i ON i.order_id = o.id AND i.qty <= o.id
SELECT c.name, o.amount FROM customers AS c GLOBAL JOIN orders AS o ON o.customer_id = c.id
SELECT c.name, e.salary FROM customers AS c JOIN hr.employees AS e ON e.id = c.id

-- Subqueries in FROM.
SELECT x FROM (SELECT id AS x FROM customers)
SELECT t.x FROM (SELECT id AS x FROM customers) AS t
SELECT * FROM (SELECT id, name FROM customers) AS t
SELECT total FROM (SELECT customer_id, sum(amount) AS total FROM orders GROUP BY customer_id)
SELECT n FROM (SELECT upper(name) AS n FROM customers) WHERE n != ''
SELECT a, b FROM (SELECT id, name FROM customers) AS t(a, b)
SELECT y FROM (SELECT x AS y FROM (SELECT amount AS x FROM orders))
SELECT * FROM (SELECT * FROM (SELECT id, email FROM customers))
SELECT t.amount FROM (SELECT * FROM orders) AS t
SELECT `plus(id, 1)` FROM (SELECT id + 1 FROM customers)
SELECT s.cnt FROM (SELECT count() AS cnt FROM items) AS s
SELECT `a.b` FROM (SELECT 1 AS `a.b`)

-- WITH: CTEs and expression aliases.
WITH c AS (SELECT id, name FROM customers) SELECT name FROM c
WITH c AS (SELECT id, name FROM customers) SELECT * FROM c
WITH c AS (SELECT id, name FROM customers), o AS (SELECT customer_id, amount FROM orders) SELECT c.name, o.amount FROM c JOIN o ON o.customer_id = c.id
WITH big AS (SELECT * FROM orders WHERE amount > 100) SELECT id, amount FROM big
WITH t AS (SELECT sum(amount) AS total FROM orders) SELECT total * 2 FROM t
WITH a AS (SELECT id FROM customers), b AS (SELECT id FROM a) SELECT id FROM b
WITH c AS (SELECT name FROM customers) SELECT x.name FROM c AS x
WITH 10 AS threshold SELECT id FROM orders WHERE amount > threshold
WITH 10 AS threshold SELECT threshold, id FROM orders
WITH amount * 2 AS doubled SELECT doubled FROM orders
WITH (SELECT max(amount) FROM orders) AS top SELECT id, top FROM orders
WITH upper(name) AS n SELECT n, email FROM customers
WITH c AS (SELECT id, name FROM customers) SELECT name FROM c UNION ALL SELECT sku FROM items
WITH customers AS (SELECT order_id AS id FROM items) SELECT id FROM customers
WITH c AS (SELECT id, name FROM customers) SELECT count() FROM c WHERE name != ''
WITH 1 AS x, c AS (SELECT x) SELECT * FROM c
WITH 2 AS k, c AS (SELECT id * k AS n FROM customers) SELECT n FROM c
WITH 1 AS x SELECT * FROM (SELECT x)
WITH 3 AS k SELECT n FROM (SELECT id * k AS n FROM customers)
WITH RECURSIVE r AS (SELECT 1 AS n UNION ALL SELECT n + 1 FROM r WHERE n < 3) SELECT n FROM r
WITH RECURSIVE r AS (SELECT id, customer_id AS parent FROM orders UNION ALL SELECT o.id, r.parent FROM orders AS o JOIN r ON o.customer_id = r.id) SELECT * FROM r

-- Set operations.
SELECT id FROM customers UNION ALL SELECT id FROM orders
SELECT name FROM customers UNION DISTINCT SELECT sku FROM items
SELECT id, name FROM customers UNION ALL SELECT id, status FROM orders
SELECT id FROM customers INTERSECT SELECT customer_id FROM orders
SELECT id FROM customers EXCEPT SELECT customer_id FROM orders
SELECT id FROM customers UNION ALL SELECT id FROM orders UNION ALL SELECT order_id FROM items
SELECT x FROM (SELECT id AS x FROM customers UNION ALL SELECT customer_id FROM orders)
SELECT * FROM (SELECT name FROM customers UNION ALL SELECT status FROM orders) AS u
(SELECT id FROM customers) UNION ALL (SELECT id FROM orders)
SELECT 1 AS n UNION ALL SELECT id FROM customers
SELECT name FROM customers UNION ALL SELECT name FROM hr.employees

-- Subqueries in expressions.
SELECT (SELECT max(amount) FROM orders)
SELECT (SELECT max(amount) FROM orders) AS top, name FROM customers
SELECT name FROM customers WHERE id IN (SELECT customer_id FROM orders)
SELECT name, id IN (SELECT customer_id FROM orders) AS has_orders FROM customers
SELECT name FROM customers WHERE id NOT IN (SELECT customer_id FROM orders WHERE amount > 10)
SELECT name FROM customers WHERE EXISTS (SELECT 1 FROM orders)
SELECT name, NOT EXISTS (SELECT 1 FROM orders WHERE amount > 0) AS none FROM customers
SELECT name FROM customers WHERE id GLOBAL IN (SELECT customer_id FROM orders)
WITH ids AS (SELECT customer_id FROM orders) SELECT name FROM customers WHERE id IN ids
WITH ids AS (SELECT customer_id FROM orders) SELECT name, id IN ids AS has_orders FROM customers
SELECT name FROM customers WHERE id IN nope
SELECT id, (SELECT count() FROM items) AS n FROM orders
SELECT name FROM customers WHERE id = (SELECT min(customer_id) FROM orders)
SELECT id FROM orders WHERE customer_id IN (SELECT id FROM customers WHERE country IN (SELECT country FROM customers WHERE id = 1))

-- Views.
SELECT * FROM big_orders
SELECT id, amount FROM big_orders
SELECT amount * 2 FROM big_orders
SELECT * FROM customer_totals
SELECT name, total FROM customer_totals
SELECT * FROM vip
SELECT name FROM vip
SELECT v.name, o.amount FROM vip AS v JOIN orders AS o ON o.customer_id = v.customer_id
SELECT b.id FROM big_orders AS b JOIN items AS i ON i.order_id = b.id
SELECT count() FROM big_orders
SELECT customer_id FROM big_orders UNION ALL SELECT customer_id FROM customer_totals
WITH v AS (SELECT * FROM vip) SELECT name FROM v
SELECT name FROM (SELECT * FROM customer_totals) WHERE total > 0
SELECT * FROM hr.payroll
SELECT dept, sum(yearly) FROM hr.payroll GROUP BY dept
SELECT id FROM customers WHERE id IN (SELECT customer_id FROM vip)

-- Arrays, ARRAY JOIN and lambdas.
SELECT id, tag FROM orders ARRAY JOIN tags AS tag
SELECT id, tags FROM orders ARRAY JOIN tags
SELECT id, tag FROM orders LEFT ARRAY JOIN tags AS tag
SELECT * FROM orders ARRAY JOIN tags
SELECT * FROM orders ARRAY JOIN tags AS tag
SELECT tag, n FROM orders ARRAY JOIN tags AS tag, arrayEnumerate(tags) AS n
SELECT arrayMap(x -> upper(x), tags) FROM orders
SELECT arrayFilter(x -> x != status, tags) FROM orders
SELECT arrayMap(x -> x || status, tags) AS t FROM orders
SELECT length(tags), has(tags, 'gift') FROM orders
SELECT arrayJoin(tags) AS tag FROM orders
SELECT arrayJoin([1, 2, 3]) AS n
SELECT groupArray(name) AS names FROM customers
SELECT arraySum(x -> x * 2, [qty, order_id]) FROM items

-- Window functions.
SELECT id, row_number() OVER (ORDER BY id) FROM orders
SELECT id, sum(amount) OVER (PARTITION BY customer_id) FROM orders
SELECT id, sum(amount) OVER (PARTITION BY customer_id ORDER BY id) AS running FROM orders
SELECT id, rank() OVER (PARTITION BY status ORDER BY amount DESC) FROM orders
SELECT id, lagInFrame(amount) OVER (ORDER BY id ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) FROM orders
SELECT id, count() OVER w FROM orders WINDOW w AS (PARTITION BY customer_id)
SELECT id FROM orders QUALIFY row_number() OVER (PARTITION BY customer_id ORDER BY amount DESC) = 1

-- Star transformers.
SELECT * EXCEPT (email) FROM customers
SELECT * EXCEPT email FROM customers
SELECT * EXCEPT (email, country) FROM customers
SELECT * REPLACE (upper(name) AS name) FROM customers
SELECT c.* EXCEPT (id) FROM customers AS c
SELECT * APPLY (toString) FROM items
SELECT * APPLY (x -> concat(toString(x), status)) FROM orders
SELECT * EXCEPT (sku) APPLY (max) FROM items
SELECT * EXCEPT ('mail') FROM customers
SELECT untuple((id, email)) FROM customers
SELECT COLUMNS('^c') FROM customers
SELECT concat(COLUMNS('^(name|email)$')) FROM customers

-- Aliases resolve before columns, anywhere in the query.
SELECT amount * 2 AS doubled, doubled + 1 FROM orders
SELECT id + 1 AS id FROM orders
SELECT name AS id, id FROM customers
SELECT email AS e FROM customers WHERE e != ''
SELECT upper(name) AS n FROM customers ORDER BY n
SELECT country AS c, count() FROM customers GROUP BY c
SELECT (amount + 1 AS a) * 2, a FROM orders
SELECT id AS x, x * 2 AS y, y + 1 FROM orders
SELECT customer_id AS amount, amount FROM orders
SELECT x FROM customers WHERE (email AS x) != ''
SELECT toString(id) AS status, concat(status, '') AS s FROM orders SETTINGS prefer_column_name_to_alias = 1

-- ClickHouse's own clauses.
FROM customers SELECT name
FROM customers SELECT *
SELECT name FROM customers PREWHERE id > 1
SELECT order_id, qty FROM items FINAL
SELECT id FROM orders LIMIT 1 BY customer_id
SELECT id FROM orders ORDER BY id WITH FILL
SELECT name FROM customers SETTINGS max_threads = 1
SELECT id FROM orders WHERE _part = 'all_0_0_0'
SELECT _part, id FROM orders
SELECT name FROM customers WHERE id > 0 ORDER BY id LIMIT 3
SELECT DISTINCT ON (country) country, name FROM customers
SELECT TOP 3 name FROM customers

-- Table functions and system tables.
SELECT number FROM numbers(10)
SELECT number * 2 AS n FROM numbers(3)
SELECT * FROM numbers(2)
SELECT c.name, n.number FROM customers AS c, numbers(2) AS n
SELECT name FROM system.tables
SELECT database, name FROM system.tables WHERE database = 'shop'
SELECT name, type FROM system.columns WHERE table = 'customers'
SELECT count() FROM system.databases

-- Another database.
SELECT name, salary FROM hr.employees
SELECT * FROM hr.depts
SELECT e.name, d.title FROM hr.employees AS e JOIN hr.depts AS d ON d.id = e.dept_id
SELECT e.name, m.name FROM hr.employees AS e LEFT JOIN hr.employees AS m ON m.id = e.manager_id
SELECT hr.employees.name FROM hr.employees

-- Refused.
SELECT nope FROM customers
SELECT c.nope FROM customers AS c
SELECT * FROM nope
SELECT * FROM nodb.customers
SELECT id FROM customers WHERE nope = 1
SELECT x.id FROM customers AS c
SELECT * FROM customers WHERE
SELECT id, FROM customers
SELECT name FROM customers WHERE id IN (SELECT nope FROM orders)
SELECT name FROM hr.customers

-- session: hr
SELECT name FROM employees
SELECT * FROM employees
SELECT e.name, d.title FROM employees AS e JOIN depts AS d ON d.id = e.dept_id
SELECT name, yearly FROM payroll
SELECT name FROM shop.customers
SELECT e.name, c.email FROM employees AS e JOIN shop.customers AS c ON c.name = e.name
SELECT dept_id, avg(salary) FROM employees GROUP BY dept_id
WITH d AS (SELECT id, title FROM depts) SELECT e.name, d.title FROM employees AS e JOIN d ON d.id = e.dept_id
SELECT * FROM shop.vip
SELECT name FROM customers

-- fixture: nested
-- session: lab
SELECT attrs.k FROM events
SELECT events.attrs.v FROM events
SELECT * FROM events
SELECT id, k FROM events ARRAY JOIN attrs.k AS k
SELECT id, attrs.k FROM events ARRAY JOIN attrs
