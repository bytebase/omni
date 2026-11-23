-- One query per line, run in the session (database.schema) and on the
-- fixture the last "-- session:" and "-- fixture:" lines name.
-- fixture: shop
-- session: shop.dbo

-- Plain columns, aliases and expressions.
SELECT id, amount FROM orders
SELECT o.id AS order_id, o.amount FROM orders o
SELECT id + customer_id, UPPER(status) FROM orders
SELECT CAST(amount AS varchar(20)) AS a, TRY_CAST(status AS int) AS b FROM orders
SELECT 1 AS one, 'a' AS b
SELECT x = id, y = amount FROM orders
SELECT amount AS id, id AS amount FROM orders
SELECT 'lit' AS id FROM orders
SELECT NULL AS n, GETDATE() AS now
SELECT CASE WHEN region = 'vip' THEN name ELSE email END AS contact FROM customers
SELECT COALESCE(name, email) AS c, ISNULL(region, 'none') AS r FROM customers
SELECT IIF(amount > 10, status, NULL) AS s FROM orders
SELECT CONCAT(name, '-', email) AS ne, name + email AS ne2 FROM customers
SELECT id AS [select], name AS "from" FROM customers
SELECT DISTINCT region FROM customers
SELECT TOP 5 id, amount FROM orders ORDER BY amount DESC
SELECT TOP (10) PERCENT name FROM customers ORDER BY name
SELECT id FROM orders ORDER BY id OFFSET 5 ROWS FETCH NEXT 5 ROWS ONLY
SELECT id FROM orders WITH (NOLOCK)
SELECT id FROM dbo.orders
SELECT id FROM shop.dbo.orders
SELECT id FROM shop..orders
SELECT orders.id, dbo.orders.amount FROM dbo.orders

-- Star expansion.
SELECT * FROM customers
SELECT c.* FROM customers c
SELECT *, id AS again FROM customers
SELECT o.*, c.name FROM orders o JOIN customers c ON c.id = o.customer_id
SELECT * FROM orders o JOIN customers c ON c.id = o.customer_id
SELECT o.*, c.* FROM orders o, customers c WHERE c.id = o.customer_id
SELECT * FROM hr.employees
SELECT TOP 3 * FROM orders ORDER BY amount DESC

-- Joins. T-SQL has no USING or NATURAL join.
SELECT o.id, c.name FROM orders o INNER JOIN customers c ON o.customer_id = c.id
SELECT o.id, c.name FROM orders o LEFT JOIN customers c ON o.customer_id = c.id
SELECT o.id, c.name FROM orders o RIGHT OUTER JOIN customers c ON o.customer_id = c.id
SELECT o.id, c.name FROM orders o FULL JOIN customers c ON o.customer_id = c.id
SELECT o.id, p.name FROM orders o CROSS JOIN products p
SELECT e.name, m.name AS manager FROM hr.employees e LEFT JOIN hr.employees m ON m.id = e.manager_id
SELECT c.name, p.name AS product, i.qty FROM customers c JOIN orders o ON o.customer_id = c.id JOIN order_items i ON i.order_id = o.id JOIN products p ON p.id = i.product_id
SELECT o.id, c.name FROM orders o JOIN (customers c JOIN hr.customers h ON h.id = c.id) ON c.id = o.customer_id
SELECT c.name, x.amount FROM customers c CROSS APPLY (SELECT TOP 1 o.amount FROM orders o WHERE o.customer_id = c.id ORDER BY o.amount DESC) x
SELECT c.name, x.n FROM customers c OUTER APPLY (SELECT COUNT(*) AS n FROM orders o WHERE o.customer_id = c.id) x

-- Subqueries.
SELECT id, (SELECT name FROM customers c WHERE c.id = o.customer_id) AS cname FROM orders o
SELECT (SELECT COUNT(*) FROM orders) AS n
SELECT * FROM (SELECT id, amount FROM orders) AS t
SELECT t.x FROM (SELECT id AS x FROM orders) t
SELECT a, b FROM (SELECT id, name FROM customers) AS t(a, b)
SELECT * FROM (SELECT * FROM (SELECT id, name FROM customers) a) b
SELECT name FROM customers WHERE id IN (SELECT customer_id FROM orders)
SELECT name FROM customers c WHERE EXISTS (SELECT 1 FROM orders o WHERE o.customer_id = c.id)
SELECT name FROM customers WHERE id = (SELECT MAX(customer_id) FROM orders)
SELECT id, CASE WHEN EXISTS (SELECT 1 FROM order_items i WHERE i.order_id = o.id) THEN 'y' ELSE 'n' END AS has_items FROM orders o
SELECT v.k FROM (VALUES (1), (2)) AS v(k)
SELECT u.x FROM (SELECT id AS x FROM orders UNION ALL SELECT id FROM customers) u

