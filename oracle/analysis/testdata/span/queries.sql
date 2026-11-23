-- One query per line, run in the session (the current schema) the last
-- "-- session:" line names, on the fixture the last "-- fixture:" line names.
-- fixture: shop
-- session: SHOP

-- plain columns, aliases, star
SELECT id, amount FROM orders
SELECT o.id, o.amount AS total FROM orders o
SELECT * FROM customers
SELECT c.* FROM customers c
SELECT customers.name FROM customers
SELECT shop.customers.email FROM shop.customers
SELECT o.*, c.name FROM orders o, customers c WHERE o.customer_id = c.id
SELECT id AS "MixedAlias", name "lower alias" FROM customers
SELECT id, id FROM orders
SELECT orders.id FROM customers orders

-- expressions and functions
SELECT id + 1, amount * 2 AS doubled FROM orders
SELECT name || ' <' || email || '>' AS contact FROM customers
SELECT UPPER(name), NVL(email, 'none') AS email FROM customers
SELECT NVL(email, name) AS reach FROM customers
SELECT CASE WHEN amount > 100 THEN 'big' ELSE status END AS kind FROM orders
SELECT CASE WHEN status = 'paid' THEN amount END AS paid FROM orders
SELECT DECODE(status, 'paid', amount, 0) AS paid FROM orders
SELECT COALESCE(email, region, name) AS c FROM customers
SELECT CAST(amount AS NUMBER(12,2)) AS a, TO_CHAR(created, 'YYYY') AS y FROM orders
SELECT 1 AS one, 'x' AS lit, NULL AS n FROM orders
SELECT 'abc', 'a B', n'x', 1e3 FROM orders
SELECT "Id" + 1, LOWER(plain), "Mixed Case".plain || 'Z', "Full Name" || 'y' FROM "Mixed Case"
SELECT amount /* note */ + 1, -amount, q'[it's]' FROM orders
SELECT id, (SELECT name FROM customers c WHERE c.id = o.customer_id) AS who FROM orders o
SELECT (SELECT MAX(amount) FROM orders) AS top FROM customers

-- Oracle pseudo-columns
SELECT ROWNUM, id FROM orders
SELECT ROWID, name FROM customers
SELECT ROWNUM AS rn, o.* FROM orders o WHERE ROWNUM <= 10
SELECT * FROM (SELECT id, amount FROM orders ORDER BY amount DESC) WHERE ROWNUM <= 3
SELECT ORA_ROWSCN, id FROM orders
SELECT USER, SYSDATE - created AS age FROM orders

-- DUAL
SELECT 1 FROM DUAL
SELECT 1 AS x, 'a' AS y FROM DUAL
SELECT dummy FROM dual
SELECT (SELECT COUNT(*) FROM orders) AS n FROM DUAL
SELECT o.id FROM orders o, dual

-- joins
SELECT o.id, c.name FROM orders o JOIN customers c ON c.id = o.customer_id
SELECT o.id, c.name FROM orders o LEFT JOIN customers c ON c.id = o.customer_id
SELECT o.id, c.name FROM orders o RIGHT OUTER JOIN customers c ON c.id = o.customer_id
SELECT o.id, c.id AS cid FROM orders o FULL JOIN customers c ON c.id = o.customer_id
SELECT o.id, c.name FROM orders o, customers c WHERE c.id = o.customer_id(+)
SELECT id, name, amount FROM customers JOIN orders USING (id)
SELECT * FROM customers JOIN orders USING (id)
SELECT * FROM customers NATURAL JOIN people.customers
SELECT id FROM customers NATURAL JOIN people.customers
SELECT o.id, i.sku FROM orders o CROSS JOIN items i
SELECT o.id, x.total FROM orders o CROSS APPLY (SELECT SUM(qty) AS total FROM items i WHERE i.order_id = o.id) x
SELECT o.id, x.sku FROM orders o OUTER APPLY (SELECT sku FROM items i WHERE i.order_id = o.id) x
SELECT o.id, l.qty FROM orders o, LATERAL (SELECT qty FROM items i WHERE i.order_id = o.id) l
SELECT e.ename, m.ename AS boss FROM emp e LEFT JOIN emp m ON m.empno = e.mgr

-- subqueries
SELECT x.id, x.total FROM (SELECT id, amount AS total FROM orders) x
SELECT total FROM (SELECT amount + 1 AS total FROM orders)
SELECT * FROM (SELECT id, name FROM customers) WHERE id > 1
SELECT id FROM orders WHERE customer_id IN (SELECT id FROM customers WHERE region = 'eu')
SELECT id FROM customers c WHERE EXISTS (SELECT 1 FROM orders o WHERE o.customer_id = c.id)
SELECT name FROM customers WHERE id = (SELECT MAX(customer_id) FROM orders)

-- CTEs
WITH t AS (SELECT id, amount FROM orders) SELECT id, amount FROM t
WITH t (k, v) AS (SELECT id, amount FROM orders) SELECT k, v FROM t
WITH a AS (SELECT id, name FROM customers), b AS (SELECT id FROM a) SELECT * FROM b
WITH orders AS (SELECT id, name FROM customers) SELECT * FROM orders
WITH orders AS (SELECT id FROM customers) SELECT o.amount FROM shop.orders o JOIN orders c ON c.id = o.customer_id
WITH r (n, nm) AS (SELECT empno, ename FROM emp WHERE mgr IS NULL UNION ALL SELECT e.empno, e.ename FROM emp e JOIN r ON e.mgr = r.n) SELECT n, nm FROM r