-- Common table expressions.
WITH c AS (SELECT id, name FROM customers) SELECT * FROM c
WITH c(x, y) AS (SELECT id, name FROM customers) SELECT x, y FROM c
WITH a AS (SELECT id, customer_id FROM orders), b AS (SELECT a.id, c.name FROM a JOIN customers c ON c.id = a.customer_id) SELECT name, id FROM b
WITH c AS (SELECT id FROM orders) SELECT x.id FROM c x JOIN c y ON x.id = y.id
WITH orders AS (SELECT id, name FROM customers) SELECT * FROM orders
WITH r AS (SELECT id, manager_id, name FROM hr.employees WHERE manager_id IS NULL UNION ALL SELECT e.id, e.manager_id, e.name FROM hr.employees e JOIN r ON e.manager_id = r.id) SELECT * FROM r
WITH r(n, label) AS (SELECT 1, CAST(name AS nvarchar(100)) FROM products UNION ALL SELECT n + 1, label FROM r WHERE n < 3) SELECT n, label FROM r

-- Set operations.
SELECT id FROM orders UNION SELECT id FROM customers
SELECT id AS k, amount FROM orders UNION ALL SELECT order_id, price FROM order_items
SELECT name FROM customers INTERSECT SELECT name FROM hr.employees
SELECT name FROM customers EXCEPT SELECT name FROM products
SELECT name FROM customers UNION SELECT name FROM hr.employees ORDER BY name
SELECT id, name FROM customers UNION SELECT id, name FROM products UNION SELECT id, name FROM hr.departments

-- Views, nested views, views over other schemas and databases.
SELECT * FROM order_totals
SELECT id, gross FROM order_totals
SELECT customer_name FROM dbo.order_totals WHERE customer_id = 1
SELECT * FROM big_orders
SELECT b.customer_name, o.status FROM big_orders b JOIN orders o ON o.id = b.id
SELECT * FROM hr.staff
SELECT staff_name, dept FROM hr.staff WHERE staff_id = 1
SELECT * FROM all_people
SELECT * FROM v_products
SELECT * FROM archived
SELECT customer_name FROM order_totals WHERE id IN (SELECT id FROM big_orders)
SELECT * FROM [Vip Names]
SELECT VipName FROM dbo.[Vip Names]
SELECT t.id, s.staff_name FROM order_totals t JOIN hr.staff s ON s.staff_id = t.customer_id

-- Other databases.
SELECT * FROM archive.dbo.old_orders
SELECT * FROM archive.dbo.recent
SELECT c.name, a.amount FROM customers c JOIN archive.dbo.old_orders a ON a.customer_id = c.id

-- Quoted and mixed-case names.
SELECT [Order ID], [Unit Price] * Qty AS total FROM [Order Details]
SELECT "Order ID", Qty FROM "Order Details"
SELECT d.[Order ID] FROM dbo.[Order Details] AS d
SELECT * FROM [Order Details]
SELECT id, camelname FROM mixedcase
SELECT ID, NAME FROM CUSTOMERS
SELECT [Id], [CamelName] FROM [dbo].[MixedCase]
SELECT c.NAME AS [Customer Name] FROM Customers AS C
SELECT * FROM MIXEDCASE

-- Aliases that shadow table names.
SELECT customers.id FROM orders AS customers
SELECT orders.name FROM customers AS orders
SELECT customers.nickname FROM hr.customers
SELECT c.id, d.name FROM hr.customers c JOIN dbo.customers d ON d.id = c.id
SELECT customers.name FROM hr.customers AS h JOIN customers ON customers.id = h.id

-- Window functions.
SELECT id, ROW_NUMBER() OVER (ORDER BY amount) AS rn FROM orders
SELECT id, SUM(amount) OVER (PARTITION BY customer_id) AS total FROM orders
SELECT id, LAG(amount) OVER (ORDER BY created) AS prev FROM orders
SELECT name, RANK() OVER (PARTITION BY dept_id ORDER BY salary DESC) AS r FROM hr.employees
SELECT id, FIRST_VALUE(status) OVER (PARTITION BY customer_id ORDER BY created) AS first_status FROM orders

-- Aggregates.
SELECT COUNT(*) FROM orders
SELECT customer_id, SUM(amount) AS total FROM orders GROUP BY customer_id
SELECT customer_id, COUNT(*) FROM orders GROUP BY customer_id HAVING COUNT(*) > 1
SELECT MAX(name) AS m FROM customers
SELECT region, COUNT(DISTINCT id) AS n FROM customers GROUP BY region
SELECT STRING_AGG(name, ',') AS names FROM customers
SELECT region, SUM(LEN(name)) AS l FROM customers GROUP BY ROLLUP (region)

-- Table-valued sources and other forms.
SELECT * FROM (SELECT customer_id, status, amount FROM orders) s PIVOT (SUM(amount) FOR status IN ([new], [paid])) p
SELECT value FROM STRING_SPLIT('a,b', ',')
SELECT id FROM orders FOR JSON PATH

-- Queries SQL Server refuses.
SELECT nope FROM orders
SELECT * FROM missing_table
SELECT id FROM orders o JOIN customers c ON c.id = o.customer_id

-- session: archive.dbo
SELECT * FROM old_orders
SELECT id, amount FROM recent
SELECT c.name FROM shop.dbo.customers c JOIN old_orders o ON o.customer_id = c.id
SELECT * FROM shop.dbo.order_totals