-- CONNECT BY
SELECT empno, ename, LEVEL FROM emp START WITH mgr IS NULL CONNECT BY PRIOR empno = mgr
SELECT ename, PRIOR ename AS boss FROM emp CONNECT BY PRIOR empno = mgr
SELECT ename, SYS_CONNECT_BY_PATH(ename, '/') AS path, CONNECT_BY_ROOT ename AS root FROM emp START WITH mgr IS NULL CONNECT BY PRIOR empno = mgr
SELECT LEVEL AS n FROM DUAL CONNECT BY LEVEL <= 5

-- set operations
SELECT id, name FROM customers UNION SELECT id, name FROM people.employees
SELECT id FROM customers UNION ALL SELECT id FROM orders
SELECT id FROM customers INTERSECT SELECT customer_id FROM orders
SELECT id FROM customers MINUS SELECT customer_id FROM orders
SELECT name FROM customers UNION SELECT name FROM customers
SELECT id, name FROM customers UNION ALL SELECT id, 'x' FROM orders

-- aggregates and windows
SELECT COUNT(*) FROM orders
SELECT customer_id, SUM(amount) AS total, COUNT(*) AS n FROM orders GROUP BY customer_id
SELECT status, MAX(amount) FROM orders GROUP BY status HAVING COUNT(*) > 1
SELECT MIN(name) AS first FROM customers
SELECT region, LISTAGG(name, ',') WITHIN GROUP (ORDER BY name) AS names FROM customers GROUP BY region
SELECT id, amount, SUM(amount) OVER (PARTITION BY customer_id ORDER BY id) AS running FROM orders
SELECT id, ROW_NUMBER() OVER (ORDER BY amount) AS rn FROM orders
SELECT id, LAG(amount) OVER (ORDER BY id) AS prev FROM orders
SELECT DISTINCT region FROM customers
SELECT region, COUNT(DISTINCT id) FROM customers GROUP BY ROLLUP (region)

-- PIVOT, UNPIVOT, MODEL
SELECT * FROM sales PIVOT (SUM(amount) FOR quarter IN ('Q1' AS q1, 'Q2' AS q2))
SELECT region, q1 FROM sales PIVOT (SUM(amount) FOR quarter IN ('Q1' AS q1, 'Q2' AS q2))
SELECT * FROM sales PIVOT (SUM(amount) AS total FOR quarter IN ('q1', 'Q2' AS b, 3))
SELECT * FROM sales PIVOT (SUM(amount) FOR quarter IN ('q1', 'a b'))
SELECT * FROM quarterly UNPIVOT (amount FOR quarter IN (q1 AS 'Q1', q2 AS 'Q2'))
SELECT region, amount FROM quarterly UNPIVOT (amount FOR quarter IN (q1, q2))
SELECT region, quarter, amount FROM sales MODEL PARTITION BY (region) DIMENSION BY (quarter) MEASURES (amount) RULES (amount['Q9'] = amount['Q1'])

-- views
SELECT * FROM order_totals
SELECT id, gross FROM order_totals
SELECT customer_name FROM shop.order_totals WHERE customer_id = 1
SELECT * FROM big_orders
SELECT b.customer_name, o.status FROM big_orders b JOIN orders o ON o.id = b.id
SELECT * FROM all_people
SELECT * FROM people.staff
SELECT staff_name, dept FROM people.staff WHERE staff_id = 1
SELECT * FROM "Vip Names"
SELECT "VipName" FROM shop."Vip Names"
SELECT * FROM region_counts
SELECT region, n FROM region_counts

-- synonyms and database links
SELECT id, amount FROM syn_orders
SELECT syn_orders.id, s2.status FROM syn_orders, shop.syn_orders s2 WHERE s2.id = syn_orders.id
SELECT name FROM remote_emp
SELECT * FROM staff_syn
SELECT name FROM omni_pub_emp
SELECT id, name FROM employees@omni_loop
SELECT e.name, c.name AS customer FROM people.employees@omni_loop e JOIN customers c ON c.id = e.id

-- quoted and upper-case names, invisible columns
SELECT * FROM "Mixed Case"
SELECT "Id", "Full Name", PLAIN FROM "Mixed Case"
SELECT m."Full Name" AS fn FROM shop."Mixed Case" m
SELECT "id", "val" FROM "lower"
SELECT * FROM "lower"
SELECT ID, NAME FROM SHOP.CUSTOMERS
SELECT "ID", "NAME" FROM "SHOP"."CUSTOMERS"
SELECT * FROM secrets
SELECT id, token FROM secrets

-- refused
SELECT nope FROM orders
SELECT id FROM missing_table
SELECT id FROM orders, customers
SELECT id FROM "orders"

-- session: PEOPLE
SELECT * FROM employees
SELECT id, name FROM customers
SELECT * FROM staff
SELECT * FROM shop.order_totals
SELECT e.name, d.name AS dept FROM employees e JOIN departments d ON d.id = e.dept_id
SELECT name FROM shop.customers UNION SELECT name FROM customers
